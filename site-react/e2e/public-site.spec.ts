import { expect, test } from "@playwright/test";

const siteLocaleCodes = [
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
];

test("landing keeps goal-course product flow with Poliglot features", async ({ page }) => {
  test.setTimeout(60_000);
  await page.goto("/poliglot-ai.html");

  await expect(page.locator(".public-nav")).toBeVisible();
  await expect(page.locator(".public-nav a", { hasText: /^v2$/i })).toHaveCount(0);
  await expect(page.locator(".public-brand__logo img")).toBeVisible();
  const appLinks = await page.locator('a[href^="/app"]').evaluateAll((links) => [...new Set(links.map((link) => (link as HTMLAnchorElement).getAttribute("href")))].sort());
  expect(appLinks).toEqual(["/app/"]);
  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  const lightHero = await page.locator(".landing-hero").evaluate((hero) => {
    document.documentElement.dataset.siteTheme = "light";
    const heroBox = hero.getBoundingClientRect();
    const matter = hero.querySelector(".landing-hero__matter") as HTMLElement | null;
    const matterBox = matter?.getBoundingClientRect();
    const h1 = hero.querySelector("h1") as HTMLElement;
    const firstGoal = hero.querySelector(".hero-goals a strong") as HTMLElement | null;
    return {
      heroWidth: heroBox.width,
      heroHeight: heroBox.height,
      matterWidth: matterBox?.width || 0,
      matterHeight: matterBox?.height || 0,
      h1Color: getComputedStyle(h1).color,
      firstGoalColor: firstGoal ? getComputedStyle(firstGoal).color : "",
    };
  });
  expect(lightHero.matterWidth).toBeGreaterThan(lightHero.heroWidth * 0.9);
  expect(lightHero.matterHeight).toBeGreaterThan(lightHero.heroHeight * 0.9);
  expect(lightHero.h1Color).not.toBe("rgb(7, 17, 31)");
  expect(lightHero.firstGoalColor).not.toBe("rgb(7, 17, 31)");
  const lightEyebrow = await page.locator(".workflow-section .eyebrow").evaluate((node) => {
    const style = getComputedStyle(node);
    return {
      color: style.color,
      textFill: style.getPropertyValue("-webkit-text-fill-color"),
      background: style.backgroundImage || style.backgroundColor,
    };
  });
  expect(lightEyebrow.color).toBe("rgb(15, 23, 42)");
  expect(lightEyebrow.textFill).toBe("rgb(15, 23, 42)");
  expect(lightEyebrow.background).not.toBe("rgba(243, 184, 75, 0.1)");

  await expect(page.locator("h1")).toContainText("Выберите цель");
  const heroGoals = page.locator(".hero-goals a strong");
  await expect(heroGoals).toHaveCount(4);
  await expect(heroGoals).toHaveText(["Путешествия", "Работа", "Экзамен", "Разговорная речь"]);
  await expect(page.locator(".hero-language-picker a")).toHaveCount(8);
  await expect(page.locator(".course-card")).toHaveCount(3);
  await expect(page.locator(".feature-card")).toHaveCount(6);
  await expect(page.locator(".review-card")).toHaveCount(3);
  await expect(page.locator(".faq-grid article")).toHaveCount(4);
  await expect(page.locator(".hero-proof")).toContainText("35");
  await expect(page.locator(".hero-proof")).toContainText("A1-C2");
  await expect(page.locator(".hero-proof")).toContainText("Free");

  await expect(page.locator(".hero-demo__tabs button")).toHaveText(["Урок", "Диалог", "Голос", "Фото"]);
  await page.locator(".hero-demo__tabs button").nth(1).click();
  await expect(page.locator(".demo-output p")).toContainText("Could you help me check in?");
  await page.locator(".hero-demo__tabs button").nth(3).click();
  await expect(page.locator(".demo-output p")).toContainText("No peanuts");
  await expect(page.locator(".demo-wave i")).toHaveCount(22);

  const oldPrices = page.locator(".plan-old-price");
  await expect(oldPrices).toHaveCount(2);
  await expect(oldPrices.first()).toHaveCSS("text-decoration-line", /line-through/);
  const oldPriceColor = await oldPrices.first().evaluate((node) => getComputedStyle(node).color);
  expect(oldPriceColor).not.toBe("rgb(255, 255, 255)");
  expect(oldPriceColor).not.toBe("rgba(255, 255, 255, 0.5)");
  await expect(page.locator(".payment-methods")).toContainText("Stars");
  await expect(page.locator(".payment-methods")).toContainText("YooKassa");
  await expect(page.locator(".payment-methods")).toContainText("TON");
  await expect(page.locator(".payment-methods")).toContainText("USDT");
});

