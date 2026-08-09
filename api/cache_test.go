package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// serproCountingServer is a fake Serpro endpoint that counts incoming requests.
type serproCountingServer struct {
	mu    sync.Mutex
	token int
	cnpj  int
	other int
}

func (s *serproCountingServer) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		switch {
		case r.URL.Path == "/token":
			s.token++
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serproToken{AccessToken: "tok", TokenType: "Bearer", ExpiresIn: 3600})
		case strings.Contains(r.URL.Path, "/consulta-cnpj-df/"):
			s.cnpj++
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(CNPJResult{Number: "34238864000168", LegalName: "SERPRO TESTE"})
		default:
			s.other++
			http.NotFound(w, r)
		}
	})
}

func (s *serproCountingServer) counts() (token, cnpj, other int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.token, s.cnpj, s.other
}

func clientForSerpro(server *httptest.Server, opts ...Option) *Client {
	base := []Option{WithSerproCredentials("key", "secret"), WithBaseURL("serpro", server.URL)}
	return NewClient(append(base, opts...)...)
}

func TestClient_CacheDisabledByDefault(t *testing.T) {
	srv := &serproCountingServer{}
	server := httptest.NewServer(srv.handler())
	defer server.Close()

	client := clientForSerpro(server)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		_, err := client.QueryCNPJ(ctx, "34.238.864/0001-68")
		require.NoError(t, err)
	}

	token, cnpj, other := srv.counts()
	assert.Equal(t, 1, token, "OAuth token should be reused across calls")
	assert.Equal(t, 2, cnpj, "without cache every query should hit the server")
	assert.Equal(t, 0, other)
}

func TestClient_CacheEnabled(t *testing.T) {
	srv := &serproCountingServer{}
	server := httptest.NewServer(srv.handler())
	defer server.Close()

	client := clientForSerpro(server, WithCacheTTL(time.Minute))
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		_, err := client.QueryCNPJ(ctx, "34.238.864/0001-68")
		require.NoError(t, err)
	}

	token, cnpj, other := srv.counts()
	assert.Equal(t, 1, token)
	assert.Equal(t, 1, cnpj, "the second query should be served from cache")
	assert.Equal(t, 0, other)
}

func TestClient_CacheExpires(t *testing.T) {
	srv := &serproCountingServer{}
	server := httptest.NewServer(srv.handler())
	defer server.Close()

	client := clientForSerpro(server, WithCacheTTL(120*time.Millisecond))
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		_, err := client.QueryCNPJ(ctx, "34.238.864/0001-68")
		require.NoError(t, err)
	}
	time.Sleep(250 * time.Millisecond)
	_, err := client.QueryCNPJ(ctx, "34.238.864/0001-68")
	require.NoError(t, err)

	token, cnpj, other := srv.counts()
	assert.Equal(t, 1, token)
	assert.Equal(t, 2, cnpj, "expired cache entries should be re-fetched")
	assert.Equal(t, 0, other)
}

func TestClient_CacheDistinctDocuments(t *testing.T) {
	srv := &serproCountingServer{}
	server := httptest.NewServer(srv.handler())
	defer server.Close()

	client := clientForSerpro(server, WithCacheTTL(time.Minute))
	ctx := context.Background()

	_, err := client.QueryCNPJ(ctx, "34.238.864/0001-68")
	require.NoError(t, err)
	_, err = client.QueryCNPJ(ctx, "11.222.333/0001-81")
	require.NoError(t, err, "a different document is a different cache key and must hit the server")

	token, cnpj, other := srv.counts()
	assert.Equal(t, 1, token)
	assert.Equal(t, 2, cnpj, "distinct documents must not share a cache entry")
	assert.Equal(t, 0, other)
}
