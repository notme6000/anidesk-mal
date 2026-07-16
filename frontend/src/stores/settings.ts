import { writable } from 'svelte/store'

interface SettingsState {
  cacheEnabled: boolean
  notificationsEnabled: boolean
  syncInterval: number
  [key: string]: unknown
}

function createSettingsStore() {
  const { subscribe, set, update } = writable<SettingsState>({
    cacheEnabled: true,
    notificationsEnabled: true,
    syncInterval: 5,
  })

  return {
    subscribe,

    async load() {
      try {
        const { GetSettings } = await import('../../wailsjs/go/main/App')
        const raw = await GetSettings()
        const data = JSON.parse(raw)
        set({ ...data } as SettingsState)
      } catch { /* use defaults */ }
    },

    async set(key: string, value: string) {
      try {
        const { SetSetting } = await import('../../wailsjs/go/main/App')
        await SetSetting(key, value)
        update(s => ({ ...s, [key]: value }))
      } catch { /* ignore */ }
    },
  }
}

export const settingsStore = createSettingsStore()
