package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AlGrushino/uncomment_cli/internal/fs"
)

func TestIsFile(t *testing.T) {
	curr, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get curr working dir: %v", err)
	}

	tempDirectory, err := os.MkdirTemp(curr, "temp_dir_pattern*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	tempPath, err := filepath.Abs(tempDirectory)
	if err != nil {
		t.Fatalf("failed to get temp path: %v", err)
	}

	defer func(name string) {
		err = os.RemoveAll(name)
		if err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}(tempPath)

	tempFile, err := os.CreateTemp(tempPath, "temp_file.txt")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tepmFilePath, err := filepath.Abs(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to get temp file path: %v", err)
	}

	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			t.Fatalf("failed to close test file: %v", err)
		}
	}(tempFile)

	defer func(name string) {
		err = os.RemoveAll(name)
		if err != nil {
			t.Fatalf("failed to remove temp file: %v", err)
		}
	}(tepmFilePath)

	tests := []struct {
		name    string
		fs      fs.FS
		path    string
		wantRes bool
		wantErr bool
	}{
		{
			name:    "regular file",
			fs:      fs.OSFS{},
			path:    tepmFilePath,
			wantRes: true,
			wantErr: false,
		},
		{
			name:    "not file 1",
			fs:      fs.OSFS{},
			path:    "./",
			wantRes: false,
			wantErr: false},
		{
			name:    "file does not exist",
			fs:      fs.OSFS{},
			path:    "./another_file.txt",
			wantRes: false,
			wantErr: true,
		},
		{
			name: "not a file mock dir",
			fs: fs.MockFS{
				StatFunc: func(name string) (os.FileInfo, error) {
					return &fs.MockFileInfo{Dir: true}, nil
				},
			},
			path:    tempDirectory,
			wantRes: false,
			wantErr: false,
		},
		{
			name: "not a file mock regular",
			fs: fs.MockFS{
				StatFunc: func(s string) (os.FileInfo, error) {
					return &fs.MockFileInfo{Dir: false}, nil
				},
			},
			path:    "any",
			wantRes: true,
			wantErr: false,
		},
		{
			name: "stat error",
			fs: fs.MockFS{
				StatFunc: func(s string) (os.FileInfo, error) {
					return nil, os.ErrNotExist
				},
			},
			path:    "missing",
			wantRes: false,
			wantErr: true,
		},
		{
			name: "other stat error",
			fs: fs.MockFS{
				StatFunc: func(string) (os.FileInfo, error) {
					return nil, os.ErrPermission
				},
			},
			path:    "some",
			wantRes: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := IsFile(tt.fs, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf(`IsFile(%s) error = %v, wantErr = %v`, tt.path, err, tt.wantErr)
			}
			if res != tt.wantRes {
				t.Errorf(`IsFile(%s) = %t, want match for %t`, tt.path, res, tt.wantRes)
			}
		})
	}
}
