import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import {
  articlesForCategory,
  categoryRouteFor,
  knowledgeArticles,
  knowledgeCategories,
  knowledgeCategoryRoutes,
  knowledgeHubPath,
  knowledgeOrigin,
  relatedArticlesFor,
} from "./knowledge-base-data.mjs";

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const siteRoot = path.resolve(scriptDir, "..");
const publicRoot = path.join(siteRoot, "public");
const knowledgeRoot = path.join(publicRoot, "knowledge");

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function escapeAttr(value) {
  return escapeHtml(value).replaceAll("\n", " ");
}

function canonical(pathname) {
  return `${knowledgeOrigin}${pathname}`;
}

function wordsForSearch(article) {
  return [
    article.title,
    article.category,
    article.cluster,
    article.excerpt,
    ...article.tags,
    ...article.keywords,
  ]
    .join(" ")
    .toLowerCase()
    .replaceAll("ё", "е")
    .match(/[\p{L}\p{N}]+/gu)
    ?.join(" ") || "";
}

function buildBaseJsonLd({ path: pathname, title, description }) {
  const url = canonical(pathname);
  const organizationId = `${knowledgeOrigin}/#organization`;
  const websiteId = `${knowledgeOrigin}/#website`;

  return [
    {
      "@type": "Organization",
      "@id": organizationId,
      name: "NERIVA",
      url: `${knowledgeOrigin}/`,
      logo: `${knowledgeOrigin}/assets/brand-logo-mini.png`,
      sameAs: [
        "https://www.youtube.com/@neriva_app",
        "https://www.instagram.com/neriva.ru",
        "https://tiktok.com/@nerivaru",
        "https://t.me/NERIVAapp_bot",
      ],
    },
    {
      "@type": "WebSite",
      "@id": websiteId,
      name: "NERIVA",
      url: `${knowledgeOrigin}/`,
      inLanguage: "ru",
      publisher: { "@id": organizationId },
      potentialAction: {
        "@type": "SearchAction",
        target: `${knowledgeOrigin}${knowledgeHubPath}?q={search_term_string}`,
        "query-input": "required name=search_term_string",
      },
    },
    {
      "@type": "WebPage",
      "@id": `${url}#webpage`,
      url,
      name: title,
      description,
      inLanguage: "ru",
      isPartOf: { "@id": websiteId },
    },
  ];
}

function renderHead({ title, description, path: pathname, jsonLd }) {
  const url = canonical(pathname);
  return `<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>${escapeHtml(title)}</title>
    <meta name="description" content="${escapeAttr(description)}" />
    <meta name="robots" content="index, follow" />
    <link rel="canonical" href="${url}" />
    <meta property="og:title" content="${escapeAttr(title)}" />
    <meta property="og:description" content="${escapeAttr(description)}" />
    <meta property="og:type" content="article" />
    <meta property="og:url" content="${url}" />
    <meta property="og:image" content="${knowledgeOrigin}/assets/seo/neriva-social-preview-v2-ru.jpg" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="${escapeAttr(title)}" />
    <meta name="twitter:description" content="${escapeAttr(description)}" />
    <meta name="twitter:image" content="${knowledgeOrigin}/assets/seo/neriva-social-preview-v2-ru.jpg" />
    <meta name="theme-color" content="#0b1422" />
    <link rel="icon" href="/assets/brand-logo-mini.png?v=20260522-react" />
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@500;700&family=Sora:wght@500;600;700;800&display=swap" rel="stylesheet" />
    <link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap" rel="stylesheet" />
    <link rel="stylesheet" href="/assets/knowledge-base.css" />
    <script type="application/ld+json">${JSON.stringify({ "@context": "https://schema.org", "@graph": jsonLd })}</script>
  </head>`;
}

