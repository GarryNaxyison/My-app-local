# Premium AI Tutor Cockpit Landing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rework the public React landing into a premium AI tutor cockpit while preserving the existing full-screen hero animation, legal links, `/app/` CTAs, and 35-language support.

**Architecture:** Keep the landing in `site-react` and treat `Сайт полиглота для бота/` as generated build output. Add tests first for the new premium cockpit contract, then update React content, then CSS desktop/mobile presentation, then rebuild and verify the generated static site.

**Tech Stack:** React 19, Vite, TypeScript, Framer Motion, lucide-react, Tailwind CSS import, Playwright, Nx CLI fallback.

---

## File Structure

- Modify `site-react/e2e/public-site.spec.ts`: encode the premium cockpit contract, `/app/` target, legal links, 35 locales, desktop/mobile text safety, and preserved hero canvas.
- Modify `site-react/src/PublicSiteApp.tsx`: update landing arrays, hero copy, product demo, route cards, pricing copy, FAQ, and footer copy while keeping legal pages intact.
- Modify `site-react/src/styles.css`: keep the full-screen animation placement, restore dark premium hero treatment in light and dark themes, and add intentional desktop/mobile cockpit layouts.
- Generate `Сайт полиглота для бота/poliglot-ai.html`, `privacy.html`, `terms.html`, and `assets/site-react/*` through `npm --prefix site-react run build`; do not edit generated output by hand.
- Do not modify `web-react`, `web`, Go backend files, payment behavior, or `GenerativeArtScene` internals.

## Task 1: Update Landing Tests First

**Files:**
- Modify: `site-react/e2e/public-site.spec.ts`

- [ ] **Step 1: Replace the main landing test contract**

In `site-react/e2e/public-site.spec.ts`, replace the first test with:

```ts
test("landing presents the premium AI tutor cockpit without losing public contracts", async ({ page }) => {
  test.setTimeout(60_000);
  await page.goto("/poliglot-ai.html");

  await expect(page.locator(".public-nav")).toBeVisible();
  await expect(page.locator(".public-nav a", { hasText: /^v2$/i })).toHaveCount(0);
  await expect(page.locator(".public-brand__logo img")).toBeVisible();
  await expect(page.locator('.public-nav a[href="/privacy.html"]')).toBeVisible();
  await expect(page.locator('.public-nav a[href="/terms.html"]')).toBeVisible();

  const appLinks = await page.locator('a[href^="/app"]').evaluateAll((links) => [...new Set(links.map((link) => (link as HTMLAnchorElement).getAttribute("href")))].sort());
  expect(appLinks).toEqual(["/app/"]);

  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  const lightHero = await page.locator(".landing-hero").evaluate((hero) => {
    document.documentElement.dataset.siteTheme = "light";
    const heroBox = hero.getBoundingClientRect();
    const matter = hero.querySelector(".landing-hero__matter") as HTMLElement | null;
    const matterBox = matter?.getBoundingClientRect();
    const h1 = hero.querySelector("h1") as HTMLElement;
    const firstGoal = hero.querySelector(".hero-goals a strong") as HTMLElement;
    return {
      heroWidth: heroBox.width,
      heroHeight: heroBox.height,
      matterWidth: matterBox?.width || 0,
      matterHeight: matterBox?.height || 0,
      h1Color: getComputedStyle(h1).color,
      firstGoalColor: getComputedStyle(firstGoal).color,
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

  await expect(page.locator("h1")).toContainText("AI-репетитором");
  await expect(page.locator(".hero-goals a")).toHaveCount(4);
  await expect(page.locator(".hero-goals")).toContainText("Путешествия");
  await expect(page.locator(".hero-goals")).toContainText("Работа");
  await expect(page.locator(".hero-goals")).toContainText("Экзамен");
  await expect(page.locator(".hero-goals")).toContainText("Разговорная речь");
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

  await expect(page.locator('.site-footer a[href="/privacy.html"]')).toBeVisible();
  await expect(page.locator('.site-footer a[href="/terms.html"]')).toBeVisible();
});
```

