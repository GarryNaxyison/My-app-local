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

const ruCopy = {
  privacy: "\u041f\u043e\u043b\u0438\u0442\u0438\u043a\u0430",
  terms: "\u0423\u0441\u043b\u043e\u0432\u0438\u044f",
  aiTutor: "AI-\u0440\u0435\u043f\u0435\u0442\u0438\u0442\u043e\u0440\u043e\u043c",
  lesson: "\u0423\u0440\u043e\u043a",
  dialogue: "\u0414\u0438\u0430\u043b\u043e\u0433",
  voice: "\u0413\u043e\u043b\u043e\u0441",
  photo: "\u0424\u043e\u0442\u043e",
  fastPractice: "\u0431\u044b\u0441\u0442\u0440\u0430\u044f \u043f\u0440\u0430\u043a\u0442\u0438\u043a\u0430",
  freePrice: "0 \u20bd",
  premiumPrice: "300 \u20bd",
  platinumPrice: "590 \u20bd",
  oldStart: "\u041d\u0430\u0447\u0430\u0442\u044c \u043e\u0431\u0443\u0447\u0435\u043d\u0438\u0435",
  oldCourses: "\u041a\u0443\u0440\u0441\u044b \u043d\u0430 \u0432\u044b\u0431\u043e\u0440",
} as const;

