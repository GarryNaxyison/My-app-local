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
  await page.goto("/neriva.html?lang=en");

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

test("landing v3 Russian copy avoids removed Telegram block and awkward AI wording", async ({ page }) => {
  await page.goto("/neriva.html?lang=ru");

  await expect(page.locator("body")).not.toContainText("Быстрая практика без второго продукта");
  await expect(page.locator("body")).not.toContainText("Запустите задание, отправьте голос или фото");
  await expect(page.locator("body")).not.toContainText("чанки");
  await expect(page.locator("body")).not.toContainText("production task");
  await expect(page.locator("body")).not.toContainText("AI-ответ");
  await expect(page.locator("body")).not.toContainText("speaking-попытку");
  await expect(page.locator("body")).not.toContainText("streak");
  await expect(page.locator("body")).toContainText("Короткий урок с понятным результатом");
  await expect(page.locator("body")).toContainText("серию занятий");
});

test("landing v3 hero uses filled product proof and transparent dark animation", async ({ page }) => {
  await page.goto("/neriva.html?lang=en");

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
  expect(heroState.matterTopOffset).toBeLessThanOrEqual(24);
  expect(heroState.productFrameTopOffset).toBeGreaterThanOrEqual(160);
  expect(heroState.emptyProofCells).toBe(0);

  const premiumSurfaceState = await page.locator(".english-spark-landing").evaluate((landing) => {
    const surface = landing as HTMLElement;
    const card = document.querySelector(".hero-proof > span") as HTMLElement;
    return {
      backgroundImage: getComputedStyle(surface).backgroundImage,
      cardAnimation: getComputedStyle(card).animationName,
      cardBackground: getComputedStyle(card).backgroundImage,
      cardBackgroundSize: getComputedStyle(card).backgroundSize,
    };
  });
  expect(premiumSurfaceState.backgroundImage).toContain("radial-gradient");
  expect(premiumSurfaceState.backgroundImage).not.toContain("rgba(255, 255, 255, 0.025) 1px");
  expect(premiumSurfaceState.cardAnimation).toContain("nerivaBorderOrbit");
  expect(premiumSurfaceState.cardBackground).toContain("linear-gradient");
  expect(premiumSurfaceState.cardBackgroundSize).toContain("260%");

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
  await page.goto("/neriva.html?lang=en");

  await expect(page.locator(".stitch-metrics")).toBeVisible();
  await expect(page.locator(".stitch-metric-card")).toHaveCount(3);
  await expect(page.locator(".stitch-metrics")).toContainText("score");
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

test("landing v3 blocks reveal as the visitor scrolls", async ({ page }) => {
  await page.goto("/neriva.html?lang=en");

  const featurePanel = page.locator(".feature-panel");
  const initialState = await featurePanel.evaluate((element) => {
    const style = getComputedStyle(element);
    return {
      top: Math.round(element.getBoundingClientRect().top),
      opacity: style.opacity,
      transform: style.transform,
      filter: style.filter,
    };
  });
  expect(initialState.top).toBeGreaterThan(900);
  expect(Number(initialState.opacity)).toBeLessThan(0.05);
  expect(initialState.transform).not.toBe("none");
  expect(initialState.filter).toContain("blur");

  await featurePanel.scrollIntoViewIfNeeded();

  await expect.poll(async () =>
    featurePanel.evaluate((element) => Number(getComputedStyle(element).opacity)),
  ).toBeGreaterThan(0.95);
  await expect.poll(async () =>
    featurePanel.evaluate((element) => {
      const transform = getComputedStyle(element).transform;
      return Math.abs(transform === "none" ? 0 : Number(transform.match(/matrix\([^,]+,[^,]+,[^,]+,[^,]+,[^,]+,\s*([^)]+)\)/)?.[1] ?? 0));
    }),
  ).toBeLessThan(0.5);
  await expect.poll(async () =>
    featurePanel.evaluate((element) => {
      const filter = getComputedStyle(element).filter;
      return Number(filter.match(/^blur\(([^)]+)px\)$/)?.[1] ?? 0);
    }),
  ).toBeLessThan(0.05);

  const revealedState = await featurePanel.evaluate((element) => {
    const style = getComputedStyle(element);
    return {
      opacity: style.opacity,
      transform: style.transform,
      filter: style.filter,
    };
  });
  expect(Number(revealedState.opacity)).toBeGreaterThan(0.95);
  const translateY = revealedState.transform === "none" ? 0 : Number(revealedState.transform.match(/matrix\([^,]+,[^,]+,[^,]+,[^,]+,[^,]+,\s*([^)]+)\)/)?.[1] ?? 0);
  expect(Math.abs(translateY)).toBeLessThan(0.5);
  const blurPx = Number(revealedState.filter.match(/^blur\(([^)]+)px\)$/)?.[1] ?? 0);
  expect(blurPx).toBeLessThan(0.05);
});

