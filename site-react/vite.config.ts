import path from "node:path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  base: "/",
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
  build: {
    outDir: "../Сайт полиглота для бота",
    emptyOutDir: false,
    sourcemap: false,
    chunkSizeWarningLimit: 1600,
    rollupOptions: {
      input: {
        main: path.resolve(__dirname, "poliglot-ai.html"),
        privacy: path.resolve(__dirname, "privacy.html"),
        terms: path.resolve(__dirname, "terms.html"),
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
