export const landingLocaleCodes = ["ru", "en"] as const;

export type LandingLocale = (typeof landingLocaleCodes)[number];

export type LandingFeature = {
  title: string;
  short: string;
  body: string;
  stat: string;
  imageKey: keyof LandingContent["images"];
};

export type LandingMobileExample = {
  title: string;
  short: string;
  body: string;
  stat: string;
  imageKey: keyof LandingContent["images"];
};

export type LandingContent = {
  nav: {
    features: string;
    pricing: string;
    faq: string;
    privacy: string;
    terms: string;
    agreement: string;
    consent: string;
    menu: string;
    light: string;
    dark: string;
    webApp: string;
    telegram: string;
  };
  hero: {
    eyebrow: string;
    title: string;
    lead: string;
    announcement: string;
    proof: readonly [string, string][];
    badges: readonly string[];
  };
  metrics: readonly {
    label: string;
    value: string;
    body: string;
  }[];
  outcomes: {
    eyebrow: string;
    title: string;
    body: string;
    items: readonly {
      title: string;
      body: string;
    }[];
  };
  beforeAfter: {
    eyebrow: string;
    title: string;
    beforeTitle: string;
    before: readonly string[];
    afterTitle: string;
    after: readonly string[];
  };
  scenario: {
    eyebrow: string;
    title: string;
    body: string;
    tags: readonly string[];
    captionLabel: string;
    caption: string;
  };
  useCases: {
    eyebrow: string;
    title: string;
    body: string;
    items: readonly {
      title: string;
      body: string;
      cta: string;
      imageKey: keyof LandingContent["images"];
    }[];
  };
  proof: {
    eyebrow: string;
    title: string;
    body: string;
    steps: readonly {
      title: string;
      body: string;
    }[];
  };
  features: {
    eyebrow: string;
    title: string;
    items: readonly LandingFeature[];
  };
  pricing: {
    eyebrow: string;
    title: string;
    lead: string;
    methods: readonly string[];
    choose: string;
    plans: readonly {
      name: string;
      label: string;
      oldPrice?: string;
      price: string;
      period: string;
      body: string;
      limits: readonly string[];
      access: readonly string[];
    }[];
  };
  telegram: {
    eyebrow: string;
    title: string;
    body: string;
    note: string;
    cta: string;
    chips: readonly string[];
  };
  mobile: {
    eyebrow: string;
    title: string;
    body: string;
    items: readonly LandingMobileExample[];
  };
  ecosystem: {
    eyebrow: string;
    title: string;
    body: string;
    items: readonly {
      title: string;
      body: string;
      imageKey: keyof LandingContent["images"];
    }[];
  };
  methodology: {
    eyebrow: string;
    title: string;
    body: string;
    steps: readonly {
      label: string;
      title: string;
      body: string;
    }[];
  };
  community: {
    eyebrow: string;
    title: string;
    body: string;
    follow: string;
    note: string;
    chips: readonly string[];
  };
  faq: {
    eyebrow: string;
    title: string;
    items: readonly { question: string; answer: string }[];
  };
  finalCta: {
    title: string;
    body: string;
    visualLabel: string;
    visualBody: string;
    flow: readonly string[];
  };
  images: {
    dashboard: string;
    mobileHome: string;
    aiTutor: string;
    mistakes: string;
    voice: string;
    photo: string;
    notes: string;
    telegram: string;
    mobileLesson: string;
    mobileMistakes: string;
    mobileVoice: string;
  };
};

