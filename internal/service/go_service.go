package service

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"uncomment-cli/internal/utils"
)

func Uncomment(originalPath string) error {
	filename, err := utils.GetFileName(originalPath)

	if err != nil {
		return fmt.Errorf("failed to get name of file: %w", err)
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, originalPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	file.Comments = nil

	newPath, err := utils.GetFilePath(originalPath)
	if err != nil {
		return fmt.Errorf("failed to get file path: %w", err)
	}

	outputFile, err := os.CreateTemp(newPath, filename)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(outputFile.Name())

	err = format.Node(outputFile, fset, file)
	if err != nil {
		return fmt.Errorf("failed to format and write code into temp file: %w", err)
	}

	err = os.Truncate(originalPath, 0)
	if err != nil {
		return fmt.Errorf("failed to clear original file: %w", err)
	}

	dst, err := os.Open(originalPath)
	if err != nil {
		return fmt.Errorf("failed to open original: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, outputFile)
	if err != nil {
		return fmt.Errorf("failed to copy src to dst: %w", err)
	}

	return nil
}
