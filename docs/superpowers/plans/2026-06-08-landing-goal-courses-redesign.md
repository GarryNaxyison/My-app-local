# Landing Goal Courses Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rebuild the public Poliglot AI landing copy and hero presentation around the approved `Goal Courses` direction while keeping the current block list and using the old full-screen hero animation.

**Architecture:** Keep the public landing in `site-react`; update React content/data first, then CSS, then generated static output in `Сайт полиглота для бота/`. The `/app/` web application in `web-react` is out of scope and must not be edited or rebuilt for this task.

**Tech Stack:** React 19, Vite, TypeScript, Tailwind CSS import, Framer Motion, lucide-react, Playwright.

---

## File Structure

- Modify `site-react/e2e/public-site.spec.ts`: update public landing assertions so tests describe the approved goal-based landing.
- Modify `site-react/src/PublicSiteApp.tsx`: update landing data arrays, hero copy, CTAs, course cards, features, pricing, reviews, progress, FAQ, and final CTA.
- Modify `site-react/src/styles.css`: make the hero animation first-viewport/full-screen, add goal selector styling, keep text readable, and preserve responsive behavior.
- Build output under `Сайт полиглота для бота/`: produced by `npm run site:build`; do not edit generated HTML/JS/CSS by hand.

## Task 1: Update Public Site E2E Assertions

**Files:**
- Modify: `site-react/e2e/public-site.spec.ts`

- [ ] **Step 1: Update the main landing smoke test expectations**

In `site-react/e2e/public-site.spec.ts`, update the first test body so the landing assertions target the approved copy and preserved counts:

```ts
test("landing keeps goal-course product flow with Poliglot features", async ({ page }) => {
  test.setTimeout(60_000);
  await page.goto("/poliglot-ai.html");

  await expect(page.locator(".public-nav")).toBeVisible();
  await expect(page.locator(".public-nav a", { hasText: /^v2$/i })).toHaveCount(0);
  await expect(page.locator(".public-brand__logo img")).toBeVisible();
  const appLinks = await page.locator('a[href^="/app"]').evaluateAll((links) => [...new Set(links.map((link) => (link as HTMLAnchorElement).getAttribute("href")))].sort());
  expect(appLinks).toEqual(["/app/v2/"]);
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

  await expect(page.locator("h1")).toContainText("Выберите цель");
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
  await expect(page.locator(".payment-methods")).toContainText("Stars");
  await expect(page.locator(".payment-methods")).toContainText("YooKassa");
  await expect(page.locator(".payment-methods")).toContainText("TON");
  await expect(page.locator(".payment-methods")).toContainText("USDT");
});
```

- [ ] **Step 2: Run the focused E2E test and confirm it fails before implementation**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "goal-course product flow" --project=desktop-chromium
```

Expected: FAIL because `.hero-goals`, the new H1, the updated demo tab labels, and `.payment-methods` are not implemented yet.

- [ ] **Step 3: Commit the failing test**

Run:

```powershell
git add site-react/e2e/public-site.spec.ts
git commit -m "test: specify goal-course landing"
```

## Task 2: Update Landing Data and Hero Structure

**Files:**
- Modify: `site-react/src/PublicSiteApp.tsx`

- [ ] **Step 1: Replace landing data arrays with goal-based content**

In `site-react/src/PublicSiteApp.tsx`, update the landing arrays to this content:

```ts
const features = [
  {
    icon: BookOpen,
    title: "AI Tutor маршрут",
    body: "Сценарий, teaching point, пример ответа, тренировка слотов и финальная проверка слова собирают урок в понятную сессию.",
  },
  {
    icon: MessageCircle,
    title: "Диалоги под цель",
    body: "Поездка, работа, экзамен или разговорная речь: AI ведет роль, исправляет ответ и предлагает модельную фразу.",
  },
  {
    icon: Mic,
    title: "Произношение и shadowing",
    body: "Голосовая тренировка оценивает речь, подсвечивает слабые слова и звуки, дает короткий совет и следующий повтор.",
  },
  {
    icon: Camera,
    title: "Фото и перевод",
    body: "Меню, вывеска, задание или заметка превращаются в перевод, полезную фразу и практику по контексту.",
  },
  {
    icon: WifiOff,
    title: "Словарь, ошибки и offline",
    body: "Learn Words, Review, Spelling, phrasebook, mistakes и offline decks помогают повторять то, что реально проседает.",
  },
  {
    icon: Trophy,
    title: "Прогресс и мотивация",
    body: "XP, streak, уровни, награды, лидеры, ежедневный бонус и история занятий остаются в одном профиле web + Telegram.",
  },
];

