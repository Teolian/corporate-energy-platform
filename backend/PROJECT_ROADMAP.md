# Project Roadmap — Japan Energy Dashboard

**Branch:** `feat/japan-dashboard` (main development branch)
**GitHub:** https://github.com/Teolian/corporate-energy-platform
**Live Demo:** [Preview URL from feat/japan-dashboard deployment]

---

## Overview

Japanese Energy Dashboard aligned with **Digital Grid** services:
- **DGP (Digital Grid Platform)** — demand, reserve, prices, cost calculation
- **Econohashi** — certificates calendar
- **RE Bridge** — PPA digest aggregator
- **Battery Aggregation** — arbitrage simulator

---

## Current Status ✅

### Backend (Go 1.23, Mode A: Static Pipeline)
- ✅ **Phase 1** — Tokyo demand adapter (TEPCO CSV → JSON)
- ✅ **Phase 2** — Kansai adapter + HTTP fetcher with retry
- ✅ Golden tests for both adapters
- ✅ CLI: `go run cmd/fetch-demand/main.go -area tokyo|kansai -date YYYY-MM-DD`
- ✅ JSON artifacts: `public/data/jp/{tokyo,kansai}/demand-YYYY-MM-DD.json`
- ✅ Time utilities (Asia/Tokyo, no DST)
- ✅ HTTP fetcher with exponential backoff

### Frontend (Vue 3.5 + TypeScript + TailwindCSS)
- ✅ **Phase 3A** — Dashboard MVP
  - Tokyo + Kansai charts side-by-side
  - Dark mode with localStorage persistence
  - Date navigation (prev/next)
  - Metrics: Peak MW, Average MW, Forecast Accuracy
  - Dual-mode data client (MOCK/LIVE with graceful fallback)
  - Chart.js visualizations (actual vs forecast)
  - Pinia stores (demand store)

**Dev Server:** `http://localhost:5174/`

---

## Priority 1 — DGP Core Features (1-2 weeks)

### Backend Tasks

#### 1.1 OCCTO Reserve Margin Adapter
**Goal:** Implement AGENT_TECH_SPEC §3.2

**Tasks:**
- [ ] Create `internal/adapters/occto.go`
- [ ] Parse OCCTO HTML/CSV → normalized JSON
- [ ] Implement status mapping:
  - `>= 8%` → "stable"
  - `5-8%` → "watch"
  - `< 5%` → "tight"
- [ ] Golden tests with sample data
- [ ] CLI flag: `--provider occto --date YYYY-MM-DD`
- [ ] Generate: `public/data/jp/system/reserve-YYYY-MM-DD.json`

**Acceptance:**
- JSON validates against schema
- Thresholds correctly mapped
- Missing data → `meta.warning` + 200 status

#### 1.2 JEPX Spot Prices Adapter
**Goal:** Implement AGENT_TECH_SPEC §3.3

**Tasks:**
- [ ] Create `internal/adapters/jepx.go`
- [ ] Parse JEPX historical day-ahead prices
- [ ] Normalize to 24 hourly JPY/kWh slots
- [ ] Golden tests
- [ ] CLI: `--provider jepx --area tokyo --date YYYY-MM-DD`
- [ ] Generate: `public/data/jp/jepx/spot-{area}-YYYY-MM-DD.json`

**Notes:**
- Check JEPX data availability (may require registration)
- Implement fallback to MOCK if unavailable

**Acceptance:**
- 24 timestamps in Asia/Tokyo timezone
- Source attribution present
- Graceful handling of gaps

#### 1.3 Real HTTP Fetching Integration
**Goal:** Replace testdata with live HTTP requests

**Tasks:**
- [ ] Integrate `pkg/http/fetcher.go` into adapters
- [ ] Add environment variables for source URLs:
  - `TEPCO_URL`
  - `KANSAI_URL`
  - `OCCTO_URL`
  - `JEPX_URL`
- [ ] Implement circuit breaker for repeated failures
- [ ] Add structured logging (JSON format):
  - `source`, `duration_ms`, `status`, `artifact`
- [ ] Test with real sources

**Acceptance:**
- Real HTTP requests work in CLI
- Graceful degradation on source unavailability
- Polite rate limiting (no hot-loop)

#### 1.4 Settlement-lite Compute
**Goal:** Implement AGENT_TECH_SPEC §3.5

**Tasks:**
- [ ] Create endpoint or CLI: `POST /api/settlements/run`
- [ ] Input validation:
  - `profile[]` (hourly kWh)
  - `prices` (area + date for JEPX fetch)
  - `pv_offset_pct` (optional)
- [ ] Computation: `cost = Σ(kWh_slot × price_slot × (1 - pv%))`
- [ ] Deterministic rounding: 0.1 JPY, 0.1 kWh
- [ ] Golden test with flat profile × flat price
- [ ] Output: totals + by-hour breakdown

