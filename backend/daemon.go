package backend

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetSocketPath() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = "."
	}
	sockDir := filepath.Join(dir, ".local/share/beam")
	if mkErr := os.MkdirAll(sockDir, 0755); mkErr != nil {
		return filepath.Join(os.TempDir(), "beam.sock")
	}
	return filepath.Join(sockDir, "beam.sock")
}

// TrySignalExistingInstance checks if Beam is already running and sends toggle signal
func TrySignalExistingInstance(sockPath string) bool {
	conn, err := net.DialTimeout("unix", sockPath, 200*time.Millisecond)
	if err != nil {
		// Clean up stale socket file if it exists but no process is listening
		if rmErr := os.Remove(sockPath); rmErr != nil && !os.IsNotExist(rmErr) {
			return false
		}
		return false
	}
	defer conn.Close()

	if _, writeErr := conn.Write([]byte("toggle\n")); writeErr != nil {
		return false
	}
	return true
}

// StartDaemonServer listens for socket commands like toggle, show, hide
func StartDaemonServer(ctx context.Context, sockPath string, onCommand func(cmd string)) error {
	if rmErr := os.Remove(sockPath); rmErr != nil && !os.IsNotExist(rmErr) {
		return fmt.Errorf("failed to clear existing socket: %w", rmErr)
	}

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}

	go func() {
		<-ctx.Done()
		if closeErr := listener.Close(); closeErr != nil {
			return
		}
		if rmErr := os.Remove(sockPath); rmErr != nil && !os.IsNotExist(rmErr) {
			return
		}
	}()

	go func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				select {
				case <-ctx.Done():
					return
				default:
					return
				}
			}

			go handleConnection(conn, onCommand)
		}
	}()

	return nil
}

func handleConnection(conn net.Conn, onCommand func(cmd string)) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		cmd := strings.TrimSpace(scanner.Text())
		if onCommand != nil && cmd != "" {
			onCommand(cmd)
		}
	}
	if err := scanner.Err(); err != nil {
		return
	}
}
