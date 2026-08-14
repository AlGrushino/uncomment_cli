package utils

import (
	"errors"
	"fmt"
	"os"
	"uncomment-cli/internal/fs"
)

// IsFile reports whether the given path refers to a regular file (not a
// directory). It returns an error if the path does not exist.
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