- [ ] **Step 2: Add mobile layout safety test**

Append this test after the hero frame-budget test:

```ts
test("premium cockpit landing stays readable on mobile", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/poliglot-ai.html");

  await expect(page.locator(".landing-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator("h1")).toContainText("AI-репетитором");
  await expect(page.locator(".hero-action--primary").first()).toBeVisible();
  await expect(page.locator(".hero-goals a")).toHaveCount(4);
  await expect(page.locator(".hero-demo")).toBeVisible();
  await expect(page.locator('.site-footer a[href="/privacy.html"]')).toBeVisible();
  await expect(page.locator('.site-footer a[href="/terms.html"]')).toBeVisible();

  const overflow = await page.evaluate(() => {
    const nodes = Array.from(document.querySelectorAll("h1, h2, h3, p, a, button, .hero-goals a, .hero-proof span, .pricing-grid article, .review-card, .faq-grid article"));
    return nodes
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
  });
  expect(overflow).toEqual([]);
});
```

- [ ] **Step 3: Run focused tests to verify RED**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "premium AI tutor cockpit|stays readable on mobile" --project=desktop-chromium
```

Expected: FAIL because `.hero-goals`, updated hero copy, updated `/app/` expectation, `.payment-methods`, and mobile cockpit layout are not implemented yet.

- [ ] **Step 4: Keep the RED test changes for the implementation task**

Do not commit the failing tests by themselves. Confirm they are the only current worktree change for this feature:

```powershell
git status --short
```

Expected: `site-react/e2e/public-site.spec.ts` is modified and the focused Playwright command above has failed for the expected missing-cockpit reasons.

## Task 2: Update Landing Content And Structure

**Files:**
- Modify: `site-react/src/PublicSiteApp.tsx`

- [ ] **Step 1: Replace top-level landing data arrays**

In `site-react/src/PublicSiteApp.tsx`, replace `features`, `workflow`, `testimonials`, `faqs`, `demoScreens`, and add `heroGoals` plus `courseRoutes` after `workflow`:

```ts
const features = [
  {
    icon: BookOpen,
    title: "AI Tutor маршрут",
    body: "Короткая сессия ведет от фразы и примера к ответу, проверке, XP и следующему повторению.",
  },
  {
    icon: MessageCircle,
    title: "Ролевые диалоги",
    body: "Поездка, работа, экзамен или разговорная речь: AI держит сценарий, исправляет ответ и дает модельную фразу.",
  },
  {
    icon: Mic,
    title: "Произношение и shadowing",
    body: "Голосовая тренировка показывает оценку, слабые слова, звуки, историю прогресса и следующий повтор.",
  },
  {
    icon: Camera,
    title: "Фото и перевод",
    body: "Меню, вывеска, задание или файл превращаются в перевод, заметку и короткую практику по контексту.",
  },
  {
    icon: WifiOff,
    title: "Ошибки, словарь и offline",
    body: "Review, Spelling, Notes, Mistakes и offline decks возвращают к тому, что реально проседает.",
  },
  {
    icon: Trophy,
    title: "Единый прогресс",
    body: "XP, streak, уровни, награды, daily bonus и история занятий остаются в одном профиле web, PWA и Telegram.",
  },
];

const workflow = [
  ["01", "Выберите цель, язык и уровень", "Путешествие, работа, экзамен или разговорная речь превращаются в маршрут под ваш уровень A1-C2."],
  ["02", "Пройдите короткую сессию", "AI Tutor, практика, роль, голос или фото работают в web app, PWA и Telegram."],
  ["03", "Повторите слабые места", "Ошибки, заметки, слабые слова, voice history и offline decks закрепляются за одним профилем."],
];

