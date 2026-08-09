package cpf_test

import (
	"fmt"

	"github.com/LeonardoSilva-10/brdocs/cpf"
)

func ExampleValid() {
	fmt.Println(cpf.Valid("529.982.247-25"))
	fmt.Println(cpf.Valid("111.111.111-11"))
	// Output:
	// true
	// false
}

func ExampleFormat() {
	fmt.Println(cpf.Format("52998224725"))
	fmt.Println(cpf.Format("11111111111"))
	// Output:
	// 529.982.247-25
	//
}

func ExampleMask() {
	fmt.Println(cpf.Mask("52998224725"))
	// Output:
	// ***.***.***-25
}

func ExampleState() {
	// CPF 52998224806 -> 9th digit = 8 -> SP
	fmt.Println(cpf.State("52998224806"))
	// Output:
	// [SP]
}
