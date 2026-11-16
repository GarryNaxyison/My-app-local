# Changelog

## 2026-11-16 — Правки оформления и экрана «Сегодня»

- Убрана белая подложка под прозрачным логотипом и знаком на сайте. Подключён фирменный фавикон по отдельному адресу.
- Возвращена вращающаяся планетарная анимация первого экрана. При системном уменьшении движения фон остаётся статичным.
- Сокращён верхний блок «Сегодня». Действие урока, прогресс и задания остаются на виду, планы и календарь раскрываются по запросу. Витрина обновлена двумя снимками этого экрана.

## 2026-11-16 — Фирменное оформление NERIVA и удобство уроков

- Официальный логотип добавлен в веб-приложение, публичный сайт, иконки установки, страницы базы знаний и превью ссылок. Применена фирменная палитра.
- На сайте показаны 10 актуальных экранов приложения с учебными примерами. Галерея поддерживает масштабирование с клавиатуры и возврат фокуса после закрытия.
- В мобильном каталоге вместо повторов и двоеточий отображаются содержательные пояснения. В длинном уроке кнопка «К ответу» переводит к полю ввода. Подпись правильного варианта переведена на русский.
- Установлены frontend-design, Playwright, screenshot и security-best-practices. Дизайн и браузерные проверки применены к обновлению.
- Статические файлы опубликованы без перезапуска сервера. Сравнение базы подтвердило сохранение прогресса 66 профилей. Проверки описаны в [спецификации](docs/specs/official-brand/spec.md).

## 2026-11-14 — Публикация учебного интерфейса

- На neriva.ru опубликованы восемь файлов веб-сборки: ассеты, service worker и HTML. Ассеты размещены первыми, HTML заменён атомарно.
- Серверный бинарный файл, процесс приложения, настройки, изображения и базы не изменялись. Сравнение записей пользователей, слов, ошибок, привычек, бонусов и активных Tutor-сессий до и после публикации подтвердило полное совпадение; в основной базе 65 профилей.
- Уточнена инструкция деплоя: исторический сброс прогресса запрещён для обновления интерфейса.
- На опубликованной версии проверены контрольные суммы семи статических файлов, экран входа Chromium 390/1440 px, healthz/приложение/лендинг и наличие 38 клиентских API-маршрутов без изменения пользовательских данных.

## 2026-11-14 — Адаптивное учебное пространство NERIVA

- Горизонтальная лента функций заменена компактной боковой навигацией на ПК; на телефоне добавлены четыре основных перехода и встроенный каталог с поиском.
- Главная получила заметный вход в занятие. Сохранённые настройки порядка меню читаются без перезаписи.
- Учебный материал и ответ следуют друг за другом в общей прокрутке страницы. Убраны вложенные рамки и отдельная прокрутка чата.
- План AI Tutor сворачивается; завершённые уроки открываются в основной области с прежней пагинацией и повторным запуском.
- Сохранены сгенерированные визуальные ассеты, учебные API, бизнес-правила, прогресс и черновики. Обновлена версия кэша оболочки PWA.
- Критерии и план проверки: `docs/specs/learning-workspace/spec.md`.
- Проверено: сборка React; полный Playwright — 198 успешно и 52 условных пропуска; адресные Go-проверки оболочки, PWA, ассетов и API-контрактов. Изображения и серверное хранилище прогресса не менялись.

## 2026-11-14

- Fixed stale AI Tutor web answers by sending and checking the expected stage; serialized session submissions within the process.
- Committed SQLite lesson completion, XP, counters and habit progress in one transaction. Added rollback/replay tests; kept JSON completion under its store lock.
- Made post-lesson quality assessment best effort with a five-second deadline after completion; included stored learner answers in its input.
- Released microphone tracks on navigation and late permission resolution, handled recorder setup errors, and used a localized native recording button.
- Preserved tutor drafts in session storage, suppressed duplicate actions, retained pending status for overlapping actions, reset explicit retry timing, and added an in-memory timer fallback.
- Prevented automatic PWA reload during input or pending operations. React build, full Go tests and 183 Playwright cases passed; 51 platform skips. Two extra PWA regression cases passed separately.

- Added a desktop/mobile web audit with code findings, UI recommendations, AI Tutor logic proposals, and explicit verification limits in `WEB_APP_AUDIT_2026-11-14.md`.
- Added the owner's current requests in root `USER_REQUESTS.md`. Installed three React/UI skills in the user Codex directory; application behavior was not changed.

Project change history is maintained in [`docs/tracking/CHANGELOG.md`](docs/tracking/CHANGELOG.md), alongside the user request log and release audit.

## 2026-10-17

- Clarified AI Tutor preparation status, added an elapsed timer that survives page refresh, and extended the lesson-start timeout to three minutes.
- Enabled AI Tutor onboarding with no pre-learned words; successful lessons seed deduplicated single-word entries into the shared SQLite dictionary.
- Replaced misleading dictionary-based lesson loading copy and improved responsive desktop/mobile layouts across learning surfaces.
- Switched text generation to OpenRouter `openai/gpt-6-luna` and pronunciation audio assessment to the audio-capable, lower-token-rate `google/gemini-3-flash-preview` model.
- Reset all users' learning progress once while preserving all 57 user profiles and the 15 web accounts with their login/profile settings; deployed and verified service health and low host load.
