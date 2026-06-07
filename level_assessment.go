package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type levelQuestion struct {
	Question     string
	Options      []string
	CorrectIndex int
	Points       int
}

var vocabularyLevelAssessmentCache = struct {
	mu        sync.RWMutex
	questions map[string][]levelQuestion
}{
	questions: map[string][]levelQuestion{},
}

var levelAssessmentQuestions = []levelQuestion{
	{
		Question:     "Выбери правильное предложение:",
		Options:      []string{"I am student.", "I am a student.", "I student.", "I be student."},
		CorrectIndex: 1,
		Points:       1,
	},
	{
		Question:     "Выбери правильный перевод: Я обычно завтракаю в 8.",
		Options:      []string{"I usually have breakfast at 8.", "I have usually breakfast in 8.", "I am usually breakfast at 8.", "I breakfast usually on 8."},
		CorrectIndex: 0,
		Points:       2,
	},
	{
		Question:     "Заполни пропуск: She _____ to London last year.",
		Options:      []string{"go", "goes", "went", "has go"},
		CorrectIndex: 2,
		Points:       2,
	},
	{
		Question:     "Выбери лучший вариант: I have lived here _____ 2020.",
		Options:      []string{"for", "since", "during", "from"},
		CorrectIndex: 1,
		Points:       3,
	},
	{
		Question:     "Какое предложение звучит естественно?",
		Options:      []string{"I look forward to hear from you.", "I look forward hearing from you.", "I look forward to hearing from you.", "I look forward hear from you."},
		CorrectIndex: 2,
		Points:       4,
	},
	{
		Question:     "Выбери ближайшее значение слова: reliable",
		Options:      []string{"easy to break", "able to be trusted", "very expensive", "hard to understand"},
		CorrectIndex: 1,
		Points:       3,
	},
	{
		Question:     "Заполни пропуск: If I _____ more time, I would travel more.",
		Options:      []string{"have", "had", "will have", "would have"},
		CorrectIndex: 1,
		Points:       4,
	},
	{
		Question:     "Выбери подходящую связку: The task was difficult; _____, we finished it on time.",
		Options:      []string{"therefore", "however", "because", "unless"},
		CorrectIndex: 1,
		Points:       4,
	},
	{
		Question:     "Выбери правильную форму пассива для предложения: They are building a new bridge.",
		Options:      []string{"A new bridge is built.", "A new bridge is being built.", "A new bridge was building.", "A new bridge has building."},
		CorrectIndex: 1,
		Points:       5,
	},
	{
		Question:     "Выбери самый естественный вариант:",
		Options:      []string{"make a photo", "do a decision", "take a photo", "take a homework"},
		CorrectIndex: 2,
		Points:       4,
	},
	{
		Question:     "Заполни пропуск: Had I known about the delay, I _____ earlier.",
		Options:      []string{"would leave", "would have left", "will leave", "left"},
		CorrectIndex: 1,
		Points:       6,
	},
	{
		Question:     "Выбери ближайшее значение выражения: to put up with something",
		Options:      []string{"to tolerate it", "to build it", "to postpone it", "to improve it"},
		CorrectIndex: 0,
		Points:       5,
	},
	{
		Question:     "Заполни пропуск: I _____ coffee every morning.",
		Options:      []string{"drink", "drinks", "am drink", "drinking"},
		CorrectIndex: 0,
		Points:       1,
	},
	{
		Question:     "Выбери правильный вопрос:",
		Options:      []string{"Where you live?", "Where do you live?", "Where are you live?", "Where lives you?"},
		CorrectIndex: 1,
		Points:       2,
	},
	{
		Question:     "Выбери правильный перевод: Я уже сделал это.",
		Options:      []string{"I already did it.", "I have already done it.", "I already make it.", "I was already doing it."},
		CorrectIndex: 1,
		Points:       3,
	},
	{
		Question:     "Заполни пропуск: This book is _____ than the last one.",
		Options:      []string{"interesting", "more interesting", "most interesting", "interestinger"},
		CorrectIndex: 1,
		Points:       2,
	},
	{
		Question:     "Выбери лучший вариант: By the time we arrived, the film _____.",
		Options:      []string{"started", "has started", "had started", "was start"},
		CorrectIndex: 2,
		Points:       5,
	},
	{
		Question:     "Выбери самый естественный вариант:",
		Options:      []string{"strong rain", "heavy rain", "big rain", "hard rain"},
		CorrectIndex: 1,
		Points:       4,
	},
	{
		Question:     "Заполни пропуск: Not only _____ the report, but she also presented it brilliantly.",
		Options:      []string{"she wrote", "did she write", "she has written", "was she writing"},
		CorrectIndex: 1,
		Points:       6,
	},
	{
		Question:     "Выбери ближайшее значение слова: nuanced",
		Options:      []string{"simple and obvious", "full of small differences", "angry and direct", "impossible to prove"},
		CorrectIndex: 1,
		Points:       5,
	},
}

