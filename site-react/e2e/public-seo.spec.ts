import { expect, test, type Page } from "@playwright/test";

const englishAnswerCopy = [
  "Search answers",
  "Answers for search and AI assistants",
  "Short direct answers about NERIVA for people comparing AI tutors, Telegram language bots, speaking practice, voice, and photo translation.",
  "What is NERIVA?",
  "NERIVA is an AI language tutor in a web app and Telegram bot. It combines short lessons, speaking practice, pronunciation, photo translation, mistakes, notes, and progress in one profile.",
  "Can I practice English with an AI tutor in Telegram?",
  "Yes. In Telegram you can start practice, receive tasks, send answers and voice messages, while progress stays synced with the web app.",
  "How is an AI tutor different from a vocabulary app?",
  "NERIVA is not only a word list. It gives a phrase in context, asks for your answer, corrects the mistake, and brings the weak spot back for review.",
  "Can I practice pronunciation and speaking?",
  "Yes. Voice Coach and shadowing help you practice speech, see weak words, get a score, and repeat a more natural phrase.",
  "Can I translate text from photos?",
  "Yes. A photo of a menu, sign, or exercise becomes a translation, note, and short practice prompt in that context.",
  "Is there a free plan?",
  "Yes. Free gives starter lessons and practice without payment, while Premium and Platinum unlock higher limits, voice, photo tools, and intensive daily study.",
] as const;

const englishComparisonCopy = [
  "Compare options",
  "NERIVA compared with the tools people usually search for",
  "See when a vocabulary app, AI tutor, Telegram bot, or web app is the better fit for language practice.",
  "NERIVA vs vocabulary app",
  "Word lists, flashcards, spaced repetition, and isolated meanings.",
  "Context phrases, answers, corrections, weak-spot review, voice, and photo practice.",
  "Best when you need correction and context, not only memorization.",
  "AI tutor vs language bot",
  "Quick chat prompts in Telegram with simple answers.",
  "Guided lessons, roleplay, mistake explanations, progress, and next repetition.",
  "Best when Telegram convenience needs tutor-level feedback.",
  "Telegram bot vs web app",
  "Fast tasks, voice messages, reminders, and short daily practice.",
  "Longer sessions, pricing, notes, progress, mistakes, and profile control.",
  "Best when quick mobile practice and deeper desktop study should stay synced.",
] as const;

const russianLandingTitle = "NERIVA - изучение языков и AI-репетитор английского в Telegram и приложении";
const russianLandingDescription =
  "NERIVA - изучение языков с AI-репетитором: английский и другие языки в Telegram и приложении NERIVA, уроки, разговорная практика, произношение, фото-перевод и ошибки.";
const russianLandingSocialDescription =
  "Изучение языков с AI-репетитором: английский и другие языки в Telegram и приложении NERIVA, уроки, speaking, произношение, фото-перевод и ошибки.";
const mojibakeCodepointPattern = /[\u0080-\u009f\u00b5\u0402\u0403\u0453\u0457\u201a\u2026]/;

function staticMetaContent(html: string, selector: "description" | "og:title" | "og:description" | "twitter:title" | "twitter:description") {
  if (selector === "description") {
    return html.match(/<meta\s+name="description"\s+content="([^"]+)"/i)?.[1] || "";
  }
  if (selector.startsWith("og:")) {
    return html.match(new RegExp(`<meta\\s+property="${selector}"\\s+content="([^"]+)"`, "i"))?.[1] || "";
  }
  return html.match(new RegExp(`<meta\\s+name="${selector}"\\s+content="([^"]+)"`, "i"))?.[1] || "";
}

async function selectLandingLanguage(page: Page, code: string) {
  const languageSelect = page.locator("[data-site-language-select]").first();
  await languageSelect.selectOption(code);
  await expect(languageSelect).toHaveValue(code);
  await expect.poll(async () => page.locator("html").getAttribute("lang")).toBe(code);
}

