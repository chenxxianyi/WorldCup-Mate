import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'
import { competitionCodes, isCompetitionCode, leagueThemes, type CompetitionCode } from '@/data/leagueTheme'
import { useCompetitionStore } from '@/stores/useCompetitionStore'
import { useSettingStore } from '@/stores/useSettingStore'

const COMPETITION_KEY = 'wm-competition'

export const useLeagueThemeStore = defineStore('league-theme', () => {
  const stored = localStorage.getItem(COMPETITION_KEY)
  const currentCode = ref<CompetitionCode>(isCompetitionCode(stored) ? stored : 'WC')
  const competitionDialogOpen = ref(false)
  const toast = ref('')
  let toastTimer = 0

  const settings = useSettingStore()
  const competition = useCompetitionStore()
  const current = computed(() => leagueThemes[currentCode.value])
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
    const code = competition.currentCode
    if (isCompetitionCode(code)) currentCode.value = code
    else competition.setCurrent(currentCode.value)
    applyDocumentTheme()
  }

  return {
    competitionCodes, currentCode, current, progress, competitionDialogOpen, toast,
    settings, competition, initialize, applyDocumentTheme, setCompetition, toggleTheme, showToast,
  }
})
