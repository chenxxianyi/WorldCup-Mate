import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  apiAddFavoriteTeam,
  apiRemoveFavoriteTeam,
  apiListFavoriteTeams,
  apiAddFavoriteMatch,
  apiRemoveFavoriteMatch,
  apiListFavoriteMatches,
} from '@/api/favorites'
import { normalizeMatch, type Match } from '@/types/match'

export const useFavoriteStore = defineStore('favorite', () => {
  const followedTeamIds = ref<number[]>([])
  const favoriteMatchIds = ref<number[]>([])
  const favoriteMatches = ref<Match[]>([])

  function isTeamFollowed(teamId: number) {
    return followedTeamIds.value.includes(teamId)
  }

  function isMatchFavorite(matchId: number) {
    return favoriteMatchIds.value.includes(matchId)
  }

  async function fetchFavoriteTeams() {
    try {
      const res = await apiListFavoriteTeams() as any[]
      followedTeamIds.value = res.map((f: any) => f.team_id)
    } catch {
      followedTeamIds.value = []
    }
  }

  async function fetchFavoriteMatches() {
    try {
      const res = await apiListFavoriteMatches() as any[]
      favoriteMatchIds.value = res.map((f: any) => f.match_id)
      favoriteMatches.value = res.filter((f: any) => f.match).map((f: any) => normalizeMatch(f.match))
    } catch {
      favoriteMatchIds.value = []
      favoriteMatches.value = []
    }
  }

  async function toggleTeamFollow(teamId: number) {
    const wasFollowed = followedTeamIds.value.includes(teamId)
    if (wasFollowed) {
      followedTeamIds.value = followedTeamIds.value.filter((id) => id !== teamId)
    } else {
      followedTeamIds.value.push(teamId)
    }
    try {
      if (wasFollowed) {
        await apiRemoveFavoriteTeam(teamId)
      } else {
        await apiAddFavoriteTeam(teamId)
      }
    } catch {
      if (wasFollowed) followedTeamIds.value.push(teamId)
      else followedTeamIds.value = followedTeamIds.value.filter((id) => id !== teamId)
    }
  }

  async function toggleMatchFavorite(matchId: number) {
    const wasFav = favoriteMatchIds.value.includes(matchId)
    if (wasFav) {
      favoriteMatchIds.value = favoriteMatchIds.value.filter((id) => id !== matchId)
    } else {
      favoriteMatchIds.value.push(matchId)
    }
    try {
      if (wasFav) {
        await apiRemoveFavoriteMatch(matchId)
      } else {
        await apiAddFavoriteMatch(matchId)
      }
    } catch {
      if (wasFav) favoriteMatchIds.value.push(matchId)
      else favoriteMatchIds.value = favoriteMatchIds.value.filter((id) => id !== matchId)
    }
  }

  return {
    followedTeamIds, favoriteMatchIds, favoriteMatches,
    isTeamFollowed, isMatchFavorite, fetchFavoriteTeams, fetchFavoriteMatches,
    toggleTeamFollow, toggleMatchFavorite,
  }
})
