<script lang="ts">
  import { onMount } from 'svelte'
  import { searchStore } from '../stores/search'
  import AnimeCard from '../components/AnimeCard.svelte'

  let query = ''

  onMount(() => {
    query = $searchStore.query
  })

  async function handleSearch() {
    await searchStore.search(query)
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter') handleSearch()
  }
</script>

<div class="p-6 max-w-7xl mx-auto">
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-surface-100 mb-4">Search</h1>
    <div class="flex gap-3">
      <input
        type="text"
        bind:value={query}
        on:keydown={handleKeydown}
        placeholder="Search anime, manga, characters..."
        class="input flex-1"
      />
      <button class="btn-primary" on:click={handleSearch}>
        Search
      </button>
    </div>
  </div>

  {#if $searchStore.loading}
    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
      {#each Array(12) as _}
        <div class="aspect-[3/4] bg-surface-800 rounded-2xl animate-pulse" />
      {/each}
    </div>
  {:else if $searchStore.results.length > 0}
    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
      {#each $searchStore.results as item}
        <AnimeCard anime={item} />
      {/each}
    </div>
  {:else if query && !$searchStore.loading}
    <div class="text-center py-20 text-surface-500">
      No results found for "{query}"
    </div>
  {:else}
    <div class="text-center py-20 text-surface-500">
      Start typing to search
    </div>
  {/if}
</div>
