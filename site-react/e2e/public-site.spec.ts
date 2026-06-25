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

const cyrillicInterfaceLocales = new Set(["ru", "tg", "tt", "kk", "ky", "uk"]);
const englishLandingCopyThatMustLocalize = [
  "Travel without panic",
  "Short hotel, cafe, airport, and doctor phrases are practiced before the moment gets stressful.",
  "Work calls and messages",
  "Emails, calls, self-intros, and deadline questions become rehearsals that make speaking easier.",
  "Exam and level progress",
  "A1-C2 vocabulary, grammar, listening, and speaking stay in a route with visible progress steps.",
  "Everyday conversation",
  "The tutor gives a natural version and a next repetition until your answer sounds more confident.",
  "Learn the phrase",
  "Get the meaning, grammar hint, natural example, and one focused prompt.",
  "Use it in context",
  "Practice through a short roleplay, voice answer, or photo-based task.",
  "Review the weak spot",
  "Mistakes, notes, weak words, XP, and streak point to the next repetition.",
  "Travel practice",
  "I rehearsed check-in, cafe orders, and transport before the trip. The phrases stayed in notes for quick review.",
  "English for work",
  "I use the web app for longer lessons and Telegram for weak words before calls. The same profile keeps it simple.",
  "Pronunciation",
  "Voice practice shows which words sound weak and gives a better sentence to repeat right away.",
] as const;

const curatedRussianLandingCopy = [
  "Поездка без паники",
  "Рабочие звонки и переписка",
  "Изучите фразу",
  "Используйте её в контексте",
  "Разберите слабое место",
  "Мобильная версия",
  "Ежедневный цикл",
  "Практика перед поездкой",
  "Открыть веб-приложение",
] as const;

const awkwardRussianLandingCopy = [
  "Мобильный Интернет",
  "План выступлений на день",
  "AI Репетитор",
  "более длительную AI репетиторскую работу",
  "Более высокие лимиты на поездки",
] as const;

