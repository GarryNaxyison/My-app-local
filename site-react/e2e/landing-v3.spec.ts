import { expect, test, type Page } from "@playwright/test";

async function hideCookieBanner(page: Page) {
  await page.addInitScript(() => {
    localStorage.setItem(
      "poliglot-cookie-consent",
      JSON.stringify({ version: 1, necessary: true, analyticsMarketing: false }),
    );
    localStorage.setItem("poliglot-site-theme", "dark");
  });
}

test.beforeEach(async ({ page }) => {
  await page.route("https://mc.yandex.ru/**", (route) => route.abort());
  await hideCookieBanner(page);
});

test("landing v3 uses the premium dark reference and exposes only Russian and English", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator("html")).toHaveAttribute("data-site-theme", "dark");
  await expect(page.locator(".nav-theme-toggle")).toHaveCount(0);
  await expect(page.locator(".landing-hero")).toHaveAttribute("data-hero-preset", "dark");
  await expect(page.locator(".english-spark-landing")).toHaveAttribute("data-visual-anchor", "premium-tech-narrative");
  await expect(page.locator(".seo-guides")).toHaveCount(1);
  await expect(page.locator(".landing-hero .seo-guides")).toHaveCount(0);
  await expect(page.locator(".public-nav .seo-guides")).toHaveCount(0);

  const optionValues = await page.locator("[data-site-language-select] option").evaluateAll((options) =>
    options.map((option) => (option as HTMLOptionElement).value),
  );
  expect(optionValues).toEqual(["ru", "en"]);
});

test("landing v3 hero uses filled product proof and transparent dark animation", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator("h1")).toContainText("NERIVA");
  await expect(page.locator(".premium-hero-animation__shader")).toHaveCount(0);
  await expect(page.locator(".premium-hero-animation__three canvas")).toHaveCount(1);
  await expect(page.locator(".hero-proof")).toContainText("Browser + Telegram");
  await expect(page.locator(".hero-proof")).toContainText("Story");
  await expect(page.locator(".hero-proof")).not.toContainText("35 languages");

  const heroState = await page.evaluate(() => {
    const hero = document.querySelector(".landing-hero") as HTMLElement;
    const matter = document.querySelector(".landing-hero__matter") as HTMLElement;
    const productFrame = document.querySelector(".hero-product-frame") as HTMLElement;
    const proofCells = Array.from(document.querySelectorAll(".hero-proof span")) as HTMLElement[];
    const heroBox = hero.getBoundingClientRect();
    const matterBox = matter.getBoundingClientRect();
    const productFrameBox = productFrame.getBoundingClientRect();
    return {
      heroBackground: getComputedStyle(hero).backgroundColor,
      heroBackgroundImage: getComputedStyle(hero).backgroundImage,
      matterBackground: getComputedStyle(matter).backgroundColor,
      matterTopOffset: Math.round(matterBox.top - heroBox.top),
      productFrameTopOffset: Math.round(productFrameBox.top - heroBox.top),
      proofTexts: proofCells.map((cell) => (cell.textContent || "").trim()),
      emptyProofCells: proofCells.filter((cell) => !(cell.textContent || "").trim()).length,
    };
  });

  expect(heroState.heroBackground).toBe("rgba(0, 0, 0, 0)");
  expect(heroState.matterBackground).toBe("rgba(0, 0, 0, 0)");
  expect(heroState.matterTopOffset).toBeGreaterThanOrEqual(70);
  expect(heroState.productFrameTopOffset).toBeGreaterThanOrEqual(210);
  expect(heroState.emptyProofCells).toBe(0);

  const heroCanvasState = await page.locator(".premium-hero-animation__three canvas").evaluate((canvas) => {
    const element = canvas as HTMLCanvasElement;
    const box = element.getBoundingClientRect();
    return { width: element.width, height: element.height, boxWidth: box.width, boxHeight: box.height };
  });
  expect(heroCanvasState.width).toBeGreaterThan(100);
  expect(heroCanvasState.height).toBeGreaterThan(100);
  expect(heroCanvasState.boxWidth).toBeGreaterThan(100);
  expect(heroCanvasState.boxHeight).toBeGreaterThan(100);
});

test("landing v3 uses before-after bridge and a controlled feature carousel", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".stitch-metrics")).toBeVisible();
  await expect(page.locator(".stitch-metric-card")).toHaveCount(3);
  await expect(page.locator(".stitch-metrics")).toContainText("98%");
  await expect(page.locator(".stitch-metrics")).toContainText("A1-C2");

  await expect(page.locator(".before-after-bridge")).toBeVisible();
  await expect(page.locator(".before-after-bridge")).toContainText("Before");
  await expect(page.locator(".before-after-bridge")).toContainText("After");

  await expect(page.locator(".feature-carousel")).toBeVisible();
  await expect(page.locator(".feature-carousel__tab")).toHaveCount(4);
  await expect(page.locator(".feature-panel")).toHaveCount(1);
  await expect(page.locator(".feature-grid")).toHaveCount(0);
  await expect(page.locator(".feature-carousel")).toContainText("AI Tutor");
  await expect(page.locator(".feature-carousel")).toContainText("Voice");
  await expect(page.locator(".feature-carousel")).toContainText("Photo");
  await expect(page.locator(".feature-carousel")).toContainText("Mistakes");

  await expect(page.locator(".methodology-timeline")).toBeVisible();
  await expect(page.locator(".methodology-step")).toHaveCount(3);
  await expect(page.locator(".methodology-timeline")).toContainText("Discovery");
  await expect(page.locator(".methodology-timeline")).toContainText("Review");
});

