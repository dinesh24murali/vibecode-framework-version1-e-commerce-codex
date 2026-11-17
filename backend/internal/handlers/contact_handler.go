package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// ContactHandler serves the public contact form endpoint.
type ContactHandler struct {
	svc *services.ContactService
}

// NewContactHandler returns a ContactHandler.
func NewContactHandler(svc *services.ContactService) *ContactHandler {
	return &ContactHandler{svc: svc}
}

type contactRequest struct {
	Name    string  `json:"name" binding:"required"`
	Email   string  `json:"email" binding:"required"`
	Phone   *string `json:"phone"`
	Subject string  `json:"subject" binding:"required"`
	Message string  `json:"message" binding:"required"`
}

// Submit godoc POST /api/v1/contact
func (h *ContactHandler) Submit(c *gin.Context) {
	var req contactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "name, email, subject and message are required")
		return
	}
	err := h.svc.Submit(c.Request.Context(), services.ContactInput{
		Name: req.Name, Email: req.Email, Phone: req.Phone, Subject: req.Subject, Message: req.Message,
	})
	switch {
	case err == nil:
		c.JSON(http.StatusOK, gin.H{"message": "Thank you. We have received your enquiry and will get back to you soon."})
	case errors.Is(err, services.ErrContactInvalid):
		apiError(c, http.StatusUnprocessableEntity, err.Error())
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}
