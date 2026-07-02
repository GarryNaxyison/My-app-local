# Static SEO Pages Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a bilingual static SEO page cluster for NERIVA that positions the product as a full web app first, with Telegram as a companion channel.

**Architecture:** Add generated static HTML pages under the public-site source, driven by a single data module. The generator owns page HTML, JSON-LD, language alternates, and sitemap entries so the ten canonical SEO pages stay consistent. The React landing only receives a compact footer/guides link area.

**Tech Stack:** Vite public site, Node.js ESM generator scripts, static HTML/CSS, Playwright e2e, existing `site-react` build and preview workflow.

---

## File Structure

- Create `site-react/e2e/static-seo-pages.spec.ts` for failing TDD coverage of SEO pages, metadata, sitemap, and landing guide links.
- Create `site-react/scripts/static-seo-pages-data.mjs` as the single content source for page slugs, localized copy, FAQ, internal links, social links, legal links, and canonical host constants.
- Create `site-react/scripts/generate-static-seo-pages.mjs` to render HTML files and `site-react/public/sitemap.xml`.
- Create generated static pages:
  - `site-react/public/ai-english-tutor.html`
  - `site-react/public/english-speaking-practice.html`
  - `site-react/public/english-pronunciation-trainer.html`
  - `site-react/public/english-for-work-and-travel.html`
  - `site-react/public/language-learning-web-app.html`
  - `site-react/public/en/ai-english-tutor.html`
  - `site-react/public/en/english-speaking-practice.html`
  - `site-react/public/en/english-pronunciation-trainer.html`
  - `site-react/public/en/english-for-work-and-travel.html`
  - `site-react/public/en/language-learning-web-app.html`
- Create `site-react/public/assets/seo-pages.css` for shared static SEO page styling.
- Modify `site-react/package.json` to run the SEO generator before the public site build.
- Modify `site-react/src/landingSeoContent.ts` to export the SEO guide links for the landing.
- Modify `site-react/src/EnglishSparkLanding.tsx` to render a compact guides area in the footer only.
- Modify `site-react/e2e/public-site.spec.ts` only if needed to assert that the main landing exposes guide links without adding hero/nav clutter.

## Task 1: Add Failing SEO Page Coverage

**Files:**
- Create: `site-react/e2e/static-seo-pages.spec.ts`

- [ ] **Step 1: Write failing Playwright tests**

Create `site-react/e2e/static-seo-pages.spec.ts`:

```ts
import { expect, test } from "@playwright/test";

const pages = [
  {
    slug: "ai-english-tutor",
    ruPath: "/ai-english-tutor.html",
    enPath: "/en/ai-english-tutor.html",
    ruTitle: "AI-репетитор английского онлайн - NERIVA",
    enTitle: "AI English Tutor Online - NERIVA",
    ruH1: "AI-репетитор английского в веб-приложении NERIVA",
    enH1: "AI English tutor inside the NERIVA web app",
  },
  {
    slug: "english-speaking-practice",
    ruPath: "/english-speaking-practice.html",
    enPath: "/en/english-speaking-practice.html",
    ruTitle: "Разговорная практика английского с AI - NERIVA",
    enTitle: "English Speaking Practice With AI - NERIVA",
    ruH1: "Разговорная практика английского с AI",
    enH1: "English speaking practice with AI",
  },
  {
    slug: "english-pronunciation-trainer",
    ruPath: "/english-pronunciation-trainer.html",
    enPath: "/en/english-pronunciation-trainer.html",
    ruTitle: "Тренажер произношения английского с AI - NERIVA",
    enTitle: "English Pronunciation Trainer With AI - NERIVA",
    ruH1: "Тренажер произношения английского с AI",
    enH1: "English pronunciation trainer with AI",
  },
  {
    slug: "english-for-work-and-travel",
    ruPath: "/english-for-work-and-travel.html",
    enPath: "/en/english-for-work-and-travel.html",
    ruTitle: "Английский для работы и путешествий - NERIVA",
    enTitle: "English for Work and Travel - NERIVA",
    ruH1: "Английский для работы и путешествий",
    enH1: "English for work and travel",
  },
  {
    slug: "language-learning-web-app",
    ruPath: "/language-learning-web-app.html",
    enPath: "/en/language-learning-web-app.html",
    ruTitle: "Веб-приложение для изучения языков - NERIVA",
    enTitle: "Language Learning Web App - NERIVA",
    ruH1: "Веб-приложение для изучения языков NERIVA",
    enH1: "NERIVA language learning web app",
  },
] as const;

test.describe("static SEO pages", () => {
  test("serves Russian static SEO pages with real HTML content and metadata", async ({ page }) => {
    for (const item of pages) {
      const response = await page.goto(item.ruPath);
      expect(response?.ok(), `${item.ruPath} should be served`).toBe(true);
      await expect(page).toHaveTitle(item.ruTitle);
      await expect(page.locator("html")).toHaveAttribute("lang", "ru");
      await expect(page.locator("h1")).toHaveText(item.ruH1);
      await expect(page.locator("#root")).toHaveCount(0);
      await expect(page.locator('meta[name="description"]')).toHaveAttribute("content", /NERIVA/);
      await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", `https://poliglotai.ru${item.ruPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="ru"]')).toHaveAttribute("href", `https://poliglotai.ru${item.ruPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="en"]')).toHaveAttribute("href", `https://poliglotai.online${item.enPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="x-default"]')).toHaveAttribute("href", `https://poliglotai.online${item.enPath}`);
      await expect(page.locator(".seo-faq__item")).toHaveCount(5);
      await expect(page.locator(".seo-related a")).toHaveCount(4);
      await expect(page.locator('a[data-entry="web-app"]').first()).toHaveAttribute("href", "/app/");

      const jsonLdText = await page.locator('script[type="application/ld+json"]').textContent();
      expect(jsonLdText).toBeTruthy();
      const jsonLd = JSON.parse(jsonLdText || "{}");
      const types = jsonLd["@graph"].map((node: { "@type": string }) => node["@type"]);
      expect(types).toEqual(expect.arrayContaining(["Organization", "WebSite", "SoftwareApplication", "WebPage", "BreadcrumbList", "FAQPage"]));
      const faqNode = jsonLd["@graph"].find((node: { "@type": string }) => node["@type"] === "FAQPage");
      expect(faqNode.mainEntity).toHaveLength(5);
    }
  });

  test("serves English static SEO pages with international canonical metadata", async ({ page }) => {
    for (const item of pages) {
      const response = await page.goto(item.enPath);
      expect(response?.ok(), `${item.enPath} should be served`).toBe(true);
      await expect(page).toHaveTitle(item.enTitle);
      await expect(page.locator("html")).toHaveAttribute("lang", "en");
      await expect(page.locator("h1")).toHaveText(item.enH1);
      await expect(page.locator("#root")).toHaveCount(0);
      await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", `https://poliglotai.online${item.enPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="ru"]')).toHaveAttribute("href", `https://poliglotai.ru${item.ruPath}`);
      await expect(page.locator('link[rel="alternate"][hreflang="en"]')).toHaveAttribute("href", `https://poliglotai.online${item.enPath}`);
      await expect(page.locator(".seo-page")).toContainText("web app");
      await expect(page.locator(".seo-page")).toContainText("Telegram");
    }
  });

  test("sitemap exposes the bilingual SEO page cluster with alternates", async ({ page }) => {
    const response = await page.goto("/sitemap.xml");
    expect(response?.ok()).toBe(true);
    const body = await response?.text();
    expect(body).toBeTruthy();

    for (const item of pages) {
      expect(body).toContain(`https://poliglotai.ru${item.ruPath}`);
      expect(body).toContain(`https://poliglotai.online${item.enPath}`);
      expect(body).toContain(`hreflang="ru" href="https://poliglotai.ru${item.ruPath}"`);
      expect(body).toContain(`hreflang="en" href="https://poliglotai.online${item.enPath}"`);
      expect(body).toContain(`hreflang="x-default" href="https://poliglotai.online${item.enPath}"`);
    }
  });

  test("landing exposes SEO guides only outside the primary hero and navigation", async ({ page }) => {
    await page.goto("/poliglot-ai.html?lang=en");

    const guideLinks = page.locator(".seo-guides a");
    await expect(guideLinks).toHaveCount(5);
    await expect(page.locator(".landing-hero .seo-guides")).toHaveCount(0);
    await expect(page.locator(".public-nav .seo-guides")).toHaveCount(0);
    await expect(guideLinks.first()).toHaveAttribute("href", "/en/ai-english-tutor.html");

    await page.goto("/poliglot-ai.html?lang=ru");
    const ruGuideLinks = page.locator(".seo-guides a");
    await expect(ruGuideLinks).toHaveCount(5);
    await expect(ruGuideLinks.first()).toHaveAttribute("href", "/ai-english-tutor.html");
  });
});
```

- [ ] **Step 2: Run the test and verify RED**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "static SEO pages"
```

