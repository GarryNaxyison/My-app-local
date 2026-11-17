import { useState } from "react";
import { ArrowDown, BookOpen, Camera, Check, CircleHelp, LayoutDashboard, MapPin, MessageCircle, Mic, NotebookPen, Repeat2 } from "lucide-react";
import type { LandingContent, LandingLocale } from "./landingContent";
import "./productWindows.css";

type WindowKind = keyof LandingContent["images"];
type Example = { title: string; context: string; phrase: string; action: string; result: string };

const examples = {
  ru: {
    dashboard: { title: "Сегодня", context: "Английский · A1", phrase: "Продолжить урок", action: "Посмотреть план", result: "Урок → практика → повторение слов." },
    aiTutor: { title: "AI-репетитор", context: "Английский · A1", phrase: "Could you show me your passport?", action: "Посмотреть пример ответа", result: "Yes, here it is." },
    mistakes: { title: "Разбор ошибок", context: "Артикли", phrase: "I have reservation.", action: "Показать правку", result: "I have a reservation. Перед reservation нужен артикль a." },
    voice: { title: "Произношение", context: "Интонация в вопросе", phrase: "Could I check in now?", action: "Открыть подсказку", result: "Выделите голосом check in, а now произнесите мягче." },
    photo: { title: "Перевод фото", context: "Фрагмент меню", phrase: "Breakfast until 10:30", action: "Показать перевод", result: "Завтрак до 10:30." },
    notes: { title: "Мои фразы", context: "Поездки", phrase: "Could you hold my luggage?", action: "Открыть заметку", result: "Не могли бы вы оставить у себя мой багаж?" },
    telegram: { title: "NERIVA в Telegram", context: "Пример сообщения бота", phrase: "How do you ask for the Wi-Fi password?", action: "Показать ответ", result: "Could you tell me the Wi-Fi password?" },
  },
  en: {
    dashboard: { title: "Today", context: "English · A1", phrase: "Continue your lesson", action: "See today’s plan", result: "Lesson → practice → word review." },
    aiTutor: { title: "AI Tutor", context: "English · A1", phrase: "Could you show me your passport?", action: "See an example answer", result: "Yes, here it is." },
    mistakes: { title: "Mistake review", context: "Articles", phrase: "I have reservation.", action: "Show the correction", result: "I have a reservation. Use a before reservation." },
    voice: { title: "Pronunciation", context: "Question intonation", phrase: "Could I check in now?", action: "See a pronunciation tip", result: "Stress check in and say now more softly." },
    photo: { title: "Photo translation", context: "From a restaurant menu", phrase: "Petit-déjeuner jusqu’à 10 h 30", action: "Show the translation", result: "Breakfast until 10:30." },
    notes: { title: "My phrases", context: "Travel", phrase: "Could you hold my luggage?", action: "Open the note", result: "Use this to ask the hotel to store your bags." },
    telegram: { title: "NERIVA on Telegram", context: "Example bot message", phrase: "How do you ask for the Wi-Fi password?", action: "See the answer", result: "Could you tell me the Wi-Fi password?" },
  },
} satisfies Record<LandingLocale, Record<string, Example>>;

const windowExamples = {
  dashboard: "dashboard", mobileHome: "dashboard", aiTutor: "aiTutor", mobileLesson: "aiTutor",
  mistakes: "mistakes", mobileMistakes: "mistakes", voice: "voice", mobileVoice: "voice",
  photo: "photo", notes: "notes", telegram: "telegram",
} as const satisfies Record<WindowKind, keyof typeof examples.en>;

const icons = { dashboard: LayoutDashboard, aiTutor: BookOpen, mistakes: Repeat2, voice: Mic, photo: Camera, notes: NotebookPen, telegram: MessageCircle };

export function ProductWindow({ kind, language }: { kind: WindowKind; language: LandingLocale }) {
  const preset = windowExamples[kind];
  const copy = examples[language][preset];
  const Icon = icons[preset];
  const phraseLanguage = preset === "dashboard" ? language : preset === "photo" && language === "en" ? "fr" : "en";

  return (
    <figure className="product-window" data-window={preset}>
      <div className="product-window__title"><Icon size={16} aria-hidden="true" /> {copy.title}</div>
      <div className="product-window__body">
        <span className="product-window__context">{copy.context}</span>
        <p className="product-window__phrase" lang={phraseLanguage}>{copy.phrase}</p>
        {preset === "voice" && <div className="product-window__wave" aria-hidden="true">{[5, 10, 18, 12, 23, 17, 9, 15, 21, 12, 6].map((height, index) => <i key={index} style={{ height }} />)}</div>}
        <details>
          <summary>{copy.action}<ArrowDown size={14} aria-hidden="true" /></summary>
          <p className="product-window__result">{copy.result}</p>
        </details>
      </div>
    </figure>
  );
}

export function HeroLesson({ language }: { language: LandingLocale }) {
  const [answer, setAnswer] = useState<string>();
  const russian = language === "ru";
  const correctAnswer = "I have a reservation.";
  const correct = answer === correctAnswer;
  const explanation = russian
    ? "Нужен артикль a перед reservation в единственном числе."
    : "Use a before the singular countable noun reservation.";

  return (
    <div className="hero-product-frame hero-windows">
      <div className="product-window hero-windows__lesson" data-window="aiTutor">
        <div className="product-window__title">
          <BookOpen size={19} aria-hidden="true" />
          <strong>{russian ? "Пробный урок" : "Sample lesson"}</strong>
          <span className="hero-windows__level">A1</span>
        </div>
        <fieldset className="lesson-example">
          <legend><span><MapPin size={15} aria-hidden="true" />{russian ? "Заселение в отель" : "Hotel check-in"}</span></legend>
          <p>{russian ? "Сообщите на ресепшене, что у вас забронирован номер." : "Tell the receptionist you have a reservation."}</p>
          <div className="lesson-example__answers">
            {["I have reservation.", correctAnswer].map((option) => (
              <button type="button" key={option} lang="en" aria-pressed={answer === option} onClick={() => setAnswer(option)}>
                {option}
                {answer === option && (correct ? <Check size={17} aria-hidden="true" /> : <CircleHelp size={17} aria-hidden="true" />)}
              </button>
            ))}
          </div>
          <p className="lesson-example__feedback" role="status">
            {correct && <strong>{russian ? "Верно. " : "Correct. "}</strong>}
            {answer && explanation}
          </p>
        </fieldset>
      </div>
      <div className="product-window hero-windows__review" data-window="notes">
        <div className="product-window__title"><Repeat2 size={16} aria-hidden="true" />{russian ? "К повторению" : "Review"}</div>
        <div className="product-window__body">
          {correct ? <p className="product-window__phrase" lang="en">{correctAnswer}</p> : <p>{russian ? "Здесь появится фраза из урока." : "Your lesson phrase will appear here."}</p>}
        </div>
      </div>
    </div>
  );
}
