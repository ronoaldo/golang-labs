package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var filename string
var outfile string

func init() {
	flag.StringVar(&filename, "s", "gs://arki1/gopher.txt", "The `GCS storage path` to read from")
	flag.StringVar(&outfile, "o", "", "The `output` to record to, defaults to os.Stdout")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		log.Printf("Error initializing storage client: %v", err)
		return
	}

	parts := strings.Split(filename, "//")
	if len(parts) != 2 {
		log.Printf("Invalid gcs path: %v; expected gs://bucket/object", filename)
		return
	}

	path := strings.Split(parts[1], "/")
	if len(path) < 2 {
		log.Printf("Expected bucket/object path, got: %v", filename)
	}
	bucket, obj := path[0], strings.Join(path[1:], "/")

	fd, err := client.Bucket(bucket).Object(obj).NewReader(ctx)
	if err != nil {
		log.Printf("Error opening object '%s' in bucket '%s':\n%v", obj, bucket, err)
		return
	}

	var out io.Writer = os.Stdout
	if outfile != "" {
		out, err = os.Create(outfile)
		if err != nil {
			log.Printf("Error opening file for write '%v': %v", outfile, err)
			return
		}
	}

	if _, err := io.Copy(out, fd); err != nil {
		log.Printf("Error reading object: %v", err)
	}
}