Expected: FAIL because `/ai-english-tutor.html` and `/en/ai-english-tutor.html` are not generated yet, and `.seo-guides` does not exist on the landing.

## Task 2: Add Static SEO Content Source

**Files:**
- Create: `site-react/scripts/static-seo-pages-data.mjs`

- [ ] **Step 1: Create the shared page data module**

Create `site-react/scripts/static-seo-pages-data.mjs` with:

```js
export const russianOrigin = "https://poliglotai.ru";
export const englishOrigin = "https://poliglotai.online";

export const socialProfileUrls = [
  "https://www.youtube.com/@NERIVA",
  "https://www.instagram.com/poliglotai.online/",
  "https://www.tiktok.com/@poliglotai.online",
];

export const legalLinks = {
  ru: [
    { href: "/privacy.html", label: "Политика обработки данных" },
    { href: "/terms.html", label: "Условия использования" },
    { href: "/agreement.html", label: "Пользовательское соглашение" },
    { href: "/consent.html", label: "Согласие на обработку данных" },
  ],
  en: [
    { href: "/privacy.html?lang=en", label: "Privacy" },
    { href: "/terms.html?lang=en", label: "Terms" },
    { href: "/agreement.html?lang=en", label: "User Agreement" },
    { href: "/consent.html?lang=en", label: "Data Consent" },
  ],
};

export const seoPages = [
  {
    slug: "ai-english-tutor",
    ru: {
      path: "/ai-english-tutor.html",
      title: "AI-репетитор английского онлайн - NERIVA",
      description: "NERIVA - веб-приложение с AI-репетитором английского: уроки, практика, исправление ошибок, повторение и прогресс в одном профиле.",
      h1: "AI-репетитор английского в веб-приложении NERIVA",
      intro: "NERIVA помогает заниматься английским в полноценном web app: вы проходите короткие уроки, отвечаете на задания, видите исправления и возвращаетесь к слабым местам.",
      sections: [
        ["Как проходит занятие", "Урок начинается с цели и фразы в контексте. Затем web app просит ответить, проверяет смысл, показывает более естественную версию и сохраняет результат в прогрессе."],
        ["Почему это не просто чат", "NERIVA связывает уроки, roleplay, словарь ошибок, заметки, XP и повторение. Telegram можно использовать для быстрого доступа, но основная учебная среда находится в web app."],
        ["Кому подходит", "Формат полезен тем, кто хочет регулярно практиковать английский без расписания с преподавателем и видеть, что именно нужно повторить дальше."],
      ],
      faq: [
        ["Что такое AI-репетитор NERIVA?", "Это веб-приложение для изучения английского, где AI ведет урок, проверяет ответ, объясняет ошибку и сохраняет прогресс."],
        ["Можно ли заниматься бесплатно?", "Да. Free-режим дает стартовые уроки и практику, а Premium и Platinum открывают больше лимитов и голосовые инструменты."],
        ["Нужен ли Telegram?", "Нет. Основной сценарий работает в web app. Telegram остается дополнительным каналом для быстрого запуска и уведомлений."],
        ["Подходит ли для начинающих?", "Да. Маршрут можно использовать для базовой практики, коротких ответов, словаря и повторения ошибок."],
        ["Заменяет ли это живого преподавателя?", "Нет. NERIVA помогает заниматься самостоятельно и чаще практиковаться, но не обещает заменить все задачи живого преподавателя."],
      ],
    },
    en: {
      path: "/en/ai-english-tutor.html",
      title: "AI English Tutor Online - NERIVA",
      description: "NERIVA is a web app with an AI English tutor for lessons, speaking practice, corrections, review, and progress in one profile.",
      h1: "AI English tutor inside the NERIVA web app",
      intro: "NERIVA helps you study English in a full web app: take short lessons, answer prompts, review corrections, and return to weak spots.",
      sections: [
        ["How a session works", "A lesson starts with a goal and a phrase in context. The web app asks for your answer, checks meaning, shows a more natural version, and saves progress."],
        ["More than a chat", "NERIVA connects lessons, roleplay, mistake review, notes, XP, and repetition. Telegram can help with quick access, while the web app remains the main study workspace."],
        ["Who it helps", "The format is useful when you want regular English practice without scheduling a human tutor and need to see what to review next."],
      ],
      faq: [
        ["What is the NERIVA tutor?", "It is a language learning web app where AI guides a lesson, checks your answer, explains mistakes, and saves progress."],
        ["Can I start for free?", "Yes. Free gives starter lessons and practice, while Premium and Platinum unlock higher limits and voice tools."],
        ["Do I need Telegram?", "No. The main experience works in the web app. Telegram is an optional companion channel for quick starts and notifications."],
        ["Is it good for beginners?", "Yes. You can use it for basic practice, short answers, vocabulary, and mistake review."],
        ["Does it replace a human teacher?", "No. NERIVA helps you practice independently and more often, but it does not claim to replace every role of a human teacher."],
      ],
    },
  },
  {
    slug: "english-speaking-practice",
    ru: {
      path: "/english-speaking-practice.html",
      title: "Разговорная практика английского с AI - NERIVA",
      description: "Разговорная практика английского в NERIVA: roleplay, свободные ответы, исправление фраз, заметки, ошибки и повторение в web app.",
      h1: "Разговорная практика английского с AI",
      intro: "NERIVA помогает говорить чаще: вы тренируете короткие диалоги, рабочие и бытовые ситуации, получаете исправления и сохраняете полезные фразы.",
      sections: [
        ["Практика без расписания", "Можно открыть web app и начать короткую сессию в любое время: ответить на вопрос, продолжить roleplay или разобрать фразу."],
        ["Исправления по делу", "AI показывает, где фраза звучит неестественно, предлагает более живой вариант и возвращает слабое место в повторение."],
        ["Один профиль", "Диалоги, заметки, ошибки и прогресс остаются в одном аккаунте. Telegram можно использовать как быстрый дополнительный вход к той же практике."],
      ],
      faq: [
        ["Можно ли практиковать speaking без преподавателя?", "Да. NERIVA дает сценарии, вопросы и исправления, чтобы вы чаще строили ответы на английском."],
        ["Какие темы есть для разговоров?", "Подходят поездки, работа, everyday small talk, экзамены, кафе, отель, аэропорт и другие практические ситуации."],
        ["Сохраняются ли ошибки?", "Да. Слабые места попадают в ошибки, заметки и повторение, чтобы к ним можно было вернуться."],
        ["Можно ли заниматься с телефона?", "Да. Web app работает на телефоне, а Telegram можно использовать для быстрого доступа к коротким задачам."],
        ["Это подходит для ежедневной практики?", "Да. Короткие сессии специально рассчитаны на регулярную практику без длинной подготовки."],
      ],
    },
    en: {
      path: "/en/english-speaking-practice.html",
      title: "English Speaking Practice With AI - NERIVA",
      description: "Practice English speaking with NERIVA: roleplay, free answers, corrections, notes, mistakes, and review inside a web app.",
      h1: "English speaking practice with AI",
      intro: "NERIVA helps you speak more often: practice short dialogues, work and everyday situations, get corrections, and save useful phrases.",
      sections: [
        ["Practice without scheduling", "Open the web app and start a short session whenever you have time: answer a prompt, continue roleplay, or repair a phrase."],
        ["Useful corrections", "AI shows where a phrase sounds unnatural, suggests a better version, and brings the weak spot back for review."],
        ["One profile", "Dialogues, notes, mistakes, and progress stay in one account. Telegram can be used as a quick companion entry to the same practice."],
      ],
      faq: [
        ["Can I practice speaking without a teacher?", "Yes. NERIVA gives scenarios, questions, and corrections so you build English answers more often."],
        ["What topics can I practice?", "Travel, work, everyday small talk, exams, cafes, hotels, airports, and other practical situations are a good fit."],
        ["Are mistakes saved?", "Yes. Weak spots go into mistakes, notes, and review so you can return to them."],
        ["Can I practice on my phone?", "Yes. The web app works on mobile, and Telegram can help with quick short tasks."],
        ["Is it useful for daily practice?", "Yes. Short sessions are designed for regular practice without long preparation."],
      ],
    },
  },
  {
    slug: "english-pronunciation-trainer",
    ru: {
      path: "/english-pronunciation-trainer.html",
      title: "Тренажер произношения английского с AI - NERIVA",
      description: "Тренируйте произношение английского в NERIVA: voice practice, shadowing, слабые слова, повторение и прогресс в web app.",
      h1: "Тренажер произношения английского с AI",
      intro: "NERIVA помогает тренировать речь через голосовые ответы, shadowing, слабые слова и повторение фраз в контексте.",
      sections: [
        ["Голосовая практика", "Вы произносите фразу, получаете расшифровку, видите слабые слова и повторяете более естественную версию."],
        ["Shadowing и listening", "Короткие аудиоциклы помогают услышать фразу, повторить ее и закрепить произношение без отдельного приложения."],
        ["История прогресса", "Результаты, слабые слова и заметки остаются в web app, чтобы следующая тренировка начиналась с реальной проблемы."],
      ],
      faq: [
        ["Что тренирует Voice Coach?", "Он помогает повторять фразы, замечать слабые слова и возвращаться к проблемным местам."],
        ["Можно ли использовать shadowing?", "Да. Shadowing помогает услышать модельную фразу, повторить ее и сравнить результат."],
        ["Нужен ли микрофон?", "Для голосовой практики нужен микрофон. Текстовые уроки и заметки доступны без него."],
        ["Сохраняется ли история?", "Да. Прогресс и слабые места остаются в профиле web app."],
        ["Подходит ли для speaking?", "Да. Произношение связано с speaking-практикой, roleplay и повторением естественных фраз."],
      ],
    },
    en: {
      path: "/en/english-pronunciation-trainer.html",
      title: "English Pronunciation Trainer With AI - NERIVA",
      description: "Train English pronunciation with NERIVA: voice practice, shadowing, weak words, review, and progress inside a web app.",
      h1: "English pronunciation trainer with AI",
      intro: "NERIVA helps train speech through voice answers, shadowing, weak words, and phrase repetition in context.",
      sections: [
        ["Voice practice", "Say a phrase, get a transcript, see weak words, and repeat a more natural version."],
        ["Shadowing and listening", "Short audio loops help you hear a phrase, repeat it, and reinforce pronunciation without a separate app."],
        ["Progress history", "Results, weak words, and notes stay in the web app, so the next session starts from a real issue."],
      ],
      faq: [
        ["What does Voice Coach train?", "It helps you repeat phrases, notice weak words, and return to problem spots."],
        ["Can I use shadowing?", "Yes. Shadowing helps you hear a model phrase, repeat it, and compare the result."],
        ["Do I need a microphone?", "Voice practice needs a microphone. Text lessons and notes work without it."],
        ["Is history saved?", "Yes. Progress and weak spots stay in your web app profile."],
        ["Does this help speaking?", "Yes. Pronunciation connects with speaking practice, roleplay, and natural phrase review."],
      ],
    },
  },
  {
    slug: "english-for-work-and-travel",
    ru: {
      path: "/english-for-work-and-travel.html",
      title: "Английский для работы и путешествий - NERIVA",
      description: "NERIVA помогает готовить английский для работы и путешествий: звонки, письма, отель, аэропорт, кафе, roleplay и заметки.",
      h1: "Английский для работы и путешествий",
      intro: "NERIVA превращает рабочие и туристические ситуации в короткие тренировки внутри web app: фразы, roleplay, фото-перевод и повторение.",
      sections: [
        ["Рабочие сценарии", "Можно тренировать созвоны, письма, self-intro, уточнение сроков и деловые фразы, которые нужны до реальной встречи."],
        ["Поездки без паники", "Отель, аэропорт, кафе, транспорт и врач превращаются в практические диалоги и заметки для быстрого повторения."],
        ["Фото и контекст", "Меню, вывеска или задание могут стать переводом, заметкой и короткой практикой по тому же контексту."],
      ],
      faq: [
        ["Какие рабочие темы можно тренировать?", "Созвоны, письма, самопрезентацию, small talk, дедлайны и уточняющие вопросы."],
        ["Поможет ли перед поездкой?", "Да. Можно отработать отель, кафе, аэропорт, транспорт и другие бытовые ситуации."],
        ["Можно ли переводить фото?", "Да. Фото меню, вывески или задания можно превратить в перевод и практику."],
        ["Сохраняются ли фразы?", "Да. Полезные фразы можно держать в заметках и возвращать в повторение."],
        ["Telegram обязателен?", "Нет. Основная работа идет в web app, Telegram остается дополнительным быстрым каналом."],
      ],
    },
    en: {
      path: "/en/english-for-work-and-travel.html",
      title: "English for Work and Travel - NERIVA",
      description: "Use NERIVA to prepare English for work and travel: calls, emails, hotels, airports, cafes, roleplay, and notes.",
      h1: "English for work and travel",
      intro: "NERIVA turns work and travel situations into short web app practice: phrases, roleplay, photo translation, and review.",
      sections: [
        ["Work scenarios", "Practice calls, emails, self-intros, deadline questions, and useful business phrases before the real meeting."],
        ["Travel without panic", "Hotels, airports, cafes, transport, and doctor visits become practical dialogues and notes for quick review."],
        ["Photos and context", "A menu, sign, or exercise can become a translation, note, and short practice in the same context."],
      ],
      faq: [
        ["What work topics can I train?", "Calls, emails, self-intros, small talk, deadlines, and clarifying questions."],
        ["Can it help before a trip?", "Yes. You can rehearse hotels, cafes, airports, transport, and other daily situations."],
        ["Can I translate photos?", "Yes. A photo of a menu, sign, or exercise can become a translation and practice prompt."],
        ["Are phrases saved?", "Yes. Useful phrases can stay in notes and return for review."],
        ["Is Telegram required?", "No. The main workflow is in the web app, while Telegram remains an optional quick channel."],
      ],
    },
  },
  {
    slug: "language-learning-web-app",
    ru: {
      path: "/language-learning-web-app.html",
      title: "Веб-приложение для изучения языков - NERIVA",
      description: "NERIVA - web app для изучения языков: AI-уроки, speaking, произношение, словарь ошибок, заметки, прогресс и тарифы.",
      h1: "Веб-приложение для изучения языков NERIVA",
      intro: "NERIVA - полноценное web app для изучения языков с AI-уроками, speaking-практикой, произношением, ошибками, заметками и прогрессом.",
      sections: [
        ["Учебное пространство", "В одном профиле собраны уроки, roleplay, словарь ошибок, заметки, прогресс, тарифы и ежедневные лимиты."],
        ["Web app как основной продукт", "Большой экран удобен для длинных сессий, просмотра прогресса, тарифов и управления профилем. Telegram дополняет этот сценарий короткой практикой."],
        ["Для регулярной привычки", "Короткие задания, XP, streak, повторение и история занятий помогают возвращаться к языку чаще."],
      ],
      faq: [
        ["Что входит в web app?", "Уроки, практика, произношение, заметки, ошибки, прогресс, тарифы и профиль пользователя."],
        ["Чем web app отличается от Telegram?", "Web app удобнее для длинных сессий, прогресса, тарифов и управления профилем. Telegram быстрее для коротких задач."],
        ["Можно ли учить не только английский?", "Интерфейс поддерживает несколько языков, а основной фокус первой SEO-волны - английский."],
        ["Есть ли подписка?", "Да. Есть Free, Premium и Platinum с разными дневными лимитами и доступом к голосовым и фото-инструментам."],
        ["Нужно ли устанавливать приложение?", "Можно пользоваться через браузер. На телефоне web app можно открыть как мобильную страницу или PWA-сценарий."],
      ],
    },
    en: {
      path: "/en/language-learning-web-app.html",
      title: "Language Learning Web App - NERIVA",
      description: "NERIVA is a language learning web app with AI lessons, speaking, pronunciation, mistakes, notes, progress, and plans.",
      h1: "NERIVA language learning web app",
      intro: "NERIVA is a full language learning web app with AI lessons, speaking practice, pronunciation, mistakes, notes, and progress.",
      sections: [
        ["Study workspace", "One profile brings together lessons, roleplay, mistake review, notes, progress, plans, and daily limits."],
        ["Web app first", "A larger screen helps with longer sessions, progress, pricing, and profile control. Telegram adds quick companion practice."],
        ["Built for a regular habit", "Short tasks, XP, streak, review, and lesson history help you return to language practice more often."],
      ],
      faq: [
        ["What is included in the web app?", "Lessons, practice, pronunciation, notes, mistakes, progress, plans, and user profile tools."],
        ["How is the web app different from Telegram?", "The web app is better for longer sessions, progress, pricing, and profile control. Telegram is faster for short tasks."],
        ["Can I learn languages beyond English?", "The interface supports multiple languages, while the first SEO wave focuses on English search demand."],
        ["Are there paid plans?", "Yes. Free, Premium, and Platinum provide different daily limits and access to voice and photo tools."],
        ["Do I need to install an app?", "You can use it in a browser. On mobile, the web app works as a mobile page or PWA-style experience."],
      ],
    },
  },
];
```

