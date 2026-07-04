import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import {
  Activity,
  ArrowRight,
  BookOpen,
  BrainCircuit,
  Camera,
  CheckCircle,
  ChevronRight,
  FileText,
  Headphones,
  Languages,
  Laptop,
  LineChart,
  LockKeyhole,
  Menu,
  MessageCircle,
  MessagesSquare,
  Mic,
  Moon,
  Repeat2,
  ScanText,
  ShieldCheck,
  Sparkles,
  Star,
  Smartphone,
  Sun,
  Trophy,
  Zap,
  WifiOff,
  X,
} from "lucide-react";
import { PremiumHeroAnimation } from "@/components/ui/premium-hero-animation";
import { SparklesCore } from "@/components/ui/sparkles";
import { personalDataConsentDocumentHtml, privacyDocumentHtml, termsDocumentHtml, userAgreementDocumentHtml } from "./legacyLegalContent";
import { EnglishSparkLanding } from "./EnglishSparkLanding";
import { getLandingContent, normalizeLandingLocale, type LandingLocale } from "./landingContent";
import { staticSeoGuideLinks } from "./landingSeoContent";
import { SiteFooterSocial } from "./components/SocialLinks";

type PageId = "landing" | "privacy" | "terms" | "agreement" | "consent";
type SiteTheme = "light" | "dark";
type LegalPageId = Exclude<PageId, "landing">;

declare global {
  interface Window {
    dataLayer?: unknown[];
    ym?: (...args: unknown[]) => void;
    poliglotSiteI18n?: {
      apply: () => void;
      currentLanguage: () => string;
    };
    poliglotLegalDocumentsI18n?: {
      apply: () => void;
    };
  }
}

type PlanCardProps = {
  name: string;
  label: string;
  oldPrice?: string;
  price: string;
  body: string;
  limits: readonly (readonly [string, string])[];
  included: readonly string[];
  locked?: readonly string[];
  note: string;
  featured?: boolean;
};

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
  ["01", "Выберите цель", "Travel, work, exam или разговорная речь меняют маршрут и первые задания."],
  ["02", "Сделайте короткую сессию", "AI ведет через урок, диалог, голос, фото или словарь ошибок."],
  ["03", "Закрепите слабое место", "Notes, Mistakes, weak words, XP и streak подсказывают следующий повтор."],
];

const heroGoals = [
  ["Путешествия", "отель, кафе, аэропорт"],
  ["Работа", "созвоны, письма, small talk"],
  ["Экзамен", "лексика, грамматика, speaking"],
  ["Разговорная речь", "живые ответы без заучивания"],
] as const;

const courseRoutes = [
  ["AI Tutor Core", "Полный маршрут A1-C2 с уроком, практикой, повторением и XP.", "A1-C2"],
  ["Travel & Work", "Ролевые сценарии, фото-перевод, быстрые фразы и диалоги.", "роль + фото"],
  ["Voice Coach", "Shadowing, слабые слова, pronunciation score и история голоса.", "voice"],
] as const;

const cockpitLanes = [
  {
    icon: BrainCircuit,
    title: "AI Tutor",
    metric: "урок -> проверка",
    body: "Объясняет фразу, дает пример, просит ответ и сразу показывает, что улучшить дальше.",
  },
  {
    icon: MessagesSquare,
    title: "Roleplay",
    metric: "контекст роли",
    body: "Держит живой сценарий поездки, работы, экзамена или casual speaking без пустых реплик.",
  },
  {
    icon: Headphones,
    title: "Voice Coach",
    metric: "score + weak words",
    body: "Голосовые проверки, shadowing и TTS помогают слышать и исправлять произношение на ходу.",
  },
  {
    icon: ScanText,
    title: "Photo Tools",
    metric: "текст с картинки",
    body: "Меню, вывеска или задание становятся переводом, заметкой и практикой по контексту.",
  },
  {
    icon: Repeat2,
    title: "Mistake Loop",
    metric: "ошибка -> повтор",
    body: "Mistakes, Notes, weak words, XP и streak возвращают именно к слабым местам.",
  },
] as const;

const scenarioCards = [
  {
    title: "Путешествия без паники",
    tag: "Travel",
    body: "В аэропорту, в отеле, в кафе или у врача отрабатываются короткие фразы, которые нужны прямо сейчас.",
    modules: ["Roleplay", "Photo", "Phrasebook"],
    prompt: "Could you help me check in?",
    image: "/assets/scenarios/travel-ai-tutor.jpg",
  },
  {
    title: "Работа и созвоны",
    tag: "Work",
    body: "Письма, созвоны, self-intro и уточняющие вопросы превращаются в репетиции, после которых легче говорить по делу.",
    modules: ["AI Tutor", "Review", "Notes"],
    prompt: "Let me clarify the deadline.",
    image: "/assets/scenarios/work-ai-tutor.jpg",
  },
  {
    title: "Экзамен и уровень",
    tag: "Exam",
    body: "A1-C2, грамматика, listening и speaking собираются в маршрут с понятным прогрессом и видимыми шагами.",
    modules: ["Lesson", "Listening", "Mistakes"],
    prompt: "I agree with the statement because...",
    image: "/assets/scenarios/exam-ai-tutor.jpg",
  },
  {
    title: "Разговорная речь",
    tag: "Speaking",
    body: "AI не просто оценивает фразу, а предлагает естественную версию и следующий повтор, пока речь не станет увереннее.",
    modules: ["Voice", "Shadowing", "XP"],
    prompt: "I have been trying to say...",
    image: "/assets/scenarios/speaking-ai-tutor.jpg",
  },
] as const;

const deviceNodes = [
  {
    icon: Laptop,
    title: "Web app",
    body: "Большой экран для уроков, прогресса, тарифов, словаря ошибок и длинных сессий.",
  },
  {
    icon: Smartphone,
    title: "PWA и mobile web",
    body: "На телефоне удобно пройти короткий урок, повторить слова или потренировать фразу в дороге.",
  },
  {
    icon: MessageCircle,
    title: "Telegram",
    body: "Тот же профиль, быстрый старт, уведомления, голос и практика без потери данных.",
  },
] as const;

