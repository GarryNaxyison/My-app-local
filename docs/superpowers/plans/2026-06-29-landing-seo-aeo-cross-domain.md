# Landing SEO/AEO Cross-Domain Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add SEO/AEO metadata, JSON-LD, visible answer content, `robots.txt`, and `sitemap.xml` to the current NERIVA landing without redesigning it.

**Architecture:** Keep the landing UI in `EnglishSparkLanding`, move reusable SEO/AEO copy into one small content module, and centralize runtime head updates in `seoMetadata.ts`. The Russian domain canonicalizes Russian pages to `poliglotai.ru`; all non-Russian interface languages canonicalize to `poliglotai.online` with `?lang=<code>`.

**Tech Stack:** React 19, Vite 7, TypeScript, Playwright, static public assets, schema.org JSON-LD.

---

## File Structure

- Create `site-react/e2e/public-seo.spec.ts`: focused Playwright tests for canonical URLs, `hreflang`, JSON-LD, visible AEO answers, `robots.txt`, `sitemap.xml`, and relative `/app/` links.
- Create `site-react/src/landingSeoContent.ts`: shared language codes, localized SEO copy, canonical path helper, and FAQ/AEO question data.
- Create `site-react/src/seoMetadata.ts`: runtime metadata helper that computes canonical URLs, alternate URLs, Open Graph/Twitter tags, and JSON-LD.
- Modify `site-react/src/main.tsx`: import and initialize the metadata helper.
- Modify `site-react/poliglot-ai.html`: improve static fallback metadata for crawlers before React runs.
- Modify `site-react/src/EnglishSparkLanding.tsx`: render the compact answer-oriented FAQ/AEO section.
- Modify `site-react/src/englishSparkLanding.css`: style the new answer section using existing section/card patterns.
- Modify `site-react/public/assets/site-phrases.js`: add Russian translations for the new English base strings.
- Create `site-react/public/robots.txt`: allow crawl and advertise both controlled sitemap URLs.
- Create `site-react/public/sitemap.xml`: list `.ru` and `.online` public URLs and landing alternates.

## Task 1: Add Failing SEO/AEO Playwright Coverage

**Files:**
- Create: `site-react/e2e/public-seo.spec.ts`

- [ ] **Step 1: Create the SEO/AEO test file**

Create `site-react/e2e/public-seo.spec.ts` with this exact content:

