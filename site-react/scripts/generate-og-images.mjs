import { mkdir, readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { chromium } from "playwright";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const outputDir = path.resolve(__dirname, "../public/assets/seo");
const wordmark = await readFile(path.resolve(__dirname, "../public/assets/brand/neriva-wordmark.png"));
const wordmarkSrc = `data:image/png;base64,${wordmark.toString("base64")}`;

const variants = [
  {
    filename: "neriva-social-preview-v3-ru.jpg",
    lang: "ru",
    kicker: "AI-репетитор в Telegram и web app",
    title: "NERIVA",
    subtitle: "Уроки, разговорная практика, голос, фото-перевод и ошибки в одном профиле.",
    signals: [
      ["AI-уроки", "Короткие задания и объяснение ошибок"],
      ["Speaking", "Диалоги, roleplay и произношение"],
      ["Telegram + web app", "Один профиль и общий прогресс"],
    ],
  },
  {
    filename: "neriva-social-preview-v3-en.jpg",
    lang: "en",
    kicker: "AI language tutor in Telegram + web app",
    title: "NERIVA",
    subtitle: "Lessons, speaking practice, voice, photo translation, mistakes, and progress in one profile.",
    signals: [
      ["AI lessons", "Short tasks and mistake explanations"],
      ["Speaking", "Dialogues, roleplay, and pronunciation"],
      ["Telegram + web app", "One profile with shared progress"],
    ],
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
          radial-gradient(circle at 84% 20%, rgba(139, 120, 222, 0.38), transparent 40%),
          linear-gradient(135deg, #0d1026 0%, #191c39 100%);
        color: #f8fbff;
        font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      }
      .frame {
        position: relative;
        display: grid;
        grid-template-columns: 1fr 0.9fr;
        gap: 48px;
        width: 100%;
        height: 100%;
        padding: 56px 64px;
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
        justify-content: center;
      }
      .kicker {
        width: fit-content;
        max-width: 100%;
        margin-top: 26px;
        color: #f3bc7b;
        font-size: 22px;
        font-weight: 650;
      }
      h1 {
        margin: 0;
        width: min(100%, 440px);
        padding: 20px 22px;
        border-radius: 18px;
        background: #fff;
      }
      h1 img {
        display: block;
        width: 100%;
        height: auto;
      }
      p {
        max-width: 730px;
        margin: 20px 0 0;
        color: rgba(238, 247, 255, 0.84);
        font-size: 28px;
        line-height: 1.35;
      }
      .signal {
        border: 1px solid rgba(255, 255, 255, 0.13);
        background: rgba(255, 255, 255, 0.08);
        color: rgba(247, 251, 255, 0.9);
        font-weight: 750;
      }
      .panel {
        position: relative;
        z-index: 1;
        align-self: center;
        border: 1px solid rgba(255, 255, 255, 0.14);
        border-radius: 24px;
        padding: 26px;
        background: rgba(255, 255, 255, 0.08);
        box-shadow: 0 36px 90px rgba(0, 0, 0, 0.32);
      }
      .panel h2 {
        margin: 0;
        font-size: 34px;
        line-height: 1;
        letter-spacing: 0;
      }
      .panel > p {
        margin: 12px 0 22px;
        color: rgba(238, 247, 255, 0.75);
        font-size: 24px;
        line-height: 1.18;
      }
      .signals {
        display: grid;
        gap: 12px;
      }
      .signal {
        border-radius: 18px;
        padding: 15px 18px;
      }
      .signal strong {
        display: block;
        color: #c2c5ff;
        font-size: 22px;
      }
      .signal span {
        display: block;
        margin-top: 6px;
        color: rgba(247, 251, 255, 0.82);
        font-size: 18px;
        line-height: 1.2;
      }
    </style>
  </head>
  <body>
    <main class="frame">
      <section class="copy">
        <div>
          <h1><img src="${wordmarkSrc}" alt="${escapeHtml(variant.title)}" width="1800" height="329"></h1>
          <div class="kicker">${escapeHtml(variant.kicker)}</div>
          <p>${escapeHtml(variant.subtitle)}</p>
        </div>
      </section>
      <aside class="panel">
        <h2>${escapeHtml(variant.lang === "ru" ? "AI-репетитор иностранных языков" : "AI language tutor")}</h2>
        <p>${escapeHtml(variant.lang === "ru" ? "Практика языка в одном профиле." : "Language practice in one profile.")}</p>
        <div class="signals">${variant.signals.map(([title, body]) => `<div class="signal"><strong>${escapeHtml(title)}</strong><span>${escapeHtml(body)}</span></div>`).join("")}</div>
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
