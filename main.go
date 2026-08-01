package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"

	"github.com/KaviV23/almari/internal/server"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	dir := flag.String("dir", ".", "Root directory to manage")
	flag.Parse()

	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	srv, err := server.NewServer(*dir, subFS)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	log.Printf("Managing filesystem root: %s", *dir)

	if err := srv.Start(":" + *port); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}