#!/usr/bin/env node

import { writeFile } from "node:fs/promises";
import { join } from "node:path";

const outputPath = join(process.cwd(), "site-react", "public", "assets", "legal-documents-i18n.js");

const languages = [
  ["ru", "Русский"], ["en", "English"], ["es", "Español"], ["de", "Deutsch"],
  ["fr", "Français"], ["it", "Italiano"], ["zh", "中文"], ["ja", "日本語"], ["ko", "한국어"],
  ["tg", "Тоҷикӣ"], ["uz", "O‘zbekcha"], ["tt", "Татарча"], ["hy", "Հայերեն"],
  ["kk", "Қазақша"], ["ky", "Кыргызча"], ["ka", "ქართული"], ["uk", "Українська"],
  ["pl", "Polski"], ["ro", "Română"], ["pt", "Português"], ["ar", "العربية"],
  ["bn", "বাংলা"], ["cs", "Čeština"], ["el", "Ελληνικά"], ["hi", "हिंदी"],
  ["hu", "Magyar"], ["id", "Bahasa Indonesia"], ["nl", "Nederlands"], ["sv", "svenska"],
  ["ta", "தமிழ்"], ["te", "తెలుగు"], ["th", "ภาษาไทย"], ["tl", "Tagalog"],
  ["tr", "Türkçe"], ["vi", "Tiếng Việt"],
];

const operator = {
  title: "Оператор и реквизиты",
  name: "Самозанятый Чебан Денис Игоревич",
  inn: "ИНН 505017471160",
  address: "г. Щёлково, ул. Сиреневая, 9к1, кв. 9",
  email: "supportpoliglotai@gmail.com",
  bot: "@poliglot_ai_bot",
};

const nav = {
  home: "Poliglot AI",
  privacy: "Политика",
  terms: "Условия",
  agreement: "Пользовательское соглашение",
  consent: "Согласие на обработку персональных данных",
};

const cookie = {
  title: "Cookie и технические данные",
  body: "Мы используем необходимые cookie и локальное хранилище для входа, языка, темы, безопасности, сохранения согласий и корректной работы сайта. Необязательные cookie применяются только после согласия.",
  accept: "Принять все cookie",
  necessary: "Только необходимые",
  privacy: "Политика",
  consent: "Согласие",
};

