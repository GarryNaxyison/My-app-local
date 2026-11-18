import { useId, useState } from "react";
import { ArrowRight, Blocks, BookOpen, Check, Link2, RotateCcw } from "lucide-react";
import { api } from "../lib/api";
import type { VocabularyItem } from "../lib/types";
import { Button } from "./ui/button";
import "./learning-games.css";

type Game = "context" | "pairs" | "spelling";
type ContextQuestion = { context: string; options: string[]; correct_index: number; explanation: string };
type VocabularyCard = { word: string; translation: string };
type Round = { cards: VocabularyCard[]; translations: number[]; tiles: string[][] };

function shuffled<T>(items: T[]) {
  const result = [...items];
  for (let i = result.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [result[i], result[j]] = [result[j], result[i]];
  }
  return result;
}

function letters(word: string) {
  return word.normalize("NFC").match(/\P{Mark}\p{Mark}*|\p{Mark}+/gu) || [];
}

function makeRound(items: VocabularyItem[], game: Game): Round {
  const words = new Set<string>();
  const translations = new Set<string>();
  const cards: VocabularyCard[] = [];
  for (const item of shuffled(items)) {
    const word = item.word?.trim();
    const translation = item.translation?.trim();
    if (!word || !translation || words.has(word.toLowerCase()) || translations.has(translation.toLowerCase())) continue;
    if (game === "spelling" && (!/^\p{L}[\p{L}\p{M}]+$/u.test(word) || letters(word).length > 20)) continue;
    words.add(word.toLowerCase());
    translations.add(translation.toLowerCase());
    cards.push({ word, translation });
    if (cards.length === 5) break;
  }
  return { cards, translations: shuffled(cards.map((_, index) => index)), tiles: cards.map(({ word }) => shuffled(letters(word))) };
}

