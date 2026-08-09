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
		{"active by code", "2", "ATIVA", true},
		{"active by description", "", "ATIVA", true},
		{"active by lowercase description", "", "registration status ativa", true},
		{"null", "1", "NULA", false},
		{"suspended", "3", "SUSPENSA", false},
		{"inapt", "4", "INAPTA", false},
		{"deregistered", "8", "BAIXADA", false},
		{"no data", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CNPJResult{RegistrationStatus: RegistrationStatus{Code: tt.code, Description: tt.desc}}
			assert.Equal(t, tt.want, r.IsActive())
		})
	}
}
