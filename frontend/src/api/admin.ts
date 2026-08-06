import request from './request'
import type { ApiTeam } from '@/types/team'
import type { ApiMatch } from '@/types/match'
import type { ApiStanding } from '@/types/standing'
import type { LeagueSyncResult } from '@/types/sync'
import type { PaginatedData } from '@/types/common'

export interface AdminDashboardStats {
  total_matches: number
  total_groups: number
  total_teams: number
  total_users: number
  total_competitions: number
  total_reminders: number
}

export function apiAdminDashboard() {
  return request.get('/admin/dashboard') as Promise<AdminDashboardStats>
}

export function apiAdminListTeams(params?: Record<string, any>) {
  return request.get('/admin/teams', { params }) as Promise<PaginatedData<ApiTeam>>
}

export function apiAdminListMatches(params?: Record<string, any>) {
  return request.get('/admin/matches', { params }) as Promise<PaginatedData<ApiMatch>>
}

export function apiAdminListStandings(params?: Record<string, any>) {
  return request.get('/admin/standings', { params }) as Promise<ApiStanding[]>
}

export function apiAdminRecalculateStandings() {
  return request.post('/admin/standings/recalculate') as Promise<{ recalculated: boolean }>
}

export function apiAdminSyncMatches(params?: Record<string, any>) {
  return request.post('/admin/sync/matches', null, { params }) as Promise<LeagueSyncResult>
}

export interface LeagueRecalculateResult {
  recalculated: boolean
  competition_id: number
  season: number
}

export function apiAdminRecalculateLeagueStanding(data: Record<string, any>) {
  return request.post('/admin/standings/league/recalculate', data) as Promise<LeagueRecalculateResult>
}

export interface AdminCompetition {
  id: number
  code: string
  name: string
  name_en: string
  country: string
  logo_url: string
  format: string
  season: number
  status: string
  sort_order: number
}

export interface CompetitionInput {
  code?: string
  name?: string
  name_en?: string
  country?: string
  logo_url?: string
  format?: string
  season?: number
  status?: string
  sort_order?: number
}

export function apiAdminListCompetitions() {
  return request.get('/admin/competitions') as Promise<AdminCompetition[]>
}

export function apiAdminCreateCompetition(data: CompetitionInput) {
  return request.post('/admin/competitions', data) as Promise<AdminCompetition>
}

export function apiAdminUpdateCompetition(id: number, data: CompetitionInput) {
  return request.put(`/admin/competitions/${id}`, data) as Promise<AdminCompetition>
}

// ---------- Team CRUD (ADM-11) ----------

export interface TeamInput {
  name?: string
  name_en?: string
  fifa_code?: string
  external_code?: string
  team_type?: string
  flag_url?: string
  continent?: string
  country?: string
  venue?: string
  group_id?: number
  coach?: string
  description?: string
}

export function apiAdminCreateTeam(data: TeamInput) {
  return request.post('/admin/teams', data) as Promise<ApiTeam>
}

export function apiAdminUpdateTeam(id: number, data: TeamInput) {
  return request.put(`/admin/teams/${id}`, data) as Promise<ApiTeam>
}

export function apiAdminDeleteTeam(id: number) {
  return request.delete(`/admin/teams/${id}`) as Promise<{ deleted: boolean }>
}

// ---------- Match CRUD (ADM-12) ----------

export interface MatchInput {
  match_no?: number
  home_team_id?: number
  away_team_id?: number
  group_id?: number
  stage?: string
  stadium_id?: number
  city_id?: number
  kickoff_time_utc?: string
  importance_level?: number
  recommend_tag?: string
  competition_id?: number
  season?: number
  matchday?: number
}

export interface ScoreInput {
  home_score: number
  away_score: number
}

export function apiAdminCreateMatch(data: MatchInput) {
  return request.post('/admin/matches', data) as Promise<ApiMatch>
}

export function apiAdminUpdateMatch(id: number, data: MatchInput) {
  return request.put(`/admin/matches/${id}`, data) as Promise<ApiMatch>
}

export function apiAdminDeleteMatch(id: number) {
  return request.delete(`/admin/matches/${id}`) as Promise<{ deleted: boolean }>
}

export function apiAdminUpdateMatchScore(id: number, data: ScoreInput) {
  return request.put(`/admin/matches/${id}/score`, data) as Promise<ApiMatch>
}

export function apiAdminUpdateMatchStatus(id: number, data: { status: string }) {
  return request.put(`/admin/matches/${id}/status`, data) as Promise<ApiMatch>
}

// ---------- Sync & Reminder operations (ADM-13) ----------

export interface SyncHistoryItem {
  id: number
  provider: string
  resource: string
  status: string
  started_at: string
  finished_at: string
  total: number
  created: number
  updated: number
  skipped: number
  error_message: string
  reason: string
}

export function apiAdminSyncHistory(params?: Record<string, any>) {
  return request.get('/admin/sync/history', { params }) as Promise<PaginatedData<SyncHistoryItem>>
}

export interface AdminReminderItem {
  id: number
  user_id: number
  match_id: number
  remind_at: string
  channel: string
  status: string
  retry_count: number
  last_error: string
  next_retry_at: string
  match?: ApiMatch
}

export function apiAdminListReminders(params?: Record<string, any>) {
  return request.get('/admin/reminders', { params }) as Promise<PaginatedData<AdminReminderItem>>
}

export function apiAdminRetryReminder(id: number) {
  return request.post(`/admin/reminders/${id}/retry`) as Promise<{ retried: boolean }>
}
