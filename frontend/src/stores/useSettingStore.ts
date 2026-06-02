import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getStoredTheme, setTheme as applyTheme, type Theme } from '@/utils/theme'

export const useSettingStore = defineStore('setting', () => {
  const theme = ref<Theme>(getStoredTheme())
  const timezone = ref<string>(
    localStorage.getItem('wm-timezone') || Intl.DateTimeFormat().resolvedOptions().timeZone
  )
  const language = ref<string>(localStorage.getItem('wm-language') || 'zh-CN')

  function toggleTheme() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    applyTheme(theme.value)
  }

  function setTheme(t: Theme) {
    theme.value = t
    applyTheme(t)
  }

  function setTimezone(tz: string) {
    timezone.value = tz
    localStorage.setItem('wm-timezone', tz)
  }

  function setLanguage(lang: string) {
    language.value = lang
    localStorage.setItem('wm-language', lang)
  }

  return { theme, timezone, language, toggleTheme, setTheme, setTimezone, setLanguage }
})
