import { describe, it, expect } from 'vitest'
import { intervalFor, pollingStatusFrom, LIVE_POLL_MS, NEAR_POLL_MS, IDLE_POLL_MS } from './useLivePolling'

describe('intervalFor (LIVE-01)', () => {
  it('live matches poll every 30s', () => {
    expect(intervalFor({ hasLive: true, nextKickoffInMinutes: 90 })).toBe(LIVE_POLL_MS)
    expect(intervalFor({ hasLive: true, nextKickoffInMinutes: null })).toBe(LIVE_POLL_MS)
  })

  it('kickoff within an hour polls every 60s', () => {
    expect(intervalFor({ hasLive: false, nextKickoffInMinutes: 60 })).toBe(NEAR_POLL_MS)
    expect(intervalFor({ hasLive: false, nextKickoffInMinutes: 0.5 })).toBe(NEAR_POLL_MS)
  })

  it('idle state polls every 5 minutes', () => {
    expect(intervalFor({ hasLive: false, nextKickoffInMinutes: null })).toBe(IDLE_POLL_MS)
    expect(intervalFor({ hasLive: false, nextKickoffInMinutes: 61 })).toBe(IDLE_POLL_MS)
  })
})

describe('pollingStatusFrom (LIVE-01)', () => {
  it('detects live matches', () => {
    const status = pollingStatusFrom([
      { status: 'scheduled', kickoff_time_utc: '2099-01-01T00:00:00Z' },
      { status: 'live', kickoff_time_utc: null },
    ])
    expect(status.hasLive).toBe(true)
  })

  it('finds the nearest upcoming kickoff in minutes', () => {
    const in10m = new Date(Date.now() + 10 * 60_000).toISOString()
    const in90m = new Date(Date.now() + 90 * 60_000).toISOString()
    const status = pollingStatusFrom([
      { status: 'scheduled', kickoff_time_utc: in90m },
      { status: 'scheduled', kickoff_time_utc: in10m },
    ])
    expect(status.hasLive).toBe(false)
    expect(status.nextKickoffInMinutes).not.toBeNull()
    expect(status.nextKickoffInMinutes!).toBeGreaterThan(9)
    expect(status.nextKickoffInMinutes!).toBeLessThan(11)
  })

  it('ignores finished/cancelled and past kickoffs', () => {
    const past = new Date(Date.now() - 60_000).toISOString()
    const status = pollingStatusFrom([
      { status: 'finished', kickoff_time_utc: '2099-01-01T00:00:00Z' },
      { status: 'cancelled', kickoff_time_utc: '2099-01-01T00:00:00Z' },
      { status: 'scheduled', kickoff_time_utc: past },
    ])
    expect(status).toEqual({ hasLive: false, nextKickoffInMinutes: null })
  })

  it('handles empty lists', () => {
    expect(pollingStatusFrom([])).toEqual({ hasLive: false, nextKickoffInMinutes: null })
  })
})
