import request from './request'

export function apiAdminListTeams(params?: Record<string, any>) {
  return request.get('/admin/teams', { params }) as Promise<any>
}

export function apiAdminListMatches(params?: Record<string, any>) {
  return request.get('/admin/matches', { params }) as Promise<any>
}

export function apiAdminListStandings(params?: Record<string, any>) {
  return request.get('/admin/standings', { params }) as Promise<any>
}

export function apiAdminRecalculateStandings() {
  return request.post('/admin/standings/recalculate') as Promise<any>
}

export function apiAdminSyncMatches() {
  return request.post('/admin/sync/matches') as Promise<any>
}

export function apiAdminSyncLiveWindowLineups() {
  return request.post('/admin/sync/lineups/live-window') as Promise<any>
}

export function apiAdminSyncMatchLineups(matchId: number) {
  return request.post(`/admin/matches/${matchId}/sync-lineups`) as Promise<any>
}

export function apiAdminGetMatchExternalMapping(matchId: number, provider = 'api-football') {
  return request.get(`/admin/matches/${matchId}/external-match-mapping`, { params: { provider } }) as Promise<any>
}

export function apiAdminUpsertMatchExternalMapping(matchId: number, data: Record<string, any>) {
  return request.put(`/admin/matches/${matchId}/external-match-mapping`, data) as Promise<any>
}
