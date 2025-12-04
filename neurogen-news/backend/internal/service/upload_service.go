package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UploadService interface {
	UploadImage(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader, uploadType string) (*UploadResult, error)
	DeleteFile(ctx context.Context, userID uuid.UUID, fileURL string) error
}

type UploadResult struct {
	URL       string `json:"url"`
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mimeType"`
}

type uploadService struct {
	logger *zap.Logger
	// TODO: Add S3 client for production
}

func NewUploadService(logger *zap.Logger) UploadService {
	return &uploadService{
		logger: logger,
	}
}

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var maxImageSize int64 = 10 * 1024 * 1024 // 10MB

func (s *uploadService) UploadImage(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader, uploadType string) (*UploadResult, error) {
	// Validate file size
	if file.Size > maxImageSize {
		return nil, &AppError{Code: "FILE_TOO_LARGE", Message: "File size exceeds 10MB limit"}
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return nil, &AppError{Code: "INVALID_FILE_TYPE", Message: "Only JPEG, PNG, GIF and WebP images are allowed"}
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Read file content
	content, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		}
	}

	filename := fmt.Sprintf("%s/%s/%s%s",
		uploadType,
		time.Now().Format("2006/01"),
		uuid.New().String(),
		strings.ToLower(ext),
	)

	// TODO: Upload to S3 in production
	// For now, just return a placeholder URL
	_ = content

	url := "/uploads/" + filename

	return &UploadResult{
		URL:      url,
		Filename: filename,
		Size:     file.Size,
		MimeType: contentType,
	}, nil
}

func (s *uploadService) DeleteFile(ctx context.Context, userID uuid.UUID, fileURL string) error {
	// TODO: Implement file deletion from S3
	return nil
}

