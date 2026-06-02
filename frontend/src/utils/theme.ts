export type Theme = 'light' | 'dark'

export function getStoredTheme(): Theme {
  const stored = localStorage.getItem('wm-theme')
  if (stored === 'dark' || stored === 'light') return stored
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function setTheme(theme: Theme) {
  document.documentElement.dataset.theme = theme
  localStorage.setItem('wm-theme', theme)
}
