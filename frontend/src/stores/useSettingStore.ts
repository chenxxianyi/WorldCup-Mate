import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getStoredTheme, setTheme as applyTheme, type Theme } from '@/utils/theme'
import { DEFAULT_TIMEZONE } from '@/utils/datetime'

export const useSettingStore = defineStore('setting', () => {
  const theme = ref<Theme>(getStoredTheme())
  const timezone = ref<string>(
    localStorage.getItem('wm-timezone') || DEFAULT_TIMEZONE
  )
  const language = ref<string>(localStorage.getItem('wm-language') || 'zh-CN')
  const defaultReminderChannel = ref<string>(localStorage.getItem('wm-reminder-channel') || 'site')

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

  function setDefaultReminderChannel(ch: string) {
    defaultReminderChannel.value = ch
    localStorage.setItem('wm-reminder-channel', ch)
  }

  return { theme, timezone, language, defaultReminderChannel, toggleTheme, setTheme, setTimezone, setLanguage, setDefaultReminderChannel }
})
