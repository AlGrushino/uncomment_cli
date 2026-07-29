package utils

import (
	"log"
	"os"
	"testing"
)

func TestIsFile(t *testing.T) {
	tempDirectory := "some_directory"

	tests := []struct {
		name    string
		path    string
		wantRes bool
		wantErr bool
	}{
		{"regular file", "./file.txt", true, false},
		{"not file 1", "./", false, false},
		{"file does not exist", "./another_file.txt", false, true},
		{"not a file 2", tempDirectory, false, false},
	}

	err := os.Mkdir(tempDirectory, 0700)
	if err != nil {
		log.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.Remove(tempDirectory)

	f, err := os.Create("file.txt")
	if err != nil {
		log.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(f.Name())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := IsFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf(`IsFile(%s) error = %v, wantErr = %v`, tt.path, err, tt.wantErr)
			}
			if res != tt.wantRes {
				t.Errorf(`IsFile(%s) = %t, want match for %t`, tt.path, res, tt.wantRes)
			}
		})
	}
}
