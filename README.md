# Almari 🗄️

Almari is a lightweight, easy-to-deploy web application designed to manage a machine's filesystem. Deployed on a remote machine, it allows other machines on the network to securely view, navigate, upload, download, edit, and manage files and directories via a modern web browser.

---

## Features

- **Responsive Web Dashboard:** Built with Tailwind CSS and Alpine.js for a clean, fast, and intuitive user interface.
- **Zero Heavy Dependencies:** Written in Go with an embedded frontend — compiles into a single standalone binary.
- **Filesystem Operations:**
  - Directory listing with sorting and metadata (size, permissions, modification time).
  - File reading and in-browser text editing.
  - File and folder creation (`mkdir`).
  - Multipart file uploads.
  - File downloads & directory downloads (auto-zipped on the fly).
  - Secure deletion of files and directories.
- **Media Viewer:** In-browser modal viewer for images and videos with buffered HTTP Range streaming.
- **Thumbnails:** Automatically generated image and video thumbnails with caching. Video thumbnails require `ffmpeg` (falls back to a placeholder if unavailable).
- **Security:** Robust path sanitization to prevent path traversal attacks (`../` escapes).

---

## Quick Start

### Prerequisites
- Go 1.21+ (for building from source)
- **Optional:** `ffmpeg` — recommended for video thumbnail generation. If not installed, Almari generates a styled placeholder thumbnail instead.
  - **macOS:** `brew install ffmpeg`
  - **Ubuntu/Debian:** `sudo apt install ffmpeg`
  - **Windows:** Download from [ffmpeg.org](https://ffmpeg.org/download.html)

### Build & Run
```bash
# Clone the repository
git clone https://github.com/KaviV23/almari.git
cd almari

# Build binary
make build

# Run server (defaults to current directory on port 8080)
./bin/almari -port 8080 -dir /path/to/manage
```

### Command-Line Flags
- `-port`: Port to listen on (default `8080`)
- `-dir`: Root directory path to manage (default `.` current directory)

---

## Testing

Run the unit test suite (including path traversal security tests):
```bash
make test
```
