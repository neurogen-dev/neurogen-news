package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UploadType represents the type of upload
type UploadType string

const (
	UploadTypeArticle UploadType = "article"
	UploadTypeAvatar  UploadType = "avatar"
	UploadTypeCover   UploadType = "cover"
)

// UploadResult represents the result of an upload
type UploadResult struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
}

// S3Config holds S3 configuration
type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string // For S3-compatible services (MinIO, DigitalOcean Spaces, etc.)
	CDNBaseURL      string // CDN URL for public access
	MaxFileSize     int64  // Maximum file size in bytes
}

// S3Client wraps AWS S3 client
type S3Client struct {
	client     *s3.Client
	bucket     string
	cdnBaseURL string
	maxSize    int64
	logger     *zap.Logger
}

// Allowed image MIME types
var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// NewS3Client creates a new S3 client
func NewS3Client(cfg S3Config, logger *zap.Logger) (*S3Client, error) {
	// Create AWS config
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client options
	var clientOpts []func(*s3.Options)

	// Use custom endpoint if provided (for S3-compatible services)
	if cfg.Endpoint != "" {
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		})
	}

	client := s3.NewFromConfig(awsCfg, clientOpts...)

	// Set default max size (10MB)
	maxSize := cfg.MaxFileSize
	if maxSize <= 0 {
		maxSize = 10 * 1024 * 1024
	}

	return &S3Client{
		client:     client,
		bucket:     cfg.Bucket,
		cdnBaseURL: cfg.CDNBaseURL,
		maxSize:    maxSize,
		logger:     logger,
	}, nil
}

// UploadImage uploads an image to S3
func (c *S3Client) UploadImage(ctx context.Context, file *multipart.FileHeader, uploadType UploadType) (*UploadResult, error) {
	// Validate file size
	if file.Size > c.maxSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", c.maxSize)
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file content
	content, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Detect content type
	contentType := http.DetectContentType(content)

	// Validate image type
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		return nil, fmt.Errorf("unsupported image type: %s", contentType)
	}

	// Generate unique filename
	filename := generateFilename(uploadType, ext)

	// Upload to S3
	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(c.bucket),
		Key:          aws.String(filename),
		Body:         bytes.NewReader(content),
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("public, max-age=31536000"), // 1 year cache
	})
	if err != nil {
		c.logger.Error("Failed to upload to S3", zap.Error(err))
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Build URL
	url := c.buildURL(filename)

	return &UploadResult{
		URL:      url,
		Key:      filename,
		Filename: filepath.Base(filename),
		Size:     file.Size,
		MimeType: contentType,
	}, nil
}

// UploadImageFromBytes uploads image from byte slice
func (c *S3Client) UploadImageFromBytes(ctx context.Context, content []byte, originalFilename string, uploadType UploadType) (*UploadResult, error) {
	// Validate file size
	if int64(len(content)) > c.maxSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", c.maxSize)
	}

	// Detect content type
	contentType := http.DetectContentType(content)

	// Validate image type
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		return nil, fmt.Errorf("unsupported image type: %s", contentType)
	}

	// Generate unique filename
	filename := generateFilename(uploadType, ext)

	// Upload to S3
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(c.bucket),
		Key:          aws.String(filename),
		Body:         bytes.NewReader(content),
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("public, max-age=31536000"),
	})
	if err != nil {
		c.logger.Error("Failed to upload to S3", zap.Error(err))
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	url := c.buildURL(filename)

	return &UploadResult{
		URL:      url,
		Key:      filename,
		Filename: filepath.Base(filename),
		Size:     int64(len(content)),
		MimeType: contentType,
	}, nil
}

// DeleteFile deletes a file from S3
func (c *S3Client) DeleteFile(ctx context.Context, key string) error {
	// Extract key from URL if full URL is provided
	if strings.HasPrefix(key, "http") {
		key = extractKeyFromURL(key, c.cdnBaseURL, c.bucket)
	}

	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		c.logger.Error("Failed to delete from S3", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// DeleteFiles deletes multiple files from S3
func (c *S3Client) DeleteFiles(ctx context.Context, keys []string) error {
	for _, key := range keys {
		if err := c.DeleteFile(ctx, key); err != nil {
			c.logger.Error("Failed to delete file", zap.String("key", key), zap.Error(err))
			// Continue with other files
		}
	}
	return nil
}

// GetPresignedURL generates a presigned URL for direct upload
func (c *S3Client) GetPresignedURL(ctx context.Context, filename string, contentType string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.client)

	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(filename),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

// buildURL constructs the public URL for a file
func (c *S3Client) buildURL(key string) string {
	if c.cdnBaseURL != "" {
		return strings.TrimSuffix(c.cdnBaseURL, "/") + "/" + key
	}
	// Default S3 URL format
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.bucket, key)
}

// generateFilename creates a unique filename with path structure
func generateFilename(uploadType UploadType, ext string) string {
	now := time.Now()
	id := uuid.New().String()

	// Structure: type/year/month/uuid.ext
	return fmt.Sprintf("%s/%d/%02d/%s%s",
		uploadType,
		now.Year(),
		now.Month(),
		id,
		ext,
	)
}

// extractKeyFromURL extracts the S3 key from a full URL
func extractKeyFromURL(url, cdnBaseURL, bucket string) string {
	// Try CDN URL first
	if cdnBaseURL != "" && strings.HasPrefix(url, cdnBaseURL) {
		return strings.TrimPrefix(url, strings.TrimSuffix(cdnBaseURL, "/")+"/")
	}

	// Try S3 URL
	s3Prefix := fmt.Sprintf("https://%s.s3.amazonaws.com/", bucket)
	if strings.HasPrefix(url, s3Prefix) {
		return strings.TrimPrefix(url, s3Prefix)
	}

	// Return as-is if no match
	return url
}

// MockS3Client provides a mock implementation for development
type MockS3Client struct {
	logger *zap.Logger
}

// NewMockS3Client creates a mock S3 client for development
func NewMockS3Client(logger *zap.Logger) *MockS3Client {
	return &MockS3Client{logger: logger}
}

// UploadImage mocks image upload
func (c *MockS3Client) UploadImage(ctx context.Context, file *multipart.FileHeader, uploadType UploadType) (*UploadResult, error) {
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	ext := ".jpg"
	if e, ok := allowedImageTypes[contentType]; ok {
		ext = e
	}

	filename := generateFilename(uploadType, ext)
	url := "/uploads/" + filename

	c.logger.Info("Mock upload",
		zap.String("filename", filename),
		zap.Int64("size", file.Size),
		zap.String("contentType", contentType))

	return &UploadResult{
		URL:      url,
		Key:      filename,
		Filename: filepath.Base(filename),
		Size:     file.Size,
		MimeType: contentType,
	}, nil
}

// DeleteFile mocks file deletion
func (c *MockS3Client) DeleteFile(ctx context.Context, key string) error {
	c.logger.Info("Mock delete", zap.String("key", key))
	return nil
}

// Storage interface for dependency injection
type Storage interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, uploadType UploadType) (*UploadResult, error)
	DeleteFile(ctx context.Context, key string) error
}


