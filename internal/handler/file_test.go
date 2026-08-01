package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KaviV23/almari/internal/model"
	"github.com/KaviV23/almari/internal/service"
)

func setupTestServer(t *testing.T) (*httptest.Server, *service.FileService) {
	fsSvc, err := service.NewFilesystemService(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create filesystem service: %v", err)
	}
	as := service.NewArchiveService()
	ts, err := service.NewThumbnailService()
	if err != nil {
		t.Fatalf("failed to create thumbnail service: %v", err)
	}
	fileSvc := service.NewFileService(fsSvc, as, ts)

	// Create test file
	if err := fileSvc.WriteFile("/hello.txt", "hello world"); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	fileHandler := NewFileHandler(fileSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/list", fileHandler.HandleList)
	mux.HandleFunc("/api/read", fileHandler.HandleRead)
	mux.HandleFunc("/api/write", fileHandler.HandleWrite)
	mux.HandleFunc("/api/mkdir", fileHandler.HandleMkdir)
	mux.HandleFunc("/api/delete", fileHandler.HandleDelete)

	return httptest.NewServer(mux), fileSvc
}

func TestFileHandlerList(t *testing.T) {
	ts, _ := setupTestServer(t)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/list?path=/")
	if err != nil {
		t.Fatalf("failed to GET /api/list: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var files []model.FileInfo
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
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
}

func TestFileHandlerRead(t *testing.T) {
	ts, _ := setupTestServer(t)
	defer ts.Close()

	// Write first
	writeBody := `{"path": "/newfile.txt", "content": "brand new content"}`
	resp, err := http.Post(ts.URL+"/api/write", "application/json", strings.NewReader(writeBody))
	if err != nil {
		t.Fatalf("failed to POST /api/write: %v", err)
	}
	resp.Body.Close()

	// Read
	resp, err = http.Get(ts.URL + "/api/read?path=/newfile.txt")
	if err != nil {
		t.Fatalf("failed to GET /api/read: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestFileHandlerMkdir(t *testing.T) {
	ts, _ := setupTestServer(t)
	defer ts.Close()

	mkdirBody := `{"path": "/newdir"}`
	resp, err := http.Post(ts.URL+"/api/mkdir", "application/json", strings.NewReader(mkdirBody))
	if err != nil {
		t.Fatalf("failed to POST /api/mkdir: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 on mkdir, got %d", resp.StatusCode)
	}
}