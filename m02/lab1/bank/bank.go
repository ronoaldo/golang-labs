package bank

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNegativeBallance = errors.New("bank: ballance bellow zero")
)

type Account interface {
	Deposit(reason string, value float64)
	Withdraw(reason string, value float64)
	Close() error
}

type account float64

func (a *account) Deposit(reason string, value float64) {
	logOp("DEPOSIT", reason, float64(*a), value)
	*a += account(value)
}

func (a *account) Withdraw(reason string, value float64) {
	logOp("WITHDRAW", reason, float64(*a), -value)
	*a -= account(value)
}

// https://go.dev/play/p/5TJveEXzJ9m
// This is not a performant aproach, using []byte instead of string
// performs better.
func (a *account) humanize() string {
	// Round and split integer/fraction parts
	aux := fmt.Sprintf("%.02f", *a)
	in, frac := strings.Split(aux, ".")[0], strings.Split(aux, ".")[1]

	// From end to first, add a '.' each 3 digits
	aux = ""
	i := len(in)
	for ; i > 3; i -= 3 {
		aux = "." + in[i-3:i] + aux
	}
	// Add remaining digits at the begining
	aux = in[:i] + aux

	return aux + "," + frac
}

func (a *account) String() string {
	return "Balance R$ " + a.humanize()
}

func (a *account) Close() error {
	logOp("CLOSE", "Account closing", float64(*a), 0)
	if *a < 0 {
		return ErrNegativeBallance
	}
	return nil
}

func NewAccount() Account {
	logOp("OPEN", "Account openning", 0, 0)
	var a account
	return &a
}

func fstr(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func logOp(op, desc string, prevBallance, value float64) {
	hhmm := time.Now().Format("03:04")
	ballance := prevBallance + value
	fmt.Printf(
		"%s %-6.6s %-12.12s %8s %8s %8s\n",
		hhmm,
		op,
		desc,
		fstr(prevBallance),
		fstr(value),
		fstr(ballance),
	)
}
