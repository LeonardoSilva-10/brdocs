package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// serproToken holds an OAuth2 access token for Serpro APIs.
type serproToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	expiresAt   time.Time
}

func (t *serproToken) isExpired() bool {
	return time.Now().After(t.expiresAt.Add(-2 * time.Minute))
}

// serproClient handles Serpro/ConectaGov API communication.
type serproClient struct {
	parent     *Client
	token      *serproToken
	tokenMutex sync.RWMutex
}

// newSerproClient creates a new Serpro client.
func newSerproClient(parent *Client) *serproClient {
	return &serproClient{parent: parent}
}

// getToken obtains or refreshes the OAuth2 access token.
func (s *serproClient) getToken(ctx context.Context) (string, error) {
	s.tokenMutex.RLock()
	if s.token != nil && !s.token.isExpired() {
		token := s.token.AccessToken
		s.tokenMutex.RUnlock()
		return token, nil
	}
	s.tokenMutex.RUnlock()

	s.tokenMutex.Lock()
	defer s.tokenMutex.Unlock()

	// Double-check after acquiring write lock
	if s.token != nil && !s.token.isExpired() {
		return s.token.AccessToken, nil
	}

	if s.parent.serproConsumerKey == "" || s.parent.serproConsumerSecret == "" {
		return "", fmt.Errorf("serpro: consumer key and secret are required")
	}

	auth := base64.StdEncoding.EncodeToString([]byte(
		s.parent.serproConsumerKey + ":" + s.parent.serproConsumerSecret,
	))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		s.parent.baseURLs["serpro"]+"/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("serpro: failed to create token request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.parent.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("serpro: token request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("serpro: token request returned status %d", resp.StatusCode)
	}

	var token serproToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("serpro: failed to decode token response: %w", err)
	}

	token.expiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	s.token = &token

	return token.AccessToken, nil
}

// QueryCNPJ queries CNPJ data from Serpro/ConectaGov.
func (s *serproClient) QueryCNPJ(ctx context.Context, cnpj string) (*CNPJResult, error) {
	cnpj = sanitizeDocument(cnpj)
	if len(cnpj) != 14 {
		return nil, fmt.Errorf("serpro: invalid CNPJ length, expected 14 digits")
	}

	cacheKey := "serpro:" + cnpj
	if cached := s.parent.cacheGet(cacheKey); cached != nil {
		return cached.(*CNPJResult), nil
	}

	token, err := s.getToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		s.parent.baseURLs["serpro"]+"/consulta-cnpj-df/v2/basica/"+cnpj,
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("serpro: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	var result CNPJResult
	if err := s.parent.doRequest(ctx, req, &result); err != nil {
		return nil, fmt.Errorf("serpro: %w", err)
	}

	s.parent.cacheSet(cacheKey, &result)
	return &result, nil
}

// QueryCNPJ is a convenience method on Client.
func (c *Client) QueryCNPJ(ctx context.Context, cnpj string) (*CNPJResult, error) {
	sc := c.serpro
	if sc == nil {
		sc = newSerproClient(c)
	}
	return sc.QueryCNPJ(ctx, cnpj)
}
