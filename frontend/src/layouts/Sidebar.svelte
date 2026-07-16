<script lang="ts">
  import { link, location } from 'svelte-spa-router'
  import { onMount } from 'svelte'
  import Icon from '../lib/Icon.svelte'
  import { authStore, isLoggedIn, currentUser } from '../stores/auth'

  const navItems = [
    { href: '/', label: 'Home', icon: 'home' },
    { href: '/search', label: 'Search', icon: 'search' },
    { href: '/library', label: 'Library', icon: 'library' },
    { href: '/settings', label: 'Settings', icon: 'settings' },
  ] as const

  $: activeRoute = $location

  function navClass(href: string): string {
    const active = activeRoute === href
    const base = 'flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200'
    if (active) return `${base} text-surface-100 bg-surface-800/50`
    return `${base} text-surface-400 hover:text-surface-100 hover:bg-surface-800/50`
  }

  async function handleLogin() {
    await authStore.login()
  }

  async function handleLogout() {
    await authStore.logout()
  }
</script>

<aside class="flex flex-col w-60 h-full bg-surface-950/90 backdrop-blur-md border-r border-surface-800/30 select-none">
  <div class="flex items-center gap-3 px-5 h-16 shrink-0 border-b border-surface-800/30">
    <div class="w-8 h-8 rounded-lg bg-brand-600 flex items-center justify-center">
      <span class="text-white font-bold text-sm">A</span>
    </div>
    <span class="font-semibold text-surface-100 text-base tracking-tight">AniDesk</span>
  </div>

  <nav class="flex-1 py-3 px-3 space-y-1">
    {#each navItems as item}
      <a
        href={item.href}
        use:link
        class={navClass(item.href)}
      >
        <Icon name={item.icon} className="w-5 h-5 shrink-0" />
        <span>{item.label}</span>
      </a>
    {/each}
  </nav>

  <div class="p-3 border-t border-surface-800/30">
    {#if $isLoggedIn && $currentUser}
      <div class="flex items-center gap-3 px-3 py-2">
        {#if $currentUser.avatar_url}
          <img
            src={$currentUser.avatar_url}
            alt={$currentUser.username}
            class="w-8 h-8 rounded-full"
          />
        {:else}
          <div class="w-8 h-8 rounded-full bg-surface-700 flex items-center justify-center">
            <span class="text-xs text-surface-300">{(($currentUser.username) || '?')[0]}</span>
          </div>
        {/if}
        <div class="flex-1 min-w-0">
          <p class="text-sm text-surface-100 truncate">{$currentUser.username}</p>
          <button class="text-xs text-surface-500 hover:text-surface-300 transition-colors" on:click={handleLogout}>
            Logout
          </button>
        </div>
      </div>
    {:else if $authStore.loading}
      <div class="flex items-center gap-3 px-3 py-2">
        <div class="w-8 h-8 rounded-full bg-surface-800 animate-pulse" />
        <div class="flex-1 space-y-1.5">
          <div class="h-3 w-20 bg-surface-800 rounded animate-pulse" />
          <div class="h-2 w-14 bg-surface-800 rounded animate-pulse" />
        </div>
      </div>
    {:else}
      <button
        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-surface-400 hover:text-surface-100 hover:bg-surface-800/50 transition-all duration-200"
        on:click={handleLogin}
      >
        <Icon name="home" className="w-5 h-5 shrink-0" />
        <span>Sign In</span>
      </button>
    {/if}
  </div>
</aside>
