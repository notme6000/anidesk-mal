import { writable } from 'svelte/store'

export interface SearchResult {
  id: number
  title: string
  main_picture?: { medium: string; large: string }
  mean?: number
  rank?: number
  media_type?: string
  num_episodes?: number
  synopsis?: string
  genres?: { id: number; name: string }[]
  status?: string
}

interface SearchState {
  query: string
  results: SearchResult[]
  loading: boolean
  type: 'anime' | 'manga' | 'character' | 'studio'
}

function createSearchStore() {
  const { subscribe, set, update } = writable<SearchState>({
    query: '',
    results: [],
    loading: false,
    type: 'anime',
  })

  return {
    subscribe,

    setType(type: SearchState['type']) {
      update(s => ({ ...s, type }))
    },

    setQuery(query: string) {
      update(s => ({ ...s, query }))
    },

    async search(query: string) {
      if (!query.trim()) {
        update(s => ({ ...s, results: [], loading: false }))
        return
      }

      update(s => ({ ...s, loading: true, query }))

      try {
        const { SearchAnime } = await import('../../wailsjs/go/main/App')
        const raw = await SearchAnime(query, 30)
        const results = JSON.parse(raw)
        update(s => ({ ...s, loading: false, results: Array.isArray(results) ? results : [] }))
      } catch {
        update(s => ({ ...s, loading: false, results: [] }))
      }
    },

    clear() {
      set({ query: '', results: [], loading: false, type: 'anime' })
    },
  }
}

export const searchStore = createSearchStore()
