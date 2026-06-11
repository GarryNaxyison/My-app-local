package main

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tutorLessonWord struct {
	ID           string `json:"id"`
	Word         string `json:"word"`
	Translation  string `json:"translation"`
	Context      string `json:"context,omitempty"`
	Example      string `json:"example"`
	ExampleRu    string `json:"example_translation,omitempty"`
	PartOfSpeech string `json:"part_of_speech,omitempty"`
	Topic        string `json:"topic"`
	Level        string `json:"level"`
	Frequency    int    `json:"frequency,omitempty"`
	AudioText    string `json:"audio_text"`
}

type tutorLessonChoice struct {
	Prompt          string              `json:"prompt"`
	Options         []tutorChoiceOption `json:"options"`
	CorrectAnswerID string              `json:"correct_answer_id"`
	Feedback        string              `json:"feedback,omitempty"`
}

type tutorChoiceOption struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Label   string `json:"label,omitempty"`
	Why     string `json:"why,omitempty"`
	Skill   string `json:"skill,omitempty"`
	Quality string `json:"quality,omitempty"`
	Avoid   bool   `json:"avoid,omitempty"`
}

type tutorAnswerVariant struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Text    string `json:"text"`
	Why     string `json:"why"`
	Level   string `json:"level,omitempty"`
	UseCase string `json:"use_case,omitempty"`
	Avoid   bool   `json:"avoid,omitempty"`
}

type tutorScenarioSlots struct {
	Role           string   `json:"role"`
	Situation      string   `json:"situation"`
	RequiredAction string   `json:"required_action"`
	Item           string   `json:"item"`
	Detail         string   `json:"detail"`
	Politeness     []string `json:"politeness,omitempty"`
	ModelAnswer    string   `json:"model_answer"`
}

type tutorTeachingPoint struct {
	Title       string   `json:"title"`
	Pattern     string   `json:"pattern"`
	SceneOrder  []string `json:"scene_order"`
	Explanation string   `json:"explanation"`
	ModelAnswer string   `json:"model_answer"`
}

type tutorFinalWordCheck struct {
	Prompt          string              `json:"prompt"`
	Items           []tutorLessonChoice `json:"items"`
	RequiredCorrect int                 `json:"required_correct"`
	SummaryPass     string              `json:"summary_pass"`
	SummaryRetry    string              `json:"summary_retry"`
}

type tutorSummaryBlock struct {
	CanSay      string   `json:"can_say"`
	StrongItems []string `json:"strong_items"`
	WeakItems   []string `json:"weak_items"`
	NextReview  string   `json:"next_review"`
}

type tutorLessonSRS struct {
	Again string `json:"again"`
	Hard  string `json:"hard"`
	Good  string `json:"good"`
	Easy  string `json:"easy"`
}

type tutorLessonStep struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type tutorLesson struct {
	ID                string               `json:"id"`
	Title             string               `json:"title"`
	Level             string               `json:"level"`
	Topic             string               `json:"topic"`
	VariantCode       string               `json:"variant_code,omitempty"`
	VariantTitle      string               `json:"variant_title,omitempty"`
	FocusSkill        string               `json:"focus_skill,omitempty"`
	SuccessCriteria   []string             `json:"success_criteria,omitempty"`
	LessonNumber      int                  `json:"lesson_number"`
	CourseSize        int                  `json:"course_size"`
	LearningLanguage  string               `json:"learning_language"`
	InterfaceLanguage string               `json:"interface_language"`
	DurationMinutes   int                  `json:"duration_minutes"`
	Goal              string               `json:"goal"`
	CanDo             string               `json:"can_do"`
	Scenario          string               `json:"scenario"`
	ScenarioSlots     tutorScenarioSlots   `json:"scenario_slots,omitempty"`
	TeachingPoint     tutorTeachingPoint   `json:"teaching_point,omitempty"`
	FinalWordCheck    tutorFinalWordCheck  `json:"final_word_check,omitempty"`
	TutorSummary      tutorSummaryBlock    `json:"tutor_summary,omitempty"`
	Steps             []tutorLessonStep    `json:"steps"`
	Words             []tutorLessonWord    `json:"words"`
	GrammarTitle      string               `json:"grammar_title"`
	Grammar           string               `json:"grammar"`
	MiniExplanation   string               `json:"mini_explanation"`
	Choice            tutorLessonChoice    `json:"choice"`
	Checks            []tutorLessonChoice  `json:"checks,omitempty"`
	WritingTask       string               `json:"writing_task"`
	WritingExpected   []string             `json:"writing_expected,omitempty"`
	AnswerVariants    []tutorAnswerVariant `json:"answer_variants,omitempty"`
	ListeningTask     string               `json:"listening_task"`
	ListeningText     string               `json:"listening_text"`
	ListeningQuestion string               `json:"listening_question,omitempty"`
	ListeningExpected []string             `json:"listening_expected,omitempty"`
	PronunciationText string               `json:"pronunciation_text,omitempty"`
	Dialogue          []string             `json:"dialogue"`
	DialoguePrompt    string               `json:"dialogue_prompt,omitempty"`
	DialogueGoal      string               `json:"dialogue_goal,omitempty"`
	DialogueVariants  []tutorAnswerVariant `json:"dialogue_variants,omitempty"`
	Review            []tutorLessonWord    `json:"review"`
	ReviewSummary     []string             `json:"review_summary,omitempty"`
	SRS               tutorLessonSRS       `json:"srs"`
	NextActions       []string             `json:"next_actions"`
	Source            string               `json:"source"`
}

type tutorTopic struct {
	Code       string
	TitleRU    string
	TitleEN    string
	ScenarioRU string
	ScenarioEN string
	Words      []string
}

type tutorCourseFunction struct {
	Code    string
	TitleRU string
	TitleEN string
	CanDoRU string
	CanDoEN string
}

type tutorCourseVariant struct {
	Code         string
	TitleRU      string
	TitleEN      string
	FocusRU      string
	FocusEN      string
	ConstraintRU string
	ConstraintEN string
}

type tutorScenarioChoiceLine struct {
	Text    string
	LabelRU string
	LabelEN string
	WhyRU   string
	WhyEN   string
	Quality string
}

type tutorScenarioVariantLine struct {
	LabelRU   string
	LabelEN   string
	Text      string
	WhyRU     string
	WhyEN     string
	Level     string
	UseCaseRU string
	UseCaseEN string
	Avoid     bool
}

type tutorScenarioTemplate struct {
	TopicCode           string
	Role                string
	SituationRU         string
	SituationEN         string
	RequiredAction      string
	Item                string
	Detail              string
	Politeness          []string
	ModelAnswer         string
	ShortAnswer         string
	WrongSceneAnswer    string
	IncompleteAnswer    string
	DetailPromptRU      string
	DetailPromptEN      string
	DetailDistractors   []tutorScenarioChoiceLine
	NextPromptRU        string
	NextPromptEN        string
	NextAnswer          string
	NextDistractors     []tutorScenarioChoiceLine
	MiniExplanationRU   string
	MiniExplanationEN   string
	WritingTaskRU       string
	WritingTaskEN       string
	WritingExpected     []string
	ListeningText       string
	ListeningQuestionRU string
	ListeningQuestionEN string
	ListeningExpected   []string
	DialogueLines       []string
	DialoguePromptRU    string
	DialoguePromptEN    string
	DialogueGoalRU      string
	DialogueGoalEN      string
	WritingVariants     []tutorScenarioVariantLine
	DialogueVariants    []tutorScenarioVariantLine
}

const tutorCourseSize = 300

var tutorScenarioCourseTopics = []tutorTopic{
	{Code: "hotel", TitleRU: "Отель и поездка", TitleEN: "Hotel and travel", ScenarioRU: "Вы на ресепшене отеля: нужно подтвердить бронь, уточнить номер и показать документ.", ScenarioEN: "You are at a hotel reception: confirm a reservation, ask about the room, and show a document.", Words: []string{"hotel", "room", "reservation", "passport", "ticket", "station", "airport", "street", "address"}},
	{Code: "food", TitleRU: "Еда и заказ", TitleEN: "Food and ordering", ScenarioRU: "Вы в кафе: нужно попросить меню, заказать напиток и уточнить завтрак или ужин.", ScenarioEN: "You are in a cafe: ask for the menu, order a drink, and clarify breakfast or dinner.", Words: []string{"food", "water", "coffee", "tea", "bread", "menu", "restaurant", "breakfast", "dinner"}},
	{Code: "city", TitleRU: "Город и маршруты", TitleEN: "City and directions", ScenarioRU: "Вы в городе: нужно спросить дорогу, назвать адрес и понять, где находится нужное место.", ScenarioEN: "You are in the city: ask for directions, give an address, and understand where a place is.", Words: []string{"city", "street", "house", "shop", "market", "bank", "school", "office", "park"}},
	{Code: "work", TitleRU: "Работа и встреча", TitleEN: "Work and meetings", ScenarioRU: "Вы на работе: нужно подтвердить встречу, время, имя и контакт.", ScenarioEN: "You are at work: confirm a meeting, time, name, and contact.", Words: []string{"work", "job", "office", "meeting", "email", "phone", "name", "number", "time"}},
	{Code: "doctor", TitleRU: "Врач и самочувствие", TitleEN: "Doctor and health", ScenarioRU: "Вы в клинике: нужно объяснить симптом, попросить врача и договориться о времени.", ScenarioEN: "You are at a clinic: explain a symptom, ask for a doctor, and agree on a time.", Words: []string{"doctor", "hospital", "pain", "medicine", "water", "help", "today", "tomorrow", "morning"}},
	{Code: "shopping", TitleRU: "Магазин и покупка", TitleEN: "Shopping", ScenarioRU: "Вы в магазине: нужно спросить цену, размер, оплату и чек.", ScenarioEN: "You are in a shop: ask about price, size, payment, and a receipt.", Words: []string{"shop", "market", "price", "money", "card", "bag", "size", "new", "good"}},
	{Code: "home", TitleRU: "Дом и быт", TitleEN: "Home and daily life", ScenarioRU: "Вы говорите о доме: нужно назвать адрес, комнату, ключ и бытовую проблему.", ScenarioEN: "You are talking about home: give an address, mention a room, a key, and a daily problem.", Words: []string{"home", "house", "room", "key", "address", "water", "door", "family", "today"}},
	{Code: "study", TitleRU: "Учёба и экзамен", TitleEN: "Study and exams", ScenarioRU: "Вы на занятии: нужно спросить задание, уточнить время и подготовиться к экзамену.", ScenarioEN: "You are in class: ask about a task, clarify time, and prepare for an exam.", Words: []string{"school", "lesson", "exam", "question", "answer", "book", "time", "help", "today"}},
	{Code: "travel", TitleRU: "Транспорт и билеты", TitleEN: "Transport and tickets", ScenarioRU: "Вы на станции или в аэропорту: нужно купить билет, уточнить время и направление.", ScenarioEN: "You are at a station or airport: buy a ticket, clarify time, and confirm direction.", Words: []string{"ticket", "station", "airport", "train", "bus", "time", "street", "city", "tomorrow"}},
	{Code: "services", TitleRU: "Сервисы и документы", TitleEN: "Services and documents", ScenarioRU: "Вы обращаетесь в сервис: нужно назвать номер, документ, проблему и следующий шаг.", ScenarioEN: "You are at a service desk: give a number, show a document, explain a problem, and agree on the next step.", Words: []string{"bank", "phone", "number", "name", "passport", "help", "problem", "address", "today"}},
}

var tutorCourseFunctions = []tutorCourseFunction{
	{Code: "ask-help", TitleRU: "Попросить помощь", TitleEN: "Ask for help", CanDoRU: "попросить помощь в конкретной ситуации", CanDoEN: "ask for help in a specific situation"},
	{Code: "give-details", TitleRU: "Дать детали", TitleEN: "Give details", CanDoRU: "дать короткие детали: имя, время, место или предмет", CanDoEN: "give short details: name, time, place, or item"},
	{Code: "confirm-correct", TitleRU: "Подтвердить и исправить", TitleEN: "Confirm and correct", CanDoRU: "подтвердить информацию и вежливо исправить ошибку", CanDoEN: "confirm information and politely correct a mistake"},
	{Code: "solve-problem", TitleRU: "Решить проблему", TitleEN: "Solve a problem", CanDoRU: "объяснить простую проблему и попросить следующий шаг", CanDoEN: "explain a simple problem and ask for the next step"},
	{Code: "close-next", TitleRU: "Закрыть диалог", TitleEN: "Close and move on", CanDoRU: "закрыть короткий диалог и договориться о следующем действии", CanDoEN: "close a short dialogue and agree on the next action"},
}

