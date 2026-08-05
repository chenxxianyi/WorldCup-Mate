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
