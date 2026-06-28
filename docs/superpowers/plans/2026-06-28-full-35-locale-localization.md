# Full 35-Locale Localization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace generated and mixed-language user copy with explicit literary localization for all 35 interface languages across the WebApp, public site/legal pages, static public output, and Telegram bot.

**Architecture:** Add failing coverage gates first, then replace each fallback path with explicit locale dictionaries while preserving the current runtime shape. Keep implementation slices independent: shared checks, WebApp, public/legal, Go bot, final integration.

**Tech Stack:** Go, TypeScript, React, Vite, Playwright, Node scripts, Nx run-commands, GitHub remote sync.

---

## File Map

- Create `tools/i18n/allowed_invariant_terms.json`: shared allow-list for brands, payment tokens, URLs, legal identifiers, and technical names that may remain untranslated.
- Create `tools/i18n/check_web_i18n_coverage.mjs`: static WebApp coverage checker for `copy("key", "fallback")` calls and explicit locale dictionaries.
- Create `web-react/src/lib/fullLocaleCopy.ts`: explicit WebApp copy dictionaries for all 35 locale codes.
- Modify `web-react/src/lib/i18n.ts`: import explicit dictionaries, remove guide/prose derived fallback from production copy, and tighten fallback leak checks.
- Modify `web-react/src/App.tsx`: route missing hardcoded copy through `appCopy` where needed.
- Modify `web-react/src/components/ui/sign-in.tsx`: require localized labels from props rather than English fallback text.
- Modify `web-react/src/components/ui/otpdialog.tsx`: require localized labels from props rather than English fallback text.
- Modify `web_app_shell_test.go`: add static source tests for full WebApp localization coverage and removed derived guide behavior.
- Modify `web-react/e2e/web-smoke.spec.ts`: add rendered all-locale smoke checks for guide, auth, premium, tools, tutor, onboarding, and errors.
- Modify `site-react/public/assets/site-i18n.js`: make landing copy explicit for all 35 locales and remove English/Russian fallback for product prose.
- Modify `site-react/public/assets/legal-documents-i18n.js`: keep all four legal documents covered for all 35 locales.
- Modify `site-react/public/assets/privacy-policy-i18n.js`: remove runtime conflict or convert it to a compatibility no-op when `legal-documents-i18n.js` owns the page.
- Modify `site-react/privacy.html`: load only the 35-locale legal renderer for privacy.
- Modify `tools/rebuild_public_site_translations.mjs`: make generated phrase output deterministic and coverage-checked for all 35 locales.
- Modify `tools/generate_legal_documents_i18n.mjs`: fail if any legal page lacks a 35-locale translation pack.
- Modify `site-react/e2e/public-site.spec.ts`: add all-locale landing and legal checks and assert no competing privacy renderer.
- Modify `site-react/e2e/deploy-output.spec.ts`: assert static output contains the verified 35-locale assets.
- Modify `Сайт полиглота для бота/*.html` and `Сайт полиглота для бота/assets/*.js`: mirror verified public-site build output.
- Create `i18n_full_locales.go`: explicit Telegram menu/tool UI copy overlays for all 35 locales.
- Create `system_i18n_full_locales.go`: explicit system/reminder/level/practice copy overlays for all 35 locales.
- Create `premium_i18n_full_locales.go`: explicit premium, payment, referral, and invite copy overlays for all 35 locales.
- Create `language_names_full.go`: full helper language-name matrix for all 35 interface languages.
- Modify `i18n.go`: apply full locale copy before generated aliases and remove production English fallback for supported locales.
- Modify `system_i18n.go`: apply full system copy and stop compact fallback from producing user-visible prose for supported locales.
- Modify `premium_i18n.go`: apply full premium copy and cover referral/invite text for all locales.
- Modify `premium_landing_copy.go`: cover `premiumLandingTerms` for all 35 locales.
- Modify `language.go`: localize language-selection headers.
- Modify `language_names.go`: use the full helper-name matrix before native/English fallback.
- Modify `bot.go`: replace Russian/English vocabulary fallback behavior with localized missing-translation behavior.
- Modify `i18n_test.go`, `premium_test.go`, `telegram_test.go`, `reminders_test.go`, and `vocabulary_test.go`: add and update full-locale coverage tests.

## Shared Locale Rules

Use this exact locale list in every new checker and dictionary:

```text
ru,en,es,de,fr,it,zh,ja,ko,tg,uz,tt,hy,kk,ky,ka,uk,pl,ro,pt,ar,bn,cs,el,hi,hu,id,nl,sv,ta,te,th,tl,tr,vi
```

Allowed invariant terms:

```json
[
  "Poliglot AI",
  "Telegram",
  "AI Tutor",
  "Free",
  "Premium",
  "Platinum",
  "Stars",
  "TON",
  "USDT",
  "RUB",
  "YooKassa",
  "SBP",
  "XP",
  "CEFR",
  "OCR",
  "PWA",
  "@poliglot_ai_bot",
  "@AsaselD",
  "supportpoliglotai@gmail.com",
  "poliglotai.ru",
  "poliglotai.online"
]
```

## Tasks

### Task 1: Shared Coverage Gate

**Files:**
- Create: `tools/i18n/allowed_invariant_terms.json`
- Create: `tools/i18n/check_web_i18n_coverage.mjs`
- Modify: `web_app_shell_test.go`

- [ ] **Step 1: Write the failing WebApp source gate**

Add a Go test that requires the checker script and rejects production guide fallback functions.

```go
func TestReactFullLocalizationCoverageGateExists(t *testing.T) {
	source, err := os.ReadFile(filepath.Join("tools", "i18n", "check_web_i18n_coverage.mjs"))
	if err != nil {
		t.Fatalf("missing WebApp i18n coverage checker: %v", err)
	}
	text := string(source)
	for _, want := range []string{
		"appLocaleCodes",
		"fullLocaleCopy",
		"copy(",
		"allowed_invariant_terms.json",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("coverage checker is missing marker %q", want)
		}
	}
}

func TestReactGuideDoesNotUseDerivedProductionCopy(t *testing.T) {
	source, err := os.ReadFile(filepath.Join("web-react", "src", "lib", "i18n.ts"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	for _, forbidden := range []string{
		"function buildDerivedGuideCopy",
		"key.startsWith(\"app_guide_\") return false",
		"localizedGuideCopy[code as LocalizedGuideCode] || buildDerivedGuideCopy",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("WebApp still allows derived guide copy: %q", forbidden)
		}
	}
}
```

- [ ] **Step 2: Run the failing source gate**

Run:

```powershell
go test ./... -run "TestReactFullLocalizationCoverageGateExists|TestReactGuideDoesNotUseDerivedProductionCopy"
```

Expected: fail because `tools/i18n/check_web_i18n_coverage.mjs` does not exist and derived guide copy still exists.

- [ ] **Step 3: Add the invariant term allow-list**

Create `tools/i18n/allowed_invariant_terms.json` with the JSON array from the Shared Locale Rules section.

- [ ] **Step 4: Add the WebApp coverage checker**

Create `tools/i18n/check_web_i18n_coverage.mjs`:

```js
#!/usr/bin/env node

import fs from "node:fs";
import path from "node:path";

const root = process.cwd();
const localeCodes = "ru,en,es,de,fr,it,zh,ja,ko,tg,uz,tt,hy,kk,ky,ka,uk,pl,ro,pt,ar,bn,cs,el,hi,hu,id,nl,sv,ta,te,th,tl,tr,vi".split(",");
const appPath = path.join(root, "web-react", "src", "App.tsx");
const copyPath = path.join(root, "web-react", "src", "lib", "fullLocaleCopy.ts");
const allowedPath = path.join(root, "tools", "i18n", "allowed_invariant_terms.json");

const appSource = fs.readFileSync(appPath, "utf8");
const copySource = fs.readFileSync(copyPath, "utf8");
const allowedTerms = JSON.parse(fs.readFileSync(allowedPath, "utf8"));

const keyPattern = /\bcopy\(\s*["'`]([^"'`]+)["'`]\s*,/g;
const keys = [...appSource.matchAll(keyPattern)].map((match) => match[1]);
const uniqueKeys = [...new Set(keys)].sort();

const missing = [];
for (const code of localeCodes) {
  if (!copySource.includes(`${code}: {`)) missing.push(`locale:${code}`);
  for (const key of uniqueKeys) {
    if (!copySource.includes(`${key}:`)) missing.push(`${code}:${key}`);
  }
}

for (const term of allowedTerms) {
  if (typeof term !== "string" || !term.trim()) {
    throw new Error("allowed invariant terms must be non-empty strings");
  }
}

if (missing.length) {
  console.error(`Missing explicit WebApp locale entries (${missing.length}):`);
  for (const item of missing.slice(0, 120)) console.error(`- ${item}`);
  process.exit(1);
}
```

- [ ] **Step 5: Wire checker into the Go source gate**

Add this assertion to `TestReactFullLocalizationCoverageGateExists`:

```go
cmd := exec.Command("node", filepath.Join("tools", "i18n", "check_web_i18n_coverage.mjs"))
cmd.Dir = "."
output, err := cmd.CombinedOutput()
if err != nil {
	t.Fatalf("WebApp i18n coverage checker failed:\n%s", output)
}
```

- [ ] **Step 6: Run the source gate again**

Run:

```powershell
go test ./... -run "TestReactFullLocalizationCoverageGateExists|TestReactGuideDoesNotUseDerivedProductionCopy"
```

Expected: fail only on missing `web-react/src/lib/fullLocaleCopy.ts` and existing derived guide logic. This proves the gate catches the implementation gap.

### Task 2: WebApp Explicit 35-Locale Copy

**Files:**
- Create: `web-react/src/lib/fullLocaleCopy.ts`
- Modify: `web-react/src/lib/i18n.ts`
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/components/ui/sign-in.tsx`
- Modify: `web-react/src/components/ui/otpdialog.tsx`
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Create the explicit WebApp dictionary**

Create `web-react/src/lib/fullLocaleCopy.ts` with this shape:

```ts
export const fullLocaleCopy = {
  ru: {
    guide: "Гайд",
    app_guide_title: "Как начать",
    app_guide_body: "Пять страниц по главным сценариям: уроки, тарифы, тема, словарь, заметки, аудио и отчеты.",
  },
  en: {
    guide: "Guide",
    app_guide_title: "Quick start guide",
    app_guide_body: "Five pages for the main app flows: lessons, plans, theme, vocabulary, notes, audio, and reports.",
  },
} as const satisfies Record<string, Record<string, string>>;
```

Then expand every locale block to the 35-code list and include every key reported by `node tools/i18n/check_web_i18n_coverage.mjs`. Use the existing Russian guide as the tone source. Preserve allowed invariant terms exactly.

- [ ] **Step 2: Import explicit copy into `i18n.ts`**

Add near the top of `web-react/src/lib/i18n.ts`:

```ts
import { fullLocaleCopy } from "./fullLocaleCopy";
```

Merge it after existing narrow overrides and before `appCopy` returns user text:

```ts
appLocaleCodes.forEach((code) => {
  localeOverrides[code] = {
    ...(localeOverrides[code] || {}),
    ...(fullLocaleCopy[code] || {}),
  };
});
```

- [ ] **Step 3: Remove guide derived production behavior**

Delete `buildDerivedGuideCopy` and change guide merge logic to require explicit `fullLocaleCopy` values:

```ts
const guideValues: Record<string, string> = code === "ru" || code === "en"
  ? {}
  : fullLocaleCopy[code] || {};
```

Remove the `if (key.startsWith("app_guide_")) return false;` exemption from `isLikelyFallbackLeak`.

- [ ] **Step 4: Localize hardcoded auth and OTP component fallbacks**

In `sign-in.tsx`, replace English fallback expressions such as:

```tsx
placeholder={labels.loginPlaceholder || "Enter your email address"}
```

with required label values:

```tsx
placeholder={labels.loginPlaceholder}
```

In `otpdialog.tsx`, use localized props for titles, descriptions, button text, and input labels. If a required prop is missing, render an empty string only in development-safe paths and let the coverage checker catch missing copy.

- [ ] **Step 5: Add rendered WebApp all-locale smoke**

In `web-react/e2e/web-smoke.spec.ts`, add a locale sweep:

```ts
const appLocaleCodes = ["ru", "en", "es", "de", "fr", "it", "zh", "ja", "ko", "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt", "ar", "bn", "cs", "el", "hi", "hu", "id", "nl", "sv", "ta", "te", "th", "tl", "tr", "vi"];

test("WebApp renders explicit localized copy across all interface locales", async ({ page }) => {
  test.setTimeout(240_000);
  for (const code of appLocaleCodes) {
    await page.goto(`/app/?lang=${code}`);
    await expect(page.locator("body")).toBeVisible();
    const text = await page.locator("body").innerText();
    expect(text).not.toMatch(/Рџ|РЎ|Р’|Рќ|Рњ/);
    if (code !== "en") {
      expect(text).not.toContain("Quick start guide");
      expect(text).not.toContain("Enter your email address");
      expect(text).not.toContain("Request failed");
    }
    if (!["ru", "tg", "tt", "kk", "ky", "uk"].includes(code) && code !== "en") {
      expect(text).not.toMatch(/\p{Script=Cyrillic}/u);
    }
  }
});
```

Adjust page URL and selectors to the existing WebApp routing if the app uses hash or auth state in the current test harness.

- [ ] **Step 6: Run RED then GREEN checks**

Run before implementation is complete:

```powershell
node tools/i18n/check_web_i18n_coverage.mjs
go test ./... -run "TestReactFullLocalizationCoverageGateExists|TestReactGuideDoesNotUseDerivedProductionCopy"
```

Expected before full dictionary: fail with missing keys.

Run after dictionary and `i18n.ts` changes:

```powershell
node tools/i18n/check_web_i18n_coverage.mjs
go test ./... -run "TestReactFullLocalizationCoverageGateExists|TestReactGuideDoesNotUseDerivedProductionCopy"
npm --prefix web-react run build
```

Expected after implementation: all commands exit 0.

### Task 3: Public Site And Legal 35-Locale Runtime

**Files:**
- Modify: `site-react/privacy.html`
- Modify: `site-react/public/assets/site-i18n.js`
- Modify: `site-react/public/assets/site-phrases.js`
- Modify: `site-react/public/assets/legal-documents-i18n.js`
- Modify: `site-react/public/assets/privacy-policy-i18n.js`
- Modify: `tools/rebuild_public_site_translations.mjs`
- Modify: `tools/generate_legal_documents_i18n.mjs`
- Modify: `site-react/e2e/public-site.spec.ts`
- Modify: `site-react/e2e/deploy-output.spec.ts`

- [ ] **Step 1: Write failing public-site tests**

In `site-react/e2e/public-site.spec.ts`, add assertions that `privacy.html` does not load the 20-locale privacy renderer:

```ts
test("privacy uses the 35-locale legal renderer only", async ({ page }) => {
  const requests: string[] = [];
  page.on("request", (request) => requests.push(new URL(request.url()).pathname));
  await page.goto("/privacy.html?lang=vi");
  await expect(page.locator(".legal-document-shell")).toBeVisible();
  expect(requests).toContain("/assets/legal-documents-i18n.js");
  expect(requests).not.toContain("/assets/privacy-policy-i18n.js");
  const text = await page.locator(".legal-document-shell").innerText();
  expect(text).not.toContain("Политика обработки персональных данных");
});
```

Extend the existing all-locale landing test to include legal pages:

```ts
for (const legalPath of ["/privacy.html", "/terms.html", "/agreement.html", "/consent.html"]) {
  for (const code of siteLocaleCodes) {
    await page.goto(`${legalPath}?lang=${code}`);
    await expect(page.locator(".legal-document-shell")).toBeVisible();
    const text = await page.locator(".legal-document-shell").innerText();
    expect(text).toContain("supportpoliglotai@gmail.com");
    expect(text).toContain("505017471160");
    if (code !== "ru") {
      expect(text).not.toContain("Пользовательское соглашение");
      expect(text).not.toContain("Согласие на обработку персональных данных");
    }
  }
}
```

- [ ] **Step 2: Run the failing public-site test**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "privacy uses the 35-locale legal renderer only|legal pages"
```

Expected: fail because `privacy.html` currently loads both privacy and legal i18n assets.

- [ ] **Step 3: Remove the competing privacy renderer from runtime**

In `site-react/privacy.html`, remove:

```html
<script defer src="/assets/privacy-policy-i18n.js?v=20260527-data-policy"></script>
```

Keep:

```html
<script defer src="/assets/legal-documents-i18n.js?v=20260624-legal"></script>
```

If a static deployed `privacy.html` mirrors this file, apply the same change after build output is generated.

- [ ] **Step 4: Make generators fail on missing 35-locale output**

In `tools/generate_legal_documents_i18n.mjs`, after building each document payload, add:

```js
for (const code of SITE_LANG_CODES) {
  if (!payload.localized[code]) {
    throw new Error(`Missing legal translation for ${pageId}:${code}`);
  }
}
```

In `tools/rebuild_public_site_translations.mjs`, add an equivalent check for every phrase that appears on landing or legal pages.

- [ ] **Step 5: Rebuild public translations**

Run:

```powershell
node tools/rebuild_public_site_translations.mjs
$env:POLIGLOT_TRANSLATE_LEGAL='1'; node tools/generate_legal_documents_i18n.mjs; Remove-Item Env:\POLIGLOT_TRANSLATE_LEGAL
```

Expected: generated assets update and every document has 35 locale entries.

- [ ] **Step 6: Build and verify public site**

Run:

```powershell
npm --prefix site-react run build
npm --prefix site-react run e2e
node tools/check_encoding_artifacts.mjs
```

Expected: all commands exit 0, with no mojibake and no `20 languages` copy in rendered pages.

### Task 4: Go Bot Failing Matrix Tests

**Files:**
- Modify: `i18n_test.go`
- Modify: `premium_test.go`
- Modify: `vocabulary_test.go`
- Modify: `telegram_test.go`
- Modify: `reminders_test.go`

- [ ] **Step 1: Add shared test helpers**

Add to `i18n_test.go`:

```go
var allowedLocalizationTerms = []string{
	"Poliglot AI", "Telegram", "AI Tutor", "Free", "Premium", "Platinum",
	"Stars", "TON", "USDT", "RUB", "YooKassa", "SBP", "XP", "CEFR", "OCR", "PWA",
	"@poliglot_ai_bot", "@AsaselD", "supportpoliglotai@gmail.com",
}

func stripAllowedLocalizationTerms(text string) string {
	out := text
	for _, term := range allowedLocalizationTerms {
		out = strings.ReplaceAll(out, term, "")
	}
	return out
}

func assertNoEnglishOrRussianFallbackLeak(t *testing.T, code string, field string, got string, english string, russian string) {
	t.Helper()
	cleaned := stripAllowedLocalizationTerms(got)
	if code != "en" && english != "" && got == english {
		t.Fatalf("%s %s fell back to English: %q", code, field, got)
	}
	if code != "ru" && russian != "" && got == russian {
		t.Fatalf("%s %s fell back to Russian: %q", code, field, got)
	}
	if strings.Contains(cleaned, "Choose bot language / Выбери язык бота") {
		t.Fatalf("%s %s contains mixed language-selection header: %q", code, field, got)
	}
}
```

- [ ] **Step 2: Add full-locale UI/system/premium test**

Add to `i18n_test.go`:

```go
func TestAllInterfaceLanguagesUseExplicitLiteraryBotCopy(t *testing.T) {
	englishUI := ui(userState{InterfaceLanguage: "en", InterfaceSelected: true})
	russianUI := ui(userState{InterfaceLanguage: "ru", InterfaceSelected: true})
	englishSystem := systemUI(userState{InterfaceLanguage: "en"})
	russianSystem := systemUI(userState{InterfaceLanguage: "ru"})
	englishPremium := premiumUI(userState{InterfaceLanguage: "en"})
	russianPremium := premiumUI(userState{InterfaceLanguage: "ru"})

	for _, language := range interfaceLanguages() {
		code := language.Code
		user := userState{InterfaceLanguage: code, InterfaceSelected: true, LearningLanguage: "en"}
		copy := ui(user)
		systemCopy := systemUI(user)
		premiumCopy := premiumUI(user)

		checks := map[string][3]string{
			"ChooseBotLang":       {copy.ChooseBotLang, englishUI.ChooseBotLang, russianUI.ChooseBotLang},
			"ChooseLearnLang":     {copy.ChooseLearnLang, englishUI.ChooseLearnLang, russianUI.ChooseLearnLang},
			"Tool.VoicePrompt":    {copy.Tool.VoicePrompt, englishUI.Tool.VoicePrompt, russianUI.Tool.VoicePrompt},
			"System.LevelStart":   {systemCopy.LevelStartText, englishSystem.LevelStartText, russianSystem.LevelStartText},
			"System.Practice":     {systemCopy.PracticeStarted, englishSystem.PracticeStarted, russianSystem.PracticeStarted},
			"Premium.Invite":      {premiumCopy.InviteLinkText, englishPremium.InviteLinkText, russianPremium.InviteLinkText},
			"Premium.Referral":    {premiumCopy.ReferralShareText, englishPremium.ReferralShareText, russianPremium.ReferralShareText},
			"Premium.BalanceHint": {premiumCopy.ReferralBalanceHint, englishPremium.ReferralBalanceHint, russianPremium.ReferralBalanceHint},
		}
		for field, values := range checks {
			if values[0] == "" {
				t.Fatalf("missing %s for %s", field, code)
			}
			assertNoEnglishOrRussianFallbackLeak(t, code, field, values[0], values[1], values[2])
		}
	}
}
```

- [ ] **Step 3: Rewrite vocabulary fallback expectations**

In `vocabulary_test.go`, replace tests that expect Russian fallback with a localized missing-translation behavior:

```go
func TestVocabularyFallbackPromptUsesInterfaceLanguageMessage(t *testing.T) {
	word := vocabularyWord{English: "station", Russian: "станция"}
	user := userState{InterfaceLanguage: "es"}
	got := vocabularyFallbackPromptForUser(word, user)
	if strings.Contains(got, "станция") || strings.Contains(got, "station") {
		t.Fatalf("expected localized missing-translation prompt, got %q", got)
	}
	if !strings.Contains(strings.ToLower(got), "traducción") {
		t.Fatalf("expected Spanish translation guidance, got %q", got)
	}
}
```

If the current function does not accept `userState`, add a new wrapper in Task 5 and keep the old function private to non-user-visible internals.

- [ ] **Step 4: Run failing Go matrix**

Run:

```powershell
go test ./... -run "TestAllInterfaceLanguagesUseExplicitLiteraryBotCopy|TestVocabularyFallbackPromptUsesInterfaceLanguageMessage"
```

Expected: fail because current code still uses generated/English/Russian fallbacks.

### Task 5: Go Bot Explicit 35-Locale Copy

**Files:**
- Create: `i18n_full_locales.go`
- Create: `system_i18n_full_locales.go`
- Create: `premium_i18n_full_locales.go`
- Create: `language_names_full.go`
- Modify: `i18n.go`
- Modify: `system_i18n.go`
- Modify: `premium_i18n.go`
- Modify: `premium_landing_copy.go`
- Modify: `language.go`
- Modify: `language_names.go`
- Modify: `bot.go`

- [ ] **Step 1: Add explicit Telegram UI overlays**

Create `i18n_full_locales.go` with:

```go
package main

var fullInterfaceUICopy = map[string]uiCopy{
	"ru": uiCopies["ru"],
	"en": englishUICopy(),
}

var fullToolUICopy = map[string]toolUICopy{
	"ru": uiCopies["ru"].Tool,
	"en": englishToolUICopy(),
}
```

Expand the maps to all 35 locale codes. Every locale entry must contain natural text for menu labels, language selection, tools, errors, and word training messages. Preserve only allowed invariant terms.

- [ ] **Step 2: Apply explicit UI overlays in `ui()`**

Modify `ui(user userState)` so supported locales use explicit full copy before any compact/generated fallback:

```go
if full, ok := fullInterfaceUICopy[code]; ok {
	copy = mergeUICopy(copy, full)
}
```

Do not allow `englishUICopy()` to be the final source for any supported non-English locale.

- [ ] **Step 3: Add explicit system copy**

Create `system_i18n_full_locales.go` with:

```go
package main

var fullSystemUICopy = map[string]systemUICopy{
	"ru": systemUICopies["ru"],
	"en": systemUICopies["en"],
}
```

Expand all 35 locale entries with natural text for level flow, practice mode, reminders, limits, and vocabulary-level questions.

Modify `systemUI(user userState)` to merge `fullSystemUICopy[code]` after the current base and before returning.

- [ ] **Step 4: Add explicit premium/referral copy**

Create `premium_i18n_full_locales.go` with:

```go
package main

var fullPremiumUICopy = map[string]premiumUICopy{
	"ru": premiumUICopyOverrides["ru"],
	"en": englishPremiumUICopy(),
}
```

Expand all 35 locale entries and include `InviteLinkText`, `ReferralShareText`, `ReferralBalanceHint`, and all payment messages.

Modify `premiumUI(user userState)` so `fullPremiumUICopy[code]` is applied for every supported locale.

- [ ] **Step 5: Complete premium landing terms**

In `premium_landing_copy.go`, expand `premiumLandingTerms` from 5 locale entries to all 35 locale entries. For each locale, localize:

```go
StarterLabel
AITutorLessons
ListeningPronunciation
VoicePhotoTools
VoiceReview
ImageTools
AIGuidedLessons
Roleplay
Review
VoicePronunciation
VoicePhotoContext
HeavyStudy
```

- [ ] **Step 6: Complete helper language names**

Create `language_names_full.go`:

```go
package main

var fullLocalizedLanguageNames = map[string]map[string]string{
	"ru": localizedLanguageNames["ru"],
	"en": localizedLanguageNames["en"],
}
```

Expand to a 35-by-35 matrix. Modify `localizedLanguageNameForInterface` so it checks `fullLocalizedLanguageNames` before `localizedLanguageNames` and before native-name fallback.

- [ ] **Step 7: Localize language-selection headers**

In `language.go`, replace hardcoded mixed headers with localized values:

```go
func interfaceLanguageSelectionText() string {
	user := userState{InterfaceLanguage: "ru", InterfaceSelected: true}
	copy := ui(user)
	var builder strings.Builder
	builder.WriteString(copy.ChooseBotLang)
	// existing per-language rows stay below
}
```

For user-specific learning-language selection, use `ui(user).ChooseLearnLang` instead of mixed English/Russian text.

- [ ] **Step 8: Replace vocabulary Russian/English fallback**

In `bot.go`, add:

```go
func vocabularyFallbackPromptForUser(word vocabularyWord, user userState) string {
	code := normalizeInterfaceLanguage(user.InterfaceLanguage)
	if translation := wordTranslationForLanguage(word, code); translation != "" {
		return translation
	}
	return localizedMissingVocabularyTranslation(code, word)
}
```

Add `localizedMissingVocabularyTranslation` in an appropriate i18n file with all 35 locale messages. Use it only for user-visible copy.

- [ ] **Step 9: Run Go GREEN checks**

Run:

```powershell
go test ./... -run "TestAllInterfaceLanguagesUseExplicitLiteraryBotCopy|TestVocabularyFallbackPromptUsesInterfaceLanguageMessage|TestEveryNonEnglishInterfaceLanguageAvoidsEnglishTelegramFallback|TestCJKInterfaceLanguagesHaveLocalizedCopies|TestPremiumPlanLandingCopyCoversAllInterfaceLanguages"
go test ./...
```

Expected: all commands exit 0.

### Task 6: Public Static Output Mirror

**Files:**
- Modify: `Сайт полиглота для бота/poliglot-ai.html`
- Modify: `Сайт полиглота для бота/privacy.html`
- Modify: `Сайт полиглота для бота/terms.html`
- Modify: `Сайт полиглота для бота/agreement.html`
- Modify: `Сайт полиглота для бота/consent.html`
- Modify: `Сайт полиглота для бота/assets/site-i18n.js`
- Modify: `Сайт полиглота для бота/assets/site-phrases.js`
- Modify: `Сайт полиглота для бота/assets/legal-documents-i18n.js`
- Modify: `Сайт полиглота для бота/assets/privacy-policy-i18n.js`

- [ ] **Step 1: Build public static output**

Run:

```powershell
npm --prefix site-react run build
```

Expected: generated public files update in `Сайт полиглота для бота`.

- [ ] **Step 2: Assert deploy output contains locale assets**

Run:

```powershell
npm --prefix site-react run e2e -- --grep "deploy|assets|language selector"
```

Expected: static output includes all required assets and all 35 locale options.

- [ ] **Step 3: Check mirrors are intentional**

Run:

```powershell
git diff -- site-react/public/assets Сайт` полиглота` для` бота/assets
```

Expected: diffs show the same logical locale content in source and static output. If PowerShell escaping is awkward, use Git Bash:

```powershell
& 'C:\Program Files\Git\bin\bash.exe' -lc "git diff -- site-react/public/assets 'Сайт полиглота для бота/assets'"
```

### Task 7: Browser Verification

**Files:**
- No source file ownership; verification only.

- [ ] **Step 1: Start WebApp preview**

Run:

```powershell
npm --prefix web-react run build
Start-Process -FilePath 'C:\Program Files\nodejs\npm.cmd' -ArgumentList @('--prefix','web-react','run','preview','--','--host','127.0.0.1','--port','4174') -WorkingDirectory 'E:\PROJECTS\New project\My app local' -WindowStyle Hidden
```

- [ ] **Step 2: Use Playwright/browser MCP for WebApp smoke**

Open:

```text
http://127.0.0.1:4174/app/?lang=es
http://127.0.0.1:4174/app/?lang=vi
http://127.0.0.1:4174/app/?lang=ar
```

Check guide, auth, premium, tools, and tutor surfaces for rendered localized copy and console errors.

- [ ] **Step 3: Start public-site preview**

Run:

```powershell
npm --prefix site-react run build
Start-Process -FilePath 'C:\Program Files\nodejs\npm.cmd' -ArgumentList @('--prefix','site-react','run','preview','--','--host','127.0.0.1','--port','4175') -WorkingDirectory 'E:\PROJECTS\New project\My app local' -WindowStyle Hidden
```

- [ ] **Step 4: Use Playwright/browser MCP for public/legal smoke**

Open:

```text
http://127.0.0.1:4175/poliglot-ai.html?lang=ja
http://127.0.0.1:4175/privacy.html?lang=vi
http://127.0.0.1:4175/terms.html?lang=ar
http://127.0.0.1:4175/agreement.html?lang=bn
http://127.0.0.1:4175/consent.html?lang=tl
```

Check rendered text, legal identifiers, language persistence, and absence of console errors.

### Task 8: Final Verification, Review, Commit, Push

**Files:**
- All changed files from Tasks 1-7.

- [ ] **Step 1: Run full verification**

Run:

```powershell
go test ./...
npm --prefix web-react run build
npm --prefix web-react run e2e
npm --prefix site-react run build
npm --prefix site-react run e2e
node tools/check_encoding_artifacts.mjs
npm run check
```

Expected: every command exits 0. If `npm run check` repeats earlier Go tests and encoding checks, still read the output and confirm it exits 0.

- [ ] **Step 2: Inspect git status**

Run:

```powershell
git status --short
git diff --stat
git diff --check
```

Expected: only localization-related files changed, no whitespace errors, and unrelated existing changes are not staged.

- [ ] **Step 3: Request code review**

Use `requesting-code-review` for the final diff. Fix actionable findings with tests first.

- [ ] **Step 4: Stage only localization files**

Run:

```powershell
git add tools/i18n web-react site-react i18n_full_locales.go system_i18n_full_locales.go premium_i18n_full_locales.go language_names_full.go i18n.go system_i18n.go premium_i18n.go premium_landing_copy.go language.go language_names.go bot.go i18n_test.go premium_test.go telegram_test.go reminders_test.go vocabulary_test.go "Сайт полиглота для бота"
```

Do not stage unrelated existing files such as `site-react/playwright-report/index.html` unless the final implementation actually updates and verifies that report intentionally.

- [ ] **Step 5: Commit and push**

Run:

```powershell
git commit -m "feat: add full 35-locale localization"
git push
```

Expected: commit succeeds and pushes to the configured GitHub remote.

## Execution Strategy

Use subagent-driven development for implementation:

- Worker A owns WebApp files only: `web-react`, `web_app_shell_test.go`, `tools/i18n/check_web_i18n_coverage.mjs`.
- Worker B owns public site/legal files only: `site-react`, `tools/rebuild_public_site_translations.mjs`, `tools/generate_legal_documents_i18n.mjs`, `Сайт полиглота для бота`.
- Worker C owns Go bot i18n files only: root Go i18n files and Go tests.
- Tester/QA owns verification commands and rendered smoke checks after workers finish.
- Reviewer is read-only and runs after the integrated diff exists.

The orchestrator integrates worker outputs, resolves overlap, runs final verification, commits, and pushes.

## Plan Self-Review

- Spec coverage: covers WebApp, public landing, legal pages, static output, Telegram bot, backend user-facing strings, fallback bans, tests, browser smoke, commit, and push.
- Placeholder scan: no unresolved implementation markers remain; translation work is expressed as explicit dictionary files with required schemas and locale set.
- Type consistency: new WebApp dictionary uses `fullLocaleCopy`; Go dictionary files use existing `uiCopy`, `toolUICopy`, `systemUICopy`, and `premiumUICopy` types; tests reference functions named in implementation steps.
