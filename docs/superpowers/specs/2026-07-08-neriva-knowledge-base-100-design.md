# NERIVA Knowledge Base 100 Articles Design

## Goal

Expand the public NERIVA Knowledge Base from 25 to 100 Russian articles while keeping the Stitch-inspired dark visual system, searchable hub, article navigation, FAQ markup, sitemap coverage, and the production deploy flow intact.

## Source Of Truth

The existing first 25 articles stay as-is. The DOCX strategy file lists topics 1-90; therefore articles 26-90 must follow that list exactly. To satisfy the requested "remaining 75 articles", add topics 91-100 from the same strategy scope: exam preparation and AI/technology in language learning.

## Content Standard

Each new article must read like a human educational article, not a keyword page. Every article should have:

- a direct answer that gives the user a usable conclusion immediately;
- explanatory paragraphs with practical linguistic reasoning;
- a concrete plan with several steps;
- typical mistakes;
- a soft NERIVA mention;
- at least 5 FAQ entries through the generated shared FAQ structure;
- unique title, slug, excerpt, category, keywords, and article URL.

The tone is expert, calm, literary Russian: concise where needed, but not dry; practical, not inflated; no placeholder names or Stitch demo copy.

## Architecture

Keep `site-react/scripts/knowledge-base-data.mjs` as the source of truth. Add reusable helpers for extended articles so the 75 new entries do not duplicate the long article object shape manually. The generator continues to emit static `/knowledge/index.html`, `/knowledge/<slug>.html`, and sitemap URLs from `knowledgeArticles`.

Add two categories for the extra 10 topics:

- `exam`: `Экзамены`
- `technology`: `AI и технологии`

Update the hub count from a hardcoded `25 материалов` to the actual `knowledgeArticles.length`.

## Verification

Focused e2e must prove:

- hub renders 100 cards;
- filter count reflects all categories;
- vocabulary filter shows 15 articles;
- CEFR filter shows 8 articles;
- exam and technology article URLs exist;
- the last article `/knowledge/how-to-choose-language-learning-app.html` renders as a full article;
- sitemap includes `/knowledge/how-language-memory-works.html`, `/knowledge/how-to-prepare-for-language-exam.html`, and `/knowledge/how-to-choose-language-learning-app.html`.

Final verification also requires production build, existing SEO tests, encoding artifact check, visual smoke for hub/article/mobile, git push, deploy, and production URL checks.
