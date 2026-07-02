import { expect, test } from "@playwright/test";

const pages = [
  {
    slug: "ai-english-tutor",
    ruPath: "/ai-english-tutor.html",
    enPath: "/en/ai-english-tutor.html",
    ruTitle: "AI-репетитор английского онлайн - NERIVA",
    enTitle: "AI English Tutor Online - NERIVA",
    ruH1: "AI-репетитор английского в веб-приложении NERIVA",
    enH1: "AI English tutor inside the NERIVA web app",
  },
  {
    slug: "english-speaking-practice",
    ruPath: "/english-speaking-practice.html",
    enPath: "/en/english-speaking-practice.html",
    ruTitle: "Разговорная практика английского с AI - NERIVA",
    enTitle: "English Speaking Practice With AI - NERIVA",
    ruH1: "Разговорная практика английского с AI",
    enH1: "English speaking practice with AI",
  },
  {
    slug: "english-pronunciation-trainer",
    ruPath: "/english-pronunciation-trainer.html",
    enPath: "/en/english-pronunciation-trainer.html",
    ruTitle: "Тренажер произношения английского с AI - NERIVA",
    enTitle: "English Pronunciation Trainer With AI - NERIVA",
    ruH1: "Тренажер произношения английского с AI",
    enH1: "English pronunciation trainer with AI",
  },
  {
    slug: "english-for-work-and-travel",
    ruPath: "/english-for-work-and-travel.html",
    enPath: "/en/english-for-work-and-travel.html",
    ruTitle: "Английский для работы и путешествий - NERIVA",
    enTitle: "English for Work and Travel - NERIVA",
    ruH1: "Английский для работы и путешествий",
    enH1: "English for work and travel",
  },
  {
    slug: "language-learning-web-app",
    ruPath: "/language-learning-web-app.html",
    enPath: "/en/language-learning-web-app.html",
    ruTitle: "Веб-приложение для изучения языков - NERIVA",
    enTitle: "Language Learning Web App - NERIVA",
    ruH1: "Веб-приложение для изучения языков NERIVA",
    enH1: "NERIVA language learning web app",
  },
] as const;

const rebrandedPublicPaths = [
  "/poliglot-ai.html?lang=ru",
  "/poliglot-ai.html?lang=en",
  "/privacy.html?lang=ru",
  "/terms.html?lang=ru",
  "/agreement.html?lang=ru",
  "/consent.html?lang=ru",
  "/ai-english-tutor.html",
  "/english-speaking-practice.html",
  "/english-pronunciation-trainer.html",
  "/english-for-work-and-travel.html",
  "/language-learning-web-app.html",
  "/en/ai-english-tutor.html",
  "/en/english-speaking-practice.html",
  "/en/english-pronunciation-trainer.html",
  "/en/english-for-work-and-travel.html",
  "/en/language-learning-web-app.html",
] as const;

const legalPaths = new Set(["/privacy.html", "/terms.html", "/agreement.html", "/consent.html"]);

const oldPublicBrandAndContacts = new RegExp(
  [
    ["Poliglot", "\\s+AI"].join(""),
    ["Poliglot", "AI"].join(""),
    ["Polyglot", "\\s+AI"].join(""),
    ["AI", "\\s+Polyglot"].join(""),
    ["support", "poliglotai@gmail\\.com"].join(""),
    ["Poliglot", "_AI_bot"].join(""),
    ["poliglot", "_ai_bot"].join(""),
    "AsaselD",
    "t\\.me/neriva_app",
    ["instagram\\.com/", "poliglotai"].join(""),
    ["tiktok\\.com/@", "poliglotai"].join(""),
  ].join("|"),
  "i",
);

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    const originalMatchMedia = window.matchMedia.bind(window);
    window.matchMedia = ((query: string) => {
      if (query === "(prefers-reduced-motion: reduce)") {
        return {
          matches: true,
          media: query,
          onchange: null,
          addListener: () => undefined,
          removeListener: () => undefined,
          addEventListener: () => undefined,
          removeEventListener: () => undefined,
          dispatchEvent: () => false,
        } as MediaQueryList;
      }
      return originalMatchMedia(query);
    }) as typeof window.matchMedia;
  });
});

