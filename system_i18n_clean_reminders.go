package main

var cleanReminderSystemUICopyOverrides = map[string]systemUICopy{
	"en": {
		ReminderTitle: "Notification settings", ReminderStatus: "Status", ReminderEnabled: "enabled", ReminderDisabled: "disabled",
		ReminderTimezone: "Time zone", ReminderDefault: "default", ReminderHour: "Send time", ReminderHint: "You can enable or disable daily reminders and choose a time zone.",
		ReminderOnButton: "Enable reminders", ReminderOffButton: "Disable reminders", ReminderFooter: "Open %s and choose the next step. You can turn reminders off with /reminderoff.",
		ReminderMessages: []string{"A small language step today is better than a perfect plan for tomorrow.", "Come back for a few minutes and keep your language rhythm alive.", "One word, one phrase, one small win. Your progress is waiting."},
	},
	"ru": {
		ReminderTitle: "Настройки уведомлений", ReminderStatus: "Статус", ReminderEnabled: "включены", ReminderDisabled: "выключены",
		ReminderTimezone: "Часовой пояс", ReminderDefault: "по умолчанию", ReminderHour: "Время отправки", ReminderHint: "Можно включить или отключить ежедневные напоминания и выбрать часовой пояс.",
		ReminderOnButton: "Включить уведомления", ReminderOffButton: "Отключить уведомления", ReminderFooter: "Открой %s и выбери следующий шаг. Напоминания можно выключить командой /reminderoff.",
		ReminderMessages: []string{"Маленький языковой шаг сегодня лучше идеального плана на завтра.", "Вернись на несколько минут и сохрани ритм обучения.", "Одно слово, одна фраза, одна маленькая победа. Прогресс ждёт тебя."},
	},
	"es": {
		ReminderTitle: "Ajustes de notificaciones", ReminderStatus: "Estado", ReminderEnabled: "activadas", ReminderDisabled: "desactivadas",
		ReminderTimezone: "Zona horaria", ReminderDefault: "por defecto", ReminderHour: "Hora de envío", ReminderHint: "Puedes activar o desactivar los recordatorios diarios y elegir una zona horaria.",
		ReminderOnButton: "Activar recordatorios", ReminderOffButton: "Desactivar recordatorios", ReminderFooter: "Abre %s y elige el siguiente paso. Puedes desactivar los recordatorios con /reminderoff.",
		ReminderMessages: []string{"Un pequeño paso de idioma hoy vale más que un plan perfecto para mañana.", "Vuelve unos minutos y mantén vivo tu ritmo de aprendizaje.", "Una palabra, una frase, una pequeña victoria. Tu progreso te espera."},
	},
	"de": {
		ReminderTitle: "Benachrichtigungseinstellungen", ReminderStatus: "Status", ReminderEnabled: "aktiviert", ReminderDisabled: "deaktiviert",
		ReminderTimezone: "Zeitzone", ReminderDefault: "Standard", ReminderHour: "Sendezeit", ReminderHint: "Du kannst tägliche Erinnerungen aktivieren oder deaktivieren und die Zeitzone wählen.",
		ReminderOnButton: "Erinnerungen aktivieren", ReminderOffButton: "Erinnerungen deaktivieren", ReminderFooter: "Öffne %s und wähle den nächsten Schritt. Erinnerungen kannst du mit /reminderoff ausschalten.",
		ReminderMessages: []string{"Ein kleiner Sprachschritt heute ist besser als ein perfekter Plan für morgen.", "Komm für ein paar Minuten zurück und halte deinen Lernrhythmus.", "Ein Wort, ein Satz, ein kleiner Sieg. Dein Fortschritt wartet."},
	},
	"fr": {
		ReminderTitle: "Réglages des notifications", ReminderStatus: "Statut", ReminderEnabled: "activées", ReminderDisabled: "désactivées",
		ReminderTimezone: "Fuseau horaire", ReminderDefault: "par défaut", ReminderHour: "Heure d'envoi", ReminderHint: "Tu peux activer ou désactiver les rappels quotidiens et choisir un fuseau horaire.",
		ReminderOnButton: "Activer les rappels", ReminderOffButton: "Désactiver les rappels", ReminderFooter: "Ouvre %s et choisis la suite. Tu peux désactiver les rappels avec /reminderoff.",
		ReminderMessages: []string{"Un petit pas linguistique aujourd'hui vaut mieux qu'un plan parfait pour demain.", "Reviens quelques minutes et garde ton rythme d'apprentissage.", "Un mot, une phrase, une petite victoire. Ton progrès t'attend."},
	},
	"it": {
		ReminderTitle: "Impostazioni notifiche", ReminderStatus: "Stato", ReminderEnabled: "attive", ReminderDisabled: "disattivate",
		ReminderTimezone: "Fuso orario", ReminderDefault: "predefinito", ReminderHour: "Ora di invio", ReminderHint: "Puoi attivare o disattivare i promemoria giornalieri e scegliere il fuso orario.",
		ReminderOnButton: "Attiva promemoria", ReminderOffButton: "Disattiva promemoria", ReminderFooter: "Apri %s e scegli il prossimo passo. Puoi disattivare i promemoria con /reminderoff.",
		ReminderMessages: []string{"Un piccolo passo linguistico oggi vale più di un piano perfetto per domani.", "Torna per qualche minuto e mantieni vivo il ritmo.", "Una parola, una frase, una piccola vittoria. Il progresso ti aspetta."},
	},
	"zh": {
		ReminderTitle: "通知设置", ReminderStatus: "状态", ReminderEnabled: "已开启", ReminderDisabled: "已关闭",
		ReminderTimezone: "时区", ReminderDefault: "默认", ReminderHour: "发送时间", ReminderHint: "你可以开启或关闭每日提醒，并选择自己的时区。",
		ReminderOnButton: "开启提醒", ReminderOffButton: "关闭提醒", ReminderFooter: "打开 %s，选择下一步。你可以用 /reminderoff 关闭提醒。",
		ReminderMessages: []string{"今天迈出一个小小的语言步伐，比明天的完美计划更有用。", "回来练习几分钟，让语言节奏继续保持。", "一个单词，一个短语，一个小胜利。你的进度在等你。"},
	},
	"ja": {
		ReminderTitle: "通知設定", ReminderStatus: "状態", ReminderEnabled: "オン", ReminderDisabled: "オフ",
		ReminderTimezone: "タイムゾーン", ReminderDefault: "既定", ReminderHour: "送信時刻", ReminderHint: "毎日の通知をオン/オフにし、タイムゾーンを選べます。",
		ReminderOnButton: "通知をオンにする", ReminderOffButton: "通知をオフにする", ReminderFooter: "%s を開いて次のステップを選んでください。/reminderoff で通知をオフにできます。",
		ReminderMessages: []string{"今日の小さな一歩は、明日の完璧な計画より役に立ちます。", "数分だけ戻って、学習のリズムを保ちましょう。", "一つの単語、一つのフレーズ、一つの小さな勝利。進歩が待っています。"},
	},
	"ko": {
		ReminderTitle: "알림 설정", ReminderStatus: "상태", ReminderEnabled: "켜짐", ReminderDisabled: "꺼짐",
		ReminderTimezone: "시간대", ReminderDefault: "기본값", ReminderHour: "보낼 시간", ReminderHint: "매일 알림을 켜거나 끄고 시간대를 선택할 수 있습니다.",
		ReminderOnButton: "알림 켜기", ReminderOffButton: "알림 끄기", ReminderFooter: "%s 를 열고 다음 단계를 선택하세요. /reminderoff 로 알림을 끌 수 있습니다.",
		ReminderMessages: []string{"오늘의 작은 언어 한 걸음이 내일의 완벽한 계획보다 낫습니다.", "몇 분만 돌아와서 학습 리듬을 이어 가세요.", "단어 하나, 구문 하나, 작은 승리 하나. 진도가 기다리고 있습니다."},
	},
	"tg": {
		ReminderTitle: "Танзимоти огоҳинома", ReminderStatus: "Ҳолат", ReminderEnabled: "фаъол", ReminderDisabled: "хомӯш",
		ReminderTimezone: "Минтақаи вақт", ReminderDefault: "пешфарз", ReminderHour: "Вақти фиристодан", ReminderHint: "Метавонӣ ёдраскуниҳои ҳаррӯзаро фаъол ё хомӯш кунӣ ва минтақаи вақтро интихоб намоӣ.",
		ReminderOnButton: "Фаъол кардани ёдраскуниҳо", ReminderOffButton: "Хомӯш кардани ёдраскуниҳо", ReminderFooter: "%s-ро кушо ва қадами навбатиро интихоб кун. Ёдраскуниҳоро бо /reminderoff хомӯш кардан мумкин аст.",
		ReminderMessages: []string{"Як қадами хурди забонӣ имрӯз аз нақшаи комили фардо беҳтар аст.", "Барои чанд дақиқа баргард ва ритми омӯзишро нигоҳ дор.", "Як калима, як ибора, як пирӯзии хурд. Пешрафт интизор аст."},
	},
	"uz": {
		ReminderTitle: "Bildirishnoma sozlamalari", ReminderStatus: "Holat", ReminderEnabled: "yoqilgan", ReminderDisabled: "o'chirilgan",
		ReminderTimezone: "Vaqt mintaqasi", ReminderDefault: "standart", ReminderHour: "Yuborish vaqti", ReminderHint: "Kundalik eslatmalarni yoqish yoki o'chirish va vaqt mintaqasini tanlash mumkin.",
		ReminderOnButton: "Eslatmalarni yoqish", ReminderOffButton: "Eslatmalarni o'chirish", ReminderFooter: "%s ni oching va keyingi qadamni tanlang. Eslatmalarni /reminderoff bilan o'chirish mumkin.",
		ReminderMessages: []string{"Bugungi kichik til qadami ertangi mukammal rejadan foydaliroq.", "Bir necha daqiqaga qayting va o'qish ritmini saqlang.", "Bitta so'z, bitta ibora, bitta kichik g'alaba. Progress sizni kutmoqda."},
	},
	"tt": {
		ReminderTitle: "Искәртмә көйләүләре", ReminderStatus: "Хәл", ReminderEnabled: "кабызылган", ReminderDisabled: "сүндерелгән",
		ReminderTimezone: "Вакыт поясы", ReminderDefault: "килешү буенча", ReminderHour: "Җибәрү вакыты", ReminderHint: "Көндәлек искәртмәләрне кабызып яки сүндереп, вакыт поясын сайлый аласың.",
		ReminderOnButton: "Искәртмәләрне кабызу", ReminderOffButton: "Искәртмәләрне сүндерү", ReminderFooter: "%s ач һәм киләсе адымны сайла. Искәртмәләрне /reminderoff белән сүндереп була.",
		ReminderMessages: []string{"Бүгенге кечкенә тел адымы иртәгәге камил планнан яхшырак.", "Берничә минутка әйләнеп кайт һәм өйрәнү ритмын сакла.", "Бер сүз, бер фраза, бер кечкенә җиңү. Алга китеш сине көтә."},
	},
	"hy": {
		ReminderTitle: "Ծանուցումների կարգավորումներ", ReminderStatus: "Կարգավիճակ", ReminderEnabled: "միացված", ReminderDisabled: "անջատված",
		ReminderTimezone: "Ժամային գոտի", ReminderDefault: "լռելյայն", ReminderHour: "Ուղարկման ժամ", ReminderHint: "Կարող ես միացնել կամ անջատել ամենօրյա հիշեցումները և ընտրել ժամային գոտին։",
		ReminderOnButton: "Միացնել հիշեցումները", ReminderOffButton: "Անջատել հիշեցումները", ReminderFooter: "Բացիր %s և ընտրիր հաջորդ քայլը։ Հիշեցումները կարելի է անջատել /reminderoff հրամանով։",
		ReminderMessages: []string{"Այսօրվա փոքր լեզվական քայլը ավելի լավ է, քան վաղվա կատարյալ ծրագիրը։", "Վերադարձիր մի քանի րոպեով և պահպանիր ուսուցման ռիթմը։", "Մեկ բառ, մեկ արտահայտություն, մեկ փոքր հաղթանակ։ Քո առաջընթացը սպասում է։"},
	},
	"kk": {
		ReminderTitle: "Хабарлама баптаулары", ReminderStatus: "Күй", ReminderEnabled: "қосулы", ReminderDisabled: "өшірулі",
		ReminderTimezone: "Уақыт белдеуі", ReminderDefault: "әдепкі", ReminderHour: "Жіберу уақыты", ReminderHint: "Күнделікті еске салуды қосуға немесе өшіруге және уақыт белдеуін таңдауға болады.",
		ReminderOnButton: "Еске салуды қосу", ReminderOffButton: "Еске салуды өшіру", ReminderFooter: "%s ашып, келесі қадамды таңда. Еске салуды /reminderoff арқылы өшіруге болады.",
		ReminderMessages: []string{"Бүгінгі шағын тіл қадамы ертеңгі мінсіз жоспардан пайдалырақ.", "Бірнеше минутқа оралып, оқу ырғағын сақта.", "Бір сөз, бір фраза, бір кішкентай жеңіс. Прогресс сені күтеді."},
	},
	"ky": {
		ReminderTitle: "Эскертме жөндөөлөрү", ReminderStatus: "Абалы", ReminderEnabled: "күйүк", ReminderDisabled: "өчүк",
		ReminderTimezone: "Убакыт алкагы", ReminderDefault: "демейки", ReminderHour: "Жөнөтүү убактысы", ReminderHint: "Күн сайынкы эскертмелерди күйгүзүп же өчүрүп, убакыт алкагын тандай аласың.",
		ReminderOnButton: "Эскертмелерди күйгүзүү", ReminderOffButton: "Эскертмелерди өчүрүү", ReminderFooter: "%s ачып, кийинки кадамды танда. Эскертмелерди /reminderoff менен өчүрсө болот.",
		ReminderMessages: []string{"Бүгүнкү кичинекей тил кадамы эртеңки идеалдуу пландан жакшы.", "Бир нече мүнөткө кайтып, окуу ыргагын сакта.", "Бир сөз, бир фраза, бир кичинекей жеңиш. Ийгилик сени күтөт."},
	},
	"ka": {
		ReminderTitle: "შეტყობინებების პარამეტრები", ReminderStatus: "სტატუსი", ReminderEnabled: "ჩართული", ReminderDisabled: "გამორთული",
		ReminderTimezone: "დროის სარტყელი", ReminderDefault: "ნაგულისხმევი", ReminderHour: "გაგზავნის დრო", ReminderHint: "შეგიძლია ყოველდღიური შეხსენებები ჩართო ან გამორთო და დროის სარტყელი აირჩიო.",
		ReminderOnButton: "შეხსენებების ჩართვა", ReminderOffButton: "შეხსენებების გამორთვა", ReminderFooter: "გახსენი %s და აირჩიე შემდეგი ნაბიჯი. შეხსენებების გამორთვა შეიძლება /reminderoff ბრძანებით.",
		ReminderMessages: []string{"დღევანდელი პატარა ენობრივი ნაბიჯი ხვალინდელ სრულყოფილ გეგმაზე უკეთესია.", "დაბრუნდი რამდენიმე წუთით და სწავლის რიტმი შეინარჩუნე.", "ერთი სიტყვა, ერთი ფრაზა, ერთი პატარა გამარჯვება. პროგრესი გელოდება."},
	},
	"uk": {
		ReminderTitle: "Налаштування сповіщень", ReminderStatus: "Статус", ReminderEnabled: "увімкнені", ReminderDisabled: "вимкнені",
		ReminderTimezone: "Часовий пояс", ReminderDefault: "за замовчуванням", ReminderHour: "Час надсилання", ReminderHint: "Можна увімкнути або вимкнути щоденні нагадування й вибрати часовий пояс.",
		ReminderOnButton: "Увімкнути нагадування", ReminderOffButton: "Вимкнути нагадування", ReminderFooter: "Відкрий %s і вибери наступний крок. Нагадування можна вимкнути командою /reminderoff.",
		ReminderMessages: []string{"Маленький мовний крок сьогодні кращий за ідеальний план на завтра.", "Повернися на кілька хвилин і збережи ритм навчання.", "Одне слово, одна фраза, одна маленька перемога. Прогрес чекає на тебе."},
	},
	"pl": {
		ReminderTitle: "Ustawienia powiadomień", ReminderStatus: "Status", ReminderEnabled: "włączone", ReminderDisabled: "wyłączone",
		ReminderTimezone: "Strefa czasowa", ReminderDefault: "domyślna", ReminderHour: "Godzina wysyłki", ReminderHint: "Możesz włączyć lub wyłączyć codzienne przypomnienia i wybrać strefę czasową.",
		ReminderOnButton: "Włącz przypomnienia", ReminderOffButton: "Wyłącz przypomnienia", ReminderFooter: "Otwórz %s i wybierz następny krok. Przypomnienia możesz wyłączyć komendą /reminderoff.",
		ReminderMessages: []string{"Mały krok językowy dziś jest lepszy niż idealny plan na jutro.", "Wróć na kilka minut i utrzymaj rytm nauki.", "Jedno słowo, jedna fraza, jedno małe zwycięstwo. Postęp czeka."},
	},
	"ro": {
		ReminderTitle: "Setări notificări", ReminderStatus: "Stare", ReminderEnabled: "activate", ReminderDisabled: "dezactivate",
		ReminderTimezone: "Fus orar", ReminderDefault: "implicit", ReminderHour: "Ora trimiterii", ReminderHint: "Poți activa sau dezactiva mementourile zilnice și poți alege fusul orar.",
		ReminderOnButton: "Activează mementourile", ReminderOffButton: "Dezactivează mementourile", ReminderFooter: "Deschide %s și alege următorul pas. Poți dezactiva mementourile cu /reminderoff.",
		ReminderMessages: []string{"Un pas mic de limbă astăzi valorează mai mult decât un plan perfect pentru mâine.", "Revino câteva minute și păstrează ritmul de învățare.", "Un cuvânt, o frază, o mică victorie. Progresul te așteaptă."},
	},
	"pt": {
		ReminderTitle: "Configurações de notificações", ReminderStatus: "Estado", ReminderEnabled: "ativadas", ReminderDisabled: "desativadas",
		ReminderTimezone: "Fuso horário", ReminderDefault: "padrão", ReminderHour: "Hora de envio", ReminderHint: "Você pode ativar ou desativar lembretes diários e escolher o fuso horário.",
		ReminderOnButton: "Ativar lembretes", ReminderOffButton: "Desativar lembretes", ReminderFooter: "Abra %s e escolha o próximo passo. Você pode desativar os lembretes com /reminderoff.",
		ReminderMessages: []string{"Um pequeno passo de idioma hoje vale mais que um plano perfeito para amanhã.", "Volte por alguns minutos e mantenha o ritmo de estudo.", "Uma palavra, uma frase, uma pequena vitória. Seu progresso espera por você."},
	},
}

var cleanReminderMenuButtons = map[string]string{
	"en": "Menu", "ru": "Меню", "es": "Menú", "de": "Menü", "fr": "Menu", "it": "Menu",
	"zh": "菜单", "ja": "メニュー", "ko": "메뉴", "tg": "Меню", "uz": "Menyu", "tt": "Меню",
	"hy": "Մենյու", "kk": "Мәзір", "ky": "Меню", "ka": "მენიუ", "uk": "Меню", "pl": "Menu",
	"ro": "Meniu", "pt": "Menu",
}

func cleanReminderMenuButton(code string) string {
	code = normalizeInterfaceLanguage(code)
	if label, ok := cleanReminderMenuButtons[code]; ok {
		return label
	}
	return cleanReminderMenuButtons["ru"]
}