const english: LandingContent = {
  nav: {
    features: "Features",
    pricing: "Pricing",
    faq: "FAQ",
    privacy: "Privacy",
    terms: "Terms",
    agreement: "User Agreement",
    consent: "Personal data consent",
    menu: "Menu",
    light: "Use light theme",
    dark: "Use dark theme",
    webApp: "Browser",
    telegram: "Telegram",
  },
  hero: {
    eyebrow: "AI tutor in browser and Telegram",
    title: "NERIVA turns every answer into a real lesson",
    lead: "A strict AI tutor gives a short story, useful chunks, questions, your answer, a correction and a review date. Browser and Telegram keep one learning profile.",
    announcement: "New lesson in 3 minutes",
    proof: [
      ["Browser + Telegram", "One learning profile"],
      ["Story -> answer", "Practice starts from context"],
      ["Correction -> repeat", "Mistakes return on time"],
    ],
    badges: ["CEFR level", "6 chunks", "review date"],
  },
  metrics: [
    { label: "Voice check", value: "score", body: "Pronunciation is checked on phrases you actually practise" },
    { label: "Route depth", value: "A1-C2", body: "Level-aware lessons, voice and review" },
    { label: "Profile", value: "1", body: "Browser and Telegram share the same learning memory" },
  ],
  outcomes: {
    eyebrow: "Why it helps",
    title: "If you study but still avoid speaking, NERIVA gives the next small step",
    body: "The product is built for the common problem: you know words, but in a real moment you cannot quickly put them into a phrase. NERIVA gives a situation, asks for an answer, corrects it and keeps the weak place for review.",
    items: [
      {
        title: "No blank page",
        body: "You start from a concrete travel, work, exam or speaking situation, not from an empty chat.",
      },
      {
        title: "No lost correction",
        body: "A corrected phrase goes to mistakes and notes, so one useful answer does not disappear.",
      },
      {
        title: "No separate tools",
        body: "Text, voice, photo and Telegram feed the same profile instead of splitting your practice.",
      },
    ],
  },
  beforeAfter: {
    eyebrow: "Before / After",
    title: "More than a vocabulary app: NERIVA closes the loop",
    beforeTitle: "Before",
    before: ["words saved without use", "chat replies forgotten", "no review schedule"],
    afterTitle: "After",
    after: ["context lesson", "teacher-style correction", "review on time"],
  },
  scenario: {
    eyebrow: "Live lesson",
    title: "One micro-lesson has structure, not random chat",
    body: "Example: hotel check-in. NERIVA builds a level-matched story, pulls useful chunks, asks questions, checks your answer and turns weak phrases into review.",
    tags: ["Story", "Chunks", "Answer", "Review"],
    captionLabel: "Hotel check-in",
    caption: "The task, answer, correction and review options stay in one lesson.",
  },
  useCases: {
    eyebrow: "Use cases",
    title: "Why people open NERIVA",
    body: "Each path starts from a real reason to speak or understand a language today.",
    items: [
      {
        title: "Travel",
        body: "Practise check-in, transport, cafe orders and short questions before you need them.",
        cta: "Start travel lesson",
        imageKey: "aiTutor",
      },
      {
        title: "Work",
        body: "Rehearse self-intros, calls, messages and deadline questions without waiting for a teacher.",
        cta: "Practise work phrase",
        imageKey: "notes",
      },
      {
        title: "Exam",
        body: "Turn weak grammar and speaking answers into repeatable tasks with a clear review date.",
        cta: "Train exam answer",
        imageKey: "mistakes",
      },
      {
        title: "Pronunciation",
        body: "Record a phrase, see weak words and repeat while the context is still fresh.",
        cta: "Check voice",
        imageKey: "voice",
      },
    ],
  },
  proof: {
    eyebrow: "Product proof",
    title: "How it looks inside the product",
    body: "The screenshots show the actual loop: open a lesson, answer, get a correction, then come back to weak phrases.",
    steps: [
      { title: "Pick a task", body: "Choose the goal and level, then start with a concrete situation." },
      { title: "Answer in your words", body: "Type, speak or use a photo when the context comes from real life." },
      { title: "Get the correction", body: "NERIVA explains what to change and shows a better phrase." },
      { title: "Repeat later", body: "Mistakes, notes and reminders bring the phrase back when it matters." },
    ],
  },
  features: {
    eyebrow: "Four functions",
    title: "Four tools that feed the same learning memory",
    items: [
      {
        title: "AI Tutor",
        short: "lesson",
        body: "A reusable micro-lesson: story, six useful chunks, questions, production task and review options matched to your level.",
        stat: "structured lesson",
        imageKey: "aiTutor",
      },
      {
        title: "Voice",
        short: "pronunciation",
        body: "Score, weak words and shadowing for phrases from your real practice, not isolated pronunciation drills.",
        stat: "speech connected",
        imageKey: "voice",
      },
      {
        title: "Photo",
        short: "OCR",
        body: "OCR pulls text from a menu, sign or exercise, then turns it into a usable phrase for practice.",
        stat: "context captured",
        imageKey: "photo",
      },
      {
        title: "Mistakes",
        short: "review",
        body: "Wrong phrases become review tasks with a date, so mistakes do not disappear after one AI reply.",
        stat: "repeat planned",
        imageKey: "mistakes",
      },
    ],
  },
  pricing: {
    eyebrow: "Pricing",
    title: "Start free. Add limits when practice becomes regular.",
    lead: "Payment works through Telegram Stars and YooKassa/SBP.",
    methods: ["Telegram Stars", "YooKassa/SBP"],
    choose: "Choose plan",
    plans: [
      {
        name: "Free",
        label: "Try it",
        price: "0 ₽",
        period: "starter access",
        body: "For checking the lesson loop without paying.",
        limits: ["5 lessons a day", "15 practice messages", "Basic phrasebook"],
        access: ["text practice", "notes", "progress"],
      },
      {
        name: "Premium",
        label: "Daily practice",
        oldPrice: "1000 ₽",
        price: "300 ₽",
        period: "per month",
        body: "The main plan for AI Tutor, voice, photo and saved mistakes.",
        limits: ["50 lessons a day", "200 practice messages", "20 voice checks"],
        access: ["AI Tutor", "Voice", "Photo", "Mistakes"],
      },
      {
        name: "Platinum",
        label: "Intensive",
        oldPrice: "2000 ₽",
        price: "590 ₽",
        period: "per month",
        body: "Higher limits for travel, work and exam preparation.",
        limits: ["100 lessons a day", "500 practice messages", "60 voice checks"],
        access: ["higher limits", "voice intensive", "photo practice"],
      },
    ],
  },
  telegram: {
    eyebrow: "Telegram companion",
    title: "Quick practice without a second product",
    body: "Start a task, send voice or photo, get a reminder and continue the same profile in the browser.",
    note: "Lesson, voice, photo, mistakes and Premium open in one tap.",
    cta: "Open Telegram",
    chips: ["Voice", "Photo", "Reminder"],
  },
  mobile: {
    eyebrow: "Mobile browser",
    title: "A short check-in, not a long scroll",
    body: "Progress, lesson, mistakes and voice practice stay compact on the phone.",
    items: [
      {
        title: "Progress",
        short: "today",
        body: "Open the phone and see level, streak, XP and the next useful task without digging through the app.",
        stat: "today is clear",
        imageKey: "mobileHome",
      },
      {
        title: "Lesson",
        short: "correction",
        body: "Answer a travel or work phrase, get the correction and keep the same learning profile.",
        stat: "answer checked",
        imageKey: "mobileLesson",
      },
      {
        title: "Mistakes",
        short: "review",
        body: "Weak phrases return as short repeat tasks, so the phone session has a clear finish.",
        stat: "repeat planned",
        imageKey: "mobileMistakes",
      },
      {
        title: "Voice",
        short: "score",
        body: "Record a phrase, see weak words and repeat it while the context is still fresh.",
        stat: "weak words found",
        imageKey: "mobileVoice",
      },
    ],
  },
  ecosystem: {
    eyebrow: "Ecosystem",
    title: "Web, mobile and Telegram stay in one loop",
    body: "The Stitch direction uses a device triptych: the browser for full lessons, mobile web for quick review, and Telegram for fast capture.",
    items: [
      { title: "Web app", body: "Full dashboard, lesson correction, pricing and longer sessions.", imageKey: "dashboard" },
      { title: "Mobile web", body: "Short lesson, review and voice practice from the phone.", imageKey: "mobileHome" },
      { title: "Telegram", body: "Fast start, reminders, voice, photo and the same profile.", imageKey: "telegram" },
    ],
  },
  methodology: {
    eyebrow: "Methodology",
    title: "Discovery, lesson, review",
    body: "The Stitch screen adds a vertical operating model. In NERIVA it maps to how a weak phrase becomes a timed repeat.",
    steps: [
      { label: "Phase 01", title: "Discovery", body: "NERIVA reads the goal, level and context before choosing the task." },
      { label: "Phase 02", title: "Correction", body: "The answer is checked like a teacher would check a short speaking attempt." },
      { label: "Phase 03", title: "Review", body: "Weak phrases return through mistakes, notes, voice and Telegram reminders." },
    ],
  },
  community: {
    eyebrow: "Community",
    title: "NERIVA channel: lessons, updates and Premium giveaways",
    body: "Follow real product updates, short practice posts and Premium key giveaways.",
    follow: "Follow NERIVA",
    note: "Short lessons and release notes.",
    chips: ["mini lesson", "release note", "Premium giveaway"],
  },
  faq: {
    eyebrow: "FAQ",
    title: "Before you start",
    items: [
      { question: "Can I start free?", answer: "Yes. Free is enough to try the lesson loop." },
      { question: "Do I need Telegram?", answer: "No. The browser version works on its own. Telegram is the fast channel." },
      { question: "Does voice work on mobile?", answer: "Yes. You can practise pronunciation from the phone." },
      { question: "Are mistakes saved?", answer: "Yes. Weak phrases return for review." },
    ],
  },
  finalCta: {
    title: "Start the first lesson today",
    body: "Full session in the browser. Quick practice in Telegram.",
    visualLabel: "Lesson loop",
    visualBody: "Pick a phrase, answer, see the correction and repeat tomorrow.",
    flow: ["Phrase", "Answer", "Correction", "Repeat"],
  },
  images: {
    dashboard: "NERIVA dashboard with progress, level, XP and streak",
    mobileHome: "NERIVA mobile dashboard with today's progress",
    aiTutor: "AI tutor lesson with a corrected hotel check-in answer",
    mistakes: "Mistakes list with corrected English phrases",
    voice: "Pronunciation score and weak words",
    photo: "Photo translation with OCR result",
    notes: "Phrasebook with saved useful phrases",
    telegram: "NERIVA Telegram bot in light theme",
    mobileLesson: "Mobile lesson correction",
    mobileMistakes: "Mobile mistakes review",
    mobileVoice: "Mobile pronunciation practice",
  },
};

