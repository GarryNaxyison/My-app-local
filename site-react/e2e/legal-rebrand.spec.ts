import { expect, test } from "@playwright/test";

const legalPages = [
  { path: "/privacy.html", title: "Политика обработки персональных данных" },
  { path: "/terms.html", title: "Условия использования" },
  { path: "/agreement.html", title: "Пользовательское соглашение" },
  { path: "/consent.html", title: "Согласие на обработку персональных данных" },
] as const;

test.describe("NERIVA legal pages", () => {
  for (const legalPage of legalPages) {
    test(`${legalPage.path} uses current brand, bot, and support contacts`, async ({ page }) => {
      await page.goto(legalPage.path, { waitUntil: "domcontentloaded" });
      await expect(page.getByRole("heading", { name: legalPage.title, level: 1 })).toBeVisible();

      const bodyText = await page.locator("body").innerText();
      expect(bodyText).toContain("NERIVA");
      expect(bodyText).toContain("support@neriva.ru");
      expect(bodyText).toContain("@NERIVAapp_bot");
      expect(bodyText).toContain("neriva.ru/app");
      expect(bodyText).toContain("35 языков интерфейса");

      expect(bodyText).not.toMatch(/POLIGLOT|Poliglot AI|Poliglot_AI_bot|@Poliglot/i);
      expect(bodyText).not.toMatch(/support@poliglot|support@poliglotai|poliglotai@gmail|poliglot@gmail/i);
      expect(bodyText).not.toMatch(/poliglotai\.ru|poliglotai\.online/i);
      expect(bodyText).not.toContain("neriva.ru и neriva.ru");
      expect(bodyText).not.toContain("20 языков");
    });
  }
});
