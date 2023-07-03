package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"reflect"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/bigqueryio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

type CommonLog struct {
	UserID string `json:"user_id" bigquery:"user_id"`
	IP     string `json:"ip"      bigquery:"ip"`

	Lat float64 `json:"lat" bigquery:"lat"`
	Lng float64 `json:"lng" bigquery:"lng"`

	Timestamp string `json:"timestamp" bigquery:"timestamp"`

	HttpRequest  string `json:"http_request" bigquery:"http_request"`
	UserAgent    string `json:"user_agent" bigquery:"user_agent"`
	HttpResponse int    `json:"http_response" bigquery:"http_response"`
	NumBytes     int    `json:"num_bytes" bigquery:"num_bytes"`
}

func JSONToCommonLog(line string, emit func(CommonLog)) {
	dst := &CommonLog{}
	if err := json.Unmarshal([]byte(line), dst); err != nil {
		log.Fatalf("Error parsing line: %v", err)
	}
	emit(*dst)
}

func init() {
	// DoFn/Types registration
	beam.RegisterType(reflect.TypeOf((*CommonLog)(nil)).Elem())
	beam.RegisterDoFn(JSONToCommonLog)
}

func main() {
	flag.Parse()
	beam.Init()

	p := beam.NewPipeline()

	// TODO: altere este valor antes de executar o pipeline
	project := "YOUR-PROJECT-ID-HERE"

	input := "gs://" + project + "/events.json"
	output := project + ":logs.logs"

	s := p.Root()
	lines := textio.Read(s, input)
	commonLogs := beam.ParDo(s, JSONToCommonLog, lines)
	bigqueryio.Write(s, project, output, commonLogs)

	log.Println("Building pipeline ...")
	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Error running pipeline: %v", err)
	}
}
