// Package service contains the core logic for stripping comments from Go
// source files, either one file at a time or many files at once.
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

	"github.com/AlGrushino/uncomment_cli/pkg/path"
)

var (
	parseFileFunc   = parser.ParseFile
	formatNodeFunc  = format.Node
	allowedPathFunc = path.AllowedPath
	openFileFunc    = os.OpenFile
	copyFunc        = io.Copy
)

// Uncomment removes all comments from the Go source file at originalPath
// and writes the cleaned source back to the same file.
//
// The file is parsed into an AST, every comment node is dropped and the
// resulting source is reformatted before the original file is overwritten.
// The path must point inside the home or temp directory (see path.AllowedPath).
func Uncomment(originalPath string) (err error) {
	fset := token.NewFileSet()

	file, err := parseFileFunc(fset, originalPath, nil, parser.ParseComments)
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

	err = formatNodeFunc(&buf, fset, file)
	if err != nil {
		return fmt.Errorf("failed to format and write code into temp file: %w", err)
	}

	safePath, err := allowedPathFunc(originalPath)
	if err != nil {
		return err
	}

	dst, err := openFileFunc(safePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open original: %w", err)
	}
	defer func(f *os.File) {
		if closeErr := dst.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %v", closeErr)
		}
	}(dst)

	_, err = copyFunc(dst, &buf)
	if err != nil {
		return fmt.Errorf("failed to copy src to dst: %w", err)
	}

	return nil
}
