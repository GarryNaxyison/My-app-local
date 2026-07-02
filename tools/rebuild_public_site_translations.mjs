import fs from "node:fs";
import path from "node:path";
import vm from "node:vm";

const ROOT = process.cwd();
const SITE_DIR = path.join(ROOT, "Сайт полиглота для бота");
const PHRASES_PATH = path.join(SITE_DIR, "assets", "site-phrases.js");
const I18N_PATH = path.join(SITE_DIR, "assets", "site-i18n.js");
const MIRROR_PHRASES_PATHS = [
  path.join(ROOT, "site-react", "public", "assets", "site-phrases.js"),
];
const PAGES = ["poliglot-ai.html", "terms.html", "privacy.html"];
const REACT_SOURCES = [
  path.join(ROOT, "site-react", "src", "PublicSiteApp.tsx"),
  path.join(ROOT, "site-react", "src", "EnglishSparkLanding.tsx"),
  path.join(ROOT, "site-react", "src", "legacyLegalContent.ts"),
];
const GENERATED_START = "// <public-site-translations-generated>";
const GENERATED_END = "// </public-site-translations-generated>";

const LANGS = [
  ["ru", "ru"],
  ["en", "en"],
  ["es", "es"],
  ["de", "de"],
  ["fr", "fr"],
  ["it", "it"],
  ["zh", "zh-CN"],
  ["ja", "ja"],
  ["ko", "ko"],
  ["tg", "tg"],
  ["uz", "uz"],
  ["tt", "tt"],
  ["hy", "hy"],
  ["kk", "kk"],
  ["ky", "ky"],
  ["ka", "ka"],
  ["uk", "uk"],
  ["pl", "pl"],
  ["ro", "ro"],
  ["pt", "pt"],
  ["ar", "ar"],
  ["bn", "bn"],
  ["cs", "cs"],
  ["el", "el"],
  ["hi", "hi"],
  ["hu", "hu"],
  ["id", "id"],
  ["nl", "nl"],
  ["sv", "sv"],
  ["ta", "ta"],
  ["te", "te"],
  ["th", "th"],
  ["tl", "tl"],
  ["tr", "tr"],
  ["vi", "vi"],
];

const SITE_LANG_CODES = LANGS.map(([code]) => code);
const TRANSLATION_CONCURRENCY = Math.max(1, Number(process.env.TRANSLATION_CONCURRENCY || 4));
const TRANSLATION_PAUSE_MS = Math.max(0, Number(process.env.TRANSLATION_PAUSE_MS || 0));
const TRANSLATION_TIMEOUT_MS = Math.max(1_000, Number(process.env.TRANSLATION_TIMEOUT_MS || 12_000));
const NATIVE_LANGUAGE_NAMES = new Set([
  "English",
  "Español",
  "Deutsch",
  "Français",
  "Italiano",
  "Тоҷикӣ",
  "O'zbekcha",
  "Татарча",
  "Հայերեն",
  "Қазақша",
  "Українська",
  "Polski",
  "Română",
  "Português",
  "Кыргызча",
  "ქართული",
  "العربية",
  "বাংলা",
  "Čeština",
  "Ελληνικά",
  "हिंदी",
  "Magyar",
  "Bahasa Indonesia",
  "Nederlands",
  "svenska",
  "தமிழ்",
  "తెలుగు",
  "ภาษาไทย",
  "Tagalog",
  "Türkçe",
  "Tiếng Việt",
]);

const PROTECTED_PATTERNS = [
  /<[^>]+>/g,
  /@[A-Za-z0-9_]+/g,
  /[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}/gi,
  /\bpoliglotai\.(?:ru|online)\/app\b/gi,
  /\b(?:NERIVA|NERIVA)\b/g,
  /\b(?:Telegram|Premium|Platinum|Free|Stars|TON|USDT|RUB|YooKassa|SBP|CEFR|OCR|AI)\b/g,
  /\b\d+(?:[.,]\d+)?\s*(?:₽|RUB|Stars|USDT|TON)\b/gi,
  /≈\s*\d[\d\s,.]*(?:₽|RUB|Stars|USDT|TON)?(?:\/год|\/year)?/gi,
];

