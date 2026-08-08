// Package cpf provides Brazilian CPF (Cadastro de Pessoas Físicas) validation.
//
// CPF is an 11-digit number where the last 2 digits are check digits calculated
// using the modulo 11 algorithm, as defined by Receita Federal.
//
// Reference: https://www.gov.br/receitafederal/pt-br
package cpf

import (
	cryptorand "crypto/rand"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	// ErrInvalidLength is returned when the CPF does not have 11 digits.
	ErrInvalidLength = errors.New("cpf: invalid length, expected 11 digits")
	// ErrRepeatedDigits is returned when all digits are the same (e.g., 111.111.111-11).
	ErrRepeatedDigits = errors.New("cpf: all digits are repeated")
	// ErrInvalidCheckDigit is returned when the check digits do not match.
	ErrInvalidCheckDigit = errors.New("cpf: invalid check digits")
)

// Valid checks if a CPF string is valid.
// It sanitizes the input, checks length, repeated digits, and verifies check digits.
func Valid(cpf string) bool {
	_, err := Parse(cpf)
	return err == nil
}

// Parse sanitizes and validates a CPF, returning the clean 11-digit string.
func Parse(cpf string) (string, error) {
	clean := Sanitize(cpf)

	if len(clean) != 11 {
		return "", ErrInvalidLength
	}

	if hasRepeatedDigits(clean) {
		return "", ErrRepeatedDigits
	}

	if !verifyCheckDigits(clean) {
		return "", ErrInvalidCheckDigit
	}

	return clean, nil
}

// Sanitize removes all non-numeric characters from the CPF string.
func Sanitize(cpf string) string {
	var b strings.Builder
	b.Grow(len(cpf))
	for _, r := range cpf {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Format returns the CPF in the standard format: XXX.XXX.XXX-XX.
// Returns empty string if the CPF is invalid.
func Format(cpf string) string {
	clean, err := Parse(cpf)
	if err != nil {
		return ""
	}
	return clean[:3] + "." + clean[3:6] + "." + clean[6:9] + "-" + clean[9:11]
}

// Mask partially masks a CPF, showing only the last 4 digits.
// Example: ***.***.***-25
func Mask(cpf string) string {
	clean, err := Parse(cpf)
	if err != nil {
		return ""
	}
	return "***.***.***-" + clean[9:11]
}

// Generate creates a valid random CPF number.
func Generate() string {
	digits := randomDigits(9)

	// Ensure the generated base is not all equal digits (e.g. 111.111.111),
	// which would produce an invalid CPF.
	if allDigitsEqual(digits) {
		digits[0] = (digits[0] + 1) % 10
	}

	// Calculate first check digit
	d1 := calculateCheckDigit(digits, 10)
	withD1 := append(append([]int(nil), digits...), d1)

	// Calculate second check digit
	d2 := calculateCheckDigit(withD1, 11)

	return digitsToString(append(withD1, d2))
}

// randomDigits returns n random decimal digits (0-9).
// It uses crypto/rand and falls back to a time-seeded PRNG if the OS entropy
// source is unavailable.
func randomDigits(n int) []int {
	digits := make([]int, n)
	buf := make([]byte, n)
	if _, err := cryptorand.Read(buf); err != nil {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := range digits {
			digits[i] = rng.Intn(10)
		}
		return digits
	}
	for i := range digits {
		digits[i] = int(buf[i] % 10)
	}
	return digits
}

// allDigitsEqual reports whether every digit in the slice has the same value.
func allDigitsEqual(digits []int) bool {
	first := digits[0]
	for _, d := range digits[1:] {
		if d != first {
			return false
		}
	}
	return true
}

// digitsToString converts a slice of digits into its string representation.
func digitsToString(digits []int) string {
	b := make([]byte, len(digits))
	for i, d := range digits {
		b[i] = byte('0' + d)
	}
	return string(b)
}

// hasRepeatedDigits checks if all digits in the CPF are the same.
func hasRepeatedDigits(cpf string) bool {
	first := cpf[0]
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != first {
			return false
		}
	}
	return true
}

// verifyCheckDigits validates the two check digits using modulo 11.
func verifyCheckDigits(cpf string) bool {
	digits := make([]int, 11)
	for i := 0; i < 11; i++ {
		digits[i] = int(cpf[i] - '0')
	}

	// First check digit
	d1 := calculateCheckDigit(digits[:9], 10)
	if d1 != digits[9] {
		return false
	}

	// Second check digit
	d2 := calculateCheckDigit(digits[:10], 11)
	return d2 == digits[10]
}

// calculateCheckDigit computes a single check digit for the given digits.
// weightStart is the initial weight (10 for first digit, 11 for second).
func calculateCheckDigit(digits []int, weightStart int) int {
	sum := 0
	for i, d := range digits {
		sum += d * (weightStart - i)
	}
	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

// State returns the Brazilian state (UF) associated with the CPF's 9th digit.
// Returns empty string if CPF is invalid.
//
// Mapping:
// 1: DF, GO, MS, MT, TO
// 2: AC, AM, AP, PA, RO, RR
// 3: CE, MA, PI
// 4: AL, PB, PE, RN
// 5: BA, SE
// 6: MG
// 7: ES, RJ
// 8: SP
// 9: PR, SC
// 0: RS
func State(cpf string) []string {
	clean, err := Parse(cpf)
	if err != nil {
		return nil
	}

	digit, _ := strconv.Atoi(string(clean[8]))
	states := map[int][]string{
		0: {"RS"},
		1: {"DF", "GO", "MS", "MT", "TO"},
		2: {"AC", "AM", "AP", "PA", "RO", "RR"},
		3: {"CE", "MA", "PI"},
		4: {"AL", "PB", "PE", "RN"},
		5: {"BA", "SE"},
		6: {"MG"},
		7: {"ES", "RJ"},
		8: {"SP"},
		9: {"PR", "SC"},
	}

	if s, ok := states[digit]; ok && len(s) > 0 {
		return s // return all possible states for that digit
	}
	return nil
}
