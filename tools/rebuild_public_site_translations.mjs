import fs from "node:fs";
import path from "node:path";
import vm from "node:vm";

const ROOT = process.cwd();
const SITE_DIR = path.join(ROOT, "Сайт полиглота для бота");
const PHRASES_PATH = path.join(SITE_DIR, "assets", "site-phrases.js");
const I18N_PATH = path.join(SITE_DIR, "assets", "site-i18n.js");
const MIRROR_PHRASES_PATHS = [
  path.join(ROOT, "site-react", "public", "assets", "site-phrases.js"),
];
const PAGES = ["poliglot-ai.html", "terms.html", "privacy.html"];
const REACT_SOURCES = [
  path.join(ROOT, "site-react", "src", "PublicSiteApp.tsx"),
  path.join(ROOT, "site-react", "src", "legacyLegalContent.ts"),
];
const GENERATED_START = "// <public-site-translations-generated>";
const GENERATED_END = "// </public-site-translations-generated>";

const LANGS = [
  ["en", "en"],
  ["es", "es"],
  ["de", "de"],
  ["fr", "fr"],
  ["it", "it"],
  ["zh", "zh-CN"],
  ["ja", "ja"],
  ["ko", "ko"],
  ["tg", "tg"],
  ["uz", "uz"],
  ["tt", "tt"],
  ["hy", "hy"],
  ["kk", "kk"],
  ["ky", "ky"],
  ["ka", "ka"],
  ["uk", "uk"],
  ["pl", "pl"],
  ["ro", "ro"],
  ["pt", "pt"],
  ["ar", "ar"],
  ["bn", "bn"],
  ["cs", "cs"],
  ["el", "el"],
  ["hi", "hi"],
  ["hu", "hu"],
  ["id", "id"],
  ["nl", "nl"],
  ["sv", "sv"],
  ["ta", "ta"],
  ["te", "te"],
  ["th", "th"],
  ["tl", "tl"],
  ["tr", "tr"],
  ["vi", "vi"],
];

const SITE_LANG_CODES = LANGS.map(([code]) => code);
const TRANSLATION_CONCURRENCY = Math.max(1, Number(process.env.TRANSLATION_CONCURRENCY || 4));
const TRANSLATION_PAUSE_MS = Math.max(0, Number(process.env.TRANSLATION_PAUSE_MS || 0));
const NATIVE_LANGUAGE_NAMES = new Set([
  "English",
  "Español",
  "Deutsch",
  "Français",
  "Italiano",
  "Тоҷикӣ",
  "O'zbekcha",
  "Татарча",
  "Հայերեն",
  "Қазақша",
  "Українська",
  "Polski",
  "Română",
  "Português",
  "Кыргызча",
  "ქართული",
  "العربية",
  "বাংলা",
  "Čeština",
  "Ελληνικά",
  "हिंदी",
  "Magyar",
  "Bahasa Indonesia",
  "Nederlands",
  "svenska",
  "தமிழ்",
  "తెలుగు",
  "ภาษาไทย",
  "Tagalog",
  "Türkçe",
  "Tiếng Việt",
]);

const PROTECTED_PATTERNS = [
  /<[^>]+>/g,
  /@[A-Za-z0-9_]+/g,
  /[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}/gi,
  /\bpoliglotai\.(?:ru|online)\/app\b/gi,
  /\b(?:Poliglot AI|Полиглот AI)\b/g,
  /\b(?:Telegram|Premium|Platinum|Free|Stars|TON|USDT|RUB|YooKassa|SBP|CEFR|OCR|AI)\b/g,
  /\b\d+(?:[.,]\d+)?\s*(?:₽|RUB|Stars|USDT|TON)\b/gi,
  /≈\s*\d[\d\s,.]*(?:₽|RUB|Stars|USDT|TON)?(?:\/год|\/year)?/gi,
];

function stripGeneratedBlock(source) {
  const start = source.indexOf(GENERATED_START);
  const end = source.indexOf(GENERATED_END);
  if (start === -1 || end === -1 || end < start) return source;
  return `${source.slice(0, start).trimEnd()}\n`;
}

