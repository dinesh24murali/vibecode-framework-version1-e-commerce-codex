package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// FacebookProfile holds the fields returned by the Graph API /me endpoint.
type FacebookProfile struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// FacebookClient validates a Facebook access token by calling the Graph API.
// MVP: calls /me directly (does not call /debug_token for app verification).
// Known limitation: a token stolen from another Facebook app would be accepted.
type FacebookClient struct {
	httpClient *http.Client
}

// NewFacebookClient returns a FacebookClient with a sensible timeout.
func NewFacebookClient() *FacebookClient {
	return &FacebookClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetProfile exchanges a Facebook access token for the user's profile.
func (c *FacebookClient) GetProfile(ctx context.Context, accessToken string) (*FacebookProfile, error) {
	url := fmt.Sprintf(
		"https://graph.facebook.com/me?fields=id,email,first_name,last_name&access_token=%s",
		accessToken,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("facebook: build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("facebook: graph api call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("facebook: invalid or expired access token")
	}

	var profile FacebookProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("facebook: decode response: %w", err)
	}
	if profile.Email == "" {
		return nil, errors.New("facebook: email not available in token scope")
	}
	return &profile, nil
}
