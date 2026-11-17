import { useState } from "react";
import { Button } from "./ui/button";
import { WordCardPrint } from "./WordCardPrint";

const cards = [
  { id: "reservation", word: "reservation", translation: "бронирование", example: "I have a reservation." },
  { id: "ticket", word: "ticket", translation: "билет", example: "I need a ticket." },
  { id: "menu", word: "menu", translation: "меню", example: "Could I see the menu?" },
];

const lessons = [
  { ru: "На стойке отеля", en: "At the hotel", prompt: "I have a ___.", choices: ["reservation", "menu", "platform"] },
  { ru: "На вокзале", en: "At the station", prompt: "I need a ___ for the next train.", choices: ["menu", "ticket", "room"] },
  { ru: "В ресторане", en: "At the restaurant", prompt: "Could I see the ___ before I order?", choices: ["platform", "reservation", "menu"] },
] as const;

export function TravelRoute({ russian, owner, onRoleplay }: { russian: boolean; owner: string; onRoleplay: () => void }) {
  const storageKey = `neriva-travel-route-v1:${owner}`;
  const [step, setStep] = useState(() => {
    try {
      const saved = Number(localStorage.getItem(storageKey));
      return Number.isInteger(saved) && saved >= 0 && saved <= lessons.length ? saved : 0;
    } catch {
      return 0;
    }
  });
  const [wrong, setWrong] = useState(false);
  const lesson = lessons[step];
  function saveStep(next: number) {
    setStep(next);
    try {
      localStorage.setItem(storageKey, String(next));
    } catch {
      // The current session can continue when browser storage is unavailable.
    }
  }
  return (
    <details className="v2-panel travel-route">
      <summary>{russian ? "Маршрут «Поездка»" : "Travel practice route"}</summary>
      <p>{russian ? "Три коротких задания на английском" : "Three short exercises in English"} · {step}/{lessons.length}</p>
      {lesson ? (
        <>
          <h3>{russian ? lesson.ru : lesson.en}</h3>
          <p lang="en">{lesson.prompt}</p>
          <div className="home-insights-v2__actions">
            {lesson.choices.map((choice) => (
              <Button key={choice} lang="en" type="button" variant="outline" onClick={() => {
                if (choice !== cards[step].word) {
                  setWrong(true);
                  return;
                }
                setWrong(false);
                saveStep(step + 1);
              }}>{choice}</Button>
            ))}
          </div>
          {wrong && <p role="status">{russian ? "Попробуйте ещё раз." : "Try again."}</p>}
        </>
      ) : (
        <>
          <p role="status">{russian ? "Все задания выполнены. Повторите фразы или попробуйте сцену «Ресторан»." : "All exercises complete. Review the phrases or try the Restaurant scene."}</p>
          <Button type="button" onClick={onRoleplay}>{russian ? "Начать сцену «Ресторан»" : "Start the Restaurant scene"}</Button>
          <Button type="button" variant="outline" onClick={() => saveStep(0)}>{russian ? "Повторить маршрут" : "Repeat the route"}</Button>
          <WordCardPrint items={cards} russian={russian} />
        </>
      )}
    </details>
  );
}
