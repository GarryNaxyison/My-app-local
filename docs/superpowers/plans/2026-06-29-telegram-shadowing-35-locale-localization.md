# Telegram Shadowing 35-Locale Localization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace damaged Telegram listening/shadowing UI strings with complete UTF-8 translations for all 35 interface languages.

**Architecture:** Keep the existing `shadowingUICopy` data model and message builders. Update only the localization map and tests, using keyed struct literals to protect field order.

**Tech Stack:** Go, standard `testing` package, existing Telegram bot code.

---

### Task 1: Add Regression Tests

**Files:**
- Modify: `shadowing_test.go`

- [ ] **Step 1: Add a test that validates all shadowing copies**

Add tests that iterate over `interfaceLanguages()`, call `shadowingCopy(userState{InterfaceLanguage: code})`, and assert every field is non-empty and free of mojibake markers.

- [ ] **Step 2: Verify the test fails before the localization patch**

Run:

```bash
go test ./... -run 'TestShadowingUICopies'
```

Expected: FAIL because existing entries contain stored `????` strings.

### Task 2: Replace Shadowing Translations

**Files:**
- Modify: `shadowing.go`

- [ ] **Step 1: Replace `shadowingUICopies`**

Use complete keyed `shadowingUICopy` entries for `ru`, `en`, `es`, `de`, `fr`, `it`, `zh`, `ja`, `ko`, `tg`, `uz`, `tt`, `hy`, `kk`, `ky`, `ka`, `uk`, `pl`, `ro`, `pt`, `ar`, `bn`, `cs`, `el`, `hi`, `hu`, `id`, `nl`, `sv`, `ta`, `te`, `th`, `tl`, `tr`, and `vi`.

- [ ] **Step 2: Simplify `shadowingCopy`**

Return the map entry directly when present, otherwise return English copy. Do not keep runtime patches for Russian fields.

- [ ] **Step 3: Verify focused tests**

Run:

```bash
go test ./... -run 'TestShadowing'
```

Expected: PASS.

### Task 3: Final Verification

**Files:**
- No new files beyond spec, plan, `shadowing.go`, `shadowing_test.go`.

- [ ] **Step 1: Run Go tests**

Run:

```bash
go test ./...
```

Expected: PASS or report any unrelated timeout/failure clearly.

- [ ] **Step 2: Review diff**

Run:

```bash
git diff -- shadowing.go shadowing_test.go docs/superpowers/specs/2026-06-29-telegram-shadowing-35-locale-localization-design.md docs/superpowers/plans/2026-06-29-telegram-shadowing-35-locale-localization.md
```

Expected: only the planned Telegram localization/test/spec changes.
