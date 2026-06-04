export type MatchStatus = 'scheduled' | 'upcoming' | 'live' | 'finished' | 'postponed' | 'cancelled'

export type Stage =
  | 'group'
  | 'group_stage'
  | 'round_of_32'
  | 'round_of_16'
  | 'quarter_final'
  | 'semi_final'
  | 'third_place'
  | 'final'

/** Backend raw match object */
export interface ApiMatch {
  id: number
  match_no: number
  stage: Stage
  group_id: number | null
  group: { id: number; name: string } | null
  home_team_id: number
  away_team_id: number
  home_team: { id: number; name: string; name_en: string; fifa_code: string; flag_url: string }
  away_team: { id: number; name: string; name_en: string; fifa_code: string; flag_url: string }
  home_score: number | null
  away_score: number | null
  status: MatchStatus
  kickoff_time_utc: string
  stadium: { id: number; name: string; city_id: number } | null
  city: { id: number; name: string } | null
  importance_level: number
  recommend_tag: string
  recommend_reason: string
  winner_team_id: number | null
  created_at: string
  updated_at: string
}

/** Frontend normalized match */
export interface Match {
  id: number
  match_number: number
  stage: Stage
  group_name: string | null
  home_team_id: number
  away_team_id: number
  home_team_name: string
  away_team_name: string
  home_team_code: string
  away_team_code: string
  home_flag: string
  away_flag: string
  home_score: number | null
  away_score: number | null
  status: MatchStatus
  kickoff_time_utc: string
  local_kickoff_time: string
  city: string
  stadium: string
  importance_level: number
  is_featured: boolean
  minute: number | null
}

/** Convert UTC time string to Beijing time (UTC+8), format: "MM-DD HH:mm" */
function formatBeijingTime(utcStr: string): string {
  if (!utcStr) return ''
  const d = new Date(utcStr)
  if (isNaN(d.getTime())) return ''

  const parts = new Intl.DateTimeFormat('en-US', {
    timeZone: 'Asia/Shanghai',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hourCycle: 'h23',
  }).formatToParts(d)
  const pick = (type: string) => parts.find((p) => p.type === type)?.value || '00'

  const mm = pick('month')
  const dd = pick('day')
  const hh = pick('hour')
  const mi = pick('minute')
  return `${mm}-${dd} ${hh}:${mi}`
}

/** Transform backend match to frontend match */
export function normalizeMatch(m: ApiMatch): Match {
  return {
    id: m.id,
    match_number: m.match_no,
    stage: m.stage,
    group_name: m.group?.name || null,
    home_team_id: m.home_team_id,
    away_team_id: m.away_team_id,
    home_team_name: m.home_team?.name || '',
    away_team_name: m.away_team?.name || '',
    home_team_code: m.home_team?.fifa_code || '',
    away_team_code: m.away_team?.fifa_code || '',
    home_flag: m.home_team?.flag_url || '',
    away_flag: m.away_team?.flag_url || '',
    home_score: m.home_score,
    away_score: m.away_score,
    status: m.status,
    kickoff_time_utc: m.kickoff_time_utc,
    local_kickoff_time: formatBeijingTime(m.kickoff_time_utc),
    city: m.city?.name || '',
    stadium: m.stadium?.name || '',
    importance_level: m.importance_level,
    is_featured: m.importance_level >= 2,
    minute: null,
  }
}
