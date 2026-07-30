package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AllowedPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to return abs path: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to return home dir: %v", err)
	}
	absHome, err := filepath.Abs(home)
	if err != nil {
		return "", fmt.Errorf("failed to return abs home dir: %v", err)
	}

	if !strings.HasPrefix(abs, absHome+string(filepath.Separator)) && abs != absHome {
		return "", fmt.Errorf("path %q is outside allowed directory %q", path, home)
	}

	return abs, nil
}
