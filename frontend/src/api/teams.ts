import request from './request'
import type { ApiTeam } from '@/types/team'
import type { ApiMatch } from '@/types/match'
import type { PaginatedData } from '@/types/common'

export function apiListTeams(params?: Record<string, any>) {
  return request.get('/teams', { params }) as Promise<PaginatedData<ApiTeam>>
}

export function apiGetTeamDetail(id: number) {
  return request.get(`/teams/${id}`) as Promise<ApiTeam>
}

export function apiGetTeamMatches(id: number) {
  return request.get(`/teams/${id}/matches`) as Promise<ApiMatch[]>
}