const heroGoals = [
  ["Путешествия", "меню, отель, аэропорт", "роль + фото"],
  ["Работа", "small talk, встречи, интервью", "диалоги"],
  ["Экзамен", "структура ответа и точность", "A1-C2"],
  ["Разговорная речь", "голос, shadowing, фразы", "voice coach"],
] as const;

const courseRoutes = [
  ["Для поездки", "Соберите фразы для меню, отеля, аэропорта и неожиданных ситуаций. Фото-перевод и roleplay помогают сразу применить слова.", "путешествия"],
  ["Для работы", "Тренируйте small talk, встречи, собеседование и деловые ответы. AI исправляет формулировки и дает модельные фразы.", "работа"],
  ["Для экзамена и речи", "Держите структуру ответа, произношение, shadowing и слабые слова в одном маршруте до уверенного повторения.", "экзамен + голос"],
] as const;

const testimonials = [
  {
    name: "Анна",
    role: "готовится к поездке",
    text: "Я тренирую меню и отель через роль, а фото-перевод сразу сохраняет нужные фразы. Это больше похоже на личный маршрут, чем на список слов.",
  },
  {
    name: "Марат",
    role: "учит английский для работы",
    text: "В Telegram удобно отвечать быстро, а в web app видно ошибки, заметки и прогресс. Один профиль между ними сильно экономит время.",
  },
  {
    name: "София",
    role: "прокачивает произношение",
    text: "После голосовой практики понятно, какие слова повторить. Мне важен не балл ради балла, а конкретный следующий шаг.",
  },
];

const faqs = [
  ["Можно заниматься только в Telegram?", "Да. Telegram работает как отдельный учебный вход, а web app и PWA дают больше места для прогресса, ошибок, заметок и offline decks."],
  ["Что общее между web app и Telegram?", "Профиль, Premium, прогресс, ошибки, заметки и учебная история связаны с одним аккаунтом."],
  ["Чем Premium и Platinum отличаются от Free?", "Free подходит для старта. Premium и Platinum дают больше уроков и практики, голосовые функции, фото-перевод и расширенные лимиты для регулярной учебы."],
  ["Можно ли учиться на телефоне и без сети?", "Да. Web app работает как PWA, а offline decks помогают повторять сохраненные слова, фразы и ошибки без стабильного соединения."],
];
```

- [ ] **Step 2: Replace demo screens**

Replace `demoScreens` with:

```ts
const demoScreens = [
  {
    tab: "Урок",
    avatar: "УР",
    intro: "Начните с цели: AI объяснит фразу, контекст и что сказать дальше.",
    label: "Путешествие",
    output: "I would like to check in, please.",
    hint: "Фраза сразу попадает в маршрут: пример, перевод, повторение и роль.",
  },
  {
    tab: "Диалог",
    avatar: "AI",
    intro: "Ответьте как в реальной ситуации, а AI поправит тон, грамматику и естественность.",
    label: "Roleplay",
    output: "Could you help me check in? I have a reservation.",
    hint: "После ответа появляется короткий разбор и модельная фраза.",
  },
  {
    tab: "Голос",
    avatar: "ГС",
    intro: "Произнесите фразу и получите оценку речи, слабые слова и следующий повтор.",
    label: "Shadowing",
    output: "Please speak a little slower.",
    hint: "Pronunciation и shadowing показывают, что улучшить в речи.",
  },
  {
    tab: "Фото",
    avatar: "ФО",
    intro: "Сфотографируйте меню, вывеску или задание и превратите это в практику.",
    label: "Фото-перевод",
    output: "No peanuts, please. How spicy is this dish?",
    hint: "Фото превращается в перевод, заметку и учебный сценарий.",
  },
] as const;
```

- [ ] **Step 3: Update hero JSX**

Inside `LandingPage`, replace the hero eyebrow, H1, body, primary CTA text, language-picker heading, and add `heroGoals` after the hero actions:

```tsx
<span className="eyebrow">Premium AI-репетитор для ежедневной практики</span>
<h1>Говорите на новом языке с AI-репетитором, который помнит ваши ошибки</h1>
<p>
  Урок, диалог, произношение, фото-перевод, словарь ошибок и прогресс собираются в один профиль web, PWA и Telegram.
