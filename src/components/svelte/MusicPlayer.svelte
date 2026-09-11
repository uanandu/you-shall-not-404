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
    <svg class="music__disk" viewBox="0 0 24 24" aria-hidden="true">
      <circle cx="12" cy="12" r="10" fill="currentColor" />
      <circle cx="12" cy="12" r="8" fill="none" stroke="var(--clr-bg)" stroke-width="0.5" />
      <circle cx="12" cy="12" r="6.2" fill="none" stroke="var(--clr-bg)" stroke-width="0.5" />
      <circle cx="12" cy="12" r="4.6" fill="none" stroke="var(--clr-bg)" stroke-width="0.5" />
      <circle cx="12" cy="12" r="3" fill="var(--clr-rule)" />
      <circle cx="12" cy="12" r="0.7" fill="var(--clr-bg)" />
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

  .music__disk {
    width: 1.05rem;
    height: 1.05rem;
    animation: music-spin 3s linear infinite;
    animation-play-state: paused;
  }

  .music--playing .music__disk { animation-play-state: running; }

  @keyframes music-spin {
    to { transform: rotate(360deg); }
  }

  @media (prefers-reduced-motion: reduce) {
    .music__disk { animation: none; }
    .music__box { transition: none; }
  }
</style>
