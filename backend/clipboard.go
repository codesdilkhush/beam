package backend

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	_ "modernc.org/sqlite"
)

type ClipboardItem struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Preview   string    `json:"preview"`
	Timestamp time.Time `json:"timestamp"`
	CharCount int       `json:"charCount"`
}

type ClipboardManager struct {
	mu      sync.RWMutex
	db      *sql.DB
	items   []ClipboardItem
	maxSize int
}

var (
	clipboardInstance *ClipboardManager
	clipboardOnce     sync.Once
)

func GetClipboardManager() *ClipboardManager {
	clipboardOnce.Do(func() {
		dir, err := os.UserHomeDir()
		if err != nil {
			dir = "."
		}
		dataDir := filepath.Join(dir, ".local/share/beam")
		if mkErr := os.MkdirAll(dataDir, 0755); mkErr != nil {
			dataDir = "."
		}

		dbPath := filepath.Join(dataDir, "beam.db")
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			db = nil
		}

		cm := &ClipboardManager{
			db:      db,
			items:   make([]ClipboardItem, 0, 100),
			maxSize: 500,
		}

		if cm.db != nil {
			cm.initDatabase()
			cm.loadHistoryFromDB()
		}

		clipboardInstance = cm
	})
	return clipboardInstance
}

func (cm *ClipboardManager) initDatabase() {
	if cm.db == nil {
		return
	}
	cm.db.SetMaxOpenConns(1)

	query := `
	PRAGMA busy_timeout = 5000;
	CREATE TABLE IF NOT EXISTS clipboard_history (
		id TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		preview TEXT NOT NULL,
		char_count INTEGER NOT NULL,
		created_at DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_clipboard_created ON clipboard_history(created_at DESC);
	`
	if _, err := cm.db.Exec(query); err != nil {
		return
	}
	// Clean up any legacy corrupt binary or PNG rows
	_, _ = cm.db.Exec(`DELETE FROM clipboard_history WHERE content LIKE '%PNG%' OR char_count > 500000;`)
}

func isTextContent(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	if bytes.IndexByte(data, 0) != -1 {
		return false
	}
	return utf8.Valid(data)
}

func (cm *ClipboardManager) loadHistoryFromDB() {
	if cm.db == nil {
		return
	}

	rows, err := cm.db.Query("SELECT id, content, preview, char_count, created_at FROM clipboard_history ORDER BY created_at DESC LIMIT 100")
	if err != nil {
		return
	}
	defer rows.Close()

	var loaded []ClipboardItem
	var corruptIDs []string
	for rows.Next() {
		var item ClipboardItem
		if scanErr := rows.Scan(&item.ID, &item.Content, &item.Preview, &item.CharCount, &item.Timestamp); scanErr == nil {
			if !isTextContent([]byte(item.Content)) || !utf8.ValidString(item.Preview) {
				corruptIDs = append(corruptIDs, item.ID)
				continue
			}
			loaded = append(loaded, item)
		}
	}
	if err := rows.Err(); err != nil {
		return
	}

	if len(corruptIDs) > 0 {
		go func(ids []string) {
			for _, id := range ids {
				_, _ = cm.db.Exec("DELETE FROM clipboard_history WHERE id = ?", id)
			}
		}(corruptIDs)
	}

	cm.mu.Lock()
	cm.items = loaded
	cm.mu.Unlock()
}

func (cm *ClipboardManager) StartWatcher(ctx context.Context) {
	// Background polling of wl-paste on Wayland without window focus is intentionally disabled:
	// it creates compositor race conditions, destroys the active selection owner ("on/off" bug),
	// and floods GNOME Mutter with stack assertion failures.
	// Clipboard capture is handled on-demand when Beam is summoned or history queried.
}

