import { expect, test, type Page } from "@playwright/test";

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
  "NERIVA: Practice speaking before the moment matters",
  "Live product scenario",
  "Every session ends with a visible next step",
  "NERIVA explains a phrase, asks for an answer, checks it, saves the mistake, and gives a repeatable prompt for the next review.",
  "Four tools, one learning profile",
  "Payment is available through Telegram Stars and YooKassa/SBP.",
  "Premium key giveaways and product updates in the NERIVA channels",
  "Follow the community for short lessons, product updates, and monthly Premium key giveaways.",
  "Questions before starting",
  "Open the loop and run the first lesson today",
] as const;

const curatedLandingCopySources = [
  "NERIVA: Practice speaking before the moment matters",
  "Live product scenario",
  "Every session ends with a visible next step",
  "NERIVA explains a phrase, asks for an answer, checks it, saves the mistake, and gives a repeatable prompt for the next review.",
  "AI Tutor is checking the answer",
  "weak word fixed",
  "translation saved",
  "Four tools, one learning profile",
  "Payment is available through Telegram Stars and YooKassa/SBP.",
  "Premium key giveaways and product updates in the NERIVA channels",
  "Follow the community for short lessons, product updates, and monthly Premium key giveaways.",
  "Questions before starting",
  "Open the loop and run the first lesson today",
] as const;

const curatedRussianLandingCopy = [
  "NERIVA: Потренируйте речь до важного момента",
  "Живой сценарий продукта",
  "Каждая сессия заканчивается понятным следующим шагом",
  "Четыре инструмента, один учебный профиль",
  "Оплата доступна через Telegram Stars и YooKassa/SBP.",
  "Розыгрыши Premium-ключей",
  "Вопросы перед стартом",
] as const;

const currentRussianLandingCopy = [
  "NERIVA \u043f\u043e\u043c\u043e\u0433\u0430\u0435\u0442 \u0433\u043e\u0432\u043e\u0440\u0438\u0442\u044c \u0443\u0432\u0435\u0440\u0435\u043d\u043d\u0435\u0435",
  "\u0416\u0438\u0432\u043e\u0439 \u0443\u0440\u043e\u043a",
  "\u041e\u0442\u0432\u0435\u0442\u0438\u043b\u0438, \u0443\u0432\u0438\u0434\u0435\u043b\u0438 \u043e\u0448\u0438\u0431\u043a\u0443, \u043f\u043e\u0432\u0442\u043e\u0440\u0438\u043b\u0438",
  "\u0427\u0435\u0442\u044b\u0440\u0435 \u0444\u0443\u043d\u043a\u0446\u0438\u0438 \u0432 \u043e\u0434\u043d\u043e\u043c \u043f\u0440\u043e\u0444\u0438\u043b\u0435",
  "\u041e\u043f\u043b\u0430\u0442\u0430 \u0447\u0435\u0440\u0435\u0437 Telegram Stars \u0438 YooKassa/SBP.",
  "\u0420\u043e\u0437\u044b\u0433\u0440\u044b\u0448\u0438 \u0438 \u043d\u043e\u0432\u043e\u0441\u0442\u0438 \u043f\u0440\u043e\u0434\u0443\u043a\u0442\u0430",
  "\u041f\u0435\u0440\u0435\u0434 \u0441\u0442\u0430\u0440\u0442\u043e\u043c",
  "\u041d\u0430\u0447\u043d\u0438\u0442\u0435 \u043f\u0435\u0440\u0432\u044b\u0439 \u0443\u0440\u043e\u043a \u0441\u0435\u0433\u043e\u0434\u043d\u044f",
  "\u0412\u0435\u0431-\u043f\u0440\u0438\u043b\u043e\u0436\u0435\u043d\u0438\u0435",
] as const;

const awkwardRussianLandingCopy = [
  "Мобильный Интернет",
  "План выступлений на день",
  "AI Репетитор",
  "более длительную AI репетиторскую работу",
  "Более высокие лимиты на поездки",
] as const;

test.beforeEach(async ({ page }) => {
  await page.route("https://mc.yandex.ru/**", (route) => route.abort());
});

async function keepCookieBannerHidden(page: Page) {
  await page.addInitScript(() => {
    localStorage.setItem(
      "poliglot-cookie-consent",
      JSON.stringify({ version: 1, necessary: true, analyticsMarketing: false }),
    );
  });
}

