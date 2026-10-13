# Beam Engineering Rules & Coding Standards

This document establishes the mandatory development rules, coding standards, and architectural conventions for **Beam**. All contributors and AI coding agents must strictly adhere to these rules.

---

## 1. Core Philosophy: Simplicity over Cleverness (KISS & YAGNI)

- **Write Simple, Readable Code**: Code must be explicit, easy to follow, and straightforward. Avoid clever tricks, obscure language features, or complex metaprogramming.
- **No Premature Abstraction**: Do not create generic layers, abstract factory patterns, or speculative utility functions for hypothetical future use cases. Solve the immediate problem directly.
- **Keep Functions Small and Focused**: Each function or method should do one thing clearly. Avoid massive multi-hundred-line functions with deeply nested branches.
- **Prefer Flat Hierarchies**: Keep package and folder trees flat. Deep module nesting complicates imports and obfuscates logic.
- **Self-Documenting Code**: Prefer clear variable and function names over excessive comments. Comments should explain *why* something unusual is done, not *what* basic code is doing.

---

## 2. Strict Linting & Zero-Bypass Policy

> [!IMPORTANT]
> **Zero Tolerance for Lint Bypasses.** Never bypass, suppress, or disable compiler or linter checks.

- **No Suppression Directives**:
  - **Go**: Never use `//nolint` or dummy discards like `_ = unusedVar` to silence warnings.
  - **TypeScript/Svelte**: Never use `@ts-ignore`, `@ts-expect-error`, or `// eslint-disable` comments.
- **Fix the Root Cause**: If a linter or compiler warns about a type mismatch, unhandled error, or unused variable, fix the underlying code rather than silencing the tool.
- **Treat Warnings as Errors**: Code with active warnings must not be committed or merged.
- **Strict Go Error Handling**:
  - Always check returned `error` values.
  - Never swallow errors silently. Return them, wrap them with descriptive context (`fmt.Errorf("...: %w", err)`), or log them appropriately.
- **Strict TypeScript Types**:
  - Avoid `any`. Use explicit types, interfaces, or type narrowing.
  - Ensure `svelte-check` passes with zero errors and zero warnings.

---

## 3. Dead and Unused Code Policy (Purge Before Push)

> [!CAUTION]
> Dead code increases cognitive load and causes technical debt. Clean your workspace before committing and pushing.

- **Remove All Dead Code Before Commits & Pushes**:
  - Delete unused functions, methods, structs, types, interfaces, and variables.
  - Delete unreferenced imports.
  - Remove unreachable branches, dead switch cases, and obsolete feature flags.
  - Remove orphaned components, styles, and assets in `frontend/`.
- **No "Zombie" or Commented-Out Code**:
  - Never leave commented-out blocks of code in files.
  - Git history preserves previous implementations; do not keep old code around "just in case".
- **No Stale Debug Statements**:
  - Remove temporary `fmt.Println()`, `console.log()`, `debugger;`, or ad-hoc print statements before committing.

---

## 4. Mandatory Pre-Commit & Pre-Push Workflow

Every commit and push must satisfy formatting, linting, and testing verification.

### Backend Verification (Go)
Run from the repository root:
```bash
# 1. Format code according to official Go standard
gofmt -s -w .

# 2. Run static analysis and vet
go vet ./...

# 3. Run unit tests
go test -v ./backend/...
```

### Frontend Verification (Bun + Svelte 5 + TypeScript)
Run from the `frontend/` directory:
```bash
cd frontend

# 1. Verify Svelte & TypeScript diagnostics (must report 0 errors, 0 warnings)
bun run check

# 2. Verify Vite production bundle builds cleanly
bun run build

cd ..
```

### Clean Git Status Before Push
Before pushing to GitHub:
```bash
git status
```
- Ensure no accidental temporary files, scratch scripts, or unintended artifacts are staged.
- Ensure untracked files are either added deliberately or excluded in `.gitignore`.

---

## 5. Agent & Contributor Guidelines (AI Agent Rules)

All automated agents working on Beam must strictly abide by these conventions:

### A. Folder Architecture & Structural Discipline
Respect and maintain the established project layout. Never introduce random directories or scatter files across the workspace.