</p>
<div className="hero-actions">
  <a className="hero-action hero-action--primary" href="/app/">
    Начать бесплатно <ArrowRight size={18} />
  </a>
  <a className="hero-action hero-action--secondary" href="https://t.me/poliglot_ai_bot">
    Открыть Telegram
  </a>
</div>
<div className="hero-goals" aria-label="Цели обучения">
  {heroGoals.map(([title, detail, badge]) => (
    <a key={title} href="/app/">
      <strong>{title}</strong>
      <span>{detail}</span>
      <small>{badge}</small>
    </a>
  ))}
</div>
```

Then change the language picker heading to:

```tsx
<h3>Популярные языки обучения:</h3>
```

Then change the hero path labels to:

```tsx
<span><b>1</b> Цель</span>
<span><b>2</b> Урок</span>
<span><b>3</b> Речь</span>
<span><b>4</b> Повтор</span>
```

Then change the third proof chip to:

```tsx
<strong>Free</strong> старт без оплаты
```

- [ ] **Step 4: Update course strip JSX**

Replace the inline course array in `.course-cards` with `courseRoutes`:

```tsx
{courseRoutes.map(([title, body, badge], index) => (
  <motion.a key={title} href="/app/" className="course-card" initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} transition={{ duration: 0.42, delay: index * 0.05 }}>
    <span>{badge}</span>
    <h3>{title}</h3>
    <p>{body}</p>
    <b>
      Начать маршрут <ChevronRight size={16} />
    </b>
  </motion.a>
))}
```

Also update the course strip intro to:

```tsx
<span className="eyebrow">Premium cockpit</span>
<h2>Каждая цель превращается в маршрут: урок, роль, голос, ошибки и повторение</h2>
<p>Poliglot AI показывает, что делать сегодня, где вы ошиблись и что повторить дальше. Web app, PWA и Telegram держат один учебный профиль.</p>
```

- [ ] **Step 5: Update pricing and final CTA**

Replace the pricing section intro and cards with:

```tsx
<span className="eyebrow">Premium</span>
<h2>Начните бесплатно, а для серьезной практики откройте голос, фото и расширенные лимиты</h2>
<p className="payment-methods">Оплата: Telegram Stars, YooKassa/SBP, TON и USDT.</p>
```

Use these `PlanCard` calls:

```tsx
<PlanCard name="Free" label="Попробовать маршрут" price="0 ₽" body="Для первого знакомства с AI Tutor и ежедневной привычкой." items={["стартовые уроки", "базовая практика", "заметки и прогресс"]} />
<PlanCard name="Premium" label="Регулярная учеба" oldPrice="1000 ₽" price="300 ₽" body="Основной режим для ежедневной практики с голосом, фото и большим числом сессий." items={["до 50 уроков в день", "до 200 сообщений практики", "голос, фото и словарь ошибок"]} featured />
<PlanCard name="Platinum" label="Интенсив" oldPrice="2000 ₽" price="590 ₽" body="Максимум для поездки, работы, экзамена или активной разговорной практики." items={["до 100 уроков в день", "до 500 сообщений практики", "больше ролей и голосовых повторов"]} />
```

Update `.start-panel` copy to:

```tsx
<span className="eyebrow">Первый маршрут</span>
<h2>Откройте Poliglot AI и начните говорить с разбором ошибок уже сегодня</h2>
```

Update the final CTA text to:

```tsx
Начать бесплатно <ArrowRight size={18} />
```

- [ ] **Step 6: Run tests and record remaining CSS/layout failures**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "premium AI tutor cockpit" --project=desktop-chromium
```

