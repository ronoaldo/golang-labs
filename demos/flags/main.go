package main

import (
	"fmt"
	"os"

	"flag"
)

var file string

func main() {
	flag.StringVar(&file, "file", "default.json", "O `arquivo` para ser lido.")
	flag.Parse()

	fmt.Printf("%#v\n", os.Args)
	fmt.Printf("file=%v\n", file)
}
