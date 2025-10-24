# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Corporate Energy Benchmark Platform - a Vue 3 application for comparing corporate energy consumption and efficiency. Built as a portfolio project demonstrating enterprise-grade frontend development with energy industry domain knowledge.

**Target audience**: Datazora (Vue 3/TypeScript/Pinia skills) and Digital Grid (energy tech expertise)

## Tech Stack

### Frontend (Primary - 80% focus)
- **Vue 3.5.x** with Composition API and `<script setup>` syntax
- **TypeScript 5.6.x** - strict mode enabled
- **Pinia 2.2.x** - state management with typed stores
- **Vue Router 4.5.x** - file-based routing preferred
- **TailwindCSS 3.4.x** - utility-first styling, custom design system
- **Chart.js 4.4.x + vue-chartjs** - data visualization
- **@tanstack/vue-table 8.x** - advanced table features

### Backend (Secondary - 20%)
- **Go 1.23.x** with Gin framework for REST API
- **Python 3.12 + FastAPI** for forecasting calculations
- **PostgreSQL 16** for structured data
- **Redis 7.x** for API response caching

### Documentation
Always use **Context7 MCP** (`mcp__context7__*` tools) for up-to-date library documentation and code examples. Library IDs:
- Vue: `/vuejs/docs`
- Pinia: `/vuejs/pinia`
- TailwindCSS: `/websites/tailwindcss`

## Development Commands

### Frontend (in `frontend/` directory)
```bash
# Install dependencies
npm install

# Development server with hot reload
npm run dev

# Type checking
npm run type-check

# Build for production
npm run build

# Preview production build
npm run preview

# Lint and format
npm run lint
npm run format
```

### Backend (in `backend/` directory)
```bash
# Run Go API server
go run cmd/api/main.go

# Run tests
go test ./...

# Build binary
go build -o bin/api cmd/api/main.go

# Python forecasting service
cd internal/forecast
python -m uvicorn main:app --reload
```

### Database
```bash
# Start PostgreSQL (Docker)
docker-compose up -d postgres

# Run migrations
go run cmd/migrate/main.go up

# Seed sample data
go run cmd/seed/main.go
```

## Architecture

### Data Flow
1. **API Layer** (Go): Handles CRUD operations, caching, authentication
2. **Forecast Service** (Python): Statistical calculations (moving average, linear regression - no complex ML)
3. **Frontend Store** (Pinia): Reactive state management with API integration
4. **Components** (Vue 3): Composables-based, reusable UI elements

### Key Design Patterns
- **Composition API**: Prefer `<script setup>` for all components
- **Composables**: Extract reusable logic into `composables/` (e.g., `useCompanyData`, `useForecast`)
- **Type Safety**: Define all API responses in `types/` directory
- **Store Organization**: One store per domain (companies, forecasts, filters)

### Component Structure
```
components/
├── common/          # Reusable UI (buttons, inputs, cards)
├── company/         # Company-specific (table, detail, filters)
├── charts/          # Chart wrappers (consumption, comparison)
└── layout/          # Layout components (header, sidebar)
```

## Code Style & Conventions

### Vue Components
- Use `<script setup lang="ts">` syntax
- Props with TypeScript interfaces: `defineProps<{ company: Company }>()`
- Emits with type annotations: `defineEmits<{ submit: [data: FormData] }>()`
- Composables for shared logic, not mixins

### TypeScript
- Strict mode enabled - no implicit `any`
- Use interfaces for objects, types for unions/primitives
- API response types in `types/api.ts`
- Domain models in `types/models.ts`

### Styling
- TailwindCSS utility classes in templates
- Custom design tokens in `tailwind.config.js`
- No `<style>` blocks unless absolutely necessary
- Dark mode support via Tailwind's `dark:` variant

### Git Workflow
- Feature branches: `feature/company-detail-page`
- Commit format: `feat: add company comparison chart` (conventional commits)
- Never commit `.env` files or credentials

## Data Model

### Core Entities

```typescript
interface Company {
  id: string
  name: string
  industry: string
  employees: number
  energyConsumption: number // MWh/year
  efficiencyScore: number // 0-100
  renewablePercentage: number // 0-100
  location: string
  tags: string[]
}

interface EnergyForecast {
  companyId: string
  month: string
  predicted: number
  actual?: number
  accuracy?: number
  confidence: number
}

interface Benchmark {
  industry: string
  avgConsumption: number
  avgEfficiency: number
  topPerformers: Company[]
}
```

## Important Notes

### Real Data Sources
Use actual data from:
- Fortune 500 sustainability reports (CSV format in `data/`)
- ESG disclosure databases
- Public energy consumption records

### UI Inspiration
- **Datazora-like table**: Searchable, filterable, with tags
- **Clean charts**: 2-3 meaningful visualizations, not quantity
- **Professional design**: Enterprise-grade, not flashy

### Performance
- Lazy load company detail charts
- Virtual scrolling for large tables (TanStack Table)
- Redis caching for frequently accessed company data
- Debounce search inputs (300ms)

### Common Pitfalls to Avoid
- ❌ Don't create 10+ mediocre charts - make 2-3 excellent ones
- ❌ Don't use `any` type - define proper interfaces
- ❌ Don't inline API calls in components - use composables or stores
- ❌ Don't create new components for simple variations - use props

### Testing Strategy
- Unit tests for composables and utilities
- Integration tests for stores (API mocking)
- E2E tests for critical flows (search, detail view)
- Use Vitest for unit/integration tests

## Project Goals

1. **Showcase Vue 3 expertise** - Modern Composition API patterns
2. **Demonstrate TypeScript proficiency** - Full type safety
3. **Clean, maintainable code** - Easy to understand and extend
4. **Real-world application** - Solve actual business problem
5. **Portfolio quality** - Production-ready code standards

## Development Philosophy

- **Quality over quantity**: 2-3 polished features > 10 half-done ones
- **Real data over mocks**: Use actual company data for authenticity
- **Simple backend**: Backend exists to serve frontend, keep it minimal
- **User-focused**: Every feature must solve a real user need
- **Type safety first**: If TypeScript complains, fix the types, don't use `any`
