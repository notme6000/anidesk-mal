<script lang="ts">
  import { onMount } from 'svelte'
  import { link, location } from 'svelte-spa-router'
  import { isLoggedIn } from '../stores/auth'
  import LoginButton from '../components/LoginButton.svelte'

  const categories = [
    { key: '', label: 'All' },
    { key: 'watching', label: 'Watching' },
    { key: 'completed', label: 'Completed' },
    { key: 'on_hold', label: 'On Hold' },
    { key: 'dropped', label: 'Dropped' },
    { key: 'plan_to_watch', label: 'Plan to Watch' },
  ]

  let entries: any[] = []
  let loading = true
  let activeCategory = ''

  $: activeCategory = $location.split('/').pop() || ''

  onMount(async () => {
    if (!$isLoggedIn) {
      loading = false
      return
    }

    await loadLibrary()
  })

  async function loadLibrary() {
    loading = true
    try {
      const { GetLibrary } = await import('../../wailsjs/go/main/App')
      const raw = await GetLibrary()
      const data = JSON.parse(raw)
      entries = data.error ? [] : data
    } catch {
      entries = []
    } finally {
      loading = false
    }
  }

  $: filtered = activeCategory
    ? entries.filter(e => e.status === activeCategory)
    : entries

  async function removeEntry(animeId: number) {
    try {
      const { RemoveAnimeFromLibrary } = await import('../../wailsjs/go/main/App')
      const raw = await RemoveAnimeFromLibrary(animeId)
      const data = JSON.parse(raw)
      if (data.ok) {
        entries = entries.filter(e => e.anime_id !== animeId)
      }
    } catch { /* ignore */ }
  }
</script>

<div class="p-6 max-w-7xl mx-auto">
  {#if !$isLoggedIn && !loading}
    <div class="flex flex-col items-center justify-center min-h-[60vh] gap-4">
      <p class="text-surface-400">Login to view your library</p>
      <LoginButton />
    </div>
  {:else}
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-surface-100">My Library</h1>
      <span class="text-sm text-surface-400">{entries.length} entries</span>
    </div>

    <div class="flex gap-2 flex-wrap mb-6">
      {#each categories as cat}
        <a
          href={cat.key ? `/library/${cat.key}` : '/library'}
          use:link
          class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-200"
          class:bg-brand-600:text-white={activeCategory === cat.key}
          class:bg-surface-800:text-surface-400={activeCategory !== cat.key}
        >
          {cat.label}
        </a>
      {/each}
    </div>

    {#if loading}
      <div class="space-y-3">
        {#each Array(8) as _}
          <div class="h-20 bg-surface-800 rounded-xl animate-pulse" />
        {/each}
      </div>
    {:else if filtered.length === 0}
      <div class="text-center py-20 text-surface-500">
        {activeCategory ? 'No entries in this category' : 'Your library is empty'}
      </div>
    {:else}
      <div class="space-y-3">
        {#each filtered as entry}
          <div class="card flex items-center gap-4 p-4">
            <a href={`/anime/${entry.anime_id}`} use:link class="flex-1 min-w-0">
              <h3 class="text-sm font-medium text-surface-100 truncate">
                Anime #{entry.anime_id}
              </h3>
              <div class="flex items-center gap-3 mt-1 text-xs text-surface-400">
                <span class="capitalize">{entry.status.replace(/_/g, ' ')}</span>
                {#if entry.score > 0}
                  <span>Score: {entry.score}</span>
                {/if}
                {#if entry.episodes_watched > 0}
                  <span>Episodes: {entry.episodes_watched}</span>
                {/if}
              </div>
            </a>
            {#if entry.notes}
              <p class="text-xs text-surface-500 max-w-xs truncate hidden md:block">{entry.notes}</p>
            {/if}
            <button
              class="text-surface-500 hover:text-red-400 transition-colors text-xs px-2 py-1"
              on:click={() => removeEntry(entry.anime_id)}
            >
              Remove
            </button>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>