```ts
import { expect, test } from "@playwright/test";

const expectedHreflang = [
  "ru",
  "en",
  "es",
  "de",
  "fr",
  "it",
  "zh",
  "ja",
  "ko",
  "tg",
  "uz",
  "tt",
  "hy",
  "kk",
  "ky",
  "ka",
  "uk",
  "pl",
  "ro",
  "pt",
  "ar",
  "bn",
  "cs",
  "el",
  "hi",
  "hu",
  "id",
  "nl",
  "sv",
  "ta",
  "te",
  "th",
  "tl",
  "tr",
  "vi",
  "x-default",
] as const;

test.describe("public landing SEO and AEO metadata", () => {
  test("sets Russian canonical metadata, alternates, JSON-LD, and visible answers", async ({ page }) => {
    await page.goto("/poliglot-ai.html?lang=ru");

    await expect(page).toHaveTitle("NERIVA - AI-репетитор английского и языков в Telegram");
    await expect(page.locator('meta[name="description"]')).toHaveAttribute(
      "content",
      /AI-уроки, разговорная практика, Telegram, web app, произношение, фото-перевод, ошибки и Premium/,
    );
    await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", "https://poliglotai.ru/poliglot-ai.html");

    const alternates = await page.locator('link[rel="alternate"]').evaluateAll((links) =>
      links.map((link) => ({
        hreflang: link.getAttribute("hreflang"),
        href: link.getAttribute("href"),
      })),
    );
    expect([...new Set(alternates.map((link) => link.hreflang))].sort()).toEqual([...expectedHreflang].sort());
    expect(alternates).toEqual(
      expect.arrayContaining([
        { hreflang: "ru", href: "https://poliglotai.ru/poliglot-ai.html" },
        { hreflang: "en", href: "https://poliglotai.online/poliglot-ai.html?lang=en" },
        { hreflang: "x-default", href: "https://poliglotai.online/poliglot-ai.html?lang=en" },
      ]),
    );

    const jsonLdText = await page.locator("#poliglot-seo-jsonld").textContent();
    expect(jsonLdText).toBeTruthy();
    const jsonLd = JSON.parse(jsonLdText || "{}");
    const graphTypes = jsonLd["@graph"].map((node: { "@type": string }) => node["@type"]);
    expect(graphTypes).toEqual(["Organization", "WebSite", "SoftwareApplication", "FAQPage"]);
    const faqNode = jsonLd["@graph"].find((node: { "@type": string }) => node["@type"] === "FAQPage");
    expect(faqNode.mainEntity).toHaveLength(6);
    expect(faqNode.mainEntity[0].name).toBe("Что такое NERIVA?");

    await expect(page.getByRole("heading", { name: "Ответы для поиска и AI" })).toBeVisible();
    await expect(page.locator(".answer-card")).toHaveCount(6);
    await expect(page.locator(".answer-card").first()).toContainText("AI-репетитор");
  });

  test("sets English international metadata without forcing app links away from the current host", async ({ page }) => {
    await page.goto("/poliglot-ai.html?lang=en");

    await expect(page).toHaveTitle("NERIVA - AI language tutor in Telegram and web app");
    await expect(page.locator('meta[name="description"]')).toHaveAttribute(
      "content",
      /AI language tutor, speaking practice, Telegram bot, web app, voice, photo translation, mistakes, and premium plans/,
    );
    await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", "https://poliglotai.online/poliglot-ai.html?lang=en");
    await expect(page.locator('meta[property="og:url"]')).toHaveAttribute("content", "https://poliglotai.online/poliglot-ai.html?lang=en");
    await expect(page.locator('meta[name="twitter:title"]')).toHaveAttribute("content", "NERIVA - AI language tutor in Telegram and web app");

    const jsonLdText = await page.locator("#poliglot-seo-jsonld").textContent();
    const jsonLd = JSON.parse(jsonLdText || "{}");
    const appNode = jsonLd["@graph"].find((node: { "@type": string }) => node["@type"] === "SoftwareApplication");
    expect(appNode.offers).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ name: "Free", price: "0", priceCurrency: "RUB" }),
        expect.objectContaining({ name: "Premium", price: "300", priceCurrency: "RUB" }),
        expect.objectContaining({ name: "Platinum", price: "590", priceCurrency: "RUB" }),
      ]),
    );

    await expect(page.getByRole("heading", { name: "Answers for search and AI assistants" })).toBeVisible();
    await expect(page.locator(".answer-card").first()).toContainText("AI language tutor");

    const webEntryHrefs = await page.locator('a[data-entry="web-app"]').evaluateAll((links) =>
      links.map((link) => (link as HTMLAnchorElement).getAttribute("href")),
    );
    expect(webEntryHrefs).toEqual(["/app/", "/app/", "/app/"]);
  });

  test("serves robots and sitemap search files", async ({ page }) => {
    const robotsResponse = await page.goto("/robots.txt");
    const robotsBody = await robotsResponse?.text();
    expect(robotsBody).toContain("User-agent: *");
    expect(robotsBody).toContain("Sitemap: https://poliglotai.ru/sitemap.xml");
    expect(robotsBody).toContain("Sitemap: https://poliglotai.online/sitemap.xml");

    const sitemapResponse = await page.goto("/sitemap.xml");
    const sitemapBody = await sitemapResponse?.text();
    expect(sitemapBody).toContain("https://poliglotai.ru/poliglot-ai.html");
    expect(sitemapBody).toContain("https://poliglotai.online/poliglot-ai.html?lang=en");
    expect(sitemapBody).toContain('hreflang="x-default"');
  });
});
```

- [ ] **Step 2: Run the new tests and confirm they fail before implementation**

Run:

```bash
npm --prefix site-react run build
npm --prefix site-react run e2e -- --grep "public landing SEO and AEO metadata"
```

Expected result: `build` succeeds, then the Playwright grep fails because `site-react/e2e/public-seo.spec.ts` expects metadata, JSON-LD, answer cards, `robots.txt`, and `sitemap.xml` that do not exist yet.

- [ ] **Step 3: Commit the failing tests**

Run:

```bash
git add site-react/e2e/public-seo.spec.ts
git commit -m "test: cover public landing seo metadata"
```

## Task 2: Add Shared SEO/AEO Content And Runtime Metadata

**Files:**
- Create: `site-react/src/landingSeoContent.ts`
- Create: `site-react/src/seoMetadata.ts`
- Modify: `site-react/src/main.tsx`
- Modify: `site-react/poliglot-ai.html`

- [ ] **Step 1: Create shared SEO/AEO content**

Create `site-react/src/landingSeoContent.ts` with this exact content:

