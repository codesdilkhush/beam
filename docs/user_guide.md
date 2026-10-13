# Beam — Comprehensive User Guide & Manual

Beam is an ultra-fast, keyboard-driven application launcher and productivity spotlight for Linux. Crafted with **Go**, **Wails v2**, and **Svelte 5 Runes**, Beam delivers a premium glassmorphic interface inspired by Raycast and macOS Spotlight, engineered specifically for modern Linux desktop environments.

---

## Table of Contents

1. [Architecture & Philosophy](#1-architecture--philosophy)
2. [Prerequisites & Installation](#2-prerequisites--installation)
   - [System Dependencies](#system-dependencies)
   - [Building Beam](#building-beam)
   - [Running as a Background Daemon](#running-as-a-background-daemon)
3. [Global Hotkey Setup](#3-global-hotkey-setup)
   - [GNOME](#gnome-wayland--x11)
   - [KDE Plasma](#kde-plasma)
   - [Hyprland](#hyprland)
   - [Sway](#sway)
   - [i3 Window Manager](#i3-window-manager)
4. [Core Features & Usage](#4-core-features--usage)
   - [Application Launcher](#application-launcher)
   - [Clipboard History Manager (`/clip`)](#clipboard-history-manager-clip)
   - [To-Do & Task Manager (`/todo`)](#to-do--task-manager-todo)
   - [Emoji Picker & Grid (`/emoji`, `:tag`)](#emoji-picker--grid-emoji-tag)
   - [Universal Converter & Calendar](#universal-converter--calendar)
   - [Math Expression Evaluator](#math-expression-evaluator)
   - [Raycast-Style Action Palette (`Ctrl+K`)](#raycast-style-action-palette-ctrlk)
   - [Themes & Customization](#themes--customization)
   - [Custom Script Commands](#custom-script-commands)
   - [In-App User Guide (`/guide`, `/help`, `/?`)](#in-app-user-guide-guide-help-)
5. [Master Keyboard Shortcuts](#5-master-keyboard-shortcuts)
6. [Data Storage & Directory Structure](#6-data-storage--directory-structure)
7. [Troubleshooting & FAQ](#7-troubleshooting--faq)

---

## 1. Architecture & Philosophy

Beam is architected on three core principles:
* **Sub-Millisecond Responsiveness:** Written in Go with an event-driven IPC daemon (`beam-daemon.sock`), Beam stays resident in memory. When summoned with a global hotkey, the window appears instantly without process initialization latency.
* **Pure Keyboard Control:** Every feature—from application launching and clipboard restoration to task completion and theme switching—is fully navigable without a mouse.
* **Zero-Emoji UI Chrome:** In accordance with Beam's design guidelines, the launcher chrome, buttons, tab badges, and status panels use strict, pixel-perfect vector SVGs. Emojis are strictly handled as data within the dedicated Emoji utility.

### Technical Stack
* **Core Daemon & Backend:** Go 1.22+, `modernc.org/sqlite` (pure Go, CGO-free SQLite3).
* **Desktop UI Shell:** Wails v2 (WebKitGTK on Linux).
* **Frontend Framework:** Svelte 5 with reactive Runes (`$state`, `$derived`, `$effect`).
* **Package Manager & Bundler:** Bun + Vite.
* **Style Engine:** Native scoped CSS with translucent glassmorphism tokens, backdrop filters, and CSS variables.

---

## 2. Prerequisites & Installation

### System Dependencies

Beam requires the following packages installed on your Linux distribution:

#### Debian / Ubuntu / Pop!_OS
```bash
sudo apt update
sudo apt install -y build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev \
    wl-clipboard xclip sqlite3 curl git
```

#### Arch Linux / Manjaro / EndeavourOS
```bash
sudo pacman -Syu --needed base-devel pkgconf gtk3 webkit2gtk-4.1 \
    wl-clipboard xclip sqlite curl git
```

#### Fedora / RHEL
```bash
sudo dnf groupinstall -y "Development Tools"
sudo dnf install -y pkg-config gtk3-devel webkit2gtk4.1-devel \
    wl-clipboard xclip sqlite curl git
```

#### Runtime & Toolchains
1. **Go 1.22+**: [golang.org/dl](https://golang.org/dl/)
2. **Bun**:
   ```bash
   curl -fsSL https://bun.sh/install | bash
   ```
3. **Wails CLI v2**:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
   Ensure `~/go/bin` is in your `$PATH`.

---

### Building Beam

Clone the repository and compile using Wails and Bun:

```bash
# Clone Beam
git clone https://github.com/dilkhush-raj/beam.git
cd beam

# Install frontend dependencies using Bun
cd frontend
bun install
cd ..

# Build the release binary
wails build -clean -s
```

The compiled binary will be placed at `build/bin/beam`.

To install Beam system-wide or for your user:
```bash
mkdir -p ~/.local/bin
cp build/bin/beam ~/.local/bin/beam
chmod +x ~/.local/bin/beam
```

---

### Running as a Background Daemon

Beam includes a background daemon mode with a Unix domain socket server (`/tmp/beam-daemon.sock` or `~/.local/share/beam/beam.sock`).

To start Beam in daemon mode automatically upon desktop login:

1. Create a systemd user service:
```ini
# ~/.config/systemd/user/beam.service
[Unit]
Description=Beam Application Launcher Daemon
After=graphical-session.target

[Service]
Type=simple
ExecStart=%h/.local/bin/beam --daemon
Restart=on-failure
RestartSec=3

[Install]
WantedBy=default.target
```

2. Enable and start the service:
```bash
systemctl --user daemon-reload
systemctl --user enable --now beam.service
```

Alternatively, you can add `beam --daemon` to your desktop environment's autostart applications or your window manager config.

---

## 3. Global Hotkey Setup

To toggle Beam from any application, bind your preferred shortcut (commonly `Super + Space` or `Alt + Space`) to invoke `beam --toggle`.

### GNOME (Wayland & X11)
1. Open **Settings** → **Keyboard** → **View and Customize Shortcuts** → **Custom Shortcuts**.
2. Click **+** to add a new shortcut:
   - **Name:** Beam Launcher
   - **Command:** `beam --toggle`
   - **Shortcut:** `Super + Space` (or `Alt + Space`)

*(If Super+Space is used by input source switching, reassign "Switch input source" in Settings → Keyboard).*

### KDE Plasma
1. Open **System Settings** → **Shortcuts** → **Custom Shortcuts**.
2. Click **Edit** → **New** → **Global Shortcut** → **Command/URL**:
   - **Name:** Beam Launcher
   - **Trigger:** `Meta + Space`
   - **Action:** `beam --toggle`

### Hyprland
Add the following line to `~/.config/hypr/hyprland.conf`:
```ini
# Toggle Beam Spotlight
bind = SUPER, SPACE, exec, beam --toggle
```

### Sway
Add the following line to `~/.config/sway/config`:
```ini
# Toggle Beam Spotlight
bindsym $mod+space exec beam --toggle
```

### i3 Window Manager
Add the following line to `~/.config/i3/config`:
```ini
# Toggle Beam Spotlight
bindsym $mod+space exec beam --toggle
```

---

## 4. Core Features & Usage

### Application Launcher
Beam indexes system applications from XDG desktop entry directories (`/usr/share/applications`, `~/.local/share/applications`, and Flatpak paths).

* **Fuzzy Search:** Type partial words (e.g., `term` finds Alacritty, GNOME Terminal, Kitty).
* **Acronym Matching:** Type initials to jump straight to tools (e.g., `gc` finds Google Chrome, `vs` finds Visual Studio Code, `mpv` finds MPV).
* **Adaptive Ranking:** Applications you launch frequently are automatically boosted to the top of results.
* **Terminal Application Detection:** CLI applications configured with `Terminal=true` in their `.desktop` files automatically open inside your system terminal emulator.

---

### Clipboard History Manager (`/clip`)
Beam monitors your Wayland/X11 clipboard in real-time, persisting history securely in an embedded SQLite database.

* **Trigger:** Type `/clip` into the search bar, or press `Tab` until the **Clipboard** mode badge activates.
* **Search History:** Type keywords to filter through copied snippets.
* **Restore & Paste:** Press `Enter` on any item to write it to your system clipboard and instantly close Beam.
* **Binary Safety:** Beam filters out heavy binary data, preserving clean text snippets, code, and URLs.

---

### To-Do & Task Manager (`/todo`)
A lightweight, keyboard-driven checklist that lives directly inside your spotlight launcher.

* **Trigger:** Type `/todo` or `/t `, or press `Tab` to cycle into **To-Do** mode.
* **Quick-Add Task:** While in To-Do mode, type your task title and hit `Enter` (e.g., `/t Finish quarterly review` + `↵`).
* **Toggle Task Completion:** Highlight a task using `↑` / `↓` and press `Space` or `Enter`.
* **Delete Task:** Press `Ctrl+D` or `Delete` on an empty search prompt to delete the selected task.
* **Clear Completed:** Open the Action Palette (`Ctrl+K`) and select **Clear Completed Tasks** (`Shift+Del`).

---

### Emoji Picker & Grid (`/emoji`, `:tag`)
An 8-column visual emoji grid providing fast character discovery with zero UI lag.

* **Trigger:**
  - Prefix with colon: `:rocket`, `:sparkles`, `:fire`
  - Slash command: `/emoji` or `/e `
  - Press `Tab` until the **Emoji** mode badge activates.
* **2D Grid Navigation:**
  - `←` / `→`: Move left or right one column.
  - `↑` / `↓`: Move up or down across rows (jumps 8 items).
  - `Page Up` / `Page Down`: Jump 3 full rows at a time.
* **Instant Insert:** Press `Enter` to copy the selected Unicode emoji to your clipboard.
* **Detail Inspector:** The right-hand panel displays the Unicode codepoint, official CLDR name, and search tags.

---

### Universal Converter & Calendar
Beam features an intelligent calculation and conversion parser that evaluates queries directly from the main search bar without requiring any mode switch:

* **Physical Units:**
  - `10 km to mi` → `6.2137 mi`
  - `100 m in ft` → `328.0840 ft`
  - `50 kg to lb` → `110.2311 lb`
  - `16 gb in mb` → `16,384 MB`
  - `100 km/h to mph` → `62.1371 mph`
* **Temperature:**
  - `100 c to f` → `212.0000 °F`
  - `32 f to c` → `0.0000 °C`
* **Color Code Translations:**
  - `#3875f6` → `rgb(56, 117, 246) | hsl(221, 91%, 59%)`
  - `rgb(56, 117, 246)` → `#3875f6 | hsl(221, 91%, 59%)`
* **CSS Units:**
  - `16px to rem` → `1.0000 rem (16px base)`
  - `2rem to px` → `32.0000 px (16px base)`
* **Currencies:**
  - `100 usd to eur`
  - `50 eur to gbp`
* **Calendar, Date & Time:**
  - `days in march` → `31 days (March)`
  - `days in current month` → Dynamic days count for this month
  - `what is the weekday?` or `weekday` → e.g., `Wednesday`
  - `weekday of 2026-10-14` → Weekday for any ISO date
  - `today` → e.g., `Wednesday, Oct 14, 2026`
  - `week number` → Current ISO week number
  - `day of year` → e.g., `Day 287 of 365`
  - `days in 2024` → `366 days (2024 is a leap year)`
* **Unix Timestamps:**
  - `now to epoch` → e.g., `1773456000`
  - `1700000000 to date` → `Wednesday, Nov 15, 2023 22:13:20 UTC`

Press `Enter` while a conversion banner is active to copy the formatted result immediately.

---

### Math Expression Evaluator
Type any mathematical equation into the search bar:
* Supported operators: `+`, `-`, `*`, `/`, `%`, `^`, `(`, `)`
* Examples:
  - `(128 * 4) + 512` → `1024`
  - `2^16` → `65536`
  - `15% of 850` → `127.5`
* Press `Enter` when the math preview banner appears to copy the result to your clipboard.

---

### Raycast-Style Action Palette (`Ctrl+K`)
Whenever an item is selected or you need secondary commands, press `Ctrl+K` to reveal the contextual Action Palette.

* **Contextual Options:**
  - Applications: Launch, Copy Exec Path, Open in Terminal.
  - Clipboard: Copy, Remove from History.
  - To-Do: Mark Complete/Incomplete, Delete Task, Clear Completed Tasks.
  - Universal: Switch themes, Open User Guide, Quit Beam.
* **Filter Actions:** Type directly inside the palette to filter available actions.
* **Escape:** Dismiss the Action Palette with `Escape`.

---

### Themes & Customization
Beam features built-in glassmorphism themes with translucent materials, subtle borders, and saturated backdrops:

1. **Default Dark Glass:** Modern deep slate with neutral glass highlights.
2. **Catppuccin Mocha:** Warm pastel accents with soothing mauve and lavender tones.
3. **Tokyo Night:** Deep cyber dark with neon cyan, blue, and purple tones.
4. **Nord:** Elegant arctic ice and frosted slate palette.

#### Switching Themes:
- Press `Ctrl+K` → select any theme (`Switch to Catppuccin Mocha`, etc.).
- Or press `T1` (Default), `T2` (Catppuccin), `T3` (Tokyo Night), `T4` (Nord) in the Action Palette.
- Your selection is automatically persisted across app restarts.

---

### Custom Script Commands
Extend Beam with your own executable scripts written in Bash, Python, Node.js, or Go.

#### Script Directory:
```bash
mkdir -p ~/.config/beam/scripts/
```

#### Metadata Header Format:
Add metadata comments to the top of your script file:
```bash
#!/usr/bin/env bash
# @name System Monitor
# @description Launch btop in terminal
# @icon utilities-system-monitor
# @terminal true

btop
```

Make your script executable:
```bash
chmod +x ~/.config/beam/scripts/sysmon.sh
```

Beam automatically discovers and lists scripts alongside system applications.

---

### In-App User Guide (`/guide`, `/help`, `/?`)
Beam provides an interactive, split-pane reference guide built right into the application.

* **Trigger:** Type `/guide`, `/help`, or `/?`, or select **Open In-App User Guide** from `Ctrl+K`.
* **Split-Pane Browser:**
  - **Left Pane:** Categorized topic list (Getting Started, Launcher, Clipboard, Tasks, Emoji, Converter, Actions, Themes, Scripts, Shortcuts).
  - **Right Pane:** Detailed instructions, syntax examples, and interactive `<kbd>` shortcut cheat sheets.
* **Click-to-Copy:** Click any example query in the guide to copy it directly to your clipboard.
* **Keyboard Browsing:** Use `↑` / `↓` to move through topics, and `Page Up` / `Page Down` to scroll quickly.

---

## 5. Master Keyboard Shortcuts

| Shortcut | Context | Action |
| :--- | :--- | :--- |
| `Super + Space` | Global | Summon / Dismiss Beam Spotlight |
| `↑` / `↓` | Any List | Move selection up or down |
| `←` / `→` | Emoji Grid | Move left or right across columns |
| `Page Up` / `Page Down` | Any List | Scroll by page / multi-row jump |
| `Home` / `End` | Any List | Jump to first / last item |
| `Tab` | Search Bar | Cycle between modes (All → Clipboard → To-Do → Emoji → Guide → Apps) |
| `↵` (Enter) | Launcher | Open selected application or run URL search |
| `↵` (Enter) | Clipboard | Copy selected clip to clipboard and hide |
| `↵` (Enter) | Calculator / Converter | Copy calculation result to clipboard |
| `↵` (Enter) | Emoji Grid | Copy selected emoji to clipboard |
| `Space` / `↵` | To-Do | Toggle task completion |
| `Ctrl + D` | To-Do | Delete selected task |
| `Ctrl + K` | Any Mode | Open / close contextual Action Palette |
| `Escape` | Search Bar | Clear search text; if already empty, hide Beam |
| `Escape` | Action Palette | Close Action Palette and return to spotlight |

---

## 6. Data Storage & Directory Structure

Beam keeps all user data and configurations neatly organized following XDG specifications:

```text
~/.local/share/beam/
├── beam.db          # SQLite3 database (clipboard history, to-do items, usage stats)
└── beam.sock        # Local IPC daemon socket (when running daemon mode)

~/.config/beam/
├── config.json      # User preferences and keybindings
└── scripts/         # User-defined executable scripts (*.sh, *.py, etc.)

~/.cache/beam/
└── icons/           # Cached rasterized desktop icons (.png)
```

---

## 7. Troubleshooting & FAQ

### Q: Beam opens slowly on Wayland.
**A:** Ensure `beam --daemon` is running in your background session. When daemon mode is active, Beam keeps its WebKitGTK instance warm in memory, allowing window visibility to toggle in less than 20 milliseconds.

### Q: Application icons are missing or generic.
**A:** Beam reads icon themes from `/usr/share/icons` and `~/.local/share/icons`. Ensure an XDG-compliant icon theme (such as Papirus, Adwaita, or Breeze) is installed on your system.

### Q: Clipboard monitoring is not picking up copied items.
**A:** On Wayland, verify `wl-clipboard` (`wl-paste`) is installed. On X11, ensure `xclip` is installed. Run `wl-paste` or `xclip -o` in your terminal to verify clipboard accessibility.

### Q: How do I backup or reset my tasks and clipboard history?
**A:** All persistent records reside in `~/.local/share/beam/beam.db`. To backup, copy the file:
```bash
cp ~/.local/share/beam/beam.db ~/backup-beam.db
```
To start fresh, simply delete the file while Beam is closed; Beam will automatically recreate a clean database on next launch.

---

*Beam is open-source software crafted for productivity and speed.*
