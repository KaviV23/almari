package handler

import (
	"net/http"

	"github.com/KaviV23/almari/internal/service"
)

// ThumbnailHandler handles thumbnail requests.
type ThumbnailHandler struct {
	svc *service.FileService
}

// NewThumbnailHandler creates a new thumbnail handler.
func NewThumbnailHandler(svc *service.FileService) *ThumbnailHandler {
	return &ThumbnailHandler{svc: svc}
}

// HandleThumbnail serves cached or generated thumbnails.
func (h *ThumbnailHandler) HandleThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := r.URL.Query().Get("path")
	thumbPath, err := h.svc.GetThumbnail(reqPath)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, thumbPath)
}