```ts
export const landingLanguageCodes = [
  "ru",
  "en",
  "es",
  "de",
  "fr",
  "it",
  "zh",
  "ja",
  "ko",
  "tg",
  "uz",
  "tt",
  "hy",
  "kk",
  "ky",
  "ka",
  "uk",
  "pl",
  "ro",
  "pt",
  "ar",
  "bn",
  "cs",
  "el",
  "hi",
  "hu",
  "id",
  "nl",
  "sv",
  "ta",
  "te",
  "th",
  "tl",
  "tr",
  "vi",
] as const;

export type LandingLanguageCode = (typeof landingLanguageCodes)[number];
export type LandingSeoLocale = "ru" | "en";

export const russianSiteOrigin = "https://poliglotai.ru";
export const internationalSiteOrigin = "https://poliglotai.online";

export const socialProfileUrls = [
  "https://www.youtube.com/@NERIVA",
  "https://www.instagram.com/poliglotai.online/",
  "https://www.tiktok.com/@poliglotai.online",
] as const;

export const landingSeoCopy = {
  ru: {
    title: "NERIVA - AI-репетитор английского и языков в Telegram",
    description:
      "NERIVA: AI-уроки, разговорная практика, Telegram, web app, произношение, фото-перевод, ошибки и Premium в одном профиле.",
    sectionEyebrow: "Поисковые ответы",
    sectionTitle: "Ответы для поиска и AI",
    sectionIntro:
      "Короткие прямые ответы о NERIVA для людей, которые сравнивают AI-репетиторов, языковые Telegram-боты, speaking practice, голос и фото-перевод.",
    questions: [
      {
        question: "Что такое NERIVA?",
        answer:
          "NERIVA — это AI-репетитор языков в web app и Telegram-боте. Он объединяет короткие уроки, разговорную практику, произношение, перевод текста с фото, ошибки, заметки и прогресс в одном профиле.",
      },
      {
        question: "Можно ли учить английский с ИИ в Telegram?",
        answer:
          "Да. В Telegram можно запускать практику, получать задания, отправлять ответы и голосовые сообщения, а прогресс сохраняется вместе с web app.",
      },
      {
        question: "Чем AI-репетитор отличается от обычного приложения со словами?",
        answer:
          "NERIVA не ограничивается списком слов: он дает фразу в контексте, просит ответить, исправляет ошибку и возвращает слабое место в повторение.",
      },
      {
        question: "Можно ли тренировать произношение и speaking?",
        answer:
          "Да. Voice Coach и shadowing помогают тренировать речь, видеть слабые слова, получать score и повторять более естественную фразу.",
      },
      {
        question: "Можно ли переводить текст с фото?",
        answer:
          "Да. Фото меню, вывески или задания превращается в перевод, заметку и короткую практику по этому контексту.",
      },
      {
        question: "Есть ли бесплатный тариф?",
        answer:
          "Да. Free дает стартовые уроки и практику без оплаты, а Premium и Platinum открывают больше лимитов, голос, фото и интенсивную ежедневную учебу.",
      },
    ],
  },
  en: {
    title: "NERIVA - AI language tutor in Telegram and web app",
    description:
      "NERIVA: AI language tutor, speaking practice, Telegram bot, web app, voice, photo translation, mistakes, and premium plans in one profile.",
    sectionEyebrow: "Search answers",
    sectionTitle: "Answers for search and AI assistants",
    sectionIntro:
      "Short direct answers about NERIVA for people comparing AI tutors, Telegram language bots, speaking practice, voice, and photo translation.",
    questions: [
      {
        question: "What is NERIVA?",
        answer:
          "NERIVA is an AI language tutor in a web app and Telegram bot. It combines short lessons, speaking practice, pronunciation, photo translation, mistakes, notes, and progress in one profile.",
      },
      {
        question: "Can I practice English with an AI tutor in Telegram?",
        answer:
          "Yes. In Telegram you can start practice, receive tasks, send answers and voice messages, while progress stays synced with the web app.",
      },
      {
        question: "How is an AI tutor different from a vocabulary app?",
        answer:
          "NERIVA is not only a word list. It gives a phrase in context, asks for your answer, corrects the mistake, and brings the weak spot back for review.",
      },
      {
        question: "Can I practice pronunciation and speaking?",
        answer:
          "Yes. Voice Coach and shadowing help you practice speech, see weak words, get a score, and repeat a more natural phrase.",
      },
      {
        question: "Can I translate text from photos?",
        answer:
          "Yes. A photo of a menu, sign, or exercise becomes a translation, note, and short practice prompt in that context.",
      },
      {
        question: "Is there a free plan?",
        answer:
          "Yes. Free gives starter lessons and practice without payment, while Premium and Platinum unlock higher limits, voice, photo tools, and intensive daily study.",
      },
    ],
  },
} as const;

export function seoLocaleForLanguage(language: string): LandingSeoLocale {
  return language === "ru" ? "ru" : "en";
}

export function normalizeLandingLanguage(value: string | null | undefined): LandingLanguageCode {
  const normalized = String(value || "").toLowerCase().trim();
  return landingLanguageCodes.includes(normalized as LandingLanguageCode) ? (normalized as LandingLanguageCode) : "ru";
}

export function canonicalUrlForLanguage(language: LandingLanguageCode): string {
  const origin = language === "ru" ? russianSiteOrigin : internationalSiteOrigin;
  const url = new URL("/poliglot-ai.html", origin);
  if (language !== "ru") url.searchParams.set("lang", language);
  return url.toString();
}

export function canonicalOriginForLanguage(language: LandingLanguageCode): string {
  return language === "ru" ? russianSiteOrigin : internationalSiteOrigin;
}
```

- [ ] **Step 2: Create the runtime metadata helper**

Create `site-react/src/seoMetadata.ts` with this exact content:

```ts
import {
  canonicalOriginForLanguage,
  canonicalUrlForLanguage,
  landingLanguageCodes,
  landingSeoCopy,
  normalizeLandingLanguage,
  seoLocaleForLanguage,
  socialProfileUrls,
} from "./landingSeoContent";

declare global {
  interface Window {
    poliglotSiteI18n?: {
      currentLanguage: () => string;
    };
  }
}

const jsonLdScriptId = "poliglot-seo-jsonld";
const supportEmail = "support@neriva.ru";

function currentLandingLanguage() {
  const i18nLanguage = window.poliglotSiteI18n?.currentLanguage?.();
  if (i18nLanguage) return normalizeLandingLanguage(i18nLanguage);

  const params = new URLSearchParams(window.location.search);
  return normalizeLandingLanguage(params.get("lang") || localStorage.getItem("poliglot_site_language") || document.documentElement.lang || "ru");
}

function upsertMetaByName(name: string, content: string) {
  let meta = document.head.querySelector<HTMLMetaElement>(`meta[name="${name}"]`);
  if (!meta) {
    meta = document.createElement("meta");
    meta.setAttribute("name", name);
    document.head.append(meta);
  }
  meta.setAttribute("content", content);
}

function upsertMetaByProperty(property: string, content: string) {
  let meta = document.head.querySelector<HTMLMetaElement>(`meta[property="${property}"]`);
  if (!meta) {
    meta = document.createElement("meta");
    meta.setAttribute("property", property);
    document.head.append(meta);
  }
  meta.setAttribute("content", content);
}

function upsertLink(selector: string, attributes: Record<string, string>) {
  let link = document.head.querySelector<HTMLLinkElement>(selector);
  if (!link) {
    link = document.createElement("link");
    document.head.append(link);
  }
  Object.entries(attributes).forEach(([name, value]) => link.setAttribute(name, value));
}

function absoluteAssetUrl(path: string, origin: string) {
  return new URL(path, origin).toString();
}

function buildJsonLd(language: ReturnType<typeof currentLandingLanguage>) {
  const locale = seoLocaleForLanguage(language);
  const copy = landingSeoCopy[locale];
  const canonicalUrl = canonicalUrlForLanguage(language);
  const origin = canonicalOriginForLanguage(language);
  const organizationId = `${origin}/#organization`;
  const websiteId = `${origin}/#website`;
  const applicationId = `${origin}/#software-application`;

  return {
    "@context": "https://schema.org",
    "@graph": [
      {
        "@type": "Organization",
        "@id": organizationId,
        name: "NERIVA",
        url: `${origin}/`,
        logo: absoluteAssetUrl("/assets/brand-logo-mini.png", origin),
        sameAs: [...socialProfileUrls],
        email: supportEmail,
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
        "@id": applicationId,
        name: "NERIVA",
        applicationCategory: "EducationalApplication",
        operatingSystem: "Web, Telegram",
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
        "@type": "FAQPage",
        "@id": `${canonicalUrl}#faq`,
        url: canonicalUrl,
        inLanguage: language,
        mainEntity: copy.questions.map((item) => ({
          "@type": "Question",
          name: item.question,
          acceptedAnswer: {
            "@type": "Answer",
            text: item.answer,
          },
        })),
      },
    ],
  };
}

export function updateLandingSeoMetadata() {
  const language = currentLandingLanguage();
  const locale = seoLocaleForLanguage(language);
  const copy = landingSeoCopy[locale];
  const canonicalUrl = canonicalUrlForLanguage(language);
  const origin = canonicalOriginForLanguage(language);
  const imageUrl = absoluteAssetUrl("/assets/brand-logo-mini.png", origin);

  document.title = copy.title;
  upsertMetaByName("description", copy.description);
  upsertMetaByName("robots", "index, follow");

  upsertLink('link[rel="canonical"]', { rel: "canonical", href: canonicalUrl });
  landingLanguageCodes.forEach((code) => {
    upsertLink(`link[rel="alternate"][hreflang="${code}"]`, {
      rel: "alternate",
      hreflang: code,
      href: canonicalUrlForLanguage(code),
    });
  });
  upsertLink('link[rel="alternate"][hreflang="x-default"]', {
    rel: "alternate",
    hreflang: "x-default",
    href: canonicalUrlForLanguage("en"),
  });

  upsertMetaByProperty("og:title", copy.title);
  upsertMetaByProperty("og:description", copy.description);
  upsertMetaByProperty("og:type", "website");
  upsertMetaByProperty("og:url", canonicalUrl);
  upsertMetaByProperty("og:image", imageUrl);

  upsertMetaByName("twitter:card", "summary");
  upsertMetaByName("twitter:title", copy.title);
  upsertMetaByName("twitter:description", copy.description);
  upsertMetaByName("twitter:image", imageUrl);

  let script = document.getElementById(jsonLdScriptId) as HTMLScriptElement | null;
  if (!script) {
    script = document.createElement("script");
    script.id = jsonLdScriptId;
    script.type = "application/ld+json";
    document.head.append(script);
  }
  script.textContent = JSON.stringify(buildJsonLd(language));
}

