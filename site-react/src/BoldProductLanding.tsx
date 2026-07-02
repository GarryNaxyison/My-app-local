import { useEffect, useState } from "react";
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

const WEB_APP_HREF = "/app/";
const TELEGRAM_HREF = "https://t.me/NERIVAapp_bot";

const heroScreens = [
  {
    tab: "Урок",
    icon: BookOpen,
    label: "AI Tutor lesson",
    prompt: "I would like to book a table for tonight.",
    output: "NERIVA проверяет смысл, грамматику и дает короткое объяснение.",
    detail: "XP +12 · ошибка сохранена",
  },
  {
    tab: "Диалог",
    icon: MessageCircle,
    label: "Hotel check-in",
    prompt: "Could you help me check in? I have a reservation.",
    output: "Сценарий держит роль, исправляет фразу и предлагает естественный ответ.",
    detail: "Roleplay · travel",
  },
  {
    tab: "Голос",
    icon: Mic,
    label: "Voice Coach",
    prompt: "Please speak a little slower.",
    output: "Score 84/100. Weak words: speak, slower. Следующая фраза для shadowing.",
    detail: "18 voice frames",
  },
  {
    tab: "Фото",
    icon: Camera,
    label: "Photo Practice",
    prompt: "No peanuts, please. How spicy is this dish?",
    output: "Фото меню превращается в перевод, заметку и короткую практику по контексту.",
    detail: "OCR · phrasebook",
  },
] as const;

const dailySteps = [
  {
    title: "Получите фокус",
    body: "Цель занятия задает короткий урок: поездка, работа, экзамен или разговорная речь.",
    marker: "01",
  },
  {
    title: "Ответьте как удобно",
    body: "Введите текст, скажите голосом или добавьте фото с меню, вывеской или заданием.",
    marker: "02",
  },
  {
    title: "Повторите слабое место",
    body: "NERIVA сохраняет ошибку, weak words, заметку и следующий шаг для review.",
    marker: "03",
  },
] as const;

const modules = [
  {
    title: "AI Tutor",
    icon: BrainCircuit,
    body: "Урок ведет от примера к ответу, проверке, XP и следующему повторению.",
    visual: ["Lesson", "answer check", "XP +12"],
    sample: "I would like to book...",
  },
  {
    title: "Roleplay",
    icon: Bot,
    body: "Travel, work, exam и speaking сценарии держат контекст как живой диалог.",
    visual: ["Hotel", "work call", "exam speaking"],
    sample: "Could you help me check in?",
  },
  {
    title: "Voice Coach",
    icon: Headphones,
    body: "Произношение получает score, weak words, shadowing и следующую фразу.",
    visual: ["score 84", "weak words", "shadowing"],
    sample: "weak words: speak, slower",
  },
  {
    title: "Photo Practice",
    icon: ScanText,
    body: "Текст с картинки становится переводом, заметкой и практикой по контексту.",
    visual: ["photo", "translation", "practice"],
    sample: "No peanuts, please.",
  },
] as const;

const memoryNodes = [
  ["Mistakes", "исправленные фразы"],
  ["Weak words", "слова из голоса"],
  ["Notes", "личный phrasebook"],
  ["Review", "повторение"],
  ["Spelling", "точность письма"],
  ["Offline decks", "карточки без сети"],
] as const;

type Plan = {
  name: string;
  price: string;
  body: string;
  limits: readonly string[];
  featured?: boolean;
};

const plans: readonly Plan[] = [
  {
    name: "Free",
    price: "0 ₽",
    body: "Старт без оплаты: базовые уроки, сообщения практики, notes и обзор прогресса.",
    limits: ["5 уроков в день", "15 сообщений практики", "без голосовых проверок"],
  },
  {
    name: "Premium",
    price: "300 ₽",
    body: "Регулярная учеба с AI Tutor, voice, photo, расширенными лимитами и mistake loop.",
    limits: ["50 уроков в день", "200 сообщений практики", "20 голосовых до 30 секунд"],
    featured: true,
  },
  {
    name: "Platinum",
    price: "590 ₽",
    body: "Интенсивный режим с максимальными дневными лимитами для поездки, работы или экзамена.",
    limits: ["100 уроков в день", "500 сообщений практики", "60 голосовых до 30 секунд"],
  },
] as const;

