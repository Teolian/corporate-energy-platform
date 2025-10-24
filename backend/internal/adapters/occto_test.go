package adapters

import (
	"os"
	"testing"

	"github.com/teo/aversome/backend/internal/reserve"
)

func TestOCCTOAdapter_ParseCSV(t *testing.T) {
	tests := []struct {
		name        string
		csvFile     string
		date        string
		wantAreas   int
		wantWarning bool
		wantError   bool
	}{
		{
			name:      "valid CSV with Tokyo and Kansai",
			csvFile:   "testdata/occto-sample.csv",
			date:      "2025-10-24",
			wantAreas: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewOCCTOAdapter()

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
			if resp.Date != tt.date {
				t.Errorf("Date = %v, want %v", resp.Date, tt.date)
			}

			if len(resp.Areas) != tt.wantAreas {
				t.Errorf("Areas length = %d, want %d", len(resp.Areas), tt.wantAreas)
			}

			// Check source attribution
			if resp.Source.Name != "OCCTO" {
				t.Errorf("Source.Name = %v, want OCCTO", resp.Source.Name)
			}

			if resp.Source.URL == "" {
				t.Error("Source.URL is empty")
			}

			// Check specific area data
			if len(resp.Areas) >= 1 {
				// First area should be Tokyo
				tokyo := resp.Areas[0]
				if tokyo.Area != "tokyo" {
					t.Errorf("First area = %v, want tokyo", tokyo.Area)
				}
				if tokyo.ReserveMarginPct != 5.1 {
					t.Errorf("Tokyo ReserveMarginPct = %v, want 5.1", tokyo.ReserveMarginPct)
				}
				// 5.1% is in [5, 8) range → should be "watch"
				if tokyo.Status != reserve.StatusWatch {
					t.Errorf("Tokyo Status = %v, want %v", tokyo.Status, reserve.StatusWatch)
				}
			}

			if len(resp.Areas) >= 2 {
				// Second area should be Kansai
				kansai := resp.Areas[1]
				if kansai.Area != "kansai" {
					t.Errorf("Second area = %v, want kansai", kansai.Area)
				}
				if kansai.ReserveMarginPct != 8.9 {
					t.Errorf("Kansai ReserveMarginPct = %v, want 8.9", kansai.ReserveMarginPct)
				}
				// 8.9% is >= 8 → should be "stable"
				if kansai.Status != reserve.StatusStable {
					t.Errorf("Kansai Status = %v, want %v", kansai.Status, reserve.StatusStable)
				}
			}

			// Check warning expectation
			hasWarning := resp.Meta != nil && resp.Meta.Warning != ""
			if hasWarning != tt.wantWarning {
				t.Errorf("Warning present = %v, want %v", hasWarning, tt.wantWarning)
			}
		})
	}
}

func TestOCCTOAdapter_normalizeArea(t *testing.T) {
	adapter := NewOCCTOAdapter()

	tests := []struct {
		input string
		want  string
	}{
		{"東京", "tokyo"},
		{"Tokyo", "tokyo"},
		{"TOKYO", "tokyo"},
		{"関西", "kansai"},
		{"Kansai", "kansai"},
		{"KANSAI", "kansai"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := adapter.normalizeArea(tt.input)
			if got != tt.want {
				t.Errorf("normalizeArea(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestReserveDeriveStatus(t *testing.T) {
	tests := []struct {
		pct  float64
		want reserve.Status
	}{
		{10.0, reserve.StatusStable},
		{8.0, reserve.StatusStable},
		{7.9, reserve.StatusWatch},
		{5.0, reserve.StatusWatch},
		{4.9, reserve.StatusTight},
		{2.0, reserve.StatusTight},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := reserve.DeriveStatus(tt.pct)
			if got != tt.want {
				t.Errorf("DeriveStatus(%v) = %v, want %v", tt.pct, got, tt.want)
			}
		})
	}
}
