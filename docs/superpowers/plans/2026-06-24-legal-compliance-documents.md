# Legal Compliance Documents Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add complete legal-document coverage, cookie opt-in, and consent links across public site, WebApp, and Telegram.

**Architecture:** Keep the existing React public-site legal shell and add two new HTML entrypoints. Keep Telegram and WebApp behavior backward-compatible by adding links and payload fields rather than replacing flows.

**Tech Stack:** Go Telegram bot/backend, React/Vite public site, React/Vite WebApp, Playwright, Go tests.

---

## File Map

- Modify `site-react/src/legacyLegalContent.ts`: add exported legal document HTML strings with operator details.
- Modify `site-react/src/PublicSiteApp.tsx`: add `agreement` and `consent` page IDs, nav/footer/asides, and public cookie banner.
- Add `site-react/agreement.html` and `site-react/consent.html`: static Vite entrypoints.
- Modify `site-react/vite.config.ts`: include the new entrypoints in Rollup input.
- Modify `site-react/e2e/public-site.spec.ts`: assert legal pages and cookie banner.
- Modify `site-react/e2e/deploy-output.spec.ts`: assert built new HTML files.
- Modify `telegram.go`: add document URLs and three legal buttons.
- Modify `telegram_test.go`: assert Telegram prompt links all documents.
- Modify `web-react/src/App.tsx`: add WebApp cookie banner and expand auth/legal consent links.
- Modify `web-react/e2e/web-smoke.spec.ts`: assert auth consent links and cookie banner.

## Tasks

### Task 1: Public Legal Pages

- [ ] Add `/agreement.html` and `/consent.html` entry files mirroring existing legal HTML.
- [ ] Add `agreementDocumentHtml` and `personalDataConsentHtml`.
- [ ] Extend `PageId`, `getPage`, `LegalPage`, nav drawer, legal aside, and footer.
- [ ] Add public-site Playwright assertions for all four routes.

### Task 2: Cookie Opt-In UI

- [ ] Add a cookie banner component to public site using localStorage key `poliglot-cookie-consent`.
- [ ] Add the same banner pattern to WebApp using localStorage key `poliglot-app-cookie-consent`.
- [ ] Add CSS for compact fixed banners that do not overlap content on mobile.
- [ ] Add Playwright assertions for banner visibility, links, and accept/dismiss behavior.

### Task 3: Telegram Legal Prompt

- [ ] Add constants for `/agreement.html` and `/consent.html`.
- [ ] Update `privacyPromptKeyboard()` to show continue plus privacy, consent, and agreement URL buttons.
- [ ] Update `telegram_test.go` expectations for the new rows.

### Task 4: WebApp Consent Links

- [ ] Add legal URL helpers for privacy, consent, and agreement with current interface language query.
- [ ] Expand registration and Telegram-link checkbox text to show all three links.
- [ ] Add consent/agreement URLs to JSON payloads when starting registration and Telegram linking.
- [ ] Add focused WebApp smoke assertions.

### Task 5: Verification, Review, Commit, Push

- [ ] Run focused Go test for Telegram prompt.
- [ ] Run public-site build.
- [ ] Run WebApp build.
- [ ] Run focused Playwright tests if preview build succeeds.
- [ ] Run repository check command where feasible.
- [ ] Request code review and fix actionable findings.
- [ ] Commit and push to the configured GitHub remote.
