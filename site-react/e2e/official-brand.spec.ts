import { expect, test } from "@playwright/test";

test('legal pages preserve full documents and offer Russian and English', async ({ page }) => {
  for (const path of ['privacy', 'terms', 'agreement', 'consent']) {
    await page.goto(`/${path}.html?lang=en`);
    await expect(page.locator('html')).toHaveAttribute('lang', 'en');
    const picker = page.locator('[data-site-language-select]').first();
    await expect(picker.locator('option')).toHaveCount(2);
    await expect(picker).toHaveValue('en');
    const documentBody = page.locator('.legal-document-shell');
    await expect(documentBody).toContainText(/personal data|agreement|service|privacy/i);
    await picker.selectOption('ru');
    await expect(page.locator('html')).toHaveAttribute('lang', 'ru');
    await expect(documentBody).toContainText(/[А-Яа-яЁё]{5}/);
  }
});

test('knowledge navigator opens the existing English guides and returns to Russian articles', async ({ page }) => {
  await page.goto('/knowledge/');
  await page.locator('.knowledge-language').selectOption('en');
  await expect(page.locator('html')).toHaveAttribute('lang', 'en');
  await expect(page.locator('main a[href^="/en/"]')).not.toHaveCount(0);
  await page.locator('.knowledge-language').selectOption('ru');
  await expect(page).toHaveURL(/\/knowledge\/$/);
});

test('visitors see a lesson example and can jump to pricing', async ({ page }) => {
  await page.emulateMedia({ reducedMotion: 'reduce' });
  await page.goto('/neriva.html?lang=ru');
  await expect(page.locator('.hero-actions a[data-entry="web-app"]')).toHaveText('Начать урок');
  await expect(page.locator('.hero-actions a[data-entry="web-app"]')).toBeInViewport();
  await expect(page.locator('main > section').nth(1)).toHaveClass('lesson-scenario');
  await expect(page.locator('.lesson-example')).toContainText('I have a reservation.');
  await expect(page.locator('main > section').nth(3)).toHaveAttribute('id', 'pricing');
  await page.locator('.nav-pricing-link').click();
  await expect(page).toHaveURL(/#pricing$/);
  await expect(page.locator('#pricing h2')).toBeInViewport();
});

test("official brand and current application screenshots are visible", async ({ page }) => {
  await page.addInitScript(() => localStorage.setItem("poliglot-cookie-consent", JSON.stringify({ version: 1, necessary: true, analyticsMarketing: false })));
  await page.goto("/neriva.html?lang=ru");
  const icon = await page.locator('head link[rel="icon"]').getAttribute('href');
  expect(icon).toContain('/assets/brand/neriva-icon-192.png');
  expect((await page.request.get(icon!)).headers()['content-type']).toContain('image/png');
  const logo = page.locator('.public-brand img').first();
  await expect(logo).toHaveAttribute("src", /\/assets\/brand\/neriva-/);
  await expect(page.locator('.public-brand .official-brand').first()).toHaveCSS('background-color', 'rgba(0, 0, 0, 0)');
  await expect(page.locator('.hero-brand-emblem')).toHaveCSS('background-color', 'rgba(0, 0, 0, 0)');
  await expect(page.locator('.hero-product-frame img').first()).toHaveAttribute("src", /\/assets\/product\/2026-11-15\//);
  await expect.poll(() => logo.evaluate((img: HTMLImageElement) => img.complete && img.naturalWidth > 0)).toBe(true);
  const aspect = await logo.evaluate((img: HTMLImageElement) => ({ actual: img.getBoundingClientRect().width / img.getBoundingClientRect().height, natural: img.naturalWidth / img.naturalHeight }));
  expect(Math.abs(aspect.actual / aspect.natural - 1)).toBeLessThan(.05);
  for (const width of [390, 820, 1440]) {
    await page.setViewportSize({ width, height: 900 });
    await expect(logo).toBeVisible();
    const size = await page.evaluate(() => ({ content: document.documentElement.scrollWidth, width: innerWidth }));
    expect(size.content).toBeLessThanOrEqual(size.width + 1);
  }
});

test('official brand export matches the recorded source manifest in both products', async () => {
  const fs = await import('node:fs/promises');
  const path = await import('node:path');
  const crypto = await import('node:crypto');
  const root = path.resolve(process.cwd(), '..');
  const manifest = JSON.parse(await fs.readFile(path.join(root, 'docs/specs/official-brand/brand-manifest.json'), 'utf8'));
  expect(manifest.sources).toHaveLength(3);
  for (const asset of manifest.outputs) {
    for (const dir of ['web/assets', 'site-react/public/assets/brand']) {
      const bytes = await fs.readFile(path.join(root, dir, asset.filename));
      expect(crypto.createHash('sha256').update(bytes).digest('hex')).toBe(asset.sha256);
    }
  }
});
