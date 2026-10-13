# 📋 Beam Product & Technical Plan

## 1. Vision & Core Principles

**Beam** is designed to bring the polished, instantaneous, and keyboard-driven productivity workflow of **Raycast** to Linux desktop environments (Wayland & X11).

### Core Design Principles
1. **Sub-Millisecond Responsiveness:** Zero perceived latency on keystrokes, filter updates, and navigation.
2. **Keyboard-Driven:** Every single action must be executable without touching the mouse.
3. **Low Footprint:** Native compiled binary with minimal idle CPU and <50 MB memory usage.
4. **Extensibility:** Users should be able to write custom scripts in any language (Bash, Python, Node.js) and see rich UI results.

---

## 2. Core Feature Pillars

### Pillar 1: Application Launcher & Fuzzy Search
* Parse all installed applications from system and user directories.
* Advanced fuzzy ranking:
  * Exact prefix match on application name gets highest priority.
  * Acronym matching (e.g. `gc` -> `Google Chrome`, `vsc` -> `Visual Studio Code`).
  * Secondary match on description, keywords, and executable name.
* Recency & frequency boosting (MRU/MFU scoring).

### Pillar 2: Persistent Clipboard Manager
* Live background capture for text, code snippets, and URLs.
* Local persistent storage using **SQLite** with encrypted/private item exclusion options.
* Full-text search across clipboard history with syntax highlighting on code snippets.
* Paging and cleanup policies (e.g., retain last 500 items or 30 days).

### Pillar 3: Calculations, Conversions & Tools
* **Inline Math:** Instant evaluation without needing to hit Enter first.
* **Unit & Currency Conversion:** `100 USD in EUR`, `50 km in miles`, `16 GB in MB`.
* **Date & Timezone:** `time in tokyo`, `utc now`, `timestamp to date`.

### Pillar 4: Script Commands & Raycast Extension Compatibility
* Support Raycast's **Script Command** specification:
  * Scripts with metadata headers (e.g., `@raycast.title`, `@raycast.mode`, `@raycast.icon`).
  * Modes:
    * `compact`: Runs command and copies/shows output.
    * `fullOutput`: Renders terminal output in the detail pane.
    * `inline`: Renders JSON list items in the main results pane.
* Scripts can be written in Bash, Python, Go, Node.js, or any executable format.

### Pillar 5: System Commands
* Direct execution of desktop commands:
  * Lock Screen (`loginctl lock-session`)
  * Suspend / Shutdown / Restart
  * Empty Trash
  * Toggle Dark/Light Mode

---

## 3. UI/UX Specification

* **Window Dimensions:** 760px width × 500px height (floating centered).
* **Search Bar:** Large typography (19px), clear icon, dynamic mode indicators.
* **Results Pane (60%):** Icon + Title + Subtitle + Action Badge.
* **Detail / Context Pane (40%):**
  * Application: Desktop metadata, command, icon, category badges.
  * Clipboard: Complete snippet preview with line numbers and character count.
  * Calculation: Step-by-step breakdown or currency rates.
* **Action Footer:** Left-aligned primary action label (`↵ Open`), right-aligned keyboard shortcut hints (`Tab`, `↑↓`, `Esc`).

---

## 4. Technical Constraints & Standards

* **Platform:** Linux (Targeting Arch, Fedora, Ubuntu/Debian).
* **Display Server:** Wayland-native first, full X11 compatibility.
* **Language/Stack:** Go 1.22+ backend, Wails v2 runtime, Svelte 5 frontend, TypeScript.
* **Build Targets:** Single executable binary distribution (`build/bin/beam`).