var levelAssessmentQuestionsByLanguage = map[string][]levelQuestion{
	"en": levelAssessmentQuestions,
	"de": {
		{Question: "Выбери правильный артикль: ___ Haus ist groß.", Options: []string{"Der", "Die", "Das", "Den"}, CorrectIndex: 2, Points: 8},
		{Question: "Как сказать по-немецки: Я живу в Берлине.", Options: []string{"Ich wohne in Berlin.", "Ich wohnen in Berlin.", "Ich lebe zu Berlin.", "Ich bin wohnen Berlin."}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери правильный Perfekt: Ich ___ gestern Deutsch gelernt.", Options: []string{"habe", "bin", "hat", "war"}, CorrectIndex: 0, Points: 9},
		{Question: "Какой порядок слов правильный?", Options: []string{"Morgen ich gehe ins Kino.", "Morgen gehe ich ins Kino.", "Ich gehe morgen ins Kino.", "Gehe ich morgen ins Kino."}, CorrectIndex: 1, Points: 9},
		{Question: "Выбери правильный вариант с дательным падежом:", Options: []string{"mit der Freund", "mit dem Freund", "mit den Freund", "mit die Freund"}, CorrectIndex: 1, Points: 10},
		{Question: "Заполни пропуск: Obwohl es regnet, ___ wir spazieren.", Options: []string{"gehen", "wir gehen", "gehen wir", "gingen wir"}, CorrectIndex: 0, Points: 10},
		{Question: "Выбери естественный вариант:", Options: []string{"Ich freue mich auf dich zu sehen.", "Ich freue mich darauf, dich zu sehen.", "Ich freue mich dich sehen.", "Ich freue auf dich sehen."}, CorrectIndex: 1, Points: 12},
		{Question: "Выбери наиболее точное значение слова dennoch:", Options: []string{"deshalb", "jedoch / trotzdem", "weil", "sofort"}, CorrectIndex: 1, Points: 14},
		{Question: "Выбери правильное место глагола: Heute ___ ich zu Hause.", Options: []string{"bleibe", "ich bleibe", "bleiben", "bin bleibe"}, CorrectIndex: 0, Points: 4},
		{Question: "Как сказать: У меня есть время.", Options: []string{"Ich habe Zeit.", "Ich bin Zeit.", "Ich mache Zeit.", "Ich werde Zeit."}, CorrectIndex: 0, Points: 4},
		{Question: "Выбери plural: ein Kind - zwei ___", Options: []string{"Kinder", "Kinds", "Kinde", "Kindern"}, CorrectIndex: 0, Points: 4},
		{Question: "Заполни пропуск: Ich interessiere mich ___ Musik.", Options: []string{"für", "auf", "an", "mit"}, CorrectIndex: 0, Points: 5},
		{Question: "Выбери правильный вариант: Я должен работать.", Options: []string{"Ich muss arbeiten.", "Ich muss arbeite.", "Ich soll arbeiten müssen.", "Ich habe arbeiten."}, CorrectIndex: 0, Points: 5},
		{Question: "Что значит trotzdem?", Options: []string{"dennoch", "deshalb", "niemals", "gestern"}, CorrectIndex: 0, Points: 6},
		{Question: "Выбери Konjunktiv II:", Options: []string{"Ich hätte gern einen Kaffee.", "Ich habe gern einen Kaffee.", "Ich würde gern einen Kaffee gehabt.", "Ich hatte gern Kaffee."}, CorrectIndex: 0, Points: 6},
		{Question: "Заполни пропуск: Das ist der Mann, ___ ich gesehen habe.", Options: []string{"den", "dem", "der", "das"}, CorrectIndex: 0, Points: 7},
		{Question: "Выбери правильный Passiv:", Options: []string{"Der Brief wird geschrieben.", "Der Brief schreibt.", "Der Brief ist schreiben.", "Der Brief wurde schreiben."}, CorrectIndex: 0, Points: 7},
		{Question: "Что значит etwas in Kauf nehmen?", Options: []string{"etwas akzeptieren", "etwas kaufen", "etwas verschieben", "etwas vergessen"}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери правильный вариант с инфинитивом:", Options: []string{"Es ist wichtig, pünktlich zu sein.", "Es ist wichtig, pünktlich sein.", "Es ist wichtig zu pünktlich sein.", "Es ist wichtig, dass pünktlich sein."}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери самое естественное:", Options: []string{"Das hängt davon ab.", "Das hängt von das ab.", "Das hängt ab davon.", "Das hängt an davon."}, CorrectIndex: 0, Points: 10},
	},
	"es": {
		{Question: "Выбери правильный артикль: ___ casa es grande.", Options: []string{"El", "La", "Los", "Un"}, CorrectIndex: 1, Points: 8},
		{Question: "Как сказать по-испански: Я говорю по-испански.", Options: []string{"Yo hablo español.", "Yo habla español.", "Yo hablar español.", "Yo soy hablo español."}, CorrectIndex: 0, Points: 8},
		{Question: "Заполни пропуск: Ayer ___ al cine.", Options: []string{"voy", "fui", "iba", "iré"}, CorrectIndex: 1, Points: 9},
		{Question: "Выбери правильный вариант с gustar:", Options: []string{"Yo gusto el café.", "Me gusta el café.", "Me gusto café.", "A mí gusto el café."}, CorrectIndex: 1, Points: 9},
		{Question: "Выбери правильный subjuntivo:", Options: []string{"Espero que vienes.", "Espero que vengas.", "Espero que venir.", "Espero que vendrás."}, CorrectIndex: 1, Points: 10},
		{Question: "Выбери ser или estar: La puerta ___ abierta.", Options: []string{"es", "está", "son", "están"}, CorrectIndex: 1, Points: 10},
		{Question: "Выбери естественный вариант:", Options: []string{"Si tendría tiempo, viajaría.", "Si tuviera tiempo, viajaría.", "Si tengo tiempo, viajaría.", "Si tuve tiempo, viajaría."}, CorrectIndex: 1, Points: 12},
		{Question: "Что значит sin embargo?", Options: []string{"no obstante / aun así", "porque", "en lugar de", "de inmediato"}, CorrectIndex: 0, Points: 14},
		{Question: "Выбери правильное место прилагательного: una casa ___", Options: []string{"grande", "grando", "granmente", "grandes"}, CorrectIndex: 0, Points: 4},
		{Question: "Как сказать: У меня есть время.", Options: []string{"Tengo tiempo.", "Soy tiempo.", "Estoy tiempo.", "Hago tiempo."}, CorrectIndex: 0, Points: 4},
		{Question: "Выбери plural: el niño - los ___", Options: []string{"niños", "niñes", "niño", "niñas"}, CorrectIndex: 0, Points: 4},
		{Question: "Заполни пропуск: Voy ___ casa.", Options: []string{"a", "en", "de", "por"}, CorrectIndex: 0, Points: 5},
		{Question: "Выбери правильный вариант: Мне нужно работать.", Options: []string{"Tengo que trabajar.", "Tengo trabajar.", "Debo de trabajo.", "Necesito que trabajo."}, CorrectIndex: 0, Points: 5},
		{Question: "Что значит aunque?", Options: []string{"si bien / a pesar de que", "porque", "nunca", "ayer"}, CorrectIndex: 0, Points: 6},
		{Question: "Выбери pretérito perfecto:", Options: []string{"He comido.", "Comí he.", "Soy comido.", "Había comer."}, CorrectIndex: 0, Points: 6},
		{Question: "Заполни пропуск: La persona ___ vi ayer.", Options: []string{"que", "quien", "cuyo", "donde"}, CorrectIndex: 0, Points: 7},
		{Question: "Выбери правильный pasivo/reflexivo:", Options: []string{"Se venden casas.", "Se vende casas.", "Son venden casas.", "Casas venden se."}, CorrectIndex: 0, Points: 7},
		{Question: "Что значит darse cuenta de algo?", Options: []string{"comprender / percatarse", "dar algo", "contar algo", "olvidar algo"}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери правильный infinitivo:", Options: []string{"Es importante estudiar.", "Es importante estudia.", "Es importante que estudiar.", "Es importante estudiando."}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери самое естественное:", Options: []string{"Depende de la situación.", "Depende la situación.", "Depende en la situación.", "Depende con la situación."}, CorrectIndex: 0, Points: 10},
	},
	"fr": {
		{Question: "Выбери правильный артикль: ___ maison est grande.", Options: []string{"Le", "La", "Les", "Un"}, CorrectIndex: 1, Points: 8},
		{Question: "Как сказать по-французски: Я говорю по-французски.", Options: []string{"Je parle français.", "Je parles français.", "Je parler français.", "Je suis parle français."}, CorrectIndex: 0, Points: 8},
		{Question: "Заполни пропуск: Hier, je ___ allé au cinéma.", Options: []string{"suis", "ai", "vais", "étais"}, CorrectIndex: 0, Points: 9},
		{Question: "Выбери правильный вариант отрицания:", Options: []string{"Je ne sais pas.", "Je sais ne pas.", "Je pas sais.", "Je ne pas sais."}, CorrectIndex: 0, Points: 9},
		{Question: "Выбери правильный pronom:", Options: []string{"Je lui vois.", "Je le vois.", "Je y vois.", "Je en vois."}, CorrectIndex: 1, Points: 10},
		{Question: "Выбери subjonctif:", Options: []string{"Il faut que tu viens.", "Il faut que tu viennes.", "Il faut que tu venir.", "Il faut que tu viendras."}, CorrectIndex: 1, Points: 10},
		{Question: "Выбери естественный вариант:", Options: []string{"Si j'aurais le temps, je voyagerais.", "Si j'avais le temps, je voyagerais.", "Si j'ai le temps, je voyagerais.", "Si j'avais le temps, je voyage."}, CorrectIndex: 1, Points: 12},
		{Question: "Что значит néanmoins?", Options: []string{"pourtant / malgré cela", "parce que", "maintenant", "ensemble"}, CorrectIndex: 0, Points: 14},
		{Question: "Выбери правильное прилагательное: une maison ___", Options: []string{"grande", "grand", "grands", "grandement"}, CorrectIndex: 0, Points: 4},
		{Question: "Как сказать: У меня есть время.", Options: []string{"J'ai le temps.", "Je suis le temps.", "Je fais le temps.", "Je vais temps."}, CorrectIndex: 0, Points: 4},
		{Question: "Выбери plural: le livre - les ___", Options: []string{"livres", "livre", "livreses", "livraux"}, CorrectIndex: 0, Points: 4},
		{Question: "Заполни пропуск: Je vais ___ Paris.", Options: []string{"à", "en", "de", "pour"}, CorrectIndex: 0, Points: 5},
		{Question: "Выбери правильный вариант: Мне нужно работать.", Options: []string{"Je dois travailler.", "Je dois travaille.", "J'ai travailler.", "Je suis travailler."}, CorrectIndex: 0, Points: 5},
		{Question: "Что значит pourtant?", Options: []string{"néanmoins / malgré cela", "parce que", "jamais", "demain"}, CorrectIndex: 0, Points: 6},
		{Question: "Выбери passé composé:", Options: []string{"J'ai mangé.", "Je suis mangé.", "J'ai manger.", "Je mangeais ai."}, CorrectIndex: 0, Points: 6},
		{Question: "Заполни пропуск: La personne ___ j'ai vue hier.", Options: []string{"que", "qui", "dont", "où"}, CorrectIndex: 0, Points: 7},
		{Question: "Выбери правильный passive:", Options: []string{"La lettre est écrite.", "La lettre écrit.", "La lettre est écrire.", "La lettre a écrire."}, CorrectIndex: 0, Points: 7},
		{Question: "Что значит se rendre compte de quelque chose?", Options: []string{"comprendre / prendre conscience", "se rendre", "rendre quelque chose", "aller quelque part"}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери правильный infinitif:", Options: []string{"Il est important d'étudier.", "Il est important étudier.", "Il est important que étudier.", "Il est important étudié."}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери самое естественное:", Options: []string{"Cela dépend de la situation.", "Cela dépend la situation.", "Cela dépend à la situation.", "Cela dépend avec la situation."}, CorrectIndex: 0, Points: 10},
	},
}

var italianLevelAssessmentQuestions = []levelQuestion{
	{Question: "Выбери правильный артикль: ___ casa è grande.", Options: []string{"Il", "La", "Lo", "Le"}, CorrectIndex: 1, Points: 4},
	{Question: "Как сказать по-итальянски: Я говорю по-итальянски.", Options: []string{"Io parlo italiano.", "Io parla italiano.", "Io parlare italiano.", "Io sono parlo italiano."}, CorrectIndex: 0, Points: 4},
	{Question: "Заполни пропуск: Oggi ___ a casa.", Options: []string{"sono", "sei", "è", "siamo"}, CorrectIndex: 0, Points: 4},
	{Question: "Выбери правильное множественное число: il libro - i ___", Options: []string{"libri", "libros", "libre", "libra"}, CorrectIndex: 0, Points: 4},
	{Question: "Как сказать: У меня есть время.", Options: []string{"Ho tempo.", "Sono tempo.", "Faccio tempo.", "Sto tempo."}, CorrectIndex: 0, Points: 5},
	{Question: "Заполни пропуск: Vado ___ Roma.", Options: []string{"a", "in", "di", "per"}, CorrectIndex: 0, Points: 5},
	{Question: "Выбери правильный вариант: Мне нужно работать.", Options: []string{"Devo lavorare.", "Devo lavoro.", "Ho lavorare.", "Sono lavorare."}, CorrectIndex: 0, Points: 5},
	{Question: "Что значит però?", Options: []string{"tuttavia", "perché", "mai", "ieri"}, CorrectIndex: 0, Points: 6},
	{Question: "Выбери passato prossimo:", Options: []string{"Ho mangiato.", "Sono mangiare.", "Ho mangiare.", "Mangiato ho."}, CorrectIndex: 0, Points: 6},
	{Question: "Заполни пропуск: La persona ___ ho visto ieri.", Options: []string{"che", "cui", "dove", "quale"}, CorrectIndex: 0, Points: 7},
	{Question: "Выбери правильный вариант с piacere:", Options: []string{"Mi piace il caffè.", "Io piaccio il caffè.", "Mi piaccio caffè.", "A me piace il caffè io."}, CorrectIndex: 0, Points: 8},
	{Question: "Заполни пропуск: Ieri ___ al cinema.", Options: []string{"sono andato", "ho andato", "vado", "andavo"}, CorrectIndex: 0, Points: 8},
	{Question: "Выбери правильное согласование: una macchina ___", Options: []string{"nuova", "nuovo", "nuovi", "nuove"}, CorrectIndex: 0, Points: 8},
	{Question: "Что значит anche se?", Options: []string{"sebbene / pur se", "dopo che", "affinché", "subito dopo"}, CorrectIndex: 0, Points: 8},
	{Question: "Выбери правильный порядок местоимений:", Options: []string{"Glielo do domani.", "Lo gli do domani.", "Gli lo do domani.", "Do glielo domani."}, CorrectIndex: 0, Points: 9},
	{Question: "Заполни пропуск: Penso che lui ___ ragione.", Options: []string{"abbia", "ha", "avrà", "avere"}, CorrectIndex: 0, Points: 9},
	{Question: "Выбери естественный вариант:", Options: []string{"Quanto più studi, tanto meglio parli.", "Più studi, più meglio parli.", "Quanto studi, tanto parli meglio futuro.", "Più studiare, più parlare."}, CorrectIndex: 0, Points: 9},
	{Question: "Что значит rendersi conto di qualcosa?", Options: []string{"accorgersi / capire", "consegnare qualcosa", "restituire qualcosa", "andare da qualche parte"}, CorrectIndex: 0, Points: 10},
	{Question: "Выбери условное предложение:", Options: []string{"Se avessi tempo, viaggerei.", "Se avrei tempo, viaggerei.", "Se ho tempo, viaggerei.", "Se avevo tempo, viaggio."}, CorrectIndex: 0, Points: 10},
	{Question: "Заполни пропуск: Non credo che lui ___ qui.", Options: []string{"sia", "è", "sarà", "essere"}, CorrectIndex: 0, Points: 10},
	{Question: "Выбери правильный condizionale composto:", Options: []string{"Sarei venuto se avessi potuto.", "Venirei se avrei potuto.", "Sarei venire se potevo.", "Sono venuto se potrei."}, CorrectIndex: 0, Points: 10},
	{Question: "Заполни пропуск: È la persona di ___ ti ho parlato.", Options: []string{"cui", "che", "dove", "quale"}, CorrectIndex: 0, Points: 10},
	{Question: "Выбери infinitivo passato:", Options: []string{"Dopo aver finito, siamo usciti.", "Dopo finito avere, siamo usciti.", "Dopo avere finire, siamo usciti.", "Dopo di finire, siamo usciti."}, CorrectIndex: 0, Points: 10},
	{Question: "Что значит mettere in discussione?", Options: []string{"mettere in dubbio", "iniziare una discussione a tavola", "approvare", "spiegare semplicemente"}, CorrectIndex: 0, Points: 11},
	{Question: "Выбери правильную конструкцию:", Options: []string{"Se l'avessi saputo prima, avrei agito diversamente.", "Se lo avrei saputo prima, avrei agito diversamente.", "Averlo sapere prima, agirei diverso.", "Di avrei saputo prima, agivo diverso."}, CorrectIndex: 0, Points: 11},
	{Question: "Выбери наиболее точный перевод: Вопрос остаётся нерешённым.", Options: []string{"La questione resta irrisolta.", "La domanda resta non risolvere.", "Il tema continua senza risolverlo.", "La questione è restare risolta."}, CorrectIndex: 0, Points: 11},
	{Question: "Что значит nonostante?", Options: []string{"malgrado / sebbene", "a causa di", "quindi", "invece di"}, CorrectIndex: 0, Points: 11},
	{Question: "Выбери формальный и точный вариант:", Options: []string{"La misura risulta parzialmente efficace.", "La misura è un po' buona.", "La misura fa efficacia parziale.", "La misura funziona qualcosa."}, CorrectIndex: 0, Points: 12},
	{Question: "Заполни пропуск: Se lo ___ saputo, non sarei venuto.", Options: []string{"avessi", "avrei", "ho", "abbia"}, CorrectIndex: 0, Points: 12},
	{Question: "Выбери корректную сложную конструкцию:", Options: []string{"Avendolo saputo troppo tardi, non ha potuto reagire.", "Averlo sapere troppo tardi, non ha potuto reagire.", "Avendo saperlo troppo tardi, non può reagito.", "Saputo avendo troppo tardi, non ha potere reagire."}, CorrectIndex: 0, Points: 12},
	{Question: "Что значит per quanto riguarda?", Options: []string{"riguardo a", "per quanto guarda", "dopo questo", "al contrario"}, CorrectIndex: 0, Points: 12},
	{Question: "Выбери правильный discorso indiretto:", Options: []string{"Ha detto che sarebbe venuto il giorno dopo.", "Ha detto che verrà domani ieri.", "Ha detto che viene il giorno dopo.", "Ha detto che sarebbe venuto domani passato."}, CorrectIndex: 0, Points: 13},
	{Question: "Выбери правильный congiuntivo trapassato:", Options: []string{"Pensavo che fosse già partito.", "Pensavo che era già partito.", "Pensavo che sarebbe già partito ieri.", "Pensavo che essere già partito."}, CorrectIndex: 0, Points: 13},
	{Question: "Что значит a prescindere da?", Options: []string{"indipendentemente da", "a causa di", "fuori", "subito dopo"}, CorrectIndex: 0, Points: 13},
	{Question: "Выбери наиболее естественный вариант:", Options: []string{"La proposta, pur essendo complessa, è realizzabile.", "La proposta, anche è complessa, può fare.", "La proposta è complessa ma realizzabile sicuro tutto.", "La proposta fa realizzazione nonostante."}, CorrectIndex: 0, Points: 13},
	{Question: "Что значит mettere in evidenza?", Options: []string{"sottolineare / evidenziare", "nascondere", "discutere", "trasferire denaro"}, CorrectIndex: 0, Points: 14},
}

func init() {
	levelAssessmentQuestions = append(levelAssessmentQuestions, []levelQuestion{
		{Question: "Choose the correct sentence:", Options: []string{"I have never been to Canada.", "I never was to Canada.", "I have never gone in Canada.", "I never have been Canada."}, CorrectIndex: 0, Points: 4},
		{Question: "Fill the gap: I wish I _____ more time to practise.", Options: []string{"have", "had", "will have", "am having"}, CorrectIndex: 1, Points: 5},
		{Question: "Choose the most natural option:", Options: []string{"Despite it was raining, we went out.", "Although it was raining, we went out.", "Because it was raining, but we went out.", "In spite it was raining, we went out."}, CorrectIndex: 1, Points: 5},
		{Question: "What does 'to bring up a topic' mean?", Options: []string{"to mention it", "to lift it", "to forget it", "to prove it"}, CorrectIndex: 0, Points: 6},
		{Question: "Choose the correct reported speech:", Options: []string{"She said she was tired.", "She said she is tired yesterday.", "She said that she tired.", "She said she has tired."}, CorrectIndex: 0, Points: 6},
		{Question: "Fill the gap: The proposal, _____ was rejected at first, later became popular.", Options: []string{"which", "what", "who", "whose"}, CorrectIndex: 0, Points: 7},
		{Question: "Choose the best collocation:", Options: []string{"raise awareness", "lift awareness", "grow awareness", "up awareness"}, CorrectIndex: 0, Points: 7},
		{Question: "Choose the correct inversion:", Options: []string{"Rarely have I seen such a clear explanation.", "Rarely I have seen such a clear explanation.", "Rarely have seen I such a clear explanation.", "Rarely I seen have such a clear explanation."}, CorrectIndex: 0, Points: 8},
		{Question: "What does 'it is unlikely to be feasible' mean?", Options: []string{"it probably cannot be done", "it is easy to do", "it must be done now", "it is already finished"}, CorrectIndex: 0, Points: 8},
		{Question: "Choose the most precise rewrite: 'The data suggests a possible link, but it is not conclusive.'", Options: []string{"The evidence points to a link, though it does not prove one.", "The data proves the link completely.", "The link is impossible.", "The data is unrelated and useless."}, CorrectIndex: 0, Points: 9},
		{Question: "Choose the correct mixed conditional:", Options: []string{"If I had studied harder, I would know the answer now.", "If I studied harder, I would have known now.", "If I had study harder, I know now.", "If I would study harder, I had known."}, CorrectIndex: 0, Points: 9},
		{Question: "Choose the most natural academic phrase:", Options: []string{"This issue warrants further investigation.", "This issue wants more looking.", "This issue needs to be watched about.", "This issue asks more research."}, CorrectIndex: 0, Points: 10},
		{Question: "What does 'notwithstanding' mean in formal English?", Options: []string{"despite", "because of", "therefore", "instead of"}, CorrectIndex: 0, Points: 10},
		{Question: "Choose the grammatically correct sentence:", Options: []string{"Had the results been available earlier, the decision might have changed.", "Had the results available earlier, the decision might changed.", "If had the results been available earlier, the decision might change.", "Were the results available earlier yesterday, the decision had changed."}, CorrectIndex: 0, Points: 11},
		{Question: "Choose the sentence with the clearest nuance:", Options: []string{"The policy is well intentioned, albeit difficult to implement.", "The policy is good but hard and therefore surely wrong.", "The policy means good and has difficulty.", "The policy is nice although implementation difficultly."}, CorrectIndex: 0, Points: 11},
		{Question: "What does 'to call into question' mean?", Options: []string{"to make something seem doubtful", "to ask a phone question", "to answer confidently", "to approve officially"}, CorrectIndex: 0, Points: 12},
	}...)
	levelAssessmentQuestionsByLanguage["en"] = levelAssessmentQuestions
	levelAssessmentQuestionsByLanguage["de"] = append(levelAssessmentQuestionsByLanguage["de"], []levelQuestion{
		{Question: "Выбери правильный вариант: Я бы пришёл, если бы мог.", Options: []string{"Ich würde kommen, wenn ich könnte.", "Ich werde kommen, wenn ich konnte.", "Ich käme, wenn ich kann.", "Ich würde gekommen, wenn ich könnte."}, CorrectIndex: 0, Points: 8},
		{Question: "Заполни пропуск: Das ist das Buch, ___ ich dir erzählt habe.", Options: []string{"von dem", "das", "deren", "wo"}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери правильный вариант с zu + Infinitiv:", Options: []string{"Er hat vergessen, die Tür zu schließen.", "Er hat vergessen, die Tür schließen.", "Er hat vergessen zu die Tür schließen.", "Er hat vergessen, dass die Tür zu schließen."}, CorrectIndex: 0, Points: 8},
		{Question: "Что значит 'sich mit etwas auseinandersetzen'?", Options: []string{"sich intensiv mit etwas befassen", "sich neben etwas setzen", "etwas vergessen", "etwas beschleunigen"}, CorrectIndex: 0, Points: 9},
		{Question: "Выбери правильный Nominalstil:", Options: []string{"Nach Abschluss der Analyse wurden die Daten veröffentlicht.", "Nach die Analyse abschließen wurden Daten veröffentlicht.", "Nach abgeschlossen Analyse die Daten veröffentlicht.", "Nach Analyse schließen Daten wurden veröffentlicht."}, CorrectIndex: 0, Points: 9},
		{Question: "Заполни пропуск: Je früher wir anfangen, ___ besser wird das Ergebnis.", Options: []string{"desto", "weil", "obwohl", "damit"}, CorrectIndex: 0, Points: 9},
		{Question: "Выбери естественный вариант:", Options: []string{"Es kommt darauf an, wie viel Zeit wir haben.", "Es kommt an darauf, wie viel Zeit haben wir.", "Es kommt darauf, wie viel Zeit wir haben an.", "Es kommt daran, wie viel Zeit wir haben."}, CorrectIndex: 0, Points: 10},
		{Question: "Что значит 'in Frage stellen'?", Options: []string{"bezweifeln", "eine einfache Frage stellen", "sich anstellen", "um Erlaubnis bitten"}, CorrectIndex: 0, Points: 10},
		{Question: "Выбери правильный вариант с Partizip I:", Options: []string{"die zunehmende Bedeutung", "die zugenommene Bedeutung", "die Bedeutung zunehmend ist", "die Bedeutung zunehmen"}, CorrectIndex: 0, Points: 10},
		{Question: "Заполни пропуск: Ohne deine Hilfe hätte ich das nicht ___", Options: []string{"geschafft", "schaffen", "geschaffen werden", "schaffte"}, CorrectIndex: 0, Points: 11},
		{Question: "Выбери наиболее точный перевод: 'Несмотря на возражения, предложение приняли.'", Options: []string{"Trotz der Einwände wurde der Vorschlag angenommen.", "Trotz die Einwände hat man den Vorschlag nehmen.", "Obwohl Einwände, der Vorschlag wurde genommen.", "Wegen der Einwände wurde der Vorschlag angenommen."}, CorrectIndex: 0, Points: 11},
		{Question: "Что значит 'etwas lässt zu wünschen übrig'?", Options: []string{"nicht zufriedenstellend sein", "ideal sein", "verpflichtend sein", "verboten sein"}, CorrectIndex: 0, Points: 12},
		{Question: "Выбери корректную сложную конструкцию:", Options: []string{"Wäre der Bericht früher eingereicht worden, hätten wir anders entschieden.", "Würde der Bericht früher eingereicht, hätten wir anders entscheiden.", "Wäre der Bericht früher einreichen, wir hätten anders entschieden.", "Wenn wäre der Bericht früher eingereicht worden, hätten wir anders entschieden."}, CorrectIndex: 0, Points: 12},
		{Question: "Выбери самый формальный и точный вариант:", Options: []string{"Die Maßnahme erweist sich als nur bedingt wirksam.", "Die Maßnahme ist ein bisschen gut.", "Die Maßnahme macht Wirkung begrenzt.", "Die Maßnahme tut nicht ganz wirken."}, CorrectIndex: 0, Points: 12},
		{Question: "Что значит 'im Hinblick auf'?", Options: []string{"in Bezug auf", "wegen eines Blicks", "plötzlich", "draußen"}, CorrectIndex: 0, Points: 13},
		{Question: "Выбери правильную конструкцию с Konjunktiv I:", Options: []string{"Der Minister erklärte, die Lage sei stabil.", "Der Minister erklärte, die Lage ist stabil gestern.", "Der Minister erklärte, die Lage wäre stabil gewesen jetzt.", "Der Minister erklärte, dass die Lage stabil sein."}, CorrectIndex: 0, Points: 13},
	}...)
	levelAssessmentQuestionsByLanguage["es"] = append(levelAssessmentQuestionsByLanguage["es"], []levelQuestion{
		{Question: "Выбери правильный condicional compuesto:", Options: []string{"Habría venido si hubiera podido.", "Vendría si habría podido.", "Hubiera venido si podría.", "Habría venir si pudiera."}, CorrectIndex: 0, Points: 8},
		{Question: "Заполни пропуск: Es la persona de ___ te hablé.", Options: []string{"quien", "que", "donde", "cuyo"}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери правильный infinitivo compuesto:", Options: []string{"Después de haber terminado, salimos.", "Después de terminado haber, salimos.", "Después haber terminar, salimos.", "Después de terminar haber, salimos."}, CorrectIndex: 0, Points: 8},
		{Question: "Что значит 'llevar a cabo'?", Options: []string{"realizar / ejecutar", "llevar sobre la cabeza", "llegar tarde", "llamar por teléfono"}, CorrectIndex: 0, Points: 9},
		{Question: "Выбери естественный вариант:", Options: []string{"Cuanto más practiques, mejor hablarás.", "Más practiques, más mejor hablarás.", "Cuanto practicas, mejor hablas futuro.", "Mientras más practicar, más hablarás."}, CorrectIndex: 0, Points: 9},
		{Question: "Заполни пропуск: No creo que ___ razón.", Options: []string{"tenga", "tiene", "tendrá", "tener"}, CorrectIndex: 0, Points: 9},
		{Question: "Выбери правильное значение 'a pesar de que':", Options: []string{"aunque / pese a que", "porque", "para que", "justo después"}, CorrectIndex: 0, Points: 10},
		{Question: "Выбери правильный pronombre:", Options: []string{"Se lo dije ayer.", "Le lo dije ayer.", "Lo le dije ayer.", "Se le dije lo ayer."}, CorrectIndex: 0, Points: 10},
		{Question: "Выбери правильную perífrasis:", Options: []string{"Acabo de llegar.", "Acabo llegar.", "Acabo a llegar.", "Soy acabado de llegar."}, CorrectIndex: 0, Points: 10},
		{Question: "Заполни пропуск: Si lo ___ sabido, no habría venido.", Options: []string{"hubiera", "habría", "he", "haya"}, CorrectIndex: 0, Points: 11},
		{Question: "Выбери наиболее точный перевод: 'Вопрос остаётся нерешённым.'", Options: []string{"La cuestión sigue sin resolverse.", "La cuestión sigue no resolver.", "La pregunta queda sin solucionar se.", "El tema continúa sin resolverlo."}, CorrectIndex: 0, Points: 11},
		{Question: "Что значит 'poner en tela de juicio'?", Options: []string{"cuestionar", "poner tela ante un juez", "aprobar", "explicar de forma simple"}, CorrectIndex: 0, Points: 12},
		{Question: "Выбери корректную сложную конструкцию:", Options: []string{"De haberlo sabido antes, habría actuado de otra manera.", "De saberlo antes, habría actué de otra manera.", "Haberlo sabido antes, actuaría otra manera.", "De habría sabido antes, hubiera actuar."}, CorrectIndex: 0, Points: 12},
		{Question: "Выбери самый формальный и точный вариант:", Options: []string{"La medida resulta parcialmente eficaz.", "La medida es un poco buena.", "La medida hace eficacia parcial.", "La medida funciona algo."}, CorrectIndex: 0, Points: 12},
		{Question: "Что значит 'con respecto a'?", Options: []string{"en relación con", "con respeto personal", "después de eso", "al contrario"}, CorrectIndex: 0, Points: 13},
		{Question: "Выбери правильный estilo indirecto:", Options: []string{"Dijo que vendría al día siguiente.", "Dijo que vendrá mañana ayer.", "Dijo que viene al día siguiente.", "Dijo que vendría mañana pasado."}, CorrectIndex: 0, Points: 13},
	}...)
	levelAssessmentQuestionsByLanguage["fr"] = append(levelAssessmentQuestionsByLanguage["fr"], []levelQuestion{
		{Question: "Выбери plus-que-parfait + conditionnel:", Options: []string{"Je serais venu si j'avais pu.", "Je viendrais si j'aurais pu.", "Je serais venir si je pouvais.", "Je venais si j'avais pu."}, CorrectIndex: 0, Points: 8},
		{Question: "Заполни пропуск: C'est la personne ___ je t'ai parlé.", Options: []string{"dont", "que", "où", "laquelle"}, CorrectIndex: 0, Points: 8},
		{Question: "Выбери правильный infinitif passé:", Options: []string{"Après avoir terminé, nous sommes sortis.", "Après terminé avoir, nous sommes sortis.", "Après avoir terminer, nous sommes sortis.", "Après de terminer, nous sommes sortis."}, CorrectIndex: 0, Points: 8},
		{Question: "Что значит 'mettre en œuvre'?", Options: []string{"mettre en place / réaliser", "mettre dans une œuvre", "oublier", "appeler"}, CorrectIndex: 0, Points: 9},
		{Question: "Выбери естественный вариант:", Options: []string{"Plus tu pratiques, mieux tu parleras.", "Plus tu pratiques, plus mieux tu parleras.", "Le plus tu pratiqueras, mieux tu parles.", "Plus pratiquer, mieux parler."}, CorrectIndex: 0, Points: 9},
		{Question: "Заполни пропуск: Je ne pense pas qu'il ___ raison.", Options: []string{"ait", "a", "aura", "avoir"}, CorrectIndex: 0, Points: 9},
		{Question: "Что значит 'bien que'?", Options: []string{"quoique / même si", "parce que", "après que", "seulement si"}, CorrectIndex: 0, Points: 10},
		{Question: "Выбери правильные pronoms:", Options: []string{"Je le lui ai donné.", "Je lui le ai donné.", "Je l'ai lui donné.", "Je lui ai donné le."}, CorrectIndex: 0, Points: 10},
		{Question: "Выбери правильную конструкцию:", Options: []string{"Je viens de partir.", "Je viens partir.", "Je suis venu de partir.", "Je viens à partir."}, CorrectIndex: 0, Points: 10},
		{Question: "Заполни пропуск: Si je l'___ su, je ne serais pas venu.", Options: []string{"avais", "aurais", "ai", "aie"}, CorrectIndex: 0, Points: 11},
		{Question: "Выбери наиболее точный перевод: 'Вопрос остаётся нерешённым.'", Options: []string{"La question reste à résoudre.", "La question reste résoudre.", "La question est rester résolue.", "La question continue sans résoudre."}, CorrectIndex: 0, Points: 11},
		{Question: "Что значит 'remettre en question'?", Options: []string{"mettre en doute", "répéter la question", "répondre tout de suite", "approuver"}, CorrectIndex: 0, Points: 12},
		{Question: "Выбери корректную сложную конструкцию:", Options: []string{"L'ayant appris trop tard, il n'a pas pu réagir.", "Ayant l'apprendre trop tard, il n'a pas pu réagir.", "L'avoir appris trop tard, il ne peut pas réagi.", "En ayant appris trop tard, il n'a pas pouvoir réagir."}, CorrectIndex: 0, Points: 12},
		{Question: "Выбери самый формальный и точный вариант:", Options: []string{"Cette mesure s'avère partiellement efficace.", "Cette mesure est un peu bonne.", "Cette mesure fait efficacité partielle.", "Cette mesure marche quelque."}, CorrectIndex: 0, Points: 12},
		{Question: "Что значит 'à l'égard de'?", Options: []string{"envers / à propos de", "près du jardin", "après cela", "malgré tout"}, CorrectIndex: 0, Points: 13},
		{Question: "Выбери правильный discours indirect:", Options: []string{"Il a dit qu'il viendrait le lendemain.", "Il a dit qu'il viendra demain hier.", "Il a dit qu'il vient le lendemain.", "Il a dit qu'il viendrait demain passé."}, CorrectIndex: 0, Points: 13},
	}...)
	levelAssessmentQuestionsByLanguage["it"] = italianLevelAssessmentQuestions
}

func levelAssessmentQuestionsForLanguage(language string) []levelQuestion {
	language = normalizeLearningLanguage(language)
	questions := levelAssessmentQuestionsByLanguage[language]
	if len(questions) > 0 {
		return questions
	}
	return nil
}

func levelAssessmentQuestionsForUser(user userState) []levelQuestion {
	language := normalizeLearningLanguage(user.LearningLanguage)
	questions := levelAssessmentQuestionsByLanguage[language]
	if len(questions) == 0 {
		return nil
	}
	localized := make([]levelQuestion, len(questions))
	for index, question := range questions {
		localized[index] = localizedLevelAssessmentQuestion(user, question)
	}
	return localized
}

func activeLevelAssessmentQuestionsForUser(user userState) []levelQuestion {
	questions := levelAssessmentQuestionsForUser(user)
	_, _, order, ok := parseLevelTestModeState(user.Mode)
	if !ok || len(order) == 0 {
		return questions
	}
	if !validLevelTestOrder(order, len(questions)) {
		return questions
	}
	ordered := make([]levelQuestion, 0, len(questions))
	for _, index := range order {
		ordered = append(ordered, questions[index])
	}
	return ordered
}

func levelAssessmentMaxScore(language string) int {
	total := 0
	for _, question := range levelAssessmentQuestionsForLanguage(language) {
		total += question.Points
	}
	return total
}

func levelAssessmentMaxScoreForUser(user userState) int {
	total := 0
	for _, question := range levelAssessmentQuestionsForUser(user) {
		total += question.Points
	}
	return total
}

func levelFromAssessmentScore(score int, maxScore int) string {
	if maxScore == 72 {
		return levelFromEnglishAssessmentScore(score)
	}
	if maxScore <= 0 {
		return "A1"
	}
	ratio := float64(score) / float64(maxScore)
	switch {
	case ratio >= 0.86:
		return "C2"
	case ratio >= 0.71:
		return "C1"
	case ratio >= 0.53:
		return "B2"
	case ratio >= 0.33:
		return "B1"
	case ratio >= 0.15:
		return "A2"
	default:
		return "A1"
	}
}

func levelFromEnglishAssessmentScore(score int) string {
	switch {
	case score >= 62:
		return "C2"
	case score >= 51:
		return "C1"
	case score >= 38:
		return "B2"
	case score >= 24:
		return "B1"
	case score >= 11:
		return "A2"
	default:
		return "A1"
	}
}

func levelAssessmentStartTextForUser(user userState) string {
	questions := levelAssessmentQuestionsForUser(user)
	learningLanguage := userLearningLanguage(user)
	return fmt.Sprintf(systemUI(user).LevelStartText, learningLanguage.InterfaceName, len(questions))
}

func manualLevelSelectionTextForUser(user userState) string {
	return systemUI(user).ManualLevelText
}

func levelQuestionTextForUser(user userState, index int, score int) string {
	questions := activeLevelAssessmentQuestionsForUser(user)
	question := questions[index]
	copy := systemUI(user)
	return fmt.Sprintf(copy.LevelQuestionText, index+1, len(questions), question.Question) + "\n\n" + fmt.Sprintf(copy.LevelCurrentScore, score)
}

func levelAssessmentKeyboardForUser(user userState, index int) map[string]any {
	question := activeLevelAssessmentQuestionsForUser(user)[index]
	return levelAssessmentKeyboardForQuestion(question, index, systemUI(user), ui(user).BackMenu)
}

func localizedLevelAssessmentQuestion(user userState, question levelQuestion) levelQuestion {
	question.Question = localizedAssessmentPrompt(question.Question, normalizeInterfaceLanguage(user.InterfaceLanguage), normalizeLearningLanguage(user.LearningLanguage))
	return question
}

func localizedAssessmentPrompt(question string, interfaceLanguage string, learningLanguage string) string {
	task := assessmentPromptTaskForQuestion(question)
	label := localizedAssessmentTaskLabel(task.Kind, interfaceLanguage)
	if task.Payload == "" || ((containsCyrillicText(task.Payload) || containsUnicodeCyrillicText(task.Payload)) && !isCyrillicLearningLanguage(learningLanguage) && normalizeInterfaceLanguage(interfaceLanguage) != "ru") {
		return label
	}
	return label + " " + task.Payload
}

func isCyrillicLearningLanguage(language string) bool {
	switch normalizeLearningLanguage(language) {
	case "ru", "uk", "kk", "ky", "tg", "tt":
		return true
	default:
		return false
	}
}

func assessmentPromptPayload(question string) string {
	question = strings.TrimSpace(question)
	if question == "" {
		return ""
	}

	for _, prefix := range []string{
		"Заполни пропуск:",
		"Fill the gap:",
		"Выбери правильный перевод:",
		"Выбери лучший вариант:",
		"Choose the most precise rewrite:",
		"Выбери подходящую связку:",
		"Выбери правильную форму пассива для предложения:",
		"Выбери ближайшее значение слова:",
		"Выбери ближайшее значение выражения:",
		"Выбери правильный артикль:",
		"Выбери правильный Perfekt:",
		"Выбери правильный вариант с дательным падежом:",
		"Выбери правильное место глагола:",
		"Выбери plural:",
		"Выбери правильный вариант:",
		"Выбери правильное место прилагательного:",
		"Выбери правильное прилагательное:",
		"Выбери правильное множественное число:",
		"Выбери правильное согласование:",
		"Выбери наиболее точный перевод:",
		"Как сказать по-немецки:",
		"Как сказать по-испански:",
		"Как сказать по-французски:",
		"Как сказать по-итальянски:",
		"Как сказать:",
	} {
		if strings.HasPrefix(question, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(question, prefix))
		}
	}

	if strings.HasPrefix(question, "Что значит ") {
		return strings.Trim(strings.TrimSpace(strings.TrimPrefix(question, "Что значит ")), "?")
	}
	if strings.HasPrefix(question, "What does ") {
		payload := strings.TrimSpace(strings.TrimPrefix(question, "What does "))
		payload = strings.TrimSuffix(payload, " mean in formal English?")
		payload = strings.TrimSuffix(payload, " mean?")
		return strings.TrimSpace(payload)
	}

	if index := strings.Index(question, ":"); index >= 0 && index+1 < len(question) {
		return strings.TrimSpace(question[index+1:])
	}
	return ""
}

type assessmentPromptTask struct {
	Kind    string
	Payload string
}

func assessmentPromptTaskForQuestion(question string) assessmentPromptTask {
	question = strings.TrimSpace(question)
	if question == "" {
		return assessmentPromptTask{Kind: "choose_answer"}
	}
	for _, item := range []struct {
		prefix string
		kind   string
	}{
		{"Выбери правильное предложение:", "correct_sentence"},
		{"Заполни пропуск:", "fill_gap"},
		{"Выбери правильный перевод:", "correct_translation"},
		{"Выбери наиболее точный перевод:", "precise_translation"},
		{"Выбери ближайшее значение слова:", "word_meaning"},
		{"Выбери ближайшее значение выражения:", "phrase_meaning"},
		{"Выбери лучший вариант:", "best_option"},
		{"Выбери самый естественный вариант:", "natural_option"},
		{"Выбери самое естественное:", "natural_option"},
		{"Выбери наиболее естественный вариант:", "natural_option"},
		{"Выбери подходящую связку:", "connector"},
		{"Выбери правильный артикль:", "correct_article"},
		{"Выбери правильный Perfekt:", "correct_perfekt"},
		{"Выбери правильную форму пассива для предложения:", "correct_passive"},
		{"Выбери правильный Passiv:", "correct_passive"},
		{"Выбери plural:", "correct_plural"},
		{"Выбери правильное множественное число:", "correct_plural"},
		{"Выбери правильный вопрос:", "correct_question"},
		{"Выбери правильное место глагола:", "verb_position"},
		{"Выбери правильное место прилагательного:", "adjective_position"},
		{"Выбери правильное прилагательное:", "correct_adjective"},
		{"Выбери правильное согласование:", "agreement"},
		{"Выбери правильный вариант с дательным падежом:", "dative_variant"},
		{"Выбери правильный вариант:", "correct_variant"},
		{"Как сказать по-немецки:", "how_to_say"},
		{"Как сказать по-испански:", "how_to_say"},
		{"Как сказать по-французски:", "how_to_say"},
		{"Как сказать по-итальянски:", "how_to_say"},
		{"Как сказать:", "how_to_say"},
		{"Р—Р°РїРѕР»РЅРё РїСЂРѕРїСѓСЃРє:", "fill_gap"},
		{"Fill the gap:", "fill_gap"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅС‹Р№ РїРµСЂРµРІРѕРґ:", "correct_translation"},
		{"Р’С‹Р±РµСЂРё РЅР°РёР±РѕР»РµРµ С‚РѕС‡РЅС‹Р№ РїРµСЂРµРІРѕРґ:", "precise_translation"},
		{"Р’С‹Р±РµСЂРё Р±Р»РёР¶Р°Р№С€РµРµ Р·РЅР°С‡РµРЅРёРµ СЃР»РѕРІР°:", "word_meaning"},
		{"Р’С‹Р±РµСЂРё Р±Р»РёР¶Р°Р№С€РµРµ Р·РЅР°С‡РµРЅРёРµ РІС‹СЂР°Р¶РµРЅРёСЏ:", "phrase_meaning"},
		{"Р’С‹Р±РµСЂРё Р»СѓС‡С€РёР№ РІР°СЂРёР°РЅС‚:", "best_option"},
		{"Р’С‹Р±РµСЂРё СЃР°РјС‹Р№ РµСЃС‚РµСЃС‚РІРµРЅРЅС‹Р№ РІР°СЂРёР°РЅС‚:", "natural_option"},
		{"Р’С‹Р±РµСЂРё СЃР°РјРѕРµ РµСЃС‚РµСЃС‚РІРµРЅРЅРѕРµ:", "natural_option"},
		{"Р’С‹Р±РµСЂРё РЅР°РёР±РѕР»РµРµ РµСЃС‚РµСЃС‚РІРµРЅРЅС‹Р№ РІР°СЂРёР°РЅС‚:", "natural_option"},
		{"Р’С‹Р±РµСЂРё РїРѕРґС…РѕРґСЏС‰СѓСЋ СЃРІСЏР·РєСѓ:", "connector"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅС‹Р№ Р°СЂС‚РёРєР»СЊ:", "correct_article"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅС‹Р№ Perfekt:", "correct_perfekt"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅСѓСЋ С„РѕСЂРјСѓ РїР°СЃСЃРёРІР° РґР»СЏ РїСЂРµРґР»РѕР¶РµРЅРёСЏ:", "correct_passive"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅС‹Р№ Passiv:", "correct_passive"},
		{"Р’С‹Р±РµСЂРё plural:", "correct_plural"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅРѕРµ РјРЅРѕР¶РµСЃС‚РІРµРЅРЅРѕРµ С‡РёСЃР»Рѕ:", "correct_plural"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅС‹Р№ РІРѕРїСЂРѕСЃ:", "correct_question"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅРѕРµ РјРµСЃС‚Рѕ РіР»Р°РіРѕР»Р°:", "verb_position"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅРѕРµ РјРµСЃС‚Рѕ РїСЂРёР»Р°РіР°С‚РµР»СЊРЅРѕРіРѕ:", "adjective_position"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅРѕРµ РїСЂРёР»Р°РіР°С‚РµР»СЊРЅРѕРµ:", "correct_adjective"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅРѕРµ СЃРѕРіР»Р°СЃРѕРІР°РЅРёРµ:", "agreement"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅС‹Р№ РІР°СЂРёР°РЅС‚ СЃ РґР°С‚РµР»СЊРЅС‹Рј РїР°РґРµР¶РѕРј:", "dative_variant"},
		{"Choose the most precise rewrite:", "precise_rewrite"},
		{"Р’С‹Р±РµСЂРё РїСЂР°РІРёР»СЊРЅС‹Р№ РІР°СЂРёР°РЅС‚:", "correct_variant"},
		{"РљР°Рє СЃРєР°Р·Р°С‚СЊ РїРѕ-РЅРµРјРµС†РєРё:", "how_to_say"},
		{"РљР°Рє СЃРєР°Р·Р°С‚СЊ РїРѕ-РёСЃРїР°РЅСЃРєРё:", "how_to_say"},
		{"РљР°Рє СЃРєР°Р·Р°С‚СЊ РїРѕ-С„СЂР°РЅС†СѓР·СЃРєРё:", "how_to_say"},
		{"РљР°Рє СЃРєР°Р·Р°С‚СЊ РїРѕ-РёС‚Р°Р»СЊСЏРЅСЃРєРё:", "how_to_say"},
		{"РљР°Рє СЃРєР°Р·Р°С‚СЊ:", "how_to_say"},
	} {
		if strings.HasPrefix(question, item.prefix) {
			return assessmentPromptTask{Kind: item.kind, Payload: strings.TrimSpace(strings.TrimPrefix(question, item.prefix))}
		}
	}
	if strings.HasPrefix(question, "Р§С‚Рѕ Р·РЅР°С‡РёС‚ ") {
		return assessmentPromptTask{Kind: "what_means", Payload: strings.Trim(strings.TrimSpace(strings.TrimPrefix(question, "Р§С‚Рѕ Р·РЅР°С‡РёС‚ ")), "?")}
	}
	if strings.HasPrefix(question, "Что значит ") {
		return assessmentPromptTask{Kind: "what_means", Payload: strings.Trim(strings.TrimSpace(strings.TrimPrefix(question, "Что значит ")), "?")}
	}
	if strings.HasPrefix(question, "What does ") {
		payload := strings.TrimSpace(strings.TrimPrefix(question, "What does "))
		payload = strings.TrimSuffix(payload, " mean in formal English?")
		payload = strings.TrimSuffix(payload, " mean?")
		return assessmentPromptTask{Kind: "what_means", Payload: strings.TrimSpace(payload)}
	}
	if strings.HasSuffix(question, "?") {
		return assessmentPromptTask{Kind: "question", Payload: question}
	}
	if index := strings.Index(question, ":"); index >= 0 && index+1 < len(question) {
		return assessmentPromptTask{Kind: "choose_answer", Payload: strings.TrimSpace(question[index+1:])}
	}
	return assessmentPromptTask{Kind: "choose_answer"}
}

func localizedAssessmentTaskLabel(kind string, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		if label := russianAssessmentTaskLabel(kind); label != "" {
			return label
		}
	}
	labels := map[string]map[string]string{
		"ru": {
			"fill_gap": "Заполни пропуск:", "correct_translation": "Выбери правильный перевод:", "precise_translation": "Выбери наиболее точный перевод:",
			"correct_article": "Выбери правильный артикль:", "correct_perfekt": "Выбери правильный Perfekt:", "correct_passive": "Выбери правильный пассив:",
			"correct_plural": "Выбери правильное множественное число:", "correct_question": "Выбери правильный вопрос:", "connector": "Выбери подходящую связку:",
			"word_meaning": "Выбери ближайшее значение слова:", "phrase_meaning": "Выбери ближайшее значение выражения:", "precise_rewrite": "Выбери наиболее точную переформулировку:",
			"best_option": "Выбери лучший вариант:", "natural_option": "Выбери самый естественный вариант:", "verb_position": "Выбери правильное место глагола:",
			"adjective_position": "Выбери правильное место прилагательного:", "correct_adjective": "Выбери правильное прилагательное:", "agreement": "Выбери правильное согласование:",
			"dative_variant": "Выбери правильный вариант с дательным падежом:", "correct_variant": "Выбери правильную форму:", "how_to_say": "Как сказать:",
			"what_means": "Что значит", "question": "Ответь на вопрос:", "choose_answer": "Выбери правильный вариант:",
		},
		"en": {
			"fill_gap": "Fill the gap:", "correct_sentence": "Choose the correct sentence:", "correct_translation": "Choose the correct translation:", "precise_translation": "Choose the most precise translation:",
			"correct_article": "Choose the correct article:", "correct_perfekt": "Choose the correct Perfekt:", "correct_passive": "Choose the correct passive form:",
			"correct_plural": "Choose the correct plural:", "correct_question": "Choose the correct question:", "connector": "Choose the right connector:",
			"word_meaning": "Choose the closest meaning of the word:", "phrase_meaning": "Choose the closest meaning of the expression:", "precise_rewrite": "Choose the most precise rewrite:",
			"best_option": "Choose the best option:", "natural_option": "Choose the most natural option:", "verb_position": "Choose the correct verb position:",
			"adjective_position": "Choose the correct adjective position:", "correct_adjective": "Choose the correct adjective:", "agreement": "Choose the correct agreement:",
			"dative_variant": "Choose the correct dative variant:", "correct_variant": "Choose the correct variant:", "how_to_say": "How do you say:",
			"what_means": "What does", "question": "Answer the question:", "choose_answer": "Choose the correct answer:",
		},
	}
	if localized, ok := labels[normalizeInterfaceLanguage(interfaceLanguage)][kind]; ok {
		return localized
	}
	if localized := localizedAssessmentTaskLabelForAllLanguages(kind, interfaceLanguage); localized != "" {
		return localized
	}
	if prompt := strings.TrimSpace(ui(userState{InterfaceLanguage: interfaceLanguage}).ChooseAnswer); prompt != "" {
		return prompt
	}
	if localized, ok := labels["en"][kind]; ok {
		return localized
	}
	return englishUICopy().ChooseAnswer
}

func russianAssessmentTaskLabel(kind string) string {
	labels := map[string]string{
		"fill_gap":            "Заполни пропуск:",
		"correct_sentence":    "Выбери правильное предложение:",
		"correct_translation": "Выбери правильный перевод:",
		"precise_translation": "Выбери наиболее точный перевод:",
		"correct_article":     "Выбери правильный артикль:",
		"correct_perfekt":     "Выбери правильный Perfekt:",
		"correct_passive":     "Выбери правильный пассив:",
		"correct_plural":      "Выбери правильное множественное число:",
		"correct_question":    "Выбери правильный вопрос:",
		"connector":           "Выбери подходящую связку:",
		"word_meaning":        "Выбери ближайшее значение слова:",
		"phrase_meaning":      "Выбери ближайшее значение выражения:",
		"precise_rewrite":     "Выбери наиболее точную переформулировку:",
		"best_option":         "Выбери лучший вариант:",
		"natural_option":      "Выбери самый естественный вариант:",
		"verb_position":       "Выбери правильное место глагола:",
		"adjective_position":  "Выбери правильное место прилагательного:",
		"correct_adjective":   "Выбери правильное прилагательное:",
		"agreement":           "Выбери правильное согласование:",
		"dative_variant":      "Выбери правильный вариант с дательным падежом:",
		"correct_variant":     "Выбери правильную форму:",
		"how_to_say":          "Как сказать:",
		"what_means":          "Что значит",
		"question":            "Ответь на вопрос:",
		"choose_answer":       "Выбери правильный вариант:",
	}
	return labels[kind]
}

func localizedAssessmentTaskLabelForAllLanguages(kind string, interfaceLanguage string) string {
	code := normalizeInterfaceLanguage(interfaceLanguage)
	labels := map[string]map[string]string{
		"es": {"fill_gap": "Completa el espacio:", "correct_sentence": "Elige la frase correcta:", "correct_translation": "Elige la traducción correcta:", "precise_translation": "Elige la traducción más precisa:", "word_meaning": "Elige el significado más cercano de la palabra:", "phrase_meaning": "Elige el significado más cercano de la expresión:", "how_to_say": "¿Cómo se dice?:", "what_means": "Qué significa", "natural_option": "Elige la opción más natural:", "best_option": "Elige la mejor opción:", "correct_question": "Elige la pregunta correcta:", "choose_answer": "Elige la opción correcta:"},
		"de": {"fill_gap": "Fülle die Lücke:", "correct_sentence": "Wähle den richtigen Satz:", "correct_translation": "Wähle die richtige Übersetzung:", "precise_translation": "Wähle die genaueste Übersetzung:", "word_meaning": "Wähle die nächstliegende Bedeutung des Wortes:", "phrase_meaning": "Wähle die nächstliegende Bedeutung des Ausdrucks:", "how_to_say": "Wie sagt man?:", "what_means": "Was bedeutet", "natural_option": "Wähle die natürlichste Variante:", "best_option": "Wähle die beste Variante:", "correct_question": "Wähle die richtige Frage:", "choose_answer": "Wähle die richtige Antwort:"},
		"fr": {"fill_gap": "Complète le blanc :", "correct_sentence": "Choisis la phrase correcte :", "correct_translation": "Choisis la bonne traduction :", "precise_translation": "Choisis la traduction la plus précise :", "word_meaning": "Choisis le sens le plus proche du mot :", "phrase_meaning": "Choisis le sens le plus proche de l’expression :", "how_to_say": "Comment dire :", "what_means": "Que signifie", "natural_option": "Choisis la formulation la plus naturelle :", "best_option": "Choisis la meilleure option :", "correct_question": "Choisis la bonne question :", "choose_answer": "Choisis la bonne réponse :"},
		"it": {"fill_gap": "Completa lo spazio:", "correct_sentence": "Scegli la frase corretta:", "correct_translation": "Scegli la traduzione corretta:", "precise_translation": "Scegli la traduzione più precisa:", "word_meaning": "Scegli il significato più vicino della parola:", "phrase_meaning": "Scegli il significato più vicino dell’espressione:", "how_to_say": "Come si dice:", "what_means": "Che cosa significa", "natural_option": "Scegli l’opzione più naturale:", "best_option": "Scegli l’opzione migliore:", "correct_question": "Scegli la domanda corretta:", "choose_answer": "Scegli la risposta corretta:"},
		"pt": {"fill_gap": "Preenche a lacuna:", "correct_sentence": "Escolhe a frase correta:", "correct_translation": "Escolhe a tradução correta:", "precise_translation": "Escolhe a tradução mais precisa:", "word_meaning": "Escolhe o significado mais próximo da palavra:", "phrase_meaning": "Escolhe o significado mais próximo da expressão:", "how_to_say": "Como se diz:", "what_means": "O que significa", "natural_option": "Escolhe a opção mais natural:", "best_option": "Escolhe a melhor opção:", "correct_question": "Escolhe a pergunta correta:", "choose_answer": "Escolhe a resposta correta:"},
		"ro": {"fill_gap": "Completează spațiul liber:", "correct_sentence": "Alege propoziția corectă:", "correct_translation": "Alege traducerea corectă:", "precise_translation": "Alege cea mai precisă traducere:", "word_meaning": "Alege sensul cel mai apropiat al cuvântului:", "phrase_meaning": "Alege sensul cel mai apropiat al expresiei:", "how_to_say": "Cum se spune:", "what_means": "Ce înseamnă", "natural_option": "Alege varianta cea mai naturală:", "best_option": "Alege cea mai bună variantă:", "correct_question": "Alege întrebarea corectă:", "choose_answer": "Alege răspunsul corect:"},
		"pl": {"fill_gap": "Uzupełnij lukę:", "correct_sentence": "Wybierz poprawne zdanie:", "correct_translation": "Wybierz poprawne tłumaczenie:", "precise_translation": "Wybierz najdokładniejsze tłumaczenie:", "word_meaning": "Wybierz najbliższe znaczenie słowa:", "phrase_meaning": "Wybierz najbliższe znaczenie wyrażenia:", "how_to_say": "Jak powiedzieć:", "what_means": "Co znaczy", "natural_option": "Wybierz najbardziej naturalną opcję:", "best_option": "Wybierz najlepszą opcję:", "correct_question": "Wybierz poprawne pytanie:", "choose_answer": "Wybierz poprawną odpowiedź:"},
		"uk": {"fill_gap": "Заповни пропуск:", "correct_sentence": "Вибери правильне речення:", "correct_translation": "Вибери правильний переклад:", "precise_translation": "Вибери найточніший переклад:", "word_meaning": "Вибери найближче значення слова:", "phrase_meaning": "Вибери найближче значення виразу:", "how_to_say": "Як сказати:", "what_means": "Що означає", "natural_option": "Вибери найприродніший варіант:", "best_option": "Вибери найкращий варіант:", "correct_question": "Вибери правильне питання:", "choose_answer": "Вибери правильну відповідь:"},
		"zh": {"fill_gap": "填空：", "correct_sentence": "选择正确的句子：", "correct_translation": "选择正确的翻译：", "precise_translation": "选择最准确的翻译：", "word_meaning": "选择这个词最接近的意思：", "phrase_meaning": "选择这个表达最接近的意思：", "how_to_say": "怎么说：", "what_means": "是什么意思", "natural_option": "选择最自然的说法：", "best_option": "选择最佳选项：", "correct_question": "选择正确的问题：", "choose_answer": "选择正确答案："},
		"ja": {"fill_gap": "空欄を埋めてください：", "correct_sentence": "正しい文を選んでください：", "correct_translation": "正しい訳を選んでください：", "precise_translation": "最も正確な訳を選んでください：", "word_meaning": "単語に最も近い意味を選んでください：", "phrase_meaning": "表現に最も近い意味を選んでください：", "how_to_say": "どう言いますか：", "what_means": "の意味は", "natural_option": "最も自然な表現を選んでください：", "best_option": "最もよい選択肢を選んでください：", "correct_question": "正しい質問を選んでください：", "choose_answer": "正しい答えを選んでください："},
		"ko": {"fill_gap": "빈칸을 채우세요:", "correct_sentence": "올바른 문장을 고르세요:", "correct_translation": "올바른 번역을 고르세요:", "precise_translation": "가장 정확한 번역을 고르세요:", "word_meaning": "단어의 가장 가까운 뜻을 고르세요:", "phrase_meaning": "표현의 가장 가까운 뜻을 고르세요:", "how_to_say": "어떻게 말하나요:", "what_means": "무슨 뜻인가요", "natural_option": "가장 자연스러운 표현을 고르세요:", "best_option": "가장 좋은 선택지를 고르세요:", "correct_question": "올바른 질문을 고르세요:", "choose_answer": "정답을 고르세요:"},
		"ar": {"fill_gap": "املأ الفراغ:", "correct_sentence": "اختر الجملة الصحيحة:", "correct_translation": "اختر الترجمة الصحيحة:", "precise_translation": "اختر الترجمة الأدق:", "word_meaning": "اختر أقرب معنى للكلمة:", "phrase_meaning": "اختر أقرب معنى للتعبير:", "how_to_say": "كيف تقول:", "what_means": "ما معنى", "natural_option": "اختر الصياغة الأكثر طبيعية:", "best_option": "اختر أفضل خيار:", "correct_question": "اختر السؤال الصحيح:", "choose_answer": "اختر الإجابة الصحيحة:"},
		"bn": {"fill_gap": "শূন্যস্থান পূরণ করুন:", "correct_sentence": "সঠিক বাক্যটি বেছে নিন:", "correct_translation": "সঠিক অনুবাদ বেছে নিন:", "precise_translation": "সবচেয়ে নির্ভুল অনুবাদ বেছে নিন:", "word_meaning": "শব্দটির কাছাকাছি অর্থ বেছে নিন:", "phrase_meaning": "অভিব্যক্তিটির কাছাকাছি অর্থ বেছে নিন:", "how_to_say": "কীভাবে বলবেন:", "what_means": "এর অর্থ কী", "natural_option": "সবচেয়ে স্বাভাবিক বিকল্প বেছে নিন:", "best_option": "সেরা বিকল্প বেছে নিন:", "correct_question": "সঠিক প্রশ্নটি বেছে নিন:", "choose_answer": "সঠিক উত্তর বেছে নিন:"},
		"hi": {"fill_gap": "रिक्त स्थान भरें:", "correct_sentence": "सही वाक्य चुनें:", "correct_translation": "सही अनुवाद चुनें:", "precise_translation": "सबसे सटीक अनुवाद चुनें:", "word_meaning": "शब्द का सबसे निकट अर्थ चुनें:", "phrase_meaning": "अभिव्यक्ति का सबसे निकट अर्थ चुनें:", "how_to_say": "कैसे कहेंगे:", "what_means": "का क्या अर्थ है", "natural_option": "सबसे स्वाभाविक विकल्प चुनें:", "best_option": "सबसे अच्छा विकल्प चुनें:", "correct_question": "सही प्रश्न चुनें:", "choose_answer": "सही उत्तर चुनें:"},
		"ta": {"fill_gap": "வெற்றிடத்தை நிரப்புக:", "correct_sentence": "சரியான வாக்கியத்தைத் தேர்வு செய்க:", "correct_translation": "சரியான மொழிபெயர்ப்பைத் தேர்வு செய்க:", "precise_translation": "மிகத் துல்லியமான மொழிபெயர்ப்பைத் தேர்வு செய்க:", "word_meaning": "சொல்லின் அருகிய பொருளைத் தேர்வு செய்க:", "phrase_meaning": "சொற்றொடரின் அருகிய பொருளைத் தேர்வு செய்க:", "how_to_say": "எப்படி சொல்வது:", "what_means": "என்ன அர்த்தம்", "natural_option": "மிக இயல்பான விருப்பத்தைத் தேர்வு செய்க:", "best_option": "சிறந்த விருப்பத்தைத் தேர்வு செய்க:", "correct_question": "சரியான கேள்வியைத் தேர்வு செய்க:", "choose_answer": "சரியான பதிலைத் தேர்வு செய்க:"},
		"te": {"fill_gap": "ఖాళీని పూరించండి:", "correct_sentence": "సరైన వాక్యాన్ని ఎంచుకోండి:", "correct_translation": "సరైన అనువాదాన్ని ఎంచుకోండి:", "precise_translation": "అత్యంత ఖచ్చితమైన అనువాదాన్ని ఎంచుకోండి:", "word_meaning": "పదానికి దగ్గరైన అర్థాన్ని ఎంచుకోండి:", "phrase_meaning": "వ్యక్తీకరణకు దగ్గరైన అర్థాన్ని ఎంచుకోండి:", "how_to_say": "ఎలా చెబుతారు:", "what_means": "అర్థం ఏమిటి", "natural_option": "అత్యంత సహజమైన ఎంపికను ఎంచుకోండి:", "best_option": "ఉత్తమ ఎంపికను ఎంచుకోండి:", "correct_question": "సరైన ప్రశ్నను ఎంచుకోండి:", "choose_answer": "సరైన సమాధానం ఎంచుకోండి:"},
		"th": {"fill_gap": "เติมคำในช่องว่าง:", "correct_sentence": "เลือกประโยคที่ถูกต้อง:", "correct_translation": "เลือกคำแปลที่ถูกต้อง:", "precise_translation": "เลือกคำแปลที่แม่นยำที่สุด:", "word_meaning": "เลือกความหมายที่ใกล้ที่สุดของคำ:", "phrase_meaning": "เลือกความหมายที่ใกล้ที่สุดของสำนวน:", "how_to_say": "พูดว่าอย่างไร:", "what_means": "หมายความว่าอะไร", "natural_option": "เลือกตัวเลือกที่เป็นธรรมชาติที่สุด:", "best_option": "เลือกตัวเลือกที่ดีที่สุด:", "correct_question": "เลือกคำถามที่ถูกต้อง:", "choose_answer": "เลือกคำตอบที่ถูกต้อง:"},
		"tr": {"fill_gap": "Boşluğu doldur:", "correct_sentence": "Doğru cümleyi seç:", "correct_translation": "Doğru çeviriyi seç:", "precise_translation": "En doğru çeviriyi seç:", "word_meaning": "Kelimenin en yakın anlamını seç:", "phrase_meaning": "İfadenin en yakın anlamını seç:", "how_to_say": "Nasıl söylenir:", "what_means": "ne demek", "natural_option": "En doğal seçeneği seç:", "best_option": "En iyi seçeneği seç:", "correct_question": "Doğru soruyu seç:", "choose_answer": "Doğru cevabı seç:"},
		"vi": {"fill_gap": "Điền vào chỗ trống:", "correct_sentence": "Chọn câu đúng:", "correct_translation": "Chọn bản dịch đúng:", "precise_translation": "Chọn bản dịch chính xác nhất:", "word_meaning": "Chọn nghĩa gần nhất của từ:", "phrase_meaning": "Chọn nghĩa gần nhất của cụm từ:", "how_to_say": "Nói như thế nào:", "what_means": "có nghĩa là gì", "natural_option": "Chọn cách nói tự nhiên nhất:", "best_option": "Chọn phương án tốt nhất:", "correct_question": "Chọn câu hỏi đúng:", "choose_answer": "Chọn đáp án đúng:"},
		"tg": {"fill_gap": "Ҷойи холиро пур кунед:", "correct_sentence": "Ҷумлаи дурустро интихоб кунед:", "correct_translation": "Тарҷумаи дурустро интихоб кунед:", "precise_translation": "Тарҷумаи дақиқтаринро интихоб кунед:", "word_meaning": "Маънои наздиктарини калимаро интихоб кунед:", "phrase_meaning": "Маънои наздиктарини ибораро интихоб кунед:", "how_to_say": "Чӣ тавр мегӯянд:", "what_means": "чӣ маъно дорад", "natural_option": "Варианти табиитаринро интихоб кунед:", "best_option": "Варианти беҳтаринро интихоб кунед:", "correct_question": "Саволи дурустро интихоб кунед:", "choose_answer": "Ҷавоби дурустро интихоб кунед:"},
		"uz": {"fill_gap": "Bo‘sh joyni to‘ldiring:", "correct_sentence": "To‘g‘ri gapni tanlang:", "correct_translation": "To‘g‘ri tarjimani tanlang:", "precise_translation": "Eng aniq tarjimani tanlang:", "word_meaning": "So‘zning eng yaqin ma’nosini tanlang:", "phrase_meaning": "Iboraning eng yaqin ma’nosini tanlang:", "how_to_say": "Qanday aytiladi:", "what_means": "nimani anglatadi", "natural_option": "Eng tabiiy variantni tanlang:", "best_option": "Eng yaxshi variantni tanlang:", "correct_question": "To‘g‘ri savolni tanlang:", "choose_answer": "To‘g‘ri javobni tanlang:"},
		"tt": {"fill_gap": "Буш урынны тутыр:", "correct_sentence": "Дөрес җөмләне сайла:", "correct_translation": "Дөрес тәрҗемәне сайла:", "precise_translation": "Иң төгәл тәрҗемәне сайла:", "word_meaning": "Сүзнең иң якын мәгънәсен сайла:", "phrase_meaning": "Гыйбарәнең иң якын мәгънәсен сайла:", "how_to_say": "Ничек әйтергә:", "what_means": "нәрсә аңлата", "natural_option": "Иң табигый вариантны сайла:", "best_option": "Иң яхшы вариантны сайла:", "correct_question": "Дөрес сорауны сайла:", "choose_answer": "Дөрес җавапны сайла:"},
		"hy": {"fill_gap": "Լրացրեք բաց թողնվածը:", "correct_sentence": "Ընտրեք ճիշտ նախադասությունը:", "correct_translation": "Ընտրեք ճիշտ թարգմանությունը:", "precise_translation": "Ընտրեք ամենաճշգրիտ թարգմանությունը:", "word_meaning": "Ընտրեք բառի ամենամոտ իմաստը:", "phrase_meaning": "Ընտրեք արտահայտության ամենամոտ իմաստը:", "how_to_say": "Ինչպես ասել:", "what_means": "ինչ է նշանակում", "natural_option": "Ընտրեք ամենաբնական տարբերակը:", "best_option": "Ընտրեք լավագույն տարբերակը:", "correct_question": "Ընտրեք ճիշտ հարցը:", "choose_answer": "Ընտրեք ճիշտ պատասխանը:"},
		"kk": {"fill_gap": "Бос орынды толтыр:", "correct_sentence": "Дұрыс сөйлемді таңда:", "correct_translation": "Дұрыс аударманы таңда:", "precise_translation": "Ең дәл аударманы таңда:", "word_meaning": "Сөздің ең жақын мағынасын таңда:", "phrase_meaning": "Тіркестің ең жақын мағынасын таңда:", "how_to_say": "Қалай айтылады:", "what_means": "нені білдіреді", "natural_option": "Ең табиғи нұсқаны таңда:", "best_option": "Ең жақсы нұсқаны таңда:", "correct_question": "Дұрыс сұрақты таңда:", "choose_answer": "Дұрыс жауапты таңда:"},
		"ky": {"fill_gap": "Бош жерди толтур:", "correct_sentence": "Туура сүйлөмдү танда:", "correct_translation": "Туура котормону танда:", "precise_translation": "Эң так котормону танда:", "word_meaning": "Сөздүн эң жакын маанисин танда:", "phrase_meaning": "Ибаранын эң жакын маанисин танда:", "how_to_say": "Кантип айтылат:", "what_means": "эмнени билдирет", "natural_option": "Эң табигый вариантты танда:", "best_option": "Эң жакшы вариантты танда:", "correct_question": "Туура суроону танда:", "choose_answer": "Туура жоопту танда:"},
		"ka": {"fill_gap": "შეავსე გამოტოვებული ადგილი:", "correct_sentence": "აირჩიე სწორი წინადადება:", "correct_translation": "აირჩიე სწორი თარგმანი:", "precise_translation": "აირჩიე ყველაზე ზუსტი თარგმანი:", "word_meaning": "აირჩიე სიტყვის უახლოესი მნიშვნელობა:", "phrase_meaning": "აირჩიე გამოთქმის უახლოესი მნიშვნელობა:", "how_to_say": "როგორ ითქმის:", "what_means": "რას ნიშნავს", "natural_option": "აირჩიე ყველაზე ბუნებრივი ვარიანტი:", "best_option": "აირჩიე საუკეთესო ვარიანტი:", "correct_question": "აირჩიე სწორი კითხვა:", "choose_answer": "აირჩიე სწორი პასუხი:"},
		"cs": {"fill_gap": "Doplň mezeru:", "correct_sentence": "Vyber správnou větu:", "correct_translation": "Vyber správný překlad:", "precise_translation": "Vyber nejpřesnější překlad:", "word_meaning": "Vyber nejbližší význam slova:", "phrase_meaning": "Vyber nejbližší význam výrazu:", "how_to_say": "Jak se řekne:", "what_means": "co znamená", "natural_option": "Vyber nejpřirozenější možnost:", "best_option": "Vyber nejlepší možnost:", "correct_question": "Vyber správnou otázku:", "choose_answer": "Vyber správnou odpověď:"},
		"el": {"fill_gap": "Συμπλήρωσε το κενό:", "correct_sentence": "Διάλεξε τη σωστή πρόταση:", "correct_translation": "Διάλεξε τη σωστή μετάφραση:", "precise_translation": "Διάλεξε την πιο ακριβή μετάφραση:", "word_meaning": "Διάλεξε την πιο κοντινή σημασία της λέξης:", "phrase_meaning": "Διάλεξε την πιο κοντινή σημασία της έκφρασης:", "how_to_say": "Πώς λέγεται:", "what_means": "τι σημαίνει", "natural_option": "Διάλεξε την πιο φυσική επιλογή:", "best_option": "Διάλεξε την καλύτερη επιλογή:", "correct_question": "Διάλεξε τη σωστή ερώτηση:", "choose_answer": "Διάλεξε τη σωστή απάντηση:"},
		"hu": {"fill_gap": "Töltsd ki a hiányzó részt:", "correct_sentence": "Válaszd ki a helyes mondatot:", "correct_translation": "Válaszd ki a helyes fordítást:", "precise_translation": "Válaszd ki a legpontosabb fordítást:", "word_meaning": "Válaszd ki a szó legközelebbi jelentését:", "phrase_meaning": "Válaszd ki a kifejezés legközelebbi jelentését:", "how_to_say": "Hogyan mondjuk:", "what_means": "mit jelent", "natural_option": "Válaszd ki a legtermészetesebb változatot:", "best_option": "Válaszd ki a legjobb lehetőséget:", "correct_question": "Válaszd ki a helyes kérdést:", "choose_answer": "Válaszd ki a helyes választ:"},
		"id": {"fill_gap": "Isi bagian yang kosong:", "correct_sentence": "Pilih kalimat yang benar:", "correct_translation": "Pilih terjemahan yang benar:", "precise_translation": "Pilih terjemahan yang paling tepat:", "word_meaning": "Pilih arti kata yang paling dekat:", "phrase_meaning": "Pilih arti ungkapan yang paling dekat:", "how_to_say": "Bagaimana mengatakannya:", "what_means": "apa artinya", "natural_option": "Pilih opsi yang paling alami:", "best_option": "Pilih opsi terbaik:", "correct_question": "Pilih pertanyaan yang benar:", "choose_answer": "Pilih jawaban yang benar:"},
		"nl": {"fill_gap": "Vul de lege plek in:", "correct_sentence": "Kies de juiste zin:", "correct_translation": "Kies de juiste vertaling:", "precise_translation": "Kies de meest precieze vertaling:", "word_meaning": "Kies de dichtstbijzijnde betekenis van het woord:", "phrase_meaning": "Kies de dichtstbijzijnde betekenis van de uitdrukking:", "how_to_say": "Hoe zeg je:", "what_means": "wat betekent", "natural_option": "Kies de meest natuurlijke optie:", "best_option": "Kies de beste optie:", "correct_question": "Kies de juiste vraag:", "choose_answer": "Kies het juiste antwoord:"},
		"sv": {"fill_gap": "Fyll i luckan:", "correct_sentence": "Välj rätt mening:", "correct_translation": "Välj rätt översättning:", "precise_translation": "Välj den mest exakta översättningen:", "word_meaning": "Välj ordets närmaste betydelse:", "phrase_meaning": "Välj uttryckets närmaste betydelse:", "how_to_say": "Hur säger man:", "what_means": "vad betyder", "natural_option": "Välj det mest naturliga alternativet:", "best_option": "Välj det bästa alternativet:", "correct_question": "Välj rätt fråga:", "choose_answer": "Välj rätt svar:"},
		"tl": {"fill_gap": "Punan ang patlang:", "correct_sentence": "Piliin ang tamang pangungusap:", "correct_translation": "Piliin ang tamang salin:", "precise_translation": "Piliin ang pinakatumpak na salin:", "word_meaning": "Piliin ang pinakamalapit na kahulugan ng salita:", "phrase_meaning": "Piliin ang pinakamalapit na kahulugan ng parirala:", "how_to_say": "Paano sabihin:", "what_means": "ano ang ibig sabihin", "natural_option": "Piliin ang pinakanatural na opsyon:", "best_option": "Piliin ang pinakamahusay na opsyon:", "correct_question": "Piliin ang tamang tanong:", "choose_answer": "Piliin ang tamang sagot:"},
	}
	return assessmentTaskLabelWithFallback(labels, code, kind)
}

func assessmentTaskLabelWithFallback(labels map[string]map[string]string, code string, kind string) string {
	if localized, ok := labels[code][kind]; ok {
		return localized
	}
	if localized := localizedGrammarAssessmentTaskLabel(code); localized != "" {
		return localized
	}
	if localized, ok := labels[code]["correct_variant"]; ok {
		return localized
	}
	if localized, ok := labels[code]["choose_answer"]; ok {
		return localized
	}
	return ""
}

func localizedGrammarAssessmentTaskLabel(code string) string {
	labels := map[string]string{
		"es": "Elige la forma correcta:", "de": "Wähle die richtige Form:", "fr": "Choisis la bonne forme :", "it": "Scegli la forma corretta:", "pt": "Escolhe a forma correta:", "ro": "Alege forma corectă:", "pl": "Wybierz poprawną formę:", "uk": "Вибери правильну форму:",
		"zh": "选择正确形式：", "ja": "正しい形を選んでください：", "ko": "올바른 형태를 고르세요:", "tg": "Шакли дурустро интихоб кунед:", "uz": "To‘g‘ri shaklni tanlang:", "tt": "Дөрес форманы сайла:", "hy": "Ընտրեք ճիշտ ձևը:", "kk": "Дұрыс форманы таңда:", "ky": "Туура форманы танда:", "ka": "აირჩიე სწორი ფორმა:",
		"ar": "اختر الصيغة الصحيحة:", "bn": "সঠিক রূপটি বেছে নিন:", "cs": "Vyber správný tvar:", "el": "Διάλεξε τον σωστό τύπο:", "hi": "सही रूप चुनें:", "hu": "Válaszd ki a helyes alakot:", "id": "Pilih bentuk yang benar:", "nl": "Kies de juiste vorm:", "sv": "Välj rätt form:", "ta": "சரியான வடிவத்தைத் தேர்வு செய்க:", "te": "సరైన రూపాన్ని ఎంచుకోండి:", "th": "เลือกรูปแบบที่ถูกต้อง:", "tl": "Piliin ang tamang anyo:", "tr": "Doğru biçimi seç:", "vi": "Chọn dạng đúng:",
	}
	return labels[code]
}

func containsCyrillicText(text string) bool {
	for _, r := range text {
		if (r >= 'А' && r <= 'я') || r == 'Ё' || r == 'ё' {
			return true
		}
	}
	return false
}

func containsUnicodeCyrillicText(text string) bool {
	for _, r := range text {
		if (r >= 'А' && r <= 'я') || r == 'Ё' || r == 'ё' {
			return true
		}
	}
	return false
}

func levelAssessmentKeyboardForQuestion(question levelQuestion, index int, copy systemUICopy, backText string) map[string]any {
	rows := make([][]map[string]any, 0, len(question.Options)+2)
	for optionIndex, option := range question.Options {
		callback := fmt.Sprintf("lt|%d|%d", index, optionIndex)
		rows = append(rows, []map[string]any{{"text": option, "callback_data": callback}})
	}
	rows = append(rows, []map[string]any{{"text": copy.LevelDontKnow, "callback_data": fmt.Sprintf("lt|%d|-1", index)}})
	rows = append(rows, []map[string]any{{"text": backText, "callback_data": "back_menu"}})
	return map[string]any{"inline_keyboard": rows}
}

func vocabularyLevelAssessmentQuestions(user userState) []levelQuestion {
	language := userLearningLanguage(user)
	interfaceLanguage := normalizeInterfaceLanguage(user.InterfaceLanguage)
	cacheKey := language.Code + "|" + interfaceLanguage
	vocabularyLevelAssessmentCache.mu.RLock()
	if cached, ok := vocabularyLevelAssessmentCache.questions[cacheKey]; ok {
		vocabularyLevelAssessmentCache.mu.RUnlock()
		return cloneLevelQuestions(cached)
	}
	vocabularyLevelAssessmentCache.mu.RUnlock()

	wordsByLevel := map[string][]vocabWord{}
	eligibleWords := []vocabWord{}
	const maxWordsPerLevel = 16
	const maxFallbackWords = 128
	if err := forEachVocabularyWord(language.Code, func(word vocabWord) bool {
		level := normalizeCEFRLevel(word.Level)
		clue := levelAssessmentWordClue(word, interfaceLanguage)
		if clue == "" || strings.EqualFold(clue, word.English) {
			return true
		}
		if len(eligibleWords) < maxFallbackWords {
			eligibleWords = append(eligibleWords, word)
		}
		if len(wordsByLevel[level]) < maxWordsPerLevel {
			wordsByLevel[level] = append(wordsByLevel[level], word)
		}
		if len(eligibleWords) >= maxFallbackWords {
			return false
		}
		return true
	}); err != nil {
		panic(err)
	}

	const questionsPerCEFRLevel = 6
	questions := make([]levelQuestion, 0, len(cefrLevels)*questionsPerCEFRLevel)
	for levelIndex, level := range cefrLevels {
		words := wordsByLevel[level]
		if len(words) < 4 {
			words = eligibleWords
		}
		if len(words) < 4 {
			continue
		}
		for offset := 0; offset < questionsPerCEFRLevel && offset < len(words); offset++ {
			correctWord := words[(offset*7+levelIndex*11)%len(words)]
			options := make([]string, 0, 4)
			for optionOffset := 0; len(options) < 4 && optionOffset < len(words)*2; optionOffset++ {
				candidate := words[(offset*7+levelIndex*11+optionOffset*5)%len(words)]
				if candidate.English == "" || containsString(options, candidate.English) {
					continue
				}
				options = append(options, candidate.English)
			}
			if len(options) < 4 || !containsString(options, correctWord.English) {
				continue
			}
			clue := levelAssessmentWordClue(correctWord, interfaceLanguage)
			questions = append(questions, levelQuestion{
				Question:     vocabularyLevelAssessmentQuestionText(user, language, clue),
				Options:      options,
				CorrectIndex: indexOfString(options, correctWord.English),
				Points:       levelIndex + 1,
			})
		}
	}
	vocabularyLevelAssessmentCache.mu.Lock()
	vocabularyLevelAssessmentCache.questions[cacheKey] = cloneLevelQuestions(questions)
	vocabularyLevelAssessmentCache.mu.Unlock()
	return questions
}

func vocabularyLevelAssessmentQuestionText(user userState, language learningLanguage, clue string) string {
	copy := systemUI(user)
	if template := strings.TrimSpace(copy.VocabLevelQuestion); template != "" && template != "%s: %s" {
		return fmt.Sprintf(template, language.InterfaceName, clue)
	}
	wordQuestion := strings.TrimSpace(ui(user).WordQuestion)
	if wordQuestion == "" {
		wordQuestion = englishUICopy().WordQuestion
	}
	return fmt.Sprintf(wordQuestion, clue)
}

func cloneLevelQuestions(questions []levelQuestion) []levelQuestion {
	if len(questions) == 0 {
		return nil
	}
	cloned := make([]levelQuestion, len(questions))
	for index, question := range questions {
		cloned[index] = question
		cloned[index].Options = append([]string(nil), question.Options...)
	}
	return cloned
}

func levelAssessmentWordClue(word vocabWord, interfaceLanguage string) string {
	interfaceLanguage = normalizeInterfaceLanguage(interfaceLanguage)
	if word.Translations != nil {
		if clue := strings.TrimSpace(word.Translations[interfaceLanguage]); clue != "" {
			if !strings.EqualFold(clue, word.English) {
				return clue
			}
		}
		if interfaceLanguage != "ru" {
			if clue := strings.TrimSpace(word.Translations["en"]); clue != "" && !strings.EqualFold(clue, word.English) {
				return clue
			}
		}
	}
	if interfaceLanguage == "ru" {
		return strings.TrimSpace(word.Russian)
	}
	return ""
}

func containsString(items []string, value string) bool {
	return indexOfString(items, value) >= 0
}

func indexOfString(items []string, value string) int {
	for index, item := range items {
		if item == value {
			return index
		}
	}
	return -1
}

func parseLevelTestCallback(data string) (int, int, bool) {
	parts := strings.Split(data, "|")
	if len(parts) != 3 {
		return 0, 0, false
	}
	index, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, false
	}
	answer, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, false
	}
	if index < 0 || answer < -1 || answer > 3 {
		return 0, 0, false
	}
	return index, answer, true
}

func parseManualLevelCallback(data string) (string, bool) {
	parts := strings.Split(data, "|")
	if len(parts) != 2 || parts[0] != "level" {
		return "", false
	}
	level := strings.ToUpper(strings.TrimSpace(parts[1]))
	for _, item := range cefrLevels {
		if level == item {
			return item, true
		}
	}
	return "", false
}

func levelTestMode(index int, score int) string {
	return fmt.Sprintf("%s%d:%d", modeLevelTestPrefix, index, score)
}

func newLevelTestModeForUser(user userState) string {
	questions := levelAssessmentQuestionsForUser(user)
	return levelTestModeWithOrder(0, 0, shuffledLevelTestOrder(len(questions), user))
}

func levelTestModeWithOrder(index int, score int, order []int) string {
	if len(order) == 0 {
		return levelTestMode(index, score)
	}
	parts := make([]string, 0, len(order))
	for _, item := range order {
		parts = append(parts, strconv.Itoa(item))
	}
	return fmt.Sprintf("%s%d:%d:%s", modeLevelTestPrefix, index, score, strings.Join(parts, ","))
}

func shuffledLevelTestOrder(count int, user userState) []int {
	if count <= 0 {
		return nil
	}
	seed := time.Now().UnixNano() ^ user.TelegramID ^ int64(user.XP+user.WordGameCount+user.PracticeCount+user.LessonCount)
	order := rand.New(rand.NewSource(seed)).Perm(count)
	if count > 1 {
		identity := true
		for index, item := range order {
			if index != item {
				identity = false
				break
			}
		}
		if identity {
			order[0], order[1] = order[1], order[0]
		}
	}
	return order
}

func parseLevelTestMode(mode string) (int, int, bool) {
	index, score, _, ok := parseLevelTestModeState(mode)
	return index, score, ok
}

func parseLevelTestModeState(mode string) (int, int, []int, bool) {
	if !strings.HasPrefix(mode, modeLevelTestPrefix) {
		return 0, 0, nil, false
	}
	parts := strings.Split(strings.TrimPrefix(mode, modeLevelTestPrefix), ":")
	if len(parts) != 2 && len(parts) != 3 {
		return 0, 0, nil, false
	}
	index, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, nil, false
	}
	score, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, nil, false
	}
	if index < 0 || score < 0 {
		return 0, 0, nil, false
	}
	var order []int
	if len(parts) == 3 {
		for _, raw := range strings.Split(parts[2], ",") {
			item, err := strconv.Atoi(strings.TrimSpace(raw))
			if err != nil || item < 0 {
				return 0, 0, nil, false
			}
			order = append(order, item)
		}
		if len(order) == 0 {
			return 0, 0, nil, false
		}
	}
	return index, score, order, true
}

