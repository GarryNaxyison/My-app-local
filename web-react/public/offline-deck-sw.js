const CACHE_NAME = "neriva-app-static-v3";
const SHELL_PATH = "/app/";
const BRAND_ASSETS = ["neriva-symbol.png", "neriva-wordmark.png", "neriva-icon-192.png"];
const NETWORK_SHELL_PATHS = new Set([SHELL_PATH, "/app/index.html", "/app/mobile", "/app/desktop"]);

function isBuildAsset(url) {
  return /^\/app\/assets\/[\w-]+-[\w-]{8}\.(?:js|css)$/.test(url.pathname);
}

function isBrandAsset(url) {
  return BRAND_ASSETS.some((name) => url.pathname === `/app/assets/${name}`);
}

async function storeAsset(cache, request, response) {
  if (response.ok && !response.redirected) {
    await cache.put(request, response.clone()).catch(() => undefined);
  }
  return response;
}

async function cacheShell(cache, response) {
  if (!response.ok || response.redirected || !response.headers.get("content-type")?.includes("text/html")) return;
  const html = await response.clone().text();
  const entryAssets = new Set();
  for (const [, attribute] of html.matchAll(/(?:src|href)=["']([^"']+)["']/g)) {
    const url = new URL(attribute, self.location.origin);
    if (url.origin === self.location.origin && isBuildAsset(url)) entryAssets.add(url.pathname);
  }
  if (!entryAssets.size) return;

  const assetsReady = await Promise.all([...entryAssets].map(async (path) => {
    if (await cache.match(path)) return true;
    const asset = await fetch(path);
    if (!asset.ok || asset.redirected) return false;
    await cache.put(path, asset);
    return true;
  }));
  if (assetsReady.every(Boolean)) await cache.put(SHELL_PATH, response);
}

self.addEventListener("install", (event) => {
  event.waitUntil((async () => {
    const cache = await caches.open(CACHE_NAME);
    await Promise.allSettled([
      fetch(SHELL_PATH, { cache: "reload" }).then((response) => cacheShell(cache, response)),
      ...BRAND_ASSETS.map((name) => cache.add(`/app/assets/${name}`)),
    ]);
    await self.skipWaiting();
  })());
});

self.addEventListener("activate", (event) => {
  event.waitUntil((async () => {
    if (self.registration.navigationPreload) {
      await self.registration.navigationPreload.enable().catch(() => undefined);
    }
    // Drop only earlier NERIVA shell caches. Offline decks live separately.
    const keys = await caches.keys();
    await Promise.all(keys.filter((key) => key !== CACHE_NAME && /^(?:neriva-app-static-|poliglot-v2-static-shell-)/.test(key)).map((key) => caches.delete(key)));
    await self.clients.claim();
  })());
});

self.addEventListener("fetch", (event) => {
  const request = event.request;
  const url = new URL(request.url);
  if (request.method !== "GET" || url.origin !== self.location.origin || url.pathname.startsWith("/api/")) return;

  if (request.mode === "navigate" && NETWORK_SHELL_PATHS.has(url.pathname)) {
    const network = (async () => {
      const preload = await Promise.resolve(event.preloadResponse).catch(() => undefined);
      return preload || fetch(request);
    })();
    if (url.pathname === SHELL_PATH && !url.searchParams.has("shell")) {
      const offlineShell = network.then((response) => response.clone());
      event.waitUntil(offlineShell.then(async (response) => cacheShell(await caches.open(CACHE_NAME), response)).catch(() => undefined));
    }
    event.respondWith(network.catch(async () => {
      const cache = await caches.open(CACHE_NAME);
      return (await cache.match(SHELL_PATH)) || Response.error();
    }));
    return;
  }

  if (!url.search && (isBuildAsset(url) || isBrandAsset(url))) {
    event.respondWith((async () => {
      const cache = await caches.open(CACHE_NAME);
      // These same-origin URLs identify one asset, including for module requests with Origin.
      const cached = await cache.match(url.pathname, { ignoreVary: true });
      if (cached) return cached;
      return storeAsset(cache, request, await fetch(request));
    })());
  }
});

self.addEventListener("message", (event) => {
  if (event.data?.type === "SKIP_WAITING") event.waitUntil(self.skipWaiting());
});
