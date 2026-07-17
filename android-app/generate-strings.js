const fs = require('fs');
const path = require('path');

const i18nPath = 'E:\\PROJECTS\\New project\\My app local\\web-react\\src\\lib\\i18n.ts';
const outputBase = 'E:\\PROJECTS\\New project\\My app local\\android-app\\app\\src\\main\\res';

const fileContent = fs.readFileSync(i18nPath, 'utf8');

// Extract a simple Record<string, string> block
function extractDict(content, varName) {
  const regex = new RegExp(`const\\s+${varName}\\s*:\\s*Record<string,\\s*string>\\s*=\\s*\\{([\\s\\S]*?)\\n\\};`);
  const match = content.match(regex);
  if (!match) return {};
  return parseKeys(match[1]);
}

function parseKeys(block) {
  const entries = {};
  const lines = block.split('\n');
  for (const line of lines) {
    const m = line.match(/^\s*(\w+)\s*:\s*"((?:[^"\\]|\\.)*)"/);
    if (m) {
      entries[m[1]] = m[2].replace(/\\"/g, '"');
    }
  }
  return entries;
}

// Extract a locale override block - more robust approach
function extractOverrideBlock(content, varName) {
  const result = {};
  
  // Find the start of the block
  const blockStart = content.indexOf(`const ${varName}`);
  if (blockStart === -1) return result;
  
  // Find the matching closing
  let depth = 0;
  let blockContent = '';
  let inBlock = false;
  let braceCount = 0;
  
  for (let i = blockStart; i < content.length; i++) {
    if (content[i] === '{') {
      braceCount++;
      if (braceCount === 1) inBlock = true;
    }
    if (inBlock) blockContent += content[i];
    if (content[i] === '}') {
      braceCount--;
      if (braceCount === 0) break;
    }
  }
  
  // Extract each language block
  // Pattern: langCode: { ... }
  const langPattern = /(\w+):\s*\{/g;
  let langMatch;
  let langStarts = [];
  
  while ((langMatch = langPattern.exec(blockContent)) !== null) {
    langStarts.push({ lang: langMatch[1], pos: langMatch.index + langMatch[0].length });
  }
  
  for (let i = 0; i < langStarts.length; i++) {
    const lang = langStarts[i].lang;
    const startPos = langStarts[i].pos;
    
    // Find matching closing brace for this language block
    let depth = 1;
    let endPos = startPos;
    while (endPos < blockContent.length && depth > 0) {
      if (blockContent[endPos] === '{') depth++;
      if (blockContent[endPos] === '}') depth--;
      endPos++;
    }
    
    const langBlock = blockContent.substring(startPos, endPos - 1);
    const entries = {};
    
    // Parse key-value pairs
    const kvPattern = /(\w+)\s*:\s*"((?:[^"\\]|\\.)*)"/g;
    let kvMatch;
    while ((kvMatch = kvPattern.exec(langBlock)) !== null) {
      entries[kvMatch[1]] = kvMatch[2].replace(/\\"/g, '"');
    }
    
    if (Object.keys(entries).length > 0) {
      result[lang] = entries;
    }
  }
  
  return result;
}

// Extract answerFieldPlaceholderFallbacks
function extractAnswerFallbacks(content) {
  const result = {};
  const blockStart = content.indexOf('const answerFieldPlaceholderFallbacks');
  if (blockStart === -1) return result;
  
  let depth = 0;
  let blockContent = '';
  let inBlock = false;
  let braceCount = 0;
  
  for (let i = blockStart; i < content.length; i++) {
    if (content[i] === '{') {
      braceCount++;
      if (braceCount === 1) inBlock = true;
    }
    if (inBlock) blockContent += content[i];
    if (content[i] === '}') {
      braceCount--;
      if (braceCount === 0) break;
    }
  }
  
  const kvPattern = /(\w+):\s*"((?:[^"\\]|\\.)*)"/g;
  let kvMatch;
  while ((kvMatch = kvPattern.exec(blockContent)) !== null) {
    result[kvMatch[1]] = kvMatch[2].replace(/\\"/g, '"');
  }
  
  return result;
}

// Extract all dictionaries
const en = extractDict(fileContent, 'en');
const ruBase = extractDict(fileContent, 'ru');

// Extract override blocks
const localeOverrides = extractOverrideBlock(fileContent, 'localeOverrides');
const authLocaleOverrides = extractOverrideBlock(fileContent, 'authLocaleOverrides');
const v2UiLocaleOverrides = extractOverrideBlock(fileContent, 'v2UiLocaleOverrides');
const v2MenuLocaleOverrides = extractOverrideBlock(fileContent, 'v2MenuLocaleOverrides');
const notesAndAuthLocaleOverrides = extractOverrideBlock(fileContent, 'notesAndAuthLocaleOverrides');
const authPageLocaleOverrides = extractOverrideBlock(fileContent, 'authPageLocaleOverrides');
const authPrivacyLocaleOverrides = extractOverrideBlock(fileContent, 'authPrivacyLocaleOverrides');