function renderKnowledgeHeader({ active = true } = {}) {
  return `<nav class="knowledge-topnav" aria-label="Основная навигация">
    <div class="knowledge-topnav__inner">
      <a class="knowledge-brand" href="/neriva.html" aria-label="NERIVA">
        <img src="/assets/brand-logo-mini.png" alt="" aria-hidden="true" />
        <span>NERIVA</span>
      </a>
      <div class="knowledge-topnav__links">
        <a href="/neriva.html#features">Возможности</a>
        <a href="/neriva.html#pricing">Тарифы</a>
        <a href="/neriva.html#faq">FAQ</a>
      </div>
      <div class="knowledge-topnav__actions">
        <a class="knowledge-nav__link--knowledge${active ? " is-active" : ""}" href="/knowledge/">База знаний</a>
        <select class="knowledge-language" aria-label="Язык">
          <option value="ru">Русский</option>
          <option value="en">English</option>
        </select>
        <a class="knowledge-app-button" href="/app/">Открыть app</a>
      </div>
    </div>
  </nav>`;
}

function renderKnowledgeFooter() {
  return `<footer class="knowledge-footer">
    <div class="knowledge-footer__inner">
      <div class="knowledge-footer__brand">
        <strong>NERIVA</strong>
        <p>AI-репетитор для уроков, roleplay, голоса, фото-практики и повтора ошибок в браузере и Telegram.</p>
        <span>© 2026 NERIVA. All rights reserved.</span>
      </div>
      <nav aria-label="Документы">
        <a href="/privacy.html">Политика</a>
        <a href="/terms.html">Условия</a>
        <a href="/agreement.html">Соглашение</a>
        <a href="/consent.html">Согласие</a>
      </nav>
    </div>
  </footer>`;
}

function renderCategoryFilters() {
  return knowledgeCategories
    .map((category) => `<button class="knowledge-filter" type="button" data-filter="${category.id}">${escapeHtml(category.label)}</button>`)
    .join("\n          ");
}

function renderArticleCard(article, featured = false) {
  return `<a class="knowledge-card${featured ? " knowledge-card--featured" : ""}" href="${article.path}" data-category="${article.categoryId}" data-search="${escapeAttr(wordsForSearch(article))}">
    <div class="knowledge-card__meta">
      <span><span class="material-symbols-outlined" aria-hidden="true">article</span>${escapeHtml(article.category)}</span>
      <span>${article.readingTime} мин чтения</span>
    </div>
    <h2>${escapeHtml(article.title)}</h2>
    <p>${escapeHtml(article.excerpt)}</p>
    <div class="knowledge-card__tags">
      <span>#${String(article.number).padStart(2, "0")}</span>
      ${article.tags.slice(0, 2).map((tag) => `<span>${escapeHtml(tag)}</span>`).join("")}
    </div>
  </a>`;
}

function renderCategoryRouteCard(route) {
  const count = articlesForCategory(route.id).length;
  return `<a class="knowledge-cluster-card" href="${route.path}">
    <span>${count} материалов</span>
    <h2>${escapeHtml(route.title)}</h2>
    <p>${escapeHtml(route.description)}</p>
  </a>`;
}

function renderCategoryRouteCards() {
  return knowledgeCategoryRoutes.map((route) => renderCategoryRouteCard(route)).join("\n        ");
}

function renderHubJsonLd() {
  return [
    ...buildBaseJsonLd({
      path: knowledgeHubPath,
      title: "База знаний NERIVA",
      description: "Подробные ответы о том, как учить иностранные языки, запоминать слова, говорить увереннее и не бросать занятия.",
    }),
    {
      "@type": "CollectionPage",
      "@id": `${canonical(knowledgeHubPath)}#collection`,
      url: canonical(knowledgeHubPath),
      name: "База знаний NERIVA",
      inLanguage: "ru",
      hasPart: [
        ...knowledgeCategoryRoutes.map((route) => ({
          "@type": "CollectionPage",
          name: route.title,
          url: canonical(route.path),
        })),
        ...knowledgeArticles.map((article) => ({
          "@type": "Article",
          headline: article.title,
          url: canonical(article.path),
        })),
      ],
    },
    {
      "@type": "BreadcrumbList",
      "@id": `${canonical(knowledgeHubPath)}#breadcrumb`,
      itemListElement: [
        { "@type": "ListItem", position: 1, name: "NERIVA", item: `${knowledgeOrigin}/neriva.html` },
        { "@type": "ListItem", position: 2, name: "База знаний", item: canonical(knowledgeHubPath) },
      ],
    },
  ];
}

