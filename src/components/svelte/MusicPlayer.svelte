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
  <span class="music__disk" aria-hidden="true">
    <span class="music__hole"></span>
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

  .music__disk {
    position: relative;
    width: 1.1rem;
    height: 1.1rem;
    border-radius: 50%;
    background: var(--clr-ink-2);
    box-shadow: 0 0 0 2px var(--clr-bg), 0 0 0 3px var(--clr-rule);
    animation: music-spin 3s linear infinite;
    animation-play-state: paused;
    transition: background var(--transition);
  }

  .music:hover .music__disk { background: var(--clr-ink); }

  .music--playing .music__disk { animation-play-state: running; }

  .music__hole {
    position: absolute;
    inset: 0;
    margin: auto;
    width: 0.3rem;
    height: 0.3rem;
    border-radius: 50%;
    background: var(--clr-bg);
  }

  @keyframes music-spin {
    to { transform: rotate(360deg); }
  }

  @media (prefers-reduced-motion: reduce) {
    .music__disk { animation: none; }
  }
</style>
