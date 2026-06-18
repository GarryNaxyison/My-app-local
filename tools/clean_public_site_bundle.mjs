import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

const PUBLIC_SITE_DIR = "\u0421\u0430\u0439\u0442 \u043f\u043e\u043b\u0438\u0433\u043b\u043e\u0442\u0430 \u0434\u043b\u044f \u0431\u043e\u0442\u0430";
const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const assetsDir = path.resolve(root, PUBLIC_SITE_DIR, "assets");
const bundleDir = path.resolve(assetsDir, "site-react");

if (!bundleDir.startsWith(`${assetsDir}${path.sep}`) || path.basename(bundleDir) !== "site-react") {
  throw new Error(`Refusing to clean unexpected public-site bundle path: ${bundleDir}`);
}

await fs.rm(bundleDir, { recursive: true, force: true });
await fs.mkdir(bundleDir, { recursive: true });
console.log(`Cleaned public-site bundle output: ${bundleDir}`);
