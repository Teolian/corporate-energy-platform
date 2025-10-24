// Package main provides a pipeline job to fetch reserve margin data from OCCTO and normalize to JSON.
// Usage: go run main.go -date 2025-10-24
// Output: /public/data/jp/system/reserve-YYYY-MM-DD.json
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

	log.Printf("Fetching OCCTO reserve margin data for %s...", date)

	// Open test CSV (in production, this would be fetched via HTTP)
	csvPath := filepath.Join("internal", "adapters", "testdata", "occto-sample.csv")
	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Failed to open OCCTO CSV: %v", err)
	}
	defer f.Close()

	// Parse CSV using OCCTO adapter
	adapter := adapters.NewOCCTOAdapter()
	resp, err := adapter.ParseCSV(f, date)
	if err != nil {
		log.Fatalf("Failed to parse OCCTO CSV: %v", err)
	}

	log.Printf("Parsed %d areas", len(resp.Areas))
	if resp.Meta != nil && resp.Meta.Warning != "" {
		log.Printf("Warning: %s", resp.Meta.Warning)
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Create output directory (system-level data)
	outputDir := filepath.Join("..", "public", "data", "jp", "system")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Write JSON file
	outputPath := filepath.Join(outputDir, fmt.Sprintf("reserve-%s.json", date))
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}

	log.Printf("✓ Successfully wrote %s (%d bytes)", outputPath, len(jsonData))
}
