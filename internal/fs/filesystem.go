package fs

import (
	"os"
	"time"
)

type FS interface {
	Stat(name string) (os.FileInfo, error)
}

type OSFS struct{}

func (OSFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

type MockFS struct {
	StatFunc func(string) (os.FileInfo, error)
}

func (f MockFS) Stat(name string) (os.FileInfo, error) {
	return f.StatFunc(name)
}

type MockFileInfo struct {
	Dir bool
}

func (m *MockFileInfo) Name() string       { return "" }
func (m *MockFileInfo) Size() int64        { return 0 }
func (m *MockFileInfo) Mode() os.FileMode  { return 0 }
func (m *MockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *MockFileInfo) IsDir() bool        { return m.Dir }
func (m *MockFileInfo) Sys() any           { return nil }
