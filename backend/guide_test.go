package backend

import (
	"strings"
	"testing"
)

func TestGetGuideTopics(t *testing.T) {
	topics := GetGuideTopics()
	if len(topics) < 8 {
		t.Fatalf("expected at least 8 guide topics, got %d", len(topics))
	}

	expectedIDs := []string{
		"getting-started",
		"app-launcher",
		"clipboard",
		"todo",
		"emoji",
		"converter",
		"actions",
		"themes",
		"scripts",
		"shortcuts",
	}

	foundIDs := make(map[string]bool)
	for _, topic := range topics {
		foundIDs[topic.ID] = true
		if topic.Title == "" {
			t.Errorf("topic %s has empty title", topic.ID)
		}
		if topic.Category == "" {
			t.Errorf("topic %s has empty category", topic.ID)
		}
		if topic.Summary == "" {
			t.Errorf("topic %s has empty summary", topic.ID)
		}
		if len(topic.Content) == 0 {
			t.Errorf("topic %s has no content lines", topic.ID)
		}
	}

	for _, id := range expectedIDs {
		if !foundIDs[id] {
			t.Errorf("expected topic ID %q to exist in guide catalog", id)
		}
	}
}

func TestSearchGuideTopics(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		wantMin     int
		shouldMatch string
	}{
		{"Empty query returns all", "", 10, "getting-started"},
		{"Search by title keyword", "emoji", 1, "emoji"},
		{"Search by category", "productivity", 2, "todo"},
		{"Search by summary keyword", "clipboard", 1, "clipboard"},
		{"Search by content keyword", "fzf", 1, "app-launcher"},
		{"Search converter", "convert", 1, "converter"},
		{"Search actions", "ctrl+k", 1, "actions"},
		{"Search non-existent", "xyznonexistentquery123", 0, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results := SearchGuideTopics(tc.query)
			if len(results) < tc.wantMin {
				t.Fatalf("for query %q, expected at least %d results, got %d", tc.query, tc.wantMin, len(results))
			}
			if tc.shouldMatch != "" {
				found := false
				for _, r := range results {
					if r.ID == tc.shouldMatch || strings.Contains(strings.ToLower(r.Title), tc.shouldMatch) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("for query %q, expected to find %q in results", tc.query, tc.shouldMatch)
				}
			}
		})
	}
}
