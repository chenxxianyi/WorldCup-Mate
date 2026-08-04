import type { DemoMatch, DemoTeam } from '@/data/leagueTheme'
import type { Match, MatchStatus } from '@/types/match'
import type { Team } from '@/types/team'

const badgeColors = ['#9f224e', '#0068b5', '#009c3b', '#d20515', '#5b2c83', '#f58220', '#173a84', '#670e36']

function hash(value: string) {
  let result = 0
  for (let index = 0; index < value.length; index += 1) result = ((result << 5) - result + value.charCodeAt(index)) | 0
  return Math.abs(result)
}

export function badgeColor(code: string) {
  return badgeColors[hash(code || 'WM') % badgeColors.length]
}

export function teamToThemeTeam(team: Pick<Team, 'name' | 'code' | 'country' | 'continent' | 'venue' | 'flag'>): DemoTeam {
  const code = team.code || team.name.slice(0, 3).toUpperCase()
  return [team.name, code, team.venue || team.country || team.continent || '', badgeColor(code), team.flag]
}

function matchTeam(name: string, code: string, city: string, crest: string): DemoTeam {
  const normalizedCode = code || name.slice(0, 3).toUpperCase()
  return [name || '待定', normalizedCode || 'TBD', city, badgeColor(normalizedCode), crest]
}

function demoStatus(status: MatchStatus): DemoMatch['status'] {
  if (status === 'live') return 'live'
  if (status === 'finished') return 'finished'
  return 'scheduled'
}

function matchScore(match: Match) {
  if (match.status !== 'live' && match.status !== 'finished') return 'VS'
  return `${match.home_score ?? 0}–${match.away_score ?? 0}`
}

function dateParts(value: string) {
  const date = value ? new Date(value) : null
  if (!date || Number.isNaN(date.getTime())) return { date: '时间待定', time: '--:--', key: 'unknown' }
  const dateText = new Intl.DateTimeFormat('zh-CN', { month: '2-digit', day: '2-digit' }).format(date).replace('/', '月') + '日'
  const time = new Intl.DateTimeFormat('zh-CN', { hour: '2-digit', minute: '2-digit', hourCycle: 'h23' }).format(date)
  const key = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  return { date: dateText, time, key }
}

export interface ThemeMatch extends DemoMatch {
  homeTeamId: number
  awayTeamId: number
  kickoffKey: string
  source: Match
}

export function matchToThemeMatch(match: Match): ThemeMatch {
  const parts = dateParts(match.kickoff_time_utc)
  return {
    id: match.id,
    home: matchTeam(match.home_team_name, match.home_team_code, match.city, match.home_flag),
    away: matchTeam(match.away_team_name, match.away_team_code, match.city, match.away_flag),
    homeTeamId: match.home_team_id,
    awayTeamId: match.away_team_id,
    time: match.status === 'live' && match.minute ? `${match.minute}'` : parts.time,
    date: parts.date,
    kickoffKey: parts.key,
    status: demoStatus(match.status),
    score: matchScore(match),
    venue: [match.city, match.stadium].filter(Boolean).join(' · ') || '比赛场地待定',
    featured: match.is_featured,
    source: match,
  }
}

export function formatMatchday(match: Match, fallback: string) {
  if (match.group_name) return match.group_name
  if (match.matchday) return `第 ${match.matchday} 轮`
  return fallback
}

export function localDateLabel(key: string) {
  if (key === 'unknown') return { weekday: '待定', day: '--', month: 'TBD' }
  const date = new Date(`${key}T12:00:00`)
  return {
    weekday: new Intl.DateTimeFormat('zh-CN', { weekday: 'short' }).format(date),
    day: String(date.getDate()).padStart(2, '0'),
    month: new Intl.DateTimeFormat('en-US', { month: 'short' }).format(date).toUpperCase(),
  }
}
