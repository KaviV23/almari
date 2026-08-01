package service

import (
	"github.com/KaviV23/almari/internal/model"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileService provides high-level file operations combining filesystem and archive services.
type FileService struct {
	fs *FilesystemService
	as *ArchiveService
	ts *ThumbnailService
}

// Filesystem returns the underlying filesystem service.
func (s *FileService) Filesystem() *FilesystemService {
	return s.fs
}

// NewFileService creates a new file service.
func NewFileService(fs *FilesystemService, as *ArchiveService, ts *ThumbnailService) *FileService {
	return &FileService{fs: fs, as: as, ts: ts}
}

// ListFiles lists files and directories in the given path.
func (s *FileService) ListFiles(reqPath string) ([]model.FileInfo, error) {
	target, err := s.fs.ResolvePath(reqPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(target)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &os.PathError{Op: "list", Path: target, Err: os.ErrInvalid}
	}

	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, err
	}

	var files []model.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		childFullPath := filepath.Join(target, entry.Name())
		rel, _ := filepath.Rel(s.fs.RootPath, childFullPath)
		relPath := "/" + filepath.ToSlash(rel)

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		isImg := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
		isVideo := ext == ".mp4" || ext == ".mkv" || ext == ".mov" || ext == ".webm" || ext == ".avi"
		hasThumb := !entry.IsDir() && (isImg || isVideo)

		files = append(files, model.FileInfo{
			Name:     entry.Name(),
			Path:     relPath,
			Size:     info.Size(),
			IsDir:    entry.IsDir(),
			ModTime:  info.ModTime().UnixMilli(),
			Mode:     info.Mode().String(),
			HasThumb: hasThumb,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	return files, nil
}

// ReadFile reads a file's content.
func (s *FileService) ReadFile(reqPath string) (*os.File, error) {
	target, err := s.fs.ResolvePath(reqPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(target)
	if err != nil || info.IsDir() {
		return nil, &os.PathError{Op: "read", Path: target, Err: os.ErrInvalid}
	}

	return os.Open(target)
}

// WriteFile writes content to a file.
func (s *FileService) WriteFile(reqPath, content string) error {
	target, err := s.fs.ResolvePath(reqPath)
	if err != nil {
		return err
	}
	return os.WriteFile(target, []byte(content), 0644)
}

// Mkdir creates a directory.
func (s *FileService) Mkdir(reqPath string) error {
	target, err := s.fs.ResolvePath(reqPath)
	if err != nil {
		return err
	}
	return os.MkdirAll(target, 0755)
}

// Delete removes a file or directory.
func (s *FileService) Delete(reqPath string) error {
	target, err := s.fs.ResolvePath(reqPath)
	if err != nil || s.fs.IsRoot(target) {
		return &os.PathError{Op: "delete", Path: reqPath, Err: os.ErrPermission}
	}
	return os.RemoveAll(target)
}

// UploadFiles uploads files to a directory.
func (s *FileService) UploadFiles(destPath string, files map[string]io.Reader) error {
	targetDir, err := s.fs.ResolvePath(destPath)
	if err != nil {
		return err
	}

	for filename, reader := range files {
		destFilePath := filepath.Join(targetDir, filename)
		if !strings.HasPrefix(destFilePath, targetDir) {
			continue
		}

		dst, err := os.Create(destFilePath)
		if err != nil {
			continue
		}
		defer dst.Close()

		io.Copy(dst, reader)
	}
	return nil
}

// CreateZipFromDir creates a zip from a directory.
func (s *FileService) CreateZipFromDir(reqPath string, writer interface{}) error {
	target, err := s.fs.ResolvePath(reqPath)
	if err != nil {
		return err
	}
	return s.as.CreateZipFromDir(target, writer.(interface{ Write([]byte) (int, error) }))
}

// GetThumbnail returns a thumbnail for a file.
func (s *FileService) GetThumbnail(reqPath string) (string, error) {
	target, err := s.fs.ResolvePath(reqPath)
	if err != nil {
		return "", err
	}
	return s.ts.GetOrGenerateThumbnail(target)
}