const planCards = [
  {
    name: "Free",
    label: "Начальный",
    price: "0 ₽",
    body: "Базовое текстовое обучение, тренировка слов, Phrasebook и обзор прогресса. AI Tutor, аудирование, произношение и голосовые проверки открываются в Premium.",
    limits: [
      ["Уроки", "5 уроков в день"],
      ["Практика", "15 сообщений практики"],
      ["Голос", "голосовые недоступны"],
    ],
    included: ["ежедневная привычка и стартовые уроки", "базовая тренировка слов", "заметки, phrasebook и обзор прогресса"],
    locked: ["Уроки с AI Tutor", "Аудирование и произношение", "Проверка голоса и фото-инструменты"],
    note: "Подходит для знакомства с продуктом без оплаты.",
  },
  {
    name: "Premium",
    label: "Регулярная учеба",
    oldPrice: "1000 ₽",
    price: "300 ₽",
    body: "Основной режим для ежедневной практики: AI Tutor, аудирование, произношение, AI-уроки, проверки голоса, фото-инструменты и расширенные дневные лимиты.",
    limits: [
      ["Уроки", "50 уроков в день"],
      ["Практика", "200 сообщений практики"],
      ["Голос", "20 голосовых до 30 секунд"],
    ],
    included: [
      "голос в текст и перевод услышанного",
      "перевод текста с картинки",
      "практика по контексту голоса или фото",
      "словарь ошибок, notes, XP и streak",
    ],
    note: "Лучший выбор для стабильного ежедневного обучения.",
    featured: true,
  },
  {
    name: "Platinum",
    label: "Интенсив",
    oldPrice: "2000 ₽",
    price: "590 ₽",
    body: "AI Tutor с максимальными дневными лимитами, глубокой ролевой практикой, интенсивным повторением и максимальной практикой голоса и произношения.",
    limits: [
      ["Уроки", "100 уроков в день"],
      ["Практика", "500 сообщений практики"],
      ["Голос", "60 голосовых до 30 секунд"],
    ],
    included: ["максимальные дневные лимиты", "больше практики по голосу или фото", "интенсивное повторение слабых мест", "лучший режим для плотной ежедневной учебы"],
    note: "Для поездки, работы, экзамена или очень плотного темпа.",
  },
] satisfies readonly PlanCardProps[];

const testimonials = [
  {
    name: "Анна",
    avatar: "/assets/testimonials/anna.jpg",
    role: "готовится к поездке",
    text: "Перед вылетом я прогнала check-in, кафе и транспорт. В заметках остались именно те фразы, которые потом пригодились в отеле.",
  },
  {
    name: "Марат",
    avatar: "/assets/testimonials/marat.jpg",
    role: "учит английский для работы",
    text: "Перед созвоном я репетирую self-intro и вопросы по срокам. После урока вижу ошибки, а слабые слова повторяю в Telegram.",
  },
  {
    name: "София",
    avatar: "/assets/testimonials/sofia.jpg",
    role: "тренирует произношение",
    text: "Я слышу, где съедаю окончания, и сразу повторяю более естественную фразу. Это спокойнее, чем просто учить список слов.",
  },
];

const faqs = [
  ["Можно заниматься только в Telegram?", "Да. Telegram и web app используют один учебный профиль, Premium и прогресс."],
  ["Чем Premium отличается от Free?", "Premium открывает больше уроков и практики, голос, фото-перевод и расширенные лимиты AI."],
  ["Сохраняются ли заметки и ошибки?", "Да. Notes, Mistakes, weak words, voice history, XP и streak закрепляются за аккаунтом."],
  ["Можно ли установить на телефон?", "Да. Web app работает как PWA и адаптирован отдельно под desktop и mobile web."],
];

const demoScreens = [
  {
    tab: "Урок",
    avatar: "AI",
    intro: "Сначала NERIVA дает фразу, смысл, пример и короткий next step.",
    label: "Урок",
    output: "I would like to book a table for tonight.",
    hint: "Фраза сразу попадает в маршрут и повторение.",
  },
  {
    tab: "Диалог",
    avatar: "DL",
    intro: "Роль держит контекст и просит ответить так, как в реальной ситуации.",
    label: "Hotel check-in",
    output: "Could you help me check in? I have a reservation.",
    hint: "AI исправляет ответ и дает более естественную версию.",
  },
  {
    tab: "Голос",
    avatar: "VO",
    intro: "Произнесите фразу, получите score, weak words и повтор для shadowing.",
    label: "Голос",
    output: "Please speak a little slower.",
    hint: "История голоса показывает, что стало лучше.",
  },
  {
    tab: "Фото",
    avatar: "PH",
    intro: "Сфотографируйте меню, вывеску или задание и получите учебный сценарий.",
    label: "Фото",
    output: "No peanuts, please. How spicy is this dish?",
    hint: "Фото превращается в перевод, заметку и practice prompt.",
  },
] as const;

const demoTabs = demoScreens.map((screen) => screen.tab);

const expandedInterfaceLanguageText =
  "доступны 35 языков интерфейса: العربية, বাংলা, Čeština, Deutsch, Ελληνικά, English, Español, Français, हिंदी, Magyar, Bahasa Indonesia, Italiano, 日本語, 한국어, Nederlands, Polski, Português, Română, Русский, svenska, தமிழ், తెలుగు, ภาษาไทย, Tagalog, Türkçe, Українською, Tiếng Việt, 中文, Тоҷикӣ, O‘zbekcha, Татарча, Հայերեն, Қазақша, Кыргызча и ქართული. Учебные словари расширяются отдельно.";

const learningLanguageShowcase = [
  ["English", "A1-C2", "маршруты для речи"],
  ["Español", "A1-C2", "диалоги для поездок"],
  ["Deutsch", "A1-C2", "работа и экзамены"],
  ["Français", "A1-C2", "аудирование и речь"],
  ["Italiano", "A1-C2", "живые сценарии"],
  ["中文", "A1-C2", "слова и практика"],
  ["日本語", "A1-C2", "уроки и повторение"],
  ["한국어", "A1-C2", "лексика и речь"],
] as const;

function expandLegalLanguageCopy(html: string) {
  return html
    .replace(/доступны 20 языков:[^;]+; перечень языков может расширяться\./g, expandedInterfaceLanguageText)
    .replace(/Сейчас в сервисе доступны 20 языков:[^;]+; перечень языков может расширяться\./g, `Сейчас ${expandedInterfaceLanguageText}`)
    .replace(/20 языков сейчас/g, "35 языков интерфейса")
    .replace(/20 языков/g, "35 языков интерфейса");
}

function getPage(): PageId {
  const raw = document.documentElement.dataset.sitePage || "";
  if (raw === "privacy" || location.pathname.includes("privacy")) return "privacy";
  if (raw === "terms" || location.pathname.includes("terms")) return "terms";
  if (raw === "agreement" || location.pathname.includes("agreement")) return "agreement";
  if (raw === "consent" || location.pathname.includes("consent")) return "consent";
  return "landing";
}

