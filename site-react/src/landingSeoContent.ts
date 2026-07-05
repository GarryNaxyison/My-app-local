export const landingLanguageCodes = ["ru", "en"] as const;

export type LandingLanguageCode = (typeof landingLanguageCodes)[number];
export type LandingSeoLocale = "ru" | "en";

export const russianSiteOrigin = "https://neriva.ru";
export const internationalSiteOrigin = "https://neriva.ru";

export const socialProfileUrls = [
  "https://www.youtube.com/@neriva_app",
  "https://www.instagram.com/neriva.ru",
  "https://tiktok.com/@nerivaru",
  "https://t.me/NERIVAapp_bot",
] as const;

export const staticSeoGuideLinks = [
  {
    slug: "ai-tutor",
    ru: { href: "/ai-tutor.html", label: "AI-репетитор" },
    en: { href: "/en/ai-tutor.html", label: "AI tutor" },
  },
  {
    slug: "speaking-practice",
    ru: { href: "/speaking-practice.html", label: "Разговорная практика" },
    en: { href: "/en/speaking-practice.html", label: "Speaking practice" },
  },
  {
    slug: "pronunciation",
    ru: { href: "/pronunciation.html", label: "Произношение" },
    en: { href: "/en/pronunciation.html", label: "Pronunciation" },
  },
  {
    slug: "photo-translation",
    ru: { href: "/photo-translation.html", label: "Фото-перевод" },
    en: { href: "/en/photo-translation.html", label: "Photo translation" },
  },
  {
    slug: "telegram-language-bot",
    ru: { href: "/telegram-language-bot.html", label: "Telegram-бот" },
    en: { href: "/en/telegram-language-bot.html", label: "Telegram bot" },
  },
  {
    slug: "ai-english-tutor",
    ru: { href: "/ai-english-tutor.html", label: "AI-репетитор английского" },
    en: { href: "/en/ai-english-tutor.html", label: "AI English tutor" },
  },
  {
    slug: "english-speaking-practice",
    ru: { href: "/english-speaking-practice.html", label: "Разговорный английский" },
    en: { href: "/en/english-speaking-practice.html", label: "English speaking practice" },
  },
  {
    slug: "english-pronunciation-trainer",
    ru: { href: "/english-pronunciation-trainer.html", label: "Тренажер произношения" },
    en: { href: "/en/english-pronunciation-trainer.html", label: "Pronunciation trainer" },
  },
  {
    slug: "language-learning-web-app",
    ru: { href: "/language-learning-web-app.html", label: "Web app для языков" },
    en: { href: "/en/language-learning-web-app.html", label: "Language web app" },
  },
] as const;

