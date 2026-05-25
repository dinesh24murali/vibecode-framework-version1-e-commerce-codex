// Package services contains all business logic for the e-commerce API.
// Services depend on repository interfaces and infrastructure clients; they
// never import gin or interact with HTTP directly.
package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/auth"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Sentinel errors returned by AuthService methods.
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
	ErrInvalidOTPToken    = errors.New("invalid or expired OTP token")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrTokenReuse         = errors.New("refresh token reuse detected — all sessions invalidated")
	ErrInvalidToken       = errors.New("invalid or expired refresh token")
)

// OTP purpose constants — must match what is stored in the Redis key and JWT claim.
const (
	OTPPurposeRegister      = "register_otp"
	OTPPurposePasswordReset = "password_reset_otp"
)

// Argon2id parameters (time=1, memory=64MB, threads=4, keyLen=32).
const (
	argon2Time    = 1
	argon2Memory  = 64 * 1024 // 64 MB in KiB
	argon2Threads = 4
	argon2KeyLen  = 32
)

// AuthTokens holds the values returned to the handler after a successful login.
type AuthTokens struct {
	AccessToken  string
	RefreshToken string // raw opaque value (before hashing)
	ExpiresIn    int    // seconds
	User         *dbsqlc.User
}

// AuthService implements all authentication flows.
type AuthService struct {
	users    repositories.UserRepository
	tokens   repositories.RefreshTokenRepository
	jwtSvc   *auth.JWTService
	otpSvc   *auth.OTPService
	fbClient *auth.FacebookClient
}

// NewAuthService constructs an AuthService with all dependencies.
func NewAuthService(
	users repositories.UserRepository,
	tokens repositories.RefreshTokenRepository,
	jwtSvc *auth.JWTService,
	otpSvc *auth.OTPService,
	fbClient *auth.FacebookClient,
) *AuthService {
	return &AuthService{
		users:    users,
		tokens:   tokens,
		jwtSvc:   jwtSvc,
		otpSvc:   otpSvc,
		fbClient: fbClient,
	}
}

// RegisterInit sends (or logs) a registration OTP for the given email.
// Always returns nil — email enumeration prevention.
func (s *AuthService) RegisterInit(ctx context.Context, email string) {
	otp, err := s.otpSvc.GenerateAndStore(ctx, email, OTPPurposeRegister)
	if err != nil {
		slog.Error("register init: store OTP", "error", err)
		return
	}
	// DEV: log OTP to stdout. Replace with SMTP/SES in production.
	slog.Info("register OTP", "email", email, "otp", otp)
}

// RegisterVerifyOTP validates the registration OTP and returns a short-lived otp_token.
func (s *AuthService) RegisterVerifyOTP(ctx context.Context, email, otp string) (string, error) {
	ok, err := s.otpSvc.Verify(ctx, email, OTPPurposeRegister, otp)
	if err != nil {
		return "", fmt.Errorf("register verify otp: %w", err)
	}
	if !ok {
		return "", ErrInvalidOTP
	}
	token, err := s.jwtSvc.SignOTPToken(email, OTPPurposeRegister)
	if err != nil {
		return "", fmt.Errorf("register verify otp: sign token: %w", err)
	}
	return token, nil
}

// RegisterComplete creates the user account and issues auth tokens.
func (s *AuthService) RegisterComplete(ctx context.Context, otpToken, firstName, lastName, password, clientIP, userAgent string) (*AuthTokens, error) {
	claims, err := s.jwtSvc.VerifyOTPToken(otpToken)
	fmt.Println("What calims!!!")
	if err != nil {
		return nil, ErrInvalidOTPToken
	}
	if claims.Purpose != OTPPurposeRegister {
		return nil, ErrInvalidOTPToken
	}
	email := claims.Subject

	// Prevent duplicate registrations.
	existing, _ := s.users.GetByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("register complete: hash password: %w", err)
	}
	fmt.Println("What service!!!")

	user, err := s.users.Create(ctx, dbsqlc.CreateUserParams{
		Email:        email,
		PasswordHash: sql.NullString{String: hash, Valid: true},
		FullName:     firstName + " " + lastName,
		Role:         dbsqlc.UserRoleCustomer,
	})
	if err != nil {
		return nil, fmt.Errorf("register complete: create user: %w", err)
	}
	fmt.Println("Whee is it?")
	fmt.Println("Registraion complete!")

	return s.issueTokens(ctx, user, clientIP, userAgent)
}