function getInitialSiteLanguage() {
  const requested = new URLSearchParams(location.search).get("lang");
  return requested || localStorage.getItem("poliglot_site_language") || document.documentElement.lang || "ru";
}

function useLegacySiteI18n(page: PageId, theme: SiteTheme) {
  useEffect(() => {
    document.body.dataset.page = page;
    const apply = () => {
      window.poliglotSiteI18n?.apply();
      window.poliglotLegalDocumentsI18n?.apply();
    };
    const timers = [0, 120, 320, 700, 1200].map((delay) => window.setTimeout(apply, delay));
    window.addEventListener("poliglot-language-change", apply);
    return () => {
      timers.forEach((timer) => window.clearTimeout(timer));
      window.removeEventListener("poliglot-language-change", apply);
    };
  }, [page, theme]);
}

function useSiteTheme() {
  const [theme, setTheme] = useState<SiteTheme>(() => (localStorage.getItem("poliglot-site-theme") === "dark" ? "dark" : "light"));
  useEffect(() => {
    document.documentElement.dataset.siteTheme = theme;
    localStorage.setItem("poliglot-site-theme", theme);
  }, [theme]);
  return [theme, setTheme] as const;
}

function useSiteLanguage() {
  const [language, setLanguage] = useState<LandingLocale>(() => normalizeLandingLocale(getInitialSiteLanguage()));

  useEffect(() => {
    const sync = () => {
      const nextLanguage = normalizeLandingLocale(window.poliglotSiteI18n?.currentLanguage?.() || getInitialSiteLanguage());
      document.documentElement.lang = nextLanguage;
      localStorage.setItem("poliglot_site_language", nextLanguage);
      setLanguage(nextLanguage);
    };

    sync();
    window.addEventListener("poliglot-language-change", sync);
    window.addEventListener("storage", sync);
    return () => {
      window.removeEventListener("poliglot-language-change", sync);
      window.removeEventListener("storage", sync);
    };
  }, []);

  return language;
}

export function PublicSiteApp() {
  const page = getPage();
  const [theme, setTheme] = useSiteTheme();
  const language = useSiteLanguage();
  const effectiveTheme = page === "landing" ? "dark" : theme;

  useEffect(() => {
    if (page === "landing") {
      document.documentElement.dataset.siteTheme = "dark";
    }
  }, [page]);

  useLegacySiteI18n(page, effectiveTheme);
  const landingCopy = getLandingContent(language);

  return (
    <div className="public-shell">
      <SiteNavDrawer page={page} theme={effectiveTheme} copy={landingCopy.nav} showThemeToggle={page !== "landing"} onThemeToggle={() => setTheme(theme === "dark" ? "light" : "dark")} />
      {page === "landing" ? <EnglishSparkLanding language={language} /> : <LegalPageV2 page={page} />}
      <SiteFooterEnglish showSeoGuides={page === "landing"} copy={landingCopy.nav} language={language} />
      <CookieConsentBanner />
    </div>
  );
}

function LandingLanguageOptions() {
  return (
    <>
      <option value="ru">Русский</option>
      <option value="en">English</option>
    </>
  );
}

function SiteNavDrawer({ page, theme, copy, showThemeToggle, onThemeToggle }: { page: PageId; theme: SiteTheme; copy: ReturnType<typeof getLandingContent>["nav"]; showThemeToggle: boolean; onThemeToggle: () => void }) {
  const [isOpen, setIsOpen] = useState(false);
  const themeLabel = theme === "dark" ? copy.light : copy.dark;

  return (
    <header className="public-nav public-nav--drawer">
      <a className="public-brand" href="/poliglot-ai.html" aria-label="NERIVA">
        <span className="public-brand__logo">
          <img src="/assets/brand-logo-mini.png" alt="" aria-hidden="true" />
        </span>
        <strong>NERIVA</strong>
      </a>
      <div className="public-nav__actions">
        <select data-site-language-select aria-label="Language">
          {page === "landing" ? <LandingLanguageOptions /> : null}
        </select>
        {showThemeToggle ? (
          <button className="nav-theme-toggle" type="button" onClick={onThemeToggle} aria-label={themeLabel}>
            {theme === "dark" ? <Sun size={18} /> : <Moon size={18} />}
          </button>
        ) : null}
        <button className="nav-burger" type="button" aria-expanded={isOpen} aria-controls="public-nav-drawer" onClick={() => setIsOpen((value) => !value)}>
          {isOpen ? <X size={20} /> : <Menu size={20} />}
          <span>{copy.menu}</span>
        </button>
      </div>
      <div className={isOpen ? "nav-drawer is-open" : "nav-drawer"} id="public-nav-drawer" aria-hidden={!isOpen}>
        <nav aria-label="Primary navigation">
          <a href="/poliglot-ai.html#features" onClick={() => setIsOpen(false)}>
            {copy.features}
          </a>
          <a href="/poliglot-ai.html#pricing" onClick={() => setIsOpen(false)}>
            {copy.pricing}
          </a>
          <a href="/poliglot-ai.html#faq" onClick={() => setIsOpen(false)}>
            {copy.faq}
          </a>
          <a className={page === "privacy" ? "is-active" : undefined} href="/privacy.html" onClick={() => setIsOpen(false)}>
            <span data-legal-nav="privacy">{copy.privacy}</span>
          </a>
          <a className={page === "terms" ? "is-active" : undefined} href="/terms.html" onClick={() => setIsOpen(false)}>
            <span data-legal-nav="terms">{copy.terms}</span>
          </a>
          <a className={page === "agreement" ? "is-active" : undefined} href="/agreement.html" onClick={() => setIsOpen(false)}>
            <span data-legal-nav="agreement">{copy.agreement}</span>
          </a>
          <a className={page === "consent" ? "is-active" : undefined} href="/consent.html" onClick={() => setIsOpen(false)}>
            <span data-legal-nav="consent">{copy.consent}</span>
          </a>
        </nav>
        <div className="nav-drawer__actions">
          {showThemeToggle ? (
            <button type="button" onClick={onThemeToggle}>
              {theme === "dark" ? <Sun size={18} /> : <Moon size={18} />}
              {themeLabel}
            </button>
          ) : null}
          <a href="/app/" onClick={() => setIsOpen(false)}>
            <Laptop size={18} />
            {copy.webApp}
          </a>
          <a href="https://t.me/NERIVAapp_bot" onClick={() => setIsOpen(false)}>
            <MessageCircle size={18} />
            {copy.telegram}
          </a>
        </div>
      </div>
    </header>
  );
}

