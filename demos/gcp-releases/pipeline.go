package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"reflect"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/bigqueryio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/options/gcpopts"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/transforms/stats"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

const (
	PublicDataReleaseNotes = "bigquery-public-data:google_cloud_release_notes.release_notes"
)

var (
	output string
)

/* Types */

type ReleaseNotes struct {
	Type        bigquery.NullString `bigquery:"release_note_type" json:"release_note_type"`
	PublishedAt civil.Date          `bigquery:"published_at" json:"published_at"`

	ProductID      bigquery.NullInt64  `bigquery:"product_id" json:"product_id"`
	ProductName    bigquery.NullString `bigquery:"product_name" json:"product_name"`
	ProductVersion bigquery.NullString `bigquery:"product_version_name" json:"product_version_name"`

	Description bigquery.NullString `bigquery:"description" json:"description"`
}

var ReleaseNotesType = reflect.TypeOf(ReleaseNotes{})

func (r *ReleaseNotes) String() string {
	if r == nil {
		return "ReleaseNotes<nil>"
	}

	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

/* DoFNs */

func AsJSON(line ReleaseNotes, emit func(string)) {
	emit(line.String())
}

func ExtractType(r ReleaseNotes) string {
	return r.Type.String()
}

func FormatCount(key string, value int) string {
	return fmt.Sprintf("{\"%v\": %v}", key, value)
}

type FilterReleses struct {
	Type string
}

func (fn *FilterReleses) ProcessElement(r ReleaseNotes, emit func(ReleaseNotes)) {
	if r.Type.String() == fn.Type {
		emit(r)
	}
}

/* Transforms */

func ReadReleases(s beam.Scope, project, table string) beam.PCollection {
	return bigqueryio.Read(s.Scope("ReadReleases"), project, table, ReleaseNotesType)
}

func WriteReleases(s beam.Scope, releases beam.PCollection) {
	lines := beam.ParDo(s.Scope("AsJSON"), AsJSON, releases)
	textio.Write(s.Scope("WriteReleases"), output+"releases.jsonl", lines)
}

func WriteReportCountByType(s beam.Scope, releases beam.PCollection) {
	types := beam.ParDo(s.Scope("ReportExtractType"), ExtractType, releases)
	counted := stats.Count(s.Scope("ReportCount"), types)
	lines := beam.ParDo(s.Scope("FormatCount"), FormatCount, counted)
	textio.Write(s.Scope("ReportWrite"), output+"report_count_by_type.jsonl", lines)
}

func WriteSecurityBulletin(s beam.Scope, releases beam.PCollection) {
	securityNotes := beam.ParDo(
		s.Scope("FilterBySecurityBulletin"),
		&FilterReleses{Type: "SECURITY_BULLETIN"}, releases)
	lines := beam.ParDo(s.Scope("AsJSON"), AsJSON, securityNotes)
	textio.Write(s.Scope("WriteSecurityBulletin"), output+"security_bulletin.jsonl", lines)
}

/* Pipeline */

func init() {
	flag.StringVar(&output, "output-dir", "./", "The `output directory` to write to.")

	beam.RegisterType(ReleaseNotesType)

	beam.RegisterDoFn(AsJSON)
	beam.RegisterDoFn(ExtractType)
	beam.RegisterDoFn(FormatCount)
	beam.RegisterDoFn(&FilterReleses{})
}

func main() {
	flag.Parse()
	beam.Init()

	ctx := context.Background()
	project := gcpopts.GetProjectFromFlagOrEnvironment(ctx)

	p := beam.NewPipeline()
	s := p.Root()
	// Read all release notes
	releases := ReadReleases(s, project, PublicDataReleaseNotes)
	// 1) Save them as a JSONL file
	WriteReleases(s, releases)
	// 2) Save a report of count by type
	WriteReportCountByType(s, releases)
	// 3) Create a Security Bulletin text file for "SECURITY_BULLETIN"
	WriteSecurityBulletin(s, releases)

	if err := beamx.Run(context.Background(), p); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	}
}
