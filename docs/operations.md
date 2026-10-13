# ⚙️ Beam Operations & Developer Guide

## 1. System Requirements & Dependencies

### Arch Linux
```bash
sudo pacman -S go nodejs npm pkgconf gtk3 webkit2gtk-4.1 wl-clipboard
```

### Ubuntu / Debian (22.04+)
```bash
sudo apt update
sudo apt install golang-go nodejs npm pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev wl-clipboard xclip
```

### Fedora
```bash
sudo dnf install golang nodejs npm pkgconf-pkg-config gtk3-devel webkit2gtk4.1-devel wl-clipboard
```

---

## 2. Local Development

### First-Time Project Setup
```bash
cd /home/dilkhush/Desktop/Mine/beam
cd frontend && npm install && cd ..
```

### Running with Live Reload
Run Wails development server. Any changes to Go backend files will trigger a background rebuild, and Svelte frontend changes reload instantly via Vite HMR:
```bash
wails dev -tags webkit2_41
```

---

## 3. Production Build & Optimization

To produce a production-ready, stripped standalone binary:

```bash
# Optimized build with stripped debug symbols
wails build -tags webkit2_41 -ldflags="-s -w"
```

The compiled binary will be generated at:
```bash
build/bin/beam
```

Verify binary size and libraries:
```bash
ls -lh build/bin/beam
ldd build/bin/beam | grep -E "webkit|gtk"
```

---

## 4. System-Wide Installation & Desktop Integration

### 1. Install Binary to User Path
```bash
mkdir -p ~/.local/bin
cp build/bin/beam ~/.local/bin/beam
chmod +x ~/.local/bin/beam
```

### 2. Create Desktop Entry
Create `~/.local/share/applications/beam.desktop`:
```ini
[Desktop Entry]
Name=Beam
Comment=Raycast-like launcher for Linux
Exec=/home/dilkhush/.local/bin/beam
Icon=utilities-terminal
Terminal=false
Type=Application
Categories=Utility;System;
StartupNotify=false
```

### 3. Global Hotkey Configuration (GNOME Wayland)
1. Open **GNOME Settings** (`gnome-control-center`).
2. Navigate to **Keyboard** → **Keyboard Shortcuts** → **View and Customize Shortcuts** → **Custom Shortcuts**.
3. Click **Add Shortcut (+)**:
   * **Name:** `Beam Launcher`
   * **Command:** `/home/dilkhush/.local/bin/beam`
   * **Shortcut:** `Alt + Space` (or `Super + Space`)

---

## 5. Troubleshooting & FAQ

### Issue: "Required dependencies missing: libwebkit"
* **Cause:** Wails looks for `webkit2gtk-4.0.pc` by default, but modern distributions (Arch, Fedora) ship `webkit2gtk-4.1`.
* **Fix:** Always provide `-tags webkit2_41` to `wails build` and `wails dev`.

### Issue: Clipboard history not updating on Wayland
* **Cause:** Missing Wayland clipboard CLI tool.
* **Fix:** Ensure `wl-clipboard` is installed (`sudo pacman -S wl-clipboard`). Test manually using `wl-paste`.

### Issue: Window not appearing on top in Wayland
* **Cause:** Wayland compositors control window stacking strictly.
* **Fix:** In `main.go`, ensure `AlwaysOnTop: true` and `DisableResize: true` are set. For GNOME, `HideWindowOnClose` allows the process to stay alive and re-focus instantly.
