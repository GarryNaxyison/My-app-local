# Bold Product Cockpit Landing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rebuild the public React landing as a Bold Product Cockpit page for Poliglot AI while preserving the first-block hero animation, legal links, 35 interface languages, Web app entry, Telegram entry, and pricing contracts.

**Architecture:** Keep legal pages and existing public-site shell in `site-react/src/PublicSiteApp.tsx`, but render a new landing component from `site-react/src/BoldProductLanding.tsx`. Scope the redesign styles in `site-react/src/boldLanding.css` and import them after the legacy stylesheet, so the new landing is authoritative without destabilizing legal pages. Rebuild generated static output only through `npm --prefix site-react run build`.

**Tech Stack:** React 19, Vite, TypeScript, Framer Motion, lucide-react, CSS, Google Fonts (`Unbounded`, `Manrope`), Playwright, local Nx/npm fallback.

---

## File Structure

- Create `site-react/src/BoldProductLanding.tsx`: new public landing component with code-built product mockups, equal Web app and Telegram CTAs, pricing cards, and final CTA.
- Create `site-react/src/boldLanding.css`: scoped Bold Product Cockpit visual system, font imports, responsive layout, reduced-motion handling, and overrides for `.landing-hero`.
- Modify `site-react/src/main.tsx`: import `boldLanding.css` after `styles.css`.
- Modify `site-react/src/PublicSiteApp.tsx`: import `BoldProductLanding` and render it for `page === "landing"`.
- Modify `site-react/e2e/public-site.spec.ts`: update Playwright tests from the previous premium cockpit contract to the new Bold Product Cockpit contract.
- Modify generated files under `Сайт полиглота для бота/` only by running `npm --prefix site-react run build`.
- Do not modify `web-react`, Go backend behavior, payment logic, or `GenerativeArtScene` internals.

## Task 1: Update Landing Contract Tests First

**Files:**
- Modify: `site-react/e2e/public-site.spec.ts`

- [ ] **Step 1: Replace the main landing test**

In `site-react/e2e/public-site.spec.ts`, replace the current test named `landing presents the premium AI tutor cockpit without losing public contracts` with this test:

```ts
test("landing presents the bold product cockpit without losing public contracts", async ({ page }) => {
  test.setTimeout(60_000);
  await page.goto("/poliglot-ai.html");

  await expect(page.locator(".public-nav")).toBeVisible();
  await expect(page.locator(".public-brand__logo img")).toBeVisible();
  await expect(page.locator(".public-nav nav a", { hasText: "Политика" })).toBeVisible();
  await expect(page.locator(".public-nav nav a", { hasText: "Условия" })).toBeVisible();
  await expect(page.locator("[data-site-language-select]")).toBeVisible();

  await expect(page.locator(".bold-landing")).toBeVisible();
  await expect(page.locator(".bold-hero")).toBeVisible();
  await expect(page.locator(".bold-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator(".bold-landing section:not(.bold-hero) canvas")).toHaveCount(0);
  await expect(page.locator(".spark-hero")).toHaveCount(0);

  const rootFont = await page.locator("body").evaluate((node) => getComputedStyle(node).fontFamily);
  expect(rootFont).not.toContain("Inter");
  expect(rootFont).not.toContain("Inrer");

  const webEntry = page.locator('a[data-entry="web-app"]');
  const telegramEntry = page.locator('a[data-entry="telegram"]');
  await expect(webEntry).toHaveCount(3);
  await expect(telegramEntry).toHaveCount(3);
  await expect(webEntry.first()).toHaveAttribute("href", "/app/");
  await expect(telegramEntry.first()).toHaveAttribute("href", "https://t.me/poliglot_ai_bot");

  const ctaShape = await page.locator(".bold-hero .entry-cta").evaluateAll((nodes) =>
    nodes.map((node) => {
      const style = getComputedStyle(node as HTMLElement);
      const rect = (node as HTMLElement).getBoundingClientRect();
      return {
        minHeight: rect.height,
        borderRadius: style.borderRadius,
        background: style.backgroundImage || style.backgroundColor,
        color: style.color,
      };
    }),
  );
  expect(ctaShape).toHaveLength(2);
  expect(ctaShape[0].minHeight).toBeGreaterThanOrEqual(48);
  expect(ctaShape[1].minHeight).toBeGreaterThanOrEqual(48);
  expect(ctaShape[0].background).not.toBe("none");
  expect(ctaShape[1].background).not.toBe("none");
  expect(ctaShape[0].borderRadius).toBe(ctaShape[1].borderRadius);

  await expect(page.locator("h1")).toContainText("AI-репетитором");
  await expect(page.locator(".bold-hero__proof")).toContainText("35");
  await expect(page.locator(".bold-hero__proof")).toContainText("A1-C2");
  await expect(page.locator(".hero-product-tabs button")).toHaveText(["Урок", "Диалог", "Голос", "Фото"]);
  await page.locator(".hero-product-tabs button").nth(1).click();
  await expect(page.locator(".hero-product-output")).toContainText("Could you help me check in?");
  await page.locator(".hero-product-tabs button").nth(2).click();
  await expect(page.locator(".voice-bars i")).toHaveCount(18);

  await expect(page.locator(".daily-step")).toHaveCount(3);
  await expect(page.locator(".module-card")).toHaveCount(4);
  await expect(page.locator(".module-card", { hasText: "Voice Coach" })).toContainText("weak words");
  await expect(page.locator(".module-card", { hasText: "Photo Practice" })).toContainText("No peanuts");
  await expect(page.locator(".memory-node")).toHaveCount(6);
  await expect(page.locator(".entry-panel")).toHaveCount(2);
  await expect(page.locator(".entry-panel", { hasText: "Web app" })).toContainText("dashboard");
  await expect(page.locator(".entry-panel", { hasText: "Telegram" })).toContainText("быстрая практика");

  await expect(page.locator(".plan-card")).toHaveCount(3);
  await expect(page.locator(".plan-card", { hasText: "Free" })).toContainText("0 ₽");
  await expect(page.locator(".plan-card", { hasText: "Premium" })).toContainText("300 ₽");
  await expect(page.locator(".plan-card", { hasText: "Platinum" })).toContainText("590 ₽");
  await expect(page.locator(".payment-methods")).toContainText("Stars");
  await expect(page.locator(".payment-methods")).toContainText("YooKassa");
  await expect(page.locator(".payment-methods")).toContainText("TON");
  await expect(page.locator(".payment-methods")).toContainText("USDT");

  await expect(page.locator(".review-card")).toHaveCount(3);
  await expect(page.locator(".review-stars")).toHaveCount(0);
  await expect(page.locator(".final-cta-section")).toContainText("Web app");
  await expect(page.locator(".final-cta-section")).toContainText("Telegram");
  await expect(page.locator(".site-footer a", { hasText: "Политика" })).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: "Условия" })).toBeVisible();
});
```

