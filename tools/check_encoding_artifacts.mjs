#!/usr/bin/env node

import { readFileSync } from 'node:fs';
import { join } from 'node:path';

const root = process.cwd();

const targetFiles = [
  'web/index.html',
  'web_api.go',
  'bot.go',
  'i18n.go',
  'premium_i18n.go',
  'prompts.go',
  'CHANGELOG.md',
  'DOCUMENTATION.md',
  'PROJECT_TRACKING.md',
  'web-react/src/App.tsx',
  'web-react/src/lib/i18n.ts',
  'Сайт полиглота для бота/poliglot-ai.html',
  'Сайт полиглота для бота/privacy.html',
  'Сайт полиглота для бота/terms.html',
  'Сайт полиглота для бота/assets/site-i18n.js'
];

const badFragments = [
  '\u0420\u2014', '\u0420\u201c', '\u0420\u2019', '\u0420\u040e',
  '\u0420\u0406', '\u0420\u045c', '\u0420\u045b', '\u0420\u0491',
  '\u0420\u00a0', '\u0420\u00b0', '\u0420\u00b5', '\u0420\u00bb',
  '\u0420\u0457', '\u0420\u0455', '\u0420\u0451', '\u0420\u0405',
  '\u0420\u2116', '\u0420\u00b1', '\u0420\u00bc', '\u0420\u00bd',
  '\u0421\u0403', '\u0421\u201a', '\u0421\u2039',
  '\u0421\u0402', '\u0421\u040c', '\u0421\u2018', '\u0421\u201c',
  '\u0421\u2021', '\u0421\u2030', '\u0421\u0405', '\u0421\u0459',
  '\u0421\u045a', '\u0421\u045c', '\u0421\u0458', '\u0421\u045f',
  '\u0432\u0402', '\u0440\u045f', '\ufffd'
];

function printable(value) {
  return value.replace(/[^\x20-\x7e]/g, (char) => `\\u${char.codePointAt(0).toString(16).padStart(4, '0')}`);
}

function shouldSkipLine(file, line) {
  if (file === 'web/index.html' && (
    line.includes('function mojibakeScore') ||
    line.includes('function hasEncodingResidue') ||
    line.includes('repairMojibakeText') ||
    line.includes('/[\ufffd]/') ||
    line.includes('windows-1251') ||
    line.includes('knownMojibake')
  )) {
    return true;
  }
  if (file.endsWith('site-i18n.js') && /mojibake|windows-1251|[\u00c2\u00c3\u00d0\u00d1]/i.test(line)) {
    return true;
  }
  if (file.endsWith('web-react/src/lib/i18n.ts') && (
    line.includes('mojibake') ||
    line.includes('windows-1251') ||
    line.includes('decodeMojibakeRun') ||
    line.includes('mojibakeTail') ||
    line.includes('pairRun') ||
    line.includes('singlePair') ||
    line.includes('mojibakeScore') ||
    line.includes('\\u0420') ||
    line.includes('\\u0421')
  )) {
    return true;
  }
  return false;
}

const hits = [];

for (const file of targetFiles) {
  let text = '';
  try {
    text = readFileSync(join(root, file), 'utf8');
  } catch {
    continue;
  }
  text.split(/\r?\n/).forEach((line, index) => {
    if (shouldSkipLine(file, line)) return;
    const found = badFragments.filter((fragment) => line.includes(fragment));
    if (found.length) {
      hits.push(`${file}:${index + 1}: ${printable(line.trim()).slice(0, 260)}`);
    }
  });
}

if (hits.length) {
  console.error('Encoding artifacts found:');
  hits.forEach((hit) => console.error(hit));
  process.exit(1);
}

console.log('Encoding artifacts check passed.');
