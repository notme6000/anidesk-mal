import { writable, derived } from 'svelte/store'

export interface User {
  id: number
  mal_id: number
  username: string
  avatar_url: string
  created_at: string
}

interface AuthState {
  user: User | null
  loading: boolean
  error: string | null
}

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>({
    user: null,
    loading: true,
    error: null,
  })

  return {
    subscribe,

    async checkSession() {
      update(s => ({ ...s, loading: true }))
      try {
        const { GetActiveUser } = await import('../../wailsjs/go/main/App')
        const raw = await GetActiveUser()
        const data = JSON.parse(raw)
        if (data.error) {
          update(s => ({ ...s, user: null, loading: false }))
        } else {
          update(s => ({ ...s, user: data as User, loading: false }))
        }
      } catch {
        update(s => ({ ...s, user: null, loading: false }))
      }
    },

    async login() {
      update(s => ({ ...s, loading: true, error: null }))
      try {
        const { Login } = await import('../../wailsjs/go/main/App')
        const raw = await Login()
        const data = JSON.parse(raw)
        if (data.error) {
          update(s => ({ ...s, loading: false, error: data.error }))
        } else {
          update(s => ({ ...s, user: data as User, loading: false, error: null }))
        }
      } catch (e) {
        update(s => ({ ...s, loading: false, error: String(e) }))
      }
    },

    async logout() {
      try {
        const { Logout } = await import('../../wailsjs/go/main/App')
        await Logout()
      } catch { /* ignore */ }
      set({ user: null, loading: false, error: null })
    },
  }
}

export const authStore = createAuthStore()
export const isLoggedIn = derived(authStore, $auth => $auth.user !== null)
export const currentUser = derived(authStore, $auth => $auth.user)
