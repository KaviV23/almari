package model

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