const documents = {
  privacy: {
    title: "Политика обработки персональных данных",
    badge: "Защита данных",
    description: "Как Poliglot AI обрабатывает данные пользователей сайта, веб-приложения и Telegram-бота.",
    updated: "Обновлено 24 июня 2026",
    sections: [
      {
        id: "scope",
        title: "1. Общие положения",
        blocks: [
          "Политика действует для сервиса Poliglot AI: сайт poliglotai.ru, веб-приложение poliglotai.online/app и Telegram-бот @poliglot_ai_bot.",
          "Сервис является единой системой для изучения иностранных языков. Используя сервис, пользователь подтверждает, что ознакомился с настоящей Политикой, Пользовательским соглашением и Согласием на обработку персональных данных.",
        ],
      },
      {
        id: "operator",
        title: "2. Оператор персональных данных",
        blocks: [
          "Оператор персональных данных: Самозанятый Чебан Денис Игоревич, ИНН 505017471160, адрес: г. Щёлково, ул. Сиреневая, 9к1, кв. 9. Email для обращений: supportpoliglotai@gmail.com.",
        ],
      },
      {
        id: "data",
        title: "3. Какие данные обрабатываются",
        blocks: [
          "Оператор может обрабатывать email или логин web-аккаунта, Telegram ID, username, имя в Telegram, выбранные языки, уровень обучения, учебные ответы, сообщения, голосовые записи, изображения, результаты распознавания, переводы, прогресс, платежные статусы и обращения в поддержку.",
          "Также могут обрабатываться технические данные: IP-адрес, cookie, локальное хранилище браузера, данные устройства, браузера, операционной системы, времени доступа, источника перехода, событий безопасности и ошибок.",
          "Оператор не просит передавать паспортные данные, медицинские сведения, биометрию, специальные категории персональных данных и полные платежные реквизиты банковских карт.",
        ],
      },
      {
        id: "purposes",
        title: "4. Цели и правовые основания",
        blocks: [
          "Данные нужны для регистрации, авторизации, ведения единого профиля web и Telegram, предоставления уроков и практики, сохранения прогресса, настройки интерфейса, оплаты тарифов, поддержки пользователя, безопасности, предотвращения злоупотреблений, диагностики ошибок и улучшения сервиса.",
          "Обработка выполняется на основании согласия пользователя, исполнения пользовательского соглашения, выполнения требований закона и законного интереса оператора по обеспечению безопасности сервиса.",
        ],
      },
      {
        id: "cookies",
        title: "5. Cookie и локальное хранилище",
        blocks: [
          "Сайт и веб-приложение используют необходимые cookie и локальное хранилище для входа, сохранения языка, темы, согласий, сессии, защиты аккаунта и стабильной работы интерфейса.",
          "Необязательные аналитические или маркетинговые cookie применяются только после отдельного согласия, если такие технологии включены в текущей версии сайта.",
        ],
      },
      {
        id: "transfer",
        title: "6. Хранение, передача и инфраструктура",
        blocks: [
          "Первичный сбор и хранение персональных данных граждан Российской Федерации должны осуществляться с учетом требований российского законодательства о локализации персональных данных.",
          "Для технической работы сервиса могут использоваться хостинг, Telegram, платежные провайдеры, аналитика, AI-, OCR-, speech-to-text- и text-to-speech-поставщики. Если данные передаются третьим лицам, это делается в пределах целей обработки и необходимых функций сервиса.",
        ],
      },
      {
        id: "minors",
        title: "7. Несовершеннолетние пользователи",
        blocks: [
          "Если сервисом пользуется несовершеннолетний, согласие на обработку его персональных данных дает законный представитель: родитель, усыновитель, опекун или попечитель. Законный представитель отвечает за контроль использования сервиса и оплату платных функций.",
        ],
      },
      {
        id: "rights",
        title: "8. Права пользователя и контакты",
        blocks: [
          "Пользователь может запросить доступ к данным, уточнение, блокирование, удаление, ограничение обработки или отзыв согласия. Обращения принимаются на supportpoliglotai@gmail.com.",
          "После отзыва согласия часть функций сервиса может стать недоступной, если обработка данных необходима для их работы.",
        ],
      },
    ],
  },
  terms: {
    title: "Условия использования",
    badge: "Правила сервиса",
    description: "Правила использования сайта, веб-приложения и Telegram-бота Poliglot AI.",
    updated: "Обновлено 24 июня 2026",
    sections: [
      {
        id: "acceptance",
        title: "1. Принятие условий",
        blocks: [
          "Используя сайт, открывая Telegram-бота, проходя регистрацию, отправляя сообщения, голосовые записи, изображения, оплачивая тариф или выполняя иные действия в сервисе, пользователь принимает настоящие Условия, Пользовательское соглашение, Политику обработки персональных данных и Согласие на обработку персональных данных.",
          "Если пользователь не согласен с документами, он должен прекратить использование сайта, веб-приложения и Telegram-бота.",
        ],
      },
      {
        id: "service",
        title: "2. Описание сервиса",
        blocks: [
          "Poliglot AI предоставляет учебные материалы и автоматизированные AI-функции для изучения иностранных языков: уроки, ролевые диалоги, словарь ошибок, голосовую практику, распознавание речи, перевод текста с изображений, статистику прогресса, тарифы и поддержку.",
          "Сервис не является образовательной организацией, не выдает документы об образовании и не гарантирует конкретный уровень владения языком, экзаменационный результат или трудоустройство.",
        ],
      },
      {
        id: "account",
        title: "3. Аккаунт, Telegram и безопасность",
        blocks: [
          "Для отдельных функций нужен web-аккаунт, Telegram-аккаунт или привязка Telegram к web-профилю. Пользователь отвечает за сохранность логина, пароля, устройства, номера телефона, Telegram-аккаунта и кодов подтверждения.",
          "Пользователь передает только собственные данные либо данные, в отношении которых у него есть законное основание и необходимые согласия.",
        ],
      },
      {
        id: "plans",
        title: "4. Тарифы и платежи",
        blocks: [
          "В сервисе могут быть бесплатные и платные функции. Стоимость, период доступа, лимиты и состав тарифа отображаются в интерфейсе сайта, веб-приложения, Telegram-бота или платежного провайдера на момент оплаты.",
          "Оплата подтверждает намерение пользователя получить доступ к выбранным платным функциям. Возвраты и спорные платежи обрабатываются по применимому законодательству, правилам платежного провайдера и опубликованным условиям сервиса.",
        ],
      },
      {
        id: "ai",
        title: "5. AI-функции и учебные материалы",
        blocks: [
          "AI-ответы, переводы, оценки произношения и учебные рекомендации являются автоматизированными учебными подсказками. Они могут содержать неточности и требуют самостоятельной оценки пользователем.",
          "Запрещено отправлять через сервис незаконные материалы, персональные данные третьих лиц без основания, вредоносный код, спам, оскорбления и контент, нарушающий права других лиц.",
        ],
      },
      {
        id: "data",
        title: "6. Персональные данные и cookie",
        blocks: [
          "Обработка персональных данных осуществляется по Политике обработки персональных данных и отдельному Согласию на обработку персональных данных.",
          "Cookie и локальное хранилище используются для авторизации, безопасности, сохранения настроек, языка, темы, согласий, диагностики ошибок и улучшения работы сервиса.",
        ],
      },
      {
        id: "liability",
        title: "7. Ответственность",
        blocks: [
          "Сервис предоставляется на условиях \"как есть\" и \"как доступно\". Оператор не гарантирует бесперебойную работу, отсутствие ошибок, достижение конкретного учебного результата или постоянную доступность сторонних платформ.",
          "Оператор не отвечает за действия Telegram, платежных систем, хостинга, интернет-провайдеров, AI-поставщиков, банков и иных третьих лиц, если иное не предусмотрено обязательными нормами закона.",
        ],
      },
      {
        id: "contacts",
        title: "8. Изменение условий и контакты",
        blocks: [
          "Оператор вправе обновлять условия. Новая редакция вступает в силу с момента публикации на сайте, если в ней не указан иной срок. Продолжение использования сервиса означает принятие обновленных условий.",
          "Контакты для обращений: supportpoliglotai@gmail.com, Telegram @AsaselD, Telegram-бот @poliglot_ai_bot.",
        ],
      },
    ],
  },
  agreement: {
    title: "Пользовательское соглашение",
    badge: "Публичная оферта",
    description: "Публичное пользовательское соглашение для сайта, веб-приложения и Telegram-бота Poliglot AI.",
    updated: "Обновлено 24 июня 2026",
    sections: [
      {
        id: "offer",
        title: "1. Общие положения и акцепт",
        blocks: [
          "Настоящее соглашение является публичным предложением оператора заключить договор об использовании сервиса Poliglot AI на изложенных ниже условиях.",
          "Акцептом считается регистрация, вход в web-аккаунт, нажатие кнопки продолжения в Telegram-боте, оплата тарифа, отправка сообщений, голосовых записей, изображений или иное фактическое использование сервиса.",
          "С момента акцепта договор считается заключенным. Пользователь подтверждает, что прочитал, понял и принимает условия без изъятий и ограничений.",
        ],
      },
      {
        id: "subject",
        title: "2. Предмет договора",
        blocks: [
          "Оператор предоставляет пользователю доступ к сайту, веб-приложению и Telegram-боту для изучения иностранных языков, а пользователь использует сервис в соответствии с настоящим соглашением и законодательством Российской Федерации.",
          "Доступ обеспечивается через web-интерфейс, PWA, Telegram-бот @poliglot_ai_bot и связанные технические функции.",
        ],
      },
      {
        id: "plans",
        title: "3. Услуги, тарифы и оплата",
        blocks: [
          "Параметры бесплатного доступа, платных тарифов, лимитов, периода действия и стоимости определяются оператором и отображаются в интерфейсе сервиса на момент использования или оплаты.",
          "Оплата производится в форме предварительной оплаты через доступные платежные инструменты: платежного провайдера, Telegram Stars, YooKassa/СБП, TON, USDT или иной способ, если он показан пользователю в интерфейсе.",
          "Оператор вправе изменять стоимость и состав тарифов для будущих периодов. Изменения не ухудшают уже оплаченный период, если иное не требуется законом.",
        ],
      },
      {
        id: "duties",
        title: "4. Права и обязанности сторон",
        blocks: [
          "Оператор обязуется предоставлять доступ к сервису в пределах технической доступности, принимать разумные меры защиты данных, публиковать актуальные документы и отвечать на обращения пользователя.",
          "Пользователь обязуется использовать сервис законно, не передавать доступ третьим лицам, не нарушать работу сервиса, не отправлять незаконный контент и своевременно обновлять свои контактные данные.",
        ],
      },
      {
        id: "restrictions",
        title: "5. Ограничения использования",
        blocks: [
          "Запрещено использовать сервис для противоправных действий, взлома аккаунтов, DDoS-атак, спама, распространения вредоносного кода, нарушения интеллектуальных прав, передачи незаконного контента или персональных данных третьих лиц без основания.",
          "При нарушении условий оператор вправе ограничить или заблокировать доступ без предварительного уведомления, если это необходимо для защиты сервиса, пользователей или третьих лиц.",
        ],
      },
      {
        id: "minors",
        title: "6. Несовершеннолетние",
        blocks: [
          "Если сервис использует несовершеннолетний, его законный представитель принимает настоящее соглашение, дает согласие на обработку персональных данных несовершеннолетнего, контролирует использование сервиса и отвечает за оплату платных функций.",
        ],
      },
      {
        id: "liability",
        title: "7. Ответственность сторон",
        blocks: [
          "Стороны несут ответственность по настоящему соглашению и законодательству Российской Федерации. Сервис предоставляется на условиях \"как есть\" и \"как доступно\".",
          "Оператор не отвечает за качество интернет-соединения пользователя, сбои третьих лиц, ограничения Telegram, платежных провайдеров, банков, хостинга, AI-поставщиков и иных платформ.",
        ],
      },
      {
        id: "disputes",
        title: "8. Споры, срок действия и изменения",
        blocks: [
          "Стороны стремятся урегулировать споры путем переговоров и письменного обращения в поддержку. Срок рассмотрения обращения составляет до 15 рабочих дней с момента получения.",
          "Соглашение действует до прекращения использования сервиса или полного исполнения обязательств. Оператор может обновлять соглашение, а новая редакция действует с момента публикации на сайте.",
        ],
      },
    ],
  },
  consent: {
    title: "Согласие на обработку персональных данных",
    badge: "152-ФЗ",
    description: "Отдельное согласие пользователя на обработку персональных данных в Poliglot AI.",
    updated: "Обновлено 24 июня 2026",
    sections: [
      {
        id: "consent",
        title: "1. Согласие пользователя",
        blocks: [
          "Пользователь свободно, своей волей и в своем интересе дает оператору согласие на обработку персональных данных при использовании сайта poliglotai.ru, веб-приложения poliglotai.online/app и Telegram-бота @poliglot_ai_bot.",
          "Если данные относятся к несовершеннолетнему, согласие дает его законный представитель: родитель, усыновитель, опекун или попечитель.",
        ],
      },
      {
        id: "data",
        title: "2. Состав персональных данных",
        blocks: [
          "Согласие распространяется на email или логин web-аккаунта, Telegram ID, username, имя в Telegram, настройки языка, уровень обучения, учебные ответы, сообщения, голосовые записи, изображения, результаты распознавания и перевода, историю занятий, прогресс, платежные статусы и обращения в поддержку.",
          "Согласие также распространяется на технические данные: IP-адрес, cookie, локальное хранилище браузера, сведения об устройстве, браузере, операционной системе, времени доступа, источнике перехода, событиях безопасности и ошибках.",
        ],
      },
      {
        id: "purposes",
        title: "3. Цели обработки",
        blocks: [
          "Данные обрабатываются для регистрации, авторизации, единого профиля web и Telegram, уроков, практики, сохранения прогресса, проверки ответов, распознавания речи и изображений, работы AI-функций, оплаты тарифов, поддержки, безопасности, предотвращения злоупотреблений, диагностики ошибок и улучшения сервиса.",
        ],
      },
      {
        id: "actions",
        title: "4. Действия с данными",
        blocks: [
          "Согласие дается на сбор, запись, систематизацию, накопление, хранение, уточнение, использование, передачу в пределах целей обработки, обезличивание, блокирование, удаление и уничтожение данных с использованием средств автоматизации и без них.",
          "Для работы сервиса данные могут передаваться техническим подрядчикам и платформам, включая хостинг, Telegram, платежных провайдеров, аналитику, AI-, OCR-, speech-to-text- и text-to-speech-поставщиков.",
        ],
      },
      {
        id: "localization",
        title: "5. Хранение и трансграничная обработка",
        blocks: [
          "Первичный сбор и хранение персональных данных граждан Российской Федерации должны осуществляться с учетом требований российского законодательства о локализации персональных данных.",
          "Зарубежная инфраструктура может использоваться для технической обработки в пределах опубликованных целей, если это необходимо для работы сайта, веб-приложения, Telegram-бота, AI-функций, платежей или поддержки.",
        ],
      },
      {
        id: "term",
        title: "6. Срок действия и отзыв",
        blocks: [
          "Согласие действует с момента предоставления до достижения целей обработки, удаления аккаунта, прекращения использования сервиса или отзыва согласия, если более длительное хранение не требуется законом или для защиты прав оператора.",
          "Пользователь может отозвать согласие, запросить доступ, исправление, блокирование или удаление данных по адресу supportpoliglotai@gmail.com. После отзыва часть функций может стать недоступной.",
        ],
      },
    ],
  },
};

