import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { apiListCompetitions } from '@/api/competitions'
import { WC_CODE, type Competition } from '@/types/competition'

const STORAGE_KEY = 'wm-competition'
const STORAGE_KEY_ID = 'wm-competition-id'

export const useCompetitionStore = defineStore('competition', () => {
  const competitions = ref<Competition[]>([])
  const currentCode = ref<string>(localStorage.getItem(STORAGE_KEY) || WC_CODE)
  const loaded = ref(false)

  const current = computed<Competition | null>(
    () => competitions.value.find((c) => c.code === currentCode.value) || null,
  )
  const isWorldCup = computed(() => currentCode.value === WC_CODE)
  const isLeague = computed(() => current.value?.format === 'league')

  async function fetchCompetitions() {
    if (loaded.value) return
    try {
      const res = await apiListCompetitions()
      competitions.value = (res || []).filter((c) => c.status === 'active')
      loaded.value = true
    } catch {
      loaded.value = false // allow retry on next mount instead of silent degradation
    }
    syncStoredId()
    // Clean up stale selections: if the stored code no longer resolves to
    // an enabled competition, fall back to the first enabled one. When
    // NOTHING is enabled, keep the current code as-is and never write a
    // stale WC default — a disabled World Cup must stay hidden.
    if (
      loaded.value &&
      !competitions.value.some((c) => c.code === currentCode.value)
    ) {
      if (competitions.value.length > 0) {
        const fallback = competitions.value[0].code
        currentCode.value = fallback
        localStorage.setItem(STORAGE_KEY, fallback)
      }
      localStorage.removeItem(STORAGE_KEY_ID)
    }
  }

  function syncStoredId() {
    const current = competitions.value.find((c) => c.code === currentCode.value)
    if (current) localStorage.setItem(STORAGE_KEY_ID, String(current.id))
  }

  function setCurrent(code: string) {
    currentCode.value = code
    localStorage.setItem(STORAGE_KEY, code)
    syncStoredId()
  }

  return {
    competitions,
    currentCode,
    current,
    loaded,
    isWorldCup,
    isLeague,
    fetchCompetitions,
    setCurrent,
  }
})
