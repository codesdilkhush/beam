package backend

import (
	"testing"
)

func TestClipboardManager_AddItem(t *testing.T) {
	cm := &ClipboardManager{
		items:   make([]ClipboardItem, 0, 50),
		hashes:  make(map[string]bool),
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
		hashes:  make(map[string]bool),
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
