# Full 35-Locale Localization Design

**Goal:** replace fallback-generated and mixed-language user copy with explicit, literary localization for all 35 interface languages across the web app, public site, legal pages, and Telegram bot.

**Approved scope:** everything user-facing: `web-react`, `site-react`, generated static public output, legal documents, Telegram bot menus/messages/payments/referrals/tools/onboarding, and localized helper names used in bot/web text.

**Locale set:** `ru`, `en`, `es`, `de`, `fr`, `it`, `zh`, `ja`, `ko`, `tg`, `uz`, `tt`, `hy`, `kk`, `ky`, `ka`, `uk`, `pl`, `ro`, `pt`, `ar`, `bn`, `cs`, `el`, `hi`, `hu`, `id`, `nl`, `sv`, `ta`, `te`, `th`, `tl`, `tr`, `vi`.

---

## Problem

The application exposes 35 interface languages, but several localization layers currently use generated or derived copy instead of complete human-quality translations.

In `web-react/src/lib/i18n.ts`, `buildDerivedGuideCopy`, `localizedFallbackForEnglishCopy`, `fallbackLeakAllowedKeys`, `isLikelyFallbackLeak`, and `normalizePremiumAndTutorCopy` can synthesize user-visible text from English, Russian, or short UI terms. That prevents Russian/English placeholder leaks in many cases, but it is not a full literary translation.

In the public site, `site-react/public/assets/site-i18n.js`, `site-phrases.js`, `legal-documents-i18n.js`, and `privacy-policy-i18n.js` overlap. `privacy-policy-i18n.js` still has a 20-locale legacy renderer, while `legal-documents-i18n.js` has 35-locale coverage. `privacy.html` loads both paths, so the privacy page can drift or fall back unexpectedly.

In the Go bot, multiple functions start from English or compact generated fragments: `ui`, `toolUICopyFor`, `systemUI`, `premiumUI`, `premiumLandingTerms`, `vocabularyFallbackPrompt`, and `localizedLanguageNameForInterface`. Current tests catch some fallback leaks, but they still allow or encode English/Russian fallback behavior in several areas.

## Design Options

### Option A: Keep generators, improve heuristics

This is the smallest code change. We would expand `localizedFallbackForEnglishCopy`, compact Go fallbacks, and public-site phrase generators.

Rejected because it keeps the current problem: the interface remains assembled from fragments, not translated as natural product copy.

### Option B: External translation files per surface

Move web app, public site, legal, and bot copy into standalone JSON or TS/Go generated files. Runtime code imports strict locale dictionaries, and scripts can validate coverage.

This is clean long term, but a large migration if done before filling the missing translations. It also creates churn across the app, bot, and public site at the same time.

### Option C: Explicit dictionaries first, migration later

Keep the existing runtime shape, but replace user-visible fallback generation with explicit 35-locale dictionaries and tests that reject missing or fallback-generated copy. Generated scripts may be used to draft translations, but committed output must be explicit and reviewable.

Recommended. It delivers the user-visible requirement with less architectural risk. After the translation release is stable, the dictionaries can be moved into separate files in a smaller refactor.

## Architecture

The implementation will define one strict policy:

User-facing copy must resolve from an explicit locale dictionary for every supported interface language. Runtime fallbacks are allowed only as development safety for missing keys and must be treated as test failures when reachable from production UI.

Allowed untranslated invariant terms:

- Brand and platform names: `Poliglot AI`, `Telegram`, `AI Tutor`.
- Plan and payment names: `Free`, `Premium`, `Platinum`, `Stars`, `TON`, `USDT`, `RUB`, `YooKassa`, `SBP`.
- Technical labels that users recognize as product/protocol tokens: `XP`, `CEFR`, `OCR`, `PWA`, email addresses, URLs, bot handles, legal identifiers, prices, dates, and exact support contacts.

Everything else should be natural in the selected interface language, including explanatory sentences, button labels, onboarding text, guide text, payment/referral messages, errors, tool prompts, legal navigation, and helper language names.

## Web App

Primary files:

- `web-react/src/lib/i18n.ts`
- `web-react/src/App.tsx`
- `web-react/src/components/ui/sign-in.tsx`
- `web-react/src/components/ui/otpdialog.tsx`
- `web-react/e2e/web-smoke.spec.ts`
- `web_app_shell_test.go`

Required behavior:

- `app_guide_*` must be full translated prose for every non-Russian locale, not output from `buildDerivedGuideCopy`.
- `localizedFallbackForEnglishCopy` must not be used to repair production copy for `auth_*`, `payment_*`, `referral_*`, `tools_*`, `tutor_*`, `pronunciation_*`, `roleplay_*`, `onboarding_*`, `cookie_*`, `today_*`, `premium_*`, `platinum_*`, `free_*`, `not_found_*`, and core error keys.
- `fallbackLeakAllowedKeys` must not exempt guide keys or other prose keys.
- Hardcoded placeholders in `sign-in.tsx` and `otpdialog.tsx` must receive localized labels from the parent translation layer or local dictionary props.
- Unknown locale normalization can still fall back to `ru` for invalid input, but every listed supported locale must be fully covered.

Acceptance checks:

- Static test that all literal `copy("key", "fallback")` keys used by `App.tsx` are represented in the translation inventory.
- Browser test that sweeps all 35 locales through high-risk screens and rejects English/Russian source phrases outside the allowed invariant list.
- Source test that fails when `buildDerivedGuideCopy` or production guide exemptions remain.

