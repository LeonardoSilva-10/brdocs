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

func TestThirdPartyClient_QueryCPFCNPJ(t *testing.T) {
	t.Run("returns error without API key", func(t *testing.T) {
		client := NewClient()
		tp := newThirdPartyClient(client)
		_, err := tp.QueryCPFCNPJ(context.Background(), "52998224725")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API key is required")
	})

	t.Run("successfully queries CPF", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "52998224725")
			assert.Contains(t, r.URL.Path, "test-api-key")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(CPFCNPJResult{
				Status: "OK",
				Name:   "JOÃO DA SILVA",
			})
		}))
		defer server.Close()

		client := NewClient(
			WithCPFCNPJAPIKey("test-api-key"),
			WithBaseURL("cpfcnpj", server.URL),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.QueryCPFCNPJ(ctx, "529.982.247-25")
		require.NoError(t, err)
		assert.Equal(t, "OK", result.Status)
		assert.Equal(t, "JOÃO DA SILVA", result.Name)
	})

	t.Run("returns error on failed query", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(CPFCNPJResult{
				Status:  "ERROR",
				Message: "registro não encontrado",
			})
		}))
		defer server.Close()

		client := NewClient(
			WithCPFCNPJAPIKey("test-api-key"),
			WithBaseURL("cpfcnpj", server.URL),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := client.QueryCPFCNPJ(ctx, "529.982.247-25")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query failed")
	})
}

func TestThirdPartyClient_QueryCPFBrasil(t *testing.T) {
	t.Run("returns error without API key", func(t *testing.T) {
		client := NewClient()
		tp := newThirdPartyClient(client)
		_, err := tp.QueryCPFBrasil(context.Background(), "52998224725")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API key is required")
	})

	t.Run("successfully queries CPF", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/cpf/52998224725", r.URL.Path)
			assert.Equal(t, "test-key", r.Header.Get("X-API-Key"))
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(struct {
				Success bool `json:"success"`
				Data    struct {
					CPF  string `json:"CPF"`
					Name string `json:"NOME"`
					Sexo string `json:"SEXO"`
					Nasc string `json:"NASC"`
				} `json:"data"`
			}{
				Success: true,
				Data: struct {
					CPF  string `json:"CPF"`
					Name string `json:"NOME"`
					Sexo string `json:"SEXO"`
					Nasc string `json:"NASC"`
				}{
					CPF:  "52998224725",
					Name: "JOÃO DA SILVA",
					Sexo: "M",
					Nasc: "1990-01-01",
				},
			})
		}))
		defer server.Close()

		client := NewClient(
			WithCPFBrasilAPIKey("test-key"),
			WithBaseURL("cpfbrasil", server.URL),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.QueryCPFBrasil(ctx, "529.982.247-25")
		require.NoError(t, err)
		assert.Equal(t, "52998224725", result.CPF)
		assert.Equal(t, "JOÃO DA SILVA", result.Name)
		assert.Equal(t, "1990-01-01", result.BirthDate)
	})

	t.Run("returns error on failed query", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(struct {
				Success bool `json:"success"`
			}{
				Success: false,
			})
		}))
		defer server.Close()

		client := NewClient(
			WithCPFBrasilAPIKey("test-key"),
			WithBaseURL("cpfbrasil", server.URL),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := client.QueryCPFBrasil(ctx, "529.982.247-25")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query failed")
	})
}
