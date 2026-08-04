import request from './request'
import type { ApiMatch } from '@/types/match'
import type { PaginatedData } from '@/types/common'

export function apiListMatches(params?: Record<string, any>) {
  return request.get('/matches', { params }) as Promise<PaginatedData<ApiMatch>>
}

export function apiGetTodayMatches() {
  return request.get('/matches/today') as Promise<ApiMatch[]>
}

export function apiGetTomorrowMatches() {
  return request.get('/matches/tomorrow') as Promise<ApiMatch[]>
}

export function apiGetUpcomingMatches() {
  return request.get('/matches/upcoming') as Promise<ApiMatch[]>
}

export function apiGetRecommendedMatches() {
  return request.get('/matches/recommended') as Promise<ApiMatch[]>
}

export function apiGetMatchDetail(id: number) {
  return request.get(`/matches/${id}`) as Promise<ApiMatch>
}

export function apiGetMatchesByTeam(teamId: number) {
  return request.get(`/matches/by-team/${teamId}`) as Promise<ApiMatch[]>
}

export function apiGetMatchesByGroup(groupId: number) {
  return request.get(`/matches/by-group/${groupId}`) as Promise<ApiMatch[]>
}
