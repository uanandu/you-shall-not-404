<script lang="ts">
  import { onMount } from 'svelte';

  interface Props { lang?: 'en' | 'fr' }
  let { lang = 'en' }: Props = $props();

  let time = $state('');
  let tz   = $state('');

  function tick() {
    const now = new Date();
    time = now.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    });
    const parts = new Intl.DateTimeFormat('en-US', { timeZoneName: 'short' }).formatToParts(now);
    tz = parts.find(p => p.type === 'timeZoneName')?.value ?? '';
  }

  onMount(() => {
    tick();
    const id = setInterval(tick, 10_000);
    return () => clearInterval(id);
  });
</script>

<span role="status" aria-live="polite" aria-label={lang === 'fr' ? 'Heure actuelle à Montréal' : 'Current time in Montréal'}>
  Montréal <span aria-hidden="true">|</span> {time} {tz}
</span>
