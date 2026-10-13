package backend

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type AppInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	Exec        string `json:"exec"`
	Icon        string `json:"icon"`
	IconPath    string `json:"iconPath"`
	Terminal    bool   `json:"terminal"`
	DesktopPath string `json:"desktopPath"`
}

var (
	appsCache []AppInfo
	appsMutex sync.RWMutex
	percentRe = regexp.MustCompile(`%[a-zA-Z]`)
)

func GetDesktopDirectories() []string {
	dirs := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		"/var/lib/flatpak/exports/share/applications",
	}

	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs,
			filepath.Join(home, ".local/share/applications"),
			filepath.Join(home, ".local/share/flatpak/exports/share/applications"),
		)
	}

	return dirs
}

func ScanApplications() []AppInfo {
	appsMutex.Lock()
	defer appsMutex.Unlock()

	var apps []AppInfo
	seen := make(map[string]bool)

	dirs := GetDesktopDirectories()
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".desktop") {
				continue
			}

			fullPath := filepath.Join(dir, entry.Name())
			app, err := parseDesktopFile(fullPath)
			if err != nil || app.Name == "" || app.Exec == "" {
				continue
			}

			// Deduplicate by Name or Exec command
			key := strings.ToLower(app.Name)
			if seen[key] {
				continue
			}
			seen[key] = true

			app.IconPath = findIconPath(app.Icon)
			apps = append(apps, app)
		}
	}

	appsCache = apps
	return apps
}

func GetCachedApplications() []AppInfo {
	appsMutex.RLock()
	defer appsMutex.RUnlock()
	if len(appsCache) == 0 {
		appsMutex.RUnlock()
		return ScanApplications()
	}
	return appsCache
}

func parseDesktopFile(path string) (AppInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return AppInfo{}, err
	}
	defer file.Close()

	var app AppInfo
	app.DesktopPath = path
	app.ID = filepath.Base(path)

	scanner := bufio.NewScanner(file)
	inDesktopEntry := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") {
			inDesktopEntry = line == "[Desktop Entry]"
			continue
		}

		if !inDesktopEntry || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "Type":
			if val != "Application" {
				return AppInfo{}, fmt.Errorf("not an application")
			}
		case "NoDisplay", "Hidden":
			if strings.EqualFold(val, "true") {
				return AppInfo{}, fmt.Errorf("hidden application")
			}
		case "Name":
			if app.Name == "" {
				app.Name = val
			}
		case "Comment":
			if app.Comment == "" {
				app.Comment = val
			}
		case "Exec":
			// Clean Exec parameters like %u, %F, %U
			clean := percentRe.ReplaceAllString(val, "")
			app.Exec = strings.TrimSpace(clean)
		case "Icon":
			if app.Icon == "" {
				app.Icon = val
			}
		case "Terminal":
			app.Terminal = strings.EqualFold(val, "true")
		}
	}

	return app, nil
}

func findIconPath(iconName string) string {
	if iconName == "" {
		return ""
	}

	if filepath.IsAbs(iconName) {
		if _, err := os.Stat(iconName); err == nil {
			return iconName
		}
	}

	// Search common icon locations
	iconSearchPaths := []string{
		"/usr/share/pixmaps",
		"/usr/share/icons/hicolor/scalable/apps",
		"/usr/share/icons/hicolor/48x48/apps",
		"/usr/share/icons/hicolor/64x64/apps",
		"/usr/share/icons/hicolor/128x128/apps",
		"/usr/share/icons/hicolor/256x256/apps",
	}

	extensions := []string{".svg", ".png", ".xpm"}

	for _, basePath := range iconSearchPaths {
		for _, ext := range extensions {
			full := filepath.Join(basePath, iconName+ext)
			if _, err := os.Stat(full); err == nil {
				return full
			}
		}
		// Direct name match
		direct := filepath.Join(basePath, iconName)
		if _, err := os.Stat(direct); err == nil {
			return direct
		}
	}

	return ""
}

func Launch(execCmd string, terminal bool) error {
	if execCmd == "" {
		return fmt.Errorf("empty command")
	}

	parts := strings.Fields(execCmd)
	if len(parts) == 0 {
		return fmt.Errorf("invalid command")
	}

	cmdName := parts[0]
	args := parts[1:]

	var cmd *exec.Cmd
	if terminal {
		// Try launching in a default terminal
		terms := []string{"x-terminal-emulator", "gnome-terminal", "alacritty", "kitty", "foot", "xterm"}
		var termFound string
		for _, t := range terms {
			if path, err := exec.LookPath(t); err == nil {
				termFound = path
				break
			}
		}
		if termFound != "" {
			cmd = exec.Command(termFound, append([]string{"-e", cmdName}, args...)...)
		} else {
			cmd = exec.Command(cmdName, args...)
		}
	} else {
		cmd = exec.Command(cmdName, args...)
	}

	return cmd.Start()
}
