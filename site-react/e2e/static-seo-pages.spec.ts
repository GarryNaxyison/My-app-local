import { expect, test } from "@playwright/test";

const pages = [
  {
    slug: "ai-tutor",
    ruPath: "/ai-tutor.html",
    enPath: "/en/ai-tutor.html",
    ruTitle: "AI-репетитор языков онлайн - NERIVA",
    enTitle: "AI Language Tutor Online - NERIVA",
    ruH1: "AI-репетитор NERIVA для коротких уроков и повторения",
    enH1: "NERIVA AI tutor for short lessons and review",
  },
  {
    slug: "speaking-practice",
    ruPath: "/speaking-practice.html",
    enPath: "/en/speaking-practice.html",
    ruTitle: "Разговорная практика английского с AI - NERIVA",
    enTitle: "English Speaking Practice With AI - NERIVA",
    ruH1: "Разговорная практика английского с AI",
    enH1: "English speaking practice with AI",
  },
  {
    slug: "pronunciation",
    ruPath: "/pronunciation.html",
    enPath: "/en/pronunciation.html",
    ruTitle: "Тренажер произношения английского с AI - NERIVA",
    enTitle: "English Pronunciation Trainer With AI - NERIVA",
    ruH1: "Тренажер произношения английского с AI",
    enH1: "English pronunciation trainer with AI",
  },
  {
    slug: "photo-translation",
    ruPath: "/photo-translation.html",
    enPath: "/en/photo-translation.html",
    ruTitle: "Фото-перевод для изучения языков - NERIVA",
    enTitle: "Photo Translation for Language Learning - NERIVA",
    ruH1: "Фото-перевод меню, вывесок и заданий в NERIVA",
    enH1: "Photo translation for menus, signs, and tasks in NERIVA",
  },
  {
    slug: "telegram-language-bot",
    ruPath: "/telegram-language-bot.html",
    enPath: "/en/telegram-language-bot.html",
    ruTitle: "Telegram-бот для изучения языков - NERIVA",
    enTitle: "Telegram Language Bot - NERIVA",
    ruH1: "Telegram-бот NERIVA как быстрый вход в языковую практику",
    enH1: "NERIVA Telegram bot for fast language practice",
  },
  {
    slug: "ai-english-tutor",
    ruPath: "/ai-english-tutor.html",
    enPath: "/en/ai-english-tutor.html",
    ruTitle: "AI-репетитор английского онлайн - NERIVA",
    enTitle: "AI English Tutor Online - NERIVA",
    ruH1: "AI-репетитор английского для speaking, слов и ошибок",
    enH1: "AI English tutor for speaking, vocabulary, and mistakes",
  },
  {
    slug: "english-speaking-practice",
    ruPath: "/english-speaking-practice.html",
    enPath: "/en/english-speaking-practice.html",
    ruTitle: "Практика разговорного английского онлайн - NERIVA",
    enTitle: "English Speaking Practice Online - NERIVA",
    ruH1: "Практика разговорного английского онлайн без расписания",
    enH1: "English speaking practice online without scheduling",
  },
  {
    slug: "english-pronunciation-trainer",
    ruPath: "/english-pronunciation-trainer.html",
    enPath: "/en/english-pronunciation-trainer.html",
    ruTitle: "Тренажер произношения английских слов - NERIVA",
    enTitle: "English Pronunciation Trainer Online - NERIVA",
    ruH1: "Тренажер произношения английских слов и фраз",
    enH1: "English pronunciation trainer for words and phrases",
  },
  {
    slug: "language-learning-web-app",
    ruPath: "/language-learning-web-app.html",
    enPath: "/en/language-learning-web-app.html",
    ruTitle: "Web app для изучения языков - NERIVA",
    enTitle: "Language Learning Web App - NERIVA",
    ruH1: "Web app для изучения языков с AI и Telegram",
    enH1: "Language learning web app with AI and Telegram",
  },
] as const;

