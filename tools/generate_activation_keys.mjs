#!/usr/bin/env node

import { appendFileSync, existsSync, readFileSync } from 'node:fs';
import { randomInt } from 'node:crypto';

const alphabet = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789';

function parseArgs(argv) {
  const args = {
    count: 1,
    days: 30,
    tier: 'premium',
    out: '',
    note: ''
  };
  for (let index = 0; index < argv.length; index += 1) {
    const arg = argv[index];
    const next = argv[index + 1];
    if (arg === '--count' || arg === '-n') {
      args.count = Number.parseInt(next, 10);
      index += 1;
    } else if (arg === '--days') {
      args.days = Number.parseInt(next, 10);
      index += 1;
    } else if (arg === '--period') {
      args.days = periodToDays(next);
      index += 1;
    } else if (arg === '--tier') {
      args.tier = String(next || 'premium').toLowerCase() === 'platinum' ? 'platinum' : 'premium';
      index += 1;
    } else if (arg === '--out' || arg === '-o') {
      args.out = next || '';
      index += 1;
    } else if (arg === '--note') {
      args.note = next || '';
      index += 1;
    } else if (arg === '--help' || arg === '-h') {
      printHelp();
      process.exit(0);
    } else {
      throw new Error(`Unknown argument: ${arg}`);
    }
  }
  if (!Number.isInteger(args.count) || args.count < 1 || args.count > 100000) {
    throw new Error('--count must be between 1 and 100000');
  }
  if (![30, 365].includes(args.days)) {
    throw new Error('--days must be 30 or 365, or use --period month|year');
  }
  return args;
}

function periodToDays(value) {
  const period = String(value || '').trim().toLowerCase();
  if (['month', 'monthly', '30', '30d', '1m'].includes(period)) return 30;
  if (['year', 'yearly', '365', '365d', '1y'].includes(period)) return 365;
  throw new Error('--period must be month or year');
}

function formatKey(compact) {
  return compact.match(/.{1,4}/g).join('-');
}

function generateKey() {
  let compact = '';
  for (let index = 0; index < 16; index += 1) {
    compact += alphabet[randomInt(alphabet.length)];
  }
  return formatKey(compact);
}

function readExistingKeys(path) {
  const keys = new Set();
  if (!path || !existsSync(path)) return keys;
  const text = readFileSync(path, 'utf8');
  for (const line of text.split(/\r?\n/)) {
    const match = line.toUpperCase().match(/[A-HJ-NP-Z2-9]{4}(?:-[A-HJ-NP-Z2-9]{4}){3}/);
    if (match) keys.add(match[0]);
  }
  return keys;
}

function printHelp() {
  console.log(`Usage:
  node tools/generate_activation_keys.mjs --count 100 --period month --out activation_keys.txt
  node tools/generate_activation_keys.mjs --count 25 --days 365 --tier premium
  node tools/generate_activation_keys.mjs --count 10 --period year --tier platinum

Output format:
  XXXX-XXXX-XXXX-XXXX 30 premium optional-note
  XXXX-XXXX-XXXX-XXXX 365 premium optional-note
  XXXX-XXXX-XXXX-XXXX 365 platinum optional-note`);
}

try {
  const args = parseArgs(process.argv.slice(2));
  const existing = readExistingKeys(args.out);
  const generated = [];
  while (generated.length < args.count) {
    const key = generateKey();
    if (existing.has(key)) continue;
    existing.add(key);
    const note = args.note ? ` ${args.note}` : '';
    generated.push(`${key} ${args.days} ${args.tier}${note}`);
  }
  const output = `${generated.join('\n')}\n`;
  if (args.out) {
    appendFileSync(args.out, output, 'utf8');
    console.log(`Generated ${generated.length} keys into ${args.out}`);
  } else {
    process.stdout.write(output);
  }
} catch (error) {
  console.error(error.message);
  console.error('Run with --help for usage.');
  process.exit(1);
}