const englishNav = {
  home: "Poliglot AI",
  privacy: "Privacy",
  terms: "Terms",
  agreement: "User Agreement",
  consent: "Personal data processing consent",
};

const englishCookie = {
  title: "Cookie and technical data",
  body: "We use necessary cookies and local storage for login, language, theme, security, saved consents, and stable site operation. Optional cookies are used only after consent.",
  accept: "Accept all cookies",
  necessary: "Necessary only",
  privacy: "Privacy",
  consent: "Consent",
};

const operatorLabels = {
  ru: {
    operator: "Оператор",
    inn: "ИНН",
    address: "Адрес",
    email: "Email",
    bot: "Telegram bot",
  },
  en: {
    operator: "Operator",
    inn: "INN",
    address: "Address",
    email: "Email",
    bot: "Telegram bot",
  },
};

const englishDocuments = {
  privacy: {
    title: "Personal Data Processing Policy",
    badge: "Data protection",
    description: "How Poliglot AI processes data of website, web app, and Telegram bot users.",
    updated: "Updated 24 June 2026",
    sections: [
      {
        id: "scope",
        title: "1. General provisions",
        blocks: [
          "This Policy applies to Poliglot AI: poliglotai.ru, the web app at poliglotai.online/app, and the Telegram bot @poliglot_ai_bot.",
          "The service is a single system for learning foreign languages. By using it, the user confirms that they have read this Policy, the User Agreement, and the Consent to Personal Data Processing.",
        ],
      },
      {
        id: "operator",
        title: "2. Personal data operator",
        blocks: [
          "Personal data operator: self-employed Cheban Denis Igorevich, INN 505017471160, address: Shchyolkovo, Sirenevaya St., 9k1, apt. 9. Contact email: supportpoliglotai@gmail.com.",
        ],
      },
      {
        id: "data",
        title: "3. Data processed",
        blocks: [
          "The operator may process an email or web account login, Telegram ID, username, Telegram name, selected languages, learning level, learning answers, messages, voice recordings, images, recognition results, translations, progress, payment statuses, and support requests.",
          "Technical data may also be processed: IP address, cookies, browser local storage, device data, browser, operating system, access time, referrer, security events, and error data.",
          "The operator does not request passport data, medical data, biometric data, special categories of personal data, or full bank card details.",
        ],
      },
      {
        id: "purposes",
        title: "4. Purposes and legal bases",
        blocks: [
          "Data is needed for registration, login, a shared web and Telegram profile, lessons and practice, progress storage, interface settings, plan payments, user support, security, abuse prevention, error diagnostics, and service improvement.",
          "Processing is based on user consent, performance of the user agreement, legal requirements, and the operator's legitimate interest in keeping the service secure.",
        ],
      },
      {
        id: "cookies",
        title: "5. Cookies and local storage",
        blocks: [
          "The website and web app use necessary cookies and local storage for login, language, theme, consent records, session, account protection, and stable interface operation.",
          "Optional analytics or marketing cookies are used only after separate consent if such technologies are enabled in the current version of the site.",
        ],
      },
      {
        id: "transfer",
        title: "6. Storage, transfer, and infrastructure",
        blocks: [
          "Initial collection and storage of personal data of Russian Federation citizens must be carried out with due regard to Russian personal data localization requirements.",
          "The service may use hosting, Telegram, payment providers, analytics, AI, OCR, speech-to-text, and text-to-speech providers. If data is transferred to third parties, this is done within the stated purposes and necessary service functions.",
        ],
      },
      {
        id: "minors",
        title: "7. Minors",
        blocks: [
          "If a minor uses the service, consent to personal data processing is given by the legal representative: parent, adoptive parent, guardian, or custodian. The legal representative controls service use and paid functions.",
        ],
      },
      {
        id: "rights",
        title: "8. User rights and contacts",
        blocks: [
          "The user may request access to data, correction, blocking, deletion, restriction of processing, or withdrawal of consent. Requests are accepted at supportpoliglotai@gmail.com.",
          "After consent is withdrawn, some service functions may become unavailable if data processing is necessary for their operation.",
        ],
      },
    ],
  },
  terms: {
    title: "Terms of Use",
    badge: "Service rules",
    description: "Rules for using the Poliglot AI website, web app, and Telegram bot.",
    updated: "Updated 24 June 2026",
    sections: [
      {
        id: "acceptance",
        title: "1. Acceptance of terms",
        blocks: [
          "By using the website, opening the Telegram bot, registering, sending messages, voice recordings or images, paying for a plan, or otherwise using the service, the user accepts these Terms, the User Agreement, the Personal Data Processing Policy, and the Consent to Personal Data Processing.",
          "If the user does not agree with the documents, they must stop using the website, web app, and Telegram bot.",
        ],
      },
      {
        id: "service",
        title: "2. Service description",
        blocks: [
          "Poliglot AI provides learning materials and automated AI functions for foreign language learning: lessons, roleplay, mistake vocabulary, voice practice, speech recognition, image text translation, progress statistics, plans, and support.",
          "The service is not an educational organization, does not issue education certificates, and does not guarantee a specific language level, exam result, or employment outcome.",
        ],
      },
      {
        id: "account",
        title: "3. Account, Telegram, and security",
        blocks: [
          "Some functions require a web account, Telegram account, or linking Telegram to the web profile. The user is responsible for protecting the login, password, device, phone number, Telegram account, and confirmation codes.",
          "The user provides only their own data or data for which they have a lawful basis and the necessary consents.",
        ],
      },
      {
        id: "plans",
        title: "4. Plans and payments",
        blocks: [
          "The service may include free and paid functions. Price, access period, limits, and plan contents are shown in the website, web app, Telegram bot, or payment provider interface at the time of payment.",
          "Payment confirms the user's intent to receive access to selected paid functions. Refunds and disputed payments are handled under applicable law, payment provider rules, and published service terms.",
        ],
      },
      {
        id: "ai",
        title: "5. AI functions and learning materials",
        blocks: [
          "AI answers, translations, pronunciation scores, and learning recommendations are automated learning hints. They may contain inaccuracies and require the user's independent assessment.",
          "It is prohibited to send illegal materials, third-party personal data without a lawful basis, malware, spam, insults, or content that violates third-party rights through the service.",
        ],
      },
      {
        id: "data",
        title: "6. Personal data and cookies",
        blocks: [
          "Personal data is processed under the Personal Data Processing Policy and the separate Consent to Personal Data Processing.",
          "Cookies and local storage are used for login, security, saved settings, language, theme, consent records, error diagnostics, and service improvement.",
        ],
      },
      {
        id: "liability",
        title: "7. Liability",
        blocks: [
          "The service is provided as is and as available. The operator does not guarantee uninterrupted operation, absence of errors, a specific learning result, or permanent availability of third-party platforms.",
          "The operator is not responsible for actions of Telegram, payment systems, hosting, internet providers, AI providers, banks, and other third parties, unless mandatory law provides otherwise.",
        ],
      },
      {
        id: "contacts",
        title: "8. Changes and contacts",
        blocks: [
          "The operator may update the terms. A new version takes effect upon publication on the site unless it specifies another effective date. Continued service use means acceptance of the updated terms.",
          "Contacts: supportpoliglotai@gmail.com, Telegram @AsaselD, Telegram bot @poliglot_ai_bot.",
        ],
      },
    ],
  },
  agreement: {
    title: "User Agreement",
    badge: "Public offer",
    description: "Public user agreement for the Poliglot AI website, web app, and Telegram bot.",
    updated: "Updated 24 June 2026",
    sections: [
      {
        id: "offer",
        title: "1. General provisions and acceptance",
        blocks: [
          "This agreement is a public offer by the operator to conclude an agreement for using Poliglot AI on the terms below.",
          "Acceptance includes registration, signing in to a web account, pressing the continue button in the Telegram bot, paying for a plan, sending messages, voice recordings, images, or otherwise using the service.",
          "From acceptance, the agreement is considered concluded. The user confirms that they have read, understood, and accept the terms without exclusions or limitations.",
        ],
      },
      {
        id: "subject",
        title: "2. Subject of the agreement",
        blocks: [
          "The operator provides access to the website, web app, and Telegram bot for foreign language learning, and the user uses the service according to this agreement and Russian Federation law.",
          "Access is provided through the web interface, PWA, Telegram bot @poliglot_ai_bot, and related technical functions.",
        ],
      },
      {
        id: "plans",
        title: "3. Services, plans, and payment",
        blocks: [
          "Free access, paid plans, limits, access period, and price are determined by the operator and shown in the service interface at the time of use or payment.",
          "Payment is made as prepayment through available tools: payment provider, Telegram Stars, YooKassa/SBP, TON, USDT, or another method shown in the interface.",
          "The operator may change prices and plan contents for future periods. Changes do not worsen an already paid period unless required by law.",
        ],
      },
      {
        id: "duties",
        title: "4. Rights and obligations",
        blocks: [
          "The operator undertakes to provide service access within technical availability, take reasonable data protection measures, publish current documents, and respond to user requests.",
          "The user undertakes to use the service lawfully, not transfer access to third parties, not disrupt the service, not send illegal content, and keep contact data current.",
        ],
      },
      {
        id: "restrictions",
        title: "5. Use restrictions",
        blocks: [
          "The service may not be used for illegal actions, account hacking, DDoS attacks, spam, malware distribution, intellectual property violations, illegal content, or third-party personal data without a lawful basis.",
          "If the terms are violated, the operator may restrict or block access without prior notice when necessary to protect the service, users, or third parties.",
        ],
      },
      {
        id: "minors",
        title: "6. Minors",
        blocks: [
          "If a minor uses the service, the legal representative accepts this agreement, gives consent to processing the minor's personal data, controls service use, and is responsible for paid functions.",
        ],
      },
      {
        id: "liability",
        title: "7. Liability of the parties",
        blocks: [
          "The parties are liable under this agreement and Russian Federation law. The service is provided as is and as available.",
          "The operator is not responsible for the user's internet connection quality, failures of third parties, Telegram restrictions, payment providers, banks, hosting, AI providers, and other platforms.",
        ],
      },
      {
        id: "disputes",
        title: "8. Disputes, term, and changes",
        blocks: [
          "The parties seek to resolve disputes by negotiation and written support requests. The response period is up to 15 business days from receipt.",
          "The agreement remains in force until service use stops or obligations are fully performed. The operator may update the agreement, and the new version applies from publication on the site.",
        ],
      },
    ],
  },
  consent: {
    title: "Consent to Personal Data Processing",
    badge: "152-FZ",
    description: "Separate user consent to personal data processing in Poliglot AI.",
    updated: "Updated 24 June 2026",
    sections: [
      {
        id: "consent",
        title: "1. User consent",
        blocks: [
          "The user freely, by their own will and in their own interest, gives the operator consent to process personal data while using poliglotai.ru, the web app at poliglotai.online/app, and the Telegram bot @poliglot_ai_bot.",
          "If data relates to a minor, consent is given by the legal representative: parent, adoptive parent, guardian, or custodian.",
        ],
      },
      {
        id: "data",
        title: "2. Personal data scope",
        blocks: [
          "Consent covers an email or web account login, Telegram ID, username, Telegram name, language settings, learning level, learning answers, messages, voice recordings, images, recognition and translation results, lesson history, progress, payment statuses, and support requests.",
          "Consent also covers technical data: IP address, cookies, browser local storage, device, browser, operating system, access time, referrer, security events, and errors.",
        ],
      },
      {
        id: "purposes",
        title: "3. Processing purposes",
        blocks: [
          "Data is processed for registration, login, a shared web and Telegram profile, lessons, practice, progress storage, answer checking, speech and image recognition, AI functions, plan payment, support, security, abuse prevention, error diagnostics, and service improvement.",
        ],
      },
      {
        id: "actions",
        title: "4. Data operations",
        blocks: [
          "Consent is given for collection, recording, systematization, accumulation, storage, clarification, use, transfer within the processing purposes, depersonalization, blocking, deletion, and destruction of data with and without automation tools.",
          "For service operation, data may be transferred to technical contractors and platforms including hosting, Telegram, payment providers, analytics, AI, OCR, speech-to-text, and text-to-speech providers.",
        ],
      },
      {
        id: "localization",
        title: "5. Storage and cross-border processing",
        blocks: [
          "Initial collection and storage of personal data of Russian Federation citizens must be carried out with due regard to Russian personal data localization requirements.",
          "Foreign infrastructure may be used for technical processing within the published purposes if needed for the website, web app, Telegram bot, AI functions, payments, or support.",
        ],
      },
      {
        id: "term",
        title: "6. Term and withdrawal",
        blocks: [
          "Consent is valid from the moment it is given until processing purposes are achieved, account deletion, service use termination, or consent withdrawal, unless longer storage is required by law or to protect the operator's rights.",
          "The user may withdraw consent or request access, correction, blocking, or deletion at supportpoliglotai@gmail.com. After withdrawal, some functions may become unavailable.",
        ],
      },
    ],
  },
};