test("landing presents the bold product cockpit without losing public contracts", async ({ page }) => {
  test.setTimeout(60_000);
  await page.goto("/poliglot-ai.html");

  await expect(page.locator(".public-nav")).toBeVisible();
  await expect(page.locator(".public-brand__logo img")).toBeVisible();
  await expect(page.locator(".public-nav nav a", { hasText: ruCopy.privacy })).toBeVisible();
  await expect(page.locator(".public-nav nav a", { hasText: ruCopy.terms })).toBeVisible();
  await expect(page.locator("[data-site-language-select]")).toBeVisible();

  await expect(page.locator(".bold-landing")).toBeVisible();
  await expect(page.locator(".bold-hero")).toBeVisible();
  await expect(page.locator(".bold-hero__matter canvas")).toHaveCount(1);
  const boldCanvasLocations = await page.locator(".bold-landing canvas").evaluateAll((canvases) =>
    canvases.map((canvas) => Boolean(canvas.closest(".bold-hero"))),
  );
  expect(boldCanvasLocations).toEqual([true]);
  await expect(page.locator(".spark-hero")).toHaveCount(0);

  const bodyFontFamily = await page.locator("body").evaluate((node) => getComputedStyle(node).fontFamily);
  expect(bodyFontFamily).not.toContain("Inter");
  expect(bodyFontFamily).not.toContain("Inrer");
  const forbiddenLandingFonts = await page.locator(".bold-landing, .bold-landing *").evaluateAll((nodes) =>
    nodes
      .map((node) => getComputedStyle(node as HTMLElement).fontFamily)
      .filter((fontFamily) => /Inter|Inrer/.test(fontFamily)),
  );
  expect(forbiddenLandingFonts).toEqual([]);

  await expect(page.locator('a[data-entry="web-app"]')).toHaveCount(3);
  const webEntryHrefs = await page.locator('a[data-entry="web-app"]').evaluateAll((links) =>
    links.map((link) => (link as HTMLAnchorElement).getAttribute("href")),
  );
  expect(webEntryHrefs).toEqual(["/app/", "/app/", "/app/"]);
  await expect(page.locator('a[data-entry="telegram"]')).toHaveCount(3);
  const telegramEntryHrefs = await page.locator('a[data-entry="telegram"]').evaluateAll((links) =>
    links.map((link) => (link as HTMLAnchorElement).getAttribute("href")),
  );
  expect(telegramEntryHrefs).toEqual([
    "https://t.me/poliglot_ai_bot",
    "https://t.me/poliglot_ai_bot",
    "https://t.me/poliglot_ai_bot",
  ]);

  const heroCtas = page.locator(".bold-hero .entry-cta");
  await expect(heroCtas).toHaveCount(2);
  const ctaShape = await page.locator(".bold-hero .entry-cta").evaluateAll((nodes) =>
    nodes.map((node) => {
      const element = node as HTMLElement;
      const style = getComputedStyle(element);
      const rect = element.getBoundingClientRect();
      return {
        minHeight: rect.height,
        width: Math.round(rect.width),
        borderRadius: style.borderRadius,
        background: style.backgroundImage || style.backgroundColor,
        color: style.color,
      };
    }),
  );
  expect(ctaShape).toHaveLength(2);
  expect(ctaShape[0].minHeight).toBeGreaterThanOrEqual(48);
  expect(ctaShape[1].minHeight).toBeGreaterThanOrEqual(48);
  expect(Math.abs(ctaShape[0].width - ctaShape[1].width)).toBeLessThanOrEqual(28);
  expect(ctaShape[0].background).not.toBe("none");
  expect(ctaShape[1].background).not.toBe("none");
  expect(ctaShape[0].borderRadius).toBe(ctaShape[1].borderRadius);

  await expect(page.locator("h1")).toContainText(ruCopy.aiTutor);
  await expect(page.locator(".bold-hero__proof")).toContainText("35");
  await expect(page.locator(".bold-hero__proof")).toContainText("A1-C2");

  await expect(page.locator(".hero-product-tabs button")).toHaveText([ruCopy.lesson, ruCopy.dialogue, ruCopy.voice, ruCopy.photo]);
  await page.locator(".hero-product-tabs button").nth(1).click();
  await expect(page.locator(".hero-product-output")).toContainText("Could you help me check in?");
  await page.locator(".hero-product-tabs button").nth(2).click();
  await expect(page.locator(".voice-bars i")).toHaveCount(18);
  await page.locator(".hero-product-tabs button").nth(3).click();
  await expect(page.locator(".hero-product-output")).toContainText("No peanuts");

  await expect(page.locator(".daily-step")).toHaveCount(3);
  await expect(page.locator(".module-card")).toHaveCount(4);
  await expect(page.locator(".module-card", { hasText: "Voice Coach" })).toContainText("weak words");
  await expect(page.locator(".module-card", { hasText: "Photo Practice" })).toContainText("No peanuts");
  await expect(page.locator(".memory-node")).toHaveCount(6);

  await expect(page.locator(".entry-panel")).toHaveCount(2);
  await expect(page.locator(".entry-panel", { hasText: "Web app" })).toContainText("dashboard");
  await expect(page.locator(".entry-panel", { hasText: "Telegram" })).toContainText(ruCopy.fastPractice);

  await expect(page.locator(".plan-card")).toHaveCount(3);
  await expect(page.locator(".plan-card", { hasText: "Free" })).toContainText(ruCopy.freePrice);
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText(ruCopy.premiumPrice);
  await expect(page.locator(".plan-card", { hasText: "Platinum" })).toContainText(ruCopy.platinumPrice);
  await expect(page.locator(".payment-methods")).toContainText("Stars");
  await expect(page.locator(".payment-methods")).toContainText("YooKassa");
  await expect(page.locator(".payment-methods")).toContainText("TON");
  await expect(page.locator(".payment-methods")).toContainText("USDT");

  await expect(page.locator(".review-card")).toHaveCount(3);
  await expect(page.locator(".review-stars")).toHaveCount(0);

  await expect(page.locator(".final-cta-section")).toContainText("Web app");
  await expect(page.locator(".final-cta-section")).toContainText("Telegram");

  const footerPrivacyLink = page.locator(".site-footer a", { hasText: ruCopy.privacy });
  const footerTermsLink = page.locator(".site-footer a", { hasText: ruCopy.terms });
  await expect(footerPrivacyLink).toBeVisible();
  await expect(footerPrivacyLink).toHaveAttribute("href", /privacy\.html(\?.*)?$/);
  await expect(footerTermsLink).toBeVisible();
  await expect(footerTermsLink).toHaveAttribute("href", /terms\.html(\?.*)?$/);
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
  await expect(page.locator(".bold-hero__matter canvas")).toHaveCount(1);
  await page.waitForTimeout(1800);

  const frameRequests = await page.evaluate(() => (window as Window & { __poliglotFrameRequests?: number }).__poliglotFrameRequests ?? 0);
  expect(frameRequests).toBeGreaterThan(15);
  expect(frameRequests).toBeLessThan(95);
});