test.describe("public landing SEO and AEO metadata", () => {
  test.describe.configure({ timeout: 90_000 });

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

  test("serves valid UTF-8 Russian social metadata before client JavaScript", async ({ request }) => {
    const response = await request.get("/neriva.html");
    expect(response.ok()).toBeTruthy();
    expect(response.headers()["content-type"]).toContain("text/html");
    const html = await response.text();

    expect(html.match(/<title>([^<]+)<\/title>/i)?.[1]).toBe(russianLandingTitle);
    expect(staticMetaContent(html, "description")).toBe(russianLandingDescription);
    expect(staticMetaContent(html, "og:title")).toBe(russianLandingTitle);
    expect(staticMetaContent(html, "og:description")).toBe(russianLandingSocialDescription);
    expect(staticMetaContent(html, "twitter:title")).toBe(russianLandingTitle);
    expect(staticMetaContent(html, "twitter:description")).toBe(russianLandingSocialDescription);
    [
      russianLandingTitle,
      russianLandingDescription,
      russianLandingSocialDescription,
      staticMetaContent(html, "og:title"),
      staticMetaContent(html, "og:description"),
      staticMetaContent(html, "twitter:title"),
      staticMetaContent(html, "twitter:description"),
    ].forEach((value) => expect(value).not.toMatch(mojibakeCodepointPattern));
  });

  test("sets Russian canonical metadata, alternates, JSON-LD, and visible answers", async ({ page }) => {
    await page.goto("/neriva.html?lang=ru");

    await expect(page).toHaveTitle(russianLandingTitle);
    await expect(page.locator('meta[name="description"]')).toHaveAttribute(
      "content",
      russianLandingDescription,
    );
    await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", "https://neriva.ru/neriva.html");

    const alternates = await page.locator('link[rel="alternate"]').evaluateAll((links) =>
      links.map((link) => ({
        hreflang: link.getAttribute("hreflang"),
        href: link.getAttribute("href"),
      })),
    );
    const supportedLanguages = await page.locator("[data-site-language-select]").first().evaluate((select) =>
      Array.from((select as HTMLSelectElement).options).map((option) => option.value),
    );
    expect([...new Set(alternates.map((link) => link.hreflang))].sort()).toEqual([...supportedLanguages, "x-default"].sort());
    expect(alternates).toEqual(
      expect.arrayContaining([
        { hreflang: "ru", href: "https://neriva.ru/neriva.html" },
        { hreflang: "en", href: "https://neriva.ru/neriva.html?lang=en" },
        { hreflang: "x-default", href: "https://neriva.ru/neriva.html?lang=en" },
      ]),
    );

    const jsonLdText = await page.locator("#poliglot-seo-jsonld").textContent();
    expect(jsonLdText).toBeTruthy();
    const jsonLd = JSON.parse(jsonLdText || "{}");
    const graphTypes = jsonLd["@graph"].map((node: { "@type": string }) => node["@type"]);
    expect(graphTypes).toHaveLength(4);
    expect(graphTypes).toEqual(
      expect.arrayContaining(["Organization", "WebSite", "SoftwareApplication", "FAQPage"]),
    );
    const faqNode = jsonLd["@graph"].find((node: { "@type": string }) => node["@type"] === "FAQPage");
    expect(faqNode.mainEntity).toHaveLength(4);
    expect(faqNode.mainEntity.map((entity: { name: string }) => entity.name)).toContain("Что такое NERIVA?");
    expect(jsonLd["@graph"].some((node: { "@type": string }) => node["@type"] === "ItemList")).toBe(false);

    await expect(page.locator('meta[property="og:image"]')).toHaveAttribute("content", "https://neriva.ru/assets/seo/neriva-social-preview-v2-ru.jpg");
    await expect(page.locator('meta[property="og:image:width"]')).toHaveAttribute("content", "1200");
    await expect(page.locator('meta[property="og:image:height"]')).toHaveAttribute("content", "630");
    await expect(page.locator('meta[name="twitter:card"]')).toHaveAttribute("content", "summary_large_image");

    await expect(page.locator(".landing-faq h2")).toBeVisible();
    await expect(page.locator(".landing-faq__item")).toHaveCount(4);
    await expect(page.locator(".landing-faq__item").filter({ hasText: "Можно начать бесплатно?" })).toBeVisible();
    await expect(page.locator(".comparison-card")).toHaveCount(0);
    await expect(page.locator(".answer-card")).toHaveCount(0);
  });

  test("sets English international metadata without forcing app links away from the current host", async ({ page }) => {
    await page.goto("/neriva.html?lang=en");

    await expect(page).toHaveTitle("NERIVA - AI language tutor in Telegram and web app");
    await expect(page.locator('meta[name="description"]')).toHaveAttribute(
      "content",
      /AI language tutor, speaking practice, Telegram bot, web app, voice, photo translation, mistakes, and premium plans/,
    );
    await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", "https://neriva.ru/neriva.html?lang=en");
    await expect(page.locator('meta[property="og:url"]')).toHaveAttribute("content", "https://neriva.ru/neriva.html?lang=en");
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

    expect(jsonLd["@graph"].some((node: { "@type": string }) => node["@type"] === "ItemList")).toBe(false);

    await expect(page.locator('meta[property="og:image"]')).toHaveAttribute("content", "https://neriva.ru/assets/seo/neriva-social-preview-v2-en.jpg");
    await expect(page.locator('meta[property="og:image:width"]')).toHaveAttribute("content", "1200");
    await expect(page.locator('meta[property="og:image:height"]')).toHaveAttribute("content", "630");
    await expect(page.locator('meta[name="twitter:card"]')).toHaveAttribute("content", "summary_large_image");

    await expect(page.locator(".landing-faq h2")).toBeVisible();
    await expect(page.locator(".landing-faq__item")).toHaveCount(4);
    await expect(page.locator(".comparison-card")).toHaveCount(0);
    await expect(page.locator(".answer-card")).toHaveCount(0);

    const webEntryHrefs = await page.locator('a[data-entry="web-app"]').evaluateAll((links) =>
      links.map((link) => (link as HTMLAnchorElement).getAttribute("href")),
    );
    expect(webEntryHrefs.length).toBeGreaterThan(0);
    expect(webEntryHrefs.every((href) => href === "/app/")).toBe(true);
  });

  test("localizes visible AEO answers across landing languages", async ({ page }) => {
    await page.goto("/neriva.html?lang=en");

    const languageCodes = await page.locator("[data-site-language-select] option").evaluateAll((options) =>
      options.map((option) => (option as HTMLOptionElement).value),
    );
    expect(languageCodes).toEqual(["ru", "en"]);

    for (const code of languageCodes) {
      await selectLandingLanguage(page, code);

      await expect(page.locator(".landing-faq__item")).toHaveCount(4);

      if (code === "en") {
        continue;
      }

      const sectionText = await page.locator(".landing-faq").innerText();
      for (const englishCopy of englishAnswerCopy) {
        expect(sectionText).not.toContain(englishCopy);
      }
    }
  });

  test("keeps comparison answers off the main landing across landing languages", async ({ page }) => {
    await page.goto("/neriva.html?lang=en");

    const languageCodes = await page.locator("[data-site-language-select] option").evaluateAll((options) =>
      options.map((option) => (option as HTMLOptionElement).value),
    );
    expect(languageCodes).toEqual(["ru", "en"]);

    for (const code of languageCodes) {
      await selectLandingLanguage(page, code);

      await expect(page.locator(".comparison-card")).toHaveCount(0);
      await expect(page.locator(".comparison-section")).toHaveCount(0);
      for (const englishCopy of englishComparisonCopy) {
        await expect(page.locator(".english-spark-landing")).not.toContainText(englishCopy);
      }
    }
  });

  test("serves robots and sitemap search files", async ({ page }) => {
    const robotsResponse = await page.goto("/robots.txt");
    expect(robotsResponse).toBeTruthy();
    expect(robotsResponse?.ok()).toBeTruthy();
    const robotsBody = await robotsResponse?.text();
    expect(robotsBody).toContain("User-agent: *");
    expect(robotsBody).toContain("Sitemap: https://neriva.ru/sitemap.xml");

    const sitemapResponse = await page.goto("/sitemap.xml");
    expect(sitemapResponse).toBeTruthy();
    expect(sitemapResponse?.ok()).toBeTruthy();
    const sitemapBody = await sitemapResponse?.text();
    expect(sitemapBody).toContain("https://neriva.ru/neriva.html");
    expect(sitemapBody).toContain("https://neriva.ru/neriva.html?lang=en");
    expect(sitemapBody).not.toContain("https://neriva.ru/poliglot-ai.html");
    expect(sitemapBody).toContain('hreflang="x-default"');
  });
});