const englishLegalPack = {
  nav: englishNav,
  cookie: englishCookie,
  operator: {
    ...operator,
    title: "Operator and details",
    labels: operatorLabels.en,
  },
  documents: englishDocuments,
};

const immutablePatterns = [
  /Poliglot AI/g,
  /poliglotai\.ru/g,
  /poliglotai\.online\/app/g,
  /@poliglot_ai_bot/g,
  /@AsaselD/g,
  /supportpoliglotai@gmail\.com/g,
  /505017471160/g,
  /Telegram/g,
  /YooKassa/g,
  /TON/g,
  /USDT/g,
  /PWA/g,
  /AI/g,
  /OCR/g,
  /speech-to-text/g,
  /text-to-speech/g,
  /cookie/g,
];

function collectStrings() {
  const values = new Set();
  for (const item of [...Object.values(nav), ...Object.values(cookie), operator.title]) values.add(item);
  for (const doc of Object.values(documents)) {
    values.add(doc.title);
    values.add(doc.badge);
    values.add(doc.description);
    values.add(doc.updated);
    for (const section of doc.sections) {
      values.add(section.title);
      section.blocks.forEach((block) => values.add(block));
    }
  }
  return [...values];
}

function protectText(text, prefix = "") {
  const tokens = [];
  let next = text;
  immutablePatterns.forEach((pattern) => {
    next = next.replace(pattern, (match) => {
      const token = `__POLIGLOT_${prefix}_${tokens.length}__`;
      tokens.push([token, match]);
      return token;
    });
  });
  return { text: next, tokens };
}

