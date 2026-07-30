package utils

import (
	"errors"
	"fmt"
	"os"
	"uncomment-cli/internal/fs"
)

func IsFile(fs fs.FS, path string) (bool, error) {
	info, err := fs.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("file does not exist: %w", err)
		}
		return false, fmt.Errorf("failed to check file: %w", err)
	}
	return !info.IsDir(), nil
}
