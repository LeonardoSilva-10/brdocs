package brdocs_test

import (
	"fmt"

	"github.com/LeonardoSilva-10/brdocs"
)

func ExampleIsValidCPF() {
	fmt.Println(brdocs.IsValidCPF("529.982.247-25"))
	fmt.Println(brdocs.IsValidCPF("111.111.111-11"))
	// Output:
	// true
	// false
}

func ExampleIsValidCNPJ() {
	fmt.Println(brdocs.IsValidCNPJ("11.222.333/0001-81"))
	fmt.Println(brdocs.IsValidCNPJ("00.000.000/0000-00"))
	// Output:
	// true
	// false
}

func ExampleFormatCPF() {
	fmt.Println(brdocs.FormatCPF("52998224725"))
	// Output:
	// 529.982.247-25
}

func ExampleFormatCNPJ() {
	fmt.Println(brdocs.FormatCNPJ("11222333000181"))
	// Output:
	// 11.222.333/0001-81
}
