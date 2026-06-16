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
  oldStart: "\u041d\u0430\u0447\u0430\u0442\u044c \u043e\u0431\u0443\u0447\u0435\u043d\u0438\u0435",
  oldCourses: "\u041a\u0443\u0440\u0441\u044b \u043d\u0430 \u0432\u044b\u0431\u043e\u0440",
} as const;

test("landing presents the approved English spark hero product site", async ({ page }) => {
  test.setTimeout(60_000);
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".public-nav")).toBeVisible();
  await expect(page.locator(".public-brand__logo img")).toBeVisible();
  await expect(page.locator(".nav-burger")).toBeVisible();
  await expect(page.locator("[data-site-language-select]")).toBeVisible();

  const navPosition = await page.locator(".public-nav").evaluate((node) => getComputedStyle(node).position);
  expect(["fixed", "sticky"]).not.toContain(navPosition);

  await expect(page.locator(".english-spark-landing")).toBeVisible();
  await expect(page.locator(".landing-hero")).toBeVisible();
  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator(".spark-hero")).toHaveCount(0);
  await expect(page.locator(".bold-hero")).toHaveCount(0);
  const heroHeight = await page.locator(".landing-hero").evaluate((node) => node.getBoundingClientRect().height);
  const viewportHeight = await page.evaluate(() => window.innerHeight);
  expect(heroHeight).toBeGreaterThanOrEqual(viewportHeight * 0.9);

  await expect(page.locator("h1")).toContainText("Practice speaking before the moment matters");
  await expect(page.locator(".hero-proof")).toContainText("35 languages");
  await expect(page.locator(".hero-proof")).toContainText("A1-C2");

  const heroCtas = page.locator(".landing-hero .entry-cta");
  await expect(heroCtas).toHaveCount(2);
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

  await expect(page.locator(".skeleton-card")).toHaveCount(2);
  await expect(page.locator(".skeleton-line")).toHaveCount(6);
  await expect(page.locator(".hero-demo")).toBeVisible();
  await expect(page.locator(".hero-demo__tabs button")).toHaveText(["Lesson", "Dialogue", "Voice", "Photo"]);
  await page.locator(".hero-demo__tabs button").nth(1).click();
  await expect(page.locator(".hero-demo__panel")).toContainText("Could you help me check in?");
  await page.locator(".hero-demo__tabs button").nth(2).click();
  await expect(page.locator(".demo-wave i")).toHaveCount(22);
  await page.locator(".hero-demo__tabs button").nth(3).click();
  await expect(page.locator(".hero-demo__panel")).toContainText("No peanuts");

  await expect(page.locator(".daily-step")).toHaveCount(3);
  await expect(page.locator(".module-card")).toHaveCount(4);
  await expect(page.locator(".module-card", { hasText: "Voice Coach" })).toContainText("weak words");
  await expect(page.locator(".module-card", { hasText: "Photo Practice" })).toContainText("menu or sign");
  await expect(page.locator(".memory-node")).toHaveCount(6);

  await expect(page.locator(".entry-panel")).toHaveCount(2);
  await expect(page.locator(".entry-panel", { hasText: "Web app" })).toContainText("deep sessions");
  await expect(page.locator(".entry-panel", { hasText: "Telegram" })).toContainText("quick practice");

  await expect(page.locator(".plan-card")).toHaveCount(3);
  await expect(page.locator(".plan-card", { hasText: "Free" })).toContainText("0 ₽");
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText("300 ₽");
  await expect(page.locator(".plan-card", { hasText: "Platinum" })).toContainText("590 ₽");
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText("Daily speaking plan");
  await expect(page.locator(".plan-card", { hasText: "Platinum" })).toContainText("Intensive preparation");
  await expect(page.locator(".plan-card s, .plan-card del")).toHaveCount(2);
  const oldPriceDecorations = await page.locator(".plan-card s, .plan-card del").evaluateAll((nodes) =>
    nodes.map((node) => {
      const style = getComputedStyle(node);
      return { text: node.textContent?.trim(), decoration: style.textDecorationLine };
    }),
  );
  expect(oldPriceDecorations).toEqual([
    { text: "1000 ₽", decoration: expect.stringContaining("line-through") },
    { text: "2000 ₽", decoration: expect.stringContaining("line-through") },
  ]);
  await expect(page.locator(".payment-methods")).toContainText("Stars");
  await expect(page.locator(".payment-methods")).toContainText("YooKassa");
  await expect(page.locator(".payment-methods")).toContainText("TON");
  await expect(page.locator(".payment-methods")).toContainText("USDT");

  await expect(page.locator(".review-card")).toHaveCount(3);
  await expect(page.locator(".review-card img")).toHaveCount(3);
  const reviewImages = await page.locator(".review-card img").evaluateAll((imgs) =>
    imgs.map((img) => ({ src: (img as HTMLImageElement).getAttribute("src"), alt: (img as HTMLImageElement).getAttribute("alt") })),
  );
  expect(reviewImages).toEqual([
    { src: "/assets/testimonials/anna.jpg", alt: "Anna" },
    { src: "/assets/testimonials/marat.jpg", alt: "Marat" },
    { src: "/assets/testimonials/sofia.jpg", alt: "Sofia" },
  ]);
  await expect(page.locator(".review-stars")).toHaveCount(0);

  await expect(page.locator(".final-cta-section")).toContainText("Web app");
  await expect(page.locator(".final-cta-section")).toContainText("Telegram");

  await page.locator(".nav-burger").click();
  await expect(page.locator(".nav-drawer")).toBeVisible();
  await expect(page.locator(".nav-drawer")).toContainText("Features");
  await expect(page.locator(".nav-drawer")).toContainText("Pricing");
  await expect(page.locator(".nav-drawer")).toContainText("Reviews");
  await expect(page.locator(".nav-drawer")).toContainText("Privacy");
  await expect(page.locator(".nav-drawer")).toContainText("Terms");
  await expect(page.locator(".nav-drawer")).toContainText("Web app");
  await expect(page.locator(".nav-drawer")).toContainText("Telegram");

  const landingCyrillic = await page.locator(".english-spark-landing").evaluate((node) => (node.textContent || "").match(/\p{Script=Cyrillic}+/gu) || []);
  expect(landingCyrillic).toEqual([]);

  const footerPrivacyLink = page.locator(".site-footer a", { hasText: "Privacy" });
  const footerTermsLink = page.locator(".site-footer a", { hasText: "Terms" });
  await expect(footerPrivacyLink).toBeVisible();
  await expect(footerPrivacyLink).toHaveAttribute("href", /privacy\.html(\?.*)?$/);
  await expect(footerTermsLink).toBeVisible();
  await expect(footerTermsLink).toHaveAttribute("href", /terms\.html(\?.*)?$/);
});