function renderHubPage() {
  const featured = knowledgeArticles[0];
  const cards = knowledgeArticles.map((article) => renderArticleCard(article, article.slug === featured.slug)).join("\n        ");
  const popular = knowledgeArticles.slice(15, 20);

  return `<!doctype html>
<html class="dark" lang="ru">
  ${renderHead({
    title: "База знаний NERIVA - как учить иностранные языки",
    description: "Подробные ответы о том, как учить иностранные языки, запоминать слова, говорить увереннее и не бросать занятия.",
    path: knowledgeHubPath,
    jsonLd: renderHubJsonLd(),
  })}
  <body class="knowledge-shell" data-page="knowledge-hub">
    <div class="knowledge-bg" aria-hidden="true"></div>
    ${renderKnowledgeHeader()}
    <main class="knowledge-main knowledge-hub">
      <section class="knowledge-hero" aria-labelledby="knowledge-title">
        <div class="knowledge-hero__copy">
          <span class="knowledge-eyebrow">NERIVA learning intelligence</span>
          <h1 id="knowledge-title">База знаний NERIVA</h1>
          <p>Подробные ответы о том, как учить иностранные языки, запоминать слова, говорить увереннее и не бросать занятия.</p>
          <div class="knowledge-search">
            <span class="material-symbols-outlined" aria-hidden="true">search</span>
            <input class="knowledge-search__input" type="search" placeholder="Найти статью, вопрос или тему" autocomplete="off" />
          </div>
          <div class="knowledge-filterbar" aria-label="Фильтры базы знаний">
            ${renderCategoryFilters()}
            <button class="knowledge-reset" type="button">
              <span class="material-symbols-outlined" aria-hidden="true">close</span>
              Сбросить фильтры
            </button>
          </div>
          <div class="knowledge-results-count" aria-live="polite">${knowledgeArticles.length} материалов</div>
        </div>
        <aside class="knowledge-rail" aria-label="Навигация по базе знаний">
          <div>
            <span class="knowledge-rail__label">Маршрут для старта</span>
            <strong>15 статей для первых 30 дней</strong>
            <p>Начните с самостоятельного обучения, ежедневной привычки, плана и типичных ошибок новичков.</p>
          </div>
          <div>
            <span class="knowledge-rail__label">Популярные вопросы</span>
            ${popular.map((article) => `<a href="${article.path}">${escapeHtml(article.title)}</a>`).join("")}
          </div>
          <div class="knowledge-plan">
            <span class="knowledge-rail__label">План публикации</span>
            <strong>500 материалов</strong>
            <p>5-10 статей в неделю: старт, лексика, английский, методика, грамматика, listening, speaking и CEFR.</p>
          </div>
        </aside>
      </section>
      <section class="knowledge-clusters" aria-label="Маршруты базы знаний">
        <div class="knowledge-section-heading">
          <span class="knowledge-eyebrow">SEO/AEO routes</span>
          <h2>Маршруты по темам</h2>
          <p>Каждый маршрут собирает статьи в понятный путь: от первого вопроса до практического действия в NERIVA.</p>
        </div>
        <div class="knowledge-cluster-grid">
          ${renderCategoryRouteCards()}
        </div>
      </section>
      <section class="knowledge-list" aria-label="Статьи базы знаний">
        ${cards}
      </section>
    </main>
    ${renderKnowledgeFooter()}
    <script>${renderHubScript()}</script>
  </body>
</html>`;
}

function renderHubScript() {
  return `
(() => {
  const input = document.querySelector(".knowledge-search__input");
  const cards = Array.from(document.querySelectorAll(".knowledge-card"));
  const filters = Array.from(document.querySelectorAll(".knowledge-filter"));
  const reset = document.querySelector(".knowledge-reset");
  const count = document.querySelector(".knowledge-results-count");
  let activeCategory = "";

  function tokens(value) {
    return (String(value || "").toLowerCase().replaceAll("ё", "е").match(/[\\p{L}\\p{N}]+/gu) || []);
  }

  function applyFilters() {
    const queryTokens = tokens(input?.value || "");
    let visible = 0;
    cards.forEach((card) => {
      const searchText = tokens(card.dataset.search || "").join(" ");
      const queryMatch = queryTokens.length === 0 || queryTokens.every((token) => searchText.includes(token));
      const categoryMatch = !activeCategory || card.dataset.category === activeCategory;
      const show = queryMatch && categoryMatch;
      card.hidden = !show;
      if (show) visible += 1;
    });
    const query = input?.value.trim() || "";
    count.textContent = query ? visible + ' материалов по запросу "' + query + '"' : visible + " материалов";
  }

  input?.addEventListener("input", applyFilters);
  filters.forEach((button) => {
    button.addEventListener("click", () => {
      activeCategory = button.dataset.filter || "";
      if (input) input.value = "";
      filters.forEach((item) => item.classList.toggle("is-active", item === button));
      applyFilters();
    });
  });
  reset?.addEventListener("click", () => {
    activeCategory = "";
    if (input) input.value = "";
    filters.forEach((item) => item.classList.remove("is-active"));
    applyFilters();
  });
  const params = new URLSearchParams(window.location.search);
  const initialQuery = params.get("q");
  if (initialQuery && input) {
    input.value = initialQuery;
  }
  applyFilters();
})();`;
}