**Acceptance:**
- Golden test matches manual calculation
- Error messages non-blocking
- Assumptions clearly stated in response

---

### Frontend Tasks

#### 1.5 Reserve Badge/Gauge
**Goal:** Display OCCTO reserve status on dashboard

**Tasks:**
- [ ] Create store: `stores/reserve.ts`
  - `getReserve(date)` via `dataClient`
  - Computed: `status`, `color`, `icon`
- [ ] Create components:
  - `ReserveBadge.vue` (inline badge with status)
  - `ReserveGauge.vue` (optional gauge/bar widget)
- [ ] Wire to Dashboard:
  - Badge under each area chart
  - Tooltip with threshold explanation
- [ ] Handle empty/partial data with neutral badge

**Acceptance:**
- Badge visible for selected date
- Thresholds correctly color-coded (green/yellow/red)
- Source attribution (OCCTO) in footer

#### 1.6 Price Overlay
**Goal:** Add JEPX spot prices to demand charts

**Tasks:**
- [ ] Create store: `stores/prices.ts`
  - `getSpotPrice(area, date)` via `dataClient`
  - Computed: chart-ready price series
- [ ] Extend `DemandChart.vue`:
  - Add second Y-axis (JPY/kWh, right side)
  - Toggle button: "Show price overlay"
  - Legend shows units and source (JEPX)
  - Handle gaps → dashed segments + warning icon
- [ ] Create route `/prices`:
  - Dedicated `SpotPriceChart.vue`
  - Price-only and price+demand modes

**Acceptance:**
- 24 hourly points aligned to demand
- Overlay toggle works smoothly
- Gaps handled gracefully with visual feedback

#### 1.7 Cost Card
**Goal:** Display settlement-lite calculation

**Tasks:**
- [ ] Create store: `stores/settlement.ts`
  - `runSettlement(profile, area, date, pvOffsetPct?)`
  - Computed: totals, by-hour breakdown
- [ ] Create components:
  - `CostCard.vue` (totals: kWh + cost JPY)
  - `CostBreakdownTable.vue` (optional per-hour table)
- [ ] Wire to Dashboard:
  - Card in sidebar or below charts
  - Show assumptions (area, date, PV%)
  - CSV export button (see P2)

**Acceptance:**
- Totals match backend result for test payload
- Assumptions + source attribution visible
- Non-blocking warnings for partial data

#### 1.8 Feature Flags Setup
**Goal:** Implement feature toggles for P1 widgets

**Tasks:**
- [ ] Create `.env.example` with flags:
  ```
  VITE_DATA_MODE=mock
  VITE_FEAT_RESERVE=false
  VITE_FEAT_JEPX=false
  VITE_FEAT_SETTLEMENT=false
  VITE_FEAT_COMPARISON=false
  ```
- [ ] Implement flag logic in `main.ts` or composable
- [ ] Conditional rendering:
  - Reserve badge shown only if `VITE_FEAT_RESERVE=true`
  - Price overlay if `VITE_FEAT_JEPX=true`
  - Cost card if `VITE_FEAT_SETTLEMENT=true`
- [ ] Default: flags OFF on `main`, ON on preview

**Acceptance:**
- Flags control widget visibility
- Easy to flip via env variable
- No UI breaking when flags off

---

### P1 Acceptance Criteria
- ✅ Dashboard shows: demand + **reserve badge** + **price overlay** + **cost card** in LIVE mode
- ✅ All widgets handle empty/partial data gracefully
- ✅ Source attribution visible (TEPCO, Kansai, JEPX, OCCTO)
- ✅ Golden tests pass for all new adapters
- ✅ Schema validation in CI
- ✅ Preview deployment renders without regressions

---

## Priority 2 — Enhanced UX (1 week)

### Backend Tasks

#### 2.1 GitHub Actions Pipeline
**Goal:** Automate JSON generation

**Tasks:**
- [ ] Create `.github/workflows/fetch-data.yml`
- [ ] Cron schedule: hourly or daily
- [ ] Run CLI for all providers:
  - Tokyo demand
  - Kansai demand
  - OCCTO reserve
  - JEPX prices (both areas)
- [ ] Commit updated JSON to `public/data/`
- [ ] Add schema validation step
- [ ] Manual dispatch option for backfill

**Acceptance:**
- Pipeline runs on schedule
- Artifacts committed automatically
- Build fails if schema invalid

---

### Frontend Tasks

#### 2.2 Comparison View
**Goal:** Tokyo vs Kansai side-by-side overlay

**Tasks:**
- [ ] Create route `/comparison`
- [ ] Create component: `ComparisonChart.vue`
  - Overlay both regions on one chart
  - Shared X-axis (24 hours)
  - Option: shared Y-axis or dual Y-axis
  - Toggle for demand-only vs demand+price