function extractGeneratedTranslations(source) {
  const marker = "const publicSiteGeneratedTranslations = ";
  const start = source.indexOf(marker);
  if (start === -1) return {};
  const jsonStart = start + marker.length;
  const jsonEnd = source.indexOf(";\nObject.entries(publicSiteGeneratedTranslations)", jsonStart);
  if (jsonEnd === -1) return {};
  try {
    return JSON.parse(source.slice(jsonStart, jsonEnd));
  } catch {
    return {};
  }
}

function normalize(value) {
  return String(value || "").replace(/\s+/g, " ").trim();
}

function hasLetter(value) {
  return /\p{L}/u.test(value);
}

function hasCyrillic(value) {
  return /[А-Яа-яЁё]/.test(value);
}

function hasLatinWord(value) {
  return /[A-Za-z]{3,}/.test(value);
}

function isOnlyProtectedName(value) {
  return NATIVE_LANGUAGE_NAMES.has(value)
    || /^(?:Poliglot AI|Premium|Platinum|Free|Telegram|Stars|TON|USDT|YooKassa|SBP|AI|FAQ|CEFR|OCR|RUB)$/i.test(value)
    || /^[@/#.]/.test(value)
    || /^[\w.-]+@[\w.-]+$/.test(value);
}

function isMeaningful(value) {
  const text = normalize(value);
  if (!text || text.length < 2) return false;
  if (!hasLetter(text)) return false;
  if (isOnlyProtectedName(text)) return false;
  if (text.includes("${") || text.includes("querySelector") || text.includes("classList")) return false;
  if (/^(?:target|href|src|class|id|dataset|style|function|return|const|let|var)\b/.test(text)) return false;
  if (/^(?:https?:|mailto:|tel:|\/|#|\.)/.test(text)) return false;
  if (/^[\w-]+(?:\s+[\w-]+){2,}$/.test(text) && /\b(?:bg|text|rounded|flex|grid|shadow|border|hover|items|justify|mx|mb|mt|px|py)\b/.test(text)) return false;
  return true;
}

function addCandidate(set, value) {
  const text = normalize(value);
  if (isMeaningful(text)) set.add(text);
}

function decodeHtmlEntities(value) {
  return String(value)
    .replace(/&copy;/g, "©")
    .replace(/&nbsp;/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, "\"");
}

function collectCandidates() {
  const candidates = new Set();

  const collectHtml = (html) => {
    const withoutScripts = html
      .replace(/<script[\s\S]*?<\/script>/gi, " ")
      .replace(/<style[\s\S]*?<\/style>/gi, " ")
      .replace(/<!--[\s\S]*?-->/g, " ");

    for (const match of withoutScripts.matchAll(/>([^<>]+)</g)) {
      addCandidate(candidates, decodeHtmlEntities(match[1]));
    }

    for (const match of html.matchAll(/(?:alt|aria-label|title|placeholder|content)="([^"]+)"/g)) {
      addCandidate(candidates, decodeHtmlEntities(match[1]));
    }
  };

  for (const page of PAGES) {
    const html = fs.readFileSync(path.join(SITE_DIR, page), "utf8");
    collectHtml(html);

    for (const block of html.matchAll(/<script[\s\S]*?>([\s\S]*?)<\/script>/gi)) {
      for (const match of block[1].matchAll(/"((?:\\.|[^"\\])*)"/g)) {
        let value = match[1];
        try {
          value = JSON.parse(`"${value}"`);
        } catch {
          // Keep raw string literal content.
        }
        const text = normalize(decodeHtmlEntities(value));
        if (hasCyrillic(text)) addCandidate(candidates, text);
      }
    }
  }

  for (const file of REACT_SOURCES) {
    if (!fs.existsSync(file)) continue;
    const source = fs.readFileSync(file, "utf8");
    if (file.endsWith("legacyLegalContent.ts")) {
      for (const match of source.matchAll(/=\s*"((?:\\.|[^"\\])*)"/g)) {
        try {
          collectHtml(JSON.parse(`"${match[1]}"`));
        } catch {
          // Keep going; malformed strings are not useful translation candidates.
        }
      }
      continue;
    }

    const withoutComments = source
      .replace(/\/\*[\s\S]*?\*\//g, " ")
      .replace(/\/\/.*$/gm, " ");

    for (const match of withoutComments.matchAll(/>([^<>]+)</g)) {
      addCandidate(candidates, decodeHtmlEntities(match[1]));
    }

    for (const match of withoutComments.matchAll(/["'`]((?:\\.|[^"'`\\]){2,})["'`]/g)) {
      let value = match[1];
      try {
        value = JSON.parse(`"${value.replace(/"/g, '\\"')}"`);
      } catch {
        // Keep raw string literal content.
      }
      addCandidate(candidates, decodeHtmlEntities(value));
    }
  }

  const i18nSource = fs.readFileSync(I18N_PATH, "utf8");
  const landingBlock = i18nSource.slice(
    i18nSource.indexOf("const landingEnglishPhrases = ["),
    i18nSource.indexOf("landingEnglishPhrases.forEach")
  );
  for (const match of landingBlock.matchAll(/\["((?:\\.|[^"\\])*)",\s*"((?:\\.|[^"\\])*)"/g)) {
    try {
      addCandidate(candidates, JSON.parse(`"${match[1]}"`));
    } catch {
      addCandidate(candidates, match[1]);
    }
  }

  return [...candidates].sort((a, b) => a.localeCompare(b, "ru"));
}

