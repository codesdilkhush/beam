# ⚙️ Beam Operations & Developer Guide (Bun & Tailwind CSS)

## 1. System Requirements & Dependencies

### Arch Linux
```bash
sudo pacman -S go bun pkgconf gtk3 webkit2gtk-4.1 wl-clipboard
```

### Ubuntu / Debian (22.04+)
```bash
# Install Bun
curl -fsSL https://bun.sh/install | bash

# Install system dependencies
sudo apt update
sudo apt install golang-go pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev wl-clipboard xclip
```

### Fedora
```bash
curl -fsSL https://bun.sh/install | bash
sudo dnf install golang pkgconf-pkg-config gtk3-devel webkit2gtk4.1-devel wl-clipboard
```

---

## 2. Local Development Workflow (Bun-Powered)

Beam uses **Bun** as its package manager and frontend build runner.

### Initializing Frontend Dependencies
```bash
cd /home/dilkhush/Desktop/Mine/beam/frontend
bun install
```

### Running with Hot Reload (Wails Dev)
Wails monitors both the Go backend and Svelte frontend. Changes to `.go` files trigger background rebuilds, while Svelte and Tailwind CSS changes update instantly via Vite HMR:

```bash
cd /home/dilkhush/Desktop/Mine/beam
wails dev -tags webkit2_41
```

### Frontend-Only Development (Browser Mode)
To work strictly on UI layouts with fast browser refresh:
```bash
cd /home/dilkhush/Desktop/Mine/beam/frontend
bun run dev
```

---

## 3. Tailwind CSS v4 Configuration

Tailwind CSS v4 is integrated directly via `@tailwindcss/vite` in `frontend/vite.config.ts`:

```typescript
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    svelte(),
    tailwindcss(),
  ],
})
```

And imported at the root of `frontend/src/style.css`:
```css
@import "tailwindcss";
```

---

## 4. Production Build & Optimization

To compile the final standalone native binary:

```bash
cd /home/dilkhush/Desktop/Mine/beam
wails build -tags webkit2_41 -ldflags="-s -w"
```

The resulting binary will be output to `build/bin/beam` (~9 MB).

---

## 5. Desktop Shortcut Integration (GNOME Wayland)

1. Open **GNOME Settings** → **Keyboard** → **View and Customize Shortcuts** → **Custom Shortcuts**.
2. Click **+** to add a new shortcut:
   * **Name:** `Beam Launcher`
   * **Command:** `/home/dilkhush/Desktop/Mine/beam/build/bin/beam`
   * **Shortcut:** `Alt + Space` or `Super + Space`
