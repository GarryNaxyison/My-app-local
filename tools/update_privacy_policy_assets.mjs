import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const locales = ["ru", "en", "es", "de", "fr", "it", "zh", "ja", "ko", "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt"];
const contacts = {
  email: "support@neriva.ru",
  admin: "@AsaselD",
  bot: "@NERIVAapp_bot",
};

const policy = {
  ru: {
    title: "ПОЛИТИКА ОБРАБОТКИ ПЕРСОНАЛЬНЫХ ДАННЫХ (NERIVA)",
    updated: "Обновлено 27 мая 2026",
    intro: "Настоящая Политика действует в отношении сервиса NERIVA, доступного через сайты poliglotai.ru, poliglotai.online и Telegram-бот @NERIVAapp_bot. Сервис представляет собой единую систему независимо от способа доступа. Используя сервис, Пользователь подтверждает согласие с данной Политикой.",
    sections: [
      ["general", "1. Общие положения", ["Настоящая Политика действует в отношении сервиса NERIVA, доступного через сайты poliglotai.ru, poliglotai.online и Telegram-бот @NERIVAapp_bot.", "Сервис представляет собой единую систему, независимо от способа доступа.", "Используя сервис, Пользователь подтверждает согласие с данной Политикой."]],
      ["operator", "2. Оператор персональных данных", ["Оператором персональных данных является владелец сервиса NERIVA (далее — “Оператор”)."]],
      ["data", "3. Какие данные обрабатываются", ["Мы можем обрабатывать:", ["идентификатор пользователя в Telegram (Telegram ID)", "имя пользователя (username, если доступен)", "технические данные (IP-адрес, cookies, информация об устройстве)"], "Мы не обрабатываем паспортные данные, платёжные реквизиты и другие специальные категории данных."],
      ],
      ["sources", "4. Источники данных", ["Данные предоставляются:", ["через авторизацию в Telegram", "при использовании сайта или Telegram-бота", "автоматически при использовании сервиса"]]],
      ["purposes", "5. Цели обработки", ["Данные используются для:", ["создания и поддержки аккаунта", "предоставления доступа к сервису", "обеспечения работы сайта и бота", "улучшения качества сервиса"]]],
      ["storage", "6. Хранение и передача данных", ["Данные хранятся на серверах, расположенных за пределами Российской Федерации (включая Нидерланды).", "Данные могут обрабатываться с использованием сторонних технических сервисов, необходимых для работы инфраструктуры."]],
      ["protection", "7. Защита данных", ["Оператор принимает разумные технические и организационные меры для защиты данных от несанкционированного доступа и утечек."]],
      ["retention", "8. Срок хранения", ["Данные хранятся до удаления аккаунта пользователем или прекращения работы сервиса."]],
      ["rights", "9. Права пользователя", ["Пользователь имеет право:", ["запросить доступ к данным", "исправить или удалить данные", "отозвать согласие на обработку"]]],
      ["consent", "10. Согласие", ["Используя сервис или проходя авторизацию через Telegram, Пользователь выражает согласие с данной Политикой обработки персональных данных."]],
      ["contacts", "11. Контакты", [`Email: ${contacts.email}`, `Telegram: ${contacts.admin}`, `Bot: ${contacts.bot}`]],
    ],
  },
  en: {
    title: "PERSONAL DATA PROCESSING POLICY (NERIVA)",
    updated: "Updated 27 May 2026",
    intro: "This Policy applies to NERIVA on poliglotai.ru, poliglotai.online, and the Telegram bot @NERIVAapp_bot. The service is a single system regardless of access method. By using it, the User agrees to this Policy.",
    sections: [
      ["general", "1. General provisions", ["This Policy applies to NERIVA on poliglotai.ru, poliglotai.online, and the Telegram bot @NERIVAapp_bot.", "The service is a single system regardless of access method.", "By using the service, the User agrees to this Policy."]],
      ["operator", "2. Data operator", ["The personal data operator is the owner of the NERIVA service (the “Operator”)."]],
      ["data", "3. Data we process", ["We may process:", ["Telegram user identifier (Telegram ID)", "username, if available", "technical data such as IP address, cookies, and device information"], "We do not process passport data, payment requisites, or special categories of personal data."]],
      ["sources", "4. Data sources", ["Data is provided:", ["through Telegram authorization", "when using the website or Telegram bot", "automatically while using the service"]]],
      ["purposes", "5. Processing purposes", ["Data is used to:", ["create and maintain an account", "provide access to the service", "operate the website and bot", "improve service quality"]]],
      ["storage", "6. Storage and transfer", ["Data is stored on servers located outside the Russian Federation, including the Netherlands.", "Data may be processed with third-party technical services required for infrastructure operation."]],
      ["protection", "7. Data protection", ["The Operator applies reasonable technical and organizational measures to protect data from unauthorized access and leaks."]],
      ["retention", "8. Retention period", ["Data is stored until the user deletes the account or the service stops operating."]],
      ["rights", "9. User rights", ["The User has the right to:", ["request access to data", "correct or delete data", "withdraw consent to processing"]]],
      ["consent", "10. Consent", ["By using the service or authorizing through Telegram, the User consents to this Personal Data Processing Policy."]],
      ["contacts", "11. Contacts", [`Email: ${contacts.email}`, `Telegram: ${contacts.admin}`, `Bot: ${contacts.bot}`]],
    ],
  },
};

