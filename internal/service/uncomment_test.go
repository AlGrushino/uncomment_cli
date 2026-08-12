package service

import (
	"fmt"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUncomment(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.go")

	srcContent := `package main

import "fmt"

// this is a type comment
type MyType int

// this is a struct comment
type MyStruct struct {
    // field comment
    Field int
}

// this is a var comment
var (
    // var1 comment
    Var1 int
    // var2 comment
    Var2 string
)

// this is a function comment
func main() {
    fmt.Println("Hello")
}
`
	if err := os.WriteFile(tempFile, []byte(srcContent), 0o600); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	if err := Uncomment(tempFile); err != nil {
		t.Fatalf("Uncomment() error = %v", err)
	}

	got, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("failed to read result file: %v", err)
	}
	gotStr := string(got)

	if strings.Contains(gotStr, "//") {
		t.Errorf("result still contains comments: %q", gotStr)
	}

	checks := []string{
		"type MyType int",
		"type MyStruct struct",
		"Field int",
		"Var1 int",
		"Var2 string",
		"fmt.Println",
	}
	for _, check := range checks {
		if !strings.Contains(gotStr, check) {
			t.Errorf("result lost code: %s", check)
		}
	}
}

func TestUncomment_Errors(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		setup   func() error
		wantErr bool
	}{
		{
			name: "file does not exist",
			path: filepath.Join(dir, "nonexistent.go"),
			setup: func() error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "invalid Go file",
			path: filepath.Join(dir, "invalid.go"),
			setup: func() error {
				return os.WriteFile(filepath.Join(dir, "invalid.go"), []byte("not a go file"), 0o600)
			},
			wantErr: true,
		},
		{
			name: "path outside allowed directory",
			path: "/etc/passwd",
			setup: func() error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "read-only file",
			path: filepath.Join(dir, "readonly.go"),
			setup: func() error {
				content := `package main\nfunc main() {}\n`
				if err := os.WriteFile(filepath.Join(dir, "readonly.go"), []byte(content), 0444); err != nil {
					return err
				}
				return nil
			},
			wantErr: true,
		},
		{
			name: "path is a directory",
			path: dir,
			setup: func() error {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(); err != nil {
				t.Fatalf("setup error = %v", err)
			}

			err := Uncomment(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uncomment() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
func TestUncomment_FormatError(t *testing.T) {
	old := formatNodeFunc
	defer func() { formatNodeFunc = old }()

	formatNodeFunc = func(w io.Writer, fset *token.FileSet, node any) error {
		return fmt.Errorf("mock format error")
	}

	filePath := createTestFile(t)
	err := Uncomment(filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestUncomment_AllowedPathError(t *testing.T) {
	old := allowedPathFunc
	defer func() { allowedPathFunc = old }()

	allowedPathFunc = func(path string) (string, error) {
		return "", fmt.Errorf("mock allowed path error")
	}

	filePath := createTestFile(t)
	err := Uncomment(filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "mock allowed path error") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestUncomment_OpenFileError(t *testing.T) {
	old := openFileFunc
	defer func() { openFileFunc = old }()

	openFileFunc = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, fmt.Errorf("mock open file error")
	}

	filePath := createTestFile(t)
	err := Uncomment(filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to open original") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestUncomment_CopyError(t *testing.T) {
	old := copyFunc
	defer func() { copyFunc = old }()

	copyFunc = func(dst io.Writer, src io.Reader) (int64, error) {
		return 0, fmt.Errorf("mock copy error")
	}

	filePath := createTestFile(t)
	err := Uncomment(filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to copy src to dst") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestUncomment_CloseError(t *testing.T) {
	t.Skip("skipping close error test: requires mocking *os.File.Close, which is complex")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || contains(s[1:], substr)))
}

func createTestFile(t *testing.T) string {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.go")
	content := `package main

import "fmt"

func main() {
    fmt.Println("Hello")
}
`
	if err := os.WriteFile(tempFile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return tempFile
}