- [ ] **Step 2: Update the motion frame-budget test selectors**

In the test named `landing hero keeps animated motion within a bounded frame budget`, change the canvas assertion from `.landing-hero__matter canvas` to `.bold-hero__matter canvas`:

```ts
await expect(page.locator(".bold-hero__matter canvas")).toHaveCount(1);
```

Keep the existing frame request assertions:

```ts
expect(frameRequests).toBeGreaterThan(15);
expect(frameRequests).toBeLessThan(95);
```

- [ ] **Step 3: Replace the mobile readability test**

Replace the current test named `premium cockpit landing stays readable on mobile` with this test:

```ts
test("bold product cockpit landing stays readable on mobile", async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/poliglot-ai.html");

  await expect(page.locator(".bold-hero__matter canvas")).toHaveCount(1);
  await expect(page.locator("h1")).toContainText("AI-репетитором");
  await expect(page.locator(".bold-hero .entry-cta")).toHaveCount(2);
  await expect(page.locator(".hero-product-mockup")).toBeVisible();
  await expect(page.locator(".entry-panel")).toHaveCount(2);
  await expect(page.locator(".site-footer a", { hasText: "Политика" })).toBeVisible();
  await expect(page.locator(".site-footer a", { hasText: "Условия" })).toBeVisible();

  const overflow = await page.evaluate(() => {
    const nodes = Array.from(
      document.querySelectorAll(
        "h1, h2, h3, p, a, button, .bold-hero__proof span, .daily-step, .module-card, .memory-node, .entry-panel, .plan-card, .review-card",
      ),
    );
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

- [ ] **Step 4: Update the Chinese localization selectors**

In the test named `premium cockpit landing localizes generated marketing copy for Chinese`, replace the `selectors` array with:

```ts
const selectors = [
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
```

Keep the existing Cyrillic scan and pricing numeric assertions.

- [ ] **Step 5: Run the focused tests to verify RED**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "bold product cockpit|bounded frame budget|localizes generated marketing copy" --project=desktop-chromium
```

Expected result: FAIL because `.bold-landing`, `.bold-hero`, `.module-card`, `.memory-node`, equal entry CTAs, and the new font stack are not implemented yet.

- [ ] **Step 6: Commit the RED tests**

Run:

```powershell
git add site-react/e2e/public-site.spec.ts
git commit -m "test: specify bold product cockpit landing"
```

Expected result: one commit containing only the Playwright contract update.

## Task 2: Add The Bold Product Landing Component

**Files:**
- Create: `site-react/src/BoldProductLanding.tsx`
- Modify: `site-react/src/PublicSiteApp.tsx`

- [ ] **Step 1: Create `site-react/src/BoldProductLanding.tsx`**

Create the file with this content:

```tsx
import { useState } from "react";
import { motion } from "framer-motion";
import {
  ArrowRight,
  BookOpen,
  Bot,
  BrainCircuit,
  Camera,
  CheckCircle,
  ChevronRight,
  Crown,
  Headphones,
  Laptop,
  MessageCircle,
  Mic,
  Repeat2,
  ScanText,
  ShieldCheck,
  Sparkles,
  Trophy,
  Zap,
} from "lucide-react";
import { GenerativeArtScene } from "@/components/ui/anomalous-matter-hero";

const telegramUrl = "https://t.me/poliglot_ai_bot";

const heroScreens = [
  {
    tab: "Урок",
    label: "AI Tutor lesson",
    output: "I would like to book a table for tonight.",
    note: "Poliglot AI объясняет фразу, проверяет ответ и сохраняет повторение.",
  },
  {
    tab: "Диалог",
    label: "Hotel check-in",
    output: "Could you help me check in? I have a reservation.",
    note: "Roleplay держит ситуацию и предлагает более естественную фразу.",
  },
  {
    tab: "Голос",
    label: "Voice Coach",
    output: "Please speak a little slower.",
    note: "Score, weak words и shadowing показывают, что повторить дальше.",
  },
  {
    tab: "Фото",
    label: "Photo Practice",
    output: "No peanuts, please. How spicy is this dish?",
    note: "Фото превращается в перевод, заметку и practice prompt.",
  },
] as const;

const dailySteps = [
  ["01", "Получите задание", "Цель, уровень и слабые слова собираются в короткий daily route."],
  ["02", "Ответьте как удобно", "Текст, голос, диалог или фото работают как один учебный поток."],
  ["03", "Повторите слабое место", "Ошибка уходит в Mistakes, Notes, Review и offline decks."],
] as const;

const modules = [
  {
    icon: BrainCircuit,
    name: "AI Tutor",
    badge: "lesson + check",
    body: "Ведет от объяснения к ответу, проверке, XP и следующему повторению.",
    sample: "Try: I have been working on my English every day.",
  },
  {
    icon: MessageCircle,
    name: "Roleplay",
    badge: "travel / work / exam",
    body: "Держит реальный сценарий и исправляет фразу без пустого чата.",
    sample: "Could you help me check in?",
  },
  {
    icon: Headphones,
    name: "Voice Coach",
    badge: "score + weak words",
    body: "Показывает pronunciation score, weak words и следующую фразу для shadowing.",
    sample: "Weak words: reservation, slower, tonight",
  },
  {
    icon: ScanText,
    name: "Photo Practice",
    badge: "image -> practice",
    body: "Меню, вывеска или задание становятся переводом, заметкой и упражнением.",
    sample: "No peanuts, please.",
  },
] as const;

const memoryNodes = ["Mistakes", "Weak words", "Notes", "Review", "Spelling", "Offline decks"] as const;

const plans = [
  {
    name: "Free",
    label: "Старт",
    price: "0 ₽",
    body: "Понять маршрут, попробовать базовые уроки, слова, заметки и прогресс.",
    limits: ["5 уроков в день", "15 сообщений практики", "голосовые недоступны"],
  },
  {
    name: "Premium",
    label: "Регулярная учеба",
    oldPrice: "1000 ₽",
    price: "300 ₽",
    body: "AI Tutor, voice, photo practice, roleplay и расширенные дневные лимиты.",
    limits: ["50 уроков в день", "200 сообщений практики", "20 голосовых до 30 секунд"],
    featured: true,
  },
  {
    name: "Platinum",
    label: "Интенсив",
    oldPrice: "2000 ₽",
    price: "590 ₽",
    body: "Максимальный режим для поездки, работы, экзамена или плотной daily practice.",
    limits: ["100 уроков в день", "500 сообщений практики", "60 голосовых до 30 секунд"],
  },
] as const;

const reviews = [
  ["Анна", "готовится к поездке", "Я прогнала check-in, кафе и транспорт. В заметках остались фразы, которые потом пригодились в отеле."],
  ["Марат", "учит английский для работы", "Перед созвоном репетирую self-intro и вопросы по срокам. Ошибки потом повторяю в Telegram."],
  ["София", "тренирует произношение", "Я вижу weak words и сразу повторяю фразу. Это спокойнее, чем просто учить список слов."],
] as const;

export function BoldProductLanding() {
  const [activeTab, setActiveTab] = useState(0);
  const activeScreen = heroScreens[activeTab] ?? heroScreens[0];

  return (
    <main className="bold-landing">
      <section className="bold-hero landing-hero">
        <div className="bold-hero__matter landing-hero__matter" aria-hidden="true">
          <GenerativeArtScene animate color="#00d4ff" particleColor="#f6c84c" />
        </div>
        <div className="bold-hero__veil landing-hero__veil" aria-hidden="true" />

        <div className="bold-hero__inner">
          <motion.div className="bold-hero__copy" initial={{ opacity: 0, y: 24 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.55 }}>
            <span className="bold-kicker">Poliglot AI</span>
            <h1>Говорите с AI-репетитором, который помнит ваши ошибки</h1>
            <p>Урок, диалог, голос, фото-перевод, словарь ошибок и прогресс собираются в один профиль для web app и Telegram.</p>
            <div className="entry-cta-row">
              <a className="entry-cta entry-cta--web" data-entry="web-app" href="/app/">
                Открыть Web app <ArrowRight size={18} />
              </a>
              <a className="entry-cta entry-cta--telegram" data-entry="telegram" href={telegramUrl}>
                Открыть Telegram <Bot size={18} />
              </a>
            </div>
            <div className="bold-hero__proof" aria-label="Ключевые факты Poliglot AI">
              <span>
                <strong>35</strong>
                языков интерфейса
              </span>
              <span>
                <strong>A1-C2</strong>
                маршруты обучения
              </span>
              <span>
                <strong>1 профиль</strong>
                web app + Telegram
              </span>
            </div>
          </motion.div>

          <motion.div className="hero-product-mockup" initial={{ opacity: 0, x: 32 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.6, delay: 0.08 }}>
            <div className="mockup-topline">
              <span>Today session</span>
              <strong>84/100</strong>
            </div>
            <div className="hero-product-tabs" role="tablist" aria-label="Демо Poliglot AI">
              {heroScreens.map((screen, index) => (
                <button key={screen.tab} type="button" className={activeTab === index ? "is-active" : undefined} onClick={() => setActiveTab(index)}>
                  {screen.tab}
                </button>
              ))}
            </div>
            <div className="hero-product-output">
              <span>{activeScreen.label}</span>
              <p>{activeScreen.output}</p>
              <small>{activeScreen.note}</small>
            </div>
            <div className="voice-bars" aria-hidden="true">
              {Array.from({ length: 18 }, (_, index) => (
                <i key={index} style={{ animationDelay: `${index * 0.04}s` }} />
              ))}
            </div>
          </motion.div>
        </div>
      </section>

      <section className="daily-loop-section bold-section">
        <div className="section-heading">
          <span className="bold-kicker">What you do today</span>
          <h2>Не ищите упражнение. Откройте маршрут и сделайте следующий шаг.</h2>
        </div>
        <div className="daily-loop-grid">
          <div className="daily-route-mockup">
            <div className="route-header">
              <BookOpen size={22} />
              <strong>Daily Route</strong>
              <span>12 min</span>
            </div>
            <div className="route-task is-active">
              <Zap size={18} />
              <p>Roleplay: hotel check-in</p>
            </div>
            <div className="route-task">
              <Mic size={18} />
              <p>Repeat weak word: reservation</p>
            </div>
            <div className="route-task">
              <Camera size={18} />
              <p>Photo practice: menu request</p>
            </div>
          </div>
          <div className="daily-steps">
            {dailySteps.map(([number, title, body]) => (
              <article className="daily-step" key={number}>
                <strong>{number}</strong>
                <div>
                  <h3>{title}</h3>
                  <p>{body}</p>
                </div>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section id="features" className="modules-section bold-section">
        <div className="section-heading">
          <span className="bold-kicker">Product modules</span>
          <h2>Каждая функция выглядит как часть одного AI-репетитора</h2>
        </div>
        <div className="module-grid">
          {modules.map((item) => {
            const Icon = item.icon;
            return (
              <article className="module-card" key={item.name}>
                <div className="module-card__top">
                  <Icon size={24} />
                  <span>{item.badge}</span>
                </div>
                <h3>{item.name}</h3>
                <p>{item.body}</p>
                <blockquote>{item.sample}</blockquote>
              </article>
            );
          })}
        </div>
      </section>

      <section className="memory-loop-section bold-section">
        <div className="section-heading">
          <span className="bold-kicker">Mistake memory loop</span>
          <h2>Ошибки не исчезают после урока. Они становятся следующей тренировкой.</h2>
          <p>Poliglot AI связывает mistakes, weak words, notes, review, spelling и offline decks в один цикл повторения.</p>
        </div>
        <div className="memory-board" aria-label="Цикл памяти Poliglot AI">
          {memoryNodes.map((node, index) => (
            <span className="memory-node" key={node}>
              <b>{String(index + 1).padStart(2, "0")}</b>
              {node}
            </span>
          ))}
          <div className="memory-core">
            <Trophy size={34} />
            <strong>XP + streak</strong>
          </div>
        </div>
      </section>

      <section className="entry-section bold-section">
        <div className="section-heading">
          <span className="bold-kicker">Web app + Telegram</span>
          <h2>Два входа. Один профиль. Одна практика.</h2>
        </div>
        <div className="entry-panels">
          <article className="entry-panel">
            <Laptop size={28} />
            <h3>Web app</h3>
            <p>Для focused lessons, dashboard, тарифов, длинных сессий, прогресса и словаря ошибок.</p>
            <a className="entry-cta entry-cta--web" data-entry="web-app" href="/app/">
              Открыть Web app <ChevronRight size={18} />
            </a>
          </article>
          <article className="entry-panel">
            <Bot size={28} />
            <h3>Telegram</h3>
            <p>Для быстрой практики, voice checks, напоминаний и продолжения без потери Premium.</p>
            <a className="entry-cta entry-cta--telegram" data-entry="telegram" href={telegramUrl}>
              Открыть Telegram <ChevronRight size={18} />
            </a>
          </article>
        </div>
      </section>

      <section id="pricing" className="pricing-section bold-pricing">
        <div className="section-heading">
          <span className="bold-kicker">Pricing</span>
          <h2>Начните бесплатно. Premium открывает голос, фото и больше ежедневной практики.</h2>
          <p className="payment-methods">Оплата: Telegram Stars, YooKassa/SBP, TON и USDT.</p>
        </div>
        <div className="pricing-grid">
          {plans.map((plan) => (
            <article className={plan.featured ? "plan-card is-featured" : "plan-card"} key={plan.name}>
              <span className="plan-label">{plan.label}</span>
              <h3>{plan.name}</h3>
              <div className="price-line">
                {"oldPrice" in plan ? <span className="plan-old-price">{plan.oldPrice}</span> : null}
                <strong>{plan.price}</strong>
              </div>
              <p>{plan.body}</p>
              <ul className="plan-included">
                {plan.limits.map((limit) => (
                  <li key={limit}>
                    <CheckCircle size={16} />
                    {limit}
                  </li>
                ))}
              </ul>
              <a href="/app/">Выбрать <ChevronRight size={16} /></a>
            </article>
          ))}
        </div>
      </section>

      <section id="reviews" className="reviews-section bold-section">
        <div className="section-heading">
          <span className="bold-kicker">User stories</span>
          <h2>Короткие истории без фейковых рейтингов</h2>
        </div>
        <div className="reviews-grid">
          {reviews.map(([name, role, text]) => (
            <article className="review-card" key={name}>
              <strong>{name}</strong>
              <small>{role}</small>
              <p>{text}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="final-cta-section">
        <Sparkles size={26} />
        <h2>Откройте Poliglot AI там, где вам удобнее начать</h2>
        <div className="entry-cta-row">
          <a className="entry-cta entry-cta--web" data-entry="web-app" href="/app/">
            Web app <ArrowRight size={18} />
          </a>
          <a className="entry-cta entry-cta--telegram" data-entry="telegram" href={telegramUrl}>
            Telegram <Bot size={18} />
          </a>
        </div>
      </section>
    </main>
  );
}
```

- [ ] **Step 2: Render the new component from `PublicSiteApp.tsx`**

Add this import near the other local imports:

```ts
import { BoldProductLanding } from "./BoldProductLanding";
```

Change the landing render line from:

```tsx
{page === "landing" ? <LandingPage /> : <LegalPage page={page} />}
```

to:

```tsx
{page === "landing" ? <BoldProductLanding /> : <LegalPage page={page} />}
```

- [ ] **Step 3: Run the focused test and confirm it still fails on styling**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "bold product cockpit" --project=desktop-chromium
```

Expected result: FAIL because the new component exists but `boldLanding.css` has not been created or imported yet, so font and visual contract assertions fail.

- [ ] **Step 4: Commit the component wiring**

Run:

```powershell
git add site-react/src/BoldProductLanding.tsx site-react/src/PublicSiteApp.tsx
git commit -m "feat: add bold product landing component"
```

Expected result: one commit containing the new component and landing render wiring.

## Task 3: Add The Bold Product Visual System

**Files:**
- Create: `site-react/src/boldLanding.css`
- Modify: `site-react/src/main.tsx`
- Modify: `site-react/src/styles.css`

- [ ] **Step 1: Import the new stylesheet after the legacy stylesheet**

In `site-react/src/main.tsx`, change:

```ts
import "./styles.css";
```

to:

```ts
import "./styles.css";
import "./boldLanding.css";
```

- [ ] **Step 2: Remove Inter from the root font stack**

In `site-react/src/styles.css`, replace the root font declaration:

```css
font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
```

with:

```css
font-family: Manrope, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
```

- [ ] **Step 3: Create `site-react/src/boldLanding.css`**

Create the file with this content:

```css
@import url("https://fonts.googleapis.com/css2?family=Manrope:wght@400..900&family=Unbounded:wght@600..900&display=swap");

:root {
  --bold-ink: #050711;
  --bold-paper: #f5f8ff;
  --bold-blue: #2458ff;
  --bold-cyan: #00d4ff;
  --bold-gold: #f6c84c;
  --bold-lime: #b7ff4a;
  --bold-line: rgba(5, 7, 17, 0.14);
  font-family: Manrope, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

.bold-landing {
  overflow-x: clip;
  background: var(--bold-paper);
  color: var(--bold-ink);
  font-family: Manrope, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

.bold-landing h1,
.bold-landing h2,
.bold-landing h3,
.bold-kicker,
.entry-cta,
.plan-card h3 {
  font-family: Unbounded, Manrope, ui-sans-serif, system-ui, sans-serif;
  letter-spacing: 0;
}

.bold-hero {
  position: relative;
  min-height: min(980px, calc(100dvh - 0px));
  overflow: hidden;
  background:
    linear-gradient(90deg, rgba(5, 7, 17, 0.98) 0%, rgba(5, 7, 17, 0.84) 48%, rgba(5, 7, 17, 0.5) 100%),
    #050711;
  color: #fff;
}

.bold-hero__matter {
  position: absolute;
  inset: -28% -22% -18% -18%;
  opacity: 0.78;
  mix-blend-mode: screen;
  pointer-events: none;
}

.bold-hero__matter .generative-art-scene,
.bold-hero__matter canvas {
  width: 100% !important;
  height: 100% !important;
}

.bold-hero__veil {
  position: absolute;
  inset: 0;
  background:
    linear-gradient(180deg, rgba(5, 7, 17, 0.08), rgba(5, 7, 17, 0.94)),
    linear-gradient(90deg, rgba(5, 7, 17, 0.98), rgba(5, 7, 17, 0.5));
}

.bold-hero__inner {
  position: relative;
  z-index: 2;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(340px, 500px);
  gap: clamp(28px, 5vw, 76px);
  align-items: center;
  max-width: 1240px;
  min-height: min(980px, 100dvh);
  margin: 0 auto;
  padding: clamp(78px, 10vw, 128px) 24px 64px;
}

.bold-kicker {
  display: inline-flex;
  width: fit-content;
  align-items: center;
  min-height: 32px;
  border: 1px solid rgba(0, 212, 255, 0.32);
  border-radius: 8px;
  padding: 6px 10px;
  background: rgba(0, 212, 255, 0.1);
  color: var(--bold-cyan);
  font-size: 0.76rem;
  font-weight: 800;
  text-transform: uppercase;
}

.bold-hero h1 {
  max-width: 800px;
  margin: 22px 0 0;
  color: #fff;
  font-size: clamp(3.4rem, 7.6vw, 7.4rem);
  line-height: 0.94;
}

.bold-hero__copy > p {
  max-width: 720px;
  margin: 24px 0 0;
  color: rgba(255, 255, 255, 0.78);
  font-size: clamp(1rem, 1.5vw, 1.24rem);
  line-height: 1.7;
}

.entry-cta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 30px;
}

.entry-cta {
  min-height: 52px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 9px;
  border: 1px solid rgba(5, 7, 17, 0.16);
  border-radius: 8px;
  padding: 0 18px;
  color: #050711;
  font-size: 0.9rem;
  font-weight: 800;
  text-decoration: none;
  transition: transform 180ms ease, filter 180ms ease;
}

.entry-cta:hover {
  transform: translateY(-2px);
  filter: saturate(1.08);
}

.entry-cta--web {
  background: linear-gradient(135deg, var(--bold-gold), var(--bold-lime));
}

.entry-cta--telegram {
  background: linear-gradient(135deg, var(--bold-cyan), #91f7ff);
}

.bold-hero__proof {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  max-width: 720px;
  margin-top: 24px;
}

.bold-hero__proof span {
  min-height: 74px;
  display: grid;
  gap: 4px;
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 8px;
  padding: 13px;
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.72);
}

.bold-hero__proof strong {
  color: #fff;
}

.hero-product-mockup,
.daily-route-mockup,
.module-card,
.entry-panel,
.plan-card,
.review-card {
  border: 1px solid var(--bold-line);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 24px 70px rgba(5, 7, 17, 0.12);
}

.hero-product-mockup {
  overflow: hidden;
  padding: 18px;
  border-color: rgba(255, 255, 255, 0.18);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.16), rgba(255, 255, 255, 0.07));
  color: #fff;
  backdrop-filter: blur(18px);
}

