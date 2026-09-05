<script lang="ts">
  import { tick } from 'svelte';

  type Line = { kind: 'prompt' | 'output' | 'error'; text: string };

  let lines  = $state<Line[]>([]);
  let input  = $state('');
  let ready  = $state(false);
  let status = $state<'idle' | 'loading' | 'ready' | 'error'>('idle');
  let cwd    = $state('/home/anandu');

  let history: string[] = [];
  let histIdx = -1;
  let stashedInput = '';

  let inputEl:  HTMLInputElement;
  let scrollEl: HTMLDivElement;

  function shortCwd(p: string) {
    return p.replace('/home/anandu', '~');
  }

  function ps1Text() {
    return `anandu@portfolio:${shortCwd(cwd)}$ `;
  }

  function syncCwd() {
    const w = window as any;
    if (typeof w.wasmGetCwd === 'function') cwd = w.wasmGetCwd();
  }

  function esc(s: string) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  // Render a stored prompt-line string as coloured spans
  function renderPromptLine(text: string): string {
    const di = text.indexOf('$ ');
    if (di === -1) return `<span class="ps1__cmd">${esc(text)}</span>`;
    const prefix = text.slice(0, di);
    const cmd    = text.slice(di + 2);
    const ci = prefix.lastIndexOf(':');
    const ai = prefix.indexOf('@');
    if (ai === -1 || ci === -1) return esc(text);
    const user = prefix.slice(0, ai);
    const host = prefix.slice(ai + 1, ci);
    const path = prefix.slice(ci + 1);
    return (
      `<span class="ps1__user">${esc(user)}</span>` +
      `<span class="ps1__at">@${esc(host)}</span>` +
      `<span class="ps1__sep">:</span>` +
      `<span class="ps1__path">${esc(path)}</span>` +
      `<span class="ps1__dollar">$</span>` +
      `<span class="ps1__cmd"> ${esc(cmd)}</span>`
    );
  }

  function addLine(kind: Line['kind'], text: string) {
    lines = [...lines, { kind, text }];
  }

  async function scrollBottom() {
    await tick();
    if (scrollEl) scrollEl.scrollTop = scrollEl.scrollHeight;
  }

  function runCommand(cmd: string) {
    const w = window as any;
    if (typeof w.wasmRunCommand !== 'function') return;
    const r = w.wasmRunCommand(cmd);
    syncCwd();
    if (r.clear) {
      lines = [];
    } else if (r.output) {
      addLine(r.err ? 'error' : 'output', r.output);
    }
  }

  async function submit() {
    const cmd = input.trim();
    addLine('prompt', ps1Text() + input);
    input = '';
    histIdx = -1;
    stashedInput = '';
    if (cmd) {
      if (history.at(-1) !== cmd) history = [...history, cmd];
      runCommand(cmd);
    }
    await scrollBottom();
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter') {
      e.preventDefault();
      submit();
      return;
    }
    if (e.key === 'ArrowUp') {
      e.preventDefault();
      if (history.length === 0) return;
      if (histIdx === -1) { stashedInput = input; histIdx = history.length - 1; }
      else if (histIdx > 0) histIdx--;
      input = history[histIdx];
      return;
    }
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (histIdx === -1) return;
      if (histIdx < history.length - 1) { histIdx++; input = history[histIdx]; }
      else { histIdx = -1; input = stashedInput; }
      return;
    }
    if (e.key === 'Tab') {
      e.preventDefault();
      const w = window as any;
      if (typeof w.wasmGetCompletions !== 'function') return;
      const completions: string[] = w.wasmGetCompletions(input);
      if (completions.length === 1) {
        const tokens = input.split(/\s+/);
        tokens[tokens.length - 1] = completions[0];
        input = tokens.join(' ');
      } else if (completions.length > 1) {
        addLine('output', completions.join('  '));
        scrollBottom();
      }
      return;
    }
    if (e.key === 'c' && e.ctrlKey) {
      e.preventDefault();
      addLine('prompt', ps1Text() + input + '^C');
      input = '';
      histIdx = -1;
      scrollBottom();
      return;
    }
    if (e.key === 'l' && e.ctrlKey) {
      e.preventDefault();
      lines = [];
      return;
    }
  }

  function focusInput() { inputEl?.focus(); }

  async function boot() {
    if (status !== 'idle') return;
    status = 'loading';

    if (!(window as any).Go) {
      try {
        await new Promise<void>((resolve, reject) => {
          const s = document.createElement('script');
          s.src = '/wasm_exec.js';
          s.onload = () => resolve();
          s.onerror = () => reject();
          document.head.appendChild(s);
        });
      } catch {
        status = 'error';
        addLine('error', 'wasm_exec.js not found.\nRun:  cd wasm/shell && make build');
        return;
      }
    }
    try {
      const go = new (window as any).Go();
      const res = await WebAssembly.instantiateStreaming(fetch('/shell.wasm'), go.importObject);
      go.run(res.instance);
      status = 'ready';
      ready  = true;
      syncCwd();
      runCommand('__welcome__');
      await scrollBottom();
      inputEl?.focus();
    } catch {
      status = 'error';
      addLine('error', 'shell.wasm not found.\nBuild it first:\n\n  cd wasm/shell && make build');
    }
  }

  function onTerminalClick() {
    if (status === 'idle') boot();
    else focusInput();
  }

  function onTerminalKeydown(e: KeyboardEvent) {
    if (status === 'idle' && (e.key === 'Enter' || e.key === ' ')) {
      e.preventDefault();
      boot();
    }
  }
