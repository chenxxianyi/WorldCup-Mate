import type { ApiMatch } from './match'
import type { ApiTeam } from './team'

/** Backend favorite-team record (UserFavoriteTeam JSON). */
export interface FavoriteTeamRecord {
  id: number
  user_id: number
  team_id: number
  team?: ApiTeam
  created_at: string
}

/** Backend favorite-match record (UserFavoriteMatch JSON). */
export interface FavoriteMatchRecord {
  id: number
  user_id: number
  match_id: number
  match?: ApiMatch
  created_at: string
}
