import { expect, test } from "@playwright/test";

test("install link carries through sign-in and opens the browser prompt on click", async ({ page }) => {
  await page.route("**/api/session", (route) => route.fulfill({ json: { authenticated: false } }));
  await page.goto("/app/?install=1&lang=en");
  await expect(page).toHaveURL(/\/app\/login\?install=1&lang=en/);
  const panel = page.getByRole("region", { name: "Install NERIVA" });
  await expect(panel).toContainText("browser menu");
  await page.evaluate(() => {
    const event = Object.assign(new Event("beforeinstallprompt", { cancelable: true }), {
      prompt: async () => window.dispatchEvent(new Event("appinstalled")),
    });
    window.dispatchEvent(event);
  });
  await panel.getByRole("button", { name: "Install NERIVA" }).click();
  await expect(panel).toHaveCount(0);
});

test("installation shows the iPhone share-menu instructions", async ({ page }) => {
  await page.addInitScript(() => Object.defineProperty(navigator, "platform", { value: "iPhone" }));
  await page.route("**/api/session", (route) => route.fulfill({ json: { authenticated: false } }));
  await page.goto("/app/login?install=1&lang=ru");
  const panel = page.getByRole("region", { name: "Установить NERIVA" });
  await expect(panel).toContainText("Поделиться");
  await expect(panel).toContainText("На экран «Домой»");
  await expect(panel.getByRole("button", { name: "Установить NERIVA" })).toHaveCount(0);
  const close = panel.getByRole("button", { name: "Закрыть" });
  expect((await close.boundingBox())!.width).toBeGreaterThanOrEqual(44);
  await close.focus();
  await page.keyboard.press("Escape");
  await expect(panel).toHaveCount(0);
  await expect(page).not.toHaveURL(/install=1/);
});

test("the service worker reuses entry assets and stores a single public app shell", async ({ page, context }) => {
  const shell = await page.request.get("/app/");
  const asset = (await shell.text()).match(/src="(\/app\/assets\/index-[^"]+\.js)"/)?.[1];
  expect(asset).toBeTruthy();
  await page.route("**/api/session", (route) => route.fulfill({ json: { authenticated: false } }));
  await page.goto("/app/login");
  await page.evaluate(async () => {
    await navigator.serviceWorker.register("/app/offline-deck-sw.js");
    await navigator.serviceWorker.ready;
  });
  await expect.poll(() => page.evaluate(() => Boolean(navigator.serviceWorker.controller))).toBe(true);
  const cacheKeys = () => page.evaluate(async () => {
    const cache = await caches.open("neriva-app-static-v3");
    return (await cache.keys()).map((entry) => new URL(entry.url).pathname + new URL(entry.url).search);
  });
  await expect.poll(cacheKeys).toContain(asset);
  const served = page.waitForResponse((response) => response.url().endsWith(asset!) && response.fromServiceWorker());
  await page.evaluate(async (url) => { await fetch(url); }, asset!);
  expect((await served).ok()).toBe(true);
  await page.unroute("**/api/session");
  await context.setOffline(true);
  const response = await page.goto("/app/?install=1&lang=en");
  expect(response?.fromServiceWorker()).toBe(true);
  await expect(page.getByText("You’re offline", { exact: true })).toBeVisible();
  const keys = await cacheKeys();
  expect(keys.filter((key) => key === "/app/")).toHaveLength(1);
  expect(keys.every((key) => !key.includes("?") && !key.startsWith("/api/"))).toBe(true);
});
