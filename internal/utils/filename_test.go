package utils

import (
	"testing"
)

func TestGetFileName(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantRes string
		wantErr bool
	}{
		{
			name:    "regular file",
			path:    "/some/new/path/file.txt",
			wantRes: "file.txt",
			wantErr: false,
		},
		{
			name:    "file without directory",
			path:    "file.txt",
			wantRes: "file.txt",
			wantErr: false,
		},
		{
			name:    "directory, no file",
			path:    "/some/new/path/",
			wantRes: "",
			wantErr: false,
		},
		{
			name:    "no source directory",
			path:    "some/new/path",
			wantRes: "path",
			wantErr: false,
		},
		{
			name:    "root path",
			path:    "/",
			wantRes: "",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantRes: "",
			wantErr: true,
		},
		{
			name:    "single file in root",
			path:    "/file.txt",
			wantRes: "file.txt",
			wantErr: false,
		},
		{
			name:    "single char file",
			path:    "a",
			wantRes: "a",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GetFileName(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileName(%q) error = %v, wantErr = %v", tt.path, err, tt.wantErr)
			}
			if res != tt.wantRes {
				t.Errorf("GetFileName(%q) = %q, want %q", tt.path, res, tt.wantRes)
			}
		})
	}
}
