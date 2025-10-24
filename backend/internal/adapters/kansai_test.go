package adapters

import (
	"os"
	"testing"

	"github.com/teo/aversome/backend/internal/demand"
)

func TestKansaiAdapter_ParseCSV(t *testing.T) {
	tests := []struct {
		name        string
		csvFile     string
		date        string
		wantArea    demand.Area
		wantSeries  int
		wantWarning bool
		wantError   bool
	}{
		{
			name:       "valid CSV with forecast",
			csvFile:    "testdata/kansai-sample.csv",
			date:       "2025-10-24",
			wantArea:   demand.AreaKansai,
			wantSeries: 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewKansaiAdapter()

			// Read test CSV
			f, err := os.Open(tt.csvFile)
			if err != nil {
				t.Fatalf("failed to open test file: %v", err)
			}
			defer f.Close()

			// Parse CSV
			resp, err := adapter.ParseCSV(f, tt.date)

			// Check error expectation
			if (err != nil) != tt.wantError {
				t.Errorf("ParseCSV() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError {
				return
			}

			// Validate response
			if resp.Area != tt.wantArea {
				t.Errorf("Area = %v, want %v", resp.Area, tt.wantArea)
			}

			if resp.Date != tt.date {
				t.Errorf("Date = %v, want %v", resp.Date, tt.date)
			}

			if resp.Timezone != "Asia/Tokyo" {
				t.Errorf("Timezone = %v, want Asia/Tokyo", resp.Timezone)
			}

			if len(resp.Series) != tt.wantSeries {
				t.Errorf("Series length = %d, want %d", len(resp.Series), tt.wantSeries)
			}

			// Check first and last points
			if len(resp.Series) > 0 {
				first := resp.Series[0]
				if first.DemandMW != 12345.0 { // 1234.5 * 10
					t.Errorf("First DemandMW = %v, want 12345.0", first.DemandMW)
				}
				if first.ForecastMW == nil {
					t.Error("First ForecastMW is nil, expected value")
				} else if *first.ForecastMW != 12562.0 { // 1256.2 * 10
					t.Errorf("First ForecastMW = %v, want 12562.0", *first.ForecastMW)
				}

				last := resp.Series[len(resp.Series)-1]
				if last.DemandMW != 12895.0 { // 1289.5 * 10
					t.Errorf("Last DemandMW = %v, want 12895.0", last.DemandMW)
				}
			}

			// Check source attribution
			if resp.Source.Name != "Kansai" {
				t.Errorf("Source.Name = %v, want Kansai", resp.Source.Name)
			}
			if resp.Source.URL == "" {
				t.Error("Source.URL is empty")
			}

			// Check warning
			hasWarning := resp.Meta != nil && resp.Meta.Warning != ""
			if hasWarning != tt.wantWarning {
				t.Errorf("Has warning = %v, want %v", hasWarning, tt.wantWarning)
			}
		})
	}
}
