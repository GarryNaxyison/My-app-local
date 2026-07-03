export const landingLocaleCodes = [
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

export type LandingLocale = (typeof landingLocaleCodes)[number];

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
    proof: readonly [string, string][];
  };
  scenario: {
    eyebrow: string;
    title: string;
    body: string;
    tags: readonly string[];
    captionLabel: string;
    caption: string;
  };
  features: {
    eyebrow: string;
    title: string;
    items: readonly { title: string; body: string }[];
  };
  notes: {
    eyebrow: string;
    title: string;
    body: string;
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
    }[];
  };
  telegram: {
    eyebrow: string;
    title: string;
    body: string;
    note: string;
    cta: string;
  };
  mobile: {
    eyebrow: string;
    title: string;
    body: string;
    rail: readonly string[];
  };
  community: {
    eyebrow: string;
    title: string;
    body: string;
    follow: string;
    note: string;
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
  };
  images: Record<string, string>;
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
    webApp: "Web app",
    telegram: "Telegram",
  },
  hero: {
    eyebrow: "AI tutor in web app and Telegram",
    title: "NERIVA helps you speak with more confidence",
    lead: "One short lesson gives you a phrase, your answer, a correction and the next repeat. Use the web app for a full session or Telegram for a quick check.",
    proof: [
      ["35 languages", "Interface and landing"],
      ["A1-C2", "Levels for daily practice"],
      ["Free, Premium, Platinum", "Clear daily limits"],
    ],
  },
  scenario: {
    eyebrow: "Live lesson",
    title: "Practice, get corrected, repeat",
    body: "NERIVA turns a real situation into a small loop: understand the phrase, answer, see the fix, save the weak point.",
    tags: ["phrase", "answer", "correction", "review"],
    captionLabel: "Hotel check-in",
    caption: "A correction, explanation, audio sample and your answer stay in one lesson.",
  },
  features: {
    eyebrow: "Core tools",
    title: "Four tools in one profile",
    items: [
      { title: "AI Tutor", body: "Short guided lessons with a clear next step." },
      { title: "Voice", body: "Pronunciation score, weak words and shadowing." },
      { title: "Photo", body: "Translate a menu, sign or task and practice it." },
      { title: "Mistakes", body: "Saved errors return when it is time to repeat." },
    ],
  },
  notes: {
    eyebrow: "Notes / Phrasebook",
    title: "Useful phrases do not disappear",
    body: "Save phrases, translations and weak words. Return to them when they matter.",
  },
  pricing: {
    eyebrow: "Pricing",
    title: "Start free. Upgrade when you need more practice.",
    lead: "Payment works through Telegram Stars and YooKassa/SBP.",
    methods: ["Telegram Stars", "YooKassa/SBP"],
    choose: "Choose plan",
    plans: [
      {
        name: "Free",
        label: "Try it",
        price: "0 ₽",
        period: "starter access",
        body: "Basic text practice, notes and progress.",
        limits: ["5 lessons a day", "15 practice messages", "Basic phrasebook"],
      },
      {
        name: "Premium",
        label: "Daily practice",
        oldPrice: "1000 ₽",
        price: "300 ₽",
        period: "per month",
        body: "The main plan for AI Tutor, voice and photo tools.",
        limits: ["50 lessons a day", "200 practice messages", "20 voice checks"],
      },
      {
        name: "Platinum",
        label: "Intensive",
        oldPrice: "2000 ₽",
        price: "590 ₽",
        period: "per month",
        body: "Higher limits for travel, work and exam preparation.",
        limits: ["100 lessons a day", "500 practice messages", "60 voice checks"],
      },
    ],
  },
  telegram: {
    eyebrow: "Telegram",
    title: "Quick practice without opening a second product",
    body: "Start a task, send voice or photo, get a reminder and continue the same profile in the web app.",
    note: "Lesson, voice, photo, mistakes and Premium are one tap away.",
    cta: "Open Telegram",
  },
  mobile: {
    eyebrow: "Mobile web app",
    title: "Phone sessions stay short",
    body: "Dashboard, lesson, mistakes and voice practice are compact on mobile.",
    rail: ["Progress", "Lesson", "Voice"],
  },
  community: {
    eyebrow: "Community",
    title: "Giveaways and product updates",
    body: "Follow NERIVA for short lessons, release notes and Premium key giveaways.",
    follow: "Follow NERIVA",
    note: "Short lessons and product updates.",
  },
  faq: {
    eyebrow: "FAQ",
    title: "Before you start",
    items: [
      { question: "Can I start free?", answer: "Yes. Free is enough to try the lesson loop." },
      { question: "Do I need Telegram?", answer: "No. The web app works on its own. Telegram is the fast entry." },
      { question: "Does voice practice work on mobile?", answer: "Yes. You can practise pronunciation from the phone." },
      { question: "Are mistakes saved?", answer: "Yes. Weak phrases return for review." },
    ],
  },
  finalCta: {
    title: "Start the first lesson today",
    body: "Full session in the web app. Quick practice in Telegram.",
    visualLabel: "Quick start",
    visualBody: "Choose a phrase, answer, get a correction, repeat tomorrow.",
  },
  images: {
    dashboard: "NERIVA dashboard with progress, level, XP and streak",
    mobileHome: "NERIVA mobile dashboard with today's progress",
    aiTutor: "AI tutor lesson with a corrected hotel check-in answer",
    mistakes: "Mistakes list with corrected English phrases",
    voice: "Pronunciation score and weak words",
    photo: "Photo translation with OCR result",
    notes: "Phrasebook with saved useful phrases",
    premium: "Free Premium and Platinum plan limits",
    telegram: "NERIVA Telegram bot in light theme",
    mobileLesson: "Mobile lesson correction",
    mobileMistakes: "Mobile mistakes review",
    mobileVoice: "Mobile pronunciation practice",
  },
};

