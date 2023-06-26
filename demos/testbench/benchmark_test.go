package humanize

import (
	"flag"
	"math"
	"testing"
)

var testNum = math.MaxFloat64

func init() {
	flag.Float64Var(&testNum, "testnum", math.MaxFloat64, "The maximum number to use when benchmarking")
}

func BenchmarkHumanize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Humanize(testNum)
	}
}

func BenchmarkHumanize2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Humanize2(testNum)
	}
}
