package cnpj

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name     string
		cnpj     string
		expected bool
	}{
		{"valid formatted", "11.222.333/0001-81", true},
		{"valid unformatted", "11222333000181", true},
		{"valid with spaces", "11 222 333 0001 81", true},
		{"invalid repeated", "00.000.000/0000-00", false},
		{"invalid repeated 2", "11.111.111/1111-11", false},
		{"invalid repeated 3", "99.999.999/9999-99", false},
		{"invalid length short", "1122233300018", false},
		{"invalid length long", "112223330001811", false},
		{"invalid letters", "ab.cde.fgh/ijkl-mn", false},
		{"empty", "", false},
		{"valid known", "34.238.864/0001-68", true}, // Serpro test CNPJ
		{"invalid check digit 1", "11.222.333/0001-82", false},
		{"invalid check digit 2", "11.222.333/0001-80", false},
		{"valid with extra chars", "CNPJ: 11.222.333/0001-81", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Valid(tt.cnpj))
		})
	}
}

func TestParse(t *testing.T) {
	t.Run("returns clean CNPJ on success", func(t *testing.T) {
		clean, err := Parse("11.222.333/0001-81")
		require.NoError(t, err)
		assert.Equal(t, "11222333000181", clean)
	})

	t.Run("returns error for invalid CNPJ", func(t *testing.T) {
		_, err := Parse("00.000.000/0000-00")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrRepeatedDigits)
	})

	t.Run("returns error for wrong length", func(t *testing.T) {
		_, err := Parse("123")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidLength)
	})
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"11.222.333/0001-81", "11222333000181"},
		{"11222333000181", "11222333000181"},
		{"11 222 333 0001 81", "11222333000181"},
		{"CNPJ: 11.222.333/0001-81", "11222333000181"},
		{"abc11222333000181def", "11222333000181"},
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
		{"11222333000181", "11.222.333/0001-81"},
		{"11.222.333/0001-81", "11.222.333/0001-81"},
		{"00000000000000", ""},
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
		{"11222333000181", "**.*2.***/****-81"},
		{"11.222.333/0001-81", "**.*2.***/****-81"},
		{"00000000000000", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, Mask(tt.input))
		})
	}
}

func TestBranchCode(t *testing.T) {
	// 11222333000181 -> headquarters (0001)
	// 11222333000262 -> branch (0002), valid CNPJ
	tests := []struct {
		input    string
		expected string
	}{
		{"11222333000181", "0001"},
		{"11222333000262", "0002"},
		{"invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, BranchCode(tt.input))
		})
	}
}

func TestIsHeadquarters(t *testing.T) {
	assert.True(t, IsHeadquarters("11222333000181"))
	assert.False(t, IsHeadquarters("11222333000262"))
	assert.False(t, IsHeadquarters("invalid"))
}

func TestCalculateCheckDigit(t *testing.T) {
	tests := []struct {
		digits   []int
		expected int
	}{
		// CNPJ 11.222.333/0001-81 -> first 12 digits: 1,1,2,2,2,3,3,3,0,0,0,1
		// weights: 5,4,3,2,9,8,7,6,5,4,3,2
		// sum = 1*5 + 1*4 + 2*3 + 2*2 + 2*9 + 3*8 + 3*7 + 3*6 + 0*5 + 0*4 + 0*3 + 1*2
		//     = 5 + 4 + 6 + 4 + 18 + 24 + 21 + 18 + 0 + 0 + 0 + 2 = 102
		// 102 % 11 = 102 - 9*11 = 102 - 99 = 3
		// 11 - 3 = 8 -> first check digit = 8
		{[]int{1, 1, 2, 2, 2, 3, 3, 3, 0, 0, 0, 1}, 8},
		// Now with first check digit: 1,1,2,2,2,3,3,3,0,0,0,1,8
		// weights: 6,5,4,3,2,9,8,7,6,5,4,3,2
		// sum = 1*6 + 1*5 + 2*4 + 2*3 + 2*2 + 3*9 + 3*8 + 3*7 + 0*6 + 0*5 + 0*4 + 1*3 + 8*2
		//     = 6 + 5 + 8 + 6 + 4 + 27 + 24 + 21 + 0 + 0 + 0 + 3 + 16 = 120
		// 120 % 11 = 120 - 10*11 = 120 - 110 = 10
		// 11 - 10 = 1 -> second check digit = 1
		{[]int{1, 1, 2, 2, 2, 3, 3, 3, 0, 0, 0, 1, 8}, 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.expected, calculateCheckDigit(tt.digits))
		})
	}
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Valid("11.222.333/0001-81")
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse("11.222.333/0001-81")
	}
}

func BenchmarkSanitize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sanitize("11.222.333/0001-81")
	}
}

func TestAllRepeatedDigits(t *testing.T) {
	for d := byte('0'); d <= '9'; d++ {
		cnpj := strings.Repeat(string(d), 14)
		assert.False(t, Valid(cnpj), "CNPJ %s should be invalid", cnpj)
	}
}
