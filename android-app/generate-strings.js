const fs = require('fs');
const path = require('path');
const ts = require('../web-react/node_modules/typescript');
const vm = require('vm');

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

function loadCanonicalWebLocales() {
  const source = `${fileContent}\nexport { en, localeOverrides };`;
  const compiled = ts.transpileModule(source, {
    compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020 },
  }).outputText;
  const module = { exports: {} };
  vm.runInNewContext(compiled, {
    module,
    exports: module.exports,
    console,
    Object,
    Array,
    String,
    Number,
    Boolean,
    Math,
    Set,
    Map,
  });
  return module.exports;
}

const { en, localeOverrides } = loadCanonicalWebLocales();
const allLanguages = ['ru','en','es','de','fr','it','zh','ja','ko','tg','uz','tt','hy','kk','ky','ka','uk','pl','ro','pt','ar','bn','cs','el','hi','hu','id','nl','sv','ta','te','th','tl','tr','vi'];
const baseDictionary = { ...en };
for (const locale of allLanguages) {
  for (const [key, value] of Object.entries(localeOverrides[locale] || {})) {
    if (baseDictionary[key] == null && value) {
      baseDictionary[key] = value;
    }
  }
}

function escapeXml(str) {
  if (!str) return '';
  return str
    .replace(/\\/g, '\\\\')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/'/g, "\\" + "'")
    .replace(/"/g, '&quot;')
    .replace(/\n/g, '\\n');
}

function buildMergedDict(lang) {
  return { ...baseDictionary, ...(localeOverrides[lang] || {}) };
}

function generateXml(dict) {
  const sortedKeys = Object.keys(dict).sort();
  const androidNames = new Map();
  let xml = '<?xml version="1.0" encoding="utf-8"?>\n';
  xml += '<resources>\n';
  xml += '    <string name="app_name">NERIVA</string>\n';
  
  for (const key of sortedKeys) {
    const androidKey = key.replace(/[^a-z0-9_]/g, '_').replace(/^[^a-z_]/, '_');
    const previousKey = androidNames.get(androidKey);
    if (previousKey && previousKey !== key) {
      throw new Error(`Android resource key collision: ${previousKey} and ${key} -> ${androidKey}`);
    }
    androidNames.set(androidKey, key);
    const value = escapeXml(dict[key]);
    if (value) {
      xml += `    <string name="${androidKey}">${value}</string>\n`;
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
