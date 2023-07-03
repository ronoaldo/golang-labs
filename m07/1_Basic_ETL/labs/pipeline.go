package main

import (
	"context"
	"flag"
	"log"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

var (
// Global flag variables
)

func init() {
	// Setup flag for parsing
}

func main() {
	flag.Parse()
	beam.Init()

	p := beam.NewPipeline()

	/*
	* Steps;
	* 1) Read something
	* 2) Transform something
	* 3) Write something
	 */

	log.Println("Building pipeline ...")
	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Error running pipeline: %v", err)
	}
}
