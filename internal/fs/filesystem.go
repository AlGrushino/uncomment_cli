// Package fs defines a small filesystem abstraction used to keep file checks
// testable with mocks.
package fs

import (
	"os"
	"time"
)

var _ FS = OSFS{}
var _ FS = MockFS{}

// FS abstracts access to the filesystem. It can be backed by OSFS for real
// files or by MockFS in tests.
type FS interface {
	Stat(name string) (os.FileInfo, error)
}

// OSFS is an FS implementation backed by the real os package.
type OSFS struct{}

// Stat returns the os.FileInfo for the given name.
func (OSFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// MockFS is an FS implementation that delegates Stat to a user-provided
// function, which makes it useful in tests.
type MockFS struct {
	StatFunc func(string) (os.FileInfo, error)
}

// Stat calls the configured StatFunc.
func (f MockFS) Stat(name string) (os.FileInfo, error) {
	return f.StatFunc(name)
}

// MockFileInfo is a minimal os.FileInfo implementation for tests. All methods
// except IsDir return zero values.
type MockFileInfo struct {
	Dir bool
}

// Name returns an empty base name.
func (m *MockFileInfo) Name() string { return "" }

// Size returns 0.
func (m *MockFileInfo) Size() int64 { return 0 }

// Mode returns the zero file mode.
func (m *MockFileInfo) Mode() os.FileMode { return 0 }

// ModTime returns the zero time.
func (m *MockFileInfo) ModTime() time.Time { return time.Time{} }

// IsDir reports whether the mocked entry is a directory.
func (m *MockFileInfo) IsDir() bool { return m.Dir }

// Sys returns nil.
func (m *MockFileInfo) Sys() any { return nil }
