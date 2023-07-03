package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"reflect"

	"cloud.google.com/go/bigquery"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/bigqueryio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/options/gcpopts"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

type CommonLog struct {
	UserID string `json:"user_id" bigquery:"user_id"`
	IP     string `json:"ip"      bigquery:"ip"`

	// Allow nulls
	Lat bigquery.NullFloat64 `json:"lat" bigquery:"lat"`
	Lng bigquery.NullFloat64 `json:"lng" bigquery:"lng"`

	Timestamp string `json:"timestamp" bigquery:"timestamp"`

	HttpRequest string `json:"http_request" bigquery:"http_request"`

	// Avoid writting to bigquery
	UserAgent string `json:"user_agent" bigquery:"-"`

	HttpResponse int `json:"http_response" bigquery:"http_response"`
	NumBytes     int `json:"num_bytes" bigquery:"num_bytes"`
}

func JSONToCommonLog(line string, emit func(CommonLog)) {
	dst := CommonLog{}
	if err := json.Unmarshal([]byte(line), &dst); err != nil {
		log.Fatalf("Error parsing line: %v", err)
	}
	emit(dst)
}

func FilterSmallRequests(row CommonLog, emit func(CommonLog)) {
	if row.NumBytes < 120 {
		emit(row)
	}
}

var (
	inputPath  string
	outputPath string
	tableName  string
)

func init() {
	// Beam type/dofn registration
	beam.RegisterType(reflect.TypeOf((*CommonLog)(nil)).Elem())
	beam.RegisterDoFn(JSONToCommonLog)
	beam.RegisterDoFn(FilterSmallRequests)

	// Flag setup
	flag.StringVar(&inputPath, "input_path", "", "Path to events.json")
	flag.StringVar(&outputPath, "output_path", "", "Path to coldline storage bucket")
	flag.StringVar(&tableName, "table_name", "", "Bigquery table name")
}

func main() {
	flag.Parse()
	beam.Init()

	ctx := context.Background()
	project := gcpopts.GetProject(ctx)

	p := beam.NewPipeline()
	s := p.Root()

	// Read from GCS
	lines := textio.Read(s.Scope("ReadFromGCS"), inputPath)

	// Branch 1: Write data as it was read to GCS
	textio.Write(s.Scope("WriteRawToGCS"), outputPath, lines)

	// Branch 2: Parse GCS and ingest to Bigquery
	commonLogs := beam.ParDo(s.Scope("ParseJson"), JSONToCommonLog, lines)
	filteredLogs := beam.ParDo(s.Scope("FilterFn"), FilterSmallRequests, commonLogs)
	bigqueryio.Write(s.Scope("WriteToBQ"), project, tableName, filteredLogs)

	// Main pipeline
	log.Println("Building pipeline ...")
	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Error running pipeline: %v", err)
	}
}

/*
To execute, run:

	cd $BASE_DIR

	export PROJECT_ID=$(gcloud config get-value project)
	export REGION='us-central1'
	export BUCKET=gs://${PROJECT_ID}
	export COLDLINE_BUCKET=${BUCKET}-coldline
	export PIPELINE_FOLDER=${BUCKET}
	export RUNNER=DataflowRunner
	export INPUT_PATH=${PIPELINE_FOLDER}/events.json
	export OUTPUT_PATH=${PIPELINE_FOLDER}-coldline/output
	export TABLE_NAME=${PROJECT_ID}:logs.logs_filtered

	go run pipeline.go \
		--project=$PROJECT_ID \
		--region=$REGION \
		--staging_location=$PIPELINE_FOLDER/staging \
		--temp_location=$PIPELINE_FOLDER/temp \
		--job_name=my-pipeline-$(date +%Y%m%d%H%M%S) \
		--runner=$RUNNER \
		--input_path=$INPUT_PATH \
		--output_path=$OUTPUT_PATH \
		--table_name=${TABLE_NAME}
*/