const localized = {
  es: ["POLÍTICA DE TRATAMIENTO DE DATOS PERSONALES (NERIVA)", "Actualizado el 27 de mayo de 2026", "Esta Política se aplica a NERIVA en poliglotai.ru, poliglotai.online y al bot de Telegram @NERIVAapp_bot. El servicio es un sistema único sin importar el modo de acceso. Al usarlo, el Usuario acepta esta Política."],
  de: ["RICHTLINIE ZUR VERARBEITUNG PERSONENBEZOGENER DATEN (NERIVA)", "Aktualisiert am 27. Mai 2026", "Diese Richtlinie gilt für NERIVA auf poliglotai.ru, poliglotai.online und den Telegram-Bot @NERIVAapp_bot. Der Dienst ist unabhängig vom Zugang eine einheitliche Plattform. Mit der Nutzung stimmt der Nutzer dieser Richtlinie zu."],
  fr: ["POLITIQUE DE TRAITEMENT DES DONNÉES PERSONNELLES (NERIVA)", "Mise à jour le 27 mai 2026", "La présente Politique s’applique à NERIVA sur poliglotai.ru, poliglotai.online et au bot Telegram @NERIVAapp_bot. Le service constitue un système unique, quel que soit le mode d’accès. En l’utilisant, l’Utilisateur accepte cette Politique."],
  it: ["POLITICA DI TRATTAMENTO DEI DATI PERSONALI (NERIVA)", "Aggiornata il 27 maggio 2026", "La presente Politica si applica a NERIVA su poliglotai.ru, poliglotai.online e al bot Telegram @NERIVAapp_bot. Il servizio è un sistema unico indipendentemente dal metodo di accesso. Usandolo, l’Utente accetta questa Politica."],
  zh: ["个人数据处理政策 (NERIVA)", "更新日期：2026年5月27日", "本政策适用于 poliglotai.ru、poliglotai.online 以及 Telegram 机器人 @NERIVAapp_bot 上的 NERIVA。无论访问方式如何，服务均为统一系统。使用服务即表示用户同意本政策。"],
  ja: ["個人データ処理ポリシー (NERIVA)", "2026年5月27日更新", "本ポリシーは、poliglotai.ru、poliglotai.online、およびTelegramボット @NERIVAapp_bot のNERIVAに適用されます。サービスはアクセス方法にかかわらず単一のシステムです。利用によりユーザーは本ポリシーに同意します。"],
  ko: ["개인정보 처리 정책 (NERIVA)", "2026년 5월 27일 업데이트", "본 정책은 poliglotai.ru, poliglotai.online 및 Telegram 봇 @NERIVAapp_bot의 NERIVA에 적용됩니다. 서비스는 접속 방식과 관계없이 하나의 시스템입니다. 사용자는 서비스를 이용함으로써 본 정책에 동의합니다."],
  tg: ["СИЁСАТИ КОРКАРДИ МАЪЛУМОТИ ШАХСӢ (NERIVA)", "Навсозӣ: 27 майи 2026", "Ин Сиёсат ба NERIVA дар poliglotai.ru, poliglotai.online ва боти Telegram @NERIVAapp_bot татбиқ мешавад. Сервис новобаста аз роҳи дастрасӣ як низоми ягона аст. Бо истифода аз сервис, Корбар ба ин Сиёсат розӣ мешавад."],
  uz: ["SHAXSIY MA’LUMOTLARNI QAYTA ISHLASH SIYOSATI (NERIVA)", "2026-yil 27-mayda yangilangan", "Ushbu Siyosat poliglotai.ru, poliglotai.online va Telegram boti @NERIVAapp_bot orqali ishlaydigan NERIVA xizmatiga taalluqlidir. Xizmat kirish usulidan qat’i nazar yagona tizimdir. Foydalanish orqali Foydalanuvchi ushbu Siyosatga rozilik bildiradi."],
  tt: ["ШӘХСИ МӘГЪЛҮМАТЛАРНЫ ЭШКӘРТҮ СӘЯСӘТЕ (NERIVA)", "2026 елның 27 маенда яңартылды", "Әлеге Сәясәт poliglotai.ru, poliglotai.online һәм Telegram-бот @NERIVAapp_bot аша эшләүче NERIVA сервисына кагыла. Сервис керү ысулына карамастан бердәм система булып тора. Сервистан файдалану белән Кулланучы әлеге Сәясәт белән килешә."],
  hy: ["ԱՆՁՆԱԿԱՆ ՏՎՅԱԼՆԵՐԻ ՄՇԱԿՄԱՆ ՔԱՂԱՔԱԿԱՆՈՒԹՅՈՒՆ (NERIVA)", "Թարմացվել է 2026 թվականի մայիսի 27-ին", "Սույն Քաղաքականությունը վերաբերում է NERIVA ծառայությանը poliglotai.ru, poliglotai.online կայքերում և Telegram բոտում՝ @NERIVAapp_bot։ Ծառայությունը մեկ միասնական համակարգ է՝ անկախ մուտքի եղանակից։ Օգտագործելով ծառայությունը՝ Օգտատերը համաձայնում է սույն Քաղաքականությանը։"],
  kk: ["ЖЕКЕ ДЕРЕКТЕРДІ ӨҢДЕУ САЯСАТЫ (NERIVA)", "2026 жылғы 27 мамырда жаңартылды", "Осы Саясат poliglotai.ru, poliglotai.online сайттарындағы және @NERIVAapp_bot Telegram-ботындағы NERIVA сервисіне қолданылады. Сервис қол жеткізу тәсіліне қарамастан біртұтас жүйе болып табылады. Сервисті пайдалану арқылы Пайдаланушы осы Саясатпен келіседі."],
  ky: ["ЖЕКЕ МААЛЫМАТТАРДЫ ИШТЕТҮҮ САЯСАТЫ (NERIVA)", "2026-жылдын 27-майында жаңыртылды", "Бул Саясат poliglotai.ru, poliglotai.online сайттарындагы жана @NERIVAapp_bot Telegram ботундагы NERIVA сервисине колдонулат. Сервис кирүү ыкмасына карабастан бирдиктүү система болуп саналат. Сервисти колдонуу менен Колдонуучу ушул Саясатка макул болот."],
  ka: ["პერსონალური მონაცემების დამუშავების პოლიტიკა (NERIVA)", "განახლდა 2026 წლის 27 მაისს", "ეს პოლიტიკა ვრცელდება NERIVA სერვისზე poliglotai.ru, poliglotai.online და Telegram ბოტში @NERIVAapp_bot. სერვისი არის ერთიანი სისტემა, წვდომის მეთოდის მიუხედავად. სერვისის გამოყენებით მომხმარებელი ეთანხმება ამ პოლიტიკას."],
  uk: ["ПОЛІТИКА ОБРОБКИ ПЕРСОНАЛЬНИХ ДАНИХ (NERIVA)", "Оновлено 27 травня 2026", "Ця Політика діє щодо сервісу NERIVA, доступного через сайти poliglotai.ru, poliglotai.online і Telegram-бот @NERIVAapp_bot. Сервіс є єдиною системою незалежно від способу доступу. Використовуючи сервіс, Користувач погоджується з цією Політикою."],
  pl: ["POLITYKA PRZETWARZANIA DANYCH OSOBOWYCH (NERIVA)", "Zaktualizowano 27 maja 2026", "Niniejsza Polityka dotyczy NERIVA w serwisach poliglotai.ru, poliglotai.online oraz w bocie Telegram @NERIVAapp_bot. Usługa jest jednym systemem niezależnie od sposobu dostępu. Korzystając z usługi, Użytkownik akceptuje tę Politykę."],
  ro: ["POLITICA DE PRELUCRARE A DATELOR CU CARACTER PERSONAL (NERIVA)", "Actualizat la 27 mai 2026", "Această Politică se aplică serviciului NERIVA disponibil prin poliglotai.ru, poliglotai.online și botul Telegram @NERIVAapp_bot. Serviciul este un sistem unic indiferent de metoda de acces. Prin utilizare, Utilizatorul acceptă această Politică."],
  pt: ["POLÍTICA DE TRATAMENTO DE DADOS PESSOAIS (NERIVA)", "Atualizada em 27 de maio de 2026", "Esta Política aplica-se ao NERIVA em poliglotai.ru, poliglotai.online e ao bot do Telegram @NERIVAapp_bot. O serviço é um sistema único independentemente do método de acesso. Ao utilizar o serviço, o Utilizador aceita esta Política."],
};

