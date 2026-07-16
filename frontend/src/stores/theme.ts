import { writable, derived } from 'svelte/store'

type Theme = 'dark' | 'light'
type AccentColor = 'blue' | 'purple' | 'green' | 'orange' | 'pink'

interface ThemeState {
  theme: Theme
  accentColor: AccentColor
}

function createThemeStore() {
  const stored = typeof localStorage !== 'undefined'
    ? localStorage.getItem('anidesk-theme')
    : null
  const initial: ThemeState = stored ? JSON.parse(stored) : { theme: 'dark', accentColor: 'blue' }

  const { subscribe, set, update } = writable<ThemeState>(initial)

  return {
    subscribe,

    toggleTheme() {
      update(s => {
        const next = { ...s, theme: s.theme === 'dark' ? 'light' : 'dark' as Theme }
        localStorage.setItem('anidesk-theme', JSON.stringify(next))
        return next
      })
    },

    setAccentColor(color: AccentColor) {
      update(s => {
        const next = { ...s, accentColor: color }
        localStorage.setItem('anidesk-theme', JSON.stringify(next))
        return next
      })
    },
  }
}

export const themeStore = createThemeStore()
export const isDark = derived(themeStore, $t => $t.theme === 'dark')