function renderCategoryRouteJsonLd(route) {
  const articles = articlesForCategory(route.id);
  const routeUrl = canonical(route.path);
  return [
    ...buildBaseJsonLd({
      path: route.path,
      title: `${route.title} - База знаний NERIVA`,
      description: route.description,
    }),
    {
      "@type": "CollectionPage",
      "@id": `${routeUrl}#collection`,
      url: routeUrl,
      name: route.title,
      description: route.description,
      inLanguage: "ru",
      isPartOf: { "@id": `${knowledgeOrigin}/#website` },
      hasPart: articles.map((article) => ({
        "@type": "Article",
        headline: article.title,
        url: canonical(article.path),
      })),
    },
    {
      "@type": "ItemList",
      "@id": `${routeUrl}#items`,
      name: `${route.title}: статьи`,
      itemListElement: articles.map((article, index) => ({
        "@type": "ListItem",
        position: index + 1,
        name: article.title,
        url: canonical(article.path),
      })),
    },
    {
      "@type": "BreadcrumbList",
      "@id": `${routeUrl}#breadcrumb`,
      itemListElement: [
        { "@type": "ListItem", position: 1, name: "NERIVA", item: `${knowledgeOrigin}/neriva.html` },
        { "@type": "ListItem", position: 2, name: "База знаний", item: canonical(knowledgeHubPath) },
        { "@type": "ListItem", position: 3, name: route.title, item: routeUrl },
      ],
    },
  ];
}

function renderCategoryRouteSteps(route) {
  return route.route
    .map(([title, body], index) => `<div class="cluster-route__step">
      <span>${index + 1}</span>
      <div>
        <h3>${escapeHtml(title)}</h3>
        <p>${escapeHtml(body)}</p>
      </div>
    </div>`)
    .join("\n            ");
}

function renderCategoryRoutePage(route) {
  const articles = articlesForCategory(route.id);
  const cards = articles.map((article, index) => renderArticleCard(article, index === 0)).join("\n        ");

  return `<!doctype html>
<html class="dark" lang="ru">
  ${renderHead({
    title: `${route.title} - База знаний NERIVA`,
    description: route.description,
    path: route.path,
    jsonLd: renderCategoryRouteJsonLd(route),
  })}
  <body class="knowledge-shell" data-page="knowledge-cluster">
    <div class="knowledge-bg" aria-hidden="true"></div>
    ${renderKnowledgeHeader()}
    ${renderMobileArticleHeader()}
    <main class="knowledge-main knowledge-cluster">
      <section class="knowledge-hero knowledge-cluster-hero" aria-labelledby="cluster-title">
        <div class="knowledge-hero__copy">
          <nav aria-label="Breadcrumb" class="article-breadcrumb">
            <ol>
              <li><a href="/neriva.html">Главная</a></li>
              <li><span>/</span></li>
              <li><a href="/knowledge/">База знаний</a></li>
              <li><span>/</span></li>
              <li aria-current="page">${escapeHtml(route.title)}</li>
            </ol>
          </nav>
          <span class="knowledge-eyebrow">NERIVA topic route</span>
          <h1 id="cluster-title">${escapeHtml(route.title)}</h1>
          <p>${escapeHtml(route.description)}</p>
        </div>
        <aside class="cluster-cta" aria-label="Практика в NERIVA">
          <span class="knowledge-rail__label">Практика</span>
          <strong>${escapeHtml(route.cta.title)}</strong>
          <p>${escapeHtml(route.cta.text)}</p>
          <a href="/app/">${escapeHtml(route.cta.button)} <span class="material-symbols-outlined" aria-hidden="true">arrow_forward</span></a>
        </aside>
      </section>
      <section class="cluster-route" aria-label="Маршрут обучения">
        ${renderCategoryRouteSteps(route)}
      </section>
      <section class="knowledge-list" aria-label="Статьи маршрута">
        ${cards}
      </section>
    </main>
    ${renderKnowledgeFooter()}
  </body>
</html>`;
}

