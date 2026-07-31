package utils

import (
	"errors"
	"fmt"
)

func GetFilePath(path string) (string, error) {
	if len(path) == 0 {
		return "", errors.New("failed to get path: length of path equals 0")
	}
	if path[0] != '/' {
		return "", errors.New(`path must start with "/"`)
	}

	if path[len(path)-1] == '/' {
		return "", fmt.Errorf("failed to get path, path ends with directory not file: %s", path)
	}

	end := len(path) - 1

	for path[end] != '/' {
		end--
	}
	return string(path[:end+1]), nil
}
