// Package brdocs provides a unified interface for Brazilian CPF and CNPJ validation.
//
// This package re-exports the most commonly used functions from the cpf and cnpj
// subpackages for convenience.
//
// For API queries, use the api subpackage directly.
//
// Example:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/LeonardoSilva-10/brdocs"
//	)
//
//	func main() {
//	    fmt.Println(brdocs.IsValidCPF("529.982.247-25"))
//	    fmt.Println(brdocs.IsValidCNPJ("11.222.333/0001-81"))
//	}
package brdocs

import (
	"github.com/LeonardoSilva-10/brdocs/cnpj"
	"github.com/LeonardoSilva-10/brdocs/cpf"
)

// IsValidCPF checks if a CPF string is valid.
func IsValidCPF(s string) bool {
	return cpf.Valid(s)
}

// IsValidCNPJ checks if a CNPJ string is valid.
func IsValidCNPJ(s string) bool {
	return cnpj.Valid(s)
}

// FormatCPF returns the CPF in standard format (XXX.XXX.XXX-XX).
func FormatCPF(s string) string {
	return cpf.Format(s)
}

// FormatCNPJ returns the CNPJ in standard format (XX.XXX.XXX/XXXX-XX).
func FormatCNPJ(s string) string {
	return cnpj.Format(s)
}

// MaskCPF partially masks a CPF, showing only the last 4 digits.
func MaskCPF(s string) string {
	return cpf.Mask(s)
}

// MaskCNPJ partially masks a CNPJ, showing only the last 4 digits.
func MaskCNPJ(s string) string {
	return cnpj.Mask(s)
}

// SanitizeCPF removes all non-numeric characters from a CPF string.
func SanitizeCPF(s string) string {
	return cpf.Sanitize(s)
}

// SanitizeCNPJ removes all non-numeric characters from a CNPJ string.
func SanitizeCNPJ(s string) string {
	return cnpj.Sanitize(s)
}

// CPFState returns the Brazilian state (UF) associated with the CPF.
func CPFState(s string) []string {
	return cpf.State(s)
}

// GenerateCPF returns a randomly generated valid CPF number.
func GenerateCPF() string {
	return cpf.Generate()
}

// GenerateCNPJ returns a randomly generated valid CNPJ number.
func GenerateCNPJ() string {
	return cnpj.Generate()
}

// CNPJIsHeadquarters checks if the CNPJ represents the headquarters (matriz).
func CNPJIsHeadquarters(s string) bool {
	return cnpj.IsHeadquarters(s)
}

// CNPJBranchCode returns the branch code from the CNPJ.
func CNPJBranchCode(s string) string {
	return cnpj.BranchCode(s)
}
