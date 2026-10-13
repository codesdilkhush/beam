package backend

import (
	"testing"
)

func TestFuzzyScore(t *testing.T) {
	tests := []struct {
		query   string
		target  string
		matched bool
	}{
		{"gc", "Google Chrome", true},
		{"vsc", "Visual Studio Code", true},
		{"term", "Terminal", true},
		{"fire", "Firefox Web Browser", true},
		{"calc", "Calculator", true},
		{"xyz123", "Terminal", false},
	}

	for _, tt := range tests {
		res := FuzzyScore(tt.query, tt.target)
		if res.Matched != tt.matched {
			t.Errorf("FuzzyScore(%q, %q).Matched = %v; want %v", tt.query, tt.target, res.Matched, tt.matched)
		}
	}
}
