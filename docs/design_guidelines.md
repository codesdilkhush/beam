# 🎨 Beam Design Guidelines & Visual Language

This document defines the visual design system, glassmorphism principles, color palette, typography, and micro-interaction specifications for **Beam**, inspired by the signature aesthetic of **Raycast**.

---

## 1. Visual Identity & Aesthetic Principles

1. **Frosted Glass (Glassmorphism):**
   * Multi-layered translucent surfaces with backdrop blur.
   * Subtle semi-transparent borders giving the appearance of etched glass.
2. **Electric Glass Blue Accent:**
   * High-contrast electric blue (`#3875f6`) used strategically for active selections, interactive focus, and key actions.
3. **Information Density & Razor-Sharp Layout:**
   * Split-pane layout (60% list / 40% details preview).
   * Strict typographic hierarchy with monospaced accents for technical data (keybindings, shortcuts, commands).
4. **Sub-Pixel Polish:**
   * 1-pixel border dividers with delicate opacity gradients (`rgba(255, 255, 255, 0.07)`).
   * Deep diffused ambient drop-shadows to separate the launcher from desktop background content.

---

## 2. Color System

### 2.1 Surfaces & Backgrounds
| Token | Hex / RGBA | Tailwind Equiv. / Description |
| :--- | :--- | :--- |
| `--surface-canvas` | `rgba(16, 18, 24, 0.85)` | Base window background with backdrop blur |
| `--surface-header` | `rgba(22, 25, 35, 0.90)` | Top search bar container |
| `--surface-card` | `rgba(255, 255, 255, 0.04)` | Inset detail panels & secondary cards |
| `--surface-hover` | `rgba(255, 255, 255, 0.07)` | Hover state for unselected list items |
| `--surface-selected`| `rgba(56, 117, 246, 0.16)` | Selected list item with electric blue tint |

### 2.2 Glass Borders & Dividers
| Token | Hex / RGBA | Description |
| :--- | :--- | :--- |
| `--border-subtle` | `rgba(255, 255, 255, 0.08)` | Default card and divider borders |
| `--border-glass` | `rgba(255, 255, 255, 0.14)` | Window outer frame border |
| `--border-focus` | `rgba(56, 117, 246, 0.50)` | Active selection highlight border |

### 2.3 Raycast Accent Palette
| Purpose | Hex | Name | Usage |
| :--- | :--- | :--- | :--- |
| **Primary Accent** | `#3875F6` | **Electric Glass Blue** | Selection pills, active icons, primary enter action |
| **Secondary Accent** | `#60A5FA` | **Sky Frost** | Badges, secondary shortcut indicators |
| **Success / Math** | `#10B981` | **Emerald Glow** | Math calculation results, successful copy feedback |
| **Extension / Tag** | `#A855F7` | **Neon Purple** | Custom scripts, plugins, and category tags |
| **Warning** | `#F59E0B` | **Amber** | Confirmation required, dangerous commands |
| **Danger** | `#EF4444` | **Crimson** | Process termination, clipboard purge |

### 2.4 Typography Colors
| Token | Hex / RGBA | Usage |
| :--- | :--- | :--- |
| `--text-primary` | `#F9FAFB` | App names, search query, primary labels |
| `--text-secondary` | `#9CA3AF` | Subtitles, command flags, descriptions |
| `--text-muted` | `#6B7280` | Shortcut hints, character counts, inactive tabs |

---

## 3. Typography & Hierarchy

* **Primary Font Family:** System Sans-Serif (`Inter`, `system-ui`, `-apple-system`, `BlinkMacSystemFont`, `Segoe UI`, `Roboto`)
* **Monospace Font Family:** `JetBrains Mono`, `Fira Code`, `ui-monospace`, `monospace`

| Element | Size | Weight | Line Height |
| :--- | :--- | :--- | :--- |
| **Search Input** | `19px` | Medium (500) | `26px` |
| **List Item Title** | `14.5px` | Medium (500) | `20px` |
| **List Item Subtitle** | `12px` | Regular (400) | `16px` |
| **Detail Header** | `16px` | SemiBold (600) | `22px` |
| **Detail Body / Code** | `12.5px` | Regular (400) | `18px` |
| **Action Pill (Footer)**| `12px` | SemiBold (600) | `16px` |
| **Key Shortcut Badge** | `11px` | Medium (500) | `14px` |

---

## 4. Spacing, Geometry & Radii

* **Window Dimensions:** `760px` width × `500px` height
* **Window Corner Radius:** `14px` (`rounded-2xl`)
* **List Item Corner Radius:** `8px` (`rounded-lg`)
* **Badge / Key Pill Radius:** `6px` (`rounded-md`)
* **Window Elevation Shadow:**
  ```css
  box-shadow: 0 24px 60px -12px rgba(0, 0, 0, 0.75), 0 0 0 1px rgba(255, 255, 255, 0.1);
  ```

---

## 5. Component Specifications

### 5.1 Search Bar (Header)
* Height: `60px`
* Background: `--surface-header` with `backdrop-filter: blur(16px)`
* Magnifying glass icon: 20px, `--text-muted`
* Search input: Full width, zero border, transparent background, auto-focused.
* Right side contains mode badge pills (e.g. `📋 Clipboard`, `🚀 Apps`) with frosted glass background.

### 5.2 Results List (Left 60%)
* Padding: `8px`
* Row Height: `44px`
* Selection Indicator:
  * Background: `rgba(56, 117, 246, 0.16)`
  * Border: `1px solid rgba(56, 117, 246, 0.35)`
  * Left border accent glow: `2px solid #3875F6`
* Icon: `30px` rounded square with high-res PNG/SVG or letter monogram avatar.

### 5.3 Detail / Context Panel (Right 40%)
* Background: `--surface-card` with left 1px border divider (`--border-subtle`).
* Content Header: Large 44px icon with application title and category pill.
* Metadata Grid: Command line, desktop file name, launch mode in clean key-value rows.
* Clipboard View: Formatted `<pre>` snippet box with character count badge and word wrapping.

### 5.4 Action Bar (Footer)
* Height: `38px`
* Background: `rgba(14, 16, 22, 0.95)`
* Left: Primary selected action pill (e.g., `↵ Open Application` or `↵ Copy to Clipboard`).
* Right: Keyboard hints (`Tab: Mode`, `↑↓: Navigate`, `Esc: Close`).

---

## 6. Micro-Interactions & Animation Standards

* **Selection Transitions:** `100ms ease-out` for row hover and arrow key navigation.
* **Mode Switches:** Immediate (`0ms` state change) without layout shift.
* **Calculation Banner:** Smooth entrance when an arithmetic query is recognized.
