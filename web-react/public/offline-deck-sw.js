const CACHE_NAME = "poliglot-v2-static-shell-20260721";
const APP_SHELL = ["/app/", "/app/index.html", "/app/assets/brand-logo-mini.png"];

function isStaticAsset(request, url) {
  return request.method === "GET"
    && (url.pathname.startsWith("/app/assets/") || /\.(?:css|js|mjs|png|jpe?g|webp|svg|woff2?)$/i.test(url.pathname));
}

async function cacheResponse(request, response) {
  if (!response || !response.ok) return response;
  const cache = await caches.open(CACHE_NAME);
  await cache.put(request, response.clone());
  return response;
}

async function fetchAndCache(request) {
  return cacheResponse(request, await fetch(request));
}

async function precacheAppShell() {
  const cache = await caches.open(CACHE_NAME);
  await cache.addAll(APP_SHELL);
  const response = await fetch("/app/", { cache: "reload" });
  if (!response.ok) return;
  await cache.put("/app/", response.clone());
  const html = await response.text();
  const assets = Array.from(html.matchAll(/(?:src|href)=["']([^"']+)["']/g), (match) => match[1])
    .map((asset) => new URL(asset, self.location.origin))
    .filter((asset) => asset.origin === self.location.origin && (asset.pathname.startsWith("/app/") || asset.pathname.startsWith("/app/assets/")))
    .map((asset) => asset.href);
  await Promise.all(assets.map((asset) => cache.add(asset).catch(() => undefined)));
}

self.addEventListener("install", (event) => {
  event.waitUntil(precacheAppShell().catch(() => undefined));
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(keys.filter((key) => key !== CACHE_NAME && key.startsWith("poliglot-v2-")).map((key) => caches.delete(key))),
    ),
  );
  self.clients.claim();
});

self.addEventListener("fetch", (event) => {
  const request = event.request;
  if (request.method !== "GET") return;
  const url = new URL(request.url);
  if (url.pathname.startsWith("/api/")) return;
  if (!url.pathname.startsWith("/app/") && !url.pathname.startsWith("/app/assets/")) return;
  if (url.pathname === "/app/offline-deck-sw.js" || url.pathname === "/app/manifest.webmanifest") return;

  if (request.mode === "navigate") {
    event.respondWith(fetchAndCache(request).catch(() => caches.match("/app/").then((cached) => cached || Response.error())));
    return;
  }

  if (isStaticAsset(request, url)) {
    event.respondWith(caches.match(request).then((cached) => cached || fetchAndCache(request)));
  }
});

self.addEventListener("message", (event) => {
  if (event.data && event.data.type === "SKIP_WAITING") self.skipWaiting();
});