const labels = {
  es: ["1. Disposiciones generales", "2. Operador de datos", "3. Datos tratados", "4. Fuentes de datos", "5. Finalidades", "6. Almacenamiento y transferencia", "7. Protección", "8. Plazo de conservación", "9. Derechos del usuario", "10. Consentimiento", "11. Contactos"],
  de: ["1. Allgemeine Bestimmungen", "2. Datenbetreiber", "3. Verarbeitete Daten", "4. Datenquellen", "5. Zwecke", "6. Speicherung und Übermittlung", "7. Schutz", "8. Aufbewahrungsdauer", "9. Rechte des Nutzers", "10. Einwilligung", "11. Kontakte"],
  fr: ["1. Dispositions générales", "2. Responsable du traitement", "3. Données traitées", "4. Sources des données", "5. Finalités", "6. Conservation et transfert", "7. Protection", "8. Durée de conservation", "9. Droits de l’utilisateur", "10. Consentement", "11. Contacts"],
  it: ["1. Disposizioni generali", "2. Titolare del trattamento", "3. Dati trattati", "4. Fonti dei dati", "5. Finalità", "6. Conservazione e trasferimento", "7. Protezione", "8. Periodo di conservazione", "9. Diritti dell’utente", "10. Consenso", "11. Contatti"],
  zh: ["1. 一般条款", "2. 数据运营者", "3. 处理的数据", "4. 数据来源", "5. 处理目的", "6. 存储与传输", "7. 数据保护", "8. 保存期限", "9. 用户权利", "10. 同意", "11. 联系方式"],
  ja: ["1. 一般規定", "2. データ管理者", "3. 処理するデータ", "4. データの取得元", "5. 処理目的", "6. 保存と移転", "7. データ保護", "8. 保存期間", "9. ユーザーの権利", "10. 同意", "11. 連絡先"],
  ko: ["1. 일반 조항", "2. 데이터 운영자", "3. 처리하는 데이터", "4. 데이터 출처", "5. 처리 목적", "6. 저장 및 이전", "7. 보호", "8. 보관 기간", "9. 사용자 권리", "10. 동의", "11. 연락처"],
  tg: ["1. Муқаррароти умумӣ", "2. Оператори маълумот", "3. Маълумоти коркардшаванда", "4. Манбаъҳои маълумот", "5. Ҳадафҳо", "6. Нигоҳдорӣ ва интиқол", "7. Муҳофизат", "8. Муҳлати нигоҳдорӣ", "9. Ҳуқуқҳои корбар", "10. Розигӣ", "11. Тамос"],
  uz: ["1. Umumiy qoidalar", "2. Ma’lumotlar operatori", "3. Qayta ishlanadigan ma’lumotlar", "4. Ma’lumot manbalari", "5. Maqsadlar", "6. Saqlash va uzatish", "7. Himoya", "8. Saqlash muddati", "9. Foydalanuvchi huquqlari", "10. Rozilik", "11. Aloqa"],
  tt: ["1. Гомуми нигезләмәләр", "2. Мәгълүмат операторы", "3. Эшкәртелә торган мәгълүмат", "4. Мәгълүмат чыганаклары", "5. Максатлар", "6. Саклау һәм тапшыру", "7. Саклау чаралары", "8. Саклау вакыты", "9. Кулланучы хокуклары", "10. Ризалык", "11. Контактлар"],
  hy: ["1. Ընդհանուր դրույթներ", "2. Տվյալների օպերատոր", "3. Մշակվող տվյալներ", "4. Տվյալների աղբյուրներ", "5. Նպատակներ", "6. Պահպանում և փոխանցում", "7. Պաշտպանություն", "8. Պահպանման ժամկետ", "9. Օգտատիրոջ իրավունքներ", "10. Համաձայնություն", "11. Կոնտակտներ"],
  kk: ["1. Жалпы ережелер", "2. Деректер операторы", "3. Өңделетін деректер", "4. Деректер көздері", "5. Мақсаттар", "6. Сақтау және беру", "7. Қорғау", "8. Сақтау мерзімі", "9. Пайдаланушы құқықтары", "10. Келісім", "11. Байланыс"],
  ky: ["1. Жалпы жоболор", "2. Маалымат оператору", "3. Иштетилген маалыматтар", "4. Маалымат булактары", "5. Максаттар", "6. Сактоо жана берүү", "7. Коргоо", "8. Сактоо мөөнөтү", "9. Колдонуучунун укуктары", "10. Макулдук", "11. Байланыш"],
  ka: ["1. ზოგადი დებულებები", "2. მონაცემთა ოპერატორი", "3. დამუშავებული მონაცემები", "4. მონაცემთა წყაროები", "5. მიზნები", "6. შენახვა და გადაცემა", "7. დაცვა", "8. შენახვის ვადა", "9. მომხმარებლის უფლებები", "10. თანხმობა", "11. კონტაქტები"],
  uk: ["1. Загальні положення", "2. Оператор даних", "3. Які дані обробляються", "4. Джерела даних", "5. Цілі обробки", "6. Зберігання та передача", "7. Захист даних", "8. Строк зберігання", "9. Права користувача", "10. Згода", "11. Контакти"],
  pl: ["1. Postanowienia ogólne", "2. Operator danych", "3. Przetwarzane dane", "4. Źródła danych", "5. Cele", "6. Przechowywanie i przekazywanie", "7. Ochrona", "8. Okres przechowywania", "9. Prawa użytkownika", "10. Zgoda", "11. Kontakt"],
  ro: ["1. Dispoziții generale", "2. Operatorul datelor", "3. Date prelucrate", "4. Surse de date", "5. Scopuri", "6. Stocare și transfer", "7. Protecție", "8. Perioada de păstrare", "9. Drepturile utilizatorului", "10. Consimțământ", "11. Contacte"],
  pt: ["1. Disposições gerais", "2. Operador dos dados", "3. Dados tratados", "4. Fontes de dados", "5. Finalidades", "6. Armazenamento e transferência", "7. Proteção", "8. Prazo de conservação", "9. Direitos do utilizador", "10. Consentimento", "11. Contactos"],
};

