(function () {
    const languages = [
        ["ru", "Русский"], ["en", "English"], ["es", "Español"], ["de", "Deutsch"],
        ["fr", "Français"], ["it", "Italiano"], ["zh", "中文"], ["ja", "日本語"], ["ko", "한국어"], ["tg", "Тоҷикӣ"], ["uz", "O‘zbekcha"],
        ["tt", "Татарча"], ["hy", "Հայերեն"], ["kk", "Қазақша"], ["ky", "Кыргызча"],
        ["ka", "ქართული"], ["uk", "Українська"], ["pl", "Polski"], ["ro", "Română"],
        ["pt", "Português"], ["ar", "العربية"], ["bn", "বাংলা"], ["cs", "Čeština"],
        ["el", "Ελληνικά"], ["hi", "हिंदी"], ["hu", "Magyar"], ["id", "Bahasa Indonesia"],
        ["nl", "Nederlands"], ["sv", "svenska"], ["ta", "தமிழ்"], ["te", "తెలుగు"],
        ["th", "ภาษาไทย"], ["tl", "Tagalog"], ["tr", "Türkçe"], ["vi", "Tiếng Việt"]
    ];

    const copy = {
        ru: {
            language_label: "Язык",
            brand_subtitle: "языки онлайн и в Telegram",
            brand_subtitle_short: "сайт и Telegram",
            nav_home: "Главная",
            nav_features: "Возможности",
            nav_online: "Онлайн",
            nav_compare: "Сравнение",
            nav_pricing: "Тарифы",
            nav_faq: "FAQ",
            nav_bot: "Онлайн-бот",
            nav_policy: "Политика",
            nav_terms: "Условия",
            nav_policy_full: "Политика конфиденциальности",
            nav_terms_full: "Условия использования",
            open_online: "Открыть онлайн",
            open_app: "Открыть онлайн-приложение",
            open_bot: "Открыть бота",
            landing_badge: "AI-репетитор в Telegram и web app: уроки, голос, фото, ошибки и словарь",
            landing_title: "Заговорите на новом языке быстрее: AI-репетитор уже внутри Telegram",
            landing_description: "NERIVA превращает привычный мессенджер в личного языкового тренера: даёт короткие уроки, ведёт диалог, разбирает ошибки, тренирует слова, понимает голос и переводит текст с фото.",
            landing_pricing_line: "Вы сразу понимаете, что улучшить: точность фразы, слабые слова, плавность речи и акцент. Без сложных терминов — короткий совет и следующий шаг.",
            terms_badge: "Правила использования сервиса",
            terms_title: "Условия использования",
            terms_description: "Пользовательское соглашение регулирует использование сайта, веб-приложения и Telegram-бота NERIVA, включая тарифы, оплату, персональные данные, ограничения сервиса и ответственность сторон.",
            privacy_badge: "Защита данных",
            privacy_title: "Политика конфиденциальности",
            privacy_description: "Здесь описано, какие данные может обрабатывать NERIVA, зачем они нужны для обучения, оплаты, поддержки, защиты аккаунта и реализации прав пользователя.",
            legal_pill_languages: "20 языков сейчас",
            legal_pill_interface: "Мультиязычный интерфейс",
            legal_pill_ai: "AI-практика",
            legal_pill_tariffs: "Free, Premium и Platinum"
        },
        en: {
            language_label: "Language",
            brand_subtitle: "languages online and in Telegram",
            brand_subtitle_short: "site and Telegram",
            nav_home: "Home",
            nav_features: "Features",
            nav_online: "Online",
            nav_compare: "Compare",
            nav_pricing: "Pricing",
            nav_faq: "FAQ",
            nav_bot: "Online bot",
            nav_policy: "Privacy",
            nav_terms: "Terms",
            nav_policy_full: "Privacy Policy",
            nav_terms_full: "Terms of Use",
            open_online: "Open online",
            open_app: "Open web app",
            open_bot: "Open bot",
            landing_badge: "AI tutor in Telegram and web app: lessons, voice, photo, mistakes and vocabulary",
            landing_title: "Speak a new language faster: your AI tutor is already inside Telegram",
            landing_description: "NERIVA turns your messenger into a personal language coach: short lessons, conversation practice, mistake review, vocabulary training, voice understanding and photo text translation.",
            landing_pricing_line: "You immediately see what to improve: phrase accuracy, weak words, speech flow and accent. No technical noise, just a short tip and the next step.",
            terms_badge: "Service rules",
            terms_title: "Terms of Use",
            terms_description: "These Terms govern the use of the NERIVA website, web app and Telegram bot, including plans, payments, personal data, service limits and liability.",
            privacy_badge: "Data protection",
            privacy_title: "Privacy Policy",
            privacy_description: "This page explains what data NERIVA may process and why it is needed for learning, payments, support, account protection and user rights.",
            legal_pill_languages: "20 languages now",
            legal_pill_interface: "Multilingual interface",
            legal_pill_ai: "AI practice",
            legal_pill_tariffs: "Free, Premium and Platinum"
        },
        es: {
            language_label: "Idioma", brand_subtitle: "idiomas online y en Telegram", brand_subtitle_short: "sitio y Telegram",
            nav_home: "Inicio", nav_features: "Funciones", nav_online: "Online", nav_compare: "Comparación", nav_pricing: "Tarifas", nav_faq: "FAQ", nav_bot: "Bot online", nav_policy: "Privacidad", nav_terms: "Términos", nav_policy_full: "Política de privacidad", nav_terms_full: "Términos de uso", open_online: "Abrir online", open_app: "Abrir app web", open_bot: "Abrir bot",
            landing_badge: "Servicio de IA multilingüe: 20 idiomas, app web, entrada por Telegram y traductor", landing_title: "Aprende idiomas online o en Telegram, en el idioma que prefieras", landing_description: "NERIVA funciona como un servicio único: app web, bot de Telegram, acceso por Telegram, progreso sincronizado y elección separada del idioma de interfaz y de aprendizaje.",
            terms_badge: "Reglas del servicio", terms_title: "Términos de uso", terms_description: "Estos términos regulan el uso del sitio, la app web y el bot de Telegram de NERIVA, incluidos planes, pagos, datos personales, límites del servicio y responsabilidad.",
            privacy_badge: "Protección de datos", privacy_title: "Política de privacidad", privacy_description: "Esta página explica qué datos puede procesar NERIVA y por qué son necesarios para aprender, pagar, recibir soporte, proteger la cuenta y ejercer derechos.",
            legal_pill_languages: "20 idiomas ahora", legal_pill_interface: "Interfaz multilingüe", legal_pill_ai: "Práctica con IA", legal_pill_tariffs: "Free y Premium"
        },
        de: {
            language_label: "Sprache", brand_subtitle: "Sprachen online und in Telegram", brand_subtitle_short: "Website und Telegram",
            nav_home: "Start", nav_features: "Funktionen", nav_online: "Online", nav_compare: "Vergleich", nav_pricing: "Preise", nav_faq: "FAQ", nav_bot: "Online-Bot", nav_policy: "Datenschutz", nav_terms: "Bedingungen", nav_policy_full: "Datenschutzerklärung", nav_terms_full: "Nutzungsbedingungen", open_online: "Online öffnen", open_app: "Web-App öffnen", open_bot: "Bot öffnen",
            landing_badge: "Mehrsprachiger KI-Service: 20 Sprachen, Web-App, Telegram-Login und Übersetzer", landing_title: "Lerne Sprachen online oder in Telegram, in deiner bevorzugten Sprache", landing_description: "NERIVA ist ein gemeinsamer Service: Web-App, Telegram-Bot, Telegram-Login, synchroner Fortschritt und getrennte Wahl von Oberfläche und Lernsprache.",
            terms_badge: "Serviceregeln", terms_title: "Nutzungsbedingungen", terms_description: "Diese Bedingungen regeln die Nutzung der NERIVA Website, Web-App und des Telegram-Bots, einschließlich Tarife, Zahlungen, personenbezogene Daten, Servicegrenzen und Haftung.",
            privacy_badge: "Datenschutz", privacy_title: "Datenschutzerklärung", privacy_description: "Diese Seite erklärt, welche Daten NERIVA verarbeiten kann und warum sie für Lernen, Zahlungen, Support, Kontoschutz und Nutzerrechte benötigt werden.",
            legal_pill_languages: "20 Sprachen jetzt", legal_pill_interface: "Mehrsprachige Oberfläche", legal_pill_ai: "KI-Praxis", legal_pill_tariffs: "Free und Premium"
        },
        fr: {
            language_label: "Langue", brand_subtitle: "langues en ligne et sur Telegram", brand_subtitle_short: "site et Telegram",
            nav_home: "Accueil", nav_features: "Fonctions", nav_online: "En ligne", nav_compare: "Comparaison", nav_pricing: "Tarifs", nav_faq: "FAQ", nav_bot: "Bot en ligne", nav_policy: "Confidentialité", nav_terms: "Conditions", nav_policy_full: "Politique de confidentialité", nav_terms_full: "Conditions d'utilisation", open_online: "Ouvrir en ligne", open_app: "Ouvrir l'app web", open_bot: "Ouvrir le bot",
            landing_badge: "Service IA multilingue : 20 langues, app web, connexion Telegram et traducteur", landing_title: "Apprenez les langues en ligne ou sur Telegram, dans la langue qui vous convient", landing_description: "NERIVA fonctionne comme un service unique : app web, bot Telegram, connexion Telegram, progression synchronisée et choix séparé de la langue d'interface et d'apprentissage.",
            terms_badge: "Règles du service", terms_title: "Conditions d'utilisation", terms_description: "Ces conditions régissent l'utilisation du site, de l'app web et du bot Telegram NERIVA, y compris les offres, paiements, données personnelles, limites du service et responsabilités.",
            privacy_badge: "Protection des données", privacy_title: "Politique de confidentialité", privacy_description: "Cette page explique quelles données NERIVA peut traiter et pourquoi elles sont nécessaires pour apprendre, payer, obtenir de l'aide, protéger le compte et exercer vos droits.",
            legal_pill_languages: "20 langues maintenant", legal_pill_interface: "Interface multilingue", legal_pill_ai: "Pratique IA", legal_pill_tariffs: "Free et Premium"
        },
        it: {
            language_label: "Lingua", brand_subtitle: "lingue online e su Telegram", brand_subtitle_short: "sito e Telegram",
            nav_home: "Home", nav_features: "Funzioni", nav_online: "Online", nav_compare: "Confronto", nav_pricing: "Tariffe", nav_faq: "FAQ", nav_bot: "Bot online", nav_policy: "Privacy", nav_terms: "Termini", nav_policy_full: "Informativa privacy", nav_terms_full: "Termini d'uso", open_online: "Apri online", open_app: "Apri app web", open_bot: "Apri bot",
            landing_badge: "Servizio IA multilingue: 20 lingue, app web, accesso Telegram e traduttore", landing_title: "Impara lingue online o su Telegram, nella lingua più comoda per te", landing_description: "NERIVA funziona come un unico servizio: app web, bot Telegram, accesso via Telegram, progresso sincronizzato e scelta separata della lingua dell'interfaccia e di studio.",
            terms_badge: "Regole del servizio", terms_title: "Termini d'uso", terms_description: "Questi termini regolano l'uso del sito, dell'app web e del bot Telegram NERIVA, inclusi piani, pagamenti, dati personali, limiti del servizio e responsabilità.",
            privacy_badge: "Protezione dei dati", privacy_title: "Informativa privacy", privacy_description: "Questa pagina spiega quali dati NERIVA può trattare e perché servono per studio, pagamenti, supporto, protezione dell'account e diritti dell'utente.",
            legal_pill_languages: "20 lingue ora", legal_pill_interface: "Interfaccia multilingue", legal_pill_ai: "Pratica IA", legal_pill_tariffs: "Free e Premium"
        },
        uk: {
            language_label: "Мова", brand_subtitle: "мови онлайн і в Telegram", brand_subtitle_short: "сайт і Telegram",
            nav_home: "Головна", nav_features: "Можливості", nav_online: "Онлайн", nav_compare: "Порівняння", nav_pricing: "Тарифи", nav_faq: "FAQ", nav_bot: "Онлайн-бот", nav_policy: "Політика", nav_terms: "Умови", nav_policy_full: "Політика конфіденційності", nav_terms_full: "Умови використання", open_online: "Відкрити онлайн", open_app: "Відкрити веб-додаток", open_bot: "Відкрити бота",
            landing_badge: "Багатомовний AI-сервіс: 20 мов, web app, Telegram-вхід і перекладач", landing_title: "Вивчайте мови онлайн або в Telegram зручною для вас мовою", landing_description: "NERIVA працює як єдиний сервіс: веб-додаток, Telegram-бот, вхід через Telegram, синхронізація прогресу та окремий вибір мови інтерфейсу й навчання.",
            terms_badge: "Правила сервісу", terms_title: "Умови використання", terms_description: "Ці умови регулюють використання сайту, веб-додатка й Telegram-бота NERIVA, включно з тарифами, оплатою, персональними даними, обмеженнями сервісу й відповідальністю.",
            privacy_badge: "Захист даних", privacy_title: "Політика конфіденційності", privacy_description: "Тут описано, які дані може обробляти NERIVA і навіщо вони потрібні для навчання, оплати, підтримки, захисту акаунта та прав користувача.",
            legal_pill_languages: "20 мов зараз", legal_pill_interface: "Багатомовний інтерфейс", legal_pill_ai: "AI-практика", legal_pill_tariffs: "Free і Premium"
        },
        pl: {
            language_label: "Język", brand_subtitle: "języki online i w Telegramie", brand_subtitle_short: "strona i Telegram",
            nav_home: "Główna", nav_features: "Funkcje", nav_online: "Online", nav_compare: "Porównanie", nav_pricing: "Cennik", nav_faq: "FAQ", nav_bot: "Bot online", nav_policy: "Prywatność", nav_terms: "Warunki", nav_policy_full: "Polityka prywatności", nav_terms_full: "Warunki korzystania", open_online: "Otwórz online", open_app: "Otwórz aplikację web", open_bot: "Otwórz bota",
            landing_badge: "Wielojęzyczny serwis AI: 20 języków, aplikacja web, logowanie Telegram i tłumacz", landing_title: "Ucz się języków online albo w Telegramie, w wygodnym dla Ciebie języku", landing_description: "NERIVA działa jako jeden serwis: aplikacja web, bot Telegram, logowanie przez Telegram, synchronizacja postępów oraz osobny wybór języka interfejsu i nauki.",
            terms_badge: "Zasady serwisu", terms_title: "Warunki korzystania", terms_description: "Te warunki regulują korzystanie ze strony, aplikacji web i bota Telegram NERIVA, w tym taryfy, płatności, dane osobowe, limity usługi i odpowiedzialność.",
            privacy_badge: "Ochrona danych", privacy_title: "Polityka prywatności", privacy_description: "Ta strona wyjaśnia, jakie dane może przetwarzać NERIVA i dlaczego są potrzebne do nauki, płatności, wsparcia, ochrony konta i praw użytkownika.",
            legal_pill_languages: "20 języków teraz", legal_pill_interface: "Wielojęzyczny interfejs", legal_pill_ai: "Praktyka AI", legal_pill_tariffs: "Free i Premium"
        },
        pt: {
            language_label: "Idioma", brand_subtitle: "idiomas online e no Telegram", brand_subtitle_short: "site e Telegram",
            nav_home: "Início", nav_features: "Recursos", nav_online: "Online", nav_compare: "Comparação", nav_pricing: "Preços", nav_faq: "FAQ", nav_bot: "Bot online", nav_policy: "Privacidade", nav_terms: "Termos", nav_policy_full: "Política de privacidade", nav_terms_full: "Termos de uso", open_online: "Abrir online", open_app: "Abrir app web", open_bot: "Abrir bot",
            landing_badge: "Serviço de IA multilíngue: 20 idiomas, app web, login Telegram e tradutor", landing_title: "Aprenda idiomas online ou no Telegram, no idioma mais confortável para você", landing_description: "NERIVA funciona como um único serviço: app web, bot Telegram, login via Telegram, progresso sincronizado e escolha separada do idioma da interface e de estudo.",
            terms_badge: "Regras do serviço", terms_title: "Termos de uso", terms_description: "Estes termos regulam o uso do site, app web e bot Telegram NERIVA, incluindo planos, pagamentos, dados pessoais, limites do serviço e responsabilidade.",
            privacy_badge: "Proteção de dados", privacy_title: "Política de privacidade", privacy_description: "Esta página explica quais dados o NERIVA pode processar e por que eles são necessários para estudo, pagamentos, suporte, proteção da conta e direitos do usuário.",
            legal_pill_languages: "20 idiomas agora", legal_pill_interface: "Interface multilíngue", legal_pill_ai: "Prática com IA", legal_pill_tariffs: "Free e Premium"
        },
        ro: {
            language_label: "Limbă", brand_subtitle: "limbi online și în Telegram", brand_subtitle_short: "site și Telegram",
            nav_home: "Acasă", nav_features: "Funcții", nav_online: "Online", nav_compare: "Comparație", nav_pricing: "Tarife", nav_faq: "FAQ", nav_bot: "Bot online", nav_policy: "Confidențialitate", nav_terms: "Termeni", nav_policy_full: "Politica de confidențialitate", nav_terms_full: "Termeni de utilizare", open_online: "Deschide online", open_app: "Deschide aplicația web", open_bot: "Deschide botul",
            landing_badge: "Serviciu AI multilingv: 20 limbi, aplicație web, login Telegram și traducător", landing_title: "Învață limbi online sau în Telegram, în limba care îți este comodă", landing_description: "NERIVA funcționează ca un singur serviciu: aplicație web, bot Telegram, login prin Telegram, progres sincronizat și alegere separată pentru limba interfeței și limba de studiu.",
            terms_badge: "Regulile serviciului", terms_title: "Termeni de utilizare", terms_description: "Acești termeni reglementează utilizarea site-ului, aplicației web și botului Telegram NERIVA, inclusiv tarife, plăți, date personale, limite ale serviciului și răspundere.",
            privacy_badge: "Protecția datelor", privacy_title: "Politica de confidențialitate", privacy_description: "Această pagină explică ce date poate prelucra NERIVA și de ce sunt necesare pentru învățare, plăți, suport, protecția contului și drepturile utilizatorului.",
            legal_pill_languages: "20 limbi acum", legal_pill_interface: "Interfață multilingvă", legal_pill_ai: "Practică AI", legal_pill_tariffs: "Free și Premium"
        },
        zh: {
            language_label: "语言", brand_subtitle: "在线和 Telegram 中的语言", brand_subtitle_short: "网站和 Telegram",
            nav_home: "首页", nav_features: "功能", nav_online: "在线", nav_compare: "对比", nav_pricing: "价格", nav_faq: "FAQ", nav_bot: "在线机器人", nav_policy: "隐私", nav_terms: "条款", nav_policy_full: "隐私政策", nav_terms_full: "使用条款", open_online: "在线打开", open_app: "打开 Web 应用", open_bot: "打开机器人",
            landing_badge: "多语言 AI 服务：20 种语言、Web 应用、Telegram 登录和翻译器", landing_title: "在线或在 Telegram 中用你熟悉的语言学习语言", landing_description: "NERIVA 是一个统一服务：Web 应用、Telegram 机器人、Telegram 登录、进度同步，并可分别选择界面语言和学习语言。",
            terms_badge: "服务规则", terms_title: "使用条款", terms_description: "本条款规范 NERIVA 网站、Web 应用和 Telegram 机器人的使用，包括套餐、付款、个人数据、服务限制和责任。",
            privacy_badge: "数据保护", privacy_title: "隐私政策", privacy_description: "本页面说明 NERIVA 可能处理哪些数据，以及这些数据为何用于学习、付款、支持、账户保护和用户权利。",
            legal_pill_languages: "目前 20 种语言", legal_pill_interface: "多语言界面", legal_pill_ai: "AI 练习", legal_pill_tariffs: "Free 和 Premium"
        },
        ja: {
            language_label: "言語", brand_subtitle: "オンラインと Telegram の言語", brand_subtitle_short: "サイトと Telegram",
            nav_home: "ホーム", nav_features: "機能", nav_online: "オンライン", nav_compare: "比較", nav_pricing: "料金", nav_faq: "FAQ", nav_bot: "オンラインボット", nav_policy: "プライバシー", nav_terms: "規約", nav_policy_full: "プライバシーポリシー", nav_terms_full: "利用規約", open_online: "オンラインで開く", open_app: "Webアプリを開く", open_bot: "ボットを開く",
            landing_badge: "20言語対応のAIサービス：Webアプリ、Telegramログイン、翻訳ツール", landing_title: "オンラインでも Telegram でも、使いやすい言語で学習", landing_description: "NERIVA は、Webアプリ、Telegramボット、Telegramログイン、進捗同期、インターフェース言語と学習言語の個別選択を備えた一つのサービスです。",
            terms_badge: "サービス規約", terms_title: "利用規約", terms_description: "本規約は、NERIVA のサイト、Webアプリ、Telegramボットの利用、料金、支払い、個人データ、サービス制限、責任について定めます。",
            privacy_badge: "データ保護", privacy_title: "プライバシーポリシー", privacy_description: "このページでは、NERIVA が処理する可能性のあるデータと、それが学習、支払い、サポート、アカウント保護、ユーザーの権利のために必要な理由を説明します。",
            legal_pill_languages: "現在20言語", legal_pill_interface: "多言語インターフェース", legal_pill_ai: "AI練習", legal_pill_tariffs: "Free と Premium"
        },
        ko: {
            language_label: "언어", brand_subtitle: "온라인 및 Telegram 언어", brand_subtitle_short: "사이트와 Telegram",
            nav_home: "홈", nav_features: "기능", nav_online: "온라인", nav_compare: "비교", nav_pricing: "요금", nav_faq: "FAQ", nav_bot: "온라인 봇", nav_policy: "개인정보", nav_terms: "약관", nav_policy_full: "개인정보 처리방침", nav_terms_full: "이용 약관", open_online: "온라인 열기", open_app: "웹 앱 열기", open_bot: "봇 열기",
            landing_badge: "20개 언어를 지원하는 다국어 AI 서비스: 웹 앱, Telegram 로그인, 번역기", landing_title: "온라인이나 Telegram에서 편한 언어로 학습하세요", landing_description: "NERIVA는 웹 앱, Telegram 봇, Telegram 로그인, 진도 동기화, 인터페이스 언어와 학습 언어의 별도 선택을 제공하는 하나의 서비스입니다.",
            terms_badge: "서비스 규칙", terms_title: "이용 약관", terms_description: "본 약관은 NERIVA 웹사이트, 웹 앱, Telegram 봇의 이용, 요금제, 결제, 개인정보, 서비스 제한 및 책임을 규정합니다.",
            privacy_badge: "데이터 보호", privacy_title: "개인정보 처리방침", privacy_description: "이 페이지는 NERIVA가 어떤 데이터를 처리할 수 있으며, 그 데이터가 학습, 결제, 지원, 계정 보호 및 사용자 권리를 위해 왜 필요한지 설명합니다.",
            legal_pill_languages: "현재 20개 언어", legal_pill_interface: "다국어 인터페이스", legal_pill_ai: "AI 연습", legal_pill_tariffs: "Free 및 Premium"
        }
    };

    ["tg", "uz", "tt", "hy", "kk", "ky", "ka"].forEach((code) => {
        copy[code] = {
            ...copy.en,
            language_label: code === "zh" ? "语言" : code === "ja" ? "言語" : code === "ko" ? "언어" : code === "hy" ? "Լեզու" : code === "ka" ? "ენა" : code === "kk" ? "Тіл" : code === "ky" ? "Тил" : code === "tg" ? "Забон" : code === "tt" ? "Тел" : "Til",
            brand_subtitle: code === "zh" ? "在线和 Telegram 中的语言" : code === "ja" ? "オンラインと Telegram の言語" : code === "ko" ? "온라인 및 Telegram 언어" : code === "hy" ? "լեզուներ առցանց և Telegram-ում" : code === "ka" ? "ენები ონლაინ და Telegram-ში" : code === "kk" ? "тілдер онлайн және Telegram-да" : code === "ky" ? "тилдер онлайн жана Telegramда" : code === "tg" ? "забонҳо онлайн ва дар Telegram" : code === "tt" ? "телләр онлайн һәм Telegramда" : "tillar onlayn va Telegramda",
            brand_subtitle_short: code === "zh" ? "网站和 Telegram" : code === "ja" ? "サイトと Telegram" : code === "ko" ? "사이트와 Telegram" : code === "hy" ? "կայք և Telegram" : code === "ka" ? "საიტი და Telegram" : code === "kk" ? "сайт және Telegram" : code === "ky" ? "сайт жана Telegram" : code === "tg" ? "сайт ва Telegram" : code === "tt" ? "сайт һәм Telegram" : "sayt va Telegram"
        };
    });

    const launchLandingCopy = {
        es: {
            landing_badge: "Tutor de IA en Telegram y web app: lecciones, voz, fotos, errores y vocabulario",
            landing_title: "Habla un nuevo idioma antes: tu tutor de IA ya vive en Telegram",
            landing_description: "NERIVA convierte tu mensajero en un entrenador personal de idiomas: lecciones cortas, práctica de conversación, revisión de errores, vocabulario, voz y traducción de texto desde fotos.",
            legal_pill_tariffs: "Free, Premium y Platinum"
        },
        de: {
            landing_badge: "KI-Tutor in Telegram und Web-App: Lektionen, Stimme, Fotos, Fehler und Vokabeln",
            landing_title: "Sprich schneller eine neue Sprache: dein KI-Tutor ist schon in Telegram",
            landing_description: "NERIVA macht deinen Messenger zum persönlichen Sprachcoach: kurze Lektionen, Dialogpraxis, Fehleranalyse, Vokabeltraining, Sprachverständnis und Fototext-Übersetzung.",
            legal_pill_tariffs: "Free, Premium und Platinum"
        },
        fr: {
            landing_badge: "Tuteur IA dans Telegram et l'app web : leçons, voix, photos, erreurs et vocabulaire",
            landing_title: "Parlez plus vite une nouvelle langue : votre tuteur IA est déjà dans Telegram",
            landing_description: "NERIVA transforme votre messagerie en coach linguistique personnel : leçons courtes, conversations, correction des erreurs, vocabulaire, voix et traduction de texte depuis les photos.",
            legal_pill_tariffs: "Free, Premium et Platinum"
        },
        it: {
            landing_badge: "Tutor IA in Telegram e web app: lezioni, voce, foto, errori e vocabolario",
            landing_title: "Parla prima una nuova lingua: il tuo tutor IA è già in Telegram",
            landing_description: "NERIVA trasforma il messenger in un coach linguistico personale: lezioni brevi, conversazioni, revisione degli errori, vocabolario, voce e traduzione del testo dalle foto.",
            legal_pill_tariffs: "Free, Premium e Platinum"
        },
        uk: {
            landing_badge: "AI-репетитор у Telegram і web app: уроки, голос, фото, помилки та словник",
            landing_title: "Заговоріть новою мовою швидше: AI-репетитор уже в Telegram",
            landing_description: "NERIVA перетворює месенджер на особистого мовного тренера: короткі уроки, діалоги, розбір помилок, тренування слів, розуміння голосу й переклад тексту з фото.",
            legal_pill_tariffs: "Free, Premium і Platinum"
        },
        pl: {
            landing_badge: "Tutor AI w Telegramie i aplikacji web: lekcje, głos, zdjęcia, błędy i słownictwo",
            landing_title: "Zacznij szybciej mówić w nowym języku: tutor AI jest już w Telegramie",
            landing_description: "NERIVA zmienia komunikator w osobistego trenera językowego: krótkie lekcje, dialogi, analiza błędów, słownictwo, głos i tłumaczenie tekstu ze zdjęć.",
            legal_pill_tariffs: "Free, Premium i Platinum"
        },
        pt: {
            landing_badge: "Tutor de IA no Telegram e web app: aulas, voz, fotos, erros e vocabulário",
            landing_title: "Fale um novo idioma mais rápido: seu tutor de IA já está no Telegram",
            landing_description: "NERIVA transforma seu mensageiro em um treinador pessoal de idiomas: aulas curtas, conversa, revisão de erros, vocabulário, voz e tradução de texto em fotos.",
            legal_pill_tariffs: "Free, Premium e Platinum"
        },
        ro: {
            landing_badge: "Tutor AI în Telegram și web app: lecții, voce, fotografii, greșeli și vocabular",
            landing_title: "Vorbește mai repede o limbă nouă: tutorul tău AI este deja în Telegram",
            landing_description: "NERIVA transformă messengerul într-un antrenor lingvistic personal: lecții scurte, dialoguri, analiza greșelilor, vocabular, voce și traducerea textului din fotografii.",
            legal_pill_tariffs: "Free, Premium și Platinum"
        },
        zh: {
            landing_badge: "Telegram 和 Web 应用中的 AI 导师：课程、语音、照片、错误和词汇",
            landing_title: "更快开口说新语言：你的 AI 导师已经在 Telegram 里",
            landing_description: "NERIVA 把常用聊天工具变成私人语言教练：短课、对话练习、错误复盘、词汇训练、语音理解和照片文字翻译。",
            legal_pill_tariffs: "Free、Premium 和 Platinum"
        },
        ja: {
            landing_badge: "Telegram と Web アプリのAIチューター：レッスン、音声、写真、ミス、語彙",
            landing_title: "新しい言語をもっと早く話そう：AIチューターはすでに Telegram の中に",
            landing_description: "NERIVA はいつものメッセンジャーを個人語学コーチに変えます。短いレッスン、会話練習、ミスの復習、語彙、音声理解、写真内テキスト翻訳に対応します。",
            legal_pill_tariffs: "Free、Premium、Platinum"
        },
        ko: {
            landing_badge: "Telegram과 웹 앱의 AI 튜터: 수업, 음성, 사진, 오류, 어휘",
            landing_title: "새 언어를 더 빠르게 말하세요: AI 튜터가 이미 Telegram 안에 있습니다",
            landing_description: "NERIVA는 메신저를 개인 언어 코치로 바꿉니다. 짧은 수업, 대화 연습, 오류 복습, 어휘 훈련, 음성 이해, 사진 속 텍스트 번역을 제공합니다.",
            legal_pill_tariffs: "Free, Premium 및 Platinum"
        },
        tg: {
            landing_badge: "Омӯзгори AI дар Telegram ва web app: дарсҳо, овоз, акс, хатоҳо ва луғат",
            landing_title: "Забони навро тезтар гап занед: омӯзгори AI аллакай дар Telegram аст",
            landing_description: "NERIVA паёмрасони шуморо ба мураббии шахсии забон табдил медиҳад: дарсҳои кӯтоҳ, муколама, таҳлили хатоҳо, луғат, фаҳмиши овоз ва тарҷумаи матн аз акс.",
            legal_pill_tariffs: "Free, Premium ва Platinum"
        },
        uz: {
            landing_badge: "Telegram va web app ichida AI repetitor: darslar, ovoz, foto, xatolar va lug'at",
            landing_title: "Yangi tilda tezroq gapiring: AI repetitor allaqachon Telegram ichida",
            landing_description: "NERIVA odatiy messenjerni shaxsiy til murabbiyiga aylantiradi: qisqa darslar, dialoglar, xatolar tahlili, lug'at mashqi, ovozni tushunish va fotodagi matn tarjimasi.",
            legal_pill_tariffs: "Free, Premium va Platinum"
        },
        tt: {
            landing_badge: "Telegram һәм web app эчендә AI-репетитор: дәресләр, тавыш, фото, хаталар һәм сүзлек",
            landing_title: "Яңа телдә тизрәк сөйләшә башлагыз: AI-репетитор инде Telegram эчендә",
            landing_description: "NERIVA гадәти мессенджерны шәхси тел тренерына әйләндерә: кыска дәресләр, диалоглар, хаталарны тикшерү, сүзлек, тавышны аңлау һәм фотодагы текстны тәрҗемә итү.",
            legal_pill_tariffs: "Free, Premium һәм Platinum"
        },
        hy: {
            landing_badge: "AI դասավանդող Telegram-ում և web app-ում՝ դասեր, ձայն, լուսանկար, սխալներ և բառապաշար",
            landing_title: "Ավելի արագ խոսեք նոր լեզվով. AI դասավանդողը արդեն Telegram-ում է",
            landing_description: "NERIVA-ը սովորական մեսենջերը դարձնում է անձնական լեզվի մարզիչ՝ կարճ դասեր, երկխոսություն, սխալների վերլուծություն, բառապաշար, ձայնի ընկալում և լուսանկարից տեքստի թարգմանություն:",
            legal_pill_tariffs: "Free, Premium և Platinum"
        },
        kk: {
            landing_badge: "Telegram және web app ішіндегі AI-репетитор: сабақтар, дауыс, фото, қателер және сөздік",
            landing_title: "Жаңа тілде тезірек сөйлеңіз: AI-репетитор Telegram ішінде",
            landing_description: "NERIVA әдеттегі мессенджерді жеке тіл жаттықтырушысына айналдырады: қысқа сабақтар, диалог, қателерді талдау, сөздік, дауысты түсіну және фотодағы мәтінді аудару.",
            legal_pill_tariffs: "Free, Premium және Platinum"
        },
        ky: {
            landing_badge: "Telegram жана web app ичиндеги AI-репетитор: сабактар, үн, фото, каталар жана сөздүк",
            landing_title: "Жаңы тилде тезирээк сүйлөңүз: AI-репетитор Telegram ичинде",
            landing_description: "NERIVA кадимки мессенжерди жеке тил машыктыруучусуна айлантат: кыска сабактар, диалог, каталарды талдоо, сөздүк, үндү түшүнүү жана фотодогу текстти которуу.",
            legal_pill_tariffs: "Free, Premium жана Platinum"
        },
        ka: {
            landing_badge: "AI რეპეტიტორი Telegram-სა და web app-ში: გაკვეთილები, ხმა, ფოტო, შეცდომები და ლექსიკა",
            landing_title: "უფრო სწრაფად ალაპარაკდით ახალ ენაზე: AI რეპეტიტორი უკვე Telegram-შია",
            landing_description: "NERIVA ჩვეულებრივ მესენჯერს პირად ენის მწვრთნელად აქცევს: მოკლე გაკვეთილები, დიალოგი, შეცდომების გარჩევა, ლექსიკა, ხმის გაგება და ფოტოდან ტექსტის თარგმნა.",
            legal_pill_tariffs: "Free, Premium და Platinum"
        }
    };

    Object.entries(launchLandingCopy).forEach(([code, values]) => {
        copy[code] = { ...copy[code], ...values };
    });

    const landingPricingLineCopy = {
        es: "La web app se sincroniza con Telegram, y Premium cuesta ahora 300 RUB, 150 Stars o unos 4.15 USDT al mes con el 70% de descuento de lanzamiento incluido. El descuento está disponible solo durante el primer mes de lanzamiento.",
        de: "Die Web-App synchronisiert sich mit Telegram, und Premium kostet jetzt 300 RUB, 150 Stars oder etwa 4.15 USDT pro Monat inklusive 70% Start-Rabatt. Der Rabatt gilt nur im ersten Launch-Monat.",
        fr: "L'app web se synchronise avec Telegram, et Premium coûte maintenant 300 RUB, 150 Stars ou environ 4.15 USDT par mois avec la remise de lancement de 70% incluse. La remise est disponible uniquement pendant le premier mois de lancement.",
        it: "La web app si sincronizza con Telegram, e Premium ora costa 300 RUB, 150 Stars o circa 4.15 USDT al mese con lo sconto lancio del 70% incluso. Lo sconto è disponibile solo durante il primo mese di lancio.",
        uk: "Веб-додаток синхронізується з Telegram, а Premium зараз коштує 300 RUB, 150 Stars або близько 4.15 USDT на місяць з урахуванням стартової знижки 70%. Знижка діє лише протягом першого місяця запуску.",
        pl: "Aplikacja web synchronizuje się z Telegramem, a Premium kosztuje teraz 300 RUB, 150 Stars lub około 4.15 USDT miesięcznie z uwzględnioną 70% zniżką startową. Zniżka obowiązuje tylko w pierwszym miesiącu startu.",
        pt: "O app web sincroniza com o Telegram, e o Premium agora custa 300 RUB, 150 Stars ou cerca de 4.15 USDT por mês com o desconto inicial de 70% incluído. O desconto vale apenas durante o primeiro mês de lançamento.",
        ro: "Aplicația web se sincronizează cu Telegram, iar Premium costă acum 300 RUB, 150 Stars sau aproximativ 4.15 USDT pe lună cu reducerea de lansare de 70% inclusă. Reducerea este disponibilă doar în prima lună de lansare.",
        zh: "Web 应用会与 Telegram 同步，Premium 现价为每月 300 卢布、150 Stars 或约 4.15 USDT，已包含 70% 上线折扣。折扣仅在上线第一个月有效。",
        ja: "Webアプリは Telegram と同期し、Premium は現在、70%のローンチ割引込みで月額300 RUB、150 Stars、または約4.15 USDTです。割引はローンチ最初の1か月だけ有効です。",
        ko: "웹 앱은 Telegram과 동기화되며, Premium은 출시 70% 할인이 적용되어 월 300 RUB, 150 Stars 또는 약 4.15 USDT입니다. 할인은 출시 첫 달에만 적용됩니다.",
        tg: "Web app бо Telegram ҳамоҳанг мешавад, Premium ҳоло бо тахфифи оғози 70% дар як моҳ 300 RUB, 150 Stars ё тақрибан 4.15 USDT аст. Тахфиф танҳо дар моҳи аввали оғоз амал мекунад.",
        uz: "Web app Telegram bilan sinxronlashadi, Premium esa 70% start chegirmasi bilan oyiga 300 RUB, 150 Stars yoki taxminan 4.15 USDT turadi. Chegirma faqat ishga tushirishning birinchi oyida amal qiladi.",
        tt: "Web app Telegram белән синхронлаша, ә Premium хәзер 70% старт ташламасы белән аена 300 RUB, 150 Stars яки якынча 4.15 USDT тора. Ташлама стартның беренче аенда гына гамәлдә.",
        hy: "Web app-ը համաժամացվում է Telegram-ի հետ, իսկ Premium-ը հիմա արժե ամսական 300 RUB, 150 Stars կամ մոտ 4.15 USDT՝ 70% մեկնարկային զեղչով: Զեղչը գործում է միայն մեկնարկի առաջին ամսում:",
        kk: "Web app Telegram-мен синхрондалады, ал Premium қазір 70% бастапқы жеңілдікпен айына 300 RUB, 150 Stars немесе шамамен 4.15 USDT тұрады. Жеңілдік іске қосылған алғашқы айда ғана қолданылады.",
        ky: "Web app Telegram менен синхрондолот, ал Premium азыр 70% старттык арзандатуу менен айына 300 RUB, 150 Stars же болжол менен 4.15 USDT турат. Арзандатуу башталган биринчи айда гана жарактуу.",
        ka: "Web app სინქრონდება Telegram-თან, ხოლო Premium ახლა 70% საწყისი ფასდაკლებით თვეში 300 RUB, 150 Stars ან დაახლოებით 4.15 USDT ღირს. ფასდაკლება მოქმედებს მხოლოდ გაშვების პირველ თვეში."
    };

    Object.entries(landingPricingLineCopy).forEach(([code, value]) => {
        copy[code] = { ...copy[code], landing_pricing_line: value };
    });

    const refreshedLandingCopy = {
        ru: {
            nav_reviews: "Отзывы",
            landing_badge: "AI-репетитор в Telegram и web app: уроки, аудирование, голосовая оценка, фото и словарь",
            landing_title: "AI-репетитор в Telegram и web app",
            landing_description: "NERIVA объединяет уроки, практику, аудирование, словарь, переводчик и оценку произношения. Пользователь занимается в Telegram или web app, получает структуру задания, голосовую обратную связь и понятный прогресс.",
            landing_pricing_line: "После голосового ответа пользователь видит понятный учебный разбор: примерную оценку, слабые слова, плавность речи и короткий совет, который можно сразу повторить вслух."
        },
        en: {
            nav_reviews: "Reviews",
            landing_badge: "AI tutor in Telegram and web app: lessons, listening, pronunciation scoring, photos and vocabulary",
            landing_title: "AI tutor in Telegram and web app",
            landing_description: "NERIVA combines lessons, practice, listening, vocabulary, translator tools and pronunciation scoring. Learn in Telegram or the web app, get structured tasks, voice feedback and visible progress.",
            landing_pricing_line: "After a voice answer, the learner sees clear coaching feedback: an approximate score, weak words, speech flow and a short tip they can repeat aloud right away."
        }
    };

    const navReviewsCopy = {
        es: "Reseñas", de: "Bewertungen", fr: "Avis", it: "Recensioni", zh: "评价", ja: "レビュー",
        ko: "후기", tg: "Шарҳҳо", uz: "Sharhlar", tt: "Фикерләр", hy: "Կարծիքներ",
        kk: "Пікірлер", ky: "Пикирлер", ka: "მიმოხილვები", uk: "Відгуки", pl: "Opinie",
        ro: "Recenzii", pt: "Avaliações"
    };

    languages.forEach(([code]) => {
        copy[code] = {
            ...copy[code],
            ...(code === "ru" ? refreshedLandingCopy.ru : refreshedLandingCopy.en),
            nav_reviews: code === "ru" ? refreshedLandingCopy.ru.nav_reviews : (navReviewsCopy[code] || refreshedLandingCopy.en.nav_reviews)
        };
    });

    const expandedLanguageCountCopy = {
        ru: "35 языков интерфейса", en: "35 interface languages", es: "35 idiomas de interfaz",
        de: "35 UI-Sprachen", fr: "35 langues d'interface", it: "35 lingue UI",
        zh: "35 种界面语言", ja: "35 UI 言語", ko: "35개 인터페이스 언어",
        tg: "35 забони интерфейс", uz: "35 interfeys tili", tt: "35 интерфейс теле",
        hy: "35 ինտերֆեյսի լեզու", kk: "35 интерфейс тілі", ky: "35 интерфейс тили",
        ka: "35 ინტერფეისის ენა", uk: "35 мов інтерфейсу", pl: "35 języków interfejsu",
        ro: "35 de limbi de interfață", pt: "35 idiomas de interface", ar: "35 لغة واجهة",
        bn: "৩৫টি ইন্টারফেস ভাষা", cs: "35 jazyků rozhraní", el: "35 γλώσσες διεπαφής",
        hi: "35 इंटरफ़ेस भाषाएं", hu: "35 felületi nyelv", id: "35 bahasa antarmuka",
        nl: "35 interfacetalen", sv: "35 gränssnittsspråk", ta: "35 இடைமுக மொழிகள்",
        te: "35 ఇంటర్‌ఫేస్ భాషలు", th: "35 ภาษาอินเทอร์เฟซ", tl: "35 wika ng interface",
        tr: "35 arayüz dili", vi: "35 ngôn ngữ giao diện"
    };

    languages.forEach(([code]) => {
        copy[code] = { ...copy.en, ...copy[code], legal_pill_languages: expandedLanguageCountCopy[code] || expandedLanguageCountCopy.en };
    });

    const legalCardCopy = {
        ru: { legal_card_site: "Сайт", legal_card_service: "Сервис", legal_card_bot: "Бот", legal_card_updated: "Обновлено", legal_updated_date: "22 мая 2026", legal_site_scope: "Сайты сервиса: neriva.ru и neriva.ru." },
        en: { legal_card_site: "Site", legal_card_service: "Service", legal_card_bot: "Bot", legal_card_updated: "Updated", legal_updated_date: "22 May 2026", legal_site_scope: "Service websites: neriva.ru and neriva.ru." },
        es: { legal_card_site: "Sitio", legal_card_service: "Servicio", legal_card_bot: "Bot", legal_card_updated: "Actualizado", legal_updated_date: "22 de mayo de 2026", legal_site_scope: "Sitios del servicio: neriva.ru y neriva.ru." },
        de: { legal_card_site: "Website", legal_card_service: "Dienst", legal_card_bot: "Bot", legal_card_updated: "Aktualisiert", legal_updated_date: "22. Mai 2026", legal_site_scope: "Websites des Dienstes: neriva.ru und neriva.ru." },
        fr: { legal_card_site: "Site", legal_card_service: "Service", legal_card_bot: "Bot", legal_card_updated: "Mis à jour", legal_updated_date: "22 mai 2026", legal_site_scope: "Sites du service : neriva.ru et neriva.ru." },
        it: { legal_card_site: "Sito", legal_card_service: "Servizio", legal_card_bot: "Bot", legal_card_updated: "Aggiornato", legal_updated_date: "22 maggio 2026", legal_site_scope: "Siti del servizio: neriva.ru e neriva.ru." },
        zh: { legal_card_site: "网站", legal_card_service: "服务", legal_card_bot: "机器人", legal_card_updated: "更新日期", legal_updated_date: "2026年5月22日", legal_site_scope: "服务网站：neriva.ru 和 neriva.ru。" },
        ja: { legal_card_site: "サイト", legal_card_service: "サービス", legal_card_bot: "ボット", legal_card_updated: "更新日", legal_updated_date: "2026年5月22日", legal_site_scope: "サービスのサイト：neriva.ru と neriva.ru。" },
        ko: { legal_card_site: "사이트", legal_card_service: "서비스", legal_card_bot: "봇", legal_card_updated: "업데이트", legal_updated_date: "2026년 5월 22일", legal_site_scope: "서비스 사이트: neriva.ru 및 neriva.ru." },
        tg: { legal_card_site: "Сайт", legal_card_service: "Сервис", legal_card_bot: "Бот", legal_card_updated: "Навсозӣ шуд", legal_updated_date: "22 майи 2026", legal_site_scope: "Сайтҳои сервис: neriva.ru ва neriva.ru." },
        uz: { legal_card_site: "Sayt", legal_card_service: "Servis", legal_card_bot: "Bot", legal_card_updated: "Yangilandi", legal_updated_date: "2026-yil 22-may", legal_site_scope: "Servis saytlari: neriva.ru va neriva.ru." },
        tt: { legal_card_site: "Сайт", legal_card_service: "Сервис", legal_card_bot: "Бот", legal_card_updated: "Яңартылды", legal_updated_date: "2026 елның 22 мае", legal_site_scope: "Сервис сайтлары: neriva.ru һәм neriva.ru." },
        hy: { legal_card_site: "Կայք", legal_card_service: "Ծառայություն", legal_card_bot: "Բոտ", legal_card_updated: "Թարմացվել է", legal_updated_date: "2026 մայիսի 22", legal_site_scope: "Ծառայության կայքերը՝ neriva.ru և neriva.ru։" },
        kk: { legal_card_site: "Сайт", legal_card_service: "Сервис", legal_card_bot: "Бот", legal_card_updated: "Жаңартылды", legal_updated_date: "2026 жылғы 22 мамыр", legal_site_scope: "Сервис сайттары: neriva.ru және neriva.ru." },
        ky: { legal_card_site: "Сайт", legal_card_service: "Сервис", legal_card_bot: "Бот", legal_card_updated: "Жаңыртылды", legal_updated_date: "2026-жылдын 22-майы", legal_site_scope: "Сервис сайттары: neriva.ru жана neriva.ru." },
        ka: { legal_card_site: "საიტი", legal_card_service: "სერვისი", legal_card_bot: "ბოტი", legal_card_updated: "განახლდა", legal_updated_date: "2026 წლის 22 მაისი", legal_site_scope: "სერვისის საიტები: neriva.ru და neriva.ru." },
        uk: { legal_card_site: "Сайт", legal_card_service: "Сервіс", legal_card_bot: "Бот", legal_card_updated: "Оновлено", legal_updated_date: "22 травня 2026", legal_site_scope: "Сайти сервісу: neriva.ru і neriva.ru." },
        pl: { legal_card_site: "Strona", legal_card_service: "Usługa", legal_card_bot: "Bot", legal_card_updated: "Zaktualizowano", legal_updated_date: "22 maja 2026", legal_site_scope: "Strony usługi: neriva.ru i neriva.ru." },
        ro: { legal_card_site: "Site", legal_card_service: "Serviciu", legal_card_bot: "Bot", legal_card_updated: "Actualizat", legal_updated_date: "22 mai 2026", legal_site_scope: "Site-urile serviciului: neriva.ru și neriva.ru." },
        pt: { legal_card_site: "Site", legal_card_service: "Serviço", legal_card_bot: "Bot", legal_card_updated: "Atualizado", legal_updated_date: "22 de maio de 2026", legal_site_scope: "Sites do serviço: neriva.ru e neriva.ru." }
    };

    Object.entries(legalCardCopy).forEach(([code, values]) => {
        copy[code] = { ...copy[code], ...values };
    });

    const phraseTranslations = window.poliglotPhraseTranslations || {};
    const phraseLanguageCodes = languages.map(([code]) => code).filter((code) => code !== "ru");

    function legalTranslations(en, overrides = {}) {
        const values = {};
        phraseLanguageCodes.forEach((code) => {
            values[code] = overrides[code] || en;
        });
        return values;
    }

    const englishSparkLiveCopy = {
        ru: {
            follow: "Соцсети NERIVA", socialBody: "Короткие уроки, обновления продукта и советы по обучению.", socialChannels: "Соцсети NERIVA",
            openYouTube: "Открыть NERIVA на YouTube", openInstagram: "Открыть NERIVA в Instagram", openTikTok: "Открыть NERIVA в TikTok",
            today: "Сегодня", cockpit: "Пульт AI-репетитора", featureDemo: "Демо функций", lesson: "Урок", dialogue: "Диалог", voice: "Голос", photo: "Фото",
            photoIntro: "Сфотографируйте меню, вывеску или задание и превратите это в учебный сценарий.", noPeanuts: "Без арахиса, пожалуйста. Насколько острое это блюдо?", photoHint: "Фото превращается в перевод, заметку и практическую подсказку.",
        },
        es: {
            follow: "Sigue a NERIVA", socialBody: "Lecciones breves, novedades del producto y consejos de aprendizaje.", socialChannels: "Redes sociales de NERIVA",
            openYouTube: "Abrir NERIVA en YouTube", openInstagram: "Abrir NERIVA en Instagram", openTikTok: "Abrir NERIVA en TikTok",
            today: "Hoy", cockpit: "Panel del tutor IA", featureDemo: "Demo de funciones", lesson: "Lección", dialogue: "Diálogo", voice: "Voz", photo: "Foto",
            photoIntro: "Fotografía un menú, cartel o tarea y conviértelo en un escenario de aprendizaje.", noPeanuts: "Sin cacahuetes, por favor. ¿Qué tan picante es este plato?", photoHint: "La foto se convierte en traducción, nota y práctica.",
        },
        de: {
            follow: "Folge NERIVA", socialBody: "Kurze Lektionen, Produktupdates und Lerntipps.", socialChannels: "Social Media von NERIVA",
            openYouTube: "NERIVA auf YouTube öffnen", openInstagram: "NERIVA auf Instagram öffnen", openTikTok: "NERIVA auf TikTok öffnen",
            today: "Heute", cockpit: "KI-Tutor-Cockpit", featureDemo: "Funktionsdemo", lesson: "Lektion", dialogue: "Dialog", voice: "Stimme", photo: "Foto",
            photoIntro: "Fotografiere ein Menü, Schild oder eine Aufgabe und mache daraus ein Lernszenario.", noPeanuts: "Bitte ohne Erdnüsse. Wie scharf ist dieses Gericht?", photoHint: "Das Foto wird zu Übersetzung, Notiz und Übungsimpuls.",
        },
        fr: {
            follow: "Suivre NERIVA", socialBody: "Courtes leçons, nouveautés produit et conseils d’apprentissage.", socialChannels: "Réseaux sociaux de NERIVA",
            openYouTube: "Ouvrir NERIVA sur YouTube", openInstagram: "Ouvrir NERIVA sur Instagram", openTikTok: "Ouvrir NERIVA sur TikTok",
            today: "Aujourd’hui", cockpit: "Poste du tuteur IA", featureDemo: "Démo des fonctions", lesson: "Leçon", dialogue: "Dialogue", voice: "Voix", photo: "Photo",
            photoIntro: "Photographiez un menu, un panneau ou une tâche et transformez-le en scénario d’apprentissage.", noPeanuts: "Sans cacahuètes, s’il vous plaît. Ce plat est-il très épicé ?", photoHint: "La photo devient une traduction, une note et une consigne de pratique.",
        },
        it: {
            follow: "Segui NERIVA", socialBody: "Lezioni brevi, aggiornamenti prodotto e consigli di studio.", socialChannels: "Social di NERIVA",
            openYouTube: "Apri NERIVA su YouTube", openInstagram: "Apri NERIVA su Instagram", openTikTok: "Apri NERIVA su TikTok",
            today: "Oggi", cockpit: "Console tutor IA", featureDemo: "Demo funzioni", lesson: "Lezione", dialogue: "Dialogo", voice: "Voce", photo: "Foto",
            photoIntro: "Fotografa un menu, un cartello o un compito e trasformalo in uno scenario di apprendimento.", noPeanuts: "Senza arachidi, per favore. Quanto è piccante questo piatto?", photoHint: "La foto diventa traduzione, nota e prompt di pratica.",
        },
        zh: {
            follow: "关注 NERIVA", socialBody: "短课、产品更新和学习技巧。", socialChannels: "NERIVA 社交媒体",
            openYouTube: "在 YouTube 打开 NERIVA", openInstagram: "在 Instagram 打开 NERIVA", openTikTok: "在 TikTok 打开 NERIVA",
            today: "今天", cockpit: "AI 导师驾驶舱", featureDemo: "功能演示", lesson: "课程", dialogue: "对话", voice: "语音", photo: "照片",
            photoIntro: "拍下菜单、标牌或任务，把它变成学习场景。", noPeanuts: "请不要花生。这道菜有多辣？", photoHint: "照片会变成翻译、笔记和练习提示。",
        },
        ja: {
            follow: "NERIVA をフォロー", socialBody: "短いレッスン、製品アップデート、学習のヒント。", socialChannels: "NERIVA のSNS",
            openYouTube: "YouTubeで NERIVA を開く", openInstagram: "Instagramで NERIVA を開く", openTikTok: "TikTokで NERIVA を開く",
            today: "今日", cockpit: "AIチューター操作画面", featureDemo: "機能デモ", lesson: "レッスン", dialogue: "会話", voice: "音声", photo: "写真",
            photoIntro: "メニュー、標識、課題を撮影して学習シナリオに変えます。", noPeanuts: "ピーナッツ抜きでお願いします。この料理はどのくらい辛いですか？", photoHint: "写真は翻訳、メモ、練習プロンプトになります。",
        },
        ko: {
            follow: "NERIVA 팔로우", socialBody: "짧은 레슨, 제품 업데이트, 학습 팁.", socialChannels: "NERIVA 소셜 채널",
            openYouTube: "YouTube에서 NERIVA 열기", openInstagram: "Instagram에서 NERIVA 열기", openTikTok: "TikTok에서 NERIVA 열기",
            today: "오늘", cockpit: "AI 튜터 조종석", featureDemo: "기능 데모", lesson: "레슨", dialogue: "대화", voice: "음성", photo: "사진",
            photoIntro: "메뉴, 표지판 또는 과제를 촬영해 학습 시나리오로 바꾸세요.", noPeanuts: "땅콩은 빼 주세요. 이 요리는 얼마나 매운가요?", photoHint: "사진은 번역, 메모, 연습 프롬프트가 됩니다.",
        },
        tg: {
            follow: "NERIVA-ро пайгирӣ кунед", socialBody: "Дарсҳои кӯтоҳ, навигариҳои маҳсулот ва маслиҳатҳои омӯзишӣ.", socialChannels: "Шабакаҳои иҷтимоии NERIVA",
            openYouTube: "NERIVA-ро дар YouTube кушоед", openInstagram: "NERIVA-ро дар Instagram кушоед", openTikTok: "NERIVA-ро дар TikTok кушоед",
            today: "Имрӯз", cockpit: "Панели AI-омӯзгор", featureDemo: "Намоиши функсияҳо", lesson: "Дарс", dialogue: "Муколама", voice: "Овоз", photo: "Акс",
            photoIntro: "Меню, лавҳа ё вазифаро акс гиред ва онро ба сенарияи омӯзишӣ табдил диҳед.", noPeanuts: "Лутфан, бе чормағзи заминӣ. Ин таом чӣ қадар тез аст?", photoHint: "Акс ба тарҷума, ёддошт ва супориши машқ табдил меёбад.",
        },
        uz: {
            follow: "NERIVA-ni kuzating", socialBody: "Qisqa darslar, mahsulot yangiliklari va o‘rganish maslahatlari.", socialChannels: "NERIVA ijtimoiy tarmoqlari",
            openYouTube: "NERIVA-ni YouTube’da ochish", openInstagram: "NERIVA-ni Instagram’da ochish", openTikTok: "NERIVA-ni TikTok’da ochish",
            today: "Bugun", cockpit: "AI repetitor paneli", featureDemo: "Funksiyalar demosı", lesson: "Dars", dialogue: "Dialog", voice: "Ovoz", photo: "Foto",
            photoIntro: "Menyu, belgi yoki topshiriqni suratga oling va uni o‘quv ssenariysiga aylantiring.", noPeanuts: "Yong‘oqsiz, iltimos. Bu taom qanchalik achchiq?", photoHint: "Foto tarjima, eslatma va mashq topshirig‘iga aylanadi.",
        },
        tt: {
            follow: "NERIVA-ны күзәтегез", socialBody: "Кыска дәресләр, продукт яңалыклары һәм уку киңәшләре.", socialChannels: "NERIVA социаль челтәрләре",
            openYouTube: "NERIVA-ны YouTube-та ачу", openInstagram: "NERIVA-ны Instagram-да ачу", openTikTok: "NERIVA-ны TikTok-та ачу",
            today: "Бүген", cockpit: "AI укытучы панели", featureDemo: "Функцияләр демосы", lesson: "Дәрес", dialogue: "Диалог", voice: "Тавыш", photo: "Фото",
            photoIntro: "Менюны, билгене яки биремне фотога төшереп, аны уку сценарие итегез.", noPeanuts: "Арахиссыз, зинһар. Бу ризык никадәр әче?", photoHint: "Фото тәрҗемәгә, язмага һәм күнегү биременә әйләнә.",
        },
        hy: {
            follow: "Հետևեք NERIVA-ին", socialBody: "Կարճ դասեր, արտադրանքի թարմացումներ և ուսուցման խորհուրդներ։", socialChannels: "NERIVA-ի սոցիալական ալիքները",
            openYouTube: "Բացել NERIVA-ը YouTube-ում", openInstagram: "Բացել NERIVA-ը Instagram-ում", openTikTok: "Բացել NERIVA-ը TikTok-ում",
            today: "Այսօր", cockpit: "AI ուսուցչի վահանակ", featureDemo: "Գործառույթների ցուցադրում", lesson: "Դաս", dialogue: "Երկխոսություն", voice: "Ձայն", photo: "Լուսանկար",
            photoIntro: "Լուսանկարեք մենյուն, ցուցանակը կամ առաջադրանքը և դարձրեք ուսումնական սցենար։", noPeanuts: "Առանց գետնանուշի, խնդրում եմ։ Այս ուտեստը որքան կծու է՞։", photoHint: "Լուսանկարը դառնում է թարգմանություն, նշում և վարժության հուշում։",
        },
        kk: {
            follow: "NERIVA-ға жазылыңыз", socialBody: "Қысқа сабақтар, өнім жаңалықтары және оқу кеңестері.", socialChannels: "NERIVA әлеуметтік желілері",
            openYouTube: "NERIVA-ды YouTube-та ашу", openInstagram: "NERIVA-ды Instagram-да ашу", openTikTok: "NERIVA-ды TikTok-та ашу",
            today: "Бүгін", cockpit: "AI тәлімгер панелі", featureDemo: "Функциялар демосы", lesson: "Сабақ", dialogue: "Диалог", voice: "Дауыс", photo: "Фото",
            photoIntro: "Мәзірді, белгіні немесе тапсырманы суретке түсіріп, оны оқу сценарийіне айналдырыңыз.", noPeanuts: "Жержаңғақсыз, өтінемін. Бұл тағам қаншалықты ащы?", photoHint: "Фото аудармаға, жазбаға және практика нұсқауына айналады.",
        },
        ky: {
            follow: "NERIVA'га жазылыңыз", socialBody: "Кыска сабактар, продукт жаңылыктары жана окуу кеңештери.", socialChannels: "NERIVA социалдык тармактары",
            openYouTube: "NERIVA'ды YouTube'та ачуу", openInstagram: "NERIVA'ды Instagram'да ачуу", openTikTok: "NERIVA'ды TikTok'та ачуу",
            today: "Бүгүн", cockpit: "AI мугалим панели", featureDemo: "Функциялар демосу", lesson: "Сабак", dialogue: "Диалог", voice: "Үн", photo: "Сүрөт",
            photoIntro: "Менюну, белгини же тапшырманы сүрөткө тартып, аны окуу сценарийине айлантыңыз.", noPeanuts: "Жер жаңгаксыз, сураныч. Бул тамак канчалык ачуу?", photoHint: "Сүрөт котормого, жазууга жана машыгуу тапшырмасына айланат.",
        },
        ka: {
            follow: "გამოიწერეთ NERIVA", socialBody: "მოკლე გაკვეთილები, პროდუქტის განახლებები და სწავლის რჩევები.", socialChannels: "NERIVA-ის სოციალური არხები",
            openYouTube: "NERIVA-ის გახსნა YouTube-ზე", openInstagram: "NERIVA-ის გახსნა Instagram-ზე", openTikTok: "NERIVA-ის გახსნა TikTok-ზე",
            today: "დღეს", cockpit: "AI მასწავლებლის პანელი", featureDemo: "ფუნქციების დემო", lesson: "გაკვეთილი", dialogue: "დიალოგი", voice: "ხმა", photo: "ფოტო",
            photoIntro: "გადაუღეთ ფოტო მენიუს, ნიშანს ან დავალებას და აქციეთ სასწავლო სცენარად.", noPeanuts: "არაქისის გარეშე, გთხოვთ. რამდენად ცხარეა ეს კერძი?", photoHint: "ფოტო გადაიქცევა თარგმანად, ჩანაწერად და სავარჯიშო მინიშნებად.",
        },
        uk: {
            follow: "Стежте за NERIVA", socialBody: "Короткі уроки, оновлення продукту й поради для навчання.", socialChannels: "Соцмережі NERIVA",
            openYouTube: "Відкрити NERIVA на YouTube", openInstagram: "Відкрити NERIVA в Instagram", openTikTok: "Відкрити NERIVA у TikTok",
            today: "Сьогодні", cockpit: "Панель AI-репетитора", featureDemo: "Демо функцій", lesson: "Урок", dialogue: "Діалог", voice: "Голос", photo: "Фото",
            photoIntro: "Сфотографуйте меню, вивіску або завдання й перетворіть це на навчальний сценарій.", noPeanuts: "Без арахісу, будь ласка. Наскільки гостра ця страва?", photoHint: "Фото стає перекладом, нотаткою та підказкою для практики.",
        },
        pl: {
            follow: "Obserwuj NERIVA", socialBody: "Krótkie lekcje, aktualizacje produktu i wskazówki do nauki.", socialChannels: "Kanały społecznościowe NERIVA",
            openYouTube: "Otwórz NERIVA na YouTube", openInstagram: "Otwórz NERIVA na Instagramie", openTikTok: "Otwórz NERIVA na TikToku",
            today: "Dziś", cockpit: "Panel tutora AI", featureDemo: "Demo funkcji", lesson: "Lekcja", dialogue: "Dialog", voice: "Głos", photo: "Zdjęcie",
            photoIntro: "Zrób zdjęcie menu, znaku lub zadania i zamień je w scenariusz nauki.", noPeanuts: "Bez orzeszków ziemnych, proszę. Jak ostre jest to danie?", photoHint: "Zdjęcie staje się tłumaczeniem, notatką i zadaniem do ćwiczenia.",
        },
        ro: {
            follow: "Urmărește NERIVA", socialBody: "Lecții scurte, actualizări de produs și sfaturi de învățare.", socialChannels: "Rețelele sociale NERIVA",
            openYouTube: "Deschide NERIVA pe YouTube", openInstagram: "Deschide NERIVA pe Instagram", openTikTok: "Deschide NERIVA pe TikTok",
            today: "Astăzi", cockpit: "Panou tutore AI", featureDemo: "Demo funcții", lesson: "Lecție", dialogue: "Dialog", voice: "Voce", photo: "Foto",
            photoIntro: "Fotografiază un meniu, indicator sau exercițiu și transformă-l într-un scenariu de învățare.", noPeanuts: "Fără arahide, vă rog. Cât de picant este acest fel?", photoHint: "Fotografia devine traducere, notă și prompt de practică.",
        },
        pt: {
            follow: "Segue o NERIVA", socialBody: "Lições curtas, novidades do produto e dicas de aprendizagem.", socialChannels: "Redes sociais do NERIVA",
            openYouTube: "Abrir o NERIVA no YouTube", openInstagram: "Abrir o NERIVA no Instagram", openTikTok: "Abrir o NERIVA no TikTok",
            today: "Hoje", cockpit: "Painel do tutor IA", featureDemo: "Demonstração de funções", lesson: "Lição", dialogue: "Diálogo", voice: "Voz", photo: "Foto",
            photoIntro: "Fotografa um menu, placa ou tarefa e transforma isso num cenário de aprendizagem.", noPeanuts: "Sem amendoins, por favor. Quão picante é este prato?", photoHint: "A foto vira tradução, nota e prompt de prática.",
        },
        ar: {
            follow: "تابع NERIVA", socialBody: "دروس قصيرة وتحديثات المنتج ونصائح للتعلم.", socialChannels: "قنوات NERIVA الاجتماعية",
            openYouTube: "فتح NERIVA على YouTube", openInstagram: "فتح NERIVA على Instagram", openTikTok: "فتح NERIVA على TikTok",
            today: "اليوم", cockpit: "لوحة المدرّس الذكي", featureDemo: "عرض الميزات", lesson: "درس", dialogue: "حوار", voice: "صوت", photo: "صورة",
            photoIntro: "صوّر قائمة أو لافتة أو مهمة وحوّلها إلى سيناريو تعلم.", noPeanuts: "بدون فول سوداني من فضلك. ما مدى حرارة هذا الطبق؟", photoHint: "تتحول الصورة إلى ترجمة وملاحظة ومهمة تدريب.",
        },
        bn: {
            follow: "NERIVA অনুসরণ করুন", socialBody: "ছোট পাঠ, পণ্যের আপডেট এবং শেখার টিপস।", socialChannels: "NERIVA সামাজিক চ্যানেল",
            openYouTube: "YouTube-এ NERIVA খুলুন", openInstagram: "Instagram-এ NERIVA খুলুন", openTikTok: "TikTok-এ NERIVA খুলুন",
            today: "আজ", cockpit: "AI টিউটর প্যানেল", featureDemo: "ফিচার ডেমো", lesson: "পাঠ", dialogue: "সংলাপ", voice: "ভয়েস", photo: "ছবি",
            photoIntro: "মেনু, সাইন বা কাজের ছবি তুলে সেটিকে শেখার দৃশ্যে বদলে দিন।", noPeanuts: "বাদাম ছাড়া, দয়া করে। এই খাবারটি কতটা ঝাল?", photoHint: "ছবিটি অনুবাদ, নোট এবং অনুশীলন প্রম্পটে বদলে যায়।",
        },
        cs: {
            follow: "Sledujte NERIVA", socialBody: "Krátké lekce, novinky produktu a tipy k učení.", socialChannels: "Sociální sítě NERIVA",
            openYouTube: "Otevřít NERIVA na YouTube", openInstagram: "Otevřít NERIVA na Instagramu", openTikTok: "Otevřít NERIVA na TikToku",
            today: "Dnes", cockpit: "Panel AI tutora", featureDemo: "Ukázka funkcí", lesson: "Lekce", dialogue: "Dialog", voice: "Hlas", photo: "Foto",
            photoIntro: "Vyfoťte menu, ceduli nebo úkol a proměňte je ve výukový scénář.", noPeanuts: "Bez arašídů, prosím. Jak pálivé je toto jídlo?", photoHint: "Fotka se změní v překlad, poznámku a cvičný úkol.",
        },
        el: {
            follow: "Ακολουθήστε το NERIVA", socialBody: "Σύντομα μαθήματα, ενημερώσεις προϊόντος και συμβουλές μάθησης.", socialChannels: "Κοινωνικά κανάλια NERIVA",
            openYouTube: "Άνοιγμα NERIVA στο YouTube", openInstagram: "Άνοιγμα NERIVA στο Instagram", openTikTok: "Άνοιγμα NERIVA στο TikTok",
            today: "Σήμερα", cockpit: "Πίνακας AI δασκάλου", featureDemo: "Επίδειξη λειτουργιών", lesson: "Μάθημα", dialogue: "Διάλογος", voice: "Φωνή", photo: "Φωτογραφία",
            photoIntro: "Φωτογραφίστε μενού, πινακίδα ή εργασία και μετατρέψτε το σε σενάριο μάθησης.", noPeanuts: "Χωρίς φιστίκια, παρακαλώ. Πόσο καυτερό είναι αυτό το πιάτο;", photoHint: "Η φωτογραφία γίνεται μετάφραση, σημείωση και άσκηση.",
        },
        hi: {
            follow: "NERIVA को फ़ॉलो करें", socialBody: "छोटे पाठ, उत्पाद अपडेट और सीखने की सलाह।", socialChannels: "NERIVA सोशल चैनल",
            openYouTube: "YouTube पर NERIVA खोलें", openInstagram: "Instagram पर NERIVA खोलें", openTikTok: "TikTok पर NERIVA खोलें",
            today: "आज", cockpit: "AI ट्यूटर पैनल", featureDemo: "फ़ीचर डेमो", lesson: "पाठ", dialogue: "संवाद", voice: "आवाज़", photo: "फ़ोटो",
            photoIntro: "मेन्यू, संकेत या कार्य की फोटो लें और उसे सीखने के परिदृश्य में बदलें।", noPeanuts: "कृपया मूंगफली नहीं। यह डिश कितनी तीखी है?", photoHint: "फ़ोटो अनुवाद, नोट और अभ्यास संकेत बन जाती है।",
        },
        hu: {
            follow: "Kövesd a NERIVA-t", socialBody: "Rövid leckék, termékfrissítések és tanulási tippek.", socialChannels: "NERIVA közösségi csatornák",
            openYouTube: "NERIVA megnyitása YouTube-on", openInstagram: "NERIVA megnyitása Instagramon", openTikTok: "NERIVA megnyitása TikTokon",
            today: "Ma", cockpit: "AI tutor vezérlőpult", featureDemo: "Funkcióbemutató", lesson: "Lecke", dialogue: "Párbeszéd", voice: "Hang", photo: "Fotó",
            photoIntro: "Fotózz le egy menüt, táblát vagy feladatot, és alakítsd tanulási helyzetté.", noPeanuts: "Mogyoró nélkül kérem. Mennyire csípős ez az étel?", photoHint: "A fotóból fordítás, jegyzet és gyakorló feladat lesz.",
        },
        id: {
            follow: "Ikuti NERIVA", socialBody: "Pelajaran singkat, pembaruan produk, dan tips belajar.", socialChannels: "Kanal sosial NERIVA",
            openYouTube: "Buka NERIVA di YouTube", openInstagram: "Buka NERIVA di Instagram", openTikTok: "Buka NERIVA di TikTok",
            today: "Hari ini", cockpit: "Panel tutor AI", featureDemo: "Demo fitur", lesson: "Pelajaran", dialogue: "Dialog", voice: "Suara", photo: "Foto",
            photoIntro: "Foto menu, tanda, atau tugas lalu ubah menjadi skenario belajar.", noPeanuts: "Tanpa kacang, ya. Seberapa pedas hidangan ini?", photoHint: "Foto menjadi terjemahan, catatan, dan latihan.",
        },
        nl: {
            follow: "Volg NERIVA", socialBody: "Korte lessen, productupdates en leertips.", socialChannels: "Sociale kanalen van NERIVA",
            openYouTube: "NERIVA openen op YouTube", openInstagram: "NERIVA openen op Instagram", openTikTok: "NERIVA openen op TikTok",
            today: "Vandaag", cockpit: "AI-tutor cockpit", featureDemo: "Functiedemo", lesson: "Les", dialogue: "Dialoog", voice: "Stem", photo: "Foto",
            photoIntro: "Fotografeer een menu, bord of taak en maak er een leerscenario van.", noPeanuts: "Geen pinda's, alstublieft. Hoe pittig is dit gerecht?", photoHint: "De foto wordt een vertaling, notitie en oefenprompt.",
        },
        sv: {
            follow: "Följ NERIVA", socialBody: "Korta lektioner, produktnyheter och lärandetips.", socialChannels: "NERIVA i sociala medier",
            openYouTube: "Öppna NERIVA på YouTube", openInstagram: "Öppna NERIVA på Instagram", openTikTok: "Öppna NERIVA på TikTok",
            today: "I dag", cockpit: "AI-lärarens panel", featureDemo: "Funktionsdemo", lesson: "Lektion", dialogue: "Dialog", voice: "Röst", photo: "Foto",
            photoIntro: "Fotografera en meny, skylt eller uppgift och gör den till ett lärscenario.", noPeanuts: "Inga jordnötter, tack. Hur stark är den här rätten?", photoHint: "Fotot blir en översättning, anteckning och övningsprompt.",
        },
        ta: {
            follow: "NERIVA-ஐ பின்தொடருங்கள்", socialBody: "குறுகிய பாடங்கள், தயாரிப்பு புதுப்பிப்புகள், கற்றல் குறிப்புகள்.", socialChannels: "NERIVA சமூக சேனல்கள்",
            openYouTube: "YouTube-ல் NERIVA-ஐ திறக்கவும்", openInstagram: "Instagram-ல் NERIVA-ஐ திறக்கவும்", openTikTok: "TikTok-ல் NERIVA-ஐ திறக்கவும்",
            today: "இன்று", cockpit: "AI ஆசிரியர் பலகை", featureDemo: "அம்ச விளக்கம்", lesson: "பாடம்", dialogue: "உரையாடல்", voice: "குரல்", photo: "புகைப்படம்",
            photoIntro: "மெனு, பலகை அல்லது பணியைப் புகைப்படம் எடுத்து அதை கற்றல் சூழலாக மாற்றுங்கள்.", noPeanuts: "வேர்க்கடலை வேண்டாம். இந்த உணவு எவ்வளவு காரம்?", photoHint: "புகைப்படம் மொழிபெயர்ப்பு, குறிப்பு, பயிற்சி தூண்டுதலாக மாறும்.",
        },
        te: {
            follow: "NERIVAని అనుసరించండి", socialBody: "చిన్న పాఠాలు, ఉత్పత్తి నవీకరణలు, నేర్చుకునే చిట్కాలు.", socialChannels: "NERIVA సామాజిక ఛానెల్లు",
            openYouTube: "YouTubeలో NERIVA తెరవండి", openInstagram: "Instagramలో NERIVA తెరవండి", openTikTok: "TikTokలో NERIVA తెరవండి",
            today: "ఈ రోజు", cockpit: "AI ట్యూటర్ ప్యానెల్", featureDemo: "ఫీచర్ డెమో", lesson: "పాఠం", dialogue: "సంభాషణ", voice: "వాయిస్", photo: "ఫోటో",
            photoIntro: "మెనూ, బోర్డు లేదా పనిని ఫోటో తీసి దాన్ని నేర్చుకునే సన్నివేశంగా మార్చండి.", noPeanuts: "పల్లీలు వద్దు, దయచేసి. ఈ వంటకం ఎంత కారం?", photoHint: "ఫోటో అనువాదం, గమనిక, సాధన సూచనగా మారుతుంది.",
        },
        th: {
            follow: "ติดตาม NERIVA", socialBody: "บทเรียนสั้น อัปเดตผลิตภัณฑ์ และเคล็ดลับการเรียน.", socialChannels: "ช่องทางโซเชียลของ NERIVA",
            openYouTube: "เปิด NERIVA บน YouTube", openInstagram: "เปิด NERIVA บน Instagram", openTikTok: "เปิด NERIVA บน TikTok",
            today: "วันนี้", cockpit: "แผงติวเตอร์ AI", featureDemo: "เดโมฟีเจอร์", lesson: "บทเรียน", dialogue: "บทสนทนา", voice: "เสียง", photo: "รูปภาพ",
            photoIntro: "ถ่ายรูปเมนู ป้าย หรืองาน แล้วเปลี่ยนเป็นสถานการณ์การเรียนรู้.", noPeanuts: "ไม่ใส่ถั่วลิสงครับ/ค่ะ จานนี้เผ็ดแค่ไหน?", photoHint: "รูปภาพจะกลายเป็นคำแปล โน้ต และโจทย์ฝึก.",
        },
        tl: {
            follow: "Sundan ang NERIVA", socialBody: "Maiikling aralin, update ng produkto, at tips sa pag-aaral.", socialChannels: "Mga social channel ng NERIVA",
            openYouTube: "Buksan ang NERIVA sa YouTube", openInstagram: "Buksan ang NERIVA sa Instagram", openTikTok: "Buksan ang NERIVA sa TikTok",
            today: "Ngayon", cockpit: "Panel ng AI tutor", featureDemo: "Demo ng mga feature", lesson: "Aralin", dialogue: "Dayalogo", voice: "Boses", photo: "Larawan",
            photoIntro: "Kunan ng larawan ang menu, karatula, o gawain at gawing senaryo sa pag-aaral.", noPeanuts: "Walang mani, pakiusap. Gaano kaanghang ang pagkaing ito?", photoHint: "Ang larawan ay nagiging salin, tala, at prompt sa pagsasanay.",
        },
        tr: {
            follow: "NERIVA'ı takip et", socialBody: "Kısa dersler, ürün güncellemeleri ve öğrenme ipuçları.", socialChannels: "NERIVA sosyal kanalları",
            openYouTube: "NERIVA'ı YouTube'da aç", openInstagram: "NERIVA'ı Instagram'da aç", openTikTok: "NERIVA'ı TikTok'ta aç",
            today: "Bugün", cockpit: "AI öğretmen paneli", featureDemo: "Özellik demosu", lesson: "Ders", dialogue: "Diyalog", voice: "Ses", photo: "Fotoğraf",
            photoIntro: "Bir menüyü, tabelayı veya görevi fotoğraflayıp öğrenme senaryosuna dönüştür.", noPeanuts: "Fıstık olmasın lütfen. Bu yemek ne kadar acı?", photoHint: "Fotoğraf çeviri, not ve pratik istemine dönüşür.",
        },
        vi: {
            follow: "Theo dõi NERIVA", socialBody: "Bài học ngắn, cập nhật sản phẩm và mẹo học tập.", socialChannels: "Kênh mạng xã hội của NERIVA",
            openYouTube: "Mở NERIVA trên YouTube", openInstagram: "Mở NERIVA trên Instagram", openTikTok: "Mở NERIVA trên TikTok",
            today: "Hôm nay", cockpit: "Bảng điều khiển gia sư AI", featureDemo: "Demo tính năng", lesson: "Bài học", dialogue: "Đối thoại", voice: "Giọng nói", photo: "Ảnh",
            photoIntro: "Chụp thực đơn, biển báo hoặc bài tập và biến nó thành tình huống học.", noPeanuts: "Không đậu phộng, làm ơn. Món này cay thế nào?", photoHint: "Ảnh trở thành bản dịch, ghi chú và gợi ý luyện tập.",
        },
    };

    function addEnglishSparkPhrase(source, key) {
        const values = { en: source };
        Object.entries(englishSparkLiveCopy).forEach(([code, phrases]) => {
            values[code] = phrases[key] || source;
        });
        const normalized = normalizePhrase(source);
        phraseTranslations[normalized] = { ...(phraseTranslations[normalized] || {}), ...values };
    }

    [
        ["Follow NERIVA", "follow"],
        ["Short lessons, product updates, and learning tips.", "socialBody"],
        ["NERIVA social channels", "socialChannels"],
        ["Open NERIVA on YouTube", "openYouTube"],
        ["Open NERIVA on Instagram", "openInstagram"],
        ["Open NERIVA on TikTok", "openTikTok"],
        ["Social", "socialChannels"],
        ["Соцсети", "socialChannels"],
        ["Today", "today"],
        ["AI Tutor Cockpit", "cockpit"],
        ["Feature demo", "featureDemo"],
        ["Lesson", "lesson"],
        ["Dialogue", "dialogue"],
        ["Voice", "voice"],
        ["Photo", "photo"],
        ["Photograph a menu, sign, or task and turn it into a learning scenario.", "photoIntro"],
        ["No peanuts, please. How spicy is this dish?", "noPeanuts"],
        ["The photo becomes translation, a note, and a practice prompt.", "photoHint"],
    ].forEach(([source, key]) => addEnglishSparkPhrase(source, key));

    const legalPhraseTranslations = {
        "4. Free, Premium и Platinum": legalTranslations("4. Free, Premium and Platinum", {
            es: "4. Free, Premium y Platinum", de: "4. Free, Premium und Platinum", fr: "4. Free, Premium et Platinum",
            it: "4. Free, Premium e Platinum", uk: "4. Free, Premium і Platinum", pl: "4. Free, Premium i Platinum",
            pt: "4. Free, Premium e Platinum", ro: "4. Free, Premium și Platinum", zh: "4. Free、Premium 和 Platinum",
            ja: "4. Free、Premium、Platinum", ko: "4. Free, Premium 및 Platinum"
        }),
        "4. Тарифы Free, Premium и Platinum": legalTranslations("4. Free, Premium and Platinum plans", {
            es: "4. Tarifas Free, Premium y Platinum", de: "4. Tarife Free, Premium und Platinum", fr: "4. Offres Free, Premium et Platinum",
            it: "4. Piani Free, Premium e Platinum", uk: "4. Тарифи Free, Premium і Platinum", pl: "4. Taryfy Free, Premium i Platinum",
            pt: "4. Planos Free, Premium e Platinum", ro: "4. Tarife Free, Premium și Platinum", zh: "4. Free、Premium 和 Platinum 套餐",
            ja: "4. Free、Premium、Platinum プラン", ko: "4. Free, Premium 및 Platinum 요금제"
        }),
        "5. Оплата и Stars": legalTranslations("5. Payments and Stars", {
            es: "5. Pagos y Stars", de: "5. Zahlung und Stars", fr: "5. Paiement et Stars", it: "5. Pagamenti e Stars",
            uk: "5. Оплата і Stars", pl: "5. Płatności i Stars", pt: "5. Pagamentos e Stars", ro: "5. Plăți și Stars",
            zh: "5. 付款和 Stars", ja: "5. 支払いと Stars", ko: "5. 결제 및 Stars"
        }),
        "5. Оплата, Telegram Stars, TON, USDT и срок доступа": legalTranslations("5. Payment, Telegram Stars, TON, USDT and access period", {
            es: "5. Pago, Telegram Stars, TON, USDT y periodo de acceso", de: "5. Zahlung, Telegram Stars, TON, USDT und Zugriffszeitraum",
            fr: "5. Paiement, Telegram Stars, TON, USDT et durée d'accès", it: "5. Pagamento, Telegram Stars, TON, USDT e periodo di accesso",
            uk: "5. Оплата, Telegram Stars, TON, USDT і строк доступу", pl: "5. Płatność, Telegram Stars, TON, USDT i okres dostępu",
            pt: "5. Pagamento, Telegram Stars, TON, USDT e período de acesso", ro: "5. Plată, Telegram Stars, TON, USDT și perioada de acces",
            zh: "5. 付款、Telegram Stars、TON、USDT 和访问期限", ja: "5. 支払い、Telegram Stars、TON、USDT、利用期間",
            ko: "5. 결제, Telegram Stars, TON, USDT 및 이용 기간"
        }),
        "6. Оплаты и Premium": legalTranslations("6. Payments and Premium", {
            es: "6. Pagos y Premium", de: "6. Zahlungen und Premium", fr: "6. Paiements et Premium", it: "6. Pagamenti e Premium",
            uk: "6. Оплати і Premium", pl: "6. Płatności i Premium", pt: "6. Pagamentos e Premium", ro: "6. Plăți și Premium",
            zh: "6. 付款和 Premium", ja: "6. 支払いと Premium", ko: "6. 결제 및 Premium"
        }),
        "6. Оплаты, Telegram Stars и платный доступ": legalTranslations("6. Payments, Telegram Stars and paid access", {
            es: "6. Pagos, Telegram Stars y acceso de pago", de: "6. Zahlungen, Telegram Stars und bezahlter Zugang",
            fr: "6. Paiements, Telegram Stars et accès payant", it: "6. Pagamenti, Telegram Stars e accesso a pagamento",
            uk: "6. Оплати, Telegram Stars і платний доступ", pl: "6. Płatności, Telegram Stars i płatny dostęp",
            pt: "6. Pagamentos, Telegram Stars e acesso pago", ro: "6. Plăți, Telegram Stars și acces plătit",
            zh: "6. 付款、Telegram Stars 和付费访问", ja: "6. 支払い、Telegram Stars、有料アクセス",
            ko: "6. 결제, Telegram Stars 및 유료 이용"
        }),
        "Бесплатно": legalTranslations("Free", { es: "Gratis", de: "Kostenlos", fr: "Gratuit", it: "Gratis", uk: "Безкоштовно", pl: "Bezpłatnie", pt: "Grátis", ro: "Gratuit", zh: "免费", ja: "無料", ko: "무료" }),
        "70% скидка": legalTranslations("70% discount", { es: "70% de descuento", de: "70% Rabatt", fr: "70% de remise", it: "70% di sconto", uk: "знижка 70%", pl: "70% zniżki", pt: "70% de desconto", ro: "reducere 70%", zh: "70% 折扣", ja: "70%割引", ko: "70% 할인" }),
        "Лучший старт": legalTranslations("Best start", { es: "Mejor inicio", de: "Bester Start", fr: "Meilleur départ", it: "Miglior inizio", uk: "Найкращий старт", pl: "Najlepszy start", pt: "Melhor começo", ro: "Cel mai bun început", zh: "最佳开始", ja: "最高のスタート", ko: "최고의 시작" }),
        "Для интенсивной практики": legalTranslations("For intensive practice", { es: "Para práctica intensiva", de: "Für intensives Üben", fr: "Pour une pratique intensive", it: "Per pratica intensiva", uk: "Для інтенсивної практики", pl: "Do intensywnej praktyki", pt: "Para prática intensiva", ro: "Pentru practică intensivă", zh: "适合高强度练习", ja: "集中練習向け", ko: "집중 연습용" }),
        "5 уроков в день": legalTranslations("5 lessons per day", { es: "5 lecciones al día", de: "5 Lektionen pro Tag", fr: "5 leçons par jour", it: "5 lezioni al giorno", uk: "5 уроків на день", pl: "5 lekcji dziennie", pt: "5 aulas por dia", ro: "5 lecții pe zi", zh: "每天 5 节课", ja: "1日5レッスン", ko: "하루 5레슨" }),
        "15 сообщений практики в день": legalTranslations("15 practice messages per day", { es: "15 mensajes de práctica al día", de: "15 Übungsnachrichten pro Tag", fr: "15 messages de pratique par jour", it: "15 messaggi di pratica al giorno", uk: "15 повідомлень практики на день", pl: "15 wiadomości ćwiczeniowych dziennie", pt: "15 mensagens de prática por dia", ro: "15 mesaje de practică pe zi", zh: "每天 15 条练习消息", ja: "1日15件の練習メッセージ", ko: "하루 15개의 연습 메시지" }),
        "20 голосовых сообщений до 30 секунд": legalTranslations("20 voice messages up to 30 seconds", { es: "20 mensajes de voz de hasta 30 segundos", de: "20 Sprachnachrichten bis 30 Sekunden", fr: "20 messages vocaux jusqu'à 30 secondes", it: "20 messaggi vocali fino a 30 secondi", uk: "20 голосових повідомлень до 30 секунд", pl: "20 wiadomości głosowych do 30 sekund", pt: "20 mensagens de voz de até 30 segundos", ro: "20 mesaje vocale de până la 30 de secunde", zh: "20 条最长 30 秒的语音消息", ja: "30秒までの音声メッセージ20件", ko: "최대 30초 음성 메시지 20개" }),
        "60 голосовых сообщений до 30 секунд": legalTranslations("60 voice messages up to 30 seconds", { es: "60 mensajes de voz de hasta 30 segundos", de: "60 Sprachnachrichten bis 30 Sekunden", fr: "60 messages vocaux jusqu'à 30 secondes", it: "60 messaggi vocali fino a 30 secondi", uk: "60 голосових повідомлень до 30 секунд", pl: "60 wiadomości głosowych do 30 sekund", pt: "60 mensagens de voz de até 30 segundos", ro: "60 mesaje vocale de până la 30 de secunde", zh: "60 条最长 30 秒的语音消息", ja: "30秒までの音声メッセージ60件", ko: "최대 30초 음성 메시지 60개" }),
        "голос в текст, перевод услышанного и перевод текста с картинки": legalTranslations("voice to text, heard speech translation and image text translation", { es: "voz a texto, traducción de lo escuchado y traducción de texto en imágenes", de: "Sprache zu Text, Übersetzung des Gehörten und Bildtext-Übersetzung", fr: "voix en texte, traduction de l'audio entendu et traduction du texte d'une image", it: "voce in testo, traduzione dell'audio ascoltato e del testo nelle immagini", uk: "голос у текст, переклад почутого й переклад тексту з картинки", pl: "głos na tekst, tłumaczenie usłyszanego i tekstu z obrazu", pt: "voz para texto, tradução do que foi ouvido e texto em imagem", ro: "voce în text, traducerea audio și a textului din imagine", zh: "语音转文字、听到内容翻译和图片文字翻译", ja: "音声のテキスト化、聞き取った内容の翻訳、画像内テキスト翻訳", ko: "음성을 텍스트로 변환, 들은 내용 번역, 이미지 텍스트 번역" }),
        "100 уроков в день": legalTranslations("100 lessons per day", { es: "100 lecciones al día", de: "100 Lektionen pro Tag", fr: "100 leçons par jour", it: "100 lezioni al giorno", uk: "100 уроків на день", pl: "100 lekcji dziennie", pt: "100 aulas por dia", ro: "100 lecții pe zi", zh: "每天 100 节课", ja: "1日100レッスン", ko: "하루 100레슨" }),
        "500 сообщений практики в день": legalTranslations("500 practice messages per day", { es: "500 mensajes de práctica al día", de: "500 Übungsnachrichten pro Tag", fr: "500 messages de pratique par jour", it: "500 messaggi di pratica al giorno", uk: "500 повідомлень практики на день", pl: "500 wiadomości ćwiczeniowych dziennie", pt: "500 mensagens de prática por dia", ro: "500 mesaje de practică pe zi", zh: "每天 500 条练习消息", ja: "1日500件の練習メッセージ", ko: "하루 500개의 연습 메시지" }),
        "максимальный доступ к AI-диалогам, voice/photo-практике и переводчику": legalTranslations("maximum access to AI dialogues, voice/photo practice and translator", { es: "acceso máximo a diálogos IA, práctica de voz/foto y traductor", de: "maximaler Zugriff auf KI-Dialoge, Sprach-/Fotoübungen und Übersetzer", fr: "accès maximal aux dialogues IA, pratique voix/photo et traducteur", it: "accesso massimo a dialoghi IA, pratica voce/foto e traduttore", uk: "максимальний доступ до AI-діалогів, voice/photo-практики й перекладача", pl: "maksymalny dostęp do dialogów AI, praktyki voice/photo i tłumacza", pt: "acesso máximo a diálogos de IA, prática por voz/foto e tradutor", ro: "acces maxim la dialoguri AI, practică voce/foto și traducător", zh: "最大权限使用 AI 对话、语音/照片练习和翻译器", ja: "AI対話、音声/写真練習、翻訳機能への最大アクセス", ko: "AI 대화, 음성/사진 연습 및 번역기에 대한 최대 접근" }),
        "NERIVA предлагает бесплатный тариф Free и два платных тарифа Premium и Platinum:": legalTranslations("NERIVA offers the free Free plan and two paid plans, Premium and Platinum:", { es: "NERIVA ofrece el plan gratuito Free y dos planes de pago, Premium y Platinum:", de: "NERIVA bietet den kostenlosen Free-Tarif und zwei kostenpflichtige Tarife, Premium und Platinum:", fr: "NERIVA propose l'offre gratuite Free et deux offres payantes, Premium et Platinum :", it: "NERIVA offre il piano gratuito Free e due piani a pagamento, Premium e Platinum:", uk: "NERIVA пропонує безкоштовний тариф Free і два платні тарифи Premium та Platinum:", pl: "NERIVA oferuje bezpłatny plan Free oraz dwa płatne plany: Premium i Platinum:", pt: "O NERIVA oferece o plano gratuito Free e dois planos pagos, Premium e Platinum:", ro: "NERIVA oferă planul gratuit Free și două planuri plătite, Premium și Platinum:", zh: "NERIVA 提供免费 Free 套餐以及两个付费套餐 Premium 和 Platinum：", ja: "NERIVA は無料の Free プランと、有料の Premium / Platinum プランを提供します：", ko: "NERIVA는 무료 Free 요금제와 두 가지 유료 요금제 Premium 및 Platinum을 제공합니다:" }),
        "Тариф": legalTranslations("Plan", { es: "Plan", de: "Tarif", fr: "Offre", it: "Piano", uk: "Тариф", pl: "Plan", pt: "Plano", ro: "Plan", zh: "套餐", ja: "プラン", ko: "요금제" }),
        "Возможности": legalTranslations("Features", { es: "Funciones", de: "Funktionen", fr: "Fonctionnalités", it: "Funzioni", uk: "Можливості", pl: "Funkcje", pt: "Recursos", ro: "Funcții", zh: "功能", ja: "機能", ko: "기능" }),
        "Стоимость": legalTranslations("Price", { es: "Precio", de: "Preis", fr: "Prix", it: "Prezzo", uk: "Вартість", pl: "Cena", pt: "Preço", ro: "Preț", zh: "价格", ja: "料金", ko: "가격" }),
        "Тариф и период": legalTranslations("Plan and period", { es: "Plan y periodo", de: "Tarif und Zeitraum", fr: "Offre et période", it: "Piano e periodo", uk: "Тариф і період", pl: "Plan i okres", pt: "Plano e período", ro: "Plan și perioadă", zh: "套餐和期限", ja: "プランと期間", ko: "요금제 및 기간" }),
        "Цена в рублях": legalTranslations("Price in RUB", { es: "Precio en RUB", de: "Preis in RUB", fr: "Prix en RUB", it: "Prezzo in RUB", uk: "Ціна в RUB", pl: "Cena w RUB", pt: "Preço em RUB", ro: "Preț în RUB", zh: "卢布价格", ja: "RUB価格", ko: "RUB 가격" }),
        "Цена в Telegram Stars": legalTranslations("Price in Telegram Stars", { es: "Precio en Telegram Stars", de: "Preis in Telegram Stars", fr: "Prix en Telegram Stars", it: "Prezzo in Telegram Stars", uk: "Ціна в Telegram Stars", pl: "Cena w Telegram Stars", pt: "Preço em Telegram Stars", ro: "Preț în Telegram Stars", zh: "Telegram Stars 价格", ja: "Telegram Stars価格", ko: "Telegram Stars 가격" }),
        "около 4.15 USDT; TON по счёту": legalTranslations("about 4.15 USDT; TON by invoice", { es: "unos 4.15 USDT; TON según factura", de: "ca. 4.15 USDT; TON per Rechnung", fr: "environ 4.15 USDT ; TON selon facture", it: "circa 4.15 USDT; TON su fattura", uk: "близько 4.15 USDT; TON за рахунком", pl: "około 4.15 USDT; TON według faktury", pt: "cerca de 4.15 USDT; TON por fatura", ro: "aprox. 4.15 USDT; TON pe factură", zh: "约 4.15 USDT；TON 按账单", ja: "約4.15 USDT、TONは請求額", ko: "약 4.15 USDT; TON은 청구서 기준" }),
        "около 41.50 USDT; TON по счёту": legalTranslations("about 41.50 USDT; TON by invoice", { ko: "약 41.50 USDT; TON은 청구서 기준", zh: "约 41.50 USDT；TON 按账单", ja: "約41.50 USDT、TONは請求額" }),
        "около 8.15 USDT; TON по счёту": legalTranslations("about 8.15 USDT; TON by invoice", { ko: "약 8.15 USDT; TON은 청구서 기준", zh: "约 8.15 USDT；TON 按账单", ja: "約8.15 USDT、TONは請求額" }),
        "около 81.50 USDT; TON по счёту": legalTranslations("about 81.50 USDT; TON by invoice", { ko: "약 81.50 USDT; TON은 청구서 기준", zh: "约 81.50 USDT；TON 按账单", ja: "約81.50 USDT、TONは請求額" }),
        "Месячные цены Premium 300 ₽ / 150 Stars и Platinum 590 ₽ / 300 Stars уже указаны с учётом стартовой скидки 70%. Годовая оплата доступна на 365 дней и дополнительно снижает стоимость по сравнению с помесячной оплатой.": legalTranslations("Monthly prices Premium 300 RUB / 150 Stars and Platinum 590 RUB / 300 Stars already include the 70% launch discount. Annual payment is available for 365 days and additionally lowers the cost compared with monthly payment.", { ko: "월간 가격 Premium 300 RUB / 150 Stars 및 Platinum 590 RUB / 300 Stars에는 이미 출시 70% 할인이 포함되어 있습니다. 연간 결제는 365일 이용권이며 월 결제보다 비용이 더 낮습니다." }),
        "50 уроков в день, 200 сообщений практики, 20 голосовых до 30 секунд, голос в текст, перевод услышанного, перевод текста с картинки, практика по голосу или фото.": legalTranslations("50 lessons per day, 200 practice messages, 20 voice messages up to 30 seconds, voice to text, heard speech translation, image text translation and practice from voice or photo context.", { ko: "하루 50레슨, 연습 메시지 200개, 최대 30초 음성 20개, 음성 텍스트 변환, 들은 내용 번역, 이미지 텍스트 번역, 음성 또는 사진 맥락 연습." }),
        "100 уроков в день, 500 сообщений практики, 60 голосовых до 30 секунд, максимальный доступ к AI-диалогам, voice/photo-практике и переводчику.": legalTranslations("100 lessons per day, 500 practice messages, 60 voice messages up to 30 seconds, maximum access to AI dialogues, voice/photo practice and translator.", { ko: "하루 100레슨, 연습 메시지 500개, 최대 30초 음성 60개, AI 대화, 음성/사진 연습 및 번역기에 대한 최대 접근." }),
        "Платежные данные банковских карт и платежных инструментов обрабатываются платежными провайдерами, YooKassa, Telegram, банками, магазинами приложений, блокчейн-сетями или иными платежными участниками. NERIVA не запрашивает и не хранит полный номер банковской карты, CVV/CVC-код и иные полные реквизиты платежного инструмента.": legalTranslations("Payment data for bank cards and payment instruments is processed by payment providers, YooKassa, Telegram, banks, app stores, blockchain networks or other payment participants. NERIVA does not request or store the full bank card number, CVV/CVC code or other full payment instrument details.", { ko: "은행 카드 및 결제 수단 데이터는 결제 제공업체, YooKassa, Telegram, 은행, 앱 스토어, 블록체인 네트워크 또는 기타 결제 참여자가 처리합니다. NERIVA는 전체 카드 번호, CVV/CVC 코드 또는 기타 전체 결제 수단 정보를 요청하거나 저장하지 않습니다." }),
        "Для подтверждения Premium- или Platinum-доступа Администрация может обрабатывать платежный статус, идентификатор операции, выбранный тариф, срок доступа, сумму, валюту, количество Telegram Stars, TON/USDT-сеть, хеш транзакции и сведения, которые Пользователь сам передает в поддержку для проверки оплаты.": legalTranslations("To confirm Premium or Platinum access, the Administration may process payment status, transaction identifier, selected plan, access period, amount, currency, number of Telegram Stars, TON/USDT network, transaction hash and information the User provides to support for payment verification.", { ko: "Premium 또는 Platinum 이용을 확인하기 위해 운영자는 결제 상태, 거래 식별자, 선택한 요금제, 이용 기간, 금액, 통화, Telegram Stars 수량, TON/USDT 네트워크, 거래 해시 및 사용자가 결제 확인을 위해 지원팀에 제공한 정보를 처리할 수 있습니다." }),
        "сведения о тарифе: Free, Premium или Platinum, дата активации, срок доступа, лимиты, факт оплаты, платежный статус, выбранный способ оплаты, Telegram Stars, TON/USDT и данные для сверки платежа;": legalTranslations("plan information: Free, Premium or Platinum, activation date, access period, limits, payment fact, payment status, selected payment method, Telegram Stars, TON/USDT and payment verification data;", { ko: "요금제 정보: Free, Premium 또는 Platinum, 활성화 날짜, 이용 기간, 한도, 결제 사실, 결제 상태, 선택한 결제 방법, Telegram Stars, TON/USDT 및 결제 확인 데이터;" }),
        "учет дневных лимитов Free, Premium и Platinum: уроки, сообщения практики, голосовые и иные доступные функции;": legalTranslations("accounting for daily Free, Premium and Platinum limits: lessons, practice messages, voice messages and other available features;", { ko: "Free, Premium 및 Platinum의 일일 한도 계산: 레슨, 연습 메시지, 음성 및 기타 사용 가능한 기능;" })
    };

    Object.entries(legalPhraseTranslations).forEach(([source, values]) => {
        const key = normalizePhrase(source);
        phraseTranslations[key] = { ...values, ...(phraseTranslations[key] || {}) };
    });

    const landingEnglishPhrases = [
        ["NERIVA", "NERIVA"],
        ["уроки и диалоги", "lessons and dialogs"],
        ["голос в текст", "voice to text"],
        ["перевод с камеры", "camera translation"],
        ["языков", "languages"],
        ["Открыть онлайн-приложение", "Open web app"],
        ["Открыть Telegram", "Open Telegram"],
        ["уроков/день Free", "lessons/day Free"],
        ["уроков/день Platinum", "lessons/day Platinum"],
        ["Premium в месяц", "Premium monthly"],
        ["переводчик", "translator"],
        ["Настройте обучение:", "Set up learning:"],
        ["Интерфейс: RU", "Interface: RU"],
        ["Учить: EN", "Learn: EN"],
        ["Интерфейс: UZ", "Interface: UZ"],
        ["Учить: IT", "Learn: IT"],
        ["Урок", "Lesson"],
        ["Диалог", "Dialog"],
        ["Голос", "Voice"],
        ["Переводчик", "Translator"],
        ["Как звучит слово", "How does the word"],
        ["? Хочу услышать произношение.", "sound? I want to hear the pronunciation."],
        ["Готово!", "Done!"],
        ["Отправляю голосовое с произношением и короткий пример в контексте.", "Sending a voice pronunciation and a short example in context."],
        ["Пример голосовой озвучки", "Voice pronunciation example"],
        ["Напишите сообщение...", "Write a message..."],
        ["7 дней бесплатно", "7 days free"],
        ["Мультиязычный продукт", "Multilingual product"],
        ["Сайт и бот для разных стран", "Website and bot for different countries"],
        ["Пользователь может выбрать язык интерфейса онлайн-приложения и бота, а отдельно выбрать язык обучения. Например, пользоваться сервисом на узбекском, а учить английский; или пользоваться на украинском, а учить итальянский.", "Users can choose the interface language for the web app and bot, and separately choose the language they want to learn. For example, use the service in Uzbek while learning English, or use it in Ukrainian while learning Italian."],
        ["Язык сайта и бота", "Website and bot language"],
        ["Интерфейс можно сделать понятным для пользователя на его языке, без привязки к русскому.", "The interface can be clear in the user's own language, without being tied to Russian."],
        ["Язык обучения", "Learning language"],
        ["Язык, который человек изучает, выбирается отдельно: это делает сценарии гибкими для разных стран.", "The language a person studies is selected separately, which keeps scenarios flexible for different countries."],
        ["английский", "English"],
        ["испанский", "Spanish"],
        ["немецкий", "German"],
        ["французский", "French"],
        ["итальянский", "Italian"],
        ["таджикский", "Tajik"],
        ["узбекский", "Uzbek"],
        ["татарский", "Tatar"],
        ["армянский", "Armenian"],
        ["казахский", "Kazakh"],
        ["украинский", "Ukrainian"],
        ["польский", "Polish"],
        ["румынский", "Romanian"],
        ["португальский", "Portuguese"],
        ["кыргызский", "Kyrgyz"],
        ["грузинский", "Georgian"],
        ["AI-платформа для языковой практики", "AI platform for language practice"],
        ["Один сервис для уроков, диалогов, перевода, голоса и фото", "One service for lessons, dialogs, translation, voice and photos"],
        ["NERIVA закрывает главную боль языкового обучения: человек не просто читает теорию, а каждый день говорит, пишет, слушает, получает исправления и сразу видит свои ошибки. Это не отдельное тяжёлое приложение, а связка Telegram-бота и web app с одним прогрессом, оплатой в рублях, Stars, TON и USDT.", "NERIVA solves the main pain of language learning: people do not just read theory, they speak, write, listen, get corrections and see their mistakes every day. It is not a heavy separate app, but a Telegram bot and web app with shared progress and payments in RUB, Stars, TON and USDT."],
        ["20 языков", "20 languages"],
        ["Telegram-вход", "Telegram login"],
        ["Мультиязычный сайт", "Multilingual website"],
        ["Фото", "Photo"],
        ["Скидка только первый месяц запуска", "Discount only for the first launch month"],
        ["Premium на 30 дней", "Premium for 30 days"],
        ["500 Stars · около 14 USDT", "500 Stars · about 14 USDT"],
        ["150 Stars · около 4.15 USDT · TON по счёту", "150 Stars · about 4.15 USDT · TON by invoice"],
        ["Получить Premium онлайн", "Get Premium online"],
        ["Стартовая цена", "Launch price"],
        ["70% скидки только первый месяц запуска", "70% discount only for the first launch month"],
        ["Сейчас выгодный момент подключить голос, фото, расширенные уроки и практику: цены уже показаны со скидкой. После стартового месяца базовые цены вернутся к округлённому уровню около 1 000 ₽ за Premium и около 2 000 ₽ за Platinum в месяц.", "Now is the best moment to unlock voice, photos, expanded lessons and practice: prices are already shown with the discount. After the launch month, base prices return to the rounded level of about 1,000 RUB for Premium and about 2,000 RUB for Platinum per month."],
        ["Забрать стартовую цену", "Claim the launch price"],
        ["Premium месяц", "Premium monthly"],
        ["Platinum месяц", "Platinum monthly"],
        ["Годовой доступ", "Annual access"],
        ["AI-ядро обучения", "AI learning core"],
        ["Бот не просто переводит, а собирает практику вокруг слова", "The bot does not just translate; it builds practice around each word"],
        ["Новое слово можно услышать голосом, разобрать в контексте, потренировать в диалоге и закрепить коротким заданием. Всё доступно на сайте и в Telegram, без установки отдельного приложения.", "A new word can be heard aloud, understood in context, practiced in a dialog and reinforced with a short task. Everything is available on the website and in Telegram without installing a separate app."],
        ["Слушать произношение", "Hear pronunciation"],
        ["AI отправляет голосовое, чтобы пользователь услышал реальное звучание слова.", "AI sends a voice message so the user can hear how the word really sounds."],
        ["Диалоговая практика", "Dialog practice"],
        ["Слова сразу попадают в живые фразы, вопросы и ответы под уровень пользователя.", "Words immediately become live phrases, questions and answers matched to the user's level."],
        ["Голосовые задания", "Voice tasks"],
        ["Premium открывает голосовую практику, распознавание речи и перевод услышанного.", "Premium unlocks voice practice, speech recognition and translation of heard speech."],
        ["Фото и контекст", "Photo and context"],
        ["Фото меню, страницы или задания превращается в перевод и полезные упражнения.", "A photo of a menu, page or task turns into a translation and useful exercises."],
        ["Возможности", "Features"],
        ["Что умеет NERIVA", "What NERIVA can do"],
        ["Бот помогает учить 20 языков через короткие уроки, переписку, голос, озвучку слов и материалы, которые вы отправляете сами. Интерфейс сайта и бота тоже стал мультиязычным.", "The bot helps people learn 20 languages through short lessons, chat practice, voice, word pronunciation and the materials they send themselves. The website and bot interface is multilingual too."],
        ["Сравнение", "Comparison"],
        ["Free, Premium или Platinum?", "Free, Premium or Platinum?"],
        ["Free — чтобы попробовать. Premium — лучший выбор на каждый день. Platinum — максимум лимитов для интенсивных рывков.", "Free is for trying the service. Premium is the best everyday choice. Platinum gives maximum limits for intensive learning sprints."],
        ["Функция", "Feature"],
        ["Для чего подходит", "Use cases"],
        ["Практика под реальные задачи", "Practice for real tasks"],
        ["Старт", "Start"],
        ["Как начать обучение", "How to start learning"],
        ["Три шага — и можно заниматься на сайте или в Telegram.", "Three steps and you can study on the website or in Telegram."],
        ["Тарифы", "Plans"],
        ["Прозрачные лимиты и цены", "Clear limits and prices"],
        ["Начните бесплатно, подключите Premium для ежедневной практики или Platinum для максимальных лимитов.", "Start for free, upgrade to Premium for daily practice or Platinum for maximum limits."],
        ["Premium бесплатно на 7 дней", "Premium free for 7 days"],
        ["Пригласите друга в NERIVA и получите 7 дней Premium бесплатно. Отличный способ протестировать голосовые, перевод услышанного и работу с фото без оплаты.", "Invite a friend to NERIVA and get 7 days of Premium for free. A great way to test voice tasks, heard speech translation and photo features without paying."],
        ["Пригласить друга", "Invite a friend"],
        ["Попробуйте NERIVA сегодня", "Try NERIVA today"],
        ["Начните с бесплатного тарифа: 5 уроков и 15 сообщений практики в день. Когда захотите больше — подключите Premium или Platinum на 30 дней либо на год.", "Start with the free plan: 5 lessons and 15 practice messages per day. When you want more, upgrade to Premium or Platinum for 30 days or a year."],
        ["Сайт и Telegram", "Website and Telegram"],
        ["YooKassa/СБП, Stars, TON или USDT", "YooKassa/SBP, Stars, TON or USDT"],
        ["7 дней за приглашение друга", "7 days for inviting a friend"],
        ["Частые вопросы", "FAQ"],
        ["Коротко о доступе, лимитах, оплате и Premium-функциях.", "Quick answers about access, limits, payments and Premium features."],
        ["Онлайн-платформа и Telegram-бот", "Web platform and Telegram bot"],
        ["Ежедневные уроки, AI-практика, Telegram-вход, синхронизация прогресса, переводчик, голосовые задания, озвучка слов, перевод услышанного и текста с фото для 20 языков.", "Daily lessons, AI practice, Telegram login, progress sync, translator, voice tasks, word pronunciation, heard speech translation and photo text translation for 20 languages."],
        ["Telegram бот", "Telegram bot"],
        ["Навигация", "Navigation"],
        ["Контакты", "Contacts"],
        ["Онлайн:", "Online:"],
        ["Бот:", "Bot:"],
        ["Поддержка в Telegram:", "Telegram support:"],
        ["Открыть онлайн", "Open online"],
        ["© 2026 NERIVA. Все права защищены.", "© 2026 NERIVA. All rights reserved."],
        ["Онлайн-приложение", "Web app"],
        ["Политика конфиденциальности", "Privacy policy"],
        ["Условия использования", "Terms of use"],
        ["Главная", "Home"],
        ["AI-репетитор, а не просто чат", "AI tutor, not just a chat"],
        ["Бот ведёт по короткому циклу: объяснил, дал задание, проверил, сохранил ошибку и вернул слово в тренировку.", "The bot follows a short cycle: explains, gives a task, checks it, saves the mistake and returns the word to practice."],
        ["Голос, фото и камера", "Voice, photos and camera"],
        ["Можно отправить голосовое, фото из файлов или снимок с камеры: AI распознаёт, переводит и превращает материал в практику.", "You can send a voice message, a file photo or a camera shot: AI recognizes it, translates it and turns the material into practice."],
        ["20 языков уже доступны", "20 languages already available"],
        ["Выбирайте язык интерфейса и язык обучения отдельно: от English, Español и Deutsch до Қазақша, O'zbekcha, Кыргызча и ქართული.", "Choose the interface language and learning language separately: from English, Spanish and German to Kazakh, Uzbek, Kyrgyz and Georgian."],
        ["Мультиязычный интерфейс", "Multilingual interface"],
        ["Сайт и бот больше не привязаны к русскому: пользователю проще начать, если навигация и подсказки доступны на понятном ему языке.", "The website and bot are no longer tied to Russian: it is easier to start when navigation and hints are available in a language the user understands."],
        ["Единый вход через Telegram", "Single Telegram login"],
        ["Можно войти на сайте через Telegram и связать прогресс веб-аккаунта с ботом. Если данные уже есть в Telegram, они становятся основными.", "Users can log in on the website through Telegram and link web account progress with the bot. If Telegram already has data, it becomes the main source."],
        ["AI-переводчик с контекстом", "AI translator with context"],
        ["Переводчик работает с выбранной парой языков, принимает текст, голос и фото, а готовый перевод можно прослушать.", "The translator works with the selected language pair, accepts text, voice and photos, and the finished translation can be played aloud."],
        ["Диалоги вместо зубрёжки", "Dialogs instead of rote learning"],
        ["Тренируйте живые ответы в переписке: 15 сообщений в день на Free, 200 на Premium и 500 на Platinum.", "Practice live replies in chat: 15 messages per day on Free, 200 on Premium and 500 on Platinum."],
        ["Ошибки превращаются в план", "Mistakes turn into a plan"],
        ["Сервис сохраняет частые ошибки и помогает возвращаться к ним через отдельный словарь и практику исправлений.", "The service saves frequent mistakes and helps users return to them through a separate dictionary and correction practice."],
        ["Premium открывает 20 голосовых в день, Platinum — 60 голосовых в день до 30 секунд.", "Premium unlocks 20 voice messages per day, Platinum unlocks 60 voice messages per day up to 30 seconds."],
        ["Голос в текст", "Voice to text"],
        ["Бот переводит услышанное в текст, чтобы вы могли разобрать фразы, слова и ошибки.", "The bot converts heard speech to text so you can analyze phrases, words and mistakes."],
        ["Произношение слов голосом", "Voice pronunciation for words"],
        ["При изучении новых слов бот может отправить голосовое, чтобы вы услышали, как слово звучит на выбранном языке.", "When learning new words, the bot can send a voice message so you hear how the word sounds in the selected language."],
        ["Практика по контексту", "Context practice"],
        ["Бот может строить задания вокруг вашего голоса или фото, поэтому обучение ощущается ближе к реальной жизни.", "The bot can build tasks around your voice or photo, so learning feels closer to real life."],
        ["Можно выбрать", "Selectable"],
        ["Вход через Telegram и синхронизация", "Telegram login and sync"],
        ["Да", "Yes"],
        ["Уроки в день", "Lessons per day"],
        ["Сообщения практики в день", "Practice messages per day"],
        ["AI-переводчик текста", "AI text translator"],
        ["Недоступны", "Unavailable"],
        ["20 в день", "20 per day"],
        ["60 в день", "60 per day"],
        ["Голос в текст и перевод услышанного", "Voice to text and heard speech translation"],
        ["Нет", "No"],
        ["Озвучка произношения слов", "Word pronunciation audio"],
        ["Фото, камера и OCR-перевод", "Photo, camera and OCR translation"],
        ["Практика по контексту голоса или фото", "Practice from voice or photo context"],
        ["Цена за месяц", "Monthly price"],
        ["Цена за год", "Annual price"],
        ["Разговорная практика", "Speaking practice"],
        ["Пишите и отправляйте голосовые, чтобы тренировать реальные диалоги.", "Write and send voice messages to practice real dialogs."],
        ["Учёба и задания", "Study and assignments"],
        ["Разбирайте тексты, упражнения и незнакомые слова на выбранном языке.", "Break down texts, exercises and unfamiliar words in the selected language."],
        ["Путешествия", "Travel"],
        ["Переводите вывески, меню и фразы, которые встречаете вокруг.", "Translate signs, menus and phrases you encounter around you."],
        ["Ежедневная привычка", "Daily habit"],
        ["Короткие лимиты помогают заниматься регулярно, без перегруза.", "Short limits help users study regularly without overload."],
        ["Откройте онлайн или Telegram", "Open online or Telegram"],
        ["Перейдите в веб-приложение, войдите через Telegram или найдите @NERIVAapp_bot в Telegram.", "Go to the web app, log in through Telegram or find @NERIVAapp_bot in Telegram."],
        ["Выберите языки", "Choose languages"],
        ["Настройте язык интерфейса бота и язык, который хотите изучать.", "Set the bot interface language and the language you want to learn."],
        ["Учитесь и переводите каждый день", "Learn and translate every day"],
        ["Используйте уроки, практику, слова и переводчик; Premium добавит голос, фото, камеру и расширенные лимиты.", "Use lessons, practice, words and the translator; Premium adds voice, photos, camera and expanded limits."],
        ["Для старта", "For starters"],
        ["навсегда", "forever"],
        ["Бесплатный тариф, чтобы познакомиться с ботом и выработать привычку заниматься.", "A free plan for getting to know the bot and building a learning habit."],
        ["выбор языка интерфейса и обучения", "choice of interface and learning language"],
        ["Голосовые недоступны", "Voice features unavailable"],
        ["Подходит для лёгкого ежедневного обучения", "Good for light daily learning"],
        ["Начать бесплатно", "Start free"],
        ["Лучший выбор · скидка 70%", "Best choice · 70% discount"],
        ["в месяц", "per month"],
        ["150 Stars · около 4.15 USDT", "150 Stars · about 4.15 USDT"],
        ["≈10 000 ₽/год", "≈10,000 RUB/year"],
        ["3 000 ₽/год · 1 500 Stars · около 41.50 USDT", "3,000 RUB/year · 1,500 Stars · about 41.50 USDT"],
        ["70% уже в цене · только первый месяц запуска", "70% already included · only for the first launch month"],
        ["Для ежедневного прогресса: больше уроков, диалогов, голос, фото и словарь ошибок.", "For daily progress: more lessons, dialogs, voice, photos and a mistake dictionary."],
        ["50 уроков в день", "50 lessons per day"],
        ["200 сообщений практики в день", "200 practice messages per day"],
        ["20 голосовых в день до 30 секунд", "20 voice messages per day up to 30 seconds"],
        ["AI-переводчик с выбранными языками", "AI translator with selected languages"],
        ["Перевод текста с фото из проводника или камеры", "Text translation from a file photo or camera"],
        ["Подключить Premium", "Upgrade to Premium"],
        ["Максимум доступа · скидка 70%", "Maximum access · 70% discount"],
        ["1 000 Stars · около 27 USDT", "1,000 Stars · about 27 USDT"],
        ["300 Stars · около 8.15 USDT", "300 Stars · about 8.15 USDT"],
        ["≈20 000 ₽/год", "≈20,000 RUB/year"],
        ["5 900 ₽/год · 3 000 Stars · около 81.50 USDT", "5,900 RUB/year · 3,000 Stars · about 81.50 USDT"],
        ["Для интенсивных рывков, поездок, экзаменов и тех, кто хочет максимум лимитов.", "For intensive sprints, trips, exams and users who want maximum limits."],
        ["Все возможности Premium", "All Premium features"],
        ["60 голосовых в день до 30 секунд", "60 voice messages per day up to 30 seconds"],
        ["Годовая оплата: 5900 ₽ / 3000 Stars / около 81.50 USDT", "Annual payment: 5,900 RUB / 3,000 Stars / about 81.50 USDT"],
        ["Выбрать Platinum", "Choose Platinum"],
        ["Что входит в бесплатный тариф Free?", "What is included in the free Free plan?"],
        ["Free включает 5 уроков в день и 15 сообщений практики в день. Голосовые функции на бесплатном тарифе недоступны.", "Free includes 5 lessons per day and 15 practice messages per day. Voice features are unavailable on the free plan."],
        ["Что добавляет Premium?", "What does Premium add?"],
        ["Premium открывает 50 уроков в день, 200 сообщений практики в день, 20 голосовых в день до 30 секунд, голос в текст, перевод услышанного, озвучку произношения слов, переводчик с голосом, фото и камерой, а также практику по контексту голоса или фото.", "Premium unlocks 50 lessons per day, 200 practice messages per day, 20 voice messages per day up to 30 seconds, voice to text, heard speech translation, word pronunciation audio, translator with voice, photos and camera, plus practice from voice or photo context."],
        ["Чем Platinum отличается от Premium?", "How is Platinum different from Premium?"],
        ["Platinum создан для интенсивного обучения: 100 уроков, 500 сообщений практики и 60 голосовых в день. Это тариф для учебных рывков, поездок, подготовки к собеседованию или экзамену.", "Platinum is built for intensive learning: 100 lessons, 500 practice messages and 60 voice messages per day. It is a plan for learning sprints, trips, interview preparation or exams."],
        ["Можно ли войти на сайте через Telegram?", "Can I log in on the website through Telegram?"],
        ["Да. Веб-версия поддерживает вход через Telegram: прогресс сайта и бота связывается с одной учёткой. Если Telegram уже содержит прогресс, он берётся за основу; если Telegram был пустой, данные сайта переезжают в Telegram.", "Yes. The web version supports Telegram login: website and bot progress are linked to one account. If Telegram already contains progress, it becomes the base; if Telegram was empty, website data moves to Telegram."],
        ["Что умеет переводчик в инструментах?", "What can the translator in tools do?"],
        ["Переводчик позволяет выбрать языки, перевести текст, голос, фото из проводника или снимок с камеры и прослушать готовый перевод. Текстовый перевод доступен сразу, голос и фото относятся к Premium-функциям.", "The translator lets you choose languages, translate text, voice, a file photo or a camera shot, and listen to the finished translation. Text translation is available immediately; voice and photo are Premium features."],
        ["Можно ли слушать произношение слов?", "Can I listen to word pronunciation?"],
        ["Да. Во время изучения слов бот может отправлять голосовое с произношением, чтобы вы слышали, как слово звучит на выбранном языке.", "Yes. While learning words, the bot can send a voice message with pronunciation so you hear how the word sounds in the selected language."],
        ["Какие языки доступны?", "Which languages are available?"],
        ["Сейчас доступны English, Русский, Español, Deutsch, Français, Italiano, 中文, 日本語, 한국어, Тоҷикӣ, O'zbekcha, Татарча, Հայերեն, Қазақша, Кыргызча, ქართული, Українська, Polski, Română и Português.", "Currently available: English, Russian, Spanish, German, French, Italian, Chinese, Japanese, Korean, Tajik, Uzbek, Tatar, Armenian, Kazakh, Kyrgyz, Georgian, Ukrainian, Polish, Romanian and Portuguese."],
        ["Можно ли выбрать язык бота отдельно от языка обучения?", "Can I choose the bot language separately from the learning language?"],
        ["Да. Пользователь может выбрать язык сайта, онлайн-приложения и интерфейса бота отдельно от языка, который он хочет изучать. Например, пользоваться сервисом на O'zbekcha, а учить English.", "Yes. Users can choose the website, web app and bot interface language separately from the language they want to learn. For example, use the service in Uzbek while learning English."],
        ["Сколько стоит Premium?", "How much does Premium cost?"],
        ["Premium стоит 300 ₽ в месяц, 150 Telegram Stars или около 4.15 USDT. Годовой доступ стоит 3000 ₽, 1500 Stars или около 41.50 USDT. Platinum стоит 590 ₽ в месяц или 5900 ₽ в год.", "Premium costs 300 RUB per month, 150 Telegram Stars or about 4.15 USDT. Annual access costs 3,000 RUB, 1,500 Stars or about 41.50 USDT. Platinum costs 590 RUB per month or 5,900 RUB per year."],
        ["Как получить Premium бесплатно на 7 дней?", "How do I get Premium free for 7 days?"],
        ["Пригласите друга в NERIVA — за приглашение можно получить 7 дней Premium бесплатно.", "Invite a friend to NERIVA and get 7 days of Premium for free."],
        ["Нужно ли устанавливать отдельное приложение?", "Do I need to install a separate app?"],
        ["Нет. Можно заниматься в онлайн-приложении на сайте или через Telegram: откройте", "No. You can study in the web app on the website or through Telegram: open"],
        ["либо @NERIVAapp_bot.", "or @NERIVAapp_bot."],
        ["Для кого подходит бот?", "Who is the bot for?"],
        ["Для пользователей из разных стран, которые хотят учить языки короткими ежедневными сессиями: через уроки, переписку, голосовые, фото и практику реальных ситуаций.", "For users from different countries who want to learn languages in short daily sessions through lessons, chat, voice messages, photos and real-life practice."],
        ["100 уроков в Platinum", "100 lessons in Platinum"],
        ["500 сообщений практики", "500 practice messages"],
        ["Вход через Telegram", "Telegram login"],
        ["AI-переводчик", "AI translator"],
        ["Фото и камера", "Photo and camera"]
    ];

    landingEnglishPhrases.push(
        ["&copy; 2026 NERIVA. Все права защищены.", "© 2026 NERIVA. All rights reserved."],
        ["Язык", "Language"],
        ["языки онлайн и в Telegram", "languages online and in Telegram"],
        ["Онлайн", "Online"],
        ["Открыть меню", "Open menu"],
        ["AI-репетитор в Telegram и web app: уроки, голос, фото, ошибки и словарь", "AI tutor in Telegram and web app: lessons, voice, photos, mistakes and vocabulary"],
        ["Заговорите на новом языке быстрее: AI-репетитор уже внутри Telegram", "Speak a new language faster: your AI tutor is already inside Telegram"],
        ["NERIVA превращает привычный мессенджер в личного языкового тренера: даёт короткие уроки, ведёт диалог, разбирает ошибки, тренирует слова, понимает голос и переводит текст с фото.", "NERIVA turns a familiar messenger into a personal language coach: it gives short lessons, holds dialogs, explains mistakes, trains words, understands voice and translates text from photos."],
        ["Веб-приложение синхронизируется с Telegram, а Premium сейчас стоит 300 ₽, 150 Stars или около 4.15 USDT в месяц с учётом стартовой скидки 70%. Скидка действует только первый месяц запуска.", "The web app syncs with Telegram, and Premium now costs 300 RUB, 150 Stars or about 4.15 USDT per month with the 70% launch discount included. The discount is available only during the first launch month."],
        ["20 языков, единый web/Telegram-аккаунт, AI-уроки, диалоги, голос, фото, камера, словарь ошибок и тарифы Free, Premium, Platinum.", "20 languages, one web/Telegram account, AI lessons, dialogs, voice, photos, camera, mistake dictionary and Free, Premium, Platinum plans."],
        ["NERIVA — AI-репетитор в Telegram и web app", "NERIVA — AI tutor in Telegram and web app"],
        ["NERIVA — AI-репетитор в Telegram и web app для 20 языков", "NERIVA — AI tutor in Telegram and web app for 20 languages"],
        ["NERIVA — AI-репетитор в Telegram и веб-приложении: уроки, диалоги, голос, фото-перевод, словарь, ошибки, CEFR-тест, синхронизация прогресса и тарифы Free, Premium, Platinum. Premium от 300 ₽, 150 Stars или около 4.15 USDT в месяц.", "NERIVA is an AI tutor in Telegram and the web app: lessons, dialogs, voice, photo translation, vocabulary, mistakes, CEFR test, progress sync and Free, Premium, Platinum plans. Premium starts from 300 RUB, 150 Stars or about 4.15 USDT per month."],
        ["изучение языков, веб-приложение для языков, Telegram бот, вход через Telegram, AI переводчик, переводчик с фото, переводчик голосом, NERIVA, AI обучение, английский язык, русский язык, испанский язык, немецкий язык, французский язык, итальянский язык, таджикский язык, узбекский язык, татарский язык, армянский язык, казахский язык, украинский язык, польский язык, румынский язык, португальский язык, кыргызский язык, грузинский язык, голосовая практика, произношение слов, YooKassa, СБП, Telegram Stars", "language learning, language web app, Telegram bot, Telegram login, AI translator, photo translator, voice translator, NERIVA, AI learning, English, Russian, Spanish, German, French, Italian, Tajik, Uzbek, Tatar, Armenian, Kazakh, Ukrainian, Polish, Romanian, Portuguese, Kyrgyz, Georgian, voice practice, word pronunciation, YooKassa, SBP, Telegram Stars"],
        ["Нет. Можно заниматься в онлайн-приложении на сайте или через Telegram: откройте <span data-app-url-text>neriva.ru/app</span> либо @NERIVAapp_bot.", "No. You can study in the web app on the website or through Telegram: open <span data-app-url-text>neriva.ru/app</span> or @NERIVAapp_bot."],
        ["AI-репетитор в Telegram и web app: уроки, аудирование, голосовая оценка, фото и словарь", "AI tutor in Telegram and web app: lessons, listening, pronunciation scoring, photos and vocabulary"],
        ["AI-репетитор, который слышит речь и ведёт к живому языку", "An AI tutor that hears your speech and guides you toward real language"],
        ["NERIVA объединяет уроки, практику, аудирование, словарь, переводчик и оценку произношения. Пользователь занимается в Telegram или web app, получает структуру задания, голосовую обратную связь и понятный прогресс.", "NERIVA combines lessons, practice, listening, vocabulary, translator tools and pronunciation scoring. Learn in Telegram or the web app, get structured tasks, voice feedback and visible progress."],
        ["После голосового ответа пользователь видит понятный учебный разбор: примерную оценку, слабые слова, плавность речи и короткий совет, который можно сразу повторить вслух.", "After a voice answer, the learner sees clear coaching feedback: an approximate score, weak words, speech flow and a short tip they can repeat aloud right away."],
        ["структура и задания", "structure and tasks"],
        ["10 000 фраз", "10,000 phrases"],
        ["оценка произношения", "pronunciation scoring"],
        ["языков интерфейса", "interface languages"],
        ["фраз аудирования", "listening phrases"],
        ["оценка произношения", "pronunciation score"],
        ["уроки и практика", "lessons and practice"],
        ["Практика", "Practice"],
        ["Аудирование", "Listening"],
        ["Повторите фразу:", "Repeat the sentence:"],
        ["Оценка 84/100.", "Score 84/100."],
        ["Слово nearest прозвучало неясно. Скажите его медленнее и не проглатывайте финальный звук.", "The word nearest sounded unclear. Say it slower and do not swallow the final sound."],
        ["Голосовая оценка", "Voice scoring"],
        ["Произношение проверяется там, где пользователь говорит", "Pronunciation is checked wherever the user speaks"],
        ["В уроке, практике и аудировании пользователь может ответить голосом и получить понятную подсказку: что прозвучало уверенно, какие слова повторить и как сделать фразу плавнее.", "In lessons, practice and listening, users can answer by voice and get clear feedback: what sounded confident, which words to repeat and how to make the phrase smoother."],
        ["Эталонная фраза", "Target sentence"],
        ["Урок, практика или аудирование знают, что именно пользователь должен сказать.", "The lesson, practice or listening mode knows what the user is expected to say."],
        ["Разбор речи", "Speech feedback"],
        ["Система помогает понять, какие слова прозвучали неясно и где стоит говорить спокойнее.", "The system helps identify unclear words and where to speak more steadily."],
        ["Быстрая подсказка", "Fast coaching tip"],
        ["После ответа пользователь получает короткий следующий шаг, а не технический отчёт.", "After each answer, the learner gets a short next step, not a technical report."],
        ["Совет человеку", "Human-friendly advice"],
        ["Пользователь видит оценку, силу акцента, проблемные слова и короткий практический совет.", "The user sees the score, accent strength, problem words and a short practical tip."],
        ["Отзывы", "Reviews"],
        ["Пользователи видят прогресс, а не просто чат", "Users see progress, not just a chat"],
        ["Новый сайт показывает реальные учебные сценарии: урок, аудирование, голосовая оценка, словарь с контекстом и ежедневную мотивацию через награды.", "The new site shows real learning scenarios: lesson, listening, voice scoring, vocabulary with context and daily motivation through awards."],
        ["Нравится, что после голосового ответа видно, где именно речь съехала. Не просто “молодец”, а нормальный совет.", "I like that after a voice answer I can see where my speech slipped. Not just “good job”, but a useful tip."],
        ["Аудирование стало как мини-тренажёр: короткая фраза, запись, оценка, повтор. Очень быстро входит в привычку.", "Listening became a mini trainer: short sentence, recording, score, repeat. It quickly becomes a habit."],
        ["Словарь наконец даёт пример в контексте. Я понимаю не перевод слова, а как оно реально живёт в предложении.", "The vocabulary finally gives an example in context. I understand not just the translation, but how the word lives in a sentence."],
        ["Privacy и Terms встроены в продуктовую историю", "Privacy and Terms are part of the product story"],
        ["Голосовая AI-проверка описана честно", "AI voice checking is explained honestly"],
        ["Privacy и Terms объясняют функцию простым языком: голос используется для учебной обратной связи по произношению, а оценка является тренировочной подсказкой, не экзаменом.", "Privacy and Terms explain the feature in plain language: voice is used for educational pronunciation feedback, and the score is a training hint, not an exam."],
        ["Обработка голоса, платежей, аккаунта и учебного прогресса.", "Processing of voice, payments, account data and learning progress."],
        ["Правила сервиса, лимиты, тарифы и роль AI-оценки.", "Service rules, limits, plans and the role of AI scoring."],
        ["Оценка произношения в уроках и практике", "Pronunciation scoring in lessons and practice"],
        ["Голосовой ввод теперь не просто превращается в текст: пользователь получает оценку, силу акцента, проблемные слова и понятный совет.", "Voice input is no longer just turned into text: the user receives a score, accent strength, problem words and a clear tip."],
        ["Аудирование вместо Shadowing", "Listening instead of Shadowing"],
        ["Режим аудирования работает на выбранном языке обучения, использует большой банк фраз и проверяет произношение по эталону.", "The listening mode works in the selected learning language, uses a large phrase bank and checks pronunciation against the target."],
        ["Оценка произношения без технического шума", "Pronunciation scoring without technical noise"],
        ["Аудирование с 10 000 фраз", "Listening with 10,000 phrases"],
        ["Как работает оценка произношения?", "How does pronunciation scoring work?"],
        ["Пользователь отправляет голосовой ответ и получает понятную учебную подсказку: примерную оценку, слова для повторения и совет, как произнести фразу увереннее.", "The learner sends a voice answer and gets clear coaching feedback: an approximate score, words to repeat and a tip for saying the phrase more confidently."],
        ["Если пользователь отправляет голосовой ответ для тренировки, сервис может обработать запись, чтобы сравнить её с учебной фразой и показать, какие слова стоит произнести яснее.", "If the user sends a voice answer for training, the service may process the recording to compare it with the learning phrase and show which words should be pronounced more clearly."],
        ["Результат используется как образовательная подсказка: примерная оценка, проблемные слова и короткий совет для следующей попытки. Такая проверка не является экзаменом, медицинской или профессиональной фонетической экспертизой.", "The result is used as an educational hint: an approximate score, problem words and a short tip for the next attempt. This check is not an exam, medical assessment or professional phonetic expertise."],
        ["оценку произношения в уроках, практике и аудировании с учебным разбором похожести на фразу, понятности речи, темпа и пауз;", "pronunciation scoring in lessons, practice and listening with educational feedback on phrase match, speech clarity, pace and pauses;"],
        ["Голосовая оценка произношения", "Voice pronunciation scoring"],
        ["Если Пользователь отправляет голосовой ответ, сервис может обработать запись, чтобы дать учебную обратную связь по произношению: примерную оценку, слова для повторения и короткий совет. Такая проверка является тренировочной подсказкой и не считается профессиональной фонетической экспертизой, экзаменационным результатом или гарантией отсутствия акцента.", "If the User sends a voice answer, the service may process the recording to provide educational pronunciation feedback: an approximate score, words to repeat and a short tip. This check is a training hint and is not professional phonetic expertise, an exam result or a guarantee of no accent."],
        ["Базово", "Basic"],
        ["Расширенно", "Expanded"],
        ["Максимум", "Maximum"],
        ["Голосовая практика", "Voice practice"],
        ["Говорите, а не просто читаете ответы", "Speak instead of only reading answers"],
        ["NERIVA превращает каждую голосовую попытку в понятный следующий шаг: что прозвучало уверенно, какие слова стоит повторить и как сказать фразу плавнее. Это работает в аудировании, уроках и практике, поэтому речь тренируется там же, где идёт обучение.", "NERIVA turns every voice attempt into a clear next step: what sounded confident, which words to repeat and how to say the phrase more smoothly. It works in listening, lessons and practice, so speech is trained right where learning happens."],
        ["После голосового ответа пользователь видит:", "After a voice answer, the user sees:"],
        ["оценку произношения, слабые слова, плавность речи, силу акцента и короткий совет, который можно сразу повторить вслух.", "pronunciation score, weak words, speech flow, accent strength and a short tip that can be repeated aloud right away."],
        ["Сначала услышать", "Hear it first"],
        ["Бот даёт короткую живую фразу на выбранном языке и уровне, чтобы было понятно, как она должна звучать.", "The bot gives a short live phrase in the selected language and level, so it is clear how it should sound."],
        ["Повторить голосом", "Repeat by voice"],
        ["Пользователь отвечает как в реальном разговоре: не выбирает вариант, а тренирует рот, темп и интонацию.", "The user answers like in a real conversation: not choosing an option, but training mouth movement, pace and intonation."],
        ["Понять слабое место", "Find the weak spot"],
        ["Вместо сухого “правильно/неправильно” сервис показывает слова, темп и простую причину, почему фраза звучит неуверенно.", "Instead of a dry “right/wrong”, the service shows words, pace and a simple reason why the phrase sounds uncertain."],
        ["Сразу закрепить", "Reinforce it right away"],
        ["После разбора можно повторить фразу, перейти к новой или продолжить урок — без отдельного тренажёра и потери контекста.", "After the review, the user can repeat the phrase, move to a new one or continue the lesson without a separate trainer or losing context."],
        ["Мария", "Maria"],
        ["Тимур", "Timur"],
        ["Анна", "Anna"],
        ["A2 English · 18 дней подряд", "A2 English · 18 days in a row"],
        ["B1 Deutsch · 126 фраз аудирования", "B1 Deutsch · 126 listening phrases"],
        ["Spanish A1 · 420 слов", "Spanish A1 · 420 words"],
        ["20 языков сайта и бота", "20 website and bot languages"],
        ["Расширенное аудирование", "Expanded listening"],
        ["Максимальные лимиты на аудирование и голосовую оценку", "Maximum limits for listening and voice scoring"],
        ["Оценка произношения", "Pronunciation scoring"],
        ["Что даёт голосовая оценка произношения?", "What does voice pronunciation scoring provide?"],
        ["После голосового ответа пользователь получает понятную оценку, слабые слова, плавность речи, силу акцента и короткий совет. Это учебная обратная связь, которая помогает повторить фразу лучше уже со следующей попытки.", "After a voice answer, the user receives a clear score, weak words, speech flow, accent strength and a short tip. This learning feedback helps repeat the phrase better on the very next attempt."]
    );

    landingEnglishPhrases.forEach(([source, en, overrides]) => {
        const key = normalizePhrase(source);
        phraseTranslations[key] = { ...legalTranslations(en, overrides), ...(phraseTranslations[key] || {}) };
    });

    const publicSiteFinalOverrides = {
        "≈10 000 ₽/год": legalTranslations("≈10,000 RUB/year", {
            es: "≈10 000 ₽/año", de: "≈10.000 ₽/Jahr", fr: "≈10 000 ₽/an", it: "≈10.000 ₽/anno",
            zh: "≈10,000 ₽/年", ja: "約10,000 ₽/年", ko: "약 10,000 ₽/년", tg: "≈10 000 ₽/сол",
            uz: "≈10 000 ₽/yil", tt: "≈10 000 ₽/ел", hy: "≈10 000 ₽/տարի", kk: "≈10 000 ₽/жыл",
            ky: "≈10 000 ₽/жыл", ka: "≈10 000 ₽/წელი", uk: "≈10 000 ₽/рік", pl: "≈10 000 ₽/rok",
            ro: "≈10.000 ₽/an", pt: "≈10.000 ₽/ano"
        }),
        "≈20 000 ₽/год": legalTranslations("≈20,000 RUB/year", {
            es: "≈20 000 ₽/año", de: "≈20.000 ₽/Jahr", fr: "≈20 000 ₽/an", it: "≈20.000 ₽/anno",
            zh: "≈20,000 ₽/年", ja: "約20,000 ₽/年", ko: "약 20,000 ₽/년", tg: "≈20 000 ₽/сол",
            uz: "≈20 000 ₽/yil", tt: "≈20 000 ₽/ел", hy: "≈20 000 ₽/տարի", kk: "≈20 000 ₽/жыл",
            ky: "≈20 000 ₽/жыл", ka: "≈20 000 ₽/წელი", uk: "≈20 000 ₽/рік", pl: "≈20 000 ₽/rok",
            ro: "≈20.000 ₽/an", pt: "≈20.000 ₽/ano"
        }),
        "500 Stars · около 14 USDT": legalTranslations("500 Stars · about 14 USDT", {
            es: "500 Stars · unos 14 USDT", de: "500 Stars · ca. 14 USDT", fr: "500 Stars · environ 14 USDT", it: "500 Stars · circa 14 USDT",
            zh: "500 Stars · 约 14 USDT", ja: "500 Stars · 約14 USDT", ko: "500 Stars · 약 14 USDT", tg: "500 Stars · тақрибан 14 USDT",
            uz: "500 Stars · taxminan 14 USDT", tt: "500 Stars · якынча 14 USDT", hy: "500 Stars · մոտ 14 USDT", kk: "500 Stars · шамамен 14 USDT",
            ky: "500 Stars · болжол менен 14 USDT", ka: "500 Stars · დაახლოებით 14 USDT", uk: "500 Stars · близько 14 USDT", pl: "500 Stars · około 14 USDT",
            ro: "500 Stars · aprox. 14 USDT", pt: "500 Stars · cerca de 14 USDT"
        }),
        "150 Stars · около 4.15 USDT": legalTranslations("150 Stars · about 4.15 USDT", {
            es: "150 Stars · unos 4.15 USDT", de: "150 Stars · ca. 4.15 USDT", fr: "150 Stars · environ 4.15 USDT", it: "150 Stars · circa 4.15 USDT",
            zh: "150 Stars · 约 4.15 USDT", ja: "150 Stars · 約4.15 USDT", ko: "150 Stars · 약 4.15 USDT", tg: "150 Stars · тақрибан 4.15 USDT",
            uz: "150 Stars · taxminan 4.15 USDT", tt: "150 Stars · якынча 4.15 USDT", hy: "150 Stars · մոտ 4.15 USDT", kk: "150 Stars · шамамен 4.15 USDT",
            ky: "150 Stars · болжол менен 4.15 USDT", ka: "150 Stars · დაახლოებით 4.15 USDT", uk: "150 Stars · близько 4.15 USDT", pl: "150 Stars · około 4.15 USDT",
            ro: "150 Stars · aprox. 4.15 USDT", pt: "150 Stars · cerca de 4.15 USDT"
        }),
        "1 000 Stars · около 27 USDT": legalTranslations("1,000 Stars · about 27 USDT", {
            es: "1 000 Stars · unos 27 USDT", de: "1.000 Stars · ca. 27 USDT", fr: "1 000 Stars · environ 27 USDT", it: "1.000 Stars · circa 27 USDT",
            zh: "1,000 Stars · 约 27 USDT", ja: "1,000 Stars · 約27 USDT", ko: "1,000 Stars · 약 27 USDT", tg: "1 000 Stars · тақрибан 27 USDT",
            uz: "1 000 Stars · taxminan 27 USDT", tt: "1 000 Stars · якынча 27 USDT", hy: "1 000 Stars · մոտ 27 USDT", kk: "1 000 Stars · шамамен 27 USDT",
            ky: "1 000 Stars · болжол менен 27 USDT", ka: "1 000 Stars · დაახლოებით 27 USDT", uk: "1 000 Stars · близько 27 USDT", pl: "1 000 Stars · około 27 USDT",
            ro: "1.000 Stars · aprox. 27 USDT", pt: "1.000 Stars · cerca de 27 USDT"
        }),
        "300 Stars · около 8.15 USDT": legalTranslations("300 Stars · about 8.15 USDT", {
            es: "300 Stars · unos 8.15 USDT", de: "300 Stars · ca. 8.15 USDT", fr: "300 Stars · environ 8.15 USDT", it: "300 Stars · circa 8.15 USDT",
            zh: "300 Stars · 约 8.15 USDT", ja: "300 Stars · 約8.15 USDT", ko: "300 Stars · 약 8.15 USDT", tg: "300 Stars · тақрибан 8.15 USDT",
            uz: "300 Stars · taxminan 8.15 USDT", tt: "300 Stars · якынча 8.15 USDT", hy: "300 Stars · մոտ 8.15 USDT", kk: "300 Stars · шамамен 8.15 USDT",
            ky: "300 Stars · болжол менен 8.15 USDT", ka: "300 Stars · დაახლოებით 8.15 USDT", uk: "300 Stars · близько 8.15 USDT", pl: "300 Stars · około 8.15 USDT",
            ro: "300 Stars · aprox. 8.15 USDT", pt: "300 Stars · cerca de 8.15 USDT"
        }),
        "в месяц": legalTranslations("per month", {
            es: "al mes", de: "pro Monat", fr: "par mois", it: "al mese", zh: "每月", ja: "月額", ko: "월",
            tg: "дар як моҳ", uz: "oyiga", tt: "айга", hy: "ամսական", kk: "айына", ky: "айына", ka: "თვეში",
            uk: "на місяць", pl: "miesięcznie", ro: "pe lună", pt: "por mês"
        }),
        "либо @NERIVAapp_bot.": legalTranslations("or @NERIVAapp_bot.", {
            es: "o @NERIVAapp_bot.", de: "oder @NERIVAapp_bot.", fr: "ou @NERIVAapp_bot.", it: "o @NERIVAapp_bot.",
            zh: "或 @NERIVAapp_bot。", ja: "または @NERIVAapp_bot。", ko: "또는 @NERIVAapp_bot.", tg: "ё @NERIVAapp_bot.",
            uz: "yoki @NERIVAapp_bot.", tt: "яки @NERIVAapp_bot.", hy: "կամ @NERIVAapp_bot:", kk: "немесе @NERIVAapp_bot.",
            ky: "же @NERIVAapp_bot.", ka: "ან @NERIVAapp_bot.", uk: "або @NERIVAapp_bot.", pl: "albo @NERIVAapp_bot.",
            ro: "sau @NERIVAapp_bot.", pt: "ou @NERIVAapp_bot."
        }),
        "Интерфейс: RU": legalTranslations("Interface: RU", {
            es: "Interfaz: RU", de: "Oberfläche: RU", fr: "Interface : RU", it: "Interfaccia: RU", zh: "界面：RU", ja: "インターフェース: RU",
            ko: "인터페이스: RU", tg: "Забони UI: RU", uz: "Interfeys: RU", tt: "UI теле: RU", hy: "Ինտերֆեյս՝ RU", kk: "UI тілі: RU",
            ky: "UI тили: RU", ka: "ინტერფეისი: RU", uk: "Інтерфейс: RU", pl: "Interfejs: RU", ro: "Interfață: RU", pt: "Interface: RU"
        }),
        "Интерфейс: UZ": legalTranslations("Interface: UZ", {
            es: "Interfaz: UZ", de: "Oberfläche: UZ", fr: "Interface : UZ", it: "Interfaccia: UZ", zh: "界面：UZ", ja: "インターフェース: UZ",
            ko: "인터페이스: UZ", tg: "Забони UI: UZ", uz: "Interfeys: UZ", tt: "UI теле: UZ", hy: "Ինտերֆեյս՝ UZ", kk: "UI тілі: UZ",
            ky: "UI тили: UZ", ka: "ინტერფეისი: UZ", uk: "Інтерфейс: UZ", pl: "Interfejs: UZ", ro: "Interfață: UZ", pt: "Interface: UZ"
        }),
        "Telegram-бот": legalTranslations("Telegram bot", {
            es: "bot de Telegram", de: "Telegram-Bot", fr: "bot Telegram", it: "bot Telegram", zh: "Telegram 机器人", ja: "Telegramボット",
            ko: "Telegram 봇", tg: "боти Telegram", uz: "Telegram bot", tt: "Telegram боты", hy: "Telegram բոտ", kk: "Telegram боты",
            ky: "Telegram боту", ka: "Telegram ბოტი", uk: "Telegram-бот", pl: "bot Telegram", ro: "bot Telegram", pt: "bot do Telegram"
        }),
        "Telegram бот": legalTranslations("Telegram bot", {
            es: "bot de Telegram", de: "Telegram-Bot", fr: "bot Telegram", it: "bot Telegram", zh: "Telegram 机器人", ja: "Telegramボット",
            ko: "Telegram 봇", tg: "боти Telegram", uz: "Telegram bot", tt: "Telegram боты", hy: "Telegram բոտ", kk: "Telegram боты",
            ky: "Telegram боту", ka: "Telegram ბოტი", uk: "Telegram-бот", pl: "bot Telegram", ro: "bot Telegram", pt: "bot do Telegram"
        }),
        "Онлайн:": legalTranslations("Online:", { tg: "Онлайн:", uz: "Onlayn:", tt: "Онлайн:", kk: "Онлайн:", ky: "Онлайн:", uk: "Онлайн:" }),
        "Онлайн": legalTranslations("Online", { tg: "Онлайн", uz: "Onlayn", tt: "Онлайн", kk: "Онлайн", ky: "Онлайн", uk: "Онлайн" }),
        "Бот:": legalTranslations("Bot:", { es: "Bot:", de: "Bot:", fr: "Bot :", it: "Bot:", zh: "机器人：", ja: "ボット:", ko: "봇:", tg: "Бот:", uz: "Bot:", tt: "Бот:", kk: "Бот:", ky: "Бот:", ka: "ბოტი:", uk: "Бот:", pl: "Bot:", ro: "Bot:", pt: "Bot:" }),
        "Бот": legalTranslations("Bot", { es: "Bot", de: "Bot", fr: "Bot", it: "Bot", zh: "机器人", ja: "ボット", ko: "봇", tg: "Бот", uz: "Bot", tt: "Бот", kk: "Бот", ky: "Бот", ka: "ბოტი", uk: "Бот", pl: "Bot", ro: "Bot", pt: "Bot" }),
        "Оператор": legalTranslations("Operator", { es: "Operador", de: "Betreiber", fr: "Opérateur", it: "Operatore", zh: "运营方", ja: "運営者", ko: "운영자", tg: "Оператор", uz: "Operator", tt: "Оператор", hy: "Օպերատոր", kk: "Оператор", ky: "Оператор", ka: "ოპერატორი", uk: "Оператор", pl: "Operator", ro: "Operator", pt: "Operador" }),
        "Стоимость": legalTranslations("Price", { es: "Precio", de: "Preis", fr: "Prix", it: "Prezzo", zh: "价格", ja: "料金", ko: "가격", tg: "Нарх", uz: "Narx", tt: "Бәя", hy: "Գին", kk: "Баға", ky: "Баасы", ka: "ფასი", uk: "Вартість", pl: "Cena", ro: "Preț", pt: "Preço" }),
        "AI-переводчик": legalTranslations("AI translator", { es: "traductor IA", de: "KI-Übersetzer", fr: "traducteur IA", it: "traduttore IA", zh: "AI 翻译器", ja: "AI翻訳", ko: "AI 번역기", tg: "тарҷумони AI", uz: "AI tarjimon", tt: "AI тәрҗемәче", hy: "AI թարգմանիչ", kk: "AI аудармашы", ky: "AI котормочу", ka: "AI მთარგმნელი", uk: "AI-перекладач", pl: "tłumacz AI", ro: "traducător AI", pt: "tradutor IA" })
    };

    Object.assign(publicSiteFinalOverrides, {
        "казахский": { fr: "kazakh", ky: "казакча" },
        "французский": { ky: "французча" },
        "Администрация": {
            fr: "Administration du service", tt: "Администрациясе", ky: "Администрациясы"
        },
        "12. Контакты": { fr: "12. Coordonnées" },
        "Навигация": { kk: "Бағдарлау" },
        "Оператор": { pl: "Operator usługi" },
        "AI-практика": { ky: "AI практикасы", uk: "Практика з AI" },
        "AI-репетитор, а не просто чат": {
            uk: "AI-наставник, а не просто чат"
        },
        "Онлайн-бот": { uk: "Бот онлайн" },
        "Telegram-бот": {
            ...publicSiteFinalOverrides["Telegram-бот"],
            uk: "бот Telegram"
        },
        "20 голосовых в день до 30 секунд": {
            en: "20 voice messages per day up to 30 seconds"
        },
        "50 уроков в день, 200 сообщений практики, 20 голосовых до 30 секунд, голос в текст, перевод услышанного, перевод текста с картинки, практика по голосу или фото.": {
            en: "50 lessons per day, 200 practice messages, 20 voice messages up to 30 seconds, voice to text, heard speech translation, image text translation and practice from voice or photo context."
        },
        "Тариф": { uk: "План" },
        "Интерфейс: RU": {
            ...publicSiteFinalOverrides["Интерфейс: RU"],
            pt: "Idioma da interface: RU"
        },
        "Интерфейс: UZ": {
            ...publicSiteFinalOverrides["Интерфейс: UZ"],
            pt: "Idioma da interface: UZ"
        }
    });

    Object.entries(publicSiteFinalOverrides).forEach(([source, values]) => {
        const key = normalizePhrase(source);
        phraseTranslations[key] = { ...(phraseTranslations[key] || {}), ...values };
    });

    const textNodeSources = new WeakMap();
    const translatedAttributes = ["alt", "aria-label", "title", "placeholder", "content"];
    const skippedTextParents = new Set(["SCRIPT", "STYLE", "NOSCRIPT", "TEXTAREA", "SELECT", "OPTION"]);
    const supported = new Set(languages.map(([code]) => code));

    function normalizeLanguage(value) {
        const code = String(value || "").toLowerCase().trim();
        if (code === "rus") return "ru";
        if (code === "rom") return "ro";
        if (code === "por") return "pt";
        if (code === "kir") return "ky";
        if (code === "geo" || code === "kat") return "ka";
        return supported.has(code) ? code : "ru";
    }

    function currentLanguage() {
        const params = new URLSearchParams(location.search);
        return normalizeLanguage(params.get("lang") || localStorage.getItem("poliglot_site_language") || document.documentElement.lang || "ru");
    }

    function t(key, fallback = "") {
        const lang = currentLanguage();
        const ruText = copy.ru[key] || fallback || key;
        const localized = copy[lang]?.[key];
        if (lang === "ru") return localized || ruText;
        if (lang === "en" && localized) return localized;
        if (localized && localized !== copy.en[key]) return localized;
        return phraseTranslations[normalizePhrase(ruText)]?.[lang] || localized || copy.en[key] || ruText;
    }

    function normalizePhrase(value) {
        return String(value || "").replace(/\s+/g, " ").trim();
    }

    function mojibakeKey(value) {
        try {
            return Array.from(new TextEncoder().encode(value)).map((byte) => String.fromCharCode(byte)).join("");
        } catch (_) {
            return value;
        }
    }

    function decodeMaybeMojibake(value) {
        const text = String(value || "");
        if (!/[ÐÑ]|[\u0080-\u009F]/.test(text)) return text;
        try {
            const bytes = Uint8Array.from(Array.from(text).map((char) => char.charCodeAt(0) & 255));
            return new TextDecoder("utf-8").decode(bytes);
        } catch (_) {
            return text;
        }
    }

    function phraseEntryFor(value) {
        const normalized = normalizePhrase(value);
        const decoded = normalizePhrase(decodeMaybeMojibake(normalized));
        return phraseTranslations[normalized]
            || phraseTranslations[decoded]
            || phraseTranslations[mojibakeKey(normalized)]
            || phraseTranslations[mojibakeKey(decoded)];
    }

    function translatePhrase(source, fallback = "") {
        const lang = currentLanguage();
        const normalized = normalizePhrase(source);
        const decoded = normalizePhrase(decodeMaybeMojibake(normalized));
        if (!normalized) return fallback || source || "";
        const entry = phraseEntryFor(normalized);
        if (lang === "ru") {
            const translated = entry?.ru;
            return translated ? decodeMaybeMojibake(translated) : fallback || decoded || normalized;
        }
        const translated = entry?.[lang] || entry?.en;
        return translated ? decodeMaybeMojibake(translated) : fallback || decoded || normalized;
    }

    function hasPhraseTranslation(normalized) {
        return Boolean(phraseEntryFor(normalized));
    }

    function translateTextNode(node) {
        const raw = node.nodeValue || "";
        const normalized = normalizePhrase(raw);
        if (!normalized) return;
        if (node.parentElement?.closest("[data-i18n], [data-email-text], [data-no-translate], [data-site-language-select]")) return;
        if (skippedTextParents.has(node.parentElement?.tagName || "")) return;
        if (currentLanguage() === "ru" && node.parentElement?.closest(".legal-document-shell")) {
            const source = textNodeSources.get(node);
            if (source) {
                const leading = raw.match(/^\s*/)?.[0] || "";
                const trailing = raw.match(/\s*$/)?.[0] || "";
                node.nodeValue = `${leading}${source}${trailing}`;
            }
            return;
        }
        if (!textNodeSources.has(node)) {
            if (!hasPhraseTranslation(normalized)) return;
            textNodeSources.set(node, normalizePhrase(decodeMaybeMojibake(normalized)) || normalized);
        }
        const source = textNodeSources.get(node);
        const leading = raw.match(/^\s*/)?.[0] || "";
        const trailing = raw.match(/\s*$/)?.[0] || "";
        node.nodeValue = `${leading}${translatePhrase(source)}${trailing}`;
    }

    function translatePlainText() {
        const walker = document.createTreeWalker(document.body || document.documentElement, NodeFilter.SHOW_TEXT);
        const nodes = [];
        while (walker.nextNode()) nodes.push(walker.currentNode);
        nodes.forEach(translateTextNode);
    }

    function translateAttribute(node, attr) {
        const raw = node.getAttribute(attr);
        const normalized = normalizePhrase(raw);
        if (!normalized || node.closest?.("[data-no-translate], [data-site-language-select]")) return;
        if (attr === "content" && node.tagName === "META" && normalizePhrase(node.getAttribute("name")) === "viewport") return;
        node.__poliglotOriginalAttrs ||= {};
        if (currentLanguage() === "ru" && node.closest?.(".legal-document-shell")) {
            if (node.__poliglotOriginalAttrs[attr]) {
                node.setAttribute(attr, node.__poliglotOriginalAttrs[attr]);
            }
            return;
        }
        if (!node.__poliglotOriginalAttrs[attr]) {
            if (!hasPhraseTranslation(normalized)) return;
            node.__poliglotOriginalAttrs[attr] = normalizePhrase(decodeMaybeMojibake(normalized)) || normalized;
        }
        node.setAttribute(attr, translatePhrase(node.__poliglotOriginalAttrs[attr]));
    }

    function translateAttributes() {
        const selector = translatedAttributes.map((attr) => `[${attr}]`).join(",");
        document.querySelectorAll(selector).forEach((node) => {
            translatedAttributes.forEach((attr) => {
                if (node.hasAttribute(attr)) translateAttribute(node, attr);
            });
        });
    }

    function localizeHref(rawHref, lang) {
        if (!rawHref || rawHref.startsWith("#") || rawHref.startsWith("mailto:") || rawHref.startsWith("tel:")) return rawHref;
        if (/^https?:\/\//i.test(rawHref) && !rawHref.includes(location.host)) return rawHref;
        const url = new URL(rawHref, location.href);
        if (!url.pathname.endsWith(".html") && url.pathname !== "/" && url.pathname !== "/poliglot-ai.html") return rawHref;
        url.searchParams.set("lang", lang);
        return url.pathname.replace(/^\//, "") + url.search + url.hash;
    }

    function apply() {
        const lang = currentLanguage();
        document.documentElement.lang = lang;
        localStorage.setItem("poliglot_site_language", lang);
        document.querySelectorAll("[data-site-language-select]").forEach((select) => {
            if (!select.options.length) {
                select.innerHTML = languages.map(([code, label]) => `<option value="${code}">${label}</option>`).join("");
            }
            select.value = lang;
            select.setAttribute("aria-label", t("language_label", "Language"));
        });
        document.querySelectorAll("[data-i18n]").forEach((node) => {
            node.textContent = t(node.dataset.i18n, node.textContent);
        });
        translatePlainText();
        translateAttributes();
        document.querySelectorAll("a[href]").forEach((link) => {
            if (link.hasAttribute("data-no-localize-href")) return;
            const original = link.dataset.originalHref || link.getAttribute("href");
            link.dataset.originalHref = original;
            link.setAttribute("href", localizeHref(original, lang));
        });
        const page = document.body?.dataset?.page || "";
        if (page && t(`${page}_title`, "")) {
            document.title = `${t(`${page}_title`)} — NERIVA`;
        }
    }

    document.addEventListener("change", (event) => {
        if (!event.target.matches("[data-site-language-select]")) return;
        const lang = normalizeLanguage(event.target.value);
        localStorage.setItem("poliglot_site_language", lang);
        const url = new URL(location.href);
        url.searchParams.set("lang", lang);
        history.replaceState(null, "", url);
        apply();
        window.dispatchEvent(new CustomEvent("poliglot-language-change", { detail: { lang } }));
    });

    document.addEventListener("DOMContentLoaded", apply);
    window.poliglotSiteI18n = { languages, t, translatePhrase, apply, currentLanguage, localizeHref };
})();