test("landing light theme keeps the approved layout readable", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");
  await page.evaluate(() => {
    document.documentElement.dataset.siteTheme = "light";
    localStorage.setItem("poliglot-site-theme", "light");
  });

  await expect(page.locator(".landing-hero")).toBeVisible();
  await expect(page.locator(".hero-demo")).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: "Privacy" })).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: "Terms" })).toBeVisible();

  const issues = await page.evaluate(() => {
    const selectors = [
      ".public-nav",
      ".landing-hero h1",
      ".landing-hero__copy p",
      ".hero-demo",
      ".daily-step",
      ".module-card",
      ".entry-panel",
      ".plan-card",
      ".review-card",
      ".site-footer",
    ];
    const contrastIssues = selectors
      .map((selector) => {
        const node = document.querySelector(selector) as HTMLElement | null;
        if (!node) return null;
        const style = getComputedStyle(node);
        return { selector, color: style.color, background: style.backgroundColor };
      })
      .filter(Boolean)
      .filter((item) => item!.color === item!.background);
    const overflowX = Math.max(0, document.documentElement.scrollWidth - document.documentElement.clientWidth);
    const overflowing = Array.from(document.querySelectorAll("h1,h2,h3,p,a,button,.plan-card,.review-card,.module-card,.entry-panel"))
      .map((node) => {
        const element = node as HTMLElement;
        return { text: (element.textContent || "").trim().slice(0, 80), sw: element.scrollWidth, cw: element.clientWidth, width: element.getBoundingClientRect().width };
      })
      .filter((item) => item.width > 0 && item.sw > item.cw + 2);
    return { contrastIssues, overflowX, overflowing };
  });

  expect(issues.contrastIssues).toEqual([]);
  expect(issues.overflowX).toBe(0);
  expect(issues.overflowing).toEqual([]);
});