const workflow = [
  ["01", "Выберите цель, язык и уровень", "Путешествие, работа, экзамен или разговорная речь становятся маршрутом под ваш уровень A1-C2."],
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
```

- [ ] **Step 2: Update demo screens**

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

- [ ] **Step 3: Update reviews and FAQ**

Replace `testimonials` and `faqs` with:

```ts
const testimonials = [
  {
    name: "Анна",
    role: "маршрут для поездки",
    text: "Я выбрала цель “Путешествия”: разобрала меню по фото, прошла роль в отеле и сохранила фразы, которые точно пригодятся.",
  },
  {
    name: "Марат",
    role: "английский для работы",
    text: "В Telegram удобно тренировать ответы, а в web app я вижу ошибки, словарь и прогресс. Один профиль сильно экономит время.",
  },
  {
    name: "София",
    role: "произношение и shadowing",
    text: "Голосовая тренировка показывает конкретные слабые слова. Я понимаю, что повторить сегодня, а не просто вижу общую оценку.",
  },
];

const faqs = [
  ["Можно заниматься только в Telegram?", "Да. Telegram-бот поддерживает уроки, практику, Premium, лимиты, прогресс и основные учебные режимы. Web app использует тот же профиль."],
  ["Чем Premium и Platinum отличаются от Free?", "Free подходит для знакомства: 5 уроков и 15 сообщений практики в день. Premium расширяет лимиты до 50 уроков, 200 практик и 20 voice. Platinum дает 100 уроков, 500 практик и 60 voice."],
  ["Сохраняются ли ошибки, слова и голосовой прогресс?", "Да. Заметки, phrasebook, mistakes, weak words, pronunciation history, XP, streak и награды закрепляются за аккаунтом."],
  ["Можно ли установить и повторять без стабильной сети?", "Да. Web app работает как PWA, а offline decks помогают повторять сохраненные наборы на телефоне."],
];
```

- [ ] **Step 4: Update hero JSX**

Inside `LandingPage`, update the hero copy and insert goal cards above the language picker:

```tsx
<span className="eyebrow">AI-репетитор под вашу цель в web app, PWA и Telegram</span>
<h1>Выберите цель, Poliglot AI соберет маршрут</h1>
<p>
  Путешествие, работа, экзамен или разговорная речь превращаются в короткий курс: AI Tutor, диалоги, голос,
  фото-перевод, словарь ошибок и прогресс работают в одном профиле.
</p>
<div className="hero-actions">
  <a className="hero-action hero-action--primary" href="/app/">
    Выбрать цель и начать бесплатно <ArrowRight size={18} />
  </a>
  <a className="hero-action hero-action--secondary" href="https://t.me/poliglot_ai_bot">
    Открыть Telegram-бота
  </a>
</div>
<div className="hero-goals" aria-label="Цели обучения">
  {heroGoals.map(([goal, body, badge]) => (
    <a key={goal} href="/app/">
      <strong>{goal}</strong>
      <span>{body}</span>
      <small>{badge}</small>
    </a>
  ))}
</div>
```

Keep the existing language picker after `hero-goals`. Update proof chips to:

```tsx
<div className="hero-proof">
  <span>
    <strong>35</strong> языков интерфейса
  </span>
  <span>
    <strong>A1-C2</strong> уровни обучения
  </span>
  <span>
    <strong>Free</strong> старт в web, PWA и Telegram
  </span>
</div>
```

- [ ] **Step 5: Run TypeScript build and confirm content changes compile**

Run:

```powershell
npm --prefix site-react run build
```

Expected: PASS, and generated files under `Сайт полиглота для бота/` are updated.

- [ ] **Step 6: Commit React content changes**

Run:

```powershell
git add site-react/src/PublicSiteApp.tsx
git commit -m "feat: update landing goal-course copy"
```

## Task 3: Update Course, Pricing, Progress, and Final CTA Content

**Files:**
- Modify: `site-react/src/PublicSiteApp.tsx`

- [ ] **Step 1: Use `courseRoutes` in the course strip**

Replace the inline course card array with `courseRoutes`:

```tsx
<span className="eyebrow">Курсы под цель</span>
<h2>Не просто уроки, а маршрут под вашу ситуацию</h2>
<p>Выберите цель, а Poliglot AI соединит урок, диалог, голос, фото, слова и повторение в понятный сценарий.</p>
```

Use this map:

```tsx
{courseRoutes.map(([title, body, badge], index) => (
  <motion.a key={title} href="/app/" className="course-card" initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} transition={{ duration: 0.42, delay: index * 0.05 }}>
    <span>{badge}</span>
    <h3>{title}</h3>
    <p>{body}</p>
    <b>
      Запустить маршрут <ChevronRight size={16} />
    </b>
  </motion.a>
))}
```

- [ ] **Step 2: Update pricing section**

Replace the pricing heading and plan cards with:

```tsx
<span className="eyebrow">Тарифы</span>
<h2>Начните бесплатно и расширьте лимиты, когда маршрут стал привычкой</h2>
<p className="pricing-lead">Оплата доступна через Telegram Stars, YooKassa/СБП, TON и USDT. Premium и Platinum открывают голос, фото и больше практики.</p>
```

Use these `PlanCard` calls:

```tsx
<PlanCard name="Free" label="Попробовать маршрут" price="0 ₽" body="Старт без оплаты для короткой ежедневной практики." items={["5 уроков в день", "15 сообщений практики", "слова, заметки и базовый прогресс"]} />
<PlanCard name="Premium" label="Регулярная учеба" oldPrice="1000 ₽" price="300 ₽" body="Основной режим для тех, кто занимается каждый день." items={["50 уроков в день", "200 сообщений практики", "20 voice до 30 секунд, голос и фото"]} featured />
<PlanCard name="Platinum" label="Интенсив" oldPrice="2000 ₽" price="590 ₽" body="Максимальные лимиты для поездки, работы или экзамена." items={["100 уроков в день", "500 сообщений практики", "60 voice до 30 секунд, максимум AI-диалогов"]} />
```

Add this payment method strip after `.pricing-grid`:

```tsx
<div className="payment-methods" aria-label="Способы оплаты">
  <span>Telegram Stars</span>
  <span>YooKassa/СБП</span>
  <span>TON</span>
  <span>USDT</span>