test.describe("static SEO pages", () => {
  test("public pages use the NERIVA brand and current support contacts", async ({ page }) => {
    for (const path of rebrandedPublicPaths) {
      const response = await page.goto(path);
      expect(response?.ok(), `${path} should be served`).toBe(true);

      const pageSnapshot = [await page.title(), await page.locator("body").innerText(), await page.content()].join("\n");
      expect(pageSnapshot, `${path} should mention NERIVA`).toContain("NERIVA");
      expect(pageSnapshot, `${path} should not expose old brand or contacts`).not.toMatch(oldPublicBrandAndContacts);

      const pathname = new URL(path, "https://neriva.ru").pathname;
      if (legalPaths.has(pathname)) {
        expect(pageSnapshot, `${path} should expose current support email`).toContain("support@neriva.ru");
        expect(pageSnapshot, `${path} should expose current Telegram app bot`).toContain("@NERIVAapp_bot");
        await expect(page.locator('a[href="mailto:support@neriva.ru"]').first()).toBeVisible();
        await expect(page.locator('a[href="https://t.me/NERIVAapp_bot"]').first()).toBeVisible();
      }
    }
  });

  test("serves Russian static SEO pages with real HTML content and metadata", async ({ page }) => {
    for (const item of pages) {
      const response = await page.goto(item.ruPath);
      expect(response?.ok(), `${item.ruPath} should be served`).toBe(true);
      await expect(page).toHaveTitle(item.ruTitle);
      await expect(page.locator("html")).toHaveAttribute("lang", "ru");
      await expect(page.locator("h1")).toHaveText(item.ruH1);
      await expect(page.locator("#root")).toHaveCount(0);
      await expect(page.locator('meta[name="description"]')).toHaveAttribute("content", /NERIVA/);
      await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", `https://neriva.ru${item.ruPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="ru"]')).toHaveAttribute("href", `https://neriva.ru${item.ruPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="en"]')).toHaveAttribute("href", `https://neriva.ru${item.enPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="x-default"]')).toHaveAttribute("href", `https://neriva.ru${item.enPath}`);
      await expect(page.locator(".seo-faq__item")).toHaveCount(5);
      await expect(page.locator(".seo-related a")).toHaveCount(4);
      await expect(page.locator('a[data-entry="web-app"]').first()).toHaveAttribute("href", "/app/");
      await expect(page.locator('.seo-footer a[data-entry="landing"]')).toHaveAttribute("href", "/poliglot-ai.html");
      await expect(page.locator('.seo-footer a[data-entry="web-app"]')).toHaveAttribute("href", "/app/");
      await expect(page.locator('.seo-footer a[href="https://t.me/NERIVAapp_bot"]')).toBeVisible();
      await expect(page.locator('meta[name="robots"]')).toHaveAttribute("content", "index, follow");
      await expect(page.locator('meta[property="og:type"]')).toHaveAttribute("content", "article");
      await expect(page.locator('meta[property="og:url"]')).toHaveAttribute("content", `https://neriva.ru${item.ruPath}`);
      await expect(page.locator('meta[name="twitter:card"]')).toHaveAttribute("content", "summary_large_image");
      await expect(page.locator('meta[name="twitter:title"]')).toHaveAttribute("content", item.ruTitle);

      const jsonLdText = await page.locator('script[type="application/ld+json"]').textContent();
      expect(jsonLdText).toBeTruthy();
      const jsonLd = JSON.parse(jsonLdText || "{}");
      const types = jsonLd["@graph"].map((node: { "@type": string }) => node["@type"]);
      expect(types).toEqual(expect.arrayContaining(["Organization", "WebSite", "SoftwareApplication", "WebPage", "BreadcrumbList", "FAQPage"]));
      const organizationNode = jsonLd["@graph"].find((node: { "@type": string }) => node["@type"] === "Organization");
      expect(organizationNode.sameAs).toContain("https://t.me/NERIVAapp_bot");
      const faqNode = jsonLd["@graph"].find((node: { "@type": string }) => node["@type"] === "FAQPage");
      expect(faqNode.mainEntity).toHaveLength(5);
      const visibleFaq = await page.locator(".seo-faq__item").evaluateAll((items) =>
        items.map((item) => ({
          question: item.querySelector("h3")?.textContent?.trim(),
          answer: item.querySelector("p")?.textContent?.trim(),
        })),
      );
      expect(faqNode.mainEntity.map((entry: { name: string; acceptedAnswer: { text: string } }) => ({
        question: entry.name,
        answer: entry.acceptedAnswer.text,
      }))).toEqual(visibleFaq);
    }
  });

  test("serves English static SEO pages with international canonical metadata", async ({ page }) => {
    for (const item of pages) {
      const response = await page.goto(item.enPath);
      expect(response?.ok(), `${item.enPath} should be served`).toBe(true);
      await expect(page).toHaveTitle(item.enTitle);
      await expect(page.locator("html")).toHaveAttribute("lang", "en");
      await expect(page.locator("h1")).toHaveText(item.enH1);
      await expect(page.locator("#root")).toHaveCount(0);
      await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", `https://neriva.ru${item.enPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="ru"]')).toHaveAttribute("href", `https://neriva.ru${item.ruPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="en"]')).toHaveAttribute("href", `https://neriva.ru${item.enPath}`);
      await expect(page.locator(".seo-page")).toContainText("web app");
      await expect(page.locator(".seo-page")).toContainText("Telegram");
      await expect(page.locator('.seo-footer a[data-entry="landing"]')).toHaveAttribute("href", "/poliglot-ai.html?lang=en");
      await expect(page.locator('.seo-footer a[data-entry="web-app"]')).toHaveAttribute("href", "/app/");
      await expect(page.locator('.seo-footer a[href="https://t.me/NERIVAapp_bot"]')).toBeVisible();
      await expect(page.locator('meta[name="robots"]')).toHaveAttribute("content", "index, follow");
      await expect(page.locator('meta[property="og:url"]')).toHaveAttribute("content", `https://neriva.ru${item.enPath}`);
      await expect(page.locator('meta[name="twitter:title"]')).toHaveAttribute("content", item.enTitle);
    }
  });

  test("sitemap exposes the bilingual SEO page cluster with alternates", async ({ page }) => {
    const response = await page.goto("/sitemap.xml");
    expect(response?.ok()).toBe(true);
    const body = await response?.text();
    expect(body).toBeTruthy();

    for (const item of pages) {
      expect(body).toContain(`https://neriva.ru${item.ruPath}`);
      expect(body).toContain(`https://neriva.ru${item.enPath}`);
      expect(body).toContain(`hreflang="ru" href="https://neriva.ru${item.ruPath}"`);
      expect(body).toContain(`hreflang="en" href="https://neriva.ru${item.enPath}"`);
      expect(body).toContain(`hreflang="x-default" href="https://neriva.ru${item.enPath}"`);
    }
  });

  test("landing exposes SEO guides only outside the primary hero and navigation", async ({ page }) => {
    await page.goto("/poliglot-ai.html?lang=en");

    const guideLinks = page.locator(".seo-guides a");
    await expect(guideLinks).toHaveCount(5);
    await expect(page.locator(".landing-hero .seo-guides")).toHaveCount(0);
    await expect(page.locator(".public-nav .seo-guides")).toHaveCount(0);
    await expect(guideLinks.first()).toHaveAttribute("href", "/en/ai-english-tutor.html");

    await page.goto("/poliglot-ai.html?lang=ru");
    const ruGuideLinks = page.locator(".seo-guides a");
    await expect(ruGuideLinks).toHaveCount(5);
    await expect(ruGuideLinks.first()).toHaveAttribute("href", "/ai-english-tutor.html");
  });
});