- [ ] **Step 2: Run a syntax check**

Run:

```powershell
node --check site-react/scripts/static-seo-pages-data.mjs
```

Expected: PASS with no syntax errors.

## Task 3: Implement the Static Page Generator

**Files:**
- Create: `site-react/scripts/generate-static-seo-pages.mjs`
- Create: `site-react/public/assets/seo-pages.css`
- Modify: `site-react/package.json`

- [ ] **Step 1: Create the generator script**

Create `site-react/scripts/generate-static-seo-pages.mjs`:

```js
import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { englishOrigin, legalLinks, russianOrigin, seoPages, socialProfileUrls } from "./static-seo-pages-data.mjs";

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const siteRoot = path.resolve(scriptDir, "..");
const publicRoot = path.join(siteRoot, "public");

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function canonicalFor(page, language) {
  const copy = page[language];
  return `${language === "ru" ? russianOrigin : englishOrigin}${copy.path}`;
}

function buildJsonLd(page, language) {
  const copy = page[language];
  const canonicalUrl = canonicalFor(page, language);
  const origin = language === "ru" ? russianOrigin : englishOrigin;
  const organizationId = `${origin}/#organization`;
  const websiteId = `${origin}/#website`;
  const appId = `${origin}/#software-application`;

  return {
    "@context": "https://schema.org",
    "@graph": [
      {
        "@type": "Organization",
        "@id": organizationId,
        name: "NERIVA",
        url: `${origin}/`,
        logo: `${origin}/assets/brand-logo-mini.png`,
        sameAs: socialProfileUrls,
      },
      {
        "@type": "WebSite",
        "@id": websiteId,
        name: "NERIVA",
        url: `${origin}/`,
        inLanguage: language,
        publisher: { "@id": organizationId },
      },
      {
        "@type": "SoftwareApplication",
        "@id": appId,
        name: "NERIVA",
        applicationCategory: "EducationalApplication",
        operatingSystem: "Web, PWA, Telegram",
        url: canonicalUrl,
        description: copy.description,
        publisher: { "@id": organizationId },
        offers: [
          { "@type": "Offer", name: "Free", price: "0", priceCurrency: "RUB" },
          { "@type": "Offer", name: "Premium", price: "300", priceCurrency: "RUB" },
          { "@type": "Offer", name: "Platinum", price: "590", priceCurrency: "RUB" },
        ],
      },
      {
        "@type": "WebPage",
        "@id": `${canonicalUrl}#webpage`,
        url: canonicalUrl,
        name: copy.title,
        description: copy.description,
        inLanguage: language,
        isPartOf: { "@id": websiteId },
        about: { "@id": appId },
      },
      {
        "@type": "BreadcrumbList",
        "@id": `${canonicalUrl}#breadcrumb`,
        itemListElement: [
          { "@type": "ListItem", position: 1, name: "NERIVA", item: `${origin}/poliglot-ai.html${language === "en" ? "?lang=en" : ""}` },
          { "@type": "ListItem", position: 2, name: copy.h1, item: canonicalUrl },
        ],
      },
      {
        "@type": "FAQPage",
        "@id": `${canonicalUrl}#faq`,
        url: canonicalUrl,
        inLanguage: language,
        mainEntity: copy.faq.map(([question, answer]) => ({
          "@type": "Question",
          name: question,
          acceptedAnswer: { "@type": "Answer", text: answer },
        })),
      },
    ],
  };
}

