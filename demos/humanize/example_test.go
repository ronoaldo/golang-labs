package humanize

import "fmt"

func ExampleHumanize() {
	fmt.Println(Humanize(0.1))
	fmt.Println(Humanize(1234.56))
	fmt.Println(Humanize(12345678.9))
	// Output:
	// 0,10
	// 1.234,56
	// 12.345.678,90
}