const localizedSectionBlocks = {
  es: [
    ["Esta Política se aplica al servicio NERIVA en poliglotai.ru, poliglotai.online y al bot de Telegram @NERIVAapp_bot.", "El servicio funciona como un sistema único sin importar el modo de acceso.", "Al usar el servicio, el Usuario acepta esta Política."],
    ["El operador de datos personales es el propietario del servicio NERIVA."],
    ["Podemos tratar Telegram ID, username si está disponible, IP, cookies e información del dispositivo.", "No tratamos datos de pasaporte, datos de pago ni categorías especiales de datos."],
    ["Los datos se reciben mediante autorización en Telegram, uso del sitio o bot, y automáticamente durante el uso del servicio."],
    ["Los datos se usan para crear y mantener la cuenta, dar acceso al servicio, operar el sitio y el bot, y mejorar la calidad del servicio."],
    ["Los datos se almacenan en servidores fuera de la Federación Rusa, incluidos los Países Bajos.", "También pueden ser tratados por servicios técnicos externos necesarios para la infraestructura."],
    ["El Operador aplica medidas técnicas y organizativas razonables contra accesos no autorizados y fugas."],
    ["Los datos se conservan hasta que el usuario elimine la cuenta o el servicio deje de funcionar."],
    ["El Usuario puede solicitar acceso a los datos, corregirlos o eliminarlos, y retirar el consentimiento."],
    ["Al usar el servicio o autorizarse por Telegram, el Usuario consiente esta Política de tratamiento de datos personales."],
  ],
  de: [
    ["Diese Richtlinie gilt für NERIVA auf poliglotai.ru, poliglotai.online und für den Telegram-Bot @NERIVAapp_bot.", "Der Dienst ist unabhängig vom Zugang eine einheitliche Plattform.", "Mit der Nutzung stimmt der Nutzer dieser Richtlinie zu."],
    ["Betreiber personenbezogener Daten ist der Eigentümer des Dienstes NERIVA."],
    ["Wir können Telegram ID, username sofern verfügbar, IP-Adresse, Cookies und Geräteinformationen verarbeiten.", "Passdaten, Zahlungsdaten und besondere Kategorien personenbezogener Daten verarbeiten wir nicht."],
    ["Daten werden über Telegram-Autorisierung, Nutzung der Website oder des Bots und automatisch während der Nutzung bereitgestellt."],
    ["Daten dienen der Kontoerstellung, dem Zugriff auf den Dienst, dem Betrieb von Website und Bot sowie der Qualitätsverbesserung."],
    ["Daten werden auf Servern außerhalb der Russischen Föderation gespeichert, einschließlich der Niederlande.", "Erforderliche technische Drittanbieter können an der Verarbeitung beteiligt sein."],
    ["Der Betreiber trifft angemessene technische und organisatorische Schutzmaßnahmen gegen unbefugten Zugriff und Datenlecks."],
    ["Daten werden bis zur Löschung des Kontos oder bis zur Einstellung des Dienstes gespeichert."],
    ["Nutzer können Zugang verlangen, Daten berichtigen oder löschen lassen und die Einwilligung widerrufen."],
    ["Durch Nutzung des Dienstes oder Telegram-Autorisierung willigt der Nutzer in diese Richtlinie ein."],
  ],
  fr: [
    ["La présente Politique s'applique à NERIVA sur poliglotai.ru, poliglotai.online et au bot Telegram @NERIVAapp_bot.", "Le service constitue un système unique quel que soit le mode d'accès.", "En utilisant le service, l'Utilisateur accepte cette Politique."],
    ["L'opérateur des données personnelles est le propriétaire du service NERIVA."],
    ["Nous pouvons traiter l'identifiant Telegram, le username s'il est disponible, l'adresse IP, les cookies et les informations sur l'appareil.", "Nous ne traitons pas les données de passeport, les coordonnées de paiement ni les catégories spéciales de données."],
    ["Les données proviennent de l'autorisation Telegram, de l'utilisation du site ou du bot, et automatiquement pendant l'utilisation du service."],
    ["Les données servent à créer et maintenir le compte, fournir l'accès, faire fonctionner le site et le bot, et améliorer le service."],
    ["Les données sont stockées sur des serveurs situés hors de la Fédération de Russie, y compris aux Pays-Bas.", "Des services techniques tiers nécessaires à l'infrastructure peuvent les traiter."],
    ["L'Opérateur applique des mesures techniques et organisationnelles raisonnables contre les accès non autorisés et les fuites."],
    ["Les données sont conservées jusqu'à la suppression du compte ou l'arrêt du service."],
    ["L'Utilisateur peut demander l'accès, la correction ou la suppression des données, ainsi que retirer son consentement."],
    ["En utilisant le service ou l'autorisation Telegram, l'Utilisateur consent à cette Politique."],
  ],
  it: [
    ["La presente Politica si applica a NERIVA su poliglotai.ru, poliglotai.online e al bot Telegram @NERIVAapp_bot.", "Il servizio è un sistema unico indipendentemente dal metodo di accesso.", "Usando il servizio, l'Utente accetta questa Politica."],
    ["Il titolare del trattamento è il proprietario del servizio NERIVA."],
    ["Possiamo trattare Telegram ID, username se disponibile, IP, cookie e informazioni sul dispositivo.", "Non trattiamo dati del passaporto, dati di pagamento o categorie speciali di dati."],
    ["I dati sono forniti tramite autorizzazione Telegram, uso del sito o del bot, e automaticamente durante l'uso del servizio."],
    ["I dati servono a creare e mantenere l'account, fornire l'accesso, far funzionare sito e bot e migliorare il servizio."],
    ["I dati sono conservati su server fuori dalla Federazione Russa, inclusi i Paesi Bassi.", "Possono essere trattati da servizi tecnici terzi necessari all'infrastruttura."],
    ["L'Operatore adotta misure tecniche e organizzative ragionevoli contro accessi non autorizzati e fughe."],
    ["I dati sono conservati fino alla cancellazione dell'account o alla cessazione del servizio."],
    ["L'Utente può chiedere accesso, correzione o cancellazione dei dati e revocare il consenso."],
    ["Usando il servizio o autorizzandosi tramite Telegram, l'Utente accetta questa Politica."],
  ],
  zh: [
    ["本政策适用于 poliglotai.ru、poliglotai.online 以及 Telegram 机器人 @NERIVAapp_bot 上的 NERIVA 服务。", "无论访问方式如何，服务均为统一系统。", "使用服务即表示用户同意本政策。"],
    ["个人数据运营者为 NERIVA 服务所有者。"],
    ["我们可能处理 Telegram ID、可用的 username、IP 地址、cookies 和设备信息。", "我们不处理护照数据、支付凭据或其他特殊类别个人数据。"],
    ["数据来自 Telegram 授权、网站或机器人使用过程，以及使用服务时自动产生的数据。"],
    ["数据用于创建和维护账户、提供服务访问、保障网站和机器人运行并改进服务质量。"],
    ["数据存储在俄罗斯联邦境外的服务器上，包括荷兰。", "必要的第三方技术服务可参与基础设施处理。"],
    ["运营者采取合理的技术和组织措施，防止未经授权的访问和泄露。"],
    ["数据保存至用户删除账户或服务停止运行。"],
    ["用户可以请求访问、更正或删除数据，并撤回处理同意。"],
    ["使用服务或通过 Telegram 授权即表示用户同意本个人数据处理政策。"],
  ],
  ja: [
    ["本ポリシーは poliglotai.ru、poliglotai.online、Telegram ボット @NERIVAapp_bot の NERIVA に適用されます。", "アクセス方法にかかわらず、サービスは一つのシステムです。", "サービスを利用することで、ユーザーは本ポリシーに同意します。"],
    ["個人データの運営者は NERIVA サービスの所有者です。"],
    ["Telegram ID、利用可能な username、IP アドレス、cookies、端末情報を処理する場合があります。", "パスポート情報、決済情報、特別な個人データは処理しません。"],
    ["データは Telegram 認証、サイトまたはボットの利用、サービス利用中の自動取得により提供されます。"],
    ["データはアカウント作成と維持、サービス提供、サイトとボットの運用、品質改善に使用されます。"],
    ["データはロシア連邦外のサーバーに保存され、オランダを含みます。", "インフラ運用に必要な第三者技術サービスで処理される場合があります。"],
    ["運営者は不正アクセスと漏えいを防ぐため合理的な技術的・組織的措置を講じます。"],
    ["データはアカウント削除またはサービス終了まで保存されます。"],
    ["ユーザーはデータへのアクセス、修正、削除、同意撤回を求める権利があります。"],
    ["サービス利用または Telegram 認証により、ユーザーは本ポリシーに同意します。"],
  ],
  ko: [
    ["본 정책은 poliglotai.ru, poliglotai.online 및 Telegram 봇 @NERIVAapp_bot의 NERIVA 서비스에 적용됩니다.", "서비스는 접근 방식과 관계없이 하나의 시스템입니다.", "서비스를 이용하면 사용자는 본 정책에 동의합니다."],
    ["개인정보 처리 운영자는 NERIVA 서비스의 소유자입니다."],
    ["Telegram ID, 사용 가능한 username, IP 주소, cookies, 기기 정보를 처리할 수 있습니다.", "여권 정보, 결제 정보, 특수 범주의 개인정보는 처리하지 않습니다."],
    ["데이터는 Telegram 인증, 사이트 또는 봇 사용, 서비스 사용 중 자동으로 제공됩니다."],
    ["데이터는 계정 생성과 유지, 서비스 접근 제공, 사이트와 봇 운영, 품질 개선에 사용됩니다."],
    ["데이터는 네덜란드를 포함한 러시아 연방 외부 서버에 저장됩니다.", "인프라 운영에 필요한 제3자 기술 서비스가 처리할 수 있습니다."],
    ["운영자는 무단 접근과 유출을 막기 위해 합리적인 기술적·조직적 조치를 적용합니다."],
    ["데이터는 사용자가 계정을 삭제하거나 서비스가 종료될 때까지 저장됩니다."],
    ["사용자는 데이터 접근, 수정, 삭제 및 처리 동의 철회를 요청할 수 있습니다."],
    ["서비스 이용 또는 Telegram 인증을 통해 사용자는 본 정책에 동의합니다."],
  ],
  uz: [
    ["Ushbu Siyosat poliglotai.ru, poliglotai.online va Telegram boti @NERIVAapp_bot orqali ishlaydigan NERIVA xizmatiga taalluqlidir.", "Xizmat kirish usulidan qat'i nazar yagona tizimdir.", "Foydalanish orqali Foydalanuvchi ushbu Siyosatga rozilik bildiradi."],
    ["Shaxsiy ma'lumotlar operatori NERIVA xizmati egasidir."],
    ["Telegram ID, username mavjud bo'lsa, IP manzil, cookies va qurilma ma'lumotlari qayta ishlanishi mumkin.", "Pasport ma'lumotlari, to'lov rekvizitlari va maxsus toifadagi ma'lumotlar qayta ishlanmaydi."],
    ["Ma'lumotlar Telegram avtorizatsiyasi, sayt yoki botdan foydalanish hamda xizmatdan foydalanishda avtomatik tarzda olinadi."],
    ["Ma'lumotlar akkaunt yaratish va qo'llab-quvvatlash, xizmatga kirish, sayt va bot ishlashi hamda sifatni yaxshilash uchun ishlatiladi."],
    ["Ma'lumotlar Rossiya Federatsiyasidan tashqaridagi serverlarda, jumladan Niderlandiyada saqlanadi.", "Infratuzilma uchun zarur uchinchi tomon texnik xizmatlari ishlatilishi mumkin."],
    ["Operator ruxsatsiz kirish va sizib chiqishdan himoya qilish uchun oqilona texnik va tashkiliy choralar ko'radi."],
    ["Ma'lumotlar foydalanuvchi akkauntni o'chirguncha yoki xizmat to'xtaguncha saqlanadi."],
    ["Foydalanuvchi ma'lumotlarga kirish, tuzatish yoki o'chirish hamda rozilikni qaytarib olish huquqiga ega."],
    ["Xizmatdan foydalanish yoki Telegram orqali avtorizatsiya qilish orqali Foydalanuvchi ushbu Siyosatga rozilik bildiradi."],
  ],
  uk: [
    ["Ця Політика діє щодо сервісу NERIVA на poliglotai.ru, poliglotai.online і Telegram-бота @NERIVAapp_bot.", "Сервіс є єдиною системою незалежно від способу доступу.", "Використовуючи сервіс, Користувач погоджується з цією Політикою."],
    ["Оператором персональних даних є власник сервісу NERIVA."],
    ["Ми можемо обробляти Telegram ID, username за наявності, IP-адресу, cookies та інформацію про пристрій.", "Ми не обробляємо паспортні дані, платіжні реквізити й спеціальні категорії даних."],
    ["Дані надходять через авторизацію в Telegram, використання сайту або бота, а також автоматично під час користування сервісом."],
    ["Дані використовуються для створення й підтримки акаунта, доступу до сервісу, роботи сайту й бота та покращення якості."],
    ["Дані зберігаються на серверах за межами Російської Федерації, включно з Нідерландами.", "Можлива обробка сторонніми технічними сервісами, необхідними для інфраструктури."],
    ["Оператор застосовує розумні технічні й організаційні заходи для захисту від несанкціонованого доступу та витоків."],
    ["Дані зберігаються до видалення акаунта користувачем або припинення роботи сервісу."],
    ["Користувач має право запросити доступ, виправити або видалити дані й відкликати згоду."],
    ["Використовуючи сервіс або авторизуючись через Telegram, Користувач погоджується з цією Політикою."],
  ],
  pl: [
    ["Niniejsza Polityka dotyczy NERIVA na poliglotai.ru, poliglotai.online oraz bota Telegram @NERIVAapp_bot.", "Usługa jest jednym systemem niezależnie od sposobu dostępu.", "Korzystając z usługi, Użytkownik akceptuje tę Politykę."],
    ["Operatorem danych osobowych jest właściciel usługi NERIVA."],
    ["Możemy przetwarzać Telegram ID, username jeśli dostępny, adres IP, cookies oraz informacje o urządzeniu.", "Nie przetwarzamy danych paszportowych, danych płatniczych ani szczególnych kategorii danych."],
    ["Dane są przekazywane przez autoryzację Telegram, korzystanie ze strony lub bota oraz automatycznie podczas używania usługi."],
    ["Dane służą do tworzenia i utrzymania konta, zapewnienia dostępu, działania strony i bota oraz poprawy jakości."],
    ["Dane są przechowywane na serwerach poza Federacją Rosyjską, w tym w Niderlandach.", "Mogą je przetwarzać zewnętrzne usługi techniczne potrzebne do działania infrastruktury."],
    ["Operator stosuje rozsądne środki techniczne i organizacyjne chroniące przed nieuprawnionym dostępem i wyciekiem."],
    ["Dane są przechowywane do usunięcia konta przez użytkownika albo zakończenia działania usługi."],
    ["Użytkownik ma prawo żądać dostępu, poprawienia lub usunięcia danych oraz wycofać zgodę."],
    ["Korzystając z usługi lub autoryzując się przez Telegram, Użytkownik wyraża zgodę na tę Politykę."],
  ],
  ro: [
    ["Această Politică se aplică serviciului NERIVA pe poliglotai.ru, poliglotai.online și botului Telegram @NERIVAapp_bot.", "Serviciul este un sistem unic indiferent de metoda de acces.", "Prin utilizare, Utilizatorul acceptă această Politică."],
    ["Operatorul datelor personale este proprietarul serviciului NERIVA."],
    ["Putem prelucra Telegram ID, username dacă este disponibil, IP, cookies și informații despre dispozitiv.", "Nu prelucrăm date de pașaport, date de plată sau categorii speciale de date."],
    ["Datele provin prin autorizare Telegram, folosirea site-ului sau botului și automat în timpul utilizării serviciului."],
    ["Datele sunt folosite pentru creare și suport cont, acces la serviciu, funcționarea site-ului și botului și îmbunătățirea calității."],
    ["Datele sunt stocate pe servere din afara Federației Ruse, inclusiv în Țările de Jos.", "Pot fi prelucrate de servicii tehnice terțe necesare infrastructurii."],
    ["Operatorul aplică măsuri tehnice și organizaționale rezonabile împotriva accesului neautorizat și scurgerilor."],
    ["Datele se păstrează până la ștergerea contului sau încetarea serviciului."],
    ["Utilizatorul poate cere acces, corectare sau ștergere și își poate retrage consimțământul."],
    ["Prin utilizarea serviciului sau autorizarea prin Telegram, Utilizatorul consimte la această Politică."],
  ],
  pt: [
    ["Esta Política aplica-se ao NERIVA em poliglotai.ru, poliglotai.online e ao bot do Telegram @NERIVAapp_bot.", "O serviço é um sistema único independentemente do método de acesso.", "Ao utilizar o serviço, o Utilizador aceita esta Política."],
    ["O operador dos dados pessoais é o proprietário do serviço NERIVA."],
    ["Podemos tratar Telegram ID, username se disponível, IP, cookies e informação do dispositivo.", "Não tratamos dados de passaporte, dados de pagamento nem categorias especiais de dados."],
    ["Os dados são fornecidos por autorização Telegram, uso do site ou bot e automaticamente durante a utilização do serviço."],
    ["Os dados são usados para criar e manter a conta, dar acesso ao serviço, operar site e bot e melhorar a qualidade."],
    ["Os dados são armazenados em servidores fora da Federação Russa, incluindo os Países Baixos.", "Podem ser tratados por serviços técnicos terceiros necessários à infraestrutura."],
    ["O Operador aplica medidas técnicas e organizacionais razoáveis contra acesso não autorizado e fugas."],
    ["Os dados são guardados até o utilizador eliminar a conta ou o serviço deixar de operar."],
    ["O Utilizador pode pedir acesso, correção ou eliminação dos dados e retirar o consentimento."],
    ["Ao usar o serviço ou autorizar-se via Telegram, o Utilizador consente esta Política."],
  ],
};

