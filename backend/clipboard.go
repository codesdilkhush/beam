package backend

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
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
	items   []ClipboardItem
	hashes  map[string]bool
	maxSize int
}

var (
	clipboardInstance *ClipboardManager
	clipboardOnce     sync.Once
)

func GetClipboardManager() *ClipboardManager {
	clipboardOnce.Do(func() {
		clipboardInstance = &ClipboardManager{
			items:   make([]ClipboardItem, 0, 50),
			hashes:  make(map[string]bool),
			maxSize: 50,
		}
	})
	return clipboardInstance
}

func (cm *ClipboardManager) StartWatcher(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(600 * time.Millisecond)
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

	cm.mu.Lock()
	defer cm.mu.Unlock()

	// If already in list, move to top
	newItems := make([]ClipboardItem, 0, cm.maxSize)
	for _, item := range cm.items {
		if item.ID != hash {
			newItems = append(newItems, item)
		}
	}

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

	// Prepend to top
	cm.items = append([]ClipboardItem{item}, newItems...)
	if len(cm.items) > cm.maxSize {
		cm.items = cm.items[:cm.maxSize]
	}
}

func (cm *ClipboardManager) GetHistory() []ClipboardItem {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	res := make([]ClipboardItem, len(cm.items))
	copy(res, cm.items)
	return res
}

func ReadCurrentClipboard() (string, error) {
	// Try wl-paste first (Wayland)
	if path, err := exec.LookPath("wl-paste"); err == nil {
		cmd := exec.Command(path, "--no-newline")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err == nil {
			return out.String(), nil
		}
	}

	// Fallback to xclip (X11)
	if path, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command(path, "-selection", "clipboard", "-o")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err == nil {
			return out.String(), nil
		}
	}

	return "", fmt.Errorf("no clipboard tool found")
}

func WriteClipboard(text string) error {
	// Try wl-copy first
	if path, err := exec.LookPath("wl-copy"); err == nil {
		cmd := exec.Command(path)
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			// Also add to manager
			GetClipboardManager().AddItem(text)
			return nil
		}
	}

	// Fallback to xclip
	if path, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command(path, "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			GetClipboardManager().AddItem(text)
			return nil
		}
	}

	return fmt.Errorf("failed to copy to clipboard")
}
