# Документация лендинга NERIVA

## Запуск

```bash
npm start
```

По умолчанию сайт открывается на `http://127.0.0.1:3000/poliglot-ai.html`.

## Основные файлы

- `poliglot-ai.html` — основной лендинг, стили, HTML-разметка, hero-демо, блоки возможностей, шагов, тарифов, отзывов, CTA и FAQ.
- `terms.html`, `privacy.html` — legal-страницы в общем брендовом стиле, с тёмной темой, переключателем языка и общими ассетами.
- `assets/site-i18n.js` — переключатель языка и словарь переводов для сайта.
- `assets/site-phrases.js` — расширенная база фраз для 20 языков.
- `assets/app-background-light.png`, `assets/app-background-dark.png`, `assets/brand-logo-hero-*.png`, `assets/header-*.png`, `assets/icon-*.png`, `assets/plan-*.png`, `assets/award-*.png` — визуальные ассеты нового брендированного лендинга.
- `server.mjs` — простой Node.js сервер для локального запуска и отдачи статических файлов.

## Локализация

Переключатель языка использует `data-site-language-select`. Статические строки переводятся через `data-i18n` или через словарь фраз, а обычные текстовые узлы и атрибуты `alt`, `aria-label`, `title`, `placeholder`, `content` проходят через `assets/site-i18n.js`.

После изменения текстов публичных страниц нужно пересобрать расширенный словарь из корня общего проекта:

```bash
node tools/rebuild_public_site_translations.mjs
```

Скрипт собирает фразы из `poliglot-ai.html`, `terms.html`, `privacy.html` и обновляет `assets/site-phrases.js`.

## Тема и кеш

Тёмная тема включается кнопкой `[data-theme-toggle]`. Выбор сохраняется в `localStorage` под ключом `poliglot_site_theme`; при первом открытии используется системная тема браузера.

Публичные страницы используют версионированные ссылки на ассеты и словари через `?v=20260521-brand3`. Для `.ru`/`.online` в Caddy-конфиге HTML отдаётся с `Cache-Control: no-cache, no-store, must-revalidate`, словари `site-i18n.js` и `site-phrases.js` — с `no-cache, must-revalidate`, а остальные ассеты могут кешироваться долго, потому что их URL меняется при версии.

## Верстка тарифов

Тарифы показываются карточками `plan-card` внутри `pricing-grid`: Free, Premium и Platinum. Premium и Platinum показывают старую цену с зачёркиванием, новую цену `300 ₽`/`590 ₽`, бейдж `70% экономии` и отметку `1 месяц запуска`. Детальные дневные лимиты вынесены в компактную таблицу `limits`, чтобы цены и ограничения читались без перегруженной сравнительной простыни.
