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
} from "lucide-react";
import { GenerativeArtScene } from "@/components/ui/anomalous-matter-hero";
import { SparklesCore } from "@/components/ui/sparkles";
import { privacyDocumentHtml, termsDocumentHtml } from "./legacyLegalContent";

type PageId = "landing" | "privacy" | "terms";
type SiteTheme = "light" | "dark";

declare global {
  interface Window {
    poliglotSiteI18n?: {
      apply: () => void;
      currentLanguage: () => string;
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
    body: "Объясняет фразу, дает пример, просит ответ и сразу показывает, что улучшить.",
  },
  {
    icon: MessagesSquare,
    title: "Roleplay",
    metric: "контекст роли",
    body: "Держит сценарий поездки, работы, экзамена или casual speaking без пустых реплик.",
  },
  {
    icon: Headphones,
    title: "Voice Coach",
    metric: "score + weak words",
    body: "Голосовые проверки, shadowing и TTS помогают слышать и исправлять произношение.",
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
    body: "Check-in, кафе, транспорт, врач и small talk тренируются как короткие живые сцены.",
    modules: ["Roleplay", "Photo", "Phrasebook"],
    prompt: "Could you help me check in?",
  },
  {
    title: "Работа и созвоны",
    tag: "Work",
    body: "Письма, созвоны, self-intro и уточняющие вопросы превращаются в безопасные репетиции.",
    modules: ["AI Tutor", "Review", "Notes"],
    prompt: "Let me clarify the deadline.",
  },
  {
    title: "Экзамен и уровень",
    tag: "Exam",
    body: "A1-C2, грамматика, listening и speaking собираются в маршрут с понятным прогрессом.",
    modules: ["Lesson", "Listening", "Mistakes"],
    prompt: "I agree with the statement because...",
  },
  {
    title: "Разговорная речь",
    tag: "Speaking",
    body: "AI не просто оценивает фразу, а предлагает естественную версию и следующий повтор.",
    modules: ["Voice", "Shadowing", "XP"],
    prompt: "I have been trying to say...",
  },
] as const;

const deviceNodes = [
  {
    icon: Laptop,
    title: "Web app",
    body: "Большой экран для AI Tutor, прогресса, тарифов, словаря ошибок и длинных сессий.",
  },
  {
    icon: Smartphone,
    title: "PWA и mobile web",
    body: "Отдельная mobile-верстка для быстрых уроков, повторений и практики в дороге.",
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
    label: "Попробовать маршрут",
    price: "0 ₽",
    body: "Базовое текстовое обучение, тренировка слов, Phrasebook и обзор прогресса. AI Tutor, аудирование, произношение и голосовые проверки открываются в Premium.",
    limits: [
      ["Уроки", "5 уроков в день"],
      ["Практика", "15 сообщений практики"],
      ["Голос", "голосовые недоступны"],
    ],
    included: ["ежедневная привычка и стартовые уроки", "базовая тренировка слов", "заметки, phrasebook и обзор прогресса"],
    locked: ["AI Tutor guided lessons", "Listening и pronunciation", "voice checks и photo tools"],
    note: "Подходит для знакомства с продуктом без оплаты.",
  },
  {
    name: "Premium",
    label: "Регулярная учеба",
    oldPrice: "1000 ₽",
    price: "300 ₽",
    body: "Основной режим для ежедневной практики: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools и расширенные дневные лимиты.",
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
    body: "AI Tutor с максимальными дневными лимитами, глубиной roleplay, интенсивным review и максимальной voice/pronunciation практикой.",
    limits: [
      ["Уроки", "100 уроков в день"],
      ["Практика", "500 сообщений практики"],
      ["Голос", "60 голосовых до 30 секунд"],
    ],
    included: ["максимальные дневные лимиты", "больше voice/photo-context практики", "интенсивный review слабых мест", "лучший режим для heavy daily learning"],
    note: "Для поездки, работы, экзамена или очень плотного темпа.",
  },
] satisfies readonly PlanCardProps[];

const testimonials = [
  {
    name: "Анна",
    initial: "A",
    role: "готовится к поездке",
    text: "Перед поездкой я тренирую check-in, кафе и транспорт. Нужная фраза сразу попадает в заметки и повторение.",
  },
  {
    name: "Марат",
    initial: "M",
    role: "учит английский для работы",
    text: "В web app вижу ошибки и XP, в Telegram быстро повторяю слова. Профиль один, поэтому ничего не теряется.",
  },
  {
    name: "София",
    initial: "S",
    role: "тренирует произношение",
    text: "Голосовой режим показывает конкретные слабые слова. Это ощущается как личный coach, а не обычный список упражнений.",
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
    intro: "Сначала Poliglot AI дает фразу, смысл, пример и короткий next step.",
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
  return "landing";
}

function useLegacySiteI18n(page: PageId, theme: SiteTheme) {
  useEffect(() => {
    document.body.dataset.page = page;
    const apply = () => window.poliglotSiteI18n?.apply();
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

export function PublicSiteApp() {
  const page = getPage();
  const [theme, setTheme] = useSiteTheme();
  useLegacySiteI18n(page, theme);

  return (
    <div className="public-shell">
      <SiteNav page={page} theme={theme} onThemeToggle={() => setTheme(theme === "dark" ? "light" : "dark")} />
      {page === "landing" ? <LandingPage /> : <LegalPage page={page} />}
      <SiteFooter />
    </div>
  );
}

function SiteNav({ page, theme, onThemeToggle }: { page: PageId; theme: SiteTheme; onThemeToggle: () => void }) {
  return (
    <header className="public-nav">
      <a className="public-brand" href="/poliglot-ai.html" aria-label="Poliglot AI">
        <span className="public-brand__logo">
          <img src="/assets/brand-logo-mini.png" alt="" aria-hidden="true" />
        </span>
        <strong>Poliglot AI</strong>
      </a>
      <nav aria-label="Основная навигация">
        <a href="/poliglot-ai.html#features" data-i18n="nav_features">
          Возможности
        </a>
        <a href="/poliglot-ai.html#pricing" data-i18n="nav_pricing">
          Тарифы
        </a>
        <a href="/poliglot-ai.html#reviews">Отзывы</a>
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
          <GenerativeArtScene animate color="#7bdcff" particleColor="#f5d27a" />
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
              <a className="hero-action hero-action--secondary" href="https://t.me/poliglot_ai_bot">
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
                  <strong>Poliglot AI</strong>
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

      <section className="course-strip" aria-label="Курсы Poliglot AI">
        <div className="course-strip__intro">
          <span className="eyebrow">Premium AI Tutor Cockpit</span>
          <h2>Маршруты под цель, а не бесконечная лента упражнений</h2>
          <p>Выберите режим для поездки, работы, экзамена или речи. Poliglot AI собирает урок, диалог, голос, фото и повторение в один управляемый цикл.</p>
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

      <section className="landing-band cockpit-section" aria-label="Premium AI Tutor Cockpit">
        <div className="cockpit-shell">
          <div className="cockpit-copy">
            <span className="eyebrow">AI Tutor Cockpit</span>
            <h2>Премиальная панель, где каждый модуль ведет к следующему действию</h2>
            <p>
              Сайт должен сразу показывать ценность продукта: не отдельные упражнения, а систему, где урок, речь, фото, ошибки и прогресс работают как один персональный преподаватель.
            </p>
            <div className="cockpit-kpis" aria-label="Ключевые показатели Poliglot AI">
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

          <div className="cockpit-console" aria-label="Пример панели AI Tutor">
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
          <h2>Лендинг продает не абстрактный AI, а понятные ситуации, где продукт нужен сегодня</h2>
          <p>Каждая цель раскрывает реальные функции: roleplay, фото, голос, ошибки, progress и повторение. Это помогает пользователю увидеть, за что он платит.</p>
        </div>
        <div className="scenario-grid">
          {scenarioCards.map((scenario, index) => (
            <motion.article key={scenario.title} className="scenario-card" initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true, amount: 0.3 }} transition={{ duration: 0.36, delay: index * 0.05 }}>
              <div className="scenario-card__visual" aria-hidden="true">
                <span>{scenario.tag}</span>
                <div>
                  <i />
                  <i />
                  <i />
                </div>
              </div>
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

      <section className="device-flow" aria-label="Связка устройств Poliglot AI">
        <div className="device-flow__copy">
          <span className="eyebrow">Web + mobile + Telegram</span>
          <h2>Отдельно красиво на desktop и mobile, но с одним учебным профилем</h2>
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
          <p className="payment-methods">Оплата: Telegram Stars, YooKassa/SBP, TON и USDT.</p>
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
          <h2>Премиальный опыт без лишнего шума</h2>
          <p>Лендинг показывает реальный учебный цикл: урок, диалог, голос, фото, ошибки, тарифы и единый прогресс до регистрации.</p>
        </div>
        <div className="reviews-grid">
          {testimonials.map((item, index) => (
            <motion.article key={item.name} className="review-card" initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} transition={{ duration: 0.38, delay: index * 0.06 }}>
              <div>
                <span>{item.initial}</span>
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
          <h2>Откройте Poliglot AI и начните говорить с разбором ошибок уже сегодня</h2>
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
  let next = html.replace("<small>Telegram bot</small>@poliglot_ai_bot", "<small>Telegram bot</small><strong data-no-translate>@poliglot_ai_bot</strong>");
  next = next.replace(/<small>Telegram bot<\/small>\s*<\/a>/g, "<small>Telegram bot</small><strong data-no-translate>@poliglot_ai_bot</strong></a>");
  if (!/<a href="https:\/\/t\.me\/poliglot_ai_bot"[^>]*>[\s\S]*?@poliglot_ai_bot[\s\S]*?<\/a>/.test(next)) {
    next = next.replace(
      /(<div class="legal-contact-grid">[\s\S]*?)(<\/div>)/,
      '$1<a href="https://t.me/poliglot_ai_bot"><small>Telegram bot</small><strong data-no-translate>@poliglot_ai_bot</strong></a>$2',
    );
  }
  return next;
}

function LegalPage({ page }: { page: "privacy" | "terms" }) {
  const isPrivacy = page === "privacy";
  const title = isPrivacy ? "Политика обработки персональных данных" : "Условия использования";
  const badge = isPrivacy ? "Защита данных" : "Правила сервиса";
  const description = isPrivacy
    ? "Как Poliglot AI обрабатывает Telegram ID, username и технические данные для работы сайта, PWA и Telegram-бота."
    : "Правила использования сайта, веб-приложения и Telegram-бота Poliglot AI.";
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
            <Sparkles size={16} /> Poliglot AI
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

function SiteFooter() {
  return (
    <footer className="site-footer">
      <div>
        <strong>Poliglot AI</strong>
        <p>Premium AI-репетитор для уроков, диалогов, голоса, фото-перевода и словаря ошибок в web app, PWA и Telegram.</p>
        <span>© 2026 Poliglot AI. Все права защищены.</span>
      </div>
      <nav>
        <strong>Навигация</strong>
        <a href="/poliglot-ai.html#features">Возможности</a>
        <a href="/poliglot-ai.html#pricing">Тарифы</a>
        <a href="/poliglot-ai.html#reviews">Отзывы</a>
        <a href="/app/">Приложение</a>
      </nav>
      <nav>
        <strong>Документы</strong>
        <a href="/privacy.html">Политика обработки данных</a>
        <a href="/terms.html">Условия использования</a>
      </nav>
      <address>
        <strong>Контакты</strong>
        <a href="/app/">poliglotai.ru/app</a>
        <a href="https://t.me/poliglot_ai_bot">@poliglot_ai_bot</a>
        <a href="https://t.me/AsaselD">@AsaselD</a>
        <a href="mailto:supportpoliglotai@gmail.com">supportpoliglotai@gmail.com</a>
      </address>
    </footer>
  );
}