var tutorCourseVariants = []tutorCourseVariant{
	{Code: "polite-minimum", TitleRU: "Polite minimum", TitleEN: "Polite minimum", FocusRU: "short polite line", FocusEN: "short polite line", ConstraintRU: "Say only the useful minimum: opener, key word, one detail.", ConstraintEN: "Say only the useful minimum: opener, key word, one detail."},
	{Code: "detail-upgrade", TitleRU: "Detail upgrade", TitleEN: "Detail upgrade", FocusRU: "add one concrete detail", FocusEN: "add one concrete detail", ConstraintRU: "Add one detail that helps the other person act.", ConstraintEN: "Add one detail that helps the other person act."},
	{Code: "clarify-correct", TitleRU: "Clarify and correct", TitleEN: "Clarify and correct", FocusRU: "clarify or repair information", FocusEN: "clarify or repair information", ConstraintRU: "Clarify one unclear point or politely correct one detail.", ConstraintEN: "Clarify one unclear point or politely correct one detail."},
	{Code: "time-pressure", TitleRU: "Time pressure", TitleEN: "Time pressure", FocusRU: "answer fast with time words", FocusEN: "answer fast with time words", ConstraintRU: "Use a time word and keep the line direct.", ConstraintEN: "Use a time word and keep the line direct."},
	{Code: "compare-options", TitleRU: "Compare options", TitleEN: "Compare options", FocusRU: "choose between two options", FocusEN: "choose between two options", ConstraintRU: "Compare two options and choose one.", ConstraintEN: "Compare two options and choose one."},
	{Code: "repair-followup", TitleRU: "Repair follow-up", TitleEN: "Repair follow-up", FocusRU: "ask the next repair question", FocusEN: "ask the next repair question", ConstraintRU: "End with a follow-up question or next step.", ConstraintEN: "End with a follow-up question or next step."},
}

var tutorCourseTopics = []tutorTopic{
	{Code: "hotel", TitleRU: "Отель и поездка", TitleEN: "Hotel and travel", Words: []string{"hotel", "room", "reservation", "passport", "ticket", "station", "airport", "street", "address"}},
	{Code: "food", TitleRU: "Еда и заказ", TitleEN: "Food and ordering", Words: []string{"food", "water", "coffee", "tea", "bread", "menu", "restaurant", "breakfast", "dinner"}},
	{Code: "city", TitleRU: "Город", TitleEN: "City", Words: []string{"city", "street", "house", "shop", "market", "bank", "school", "office", "park"}},
	{Code: "work", TitleRU: "Работа", TitleEN: "Work", Words: []string{"work", "job", "office", "meeting", "email", "phone", "name", "number", "time"}},
	{Code: "doctor", TitleRU: "Врач и самочувствие", TitleEN: "Doctor and health", Words: []string{"doctor", "hospital", "pain", "medicine", "water", "help", "today", "tomorrow", "morning"}},
}

func tutorScenarioTemplateFor(topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) tutorScenarioTemplate {
	switch topic.Code {
	case "food":
		return tutorScenarioTemplate{
			TopicCode:           "food",
			Role:                "guest",
			SituationRU:         "Вы в кафе. Нужно попросить меню, заказать напиток и уточнить завтрак или ужин.",
			SituationEN:         "You are in a cafe: ask for the menu, order a drink, and clarify breakfast or dinner.",
			RequiredAction:      "ask for the menu",
			Item:                "coffee",
			Detail:              "for breakfast",
			Politeness:          []string{"Could I...?", "I'd like..., please."},
			ModelAnswer:         "Could I see the menu and have coffee for breakfast, please?",
			ShortAnswer:         "Coffee.",
			WrongSceneAnswer:    "I need a ticket.",
			IncompleteAnswer:    "Could I see the menu, please?",
			DetailPromptRU:      "Что уточняет заказ?",
			DetailPromptEN:      "What clarifies the order?",
			DetailDistractors:   []tutorScenarioChoiceLine{tutorScenarioChoice("at the airport", "Не эта сцена", "Wrong scene", "Это не кафе.", "This is not the cafe scene.", "wrong"), tutorScenarioChoice("my passport", "Не нужно в кафе", "Not needed here", "Паспорт не нужен для заказа кофе.", "A passport is not needed for a coffee order.", "wrong"), tutorScenarioChoice("football", "Не связано с задачей", "Unrelated", "Это слово из словаря, а не деталь заказа.", "This is vocabulary, not an order detail.", "wrong")},
			NextPromptRU:        "Какой следующий шаг звучит естественно?",
			NextPromptEN:        "Which next step sounds natural?",
			NextAnswer:          "Thank you. That is all for now.",
			NextDistractors:     []tutorScenarioChoiceLine{tutorScenarioChoice("I need a passport.", "Не по сцене", "Wrong scene", "Это не нужно в кафе.", "This is not needed in a cafe.", "wrong"), tutorScenarioChoice("Coffee. Football.", "Список слов", "Word list", "Это не ответ официанту.", "This is not an answer to the server.", "wrong"), tutorScenarioChoice("Close the app.", "Пустая реакция", "Empty reaction", "Собеседник ждёт завершения заказа.", "The other person expects you to finish the order.", "wrong")},
			MiniExplanationRU:   "В кафе нужен короткий порядок: попросить меню, затем заказать coffee, затем уточнить for breakfast. Паттерн: Could I see the menu? I'd like coffee for breakfast, please.",
			MiniExplanationEN:   "In a cafe: ask for the menu, then order coffee, then clarify for breakfast. Pattern: Could I see the menu? I'd like coffee for breakfast, please.",
			WritingTaskRU:       "Напишите одну живую реплику для кафе: вежливое начало, заказ coffee и деталь for breakfast.",
			WritingTaskEN:       "Write one live cafe reply: polite opener, coffee order, and the detail for breakfast.",
			WritingExpected:     []string{"coffee", "for breakfast"},
			ListeningText:       "Good morning. Could I see the menu? I'd like coffee for breakfast, please.",
			ListeningQuestionRU: "Какие 2-3 детали важны в заказе?",
			ListeningQuestionEN: "Which 2-3 details matter in the order?",
			ListeningExpected:   []string{"menu", "coffee", "breakfast"},
			DialogueLines:       []string{"Server: Good morning. Would you like breakfast or dinner?"},
			DialoguePromptRU:    "Ответьте официанту: выберите breakfast и закажите coffee.",
			DialoguePromptEN:    "Answer the server: choose breakfast and order coffee.",
			DialogueGoalRU:      "Заказать coffee и уточнить for breakfast.",
			DialogueGoalEN:      "Order coffee and clarify for breakfast.",
			WritingVariants: []tutorScenarioVariantLine{
				tutorScenarioVariant("Базовый ответ", "Base answer", "Coffee for breakfast, please.", "Есть заказ coffee и деталь for breakfast.", "It has the coffee order and the detail for breakfast.", "A1", "минимум", "safe minimum", false),
				tutorScenarioVariant("Естественно", "Natural answer", "I'd like coffee for breakfast, please.", "Звучит естественно: есть вежливость, заказ и деталь.", "Natural: it has politeness, order, and detail.", "A2", "живой разговор", "daily conversation", false),
				tutorScenarioVariant("Сильный вариант", "Stronger answer", "Could I see the menu first? I'd like coffee for breakfast, please.", "Есть просьба о меню, заказ coffee и деталь for breakfast.", "It asks for the menu and includes the coffee order with the breakfast detail.", "B1", "лучший образец", "best practice model", false),
				tutorScenarioVariant("Не использовать", "Do not use", "cup / football.", "Это список слов, а не ответ человеку.", "This is a word list, not an answer to a person.", "avoid", "не повторять", "avoid", true),
			},
			DialogueVariants: []tutorScenarioVariantLine{
				tutorScenarioVariant("A1", "A1", "Coffee for breakfast, please.", "Хорошо: есть заказ coffee и деталь for breakfast.", "Good: it has coffee and for breakfast.", "A1", "минимум", "safe minimum", false),
				tutorScenarioVariant("A2", "A2", "I'd like coffee for breakfast, please.", "Хорошо: звучит вежливо и закрывает вопрос официанта.", "Good: it sounds polite and answers the server's question.", "A2", "живой разговор", "daily conversation", false),
				tutorScenarioVariant("B1", "B1", "Could I see the menu first? I'd like coffee for breakfast, please.", "Сильнее: есть меню, заказ coffee и деталь for breakfast.", "Stronger: it includes the menu, coffee, and for breakfast.", "B1", "лучший образец", "best practice model", false),
				tutorScenarioVariant("Не использовать", "Do not use", "cup / football.", "Это список слов, а не ответ человеку.", "This is a word list, not an answer to a person.", "avoid", "не повторять", "avoid", true),
			},
		}
	case "hotel":
		return tutorBasicScenarioTemplate(topic, "guest", "confirm a reservation", "reservation", "passport", "I have a reservation. Here is my passport, please.", "Receptionist: Good evening. May I see your passport, please?")
	case "work":
		return tutorScenarioTemplate{
			TopicCode:        "work",
			Role:             "employee",
			SituationRU:      "Вы на работе. Нужно подтвердить встречу, время, имя и контакт.",
			SituationEN:      "You are at work: confirm a meeting, time, name, and contact.",
			RequiredAction:   "confirm a meeting",
			Item:             "meeting",
			Detail:           "with Alex at 10",
			Politeness:       []string{"Could you...?", "please"},
			ModelAnswer:      "Could you help me confirm the meeting with Alex at 10 and send me the contact, please?",
			ShortAnswer:      "Meeting at 10.",
			WrongSceneAnswer: "Can I order coffee?",
			IncompleteAnswer: "Could you help me with the meeting, please?",
			DetailPromptRU:   "Что уточняет встречу?",
			DetailPromptEN:   "What clarifies the meeting?",
			DetailDistractors: []tutorScenarioChoiceLine{
				tutorScenarioChoice("for breakfast", "Не эта сцена", "Wrong scene", "Это деталь заказа в кафе, не рабочей встречи.", "This belongs to a cafe order, not a work meeting.", "wrong"),
				tutorScenarioChoice("my passport", "Не нужно здесь", "Not needed here", "Паспорт не помогает подтвердить рабочую встречу.", "A passport does not help confirm a work meeting.", "wrong"),
				tutorScenarioChoice("football", "Не связано с задачей", "Unrelated", "Это слово из словаря, а не деталь встречи.", "This is vocabulary, not a meeting detail.", "wrong"),
			},
			NextPromptRU:        "Какой следующий шаг звучит естественно?",
			NextPromptEN:        "Which next step sounds natural?",
			NextAnswer:          "Thank you. Please send the contact by email.",
			NextDistractors:     []tutorScenarioChoiceLine{tutorScenarioChoice("I need a ticket.", "Не по сцене", "Wrong scene", "Это транспорт, не рабочая встреча.", "This is transport, not a work meeting.", "wrong"), tutorScenarioChoice("Meeting. Alex. Contact.", "Список слов", "Word list", "Это не звучит как реплика человеку.", "This does not sound like a line to a person.", "wrong"), tutorScenarioChoice("Close the app.", "Пустая реакция", "Empty reaction", "Собеседник ждёт следующий шаг.", "The other person expects a next step.", "wrong")},
			MiniExplanationRU:   "В рабочей сцене держите один порядок: вежливое начало -> meeting -> время и имя -> contact. Паттерн: Could you help me confirm the meeting with Alex at 10 and send me the contact, please.",
			MiniExplanationEN:   "In a work scene, keep one order: polite opener -> meeting -> time and name -> contact. Pattern: Could you help me confirm the meeting with Alex at 10 and send me the contact, please.",
			WritingTaskRU:       "Напишите одну рабочую реплику: подтвердить meeting, назвать Alex at 10 и попросить contact.",
			WritingTaskEN:       "Write one work reply: confirm the meeting, say Alex at 10, and ask for the contact.",
			WritingExpected:     []string{"meeting", "Alex", "10", "contact"},
			ListeningText:       "Could you help me confirm the meeting with Alex at 10 and send me the contact, please?",
			ListeningQuestionRU: "Какие 2-3 детали важны для встречи?",
			ListeningQuestionEN: "Which 2-3 details matter for the meeting?",
			ListeningExpected:   []string{"meeting", "Alex", "10", "contact"},
			DialogueLines:       []string{"Coworker: Good morning. Which meeting should I confirm?"},
			DialoguePromptRU:    "Ответьте коллеге: подтвердите meeting with Alex at 10 и попросите contact.",
			DialoguePromptEN:    "Answer the coworker: confirm the meeting with Alex at 10 and ask for the contact.",
			DialogueGoalRU:      "Подтвердить meeting, время, имя и contact.",
			DialogueGoalEN:      "Confirm the meeting, time, name, and contact.",
			WritingVariants: []tutorScenarioVariantLine{
				tutorScenarioVariant("Базовый ответ", "Base answer", "Meeting with Alex at 10, please.", "Есть встреча, имя и время, но контакт ещё не запрошен.", "It has the meeting, name, and time, but not the contact yet.", "A1", "минимум", "safe minimum", false),
				tutorScenarioVariant("Естественно", "Natural answer", "Could you confirm the meeting with Alex at 10, please?", "Звучит как рабочая просьба: есть действие, meeting, имя и время.", "It sounds like a work request: action, meeting, name, and time.", "A2", "живой разговор", "daily conversation", false),
				tutorScenarioVariant("Сильный вариант", "Stronger answer", "Could you help me confirm the meeting with Alex at 10 and send me the contact, please?", "Закрывает всю задачу: meeting, время, имя и contact.", "It closes the whole task: meeting, time, name, and contact.", "B1", "лучший образец", "best practice model", false),
				tutorScenarioVariant("Не использовать", "Do not use", "Meeting. Alex. Contact.", "Это список слов, а не ответ человеку.", "This is a word list, not an answer to a person.", "avoid", "не повторять", "avoid", true),
			},
			DialogueVariants: []tutorScenarioVariantLine{
				tutorScenarioVariant("A1", "A1", "The meeting with Alex at 10, please.", "Коротко отвечает на вопрос коллеги.", "Shortly answers the coworker's question.", "A1", "минимум", "safe minimum", false),
				tutorScenarioVariant("A2", "A2", "Please confirm the meeting with Alex at 10.", "Есть действие, meeting, имя и время.", "It has action, meeting, name, and time.", "A2", "живой разговор", "daily conversation", false),
				tutorScenarioVariant("B1", "B1", "Please confirm the meeting with Alex at 10 and send me the contact.", "Полный ответ на последнюю реплику: meeting, время, имя и contact.", "A complete answer to the last line: meeting, time, name, and contact.", "B1", "лучший образец", "best practice model", false),
				tutorScenarioVariant("Не использовать", "Do not use", "Meeting. Alex. Contact.", "Это список слов, а не ответ человеку.", "This is a word list, not an answer to a person.", "avoid", "не повторять", "avoid", true),
			},
		}
	case "doctor":
		return tutorBasicScenarioTemplate(topic, "patient", "ask for a doctor", "doctor", "today", "Hello. I need to see a doctor today, please.", "Receptionist: Can you come this morning?")
	case "travel":
		return tutorBasicScenarioTemplate(topic, "traveler", "buy a ticket", "ticket", "tomorrow", "I'd like a ticket for tomorrow, please.", "Clerk: Do you need a ticket for today or tomorrow?")
	case "services":
		return tutorBasicScenarioTemplate(topic, "customer", "ask for service help", "help", "passport", "I need help with my number. Here is my passport, please.", "Agent: What document do you have?")
	default:
		return tutorBasicScenarioTemplate(topic, "learner", function.CanDoEN, "help", "today", "Could you help me with this today, please?", "Assistant: What detail should I write?")
	}
}

