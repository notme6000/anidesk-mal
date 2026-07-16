<script lang="ts">
  import { onMount } from 'svelte'
  import Sidebar from './Sidebar.svelte'
  import Router from 'svelte-spa-router'
  import Toast from '../components/Toast.svelte'
  import Dashboard from '../pages/Dashboard.svelte'
  import AnimeDetail from '../pages/AnimeDetail.svelte'
  import Search from '../pages/Search.svelte'
  import Library from '../pages/Library.svelte'
  import Settings from '../pages/Settings.svelte'
  import { authStore } from '../stores/auth'

  const routes = {
    '/': Dashboard,
    '/anime/:id': AnimeDetail,
    '/search': Search,
    '/library': Library,
    '/library/:category': Library,
    '/settings': Settings,
  }

  onMount(() => {
    authStore.checkSession()
  })
</script>

<div class="flex h-screen w-screen bg-[#0d1117] text-surface-100">
  <Sidebar />
  <main class="flex-1 overflow-y-auto">
    <Router {routes} />
  </main>
  <Toast />
</div>
