<script lang="ts">
  import { onMount } from 'svelte'
  import { authStore, isLoggedIn } from '../stores/auth'
  import LoginButton from '../components/LoginButton.svelte'
  import AnimeRow from '../components/AnimeRow.svelte'

  let seasonal: any[] = []
  let topAiring: any[] = []
  let trending: any[] = []
  let loading = true
  let error: string | null = null

  onMount(async () => {
    await authStore.checkSession()
    if (!$isLoggedIn) {
      loading = false
      return
    }

    try {
      const { GetSeasonalAnime, GetAnimeRanking } = await import('../../wailsjs/go/main/App')

      const now = new Date()
      const year = now.getFullYear()
      const month = now.getMonth()
      const season = month < 3 ? 'winter' : month < 6 ? 'spring' : month < 9 ? 'summer' : 'fall'

      const [seasonalRaw, airingRaw, trendingRaw] = await Promise.all([
        GetSeasonalAnime(year, season),
        GetAnimeRanking('airing'),
        GetAnimeRanking('bypopularity'),
      ])

      seasonal = JSON.parse(seasonalRaw) || []
      topAiring = JSON.parse(airingRaw) || []
      trending = JSON.parse(trendingRaw) || []
    } catch (e) {
      error = String(e)
    } finally {
      loading = false
    }
  })
</script>

<div class="p-6 max-w-7xl mx-auto">
  {#if !$isLoggedIn && !loading}
    <div class="flex flex-col items-center justify-center min-h-[60vh] gap-6">
      <div class="text-center max-w-md">
        <div class="w-16 h-16 rounded-2xl bg-brand-600/20 flex items-center justify-center mx-auto mb-6">
          <span class="text-3xl font-bold text-brand-400">A</span>
        </div>
        <h1 class="text-3xl font-bold text-surface-100 mb-3">Welcome to AniDesk</h1>
        <p class="text-surface-400 leading-relaxed mb-2">
          Your lightweight desktop companion for MyAnimeList.
        </p>
        <p class="text-surface-500 text-sm mb-8">
          Track your anime, discover new shows, and sync your library — all in a native desktop app.
        </p>
        <LoginButton />
      </div>
    </div>
  {:else if loading}
    <div class="space-y-8 animate-pulse">
      <div class="h-8 w-48 bg-surface-800 rounded-lg" />
      <div class="grid grid-cols-6 gap-4">
        {#each Array(6) as _}
          <div class="aspect-[3/4] bg-surface-800 rounded-2xl" />
        {/each}
      </div>
    </div>
  {:else}
    <div class="space-y-10 pb-16">
      <AnimeRow title="Seasonal Anime" {seasonal} />
      <AnimeRow title="Top Airing" items={topAiring} />
      <AnimeRow title="Trending" items={trending} />
    </div>
  {/if}

  {#if error}
    <div class="fixed bottom-6 right-6 card bg-red-900/50 border-red-700/50 text-red-200 px-4 py-3 text-sm">
      {error}
    </div>
  {/if}
</div>
