# Backend — Corporate Energy Platform

Backend services for the Corporate Energy Benchmark Platform. Implements data adapters, normalization, and API endpoints for Japanese power system data.

## Architecture

- **Go 1.24+** — Core backend services
- **Mode A (Static Pipeline)** — Current implementation: GitHub Actions → normalized JSON artifacts
- **Mode B (Serverless API)** — Future: Runtime API endpoints with caching

Following specifications:
- [AGENT_TECH_SPEC.md](./AGENT_TECH_SPEC.md) — Technical contracts
- [AGENT_BIZ_SPEC.md](./AGENT_BIZ_SPEC.md) — Business logic
- [AGENT_INTEGRATION_GUIDE.md](./AGENT_INTEGRATION_GUIDE.md) — Frontend integration

## Project Structure

```
backend/
├── cmd/
│   ├── api/              # Future: HTTP API server (chi/gin)
│   └── fetch-demand/     # Pipeline: fetch & normalize demand data
├── internal/
│   ├── demand/           # Demand data types & business logic
│   ├── adapters/         # Source adapters (TEPCO, Kansai, etc.)
│   ├── reserve/          # Reserve margin (future)
│   ├── jepx/             # JEPX price data (future)
│   ├── weather/          # Weather/solar data (future)
│   └── settlement/       # Cost calculation (future)
├── pkg/
│   ├── timeutil/         # Asia/Tokyo timezone utilities
│   └── units/            # Unit conversions (MW ↔ kWh)
└── go.mod
```

## Phase 1 Status — Tokyo Demand Adapter ✅

**Implemented:**
- ✅ TEPCO CSV adapter with auto-header detection
- ✅ Normalization to AGENT_TECH_SPEC §3.1 format
- ✅ Unit tests with golden test data
- ✅ Pipeline job: `cmd/fetch-demand/main.go`
- ✅ Sample output: `/public/data/jp/tokyo/demand-2025-10-24.json`
- ✅ JSON Schema validation: `/pipeline/schemas/demand-schema.json`

**Not yet implemented:**
- ❌ Real HTTP fetch from TEPCO website (using testdata for now)
- ❌ Kansai adapter
- ❌ Reserve, JEPX, Weather adapters
- ❌ API server endpoints

## Quick Start

### Prerequisites

- Go 1.24+ installed
- Git repository initialized

### Run Tests

```bash
cd backend
go test ./...
```

### Generate Tokyo Demand JSON

```bash
cd backend
go run cmd/fetch-demand/main.go -date 2025-10-24
```

Output: `../public/data/jp/tokyo/demand-2025-10-24.json`

### Validate Output Against Schema

```bash
# Install JSON schema validator (optional)
npm install -g ajv-cli

# Validate
ajv validate -s ../pipeline/schemas/demand-schema.json \
             -d ../public/data/jp/tokyo/demand-2025-10-24.json
```

## Data Sources

### Tokyo (TEPCO)
- **Source:** Tokyo Electric Power Company
- **URL:** https://www.tepco.co.jp/forecast/html/download-j.html
- **Format:** CSV with hourly demand (actual + forecast)
- **Units:** 万kW (10,000 kW) → converted to MW
- **Timezone:** Asia/Tokyo (UTC+09:00, no DST)

### Future Sources
- Kansai Electric Power Company
- OCCTO (reserve margin)
- JEPX (spot prices)
- Open-Meteo (weather/solar)

## API Contracts (Future)

Following AGENT_TECH_SPEC §3:

### Demand
```
GET /api/jp/{area}/demand?date=YYYY-MM-DD
area ∈ {tokyo, kansai}
```

Response:
```json
{
  "area": "tokyo",
  "date": "2025-10-24",
  "timezone": "Asia/Tokyo",
  "timescale": "hourly",
  "series": [...],
  "source": {"name": "TEPCO", "url": "..."}
}
```

## Development Workflow

### Adding a New Adapter

1. Create type definitions in `internal/{domain}/types.go`
2. Implement adapter in `internal/adapters/{source}.go`
3. Write tests with golden test data in `testdata/`
4. Create pipeline job in `cmd/fetch-{domain}/`
5. Add JSON schema in `pipeline/schemas/`
6. Update this README

### Testing Strategy

- **Unit tests:** Pure adapter functions with golden CSV/JSON
- **Golden tests:** Pinned sample files in `testdata/`
- **Schema validation:** All outputs must pass JSON Schema checks
- **Edge cases:** Empty data, missing forecast, date boundaries

### Commit Convention

```bash
# Feature
git commit -m "feat(adapters): add Kansai demand adapter"

# Fix
git commit -m "fix(tepco): handle missing forecast gracefully"

# Test
git commit -m "test(adapters): add edge cases for TEPCO parser"
```

## Phase 2 Roadmap

- [ ] Real HTTP fetching with retries
- [ ] Kansai demand adapter
- [ ] OCCTO reserve margin adapter
- [ ] JEPX spot price adapter
- [ ] Open-Meteo weather/solar adapter
- [ ] GitHub Actions cron pipeline
- [ ] API server with chi router (Mode B)
- [ ] Redis caching layer
- [ ] Settlement-lite compute endpoint

## Troubleshooting

### Tests failing
```bash
cd backend
go mod tidy
go test -v ./...
```

### Import errors
Make sure you're in the backend directory when running commands.

### Timezone issues
All timestamps use `Asia/Tokyo` (UTC+09:00). Japan has no DST.

## Contributing

1. Create feature branch from `feat/backend-core`
2. Follow AGENT_TECH_SPEC contracts exactly
3. Add tests for all new code
4. Run `go test ./...` before committing
5. Open PR with clear description

## License

Private portfolio project.
