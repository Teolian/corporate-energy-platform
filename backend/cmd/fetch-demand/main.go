// Package main provides a pipeline job to fetch Tokyo demand data and normalize to JSON.
// Usage: go run fetch-tokyo-demand.go [YYYY-MM-DD]
// Output: /public/data/jp/tokyo/demand-YYYY-MM-DD.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/teo/aversome/backend/internal/adapters"
	"github.com/teo/aversome/backend/pkg/timeutil"
)

func main() {
	var date string
	flag.StringVar(&date, "date", "", "Date in YYYY-MM-DD format (defaults to today)")
	flag.Parse()

	// Default to today if no date provided
	if date == "" {
		date = timeutil.FormatDate(time.Now())
	}

	// Validate date format
	if _, err := timeutil.ParseDate(date); err != nil {
		log.Fatalf("Invalid date format: %v", err)
	}

	log.Printf("Fetching Tokyo demand data for %s...", date)

	// For MVP: use sample CSV from testdata
	// TODO: Replace with actual HTTP fetch from TEPCO website
	csvPath := filepath.Join("..", "backend", "internal", "adapters", "testdata", "tepco-sample.csv")
	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Failed to open sample CSV: %v", err)
	}
	defer f.Close()

	// Parse CSV using TEPCO adapter
	adapter := adapters.NewTEPCOAdapter()
	resp, err := adapter.ParseCSV(f, date)
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}

	log.Printf("Parsed %d data points", len(resp.Series))

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create output directory
	outputDir := filepath.Join("..", "public", "data", "jp", "tokyo")
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
