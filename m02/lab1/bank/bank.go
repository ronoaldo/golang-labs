package bank

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNegativeBallance = errors.New("bank: ballance bellow zero.")
)

type Account interface {
	Deposit(reason string, value float64)
	Withdraw(reason string, value float64)
	Close() error
}

type account float64

func (a *account) Deposit(reason string, value float64) {
	*a += account(value)
}

func (a *account) Withdraw(reason string, value float64) {
	*a -= account(value)
}

func (a *account) String() string {
	out := fmt.Sprintf("Balance R$ %.02f", float64(*a))
	return strings.Replace(out, ".", ",", -1)
}

func (a *account) Close() error {
	if *a < 0 {
		return ErrNegativeBallance
	}
	return nil
}

func NewAccount() Account {
	var a account
	return &a
}