test("bold product cockpit landing stays readable on mobile", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/poliglot-ai.html");

  await expect(page.locator(".bold-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator("h1")).toContainText(ruCopy.aiTutor);
  await expect(page.locator(".bold-hero .entry-cta")).toHaveCount(2);
  await expect(page.locator(".hero-product-mockup")).toBeVisible();
  await expect(page.locator(".entry-panel")).toHaveCount(2);
  await expect(page.locator(".site-footer a", { hasText: ruCopy.privacy })).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: ruCopy.terms })).toBeVisible();

  const overflow = await page.evaluate(() => {
    const documentOverflow = document.documentElement.scrollWidth > document.documentElement.clientWidth + 2 || document.body.scrollWidth > document.body.clientWidth + 2;
    const nodes = Array.from(
      document.querySelectorAll(
        ".public-nav, [data-site-language-select], h1, h2, h3, p, a, button, .bold-hero__proof span, .daily-step, .module-card, .memory-node, .entry-panel, .plan-card, .review-card",
      ),
    );
    const elementOverflow = nodes
      .map((node) => {
        const element = node as HTMLElement;
        const rect = element.getBoundingClientRect();
        return {
          text: (element.textContent || "").trim().slice(0, 80),
          scrollWidth: element.scrollWidth,
          clientWidth: element.clientWidth,
          width: rect.width,
        };
      })
      .filter((item) => item.width > 0 && item.scrollWidth > item.clientWidth + 2);
    return documentOverflow ? [{ text: "document", scrollWidth: document.documentElement.scrollWidth, clientWidth: document.documentElement.clientWidth, width: window.innerWidth }, ...elementOverflow] : elementOverflow;
  });
  expect(overflow).toEqual([]);
});
test("premium cockpit landing localizes generated marketing copy for Chinese", async ({ page }) => {
  test.setTimeout(120_000);
  await page.goto("/poliglot-ai.html?lang=zh");
  await page.waitForTimeout(1400);

  const sectionsWithCyrillic = await page.evaluate(() => {
    const selectors = [
      ".public-brand",
      ".public-nav nav",
      ".public-nav__actions a",
      ".bold-hero",
      ".daily-loop-section",
      ".modules-section",
      ".memory-loop-section",
      ".entry-section",
      ".pricing-section",
      ".reviews-section",
      ".final-cta-section",
      ".site-footer",
    ];
    return selectors
      .map((selector) => {
        const node = document.querySelector(selector);
        const text = (node?.textContent || "").replace(/\s+/g, " ").trim();
        const cyrillic = text.match(/\p{Script=Cyrillic}+/gu) || [];
        return { selector, cyrillic: cyrillic.filter((fragment) => !/^[AMS]$/.test(fragment)).slice(0, 5) };
      })
      .filter((item) => item.cyrillic.length > 0);
  });

  expect(sectionsWithCyrillic).toEqual([]);
  await expect(page.locator(".pricing-section")).toContainText("50");
  await expect(page.locator(".pricing-section")).toContainText("200");
  await expect(page.locator(".pricing-section")).toContainText("20");
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
      expect(text).not.toContain(ruCopy.oldStart);
      expect(text).not.toContain(ruCopy.oldCourses);
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
