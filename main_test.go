package main

import (
	"testing"
	"time"
)

func TestCmd(t *testing.T) {
	// Test command execution
	cmd := "echo"
	args := []string{"Hello, World!"}
	timeout := 5 * time.Second

	opt := Opt{
		Command: cmd,
		Args:    args,
		Timeout: timeout,
	}

	status, duration := opt.cmd()
	if status != 0 {
		t.Fatalf("Expected exit code 0, got %d", status)
	}
	if duration <= 0 {
		t.Fatalf("Expected duration greater than 0, got %v", duration)
	}
}

// Test with unknown command to ensure it returns a non-zero exit code
func TestCmdUnknownCommand(t *testing.T) {
	cmd := "unknown_command"
	args := []string{}
	timeout := 5 * time.Second

	opt := Opt{
		Command: cmd,
		Args:    args,
		Timeout: timeout,
	}

	status, _ := opt.cmd()
	if status != UnknownCommandStatus {
		t.Fatalf("Expected non-zero exit code for unknown command, got %d", status)
	}
}