test.skip("landing presents the approved light product-first site", async ({ page }) => {
  test.setTimeout(120_000);
  await keepCookieBannerHidden(page);
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator("html")).toHaveAttribute("data-site-theme", "light");
  await expect(page.locator(".public-nav")).toBeVisible();
  await expect(page.locator(".public-brand__logo img")).toBeVisible();
  await expect(page.locator(".nav-burger")).toBeVisible();
  await expect(page.locator("[data-site-language-select]")).toBeVisible();

  const navPosition = await page.locator(".public-nav").evaluate((node) => getComputedStyle(node).position);
  expect(navPosition).toBe("sticky");

  await expect(page.locator(".english-spark-landing")).toBeVisible();
  await expect(page.locator(".landing-hero")).toBeVisible();
  await expect(page.locator(".landing-hero")).toHaveAttribute("data-hero-preset", "light");
  await expect(page.locator(".hero-demo")).toHaveCount(0);
  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator(".spark-hero")).toHaveCount(0);
  await expect(page.locator(".bold-hero")).toHaveCount(0);
  const heroHeight = await page.locator(".landing-hero").evaluate((node) => node.getBoundingClientRect().height);
  const viewportHeight = await page.evaluate(() => window.innerHeight);
  expect(heroHeight).toBeGreaterThanOrEqual(viewportHeight * 0.68);
  expect(heroHeight).toBeLessThanOrEqual(viewportHeight * 1.08);

  await expect(page.locator("h1")).toContainText("helps you speak with more confidence");
  await expect(page.locator(".hero-proof")).toContainText("35 languages");
  await expect(page.locator(".hero-proof")).toContainText("A1-C2");

  const heroCtas = page.locator(".landing-hero .entry-cta");
  await expect(heroCtas).toHaveCount(2);
  expect(await page.locator('a[data-entry="web-app"]').count()).toBeGreaterThanOrEqual(2);
  const webEntryHrefs = await page.locator('a[data-entry="web-app"]').evaluateAll((links) =>
    links.map((link) => (link as HTMLAnchorElement).getAttribute("href")),
  );
  expect(webEntryHrefs.every((href) => href === "/app/")).toBe(true);
  expect(await page.locator('a[data-entry="telegram"]').count()).toBeGreaterThanOrEqual(2);
  const telegramEntryHrefs = await page.locator('a[data-entry="telegram"]').evaluateAll((links) =>
    links.map((link) => (link as HTMLAnchorElement).getAttribute("href")),
  );
  expect(telegramEntryHrefs.every((href) => href === "https://t.me/NERIVAapp_bot")).toBe(true);

  const productShots = page.locator(".product-shot img");
  await expect(productShots).toHaveCount(13);
  await expect(page.locator('.hero-product-frame img[src="/assets/product/dashboard-progress.png"]')).toBeVisible();
  await expect(page.locator('.scenario-shot img[src="/assets/product/ai-tutor-lesson-correction.png"]')).toBeVisible();
  await expect(page.locator('.feature-panel img[src="/assets/product/voice-pronunciation-score.png"]')).toBeVisible();
  await expect(page.locator('.telegram-panel img[src="/assets/product/telegram-app-light.png"]')).toBeVisible();
  await productShots.last().scrollIntoViewIfNeeded();
  await expect.poll(() =>
    productShots.evaluateAll((imgs) =>
      imgs.every((img) => (img as HTMLImageElement).complete && (img as HTMLImageElement).naturalWidth > 0),
    ),
  ).toBe(true);
  await page.evaluate(() => window.scrollTo(0, 0));

  await expect(page.locator(".lesson-scenario")).toBeVisible();
  await expect(page.locator(".lesson-scenario")).toContainText("Practice, get corrected, repeat");
  await expect(page.locator(".lesson-scenario")).toContainText("Hotel check-in");
  await expect(page.locator(".feature-panel")).toHaveCount(4);
  await expect(page.locator(".feature-panel", { hasText: "Voice" })).toContainText("weak words");
  await expect(page.locator(".feature-panel", { hasText: "Photo" })).toContainText("menu");

  await expect(page.locator(".plan-card")).toHaveCount(3);
  await expect(page.locator(".plan-card", { hasText: "Free" })).toContainText("0 ₽");
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText("300 ₽");
  await expect(page.locator(".plan-card", { hasText: "Platinum" })).toContainText("590 ₽");
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText("Daily practice");
  await expect(page.locator(".plan-card", { hasText: "Platinum" })).toContainText("Intensive");
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
  await expect(page.locator(".payment-methods")).toContainText("YooKassa/SBP");
  await expect(page.locator(".payment-methods")).not.toContainText("TON");
  await expect(page.locator(".payment-methods")).not.toContainText("USDT");
  await expect(page.locator(".payment-methods")).not.toContainText(/crypto|blockchain|RollyPay/i);

  await expect(page.locator(".community-section")).toBeVisible();
  await expect(page.locator(".community-section")).toContainText("Giveaways");
  await expect(page.locator(".telegram-panel")).toContainText("Telegram");
  await expect(page.locator(".landing-faq")).toBeVisible();
  await expect(page.locator(".landing-faq__item")).toHaveCount(4);
  await expect(page.locator(".seo-guides a")).toHaveCount(5);

  await expect(page.locator(".comparison-card")).toHaveCount(0);
  await expect(page.locator(".answer-card")).toHaveCount(0);
  await expect(page.locator(".course-strip")).toHaveCount(0);
  await expect(page.locator(".cockpit-section")).toHaveCount(0);
  await expect(page.locator(".scenario-card")).toHaveCount(0);
  await expect(page.locator(".device-flow")).toHaveCount(0);
  await expect(page.locator(".memory-loop-section")).toHaveCount(0);
  await expect(page.locator(".review-card")).toHaveCount(0);
  await expect(page.locator(".landing-hero .seo-guides, .public-nav .seo-guides")).toHaveCount(0);

  await expect(page.locator(".final-cta-section")).toContainText("Browser");
  await expect(page.locator(".final-cta-section")).toContainText("Telegram");

  await page.locator(".nav-burger").click();
  const drawer = page.locator(".nav-drawer");
  await expect(drawer).toBeVisible();
  await expect(drawer).toContainText("Features");
  await expect(drawer).toContainText("Pricing");
  await expect(drawer).toContainText("FAQ");
  await expect(drawer).toContainText("Privacy");
  await expect(drawer).toContainText("Terms");
  await expect(drawer).toContainText("Browser");
  await expect(drawer).toContainText("Telegram");

  const navLayout = await page.evaluate(() => {
    const nav = document.querySelector(".public-nav") as HTMLElement;
    const burger = document.querySelector(".nav-burger") as HTMLElement;
    const drawerElement = document.querySelector(".nav-drawer") as HTMLElement;
    const navRect = nav.getBoundingClientRect();
    const burgerRect = burger.getBoundingClientRect();
    const drawerRect = drawerElement.getBoundingClientRect();
    return {
      navHeight: navRect.height,
      burgerHeight: burgerRect.height,
      burgerLines: burger.querySelectorAll(".nav-burger__lines span").length,
      drawerLeft: drawerRect.left,
      drawerRight: drawerRect.right,
      drawerTop: drawerRect.top,
      viewportWidth: window.innerWidth,
    };
  });
  expect(navLayout.navHeight).toBeLessThanOrEqual(82);
  expect(navLayout.burgerHeight).toBe(44);
  expect(navLayout.burgerLines).toBe(3);
  expect(navLayout.drawerLeft).toBeGreaterThanOrEqual(0);
  expect(navLayout.drawerRight).toBeLessThanOrEqual(navLayout.viewportWidth);
  expect(navLayout.drawerTop).toBeGreaterThan(navLayout.burgerHeight);

  const landingCyrillic = await page.locator(".english-spark-landing").evaluate((node) => (node.textContent || "").match(/\p{Script=Cyrillic}+/gu) || []);
  expect(landingCyrillic).toEqual([]);

  const footerPrivacyLink = page.locator(".site-footer a", { hasText: "Privacy" });
  const footerTermsLink = page.locator(".site-footer a", { hasText: "Terms" });
  await expect(footerPrivacyLink).toBeVisible();
  await expect(footerPrivacyLink).toHaveAttribute("href", /privacy\.html(\?.*)?$/);
  await expect(footerTermsLink).toBeVisible();
  await expect(footerTermsLink).toHaveAttribute("href", /terms\.html(\?.*)?$/);
});