Expected: content-specific assertions for H1, `.hero-goals`, `/app/`, `.payment-methods`, and demo tabs should pass. Color/layout assertions may still fail until Task 3.

- [ ] **Step 7: Keep content changes uncommitted until CSS is green**

Run:

```powershell
git status --short
```

Expected: `site-react/e2e/public-site.spec.ts` and `site-react/src/PublicSiteApp.tsx` are modified. Do not commit yet; Task 3 creates the first green source/test commit.

## Task 3: Implement Premium Desktop And Mobile Styling

**Files:**
- Modify: `site-react/src/styles.css`

- [ ] **Step 1: Add premium cockpit overrides after the current landing overrides**

Append this block near the end of `site-react/src/styles.css`, after the current landing/mobile overrides:

```css
/* 2026-06-15 premium AI tutor cockpit refresh. */
.landing-hero {
  background: linear-gradient(135deg, #050914 0%, #07111f 54%, #0a172a 100%);
  color: #fff;
}

html[data-site-theme="light"] .landing-hero,
html[data-site-theme="dark"] .landing-hero {
  background: linear-gradient(135deg, #050914 0%, #07111f 54%, #0a172a 100%);
}

.landing-hero__matter {
  inset: -16% -8% -6% 33%;
  opacity: 0.95;
  mix-blend-mode: normal;
}

html[data-site-theme="light"] .landing-hero__matter,
html[data-site-theme="dark"] .landing-hero__matter {
  opacity: 0.95;
  mix-blend-mode: normal;
}

.landing-hero__veil,
html[data-site-theme="light"] .landing-hero__veil,
html[data-site-theme="dark"] .landing-hero__veil {
  background:
    radial-gradient(circle at 74% 30%, rgba(122, 220, 255, 0.2), transparent 36%),
    radial-gradient(circle at 54% 74%, rgba(245, 210, 122, 0.14), transparent 34%),
    linear-gradient(90deg, rgba(5, 9, 20, 0.99) 0%, rgba(5, 9, 20, 0.84) 42%, rgba(5, 9, 20, 0.28) 100%),
    linear-gradient(180deg, rgba(5, 9, 20, 0.02), rgba(5, 9, 20, 0.9));
}

.landing-hero h1,
html[data-site-theme="light"] .landing-hero h1,
html[data-site-theme="dark"] .landing-hero h1 {
  max-width: 900px;
  color: #fff;
  font-size: clamp(3.35rem, 7.8vw, 7.2rem);
  line-height: 0.9;
}

.landing-hero__copy > p,
html[data-site-theme="light"] .landing-hero__copy > p,
html[data-site-theme="dark"] .landing-hero__copy > p {
  max-width: 720px;
  color: rgba(255, 255, 255, 0.82);
}

.hero-goals {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
  margin-top: 24px;
  max-width: 820px;
}

.hero-goals a {
  display: grid;
  gap: 6px;
  min-height: 118px;
  padding: 14px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 8px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.12), rgba(255, 255, 255, 0.055)),
    rgba(6, 12, 26, 0.5);
  box-shadow: 0 20px 52px rgba(0, 0, 0, 0.22);
  color: rgba(255, 255, 255, 0.74);
  transition:
    transform 180ms ease,
    border-color 180ms ease,
    background 180ms ease;
}

.hero-goals a:hover {
  border-color: rgba(123, 220, 255, 0.42);
  background:
    linear-gradient(180deg, rgba(123, 220, 255, 0.15), rgba(255, 255, 255, 0.065)),
    rgba(6, 12, 26, 0.62);
  transform: translateY(-3px);
}

.hero-goals strong {
  color: #fff;
  font-size: 1rem;
}

.hero-goals span {
  color: rgba(255, 255, 255, 0.68);
  font-size: 0.82rem;
  line-height: 1.35;
}

.hero-goals small {
  width: fit-content;
  padding: 5px 8px;
  border-radius: 999px;
  background: rgba(245, 210, 122, 0.14);
  color: #f5d27a;
  font-size: 0.72rem;
  font-weight: 900;
}

.hero-language-picker h3,
html[data-site-theme="light"] .hero-language-picker h3 {
  color: rgba(255, 255, 255, 0.88);
}

.hero-language-picker a,
html[data-site-theme="light"] .hero-language-picker a {
  border-color: rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 22px 46px rgba(0, 0, 0, 0.2);
}

.hero-language-picker strong,
html[data-site-theme="light"] .hero-language-picker strong {
  color: #fff;
}

.hero-language-picker span,
html[data-site-theme="light"] .hero-language-picker span {
  color: #7bdcff;
}

.hero-language-picker small,
html[data-site-theme="light"] .hero-language-picker small {
  color: rgba(255, 255, 255, 0.66);
}

.hero-proof span,
html[data-site-theme="light"] .hero-proof span {
  border-color: rgba(255, 255, 255, 0.14);
  background: rgba(255, 255, 255, 0.07);
  color: rgba(255, 255, 255, 0.72);
  box-shadow: 0 18px 38px rgba(0, 0, 0, 0.18);
}

.hero-proof strong,
html[data-site-theme="light"] .hero-proof strong {
  color: #fff;
}

.payment-methods {
  margin: 10px auto 0;
  max-width: 760px;
  color: var(--muted);
  font-weight: 800;
  text-align: center;
}
```

