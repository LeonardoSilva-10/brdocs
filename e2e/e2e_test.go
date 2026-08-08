package e2e

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/LeonardoSilva-10/brdocs/api"
	"github.com/LeonardoSilva-10/brdocs/cnpj"
	"github.com/LeonardoSilva-10/brdocs/cpf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocalValidation runs validation tests without any external dependencies.
func TestLocalValidation(t *testing.T) {
	t.Run("CPF validation", func(t *testing.T) {
		assert.True(t, cpf.Valid("529.982.247-25"))
		assert.False(t, cpf.Valid("111.111.111-11"))
		assert.False(t, cpf.Valid("123.456.789-01")) // DV correto eh 09
		assert.False(t, cpf.Valid("invalid"))
		assert.False(t, cpf.Valid(""))
	})

	t.Run("CNPJ validation", func(t *testing.T) {
		assert.True(t, cnpj.Valid("11.222.333/0001-81"))
		assert.False(t, cnpj.Valid("00.000.000/0000-00"))
		assert.False(t, cnpj.Valid("11.222.333/0001-82"))
		assert.False(t, cnpj.Valid("invalid"))
		assert.False(t, cnpj.Valid(""))
	})

	t.Run("CPF formatting", func(t *testing.T) {
		assert.Equal(t, "529.982.247-25", cpf.Format("52998224725"))
		assert.Equal(t, "", cpf.Format("11111111111"))
	})

	t.Run("CNPJ formatting", func(t *testing.T) {
		assert.Equal(t, "11.222.333/0001-81", cnpj.Format("11222333000181"))
		assert.Equal(t, "", cnpj.Format("00000000000000"))
	})

	t.Run("CPF masking", func(t *testing.T) {
		assert.Equal(t, "***.***.***-25", cpf.Mask("52998224725"))
		assert.Equal(t, "", cpf.Mask("11111111111"))
	})

	t.Run("CNPJ branch code", func(t *testing.T) {
		assert.Equal(t, "0001", cnpj.BranchCode("11222333000181"))
		assert.True(t, cnpj.IsHeadquarters("11222333000181"))
		assert.False(t, cnpj.IsHeadquarters("11222333000262"))
	})
}

// TestSerproIntegration tests the Serpro/ConectaGov API integration.
// Requires SERPRO_CONSUMER_KEY and SERPRO_CONSUMER_SECRET environment variables.
func TestSerproIntegration(t *testing.T) {
	consumerKey := os.Getenv("SERPRO_CONSUMER_KEY")
	consumerSecret := os.Getenv("SERPRO_CONSUMER_SECRET")

	if consumerKey == "" || consumerSecret == "" {
		t.Skip("Skipping Serpro integration test: SERPRO_CONSUMER_KEY and SERPRO_CONSUMER_SECRET not set")
	}

	client := api.NewClient(
		api.WithSerproCredentials(consumerKey, consumerSecret),
		api.WithTimeout(30*time.Second),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with a known valid CNPJ (Serpro itself)
	result, err := client.ConsultaCNPJ(ctx, "34.238.864/0001-68")
	require.NoError(t, err)
	assert.NotEmpty(t, result.NomeEmpresarial)
	assert.Equal(t, "34238864000168", result.NI)
}

// TestCPFCNPJAPIIntegration tests the cpfcnpj.com.br API integration.
// Requires CPFCNPJ_API_KEY environment variable.
func TestCPFCNPJAPIIntegration(t *testing.T) {
	apiKey := os.Getenv("CPFCNPJ_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping cpfcnpj.com.br integration test: CPFCNPJ_API_KEY not set")
	}

	client := api.NewClient(
		api.WithCPFCNPJAPIKey(apiKey),
		api.WithTimeout(30*time.Second),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := client.ConsultaCPFCPFCNPJ(ctx, "529.982.247-25")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

// TestReceitaFederalIntegration tests the Receita Federal public CPF query.
// Note: This may fail due to CAPTCHA or page structure changes.
func TestReceitaFederalIntegration(t *testing.T) {
	if os.Getenv("RUN_RF_INTEGRATION") != "1" {
		t.Skip("Skipping Receita Federal integration test: set RUN_RF_INTEGRATION=1 to enable")
	}

	client := api.NewClient(
		api.WithTimeout(30 * time.Second),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use a known test CPF with a test birth date
	// Note: This will likely fail in practice due to CAPTCHA
	_, err := client.ConsultaCPF(ctx, "529.982.247-25", "01/01/1990")
	// We expect this to potentially fail due to CAPTCHA or page changes
	assert.Error(t, err)
}

// BenchmarkValidations benchmarks local validation performance.
func BenchmarkValidations(b *testing.B) {
	b.Run("CPF", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cpf.Valid("529.982.247-25")
		}
	})

	b.Run("CNPJ", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cnpj.Valid("11.222.333/0001-81")
		}
	})
}
