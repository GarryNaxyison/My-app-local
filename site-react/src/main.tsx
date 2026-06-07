import React from "react";
import { createRoot } from "react-dom/client";
import { PublicSiteApp } from "./PublicSiteApp";
import "./styles.css";

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <PublicSiteApp />
  </React.StrictMode>,
);
