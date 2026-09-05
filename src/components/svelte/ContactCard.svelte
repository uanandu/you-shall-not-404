<script lang="ts">
  import { fade } from 'svelte/transition';

  interface Strings {
    q1: string; q2: string; q3: string; q4: string;
    opt1: string; opt2: string; opt3: string; opt4: string;
    ph1: string; ph2: string; ph4: string;
    back: string; next: string;
    submit: string; sending: string; sent: string;
    err_required: string; err_email: string; err_min: string;
  }

  interface Props { strings?: Partial<Strings>; }

  const defaults: Strings = {
    q1: "What's your name?",
    q2: "What's your email?",
    q3: "What are you working on?",
    q4: "Tell me about it.",
    opt1: "New product",
    opt2: "Infrastructure / DevOps",
    opt3: "Consulting",
    opt4: "Something else",
    ph1: "Your name",
    ph2: "you@example.com",
    ph4: "As much or as little as you'd like…",
    back: "← Back",
    next: "Next →",
    submit: "Send message →",
    sending: "Sending…",
    sent: "Message received. I'll be in touch.",
    err_required: "Required.",
    err_email: "Invalid email.",
    err_min: "At least 10 characters.",
  };

  const { strings: overrides = {} }: Props = $props();
  const s: Strings = { ...defaults, ...overrides };

  const TOTAL = 4;
  let step   = $state(0);
  let status = $state<'idle' | 'sending' | 'sent'>('idle');

  let name    = $state('');
  let email   = $state('');
  let topic   = $state('');
  let message = $state('');
  let error   = $state('');

  const questions = $derived([s.q1, s.q2, s.q3, s.q4]);
  const options   = $derived([s.opt1, s.opt2, s.opt3, s.opt4]);

  function validate(): boolean {
    error = '';
    if (step === 0 && !name.trim()) { error = s.err_required; return false; }
    if (step === 1) {
      if (!email.trim()) { error = s.err_required; return false; }
      if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) { error = s.err_email; return false; }
    }
    if (step === 2 && !topic) { error = s.err_required; return false; }
    if (step === 3 && message.trim().length < 10) { error = s.err_min; return false; }
    return true;
  }

  function next() { if (validate()) step = Math.min(step + 1, TOTAL); }
  function back() { error = ''; step = Math.max(step - 1, 0); }

  function handleKey(e: KeyboardEvent) {
    if (e.key === 'Enter') { e.preventDefault(); next(); }
  }

  async function submit() {
    status = 'sending';
    await new Promise(r => setTimeout(r, 1000));
    status = 'sent';
  }

  function focusOnMount(node: HTMLElement) {
    node.focus();
  }
</script>