func tutorBasicScenarioTemplate(topic tutorTopic, role string, action string, item string, detail string, modelAnswer string, lastTutorLine string) tutorScenarioTemplate {
	situationRU := topic.ScenarioRU
	if strings.TrimSpace(situationRU) == "" {
		situationRU = "Сцена: " + topic.TitleRU + "."
	}
	situationEN := topic.ScenarioEN
	if strings.TrimSpace(situationEN) == "" {
		situationEN = "Scene: " + topic.TitleEN + "."
	}
	return tutorScenarioTemplate{
		TopicCode:           topic.Code,
		Role:                role,
		SituationRU:         situationRU,
		SituationEN:         situationEN,
		RequiredAction:      action,
		Item:                item,
		Detail:              detail,
		Politeness:          []string{"Could I...?", "please"},
		ModelAnswer:         modelAnswer,
		ShortAnswer:         item + ".",
		WrongSceneAnswer:    "I need a coffee.",
		IncompleteAnswer:    "Could you help me, please?",
		DetailPromptRU:      "Какая деталь двигает сцену дальше?",
		DetailPromptEN:      "Which detail moves the scene forward?",
		DetailDistractors:   []tutorScenarioChoiceLine{tutorScenarioChoice("at the cafe", "Не эта сцена", "Wrong scene", "Это из другой сцены.", "This belongs to another scene.", "wrong"), tutorScenarioChoice("football", "Не связано с задачей", "Unrelated", "Это не помогает собеседнику.", "This does not help the other person.", "wrong"), tutorScenarioChoice("say nothing", "Пустая реакция", "Empty reaction", "Диалог остановится.", "The dialogue stops.", "wrong")},
		NextPromptRU:        "Какой следующий шаг звучит естественно?",
		NextPromptEN:        "Which next step sounds natural?",
		NextAnswer:          "Thank you. What is the next step?",
		NextDistractors:     []tutorScenarioChoiceLine{tutorScenarioChoice("Repeat the word list.", "Список слов", "Word list", "Это не реплика.", "This is not a reply.", "wrong"), tutorScenarioChoice("Close the app.", "Пустая реакция", "Empty reaction", "Собеседник ждёт ответ.", "The other person expects an answer.", "wrong"), tutorScenarioChoice("I need a coffee.", "Не по сцене", "Wrong scene", "Это другая ситуация.", "This is another situation.", "wrong")},
		MiniExplanationRU:   fmt.Sprintf("Порядок сцены: действие, затем %s, затем деталь %s. Паттерн: %s", item, detail, modelAnswer),
		MiniExplanationEN:   fmt.Sprintf("Scene order: action, then %s, then detail %s. Pattern: %s", item, detail, modelAnswer),
		WritingTaskRU:       fmt.Sprintf("Напишите одну рабочую реплику: используйте %s и деталь %s.", item, detail),
		WritingTaskEN:       fmt.Sprintf("Write one usable reply: use %s and the detail %s.", item, detail),
		WritingExpected:     []string{item, detail},
		ListeningText:       modelAnswer,
		ListeningQuestionRU: "Какие 2-3 детали важны в аудио?",
		ListeningQuestionEN: "Which 2-3 details matter in the audio?",
		ListeningExpected:   []string{item, detail},
		DialogueLines:       []string{lastTutorLine},
		DialoguePromptRU:    "Ответьте на последнюю реплику одной живой фразой.",
		DialoguePromptEN:    "Answer the last line with one live reply.",
		DialogueGoalRU:      fmt.Sprintf("Использовать %s и деталь %s.", item, detail),
		DialogueGoalEN:      fmt.Sprintf("Use %s and the detail %s.", item, detail),
		WritingVariants: []tutorScenarioVariantLine{
			tutorScenarioVariant("Базовый ответ", "Base answer", item+".", "Слишком коротко, но предмет понятен.", "Short, but the item is clear.", "A1", "минимум", "safe minimum", false),
			tutorScenarioVariant("Естественно", "Natural answer", modelAnswer, "Есть предмет и деталь.", "It has the item and detail.", "A2", "живой разговор", "daily conversation", false),
			tutorScenarioVariant("Сильный вариант", "Stronger answer", "Could you help me? "+modelAnswer, "Есть вежливое начало, предмет и деталь.", "It has a polite opener, item, and detail.", "B1", "лучший образец", "best practice model", false),
			tutorScenarioVariant("Не использовать", "Do not use", item+" / football.", "Это список слов, а не ответ человеку.", "This is a word list, not an answer to a person.", "avoid", "не повторять", "avoid", true),
		},
		DialogueVariants: []tutorScenarioVariantLine{
			tutorScenarioVariant("A1", "A1", modelAnswer, "Ответ закрывает последнюю реплику.", "The reply answers the last line.", "A1", "минимум", "safe minimum", false),
			tutorScenarioVariant("A2", "A2", "Yes. "+modelAnswer, "Звучит естественнее и остаётся коротко.", "It sounds more natural and stays short.", "A2", "живой разговор", "daily conversation", false),
			tutorScenarioVariant("B1", "B1", "Yes, thank you. "+modelAnswer, "Есть подтверждение, вежливость и нужная деталь.", "It has confirmation, politeness, and the needed detail.", "B1", "лучший образец", "best practice model", false),
			tutorScenarioVariant("Не использовать", "Do not use", item+" / football.", "Это список слов, а не ответ человеку.", "This is a word list, not an answer to a person.", "avoid", "не повторять", "avoid", true),
		},
	}
}

func tutorScenarioChoice(text string, labelRU string, labelEN string, whyRU string, whyEN string, quality string) tutorScenarioChoiceLine {
	return tutorScenarioChoiceLine{Text: text, LabelRU: labelRU, LabelEN: labelEN, WhyRU: whyRU, WhyEN: whyEN, Quality: quality}
}

func tutorScenarioVariant(labelRU string, labelEN string, text string, whyRU string, whyEN string, level string, useCaseRU string, useCaseEN string, avoid bool) tutorScenarioVariantLine {
	return tutorScenarioVariantLine{LabelRU: labelRU, LabelEN: labelEN, Text: text, WhyRU: whyRU, WhyEN: whyEN, Level: level, UseCaseRU: useCaseRU, UseCaseEN: useCaseEN, Avoid: avoid}
}

func buildTutorLesson(user userState) (tutorLesson, error) {
	return buildTutorLessonForSequence(user, tutorLessonSequence(user))
}