test("landing keeps rich product sections below the restored hero", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".landing-hero")).toBeVisible();
  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator(".course-strip")).toBeVisible();
  await expect(page.locator(".course-card")).toHaveCount(3);
  await expect(page.locator(".course-card", { hasText: "AI Tutor Core" })).toContainText("Start route");
  await expect(page.locator(".course-stats")).toContainText("35");

  await expect(page.locator(".cockpit-section")).toBeVisible();
  await expect(page.locator(".cockpit-lane")).toHaveCount(5);
  await expect(page.locator(".cockpit-console")).toContainText("Live learning loop");
  await expect(page.locator(".cockpit-console")).toContainText("weak word fixed");

  await expect(page.locator(".scenario-section")).toBeVisible();
  await expect(page.locator(".scenario-card")).toHaveCount(4);
  await expect(page.locator(".scenario-card img")).toHaveCount(4);
  const scenarioImages = await page.locator(".scenario-card img").evaluateAll((imgs) =>
    imgs.map((img) => ({ src: (img as HTMLImageElement).getAttribute("src"), alt: (img as HTMLImageElement).getAttribute("alt") })),
  );
  expect(scenarioImages).toEqual([
    { src: "/assets/scenarios/travel-ai-tutor.jpg", alt: "Travel practice scene" },
    { src: "/assets/scenarios/work-ai-tutor.jpg", alt: "Work practice scene" },
    { src: "/assets/scenarios/exam-ai-tutor.jpg", alt: "Exam practice scene" },
    { src: "/assets/scenarios/speaking-ai-tutor.jpg", alt: "Speaking practice scene" },
  ]);

  await expect(page.locator(".device-flow")).toBeVisible();
  await expect(page.locator(".device-flow__node")).toHaveCount(3);
  await expect(page.locator(".device-flow__node", { hasText: "Web app" })).toContainText("longer sessions");
  await expect(page.locator(".device-flow__node", { hasText: "Telegram" })).toContainText("same profile");
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

  await page.goto("/poliglot-ai.html?lang=en");
  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  await page.waitForTimeout(1800);

  const frameRequests = await page.evaluate(() => (window as Window & { __poliglotFrameRequests?: number }).__poliglotFrameRequests ?? 0);
  expect(frameRequests).toBeGreaterThan(15);
  expect(frameRequests).toBeLessThan(95);
});

test("English spark landing stays readable on mobile", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator("h1")).toContainText("Practice speaking before the moment matters");
  await expect(page.locator(".landing-hero .entry-cta")).toHaveCount(2);
  await expect(page.locator(".hero-proof")).toBeVisible();
  await expect(page.locator(".entry-panel")).toHaveCount(2);
  await expect(page.locator(".site-footer a", { hasText: "Privacy" })).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: "Terms" })).toBeVisible();

  const heroHeight = await page.locator(".landing-hero").evaluate((node) => node.getBoundingClientRect().height);
  const viewportHeight = await page.evaluate(() => window.innerHeight);
  expect(heroHeight).toBeGreaterThanOrEqual(viewportHeight * 0.9);

  const overflow = await page.evaluate(() => {
    const documentOverflow = document.documentElement.scrollWidth > document.documentElement.clientWidth + 2 || document.body.scrollWidth > document.body.clientWidth + 2;
    const nodes = Array.from(
      document.querySelectorAll(
        ".public-nav, .nav-drawer, [data-site-language-select], h1, h2, h3, p, a, button, .hero-proof span, .daily-step, .module-card, .memory-node, .entry-panel, .plan-card, .review-card",
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

test("public site language selector exposes all interface locales", async ({ page }) => {
  test.setTimeout(120_000);
  await page.goto("/poliglot-ai.html");
  const select = page.locator("[data-site-language-select]");
  await expect(select).toBeVisible();

  const optionValues = await select.locator("option").evaluateAll((options) => options.map((option) => (option as HTMLOptionElement).value));
  expect(optionValues).toEqual(siteLocaleCodes);

  for (const code of siteLocaleCodes.filter((locale) => locale !== "ru")) {
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
      const landingText = await page.locator(".english-spark-landing").innerText();
      expect(landingText).not.toMatch(/\p{Script=Cyrillic}/u);
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