function renderBreadcrumb(article) {
  return `<nav aria-label="Breadcrumb" class="article-breadcrumb">
    <ol>
      <li><a href="/neriva.html">Главная</a></li>
      <li><span>/</span></li>
      <li><a href="/knowledge/">База знаний</a></li>
      <li><span>/</span></li>
      <li aria-current="page">${escapeHtml(article.category)}</li>
    </ol>
  </nav>`;
}

function renderArticleSections(article) {
  const why = article.why.map((paragraph) => `<p>${escapeHtml(paragraph)}</p>`).join("\n            ");
  const plan = article.plan
    .map(([title, body], index) => `<div class="article-step">
      <span>${index + 1}</span>
      <div>
        <h3>${escapeHtml(title)}</h3>
        <p>${escapeHtml(body)}</p>
      </div>
    </div>`)
    .join("\n            ");
  const mistakes = article.mistakes
    .map(([title, body]) => `<li><span class="material-symbols-outlined" aria-hidden="true">error</span><span><strong>${escapeHtml(title)}:</strong> ${escapeHtml(body)}</span></li>`)
    .join("\n              ");

  return `<section class="article-section reveal-element" id="why-it-works">
            <h2>Почему это работает именно так</h2>
            ${why}
          </section>
          <section class="article-section reveal-element" id="practical-plan">
            <h2>Практический план</h2>
            <div class="article-steps">${plan}</div>
          </section>
          <section class="article-section reveal-element" id="common-mistakes">
            <h2>Типичные ошибки</h2>
            <ul class="article-mistakes">${mistakes}</ul>
          </section>
          <section class="article-section reveal-element" id="how-neriva-helps">
            <h2>Как NERIVA может помочь</h2>
            <p>${escapeHtml(article.neriva)}</p>
            <div class="article-platform-box">
              <h3>Ключевой учебный цикл</h3>
              <ul>
                <li>короткое объяснение и пример;</li>
                <li>самостоятельный ответ текстом или голосом;</li>
                <li>исправление ошибки и сохранение слабого места;</li>
                <li>возврат материала в повторение.</li>
              </ul>
            </div>
          </section>`;
}

function renderArticleFaq(article) {
  return article.faq
    .map(([question, answer], index) => `<details class="article-faq__item"${index === 0 ? " open" : ""}>
      <summary>${escapeHtml(question)}<span class="material-symbols-outlined" aria-hidden="true">expand_more</span></summary>
      <p>${escapeHtml(answer)}</p>
    </details>`)
    .join("\n              ");
}

function renderRelatedDesktop(article) {
  return relatedArticlesFor(article)
    .map((related) => `<a class="article-related__card" href="${related.path}">
      <span>${escapeHtml(related.category)}</span>
      <h3>${escapeHtml(related.title)}</h3>
      <p>${escapeHtml(related.excerpt)}</p>
    </a>`)
    .join("\n          ");
}

function renderRelatedMobile(article) {
  return relatedArticlesFor(article)
    .map((related) => `<a class="article-related__row" href="${related.path}">
      <span class="material-symbols-outlined" aria-hidden="true">article</span>
      <div>
        <h4>${escapeHtml(related.title)}</h4>
        <small>${related.readingTime} мин чтения</small>
      </div>
    </a>`)
    .join("\n          ");
}

