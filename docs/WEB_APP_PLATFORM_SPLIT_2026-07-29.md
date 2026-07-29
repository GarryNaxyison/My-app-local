# Разделение ПК- и мобильной версий NERIVA Web

Дата: 2026-07-29
Статус: внедрено
Связанные документы: [WEB_APP_AUDIT_2026-07-29.md](WEB_APP_AUDIT_2026-07-29.md)

---

## Принцип: единый источник истины

Платформа определяется **только шириной вьюпорта** — breakpoint **760px**,
одинаково в JS и CSS. Никаких `pointer: coarse` / UA-эвристик в JS:
именно они раньше давали рассинхрон на планшетах 761–1180px и тач-ноутбуках
(JS добавлял мобильные классы, а CSS-правила в `@media (max-width: 760px)` не применялись).

```
ширина ≤ 760px  →  mobile   (нижняя навигация, quick-контролы, мобильные раскладки)
ширина > 760px  →  desktop  (верхняя FunctionRibbon, десктопные сетки)
```

## Архитектура

### JS: `web-react/src/lib/platform.ts`

| Экспорт | Назначение |
|---|---|
| `getPlatformId()` | `"mobile" \| "desktop"` по matchMedia/innerWidth |
| `usePlatform()` | Реактивный хук платформы (matchMedia + resize) |
| `useIsMobile()` | Булев вариант |
| `usePlatformAttribute()` | Выставляет `data-platform` на `<html>` и возвращает платформу |
| `usePrefersReducedMotion()` | Реактивный `prefers-reduced-motion` |

`useMobileUiLayout()` в `App.tsx` теперь просто делегирует `usePlatform()` —
все существующие вызовы (roleplay, pronunciation) автоматически синхронизированы с CSS.

### Условный рендер платформенной навигации (App.tsx)

- **desktop**: `FunctionRibbon` (верхняя лента функций)
- **mobile**: `MobileBottomNav` (нижняя навигация) + `MobileQuickControls` (на home)

Компоненты другой платформы **не монтируются** в DOM вовсе (раньше рендерились
и скрывались через `display: none` — лишняя работа эффектов, localStorage, pointer listeners).

Корневой шелл помечен классом `v2-shell--mobile` / `v2-shell--desktop`
и атрибутом `data-platform` — для платформенных CSS-переопределений.

### CSS

- `src/styles/app.css` — общая база + существующие `@media (max-width: 760px)` блоки
  (они и есть мобильный слой, синхронизирован с JS по breakpoint).
- `src/styles/platform.css` — платформенный слой (подключается после app.css):
  - статичная сцена входа (замена WebGL на мобильных/reduced-motion);
  - глобальный `@media (prefers-reduced-motion: reduce)`;
  - хуки `html[data-platform="..."]` для правил, следующих за JS-определением платформы.

### URL-маршруты

- `/app` — универсальный вход, адаптируется под устройство без редиректа.
- `/app/mobile`, `/app/desktop` — явные платформенные входы; если путь не совпадает
  с текущей платформой, нормализуется через `history.replaceState` (без перезагрузки).
  Сервер отдаёт SPA на `/app/` catch-all — дополнительной серверной логики не требуется.

## Производительность

- Страница входа (`SignInPage`) + WebGL/three.js (~879 KB) вынесены в ленивый чанк:
  `React.lazy` в `App.tsx` → `sign-in.tsx`, внутри — ленивый `AuthGenerativeScene`.
  На мобильных и при reduced-motion WebGL не запрашивается вообще — рендерится
  статичный CSS-градиент `.sign-in-page-v2__static-scene`.
- Удалён мёртвый код: `CanvasRevealEffect.tsx`, `background-paper-shaders.tsx`,
  `shaders-hero-section.tsx` и пакет `@paper-design/shaders-react` (чанк `shaders`).

### Важно: ref-эффекты внутри ленивых поддеревьев

Ленивая `SignInPage` сломала рендер Turnstile-капчи: слот капчи передаётся
внутрь ленивого поддерева (`captchaSlot`), поэтому монтируется после первого
коммита `AuthStandaloneView`. Эффект с `useRef` срабатывал раньше появления
элемента (`ref.current === null`) и не перезапускался — капча не рендерилась,
логин блокировался. Лечится **callback-ref + состоянием**
(`ref={setCaptchaEl}`, элемент в deps эффекта) — эффект перезапускается,
когда слот реально появляется в DOM. Правило: любой эффект, зависящий от
DOM-элемента внутри `React.lazy`-поддерева, должен использовать callback-ref.

## Доступность

- `prefers-reduced-motion`: CSS-анимации отключены глобально (platform.css),
  framer-motion — через `MotionConfig reducedMotion="user"` (main.tsx),
  WebGL-сцена заменяется статикой (sign-in.tsx), бут-экран — inline media-query в index.html.
- `viewport-fit=cover` — safe-area инсеты работают на iPhone (нижняя навигация выше home-индикатора).
- Поля формы входа — 16px на мобильных (`text-base md:text-sm`) — нет авто-зума iOS.

## Проверка

- `npm run build` — чисто (tsc + vite).
- `npx playwright test` (Desktop Chrome 1440×900 + Pixel 5 390×844):
  **172 passed / 0 failed / 44 skipped** (skip — платформенные пары по `test.skip(isMobile)`).
  Для сравнения до правок: 169 passed / 3 failed.
- Два e2e-теста обновлены под новую архитектуру (это осознанное изменение поведения,
  а не ослабление проверок):
  - `AI Tutor server-driven lesson…` — селектор навигации стал платформенным
    (`.function-ribbon` на ПК / `.mobile-bottom-nav-v2` на мобильных), как в соседних тестах;
  - `auth login page uses React sign-in component…` — на мобильных ожидается статичная
    сцена `.sign-in-page-v2__static-scene` вместо WebGL-canvas; на ПК canvas ожидается
    через `toHaveCount(1)` с авто-ожиданием ленивого чанка.
