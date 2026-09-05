<script lang="ts">
  import { fade } from 'svelte/transition';
  import defaultProjects from '../../data/projects.json';

  type Project = {
    title: string;
    year: string;
    category: string;
    status: string;
    /** Language-independent — status is a display string (English or French)
     *  and shouldn't be string-matched to decide behaviour. */
    shipped?: boolean;
    description: string;
    stack: string[];
    link?: string;
    repo?: string;
    install?: string;
  };

  interface Props {
    projects?: typeof defaultProjects;
    lang?: 'en' | 'fr';
  }

  const { projects = defaultProjects, lang = 'en' }: Props = $props();

  const visible = projects as Project[];
  const nothingHereYet = lang === 'fr' ? 'Rien à voir pour le moment' : 'Nothing to see here for now';
</script>

<div class="projects">
  <div class="projects__list" role="list">
    {#each visible as project (project.title)}
      <div
        class="project"
        role="listitem"
        in:fade={{ duration: 200 }}
        out:fade={{ duration: 120 }}
      >
        <div class="project__meta">
          <span class="label project__year">{project.year}</span>
          <span class="label project__status project__status--{project.status.toLowerCase().replace(' ', '-')}">{project.status}</span>
        </div>
        <div class="project__body">
          <h3 class="project__title">{project.title}</h3>
          <p class="project__desc">{project.description}</p>
          <ul class="project__stack" aria-label="Stack">
            {#each project.stack as tech}
              <li class="label project__tech">{tech}</li>
            {/each}
          </ul>
          {#if project.shipped && project.install}
            <p class="project__install">
              <span class="label">Install</span>
              <code>{project.install}</code>
            </p>
          {/if}
        </div>
        <div class="project__links">
          {#if project.link}
            <a href={project.link} target="_blank" rel="noopener" class="project__link label">↗ Live</a>
          {/if}
          {#if project.shipped}
            {#if project.repo}
              <a href={project.repo} target="_blank" rel="noopener" class="project__link label">↗ Code</a>
            {/if}
          {:else}
            <span class="project__placeholder label">{nothingHereYet}</span>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .projects__list {
    display: flex;
    flex-direction: column;
  }

  .project {
    display: grid;
    grid-template-columns: 10rem 1fr auto;
    gap: 1.5rem;
    align-items: start;
    padding-block: 1.75rem;
    padding-inline: 0.75rem;
    border-top: 1px solid var(--clr-rule);
    transition: background var(--transition), box-shadow var(--transition);
  }

  .project:last-child { border-bottom: 1px solid var(--clr-rule); }

  .project:hover {
    background: var(--clr-surface);
    box-shadow: inset 2px 0 0 var(--clr-ink);
  }

  .project__meta {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding-top: 0.2rem;
  }

  .project__year                { color: var(--clr-ink); }
  .project__status              { color: var(--clr-muted); }
  .project__status--in-progress { color: var(--clr-ink-2); font-style: italic; }
  .project__status--planning    { color: var(--clr-muted); font-style: italic; }

  .project__body { display: flex; flex-direction: column; gap: 0.5rem; }

  .project__title {
    font-family: var(--font-display);
    font-size: clamp(1.1rem, 2vw, 1.4rem);
    font-weight: 700;
    color: var(--clr-ink);
  }

  .project__desc {
    font-size: 1rem;
    color: var(--clr-muted);
    line-height: 1.65;
  }

  .project__stack {
    display: flex;
    flex-wrap: wrap;
    gap: 0.2rem 1rem;
    margin-top: 0.25rem;
  }

  .project__tech { color: var(--clr-muted); }

  .project__install {
    display: flex;
    align-items: baseline;
    gap: 0.6rem;
    margin-top: 0.4rem;
  }

  .project__install .label { color: var(--clr-muted); flex-shrink: 0; }

  .project__install code {
    font-family: var(--font-mono);
    font-size: 0.82rem;
    color: var(--clr-ink-2);
    overflow-wrap: anywhere;
  }

  .project__links {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    align-items: flex-end;
    padding-top: 0.2rem;
    white-space: nowrap;
  }

  .project__link {
    color: var(--clr-muted);
    transition: color var(--transition);
  }

  .project__link:hover { color: var(--clr-ink); }

  .project__placeholder {
    color: var(--clr-muted);
    font-style: italic;
    white-space: normal;
    text-align: right;
    max-width: 9rem;
  }

  @media (max-width: 640px) {
    .project {
      grid-template-columns: 1fr;
      gap: 0.75rem;
    }
    .project__links { flex-direction: row; align-items: flex-start; }
  }
</style>
