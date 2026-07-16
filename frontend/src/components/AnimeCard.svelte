<script lang="ts">
  import { link } from 'svelte-spa-router'

  export let anime: {
    id: number
    title: string
    main_picture?: { medium: string; large: string }
    mean?: number
    rank?: number
    media_type?: string
    num_episodes?: number
    genres?: { id: number; name: string }[]
  }
</script>

<a href={`/anime/${anime.id}`} use:link class="card-hover group block overflow-hidden">
  <div class="aspect-[3/4] bg-surface-800/30 overflow-hidden rounded-t-2xl">
    {#if anime.main_picture?.medium}
      <img
        src={anime.main_picture.medium}
        alt={anime.title}
        class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
        loading="lazy"
      />
    {:else}
      <div class="w-full h-full flex items-center justify-center text-surface-600 text-sm">
        No Image
      </div>
    {/if}
  </div>
  <div class="p-3">
    <h3 class="text-sm font-medium text-surface-100 line-clamp-2 leading-snug">
      {anime.title}
    </h3>
    <div class="flex items-center gap-2 mt-1.5">
      {#if anime.mean}
        <span class="text-xs text-yellow-400 font-medium">{anime.mean.toFixed(1)}</span>
      {/if}
      {#if anime.media_type}
        <span class="text-xs text-surface-500 uppercase">{anime.media_type}</span>
      {/if}
      {#if anime.num_episodes}
        <span class="text-xs text-surface-500">{anime.num_episodes} eps</span>
      {/if}
    </div>
  </div>
</a>
