<script lang="ts">
  let active = $state<string | null>(null);

  // A pipeline, not a hub-and-spoke: four tools, each extending trust one
  // layer further out, closing back on whether the plan was actually kept.
  const nodes = [
    {
      id: 'shadow',
      label: 'Shadow Mode',
      sub: "don't trust one action",
      x: 28, y: 150, w: 160, h: 60,
      desc: 'A Claude Code Plugin that intercepts risky tool calls and logs exactly what would have happened instead of executing it. A live toggle command, and a skill that summarizes a session’s shadow log in plain English.',
    },
    {
      id: 'drift',
      label: 'Spec Drift Detector',
      sub: "don't trust one session",
      x: 223, y: 150, w: 160, h: 60,
      desc: 'Long agent sessions quietly drift from the constraints they were given at the start. A Stop hook checks every response against a written invariant checklist — deterministically first, then a real model judgment — and forces a correction before it compounds.',
    },
    {
      id: 'pool',
      label: 'Agent Worker Pool',
      sub: "don't trust unbounded concurrency",
      x: 418, y: 150, w: 160, h: 60,
      desc: 'A fixed-size Go worker pool for running many agent calls at once without melting a rate limit or the machine underneath — bounded backpressure and clean shutdown, benchmarked against the naive unbounded version with real numbers.',
    },
    {
      id: 'conform',
      label: 'Conformance Checker',
      sub: "don't trust the plan was followed",
      x: 613, y: 150, w: 160, h: 60,
      desc: 'OpenSpec’s own validate command checks a change’s file format — never whether the code it produced actually matches. This closes that gap: parses the plan, diffs it against real git history, flags drift deterministically first, a model’s judgment only where that’s not enough.',
    },
  ];

  // Three short hops carry the pipeline forward; one long arc closes the
  // loop back to the plan the first tool started from.
  const edges = [
    { id: 'shadow-drift', from: 'shadow', to: 'drift', d: 'M 188 180 L 223 180', duration: 1, delay: 0 },
    { id: 'drift-pool', from: 'drift', to: 'pool', d: 'M 383 180 L 418 180', duration: 1, delay: 0.3 },
    { id: 'pool-conform', from: 'pool', to: 'conform', d: 'M 578 180 L 613 180', duration: 1, delay: 0.6 },
    {
      id: 'conform-shadow',
      from: 'conform', to: 'shadow',
      d: 'M 693 150 C 693 24 108 24 108 150',
      label: 'closes the loop — did the code match the plan?',
      lx: 400, ly: 18,
      anchor: 'middle',
      duration: 2.6,
      delay: 0,
    },
  ];

  function nodeState(id: string) {
    if (!active) return 'idle';
    if (active === id) return 'self';
    const linked = edges.some(e =>
      (e.from === active || e.to === active) && (e.from === id || e.to === id)
    );
    return linked ? 'linked' : 'dim';
  }

  function edgeState(edge: typeof edges[0]) {
    if (!active) return 'idle';
    return (edge.from === active || edge.to === active) ? 'active' : 'dim';
  }

  const activeNode = $derived(nodes.find(n => n.id === active) ?? null);
</script>