.mockup-topline,
.module-card__top,
.route-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.mockup-topline span,
.hero-product-output span {
  color: rgba(255, 255, 255, 0.62);
  font-size: 0.78rem;
  font-weight: 850;
  text-transform: uppercase;
}

.mockup-topline strong {
  border-radius: 8px;
  padding: 7px 10px;
  background: linear-gradient(135deg, var(--bold-gold), var(--bold-cyan));
  color: #050711;
}

.hero-product-tabs {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8px;
  margin-top: 18px;
}

.hero-product-tabs button {
  min-height: 40px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.72);
  cursor: pointer;
  font-weight: 800;
}

.hero-product-tabs button.is-active {
  background: rgba(0, 212, 255, 0.22);
  color: #fff;
}

.hero-product-output {
  display: grid;
  gap: 8px;
  margin-top: 14px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 8px;
  padding: 16px;
  background: rgba(5, 7, 17, 0.46);
}

.hero-product-output p {
  margin: 0;
  color: #fff;
  font-size: 1.1rem;
  line-height: 1.45;
}

.hero-product-output small {
  color: rgba(255, 255, 255, 0.68);
  line-height: 1.5;
}

.voice-bars {
  display: flex;
  align-items: center;
  gap: 5px;
  min-height: 52px;
  margin-top: 14px;
}

