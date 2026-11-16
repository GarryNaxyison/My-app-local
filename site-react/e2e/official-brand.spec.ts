import { expect, test } from "@playwright/test";

test("official brand and current application screenshots are visible", async ({ page }) => {
  await page.addInitScript(() => localStorage.setItem("poliglot-cookie-consent", JSON.stringify({ version: 1, necessary: true, analyticsMarketing: false })));
  await page.goto("/neriva.html?lang=ru");
  const logo = page.locator('.public-brand img').first();
  await expect(logo).toHaveAttribute("src", /\/assets\/brand\/neriva-/);
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