<div class="arch">
  <svg viewBox="0 0 800 260" width="100%" aria-label="Agent Tooling Sprint pipeline diagram" role="img">
    <defs>
      <marker id="arr" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
        <path d="M 0 1 L 5 3 L 0 5 Z" fill="var(--clr-muted)" />
      </marker>
      <marker id="arr-loop" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
        <path d="M 0 1 L 5 3 L 0 5 Z" fill="var(--clr-ink-2)" />
      </marker>
    </defs>

    <!-- ── Edge rails + flow lines ────────────────────────── -->
    {#each edges as edge}
      {@const state = edgeState(edge)}
      {@const isLoop = edge.id === 'conform-shadow'}

      <path
        d={edge.d}
        fill="none"
        stroke="var(--clr-rule)"
        stroke-width="1"
        opacity={state === 'dim' ? 0.25 : 1}
        marker-end={isLoop ? 'url(#arr-loop)' : 'url(#arr)'}
        class="rail"
      />

      <path
        d={edge.d}
        fill="none"
        stroke="var(--clr-ink)"
        stroke-width="1.5"
        stroke-dasharray={isLoop ? '4 60' : '3 20'}
        stroke-dashoffset="0"
        opacity={state === 'dim' ? 0.08 : state === 'active' ? 1 : isLoop ? 0.5 : 0.45}
        class="flow"
        style="animation-duration: {edge.duration}s; animation-delay: -{edge.delay}s;"
      />

      {#if edge.label}
        <text
          x={edge.lx} y={edge.ly}
          text-anchor={edge.anchor}
          class="edge-label"
          opacity={state === 'dim' ? 0.15 : state === 'active' ? 1 : 0.6}
        >{edge.label}</text>
      {/if}
    {/each}

    <!-- ── Nodes ───────────────────────────────────────────── -->
    {#each nodes as node}
      {@const state = nodeState(node.id)}
      <g
        class="node"
        role="button"
        tabindex="0"
        aria-label={node.label}
        onmouseenter={() => active = node.id}
        onmouseleave={() => active = null}
        onfocus={() => active = node.id}
        onblur={() => active = null}
      >
        <rect
          x={node.x} y={node.y}
          width={node.w} height={node.h}
          fill="var(--clr-bg)"
          stroke={state === 'self' ? 'var(--clr-ink)' : 'var(--clr-rule)'}
          stroke-width={state === 'self' ? 1.5 : 1}
          opacity={state === 'dim' ? 0.3 : 1}
          class="node-rect"
        />
        <text
          x={node.x + node.w / 2} y={node.y + 25}
          text-anchor="middle"
          class="node-label"
          opacity={state === 'dim' ? 0.25 : 1}
        >{node.label}</text>
        <text
          x={node.x + node.w / 2} y={node.y + 43}
          text-anchor="middle"
          class="node-sub"
          opacity={state === 'dim' ? 0.15 : state === 'self' ? 0.7 : 0.45}
        >{node.sub}</text>
      </g>
    {/each}
  </svg>

  <!-- Description strip -->
  <div class="arch__desc" class:visible={!!activeNode}>
    <span class="label arch__name">{activeNode?.label ?? ''}</span>
    <p class="arch__body">{activeNode?.desc ?? ''}</p>
  </div>
</div>

<style>
  .arch {
    border: 1px solid var(--clr-rule);
    margin-bottom: 2.5rem;
    padding: 1.5rem 1.5rem 0;
  }

  svg { display: block; overflow: visible; }

  @keyframes flow-along {
    from { stroke-dashoffset: 24; }
    to   { stroke-dashoffset:  0; }
  }

  .flow {
    animation: flow-along linear infinite;
    pointer-events: none;
    transition: opacity 0.25s ease;
  }

  .rail {
    transition: opacity 0.22s ease;
    pointer-events: none;
  }

  .edge-label {
    font-family: var(--font-mono);
    font-size: 11px;
    fill: var(--clr-muted);
    letter-spacing: 0.03em;
    pointer-events: none;
    transition: opacity 0.22s ease;
  }

  .node { cursor: default; }

  .node-rect {
    transition: stroke 0.18s ease, stroke-width 0.18s ease, opacity 0.18s ease;
  }

  .node-label {
    font-family: var(--font-display);
    font-size: 14px;
    font-weight: 700;
    fill: var(--clr-ink);
    pointer-events: none;
    transition: opacity 0.18s ease;
  }

  .node-sub {
    font-family: var(--font-mono);
    font-size: 10px;
    fill: var(--clr-muted);
    letter-spacing: 0.02em;
    pointer-events: none;
    transition: opacity 0.18s ease;
  }

  .arch__desc {
    display: flex;
    align-items: baseline;
    gap: 1.25rem;
    border-top: 1px solid var(--clr-rule);
    padding: 0.85rem 0.25rem;
    min-height: 4rem;
    opacity: 0;
    transition: opacity 0.2s ease;
  }

  .arch__desc.visible { opacity: 1; }

  .arch__name {
    color: var(--clr-muted);
    white-space: nowrap;
    flex-shrink: 0;
  }

  .arch__body {
    font-size: 0.875rem;
    color: var(--clr-ink-2);
    line-height: 1.65;
  }

  @media (prefers-reduced-motion: reduce) {
    .flow { animation: none; }
  }

  @media (max-width: 580px) {
    .arch { padding: 1rem 1rem 0; }
    .arch__desc { flex-direction: column; gap: 0.3rem; }
  }
</style>