.voice-bars i {
  width: 5px;
  height: 20px;
  border-radius: 999px;
  background: linear-gradient(180deg, #fff, var(--bold-cyan));
  animation: boldWave 1.4s ease-in-out infinite;
}

.voice-bars i:nth-child(3n) {
  height: 34px;
}

.voice-bars i:nth-child(4n) {
  height: 42px;
}

.bold-section,
.bold-pricing,
.final-cta-section {
  max-width: 1240px;
  margin: 0 auto;
  padding: clamp(72px, 9vw, 116px) 24px;
}

.section-heading {
  max-width: 880px;
}

.section-heading h2,
.final-cta-section h2 {
  margin: 16px 0 0;
  color: var(--bold-ink);
  font-size: clamp(2.35rem, 5vw, 5.4rem);
  line-height: 1;
}

.section-heading p {
  margin: 18px 0 0;
  color: #4b5870;
  font-size: 1.08rem;
  line-height: 1.7;
}

.daily-loop-grid,
.entry-panels {
  display: grid;
  grid-template-columns: minmax(0, 0.9fr) minmax(0, 1.1fr);
  gap: 18px;
  margin-top: 34px;
}

.daily-route-mockup {
  display: grid;
  align-self: start;
  gap: 12px;
  padding: 20px;
}

.route-header {
  min-height: 46px;
}

.route-header span {
  border-radius: 8px;
  padding: 7px 9px;
  background: rgba(36, 88, 255, 0.1);
  color: var(--bold-blue);
  font-weight: 850;
}

.route-task,
.daily-step {
  display: grid;
  grid-template-columns: 44px minmax(0, 1fr);
  gap: 12px;
  align-items: center;
  border: 1px solid rgba(5, 7, 17, 0.1);
  border-radius: 8px;
  padding: 14px;
  background: #fff;
}

.route-task.is-active {
  background: linear-gradient(135deg, rgba(36, 88, 255, 0.12), rgba(0, 212, 255, 0.12));
}

.route-task p,
.daily-step p {
  margin: 0;
  color: #4b5870;
  line-height: 1.55;
}

.daily-steps {
  display: grid;
  gap: 12px;
}

.daily-step > strong {
  display: grid;
  width: 44px;
  height: 44px;
  place-items: center;
  border-radius: 8px;
  background: var(--bold-ink);
  color: #fff;
}

.daily-step h3,
.module-card h3,
.entry-panel h3,
.review-card strong {
  margin: 0;
}

.module-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-top: 34px;
}

