import { createRequire } from 'node:module';
import { readFile, writeFile, mkdir } from 'node:fs/promises';
import { createHash } from 'node:crypto';
import { homedir } from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const sharp = createRequire(path.join(root, 'site-react/package.json'))('sharp');
const sourceDir = process.argv[2] || path.join(homedir(), 'Downloads', 'Telegram Desktop');
const targets = ['web/assets', 'site-react/public/assets/brand'];
const sources = [
  ['Neriva3-02 (2).png', 'neriva-symbol.png', 1024],
  ['Neriva2_Монтажная область 1 (2).png', 'neriva-wordmark.png', 1800],
  ['Neriva квадрат пнг-02 (2).png', 'neriva-stacked.png', 1024],
];
const sha = (buffer) => createHash('sha256').update(buffer).digest('hex');
const manifest = { date: '2026-11-15', processing: 'Trim transparent margins and resize only; original artwork and colours retained.', sources: [], outputs: [] };
const outputs = new Map();
for (const [source, filename, width] of sources) {
  const original = await readFile(path.join(sourceDir, source));
  const metadata = await sharp(original).metadata();
  const output = await sharp(original).trim().resize({ width, withoutEnlargement: true }).png({ compressionLevel: 9 }).toBuffer();
  outputs.set(filename, output);
  manifest.sources.push({ filename: source, sha256: sha(original), width: metadata.width, height: metadata.height });
}
for (const size of [192, 512]) {
  const symbol = await sharp(outputs.get('neriva-symbol.png')).resize({ width: Math.round(size * .86) }).toBuffer();
  const output = await sharp(symbol).resize(size, size, { fit: 'contain', background: { r: 0, g: 0, b: 0, alpha: 0 } }).png({ compressionLevel: 9 }).toBuffer();
  outputs.set(`neriva-icon-${size}.png`, output);
}
for (const [filename, buffer] of outputs) {
  const metadata = await sharp(buffer).metadata();
  manifest.outputs.push({ filename, sha256: sha(buffer), width: metadata.width, height: metadata.height, bytes: buffer.length });
  for (const dir of targets) {
    await mkdir(path.join(root, dir), { recursive: true });
    await writeFile(path.join(root, dir, filename), buffer);
  }
}
await writeFile(path.join(root, 'docs/specs/official-brand/brand-manifest.json'), JSON.stringify(manifest, null, 2) + '\n');
console.log(`Prepared ${outputs.size} official brand assets in app and site.`);
