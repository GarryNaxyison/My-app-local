import { readFile, writeFile } from "node:fs/promises";
import path from "node:path";
import sharp from "sharp";

const publicDir = path.resolve("public");
const sourcePath = path.join(publicDir, "assets", "brand-logo-mini.png");
const sourceBuffer = await readFile(sourcePath);
const themeColor = "#08111f";

async function transparentPng(size, filename) {
  await sharp(sourceBuffer)
    .resize(size, size, { fit: "contain" })
    .png()
    .toFile(path.join(publicDir, filename));
}

async function paddedPng(size, filename) {
  const innerSize = Math.round(size * 0.76);
  await sharp({
    create: {
      width: size,
      height: size,
      channels: 4,
      background: themeColor,
    },
  })
    .composite([
      {
        input: await sharp(sourceBuffer)
          .resize(innerSize, innerSize, { fit: "contain" })
          .png()
          .toBuffer(),
        gravity: "center",
      },
    ])
    .png()
    .toFile(path.join(publicDir, filename));
}

async function pngBuffer(size) {
  return sharp(sourceBuffer)
    .resize(size, size, { fit: "contain" })
    .png()
    .toBuffer();
}

function buildIco(entries) {
  const headerSize = 6;
  const directorySize = 16 * entries.length;
  let imageOffset = headerSize + directorySize;
  const directory = Buffer.alloc(directorySize);
  const header = Buffer.alloc(headerSize);

  header.writeUInt16LE(0, 0);
  header.writeUInt16LE(1, 2);
  header.writeUInt16LE(entries.length, 4);

  entries.forEach(({ size, buffer }, index) => {
    const offset = index * 16;
    directory.writeUInt8(size === 256 ? 0 : size, offset);
    directory.writeUInt8(size === 256 ? 0 : size, offset + 1);
    directory.writeUInt8(0, offset + 2);
    directory.writeUInt8(0, offset + 3);
    directory.writeUInt16LE(1, offset + 4);
    directory.writeUInt16LE(32, offset + 6);
    directory.writeUInt32LE(buffer.length, offset + 8);
    directory.writeUInt32LE(imageOffset, offset + 12);
    imageOffset += buffer.length;
  });

  return Buffer.concat([header, directory, ...entries.map((entry) => entry.buffer)]);
}

await transparentPng(16, "favicon-16x16.png");
await transparentPng(32, "favicon-32x32.png");
await transparentPng(48, "favicon-48x48.png");
await paddedPng(180, "apple-touch-icon.png");
await paddedPng(192, "android-chrome-192x192.png");
await paddedPng(512, "android-chrome-512x512.png");

const icoEntries = await Promise.all(
  [16, 32, 48].map(async (size) => ({ size, buffer: await pngBuffer(size) })),
);
await writeFile(path.join(publicDir, "favicon.ico"), buildIco(icoEntries));
