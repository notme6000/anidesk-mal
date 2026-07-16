<script lang="ts">
  import AnimeCard from './AnimeCard.svelte'

  export let title: string
  export let items: Array<{
    id: number
    title: string
    main_picture?: { medium: string; large: string }
    mean?: number
    rank?: number
    media_type?: string
    num_episodes?: number
  }>
</script>

{#if items && items.length > 0}
  <section>
    <h2 class="text-lg font-semibold text-surface-100 mb-4">{title}</h2>
    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
      {#each items as item}
        <AnimeCard anime={item} />
      {/each}
    </div>
  </section>
{/if}