const russian: LandingContent = {
  ...english,
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
    dark: "Включить тёмную тему",
    webApp: "Веб-приложение",
    telegram: "Telegram",
  },
  hero: {
    eyebrow: "AI-репетитор в вебе и Telegram",
    title: "NERIVA помогает говорить увереннее",
    lead: "Короткий урок: фраза, ваш ответ, исправление и следующий повтор. Веб-приложение для полной сессии, Telegram для быстрой практики.",
    proof: [
      ["35 языков", "Интерфейс и лендинг"],
      ["A1-C2", "Уровни для ежедневной практики"],
      ["Free, Premium, Platinum", "Понятные дневные лимиты"],
    ],
  },
  scenario: {
    eyebrow: "Живой урок",
    title: "Ответили, увидели ошибку, повторили",
    body: "NERIVA берёт реальную ситуацию и ведёт по короткому циклу: понять фразу, ответить, увидеть правку, сохранить слабое место.",
    tags: ["фраза", "ответ", "правка", "повтор"],
    captionLabel: "Заселение в отель",
    caption: "Исправление, объяснение, аудио и ваш ответ остаются в одном уроке.",
  },
  features: {
    eyebrow: "Главные инструменты",
    title: "Четыре функции в одном профиле",
    items: [
      { title: "AI-репетитор", body: "Короткие уроки с понятным следующим шагом." },
      { title: "Голос", body: "Оценка произношения, слабые слова и shadowing." },
      { title: "Фото", body: "Перевод меню, вывески или задания с практикой." },
      { title: "Ошибки", body: "Сохранённые ошибки возвращаются на повтор." },
    ],
  },
  notes: {
    eyebrow: "Заметки / Phrasebook",
    title: "Полезные фразы не теряются",
    body: "Сохраняйте фразы, переводы и слабые слова. Возвращайтесь к ним, когда они нужны.",
  },
  pricing: {
    ...english.pricing,
    eyebrow: "Тарифы",
    title: "Начните бесплатно. Расширяйте лимиты, когда понадобится.",
    lead: "Оплата через Telegram Stars и YooKassa/SBP.",
    choose: "Выбрать тариф",
    plans: [
      { ...english.pricing.plans[0], label: "Попробовать", body: "Базовая текстовая практика, заметки и прогресс.", limits: ["5 уроков в день", "15 сообщений практики", "Базовый phrasebook"] },
      { ...english.pricing.plans[1], label: "Ежедневно", body: "Основной тариф для AI-репетитора, голоса и фото.", limits: ["50 уроков в день", "200 сообщений практики", "20 голосовых проверок"] },
      { ...english.pricing.plans[2], label: "Интенсив", body: "Больше лимитов для поездки, работы и экзамена.", limits: ["100 уроков в день", "500 сообщений практики", "60 голосовых проверок"] },
    ],
  },
  telegram: {
    eyebrow: "Telegram",
    title: "Быстрая практика без второго продукта",
    body: "Запустите задание, отправьте голос или фото, получите напоминание и продолжите тот же профиль в веб-приложении.",
    note: "Уроки, голос, фото, ошибки и Premium открываются в один тап.",
    cta: "Открыть Telegram",
  },
  mobile: {
    eyebrow: "Мобильная версия",
    title: "На телефоне всё короче",
    body: "Прогресс, урок, ошибки и голосовая практика собраны компактно.",
    rail: ["Прогресс", "Урок", "Голос"],
  },
  community: {
    eyebrow: "Сообщество",
    title: "Розыгрыши и новости продукта",
    body: "Подпишитесь на NERIVA: короткие уроки, обновления и розыгрыши Premium-ключей.",
    follow: "Следить за NERIVA",
    note: "Короткие уроки и обновления.",
  },
  faq: {
    eyebrow: "FAQ",
    title: "Перед стартом",
    items: [
      { question: "Можно начать бесплатно?", answer: "Да. Free достаточно, чтобы попробовать цикл урока." },
      { question: "Telegram обязателен?", answer: "Нет. Веб-приложение работает отдельно. Telegram нужен для быстрого входа." },
      { question: "Голос работает на телефоне?", answer: "Да. Произношение можно тренировать с телефона." },
      { question: "Ошибки сохраняются?", answer: "Да. Слабые фразы возвращаются на повтор." },
    ],
  },
  finalCta: {
    title: "Начните первый урок сегодня",
    body: "Полная сессия в веб-приложении. Быстрая практика в Telegram.",
    visualLabel: "Быстрый старт",
    visualBody: "Выберите фразу, ответьте, получите правку и повторите завтра.",
  },
  images: {
    dashboard: "Дашборд NERIVA с прогрессом, уровнем, XP и streak",
    mobileHome: "Мобильный дашборд NERIVA с прогрессом дня",
    aiTutor: "Урок AI-репетитора с исправлением ответа",
    mistakes: "Список ошибок и исправленных фраз",
    voice: "Оценка произношения и слабые слова",
    photo: "Фото-перевод с OCR",
    notes: "Phrasebook с сохранёнными фразами",
    premium: "Лимиты Free Premium и Platinum",
    telegram: "Telegram-бот NERIVA в светлой теме",
    mobileLesson: "Мобильный урок с исправлением",
    mobileMistakes: "Мобильный список ошибок",
    mobileVoice: "Мобильная голосовая практика",
  },
};

