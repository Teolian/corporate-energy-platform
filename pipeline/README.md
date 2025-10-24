# Pipeline — Data Normalization & Validation

Pipeline jobs for fetching, normalizing, and validating energy data artifacts.

## Purpose

Generate static JSON files from public data sources following AGENT_TECH_SPEC contracts. These artifacts are consumed by the frontend in LIVE mode.

## Structure

```
pipeline/
├── schemas/              # JSON Schema definitions
│   └── demand-schema.json
└── README.md
```

**Note:** Pipeline jobs are located in `backend/cmd/` to access internal packages.

## JSON Schemas

### Demand Schema (`schemas/demand-schema.json`)

Validates demand response format per AGENT_TECH_SPEC §3.1:
- Required fields: `area`, `date`, `timezone`, `timescale`, `series`, `source`
- Series: 1-24 hourly data points
- Timestamps: ISO8601 with `+09:00` offset
- Optional: `meta.warning` for partial data

### Validation (Optional)

Using `ajv-cli`:

```bash
npm install -g ajv-cli

ajv validate -s schemas/demand-schema.json \
             -d ../public/data/jp/tokyo/demand-2025-10-24.json
```

## Pipeline Jobs

Jobs are Go programs in `backend/cmd/`:

### Fetch Tokyo Demand

```bash
cd backend
go run cmd/fetch-demand/main.go -date 2025-10-24
```

**Output:** `../public/data/jp/tokyo/demand-2025-10-24.json`

### Future Jobs

- `fetch-kansai-demand` — Kansai Electric Power
- `fetch-reserve` — OCCTO reserve margin
- `fetch-jepx-prices` — JEPX day-ahead prices
- `fetch-weather` — Open-Meteo solar/weather

## GitHub Actions (Future)

Planned cron workflow:

```yaml
name: Update Data Artifacts
on:
  schedule:
    - cron: '0 */6 * * *'  # Every 6 hours
  workflow_dispatch:        # Manual trigger

jobs:
  fetch:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Fetch Tokyo demand
        run: |
          cd backend
          go run cmd/fetch-demand/main.go -date $(date +%Y-%m-%d)

      - name: Validate output
        run: |
          npm install -g ajv-cli
          ajv validate -s pipeline/schemas/demand-schema.json \
                       -d public/data/jp/tokyo/demand-*.json

      - name: Commit artifacts
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add public/data/
          git commit -m "chore(data): update artifacts $(date +%Y-%m-%d)" || true
          git push
```

## Output Directory Structure

```
public/data/
├── jp/
│   ├── tokyo/
│   │   ├── demand-2025-10-24.json
│   │   ├── demand-2025-10-25.json
│   │   └── ...
│   ├── kansai/
│   │   └── demand-YYYY-MM-DD.json
│   ├── system/
│   │   └── reserve-YYYY-MM-DD.json
│   └── jepx/
│       └── spot-AREA-YYYY-MM-DD.json
└── weather/
    └── solar-LAT-LON-YYYY-MM-DD.json
```

## Artifact Requirements

Per AGENT_TECH_SPEC §4:

1. **Timestamps:** ISO8601 with `+09:00` offset
2. **Source attribution:** Every artifact includes `source.name` and `source.url`
3. **Schema version:** Top-level `schema_version` field (future)
4. **Warnings:** Use `meta.warning` for partial/missing data
5. **Deterministic:** Same input → same output

## Development

### Adding a New Schema

1. Create `schemas/{domain}-schema.json`
2. Follow JSON Schema Draft 07 format
3. Document required fields and constraints
4. Add validation to CI/CD

### Testing Schemas

```bash
# Good: Valid data
ajv validate -s schemas/demand-schema.json -d testdata/valid.json

# Bad: Invalid data (should fail)
ajv validate -s schemas/demand-schema.json -d testdata/invalid.json
```

## Troubleshooting

### Schema validation fails
- Check timestamp format: `YYYY-MM-DDTHH:MM:SS+09:00`
- Ensure `source.name` and `source.url` are present
- Verify all required fields exist

### Stale data
- Run pipeline jobs manually with `workflow_dispatch`
- Check cron schedule in GitHub Actions
- Review job logs for errors

## License

Private portfolio project.