function renderArticleJsonLd(article) {
  const articleUrl = canonical(article.path);
  return [
    ...buildBaseJsonLd({
      path: article.path,
      title: `${article.title} - База знаний NERIVA`,
      description: article.excerpt,
    }),
    {
      "@type": "Article",
      "@id": `${articleUrl}#article`,
      headline: article.title,
      description: article.excerpt,
      url: articleUrl,
      inLanguage: "ru",
      datePublished: article.updatedIso,
      dateModified: article.updatedIso,
      author: { "@type": "Organization", name: "NERIVA" },
      publisher: { "@id": `${knowledgeOrigin}/#organization` },
      articleSection: article.category,
      keywords: [...article.tags, ...article.keywords].join(", "),
    },
    {
      "@type": "BreadcrumbList",
      "@id": `${articleUrl}#breadcrumb`,
      itemListElement: [
        { "@type": "ListItem", position: 1, name: "NERIVA", item: `${knowledgeOrigin}/neriva.html` },
        { "@type": "ListItem", position: 2, name: "База знаний", item: canonical(knowledgeHubPath) },
        { "@type": "ListItem", position: 3, name: article.category, item: articleUrl },
      ],
    },
    {
      "@type": "FAQPage",
      "@id": `${articleUrl}#faq`,
      url: articleUrl,
      inLanguage: "ru",
      mainEntity: article.faq.map(([question, answer]) => ({
        "@type": "Question",
        name: question,
        acceptedAnswer: { "@type": "Answer", text: answer },
      })),
    },
  ];
}

function renderMobileArticleHeader(article) {
  return `<header class="knowledge-mobile-header">
    <div class="knowledge-mobile-header__bar">
      <a href="/knowledge/" aria-label="Назад к базе знаний"><span class="material-symbols-outlined" aria-hidden="true">arrow_back</span></a>
      <strong>NERIVA</strong>
      <a href="/app/" aria-label="Открыть NERIVA"><span class="material-symbols-outlined" aria-hidden="true">bookmark_border</span></a>
    </div>
    <form class="knowledge-mobile-header__search" action="/knowledge/" role="search">
      <span class="material-symbols-outlined" aria-hidden="true">search</span>
      <input class="knowledge-search__input" name="q" type="search" placeholder="Найти статью..." />
    </form>
  </header>`;
}

function renderArticleRouteLinks(article) {
  const route = categoryRouteFor(article.categoryId);
  const relatedLinks = relatedArticlesFor(article);

  return `<section class="article-route-links reveal-element" aria-label="Продолжить маршрут">
        <div class="article-route-links__copy">
          <span class="knowledge-eyebrow">Маршрут по теме</span>
          <h2>Что изучить дальше</h2>
          <p>Чтобы материал не остался отдельной статьей, двигайтесь по соседним вопросам и возвращайтесь к практике.</p>
        </div>
        <div class="article-route-links__grid">
          <a class="article-route-links__pillar" href="${route.path}">
            <span class="material-symbols-outlined" aria-hidden="true">route</span>
            <strong>${escapeHtml(route.title)}</strong>
            <small>Открыть весь маршрут</small>
          </a>
          ${relatedLinks
            .map((related) => `<a href="${related.path}">
              <span>${escapeHtml(related.category)}</span>
              <strong>${escapeHtml(related.title)}</strong>
            </a>`)
            .join("\n          ")}
        </div>
      </section>`;
}

function renderArticleCta(article) {
  const route = categoryRouteFor(article.categoryId);
  return `<section class="article-cta article-cta--${article.categoryId} reveal-element">
            <h2>${escapeHtml(route.cta.title)}</h2>
            <p>${escapeHtml(route.cta.text)}</p>
            <a href="/app/">${escapeHtml(route.cta.button)} <span class="material-symbols-outlined" aria-hidden="true">arrow_forward</span></a>
          </section>`;
}

