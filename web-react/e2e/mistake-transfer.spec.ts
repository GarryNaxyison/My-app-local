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

for (const [original, correction, valid, invalid] of [
  ["I have reservation.", "I have a reservation.", "I need an umbrella.", "I need a umbrella."],
  ["They is outside.", "They are outside.", "We are ready.", "We is ready."],
  ["I didn't went home.", "I didn't go home.", "I did not call Mia yesterday.", "I didn’t called Mia yesterday."],
] as const) {
  test(`transfer recognizes the correction in "${original}"`, () => {
    const task = buildMistakeTransfer(original, correction, "en");
    expect(task).not.toBeNull();
    expect(checkMistakeTransfer(task!, valid)).toBe(true);
    expect(checkMistakeTransfer(task!, invalid)).toBe(false);
  });
}

test("transfer skips a correction with unrelated changes", () => {
  expect(buildMistakeTransfer("They is outside.", "They are waiting inside.", "en")).toBeNull();
  expect(buildMistakeTransfer("I didn't went home.", "I didn't eat dinner.", "en")).toBeNull();
});
