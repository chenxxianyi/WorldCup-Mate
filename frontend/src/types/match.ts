import { formatKickoff } from '@/utils/datetime'

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
  | 'regular_season'

/** Backend raw match object */
export interface ApiMatch {
  id: number
  match_no: number
  stage: Stage
  group_id: number | null
  group: { id: number; name: string } | null
  competition_id?: number | null
  season?: number | null
  matchday?: number | null
  home_team_id: number
  away_team_id: number
  home_team: { id: number; name: string; name_en: string; fifa_code: string | null; external_code?: string | null; flag_url: string }
  away_team: { id: number; name: string; name_en: string; fifa_code: string | null; external_code?: string | null; flag_url: string }
  home_score: number | null
  away_score: number | null
  status: MatchStatus
  live_minute?: number | null
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
  group_id?: number | null
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
  matchday: number | null
  competition_id: number | null
}

/** Transform backend match to frontend match */
export function normalizeMatch(m: ApiMatch): Match {
  return {
    id: m.id,
    match_number: m.match_no,
    stage: m.stage,
    group_id: m.group_id ?? null,
    group_name: m.group?.name || null,
    home_team_id: m.home_team_id,
    away_team_id: m.away_team_id,
    home_team_name: m.home_team?.name || '',
    away_team_name: m.away_team?.name || '',
    home_team_code: m.home_team?.fifa_code || m.home_team?.external_code || '',
    away_team_code: m.away_team?.fifa_code || m.away_team?.external_code || '',
    home_flag: m.home_team?.flag_url || '',
    away_flag: m.away_team?.flag_url || '',
    home_score: m.home_score,
    away_score: m.away_score,
    status: m.status,
    kickoff_time_utc: m.kickoff_time_utc,
    local_kickoff_time: formatKickoff(m.kickoff_time_utc),
    city: m.city?.name || '',
    stadium: m.stadium?.name || '',
    importance_level: m.importance_level,
    is_featured: m.importance_level >= 2,
    minute: m.live_minute ?? null,
    matchday: m.matchday ?? null,
    competition_id: m.competition_id ?? null,
  }
}

// Unified match status mapping (DATA-04): every card/detail page must use
// these helpers so postponed/cancelled/live-minute semantics stay consistent.
export function matchStatusLabel(status: string, minute: number | null): string {
  if (status === 'live') {
    // minute may be null when the provider has no live clock data:
    // never render a fake "0'".
    return minute != null ? `${minute}'` : '直播中'
  }
  if (status === 'finished') return '已结束'
  if (status === 'postponed') return '已延期'
  if (status === 'cancelled') return '已取消'
  return '未开始'
}

export function matchStatusClass(status: string): string {
  if (status === 'live') return 'live'
  if (status === 'finished') return 'green'
  if (status === 'postponed' || status === 'cancelled') return 'warn'
  return ''
}