function createContext(lang, phrasesSource, i18nSource) {
  const ctx = {
    console,
    TextEncoder,
    TextDecoder,
    URL,
    URLSearchParams,
    CustomEvent: class CustomEvent {
      constructor(type, init) {
        this.type = type;
        this.detail = init?.detail;
      }
    },
    NodeFilter: { SHOW_TEXT: 4 },
    location: {
      search: `?lang=${lang}`,
      href: `https://example.com/poliglot-ai.html?lang=${lang}`,
      host: "example.com",
    },
    localStorage: {
      getItem: () => lang,
      setItem: () => {},
    },
    document: {
      documentElement: { lang },
      body: { dataset: { page: "landing" } },
      addEventListener: () => {},
      querySelectorAll: () => [],
      createTreeWalker: () => ({ nextNode: () => false }),
    },
    window: {
      addEventListener: () => {},
      dispatchEvent: () => {},
    },
  };
  ctx.window.window = ctx.window;
  ctx.window.document = ctx.document;
  ctx.window.location = ctx.location;
  ctx.window.localStorage = ctx.localStorage;
  vm.runInNewContext(phrasesSource, ctx, { filename: PHRASES_PATH });
  vm.runInNewContext(i18nSource, ctx, { filename: I18N_PATH });
  return ctx.window.poliglotSiteI18n;
}

function protectText(source) {
  const placeholders = [];
  let text = source.replace(/\bПолиглот AI\b/g, "Poliglot AI");

  for (const pattern of PROTECTED_PATTERNS) {
    text = text.replace(pattern, (match) => {
      const token = `__P${placeholders.length}__`;
      placeholders.push(match === "Полиглот AI" ? "Poliglot AI" : match);
      return token;
    });
  }

  return {
    text,
    restore(value) {
      let restored = value;
      placeholders.forEach((original, index) => {
        const tokenPattern = new RegExp(`__\\s*P\\s*${index}\\s*__`, "gi");
        restored = restored.replace(tokenPattern, original);
      });
      return restored
        .replace(/(\p{L})(Telegram|Premium|Platinum|Free|Stars|TON|USDT|RUB|YooKassa|SBP|CEFR|OCR|AI)\b/gu, "$1 $2")
        .replace(/\b(Telegram|Premium|Platinum|Free|Stars|TON|USDT|RUB|YooKassa|SBP|CEFR|OCR|AI)(\p{L})/gu, "$1 $2")
        .replace(/\s+([,.!?;:%])/g, "$1")
        .replace(/([([{])\s+/g, "$1")
        .replace(/\s+([)\]}])/g, "$1")
        .replace(/\s+·\s+/g, " · ")
        .trim();
    },
  };
}

