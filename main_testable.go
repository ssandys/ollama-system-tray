package main

import (
	"context"
	"os/exec"
	"time"
)

// CommandRunner interface for mocking exec commands
type CommandRunner interface {
	RunCommand(ctx context.Context, name string, args ...string) error
	StartCommand(name string, args ...string) error
}

// RealCommandRunner implements CommandRunner using real exec commands
type RealCommandRunner struct{}

func (r *RealCommandRunner) RunCommand(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Run()
}

func (r *RealCommandRunner) StartCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Start()
}


// TestableIsOllamaRunning allows for dependency injection
func TestableIsOllamaRunning(runner CommandRunner) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	err := runner.RunCommand(ctx, "pgrep", "-f", "ollama serve")
	return err == nil
}

// TestableStartOllama allows for dependency injection
func TestableStartOllama(runner CommandRunner) error {
	return runner.StartCommand("pkexec", "systemctl", "start", "ollama")
}

// TestableStopOllama allows for dependency injection
func TestableStopOllama(runner CommandRunner) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return runner.RunCommand(ctx, "pkexec", "systemctl", "stop", "ollama")
}