package humanize

import "fmt"

func ExampleHumanize2() {
	fmt.Println(Humanize2(0.1))
	fmt.Println(Humanize2(1234.56))
	fmt.Println(Humanize2(12345678.9))
	// Output:
	// 0,10
	// 1.234,56
	// 12.345.678,90
}