function googleSourceCode(source) {
  return hasCyrillic(source) ? "ru" : "en";
}

function needsTranslation(source, lang, current, english) {
  if (lang === "ru") return false;
  if (!isMeaningful(source)) return false;
  if (isOnlyProtectedName(source)) return false;
  if (lang === "en") return hasCyrillic(source) && current === source;
  if (current === source) return true;
  if (english && current === english && hasLatinWord(english) && !isOnlyProtectedName(english)) return true;
  return false;
}

async function translateBatch(items, targetGoogleCode) {
  const groups = new Map();
  for (const item of items) {
    const sl = googleSourceCode(item.source);
    if (sl === targetGoogleCode) {
      item.result = item.source.replace(/\bПолиглот AI\b/g, "Poliglot AI");
      continue;
    }
    const key = sl;
    if (!groups.has(key)) groups.set(key, []);
    groups.get(key).push(item);
  }

  for (const [sourceGoogleCode, group] of groups.entries()) {
    let offset = 0;
    while (offset < group.length) {
      const batch = [];
      let size = 0;
      while (offset < group.length && batch.length < 10 && size < 2800) {
        const item = group[offset++];
        const protectedItem = protectText(item.source);
        item.protectedItem = protectedItem;
        size += protectedItem.text.length;
        batch.push(item);
      }

      const parts = await translateMarkedBatch(batch, sourceGoogleCode, targetGoogleCode);

      for (let index = 0; index < batch.length; index += 1) {
        const raw = parts.get(index) || await fetchGoogle(batch[index].protectedItem.text, sourceGoogleCode, targetGoogleCode);
        batch[index].result = batch[index].protectedItem.restore(raw.replace(/^\[\[\d+\]\]\s*/, ""));
      }

      if (TRANSLATION_PAUSE_MS > 0) {
        await new Promise((resolve) => setTimeout(resolve, TRANSLATION_PAUSE_MS));
      }
    }
  }
}

async function translateMarkedBatch(batch, sourceGoogleCode, targetGoogleCode) {
  const marked = batch.map((item, index) => `[[${index}]] ${item.protectedItem.text}`).join("\n");
  try {
    const translated = await fetchGoogle(marked, sourceGoogleCode, targetGoogleCode);
    return splitMarkedTranslation(translated, batch.length);
  } catch (error) {
    if (!String(error?.message || error).includes("HTTP 413") || batch.length <= 1) throw error;
    const midpoint = Math.ceil(batch.length / 2);
    const left = await translateMarkedBatch(batch.slice(0, midpoint), sourceGoogleCode, targetGoogleCode);
    const right = await translateMarkedBatch(batch.slice(midpoint), sourceGoogleCode, targetGoogleCode);
    const merged = new Map(left);
    for (const [relativeIndex, value] of right.entries()) merged.set(relativeIndex + midpoint, value);
    return merged;
  }
}

async function fetchGoogle(text, sl, tl) {
  const url = new URL("https://translate.googleapis.com/translate_a/single");
  url.searchParams.set("client", "gtx");
  url.searchParams.set("sl", sl);
  url.searchParams.set("tl", tl);
  url.searchParams.set("dt", "t");
  url.searchParams.set("q", text);

  for (let attempt = 1; attempt <= 4; attempt += 1) {
    try {
      const response = await fetch(url);
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      const data = await response.json();
      return (data?.[0] || []).map((part) => part?.[0] || "").join("");
    } catch (error) {
      if (attempt === 4) throw error;
      await new Promise((resolve) => setTimeout(resolve, attempt * 500));
    }
  }
  return text;
}

