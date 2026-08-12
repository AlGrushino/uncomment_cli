package service

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func checkNoComments(path string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", path, err)
	}

	if len(file.Comments) > 0 {
		return fmt.Errorf("file %s contains upper-level comments", path)
	}

	var found bool
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		switch node := n.(type) {
		case *ast.GenDecl:
			if node.Doc != nil {
				found = true
			}
		case *ast.FuncDecl:
			if node.Doc != nil {
				found = true
			}
		case *ast.TypeSpec:
			if node.Doc != nil || node.Comment != nil {
				found = true
			}
		case *ast.ValueSpec:
			if node.Doc != nil || node.Comment != nil {
				found = true
			}
		case *ast.Field:
			if node.Doc != nil || node.Comment != nil {
				found = true
			}
		case *ast.ImportSpec:
			if node.Doc != nil || node.Comment != nil {
				found = true
			}
		}
		return !found
	})
	if found {
		return fmt.Errorf("file %s contains comments in AST nodes", path)
	}
	return nil
}

func TestUncommentMany(t *testing.T) {
	cores := runtime.NumCPU()
	tempDir := t.TempDir()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get path to current file")
	}
	testdataDir := filepath.Join(filepath.Dir(filename), "testdata")

	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatal(err)
	}

	fileNames := []string{}

	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}

	if len(fileNames) == 0 {
		t.Fatal("no files in testdata")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	copiedPaths := make([]string, 0, cores)

	for range cores {
		idx := rng.Intn(len(fileNames))
		srcName := fileNames[idx]
		srcPath := filepath.Join(testdataDir, srcName)

		content, err := os.ReadFile(srcPath)
		if err != nil {
			t.Fatal(err)
		}

		tmpFile, err := os.CreateTemp(tempDir, filepath.Base(srcName)+"-*.go")
		if err != nil {
			t.Fatal(err)
		}
		dstPath := tmpFile.Name()
		err = tmpFile.Close()
		if err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(dstPath, content, 0644); err != nil {
			t.Fatal(err)
		}

		copiedPaths = append(copiedPaths, dstPath)
	}

	resCh := UncommentMany(copiedPaths...)

	errorsList := []error{}
	for pair := range resCh {
		if pair.err != nil {
			errorsList = append(errorsList, fmt.Errorf("failed %s: %w", pair.path, pair.err))
		}
	}

	if len(errorsList) > 0 {
		for _, e := range errorsList {
			t.Error(e)
		}
		t.Fatalf("failed in handling file")
	}

	for _, path := range copiedPaths {
		if err := checkNoComments(path); err != nil {
			t.Errorf("comments found: %v", err)
		}
	}
}

func TestUncommentManyMoreFilesThanCores(t *testing.T) {
	cores := runtime.NumCPU()
	filesCount := max(cores*3+1, 10)

	t.Logf("cores=%d, files to process=%d", cores, filesCount)

	tempDir := t.TempDir()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get path to current file")
	}
	testdataDir := filepath.Join(filepath.Dir(filename), "testdata")

	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no files in testdata")
	}

	fileNames := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	if len(fileNames) == 0 {
		t.Fatal("no regular files in testdata")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	copiedPaths := make([]string, 0, filesCount)

	for range filesCount {
		idx := rng.Intn(len(fileNames))
		srcName := fileNames[idx]
		srcPath := filepath.Join(testdataDir, srcName)

		content, err := os.ReadFile(srcPath)
		if err != nil {
			t.Fatal(err)
		}

		tmpFile, err := os.CreateTemp(tempDir, filepath.Base(srcName)+"-*.go")
		if err != nil {
			t.Fatal(err)
		}
		dstPath := tmpFile.Name()
		if err := tmpFile.Close(); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(dstPath, content, 0644); err != nil {
			t.Fatal(err)
		}

		copiedPaths = append(copiedPaths, dstPath)
	}

	resCh := UncommentMany(copiedPaths...)

	errorsList := []error{}
	for pair := range resCh {
		if pair.err != nil {
			errorsList = append(errorsList, fmt.Errorf("failed %s: %w", pair.path, pair.err))
		}
	}

	if len(errorsList) > 0 {
		for _, e := range errorsList {
			t.Error(e)
		}
		t.Fatalf("failed in handling %d files", len(errorsList))
	}

	for _, path := range copiedPaths {
		if err := checkNoComments(path); err != nil {
			t.Errorf("comments found in %s: %v", path, err)
		}
	}
}
