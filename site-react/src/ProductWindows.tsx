import { useState } from "react";
import { ArrowDown, BookOpen, Camera, LayoutDashboard, MessageCircle, Mic, NotebookPen, Repeat2 } from "lucide-react";
import type { LandingContent, LandingLocale } from "./landingContent";
import "./productWindows.css";

type WindowKind = keyof LandingContent["images"];
const windows = {
  dashboard: { view: "home", label: { ru: "Сегодня", en: "Today" }, icon: LayoutDashboard },
  mobileHome: { view: "home", label: { ru: "Сегодня", en: "Today" }, icon: LayoutDashboard },
  aiTutor: { view: "tutor", label: { ru: "AI-репетитор", en: "AI Tutor" }, icon: BookOpen },
  mobileLesson: { view: "lesson", label: { ru: "Урок", en: "Lesson" }, icon: BookOpen },
  mistakes: { view: "mistakes", label: { ru: "Разбор ошибок", en: "Mistake review" }, icon: Repeat2 },
  mobileMistakes: { view: "mistakes", label: { ru: "Разбор ошибок", en: "Mistake review" }, icon: Repeat2 },
  voice: { view: "pronunciation", label: { ru: "Произношение", en: "Pronunciation" }, icon: Mic },
  mobileVoice: { view: "pronunciation", label: { ru: "Произношение", en: "Pronunciation" }, icon: Mic },
  photo: { view: "tools", label: { ru: "Перевод фото", en: "Photo translation" }, icon: Camera },
  notes: { view: "phrasebook", label: { ru: "Мои фразы", en: "My phrases" }, icon: NotebookPen },
  telegram: { view: "lesson", label: { ru: "Урок в веб-приложении", en: "Web app lesson" }, icon: MessageCircle },
} satisfies Record<WindowKind, unknown>;

function PreviewFrame({ kind, language }: { kind: WindowKind; language: LandingLocale }) {
  const window = windows[kind];
  return <iframe className="product-window__frame" src={`/product-preview.html?view=${window.view}&lang=${language}`} title={window.label[language]} sandbox="allow-scripts allow-same-origin allow-forms" allow="camera 'none'; microphone 'none'" />;
}

export function ProductWindow({ kind, language }: { kind: WindowKind; language: LandingLocale }) {
  const [loaded, setLoaded] = useState(false);
  const { label, icon: Icon } = windows[kind];
  return (
    <details className="product-window" data-window={kind} onToggle={(event) => { if (event.currentTarget.open) setLoaded(true); }}>
      <summary className="product-window__title"><Icon size={17} aria-hidden="true" />{label[language]}<ArrowDown className="product-window__arrow" size={15} aria-hidden="true" /></summary>
      {loaded && <PreviewFrame kind={kind} language={language} />}
    </details>
  );
}

export function HeroLesson({ language }: { language: LandingLocale }) {
  return (
    <div className="hero-product-frame hero-windows">
      <div className="product-window hero-windows__lesson" data-window="tutor">
        <div className="product-window__title"><BookOpen size={19} aria-hidden="true" /><strong>{language === "ru" ? "Пробный урок" : "Sample lesson"}</strong></div>
        <PreviewFrame kind="aiTutor" language={language} />
      </div>
      <ProductWindow kind="mistakes" language={language} />
    </div>
  );
}
