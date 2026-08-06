import request from './request'
import type { Competition } from '@/types/competition'
import type { ApiLeagueStanding } from '@/types/standing'

export function apiListCompetitions() {
  return request.get('/competitions') as Promise<Competition[]>
}

export interface CompetitionOverview {
  competition: Competition
  seasons: number[]
  season: number
  matchday: number | null
  match_count: number
}

export function apiCompetitionOverview(code: string, season?: number) {
  const params = season ? { season } : undefined
  return request.get(`/competitions/${encodeURIComponent(code)}/overview`, { params }) as Promise<CompetitionOverview>
}

export function apiGetCompetitionStandings(code: string, params?: Record<string, any>) {
  return request.get(`/competitions/${encodeURIComponent(code)}/standings`, { params }) as Promise<ApiLeagueStanding[]>
}
