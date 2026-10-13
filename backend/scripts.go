package backend

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ScriptCommand struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	Mode        string `json:"mode"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

func GetScriptsDirectory() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = "."
	}
	scriptsDir := filepath.Join(dir, ".config/beam/scripts")
	_ = os.MkdirAll(scriptsDir, 0755)
	return scriptsDir
}

func DiscoverScripts() []ScriptCommand {
	dir := GetScriptsDirectory()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []ScriptCommand{}
	}

	var scripts []ScriptCommand
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())
		cmd, parseErr := ParseScriptHeader(fullPath)
		if parseErr == nil && cmd.Title != "" {
			scripts = append(scripts, cmd)
		}
	}
	return scripts
}

func ParseScriptHeader(path string) (ScriptCommand, error) {
	file, err := os.Open(path)
	if err != nil {
		return ScriptCommand{}, err
	}
	defer file.Close()

	cmd := ScriptCommand{
		ID:   filepath.Base(path),
		Path: path,
		Mode: "fullOutput",
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() && lineCount < 40 {
		line := strings.TrimSpace(scanner.Text())
		lineCount++

		if strings.HasPrefix(line, "#") && strings.Contains(line, "@raycast.") {
			idx := strings.Index(line, "@raycast.")
			tagContent := strings.TrimSpace(line[idx+9:])
			parts := strings.SplitN(tagContent, " ", 2)
			if len(parts) < 2 {
				continue
			}

			key := strings.ToLower(strings.TrimSpace(parts[0]))
			val := strings.TrimSpace(parts[1])

			switch key {
			case "title":
				cmd.Title = val
			case "mode":
				cmd.Mode = val
			case "description":
				cmd.Description = val
			case "author":
				cmd.Author = val
			}
		}
	}

	if cmd.Title == "" {
		cmd.Title = filepath.Base(path)
	}

	return cmd, nil
}

func RunScript(path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, path)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("script failed: %s (%w)", strings.TrimSpace(stderr.String()), err)
		}
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}
