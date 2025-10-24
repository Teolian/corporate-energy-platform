# AGENT_TECH_SPEC.md
**Project:** Corporate Energy Platform — Backend MVP  
**Audience:** AI coding agent (e.g., Claude/Cursor) + human reviewers  
**Purpose:** Precise technical contract and engineering guardrails. The agent MUST follow this document when generating code, tests, and infra.

---

## 0. How to Use This Spec
- Treat this file as a **System Prompt** for the coding agent.  
- Always cross-check with external docs index **Context7** (latest patterns for Go/FastAPI, sqlc, pydantic v2, ruff, mypy, chi/gin, testing, CI). If Context7 conflicts with this file, **ask to reconcile** and prefer this file unless a security bug or breaking change is identified.
- Never invent fields or change schemas without updating this spec and the companion business file **AGENT_BIZ_SPEC.md**.

---

## 1. Scope (MVP)
Implement a read-only data backend + a simple compute endpoint (“settlement-lite”) for the existing Vercel frontend.

**Functional endpoints (JSON):**
1. `GET /api/jp/tokyo/demand?date=YYYY-MM-DD`
2. `GET /api/jp/kansai/demand?date=YYYY-MM-DD`
3. `GET /api/jp/system/reserve?date=YYYY-MM-DD`
4. `GET /api/jp/jepx/spot?date=YYYY-MM-DD&area=tokyo`
5. `GET /api/weather/solar?lat=..&lon=..&date=YYYY-MM-DD`
6. `POST /api/settlements/run`  (compute cost from profile × JEPX prices)

**Data sources (free):** TEPCO CSV, Kansai CSV, OCCTO reserve (JP), JEPX day-ahead (history), Open‑Meteo (JMA + Solar).

Two delivery modes (choose A first, keep B-compatible):
- **Mode A (Static Pipeline):** GitHub Actions normalize → commit JSON into `/public/data/*` for the frontend.
- **Mode B (Serverless API):** Fast API/Go handlers fetch/cache on demand, return normalized JSON.

---

## 2. Architecture & Directory Layout
Monorepo layout suggestion (adapt if the repo already differs):
```
/backend
  /cmd/api                # Go HTTP entry (chi)
  /internal               # domain packages
    /demand, /reserve, /jepx, /weather, /settlement
    /adapters             # CSV/HTML/JSON adapters (pure, deterministic)
    /store                # persistence/caching abstraction
  /pkg
    /http                 # middlewares (CORS, logging, panic-recover)
    /timeutil, /units     # helpers (TZ, MW↔kWh)
  /migrations             # if Postgres is used (Mode B+)
/pipeline
  /jobs                   # GH Actions scripts to fetch→normalize→emit JSON
  /schemas                # JSON schema files for emitted artifacts
/docs
  AGENT_TECH_SPEC.md
  AGENT_BIZ_SPEC.md
```

**Languages & core libs**
- Go 1.22+, `chi` (router), `pgx` if DB is used, `sqlc` for queries (Mode B).  
- Python 3.11+ for ML/solar/weather glue (optional): FastAPI, pydantic v2, `ruff`, `mypy --strict`.
- Static first: allow running with **no DB** (read JSON artifacts from `/public/data/`).

---

## 3. API Contracts (immutable without spec change)
### 3.1 Demand (Tokyo/Kansai)
**Request:** `GET /api/jp/{area}/demand?date=YYYY-MM-DD`  
`area ∈ {tokyo, kansai}`  
**Response:**
```json
{
  "area": "tokyo",
  "date": "2025-10-24",
  "timezone": "Asia/Tokyo",
  "timescale": "hourly",
  "series": [
    {"ts": "2025-10-24T00:00:00+09:00", "demand_mw": 26654, "forecast_mw": 27010}
  ],
  "source": {"name": "TEPCO", "url": "..."}
}
```
Rules:
- `forecast_mw` MAY be null/omitted if source lacks forecast for the day.
- `ts` must be ISO8601 with offset `+09:00` (no DST in JP).

### 3.2 Reserve Margin (System)
**Request:** `GET /api/jp/system/reserve?date=YYYY-MM-DD`  
**Response:**
```json
{
  "date": "2025-10-24",
  "areas": [
    {"area":"tokyo","reserve_margin_pct":5.1,"status":"tight"},
    {"area":"kansai","reserve_margin_pct":8.9,"status":"stable"}
  ],
  "source":{"name":"OCCTO","url":"..."}
}
```
Status mapping (configurable): `>=8 → "stable"`, `5–8 → "watch"`, `<5 → "tight"`.

### 3.3 JEPX Spot (Day‑Ahead)
**Request:** `GET /api/jp/jepx/spot?date=YYYY-MM-DD&area=tokyo`  
**Response:**
```json
{
  "date": "2025-10-23",
  "area": "tokyo",
  "timescale": "hourly",
  "price_yen_per_kwh": [
    {"ts":"2025-10-23T00:00:00+09:00","price":24.32}
  ],
  "source":{"name":"JEPX","url":"..."}
}
```

