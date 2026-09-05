<script lang="ts">
  import { load as yamlLoad } from 'js-yaml';
  import Terminal from './Terminal.svelte';

  const { lang = 'en' }: { lang?: 'en' | 'fr' } = $props();

  const s = {
    en: {
      tag04:       '04 — K8s',
      hint01:      'type help to start',
      hintRegex:   'type a pattern',
      hintStack:   'type your stack',
      hint04:      'paste a manifest',
      regexAria:   'Regular expression pattern',
      samplesAria: 'Test strings',
      matchedAria: 'Matched skills',
      gapsAria:    'Skills not in my stack',
      stackPh:     'Go, Kubernetes, Postgres…',
      stackAria:   'Enter your tech stack',
      stackEmpty:  'Your stack against mine.',
      matches:     (n: number) => `${n} match${n === 1 ? '' : 'es'}`,
      noMatches:   'no matches',
    },
    fr: {
      tag04:       '04 — K8s',
      hint01:      'tapez help pour débuter',
      hintRegex:   'entrez un motif',
      hintStack:   'entrez votre pile',
      hint04:      'collez un manifeste',
      regexAria:   'Expression régulière',
      samplesAria: 'Chaînes de test',
      matchedAria: 'Compétences correspondantes',
      gapsAria:    'Compétences absentes de ma pile',
      stackPh:     'Go, Kubernetes, Postgres…',
      stackAria:   'Entrez votre pile technologique',
      stackEmpty:  'Votre pile contre la mienne.',
      matches:     (n: number) => `${n} correspondance${n === 1 ? '' : 's'}`,
      noMatches:   'aucune correspondance',
    },
  } as const;

  const t = s[lang];

  // ── 02  Regex tester ────────────────────────────────────────
  const SAMPLES = [
    '/api/v1/users/42/posts',
    'GET /health HTTP/1.1 200 OK',
    '2026-06-24T15:04:05Z ERROR connection refused',
    'my-deployment-7d9f8c-x4k2m',
    'v1.24.3+build.2026',
    '550e8400-e29b-41d4-a716-446655440000',
  ];

  let regexInput = $state('');

  type RegexState = { ok: true; re: RegExp | null } | { ok: false; error: string };

  const regexState = $derived<RegexState>((() => {
    const pat = regexInput.trim();
    if (!pat) return { ok: true, re: null };
    try { return { ok: true, re: new RegExp(pat, 'g') }; }
    catch (e: any) { return { ok: false, error: e.message }; }
  })());

  const totalMatches = $derived(
    regexState.ok && regexState.re
      ? SAMPLES.reduce((n, s) => n + (s.match(new RegExp((regexState.re as RegExp).source, 'g'))?.length ?? 0), 0)
      : 0
  );

  function esc(s: string) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  function highlight(sample: string, re: RegExp | null): string {
    if (!re) return esc(sample);
    const g = new RegExp(re.source, 'g');
    let out = '', last = 0, m: RegExpExecArray | null;
    while ((m = g.exec(sample)) !== null) {
      out += esc(sample.slice(last, m.index));
      out += `<mark class="regex__mark">${esc(m[0])}</mark>`;
      last = g.lastIndex;
      if (m[0].length === 0) g.lastIndex++;
    }
    return out + esc(sample.slice(last));
  }

  // ── 03  Stack matcher ───────────────────────────────────────
  const MY_SKILLS = [
    { label: 'Go',         keys: ['go', 'golang'] },
    { label: 'Rust',       keys: ['rust'] },
    { label: 'TypeScript', keys: ['typescript', 'ts'] },
    { label: 'JavaScript', keys: ['javascript', 'js'] },
    { label: 'Node.js',    keys: ['node', 'nodejs', 'node.js'] },
    { label: 'Docker',     keys: ['docker'] },
    { label: 'Kubernetes', keys: ['kubernetes', 'k8s'] },
    { label: 'Helm',       keys: ['helm'] },
    { label: 'ArgoCD',     keys: ['argocd', 'argo'] },
    { label: 'Azure',      keys: ['azure'] },
    { label: 'Svelte',     keys: ['svelte'] },
    { label: 'Astro',      keys: ['astro'] },
    { label: 'Shopify',    keys: ['shopify', 'liquid'] },
    { label: 'PostgreSQL', keys: ['postgresql', 'postgres', 'psql', 'sql'] },
  ];

  let stackInput = $state('');

  function parseTokens(raw: string): string[] {
    return raw.split(/[\s,;/|]+/).map(s => s.toLowerCase().trim()).filter(s => s.length > 1);
  }

  const tokens = $derived(parseTokens(stackInput));
  const matched = $derived(MY_SKILLS.filter(skill => tokens.some(t => skill.keys.includes(t))));
  const gaps    = $derived(tokens.filter(t => !MY_SKILLS.some(s => s.keys.includes(t))));
  const pct     = $derived(tokens.length > 0 ? Math.round((matched.length / tokens.length) * 100) : null);

  // ── 04  Kubernetes manifest scorer ─────────────────────────
  const SAMPLE_MANIFEST = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-server
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: api-server
  template:
    metadata:
      labels:
        app: api-server
    spec:
      containers:
      - name: api-server
        image: nginx:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_PASSWORD
          value: "hunter2"`;

  interface Check {
    id:          string;
    title:       string;
    description: string;
    severity:    'error' | 'warning';
    pass:        boolean;
  }

  interface K8sResult {
    checks:  Check[];
    passed:  number;
    total:   number;
    score:   number;
    error:   string | null;
  }

  function containers(m: any): any[] {
    return m?.spec?.template?.spec?.containers ?? m?.spec?.containers ?? [];
  }

  function runChecks(m: any): Check[] {
    const ctrs = containers(m);
    const SECRET_RE = /password|secret|token|key|credential|api_key/i;
    return [
      {
        id: 'image-tag',
        title: 'Image tag pinned',
        description: 'Avoid :latest — use a specific version for reproducible deployments.',
        severity: 'error',
        pass: ctrs.length > 0 && ctrs.every((c: any) =>
          c.image && !c.image.endsWith(':latest') && c.image.includes(':')
        ),
      },
      {
        id: 'resource-limits',
        title: 'Resource limits set',
        description: 'CPU and memory limits prevent runaway containers.',
        severity: 'error',
        pass: ctrs.every((c: any) => c.resources?.limits?.cpu && c.resources?.limits?.memory),
      },
      {
        id: 'resource-requests',
        title: 'Resource requests set',
        description: 'Requests let the scheduler place pods correctly.',
        severity: 'warning',
        pass: ctrs.every((c: any) => c.resources?.requests?.cpu && c.resources?.requests?.memory),
      },
      {
        id: 'liveness',
        title: 'Liveness probe defined',
        description: 'Restarts containers that are stuck or deadlocked.',
        severity: 'warning',
        pass: ctrs.every((c: any) => c.livenessProbe != null),
      },
      {
        id: 'readiness',
        title: 'Readiness probe defined',
        description: 'Prevents traffic reaching the container before it is ready.',
        severity: 'warning',
        pass: ctrs.every((c: any) => c.readinessProbe != null),
      },
      {
        id: 'replicas',
        title: 'Replicas > 1',
        description: 'A single replica is a single point of failure.',
        severity: 'warning',
        pass: (m?.spec?.replicas ?? 1) > 1,
      },
      {
        id: 'non-root',
        title: 'Non-root user',
        description: 'Running as root inside a container is a security risk.',
        severity: 'error',
        pass: ctrs.every((c: any) =>
          c.securityContext?.runAsNonRoot === true ||
          (c.securityContext?.runAsUser != null && c.securityContext.runAsUser !== 0)
        ),
      },
      {
        id: 'namespace',
        title: 'Custom namespace',
        description: 'Avoid the default namespace in production workloads.',
        severity: 'warning',
        pass: !!m?.metadata?.namespace && m.metadata.namespace !== 'default',
      },
      {
        id: 'secret-env',
        title: 'No secrets in env vars',
        description: 'Use Kubernetes Secrets or an external vault instead.',
        severity: 'error',
        pass: ctrs.every((c: any) =>
          (c.env ?? []).every((e: any) => !e.value || !SECRET_RE.test(e.name))
        ),
      },
      {
        id: 'labels',
        title: 'App label present',
        description: 'Labels are required for selectors and observability tooling.',
        severity: 'warning',
        pass: !!(m?.metadata?.labels?.app ?? m?.spec?.template?.metadata?.labels?.app),
      },
    ];
  }

  let k8sInput = $state(SAMPLE_MANIFEST);

  const k8sResult = $derived<K8sResult | null>((() => {
    if (!k8sInput.trim()) return null;
    try {
      const m = yamlLoad(k8sInput) as any;
      if (!m?.apiVersion || !m?.kind) {
        return { checks: [], passed: 0, total: 0, score: 0, error: 'Missing apiVersion or kind — paste a Kubernetes manifest.' };
      }
      const checks = runChecks(m);
      const passed = checks.filter(c => c.pass).length;
      return { checks, passed, total: checks.length, score: Math.round((passed / checks.length) * 100), error: null };
    } catch (e: any) {
      return { checks: [], passed: 0, total: 0, score: 0, error: `YAML error: ${e.message}` };
    }
  })());
</script>

<!-- ── 01 TERMINAL ───────────────────────────────────────────── -->
<div class="module">
  <span class="module__tag label">01 — Terminal</span>
  <div class="term-wrap">
    <Terminal />
  </div>
  <span class="module__hint label">{t.hint01}</span>
</div>

<!-- ── BOTTOM ROW ────────────────────────────────────────────── -->
<div class="bottom">

  <!-- 02 Regex -->
  <div class="module">
    <span class="module__tag label">02 — Regex</span>
    <div class="regex">
      <input
        class="regex__input label"
        type="text"
        bind:value={regexInput}
        placeholder="\d+"
        aria-label={t.regexAria}
        spellcheck="false"
        autocomplete="off"
      />
      {#if !regexState.ok}
        <p class="regex__error label">{regexState.error}</p>
      {:else}
        <ul class="regex__samples" aria-label={t.samplesAria}>
          {#each SAMPLES as sample}
            <li class="regex__sample">{@html highlight(sample, regexState.re)}</li>
          {/each}
        </ul>
        {#if totalMatches > 0}
          <p class="regex__score label">{t.matches(totalMatches)}</p>
        {:else if regexState.re}
          <p class="regex__score regex__score--none label">{t.noMatches}</p>
        {/if}
      {/if}
    </div>
    <span class="module__hint label">{t.hintRegex}</span>
  </div>

  <!-- 03 Stack matcher -->
  <div class="module">
    <span class="module__tag label">03 — Stack</span>
    <div class="matcher">
      <input
        class="matcher__input label"
        type="text"
        bind:value={stackInput}
        placeholder={t.stackPh}
        aria-label={t.stackAria}
        spellcheck="false"
      />
      {#if tokens.length > 0}
        <div class="matcher__results">
          <p class="matcher__score label">
            {matched.length} / {tokens.length} match
            {#if pct !== null}<span class="matcher__pct">— {pct}%</span>{/if}
          </p>
          {#if matched.length > 0}
            <ul class="matcher__list" aria-label={t.matchedAria}>
              {#each matched as skill}
                <li class="matcher__item matcher__item--match label">{skill.label}</li>
              {/each}
            </ul>
          {/if}
          {#if gaps.length > 0}
            <ul class="matcher__list" aria-label={t.gapsAria}>
              {#each gaps as gap}
                <li class="matcher__item matcher__item--gap label">{gap}</li>
              {/each}
            </ul>
          {/if}
        </div>
      {:else}
        <p class="matcher__empty label">{t.stackEmpty}</p>
      {/if}
    </div>
    <span class="module__hint label">{t.hintStack}</span>
  </div>

</div>

<!-- ── 04 K8S MANIFEST SCORER ────────────────────────────────── -->
<div class="module">
  <div class="module__header">
    <span class="module__tag label" style="position:static;">{t.tag04}</span>
    <span class="k8s__title label">Manifest Scorer</span>
  </div>
  <div class="k8s">

    <textarea
      class="k8s__textarea"
      bind:value={k8sInput}
      spellcheck="false"
      autocomplete="off"
      aria-label="Kubernetes manifest YAML"
    ></textarea>

    <div class="k8s__panel">
      {#if k8sResult}
        {#if k8sResult.error}
          <p class="k8s__error label">{k8sResult.error}</p>
        {:else}
          <ul class="k8s__checks">
            {#each k8sResult.checks as check}
              <li class="k8s__check" class:k8s__check--pass={check.pass}>
                <span class="k8s__icon">{check.pass ? '✓' : check.severity === 'error' ? '✗' : '▲'}</span>
                <span class="k8s__check-body">
                  <span class="k8s__check-title">{check.title}</span>
                  {#if !check.pass}
                    <span class="k8s__check-desc">{check.description}</span>
                  {/if}
                </span>
                {#if !check.pass}
                  <span class="k8s__badge k8s__badge--{check.severity} label">{check.severity}</span>
                {/if}
              </li>
            {/each}
          </ul>
          <div class="k8s__score">
            <span class="k8s__score-pct">{k8sResult.score}%</span>
            <span class="k8s__score-sub label">{k8sResult.passed} / {k8sResult.total} checks passed</span>
          </div>
        {/if}
      {/if}
    </div>

  </div>
  <span class="module__hint label">{t.hint04}</span>
</div>

<style>
  /* ── Module shell ─────────────────────────────────────────── */
  .module {
    position: relative;
    border: 1px solid var(--clr-rule);
    overflow: hidden;
    margin-bottom: 1.5rem;
  }

  .module__tag {
    position: absolute;
    top: 0.7rem;
    left: 0.85rem;
    color: var(--clr-muted);
    z-index: 2;
    pointer-events: none;
  }

  .module__hint {
    position: absolute;
    bottom: 0.7rem;
    right: 0.85rem;
    color: var(--clr-muted);
    pointer-events: none;
  }

  /* ── Bottom row ───────────────────────────────────────────── */
  .bottom {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
  }

  @media (max-width: 640px) {
    .bottom { grid-template-columns: 1fr; }
  }

  /* ── Terminal ─────────────────────────────────────────────── */
  .term-wrap { margin-top: 2rem; }

  /* ── Regex tester ─────────────────────────────────────────── */
  .regex {
    padding: 2.5rem 0.85rem 2.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.9rem;
    min-height: 9rem;
  }

  .regex__input {
    background: none;
    border: none;
    border-bottom: 1px solid var(--clr-rule);
    padding: 0.3rem 0;
    font-family: var(--font-mono);
    font-size: 0.75rem;
    color: var(--clr-ink);
    width: 100%;
    transition: border-color var(--transition);
    letter-spacing: 0.04em;
  }

  .regex__input::placeholder { color: var(--clr-muted); }
  .regex__input:focus { outline: none; border-bottom-color: var(--clr-ink); }

  .regex__samples { display: flex; flex-direction: column; gap: 0.3rem; }

  .regex__sample {
    font-family: var(--font-mono);
    font-size: 0.68rem;
    line-height: 1.55;
    color: var(--clr-ink-2);
    word-break: break-all;
  }

  :global(.regex__mark) {
    background: var(--clr-ink);
    color: var(--clr-bg);
    font-style: normal;
  }

  .regex__error { color: var(--clr-error); font-family: var(--font-mono); font-size: 0.68rem; }
  .regex__score { color: var(--clr-ink); font-size: 0.68rem; letter-spacing: 0.06em; }
  .regex__score--none { color: var(--clr-muted); }

  /* ── Stack matcher ────────────────────────────────────────── */
  .matcher {
    padding: 2.5rem 0.85rem 2.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    min-height: 9rem;
  }

  .matcher__input {
    background: none;
    border: none;
    border-bottom: 1px solid var(--clr-rule);
    padding: 0.3rem 0;
    font-family: var(--font-mono);
    font-size: 0.75rem;
    color: var(--clr-ink);
    width: 100%;
    transition: border-color var(--transition);
    letter-spacing: 0.04em;
  }

  .matcher__input::placeholder { color: var(--clr-muted); }
  .matcher__input:focus { outline: none; border-bottom-color: var(--clr-ink); }

  .matcher__results { display: flex; flex-direction: column; gap: 0.5rem; }
  .matcher__score { color: var(--clr-ink); letter-spacing: 0.06em; }
  .matcher__pct { color: var(--clr-muted); }
  .matcher__empty { color: var(--clr-muted); font-style: italic; }

  .matcher__list { display: flex; flex-wrap: wrap; gap: 0.3rem 0.6rem; }

  .matcher__item {
    font-size: 0.68rem;
    letter-spacing: 0.06em;
    padding: 0.2rem 0.5rem;
    border: 1px solid var(--clr-rule);
  }

  .matcher__item--match { color: var(--clr-ink); border-color: var(--clr-ink); }
  .matcher__item--gap   { color: var(--clr-muted); text-decoration: line-through; }

  /* ── Module header (tag + title inline) ──────────────────── */
  .module__header {
    display: flex;
    align-items: baseline;
    gap: 0.75rem;
    padding: 0.7rem 0.85rem 0;
  }

  /* ── K8s scorer ───────────────────────────────────────────── */
  .k8s {
    display: grid;
    grid-template-columns: 1fr 1fr;
    min-height: 300px;
    padding-top: 0;
  }

  @media (max-width: 680px) {
    .k8s { grid-template-columns: 1fr; }
    .k8s__textarea {
      border-right: none;
      min-height: 180px;
    }
    .k8s__panel {
      border-top: 1px solid var(--clr-rule);
    }
    .module__hint { display: none; }
  }

  .k8s__title {
    color: var(--clr-muted);
  }

  .k8s__textarea {
    width: 100%;
    height: 100%;
    min-height: 280px;
    background: none;
    border: none;
    border-right: 1px solid var(--clr-rule);
    padding: 0.85rem;
    font-family: var(--font-mono);
    font-size: 0.62rem;
    line-height: 1.6;
    color: var(--clr-ink-2);
    resize: none;
    outline: none;
  }

  .k8s__textarea:focus { color: var(--clr-ink); }

  .k8s__panel {
    padding: 0.85rem 0.85rem 2.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.85rem;
    overflow-y: auto;
  }

  .k8s__error {
    color: var(--clr-error);
    font-family: var(--font-mono);
    font-size: 0.68rem;
    line-height: 1.5;
  }

  .k8s__checks {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .k8s__check {
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
  }

  .k8s__icon {
    font-family: var(--font-mono);
    font-size: 0.68rem;
    flex-shrink: 0;
    width: 1rem;
    color: var(--clr-muted);
    margin-top: 0.05rem;
  }

  .k8s__check--pass .k8s__icon { color: var(--clr-ink-2); }

  .k8s__check-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
    min-width: 0;
  }

  .k8s__check-title {
    font-family: var(--font-mono);
    font-size: 0.68rem;
    color: var(--clr-ink);
    line-height: 1.4;
  }

  .k8s__check--pass .k8s__check-title { color: var(--clr-muted); }

  .k8s__check-desc {
    font-family: var(--font-body);
    font-size: 0.62rem;
    color: var(--clr-muted);
    line-height: 1.45;
  }

  .k8s__badge {
    font-size: 0.55rem;
    letter-spacing: 0.08em;
    padding: 0.1rem 0.35rem;
    border: 1px solid var(--clr-rule);
    color: var(--clr-muted);
    flex-shrink: 0;
    align-self: flex-start;
    margin-top: 0.1rem;
  }

  .k8s__badge--error {
    border-color: var(--clr-error);
    color: var(--clr-error);
  }

  .k8s__score {
    display: flex;
    align-items: baseline;
    gap: 0.6rem;
    border-top: 1px solid var(--clr-rule);
    padding-top: 0.75rem;
    margin-top: auto;
  }

  .k8s__score-pct {
    font-family: var(--font-display);
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--clr-ink);
    line-height: 1;
  }

  .k8s__score-sub {
    color: var(--clr-muted);
    font-size: 0.65rem;
  }
</style>
