// Package cnpj provides Brazilian CNPJ (Cadastro Nacional da Pessoa Jurídica) validation.
//
// CNPJ is a 14-digit number where the last 2 digits are check digits calculated
// using the modulo 11 algorithm with weights 5,4,3,2,9,8,7,6,5,4,3,2.
//
// Reference: https://www.gov.br/receitafederal/pt-br
package cnpj

import (
	cryptorand "crypto/rand"
	"errors"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

var (
	// ErrInvalidLength is returned when the CNPJ does not have 14 digits.
	ErrInvalidLength = errors.New("cnpj: invalid length, expected 14 digits")
	// ErrRepeatedDigits is returned when all digits are the same.
	ErrRepeatedDigits = errors.New("cnpj: all digits are repeated")
	// ErrInvalidCheckDigit is returned when the check digits do not match.
	ErrInvalidCheckDigit = errors.New("cnpj: invalid check digits")
)

// Valid checks if a CNPJ string is valid.
func Valid(cnpj string) bool {
	_, err := Parse(cnpj)
	return err == nil
}

// Parse sanitizes and validates a CNPJ, returning the clean 14-digit string.
func Parse(cnpj string) (string, error) {
	clean := Sanitize(cnpj)

	if len(clean) != 14 {
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

// Sanitize removes all non-numeric characters from the CNPJ string.
func Sanitize(cnpj string) string {
	var b strings.Builder
	b.Grow(len(cnpj))
	for _, r := range cnpj {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Format returns the CNPJ in the standard format: XX.XXX.XXX/XXXX-XX.
// Returns empty string if the CNPJ is invalid.
func Format(cnpj string) string {
	clean, err := Parse(cnpj)
	if err != nil {
		return ""
	}
	return clean[:2] + "." + clean[2:5] + "." + clean[5:8] + "/" + clean[8:12] + "-" + clean[12:14]
}

// Mask partially masks a CNPJ, showing only the last 4 digits.
// Example: **.***.***/****-XX
func Mask(cnpj string) string {
	clean, err := Parse(cnpj)
	if err != nil {
		return ""
	}
	return "**.*" + clean[4:5] + ".***/****-" + clean[12:14]
}

// BranchCode returns the branch code from the CNPJ.
// The branch code is positions 9 to 12. "0001" indicates the headquarters.
func BranchCode(cnpj string) string {
	clean, err := Parse(cnpj)
	if err != nil {
		return ""
	}
	return clean[8:12]
}

// IsHeadquarters checks if the CNPJ represents the headquarters.
func IsHeadquarters(cnpj string) bool {
	return BranchCode(cnpj) == "0001"
}

// Generate creates a valid random CNPJ number.
func Generate() string {
	digits := randomDigits(12)

	// Ensure the generated base is not all equal digits (e.g. 11.111.111/1111-11),
	// which would produce an invalid CNPJ.
	if allDigitsEqual(digits) {
		digits[0] = (digits[0] + 1) % 10
	}

	// First check digit (12 digits -> weights 5,4,3,2,9,8,7,6,5,4,3,2)
	d1 := calculateCheckDigit(digits)
	withD1 := append(append([]int(nil), digits...), d1)

	// Second check digit (13 digits -> weights 6,5,4,3,2,9,8,7,6,5,4,3,2)
	d2 := calculateCheckDigit(withD1)

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

// hasRepeatedDigits checks if all digits in the CNPJ are the same.
func hasRepeatedDigits(cnpj string) bool {
	first := cnpj[0]
	for i := 1; i < len(cnpj); i++ {
		if cnpj[i] != first {
			return false
		}
	}
	return true
}

// verifyCheckDigits validates the two check digits using modulo 11.
func verifyCheckDigits(cnpj string) bool {
	digits := make([]int, 14)
	for i := 0; i < 14; i++ {
		digits[i] = int(cnpj[i] - '0')
	}

	// First check digit
	d1 := calculateCheckDigit(digits[:12])
	if d1 != digits[12] {
		return false
	}

	// Second check digit
	d2 := calculateCheckDigit(digits[:13])
	return d2 == digits[13]
}

// calculateCheckDigit computes a single check digit for the given digits.
// Uses weights: 5,4,3,2,9,8,7,6,5,4,3,2 for first digit
// and 6,5,4,3,2,9,8,7,6,5,4,3,2 for second digit.
func calculateCheckDigit(digits []int) int {
	var weights []int
	if len(digits) == 12 {
		weights = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	} else {
		weights = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	}

	sum := 0
	for i, d := range digits {
		sum += d * weights[i]
	}
	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}