function SiteNav({ page, theme, onThemeToggle }: { page: PageId; theme: SiteTheme; onThemeToggle: () => void }) {
  return (
    <header className="public-nav">
      <a className="public-brand" href="/poliglot-ai.html" aria-label="NERIVA">
        <span className="public-brand__logo">
          <img src="/assets/brand-logo-mini.png" alt="" aria-hidden="true" />
        </span>
        <strong>NERIVA</strong>
      </a>
      <nav aria-label="Основная навигация">
        <a href="/poliglot-ai.html#features" data-i18n="nav_features">
          Возможности
        </a>
        <a href="/poliglot-ai.html#pricing" data-i18n="nav_pricing">
          Тарифы
        </a>
        <a href="/poliglot-ai.html#faq">FAQ</a>
        <a className={page === "privacy" ? "is-active" : undefined} href="/privacy.html" data-i18n="nav_policy">
          Политика
        </a>
        <a className={page === "terms" ? "is-active" : undefined} href="/terms.html" data-i18n="nav_terms">
          Условия
        </a>
      </nav>
      <div className="public-nav__actions">
        <select data-site-language-select aria-label="Язык" />
        <button className="nav-theme-toggle" type="button" onClick={onThemeToggle} aria-label={theme === "dark" ? "Включить светлую тему" : "Включить темную тему"}>
          {theme === "dark" ? <Sun size={18} /> : <Moon size={18} />}
        </button>
        <a className="nav-pill" href="/app/" data-i18n="open_app">
          Открыть онлайн-приложение
        </a>
      </div>
    </header>
  );
}