const localeOverrides: Partial<Record<LandingLocale, Partial<LandingContent>>> = {
  es: { hero: { ...english.hero, title: "NERIVA te ayuda a hablar con más confianza", lead: "Una lección breve: frase, respuesta, corrección y siguiente repaso. Web app para sesiones completas, Telegram para practicar rápido." }, finalCta: { ...english.finalCta, title: "Empieza tu primera lección hoy", body: "Sesión completa en web. Práctica rápida en Telegram.", visualLabel: "Inicio rápido", visualBody: "Elige una frase, responde, corrige y repite mañana." } },
  de: { hero: { ...english.hero, title: "NERIVA hilft dir, sicherer zu sprechen", lead: "Eine kurze Lektion: Satz, Antwort, Korrektur und nächste Wiederholung. Web app für volle Sitzungen, Telegram für schnelle Praxis." }, finalCta: { ...english.finalCta, title: "Starte heute deine erste Lektion", body: "Volle Sitzung in der Web app. Schnelle Übung in Telegram.", visualLabel: "Schnellstart", visualBody: "Satz wählen, antworten, korrigieren, morgen wiederholen." } },
  fr: { hero: { ...english.hero, title: "NERIVA vous aide à parler avec plus d'assurance", lead: "Une leçon courte : phrase, réponse, correction et prochaine révision. Web app pour une vraie session, Telegram pour un entraînement rapide." } },
  it: { hero: { ...english.hero, title: "NERIVA ti aiuta a parlare con più sicurezza", lead: "Una lezione breve: frase, risposta, correzione e prossimo ripasso. Web app per sessioni complete, Telegram per pratica rapida." } },
  zh: { hero: { ...english.hero, title: "NERIVA 让你开口更有信心", lead: "一节短课包含短语、你的回答、纠正和下一次复习。网页应用适合完整学习，Telegram 适合快速练习。" } },
  ja: { hero: { ...english.hero, title: "NERIVAで、もっと自信を持って話す", lead: "短いレッスンで、表現、回答、訂正、次の復習まで進みます。Web app はしっかり学習、Telegram はすばやい練習に。" } },
  ko: { hero: { ...english.hero, title: "NERIVA로 더 자신 있게 말하세요", lead: "짧은 수업 안에 표현, 답변, 교정, 다음 복습이 들어 있습니다. Web app은 전체 학습, Telegram은 빠른 연습용입니다." } },
  uk: { hero: { ...russian.hero, title: "NERIVA допомагає говорити впевненіше", lead: "Короткий урок: фраза, ваша відповідь, виправлення і наступне повторення. Web app для повної сесії, Telegram для швидкої практики." } },
  pl: { hero: { ...english.hero, title: "NERIVA pomaga mówić pewniej", lead: "Krótka lekcja: fraza, odpowiedź, poprawka i następna powtórka. Web app do pełnej sesji, Telegram do szybkiej praktyki." } },
  pt: { hero: { ...english.hero, title: "NERIVA ajuda você a falar com mais confiança", lead: "Uma lição curta: frase, resposta, correção e próxima revisão. Web app para sessão completa, Telegram para prática rápida." } },
  tr: { hero: { ...english.hero, title: "NERIVA daha güvenli konuşmana yardım eder", lead: "Kısa ders: ifade, cevabın, düzeltme ve sonraki tekrar. Web app tam oturum, Telegram hızlı pratik içindir." } },
  vi: { hero: { ...english.hero, title: "NERIVA giúp bạn nói tự tin hơn", lead: "Một bài ngắn: cụm từ, câu trả lời, sửa lỗi và lần ôn tiếp theo. Web app cho buổi học đầy đủ, Telegram để luyện nhanh." } },
};