{#if status === 'sent'}
  <div in:fade={{ duration: 250 }}>
    <p class="sent-msg">{s.sent}</p>
  </div>

{:else if step < TOTAL}
  {#key step}
    <div class="card" in:fade={{ duration: 180, delay: 80 }} out:fade={{ duration: 120 }}>

      <div class="card__head">
        <span class="label card__counter">
          {String(step + 1).padStart(2, '0')} / {String(TOTAL).padStart(2, '0')}
        </span>
        <div class="card__bar" aria-hidden="true">
          <div class="card__bar-fill" style="width: {((step + 1) / TOTAL) * 100}%"></div>
        </div>
      </div>

      <p class="card__q">{questions[step]}</p>

      <div class="card__answer">
        {#if step === 0}
          <input class="card__input" type="text" bind:value={name}
            placeholder={s.ph1} onkeydown={handleKey} use:focusOnMount />
        {:else if step === 1}
          <input class="card__input" type="email" bind:value={email}
            placeholder={s.ph2} onkeydown={handleKey} use:focusOnMount />
        {:else if step === 2}
          <div class="card__opts" role="group">
            {#each options as opt}
              <button type="button" class="card__opt label"
                class:card__opt--on={topic === opt}
                onclick={() => { topic = opt; error = ''; }}>
                {opt}
              </button>
            {/each}
          </div>
        {:else}
          <textarea class="card__textarea" bind:value={message}
            placeholder={s.ph4} rows="4" use:focusOnMount></textarea>
        {/if}
      </div>

      {#if error}
        <p class="card__err label" role="alert" in:fade={{ duration: 150 }}>{error}</p>
      {/if}

      <div class="card__nav">
        {#if step > 0}
          <button type="button" class="card__back label" onclick={back}>{s.back}</button>
        {:else}
          <span></span>
        {/if}
        <button type="button" class="card__next label" onclick={next}>{s.next}</button>
      </div>

    </div>
  {/key}

{:else}
  <div class="review" in:fade={{ duration: 180 }}>

    <p class="label review__label">Review</p>

    <dl class="review__dl">
      {#each [
        { label: s.q1, value: name },
        { label: s.q2, value: email },
        { label: s.q3, value: topic },
        { label: s.q4, value: message },
      ] as row}
        <div class="review__row">
          <dt class="label review__dt">{row.label}</dt>
          <dd class="review__dd">{row.value}</dd>
        </div>
      {/each}
    </dl>

    <div class="card__nav">
      <button type="button" class="card__back label" onclick={back}>{s.back}</button>
      <button type="button" class="card__next label"
        onclick={submit} disabled={status === 'sending'}>
        {status === 'sending' ? s.sending : s.submit}
      </button>
    </div>

  </div>
{/if}

<style>
  /* ── Progress header ─────────────────────────────────────── */
  .card__head {
    display: flex;
    flex-direction: column;
    gap: 0.65rem;
    margin-bottom: 2rem;
  }

  .card__counter { color: var(--clr-muted); }

  .card__bar {
    height: 1px;
    background: var(--clr-rule);
    position: relative;
  }

  .card__bar-fill {
    position: absolute;
    inset-block: 0;
    left: 0;
    background: var(--clr-ink);
    transition: width 0.38s cubic-bezier(0.4, 0, 0.2, 1);
  }

  /* ── Question ─────────────────────────────────────────────── */
  .card__q {
    font-family: var(--font-display);
    font-size: clamp(1.4rem, 3vw, 2rem);
    font-weight: 700;
    letter-spacing: -0.02em;
    line-height: 1.15;
    color: var(--clr-ink);
    margin-bottom: 1.75rem;
  }

  /* ── Inputs ───────────────────────────────────────────────── */
  .card__input,
  .card__textarea {
    display: block;
    width: 100%;
    background: transparent;
    border: none;
    border-bottom: 1px solid var(--clr-rule);
    padding: 0.5rem 0;
    font-family: var(--font-body);
    font-size: 1rem;
    color: var(--clr-ink);
    resize: none;
    transition: border-color var(--transition);
  }

  .card__input::placeholder,
  .card__textarea::placeholder { color: var(--clr-muted); }

  .card__input:focus,
  .card__textarea:focus {
    outline: none;
    border-bottom-color: var(--clr-ink);
  }

  /* ── Options ──────────────────────────────────────────────── */
  .card__opts {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.5rem;
  }

  .card__opt {
    padding: 0.75rem 1rem;
    border: 1px solid var(--clr-rule);
    background: transparent;
    color: var(--clr-muted);
    text-align: left;
    cursor: pointer;
    transition: border-color var(--transition), color var(--transition);
  }

  .card__opt:hover {
    border-color: var(--clr-ink-2);
    color: var(--clr-ink-2);
  }

  .card__opt--on {
    border-color: var(--clr-ink) !important;
    color: var(--clr-ink) !important;
  }

  /* ── Error ────────────────────────────────────────────────── */
  .card__err {
    margin-top: 0.6rem;
    color: var(--clr-error);
  }

  /* ── Navigation ───────────────────────────────────────────── */
  .card__nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 2rem;
    padding-top: 1.25rem;
    border-top: 1px solid var(--clr-rule);
  }

  .card__back {
    background: none;
    border: none;
    cursor: pointer;
    color: var(--clr-muted);
    font-family: var(--font-body);
    padding: 0;
    transition: color var(--transition);
  }

  .card__back:hover { color: var(--clr-ink); }

  .card__next {
    background: none;
    border: none;
    border-bottom: 1px solid var(--clr-ink);
    cursor: pointer;
    color: var(--clr-ink);
    font-family: var(--font-body);
    padding: 0 0 0.15rem;
    transition: opacity var(--transition);
  }

  .card__next:hover { opacity: 0.5; }
  .card__next:disabled { opacity: 0.4; cursor: not-allowed; }

  /* ── Review ───────────────────────────────────────────────── */
  .review__label {
    color: var(--clr-muted);
    display: block;
    margin-bottom: 1rem;
  }

  .review__dl {
    display: flex;
    flex-direction: column;
    border-top: 1px solid var(--clr-rule);
  }

  .review__row {
    display: grid;
    grid-template-columns: 9rem 1fr;
    gap: 1rem;
    padding-block: 0.85rem;
    border-bottom: 1px solid var(--clr-rule);
    align-items: baseline;
  }

  .review__dt { color: var(--clr-muted); }

  .review__dd {
    color: var(--clr-ink);
    font-size: 0.95rem;
    line-height: 1.6;
    word-break: break-word;
  }

  /* ── Sent ─────────────────────────────────────────────────── */
  .sent-msg {
    font-family: var(--font-display);
    font-style: italic;
    font-size: clamp(1.25rem, 3vw, 2rem);
    color: var(--clr-ink);
  }

  @media (max-width: 640px) {
    .card__opts { grid-template-columns: 1fr; }
    .review__row { grid-template-columns: 1fr; gap: 0.2rem; }
  }
</style>
