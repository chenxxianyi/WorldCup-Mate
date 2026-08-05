// @vitest-environment jsdom
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useCompetitionStore } from '@/stores/useCompetitionStore'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import type { Competition } from '@/types/competition'

vi.mock('@/api/competitions', () => ({
  apiListCompetitions: vi.fn(),
  apiGetCompetitionStandings: vi.fn(),
}))

import { apiListCompetitions } from '@/api/competitions'

function comp(overrides: Partial<Competition>): Competition {
  return {
    id: 1, code: 'PL', name: '英超', name_en: 'Premier League', country: 'England',
    logo_url: '', format: 'league', season: 2026, status: 'active', sort_order: 1, ...overrides,
  }
}

describe('competition config driven switcher', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
    vi.mocked(apiListCompetitions).mockReset()
    // jsdom lacks matchMedia; useSettingStore reads it at init.
    window.matchMedia = vi.fn().mockReturnValue({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn() }) as any
  })

  it('hides the disabled World Cup from the switcher options', async () => {
    // WC disabled in the admin panel, PL enabled.
    vi.mocked(apiListCompetitions).mockResolvedValue([
      comp({ code: 'WC', name: '世界杯', format: 'cup', status: 'inactive', sort_order: 0 }),
      comp({ code: 'PL', name: '英超', format: 'league', status: 'active', sort_order: 1 }),
    ])
    const store = useCompetitionStore()
    await store.fetchCompetitions()

    const theme = useLeagueThemeStore()
    expect(theme.competitionCodes).toEqual(['PL'])
    expect(theme.isCompetitionCode('WC')).toBe(false)
    expect(theme.isCompetitionCode('PL')).toBe(true)
  })

  it('falls back to the first enabled league when the stored selection was disabled', async () => {
    localStorage.setItem('wm-competition', 'WC')
    vi.mocked(apiListCompetitions).mockResolvedValue([
      comp({ code: 'PL', name: '英超', status: 'active', sort_order: 1 }),
      comp({ code: 'PD', name: '西甲', status: 'active', sort_order: 2 }),
    ])
    const store = useCompetitionStore()
    await store.fetchCompetitions()

    expect(store.currentCode).toBe('PL')
    expect(localStorage.getItem('wm-competition')).toBe('PL')
  })

  it('keeps the World Cup when it is enabled', async () => {
    vi.mocked(apiListCompetitions).mockResolvedValue([
      comp({ code: 'WC', name: '世界杯', format: 'cup', status: 'active', sort_order: 0 }),
      comp({ code: 'PL', name: '英超', status: 'active', sort_order: 1 }),
    ])
    const store = useCompetitionStore()
    await store.fetchCompetitions()

    const theme = useLeagueThemeStore()
    expect(theme.competitionCodes).toEqual(['WC', 'PL'])
    expect(theme.isCompetitionCode('WC')).toBe(true)
  })

  it('leaves the switcher empty when every competition is disabled (no stale WC fallback)', async () => {
    // A stale non-WC selection must be preserved, not rewritten to WC.
    localStorage.setItem('wm-competition', 'SA')
    vi.mocked(apiListCompetitions).mockResolvedValue([
      comp({ code: 'WC', name: '世界杯', format: 'cup', status: 'inactive', sort_order: 0 }),
      comp({ code: 'PL', name: '英超', status: 'inactive', sort_order: 1 }),
    ])
    const store = useCompetitionStore()
    await store.fetchCompetitions()

    const theme = useLeagueThemeStore()
    expect(theme.competitionCodes).toEqual([])
    expect(theme.isCompetitionCode('WC')).toBe(false)
    // The stored selection is left untouched — no stale WC written back.
    expect(localStorage.getItem('wm-competition')).toBe('SA')
  })

  it('syncs the theme store when the fallback rewrites the selection', async () => {
    localStorage.setItem('wm-competition', 'SA')
    vi.mocked(apiListCompetitions).mockResolvedValue([
      comp({ code: 'PL', name: '英超', status: 'active', sort_order: 1 }),
    ])
    const store = useCompetitionStore()
    const theme = useLeagueThemeStore()
    await store.fetchCompetitions()

    // The disabled SA selection falls back to PL and the theme store
    // follows via the watch on competition.currentCode.
    expect(store.currentCode).toBe('PL')
    expect(theme.currentCode).toBe('PL')
  })
})
