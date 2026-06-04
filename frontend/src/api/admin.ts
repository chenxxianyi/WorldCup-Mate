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
