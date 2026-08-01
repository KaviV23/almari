package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/KaviV23/almari/internal/service"
)

// DownloadHandler handles file download requests.
type DownloadHandler struct {
	svc *service.FileService
}

// NewDownloadHandler creates a new download handler.
func NewDownloadHandler(svc *service.FileService) *DownloadHandler {
	return &DownloadHandler{svc: svc}
}

// HandleDownload handles file downloads or directory zip downloads.
func (h *DownloadHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Query().Get("path")
	// We need to check if it's a file or directory first
	// For simplicity, we'll let the service handle it via http.ServeFile for files
	// and create zip for directories
	
	// Get the resolved path
	target, err := h.svc.Filesystem().ResolvePath(reqPath)
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

	if err := h.svc.CreateZipFromDir(reqPath, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}