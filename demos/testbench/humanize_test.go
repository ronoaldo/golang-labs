package humanize

import (
	"testing"
)

func TestHumanize(t *testing.T) {
	testCases := []struct {
		in  float64
		out string
	}{
		{0.1, "0,10"},
		{123.45, "123,45"},
		{12345.67, "12.345,67"},
	}

	for _, tc := range testCases {
		out := Humanize(tc.in)
		t.Logf("humanize(%v) => out=%#v", tc.in, out)

		if out != tc.out {
			t.Errorf("Unexpected output: want=%v, got=%v", tc.out, out)
		}
	}
}
