package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUncomment(t *testing.T) {
	testFileName := "test_file.go"

	file, err := os.Create(testFileName)
	if err != nil {
		t.Fatalf("failed to create temp file for tests: %v", err)
	}
	defer file.Close()
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			t.Fatalf("failed to remove test_file.go: %v", err)
		}
	}(testFileName)

	srcContent := `package main

import "fmt"

// This is a comment
func main() {
    // Another comment
    fmt.Println("Hello")
}
`

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get curr dir for tests: %v", err)
	}

	srcPath := filepath.Join(dir, testFileName)
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

	gotStr := string(got)

	if strings.Contains(gotStr, "// This is a comment") {
		t.Errorf("result still contains comment: %q", gotStr)
	}
	if strings.Contains(gotStr, "// Another comment") {
		t.Errorf("result still contains comment: %q", gotStr)
	}

	if !strings.Contains(gotStr, `fmt.Println("Hello")`) {
		t.Errorf("result lost code: %q", gotStr)
	}
	if !strings.Contains(gotStr, "func main()") {
		t.Errorf("result lost code: %q", gotStr)
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
