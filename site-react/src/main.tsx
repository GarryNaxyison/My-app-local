import React from "react";
import { createRoot } from "react-dom/client";
import { ErrorPageApp, type ErrorPageKind } from "./ErrorPageApp";
import { PublicSiteApp } from "./PublicSiteApp";
import { setupLandingSeoMetadata } from "./seoMetadata";
import "./styles.css";
import "./englishSparkLanding.css";
import "./errorPages.css";

const rootElement = document.getElementById("root")!;
const pageKind = rootElement.dataset.pageKind as ErrorPageKind | undefined;

if (pageKind === "not-found" || pageKind === "maintenance") {
  createRoot(rootElement).render(
    <React.StrictMode>
      <ErrorPageApp kind={pageKind} />
    </React.StrictMode>,
  );
} else {
  createRoot(rootElement).render(
    <React.StrictMode>
      <PublicSiteApp />
    </React.StrictMode>,
  );

  setupLandingSeoMetadata();
}
