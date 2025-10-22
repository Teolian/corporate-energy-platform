/**
 * Application logger service
 *
 * Provides a centralized logging interface that can be easily replaced
 * with external logging services (Sentry, LogRocket, etc.) in production.
 *
 * @example
 * ```ts
 * logger.error('Failed to fetch data', { userId: '123', error })
 * logger.warn('Deprecated API usage')
 * logger.info('User logged in', { userId: '123' })
 * ```
 */

interface LogContext {
  [key: string]: unknown
}

class Logger {
  private isDevelopment = import.meta.env.MODE === 'development'

  /**
   * Log debug information (only in development)
   */
  debug(message: string, context?: LogContext): void {
    if (this.isDevelopment) {
      console.debug(`[DEBUG] ${message}`, context || '')
    }
  }

  /**
   * Log informational messages
   */
  info(message: string, context?: LogContext): void {
    if (this.isDevelopment) {
      console.info(`[INFO] ${message}`, context || '')
    }
    // In production, send to analytics service
    // e.g., analytics.track(message, context)
  }

  /**
   * Log warning messages
   */
  warn(message: string, context?: LogContext): void {
    if (this.isDevelopment) {
      console.warn(`[WARN] ${message}`, context || '')
    }
    // In production, send to monitoring service
    // e.g., Sentry.captureMessage(message, 'warning')
  }

  /**
   * Log error messages
   */
  error(message: string, error?: Error | unknown, context?: LogContext): void {
    if (this.isDevelopment) {
      console.error(`[ERROR] ${message}`, error || '', context || '')
    }
    // In production, send to error tracking service
    // e.g., Sentry.captureException(error, { tags: context })
  }
}

/**
 * Singleton logger instance
 */
export const logger = new Logger()