```text
beam/
├── app.go                 # Wails App bridge & lifecycle bindings
├── main.go                # Application entrypoint & Wails window config
├── wails.json             # Wails project configuration
├── Makefile               # Standardized build and maintenance tasks
├── backend/               # Go backend engines & subsystems
│   ├── apps.go            # Linux .desktop launcher & icon parser
│   ├── apps_test.go       # Applications scanner unit tests
│   ├── calc.go            # Calculator & AST evaluator
│   ├── calc_test.go       # Calculator unit tests
│   ├── clipboard.go       # System clipboard monitor & history
│   ├── clipboard_test.go  # Clipboard manager unit tests
│   ├── daemon.go          # Single-instance UNIX socket server & client
│   ├── daemon_test.go     # Daemon socket unit tests
│   ├── fuzzy.go           # High-speed fuzzy search algorithm
│   ├── fuzzy_test.go      # Fuzzy scoring & acronym unit tests
│   ├── icon_cache.go      # Two-tier memory and disk icon cache
│   ├── icon_cache_test.go # Icon cache unit tests
│   ├── scripts.go         # Raycast script command discovery & runner
│   └── scripts_test.go    # Script metadata & execution unit tests
├── frontend/              # Svelte 5 + Bun + Vite frontend application
│   ├── index.html         # Webview host HTML
│   ├── package.json       # Bun dependencies & scripts
│   ├── svelte.config.js   # Svelte configuration
│   ├── tsconfig.json      # TypeScript compiler settings
│   ├── vite.config.ts     # Vite 7 + Tailwind CSS v4 bundler config
│   └── src/
│       ├── App.svelte     # Main launcher interface & view coordinator
│       ├── main.ts        # Svelte 5 mount entrypoint
│       ├── style.css      # Tailwind CSS v4 entrypoint
│       └── wailsjs/       # Auto-generated Wails TypeScript bindings
├── README.md              # Project overview, quickstart & documentation index
├── docs/                  # Architecture, plan, and operational specifications
│   ├── architecture.md    # Subsystem design & Mermaid architecture
│   ├── design_guidelines.md # Glassmorphism & UI design system
│   ├── operations.md      # Bun & developer workflow guide
│   ├── plan.md            # Feature roadmap & priority matrix
│   └── rules.md           # Engineering rules & quality standards (this file)
└── build/                 # Wails build templates, icons, and output binaries
```

- **Backend Logic**: Put Go services and system routines inside `backend/` or `app.go`. Keep individual modules modular and self-contained.
- **Frontend UI**: Put UI components and client logic inside `frontend/src/`.
- **Docs**: Put architectural, product, and guideline updates in `docs/`.

### B. Consistent Styling & Patterns
- **Go Conventions**:
  - Follow idiomatic Go naming (`camelCase` unexported, `PascalCase` exported).
  - Short receiver names (e.g., `func (c *ClipboardManager) Add(...)`).
  - Return early upon errors to reduce indentation depth.
  - All Go files must pass `gofmt -s`.
- **Svelte 5 Conventions**:
  - Use modern Svelte 5 runes (`$state`, `$derived`, `$effect`).
  - Do not mix legacy Svelte 3/4 stores or syntax when runes apply.
  - Follow the glassmorphism design system in [design_guidelines.md](design_guidelines.md) (electric blue accents, translucent dark backdrops, smooth transitions).
- **Custom Glassmorphism CSS Conventions**:
  - Pure custom CSS with centralized design tokens in `frontend/src/style.css`.
  - Zero Tailwind CSS. Use clean, scoped, semantic CSS in component `<style>` blocks.
  - Consistent colors (`--surface-canvas`, `--surface-header`, `--surface-card`, `--accent-blue`, `--accent-sky`, `--border-glass`, `--border-subtle`).
  - Follow the glassmorphism design system in [design_guidelines.md](design_guidelines.md) (electric blue accents, translucent dark backdrops, smooth transitions).
- **Strict Zero-Emoji Policy (Applications & Docs)**:
  - Never use emojis anywhere in the application code, UI components, templates, or documentation files.
  - Use vector SVGs, Lucide icons, or semantic plain text for all visual indicators and badges.

### C. Minimal & Surgical Changes
- **No Scope Creep**: Only touch files directly related to the assigned task.
- **Preserve Unrelated Work**: Preserve existing comments, docstrings, and structures in untouched areas.
- **No Unsolicited Dependencies**: Do not install new Go packages or npm/bun libraries without explicit justification or user instruction.

