import request from './request'
import type { ApiMatchLineups } from '@/types/lineup'

export function apiGetMatchLineups(matchId: number) {
  return request.get(`/matches/${matchId}/lineups`) as Promise<ApiMatchLineups>
}
