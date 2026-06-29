export const landingLanguageCodes = [
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
] as const;

export type LandingLanguageCode = (typeof landingLanguageCodes)[number];
export type LandingSeoLocale = "ru" | "en";

export const russianSiteOrigin = "https://poliglotai.ru";
export const internationalSiteOrigin = "https://poliglotai.online";

export const socialProfileUrls = [
  "https://www.youtube.com/@PoliglotAI",
  "https://www.instagram.com/poliglotai.online/",
  "https://www.tiktok.com/@poliglotai.online",
] as const;

export const landingSeoCopy = {
  ru: {
    title: "Poliglot AI - AI-репетитор английского и языков в Telegram",
    description: "Poliglot AI: AI-уроки, разговорная практика, Telegram, web app, произношение, фото-перевод, ошибки и Premium в одном профиле.",
    sectionEyebrow: "Поисковые ответы",
    sectionTitle: "Ответы для поиска и AI",
    sectionIntro: "Короткие прямые ответы о Poliglot AI для людей, которые сравнивают AI-репетиторов, языковые Telegram-боты, speaking practice, голос и фото-перевод.",
    questions: [
      {
        question: "Что такое Poliglot AI?",
        answer: "Poliglot AI — это AI-репетитор языков в web app и Telegram-боте. Он объединяет короткие уроки, разговорную практику, произношение, перевод текста с фото, ошибки, заметки и прогресс в одном профиле.",
      },
      {
        question: "Можно ли учить английский с ИИ в Telegram?",
        answer: "Да. В Telegram можно запускать практику, получать задания, отправлять ответы и голосовые сообщения, а прогресс сохраняется вместе с web app.",
      },
      {
        question: "Чем AI-репетитор отличается от обычного приложения со словами?",
        answer: "Poliglot AI не ограничивается списком слов: он дает фразу в контексте, просит ответить, исправляет ошибку и возвращает слабое место в повторение.",
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
    title: "Poliglot AI - AI language tutor in Telegram and web app",
    description: "Poliglot AI: AI language tutor, speaking practice, Telegram bot, web app, voice, photo translation, mistakes, and premium plans in one profile.",
    sectionEyebrow: "Search answers",
    sectionTitle: "Answers for search and AI assistants",
    sectionIntro: "Short direct answers about Poliglot AI for people comparing AI tutors, Telegram language bots, speaking practice, voice, and photo translation.",
    questions: [
      {
        question: "What is Poliglot AI?",
        answer: "Poliglot AI is an AI language tutor in a web app and Telegram bot. It combines short lessons, speaking practice, pronunciation, photo translation, mistakes, notes, and progress in one profile.",
      },
      {
        question: "Can I practice English with an AI tutor in Telegram?",
        answer: "Yes. In Telegram you can start practice, receive tasks, send answers and voice messages, while progress stays synced with the web app.",
      },
      {
        question: "How is an AI tutor different from a vocabulary app?",
        answer: "Poliglot AI is not only a word list. It gives a phrase in context, asks for your answer, corrects the mistake, and brings the weak spot back for review.",
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
