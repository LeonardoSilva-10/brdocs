package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCPFResponse(t *testing.T) {
	t.Run("extracts data from HTML", func(t *testing.T) {
		html := `
			<html>
			<body>
				<span class="nome">JOÃO DA SILVA</span>
				<div>Situação Cadastral</div><div>REGULAR</div>
				<div>Data de Inscrição</div><div>01/01/1990</div>
				<div>Dígito Verificador</div><div>91</div>
			</body>
			</html>
		`
		result, err := parseCPFResponse(html, "12345678901")
		require.NoError(t, err)
		assert.Equal(t, "JOÃO DA SILVA", result.Nome)
		assert.Equal(t, "REGULAR", result.SituacaoCadastral)
		assert.Equal(t, "01/01/1990", result.DataInscricao)
		assert.Equal(t, "91", result.DigitoVerificador)
		assert.Equal(t, "12345678901", result.CPF)
		assert.False(t, result.ConsultaDataHora.IsZero())
	})

	t.Run("returns error for empty response", func(t *testing.T) {
		_, err := parseCPFResponse("<html></html>", "12345678901")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to parse response")
	})

	t.Run("detects error in response", func(t *testing.T) {
		html := `<html><body>CPF inválido ou não encontrado</body></html>`
		_, err := parseCPFResponse(html, "12345678901")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CPF not found or invalid")
	})
}

func TestReceitaFederalClient_ConsultaCPF(t *testing.T) {
	t.Run("returns error for invalid CPF length", func(t *testing.T) {
		client := NewClient()
		rf := newReceitaFederalClient(client)
		_, err := rf.ConsultaCPF(context.Background(), "123", "01/01/1990")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid CPF length")
	})
}
