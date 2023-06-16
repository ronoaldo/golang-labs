package main

import (
	"fmt"
	"strconv"
	"time"
)

var saldo, saldoAnterior float64

func AbreConta(saldoInicial float64) {
	saldo = saldoInicial
	saldoAnterior = saldo
}

func Credito(desc string, valor float64) {
	saldoAnterior = saldo
	saldo += valor
	Log("CREDITO", desc, valor)
}

func Debito(desc string, valor float64) {
	saldoAnterior = saldo
	saldo -= valor
	Log("DEBITO", desc, valor)
}

func fstr(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func Log(op, desc string, valor float64) {
	hhmm := time.Now().Format("03:04")
	fmt.Printf(
		"%s %-6.6s %-12.12s %8s %8s %8s\n",
		hhmm,
		op,
		desc,
		fstr(saldoAnterior),
		fstr(valor),
		fstr(saldo),
	)
}

func main() {
	AbreConta(100.00)
	Credito("Depósito à vista", 50.00)
	Debito("Saque", 70.00)
}
