// Package adapters provides data source adapters for normalization.
package adapters

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/teo/aversome/backend/internal/reserve"
)

// OCCTOAdapter normalizes OCCTO (Organization for Cross-regional Coordination of
// Transmission Operators) reserve margin data.
type OCCTOAdapter struct {
	sourceURL string
}

// NewOCCTOAdapter creates a new OCCTO data adapter.
func NewOCCTOAdapter() *OCCTOAdapter {
	return &OCCTOAdapter{
		sourceURL: "https://www.occto.or.jp/",
	}
}

// ParseCSV parses OCCTO CSV data into normalized reserve.Response.
// CSV format (expected):
//   日付,エリア,予備率
//   2025-10-24,東京,8.5
//   2025-10-24,関西,10.2
//
// Or English headers:
//   DATE,AREA,RESERVE_MARGIN
//   2025-10-24,tokyo,8.5
//   2025-10-24,kansai,10.2
//
// Notes:
// - Header detection by column names (order may vary)
// - Reserve margin is in percentage (%)
// - Multiple areas per date
func (a *OCCTOAdapter) ParseCSV(reader io.Reader, date string) (*reserve.Response, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Auto-detect column indices
	colIndices := a.detectColumns(header)
	if colIndices["date"] == -1 || colIndices["area"] == -1 || colIndices["reserve_margin"] == -1 {
		return nil, fmt.Errorf("required columns not found in header: %v", header)
	}

	resp := reserve.NewResponse(date)
	resp.Source = reserve.Source{
		Name: "OCCTO",
		URL:  a.sourceURL,
	}

	lineNum := 1
	areasFound := make(map[string]bool)

	// Read data rows
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV line %d: %w", lineNum, err)
		}
		lineNum++

		// Extract date
		rowDate := strings.TrimSpace(record[colIndices["date"]])

		// Only include rows matching the requested date
		if rowDate != date {
			continue
		}

		// Extract area
		areaStr := strings.TrimSpace(record[colIndices["area"]])
		area := a.normalizeArea(areaStr)
		if area == "" {
			return nil, fmt.Errorf("invalid area at line %d: %s", lineNum, areaStr)
		}

		// Parse reserve margin percentage
		reserveStr := strings.TrimSpace(record[colIndices["reserve_margin"]])
		reservePct, err := strconv.ParseFloat(reserveStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid reserve margin at line %d: %s", lineNum, reserveStr)
		}

		// Derive status from percentage
		status := reserve.DeriveStatus(reservePct)

		areaReserve := reserve.AreaReserve{
			Area:             area,
			ReserveMarginPct: reservePct,
			Status:           status,
		}
		resp.Areas = append(resp.Areas, areaReserve)
		areasFound[area] = true
	}

	// Validate we have data
	if len(resp.Areas) == 0 {
		return nil, fmt.Errorf("no data found for date %s", date)
	}

	// Add warning if missing expected areas
	if !areasFound["tokyo"] || !areasFound["kansai"] {
		resp.Meta = &reserve.Meta{
			Warning: "Data for some areas not available for this date",
		}
	}

	return resp, nil
}

// detectColumns finds column indices by header names.
// Returns map with keys: date, area, reserve_margin.
func (a *OCCTOAdapter) detectColumns(header []string) map[string]int {
	indices := map[string]int{
		"date":           -1,
		"area":           -1,
		"reserve_margin": -1,
	}

	for i, col := range header {
		col = strings.ToLower(strings.TrimSpace(col))
		switch {
		case strings.Contains(col, "date") || col == "日付":
			indices["date"] = i
		case strings.Contains(col, "area") || col == "エリア" || col == "地域":
			indices["area"] = i
		case strings.Contains(col, "reserve") || strings.Contains(col, "margin") || col == "予備率":
			indices["reserve_margin"] = i
		}
	}

	return indices
}

// normalizeArea converts various area representations to canonical form.
// Handles both English and Japanese names.
func (a *OCCTOAdapter) normalizeArea(areaStr string) string {
	areaStr = strings.ToLower(strings.TrimSpace(areaStr))

	switch {
	case areaStr == "tokyo" || areaStr == "東京" || strings.Contains(areaStr, "tokyo"):
		return "tokyo"
	case areaStr == "kansai" || areaStr == "関西" || strings.Contains(areaStr, "kansai"):
		return "kansai"
	default:
		// Return as-is for other areas (future expansion)
		return areaStr
	}
}
