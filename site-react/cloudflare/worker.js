const DEFAULT_ORIGIN_BASE_URL = "https://origin.poliglotai.ru";
const DEFAULT_MAINTENANCE_PATH = "/maintenance.html";

const exactProxyPaths = new Set([
  "/app",
  "/login",
  "/healthz",
  "/payment/success",
  "/tonapi/webhook",
  "/yookassa/webhook",
  "/rollypay/webhook/site",
  "/rollypay/webhook/bot",
]);

const prefixProxyPaths = [
  "/app/",
  "/login/",
  "/api/",
  "/payment/success/",
  "/tonapi/webhook/",
  "/yookassa/webhook/",
  "/rollypay/webhook/",
];

const fallbackStatuses = new Set([502, 503, 504]);

export function shouldProxyPath(pathname) {
  return exactProxyPaths.has(pathname) || prefixProxyPaths.some((prefix) => pathname.startsWith(prefix));
}

export function resolveOriginUrl(requestUrl, originBaseUrl = DEFAULT_ORIGIN_BASE_URL) {
  const url = new URL(requestUrl);
  const origin = new URL(originBaseUrl);
  origin.pathname = url.pathname;
  origin.search = url.search;
  return origin.toString();
}

function maintenanceUrlFor(request, env) {
  const configuredUrl = env?.MAINTENANCE_PAGE_URL || DEFAULT_MAINTENANCE_PATH;
  return new URL(configuredUrl, request.url).toString();
}

function inlineMaintenanceHtml() {
  return `<!doctype html>
<html lang="ru">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="robots" content="noindex, nofollow" />
    <title>Технические работы - Poliglot AI</title>
  </head>
  <body>
    <main>
      <h1>Технические работы</h1>
      <p>Poliglot AI временно недоступен. Попробуйте обновить страницу через пару минут.</p>
      <p>Maintenance is in progress. Please retry in a few minutes.</p>
      <a href="/poliglot-ai.html">На главную</a>
    </main>
  </body>
</html>`;
}

export async function maintenanceResponse(request, env, status = 503) {
  try {
    const response = await fetch(maintenanceUrlFor(request, env), {
      cf: { cacheTtl: 60, cacheEverything: true },
    });

    if (response.ok) {
      const headers = new Headers(response.headers);
      headers.set("content-type", headers.get("content-type") || "text/html; charset=utf-8");
      headers.set("cache-control", "no-cache, no-store, must-revalidate");
      return new Response(await response.text(), { status, headers });
    }
  } catch {
    // Fall through to the inline backup so origin outages never produce an empty body.
  }

  return new Response(inlineMaintenanceHtml(), {
    status,
    headers: {
      "content-type": "text/html; charset=utf-8",
      "cache-control": "no-cache, no-store, must-revalidate",
    },
  });
}

async function proxyToOrigin(request, env) {
  const originUrl = resolveOriginUrl(request.url, env?.ORIGIN_BASE_URL || DEFAULT_ORIGIN_BASE_URL);
  const headers = new Headers(request.headers);
  headers.set("x-forwarded-host", new URL(request.url).host);
  headers.set("x-poliglot-edge", "cloudflare-worker");

  return fetch(
    new Request(originUrl, {
      method: request.method,
      headers,
      body: request.method === "GET" || request.method === "HEAD" ? undefined : request.body,
      redirect: "manual",
    }),
  );
}

export default {
  async fetch(request, env) {
    const url = new URL(request.url);

    if (!shouldProxyPath(url.pathname)) {
      return fetch(request);
    }

    try {
      const response = await proxyToOrigin(request, env);
      if (fallbackStatuses.has(response.status)) {
        return maintenanceResponse(request, env, 503);
      }
      return response;
    } catch {
      return maintenanceResponse(request, env, 503);
    }
  },
};
