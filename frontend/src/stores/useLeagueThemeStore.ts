import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'
import {
  leagueThemes,
  buildThemeFromCompetition,
  type CompetitionCode,
  type LeagueTheme,
} from '@/data/leagueTheme'
import { WC_CODE } from '@/types/competition'
import { apiGetFeatured, type FeaturedConfig } from '@/api/featured'
import { useCompetitionStore } from '@/stores/useCompetitionStore'
import { useSettingStore } from '@/stores/useSettingStore'

const COMPETITION_KEY = 'wm-competition'

export const useLeagueThemeStore = defineStore('league-theme', () => {
  const stored = localStorage.getItem(COMPETITION_KEY)
  const currentCode = ref<CompetitionCode>(stored || WC_CODE)
  const competitionDialogOpen = ref(false)
  const toast = ref('')
  let toastTimer = 0

  const settings = useSettingStore()
  const competition = useCompetitionStore()

  // Admin-configurable switcher: driven entirely by the backend
  // `competitions` table (active + sort_order). The World Cup is a row
  // like any other (seeded as active) — disabling it hides it too.
  const competitionCodes = computed<CompetitionCode[]>(() =>
    competition.competitions.map((c) => c.code),
  )

  function isCompetitionCode(value: string | null): value is CompetitionCode {
    if (!value) return false
    return competition.competitions.some((c) => c.code === value)
  }

  // Admin-configurable hero copy (featured configs) merged over the static
  // theme identity: tagline/description/stage label come from the backend.
  const featuredConfigs = ref<Record<string, FeaturedConfig>>({})
  const featuredLoaded = ref(false)

  // Theme resolution: static identity for known leagues, generated theme
  // for leagues added through the admin panel (cached per code).
  const dynamicThemes = ref<Record<string, LeagueTheme>>({})

  function themeFor(code: string): LeagueTheme {
    if (leagueThemes[code]) return leagueThemes[code]
    if (dynamicThemes.value[code]) return dynamicThemes.value[code]
    const comp = competition.competitions.find((c) => c.code === code)
    const theme = comp
      ? buildThemeFromCompetition(comp)
      : { ...leagueThemes[WC_CODE], code }
    dynamicThemes.value[code] = theme
    return theme
  }

  // The copy shown on the homepage hero: backend config wins, static theme
  // is the fallback so nothing breaks before the config loads.
  const currentCopy = computed(() => {
    const cfg = featuredConfigs.value[currentCode.value]
    return {
      tagline: cfg?.tagline || current.value.tagline,
      description: cfg?.description || current.value.description,
      stage: cfg?.stage_label || current.value.stage,
      pinnedMatchId: cfg?.match_id ?? null,
    }
  })

  const current = computed(() => themeFor(currentCode.value))
  const progress = computed(() => Math.round((current.value.played / current.value.total) * 100))

  watch(() => competition.currentCode, (code) => {
    if (isCompetitionCode(code) && code !== currentCode.value) {
      currentCode.value = code
      applyDocumentTheme()
    }
  })

  function applyDocumentTheme() {
    document.documentElement.dataset.competition = current.value.slug
    document.documentElement.dataset.theme = settings.theme
    document.title = `WorldCup Mate · ${current.value.name}`
  }

  function setCompetition(code: CompetitionCode) {
    currentCode.value = code
    competition.setCurrent(code)
    competitionDialogOpen.value = false
    applyDocumentTheme()
    showToast(`已切换到${current.value.name} · ${current.value.season}`)
  }

  function toggleTheme() {
    settings.toggleTheme()
    applyDocumentTheme()
    showToast(settings.theme === 'dark' ? '已切换深色模式' : '已切换浅色模式')
  }

  function showToast(message: string) {
    window.clearTimeout(toastTimer)
    toast.value = message
    toastTimer = window.setTimeout(() => { toast.value = '' }, 1900)
  }

  async function initialize() {
    await competition.fetchCompetitions()
    // Hero copy from the backend; failure keeps the static theme.
    try {
      featuredConfigs.value = (await apiGetFeatured()) || {}
    } catch {
      featuredConfigs.value = {}
    }
    featuredLoaded.value = true
    const code = competition.currentCode
    if (isCompetitionCode(code)) {
      currentCode.value = code
    } else if (competition.competitions.length > 0) {
      // Stored code is disabled: adopt the store's fallback (first enabled).
      // When NOTHING is enabled, leave the selection untouched — never
      // write a stale WC key back to localStorage.
      competition.setCurrent(currentCode.value)
    }
    applyDocumentTheme()
  }

  return {
    competitionCodes, currentCode, current, currentCopy, progress, competitionDialogOpen, toast,
    settings, competition, featuredConfigs, featuredLoaded, isCompetitionCode, themeFor,
    initialize, applyDocumentTheme, setCompetition, toggleTheme, showToast,
  }
})
