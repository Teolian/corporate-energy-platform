// Package main provides a pipeline job to fetch demand data for Tokyo/Kansai and normalize to JSON.
// Usage: go run main.go -area tokyo -date 2025-10-24
// Output: /public/data/jp/{area}/demand-YYYY-MM-DD.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/teo/aversome/backend/internal/adapters"
	"github.com/teo/aversome/backend/internal/demand"
	"github.com/teo/aversome/backend/pkg/timeutil"
)

func main() {
	var date, area string
	flag.StringVar(&date, "date", "", "Date in YYYY-MM-DD format (defaults to today)")
	flag.StringVar(&area, "area", "tokyo", "Area: tokyo or kansai")
	flag.Parse()

	// Default to today if no date provided
	if date == "" {
		date = timeutil.FormatDate(time.Now())
	}

	// Validate date format
	if _, err := timeutil.ParseDate(date); err != nil {
		log.Fatalf("Invalid date format: %v", err)
	}

	// Validate area
	if area != "tokyo" && area != "kansai" {
		log.Fatalf("Invalid area: %s (must be 'tokyo' or 'kansai')", area)
	}

	log.Printf("Fetching %s demand data for %s...", area, date)

	// Get CSV path and adapter based on area
	var csvPath string
	var resp *demand.Response
	var err error

	switch area {
	case "tokyo":
		csvPath = filepath.Join("internal", "adapters", "testdata", "tepco-sample.csv")
		f, err := os.Open(csvPath)
		if err != nil {
			log.Fatalf("Failed to open TEPCO CSV: %v", err)
		}
		defer f.Close()

		adapter := adapters.NewTEPCOAdapter()
		resp, err = adapter.ParseCSV(f, date)
		if err != nil {
			log.Fatalf("Failed to parse TEPCO CSV: %v", err)
		}

	case "kansai":
		csvPath = filepath.Join("internal", "adapters", "testdata", "kansai-sample.csv")
		f, err := os.Open(csvPath)
		if err != nil {
			log.Fatalf("Failed to open Kansai CSV: %v", err)
		}
		defer f.Close()

		adapter := adapters.NewKansaiAdapter()
		resp, err = adapter.ParseCSV(f, date)
		if err != nil {
			log.Fatalf("Failed to parse Kansai CSV: %v", err)
		}
	}

	log.Printf("Parsed %d data points", len(resp.Series))

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create output directory
	outputDir := filepath.Join("..", "public", "data", "jp", area)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write JSON file
	outputPath := filepath.Join(outputDir, fmt.Sprintf("demand-%s.json", date))
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}

	log.Printf("✓ Successfully wrote %s (%d bytes)", outputPath, len(jsonData))
}

// fetchCSV fetches CSV data from a URL (future implementation).
func fetchCSV(url string) (io.ReadCloser, error) {
	// TODO: Implement HTTP fetch with retry logic
	return nil, fmt.Errorf("HTTP fetching not yet implemented")
}
