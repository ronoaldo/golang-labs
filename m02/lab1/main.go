package main

import (
	"fmt"
	"labsolutions/m02/lab1/bank"
)

func main() {
	a := bank.NewAccount()
	a.Deposit("Abertura de conta", 100.00)
	a.Withdraw("Compra no débito", 50.00)
	fmt.Printf("%v\n", a)
	a.Withdraw("Compra no débito", 100.00)
	fmt.Printf("%v\n", a)

	if err := a.Close(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	a.Deposit("Encerramento de conta", 50.00)
	if err := a.Close(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
