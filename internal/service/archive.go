package service

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// ArchiveService handles archive operations (zip creation).
type ArchiveService struct{}

// NewArchiveService creates a new archive service.
func NewArchiveService() *ArchiveService {
	return &ArchiveService{}
}

// CreateZipFromDir creates a zip archive from a directory and writes it to the provided writer.
func (s *ArchiveService) CreateZipFromDir(sourceDir string, writer io.Writer) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if rel == "." {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = rel
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		entryWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(entryWriter, file)
		return err
	})
}