package main

import (
	"fmt"
	"math/big"
)

func Fatorial(n int) {
	fmt.Printf("%d! = ", n)
	var fat int64
	for fat = 1; n >= 2; n = n - 1 {
		fmt.Printf("%d * ", n)
		fat = fat * int64(n)
	}
	fmt.Printf("1 = %v\n", fat)
}

func FatorialBig(num int64) {
	fmt.Printf("%d! = ", num)
	fat, n := big.NewInt(1), big.NewInt(num)
	one, two := big.NewInt(1), big.NewInt(2)
	for n.Cmp(two) >= 0 {
		fmt.Printf("%d * ", n)
		fat.Mul(fat, n)
		n.Sub(n, one)
	}
	fmt.Printf("1 = %v\n", fat)
}

func main() {
	Fatorial(3)
	Fatorial(1)
	Fatorial(5)
	Fatorial(10)
	Fatorial(30)
	FatorialBig(30)
}