func buildTutorLessonForSequence(user userState, lessonSequence int) (tutorLesson, error) {
	language := normalizeLearningLanguage(user.LearningLanguage)
	interfaceLanguage := normalizeInterfaceLanguage(user.InterfaceLanguage)
	level := normalizeTutorLevel(user.Level)
	if lessonSequence < 0 {
		lessonSequence = 0
	}
	topicIndex := lessonSequence % len(tutorScenarioCourseTopics)
	functionIndex := (lessonSequence / len(tutorScenarioCourseTopics)) % len(tutorCourseFunctions)
	variantIndex := (lessonSequence / (len(tutorScenarioCourseTopics) * len(tutorCourseFunctions))) % len(tutorCourseVariants)
	topic := tutorScenarioCourseTopics[topicIndex]
	courseFunction := tutorCourseFunctions[functionIndex]
	variant := tutorCourseVariants[variantIndex]
	lessonNumber := tutorCourseLessonNumber(lessonSequence)
	words, err := selectTutorLessonWords(user, topic, level, 10, lessonSequence, variant.Code)
	if err != nil {
		return tutorLesson{}, err
	}
	if len(words) < 4 {
		return tutorLesson{}, fmt.Errorf("not enough local A1-A2 vocabulary for %s", language)
	}
	if len(words) > 8 {
		words = words[:8]
	}

	scenarioTemplate := tutorScenarioTemplateFor(topic, courseFunction, interfaceLanguage)
	review := selectTutorReviewWords(user, interfaceLanguage, language, 3)
	checks := tutorScenarioChecks(words, topic, courseFunction, interfaceLanguage)
	if language != "en" {
		checks = tutorTargetWordRecallChecks(words, interfaceLanguage)
	}
	choice := checks[0]
	languageName := learningLanguageByCode(language).ExplanationName
	topicTitle := tutorTopicTitle(topic, interfaceLanguage)
	first := words[0]
	second := words[tutorMinInt(1, len(words)-1)]
	listeningText, listeningQuestion, listeningExpected := tutorListeningBlock(words, topic, courseFunction, language, interfaceLanguage)
	dialogue, dialoguePrompt, dialogueGoal := tutorScenarioDialogue(words, topic, courseFunction, language, interfaceLanguage)
	title := tutorLocalized(interfaceLanguage, "AI Репетитор", "AI Tutor")

	if interfaceLanguage == "ru" {
		title = "AI Репетитор"
	}

	return tutorLesson{
		ID:                tutorLessonID(language, interfaceLanguage, level, lessonSequence, topic, courseFunction, variant),
		Title:             title,
		Level:             level,
		Topic:             topicTitle,
		VariantCode:       variant.Code,
		VariantTitle:      tutorVariantTitle(variant, interfaceLanguage),
		FocusSkill:        tutorVariantFocus(variant, interfaceLanguage),
		SuccessCriteria:   tutorSuccessCriteria(topic, courseFunction, variant, interfaceLanguage),
		LessonNumber:      lessonNumber,
		CourseSize:        tutorCourseSize,
		LearningLanguage:  language,
		InterfaceLanguage: interfaceLanguage,
		DurationMinutes:   8,
		Goal:              tutorLessonGoal(topic, courseFunction, interfaceLanguage) + " " + tutorVariantFocusLine(variant, interfaceLanguage),
		CanDo:             tutorCanDo(level, topic, courseFunction, interfaceLanguage),
		Scenario:          tutorTopicScenario(topic, interfaceLanguage) + " " + tutorVariantConstraint(variant, interfaceLanguage),
		ScenarioSlots:     tutorScenarioSlotsFromTemplate(scenarioTemplate, interfaceLanguage),
		TeachingPoint:     tutorTeachingPointForLanguage(scenarioTemplate, words, language, interfaceLanguage),
		FinalWordCheck:    buildTutorFinalWordCheck(words, tutorFinalWordTemplateForLanguage(scenarioTemplate, language), interfaceLanguage),
		TutorSummary:      tutorSummaryForLesson(words, scenarioTemplate, courseFunction, interfaceLanguage),
		Steps:             tutorLessonSteps(interfaceLanguage),
		Words:             words,
		GrammarTitle:      tutorScenarioGrammarTitle(level, courseFunction, interfaceLanguage),
		Grammar:           tutorScenarioGrammarBody(level, languageName, topic, courseFunction, interfaceLanguage),
		MiniExplanation:   tutorMiniExplanationForLanguage(level, topic, courseFunction, words, language, interfaceLanguage),
		Choice:            choice,
		Checks:            checks,
		WritingTask:       tutorWritingTaskForLanguage(first, second, topic, courseFunction, language, interfaceLanguage),
		WritingExpected:   tutorExpectedWordsForLanguage(words, scenarioTemplate, language),
		AnswerVariants:    tutorAnswerVariants(words, topic, courseFunction, language, interfaceLanguage, false),
		ListeningTask:     tutorListeningTask(topic, courseFunction, interfaceLanguage),
		ListeningText:     listeningText,
		ListeningQuestion: listeningQuestion,
		ListeningExpected: listeningExpected,
		PronunciationText: tutorPronunciationLine(words, topic, language),
		Dialogue:          dialogue,
		DialoguePrompt:    dialoguePrompt,
		DialogueGoal:      dialogueGoal,
		DialogueVariants:  tutorAnswerVariants(words, topic, courseFunction, language, interfaceLanguage, true),
		Review:            review,
		ReviewSummary:     tutorReviewSummary(words, topic, courseFunction, interfaceLanguage),
		SRS: tutorLessonSRS{
			Again: tutorLocalized(interfaceLanguage, "ошибка: повторить сегодня", "missed: repeat today"),
			Hard:  tutorLocalized(interfaceLanguage, "трудно: завтра", "hard: tomorrow"),
			Good:  tutorLocalized(interfaceLanguage, "верно: через 2 дня", "correct: in 2 days"),
			Easy:  tutorLocalized(interfaceLanguage, "легко: через неделю", "easy: in one week"),
		},
		NextActions: []string{
			tutorLocalized(interfaceLanguage, "Повторить новые слова в Review", "Review the new words"),
			tutorLocalized(interfaceLanguage, "Проиграть аудио ещё раз", "Replay the listening model"),
			tutorLocalized(interfaceLanguage, "Сохранить лучшую фразу в разговорник", "Save the best phrase to the phrasebook"),
		},
		Source: "local-a1-a2-course-bank",
	}, nil
}

func tutorLessonSteps(interfaceLanguage string) []tutorLessonStep {
	ru := normalizeInterfaceLanguage(interfaceLanguage) == "ru"
	if ru {
		return []tutorLessonStep{
			{Code: "words", Title: "Новые слова", Summary: "слова темы"},
			{Code: "explain", Title: "Мини-объяснение", Summary: "одно правило"},
			{Code: "choice", Title: "Выбор ответа", Summary: "узнать перевод"},
			{Code: "writing", Title: "Письмо", Summary: "собрать фразу"},
			{Code: "listening", Title: "Аудирование", Summary: "повторить вслух"},
			{Code: "dialogue", Title: "Мини-диалог", Summary: "контекст"},
			{Code: "review", Title: "SRS", Summary: "расписание памяти"},
		}
	}
	return []tutorLessonStep{
		{Code: "words", Title: "New words", Summary: "topic vocabulary"},
		{Code: "explain", Title: "Mini explanation", Summary: "one pattern"},
		{Code: "choice", Title: "Choice", Summary: "recognize meaning"},
		{Code: "writing", Title: "Writing", Summary: "build a sentence"},
		{Code: "listening", Title: "Listening", Summary: "catch details"},
		{Code: "pronunciation", Title: "Pronunciation", Summary: "repeat aloud"},
		{Code: "dialogue", Title: "Mini dialogue", Summary: "context"},
		{Code: "final-check", Title: "Final word check", Summary: "lesson words"},
		{Code: "review", Title: "SRS", Summary: "memory schedule"},
	}
}

func normalizeTutorLevel(level string) string {
	level = normalizeCEFRLevel(level)
	if level == "" {
		return "A1"
	}
	return level
}

func selectTutorLessonWords(user userState, topic tutorTopic, level string, count int, lessonSequence int, variantCode string) ([]tutorLessonWord, error) {
	language := normalizeLearningLanguage(user.LearningLanguage)
	interfaceLanguage := normalizeInterfaceLanguage(user.InterfaceLanguage)
	learned := map[string]bool{}
	for _, learnedWord := range user.LearnedWords {
		if normalizeLearningLanguage(learnedWord.Language) == language {
			learned[learnedWord.ID] = true
			learned[legacyVocabID(learnedWord.ID)] = true
		}
	}
	candidates := make([]vocabWord, 0, 64)
	if err := forEachVocabularyWord(language, func(word vocabWord) bool {
		if learned[word.ID] || learned[legacyVocabID(word.ID)] {
			return true
		}
		if !tutorAcceptWord(word, user, level) {
			return true
		}
		candidates = append(candidates, word)
		return len(candidates) < 1600
	}); err != nil {
		return nil, err
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		left := tutorWordScore(candidates[i], topic, level)
		right := tutorWordScore(candidates[j], topic, level)
		if left != right {
			return left > right
		}
		return tutorStableWordKey(candidates[i]) < tutorStableWordKey(candidates[j])
	})
	if len(candidates) > 1 {
		window := tutorMinInt(len(candidates), tutorMaxInt(48, count*8))
		offset := tutorStableIndex(fmt.Sprintf("%s:%s:%s:%s:%d", language, level, topic.Code, variantCode, lessonSequence), window)
		if offset > 0 {
			rotated := make([]vocabWord, 0, len(candidates))
			rotated = append(rotated, candidates[offset:window]...)
			rotated = append(rotated, candidates[:offset]...)
			rotated = append(rotated, candidates[window:]...)
			candidates = rotated
		}
	}
	words := make([]tutorLessonWord, 0, count)
	seen := map[string]bool{}
	for _, candidate := range candidates {
		item := tutorLessonWordFromVocab(candidate, interfaceLanguage, topic)
		if item.Word == "" || item.Translation == "" || seen[strings.ToLower(item.Word)] {
			continue
		}
		seen[strings.ToLower(item.Word)] = true
		words = append(words, item)
		if len(words) >= count {
			break
		}
	}
	return words, nil
}

func selectTutorReviewWords(user userState, interfaceLanguage string, language string, count int) []tutorLessonWord {
	review := make([]tutorLessonWord, 0, count)
	for _, learned := range user.LearnedWords {
		if len(review) >= count {
			break
		}
		if normalizeLearningLanguage(learned.Language) != language || learnedWordMastered(learned) {
			continue
		}
		word := learnedEntryVocabWord(learned)
		if !tutorLooksClean(word.English) || !tutorLooksClean(tutorTranslationForUser(word, interfaceLanguage)) {
			continue
		}
		item := tutorLessonWordFromVocab(word, interfaceLanguage, tutorCourseTopics[0])
		if item.Word != "" && item.Translation != "" {
			review = append(review, item)
		}
	}
	return review
}

func tutorAcceptWord(word vocabWord, user userState, level string) bool {
	wordLevel := normalizeCEFRLevel(word.Level)
	if wordLevel != "A1" && wordLevel != "A2" {
		return false
	}
	if cefrRank(wordLevel) > cefrRank(level)+1 {
		return false
	}
	if !tutorLooksClean(word.English) || !tutorLooksClean(tutorTranslationForUser(word, user.InterfaceLanguage)) {
		return false
	}
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" && normalizeLearningLanguage(word.Language) != "ru" && looksLikeEnglishOnly(tutorTranslationForUser(word, "ru")) {
		return false
	}
	if strings.Count(cleanDictionaryDisplay(word.English), " ") > 3 {
		return false
	}
	return true
}

func tutorWordScore(word vocabWord, topic tutorTopic, level string) int {
	score := 0
	if normalizeCEFRLevel(word.Level) == level {
		score += 50
	}
	if strings.EqualFold(word.Topic, "curated-core") || strings.HasPrefix(strings.ToLower(word.Source), "curated") {
		score += 25
	}
	if word.FrequencyRank > 0 {
		score += tutorMaxInt(0, 20-word.FrequencyRank/100)
	}
	needle := strings.ToLower(cleanDictionaryDisplay(word.English))
	for _, expected := range topic.Words {
		if needle == expected {
			score += 80
			break
		}
		if strings.Contains(needle, expected) || strings.Contains(strings.ToLower(word.Topic), expected) {
			score += 20
		}
	}
	return score
}

func tutorLessonWordFromVocab(word vocabWord, interfaceLanguage string, topic tutorTopic) tutorLessonWord {
	target := firstDictionaryValue(word.English)
	translation := tutorTranslationForUser(word, interfaceLanguage)
	context := wordContext(word, interfaceLanguage)
	if context == "" {
		context = firstDictionaryValue(word.Russian)
	}
	example, exampleTranslation := tutorScenarioExample(target, translation, interfaceLanguage, topic, normalizeLearningLanguage(word.Language))
	return tutorLessonWord{
		ID:           word.ID,
		Word:         target,
		Translation:  translation,
		Context:      context,
		Example:      example,
		ExampleRu:    exampleTranslation,
		PartOfSpeech: strings.TrimSpace(word.PartOfSpeech),
		Topic:        tutorTopicTitle(topic, interfaceLanguage),
		Level:        normalizeTutorLevel(word.Level),
		Frequency:    word.FrequencyRank,
		AudioText:    target,
	}
}

func tutorChoiceFromWords(words []tutorLessonWord, interfaceLanguage string) tutorLessonChoice {
	correct := words[0]
	options := make([]tutorChoiceOption, 0, tutorMinInt(4, len(words)))
	for index, word := range words {
		if index >= 4 {
			break
		}
		quality := "distractor"
		label := tutorLocalized(interfaceLanguage, "Похоже, но не здесь", "Close, but not here")
		why := tutorLocalized(interfaceLanguage, "Смотрите на слово в задании, а не на знакомый перевод рядом.", "Match the target word, not just a familiar meaning nearby.")
		if word.ID == correct.ID {
			quality = "correct"
			label = tutorLocalized(interfaceLanguage, "Точный перевод", "Exact meaning")
			why = tutorLocalized(interfaceLanguage, "Это значение нужно узнать перед тем, как собрать реплику.", "You need this meaning before building the line.")
		}
		options = append(options, tutorChoiceOption{ID: word.ID, Text: word.Translation, Label: label, Why: why, Skill: "meaning", Quality: quality, Avoid: word.ID != correct.ID})
	}
	return tutorLessonChoice{
		Prompt:          tutorLocalized(interfaceLanguage, "Выбери перевод: "+correct.Word, "Choose the translation: "+correct.Word),
		Options:         options,
		CorrectAnswerID: correct.ID,
	}
}

