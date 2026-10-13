package backend

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestTodoManager_Operations(t *testing.T) {
	tm := &TodoManager{
		items: make([]TodoItem, 0),
	}

	// 1. Adding empty should error
	if _, err := tm.AddTodo(""); err == nil {
		t.Errorf("expected error adding empty todo, got nil")
	}
	if _, err := tm.AddTodo("   "); err == nil {
		t.Errorf("expected error adding whitespace todo, got nil")
	}

	// 2. Add valid todos
	item1, err := tm.AddTodo("Buy groceries")
	if err != nil {
		t.Fatalf("failed to add todo: %v", err)
	}
	if item1.Title != "Buy groceries" {
		t.Errorf("expected title 'Buy groceries', got %s", item1.Title)
	}
	if item1.Completed {
		t.Errorf("expected new todo to be incomplete")
	}

	item2, err := tm.AddTodo("Review PR")
	if err != nil {
		t.Fatalf("failed to add todo: %v", err)
	}

	todos := tm.GetTodos()
	if len(todos) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(todos))
	}
	// Newest should be on top
	if todos[0].ID != item2.ID {
		t.Errorf("expected newest item first, got %s", todos[0].Title)
	}

	// 3. Toggle todo
	if err := tm.ToggleTodo(item1.ID); err != nil {
		t.Fatalf("failed to toggle todo: %v", err)
	}
	todos = tm.GetTodos()
	var toggledItem *TodoItem
	for _, it := range todos {
		if it.ID == item1.ID {
			toggledItem = &it
			break
		}
	}
	if toggledItem == nil || !toggledItem.Completed {
		t.Errorf("expected item1 to be completed")
	}
	if toggledItem.CompletedAt == nil {
		t.Errorf("expected completedAt to be set")
	}

	// Active should come before completed
	if todos[0].ID != item2.ID {
		t.Errorf("expected active item first, got %s", todos[0].Title)
	}

	// 4. Toggle back to incomplete
	if err := tm.ToggleTodo(item1.ID); err != nil {
		t.Fatalf("failed to toggle back todo: %v", err)
	}
	todos = tm.GetTodos()
	for _, it := range todos {
		if it.ID == item1.ID {
			if it.Completed || it.CompletedAt != nil {
				t.Errorf("expected item1 to be incomplete with nil CompletedAt")
			}
		}
	}

	// 5. Delete todo
	if err := tm.DeleteTodo(item2.ID); err != nil {
		t.Fatalf("failed to delete todo: %v", err)
	}
	if len(tm.GetTodos()) != 1 {
		t.Errorf("expected 1 todo remaining after delete, got %d", len(tm.GetTodos()))
	}
	if err := tm.DeleteTodo("non-existent"); err == nil {
		t.Errorf("expected error deleting non-existent todo")
	}

	// 6. Clear completed
	_, _ = tm.AddTodo("Task to complete")
	todos = tm.GetTodos()
	_ = tm.ToggleTodo(todos[0].ID)
	if err := tm.ClearCompleted(); err != nil {
		t.Fatalf("failed to clear completed: %v", err)
	}
	for _, it := range tm.GetTodos() {
		if it.Completed {
			t.Errorf("found completed item after ClearCompleted: %s", it.Title)
		}
	}
}

func TestTodoManager_SQLitePersistence(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_todo.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test sqlite db: %v", err)
	}
	defer db.Close()

	tm := &TodoManager{
		db:    db,
		items: make([]TodoItem, 0),
	}
	tm.initDatabase()

	it1, err := tm.AddTodo("Persistent Task 1")
	if err != nil {
		t.Fatalf("failed to add: %v", err)
	}
	_, err = tm.AddTodo("Persistent Task 2")
	if err != nil {
		t.Fatalf("failed to add: %v", err)
	}

	// Wait for async write
	var count int
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		if err := db.QueryRow("SELECT COUNT(*) FROM todo_items").Scan(&count); err == nil && count == 2 {
			break
		}
	}
	if count != 2 {
		t.Fatalf("expected 2 items in sqlite, got %d", count)
	}

	// Toggle in memory and DB
	_ = tm.ToggleTodo(it1.ID)
	for i := 0; i < 20; i++ {
		time.Sleep(20 * time.Millisecond)
		var completed bool
		if err := db.QueryRow("SELECT completed FROM todo_items WHERE id = ?", it1.ID).Scan(&completed); err == nil && completed {
			break
		}
	}

	// Reload from fresh TodoManager with same DB
	tm2 := &TodoManager{
		db:    db,
		items: make([]TodoItem, 0),
	}
	tm2.loadTodosFromDB()

	loaded := tm2.GetTodos()
	if len(loaded) != 2 {
		t.Fatalf("expected 2 loaded items, got %d", len(loaded))
	}
}