</div>
```

- [ ] **Step 3: Update progress and final CTA copy**

Use this progress copy:

```tsx
<span className="eyebrow">Прогресс</span>
<h2>Сервис показывает, что повторить дальше</h2>
<p>Ошибки, заметки, weak words, pronunciation history, XP, streak, награды и лидеры остаются в одном профиле между web app, PWA и Telegram.</p>
```

Use this final CTA copy:

```tsx
<span className="eyebrow">Старт</span>
<h2>Выберите цель и пройдите первый маршрут</h2>
```

Final CTA button:

```tsx
<a className="hero-action hero-action--primary" href="/app/">
  Начать бесплатно <ArrowRight size={18} />
</a>
```

- [ ] **Step 4: Run build**

Run:

```powershell
npm --prefix site-react run build
```

Expected: PASS.

- [ ] **Step 5: Commit course/pricing/progress content**

Run:

```powershell
git add site-react/src/PublicSiteApp.tsx
git commit -m "feat: clarify landing pricing and routes"
```

## Task 4: Update CSS for Full-Screen Animated Hero and Goal Cards

**Files:**
- Modify: `site-react/src/styles.css`

- [ ] **Step 1: Make the hero animation fill the first viewport**

At the bottom of the existing landing hero override section in `site-react/src/styles.css`, add:

```css
.landing-hero {
  min-height: 100svh;
  display: grid;
  align-items: stretch;
}

.landing-hero__inner {
  min-height: calc(100svh - 72px);
  align-items: center;
  padding-top: clamp(96px, 12vh, 132px);
  padding-bottom: clamp(42px, 8vh, 82px);
}

.landing-hero__matter {
  position: absolute;
  inset: -28% -24% -20% -24%;
}

.landing-hero__veil {
  pointer-events: none;
}
```

- [ ] **Step 2: Add goal selector styles**

Add:

```css
.hero-goals {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
  max-width: 820px;
  margin-top: 26px;
}

