// Package cnpj provides validation utilities for Brazilian CNPJ numbers.
//
// A CNPJ (Cadastro Nacional da Pessoa Jurídica) is a 14-digit number issued by
// Receita Federal. The last two digits are check digits computed using
// the modulo 11 algorithm with specific weights.
//
// Basic usage:
//
//	if cnpj.Valid("11.222.333/0001-81") {
//	    fmt.Println("Valid CNPJ")
//	}
//
// The package also provides branch identification:
//
//	fmt.Println(cnpj.BranchCode("11.222.333/0001-81")) // "0001"
//	fmt.Println(cnpj.IsHeadquarters("11.222.333/0001-81")) // true
package cnpj
