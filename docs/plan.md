# Beam Priority-Based Product & Technical Plan

## 1. Vision & Strategy

Beam delivers the **Raycast experience on Linux** through a tiered, priority-driven execution model. Features are grouped into distinct priority tiers ensuring immediate usability while laying the foundation for advanced extensibility.

---

## 2. Priority Matrix

| Priority | Focus Area | Key Deliverables | Target Latency |
| :--- | :--- | :--- | :--- |
| **P0 (Critical)** | Core Launcher & Window | • Frameless spotlight window<br>• `.desktop` app discovery & launch<br>• Wayland live clipboard tracker<br>• Bun + Tailwind CSS UI | < 50ms |
| **P1 (High)** | Search & Persistence | • `fzf`-style fuzzy scoring & acronym matching<br>• Persistent SQLite database (`beam.db`)<br>• AST math evaluation & converters<br>• High-res icon cache | < 10ms |
| **P2 (Medium)** | Extensibility & Daemon | • UNIX socket single-instance daemon (`beam --toggle`)<br>• Raycast-compatible script commands<br>• Image/screenshot clipboard storage | < 5ms (toggle) |
| **P3 (Enhancement)**| Customization & Polish | • `Ctrl+K` action palette<br>• Community themes (Catppuccin, Tokyo Night)<br>• Preferences configuration UI | Instant |

---

## 3. Detailed Priority Specifications

### P0: Core Operational Foundation (Completed Baseline)
* **Goal:** A usable keyboard-centric application runner and live clipboard inspector.
* **Requirements:**
  * Clean parsing of Linux `.desktop` files in system and local paths.
  * Reliable detached execution of GUI and CLI apps.
  * Live monitoring of Wayland clipboard via `wl-paste`.
  * Bun-driven frontend bundling with Svelte 5 and Tailwind CSS v4.
  * Instant keyboard controls (`↑`/`↓`, `↵`, `Tab`, `Esc`).

### P1: Intelligence, Fuzzy Search & Persistent Storage
* **Goal:** Zero-friction typing and historical clipboard recovery.
* **Requirements:**
  * Replace simple substring matching with a scoring algorithm that rewards:
    * Exact prefix matches.
    * Acronyms (`vsc` for Visual Studio Code, `gc` for Google Chrome).
    * Word boundaries.
  * Local SQLite storage for clipboard entries with FTS5 search.
  * Real-time calculation banner with one-click copy.
  * Icon caching in `~/.cache/beam/icons/` to eliminate decoding overhead.

### P2: Desktop Integration & Script Commands
* **Goal:** Raycast-grade instant invocation and power-user customization.
* **Requirements:**
  * Background daemon with UNIX domain socket communication (`/tmp/beam.sock`).
  * Instant window toggling via `beam --toggle` without re-initializing the runtime.
  * Script runner: scan `~/.config/beam/scripts/` for scripts declaring `@raycast.*` metadata.
  * Clipboard support for binary image buffers (screenshots).

### P3: Actions Palette & Theming
* **Goal:** Contextual deep actions and visual personalization.
* **Requirements:**
  * `Ctrl+K` popover showing secondary actions for selected item.
  * Theme switcher powered by Tailwind CSS classes.
  * User settings stored in `~/.config/beam/config.json`.

---

## 4. Technology Stack Specification

* **Backend Engine:** Go 1.22+
* **Frontend Runtime:** Bun 1.4+
* **Frontend Framework:** Svelte 5 (Runes)
* **CSS Framework:** Tailwind CSS v4 (`@tailwindcss/vite`)
* **Bundler:** Vite 7
* **Desktop Runtime:** Wails v2 with WebKit2GTK 4.1
