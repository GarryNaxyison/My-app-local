import { useState } from "react";
import { Button } from "./ui/button";
import { api } from "../lib/api";

type Card = { id: string; title: string; answer?: string };
const ownerKey = "neriva-offline-owner";

export async function syncOfflineCards(owner: string, signal: AbortSignal) {
  const reviewedKey = `neriva-offline-reviewed:${owner}`;
  const syncedKey = `neriva-offline-synced:${owner}`;
  try {
    const reviewed: string[] = JSON.parse(localStorage.getItem(reviewedKey) || "[]");
    const synced = new Set<string>(JSON.parse(localStorage.getItem(syncedKey) || "[]"));
    const pending = [...new Set(reviewed)].filter((id) => typeof id === "string" && !synced.has(id));
    for (let offset = 0; offset < pending.length && !signal.aborted; offset += 100) {
      const batch = pending.slice(offset, offset + 100);
      const digests = await Promise.all(batch.map(async (id) => {
        const hash = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(id));
        return Array.from(new Uint8Array(hash), (byte) => byte.toString(16).padStart(2, "0")).join("");
      }));
      if (signal.aborted || localStorage.getItem(ownerKey) !== owner) return;
      await api("/api/offline-reviews", { method: "POST", body: { card_ids: digests }, signal });
      for (const id of batch) synced.add(id);
      localStorage.setItem(syncedKey, JSON.stringify([...synced]));
    }
  } catch {
    // Keep pending local marks for the next authenticated session or reconnect.
  }
}

export function saveOfflineCards(owner: string, cards: Card[], replaceOnly = false) {
  if (!owner) return;
  try {
    const key = `neriva-offline-cards:${owner}`;
    const selected = replaceOnly ? JSON.parse(localStorage.getItem(key) || "[]") as Card[] : cards;
    const allowed = new Set(cards.map((card) => card.id));
    localStorage.setItem(key, JSON.stringify(selected.filter((card) => card.answer && allowed.has(card.id)).slice(0, 100)));
    localStorage.setItem(ownerKey, owner);
  } catch {
    return;
  }
}

export function selectOfflineOwner(owner: string) {
  try {
    if (owner && localStorage.getItem(`neriva-offline-cards:${owner}`)) localStorage.setItem(ownerKey, owner);
    else localStorage.removeItem(ownerKey);
  } catch {
    return;
  }
}

export function clearOfflineOwner() {
  try {
    localStorage.removeItem(ownerKey);
  } catch {
    return;
  }
}

function readOfflineCards() {
  try {
    const owner = localStorage.getItem(ownerKey);
    const cards = owner ? JSON.parse(localStorage.getItem(`neriva-offline-cards:${owner}`) || "[]") : [];
    const key = `neriva-offline-reviewed:${owner}`;
    const attempts = JSON.parse(localStorage.getItem(key) || "[]");
    return {
      key,
      cards: Array.isArray(cards) ? cards.filter((card): card is Card => typeof card?.id === "string" && typeof card.title === "string" && typeof card.answer === "string") : [],
      reviewed: Array.isArray(attempts) ? attempts.filter((id): id is string => typeof id === "string") : [],
    };
  } catch {
    return { key: "", cards: [], reviewed: [] };
  }
}

export function OfflineReview({ russian }: { russian: boolean }) {
  const [deck] = useState(readOfflineCards);
  const [reviewed, setReviewed] = useState(deck.reviewed);
  const [index, setIndex] = useState(0);
  const [revealed, setRevealed] = useState(false);
  const [saveFailed, setSaveFailed] = useState(false);
  const card = deck.cards[index];
  if (!card) return null;

  function advance(remembered: boolean) {
    if (remembered && !reviewed.includes(card.id)) {
      const next = [...reviewed, card.id];
      try {
        localStorage.setItem(deck.key, JSON.stringify(next));
      } catch {
        setSaveFailed(true);
        return;
      }
      setReviewed(next);
    }
    setSaveFailed(false);
    setRevealed(false);
    setIndex((index + 1) % deck.cards.length);
  }

  return (
    <div className="offline-review">
      <p>{russian ? "Отработано на этом устройстве" : "Reviewed on this device"}: {reviewed.filter((id) => deck.cards.some((item) => item.id === id)).length}</p>
      <h2 dir="auto">{card.title}</h2>
      {revealed ? (
        <>
          <p dir="auto">{card.answer}</p>
          <Button type="button" onClick={() => advance(true)}>{russian ? "Вспомнил" : "I remembered"}</Button>{" "}
          <Button type="button" variant="outline" onClick={() => advance(false)}>{russian ? "Повторить позже" : "Review later"}</Button>
        </>
      ) : <Button type="button" onClick={() => setRevealed(true)}>{russian ? "Показать ответ" : "Show answer"}</Button>}
      {saveFailed && <p role="status">{russian ? "Не удалось сохранить результат." : "Could not save the result."}</p>}
    </div>
  );
}