const cp1251Decoder = new TextDecoder("windows-1251");
const utf8Decoder = new TextDecoder("utf-8", { fatal: true });
const cp1251EncodeMap = new Map();
for (let byte = 0; byte <= 255; byte += 1) {
  const decoded = cp1251Decoder.decode(Uint8Array.of(byte));
  cp1251EncodeMap.set(decoded, byte);
  cp1251EncodeMap.set(String.fromCharCode(byte), byte);
}

function encodeCP1251Like(value) {
  const bytes = [];
  for (const char of value) {
    const byte = cp1251EncodeMap.get(char);
    if (byte === undefined) return null;
    bytes.push(byte);
  }
  return Uint8Array.from(bytes);
}

function mojibakeScore(value) {
  const markers = [
    /\u0420/g,
    /\u0421[\u0400-\u04ff]/g,
    /\u0413./g,
    /\u0432\u0402/g,
    /\u0434[\u0451\u0454\u0456]/g,
    /\u0431\u0453/g,
    /\u043c[\u045c\u045d]/g,
    /\u0436[\u040c\u045c]/g,
    /\u0424./g,
    /\u0425./g,
  ];
  return markers.reduce((score, pattern) => score + (value.match(pattern) || []).length, 0);
}

function repairMojibake(value) {
  if (typeof value !== "string" || mojibakeScore(value) === 0) return value;
  const bytes = encodeCP1251Like(value);
  if (!bytes) return value;
  try {
    const repaired = utf8Decoder.decode(bytes);
    return mojibakeScore(repaired) < mojibakeScore(value) ? repaired : value;
  } catch {
    return value;
  }
}

