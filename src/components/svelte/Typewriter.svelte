<script lang="ts">
  import { onMount } from 'svelte';

  let { line1, line2, charDelay = 38 }: { line1: string; line2: string; charDelay?: number } = $props();

  let displayed1 = $state('');
  let displayed2 = $state('');
  let showLine2  = $state(false);

  onMount(() => {
    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    if (reduced) {
      displayed1 = line1;
      displayed2 = line2;
      showLine2  = true;
      return;
    }

    let i = 0;
    const typeLine1 = () => {
      if (i < line1.length) {
        displayed1 += line1[i++];
        setTimeout(typeLine1, charDelay);
      } else {
        setTimeout(startLine2, 180);
      }
    };

    const startLine2 = () => {
      showLine2 = true;
      let j = 0;
      const typeLine2 = () => {
        if (j < line2.length) {
          displayed2 += line2[j++];
          setTimeout(typeLine2, charDelay);
        }
      };
      typeLine2();
    };

    setTimeout(typeLine1, 700);
  });
</script>

<p class="statement">
  {displayed1}<br>
  {#if showLine2}<em>{displayed2}</em>{/if}
</p>

<style>
  .statement {
    font-family: var(--font-display);
    font-size: clamp(0.95rem, 1.4vw, 1.15rem);
    font-weight: 700;
    line-height: 1.4;
    color: var(--clr-ink);
    min-height: 2.8em;
  }

  .statement em {
    font-weight: 400;
    font-style: italic;
  }
</style>