.module-card {
  display: grid;
  gap: 14px;
  min-height: 320px;
  padding: 20px;
}

.module-card__top span {
  border-radius: 8px;
  padding: 7px 9px;
  background: rgba(246, 200, 76, 0.18);
  color: #775600;
  font-size: 0.74rem;
  font-weight: 900;
}

.module-card p,
.entry-panel p,
.plan-card p,
.review-card p {
  margin: 0;
  color: #4b5870;
  line-height: 1.62;
}

.module-card blockquote {
  align-self: end;
  margin: 0;
  border-left: 4px solid var(--bold-blue);
  padding: 12px;
  background: rgba(36, 88, 255, 0.08);
  color: var(--bold-ink);
  font-weight: 800;
}

.memory-loop-section {
  max-width: none;
  background: #050711;
  color: #fff;
}

.memory-loop-section .section-heading,
.memory-board {
  max-width: 1240px;
  margin-inline: auto;
}

.memory-loop-section h2 {
  color: #fff;
}

.memory-loop-section .section-heading p {
  color: rgba(255, 255, 255, 0.72);
}

.memory-board {
  position: relative;
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 12px;
  margin-top: 34px;
}

.memory-node,
.memory-core {
  display: grid;
  gap: 8px;
  min-height: 124px;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.08);
  text-align: center;
  font-weight: 900;
}

