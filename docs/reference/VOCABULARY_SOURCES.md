# Vocabulary Sources

The bot ships local JSON learning dictionaries for every selectable learning language:

- `vocabulary_words.json` for English: 105,197 entries.
- `vocabulary_words_ru.json` for Russian: 54,831 entries.
- `vocabulary_words_es.json` for Spanish: 18,890 entries.
- `vocabulary_words_de.json` for German: 65,811 entries.
- `vocabulary_words_fr.json` for French: 65,673 entries.
- `vocabulary_words_it.json` for Italian: 25,799 entries.
- `vocabulary_words_zh.json` for Chinese: 8,781 entries.
- `vocabulary_words_ja.json` for Japanese: 12,094 entries.
- `vocabulary_words_ko.json` for Korean: 26,801 entries.
- `vocabulary_words_tg.json` for Tajik: 2,156 entries.
- `vocabulary_words_uz.json` for Uzbek: 2,765 entries.
- `vocabulary_words_tt.json` for Tatar: 1,516 entries.
- `vocabulary_words_hy.json` for Armenian: 13,592 entries.
- `vocabulary_words_kk.json` for Kazakh: 9,093 entries.
- `vocabulary_words_ky.json` for Kyrgyz: 2,139 entries.
- `vocabulary_words_ka.json` for Georgian: 16,839 entries.
- `vocabulary_words_uk.json` for Ukrainian: 23,318 entries.
- `vocabulary_words_pl.json` for Polish: 31,232 entries.
- `vocabulary_words_ro.json` for Romanian: 30,000 entries.
- `vocabulary_words_pt.json` for Portuguese: 11,789 entries.
- `vocabulary_words_ar.json` for Arabic: 6,203 entries.
- `vocabulary_words_bn.json` for Bengali: 6,919 entries.
- `vocabulary_words_cs.json` for Czech: 30,000 entries.
- `vocabulary_words_el.json` for Greek: 24,143 entries.
- `vocabulary_words_hi.json` for Hindi: 15,403 entries.
- `vocabulary_words_hu.json` for Hungarian: 30,000 entries.
- `vocabulary_words_id.json` for Indonesian: 20,493 entries.
- `vocabulary_words_nl.json` for Dutch: 29,999 entries.
- `vocabulary_words_sv.json` for Swedish: 30,000 entries.
- `vocabulary_words_ta.json` for Tamil: 8,308 entries.
- `vocabulary_words_te.json` for Telugu: 11,495 entries.
- `vocabulary_words_th.json` for Thai: 13,379 entries.
- `vocabulary_words_tl.json` for Tagalog: 19,469 entries.
- `vocabulary_words_tr.json` for Turkish: 16,778 entries.
- `vocabulary_words_vi.json` for Vietnamese: 9,765 entries.

Current open sources:

- English-Russian base: FreeDict/WikDict `eng-rus` TEI from <https://www.wikdict.com/page/download>, CC BY-SA 3.0.
- English learning dictionary: rebuilt from the English-Russian base and expanded through the multilingual open-dictionary pivots already downloaded for the other learning languages.
- Russian learning dictionary: rebuilt from the English-Russian base plus the Russian Kaikki/Wiktextract dump, so Russian entries keep the Wiktionary context instead of relying only on a trimmed bilingual pivot.
- Spanish, German, French, Italian, Polish, Portuguese, Chinese, and Japanese: WikDict/FreeDict TEI target-English and English-target dictionaries from <https://www.wikdict.com/page/download>, CC BY-SA.
- Korean, Tajik, Uzbek, Tatar, Armenian, Kazakh, Kyrgyz, Georgian, Ukrainian, Romanian, Russian, Arabic, Bengali, Czech, Greek, Hindi, Hungarian, Indonesian, Dutch, Swedish, Tamil, Telugu, Thai, Tagalog, Turkish, and Vietnamese: Kaikki/Wiktextract JSONL dumps from <https://kaikki.org/dictionary/>, derived from Wiktionary and distributed under the same licenses as Wiktionary.
- Learning order and first-screen quality: a local multilingual A1 core is placed before the imported dictionary tail, using open dictionary pivots plus curated common-word overrides for ambiguous high-frequency entries. The imported tail is ranked into A1-C2 with `wordfreq` Zipf frequency scores where available and an English-pivot fallback where a target-language frequency list is not available.
- Study-card display: English pivots are treated as internal build data. The bot and web app choose a prompt in the user's interface language when possible, allow multi-value dictionary prompts such as `ser; estar`, and suppress raw pivot context when it would leak English or reveal the answer.
- Original source files can be downloaded into `tools/sources/` (`*.tei` and `kaikki-*.jsonl`) when a full dictionary rebuild is needed. This directory is rebuild-only cache, is ignored by git, and is not required at runtime.

Cross-language behavior:

- Each `vocabWord` can contain a `translations` map keyed by interface language code.
- The importer builds these maps through an English pivot plus POS-aware direct pairs, so a Polish interface user can study German, a Portuguese interface user can study Ukrainian, and so on.
- If the selected interface language is the same as the learning language, the bot avoids showing the answer as the prompt and falls back to another available translation.
- If a specific word/interface-language prompt is missing, the runtime can ask OpenRouter for a precise dictionary translation. Successful AI translations are persisted in SQLite in `vocabulary_ai_translations` and mirrored into `vocabulary_translations`, so the next request uses local data instead of another AI call.
- If the local seed dictionary cannot provide a usable unlearned word/prompt for a word lesson, the runtime asks OpenRouter for one CEFR-level vocabulary card, validates that it is not a duplicate and does not reveal the answer, persists it in `vocabulary_ai_words`, and mirrors it into `vocabulary_words` plus `vocabulary_translations`. Duplicate checks use the deterministic vocabulary ID, so 30k-word language dictionaries remain fast in SQLite.
- The AI prompt is intentionally dictionary-style: strict JSON, common sense of the target lemma, part-of-speech/CEFR/topic aware, up to three short native-language equivalents separated by `; `, and no target-language answer, transliteration, examples, grammar notes, or Markdown.
- All 35 interface languages are now available as dictionary-backed learning languages. Before exposing another language, download source files into `tools/sources/`, document the license here, rebuild with the multilingual importer, and add regression coverage for word selection/options.

Generators:

- `scripts/generate_vocabulary_from_freedict.ps1` rebuilds the English base from `freedict-eng-rus.tei`.
- `tools/build_multilingual_vocab.py` rebuilds every multilingual JSON dictionary from `tools/sources/*` and preserves multi-gloss dictionary context in the generated `context` fields. Its rebuild-time Python dependencies are pinned in `tools/requirements-vocab.txt`.
- `tools/build_freedict_vocab.py` is kept as the old Polish/Portuguese-only helper; the multilingual importer supersedes it for production dictionaries.

The application does not copy Cambridge, Oxford, or other proprietary dictionary entries.
