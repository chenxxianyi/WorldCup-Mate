import request from './request'
import type { FavoriteTeamRecord, FavoriteMatchRecord } from '@/types/favorite'

export function apiAddFavoriteTeam(teamId: number) {
  return request.post(`/favorites/teams/${teamId}`) as Promise<{ added: boolean }>
}

export function apiRemoveFavoriteTeam(teamId: number) {
  return request.delete(`/favorites/teams/${teamId}`) as Promise<{ removed: boolean }>
}

export function apiListFavoriteTeams() {
  return request.get('/favorites/teams') as Promise<FavoriteTeamRecord[]>
}

export function apiAddFavoriteMatch(matchId: number) {
  return request.post(`/favorites/matches/${matchId}`) as Promise<{ added: boolean }>
}

export function apiRemoveFavoriteMatch(matchId: number) {
  return request.delete(`/favorites/matches/${matchId}`) as Promise<{ removed: boolean }>
}

export function apiListFavoriteMatches() {
  return request.get('/favorites/matches') as Promise<FavoriteMatchRecord[]>
}
