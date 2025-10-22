# Corporate Energy Benchmark Platform

> A Vue 3 application for comparing corporate energy consumption and efficiency metrics.

**Live Demo:** [https://corporate-energy-platform.vercel.app](https://corporate-energy-platform.vercel.app)

**Status:** Frontend-only implementation with mocked data. Backend API is planned for future development.

## Overview

This application simulates an energy benchmarking system using Fortune 500 company data. Built with Vue 3, TypeScript, and TailwindCSS.

**Tech stack:**
- Vue 3.5 with Composition API
- TypeScript (strict mode)
- Pinia for state management
- TanStack Table for data tables
- Chart.js for visualizations
- TailwindCSS for styling
- Dark mode support
- Responsive design

**Note:** Currently uses mocked data. Backend API integration is planned for future development.

### Features

- Company database with 50+ Fortune 500 entries
- Search and filter by industry, efficiency score, renewable percentage
- 12-month energy consumption forecasts (mocked data)
- Industry benchmark comparisons
- Interactive charts for trends and analytics
- Dark mode theme

## Screenshots

### Home Page

![Home Page](./screenshots/01-home.png)

### Companies Database

![Companies Database](./screenshots/02-companies.png)

### Company Detail

![Company Detail](./screenshots/03-company-detail.png)

### Dark Mode

<table>
  <tr>
    <td><img src="./screenshots/04-home-dark.png" alt="Home Dark" /></td>
    <td><img src="./screenshots/05-companies-dark.png" alt="Companies Dark" /></td>
    <td><img src="./screenshots/06-company-detail-dark.png" alt="Detail Dark" /></td>
  </tr>
  <tr>
    <td align="center">Home Page (Dark)</td>
    <td align="center">Companies Table (Dark)</td>
    <td align="center">Company Detail (Dark)</td>
  </tr>
</table>

## Tech Stack

### Frontend
- **Vue 3.5** - Progressive JavaScript framework with Composition API
- **TypeScript 5.6** - Type-safe development with strict mode
- **Pinia 2.2** - State management library
- **Vue Router 4.5** - Client-side routing
- **TailwindCSS 3.4** - Utility-first CSS framework
- **Chart.js 4.4** - Data visualization library
- **TanStack Table 8** - Advanced table features
- **VueUse** - Collection of Vue composition utilities

### Development Tools
- **Vite 6** - Fast build tool and dev server
- **Playwright** - End-to-end testing and screenshots
- **ESLint** - Code linting
- **Prettier** - Code formatting
- **TypeScript ESLint** - TypeScript-specific linting rules

## Getting Started

### Prerequisites

- Node.js 20.x or higher
- npm 10.x or higher

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The application will be available at `http://localhost:5173`

### Available Scripts

```bash
npm run dev          # Start development server
npm run build        # Build for production
npm run preview      # Preview production build
npm run type-check   # Run TypeScript type checking
npm run lint         # Lint code with ESLint
npm run format       # Format code with Prettier
```

## Project Structure

```
frontend/
├── src/
│   ├── components/     # Reusable Vue components
│   ├── views/          # Page components
│   ├── stores/         # Pinia state management
│   ├── services/       # API services
│   ├── types/          # TypeScript type definitions
│   ├── data/           # Mock data (Fortune 500 companies)
│   ├── constants/      # Application constants
│   ├── utils/          # Utility functions
│   └── router/         # Vue Router configuration
├── public/             # Static assets
└── screenshots/        # Application screenshots
```

## Pages

- **`/`** - Home page with feature overview
- **`/companies`** - Searchable table with filters and sorting
- **`/companies/:id`** - Individual company details with charts
- **`/compare`** - Coming soon

## Data Model

The application uses real data from Fortune 500 company sustainability reports. Each company entry includes:

- Company name and industry
- Location and employee count
- Annual energy consumption (MWh)
- Efficiency score (0-100)
- Renewable energy percentage
- Energy trend indicators
- ESG tags

## Development

- Vue 3 Composition API with `<script setup>`
- TypeScript strict mode
- TailwindCSS for styling
- Type definitions in `src/types/`

## Deployment

### Docker

Build and run using Docker:

```bash
# Build production image
docker build -t energybench-frontend .

# Run container
docker run -p 3000:80 energybench-frontend
```

Or use Docker Compose:

```bash
# Production mode
docker-compose up frontend

# Development mode
docker-compose --profile development up
```

### Environment Variables

Create a `.env` file for configuration:

```env
VITE_API_URL=http://localhost:8080/api
```

## Future Plans

- Backend API (Go + PostgreSQL)
- Company comparison features
- Data export (PDF, CSV)
- User authentication

## License

This project is a portfolio demonstration and is not licensed for commercial use.

## Contact

For questions or collaboration opportunities, please reach out through GitHub.