test("public footer exposes NERIVA social channels", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  const footerSocial = page.locator(".site-footer .site-footer-social");
  await expect(footerSocial).toBeVisible();
  await expect(footerSocial).toContainText("Social");
  await expect(footerSocial.locator("a")).toHaveCount(4);
  await expect(footerSocial.locator('a[href="https://www.youtube.com/@neriva_app"]')).toHaveAttribute("aria-label", "Open NERIVA on YouTube");
  await expect(footerSocial.locator('a[href="https://www.instagram.com/neriva.ru"]')).toHaveAttribute("aria-label", "Open NERIVA on Instagram");
  await expect(footerSocial.locator('a[href="https://tiktok.com/@nerivaru"]')).toHaveAttribute("aria-label", "Open NERIVA on TikTok");
  await expect(footerSocial.locator('a[href="https://t.me/NERIVAapp_bot"]')).toHaveAttribute("aria-label", "Open NERIVA on Telegram");
  await expect(footerSocial.locator('a[href="https://t.me/NERIVAapp_bot"]')).toHaveAttribute("target", "_blank");
  await expect(footerSocial.locator('a[href="https://t.me/NERIVAapp_bot"]')).toHaveAttribute("rel", "noreferrer");
  await expect(footerSocial.locator('a[href="https://t.me/NERIVAapp_bot"] img[src="/assets/social/telegram-logo.png"]')).toBeVisible();

  await page.goto("/poliglot-ai.html?lang=ru");
  const ruFooterSocial = page.locator(".site-footer .site-footer-social");
  await expect(ruFooterSocial).toContainText("Соцсети");
  await expect(ruFooterSocial.locator("a")).toHaveCount(4);
  await expect(ruFooterSocial.locator('a[href="https://www.youtube.com/@neriva_app"]')).toBeVisible();
  await expect(ruFooterSocial.locator('a[href="https://www.instagram.com/neriva.ru"]')).toBeVisible();
  await expect(ruFooterSocial.locator('a[href="https://tiktok.com/@nerivaru"]')).toBeVisible();
  await expect(ruFooterSocial.locator('a[href="https://t.me/NERIVAapp_bot"]')).toBeVisible();
});

