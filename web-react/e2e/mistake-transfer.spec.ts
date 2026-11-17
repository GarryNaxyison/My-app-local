import { expect, test } from "@playwright/test";
import { buildMistakeTransfer, checkMistakeTransfer } from "../src/lib/mistakeTransfer";

test("transfer tests the corrected construction with a different verb and situation", () => {
  const task = buildMistakeTransfer("I go yesterday", "I went yesterday", "en");
  expect(task?.prompt.en).toContain("station");
  expect(checkMistakeTransfer(task!, "I walked to the station yesterday.")).toBe(true);
  expect(checkMistakeTransfer(task!, "Yesterday, I walked to the station.")).toBe(true);
  expect(checkMistakeTransfer(task!, "I walk to the station yesterday.")).toBe(false);
  expect(checkMistakeTransfer(task!, "I went yesterday.")).toBe(false);
  expect(buildMistakeTransfer("I go", "I went", "fr")).toBeNull();
  expect(buildMistakeTransfer("Hello", "Good morning", "en")).toBeNull();
});