const spanishFallback: LandingContent = mergeContent(english, {
  nav: {
    features: "Funciones",
    pricing: "Planes",
    faq: "FAQ",
    privacy: "Privacidad",
    terms: "Condiciones",
    agreement: "Acuerdo de usuario",
    consent: "Consentimiento de datos",
    menu: "Menú",
    light: "Tema claro",
    dark: "Tema oscuro",
    webApp: "Web app",
    telegram: "Telegram",
  },
  hero: {
    eyebrow: "Tutor AI en web app y Telegram",
    title: "NERIVA te ayuda a hablar con más confianza",
    lead: "Una lección breve: frase, respuesta, corrección y siguiente repaso. Web app para sesiones completas, Telegram para practicar rápido.",
    proof: [
      ["35 idiomas", "Interfaz y landing"],
      ["A1-C2", "Niveles de práctica"],
      ["Free, Premium, Platinum", "Límites diarios claros"],
    ],
  },
  scenario: {
    eyebrow: "Lección real",
    title: "Practica, corrige y repite",
    body: "NERIVA convierte una situación real en un ciclo corto: entiende la frase, responde, mira la corrección y guarda el punto débil.",
    tags: ["frase", "respuesta", "corrección", "repaso"],
    captionLabel: "Check-in en hotel",
    caption: "Corrección, explicación, audio y respuesta quedan en una sola lección.",
  },
  features: {
    eyebrow: "Herramientas",
    title: "Cuatro funciones en un perfil",
    items: [
      { title: "Tutor AI", body: "Lecciones cortas con el siguiente paso claro." },
      { title: "Voz", body: "Puntuación, palabras débiles y shadowing." },
      { title: "Foto", body: "Traduce menú, señal o tarea y practícalo." },
      { title: "Errores", body: "Los errores guardados vuelven para repasar." },
    ],
  },
  notes: {
    eyebrow: "Notas / Phrasebook",
    title: "Las frases útiles no se pierden",
    body: "Guarda frases, traducciones y palabras débiles. Vuelve cuando las necesites.",
  },
  pricing: {
    ...english.pricing,
    eyebrow: "Planes",
    title: "Empieza gratis. Amplía límites cuando haga falta.",
    lead: "Pago por Telegram Stars y YooKassa/SBP.",
    choose: "Elegir plan",
    plans: [
      { ...english.pricing.plans[0], label: "Probar", body: "Práctica básica, notas y progreso.", limits: ["5 lecciones al día", "15 mensajes de práctica", "Phrasebook básico"] },
      { ...english.pricing.plans[1], label: "Diario", body: "El plan principal para tutor AI, voz y foto.", limits: ["50 lecciones al día", "200 mensajes de práctica", "20 revisiones de voz"] },
      { ...english.pricing.plans[2], label: "Intensivo", body: "Más límites para viaje, trabajo y examen.", limits: ["100 lecciones al día", "500 mensajes de práctica", "60 revisiones de voz"] },
    ],
  },
  telegram: {
    eyebrow: "Telegram",
    title: "Práctica rápida sin otro producto",
    body: "Inicia una tarea, envía voz o foto, recibe recordatorio y sigue el mismo perfil en la web app.",
    note: "Lección, voz, foto, errores y Premium están a un toque.",
    cta: "Abrir Telegram",
  },
  mobile: {
    eyebrow: "Mobile web app",
    title: "En el teléfono todo queda corto",
    body: "Progreso, lección, errores y voz se muestran de forma compacta.",
    rail: ["Progreso", "Lección", "Voz"],
  },
  community: {
    eyebrow: "Comunidad",
    title: "Sorteos y novedades",
    body: "Sigue NERIVA para lecciones breves, lanzamientos y sorteos de Premium.",
    follow: "Sigue a NERIVA",
    note: "Lecciones breves y novedades.",
  },
  faq: {
    eyebrow: "FAQ",
    title: "Antes de empezar",
    items: [
      { question: "¿Puedo empezar gratis?", answer: "Sí. Free basta para probar el ciclo." },
      { question: "¿Telegram es obligatorio?", answer: "No. La web app funciona sola. Telegram es entrada rápida." },
      { question: "¿La voz funciona en móvil?", answer: "Sí. Puedes practicar pronunciación desde el teléfono." },
      { question: "¿Se guardan los errores?", answer: "Sí. Las frases débiles vuelven al repaso." },
    ],
  },
  finalCta: {
    title: "Empieza tu primera lección hoy",
    body: "Sesión completa en web. Práctica rápida en Telegram.",
    visualLabel: "Inicio rápido",
    visualBody: "Elige una frase, responde, corrige y repite mañana.",
  },
});