function renderPage(page, language) {
  const copy = page[language];
  const otherLanguage = language === "ru" ? "en" : "ru";
  const canonicalUrl = canonicalFor(page, language);
  const ruUrl = canonicalFor(page, "ru");
  const enUrl = canonicalFor(page, "en");
  const rootPrefix = language === "en" ? "../" : "";
  const landingHref = language === "en" ? "/poliglot-ai.html?lang=en" : "/poliglot-ai.html";
  const relatedPages = seoPages.filter((item) => item.slug !== page.slug);
  const ogImage = `${language === "ru" ? russianOrigin : englishOrigin}/assets/seo/${language === "ru" ? "poliglot-ai-og-ru.jpg" : "poliglot-ai-og-en.jpg"}`;
  const appCta = language === "ru" ? "Открыть web app" : "Open web app";
  const landingLabel = language === "ru" ? "Главная" : "Landing";
  const guidesLabel = language === "ru" ? "Материалы" : "Guides";
  const switchLabel = language === "ru" ? "English" : "Русский";
  const relatedTitle = language === "ru" ? "Другие материалы" : "More guides";
  const faqTitle = language === "ru" ? "Вопросы и ответы" : "Questions and answers";

  return `<!doctype html>
<html lang="${language}" data-site-page="seo-guide">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>${escapeHtml(copy.title)}</title>
    <meta name="description" content="${escapeHtml(copy.description)}" />
    <meta name="robots" content="index, follow" />
    <link rel="canonical" href="${canonicalUrl}" />
    <link rel="alternate" hreflang="ru" href="${ruUrl}" />
    <link rel="alternate" hreflang="en" href="${enUrl}" />
    <link rel="alternate" hreflang="x-default" href="${enUrl}" />
    <meta property="og:title" content="${escapeHtml(copy.title)}" />
    <meta property="og:description" content="${escapeHtml(copy.description)}" />
    <meta property="og:type" content="article" />
    <meta property="og:url" content="${canonicalUrl}" />
    <meta property="og:image" content="${ogImage}" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="${escapeHtml(copy.title)}" />
    <meta name="twitter:description" content="${escapeHtml(copy.description)}" />
    <meta name="twitter:image" content="${ogImage}" />
    <meta name="theme-color" content="#08111f" />
    <link rel="icon" href="/assets/brand-logo-mini.png?v=20260522-react" />
    <link rel="stylesheet" href="/assets/seo-pages.css" />
    <script type="application/ld+json">${JSON.stringify(buildJsonLd(page, language))}</script>
  </head>
  <body>
    <main class="seo-page">
      <header class="seo-header">
        <a class="seo-brand" href="${landingHref}" aria-label="NERIVA">
          <img src="/assets/brand-logo-mini.png" alt="" />
          <span>NERIVA</span>
        </a>
        <nav class="seo-nav" aria-label="${escapeHtml(guidesLabel)}">
          <a href="${landingHref}">${escapeHtml(landingLabel)}</a>
          <a data-entry="web-app" class="seo-nav__cta" href="/app/">${escapeHtml(appCta)}</a>
          <a href="${page[otherLanguage].path}">${escapeHtml(switchLabel)}</a>
        </nav>
      </header>
      <article class="seo-article">
        <section class="seo-hero">
          <p class="seo-eyebrow">NERIVA web app</p>
          <h1>${escapeHtml(copy.h1)}</h1>
          <p>${escapeHtml(copy.intro)}</p>
          <a data-entry="web-app" class="seo-button" href="/app/">${escapeHtml(appCta)}</a>
        </section>
        <section class="seo-sections" aria-label="${escapeHtml(copy.h1)}">
          ${copy.sections.map(([title, body]) => `<section class="seo-section"><h2>${escapeHtml(title)}</h2><p>${escapeHtml(body)}</p></section>`).join("\n          ")}
        </section>
        <section class="seo-faq" id="faq">
          <h2>${escapeHtml(faqTitle)}</h2>
          ${copy.faq.map(([question, answer]) => `<article class="seo-faq__item"><h3>${escapeHtml(question)}</h3><p>${escapeHtml(answer)}</p></article>`).join("\n          ")}
        </section>
        <section class="seo-related" aria-label="${escapeHtml(relatedTitle)}">
          <h2>${escapeHtml(relatedTitle)}</h2>
          <div class="seo-related__grid">
            ${relatedPages.map((item) => `<a href="${item[language].path}">${escapeHtml(item[language].h1)}</a>`).join("\n            ")}
          </div>
        </section>
      </article>
      <footer class="seo-footer">
        <nav aria-label="${escapeHtml(language === "ru" ? "Документы" : "Documents")}">
          ${legalLinks[language].map((link) => `<a href="${link.href}">${escapeHtml(link.label)}</a>`).join("\n          ")}
        </nav>
        <nav aria-label="Social">
          ${socialProfileUrls.map((url) => `<a href="${url}" target="_blank" rel="noreferrer">${new URL(url).hostname.replace("www.", "")}</a>`).join("\n          ")}
        </nav>
      </footer>
    </main>
  </body>
</html>
`;
}