const rebrandedPublicPaths = [
  "/neriva.html?lang=ru",
  "/neriva.html?lang=en",
  "/privacy.html?lang=ru",
  "/terms.html?lang=ru",
  "/agreement.html?lang=ru",
  "/consent.html?lang=ru",
  "/ai-tutor.html",
  "/speaking-practice.html",
  "/pronunciation.html",
  "/photo-translation.html",
  "/telegram-language-bot.html",
  "/ai-english-tutor.html",
  "/english-speaking-practice.html",
  "/english-pronunciation-trainer.html",
  "/language-learning-web-app.html",
  "/en/ai-tutor.html",
  "/en/speaking-practice.html",
  "/en/pronunciation.html",
  "/en/photo-translation.html",
  "/en/telegram-language-bot.html",
  "/en/ai-english-tutor.html",
  "/en/english-speaking-practice.html",
  "/en/english-pronunciation-trainer.html",
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
        await expect(page.locator('a[href="https://t.me/NERIVAapp_bot"]').filter({ visible: true }).first()).toBeVisible();
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
      await expect(page.locator(".seo-related a")).toHaveCount(pages.length - 1);
      await expect(page.locator('a[data-entry="web-app"]').first()).toHaveAttribute("href", "/app/");
      await expect(page.locator('.seo-footer a[data-entry="landing"]')).toHaveAttribute("href", "/neriva.html");
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
      await expect(page.locator('.seo-footer a[data-entry="landing"]')).toHaveAttribute("href", "/neriva.html?lang=en");
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
    expect(body).not.toMatch(/hreflang="(?!ru|en|x-default)[a-z-]+"/);
    expect(body).not.toMatch(/https:\/\/neriva\.ru\/(?:es|de|fr|it|zh|ja|ko|pt|pl|uk)\//);
  });

  test("keeps the static SEO and AEO page cluster Russian and English only", async ({ request }) => {
    expect(pages).toHaveLength(9);

    for (const item of pages) {
      await expect.poll(async () => (await request.get(item.ruPath)).status(), { message: `${item.ruPath} should stay available` }).toBe(200);
      await expect.poll(async () => (await request.get(item.enPath)).status(), { message: `${item.enPath} should stay available` }).toBe(200);
      for (const locale of ["es", "de", "fr", "pt"] as const) {
        const response = await request.get(`/${locale}/${item.slug}.html`);
        expect(response.status(), `/${locale}/${item.slug}.html should not be generated for SEO/AEO`).toBe(404);
      }
    }
  });

  test("landing exposes SEO guides only outside the primary hero and navigation", async ({ page }) => {
    await page.goto("/neriva.html?lang=en");

    const guideLinks = page.locator(".seo-guides a");
    await expect(guideLinks).toHaveCount(pages.length);
    await expect(page.locator(".landing-hero .seo-guides")).toHaveCount(0);
    await expect(page.locator(".public-nav .seo-guides")).toHaveCount(0);
    await expect(guideLinks.first()).toHaveAttribute("href", "/en/ai-tutor.html");

    await page.goto("/neriva.html?lang=ru");
    const ruGuideLinks = page.locator(".seo-guides a");
    await expect(ruGuideLinks).toHaveCount(pages.length);
    await expect(ruGuideLinks.first()).toHaveAttribute("href", "/ai-tutor.html");
  });

  test("landing social links use downloaded color icon assets", async ({ page }) => {
    await page.goto("/neriva.html?lang=en");

    const expectedIcons = [
      "/assets/social/youtube-full-color.svg",
      "/assets/social/instagram-logo-2022.svg",
      "/assets/social/tiktok-icon.svg",
      "/assets/social/telegram-logo.png",
    ];

    for (const src of expectedIcons) {
      await expect(page.locator(`.social-icon-link img[src="${src}"]`).first()).toBeVisible();
      const response = await page.request.get(src);
      expect(response.ok(), `${src} should be a served downloaded social icon asset`).toBe(true);
    }
  });
});