- [ ] Date navigation synced for both areas

**Acceptance:**
- Renders live data for both areas
- Date nav updates both series
- Legend clearly distinguishes regions

#### 2.3 Dedicated Price View
**Goal:** Full-page JEPX spot price visualization

**Tasks:**
- [ ] Enhance route `/prices`
- [ ] Modes:
  - Price-only (line chart JPY/kWh)
  - Price+demand correlation (dual-axis)
- [ ] Area selector (Tokyo/Kansai)
- [ ] Date range selector (optional)

**Acceptance:**
- Shows 24 hourly prices clearly
- Correlation mode useful for analysis

#### 2.4 Export Features
**Goal:** CSV and PNG downloads

**Tasks:**
- [ ] Create composable: `useExport.ts`
- [ ] Functions:
  - `exportDemandCSV(data, area, date)` → download CSV
  - `exportChartPNG(chartRef)` → Chart.js `.toBase64Image()`
- [ ] Add buttons to Dashboard:
  - "Export CSV" near date picker
  - "Export PNG" on each chart
- [ ] Filenames: `demand-{area}-{date}.{csv|png}`

**Acceptance:**
- CSV contains all 24 hourly points
- PNG captures current chart state
- Files download with sensible names

#### 2.5 i18n Scaffolding
**Goal:** Prepare for English/Japanese localization

**Tasks:**
- [ ] Create `src/i18n/messages.ts`
- [ ] Define keys:
  - `reserve.status.{stable|watch|tight}`
  - `labels.{price_yen_per_kwh, cost_yen, kwh, mw}`
  - `tooltips.{gaps_detected, reserve_explain, forecast_unavailable}`
  - `buttons.{export_csv, export_png, toggle_price}`
- [ ] Stub translations (EN + JA)
- [ ] Wire to at least one component (proof of concept)

**Acceptance:**
- i18n structure ready
- Easy to add more translations
- No hardcoded strings in critical UI

---

### P2 Acceptance Criteria
- ✅ `/prices` and `/comparison` routes functional
- ✅ CSV/PNG export works for demand data
- ✅ i18n scaffolding in place (EN + JA stub)
- ✅ GitHub Actions pipeline running (optional for P2)

---

## Priority 3 — DG-Adjacent Features (Optional, 2-3 weeks)

### 3.1 Portfolio Mix Simulator (DGP)
**Goal:** Compare fixed vs spot vs RE contract mix

**Tasks:**
- [ ] Backend: `POST /api/portfolio/simulate`
  - Input: mix percentages (fixed%, spot%, RE%)
  - Compute weighted average cost
  - Scenario comparison (Base vs Alt)
- [ ] Frontend: `/portfolio` route
  - Input sliders for mix %
  - Cost breakdown visualization
  - Savings calculator

**Business Value:** Aligns with DGP service (fixed/spot mix optimization)

---

### 3.2 Certificates Calendar (Econohashi)
**Goal:** Display renewable certificate schedule

**Tasks:**
- [ ] Backend: `GET /api/certs/calendar?year=YYYY`
  - Read-only calendar view
  - Public data aggregation
- [ ] Frontend: `/certificates` route
  - Calendar grid (monthly view)
  - Tooltips with certificate details
  - Explainer: what are RE certificates

**Business Value:** Aligns with Econohashi (certificate tracking)

---

### 3.3 PPA Digest (RE Bridge)
**Goal:** Aggregate PPA announcements from public sources

**Tasks:**
- [ ] Backend: `GET /api/rebridge/digest`
  - Scrape/aggregate press releases
  - Metrics: total MW, avg price, top companies
  - Update frequency: weekly
- [ ] Frontend: `/ppa` route
  - Dashboard with key metrics
  - List of recent announcements
  - Link to sources

**Business Value:** Aligns with RE Bridge (PPA tracking)

---

### 3.4 Battery Arbitrage Simulator (Battery Aggregation)
**Goal:** Simple charge/discharge optimization vs spot prices

**Tasks:**
- [ ] Backend: `POST /api/battery/sim`
  - Input: battery capacity (kWh), max charge rate (kW)
  - Logic: charge at low price, discharge at high price
  - Output: potential profit, optimal schedule
- [ ] Frontend: `/battery` route
  - Input fields for battery specs
  - Visualization: charge schedule vs price
  - ROI calculator

**Business Value:** Aligns with Battery Aggregation service

---

### P3 Acceptance Criteria
- ✅ At least 2 DG-adjacent features live
- ✅ Clear value proposition for each feature
- ✅ Production-ready deployment

---

## Technical References

