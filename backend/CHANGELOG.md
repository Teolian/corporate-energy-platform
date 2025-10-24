# Changelog

All notable changes to the backend will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Phase 1 — Tokyo Demand Adapter MVP (2025-10-24)

#### Added
- **Project initialization**
  - Go module setup (`go.mod`)
  - Directory structure following AGENT_TECH_SPEC §2
  - Git repository with `feat/backend-core` branch

- **Core types** (`internal/demand/types.go`)
  - `Response` structure matching AGENT_TECH_SPEC §3.1
  - `Area`, `Timescale`, `Source`, `SeriesPoint`, `Meta` types
  - Type safety for demand data contracts

- **Timezone utilities** (`pkg/timeutil/jp.go`)
  - Asia/Tokyo timezone handling (UTC+09:00, no DST)
  - Date parsing/formatting helpers
  - Hourly slot generation
  - ISO8601 formatting with +09:00 offset

- **TEPCO adapter** (`internal/adapters/tepco.go`)
  - CSV parser with auto-header detection
  - Unit conversion: 万kW (10,000 kW) → MW
  - Graceful handling of missing forecast data
  - Source attribution (name + URL)
  - Deterministic normalization

- **Tests** (`internal/adapters/tepco_test.go`)
  - Golden test with sample TEPCO CSV
  - Unit tests for hour parsing
  - Validation of all response fields
  - Edge case coverage

- **Pipeline job** (`cmd/fetch-demand/main.go`)
  - CLI tool for generating JSON artifacts
  - Date parameter support (`-date YYYY-MM-DD`)
  - Output to `/public/data/jp/tokyo/demand-YYYY-MM-DD.json`
  - Logging and error handling

- **JSON Schema** (`pipeline/schemas/demand-schema.json`)
  - Validation schema for demand responses
  - JSON Schema Draft 07 format
  - Field constraints and patterns

- **Documentation**
  - `backend/README.md` — Quick start, structure, workflow
  - `.env.example` — Configuration template
  - `.gitignore` — Backend-specific ignores
  - This CHANGELOG

#### Sample Output
- `public/data/jp/tokyo/demand-2025-10-24.json` (2.8 KB)
  - 24 hourly data points
  - Actual + forecast values
  - Full source attribution

#### Technical Decisions
1. **Go 1.24** — Latest stable version
2. **Pure functions** — Adapters are deterministic, testable
3. **Auto-header detection** — Resilient to CSV format changes
4. **Mode A (Static)** — Simplest MVP, zero runtime cost
5. **Internal packages** — Proper Go module boundaries

#### Testing
- ✅ All tests passing
- ✅ Sample JSON validates against schema
- ✅ Code follows AGENT_TECH_SPEC contracts

#### Next Steps (Phase 2)
- [ ] Real HTTP fetch from TEPCO website
- [ ] Kansai demand adapter
- [ ] OCCTO reserve margin adapter
- [ ] GitHub Actions cron pipeline
- [ ] Error handling improvements

---

## How to Read This Changelog

- **Added** — New features
- **Changed** — Changes to existing functionality
- **Deprecated** — Features marked for removal
- **Removed** — Features removed
- **Fixed** — Bug fixes
- **Security** — Security fixes

Each phase corresponds to sections in AGENT_BIZ_SPEC.md.
