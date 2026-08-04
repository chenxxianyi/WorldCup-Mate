import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Match } from '@/types/match'
import { normalizeMatch } from '@/types/match'
import {
  apiListMatches,
  apiGetTodayMatches,
  apiGetTomorrowMatches,
  apiGetRecommendedMatches,
  apiGetMatchDetail,
  apiGetMatchesByTeam,
} from '@/api/matches'

export const useMatchStore = defineStore('match', () => {
  const matches = ref<Match[]>([])
  const todayMatches = ref<Match[]>([])
  const tomorrowMatches = ref<Match[]>([])
  const recommendedMatches = ref<Match[]>([])
  const currentMatch = ref<Match | null>(null)
  const loading = ref(false)

  async function fetchMatches(params?: Record<string, any>) {
    loading.value = true
    try {
      const res = await apiListMatches(params)
      matches.value = res.list.map(normalizeMatch)
    } finally {
      loading.value = false
    }
  }

  async function fetchTodayMatches() {
    try {
      const res = await apiGetTodayMatches()
      todayMatches.value = (res || []).map(normalizeMatch)
    } catch {
      todayMatches.value = []
    }
  }

  async function fetchTomorrowMatches() {
    try {
      const res = await apiGetTomorrowMatches()
      tomorrowMatches.value = (res || []).map(normalizeMatch)
    } catch {
      tomorrowMatches.value = []
    }
  }

  async function fetchRecommendedMatches() {
    try {
      const res = await apiGetRecommendedMatches()
      recommendedMatches.value = (res || []).map(normalizeMatch)
    } catch {
      recommendedMatches.value = []
    }
  }

  async function fetchMatchDetail(id: number) {
    loading.value = true
    try {
      const res = await apiGetMatchDetail(id)
      currentMatch.value = normalizeMatch(res)
      return currentMatch.value
    } finally {
      loading.value = false
    }
  }

  async function fetchMatchesByTeam(teamId: number) {
    try {
      const res = await apiGetMatchesByTeam(teamId)
      return (res || []).map(normalizeMatch)
    } catch {
      return []
    }
  }

  function getMatchById(id: number) {
    return matches.value.find((m) => m.id === id) || currentMatch.value
  }

  return {
    matches, todayMatches, tomorrowMatches, recommendedMatches, currentMatch, loading,
    fetchMatches, fetchTodayMatches, fetchTomorrowMatches, fetchRecommendedMatches,
    fetchMatchDetail, fetchMatchesByTeam, getMatchById,
  }
})
