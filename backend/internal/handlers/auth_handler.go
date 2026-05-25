package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/auth"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// AuthHandler wires HTTP to AuthService. Each method corresponds to one
// endpoint defined in docs/02_outputs/04_api_spec.yaml.
type AuthHandler struct {
	svc *services.AuthService
}

// NewAuthHandler returns an AuthHandler backed by the given service.
func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// ── Request / Response types ─────────────────────────────────────────────────

type registerInitRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type verifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type registerCompleteRequest struct {
	OTPToken  string `json:"otp_token" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type facebookLoginRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type passwordResetRequestBody struct {
	Email string `json:"email" binding:"required,email"`
}

type passwordResetCompleteRequest struct {
	OTPToken           string `json:"otp_token" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
}

// ── Handlers ─────────────────────────────────────────────────────────────────

// RegisterInit godoc POST /api/v1/auth/register/init
func (h *AuthHandler) RegisterInit(c *gin.Context) {
	var req registerInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.svc.RegisterInit(c.Request.Context(), req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "If this email is not already registered, an OTP has been sent."})
}

// RegisterVerifyOTP godoc POST /api/v1/auth/register/verify-otp
func (h *AuthHandler) RegisterVerifyOTP(c *gin.Context) {
	var req verifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	otpToken, err := h.svc.RegisterVerifyOTP(c.Request.Context(), req.Email, req.OTP)
	if err != nil {
		if errors.Is(err, services.ErrInvalidOTP) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid or expired OTP"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"otp_token":  otpToken,
		"expires_in": int(auth.OTPTokenTTL.Seconds()),
	})
}

// RegisterComplete godoc POST /api/v1/auth/register/complete
func (h *AuthHandler) RegisterComplete(c *gin.Context) {
	fmt.Println("What handler!!!")
	var req registerCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("What handler!!!")
	fmt.Println(c.ClientIP())
	fmt.Println(c.GetHeader("User-Agent"))
	tokens, err := h.svc.RegisterComplete(
		c.Request.Context(),
		req.OTPToken, req.FirstName, req.LastName, req.Password,
		c.ClientIP(), c.GetHeader("User-Agent"),
	)
	fmt.Printf("%s", err.Error())

	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidOTPToken):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid or expired OTP token"})
		case errors.Is(err, services.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	h.setRefreshCookie(c, tokens.RefreshToken)
	c.JSON(http.StatusCreated, buildAuthResponse(tokens))
}

// Login godoc POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := h.svc.Login(
		c.Request.Context(),
		req.Email, req.Password,
		c.ClientIP(), c.GetHeader("User-Agent"),
	)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	h.setRefreshCookie(c, tokens.RefreshToken)
	c.JSON(http.StatusOK, buildAuthResponse(tokens))
}

// FacebookLogin godoc POST /api/v1/auth/login/facebook
func (h *AuthHandler) FacebookLogin(c *gin.Context) {
	var req facebookLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := h.svc.FacebookLogin(
		c.Request.Context(),
		req.AccessToken,
		c.ClientIP(), c.GetHeader("User-Agent"),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Facebook authentication failed"})
		return
	}
	h.setRefreshCookie(c, tokens.RefreshToken)
	c.JSON(http.StatusOK, buildAuthResponse(tokens))
}

// RefreshToken godoc POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := h.svc.Refresh(
		c.Request.Context(),
		req.RefreshToken,
		c.ClientIP(), c.GetHeader("User-Agent"),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	h.setRefreshCookie(c, tokens.RefreshToken)
	c.JSON(http.StatusOK, buildAuthResponse(tokens))
}

// Logout godoc POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req logoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_ = h.svc.Logout(c.Request.Context(), req.RefreshToken)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// PasswordResetRequest godoc POST /api/v1/auth/password-reset/request
func (h *AuthHandler) PasswordResetRequest(c *gin.Context) {
	var req passwordResetRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.svc.PasswordResetInit(c.Request.Context(), req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "If this email is registered, a password-reset OTP has been sent."})
}

// PasswordResetVerifyOTP godoc POST /api/v1/auth/password-reset/verify-otp
func (h *AuthHandler) PasswordResetVerifyOTP(c *gin.Context) {
	var req verifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	otpToken, err := h.svc.PasswordResetVerifyOTP(c.Request.Context(), req.Email, req.OTP)
	if err != nil {
		if errors.Is(err, services.ErrInvalidOTP) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid or expired OTP"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"otp_token":  otpToken,
		"expires_in": int(auth.OTPTokenTTL.Seconds()),
	})
}

// PasswordResetComplete godoc POST /api/v1/auth/password-reset/complete
func (h *AuthHandler) PasswordResetComplete(c *gin.Context) {
	var req passwordResetCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.NewPassword != req.ConfirmNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}
	if err := h.svc.PasswordResetComplete(c.Request.Context(), req.OTPToken, req.NewPassword); err != nil {
		if errors.Is(err, services.ErrInvalidOTPToken) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid or expired OTP token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Your password has been reset. Please log in with your new password."})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// setRefreshCookie sets the HttpOnly refresh token cookie.
func (h *AuthHandler) setRefreshCookie(c *gin.Context, rawToken string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token",
		rawToken,
		int(auth.RefreshTokenTTL.Seconds()),
		"/",
		"",   // domain — empty uses request host
		true, // secure
		true, // httpOnly
	)
}

// buildAuthResponse constructs the AuthTokenResponse JSON shape.
// full_name is split on the first space to produce first_name / last_name
// as required by the API spec UserProfile schema.
func buildAuthResponse(tokens *services.AuthTokens) gin.H {
	firstName, lastName := splitFullName(tokens.User.FullName)
	return gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"token_type":    "bearer",
		"expires_in":    tokens.ExpiresIn,
		"user": gin.H{
			"id":         tokens.User.ID,
			"email":      tokens.User.Email,
			"first_name": firstName,
			"last_name":  lastName,
			"role":       tokens.User.Role,
			"created_at": tokens.User.CreatedAt,
			"updated_at": tokens.User.UpdatedAt,
		},
	}
}

// splitFullName splits "Alice Smith" → ("Alice", "Smith").
// If there is no space, the entire string is the first name.
func splitFullName(fullName string) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(fullName), " ", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return fullName, ""
}
