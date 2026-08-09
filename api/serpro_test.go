package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerproClient_GetToken(t *testing.T) {
	t.Run("returns error without credentials", func(t *testing.T) {
		client := NewClient()
		sc := newSerproClient(client)
		_, err := sc.getToken(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "consumer key and secret are required")
	})

	t.Run("successfully obtains token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/token", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.Header.Get("Authorization"), "Basic ")
			assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serproToken{
				AccessToken: "test-token-123",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			})
		}))
		defer server.Close()

		client := NewClient(
			WithSerproCredentials("key", "secret"),
			WithBaseURL("serpro", server.URL),
		)
		sc := newSerproClient(client)

		token, err := sc.getToken(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "test-token-123", token)
	})

	t.Run("caches token and reuses", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serproToken{
				AccessToken: "cached-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			})
		}))
		defer server.Close()

		client := NewClient(
			WithSerproCredentials("key", "secret"),
			WithBaseURL("serpro", server.URL),
		)
		sc := newSerproClient(client)

		_, err := sc.getToken(context.Background())
		require.NoError(t, err)
		_, err = sc.getToken(context.Background())
		require.NoError(t, err)

		assert.Equal(t, 1, callCount)
	})
}

func TestSerproClient_ConsultaCNPJ(t *testing.T) {
	t.Run("successfully queries CNPJ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/token":
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(serproToken{
					AccessToken: "test-token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				})
			case "/consulta-cnpj-df/v2/basica/34238864000168":
				assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(CNPJResult{
					NI:              "34238864000168",
					NomeEmpresarial: "SERVICO DE PROCESSAMENTO DE DADOS SERPRO",
					SituacaoCadastral: SituacaoCadastral{
						Codigo:    "2",
						Descricao: "ATIVA",
					},
				})
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		client := NewClient(
			WithSerproCredentials("key", "secret"),
			WithBaseURL("serpro", server.URL),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.ConsultaCNPJ(ctx, "34.238.864/0001-68")
		require.NoError(t, err)
		assert.Equal(t, "34238864000168", result.NI)
		assert.Equal(t, "SERVICO DE PROCESSAMENTO DE DADOS SERPRO", result.NomeEmpresarial)
		assert.True(t, result.IsActive())
	})

	t.Run("returns error for invalid CNPJ response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(serproToken{
					AccessToken: "test-token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		client := NewClient(
			WithSerproCredentials("key", "secret"),
			WithBaseURL("serpro", server.URL),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := client.ConsultaCNPJ(ctx, "00.000.000/0000-00")
		assert.Error(t, err)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewClient(
			WithSerproCredentials("key", "secret"),
			WithBaseURL("serpro", server.URL),
			WithTimeout(50*time.Millisecond),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := client.ConsultaCNPJ(ctx, "34.238.864/0001-68")
		assert.Error(t, err)
	})
}
