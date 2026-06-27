# API Security Hardening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Harden API replay and abuse controls without memory-heavy tooling or broad rewrites.

**Architecture:** Keep the existing Go HTTP stack and in-process limiter. Add trusted proxy parsing, stable payment idempotency hooks, HTTP server timeouts, and deploy guidance for edge limits.

**Tech Stack:** Go `net/http`, existing SQLite/JSON stores, existing payment clients, Caddy config snapshots.

---

### Task 1: Trusted Client IP And Server Timeouts

**Files:**
- Modify: `config.go`
- Modify: `main.go`
- Modify: `web_api.go`
- Test: `web_api_feature_test.go`

- [ ] Add failing tests for untrusted `X-Forwarded-For`, trusted proxy CIDR behavior, and HTTP server timeout values.
- [ ] Implement `TRUSTED_PROXY_CIDRS`, `clientRateKey`, and server timeout helpers.
- [ ] Run: `go test -run 'Test(ClientRateKey|HTTPServerTimeouts)' .`

### Task 2: Payment Creation Idempotency

**Files:**
- Modify: `web_api.go`
- Modify: `yookassa.go`
- Modify: `rollypay.go`
- Modify: `crypto_payments.go`
- Test: `premium_test.go`
- Test: `crypto_payments_test.go`

- [ ] Add failing tests for stable YooKassa `Idempotence-Key` and direct crypto retry returning the same pending payment.
- [ ] Implement optional `Idempotency-Key` support for web payment creation.
- [ ] Run: `go test -run 'Test(YooKassa.*Idempot|DirectCrypto.*Idempot)' .`

### Task 3: Edge Limit Documentation

**Files:**
- Modify: `deploy/caddy/Caddyfile.updated`
- Modify: `docs/technical/DOCUMENTATION.md`

- [ ] Add explicit production note that Caddy snapshot still requires edge/WAF rate limiting for volumetric DDoS.
- [ ] Do not add unverified Caddy plugins or directives that may fail on stock Caddy.

### Task 4: Verification And Git Sync

**Files:**
- Review all changed files.

- [ ] Run targeted Go tests only.
- [ ] Run `gofmt` on changed Go files.
- [ ] Run `git diff --check`.
- [ ] Commit and push if verification passes.
