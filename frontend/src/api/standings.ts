import request from './request'
import type { ApiStanding } from '@/types/standing'

export function apiGetAllStandings() {
  return request.get('/standings') as Promise<ApiStanding[]>
}

export function apiGetGroupStandings(groupId: number) {
  return request.get(`/groups/${groupId}/standings`) as Promise<ApiStanding[]>
}

export function apiGetBestThird() {
  return request.get('/standings/best-third') as Promise<ApiStanding[]>
}