function restoreText(text, tokens) {
  let next = text;
  for (const [token, value] of tokens) {
    next = next.replaceAll(token, value);
    next = next.replaceAll(token.toLowerCase(), value);
  }
  return next;
}

async function translateOne(text, target) {
  if (target === "ru") return text;
  const protectedValue = protectText(text);
  const params = new URLSearchParams({ client: "gtx", sl: "ru", tl: target, dt: "t", q: protectedValue.text });
  const response = await fetchWithRetry(`https://translate.googleapis.com/translate_a/single?${params.toString()}`);
  if (!response.ok) throw new Error(`translate ${target} failed with ${response.status}`);
  const json = await response.json();
  const translated = (json?.[0] || []).map((part) => part?.[0] || "").join("");
  return restoreText(translated || text, protectedValue.tokens);
}

async function fetchWithRetry(url, attempts = 4) {
  let lastError;
  for (let attempt = 1; attempt <= attempts; attempt += 1) {
    try {
      const response = await fetch(url);
      if (response.ok || response.status < 500) return response;
      lastError = new Error(`HTTP ${response.status}`);
    } catch (error) {
      lastError = error;
    }
    await new Promise((resolve) => setTimeout(resolve, 750 * attempt));
  }
  throw lastError;
}

async function translateBatch(values, target) {
  if (target === "ru") return values;
  if (values.length === 1) return [await translateOne(values[0], target)];

  const delimiter = "@@POLIGLOT_SPLIT_20260624@@";
  const protectedValues = values.map((value, index) => protectText(value, String(index)));
  const tokens = protectedValues.flatMap((entry) => entry.tokens);
  const source = protectedValues.map((entry) => entry.text).join(`\n${delimiter}\n`);
  const params = new URLSearchParams({ client: "gtx", sl: "ru", tl: target, dt: "t", q: source });
  const response = await fetchWithRetry(`https://translate.googleapis.com/translate_a/single?${params.toString()}`);

  if (!response.ok) {
    const midpoint = Math.ceil(values.length / 2);
    return [
      ...(await translateBatch(values.slice(0, midpoint), target)),
      ...(await translateBatch(values.slice(midpoint), target)),
    ];
  }

  const json = await response.json();
  const translated = (json?.[0] || []).map((part) => part?.[0] || "").join("");
  const parts = translated.split(delimiter).map((part) => part.trim());
  if (parts.length !== values.length) {
    const midpoint = Math.ceil(values.length / 2);
    return [
      ...(await translateBatch(values.slice(0, midpoint), target)),
      ...(await translateBatch(values.slice(midpoint), target)),
    ];
  }
  return parts.map((part) => restoreText(part, tokens));
}

