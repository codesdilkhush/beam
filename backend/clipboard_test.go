package backend

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestClipboardManager_AddItem(t *testing.T) {
	cm := &ClipboardManager{
		items:   make([]ClipboardItem, 0, 50),
		maxSize: 10,
	}

	// 1. Adding empty content should be ignored
	cm.AddItem("")
	cm.AddItem("   ")
	if len(cm.GetHistory()) != 0 {
		t.Fatalf("expected 0 items, got %d", len(cm.GetHistory()))
	}

	// 2. Add valid item
	cm.AddItem("Hello World")
	history := cm.GetHistory()
	if len(history) != 1 {
		t.Fatalf("expected 1 item, got %d", len(history))
	}
	if history[0].Content != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", history[0].Content)
	}
	if history[0].CharCount != 11 {
		t.Errorf("expected CharCount 11, got %d", history[0].CharCount)
	}

	// 3. Add another item
	cm.AddItem("Second Item")
	history = cm.GetHistory()
	if len(history) != 2 {
		t.Fatalf("expected 2 items, got %d", len(history))
	}
	if history[0].Content != "Second Item" {
		t.Errorf("expected 'Second Item' on top, got %q", history[0].Content)
	}

	// 4. Duplicate item moves to top without growing count
	cm.AddItem("Hello World")
	history = cm.GetHistory()
	if len(history) != 2 {
		t.Fatalf("expected 2 items after duplicate re-add, got %d", len(history))
	}
	if history[0].Content != "Hello World" {
		t.Errorf("expected 'Hello World' back on top, got %q", history[0].Content)
	}
}

func TestClipboardManager_MaxSize(t *testing.T) {
	cm := &ClipboardManager{
		items:   make([]ClipboardItem, 0, 5),
		maxSize: 3,
	}

	cm.AddItem("one")
	cm.AddItem("two")
	cm.AddItem("three")
	cm.AddItem("four")

	history := cm.GetHistory()
	if len(history) != 3 {
		t.Fatalf("expected max 3 items, got %d", len(history))
	}
	if history[0].Content != "four" {
		t.Errorf("expected newest item 'four' on top, got %q", history[0].Content)
	}
	if history[2].Content != "two" {
		t.Errorf("expected oldest retained item 'two', got %q", history[2].Content)
	}
}

func TestClipboardManager_SQLitePersistence(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_clipboard.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test sqlite db: %v", err)
	}
	defer db.Close()

	cm := &ClipboardManager{
		db:      db,
		items:   make([]ClipboardItem, 0, 10),
		maxSize: 10,
	}
	cm.initDatabase()

	cm.AddItem("Persistent Note 1")
	cm.AddItem("Persistent Note 2")

	// Allow goroutine persist to write
	var count int
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		if err := db.QueryRow("SELECT COUNT(*) FROM clipboard_history").Scan(&count); err == nil && count == 2 {
			break
		}
	}

	if count != 2 {
		t.Errorf("expected 2 rows in sqlite, got %d", count)
	}

	// Test ClearHistory
	if err := cm.ClearHistory(); err != nil {
		t.Fatalf("ClearHistory failed: %v", err)
	}
	if len(cm.GetHistory()) != 0 {
		t.Errorf("expected 0 in-memory items after clear, got %d", len(cm.GetHistory()))
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM clipboard_history").Scan(&count); err != nil {
		t.Fatalf("failed to query count after clear: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 rows in sqlite after clear, got %d", count)
	}
}