function repairDeep(value) {
  if (Array.isArray(value)) {
    for (let index = 0; index < value.length; index += 1) {
      value[index] = repairDeep(value[index]);
    }
    return value;
  }
  if (value && typeof value === "object") {
    for (const key of Object.keys(value)) {
      value[key] = repairDeep(value[key]);
    }
    return value;
  }
  return repairMojibake(value);
}

repairDeep(policy);
repairDeep(localized);
repairDeep(labels);

function cloneEnglish(locale) {
  const [title, updated, intro] = localized[locale];
  const next = JSON.parse(JSON.stringify(policy.en));
  next.title = title;
  next.updated = updated;
  next.intro = intro;
  const localizedBlocks = localizedSectionBlocks[locale];
  next.sections = next.sections.map((section, index) => [
    section[0],
    labels[locale][index],
    index === 10
      ? [`Email: ${contacts.email}`, `Telegram: ${contacts.admin}`, `Bot: ${contacts.bot}`]
      : localizedBlocks?.[index] || [intro],
  ]);
  return next;
}

for (const locale of locales) {
  if (!policy[locale]) policy[locale] = cloneEnglish(locale);
}

function escapeHTML(value) {
  return String(value)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

function renderBlock(block) {
  if (Array.isArray(block)) {
    return `<ul>${block.map((item) => `<li>${escapeHTML(item)}</li>`).join("")}</ul>`;
  }
  return `<p>${escapeHTML(block)}</p>`;
}

function renderContactGrid(blocks) {
  const values = blocks.map((item) => String(item).replace(/^[^:]+:\s*/, ""));
  const email = values[0] || contacts.email;
  const admin = values[1] || contacts.admin;
  const bot = values[2] || contacts.bot;
  return `<div class="legal-contact-grid">
    <a href="mailto:${escapeHTML(email)}"><small>Email</small>${escapeHTML(email)}</a>
    <a href="https://t.me/${escapeHTML(admin.replace(/^@/, ""))}"><small>Telegram</small>${escapeHTML(admin)}</a>
    <a href="https://t.me/${escapeHTML(bot.replace(/^@/, ""))}"><small>Telegram bot</small>${escapeHTML(bot)}</a>
  </div>`;
}

function renderHTML(entry) {
  return `<div data-privacy-policy-card>
    <header class="mb-10">
      <p class="eyebrow">${escapeHTML(entry.updated)}</p>
      <h2 class="text-2xl md:text-3xl font-black text-dark mb-4 break-words">${escapeHTML(entry.title)}</h2>
      <p class="break-words">${escapeHTML(entry.intro)}</p>
    </header>
    ${entry.sections.map(([id, title, blocks]) => `<section id="${escapeHTML(id)}" class="mb-10"><h2 class="text-2xl md:text-3xl font-black text-dark mb-4 break-words">${escapeHTML(title)}</h2>${id === "contacts" ? renderContactGrid(blocks) : blocks.map(renderBlock).join("")}</section>`).join("")}
  </div>`;
}

function generatedAsset() {
  return `(function () {
  const policy = ${JSON.stringify(policy, null, 2)};
  const locales = ${JSON.stringify(locales)};

  function escapeHTML(value) {
    return String(value)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;");
  }

  function renderBlock(block) {
    if (Array.isArray(block)) {
      return "<ul>" + block.map((item) => "<li>" + escapeHTML(item) + "</li>").join("") + "</ul>";
    }
    return "<p>" + escapeHTML(block) + "</p>";
  }

  function renderContactGrid(blocks) {
    const values = blocks.map((item) => String(item).replace(/^[^:]+:\\s*/, ""));
    const email = values[0] || "${contacts.email}";
    const admin = values[1] || "${contacts.admin}";
    const bot = values[2] || "${contacts.bot}";
    return '<div class="legal-contact-grid">' +
      '<a href="mailto:' + escapeHTML(email) + '"><small>Email</small>' + escapeHTML(email) + '</a>' +
      '<a href="https://t.me/' + escapeHTML(admin.replace(/^@/, "")) + '"><small>Telegram</small>' + escapeHTML(admin) + '</a>' +
      '<a href="https://t.me/' + escapeHTML(bot.replace(/^@/, "")) + '"><small>Telegram bot</small>' + escapeHTML(bot) + '</a>' +
      '</div>';
  }

  function renderHTML(entry) {
    return '<div data-privacy-policy-card>' +
      '<header class="mb-10"><p class="eyebrow">' + escapeHTML(entry.updated) + '</p><h2 class="text-2xl md:text-3xl font-black text-dark mb-4 break-words">' + escapeHTML(entry.title) + '</h2><p class="break-words">' + escapeHTML(entry.intro) + '</p></header>' +
      entry.sections.map(([id, title, blocks]) => '<section id="' + escapeHTML(id) + '" class="mb-10"><h2 class="text-2xl md:text-3xl font-black text-dark mb-4 break-words">' + escapeHTML(title) + '</h2>' + (id === "contacts" ? renderContactGrid(blocks) : blocks.map(renderBlock).join('')) + '</section>').join('') +
      '</div>';
  }

  function currentLocale() {
    const params = new URLSearchParams(window.location.search);
    const fromQuery = params.get("lang");
    const fromStorage = localStorage.getItem("poliglot-site-language") || localStorage.getItem("poliglot-auth-language");
    const fromHTML = document.documentElement.lang;
    const code = (fromQuery || fromStorage || fromHTML || "ru").toLowerCase().slice(0, 2);
    return locales.includes(code) ? code : "ru";
  }

  function applyPolicy() {
    const locale = currentLocale();
    const entry = policy[locale] || policy.ru;
    document.documentElement.lang = locale;
    document.title = entry.title + " - NERIVA";
    const shell = document.querySelector("[data-privacy-policy-card]") || document.querySelector(".legal-document-shell");
    if (!shell) return false;
    shell.innerHTML = renderHTML(entry);
    document.querySelectorAll('a[href^="/privacy.html"]').forEach((link) => {
      const url = new URL(link.getAttribute("href"), window.location.origin);
      url.searchParams.set("lang", locale);
      link.setAttribute("href", url.pathname + url.search);
    });
    return true;
  }

  function waitAndApply(tries = 80) {
    if (applyPolicy() || tries <= 0) return;
    requestAnimationFrame(() => waitAndApply(tries - 1));
  }

  window.poliglotPrivacyPolicy = { policy, locales, apply: applyPolicy };
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", () => waitAndApply(), { once: true });
  } else {
    waitAndApply();
  }
})();\n`;
}

const asset = generatedAsset();
for (const file of [
  path.join(root, "site-react", "public", "assets", "privacy-policy-i18n.js"),
  path.join(root, "Сайт полиглота для бота", "assets", "privacy-policy-i18n.js"),
]) {
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, asset, "utf8");
}

const legacyPath = path.join(root, "site-react", "src", "legacyLegalContent.ts");
const legacy = fs.readFileSync(legacyPath, "utf8");
const replacement = `export const privacyDocumentHtml = ${JSON.stringify(renderHTML(policy.ru))};\n`;
const updated = legacy.replace(/export const privacyDocumentHtml = "[\s\S]*";\s*$/m, replacement);
if (updated === legacy && !legacy.includes("ПОЛИТИКА ОБРАБОТКИ ПЕРСОНАЛЬНЫХ ДАННЫХ")) {
  throw new Error("Could not replace privacyDocumentHtml");
}
if (updated !== legacy) {
  fs.writeFileSync(legacyPath, updated, "utf8");
}
