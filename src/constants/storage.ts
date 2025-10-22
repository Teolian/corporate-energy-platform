/**
 * LocalStorage and SessionStorage key constants
 *
 * Centralized definition of storage keys to avoid typos
 * and make it easy to update keys across the application.
 */

/**
 * LocalStorage keys for persistent user preferences
 */
export const STORAGE_KEYS = {
  /**
   * Dark mode preference ('true' | 'false')
   */
  DARK_MODE: 'darkMode',

  /**
   * User's preferred theme ('light' | 'dark' | 'system')
   */
  THEME_PREFERENCE: 'themePreference',

  /**
   * Last viewed company ID for quick access
   */
  LAST_VIEWED_COMPANY: 'lastViewedCompany',

  /**
   * Saved filter preferences for company list
   */
  SAVED_FILTERS: 'savedFilters'
} as const

/**
 * Storage values for dark mode
 */
export const DARK_MODE_VALUES = {
  ENABLED: 'true',
  DISABLED: 'false'
} as const
