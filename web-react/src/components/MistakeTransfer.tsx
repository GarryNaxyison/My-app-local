import { useId, useState } from "react";
import { Button } from "./ui/button";
import { checkMistakeTransfer, type MistakeTransferTask } from "../lib/mistakeTransfer";

export function MistakeTransfer({ task, russian, onClose }: {
  task: MistakeTransferTask;
  russian: boolean;
  onClose: () => void;
}) {
  const id = useId();
  const locale = russian ? "ru" : "en";
  const [answer, setAnswer] = useState("");
  const [result, setResult] = useState<"retry" | "correct" | null>(null);

  return (
    <section className="v2-panel mistake-transfer" aria-labelledby={id}>
      <div className="panel-head">
        <h2 id={id}>{russian ? "Та же конструкция в новой ситуации" : "Use the pattern in a new situation"}</h2>
        <Button type="button" variant="outline" onClick={onClose}>{russian ? "К списку ошибок" : "Back to mistakes"}</Button>
      </div>
      <p>{task.prompt[locale]}</p>
      <form className="inline-form-v2" onSubmit={(event) => {
        event.preventDefault();
        if (answer.trim()) setResult(checkMistakeTransfer(task, answer) ? "correct" : "retry");
      }}>
        <input aria-label={russian ? "Ответ на английском" : "Answer in English"} lang="en" value={answer} onChange={(event) => {
          setAnswer(event.target.value);
          setResult(null);
        }} />
        <Button type="submit" disabled={!answer.trim()}>{russian ? "Проверить" : "Check"}</Button>
      </form>
      <div role="status">
        {result === "correct" && <p>{russian ? "Верно. Вы применили конструкцию в новой ситуации." : "Correct. You used the pattern in a new situation."}</p>}
        {result === "retry" && <p>{russian ? "Попробуйте ещё раз." : "Try again."} {task.hint[locale]}</p>}
      </div>
      <details>
        <summary>{russian ? "Показать пример" : "Show an example"}</summary>
        <p lang="en">{task.example}</p>
      </details>
    </section>
  );
}