export function setupLandingSeoMetadata() {
  const apply = () => updateLandingSeoMetadata();
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", apply, { once: true });
  } else {
    apply();
  }
  window.addEventListener("poliglot-language-change", apply);
}
```

- [ ] **Step 3: Initialize the metadata helper**

Modify `site-react/src/main.tsx` so the full file is:

```ts
import React from "react";
import { createRoot } from "react-dom/client";
import { PublicSiteApp } from "./PublicSiteApp";
import { setupLandingSeoMetadata } from "./seoMetadata";
import "./styles.css";
import "./englishSparkLanding.css";

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <PublicSiteApp />
  </React.StrictMode>,
);

setupLandingSeoMetadata();
```

- [ ] **Step 4: Update static landing metadata**

Modify only the `<head>` of `site-react/poliglot-ai.html` so it contains these metadata tags while keeping the existing script order:

```html
    <title>NERIVA - AI-репетитор английского и языков в Telegram</title>
    <meta name="description" content="NERIVA: AI-уроки, разговорная практика, Telegram, web app, произношение, фото-перевод, ошибки и Premium в одном профиле." />
    <meta name="robots" content="index, follow" />
    <link rel="canonical" href="https://poliglotai.ru/poliglot-ai.html" />
    <link rel="alternate" hreflang="ru" href="https://poliglotai.ru/poliglot-ai.html" />
    <link rel="alternate" hreflang="en" href="https://poliglotai.online/poliglot-ai.html?lang=en" />
    <link rel="alternate" hreflang="x-default" href="https://poliglotai.online/poliglot-ai.html?lang=en" />
    <meta property="og:title" content="NERIVA - AI-репетитор английского и языков в Telegram" />
    <meta property="og:description" content="AI-уроки, разговорная практика, Telegram, web app, произношение, фото-перевод, ошибки и Premium в одном профиле." />
    <meta property="og:type" content="website" />
    <meta property="og:url" content="https://poliglotai.ru/poliglot-ai.html" />
    <meta property="og:image" content="https://poliglotai.ru/assets/brand-logo-mini.png" />
    <meta name="twitter:card" content="summary" />
    <meta name="twitter:title" content="NERIVA - AI-репетитор английского и языков в Telegram" />
    <meta name="twitter:description" content="AI-уроки, разговорная практика, Telegram, web app, произношение, фото-перевод, ошибки и Premium в одном профиле." />
    <meta name="twitter:image" content="https://poliglotai.ru/assets/brand-logo-mini.png" />
```

- [ ] **Step 5: Run metadata tests and confirm only content/search-file assertions still fail**

Run:

```bash
npm --prefix site-react run build
npm --prefix site-react run e2e -- --grep "public landing SEO and AEO metadata"
```

Expected result: TypeScript and Vite build pass. Playwright still fails on `.answer-card`, `robots.txt`, and `sitemap.xml` assertions because those files and visible answer cards are added in later tasks.

- [ ] **Step 6: Commit metadata implementation**

Run:

```bash
git add site-react/src/landingSeoContent.ts site-react/src/seoMetadata.ts site-react/src/main.tsx site-react/poliglot-ai.html
git commit -m "feat: add landing seo metadata"
```

## Task 3: Add Visible FAQ/AEO Answer Section

**Files:**
- Modify: `site-react/src/EnglishSparkLanding.tsx`
- Modify: `site-react/src/englishSparkLanding.css`
- Modify: `site-react/public/assets/site-phrases.js`

- [ ] **Step 1: Import shared SEO copy in the landing**

In `site-react/src/EnglishSparkLanding.tsx`, add this import after the existing local imports:

```ts
import { landingSeoCopy } from "./landingSeoContent";
```

- [ ] **Step 2: Render the answer section after pricing and before reviews**

In `site-react/src/EnglishSparkLanding.tsx`, insert this JSX block immediately after the closing `</section>` for `pricing-section` and before `<section id="reviews" className="spark-section reviews-section">`:

```tsx
      <section className="spark-section answer-section" aria-labelledby="answer-section-heading">
        <div className="section-copy">
          <span className="eyebrow">{landingSeoCopy.en.sectionEyebrow}</span>
          <h2 id="answer-section-heading">{landingSeoCopy.en.sectionTitle}</h2>
          <p>{landingSeoCopy.en.sectionIntro}</p>
        </div>
        <div className="answer-grid">
          {landingSeoCopy.en.questions.map((item) => (
            <article className="answer-card" key={item.question}>
              <h3>{item.question}</h3>
              <p>{item.answer}</p>
            </article>
          ))}
        </div>
      </section>
```

- [ ] **Step 3: Add answer section CSS to shared section selectors**

In `site-react/src/englishSparkLanding.css`, update the shared section selector near the top to include `.answer-section`:

```css
.spark-section,
.pricing-section,
.memory-loop-section,
.final-cta-section,
.answer-section,
.course-strip,
.device-flow {
  max-width: 1180px;
  margin: 0 auto;
}

.spark-section,
.pricing-section,
.memory-loop-section,
.final-cta-section,
.answer-section {
  padding: 104px 24px;
}
```

- [ ] **Step 4: Add answer card CSS**

In `site-react/src/englishSparkLanding.css`, add these blocks after the existing `.reviews-grid` and card grid definitions:

```css
.answer-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin-top: 34px;
}

