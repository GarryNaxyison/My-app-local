import fs from "node:fs";
import path from "node:path";

const uiRoot = path.resolve("app/src/main/java/ru/neriva/app/ui");
const violations = [];

function visit(directory) {
  for (const entry of fs.readdirSync(directory, { withFileTypes: true })) {
    const file = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      visit(file);
      continue;
    }
    if (!entry.name.endsWith(".kt")) continue;
    const source = fs.readFileSync(file, "utf8");
    const matcher = /\b(?:Text|ScreenScaffold)\(\s*"((?:\\.|[^"\\])*)"/g;
    for (const match of source.matchAll(matcher)) {
      const value = match[1];
      if (value === "•" || value.includes("$") || /^X+(?:-X+)+$/.test(value) || value.startsWith("Demo Mode") || value.startsWith("Telegram login:")) continue;
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