export function LearningGames({ russian, learningLanguage, openVocabulary }: {
  russian: boolean;
  learningLanguage?: string;
  openVocabulary: () => void;
}) {
  const t = (ru: string, en: string) => russian ? ru : en;
  const sourceId = useId();
  const [game, setGame] = useState<Game>("context");
  const [text, setText] = useState("");
  const [questions, setQuestions] = useState<ContextQuestion[]>([]);
  const [question, setQuestion] = useState(0);
  const [selected, setSelected] = useState<number | null>(null);
  const [contextScore, setContextScore] = useState(0);
  const [contextSource, setContextSource] = useState("");
  const [lessonId, setLessonId] = useState("");
  const [sourceBody, setSourceBody] = useState<{ preset: string } | { text: string } | null>(null);
  const [round, setRound] = useState<Round | null>(null);
  const [wordIndex, setWordIndex] = useState<number | null>(null);
  const [matched, setMatched] = useState<number[]>([]);
  const [puzzle, setPuzzle] = useState(0);
  const [picked, setPicked] = useState<number[]>([]);
  const [message, setMessage] = useState("");
  const [busy, setBusy] = useState(false);
  const modes = [
    { id: "context", title: t("Контекст", "Context"), caption: t("Выберите слово по смыслу", "Choose the word that fits"), Icon: BookOpen },
    { id: "pairs", title: t("Пары", "Pairs"), caption: t("Свяжите слово и перевод", "Match words and meanings"), Icon: Link2 },
    { id: "spelling", title: t("Собери слово", "Build a word"), caption: t("Восстановите написание", "Rebuild the spelling"), Icon: Blocks },
  ] as const;

  async function loadContext(body: { preset: string } | { text: string }, label: string) {
    if (busy) return;
    setBusy(true);
    setMessage("");
    try {
      const result = await api<{ lesson_id: string; questions: ContextQuestion[] }>("/api/games/context", { method: "POST", body });
      setQuestions(result.questions);
      setLessonId(result.lesson_id);
      setSourceBody(body);
      setContextSource(label);
      setQuestion(0);
      setSelected(null);
      setContextScore(0);
    } catch {
      setMessage(t("Не удалось подготовить задания. Попробуйте ещё раз.", "Could not prepare the exercises. Try again."));
    } finally {
      setBusy(false);
    }
  }

  async function nextQuestion() {
    if (busy) return;
    if (question + 1 === questions.length && sourceBody && lessonId) {
      setBusy(true);
      setMessage("");
      try {
        await api("/api/games/context", { method: "POST", body: { ...sourceBody, completed_id: lessonId, complete_only: true } });
      } catch {
        setMessage(t("Не удалось сохранить результат. Нажмите «Дальше» ещё раз.", "Could not save progress. Press Next again."));
        return;
      } finally { setBusy(false); }
    }
    setQuestion(question + 1);
    setSelected(null);
  }

  async function loadRound() {
    if (busy) return;
    setBusy(true);
    setMessage("");
    try {
      const data = await api<{ items: VocabularyItem[] }>("/api/vocabulary?page=0");
      setRound(makeRound(data.items || [], game));
      setWordIndex(null);
      setMatched([]);
      setPuzzle(0);
      setPicked([]);
    } catch {
      setMessage(t("Не удалось загрузить словарь.", "Could not load your vocabulary."));
    } finally {
      setBusy(false);
    }
  }

  const current = questions[question];
  const card = round?.cards[puzzle];
  const tiles = round?.tiles[puzzle] || [];
  const built = picked.map((index) => tiles[index]).join("");
  const solved = Boolean(card && built === card.word.normalize("NFC"));
  const complete = game === "pairs" ? Boolean(round?.cards.length && matched.length === round.cards.length) : Boolean(round?.cards.length && puzzle === round.cards.length);

  return (
    <div className="learning-games">
      <header className="learning-games__header">
        <span className="learning-games__art" aria-hidden="true" />
        <div><h2>{t("Учебные игры", "Learning games")}</h2><p>{t("Выберите, что потренировать", "Choose what to practise")}</p></div>
      </header>
      <nav className="learning-games__modes" aria-label={t("Режимы игры", "Game modes")}>
        {modes.map(({ id, title, caption, Icon }) => <button type="button" key={id} aria-pressed={game === id} disabled={busy} onClick={() => {
          setGame(id); setRound(null); setMessage("");
        }}><Icon size={23} aria-hidden="true" /><span><strong>{title}</strong><small>{caption}</small></span></button>)}
      </nav>
      {game === "context" ? (
        <section className="v2-panel learning-games__board" aria-busy={busy}>
          <div className="learning-games__presets">
            {([["travel", t("Поездка", "Travel")], ["restaurant", t("Ресторан", "Restaurant")], ["work", t("Работа", "Work")]] as const)
              .map(([id, title]) => <Button type="button" variant="outline" key={id} disabled={busy} onClick={() => void loadContext({ preset: id }, title)}>{title}</Button>)}
          </div>
          <details>
            <summary>{t("Свой текст", "Your own text")}</summary>
            <form onSubmit={(event) => {
              event.preventDefault();
              if (text.trim()) void loadContext({ text: text.trim() }, t("Свой текст", "Your own text"));
            }}>
              <label htmlFor={sourceId}>{t("Вставьте короткий текст", "Paste a short text")}</label>
              <textarea id={sourceId} value={text} maxLength={1200} disabled={busy} onChange={(event) => setText(event.target.value)} />
              <Button type="submit" disabled={busy || !text.trim()}>{t("Подготовить задания", "Prepare exercises")}</Button>
            </form>
          </details>
          {current && <article key={`${lessonId}-${question}`} className="learning-games__question">
            <header><strong>{contextSource}</strong><span>{question + 1} / {questions.length}</span></header>
            <p>{current.context}</p>
            <div className="learning-games__choices" lang={learningLanguage} dir="auto">
              {current.options.map((option, index) => <button type="button" key={option} aria-pressed={selected === index} disabled={selected !== null || busy}
                data-result={selected === null ? undefined : index === current.correct_index ? "correct" : selected === index ? "wrong" : undefined}
                onClick={() => {
                  setSelected(index);
                  if (index === current.correct_index) setContextScore((score) => score + 1);
                }}>{option}</button>)}
            </div>
            {selected !== null && <>
              <p role="status"><strong>{selected === current.correct_index ? t("Верно. ", "Correct. ") : t("Правильный вариант: ", "Correct option: ") + current.options[current.correct_index] + ". "}</strong>{current.explanation}</p>
              <Button type="button" disabled={busy} onClick={() => void nextQuestion()}>{t("Дальше", "Next")} <ArrowRight size={16} /></Button>
            </>}
          </article>}
          {questions.length > 0 && !current && <div className="learning-games__complete">
            <Check size={36} aria-hidden="true" />
            <p role="status">{t("Результат", "Score")}: {contextScore} / {questions.length}</p>
            <Button type="button" disabled={busy || !sourceBody} onClick={() => sourceBody && void loadContext(sourceBody, contextSource)}>{t("Следующая игра", "Next game")} <ArrowRight size={18} /></Button>
          </div>}
        </section>
      ) : (
        <section className="v2-panel learning-games__board" aria-busy={busy}>
          {(!round || complete || round.cards.length < (game === "pairs" ? 2 : 1)) && <div className="learning-games__presets">
            <Button type="button" disabled={busy} onClick={() => void loadRound()}><RotateCcw size={16} />{round ? t("Новый раунд", "New round") : t("Начать игру", "Start game")}</Button>
            {round && !complete && <Button type="button" variant="outline" onClick={openVocabulary}>{t("Открыть словарь", "Open vocabulary")}</Button>}
          </div>}
          {round && round.cards.length < (game === "pairs" ? 2 : 1) && <p>{t("Добавьте подходящие слова с переводом в словарь.", "Add suitable words with translations to your vocabulary.")}</p>}
          {game === "pairs" && round && round.cards.length >= 2 && <>
            <strong>{matched.length} / {round.cards.length}</strong>
            <div className="learning-games__pairs">
              <svg className="learning-games__connections" viewBox={`0 0 100 ${round.cards.length * 100}`} preserveAspectRatio="none" aria-hidden="true">
                {matched.map(index => <path key={index} pathLength="1" d={`M 0 ${index * 100 + 50} C 45 ${index * 100 + 50}, 55 ${round.translations.indexOf(index) * 100 + 50}, 100 ${round.translations.indexOf(index) * 100 + 50}`} />)}
              </svg>
              <div lang={learningLanguage}>{round.cards.map(({ word }, index) => <button type="button" key={word} aria-pressed={wordIndex === index} disabled={matched.includes(index)}
                data-result={matched.includes(index) ? "correct" : undefined} onClick={() => { setWordIndex(index); setMessage(""); }}>{word}</button>)}</div>
              <div>{round.translations.map((index) => <button type="button" key={index} disabled={wordIndex === null || matched.includes(index)}
                data-result={matched.includes(index) ? "correct" : undefined} onClick={() => {
                  if (wordIndex === index) { setMatched([...matched, index]); setMessage(t("Пара найдена", "Match found")); }
                  else setMessage(t("Попробуйте другой перевод.", "Try another meaning."));
                  setWordIndex(null);
                }}>{round.cards[index].translation}</button>)}</div>
            </div>
          </>}
          {game === "spelling" && card && <>
            <strong>{puzzle + 1} / {round?.cards.length} · {card.translation}</strong>
            <div className="learning-games__puzzle" lang={learningLanguage} aria-live="polite" dir="auto">{built || "· · ·"}</div>
            <div className="learning-games__choices" lang={learningLanguage} dir="auto">
              {tiles.map((letter, index) => <button type="button" key={index} disabled={solved || picked.includes(index)}
                onClick={() => { setPicked([...picked, index]); setMessage(""); }}>{letter}</button>)}
            </div>
            <div className="learning-games__presets">
              {!solved && <Button type="button" variant="outline" disabled={!picked.length} onClick={() => setPicked(picked.slice(0, -1))}>{t("Убрать букву", "Remove letter")}</Button>}
              {solved && <Button type="button" onClick={() => { setPuzzle(puzzle + 1); setPicked([]); }}><Check size={16} />{t("Дальше", "Next")}</Button>}
            </div>
            {picked.length === tiles.length && !solved && <p role="status">{t("Проверьте порядок букв.", "Check the order of the letters.")}</p>}
          </>}
          {complete && <p role="status">{t("Раунд пройден", "Round complete")}</p>}
        </section>
      )}
      {busy && <p role="status">{t("Загрузка…", "Loading…")}</p>}
      {message && <p role="status">{message}</p>}
    </div>
  );
}
