#!/usr/bin/env bash
# ==============================================================================
# Beam Launcher - Automated Installation Script
# https://github.com/codesdilkhush/beam
# ==============================================================================

set -euo pipefail

# Color palette
BOLD="\033[1m"
DIM="\033[2m"
CYAN="\033[36m"
GREEN="\033[32m"
YELLOW="\033[33m"
RED="\033[31m"
RESET="\033[0m"

log_info() {
    printf "${CYAN}==>${RESET} ${BOLD}%s${RESET}\n" "$1"
}

log_success() {
    printf "${GREEN}==>${RESET} ${BOLD}%s${RESET}\n" "$1"
}

log_warn() {
    printf "${YELLOW}==> WARNING:${RESET} %s\n" "$1"
}

log_error() {
    printf "${RED}==> ERROR:${RESET} %s\n" "$1" >&2
}

REPO_OWNER="codesdilkhush"
REPO_NAME="beam"
GITHUB_REPO="${REPO_OWNER}/${REPO_NAME}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
DESKTOP_DIR="$HOME/.local/share/applications"
ICON_DIR="$HOME/.local/share/icons/hicolor/512x512/apps"
SYSTEMD_USER_DIR="$HOME/.config/systemd/user"
DATA_DIR="$HOME/.local/share/beam"

print_banner() {
    printf "\n"
    printf "  ${BOLD}${CYAN} ____  _____    _    __  __${RESET}\n"
    printf "  ${BOLD}${CYAN}| __ )| ____|  / \\\\  |  \\/  |${RESET}\n"
    printf "  ${BOLD}${CYAN}|  _ \\|  _|   / _ \\ | |\\/| |${RESET}\n"
    printf "  ${BOLD}${CYAN}| |_) | |___ / ___ \\| |  | |${RESET}\n"
    printf "  ${BOLD}${CYAN}|____/|_____/_/   \\_\\_|  |_|${RESET}\n"
    printf "  ${DIM}Ultra-fast Spotlight & Launcher for Linux${RESET}\n\n"
}

check_system() {
    OS="$(uname -s)"
    ARCH="$(uname -m)"

    if [ "$OS" != "Linux" ]; then
        log_error "Beam is designed specifically for Linux desktop environments. Detected OS: $OS"
        exit 1
    fi

    case "$ARCH" in
        x86_64|amd64)
            TARGET_ARCH="amd64"
            ;;
        aarch64|arm64)
            TARGET_ARCH="arm64"
            ;;
        *)
            log_error "Unsupported CPU architecture: $ARCH. Beam supports x86_64 and aarch64."
            exit 1
            ;;
    esac

    log_info "Detected platform: Linux ($TARGET_ARCH)"
}

check_prerequisites() {
    log_info "Verifying desktop utilities..."

    # Check clipboard utility for Wayland or X11
    if [ "${XDG_SESSION_TYPE:-}" = "wayland" ]; then
        if ! command -v wl-paste >/dev/null 2>&1; then
            log_warn "wl-clipboard is recommended on Wayland for clipboard history support."
            log_warn "Install via: sudo apt install wl-clipboard (or distro equivalent)"
        fi
    else
        if ! command -v xclip >/dev/null 2>&1; then
            log_warn "xclip is recommended on X11 for clipboard history support."
            log_warn "Install via: sudo apt install xclip (or distro equivalent)"
        fi
    fi
}

install_from_source() {
    log_info "Building Beam from local source repository..."

    if ! command -v go >/dev/null 2>&1; then
        log_error "Go 1.22+ is required to build from source. Please install Go from https://golang.org/dl/"
        exit 1
    fi

    if ! command -v bun >/dev/null 2>&1; then
        log_error "Bun is required to build frontend assets. Install via: curl -fsSL https://bun.sh/install | bash"
        exit 1
    fi

    if ! command -v wails >/dev/null 2>&1; then
        log_info "Installing Wails v2 CLI via go install..."
        go install github.com/wailsapp/wails/v2/cmd/wails@latest
        export PATH="$PATH:$HOME/go/bin"
    fi

    # Install frontend dependencies and build
    (cd frontend && bun install)
    log_info "Compiling optimized production binary (tags: webkit2_41)..."
    wails build -tags webkit2_41 -ldflags="-s -w" -clean

    mkdir -p "$INSTALL_DIR"
    cp "build/bin/beam" "$INSTALL_DIR/beam"
    chmod +x "$INSTALL_DIR/beam"
}

