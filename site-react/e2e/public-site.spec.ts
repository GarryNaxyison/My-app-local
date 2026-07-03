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
  "AI Tutor is checking the answer",
  "weak word fixed",
  "translation saved",
  "Four tools, one learning profile",
  "Payment is available through Telegram Stars and YooKassa/SBP.",
  "Premium key giveaways and product updates in the NERIVA channels",
  "Follow the community for short lessons, product updates, and monthly Premium key giveaways.",
  "Questions before starting",
  "Open the loop and run the first lesson today",
  "Follow NERIVA",
  "Short lessons, product updates, and learning tips.",
  "AI Tutor Cockpit",
  "Photograph a menu, sign, or task and turn it into a learning scenario.",
  "No peanuts, please. How spicy is this dish?",
  "The photo becomes translation, a note, and a practice prompt.",
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

const currentRussianLandingCopy = [
  "NERIVA: \u041f\u043e\u0442\u0440\u0435\u043d\u0438\u0440\u0443\u0439\u0442\u0435 \u0440\u0435\u0447\u044c \u0434\u043e \u0432\u0430\u0436\u043d\u043e\u0433\u043e \u043c\u043e\u043c\u0435\u043d\u0442\u0430",
  "\u0416\u0438\u0432\u043e\u0439 \u0441\u0446\u0435\u043d\u0430\u0440\u0438\u0439 \u043f\u0440\u043e\u0434\u0443\u043a\u0442\u0430",
  "\u041a\u0430\u0436\u0434\u0430\u044f \u0441\u0435\u0441\u0441\u0438\u044f \u0437\u0430\u043a\u0430\u043d\u0447\u0438\u0432\u0430\u0435\u0442\u0441\u044f \u043f\u043e\u043d\u044f\u0442\u043d\u044b\u043c \u0441\u043b\u0435\u0434\u0443\u044e\u0449\u0438\u043c \u0448\u0430\u0433\u043e\u043c",
  "\u0427\u0435\u0442\u044b\u0440\u0435 \u0438\u043d\u0441\u0442\u0440\u0443\u043c\u0435\u043d\u0442\u0430, \u043e\u0434\u0438\u043d \u0443\u0447\u0435\u0431\u043d\u044b\u0439 \u043f\u0440\u043e\u0444\u0438\u043b\u044c",
  "\u041e\u043f\u043b\u0430\u0442\u0430 \u0434\u043e\u0441\u0442\u0443\u043f\u043d\u0430 \u0447\u0435\u0440\u0435\u0437 Telegram Stars \u0438 YooKassa/SBP.",
  "\u0420\u043e\u0437\u044b\u0433\u0440\u044b\u0448\u0438 Premium-\u043a\u043b\u044e\u0447\u0435\u0439",
  "\u0412\u043e\u043f\u0440\u043e\u0441\u044b \u043f\u0435\u0440\u0435\u0434 \u0441\u0442\u0430\u0440\u0442\u043e\u043c",
  "\u0417\u0430\u043f\u0443\u0441\u0442\u0438\u0442\u0435 \u0446\u0438\u043a\u043b \u0438 \u043d\u0430\u0447\u043d\u0438\u0442\u0435 \u043f\u0435\u0440\u0432\u044b\u0439 \u0443\u0440\u043e\u043a \u0441\u0435\u0433\u043e\u0434\u043d\u044f",
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

test("landing presents the approved English spark hero product site", async ({ page }) => {
  test.setTimeout(120_000);
  await keepCookieBannerHidden(page);
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
  expect(heroHeight).toBeGreaterThanOrEqual(viewportHeight * 0.72);

  await expect(page.locator("h1")).toContainText("Practice speaking before the moment matters");
  await expect(page.locator(".hero-proof")).toContainText("35 languages");
  await expect(page.locator(".hero-proof")).toContainText("A1-C2");
  const heroSocialLinks = page.locator(".landing-social-proof a");
  await expect(page.locator(".landing-social-proof")).toContainText("Follow NERIVA");
  await expect(page.locator(".landing-social-proof")).toContainText("Short lessons, product updates, and learning tips.");
  await expect(heroSocialLinks).toHaveCount(4);
  await expect(heroSocialLinks.nth(0)).toHaveAttribute("href", "https://www.youtube.com/@neriva_app");
  await expect(heroSocialLinks.nth(0)).toHaveAttribute("aria-label", "Open NERIVA on YouTube");
  await expect(heroSocialLinks.nth(0)).toHaveAttribute("target", "_blank");
  await expect(heroSocialLinks.nth(0)).toHaveAttribute("rel", "noreferrer");
  await expect(heroSocialLinks.nth(1)).toHaveAttribute("href", "https://www.instagram.com/neriva.ru");
  await expect(heroSocialLinks.nth(1)).toHaveAttribute("aria-label", "Open NERIVA on Instagram");
  await expect(heroSocialLinks.nth(2)).toHaveAttribute("href", "https://tiktok.com/@nerivaru");
  await expect(heroSocialLinks.nth(2)).toHaveAttribute("aria-label", "Open NERIVA on TikTok");
  await expect(heroSocialLinks.nth(3)).toHaveAttribute("href", "https://t.me/NERIVAapp_bot");
  await expect(heroSocialLinks.nth(3)).toHaveAttribute("aria-label", "Open NERIVA on Telegram");
  await expect(heroSocialLinks.nth(3).locator('img[src="/assets/social/telegram-logo.png"]')).toBeVisible();

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

  await expect(page.locator(".lesson-scenario")).toBeVisible();
  await expect(page.locator(".lesson-scenario")).toContainText("Every session ends with a visible next step");
  await expect(page.locator(".lesson-scenario img")).toHaveCount(1);
  await expect(page.locator(".module-card")).toHaveCount(4);
  await expect(page.locator(".module-card", { hasText: "Voice Coach" })).toContainText("weak words");
  await expect(page.locator(".module-card", { hasText: "Photo Practice" })).toContainText("menu or sign");

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
  await expect(page.locator(".payment-methods")).toContainText("YooKassa/SBP");
  await expect(page.locator(".payment-methods")).not.toContainText("TON");
  await expect(page.locator(".payment-methods")).not.toContainText("USDT");
  await expect(page.locator(".payment-methods")).not.toContainText(/crypto|blockchain|RollyPay/i);

  await expect(page.locator(".community-section")).toBeVisible();
  await expect(page.locator(".community-section")).toContainText("Premium key giveaways");
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

  await expect(page.locator(".final-cta-section")).toContainText("Web app");
  await expect(page.locator(".final-cta-section")).toContainText("Telegram");

  await page.locator(".nav-burger").click();
  await expect(page.locator(".nav-drawer")).toBeVisible();
  await expect(page.locator(".nav-drawer")).toContainText("Features");
  await expect(page.locator(".nav-drawer")).toContainText("Pricing");
  await expect(page.locator(".nav-drawer")).toContainText("FAQ");
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

test("cookie banner stays hidden with legacy saved consent values", async ({ page }) => {
  for (const legacyValue of ["necessary", "accepted"]) {
    await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
    await page.evaluate((value) => localStorage.setItem("poliglot-cookie-consent", value), legacyValue);
    await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
    await expect(page.locator(".cookie-consent-banner")).toBeHidden();
  }
});

test("landing light theme keeps the approved layout readable", async ({ page }) => {
  test.setTimeout(90_000);
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
  await expect(page.locator(".community-section .eyebrow")).toHaveCSS("color", "rgb(15, 23, 42)");
  await expect(page.locator(".plan-label").first()).toHaveCSS("color", "rgb(45, 91, 255)");
  await expect(page.locator(".lesson-scenario__tag").first()).toHaveCSS("color", "rgb(0, 47, 167)");
  await expect(page.locator(".landing-faq__item h3").first()).toHaveCSS("color", "rgb(5, 7, 10)");
  await expect(page.locator(".landing-social-proof strong")).toHaveCSS("color", "rgb(7, 17, 31)");
  await expect(page.locator(".landing-social-proof span")).toHaveCSS("color", "rgba(7, 17, 31, 0.68)");
  await expect(page.locator(".landing-social-proof .social-icon-link--tiktok")).toHaveCSS("color", "rgb(7, 17, 31)");

  const issues = await page.evaluate(() => {
    const selectors = [
      ".public-nav",
      ".landing-hero h1",
      ".landing-hero__copy p",
      ".hero-demo",
      ".lesson-scenario",
      ".module-card",
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
    const overflowing = Array.from(document.querySelectorAll("h1,h2,h3,p,a,button,.landing-social-proof,.social-icon-links,.site-footer-social,.plan-card,.module-card,.lesson-scenario,.landing-faq__item,.community-section"))
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

test("landing keeps a short product flow below the hero", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  await expect(page.locator(".landing-hero")).toBeVisible();
  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);

  await expect(page.locator(".lesson-scenario")).toBeVisible();
  await expect(page.locator(".lesson-scenario")).toContainText("AI Tutor is checking the answer");
  await expect(page.locator(".lesson-scenario")).toContainText("XP +12");
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
    { src: "/assets/scenarios/speaking-ai-tutor.jpg", alt: "Speaking practice lesson preview", naturalWidth: expect.any(Number), complete: true },
  ]);
  expect(scenarioImages.every((img) => img.naturalWidth >= 480)).toBe(true);

  await expect(page.locator(".module-card")).toHaveCount(4);
  await expect(page.locator(".landing-faq__item")).toHaveCount(4);
  await expect(page.locator(".seo-guides a")).toHaveCount(5);
  await expect(page.locator(".course-strip, .cockpit-section, .scenario-section, .device-flow, .memory-loop-section")).toHaveCount(0);
});

test("Russian landing uses edited copy and keeps compact CTAs aligned", async ({ page }) => {
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
    };

    return { planMetrics, sectionCounts };
  });

  expect(alignment.planMetrics).toHaveLength(3);
  expect(alignment.planMetrics.every((metric) => metric.centerDelta <= 1)).toBe(true);
  expect(new Set(alignment.planMetrics.map((metric) => metric.bottomGap)).size).toBe(1);
  expect(alignment.sectionCounts).toEqual({ faq: 4, community: 1, oldDeviceFlow: 0, oldEntryPanels: 0 });
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
  await expect(page.locator(".lesson-scenario")).toBeVisible();
  await expect(page.locator(".module-card")).toHaveCount(4);
  await expect(page.locator(".landing-faq__item")).toHaveCount(4);
  await expect(page.locator(".site-footer a", { hasText: "Privacy" })).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: "Terms" })).toBeVisible();

  const heroHeight = await page.locator(".landing-hero").evaluate((node) => node.getBoundingClientRect().height);
  const viewportHeight = await page.evaluate(() => window.innerHeight);
  expect(heroHeight).toBeGreaterThanOrEqual(viewportHeight * 0.9);

  const overflow = await page.evaluate(() => {
    const documentOverflow = document.documentElement.scrollWidth > document.documentElement.clientWidth + 2 || document.body.scrollWidth > document.body.clientWidth + 2;
    const nodes = Array.from(
      document.querySelectorAll(
        ".public-nav, .nav-drawer, [data-site-language-select], h1, h2, h3, p, a, button, .hero-proof span, .landing-social-proof, .social-icon-links, .site-footer-social, .lesson-scenario, .landing-faq__item, .community-section, .module-card, .plan-card",
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
  for (const image of await page.locator(".lesson-scenario img").all()) {
    await image.scrollIntoViewIfNeeded();
  }
  await page.waitForFunction(() =>
    Array.from(document.images)
      .filter((image) => image.src.includes("/assets/scenarios/"))
      .every((image) => image.complete && image.naturalWidth > 0),
  );
  const imageState = await page.evaluate(() =>
    Array.from(document.images)
      .filter((image) => image.src.includes("/assets/scenarios/"))
      .map((image) => ({
        src: new URL(image.src).pathname,
        naturalWidth: image.naturalWidth,
        naturalHeight: image.naturalHeight,
        complete: image.complete,
      })),
  );

  expect(imageState).toHaveLength(1);
  expect(imageState.every((image) => image.complete && image.naturalWidth >= 240 && image.naturalHeight >= 160)).toBe(true);
  expect(imageState.every((image) => image.naturalWidth >= 480)).toBe(true);
});
