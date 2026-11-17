import { useId, useRef, useState } from "react";
import { Button } from "./ui/button";
import { api } from "../lib/api";

const maxLength = 1200;

export function TextPractice({ russian }: { russian: boolean }) {
  const id = useId();
  const [source, setSource] = useState("");
  const [material, setMaterial] = useState("");
  const [question, setQuestion] = useState("");
  const [answer, setAnswer] = useState("");
  const [feedback, setFeedback] = useState("");
  const [busy, setBusy] = useState(false);
  const [generated, setGenerated] = useState<string[]>([]);
  const [usingFallback, setUsingFallback] = useState(false);
  const pending = useRef(false);
  const questions = russian
    ? ["Сформулируйте главную мысль в одном предложении.", "Выберите важную деталь и объясните, почему она важна."]
    : ["Summarise the main idea in one sentence.", "Choose a key detail and explain why it matters."];
  return (
    <details className="v2-panel text-practice">
      <summary>{russian ? "Практика из своего текста" : "Practise with your own text"}</summary>
      <form onSubmit={async (event) => {
        event.preventDefault();
        const text = source.trim();
        if (!text || source.length > maxLength || pending.current) return;
        pending.current = true;
        setBusy(true);
        setQuestion("");
        setAnswer("");
        setFeedback("");
        setUsingFallback(false);
        try {
          const result = await api<{ questions: string[] }>("/api/practice/text-questions", { method: "POST", body: { text } });
          if (!Array.isArray(result.questions) || result.questions.length !== 2 || !result.questions.every((value) => typeof value === "string" && value.trim())) throw new Error("Invalid questions");
          setGenerated(result.questions);
        } catch {
          setGenerated(questions);
          setUsingFallback(true);
        } finally {
          setMaterial(text);
          pending.current = false;
          setBusy(false);
        }
      }}>
        <label htmlFor={id}>{russian ? "Вставьте короткий текст" : "Paste a short text"}</label>
        <textarea id={id} value={source} aria-describedby={`${id}-length`} onChange={(event) => {
          setSource(event.target.value);
          setMaterial("");
          setQuestion("");
          setAnswer("");
          setFeedback("");
        }} disabled={busy} />
        <p id={`${id}-length`}>{source.length} / {maxLength}</p>
        <Button type="submit" disabled={busy || !source.trim() || source.length > maxLength}>
          {russian ? "Подготовить вопросы" : "Prepare questions"}
        </Button>
      </form>
      {usingFallback && material && <p>{russian ? "Вопросы по тексту недоступны. Можно выбрать готовый вопрос." : "Text-specific questions are unavailable. Choose a standard question below."}</p>}
      {material && <ul>{generated.map((question) => <li key={question}>
        <Button type="button" variant="outline" disabled={busy} onClick={() => { setQuestion(question); setAnswer(""); setFeedback(""); }}>{question}</Button>
      </li>)}</ul>}
      {question && <form onSubmit={async (event) => {
        event.preventDefault();
        if (!answer.trim() || pending.current) return;
        pending.current = true;
        setBusy(true);
        setFeedback("");
        try {
          const response = await api<{ reply?: string }>("/api/practice", {
            method: "POST",
            body: { text: `Review the learner's answer in the context of this reading exercise.\n${JSON.stringify({ material, question, answer: answer.trim() })}` },
          });
          const reply = response.reply?.trim();
          if (!reply) throw new Error("empty practice response");
          setFeedback(reply);
        } catch {
          setFeedback(russian ? "Не удалось проверить ответ. Попробуйте снова." : "Could not check your answer. Try again.");
        } finally {
          pending.current = false;
          setBusy(false);
        }
      }}>
        <p>{question}</p>
        <label htmlFor={`${id}-answer`}>{russian ? "Ваш ответ" : "Your answer"}</label>
        <textarea id={`${id}-answer`} value={answer} maxLength={maxLength} onChange={(event) => setAnswer(event.target.value)} disabled={busy} />
        <Button type="submit" disabled={busy || !answer.trim()}>{russian ? "Проверить ответ" : "Check answer"}</Button>
      </form>}
      {feedback && <p role="status" style={{ whiteSpace: "pre-wrap" }}>{feedback}</p>}
    </details>
  );
}