test("landing presents the approved English spark hero product site", async ({ page }) => {
  test.setTimeout(120_000);
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator("html")).toHaveAttribute("data-site-theme", "dark");
  await expect(page.locator(".public-nav")).toBeVisible();
  await expect(page.locator(".public-brand__logo img")).toBeVisible();
  await expect(page.locator(".nav-burger")).toBeVisible();
  await expect(page.locator("[data-site-language-select]")).toBeVisible();

  const navPosition = await page.locator(".public-nav").evaluate((node) => getComputedStyle(node).position);
  expect(["fixed", "sticky"]).not.toContain(navPosition);

  await expect(page.locator(".english-spark-landing")).toBeVisible();
  await expect(page.locator(".landing-hero")).toBeVisible();
  await expect(page.locator(".landing-hero")).toHaveAttribute("data-hero-preset", "dark");
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
  await page.locator(".review-card").last().scrollIntoViewIfNeeded();
  await page.waitForFunction(() =>
    Array.from(document.querySelectorAll<HTMLImageElement>(".review-card img")).every((image) => image.complete && image.naturalWidth > 0),
  );
  const reviewImages = await page.locator(".review-card img").evaluateAll((imgs) =>
    imgs.map((img) => ({
      src: (img as HTMLImageElement).getAttribute("src"),
      alt: (img as HTMLImageElement).getAttribute("alt"),
      naturalWidth: (img as HTMLImageElement).naturalWidth,
      complete: (img as HTMLImageElement).complete,
    })),
  );
  expect(reviewImages).toEqual([
    { src: "/assets/testimonials/anna.jpg", alt: "Anna", naturalWidth: expect.any(Number), complete: true },
    { src: "/assets/testimonials/marat.jpg", alt: "Marat", naturalWidth: expect.any(Number), complete: true },
    { src: "/assets/testimonials/sofia.jpg", alt: "Sofia", naturalWidth: expect.any(Number), complete: true },
  ]);
  expect(reviewImages.every((img) => img.naturalWidth >= 240)).toBe(true);
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
  await page.addInitScript(() => {
    localStorage.setItem("poliglot-site-theme", "light");
  });
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator("html")).toHaveAttribute("data-site-theme", "light");
  await expect(page.locator(".landing-hero")).toBeVisible();
  await expect(page.locator(".landing-hero")).toHaveAttribute("data-hero-preset", "light");
  await expect(page.locator(".hero-demo")).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: "Privacy" })).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: "Terms" })).toBeVisible();

  const lightHeroHeader = await page.evaluate(() => {
    const parseColor = (value: string) => {
      const match = value.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)(?:,\s*([0-9.]+))?\)/);
      if (!match) return null;
      return {
        red: Number(match[1]),
        green: Number(match[2]),
        blue: Number(match[3]),
        alpha: match[4] ? Number(match[4]) : 1,
      };
    };

    const nav = document.querySelector(".public-nav--drawer") as HTMLElement;
    const hero = document.querySelector(".landing-hero") as HTMLElement;
    const matter = document.querySelector(".landing-hero__matter") as HTMLElement;
    const heading = document.querySelector(".landing-hero h1") as HTMLElement;
    const navBackground = parseColor(getComputedStyle(nav).backgroundColor);
    const matterStyle = getComputedStyle(matter);
    return {
      navBackground,
      heroColor: getComputedStyle(hero).color,
      heroBackgroundImage: getComputedStyle(hero).backgroundImage,
      headingColor: getComputedStyle(heading).color,
      matterOpacity: Number(matterStyle.opacity),
      matterBlendMode: matterStyle.mixBlendMode,
      matterFilter: matterStyle.filter,
    };
  });

  expect(lightHeroHeader.navBackground).toEqual(expect.objectContaining({ red: expect.any(Number), green: expect.any(Number), blue: expect.any(Number), alpha: expect.any(Number) }));
  expect(Math.max(lightHeroHeader.navBackground!.red, lightHeroHeader.navBackground!.green, lightHeroHeader.navBackground!.blue)).toBeLessThan(250);
  expect(lightHeroHeader.navBackground!.alpha).toBeLessThanOrEqual(0.78);
  expect(lightHeroHeader.heroColor).toBe("rgb(7, 17, 31)");
  expect(lightHeroHeader.headingColor).toBe("rgb(7, 17, 31)");
  expect(lightHeroHeader.heroBackgroundImage).not.toContain("rgb(5, 9, 20)");
  expect(lightHeroHeader.matterOpacity).toBeGreaterThanOrEqual(0.62);
  expect(lightHeroHeader.matterBlendMode).not.toBe("multiply");
  expect(lightHeroHeader.matterFilter).not.toContain("contrast(0.");
  await expect(page.locator(".pricing-section .eyebrow")).toHaveCSS("color", "rgb(15, 23, 42)");
  await expect(page.locator(".reviews-section .eyebrow")).toHaveCSS("color", "rgb(15, 23, 42)");
  await expect(page.locator(".plan-label").first()).toHaveCSS("color", "rgb(45, 91, 255)");
  await expect(page.locator(".course-card span").first()).toHaveCSS("color", "rgb(45, 91, 255)");
  await expect(page.locator(".scenario-modules span").first()).toHaveCSS("color", "rgb(45, 91, 255)");

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
  for (const image of await page.locator(".scenario-card img").all()) {
    await image.scrollIntoViewIfNeeded();
  }
  await page.waitForFunction(() =>
    Array.from(document.querySelectorAll<HTMLImageElement>(".scenario-card img")).every((image) => image.complete && image.naturalWidth > 0),
  );
  const scenarioImages = await page.locator(".scenario-card img").evaluateAll((imgs) =>
    imgs.map((img) => ({
      src: (img as HTMLImageElement).getAttribute("src"),
      alt: (img as HTMLImageElement).getAttribute("alt"),
      naturalWidth: (img as HTMLImageElement).naturalWidth,
      complete: (img as HTMLImageElement).complete,
    })),
  );
  expect(scenarioImages).toEqual([
    { src: "/assets/scenarios/travel-ai-tutor.jpg", alt: "Travel practice scene", naturalWidth: expect.any(Number), complete: true },
    { src: "/assets/scenarios/work-ai-tutor.jpg", alt: "Work practice scene", naturalWidth: expect.any(Number), complete: true },
    { src: "/assets/scenarios/exam-ai-tutor.jpg", alt: "Exam practice scene", naturalWidth: expect.any(Number), complete: true },
    { src: "/assets/scenarios/speaking-ai-tutor.jpg", alt: "Speaking practice scene", naturalWidth: expect.any(Number), complete: true },
  ]);
  expect(scenarioImages.every((img) => img.naturalWidth >= 480)).toBe(true);

  await expect(page.locator(".device-flow")).toBeVisible();
  await expect(page.locator(".device-flow__node")).toHaveCount(3);
  await expect(page.locator(".device-flow__node", { hasText: "Web app" })).toContainText("longer sessions");
  await expect(page.locator(".device-flow__node", { hasText: "Telegram" })).toContainText("same profile");
});

