package main

import "strings"

type premiumLandingPlanCopy struct {
	Label       string
	Description string
	Included    []string
	Locked      []string
	Note        string
}

type premiumLandingCopySet struct {
	Free     premiumLandingPlanCopy
	Premium  premiumLandingPlanCopy
	Platinum premiumLandingPlanCopy
}

func premiumLandingPlan(label, description string, included, locked []string, note string) premiumLandingPlanCopy {
	return premiumLandingPlanCopy{
		Label:       label,
		Description: description,
		Included:    included,
		Locked:      locked,
		Note:        note,
	}
}

var premiumLandingCopyByLanguage = map[string]premiumLandingCopySet{
	"ru": {
		Free: premiumLandingPlan(`Попробовать маршрут`, `Базовое текстовое обучение, тренировка слов, Phrasebook и обзор прогресса. AI Tutor, аудирование, произношение и голосовые проверки открываются в Premium.`,
			[]string{`ежедневная привычка и стартовые уроки`, `базовая тренировка слов`, `заметки, phrasebook и обзор прогресса`},
			[]string{`AI Tutor guided lessons`, `Listening и pronunciation`, `voice checks и photo tools`},
			`Подходит для знакомства с продуктом без оплаты.`),
		Premium: premiumLandingPlan(`Регулярная учеба`, `Основной режим для ежедневной практики: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools и расширенные дневные лимиты.`,
			[]string{`голос в текст и перевод услышанного`, `перевод текста с картинки`, `практика по контексту голоса или фото`, `словарь ошибок, notes, XP и streak`},
			nil,
			`Лучший выбор для стабильного ежедневного обучения.`),
		Platinum: premiumLandingPlan(`Интенсив`, `AI Tutor с максимальными дневными лимитами, глубиной roleplay, интенсивным review и максимальной voice/pronunciation практикой.`,
			[]string{`максимальные дневные лимиты`, `больше voice/photo-context практики`, `интенсивный review слабых мест`, `лучший режим для heavy daily learning`},
			nil,
			`Для поездки, работы, экзамена или очень плотного темпа.`),
	},
	"en": {
		Free: premiumLandingPlan(`Try the path`, `Basic text learning, word training, Phrasebook, and progress overview. AI Tutor, listening, pronunciation, and voice checks open in Premium.`,
			[]string{`daily habit and starter lessons`, `basic word training`, `notes, phrasebook, and progress overview`},
			[]string{`AI Tutor guided lessons`, `Listening and pronunciation`, `voice checks and photo tools`},
			`Best for getting to know the product without paying.`),
		Premium: premiumLandingPlan(`Regular study`, `The main mode for daily practice: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools, and expanded daily limits.`,
			[]string{`voice to text and translation of what you heard`, `translate text from an image`, `practice from voice or photo context`, `mistake dictionary, notes, XP, and streak`},
			nil,
			`The best choice for steady daily learning.`),
		Platinum: premiumLandingPlan(`Intensive`, `AI Tutor with maximum daily limits, deeper roleplay, intensive review, and maximum voice/pronunciation practice.`,
			[]string{`maximum daily limits`, `more voice/photo-context practice`, `intensive review of weak spots`, `best mode for heavy daily learning`},
			nil,
			`For travel, work, an exam, or a very dense pace.`),
	},
	"es": {
		Free: premiumLandingPlan(`Probar la ruta`, `Aprendizaje básico con texto, entrenamiento de palabras, Phrasebook y resumen de progreso. AI Tutor, listening, pronunciation y revisiones de voz se abren en Premium.`,
			[]string{`hábito diario y lecciones iniciales`, `entrenamiento básico de palabras`, `notas, phrasebook y resumen de progreso`},
			[]string{`lecciones guiadas de AI Tutor`, `Listening y pronunciation`, `revisiones de voz y herramientas de foto`},
			`Adecuado para conocer el producto sin pagar.`),
		Premium: premiumLandingPlan(`Estudio regular`, `Modo principal para la práctica diaria: AI Tutor, listening, pronunciation, lecciones guiadas de IA, revisiones de voz, photo tools y límites diarios ampliados.`,
			[]string{`voz a texto y traducción de lo escuchado`, `traducción de texto desde una imagen`, `práctica con contexto de voz o foto`, `diccionario de errores, notes, XP y streak`},
			nil,
			`La mejor opción para un aprendizaje diario estable.`),
		Platinum: premiumLandingPlan(`Intensivo`, `AI Tutor con límites diarios máximos, roleplay más profundo, review intensivo y máxima práctica de voice/pronunciation.`,
			[]string{`límites diarios máximos`, `más práctica con contexto de voz/foto`, `review intensivo de puntos débiles`, `mejor modo para heavy daily learning`},
			nil,
			`Para un viaje, trabajo, examen o un ritmo muy denso.`),
	},
	"de": {
		Free: premiumLandingPlan(`Route testen`, `Grundlegendes Textlernen, Worttraining, Phrasebook und Fortschrittsüberblick. AI Tutor, listening, pronunciation und Sprachchecks werden in Premium freigeschaltet.`,
			[]string{`tägliche Routine und Startlektionen`, `grundlegendes Worttraining`, `Notizen, phrasebook und Fortschrittsüberblick`},
			[]string{`geführte AI Tutor-Lektionen`, `Listening und pronunciation`, `Sprachchecks und Foto-Tools`},
			`Geeignet, um das Produkt ohne Zahlung kennenzulernen.`),
		Premium: premiumLandingPlan(`Regelmäßiges Lernen`, `Der Hauptmodus für tägliche Praxis: AI Tutor, listening, pronunciation, geführte KI-Lektionen, voice checks, photo tools und erweiterte Tageslimits.`,
			[]string{`Sprache zu Text und Übersetzung des Gehörten`, `Text aus einem Bild übersetzen`, `Praxis aus Sprach- oder Fotokontext`, `Fehlerwörterbuch, notes, XP und streak`},
			nil,
			`Die beste Wahl für stabiles tägliches Lernen.`),
		Platinum: premiumLandingPlan(`Intensiv`, `AI Tutor mit maximalen Tageslimits, tieferem roleplay, intensivem review und maximaler voice/pronunciation-Praxis.`,
			[]string{`maximale Tageslimits`, `mehr voice/photo-context Praxis`, `intensiver review schwacher Stellen`, `bester Modus für heavy daily learning`},
			nil,
			`Für Reise, Arbeit, Prüfung oder ein sehr dichtes Tempo.`),
	},
	"fr": {
		Free: premiumLandingPlan(`Essayer le parcours`, `Apprentissage texte de base, entraînement des mots, Phrasebook et vue d'ensemble des progrès. AI Tutor, listening, pronunciation et vérifications vocales s'ouvrent avec Premium.`,
			[]string{`habitude quotidienne et leçons de départ`, `entraînement de mots de base`, `notes, phrasebook et vue d'ensemble des progrès`},
			[]string{`leçons guidées AI Tutor`, `Listening et pronunciation`, `vérifications vocales et outils photo`},
			`Convient pour découvrir le produit sans payer.`),
		Premium: premiumLandingPlan(`Étude régulière`, `Mode principal pour la pratique quotidienne : AI Tutor, listening, pronunciation, leçons guidées par l'IA, voice checks, photo tools et limites quotidiennes élargies.`,
			[]string{`voix en texte et traduction de ce qui est entendu`, `traduction de texte depuis une image`, `pratique selon le contexte vocal ou photo`, `dictionnaire d'erreurs, notes, XP et streak`},
			nil,
			`Le meilleur choix pour un apprentissage quotidien stable.`),
		Platinum: premiumLandingPlan(`Intensif`, `AI Tutor avec limites quotidiennes maximales, roleplay plus profond, review intensif et pratique voice/pronunciation maximale.`,
			[]string{`limites quotidiennes maximales`, `plus de pratique voice/photo-context`, `review intensif des points faibles`, `meilleur mode pour heavy daily learning`},
			nil,
			`Pour un voyage, le travail, un examen ou un rythme très dense.`),
	},
	"it": {
		Free: premiumLandingPlan(`Prova il percorso`, `Apprendimento testuale di base, allenamento delle parole, Phrasebook e panoramica dei progressi. AI Tutor, listening, pronunciation e controlli vocali si aprono con Premium.`,
			[]string{`abitudine quotidiana e lezioni iniziali`, `allenamento base delle parole`, `note, phrasebook e panoramica dei progressi`},
			[]string{`lezioni guidate AI Tutor`, `Listening e pronunciation`, `controlli vocali e strumenti foto`},
			`Adatto per conoscere il prodotto senza pagamento.`),
		Premium: premiumLandingPlan(`Studio regolare`, `La modalità principale per la pratica quotidiana: AI Tutor, listening, pronunciation, lezioni guidate dall'IA, voice checks, photo tools e limiti giornalieri estesi.`,
			[]string{`voce in testo e traduzione dell'ascoltato`, `traduzione del testo da un'immagine`, `pratica dal contesto voce o foto`, `dizionario degli errori, notes, XP e streak`},
			nil,
			`La scelta migliore per un apprendimento quotidiano stabile.`),
		Platinum: premiumLandingPlan(`Intensivo`, `AI Tutor con limiti giornalieri massimi, roleplay più profondo, review intensivo e massima pratica voice/pronunciation.`,
			[]string{`limiti giornalieri massimi`, `più pratica voice/photo-context`, `review intensivo dei punti deboli`, `modalità migliore per heavy daily learning`},
			nil,
			`Per un viaggio, lavoro, esame o ritmo molto intenso.`),
	},
	"zh": {
		Free: premiumLandingPlan(`试用路线`, `基础文本学习、单词训练、Phrasebook 和进度概览。AI Tutor、listening、pronunciation 和语音检查在 Premium 中开放。`,
			[]string{`每日习惯和入门课程`, `基础单词训练`, `notes、phrasebook 和进度概览`},
			[]string{`AI Tutor 引导课程`, `Listening 和 pronunciation`, `语音检查和照片工具`},
			`适合在不付费的情况下了解产品。`),
		Premium: premiumLandingPlan(`日常学习`, `每日练习的主要模式：AI Tutor、listening、pronunciation、AI 引导课程、voice checks、photo tools 和更高的每日限额。`,
			[]string{`语音转文字并翻译听到的内容`, `翻译图片中的文字`, `根据语音或照片上下文练习`, `错误词典、notes、XP 和 streak`},
			nil,
			`稳定日常学习的最佳选择。`),
		Platinum: premiumLandingPlan(`强化`, `AI Tutor 配备最高每日限额、更深入的 roleplay、密集 review 和最大 voice/pronunciation 练习。`,
			[]string{`最高每日限额`, `更多 voice/photo-context 练习`, `密集 review 薄弱点`, `heavy daily learning 的最佳模式`},
			nil,
			`适合旅行、工作、考试或非常紧凑的节奏。`),
	},
	"ja": {
		Free: premiumLandingPlan(`ルートを試す`, `基本的なテキスト学習、単語トレーニング、Phrasebook、進捗の概要。AI Tutor、listening、pronunciation、音声チェックは Premium で使えます。`,
			[]string{`毎日の習慣と入門レッスン`, `基本の単語トレーニング`, `notes、phrasebook、進捗の概要`},
			[]string{`AI Tutor のガイド付きレッスン`, `Listening と pronunciation`, `音声チェックと写真ツール`},
			`支払いなしで製品を試すのに向いています。`),
		Premium: premiumLandingPlan(`定期学習`, `毎日の練習の基本モード：AI Tutor、listening、pronunciation、AI ガイドレッスン、voice checks、photo tools、拡張された日次上限。`,
			[]string{`音声をテキスト化し聞いた内容を翻訳`, `画像内テキストの翻訳`, `音声または写真の文脈で練習`, `間違い辞書、notes、XP、streak`},
			nil,
			`安定した毎日の学習に最適です。`),
		Platinum: premiumLandingPlan(`集中`, `AI Tutor に最大の日次上限、深い roleplay、集中的な review、最大の voice/pronunciation 練習を追加。`,
			[]string{`最大の日次上限`, `より多い voice/photo-context 練習`, `弱点の集中的な review`, `heavy daily learning に最適なモード`},
			nil,
			`旅行、仕事、試験、または非常に密なペース向けです。`),
	},
	"ko": {
		Free: premiumLandingPlan(`학습 경로 체험`, `기본 텍스트 학습, 단어 훈련, Phrasebook, 진행 상황 개요. AI Tutor, listening, pronunciation, 음성 평가는 Premium에서 열립니다.`,
			[]string{`매일 습관과 시작 레슨`, `기본 단어 훈련`, `notes, phrasebook, 진행 상황 개요`},
			[]string{`AI Tutor 가이드 레슨`, `Listening 및 pronunciation`, `음성 평가와 사진 도구`},
			`결제 없이 제품을 알아보기에 적합합니다.`),
		Premium: premiumLandingPlan(`정기 학습`, `매일 연습을 위한 기본 모드: AI Tutor, listening, pronunciation, AI 가이드 레슨, voice checks, photo tools, 확장된 일일 한도.`,
			[]string{`음성을 텍스트로 바꾸고 들은 내용 번역`, `이미지의 텍스트 번역`, `음성 또는 사진 맥락으로 연습`, `실수 사전, notes, XP, streak`},
			nil,
			`꾸준한 매일 학습에 가장 좋은 선택입니다.`),
		Platinum: premiumLandingPlan(`집중`, `AI Tutor에 최대 일일 한도, 더 깊은 roleplay, 집중 review, 최대 voice/pronunciation 연습을 제공합니다.`,
			[]string{`최대 일일 한도`, `더 많은 voice/photo-context 연습`, `약점 집중 review`, `heavy daily learning에 가장 좋은 모드`},
			nil,
			`여행, 업무, 시험 또는 매우 촘촘한 속도에 적합합니다.`),
	},
	"tg": {
		Free: premiumLandingPlan(`Роҳро санҷед`, `Омӯзиши асосии матнӣ, машқи калимаҳо, Phrasebook ва шарҳи пешрафт. AI Tutor, listening, pronunciation ва санҷишҳои овозӣ дар Premium кушода мешаванд.`,
			[]string{`одати ҳаррӯза ва дарсҳои оғозӣ`, `машқи асосии калимаҳо`, `notes, phrasebook ва шарҳи пешрафт`},
			[]string{`дарсҳои роҳнамоии AI Tutor`, `Listening ва pronunciation`, `санҷишҳои овозӣ ва абзорҳои акс`},
			`Барои шиносоӣ бо маҳсулот бе пардохт мувофиқ аст.`),
		Premium: premiumLandingPlan(`Омӯзиши мунтазам`, `Режими асосӣ барои машқи ҳаррӯза: AI Tutor, listening, pronunciation, дарсҳои роҳнамоии AI, voice checks, photo tools ва лимитҳои рӯзонаи васеъ.`,
			[]string{`овоз ба матн ва тарҷумаи шунидашуда`, `тарҷумаи матн аз акс`, `машқ аз контексти овоз ё акс`, `луғати хатогиҳо, notes, XP ва streak`},
			nil,
			`Беҳтарин интихоб барои омӯзиши устувори ҳаррӯза.`),
		Platinum: premiumLandingPlan(`Интенсив`, `AI Tutor бо лимитҳои рӯзонаи ҳадди аксар, roleplay амиқ, review шадид ва машқи максималии voice/pronunciation.`,
			[]string{`лимитҳои рӯзонаи ҳадди аксар`, `бештар машқи voice/photo-context`, `review шадиди нуқтаҳои суст`, `беҳтарин режим барои heavy daily learning`},
			nil,
			`Барои сафар, кор, имтиҳон ё суръати бисёр зич.`),
	},
	"uz": {
		Free: premiumLandingPlan(`Marshrutni sinab ko'rish`, `Asosiy matnli o'rganish, so'z mashqi, Phrasebook va progress sharhi. AI Tutor, listening, pronunciation va ovoz tekshiruvlari Premiumda ochiladi.`,
			[]string{`kundalik odat va boshlang'ich darslar`, `asosiy so'z mashqi`, `notes, phrasebook va progress sharhi`},
			[]string{`AI Tutor guided lessons`, `Listening va pronunciation`, `voice checks va photo tools`},
			`Mahsulotni to'lovsiz tanib olish uchun mos.`),
		Premium: premiumLandingPlan(`Muntazam o'qish`, `Kundalik amaliyot uchun asosiy rejim: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools va kengaytirilgan kunlik limitlar.`,
			[]string{`ovozni matnga aylantirish va eshitilganni tarjima qilish`, `rasmdagi matnni tarjima qilish`, `ovoz yoki foto konteksti bo'yicha amaliyot`, `xatolar lug'ati, notes, XP va streak`},
			nil,
			`Barqaror kundalik ta'lim uchun eng yaxshi tanlov.`),
		Platinum: premiumLandingPlan(`Intensiv`, `AI Tutor maksimal kunlik limitlar, chuqur roleplay, intensiv review va maksimal voice/pronunciation amaliyoti bilan.`,
			[]string{`maksimal kunlik limitlar`, `ko'proq voice/photo-context amaliyoti`, `zaif joylarni intensiv review qilish`, `heavy daily learning uchun eng yaxshi rejim`},
			nil,
			`Safar, ish, imtihon yoki juda zich temp uchun.`),
	},
	"tt": {
		Free: premiumLandingPlan(`Маршрутны сынау`, `Төп текстлы уку, сүзләр күнегүе, Phrasebook һәм үсеш күзәтүе. AI Tutor, listening, pronunciation һәм тавыш тикшерүләре Premiumда ачыла.`,
			[]string{`көндәлек гадәт һәм башлангыч дәресләр`, `төп сүз күнегүе`, `notes, phrasebook һәм үсеш күзәтүе`},
			[]string{`AI Tutor guided lessons`, `Listening һәм pronunciation`, `voice checks һәм photo tools`},
			`Продукт белән түләүсез танышу өчен туры килә.`),
		Premium: premiumLandingPlan(`Даими уку`, `Көндәлек практика өчен төп режим: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools һәм киңәйтелгән көнлек лимитлар.`,
			[]string{`тавышны текстка һәм ишетелгәнне тәрҗемә итү`, `рәсемдәге текстны тәрҗемә итү`, `тавыш яки фото контексты буенча практика`, `хаталар сүзлеге, notes, XP һәм streak`},
			nil,
			`Тотрыклы көндәлек уку өчен иң яхшы сайлау.`),
		Platinum: premiumLandingPlan(`Интенсив`, `AI Tutor максималь көнлек лимитлар, тирән roleplay, интенсив review һәм максималь voice/pronunciation практикасы белән.`,
			[]string{`максималь көнлек лимитлар`, `күбрәк voice/photo-context практикасы`, `зәгыйфь урыннарны интенсив review`, `heavy daily learning өчен иң яхшы режим`},
			nil,
			`Сәфәр, эш, имтихан яки бик тыгыз темп өчен.`),
	},
	"hy": {
		Free: premiumLandingPlan(`Փորձել ուղին`, `Հիմնական տեքստային ուսուցում, բառերի մարզում, Phrasebook և առաջընթացի տեսություն։ AI Tutor-ը, listening-ը, pronunciation-ը և ձայնային ստուգումները բացվում են Premium-ում։`,
			[]string{`ամենօրյա սովորություն և մեկնարկային դասեր`, `հիմնական բառերի մարզում`, `notes, phrasebook և առաջընթացի տեսություն`},
			[]string{`AI Tutor guided lessons`, `Listening և pronunciation`, `voice checks և photo tools`},
			`Հարմար է ապրանքը առանց վճարման ճանաչելու համար։`),
		Premium: premiumLandingPlan(`Կանոնավոր ուսուցում`, `Ամենօրյա պրակտիկայի հիմնական ռեժիմը՝ AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools և ընդլայնված օրական սահմաններ։`,
			[]string{`ձայնը տեքստի և լսածի թարգմանություն`, `նկարից տեքստի թարգմանություն`, `պրակտիկա ձայնի կամ նկարի համատեքստով`, `սխալների բառարան, notes, XP և streak`},
			nil,
			`Լավագույն ընտրությունը կայուն ամենօրյա ուսուցման համար։`),
		Platinum: premiumLandingPlan(`Ինտենսիվ`, `AI Tutor՝ առավելագույն օրական սահմաններով, խոր roleplay-ով, ինտենսիվ review-ով և առավելագույն voice/pronunciation պրակտիկայով։`,
			[]string{`առավելագույն օրական սահմաններ`, `ավելի շատ voice/photo-context պրակտիկա`, `թույլ կողմերի ինտենսիվ review`, `լավագույն ռեժիմ heavy daily learning-ի համար`},
			nil,
			`Ճամփորդության, աշխատանքի, քննության կամ շատ խիտ տեմպի համար։`),
	},
	"kk": {
		Free: premiumLandingPlan(`Бағытты байқап көру`, `Негізгі мәтіндік оқу, сөз жаттығуы, Phrasebook және прогресс шолуы. AI Tutor, listening, pronunciation және дауыс тексерулері Premium ішінде ашылады.`,
			[]string{`күнделікті әдет және бастапқы сабақтар`, `негізгі сөз жаттығуы`, `notes, phrasebook және прогресс шолуы`},
			[]string{`AI Tutor guided lessons`, `Listening және pronunciation`, `voice checks және photo tools`},
			`Өніммен төлемсіз танысуға жарайды.`),
		Premium: premiumLandingPlan(`Тұрақты оқу`, `Күнделікті практикаға арналған негізгі режим: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools және кеңейтілген күндік лимиттер.`,
			[]string{`дауысты мәтінге және естігенді аудару`, `суреттегі мәтінді аудару`, `дауыс немесе фото контексі бойынша практика`, `қателер сөздігі, notes, XP және streak`},
			nil,
			`Тұрақты күнделікті оқу үшін ең жақсы таңдау.`),
		Platinum: premiumLandingPlan(`Интенсив`, `AI Tutor ең жоғары күндік лимиттермен, терең roleplay, интенсивті review және максималды voice/pronunciation практикасымен.`,
			[]string{`ең жоғары күндік лимиттер`, `көбірек voice/photo-context практикасы`, `әлсіз жерлерді интенсивті review`, `heavy daily learning үшін ең жақсы режим`},
			nil,
			`Сапар, жұмыс, емтихан немесе өте тығыз қарқын үшін.`),
	},
	"ky": {
		Free: premiumLandingPlan(`Маршрутту сынап көрүү`, `Негизги тексттик окуу, сөз машыгуусу, Phrasebook жана прогресс обзору. AI Tutor, listening, pronunciation жана үн текшерүүлөр Premium ичинде ачылат.`,
			[]string{`күнүмдүк адат жана баштапкы сабактар`, `негизги сөз машыгуусу`, `notes, phrasebook жана прогресс обзору`},
			[]string{`AI Tutor guided lessons`, `Listening жана pronunciation`, `voice checks жана photo tools`},
			`Продукт менен акысыз таанышууга ылайыктуу.`),
		Premium: premiumLandingPlan(`Туруктуу окуу`, `Күнүмдүк практика үчүн негизги режим: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools жана кеңейтилген күндүк лимиттер.`,
			[]string{`үндү текстке жана укканды которуу`, `сүрөттөгү текстти которуу`, `үн же фото контексти боюнча практика`, `каталар сөздүгү, notes, XP жана streak`},
			nil,
			`Туруктуу күнүмдүк окуу үчүн эң жакшы тандоо.`),
		Platinum: premiumLandingPlan(`Интенсив`, `AI Tutor эң жогорку күндүк лимиттер, терең roleplay, интенсивдүү review жана максималдуу voice/pronunciation практикасы менен.`,
			[]string{`эң жогорку күндүк лимиттер`, `көбүрөөк voice/photo-context практикасы`, `алсыз жерлерди интенсивдүү review`, `heavy daily learning үчүн эң жакшы режим`},
			nil,
			`Сапар, иш, экзамен же абдан тыгыз темп үчүн.`),
	},
	"ka": {
		Free: premiumLandingPlan(`მარშრუტის მოსინჯვა`, `საბაზო ტექსტური სწავლა, სიტყვების ვარჯიში, Phrasebook და პროგრესის მიმოხილვა. AI Tutor, listening, pronunciation და ხმოვანი შემოწმებები Premium-ში იხსნება.`,
			[]string{`ყოველდღიური ჩვევა და საწყისი გაკვეთილები`, `საბაზო სიტყვების ვარჯიში`, `notes, phrasebook და პროგრესის მიმოხილვა`},
			[]string{`AI Tutor guided lessons`, `Listening და pronunciation`, `voice checks და photo tools`},
			`შესაფერისია პროდუქტის გასაცნობად გადახდის გარეშე.`),
		Premium: premiumLandingPlan(`რეგულარული სწავლა`, `ყოველდღიური პრაქტიკის მთავარი რეჟიმი: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools და გაფართოებული დღიური ლიმიტები.`,
			[]string{`ხმის ტექსტად გადაყვანა და მოსმენილის თარგმნა`, `ტექსტის თარგმნა სურათიდან`, `პრაქტიკა ხმის ან ფოტოს კონტექსტით`, `შეცდომების ლექსიკონი, notes, XP და streak`},
			nil,
			`საუკეთესო არჩევანი სტაბილური ყოველდღიური სწავლისთვის.`),
		Platinum: premiumLandingPlan(`ინტენსივი`, `AI Tutor მაქსიმალური დღიური ლიმიტებით, ღრმა roleplay-ით, ინტენსიური review-ით და მაქსიმალური voice/pronunciation პრაქტიკით.`,
			[]string{`მაქსიმალური დღიური ლიმიტები`, `მეტი voice/photo-context პრაქტიკა`, `სუსტი ადგილების ინტენსიური review`, `საუკეთესო რეჟიმი heavy daily learning-ისთვის`},
			nil,
			`მოგზაურობისთვის, სამუშაოსთვის, გამოცდისთვის ან ძალიან მკვრივი ტემპისთვის.`),
	},
	"uk": {
		Free: premiumLandingPlan(`Спробувати маршрут`, `Базове текстове навчання, тренування слів, Phrasebook і огляд прогресу. AI Tutor, listening, pronunciation і голосові перевірки відкриваються в Premium.`,
			[]string{`щоденна звичка і стартові уроки`, `базове тренування слів`, `notes, phrasebook і огляд прогресу`},
			[]string{`AI Tutor guided lessons`, `Listening і pronunciation`, `voice checks і photo tools`},
			`Підходить для знайомства з продуктом без оплати.`),
		Premium: premiumLandingPlan(`Регулярне навчання`, `Основний режим для щоденної практики: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools і розширені денні ліміти.`,
			[]string{`голос у текст і переклад почутого`, `переклад тексту з картинки`, `практика за контекстом голосу або фото`, `словник помилок, notes, XP і streak`},
			nil,
			`Найкращий вибір для стабільного щоденного навчання.`),
		Platinum: premiumLandingPlan(`Інтенсив`, `AI Tutor з максимальними денними лімітами, глибиною roleplay, інтенсивним review і максимальною voice/pronunciation практикою.`,
			[]string{`максимальні денні ліміти`, `більше voice/photo-context практики`, `інтенсивний review слабких місць`, `найкращий режим для heavy daily learning`},
			nil,
			`Для поїздки, роботи, іспиту або дуже щільного темпу.`),
	},
	"pl": {
		Free: premiumLandingPlan(`Wypróbuj ścieżkę`, `Podstawowa nauka tekstowa, trening słów, Phrasebook i przegląd postępów. AI Tutor, listening, pronunciation i sprawdzanie głosu otwierają się w Premium.`,
			[]string{`codzienny nawyk i lekcje startowe`, `podstawowy trening słów`, `notes, phrasebook i przegląd postępów`},
			[]string{`lekcje prowadzone przez AI Tutor`, `Listening i pronunciation`, `sprawdzanie głosu i narzędzia foto`},
			`Dobre do poznania produktu bez płatności.`),
		Premium: premiumLandingPlan(`Regularna nauka`, `Główny tryb do codziennej praktyki: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools i rozszerzone dzienne limity.`,
			[]string{`głos na tekst i tłumaczenie usłyszanego`, `tłumaczenie tekstu z obrazu`, `praktyka z kontekstu głosu lub zdjęcia`, `słownik błędów, notes, XP i streak`},
			nil,
			`Najlepszy wybór do stabilnej codziennej nauki.`),
		Platinum: premiumLandingPlan(`Intensywnie`, `AI Tutor z maksymalnymi dziennymi limitami, głębszym roleplay, intensywnym review i maksymalną praktyką voice/pronunciation.`,
			[]string{`maksymalne dzienne limity`, `więcej praktyki voice/photo-context`, `intensywny review słabych miejsc`, `najlepszy tryb dla heavy daily learning`},
			nil,
			`Na podróż, pracę, egzamin lub bardzo gęste tempo.`),
	},
	"ro": {
		Free: premiumLandingPlan(`Încearcă traseul`, `Învățare text de bază, antrenament de cuvinte, Phrasebook și privire de ansamblu asupra progresului. AI Tutor, listening, pronunciation și verificările vocale se deschid în Premium.`,
			[]string{`obicei zilnic și lecții de început`, `antrenament de cuvinte de bază`, `notes, phrasebook și privire de ansamblu asupra progresului`},
			[]string{`lecții ghidate AI Tutor`, `Listening și pronunciation`, `verificări vocale și instrumente foto`},
			`Potrivit pentru a cunoaște produsul fără plată.`),
		Premium: premiumLandingPlan(`Studiu regulat`, `Modul principal pentru practică zilnică: AI Tutor, listening, pronunciation, lecții AI ghidate, voice checks, photo tools și limite zilnice extinse.`,
			[]string{`voce în text și traducerea a ceea ce ai auzit`, `traducerea textului dintr-o imagine`, `practică după context vocal sau foto`, `dicționar de greșeli, notes, XP și streak`},
			nil,
			`Cea mai bună alegere pentru învățare zilnică stabilă.`),
		Platinum: premiumLandingPlan(`Intensiv`, `AI Tutor cu limite zilnice maxime, roleplay mai profund, review intensiv și practică voice/pronunciation maximă.`,
			[]string{`limite zilnice maxime`, `mai multă practică voice/photo-context`, `review intensiv al punctelor slabe`, `cel mai bun mod pentru heavy daily learning`},
			nil,
			`Pentru călătorie, muncă, examen sau un ritm foarte dens.`),
	},
	"pt": {
		Free: premiumLandingPlan(`Experimentar rota`, `Aprendizagem básica por texto, treino de palavras, Phrasebook e visão geral do progresso. AI Tutor, listening, pronunciation e verificações de voz abrem no Premium.`,
			[]string{`hábito diário e lições iniciais`, `treino básico de palavras`, `notes, phrasebook e visão geral do progresso`},
			[]string{`lições guiadas do AI Tutor`, `Listening e pronunciation`, `verificações de voz e ferramentas de foto`},
			`Indicado para conhecer o produto sem pagar.`),
		Premium: premiumLandingPlan(`Estudo regular`, `Modo principal para prática diária: AI Tutor, listening, pronunciation, aulas guiadas por IA, voice checks, photo tools e limites diários ampliados.`,
			[]string{`voz em texto e tradução do que foi ouvido`, `tradução de texto de uma imagem`, `prática pelo contexto de voz ou foto`, `dicionário de erros, notes, XP e streak`},
			nil,
			`A melhor escolha para aprendizagem diária estável.`),
		Platinum: premiumLandingPlan(`Intensivo`, `AI Tutor com limites diários máximos, roleplay mais profundo, review intensivo e prática máxima de voice/pronunciation.`,
			[]string{`limites diários máximos`, `mais prática voice/photo-context`, `review intensivo dos pontos fracos`, `melhor modo para heavy daily learning`},
			nil,
			`Para viagem, trabalho, exame ou um ritmo muito intenso.`),
	},
	"ar": {
		Free: premiumLandingPlan(`جرّب المسار`, `تعلّم نصي أساسي، تدريب كلمات، Phrasebook ونظرة على التقدم. AI Tutor وlistening وpronunciation وفحوصات الصوت تفتح في Premium.`,
			[]string{`عادة يومية ودروس بداية`, `تدريب كلمات أساسي`, `notes وphrasebook ونظرة على التقدم`},
			[]string{`دروس AI Tutor الموجهة`, `Listening وpronunciation`, `فحوصات صوت وأدوات صور`},
			`مناسب للتعرّف على المنتج من دون دفع.`),
		Premium: premiumLandingPlan(`دراسة منتظمة`, `الوضع الأساسي للتدريب اليومي: AI Tutor وlistening وpronunciation ودروس AI موجهة وvoice checks وphoto tools وحدود يومية موسعة.`,
			[]string{`تحويل الصوت إلى نص وترجمة ما سمعته`, `ترجمة النص من الصورة`, `تدريب حسب سياق الصوت أو الصورة`, `قاموس الأخطاء وnotes وXP وstreak`},
			nil,
			`الخيار الأفضل لتعلّم يومي ثابت.`),
		Platinum: premiumLandingPlan(`مكثف`, `AI Tutor مع أعلى حدود يومية وroleplay أعمق وreview مكثف وأقصى تدريب voice/pronunciation.`,
			[]string{`أعلى حدود يومية`, `مزيد من تدريب voice/photo-context`, `review مكثف لنقاط الضعف`, `أفضل وضع لـ heavy daily learning`},
			nil,
			`للسفر أو العمل أو الامتحان أو الوتيرة المكثفة جدًا.`),
	},
	"bn": {
		Free: premiumLandingPlan(`পথটি চেষ্টা করুন`, `মৌলিক টেক্সট শেখা, শব্দ অনুশীলন, Phrasebook এবং অগ্রগতি দেখা। AI Tutor, listening, pronunciation এবং ভয়েস চেক Premium-এ খুলবে।`,
			[]string{`দৈনিক অভ্যাস এবং শুরুর পাঠ`, `মৌলিক শব্দ অনুশীলন`, `notes, phrasebook এবং অগ্রগতি দেখা`},
			[]string{`AI Tutor guided lessons`, `Listening এবং pronunciation`, `voice checks এবং photo tools`},
			`পেমেন্ট ছাড়া পণ্যটি চিনে নেওয়ার জন্য উপযুক্ত।`),
		Premium: premiumLandingPlan(`নিয়মিত শেখা`, `দৈনিক অনুশীলনের প্রধান মোড: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools এবং বাড়তি দৈনিক সীমা।`,
			[]string{`ভয়েস থেকে টেক্সট এবং শোনা কথার অনুবাদ`, `ছবি থেকে টেক্সট অনুবাদ`, `ভয়েস বা ছবির প্রসঙ্গে অনুশীলন`, `ভুলের অভিধান, notes, XP এবং streak`},
			nil,
			`স্থির দৈনিক শেখার জন্য সেরা পছন্দ।`),
		Platinum: premiumLandingPlan(`নিবিড়`, `AI Tutor সর্বোচ্চ দৈনিক সীমা, গভীর roleplay, নিবিড় review এবং সর্বোচ্চ voice/pronunciation অনুশীলনসহ।`,
			[]string{`সর্বোচ্চ দৈনিক সীমা`, `আরও voice/photo-context অনুশীলন`, `দুর্বল জায়গার নিবিড় review`, `heavy daily learning-এর সেরা মোড`},
			nil,
			`ভ্রমণ, কাজ, পরীক্ষা বা খুব ঘন গতির জন্য।`),
	},
	"cs": {
		Free: premiumLandingPlan(`Vyzkoušet trasu`, `Základní textové učení, trénink slov, Phrasebook a přehled pokroku. AI Tutor, listening, pronunciation a hlasové kontroly se otevřou v Premium.`,
			[]string{`denní návyk a úvodní lekce`, `základní trénink slov`, `notes, phrasebook a přehled pokroku`},
			[]string{`vedené lekce AI Tutor`, `Listening a pronunciation`, `hlasové kontroly a foto nástroje`},
			`Vhodné pro seznámení s produktem bez platby.`),
		Premium: premiumLandingPlan(`Pravidelné učení`, `Hlavní režim pro každodenní praxi: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools a rozšířené denní limity.`,
			[]string{`hlas na text a překlad slyšeného`, `překlad textu z obrázku`, `praxe podle kontextu hlasu nebo fotky`, `slovník chyb, notes, XP a streak`},
			nil,
			`Nejlepší volba pro stabilní každodenní učení.`),
		Platinum: premiumLandingPlan(`Intenzivní`, `AI Tutor s maximálními denními limity, hlubším roleplay, intenzivním review a maximální voice/pronunciation praxí.`,
			[]string{`maximální denní limity`, `více voice/photo-context praxe`, `intenzivní review slabých míst`, `nejlepší režim pro heavy daily learning`},
			nil,
			`Pro cestu, práci, zkoušku nebo velmi husté tempo.`),
	},
	"el": {
		Free: premiumLandingPlan(`Δοκιμή διαδρομής`, `Βασική μάθηση με κείμενο, εξάσκηση λέξεων, Phrasebook και επισκόπηση προόδου. AI Tutor, listening, pronunciation και φωνητικοί έλεγχοι ανοίγουν στο Premium.`,
			[]string{`καθημερινή συνήθεια και αρχικά μαθήματα`, `βασική εξάσκηση λέξεων`, `notes, phrasebook και επισκόπηση προόδου`},
			[]string{`καθοδηγούμενα μαθήματα AI Tutor`, `Listening και pronunciation`, `φωνητικοί έλεγχοι και εργαλεία φωτογραφίας`},
			`Κατάλληλο για γνωριμία με το προϊόν χωρίς πληρωμή.`),
		Premium: premiumLandingPlan(`Τακτική μελέτη`, `Ο βασικός τρόπος για καθημερινή πρακτική: AI Tutor, listening, pronunciation, καθοδηγούμενα AI lessons, voice checks, photo tools και διευρυμένα ημερήσια όρια.`,
			[]string{`φωνή σε κείμενο και μετάφραση όσων άκουσες`, `μετάφραση κειμένου από εικόνα`, `πρακτική με βάση φωνητικό ή φωτογραφικό πλαίσιο`, `λεξικό λαθών, notes, XP και streak`},
			nil,
			`Η καλύτερη επιλογή για σταθερή καθημερινή μάθηση.`),
		Platinum: premiumLandingPlan(`Εντατικό`, `AI Tutor με μέγιστα ημερήσια όρια, βαθύτερο roleplay, εντατικό review και μέγιστη voice/pronunciation πρακτική.`,
			[]string{`μέγιστα ημερήσια όρια`, `περισσότερη voice/photo-context πρακτική`, `εντατικό review αδύναμων σημείων`, `καλύτερος τρόπος για heavy daily learning`},
			nil,
			`Για ταξίδι, εργασία, εξέταση ή πολύ πυκνό ρυθμό.`),
	},
	"hi": {
		Free: premiumLandingPlan(`मार्ग आज़माएँ`, `बुनियादी टेक्स्ट सीखना, शब्द अभ्यास, Phrasebook और प्रगति अवलोकन। AI Tutor, listening, pronunciation और voice checks Premium में खुलते हैं।`,
			[]string{`दैनिक आदत और शुरुआती पाठ`, `बुनियादी शब्द अभ्यास`, `notes, phrasebook और प्रगति अवलोकन`},
			[]string{`AI Tutor guided lessons`, `Listening और pronunciation`, `voice checks और photo tools`},
			`बिना भुगतान उत्पाद को समझने के लिए उपयुक्त।`),
		Premium: premiumLandingPlan(`नियमित अध्ययन`, `दैनिक अभ्यास का मुख्य मोड: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools और बढ़ी हुई दैनिक सीमाएँ।`,
			[]string{`आवाज़ को टेक्स्ट और सुनी हुई बात का अनुवाद`, `चित्र से टेक्स्ट अनुवाद`, `आवाज़ या फोटो संदर्भ से अभ्यास`, `गलतियों का शब्दकोश, notes, XP और streak`},
			nil,
			`स्थिर दैनिक सीखने के लिए सबसे अच्छा विकल्प।`),
		Platinum: premiumLandingPlan(`गहन`, `AI Tutor अधिकतम दैनिक सीमाओं, गहरे roleplay, गहन review और अधिकतम voice/pronunciation अभ्यास के साथ।`,
			[]string{`अधिकतम दैनिक सीमाएँ`, `अधिक voice/photo-context अभ्यास`, `कमज़ोरियों का गहन review`, `heavy daily learning के लिए सबसे अच्छा मोड`},
			nil,
			`यात्रा, काम, परीक्षा या बहुत घने अभ्यास-ताल के लिए।`),
	},
	"hu": {
		Free: premiumLandingPlan(`Útvonal kipróbálása`, `Alap szöveges tanulás, szógyakorlás, Phrasebook és fejlődési áttekintés. Az AI Tutor, listening, pronunciation és hangellenőrzések Premiumban nyílnak meg.`,
			[]string{`napi szokás és kezdő leckék`, `alap szógyakorlás`, `notes, phrasebook és fejlődési áttekintés`},
			[]string{`AI Tutor guided lessons`, `Listening és pronunciation`, `voice checks és photo tools`},
			`Alkalmas a termék megismerésére fizetés nélkül.`),
		Premium: premiumLandingPlan(`Rendszeres tanulás`, `A fő mód a napi gyakorláshoz: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools és bővített napi limitek.`,
			[]string{`hangból szöveg és a hallottak fordítása`, `szöveg fordítása képről`, `gyakorlás hang- vagy fotókontekstusból`, `hibaszótár, notes, XP és streak`},
			nil,
			`A legjobb választás stabil napi tanuláshoz.`),
		Platinum: premiumLandingPlan(`Intenzív`, `AI Tutor maximális napi limitekkel, mélyebb roleplay-jel, intenzív review-val és maximális voice/pronunciation gyakorlással.`,
			[]string{`maximális napi limitek`, `több voice/photo-context gyakorlás`, `gyenge pontok intenzív review-ja`, `legjobb mód heavy daily learninghez`},
			nil,
			`Utazáshoz, munkához, vizsgához vagy nagyon sűrű tempóhoz.`),
	},
	"id": {
		Free: premiumLandingPlan(`Coba jalur`, `Belajar teks dasar, latihan kata, Phrasebook, dan ringkasan progres. AI Tutor, listening, pronunciation, dan pemeriksaan suara terbuka di Premium.`,
			[]string{`kebiasaan harian dan pelajaran awal`, `latihan kata dasar`, `notes, phrasebook, dan ringkasan progres`},
			[]string{`pelajaran terpandu AI Tutor`, `Listening dan pronunciation`, `pemeriksaan suara dan alat foto`},
			`Cocok untuk mengenal produk tanpa membayar.`),
		Premium: premiumLandingPlan(`Belajar rutin`, `Mode utama untuk praktik harian: AI Tutor, listening, pronunciation, pelajaran AI terpandu, voice checks, photo tools, dan batas harian yang diperluas.`,
			[]string{`suara ke teks dan terjemahan yang didengar`, `terjemahan teks dari gambar`, `praktik dari konteks suara atau foto`, `kamus kesalahan, notes, XP, dan streak`},
			nil,
			`Pilihan terbaik untuk belajar harian yang stabil.`),
		Platinum: premiumLandingPlan(`Intensif`, `AI Tutor dengan batas harian maksimum, roleplay lebih dalam, review intensif, dan praktik voice/pronunciation maksimal.`,
			[]string{`batas harian maksimum`, `lebih banyak praktik voice/photo-context`, `review intensif titik lemah`, `mode terbaik untuk heavy daily learning`},
			nil,
			`Untuk perjalanan, kerja, ujian, atau tempo yang sangat padat.`),
	},
	"nl": {
		Free: premiumLandingPlan(`Route proberen`, `Basis leren met tekst, woordtraining, Phrasebook en voortgangsoverzicht. AI Tutor, listening, pronunciation en stemchecks openen in Premium.`,
			[]string{`dagelijkse gewoonte en startlessen`, `basis woordtraining`, `notes, phrasebook en voortgangsoverzicht`},
			[]string{`AI Tutor guided lessons`, `Listening en pronunciation`, `stemchecks en fototools`},
			`Geschikt om het product zonder betaling te leren kennen.`),
		Premium: premiumLandingPlan(`Regelmatig leren`, `De hoofdmodus voor dagelijkse oefening: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools en ruimere daglimieten.`,
			[]string{`spraak naar tekst en vertaling van wat je hoorde`, `tekst uit een afbeelding vertalen`, `oefenen vanuit stem- of fotocontext`, `foutenwoordenboek, notes, XP en streak`},
			nil,
			`De beste keuze voor stabiel dagelijks leren.`),
		Platinum: premiumLandingPlan(`Intensief`, `AI Tutor met maximale daglimieten, diepere roleplay, intensieve review en maximale voice/pronunciation-oefening.`,
			[]string{`maximale daglimieten`, `meer voice/photo-context oefening`, `intensieve review van zwakke plekken`, `beste modus voor heavy daily learning`},
			nil,
			`Voor reizen, werk, examen of een zeer dicht tempo.`),
	},
	"sv": {
		Free: premiumLandingPlan(`Prova rutten`, `Grundläggande textinlärning, ordträning, Phrasebook och översikt över framsteg. AI Tutor, listening, pronunciation och röstkontroller öppnas i Premium.`,
			[]string{`daglig vana och startlektioner`, `grundläggande ordträning`, `notes, phrasebook och översikt över framsteg`},
			[]string{`AI Tutor guided lessons`, `Listening och pronunciation`, `röstkontroller och fotoverktyg`},
			`Passar för att lära känna produkten utan betalning.`),
		Premium: premiumLandingPlan(`Regelbunden studie`, `Huvudläget för daglig träning: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools och utökade dagliga gränser.`,
			[]string{`röst till text och översättning av det du hörde`, `översättning av text från bild`, `träning från röst- eller fotokontext`, `misstagsordbok, notes, XP och streak`},
			nil,
			`Bästa valet för stabil daglig inlärning.`),
		Platinum: premiumLandingPlan(`Intensivt`, `AI Tutor med maximala dagliga gränser, djupare roleplay, intensiv review och maximal voice/pronunciation-träning.`,
			[]string{`maximala dagliga gränser`, `mer voice/photo-context träning`, `intensiv review av svaga punkter`, `bästa läget för heavy daily learning`},
			nil,
			`För resa, arbete, examen eller mycket tätt tempo.`),
	},
	"ta": {
		Free: premiumLandingPlan(`பாதையை முயற்சி செய்`, `அடிப்படை உரை கற்றல், சொல் பயிற்சி, Phrasebook மற்றும் முன்னேற்ற கண்ணோட்டம். AI Tutor, listening, pronunciation மற்றும் குரல் சரிபார்ப்புகள் Premium-இல் திறக்கும்.`,
			[]string{`தினசரி பழக்கம் மற்றும் தொடக்கப் பாடங்கள்`, `அடிப்படை சொல் பயிற்சி`, `notes, phrasebook மற்றும் முன்னேற்ற கண்ணோட்டம்`},
			[]string{`AI Tutor guided lessons`, `Listening மற்றும் pronunciation`, `voice checks மற்றும் photo tools`},
			`கட்டணம் இல்லாமல் தயாரிப்பை அறிய ஏற்றது.`),
		Premium: premiumLandingPlan(`வழக்கமான கற்றல்`, `தினசரி பயிற்சிக்கான முக்கிய முறை: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools மற்றும் விரிவான தினசரி வரம்புகள்.`,
			[]string{`குரலை உரையாக மாற்றி கேட்டதை மொழிபெயர்ப்பு`, `படத்திலிருந்து உரை மொழிபெயர்ப்பு`, `குரல் அல்லது புகைப்பட சூழலால் பயிற்சி`, `பிழை அகராதி, notes, XP மற்றும் streak`},
			nil,
			`நிலையான தினசரி கற்றலுக்கான சிறந்த தேர்வு.`),
		Platinum: premiumLandingPlan(`தீவிரம்`, `AI Tutor அதிகபட்ச தினசரி வரம்புகள், ஆழமான roleplay, தீவிர review மற்றும் அதிகபட்ச voice/pronunciation பயிற்சியுடன்.`,
			[]string{`அதிகபட்ச தினசரி வரம்புகள்`, `மேலும் voice/photo-context பயிற்சி`, `பலவீனங்களை தீவிரமாக review செய்வது`, `heavy daily learning-க்கு சிறந்த முறை`},
			nil,
			`பயணம், வேலை, தேர்வு அல்லது மிக நெருக்கமான வேகத்திற்கு.`),
	},
	"te": {
		Free: premiumLandingPlan(`మార్గాన్ని ప్రయత్నించండి`, `ప్రాథమిక టెక్స్ట్ అభ్యాసం, పదాల సాధన, Phrasebook మరియు పురోగతి అవలోకనం. AI Tutor, listening, pronunciation మరియు వాయిస్ చెక్స్ Premiumలో తెరుచుకుంటాయి.`,
			[]string{`రోజువారీ అలవాటు మరియు ప్రారంభ పాఠాలు`, `ప్రాథమిక పదాల సాధన`, `notes, phrasebook మరియు పురోగతి అవలోకనం`},
			[]string{`AI Tutor guided lessons`, `Listening మరియు pronunciation`, `voice checks మరియు photo tools`},
			`చెల్లింపు లేకుండా ఉత్పత్తిని తెలుసుకోవడానికి సరిపోతుంది.`),
		Premium: premiumLandingPlan(`నియమిత అభ్యాసం`, `రోజువారీ సాధనకు ప్రధాన మోడ్: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools మరియు విస్తరించిన రోజువారీ పరిమితులు.`,
			[]string{`వాయిస్‌ను టెక్స్ట్‌గా మార్చి విన్నదాన్ని అనువదించడం`, `చిత్రం నుంచి టెక్స్ట్ అనువాదం`, `వాయిస్ లేదా ఫోటో సందర్భంతో సాధన`, `తప్పుల నిఘంటువు, notes, XP మరియు streak`},
			nil,
			`స్థిరమైన రోజువారీ అభ్యాసానికి ఉత్తమ ఎంపిక.`),
		Platinum: premiumLandingPlan(`ఇంటెన్సివ్`, `AI Tutor గరిష్ఠ రోజువారీ పరిమితులు, లోతైన roleplay, తీవ్ర review మరియు గరిష్ఠ voice/pronunciation సాధనతో.`,
			[]string{`గరిష్ఠ రోజువారీ పరిమితులు`, `మరింత voice/photo-context సాధన`, `బలహీన ప్రాంతాల తీవ్ర review`, `heavy daily learning కోసం ఉత్తమ మోడ్`},
			nil,
			`ప్రయాణం, పని, పరీక్ష లేదా చాలా గట్టి వేగానికి.`),
	},
	"th": {
		Free: premiumLandingPlan(`ทดลองเส้นทาง`, `การเรียนด้วยข้อความพื้นฐาน ฝึกคำศัพท์ Phrasebook และภาพรวมความก้าวหน้า AI Tutor, listening, pronunciation และการตรวจเสียงจะเปิดใน Premium`,
			[]string{`นิสัยประจำวันและบทเรียนเริ่มต้น`, `ฝึกคำศัพท์พื้นฐาน`, `notes, phrasebook และภาพรวมความก้าวหน้า`},
			[]string{`AI Tutor guided lessons`, `Listening และ pronunciation`, `voice checks และ photo tools`},
			`เหมาะสำหรับทำความรู้จักผลิตภัณฑ์โดยไม่ต้องจ่ายเงิน`),
		Premium: premiumLandingPlan(`เรียนสม่ำเสมอ`, `โหมดหลักสำหรับฝึกทุกวัน: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools และขีดจำกัดรายวันที่เพิ่มขึ้น`,
			[]string{`เสียงเป็นข้อความและแปลสิ่งที่ได้ยิน`, `แปลข้อความจากรูปภาพ`, `ฝึกตามบริบทจากเสียงหรือรูปภาพ`, `พจนานุกรมข้อผิดพลาด, notes, XP และ streak`},
			nil,
			`ตัวเลือกที่ดีที่สุดสำหรับการเรียนทุกวันที่สม่ำเสมอ`),
		Platinum: premiumLandingPlan(`เข้มข้น`, `AI Tutor พร้อมขีดจำกัดรายวันสูงสุด, roleplay ที่ลึกขึ้น, review เข้มข้น และฝึก voice/pronunciation สูงสุด`,
			[]string{`ขีดจำกัดรายวันสูงสุด`, `ฝึก voice/photo-context มากขึ้น`, `review จุดอ่อนแบบเข้มข้น`, `โหมดที่ดีที่สุดสำหรับ heavy daily learning`},
			nil,
			`สำหรับการเดินทาง งาน สอบ หรือจังหวะที่แน่นมาก`),
	},
	"tl": {
		Free: premiumLandingPlan(`Subukan ang ruta`, `Basic na pag-aaral sa text, word training, Phrasebook, at overview ng progreso. AI Tutor, listening, pronunciation, at voice checks ay bukas sa Premium.`,
			[]string{`araw-araw na habit at panimulang lessons`, `basic word training`, `notes, phrasebook, at overview ng progreso`},
			[]string{`AI Tutor guided lessons`, `Listening at pronunciation`, `voice checks at photo tools`},
			`Bagay para makilala ang produkto nang walang bayad.`),
		Premium: premiumLandingPlan(`Regular na pag-aaral`, `Pangunahing mode para sa daily practice: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools, at pinalawak na daily limits.`,
			[]string{`voice to text at translation ng narinig`, `translation ng text mula sa image`, `practice mula sa voice o photo context`, `mistake dictionary, notes, XP, at streak`},
			nil,
			`Pinakamagandang piliin para sa steady daily learning.`),
		Platinum: premiumLandingPlan(`Masinsinan`, `AI Tutor na may maximum daily limits, mas malalim na roleplay, intensive review, at maximum voice/pronunciation practice.`,
			[]string{`maximum daily limits`, `mas maraming voice/photo-context practice`, `intensive review ng weak spots`, `best mode para sa heavy daily learning`},
			nil,
			`Para sa biyahe, trabaho, exam, o napakasiksik na tempo.`),
	},
	"tr": {
		Free: premiumLandingPlan(`Rotayı dene`, `Temel metinle öğrenme, kelime çalışması, Phrasebook ve ilerleme özeti. AI Tutor, listening, pronunciation ve ses kontrolleri Premium'da açılır.`,
			[]string{`günlük alışkanlık ve başlangıç dersleri`, `temel kelime çalışması`, `notes, phrasebook ve ilerleme özeti`},
			[]string{`AI Tutor guided lessons`, `Listening ve pronunciation`, `voice checks ve photo tools`},
			`Ürünü ödeme yapmadan tanımak için uygundur.`),
		Premium: premiumLandingPlan(`Düzenli çalışma`, `Günlük pratik için ana mod: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools ve genişletilmiş günlük limitler.`,
			[]string{`sesi metne çevirme ve duyduğunu tercüme etme`, `görüntüden metin çevirisi`, `ses veya fotoğraf bağlamına göre pratik`, `hata sözlüğü, notes, XP ve streak`},
			nil,
			`İstikrarlı günlük öğrenme için en iyi seçim.`),
		Platinum: premiumLandingPlan(`Yoğun`, `AI Tutor maksimum günlük limitler, daha derin roleplay, yoğun review ve maksimum voice/pronunciation pratiğiyle.`,
			[]string{`maksimum günlük limitler`, `daha fazla voice/photo-context pratiği`, `zayıf noktaların yoğun review'u`, `heavy daily learning için en iyi mod`},
			nil,
			`Seyahat, iş, sınav veya çok yoğun tempo için.`),
	},
	"vi": {
		Free: premiumLandingPlan(`Thử lộ trình`, `Học văn bản cơ bản, luyện từ, Phrasebook và tổng quan tiến độ. AI Tutor, listening, pronunciation và kiểm tra giọng nói mở trong Premium.`,
			[]string{`thói quen hằng ngày và bài học khởi đầu`, `luyện từ cơ bản`, `notes, phrasebook và tổng quan tiến độ`},
			[]string{`bài học có AI Tutor hướng dẫn`, `Listening và pronunciation`, `kiểm tra giọng nói và công cụ ảnh`},
			`Phù hợp để làm quen với sản phẩm mà không cần trả phí.`),
		Premium: premiumLandingPlan(`Học đều đặn`, `Chế độ chính cho luyện tập hằng ngày: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools và giới hạn hằng ngày mở rộng.`,
			[]string{`giọng nói thành văn bản và dịch nội dung đã nghe`, `dịch văn bản từ ảnh`, `luyện tập theo ngữ cảnh giọng nói hoặc ảnh`, `từ điển lỗi, notes, XP và streak`},
			nil,
			`Lựa chọn tốt nhất cho việc học hằng ngày ổn định.`),
		Platinum: premiumLandingPlan(`Chuyên sâu`, `AI Tutor với giới hạn hằng ngày tối đa, roleplay sâu hơn, review chuyên sâu và luyện voice/pronunciation tối đa.`,
			[]string{`giới hạn hằng ngày tối đa`, `nhiều luyện tập voice/photo-context hơn`, `review chuyên sâu điểm yếu`, `chế độ tốt nhất cho heavy daily learning`},
			nil,
			`Cho chuyến đi, công việc, kỳ thi hoặc nhịp học rất dày.`),
	},
}

