# ⚡ Beam

A modern, fast, and keyboard-centric launcher for Linux inspired by **Raycast**, built with **Go + Wails + Svelte**.

---

## ✨ Features

- **🚀 Application Launcher**: Instantly searches and launches installed `.desktop` apps with icon resolution.
- **📋 Live Clipboard Manager**: Tracks your clipboard history in real-time (Wayland & X11). Type `/clip` or press `Tab` to search through copied items.
- **🧮 Inline Math Calculator**: Evaluates arithmetic expressions directly in the search bar (e.g., `45 * 12`, `sqrt(144) + 10`, `sin(pi/2)`). Press `Enter` to copy results.
- **🌐 Quick Web Fallback**: Seamlessly search Google, GitHub, or open URLs when no local app matches.
- **⌨️ Keyboard-Centric Workflow**:
  - `↑` / `↓`: Navigate items
  - `↵ Enter`: Open app or copy result/clipboard
  - `Tab`: Switch between modes (All / Clipboard / Apps)
  - `Esc`: Clear search query or dismiss launcher

---

## 📂 Project Structure

```
beam/
├── app.go                 # Wails app lifecycle & JS bindings
├── main.go                # Frameless floating window configuration
├── backend/               # Pure Go engine
│   ├── apps.go            # Linux .desktop parser, icon finder & launcher
│   ├── calc.go            # Safe AST math expression evaluator
│   └── clipboard.go       # Background Wayland/X11 clipboard tracker
└── frontend/              # Raycast-inspired Svelte UI
    └── src/
        ├── App.svelte     # Search, split view, keyboard handlers
        └── style.css      # Dark glassmorphic theme
```

---

## 🛠️ Build & Run

### 1. Build the Binary
```bash
wails build -tags webkit2_41
```
The compiled binary will be placed at `build/bin/beam` (~9 MB).

### 2. Run Directly
```bash
./build/bin/beam
```

### 3. Development Mode (Live Reloading)
```bash
wails dev -tags webkit2_41
```

---

## 🚀 Setting a Global Hotkey (e.g., `Super + Space` or `Alt + Space`)

To make Beam behave like Raycast:
1. Open **GNOME Settings** → **Keyboard** → **Keyboard Shortcuts** → **View and Customize Shortcuts** → **Custom Shortcuts**.
2. Click **+** to add a new shortcut:
   - **Name**: `Beam Launcher`
   - **Command**: `/home/dilkhush/Desktop/Mine/beam/build/bin/beam`
   - **Shortcut**: `Alt + Space` or `Super + Space`