export const landingSeoCopy = {
  ru: {
    title: "NERIVA - AI-репетитор иностранных языков в Telegram",
    description: "NERIVA: AI-уроки, разговорная практика, Telegram, web app, произношение, фото-перевод, ошибки и Premium в одном профиле.",
    sectionEyebrow: "Поисковые ответы",
    sectionTitle: "Ответы для поиска и AI",
    sectionIntro: "Короткие прямые ответы о NERIVA для людей, которые сравнивают AI-репетиторов, языковые Telegram-боты, speaking practice, голос и фото-перевод.",
    comparisonEyebrow: "Сравнение вариантов",
    comparisonTitle: "NERIVA в сравнении с инструментами, которые обычно ищут",
    comparisonIntro: "Коротко о том, когда словарное приложение, AI-репетитор, Telegram-бот или web app лучше подходят для языковой практики.",
    comparisons: [
      {
        title: "NERIVA vs словарное приложение",
        alternativeLabel: "Словарное приложение",
        alternative: "Списки слов, карточки, интервальные повторения и отдельные значения.",
        productLabel: "NERIVA",
        product: "Фразы в контексте, ответы, исправления, повторение слабых мест, голос и фото-практика.",
        verdict: "Подходит, когда нужны исправления и контекст, а не только запоминание.",
      },
      {
        title: "AI-репетитор vs языковой бот",
        alternativeLabel: "Языковой бот",
        alternative: "Быстрые подсказки в Telegram и простые ответы.",
        productLabel: "AI-репетитор",
        product: "Уроки с направлением, roleplay, объяснение ошибок, прогресс и следующее повторение.",
        verdict: "Подходит, когда удобство Telegram нужно соединить с обратной связью уровня репетитора.",
      },
      {
        title: "Telegram-бот vs web app",
        alternativeLabel: "Telegram-бот",
        alternative: "Быстрые задания, голосовые сообщения, напоминания и короткая ежедневная практика.",
        productLabel: "Web app",
        product: "Длинные сессии, тарифы, заметки, прогресс, ошибки и управление профилем.",
        verdict: "Подходит, когда быстрая мобильная практика и глубокая учеба за компьютером должны синхронизироваться.",
      },
    ],
    questions: [
      {
        question: "Что такое NERIVA?",
        answer: "NERIVA — это AI-репетитор иностранных языков в web app и Telegram-боте. Он объединяет короткие уроки, разговорную практику, произношение, перевод текста с фото, ошибки, заметки и прогресс в одном профиле.",
      },
      {
        question: "Можно ли учить английский с ИИ в Telegram?",
        answer: "Да. В Telegram можно запускать практику, получать задания, отправлять ответы и голосовые сообщения, а прогресс сохраняется вместе с web app.",
      },
      {
        question: "Чем AI-репетитор отличается от обычного приложения со словами?",
        answer: "NERIVA не ограничивается списком слов: он дает фразу в контексте, просит ответить, исправляет ошибку и возвращает слабое место в повторение.",
      },
      {
        question: "Можно ли тренировать произношение и speaking?",
        answer: "Да. Voice Coach и shadowing помогают тренировать речь, видеть слабые слова, получать score и повторять более естественную фразу.",
      },
      {
        question: "Можно ли переводить текст с фото?",
        answer: "Да. Фото меню, вывески или задания превращается в перевод, заметку и короткую практику по этому контексту.",
      },
      {
        question: "Есть ли бесплатный тариф?",
        answer: "Да. Free дает стартовые уроки и практику без оплаты, а Premium и Platinum открывают больше лимитов, голос, фото и интенсивную ежедневную учебу.",
      },
    ],
  },
  en: {
    title: "NERIVA - AI language tutor in Telegram and web app",
    description: "NERIVA: AI language tutor, speaking practice, Telegram bot, web app, voice, photo translation, mistakes, and premium plans in one profile.",
    sectionEyebrow: "Search answers",
    sectionTitle: "Answers for search and AI assistants",
    sectionIntro: "Short direct answers about NERIVA for people comparing AI tutors, Telegram language bots, speaking practice, voice, and photo translation.",
    comparisonEyebrow: "Compare options",
    comparisonTitle: "NERIVA compared with the tools people usually search for",
    comparisonIntro: "See when a vocabulary app, AI tutor, Telegram bot, or web app is the better fit for language practice.",
    comparisons: [
      {
        title: "NERIVA vs vocabulary app",
        alternativeLabel: "Vocabulary app",
        alternative: "Word lists, flashcards, spaced repetition, and isolated meanings.",
        productLabel: "NERIVA",
        product: "Context phrases, answers, corrections, weak-spot review, voice, and photo practice.",
        verdict: "Best when you need correction and context, not only memorization.",
      },
      {
        title: "AI tutor vs language bot",
        alternativeLabel: "Language bot",
        alternative: "Quick chat prompts in Telegram with simple answers.",
        productLabel: "AI tutor",
        product: "Guided lessons, roleplay, mistake explanations, progress, and next repetition.",
        verdict: "Best when Telegram convenience needs tutor-level feedback.",
      },
      {
        title: "Telegram bot vs web app",
        alternativeLabel: "Telegram bot",
        alternative: "Fast tasks, voice messages, reminders, and short daily practice.",
        productLabel: "Web app",
        product: "Longer sessions, pricing, notes, progress, mistakes, and profile control.",
        verdict: "Best when quick mobile practice and deeper desktop study should stay synced.",
      },
    ],
    questions: [
      {
        question: "What is NERIVA?",
        answer: "NERIVA is an AI language tutor in a web app and Telegram bot. It combines short lessons, speaking practice, pronunciation, photo translation, mistakes, notes, and progress in one profile.",
      },
      {
        question: "Can I practice English with an AI tutor in Telegram?",
        answer: "Yes. In Telegram you can start practice, receive tasks, send answers and voice messages, while progress stays synced with the web app.",
      },
      {
        question: "How is an AI tutor different from a vocabulary app?",
        answer: "NERIVA is not only a word list. It gives a phrase in context, asks for your answer, corrects the mistake, and brings the weak spot back for review.",
      },
      {
        question: "Can I practice pronunciation and speaking?",
        answer: "Yes. Voice Coach and shadowing help you practice speech, see weak words, get a score, and repeat a more natural phrase.",
      },
      {
        question: "Can I translate text from photos?",
        answer: "Yes. A photo of a menu, sign, or exercise becomes a translation, note, and short practice prompt in that context.",
      },
      {
        question: "Is there a free plan?",
        answer: "Yes. Free gives starter lessons and practice without payment, while Premium and Platinum unlock higher limits, voice, photo tools, and intensive daily study.",
      },
    ],
  },
} as const;

export function seoLocaleForLanguage(language: string): LandingSeoLocale {
  return language === "ru" ? "ru" : "en";
}

export function normalizeLandingLanguage(value: string | null | undefined): LandingLanguageCode {
  const normalized = String(value || "").toLowerCase().trim();
  return landingLanguageCodes.includes(normalized as LandingLanguageCode) ? (normalized as LandingLanguageCode) : "ru";
}

export function canonicalUrlForLanguage(language: LandingLanguageCode): string {
  const origin = language === "ru" ? russianSiteOrigin : internationalSiteOrigin;
  const url = new URL("/poliglot-ai.html", origin);
  if (language !== "ru") url.searchParams.set("lang", language);
  return url.toString();
}

export function canonicalOriginForLanguage(language: LandingLanguageCode): string {
  return language === "ru" ? russianSiteOrigin : internationalSiteOrigin;
}
