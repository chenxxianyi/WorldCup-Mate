import request from './request'
import type { ApiPlayer } from '@/types/player'

export function apiListTeams(params?: Record<string, any>) {
  return request.get('/teams', { params }) as Promise<any>
}

export function apiGetTeamDetail(id: number) {
  return request.get(`/teams/${id}`) as Promise<any>
}

export function apiGetTeamMatches(id: number) {
  return request.get(`/teams/${id}/matches`) as Promise<any[]>
}

export function apiGetTeamPlayers(id: number) {
  return request.get(`/teams/${id}/players`) as Promise<ApiPlayer[]>
}