test("Russian landing uses edited copy and keeps CTAs and connectors aligned", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=ru");
  await page.evaluate(() => {
    document.documentElement.dataset.siteTheme = "light";
    localStorage.setItem("poliglot-site-theme", "light");
  });
  await expect(page.locator(".english-spark-landing")).toBeVisible();

  const landingText = await page.locator(".english-spark-landing").innerText();
  for (const expected of curatedRussianLandingCopy) {
    expect(landingText).toContain(expected);
  }
  for (const forbidden of awkwardRussianLandingCopy) {
    expect(landingText).not.toContain(forbidden);
  }

  const alignment = await page.evaluate(() => {
    const centerDelta = (card: Element, cta: Element) => {
      const cardBox = card.getBoundingClientRect();
      const ctaBox = cta.getBoundingClientRect();
      return Math.abs((cardBox.left + cardBox.right) / 2 - (ctaBox.left + ctaBox.right) / 2);
    };

    const planMetrics = Array.from(document.querySelectorAll(".plan-card")).map((card) => {
      const cta = card.querySelector("a");
      if (!cta) return null;
      const cardBox = card.getBoundingClientRect();
      const ctaBox = cta.getBoundingClientRect();
      return {
        centerDelta: centerDelta(card, cta),
        bottomGap: Math.round(cardBox.bottom - ctaBox.bottom),
      };
    }).filter(Boolean) as Array<{ centerDelta: number; bottomGap: number }>;

    const entryCenterDeltas = Array.from(document.querySelectorAll(".entry-panel")).map((card) => {
      const cta = card.querySelector("a");
      return cta ? centerDelta(card, cta) : 999;
    });

    const connectors = Array.from(document.querySelectorAll(".device-flow__connector")).map((connector) => {
      const parent = connector.parentElement;
      const connectorBox = connector.getBoundingClientRect();
      const parentBox = parent?.getBoundingClientRect();
      const connectorStyles = getComputedStyle(connector);
      return {
        extendsPastCard: Boolean(parentBox && connectorBox.right > parentBox.right),
        display: connectorStyles.display,
        parentOverflowX: parent ? getComputedStyle(parent).overflowX : "",
      };
    });

    const connectorVisuals = Array.from(document.querySelectorAll(".device-flow__connector")).map((connector) => {
      const connectorStyles = getComputedStyle(connector);
      return {
        backgroundImage: connectorStyles.backgroundImage,
        boxShadow: connectorStyles.boxShadow,
        zIndex: connectorStyles.zIndex,
      };
    });

    return { planMetrics, entryCenterDeltas, connectors, connectorVisuals };
  });

  expect(alignment.planMetrics).toHaveLength(3);
  expect(alignment.planMetrics.every((metric) => metric.centerDelta <= 1)).toBe(true);
  expect(new Set(alignment.planMetrics.map((metric) => metric.bottomGap)).size).toBe(1);
  expect(alignment.entryCenterDeltas.every((delta) => delta <= 1)).toBe(true);
  if (alignment.connectors.every((connector) => connector.display === "none")) {
    expect(alignment.connectors).toEqual([
      { extendsPastCard: false, display: "none", parentOverflowX: "visible" },
      { extendsPastCard: false, display: "none", parentOverflowX: "visible" },
    ]);
  } else {
    expect(alignment.connectors).toEqual([
      { extendsPastCard: true, display: "block", parentOverflowX: "visible" },
      { extendsPastCard: true, display: "block", parentOverflowX: "visible" },
    ]);
    expect(alignment.connectorVisuals.every((connector) => connector.backgroundImage.includes("rgb(123, 220, 255)") && connector.backgroundImage.includes("rgb(245, 210, 122)"))).toBe(true);
    expect(alignment.connectorVisuals.every((connector) => connector.boxShadow !== "none")).toBe(true);
    expect(alignment.connectorVisuals.every((connector) => connector.zIndex === "2")).toBe(true);
  }
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
  expect(frameRequests).toBeGreaterThanOrEqual(15);
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
  test.setTimeout(180_000);
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
    expect(text).not.toContain(ruCopy.oldStart);
    expect(text).not.toContain(ruCopy.oldCourses);
    const landingText = await page.locator(".english-spark-landing").innerText();
    if (!cyrillicInterfaceLocales.has(code)) {
      expect(landingText).not.toMatch(/\p{Script=Cyrillic}/u);
    }
  }
});

