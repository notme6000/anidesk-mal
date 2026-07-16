import { writable } from 'svelte/store'

export interface Toast {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
}

function createToastStore() {
  const { subscribe, update } = writable<Toast[]>([])
  let nextId = 0

  return {
    subscribe,

    show(message: string, type: Toast['type'] = 'info') {
      const id = nextId++
      update(t => [...t, { id, message, type }])
      setTimeout(() => {
        update(t => t.filter(toast => toast.id !== id))
      }, 4000)
    },

    success(message: string) {
      this.show(message, 'success')
    },

    error(message: string) {
      this.show(message, 'error')
    },

    info(message: string) {
      this.show(message, 'info')
    },
  }
}

export const toastStore = createToastStore()
