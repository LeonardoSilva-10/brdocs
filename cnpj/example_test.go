package cnpj_test

import (
	"fmt"

	"github.com/LeonardoSilva-10/brdocs/cnpj"
)

func ExampleValid() {
	fmt.Println(cnpj.Valid("11.222.333/0001-81"))
	fmt.Println(cnpj.Valid("00.000.000/0000-00"))
	// Output:
	// true
	// false
}

func ExampleFormat() {
	fmt.Println(cnpj.Format("11222333000181"))
	fmt.Println(cnpj.Format("00000000000000"))
	// Output:
	// 11.222.333/0001-81
	//
}

func ExampleBranchCode_headquarters() {
	fmt.Println(cnpj.BranchCode("11222333000181"))
	// Output:
	// 0001
}

func ExampleBranchCode_branch() {
	fmt.Println(cnpj.BranchCode("11222333000262"))
	// Output:
	// 0002
}

func ExampleIsHeadquarters() {
	fmt.Println(cnpj.IsHeadquarters("11222333000181"))
	fmt.Println(cnpj.IsHeadquarters("11222333000262"))
	// Output:
	// true
	// false
}
