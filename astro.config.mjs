import { defineConfig } from 'astro/config';
import svelte from '@astrojs/svelte';
import sitemap from '@astrojs/sitemap';

export default defineConfig({
  site: 'https://www.anandu.ca',
  integrations: [
    svelte(),
    // /fr/* is unlinked from the nav while the French translation is under
    // review — excluded here too, so search engines can't discover and
    // index it via the sitemap in the meantime. Remove this filter when
    // the FR nav toggle comes back.
    sitemap({ filter: (page) => !page.includes('/fr/') }),
  ],
  output: 'static',
  prefetch: {
    prefetchAll: false,
    defaultStrategy: 'hover',
  },
  compressHTML: true,
  build: {
    inlineStylesheets: 'auto',
  },
});
