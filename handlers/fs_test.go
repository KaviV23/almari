package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPathTraversalSecurity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "almari_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srv, err := NewServer(tmpDir)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	// Test malicious paths
	maliciousPaths := []string{
		"../etc/passwd",
		"/../../etc/passwd",
		"foo/../../bar",
	}

	for _, p := range maliciousPaths {
		_, err := srv.resolvePath(p)
		if err == nil {
			t.Errorf("expected security error for path %q, got nil", p)
		}
	}

	// Test valid paths
	validPaths := []string{
		"",
		"/",
		"test.txt",
		"sub/dir/file.txt",
	}

	for _, p := range validPaths {
		_, err := srv.resolvePath(p)
		if err != nil {
			t.Errorf("expected valid path for %q, got error: %v", p, err)
		}
	}
}

func TestServerHandlers(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "almari_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFilePath := filepath.Join(tmpDir, "hello.txt")
	if err := os.WriteFile(testFilePath, []byte("hello world"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	srv, err := NewServer(tmpDir)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/list", srv.HandleList)
	mux.HandleFunc("/api/read", srv.HandleRead)
	mux.HandleFunc("/api/write", srv.HandleWrite)
	mux.HandleFunc("/api/mkdir", srv.HandleMkdir)
	mux.HandleFunc("/api/delete", srv.HandleDelete)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	// 1. Test List
	resp, err := http.Get(ts.URL + "/api/list?path=/")
	if err != nil {
		t.Fatalf("failed to GET /api/list: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var files []FileInfo
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

	// 2. Test Read
	resp, err = http.Get(ts.URL + "/api/read?path=/hello.txt")
	if err != nil {
		t.Fatalf("failed to GET /api/read: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// 3. Test Write
	writeBody := `{"path": "/newfile.txt", "content": "brand new content"}`
	resp, err = http.Post(ts.URL+"/api/write", "application/json", strings.NewReader(writeBody))
	if err != nil {
		t.Fatalf("failed to POST /api/write: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 on write, got %d", resp.StatusCode)
	}

	// Verify newfile.txt content
	content, err := os.ReadFile(filepath.Join(tmpDir, "newfile.txt"))
	if err != nil || string(content) != "brand new content" {
		t.Errorf("file content mismatch or error reading file: %v", err)
	}

	// 4. Test Mkdir
	mkdirBody := `{"path": "/newdir"}`
	resp, err = http.Post(ts.URL+"/api/mkdir", "application/json", strings.NewReader(mkdirBody))
	if err != nil {
		t.Fatalf("failed to POST /api/mkdir: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 on mkdir, got %d", resp.StatusCode)
	}

	info, err := os.Stat(filepath.Join(tmpDir, "newdir"))
	if err != nil || !info.IsDir() {
		t.Errorf("failed to create directory /newdir")
	}

	// 5. Test Delete
	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/delete?path=/newfile.txt", nil)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to DELETE /api/delete: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 on delete, got %d", resp.StatusCode)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "newfile.txt")); !os.IsNotExist(err) {
		t.Errorf("expected newfile.txt to be deleted")
	}
}