</script>

<div
  class="terminal"
  onclick={onTerminalClick}
  onkeydown={onTerminalKeydown}
  role={status === 'idle' ? 'button' : undefined}
  tabindex={status === 'idle' ? 0 : undefined}
  aria-label={status === 'idle' ? 'Start the interactive shell (loads a 3MB WASM binary)' : undefined}
>

  <!-- ── Chrome bar ──────────────────────────────────────────── -->
  <div class="terminal__chrome">
    <span class="terminal__chrome-path">
      {#if status === 'ready'}{shortCwd(cwd)}{:else if status === 'error'}error{:else if status === 'idle'}click to start{:else}&nbsp;{/if}
    </span>
    <span class="terminal__chrome-shell">gosh · WASM/Go</span>
  </div>

  <!-- ── Output body ─────────────────────────────────────────── -->
  <div class="terminal__body" bind:this={scrollEl}>

    {#if status === 'idle'}
      <div class="terminal__status">
        <span class="terminal__loading-text">Click to load the shell (~3MB WASM) →</span>
      </div>
    {:else if status === 'loading'}
      <div class="terminal__status">
        <span class="terminal__loading-text">Loading shell.wasm</span>
        <span class="terminal__dot terminal__dot--1">.</span><span
             class="terminal__dot terminal__dot--2">.</span><span
             class="terminal__dot terminal__dot--3">.</span>
      </div>
    {/if}

    {#each lines as line, i}
      {#if line.kind === 'prompt'}
        <!-- Colourised prompt history line -->
        <div
          class="terminal__line terminal__line--prompt"
          class:terminal__line--group-start={i > 0}
        >{@html renderPromptLine(line.text)}</div>
      {:else}
        <pre
          class="terminal__line terminal__line--{line.kind}"
        >{line.text}</pre>
      {/if}
    {/each}

    {#if ready}
      <!-- Live input row with colourised prompt parts -->
      <div
        class="terminal__input-row"
        class:terminal__input-row--first={lines.length === 0}
      >
        <span class="terminal__ps1" aria-hidden="true">
          <span class="ps1__user">anandu</span><span
               class="ps1__at">@portfolio</span><span
               class="ps1__sep">:</span><span
               class="ps1__path">{shortCwd(cwd)}</span><span
               class="ps1__dollar">$</span>
        </span>
        <input
          bind:this={inputEl}
          bind:value={input}
          class="terminal__input"
          type="text"
          spellcheck="false"
          autocomplete="off"
          autocapitalize="off"
          onkeydown={onKeydown}
          aria-label="Shell input — {ps1Text()}"
        />
      </div>
    {/if}

  </div>
</div>

<style>
  /* ── Outer shell ──────────────────────────────────────────── */
  .terminal {
    display: flex;
    flex-direction: column;
    height: 480px;
    cursor: text;
    border-top: 1px solid var(--clr-rule);
  }

  .terminal[role="button"] { cursor: pointer; }

  /* ── Chrome bar ───────────────────────────────────────────── */
  .terminal__chrome {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.45rem 0.85rem;
    /* left indent clears the absolute-positioned module__tag */
    padding-left: 7.5rem;
    background: var(--clr-surface);
    border-bottom: 1px solid var(--clr-rule);
    flex-shrink: 0;
  }

  .terminal__chrome-path {
    font-family: var(--font-mono);
    font-size: 0.65rem;
    color: var(--clr-ink-2);
    letter-spacing: 0.03em;
  }

  .terminal__chrome-shell {
    font-family: var(--font-mono);
    font-size: 0.62rem;
    color: var(--clr-muted);
    letter-spacing: 0.04em;
  }

  /* ── Body ─────────────────────────────────────────────────── */
  .terminal__body {
    flex: 1;
    overflow-y: auto;
    padding: 0.85rem 0.85rem;
    display: flex;
    flex-direction: column;
  }

  /* ── Lines ────────────────────────────────────────────────── */
  .terminal__line {
    font-family: var(--font-mono);
    font-size: 0.72rem;
    line-height: 1.6;
    color: var(--clr-ink-2);
    white-space: pre-wrap;
    word-break: break-word;
    margin: 0;
  }

  /* Gap above each new command block */
  .terminal__line--group-start { margin-top: 0.65rem; }

  .terminal__line--prompt { color: var(--clr-ink); }
  .terminal__line--error  { color: var(--clr-error); }

  /* ── Prompt colour parts (history lines & live row) ────────── */
  :global(.ps1__user)   { color: var(--clr-ink); }
  :global(.ps1__at)     { color: var(--clr-muted); }
  :global(.ps1__sep)    { color: var(--clr-muted); }
  :global(.ps1__path)   { color: var(--clr-ink-2); }
  :global(.ps1__dollar) { color: var(--clr-muted); }
  :global(.ps1__cmd)    { color: var(--clr-ink); }

  /* ── Live input row ───────────────────────────────────────── */
  .terminal__input-row {
    display: flex;
    align-items: baseline;
    margin-top: 0.65rem;
  }

  .terminal__input-row--first { margin-top: 0; }

  .terminal__ps1 {
    font-family: var(--font-mono);
    font-size: 0.72rem;
    white-space: nowrap;
    flex-shrink: 0;
    line-height: 1.6;
  }

  .terminal__input {
    flex: 1;
    background: none;
    border: none;
    outline: none;
    font-family: var(--font-mono);
    font-size: 0.72rem;
    line-height: 1.6;
    color: var(--clr-ink);
    caret-color: var(--clr-ink);
    /* space between $ and typed text */
    padding: 0 0 0 0.45em;
    min-width: 0;
  }

  /* ── Loading ──────────────────────────────────────────────── */
  .terminal__status {
    display: flex;
    align-items: baseline;
  }

  .terminal__loading-text {
    font-family: var(--font-mono);
    font-size: 0.72rem;
    color: var(--clr-muted);
  }

  .terminal__dot {
    font-family: var(--font-mono);
    font-size: 0.72rem;
    color: var(--clr-muted);
    animation: dot-blink 1.4s ease-in-out infinite;
  }

  .terminal__dot--1 { animation-delay: 0s;    }
  .terminal__dot--2 { animation-delay: 0.22s; }
  .terminal__dot--3 { animation-delay: 0.44s; }

  @keyframes dot-blink {
    0%, 80%, 100% { opacity: 0.2; }
    40%           { opacity: 1;   }
  }

  @media (prefers-reduced-motion: reduce) {
    .terminal__dot { animation: none; opacity: 1; }
  }

  @media (max-width: 640px) {
    .terminal { height: 320px; }
  }
</style>
