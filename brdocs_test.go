package brdocs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidCPF(t *testing.T) {
	assert.True(t, IsValidCPF("529.982.247-25"))
	assert.False(t, IsValidCPF("111.111.111-11"))
	assert.False(t, IsValidCPF("invalid"))
}

func TestIsValidCNPJ(t *testing.T) {
	assert.True(t, IsValidCNPJ("11.222.333/0001-81"))
	assert.False(t, IsValidCNPJ("00.000.000/0000-00"))
	assert.False(t, IsValidCNPJ("invalid"))
}

func TestFormatCPF(t *testing.T) {
	assert.Equal(t, "529.982.247-25", FormatCPF("52998224725"))
	assert.Equal(t, "", FormatCPF("11111111111"))
}

func TestFormatCNPJ(t *testing.T) {
	assert.Equal(t, "11.222.333/0001-81", FormatCNPJ("11222333000181"))
	assert.Equal(t, "", FormatCNPJ("00000000000000"))
}

func TestMaskCPF(t *testing.T) {
	assert.Equal(t, "***.***.***-25", MaskCPF("52998224725"))
	assert.Equal(t, "", MaskCPF("11111111111"))
}

func TestSanitizeCPF(t *testing.T) {
	assert.Equal(t, "52998224725", SanitizeCPF("529.982.247-25"))
	assert.Equal(t, "52998224725", SanitizeCPF("529 982 247 25"))
}

func TestSanitizeCNPJ(t *testing.T) {
	assert.Equal(t, "11222333000181", SanitizeCNPJ("11.222.333/0001-81"))
	assert.Equal(t, "11222333000181", SanitizeCNPJ("11 222 333 0001 81"))
}

func TestCPFState(t *testing.T) {
	// 52998224210 -> 9o digito = 2 -> "AC", "AM", "AP", "PA", "RO", "RR"
	// The 9th digit of the CPF indicates the state, and for digit 2, it can correspond to multiple states.
	assert.Contains(t, CPFState("52998224210"), "AM")
	assert.Equal(t, []string(nil), CPFState("invalid")) // Invalid CPF should return nil
}

func TestCNPJIsHeadquarters(t *testing.T) {
	assert.True(t, CNPJIsHeadquarters("11222333000181"))
	assert.False(t, CNPJIsHeadquarters("11222333000262"))
	assert.False(t, CNPJIsHeadquarters("invalid"))
}

func TestCNPJBranchCode(t *testing.T) {
	assert.Equal(t, "0001", CNPJBranchCode("11222333000181"))
	assert.Equal(t, "0002", CNPJBranchCode("11222333000262"))
	assert.Equal(t, "", CNPJBranchCode("invalid"))
}

func TestGenerateCPF(t *testing.T) {
	g := GenerateCPF()
	assert.Len(t, g, 11)
	assert.True(t, IsValidCPF(g))
}

func TestGenerateCNPJ(t *testing.T) {
	g := GenerateCNPJ()
	assert.Len(t, g, 14)
	assert.True(t, IsValidCNPJ(g))
}
