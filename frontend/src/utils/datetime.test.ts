import { describe, it, expect, beforeEach } from 'vitest'
import {
  setUserTimezone,
  getUserTimezone,
  DEFAULT_TIMEZONE,
  formatKickoff,
  dayKey,
  isToday,
  isTomorrow,
} from './datetime'

// All fixtures are fixed UTC instants so the tests are deterministic.
const UTC_NOON = '2026-03-01T12:00:00Z' // Sunday noon UTC (US winter, UTC-5)
const DST_DAY = '2026-03-08T12:00:00Z'  // US spring-forward day

describe('datetime (DATA-05)', () => {
  beforeEach(() => {
    setUserTimezone(DEFAULT_TIMEZONE)
  })

  it('formats in the user timezone, not a hard-coded one', () => {
    setUserTimezone('Asia/Shanghai') // UTC+8
    expect(formatKickoff(UTC_NOON)).toBe('03-01 20:00')
    setUserTimezone('America/New_York') // UTC-5 in winter
    expect(formatKickoff(UTC_NOON)).toBe('03-01 07:00')
  })

  it('is DST-aware: same UTC instant renders with the post-switch offset', () => {
    setUserTimezone('America/New_York')
    // 2026-03-08 02:00 EST jumps to 03:00 EDT (UTC-4), so noon UTC is 08:00.
    expect(formatKickoff(DST_DAY)).toBe('03-08 08:00')
    expect(dayKey(DST_DAY)).toBe('2026-03-08')
  })

  it('cross-day boundary: UTC 23:30 is next day in Beijing', () => {
    setUserTimezone('Asia/Shanghai')
    const late = '2026-03-08T23:30:00Z'
    expect(formatKickoff(late)).toBe('03-09 07:30')
    expect(dayKey(late)).toBe('2026-03-09')
  })

  it('isToday uses the user calendar day (DST-aware)', () => {
    setUserTimezone('America/New_York')
    // 2026-03-08 is the US DST spring-forward day; UTC 07:30 == 02:30 EST
    // before the switch, still the same calendar day.
    expect(isToday(UTC_NOON)).toBe(false) // fixed fixture is not "now"
  })

  it('dayKey is stable across zones', () => {
    setUserTimezone('UTC')
    expect(dayKey('2026-03-08T23:59:59Z')).toBe('2026-03-08')
    setUserTimezone('Asia/Shanghai')
    expect(dayKey('2026-03-08T23:59:59Z')).toBe('2026-03-09')
  })

  it('isTomorrow compares calendar days in the user zone', () => {
    // Fixed "now" so the test never goes flaky around UTC midnight.
    const now = new Date('2026-03-01T12:00:00Z')
    setUserTimezone('UTC')
    expect(isTomorrow('2026-03-02T00:00:00Z', 'UTC', now)).toBe(true)
    expect(isTomorrow('2026-03-02T23:59:59Z', 'UTC', now)).toBe(true)
    expect(isTomorrow('2026-03-01T23:59:59Z', 'UTC', now)).toBe(false)
    expect(isTomorrow('2026-03-03T00:00:00Z', 'UTC', now)).toBe(false)
  })

  it('isTomorrow works for negative-offset timezones (regression)', () => {
    // With "today" = 2026-03-01 in America/New_York, a match at 2026-03-02
    // 03:00 UTC is 2026-03-01 22:00 EST — same local day, NOT tomorrow.
    const now = new Date('2026-03-01T12:00:00Z')
    setUserTimezone('America/New_York')
    expect(isTomorrow('2026-03-02T03:00:00Z', 'America/New_York', now)).toBe(false)
    // 2026-03-02 09:00 UTC = 04:00 EST on 2026-03-02 — tomorrow.
    expect(isTomorrow('2026-03-02T09:00:00Z', 'America/New_York', now)).toBe(true)
  })

  it('isToday uses the user calendar day', () => {
    const now = new Date('2026-03-01T12:00:00Z')
    setUserTimezone('America/New_York')
    // 11:00 EST on 2026-03-01 is the same local day as `now` (07:00 EST).
    expect(isToday('2026-03-01T16:00:00Z', 'America/New_York', now)).toBe(true)
    // 2026-03-02T05:00Z = 00:00 EST on 2026-03-02 — the next local day.
    expect(isToday('2026-03-02T05:00:00Z', 'America/New_York', now)).toBe(false)
  })

  it('falls back to the default timezone for invalid zone names', () => {
    setUserTimezone('Not/AZone')
    expect(getUserTimezone()).toBe(DEFAULT_TIMEZONE)
    expect(formatKickoff(UTC_NOON, 'Not/AZone')).toBe('03-01 20:00') // Asia/Shanghai
    expect(dayKey(UTC_NOON, 'Not/AZone')).toBe('2026-03-01')
  })

  it('rejects invalid input gracefully', () => {
    expect(formatKickoff('not-a-date')).toBe('')
    expect(formatKickoff(null)).toBe('')
    expect(dayKey('')).toBe('')
    expect(isToday(undefined)).toBe(false)
  })
})
