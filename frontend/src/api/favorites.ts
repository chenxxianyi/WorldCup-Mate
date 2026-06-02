import request from './request'

export function apiAddFavoriteTeam(teamId: number) {
  return request.post(`/favorites/teams/${teamId}`) as Promise<any>
}

export function apiRemoveFavoriteTeam(teamId: number) {
  return request.delete(`/favorites/teams/${teamId}`) as Promise<any>
}

export function apiListFavoriteTeams() {
  return request.get('/favorites/teams') as Promise<any[]>
}

export function apiAddFavoriteMatch(matchId: number) {
  return request.post(`/favorites/matches/${matchId}`) as Promise<any>
}

export function apiRemoveFavoriteMatch(matchId: number) {
  return request.delete(`/favorites/matches/${matchId}`) as Promise<any>
}

export function apiListFavoriteMatches() {
  return request.get('/favorites/matches') as Promise<any[]>
}
