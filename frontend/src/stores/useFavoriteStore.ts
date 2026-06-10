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

  function clearFavorites() {
    followedTeamIds.value = []
    favoriteMatchIds.value = []
    favoriteMatches.value = []
  }

  function isTeamFollowed(teamId: number) {
    return followedTeamIds.value.includes(Number(teamId))
  }

  function isMatchFavorite(matchId: number) {
    const id = Number(matchId)
    return favoriteMatchIds.value.includes(id) || favoriteMatches.value.some((match) => match.id === id)
  }

  async function fetchFavoriteTeams() {
    try {
      const res = await apiListFavoriteTeams() as any[]
      followedTeamIds.value = res.map((f: any) => Number(f.team_id)).filter(Number.isFinite)
    } catch {
      followedTeamIds.value = []
    }
  }

  async function fetchFavoriteMatches() {
    try {
      const res = await apiListFavoriteMatches() as any[]
      favoriteMatchIds.value = res
        .map((f: any) => Number(f.match_id ?? f.match?.id))
        .filter(Number.isFinite)
      favoriteMatches.value = res.filter((f: any) => f.match).map((f: any) => normalizeMatch(f.match))
    } catch {
      favoriteMatchIds.value = []
      favoriteMatches.value = []
    }
  }

  async function toggleTeamFollow(teamId: number) {
    const id = Number(teamId)
    const wasFollowed = followedTeamIds.value.includes(id)
    if (wasFollowed) {
      followedTeamIds.value = followedTeamIds.value.filter((teamId) => teamId !== id)
    } else {
      followedTeamIds.value.push(id)
    }
    try {
      if (wasFollowed) {
        await apiRemoveFavoriteTeam(id)
      } else {
        await apiAddFavoriteTeam(id)
      }
    } catch {
      if (wasFollowed) followedTeamIds.value.push(id)
      else followedTeamIds.value = followedTeamIds.value.filter((teamId) => teamId !== id)
    }
  }

  async function toggleMatchFavorite(matchId: number, match?: Match) {
    const id = Number(matchId)
    const wasFav = isMatchFavorite(id)
    const prevIds = [...favoriteMatchIds.value]
    const prevMatches = [...favoriteMatches.value]
    if (wasFav) {
      favoriteMatchIds.value = favoriteMatchIds.value.filter((matchId) => matchId !== id)
      favoriteMatches.value = favoriteMatches.value.filter((match) => match.id !== id)
    } else {
      favoriteMatchIds.value.push(id)
      if (match && !favoriteMatches.value.some((item) => item.id === id)) {
        favoriteMatches.value.push(match)
      }
    }
    try {
      if (wasFav) {
        await apiRemoveFavoriteMatch(id)
      } else {
        await apiAddFavoriteMatch(id)
      }
    } catch {
      favoriteMatchIds.value = prevIds
      favoriteMatches.value = prevMatches
    }
  }

  async function addMatchFavorite(matchId: number, match?: Match) {
    const id = Number(matchId)
    if (isMatchFavorite(id)) return
    const prev = [...favoriteMatchIds.value]
    const prevMatches = [...favoriteMatches.value]
    favoriteMatchIds.value.push(id)
    if (match && !favoriteMatches.value.some((item) => item.id === id)) {
      favoriteMatches.value.push(match)
    }
    try {
      await apiAddFavoriteMatch(id)
    } catch {
      favoriteMatchIds.value = prev
      favoriteMatches.value = prevMatches
    }
  }

  return {
    followedTeamIds, favoriteMatchIds, favoriteMatches,
    clearFavorites,
    isTeamFollowed, isMatchFavorite, fetchFavoriteTeams, fetchFavoriteMatches,
    toggleTeamFollow, toggleMatchFavorite, addMatchFavorite,
  }
})
