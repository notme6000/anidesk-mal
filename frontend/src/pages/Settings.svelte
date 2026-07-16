<script lang="ts">
  import { onMount } from 'svelte'
  import { themeStore } from '../stores/theme'
  import { settingsStore } from '../stores/settings'

  let clientID = ''
  let saved = false

  onMount(async () => {
    await settingsStore.load()
  })

  async function saveClientID() {
    try {
      const { SetMALClientID } = await import('../../wailsjs/go/main/App')
      await SetMALClientID(clientID)
      saved = true
      setTimeout(() => saved = false, 2000)
    } catch { /* ignore */ }
  }
</script>

<div class="p-6 max-w-2xl mx-auto">
  <h1 class="text-2xl font-bold text-surface-100 mb-8">Settings</h1>

  <div class="space-y-8">
    <section>
      <h2 class="text-lg font-semibold text-surface-100 mb-4">Appearance</h2>
      <div class="card p-4 space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-surface-100 font-medium">Dark Mode</p>
            <p class="text-xs text-surface-400">Toggle between dark and light theme</p>
          </div>
          <button
            class="w-12 h-6 rounded-full transition-colors duration-200 relative"
            class:bg-brand-600={$themeStore.theme === 'dark'}
            class:bg-surface-700={$themeStore.theme === 'light'}
            on:click={themeStore.toggleTheme}
          >
            <div
              class="absolute top-0.5 w-5 h-5 bg-white rounded-full transition-transform duration-200 shadow-sm"
              class:translate-x={$themeStore.theme === 'dark' ? '6' : '0.5'}
            />
          </button>
        </div>

        <div>
          <p class="text-sm text-surface-100 font-medium mb-2">Accent Color</p>
          <div class="flex gap-2">
            {#each [
              { name: 'blue', hex: '#3b82f6' },
              { name: 'purple', hex: '#8b5cf6' },
              { name: 'green', hex: '#10b981' },
              { name: 'orange', hex: '#f97316' },
              { name: 'pink', hex: '#ec4899' },
            ] as color}
              <button
                class="w-8 h-8 rounded-full transition-all duration-200"
                class:ring-2={$themeStore.accentColor === color.name}
                class:ring-surface-100={$themeStore.accentColor === color.name}
                class:scale-110={$themeStore.accentColor === color.name}
                style:background-color={color.hex}
                on:click={() => themeStore.setAccentColor(color.name)}
              />
            {/each}
          </div>
        </div>
      </div>
    </section>

    <section>
      <h2 class="text-lg font-semibold text-surface-100 mb-4">MyAnimeList API</h2>
      <div class="card p-4">
        <p class="text-xs text-surface-400 mb-3">
          Enter your MAL Client ID to enable API access.
          Get one at <a href="https://myanimelist.net/apiconfig" target="_blank" class="text-brand-400 hover:text-brand-300">myanimelist.net/apiconfig</a>
        </p>
        <div class="flex gap-3">
          <input
            type="text"
            bind:value={clientID}
            placeholder="MAL Client ID"
            class="input flex-1"
          />
          <button class="btn-primary" on:click={saveClientID}>
            Save
          </button>
        </div>
        {#if saved}
          <p class="text-xs text-green-400 mt-2">Saved successfully</p>
        {/if}
      </div>
    </section>

    <section>
      <h2 class="text-lg font-semibold text-surface-100 mb-4">Cache</h2>
      <div class="card p-4 space-y-3">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-surface-100 font-medium">Image Cache</p>
            <p class="text-xs text-surface-400">Cache anime images locally for offline use</p>
          </div>
          <div class="w-12 h-6 rounded-full bg-brand-600 relative">
            <div class="absolute top-0.5 right-0.5 w-5 h-5 bg-white rounded-full shadow-sm" />
          </div>
        </div>
        <p class="text-xs text-surface-500">Cached images are stored in ~/.anidesk/cache/images/</p>
      </div>
    </section>

    <section>
      <h2 class="text-lg font-semibold text-surface-100 mb-4">About</h2>
      <div class="card p-4">
        <p class="text-sm text-surface-100 font-medium">AniDesk</p>
        <p class="text-xs text-surface-400 mt-1">Version 0.1.0</p>
        <p class="text-xs text-surface-500 mt-2">
          A lightweight desktop client for MyAnimeList built with Go, Wails, and Svelte.
        </p>
      </div>
    </section>
  </div>
</div>
