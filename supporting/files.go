package supporting

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/jackillll/uqpay-sdk-go/common"
)

// FilesClient handles file operations
type FilesClient struct {
	client *common.APIClient
}

// UploadFileParams represents file upload parameters
type UploadFileParams struct {
	File     io.Reader
	FileName string
	Notes    string
}

// UploadFileResponse represents file upload response
type UploadFileResponse struct {
	CreateTime string `json:"create_time"`
	FileID     string `json:"file_id"`
	FileName   string `json:"file_name"`
	FileType   string `json:"file_type"`
	Size       int    `json:"size"`
	Notes      string `json:"notes"`
}

// DownloadLinksRequest represents download links request
type DownloadLinksRequest struct {
	FileIDs []string `json:"file_ids"` // required
}

// FileDownloadInfo represents file download information
type FileDownloadInfo struct {
	FileID   string `json:"file_id"`
	FileType string `json:"file_type"`
	FileName string `json:"file_name"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
}

// DownloadLinksResponse represents download links response
type DownloadLinksResponse struct {
	Files       []FileDownloadInfo `json:"files"`
	AbsentFiles []string           `json:"absent_files"`
}

// Upload uploads a file to UQPAY
// POST /v1/files/upload
// Maximum file size: 20MB
// Supported types: jpeg, png, jpg, doc, docx, pdf
func (c *FilesClient) Upload(ctx context.Context, params *UploadFileParams) (*UploadFileResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create form file
	part, err := writer.CreateFormFile("file", params.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content
	if _, err := io.Copy(part, params.File); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add notes if provided
	if params.Notes != "" {
		if err := writer.WriteField("notes", params.Notes); err != nil {
			return nil, fmt.Errorf("failed to write notes field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// TODO: Implement multipart form-data POST in APIClient
	// For now, this is a placeholder implementation
	var resp UploadFileResponse
	path := "/v1/files/upload"
	if params.Notes != "" {
		path += fmt.Sprintf("?notes=%s", params.Notes)
	}

	// Note: This requires special handling in APIClient for multipart/form-data
	// The actual implementation needs to be added to common.APIClient
	_ = path // Placeholder to avoid unused variable error

	return &resp, fmt.Errorf("file upload not yet implemented - requires multipart/form-data support")
}

// GetDownloadLinks retrieves download links for specified file IDs
// POST /v1/files/download_links
func (c *FilesClient) GetDownloadLinks(ctx context.Context, req *DownloadLinksRequest) (*DownloadLinksResponse, error) {
	var resp DownloadLinksResponse
	if err := c.client.Post(ctx, "/v1/files/download_links", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to get download links: %w", err)
	}
	return &resp, nil
}