.answer-card {
  min-width: 0;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.06);
  padding: 22px;
  color: #f7fbff;
  box-shadow: none;
  overflow-wrap: anywhere;
}

.answer-card h3 {
  margin: 0;
  color: #ffffff;
  font-size: 1.05rem;
  line-height: 1.25;
}

.answer-card p {
  margin: 12px 0 0;
  color: rgba(236, 246, 255, 0.74);
}

html[data-site-theme="light"] .answer-section {
  color: #07111f;
}

html[data-site-theme="light"] .answer-card {
  border-color: rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.9);
  color: #07111f;
  box-shadow: 0 18px 42px rgba(15, 23, 42, 0.08);
}

html[data-site-theme="light"] .answer-card h3 {
  color: #07111f;
}

html[data-site-theme="light"] .answer-card p {
  color: #475569;
}
```

- [ ] **Step 5: Add the answer grid to the mobile single-column selector**

In the existing mobile media block in `site-react/src/englishSparkLanding.css`, add `.answer-grid` to the grid selector so it becomes:

```css
  .daily-steps,
  .module-grid,
  .entry-grid,
  .pricing-grid,
  .reviews-grid,
  .answer-grid,
  .memory-grid,
  .course-cards,
  .course-stats,
  .cockpit-kpis,
  .cockpit-metric-grid,
  .scenario-grid,
  .device-flow__grid {
    grid-template-columns: 1fr;
  }
```

- [ ] **Step 6: Add Russian translations for new answer strings**

In `site-react/public/assets/site-phrases.js`, add these entries near the other English landing phrases inside `window.poliglotPhraseTranslations`:

```js
  "Answers for search and AI assistants": {
    "ru": "Ответы для поиска и AI"
  },
  "Search answers": {
    "ru": "Поисковые ответы"
  },
  "Short direct answers about NERIVA for people comparing AI tutors, Telegram language bots, speaking practice, voice, and photo translation.": {
    "ru": "Короткие прямые ответы о NERIVA для людей, которые сравнивают AI-репетиторов, языковые Telegram-боты, speaking practice, голос и фото-перевод."
  },
  "What is NERIVA?": {
    "ru": "Что такое NERIVA?"
  },
  "NERIVA is an AI language tutor in a web app and Telegram bot. It combines short lessons, speaking practice, pronunciation, photo translation, mistakes, notes, and progress in one profile.": {
    "ru": "NERIVA — это AI-репетитор языков в web app и Telegram-боте. Он объединяет короткие уроки, разговорную практику, произношение, перевод текста с фото, ошибки, заметки и прогресс в одном профиле."
  },
  "Can I practice English with an AI tutor in Telegram?": {
    "ru": "Можно ли учить английский с ИИ в Telegram?"
  },
  "Yes. In Telegram you can start practice, receive tasks, send answers and voice messages, while progress stays synced with the web app.": {
    "ru": "Да. В Telegram можно запускать практику, получать задания, отправлять ответы и голосовые сообщения, а прогресс сохраняется вместе с web app."
  },
  "How is an AI tutor different from a vocabulary app?": {
    "ru": "Чем AI-репетитор отличается от обычного приложения со словами?"
  },
  "NERIVA is not only a word list. It gives a phrase in context, asks for your answer, corrects the mistake, and brings the weak spot back for review.": {
    "ru": "NERIVA не ограничивается списком слов: он дает фразу в контексте, просит ответить, исправляет ошибку и возвращает слабое место в повторение."
  },
  "Can I practice pronunciation and speaking?": {
    "ru": "Можно ли тренировать произношение и speaking?"
  },
  "Yes. Voice Coach and shadowing help you practice speech, see weak words, get a score, and repeat a more natural phrase.": {
    "ru": "Да. Voice Coach и shadowing помогают тренировать речь, видеть слабые слова, получать score и повторять более естественную фразу."
  },
  "Can I translate text from photos?": {
    "ru": "Можно ли переводить текст с фото?"
  },
  "Yes. A photo of a menu, sign, or exercise becomes a translation, note, and short practice prompt in that context.": {
    "ru": "Да. Фото меню, вывески или задания превращается в перевод, заметку и короткую практику по этому контексту."
  },
  "Is there a free plan?": {
    "ru": "Есть ли бесплатный тариф?"
  },
  "Yes. Free gives starter lessons and practice without payment, while Premium and Platinum unlock higher limits, voice, photo tools, and intensive daily study.": {
    "ru": "Да. Free дает стартовые уроки и практику без оплаты, а Premium и Platinum открывают больше лимитов, голос, фото и интенсивную ежедневную учебу."
  },