async function translateAll(strings) {
  const translations = { ru: Object.fromEntries(strings.map((value) => [value, value])) };
  for (const [code] of languages) {
    if (code === "ru") continue;
    translations[code] = {};
    for (let index = 0; index < strings.length; index += 10) {
      const batch = strings.slice(index, index + 10);
      const translated = await translateBatch(batch, code);
      batch.forEach((value, offset) => {
        translations[code][value] = translated[offset] || value;
      });
      await new Promise((resolve) => setTimeout(resolve, 30));
    }
    console.log(`translated ${code}`);
  }
  return translations;
}

function localizeShape(value, code, translations) {
  if (typeof value === "string") return translations[code]?.[value] || translations.en?.[value] || value;
  if (Array.isArray(value)) return value.map((item) => localizeShape(item, code, translations));
  if (value && typeof value === "object") {
    return Object.fromEntries(Object.entries(value).map(([key, child]) => [key, localizeShape(child, code, translations)]));
  }
  return value;
}

function buildRendererPayload(translations) {
  const localized = {};
  for (const [code] of languages) {
    if (code === "ru" && !translations) {
      localized[code] = {
        nav,
        cookie,
        operator: {
          ...operator,
          labels: operatorLabels.ru,
        },
        documents,
      };
      continue;
    }
    if (code !== "ru" && !translations) {
      localized[code] = englishLegalPack;
      continue;
    }
    localized[code] = {
      nav: localizeShape(nav, code, translations),
      cookie: localizeShape(cookie, code, translations),
      operator: {
        ...operator,
        title: localizeShape(operator.title, code, translations),
        labels: operatorLabels.ru,
      },
      documents: localizeShape(documents, code, translations),
    };
  }
  return { languages, localized };
}

