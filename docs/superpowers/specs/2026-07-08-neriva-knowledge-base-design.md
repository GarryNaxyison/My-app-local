# NERIVA Knowledge Base Design

## Context

The user approved a separate NERIVA Knowledge Base and provided three Stitch-generated reference files:

- `C:/Temp/neriva_knowledge_base_search_results.html`
- `C:/Temp/neriva_knowledge_base_article_detail.html`
- `C:/Temp/neriva_knowledge_base_article_detail_mobile.html`

The implementation must follow those layouts 1:1 in visual language: fixed glass header, dark navy canvas, Sora/Inter/JetBrains Mono typography, pale blue primary text/accent, rounded glass search/cards, dense article rows, sticky article table of contents, mobile article header with search, and compact related rows.

One Stitch artifact is intentionally corrected: the mobile file renders a white body because Tailwind `theme()` declarations in plain CSS are not processed. The production version keeps the same mobile composition but uses the dark NERIVA token system from the desktop references.

## Scope

Build a static `/knowledge/` section inside `site-react`:

- `/knowledge/index.html` as the Knowledge Base hub.
- `/knowledge/<slug>.html` for the first 25 Russian articles from the provided DOCX topic list.
- Client-side search and category filtering on the hub.
- JSON-LD for `Article`, `FAQPage`, `BreadcrumbList`, `Organization`, and `WebSite`.
- `sitemap.xml` entries for the hub and all articles.
- A public-site header link labelled `База знаний` / `Knowledge Base`, placed immediately to the left of the language selector.

Out of scope:

- 35-language Knowledge Base translation.
- Replacing the existing 9 RU/EN product SEO/AEO pages.
- Cloudflare dashboard or DNS changes.

## Content Model

The first 25 articles are the first 25 topics from `md style.docx`:

1. Можно ли выучить иностранный язык самостоятельно?
2. Как начать учить язык с нуля?
3. Сколько времени нужно для изучения языка?
4. Какой иностранный язык проще всего выучить?
5. Как выбрать язык для изучения?
6. Какой язык лучше учить первым?
7. Можно ли выучить язык после 30 лет?
8. Можно ли выучить язык после 40 лет?
9. Почему взрослым сложно учить языки?
10. Как заниматься языком каждый день?
11. Сколько минут в день нужно изучать язык?
12. Как создать привычку изучать язык?
13. Как не бросить изучение языка?
14. Какие ошибки делают новички?
15. Как составить план обучения языку?
16. Почему иностранные слова быстро забываются?
17. Как запоминать слова быстрее?
18. Как увеличить словарный запас?
19. Сколько слов нужно знать для общения?
20. Сколько слов нужно знать для уровня B1?
21. Как выучить первые 1000 слов?
22. Почему зубрежка слов не работает?
23. Как учить слова без зубрежки?
24. Почему нужно учить фразы вместо слов?
25. Что такое интервальное повторение?

Each article includes:

- Short direct answer.
- Why it works.
- Practical plan.
- Common mistakes.
- How NERIVA can help.
- Five FAQ items.
- Three related articles.

## Visual Contract

Use a local CSS file instead of CDN Tailwind. The CSS must reproduce the Stitch references:

- Background: `#0b1422`.
- Lowest surface: `#060e1c`.
- Surface container: `#17202e`.
- Low surface: `#131c2a`.
- Glass background: `rgba(8, 17, 31, 0.7)`.
- Stroke: `rgba(255, 255, 255, 0.12)`.
- Text: `#dae3f7`.
- Muted text: `#c4c5d5`.
- Primary: `#b7c4ff`.
- Electric focus: `#3d7bff`.
- Error: `#ffb4ab`.
- Desktop header height: `80px`.
- Max content width: `1280px`.
- Hub hero top padding: `120px`.
- Desktop article top padding: `128px`.
- Desktop article content split: 256px sticky TOC + max 768px article body.
- Mobile article header: sticky, brand centered, back/bookmark icons, search below.
- Mobile touch targets: at least 44px.

## Behavior

The hub search filters by title, category, cluster, tags, excerpt, and keywords. Typing `слова` must narrow the list to the vocabulary articles and show a result count. Category filters remain clickable and a reset button restores all 25 articles.

Article FAQ accordions open and close without external dependencies. Article TOC links navigate to local sections. The desktop TOC is sticky; mobile uses a collapsible details block.

## Verification

Required checks:

- `npm --prefix site-react run e2e -- e2e/knowledge-base.spec.ts`
- `npm --prefix site-react run build`
- Focused existing SEO/deploy tests after build.
- Visual screenshots of `/knowledge/`, a desktop article, and mobile article through local Playwright.
- Production verification after deploy: `/knowledge/`, first article, `/sitemap.xml`, and public landing header.
