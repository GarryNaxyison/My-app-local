import { expect, test } from "@playwright/test";

test("landing opens one example at a time and returns focus", async ({ page }) => {
  await page.emulateMedia({ reducedMotion: "reduce" });
  await page.addInitScript(() => localStorage.setItem("poliglot-cookie-consent", JSON.stringify({ version: 1, necessary: true, analyticsMarketing: false })));
  await page.goto("/neriva.html?lang=ru");

  for (const selector of [".use-case-section", ".ecosystem-showcase", ".mobile-section"]) {
    const section = page.locator(selector);
    const trigger = section.locator(".showcase-launchers button").first();
    await expect(section.locator("iframe")).toHaveCount(0);
    await trigger.click();
    const dialog = section.getByRole("dialog");
    await expect(dialog).toBeVisible();
    const tabs = dialog.getByRole("tab");
    await tabs.first().focus();
    await page.keyboard.press("ArrowRight");
    await expect(tabs.nth(1)).toBeFocused();
    await expect(tabs.nth(1)).toHaveAttribute("aria-selected", "true");
    await expect(dialog.locator("iframe")).toHaveCount(1);
    await page.keyboard.press("Escape");
    await expect(dialog).not.toBeVisible();
    await expect(trigger).toBeFocused();
  }

  const proof = page.locator(".proof-section");
  await proof.getByRole("tab", { name: "Карточки слов" }).click();
  await expect(proof.getByRole("tabpanel")).toContainText("Печать / PDF");
  await expect(proof.locator("iframe")).toHaveCount(0);
  expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBeLessThanOrEqual(await page.evaluate(() => innerWidth));
});
