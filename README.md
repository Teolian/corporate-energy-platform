# Corporate Energy Benchmark Platform

> **"Datazora UI meets Energy Industry Data"**

A modern web platform for comparing corporate energy performance with industry leaders. Built with Vue 3, TypeScript, and Go.

## Overview

This platform allows companies to:
- 📊 Compare energy consumption with competitors
- 🔍 Search and filter companies by industry and efficiency
- 📈 View forecasts and benchmark rankings
- ✅ Track forecast accuracy vs actual consumption
- 💡 Get AI-powered insights and recommendations

## Tech Stack

### Frontend (Primary Focus)
- **Vue 3.5.x** - Progressive JavaScript framework
- **TypeScript 5.6.x** - Type-safe development
- **Pinia 2.2.x** - State management
- **Vue Router 4.5.x** - Routing
- **TailwindCSS 3.4.x** - Utility-first CSS
- **Chart.js 4.4.x** - Data visualization
- **TanStack Table 8.x** - Advanced data tables

### Backend
- **Go 1.23.x** + Gin - REST API
- **Python 3.12** + FastAPI - Forecasting service
- **PostgreSQL 16** - Database
- **Redis 7.x** - Caching

## Project Structure

```
aversome/
├── frontend/          # Vue 3 application
│   ├── src/
│   │   ├── components/  # UI components
│   │   ├── views/       # Page views
│   │   ├── stores/      # Pinia stores
│   │   ├── services/    # API services
│   │   └── types/       # TypeScript types
│   └── package.json
├── backend/           # Go API
│   ├── cmd/           # Main applications
│   ├── internal/      # Private code
│   └── go.mod
├── data/              # Sample datasets
└── docs/              # Documentation
```

## Key Features

### 1. Company Database
- Searchable table with real company energy data
- Advanced filtering (industry, size, efficiency)
- Tag-based categorization
- Sortable columns with pagination

### 2. Company Detail Page
- Energy consumption trends
- Forecast vs Actual comparison
- Benchmark rankings
- AI-powered insights

### 3. Industry Comparison
- Side-by-side company comparisons
- Efficiency score rankings
- Cost per employee analysis

## Target Companies

This project was designed for:
- **Datazora** - Demonstrating Vue 3 + TypeScript + Pinia skills with enterprise search UI
- **Digital Grid** - Showcasing energy industry domain knowledge with benchmarking and forecasting

## Getting Started

See [CLAUDE.md](./CLAUDE.md) for detailed development instructions.

## Data Sources

Real data from publicly available sources:
- Fortune 500 energy reports
- Corporate ESG disclosures
- Industry benchmark databases

## License

MIT
