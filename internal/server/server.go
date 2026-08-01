package server

import (
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/KaviV23/almari/internal/handler"
	"github.com/KaviV23/almari/internal/service"
)

// Server holds all dependencies for the HTTP server.
type Server struct {
	fileHandler     *handler.FileHandler
	uploadHandler   *handler.UploadHandler
	downloadHandler *handler.DownloadHandler
	thumbHandler    *handler.ThumbnailHandler
	staticFS        fs.FS
}

// NewServer creates a new Server with all dependencies wired up.
func NewServer(rootPath string, staticFS fs.FS) (*Server, error) {
	fsSvc, err := service.NewFilesystemService(rootPath)
	if err != nil {
		return nil, err
	}
	as := service.NewArchiveService()
	ts, err := service.NewThumbnailService()
	if err != nil {
		return nil, err
	}
	fileSvc := service.NewFileService(fsSvc, as, ts)

	return &Server{
		fileHandler:     handler.NewFileHandler(fileSvc),
		uploadHandler:   handler.NewUploadHandler(fileSvc),
		downloadHandler: handler.NewDownloadHandler(fileSvc),
		thumbHandler:    handler.NewThumbnailHandler(fileSvc),
		staticFS:        staticFS,
	}, nil
}

// Routes returns an http.Handler with all routes registered.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/list", s.fileHandler.HandleList)
	mux.HandleFunc("/api/read", s.fileHandler.HandleRead)
	mux.HandleFunc("/api/write", s.fileHandler.HandleWrite)
	mux.HandleFunc("/api/mkdir", s.fileHandler.HandleMkdir)
	mux.HandleFunc("/api/delete", s.fileHandler.HandleDelete)
	mux.HandleFunc("/api/upload", s.uploadHandler.HandleUpload)
	mux.HandleFunc("/api/download", s.downloadHandler.HandleDownload)
	mux.HandleFunc("/api/thumbnail", s.thumbHandler.HandleThumbnail)

	// Static UI Routes - SPA fallback
	fileServer := http.FileServer(http.FS(s.staticFS))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		_, err := s.staticFS.Open(strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})

	return mux
}

// Start starts the HTTP server on the given address.
func (s *Server) Start(addr string) error {
	log.Printf("Starting Almari server on http://%s", addr)
	return http.ListenAndServe(addr, s.Routes())
}