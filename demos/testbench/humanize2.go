package humanize

import (
	"strconv"
	"strings"
)

func Humanize2(f float64) string {
	// Round and split integer/fraction parts
	buff := []byte(strconv.FormatFloat(f, 'f', 2, 64))
	buff[len(buff)-3] = ','

	before := len(buff)

	// Allocate space for '.' only once
	inc := int((len(buff) - 3) / 3)
	buff = append(buff, []byte(strings.Repeat("#", inc))...)

	end := len(buff) - 1
	diff := len(buff) - before
	j := end - diff
	i := end
	count := 0
	for j >= 0 {
		buff[i] = buff[j]
		i--
		j--
		count++

		if (count-3)%3 == 0 && count-3 > 0 {
			buff[i] = '.'
			i--
		}
	}
	return string(buff)
}