// Extract answerFieldPlaceholderFallbacks
const answerFallbacks = extractAnswerFallbacks(fileContent);

const coreLanguages = ['ru','en','es','de','fr','it','zh','ja','ko','tg','uz','tt','hy','kk','ky','ka','uk','pl','ro','pt'];
const fallbackLanguages = ['ar','bn','cs','el','hi','hu','id','nl','sv','ta','te','th','tl','tr','vi'];

function escapeXml(str) {
  if (!str) return '';
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/'/g, '&apos;')
    .replace(/"/g, '&quot;')
    .replace(/\n/g, '\\n');
}

function buildMergedDict(lang) {
  const merged = { ...en };
  
  // Add ru base if lang is ru
  if (lang === 'ru') {
    Object.assign(merged, ruBase);
  }
  
  // Apply localeOverrides
  if (localeOverrides[lang]) {
    Object.assign(merged, localeOverrides[lang]);
  }
  
  // Apply authLocaleOverrides
  if (authLocaleOverrides[lang]) {
    Object.assign(merged, authLocaleOverrides[lang]);
  }
  
  // Apply v2UiLocaleOverrides
  if (v2UiLocaleOverrides[lang]) {
    Object.assign(merged, v2UiLocaleOverrides[lang]);
  }
  
  // Apply v2MenuLocaleOverrides
  if (v2MenuLocaleOverrides[lang]) {
    Object.assign(merged, v2MenuLocaleOverrides[lang]);
  }
  
  // Apply notesAndAuthLocaleOverrides
  if (notesAndAuthLocaleOverrides[lang]) {
    Object.assign(merged, notesAndAuthLocaleOverrides[lang]);
  }
  
  // Apply authPageLocaleOverrides
  if (authPageLocaleOverrides[lang]) {
    Object.assign(merged, authPageLocaleOverrides[lang]);
  }
  
  // Apply authPrivacyLocaleOverrides
  if (authPrivacyLocaleOverrides[lang]) {
    Object.assign(merged, authPrivacyLocaleOverrides[lang]);
  }
  
  // Apply answerFieldPlaceholderFallbacks
  if (answerFallbacks[lang]) {
    merged.answer_field_placeholder = answerFallbacks[lang];
  }
  
  return merged;
}

function generateXml(dict) {
  const sortedKeys = Object.keys(dict).sort();
  let xml = '<?xml version="1.0" encoding="utf-8"?>\n';
  xml += '<resources>\n';
  xml += '    <string name="app_name">NERIVA</string>\n';
  
  for (const key of sortedKeys) {
    const value = escapeXml(dict[key]);
    if (value) {
      xml += `    <string name="${key}">${value}</string>\n`;
    }
  }
  
  xml += '</resources>\n';
  return xml;
}

function ensureDir(dir) {
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }
}

const keyCounts = {};

// Generate for all 35 languages
const allLanguages = [...coreLanguages, ...fallbackLanguages];

for (const lang of allLanguages) {
  const dict = buildMergedDict(lang);
  const xml = generateXml(dict);
  
  const valuesDir = lang === 'en' 
    ? path.join(outputBase, 'values')
    : path.join(outputBase, `values-${lang}`);
  
  ensureDir(valuesDir);
  
  const filePath = path.join(valuesDir, 'strings.xml');
  fs.writeFileSync(filePath, xml, 'utf8');
  
  const keyCount = Object.keys(dict).length + 1; // +1 for app_name
  keyCounts[lang] = keyCount;
  
  console.log(`${lang}: ${keyCount} keys -> ${filePath}`);
}

// Count how many keys are translated (not English) per language
console.log('\n=== TRANSLATION COVERAGE ===');
for (const lang of allLanguages) {
  if (lang === 'en') continue;
  const dict = buildMergedDict(lang);
  let translated = 0;
  let total = 0;
  for (const [key, value] of Object.entries(dict)) {
    if (en[key] && value !== en[key]) {
      translated++;
    }
    total++;
  }
  const pct = total > 0 ? Math.round(translated / total * 100) : 0;
  console.log(`  ${lang}: ${translated}/${total} translated (${pct}%)`);
}

console.log('\n=== SUMMARY ===');
console.log(`Total languages: ${allLanguages.length}`);
console.log('\nKey counts per language:');
for (const [lang, count] of Object.entries(keyCounts)) {
  console.log(`  ${lang}: ${count} strings`);
}