func (cm *ClipboardManager) AddItem(content string) {
	if !isTextContent([]byte(content)) {
		return
	}

	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return
	}

	// Guard against enormous memory exhaustion
	const maxClipLength = 200000
	if len(content) > maxClipLength {
		content = content[:maxClipLength]
		trimmed = strings.TrimSpace(content)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	preview := trimmed
	if len(preview) > 120 {
		preview = preview[:120] + "..."
	}

	item := ClipboardItem{
		ID:        hash,
		Content:   content,
		Preview:   preview,
		Timestamp: time.Now(),
		CharCount: len(content),
	}

	cm.mu.Lock()
	// Move duplicate to top if existing
	newItems := make([]ClipboardItem, 0, len(cm.items)+1)
	newItems = append(newItems, item)
	for _, existing := range cm.items {
		if existing.ID != hash {
			newItems = append(newItems, existing)
		}
	}
	if len(newItems) > cm.maxSize {
		newItems = newItems[:cm.maxSize]
	}
	cm.items = newItems
	cm.mu.Unlock()

	// Persist to SQLite
	if cm.db != nil {
		go func(it ClipboardItem) {
			stmt := `
			INSERT INTO clipboard_history (id, content, preview, char_count, created_at)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET created_at = excluded.created_at;
			`
			if _, insertErr := cm.db.Exec(stmt, it.ID, it.Content, it.Preview, it.CharCount, it.Timestamp); insertErr != nil {
				return
			}

			// Prune excess
			prune := `DELETE FROM clipboard_history WHERE id NOT IN (SELECT id FROM clipboard_history ORDER BY created_at DESC LIMIT 500);`
			if _, pruneErr := cm.db.Exec(prune); pruneErr != nil {
				return
			}
		}(item)
	}
}

func (cm *ClipboardManager) GetHistory() []ClipboardItem {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	res := make([]ClipboardItem, len(cm.items))
	copy(res, cm.items)
	return res
}

func (cm *ClipboardManager) ClearHistory() error {
	cm.mu.Lock()
	cm.items = make([]ClipboardItem, 0, 100)
	cm.mu.Unlock()

	if cm.db != nil {
		_, err := cm.db.Exec("DELETE FROM clipboard_history")
		return err
	}
	return nil
}

func ReadCurrentClipboard() (string, error) {
	if path, err := exec.LookPath("wl-paste"); err == nil {
		// Explicitly request text/plain to avoid receiving binary image formats like PNG
		cmd := exec.Command(path, "-t", "text/plain", "--no-newline")
		var out bytes.Buffer
		cmd.Stdout = &out
		if runErr := cmd.Run(); runErr == nil {
			if isTextContent(out.Bytes()) {
				return out.String(), nil
			}
		}

		// Fallback without explicit MIME filter for applications exposing standard UTF-8 text
		cmd = exec.Command(path, "--no-newline")
		out.Reset()
		if runErr := cmd.Run(); runErr == nil {
			if isTextContent(out.Bytes()) {
				return out.String(), nil
			}
		}
	}

	if path, err := exec.LookPath("xclip"); err == nil {
		// Prefer UTF8_STRING to avoid raw binary targets
		cmd := exec.Command(path, "-selection", "clipboard", "-t", "UTF8_STRING", "-o")
		var out bytes.Buffer
		cmd.Stdout = &out
		if runErr := cmd.Run(); runErr == nil {
			if isTextContent(out.Bytes()) {
				return out.String(), nil
			}
		}

		// Fallback for xclip
		cmd = exec.Command(path, "-selection", "clipboard", "-o")
		out.Reset()
		cmd.Stdout = &out
		if runErr := cmd.Run(); runErr == nil {
			if isTextContent(out.Bytes()) {
				return out.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no clipboard tool found or non-text content")
}

func WriteClipboard(text string) error {
	if path, err := exec.LookPath("wl-copy"); err == nil {
		cmd := exec.Command(path)
		cmd.Stdin = strings.NewReader(text)
		if runErr := cmd.Run(); runErr == nil {
			GetClipboardManager().AddItem(text)
			return nil
		}
	}

	if path, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command(path, "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		if runErr := cmd.Run(); runErr == nil {
			GetClipboardManager().AddItem(text)
			return nil
		}
	}

	return fmt.Errorf("failed to copy to clipboard")
}
