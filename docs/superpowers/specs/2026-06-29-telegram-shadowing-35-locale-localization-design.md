# Telegram Shadowing 35-Locale Localization Design

## Goal

Fix the Telegram listening/shadowing result copy so every supported interface language sends real UTF-8 text instead of stored question marks or mojibake.

## Scope

This change covers the Telegram bot listening feature implemented in `shadowing.go`. It updates all 35 interface languages returned by `interfaceLanguages()` and keeps the current message structure: title, target phrase, learner answer, accuracy, XP bonus, feedback, and next action.

Web mojibake found in `site-react`, `web-react`, and built `web/assets` is a separate bug and is not bundled into this Telegram fix.

## Approach

Replace the damaged positional `shadowingUICopies` entries with complete, natural translations for all supported interface language codes. Use keyed struct literals so future edits do not accidentally shift fields. Remove the Russian-only runtime patch once the map itself contains correct Russian text.

## Acceptance Criteria

- `shadowingUICopies` contains entries for all 35 interface languages.
- Each entry has every required field filled.
- Listening result messages for Russian contain readable labels such as `Фраза`, `В ответе`, `Совпадение`, and `Опыт`.
- No shadowing UI copy contains replacement characters, mojibake markers, or runs of question marks.
- Focused Go tests pass.

## Verification

Run:

```bash
go test ./... -run 'TestShadowing'
```

Then run the broader package tests if time permits:

```bash
go test ./...
```
