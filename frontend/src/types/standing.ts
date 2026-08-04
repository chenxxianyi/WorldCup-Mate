/** Backend raw standing */
export interface ApiStanding {
  id: number
  group_id: number
  group: { id: number; name: string } | null
  team_id: number
  team: { id: number; name: string; fifa_code: string; flag_url: string } | null
  played: number
  won: number
  drawn: number
  lost: number
  goals_for: number
  goals_against: number
  goal_difference: number
  points: number
  rank: number
  qualification_status: string
}

/** Frontend normalized standing */
export interface Standing {
  team_id: number
  team_name: string
  team_code: string
  flag: string
  played: number
  won: number
  drawn: number
  lost: number
  goals_for: number
  goals_against: number
  goal_difference: number
  points: number
  status: '晋级' | '待定' | '淘汰'
}

export function normalizeStanding(s: ApiStanding): Standing {
  let status: '晋级' | '待定' | '淘汰' = '淘汰'
  if (s.qualification_status === 'qualified') status = '晋级'
  else if (s.qualification_status === 'possible') status = '待定'
  return {
    team_id: s.team_id,
    team_name: s.team?.name || '',
    team_code: s.team?.fifa_code || '',
    flag: s.team?.flag_url || '',
    played: s.played,
    won: s.won,
    drawn: s.drawn,
    lost: s.lost,
    goals_for: s.goals_for,
    goals_against: s.goals_against,
    goal_difference: s.goal_difference,
    points: s.points,
    status,
  }
}

export interface GroupStanding {
  group_name: string
  standings: Standing[]
}

/** Backend raw league standing (league_standings table) */
export interface ApiLeagueStanding {
  id: number
  competition_id: number
  season: number
  team_id: number
  team: { id: number; name: string; name_en: string; fifa_code: string | null; flag_url: string } | null
  type: 'total' | 'home' | 'away'
  position: number
  played: number
  won: number
  drawn: number
  lost: number
  goals_for: number
  goals_against: number
  goal_difference: number
  points: number
  zone: string
}

/** Frontend normalized league standing */
export interface LeagueStanding {
  team_id: number
  team_name: string
  team_code: string
  flag: string
  position: number
  played: number
  won: number
  drawn: number
  lost: number
  goals_for: number
  goals_against: number
  goal_difference: number
  points: number
  zone: string
}

export function normalizeLeagueStanding(s: ApiLeagueStanding): LeagueStanding {
  return {
    team_id: s.team_id,
    team_name: s.team?.name || '',
    team_code: s.team?.fifa_code || '',
    flag: s.team?.flag_url || '',
    position: s.position,
    played: s.played,
    won: s.won,
    drawn: s.drawn,
    lost: s.lost,
    goals_for: s.goals_for,
    goals_against: s.goals_against,
    goal_difference: s.goal_difference,
    points: s.points,
    zone: s.zone,
  }
}
