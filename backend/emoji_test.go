package backend

import (
	"testing"
)

func TestEmojiManager_GetAll(t *testing.T) {
	em := GetEmojiManager()
	all := em.GetAll()
	if len(all) == 0 {
		t.Fatalf("expected non-empty emoji catalog, got 0")
	}
}

func TestEmojiManager_SearchEmojis(t *testing.T) {
	em := GetEmojiManager()

	// 1. Empty query returns top subset
	res := em.SearchEmojis("")
	if len(res) == 0 {
		t.Errorf("expected results for empty query, got 0")
	}

	// 2. Search by name
	fireRes := em.SearchEmojis("fire")
	if len(fireRes) == 0 {
		t.Fatalf("expected results searching for 'fire', got 0")
	}
	if fireRes[0].Emoji != "🔥" {
		t.Errorf("expected top result for 'fire' to be '🔥', got %s (%s)", fireRes[0].Emoji, fireRes[0].Name)
	}

	// 3. Search by shortcode
	rocketRes := em.SearchEmojis(":rocket:")
	if len(rocketRes) == 0 {
		t.Fatalf("expected results searching for ':rocket:', got 0")
	}
	if rocketRes[0].Emoji != "🚀" {
		t.Errorf("expected top result for ':rocket:' to be '🚀', got %s", rocketRes[0].Emoji)
	}

	// 4. Search with `/emoji` prefix
	prefixRes := em.SearchEmojis("/emoji laugh")
	if len(prefixRes) == 0 {
		t.Fatalf("expected results for '/emoji laugh', got 0")
	}
	foundJoyOrLaugh := false
	for _, it := range prefixRes[:3] {
		if it.Emoji == "😂" || it.Emoji == "😆" || it.Emoji == "🤣" {
			foundJoyOrLaugh = true
			break
		}
	}
	if !foundJoyOrLaugh {
		t.Errorf("expected laugh emoji in top 3 results, got %v", prefixRes[:3])
	}

	// 5. Search by keyword
	coffeeRes := em.SearchEmojis("caffeine")
	if len(coffeeRes) == 0 || coffeeRes[0].Emoji != "☕️" {
		t.Errorf("expected coffee emoji searching for 'caffeine', got %v", coffeeRes)
	}
}