func premiumPlanLandingCopy(user userState, tier string) premiumLandingPlanCopy {
	code := normalizeInterfaceLanguage(user.InterfaceLanguage)
	set, ok := premiumLandingCopyByLanguage[code]
	if !ok {
		set = premiumLandingCopyByLanguage["en"]
	}

	switch normalizedPremiumLandingTier(tier) {
	case "free":
		return clonePremiumLandingPlanCopy(set.Free)
	case "platinum":
		return clonePremiumLandingPlanCopy(set.Platinum)
	default:
		return clonePremiumLandingPlanCopy(set.Premium)
	}
}

func normalizedPremiumLandingTier(tier string) string {
	tier = strings.ToLower(strings.TrimSpace(tier))
	switch {
	case strings.Contains(tier, "free") || strings.Contains(tier, "basic"):
		return "free"
	case strings.Contains(tier, "platinum"):
		return "platinum"
	default:
		return "premium"
	}
}

func clonePremiumLandingPlanCopy(copy premiumLandingPlanCopy) premiumLandingPlanCopy {
	copy.Included = append([]string(nil), copy.Included...)
	copy.Locked = append([]string(nil), copy.Locked...)
	return copy
}

func premiumPlanLandingDescription(user userState, tier string) string {
	return premiumPlanLandingCopy(user, tier).Description
}

func premiumLandingMarkdown(copy premiumLandingPlanCopy) string {
	return premiumLandingMarkdownWithLimits(copy, nil)
}

func premiumLandingMarkdownWithLimits(copy premiumLandingPlanCopy, limits []string) string {
	var builder strings.Builder
	builder.WriteString(escapeMarkdownV2(copy.Description))
	for _, item := range limits {
		builder.WriteString("\n• ")
		builder.WriteString(escapeMarkdownV2(item))
	}
	for _, item := range copy.Included {
		builder.WriteString("\n• ")
		builder.WriteString(escapeMarkdownV2(item))
	}
	for _, item := range copy.Locked {
		builder.WriteString("\n• ")
		builder.WriteString(escapeMarkdownV2(item))
	}
	if copy.Note != "" {
		builder.WriteString("\n_")
		builder.WriteString(escapeMarkdownV2(copy.Note))
		builder.WriteString("_")
	}
	return builder.String()
}
