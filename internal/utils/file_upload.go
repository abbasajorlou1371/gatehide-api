package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gatehide/gatehide-api/config"
)

// FileUploadResult represents the result of a file upload
type FileUploadResult struct {
	FileName    string
	FilePath    string
	FileSize    int64
	ContentType string
	PublicURL   string
}

// FileUploader handles file upload operations
type FileUploader struct {
	config *config.FileStorageConfig
}

// NewFileUploader creates a new file uploader
func NewFileUploader(cfg *config.FileStorageConfig) *FileUploader {
	return &FileUploader{
		config: cfg,
	}
}

// UploadFile uploads a file and returns the result
func (fu *FileUploader) UploadFile(file *multipart.FileHeader, subfolder string) (*FileUploadResult, error) {
	// Validate file size
	if file.Size > fu.config.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", fu.config.MaxFileSize)
	}

	// Validate file type
	if !fu.isAllowedFileType(file.Filename) {
		return nil, fmt.Errorf("file type not allowed. Allowed types: %v", fu.config.AllowedTypes)
	}

	// Create upload directory if it doesn't exist
	uploadDir := filepath.Join(fu.config.UploadPath, subfolder)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(file.Filename, ext)
	uniqueName := fmt.Sprintf("%s_%d%s", baseName, time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, uniqueName)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Generate public URL
	publicURL := fmt.Sprintf("%s/uploads/%s/%s", fu.config.PublicURL, subfolder, uniqueName)

	return &FileUploadResult{
		FileName:    uniqueName,
		FilePath:    filePath,
		FileSize:    file.Size,
		ContentType: file.Header.Get("Content-Type"),
		PublicURL:   publicURL,
	}, nil
}

// DeleteFile deletes a file from the filesystem
func (fu *FileUploader) DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// isAllowedFileType checks if the file type is allowed
func (fu *FileUploader) isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range fu.config.AllowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// GetFileInfo returns information about a file
func (fu *FileUploader) GetFileInfo(filePath string) (*FileUploadResult, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Extract subfolder and filename from path
	relPath, err := filepath.Rel(fu.config.UploadPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %w", err)
	}

	parts := strings.Split(relPath, string(filepath.Separator))
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid file path structure")
	}

	subfolder := parts[0]
	filename := parts[len(parts)-1]

	// Generate public URL
	publicURL := fmt.Sprintf("%s/uploads/%s/%s", fu.config.PublicURL, subfolder, filename)

	return &FileUploadResult{
		FileName:    filename,
		FilePath:    filePath,
		FileSize:    fileInfo.Size(),
		ContentType: "", // We don't store content type in file info
		PublicURL:   publicURL,
	}, nil
}
