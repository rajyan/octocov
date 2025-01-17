package report

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/tenntenn/golden"
)

func TestCustomMetricSetTable(t *testing.T) {
	tests := []struct {
		s *CustomMetricSet
	}{
		{&CustomMetricSet{}},
		{&CustomMetricSet{
			Key:  "benchmark_0",
			Name: "Benchmark-0",
			Metrics: []*CustomMetric{
				{Key: "N", Name: "Number of iterations", Value: 1000.0, Unit: ""},
				{Key: "NsPerOp", Name: "Nanoseconds per iteration", Value: 676.5, Unit: " ns/op"},
			},
		}},
		{&CustomMetricSet{
			Key:  "benchmark_1",
			Name: "Benchmark-1",
			Metrics: []*CustomMetric{
				{Key: "N", Name: "Number of iterations", Value: 1500.0, Unit: ""},
				{Key: "NsPerOp", Name: "Nanoseconds per iteration", Value: 1340.0, Unit: " ns/op"},
			},
		}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := tt.s.Table()
			f := filepath.Join("custom_metrics", fmt.Sprintf("custom_metric_set_table.%d", i))
			if os.Getenv("UPDATE_GOLDEN") != "" {
				golden.Update(t, testdataDir(t), f, got)
				return
			}
			if diff := golden.Diff(t, testdataDir(t), f, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestCustomMetricSetOut(t *testing.T) {
	tests := []struct {
		s *CustomMetricSet
	}{
		{&CustomMetricSet{}},
		{&CustomMetricSet{
			Key:  "benchmark_0",
			Name: "Benchmark-0",
			Metrics: []*CustomMetric{
				{Key: "N", Name: "Number of iterations", Value: 1000.0, Unit: ""},
				{Key: "NsPerOp", Name: "Nanoseconds per iteration", Value: 676.5, Unit: " ns/op"},
			},
			report: &Report{
				Ref:      "main",
				Commit:   "1234567890",
				covPaths: []string{"testdata/cover.out"},
			},
		}},
		{&CustomMetricSet{
			Key:  "benchmark_1",
			Name: "Benchmark-1",
			Metrics: []*CustomMetric{
				{Key: "N", Name: "Number of iterations", Value: 1500.0, Unit: ""},
				{Key: "NsPerOp", Name: "Nanoseconds per iteration", Value: 1340.0, Unit: " ns/op"},
			},
			report: &Report{
				Ref:      "main",
				Commit:   "1234567890",
				covPaths: []string{"testdata/cover.out"},
			},
		}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := new(bytes.Buffer)
			if err := tt.s.Out(got); err != nil {
				t.Fatal(err)
			}
			f := filepath.Join("custom_metrics", fmt.Sprintf("custom_metric_set_out.%d", i))
			if os.Getenv("UPDATE_GOLDEN") != "" {
				golden.Update(t, testdataDir(t), f, got)
				return
			}
			if diff := golden.Diff(t, testdataDir(t), f, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestCustomMetricsSetValidate(t *testing.T) {
	tests := []struct {
		in      *CustomMetricSet
		wantErr bool
	}{
		{&CustomMetricSet{}, true},
		{&CustomMetricSet{
			Key: "key",
			Metrics: []*CustomMetric{
				{Key: "count", Value: 1000.0},
				{Key: "ns_per_op", Value: 676.0, Unit: "ns/op"},
			},
		}, false},
		{&CustomMetricSet{
			Key:     "key",
			Metrics: []*CustomMetric{},
		}, true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if err := tt.in.Validate(); err != nil {
				if !tt.wantErr {
					t.Error(err)
				}
				return
			}
			if tt.wantErr {
				t.Error("want error")
			}
		})
	}
}

func TestDiffCustomMetricSetTable(t *testing.T) {
	tests := []struct {
		a *CustomMetricSet
		b *CustomMetricSet
	}{
		{
			&CustomMetricSet{
				Key:  "benchmark_0",
				Name: "Benchmark-0",
				Metrics: []*CustomMetric{
					{Key: "N", Name: "Number of iterations", Value: 1000.0, Unit: ""},
					{Key: "NsPerOp", Name: "Nanoseconds per iteration", Value: 676.5, Unit: " ns/op"},
				},
				report: &Report{
					Ref:      "main",
					Commit:   "1234567890",
					covPaths: []string{"testdata/cover.out"},
				},
			},
			nil,
		},
		{
			&CustomMetricSet{
				Key:  "benchmark_0",
				Name: "Benchmark-0",
				Metrics: []*CustomMetric{
					{Key: "N", Name: "Number of iterations", Value: 1000.0, Unit: ""},
					{Key: "NsPerOp", Name: "Nanoseconds per iteration", Value: 676.5, Unit: " ns/op"},
				},
				report: &Report{
					Ref:      "main",
					Commit:   "1234567890",
					covPaths: []string{"testdata/cover.out"},
				},
			},
			&CustomMetricSet{
				Key:  "benchmark_0",
				Name: "Benchmark-0",
				Metrics: []*CustomMetric{
					{Key: "N", Name: "Number of iterations", Value: 9393.0, Unit: ""},
					{Key: "NsPerOp", Name: "Nanoseconds per iteration", Value: 456.0, Unit: " ns/op"},
				},
				report: &Report{
					Ref:      "main",
					Commit:   "2345678901",
					covPaths: []string{"testdata/cover.out"},
				},
			},
		},
	}

	t.Setenv("GITHUB_SERVER_URL", "https://github.com")
	t.Setenv("GITHUB_REPOSITORY", "owner/repo")
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			d := tt.a.Compare(tt.b)
			got := d.Table()
			f := filepath.Join("custom_metrics", fmt.Sprintf("diff_custom_metric_set_table.%d", i))
			if os.Getenv("UPDATE_GOLDEN") != "" {
				golden.Update(t, testdataDir(t), f, got)
				return
			}
			if diff := golden.Diff(t, testdataDir(t), f, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
