package service

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
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

	var buf bytes.Buffer

	err = format.Node(&buf, fset, file)
	if err != nil {
		return fmt.Errorf("failed to format and write code into temp file: %w", err)
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

	_, err = io.Copy(dst, &buf)
	if err != nil {
		return fmt.Errorf("failed to copy src to dst: %w", err)
	}

	return nil
}