```

- [ ] **Step 7: Run visible answer tests and confirm search-file assertions still fail**

Run:

```bash
npm --prefix site-react run build
npm --prefix site-react run e2e -- --grep "public landing SEO and AEO metadata"
```

Expected result: build passes, metadata and visible answer assertions pass, and the search-file test still fails because `robots.txt` and `sitemap.xml` are not present.

- [ ] **Step 8: Commit visible answer section**

Run:

```bash
git add site-react/src/EnglishSparkLanding.tsx site-react/src/englishSparkLanding.css site-react/public/assets/site-phrases.js
git commit -m "feat: add landing answer section"
```

## Task 4: Add Robots And Sitemap Files

**Files:**
- Create: `site-react/public/robots.txt`
- Create: `site-react/public/sitemap.xml`

- [ ] **Step 1: Create robots.txt**

Create `site-react/public/robots.txt` with this exact content:

```txt
User-agent: *
Allow: /

Sitemap: https://poliglotai.ru/sitemap.xml
Sitemap: https://poliglotai.online/sitemap.xml
```

- [ ] **Step 2: Create sitemap.xml**

Create `site-react/public/sitemap.xml` with this exact content:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">
  <url>
    <loc>https://poliglotai.ru/poliglot-ai.html</loc>
    <xhtml:link rel="alternate" hreflang="ru" href="https://poliglotai.ru/poliglot-ai.html" />
    <xhtml:link rel="alternate" hreflang="en" href="https://poliglotai.online/poliglot-ai.html?lang=en" />
    <xhtml:link rel="alternate" hreflang="es" href="https://poliglotai.online/poliglot-ai.html?lang=es" />
    <xhtml:link rel="alternate" hreflang="de" href="https://poliglotai.online/poliglot-ai.html?lang=de" />
    <xhtml:link rel="alternate" hreflang="fr" href="https://poliglotai.online/poliglot-ai.html?lang=fr" />
    <xhtml:link rel="alternate" hreflang="it" href="https://poliglotai.online/poliglot-ai.html?lang=it" />
    <xhtml:link rel="alternate" hreflang="zh" href="https://poliglotai.online/poliglot-ai.html?lang=zh" />
    <xhtml:link rel="alternate" hreflang="ja" href="https://poliglotai.online/poliglot-ai.html?lang=ja" />
    <xhtml:link rel="alternate" hreflang="ko" href="https://poliglotai.online/poliglot-ai.html?lang=ko" />
    <xhtml:link rel="alternate" hreflang="tg" href="https://poliglotai.online/poliglot-ai.html?lang=tg" />
    <xhtml:link rel="alternate" hreflang="uz" href="https://poliglotai.online/poliglot-ai.html?lang=uz" />
    <xhtml:link rel="alternate" hreflang="tt" href="https://poliglotai.online/poliglot-ai.html?lang=tt" />
    <xhtml:link rel="alternate" hreflang="hy" href="https://poliglotai.online/poliglot-ai.html?lang=hy" />
    <xhtml:link rel="alternate" hreflang="kk" href="https://poliglotai.online/poliglot-ai.html?lang=kk" />
    <xhtml:link rel="alternate" hreflang="ky" href="https://poliglotai.online/poliglot-ai.html?lang=ky" />
    <xhtml:link rel="alternate" hreflang="ka" href="https://poliglotai.online/poliglot-ai.html?lang=ka" />
    <xhtml:link rel="alternate" hreflang="uk" href="https://poliglotai.online/poliglot-ai.html?lang=uk" />
    <xhtml:link rel="alternate" hreflang="pl" href="https://poliglotai.online/poliglot-ai.html?lang=pl" />
    <xhtml:link rel="alternate" hreflang="ro" href="https://poliglotai.online/poliglot-ai.html?lang=ro" />
    <xhtml:link rel="alternate" hreflang="pt" href="https://poliglotai.online/poliglot-ai.html?lang=pt" />
    <xhtml:link rel="alternate" hreflang="ar" href="https://poliglotai.online/poliglot-ai.html?lang=ar" />
    <xhtml:link rel="alternate" hreflang="bn" href="https://poliglotai.online/poliglot-ai.html?lang=bn" />
    <xhtml:link rel="alternate" hreflang="cs" href="https://poliglotai.online/poliglot-ai.html?lang=cs" />
    <xhtml:link rel="alternate" hreflang="el" href="https://poliglotai.online/poliglot-ai.html?lang=el" />
    <xhtml:link rel="alternate" hreflang="hi" href="https://poliglotai.online/poliglot-ai.html?lang=hi" />
    <xhtml:link rel="alternate" hreflang="hu" href="https://poliglotai.online/poliglot-ai.html?lang=hu" />
    <xhtml:link rel="alternate" hreflang="id" href="https://poliglotai.online/poliglot-ai.html?lang=id" />
    <xhtml:link rel="alternate" hreflang="nl" href="https://poliglotai.online/poliglot-ai.html?lang=nl" />
    <xhtml:link rel="alternate" hreflang="sv" href="https://poliglotai.online/poliglot-ai.html?lang=sv" />
    <xhtml:link rel="alternate" hreflang="ta" href="https://poliglotai.online/poliglot-ai.html?lang=ta" />
    <xhtml:link rel="alternate" hreflang="te" href="https://poliglotai.online/poliglot-ai.html?lang=te" />
    <xhtml:link rel="alternate" hreflang="th" href="https://poliglotai.online/poliglot-ai.html?lang=th" />
    <xhtml:link rel="alternate" hreflang="tl" href="https://poliglotai.online/poliglot-ai.html?lang=tl" />
    <xhtml:link rel="alternate" hreflang="tr" href="https://poliglotai.online/poliglot-ai.html?lang=tr" />
    <xhtml:link rel="alternate" hreflang="vi" href="https://poliglotai.online/poliglot-ai.html?lang=vi" />
    <xhtml:link rel="alternate" hreflang="x-default" href="https://poliglotai.online/poliglot-ai.html?lang=en" />
  </url>
  <url>
    <loc>https://poliglotai.online/poliglot-ai.html?lang=en</loc>
    <xhtml:link rel="alternate" hreflang="ru" href="https://poliglotai.ru/poliglot-ai.html" />
    <xhtml:link rel="alternate" hreflang="en" href="https://poliglotai.online/poliglot-ai.html?lang=en" />
    <xhtml:link rel="alternate" hreflang="x-default" href="https://poliglotai.online/poliglot-ai.html?lang=en" />
  </url>
  <url><loc>https://poliglotai.ru/privacy.html</loc></url>
  <url><loc>https://poliglotai.ru/terms.html</loc></url>
  <url><loc>https://poliglotai.ru/agreement.html</loc></url>
  <url><loc>https://poliglotai.ru/consent.html</loc></url>
  <url><loc>https://poliglotai.online/privacy.html?lang=en</loc></url>
  <url><loc>https://poliglotai.online/terms.html?lang=en</loc></url>
  <url><loc>https://poliglotai.online/agreement.html?lang=en</loc></url>
  <url><loc>https://poliglotai.online/consent.html?lang=en</loc></url>
</urlset>
```

