# FRONTEND_GREENFIELD_ALIGNMENT.md
**Project:** Corporate Energy Platform — Japan Energy Dashboard (Greenfield FE)  
**Goal:** Align the new simple frontend with existing backend contracts and the DG-aligned roadmap, without breaking the current demo.

> This file complements: `AGENT_TECH_SPEC.md` (§3), `FRONTEND_PLAN_JAPAN_DASHBOARD.md`, `JAPAN_DASHBOARD_SPRINT_PLAN.md`, `AGENT_INTEGRATION_GUIDE.md`, `AGENT_TOPOLOGY.md`.
> Canonical data shapes remain those in `AGENT_TECH_SPEC.md §3`.

---

## 1) Adopt Greenfield FE as `frontend-live`
- Keep the original demo as `frontend-mock` (unchanged).
- Treat the **new simple FE** as `frontend-live` (live data + new features).
- Use the same dual-mode switch (MOCK/LIVE), but default **LIVE** for preview deployments only.

---

## 2) Data Access Contract (no UI refactor needed)
Expose these functions in a minimal data client for the greenfield FE (keep component props stable):

```ts
// src/services/dataClient.ts (greenfield)
export async function getDemand(area: 'tokyo' | 'kansai', date: ISODate): Promise<DemandSeries>;
export async function getReserve(date: ISODate): Promise<ReserveByArea>;
export async function getJepxSpot(area: 'tokyo' | 'kansai', date: ISODate): Promise<PriceSeries>;
export async function runSettlement(input: SettlementInput): Promise<SettlementResult>;
```

**Mapping (domain → UI):**
- Demand: `{ ts, demand_mw, forecast_mw? } → { x, yActual, yForecast? }`
- Reserve: `{ reserve_margin_pct, status } → { pct, status }` for badge/gauge
- Prices: `{ price_yen_per_kwh[].{ ts, price } } → { x, yPrice }`
- Settlement: `{ totals, by_hour[] } → cost card + table DTOs

**Sources & meta:**
- Surface `source.name` and `source.url` in a small footer line per widget.
- Show `meta.warning` non-blocking if present (e.g., missing forecast).

---

## 3) Files & Folders (suggested)
```
src/
  components/
    demand/DemandChart.vue
    prices/SpotPriceChart.vue
    reserve/ReserveBadge.vue
    reserve/ReserveGauge.vue
    settlement/CostCard.vue
    settlement/CostBreakdownTable.vue
  stores/
    demand.ts
    prices.ts
    reserve.ts
    settlement.ts
  services/
    dataClient.ts        // dual-mode fetcher + adapters
  composables/
    useExport.ts
  views/
    Dashboard.vue        // demand + reserve badge + price overlay + cost card
    Prices.vue           // dedicated JEPX view
    Comparison.vue       // Tokyo vs Kansai overlay
```

---

## 4) UI Tasks (P1/P2)
### P1 (DGP-aligned)
- **Reserve badge/gauge** (OCCTO) on the dashboard. Tooltip with thresholds (>=8 stable, 5–8 watch, <5 tight).
- **Price overlay (JEPX)** on each demand chart; dedicated `/prices` view.
- **Cost card (Settlement-lite)** with totals + optional hourly breakdown; CSV export.

### P2 (UX)
- **Comparison** route with Tokyo vs Kansai overlay.
- **Exports**: CSV of demand for current day; PNG of the chart.
- **Empty/partial states**: neutral messages; dashed segments for price gaps.

---

## 5) Chart.js Guidance
- Use a second Y-axis for prices (¥/kWh) on overlays.
- Keep **solid line** for Actual, **dashed** for Forecast; consistent colors with Tailwind tokens.
- X axis: 24 hourly ticks in Asia/Tokyo; tooltip shows local time + units.
- Handle gaps: render `null` and show a tiny warning icon with tooltip.

---

## 6) ENV & Flags
- `VITE_DATA_MODE=mock|live` (preview: live; main: mock).
- Feature flags (env or runtime config):
  - `VITE_FEAT_RESERVE`, `VITE_FEAT_JEPX`, `VITE_FEAT_SETTLEMENT`, `VITE_FEAT_COMPARISON`
- Keep flags **ON** only in preview until validated.

---

## 7) Acceptance Criteria (Front)
- Dashboard renders with **demand + reserve badge + price overlay + cost card** in LIVE.
- Warnings are non-blocking; source attribution visible.
- `/prices` and `/comparison` routes work; CSV/PNG export OK.
- No breaking changes to component props; all adaptations live in `dataClient.ts`.

---

## 8) Visual Feedback on the Current Screenshot
- Good: clear dual charts, readable legends, Peak/Average/Accuracy cards, dark-mode-ready.
- Recommended tweaks:
  1) Add **source footers** under each chart (TEPCO, Kansai EPCO; later JEPX/OCCTO).
  2) Add **toggle** “Show price overlay” (hidden until JEPX adapter ready).
  3) Show **neutral badge** “Forecast unavailable” when forecast missing.
  4) Format Y-axes with **thin grid** and consistent units (MW left; ¥/kWh right when overlay is on).
  5) Keep **date picker** the single source of truth; prev/next snap to local JP day boundaries.
  6) Expose **accuracy formula** on hover (e.g., 1 - MAPE over available overlap).

---

## 9) Rollout
- Ship P1 widgets behind flags, test on preview (LIVE).
- Keep `main` on MOCK; once P1 stable → flip flags to on in production env.
- Always provide a quick **flip back to MOCK** path if data is unavailable.

