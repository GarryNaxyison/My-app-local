import { expect, test } from "@playwright/test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { canonicalPhraseKey } from "../src/lib/phrasebook";

const testDirectory = fileURLToPath(new URL(".", import.meta.url));

test("desktop regression: phrase suggestions share one canonical key", () => {
  expect(canonicalPhraseKey(" Hello! ")).toBe(canonicalPhraseKey("hello"));
  expect(canonicalPhraseKey("Hello.")).toBe(canonicalPhraseKey("HELLO?"));
  expect(canonicalPhraseKey("Hello there!")).not.toBe(canonicalPhraseKey("Hello!"));
});

test("desktop regression: containment rules are separate from mobile rules", () => {
  const css = readFileSync(resolve(testDirectory, "../src/styles/app.css"), "utf8");
  const logout = readFileSync(resolve(testDirectory, "../src/components/ui/logout-button.tsx"), "utf8");

  expect(css).toContain("Desktop stability: one scroll owner");
  expect(css).toContain(".auth-turnstile-slot-v2 iframe");
  expect(css).toContain(".context-display--roleplay");
  expect(logout).toContain("minWidth");
});
