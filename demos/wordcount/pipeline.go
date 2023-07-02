package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/transforms/stats"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

var (
	splitWordsRe = regexp.MustCompile(`[a-zA-Z]+('[a-z])?`)
	output       string
)

func splitWords(line string, emit func(string)) {
	line = strings.ToLower(strings.TrimSpace(line))
	for _, word := range splitWordsRe.FindAllString(line, -1) {
		emit(word)
	}
}

func formatOutput(w string, c int) string {
	return fmt.Sprintf("%s: %v", w, c)
}

func init() {
	flag.StringVar(&output, "output", "wordcount.txt", "The `output path` to write to.")
}

func main() {
	flag.Parse()
	beam.Init()

	p, s := beam.NewPipelineWithRoot()

	lines := textio.Read(s.Scope("ReadFiles"), "gs://apache-beam-samples/shakespeare/*")
	words := beam.ParDo(s.Scope("SplitWords"), splitWords, lines)
	counted := stats.Count(s.Scope("CountWords"), words)
	formatted := beam.ParDo(s.Scope("FormatOutput"), formatOutput, counted)
	textio.Write(s.Scope("WriteToOutput"), output, formatted)

	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	}
}
