import { useEffect, useState } from "react";
import { Repeat2 } from "lucide-react";
import { api } from "../lib/api";
import { normalizeAppLocale } from "../lib/i18n";
import { Button } from "./ui/button";

type ReviewItem = {
  id: string;
  source: string;
  lesson_id: string;
  title: string;
  due_at: string;
  status: "due" | "upcoming";
};

type Props = {
  requestedId?: string;
  compact?: boolean;
  interfaceLanguage?: string;
  canStart: boolean;
  hasActiveLesson: boolean;
  busy: boolean;
  start: (lessonId: string) => Promise<boolean>;
  startMistakePractice: (reviewId: string) => Promise<void>;
};

export function TodayReview({ requestedId, compact = false, interfaceLanguage, canStart, hasActiveLesson, busy, start, startMistakePractice }: Props) {
  const [item, setItem] = useState<ReviewItem | null>(null);
  const [loaded, setLoaded] = useState(false);
  useEffect(() => {
    const controller = new AbortController();
    void api<{ items: ReviewItem[] }>("/api/learning-reviews", { signal: controller.signal }).then((response) => {
      if (controller.signal.aborted) return;
      const next = response.items.find((entry) =>
        (!requestedId || entry.id === requestedId) &&
        (entry.source === "ai_tutor" && entry.lesson_id || entry.source === "mistake" && entry.id) && Number.isFinite(Date.parse(entry.due_at)),
      );
      setItem(next ?? null);
      setLoaded(true);
    }).catch(() => {
      if (!controller.signal.aborted) setItem(null);
    });
    return () => controller.abort();
  }, [requestedId]);

  const locale = normalizeAppLocale(interfaceLanguage);
  const russian = locale === "ru";
  if (!item) return loaded && requestedId
    ? <section className="v2-panel today-review" role="status">{russian ? "Повторение недоступно" : "Review unavailable"}</section>
    : null;
  const due = item.status === "due";
  const mistake = item.source === "mistake";
  const formattedDate = new Intl.DateTimeFormat(locale, { dateStyle: "medium", timeStyle: "short" }).format(new Date(item.due_at));
  const Heading = compact ? "h3" : "h2";

  return (
    <section className={compact ? "today-review today-review--compact" : "v2-panel today-review"}>
      <Repeat2 size={22} aria-hidden="true" />
      <div className="today-review__details">
        <Heading>{due ? russian ? "Пора повторить" : "Review due" : russian ? "Ближайшее повторение" : "Next review"}</Heading>
        <p dir="auto">{item.title}</p>
        {!mistake && <time dateTime={item.due_at}>{formattedDate}</time>}
        {due && hasActiveLesson && <p>{russian ? "Сначала завершите текущий урок." : "Finish the current lesson first."}</p>}
      </div>
      {due && (mistake || canStart) && (
        <Button type="button" variant="outline" disabled={hasActiveLesson || busy} onClick={() => void (mistake ? startMistakePractice(item.id) : start(item.lesson_id))}>
          {mistake ? (russian ? "Исправить ошибку" : "Fix the mistake") : (russian ? "Повторить урок" : "Review lesson")}
        </Button>
      )}
    </section>
  );
}