test("Russian landing localizes metrics and methodology blocks", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=ru");

  await expect(page.locator(".stitch-metrics")).toContainText("Анализ речи");
  await expect(page.locator(".stitch-metrics")).toContainText("Глубина маршрута");
  await expect(page.locator(".stitch-metrics")).toContainText("Браузер и Telegram используют одну учебную память");
  await expect(page.locator(".stitch-metrics")).not.toContainText("Speech analysis");
  await expect(page.locator(".stitch-metrics")).not.toContainText("Browser and Telegram share");

  await expect(page.locator(".methodology-timeline")).toContainText("Методика");
  await expect(page.locator(".methodology-timeline")).toContainText("Диагностика");
  await expect(page.locator(".methodology-timeline")).toContainText("Правка");
  await expect(page.locator(".methodology-timeline")).toContainText("Повтор");
  await expect(page.locator(".methodology-timeline")).not.toContainText("Methodology");
  await expect(page.locator(".methodology-timeline")).not.toContainText("Discovery");
});

test("footer keeps contacts in the first column without duplicate Telegram links", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=ru");

  const firstColumn = page.locator(".site-footer > div").first();
  await expect(firstColumn.locator("address")).toBeVisible();
  await expect(firstColumn.locator("address")).toContainText("Контакты");
  await expect(firstColumn.locator('a[href="https://t.me/NERIVAapp_bot"]')).toHaveCount(1);
});

test("landing v3 pricing, community and final CTA are product-specific dark sections", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".plan-card")).toHaveCount(3);
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText("50 lessons");
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText("20 voice checks");
  await expect(page.locator(".pricing-section .product-shot--pricing")).toHaveCount(0);
  await expect(page.locator(".pricing-section")).toHaveCSS("background-color", "rgba(0, 0, 0, 0)");
  await expect(page.locator(".payment-methods")).toHaveCount(0);

  await expect(page.locator(".telegram-panel")).toContainText("Voice");
  await expect(page.locator(".telegram-panel")).toContainText("Photo");
  await expect(page.locator(".telegram-panel")).toContainText("Reminder");

  await expect(page.locator(".ecosystem-showcase")).toBeVisible();
  await expect(page.locator(".ecosystem-showcase .product-shot")).toHaveCount(3);
  await expect(page.locator(".ecosystem-showcase")).toContainText("Web app");
  await expect(page.locator(".ecosystem-showcase")).toContainText("Telegram");

  await expect(page.locator(".community-section")).toContainText("NERIVA channel");
  await expect(page.locator(".community-section")).not.toContainText("reviews");
  await expect(page.locator(".community-section")).not.toContainText("rating");

  await expect(page.locator(".final-cta-section")).toContainText("Phrase");
  await expect(page.locator(".final-cta-section")).toContainText("Answer");
  await expect(page.locator(".final-cta-section")).toContainText("Correction");
  await expect(page.locator(".final-cta-section")).toContainText("Repeat");
  await expect(page.locator(".final-cta-section")).toHaveCSS("background-color", "rgba(0, 0, 0, 0)");
});

test("landing v3 mobile is compact with no horizontal overflow", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".landing-hero .entry-cta")).toHaveCount(2);
  await expect(page.locator(".feature-carousel__viewport")).toBeVisible();
  await expect(page.locator(".mobile-carousel__tab")).toHaveCount(4);
  await expect(page.locator(".mobile-panel")).toHaveCount(1);

  const layout = await page.evaluate(() => {
    const heroHeight = document.querySelector(".landing-hero")?.getBoundingClientRect().height || 0;
    return {
      heroHeight,
      viewportHeight: window.innerHeight,
      overflowX: Math.max(0, document.documentElement.scrollWidth - document.documentElement.clientWidth),
    };
  });

  expect(layout.heroHeight).toBeLessThanOrEqual(layout.viewportHeight * 1.15);
  expect(layout.overflowX).toBe(0);
});

test("landing v3 product screenshots open a focused preview", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  await page.locator(".hero-product-frame .product-shot__button").first().click();
  await expect(page.locator(".image-preview")).toBeVisible();
  const previewImage = page.locator(".image-preview img");
  await expect(previewImage).toHaveAttribute("src", "/assets/product/dashboard-progress-dark.png");
  await expect(page.locator(".image-preview")).toHaveAttribute("data-preview-scale", "1.00");

  await page.locator(".image-preview__frame").hover();
  await page.mouse.wheel(0, -360);
  await expect(page.locator(".image-preview")).not.toHaveAttribute("data-preview-scale", "1.00");

  await page.locator(".image-preview__close").click();
  await expect(page.locator(".image-preview")).toHaveCount(0);
});
