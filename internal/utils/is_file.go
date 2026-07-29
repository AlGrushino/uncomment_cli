package utils

import (
	"errors"
	"fmt"
	"os"
)

func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("file does not exist: %w", err)
		}
		return false, fmt.Errorf("failed to check file: %w", err)
	}
	return !info.IsDir(), nil
}
