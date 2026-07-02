import { createWriteStream } from 'node:fs';
import { execFile } from 'node:child_process';
import { mkdir, readFile, writeFile } from 'node:fs/promises';
import { dirname, join } from 'node:path';
import { pipeline } from 'node:stream/promises';
import { promisify } from 'node:util';
import { randomUUID } from 'node:crypto';

const COMFY_BASE = process.env.COMFY_BASE || 'http://127.0.0.1:8000';
const OUTPUT_DIR = process.env.POLIGLOT_ASSET_DIR || 'web/assets';
const RAW_TRANSPARENT_DIR = process.env.POLIGLOT_RAW_TRANSPARENT_DIR || 'tmp/comfy-transparent-raw';
const REPORT_DIR = process.env.POLIGLOT_ASSET_REPORT_DIR || 'tmp/asset-validation';
const PYTHON = process.env.PYTHON || 'python';
const MODEL = process.env.COMFY_CHECKPOINT || 'sd3.5_large_fp8_scaled.safetensors';
const FLUX2_DIFFUSION_MODEL = process.env.COMFY_FLUX2_DIFFUSION || 'flux2_dev_fp8mixed.safetensors';
const FLUX2_TEXT_ENCODER = process.env.COMFY_FLUX2_TEXT_ENCODER || 'mistral_3_small_flux2_bf16.safetensors';
const FLUX2_VAE = process.env.COMFY_FLUX2_VAE || 'flux2-vae.safetensors';
const UPSCALE_MODEL = process.env.COMFY_UPSCALE_MODEL || 'RealESRGAN_x4plus.pth';
const FALLBACK_MODELS = [
  'flux2-dev.safetensors',
  'flux2_dev_fp8mixed.safetensors',
  'flux1-dev-fp8.safetensors',
  'Juggernaut-XL_v9_RunDiffusionPhoto_v2.safetensors',
  'sd_xl_base_1.0.safetensors',
];
const THEMED_GROUPS = new Set(['plans', 'headers', 'icons', 'backgrounds', 'panels']);
const TRANSPARENT_GROUPS = new Set(['icons', 'awards', 'logos']);
const execFileAsync = promisify(execFile);
const THEMES = {
  light: 'light theme variant, white gallery studio, bright ceramic glass, sapphire blue, teal, mint, restrained warm gold, clean shadows',
  dark: 'dark theme variant, graphite indigo gallery studio, luminous sapphire blue, teal, mint, restrained warm gold rim light, deep contrast',
};

function isSD3Model(model = MODEL) {
  return /(^|[_-])sd3|sd3\.5|stable[-_ ]diffusion[-_ ]3/i.test(model);
}

function isFluxModel(model = MODEL) {
  return /(^|[_-])flux/i.test(model);
}

function isFlux2Model(model = MODEL) {
  return /(^|[_-])flux2|flux\.2/i.test(model);
}

function assetSlug(file) {
  return file
    .replace(/\.png$/i, '')
    .replace(/-(light|dark)$/i, '')
    .replace(/^(plan|header|icon)-/i, '');
}

const flux2Subjects = {
  home: 'main menu home: large freestanding house-shaped learning hub, roof arch, central sapphire globe core, mint route ring wrapped around the house, instantly reads as home hub',
  lesson: 'new lesson start: large freestanding lesson gateway, forward path ribbon passing through a sapphire learning lens, one start crystal, reads as begin lesson',
  practice: 'conversation practice: two oversized blank chat-bubble sculptures stacked diagonally, upright microphone capsule between them, thick voice waveform loop around the bubbles, fills square height and width, absolutely not a horizontal strip',
  shadowing: 'shadowing speech trainer: large upright microphone capsule surrounded by open circular echo rings, one sapphire pronunciation core, mint sound waves returning back to the learner, reads as listen then repeat aloud, no rectangular echo panel',
  roleplay: 'AI roleplay scenes: two separated polished theatre-mask shaped blank speech capsules orbiting one sapphire dialogue core, mint route arcs, warm-gold scene spark, reads as situational dialogue practice, no human faces, no cube, no box',
  pronunciation: 'pronunciation analysis: one sapphire sound core with precise open waveform heatmap rings and a small detached calibration lens, mint acoustic arcs, reads as voice accuracy and speech clarity, no microphone duplicate, no square plate',
  offline: 'offline mini-decks: one open crescent-shaped sapphire learning capsule, detached mint pause ribbon, two separated pearl memory beads, lots of empty transparent gaps between parts, reads as portable saved practice, no cube, no box, no container, no square silhouette, no stack, no tile, no cards, no phone',
  phrasebook: 'favorite phrases and personal phrasebook: one freestanding vertical bookmark ribbon crossing a blank speech capsule, sapphire phrase gem, mint memory arc, open airy silhouette, no book pages, no printed marks, no cube, no box',
  progress: 'learning progress: large ascending stair-and-route sculpture, three rising crystal steps inside a broad progress ring, upward motion, no numbers',
  awards: 'achievements: single large trophy cup with broad handles, sapphire prize gem, warm-gold celebration ring, no plaque',
  vocabulary: 'vocabulary dictionary: large closed blank dictionary-vault sculpture with sapphire bookmark ribbon, memory beads orbiting around it, no letters, no pages, no drawers',
  words: 'learn words: large cluster of three blank word-card gems arranged as a flower, sapphire seed core, mint memory route arc, no letters',
  'word-game': 'word review game: two oversized interlocking puzzle pieces forming a circular repeat-arrow silhouette, reward gem in the center, no square base',
  spelling: 'spelling practice: large diagonal polished check mark inside a broad proofreading ring, small stylus crystal crossing the lower edge, reads as spelling correction',
  level: 'language level assessment: open freestanding mastery meter made of three separate floating vertical glass bars, small sapphire calibration pointer, broken mint alignment halo with transparent gaps, no solid disk, no target plate, no scale, no ticks, no numbers',
  leaderboard: 'leaderboard: three large ascending podium crystals with one central winner gem and laurel-like mint arc, no numbers, no plaques',
  premium: 'premium upgrade: wide freestanding premium shield with broad side wings, large crown integrated across the top, sapphire premium core, warm-gold neural ring wrapping around the whole shield, fills full width and height, no pass card',
  limits: 'daily limits: wide freestanding twin reservoir battery, two large side-by-side glass capacity tanks connected by a sapphire valve, broad mint usage ring wrapping around both tanks, fills full width and height, no vertical hourglass, no gauge marks, no numbers',
  mistakes: 'mistake correction: large repair loop around a polished correction gem, eraser-shaped crystal, friendly check mark, no document',
  tools: 'tools: large triangular tool cluster, microphone capsule, camera lens, translation prism, and voice waveform ribbon locked together, fills square, not a horizontal strip',
  referral: 'referrals: large connected network nodes tied into a gift-ribbon sculpture with reward gem, no invitation card',
  settings: 'settings: large premium gear dial with two oversized toggle capsules and secure profile ring, quiet utility sculpture',
  free: 'beautiful starter access portal, sapphire learning gem, polished mint route ribbon, three pearl knowledge gems, one refined check crystal, inviting premium free-tier object, not cheap, not empty, no compass, no dial',
  platinum: 'mastery crown observatory, silver and warm-gold crown, sapphire mastery gem, layered calibration rings',
  'app-background': 'full-screen web app backdrop for desktop and mobile, immersive language intelligence atmosphere made from blank frosted-glass sound capsules, pure waveform light ribbons, sapphire learning cores pushed toward the outer edges, mint translation route arcs, subtle achievement gems, calm empty center corridor, mobile-safe vertical crop, no chat messages, no paragraph marks, no interface text',
  'panel-progress': 'tiered sapphire crystal stair sculpture, mint silk route ribbon, small faceted milestone gems, soft pearl base, no dial, no gauge, no clock face, no measurement marks',
  'panel-limits': 'stacked translucent mint glass reservoir capsules, warm-gold ribbon, sapphire bead accents, calm utility sculpture, no dial, no gauge, no measuring scale, no tick marks',
  'panel-account': 'decorative account sidebar backdrop only, smooth secure pearl capsule, sapphire shield gem, premium access core, warm-gold protection ring, mint privacy halo, object detail bleeds from lower-right edge, upper-left remains quiet blank studio space for real account text, no account card, no app UI, no portrait, no avatar, no profile icon, no ID card, no menu rows',
  'brand-logo': 'NERIVA brand mark for a language-learning companion: one freestanding polyglot gem-ring sculpture, sapphire speech capsule core, smooth mint translation orbit ring, small warm-gold achievement star only, subtly suggests listening, speaking, learning path and mastery, iconic simple silhouette, no compass arrow, no triangle mark, no printed symbol, no letters, no text',
  'brand-logo-hero-bg': 'large thematic language-learning brand backdrop without any logo or typography: immersive atmosphere with blank speech capsules, smooth translation route arcs, sapphire learning cores, continuous mint waveform ribbons and achievement crystal glow around the edges, completely empty clean center reserved for a separate real transparent logo overlay, no center emblem, no isolated glyph shapes, no arrowheads, no letters, no text, no fake UI',
  'auth-login-hero': 'login page hero background for an AI language-learning service: premium language gateway with only abstract blank speech capsules, sapphire account core without any symbol inside it, mint route arcs, subtle privacy shield halo and daily learning orbit, clean empty left-side negative space for the real login form, rich right-side object detail, absolutely no text, no labels, no letters on rings, no brand names, no Telegram logo, no paper-plane symbol, no icons, no UI screenshot',
};

