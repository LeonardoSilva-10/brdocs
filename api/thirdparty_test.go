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

func TestThirdPartyClient_ConsultaCPFCPFCNPJ(t *testing.T) {
	t.Run("returns error without API key", func(t *testing.T) {
		client := NewClient()
		tp := newThirdPartyClient(client)
		_, err := tp.ConsultaCPFCPFCNPJ(context.Background(), "52998224725")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API key is required")
	})

	t.Run("successfully queries CPF", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Path, "52998224725")
			assert.Contains(t, r.URL.Path, "test-api-key")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "OK",
				"nome":   "JOÃO DA SILVA",
			})
		}))
		defer server.Close()

		client := NewClient(
			WithCPFCNPJAPIKey("test-api-key"),
			WithBaseURL("cpfcnpj", server.URL),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := client.ConsultaCPFCPFCNPJ(ctx, "529.982.247-25")
		require.NoError(t, err)
		m := result.(map[string]interface{})
		assert.Equal(t, "OK", m["status"])
		assert.Equal(t, "JOÃO DA SILVA", m["nome"])
	})
}

func TestThirdPartyClient_ConsultaCPFBrasil(t *testing.T) {
	t.Run("returns error without API key", func(t *testing.T) {
		client := NewClient()
		tp := newThirdPartyClient(client)
		_, err := tp.ConsultaCPFBrasil(context.Background(), "52998224725")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API key is required")
	})

	t.Run("successfully queries CPF", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/cpf/52998224725", r.URL.Path)
			assert.Equal(t, "test-key", r.Header.Get("X-API-Key"))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				Success bool `json:"success"`
				Data    struct {
					CPF  string `json:"CPF"`
					Nome string `json:"NOME"`
					Sexo string `json:"SEXO"`
					Nasc string `json:"NASC"`
				} `json:"data"`
			}{
				Success: true,
				Data: struct {
					CPF  string `json:"CPF"`
					Nome string `json:"NOME"`
					Sexo string `json:"SEXO"`
					Nasc string `json:"NASC"`
				}{
					CPF:  "52998224725",
					Nome: "JOÃO DA SILVA",
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

		result, err := client.ConsultaCPFBrasil(ctx, "529.982.247-25")
		require.NoError(t, err)
		assert.Equal(t, "52998224725", result.CPF)
		assert.Equal(t, "JOÃO DA SILVA", result.Nome)
		assert.Equal(t, "1990-01-01", result.DataNascimento)
	})

	t.Run("returns error on failed query", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
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

		_, err := client.ConsultaCPFBrasil(ctx, "529.982.247-25")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query failed")
	})
}