function renderSitemap() {
  const landingAlternates = [
    '<xhtml:link rel="alternate" hreflang="ru" href="https://poliglotai.ru/poliglot-ai.html" />',
    '<xhtml:link rel="alternate" hreflang="en" href="https://poliglotai.online/poliglot-ai.html?lang=en" />',
    '<xhtml:link rel="alternate" hreflang="x-default" href="https://poliglotai.online/poliglot-ai.html?lang=en" />',
  ];
  const urls = [
    `<url>
    <loc>https://poliglotai.ru/poliglot-ai.html</loc>
    ${landingAlternates.join("\n    ")}
  </url>`,
    `<url>
    <loc>https://poliglotai.online/poliglot-ai.html?lang=en</loc>
    ${landingAlternates.join("\n    ")}
  </url>`,
  ];

  for (const page of seoPages) {
    const alternates = [
      `<xhtml:link rel="alternate" hreflang="ru" href="${canonicalFor(page, "ru")}" />`,
      `<xhtml:link rel="alternate" hreflang="en" href="${canonicalFor(page, "en")}" />`,
      `<xhtml:link rel="alternate" hreflang="x-default" href="${canonicalFor(page, "en")}" />`,
    ].join("\n    ");
    urls.push(`<url>
    <loc>${canonicalFor(page, "ru")}</loc>
    ${alternates}
  </url>`);
    urls.push(`<url>
    <loc>${canonicalFor(page, "en")}</loc>
    ${alternates}
  </url>`);
  }

  for (const path of ["/privacy.html", "/terms.html", "/agreement.html", "/consent.html"]) {
    urls.push(`<url><loc>https://poliglotai.ru${path}</loc></url>`);
    urls.push(`<url><loc>https://poliglotai.online${path}?lang=en</loc></url>`);
  }

  return `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">
  ${urls.join("\n  ")}
</urlset>
`;
}

