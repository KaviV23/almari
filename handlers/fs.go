package handlers

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/gif"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// FileInfo represents metadata about a file or directory.
type FileInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"isDir"`
	ModTime  int64  `json:"modTime"`
	Mode     string `json:"mode"`
	HasThumb bool   `json:"hasThumb"`
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

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		isImg := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
		isVideo := ext == ".mp4" || ext == ".mkv" || ext == ".mov" || ext == ".webm" || ext == ".avi"
		hasThumb := !entry.IsDir() && (isImg || isVideo)

		files = append(files, FileInfo{
			Name:     entry.Name(),
			Path:     relPath,
			Size:     info.Size(),
			IsDir:    entry.IsDir(),
			ModTime:  info.ModTime().UnixMilli(),
			Mode:     info.Mode().String(),
			HasThumb: hasThumb,
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

// HandleThumbnail serves cached or generated thumbnails for images and videos.
func (s *Server) HandleThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Query().Get("path")
	target, err := s.resolvePath(reqPath)
	if err != nil {
		http.Error(w, "Access denied", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(target)
	if err != nil || info.IsDir() {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Cache directory setup in OS UserCacheDir or TempDir
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	thumbDir := filepath.Join(cacheDir, "almari", "thumbnails")
	os.MkdirAll(thumbDir, 0755)

	// Unique cache key based on absolute path and modification time
	hasher := sha256.New()
	hasher.Write([]byte(target))
	hasher.Write([]byte(fmt.Sprintf("%d", info.ModTime().UnixNano())))
	cacheKey := hex.EncodeToString(hasher.Sum(nil)) + ".jpg"
	cacheFilePath := filepath.Join(thumbDir, cacheKey)

	// Serve from cache if exists
	if cachedInfo, err := os.Stat(cacheFilePath); err == nil && !cachedInfo.IsDir() {
		http.ServeFile(w, r, cacheFilePath)
		return
	}

	// Generate thumbnail
	ext := strings.ToLower(filepath.Ext(target))
	isImg := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
	isVideo := ext == ".mp4" || ext == ".mkv" || ext == ".mov" || ext == ".webm" || ext == ".avi"

	var genErr error
	if isImg {
		genErr = generateImageThumb(target, cacheFilePath, 250, 250)
	} else if isVideo {
		genErr = generateVideoThumb(target, cacheFilePath, 250, 250)
	} else {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	if genErr != nil {
		http.Error(w, "Failed to generate thumbnail: "+genErr.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, cacheFilePath)
}

func generateImageThumb(srcPath, dstPath string, maxWidth, maxHeight int) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	// Calculate aspect ratio
	w := maxWidth
	h := maxHeight
	if srcW > srcH {
		h = int(float64(srcH) * float64(maxWidth) / float64(srcW))
	} else {
		w = int(float64(srcW) * float64(maxHeight) / float64(srcH))
	}
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}

	// Pure Go nearest-neighbor scaling
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			srcX := bounds.Min.X + (x*srcW)/w
			srcY := bounds.Min.Y + (y*srcH)/h
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, dst, &jpeg.Options{Quality: 80})
}

func generateVideoThumb(srcPath, dstPath string, width, height int) error {
	cmd := exec.Command("ffmpeg", "-ss", "00:00:01", "-i", srcPath, "-vframes", "1", "-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", width, height), "-f", "image2", dstPath)
	if err := cmd.Run(); err == nil {
		if _, err := os.Stat(dstPath); err == nil {
			return nil
		}
	}
	// Fallback to pure Go styled placeholder thumbnail if ffmpeg is unavailable or fails
	return generatePlaceholderThumb(dstPath, width, height)
}

func generatePlaceholderThumb(dstPath string, width, height int) error {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	bgCol := color.RGBA{30, 41, 59, 255} // slate-800
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(x, y, bgCol)
		}
	}

	// Draw a play button triangle
	cx, cy := width/2, height/2
	halfSize := 35
	for y := cy - halfSize; y <= cy + halfSize; y++ {
		for x := cx - halfSize/2; x <= cx + halfSize; x++ {
			dx := float64(x - (cx - halfSize/2))
			dy := float64(y - cy)
			if dx >= 0 && dx <= float64(halfSize)*1.5 && math.Abs(dy) <= dx*0.75 {
				dst.Set(x, y, color.RGBA{255, 255, 255, 230})
			}
		}
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, dst, &jpeg.Options{Quality: 80})
}
