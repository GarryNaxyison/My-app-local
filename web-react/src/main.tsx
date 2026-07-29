import React from "react";
import { createRoot } from "react-dom/client";
import { MotionConfig } from "framer-motion";
import { App } from "./App";
import "./styles/app.css";
// Платформенный слой (ПК/мобайл разделение, reduced-motion) — после базовых стилей.
import "./styles/platform.css";

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    {/* Уважаем prefers-reduced-motion во всех framer-motion анимациях. */}
    <MotionConfig reducedMotion="user">
      <App />
    </MotionConfig>
  </React.StrictMode>
);