### D. Verification Before Completion
Before completing any task, agents must run:
1. `gofmt -s -w .` and `go vet ./...`
2. `cd frontend && bun run check`
3. `make test` (or `go test ./backend/...`)

---

## 6. Pre-Commit / Pre-Push Checklist

| Check | Requirement | Command |
| :--- | :--- | :--- |
| **Simplicity** | Code is simple, readable, and free of over-engineering | Manual inspection |
| **No Dead Code** | Unused functions, imports, types, and comments removed | Manual inspection |
| **1:1 Test Files** | Every Go backend file has a corresponding `_test.go` file | Manual inspection |
| **Zero Emojis** | Zero emojis in application UI, source code, and docs | Manual / regex audit |
| **Go Format** | All Go files formatted with standard `gofmt` | `gofmt -s -w .` |
| **Go Vet** | Zero Go compiler or static analysis warnings | `go vet ./...` |
| **Go Tests** | Backend unit tests pass | `go test -v ./backend/...` |
| **Frontend Check** | Zero Svelte/TypeScript errors or warnings | `cd frontend && bun run check` |
| **No Lint Bypasses** | No `//nolint`, `@ts-ignore`, or silent swallows | `git diff` audit |
| **Clean Tree** | No unintended files or temporary artifacts staged | `git status` |

---

## 7. Mandatory 1:1 Test File Rule

Every functional Go file in `backend/` must have a dedicated corresponding test file:

- `backend/apps.go` -> `backend/apps_test.go`
- `backend/calc.go` -> `backend/calc_test.go`
- `backend/clipboard.go` -> `backend/clipboard_test.go`
- `backend/daemon.go` -> `backend/daemon_test.go`
- `backend/fuzzy.go` -> `backend/fuzzy_test.go`
- `backend/icon_cache.go` -> `backend/icon_cache_test.go`
- `backend/scripts.go` -> `backend/scripts_test.go`

### Rules for Test Files:
1. **Strict 1:1 Mapping**: Every new backend file introduced must be paired with its unit test file in the same commit.
2. **Table-Driven Tests**: Prefer table-driven testing patterns for parsers, evaluators, and search algorithms.
3. **No External Network Dependencies in Unit Tests**: Unit tests must execute locally and offline without requiring internet access or graphical window surfaces.
4. **Fast Execution**: Unit tests must run to completion within milliseconds (`make test`).

---

## 8. Frontend Modular Architecture & Pure Glassmorphism CSS Policy

> [!IMPORTANT]
> **No Monolithic Frontend Files & Pure Custom Glassmorphism CSS.**
> - Never pack all templates, state, styles, and utilities into a single giant file.
> - **Pure Custom CSS (Zero Tailwind)**: Use scoped custom CSS in component `<style>` blocks with consistent global design tokens.
> - **Rich Glassmorphism**: Maintain consistent frosted-glass visual hierarchy (`backdrop-filter: blur(...)`, electric blue highlights, subtle borders).

Frontend source code in `frontend/src/` must adhere to modular separation:

1. **Consistent Custom Glassmorphism CSS**:
   - Centralize design tokens and theme palettes (`--surface-canvas`, `--surface-card`, `--accent-blue`, etc.) in `style.css`.
   - Style components using clean, scoped CSS in each component's `<style>` block.
   - Do not use Tailwind CSS. Rely on standard, deterministic CSS flexbox and grid.

2. **Dedicated Svelte Components (`frontend/src/components/`)**:
   - Break UI features into focused, single-responsibility components.
   - Examples: `Header.svelte`, `ResultsList.svelte`, `DetailPanel.svelte`, `ActionPalette.svelte`, `Footer.svelte`, `MathBanner.svelte`.
   - Components receive typed props via Svelte 5 `$props()` and communicate callbacks cleanly.
   - Max file size recommendation: under 180 lines per component.

3. **Domain Utilities & Pure Functions (`frontend/src/utils/`)**:
   - Pure utility functions (e.g., viewport scroll calculations, date/time formatting, DOM helpers) must live in dedicated `.ts` files:
     - `utils/scroll.ts`: Viewport synchronization and scroll calculation logic.
     - `utils/format.ts`: Timestamp formatting, string truncation, and number formatting.

4. **Shared TypeScript Interfaces & Types (`frontend/src/types/`)**:
   - Shared data contracts, model interfaces, and union types must live in `types/` (e.g., `frontend/src/types/index.ts`).
   - Do not redefine identical types locally across multiple files.


