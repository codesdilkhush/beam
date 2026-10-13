# ==============================================================================
# Beam Launcher - Makefile
# ==============================================================================

SHELL := /bin/bash
BINARY_NAME := beam
BUILD_DIR := build/bin
OUTPUT_BIN := $(BUILD_DIR)/$(BINARY_NAME)
INSTALL_DIR := $(HOME)/.local/bin
TAGS := webkit2_41

.PHONY: all help dev dab build build-prod prod run install-deps deps clean install test

# Default target
all: help

## help: Display this help message with available commands
help:
	@echo "Beam Launcher - Development & Build Tasks"
	@echo "--------------------------------------------------------"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Development:"
	@echo "  make dev (or dab)   Start development server with live reload (Wails + Vite)"
	@echo "  make run            Run the compiled binary directly"
	@echo ""
	@echo "Building:"
	@echo "  make build          Build standard development binary"
	@echo "  make build-prod     Build optimized production binary (stripped -s -w)"
	@echo "  make prod           Alias for build-prod"
	@echo ""
	@echo "Dependencies & Maintenance:"
	@echo "  make install-deps   Install frontend dependencies using Bun"
	@echo "  make clean          Clean build artifacts and temporary files"
	@echo "  make install        Install compiled binary to ~/.local/bin/beam"
	@echo "  make test           Run Go backend unit tests"
	@echo "--------------------------------------------------------"

## dev: Run application in development mode with live reload
dev:
	wails dev -tags $(TAGS)

## dab: Alias for dev
dab: dev

## build: Compile standard standalone binary
build:
	wails build -tags $(TAGS)

## build-prod: Compile optimized production binary with stripped symbols
build-prod:
	wails build -tags $(TAGS) -ldflags="-s -w"

## prod: Alias for build-prod
prod: build-prod

## run: Run the compiled binary
run:
	@if [ -f "$(OUTPUT_BIN)" ]; then \
		$(OUTPUT_BIN); \
	else \
		echo "Binary not found at $(OUTPUT_BIN). Run 'make build' first."; \
		exit 1; \
	fi

## install-deps: Install frontend dependencies via Bun
install-deps:
	cd frontend && bun install

## deps: Alias for install-deps
deps: install-deps

## clean: Remove build outputs and cached bundles
clean:
	rm -rf $(BUILD_DIR)
	rm -rf frontend/dist

## install: Install binary to user's local bin (~/.local/bin)
install: build-prod
	mkdir -p $(INSTALL_DIR)
	cp $(OUTPUT_BIN) $(INSTALL_DIR)/$(BINARY_NAME)
	chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Successfully installed $(BINARY_NAME) to $(INSTALL_DIR)/$(BINARY_NAME)"

## test: Run backend Go tests
test:
	go test -v ./backend/...
