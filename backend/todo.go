package backend

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type TodoItem struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

type TodoManager struct {
	mu    sync.RWMutex
	db    *sql.DB
	items []TodoItem
}

var (
	todoInstance *TodoManager
	todoOnce     sync.Once
)

func GetTodoManager() *TodoManager {
	todoOnce.Do(func() {
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

		tm := &TodoManager{
			db:    db,
			items: make([]TodoItem, 0),
		}

		if tm.db != nil {
			tm.initDatabase()
			tm.loadTodosFromDB()
		}

		todoInstance = tm
	})
	return todoInstance
}

func (tm *TodoManager) initDatabase() {
	if tm.db == nil {
		return
	}
	tm.db.SetMaxOpenConns(1)

	query := `
	PRAGMA busy_timeout = 5000;
	CREATE TABLE IF NOT EXISTS todo_items (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL,
		completed_at DATETIME
	);
	CREATE INDEX IF NOT EXISTS idx_todo_created ON todo_items(created_at DESC);
	`
	if _, err := tm.db.Exec(query); err != nil {
		return
	}
}

func (tm *TodoManager) loadTodosFromDB() {
	if tm.db == nil {
		return
	}

	rows, err := tm.db.Query("SELECT id, title, completed, created_at, completed_at FROM todo_items ORDER BY completed ASC, created_at DESC LIMIT 200")
	if err != nil {
		return
	}
	defer rows.Close()

	var loaded []TodoItem
	for rows.Next() {
		var item TodoItem
		var compAt sql.NullTime
		if scanErr := rows.Scan(&item.ID, &item.Title, &item.Completed, &item.CreatedAt, &compAt); scanErr == nil {
			if compAt.Valid {
				t := compAt.Time
				item.CompletedAt = &t
			}
			loaded = append(loaded, item)
		}
	}
	if err := rows.Err(); err != nil {
		return
	}

	tm.mu.Lock()
	tm.items = loaded
	tm.mu.Unlock()
}

func (tm *TodoManager) GetTodos() []TodoItem {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	res := make([]TodoItem, len(tm.items))
	copy(res, tm.items)
	return res
}

func (tm *TodoManager) AddTodo(title string) (*TodoItem, error) {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return nil, fmt.Errorf("todo title cannot be empty")
	}

	b := make([]byte, 8)
	_, _ = rand.Read(b)
	id := fmt.Sprintf("%x", b)

	item := TodoItem{
		ID:        id,
		Title:     trimmed,
		Completed: false,
		CreatedAt: time.Now(),
	}

	tm.mu.Lock()
	// New pending items go to the top of active items
	newItems := make([]TodoItem, 0, len(tm.items)+1)
	newItems = append(newItems, item)
	newItems = append(newItems, tm.items...)
	tm.items = newItems
	tm.mu.Unlock()

	if tm.db != nil {
		go func(it TodoItem) {
			stmt := "INSERT INTO todo_items (id, title, completed, created_at) VALUES (?, ?, ?, ?)"
			_, _ = tm.db.Exec(stmt, it.ID, it.Title, it.Completed, it.CreatedAt)
		}(item)
	}

	return &item, nil
}

func (tm *TodoManager) ToggleTodo(id string) error {
	tm.mu.Lock()
	var found *TodoItem
	for i := range tm.items {
		if tm.items[i].ID == id {
			tm.items[i].Completed = !tm.items[i].Completed
			if tm.items[i].Completed {
				now := time.Now()
				tm.items[i].CompletedAt = &now
			} else {
				tm.items[i].CompletedAt = nil
			}
			found = &tm.items[i]
			break
		}
	}

	// Sort active first, then completed
	if found != nil {
		active := make([]TodoItem, 0, len(tm.items))
		completed := make([]TodoItem, 0, len(tm.items))
		for _, it := range tm.items {
			if it.Completed {
				completed = append(completed, it)
			} else {
				active = append(active, it)
			}
		}
		tm.items = append(active, completed...)
	}
	tm.mu.Unlock()

	if found == nil {
		return fmt.Errorf("todo not found: %s", id)
	}

	if tm.db != nil {
		go func(it TodoItem) {
			stmt := "UPDATE todo_items SET completed = ?, completed_at = ? WHERE id = ?"
			_, _ = tm.db.Exec(stmt, it.Completed, it.CompletedAt, it.ID)
		}(*found)
	}

	return nil
}

func (tm *TodoManager) DeleteTodo(id string) error {
	tm.mu.Lock()
	newItems := make([]TodoItem, 0, len(tm.items))
	found := false
	for _, it := range tm.items {
		if it.ID == id {
			found = true
			continue
		}
		newItems = append(newItems, it)
	}
	tm.items = newItems
	tm.mu.Unlock()

	if !found {
		return fmt.Errorf("todo not found: %s", id)
	}

	if tm.db != nil {
		go func(todoID string) {
			_, _ = tm.db.Exec("DELETE FROM todo_items WHERE id = ?", todoID)
		}(id)
	}

	return nil
}

func (tm *TodoManager) ClearCompleted() error {
	tm.mu.Lock()
	active := make([]TodoItem, 0, len(tm.items))
	for _, it := range tm.items {
		if !it.Completed {
			active = append(active, it)
		}
	}
	tm.items = active
	tm.mu.Unlock()

	if tm.db != nil {
		go func() {
			_, _ = tm.db.Exec("DELETE FROM todo_items WHERE completed = 1")
		}()
	}

	return nil
}