const reviews = [
  {
    name: "Анна",
    detail: "готовится к поездке",
    text: "Перед вылетом я прогнала check-in, кафе и транспорт. В заметках остались фразы, которые потом пригодились в отеле.",
  },
  {
    name: "Марат",
    detail: "учит английский для работы",
    text: "Перед созвоном я репетирую self-intro и вопросы по срокам. После урока вижу ошибки и повторяю слабые слова.",
  },
  {
    name: "София",
    detail: "тренирует произношение",
    text: "Я слышу, где съедаю окончания, и сразу повторяю более естественную фразу. Это спокойнее, чем учить список слов.",
  },
] as const;

type BoldLandingLocale = "ru" | "en" | "zh";

const boldNodeSources = new WeakMap<Text, string>();

const boldTranslations: Record<Exclude<BoldLandingLocale, "ru">, Record<string, string>> = {
  en: {
    "Урок": "Lesson",
    "NERIVA проверяет смысл, грамматику и дает короткое объяснение.": "NERIVA checks meaning, grammar, and gives a short explanation.",
    "XP +12 · ошибка сохранена": "XP +12 · mistake saved",
    "Диалог": "Dialogue",
    "Сценарий держит роль, исправляет фразу и предлагает естественный ответ.": "The scenario keeps the role, corrects the phrase, and suggests a natural reply.",
    "Голос": "Voice",
    "Score 84/100. Weak words: speak, slower. Следующая фраза для shadowing.": "Score 84/100. Weak words: speak, slower. Next phrase for shadowing.",
    "Фото": "Photo",
    "Фото меню превращается в перевод, заметку и короткую практику по контексту.": "A menu photo becomes a translation, note, and short contextual practice.",
    "Получите фокус": "Get a focus",
    "Цель занятия задает короткий урок: поездка, работа, экзамен или разговорная речь.": "The session goal sets a short lesson: travel, work, exam, or speaking practice.",
    "Ответьте как удобно": "Answer your way",
    "Введите текст, скажите голосом или добавьте фото с меню, вывеской или заданием.": "Type text, speak by voice, or add a photo of a menu, sign, or exercise.",
    "Повторите слабое место": "Repeat the weak spot",
    "NERIVA сохраняет ошибку, weak words, заметку и следующий шаг для review.": "NERIVA saves the mistake, weak words, note, and next review step.",
    "Урок ведет от примера к ответу, проверке, XP и следующему повторению.": "The lesson leads from example to answer, check, XP, and the next repetition.",
    "Travel, work, exam и speaking сценарии держат контекст как живой диалог.": "Travel, work, exam, and speaking scenarios keep context like a live dialogue.",
    "Произношение получает score, weak words, shadowing и следующую фразу.": "Pronunciation gets a score, weak words, shadowing, and the next phrase.",
    "Текст с картинки становится переводом, заметкой и практикой по контексту.": "Text from an image becomes a translation, note, and contextual practice.",
    "исправленные фразы": "corrected phrases",
    "слова из голоса": "words from voice",
    "личный phrasebook": "personal phrasebook",
    "повторение": "review",
    "точность письма": "writing accuracy",
    "карточки без сети": "offline cards",
    "Старт без оплаты: базовые уроки, сообщения практики, notes и обзор прогресса.": "Start free: basic lessons, practice messages, notes, and progress overview.",
    "5 уроков в день": "5 lessons per day",
    "15 сообщений практики": "15 practice messages",
    "без голосовых проверок": "no voice checks",
    "Регулярная учеба с AI Tutor, voice, photo, расширенными лимитами и mistake loop.": "Regular study with AI Tutor, voice, photo, expanded limits, and the mistake loop.",
    "50 уроков в день": "50 lessons per day",
    "200 сообщений практики": "200 practice messages",
    "20 голосовых до 30 секунд": "20 voice messages up to 30 seconds",
    "Интенсивный режим с максимальными дневными лимитами для поездки, работы или экзамена.": "Intensive mode with maximum daily limits for travel, work, or an exam.",
    "100 уроков в день": "100 lessons per day",
    "500 сообщений практики": "500 practice messages",
    "60 голосовых до 30 секунд": "60 voice messages up to 30 seconds",
    "Анна": "Anna",
    "готовится к поездке": "preparing for travel",
    "Перед вылетом я прогнала check-in, кафе и транспорт. В заметках остались фразы, которые потом пригодились в отеле.": "Before the flight, I practiced check-in, cafes, and transport. The notes kept phrases I later used at the hotel.",
    "Марат": "Marat",
    "учит английский для работы": "learning English for work",
    "Перед созвоном я репетирую self-intro и вопросы по срокам. После урока вижу ошибки и повторяю слабые слова.": "Before a call, I rehearse my self-intro and deadline questions. After the lesson, I see mistakes and repeat weak words.",
    "София": "Sofia",
    "тренирует произношение": "practices pronunciation",
    "Я слышу, где съедаю окончания, и сразу повторяю более естественную фразу. Это спокойнее, чем учить список слов.": "I hear where I drop endings and immediately repeat a more natural phrase. It is calmer than memorizing a word list.",
    "AI Tutor для ежедневной практики": "AI Tutor for daily practice",
    "Говорите с": "Speak with an",
    "AI-репетитором": "AI tutor",
    ", который помнит ваши ошибки": "that remembers your mistakes",
    "Урок, диалог, голос, фото, заметки и повторение слабых мест собираются в один профиль Web app и Telegram.": "Lessons, dialogues, voice, photos, notes, and weak-spot review stay in one Web app and Telegram profile.",
    "Открыть Web app": "Open Web app",
    "Открыть Telegram": "Open Telegram",
    "языков интерфейса": "interface languages",
    "маршруты уровня": "level paths",
    "старт без оплаты": "free start",
    "Что вы делаете сегодня": "What you do today",
    "Короткий учебный цикл вместо бесконечной ленты упражнений": "A short learning loop instead of an endless exercise feed",
    "Каждый день начинается с конкретной задачи и заканчивается повторением того, что реально просело.": "Each day starts with a concrete task and ends by repeating what actually slipped.",
    "Исправление сохранено в Mistakes и Notes.": "The correction is saved to Mistakes and Notes.",
    "Четыре рабочих режима в одном учебном профиле": "Four working modes in one learning profile",
    "Каждый модуль выглядит как продуктовый инструмент, а не как абстрактная карточка возможностей.": "Each module looks like a product tool, not an abstract feature card.",
    "Ошибки становятся тренировочным материалом": "Mistakes become training material",
    "NERIVA связывает исправления, weak words, spelling, notes, offline decks, XP, streak, awards и daily bonus в понятный повтор.": "NERIVA connects corrections, weak words, spelling, notes, offline decks, XP, streak, awards, and daily bonus into a clear review loop.",
    "Две равные точки входа, один профиль обучения": "Two equal entry points, one learning profile",
    "Web app удобнее для длинных занятий и dashboard, Telegram — для быстрой практики, напоминаний, голоса и продолжения на ходу.": "The Web app is better for longer sessions and the dashboard; Telegram is better for fast practice, reminders, voice, and continuing on the go.",
    "Фокусные уроки, dashboard прогресса, длинные сессии, тарифы и полный словарь ошибок.": "Focused lessons, progress dashboard, longer sessions, pricing, and the full mistake dictionary.",
    "быстрая практика, voice, reminders и продолжение того же урока без потери прогресса.": "fast practice, voice, reminders, and continuing the same lesson without losing progress.",
    "Free для старта, Premium и Platinum для ежедневной практики": "Free to start, Premium and Platinum for daily practice",
    "Оплата: Stars, YooKassa, TON и USDT.": "Payment: Stars, YooKassa, TON, and USDT.",
    "Выбрать тариф": "Choose plan",
    "Короткие истории пользователей": "Short user stories",
    "Без рейтингов и обещаний результата: только реальные сценарии, где нужен язык.": "No ratings or outcome promises: only real scenarios where language is needed.",
    "Начать сегодня": "Start today",
    "Откройте NERIVA и начните занятие с разбором ошибок": "Open NERIVA and start a lesson with mistake review",
  },
  zh: {
    "Урок": "课程",
    "NERIVA проверяет смысл, грамматику и дает короткое объяснение.": "NERIVA 检查含义和语法，并给出简短解释。",
    "XP +12 · ошибка сохранена": "XP +12 · 错误已保存",
    "Диалог": "对话",
    "Сценарий держит роль, исправляет фразу и предлагает естественный ответ.": "场景保持角色，修正句子，并给出自然回答。",
    "Голос": "语音",
    "Score 84/100. Weak words: speak, slower. Следующая фраза для shadowing.": "得分 84/100。弱项词：speak、slower。下一句用于 shadowing。",
    "Фото": "照片",
    "Фото меню превращается в перевод, заметку и короткую практику по контексту.": "菜单照片会变成翻译、笔记和简短的情境练习。",
    "Получите фокус": "获得重点",
    "Цель занятия задает короткий урок: поездка, работа, экзамен или разговорная речь.": "学习目标会生成一节短课：旅行、工作、考试或口语。",
    "Ответьте как удобно": "用适合你的方式回答",
    "Введите текст, скажите голосом или добавьте фото с меню, вывеской или заданием.": "输入文字、用语音回答，或添加菜单、标牌、作业的照片。",
    "Повторите слабое место": "复习薄弱点",
    "NERIVA сохраняет ошибку, weak words, заметку и следующий шаг для review.": "NERIVA 保存错误、weak words、笔记和下一步 review。",
    "Урок ведет от примера к ответу, проверке, XP и следующему повторению.": "课程从示例到回答、检查、XP 和下一次复习。",
    "Travel, work, exam и speaking сценарии держат контекст как живой диалог.": "Travel、work、exam 和 speaking 场景像真实对话一样保持上下文。",
    "Произношение получает score, weak words, shadowing и следующую фразу.": "发音会得到 score、weak words、shadowing 和下一句。",
    "Текст с картинки становится переводом, заметкой и практикой по контексту.": "图片文字会变成翻译、笔记和情境练习。",
    "исправленные фразы": "已修正的短语",
    "слова из голоса": "语音中的词",
    "личный phrasebook": "个人 phrasebook",
    "повторение": "复习",
    "точность письма": "书写准确度",
    "карточки без сети": "离线卡片",
    "Старт без оплаты: базовые уроки, сообщения практики, notes и обзор прогресса.": "免费开始：基础课程、练习消息、notes 和进度概览。",
    "5 уроков в день": "每天 5 节课",
    "15 сообщений практики": "15 条练习消息",
    "без голосовых проверок": "不含语音检查",
    "Регулярная учеба с AI Tutor, voice, photo, расширенными лимитами и mistake loop.": "使用 AI Tutor、voice、photo、更高限额和 mistake loop 进行日常学习。",
    "50 уроков в день": "每天 50 节课",
    "200 сообщений практики": "200 条练习消息",
    "20 голосовых до 30 секунд": "20 条最长 30 秒语音",
    "Интенсивный режим с максимальными дневными лимитами для поездки, работы или экзамена.": "用于旅行、工作或考试的高强度模式，拥有最高每日限额。",
    "100 уроков в день": "每天 100 节课",
    "500 сообщений практики": "500 条练习消息",
    "60 голосовых до 30 секунд": "60 条最长 30 秒语音",
    "Анна": "Anna",
    "готовится к поездке": "准备旅行",
    "Перед вылетом я прогнала check-in, кафе и транспорт. В заметках остались фразы, которые потом пригодились в отеле.": "出发前我练习了 check-in、咖啡馆和交通。笔记里留下了后来在酒店用到的短语。",
    "Марат": "Marat",
    "учит английский для работы": "为工作学习英语",
    "Перед созвоном я репетирую self-intro и вопросы по срокам. После урока вижу ошибки и повторяю слабые слова.": "会议前我练习 self-intro 和截止日期问题。课后我能看到错误并复习弱项词。",
    "София": "Sofia",
    "тренирует произношение": "练习发音",
    "Я слышу, где съедаю окончания, и сразу повторяю более естественную фразу. Это спокойнее, чем учить список слов.": "我能听出哪里吞掉了词尾，并马上重复更自然的句子。这比背单词列表更轻松。",
    "AI Tutor для ежедневной практики": "用于每日练习的 AI Tutor",
    "Говорите с": "和",
    "AI-репетитором": "AI 导师",
    ", который помнит ваши ошибки": "一起练习，它会记住你的错误",
    "Урок, диалог, голос, фото, заметки и повторение слабых мест собираются в один профиль Web app и Telegram.": "课程、对话、语音、照片、笔记和薄弱点复习都连接到同一个 Web app 和 Telegram 账户。",
    "Открыть Web app": "打开 Web app",
    "Открыть Telegram": "打开 Telegram",
    "языков интерфейса": "种界面语言",
    "маршруты уровня": "级别路径",
    "старт без оплаты": "免费开始",
    "Что вы делаете сегодня": "今天要做什么",
    "Короткий учебный цикл вместо бесконечной ленты упражнений": "用短学习循环代替无尽练习流",
    "Каждый день начинается с конкретной задачи и заканчивается повторением того, что реально просело.": "每天从一个具体任务开始，并以复习真正薄弱的地方结束。",
    "Исправление сохранено в Mistakes и Notes.": "修正已保存到 Mistakes 和 Notes。",
    "Четыре рабочих режима в одном учебном профиле": "一个学习账户中的四种工作模式",
    "Каждый модуль выглядит как продуктовый инструмент, а не как абстрактная карточка возможностей.": "每个模块都像真实产品工具，而不是抽象功能卡。",
    "Ошибки становятся тренировочным материалом": "错误会变成训练材料",
    "NERIVA связывает исправления, weak words, spelling, notes, offline decks, XP, streak, awards и daily bonus в понятный повтор.": "NERIVA 将修正、weak words、spelling、notes、offline decks、XP、streak、awards 和 daily bonus 连接成清晰复习循环。",
    "Две равные точки входа, один профиль обучения": "两个平等入口，一个学习账户",
    "Web app удобнее для длинных занятий и dashboard, Telegram — для быстрой практики, напоминаний, голоса и продолжения на ходу.": "Web app 更适合长课和 dashboard；Telegram 更适合快速练习、提醒、语音和随时继续。",
    "Фокусные уроки, dashboard прогресса, длинные сессии, тарифы и полный словарь ошибок.": "专注课程、进度 dashboard、长时学习、价格和完整错误词典。",
    "быстрая практика, voice, reminders и продолжение того же урока без потери прогресса.": "快速练习、voice、reminders，以及不丢进度地继续同一节课。",
    "Free для старта, Premium и Platinum для ежедневной практики": "Free 用于开始，Premium 和 Platinum 用于每日练习",
    "Оплата: Stars, YooKassa, TON и USDT.": "支付方式：Stars、YooKassa、TON 和 USDT。",
    "Выбрать тариф": "选择套餐",
    "Короткие истории пользователей": "简短用户故事",
    "Без рейтингов и обещаний результата: только реальные сценарии, где нужен язык.": "没有评分和结果承诺：只有真正需要语言的场景。",
    "Начать сегодня": "今天开始",
    "Откройте NERIVA и начните занятие с разбором ошибок": "打开 NERIVA，开始带错误解析的课程",
  },
};