## Public Site And Legal Pages

Primary files:

- `site-react/src/PublicSiteApp.tsx`
- `site-react/src/legacyLegalContent.ts`
- `site-react/public/assets/site-i18n.js`
- `site-react/public/assets/site-phrases.js`
- `site-react/public/assets/legal-documents-i18n.js`
- `site-react/public/assets/privacy-policy-i18n.js`
- `tools/rebuild_public_site_translations.mjs`
- `tools/generate_legal_documents_i18n.mjs`
- `site-react/e2e/public-site.spec.ts`
- `site-react/e2e/deploy-output.spec.ts`
- `Сайт полиглота для бота/*`

Required behavior:

- Landing page product copy must have explicit 35-locale translations, not runtime English/Russian fallback through the phrase dictionary.
- Legal pages `privacy`, `terms`, `agreement`, and `consent` must use one 35-locale renderer. The 20-locale `privacy-policy-i18n.js` path must be removed from runtime loading or made a non-conflicting compatibility asset.
- Legal document translations must preserve exact operator details, support contacts, payment/provider names, bot handles, URLs, dates, and legal identifiers.
- Generated static output in `Сайт полиглота для бота` must mirror the verified public-site assets after build.

Acceptance checks:

- Public-site Playwright locale sweep must verify all 35 languages for landing and legal pages.
- Tests must reject old `20 languages` copy in any locale.
- Tests must verify that `privacy.html` is not rendered by competing 20-locale and 35-locale paths.
- Encoding artifact check must pass after generated assets are rebuilt.

## Telegram Bot And Go Backend

Primary files:

- `i18n.go`
- `i18n_cjk.go`
- `system_i18n.go`
- `system_i18n_clean_reminders.go`
- `premium_i18n.go`
- `premium_landing_copy.go`
- `language.go`
- `language_names.go`
- `bot.go`
- `telegram.go`
- `app_prompts.json`
- `i18n_test.go`
- `premium_test.go`
- `telegram_test.go`
- `reminders_test.go`
- `vocabulary_test.go`

Required behavior:

- `ui`, `toolUICopyFor`, `systemUI`, and `premiumUI` must produce complete localized user copy for all 35 interface languages.
- `premiumLandingTerms` must have terms for all 35 locales instead of defaulting most locales to English.
- `ReferralShareText` and `InviteLinkText` must be localized for all 35 locales.
- `localizedLanguageNameForInterface` must cover the full 35-by-35 helper language name matrix where names are inserted into user-facing text.
- `interfaceLanguageSelectionText` and learning-language selection text must not use mixed English/Russian headers for localized flows.
- `vocabularyFallbackPrompt` must no longer prefer Russian then English for a non-Russian interface. If a source word has no translation in the selected interface language, the user-visible fallback must be localized or phrased as a localized missing-translation message.
- `app_prompts.json` is internal model instruction copy, not direct UI copy. It should keep instructing AI output to use the selected interface language. If any prompt text is shown directly to users, that direct surface must be added to the 35-locale copy inventory.

Acceptance checks:

- Go matrix tests over all 35 interface languages for `ui`, `systemUI`, `premiumUI`, `premiumPlanLandingCopy`, `localizedLanguageNameForInterface`, referral/invite text, reminders, tools, level flow, and vocabulary fallback behavior.
- Existing tests that encode Russian fallback must be rewritten to assert localized behavior.
- `go test ./...` must pass.

## Translation Workflow

The Russian guide and current Russian product/legal copy are the source of tone for literary translations. English can be used as an auxiliary source where it already reflects product terminology, but it is not the fallback authority for non-English locales.

For every surface:

1. Extract the source copy inventory.
2. Group keys by product context, not alphabetically, so translators preserve voice and meaning.
3. Draft or import translations for all 35 locales.
4. Commit translations as explicit code/data dictionaries.
5. Run coverage tests and rendered smoke tests.
6. Review rendered output for high-risk locales and legal pages.

Machine translation is allowed as a drafting mechanism only. The committed release must be explicit, deterministic, and testable.

## Quality Gates

Required verification before completion:

- `go test ./...`
- `npm --prefix web-react run build`
- `npm --prefix web-react run e2e`
- `npm --prefix site-react run build`
- `npm --prefix site-react run e2e`
- `node tools/check_encoding_artifacts.mjs`
- Targeted Playwright smoke in Chromium for web app and public site after local preview.

Git handling:

- Existing unrelated changes must not be reverted or included accidentally.
- Stage only files changed for this localization release.
- After verified changes, commit and push to `https://github.com/GarryNaxyison/My-app-local.git`.

## Residual Risks

This design can prove coverage and prevent fallback leaks. It cannot prove literary quality in every language without human or native-speaker review. The implementation should therefore make quality review easy by keeping explicit dictionaries and by adding rendered locale smoke tests that expose awkward copy early.

Legal translations are product copy, not a substitute for jurisdiction-specific legal review. Exact operator and payment details must remain unchanged across languages.

## Spec Self-Review

- Scope covers web app, public site, legal pages, generated static output, Telegram bot, and backend user-facing strings.
- No unresolved placeholder sections remain.
- The 35-locale set is explicit and consistent across surfaces.
- Runtime fallback behavior is constrained by acceptance tests instead of accepted as product behavior.
