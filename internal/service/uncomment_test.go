package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUncomment(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get curr dir for tests: %v", err)
	}

	const tempDirPattern = "test_dir_*"
	tempDir, err := os.MkdirTemp(dir, tempDirPattern)
	if err != nil {
		t.Fatalf("failed to create temp dir for tests: %v", err)
	}

	tempDirPath, err := filepath.Abs(tempDir)
	if err != nil {
		t.Fatalf("failed te get abs path of temp dir: %v", err)
	}

	defer func(path string) {
		err = os.RemoveAll(path)
		if err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}(tempDirPath)

	const testFilePattern = "test_file_*.go"

	file, err := os.CreateTemp(tempDirPath, testFilePattern)

	tempFilePath, err := filepath.Abs(file.Name())
	if err != nil {
		t.Fatalf("failed to get temp file path: %v", err)
	}

	defer func(f *os.File) {
		if err = file.Close(); err != nil {
			t.Fatalf("failed to close temp file: %v", err)
		}
	}(file)

	defer func(path string) {
		err = os.RemoveAll(tempFilePath)
		if err != nil {
			t.Fatalf("failed to remove test file: %v", err)
		}
	}(tempDirPath)

	srcContent := `package main

import "fmt"

// This is a comment
func main() {
    // Another comment
    fmt.Println("Hello")
}
`

	srcPath, err := filepath.Abs(file.Name())
	if err != nil {
		t.Fatalf("failed to get abs path of test file: %v", err)
	}

	if err := os.WriteFile(srcPath, []byte(srcContent), 0o600); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	if err := Uncomment(srcPath); err != nil {
		t.Fatalf("Uncomment() error = %v", err)
	}

	got, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("failed to read result file: %v", err)
	}

	if strings.Contains(string(got), "// This is a comment") {
		t.Errorf("result still contains comment: %q", got)
	}
	if strings.Contains(string(got), "// Another comment") {
		t.Errorf("result still contains comment: %q", got)
	}

	if !strings.Contains(string(got), `fmt.Println("Hello")`) {
		t.Errorf("result lost code: %q", got)
	}
	if !strings.Contains(string(got), "func main()") {
		t.Errorf("result lost code: %q", got)
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
