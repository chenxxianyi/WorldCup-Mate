import request from './request'
import type { Competition } from '@/types/competition'

export function apiListCompetitions() {
  return request.get('/competitions') as Promise<Competition[]>
}

export function apiGetCompetitionStandings(code: string, params?: Record<string, any>) {
  return request.get(`/competitions/${encodeURIComponent(code)}/standings`, { params }) as Promise<any[]>
}
