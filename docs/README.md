# 📚 Beam Documentation Index

Welcome to the engineering documentation for **Beam**, a fast, keyboard-centric Raycast-inspired launcher for Linux built with Go, Wails v2, and Svelte 5.

---

## 📑 Documentation Sections

| Document | Description |
| :--- | :--- |
| **[Architecture](architecture.md)** | System overview, Mermaid component diagrams, Go backend services, IPC bridge, and memory model. |
| **[Product & Technical Plan](plan.md)** | Product vision, core feature pillars, UI/UX specifications, and non-functional requirements. |
| **[Operations & Developer Guide](operations.md)** | System prerequisites, local development commands, production builds, and desktop hotkey setup. |
| **[Milestones & Roadmap](milestones.md)** | Phased implementation milestones (M1 to M6) detailing completed and upcoming features. |

---

## ⚡ Quick Start
To build and test the current version of Beam:

```bash
cd /home/dilkhush/Desktop/Mine/beam
wails build -tags webkit2_41
./build/bin/beam
```
