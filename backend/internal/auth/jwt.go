// Package auth provides JWT signing/verification, OTP management, and the
// Gin AuthMiddleware for the e-commerce API.
package auth

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	OTPTokenTTL     = 10 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour // 7 days
)

// JWTService signs and verifies Ed25519 JWTs.
type JWTService struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// NewJWTService parses PEM-encoded Ed25519 keys and returns a JWTService.
// privateKeyPEM must be a PKCS#8 private key; publicKeyPEM must be a PKIX public key.
func NewJWTService(privateKeyPEM, publicKeyPEM string) (*JWTService, error) {
	privKey, err := parseEd25519PrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("jwt: parse private key: %w", err)
	}
	pubKey, err := parseEd25519PublicKey(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("jwt: parse public key: %w", err)
	}
	return &JWTService{privateKey: privKey, publicKey: pubKey}, nil
}

// AccessClaims are the claims embedded in a standard access token.
type AccessClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// OTPClaims are the claims embedded in a short-lived OTP verification token.
type OTPClaims struct {
	jwt.RegisteredClaims
	// Purpose is either "register_otp" or "password_reset_otp".
	Purpose string `json:"purpose"`
}

// SignAccessToken creates a 15-minute Ed25519-signed JWT for userID.
func (s *JWTService) SignAccessToken(userID uuid.UUID, role string) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
		},
		Role: role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(s.privateKey)
}

// SignOTPToken creates a 10-minute Ed25519-signed JWT binding an email to a
// purpose ("register_otp" or "password_reset_otp").
func (s *JWTService) SignOTPToken(email, purpose string) (string, error) {
	now := time.Now()
	claims := OTPClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(OTPTokenTTL)),
		},
		Purpose: purpose,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(s.privateKey)
}

// VerifyAccessToken parses and validates an access token; returns its claims.
func (s *JWTService) VerifyAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, errors.New("jwt: invalid access token")
	}
	return claims, nil
}

// VerifyOTPToken parses and validates an OTP token; returns its claims.
func (s *JWTService) VerifyOTPToken(tokenStr string) (*OTPClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &OTPClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return s.publicKey, nil
	})
	fmt.Println("Can I see this???/")
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*OTPClaims)
	if !ok || !token.Valid {
		return nil, errors.New("jwt: invalid otp token")
	}
	return claims, nil
}

// normalisePEM trims surrounding quotes and replaces literal \n escape
// sequences with real newlines so that PEM values set as single-line
// environment variables (with or without wrapping quotes) are correctly decoded.
func normalisePEM(s string) string {
	s = strings.Trim(s, `"`)
	return strings.ReplaceAll(s, `\n`, "\n")
}

// parseEd25519PrivateKey decodes a PEM block and returns an ed25519.PrivateKey.
func parseEd25519PrivateKey(pemStr string) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode([]byte(normalisePEM(pemStr)))
	if block == nil {
		return nil, errors.New("jwt: failed to decode private key PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("jwt: parse PKCS8 private key: %w", err)
	}
	edKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("jwt: private key is not Ed25519")
	}
	return edKey, nil
}

// parseEd25519PublicKey decodes a PEM block and returns an ed25519.PublicKey.
func parseEd25519PublicKey(pemStr string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(normalisePEM(pemStr)))
	if block == nil {
		return nil, errors.New("jwt: failed to decode public key PEM block")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("jwt: parse PKIX public key: %w", err)
	}
	edKey, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("jwt: public key is not Ed25519")
	}
	return edKey, nil
}