const awardTiers = [
  'tier 1 steel trophy: simple brushed steel cup, small sapphire learning gem, beginner silhouette',
  'tier 2 dark steel trophy: stronger blackened steel cup, mint rim, slightly taller base',
  'tier 3 copper trophy: warm copper cup, sapphire center gem, cleaner handles',
  'tier 4 polished copper trophy: polished copper and mint glass handles, brighter highlights',
  'tier 5 bronze trophy: classic bronze cup, sapphire shield gem, confident entry-rank form',
  'tier 6 antique bronze trophy: antique bronze cup, layered base, subtle mint route halo',
  'tier 7 silver trophy: bright silver cup, larger handles, faceted sapphire core',
  'tier 8 sapphire silver trophy: silver cup with sapphire inlays and floating mint ring',
  'tier 9 gold trophy: warm gold cup, sapphire prize jewel, richer pedestal',
  'tier 10 royal gold trophy: royal gold cup, crown-like handles, glowing achievement halo',
  'tier 11 platinum trophy: platinum cup, clean high-end silhouette, crystal base',
  'tier 12 platinum sapphire trophy: platinum cup with large sapphire heart and elegant route ring',
  'tier 13 emerald trophy: platinum and emerald crystal accents, more complex handles',
  'tier 14 ruby trophy: gold and ruby trophy, strong heroic silhouette, layered gem base',
  'tier 15 amethyst trophy: platinum-amethyst cup, arcane language mastery halo, premium facets',
  'tier 16 diamond trophy: clear diamond crystal trophy, platinum frame, radiant sapphire core',
  'tier 17 obsidian gold trophy: black obsidian and gold cup, dramatic luxury rim light',
  'tier 18 celestial crystal trophy: translucent celestial crystal cup, floating orbit rings, luminous base',
  'tier 19 mythril trophy: mythril silver-blue trophy, intricate handles, rare mastery aura',
  'tier 20 legendary polyglot crystal trophy: grand legendary trophy, sapphire diamond crown, radiant translation halo, most prestigious tier',
];

function awardTierPrompt(level) {
  return awardTiers[Math.max(0, Math.min(19, Number(level) - 1))] || awardTiers[0];
}

function flux2AssetPrompt(asset) {
  const slug = assetSlug(asset.file);
  let subject = flux2Subjects[slug] || 'abstract speech capsule, voice ribbon, sapphire learning gem, mint route arcs';
  if (asset.group === 'plans' && slug === 'premium') {
    subject = 'sapphire learning core, two blank speech capsules, mint voice ribbon halo, polished titanium base, small warm-gold route nodes, no crown';
  }
  if (asset.group === 'awards') {
    subject = awardTierPrompt(slug.replace(/^award-/, ''));
  }
  const theme = TRANSPARENT_GROUPS.has(asset.group)
    ? (asset.theme === 'dark'
      ? 'dark theme material variant on the freestanding object only: graphite enamel, luminous sapphire, mint rim light, background remains flat #00ff00'
      : 'light theme material variant on the freestanding object only: pearl ceramic, sapphire blue, mint glass, warm-gold accents, background remains flat #00ff00')
    : (asset.theme === 'dark'
      ? 'dark graphite indigo studio, luminous sapphire and mint rim light, deep contrast, no visible writing, no fake writing, no glyph-like marks, no UI text'
      : 'bright white gallery studio, pearl ceramic glass, sapphire blue, mint, soft warm-gold accents, no visible writing, no fake writing, no glyph-like marks, no UI text');
  const transparentCanvas = 'perfectly flat solid #00ff00 chroma-key background for transparent PNG extraction, uniform green background only, no shadow, no floor, no reflection, no gradient, no texture, no white background, do not use green inside the subject';
  const canvas = asset.group === 'awards'
    ? `square trophy cutout source, exactly one freestanding centered trophy cup object, trophy fills 90-94 percent of the square canvas by height and uses broad handles to fill width, same camera distance across all twenty trophies, no white box, no rounded square app tile, no glass tile, no frame, no card, no plaque slab, no background scene, ${transparentCanvas}`
    : asset.group === 'panels'
      ? 'compact decorative right-side sidebar backdrop, pure abstract object cluster anchored to the right and bottom edges, object may extend beyond the frame for CSS masked card use, upper-left remains plain empty studio space for real HTML text, not a UI screenshot, not an account card, no card chrome, no profile rows, no menu rows'
    : asset.group === 'backgrounds'
    ? 'wide immersive web app backdrop, full-bleed composition for both desktop landscape and mobile portrait crop, readable empty center and center-left, soft abstract object detail around outer edges, no central clutter, no repeated pattern, no tiling, no chat UI bubbles, no card UI, no pseudo-text strokes, blank glass surfaces only'
    : asset.group === 'headers'
    ? 'single wide banner composition, one cohesive object cluster on the right, left side intentionally empty and clean for real HTML text, no tiled pattern, no repeated duplicate objects side by side, edges may fade softly for CSS masking'
    : asset.group === 'logos'
      ? `square transparent logo source, exactly one freestanding centered brand-mark sculpture, iconic simple silhouette, object fills 88-92 percent of the square canvas by width and height, same luxury education technology style as all menu icons, no tiny centered symbol, no empty border, no square tile, no rounded app tile, no badge, no frame, no text, no letters, no small side marks, no arrow, no triangle, no glyph-like marks, ${transparentCanvas}`
    : asset.group === 'icons'
      ? `square menu button cutout source, one freestanding function-specific centered symbol-sculpture, object fills 90-94 percent of the square canvas with both substantial width and substantial height, same scale and camera across the icon set, transparent-background cutout after chroma removal, no tiny centered subject, no empty border, no long thin horizontal strip, no row of small objects, no square backplate, no rounded square app tile, no button background, no glass tile, no frame, no card, no panel behind the subject, no object inside a square panel, no translucent rectangular pane, no rounded rectangle sheet, no slab or plate behind the object, use compact diagonal, circular, or stacked composition when the subject has multiple parts, ${transparentCanvas}`
      : '16:9 subscription card artwork, one centered hero object, breathable blank background';
  const transparentRules = TRANSPARENT_GROUPS.has(asset.group)
    ? 'final image is meant for local background removal into transparent PNG, keep the green background perfectly removable, subject has crisp closed silhouette and generous antialiased edges, the visible object must be the only non-green content'
    : '';
  return [
    'OBJECT ONLY PRODUCT RENDER, absolutely no typography anywhere',
    'all surfaces are blank and unmarked, no readable text, no pseudo text, no letters, no numbers, no pictograms, no printed icons, no symbols on surfaces, no logos, no watermark, no labels, no captions, no engravings, no plaques, no etched micro marks',
    'no screens, no dashboard, no browser window, no phone, no tablet, no documents, no papers, no folders, no file browser, no rounded square icon tile, no rectangular tile, no card background, no square plate, no app tile, no decorative backplate',
    'no measurement ticks, no clock face, no tiny secondary interface icons, no avatar silhouette, no list rows, no paragraphs, no interface chrome',
    'luxury education technology object, museum-grade studio lighting, polished ceramic, frosted glass, brushed titanium, sapphire crystal, refined industrial design',
    canvas,
    subject,
    transparentRules,
    theme,
  ].join(', ');
}