const CURATED_RU_TRANSLATIONS = {
  "Premium AI tutor for daily practice": "Премиум AI-тьютор для ежедневной практики",
  "Practice speaking before the moment matters": "Потренируйте речь до важного момента",
  "Lessons, dialogue, pronunciation, photo translation, mistake review, and progress live in one profile across web, PWA, and Telegram.": "Уроки, диалоги, произношение, перевод по фото, разбор ошибок и прогресс хранятся в одном профиле: веб, PWA и Telegram.",
  "Start free": "Начать бесплатно",
  "Open Telegram": "Открыть Telegram",
  "Travel": "Поездки",
  "hotel, cafe, airport": "отель, кафе, аэропорт",
  "Work": "Работа",
  "calls, emails, small talk": "звонки, письма, small talk",
  "Exam": "Экзамен",
  "vocabulary, grammar, speaking": "лексика, грамматика, речь",
  "Conversation": "Разговор",
  "live answers without memorizing": "живые ответы без заучивания",
  "Popular learning languages:": "Популярные языки:",
  "speaking routes": "маршруты для речи",
  "travel dialogues": "диалоги в поездке",
  "work and exams": "работа и экзамены",
  "listening and speech": "аудирование и речь",
  "real-life scenes": "живые ситуации",
  "words and practice": "слова и практика",
  "lessons and review": "уроки и повторение",
  "vocabulary and speech": "лексика и речь",
  "Goal": "Цель",
  "Speech": "Речь",
  "Review": "Повторение",
  "35 languages": "35 языков",
  "interface coverage": "интерфейс",
  "learning levels": "уровни обучения",
  "Free": "Free",
  "start without payment": "старт без оплаты",
  "Today": "Сегодня",
  "AI Tutor Cockpit": "Панель AI Tutor",
  "NERIVA gives the phrase, meaning, example, and one short next step.": "NERIVA дает фразу, смысл, пример и короткий следующий шаг.",
  "I would like to book a table for tonight.": "Я хотел бы забронировать столик на сегодня вечером.",
  "The phrase moves straight into your route and review loop.": "Фраза сразу попадает в маршрут и повторение.",
  "AI Tutor routes": "Маршруты AI Tutor",
  "Routes for the goal, not an endless exercise feed": "Маршруты под цель, а не бесконечная лента упражнений",
  "Choose travel, work, exam, or conversation practice. NERIVA connects lesson, dialogue, voice, photo, and review into one controlled cycle.": "Выберите поездки, работу, экзамен или разговорную практику. NERIVA связывает урок, диалог, голос, фото и повторение в один понятный цикл.",
  "AI Tutor Core": "База AI Tutor",
  "A complete A1-C2 route with lessons, practice, review, and XP in one controlled loop.": "Полный маршрут A1-C2 с уроками, практикой, повторением и XP в одном цикле.",
  "Start route": "Начать маршрут",
  "Travel & Work": "Поездки и работа",
  "Roleplay scenes, photo translation, quick phrases, and dialogues for moments outside the classroom.": "Ролевые сцены, перевод по фото, быстрые фразы и диалоги для ситуаций вне урока.",
  "Voice Coach": "Тренер произношения",
  "Shadowing, weak words, pronunciation score, and voice history for steady speaking practice.": "Шэдоуинг, слабые слова, оценка произношения и история голоса для регулярной практики.",
  "profile for web, PWA, Telegram": "профиль для веба, PWA и Telegram",
  "AI Tutor cockpit": "Панель AI Tutor",
  "Every session ends with a visible next step": "Каждая сессия заканчивается понятным следующим шагом",
  "NERIVA does not leave you alone with a task. It explains, asks for an answer, checks voice or photo context, and returns to mistakes until the weak spot becomes familiar.": "NERIVA не оставляет вас один на один с заданием: объясняет, просит ответить, проверяет голос или фото и возвращает к ошибкам, пока слабое место не станет понятным.",
  "interface languages": "языков интерфейса",
  "level routes": "маршруты по уровням",
  "profile everywhere": "один профиль везде",
  "AI Tutor": "AI Tutor",
  "Guided lessons turn a phrase into an answer, a correction, and a next action.": "Уроки превращают фразу в ответ, исправление и следующее действие.",
  "Explains a phrase, gives an example, asks for your answer, and shows the next correction.": "Объясняет фразу, дает пример, просит ответить и показывает следующее исправление.",
  "lesson -> check": "урок -> проверка",
  "Roleplay": "Ролевая сцена",
  "Travel, work, exam, or casual speaking scenarios stay on topic and correct your reply.": "Сценарии поездки, работы, экзамена или обычного разговора держат контекст и исправляют ответ.",
  "Keeps travel, work, exam, or casual speaking scenarios focused without empty replies.": "Держит фокус на поездке, работе, экзамене или обычном разговоре без пустых реплик.",
  "context role": "роль и контекст",
  "Voice checks, shadowing, and repeat prompts help you hear and fix pronunciation quickly.": "Проверки голоса, шэдоуинг и повторение помогают быстро услышать и исправить произношение.",
  "score + weak words": "оценка + слабые слова",
  "Photo Tools": "Фотоинструменты",
  "Menus, signs, or tasks become translation, notes, and a practice prompt in context.": "Меню, вывески и задания превращаются в перевод, заметки и практику по контексту.",
  "text from image": "текст с фото",
  "Mistake Loop": "Цикл ошибок",
  "Mistakes, notes, weak words, XP, and streaks keep weak spots visible until they become easy.": "Ошибки, заметки, слабые слова, XP и серия держат слабые места на виду, пока они не станут проще.",
  "mistake -> repeat": "ошибка -> повтор",
  "Live learning loop": "Живой цикл обучения",
  "AI Tutor is checking the answer": "AI Tutor проверяет ответ",
  "Say it naturally: \"Could you help me check in?\"": "Скажите естественно: «Could you help me check in?»",
  "Better: \"I have a reservation under my name.\"": "Лучше: «I have a reservation under my name.»",
  "weak word fixed": "слабое слово исправлено",
  "translation saved": "перевод сохранен",
  "Scenarios": "Сценарии",
  "Real situations where the language has to work today": "Реальные ситуации, где язык нужен уже сегодня",
  "At the airport, on a work call, before an exam, or in everyday conversation, NERIVA turns practice into a concrete speaking moment.": "В аэропорту, на рабочем звонке, перед экзаменом или в обычном разговоре NERIVA превращает практику в конкретный речевой момент.",
  "Travel without panic": "Поездка без паники",
  "Short hotel, cafe, airport, and doctor phrases are practiced before the moment gets stressful.": "Короткие фразы для отеля, кафе, аэропорта и врача отрабатываются заранее, до стрессового момента.",
  "Photo": "Фото",
  "Phrasebook": "Разговорник",
  "Could you help me check in?": "Could you help me check in?",
  "Work calls and messages": "Рабочие звонки и переписка",
  "Emails, calls, self-intros, and deadline questions become rehearsals that make speaking easier.": "Письма, звонки, самопрезентация и вопросы о сроках становятся репетицией, после которой говорить проще.",
  "Notes": "Заметки",
  "Let me clarify the deadline.": "Let me clarify the deadline.",
  "Exam and level progress": "Экзамен и рост уровня",
  "A1-C2 vocabulary, grammar, listening, and speaking stay in a route with visible progress steps.": "Лексика, грамматика, аудирование и речь от A1 до C2 идут по маршруту с понятными шагами прогресса.",
  "Listening": "Аудирование",
  "Mistakes": "Ошибки",
  "I agree with the statement because...": "I agree with the statement because...",
  "Speaking": "Разговор",
  "Everyday conversation": "Разговор на каждый день",
  "The tutor gives a natural version and a next repetition until your answer sounds more confident.": "Тьютор показывает естественный вариант и дает повторение, пока ответ не звучит увереннее.",
  "Voice": "Голос",
  "Shadowing": "Шэдоуинг",
  "I have been trying to say...": "I have been trying to say...",
  "Web + mobile + Telegram": "Веб, мобильная версия и Telegram",
  "One profile for web, mobile, and Telegram": "Один профиль для веба, мобильной версии и Telegram",
  "Start a longer lesson on desktop, repeat weak words from mobile, and return to Telegram without losing Premium status, progress, notes, or mistake history.": "Начните длинный урок на компьютере, повторите слабые слова с телефона и вернитесь в Telegram без потери Premium, прогресса, заметок и истории ошибок.",
  "Web app": "Веб-приложение",
  "Use the large screen for longer sessions, pricing, progress, mistake review, and AI Tutor work.": "Большой экран удобен для длинных занятий, тарифов, прогресса, разбора ошибок и работы с AI Tutor.",
  "Mobile web": "Мобильная версия",
  "Run short lessons, repeat weak words, or practice a phrase when you are away from the desk.": "Проходите короткие уроки, повторяйте слабые слова и тренируйте фразы, когда вы не за компьютером.",
  "Start quickly, receive reminders, send voice, and keep practicing from the same profile.": "Быстро начинайте практику, получайте напоминания, отправляйте голос и продолжайте из того же профиля.",
  "Daily loop": "Ежедневный цикл",
  "One short session always ends with the next useful step": "Каждая короткая сессия заканчивается понятным следующим шагом",
  "NERIVA is built around a simple loop: learn, use, review. Every module feeds the same progress profile.": "NERIVA построен на простом цикле: изучить, применить, повторить. Каждый модуль ведет в общий профиль прогресса.",
  "Learn the phrase": "Изучите фразу",
  "Get the meaning, grammar hint, natural example, and one focused prompt.": "Получите значение, подсказку по грамматике, естественный пример и один точный вопрос.",
  "Use it in context": "Используйте её в контексте",
  "Practice through a short roleplay, voice answer, or photo-based task.": "Отработайте фразу в короткой ролевой сцене, голосовом ответе или задании по фото.",
  "Review the weak spot": "Разберите слабое место",
  "Mistakes, notes, weak words, XP, and streak point to the next repetition.": "Ошибки, заметки, слабые слова, XP и серия подсказывают, что повторять дальше.",
  "Product modules": "Модули продукта",
  "Lessons, conversation, voice, and photos work as one system": "Уроки, диалоги, голос и фото работают как одна система",
  "Pronunciation practice highlights weak words, rhythm, and the next sentence to repeat.": "Тренировка произношения показывает слабые слова, ритм и следующую фразу для повтора.",
  "Photo Practice": "Практика по фото",
  "A menu or sign becomes translation, context, notes, and a short practice loop.": "Меню или вывеска превращаются в перевод, контекст, заметку и короткую практику.",
  "Review memory": "Память повторения",
  "Your weak spots stay visible until they become easy": "Слабые места остаются на виду, пока не станут проще",
  "Review is not a separate folder. It is the connective tissue between lessons, voice, notes, and Telegram practice.": "Повторение не лежит в отдельной папке. Оно связывает уроки, голос, заметки и практику в Telegram.",
  "Weak words": "Слабые слова",
  "Streak": "Серия",
  "Offline decks": "Офлайн-наборы",
  "Two entry points": "Две точки входа",
  "Use the same profile for deep sessions and quick practice": "Один профиль для длинных занятий и быстрой практики",
  "Open the dashboard for deep sessions, progress, pricing, mistakes, and longer AI Tutor work.": "Откройте веб-приложение для длинных занятий, прогресса, тарифов, ошибок и работы с AI Tutor.",
  "Open Web app": "Открыть веб-приложение",
  "Start quick practice, voice checks, reminders, and weak-word review from the same profile.": "Начинайте быструю практику, проверку голоса, напоминания и повторение слабых слов из одного профиля.",
  "Pricing": "Тарифы",
  "Start free, then unlock the daily practice limits you actually need": "Начните бесплатно, затем подключите нужные дневные лимиты",
  "Payment is available through Telegram Stars, YooKassa/SBP, TON, and USDT.": "Оплата доступна через Telegram Stars, YooKassa/СБП, TON и USDT.",
  "Try the loop": "Попробовать цикл",
  "Basic text practice for trying lessons, word training, notes, and progress without payment.": "Базовая текстовая практика: уроки, слова, заметки и прогресс без оплаты.",
  "Basic notes and phrasebook": "Заметки и разговорник",
  "Choose plan": "Выбрать тариф",
  "Daily speaking plan": "Ежедневная практика",
  "The main plan for daily progress with AI Tutor, roleplay, pronunciation, voice checks, and photo tools.": "Основной тариф: AI Tutor, ролевые диалоги, произношение, проверка голоса и работа с фото.",
  "Intensive preparation": "Интенсивная подготовка",
  "Higher limits for travel, work, exam preparation, and longer AI-dialogue sessions.": "Повышенные лимиты для поездок, работы, экзаменов и более длинных AI-диалогов.",
  "per month": "в месяц",
  "starter access": "стартовый доступ",
  "User stories": "Отзывы",
  "Built for concrete speaking moments": "Для конкретных речевых ситуаций",
  "Travel practice": "Практика перед поездкой",
  "I rehearsed check-in, cafe orders, and transport before the trip. The phrases stayed in notes for quick review.": "Перед поездкой я отработала регистрацию, заказ в кафе и транспорт. Фразы остались в заметках для быстрого повтора.",
  "English for work": "Английский для работы",
  "I use the web app for longer lessons and Telegram for weak words before calls. The same profile keeps it simple.": "Длинные уроки прохожу в веб-приложении, а слабые слова перед звонками повторяю в Telegram. Один профиль все упрощает.",
  "Pronunciation": "Произношение",
  "Voice practice shows which words sound weak and gives a better sentence to repeat right away.": "Голосовая практика показывает слабые слова и сразу дает более естественную фразу для повтора.",
  "Open the loop and run the first lesson today": "Запустите цикл и начните первый урок сегодня",
  "Use the web app for a full session or Telegram for a fast practice check.": "Используйте веб-приложение для полноценной сессии или Telegram для быстрой проверки.",
};

