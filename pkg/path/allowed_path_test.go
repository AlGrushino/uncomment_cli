package path

import (
	"errors"
	"os"
	"path/filepath"
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
			name:    "temporary directory",
			path:    os.TempDir(),
			wantOK:  true,
			wantAbs: mustAbs(os.TempDir()),
			wantErr: false,
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

			if tt.wantAbs != "" && got != tt.wantAbs {
				t.Errorf("AllowedPath(%q) = %q, want %q", tt.path, got, tt.wantAbs)
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

func TestAllowedPath_AbsPathError(t *testing.T) {
	oldAbs := absFunc
	defer func() { absFunc = oldAbs }()

	testPath := "some/path"
	absFunc = func(p string) (string, error) {
		if p == testPath {
			return "", errors.New("mock abs error")
		}
		return filepath.Abs(p)
	}

	_, err := AllowedPath(testPath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to get absolute path") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAllowedPath_HomeDirError(t *testing.T) {
	oldHome := homeDirFunc
	defer func() { homeDirFunc = oldHome }()

	homeDirFunc = func() (string, error) {
		return "", errors.New("mock home dir error")
	}

	_, err := AllowedPath("/some/path")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to get home dir") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAllowedPath_AbsHomeError(t *testing.T) {
	oldAbs := absFunc
	oldHome := homeDirFunc
	defer func() {
		absFunc = oldAbs
		homeDirFunc = oldHome
	}()

	fakeHome := "/home/testuser"
	homeDirFunc = func() (string, error) {
		return fakeHome, nil
	}

	absFunc = func(p string) (string, error) {
		if p == fakeHome {
			return "", errors.New("mock abs home error")
		}
		return filepath.Abs(p)
	}

	_, err := AllowedPath("/some/path")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to get absolute home dir") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAllowedPath_AbsTempError(t *testing.T) {
	oldAbs := absFunc
	oldHome := homeDirFunc
	oldTemp := tempDirFunc
	defer func() {
		absFunc = oldAbs
		homeDirFunc = oldHome
		tempDirFunc = oldTemp
	}()

	fakeHome := "/home/testuser"
	fakeTmp := "/tmp"
	homeDirFunc = func() (string, error) {
		return fakeHome, nil
	}
	tempDirFunc = func() string {
		return fakeTmp
	}

	absFunc = func(p string) (string, error) {
		if p == fakeTmp {
			return "", errors.New("mock abs tmp error")
		}
		return filepath.Abs(p)
	}

	_, err := AllowedPath("/some/outside/path")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "failed to get absolute temp dir") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || contains(s[1:], substr)))
}
