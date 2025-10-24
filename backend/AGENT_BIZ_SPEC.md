# AGENT_BIZ_SPEC.md
**Project:** Corporate Energy Platform — Business Logic & Plan  
**Audience:** AI product/engineering agent + human PM  
**Purpose:** Explain the *why* and *what* for the product; define stages, value, and acceptance. The tech spec (AGENT_TECH_SPEC.md) is the *how*.

---

## 0. Product Vision
A lightweight, free-to-use demo that shows **real Japanese power system signals** (demand, reserve, prices) and simple **cost estimation** and **what‑if** scenarios for companies. Goal: demonstrate domain understanding and practical value, support internship/job discussions with energy companies.

---

## 1. Target Users & Jobs-to-be-Done
- **Energy/Facility Manager:** see demand trends, peaks, anomalies; check system tightness.  
- **Finance/Procurement:** estimate cost under spot prices; compare scenarios.  
- **ESG/Sustainability:** track % renewables proxy and reporting context.

**Key JTBD statements**
- “I want to know if the grid will be **tight** tomorrow in my area.”  
- “I want to estimate **how much** our consumption would cost at spot prices.”  
- “I want to see **what changes** if we apply a PV offset or switch area proxy.”

---

## 2. Value Proposition
- **Clarity:** clean normalized data from public sources in one place.  
- **Actionability:** quick cost estimate and simple scenario comparison.  
- **Credibility:** visible attribution, conservative assumptions, explicit caveats.

Success metrics (for demo):
- **Activation:** user loads real Japan data without configuration.  
- **Time‑to‑Insight:** < 3 minutes to first insight (peak/tight/cost).  
- **Engagement:** runs one “what‑if” scenario per session.

---

## 3. Data & Calculations (Business Rules)
### 3.1 Demand & Forecast
- Use **TEPCO/Kansai** demand as official system proxies. If forecast is missing for a day, show fact only and a neutral notice.
- Units: MW hourly → convert to kWh per slot for cost math.

### 3.2 Reserve Margin → Status
- Reserve margin thresholds (configurable defaults):  
  - `>= 8%` → **stable**  
  - `5–8%` → **watch**  
  - `< 5%` → **tight**  
- Show a short tooltip explaining statuses.

### 3.3 Prices (JEPX Day‑Ahead)
- Use day‑ahead **area** prices as ¥/kWh. For demo, map “Tokyo” UI area to JEPX Tokyo area (configurable).
- Fetch at most daily; cache results. Display source and date clearly.

### 3.4 Settlement‑Lite
- Input: hourly consumption profile (kWh) and day‑ahead prices (¥/kWh) for the same period and area.  
- Computation: `cost = Σ kWh_slot × price_slot`.  
- Optional `pv_offset_pct` reduces kWh by the given percentage (simple proxy).  
- Output: totals + by‑hour breakdown; show assumptions and source of prices.

### 3.5 Weather/Solar Overlay
- Show GHI/DNI/DHI and cloud cover as context; use them in simple narratives (“High cloud cover aligns with lower PV output”). No hard ESG claims.

---

## 4. UX/Content
- **Region & Date pickers** with sane defaults (today/area=Tokyo).  
- **Badges**: source attribution (TEPCO/Kansai/JEPX/Open‑Meteo).  
- **States**: loading skeletons, empty/partial with clear short messages.  
- **Scenario Panel**: base vs. one alternative (PV offset or different area’s price).  
- **Exports**: copy CSV/PNG (optional in MVP).

Tone: concise, neutral, no investment advice.

---

## 5. Phased Plan (MVP first)
### Phase 1 — Adapter MVP (1–2 days)
- Demand Tokyo/Kansai + normalization → surface on UI.  
- Reserve margin daily status.  
**Acceptance:** endpoints deliver valid JSON with attribution; UI renders without mocks.

### Phase 2 — Weather/Solar (1 day)
- Overlay solar/weather for Tokyo.  
**Acceptance:** overlay toggles on/off; date alignment correct.

### Phase 3 — Prices & Settlement‑Lite (1–2 days)
- JEPX day‑ahead prices → cost compute + simple what‑if.  
**Acceptance:** deterministic totals for provided demo profile; clear assumptions section.

### Phase 4 — “World” Fallback (optional)
- ENTSO‑E/EIA endpoints behind a feature flag.  
**Acceptance:** UI switches to world mode and renders at least one EU/US series.

---

## 6. Risks & Mitigations
- **Source variability (CSV changes):** robust header detection; schema tests; warnings surfaced.  
- **No forecast on some days:** show fact only; do not fail UI.  
- **Legal/ToS:** attribute sources; avoid scraping; throttle politely.  
- **Performance on free tier:** prefer Static Mode; cache aggressively in Serverless Mode.

---

## 7. Non‑Goals (MVP)
- User authentication/authorization.  
- Multi‑tenant persistence and billing.  
- Full retail tariff modeling; advanced ML.

---

## 8. Communication Protocol with the Agent
When you (the agent) implement a slice:
1. **Confirm scope** referencing this file section and AGENT_TECH_SPEC §3.  
2. **Draft a mini‑design** (data flow, structs, endpoints, tests).  
3. **List artifacts** to produce (JSON path, schema version).  
4. **Identify risks** and how you’ll validate with tests.  
5. **Deliver** code + tests + short CHANGELOG entry.  
6. **Post a summary** of what changed and how to revert.

Use **Context7** for patterns (Go/FastAPI/testing). If Context7 diverges from this file, propose an update here before changing behavior.

---

## 9. Definition of Done (Per Feature)
- Contracts match AGENT_TECH_SPEC §3 exactly (no hidden fields).  
- Tests pass and cover happy/edge/empty.  
- Artifacts have source attribution and timestamps.  
- Frontend fetches live data (no mocks) and shows correct states.  
- README in the feature folder includes how to rerun the pipeline locally.

---

## 10. Roadmap (beyond MVP)
- Org/workspace model; role-based access.  
- TimescaleDB and historical analytics.  
- PV modeling with actual radiation + simple rooftop assumptions.  
- Multi-scenario compare and PDF report generation.
