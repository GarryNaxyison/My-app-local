import { expect, test, type Page } from "@playwright/test";

async function prepareLanding(page: Page) {
  await page.route("https://mc.yandex.ru/**", (route) => route.abort());
  await page.addInitScript(() => {
    localStorage.setItem(
      "poliglot-cookie-consent",
      JSON.stringify({ version: 1, necessary: true, analyticsMarketing: false }),
    );
    localStorage.setItem("poliglot-site-theme", "dark");
  });
}

test.beforeEach(async ({ page }) => {
  await prepareLanding(page);
});

test("Russian landing adds human pain, use-case and premium value blocks", async ({ page }) => {
  await page.goto("/neriva.html?lang=ru");

  await expect(page.locator(".outcome-section")).toContainText("Если вы учите язык, но всё равно молчите");
  await expect(page.locator(".outcome-section")).toContainText("не даёт ошибкам пропасть");
  await expect(page.locator(".outcome-card")).toHaveCount(3);

  await expect(page.locator(".use-case-section")).toContainText("Для чего открывают NERIVA");
  await expect(page.locator(".use-case-card")).toHaveCount(4);
  await expect(page.locator(".use-case-section")).toContainText("Поездка");
  await expect(page.locator(".use-case-section")).toContainText("Работа");
  await expect(page.locator(".use-case-section")).toContainText("Экзамен");
  await expect(page.locator(".use-case-section")).toContainText("Произношение");

  await expect(page.locator(".proof-section")).toContainText("Как это выглядит в продукте");
  await expect(page.locator(".proof-step")).toHaveCount(4);
  await expect(page.locator(".proof-section .product-shot")).toHaveCount(2);

  await expect(page.locator(".pricing-section")).toContainText("Premium нужен, когда вы занимаетесь почти каждый день");
  await expect(page.locator(".plan-card.is-featured")).toContainText("голос, фото и разбор ошибок");
});