func tutorLessonSequence(user userState) int {
	sequence := user.WordLessonCount + user.LessonCount
	if sequence < 0 {
		return 0
	}
	return sequence
}

func tutorCourseLessonNumber(lessonSequence int) int {
	if lessonSequence < 0 {
		return 1
	}
	return lessonSequence%tutorCourseSize + 1
}

func tutorLessonID(language string, interfaceLanguage string, level string, lessonSequence int, topic tutorTopic, function tutorCourseFunction, variant tutorCourseVariant) string {
	return fmt.Sprintf(
		"tutor-%s-%s-%s-%04d-%s-%s-%s",
		normalizeLearningLanguage(language),
		normalizeInterfaceLanguage(interfaceLanguage),
		strings.ToLower(normalizeTutorLevel(level)),
		lessonSequence+1,
		topic.Code,
		function.Code,
		variant.Code,
	)
}

func tutorVariantTitle(variant tutorCourseVariant, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" && strings.TrimSpace(variant.TitleRU) != "" {
		return variant.TitleRU
	}
	return variant.TitleEN
}

func tutorVariantFocus(variant tutorCourseVariant, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" && strings.TrimSpace(variant.FocusRU) != "" {
		return variant.FocusRU
	}
	return variant.FocusEN
}

func tutorVariantFocusLine(variant tutorCourseVariant, interfaceLanguage string) string {
	focus := tutorVariantFocus(variant, interfaceLanguage)
	if focus == "" {
		return ""
	}
	return tutorLocalized(interfaceLanguage, "Focus: "+focus+".", "Focus: "+focus+".")
}

func tutorVariantConstraint(variant tutorCourseVariant, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" && strings.TrimSpace(variant.ConstraintRU) != "" {
		return variant.ConstraintRU
	}
	return variant.ConstraintEN
}

func tutorScenarioSlotsFromTemplate(template tutorScenarioTemplate, interfaceLanguage string) tutorScenarioSlots {
	return tutorScenarioSlots{
		Role:           template.Role,
		Situation:      tutorScenarioSituation(template, interfaceLanguage),
		RequiredAction: template.RequiredAction,
		Item:           template.Item,
		Detail:         template.Detail,
		Politeness:     append([]string{}, template.Politeness...),
		ModelAnswer:    template.ModelAnswer,
	}
}

func tutorScenarioExpectedWords(template tutorScenarioTemplate) []string {
	expected := make([]string, 0, 3)
	for _, item := range template.WritingExpected {
		item = strings.TrimSpace(item)
		if item != "" {
			expected = append(expected, item)
		}
	}
	if len(expected) == 0 {
		expected = append(expected, template.Item, template.Detail)
	}
	return expected
}

func tutorTeachingPointFor(template tutorScenarioTemplate, interfaceLanguage string) tutorTeachingPoint {
	item := strings.TrimSpace(template.Item)
	detail := strings.TrimSpace(template.Detail)
	model := strings.TrimSpace(template.ModelAnswer)
	order := []string{"action"}
	if item != "" {
		order = append(order, item)
	}
	if detail != "" {
		order = append(order, detail)
	}
	explanation := tutorLocalized(interfaceLanguage, template.MiniExplanationRU, template.MiniExplanationEN)
	if strings.TrimSpace(explanation) == "" {
		explanation = fmt.Sprintf("Scene order: %s. Pattern: %s", strings.Join(order, " -> "), model)
	}
	return tutorTeachingPoint{
		Title:       tutorLocalized(interfaceLanguage, "One useful pattern", "One useful pattern"),
		Pattern:     model,
		SceneOrder:  order,
		Explanation: explanation,
		ModelAnswer: model,
	}
}

func tutorTeachingPointForLanguage(template tutorScenarioTemplate, words []tutorLessonWord, language string, interfaceLanguage string) tutorTeachingPoint {
	if normalizeLearningLanguage(language) == "en" {
		return tutorTeachingPointFor(template, interfaceLanguage)
	}
	model := tutorTargetPracticeLine(words, 4)
	return tutorTeachingPoint{
		Title:       tutorLocalized(interfaceLanguage, "Target-language pattern", "Target-language pattern"),
		Pattern:     model,
		SceneOrder:  tutorTargetSceneOrder(words, 3),
		Explanation: tutorTargetPracticeExplanation(words, interfaceLanguage),
		ModelAnswer: model,
	}
}

func buildTutorFinalWordCheck(words []tutorLessonWord, template tutorScenarioTemplate, interfaceLanguage string) tutorFinalWordCheck {
	items := tutorWordCheckItems(words, template, interfaceLanguage)
	required := tutorMinInt(len(items), tutorMaxInt(2, len(items)-1))
	return tutorFinalWordCheck{
		Prompt:          tutorLocalized(interfaceLanguage, "Final word check: prove you still remember the lesson words.", "Final word check: prove you still remember the lesson words."),
		Items:           items,
		RequiredCorrect: required,
		SummaryPass:     tutorLocalized(interfaceLanguage, "Word check passed. Now schedule the review.", "Word check passed. Now schedule the review."),
		SummaryRetry:    tutorLocalized(interfaceLanguage, "Repeat the weak words before finishing the lesson.", "Repeat the weak words before finishing the lesson."),
	}
}

func tutorFinalWordTemplateForLanguage(template tutorScenarioTemplate, language string) tutorScenarioTemplate {
	if normalizeLearningLanguage(language) == "en" {
		return template
	}
	return tutorScenarioTemplate{}
}

func tutorWordCheckItems(words []tutorLessonWord, template tutorScenarioTemplate, interfaceLanguage string) []tutorLessonChoice {
	items := make([]tutorLessonChoice, 0, 4)
	seen := map[string]bool{}
	addChoice := func(choice tutorLessonChoice) {
		key := strings.ToLower(strings.TrimSpace(choice.CorrectAnswerID))
		if key == "" {
			key = strings.ToLower(strings.TrimSpace(choice.Prompt))
		}
		if key == "" || seen[key] || len(choice.Options) == 0 {
			return
		}
		seen[key] = true
		items = append(items, choice)
	}
	addSlot := func(slot string) {
		slot = strings.TrimSpace(slot)
		if slot == "" {
			return
		}
		for index, word := range words {
			if strings.EqualFold(strings.TrimSpace(word.Word), slot) {
				addChoice(tutorTargetWordRecallChoice(word, words, index, interfaceLanguage))
				return
			}
		}
		addChoice(tutorSlotWordChoice(slot, words, interfaceLanguage))
	}
	addSlot(template.Item)
	addSlot(template.Detail)
	for index, word := range words {
		if strings.TrimSpace(word.ID) == "" || strings.TrimSpace(word.Word) == "" {
			continue
		}
		addChoice(tutorTargetWordRecallChoice(word, words, index, interfaceLanguage))
		if len(items) >= 4 {
			break
		}
	}
	return items
}

func tutorTargetWordRecallChecks(words []tutorLessonWord, interfaceLanguage string) []tutorLessonChoice {
	items := tutorWordCheckItems(words, tutorScenarioTemplate{}, interfaceLanguage)
	if len(items) == 0 {
		return []tutorLessonChoice{tutorChoiceFromWords(words, interfaceLanguage)}
	}
	return items
}

func tutorTargetWordRecallChoice(correct tutorLessonWord, words []tutorLessonWord, wordIndex int, interfaceLanguage string) tutorLessonChoice {
	options := make([]tutorChoiceOption, 0, 4)
	seen := map[string]bool{}
	addOption := func(word tutorLessonWord, avoid bool) {
		text := strings.TrimSpace(word.Word)
		id := strings.TrimSpace(word.ID)
		if text == "" || id == "" {
			return
		}
		key := strings.ToLower(text)
		if seen[key] {
			return
		}
		seen[key] = true
		options = append(options, tutorChoiceOption{
			ID:      id,
			Text:    text,
			Skill:   "final-word",
			Quality: map[bool]string{true: "distractor", false: "correct"}[avoid],
			Avoid:   avoid,
		})
	}
	addOption(correct, false)
	for offset := 1; offset <= len(words) && len(options) < 4; offset++ {
		candidate := words[(wordIndex+offset)%len(words)]
		if candidate.ID == correct.ID || strings.EqualFold(candidate.Word, correct.Word) {
			continue
		}
		addOption(candidate, true)
	}
	for _, candidate := range words {
		if len(options) >= 4 {
			break
		}
		if candidate.ID == correct.ID || strings.EqualFold(candidate.Word, correct.Word) {
			continue
		}
		addOption(candidate, true)
	}
	return tutorLessonChoice{
		Prompt:          tutorRecallCue(correct, wordIndex),
		Options:         options,
		CorrectAnswerID: correct.ID,
		Feedback:        tutorLocalized(interfaceLanguage, "Word recalled.", "Word recalled."),
	}
}

func tutorRecallCue(word tutorLessonWord, index int) string {
	target := strings.ToLower(strings.TrimSpace(word.Word))
	for _, candidate := range []string{word.Translation, word.Context, word.ExampleRu, word.PartOfSpeech} {
		cue := strings.TrimSpace(candidate)
		if cue == "" {
			continue
		}
		if target != "" && strings.Contains(strings.ToLower(cue), target) {
			continue
		}
		return cue
	}
	return fmt.Sprintf("Meaning %d", index+1)
}

func tutorSlotWordChoice(slot string, words []tutorLessonWord, interfaceLanguage string) tutorLessonChoice {
	options := []tutorChoiceOption{{
		ID:      "slot:" + strings.ToLower(strings.ReplaceAll(slot, " ", "-")),
		Text:    slot,
		Label:   tutorLocalized(interfaceLanguage, "Scene word", "Scene word"),
		Why:     tutorLocalized(interfaceLanguage, "This word is required by the scenario pattern.", "This word is required by the scenario pattern."),
		Skill:   "final-word",
		Quality: "correct",
	}}
	seen := map[string]bool{strings.ToLower(slot): true}
	for index, word := range words {
		text := strings.TrimSpace(word.Word)
		key := strings.ToLower(text)
		if text == "" || seen[key] {
			continue
		}
		seen[key] = true
		options = append(options, tutorChoiceOption{
			ID:      fmt.Sprintf("slot:wrong:%d", index),
			Text:    text,
			Label:   tutorWeakOptionLabel(index, interfaceLanguage),
			Why:     tutorWeakOptionWhy(index, interfaceLanguage),
			Skill:   "final-word",
			Quality: "weak",
			Avoid:   true,
		})
		if len(options) >= 4 {
			break
		}
	}
	return tutorLessonChoice{
		Prompt:          tutorLocalized(interfaceLanguage, "Выберите обязательное слово сцены.", "Choose the required scenario word."),
		Options:         options,
		CorrectAnswerID: options[0].ID,
		Feedback:        tutorLocalized(interfaceLanguage, "Correct: this word belongs to the lesson pattern.", "Correct: this word belongs to the lesson pattern."),
	}
}

func tutorSummaryForLesson(words []tutorLessonWord, template tutorScenarioTemplate, function tutorCourseFunction, interfaceLanguage string) tutorSummaryBlock {
	strong := make([]string, 0, 4)
	seen := map[string]bool{}
	add := func(value string) {
		value = strings.TrimSpace(value)
		key := strings.ToLower(value)
		if value == "" || seen[key] {
			return
		}
		seen[key] = true
		strong = append(strong, value)
	}
	add(template.Item)
	add(template.Detail)
	for _, word := range words {
		if len(strong) >= 4 {
			break
		}
		add(word.Word)
	}
	return tutorSummaryBlock{
		CanSay:      strings.TrimSpace(template.ModelAnswer),
		StrongItems: strong,
		WeakItems:   []string{},
		NextReview:  tutorLocalized(interfaceLanguage, "Choose a review button based on how hard the final check felt.", "Choose a review button based on how hard the final check felt."),
	}
}

