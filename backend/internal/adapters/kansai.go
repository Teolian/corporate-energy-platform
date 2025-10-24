package adapters

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/teo/aversome/backend/internal/demand"
	"github.com/teo/aversome/backend/pkg/timeutil"
)

// KansaiAdapter normalizes Kansai Electric Power Company demand data.
type KansaiAdapter struct {
	sourceURL string
}

// NewKansaiAdapter creates a new Kansai data adapter.
func NewKansaiAdapter() *KansaiAdapter {
	return &KansaiAdapter{
		sourceURL: "https://www.kansai-td.co.jp/denkiyoho/",
	}
}

// ParseCSV parses Kansai CSV data into normalized demand.Response.
// CSV format is similar to TEPCO:
//   DATE,TIME,実績(万kW),予測(万kW)
//   2025-10-24,0:00,1234.5,1256.2
//
// Notes:
// - Same header detection as TEPCO
// - Units: 万kW (10,000 kW) = 10 MW → convert to MW
// - Forecast may be empty for some rows
func (a *KansaiAdapter) ParseCSV(reader io.Reader, date string) (*demand.Response, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Auto-detect column indices
	colIndices := a.detectColumns(header)
	if colIndices["date"] == -1 || colIndices["time"] == -1 || colIndices["actual"] == -1 {
		return nil, fmt.Errorf("required columns not found in header: %v", header)
	}

	resp := demand.NewResponse(demand.AreaKansai, date)
	resp.Source = demand.Source{
		Name: "Kansai",
		URL:  a.sourceURL,
	}

	// Parse date for validation
	baseDate, err := timeutil.ParseDate(date)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}

	hasForecast := false
	lineNum := 1

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

		// Extract date and time
		rowDate := strings.TrimSpace(record[colIndices["date"]])
		rowTime := strings.TrimSpace(record[colIndices["time"]])

		// Only include rows matching the requested date
		if rowDate != date {
			continue
		}

		// Parse hour from time string
		hour, err := a.parseHour(rowTime)
		if err != nil {
			return nil, fmt.Errorf("invalid time format at line %d: %s", lineNum, rowTime)
		}

		// Create timestamp
		ts := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), hour, 0, 0, 0, timeutil.TokyoLocation)

		// Parse actual demand (万kW → MW)
		actualStr := strings.TrimSpace(record[colIndices["actual"]])
		actual, err := strconv.ParseFloat(actualStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid actual value at line %d: %s", lineNum, actualStr)
		}
		demandMW := actual * 10.0 // 万kW to MW

		// Parse forecast if available
		var forecastMW *float64
		if colIndices["forecast"] != -1 && colIndices["forecast"] < len(record) {
			forecastStr := strings.TrimSpace(record[colIndices["forecast"]])
			if forecastStr != "" {
				forecast, err := strconv.ParseFloat(forecastStr, 64)
				if err == nil {
					fVal := forecast * 10.0
					forecastMW = &fVal
					hasForecast = true
				}
			}
		}

		point := demand.SeriesPoint{
			Timestamp:  ts,
			DemandMW:   demandMW,
			ForecastMW: forecastMW,
		}
		resp.Series = append(resp.Series, point)
	}

	// Validate we have data
	if len(resp.Series) == 0 {
		return nil, fmt.Errorf("no data found for date %s", date)
	}

	// Add warning if no forecast data
	if !hasForecast {
		resp.Meta = &demand.Meta{
			Warning: "Forecast data not available for this date",
		}
	}

	return resp, nil
}

// detectColumns finds column indices by header names.
func (a *KansaiAdapter) detectColumns(header []string) map[string]int {
	indices := map[string]int{
		"date":     -1,
		"time":     -1,
		"actual":   -1,
		"forecast": -1,
	}

	for i, col := range header {
		col = strings.ToLower(strings.TrimSpace(col))
		switch {
		case strings.Contains(col, "date") || col == "日付":
			indices["date"] = i
		case strings.Contains(col, "time") || col == "時刻":
			indices["time"] = i
		case strings.Contains(col, "実績") || strings.Contains(col, "actual"):
			indices["actual"] = i
		case strings.Contains(col, "予測") || strings.Contains(col, "forecast"):
			indices["forecast"] = i
		}
	}

	return indices
}

// parseHour extracts hour from time string.
func (a *KansaiAdapter) parseHour(timeStr string) (int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return 0, fmt.Errorf("invalid hour: %s", timeStr)
	}
	return hour, nil
}
