import { expect, test } from "@playwright/test";

const requiredCloudflareFiles = [
  "poliglot-ai.html",
  "privacy.html",
  "terms.html",
  "agreement.html",
  "consent.html",
  "404.html",
  "maintenance.html",
  "robots.txt",
  "sitemap.xml",
  "_headers",
  "_redirects",
  "assets/site-i18n.js",
  "assets/site-phrases.js",
  "assets/seo/poliglot-ai-og-ru.jpg",
  "assets/seo/poliglot-ai-og-en.jpg",
] as const;

test("Cloudflare Pages output contains the static landing, SEO files, and fallback pages", async () => {
  const fs = await import("node:fs/promises");
  const path = await import("node:path");
  const outputRoot = path.resolve(process.cwd(), "dist");

  for (const filePath of requiredCloudflareFiles) {
    const stat = await fs.stat(path.join(outputRoot, filePath));
    expect(stat.size, `${filePath} should be present and non-empty in Cloudflare Pages output`).toBeGreaterThan(20);
  }

  const bundleDir = path.join(outputRoot, "assets", "site-react");
  const bundleFiles = await fs.readdir(bundleDir);
  expect(bundleFiles.filter((file) => /^main-[\w-]+\.js$/.test(file))).toHaveLength(1);
  expect(bundleFiles.filter((file) => /^main-[\w-]+\.css$/.test(file))).toHaveLength(1);

  const maintenanceHtml = await fs.readFile(path.join(outputRoot, "maintenance.html"), "utf8");
  expect(maintenanceHtml).toContain("data-page-kind=\"maintenance\"");
  expect(maintenanceHtml).toContain("NERIVA");

  const notFoundHtml = await fs.readFile(path.join(outputRoot, "404.html"), "utf8");
  expect(notFoundHtml).toContain("data-page-kind=\"not-found\"");
  expect(notFoundHtml).toContain("NERIVA");

  const headers = await fs.readFile(path.join(outputRoot, "_headers"), "utf8");
  expect(headers).toContain("/assets/site-react/*");
  expect(headers).toContain("Cache-Control: public, max-age=31536000, immutable");

  const redirects = await fs.readFile(path.join(outputRoot, "_redirects"), "utf8");
  expect(redirects).toContain("/ /poliglot-ai.html 200");
});