- [ ] **Step 2: Add responsive cockpit refinements**

Append this block after the previous CSS:

```css
@media (max-width: 980px) {
  .landing-hero__matter,
  html[data-site-theme="light"] .landing-hero__matter,
  html[data-site-theme="dark"] .landing-hero__matter {
    inset: -8% -26% 32% 4%;
    opacity: 0.72;
  }

  .hero-goals {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    max-width: 640px;
  }

  .hero-language-picker > div {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 680px) {
  .landing-hero {
    min-height: auto;
  }

  .landing-hero__inner {
    min-height: auto;
    padding-top: 74px;
    padding-bottom: 42px;
  }

  .landing-hero h1,
  html[data-site-theme="light"] .landing-hero h1,
  html[data-site-theme="dark"] .landing-hero h1 {
    font-size: clamp(2.85rem, 14.5vw, 4.2rem);
    line-height: 0.95;
    overflow-wrap: anywhere;
  }

  .landing-hero__copy > p {
    font-size: 1rem;
    line-height: 1.62;
  }

  .hero-actions {
    align-items: stretch;
  }

  .hero-action {
    width: 100%;
    min-width: 0;
  }

  .hero-goals,
  .hero-language-picker > div {
    grid-template-columns: 1fr;
  }

  .hero-goals a {
    min-height: 0;
  }

  .hero-language-picker {
    margin-top: 22px;
  }

  .hero-language-picker a {
    min-height: 0;
  }

  .hero-demo {
    width: 100%;
  }

  .demo-message {
    grid-template-columns: 38px minmax(0, 1fr);
  }

  .demo-wave {
    overflow: hidden;
  }

  .start-panel {
    display: grid;
    padding: 24px;
  }
}
```

