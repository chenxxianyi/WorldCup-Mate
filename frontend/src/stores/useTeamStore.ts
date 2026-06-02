import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Team } from '@/types/team'
import { normalizeTeam } from '@/types/team'
import { apiListTeams, apiGetTeamDetail } from '@/api/teams'

export const useTeamStore = defineStore('team', () => {
  const teams = ref<Team[]>([])
  const loading = ref(false)

  async function fetchTeams(params?: Record<string, any>) {
    loading.value = true
    try {
      const res = await apiListTeams(params) as any
      teams.value = (res.list || res).map(normalizeTeam)
    } finally {
      loading.value = false
    }
  }

  async function fetchTeamDetail(id: number) {
    try {
      const res = await apiGetTeamDetail(id) as any
      return normalizeTeam(res)
    } catch {
      return null
    }
  }

  function getTeamById(id: number) {
    return teams.value.find((t) => t.id === id) || null
  }

  return { teams, loading, fetchTeams, fetchTeamDetail, getTeamById }
})