test("landing localizes core product copy across all 35 interface locales", async ({ page }) => {
  test.setTimeout(120_000);
  await page.goto("/poliglot-ai.html?lang=en");
  const select = page.locator("[data-site-language-select]");
  await expect(select).toBeVisible();

  for (const code of siteLocaleCodes) {
    await page.evaluate((nextCode) => {
      const languageSelect = document.querySelector("[data-site-language-select]") as HTMLSelectElement | null;
      if (!languageSelect) throw new Error("missing public site language selector");
      languageSelect.value = nextCode;
      languageSelect.dispatchEvent(new Event("change", { bubbles: true }));
    }, code);
    await page.waitForTimeout(140);

    const htmlLang = await page.locator("html").getAttribute("lang");
    expect(htmlLang).toBe(code);

    const landingText = await page.locator(".english-spark-landing").innerText();
    expect(landingText).toContain("35");
    if (code !== "en") {
      expect(landingText).not.toContain("Practice speaking before the moment matters");
      expect(landingText).not.toContain("Real situations where the language has to work today");
      expect(landingText).not.toContain("Start free");
      for (const englishCopy of englishLandingCopyThatMustLocalize) {
        expect(landingText).not.toContain(englishCopy);
      }
    }
    if (!cyrillicInterfaceLocales.has(code) && code !== "en") {
      expect(landingText).not.toMatch(/\p{Script=Cyrillic}/u);
    }
  }
});

test("privacy and terms keep contacts and both themes readable", async ({ page }) => {
  for (const path of ["/privacy.html", "/terms.html"]) {
    await page.context().clearCookies();
    await page.addInitScript(() => localStorage.clear());
    await page.goto(path);
    await expect(page.locator("html")).toHaveAttribute("data-site-theme", "dark");

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

    for (const theme of ["dark", "light"] as const) {
      await page.evaluate((nextTheme) => {
        document.documentElement.dataset.siteTheme = nextTheme;
        localStorage.setItem("poliglot-site-theme", nextTheme);
      }, theme);
      const readability = await page.evaluate(() => {
        const selectors = [
          ".legal-page",
          ".legal-hero h1",
          ".legal-hero p",
          ".legal-pills span",
          ".legal-aside a",
          ".legal-price-tags span",
          ".legal-old-price",
          ".legal-document-shell",
          ".legal-document-shell h2",
          ".legal-document-shell p",
          ".legal-document-shell li",
          ".legal-contact-grid a",
        ];
        const nodes = selectors
          .map((selector) => {
            const node = document.querySelector(selector) as HTMLElement | null;
            if (!node) return null;
            const style = getComputedStyle(node);
            const rect = node.getBoundingClientRect();
            return {
              selector,
              color: style.color,
              background: style.backgroundColor,
              width: rect.width,
              overflow: node.scrollWidth > node.clientWidth + 2,
            };
          })
          .filter(Boolean);
        return {
          sameColors: nodes.filter((item) => item!.color === item!.background),
          overflowing: nodes.filter((item) => item!.width > 0 && item!.overflow),
          horizontalOverflow: Math.max(0, document.documentElement.scrollWidth - document.documentElement.clientWidth),
        };
      });
      expect(readability.sameColors).toEqual([]);
      expect(readability.overflowing).toEqual([]);
      expect(readability.horizontalOverflow).toBe(0);
    }
  }
});

