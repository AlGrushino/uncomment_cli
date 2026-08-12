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
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %v", err)
	}
	absHome, err := filepath.Abs(home)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute home dir: %v", err)
	}

	if strings.HasPrefix(abs, absHome+string(filepath.Separator)) || abs == absHome {
		return abs, nil
	}

	tmpDir := os.TempDir()
	absTmp, err := filepath.Abs(tmpDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute temp dir: %v", err)
	}
	if strings.HasPrefix(abs, absTmp+string(filepath.Separator)) || abs == absTmp {
		return abs, nil
	}

	return "", fmt.Errorf("path %q is outside allowed directories (%q or %q)", path, home, tmpDir)
}
