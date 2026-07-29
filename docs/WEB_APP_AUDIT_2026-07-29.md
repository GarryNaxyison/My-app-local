# Аудит веб-приложения NERIVA (web-react) — ПК и мобильная версия

Дата: 2026-07-29
Объект: `web-react/` (React 19 + Vite 7 + TS), собранный вывод в `web/`.
Методы: статический анализ кода, сборка `tsc -b && vite build`, прогон Playwright e2e (Desktop Chrome 1440×900 + Pixel 5 390×844).

---

## Итог

| Проверка | Результат |
|---|---|
| Сборка TypeScript + Vite | ✅ чисто, без ошибок |
| E2E-тесты | 169 passed / **3 failed** / 44 skipped (skip — платформенные пары, норма) |
| Адаптивность CSS | хорошая база (dvh, safe-area, media 1180/900/760/420px) |
| Критичные баги | **1** (ломает UX урока и 3 e2e-теста) |
| Важные проблемы | **5** |
| Техдолг / незначительное | **4** |

---

## 🔴 Критично

### 1. Двойная кнопка «Следующий урок» после ответа (ПК + мобайл)

**Файл:** `web-react/src/App.tsx:7303` и `web-react/src/App.tsx:7355`

После отправки ответа в уроке одновременно выполняются два условия:
- `lessonHasCompletedAnswer === true` → кнопка `.lesson-panel-next-v2` рендерится под выводом чата (строка 7303);
- `lessonComposerLocked === true` → рендерится заблокированная панель композера, в которой есть **та же** кнопка `.lesson-panel-next-v2` (строка 7355).

Итог: на экране две одинаковые кнопки «Следующий урок» подряд. Падают e2e-тесты:
- `desktop-chromium › regression: lesson tab keeps one active task until the learner submits`
- `mobile-chromium › regression: lesson tab keeps one active task until the learner submits`
- `mobile-chromium › mobile lesson keeps output readable and phrase save inside input controls`

**Фикс:** оставить один рендер. Рекомендуется убрать кнопку из области вывода (строки 7302–7307) — вариант в locked-панели информативнее (есть подпись «Ответ принят…») и на мобильных ближе к большому пальцу.

---

## 🟠 Важно

### 2. Нет `viewport-fit=cover` — safe-area не работает на iPhone (мобайл)

**Файл:** `web-react/index.html:5` (и собранный `web/index.html:5`)

```html
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
```

Весь CSS активно использует `env(safe-area-inset-top/bottom)` (44 места: нижняя навигация, модалки, композеры). Но без `viewport-fit=cover` iOS Safari отдаёт эти переменные как `0`. Итог на iPhone с «чёлкой»/home-индикатором:
- нижняя мобильная навигация уходит под home-индикатор;
- контент модалок прижимается к краям под notch.

**Фикс:**
```html
<meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover" />
```

### 3. Авто-зум iOS Safari на форме входа (мобайл)

**Файл:** `web-react/src/components/ui/sign-in.tsx:109, 117, 131, 145`

Поля логина/пароля/повтора/реф-кода используют Tailwind `text-sm` = **14px**. iOS Safari автоматически зумирует страницу при фокусе на input с font-size < 16px и не возвращает зум обратно — раскладка «ломается» до ручного действия пользователя.

**Фикс:** для этих полей на мобильных ставить `text-base` (16px), например `text-base md:text-sm`. Поля композеров и OTP в порядке (16px/18px).

### 4. Three.js/WebGL грузится всегда и везде (мобайл, производительность)

**Файлы:** `web-react/src/App.tsx:104` (статический импорт `sign-in`), `web-react/src/components/ui/sign-in.tsx:3,97`, `web-react/vite.config.ts:78-84`

- Чанк `canvas` (three + @react-three/fiber) = **879 KB (240 KB gzip)** — подключён через `<link rel="modulepreload">` в собранном `web/index.html` и парсится при каждом заходе, даже у авторизованных пользователей.
- WebGL-сцена `AuthGenerativeScene` рендерится на странице входа **включая мобильные** — лишняя нагрузка на GPU и батарею.
- Общий вес первой загрузки ≈ 2.5 MB (≈730 KB gzip): main 1.36 MB, canvas 879 KB, motion 137 KB, css 276 KB.

**Фикс:** `React.lazy`/`Suspense` для `SignInPage` (или хотя бы для `AuthGenerativeScene`); на мобильных (`pointer: coarse` / `max-width: 760px`) заменять WebGL-сцену статичным CSS-градиентом (класс-заготовка `paper-shader-bg--*` уже есть).

