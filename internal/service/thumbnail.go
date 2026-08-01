package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/gif"
	_ "image/png"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ThumbnailService handles thumbnail generation and caching.
type ThumbnailService struct {
	CacheDir string
	MaxSize  int
}

// NewThumbnailService creates a new thumbnail service.
func NewThumbnailService() (*ThumbnailService, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	thumbDir := filepath.Join(cacheDir, "almari", "thumbnails")
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return nil, err
	}
	return &ThumbnailService{
		CacheDir: thumbDir,
		MaxSize:  250,
	}, nil
}

// GetOrGenerateThumbnail returns the path to a cached thumbnail or generates a new one.
func (s *ThumbnailService) GetOrGenerateThumbnail(filePath string) (string, error) {
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write([]byte(filePath))
	hasher.Write([]byte(fmt.Sprintf("%d", info.ModTime().UnixNano())))
	cacheKey := hex.EncodeToString(hasher.Sum(nil)) + ".jpg"
	cacheFilePath := filepath.Join(s.CacheDir, cacheKey)

	if cachedInfo, err := os.Stat(cacheFilePath); err == nil && !cachedInfo.IsDir() {
		return cacheFilePath, nil
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	isImg := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
	isVideo := ext == ".mp4" || ext == ".mkv" || ext == ".mov" || ext == ".webm" || ext == ".avi"

	var genErr error
	if isImg {
		genErr = s.generateImageThumb(filePath, cacheFilePath)
	} else if isVideo {
		genErr = s.generateVideoThumb(filePath, cacheFilePath)
	} else {
		return "", fmt.Errorf("unsupported file type")
	}

	if genErr != nil {
		return "", genErr
	}

	return cacheFilePath, nil
}

func (s *ThumbnailService) generateImageThumb(srcPath, dstPath string) error {
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

	w := s.MaxSize
	h := s.MaxSize
	if srcW > srcH {
		h = int(float64(srcH) * float64(s.MaxSize) / float64(srcW))
	} else {
		w = int(float64(srcW) * float64(s.MaxSize) / float64(srcH))
	}
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}

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

func (s *ThumbnailService) generateVideoThumb(srcPath, dstPath string) error {
	cmd := exec.Command("ffmpeg", "-ss", "00:00:01", "-i", srcPath, "-vframes", "1",
		"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", s.MaxSize, s.MaxSize),
		"-f", "image2", dstPath)
	if err := cmd.Run(); err == nil {
		if _, err := os.Stat(dstPath); err == nil {
			return nil
		}
	}
	return s.generatePlaceholderThumb(dstPath)
}

func (s *ThumbnailService) generatePlaceholderThumb(dstPath string) error {
	dst := image.NewRGBA(image.Rect(0, 0, s.MaxSize, s.MaxSize))
	bgCol := color.RGBA{30, 41, 59, 255}
	for y := 0; y < s.MaxSize; y++ {
		for x := 0; x < s.MaxSize; x++ {
			dst.Set(x, y, bgCol)
		}
	}

	cx, cy := s.MaxSize/2, s.MaxSize/2
	halfSize := 35
	for y := cy - halfSize; y <= cy+halfSize; y++ {
		for x := cx - halfSize/2; x <= cx+halfSize; x++ {
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