test("landing visual blocks animate without breaking hero and ecosystem alignment", async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto("/poliglot-ai.html?lang=ru");

  for (const index of [1, 2, 3, 0]) {
    await page.locator(".feature-carousel__tab").nth(index).click();
    await page.waitForTimeout(80);
  }

  const heroState = await page.evaluate(() => {
    const matter = document.querySelector(".landing-hero__matter") as HTMLElement;
    const canvases = Array.from(document.querySelectorAll(".landing-hero__matter canvas")) as HTMLCanvasElement[];
    const rect = matter.getBoundingClientRect();
    return {
      width: rect.width,
      height: rect.height,
      background: getComputedStyle(matter).backgroundImage,
      nonEmptyCanvases: canvases.filter((canvas) => canvas.width > 10 && canvas.height > 10).length,
    };
  });
  expect(heroState.width).toBeGreaterThan(900);
  expect(heroState.height).toBeGreaterThan(500);
  expect(heroState.background).toContain("radial-gradient");
  expect(heroState.nonEmptyCanvases).toBeGreaterThanOrEqual(1);

  await page.locator(".ecosystem-showcase").scrollIntoViewIfNeeded();
  const textTops = await page.locator(".ecosystem-card > div").evaluateAll((nodes) => nodes.map((node) => Math.round((node as HTMLElement).getBoundingClientRect().top)));
  expect(Math.max(...textTops) - Math.min(...textTops)).toBeLessThanOrEqual(3);
});

test("landing shows Russian legal links and centers the cookie banner on desktop", async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  await page.evaluate(() => localStorage.removeItem("poliglot-cookie-consent"));
  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });

  await page.locator(".nav-burger").click();
  const drawer = page.locator(".nav-drawer");
  await expect(drawer.locator('a[href*="agreement.html"]')).toHaveText("Пользовательское соглашение");
  await expect(drawer.locator('a[href*="consent.html"]')).toHaveText("Согласие на обработку персональных данных");

  const footerDocuments = page.locator(".site-footer nav").filter({ hasText: "Документы" });
  await expect(footerDocuments).toBeVisible();
  await expect(footerDocuments.locator('a[href*="agreement.html"]')).toHaveText("Пользовательское соглашение");
  await expect(footerDocuments.locator('a[href*="consent.html"]')).toHaveText("Согласие на обработку персональных данных");

  const banner = page.locator(".cookie-consent-banner");
  await expect(banner).toBeVisible();
  const bannerBox = await banner.evaluate((node) => {
    const rect = node.getBoundingClientRect();
    return {
      left: rect.left,
      right: rect.right,
      width: rect.width,
      viewportWidth: window.innerWidth,
      computedLeft: getComputedStyle(node).left,
    };
  });

  const center = bannerBox.left + bannerBox.width / 2;
  expect(Math.abs(center - bannerBox.viewportWidth / 2)).toBeLessThanOrEqual(2);
  expect(bannerBox.right).toBeLessThanOrEqual(bannerBox.viewportWidth - 14);
  expect(bannerBox.computedLeft).not.toBe("auto");
});

test("cookie banner lets users configure optional cookies", async ({ page }) => {
  test.setTimeout(60_000);

  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  await page.evaluate(() => localStorage.removeItem("poliglot-cookie-consent"));
  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });

  const banner = page.locator(".cookie-consent-banner");
  await expect(banner).toBeVisible();
  await expect(banner.getByRole("button", { name: "Настроить" })).toBeVisible();

  await banner.getByRole("button", { name: "Настроить" }).click();
  await expect(banner.locator(".cookie-consent-banner__settings")).toBeVisible();
  await expect(banner).toContainText("Необходимые");
  await expect(banner).toContainText("Аналитические/маркетинговые");
  await expect.poll(() => page.evaluate(() => localStorage.getItem("poliglot-cookie-consent"))).toBeNull();

  const optionalSwitch = banner.getByRole("switch", { name: "Аналитические/маркетинговые" });
  await expect(optionalSwitch).toHaveAttribute("aria-checked", "false");
  await optionalSwitch.click();
  await expect(optionalSwitch).toHaveAttribute("aria-checked", "true");

  await banner.getByRole("button", { name: "Сохранить выбор" }).click();
  await expect(banner).toBeHidden();

  const storedConsent = await page.evaluate(() => localStorage.getItem("poliglot-cookie-consent"));
  expect(JSON.parse(storedConsent || "{}")).toEqual({
    version: 1,
    necessary: true,
    analyticsMarketing: true,
  });

  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  await expect(page.locator(".cookie-consent-banner")).toBeHidden();
});

