// Package api provides HTTP clients for querying CPF and CNPJ data from
// external APIs, including official government services and third-party providers.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is the main API client for CPF/CNPJ queries.
type Client struct {
	httpClient        *http.Client
	serproConsumerKey string
	serproConsumerSecret string
	cpfCnpjAPIKey     string
	cpfBrasilAPIKey   string
	baseURLs          map[string]string
}

// Option configures a Client.
type Option func(*Client)

// NewClient creates a new API client with the given options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURLs: map[string]string{
			"serpro":         "https://gateway.apiserpro.gov.br",
			"receitafederal": "https://servicos.receita.fazenda.gov.br",
			"cpfcnpj":        "https://api.cpfcnpj.com.br",
			"cpfbrasil":      "https://api.cpf-brasil.org",
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// WithSerproCredentials sets Serpro/ConectaGov OAuth2 credentials.
func WithSerproCredentials(consumerKey, consumerSecret string) Option {
	return func(c *Client) {
		c.serproConsumerKey = consumerKey
		c.serproConsumerSecret = consumerSecret
	}
}

// WithCPFCNPJAPIKey sets the cpfcnpj.com.br API key.
func WithCPFCNPJAPIKey(key string) Option {
	return func(c *Client) {
		c.cpfCnpjAPIKey = key
	}
}

// WithCPFBrasilAPIKey sets the cpf-brasil.org API key.
func WithCPFBrasilAPIKey(key string) Option {
	return func(c *Client) {
		c.cpfBrasilAPIKey = key
	}
}

// WithBaseURL overrides a base URL for a provider.
func WithBaseURL(provider, url string) Option {
	return func(c *Client) {
		c.baseURLs[provider] = url
	}
}

// doRequest performs an HTTP request and decodes the JSON response.
func (c *Client) doRequest(ctx context.Context, req *http.Request, result interface{}) error {
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("api: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api: unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("api: failed to decode response: %w", err)
	}

	return nil
}
