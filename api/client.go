// Package api provides HTTP clients for querying CPF and CNPJ data from
// external APIs, including official government services and third-party providers.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode"
)

// cacheEntry stores a cached query result and its expiry time.
type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

// Client is the main API client for CPF/CNPJ queries.
type Client struct {
	httpClient           *http.Client
	serpro               *serproClient
	serproConsumerKey    string
	serproConsumerSecret string
	cpfCnpjAPIKey        string
	cpfBrasilAPIKey      string
	baseURLs             map[string]string

	cacheMutex sync.RWMutex
	cache      map[string]cacheEntry
	cacheTTL   time.Duration
}

// Option configures a Client.
type Option func(*Client)

// NewClient creates a new API client with the given options.
//
// Result caching is disabled by default; enable it with WithCacheTTL.
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
	c.serpro = newSerproClient(c)
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

// WithCacheTTL enables in-memory caching of successful query results.
//
// Results are cached per provider and document for the given duration. A TTL
// of zero or less disables caching (the default). Only successful responses
// are cached; HTTP errors and parse failures are never cached.
func WithCacheTTL(ttl time.Duration) Option {
	return func(c *Client) {
		c.cacheTTL = ttl
	}
}

// sanitizeDocument removes all non-digit characters from a document string.
func sanitizeDocument(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// cacheGet returns the cached value for key, or nil when the cache is
// disabled, the key is missing, or the entry has expired.
func (c *Client) cacheGet(key string) interface{} {
	if c.cacheTTL <= 0 {
		return nil
	}
	c.cacheMutex.RLock()
	entry, ok := c.cache[key]
	c.cacheMutex.RUnlock()
	if !ok {
		return nil
	}
	if time.Now().After(entry.expiresAt) {
		c.cacheMutex.Lock()
		delete(c.cache, key)
		c.cacheMutex.Unlock()
		return nil
	}
	return entry.value
}

// cacheSet stores value in the cache under key. It is a no-op when caching
// is disabled.
func (c *Client) cacheSet(key string, value interface{}) {
	if c.cacheTTL <= 0 {
		return
	}
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()
	if c.cache == nil {
		c.cache = make(map[string]cacheEntry)
	}
	c.cache[key] = cacheEntry{value: value, expiresAt: time.Now().Add(c.cacheTTL)}
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
