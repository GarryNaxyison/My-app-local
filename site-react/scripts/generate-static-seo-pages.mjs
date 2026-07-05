import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { englishOrigin, legalLinks, russianOrigin, seoPages, socialProfileUrls } from "./static-seo-pages-data.mjs";

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const siteRoot = path.resolve(scriptDir, "..");
const publicRoot = path.join(siteRoot, "public");

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function canonicalFor(page, language) {
  const copy = page[language];
  return `${language === "ru" ? russianOrigin : englishOrigin}${copy.path}`;
}

function buildJsonLd(page, language) {
  const copy = page[language];
  const canonicalUrl = canonicalFor(page, language);
  const origin = language === "ru" ? russianOrigin : englishOrigin;
  const organizationId = `${origin}/#organization`;
  const websiteId = `${origin}/#website`;
  const appId = `${origin}/#software-application`;

  return {
    "@context": "https://schema.org",
    "@graph": [
      {
        "@type": "Organization",
        "@id": organizationId,
        name: "NERIVA",
        url: `${origin}/`,
        logo: `${origin}/assets/brand-logo-mini.png`,
        sameAs: socialProfileUrls,
      },
      {
        "@type": "WebSite",
        "@id": websiteId,
        name: "NERIVA",
        url: `${origin}/`,
        inLanguage: language,
        publisher: { "@id": organizationId },
      },
      {
        "@type": "SoftwareApplication",
        "@id": appId,
        name: "NERIVA",
        applicationCategory: "EducationalApplication",
        operatingSystem: "Web, PWA, Telegram",
        url: canonicalUrl,
        description: copy.description,
        publisher: { "@id": organizationId },
        offers: [
          { "@type": "Offer", name: "Free", price: "0", priceCurrency: "RUB" },
          { "@type": "Offer", name: "Premium", price: "300", priceCurrency: "RUB" },
          { "@type": "Offer", name: "Platinum", price: "590", priceCurrency: "RUB" },
        ],
      },
      {
        "@type": "WebPage",
        "@id": `${canonicalUrl}#webpage`,
        url: canonicalUrl,
        name: copy.title,
        description: copy.description,
        inLanguage: language,
        isPartOf: { "@id": websiteId },
        about: { "@id": appId },
      },
      {
        "@type": "BreadcrumbList",
        "@id": `${canonicalUrl}#breadcrumb`,
        itemListElement: [
          {
            "@type": "ListItem",
            position: 1,
            name: "NERIVA",
            item: `${origin}/neriva.html${language === "en" ? "?lang=en" : ""}`,
          },
          { "@type": "ListItem", position: 2, name: copy.h1, item: canonicalUrl },
        ],
      },
      {
        "@type": "FAQPage",
        "@id": `${canonicalUrl}#faq`,
        url: canonicalUrl,
        inLanguage: language,
        mainEntity: copy.faq.map(([question, answer]) => ({
          "@type": "Question",
          name: question,
          acceptedAnswer: { "@type": "Answer", text: answer },
        })),
      },
    ],
  };
}