function normalizeBoldText(value: string) {
  return value.replace(/\s+/g, " ").trim();
}

function resolveBoldLandingLocale(): BoldLandingLocale {
  const raw =
    window.poliglotSiteI18n?.currentLanguage?.()
    || new URLSearchParams(window.location.search).get("lang")
    || window.localStorage.getItem("poliglot-site-language")
    || "ru";
  if (raw === "ru") return "ru";
  if (raw === "zh") return "zh";
  return "en";
}

function useBoldLandingLocale() {
  const [locale, setLocale] = useState<BoldLandingLocale>(() => resolveBoldLandingLocale());

  useEffect(() => {
    const sync = () => setLocale(resolveBoldLandingLocale());
    window.addEventListener("poliglot-language-change", sync);
    return () => window.removeEventListener("poliglot-language-change", sync);
  }, []);

  return locale;
}

function translateBoldText(source: string, locale: BoldLandingLocale) {
  if (locale === "ru") return source;
  return boldTranslations[locale][source] || boldTranslations.en[source] || source;
}

function localizeBoldLanding(locale: BoldLandingLocale) {
  const root = document.querySelector(".bold-landing");
  if (!root) return;
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT);
  const nodes: Text[] = [];
  while (walker.nextNode()) nodes.push(walker.currentNode as Text);

  nodes.forEach((node) => {
    const raw = node.nodeValue || "";
    const normalized = normalizeBoldText(raw);
    if (!normalized) return;
    if (!boldNodeSources.has(node)) {
      boldNodeSources.set(node, normalized);
    }
    const source = boldNodeSources.get(node) || normalized;
    const translated = translateBoldText(source, locale);
    const leading = raw.match(/^\s*/)?.[0] || "";
    const trailing = raw.match(/\s*$/)?.[0] || "";
    node.nodeValue = `${leading}${translated}${trailing}`;
  });
}

