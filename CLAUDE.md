# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a Go-based system tray application for managing the Ollama service on Linux systems, specifically designed for Hyprland and waybar compatibility. It's a single-file application with minimal dependencies.

## Build and Development Commands

### Building
```bash
go build -o ollama-system-tray main.go
```

### Running
```bash
./ollama-system-tray
```

### Testing
```bash
go test ./...
```

### Module Management
```bash
go mod tidy        # Clean up dependencies
go mod download    # Download dependencies
```

## Architecture

### Core Components
- **main.go**: Single-file application containing all logic
- **System Tray Integration**: Uses `fyne.io/systray` library for cross-platform tray functionality
- **Service Management**: Uses shell commands (`pgrep`, `pkexec`, and `systemctl`) to manage Ollama systemd service

### Key Functions
- `OnReady()`: Initializes system tray UI and event handlers (main.go:29)
- `isOllamaRunning()`: Checks if Ollama service is running using `pgrep` (main.go:97)
- `UpdateStatus()`: Updates tray menu status display (main.go:111)
- `StartOllama()`: Starts Ollama service with `pkexec systemctl start ollama` (main.go:124)
- `StopOllama()`: Stops Ollama service with `pkexec systemctl stop ollama` (main.go:137)

### Application Flow
1. System tray initializes with menu items (Status, Start, Stop, Restart, Quit)
2. Status updates automatically every 5 seconds via ticker
3. Menu actions trigger corresponding service management functions
4. Signal handling for graceful shutdown on SIGTERM/SIGINT

## Dependencies
- Go 1.24.6 (as specified in go.mod)
- `fyne.io/systray v1.11.0`: Cross-platform system tray library
- System requirements: Linux with system tray support, Ollama systemd service installed

## Installation Process
The application includes a desktop entry file (`ollama-system-tray.desktop`) for system integration. Standard installation involves building the binary and copying to `~/.local/bin/`.