package path

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAllowedPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}
	absHome, err := filepath.Abs(home)
	if err != nil {
		t.Fatalf("failed to get abs home dir: %v", err)
	}

	absHomeFile := filepath.Join(absHome, "test_file.txt")
	absHomeSubDir := filepath.Join(absHome, "some", "subdir")

	tests := []struct {
		name    string
		path    string
		wantOK  bool
		wantAbs string
		wantErr bool
	}{
		{
			name:    "empty path (current dir)",
			path:    "",
			wantOK:  true,
			wantAbs: mustAbs("."),
			wantErr: false,
		},
		{
			name:    "home dir itself",
			path:    home,
			wantOK:  true,
			wantAbs: absHome,
			wantErr: false,
		},
		{
			name:    "file inside home",
			path:    filepath.Join(home, "test_file.txt"),
			wantOK:  true,
			wantAbs: absHomeFile,
			wantErr: false,
		},
		{
			name:    "subdir inside home",
			path:    filepath.Join(home, "some", "subdir"),
			wantOK:  true,
			wantAbs: absHomeSubDir,
			wantErr: false,
		},
		{
			name:    "relative path inside home",
			path:    "test_file.txt",
			wantOK:  true,
			wantAbs: mustAbs("test_file.txt"),
			wantErr: false,
		},
		{
			name:    "path outside home via ..",
			path:    filepath.Join(home, ".."),
			wantOK:  false,
			wantAbs: "",
			wantErr: true,
		},
		{
			name:    "root path",
			path:    "/",
			wantOK:  false,
			wantAbs: "",
			wantErr: true,
		},
		{
			name:    "another system dir",
			path:    "/tmp",
			wantOK:  false,
			wantAbs: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AllowedPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AllowedPath(%q) error = %v, wantErr = %v", tt.path, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if !strings.HasPrefix(got, absHome+string(filepath.Separator)) && got != absHome {
				t.Errorf("AllowedPath(%q) = %q, must be inside home %q", tt.path, got, absHome)
			}

			if tt.wantAbs != "" {
				if got != tt.wantAbs {
					t.Errorf("AllowedPath(%q) = %q, want %q", tt.path, got, tt.wantAbs)
				}
			}
		})
	}
}

func mustAbs(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}
	return abs
}
