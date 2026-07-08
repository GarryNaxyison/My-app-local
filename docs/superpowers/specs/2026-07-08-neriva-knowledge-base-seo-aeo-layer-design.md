# NERIVA Knowledge Base SEO/AEO Layer Design

## Goal

Strengthen the existing 100-article Knowledge Base for SEO, AEO, and conversion without rewriting articles one by one. The work adds a systematic layer over the generator: topic-specific calls to action, in-body internal links, and category pillar pages.

## Approved Direction

Use the existing static-generation architecture. Keep all article content indexable, fast, and consistent with the current NERIVA landing/Knowledge Base visual language. Do not change the backend or web app routes.

## Scope

- Add one pillar page per Knowledge Base category under `/knowledge/<category>.html`.
- Add those pillar URLs to `sitemap.xml`.
- Add an in-body internal-link section to every article.
- Replace the generic article CTA copy with category-specific CTA copy.
- Keep existing article pages, search, FAQ, Article JSON-LD, BreadcrumbList, and related cards.
- Add tests that prove pillar pages, sitemap entries, article CTA, and in-body internal links exist.

## Content Model

Each category gets a compact route object:

- `id`: category id, such as `vocabulary`.
- `path`: static pillar URL, such as `/knowledge/vocabulary.html`.
- `title`: human-readable pillar title.
- `description`: meta description and page intro.
- `route`: three short route steps for the pillar page.
- `cta`: title, text, and button label for category-specific conversion.

Article CTAs are selected by `categoryId`. All CTA links can keep using `/app/` for now, because deep product routes are not stable enough to publish in SEO pages.

## Page Structure

Pillar pages use the same Knowledge Base shell and header as the hub:

- hero title and intro;
- a route panel with three learning steps;
- topic-specific CTA block;
- article cards for that category;
- JSON-LD `CollectionPage`, `ItemList`, and `BreadcrumbList`.

Article pages receive a new section after the direct-answer block and before the main explanation:

- link to the category pillar page;
- three related articles chosen by the existing related-article logic;
- compact explanatory copy so links feel editorial rather than mechanical.

## Testing

Add Playwright coverage for:

- `/knowledge/vocabulary.html` loads and lists only vocabulary articles;
- the sitemap includes all category pillar URLs;
- an article has a topic-specific CTA, not only the old generic CTA;
- an article has an in-body route/internal-links block before the main content;
- mobile layout has no horizontal overflow after the new blocks.

## Deployment

Use the existing build and site-only deployment flow:

- `npm --prefix site-react run build`;
- archive `Сайт полиглота для бота`;
- upload through deploy-upload;
- run the NERIVA public-site-only deploy script;
- verify production with browser checks and read-only HTTP checks.
