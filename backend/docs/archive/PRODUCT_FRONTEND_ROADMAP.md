# PRODUCT_FRONTEND_ROADMAP.md
**Project:** Corporate Energy Platform — Japan Energy Dashboard  
**Audience:** AI coding agent + maintainer  
**Purpose:** One document with the *future plan* and *frontend instructions* (greenfield FE accepted as `frontend-live`). API contracts remain canonical in `AGENT_TECH_SPEC.md §3`.

---

## 1) Roadmap (aligned to Digital Grid services)
**P1 — System status + prices + cost (DGP-aligned)**
- Backend (Mode A static first):
  - OCCTO Reserve Margin adapter → `/public/data/jp/system/reserve-YYYY-MM-DD.json`
  - JEPX Day-Ahead Prices adapter → `/public/data/jp/jepx/spot-AREA-YYYY-MM-DD.json`
  - Settlement‑lite compute (kWh × ¥/kWh, deterministic rounding)
- Frontend (`frontend-live`):
  - ReserveBadge + optional ReserveGauge (thresholds: ≥8 stable, 5–8 watch, <5 tight)
  - Price overlay on demand charts + dedicated `/prices` view
  - CostCard with totals (+ optional hourly table) and CSV export

**P2 — Market context UX (still DGP)**
- Tokyo vs Kansai **comparison** route (`/comparison`)
- Full **SpotPriceChart** page (price-only and price+overlay)
- Export **CSV/PNG** for demand/price charts

**P3 — Adjacent to other DG services (optional)**
- **Portfolio Mix** (fixed/spot/RE share) → DGP
- **Certificates calendar** → Econohashi
- **PPA digest** (press-release aggregator) → RE Bridge
- **Battery arbitrage simulator** → Battery Aggregation

---

## 2) Frontend Instructions (what to implement)
### 2.1 Data Client (single adapter layer)
Create `src/services/dataClient.ts` that exposes **only**:
```ts
getDemand(area: 'tokyo' | 'kansai', date: ISODate): Promise<DemandSeries>
getReserve(date: ISODate): Promise<ReserveByArea>
getJepxSpot(area: 'tokyo' | 'kansai', date: ISODate): Promise<PriceSeries>
runSettlement(input: SettlementInput): Promise<SettlementResult>
```
- Map backend shapes → UI DTOs; **не меняем пропсы компонентов**.
- Показывать `source.name` + `source.url` в футере каждого виджета.
- Если есть `meta.warning` — отобразить компактно, **без блокировки UI**.

### 2.2 Stores & Components
- **stores**: `reserve.ts`, `prices.ts`, `settlement.ts`
- **components**:
  - `ReserveBadge.vue`, `ReserveGauge.vue`
  - `SpotPriceChart.vue` (+ overlay toggle в `DemandChart.vue`)
  - `CostCard.vue`, `CostBreakdownTable.vue`
- **routes**: `/` (Dashboard), `/prices`, `/comparison`

### 2.3 Chart.js правила
- Линии: Actual — solid; Forecast — dashed; Price — отдельная ось **¥/kWh** справа.
- Ось X: 24 часа **Asia/Tokyo**; тултип: локальное время + единицы.
- Пробелы в данных → `null` + dashed gap + иконка-подсказка.

### 2.4 Flags & Modes
- `VITE_DATA_MODE=mock|live` (preview=live, main=mock)
- Фичи под флагами (ON в preview, OFF на main):
  - `VITE_FEAT_RESERVE`, `VITE_FEAT_JEPX`, `VITE_FEAT_SETTLEMENT`, `VITE_FEAT_COMPARISON`

### 2.5 Acceptance (Front P1/P2)
- В LIVE видны: demand + **reserve badge** + **price overlay** + **cost card**.
- `/prices` и `/comparison` работают; экспорт CSV/PNG корректный.
- UI устойчив к пустым данным; атрибуция источников присутствует.

---

## 3) Branching & Safety
- Старый демо оставить как `frontend-mock` (MOCK only, стабильный).
- Новый простой фронт принять как `frontend-live` (LIVE + flags).
- Рабочие ветки: UI в `feat/frontend-japan`, бек в `feat/backend-core`.
- Промоутить в `main` только после успешного превью и валидных схем.

---

## 4) References
- Canonical contracts: `AGENT_TECH_SPEC.md §3`
- Детализация задач: `JAPAN_DASHBOARD_SPRINT_PLAN.md`
- Уточнения по UI: `FRONTEND_PLAN_JAPAN_DASHBOARD.md`
- Гайд для нового фронта: `FRONTEND_GREENFIELD_ALIGNMENT.md`
