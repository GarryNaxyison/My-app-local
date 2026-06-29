import React from "react";
import { createRoot } from "react-dom/client";
import { PublicSiteApp } from "./PublicSiteApp";
import { setupLandingSeoMetadata } from "./seoMetadata";
import "./styles.css";
import "./englishSparkLanding.css";

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <PublicSiteApp />
  </React.StrictMode>,
);

setupLandingSeoMetadata();
