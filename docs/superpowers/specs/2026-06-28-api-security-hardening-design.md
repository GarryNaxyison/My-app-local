# API Security Hardening Design

## Goal

Reduce replay, brute-force, and low-cost abuse risk without adding memory-heavy dependencies or broad rewrites.

## Scope

- Add safe server timeouts to the Go HTTP server.
- Make API rate-limit client identity trust proxy headers only when the request came from a trusted proxy address.
- Support stable idempotency for payment creation requests through `Idempotency-Key`.
- Reuse existing payment application dedupe for webhooks and successful payment callbacks.
- Document edge rate limiting requirements in the checked-in Caddy candidate.

## Non-Goals

- This does not replace provider-level DDoS mitigation, WAF, or Cloudflare-style edge protection.
- This does not persist a full generic idempotency cache for every mutating API endpoint.
- This does not run large load tests in the local workspace.

## Design

Rate limiting remains in-process and uses the existing bucket map. The client key changes so `X-Forwarded-For` and `X-Real-IP` are used only when `RemoteAddr` is loopback/private or explicitly configured through `TRUSTED_PROXY_CIDRS`.

Payment creation endpoints accept an optional `Idempotency-Key` header. When present, YooKassa and RollyPay use a deterministic order/idempotence value derived from user, product, channel, and key. Direct crypto uses the same deterministic key to return an existing pending payment for a retry instead of creating another pending invoice.

The HTTP server gets read/write/header/idle timeouts to reduce slow connection abuse. Caddy config documents the expected production edge limit layer because application code cannot absorb volumetric DDoS.

## Testing

Add targeted Go tests only:

- proxy spoofing does not bypass rate buckets unless the proxy is trusted;
- trusted proxy CIDRs allow `X-Forwarded-For`;
- server timeouts are configured;
- YooKassa idempotence keys are stable for the same logical retry;
- direct crypto idempotency returns the same pending payment.
