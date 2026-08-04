package service

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"uncomment-cli/internal/utils"
	"uncomment-cli/pkg/path"
)

func Uncomment(originalPath string) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, originalPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	file.Comments = nil

	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch node := n.(type) {
		case *ast.GenDecl:
			node.Doc = nil
		case *ast.FuncDecl:
			node.Doc = nil
		case *ast.TypeSpec:
			node.Doc = nil
			node.Comment = nil
		case *ast.ValueSpec:
			node.Doc = nil
			node.Comment = nil
		case *ast.Field:
			node.Doc = nil
			node.Comment = nil
		}
		return true
	})

	newPath, err := utils.GetFilePath(originalPath)
	if err != nil {
		return fmt.Errorf("failed to get file path: %w", err)
	}

	outputFile, err := os.CreateTemp(newPath, "temp_file_*.go")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(outputFile.Name())

	err = format.Node(outputFile, fset, file)
	if err != nil {
		return fmt.Errorf("failed to format and write code into temp file: %w", err)
	}

	if _, err := outputFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek temp file: %w", err)
	}

	err = os.Truncate(originalPath, 0)
	if err != nil {
		return fmt.Errorf("failed to clear original file: %w", err)
	}

	safePath, err := path.AllowedPath(originalPath)
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(safePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600) // #nosec G304
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