func tutorAnswerSlotFeedback(answer string, template tutorScenarioTemplate, interfaceLanguage string) string {
	normalized := strings.ToLower(answer)
	present := make([]string, 0, 3)
	missing := make([]string, 0, 3)
	for _, slot := range []string{template.Item, template.Detail} {
		slot = strings.TrimSpace(slot)
		if slot == "" {
			continue
		}
		if strings.Contains(normalized, strings.ToLower(slot)) {
			present = append(present, slot)
		} else {
			missing = append(missing, slot)
		}
	}
	if len(missing) == 0 {
		return tutorLocalized(interfaceLanguage, "Good: your answer has the required scene words.", "Good: your answer has the required scene words.")
	}
	return tutorLocalized(interfaceLanguage, "Good parts: ", "Good parts: ") + strings.Join(present, ", ") + ". " +
		tutorLocalized(interfaceLanguage, "Add: ", "Add: ") + strings.Join(missing, ", ") + "."
}

func tutorScenarioSituation(template tutorScenarioTemplate, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" && strings.TrimSpace(template.SituationRU) != "" {
		return template.SituationRU
	}
	return template.SituationEN
}

func tutorScenarioChoiceOption(prefix string, index int, line tutorScenarioChoiceLine, interfaceLanguage string) tutorChoiceOption {
	ru := normalizeInterfaceLanguage(interfaceLanguage) == "ru"
	label := line.LabelEN
	why := line.WhyEN
	if ru {
		label = line.LabelRU
		why = line.WhyRU
	}
	if label == "" {
		label = tutorWeakOptionLabel(index, interfaceLanguage)
	}
	if why == "" {
		why = tutorWeakOptionWhy(index, interfaceLanguage)
	}
	quality := line.Quality
	if quality == "" {
		quality = "weak"
	}
	return tutorChoiceOption{
		ID:      fmt.Sprintf("%s:wrong:%d", prefix, index),
		Text:    line.Text,
		Label:   label,
		Why:     why,
		Skill:   prefix,
		Quality: quality,
		Avoid:   true,
	}
}

func tutorScenarioVariantToAnswer(prefix string, index int, variant tutorScenarioVariantLine, interfaceLanguage string) tutorAnswerVariant {
	ru := normalizeInterfaceLanguage(interfaceLanguage) == "ru"
	label := variant.LabelEN
	why := variant.WhyEN
	useCase := variant.UseCaseEN
	if ru {
		label = variant.LabelRU
		why = variant.WhyRU
		useCase = variant.UseCaseRU
	}
	return tutorAnswerVariant{
		ID:      fmt.Sprintf("%s:%d", prefix, index+1),
		Label:   label,
		Text:    variant.Text,
		Why:     why,
		Level:   variant.Level,
		UseCase: useCase,
		Avoid:   variant.Avoid,
	}
}

func tutorSuccessCriteria(topic tutorTopic, function tutorCourseFunction, variant tutorCourseVariant, interfaceLanguage string) []string {
	template := tutorScenarioTemplateFor(topic, function, interfaceLanguage)
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return []string{
			"В ответе должно быть: просьба или действие.",
			"Предмет: " + template.Item + ".",
			"Деталь: " + template.Detail + ".",
		}
	}
	return []string{
		"Answer checklist: request or action.",
		"Item: " + template.Item + ".",
		"Detail: " + template.Detail + ".",
	}
}

func tutorLessonStorageContext(user userState) (string, string, string) {
	return normalizeLearningLanguage(user.LearningLanguage), normalizeInterfaceLanguage(user.InterfaceLanguage), normalizeTutorLevel(user.Level)
}

func tutorLessonMatchesContext(lesson tutorLesson, language string, interfaceLanguage string, level string) bool {
	return normalizeLearningLanguage(lesson.LearningLanguage) == normalizeLearningLanguage(language) &&
		normalizeInterfaceLanguage(lesson.InterfaceLanguage) == normalizeInterfaceLanguage(interfaceLanguage) &&
		normalizeTutorLevel(lesson.Level) == normalizeTutorLevel(level)
}

func tutorReusableLessonUser(user userState) userState {
	core := user
	core.TelegramID = 0
	core.LearnedWords = nil
	core.Mistakes = nil
	core.Phrasebook = nil
	core.LessonCount = 0
	core.WordLessonCount = 0
	return core
}

func tutorPersonalizeLessonForUser(lesson tutorLesson, user userState) tutorLesson {
	language := normalizeLearningLanguage(lesson.LearningLanguage)
	if language == "" {
		language = normalizeLearningLanguage(user.LearningLanguage)
	}
	interfaceLanguage := normalizeInterfaceLanguage(lesson.InterfaceLanguage)
	if interfaceLanguage == "" {
		interfaceLanguage = normalizeInterfaceLanguage(user.InterfaceLanguage)
	}
	review := selectTutorReviewWords(user, interfaceLanguage, language, 3)
	if len(review) > 0 {
		lesson.Review = review
	}
	return lesson
}

func tutorLessonGenerationAttemptLimit() int {
	return tutorMaxInt(64, tutorCourseSize)
}

func tutorLessonGoal(topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("Сценарная задача: %s в теме «%s».", strings.ToLower(function.TitleRU), tutorTopicTitle(topic, interfaceLanguage))
	}
	return fmt.Sprintf("Scenario task: %s in %s.", function.TitleEN, tutorTopicTitle(topic, interfaceLanguage))
}

func tutorCanDo(level string, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("После урока вы сможете на уровне %s %s.", level, function.CanDoRU)
	}
	return fmt.Sprintf("After this %s lesson you can %s.", level, function.CanDoEN)
}

func tutorTopicScenario(topic tutorTopic, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" && topic.ScenarioRU != "" {
		return topic.ScenarioRU
	}
	if topic.ScenarioEN != "" {
		return topic.ScenarioEN
	}
	return tutorTopicTitle(topic, interfaceLanguage)
}

func tutorScenarioChecks(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) []tutorLessonChoice {
	if len(words) < 4 {
		return []tutorLessonChoice{tutorChoiceFromWords(words, interfaceLanguage)}
	}
	template := tutorScenarioTemplateFor(topic, function, interfaceLanguage)
	promptBestLine := tutorLocalized(interfaceLanguage, "Какая реплика лучше подходит к сценарию?", "Which line best fits the scenario?")
	promptDetail := tutorLocalized(interfaceLanguage, template.DetailPromptRU, template.DetailPromptEN)
	promptNext := tutorLocalized(interfaceLanguage, template.NextPromptRU, template.NextPromptEN)
	feedback := tutorLocalized(interfaceLanguage, "Верно: это реплика из сценария, а не просто перевод слова.", "Correct: this is a scenario line, not just a word translation.")

	return []tutorLessonChoice{
		{
			Prompt: promptBestLine,
			Options: tutorScenarioCheckOptions("scenario", tutorChoiceOption{
				ID:      "scenario:correct",
				Text:    template.ModelAnswer,
				Label:   tutorLocalized(interfaceLanguage, "Рабочая реплика", "Usable line"),
				Why:     tutorLocalized(interfaceLanguage, "Есть просьба, предмет и деталь. Такую фразу можно сказать в реальном диалоге.", "It has a request, item, and detail. You can say it in a real dialogue."),
				Skill:   "scenario",
				Quality: "correct",
			}, []tutorScenarioChoiceLine{
				tutorScenarioChoice(template.ShortAnswer, "Слишком коротко", "Too short", "Нет вежливости и детали.", "It lacks politeness and detail.", "weak"),
				tutorScenarioChoice(template.WrongSceneAnswer, "Не по сцене", "Wrong scene", "Это не решает текущую ситуацию.", "This does not solve the current situation.", "weak"),
				tutorScenarioChoice(template.IncompleteAnswer, "Похоже, но неполно", "Close, but incomplete", "Хорошо, но ещё нет обязательной детали.", "Good, but the required detail is missing.", "weak"),
			}, interfaceLanguage),
			CorrectAnswerID: "scenario:correct",
			Feedback:        feedback,
		},
		{
			Prompt: promptDetail,
			Options: tutorScenarioCheckOptions("detail", tutorChoiceOption{
				ID:      "detail:correct",
				Text:    template.Detail,
				Label:   tutorLocalized(interfaceLanguage, "Верно", "Correct"),
				Why:     tutorLocalized(interfaceLanguage, "Эта деталь уточняет реплику и помогает собеседнику ответить.", "This detail clarifies the line and helps the other person respond."),
				Skill:   "detail",
				Quality: "correct",
			}, template.DetailDistractors, interfaceLanguage),
			CorrectAnswerID: "detail:correct",
			Feedback:        tutorLocalized(interfaceLanguage, "Верно: вы выбрали деталь, которая двигает диалог дальше.", "Correct: you picked the detail that moves the dialogue forward."),
		},
		{
			Prompt: promptNext,
			Options: tutorScenarioCheckOptions("next", tutorChoiceOption{
				ID:      "next:correct",
				Text:    template.NextAnswer,
				Label:   tutorLocalized(interfaceLanguage, "Рабочая реплика", "Usable line"),
				Why:     tutorLocalized(interfaceLanguage, "Это естественно закрывает следующий шаг сцены.", "This naturally closes the next step of the scene."),
				Skill:   "next",
				Quality: "correct",
			}, template.NextDistractors, interfaceLanguage),
			CorrectAnswerID: "next:correct",
			Feedback:        tutorLocalized(interfaceLanguage, "Верно: это закрывает шаг и готовит ваш ответ.", "Correct: this closes the step and prepares your answer."),
		},
	}
}

func tutorScenarioCheckOptions(prefix string, correct tutorChoiceOption, distractors []tutorScenarioChoiceLine, interfaceLanguage string) []tutorChoiceOption {
	options := []tutorChoiceOption{correct}
	seen := map[string]bool{strings.ToLower(strings.TrimSpace(correct.Text)): true}
	for index, distractor := range distractors {
		key := strings.ToLower(strings.TrimSpace(distractor.Text))
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		options = append(options, tutorScenarioChoiceOption(prefix, index, distractor, interfaceLanguage))
		if len(options) >= 4 {
			break
		}
	}
	return options
}

func tutorCheckOptions(prefix string, correct string, distractors []string, interfaceLanguage string) []tutorChoiceOption {
	options := []tutorChoiceOption{{
		ID:      prefix + ":correct",
		Text:    correct,
		Label:   tutorLocalized(interfaceLanguage, "Рабочая реплика", "Usable line"),
		Why:     tutorLocalized(interfaceLanguage, "Есть действие, деталь и понятная цель. Такую фразу можно сказать в реальном диалоге.", "It has an action, a detail, and a clear goal. You can say it in a real dialogue."),
		Skill:   prefix,
		Quality: "correct",
	}}
	for index, distractor := range distractors {
		distractor = strings.TrimSpace(distractor)
		if distractor == "" || strings.EqualFold(distractor, correct) {
			continue
		}
		options = append(options, tutorChoiceOption{
			ID:      fmt.Sprintf("%s:wrong:%d", prefix, index),
			Text:    distractor,
			Label:   tutorWeakOptionLabel(index, interfaceLanguage),
			Why:     tutorWeakOptionWhy(index, interfaceLanguage),
			Skill:   prefix,
			Quality: "weak",
			Avoid:   true,
		})
		if len(options) >= 4 {
			break
		}
	}
	return options
}

func tutorWeakOptionLabel(index int, interfaceLanguage string) string {
	labelsRU := []string{"Просто перевод", "Повтор слов", "Пустая реакция"}
	labelsEN := []string{"Translation only", "Word repeat", "Empty reaction"}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return labelsRU[tutorMinInt(index, len(labelsRU)-1)]
	}
	return labelsEN[tutorMinInt(index, len(labelsEN)-1)]
}

func tutorWeakOptionWhy(index int, interfaceLanguage string) string {
	whyRU := []string{
		"Перевод слова не решает задачу собеседника.",
		"Повторить слова легче, но это не звучит как ответ.",
		"В реальной сцене собеседник ждёт действие или деталь.",
	}
	whyEN := []string{
		"A word translation does not solve the other person's task.",
		"Repeating words is easier, but it does not sound like an answer.",
		"In a real scene, the other person expects an action or detail.",
	}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return whyRU[tutorMinInt(index, len(whyRU)-1)]
	}
	return whyEN[tutorMinInt(index, len(whyEN)-1)]
}

