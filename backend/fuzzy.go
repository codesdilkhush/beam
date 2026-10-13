package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"
)

type FuzzyMatchResult struct {
	Score   int
	Matched bool
}

type StatsTracker struct {
	mu       sync.RWMutex
	filePath string
	Counts   map[string]int `json:"counts"`
}

var (
	statsInstance *StatsTracker
	statsOnce     sync.Once
)

func GetStatsTracker() *StatsTracker {
	statsOnce.Do(func() {
		dir, err := os.UserHomeDir()
		if err != nil {
			dir = "."
		}
		storePath := filepath.Join(dir, ".local/share/beam", "stats.json")
		if mkErr := os.MkdirAll(filepath.Dir(storePath), 0755); mkErr != nil {
			// fallback if directory creation fails
			storePath = "stats.json"
		}

		st := &StatsTracker{
			filePath: storePath,
			Counts:   make(map[string]int),
		}

		if data, readErr := os.ReadFile(storePath); readErr == nil {
			if unmarshalErr := json.Unmarshal(data, &st.Counts); unmarshalErr != nil {
				st.Counts = make(map[string]int)
			}
		}

		statsInstance = st
	})
	return statsInstance
}

func (st *StatsTracker) Increment(appID string) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.Counts[appID]++

	data, err := json.Marshal(st.Counts)
	if err != nil {
		return
	}
	if writeErr := os.WriteFile(st.filePath, data, 0644); writeErr != nil {
		return
	}
}

func (st *StatsTracker) GetCount(appID string) int {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.Counts[appID]
}

// FuzzyScore calculates relevance between search query and target string
func FuzzyScore(query, target string) FuzzyMatchResult {
	q := strings.ToLower(strings.TrimSpace(query))
	t := strings.ToLower(strings.TrimSpace(target))

	if q == "" {
		return FuzzyMatchResult{Score: 0, Matched: true}
	}
	if t == "" {
		return FuzzyMatchResult{Score: 0, Matched: false}
	}

	// 1. Exact match
	if q == t {
		return FuzzyMatchResult{Score: 1000, Matched: true}
	}

	// 2. Exact prefix match
	if strings.HasPrefix(t, q) {
		return FuzzyMatchResult{Score: 500 + (len(q) * 10), Matched: true}
	}

	// 3. Acronym match (e.g. "gc" -> "Google Chrome", "vsc" -> "Visual Studio Code")
	words := strings.FieldsFunc(target, func(r rune) bool {
		return unicode.IsSpace(r) || r == '-' || r == '_' || r == '.'
	})

	var acronym strings.Builder
	for _, w := range words {
		if len(w) > 0 {
			acronym.WriteRune(unicode.ToLower([]rune(w)[0]))
		}
	}
	acronymStr := acronym.String()
	if strings.HasPrefix(acronymStr, q) {
		return FuzzyMatchResult{Score: 400 + (len(q) * 20), Matched: true}
	}
	if strings.Contains(acronymStr, q) {
		return FuzzyMatchResult{Score: 300 + (len(q) * 10), Matched: true}
	}

	// 4. Word boundary match
	for _, w := range words {
		wLower := strings.ToLower(w)
		if strings.HasPrefix(wLower, q) {
			return FuzzyMatchResult{Score: 350 + (len(q) * 10), Matched: true}
		}
	}

	// 5. Substring match
	if idx := strings.Index(t, q); idx >= 0 {
		// Earlier in the string is higher score
		score := 250 - idx*5
		if score < 150 {
			score = 150
		}
		return FuzzyMatchResult{Score: score, Matched: true}
	}

	// 6. Subsequence fuzzy match (fzf style)
	qRunes := []rune(q)
	tRunes := []rune(t)

	qIdx := 0
	score := 100
	lastMatchIdx := -1

	for tIdx := 0; tIdx < len(tRunes) && qIdx < len(qRunes); tIdx++ {
		if tRunes[tIdx] == qRunes[qIdx] {
			if lastMatchIdx != -1 {
				gap := tIdx - lastMatchIdx
				if gap == 1 {
					score += 20 // Consecutive characters bonus
				} else {
					score -= gap * 2
				}
			}
			lastMatchIdx = tIdx
			qIdx++
		}
	}

	if qIdx == len(qRunes) {
		if score < 10 {
			score = 10
		}
		return FuzzyMatchResult{Score: score, Matched: true}
	}

	return FuzzyMatchResult{Score: 0, Matched: false}
}
