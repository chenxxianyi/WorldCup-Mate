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
      const res = await apiListFavoriteTeams()
      followedTeamIds.value = res.map((f) => f.team_id)
    } catch {
      followedTeamIds.value = []
    }
  }

  async function fetchFavoriteMatches() {
    try {
      const res = await apiListFavoriteMatches()
      favoriteMatchIds.value = res.map((f) => f.match_id)
      favoriteMatches.value = res.filter((f) => f.match).map((f) => normalizeMatch(f.match!))
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
      return true
    } catch {
      if (wasFollowed) followedTeamIds.value.push(teamId)
      else followedTeamIds.value = followedTeamIds.value.filter((id) => id !== teamId)
      return false
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
      return true
    } catch {
      if (wasFav) favoriteMatchIds.value.push(matchId)
      else favoriteMatchIds.value = favoriteMatchIds.value.filter((id) => id !== matchId)
      return false
    }
  }

  return {
    followedTeamIds, favoriteMatchIds, favoriteMatches,
    isTeamFollowed, isMatchFavorite, fetchFavoriteTeams, fetchFavoriteMatches,
    toggleTeamFollow, toggleMatchFavorite,
  }
})
