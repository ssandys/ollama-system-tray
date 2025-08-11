// Package main implements a system tray application for managing the Ollama service on Linux systems.
// Specifically designed for Hyprland and waybar compatibility with minimal dependencies.
package main

import (
	"context"
	_ "embed" // Required for embedding the icon file
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/systray"
)

// iconData embeds the Ollama icon at compile time for use in the system tray
//go:embed assets/icons/ollama-icon.png
var iconData []byte

// main initializes and runs the system tray application
func main() {
	// systray.Run blocks until the application is terminated
	systray.Run(OnReady, OnExit)
}

// OnReady is called when the system tray is ready to be configured
func OnReady() {
	// Initialize system tray appearance
	systray.SetIcon(iconData)
	systray.SetTitle("Ollama")
	systray.SetTooltip("Ollama Server Manager")

	// Create menu items for the system tray
	mStatus := systray.AddMenuItem("Status: Checking...", "Check Ollama Server Status")
	systray.AddSeparator()
	mStart := systray.AddMenuItem("Start Ollama", "Start Ollama Server")
	mStop := systray.AddMenuItem("Stop Ollama", "Stop Ollama Server")
	mRestart := systray.AddMenuItem("Restart Ollama", "Restart Ollama Server")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Start goroutine to handle menu events and periodic status updates
	go func() {
		// Create ticker for periodic status updates every 5 seconds
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// Initial status update
		UpdateStatus(mStatus)

		// Event loop for handling menu clicks and periodic updates
		for {
			select {
			case <-ticker.C:
				// Periodic status update
				UpdateStatus(mStatus)
			case <-mStart.ClickedCh:
				// Start Ollama service and update status
				StartOllama()
				UpdateStatus(mStatus)
			case <-mStop.ClickedCh:
				// Stop Ollama service and update status
				StopOllama()
				UpdateStatus(mStatus)
			case <-mRestart.ClickedCh:
				// Restart Ollama service: stop, wait, then start
				StopOllama()
				time.Sleep(2 * time.Second) // Allow time for service to fully stop
				StartOllama()
				UpdateStatus(mStatus)
			case <-mQuit.ClickedCh:
				// Quit the application
				systray.Quit()
				return
			}
		}
	}()

	// Set up signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c // Wait for signal
		systray.Quit()
	}()
}

// OnExit is called when the system tray application is terminating
func OnExit() {
	log.Println("Ollama system tray exiting...")
}

// isOllamaRunning checks if the Ollama service is currently running
// Returns true if the service is running, false otherwise
func isOllamaRunning() bool {
	// Create context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use pgrep to search for the "ollama serve" process
	cmd := exec.CommandContext(ctx, "pgrep", "-f", "ollama serve")
	err := cmd.Run()
	
	// pgrep returns 0 (no error) if process is found, non-zero if not found
	return err == nil
}

// UpdateStatus updates the system tray status menu item based on service state
func UpdateStatus(mStatus *systray.MenuItem) {
	if isOllamaRunning() {
		// Service is running - show green checkmark
		mStatus.SetTitle("Status: Running ✓")
		systray.SetIcon(iconData)
	} else {
		// Service is stopped - show red X
		mStatus.SetTitle("Status: Stopped ✗")
		systray.SetIcon(iconData)
	}
}

// StartOllama starts the Ollama service using systemctl with privilege escalation
func StartOllama() {
	// Use pkexec to prompt for authentication and start the systemd service
	cmd := exec.Command("pkexec", "systemctl", "start", "ollama")
	err := cmd.Start() // Use Start() to run asynchronously
	
	if err != nil {
		log.Printf("Failed to start Ollama: %v", err)
	} else {
		log.Println("Ollama service started")
	}
}

// StopOllama stops the Ollama service using systemctl with privilege escalation
func StopOllama() {
	// Use pkexec to prompt for authentication and stop the systemd service
	cmd := exec.Command("pkexec", "systemctl", "stop", "ollama")
	err := cmd.Run() // Use Run() to wait for completion
	
	if err != nil {
		log.Printf("Failed to stop Ollama: %v", err)
	} else {
		log.Println("Ollama service stopped")
	}
}
