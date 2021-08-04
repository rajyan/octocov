package datastore

import (
	"context"
	"crypto/rand"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/k1LoW/octocov/gh"
	"github.com/k1LoW/octocov/report"
	"github.com/oklog/ulid/v2"
)

type BQ struct {
	client  *bigquery.Client
	dataset string
}

func NewBQ(client *bigquery.Client, dataset string) (*BQ, error) {
	return &BQ{
		client:  client,
		dataset: dataset,
	}, nil
}

type ReportRecord struct {
	Id                  string               `bigquery:"id"`
	Owner               string               `bigquery:"owner"`
	Repo                string               `bigquery:"repo"`
	Ref                 string               `bigquery:"ref"`
	Commit              string               `bigquery:"commit"`
	CoverageTotal       bigquery.NullInt64   `bigquery:"coverage_total"`
	CoverageCovered     bigquery.NullInt64   `bigquery:"coverage_covered"`
	CodeToTestRatioCode bigquery.NullInt64   `bigquery:"code_to_test_ratio_code"`
	CodeToTestRatioTest bigquery.NullInt64   `bigquery:"code_to_test_ratio_test"`
	TestExecutionTime   bigquery.NullFloat64 `bigquery:"test_execution_time"`
	Timestamp           time.Time            `bigquery:"timestamp"`
	Raw                 string               `bigquery:"raw"`
}

var reportsSchema = bigquery.Schema{
	&bigquery.FieldSchema{Name: "id", Type: bigquery.StringFieldType, Required: true},
	&bigquery.FieldSchema{Name: "owner", Type: bigquery.StringFieldType, Required: true},
	&bigquery.FieldSchema{Name: "repo", Type: bigquery.StringFieldType, Required: true},
	&bigquery.FieldSchema{Name: "ref", Type: bigquery.StringFieldType, Required: true},
	&bigquery.FieldSchema{Name: "commit", Type: bigquery.StringFieldType, Required: true},
	&bigquery.FieldSchema{Name: "coverage_total", Type: bigquery.IntegerFieldType, Required: false},
	&bigquery.FieldSchema{Name: "coverage_covered", Type: bigquery.IntegerFieldType, Required: false},
	&bigquery.FieldSchema{Name: "code_to_test_ratio_code", Type: bigquery.IntegerFieldType, Required: false},
	&bigquery.FieldSchema{Name: "code_to_test_ratio_test", Type: bigquery.IntegerFieldType, Required: false},
	&bigquery.FieldSchema{Name: "test_execution_time", Type: bigquery.NumericFieldType, Required: false},
	&bigquery.FieldSchema{Name: "timestamp", Type: bigquery.TimestampFieldType, Required: true},
	&bigquery.FieldSchema{Name: "raw", Type: bigquery.StringFieldType, Required: true},
}

func (b *BQ) Store(ctx context.Context, table string, r *report.Report) error {
	u := b.client.Dataset(b.dataset).Table(table).Uploader()
	owner, repo, err := gh.SplitRepository(r.Repository)
	if err != nil {
		return nil
	}
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return nil
	}
	rr := &ReportRecord{
		Id:        id.String(),
		Owner:     owner,
		Repo:      repo,
		Ref:       r.Ref,
		Commit:    r.Commit,
		Timestamp: r.Timestamp,
		Raw:       r.String(),
	}

	if r.Coverage != nil {
		rr.CoverageTotal = bigquery.NullInt64{
			Int64: int64(r.Coverage.Total),
			Valid: true,
		}
		rr.CoverageCovered = bigquery.NullInt64{
			Int64: int64(r.Coverage.Covered),
			Valid: true,
		}
	}
	if r.CodeToTestRatio != nil {
		rr.CodeToTestRatioCode = bigquery.NullInt64{
			Int64: int64(r.CodeToTestRatio.Code),
			Valid: true,
		}
		rr.CodeToTestRatioTest = bigquery.NullInt64{
			Int64: int64(r.CodeToTestRatio.Test),
			Valid: true,
		}
	}
	if r.TestExecutionTime != nil {
		rr.TestExecutionTime = bigquery.NullFloat64{
			Float64: *r.TestExecutionTime,
			Valid:   true,
		}
	}
	return u.Put(ctx, []*ReportRecord{rr})
}

func (b *BQ) CreateTable(ctx context.Context, table string) error {
	metaData := &bigquery.TableMetadata{
		Schema: reportsSchema,
	}
	tableRef := b.client.Dataset(b.dataset).Table(table)
	if err := tableRef.Create(ctx, metaData); err != nil {
		return err
	}
	return nil
}