import { expect, test } from "@playwright/test";

test("Cloudflare Worker fallback path matcher only proxies dynamic application paths", async () => {
  const worker = await import("../cloudflare/worker.js");

  expect(worker.shouldProxyPath("/app")).toBe(true);
  expect(worker.shouldProxyPath("/app/lesson")).toBe(true);
  expect(worker.shouldProxyPath("/__codex_deploy_upload/healthz")).toBe(true);
  expect(worker.shouldProxyPath("/__codex_deploy_upload/aibot-linux-amd64")).toBe(true);
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

test("Cloudflare Worker redirects legacy public domains to Neriva with path and query preserved", async () => {
  const worker = await import("../cloudflare/worker.js");

  for (const host of ["poliglotai.ru", "www.poliglotai.ru", "poliglotai.online", "www.poliglotai.online"]) {
    const response = await worker.default.fetch(new Request(`https://${host}/en/ai-english-tutor.html?lang=en&utm=test`), {});

    expect(response.status).toBe(301);
    expect(response.headers.get("location")).toBe("https://neriva.ru/en/ai-english-tutor.html?lang=en&utm=test");
  }
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
      return new Response("<!doctype html><h1>Technical maintenance</h1><p>NERIVA</p>", {
        headers: { "content-type": "text/html; charset=utf-8" },
      });
    }
    return new Response("unexpected fetch", { status: 418 });
  }) as typeof fetch;

  try {
    const response = await worker.default.fetch(new Request("https://neriva.ru/app/"), {
      ORIGIN_BASE_URL: "https://origin.example.com",
      MAINTENANCE_PAGE_URL: "https://neriva.ru/maintenance.html",
    });

    expect(response.status).toBe(503);
    expect(response.headers.get("content-type")).toContain("text/html");
    expect(await response.text()).toContain("NERIVA");
  } finally {
    globalThis.fetch = originalFetch;
  }
});

test("Cloudflare Worker serves static pages from the Pages origin", async () => {
  const worker = await import("../cloudflare/worker.js");
  const originalFetch = globalThis.fetch;
  const requestedUrls: string[] = [];

  globalThis.fetch = (async (input: RequestInfo | URL) => {
    const requestUrl = typeof input === "string" || input instanceof URL ? input.toString() : input.url;
    requestedUrls.push(requestUrl);
    return new Response("<!doctype html><h1>SEO page</h1>", {
      headers: { "content-type": "text/html; charset=utf-8" },
    });
  }) as typeof fetch;

  try {
    const response = await worker.default.fetch(new Request("https://neriva.ru/en/ai-english-tutor.html"), {
      STATIC_BASE_URL: "https://neriva.pages.dev",
    });

    expect(response.status).toBe(200);
    expect(await response.text()).toContain("SEO page");

    const homeResponse = await worker.default.fetch(new Request("https://neriva.ru/"), {
      STATIC_BASE_URL: "https://neriva.pages.dev",
    });

    expect(homeResponse.status).toBe(200);
    expect(await homeResponse.text()).toContain("SEO page");
    expect(requestedUrls).toEqual([
      "https://neriva.pages.dev/en/ai-english-tutor",
      "https://neriva.pages.dev/poliglot-ai",
    ]);
  } finally {
    globalThis.fetch = originalFetch;
  }
});
