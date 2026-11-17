package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PaymentGateway creates payment orders at the payment provider.
type PaymentGateway interface {
	// CreateOrder registers an order for amountPaise and returns the provider order id.
	CreateOrder(ctx context.Context, amountPaise int64, receipt string) (string, error)
	// KeyID is the public key the frontend passes to the provider checkout.
	KeyID() string
}

// RazorpayGateway calls the Razorpay Orders API v1.
type RazorpayGateway struct {
	client  *http.Client
	keyID   string
	secret  string
	baseURL string
}

// NewRazorpayGateway returns a gateway using HTTP basic auth with the key pair.
func NewRazorpayGateway(keyID, secret string) *RazorpayGateway {
	return &RazorpayGateway{client: &http.Client{Timeout: 15 * time.Second}, keyID: keyID, secret: secret, baseURL: "https://api.razorpay.com"}
}

// KeyID returns the public key id.
func (g *RazorpayGateway) KeyID() string { return g.keyID }

// CreateOrder implements PaymentGateway.
func (g *RazorpayGateway) CreateOrder(ctx context.Context, amountPaise int64, receipt string) (string, error) {
	body, _ := json.Marshal(map[string]any{"amount": amountPaise, "currency": "INR", "receipt": receipt})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/v1/orders", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.keyID, g.secret)
	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("razorpay: create order: status %d", resp.StatusCode)
	}
	var out struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.ID == "" {
		return "", fmt.Errorf("razorpay: unexpected create order response")
	}
	return out.ID, nil
}