install_from_release() {
    log_info "Checking latest release on GitHub (${GITHUB_REPO})..."

    TMP_DIR="$(mktemp -d)"
    trap 'rm -rf "$TMP_DIR"' EXIT

    LATEST_RELEASE_URL="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
    DOWNLOAD_URL=""

    # Attempt to fetch release data from GitHub API
    if command -v curl >/dev/null 2>&1; then
        RELEASE_JSON="$(curl -fsSL "$LATEST_RELEASE_URL" 2>/dev/null || true)"
    elif command -v wget >/dev/null 2>&1; then
        RELEASE_JSON="$(wget -qO- "$LATEST_RELEASE_URL" 2>/dev/null || true)"
    else
        log_error "Neither curl nor wget was found. Please install curl or wget."
        exit 1
    fi

    if [ -n "$RELEASE_JSON" ]; then
        DOWNLOAD_URL="$(printf '%s' "$RELEASE_JSON" | grep -o 'https://[^"]*beam-linux-amd64.tar.gz' | head -n 1 || true)"
    fi

    # If GitHub release asset is available, download and extract
    if [ -n "$DOWNLOAD_URL" ]; then
        log_info "Downloading Beam release archive from $DOWNLOAD_URL..."
        curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/beam.tar.gz"
        tar -xzf "$TMP_DIR/beam.tar.gz" -C "$TMP_DIR"

        mkdir -p "$INSTALL_DIR"
        if [ -f "$TMP_DIR/beam" ]; then
            cp "$TMP_DIR/beam" "$INSTALL_DIR/beam"
        elif [ -f "$TMP_DIR/build/bin/beam" ]; then
            cp "$TMP_DIR/build/bin/beam" "$INSTALL_DIR/beam"
        else
            log_error "Extracted archive did not contain beam executable."
            exit 1
        fi
        chmod +x "$INSTALL_DIR/beam"
    else
        # If no release tarball exists yet on GitHub, check if we're inside the repo or clone it
        log_info "No compiled release binaries found on GitHub yet. Falling back to local build..."
        if [ -f "main.go" ] && [ -d "frontend" ]; then
            install_from_source
        else
            log_info "Cloning Beam repository to compile from source..."
            git clone --depth 1 "https://github.com/${GITHUB_REPO}.git" "$TMP_DIR/beam-repo"
            (cd "$TMP_DIR/beam-repo" && install_from_source)
        fi
    fi
}

install_desktop_integration() {
    log_info "Configuring desktop integration and icons..."

    mkdir -p "$DESKTOP_DIR"
    mkdir -p "$ICON_DIR"
    mkdir -p "$DATA_DIR"

    # Install desktop entry
    cat <<EOF > "$DESKTOP_DIR/beam.desktop"
[Desktop Entry]
Name=Beam
Comment=Ultra-fast keyboard-first spotlight launcher and productivity tool
Exec=${INSTALL_DIR}/beam
Icon=beam
Terminal=false
Type=Application
Categories=Utility;System;
Keywords=spotlight;launcher;raycast;search;todo;emoji;
StartupNotify=false
EOF
    chmod 644 "$DESKTOP_DIR/beam.desktop"

    # Copy icon if present
    if [ -f "build/appicon.png" ]; then
        cp "build/appicon.png" "$ICON_DIR/beam.png"
    elif [ -f "$HOME/.local/share/icons/hicolor/512x512/apps/beam.png" ]; then
        :
    fi

    # Configure systemd user daemon service
    mkdir -p "$SYSTEMD_USER_DIR"
    cat <<EOF > "$SYSTEMD_USER_DIR/beam.service"
[Unit]
Description=Beam Application Launcher Daemon
After=graphical-session.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/beam --daemon
Restart=on-failure
RestartSec=3

[Install]
WantedBy=default.target
EOF

    if command -v systemctl >/dev/null 2>&1; then
        log_info "Enabling and starting Beam background daemon service..."
        systemctl --user daemon-reload 2>/dev/null || true
        systemctl --user enable --now beam.service 2>/dev/null || true
        log_success "Background daemon service auto-enabled (beam.service active)."
    fi

    if command -v update-desktop-database >/dev/null 2>&1; then
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    fi
}

check_path() {
    case ":${PATH}:" in
        *:"${INSTALL_DIR}":*)
            ;;
        *)
            log_warn "${INSTALL_DIR} is not currently in your system PATH."
            log_warn "Add the following line to your ~/.bashrc or ~/.zshrc:"
            printf "    ${BOLD}export PATH=\"\$PATH:%s\"${RESET}\n\n" "$INSTALL_DIR"
            ;;
    esac
}

print_summary() {
    printf "\n"
    log_success "Beam has been successfully installed to ${BOLD}${INSTALL_DIR}/beam${RESET}!"
    printf "\n"
    printf "  ${BOLD}Status:${RESET}\n"
    printf "  - Background daemon: Auto-enabled and running for instant popup.\n\n"
    printf "  ${BOLD}Assign a Global Shortcut to invoke beam --toggle:${RESET}\n"
    printf "  - GNOME: Settings -> Keyboard -> Custom Shortcuts -> Add beam --toggle (Super+Space)\n"
    printf "  - KDE Plasma: System Settings -> Shortcuts -> Custom Shortcuts -> beam --toggle\n"
    printf "  - Hyprland: Add 'bind = SUPER, SPACE, exec, beam --toggle' to hyprland.conf\n"
    printf "  - Sway / i3: Add 'bindsym \$mod+space exec beam --toggle' to config\n\n"
    printf "  ${BOLD}Interactive Guide:${RESET} Summon Beam and type ${CYAN}/guide${RESET} or ${CYAN}/?${RESET}\n"
    printf "\n"
}

main() {
    print_banner
    check_system
    check_prerequisites

    # If run inside the repo and with --build or if release not yet published
    if [ "${1:-}" = "--build" ]; then
        install_from_source
    else
        install_from_release
    fi

    install_desktop_integration
    check_path
    print_summary
}

main "$@"
