import request from './request'
import type { PostMatchSummary } from '@/types/postMatchSummary'

export function apiGetPostMatchSummary(matchId: number) {
  return request.get(`/matches/${matchId}/post-match-summary`) as Promise<PostMatchSummary | any>
}

export function apiGeneratePostMatchSummary(matchId: number, forceRefresh = false) {
  return request.post(`/matches/${matchId}/post-match-summary/generate`, {
    force_refresh: forceRefresh,
  }) as Promise<PostMatchSummary>
}
