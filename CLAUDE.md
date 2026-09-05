# Portfolio — CLAUDE.md

## Project Overview

A personal frontend portfolio website built with **AstroJS + Svelte**. The aesthetic is editorial brutalism — think print magazine meets personal web, not tech startup. Dark-only, typographically driven, high contrast white-on-black. No gradients, no glow, no light surfaces, no light/dark toggle.

Reference: `public/landing.png` — shows the *original* light/cream version of this design (superseded 2026-09-05 — see Colour Palette below). Still useful for typography, layout, and composition; ignore its colours.

## Tech Stack

| Tool | Purpose |
|------|---------|
| **Astro 6.x** | Static site generator / routing / islands architecture |
| **Svelte 5** | Interactive UI components (via Astro islands) |
| **pnpm** | Package manager — never use npm or yarn |
| **TypeScript** | Type safety across `.astro` and `.svelte` files |
| **CSS (vanilla + custom properties)** | Styling — no Tailwind, no CSS-in-JS |

## Package Manager Rules

- **Always use `pnpm`** — never `npm install`, `npm run`, `yarn`, etc.
- Install deps: `pnpm add <pkg>` / `pnpm add -D <pkg>`
- Run scripts: `pnpm dev`, `pnpm build`, `pnpm preview`

## Project Structure

```
portfolio/
├── public/
│   ├── landing.png          Structure/layout reference only — colours superseded, see Aesthetic section
│   └── favicon.svg
├── src/
│   ├── components/
│   │   └── svelte/          Interactive islands only
│   ├── layouts/
│   ├── pages/
│   └── styles/
├── astro.config.mjs
└── tsconfig.json
```

---

## Aesthetic — Editorial Brutalism (Dark)

The visual *structure* — typography, layout, composition — still comes from `public/landing.png`; only the colour direction inverted, on 2026-09-05, to a dark-only palette. Every design decision should be traceable back to the reference for structure, and to the palette below for colour.

### What the design shows

- **Pitch-black background** — literal, not a tinted near-black. No light surfaces anywhere.
- **Near-white text** — not pure `#fff` (avoids halation/eye strain on a pure-black field). Think ink reversed onto black paper.
- **Large serif display type** — dominant, editorial, takes up space unapologetically. The name is the hero.
- **Small, light-weight sans-serif** for nav and supporting text — minimal footprint.
- **Vintage/analog imagery** — a CRT television, black-and-white photo. Tactile, physical, not digital-clean.
- **Sparse layout** — deliberate negative space. Content doesn't fill the page; it lives in it.
- **No colour accent** — monochromatic (one narrow exception: a single accessible red for error/validation states, `--clr-error` — the only functional, non-decorative colour in the system).
- **Coordinate/timestamp detail** — small contextual text in corners (location, time). Adds character without noise.
- **Minimal navigation** — text links spread across the top, no backgrounds, no borders, no hover glow.

### Colour Palette

```css
--clr-bg:        #000000;   /* pitch black */
--clr-surface:   #0f0f0d;   /* barely-lifted black for hover/active states */
--clr-ink:       #fafaf9;   /* near-white — primary text, ~20:1 on bg */
--clr-ink-2:     #b6b2aa;   /* warm light grey — secondary text, ~9.9:1 on bg */
--clr-muted:     #8f8b83;   /* muted warm grey — captions, labels, ~6.2:1 on bg */
--clr-rule:      #5c584f;   /* thin rule lines, borders — ~3:1 on bg (WCAG non-text minimum) */
--clr-error:     #e5484d;   /* form/terminal error states only — ~5.4:1 on bg */
```

Contrast ratios computed against pure black (`--clr-bg`); all meet or exceed WCAG AA for their role (rules meet the 3:1 non-text minimum, everything else clears 4.5:1 body-text AA, most clear 7:1 AAA). No other colour, no gradients, no rgba glow colours. If you feel the urge to add a colour beyond `--clr-error`, don't.

### Typography

```css
--font-display: 'Playfair Display', 'Georgia', serif;   /* hero names, section titles */
--font-body:    'Inter', system-ui, sans-serif;          /* body, nav, UI */
--font-mono:    'JetBrains Mono', monospace;             /* code tags, metadata */
```

- Display type is large and loose (`letter-spacing: -0.02em`, tight leading).
- Body text is small and airy (`font-size: 0.875–1rem`, generous `line-height: 1.75`).
- Navigation uses all-caps, tracked-out small text (`letter-spacing: 0.12em`, `font-size: 0.75rem`).
- No bold body text. Weight contrast comes from size and typeface, not `font-weight: 700` everywhere.

### Layout Principles

- **Max width ~80rem**, wide margins, content breathes.
- Hero occupies the full viewport: name in large display type, an image (or image placeholder) anchored to a side, supporting text minimal and small.
- Sections separated by thin `1px` rules, not gaps or cards.
- Grid is asymmetric where possible — a 60/40 or 70/30 split feels more editorial than equal columns.
- No rounded corners on structural elements (cards, containers). Sharp edges only. `border-radius: 0` is the default.
- Borders are `1px solid var(--clr-rule)` — never coloured, never thick.

### What to Avoid — Hard Rules

| Banned | Why |
|--------|-----|
| Light backgrounds, cream, off-white surfaces | Wrong aesthetic entirely — dark-only, no light variant, no toggle |
| Pure `#fff` for body text | Halation/eye strain on pure black — use `--clr-ink` |
| Gradients (`linear-gradient`, `background-clip: text`) | Too digital, too decorative |
| Glow / box-shadow with colour | Same — belongs in the old neon theme |
| Glassmorphism / `backdrop-filter: blur` | |
| Rounded corners on containers | Softens the editorial edge |
| Any colour accent besides `--clr-error` | Monochromatic except that one functional exception |
| Heavy bold weights on body text | Use type size for hierarchy instead |
| Gradient text (`-webkit-text-fill-color: transparent`) | |
| Neon or saturated colours anywhere | |

---

## Sections / Pages

1. **Hero** — Full-viewport. Large display-type name. Vintage/analog image element. Small supporting descriptor text. Sparse.
2. **About** — Editorial two-column: long-form text left, skills/facts right as a simple list (no chips, no cards).
3. **Projects** — Text-first list or minimal grid. Filter if needed, but understated. No coloured tags.
4. **Experience** — Clean typographic timeline. Dates in mono. No dots, no coloured lines.
5. **Contact** — Ultra-minimal. Email link prominent, form optional.

## Animation Guidelines

- Animations must be **subtle and mechanical** — no springy, bouncy, or glowing transitions.
- Scroll reveals: simple `opacity` fade, no `translateY`. Slow (`0.8s ease`).
- Hover states: `opacity` shift or thin `border-bottom` underline. No colour changes, no transforms.
- Svelte transitions: `fade` only. No `fly`, no `scale`.
- `prefers-reduced-motion: reduce` — disable all transitions.
- If an animation draws attention to itself, remove it.

## Deployment

Target: **Cloudflare Pages** or **Vercel** (both free tier). Output mode: `static`. No server-side rendering needed.

## Commands

```bash
pnpm dev          # Dev server at localhost:4321
pnpm build        # Production build → dist/
pnpm preview      # Preview production build locally
```

## Code Style

- `.astro` files: frontmatter in TypeScript, template uses semantic HTML5 elements.
- `.svelte` files: `<script lang="ts">`, single-file components.
- CSS: BEM-lite naming (`.section__title`, `.card--featured`). No utility class soup.
- No inline styles unless driven by dynamic JS values.
- No comments unless the why is non-obvious.