function renderPage(page, language) {
  const copy = page[language];
  const otherLanguage = language === "ru" ? "en" : "ru";
  const canonicalUrl = canonicalFor(page, language);
  const ruUrl = canonicalFor(page, "ru");
  const enUrl = canonicalFor(page, "en");
  const origin = language === "ru" ? russianOrigin : englishOrigin;
  const landingHref = language === "en" ? "/neriva.html?lang=en" : "/neriva.html";
  const relatedPages = seoPages.filter((item) => item.slug !== page.slug);
  const ogImage = `${origin}/assets/seo/${language === "ru" ? "neriva-social-preview-v2-ru.jpg" : "neriva-social-preview-v2-en.jpg"}`;
  const appCta = language === "ru" ? "Открыть web app" : "Open web app";
  const landingLabel = language === "ru" ? "Главная" : "Landing";
  const guidesLabel = language === "ru" ? "Материалы" : "Guides";
  const switchLabel = language === "ru" ? "English" : "Русский";
  const relatedTitle = language === "ru" ? "Другие материалы" : "More guides";
  const faqTitle = language === "ru" ? "Вопросы и ответы" : "Questions and answers";
  const docsLabel = language === "ru" ? "Документы" : "Documents";
  const companionLine =
    language === "ru"
      ? "Web app - основной продукт, Telegram - быстрый companion-канал для короткой практики."
      : "Web app first, with Telegram as an optional companion channel for short practice.";

  return `<!doctype html>
<html lang="${language}" data-site-page="seo-guide">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>${escapeHtml(copy.title)}</title>
    <meta name="description" content="${escapeHtml(copy.description)}" />
    <meta name="robots" content="index, follow" />
    <link rel="canonical" href="${canonicalUrl}" />
    <link rel="alternate" hreflang="ru" href="${ruUrl}" />
    <link rel="alternate" hreflang="en" href="${enUrl}" />
    <link rel="alternate" hreflang="x-default" href="${enUrl}" />
    <meta property="og:title" content="${escapeHtml(copy.title)}" />
    <meta property="og:description" content="${escapeHtml(copy.description)}" />
    <meta property="og:type" content="article" />
    <meta property="og:url" content="${canonicalUrl}" />
    <meta property="og:image" content="${ogImage}" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="${escapeHtml(copy.title)}" />
    <meta name="twitter:description" content="${escapeHtml(copy.description)}" />
    <meta name="twitter:image" content="${ogImage}" />
    <meta name="theme-color" content="#08111f" />
    <link rel="icon" href="/assets/brand-logo-mini.png?v=20260522-react" />
    <link rel="stylesheet" href="/assets/seo-pages.css" />
    <script type="application/ld+json">${JSON.stringify(buildJsonLd(page, language))}</script>
  </head>
  <body>
    <main class="seo-page">
      <header class="seo-header">
        <a class="seo-brand" href="${landingHref}" aria-label="NERIVA">
          <img src="/assets/brand-logo-mini.png" alt="" />
          <span>NERIVA</span>
        </a>
        <nav class="seo-nav" aria-label="${escapeHtml(guidesLabel)}">
          <a href="${landingHref}">${escapeHtml(landingLabel)}</a>
          <a data-entry="web-app" class="seo-nav__cta" href="/app/">${escapeHtml(appCta)}</a>
          <a href="${page[otherLanguage].path}">${escapeHtml(switchLabel)}</a>
        </nav>
      </header>
      <article class="seo-article">
        <section class="seo-hero">
          <p class="seo-eyebrow">NERIVA web app</p>
          <h1>${escapeHtml(copy.h1)}</h1>
          <p>${escapeHtml(copy.intro)}</p>
          <p class="seo-companion">${escapeHtml(companionLine)}</p>
          <a data-entry="web-app" class="seo-button" href="/app/">${escapeHtml(appCta)}</a>
        </section>
        <section class="seo-sections" aria-label="${escapeHtml(copy.h1)}">
          ${copy.sections.map(([title, body]) => `<section class="seo-section"><h2>${escapeHtml(title)}</h2><p>${escapeHtml(body)}</p></section>`).join("\n          ")}
        </section>
        <section class="seo-faq" id="faq">
          <h2>${escapeHtml(faqTitle)}</h2>
          ${copy.faq.map(([question, answer]) => `<article class="seo-faq__item"><h3>${escapeHtml(question)}</h3><p>${escapeHtml(answer)}</p></article>`).join("\n          ")}
        </section>
        <section class="seo-related" aria-label="${escapeHtml(relatedTitle)}">
          <h2>${escapeHtml(relatedTitle)}</h2>
          <div class="seo-related__grid">
            ${relatedPages.map((item) => `<a href="${item[language].path}">${escapeHtml(item[language].h1)}</a>`).join("\n            ")}
          </div>
        </section>
      </article>
      <footer class="seo-footer">
        <nav aria-label="NERIVA">
          <a data-entry="landing" href="${landingHref}">${escapeHtml(landingLabel)}</a>
          <a data-entry="web-app" href="/app/">${escapeHtml(appCta)}</a>
        </nav>
        <nav aria-label="${escapeHtml(docsLabel)}">
          ${legalLinks[language].map((link) => `<a href="${link.href}">${escapeHtml(link.label)}</a>`).join("\n          ")}
        </nav>
        <nav aria-label="Social">
          ${socialProfileUrls.map((url) => `<a href="${url}" target="_blank" rel="noreferrer">${new URL(url).hostname.replace("www.", "")}</a>`).join("\n          ")}
        </nav>
      </footer>
    </main>
  </body>
</html>
`;
}

function renderSitemap() {
  const landingRuUrl = `${russianOrigin}/neriva.html`;
  const landingEnUrl = `${englishOrigin}/neriva.html?lang=en`;
  const landingAlternates = [
    `<xhtml:link rel="alternate" hreflang="ru" href="${landingRuUrl}" />`,
    `<xhtml:link rel="alternate" hreflang="en" href="${landingEnUrl}" />`,
    `<xhtml:link rel="alternate" hreflang="x-default" href="${landingEnUrl}" />`,
  ];
  const urls = [
    `<url>
    <loc>${landingRuUrl}</loc>
    ${landingAlternates.join("\n    ")}
  </url>`,
    `<url>
    <loc>${landingEnUrl}</loc>
    ${landingAlternates.join("\n    ")}
  </url>`,
  ];

  for (const page of seoPages) {
    const alternates = [
      `<xhtml:link rel="alternate" hreflang="ru" href="${canonicalFor(page, "ru")}" />`,
      `<xhtml:link rel="alternate" hreflang="en" href="${canonicalFor(page, "en")}" />`,
      `<xhtml:link rel="alternate" hreflang="x-default" href="${canonicalFor(page, "en")}" />`,
    ].join("\n    ");
    urls.push(`<url>
    <loc>${canonicalFor(page, "ru")}</loc>
    ${alternates}
  </url>`);
    urls.push(`<url>
    <loc>${canonicalFor(page, "en")}</loc>
    ${alternates}
  </url>`);
  }

  for (const legalPath of ["/privacy.html", "/terms.html", "/agreement.html", "/consent.html"]) {
    urls.push(`<url><loc>${russianOrigin}${legalPath}</loc></url>`);
    urls.push(`<url><loc>${englishOrigin}${legalPath}?lang=en</loc></url>`);
  }

  return `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">
  ${urls.join("\n  ")}
</urlset>
`;
}

await fs.mkdir(path.join(publicRoot, "en"), { recursive: true });

for (const page of seoPages) {
  await fs.writeFile(path.join(publicRoot, `${page.slug}.html`), renderPage(page, "ru"), "utf8");
  await fs.writeFile(path.join(publicRoot, "en", `${page.slug}.html`), renderPage(page, "en"), "utf8");
}

await fs.writeFile(path.join(publicRoot, "sitemap.xml"), renderSitemap(), "utf8");
console.log(`Generated ${seoPages.length * 2} static SEO pages and sitemap.xml`);