export function BoldProductLanding() {
  const [activeTab, setActiveTab] = useState(0);
  const locale = useBoldLandingLocale();
  const activeScreen = heroScreens[activeTab] ?? heroScreens[0];
  const ActiveIcon = activeScreen.icon;

  useEffect(() => {
    localizeBoldLanding(locale);
    const timer = window.setTimeout(() => localizeBoldLanding(locale), 80);
    return () => window.clearTimeout(timer);
  }, [locale, activeTab]);

  return (
    <main className="bold-landing">
      <section className="bold-hero">
        <div className="bold-hero__matter">
          <GenerativeArtScene animate color="#7bdcff" particleColor="#f5d27a" />
        </div>
        <div className="bold-hero__veil" aria-hidden="true" />

        <div className="bold-hero__inner">
          <motion.div className="bold-hero__copy" initial={{ opacity: 0, y: 24 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.62 }}>
            <span className="bold-kicker">
              <Sparkles size={16} />
              AI Tutor для ежедневной практики
            </span>
            <h1>
              Говорите с <span>AI-репетитором</span>, который помнит ваши ошибки
            </h1>
            <p>
              Урок, диалог, голос, фото, заметки и повторение слабых мест собираются в один профиль Web app и Telegram.
            </p>
            <div className="entry-cta-row" aria-label="Точки входа NERIVA">
              <a className="entry-cta" data-entry="web-app" href={WEB_APP_HREF}>
                <Laptop size={20} />
                Открыть Web app
                <ArrowRight size={18} />
              </a>
              <a className="entry-cta" data-entry="telegram" href={TELEGRAM_HREF}>
                <MessageCircle size={20} />
                Открыть Telegram
                <ArrowRight size={18} />
              </a>
            </div>
            <div className="bold-hero__proof" aria-label="Ключевые факты">
              <span>
                <strong>35</strong> языков интерфейса
              </span>
              <span>
                <strong>A1-C2</strong> маршруты уровня
              </span>
              <span>
                <strong>Free</strong> старт без оплаты
              </span>
            </div>
          </motion.div>

          <motion.div className="hero-product-mockup" initial={{ opacity: 0, x: 36 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.66, delay: 0.08 }}>
            <div className="hero-product-top">
              <span>NERIVA</span>
              <strong>Today · 84/100</strong>
            </div>
            <div className="hero-product-tabs" role="tablist" aria-label="Демо функций">
              {heroScreens.map((screen, index) => (
                <button
                  key={screen.tab}
                  type="button"
                  role="tab"
                  aria-selected={activeTab === index}
                  className={activeTab === index ? "is-active" : undefined}
                  onClick={() => setActiveTab(index)}
                >
                  {screen.tab}
                </button>
              ))}
            </div>
            <motion.div key={activeScreen.tab} className="hero-product-output" initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.2 }}>
              <div className="hero-product-output__head">
                <ActiveIcon size={22} />
                <div>
                  <span>{activeScreen.label}</span>
                  <strong>{activeScreen.prompt}</strong>
                </div>
              </div>
              <p>{activeScreen.output}</p>
              {activeScreen.tab === "Голос" ? (
                <div className="voice-bars" aria-hidden="true">
                  {Array.from({ length: 18 }, (_, index) => (
                    <i key={index} style={{ animationDelay: `${index * 0.04}s` }} />
                  ))}
                </div>
              ) : (
                <div className="hero-product-tags" aria-hidden="true">
                  <span>{activeScreen.detail}</span>
                  <span>next review</span>
                </div>
              )}
            </motion.div>
          </motion.div>
        </div>
      </section>

      <section className="bold-section daily-loop-section">
        <div className="section-copy">
          <span className="bold-kicker">Что вы делаете сегодня</span>
          <h2>Короткий учебный цикл вместо бесконечной ленты упражнений</h2>
          <p>Каждый день начинается с конкретной задачи и заканчивается повторением того, что реально просело.</p>
        </div>
        <div className="daily-layout">
          <div className="daily-route-mockup" aria-label="Ежедневный маршрут">
            <div className="daily-route-mockup__top">
              <strong>Daily route</strong>
              <span>Travel · speaking</span>
            </div>
            <div className="route-line">
              <span>task</span>
              <span>answer</span>
              <span>correction</span>
              <span>repeat</span>
            </div>
            <p>Could you help me check in?</p>
            <small>Исправление сохранено в Mistakes и Notes.</small>
          </div>
          <div className="daily-steps">
            {dailySteps.map((step) => (
              <article className="daily-step" key={step.title}>
                <strong>{step.marker}</strong>
                <div>
                  <h3>{step.title}</h3>
                  <p>{step.body}</p>
                </div>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section id="features" className="bold-section modules-section">
        <div className="section-copy">
          <span className="bold-kicker">AI Tutor modules</span>
          <h2>Четыре рабочих режима в одном учебном профиле</h2>
          <p>Каждый модуль выглядит как продуктовый инструмент, а не как абстрактная карточка возможностей.</p>
        </div>
        <div className="module-grid">
          {modules.map((module) => {
            const Icon = module.icon;
            return (
              <article className="module-card" key={module.title}>
                <div className="module-card__title">
                  <Icon size={24} />
                  <h3>{module.title}</h3>
                </div>
                <p>{module.body}</p>
                <div className="module-card__screen">
                  {module.visual.map((item) => (
                    <span key={item}>{item}</span>
                  ))}
                  <strong>{module.sample}</strong>
                </div>
              </article>
            );
          })}
        </div>
      </section>

      <section className="memory-loop-section">
        <div className="memory-loop-section__copy">
          <span className="bold-kicker">Mistake memory loop</span>
          <h2>Ошибки становятся тренировочным материалом</h2>
          <p>NERIVA связывает исправления, weak words, spelling, notes, offline decks, XP, streak, awards и daily bonus в понятный повтор.</p>
        </div>
        <div className="memory-board" aria-label="Цикл памяти ошибок">
          {memoryNodes.map(([title, body]) => (
            <article className="memory-node" key={title}>
              <Repeat2 size={18} />
              <strong>{title}</strong>
              <span>{body}</span>
            </article>
          ))}
          <div className="memory-core">
            <Trophy size={32} />
            <span>XP · streak · awards</span>
          </div>
        </div>
      </section>

      <section className="bold-section entry-section">
        <div className="section-copy">
          <span className="bold-kicker">Web app + Telegram</span>
          <h2>Две равные точки входа, один профиль обучения</h2>
          <p>Web app удобнее для длинных занятий и dashboard, Telegram — для быстрой практики, напоминаний, голоса и продолжения на ходу.</p>
        </div>
        <div className="entry-grid">
          <article className="entry-panel">
            <Laptop size={28} />
            <h3>Web app</h3>
            <p>Фокусные уроки, dashboard прогресса, длинные сессии, тарифы и полный словарь ошибок.</p>
            <a className="entry-cta" data-entry="web-app" href={WEB_APP_HREF}>
              Открыть Web app <ChevronRight size={18} />
            </a>
          </article>
          <article className="entry-panel">
            <MessageCircle size={28} />
            <h3>Telegram</h3>
            <p>быстрая практика, voice, reminders и продолжение того же урока без потери прогресса.</p>
            <a className="entry-cta" data-entry="telegram" href={TELEGRAM_HREF}>
              Открыть Telegram <ChevronRight size={18} />
            </a>
          </article>
        </div>
      </section>

      <section id="pricing" className="pricing-section">
        <div className="section-copy">
          <span className="bold-kicker">Pricing</span>
          <h2>Free для старта, Premium и Platinum для ежедневной практики</h2>
          <p className="payment-methods">Оплата: Stars, YooKassa, TON и USDT.</p>
        </div>
        <div className="pricing-grid">
          {plans.map((plan) => (
            <article className={plan.featured ? "plan-card is-featured" : "plan-card"} key={plan.name}>
              {plan.featured ? <Crown size={22} /> : <ShieldCheck size={22} />}
              <h3>{plan.name}</h3>
              <strong className="plan-price">{plan.price}</strong>
              <p>{plan.body}</p>
              <ul>
                {plan.limits.map((limit) => (
                  <li key={limit}>
                    <CheckCircle size={16} />
                    {limit}
                  </li>
                ))}
              </ul>
              <a href={WEB_APP_HREF}>Выбрать тариф</a>
            </article>
          ))}
        </div>
      </section>

      <section id="reviews" className="bold-section reviews-section">
        <div className="section-copy">
          <span className="bold-kicker">Reviews</span>
          <h2>Короткие истории пользователей</h2>
          <p>Без рейтингов и обещаний результата: только реальные сценарии, где нужен язык.</p>
        </div>
        <div className="reviews-grid">
          {reviews.map((review) => (
            <article className="review-card" key={review.name}>
              <div>
                <strong>{review.name}</strong>
                <small>{review.detail}</small>
              </div>
              <p>{review.text}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="final-cta-section">
        <div>
          <span className="bold-kicker">Начать сегодня</span>
          <h2>Откройте NERIVA и начните занятие с разбором ошибок</h2>
        </div>
        <div className="entry-cta-row">
          <a className="entry-cta" data-entry="web-app" href={WEB_APP_HREF}>
            Web app <ArrowRight size={18} />
          </a>
          <a className="entry-cta" data-entry="telegram" href={TELEGRAM_HREF}>
            Telegram <ArrowRight size={18} />
          </a>
        </div>
      </section>
    </main>
  );
}
