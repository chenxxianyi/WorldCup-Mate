import request from './request'
import type { ApiMatch } from '@/types/match'
import type { PaginatedData } from '@/types/common'
import type { AxiosRequestConfig } from 'axios'

export function apiListMatches(params?: Record<string, any>, config?: AxiosRequestConfig) {
  return request.get('/matches', { params, ...config }) as Promise<PaginatedData<ApiMatch>>
}

export function apiGetTodayMatches(config?: AxiosRequestConfig) {
  return request.get('/matches/today', config) as Promise<ApiMatch[]>
}

export function apiGetTomorrowMatches(config?: AxiosRequestConfig) {
  return request.get('/matches/tomorrow', config) as Promise<ApiMatch[]>
}

export function apiGetUpcomingMatches(config?: AxiosRequestConfig) {
  return request.get('/matches/upcoming', config) as Promise<ApiMatch[]>
}

export function apiGetRecommendedMatches(config?: AxiosRequestConfig) {
  return request.get('/matches/recommended', config) as Promise<ApiMatch[]>
}

export function apiGetMatchDetail(id: number, config?: AxiosRequestConfig) {
  return request.get(`/matches/${id}`, config) as Promise<ApiMatch>
}

export function apiGetMatchesByTeam(teamId: number, config?: AxiosRequestConfig) {
  return request.get(`/matches/by-team/${teamId}`, config) as Promise<ApiMatch[]>
}

export function apiGetMatchesByGroup(groupId: number, config?: AxiosRequestConfig) {
  return request.get(`/matches/by-group/${groupId}`, config) as Promise<ApiMatch[]>
}
