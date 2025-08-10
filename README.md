# Ollama System Tray

A simple system tray application for managing the Ollama service on Linux, specifically designed for Hyprland and waybar compatibility.

## Features

- Real-time status monitoring of Ollama service
- Start/Stop/Restart Ollama service from the system tray
- Visual status indicators (different icons for running/stopped states)
- Auto-refreshes status every 5 seconds
- Compatible with waybar tray-expander module

## Building

```bash
go build -o ollama-tray main.go
```

## Usage

Simply run the executable:

```bash
./ollama-system-tray
```

The application will appear in your system tray. Right-click to see the menu options:
- Status display (updates automatically)
- Start Ollama
- Stop Ollama  
- Restart Ollama
- Quit

## Waybar Integration

This tray application is compatible with waybar's tray-expander module. Add the following to your waybar configuration:

```json
"tray": {
    "icon-theme": "Adwaita",
    "show-passive-items": true,
    "spacing": 10
}
```

## Installation

1. Build the application: `go build -o ollama-system-tray main.go`
2. Copy the binary to a permanent location: `cp ollama-system-tray ~/.local/bin/`
3. Make it executable: `chmod +x ~/.local/bin/ollama-system-tray`
4. Optionally install the desktop entry: `cp ollama-system-tray.desktop ~/.local/share/applications/`

## Requirements

- Go 1.19 or later
- Linux with system tray support
- Ollama installed and in PATH

## Dependencies

- fyne.io/systray - Cross-platform system tray library