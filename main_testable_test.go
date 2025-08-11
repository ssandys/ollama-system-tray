package main

import (
	"context"
	"errors"
	"testing"
)

// MockCommandRunner for testing
type MockCommandRunner struct {
	RunError   error
	StartError error
	Commands   []MockCommand
}

type MockCommand struct {
	Name string
	Args []string
}

func (m *MockCommandRunner) RunCommand(ctx context.Context, name string, args ...string) error {
	m.Commands = append(m.Commands, MockCommand{Name: name, Args: args})
	return m.RunError
}

func (m *MockCommandRunner) StartCommand(name string, args ...string) error {
	m.Commands = append(m.Commands, MockCommand{Name: name, Args: args})
	return m.StartError
}

func TestTestableIsOllamaRunning(t *testing.T) {
	tests := []struct {
		name     string
		runError error
		expected bool
	}{
		{
			name:     "ollama is running",
			runError: nil,
			expected: true,
		},
		{
			name:     "ollama is not running",
			runError: errors.New("process not found"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockCommandRunner{RunError: tt.runError}
			result := TestableIsOllamaRunning(mock)
			
			if result != tt.expected {
				t.Errorf("TestableIsOllamaRunning() = %v, want %v", result, tt.expected)
			}
			
			// Verify correct command was called
			if len(mock.Commands) != 1 {
				t.Errorf("Expected 1 command, got %d", len(mock.Commands))
			}
			
			cmd := mock.Commands[0]
			if cmd.Name != "pgrep" {
				t.Errorf("Expected command 'pgrep', got '%s'", cmd.Name)
			}
			
			expectedArgs := []string{"-f", "ollama serve"}
			if len(cmd.Args) != len(expectedArgs) {
				t.Errorf("Expected args %v, got %v", expectedArgs, cmd.Args)
			}
			
			for i, arg := range expectedArgs {
				if cmd.Args[i] != arg {
					t.Errorf("Expected arg[%d] '%s', got '%s'", i, arg, cmd.Args[i])
				}
			}
		})
	}
}

func TestTestableStartOllama(t *testing.T) {
	tests := []struct {
		name       string
		startError error
		expectErr  bool
	}{
		{
			name:       "successful start",
			startError: nil,
			expectErr:  false,
		},
		{
			name:       "failed start",
			startError: errors.New("permission denied"),
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockCommandRunner{StartError: tt.startError}
			err := TestableStartOllama(mock)
			
			if tt.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			
			// Verify correct command was called
			if len(mock.Commands) != 1 {
				t.Errorf("Expected 1 command, got %d", len(mock.Commands))
			}
			
			cmd := mock.Commands[0]
			if cmd.Name != "pkexec" {
				t.Errorf("Expected command 'pkexec', got '%s'", cmd.Name)
			}
			
			expectedArgs := []string{"systemctl", "start", "ollama"}
			if len(cmd.Args) != len(expectedArgs) {
				t.Errorf("Expected args %v, got %v", expectedArgs, cmd.Args)
			}
		})
	}
}

func TestTestableStopOllama(t *testing.T) {
	tests := []struct {
		name      string
		runError  error
		expectErr bool
	}{
		{
			name:      "successful stop",
			runError:  nil,
			expectErr: false,
		},
		{
			name:      "failed stop",
			runError:  errors.New("service not found"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockCommandRunner{RunError: tt.runError}
			err := TestableStopOllama(mock)
			
			if tt.expectErr && err == nil {
				t.Error("Expected error but got nil")
			}
			
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			
			// Verify correct command was called
			if len(mock.Commands) != 1 {
				t.Errorf("Expected 1 command, got %d", len(mock.Commands))
			}
			
			cmd := mock.Commands[0]
			if cmd.Name != "pkexec" {
				t.Errorf("Expected command 'pkexec', got '%s'", cmd.Name)
			}
			
			expectedArgs := []string{"systemctl", "stop", "ollama"}
			if len(cmd.Args) != len(expectedArgs) {
				t.Errorf("Expected args %v, got %v", expectedArgs, cmd.Args)
			}
		})
	}
}

// Benchmark tests
func BenchmarkTestableIsOllamaRunning(b *testing.B) {
	mock := &MockCommandRunner{RunError: nil}
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		TestableIsOllamaRunning(mock)
	}
}