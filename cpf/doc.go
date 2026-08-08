// Package cpf provides validation utilities for Brazilian CPF numbers.
//
// A CPF (Cadastro de Pessoas Físicas) is an 11-digit number issued by
// Receita Federal. The last two digits are check digits computed using
// the modulo 11 algorithm.
//
// Basic usage:
//
//	if cpf.Valid("529.982.247-25") {
//	    fmt.Println("Valid CPF")
//	}
//
// The package also provides formatting and masking utilities:
//
//	fmt.Println(cpf.Format("52998224725")) // "529.982.247-25"
//	fmt.Println(cpf.Mask("52998224725"))   // "***.***.***-25"
package cpf
