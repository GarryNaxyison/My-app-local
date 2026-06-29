#!/usr/bin/env node

import fs from "node:fs";
import path from "node:path";

const root = process.cwd();
const appLocaleCodes = "ru,en,es,de,fr,it,zh,ja,ko,tg,uz,tt,hy,kk,ky,ka,uk,pl,ro,pt,ar,bn,cs,el,hi,hu,id,nl,sv,ta,te,th,tl,tr,vi".split(",");
const appPath = path.join(root, "web-react", "src", "App.tsx");
const copyPath = path.join(root, "web-react", "src", "lib", "fullLocaleCopy.ts");
const allowedPath = path.join(root, "tools", "i18n", "allowed_invariant_terms.json");

function readRequiredFile(filePath, label) {
  try {
    return fs.readFileSync(filePath, "utf8");
  } catch (error) {
    console.error(`Missing required ${label}: ${path.relative(root, filePath)}`);
    console.error(String(error?.message || error));
    process.exit(1);
  }
}

const appSource = readRequiredFile(appPath, "WebApp source");
const copySource = readRequiredFile(copyPath, "fullLocaleCopy source");
const allowedTerms = JSON.parse(readRequiredFile(allowedPath, "allowed invariant terms"));

// Coverage source: copy("key", "fallback") call sites in App.tsx.
const copyCallPattern = /\bcopy\(\s*["'`]([^"'`]+)["'`]\s*,/g;
const copiedKeys = [...appSource.matchAll(copyCallPattern)]
  .map((match) => match[1])
  .filter((key) => !key.includes("${"));
const uniqueKeys = [...new Set(copiedKeys)].sort();

function objectBlockAfter(source, marker) {
  const start = source.indexOf(marker);
  if (start < 0) return "";
  const open = source.indexOf("{", start + marker.length);
  if (open < 0) return "";
  let depth = 0;
  let quote = "";
  let escaped = false;
  for (let index = open; index < source.length; index += 1) {
    const char = source[index];
    if (quote) {
      if (escaped) {
        escaped = false;
      } else if (char === "\\") {
        escaped = true;
      } else if (char === quote) {
        quote = "";
      }
      continue;
    }
    if (char === "\"" || char === "'" || char === "`") {
      quote = char;
      continue;
    }
    if (char === "{") depth += 1;
    if (char === "}") {
      depth -= 1;
      if (depth === 0) return source.slice(open, index + 1);
    }
  }
  return "";
}

const missing = [];
for (const code of appLocaleCodes) {
  const localeBlock = objectBlockAfter(copySource, `"${code}":`) || objectBlockAfter(copySource, `  ${code}:`);
  if (!localeBlock) {
    missing.push(`locale:${code}`);
    continue;
  }
  for (const key of uniqueKeys) {
    const keyMarkers = [`${key}:`, `"${key}":`, `'${key}':`];
    if (!keyMarkers.some((marker) => localeBlock.includes(marker))) {
      missing.push(`${code}:${key}`);
    }
  }
}

const generatedMarkerLeaks = [
  "ZZZTERM",
  "ЗЗЗТЕРМ",
  "ZZZ术语",
  "?T0?",
  "?T1?",
  "?Т0?",
  "?Т1?",
  "app_guide_title: \"app_guide_title\"",
  "\"app_guide_title\": \"app_guide_title\"",
].filter((marker) => copySource.includes(marker));

for (const marker of generatedMarkerLeaks) {
  missing.push(`generated-marker:${marker}`);
}

for (const code of appLocaleCodes) {
  const localeBlock = objectBlockAfter(copySource, `"${code}":`) || objectBlockAfter(copySource, `  ${code}:`);
  const valueFor = (key) => {
    const match = localeBlock.match(new RegExp(`["']?${key}["']?\\s*:\\s*["']([^"']*)["']`));
    return match?.[1] || "";
  };
  for (const match of localeBlock.matchAll(/["']([^"']+)["']\s*:\s*["']([^"']*)["']/g)) {
    if (/^\?{2,}$/.test(match[2])) missing.push(`${code}:${match[1]} question-mark placeholder`);
  }
  if (valueFor("pay_stars") !== "Telegram Stars") missing.push(`${code}:pay_stars must stay Telegram Stars`);
  if (valueFor("ai_router") !== "AI Router") missing.push(`${code}:ai_router must stay AI Router`);
  if (code !== "en" && valueFor("dashboard") === "Dashboard") missing.push(`${code}:dashboard exact English fallback`);
}

for (const term of allowedTerms) {
  if (typeof term !== "string" || !term.trim()) {
    console.error("allowed invariant terms must be non-empty strings");
    process.exit(1);
  }
}

if (missing.length) {
  console.error(`Missing explicit WebApp locale entries (${missing.length}):`);
  for (const item of missing.slice(0, 160)) console.error(`- ${item}`);
  if (missing.length > 160) console.error(`... ${missing.length - 160} more`);
  process.exit(1);
}

console.log(`WebApp i18n coverage OK: ${appLocaleCodes.length} locales, ${uniqueKeys.length} copy keys.`);
