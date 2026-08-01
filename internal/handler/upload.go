package handler

import (
	"io"
	"net/http"

	"github.com/KaviV23/almari/internal/service"
)

// UploadHandler handles file upload requests.
type UploadHandler struct {
	svc *service.FileService
}

// NewUploadHandler creates a new upload handler.
func NewUploadHandler(svc *service.FileService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

// HandleUpload handles multipart file uploads.
func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	destPath := r.FormValue("path")
	files := make(map[string]io.Reader)

	for _, fh := range r.MultipartForm.File["files"] {
		f, err := fh.Open()
		if err != nil {
			continue
		}
		defer f.Close()
		files[fh.Filename] = f
	}

	if err := h.svc.UploadFiles(destPath, files); err != nil {
		http.Error(w, "Access denied or invalid path", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}