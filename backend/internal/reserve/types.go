// Package reserve provides types and business logic for reserve margin data.
// Follows AGENT_TECH_SPEC.md §3.2 API contract.
package reserve

// Status represents the reserve margin status based on thresholds.
type Status string

const (
	StatusStable Status = "stable" // >= 8%
	StatusWatch  Status = "watch"  // 5-8%
	StatusTight  Status = "tight"  // < 5%
)

// AreaReserve represents reserve margin data for a specific area.
type AreaReserve struct {
	Area              string  `json:"area"`                 // e.g., "tokyo", "kansai"
	ReserveMarginPct  float64 `json:"reserve_margin_pct"`   // Reserve margin percentage
	Status            Status  `json:"status"`               // Derived status (stable/watch/tight)
}

// Source contains attribution for the data source.
type Source struct {
	Name string `json:"name"` // e.g., "OCCTO"
	URL  string `json:"url"`  // Original data source URL
}

// Meta contains optional metadata and warnings.
type Meta struct {
	Warning string `json:"warning,omitempty"` // Non-blocking warning message
}

// Response is the top-level response structure for reserve margin endpoint.
// GET /api/jp/system/reserve?date=YYYY-MM-DD
type Response struct {
	Date   string         `json:"date"`            // YYYY-MM-DD format
	Areas  []AreaReserve  `json:"areas"`           // Reserve data for each area
	Source Source         `json:"source"`          // Data attribution
	Meta   *Meta          `json:"meta,omitempty"`  // Optional metadata/warnings
}

// NewResponse creates a properly initialized Response with defaults.
func NewResponse(date string) *Response {
	return &Response{
		Date:  date,
		Areas: make([]AreaReserve, 0, 2), // Typically Tokyo + Kansai
	}
}

// DeriveStatus calculates the status from reserve margin percentage.
// Thresholds per AGENT_TECH_SPEC.md §3.2:
//   >= 8%   → stable
//   5-8%    → watch
//   < 5%    → tight
func DeriveStatus(reserveMarginPct float64) Status {
	switch {
	case reserveMarginPct >= 8.0:
		return StatusStable
	case reserveMarginPct >= 5.0:
		return StatusWatch
	default:
		return StatusTight
	}
}
