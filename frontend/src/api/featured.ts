import request from './request'
import type { ApiMatch } from '@/types/match'

export interface FeaturedConfig {
  id: number
  competition_id: number
  match_id: number | null
  tagline: string
  description: string
  stage_label: string
  enabled: boolean
  updated_at: string
}

export interface FeaturedMatchPick {
  id: number
  home: string
  away: string
  kickoff_time_utc: string
  status: string
}

export interface FeaturedInput {
  match_id: number | null
  tagline: string
  description: string
  stage_label: string
  enabled: boolean
}

// Public: enabled configs keyed by competition code.
export function apiGetFeatured() {
  return request.get('/featured') as Promise<Record<string, FeaturedConfig>>
}

export function apiAdminListFeatured() {
  return request.get('/admin/featured') as Promise<FeaturedConfig[]>
}

export function apiAdminFeaturedMatches(code: string) {
  return request.get(`/admin/featured/${encodeURIComponent(code)}/matches`) as Promise<FeaturedMatchPick[]>
}

export function apiAdminUpdateFeatured(code: string, data: FeaturedInput) {
  return request.put(`/admin/featured/${encodeURIComponent(code)}`, data) as Promise<FeaturedConfig>
}