### 5. Рассинхрон детекта «мобильности» JS vs CSS (планшеты, тач-ноутбуки)

**Файлы:** `web-react/src/App.tsx:789-815` (`useMobileUiLayout`), `web-react/src/styles/app.css`

JS считает устройство мобильным по `(max-width: 760px)` **ИЛИ** `pointer: coarse` **ИЛИ** UA (`Android|iPhone|iPad|Mobile`). CSS-правила — только по `@media (max-width: 760px)`.

На планшете 761–1180px с тачем (iPad, Android-планшет, тач-ноутбук без мыши) JS добавит мобильные классы (`roleplay-view-v2--mobile-session`, `pronunciation-dashboard-v2--mobile`), но часть CSS-правил, спрятанных в `@media (max-width: 760px)` (например, `app.css:9416-9466` для roleplay mobile-session), не применится → промежуточная, не проверенная раскладка.

**Фикс:** единый источник истины — либо убрать `pointer: coarse`/UA из JS-хука (оставить только 760px), либо вынести мобильные CSS-правила из media-query в классы.

### 6. Нет поддержки `prefers-reduced-motion` (доступность, батарея)

Ни в CSS, ни в framer-motion не обрабатывается `prefers-reduced-motion`: орбиты, шейдеры, анимации крутятся всегда. Это и доступность (вестибулярные расстройства), и расход батареи на мобильных.

**Фикс:** медиа-запрос `@media (prefers-reduced-motion: reduce)` для отключения keyframe-анимаций + `useReducedMotion()` из framer-motion.

---

## 🟡 Незначительно / техдолг

7. **Мёртвый код:** `src/components/CanvasRevealEffect.tsx`, `src/components/ui/background-paper-shaders.tsx`, `src/components/ui/shaders-hero-section.tsx` нигде не импортируются. Удалить (и убрать `shaders`/`canvas` из `manualChunks`, если останутся не нужны).
8. **`width: 100dvw`** у `.sign-in-page-v2` (`app.css:4358`) — на десктопе со скроллбаром даёт горизонтальное переполнение (сейчас маскируется `overflow-x: hidden` на body). Лучше `width: 100%`.
9. **Монолит:** `App.tsx` — 9 227 строк, один main-бандл 1.36 MB. Стоит разбить по views и подключать их через `React.lazy` — ускорит первую загрузку на мобильных.
10. **44 skipped в e2e** — это платформенные пары `test.skip(isMobile)`/`skip(!isMobile)`, штатно; просто иметь в виду при чтении отчётов.

---

## ✅ Что работает хорошо

- **Сборка чистая:** `tsc -b` без ошибок типов, Vite build ~4.4 s.
- **169 e2e зелёные** на обоих вьюпортах: навигация, словарь, ролевые, произношение, аудирование, оффлайн, платежи, настройки, reorder навигации и т.д.
- **CSS-адаптив добротный:** `100dvh` с `100vh`-фолбэком; `min-width: 320px`; `overflow-x: hidden`; медиа-точки 1180/900/760/420px; `overscroll-behavior`, `touch-action`, `-webkit-overflow-scrolling: touch` в скроллящихся областях.
- **Мобильная нижняя навигация:** sticky с safe-area, touch-targets 54px, long-press + drag-to-reorder с pointer capture, автоскролл, серверная персистенция раскладки.
- **Композеры** (урок/практика/ролевая): textarea наследует 16px — авто-зума iOS нет; авто-ресайз; Enter=send, Shift+Enter=новая строка.
- **OTP-инпуты** 18px — ок.
- **PWA:** корректный `manifest.webmanifest`, service worker `offline-deck-sw.js` с аккуратным controllerchange-reload.
- **Модалки на мобильных:** `place-items: start center`, прокрутка, паддинги с safe-area сверху/снизу.

---

## Рекомендуемый порядок исправлений

1. Убрать дубль `.lesson-panel-next-v2` (App.tsx:7302–7307) → вернуть 3 падающих теста в зелень.
2. Добавить `viewport-fit=cover` в `web-react/index.html`.
3. Поднять font-size полей входа до 16px на мобильных (sign-in.tsx).
4. Ленивая загрузка WebGL-сцены + CSS-фолбэк на мобильных.
5. Унифицировать детект мобильности (JS ↔ CSS).
6. `prefers-reduced-motion`.
7. Удалить мёртвые компоненты и нарезать code-splitting.

> Примечание: в ходе аудита была выполнена сборка `npm run build` — артефакты в `web/` (index.html + assets) обновлены до актуального состояния исходников; изменений исходного кода не вносилось.