test("landing hero keeps animated motion within a bounded frame budget", async ({ page }) => {
  await page.addInitScript(() => {
    const originalRequestAnimationFrame = window.requestAnimationFrame.bind(window);
    let frameRequests = 0;
    window.requestAnimationFrame = ((callback: FrameRequestCallback) => {
      frameRequests += 1;
      return originalRequestAnimationFrame(callback);
    }) as typeof window.requestAnimationFrame;
    Object.defineProperty(window, "__poliglotFrameRequests", {
      get: () => frameRequests,
    });
  });

  await page.goto("/poliglot-ai.html");
  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  await page.waitForTimeout(1800);

  const frameRequests = await page.evaluate(() => (window as Window & { __poliglotFrameRequests?: number }).__poliglotFrameRequests ?? 0);
  expect(frameRequests).toBeGreaterThan(15);
  expect(frameRequests).toBeLessThan(95);
});

test("public site language selector exposes all interface locales", async ({ page }) => {
  test.setTimeout(120_000);
  await page.goto("/poliglot-ai.html");
  const select = page.locator("[data-site-language-select]");
  await expect(select).toBeVisible();

  const optionValues = await select.locator("option").evaluateAll((options) => options.map((option) => (option as HTMLOptionElement).value));
  expect(optionValues).toEqual(siteLocaleCodes);

  for (const code of ["en", "ar", "bn", "cs", "hi", "ta", "te", "th", "tr", "vi"]) {
    await page.evaluate((nextCode) => {
      const languageSelect = document.querySelector("[data-site-language-select]") as HTMLSelectElement | null;
      if (!languageSelect) throw new Error("missing public site language selector");
      languageSelect.value = nextCode;
      languageSelect.dispatchEvent(new Event("change", { bubbles: true }));
    }, code);
    await page.waitForTimeout(80);
    const text = await page.locator("body").innerText();
    expect(text).not.toMatch(/Рџ|РЎ|Р’|Рќ|Рњ/);
    if (code !== "ru") {
      expect(text).not.toContain("Начать обучение");
      expect(text).not.toContain("Курсы на выбор");
    }
  }
});

test("privacy and terms keep contacts and dark language picker readable", async ({ page }) => {
  for (const path of ["/privacy.html", "/terms.html"]) {
    await page.goto(path);
    await page.evaluate(() => {
      document.documentElement.dataset.siteTheme = "dark";
    });

    const select = page.locator("[data-site-language-select]");
    await expect(select).toBeVisible();
    const selectColors = await select.evaluate((node) => {
      const style = getComputedStyle(node);
      const option = (node as HTMLSelectElement).options[0];
      const optionStyle = option ? getComputedStyle(option) : style;
      return {
        background: style.backgroundColor,
        color: style.color,
        optionBackground: optionStyle.backgroundColor,
        optionColor: optionStyle.color,
      };
    });
    expect(selectColors.background).not.toBe("rgb(255, 255, 255)");
    expect(selectColors.optionBackground).not.toBe("rgb(255, 255, 255)");
    expect(selectColors.color).not.toBe(selectColors.background);
    expect(selectColors.optionColor).not.toBe(selectColors.optionBackground);

    const text = await page.locator("body").innerText();
    expect(text).toContain("@AsaselD");
    expect(text).toContain("@poliglot_ai_bot");
    expect(text).toContain("35");
  }
});

test("privacy Russian contact grid keeps Telegram bot visible", async ({ page }) => {
  await page.goto("/privacy.html?lang=ru");

  const botContact = page.locator('#contacts .legal-contact-grid a[href="https://t.me/poliglot_ai_bot"]');
  await expect(botContact).toBeVisible();
  await expect(botContact.locator("small")).toHaveText("Telegram bot");
  await expect(botContact).toContainText("@poliglot_ai_bot");
});