- [ ] **Step 3: Run focused tests to verify GREEN**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "premium AI tutor cockpit|stays readable on mobile" --project=desktop-chromium
```

Expected: PASS for the new premium cockpit tests.

- [ ] **Step 4: Commit green source, style, and test changes**

Run:

```powershell
git add site-react/e2e/public-site.spec.ts site-react/src/PublicSiteApp.tsx site-react/src/styles.css
git commit -m "feat: refresh premium cockpit landing source"
```

## Task 4: Full Site Build, Generated Output, And Regression Tests

**Files:**
- Modify: generated files under `Сайт полиглота для бота/`
- Verify: `site-react/e2e/public-site.spec.ts`

- [ ] **Step 1: Run full site build**

Run:

```powershell
npm --prefix site-react run build
```

Expected: exit code 0 and Vite writes `poliglot-ai.html`, `privacy.html`, `terms.html`, and `assets/site-react/*` under `Сайт полиглота для бота/`.

- [ ] **Step 2: Run full public-site E2E suite**

Run:

```powershell
npm --prefix site-react run e2e -- --reporter=line
```

Expected: all public-site tests pass on desktop and mobile Chromium. The old `/app/v2/` failure must be gone.

- [ ] **Step 3: Inspect generated diff**

Run:

```powershell
git status --short
git diff --stat -- site-react "Сайт полиглота для бота"
```

Expected: changes only in `site-react` source/test files and generated public-site output. No `web-react`, Go backend, `.env`, or unrelated files.

- [ ] **Step 4: Commit generated output**

Run:

```powershell
git add site-react/e2e/public-site.spec.ts site-react/src/PublicSiteApp.tsx site-react/src/styles.css "Сайт полиглота для бота"
git commit -m "feat: refresh premium cockpit landing"
```

Expected: commit succeeds and includes source, tests, and generated public-site build output.

## Task 5: Browser Verification And Final Push

**Files:**
- No source edits unless verification finds a concrete bug.

- [ ] **Step 1: Start the public-site dev server**

Run:

```powershell
npm --prefix site-react run dev
```

Expected: Vite serves the public site on `http://127.0.0.1:5175/`. If port 5175 is busy, use the URL printed by Vite.

- [ ] **Step 2: Verify desktop in the browser**

Open `http://127.0.0.1:5175/poliglot-ai.html` at 1440x900 and verify:

- hero canvas is visible and nonblank;
- hero text is readable over the animation;
- primary CTA goes to `/app/`;
- Privacy and Terms links are present;
- 35-language selector is visible;
- the page presents hero goals, product demo, course cards, features, pricing, reviews, FAQ, and footer without overlap.

- [ ] **Step 3: Verify mobile in the browser**

Resize to 390x844 and verify:

- hero animation still renders behind the mobile hero;
- H1, CTAs, goals, language cards, demo, pricing, reviews, FAQ, Terms, and Privacy links do not overflow horizontally;
- primary CTA is visible before the demo becomes too deep;
- footer legal links remain reachable.

- [ ] **Step 4: Fix any browser-found bug with TDD**

If verification finds a bug, first add or tighten a Playwright assertion in `site-react/e2e/public-site.spec.ts`, run it to fail, then update `PublicSiteApp.tsx` or `styles.css`, then rerun the focused test to pass.

Use this command for focused reruns:

```powershell
npm --prefix site-react run e2e -- --grep "premium AI tutor cockpit|stays readable on mobile" --project=desktop-chromium
```

- [ ] **Step 5: Run final verification**

Run:

```powershell
npm --prefix site-react run build
npm --prefix site-react run e2e -- --reporter=line
git status --short
```

Expected: build passes, E2E passes, and git status contains only intentional committed changes or is clean after final commit.

- [ ] **Step 6: Push the branch**

Run:

```powershell
git push origin HEAD
```

Expected: branch `codex/ai-tutor-rebuild-fix` pushes successfully to `origin`.

## Self-Review Checklist

- Spec coverage: tasks cover premium cockpit positioning, hero preservation, `/app/`, Terms/Privacy, 35 languages, desktop/mobile layouts, pricing/payment methods, generated output, Playwright verification, and GitHub push.
- Red-flag scan: no deferred implementation markers or vague future-work instructions are allowed in this plan.
- Type consistency: new arrays are `heroGoals`, `courseRoutes`, and existing `demoScreens`, `features`, `workflow`, `testimonials`, `faqs`; tests reference `.hero-goals`, `.payment-methods`, and existing selectors.
- Scope: no `web-react`, backend, payment behavior, legal content rewrite, or `GenerativeArtScene` internals are included.
