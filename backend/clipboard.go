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

	query := `
	CREATE TABLE IF NOT EXISTS clipboard_history (
		id TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		preview TEXT NOT NULL,
		char_count INTEGER NOT NULL,
		created_at DATETIME NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_clipboard_created ON clipboard_history(created_at DESC);
	`
	_, _ = cm.db.Exec(query)
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
	for rows.Next() {
		var item ClipboardItem
		if scanErr := rows.Scan(&item.ID, &item.Content, &item.Preview, &item.CharCount, &item.Timestamp); scanErr == nil {
			loaded = append(loaded, item)
		}
	}

	cm.mu.Lock()
	cm.items = loaded
	cm.mu.Unlock()
}

func (cm *ClipboardManager) StartWatcher(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		var lastContent string

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				current, err := ReadCurrentClipboard()
				if err != nil || current == "" || current == lastContent {
					continue
				}
				lastContent = current
				cm.AddItem(current)
			}
		}
	}()
}

func (cm *ClipboardManager) AddItem(content string) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return
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
			_, _ = cm.db.Exec(stmt, it.ID, it.Content, it.Preview, it.CharCount, it.Timestamp)

			// Prune excess
			prune := `DELETE FROM clipboard_history WHERE id NOT IN (SELECT id FROM clipboard_history ORDER BY created_at DESC LIMIT 500);`
			_, _ = cm.db.Exec(prune)
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
		cmd := exec.Command(path, "--no-newline")
		var out bytes.Buffer
		cmd.Stdout = &out
		if runErr := cmd.Run(); runErr == nil {
			return out.String(), nil
		}
	}

	if path, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command(path, "-selection", "clipboard", "-o")
		var out bytes.Buffer
		cmd.Stdout = &out
		if runErr := cmd.Run(); runErr == nil {
			return out.String(), nil
		}
	}

	return "", fmt.Errorf("no clipboard tool found")
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
