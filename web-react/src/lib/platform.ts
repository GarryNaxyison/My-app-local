import { useEffect, useState } from "react";

/**
 * Единый источник истины для разделения ПК- и мобильной версий веб-приложения.
 *
 * Платформа определяется ТОЛЬКО шириной вьюпорта (breakpoint 760px) —
 * ровно так же, как мобильные CSS-правила `@media (max-width: 760px)` в app.css.
 * Это устраняет рассинхрон JS ↔ CSS на планшетах 761–1180px и тач-ноутбуках:
 * если CSS показывает десктопную раскладку, JS тоже считает устройство десктопом.
 *
 * - `usePlatform()` — React-хук, реактивно возвращает "mobile" | "desktop".
 * - `useIsMobile()` — булев вариант.
 * - `syncPlatformAttribute()` — выставляет `data-platform` на <html> для CSS-хуков.
 */

export type PlatformId = "mobile" | "desktop";

export const MOBILE_LAYOUT_MAX_WIDTH = 760;
export const MOBILE_LAYOUT_QUERY = `(max-width: ${MOBILE_LAYOUT_MAX_WIDTH}px)`;
export const REDUCED_MOTION_QUERY = "(prefers-reduced-motion: reduce)";

export function getPlatformId(): PlatformId {
  if (typeof window === "undefined" || typeof window.matchMedia !== "function") return "desktop";
  if (window.matchMedia(MOBILE_LAYOUT_QUERY).matches) return "mobile";
  return window.innerWidth <= MOBILE_LAYOUT_MAX_WIDTH ? "mobile" : "desktop";
}

export function isMobilePlatform(): boolean {
  return getPlatformId() === "mobile";
}

export function prefersReducedMotion(): boolean {
  if (typeof window === "undefined" || typeof window.matchMedia !== "function") return false;
  return window.matchMedia(REDUCED_MOTION_QUERY).matches;
}

function useMediaState(queryText: string, fallback: () => boolean): boolean {
  const [state, setState] = useState(fallback);
  useEffect(() => {
    const query = window.matchMedia(queryText);
    const update = () => setState(fallback());
    update();
    query.addEventListener("change", update);
    window.addEventListener("resize", update);
    return () => {
      query.removeEventListener("change", update);
      window.removeEventListener("resize", update);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [queryText]);
  return state;
}

export function usePlatform(): PlatformId {
  return useMediaState(MOBILE_LAYOUT_QUERY, () => getPlatformId() === "mobile") ? "mobile" : "desktop";
}

export function useIsMobile(): boolean {
  return usePlatform() === "mobile";
}

export function usePrefersReducedMotion(): boolean {
  return useMediaState(REDUCED_MOTION_QUERY, prefersReducedMotion);
}

/** Проставляет data-platform="mobile|desktop" на <html> и реактивно обновляет. */
export function usePlatformAttribute(): PlatformId {
  const platform = usePlatform();
  useEffect(() => {
    document.documentElement.dataset.platform = platform;
  }, [platform]);
  return platform;
}
