import request from './request'

export function apiListMatches(params?: Record<string, any>) {
  return request.get('/matches', { params }) as Promise<any>
}

export function apiGetTodayMatches() {
  return request.get('/matches/today') as Promise<any[]>
}

export function apiGetTomorrowMatches() {
  return request.get('/matches/tomorrow') as Promise<any[]>
}

export function apiGetUpcomingMatches() {
  return request.get('/matches/upcoming') as Promise<any[]>
}

export function apiGetRecommendedMatches() {
  return request.get('/matches/recommended') as Promise<any[]>
}

export function apiGetMatchDetail(id: number) {
  return request.get(`/matches/${id}`) as Promise<any>
}

export function apiGetMatchesByTeam(teamId: number) {
  return request.get(`/matches/by-team/${teamId}`) as Promise<any[]>
}

export function apiGetMatchesByGroup(groupId: number) {
  return request.get(`/matches/by-group/${groupId}`) as Promise<any[]>
}

export function apiGetTournamentProgress() {
  return request.get('/matches/progress') as Promise<any>
}

export function apiGetTimeline(params?: { days_back?: number; days_ahead?: number }) {
  return request.get('/matches/timeline', { params }) as Promise<any[]>
}

export function apiGetMatchesByStage(stage: string) {
  return request.get(`/matches/by-stage/${stage}`) as Promise<any[]>
}
