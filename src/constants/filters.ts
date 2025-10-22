/**
 * Filter-related constants
 *
 * Centralized definition of filter ranges and thresholds
 * to avoid magic numbers throughout the codebase.
 */

/**
 * Range boundaries for numeric filters
 */
export const FILTER_RANGES = {
  EFFICIENCY: {
    MIN: 0,
    MAX: 100
  },
  RENEWABLE: {
    MIN: 0,
    MAX: 100
  }
} as const

/**
 * Efficiency score thresholds for visual indicators
 */
export const EFFICIENCY_THRESHOLDS = {
  EXCELLENT: 90, // Green indicator
  GOOD: 80, // Blue indicator
  AVERAGE: 0 // Yellow indicator
} as const

/**
 * Renewable energy percentage threshold for "green" companies
 */
export const RENEWABLE_THRESHOLD = {
  HIGH: 80 // Companies above this are considered highly sustainable
} as const

/**
 * Default number of companies to show in previews
 */
export const PREVIEW_LIMITS = {
  HOME_PAGE: 5,
  SIDEBAR: 3
} as const
