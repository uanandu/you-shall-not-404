# Portfolio

Personal frontend portfolio website built with Astro and Svelte. Fast, statically generated, and deployable for free.

---

## For Users

The site is a single-page experience divided into five sections:

**Hero** — A brief introduction with a call-to-action to explore work or make contact. An animated gradient background and scroll indicator set the tone.

**About** — A short bio and a grid of technologies. Hover over any skill chip to highlight it. The section scrolls into view with a fade-up animation.

**Projects** — A filterable card grid. Click any tag tab (e.g. "Svelte", "AI", "TypeScript") to instantly filter projects by category. Cards animate in and out smoothly. Each card shows a title, description, tags, and links to the live site and/or source code.

**Experience** — A vertical timeline of past roles. Each entry shows company, period, and a short description of work done.

**Contact** — A contact form with inline validation. Fields are checked before submission and error messages appear underneath the relevant input. On success, the form is replaced with a confirmation message.

The navigation bar is fixed at the top, blurs and compresses on scroll, and collapses to a hamburger menu on mobile. All animations respect the OS-level "Reduce Motion" preference.

---

## Technical Details

### Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Site framework | Astro | 6.x |
| UI components | Svelte | 5.x |
| Language | TypeScript | via Astro strict tsconfig |
| Styling | Vanilla CSS (custom properties) | — |
| Package manager | pnpm | 10.x |
| Output | Static (`dist/`) | — |

### Architecture — Astro Islands

Astro generates fully static HTML at build time. JavaScript is only shipped for components that explicitly need it. Each Svelte component uses an `client:*` directive to control when it hydrates:

```
index.astro (static shell)
├── Nav.svelte          client:load     → hydrates immediately (needed for scroll events)
├── ProjectFilter.svelte client:visible → hydrates when scrolled into view
└── ContactForm.svelte  client:visible → hydrates when scrolled into view
```

Components without a `client:*` directive render to static HTML only — zero JS sent to the browser for them.

### Project Structure

```
portfolio/
├── public/
│   └── favicon.svg
├── src/
│   ├── components/
│   │   └── svelte/
│   │       ├── Nav.svelte            Sticky navigation with mobile menu
│   │       ├── ProjectFilter.svelte  Filterable project card grid
│   │       └── ContactForm.svelte    Validated contact form
│   ├── layouts/
│   │   └── Base.astro               HTML shell, font loading, scroll-reveal observer
│   ├── pages/
│   │   └── index.astro              Single page — all sections and static content
│   └── styles/
│       └── global.css               Custom properties, reset, typography, utilities
├── astro.config.mjs
├── tsconfig.json
├── package.json
├── pnpm-lock.yaml
└── CLAUDE.md                        AI assistant context and project conventions
```

### Data Flow

All content is hardcoded in the source files — no CMS, no database, no API calls at runtime.

```
Build time:
  index.astro  →  reads static arrays (skills, experience)
                →  renders HTML for About and Experience sections
                →  passes ProjectFilter and ContactForm as Svelte islands

Runtime (browser):
  Nav.svelte           listens to window scroll → toggles scrolled/open state
  ProjectFilter.svelte filters projects array in $derived state → Svelte renders diff
  ContactForm.svelte   validates fields in $state → calls fetch/mailto on submit
  Base.astro <script>  IntersectionObserver fires .visible class on .reveal elements
```

### Styling System

All design tokens are CSS custom properties on `:root`:

- `--clr-*` — colour palette (background, surface, accent, text, muted)
- `--space-*` — spacing scale (xs → xl)
- `--radius-*` — border radius scale
- `--transition` / `--transition-slow` — easing curves
- `--font-sans` / `--font-mono` — font stacks

Responsive layout uses `clamp()` for fluid type and `min()` for container widths. No breakpoint utility classes — layout shifts are handled per-component with `@media` blocks.

### Animation

| Mechanism | Used for |
|-----------|---------|
| CSS `@keyframes` | Hero scroll indicator pulse |
| CSS `transition` on `.reveal` + JS `IntersectionObserver` | Scroll-triggered fade-up on every section |
| Svelte `in:fly` / `out:fade` | Project cards entering/leaving on filter change |
| Svelte `animate:flip` | Project cards reordering smoothly |
| CSS `transition` on hover | Nav links, buttons, skill chips, social icons |

All animations are disabled when `prefers-reduced-motion: reduce` is set.

### Pages & Routes

The site has one route:

| Route | File | Description |
|-------|------|-------------|
| `/` | `src/pages/index.astro` | Full portfolio — all sections on one page |

Anchor links (`#about`, `#projects`, etc.) handle in-page navigation via `scroll-behavior: smooth`.

### Build & Deploy

```bash
pnpm dev        # Dev server at localhost:4321 (or next available port)
pnpm build      # Static output → dist/
pnpm preview    # Serve dist/ locally to verify production build
```

The `dist/` folder is a self-contained set of static files. Deploy by pointing any static host (Cloudflare Pages, Vercel, Netlify, GitHub Pages) at it. No server-side runtime required.
