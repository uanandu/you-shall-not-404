<script lang="ts">
  import { fade } from 'svelte/transition';
  import { onMount } from 'svelte';

  interface Props {
    links: { label: string; href: string }[];
    langToggle?: { label: string; href: string };
    menuLabel?: string;
    closeLabel?: string;
  }

  const { links, langToggle, menuLabel = 'Menu', closeLabel = 'Close' }: Props = $props();

  let open = $state(false);

  onMount(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') open = false; };
    window.addEventListener('keydown', onKey);
    document.addEventListener('astro:before-preparation', () => { open = false; });
    return () => window.removeEventListener('keydown', onKey);
  });
</script>

<button
  class="burger label"
  onclick={() => (open = !open)}
  aria-label={open ? closeLabel : menuLabel}
  aria-expanded={open}
>
  {open ? closeLabel : menuLabel}
</button>

{#if open}
  <div
    class="drawer"
    role="dialog"
    aria-modal="true"
    transition:fade={{ duration: 220 }}
  >
    {#each links as link, i}
      <a
        href={link.href}
        class="drawer__link"
        style="animation-delay: {i * 55}ms"
        onclick={() => (open = false)}
      >
        {link.label}<span class="arrow" aria-hidden="true">↗</span>
      </a>
    {/each}
    {#if langToggle}
      <a
        href={langToggle.href}
        class="drawer__link drawer__link--lang"
        style="animation-delay: {links.length * 55}ms"
        onclick={() => (open = false)}
      >
        {langToggle.label}
      </a>
    {/if}
  </div>
{/if}

<style>
  .burger {
    display: none;
    position: fixed;
    top: 1.5rem;
    right: 1.25rem;
    z-index: 102;
    background: none;
    border: none;
    cursor: pointer;
    color: var(--clr-ink);
    padding: 0;
    font-family: var(--font-body);
    pointer-events: all;
  }

  /* ── Drawer ─────────────────────────────────────────────── */
  .drawer {
    position: fixed;
    inset: 0;
    z-index: 101;
    background: var(--clr-bg);
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: center;
    gap: 0;
    padding: 0 2rem;
    pointer-events: all;
  }

  /* ── Links ──────────────────────────────────────────────── */
  .drawer__link {
    display: flex;
    align-items: baseline;
    gap: 0.4em;
    font-family: var(--font-display);
    font-size: clamp(2.8rem, 10vw, 4.5rem);
    font-weight: 700;
    color: var(--clr-ink);
    line-height: 1.15;
    border-bottom: 1px solid var(--clr-rule);
    width: 100%;
    padding-block: 0.5rem;
    pointer-events: all;

    opacity: 0;
    animation: link-enter 0.35s cubic-bezier(0.22, 1, 0.36, 1) forwards;
  }

  .drawer__link--lang {
    font-size: clamp(1.1rem, 3vw, 1.4rem);
    font-weight: 400;
    color: var(--clr-muted);
    border-bottom: none;
    padding-top: 1.25rem;
    border-top: 1px solid var(--clr-rule);
  }

  .drawer__link--lang:active {
    color: var(--clr-ink);
    opacity: 1;
  }

  @media (hover: hover) {
    .drawer__link--lang:hover {
      color: var(--clr-ink);
      opacity: 1;
    }
  }

  .drawer__link:last-child { border-bottom: none; }

  @keyframes link-enter {
    from { opacity: 0; transform: translateY(10px); }
    to   { opacity: 1; transform: translateY(0); }
  }

  /* ── Arrow ──────────────────────────────────────────────── */
  .arrow {
    font-family: var(--font-body);
    font-size: 0.38em;
    font-weight: 400;
    color: var(--clr-muted);
    align-self: flex-start;
    margin-top: 0.6em;
    letter-spacing: 0;
    transition: color var(--transition), transform var(--transition);
  }

  .drawer__link {
    transition: opacity var(--transition);
  }

  .drawer__link:active {
    opacity: 0.55;
  }

  .drawer__link:active .arrow {
    color: var(--clr-ink);
    transform: translate(2px, -2px);
  }

  @media (hover: hover) {
    .drawer__link:hover { opacity: 0.55; }
    .drawer__link:hover .arrow {
      color: var(--clr-ink);
      transform: translate(2px, -2px);
    }
  }

  @media (max-width: 640px) {
    .burger { display: block; }
  }
</style>
