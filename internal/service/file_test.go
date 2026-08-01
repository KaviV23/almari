package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilesystemServiceResolvePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "almari_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	svc, err := NewFilesystemService(tmpDir)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	maliciousPaths := []string{
		"../etc/passwd",
		"/../../etc/passwd",
		"foo/../../bar",
	}

	for _, p := range maliciousPaths {
		_, err := svc.ResolvePath(p)
		if err == nil {
			t.Errorf("expected security error for path %q, got nil", p)
		}
	}

	validPaths := []string{
		"",
		"/",
		"test.txt",
		"sub/dir/file.txt",
	}

	for _, p := range validPaths {
		_, err := svc.ResolvePath(p)
		if err != nil {
			t.Errorf("expected valid path for %q, got error: %v", p, err)
		}
	}
}

func TestFileServiceIntegration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "almari_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFilePath := filepath.Join(tmpDir, "hello.txt")
	if err := os.WriteFile(testFilePath, []byte("hello world"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	fsSvc, err := NewFilesystemService(tmpDir)
	if err != nil {
		t.Fatalf("failed to create filesystem service: %v", err)
	}
	as := NewArchiveService()
	ts, err := NewThumbnailService()
	if err != nil {
		t.Fatalf("failed to create thumbnail service: %v", err)
	}
	fileSvc := NewFileService(fsSvc, as, ts)

	// Test ListFiles
	files, err := fileSvc.ListFiles("/")
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}
	found := false
	for _, f := range files {
		if f.Name == "hello.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected hello.txt in file list")
	}

	// Test WriteFile
	if err := fileSvc.WriteFile("/newfile.txt", "brand new content"); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "newfile.txt"))
	if err != nil || string(content) != "brand new content" {
		t.Errorf("file content mismatch or error reading file: %v", err)
	}

	// Test Mkdir
	if err := fileSvc.Mkdir("/newdir"); err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	info, err := os.Stat(filepath.Join(tmpDir, "newdir"))
	if err != nil || !info.IsDir() {
		t.Errorf("failed to create directory /newdir")
	}

	// Test Delete
	if err := fileSvc.Delete("/newfile.txt"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "newfile.txt")); !os.IsNotExist(err) {
		t.Errorf("expected newfile.txt to be deleted")
	}
}