### 3.4 Solar/Weather
**Request:** `GET /api/weather/solar?lat=35.68&lon=139.76&date=YYYY-MM-DD`  
**Response:**
```json
{
  "lat":35.68, "lon":139.76, "date":"2025-10-24",
  "series":[
    {"ts":"2025-10-24T00:00:00+09:00","ghi":120.0,"dni":200.0,"dhi":60.0,"cloudcover":0.35}
  ],
  "source":{"name":"Open-Meteo","url":"..."}
}
```

### 3.5 Settlement‑Lite (Compute)
**Request:** `POST /api/settlements/run`
```json
{
  "profile": [{"ts":"...","kwh":120.0}],      // or pass "company_id" later
  "prices": {"area":"tokyo","date":"2025-10-23"}, // backend fetches JEPX
  "pv_offset_pct": 0.15
}
```
**Response:**
```json
{
  "period": {"from":"...","to":"..."},
  "totals": {"kwh":12345.6,"cost_yen":302100.0},
  "by_hour": [{"ts":"...","kwh":...,"price":...,"cost":...}],
  "assumptions": {"pv_offset_pct":0.15,"area":"tokyo"},
  "source_prices": {"name":"JEPX","url":"..."}
}
```

**Unit rules:** Demand source is MW hourly → convert to kWh per slot; JEPX is JPY/kWh.

---

## 4. Adapters & Normalization
- **Adapters are pure functions**: input artifact (CSV/HTML/JSON) → validated typed structs → normalized JSON documents per schemas in `/pipeline/schemas`.
- Timestamps: normalize to `Asia/Tokyo` for output; internal math allowed in UTC.  
- CSV variability: column order/header changes → adapter must auto-detect using header names; hard failures if required columns missing.
- Add `meta.warning` when source known to occasionally omit forecast.

---

## 5. Caching, Storage, and Modes
- **Mode A (Static):** Write artifacts to `/public/data/{provider}/{kind}/{y}/{m}/{d}.json`. Frontend fetches static files.
- **Mode B (Serverless):** KV/Edge cache with TTLs (TEPCO/Kansai hourly, JEPX daily, OCCTO daily, Open‑Meteo 3–6h). On miss → fetch → normalize → cache → respond.
- No secrets needed for JP sources; keep `ENTSOE_TOKEN`, `EIA_API_KEY` optional for world fallback.

---

## 6. NFRs & Security
- **Cost:** Minimize to ~0 JPY (prefer Mode A).  
- **Latency:** static < 100ms; serverless < 1000ms at cache hit.  
- **CORS:** allow the Vercel frontend origin.  
- **Attribution:** every response must include `source.name` and `source.url`.  
- **Licenses:** respect source ToS; throttle politely; never hot‑loop on vendor sites.

---

## 7. Testing Strategy
- Golden tests for each adapter with pinned sample files.  
- Boundary tests for date edges; JP has no DST but test month/day boundaries.  
- Empty/partial data: forecast missing → still 200 with `meta.warning`.  
- Settlement math: deterministic sample with flat profile × flat price.  
- JSON Schema validation for public artifacts.

---

## 8. CI/CD & Quality Gates
- GitHub Actions:
  - **Pipeline cron** (Mode A): fetch→normalize→commit JSON.
  - **Build & test** on PR; fail on lint or schema drift.
- Go: `golangci-lint`, table‑driven tests, `testcontainers` for Postgres when used.
- Python: `ruff`, `mypy --strict`, `pytest` with fixtures.  
- Conventional Commits; semantic versioning of artifact schemas.

---

## 9. Observability
- Structured logs (JSON) with `source`, `artifact`, `count`, `duration_ms`.  
- Error budget: retry with backoff (max 3), circuit‑break per source on repeated failures.  
- Emit minimal metrics (counts, durations) to logs (future: OTEL).

---

## 10. Idempotency & Determinism
- Adapters must be deterministic for the same artifact.  
- Settlement endpoint must be **idempotent** if called with the same payload.  
- Where rounding is needed, specify: money to 0.1 JPY, energy to 0.1 kWh (configurable).

---

## 11. Configuration (.env / secrets)
```
OPEN_METEO_BASE=...
OCCTO_BASE=...
TEPCO_BASE=...
KANSAI_BASE=...
JEPX_BASE=...
# Optional:
ENTSOE_TOKEN=...
EIA_API_KEY=...
FRONTEND_ORIGIN=https://corporate-energy-platform.vercel.app
```

---

## 12. Agent Rules of Engagement
1. Before coding, **restate requirements** and list assumptions.  
2. Generate a **small design doc** per feature (what structures, endpoints, tests).  
3. Implement minimal increment; write tests first (or alongside).  
4. Run self‑review checklist (see §13) and produce a PR with changelog.  
5. If Context7 suggests a different pattern, propose a patch to this spec (do not silently diverge).

---

## 13. Self‑Review Checklist (Agent)
- [ ] API shape matches §3 contracts exactly.  
- [ ] Timezones and units correct; no DST assumptions for JP.  
- [ ] Errors graceful; warnings surfaced in `meta.warning` when partial.  
- [ ] Tests cover happy/edge/empty.  
- [ ] Linters green; no secrets in logs.  
- [ ] Source attribution present.

---

## 14. Future Extensions (non-goals for MVP)
- AuthN/AuthZ, multi‑tenant orgs, DB persistence, roles.  
- ENTSO‑E/EIA world endpoints.  
- Advanced ML forecast and PV modeling.  
- PDF/CSV exporters.
