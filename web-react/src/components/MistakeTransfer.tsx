import { useId, useRef, useState } from "react";
import { Button } from "./ui/button";
import { checkMistakeTransfer, type MistakeTransferTask } from "../lib/mistakeTransfer";
import { api } from "../lib/api";

export function MistakeTransfer({ task, reviewId, russian, savedAnswer, onSaved, onPremium, onClose }: {
  task: MistakeTransferTask;
  reviewId?: string;
  russian: boolean;
  savedAnswer?: string;
  onSaved?: (answer: string) => void;
  onPremium?: () => void;
  onClose: () => void;
}) {
  const id = useId();
  const locale = russian ? "ru" : "en";
  const [answer, setAnswer] = useState(savedAnswer ?? "");
  const [result, setResult] = useState<"retry" | "correct" | null>(savedAnswer ? "correct" : null);
  const [saveState, setSaveState] = useState<"saving" | "saved" | "error" | null>(savedAnswer ? "saved" : null);
  const saving = useRef(false);

  async function saveProgress() {
    if (!reviewId || saving.current) return;
    saving.current = true;
    setSaveState("saving");
    try {
      await api("/api/mistakes/transfer", { method: "POST", body: { review_id: reviewId } });
      setSaveState("saved");
      onSaved?.(answer);
    } catch {
      setSaveState("error");
    } finally {
      saving.current = false;
    }
  }

  return (
    <section className="v2-panel mistake-transfer" aria-labelledby={id}>
      <div className="panel-head">
        <h2 id={id}>{russian ? "Та же конструкция в новой ситуации" : "Use the pattern in a new situation"}</h2>
        <Button type="button" variant="outline" onClick={onClose}>{russian ? "К списку ошибок" : "Back to mistakes"}</Button>
      </div>
      <p>{task.prompt[locale]}</p>
      <form className="inline-form-v2" onSubmit={(event) => {
        event.preventDefault();
        if (!answer.trim() || result === "correct") return;
        const correct = checkMistakeTransfer(task, answer);
        setResult(correct ? "correct" : "retry");
        if (correct) void saveProgress();
      }}>
        <input aria-label={russian ? "Ответ на английском" : "Answer in English"} lang="en" value={answer} readOnly={result === "correct"} onChange={(event) => {
          setAnswer(event.target.value);
          setResult(null);
        }} />
        <Button type="submit" disabled={!answer.trim() || result === "correct"}>{russian ? "Проверить" : "Check"}</Button>
      </form>
      <div role="status">
        {result === "correct" && <p>{russian ? "Верно. Вы применили конструкцию в новой ситуации." : "Correct. You used the pattern in a new situation."}</p>}
        {result === "retry" && <p>{russian ? "Попробуйте ещё раз." : "Try again."} {task.hint[locale]}</p>}
        {saveState === "saving" && <p>{russian ? "Сохраняем результат…" : "Saving your result…"}</p>}
        {saveState === "saved" && <p>{russian ? "Результат сохранён в профиле." : "Result saved to your profile."}</p>}
        {saveState === "error" && <p>{russian ? "Не удалось сохранить результат." : "Could not save your result."}</p>}
      </div>
      {saveState === "error" && <Button type="button" variant="outline" onClick={() => void saveProgress()}>{russian ? "Повторить сохранение" : "Retry saving"}</Button>}
      {saveState === "saved" && onPremium && (
        <div className="mistake-transfer__premium">
          <p>{russian ? "С Premium можно тренировать произношение и проверять голосовые записи." : "With Premium you can practise pronunciation and get feedback on your recordings."}</p>
          <Button type="button" variant="outline" onClick={onPremium}>{russian ? "Посмотреть Premium" : "Explore Premium"}</Button>
        </div>
      )}
      <details>
        <summary>{russian ? "Показать пример" : "Show an example"}</summary>
        <p lang="en">{task.example}</p>
      </details>
    </section>
  );
}
