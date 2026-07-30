package utils

import (
	"errors"
)

func GetFileName(path string) (string, error) {
	if len(path) == 0 {
		return "", errors.New("len of path equals 0")
	}

	end := len(path) - 1
	res := []byte{}

	for i := end; i >= 0 && path[i] != '/'; i-- {
		res = append(res, path[i])
	}

	l, r := 0, len(res)-1

	for l < r {
		res[l], res[r] = res[r], res[l]
		l++
		r--
	}

	return string(res), nil
}
