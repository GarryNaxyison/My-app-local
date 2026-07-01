import { expect, test } from "@playwright/test";

test("Cloudflare Worker fallback path matcher only proxies dynamic application paths", async () => {
  const worker = await import("../cloudflare/worker.js");

  expect(worker.shouldProxyPath("/app")).toBe(true);
  expect(worker.shouldProxyPath("/app/lesson")).toBe(true);
  expect(worker.shouldProxyPath("/login")).toBe(true);
  expect(worker.shouldProxyPath("/api/profile")).toBe(true);
  expect(worker.shouldProxyPath("/healthz")).toBe(true);
  expect(worker.shouldProxyPath("/payment/success")).toBe(true);
  expect(worker.shouldProxyPath("/tonapi/webhook")).toBe(true);
  expect(worker.shouldProxyPath("/rollypay/webhook/site")).toBe(true);
  expect(worker.shouldProxyPath("/rollypay/webhook/anything")).toBe(true);
  expect(worker.shouldProxyPath("/poliglot-ai.html")).toBe(false);
  expect(worker.shouldProxyPath("/assets/site-react/main.js")).toBe(false);
});

test("Cloudflare Worker returns branded maintenance HTML when the origin is unavailable", async () => {
  const worker = await import("../cloudflare/worker.js");
  const originalFetch = globalThis.fetch;

  globalThis.fetch = (async (input: RequestInfo | URL) => {
    const url = new URL(typeof input === "string" || input instanceof URL ? input : input.url);
    if (url.hostname === "origin.example.com") {
      throw new Error("origin unavailable");
    }
    if (url.pathname === "/maintenance.html") {
      return new Response("<!doctype html><h1>Technical maintenance</h1><p>Poliglot AI</p>", {
        headers: { "content-type": "text/html; charset=utf-8" },
      });
    }
    return new Response("unexpected fetch", { status: 418 });
  }) as typeof fetch;

  try {
    const response = await worker.default.fetch(new Request("https://poliglotai.ru/app/"), {
      ORIGIN_BASE_URL: "https://origin.example.com",
      MAINTENANCE_PAGE_URL: "https://poliglotai.ru/maintenance.html",
    });

    expect(response.status).toBe(503);
    expect(response.headers.get("content-type")).toContain("text/html");
    expect(await response.text()).toContain("Poliglot AI");
  } finally {
    globalThis.fetch = originalFetch;
  }
});
