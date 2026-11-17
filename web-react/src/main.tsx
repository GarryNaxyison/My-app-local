import React from "react";
import { createRoot } from "react-dom/client";
import { MotionConfig } from "framer-motion";
import { App, ProductPreview } from "./App";
import { PwaInstall } from "./components/PwaInstall";
import "./styles/app.css";
// Платформенный слой (ПК/мобайл разделение, reduced-motion) — после базовых стилей.
import "./styles/platform.css";
import "./styles/navigation.css";
import "./styles/learning-workspace.css";
import "./styles/brand.css";

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    {/* Уважаем prefers-reduced-motion во всех framer-motion анимациях. */}
    <MotionConfig reducedMotion="user">
      {location.pathname === "/product-preview.html" ? <ProductPreview /> : <><App /><PwaInstall /></>}
    </MotionConfig>
  </React.StrictMode>
);
