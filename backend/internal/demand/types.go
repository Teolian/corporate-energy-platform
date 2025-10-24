// Package demand provides types and business logic for electricity demand data.
// Follows AGENT_TECH_SPEC.md §3.1 API contract.
package demand

import "time"

// Area represents a geographic region for demand data.
type Area string

const (
	AreaTokyo  Area = "tokyo"
	AreaKansai Area = "kansai"
)

// Timescale represents the granularity of data points.
type Timescale string

const (
	TimescaleHourly Timescale = "hourly"
)

// Source contains attribution for the data source.
type Source struct {
	Name string `json:"name"` // e.g., "TEPCO", "Kansai"
	URL  string `json:"url"`  // Original data source URL
}

// SeriesPoint represents a single time-series data point.
type SeriesPoint struct {
	Timestamp  time.Time `json:"ts"`          // ISO8601 with +09:00 offset
	DemandMW   float64   `json:"demand_mw"`   // Actual demand in megawatts
	ForecastMW *float64  `json:"forecast_mw"` // Forecasted demand (optional)
}

// Meta contains optional metadata and warnings.
type Meta struct {
	Warning string `json:"warning,omitempty"` // Non-blocking warning message
}

// Response is the top-level response structure for demand endpoints.
// GET /api/jp/{area}/demand?date=YYYY-MM-DD
type Response struct {
	Area      Area          `json:"area"`      // Geographic area
	Date      string        `json:"date"`      // YYYY-MM-DD format
	Timezone  string        `json:"timezone"`  // Always "Asia/Tokyo" for Japan
	Timescale Timescale     `json:"timescale"` // Data granularity
	Series    []SeriesPoint `json:"series"`    // Hourly data points (up to 24)
	Source    Source        `json:"source"`    // Data attribution
	Meta      *Meta         `json:"meta,omitempty"`
}

// NewResponse creates a properly initialized Response with defaults.
func NewResponse(area Area, date string) *Response {
	return &Response{
		Area:      area,
		Date:      date,
		Timezone:  "Asia/Tokyo",
		Timescale: TimescaleHourly,
		Series:    make([]SeriesPoint, 0, 24),
	}
}