const genericTitles: Partial<Record<LandingLocale, string>> = {
  tg: "NERIVA ба шумо кӯмак мекунад, ки боэътимодтар гап занед",
  uz: "NERIVA ishonchliroq gapirishga yordam beradi",
  tt: "NERIVA ышанычлырак сөйләшергә ярдәм итә",
  hy: "NERIVA-ն օգնում է խոսել ավելի վստահ",
  kk: "NERIVA сенімді сөйлеуге көмектеседі",
  ky: "NERIVA ишенимдүү сүйлөөгө жардам берет",
  ka: "NERIVA გეხმარებათ უფრო თავდაჯერებულად საუბარში",
  ro: "NERIVA te ajută să vorbești mai sigur",
  ar: "NERIVA يساعدك على التحدث بثقة أكبر",
  bn: "NERIVA আপনাকে আরও আত্মবিশ্বাসের সঙ্গে বলতে সাহায্য করে",
  cs: "NERIVA pomáhá mluvit jistěji",
  el: "Το NERIVA σας βοηθά να μιλάτε με περισσότερη σιγουριά",
  hi: "NERIVA आपको अधिक आत्मविश्वास से बोलने में मदद करता है",
  hu: "A NERIVA segít magabiztosabban beszélni",
  id: "NERIVA membantu Anda berbicara lebih percaya diri",
  nl: "NERIVA helpt je met meer vertrouwen spreken",
  sv: "NERIVA hjälper dig att tala tryggare",
  ta: "NERIVA நீங்கள் நம்பிக்கையுடன் பேச உதவுகிறது",
  te: "NERIVA మీరు మరింత నమ్మకంగా మాట్లాడటానికి సహాయపడుతుంది",
  th: "NERIVA ช่วยให้คุณพูดได้มั่นใจขึ้น",
  tl: "Tinutulungan ka ng NERIVA na magsalita nang mas may kumpiyansa",
};

function mergeContent(base: LandingContent, patch: Partial<LandingContent> = {}): LandingContent {
  return {
    ...base,
    ...patch,
    nav: { ...base.nav, ...patch.nav },
    hero: { ...base.hero, ...patch.hero },
    scenario: { ...base.scenario, ...patch.scenario },
    features: { ...base.features, ...patch.features },
    notes: { ...base.notes, ...patch.notes },
    pricing: { ...base.pricing, ...patch.pricing },
    telegram: { ...base.telegram, ...patch.telegram },
    mobile: { ...base.mobile, ...patch.mobile },
    community: { ...base.community, ...patch.community },
    faq: { ...base.faq, ...patch.faq },
    finalCta: { ...base.finalCta, ...patch.finalCta },
    images: { ...base.images, ...patch.images },
  };
}

export function normalizeLandingLocale(locale: string | null | undefined): LandingLocale {
  const normalized = (locale || "ru").toLowerCase();
  return landingLocaleCodes.includes(normalized as LandingLocale) ? (normalized as LandingLocale) : "ru";
}

export function getLandingContent(locale: string | null | undefined): LandingContent {
  const code = normalizeLandingLocale(locale);
  if (code === "ru") return russian;
  if (code === "en") return english;

  const base = genericTitles[code]
    ? mergeContent(spanishFallback, { hero: { ...spanishFallback.hero, title: genericTitles[code] } })
    : spanishFallback;

  return mergeContent(base, localeOverrides[code]);
}
