// @vitest-environment jsdom
import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { createMemoryHistory, createRouter } from 'vue-router'
import { createPinia } from 'pinia'
import MatchCard from '@/components/common/MatchCard.vue'
import type { Match } from '@/types/match'

const router = createRouter({ history: createMemoryHistory(), routes: [] })

function makeMatch(overrides: Partial<Match> = {}): Match {
  return {
    id: 1,
    match_number: 1,
    stage: 'group',
    group_name: null,
    home_team_id: 10,
    away_team_id: 20,
    home_team_name: '阿根廷',
    away_team_name: '法国',
    home_team_code: 'ARG',
    away_team_code: 'FRA',
    home_flag: '',
    away_flag: '',
    home_score: null,
    away_score: null,
    status: 'scheduled',
    kickoff_time_utc: '2026-09-02T19:00:00Z',
    local_kickoff_time: '09-03 03:00',
    city: '',
    stadium: '',
    importance_level: 1,
    is_featured: false,
    minute: null,
    matchday: null,
    competition_id: null,
    ...overrides,
  }
}

function mountCard(match: Match) {
  return mount(MatchCard, {
    props: { match },
    global: {
      plugins: [createPinia(), router],
      stubs: { ReminderControl: true, RouterLink: true },
    },
  })
}

describe('MatchCard status mapping (DATA-04)', () => {
  it('renders the live minute when present', () => {
    const w = mountCard(makeMatch({ status: 'live', minute: 67, home_score: 1, away_score: 0 }))
    expect(w.text()).toContain("67'")
    expect(w.text()).not.toContain("0'")
  })

  it('renders 直播中 instead of a fake 0 minute when the clock is missing', () => {
    const w = mountCard(makeMatch({ status: 'live', minute: null }))
    expect(w.text()).toContain('直播中')
    expect(w.text()).not.toMatch(/0'\s*LIVE/)
    expect(w.text()).not.toContain("进行中 · 0'")
  })

  it('renders postponed and cancelled labels', () => {
    expect(mountCard(makeMatch({ status: 'postponed' })).text()).toContain('已延期')
    expect(mountCard(makeMatch({ status: 'cancelled' })).text()).toContain('已取消')
  })

  it('renders the final score for finished matches', () => {
    const w = mountCard(makeMatch({ status: 'finished', home_score: 3, away_score: 1 }))
    expect(w.text()).toContain('已结束')
    expect(w.text()).toContain('3')
    expect(w.text()).toContain('1')
  })

  it('shows 时间待定 when no kickoff time is available', () => {
    const w = mountCard(makeMatch({ local_kickoff_time: '' }))
    expect(w.text()).toContain('时间待定')
  })
})