func validLevelTestOrder(order []int, questionCount int) bool {
	if len(order) != questionCount {
		return false
	}
	seen := make([]bool, questionCount)
	for _, item := range order {
		if item < 0 || item >= questionCount || seen[item] {
			return false
		}
		seen[item] = true
	}
	return true
}

type levelPromotionRule struct {
	Words    int
	Lessons  int
	Practice int
	Games    int
}

var levelPromotionRules = map[string]levelPromotionRule{
	"A1": {Words: 120, Lessons: 5, Practice: 15, Games: 5},
	"A2": {Words: 300, Lessons: 15, Practice: 40, Games: 15},
	"B1": {Words: 700, Lessons: 35, Practice: 90, Games: 40},
	"B2": {Words: 1400, Lessons: 70, Practice: 180, Games: 80},
	"C1": {Words: 2500, Lessons: 120, Practice: 320, Games: 150},
}

func shouldPromoteLearningLevel(user userState) (string, bool) {
	current := normalizeCEFRLevel(user.Level)
	next, ok := nextCEFRLevel(current)
	if !ok {
		return current, false
	}
	rule, ok := levelPromotionRules[current]
	if !ok {
		return current, false
	}
	if len(learnedWordsForLanguage(user)) >= rule.Words &&
		user.LessonCount >= rule.Lessons &&
		user.PracticeCount >= rule.Practice &&
		user.WordGameCount >= rule.Games {
		return next, true
	}
	return current, false
}
