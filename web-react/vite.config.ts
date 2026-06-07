import path from "node:path";
import fs from "node:fs";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

function localAppAssets() {
  const assetsRoot = path.resolve(__dirname, "../web/assets");
  const contentTypes: Record<string, string> = {
    ".avif": "image/avif",
    ".css": "text/css; charset=utf-8",
    ".gif": "image/gif",
    ".ico": "image/x-icon",
    ".js": "text/javascript; charset=utf-8",
    ".png": "image/png",
    ".jpg": "image/jpeg",
    ".jpeg": "image/jpeg",
    ".map": "application/json; charset=utf-8",
    ".mjs": "text/javascript; charset=utf-8",
    ".webp": "image/webp",
    ".woff": "font/woff",
    ".woff2": "font/woff2",
    ".svg": "image/svg+xml",
    ".json": "application/json; charset=utf-8",
    ".wasm": "application/wasm",
  };

  const middleware = (req: any, res: any, next: () => void) => {
    if (!req.url?.startsWith("/app/assets/")) {
      next();
      return;
    }
    const pathname = new URL(req.url, "http://localhost").pathname;
    const relative = decodeURIComponent(pathname.replace(/^\/app\/assets\//, ""));
    const target = path.resolve(assetsRoot, relative);
    if (!target.startsWith(assetsRoot) || !fs.existsSync(target) || fs.statSync(target).isDirectory()) {
      res.statusCode = 404;
      res.end("not found");
      return;
    }
    res.setHeader("Content-Type", contentTypes[path.extname(target).toLowerCase()] || "application/octet-stream");
    fs.createReadStream(target).pipe(res);
  };

  return {
    name: "local-app-assets",
    configureServer(server: any) {
      server.middlewares.use(middleware);
    },
    configurePreviewServer(server: any) {
      server.middlewares.use(middleware);
    },
  };
}

export default defineConfig({
  base: "/app/",
  plugins: [localAppAssets(), react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
  server: {
    proxy: {
      "/api": "http://127.0.0.1:8080",
      "/app/assets": "http://127.0.0.1:8080",
      "/assets": "http://127.0.0.1:8080"
    }
  },
  build: {
    outDir: "../web",
    emptyOutDir: false,
    sourcemap: false,
    chunkSizeWarningLimit: 1500,
    rollupOptions: {
      output: {
        manualChunks: {
          react: ["react", "react-dom"],
          motion: ["framer-motion"],
          shaders: ["@paper-design/shaders-react"],
          canvas: ["three", "@react-three/fiber"]
        }
      }
    }
  }
});
