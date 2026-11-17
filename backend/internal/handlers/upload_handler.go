package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// UploadHandler serves the admin image upload endpoints.
type UploadHandler struct {
	storage services.StorageService
}

// NewUploadHandler returns an UploadHandler.
func NewUploadHandler(s services.StorageService) *UploadHandler {
	return &UploadHandler{storage: s}
}

type presignRequest struct {
	Filename      string `json:"filename" binding:"required"`
	ContentType   string `json:"content_type" binding:"required"`
	FileSizeBytes int64  `json:"file_size_bytes" binding:"required"`
}

// Presign godoc POST /api/v1/admin/upload/image/presign
func (h *UploadHandler) Presign(c *gin.Context) {
	var req presignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, err.Error())
		return
	}
	up, err := h.storage.PresignUpload(c.Request.Context(), req.ContentType, req.FileSizeBytes)
	switch {
	case errors.Is(err, services.ErrInvalidImage):
		apiError(c, http.StatusBadRequest, "content_type must be image/jpeg, image/png or image/webp and file_size_bytes 1..5242880")
	case errors.Is(err, services.ErrStorageNotConfigured):
		apiError(c, http.StatusServiceUnavailable, err.Error())
	case err != nil:
		apiError(c, http.StatusInternalServerError, "internal error")
	default:
		c.JSON(http.StatusOK, gin.H{
			"presigned_url": up.PresignedURL, "public_url": up.PublicURL, "expires_at": up.ExpiresAt,
		})
	}
}
