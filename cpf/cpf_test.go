package cpf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{"valid formatted", "529.982.247-25", true},
		{"valid unformatted", "52998224725", true},
		{"valid with spaces", "529 982 247 25", true},
		{"valid known", "000.000.001-91", true},
		{"invalid repeated", "111.111.111-11", false},
		{"invalid repeated 2", "000.000.000-00", false},
		{"invalid repeated 3", "999.999.999-99", false},
		{"invalid length short", "1234567890", false},
		{"invalid length long", "123456789012", false},
		{"invalid letters", "abc.def.ghi-jk", false},
		{"empty", "", false},
		{"invalid check digit 1", "529.982.247-26", false},
		{"invalid check digit 2", "529.982.247-24", false},
		{"invalid dv wrong", "123.456.789-01", false}, // DV correto eh 09
		{"valid with extra chars", "CPF: 529.982.247-25", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Valid(tt.cpf))
		})
	}
}

func TestParse(t *testing.T) {
	t.Run("returns clean CPF on success", func(t *testing.T) {
		clean, err := Parse("529.982.247-25")
		require.NoError(t, err)
		assert.Equal(t, "52998224725", clean)
	})

	t.Run("returns error for invalid CPF", func(t *testing.T) {
		_, err := Parse("111.111.111-11")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrRepeatedDigits)
	})

	t.Run("returns error for wrong length", func(t *testing.T) {
		_, err := Parse("123")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidLength)
	})

	t.Run("returns error for non-numeric", func(t *testing.T) {
		_, err := Parse("abc.def.ghi-jk")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidLength)
	})
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"529.982.247-25", "52998224725"},
		{"52998224725", "52998224725"},
		{"529 982 247 25", "52998224725"},
		{"CPF: 529.982.247-25", "52998224725"},
		{"abc123def", "123"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, Sanitize(tt.input))
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"52998224725", "529.982.247-25"},
		{"529.982.247-25", "529.982.247-25"},
		{"11111111111", ""},
		{"invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, Format(tt.input))
		})
	}
}

func TestMask(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"52998224725", "***.***.***-25"},
		{"529.982.247-25", "***.***.***-25"},
		{"11111111111", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, Mask(tt.input))
		})
	}
}

func TestState(t *testing.T) {
	// CPFs validos com digito de UF conhecido no indice 8:
	// 52998224210 -> 9o digito = 2 -> "AC", "AM", "AP", "PA", "RO", "RR"
	// 52998224130 -> 9o digito = 1 -> "DF", "GO", "MS", "MT", "TO"
	// 52998224482 -> 9o digito = 4 -> "AL", "PB", "PE", "RN"
	tests := []struct {
		cpf      string
		expected string
	}{
		{"52998224210", "AC"},
		{"52998224130", "DF"},
		{"52998224482", "AL"},
		// {"invalid", []string(nil)}, // Invalid CPF should return nil
	}

	for _, tt := range tests {
		t.Run(tt.cpf, func(t *testing.T) {
			assert.Contains(t, State(tt.cpf), tt.expected)
		})
	}
}

func TestCalculateCheckDigit(t *testing.T) {
	tests := []struct {
		digits      []int
		weightStart int
		expected    int
	}{
		// CPF 529.982.247-25 -> first 9 digits: 5,2,9,9,8,2,2,4,7
		// weights: 10,9,8,7,6,5,4,3,2
		// sum = 5*10 + 2*9 + 9*8 + 9*7 + 8*6 + 2*5 + 2*4 + 4*3 + 7*2
		//     = 50 + 18 + 72 + 63 + 48 + 10 + 8 + 12 + 14 = 295
		// 295 % 11 = 295 - 26*11 = 295 - 286 = 9
		// 11 - 9 = 2 -> first check digit = 2
		{[]int{5, 2, 9, 9, 8, 2, 2, 4, 7}, 10, 2},
		// Now with first check digit included: 5,2,9,9,8,2,2,4,7,2
		// weights: 11,10,9,8,7,6,5,4,3,2
		// sum = 5*11 + 2*10 + 9*9 + 9*8 + 8*7 + 2*6 + 2*5 + 4*4 + 7*3 + 2*2
		//     = 55 + 20 + 81 + 72 + 56 + 12 + 10 + 16 + 21 + 4 = 347
		// 347 % 11 = 347 - 31*11 = 347 - 341 = 6
		// 11 - 6 = 5 -> second check digit = 5
		{[]int{5, 2, 9, 9, 8, 2, 2, 4, 7, 2}, 11, 5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.expected, calculateCheckDigit(tt.digits, tt.weightStart))
		})
	}
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Valid("529.982.247-25")
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse("529.982.247-25")
	}
}

func BenchmarkSanitize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sanitize("529.982.247-25")
	}
}

func TestAllRepeatedDigits(t *testing.T) {
	for d := byte('0'); d <= '9'; d++ {
		cpf := strings.Repeat(string(d), 11)
		assert.False(t, Valid(cpf), "CPF %s should be invalid", cpf)
	}
}
