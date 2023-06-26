package humanize

import (
	"fmt"
	"strings"
)

func Humanize(f float64) string {
	// Round and split integer/fraction parts
	aux := fmt.Sprintf("%.02f", f)
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
