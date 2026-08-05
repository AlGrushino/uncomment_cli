package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUncommentCmd(t *testing.T) {
	currDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	tempDir, err := os.MkdirTemp(currDir, "temp_dir_for_tests_*")
	if err != nil {
		t.Fatalf("failed to create temp dir for tests: %v", err)
	}

	tempPath, err := filepath.Abs(tempDir)
	if err != nil {
		t.Fatalf("failed to get abs path of temp dir: %v", err)
	}

	defer func(name string) {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Fatalf("failed to remove temp directory for tests %s: %v", tempDir, err)
		}
	}(tempPath)

	firstTestFile, err := os.CreateTemp(tempPath, "temp_file_*.go")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	secondTestFile, err := os.CreateTemp(tempPath, "temp_file_*.go")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			t.Fatalf("failed to close test file: %v", err)
		}

	}(firstTestFile)

	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			t.Fatalf("failed to close test file: %v", err)
		}

	}(secondTestFile)

	firstPath, err := filepath.Abs(firstTestFile.Name())
	if err != nil {
		t.Fatalf("failed to get temp file path: %v", err)
	}
	secondPath, err := filepath.Abs(secondTestFile.Name())
	if err != nil {
		t.Fatalf("failed to get temp file path: %v", err)
	}

	defer func(path string) {
		if err = os.RemoveAll(path); err != nil {
			t.Fatalf("failed to remove test file: %v", err)
		}
	}(firstPath)

	defer func(path string) {
		if err = os.RemoveAll(path); err != nil {
			t.Fatalf("failed to remove test file: %v", err)
		}
	}(secondPath)

	const mockText = `package main

import "fmt"

// This is a comment
func main() {
	// Another comment
	fmt.Println("Hello")
}
	`

	if err = os.WriteFile(firstPath, []byte(mockText), 0o600); err != nil {
		t.Fatalf("failed to write mock text into file: %v", err)
	}
	if err = os.WriteFile(secondPath, []byte(mockText), 0o600); err != nil {
		t.Fatalf("failed to write mock text into file: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		expected    string
		expectedErr error
	}{
		{
			name:        "with long flag",
			args:        []string{"--path", firstPath},
			expected:    "All comments in " + firstPath + " deleted\n",
			expectedErr: nil,
		},
		{
			name:        "with short flag",
			args:        []string{"-p", secondPath},
			expected:    "All comments in " + secondPath + " deleted\n",
			expectedErr: nil,
		},
		{
			name:        "empty path",
			args:        []string{"-p", ""},
			expected:    "",
			expectedErr: errors.New("Error: failed to parse file: open : no such file or directory"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(append([]string{"uncomment"}, tt.args...))

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)

			errBuf := new(bytes.Buffer)
			rootCmd.SetErr(errBuf)

			err := rootCmd.Execute()

			if tt.expectedErr != nil {
				gotErr := errBuf.String()
				if gotErr == "" {
					t.Error("expected error in stderr but got none")
				}
				if !strings.Contains(gotErr, "no such file or directory") {
					t.Errorf("stderr = %q, want to contain %q", gotErr, "no such file or directory")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if errBuf.Len() > 0 {
					t.Errorf("unexpected stderr: %q", errBuf.String())
				}
			}
		})
	}
}
