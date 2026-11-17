import path from "node:path";
import fs from "node:fs";
import { defineConfig, type PreviewServer, type ViteDevServer } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const defaultOutDir = "../Сайт полиглота для бота";
const outDir = process.env.PUBLIC_SITE_OUT_DIR || defaultOutDir;
const emptyOutDir = process.env.PUBLIC_SITE_EMPTY_OUT_DIR
  ? process.env.PUBLIC_SITE_EMPTY_OUT_DIR === "true"
  : false;

function productPreview() {
  const app = path.resolve(__dirname, "../web");
  const types: Record<string, string> = { ".js": "text/javascript", ".css": "text/css", ".png": "image/png", ".jpg": "image/jpeg", ".webp": "image/webp", ".svg": "image/svg+xml" };
  const previewHtml = () => {
    const html = fs.readFileSync(path.join(app, "index.html"), "utf8");
    const assets = html.match(/<(?:script|link)\b[^>]*(?:src|href)="\/app\/assets\/[\w-]+\.(?:js|css)"[^>]*>(?:<\/script>)?/g);
    if (!assets?.some((tag) => tag.startsWith("<script"))) throw new Error("Build web-react before building the product preview.");
    return `<!doctype html><html lang="ru" data-theme="dark"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="robots" content="noindex"><meta http-equiv="Content-Security-Policy" content="connect-src 'none'; form-action 'none'; object-src 'none'"><title>NERIVA — Product preview</title>${assets.join("\n")}</head><body><div id="root"></div></body></html>`;
  };
  const serve = (server: PreviewServer | ViteDevServer) => {
    server.middlewares.use((req, res, next) => {
      const url = new URL(req.url || "/", "http://localhost");
      if (url.pathname === "/product-preview.html") {
        res.setHeader("Content-Type", "text/html; charset=utf-8");
        res.end(previewHtml());
        return;
      }
      const file = /^\/app\/assets\/([\w.-]+)$/.exec(url.pathname)?.[1];
      const type = file && types[path.extname(file)];
      if (!type || !file) return next();
      const asset = path.join(app, "assets", file);
      if (!fs.existsSync(asset)) return next();
      res.setHeader("Content-Type", type);
      res.end(fs.readFileSync(asset));
    });
  };
  return {
    name: "neriva-product-preview",
    configureServer: serve,
    configurePreviewServer: serve,
    writeBundle() {
      fs.writeFileSync(path.resolve(__dirname, outDir, "product-preview.html"), previewHtml());
    },
  };
}

export default defineConfig({
  base: "/",
  plugins: [react(), tailwindcss(), productPreview()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
  build: {
    outDir,
    emptyOutDir,
    sourcemap: false,
    chunkSizeWarningLimit: 1600,
    rollupOptions: {
      input: {
        main: path.resolve(__dirname, "neriva.html"),
        privacy: path.resolve(__dirname, "privacy.html"),
        terms: path.resolve(__dirname, "terms.html"),
        agreement: path.resolve(__dirname, "agreement.html"),
        consent: path.resolve(__dirname, "consent.html"),
        notFound: path.resolve(__dirname, "404.html"),
        maintenance: path.resolve(__dirname, "maintenance.html"),
      },
      output: {
        entryFileNames: "assets/site-react/[name]-[hash].js",
        chunkFileNames: "assets/site-react/[name]-[hash].js",
        assetFileNames: "assets/site-react/[name]-[hash][extname]",
        manualChunks: {
          react: ["react", "react-dom"],
          motion: ["framer-motion"],
        },
      },
    },
  },
});
