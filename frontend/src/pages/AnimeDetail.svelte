<script lang="ts">
  import { onMount } from 'svelte'
  import { link } from 'svelte-spa-router'
  import { isLoggedIn } from '../stores/auth'

  export let params: { id: string } = { id: '' }

  let anime: any = null
  let loading = true
  let error: string | null = null
  let libraryStatus = ''

  onMount(async () => {
    if (!$isLoggedIn) {
      loading = false
      return
    }

    try {
      const { GetAnimeDetails } = await import('../../wailsjs/go/main/App')
      const raw = await GetAnimeDetails(parseInt(params.id))
      const data = JSON.parse(raw)
      if (data.error) {
        error = data.error
      } else {
        anime = data
      }
    } catch (e) {
      error = String(e)
    } finally {
      loading = false
    }
  })

  const statusOptions = ['watching', 'completed', 'on_hold', 'dropped', 'plan_to_watch']

  async function updateStatus(status: string) {
    try {
      const { UpdateAnimeStatus } = await import('../../wailsjs/go/main/App')
      const raw = await UpdateAnimeStatus(parseInt(params.id), status, 0, 0)
      const data = JSON.parse(raw)
      if (data.ok) libraryStatus = status
    } catch { /* ignore */ }
  }
</script>

<div class="min-h-screen">
  {#if loading}
    <div class="animate-pulse">
      <div class="h-64 bg-surface-800" />
      <div class="p-6 space-y-4">
        <div class="h-8 w-64 bg-surface-800 rounded-lg" />
        <div class="h-4 w-96 bg-surface-800 rounded" />
        <div class="h-32 bg-surface-800 rounded-xl" />
      </div>
    </div>
  {:else if anime}
    <div class="relative">
      {#if anime.main_picture?.large}
        <div class="h-64 md:h-80 bg-surface-900 overflow-hidden">
          <img
            src={anime.main_picture.large}
            alt={anime.title}
            class="w-full h-full object-cover opacity-40"
          />
          <div class="absolute inset-0 bg-gradient-to-t from-[#0d1117] via-[#0d1117]/60 to-transparent" />
        </div>
      {/if}

      <div class="relative -mt-32 px-6 pb-6 max-w-5xl mx-auto">
        <div class="flex gap-6">
          {#if anime.main_picture?.medium}
            <div class="w-48 shrink-0 -mt-16">
              <img
                src={anime.main_picture.medium}
                alt={anime.title}
                class="w-full rounded-2xl shadow-2xl shadow-black/40"
              />
            </div>
          {/if}

          <div class="flex-1 pt-4">
            <h1 class="text-3xl font-bold text-surface-100 mb-1">{anime.title}</h1>
            {#if anime.alternative_titles?.en && anime.alternative_titles.en !== anime.title}
              <p class="text-surface-400 text-sm mb-3">{anime.alternative_titles.en}</p>
            {/if}

            <div class="flex items-center gap-4 flex-wrap text-sm mb-4">
              {#if anime.mean}
                <span class="flex items-center gap-1 text-yellow-400 font-medium">
                  &#9733; {anime.mean.toFixed(1)}
                </span>
              {/if}
              {#if anime.rank}
                <span class="text-surface-400">Rank #{anime.rank}</span>
              {/if}
              {#if anime.popularity}
                <span class="text-surface-400">Popularity #{anime.popularity}</span>
              {/if}
              {#if anime.media_type}
                <span class="text-surface-500 uppercase">{anime.media_type}</span>
              {/if}
              {#if anime.num_episodes}
                <span class="text-surface-500">{anime.num_episodes} episodes</span>
              {/if}
              {#if anime.status}
                <span class="text-surface-500">{anime.status.replace(/_/g, ' ')}</span>
              {/if}
            </div>

            {#if $isLoggedIn}
              <div class="flex gap-2 flex-wrap">
                {#each statusOptions as status}
                  <button
                    class="px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-200"
                    class:bg-brand-600:text-white={libraryStatus === status}
                    class:bg-surface-800:text-surface-400={libraryStatus !== status}
                    class:hover:bg-surface-700={libraryStatus !== status}
                    on:click={() => updateStatus(status)}
                  >
                    {status.replace(/_/g, ' ')}
                  </button>
                {/each}
              </div>
            {/if}

            {#if anime.genres && anime.genres.length > 0}
              <div class="flex gap-2 flex-wrap mt-4">
                {#each anime.genres as genre}
                  <span class="px-2.5 py-1 text-xs rounded-lg bg-brand-600/10 text-brand-400 font-medium">
                    {genre.name}
                  </span>
                {/each}
              </div>
            {/if}
          </div>
        </div>

        {#if anime.synopsis}
          <div class="mt-8">
            <h2 class="text-lg font-semibold text-surface-100 mb-2">Synopsis</h2>
            <p class="text-surface-300 leading-relaxed text-sm max-w-3xl">{anime.synopsis}</p>
          </div>
        {/if}

        {#if anime.studios && anime.studios.length > 0}
          <div class="mt-8">
            <h2 class="text-lg font-semibold text-surface-100 mb-2">Studios</h2>
            <div class="flex gap-2 flex-wrap">
              {#each anime.studios as studio}
                <span class="px-3 py-1.5 text-sm bg-surface-800 rounded-lg text-surface-300">{studio.name}</span>
              {/each}
            </div>
          </div>
        {/if}
      </div>
    </div>
  {:else if error}
    <div class="flex items-center justify-center min-h-[60vh]">
      <div class="text-center">
        <p class="text-surface-400 mb-2">Failed to load anime</p>
        <p class="text-surface-600 text-sm">{error}</p>
      </div>
    </div>
  {/if}
</div>
