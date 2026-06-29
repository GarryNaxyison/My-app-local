import { mkdir } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { chromium } from "playwright";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const outputDir = path.resolve(__dirname, "../public/assets/seo");

const variants = [
  {
    filename: "poliglot-ai-og-ru.jpg",
    lang: "ru",
    kicker: "AI-репетитор в Telegram и web app",
    title: "Poliglot AI",
    subtitle: "Уроки, разговорная практика, голос, фото-перевод и ошибки в одном профиле.",
    chips: ["AI tutor", "Telegram", "Voice Coach"],
    compare: ["vs словарное приложение", "vs языковой бот", "Telegram + web app"],
  },
  {
    filename: "poliglot-ai-og-en.jpg",
    lang: "en",
    kicker: "AI language tutor in Telegram and web app",
    title: "Poliglot AI",
    subtitle: "Lessons, speaking practice, voice, photo translation, mistakes, and progress in one profile.",
    chips: ["AI tutor", "Telegram", "Voice Coach"],
    compare: ["vs vocabulary app", "vs language bot", "Telegram + web app"],
  },
];

function escapeHtml(value) {
  return String(value)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

function renderTemplate(variant) {
  return `<!doctype html>
<html lang="${escapeHtml(variant.lang)}">
  <head>
    <meta charset="utf-8" />
    <style>
      * { box-sizing: border-box; }
      body {
        margin: 0;
        width: 1200px;
        height: 630px;
        overflow: hidden;
        background:
          radial-gradient(circle at 18% 18%, rgba(143, 241, 208, 0.34), transparent 24%),
          radial-gradient(circle at 82% 26%, rgba(45, 91, 255, 0.36), transparent 28%),
          linear-gradient(135deg, #050914 0%, #0b1d2d 46%, #08111f 100%);
        color: #f8fbff;
        font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      }
      .frame {
        position: relative;
        display: grid;
        grid-template-columns: 1.05fr 0.95fr;
        gap: 42px;
        width: 100%;
        height: 100%;
        padding: 64px;
      }
      .frame::before {
        content: "";
        position: absolute;
        inset: 30px;
        border: 1px solid rgba(255, 255, 255, 0.13);
        border-radius: 28px;
        pointer-events: none;
      }
      .copy {
        position: relative;
        z-index: 1;
        display: flex;
        min-width: 0;
        flex-direction: column;
        justify-content: space-between;
      }
      .kicker {
        width: fit-content;
        max-width: 100%;
        border: 1px solid rgba(143, 241, 208, 0.34);
        border-radius: 999px;
        padding: 12px 18px;
        color: #8ff1d0;
        font-size: 24px;
        font-weight: 800;
      }
      h1 {
        margin: 42px 0 0;
        font-size: 118px;
        line-height: 0.88;
        letter-spacing: 0;
      }
      p {
        max-width: 730px;
        margin: 30px 0 0;
        color: rgba(238, 247, 255, 0.84);
        font-size: 35px;
        line-height: 1.14;
      }
      .chips {
        display: flex;
        flex-wrap: wrap;
        gap: 12px;
      }
      .chips span,
      .compare span {
        border: 1px solid rgba(255, 255, 255, 0.13);
        border-radius: 999px;
        padding: 12px 16px;
        background: rgba(255, 255, 255, 0.08);
        color: rgba(247, 251, 255, 0.9);
        font-size: 22px;
        font-weight: 750;
      }
      .panel {
        position: relative;
        z-index: 1;
        align-self: center;
        border: 1px solid rgba(255, 255, 255, 0.14);
        border-radius: 26px;
        padding: 28px;
        background: rgba(255, 255, 255, 0.08);
        box-shadow: 0 36px 90px rgba(0, 0, 0, 0.32);
      }
      .panel-top {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 18px;
        margin-bottom: 26px;
      }
      .panel-top strong {
        font-size: 28px;
      }
      .score {
        border-radius: 999px;
        padding: 10px 14px;
        background: rgba(143, 241, 208, 0.16);
        color: #8ff1d0;
        font-size: 22px;
        font-weight: 800;
      }
      .compare {
        display: grid;
        gap: 16px;
      }
      .compare span {
        display: flex;
        justify-content: space-between;
        border-radius: 18px;
        padding: 20px;
        font-size: 25px;
      }
      .compare span::after {
        content: "Poliglot AI";
        color: #8ff1d0;
      }
    </style>
  </head>
  <body>
    <main class="frame">
      <section class="copy">
        <div>
          <div class="kicker">${escapeHtml(variant.kicker)}</div>
          <h1>${escapeHtml(variant.title)}</h1>
          <p>${escapeHtml(variant.subtitle)}</p>
        </div>
        <div class="chips">${variant.chips.map((chip) => `<span>${escapeHtml(chip)}</span>`).join("")}</div>
      </section>
      <aside class="panel">
        <div class="panel-top">
          <strong>Compare</strong>
          <span class="score">AEO ready</span>
        </div>
        <div class="compare">${variant.compare.map((item) => `<span>${escapeHtml(item)}</span>`).join("")}</div>
      </aside>
    </main>
  </body>
</html>`;
}

await mkdir(outputDir, { recursive: true });
const browser = await chromium.launch({ headless: true });
try {
  const page = await browser.newPage({ viewport: { width: 1200, height: 630 }, deviceScaleFactor: 1 });
  for (const variant of variants) {
    await page.setContent(renderTemplate(variant), { waitUntil: "load" });
    await page.screenshot({
      path: path.join(outputDir, variant.filename),
      type: "jpeg",
      quality: 92,
    });
    console.log(`wrote ${variant.filename}`);
  }
} finally {
  await browser.close();
}
