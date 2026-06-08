import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import {
  ArrowRight,
  BookOpen,
  Camera,
  CheckCircle,
  ChevronRight,
  FileText,
  MessageCircle,
  Mic,
  Moon,
  ShieldCheck,
  Sparkles,
  Star,
  Sun,
  Trophy,
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

const demoTabs = demoScreens.map((screen) => screen.tab);

const expandedInterfaceLanguageText =
  "доступны 35 языков интерфейса: العربية, বাংলা, Čeština, Deutsch, Ελληνικά, English, Español, Français, हिंदी, Magyar, Bahasa Indonesia, Italiano, 日本語, 한국어, Nederlands, Polski, Português, Română, Русский, svenska, தமிழ், తెలుగు, ภาษาไทย, Tagalog, Türkçe, Українською, Tiếng Việt, 中文, Тоҷикӣ, O‘zbekcha, Татарча, Հայերեն, Қазақша, Кыргызча и ქართული. Учебные словари расширяются отдельно.";

const learningLanguageShowcase = [
  ["English", "A1-C2", "26 млн учащихся"],
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
            <span className="eyebrow">AI-репетитор под вашу цель в web app, PWA и Telegram</span>
            <h1>Выберите цель, Poliglot AI соберет маршрут</h1>
            <p>
              Путешествие, работа, экзамен или разговорная речь превращаются в короткий курс: AI Tutor, диалоги, голос, фото-перевод, словарь ошибок и прогресс работают в одном профиле.
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
            <div className="hero-language-picker" aria-label="Популярные языки обучения">
              <h3>Я хочу изучать:</h3>
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
                <b>1</b> Урок
              </span>
              <span>
                <b>2</b> Практика
              </span>
              <span>
                <b>3</b> Роль
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
                <strong>Free</strong> старт в web, PWA и Telegram
              </span>
            </div>
          </motion.div>

          <motion.div initial={{ opacity: 0, x: 38 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.72, delay: 0.08 }} className="hero-demo">
            <div className="hero-demo__top">
              <div>
                <span>Сегодня</span>
                <strong>Маршрут дня</strong>
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
          <span className="eyebrow">Курсы под цель</span>
          <h2>Не просто уроки, а маршрут под вашу ситуацию</h2>
          <p>Выберите цель, а Poliglot AI соединит урок, диалог, голос, фото, слова и повторение в понятный сценарий.</p>
        </div>
        <div className="course-cards">
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
            профиль web + Telegram
          </span>
        </div>
      </section>

      <section id="features" className="landing-band feature-section">
        <div className="section-copy">
          <span className="eyebrow">Возможности</span>
          <h2>Чем эффективен подход Poliglot AI</h2>
          <p>Каждый блок ведет к действию: ответить, проговорить, повторить, сохранить фразу или разобрать ошибку.</p>
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
          <h2>Один профиль связывает сайт, PWA и Telegram</h2>
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

      <section id="pricing" className="pricing-section">
        <div className="section-copy">
          <span className="eyebrow">Тарифы</span>
          <h2>Начните бесплатно и расширьте лимиты, когда маршрут стал привычкой</h2>
          <p className="pricing-lead">Оплата доступна через Telegram Stars, YooKassa/СБП, TON и USDT. Premium и Platinum открывают голос, фото и больше практики.</p>
        </div>
        <div className="pricing-grid">
          <PlanCard name="Free" label="Попробовать маршрут" price="0 ₽" body="Старт без оплаты для короткой ежедневной практики." items={["5 уроков в день", "15 сообщений практики", "слова, заметки и базовый прогресс"]} />
          <PlanCard name="Premium" label="Регулярная учеба" oldPrice="1000 ₽" price="300 ₽" body="Основной режим для тех, кто занимается каждый день." items={["50 уроков в день", "200 сообщений практики", "20 voice до 30 секунд, голос и фото"]} featured />
          <PlanCard name="Platinum" label="Интенсив" oldPrice="2000 ₽" price="590 ₽" body="Максимальные лимиты для поездки, работы или экзамена." items={["100 уроков в день", "500 сообщений практики", "60 voice до 30 секунд, максимум AI-диалогов"]} />
        </div>
        <div className="payment-methods" aria-label="Способы оплаты">
          <span>Telegram Stars</span>
          <span>YooKassa/СБП</span>
          <span>TON</span>
          <span>USDT</span>
        </div>
      </section>

      <section id="reviews" className="landing-band reviews-section">
        <div className="section-copy">
          <span className="eyebrow">Отзывы</span>
          <h2>Сервис показывает продукт, а не обещания</h2>
          <p>Главное место на лендинге занимает живой сценарий: урок, голос, роль, ошибки и тарифы видны до регистрации.</p>
        </div>
        <div className="reviews-grid">
          {testimonials.map((item, index) => (
            <motion.article key={item.name} className="review-card" initial={{ opacity: 0, y: 18 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} transition={{ duration: 0.38, delay: index * 0.06 }}>
              <div>
                <span>{item.name.slice(0, 1)}</span>
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
          <h2>Сервис показывает, что повторить дальше</h2>
          <p>Ошибки, заметки, weak words, pronunciation history, XP, streak, награды и лидеры остаются в одном профиле между web app, PWA и Telegram.</p>
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
          <span className="eyebrow">Старт</span>
          <h2>Выберите цель и пройдите первый маршрут</h2>
        </div>
        <a className="hero-action hero-action--primary" href="/app/">
          Начать бесплатно <ArrowRight size={18} />
        </a>
      </section>
    </main>
  );
}

function PlanCard({ name, label, oldPrice, price, body, items, featured }: { name: string; label: string; oldPrice?: string; price: string; body: string; items: string[]; featured?: boolean }) {
  return (
    <article className={featured ? "is-featured" : undefined}>
      <span>{label}</span>
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
      <ul>
        {items.map((item) => (
          <li key={item}>
            <CheckCircle size={16} />
            {item}
          </li>
        ))}
      </ul>
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
        <p>AI-репетитор для короткой ежедневной языковой практики в web app, PWA и Telegram.</p>
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
