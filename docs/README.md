# 📚 Beam Documentation Index

Welcome to the engineering documentation for **Beam**, a fast, keyboard-centric Raycast-inspired launcher for Linux built with **Go**, **Wails v2**, **Bun**, **Svelte 5**, and **Tailwind CSS v4**.

---

## 📑 Documentation Sections

| Document | Description |
| :--- | :--- |
| **[Architecture](architecture.md)** | Subsystem architecture prioritized by P0/P1/P2/P3, Mermaid diagrams, Go engine, and IPC bridge. |
| **[Product & Technical Plan](plan.md)** | Priority Matrix (P0 to P3), acceptance criteria, and technical stack specifications. |
| **[Operations & Developer Guide](operations.md)** | Bun setup commands, Tailwind CSS v4 Vite integration, hot reload, and desktop hotkeys. |
| **[Milestones & Roadmap](milestones.md)** | Priority-driven development milestones from P0 MVP to P3 Theming. |

---

## ⚡ Quick Start (Bun-Powered)

```bash
cd /home/dilkhush/Desktop/Mine/beam

# 1. Install frontend dependencies with Bun
cd frontend && bun install && cd ..

# 2. Build the standalone binary
wails build -tags webkit2_41

# 3. Run Beam
./build/bin/beam
```