test("Yandex Metrika waits for explicit cookie banner analytics consent", async ({ page, context }) => {
  test.setTimeout(60_000);
  await page.unroute("https://mc.yandex.ru/**");

  const metrikaRequests: string[] = [];
  await page.route("https://mc.yandex.ru/**", (route) => {
    metrikaRequests.push(route.request().url());
    return route.abort();
  });

  await page.addInitScript(() => {
    localStorage.setItem("poliglot-cookie-consent", "accepted");
  });
  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });

  await expect(page.locator(".cookie-consent-banner")).toBeHidden();
  await expect.poll(() => metrikaRequests.length).toBe(0);
  await expect.poll(() =>
    page.evaluate(() => ({
      hasMetrikaScript: Boolean(document.querySelector('script[data-yandex-metrika-id="109242081"]')),
      hasYmStub: typeof window.ym === "function",
    })),
  ).toEqual({ hasMetrikaScript: false, hasYmStub: false });

  const acceptPage = await context.newPage();
  const acceptMetrikaRequests: string[] = [];
  await acceptPage.route("https://mc.yandex.ru/**", (route) => {
    acceptMetrikaRequests.push(route.request().url());
    return route.abort();
  });
  await acceptPage.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  await acceptPage.evaluate(() => localStorage.removeItem("poliglot-cookie-consent"));
  await acceptPage.reload({ waitUntil: "domcontentloaded" });

  const banner = acceptPage.locator(".cookie-consent-banner");
  await expect(banner).toBeVisible();
  await expect.poll(() => acceptMetrikaRequests.length).toBe(0);

  await banner.getByRole("button", { name: "Принять все cookie" }).click();

  await expect.poll(() =>
    acceptPage.evaluate(() => Boolean(document.querySelector('script[data-yandex-metrika-id="109242081"]'))),
  ).toBe(true);
  await expect.poll(() => acceptMetrikaRequests.some((url) => url.includes("mc.yandex.ru/metrika/tag.js?id=109242081"))).toBe(true);
  await acceptPage.close();
});

test("cookie banner stays hidden with legacy saved consent values", async ({ page }) => {
  for (const legacyValue of ["necessary", "accepted"]) {
    await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
    await page.evaluate((value) => localStorage.setItem("poliglot-cookie-consent", value), legacyValue);
    await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
    await expect(page.locator(".cookie-consent-banner")).toBeHidden();
  }
});

test.skip("landing light theme keeps the approved layout readable", async ({ page }) => {
  test.setTimeout(90_000);
  await page.addInitScript(() => {
    localStorage.setItem("poliglot-site-theme", "light");
  });
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator("html")).toHaveAttribute("data-site-theme", "light");
  await expect(page.locator(".landing-hero")).toBeVisible();
  await expect(page.locator(".landing-hero")).toHaveAttribute("data-hero-preset", "light");
  await expect(page.locator(".hero-demo")).toHaveCount(0);
  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator(".hero-product-frame")).toBeVisible();
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
      heroBackgroundColor: getComputedStyle(hero).backgroundColor,
      matterOpacity: Number(matterStyle.opacity),
      matterPointerEvents: matterStyle.pointerEvents,
    };
  });

  expect(lightHeroHeader.navBackground).toEqual(expect.objectContaining({ red: expect.any(Number), green: expect.any(Number), blue: expect.any(Number), alpha: expect.any(Number) }));
  expect(lightHeroHeader.navBackground!.alpha).toBeGreaterThanOrEqual(0.92);
  expect(lightHeroHeader.heroColor).toBe("rgb(5, 7, 10)");
  expect(lightHeroHeader.headingColor).toBe("rgb(5, 7, 10)");
  expect(lightHeroHeader.heroBackgroundColor).toBe("rgb(255, 255, 255)");
  expect(lightHeroHeader.heroBackgroundImage).not.toContain("gradient");
  expect(lightHeroHeader.matterOpacity).toBeGreaterThanOrEqual(0.16);
  expect(lightHeroHeader.matterOpacity).toBeLessThanOrEqual(0.34);
  expect(lightHeroHeader.matterPointerEvents).toBe("none");
  await expect(page.locator(".pricing-section .eyebrow")).toHaveCSS("color", "rgb(0, 47, 167)");
  await expect(page.locator(".community-section .eyebrow")).toHaveCSS("color", "rgb(0, 47, 167)");
  await expect(page.locator(".plan-label").first()).toHaveCSS("color", "rgb(0, 47, 167)");
  await expect(page.locator(".lesson-scenario__tag").first()).toHaveCSS("color", "rgb(0, 47, 167)");
  await expect(page.locator(".landing-faq__item h3").first()).toHaveCSS("color", "rgb(5, 7, 10)");

  const issues = await page.evaluate(() => {
    const selectors = [
      ".public-nav",
      ".landing-hero h1",
      ".landing-hero__copy p",
      ".hero-product-frame",
      ".lesson-scenario",
      ".feature-panel",
      ".plan-card",
      ".landing-faq__item",
      ".community-section",
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
    const overflowing = Array.from(document.querySelectorAll("h1,h2,h3,p,a,button,.social-icon-links,.site-footer-social,.plan-card,.feature-panel,.lesson-scenario,.landing-faq__item,.community-section,.product-shot"))
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

test.skip("landing keeps a short product flow below the hero", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".landing-hero")).toBeVisible();
  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);

  await expect(page.locator(".lesson-scenario")).toBeVisible();
  await expect(page.locator(".lesson-scenario")).toContainText("NERIVA turns a real situation");
  await expect(page.locator(".lesson-scenario")).toContainText("Hotel check-in");
  await expect(page.locator(".lesson-scenario img")).toHaveCount(1);
  await page.locator(".lesson-scenario img").scrollIntoViewIfNeeded();
  await page.waitForFunction(() =>
    Array.from(document.querySelectorAll<HTMLImageElement>(".lesson-scenario img")).every((image) => image.complete && image.naturalWidth > 0),
  );
  const scenarioImages = await page.locator(".lesson-scenario img").evaluateAll((imgs) =>
    imgs.map((img) => ({
      src: (img as HTMLImageElement).getAttribute("src"),
      alt: (img as HTMLImageElement).getAttribute("alt"),
      naturalWidth: (img as HTMLImageElement).naturalWidth,
      complete: (img as HTMLImageElement).complete,
    })),
  );
  expect(scenarioImages).toEqual([
    { src: "/assets/product/ai-tutor-lesson-correction.png", alt: "AI tutor lesson with a corrected hotel check-in answer", naturalWidth: expect.any(Number), complete: true },
  ]);
  expect(scenarioImages.every((img) => img.naturalWidth >= 900)).toBe(true);

  await expect(page.locator(".feature-panel")).toHaveCount(4);
  await expect(page.locator(".landing-faq__item")).toHaveCount(4);
  await expect(page.locator(".seo-guides a")).toHaveCount(5);
  await expect(page.locator(".course-strip, .cockpit-section, .scenario-section, .device-flow, .memory-loop-section")).toHaveCount(0);
});