await fs.mkdir(path.join(publicRoot, "en"), { recursive: true });

for (const page of seoPages) {
  await fs.writeFile(path.join(publicRoot, `${page.slug}.html`), renderPage(page, "ru"), "utf8");
  await fs.writeFile(path.join(publicRoot, "en", `${page.slug}.html`), renderPage(page, "en"), "utf8");
}

await fs.writeFile(path.join(publicRoot, "sitemap.xml"), renderSitemap(), "utf8");
console.log(`Generated ${seoPages.length * 2} static SEO pages and sitemap.xml`);
```

- [ ] **Step 2: Add shared CSS**

Create `site-react/public/assets/seo-pages.css`:

```css
:root {
  color-scheme: dark;
  --seo-bg: #08111f;
  --seo-panel: #101b2c;
  --seo-panel-soft: #142238;
  --seo-text: #eef4ff;
  --seo-muted: #aebbd0;
  --seo-line: rgba(174, 187, 208, 0.24);
  --seo-accent: #5eead4;
  --seo-warm: #f9d56e;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  min-width: 320px;
  background: var(--seo-bg);
  color: var(--seo-text);
  font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  line-height: 1.6;
}

a {
  color: inherit;
}

.seo-page {
  min-height: 100vh;
}

.seo-header,
.seo-article,
.seo-footer {
  width: min(1120px, calc(100% - 32px));
  margin: 0 auto;
}

