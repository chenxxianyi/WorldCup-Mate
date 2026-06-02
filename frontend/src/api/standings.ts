import request from './request'

export function apiGetAllStandings() {
  return request.get('/standings') as Promise<any[]>
}

export function apiGetGroupStandings(groupId: number) {
  return request.get(`/groups/${groupId}/standings`) as Promise<any[]>
}

export function apiGetBestThird() {
  return request.get('/standings/best-third') as Promise<any[]>
}