function effectivePrompt(asset) {
  if (TRANSPARENT_GROUPS.has(asset.group) || asset.group === 'headers') {
    return flux2AssetPrompt(asset);
  }
  return isFlux2Model() ? flux2AssetPrompt(asset) : asset.prompt;
}

function themedFile(file, theme) {
  const dot = file.lastIndexOf('.');
  return dot === -1 ? `${file}-${theme}` : `${file.slice(0, dot)}-${theme}${file.slice(dot)}`;
}

function themedPrefix(prefix, theme) {
  return `${prefix}-${theme}`;
}

const basePositive = [
  'luxury editorial 3D product design asset for an intelligent language-learning companion',
  'high-end SaaS education visual system, refined industrial design, museum-grade studio product photography, cinematic Octane style render',
  'single authored design language, polished ceramic, frosted glass, brushed titanium, sapphire crystal, mint light, precise soft shadows',
  'language learning is expressed through abstract speech capsules, voice waveform ribbons, translation route arcs, calibration rings, crafted learning instruments, small achievement gems only when specified',
  'object-only composition, one clear hero object per asset, bespoke silhouette, luxury tactile materials, no generic interface mockups',
  'controlled palette: deep indigo, graphite, pearl white, sapphire blue, teal, mint, selective warm gold, no childish rainbow colors',
  'crisp high quality, expensive, elegant, designer-grade, restrained, memorable, not stock looking, not clipart, not cartoon',
  'clean negative space and readable composition for real HTML text overlays, surfaces are completely blank and unmarked',
  'no readable surfaces, no visible text, no letters, no words, no logo, no watermark',
].join(', ');

const negative = [
  'text, letters, numbers, words, logo, watermark, signature',
  'do not write Premium, Free, Platinum, Poliglot, AI, plan name, tier name, title, heading, caption, label',
  'engraved label, small pedestal text, nameplate, plaque, caption on object base',
  'pseudo text, fake letters, fake logo, brand name, typography, glyphs, runes, character marks, scribbles, small printed marks, AI letters, CEFR letters, infinity symbol, mathematical signs',
  'folders, file folders, file explorer, desktop folders, office folders, rows of files, binders, shelves, archive boxes, stacked documents',
  'messy details, clutter, childish cartoon, stock photo, random flags, country flags',
  'people, faces, hands, fingers, human anatomy, objects held by hands, distorted UI, unreadable typography',
  'smartphone screen, laptop screen, monitor screen, tablet, browser window, website mockup, dashboard screenshot, interface panel with rows, menu bar, buttons with text',
  'documents, certificates, forms, ID card, pass card, paper sheet, plaque, label surface, book title, printed book cover, notebook, open book, textbook pages',
  'fake interface text, gibberish characters, tiny text marks, letter-like strokes, labels, captions, inscriptions, printed symbols',
  'small centered subject, tiny icon, empty border, long horizontal strip, row of small objects, repeated identical icons, repeated tiles',
  'rounded square app icon, square tile, rectangular tile, square backplate, glass backplate, button background, card background, frame around object, object printed on a panel',
  'single orb, single sphere, ball, generic product render, unrelated abstract sculpture',
  'natural crystal cave, raw mineral cluster, geode, fantasy crystal landscape, gemstone pile, only crystals with no learning concept',
  'dark muddy image, excessive gradients, low resolution, blurry edges',
  'jpeg artifacts, bad composition, noisy background, overcomplicated scene',
].join(', ');

