import request from './request'

export function apiListTeams(params?: Record<string, any>) {
  return request.get('/teams', { params }) as Promise<any>
}

export function apiGetTeamDetail(id: number) {
  return request.get(`/teams/${id}`) as Promise<any>
}

export function apiGetTeamMatches(id: number) {
  return request.get(`/teams/${id}/matches`) as Promise<any[]>
}