.hero-goals a {
  display: grid;
  gap: 6px;
  min-height: 118px;
  padding: 15px;
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.82);
  box-shadow: 0 22px 50px rgba(0, 0, 0, 0.2);
  transition:
    transform 180ms ease,
    border-color 180ms ease,
    background 180ms ease;
}

.hero-goals a:hover {
  border-color: rgba(123, 220, 255, 0.42);
  background: rgba(123, 220, 255, 0.12);
  transform: translateY(-3px);
}

.hero-goals strong {
  color: #fff;
  font-size: 1rem;
}

.hero-goals span {
  color: rgba(255, 255, 255, 0.72);
  font-size: 0.82rem;
  line-height: 1.35;
}

.hero-goals small {
  width: fit-content;
  padding: 6px 8px;
  border-radius: 999px;
  background: rgba(245, 210, 122, 0.16);
  color: #f5d27a;
  font-size: 0.74rem;
  font-weight: 900;
}
```

- [ ] **Step 3: Add pricing support styles**

Add:

```css
.pricing-lead {
  max-width: 780px;
  margin: 18px 0 0;
  color: var(--muted);
  font-size: 1.04rem;
  line-height: 1.65;
}

.payment-methods {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 10px;
  max-width: 920px;
  margin: 22px auto 0;
}

.payment-methods span {
  border: 1px solid rgba(45, 91, 255, 0.18);
  border-radius: 999px;
  padding: 8px 11px;
  background: rgba(45, 91, 255, 0.08);
  color: var(--ink);
  font-weight: 850;
}
```

- [ ] **Step 4: Update responsive rules**

Inside the existing `@media (max-width: 980px)` block, add:

```css
.hero-goals {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}
```

Inside the existing `@media (max-width: 680px)` block, add:

```css
.landing-hero {
  min-height: auto;
}

.landing-hero__inner {
  min-height: auto;
  padding-top: 98px;
}

.hero-goals {
  grid-template-columns: 1fr;
}

.hero-goals a {
  min-height: auto;
}
```

- [ ] **Step 5: Run focused E2E test**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "goal-course product flow" --project=desktop-chromium
```

Expected: PASS.

- [ ] **Step 6: Commit CSS changes**

Run:

```powershell
git add site-react/src/styles.css
git commit -m "style: expand animated landing hero"
```

## Task 5: Build Static Landing and Run Full Verification

**Files:**
- Generated: `Сайт полиглота для бота/poliglot-ai.html`
- Generated: `Сайт полиглота для бота/privacy.html`
- Generated: `Сайт полиглота для бота/terms.html`
- Generated: `Сайт полиглота для бота/assets/site-react/*`

- [ ] **Step 1: Build public landing**

Run:

```powershell
npm run site:build
```

Expected: PASS. Vite emits updated `main-*.js` and `main-*.css` files under `Сайт полиглота для бота/assets/site-react/`.

- [ ] **Step 2: Run full public site E2E**

Run:

```powershell
npm --prefix site-react run e2e
```

Expected: PASS for desktop and mobile Chromium projects.

- [ ] **Step 3: Verify in browser MCP**

Open `http://127.0.0.1:4175/poliglot-ai.html` via the preview server from Playwright or `npm --prefix site-react run preview -- --strictPort`. Check:

```text
Desktop 1440x900:
- full-screen animated hero is visible
- H1 and CTA are readable
- goal cards do not overlap demo card
- next section is reachable without layout shift

Mobile 390x844:
- nav, H1, CTA, goal cards, demo tabs fit without text overlap
- hero animation is visible behind the content
- course, pricing, review, FAQ cards stack cleanly
```

- [ ] **Step 4: Inspect final git status**

Run:

```powershell
git status --short --untracked-files=all
```

Expected: only intended source/test/generated static files are changed. `.superpowers/` is ignored.

- [ ] **Step 5: Commit built landing output**

Run:

```powershell
git add site-react/src/PublicSiteApp.tsx site-react/src/styles.css site-react/e2e/public-site.spec.ts "Сайт полиглота для бота/poliglot-ai.html" "Сайт полиглота для бота/privacy.html" "Сайт полиглота для бота/terms.html" "Сайт полиглота для бота/assets/site-react"
git commit -m "feat: rebuild goal-course public landing"
```

- [ ] **Step 6: Push final branch**

Run:

```powershell
git push origin main
```

Expected: push succeeds to `https://github.com/GarryNaxyison/My-app-local.git`.
