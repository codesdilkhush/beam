package backend

import (
	"sort"
	"strings"
	"sync"
)

type EmojiItem struct {
	Emoji     string   `json:"emoji"`
	Name      string   `json:"name"`
	Shortcode string   `json:"shortcode"`
	Category  string   `json:"category"`
	Keywords  []string `json:"keywords"`
}

type EmojiManager struct {
	items []EmojiItem
}

var (
	emojiInstance *EmojiManager
	emojiOnce     sync.Once
)

func GetEmojiManager() *EmojiManager {
	emojiOnce.Do(func() {
		emojiInstance = &EmojiManager{
			items: loadDefaultEmojiCatalog(),
		}
	})
	return emojiInstance
}

func (em *EmojiManager) GetAll() []EmojiItem {
	res := make([]EmojiItem, len(em.items))
	copy(res, em.items)
	return res
}

type scoredEmoji struct {
	item  EmojiItem
	score int
}

func (em *EmojiManager) SearchEmojis(query string) []EmojiItem {
	trimmed := strings.TrimSpace(query)
	// Strip leading : or /emoji if present
	trimmed = strings.TrimPrefix(trimmed, "/emoji")
	trimmed = strings.TrimPrefix(trimmed, "/e")
	trimmed = strings.TrimSpace(trimmed)
	cleanQuery := strings.Trim(trimmed, ":")

	if cleanQuery == "" {
		// Return top popular subset
		if len(em.items) > 50 {
			res := make([]EmojiItem, 50)
			copy(res, em.items[:50])
			return res
		}
		return em.GetAll()
	}

	var matches []scoredEmoji
	for _, it := range em.items {
		bestScore := 0

		// Match shortcode
		if sc := FuzzyScore(cleanQuery, strings.Trim(it.Shortcode, ":")); sc.Matched && sc.Score+20 > bestScore {
			bestScore = sc.Score + 20 // shortcode bonus
		}

		// Match name
		if ns := FuzzyScore(cleanQuery, it.Name); ns.Matched && ns.Score > bestScore {
			bestScore = ns.Score
		}

		// Match category
		if cs := FuzzyScore(cleanQuery, it.Category); cs.Matched && cs.Score-10 > bestScore {
			bestScore = cs.Score - 10
		}

		// Match keywords
		for _, kw := range it.Keywords {
			if ks := FuzzyScore(cleanQuery, kw); ks.Matched && ks.Score > bestScore {
				bestScore = ks.Score
			}
		}

		if bestScore > 0 {
			matches = append(matches, scoredEmoji{
				item:  it,
				score: bestScore,
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].score != matches[j].score {
			return matches[i].score > matches[j].score
		}
		return matches[i].item.Name < matches[j].item.Name
	})

	limit := 60
	if len(matches) < limit {
		limit = len(matches)
	}

	res := make([]EmojiItem, limit)
	for i := 0; i < limit; i++ {
		res[i] = matches[i].item
	}
	return res
}

func loadDefaultEmojiCatalog() []EmojiItem {
	return []EmojiItem{
		// Smileys & Emotion
		{"😀", "Grinning Face", ":grinning:", "Smileys & Emotion", []string{"face", "smile", "happy", "joy"}},
		{"😃", "Grinning Face with Big Eyes", ":smiley:", "Smileys & Emotion", []string{"face", "happy", "joy", "haha"}},
		{"😄", "Grinning Face with Smiling Eyes", ":smile:", "Smileys & Emotion", []string{"face", "happy", "joy", "laugh"}},
		{"😁", "Beaming Face with Smiling Eyes", ":grin:", "Smileys & Emotion", []string{"face", "happy", "smile", "grin"}},
		{"😆", "Grinning Squinting Face", ":laughing:", "Smileys & Emotion", []string{"face", "happy", "haha", "laugh"}},
		{"😅", "Grinning Face with Sweat", ":sweat_smile:", "Smileys & Emotion", []string{"face", "hot", "happy", "relief"}},
		{"🤣", "Rolling on the Floor Laughing", ":rofl:", "Smileys & Emotion", []string{"face", "rotfl", "laughing", "lol"}},
		{"😂", "Face with Tears of Joy", ":joy:", "Smileys & Emotion", []string{"face", "cry", "tears", "laugh", "happy"}},
		{"🙂", "Slightly Smiling Face", ":slight_smile:", "Smileys & Emotion", []string{"face", "smile", "fine"}},
		{"🙃", "Upside-Down Face", ":upside_down:", "Smileys & Emotion", []string{"face", "silly", "sarcasm", "irony"}},
		{"😉", "Winking Face", ":wink:", "Smileys & Emotion", []string{"face", "wink", "flirt"}},
		{"😊", "Smiling Face with Smiling Eyes", ":blush:", "Smileys & Emotion", []string{"face", "smile", "happy", "blush", "proud"}},
		{"😇", "Smiling Face with Halo", ":innocent:", "Smileys & Emotion", []string{"face", "angel", "halo", "innocent"}},
		{"🥰", "Smiling Face with Hearts", ":smiling_face_with_3_hearts:", "Smileys & Emotion", []string{"face", "love", "crush", "adoring"}},
		{"😍", "Smiling Face with Heart-Eyes", ":heart_eyes:", "Smileys & Emotion", []string{"face", "love", "crush", "heart"}},
		{"🤩", "Star-Struck", ":star_struck:", "Smileys & Emotion", []string{"face", "eyes", "star", "excited"}},
		{"😘", "Face Blowing a Kiss", ":kissing_heart:", "Smileys & Emotion", []string{"face", "kiss", "love", "flirt"}},
		{"😋", "Face Savoring Food", ":yum:", "Smileys & Emotion", []string{"face", "delicious", "yum", "food"}},
		{"😛", "Face with Tongue", ":stuck_out_tongue:", "Smileys & Emotion", []string{"face", "tongue", "silly"}},
		{"😜", "Winking Face with Tongue", ":stuck_out_tongue_winking_eye:", "Smileys & Emotion", []string{"face", "wink", "tongue", "joke"}},
		{"🤪", "Zany Face", ":zany_face:", "Smileys & Emotion", []string{"face", "crazy", "wild", "goofy"}},
		{"🤑", "Money-Mouth Face", ":money_mouth_face:", "Smileys & Emotion", []string{"face", "money", "rich", "dollar"}},
		{"🤗", "Smiling Face with Open Hands", ":hugs:", "Smileys & Emotion", []string{"face", "hug", "welcome", "warm"}},
		{"🤫", "Shushing Face", ":shushing_face:", "Smileys & Emotion", []string{"face", "quiet", "secret", "shh"}},
		{"🤔", "Thinking Face", ":thinking:", "Smileys & Emotion", []string{"face", "hmmm", "think", "ponder"}},
		{"🤐", "Zipper-Mouth Face", ":zipper_mouth_face:", "Smileys & Emotion", []string{"face", "silence", "secret"}},
		{"🤨", "Face with Raised Eyebrow", ":raised_eyebrow:", "Smileys & Emotion", []string{"face", "skeptic", "suspicious"}},
		{"😐", "Neutral Face", ":neutral_face:", "Smileys & Emotion", []string{"face", "meh", "neutral"}},
		{"😑", "Expressionless Face", ":expressionless:", "Smileys & Emotion", []string{"face", "blank", "unimpressed"}},
		{"😶", "Face Without Mouth", ":no_mouth:", "Smileys & Emotion", []string{"face", "mute", "quiet"}},
		{"😏", "Smirking Face", ":smirk:", "Smileys & Emotion", []string{"face", "smug", "flirt"}},
		{"😒", "Unamused Face", ":unamused:", "Smileys & Emotion", []string{"face", "bored", "unhappy"}},
		{"🙄", "Face with Rolling Eyes", ":roll_eyes:", "Smileys & Emotion", []string{"face", "eyeroll", "frustrated"}},
		{"😬", "Grimacing Face", ":grimacing:", "Smileys & Emotion", []string{"face", "awkward", "nervous"}},
		{"🤥", "Lying Face", ":lying_face:", "Smileys & Emotion", []string{"face", "lie", "pinocchio"}},
		{"😌", "Relieved Face", ":relieved:", "Smileys & Emotion", []string{"face", "calm", "peace", "relief"}},
		{"😔", "Pensive Face", ":pensive:", "Smileys & Emotion", []string{"face", "sad", "depressed", "thoughtful"}},
		{"😪", "Sleepy Face", ":sleepy:", "Smileys & Emotion", []string{"face", "tired", "sleep"}},
		{"🤤", "Drooling Face", ":drooling_face:", "Smileys & Emotion", []string{"face", "drool", "food", "sleep"}},
		{"😴", "Sleeping Face", ":sleeping:", "Smileys & Emotion", []string{"face", "tired", "sleep", "zzz"}},
		{"😷", "Face with Medical Mask", ":mask:", "Smileys & Emotion", []string{"face", "sick", "virus", "doctor"}},
		{"🤒", "Face with Thermometer", ":face_with_thermometer:", "Smileys & Emotion", []string{"face", "sick", "fever", "ill"}},
		{"🤕", "Face with Head-Bandage", ":face_with_head_bandage:", "Smileys & Emotion", []string{"face", "hurt", "injury"}},
		{"🤢", "Nauseated Face", ":nauseated_face:", "Smileys & Emotion", []string{"face", "vomit", "gross", "sick"}},
		{"🤮", "Face Vomiting", ":vomiting_face:", "Smileys & Emotion", []string{"face", "puke", "sick", "vomit"}},
		{"🤧", "Sneezing Face", ":sneezing_face:", "Smileys & Emotion", []string{"face", "sneeze", "cold", "allergy"}},
		{"🥵", "Hot Face", ":hot_face:", "Smileys & Emotion", []string{"face", "heat", "sweat", "red"}},
		{"🥶", "Cold Face", ":cold_face:", "Smileys & Emotion", []string{"face", "freezing", "ice", "blue"}},
		{"🥴", "Woozy Face", ":woozy_face:", "Smileys & Emotion", []string{"face", "dizzy", "drunk", "intoxicated"}},
		{"😵", "Dizzy Face", ":dizzy_face:", "Smileys & Emotion", []string{"face", "dizzy", "knocked_out"}},
		{"🤯", "Exploding Head", ":exploding_head:", "Smileys & Emotion", []string{"face", "mind_blown", "shocked"}},
		{"🤠", "Cowboy Hat Face", ":cowboy_hat_face:", "Smileys & Emotion", []string{"face", "cowboy", "wild_west"}},
		{"🥳", "Partying Face", ":partying_face:", "Smileys & Emotion", []string{"face", "party", "celebration", "birthday"}},
		{"😎", "Smiling Face with Sunglasses", ":sunglasses:", "Smileys & Emotion", []string{"face", "cool", "summer", "sun"}},
		{"🤓", "Nerd Face", ":nerd_face:", "Smileys & Emotion", []string{"face", "geek", "glasses", "smart"}},
		{"🧐", "Face with Monocle", ":monocle_face:", "Smileys & Emotion", []string{"face", "inspect", "curious"}},
		{"😕", "Confused Face", ":confused:", "Smileys & Emotion", []string{"face", "what", "indifferent"}},
		{"😟", "Worried Face", ":worried:", "Smileys & Emotion", []string{"face", "nervous", "concern"}},
		{"🙁", "Slightly Frowning Face", ":slightly_frowning_face:", "Smileys & Emotion", []string{"face", "sad", "disappointed"}},
		{"😮", "Face with Open Mouth", ":open_mouth:", "Smileys & Emotion", []string{"face", "surprise", "wow"}},
		{"😯", "Hushed Face", ":hushed:", "Smileys & Emotion", []string{"face", "quiet", "surprise"}},
		{"😲", "Astonished Face", ":astonished:", "Smileys & Emotion", []string{"face", "shocked", "amazed"}},
		{"😳", "Flushed Face", ":flushed:", "Smileys & Emotion", []string{"face", "blush", "embarrassed"}},
		{"🥺", "Pleading Face", ":pleading_face:", "Smileys & Emotion", []string{"face", "puppy_eyes", "begging", "mercy"}},
		{"😦", "Frowning Face with Open Mouth", ":frowning:", "Smileys & Emotion", []string{"face", "sad", "shocked"}},
		{"😧", "Anguished Face", ":anguished:", "Smileys & Emotion", []string{"face", "pain", "suffering"}},
		{"😨", "Fearful Face", ":fearful:", "Smileys & Emotion", []string{"face", "scared", "afraid"}},
		{"😰", "Anxious Face with Sweat", ":cold_sweat:", "Smileys & Emotion", []string{"face", "nervous", "worried"}},
		{"😥", "Sad but Relieved Face", ":disappointed_relieved:", "Smileys & Emotion", []string{"face", "phew", "sweat"}},
		{"😢", "Crying Face", ":cry:", "Smileys & Emotion", []string{"face", "tear", "sad", "depressed"}},
		{"😭", "Loudly Crying Face", ":sob:", "Smileys & Emotion", []string{"face", "cry", "bawling", "heartbroken"}},
		{"😱", "Face Screaming in Fear", ":scream:", "Smileys & Emotion", []string{"face", "munch", "horror", "shock"}},
		{"😖", "Confounded Face", ":confounded:", "Smileys & Emotion", []string{"face", "confused", "sick"}},
		{"😣", "Persevering Face", ":persevere:", "Smileys & Emotion", []string{"face", "struggle", "effort"}},
		{"😞", "Disappointed Face", ":disappointed:", "Smileys & Emotion", []string{"face", "sad", "down"}},
		{"😓", "Downcast Face with Sweat", ":sweat:", "Smileys & Emotion", []string{"face", "hard_work", "exhausted"}},
		{"😩", "Weary Face", ":weary:", "Smileys & Emotion", []string{"face", "tired", "stressed"}},
		{"😫", "Tired Face", ":tired_face:", "Smileys & Emotion", []string{"face", "exhausted", "done"}},
		{"🥱", "Yawning Face", ":yawning_face:", "Smileys & Emotion", []string{"face", "tired", "sleepy", "bored"}},
		{"😤", "Face with Steam From Nose", ":triumph:", "Smileys & Emotion", []string{"face", "proud", "angry", "frustrated"}},
		{"😡", "Enraged Face", ":rage:", "Smileys & Emotion", []string{"face", "mad", "angry", "furious"}},
		{"😠", "Angry Face", ":angry:", "Smileys & Emotion", []string{"face", "mad", "annoyed"}},
		{"🤬", "Face with Symbols on Mouth", ":cursing_face:", "Smileys & Emotion", []string{"face", "swearing", "curse", "profanity"}},
		{"😈", "Smiling Face with Horns", ":smiling_imp:", "Smileys & Emotion", []string{"face", "devil", "evil", "horns"}},
		{"👿", "Angry Face with Horns", ":imp:", "Smileys & Emotion", []string{"face", "demon", "devil"}},
		{"💀", "Skull", ":skull:", "Smileys & Emotion", []string{"dead", "skeleton", "death", "danger"}},
		{"☠️", "Skull and Crossbones", ":skull_and_crossbones:", "Smileys & Emotion", []string{"dead", "pirate", "poison", "danger"}},
		{"💩", "Pile of Poo", ":poop:", "Smileys & Emotion", []string{"shit", "turd", "crap"}},
		{"🤡", "Clown Face", ":clown_face:", "Smileys & Emotion", []string{"circus", "clown", "fool"}},
		{"👻", "Ghost", ":ghost:", "Smileys & Emotion", []string{"spooky", "halloween", "spirit"}},
		{"👽", "Alien", ":alien:", "Smileys & Emotion", []string{"ufo", "space", "extraterrestrial"}},
		{"🤖", "Robot", ":robot:", "Smileys & Emotion", []string{"bot", "ai", "droid", "machine"}},

		// Gestures & People
		{"👋", "Waving Hand", ":wave:", "People & Body", []string{"hand", "hello", "goodbye", "hi"}},
		{"🤚", "Raised Back of Hand", ":raised_back_of_hand:", "People & Body", []string{"hand", "stop", "backhand"}},
		{"✋", "Raised Hand", ":raised_hand:", "People & Body", []string{"hand", "highfive", "stop"}},
		{"🖖", "Vulcan Salute", ":vulcan_salute:", "People & Body", []string{"hand", "spock", "star_trek"}},
		{"👌", "OK Hand", ":ok_hand:", "People & Body", []string{"hand", "perfect", "good", "okay"}},
		{"🤌", "Pinched Fingers", ":pinched_fingers:", "People & Body", []string{"hand", "italian", "gesture", "what"}},
		{"✌️", "Victory Hand", ":v:", "People & Body", []string{"hand", "peace", "victory", "two"}},
		{"🤞", "Crossed Fingers", ":crossed_fingers:", "People & Body", []string{"hand", "luck", "hope", "wish"}},
		{"🤟", "Love-You Gesture", ":love_you_gesture:", "People & Body", []string{"hand", "love", "ily"}},
		{"🤘", "Sign of the Horns", ":metal:", "People & Body", []string{"hand", "rock", "heavy_metal"}},
		{"🤙", "Call Me Hand", ":call_me_hand:", "People & Body", []string{"hand", "shaka", "hang_loose", "phone"}},
		{"👈", "Backhand Index Pointing Left", ":point_left:", "People & Body", []string{"hand", "direction", "left"}},
		{"👉", "Backhand Index Pointing Right", ":point_right:", "People & Body", []string{"hand", "direction", "right"}},
		{"👆", "Backhand Index Pointing Up", ":point_up_2:", "People & Body", []string{"hand", "direction", "up"}},
		{"👇", "Backhand Index Pointing Down", ":point_down:", "People & Body", []string{"hand", "direction", "down"}},
		{"👍", "Thumbs Up", ":+1:", "People & Body", []string{"hand", "like", "approve", "yes", "good"}},
		{"👎", "Thumbs Down", ":-1:", "People & Body", []string{"hand", "dislike", "disapprove", "no"}},
		{"✊", "Raised Fist", ":fist:", "People & Body", []string{"hand", "power", "strength", "solidarity"}},
		{"👊", "Oncoming Fist", ":punch:", "People & Body", []string{"hand", "fistbump", "brofist"}},
		{"🤛", "Left-Facing Fist", ":left_facing_fist:", "People & Body", []string{"hand", "fistbump"}},
		{"🤜", "Right-Facing Fist", ":right_facing_fist:", "People & Body", []string{"hand", "fistbump"}},
		{"👏", "Clapping Hands", ":clap:", "People & Body", []string{"hand", "applause", "praise", "bravo"}},
		{"🙌", "Raising Hands", ":raised_hands:", "People & Body", []string{"hand", "celebration", "praise", "hooray"}},
		{"👐", "Open Hands", ":open_hands:", "People & Body", []string{"hand", "hug", "jazz_hands"}},
		{"🤲", "Palms Up Together", ":palms_up_together:", "People & Body", []string{"hand", "prayer", "dua", "beg"}},
		{"🤝", "Handshake", ":handshake:", "People & Body", []string{"hand", "agreement", "deal", "shake"}},
		{"🙏", "Folded Hands", ":pray:", "People & Body", []string{"hand", "please", "thank_you", "prayer", "namaste"}},
		{"💪", "Flexed Biceps", ":muscle:", "People & Body", []string{"arm", "strong", "flex", "workout"}},
		{"🧠", "Brain", ":brain:", "People & Body", []string{"mind", "smart", "intellect", "thinking"}},
		{"👀", "Eyes", ":eyes:", "People & Body", []string{"look", "see", "watch", "peeping"}},

		// Hearts & Symbols
		{"❤️", "Red Heart", ":heart:", "Symbols", []string{"love", "like", "valentines"}},
		{"🧡", "Orange Heart", ":orange_heart:", "Symbols", []string{"love", "orange"}},
		{"💛", "Yellow Heart", ":yellow_heart:", "Symbols", []string{"love", "friendship"}},
		{"💚", "Green Heart", ":green_heart:", "Symbols", []string{"love", "nature", "organic"}},
		{"💙", "Blue Heart", ":blue_heart:", "Symbols", []string{"love", "blue", "trust"}},
		{"💜", "Purple Heart", ":purple_heart:", "Symbols", []string{"love", "purple", "honor"}},
		{"🖤", "Black Heart", ":black_heart:", "Symbols", []string{"love", "dark", "goth"}},
		{"🤍", "White Heart", ":white_heart:", "Symbols", []string{"love", "pure", "clean"}},
		{"💔", "Broken Heart", ":broken_heart:", "Symbols", []string{"heartbreak", "sad", "grief"}},
		{"✨", "Sparkles", ":sparkles:", "Symbols", []string{"shine", "magic", "clean", "glitter", "star"}},
		{"⭐️", "Star", ":star:", "Symbols", []string{"favorite", "night", "yellow"}},
		{"🌟", "Glowing Star", ":star2:", "Symbols", []string{"shine", "sparkle", "star"}},
		{"💥", "Collision", ":boom:", "Symbols", []string{"explode", "bang", "impact"}},
		{"🔥", "Fire", ":fire:", "Symbols", []string{"flame", "hot", "lit", "burn", "trend"}},
		{"⚡️", "High Voltage", ":zap:", "Symbols", []string{"lightning", "thunder", "fast", "electric"}},
		{"☀️", "Sun", ":sunny:", "Symbols", []string{"weather", "sunshine", "warm"}},
		{"🌙", "Crescent Moon", ":crescent_moon:", "Symbols", []string{"night", "sleep", "dark"}},
		{"🌈", "Rainbow", ":rainbow:", "Symbols", []string{"color", "pride", "weather"}},

		// Objects & Tools
		{"🚀", "Rocket", ":rocket:", "Travel & Places", []string{"ship", "launch", "space", "fast", "startup"}},
		{"🛸", "Flying Saucer", ":flying_saucer:", "Travel & Places", []string{"ufo", "alien", "space"}},
		{"💻", "Laptop", ":laptop:", "Objects", []string{"computer", "work", "tech", "pc", "mac"}},
		{"🖥️", "Desktop Computer", ":desktop_computer:", "Objects", []string{"screen", "pc", "tech"}},
		{"📱", "Mobile Phone", ":iphone:", "Objects", []string{"phone", "smartphone", "call", "ios", "android"}},
		{"⌨️", "Keyboard", ":keyboard:", "Objects", []string{"type", "computer", "input"}},
		{"🖱️", "Computer Mouse", ":computer_mouse:", "Objects", []string{"click", "hardware"}},
		{"📷", "Camera", ":camera:", "Objects", []string{"photo", "picture", "snapshot"}},
		{"🔍", "Magnifying Glass Tilted Left", ":mag:", "Objects", []string{"search", "find", "query", "zoom"}},
		{"💡", "Light Bulb", ":bulb:", "Objects", []string{"idea", "light", "energy", "bright"}},
		{"📦", "Package", ":package:", "Objects", []string{"box", "delivery", "shipping", "parcel"}},
		{"🔒", "Locked", ":lock:", "Objects", []string{"security", "privacy", "secure", "safe"}},
		{"🔓", "Unlocked", ":unlock:", "Objects", []string{"open", "security", "insecure"}},
		{"🔑", "Key", ":key:", "Objects", []string{"password", "access", "door"}},
		{"🔨", "Hammer", ":hammer:", "Objects", []string{"tool", "build", "fix"}},
		{"⚙️", "Gear", ":gear:", "Objects", []string{"settings", "cog", "preferences", "config"}},
		{"🔧", "Wrench", ":wrench:", "Objects", []string{"tool", "fix", "repair"}},
		{"☕️", "Hot Beverage", ":coffee:", "Food & Drink", []string{"drink", "tea", "caffeine", "mug", "cafe"}},
		{"🍕", "Pizza", ":pizza:", "Food & Drink", []string{"food", "cheese", "slice", "italian"}},
		{"🍔", "Hamburger", ":hamburger:", "Food & Drink", []string{"food", "burger", "beef", "fastfood"}},
		{"🍺", "Beer Mug", ":beer:", "Food & Drink", []string{"drink", "alcohol", "cheers", "pub"}},
		{"🎉", "Party Popper", ":tada:", "Activities", []string{"celebration", "party", "congrats", "hooray"}},
		{"🏆", "Trophy", ":trophy:", "Activities", []string{"win", "champion", "award", "first"}},
		{"🎯", "Direct Hit", ":dart:", "Activities", []string{"bullseye", "target", "goal", "accuracy"}},
		{"🎨", "Artist Palette", ":art:", "Activities", []string{"color", "paint", "design", "creative"}},
		{"🎵", "Musical Note", ":musical_note:", "Activities", []string{"music", "sound", "song"}},
		{"✅", "Check Mark Button", ":white_check_mark:", "Symbols", []string{"done", "complete", "yes", "ok", "todo"}},
		{"❌", "Cross Mark", ":x:", "Symbols", []string{"cancel", "no", "delete", "error"}},
		{"⚠️", "Warning", ":warning:", "Symbols", []string{"alert", "danger", "caution"}},
		{"🚨", "Police Car Light", ":rotating_light:", "Symbols", []string{"emergency", "siren", "alert", "police"}},
		{"💯", "Hundred Points", ":100:", "Symbols", []string{"score", "perfect", "full", "hundred"}},
	}
}
