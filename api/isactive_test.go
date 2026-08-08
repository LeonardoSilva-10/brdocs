package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCNPJResult_IsActive(t *testing.T) {
	tests := []struct {
		name string
		code string
		desc string
		want bool
	}{
		{"ativa por codigo", "2", "ATIVA", true},
		{"ativa por descricao", "", "ATIVA", true},
		{"ativa por descricao minúscula", "", "Situação cadastral ativa", true},
		{"nula", "1", "NULA", false},
		{"suspensa", "3", "SUSPENSA", false},
		{"inapta", "4", "INAPTA", false},
		{"baixada", "8", "BAIXADA", false},
		{"sem dados", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CNPJResult{SituacaoCadastral: SituacaoCadastral{Codigo: tt.code, Descricao: tt.desc}}
			assert.Equal(t, tt.want, r.IsActive())
		})
	}
}
