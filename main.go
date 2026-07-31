package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"

	"github.com/KaviV23/almari/handlers"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	dir := flag.String("dir", ".", "Root directory to manage")
	flag.Parse()

	srv, err := handlers.NewServer(*dir)
	if err != nil {
		log.Fatalf("Failed to initialize server root %q: %v", *dir, err)
	}

	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/list", srv.HandleList)
	mux.HandleFunc("/api/read", srv.HandleRead)
	mux.HandleFunc("/api/write", srv.HandleWrite)
	mux.HandleFunc("/api/mkdir", srv.HandleMkdir)
	mux.HandleFunc("/api/delete", srv.HandleDelete)
	mux.HandleFunc("/api/upload", srv.HandleUpload)
	mux.HandleFunc("/api/download", srv.HandleDownload)
	mux.HandleFunc("/api/thumbnail", srv.HandleThumbnail)

	// Static UI Routes
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}
	fileServer := http.FileServer(http.FS(subFS))
	mux.Handle("/", fileServer)

	log.Printf("Starting Almari server on http://localhost:%s", *port)
	log.Printf("Managing filesystem root: %s", srv.RootPath)

	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
