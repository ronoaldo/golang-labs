package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"reflect"

	_ "cloud.google.com/go/bigquery"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	_ "github.com/apache/beam/sdks/v2/go/pkg/beam/io/bigqueryio"
	_ "github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/options/gcpopts"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

type CommonLog struct {
	UserID string `json:"user_id" bigquery:"user_id"`
	IP     string `json:"ip"      bigquery:"ip"`

	// Item 1: Permitir que estes campos sejam nulos
	// Dica: consulte a documentação da biblioteca cliente do Bigquery
	// em https://cloud.google.com/go
	Lat float64 `json:"lat" bigquery:"lat"`
	Lng float64 `json:"lng" bigquery:"lng"`

	Timestamp string `json:"timestamp" bigquery:"timestamp"`

	HttpRequest string `json:"http_request" bigquery:"http_request"`

	// Item 2: Evitar que este campo seja inserido no bigquery
	UserAgent string `json:"user_agent" bigquery:"user_agent"`

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

// Item 3: Aplicar algum filtro nos dados. Pode utilizar qualquer
// critério para filtro, por exemplo, filtrar requisições com menos
// do que 100 bytes de tamanho.
func Filter(row CommonLog, emit func(CommonLog)) {
	emit(row)
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
	beam.RegisterDoFn(Filter)

	// Item 4: Configure como parâmetros de linha de comando as variáveis
	// inputPath, outputPath e tableName utilizando o pacote flag.
	// Dica: Consulte a sintaxe das funções em https://pkg.go.dev/flag
}

func main() {
	flag.Parse()
	beam.Init()

	ctx := context.Background()
	project := gcpopts.GetProject(ctx)

	p := beam.NewPipeline()
	s := p.Root()

	// Item 5: Implementar o código para ler o arquivo de entrada
	// do Cloud Storage (inputPath), utilizando o textio
	// Dica: consulte a sintaxe das funções em:
	// https://pkg.go.dev/github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio
	_ = s
	lines, _ := beam.PCollection{}, inputPath

	// Item 6: Branch 1: escrever as linhas no arquivo de saída (outputPath)
	// utilizando o textio.
	// Dica: consulte a sintaxe das funções em:
	// https://pkg.go.dev/github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio
	_, _ = lines, outputPath

	// Item 7: Branch 2: analisar os dados JSON e filtrá-los, escrevendo o resultado
	// diretamente no Bigquery, utilizando o bigqueryio.
	// https://pkg.go.dev/github.com/apache/beam/sdks/v2/go/pkg/beam/io/bigqueryio
	commonLogs := lines
	filteredLogs := commonLogs
	_, _, _ = filteredLogs, project, tableName

	// Main pipeline
	log.Println("Building pipeline ...")
	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Error running pipeline: %v", err)
	}
}

/*
Para testar seu pipeline, depois de ter configurado o diretório BASE_DIR como no exercício 1,
execute os comandos abaixo no Cloud Shell

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