.seo-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 24px 0;
}

.seo-brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  font-weight: 800;
}

.seo-brand img {
  width: 34px;
  height: 34px;
}

.seo-nav,
.seo-footer nav {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 14px;
}

.seo-nav a,
.seo-footer a,
.seo-related a {
  color: var(--seo-muted);
  text-decoration: none;
}

.seo-nav a:hover,
.seo-footer a:hover,
.seo-related a:hover {
  color: var(--seo-text);
}

.seo-nav__cta,
.seo-button {
  border: 1px solid rgba(94, 234, 212, 0.55);
  border-radius: 8px;
  background: rgba(94, 234, 212, 0.12);
  color: var(--seo-text) !important;
  font-weight: 700;
}

.seo-nav__cta {
  padding: 8px 12px;
}

.seo-hero {
  padding: 72px 0 44px;
}

.seo-eyebrow {
  margin: 0 0 14px;
  color: var(--seo-accent);
  font-size: 0.82rem;
  font-weight: 800;
  letter-spacing: 0;
  text-transform: uppercase;
}

h1,
h2,
h3,
p {
  margin-top: 0;
}

h1 {
  max-width: 820px;
  margin-bottom: 18px;
  font-size: clamp(2.4rem, 7vw, 5.8rem);
  line-height: 0.98;
  letter-spacing: 0;
}

.seo-hero p {
  max-width: 760px;
  color: var(--seo-muted);
  font-size: 1.12rem;
}

.seo-button {
  display: inline-flex;
  margin-top: 14px;
  padding: 12px 16px;
  text-decoration: none;
}

.seo-sections {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin: 28px 0 52px;
}

.seo-section,
.seo-faq__item,
.seo-related {
  border: 1px solid var(--seo-line);
  border-radius: 8px;
  background: var(--seo-panel);
}

.seo-section,
.seo-faq__item {
  padding: 22px;
}

.seo-section h2,
.seo-faq h2,
.seo-related h2 {
  font-size: 1.35rem;
  letter-spacing: 0;
}

.seo-section p,
.seo-faq__item p {
  color: var(--seo-muted);
}

.seo-faq {
  display: grid;
  gap: 14px;
  margin: 0 0 52px;
}

.seo-faq h2,
.seo-related h2 {
  margin-bottom: 4px;
}

.seo-related {
  padding: 22px;
  margin-bottom: 56px;
}

.seo-related__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.seo-related a {
  display: block;
  min-height: 54px;
  padding: 14px;
  border-radius: 8px;
  background: var(--seo-panel-soft);
}

.seo-footer {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  padding: 28px 0 40px;
  border-top: 1px solid var(--seo-line);
}