func tutorCheckLine(topic tutorTopic, first tutorLessonWord, second tutorLessonWord, ru bool) string {
	if ru {
		switch topic.Code {
		case "hotel":
			return fmt.Sprintf("I have a reservation. Can I have a %s?", second.Word)
		case "doctor":
			return fmt.Sprintf("I have %s. Can I see a %s today?", first.Word, second.Word)
		case "food":
			return fmt.Sprintf("Can I see the %s and have %s?", first.Word, second.Word)
		case "travel":
			return fmt.Sprintf("I need a %s for tomorrow. What time is it?", first.Word)
		}
	}
	switch topic.Code {
	case "hotel":
		return fmt.Sprintf("I have a reservation. Can I have a %s?", second.Word)
	case "doctor":
		return fmt.Sprintf("I have %s. Can I see a %s today?", first.Word, second.Word)
	case "food":
		return fmt.Sprintf("Can I see the %s and have %s?", first.Word, second.Word)
	case "travel":
		return fmt.Sprintf("I need a %s for tomorrow. What time is it?", first.Word)
	}
	return fmt.Sprintf("Can you help me with %s and %s?", first.Word, second.Word)
}

func tutorCheckDetailLine(topic tutorTopic, second tutorLessonWord, third tutorLessonWord, ru bool) string {
	switch topic.Code {
	case "hotel", "services":
		return fmt.Sprintf("%s: %s", third.Word, third.Translation)
	case "work", "study", "travel":
		return fmt.Sprintf("%s: %s", second.Word, second.Translation)
	case "doctor":
		return fmt.Sprintf("%s: %s", second.Word, second.Translation)
	}
	return fmt.Sprintf("%s: %s", second.Word, second.Translation)
}

func tutorCheckNextLine(topic tutorTopic, third tutorLessonWord, fourth tutorLessonWord, ru bool) string {
	switch topic.Code {
	case "hotel":
		return fmt.Sprintf("Here is my %s. Thank you.", third.Word)
	case "doctor":
		return fmt.Sprintf("Can I come this %s?", fourth.Word)
	case "food":
		return fmt.Sprintf("Thank you. Can I pay by %s?", fourth.Word)
	case "travel":
		return fmt.Sprintf("Great. I will take this %s.", third.Word)
	}
	return fmt.Sprintf("Thank you. The next step is %s.", fourth.Word)
}

func tutorDistractorLines(topic tutorTopic, first tutorLessonWord, second tutorLessonWord, ru bool) []string {
	return []string{
		fmt.Sprintf("%s / %s", first.Translation, second.Translation),
		fmt.Sprintf("I repeat: %s, %s.", first.Word, second.Word),
		"Yes. No. Maybe.",
	}
}

func tutorNextDistractors(topic tutorTopic, ru bool) []string {
	if ru {
		return []string{"Повторить список слов без ответа", "Выбрать случайный перевод", "Закрыть приложение"}
	}
	return []string{"Repeat the word list only", "Pick a random translation", "Close the app"}
}

func tutorScenarioGrammarTitle(level string, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("Паттерн %s: %s", level, strings.ToLower(function.TitleRU))
	}
	return fmt.Sprintf("%s pattern: %s", level, function.TitleEN)
}

func tutorScenarioGrammarBody(level string, languageName string, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		switch level {
		case "A1":
			return "Соберите короткую реплику: вежливое начало + ключевое слово + одна деталь. В этом уроке важна не форма слова сама по себе, а то, какую задачу она решает в сцене."
		case "A2":
			return "Добавьте уточнение: время, место, имя или документ. Ответ должен звучать как реплика в диалоге, а не как список слов."
		default:
			return "Сначала выберите коммуникативную цель, затем переформулируйте мысль естественно: причина, деталь, следующий шаг. Это тренирует контекст и нюанс, а не механический перевод."
		}
	}
	switch level {
	case "A1":
		return "Build a short " + languageName + " line: polite opener + key word + one detail. The word matters because it solves the scene task."
	case "A2":
		return "Add one clarification: time, place, name, or document. Make it a dialogue line, not a word list."
	default:
		return "Start from the communicative goal, then rephrase naturally: reason, detail, and next step."
	}
}

func tutorMiniExplanation(level string, topic tutorTopic, function tutorCourseFunction, words []tutorLessonWord, interfaceLanguage string) string {
	template := tutorScenarioTemplateFor(topic, function, interfaceLanguage)
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return template.MiniExplanationRU
	}
	return template.MiniExplanationEN
}

func tutorMiniExplanationForLanguage(level string, topic tutorTopic, function tutorCourseFunction, words []tutorLessonWord, language string, interfaceLanguage string) string {
	if normalizeLearningLanguage(language) == "en" {
		return tutorMiniExplanation(level, topic, function, words, interfaceLanguage)
	}
	return tutorTargetPracticeExplanation(words, interfaceLanguage)
}

func tutorWritingTask(first tutorLessonWord, second tutorLessonWord, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	template := tutorScenarioTemplateFor(topic, function, interfaceLanguage)
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return template.WritingTaskRU
	}
	return template.WritingTaskEN
}

func tutorWritingTaskForLanguage(first tutorLessonWord, second tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string) string {
	if normalizeLearningLanguage(language) == "en" {
		return tutorWritingTask(first, second, topic, function, interfaceLanguage)
	}
	targets := strings.Join(tutorLessonWordTexts([]tutorLessonWord{first, second}, 2), ", ")
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return "Напишите одну короткую фразу на изучаемом языке, используя: " + targets + "."
	}
	return "Write one short target-language line using: " + targets + "."
}

func tutorAnswerVariants(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string, dialogue bool) []tutorAnswerVariant {
	if normalizeLearningLanguage(language) != "en" {
		return tutorTargetAnswerVariants(words, topic, language, interfaceLanguage, dialogue)
	}
	template := tutorScenarioTemplateFor(topic, function, interfaceLanguage)
	lines := template.WritingVariants
	prefix := "writing"
	if dialogue {
		lines = template.DialogueVariants
		prefix = "dialogue"
	}
	variants := make([]tutorAnswerVariant, 0, len(lines))
	seen := map[string]bool{}
	for index, line := range lines {
		line.Text = strings.TrimSpace(line.Text)
		key := strings.ToLower(line.Text)
		if line.Text == "" || seen[key] {
			continue
		}
		seen[key] = true
		variants = append(variants, tutorScenarioVariantToAnswer(prefix, index, line, interfaceLanguage))
	}
	return variants
}

func tutorTargetAnswerVariants(words []tutorLessonWord, topic tutorTopic, language string, interfaceLanguage string, dialogue bool) []tutorAnswerVariant {
	first := tutorWordOrFallback(words, 0, "word-1")
	second := tutorWordOrFallback(words, 1, "word-2")
	third := tutorWordOrFallback(words, 2, "word-3")
	fourth := tutorWordOrFallback(words, 3, "word-4")
	lines := tutorVariantLines(topic, language, first, second, third, fourth, dialogue)
	labels := []string{
		tutorLocalized(interfaceLanguage, "Base answer", "Base answer"),
		tutorLocalized(interfaceLanguage, "Natural answer", "Natural answer"),
		tutorLocalized(interfaceLanguage, "Stronger answer", "Stronger answer"),
		tutorLocalized(interfaceLanguage, "Do not use", "Do not use"),
	}
	whys := []string{
		tutorLocalized(interfaceLanguage, "Uses the first target words.", "Uses the first target words."),
		tutorLocalized(interfaceLanguage, "Adds one more target word.", "Adds one more target word."),
		tutorLocalized(interfaceLanguage, "Uses more lesson words for a fuller answer.", "Uses more lesson words for a fuller answer."),
		tutorLocalized(interfaceLanguage, "This is a word list, not a natural answer.", "This is a word list, not a natural answer."),
	}
	levels := []string{"A1", "A2", "B1", "avoid"}
	prefix := "writing"
	if dialogue {
		prefix = "dialogue"
	}
	variants := make([]tutorAnswerVariant, 0, len(lines))
	for index, text := range lines {
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		variants = append(variants, tutorAnswerVariant{
			ID:      fmt.Sprintf("%s:%d", prefix, index+1),
			Label:   labels[tutorMinInt(index, len(labels)-1)],
			Text:    text,
			Why:     whys[tutorMinInt(index, len(whys)-1)],
			Level:   levels[tutorMinInt(index, len(levels)-1)],
			UseCase: tutorLocalized(interfaceLanguage, "target words", "target words"),
			Avoid:   index == len(lines)-1,
		})
	}
	return variants
}

func tutorVariantLines(topic tutorTopic, language string, first tutorLessonWord, second tutorLessonWord, third tutorLessonWord, fourth tutorLessonWord, dialogue bool) []string {
	if normalizeLearningLanguage(language) != "en" {
		base := strings.TrimSpace(first.Word + " " + second.Word)
		natural := strings.TrimSpace(first.Word + " " + second.Word + " " + third.Word)
		strong := strings.TrimSpace(first.Word + " " + second.Word + " " + third.Word + " " + fourth.Word)
		weak := strings.TrimSpace(first.Word + " / " + second.Word)
		return []string{base, natural, strong, weak}
	}
	if dialogue {
		switch topic.Code {
		case "hotel":
			return []string{
				fmt.Sprintf("Yes, here is my %s.", third.Word),
				fmt.Sprintf("Yes, here is my %s. The %s is under my name.", third.Word, first.Word),
				fmt.Sprintf("Yes, here is my %s. The %s is under my name, and I need a quiet %s.", third.Word, first.Word, second.Word),
				fmt.Sprintf("%s. %s. Yes.", first.Word, second.Word),
			}
		case "doctor":
			return []string{
				fmt.Sprintf("Yes, I can come %s.", fourth.Word),
				fmt.Sprintf("Yes, I can come %s. I have %s.", fourth.Word, first.Word),
				fmt.Sprintf("Yes, I can come %s. I have %s and need a %s.", fourth.Word, first.Word, second.Word),
				fmt.Sprintf("%s, %s.", first.Word, second.Word),
			}
		case "food":
			return []string{
				fmt.Sprintf("%s, please.", second.Word),
				fmt.Sprintf("I would like %s, please.", second.Word),
				fmt.Sprintf("I would like %s, please. Could I also see the %s?", second.Word, first.Word),
				fmt.Sprintf("%s / %s.", first.Word, second.Word),
			}
		case "travel":
			return []string{
				fmt.Sprintf("I need a %s.", first.Word),
				fmt.Sprintf("I need a %s for %s.", first.Word, fourth.Word),
				fmt.Sprintf("I need a %s for %s. What time does it leave?", first.Word, fourth.Word),
				fmt.Sprintf("%s. Time?", first.Word),
			}
		}
		return []string{
			fmt.Sprintf("Yes, the detail is %s.", second.Word),
			fmt.Sprintf("Yes, the detail is %s. I need help with %s.", second.Word, first.Word),
			fmt.Sprintf("Yes, the detail is %s. Could you help me with the next step for %s?", second.Word, first.Word),
			fmt.Sprintf("%s, %s, okay.", first.Word, second.Word),
		}
	}
	switch topic.Code {
	case "hotel":
		return []string{
			fmt.Sprintf("I have a %s.", first.Word),
			fmt.Sprintf("I have a %s and need a quiet %s.", first.Word, second.Word),
			fmt.Sprintf("Good evening. I have a %s and need a quiet %s for %s.", first.Word, second.Word, fourth.Word),
			fmt.Sprintf("%s, %s, %s.", first.Word, second.Word, third.Word),
		}
	case "doctor":
		return []string{
			fmt.Sprintf("I have %s.", first.Word),
			fmt.Sprintf("I have %s and need a %s.", first.Word, second.Word),
			fmt.Sprintf("Hello. I have %s and need a %s %s.", first.Word, second.Word, fourth.Word),
			fmt.Sprintf("%s. %s. Help.", first.Word, second.Word),
		}
	case "food":
		return []string{
			fmt.Sprintf("Can I see the %s?", first.Word),
			fmt.Sprintf("Can I see the %s and have %s?", first.Word, second.Word),
			fmt.Sprintf("Hello. Can I see the %s and have %s, please?", first.Word, second.Word),
			fmt.Sprintf("%s / %s / please.", first.Word, second.Word),
		}
	case "travel":
		return []string{
			fmt.Sprintf("I need a %s.", first.Word),
			fmt.Sprintf("I need a %s for %s.", first.Word, fourth.Word),
			fmt.Sprintf("Hello. I need a %s for %s. What time does it leave?", first.Word, fourth.Word),
			fmt.Sprintf("%s. %s.", first.Word, fourth.Word),
		}
	}
	return []string{
		fmt.Sprintf("Can you help me with %s?", first.Word),
		fmt.Sprintf("Can you help me with %s? The detail is %s.", first.Word, second.Word),
		fmt.Sprintf("Could you help me with %s? The detail is %s, and the next step is %s.", first.Word, second.Word, third.Word),
		fmt.Sprintf("%s, %s, maybe.", first.Word, second.Word),
	}
}

