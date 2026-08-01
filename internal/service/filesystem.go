package service

import (
	"os"
	"path/filepath"
	"strings"
)

// FilesystemService handles filesystem operations with path safety.
type FilesystemService struct {
	RootPath string
}

// NewFilesystemService creates a new service with the given root path.
func NewFilesystemService(rootPath string) (*FilesystemService, error) {
	abs, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &os.PathError{Op: "stat", Path: abs, Err: os.ErrInvalid}
	}
	return &FilesystemService{RootPath: abs}, nil
}

// ResolvePath securely resolves a user-supplied path relative to RootPath,
// preventing path traversal attacks.
func (s *FilesystemService) ResolvePath(userPath string) (string, error) {
	userPath = strings.TrimPrefix(userPath, "/")

	fullPath := filepath.Join(s.RootPath, userPath)

	rel, err := filepath.Rel(s.RootPath, fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", os.ErrPermission
	}

	return fullPath, nil
}

// IsRoot checks if the given path is the root directory.
func (s *FilesystemService) IsRoot(path string) bool {
	return path == s.RootPath
}