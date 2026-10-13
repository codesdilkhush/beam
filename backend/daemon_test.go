package backend

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestDaemonServer_Toggle(t *testing.T) {
	tempDir := t.TempDir()
	sockPath := filepath.Join(tempDir, "test_daemon.sock")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var receivedCmd string
	wg.Add(1)

	err := StartDaemonServer(ctx, sockPath, func(cmd string) {
		receivedCmd = cmd
		wg.Done()
	})
	if err != nil {
		t.Fatalf("failed to start daemon server: %v", err)
	}

	// Try signaling instance
	signaled := TrySignalExistingInstance(sockPath)
	if !signaled {
		t.Fatalf("expected TrySignalExistingInstance to succeed")
	}

	wg.Wait()
	if receivedCmd != "toggle" {
		t.Errorf("expected received command 'toggle', got %q", receivedCmd)
	}

	// Cancel server and verify cleanup
	cancel()
	time.Sleep(50 * time.Millisecond)

	// After server stops, signaling should return false
	if TrySignalExistingInstance(sockPath) {
		t.Errorf("expected TrySignalExistingInstance to fail after server closed")
	}
}