- [ ] **Step 3: Run the focused SEO/AEO tests**

Run:

```bash
npm --prefix site-react run build
npm --prefix site-react run e2e -- --grep "public landing SEO and AEO metadata"
```

Expected result: build passes and all tests in `public-seo.spec.ts` pass on desktop and mobile Playwright projects.

- [ ] **Step 4: Commit search files**

Run:

```bash
git add site-react/public/robots.txt site-react/public/sitemap.xml
git commit -m "feat: add public search files"
```

## Task 5: Full Verification, Diff Review, Push

**Files:**
- Verify all files changed by Tasks 1-4.

- [ ] **Step 1: Run full public site verification**

Run:

```bash
npm --prefix site-react run build
npm --prefix site-react run e2e
```

Expected result: build exits `0`; Playwright exits `0` for both configured projects. If Chrome remote debugging blocks Browser MCP, do not use Browser MCP as the success source; use the Playwright runner output.

- [ ] **Step 2: Review staged and unstaged changes**

Run:

```bash
git status --short
git diff -- site-react/e2e/public-seo.spec.ts site-react/src/landingSeoContent.ts site-react/src/seoMetadata.ts site-react/src/main.tsx site-react/poliglot-ai.html site-react/src/EnglishSparkLanding.tsx site-react/src/englishSparkLanding.css site-react/public/assets/site-phrases.js site-react/public/robots.txt site-react/public/sitemap.xml
```

Expected result: SEO/AEO changes are limited to the files listed in this plan. Existing unrelated local changes such as `shadowing.go`, `shadowing_test.go`, `shadowing_i18n.go`, current `site-react/e2e/public-site.spec.ts`, current `site-react/public/assets/site-i18n.js`, and `english-coach-bot` remain untouched.

- [ ] **Step 3: Commit final verification note if Task 5 introduced fixes**

If Task 5 required a correction to any SEO/AEO file, run:

```bash
git add site-react/e2e/public-seo.spec.ts site-react/src/landingSeoContent.ts site-react/src/seoMetadata.ts site-react/src/main.tsx site-react/poliglot-ai.html site-react/src/EnglishSparkLanding.tsx site-react/src/englishSparkLanding.css site-react/public/assets/site-phrases.js site-react/public/robots.txt site-react/public/sitemap.xml
git commit -m "fix: stabilize landing seo checks"
```

Expected result: a commit is created only when Task 5 changed a file. If Task 5 changed nothing, skip this commit.

- [ ] **Step 4: Push the branch**

Run:

```bash
git push origin HEAD
```

Expected result: current branch pushes to `origin`. If authentication fails, report the push failure and keep the local commits.

## Self-Review Checklist

- Spec coverage: Tasks 2-4 implement metadata, canonical URLs, `hreflang`, JSON-LD, visible FAQ/AEO content, `robots.txt`, and `sitemap.xml`.
- Domain coverage: `ru` canonical URL points to `poliglotai.ru`; every non-Russian language points to `poliglotai.online` with `?lang=<code>`.
- Testing coverage: Task 1 adds failing tests first; Tasks 2-4 make them pass; Task 5 runs full build and e2e.
- Worktree safety: every commit command stages only SEO/AEO files named in this plan.
- MCP limitation: Browser MCP was previously blocked by the system Chrome remote debugging policy, so Playwright runner verification is the required browser fallback.