test("privacy Russian contact grid keeps Telegram bot visible", async ({ page }) => {
  await page.goto("/privacy.html?lang=ru");

  const botContact = page.locator('#contacts .legal-contact-grid a[href="https://t.me/poliglot_ai_bot"]');
  await expect(botContact).toBeVisible();
  await expect(botContact.locator("small")).toHaveText("Telegram bot");
  await expect(botContact).toContainText("@poliglot_ai_bot");
});

test("legal document package exposes operator details and cookie opt-in", async ({ page }) => {
  test.setTimeout(90_000);

  for (const path of ["/privacy.html", "/terms.html", "/agreement.html", "/consent.html"]) {
    await page.goto(`${path}?lang=ru`);
    await expect(page.locator(".legal-page")).toBeVisible();
    await expect(page.locator(".legal-document-shell")).toContainText("Самозанятый Чебан Денис Игоревич");
    await expect(page.locator(".legal-document-shell")).toContainText("ИНН 505017471160");
    await expect(page.locator(".legal-document-shell")).toContainText("supportpoliglotai@gmail.com");
    await expect(page.locator(".legal-aside a[href*='privacy.html']")).toBeVisible();
    await expect(page.locator(".legal-aside a[href*='terms.html']")).toBeVisible();
    await expect(page.locator(".legal-aside a[href*='agreement.html']")).toBeVisible();
    await expect(page.locator(".legal-aside a[href*='consent.html']")).toBeVisible();
  }

  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  await page.evaluate(() => localStorage.removeItem("poliglot-cookie-consent"));
  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  const banner = page.locator(".cookie-consent-banner");
  await expect(banner).toBeVisible();
  await expect(banner).toContainText("cookie");
  await expect(banner.locator("a[href*='privacy.html']")).toBeVisible();
  await expect(banner.locator("a[href*='consent.html']")).toBeVisible();
  await banner.getByRole("button", { name: /Принять/ }).click();
  await expect(banner).toHaveCount(0);
  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  await expect(page.locator(".cookie-consent-banner")).toHaveCount(0);
});

test("production public-site output keeps all image and localization assets for deploy", async ({ page, request }) => {
  const assetPaths = [
    "/assets/site-i18n.js",
    "/assets/site-phrases.js",
    "/assets/privacy-policy-i18n.js",
    "/assets/legal-documents-i18n.js",
    "/assets/scenarios/travel-ai-tutor.jpg",
    "/assets/scenarios/work-ai-tutor.jpg",
    "/assets/scenarios/exam-ai-tutor.jpg",
    "/assets/scenarios/speaking-ai-tutor.jpg",
    "/assets/testimonials/anna.jpg",
    "/assets/testimonials/marat.jpg",
    "/assets/testimonials/sofia.jpg",
  ];

  for (const assetPath of assetPaths) {
    const response = await request.get(assetPath);
    expect(response.status(), `${assetPath} should be shipped with the built public site`).toBe(200);
    const body = await response.body();
    expect(body.length, `${assetPath} should not be empty`).toBeGreaterThan(100);
  }

  await page.goto("/poliglot-ai.html?lang=en");
  for (const image of await page.locator(".scenario-card img, .review-card img").all()) {
    await image.scrollIntoViewIfNeeded();
  }
  await page.waitForFunction(() =>
    Array.from(document.images)
      .filter((image) => image.src.includes("/assets/scenarios/") || image.src.includes("/assets/testimonials/"))
      .every((image) => image.complete && image.naturalWidth > 0),
  );
  const imageState = await page.evaluate(() =>
    Array.from(document.images)
      .filter((image) => image.src.includes("/assets/scenarios/") || image.src.includes("/assets/testimonials/"))
      .map((image) => ({
        src: new URL(image.src).pathname,
        naturalWidth: image.naturalWidth,
        naturalHeight: image.naturalHeight,
        complete: image.complete,
      })),
  );

  expect(imageState).toHaveLength(7);
  expect(imageState.every((image) => image.complete && image.naturalWidth >= 240 && image.naturalHeight >= 160)).toBe(true);
  expect(imageState.filter((image) => image.src.includes("/assets/scenarios/")).every((image) => image.naturalWidth >= 480)).toBe(true);
  expect(imageState.filter((image) => image.src.includes("/assets/testimonials/")).every((image) => image.naturalWidth >= 240)).toBe(true);
});
