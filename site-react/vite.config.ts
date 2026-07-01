import path from "node:path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const defaultOutDir = "../Сайт полиглота для бота";
const outDir = process.env.PUBLIC_SITE_OUT_DIR || defaultOutDir;
const emptyOutDir = process.env.PUBLIC_SITE_EMPTY_OUT_DIR
  ? process.env.PUBLIC_SITE_EMPTY_OUT_DIR === "true"
  : false;

export default defineConfig({
  base: "/",
  plugins: [react(), tailwindcss()],
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
        main: path.resolve(__dirname, "poliglot-ai.html"),
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
          three: ["three"],
        },
      },
    },
  },
});