test("Russian landing localizes metrics and methodology blocks", async ({ page }) => {
  await page.goto("/neriva.html?lang=ru");

  await expect(page.locator(".stitch-metrics")).toContainText("Проверка голоса");
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
  await page.goto("/neriva.html?lang=ru");

  const firstColumn = page.locator(".site-footer > div").first();
  await expect(firstColumn.locator("address")).toBeVisible();
  await expect(firstColumn.locator("address")).toContainText("Контакты");
  await expect(firstColumn.locator('a[href="https://t.me/NERIVAapp_bot"]')).toHaveCount(1);
});

test("landing v3 pricing, community and final CTA are product-specific dark sections", async ({ page }) => {
  await page.goto("/neriva.html?lang=en");

  await expect(page.locator(".plan-card")).toHaveCount(3);
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText("50 lessons");
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText("20 voice checks");
  await expect(page.locator(".pricing-section .product-shot--pricing")).toHaveCount(0);
  await expect(page.locator(".pricing-section")).toHaveCSS("background-color", "rgba(0, 0, 0, 0)");
  await expect(page.locator(".payment-methods")).toHaveCount(0);

  await expect(page.locator(".telegram-panel")).toHaveCount(0);

  await expect(page.locator(".ecosystem-showcase")).toBeVisible();
  await expect(page.locator(".ecosystem-showcase .product-shot")).toHaveCount(3);
  await expect(page.locator(".ecosystem-showcase")).toContainText("Web app");
  await expect(page.locator(".ecosystem-showcase")).toContainText("Telegram");

  await expect(page.locator(".community-section")).toContainText("NERIVA channel");
  await expect(page.locator(".community-section")).not.toContainText("reviews");
  await expect(page.locator(".community-section")).not.toContainText("rating");
  await expect(page.locator('.community-section a[href="https://t.me/neriva_app"]')).toBeVisible();

  await expect(page.locator(".final-cta-section")).toContainText("Phrase");
  await expect(page.locator(".final-cta-section")).toContainText("Answer");
  await expect(page.locator(".final-cta-section")).toContainText("Correction");
  await expect(page.locator(".final-cta-section")).toContainText("Repeat");
  await expect(page.locator(".final-cta-section")).toHaveCSS("background-color", "rgba(0, 0, 0, 0)");
});

test("landing v3 mobile is compact with no horizontal overflow", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/neriva.html?lang=en");

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

test("landing v3 mobile hero and methodology do not look clipped", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/neriva.html?lang=ru");

  const layout = await page.evaluate(() => {
    const hero = document.querySelector(".landing-hero") as HTMLElement;
    const step = document.querySelector(".methodology-step") as HTMLElement;
    const label = step.querySelector(":scope > span") as HTMLElement;
    const dot = getComputedStyle(step, "::before");
    const heroBox = hero.getBoundingClientRect();
    const stepBox = step.getBoundingClientRect();
    const labelBox = label.getBoundingClientRect();
    const dotLeft = stepBox.left + Number.parseFloat(dot.left);
    const dotRight = dotLeft + Number.parseFloat(dot.width);

    return {
      heroLeft: Math.round(heroBox.left),
      heroRight: Math.round(heroBox.right),
      viewportWidth: window.innerWidth,
      labelLeft: Math.round(labelBox.left),
      dotRight: Math.round(dotRight),
    };
  });

  expect(layout.heroLeft).toBe(0);
  expect(layout.heroRight).toBe(layout.viewportWidth);
  expect(layout.labelLeft).toBeGreaterThanOrEqual(layout.dotRight + 12);
});

test("landing v3 product screenshots open a focused preview", async ({ page }) => {
  await page.goto("/neriva.html?lang=en");

  await page.locator(".hero-product-frame .product-shot__button").first().click();
  await expect(page.locator(".image-preview")).toBeVisible();
  const previewImage = page.locator(".image-preview img");
  await expect(previewImage).toHaveAttribute("src", "/assets/product/dashboard-progress-dark.png");
  await expect(page.locator(".image-preview")).toHaveAttribute("data-preview-scale", "1.00");

  await page.locator(".image-preview__frame").hover();
  await page.mouse.wheel(0, -360);
  await expect(page.locator(".image-preview")).not.toHaveAttribute("data-preview-scale", "1.00");
  await expect(page.locator(".image-preview")).toHaveAttribute("data-preview-pan", "0,0");

  const frameOverflow = await page.locator(".image-preview__frame").evaluate((frame) => {
    const style = getComputedStyle(frame);
    return {
      overflowX: style.overflowX,
      overflowY: style.overflowY,
    };
  });
  expect(frameOverflow).toEqual({ overflowX: "hidden", overflowY: "hidden" });

  const frameBox = await page.locator(".image-preview__frame").boundingBox();
  expect(frameBox).not.toBeNull();
  await page.mouse.move(frameBox!.x + frameBox!.width / 2, frameBox!.y + frameBox!.height / 2);
  await page.mouse.down();
  await page.mouse.move(frameBox!.x + frameBox!.width / 2 + 90, frameBox!.y + frameBox!.height / 2 + 42);
  await page.mouse.up();
  await expect(page.locator(".image-preview")).not.toHaveAttribute("data-preview-pan", "0,0");

  await page.locator(".image-preview__close").click();
  await expect(page.locator(".image-preview")).toHaveCount(0);
});
