import { expect, test } from "@playwright/test";

const firstArticlePath = "/knowledge/can-you-learn-a-language-yourself.html";
const memoryArticlePath = "/knowledge/how-language-memory-works.html";
const examArticlePath = "/knowledge/how-to-prepare-for-language-exam.html";
const lastArticlePath = "/knowledge/how-to-choose-language-learning-app.html";
const vocabularyClusterPath = "/knowledge/vocabulary.html";

test.describe("NERIVA knowledge base", () => {
  test("renders the knowledge hub with Stitch-style search, filters, and 100 articles", async ({ page }) => {
    const response = await page.goto("/knowledge/");
    expect(response?.ok()).toBe(true);

    await expect(page).toHaveTitle(/База знаний NERIVA/);
    await expect(page.locator("html")).toHaveAttribute("lang", "ru");
    await expect(page.locator(".knowledge-shell")).toBeVisible();
    await expect(page.locator(".knowledge-nav__link--knowledge")).toHaveText("База знаний");
    await expect(page.locator(".knowledge-nav__link--knowledge")).toHaveCSS("border-bottom-color", "rgb(183, 196, 255)");
    await expect(page.locator(".knowledge-search__input")).toHaveAttribute("placeholder", "Найти статью, вопрос или тему");
    await expect(page.locator(".knowledge-filter")).toHaveCount(10);
    await expect(page.locator(".knowledge-card")).toHaveCount(100);
    await expect(page.locator(".knowledge-card").first()).toContainText("Можно ли выучить иностранный язык самостоятельно?");
    await expect(page.locator(".knowledge-rail")).toContainText("Маршрут для старта");
    await expect(page.locator(".knowledge-plan")).toContainText("500 материалов");

    const languageSelectBox = await page.locator(".knowledge-language").boundingBox();
    const knowledgeLinkBox = await page.locator(".knowledge-nav__link--knowledge").boundingBox();
    expect(knowledgeLinkBox).toBeTruthy();
    expect(languageSelectBox).toBeTruthy();
    expect((knowledgeLinkBox?.x || 0) + (knowledgeLinkBox?.width || 0)).toBeLessThan(languageSelectBox?.x || 0);
  });

  test("filters article cards by search query and restores the full list", async ({ page }) => {
    await page.goto("/knowledge/");

    await page.locator(".knowledge-search__input").fill("слова");
    const visibleCards = page.locator(".knowledge-card:not([hidden])");
    await expect(page.locator(".knowledge-results-count")).toContainText('20 материалов по запросу "слова"');
    await expect(visibleCards).toHaveCount(20);
    await expect(visibleCards.first()).toContainText("Почему иностранные слова быстро забываются?");
    await expect(page.locator(".knowledge-card:not([hidden])", { hasText: "Как выбрать язык для изучения?" })).toHaveCount(0);

    await page.locator(".knowledge-filter", { hasText: "Словарный запас" }).click();
    await expect(page.locator(".knowledge-filter.is-active")).toContainText("Словарный запас");
    await expect(visibleCards).toHaveCount(15);

    await page.locator(".knowledge-reset").click();
    await expect(page.locator(".knowledge-results-count")).toContainText("100 материалов");
    await expect(visibleCards).toHaveCount(100);
  });

  test("search supports partial queries and query parameters", async ({ page }) => {
    await page.goto("/knowledge/");
    const visibleCards = page.locator(".knowledge-card:not([hidden])");

    await page.locator(".knowledge-search__input").fill("пам");
    expect(await visibleCards.count()).toBeGreaterThan(0);
    await expect(page.locator(".knowledge-card:not([hidden])", { hasText: "Как работает память при изучении языков?" })).toHaveCount(1);

    await page.goto("/knowledge/?q=экзам");
    await expect(page.locator(".knowledge-search__input")).toHaveValue("экзам");
    await expect(page.locator(".knowledge-results-count")).toContainText('материалов по запросу "экзам"');
    await expect(visibleCards.first()).toContainText("Как подготовиться к языковому экзамену?");
    expect(await visibleCards.count()).toBeGreaterThanOrEqual(5);
  });

  test("article mobile search opens the filtered knowledge hub", async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 844 });
    await page.goto(memoryArticlePath);

    await page.locator(".knowledge-mobile-header .knowledge-search__input").fill("экзам");
    await page.locator(".knowledge-mobile-header .knowledge-search__input").press("Enter");

    await expect(page).toHaveURL(/\/knowledge\/\?q=/);
    await expect(page.locator(".knowledge-search__input")).toHaveValue("экзам");
    await expect(page.locator(".knowledge-card:not([hidden])").first()).toContainText("Как подготовиться к языковому экзамену?");
  });

  test("filters the expanded article clusters by category", async ({ page }) => {
    await page.goto("/knowledge/");
    const visibleCards = page.locator(".knowledge-card:not([hidden])");

    await page.locator(".knowledge-filter", { hasText: "Словарный запас" }).click();
    await expect(visibleCards).toHaveCount(15);
    await expect(visibleCards.last()).toContainText("Какие слова учить в первую очередь?");

    await page.locator(".knowledge-filter", { hasText: "CEFR уровни" }).click();
    await expect(visibleCards).toHaveCount(8);
    await expect(visibleCards.first()).toContainText("Что такое уровень A1?");

    await page.locator(".knowledge-filter", { hasText: "Экзамены" }).click();
    await expect(visibleCards).toHaveCount(5);
    await expect(visibleCards.first()).toContainText("Как подготовиться к языковому экзамену?");

    await page.locator(".knowledge-filter", { hasText: "AI и технологии" }).click();
    await expect(visibleCards).toHaveCount(5);
    await expect(visibleCards.last()).toContainText("Как выбрать приложение для изучения языка?");
  });

  test("renders a detailed article with TOC, FAQ, related materials, and JSON-LD", async ({ page }) => {
    const response = await page.goto(firstArticlePath);
    expect(response?.ok()).toBe(true);

    await expect(page).toHaveTitle("Можно ли выучить иностранный язык самостоятельно? - База знаний NERIVA");
    await expect(page.locator("html")).toHaveAttribute("lang", "ru");
    await expect(page.locator(".knowledge-article h1")).toHaveText("Можно ли выучить иностранный язык самостоятельно?");
    await expect(page.locator(".article-answer")).toContainText("Короткий ответ");
    await expect(page.locator(".article-toc")).toContainText("Содержание");
    await expect(page.locator("#why-it-works")).toContainText("Почему это работает именно так");
    await expect(page.locator("#practical-plan")).toContainText("Практический план");
    await expect(page.locator("#common-mistakes")).toContainText("Типичные ошибки");
    await expect(page.locator("#how-neriva-helps")).toContainText("Как NERIVA может помочь");
    await expect(page.locator(".article-faq__item")).toHaveCount(5);
    await expect(page.locator(".article-related__card")).toHaveCount(3);
    await expect(page.locator(".article-cta")).toHaveClass(/article-cta--start/);
    await expect(page.locator(".article-cta a")).toHaveAttribute("href", "/app/");

    const pageText = await page.locator("body").innerText();
    expect(pageText).not.toMatch(/Quantum|Dr\. Elena|Solutions|Capabilities|Stats/);

    const jsonLdText = await page.locator('script[type="application/ld+json"]').textContent();
    const jsonLd = JSON.parse(jsonLdText || "{}");
    const types = jsonLd["@graph"].map((node: { "@type": string }) => node["@type"]);
    expect(types).toEqual(expect.arrayContaining(["Organization", "WebSite", "Article", "BreadcrumbList", "FAQPage"]));
    const articleNode = jsonLd["@graph"].find((node: { "@type": string }) => node["@type"] === "Article");
    expect(articleNode.headline).toBe("Можно ли выучить иностранный язык самостоятельно?");
    const faqNode = jsonLd["@graph"].find((node: { "@type": string }) => node["@type"] === "FAQPage");
    expect(faqNode.mainEntity).toHaveLength(5);
  });

  test("renders representative new articles from the 75-article expansion", async ({ page }) => {
    await page.goto(memoryArticlePath);
    await expect(page.locator(".knowledge-article h1")).toHaveText("Как работает память при изучении языков?");
    await expect(page.locator(".article-answer")).toContainText("Короткий ответ");
    await expect(page.locator(".article-faq__item")).toHaveCount(5);

    await page.goto(lastArticlePath);
    await expect(page.locator(".knowledge-article h1")).toHaveText("Как выбрать приложение для изучения языка?");
    await expect(page.locator(".article-answer")).toContainText("приложение");
    await expect(page.locator(".article-related__card")).toHaveCount(3);
  });

  test("keeps article layout usable on mobile without horizontal overflow", async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 1200 });
    await page.goto(firstArticlePath);

    await expect(page.locator(".knowledge-mobile-header")).toBeVisible();
    await expect(page.locator(".knowledge-mobile-header .knowledge-search__input")).toBeVisible();
    await expect(page.locator(".article-toc-mobile")).toBeVisible();
    await expect(page.locator(".article-toc")).toBeHidden();
    await expect(page.locator(".article-related__row")).toHaveCount(3);
    await expect(page.locator(".knowledge-footer")).toBeVisible();

    const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth);
    expect(overflow).toBeLessThanOrEqual(1);
  });

  test("adds knowledge base URLs to sitemap", async ({ page }) => {
    const response = await page.goto("/sitemap.xml");
    expect(response?.ok()).toBe(true);
    const body = await response?.text();
    expect(body).toContain("https://neriva.ru/knowledge/");
    expect(body).toContain("https://neriva.ru/knowledge/start.html");
    expect(body).toContain("https://neriva.ru/knowledge/vocabulary.html");
    expect(body).toContain("https://neriva.ru/knowledge/english.html");
    expect(body).toContain("https://neriva.ru/knowledge/speaking.html");
    expect(body).toContain("https://neriva.ru/knowledge/exam.html");
    expect(body).toContain("https://neriva.ru/knowledge/technology.html");
    expect(body).toContain(`https://neriva.ru${firstArticlePath}`);
    expect(body).toContain("https://neriva.ru/knowledge/what-is-spaced-repetition.html");
    expect(body).toContain(`https://neriva.ru${memoryArticlePath}`);
    expect(body).toContain(`https://neriva.ru${examArticlePath}`);
    expect(body).toContain(`https://neriva.ru${lastArticlePath}`);
  });

  test("renders category pillar pages as indexable article routes", async ({ page }) => {
    const response = await page.goto(vocabularyClusterPath);
    expect(response?.ok()).toBe(true);

    await expect(page.locator("body")).toHaveAttribute("data-page", "knowledge-cluster");
    await expect(page.locator(".cluster-route__step")).toHaveCount(3);
    await expect(page.locator(".knowledge-card")).toHaveCount(15);
    await expect(page.locator(`.knowledge-card[href="${memoryArticlePath}"]`)).toHaveCount(1);
    await expect(page.locator(`.knowledge-card[href="${firstArticlePath}"]`)).toHaveCount(0);
    await expect(page.locator(".cluster-cta a")).toHaveAttribute("href", "/app/");

    const jsonLdText = await page.locator('script[type="application/ld+json"]').textContent();
    const jsonLd = JSON.parse(jsonLdText || "{}");
    const types = jsonLd["@graph"].map((node: { "@type": string }) => node["@type"]);
    expect(types).toEqual(expect.arrayContaining(["CollectionPage", "ItemList", "BreadcrumbList"]));
  });

  test("adds in-body internal links and topic-specific CTA to articles", async ({ page }) => {
    await page.goto(memoryArticlePath);

    await expect(page.locator(".article-route-links")).toBeVisible();
    await expect(page.locator(".article-route-links a[href='/knowledge/vocabulary.html']")).toHaveCount(1);
    await expect(page.locator(".article-route-links a[href^='/knowledge/']")).toHaveCount(4);
    await expect(page.locator(".article-cta")).toHaveClass(/article-cta--vocabulary/);
    await expect(page.locator(".article-cta a")).toHaveAttribute("href", "/app/");

    await page.setViewportSize({ width: 390, height: 1200 });
    const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth);
    expect(overflow).toBeLessThanOrEqual(1);
  });
});
