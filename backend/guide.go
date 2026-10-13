package backend

import (
	"strings"
)

type GuideShortcut struct {
	Key  string `json:"key"`
	Desc string `json:"desc"`
}

type GuideTopic struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Category  string          `json:"category"`
	Summary   string          `json:"summary"`
	Content   []string        `json:"content"`
	Shortcuts []GuideShortcut `json:"shortcuts"`
	Examples  []string        `json:"examples"`
}

var guideCatalog = []GuideTopic{
	{
		ID:       "getting-started",
		Title:    "Getting Started with Beam",
		Category: "Basics",
		Summary:  "Learn core navigation, keyboard interactions, and the spotlight interface.",
		Content: []string{
			"Beam is an ultra-fast, Raycast-inspired spotlight launcher for Linux designed for zero-friction keyboard navigation.",
			"Invoke Beam instantly using your global desktop shortcut (typically Super+Space or Alt+Space).",
			"Cycle effortlessly through modes (Apps, Clipboard, To-Do, Emoji, and Guide) using the Tab key.",
			"Press Esc to clear your current query or dismiss the launcher instantly.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "Super+Space", Desc: "Toggle Beam window from anywhere"},
			{Key: "Tab", Desc: "Cycle through modes (Apps, Clipboard, To-Do, Emoji, Guide)"},
			{Key: "Esc", Desc: "Clear query or hide window"},
			{Key: "Ctrl+K", Desc: "Open the Action Palette"},
		},
		Examples: []string{
			"Press Tab to switch between Apps, Clipboard, To-Do, and Emoji",
			"Type /clip, /todo, /emoji, or /guide directly into the search bar",
		},
	},
	{
		ID:       "app-launcher",
		Title:    "Application Launcher & Fuzzy Search",
		Category: "Apps",
		Summary:  "Find and launch Linux desktop applications with high-speed fuzzy search and acronym matching.",
		Content: []string{
			"Beam indexes all .desktop applications in /usr/share/applications and ~/.local/share/applications.",
			"Supports fzf-style fuzzy matching, exact prefix matches, and acronym abbreviations (e.g., 'vsc' matches Visual Studio Code, 'gc' matches Google Chrome).",
			"Applications you open more frequently receive automatic score boosting, pushing them to the top of your results.",
			"Launch any application in detached mode with Enter, or launch terminal apps via the Action Palette.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "↵ Enter", Desc: "Launch the selected application"},
			{Key: "Ctrl+K", Desc: "Open actions to copy CLI command or desktop file path"},
			{Key: "↑ / ↓", Desc: "Navigate search results"},
		},
		Examples: []string{
			"Type 'vsc' to quickly select Visual Studio Code",
			"Type 'term' to find terminal emulators",
			"Type 'calc' or math expressions directly",
		},
	},
	{
		ID:       "clipboard",
		Title:    "Persistent SQLite Clipboard Manager",
		Category: "Clipboard",
		Summary:  "Automatically capture, search, and recover clipboard history across sessions.",
		Content: []string{
			"Continuous background tracking captures Wayland clipboard text seamlessly via wl-paste.",
			"Entries are cryptographically deduplicated with SHA-256 and safely persisted in ~/.local/share/beam/beam.db.",
			"Switch to Clipboard mode by pressing Tab or typing /clip in the search bar.",
			"Inspect full content, character count, word count, and timestamps in the detail inspection panel.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "Tab", Desc: "Switch to Clipboard history mode"},
			{Key: "↵ Enter", Desc: "Copy selected history item back to active clipboard"},
			{Key: "↑ / ↓", Desc: "Scroll through clipboard history items"},
		},
		Examples: []string{
			"/clip to enter clipboard mode",
			"/clip react to search history containing 'react'",
		},
	},
	{
		ID:       "todo",
		Title:    "To-Do Task Checklist Manager",
		Category: "Productivity",
		Summary:  "Manage daily tasks with inline quick-add, spacebar completion, and SQLite persistence.",
		Content: []string{
			"Keep lightweight daily tasks in a fast, distraction-free checklist saved in your local SQLite database.",
			"Switch to To-Do mode by pressing Tab or typing /todo in the search bar.",
			"Quickly add new tasks by typing a task title and pressing Enter when no matching task exists, or via Ctrl+K.",
			"Check off completed tasks with Space or Enter. Active tasks always stay organized at the top.",
			"Delete tasks cleanly with Ctrl+D or the Action Palette.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "Space", Desc: "Toggle task completion status"},
			{Key: "↵ Enter", Desc: "Toggle completion, or add new task if not matching"},
			{Key: "Ctrl+D", Desc: "Delete selected task immediately"},
			{Key: "Ctrl+K", Desc: "Open task actions (Add, Toggle, Delete, Clear Completed)"},
		},
		Examples: []string{
			"/todo Buy groceries (press Enter to add task)",
			"Select any task and press Space to check or uncheck it",
			"Press Ctrl+D on a selected task to delete it",
		},
	},
	{
		ID:       "emoji",
		Title:    "8-Column Raycast Emoji & Symbol Picker",
		Category: "Productivity",
		Summary:  "Search thousands of Unicode emojis by name, tags, and shortcodes in an 8-column grid.",
		Content: []string{
			"Open the emoji picker by typing /emoji, /e, a colon (:), or pressing Tab.",
			"Displays an 8-column visual grid of square glyph tiles inspired by Raycast.",
			"Full 2D keyboard navigation allows moving left, right, up, and down across rows seamlessly.",
			"Press Enter on any emoji to copy it to your clipboard immediately.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "← / →", Desc: "Navigate one emoji left or right"},
			{Key: "↑ / ↓", Desc: "Navigate one row up or down (8 columns)"},
			{Key: "↵ Enter", Desc: "Copy selected emoji to clipboard"},
			{Key: "Ctrl+K", Desc: "Open actions to copy shortcode or Unicode codepoint"},
		},
		Examples: []string{
			":rocket to search rocket emojis",
			"/emoji fire to search fire and flame symbols",
			"/emoji smile to browse friendly faces",
		},
	},
	{
		ID:       "converter",
		Title:    "Universal 'Convert Anything' & Calendar Engine",
		Category: "Intelligence",
		Summary:  "Perform real-time multi-unit, currency, color, time, and natural calendar queries inline.",
		Content: []string{
			"Type queries directly in the search bar — Beam detects conversions and displays a dynamic preview banner.",
			"Physical Units: Length (km, m, mi, ft, in), Mass (kg, lb, oz, g), Data (GB, MB, TB), Speed (km/h, mph), Duration (hours, minutes, seconds).",
			"Currencies: Real-time conversion between USD, EUR, GBP, JPY, INR, CAD, AUD, CHF, CNY.",
			"Colors: Convert Hex to RGB/HSL (#3875f6) and RGB to Hex/HSL (rgb(56, 117, 246)).",
			"CSS Units: Seamless 16px-base conversion between px and rem (e.g., '16px to rem', '2rem to px').",
			"Calendar & Date Queries: Natural queries like 'days in march', 'days in current month', 'what is the weekday?', 'days left in year', 'week number', and 'is 2024 a leap year'.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "↵ Enter", Desc: "Copy conversion result banner directly to clipboard"},
		},
		Examples: []string{
			"10 km to mi or 50 kg to lb",
			"100 usd to eur or 1000 inr to usd",
			"#3875f6 or rgb(56, 117, 246)",
			"16px to rem or 2rem to px",
			"days in march or days in current month",
			"what is the weekday? or weekday of 2026-12-25",
			"now to epoch or 1700000000 to date",
		},
	},
	{
		ID:       "actions",
		Title:    "Raycast Action Palette (Ctrl+K)",
		Category: "Navigation",
		Summary:  "Context-aware popover menu providing secondary operations for any selected item.",
		Content: []string{
			"Press Ctrl+K from anywhere in Beam to open the Raycast-style Action Palette.",
			"Displays context-sensitive actions depending on what is currently selected (Apps, Clipboard, Tasks, Emojis, or Themes).",
			"For Apps: Launch, Run in Terminal, Copy Command Line, Copy Desktop Path.",
			"For Clipboard: Copy to Clipboard, Share.",
			"For Tasks: Add Task, Toggle Completion, Delete Task, Clear Completed.",
			"For Emojis: Copy Glyph, Copy Shortcode, Copy Codepoint.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "Ctrl+K", Desc: "Toggle the Action Palette"},
			{Key: "↑ / ↓", Desc: "Navigate actions"},
			{Key: "↵ Enter", Desc: "Execute the highlighted action"},
			{Key: "Esc", Desc: "Dismiss the Action Palette"},
		},
		Examples: []string{
			"Highlight any application and press Ctrl+K to inspect actions",
			"Switch themes quickly from the Action Palette",
		},
	},
	{
		ID:       "themes",
		Title:    "Glassmorphic Themes & Custom Styling",
		Category: "Customization",
		Summary:  "Switch between 4 hand-crafted dark glassmorphic color palettes with persistent storage.",
		Content: []string{
			"Beam features a pure CSS glassmorphism design system with backdrop blur, saturated borders, and radiant glows.",
			"Default Dark Glass: Deep void obsidian canvas with electric sky-blue accents.",
			"Catppuccin Mocha: Rich pastel accents on warm dark crust.",
			"Tokyo Night: High-contrast midnight blue canvas with neon violet and cyan accents.",
			"Nord: Elegant Arctic deep slate palette with frosty blue highlights.",
			"Your active theme is preserved across sessions via localStorage.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "Ctrl+K -> T1", Desc: "Switch to Default Dark Theme"},
			{Key: "Ctrl+K -> T2", Desc: "Switch to Catppuccin Mocha Theme"},
			{Key: "Ctrl+K -> T3", Desc: "Switch to Tokyo Night Theme"},
			{Key: "Ctrl+K -> T4", Desc: "Switch to Nord Theme"},
		},
		Examples: []string{
			"Open Action Palette (Ctrl+K) and select 'Switch to Catppuccin Mocha'",
		},
	},
	{
		ID:       "scripts",
		Title:    "Raycast Script Commands Engine",
		Category: "Extensibility",
		Summary:  "Execute custom Bash, Python, or Node scripts with Raycast metadata headers.",
		Content: []string{
			"Beam automatically discovers executable scripts located in ~/.config/beam/scripts/.",
			"Supports standard Raycast metadata comments (@raycast.title, @raycast.mode, @raycast.description, @raycast.author).",
			"Scripts run in safe detached execution with timeouts and output capture.",
			"Any executable file (chmod +x) with a valid @raycast.title will appear directly in your launcher search results.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "↵ Enter", Desc: "Execute the selected script command"},
		},
		Examples: []string{
			"Place a script in ~/.config/beam/scripts/ip.sh to check your public IP",
			"Use @raycast.mode compact or fullOutput to control display",
		},
	},
	{
		ID:       "shortcuts",
		Title:    "Master Keyboard Shortcuts Cheat Sheet",
		Category: "Reference",
		Summary:  "Comprehensive list of all keyboard controls and shortcuts in Beam.",
		Content: []string{
			"Beam is completely navigable without ever touching the mouse.",
			"Here is the master quick-reference for all primary key combinations.",
		},
		Shortcuts: []GuideShortcut{
			{Key: "Super+Space", Desc: "Toggle Beam spotlight window"},
			{Key: "Tab", Desc: "Cycle through active modes"},
			{Key: "Ctrl+K", Desc: "Open Action Palette"},
			{Key: "↵ Enter", Desc: "Launch app / Copy clipboard / Toggle task / Copy emoji"},
			{Key: "Space", Desc: "Toggle To-Do task completion"},
			{Key: "Ctrl+D", Desc: "Delete selected task in To-Do mode"},
			{Key: "↑ / ↓", Desc: "Navigate items up or down (1 row in emoji grid)"},
			{Key: "← / →", Desc: "Navigate items left or right (in emoji grid)"},
			{Key: "PageUp / PageDown", Desc: "Scroll 5 items in list or 3 rows in grid"},
			{Key: "Home / End", Desc: "Jump to first or last item"},
			{Key: "Esc", Desc: "Clear query, close palette, or hide Beam window"},
		},
		Examples: []string{
			"Press Tab to cycle between all launcher features",
			"Use Ctrl+K for all contextual options",
		},
	},
}

// GetGuideTopics returns the complete catalog of guide topics.
func GetGuideTopics() []GuideTopic {
	return guideCatalog
}

// SearchGuideTopics searches the guide catalog by title, category, summary, content, and shortcuts.
func SearchGuideTopics(rawQuery string) []GuideTopic {
	q := strings.ToLower(strings.TrimSpace(rawQuery))
	if q == "" {
		return guideCatalog
	}

	var results []GuideTopic
	for _, topic := range guideCatalog {
		if strings.Contains(strings.ToLower(topic.Title), q) ||
			strings.Contains(strings.ToLower(topic.Category), q) ||
			strings.Contains(strings.ToLower(topic.Summary), q) {
			results = append(results, topic)
			continue
		}

		matchedContent := false
		for _, line := range topic.Content {
			if strings.Contains(strings.ToLower(line), q) {
				results = append(results, topic)
				matchedContent = true
				break
			}
		}
		if matchedContent {
			continue
		}

		for _, ex := range topic.Examples {
			if strings.Contains(strings.ToLower(ex), q) {
				results = append(results, topic)
				break
			}
		}
	}

	return results
}