function splitMarkedTranslation(text, expectedCount) {
  const result = new Map();
  const marker = /\[\[(\d+)\]\]/g;
  const matches = [...text.matchAll(marker)];
  for (let i = 0; i < matches.length; i += 1) {
    const index = Number(matches[i][1]);
    const start = matches[i].index + matches[i][0].length;
    const end = i + 1 < matches.length ? matches[i + 1].index : text.length;
    if (index >= 0 && index < expectedCount) {
      result.set(index, text.slice(start, end).trim());
    }
  }
  return result;
}

function buildCurrentSites(phrasesSource, i18nSource) {
  const sites = {};
  for (const [code] of LANGS) sites[code] = createContext(code, phrasesSource, i18nSource);
  return sites;
}

function shouldKeepEnglishTerm(source, value) {
  if (!value) return false;
  if (isOnlyProtectedName(value)) return true;
  return false;
}

async function main() {
  const phraseFile = fs.readFileSync(PHRASES_PATH, "utf8");
  const basePhraseFile = stripGeneratedBlock(phraseFile);
  const existingGenerated = extractGeneratedTranslations(phraseFile);
  const i18nSource = fs.readFileSync(I18N_PATH, "utf8");
  const candidates = collectCandidates();
  const currentSites = buildCurrentSites(phraseFile, i18nSource);
  const englishSite = currentSites.en;
  const generated = JSON.parse(JSON.stringify(existingGenerated));
  const workByLang = new Map();

  for (const source of candidates) {
    const english = englishSite.translatePhrase(source);
    for (const [lang, googleCode] of LANGS) {
      if (lang === "ru") continue;
      const current = currentSites[lang].translatePhrase(source);
      if (!needsTranslation(source, lang, current, english)) continue;
      if (!workByLang.has(lang)) workByLang.set(lang, []);
      workByLang.get(lang).push({ source, lang, googleCode });
    }
  }

  let total = 0;
  const languageJobs = [...workByLang.entries()];
  let jobIndex = 0;
  const workers = Array.from({ length: Math.min(TRANSLATION_CONCURRENCY, languageJobs.length) }, async () => {
    while (jobIndex < languageJobs.length) {
      const currentIndex = jobIndex++;
      const [lang, items] = languageJobs[currentIndex];
      total += items.length;
      console.log(`Translating ${items.length} strings for ${lang}...`);
      const googleCode = LANGS.find(([code]) => code === lang)?.[1] || lang;
      await translateBatch(items, googleCode);
      for (const item of items) {
        if (shouldKeepEnglishTerm(item.source, item.result)) continue;
        generated[item.source] ||= {};
        generated[item.source][lang] = item.result;
      }
    }
  });
  await Promise.all(workers);

  const ordered = {};
  for (const source of Object.keys(generated).sort((a, b) => a.localeCompare(b, "ru"))) {
    ordered[source] = {};
    for (const [code] of LANGS) {
      if (code !== "ru" && generated[source][code]) ordered[source][code] = generated[source][code];
    }
  }

  const block = [
    GENERATED_START,
    "window.poliglotPhraseTranslations = window.poliglotPhraseTranslations || {};",
    `const publicSiteGeneratedTranslations = ${JSON.stringify(ordered, null, 2)};`,
    "Object.entries(publicSiteGeneratedTranslations).forEach(([source, values]) => {",
    "  window.poliglotPhraseTranslations[source] = { ...(window.poliglotPhraseTranslations[source] || {}), ...values };",
    "});",
    GENERATED_END,
    "",
  ].join("\n");

  const nextPhrases = `${basePhraseFile.trimEnd()}\n\n${block}`;
  fs.writeFileSync(PHRASES_PATH, nextPhrases, "utf8");
  for (const mirrorPath of MIRROR_PHRASES_PATHS) {
    fs.mkdirSync(path.dirname(mirrorPath), { recursive: true });
    fs.writeFileSync(mirrorPath, nextPhrases, "utf8");
  }
  console.log(`Generated ${Object.keys(ordered).length} phrase entries, ${total} language values considered.`);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
