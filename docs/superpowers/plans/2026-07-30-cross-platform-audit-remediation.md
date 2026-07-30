# Cross-Platform Audit Remediation Plan

> **For agentic workers:** Execute each task with a focused regression test before the production change.

**Goal:** Remove the confirmed release blockers and reduce the highest-risk operational gaps across the Telegram backend, Android client, and web client.

**Architecture:** Preserve server contracts as the source of truth. Android fixes change its DTO contract and report network errors explicitly. Backend configuration becomes fail-closed for production cookies and externalizes Android signing secrets. Web delivery is split at third-party module boundaries without changing routes.

**Tech Stack:** Go, Kotlin/Compose, Gradle, React, Vite, Playwright.

---

### Task 1: Repair Android mistake deletion

- [ ] Add a failing repository/API contract test proving the client sends a zero-based `index` field.
- [ ] Replace the `IdRequest` mistake-delete DTO with a dedicated `MistakeDeleteRequest(index)` contract.
- [ ] Pass `MistakeItem.index` from the UI, surface an error on failure, and run Android unit tests.

### Task 2: Make Android release signing secret-driven

- [ ] Replace embedded signing credentials with required Gradle properties/environment variables.
- [ ] Verify debug builds without secrets and release builds with ephemeral environment variables.

### Task 3: Repair the backend test gate and package discovery

- [ ] Update stale PWA/component assertions to current production assets.
- [ ] Move deployment build artifacts outside the Go module or make their directory ignored by Go package discovery.
- [ ] Run `go test .` and `go test ./...`.

### Task 4: Harden session configuration and Telegram retries

- [ ] Add configuration tests for Secure cookies and stable session-secret requirements.
- [ ] Reject insecure production configuration and avoid random session keys in production.
- [ ] Add bounded retry/error notification handling for Telegram update failures without reprocessing updates indefinitely.

### Task 5: Improve client resilience and web delivery

- [ ] Introduce visible Android error/retry states for critical learning flows and cover the state reducer/unit contract.
- [ ] Split the web vendor and authentication-scene chunks; enforce a build-output budget check.
- [ ] Run Android, Go, web build, and Playwright verification.

### Task 6: Device verification

- [ ] Add Android instrumented smoke coverage for startup, authenticated navigation, phrasebook deletion, and error presentation.
- [ ] Run it on an attached emulator/device; until then, document it as an environment blocker.
