import { expect, test } from "@playwright/test";

const requiredDeployAssets = [
  "assets/site-i18n.js",
  "assets/site-phrases.js",
  "assets/privacy-policy-i18n.js",
  "assets/legal-documents-i18n.js",
  "assets/scenarios/travel-ai-tutor.jpg",
  "assets/scenarios/work-ai-tutor.jpg",
  "assets/scenarios/exam-ai-tutor.jpg",
  "assets/scenarios/speaking-ai-tutor.jpg",
  "assets/seo/poliglot-ai-og-ru.jpg",
  "assets/seo/poliglot-ai-og-en.jpg",
  "assets/testimonials/anna.jpg",
  "assets/testimonials/marat.jpg",
  "assets/testimonials/sofia.jpg",
  "assets/site-react",
] as const;

test("built public-site folder contains every asset required by the deploy package", async () => {
  const fs = await import("node:fs/promises");
  const path = await import("node:path");
  const outputRoot = path.resolve(process.cwd(), "..", "Сайт полиглота для бота");

  for (const assetPath of requiredDeployAssets) {
    const stat = await fs.stat(path.join(outputRoot, assetPath));
    if (stat.isDirectory()) {
      const entries = await fs.readdir(path.join(outputRoot, assetPath));
      expect(entries.length, `${assetPath} directory should not be empty in build output`).toBeGreaterThan(0);
    } else {
      expect(stat.size, `${assetPath} should not be empty in build output`).toBeGreaterThan(100);
    }
  }

  const htmlFiles = ["poliglot-ai.html", "privacy.html", "terms.html", "agreement.html", "consent.html", "404.html", "maintenance.html"];
  for (const htmlFile of htmlFiles) {
    const stat = await fs.stat(path.join(outputRoot, htmlFile));
    expect(stat.size, `${htmlFile} should be present in build output`).toBeGreaterThan(100);
  }

  const referencedBundles = new Set<string>();
  for (const htmlFile of htmlFiles) {
    const html = await fs.readFile(path.join(outputRoot, htmlFile), "utf8");
    for (const match of html.matchAll(/\/assets\/site-react\/([^"'<>\s)]+)/g)) {
      referencedBundles.add(match[1]);
    }
  }

  const bundleFiles = await fs.readdir(path.join(outputRoot, "assets", "site-react"));
  expect(bundleFiles.sort(), "assets/site-react should only contain bundles referenced by the built HTML files").toEqual(
    [...referencedBundles].sort(),
  );
  expect(bundleFiles.filter((file) => /^main-[\w-]+\.js$/.test(file))).toHaveLength(1);
  expect(bundleFiles.filter((file) => /^main-[\w-]+\.css$/.test(file))).toHaveLength(1);
});