// Login verifies email/password and issues auth tokens.
func (s *AuthService) Login(ctx context.Context, email, password, clientIP, userAgent string) (*AuthTokens, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !user.PasswordHash.Valid {
		return nil, ErrInvalidCredentials
	}
	if !verifyPassword(password, user.PasswordHash.String) {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokens(ctx, user, clientIP, userAgent)
}

// FacebookLogin validates a Facebook access token and upserts the user.
func (s *AuthService) FacebookLogin(ctx context.Context, accessToken, clientIP, userAgent string) (*AuthTokens, error) {
	profile, err := s.fbClient.GetProfile(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("facebook login: %w", err)
	}

	user, err := s.users.GetByEmail(ctx, profile.Email)
	if err != nil {
		// Create a new user (no password — Facebook-only account).
		user, err = s.users.Create(ctx, dbsqlc.CreateUserParams{
			Email:    profile.Email,
			FullName: profile.FirstName + " " + profile.LastName,
			Role:     dbsqlc.UserRoleCustomer,
		})
		if err != nil {
			return nil, fmt.Errorf("facebook login: create user: %w", err)
		}
	}
	return s.issueTokens(ctx, user, clientIP, userAgent)
}

// Refresh rotates the refresh token and issues a new access token.
// Reuse of a revoked refresh token invalidates the entire token family.
func (s *AuthService) Refresh(ctx context.Context, rawToken, clientIP, userAgent string) (*AuthTokens, error) {
	hash := hashToken(rawToken)
	stored, err := s.tokens.GetByHash(ctx, hash)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Detect reuse of a revoked token.
	if stored.RevokedAt.Valid {
		// Invalidate the entire family (breach detection).
		_ = s.tokens.RevokeFamily(ctx, stored.TokenFamilyID)
		return nil, ErrTokenReuse
	}

	if time.Now().After(stored.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// Revoke the consumed token before issuing a new one.
	if err := s.tokens.Revoke(ctx, stored.ID); err != nil {
		return nil, fmt.Errorf("refresh: revoke old token: %w", err)
	}

	user, err := s.users.GetByID(ctx, stored.UserID)
	if err != nil {
		return nil, fmt.Errorf("refresh: get user: %w", err)
	}

	// Issue new token pair under the same family.
	return s.issueTokensInFamily(ctx, user, stored.TokenFamilyID, clientIP, userAgent)
}

// Logout revokes the refresh token family so all active sessions for the
// family are invalidated.
func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	hash := hashToken(rawToken)
	stored, err := s.tokens.GetByHash(ctx, hash)
	if err != nil {
		// Token not found — treat as already logged out.
		return nil
	}
	return s.tokens.RevokeFamily(ctx, stored.TokenFamilyID)
}

// PasswordResetInit sends (or logs) a password-reset OTP.
// Always returns nil — email enumeration prevention.
func (s *AuthService) PasswordResetInit(ctx context.Context, email string) {
	otp, err := s.otpSvc.GenerateAndStore(ctx, email, OTPPurposePasswordReset)
	if err != nil {
		slog.Error("password reset init: store OTP", "error", err)
		return
	}
	// DEV: log OTP to stdout. Replace with SMTP/SES in production.
	slog.Info("password reset OTP", "email", email, "otp", otp)
}

// PasswordResetVerifyOTP validates the password-reset OTP and returns an otp_token.
func (s *AuthService) PasswordResetVerifyOTP(ctx context.Context, email, otp string) (string, error) {
	ok, err := s.otpSvc.Verify(ctx, email, OTPPurposePasswordReset, otp)
	if err != nil {
		return "", fmt.Errorf("password reset verify otp: %w", err)
	}
	if !ok {
		return "", ErrInvalidOTP
	}
	token, err := s.jwtSvc.SignOTPToken(email, OTPPurposePasswordReset)
	if err != nil {
		return "", fmt.Errorf("password reset verify otp: sign token: %w", err)
	}
	return token, nil
}

// PasswordResetComplete sets a new password using the otp_token.
func (s *AuthService) PasswordResetComplete(ctx context.Context, otpToken, newPassword string) error {
	claims, err := s.jwtSvc.VerifyOTPToken(otpToken)
	if err != nil {
		return ErrInvalidOTPToken
	}
	if claims.Purpose != OTPPurposePasswordReset {
		return ErrInvalidOTPToken
	}
	email := claims.Subject

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return ErrInvalidOTPToken // do not reveal whether the email exists
	}

	hash, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("password reset complete: hash: %w", err)
	}
	if err := s.users.UpdatePasswordHash(ctx, user.ID, hash); err != nil {
		return fmt.Errorf("password reset complete: update hash: %w", err)
	}

	slog.Info("password reset complete: password updated", "user_id", user.ID)
	return nil
}

// issueTokens creates a new token family and issues access + refresh tokens.
func (s *AuthService) issueTokens(ctx context.Context, user *dbsqlc.User, clientIP, userAgent string) (*AuthTokens, error) {
	familyID := uuid.New().String()
	return s.issueTokensInFamily(ctx, user, familyID, clientIP, userAgent)
}

// issueTokensInFamily issues a new access token and refresh token within
// the given family.
func (s *AuthService) issueTokensInFamily(ctx context.Context, user *dbsqlc.User, familyID, clientIP, userAgent string) (*AuthTokens, error) {
	accessToken, err := s.jwtSvc.SignAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("issue tokens: sign access: %w", err)
	}

	rawRefresh := uuid.New().String()
	refreshHash := hashToken(rawRefresh)

	_, err = s.tokens.Create(ctx, dbsqlc.CreateRefreshTokenParams{
		UserID:        user.ID,
		TokenFamilyID: familyID,
		TokenHash:     refreshHash,
		UserAgent:     sql.NullString{String: userAgent, Valid: userAgent != ""},
		IpAddress:     repositories.ParseInet(clientIP),
		ExpiresAt:     time.Now().Add(auth.RefreshTokenTTL),
	})
	if err != nil {
		return nil, fmt.Errorf("issue tokens: store refresh token: %w", err)
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int(auth.AccessTokenTTL.Seconds()),
		User:         user,
	}, nil
}

// hashPassword creates an Argon2id hash of the plaintext password.
// Stored format: hex(salt):hex(hash)
func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("hash password: generate salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(hash), nil
}

// verifyPassword checks a plaintext password against a stored Argon2id hash.
func verifyPassword(password, stored string) bool {
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return false
	}
	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}
	expected, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	actual := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	return constantTimeEqual(actual, expected)
}

// constantTimeEqual compares two byte slices in constant time.
func constantTimeEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := range a {
		v |= a[i] ^ b[i]
	}
	return v == 0
}

// hashToken returns the SHA-256 hex digest of a raw token string.
func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