test.skip("Russian landing uses edited copy and keeps compact CTAs aligned", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=ru");
  await page.evaluate(() => {
    document.documentElement.dataset.siteTheme = "light";
    localStorage.setItem("poliglot-site-theme", "light");
  });
  await expect(page.locator(".english-spark-landing")).toBeVisible();

  const landingText = await page.locator(".english-spark-landing").innerText();
  for (const expected of currentRussianLandingCopy) {
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

    const sectionCounts = {
      faq: document.querySelectorAll(".landing-faq__item").length,
      community: document.querySelectorAll(".community-section").length,
      oldDeviceFlow: document.querySelectorAll(".device-flow__connector").length,
      oldEntryPanels: document.querySelectorAll(".entry-panel").length,
      productShots: document.querySelectorAll(".product-shot img").length,
    };

    return { planMetrics, sectionCounts };
  });

  expect(alignment.planMetrics).toHaveLength(3);
  expect(alignment.planMetrics.every((metric) => metric.centerDelta <= 1)).toBe(true);
  expect(new Set(alignment.planMetrics.map((metric) => metric.bottomGap)).size).toBe(1);
  expect(alignment.sectionCounts).toEqual({ faq: 4, community: 1, oldDeviceFlow: 0, oldEntryPanels: 0, productShots: 13 });
});

test("landing hero keeps the animation while product screenshots remain primary", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(2);
  await expect(page.locator('.landing-hero img[src="/assets/product/dashboard-progress-dark.png"]')).toBeVisible();
  const viewportWidth = page.viewportSize()?.width ?? 1440;
  const mobileHomeShot = page.locator('.landing-hero img[src="/assets/product/mobile-home-progress-dark.png"]');
  if (viewportWidth <= 860) {
    await expect(mobileHomeShot).toHaveCount(1);
  } else {
    await expect(mobileHomeShot).toBeVisible();
  }
});

test.skip("landing startup is localized, lean, and theme-aware", async ({ page, request }) => {
  const html = await (await request.get("/poliglot-ai.html")).text();
  expect(html).not.toContain("/assets/site-phrases.js");
  expect(html).not.toContain("/assets/legal-documents-i18n.js");

  await keepCookieBannerHidden(page);
  await page.goto("/poliglot-ai.html?lang=ru");
  await expect(page.locator("h1")).toContainText("NERIVA помогает говорить увереннее");
  await expect(page.locator(".hero-lead")).toContainText("Короткий урок");
  await expect(page.locator(".final-cta-section")).toContainText("Начните первый урок сегодня");
  await expect(page.locator(".final-cta-section")).toContainText("Быстрый старт");

  const lightColors = await page.locator(".english-spark-landing").evaluate((node) => {
    const style = getComputedStyle(node);
    return { background: style.backgroundColor, color: style.color };
  });

  await page.locator(".nav-theme-toggle").first().click();
  await expect(page.locator("html")).toHaveAttribute("data-site-theme", "dark");
  await expect(page.locator(".landing-hero")).toHaveAttribute("data-hero-preset", "dark");
  await expect(page.locator(".landing-hero__matter .generative-art-scene")).toHaveAttribute("data-hero-preset", "dark");
  const darkColors = await page.locator(".english-spark-landing").evaluate((node) => {
    const style = getComputedStyle(node);
    return { background: style.backgroundColor, color: style.color };
  });

  expect(darkColors.background).not.toBe(lightColors.background);
  expect(darkColors.color).not.toBe(lightColors.color);
});

