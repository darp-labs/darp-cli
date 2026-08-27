package init

import (
	"errors"
	"os"
	"path/filepath"
)

type osFileSystem struct{}

// NewOSFileSystem creates a filesystem backed by the local operating system.
func NewOSFileSystem() FileSystem {
	return osFileSystem{}
}

func (osFileSystem) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (osFileSystem) MkdirAll(path string) error {
	return os.MkdirAll(path, 0o755)
}

func (osFileSystem) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}

func (osFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (osFileSystem) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func (osFileSystem) Remove(path string) error {
	return os.Remove(path)
}

func (osFileSystem) Base(path string) string {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return filepath.Base(path)
	}
	return filepath.Base(absolutePath)
}
