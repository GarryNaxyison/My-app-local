import { expect, test } from "@playwright/test";

test.describe("public landing SEO and AEO metadata", () => {
  test("sets Russian canonical metadata, alternates, JSON-LD, and visible answers", async ({ page }) => {
    await page.goto("/poliglot-ai.html?lang=ru");

    await expect(page).toHaveTitle("Poliglot AI - AI-репетитор английского и языков в Telegram");
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
    const supportedLanguages = await page.locator("[data-site-language-select]").first().evaluate((select) =>
      Array.from((select as HTMLSelectElement).options).map((option) => option.value),
    );
    expect([...new Set(alternates.map((link) => link.hreflang))].sort()).toEqual([...supportedLanguages, "x-default"].sort());
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
    expect(graphTypes).toHaveLength(4);
    expect(graphTypes).toEqual(
      expect.arrayContaining(["Organization", "WebSite", "SoftwareApplication", "FAQPage"]),
    );
    const faqNode = jsonLd["@graph"].find((node: { "@type": string }) => node["@type"] === "FAQPage");
    expect(faqNode.mainEntity).toHaveLength(6);
    expect(faqNode.mainEntity.map((entity: { name: string }) => entity.name)).toContain("Что такое Poliglot AI?");

    await expect(page.getByRole("heading", { name: "Ответы для поиска и AI" })).toBeVisible();
    await expect(page.locator(".answer-card")).toHaveCount(6);
    await expect(page.locator(".answer-card").filter({ hasText: "Что такое Poliglot AI?" })).toBeVisible();
  });

  test("sets English international metadata without forcing app links away from the current host", async ({ page }) => {
    await page.goto("/poliglot-ai.html?lang=en");

    await expect(page).toHaveTitle("Poliglot AI - AI language tutor in Telegram and web app");
    await expect(page.locator('meta[name="description"]')).toHaveAttribute(
      "content",
      /AI language tutor, speaking practice, Telegram bot, web app, voice, photo translation, mistakes, and premium plans/,
    );
    await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", "https://poliglotai.online/poliglot-ai.html?lang=en");
    await expect(page.locator('meta[property="og:url"]')).toHaveAttribute("content", "https://poliglotai.online/poliglot-ai.html?lang=en");
    await expect(page.locator('meta[name="twitter:title"]')).toHaveAttribute("content", "Poliglot AI - AI language tutor in Telegram and web app");

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
    await expect(page.locator(".answer-card").filter({ hasText: "AI language tutor" })).toBeVisible();

    const webEntryHrefs = await page.locator('a[data-entry="web-app"]').evaluateAll((links) =>
      links.map((link) => (link as HTMLAnchorElement).getAttribute("href")),
    );
    expect(webEntryHrefs.length).toBeGreaterThan(0);
    expect(webEntryHrefs.every((href) => href === "/app/")).toBe(true);
  });

  test("serves robots and sitemap search files", async ({ page }) => {
    const robotsResponse = await page.goto("/robots.txt");
    expect(robotsResponse).toBeTruthy();
    expect(robotsResponse?.ok()).toBeTruthy();
    const robotsBody = await robotsResponse?.text();
    expect(robotsBody).toContain("User-agent: *");
    expect(robotsBody).toContain("Sitemap: https://poliglotai.ru/sitemap.xml");
    expect(robotsBody).toContain("Sitemap: https://poliglotai.online/sitemap.xml");

    const sitemapResponse = await page.goto("/sitemap.xml");
    expect(sitemapResponse).toBeTruthy();
    expect(sitemapResponse?.ok()).toBeTruthy();
    const sitemapBody = await sitemapResponse?.text();
    expect(sitemapBody).toContain("https://poliglotai.ru/poliglot-ai.html");
    expect(sitemapBody).toContain("https://poliglotai.online/poliglot-ai.html?lang=en");
    expect(sitemapBody).toContain('hreflang="x-default"');
  });
});
