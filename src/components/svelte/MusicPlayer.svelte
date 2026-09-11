<script lang="ts">
  interface Props {
    lang?: 'en' | 'fr';
  }

  const { lang = 'en' }: Props = $props();

  const s = {
    en: { play: 'Play lofi music', pause: 'Pause lofi music', playShort: 'Play', pauseShort: 'Pause' },
    fr: { play: 'Jouer la musique lofi', pause: 'Mettre la musique en pause', playShort: 'Jouer', pauseShort: 'Pause' },
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
      <circle cx="5" cy="6.6" r="0.4" fill="var(--clr-ink)" />
      <circle cx="19" cy="6.6" r="0.4" fill="var(--clr-ink)" />
      <circle cx="8" cy="12" r="3.4" fill="var(--clr-ink)" />
      <circle cx="16" cy="12" r="3.4" fill="var(--clr-ink)" />
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
    <span class="music__label label">{playing ? s.pauseShort : s.playShort}</span>
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

  /* pill shape is a deliberate one-off exception to the site's sharp-corner
     rule — this is a media control, not a structural container, and the
     rounded capsule is what visually separates it from the "Montréal"
     label sitting right next to it */
  .music__box {
    display: flex;
    align-items: center;
    gap: 0.45rem;
    height: 1.5rem;
    padding: 0 0.7rem 0 0.35rem;
    background: var(--clr-ink);
    border-radius: 999px;
    color: var(--clr-bg);
    transition: opacity var(--transition);
  }

  .music:hover .music__box,
  .music:focus-visible .music__box {
    opacity: 0.8;
  }

  .music__icon {
    width: 1.05rem;
    height: 1.05rem;
    flex-shrink: 0;
  }

  .music__label {
    color: inherit;
    line-height: 1;
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
    .music__box, .music__label { transition: none; }
  }
</style>