.memory-node b {
  color: var(--bold-cyan);
}

.memory-core {
  grid-column: 3 / 5;
  min-height: 150px;
  background: linear-gradient(135deg, var(--bold-blue), #111d46);
}

.entry-panels {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.entry-panel {
  display: grid;
  gap: 16px;
  padding: 24px;
}

.entry-panel .entry-cta {
  width: fit-content;
}

.bold-pricing {
  max-width: none;
  background: linear-gradient(180deg, #eaf1ff, #f8fbff);
}

.bold-pricing > .section-heading,
.bold-pricing > .pricing-grid {
  max-width: 1240px;
  margin-inline: auto;
}

.pricing-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 34px;
}

.plan-card {
  display: grid;
  gap: 16px;
  padding: 24px;
}

.plan-card.is-featured {
  border-color: rgba(36, 88, 255, 0.36);
  box-shadow: 0 28px 86px rgba(36, 88, 255, 0.2);
}

.plan-label {
  width: fit-content;
  border-radius: 8px;
  padding: 7px 9px;
  background: rgba(0, 212, 255, 0.12);
  color: #005a70;
  font-weight: 900;
}

.price-line {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 10px;
}

.price-line strong {
  color: var(--bold-ink);
  font-size: clamp(2.2rem, 4vw, 3.25rem);
}

.plan-old-price {
  color: #8a1f1f;
  font-weight: 900;
  text-decoration: line-through;
  text-decoration-color: #dc2626;
  text-decoration-thickness: 2px;
}

.plan-included {
  display: grid;
  gap: 10px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.plan-included li {
  display: flex;
  align-items: flex-start;
  gap: 9px;
  color: #263246;
  line-height: 1.45;
}

.plan-card > a {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  width: fit-content;
  color: var(--bold-blue);
  font-weight: 900;
}

.payment-methods {
  display: inline-flex;
  width: fit-content;
  max-width: 100%;
  border: 1px solid rgba(36, 88, 255, 0.16);
  border-radius: 8px;
  padding: 10px 12px;
  background: rgba(36, 88, 255, 0.08);
  color: #273755 !important;
  font-weight: 850;
}

.reviews-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 34px;
}

.review-card {
  display: grid;
  gap: 10px;
  min-height: 210px;
  padding: 22px;
}

.review-card small {
  color: #67748c;
}

.final-cta-section {
  display: grid;
  justify-items: center;
  text-align: center;
}

.final-cta-section h2 {
  max-width: 860px;
}

@keyframes boldWave {
  0%,
  100% {
    transform: scaleY(0.62);
    opacity: 0.62;
  }
  50% {
    transform: scaleY(1);
    opacity: 1;
  }
}

@media (prefers-reduced-motion: reduce) {
  .voice-bars i,
  .entry-cta {
    animation: none !important;
    transition: none !important;
  }
}

@media (max-width: 1080px) {
  .bold-hero__inner,
  .daily-loop-grid,
  .entry-panels {
    grid-template-columns: 1fr;
  }

  .module-grid,
  .memory-board,
  .pricing-grid,
  .reviews-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .memory-core {
    grid-column: 1 / -1;
  }
}

@media (max-width: 680px) {
  .bold-hero {
    min-height: auto;
  }

  .bold-hero__inner {
    min-height: auto;
    padding-top: 78px;
  }

  .bold-hero h1 {
    font-size: clamp(2.55rem, 13vw, 4rem);
    line-height: 1.02;
  }

  .entry-cta-row,
  .entry-cta {
    width: 100%;
  }

  .bold-hero__proof,
  .hero-product-tabs,
  .module-grid,
  .memory-board,
  .pricing-grid,
  .reviews-grid {
    grid-template-columns: 1fr;
  }

  .bold-section,
  .bold-pricing,
  .final-cta-section {
    padding: 64px 18px;
  }

  .section-heading h2,
  .final-cta-section h2 {
    font-size: clamp(2rem, 11vw, 3.1rem);
  }
}
```

- [ ] **Step 4: Run the focused tests to verify GREEN for the redesigned contract**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "bold product cockpit|bounded frame budget" --project=desktop-chromium
```

Expected result: PASS for `landing presents the bold product cockpit without losing public contracts`, `landing hero keeps animated motion within a bounded frame budget`, and `bold product cockpit landing stays readable on mobile`.

- [ ] **Step 5: Commit the visual system**

Run:

```powershell
git add site-react/src/main.tsx site-react/src/styles.css site-react/src/boldLanding.css
git commit -m "style: add bold product cockpit visual system"
```

Expected result: one commit containing the scoped CSS, font stack update, and CSS import.

## Task 4: Rebuild Public Translations And Static Site

**Files:**
- Modify: `site-react/public/assets/site-phrases.js`
- Modify: `Сайт полиглота для бота/poliglot-ai.html`
- Modify: `Сайт полиглота для бота/privacy.html`
- Modify: `Сайт полиглота для бота/terms.html`
- Modify: `Сайт полиглота для бота/assets/site-react/*`

- [ ] **Step 1: Rebuild public-site translation phrases**

Run:

```powershell
node tools/rebuild_public_site_translations.mjs
```

Expected result: command exits with code 0 and updates `site-react/public/assets/site-phrases.js` so the new Bold Product copy can be localized by the existing public-site i18n runtime.

- [ ] **Step 2: Run the Chinese localization test**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "localizes generated marketing copy for Chinese" --project=desktop-chromium
```

Expected result: PASS and no Cyrillic fragments returned by the test.

- [ ] **Step 3: Build the public React site**

Run:

```powershell
npm --prefix site-react run build
```

Expected result: command exits with code 0 and writes generated output to `Сайт полиглота для бота/`.

- [ ] **Step 4: Commit generated public-site changes**

Run:

```powershell
git add site-react/public/assets/site-phrases.js "Сайт полиглота для бота"
git commit -m "build: refresh bold product public site"
```

Expected result: one commit containing public translation phrases and generated static output.

## Task 5: Full Verification, Browser Check, And Push

**Files:**
- Verify only unless a command exposes a concrete defect.

- [ ] **Step 1: Run TypeScript and Vite build**

Run:

```powershell
npm --prefix site-react run build
```

Expected result: exit code 0.

- [ ] **Step 2: Run the full public-site Playwright suite**

Run:

```powershell
npm --prefix site-react run e2e -- --project=desktop-chromium
```

Expected result: all public-site desktop tests pass.

- [ ] **Step 3: Run the mobile project for the landing tests**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "bold product cockpit|language selector" --project=mobile-chromium
```

Expected result: all selected mobile tests pass.

- [ ] **Step 4: Start preview server for browser inspection**

Run:

```powershell
Start-Process -FilePath "C:\Program Files\nodejs\npm.cmd" -ArgumentList @("--prefix","site-react","run","preview","--","--strictPort") -WorkingDirectory "E:\PROJECTS\New project\My app local" -WindowStyle Hidden
```

Expected result: `http://127.0.0.1:4175/poliglot-ai.html` serves the rebuilt landing.

- [ ] **Step 5: Inspect desktop and mobile in the in-app browser or Playwright**

Check:

- desktop `1440x900`: hero animation visible only in first block, equal Web app and Telegram CTAs, no text overlap;
- mobile `390x844`: no horizontal scroll, CTAs are touch-safe, product mockup fits, pricing cards do not overflow;
- first section after hero has no `canvas`;
- computed body font family contains `Manrope` and not `Inter`.

Expected result: visual checks pass. If a check fails, patch `site-react/src/boldLanding.css`, rerun Steps 1-3, and include the fix in a new commit.

- [ ] **Step 6: Verify git diff scope before push**

Run:

```powershell
git status --short
git log --oneline -5
```

Expected result: only intended landing, test, translation, and generated public-site changes are present for this feature. Pre-existing unrelated working-tree changes are not staged.

- [ ] **Step 7: Push the current branch**

Run:

```powershell
git push origin codex/ai-tutor-rebuild-fix
```

Expected result: push succeeds to `https://github.com/GarryNaxyison/My-app-local.git`.

## Self-Review Checklist

- Spec coverage: Tasks 1-5 cover Bold Product Cockpit, first-block-only animation, code-built product mockups, equal Web app and Telegram CTAs, no Inter/Inrer, legal links, pricing, localization, build output, browser verification, commit, and push.
- Scope: The plan changes only the public landing, public-site tests, public-site translations, generated static public site, and landing CSS. It does not change `web-react`, Go backend, payment logic, legal source content, or `GenerativeArtScene`.
- Test-first flow: Task 1 writes RED Playwright tests before Task 2 and Task 3 implement the new landing.
- Risk control: New visual styles are scoped in `boldLanding.css`, imported after legacy styles, so legal pages keep existing styling while the landing gets an authoritative visual layer.
