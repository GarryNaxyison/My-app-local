import { api } from "./api";

type Stage = "today_opened" | "lesson_started" | "answer_submitted" | "answer_checked" | "retry_submitted" | "review_due" | "transfer_completed";

export function recordLearningEvent(stage: Stage, format: "" | "text" | "voice" = "", elapsed = 0, attemptId?: string) {
  if (!globalThis.crypto?.randomUUID) return;
  void api("/api/learning-events", {
    method: "POST",
    keepalive: true,
    body: { id: crypto.randomUUID(), attempt_id: attemptId, stage, format, elapsed_ms: Math.max(0, Math.min(86400000, Math.round(elapsed))) },
  }).catch(() => {});
}