const russian: LandingContent = {
  nav: {
    features: "Функции",
    pricing: "Тарифы",
    faq: "FAQ",
    privacy: "Политика",
    terms: "Условия",
    agreement: "Пользовательское соглашение",
    consent: "Согласие на обработку персональных данных",
    menu: "Меню",
    light: "Включить светлую тему",
    dark: "Включить темную тему",
    webApp: "Браузер",
    telegram: "Telegram",
  },
  hero: {
    eyebrow: "AI-репетитор в браузере и Telegram",
    title: "NERIVA превращает каждый ответ в настоящий урок",
    lead: "NERIVA дает короткий урок под ваш уровень: ситуацию, полезные фразы, вопросы, проверку ответа и повтор. В браузере и Telegram остается один профиль.",
    announcement: "Новый урок за 3 минуты",
    proof: [
      ["Браузер + Telegram", "Один учебный профиль"],
      ["История -> ответ", "Практика начинается с контекста"],
      ["Правка -> повтор", "Ошибки возвращаются вовремя"],
    ],
    badges: ["уровень CEFR", "6 чанков", "дата повтора"],
  },
  metrics: [
    { label: "Проверка голоса", value: "score", body: "Произношение проверяется на фразах из вашей практики" },
    { label: "Глубина маршрута", value: "A1-C2", body: "Уроки, голос и повтор учитывают уровень" },
    { label: "Профиль", value: "1", body: "Браузер и Telegram используют одну учебную память" },
  ],
  outcomes: {
    eyebrow: "Зачем это нужно",
    title: "Если вы учите язык, но всё равно молчите, нужен не ещё один список слов",
    body: "Обычно проблема не в том, что вы совсем ничего не знаете. Проблема в том, что в нужный момент фраза не собирается. NERIVA даёт ситуацию, просит ответить, исправляет и не даёт ошибкам пропасть.",
    items: [
      {
        title: "Есть с чего начать",
        body: "Вы открываете не пустой чат, а короткую задачу: поездка, работа, экзамен или разговор.",
      },
      {
        title: "Исправление не теряется",
        body: "Полезная правка попадает в ошибки и заметки, чтобы к ней можно было вернуться.",
      },
      {
        title: "Всё в одном профиле",
        body: "Текст, голос, фото, браузер и Telegram работают на одну учебную память.",
      },
    ],
  },
  beforeAfter: {
    eyebrow: "До / После",
    title: "Больше чем словарик: NERIVA закрывает цикл обучения",
    beforeTitle: "До",
    before: ["слова лежат без дела", "ответы быстро забываются", "нет графика повтора"],
    afterTitle: "После",
    after: ["урок из контекста", "правка как от преподавателя", "повтор вовремя"],
  },
  scenario: {
    eyebrow: "Живой урок",
    title: "Короткий урок с понятным результатом",
    body: "Например, заселение в отель. NERIVA подбирает ситуацию под ваш уровень, дает нужные фразы, задает вопросы, проверяет ответ и возвращает слабые места в повтор.",
    tags: ["Ситуация", "Фразы", "Ответ", "Повтор"],
    captionLabel: "Заселение в отель",
    caption: "Задание, ответ, правка и варианты повтора остаются в одном уроке.",
  },
  useCases: {
    eyebrow: "Сценарии",
    title: "Для чего открывают NERIVA",
    body: "У каждого занятия есть понятная причина: скоро поездка, рабочий созвон, экзамен или желание звучать увереннее.",
    items: [
      {
        title: "Поездка",
        body: "Прогоните заселение, кафе, транспорт и короткие вопросы до момента, когда они понадобятся.",
        cta: "Начать урок для поездки",
        imageKey: "aiTutor",
      },
      {
        title: "Работа",
        body: "Отрепетируйте self-intro, созвон, письмо или вопрос по срокам без ожидания преподавателя.",
        cta: "Разобрать рабочую фразу",
        imageKey: "notes",
      },
      {
        title: "Экзамен",
        body: "Слабая грамматика и speaking-ответы превращаются в задания, которые можно повторить.",
        cta: "Потренировать ответ",
        imageKey: "mistakes",
      },
      {
        title: "Произношение",
        body: "Запишите фразу, увидьте слабые слова и повторите её, пока контекст ещё свежий.",
        cta: "Проверить голос",
        imageKey: "voice",
      },
    ],
  },
  proof: {
    eyebrow: "Доказательство продуктом",
    title: "Как это выглядит в продукте",
    body: "На скринах не абстрактные обещания, а рабочий цикл: открыть урок, ответить, получить правку и вернуться к слабой фразе.",
    steps: [
      { title: "Выбираете задачу", body: "Цель и уровень задают короткий урок под ситуацию." },
      { title: "Отвечаете своими словами", body: "Можно писать, говорить голосом или использовать фото из реальной жизни." },
      { title: "Получаете правку", body: "NERIVA показывает, что исправить, и даёт более естественную фразу." },
      { title: "Повторяете позже", body: "Ошибки, заметки и напоминания возвращают фразу в нужный момент." },
    ],
  },
  features: {
    eyebrow: "Четыре функции",
    title: "Четыре инструмента кормят одну учебную память",
    items: [
      {
        title: "AI-репетитор",
        short: "урок",
        body: "Короткий урок: ситуация, шесть полезных фраз, вопросы, практическое задание и повтор под ваш уровень.",
        stat: "структурный урок",
        imageKey: "aiTutor",
      },
      {
        title: "Голос",
        short: "произношение",
        body: "Оценка, слабые слова и shadowing для фраз из вашей практики, а не отдельная тренировка ради тренировки.",
        stat: "речь связана",
        imageKey: "voice",
      },
      {
        title: "Фото",
        short: "OCR",
        body: "OCR достает текст из меню, вывески или задания и превращает его в фразу для практики.",
        stat: "контекст пойман",
        imageKey: "photo",
      },
      {
        title: "Ошибки",
        short: "повтор",
        body: "Неверные фразы становятся заданиями с датой повтора, а не теряются после одной проверки.",
        stat: "повтор запланирован",
        imageKey: "mistakes",
      },
    ],
  },
  pricing: {
    ...english.pricing,
    eyebrow: "Тарифы",
    title: "Free для пробы. Premium для нормальной ежедневной практики.",
    lead: "Premium нужен, когда вы занимаетесь почти каждый день: открываются голос, фото, больше уроков и нормальный запас сообщений. Оплата через Telegram Stars и YooKassa/SBP.",
    choose: "Выбрать тариф",
    plans: [
      {
        name: "Free",
        label: "Попробовать",
        price: "0 ₽",
        period: "стартовый доступ",
        body: "Для первого знакомства: понять, как работает урок, правка и сохранение фраз.",
        limits: ["5 уроков в день", "15 сообщений практики", "Базовый phrasebook"],
        access: ["текстовая практика", "заметки", "прогресс"],
      },
      {
        name: "Premium",
        label: "Ежедневно",
        oldPrice: "1000 ₽",
        price: "300 ₽",
        period: "в месяц",
        body: "Основной тариф для ежедневной практики: AI-уроки, голос, фото и разбор ошибок в одном профиле.",
        limits: ["50 уроков в день", "200 сообщений практики", "20 голосовых проверок"],
        access: ["AI-репетитор", "голос, фото и разбор ошибок", "повтор слабых мест"],
      },
      {
        name: "Platinum",
        label: "Интенсив",
        oldPrice: "2000 ₽",
        price: "590 ₽",
        period: "в месяц",
        body: "Для плотной подготовки перед поездкой, рабочим периодом или экзаменом, когда лимиты Free и Premium быстро заканчиваются.",
        limits: ["100 уроков в день", "500 сообщений практики", "60 голосовых проверок"],
        access: ["больше лимитов", "голос интенсив", "фото практика"],
      },
    ],
  },
  telegram: {
    eyebrow: "Telegram-связка",
    title: "Быстрый вход через Telegram",
    body: "Откройте короткое задание, отправьте голос или фото и вернитесь к тому же профилю в браузере.",
    note: "Уроки, голос, фото, ошибки и Premium открываются в один тап.",
    cta: "Открыть Telegram",
    chips: ["Голос", "Фото", "Напоминание"],
  },
  mobile: {
    eyebrow: "Мобильный браузер",
    title: "Короткая практика с телефона",
    body: "Прогресс, урок, ошибки и голос собраны компактно.",
    items: [
      {
        title: "Прогресс",
        short: "сегодня",
        body: "Открываете телефон и сразу видите уровень, серию занятий, XP и следующее полезное задание.",
        stat: "день понятен",
        imageKey: "mobileHome",
      },
      {
        title: "Урок",
        short: "правка",
        body: "Отвечаете на фразу для поездки или работы, получаете правку и продолжаете тот же профиль.",
        stat: "ответ проверен",
        imageKey: "mobileLesson",
      },
      {
        title: "Ошибки",
        short: "повтор",
        body: "Слабые фразы возвращаются короткими заданиями, у мобильной сессии есть понятный финиш.",
        stat: "повтор запланирован",
        imageKey: "mobileMistakes",
      },
      {
        title: "Голос",
        short: "score",
        body: "Записываете фразу, видите слабые слова и повторяете, пока контекст ещё свежий.",
        stat: "слабые слова найдены",
        imageKey: "mobileVoice",
      },
    ],
  },
  ecosystem: {
    eyebrow: "Экосистема",
    title: "Браузер, телефон и Telegram работают как один продукт",
    body: "Главный сценарий остается в web app, а быстрые входы помогают продолжать учебу без потери прогресса.",
    items: [
      {
        title: "Большой экран",
        body: "Уроки, прогресс, тарифы, заметки и ошибки удобнее разбирать в браузере.",
        imageKey: "dashboard",
      },
      {
        title: "Мобильный режим",
        body: "Короткий урок, повтор и голосовая проверка помещаются в несколько минут с телефона.",
        imageKey: "mobileHome",
      },
      {
        title: "Telegram",
        body: "Быстрый старт, фото, голос и напоминания подключены к тому же учебному профилю.",
        imageKey: "telegram",
      },
    ],
  },
  methodology: {
    eyebrow: "Методика",
    title: "Диагностика, правка, повтор",
    body: "NERIVA показывает путь фразы: нашли слабое место, исправили, вернули в повтор в нужный день.",
    steps: [
      { label: "Фаза 01", title: "Диагностика", body: "NERIVA смотрит на цель, уровень и контекст перед выбором задания." },
      { label: "Фаза 02", title: "Правка", body: "Ответ проверяется так, как преподаватель проверяет короткую устную попытку." },
      { label: "Фаза 03", title: "Повтор", body: "Слабые фразы возвращаются через ошибки, заметки, голос и Telegram-напоминания." },
    ],
  },
  community: {
    eyebrow: "Сообщество",
    title: "Канал NERIVA: уроки, обновления и Premium-ключи",
    body: "Подпишитесь на обновления продукта, короткие задания и раздачи Premium-ключей.",
    follow: "Следить за NERIVA",
    note: "Короткие уроки и заметки о релизах.",
    chips: ["мини-урок", "релиз", "розыгрыш Premium"],
  },
  faq: {
    eyebrow: "FAQ",
    title: "Перед стартом",
    items: [
      { question: "Можно начать бесплатно?", answer: "Да. Free достаточно, чтобы попробовать цикл урока." },
      { question: "Telegram обязателен?", answer: "Нет. Браузерная версия работает отдельно. Telegram нужен для быстрого входа." },
      { question: "Голос работает на телефоне?", answer: "Да. Произношение можно тренировать с телефона." },
      { question: "Ошибки сохраняются?", answer: "Да. Слабые фразы возвращаются на повтор." },
    ],
  },
  finalCta: {
    title: "Начните первый урок сегодня",
    body: "Большой урок удобнее проходить в браузере. Короткое задание можно открыть в Telegram.",
    visualLabel: "Цикл урока",
    visualBody: "Выберите фразу, ответьте, получите правку и повторите завтра.",
    flow: ["Фраза", "Ответ", "Правка", "Повтор"],
  },
  images: {
    dashboard: "Дашборд NERIVA с прогрессом, уровнем, XP и серией занятий",
    mobileHome: "Мобильный дашборд NERIVA с прогрессом дня",
    aiTutor: "Урок AI-репетитора с исправлением ответа",
    mistakes: "Список ошибок и исправленных фраз",
    voice: "Оценка произношения и слабые слова",
    photo: "Фото-перевод с OCR",
    notes: "Phrasebook с сохраненными фразами",
    telegram: "Telegram-бот NERIVA в светлой теме",
    mobileLesson: "Мобильный урок с исправлением",
    mobileMistakes: "Мобильный список ошибок",
    mobileVoice: "Мобильная голосовая практика",
  },
};

export function normalizeLandingLocale(locale: string | null | undefined): LandingLocale {
  const normalized = (locale || "ru").toLowerCase();
  return normalized === "en" ? "en" : "ru";
}

export function getLandingContent(locale: string | null | undefined): LandingContent {
  return normalizeLandingLocale(locale) === "en" ? english : russian;
}