function renderArticlePage(article) {
  const tocLinks = [
    ["why-it-works", "Почему это работает"],
    ["practical-plan", "Практический план"],
    ["common-mistakes", "Типичные ошибки"],
    ["how-neriva-helps", "Как NERIVA может помочь"],
    ["faq", "Частые вопросы"],
  ];

  return `<!doctype html>
<html class="dark" lang="ru">
  ${renderHead({
    title: `${article.title} - База знаний NERIVA`,
    description: article.excerpt,
    path: article.path,
    jsonLd: renderArticleJsonLd(article),
  })}
  <body class="knowledge-shell knowledge-article-shell" data-page="knowledge-article">
    <div class="knowledge-bg" aria-hidden="true"></div>
    ${renderKnowledgeHeader()}
    ${renderMobileArticleHeader(article)}
    <main class="knowledge-main knowledge-article">
      <header class="article-header reveal-element">
        ${renderBreadcrumb(article)}
        <div class="article-kicker">
          <span>${escapeHtml(article.category)}</span>
          <span>#${String(article.number).padStart(2, "0")}</span>
        </div>
        <h1>${escapeHtml(article.title)}</h1>
        <div class="article-meta">
          <span><span class="material-symbols-outlined" aria-hidden="true">schedule</span>${article.readingTime} мин чтения</span>
          <span></span>
          <span><span class="material-symbols-outlined" aria-hidden="true">category</span>${escapeHtml(article.category)}</span>
          <span></span>
          <span><span class="material-symbols-outlined" aria-hidden="true">update</span>Обновлено ${article.updated}</span>
        </div>
      </header>
      <section class="article-answer reveal-element">
        <div>
          <span class="material-symbols-outlined" aria-hidden="true">lightbulb</span>
          <strong>Короткий ответ</strong>
        </div>
        <p>${escapeHtml(article.directAnswer)}</p>
      </section>
      ${renderArticleRouteLinks(article)}
      <details class="article-toc-mobile">
        <summary>Содержание <span class="material-symbols-outlined" aria-hidden="true">expand_more</span></summary>
        <nav>
          ${tocLinks.map(([id, label]) => `<a href="#${id}">${escapeHtml(label)}</a>`).join("")}
        </nav>
      </details>
      <div class="article-layout">
        <aside class="article-toc">
          <div>
            <h3>Содержание</h3>
            <nav>
              ${tocLinks.map(([id, label]) => `<a class="toc-link" href="#${id}">${escapeHtml(label)}</a>`).join("")}
            </nav>
          </div>
        </aside>
        <article class="article-body">
          ${renderArticleSections(article)}
          <section class="article-section reveal-element" id="faq">
            <h2>Частые вопросы</h2>
            <div class="article-faq">${renderArticleFaq(article)}</div>
          </section>
          ${renderArticleCta(article)}
        </article>
      </div>
      <section class="article-related reveal-element">
        <h2>Связанные материалы</h2>
        <div class="article-related__grid">${renderRelatedDesktop(article)}</div>
        <div class="article-related__rows">${renderRelatedMobile(article)}</div>
        <a class="article-related__all" href="/knowledge/">Все статьи <span class="material-symbols-outlined" aria-hidden="true">arrow_forward</span></a>
      </section>
    </main>
    ${renderKnowledgeFooter()}
    <script>${renderArticleScript()}</script>
  </body>
</html>`;
}

function renderArticleScript() {
  return `
(() => {
  const revealElements = Array.from(document.querySelectorAll(".reveal-element"));
  const observer = "IntersectionObserver" in window
    ? new IntersectionObserver((entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add("visible");
            observer.unobserve(entry.target);
          }
        });
      }, { threshold: 0.1 })
    : null;
  revealElements.forEach((element) => observer ? observer.observe(element) : element.classList.add("visible"));

  const sections = Array.from(document.querySelectorAll(".article-section[id]"));
  const links = Array.from(document.querySelectorAll(".toc-link"));
  const activate = () => {
    let current = "";
    sections.forEach((section) => {
      if (window.scrollY >= section.offsetTop - 220) current = section.id;
    });
    links.forEach((link) => link.classList.toggle("active", current && link.getAttribute("href") === "#" + current));
  };
  window.addEventListener("scroll", activate, { passive: true });
  activate();
})();`;
}

await fs.mkdir(knowledgeRoot, { recursive: true });
await fs.writeFile(path.join(knowledgeRoot, "index.html"), renderHubPage(), "utf8");

for (const route of knowledgeCategoryRoutes) {
  await fs.writeFile(path.join(knowledgeRoot, path.basename(route.path)), renderCategoryRoutePage(route), "utf8");
}

for (const article of knowledgeArticles) {
  await fs.writeFile(path.join(knowledgeRoot, `${article.slug}.html`), renderArticlePage(article), "utf8");
}

console.log(`Generated ${knowledgeArticles.length} Knowledge Base articles, ${knowledgeCategoryRoutes.length} category routes, and /knowledge/index.html`);