test.skip("product-first landing stays readable on mobile", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator("h1")).toContainText("helps you speak with more confidence");
  await expect(page.locator(".landing-hero .entry-cta")).toHaveCount(2);
  await expect(page.locator(".hero-proof")).toBeHidden();
  await expect(page.locator(".lesson-scenario")).toBeVisible();
  await expect(page.locator(".feature-panel")).toHaveCount(4);
  await expect(page.locator(".landing-faq__item")).toHaveCount(4);
  await expect(page.locator(".site-footer a", { hasText: "Privacy" })).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: "Terms" })).toBeVisible();

  const heroHeight = await page.locator(".landing-hero").evaluate((node) => node.getBoundingClientRect().height);
  const viewportHeight = await page.evaluate(() => window.innerHeight);
  expect(heroHeight).toBeGreaterThanOrEqual(viewportHeight * 0.62);
  expect(heroHeight).toBeLessThanOrEqual(viewportHeight * 1.12);
  await expect(page.locator(".mobile-landing-rail")).toBeVisible();

  const overflow = await page.evaluate(() => {
    const documentOverflow = document.documentElement.scrollWidth > document.documentElement.clientWidth + 2 || document.body.scrollWidth > document.body.clientWidth + 2;
    const nodes = Array.from(
      document.querySelectorAll(
        ".public-nav, .nav-drawer, [data-site-language-select], h1, h2, h3, p, a, button, .hero-proof span, .social-icon-links, .site-footer-social, .lesson-scenario, .landing-faq__item, .community-section, .feature-panel, .plan-card, .product-shot",
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

test.skip("public site language selector exposes all interface locales", async ({ page }) => {
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

test.skip("landing localizes core product copy across all 35 interface locales", async ({ page }) => {
  test.setTimeout(240_000);
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

    await expect(page.locator("html")).toHaveAttribute("lang", code);
    const localizedLanding = await page.evaluate(() => {
      const htmlLang = document.documentElement.getAttribute("lang") || "";
      const landingText = document.querySelector(".english-spark-landing")?.textContent || "";
      const socialProofText = document.querySelector(".landing-social-proof")?.textContent || "";
      return { htmlLang, landingText, socialProofText };
    });

    expect(localizedLanding.htmlLang).toBe(code);
    const landingText = localizedLanding.landingText;
    expect(landingText).toContain("35");
    if (code !== "en") {
      expect(landingText).not.toContain("Practice speaking before the moment matters");
      expect(landingText).not.toContain("Real situations where the language has to work today");
      expect(landingText).not.toContain("Start free");
      expect(localizedLanding.socialProofText).not.toContain("Follow NERIVA");
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
    await expect(page.locator("html")).toHaveAttribute("data-site-theme", "light");

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
    expect(selectColors.color).not.toBe(selectColors.background);
    expect(selectColors.optionColor).not.toBe(selectColors.optionBackground);

    const text = await page.locator("body").innerText();
    expect(text).toContain("@NERIVAapp_bot");
    expect(text).toContain("@NERIVAapp_bot");
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

  const botContact = page.locator('.legal-contact-grid a[href="https://t.me/NERIVAapp_bot"]');
  await expect(botContact).toBeVisible();
  await expect(botContact.locator("small")).toHaveText("Telegram bot");
  await expect(botContact).toContainText("@NERIVAapp_bot");
});

test("legal document package exposes operator details and cookie opt-in", async ({ page }) => {
  test.setTimeout(90_000);

  for (const path of ["/privacy.html", "/terms.html", "/agreement.html", "/consent.html"]) {
    await page.goto(`${path}?lang=ru`);
    await expect(page.locator(".legal-page")).toBeVisible();
    await expect(page.locator(".legal-document-shell")).toContainText("Самозанятый Чебан Денис Игоревич");
    await expect(page.locator(".legal-document-shell")).toContainText("ИНН 505017471160");
    await expect(page.locator(".legal-document-shell")).toContainText("support@neriva.ru");
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
  await expect(banner.getByRole("button", { name: "Настроить" })).toBeVisible();
  await expect(banner.getByRole("button", { name: "Только необходимые" })).toBeVisible();
  await expect(banner.getByRole("button", { name: "Принять все cookie" })).toBeVisible();
  await banner.getByRole("button", { name: "Принять все cookie" }).click();
  await expect(banner).toBeHidden();
  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  await expect(page.locator(".cookie-consent-banner")).toBeHidden();
});

test("legal documents include the updated Telegram, AI, and Metrika clauses", async ({ page }) => {
  test.setTimeout(60_000);

  await page.goto("/privacy.html?lang=ru");
  await expect(page.locator(".legal-document-shell")).toContainText("уникальный идентификатор Telegram (Telegram ID)");
  await expect(page.locator(".legal-document-shell")).toContainText("Данные собираются только после нажатия пользователем кнопки согласия на куки-баннере.");
  await expect(page.locator(".legal-document-shell")).toContainText("трансграничную передачу обезличенных учебных запросов");

  await page.goto("/terms.html?lang=ru");
  await expect(page.locator(".legal-document-shell")).toContainText("без использования классической регистрации по Email");
  await expect(page.locator(".legal-document-shell")).toContainText("ввода одноразового защищенного кода");
  await expect(page.locator(".legal-document-shell")).toContainText("Запрещено отправлять через Сервис конфиденциальную информацию");

  await page.goto("/agreement.html?lang=ru");
  await expect(page.locator(".legal-document-shell")).toContainText("Telegram Stars признается цифровым продуктом экосистемы Telegram");

  await page.goto("/consent.html?lang=ru");
  await expect(page.locator(".legal-document-shell")).toContainText("Сервис не собирает и не обрабатывает Email-адреса, номера телефонов и паспортные данные.");
  await expect(page.locator(".legal-document-shell")).toContainText("страны ЕС и США");
  await expect(page.locator(".legal-document-shell")).toContainText("зарубежными провайдерами моделей искусственного интеллекта (AI)");
});

test("production public-site output keeps all image and localization assets for deploy", async ({ page, request }) => {
  const assetPaths = [
    "/assets/site-i18n.js",
    "/assets/site-phrases.js",
    "/assets/privacy-policy-i18n.js",
    "/assets/legal-documents-i18n.js",
    "/assets/product/dashboard-progress.png",
    "/assets/product/ai-tutor-lesson-correction.png",
    "/assets/product/mistakes-review.png",
    "/assets/product/voice-pronunciation-score.png",
    "/assets/product/photo-translation.png",
    "/assets/product/notes-phrasebook.png",
    "/assets/product/premium-plans.png",
    "/assets/product/telegram-app-light.png",
    "/assets/product/mobile-home-progress.png",
    "/assets/product/mobile-lesson-correction.png",
    "/assets/product/mobile-mistakes.png",
    "/assets/product/mobile-voice-pronunciation.png",
  ];

  for (const assetPath of assetPaths) {
    const response = await request.get(assetPath);
    expect(response.status(), `${assetPath} should be shipped with the built public site`).toBe(200);
    const body = await response.body();
    expect(body.length, `${assetPath} should not be empty`).toBeGreaterThan(100);
  }

  const i18nSource = await (await request.get("/assets/site-i18n.js")).text();
  expect(i18nSource).toContain("englishSparkCreativeCopy");
  for (const source of curatedLandingCopySources) {
    expect(i18nSource, `${source} should be maintained in the curated landing translation layer`).toContain(JSON.stringify(source));
  }

  const phraseSource = await (await request.get("/assets/site-phrases.js")).text();
  const generatedBlocks = Array.from(phraseSource.matchAll(/\/\/ <public-site-translations-generated>([\s\S]*?)\/\/ <\/public-site-translations-generated>/g), (match) => match[1]);
  expect(generatedBlocks).toHaveLength(1);
  for (const source of curatedLandingCopySources) {
    expect(generatedBlocks[0], `${source} should not be produced by the machine translation block`).not.toContain(JSON.stringify(source));
  }

  await page.goto("/poliglot-ai.html?lang=en");
  const heroImage = page.locator('.landing-hero img[src="/assets/product/dashboard-progress-dark.png"]');
  await heroImage.waitFor({ state: "visible" });
  await expect
    .poll(() =>
      heroImage.evaluate((image) => {
        const productImage = image as HTMLImageElement;
        return productImage.complete && productImage.naturalWidth >= 240 && productImage.naturalHeight >= 160;
      }),
    )
    .toBe(true);
  const imageState = await page.evaluate(() =>
    Array.from(document.images)
      .filter((image) => image.src.includes("/assets/product/"))
      .map((image) => ({
        src: new URL(image.src).pathname,
        naturalWidth: image.naturalWidth,
        naturalHeight: image.naturalHeight,
        complete: image.complete,
        visible: image.getClientRects().length > 0,
      })),
  );

  expect(imageState.length).toBeGreaterThanOrEqual(7);
  expect(imageState.some((image) => image.src.includes("dashboard-progress-dark.png") && image.complete && image.naturalWidth >= 240 && image.naturalHeight >= 160)).toBe(true);
  expect(imageState.some((image) => image.src.includes("telegram-app-dark.png"))).toBe(true);
});
