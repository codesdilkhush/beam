package backend

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseScriptHeader(t *testing.T) {
	tempDir := t.TempDir()
	scriptFile := filepath.Join(tempDir, "sample.sh")

	content := `#!/bin/bash
# @raycast.schemaVersion 1
# @raycast.title Check Public IP
# @raycast.mode compact
# @raycast.description Fetches external IP address
# @raycast.author Beam Team

echo "127.0.0.1"
`
	if err := os.WriteFile(scriptFile, []byte(content), 0755); err != nil {
		t.Fatalf("failed to write test script: %v", err)
	}

	cmd, err := ParseScriptHeader(scriptFile)
	if err != nil {
		t.Fatalf("ParseScriptHeader failed: %v", err)
	}

	if cmd.Title != "Check Public IP" {
		t.Errorf("expected Title 'Check Public IP', got %q", cmd.Title)
	}
	if cmd.Mode != "compact" {
		t.Errorf("expected Mode 'compact', got %q", cmd.Mode)
	}
	if cmd.Description != "Fetches external IP address" {
		t.Errorf("expected Description 'Fetches external IP address', got %q", cmd.Description)
	}
	if cmd.Author != "Beam Team" {
		t.Errorf("expected Author 'Beam Team', got %q", cmd.Author)
	}
}

func TestRunScript(t *testing.T) {
	tempDir := t.TempDir()
	scriptFile := filepath.Join(tempDir, "echo_test.sh")

	content := `#!/bin/sh
echo "Beam Test Success"
`
	if err := os.WriteFile(scriptFile, []byte(content), 0755); err != nil {
		t.Fatalf("failed to write test script: %v", err)
	}

	output, err := RunScript(scriptFile)
	if err != nil {
		t.Fatalf("RunScript failed: %v", err)
	}

	if output != "Beam Test Success" {
		t.Errorf("expected output 'Beam Test Success', got %q", output)
	}
}