function LandingPage() {
  const [activeTab, setActiveTab] = useState(0);
  const activeDemo = demoScreens[activeTab] ?? demoScreens[0];
  const activeLabel = activeDemo.tab;

  return (
    <main>
      <section className="landing-hero">
        <div className="landing-hero__matter">
          <PremiumHeroAnimation />
        </div>
        <div className="landing-hero__veil" aria-hidden="true" />

        <div className="landing-hero__inner">
          <motion.div initial={{ opacity: 0, y: 24 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.66 }} className="landing-hero__copy">
            <span className="eyebrow">Premium AI-репетитор для ежедневной практики</span>
            <h1>
              Говорите на новом языке с <span className="hero-nowrap">AI-репетитором</span>, который помнит ваши ошибки
            </h1>
            <p>
              Урок, диалог, произношение, фото-перевод, словарь ошибок и прогресс собираются в один профиль web, PWA и Telegram.
            </p>
            <div className="hero-actions">
              <a className="hero-action hero-action--primary" href="/app/">
                Начать бесплатно <ArrowRight size={18} />
              </a>
              <a className="hero-action hero-action--secondary" href="https://t.me/NERIVAapp_bot">
                Открыть Telegram
              </a>
            </div>
            <div className="hero-goals" aria-label="Цели обучения">
              {heroGoals.map(([goal, note]) => (
                <a key={goal} href="/app/">
                  <strong>{goal}</strong>
                  <span>{note}</span>
                </a>
              ))}
            </div>
            <div className="hero-language-picker" aria-label="Популярные языки обучения">
              <h3>Популярные языки обучения:</h3>
              <div>
                {learningLanguageShowcase.map(([language, level, note]) => (
                  <a key={language} href="/app/">
                    <strong>{language}</strong>
                    <span>{level}</span>
                    <small>{note}</small>
                  </a>
                ))}
              </div>
            </div>
            <div className="hero-path" aria-label="Учебный маршрут">
              <span>
                <b>1</b> Цель
              </span>
              <span>
                <b>2</b> Урок
              </span>
              <span>
                <b>3</b> Речь
              </span>
              <span>
                <b>4</b> Повтор
              </span>
            </div>
            <div className="hero-proof">
              <span>
                <strong>35</strong> языков интерфейса
              </span>
              <span>
                <strong>A1-C2</strong> уровни обучения
              </span>
              <span>
                <strong>Free</strong> старт без оплаты
              </span>
            </div>
          </motion.div>

          <motion.div initial={{ opacity: 0, x: 38 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.72, delay: 0.08 }} className="hero-demo">
            <div className="hero-demo__top">
              <div>
                <span>Сегодня</span>
                <strong>AI Tutor Cockpit</strong>
              </div>
              <span className="hero-demo__score">84/100</span>
            </div>
            <div className="hero-demo__tabs" role="tablist" aria-label="Демо функций">
              {demoTabs.map((tab, index) => (
                <button key={tab} type="button" className={activeTab === index ? "is-active" : undefined} onClick={() => setActiveTab(index)}>
                  {tab}
                </button>
              ))}
            </div>
            <motion.div key={activeLabel} initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.22 }} className="hero-demo__panel">
              <div className="demo-message">
                <div className="demo-avatar">{activeDemo.avatar}</div>
                <div>
                  <strong>NERIVA</strong>
                  <p>{activeDemo.intro}</p>
                </div>
              </div>
              <div className="demo-output">
                <span>{activeDemo.label}</span>
                <p>{activeDemo.output}</p>
              </div>
              <div className="demo-wave" aria-hidden="true">
                {Array.from({ length: 22 }, (_, index) => (
                  <i key={index} style={{ animationDelay: `${index * 0.035}s` }} />
                ))}
              </div>
              <div className="demo-hint">
                <CheckCircle size={17} />
                <span>{activeDemo.hint}</span>
              </div>
            </motion.div>
          </motion.div>
        </div>
      </section>

      <section className="course-strip" aria-label="Курсы NERIVA">
        <div className="course-strip__intro">
          <span className="eyebrow">AI Tutor</span>
          <h2>Маршруты под цель, а не бесконечная лента упражнений</h2>
          <p>Выберите режим для поездки, работы, экзамена или речи. NERIVA собирает урок, диалог, голос, фото и повторение в один управляемый цикл.</p>
        </div>
        <div className="course-cards">
          {courseRoutes.map(([title, body, badge], index) => (
            <motion.a key={title} href="/app/" className="course-card" initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} transition={{ duration: 0.42, delay: index * 0.05 }}>
              <span>{badge}</span>
              <h3>{title}</h3>
              <p>{body}</p>
              <b>
                Начать учиться <ChevronRight size={16} />
              </b>
            </motion.a>
          ))}
        </div>
        <div className="course-stats">
          <span>
            <strong>A1-C2</strong>
            уровни обучения
          </span>
          <span>
            <strong>35</strong>
            языков интерфейса
          </span>
          <span>
            <strong>1</strong>
            профиль web, PWA, Telegram
          </span>
        </div>
      </section>

      <section className="landing-band cockpit-section" aria-label="Как NERIVA ведет занятие">
        <div className="cockpit-shell">
          <div className="cockpit-copy">
            <span className="eyebrow">AI Tutor</span>
            <h2>Каждое занятие заканчивается понятным следующим шагом</h2>
            <p>
              NERIVA не оставляет вас один на один с упражнением. Он объясняет фразу, дает попробовать ее в диалоге, разбирает голос или фото и возвращает к ошибкам, пока слабое место не станет привычным.
            </p>
            <div className="cockpit-kpis" aria-label="Ключевые показатели NERIVA">
              <span>
                <strong>35</strong>
                языков интерфейса
              </span>
              <span>
                <strong>A1-C2</strong>
                маршруты уровня
              </span>
              <span>
                <strong>1</strong>
                профиль везде
              </span>
            </div>
          </div>

          <div className="cockpit-lanes">
            {cockpitLanes.map((lane, index) => {
              const Icon = lane.icon;
              return (
                <motion.article key={lane.title} className="cockpit-lane" initial={{ opacity: 0, x: -16 }} whileInView={{ opacity: 1, x: 0 }} viewport={{ once: true, amount: 0.35 }} transition={{ duration: 0.34, delay: index * 0.04 }}>
                  <Icon size={20} />
                  <div>
                    <strong>{lane.title}</strong>
                    <p>{lane.body}</p>
                  </div>
                  <span>{lane.metric}</span>
                </motion.article>
              );
            })}
          </div>

          <div className="cockpit-console" aria-label="Пример занятия AI Tutor">
            <div className="cockpit-console__top">
              <span>Live learning loop</span>
              <strong>84/100</strong>
            </div>
            <div className="cockpit-console__screen">
              <div className="cockpit-status">
                <Activity size={18} />
                <span>AI Tutor слушает ответ</span>
              </div>
              <div className="cockpit-dialogue">
                <p>Say it naturally: “Could you help me check in?”</p>
                <p>Лучше: “I have a reservation under my name.”</p>
              </div>
              <div className="cockpit-metric-grid">
                <span>
                  <Zap size={17} />
                  XP +12
                </span>
                <span>
                  <LineChart size={17} />
                  weak word fixed
                </span>
                <span>
                  <Languages size={17} />
                  перевод сохранен
                </span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="landing-band scenario-section">
        <div className="section-copy">
          <span className="eyebrow">Сценарии</span>
          <h2>Ситуации, ради которых язык нужен уже сегодня</h2>
          <p>В аэропорту, на созвоне, перед экзаменом или в обычном разговоре NERIVA дает не теорию ради теории, а короткую практику с фразами, голосом, фото и разбором ошибок.</p>
        </div>
        <div className="scenario-grid">
          {scenarioCards.map((scenario, index) => (
            <motion.article key={scenario.title} className="scenario-card" initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true, amount: 0.3 }} transition={{ duration: 0.36, delay: index * 0.05 }}>
              <figure className="scenario-card__visual">
                <img src={scenario.image} alt="" loading="lazy" />
                <span>{scenario.tag}</span>
              </figure>
              <h3>{scenario.title}</h3>
              <p>{scenario.body}</p>
              <div className="scenario-modules">
                {scenario.modules.map((module) => (
                  <span key={module}>{module}</span>
                ))}
              </div>
              <blockquote>{scenario.prompt}</blockquote>
            </motion.article>
          ))}
        </div>
      </section>

      <section id="features" className="landing-band feature-section">
        <div className="section-copy">
          <span className="eyebrow">Возможности</span>
          <h2>Не просто упражнения, а цикл AI-репетитора</h2>
          <p>Каждый блок ведет к действию: сказать, получить разбор, сохранить фразу и вернуться к слабому месту.</p>
        </div>
        <div className="feature-grid">
          {features.map((feature, index) => {
            const Icon = feature.icon;
            return (
              <motion.article key={feature.title} initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true, amount: 0.35 }} transition={{ duration: 0.42, delay: index * 0.04 }} className="feature-card">
                <Icon size={24} />
                <h3>{feature.title}</h3>
                <p>{feature.body}</p>
              </motion.article>
            );
          })}
        </div>
      </section>

      <section className="landing-band workflow-section">
        <div className="section-copy">
          <span className="eyebrow">Как это работает</span>
          <h2>Один профиль связывает web app, PWA и Telegram</h2>
          <p>Пользователь может начать на большом экране, продолжить как mobile web/PWA и вернуться в Telegram без потери Premium, прогресса и словаря ошибок.</p>
        </div>
        <div className="workflow-card">
          {workflow.map(([num, title, body]) => (
            <motion.div key={num} className="workflow-step" initial={{ opacity: 0, x: -16 }} whileInView={{ opacity: 1, x: 0 }} viewport={{ once: true }} transition={{ duration: 0.36 }}>
              <strong>{num}</strong>
              <div>
                <h3>{title}</h3>
                <p>{body}</p>
              </div>
            </motion.div>
          ))}
        </div>
      </section>

      <section className="device-flow" aria-label="Связка устройств NERIVA">
        <div className="device-flow__copy">
          <span className="eyebrow">Web + mobile + Telegram</span>
          <h2>Один профиль для web, mobile и Telegram</h2>
          <p>Начните большой урок на компьютере, повторите слабые слова с телефона и вернитесь в Telegram без потери прогресса, Premium-статуса и заметок.</p>
        </div>
        <div className="device-flow__grid">
          {deviceNodes.map((node, index) => {
            const Icon = node.icon;
            return (
              <motion.article key={node.title} className="device-flow__node" initial={{ opacity: 0, y: 16 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} transition={{ duration: 0.34, delay: index * 0.05 }}>
                <Icon size={24} />
                <h3>{node.title}</h3>
                <p>{node.body}</p>
                {index < deviceNodes.length - 1 ? <span className="device-flow__connector" aria-hidden="true" /> : null}
              </motion.article>
            );
          })}
        </div>
      </section>

      <section id="pricing" className="pricing-section">
        <div className="section-copy">
          <span className="eyebrow">Premium</span>
          <h2>Начните бесплатно, а для серьезной практики откройте голос, фото и расширенные лимиты</h2>
          <p className="payment-methods">Оплата: Telegram Stars и YooKassa/SBP.</p>
        </div>
        <div className="pricing-grid">
          {planCards.map((plan) => (
            <PlanCard key={plan.name} {...plan} />
          ))}
        </div>
      </section>

      <section id="reviews" className="landing-band reviews-section">
        <div className="section-copy">
          <span className="eyebrow">Отзывы</span>
          <h2>Истории пользователей, которым важно говорить увереннее</h2>
          <p>Короткие занятия помогают людям закрывать конкретные задачи: поездка, созвон, экзамен, произношение и повторение слабых слов.</p>
        </div>
        <div className="reviews-grid">
          {testimonials.map((item, index) => (
            <motion.article key={item.name} className="review-card" initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} transition={{ duration: 0.38, delay: index * 0.06 }}>
              <div>
                <img src={item.avatar} alt="" loading="lazy" />
                <div>
                  <strong>{item.name}</strong>
                  <small>{item.role}</small>
                </div>
              </div>
              <p>{item.text}</p>
              <div className="review-stars" aria-label="5 из 5">
                {Array.from({ length: 5 }, (_, star) => (
                  <Star key={star} size={15} fill="currentColor" />
                ))}
              </div>
            </motion.article>
          ))}
        </div>
      </section>

      <section className="progress-section">
        <div>
          <span className="eyebrow">Прогресс</span>
          <h2>Слабые места не теряются между устройствами</h2>
          <p>Ошибки, заметки, weak words, voice history, XP и streak подсказывают следующий шаг.</p>
        </div>
        <div className="progress-orbit" aria-hidden="true">
          <Trophy size={44} />
          <span>XP</span>
          <span>Streak</span>
          <span>Words</span>
        </div>
      </section>

      <section className="landing-band faq-section">
        <div className="section-copy">
          <span className="eyebrow">FAQ</span>
          <h2>Коротко по делу</h2>
        </div>
        <div className="faq-grid">
          {faqs.map(([question, answer]) => (
            <article key={question}>
              <h3>{question}</h3>
              <p>{answer}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="start-panel">
        <div>
          <span className="eyebrow">Первый маршрут</span>
          <h2>Откройте NERIVA и начните говорить с разбором ошибок уже сегодня</h2>
        </div>
        <a className="hero-action hero-action--primary" href="/app/">
          Начать бесплатно <ArrowRight size={18} />
        </a>
      </section>
    </main>
  );
}

function PlanCard({ name, label, oldPrice, price, body, limits, included, locked, note, featured }: PlanCardProps) {
  return (
    <article className={featured ? "plan-card is-featured" : "plan-card"}>
      <span className="plan-label">{label}</span>
      <h3>{name}</h3>
      <div className="price-line">
        {oldPrice ? <span className="plan-old-price">{oldPrice}</span> : null}
        <strong>{price}</strong>
        {oldPrice ? <small>в месяц</small> : null}
      </div>
      {oldPrice ? (
        <div className="price-tags">
          <span>70% экономии</span>
          <span>1 месяц запуска</span>
        </div>
      ) : null}
      <p>{body}</p>

      <div className="plan-limit-table" aria-label={`${name} лимиты`}>
        {limits.map(([metric, value]) => (
          <div key={metric}>
            <span>{metric}</span>
            <strong>{value}</strong>
          </div>
        ))}
      </div>

      <ul className="plan-included">
        {included.map((item) => (
          <li key={item}>
            <CheckCircle size={16} />
            {item}
          </li>
        ))}
      </ul>

      {locked?.length ? (
        <ul className="plan-locked">
          {locked.map((item) => (
            <li key={item}>
              <LockKeyhole size={15} />
              {item}
            </li>
          ))}
        </ul>
      ) : null}

      <small className="plan-note">{note}</small>
      <a href="/app/">
        Выбрать <ChevronRight size={16} />
      </a>
    </article>
  );
}

function ensurePrivacyBotContact(html: string) {
  let next = html.replace("<small>Telegram bot</small>@NERIVAapp_bot", "<small>Telegram bot</small><strong data-no-translate>@NERIVAapp_bot</strong>");
  next = next.replace(/<small>Telegram bot<\/small>\s*<\/a>/g, "<small>Telegram bot</small><strong data-no-translate>@NERIVAapp_bot</strong></a>");
  if (!/<a href="https:\/\/t\.me\/NERIVAapp_bot"[^>]*>[\s\S]*?@NERIVAapp_bot[\s\S]*?<\/a>/.test(next)) {
    next = next.replace(
      /(<div class="legal-contact-grid">[\s\S]*?)(<\/div>)/,
      '$1<a href="https://t.me/NERIVAapp_bot"><small>Telegram bot</small><strong data-no-translate>@NERIVAapp_bot</strong></a>$2',
    );
  }
  return next;
}

const legalPageMetaV2: Record<LegalPageId, { title: string; badge: string; description: string }> = {
  privacy: {
    title: "Политика обработки персональных данных",
    badge: "Защита данных",
    description: "Как NERIVA обрабатывает данные пользователей сайта, веб-приложения и Telegram-бота.",
  },
  terms: {
    title: "Условия использования",
    badge: "Правила сервиса",
    description: "Правила использования сайта, веб-приложения и Telegram-бота NERIVA.",
  },
  agreement: {
    title: "Пользовательское соглашение",
    badge: "Публичная оферта",
    description: "Публичное пользовательское соглашение для сайта, веб-приложения и Telegram-бота NERIVA.",
  },
  consent: {
    title: "Согласие на обработку персональных данных",
    badge: "152-ФЗ",
    description: "Отдельное согласие пользователя на обработку персональных данных в NERIVA.",
  },
};

const legalDocumentHtmlV2: Record<LegalPageId, string> = {
  privacy: ensurePrivacyBotContact(privacyDocumentHtml),
  terms: termsDocumentHtml,
  agreement: userAgreementDocumentHtml,
  consent: personalDataConsentDocumentHtml,
};

const legalOperatorNoticeHtmlV2 = `
<section id="operator-details" class="legal-operator-card">
  <h2>Оператор и реквизиты</h2>
  <div class="legal-contact-grid">
    <span><small>Оператор</small><strong>Самозанятый Чебан Денис Игоревич</strong></span>
    <span><small>ИНН</small><strong>ИНН 505017471160</strong></span>
    <span><small>Адрес</small><strong>г. Щёлково, ул. Сиреневая, 9к1, кв. 9</strong></span>
    <a href="mailto:support@neriva.ru"><small>Email</small><strong>support@neriva.ru</strong></a>
    <a href="https://t.me/NERIVAapp_bot"><small>Telegram bot</small><strong>@NERIVAapp_bot</strong></a>
  </div>
</section>`;

function LegalPageV2({ page }: { page: LegalPageId }) {
  const meta = legalPageMetaV2[page];
  const html = expandLegalLanguageCopy(`${legalOperatorNoticeHtmlV2}${legalDocumentHtmlV2[page]}`);

  return (
    <main className="legal-page">
      <section className="legal-hero">
        <div className="legal-hero__sparkles" aria-hidden="true">
          <SparklesCore background="transparent" minSize={0.3} maxSize={0.9} particleDensity={100} particleColor="#ffffff" speed={0.75} className="h-full w-full" />
        </div>
        <div className="legal-hero__content">
          <span className="eyebrow" data-legal-badge>{meta.badge}</span>
          <h1 data-legal-title>{meta.title}</h1>
          <p data-legal-description>{meta.description}</p>
          <div className="legal-pills">
            <span>35 языков интерфейса</span>
            <span>Сайт и Telegram</span>
            <span>Единый профиль</span>
            <span>Free, Premium и Platinum</span>
          </div>
        </div>
      </section>

      <section className="legal-layout">
        <aside className="legal-aside">
          <a href="/poliglot-ai.html">
            <Sparkles size={16} /> NERIVA
          </a>
          <a href="/privacy.html" className={page === "privacy" ? "is-active" : undefined}>
            <ShieldCheck size={16} /> <span data-legal-nav="privacy">Политика</span>
          </a>
          <a href="/terms.html" className={page === "terms" ? "is-active" : undefined}>
            <FileText size={16} /> <span data-legal-nav="terms">Условия</span>
          </a>
          <a href="/agreement.html" className={page === "agreement" ? "is-active" : undefined}>
            <FileText size={16} /> <span data-legal-nav="agreement">Пользовательское соглашение</span>
          </a>
          <a href="/consent.html" className={page === "consent" ? "is-active" : undefined}>
            <ShieldCheck size={16} /> <span data-legal-nav="consent">Согласие на обработку персональных данных</span>
          </a>
        </aside>
        <article className="legal-document-shell" dangerouslySetInnerHTML={{ __html: html }} />
      </section>
    </main>
  );
}

type CookieConsentChoice = {
  version: 1;
  necessary: true;
  analyticsMarketing: boolean;
};

type MetrikaFunction = ((...args: unknown[]) => void) & {
  a?: unknown[][];
  l?: number;
};

const cookieConsentStorageKey = "poliglot-cookie-consent";
const siteYandexMetrikaId = 109242081;

function serializeCookieConsent(analyticsMarketing: boolean) {
  const choice: CookieConsentChoice = {
    version: 1,
    necessary: true,
    analyticsMarketing,
  };

  return JSON.stringify(choice);
}

function parseCookieConsent(value: string | null): { saved: boolean; analyticsMarketing: boolean } {
  if (!value) return { saved: false, analyticsMarketing: false };
  if (value === "accepted") return { saved: true, analyticsMarketing: false };
  if (value === "necessary") return { saved: true, analyticsMarketing: false };

  try {
    const parsed = JSON.parse(value) as Partial<CookieConsentChoice>;
    if (parsed.version === 1 && parsed.necessary === true) {
      return { saved: true, analyticsMarketing: parsed.analyticsMarketing === true };
    }
  } catch {
    return { saved: false, analyticsMarketing: false };
  }

  return { saved: false, analyticsMarketing: false };
}

function hasSavedAnalyticsConsent() {
  if (typeof window === "undefined") return false;

  try {
    return parseCookieConsent(localStorage.getItem(cookieConsentStorageKey)).analyticsMarketing;
  } catch {
    return false;
  }
}

function loadYandexMetrika(counterId: number) {
  if (typeof window === "undefined" || typeof document === "undefined") return;
  if (!hasSavedAnalyticsConsent()) return;

  const existing = document.querySelector(`script[data-yandex-metrika-id="${counterId}"]`);
  if (existing) return;

  window.dataLayer = window.dataLayer || [];
  const metrika = window as Window & { ym?: MetrikaFunction };
  metrika.ym =
    metrika.ym ||
    function ymStub(...args: unknown[]) {
      const queue = (metrika.ym as MetrikaFunction);
      queue.a = queue.a || [];
      queue.a.push(args);
    };
  metrika.ym.l = Date.now();

  const script = document.createElement("script");
  script.async = true;
  script.src = `https://mc.yandex.ru/metrika/tag.js?id=${counterId}`;
  script.dataset.yandexMetrikaId = String(counterId);
  document.head.appendChild(script);

  metrika.ym(counterId, "init", {
    ssr: true,
    webvisor: true,
    clickmap: true,
    ecommerce: "dataLayer",
    referrer: document.referrer,
    url: location.href,
    accurateTrackBounce: true,
    trackLinks: true,
  });
}

function CookieConsentBanner() {
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [analyticsMarketing, setAnalyticsMarketing] = useState(false);
  const [visible, setVisible] = useState(() => {
    if (typeof window === "undefined") return false;
    return !parseCookieConsent(localStorage.getItem(cookieConsentStorageKey)).saved;
  });

  useEffect(() => {
    const consent = parseCookieConsent(localStorage.getItem(cookieConsentStorageKey));
    setVisible(!consent.saved);
    if (consent.analyticsMarketing) loadYandexMetrika(siteYandexMetrikaId);
  }, []);

  const saveChoice = (value: string, allowAnalytics: boolean) => {
    try {
      localStorage.setItem(cookieConsentStorageKey, value);
    } catch {
      // Ignore storage failures; the banner stays available for reading the legal links.
    }
    setVisible(false);
    if (allowAnalytics) loadYandexMetrika(siteYandexMetrikaId);
  };

  const saveNecessary = () => saveChoice(serializeCookieConsent(false), false);
  const saveAccepted = () => saveChoice(serializeCookieConsent(true), true);
  const saveSelected = () => saveChoice(serializeCookieConsent(analyticsMarketing), analyticsMarketing);

  if (!visible) return null;

  return (
    <section className="cookie-consent-banner" aria-label="Cookie consent">
      <div className="cookie-consent-banner__copy">
        <strong data-legal-cookie="title">Cookie и технические данные</strong>
        <p data-legal-cookie="body">Мы используем необходимые cookie и локальное хранилище для входа, языка, темы, безопасности, сохранения согласий и корректной работы сайта. Необязательные cookie применяются только после согласия.</p>
        <nav>
          <a href="/privacy.html" data-legal-cookie="privacy">Политика</a>
          <a href="/consent.html" data-legal-cookie="consent">Согласие</a>
        </nav>

        {settingsOpen && (
          <div className="cookie-consent-banner__settings" aria-label="Настройки cookie">
            <div className="cookie-consent-banner__setting-row">
              <div>
                <strong>Необходимые</strong>
                <span>Всегда активны для входа, безопасности, языка и сохранения согласий.</span>
              </div>
              <span className="cookie-consent-banner__required">Всегда активны</span>
            </div>

            <div className="cookie-consent-banner__setting-row">
              <div>
                <strong>Аналитические/маркетинговые</strong>
                <span>Помогают понимать работу сайта и улучшать продвижение, если такие инструменты включены.</span>
              </div>
              <button
                type="button"
                className="cookie-consent-banner__switch"
                role="switch"
                aria-checked={analyticsMarketing}
                aria-label="Аналитические/маркетинговые"
                onClick={() => setAnalyticsMarketing((value) => !value)}
              >
                <span />
              </button>
            </div>
          </div>
        )}
      </div>
      <div className="cookie-consent-banner__actions">
        <button type="button" onClick={saveNecessary} data-legal-cookie="necessary">Только необходимые</button>
        <button type="button" onClick={() => setSettingsOpen((value) => !value)} data-legal-cookie="settings">Настроить</button>
        {settingsOpen && <button type="button" onClick={saveSelected} data-legal-cookie="save-selected">Сохранить выбор</button>}
        <button type="button" onClick={saveAccepted} data-legal-cookie="accept">Принять все cookie</button>
      </div>
    </section>
  );
}

function LegalPage({ page }: { page: "privacy" | "terms" }) {
  const isPrivacy = page === "privacy";
  const title = isPrivacy ? "Политика обработки персональных данных" : "Условия использования";
  const badge = isPrivacy ? "Защита данных" : "Правила сервиса";
  const description = isPrivacy
    ? "Как NERIVA обрабатывает Telegram ID, username и технические данные для работы сайта, PWA и Telegram-бота."
    : "Правила использования сайта, веб-приложения и Telegram-бота NERIVA.";
  const sourceHtml = isPrivacy ? ensurePrivacyBotContact(privacyDocumentHtml) : termsDocumentHtml;
  const html = expandLegalLanguageCopy(sourceHtml);

  return (
    <main className="legal-page">
      <section className="legal-hero">
        <div className="legal-hero__sparkles" aria-hidden="true">
          <SparklesCore background="transparent" minSize={0.3} maxSize={0.9} particleDensity={100} particleColor="#ffffff" speed={0.75} className="h-full w-full" />
        </div>
        <div className="legal-hero__content">
          <span className="eyebrow">{badge}</span>
          <h1>{title}</h1>
          <p>{description}</p>
          <div className="legal-pills">
            <span>35 языков интерфейса</span>
            <span>Сайт и Telegram</span>
            <span>Единый профиль</span>
            <span>Free, Premium и Platinum</span>
          </div>
        </div>
      </section>

      <section className="legal-layout">
        <aside className="legal-aside">
          <a href="/poliglot-ai.html">
            <Sparkles size={16} /> NERIVA
          </a>
          <a href="/privacy.html" className={isPrivacy ? "is-active" : undefined}>
            <ShieldCheck size={16} /> Политика
          </a>
          <a href="/terms.html" className={!isPrivacy ? "is-active" : undefined}>
            <FileText size={16} /> Условия
          </a>
        </aside>
        <article className="legal-document-shell" dangerouslySetInnerHTML={{ __html: html }} />
      </section>
    </main>
  );
}

function SiteFooterEnglish({ showSeoGuides, copy, language }: { showSeoGuides: boolean; copy: ReturnType<typeof getLandingContent>["nav"]; language: LandingLocale }) {
  const isRussian = language === "ru";
  const guideLocale = isRussian ? "ru" : "en";
  return (
    <footer className="site-footer">
      <div>
        <strong>NERIVA</strong>
        <p>{isRussian ? "AI-репетитор для уроков, roleplay, голоса, фото-практики и повтора ошибок в браузере и Telegram." : "AI Tutor for lessons, roleplay, voice, photo practice, and mistake review across web app and Telegram."}</p>
        <span>© 2026 NERIVA. All rights reserved.</span>
        <address>
          <strong>{isRussian ? "Контакты" : "Contacts"}</strong>
          <a href="/app/">neriva.ru/app</a>
          <a href="https://t.me/NERIVAapp_bot">@NERIVAapp_bot</a>
          <a href="mailto:support@neriva.ru">support@neriva.ru</a>
        </address>
      </div>
      <nav>
        <strong>{isRussian ? "Навигация" : "Navigation"}</strong>
        <a href="/poliglot-ai.html#features">{copy.features}</a>
        <a href="/poliglot-ai.html#pricing">{copy.pricing}</a>
        <a href="/poliglot-ai.html#faq">{copy.faq}</a>
        <a href="/app/">{copy.webApp}</a>
      </nav>
      <SiteFooterSocial title={isRussian ? "Соцсети" : "Social"} />
      {showSeoGuides ? (
        <nav className="seo-guides" aria-label={isRussian ? "Материалы" : "Guides"}>
          <strong>{isRussian ? "Материалы" : "Guides"}</strong>
          {staticSeoGuideLinks.map((item) => {
            const link = item[guideLocale];
            return (
              <a key={item.slug} href={link.href} data-no-localize-href>
                {link.label}
              </a>
            );
          })}
        </nav>
      ) : null}
      <nav>
        <strong>{isRussian ? "Документы" : "Documents"}</strong>
        <a href="/privacy.html"><span data-legal-nav="privacy">{copy.privacy}</span></a>
        <a href="/terms.html"><span data-legal-nav="terms">{copy.terms}</span></a>
        <a href="/agreement.html"><span data-legal-nav="agreement">{copy.agreement}</span></a>
        <a href="/consent.html"><span data-legal-nav="consent">{copy.consent}</span></a>
      </nav>
    </footer>
  );
}

function SiteFooter() {
  return (
    <footer className="site-footer">
      <div>
        <strong>NERIVA</strong>
        <p>Premium AI-репетитор для уроков, диалогов, голоса, фото-перевода и словаря ошибок в web app, PWA и Telegram.</p>
        <span>© 2026 NERIVA. Все права защищены.</span>
        <address>
          <strong>Контакты</strong>
          <a href="/app/">neriva.ru/app</a>
          <a href="https://t.me/NERIVAapp_bot">@NERIVAapp_bot</a>
          <a href="mailto:support@neriva.ru">support@neriva.ru</a>
        </address>
      </div>
      <nav>
        <strong>Навигация</strong>
        <a href="/poliglot-ai.html#features">Возможности</a>
        <a href="/poliglot-ai.html#pricing">Тарифы</a>
        <a href="/poliglot-ai.html#faq">FAQ</a>
        <a href="/app/">Приложение</a>
      </nav>
      <SiteFooterSocial title="Соцсети" />
      <nav>
        <strong>Документы</strong>
        <a href="/privacy.html">Политика обработки данных</a>
        <a href="/terms.html">Условия использования</a>
      </nav>
    </footer>
  );
}
