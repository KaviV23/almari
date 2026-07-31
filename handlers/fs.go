package handlers

import (
	"archive/zip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileInfo represents metadata about a file or directory.
type FileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"isDir"`
	ModTime int64  `json:"modTime"`
	Mode    string `json:"mode"`
}

// Server holds the configuration for the filesystem manager server.
type Server struct {
	RootPath string
}

// NewServer creates a new Server instance with the given root path.
func NewServer(rootPath string) (*Server, error) {
	abs, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}
	// Ensure root exists
	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &os.PathError{Op: "stat", Path: abs, Err: os.ErrInvalid}
	}
	return &Server{RootPath: abs}, nil
}

// resolvePath securely resolves a user-supplied path relative to RootPath,
// preventing path traversal attacks.
func (s *Server) resolvePath(userPath string) (string, error) {
	userPath = strings.TrimPrefix(userPath, "/")

	// Join with RootPath
	fullPath := filepath.Join(s.RootPath, userPath)

	// Ensure fullPath is within RootPath
	rel, err := filepath.Rel(s.RootPath, fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", os.ErrPermission
	}

	return fullPath, nil
}

// HandleList lists files and directories in the given path.
func (s *Server) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Query().Get("path")
	target, err := s.resolvePath(reqPath)
	if err != nil {
		http.Error(w, "Access denied or invalid path", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Path not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !info.IsDir() {
		http.Error(w, "Path is not a directory", http.StatusBadRequest)
		return
	}

	entries, err := os.ReadDir(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Compute relative path from RootPath for client navigation
		childFullPath := filepath.Join(target, entry.Name())
		rel, _ := filepath.Rel(s.RootPath, childFullPath)
		relPath := "/" + filepath.ToSlash(rel)

		files = append(files, FileInfo{
			Name:    entry.Name(),
			Path:    relPath,
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime().UnixMilli(),
			Mode:    info.Mode().String(),
		})
	}

	// Sort directories first, then files alphabetically
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// HandleRead reads the content of a text file.
func (s *Server) HandleRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Query().Get("path")
	target, err := s.resolvePath(reqPath)
	if err != nil {
		http.Error(w, "Access denied or invalid path", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(target)
	if err != nil || info.IsDir() {
		http.Error(w, "File not found or is a directory", http.StatusNotFound)
		return
	}

	file, err := os.Open(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.Copy(w, file)
}

// HandleWrite writes content to a file.
func (s *Server) HandleWrite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	target, err := s.resolvePath(req.Path)
	if err != nil {
		http.Error(w, "Access denied or invalid path", http.StatusBadRequest)
		return
	}

	if err := os.WriteFile(target, []byte(req.Content), 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

// HandleMkdir creates a new directory.
func (s *Server) HandleMkdir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	target, err := s.resolvePath(req.Path)
	if err != nil {
		http.Error(w, "Access denied or invalid path", http.StatusBadRequest)
		return
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

// HandleDelete deletes a file or directory.
func (s *Server) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Query().Get("path")
	if reqPath == "" {
		var req struct {
			Path string `json:"path"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		reqPath = req.Path
	}

	target, err := s.resolvePath(reqPath)
	if err != nil || target == s.RootPath {
		http.Error(w, "Access denied or cannot delete root", http.StatusBadRequest)
		return
	}

	if err := os.RemoveAll(target); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

// HandleUpload handles multipart file uploads.
func (s *Server) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	destPath := r.FormValue("path")
	targetDir, err := s.resolvePath(destPath)
	if err != nil {
		http.Error(w, "Access denied or invalid path", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			continue
		}
		defer f.Close()

		destFilePath := filepath.Join(targetDir, fh.Filename)
		// Ensure file stays within targetDir
		if !strings.HasPrefix(destFilePath, targetDir) {
			continue
		}

		dst, err := os.Create(destFilePath)
		if err != nil {
			continue
		}
		defer dst.Close()

		io.Copy(dst, f)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

// HandleDownload handles file downloads or directory zip downloads.
func (s *Server) HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Query().Get("path")
	target, err := s.resolvePath(reqPath)
	if err != nil {
		http.Error(w, "Access denied or invalid path", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(target)
	if err != nil {
		http.Error(w, "Path not found", http.StatusNotFound)
		return
	}

	if !info.IsDir() {
		http.ServeFile(w, r, target)
		return
	}

	// If directory, zip it on the fly
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filepath.Base(target)+`.zip"`)

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(target, path)
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

		writer, err := zipWriter.CreateHeader(header)
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

		_, err = io.Copy(writer, file)
		return err
	})
}
