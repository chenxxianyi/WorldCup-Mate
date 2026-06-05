export type LineupStatus = 'unavailable' | 'partial' | 'available'
export type LineupSide = 'home' | 'away'
export type LineupRole = 'start_xi' | 'substitute'

export interface ApiLineupPlayer {
  player_id?: number | null
  name?: string
  name_en?: string
  shirt_number?: number
  position?: string
  position_label?: string
  role?: LineupRole
  grid?: string
  photo_url?: string
}

export interface ApiTeamLineup {
  team_id: number
  team_name: string
  team_code?: string
  side?: LineupSide
  formation?: string
  coach_name?: string
  status?: LineupStatus
  start_xi?: ApiLineupPlayer[]
  substitutes?: ApiLineupPlayer[]
}

export interface ApiMatchLineups {
  match_id: number
  status?: LineupStatus
  source?: 'football-data' | 'api-football' | string
  message?: string
  home?: ApiTeamLineup | null
  away?: ApiTeamLineup | null
}

export interface LineupPlayer {
  playerId: number | null
  name: string
  nameEn: string
  shirtNumber?: number
  position: string
  positionLabel: string
  role: LineupRole
  grid: string
  photoUrl: string
}

export interface TeamLineup {
  teamId: number
  teamName: string
  teamCode: string
  side: LineupSide
  formation: string
  coachName: string
  status: LineupStatus
  startXi: LineupPlayer[]
  substitutes: LineupPlayer[]
}

export interface MatchLineups {
  matchId: number
  status: LineupStatus
  source: string
  message: string
  home: TeamLineup | null
  away: TeamLineup | null
}

export function normalizeLineupPlayer(value: ApiLineupPlayer, role: LineupRole): LineupPlayer {
  return {
    playerId: value.player_id ?? null,
    name: value.name || '',
    nameEn: value.name_en || '',
    shirtNumber: value.shirt_number,
    position: (value.position || '').toUpperCase(),
    positionLabel: value.position_label || '',
    role,
    grid: value.grid || '',
    photoUrl: value.photo_url || '',
  }
}

export function normalizeTeamLineup(value: ApiTeamLineup | null | undefined, side: LineupSide): TeamLineup | null {
  if (!value) return null

  return {
    teamId: value.team_id,
    teamName: value.team_name || '',
    teamCode: value.team_code || '',
    side: value.side || side,
    formation: value.formation || '',
    coachName: value.coach_name || '',
    status: value.status || 'unavailable',
    startXi: (value.start_xi || []).map((player) => normalizeLineupPlayer(player, 'start_xi')),
    substitutes: (value.substitutes || []).map((player) => normalizeLineupPlayer(player, 'substitute')),
  }
}

export function normalizeMatchLineups(value: ApiMatchLineups): MatchLineups {
  return {
    matchId: value.match_id,
    status: value.status || 'unavailable',
    source: value.source || '',
    message: value.message || '',
    home: normalizeTeamLineup(value.home, 'home'),
    away: normalizeTeamLineup(value.away, 'away'),
  }
}
