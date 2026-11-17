import { useEffect, useState } from "react";
import { Repeat2 } from "lucide-react";
import { api } from "../lib/api";
import { normalizeAppLocale } from "../lib/i18n";
import { Button } from "./ui/button";

type ReviewItem = {
  source: string;
  lesson_id: string;
  title: string;
  due_at: string;
  status: "due" | "upcoming";
};

type Props = {
  interfaceLanguage?: string;
  canStart: boolean;
  hasActiveLesson: boolean;
  busy: boolean;
  start: (lessonId: string) => Promise<boolean>;
};

export function TodayReview({ interfaceLanguage, canStart, hasActiveLesson, busy, start }: Props) {
  const [item, setItem] = useState<ReviewItem | null>(null);
  useEffect(() => {
    const controller = new AbortController();
    void api<{ items: ReviewItem[] }>("/api/learning-reviews", { signal: controller.signal }).then((response) => {
      if (controller.signal.aborted) return;
      const next = response.items.find((entry) => entry.source === "ai_tutor" && entry.lesson_id && Number.isFinite(Date.parse(entry.due_at)));
      setItem(next ?? null);
    }).catch(() => {
      if (!controller.signal.aborted) setItem(null);
    });
    return () => controller.abort();
  }, []);

  if (!item) return null;
  const locale = normalizeAppLocale(interfaceLanguage);
  const russian = locale === "ru";
  const due = item.status === "due";
  const formattedDate = new Intl.DateTimeFormat(locale, { dateStyle: "medium", timeStyle: "short" }).format(new Date(item.due_at));

  return (
    <section className="v2-panel today-review">
      <Repeat2 size={22} aria-hidden="true" />
      <div className="today-review__details">
        <h2>{due ? russian ? "Пора повторить" : "Review due" : russian ? "Ближайшее повторение" : "Next review"}</h2>
        <p dir="auto">{item.title}</p>
        <time dateTime={item.due_at}>{formattedDate}</time>
        {due && hasActiveLesson && <p>{russian ? "Сначала завершите текущий урок." : "Finish the current lesson first."}</p>}
      </div>
      {due && canStart && (
        <Button type="button" variant="outline" disabled={hasActiveLesson || busy} onClick={() => void start(item.lesson_id)}>
          {russian ? "Повторить урок" : "Review lesson"}
        </Button>
      )}
    </section>
  );
}