function rendererSource(payload) {
  return `(function () {
  const payload = ${JSON.stringify(payload)};
  const supported = new Set(payload.languages.map(([code]) => code));
  const docOrder = ["privacy", "terms", "agreement", "consent"];
  const page = document.documentElement.dataset.sitePage || "";
  const storageKeys = ["poliglot_site_language", "poliglot-site-language", "poliglot-auth-language"];

  function normalizedLanguage(value) {
    const raw = String(value || "").toLowerCase().trim();
    if (!raw) return "";
    const code = raw.split(/[-_]/)[0];
    return supported.has(code) ? code : "";
  }

  function currentLanguage() {
    const params = new URLSearchParams(location.search);
    for (const value of [params.get("lang"), localStorage.getItem(storageKeys[0]), localStorage.getItem(storageKeys[1]), localStorage.getItem(storageKeys[2]), document.documentElement.lang]) {
      const code = normalizedLanguage(value);
      if (code) return code;
    }
    return "ru";
  }

  function escapeHTML(value) {
    return String(value || "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;");
  }

  function linkify(text) {
    return escapeHTML(text)
      .replace(/Политик[^,.]+персональных данных|Personal Data Processing Policy|Privacy Policy/g, '<a href="/privacy.html">$&</a>')
      .replace(/Пользовательск[^,.]+соглашени[^,.]+|User Agreement/g, '<a href="/agreement.html">$&</a>')
      .replace(/Соглас[^,.]+персональных данных|Consent to Personal Data Processing/g, '<a href="/consent.html">$&</a>')
      .replace(/Услов[^,.]+использования|Terms of Use/g, '<a href="/terms.html">$&</a>');
  }

  function renderOperator(operator) {
    const labels = operator.labels || {};
    return '<section id="operator-details" class="legal-operator-card">' +
      '<h2>' + escapeHTML(operator.title) + '</h2>' +
      '<div class="legal-contact-grid">' +
      '<span><small>' + escapeHTML(labels.operator || "Operator") + '</small><strong>' + escapeHTML(operator.name) + '</strong></span>' +
      '<span><small>' + escapeHTML(labels.inn || "INN") + '</small><strong>' + escapeHTML(operator.inn) + '</strong></span>' +
      '<span><small>' + escapeHTML(labels.address || "Address") + '</small><strong>' + escapeHTML(operator.address) + '</strong></span>' +
      '<a href="mailto:' + escapeHTML(operator.email) + '"><small>' + escapeHTML(labels.email || "Email") + '</small><strong>' + escapeHTML(operator.email) + '</strong></a>' +
      '<a href="https://t.me/' + escapeHTML(operator.bot.replace(/^@/, "")) + '"><small>' + escapeHTML(labels.bot || "Telegram bot") + '</small><strong>' + escapeHTML(operator.bot) + '</strong></a>' +
      '</div>' +
      '</section>';
  }

  function renderDocument(doc, operator) {
    return renderOperator(operator) +
      '<header class="legal-document-heading"><p class="eyebrow">' + escapeHTML(doc.updated) + '</p><h2>' + escapeHTML(doc.title) + '</h2><p>' + escapeHTML(doc.description) + '</p></header>' +
      doc.sections.map((section) => '<section id="' + escapeHTML(section.id) + '">' +
        '<h2>' + escapeHTML(section.title) + '</h2>' +
        section.blocks.map((block) => '<p>' + linkify(block) + '</p>').join('') +
      '</section>').join('');
  }

  function localHref(path, lang) {
    const url = new URL(path, location.origin);
    url.searchParams.set("lang", lang);
    return url.pathname + url.search;
  }

  function apply() {
    const lang = currentLanguage();
    const pack = payload.localized[lang] || payload.localized.ru;
    localStorage.setItem(storageKeys[0], lang);
    document.documentElement.lang = lang;

    if (docOrder.includes(page)) {
      const doc = pack.documents[page];
      const shell = document.querySelector(".legal-document-shell");
      if (shell && doc) {
        shell.innerHTML = renderDocument(doc, pack.operator);
      }
      document.title = doc ? doc.title + " - Poliglot AI" : document.title;
      document.querySelectorAll("[data-legal-title]").forEach((node) => { node.textContent = doc?.title || node.textContent; });
      document.querySelectorAll("[data-legal-badge]").forEach((node) => { node.textContent = doc?.badge || node.textContent; });
      document.querySelectorAll("[data-legal-description]").forEach((node) => { node.textContent = doc?.description || node.textContent; });
    }

    document.querySelectorAll("[data-legal-nav]").forEach((node) => {
      const key = node.getAttribute("data-legal-nav");
      node.textContent = pack.nav[key] || node.textContent;
    });

    document.querySelectorAll("[data-legal-cookie]").forEach((node) => {
      const key = node.getAttribute("data-legal-cookie");
      node.textContent = pack.cookie[key] || node.textContent;
    });

    document.querySelectorAll("a[href]").forEach((link) => {
      const original = link.dataset.legalOriginalHref || link.dataset.originalHref || link.getAttribute("href");
      if (!original || original.startsWith("#") || original.startsWith("mailto:") || original.startsWith("tel:") || /^https?:\\/\\//i.test(original)) return;
      link.dataset.legalOriginalHref = original;
      if (/\\/(privacy|terms|agreement|consent|poliglot-ai)\\.html$/.test(new URL(original, location.origin).pathname)) {
        link.setAttribute("href", localHref(original, lang));
      }
    });
  }

  window.poliglotLegalDocumentsI18n = { apply, payload };
  window.addEventListener("poliglot-language-change", apply);
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", apply, { once: true });
  } else {
    apply();
  }
})();`;
}

const translations = process.env.POLIGLOT_TRANSLATE_LEGAL === "1" ? await translateAll(collectStrings()) : undefined;
const payload = buildRendererPayload(translations);
await writeFile(outputPath, rendererSource(payload), "utf8");
console.log(`wrote ${outputPath}`);
