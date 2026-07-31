package utils

import (
	"testing"
)

func TestGetFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantRes string
		wantErr bool
	}{
		{
			name:    "regular file",
			path:    "/some/new/path",
			wantRes: "/some/new/",
			wantErr: false,
		},
		{
			name:    "file in root",
			path:    "/file.txt",
			wantRes: "/",
			wantErr: false,
		},
		{
			name:    "directory, ends with /",
			path:    "/some/new/",
			wantRes: "",
			wantErr: true,
		},
		{
			name:    "no source directory",
			path:    "some/new/path",
			wantRes: "",
			wantErr: true,
		},
		{
			name:    "no directory at all",
			path:    "file.txt",
			wantRes: "",
			wantErr: true,
		},
		{
			name:    "root path",
			path:    "/",
			wantRes: "",
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			wantRes: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GetFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilePath(%q) error = %v, wantErr = %v", tt.path, err, tt.wantErr)
			}
			if res != tt.wantRes {
				t.Errorf("GetFilePath(%q) = %q, want %q", tt.path, res, tt.wantRes)
			}
		})
	}
}
