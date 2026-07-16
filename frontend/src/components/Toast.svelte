<script lang="ts">
  import { toastStore } from '../stores/toast'
  import Icon from '../lib/Icon.svelte'

  const icons: Record<string, string> = {
    success: `<path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>`,
    error: `<path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"/>`,
    info: `<path stroke-linecap="round" stroke-linejoin="round" d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"/>`,
  }

  const bgColors: Record<string, string> = {
    success: 'bg-green-900/80 border-green-700/50',
    error: 'bg-red-900/80 border-red-700/50',
    info: 'bg-surface-800/80 border-surface-700/50',
  }

  const iconColors: Record<string, string> = {
    success: 'text-green-400',
    error: 'text-red-400',
    info: 'text-brand-400',
  }
</script>

<div class="fixed bottom-6 right-6 z-50 space-y-2">
  {#each $toastStore as toast (toast.id)}
    <div
      class="card border px-4 py-3 flex items-center gap-3 min-w-[300px] max-w-[400px] shadow-2xl shadow-black/30 animate-slide-up {bgColors[toast.type]}"
      style="animation: slide-up 0.3s ease-out"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="w-5 h-5 shrink-0 {iconColors[toast.type]}"
      >
        {@html icons[toast.type] || ''}
      </svg>
      <p class="text-sm text-surface-100 flex-1">{toast.message}</p>
    </div>
  {/each}
</div>

<style>
  @keyframes slide-up {
    from {
      opacity: 0;
      transform: translateY(1rem);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