function stripGeneratedBlock(source) {
  const start = source.indexOf(GENERATED_START);
  const end = source.indexOf(GENERATED_END);
  if (start === -1 || end === -1 || end < start) return source;
  return `${source.slice(0, start).trimEnd()}\n`;
}

function extractGeneratedTranslations(source) {
  const marker = "const publicSiteGeneratedTranslations = ";
  const start = source.indexOf(marker);
  if (start === -1) return {};
  const jsonStart = start + marker.length;
  const jsonEnd = source.indexOf(";\nObject.entries(publicSiteGeneratedTranslations)", jsonStart);
  if (jsonEnd === -1) return {};
  try {
    return JSON.parse(source.slice(jsonStart, jsonEnd));
  } catch {
    return {};
  }
}

function normalize(value) {
  return String(value || "").replace(/\s+/g, " ").trim();
}

function hasLetter(value) {
  return /\p{L}/u.test(value);
}

function hasCyrillic(value) {
  return /[А-Яа-яЁё]/.test(value);
}

function hasLatinWord(value) {
  return /[A-Za-z]{3,}/.test(value);
}

function isOnlyProtectedName(value) {
  return NATIVE_LANGUAGE_NAMES.has(value)
    || /^(?:NERIVA|Premium|Platinum|Free|Telegram|Stars|TON|USDT|YooKassa|SBP|AI|FAQ|CEFR|OCR|RUB)$/i.test(value)
    || /^[@/#.]/.test(value)
    || /^[\w.-]+@[\w.-]+$/.test(value);
}

function isMeaningful(value) {
  const text = normalize(value);
  if (!text || text.length < 2) return false;
  if (!hasLetter(text)) return false;
  if (isOnlyProtectedName(text)) return false;
  if (text.length > 700) return false;
  if (/[<>]/.test(text)) return false;
  if (/\\[sS]|\/[gimsuy]*$|\$\s*\p{Nd}+/u.test(text)) return false;
  if (/^[a-z][a-z0-9_-]{2,}$/.test(text)) return false;
  if (/[{}]|=>|className=|style=|children=|useState|useEffect|window\.|document\.|querySelector|addEventListener|removeEventListener/.test(text)) return false;
  if (/\b(?:import|export|return|const|let|var|function|type|interface)\b/.test(text) && /[{}()[\];]/.test(text)) return false;
  if (text.includes("${") || text.includes("querySelector") || text.includes("classList")) return false;
  if (/^(?:target|href|src|class|id|dataset|style|function|return|const|let|var)\b/.test(text)) return false;
  if (/^(?:https?:|mailto:|tel:|\/|#|\.)/.test(text)) return false;
  if (/^[\w-]+(?:\s+[\w-]+){2,}$/.test(text) && /\b(?:bg|text|rounded|flex|grid|shadow|border|hover|items|justify|mx|mb|mt|px|py)\b/.test(text)) return false;
  return true;
}

function addCandidate(set, value) {
  const text = normalize(value);
  if (isMeaningful(text)) set.add(text);
}

function decodeHtmlEntities(value) {
  return String(value)
    .replace(/&copy;/g, "©")
    .replace(/&nbsp;/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, "\"");
}

function decodeJsLiteral(value) {
  return String(value)
    .replace(/\\n/g, " ")
    .replace(/\\r/g, " ")
    .replace(/\\t/g, " ")
    .replace(/\\`/g, "`")
    .replace(/\\"/g, "\"")
    .replace(/\\'/g, "'")
    .replace(/\\\\/g, "\\");
}

function collectQuotedLiterals(source) {
  const literals = [];
  const literalPattern = /(["'`])((?:\\[\s\S]|(?!\1)[\s\S])*?)\1/g;

  for (const match of source.matchAll(literalPattern)) {
    const quote = match[1];
    const raw = match[2];
    if (quote === "`" && raw.includes("${")) continue;
    if (quote !== "`" && /[\r\n]/.test(raw)) continue;
    literals.push(decodeJsLiteral(raw));
  }

  return literals;
}

function currentPublicSiteSource(file, source) {
  if (!file.endsWith("PublicSiteApp.tsx")) return source;

  let next = source;
  const legacyNavStart = next.indexOf("function SiteNav(");
  const legalHelpersStart = next.indexOf("function ensurePrivacyBotContact(");
  if (legacyNavStart !== -1 && legalHelpersStart !== -1 && legalHelpersStart > legacyNavStart) {
    next = `${next.slice(0, legacyNavStart)}\n${next.slice(legalHelpersStart)}`;
  }

  const legacyFooterStart = next.indexOf("function SiteFooter() ");
  if (legacyFooterStart !== -1) {
    next = next.slice(0, legacyFooterStart);
  }

  return next;
}

function collectCandidates() {
  const candidates = new Set();

  const collectHtml = (html) => {
    const withoutScripts = html
      .replace(/<script[\s\S]*?<\/script>/gi, " ")
      .replace(/<style[\s\S]*?<\/style>/gi, " ")
      .replace(/<!--[\s\S]*?-->/g, " ");

    for (const match of withoutScripts.matchAll(/>([^<>]+)</g)) {
      addCandidate(candidates, decodeHtmlEntities(match[1]));
    }

    for (const match of html.matchAll(/(?:alt|aria-label|title|placeholder|content)="([^"]+)"/g)) {
      addCandidate(candidates, decodeHtmlEntities(match[1]));
    }
  };

  for (const page of PAGES) {
    const html = fs.readFileSync(path.join(SITE_DIR, page), "utf8");
    collectHtml(html);

    for (const block of html.matchAll(/<script[\s\S]*?>([\s\S]*?)<\/script>/gi)) {
      for (const match of block[1].matchAll(/"((?:\\.|[^"\\])*)"/g)) {
        let value = match[1];
        try {
          value = JSON.parse(`"${value}"`);
        } catch {
          // Keep raw string literal content.
        }
        const text = normalize(decodeHtmlEntities(value));
        if (hasCyrillic(text)) addCandidate(candidates, text);
      }
    }
  }

  for (const file of REACT_SOURCES) {
    if (!fs.existsSync(file)) continue;
    const source = currentPublicSiteSource(file, fs.readFileSync(file, "utf8"));
    if (file.endsWith("legacyLegalContent.ts")) {
      for (const match of source.matchAll(/=\s*"((?:\\.|[^"\\])*)"/g)) {
        try {
          collectHtml(JSON.parse(`"${match[1]}"`));
        } catch {
          // Keep going; malformed strings are not useful translation candidates.
        }
      }
      continue;
    }

    const withoutImports = source
      .replace(/\bimport[\s\S]*?\bfrom\s+["'][^"']+["'];/g, " ")
      .replace(/\bimport\s+["'][^"']+["'];/g, " ");
    const withoutComments = withoutImports
      .replace(/\/\*[\s\S]*?\*\//g, " ")
      .replace(/\/\/.*$/gm, " ");

    for (const match of withoutComments.matchAll(/>([^<>]+)</g)) {
      addCandidate(candidates, decodeHtmlEntities(match[1]));
    }

    for (const value of collectQuotedLiterals(withoutImports)) {
      addCandidate(candidates, decodeHtmlEntities(value));
    }
  }

  const i18nSource = fs.readFileSync(I18N_PATH, "utf8");
  const landingBlock = i18nSource.slice(
    i18nSource.indexOf("const landingEnglishPhrases = ["),
    i18nSource.indexOf("landingEnglishPhrases.forEach")
  );
  for (const match of landingBlock.matchAll(/\["((?:\\.|[^"\\])*)",\s*"((?:\\.|[^"\\])*)"/g)) {
    try {
      addCandidate(candidates, JSON.parse(`"${match[1]}"`));
    } catch {
      addCandidate(candidates, match[1]);
    }
  }

  return [...candidates].sort((a, b) => a.localeCompare(b, "ru"));
}

function createContext(lang, phrasesSource, i18nSource) {
  const ctx = {
    console,
    TextEncoder,
    TextDecoder,
    URL,
    URLSearchParams,
    CustomEvent: class CustomEvent {
      constructor(type, init) {
        this.type = type;
        this.detail = init?.detail;
      }
    },
    NodeFilter: { SHOW_TEXT: 4 },
    location: {
      search: `?lang=${lang}`,
      href: `https://example.com/poliglot-ai.html?lang=${lang}`,
      host: "example.com",
    },
    localStorage: {
      getItem: () => lang,
      setItem: () => {},
    },
    document: {
      documentElement: { lang },
      body: { dataset: { page: "landing" } },
      addEventListener: () => {},
      querySelectorAll: () => [],
      createTreeWalker: () => ({ nextNode: () => false }),
    },
    window: {
      addEventListener: () => {},
      dispatchEvent: () => {},
    },
  };
  ctx.window.window = ctx.window;
  ctx.window.document = ctx.document;
  ctx.window.location = ctx.location;
  ctx.window.localStorage = ctx.localStorage;
  vm.runInNewContext(phrasesSource, ctx, { filename: PHRASES_PATH });
  vm.runInNewContext(i18nSource, ctx, { filename: I18N_PATH });
  return ctx.window.poliglotSiteI18n;
}

function protectText(source) {
  const placeholders = [];
  let text = source.replace(/\bNERIVA\b/g, "NERIVA");

  for (const pattern of PROTECTED_PATTERNS) {
    text = text.replace(pattern, (match) => {
      const token = `__P${placeholders.length}__`;
      placeholders.push(match === "NERIVA" ? "NERIVA" : match);
      return token;
    });
  }

  return {
    text,
    restore(value) {
      let restored = value;
      placeholders.forEach((original, index) => {
        const tokenPattern = new RegExp(`__\\s*P\\s*${index}\\s*__`, "gi");
        restored = restored.replace(tokenPattern, original);
      });
      return restored
        .replace(/(\p{L})(Telegram|Premium|Platinum|Free|Stars|TON|USDT|RUB|YooKassa|SBP|CEFR|OCR|AI)\b/gu, "$1 $2")
        .replace(/\b(Telegram|Premium|Platinum|Free|Stars|TON|USDT|RUB|YooKassa|SBP|CEFR|OCR|AI)(\p{L})/gu, "$1 $2")
        .replace(/\s+([,.!?;:%])/g, "$1")
        .replace(/([([{])\s+/g, "$1")
        .replace(/\s+([)\]}])/g, "$1")
        .replace(/\s+·\s+/g, " · ")
        .trim();
    },
  };
}

function googleSourceCode(source) {
  return hasCyrillic(source) ? "ru" : "en";
}

function needsTranslation(source, lang, current, english) {
  if (!isMeaningful(source)) return false;
  if (isOnlyProtectedName(source)) return false;
  if (hasTranslationLeak(current)) return true;
  if (lang === "ru") return !hasCyrillic(source) && hasLatinWord(source) && current === source;
  if (lang === "en") return hasCyrillic(source) && current === source;
  if (current === source) return true;
  if (english && current === english && hasLatinWord(english) && !isOnlyProtectedName(english)) return true;
  return false;
}

function hasTranslationLeak(value) {
  const text = String(value || "");
  return /\[\[?\s*\p{Nd}+\s*\]?\]?/u.test(text)
    || /_{1,2}\s*[PП]\s*\p{Nd}+\s*_{1,2}/iu.test(text)
    || /_{2}\s*[PП]\s*\p{Nd}+/iu.test(text)
    || /[PП]\s*\p{Nd}+\s*_{1,2}/iu.test(text)
    || /_[PП]\p{Nd}+_/iu.test(text)
    || /\$\s*\p{Nd}+/u.test(text);
}

async function translateBatch(items, targetGoogleCode) {
  const groups = new Map();
  for (const item of items) {
    const sl = googleSourceCode(item.source);
    if (sl === targetGoogleCode) {
      item.result = item.source.replace(/\bNERIVA\b/g, "NERIVA");
      continue;
    }
    const key = sl;
    if (!groups.has(key)) groups.set(key, []);
    groups.get(key).push(item);
  }

  for (const [sourceGoogleCode, group] of groups.entries()) {
    let offset = 0;
    while (offset < group.length) {
      const batch = [];
      let size = 0;
      while (offset < group.length && batch.length < 10 && size < 2800) {
        const item = group[offset++];
        const protectedItem = protectText(item.source);
        item.protectedItem = protectedItem;
        size += protectedItem.text.length;
        batch.push(item);
      }

      const parts = await translateMarkedBatch(batch, sourceGoogleCode, targetGoogleCode);

      for (let index = 0; index < batch.length; index += 1) {
        const item = batch[index];
        const raw = parts.get(index);
        let restored = raw ? item.protectedItem.restore(raw.replace(/^\[\[\d+\]\]\s*/, "")) : "";
        if (!restored || hasTranslationLeak(restored)) {
          const singleRaw = await fetchGoogle(item.protectedItem.text, sourceGoogleCode, targetGoogleCode);
          restored = item.protectedItem.restore(singleRaw.replace(/^\[\[\d+\]\]\s*/, ""));
        }
        batch[index].result = restored && !hasTranslationLeak(restored) ? restored : "";
      }

      if (TRANSLATION_PAUSE_MS > 0) {
        await new Promise((resolve) => setTimeout(resolve, TRANSLATION_PAUSE_MS));
      }
    }
  }
}

async function translateMarkedBatch(batch, sourceGoogleCode, targetGoogleCode) {
  const marked = batch.map((item, index) => `[[${index}]] ${item.protectedItem.text}`).join("\n");
  try {
    const translated = await fetchGoogle(marked, sourceGoogleCode, targetGoogleCode);
    return splitMarkedTranslation(translated, batch.length);
  } catch (error) {
    if (batch.length <= 1) {
      throw error;
    }
    const midpoint = Math.ceil(batch.length / 2);
    const left = await translateMarkedBatch(batch.slice(0, midpoint), sourceGoogleCode, targetGoogleCode);
    const right = await translateMarkedBatch(batch.slice(midpoint), sourceGoogleCode, targetGoogleCode);
    const merged = new Map(left);
    for (const [relativeIndex, value] of right.entries()) merged.set(relativeIndex + midpoint, value);
    return merged;
  }
}

async function fetchGoogle(text, sl, tl) {
  const url = new URL("https://translate.googleapis.com/translate_a/single");
  url.searchParams.set("client", "gtx");
  url.searchParams.set("sl", sl);
  url.searchParams.set("tl", tl);
  url.searchParams.set("dt", "t");
  url.searchParams.set("q", text);

  for (let attempt = 1; attempt <= 4; attempt += 1) {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), TRANSLATION_TIMEOUT_MS);
    try {
      const response = await fetch(url, { signal: controller.signal });
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      const data = await response.json();
      return (data?.[0] || []).map((part) => part?.[0] || "").join("");
    } catch (error) {
      if (attempt === 4) throw error;
      await new Promise((resolve) => setTimeout(resolve, attempt * 500));
    } finally {
      clearTimeout(timer);
    }
  }
  return text;
}

function splitMarkedTranslation(text, expectedCount) {
  const result = new Map();
  const marker = /\[\[\s*(\p{Nd}+)\s*\]\]/gu;
  const matches = [...text.matchAll(marker)];
  for (let i = 0; i < matches.length; i += 1) {
    const index = parseLocalizedNumber(matches[i][1]);
    const start = matches[i].index + matches[i][0].length;
    const end = i + 1 < matches.length ? matches[i + 1].index : text.length;
    if (index >= 0 && index < expectedCount) {
      result.set(index, text.slice(start, end).trim());
    }
  }
  return result;
}

function parseLocalizedNumber(value) {
  const digits = [];
  for (const char of String(value || "")) {
    const digit = localizedDigitValue(char);
    if (digit === null) return Number.NaN;
    digits.push(String(digit));
  }
  return digits.length ? Number(digits.join("")) : Number.NaN;
}

function localizedDigitValue(char) {
  const code = char.codePointAt(0);
  if (code === undefined) return null;
  const ranges = [
    0x0030, // ASCII
    0x0660, // Arabic-Indic
    0x06f0, // Extended Arabic-Indic
    0x0966, // Devanagari
    0x09e6, // Bengali
    0x0a66, // Gurmukhi
    0x0ae6, // Gujarati
    0x0b66, // Oriya
    0x0be6, // Tamil
    0x0c66, // Telugu
    0x0ce6, // Kannada
    0x0d66, // Malayalam
    0x0e50, // Thai
    0xff10, // Fullwidth
  ];
  for (const start of ranges) {
    if (code >= start && code <= start + 9) return code - start;
  }
  return null;
}

function buildCurrentSites(phrasesSource, i18nSource) {
  const sites = {};
  for (const [code] of LANGS) sites[code] = createContext(code, phrasesSource, i18nSource);
  return sites;
}

function shouldKeepEnglishTerm(source, value) {
  if (!value) return false;
  if (isOnlyProtectedName(value)) return true;
  return false;
}

async function main() {
  const phraseFile = fs.readFileSync(PHRASES_PATH, "utf8");
  const basePhraseFile = stripGeneratedBlock(phraseFile);
  const existingGenerated = extractGeneratedTranslations(phraseFile);
  const i18nSource = fs.readFileSync(I18N_PATH, "utf8");
  const candidates = collectCandidates();
  const candidateSet = new Set(candidates);
  const currentSites = buildCurrentSites(phraseFile, i18nSource);
  const englishSite = currentSites.en;
  const generated = {};
  for (const source of Object.keys(existingGenerated)) {
    if (!candidateSet.has(source)) continue;
    for (const [lang, value] of Object.entries(existingGenerated[source] || {})) {
      if (hasTranslationLeak(value)) continue;
      generated[source] ||= {};
      generated[source][lang] = value;
    }
  }
  const workByLang = new Map();

  for (const source of candidates) {
    const english = englishSite.translatePhrase(source);
    for (const [lang, googleCode] of LANGS) {
      const curated = lang === "ru" ? CURATED_RU_TRANSLATIONS[source] : undefined;
      if (curated && !hasTranslationLeak(curated)) {
        generated[source] ||= {};
        generated[source][lang] = curated;
        continue;
      }
      const current = currentSites[lang].translatePhrase(source);
      if (!needsTranslation(source, lang, current, english)) continue;
      if (!workByLang.has(lang)) workByLang.set(lang, []);
      workByLang.get(lang).push({ source, lang, googleCode });
    }
  }

  let total = 0;
  const languageJobs = [...workByLang.entries()];
  let jobIndex = 0;
  const workers = Array.from({ length: Math.min(TRANSLATION_CONCURRENCY, languageJobs.length) }, async () => {
    while (jobIndex < languageJobs.length) {
      const currentIndex = jobIndex++;
      const [lang, items] = languageJobs[currentIndex];
      total += items.length;
      console.log(`Translating ${items.length} strings for ${lang}...`);
      const googleCode = LANGS.find(([code]) => code === lang)?.[1] || lang;
      await translateBatch(items, googleCode);
      for (const item of items) {
        if (shouldKeepEnglishTerm(item.source, item.result)) continue;
        if (!item.result || hasTranslationLeak(item.result)) continue;
        generated[item.source] ||= {};
        generated[item.source][lang] = item.result;
      }
    }
  });
  await Promise.all(workers);

  const ordered = {};
  for (const source of Object.keys(generated).filter((source) => candidateSet.has(source)).sort((a, b) => a.localeCompare(b, "ru"))) {
    ordered[source] = {};
    for (const [code] of LANGS) {
      if (generated[source][code]) ordered[source][code] = generated[source][code];
    }
  }

  const block = [
    GENERATED_START,
    "window.poliglotPhraseTranslations = window.poliglotPhraseTranslations || {};",
    `const publicSiteGeneratedTranslations = ${JSON.stringify(ordered, null, 2)};`,
    "Object.entries(publicSiteGeneratedTranslations).forEach(([source, values]) => {",
    "  window.poliglotPhraseTranslations[source] = { ...(window.poliglotPhraseTranslations[source] || {}), ...values };",
    "});",
    GENERATED_END,
    "",
  ].join("\n");

  const nextPhrases = `${basePhraseFile.trimEnd()}\n\n${block}`;
  fs.writeFileSync(PHRASES_PATH, nextPhrases, "utf8");
  for (const mirrorPath of MIRROR_PHRASES_PATHS) {
    fs.mkdirSync(path.dirname(mirrorPath), { recursive: true });
    fs.writeFileSync(mirrorPath, nextPhrases, "utf8");
  }
  console.log(`Generated ${Object.keys(ordered).length} phrase entries, ${total} language values considered.`);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