@media (max-width: 760px) {
  .seo-header,
  .seo-footer {
    align-items: flex-start;
    flex-direction: column;
  }

  .seo-hero {
    padding-top: 42px;
  }

  .seo-sections,
  .seo-related__grid {
    grid-template-columns: 1fr;
  }
}
```

- [ ] **Step 3: Hook generation into package scripts**

Modify `site-react/package.json` scripts to:

```json
{
  "scripts": {
    "dev": "vite --host 127.0.0.1 --port 5175",
    "seo:generate": "node scripts/generate-static-seo-pages.mjs",
    "build": "npm run seo:generate && node ../tools/clean_public_site_bundle.mjs && tsc -b && vite build",
    "preview": "vite preview --host 127.0.0.1 --port 4175",
    "e2e": "playwright test"
  }
}
```

- [ ] **Step 4: Run generator and verify GREEN for page tests**

Run:

```powershell
npm --prefix site-react run seo:generate
npm --prefix site-react run e2e -- --grep "static SEO pages"
```

Expected: static SEO page tests now pass or fail only because the landing `.seo-guides` links are still missing.

## Task 4: Add Compact Landing Guide Links

**Files:**
- Modify: `site-react/src/landingSeoContent.ts`
- Modify: `site-react/src/EnglishSparkLanding.tsx`
- Modify: `site-react/src/englishSparkLanding.css` if existing footer styles need a small `.seo-guides` block

- [ ] **Step 1: Export guide link data**

Add to `site-react/src/landingSeoContent.ts`:

```ts
export const staticSeoGuideLinks = [
  {
    slug: "ai-english-tutor",
    ru: { href: "/ai-english-tutor.html", label: "AI-репетитор английского" },
    en: { href: "/en/ai-english-tutor.html", label: "AI English tutor" },
  },
  {
    slug: "english-speaking-practice",
    ru: { href: "/english-speaking-practice.html", label: "Разговорная практика" },
    en: { href: "/en/english-speaking-practice.html", label: "Speaking practice" },
  },
  {
    slug: "english-pronunciation-trainer",
    ru: { href: "/english-pronunciation-trainer.html", label: "Произношение" },
    en: { href: "/en/english-pronunciation-trainer.html", label: "Pronunciation trainer" },
  },
  {
    slug: "english-for-work-and-travel",
    ru: { href: "/english-for-work-and-travel.html", label: "Работа и поездки" },
    en: { href: "/en/english-for-work-and-travel.html", label: "Work and travel" },
  },
  {
    slug: "language-learning-web-app",
    ru: { href: "/language-learning-web-app.html", label: "Web app для языков" },
    en: { href: "/en/language-learning-web-app.html", label: "Language web app" },
  },
] as const;
```

- [ ] **Step 2: Render the footer-only guide area**

In `site-react/src/EnglishSparkLanding.tsx`, import `staticSeoGuideLinks`, derive current locale the same way existing landing SEO content does, and render this block inside the existing footer, not inside `.landing-hero` or `.public-nav`:

```tsx
<nav className="seo-guides" aria-label={language === "ru" ? "Материалы" : "Guides"}>
  <strong>{language === "ru" ? "Материалы" : "Guides"}</strong>
  {staticSeoGuideLinks.map((item) => {
    const link = item[language === "ru" ? "ru" : "en"];
    return (
      <a key={item.slug} href={link.href}>
        {link.label}
      </a>
    );
  })}
</nav>
```

- [ ] **Step 3: Add restrained footer styles if needed**

If no existing footer style supports this, add to `site-react/src/englishSparkLanding.css`:

```css
.seo-guides {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 14px;
  align-items: center;
}

.seo-guides strong {
  color: var(--landing-text, #f8fbff);
  font-size: 0.92rem;
}

.seo-guides a {
  color: var(--landing-muted, rgba(248, 251, 255, 0.72));
  font-size: 0.92rem;
  text-decoration: none;
}

.seo-guides a:hover {
  color: var(--landing-text, #f8fbff);
}
```

- [ ] **Step 4: Re-run the focused test**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "static SEO pages"
```

Expected: PASS for the static SEO page suite.

## Task 5: Build, Full E2E, Browser Spot Check, Commit, Push

**Files:**
- Generated output under the public deployment directory may change when `npm --prefix site-react run build` runs.

- [ ] **Step 1: Run the build**

Run:

```powershell
npm --prefix site-react run build
```

Expected: exit code 0. Output includes `Generated 10 static SEO pages and sitemap.xml`, TypeScript build succeeds, and Vite build succeeds.

- [ ] **Step 2: Run the public site e2e suite**

Run:

```powershell
npm --prefix site-react run e2e
```

Expected: exit code 0 for all public-site Playwright tests.

- [ ] **Step 3: Use Browser/Playwright MCP for a visual spot check**

Open `http://127.0.0.1:4175/en/ai-english-tutor.html` in the in-app browser after the preview server is available. Confirm:

- H1 is visible.
- The page is not visually blank.
- CTA points to `/app/`.
- Footer links and guide links do not overlap on mobile or desktop.

- [ ] **Step 4: Review git diff**

Run:

```powershell
git status --short
git diff -- site-react docs/superpowers/plans/2026-06-29-static-seo-pages.md
```

Expected: only planned files and generated public-site output are changed, alongside pre-existing unrelated dirty files that must not be reverted.

- [ ] **Step 5: Commit only this task's changes**

Run:

```powershell
git add docs/superpowers/plans/2026-06-29-static-seo-pages.md site-react
git add "Сайт полиглота для бота"
git commit -m "feat: add static seo page cluster"
```

Expected: commit succeeds. Do not stage unrelated root binaries or unrelated generated files outside this task.

- [ ] **Step 6: Push the branch**

Run:

```powershell
git push origin codex/ai-tutor-rebuild-fix
```

Expected: branch pushes to `https://github.com/GarryNaxyison/My-app-local.git`.

## Self-Review

- Spec coverage: The plan covers static bilingual pages, web app first positioning, Telegram as companion, sitemap, canonical alternates, visible guide links, JSON-LD, tests, build, browser spot check, commit, and push.
- Completeness scan: No unresolved marker text remains.
- Type consistency: `seoPages`, `staticSeoGuideLinks`, canonical path fields, and Playwright page fixtures use the same slug/path/title names across tasks.
