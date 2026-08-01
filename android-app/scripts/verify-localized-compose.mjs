import fs from "node:fs";
import path from "node:path";

const uiRoot = path.resolve("app/src/main/java/ru/neriva/app/ui");
const violations = [];
const allowedLiterals = new Set([
  "AI",
  "XP",
  "NERIVA hero",
  "NERIVAapp_bot",
  "Telegram",
  "YouTube",
  "Instagram",
  "TikTok",
  "Demo Mode (test/test)",
  "Telegram login: @NERIVAapp_bot",
]);

function requiresLocalization(value) {
  if (!value || allowedLiterals.has(value) || /^X+(?:-X+)+$/.test(value)) return false;
  if (/^(?:https?:\/\/|[a-z][a-z0-9-]*$|[a-z]{2}(?:_[A-Z]{2})?|[A-C][1-2]|image\/|text\/|\.)/.test(value)) return false;
  return /[A-Za-z]/.test(value) && (/[A-Z]/.test(value) || /\s/.test(value));
}

function visit(directory) {
  for (const entry of fs.readdirSync(directory, { withFileTypes: true })) {
    const file = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      visit(file);
      continue;
    }
    if (!entry.name.endsWith(".kt")) continue;
    const source = fs.readFileSync(file, "utf8");
    const matcher = /"((?:\\.|[^"\\])*)"/g;
    for (const match of source.matchAll(matcher)) {
      const value = match[1];
      if (!requiresLocalization(value)) continue;
      if (value === "" || value === "•" || value.includes("$") || /^X+(?:-X+)+$/.test(value) || value.startsWith("Demo Mode") || value.startsWith("Telegram login:")) continue;
      const line = source.slice(0, match.index).split("\n").length;
      violations.push(`${path.relative(process.cwd(), file)}:${line}: ${value}`);
    }
  }
}

visit(uiRoot);
if (violations.length) {
  console.error("Hardcoded Compose UI strings found:\n" + violations.join("\n"));
  process.exit(1);
}
