<script lang="ts">
  interface Props {
    lang?: 'en' | 'fr';
  }

  const { lang = 'en' }: Props = $props();

  const s = {
    en: { play: 'Play lofi music', pause: 'Pause lofi music' },
    fr: { play: 'Jouer la musique lofi', pause: 'Mettre la musique en pause' },
  }[lang];

  let audioEl: HTMLAudioElement;
  let playing = $state(false);

  function toggle() {
    if (playing) audioEl.pause();
    else audioEl.play().catch(() => {});
  }
</script>

<button
  type="button"
  class="music"
  class:music--playing={playing}
  onclick={toggle}
  aria-pressed={playing}
  aria-label={playing ? s.pause : s.play}
>
  <span class="music__box">
    <svg class="music__icon" viewBox="0 0 24 24" aria-hidden="true">
      <rect x="2" y="5" width="20" height="14" fill="currentColor" />
      <circle cx="5" cy="6.6" r="0.4" fill="var(--clr-bg)" />
      <circle cx="19" cy="6.6" r="0.4" fill="var(--clr-bg)" />
      <circle cx="8" cy="12" r="3.4" fill="var(--clr-bg)" />
      <circle cx="16" cy="12" r="3.4" fill="var(--clr-bg)" />
      <rect x="9" y="15.6" width="6" height="1.8" fill="var(--clr-rule)" />

      <g class="music__reel music__reel--a">
        <circle cx="8" cy="12" r="2.6" fill="none" stroke="var(--clr-rule)" stroke-width="0.5" />
        <line x1="8" y1="9.8" x2="8" y2="14.2" stroke="var(--clr-rule)" stroke-width="0.5" />
        <line x1="6.1" y1="10.9" x2="9.9" y2="13.1" stroke="var(--clr-rule)" stroke-width="0.5" />
        <line x1="6.1" y1="13.1" x2="9.9" y2="10.9" stroke="var(--clr-rule)" stroke-width="0.5" />
        <circle cx="8" cy="12" r="0.6" fill="var(--clr-rule)" />
      </g>

      <g class="music__reel music__reel--b">
        <circle cx="16" cy="12" r="2.6" fill="none" stroke="var(--clr-rule)" stroke-width="0.5" />
        <line x1="16" y1="9.8" x2="16" y2="14.2" stroke="var(--clr-rule)" stroke-width="0.5" />
        <line x1="14.1" y1="10.9" x2="17.9" y2="13.1" stroke="var(--clr-rule)" stroke-width="0.5" />
        <line x1="14.1" y1="13.1" x2="17.9" y2="10.9" stroke="var(--clr-rule)" stroke-width="0.5" />
        <circle cx="16" cy="12" r="0.6" fill="var(--clr-rule)" />
      </g>
    </svg>
  </span>
</button>

<audio
  bind:this={audioEl}
  src="/audio/lofi.mp3"
  loop
  preload="none"
  onplay={() => (playing = true)}
  onpause={() => (playing = false)}
></audio>

<style>
  .music {
    appearance: none;
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    display: flex;
    align-items: center;
    line-height: 0;
  }

  .music__box {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 1.6rem;
    height: 1.6rem;
    border: 1px solid var(--clr-rule);
    color: var(--clr-ink-2);
    transform: rotate(-4deg);
    transition: border-color var(--transition), color var(--transition), transform var(--transition);
  }

  .music:hover .music__box,
  .music:focus-visible .music__box {
    border-color: var(--clr-ink);
    color: var(--clr-ink);
    transform: rotate(0deg);
  }

  .music__icon {
    width: 1.15rem;
    height: 1.15rem;
  }

  .music__reel {
    transform-box: fill-box;
    transform-origin: center;
    animation: music-spin 3.6s linear infinite;
    animation-play-state: paused;
  }

  /* reels turn at slightly different rates, like a real deck's supply
     and takeup spools rather than two mirrored copies of one motion */
  .music__reel--b { animation-duration: 3s; }

  .music--playing .music__reel { animation-play-state: running; }

  @keyframes music-spin {
    to { transform: rotate(360deg); }
  }

  @media (prefers-reduced-motion: reduce) {
    .music__reel { animation: none; }
    .music__box { transition: none; }
  }
</style>
