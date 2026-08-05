/**
 * Single source of truth for all time display and date-grouping logic
 * (DATA-05). Everything renders in the *user's* timezone — never a
 * hard-coded zone — while the backend stores UTC only.
 *
 * The active timezone is a module-level value (default Asia/Shanghai for
 * backwards compatibility) hydrated from the logged-in profile and the
 * profile page; see useAuthStore/setUserTimezone.
 */

let activeTimezone = 'Asia/Shanghai'

export const DEFAULT_TIMEZONE = 'Asia/Shanghai'

/** Rejects invalid IANA names (throws in Intl) by falling back. */
function normalizeTz(tz: string): string {
  try {
    new Intl.DateTimeFormat('en-US', { timeZone: tz })
    return tz
  } catch {
    return DEFAULT_TIMEZONE
  }
}

export function setUserTimezone(tz: string | null | undefined) {
  if (!tz || typeof tz !== 'string' || !tz.trim()) {
    // Reset to default (e.g. logout on a shared device).
    activeTimezone = DEFAULT_TIMEZONE
    try {
      localStorage.removeItem('wm-timezone')
    } catch { /* storage unavailable */ }
    return
  }
  activeTimezone = normalizeTz(tz.trim())
  try {
    localStorage.setItem('wm-timezone', activeTimezone)
  } catch { /* storage unavailable */ }
}

export function getUserTimezone(): string {
  try {
    const stored = localStorage.getItem('wm-timezone')
    if (stored) activeTimezone = normalizeTz(stored)
  } catch { /* storage unavailable */ }
  return activeTimezone
}

const MONTH_DAY: Intl.DateTimeFormatOptions = { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hourCycle: 'h23' }

function partsFor(utcStr: string, tz: string, opts: Intl.DateTimeFormatOptions) {
  const d = new Date(utcStr)
  if (Number.isNaN(d.getTime())) return null
  return new Intl.DateTimeFormat('en-US', { ...opts, timeZone: normalizeTz(tz) }).formatToParts(d)
}

function part(parts: Intl.DateTimeFormatPart[], type: string): string {
  return parts.find((p) => p.type === type)?.value ?? ''
}

/** "MM-DD HH:mm" in the user's timezone (the classic schedule format). */
export function formatKickoff(utcStr: string | null | undefined, tz = getUserTimezone()): string {
  if (!utcStr) return ''
  const parts = partsFor(utcStr, tz, MONTH_DAY)
  if (!parts) return ''
  return `${part(parts, 'month')}-${part(parts, 'day')} ${part(parts, 'hour')}:${part(parts, 'minute')}`
}

/** "MM-DD HH:mm 星期X" for detail pages. */
export function formatKickoffLong(utcStr: string | null | undefined, tz = getUserTimezone()): string {
  if (!utcStr) return ''
  const parts = partsFor(utcStr, tz, { ...MONTH_DAY, weekday: 'short' })
  if (!parts) return ''
  return `${part(parts, 'month')}-${part(parts, 'day')} ${part(parts, 'hour')}:${part(parts, 'minute')} ${part(parts, 'weekday')}`
}

/** "YYYY-MM-DD" calendar day key in the user's timezone (grouping). */
export function dayKey(utcStr: string | null | undefined, tz = getUserTimezone()): string {
  if (!utcStr) return ''
  const d = new Date(utcStr)
  if (Number.isNaN(d.getTime())) return ''
  const parts = new Intl.DateTimeFormat('en-CA', { timeZone: normalizeTz(tz), year: 'numeric', month: '2-digit', day: '2-digit' }).formatToParts(d)
  return `${part(parts, 'year')}-${part(parts, 'month')}-${part(parts, 'day')}`
}

/** Today's date key in the user's timezone. `now` is injectable for tests. */
export function todayKey(tz = getUserTimezone(), now: Date = new Date()): string {
  return dayKey(now.toISOString(), tz)
}

export function isToday(utcStr: string | null | undefined, tz = getUserTimezone(), now: Date = new Date()): boolean {
  return !!utcStr && dayKey(utcStr, tz) === todayKey(tz, now)
}

export function isTomorrow(utcStr: string | null | undefined, tz = getUserTimezone(), now: Date = new Date()): boolean {
  if (!utcStr) return false
  // Calendar-day arithmetic on the YYYY-MM-DD key itself: going through the
  // timezone again would mis-handle negative offsets (a UTC instant may
  // still be "today" locally on the target day).
  const [y, m, d] = todayKey(tz, now).split('-').map(Number)
  const tomorrow = new Date(Date.UTC(y, m - 1, d + 1)).toISOString().slice(0, 10)
  return dayKey(utcStr, tz) === tomorrow
}

/** Short "HH:mm" in the user's timezone. */
export function formatTime(utcStr: string | null | undefined, tz = getUserTimezone()): string {
  if (!utcStr) return ''
  const parts = partsFor(utcStr, tz, { hour: '2-digit', minute: '2-digit', hourCycle: 'h23' })
  if (!parts) return ''
  return `${part(parts, 'hour')}:${part(parts, 'minute')}`
}
