const DEFAULT_ORIGIN_BASE_URL = "https://api.neriva.ru";
const DEFAULT_STATIC_BASE_URL = "https://neriva.pages.dev";
const DEFAULT_MAINTENANCE_PATH = "/maintenance.html";
const CANONICAL_PUBLIC_HOST = "neriva.ru";
const LEGACY_PUBLIC_HOSTS = new Set(["poliglotai.ru", "www.poliglotai.ru", "poliglotai.online", "www.poliglotai.online"]);

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

const fallbackStatuses = new Set([520, 521, 522, 523, 524, 525, 526, 530, 502, 503, 504]);

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

export function resolveStaticUrl(requestUrl, staticBaseUrl = DEFAULT_STATIC_BASE_URL) {
  const url = new URL(requestUrl);
  const staticOrigin = new URL(staticBaseUrl);
  staticOrigin.pathname = resolvePagesPath(url.pathname);
  staticOrigin.search = url.search;
  return staticOrigin.toString();
}

export function resolvePagesPath(pathname) {
  if (pathname === "/") {
    return "/poliglot-ai";
  }
  if (pathname === "/index.html") {
    return "/poliglot-ai";
  }
  if (pathname.endsWith(".html")) {
    return pathname.slice(0, -".html".length) || "/";
  }
  return pathname;
}

export function legacyRedirectUrl(requestUrl) {
  const url = new URL(requestUrl);
  if (!LEGACY_PUBLIC_HOSTS.has(url.hostname.toLowerCase())) {
    return "";
  }
  url.protocol = "https:";
  url.hostname = CANONICAL_PUBLIC_HOST;
  url.port = "";
  return url.toString();
}

function maintenanceUrlFor(request, env) {
  const configuredUrl = env?.MAINTENANCE_PAGE_URL || DEFAULT_MAINTENANCE_PATH;
  return new URL(configuredUrl, env?.STATIC_BASE_URL || DEFAULT_STATIC_BASE_URL).toString();
}

function inlineMaintenanceHtml() {
  return `<!doctype html>
<html lang="ru">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="robots" content="noindex, nofollow" />
    <title>Технические работы - NERIVA</title>
  </head>
  <body>
    <main>
      <h1>Технические работы</h1>
      <p>NERIVA временно недоступен. Попробуйте обновить страницу через пару минут.</p>
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
  const originalHost = new URL(request.url).host;
  headers.set("host", env?.ORIGIN_HOST || new URL(originUrl).host);
  headers.set("x-forwarded-host", originalHost);
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

async function serveStatic(request, env) {
  const staticUrl = resolveStaticUrl(request.url, env?.STATIC_BASE_URL || DEFAULT_STATIC_BASE_URL);
  const headers = new Headers(request.headers);
  headers.delete("host");
  headers.set("x-forwarded-host", new URL(request.url).host);
  headers.set("x-poliglot-edge", "cloudflare-worker");

  return fetch(
    new Request(staticUrl, {
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
    const redirectUrl = legacyRedirectUrl(request.url);
    if (redirectUrl) {
      return Response.redirect(redirectUrl, 301);
    }

    if (!shouldProxyPath(url.pathname)) {
      try {
        const response = await serveStatic(request, env);
        if (fallbackStatuses.has(response.status)) {
          return maintenanceResponse(request, env, 503);
        }
        return response;
      } catch {
        return maintenanceResponse(request, env, 503);
      }
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