const baseAssets = [
  {
    group: 'plans',
    file: 'plan-free.png',
    prefix: 'poliglot_ai/plan-free',
    width: 1280,
    height: 768,
    baseWidth: 960,
    baseHeight: 576,
    seed: 2205205601,
    prompt: [
      basePositive,
      '16:9 card background for free starter access, premium welcoming starter atelier, beautiful generous free-tier visual, not cheap, not empty',
      'one large sapphire learning gem inside a smooth open portal arch, polished mint route ribbon, three pearl knowledge gems, one refined check crystal, soft golden beginner path arc, calm and accessible luxury',
      'no compass, no dial, no gauge, no scale, no ticks, no north mark, no direction letters, no book, no paper, no folder, no screen, no written content, no flashcards, no printable surfaces, no tiny decorative marks, no sad minimal object',
      'central clean object cluster with breathable empty space, elegant enough to sit next to paid plans',
    ].join(', '),
  },
  {
    group: 'plans',
    file: 'plan-premium.png',
    prefix: 'poliglot_ai/plan-premium',
    width: 1280,
    height: 768,
    baseWidth: 960,
    baseHeight: 576,
    seed: 2205203602,
    prompt: [
      basePositive,
      '16:9 card background for the advanced subscription tier, learning core concept',
      'one luminous sapphire neural core, two blank speech capsules, sculptural voice waveform halo, floating mint route nodes, polished titanium base with no engraving',
      'smart powerful daily learning upgrade, dynamic but clean, high-end subscription visual, completely blank background with no heading',
      'no cards, no folders, no documents, no frame, no phone, no screen, no written content, no inscriptions, no flat printable surfaces',
    ].join(', '),
  },
  {
    group: 'plans',
    file: 'plan-platinum.png',
    prefix: 'poliglot_ai/plan-platinum',
    width: 1280,
    height: 768,
    baseWidth: 960,
    baseHeight: 576,
    seed: 2205202603,
    prompt: [
      basePositive,
      '16:9 card background for the elite subscription tier, mastery observatory concept',
      'one silver and warm-gold achievement crown sculpture inside layered calibration rings, faceted sapphire mastery gem, elegant milestone arc, luxury pedestal',
      'elite focused ambitious mood, tasteful luxury, strong central object with elegant depth',
      'no cards, no folders, no documents, no phone, no screen, no tablet, no written content, no dashboard panels, no inscriptions, no flat printable surfaces',
    ].join(', '),
  },
  {
    group: 'headers',
    file: 'header-home.png',
    prefix: 'poliglot_ai/header-home',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205205610,
    prompt: [basePositive, 'wide home and main-menu header background, beautiful language intelligence command center object, right-side sapphire learning core with speech capsules, voice ribbon, translation route arcs, knowledge crystals, left side clean negative space, premium but friendly, no screens, no folders, no markings'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-lesson.png',
    prefix: 'poliglot_ai/header-lesson',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202611,
    prompt: [basePositive, 'wide header for new language lesson, right-side lesson atelier object: luminous learning lens, folded ceramic ribbon, three knowledge crystals, calm high-end education, left side negative space, no book, no pages, no screens, no written marks'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-practice.png',
    prefix: 'poliglot_ai/header-practice',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202612,
    prompt: [basePositive, 'wide header for conversation practice, right-side speech duet sculpture, two blank chat capsules, sculptural microphone waveform helix, tiny route light nodes, left side negative space, no screens, no written marks'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-shadowing.png',
    prefix: 'poliglot_ai/header-shadowing',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205203211,
    prompt: [basePositive, 'wide header for shadowing pronunciation training, right-side large microphone capsule facing a translucent echo-wave mirror, sapphire pronunciation core, mint sound rings and returning voice ribbon, feels like listen then repeat aloud, left side negative space, no face, no headphones, no screen, no letters, no text'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-progress.png',
    prefix: 'poliglot_ai/header-progress',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202613,
    prompt: [basePositive, 'wide header for learning progress, right-side achievement chronograph with glowing progress rings, rising crystal fins, clean milestone path, refined analytics mood, left side negative space, no dashboard panels, no numbers, no written marks'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-awards.png',
    prefix: 'poliglot_ai/header-awards',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202614,
    prompt: [basePositive, 'wide header for achievements and awards, right-side elegant trophy gem sculpture, smooth medal ring without plaque, faceted milestone gems, subtle celebration sparks, gold and deep blue accents, no inscriptions, no certificates'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-vocabulary.png',
    prefix: 'poliglot_ai/header-vocabulary',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205204615,
    prompt: [basePositive, 'wide header for vocabulary library, right-side lexicon oracle instrument with floating glass orbit rings, sapphire core, mint route halo, round unmarked glass beads, soft route glow, teal and blue, no drawers, no tablets, no plaques, no flat capsules, no folders, no card faces, no pages, no markings'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-words.png',
    prefix: 'poliglot_ai/header-words',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202621,
    prompt: [basePositive, 'wide header for learning new words, right-side word-seed constellation of knowledge gems, memory capsules, route sparks, calm luxury study mood, no folders, no letters, no pages, no markings'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-word-game.png',
    prefix: 'poliglot_ai/header-word-game',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202622,
    prompt: [basePositive, 'wide header for a vocabulary game, right-side kinetic puzzle prism, rotating blank challenge pieces, reward ring, playful but polished, no text, no symbols, no screens, no folders'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-spelling.png',
    prefix: 'poliglot_ai/header-spelling',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202623,
    prompt: [basePositive, 'wide header for spelling practice, right-side precision alignment instrument, glowing check ring, soft correction arcs, high-end focus exercise mood, no letters, no keyboard, no documents, no folders'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-level.png',
    prefix: 'poliglot_ai/header-level',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205204616,
    prompt: [basePositive, 'wide header for language level assessment, right-side smooth level calibration observatory with stepped sapphire lens, clean floating mastery rings, mint alignment halo, precise modern education mood, no gauge scale, no tick marks, no measuring arc, no screens, no text, no letters, no numbers'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-leaderboard.png',
    prefix: 'poliglot_ai/header-leaderboard',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202624,
    prompt: [basePositive, 'wide header for leaderboard and rank competition, right-side luxury podium crystals, ascending achievement gems, glowing rank rings, competitive high-end learning mood, no numbers, no labels, no plaques'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-premium.png',
    prefix: 'poliglot_ai/header-premium',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202617,
    prompt: [basePositive, 'wide header for subscription upgrade, right-side neural crown core, glowing learning ring, faceted blue gem shield, luxury upgrade, gold blue teal white glass, no pass card, no documents, no inscriptions'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-limits.png',
    prefix: 'poliglot_ai/header-limits',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202625,
    prompt: [basePositive, 'wide header for daily learning limits, right-side refined capacity gauges, nested progress rings, measured crystal reservoirs, usage concept as abstract objects only, no screens, no text, no numbers'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-mistakes.png',
    prefix: 'poliglot_ai/header-mistakes',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202618,
    prompt: [basePositive, 'wide header for mistake review and correction, right-side sculptural check mark, gentle repair loop, neural feedback rings, polished correction gems, encouraging mood, no cards, no documents, no folders, no written marks'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-tools.png',
    prefix: 'poliglot_ai/header-tools',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202619,
    prompt: [basePositive, 'wide header for intelligent tools, right-side translation instrument cluster, sculptural voice waveform, camera-lens prism, smooth language capsules, small neural control nodes, clean futuristic SaaS, no screens, no cards, no text'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-referral.png',
    prefix: 'poliglot_ai/header-referral',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202620,
    prompt: [basePositive, 'wide header for referrals and social growth, connected learning network nodes, folded ribbon gift sculpture, reward spark, friendly high-end growth design, object floating freely in studio space, no hands, no fingers, no person, no invitation card, no documents, no text, no screens'].join(', '),
  },
  {
    group: 'headers',
    file: 'header-settings.png',
    prefix: 'poliglot_ai/header-settings',
    width: 1536,
    height: 640,
    baseWidth: 1152,
    baseHeight: 480,
    seed: 2205202626,
    prompt: [basePositive, 'wide header for settings and account controls, right-side elegant control dials, smooth toggle capsules, secure profile ring, quiet luxury utility mood, no interface panel, no text, no symbols'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-home.png',
    prefix: 'poliglot_ai/icon-home',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205205630,
    prompt: [basePositive, 'square app icon for main menu button, beautiful home cockpit mark, sapphire learning hub gem inside a soft mint navigation ring, tiny polished route nodes, centered premium object, clean blank background, no text, no marks, no folders'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-lesson.png',
    prefix: 'poliglot_ai/icon-lesson',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202631,
    prompt: [basePositive, 'square app icon for new lesson action, luminous lesson lens with a folded ceramic ribbon and one mint knowledge crystal, centered object, clean blank background, no book, no pages, no title, no markings'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-practice.png',
    prefix: 'poliglot_ai/icon-practice',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202632,
    prompt: [basePositive, 'square app icon for conversation practice action, two glossy speech capsules in dialogue, subtle voice waveform arcs, neural glow, centered object, clean background, no panels, no marks'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-shadowing.png',
    prefix: 'poliglot_ai/icon-shadowing',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205203231,
    prompt: [basePositive, 'square app icon for shadowing speech trainer, one large microphone capsule wrapped by open circular echo rings and a small sapphire pronunciation core, mint sound waves returning back, instantly reads as listen and repeat aloud, centered object, clean blank background, no rectangular mirror, no panel, no slab, no headphones, no screens, no letters, no marks'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-roleplay.png',
    prefix: 'poliglot_ai/icon-roleplay',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205203241,
    prompt: [basePositive, 'square app icon for AI roleplay language scenes, two polished theatre masks made of blank speech capsules orbiting a sapphire dialogue core, mint route arcs, warm gold scene spark, centered premium object, clean blank background, no people, no face, no letters, no text, no UI'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-pronunciation.png',
    prefix: 'poliglot_ai/icon-pronunciation',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205203242,
    prompt: [basePositive, 'square app icon for pronunciation analysis, sapphire sound core with precise waveform heatmap rings and a small calibration lens, mint acoustic arcs, reads as voice accuracy and speech clarity, centered object, clean blank background, no microphone duplicate, no letters, no text, no numbers'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-offline.png',
    prefix: 'poliglot_ai/icon-offline',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205203243,
    prompt: [basePositive, 'square app icon for offline mini decks, one open crescent-shaped sapphire learning capsule, detached mint pause ribbon, two separated pearl memory beads, portable saved-practice concept, centered premium object, very airy open silhouette with transparent gaps, clean blank background, no cube, no box, no container, no square silhouette, no square stack, no filled tile, no wifi symbol, no phone, no cards with text, no letters'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-phrasebook.png',
    prefix: 'poliglot_ai/icon-phrasebook',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205203244,
    prompt: [basePositive, 'square app icon for favorite phrases and personal phrasebook, elegant bookmark ribbon wrapped around a blank speech capsule with sapphire phrase gem, mint memory arc, centered premium object, clean blank background, no book pages, no printed marks, no letters, no text'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-words.png',
    prefix: 'poliglot_ai/icon-words',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205204635,
    prompt: [basePositive, 'square app icon for learning words, word-seed constellation instrument with one sapphire seed gem, three oval glass knowledge beads, mint learning route arc, small warm-gold orbit nodes, centered object, clean background, no blank square, no folders, no letters, no text, no page lines, no printed marks, no acronym'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-word-game.png',
    prefix: 'poliglot_ai/icon-word-game',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202636,
    prompt: [basePositive, 'square app icon for vocabulary game, kinetic puzzle prism and floating reward ring, polished playful object, centered, clean background, no folders, no symbols, no letters'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-spelling.png',
    prefix: 'poliglot_ai/icon-spelling',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202637,
    prompt: [basePositive, 'square app icon for spelling practice, precision check ring with aligned blank sapphire segments, refined focus object, centered, no letters, no keyboard, no text, no documents'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-vocabulary.png',
    prefix: 'poliglot_ai/icon-vocabulary',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205204638,
    prompt: [basePositive, 'square app icon for vocabulary library, lexicon oracle instrument with floating glass orbit rings, sapphire core, mint route halo, round unmarked glass beads, centered, no drawers, no tablets, no plaques, no flat capsules, no folders, no pages, no labels, no text, no symbols, no printed marks'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-level.png',
    prefix: 'poliglot_ai/icon-level',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205204639,
    prompt: [basePositive, 'square app icon for language level assessment, smooth calibration lens, stepped sapphire mastery ring, mint alignment halo, precision object, centered, no gauge scale, no tick marks, no measuring arc, no letters, no numbers, no text'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-progress.png',
    prefix: 'poliglot_ai/icon-progress',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202640,
    prompt: [basePositive, 'square app icon for learning progress, glowing progress rings and rising crystal fins, refined analytics object, centered, no numbers, no labels, no dashboard'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-awards.png',
    prefix: 'poliglot_ai/icon-awards',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202633,
    prompt: [basePositive, 'square app icon for achievements action, refined trophy gem with blank medal ring and milestone sparkles, elegant not cartoonish, centered object, no inscription, no plaque'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-leaderboard.png',
    prefix: 'poliglot_ai/icon-leaderboard',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202641,
    prompt: [basePositive, 'square app icon for leaderboard, elegant podium crystals with rank ring and achievement gem, centered object, no numbers, no plaques, no text'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-mistakes.png',
    prefix: 'poliglot_ai/icon-mistakes',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202642,
    prompt: [basePositive, 'square app icon for mistake correction, sculptural check mark and repair loop around polished correction gem, encouraging object, centered, no cards, no letters, no text'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-referral.png',
    prefix: 'poliglot_ai/icon-referral',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202643,
    prompt: [basePositive, 'square app icon for referrals, connected network nodes with folded ribbon gift sculpture, friendly high-end growth object, centered, no invitation card, no text, no hands'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-premium.png',
    prefix: 'poliglot_ai/icon-premium',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202634,
    prompt: [basePositive, 'square app icon for subscription upgrade action, gold and silver crown core inside glowing neural ring, faceted blue gem shield, luxurious modern object, no pass card, no document, no inscription'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-limits.png',
    prefix: 'poliglot_ai/icon-limits',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202644,
    prompt: [basePositive, 'square app icon for learning limits, capacity gauge rings and measured usage crystals, refined object, centered, no numbers, no labels, no text, no dashboard'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-tools.png',
    prefix: 'poliglot_ai/icon-tools',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202645,
    prompt: [basePositive, 'square app icon for intelligent tools, voice waveform sculpture crossing a camera-lens translation prism, clean futuristic object, centered, no screens, no text, no UI'].join(', '),
  },
  {
    group: 'icons',
    file: 'icon-settings.png',
    prefix: 'poliglot_ai/icon-settings',
    width: 512,
    height: 512,
    baseWidth: 640,
    baseHeight: 640,
    refineWidth: 768,
    refineHeight: 768,
    seed: 2205202646,
    prompt: [basePositive, 'square app icon for settings, elegant control dial, smooth toggle capsule, secure profile ring, centered, no interface panel, no text, no symbols'].join(', '),
  },
  {
    group: 'panels',
    file: 'panel-progress.png',
    prefix: 'poliglot_ai/panel-progress',
    width: 960,
    height: 720,
    baseWidth: 960,
    baseHeight: 720,
    seed: 2205206910,
    prompt: [
      basePositive,
      'right-side inspector card decorative background, progress and achievement mood without any interface',
      'achievement sculptural rings, sapphire level gem, rising mint energy arc, small milestone crystals, polished glass depth',
      'object detail only on right and bottom edges, upper-left is pure blank studio space with no marks',
      'no numbers, no text, no words, no typography, no dashboard, no screen, no labels, no UI rows, no paragraphs',
    ].join(', '),
  },
  {
    group: 'panels',
    file: 'panel-limits.png',
    prefix: 'poliglot_ai/panel-limits',
    width: 960,
    height: 720,
    baseWidth: 960,
    baseHeight: 720,
    seed: 2205206911,
    prompt: [
      basePositive,
      'right-side inspector card decorative background, daily capacity and usage mood without any interface',
      'measured capacity rings, mint crystal reservoirs, warm-gold usage arc, quiet utility sculpture, refined control instrument',
      'object detail only on right and bottom edges, upper-left is pure blank studio space with no marks',
      'no numbers, no text, no words, no typography, no gauge ticks, no dashboard, no screen, no labels, no UI rows, no paragraphs',
    ].join(', '),
  },
  {
    group: 'panels',
    file: 'panel-account.png',
    prefix: 'poliglot_ai/panel-account',
    width: 960,
    height: 720,
    baseWidth: 960,
    baseHeight: 720,
    seed: 2205206912,
    prompt: [
      basePositive,
      'decorative account sidebar backdrop only, secure account and premium access mood without any interface screenshot',
      'secure pearl capsule, sapphire shield gem, premium access core, warm-gold protection ring, mint privacy halo, soft glass depth, abstract product object only',
      'one cohesive object cluster anchored to the right and bottom edges, detail can visually bleed beyond the panel frame after CSS masking, upper-left is pure blank studio space with no marks',
      'no person, no face, no ID card, no account card, no profile row, no menu row, no text, no words, no typography, no dashboard, no screen, no labels, no UI rows, no paragraphs',
    ].join(', '),
  },
  ...Array.from({ length: 20 }, (_, index) => {
    const level = index + 1;
    return {
      group: 'awards',
      file: `award-${String(level).padStart(2, '0')}.png`,
      prefix: `poliglot_ai/award-${String(level).padStart(2, '0')}`,
      width: 768,
      height: 768,
      baseWidth: 768,
      baseHeight: 768,
      seed: 2205206000 + level,
      prompt: [
        basePositive,
        `square trophy artwork for language XP level ${level}, unique premium trophy cup or medal sculpture, no visible number`,
        level <= 5
          ? 'starter bronze-sapphire learning trophy, simple polished ring, mint route spark, friendly entry rank'
          : level <= 10
            ? 'silver-sapphire polyglot trophy, layered speech capsule crown, refined progress halo'
            : level <= 15
              ? 'gold-sapphire mastery trophy, faceted language constellation ring, luminous achievement core'
              : 'platinum-sapphire grand polyglot trophy, celestial translation halo, museum-grade mastery relic',
        'centered trophy cup object with handles or medal ring silhouette, clean blank studio background, generous padding, no flat square tile, no panel, no plaque, no engraved text, no numbers, no letters, no flags, no people',
      ].join(', '),
    };
  }),
  {
    group: 'logos',
    file: 'brand-logo.png',
    prefix: 'poliglot_ai/brand-logo',
    width: 768,
    height: 768,
    baseWidth: 768,
    baseHeight: 768,
    seed: 2205207201,
    prompt: [
      basePositive,
      'transparent NERIVA logo mark source, one freestanding iconic language-learning symbol sculpture, no typography',
      'sapphire speech capsule core inside a smooth mint translation route ring, tiny warm-gold achievement star floating outside the ring, polished ceramic and glass, memorable silhouette that reads as listen, speak, translate, progress',
      'object fills almost the whole square, transparent-ready chroma background only, no white background, no square tile, no badge, no frame, no letters, no words, no triangle, no arrow, no tiny printed mark',
    ].join(', '),
  },
  {
    group: 'backgrounds',
    file: 'brand-logo-hero-bg.png',
    prefix: 'poliglot_ai/brand-logo-hero-bg',
    width: 1920,
    height: 1080,
    baseWidth: 1344,
    baseHeight: 768,
    seed: 2205207302,
    prompt: [
      basePositive,
      'large thematic background image for a language-learning app, designed to receive a separate real transparent logo composited later in the exact center',
      'premium language-learning world: blank speech capsules, smooth continuous mint waveform ribbons, sapphire learning cores, translation route arcs, soft achievement crystal glints around the outer edges',
      'center is completely empty clean luminous negative space with no object, no emblem, no label and no mark, rich edge detail, works as desktop and mobile cover background',
      'do not write Poliglot, AI, app names, words, letters or captions anywhere; no text, no letters, no logos, no UI, no chat bubbles, no cards, no panels, no people, no flags, no isolated glyph shapes, no arrowheads, no letter-like mint marks',
    ].join(', '),
  },
  {
    group: 'backgrounds',
    file: 'auth-login-hero.png',
    prefix: 'poliglot_ai/auth-login-hero',
    width: 1920,
    height: 1080,
    baseWidth: 1344,
    baseHeight: 768,
    seed: 2705202607,
    prompt: [
      basePositive,
      'login page hero background for an AI language-learning service, designed to sit beside a real sign-in form on desktop and crop cleanly on mobile',
      'premium language-learning gateway with blank frosted speech capsules, sapphire account core with no symbol inside, mint learning route arcs, subtle privacy shield halo, daily route orbit and small warm-gold achievement glint',
      'left and center-left remain calm and readable for real interface content, richer object detail on the right and lower-right, same visual language as the generated V2 menu assets',
      'all rings and surfaces are completely blank: no login form, no fake UI, no phone, no dashboard, no text, no labels, no words, no letters, no brand names, no logo, no watermark, no Telegram logo, no paper-plane symbol, no icons, no people, no flags, no glyph-like marks',
    ].join(', '),
  },
  {
    group: 'backgrounds',
    file: 'app-background.png',
    prefix: 'poliglot_ai/app-background',
    width: 1920,
    height: 1080,
    baseWidth: 1344,
    baseHeight: 768,
    seed: 2205205901,
    prompt: [
      basePositive,
      'full-screen decorative backdrop for NERIVA language-learning web app, designed to sit behind the real interface on desktop and mobile',
      'immersive premium language intelligence atmosphere made of blank translucent sound capsules, soft voice waveform rivers, sapphire learning core near the outer edges, mint translation route arcs, small knowledge gems and calm depth',
      'must work behind glass UI panels, with readable calm negative space through the center and center-left, richer object detail near the outer edges, and a safe mobile portrait crop',
      'one cohesive full-bleed backdrop, no repeated pattern, no tiling, no chat message bubbles, no UI cards, no pseudo text strokes, no screen, no phone, no dashboard, no text, no letters, no words, no logos, no people, no flags',
    ].join(', '),
  },
  {
    group: 'site',
    file: 'site-hero-cockpit.png',
    prefix: 'poliglot_ai/site-hero-cockpit',
    width: 1920,
    height: 1080,
    baseWidth: 1344,
    baseHeight: 768,
    refineWidth: 1536,
    refineHeight: 864,
    seed: 2205202701,
    prompt: [
      basePositive,
      'landing page hero background for a premium AI language coach',
      'large language intelligence cockpit sculpture made from three blank speech bubble solids, neural route arcs, voice waveform ribbon, memory tiles, progress ring, Telegram and web sync feeling without logos',
      'must clearly suggest speaking, listening, translating and daily progress through abstract objects only',
      'composition leaves clean negative space on the left for real HTML headline text',
      'deep indigo studio backdrop, teal neural glow, mint highlights, small gold learning achievement detail',
      'no screen, no phone, no dashboard, no text, no letters, no flags, no people',
    ].join(', '),
  },
  {
    group: 'site',
    file: 'site-languages-map.png',
    prefix: 'poliglot_ai/site-languages-map',
    width: 1600,
    height: 900,
    baseWidth: 1280,
    baseHeight: 720,
    seed: 2205202702,
    prompt: [
      basePositive,
      'public site section image for 20 languages and global learning routes',
      'abstract world language route map made of glowing route nodes, curved translation paths, blank speech-bubble destination gems, no country flags and no written labels',
      'must feel international and multilingual without showing maps with text',
      'international, trustworthy, spacious, premium education mood',
      'deep indigo and white with teal routes, mint nodes, selective gold destination glow',
      'no maps with labels, no text, no numbers, no flags, no people',
    ].join(', '),
  },
  {
    group: 'site',
    file: 'site-feature-tools.png',
    prefix: 'poliglot_ai/site-feature-tools',
    width: 1600,
    height: 900,
    baseWidth: 1280,
    baseHeight: 720,
    seed: 2205202703,
    prompt: [
      basePositive,
      'public site feature image for AI tools: translation, voice, photo, text understanding',
      'sculptural voice waveform, camera lens translation prism, blank speech bubble pair, smooth language tokens, neural control ring, premium tool cluster',
      'clean product shot with right side object cluster and negative space',
      'no screen, no phone, no UI panel, no text, no letters, no people, no hands',
    ].join(', '),
  },
  {
    group: 'site',
    file: 'site-progress-awards.png',
    prefix: 'poliglot_ai/site-progress-awards',
    width: 1600,
    height: 900,
    baseWidth: 1280,
    baseHeight: 720,
    seed: 2205202704,
    prompt: [
      basePositive,
      'public site image for progress, streaks, awards, motivation',
      'premium trophy sculpture, progress rings, milestone path, small achievement gems, polished gold and blue accents',
      'motivating but not childish, high-end education product feel',
      'no plaques, no numbers, no text, no certificates, no people',
    ].join(', '),
  },
  {
    group: 'site',
    file: 'site-pricing-premium.png',
    prefix: 'poliglot_ai/site-pricing-premium',
    width: 1600,
    height: 900,
    baseWidth: 1280,
    baseHeight: 720,
    seed: 2205202705,
    prompt: [
      basePositive,
      'public site pricing image for three subscription tiers',
      'three abstract plan pedestals, crown and progress ring for premium, calm free study gem, platinum achievement sculpture, blank speech-bubble gems around the ladder',
      'clear hierarchy without any labels or numbers, refined price section visual',
      'deep indigo, teal, mint, selective gold, white glass',
      'no text, no letters, no price tags, no cards with writing, no screens',
    ].join(', '),
  },
  {
    group: 'site',
    file: 'site-legal-shield.png',
    prefix: 'poliglot_ai/site-legal-shield',
    width: 1600,
    height: 900,
    baseWidth: 1280,
    baseHeight: 720,
    seed: 2205202706,
    prompt: [
      basePositive,
      'public site legal page image for trust, privacy, terms, account protection',
      'polished shield sculpture, secure neural ring, calm data vault shapes, soft blue and teal legal trust mood',
      'editorial legal hero background with negative space',
      'no documents, no forms, no locks with text, no letters, no people, no hands',
    ].join(', '),
  },
];

const assets = baseAssets.flatMap((asset) => {
  if (!THEMED_GROUPS.has(asset.group)) return [asset];
  return Object.entries(THEMES).map(([theme, themePrompt]) => ({
    ...asset,
    file: themedFile(asset.file, theme),
    prefix: themedPrefix(asset.prefix, theme),
    theme,
    prompt: [
      asset.prompt,
      themePrompt,
      'must match the paired light and dark language-learning asset composition, surfaces stay blank, no pseudo text, no letters, no logo marks',
    ].join(', '),
  }));
});

function workflow(asset) {
  if (isFlux2Model()) return workflowFlux2(asset);
  if (isFluxModel()) return workflowFlux(asset);
  return isSD3Model() ? workflowSD3(asset) : workflowSDXL(asset);
}

function workflowFlux2(asset) {
  const baseWidth = asset.baseWidth || asset.width;
  const baseHeight = asset.baseHeight || asset.height;
  const flux2Prompt = effectivePrompt(asset);
  return {
    '1': {
      class_type: 'UNETLoader',
      inputs: {
        unet_name: FLUX2_DIFFUSION_MODEL,
        weight_dtype: 'default',
      },
    },
    '2': {
      class_type: 'CLIPLoader',
      inputs: {
        clip_name: FLUX2_TEXT_ENCODER,
        type: 'flux2',
        device: 'default',
      },
    },
    '3': {
      class_type: 'VAELoader',
      inputs: { vae_name: FLUX2_VAE },
    },
    '4': {
      class_type: 'CLIPTextEncode',
      inputs: {
        clip: ['2', 0],
        text: flux2Prompt,
      },
    },
    '5': {
      class_type: 'BasicGuider',
      inputs: {
        model: ['1', 0],
        conditioning: ['4', 0],
      },
    },
    '6': {
      class_type: 'RandomNoise',
      inputs: { noise_seed: asset.seed },
    },
    '7': {
      class_type: 'KSamplerSelect',
      inputs: { sampler_name: 'euler' },
    },
    '8': {
      class_type: 'Flux2Scheduler',
      inputs: {
        steps: asset.group === 'icons' ? 16 : 20,
        width: baseWidth,
        height: baseHeight,
      },
    },
    '9': {
      class_type: 'EmptyFlux2LatentImage',
      inputs: {
        width: baseWidth,
        height: baseHeight,
        batch_size: 1,
      },
    },
    '10': {
      class_type: 'SamplerCustomAdvanced',
      inputs: {
        noise: ['6', 0],
        guider: ['5', 0],
        sampler: ['7', 0],
        sigmas: ['8', 0],
        latent_image: ['9', 0],
      },
    },
    '11': {
      class_type: 'VAEDecode',
      inputs: { samples: ['10', 0], vae: ['3', 0] },
    },
    '12': {
      class_type: 'ImageScale',
      inputs: {
        image: ['11', 0],
        upscale_method: 'lanczos',
        width: asset.width,
        height: asset.height,
        crop: 'center',
      },
    },
    '14': {
      class_type: 'SaveImage',
      inputs: { images: ['12', 0], filename_prefix: asset.prefix },
    },
  };
}

function workflowFlux(asset) {
  const baseWidth = asset.baseWidth || asset.width;
  const baseHeight = asset.baseHeight || asset.height;
  const fluxPrompt = [
    effectivePrompt(asset),
    'strict exclusions: no text, no letters, no numbers, no logo, no watermark, no folders, no screens, no documents, no UI panels, no account card, no chat message, no paragraph-like strokes, no pseudo typography, no repeated generic icons',
  ].join(', ');
  return {
    '1': {
      class_type: 'CheckpointLoaderSimple',
      inputs: { ckpt_name: MODEL },
    },
    '2': {
      class_type: 'ModelSamplingFlux',
      inputs: {
        model: ['1', 0],
        max_shift: asset.group === 'icons' ? 1.05 : 1.15,
        base_shift: asset.group === 'icons' ? 0.45 : 0.50,
        width: baseWidth,
        height: baseHeight,
      },
    },
    '3': {
      class_type: 'CLIPTextEncodeFlux',
      inputs: {
        clip: ['1', 1],
        clip_l: fluxPrompt,
        t5xxl: fluxPrompt,
        guidance: asset.group === 'icons' ? 3.8 : 3.5,
      },
    },
    '4': {
      class_type: 'CLIPTextEncodeFlux',
      inputs: {
        clip: ['1', 1],
        clip_l: '',
        t5xxl: '',
        guidance: 1,
      },
    },
    '5': {
      class_type: 'EmptyLatentImage',
      inputs: { width: baseWidth, height: baseHeight, batch_size: 1 },
    },
    '6': {
      class_type: 'KSampler',
      inputs: {
        model: ['2', 0],
        positive: ['3', 0],
        negative: ['4', 0],
        latent_image: ['5', 0],
        seed: asset.seed,
        steps: asset.group === 'icons' ? 22 : 28,
        cfg: 1,
        sampler_name: 'euler',
        scheduler: 'simple',
        denoise: 1,
      },
    },
    '7': {
      class_type: 'VAEDecode',
      inputs: { samples: ['6', 0], vae: ['1', 2] },
    },
    '8': {
      class_type: 'UpscaleModelLoader',
      inputs: { model_name: UPSCALE_MODEL },
    },
    '9': {
      class_type: 'ImageUpscaleWithModel',
      inputs: { upscale_model: ['8', 0], image: ['7', 0] },
    },
    '10': {
      class_type: 'ImageScale',
      inputs: {
        image: ['9', 0],
        upscale_method: 'lanczos',
        width: asset.width,
        height: asset.height,
        crop: 'center',
      },
    },
    '14': {
      class_type: 'SaveImage',
      inputs: { images: ['10', 0], filename_prefix: asset.prefix },
    },
  };
}

function workflowSD3(asset) {
  const baseWidth = asset.baseWidth || asset.width;
  const baseHeight = asset.baseHeight || asset.height;
  const prompt = effectivePrompt(asset);
  return {
    '1': {
      class_type: 'CheckpointLoaderSimple',
      inputs: { ckpt_name: MODEL },
    },
    '2': {
      class_type: 'ModelSamplingSD3',
      inputs: { model: ['1', 0], shift: asset.group === 'icons' ? 2.6 : 3.0 },
    },
    '3': {
      class_type: 'SkipLayerGuidanceSD3',
      inputs: {
        model: ['2', 0],
        layers: '7, 8, 9',
        scale: asset.group === 'icons' ? 2.2 : 2.8,
        start_percent: 0.01,
        end_percent: 0.16,
      },
    },
    '4': {
      class_type: 'CLIPTextEncodeSD3',
      inputs: {
        clip: ['1', 1],
        clip_l: prompt,
        clip_g: prompt,
        t5xxl: prompt,
        empty_padding: 'empty_prompt',
      },
    },
    '5': {
      class_type: 'CLIPTextEncodeSD3',
      inputs: {
        clip: ['1', 1],
        clip_l: negative,
        clip_g: negative,
        t5xxl: negative,
        empty_padding: 'empty_prompt',
      },
    },
    '6': {
      class_type: 'EmptySD3LatentImage',
      inputs: { width: baseWidth, height: baseHeight, batch_size: 1 },
    },
    '7': {
      class_type: 'KSampler',
      inputs: {
        model: ['3', 0],
        positive: ['4', 0],
        negative: ['5', 0],
        latent_image: ['6', 0],
        seed: asset.seed,
        steps: asset.group === 'icons' ? 24 : 30,
        cfg: asset.group === 'icons' ? 4.0 : 4.6,
        sampler_name: 'dpmpp_2m',
        scheduler: 'sgm_uniform',
        denoise: 1,
      },
    },
    '8': {
      class_type: 'VAEDecode',
      inputs: { samples: ['7', 0], vae: ['1', 2] },
    },
    '9': {
      class_type: 'UpscaleModelLoader',
      inputs: { model_name: UPSCALE_MODEL },
    },
    '10': {
      class_type: 'ImageUpscaleWithModel',
      inputs: { upscale_model: ['9', 0], image: ['8', 0] },
    },
    '11': {
      class_type: 'ImageScale',
      inputs: {
        image: ['10', 0],
        upscale_method: 'lanczos',
        width: asset.width,
        height: asset.height,
        crop: 'center',
      },
    },
    '14': {
      class_type: 'SaveImage',
      inputs: { images: ['11', 0], filename_prefix: asset.prefix },
    },
  };
}

function workflowSDXL(asset) {
  const baseWidth = asset.baseWidth || asset.width;
  const baseHeight = asset.baseHeight || asset.height;
  const refineWidth = asset.refineWidth || asset.width;
  const refineHeight = asset.refineHeight || asset.height;
  const prompt = effectivePrompt(asset);
  return {
    '1': {
      class_type: 'CheckpointLoaderSimple',
      inputs: { ckpt_name: MODEL },
    },
    '2': {
      class_type: 'CLIPTextEncode',
      inputs: { clip: ['1', 1], text: prompt },
    },
    '3': {
      class_type: 'CLIPTextEncode',
      inputs: { clip: ['1', 1], text: negative },
    },
    '4': {
      class_type: 'EmptyLatentImage',
      inputs: { width: baseWidth, height: baseHeight, batch_size: 1 },
    },
    '5': {
      class_type: 'KSampler',
      inputs: {
        model: ['1', 0],
        positive: ['2', 0],
        negative: ['3', 0],
        latent_image: ['4', 0],
        seed: asset.seed,
        steps: asset.group === 'icons' ? 22 : 28,
        cfg: 5.8,
        sampler_name: 'dpmpp_3m_sde_gpu',
        scheduler: 'exponential',
        denoise: 1,
      },
    },
    '6': {
      class_type: 'VAEDecode',
      inputs: { samples: ['5', 0], vae: ['1', 2] },
    },
    '7': {
      class_type: 'ImageScale',
      inputs: {
        image: ['6', 0],
        upscale_method: 'lanczos',
        width: refineWidth,
        height: refineHeight,
        crop: 'center',
      },
    },
    '8': {
      class_type: 'VAEEncode',
      inputs: { pixels: ['7', 0], vae: ['1', 2] },
    },
    '9': {
      class_type: 'KSampler',
      inputs: {
        model: ['1', 0],
        positive: ['2', 0],
        negative: ['3', 0],
        latent_image: ['8', 0],
        seed: asset.seed + 10000,
        steps: asset.group === 'icons' ? 12 : 14,
        cfg: 5.2,
        sampler_name: 'dpmpp_2m_sde_gpu',
        scheduler: 'karras',
        denoise: asset.group === 'icons' ? 0.2 : 0.24,
      },
    },
    '10': {
      class_type: 'VAEDecode',
      inputs: { samples: ['9', 0], vae: ['1', 2] },
    },
    '11': {
      class_type: 'UpscaleModelLoader',
      inputs: { model_name: UPSCALE_MODEL },
    },
    '12': {
      class_type: 'ImageUpscaleWithModel',
      inputs: { upscale_model: ['11', 0], image: ['10', 0] },
    },
    '13': {
      class_type: 'ImageScale',
      inputs: {
        image: ['12', 0],
        upscale_method: 'lanczos',
        width: asset.width,
        height: asset.height,
        crop: 'center',
      },
    },
    '14': {
      class_type: 'SaveImage',
      inputs: { images: ['13', 0], filename_prefix: asset.prefix },
    },
  };
}

async function queue(asset) {
  const response = await fetch(`${COMFY_BASE}/prompt`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      client_id: randomUUID(),
      prompt: workflow(asset),
    }),
  });
  const data = await response.json();
  if (!response.ok || data.error) {
    throw new Error(`ComfyUI rejected ${asset.file}: ${JSON.stringify(data, null, 2)}`);
  }
  return data.prompt_id;
}

async function history(promptId) {
  const response = await fetch(`${COMFY_BASE}/history/${promptId}`);
  if (!response.ok) return null;
  const data = await response.json();
  return data[promptId] || null;
}

async function waitForImage(promptId, asset) {
  const started = Date.now();
  while (Date.now() - started < 20 * 60 * 1000) {
    const item = await history(promptId);
    const images = item?.outputs?.['14']?.images;
    if (images?.length) return images[0];
    await new Promise((resolve) => setTimeout(resolve, 2500));
  }
  throw new Error(`Timed out waiting for ${asset.file}`);
}

async function saveComfyImage(image, asset) {
  const params = new URLSearchParams({
    filename: image.filename,
    subfolder: image.subfolder || '',
    type: image.type || 'output',
  });
  const response = await fetch(`${COMFY_BASE}/view?${params.toString()}`);
  if (!response.ok || !response.body) {
    throw new Error(`Could not download ${asset.file}: HTTP ${response.status}`);
  }
  const out = TRANSPARENT_GROUPS.has(asset.group)
    ? join(RAW_TRANSPARENT_DIR, asset.file)
    : join(OUTPUT_DIR, asset.file);
  await mkdir(dirname(out), { recursive: true });
  await pipeline(response.body, createWriteStream(out));
  return out;
}

async function postProcessTransparentAsset(rawPath, asset) {
  const out = join(OUTPUT_DIR, asset.file);
  const report = join(REPORT_DIR, asset.file.replace(/\.png$/i, '.json'));
  const kind = asset.group === 'awards' ? 'award' : 'icon';
  const { stdout, stderr } = await execFileAsync(PYTHON, [
    'tools/postprocess_transparent_asset.py',
    '--input', rawPath,
    '--out', out,
    '--kind', kind,
    '--report', report,
  ], { maxBuffer: 1024 * 1024 * 8 });
  if (stderr) process.stderr.write(stderr);
  if (stdout) process.stdout.write(stdout);
  return JSON.parse(await readFile(report, 'utf8'));
}

async function ensureServer() {
  const response = await fetch(`${COMFY_BASE}/system_stats`);
  if (!response.ok) throw new Error(`ComfyUI is not ready at ${COMFY_BASE}`);
}

async function readExistingManifest() {
  try {
    const text = await readFile(join(OUTPUT_DIR, 'brand-assets-manifest.json'), 'utf8');
    const parsed = JSON.parse(text);
    if (Array.isArray(parsed)) return parsed;
  } catch {
    // A missing or old manifest should not block asset generation.
  }
  return [];
}

async function main() {
  await ensureServer();
  const arg = process.argv[2] || 'all';
  const only = new Set((process.env.ONLY_ASSETS || '')
    .split(',')
    .map((value) => value.trim())
    .filter(Boolean));
  const selected = assets.filter((asset) => {
    if (only.size) return only.has(asset.file.replace(/\.png$/, '')) || only.has(asset.file);
    return arg === 'all' || asset.group === arg || asset.file.replace(/\.png$/, '') === arg;
  });
  if (!selected.length) {
    throw new Error(`No assets selected for "${arg}"`);
  }

  const manifestByFile = new Map((await readExistingManifest())
    .filter((entry) => entry?.file)
    .map((entry) => [entry.file, entry]));
  for (const asset of selected) {
    const attempts = TRANSPARENT_GROUPS.has(asset.group)
      ? Math.max(1, Number(process.env.TRANSPARENT_ASSET_ATTEMPTS || 3))
      : 1;
    let lastError = null;
    let completed = null;
    for (let attempt = 0; attempt < attempts; attempt += 1) {
      const workingAsset = { ...asset, seed: asset.seed + attempt * 100000 };
      try {
        console.log(`queue ${workingAsset.file} ${workingAsset.width}x${workingAsset.height} attempt ${attempt + 1}/${attempts}`);
        const promptId = await queue(workingAsset);
        const image = await waitForImage(promptId, workingAsset);
        const saved = await saveComfyImage(image, workingAsset);
        const validation = TRANSPARENT_GROUPS.has(workingAsset.group)
          ? await postProcessTransparentAsset(saved, workingAsset)
          : null;
        completed = { workingAsset, promptId, saved, validation };
        break;
      } catch (error) {
        lastError = error;
        console.error(`failed ${asset.file} attempt ${attempt + 1}: ${error.message}`);
      }
    }
    if (!completed) {
      throw lastError || new Error(`Could not generate ${asset.file}`);
    }

    const { workingAsset, promptId, saved, validation } = completed;
    manifestByFile.set(asset.file, {
      file: asset.file,
      group: asset.group,
      theme: asset.theme || 'fixed',
      model: MODEL,
      workflow: isFlux2Model() ? 'flux2_split_comfyui' : (isFluxModel() ? 'flux_checkpoint_comfyui' : (isSD3Model() ? 'sd3_comfyui' : 'sdxl_comfyui')),
      flux2_diffusion_model: isFlux2Model() ? FLUX2_DIFFUSION_MODEL : undefined,
      flux2_text_encoder: isFlux2Model() ? FLUX2_TEXT_ENCODER : undefined,
      flux2_vae: isFlux2Model() ? FLUX2_VAE : undefined,
      fallback_models: FALLBACK_MODELS,
      prompt_id: promptId,
      seed: workingAsset.seed,
      width: workingAsset.width,
      height: workingAsset.height,
      prompt: effectivePrompt(workingAsset),
      transparent_postprocess: validation || undefined,
      qa_status: validation?.ok ? 'generated_transparent_validated' : 'generated_visual_review_required',
    });
    console.log(`saved ${TRANSPARENT_GROUPS.has(workingAsset.group) ? join(OUTPUT_DIR, workingAsset.file) : saved}`);
  }

  const manifest = assets
    .map((asset) => manifestByFile.get(asset.file) || ({
      file: asset.file,
      group: asset.group,
      theme: asset.theme || 'fixed',
      model: MODEL,
      workflow: isFlux2Model() ? 'flux2_split_comfyui' : (isFluxModel() ? 'flux_checkpoint_comfyui' : (isSD3Model() ? 'sd3_comfyui' : 'sdxl_comfyui')),
      flux2_diffusion_model: isFlux2Model() ? FLUX2_DIFFUSION_MODEL : undefined,
      flux2_text_encoder: isFlux2Model() ? FLUX2_TEXT_ENCODER : undefined,
      flux2_vae: isFlux2Model() ? FLUX2_VAE : undefined,
      fallback_models: FALLBACK_MODELS,
      seed: asset.seed,
      width: asset.width,
      height: asset.height,
      prompt: effectivePrompt(asset),
      qa_status: 'pending_generation',
    }));
  await writeFile(join(OUTPUT_DIR, 'brand-assets-manifest.json'), JSON.stringify(manifest, null, 2), 'utf8');
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