func tutorListeningTask(topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return "Прослушайте модель и напишите, какие 2-3 ключевые детали вы услышали. Репетитор проверит, есть ли смысловые слова."
	}
	return "Listen to the model and type 2-3 key details you heard. The tutor checks the meaning words."
}

func tutorListeningBlock(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string) (string, string, []string) {
	if normalizeLearningLanguage(language) != "en" {
		text := tutorTargetPracticeLine(words, 4)
		question := tutorLocalized(interfaceLanguage, "Какие слова из урока вы услышали?", "Which lesson words did you hear?")
		return text, question, tutorExpectedWordsForLanguage(words, tutorScenarioTemplate{}, language)
	}
	template := tutorScenarioTemplateFor(topic, function, interfaceLanguage)
	question := tutorLocalized(interfaceLanguage, template.ListeningQuestionRU, template.ListeningQuestionEN)
	return template.ListeningText, question, append([]string{}, template.ListeningExpected...)
}

func tutorScenarioDialogue(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string) ([]string, string, string) {
	if normalizeLearningLanguage(language) != "en" {
		line := tutorTargetPracticeLine(words, 3)
		prompt := tutorLocalized(interfaceLanguage, "Ответьте одной короткой фразой на изучаемом языке.", "Answer with one short target-language line.")
		goal := tutorLocalized(interfaceLanguage, "Использовать слова урока: ", "Use lesson words: ") + strings.Join(tutorLessonWordTexts(words, 3), ", ")
		return []string{line}, prompt, goal
	}
	template := tutorScenarioTemplateFor(topic, function, interfaceLanguage)
	prompt := tutorLocalized(interfaceLanguage, template.DialoguePromptRU, template.DialoguePromptEN)
	goal := tutorLocalized(interfaceLanguage, template.DialogueGoalRU, template.DialogueGoalEN)
	return append([]string{}, template.DialogueLines...), prompt, goal
}

func tutorExpectedWordsForLanguage(words []tutorLessonWord, template tutorScenarioTemplate, language string) []string {
	if normalizeLearningLanguage(language) == "en" {
		return tutorScenarioExpectedWords(template)
	}
	return tutorLessonWordTexts(words, tutorMinInt(4, len(words)))
}

func tutorTargetPracticeLine(words []tutorLessonWord, limit int) string {
	text := strings.Join(tutorLessonWordTexts(words, limit), " ")
	if strings.TrimSpace(text) != "" {
		return text
	}
	return "target language practice"
}

func tutorTargetSceneOrder(words []tutorLessonWord, limit int) []string {
	order := []string{"target"}
	for _, word := range tutorLessonWordTexts(words, limit) {
		order = append(order, word)
	}
	return order
}

func tutorTargetPracticeExplanation(words []tutorLessonWord, interfaceLanguage string) string {
	pairs := make([]string, 0, tutorMinInt(4, len(words)))
	for _, word := range words {
		wordText := strings.TrimSpace(word.Word)
		translation := strings.TrimSpace(word.Translation)
		if wordText == "" {
			continue
		}
		if translation != "" {
			pairs = append(pairs, wordText+" = "+translation)
		} else {
			pairs = append(pairs, wordText)
		}
		if len(pairs) >= 4 {
			break
		}
	}
	if len(pairs) == 0 {
		return tutorLocalized(interfaceLanguage, "Use the lesson words in a short target-language line.", "Use the lesson words in a short target-language line.")
	}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return "Соберите короткую фразу на изучаемом языке из слов урока: " + strings.Join(pairs, "; ") + "."
	}
	return "Build one short target-language line from the lesson words: " + strings.Join(pairs, "; ") + "."
}

func tutorReviewSummary(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) []string {
	count := tutorMinInt(4, len(words))
	summary := make([]string, 0, count+2)
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		summary = append(summary, "Вы закрыли сценарий: "+tutorTopicTitle(topic, interfaceLanguage))
		summary = append(summary, "Цель: "+function.CanDoRU)
	} else {
		summary = append(summary, "Scenario completed: "+tutorTopicTitle(topic, interfaceLanguage))
		summary = append(summary, "Can-do: "+function.CanDoEN)
	}
	for index := 0; index < count; index++ {
		word := words[index]
		summary = append(summary, fmt.Sprintf("%s - %s", word.Word, word.Translation))
	}
	return summary
}

func tutorPronunciationLine(words []tutorLessonWord, topic tutorTopic, language string) string {
	if normalizeLearningLanguage(language) != "en" {
		for _, word := range words {
			example := strings.TrimSpace(word.Example)
			if example != "" {
				return example
			}
		}
		return strings.TrimSpace(strings.Join(tutorLessonWordTexts(words, 4), " "))
	}
	word := func(index int, fallback string) string {
		if index >= 0 && index < len(words) && strings.TrimSpace(words[index].Word) != "" {
			return strings.TrimSpace(words[index].Word)
		}
		return fallback
	}
	switch topic.Code {
	case "hotel":
		return "I have a reservation and I need a quiet room, please."
	case "doctor":
		return "I have pain here and I need to see a doctor today."
	case "food":
		return "Could I have the menu and a glass of water, please?"
	case "transport":
		return "I need a ticket for the next train, please."
	case "work":
		return "Can we check the meeting time and the main task?"
	case "study":
		return "Could you help me with this lesson question today?"
	default:
		return fmt.Sprintf("Please help me with %s and %s today.", word(0, "this"), word(1, "practice"))
	}
}

func tutorLessonWordTexts(words []tutorLessonWord, limit int) []string {
	out := make([]string, 0, tutorMinInt(limit, len(words)))
	for _, word := range words {
		text := strings.TrimSpace(word.Word)
		if text == "" {
			continue
		}
		out = append(out, text)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func tutorScenarioExample(word string, translation string, interfaceLanguage string, topic tutorTopic, language string) (string, string) {
	if normalizeLearningLanguage(language) != "en" {
		return word, translation
	}
	example := fmt.Sprintf("Can you help me with %s?", word)
	switch topic.Code {
	case "hotel":
		if strings.EqualFold(word, "hotel") {
			example = "The hotel is near the station."
		} else if strings.EqualFold(word, "room") {
			example = "Can I have a quiet room?"
		} else if strings.EqualFold(word, "passport") {
			example = "Here is my passport."
		} else if strings.EqualFold(word, "reservation") {
			example = "I have a reservation for tonight."
		}
	case "doctor":
		if strings.EqualFold(word, "pain") {
			example = "I have pain here."
		} else if strings.EqualFold(word, "doctor") {
			example = "Can I see a doctor today?"
		} else if strings.EqualFold(word, "medicine") {
			example = "Do I need medicine?"
		}
	case "food":
		if strings.EqualFold(word, "menu") {
			example = "Can I see the menu?"
		} else if strings.EqualFold(word, "water") {
			example = "A glass of water, please."
		} else if strings.EqualFold(word, "coffee") {
			example = "Can I have coffee with breakfast?"
		}
	}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return example, fmt.Sprintf("Смысл в уроке: %s.", translation)
	}
	return example, translation
}

func tutorWordOrFallback(words []tutorLessonWord, index int, fallback string) tutorLessonWord {
	if index >= 0 && index < len(words) {
		return words[index]
	}
	return tutorLessonWord{ID: fallback, Word: fallback, Translation: fallback}
}

func tutorDialogue(words []tutorLessonWord, interfaceLanguage string) []string {
	if len(words) < 3 {
		return nil
	}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return []string{
			"А: " + words[0].Example + " " + words[0].ExampleRu,
			"Б: " + words[1].Example + " " + words[1].ExampleRu,
			"А: " + words[2].Example + " " + words[2].ExampleRu,
		}
	}
	return []string{
		"A: " + words[0].Example,
		"B: " + words[1].Example,
		"A: " + words[2].Example,
	}
}

func tutorExample(word string, translation string, interfaceLanguage string) (string, string) {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("I need %s.", word), fmt.Sprintf("Мне нужно: %s.", translation)
	}
	return fmt.Sprintf("I need %s.", word), fmt.Sprintf("%s.", translation)
}

func tutorGrammarTitle(level string, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		if level == "A1" {
			return "Паттерн A1: короткое утверждение"
		}
		return "Паттерн A2: просьба и уточнение"
	}
	if level == "A1" {
		return "A1 pattern: short statement"
	}
	return "A2 pattern: request and clarification"
}

func tutorGrammarBody(level string, languageName string, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		if level == "A1" {
			return "Собери фразу по схеме: кто/что + действие или состояние + одно полезное слово. Держи предложение коротким и понятным."
		}
		return "Собери просьбу: вежливое начало + нужное слово + уточнение времени, места или количества."
	}
	if level == "A1" {
		return "Build a short " + languageName + " sentence: subject + action or state + one useful word."
	}
	return "Build a polite request: opener + target word + time, place, or quantity detail."
}

func tutorTopicTitle(topic tutorTopic, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return topic.TitleRU
	}
	return topic.TitleEN
}

func tutorLocalized(interfaceLanguage string, ru string, en string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return ru
	}
	return en
}

func tutorTranslationForUser(word vocabWord, interfaceLanguage string) string {
	translation := wordTranslation(word, interfaceLanguage)
	if translation != "" {
		return translation
	}
	if normalizeLearningLanguage(word.Language) == "ru" && normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		if translation := wordTranslationForLanguage(word, "en"); translation != "" {
			return translation
		}
		if word.English != "" {
			return firstDictionaryValue(word.English)
		}
	}
	return ""
}

func tutorLooksClean(value string) bool {
	value = cleanDictionaryDisplay(value)
	if value == "" || !utf8.ValidString(value) {
		return false
	}
	lower := strings.ToLower(value)
	badFragments := []string{"пїЅ", "Гђ", "Г‘", "вЂ", "В«", "В»", "\\u", "<", ">", "\ufffd"}
	for _, fragment := range badFragments {
		if strings.Contains(lower, strings.ToLower(fragment)) {
			return false
		}
	}
	for _, r := range value {
		switch r {
		case 'Ђ', 'Ѓ', 'Љ', 'Њ', 'Ќ', 'Ћ', 'Џ', 'ђ', 'ѓ', 'љ', 'њ', 'ќ', 'ћ', 'џ', 'ў', 'ґ':
			return false
		}
	}
	if len([]rune(value)) > 220 {
		return false
	}
	return true
}

func looksLikeEnglishOnly(value string) bool {
	value = cleanDictionaryDisplay(value)
	if value == "" {
		return false
	}
	letters := 0
	latin := 0
	cyrillic := 0
	for _, r := range value {
		if !unicode.IsLetter(r) {
			continue
		}
		letters++
		switch {
		case unicode.In(r, unicode.Latin):
			latin++
		case unicode.In(r, unicode.Cyrillic):
			cyrillic++
		}
	}
	return letters > 0 && latin == letters && cyrillic == 0
}

func tutorStableWordKey(word vocabWord) string {
	if word.FrequencyRank > 0 {
		return fmt.Sprintf("%08d:%s", word.FrequencyRank, word.ID)
	}
	return "99999999:" + word.ID
}

func tutorStableIndex(seed string, max int) int {
	if max <= 0 {
		return 0
	}
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(seed))
	return int(hash.Sum32() % uint32(max))
}

func tutorMinInt(left int, right int) int {
	if left < right {
		return left
	}
	return right
}

func tutorMaxInt(left int, right int) int {
	if left > right {
		return left
	}
	return right
}