### Canonical Specifications
- **API Contracts:** [AGENT_TECH_SPEC.md](./AGENT_TECH_SPEC.md) §3
- **Business Logic:** [AGENT_BIZ_SPEC.md](./AGENT_BIZ_SPEC.md)
- **Integration Guide:** [AGENT_INTEGRATION_GUIDE.md](./AGENT_INTEGRATION_GUIDE.md)
- **Changelog:** [CHANGELOG.md](./CHANGELOG.md)

### Archived Plans
- Old sprint plans: [docs/archive/](./docs/archive/)

---

## Branching Strategy

```
main                          # Original demo (Vue 3 corporate benchmark, stable, DON'T TOUCH)
└── feat/japan-dashboard      # Main dev branch (backend + frontend-japan)
    ├── feat/occto-adapter    # Feature: OCCTO reserve margin
    ├── feat/jepx-adapter     # Feature: JEPX spot prices
    ├── feat/settlement       # Feature: Settlement-lite compute
    └── feat/export-ui        # Feature: CSV/PNG export
```

**Rules:**
- `main` — stable old demo, **never modify**
- `feat/japan-dashboard` — primary development branch
- Feature branches: create from `feat/japan-dashboard`, merge back via PR
- Preview deployments: test on `feat/japan-dashboard` URL with LIVE mode

---

## Testing Strategy

### Backend
- **Unit tests:** Adapters with golden CSV/JSON samples
- **Integration tests:** HTTP fetcher with mock server
- **Schema validation:** JSON artifacts vs `/pipeline/schemas/*.json`
- **Edge cases:** Empty data, missing forecast, date boundaries
- **Golden tests:** Deterministic output for same input

### Frontend
- **Unit:** Composables, stores (Vitest)
- **Integration:** Data client with mocked API
- **E2E (Playwright):** Smoke tests
  - Load date → render charts
  - Toggle dark mode → verify styles
  - Switch MOCK/LIVE modes → no crash
  - Toggle price overlay
  - Export CSV/PNG → files download

---

## Deployment

### Current (Development)
- **Frontend:** `npm run dev` → `http://localhost:5174/`
- **Backend:** CLI jobs generate JSON artifacts locally

### Future (Production)
- **Frontend:** Vercel deployment from `feat/japan-dashboard`
  - Preview: every PR/branch
  - Production: merge to `main` or dedicated live branch
- **Backend:** Mode B (Serverless API) on Fly.io/Railway (optional)
  - Alternative: Keep Mode A (static) for zero cost
- **Preview:** Every PR generates preview URL with `VITE_DATA_MODE=live`

---

## Success Metrics

### Technical
- ✅ All golden tests pass
- ✅ Schema validation in CI (zero drift)
- ✅ E2E smoke tests green
- ✅ Preview deployments stable

### User Experience
- ✅ Dashboard loads in < 2s (LIVE mode)
- ✅ Zero crashes on empty/partial data
- ✅ Dark mode seamless
- ✅ Mobile responsive (TailwindCSS)

### Business Value
- ✅ Demonstrates Digital Grid service understanding
- ✅ Portfolio-ready (clean code, good docs)
- ✅ Real data (not mocks) in production

---

## Non-Goals (MVP)

- ❌ User authentication/authorization
- ❌ Multi-tenant orgs and billing
- ❌ Advanced ML forecasting
- ❌ Full retail tariff modeling
- ❌ PDF report generation
- ❌ World data (ENTSO-E, EIA) — future Phase 4

---

## FAQ

**Q: Why two frontends (old demo + japan dashboard)?**
A: The old demo (`frontend/`) showcases corporate energy benchmarking with mock data. The japan dashboard (`frontend-japan/`) demonstrates live data integration and DG-aligned features. Both serve different portfolio purposes.

**Q: Why keep `main` branch untouched?**
A: `main` is a stable demo we can show anytime. Development happens in `feat/japan-dashboard` with preview deployments. We merge to `main` only after full validation.

**Q: What if JEPX data requires login?**
A: Implement graceful fallback to MOCK for that widget. Add `meta.warning` and document in README.

**Q: Mode A vs Mode B?**
A: Mode A (static JSON) is cheaper and simpler. Mode B (serverless API) adds caching and runtime flexibility. Start with A, add B later if needed.

**Q: Where are Digital Grid services?**
A: https://digitalgrid.jp/ — Study their DGP, Econohashi, RE Bridge, Battery Aggregation services to align features.

---

## Next Steps

1. **Immediate:** Start P1 Backend (OCCTO adapter)
2. **Parallel:** Implement P1 Frontend (Reserve badge)
3. **Week 2:** JEPX adapter + Price overlay
4. **Week 2-3:** Settlement-lite + Cost card
5. **Week 3:** Feature flags + P1 complete
6. **Week 4+:** P2 UX enhancements
7. **Optional:** P3 DG-adjacent features

---

**Last Updated:** 2025-10-24
**Current Phase:** P1 Planning
**Branch:** feat/japan-dashboard
**Contributors:** Teo + Claude Code
