<script lang="ts">
  import { authStore, isLoggedIn, currentUser } from '../stores/auth'

  async function handleLogin() {
    await authStore.login()
  }

  async function handleLogout() {
    await authStore.logout()
  }
</script>

{#if $isLoggedIn && $currentUser}
  <div class="flex items-center gap-3">
    {#if $currentUser.avatar_url}
      <img
        src={$currentUser.avatar_url}
        alt={$currentUser.username}
        class="w-8 h-8 rounded-full"
      />
    {/if}
    <span class="text-sm text-surface-100">{$currentUser.username}</span>
    <button class="btn-ghost text-xs" on:click={handleLogout}>Logout</button>
  </div>
{:else if $authStore.loading}
  <div class="w-8 h-8 rounded-full bg-surface-800 animate-pulse" />
{:else}
  <button class="btn-primary" on:click={handleLogin}>
    Login with MyAnimeList
  </button>
{/if}
