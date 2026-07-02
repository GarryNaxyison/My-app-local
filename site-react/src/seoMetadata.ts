import {
  canonicalOriginForLanguage,
  canonicalUrlForLanguage,
  landingLanguageCodes,
  landingSeoCopy,
  normalizeLandingLanguage,
  seoLocaleForLanguage,
  socialProfileUrls,
} from "./landingSeoContent";

type SiteI18nWindow = Window & {
  poliglotSiteI18n?: {
    currentLanguage?: () => string;
  };
};

const jsonLdScriptId = "poliglot-seo-jsonld";
const supportEmail = "support@neriva.ru";

function isLandingPage() {
  return document.documentElement.dataset.sitePage === "landing" || location.pathname === "/" || location.pathname.endsWith("/poliglot-ai.html");
}

function currentLandingLanguage() {
  const i18nLanguage = (window as SiteI18nWindow).poliglotSiteI18n?.currentLanguage?.();
  if (i18nLanguage) return normalizeLandingLanguage(i18nLanguage);

  const params = new URLSearchParams(window.location.search);
  return normalizeLandingLanguage(params.get("lang") || localStorage.getItem("poliglot_site_language") || document.documentElement.lang || "ru");
}

function upsertMetaByName(name: string, content: string) {
  let meta = document.head.querySelector<HTMLMetaElement>(`meta[name="${name}"]`);
  if (!meta) {
    meta = document.createElement("meta");
    meta.setAttribute("name", name);
    document.head.append(meta);
  }
  if (meta.getAttribute("content") !== content) {
    meta.setAttribute("content", content);
  }
}

function upsertMetaByProperty(property: string, content: string) {
  let meta = document.head.querySelector<HTMLMetaElement>(`meta[property="${property}"]`);
  if (!meta) {
    meta = document.createElement("meta");
    meta.setAttribute("property", property);
    document.head.append(meta);
  }
  if (meta.getAttribute("content") !== content) {
    meta.setAttribute("content", content);
  }
}

function upsertLink(selector: string, attributes: Record<string, string>) {
  let link = document.head.querySelector<HTMLLinkElement>(selector);
  if (!link) {
    link = document.createElement("link");
    document.head.append(link);
  }
  Object.entries(attributes).forEach(([name, value]) => {
    if (link.getAttribute(name) !== value) {
      link.setAttribute(name, value);
    }
  });
}

function absoluteAssetUrl(path: string, origin: string) {
  return new URL(path, origin).toString();
}

function socialImagePathForLanguage(language: ReturnType<typeof currentLandingLanguage>) {
  return language === "ru" ? "/assets/seo/poliglot-ai-og-ru.jpg" : "/assets/seo/poliglot-ai-og-en.jpg";
}

function buildJsonLd(language: ReturnType<typeof currentLandingLanguage>) {
  const locale = seoLocaleForLanguage(language);
  const copy = landingSeoCopy[locale];
  const canonicalUrl = canonicalUrlForLanguage(language);
  const origin = canonicalOriginForLanguage(language);
  const organizationId = `${origin}/#organization`;
  const websiteId = `${origin}/#website`;
  const applicationId = `${origin}/#software-application`;

  return {
    "@context": "https://schema.org",
    "@graph": [
      {
        "@type": "Organization",
        "@id": organizationId,
        name: "NERIVA",
        url: `${origin}/`,
        logo: absoluteAssetUrl("/assets/brand-logo-mini.png", origin),
        sameAs: [...socialProfileUrls],
        email: supportEmail,
      },
      {
        "@type": "WebSite",
        "@id": websiteId,
        name: "NERIVA",
        url: `${origin}/`,
        inLanguage: language,
        publisher: { "@id": organizationId },
      },
      {
        "@type": "SoftwareApplication",
        "@id": applicationId,
        name: "NERIVA",
        applicationCategory: "EducationalApplication",
        operatingSystem: "Web, Telegram",
        url: canonicalUrl,
        description: copy.description,
        publisher: { "@id": organizationId },
        offers: [
          { "@type": "Offer", name: "Free", price: "0", priceCurrency: "RUB" },
          { "@type": "Offer", name: "Premium", price: "300", priceCurrency: "RUB" },
          { "@type": "Offer", name: "Platinum", price: "590", priceCurrency: "RUB" },
        ],
      },
      {
        "@type": "FAQPage",
        "@id": `${canonicalUrl}#faq`,
        url: canonicalUrl,
        inLanguage: language,
        mainEntity: copy.questions.map((item) => ({
          "@type": "Question",
          name: item.question,
          acceptedAnswer: {
            "@type": "Answer",
            text: item.answer,
          },
        })),
      },
      {
        "@type": "ItemList",
        "@id": `${canonicalUrl}#comparisons`,
        name: copy.comparisonTitle,
        description: copy.comparisonIntro,
        url: `${canonicalUrl}#compare`,
        inLanguage: language,
        itemListElement: copy.comparisons.map((item, index) => ({
          "@type": "ListItem",
          position: index + 1,
          name: item.title,
          url: `${canonicalUrl}#comparison-${index + 1}`,
          item: {
            "@type": "Thing",
            name: item.title,
            description: `${item.alternativeLabel}: ${item.alternative} ${item.productLabel}: ${item.product} ${item.verdict}`,
          },
        })),
      },
    ],
  };
}

export function updateLandingSeoMetadata() {
  if (!isLandingPage()) return;

  const language = currentLandingLanguage();
  const locale = seoLocaleForLanguage(language);
  const copy = landingSeoCopy[locale];
  const canonicalUrl = canonicalUrlForLanguage(language);
  const origin = canonicalOriginForLanguage(language);
  const imageUrl = absoluteAssetUrl(socialImagePathForLanguage(language), origin);

  if (document.title !== copy.title) {
    document.title = copy.title;
  }
  upsertMetaByName("description", copy.description);
  upsertMetaByName("robots", "index, follow");

  upsertLink('link[rel="canonical"]', { rel: "canonical", href: canonicalUrl });
  landingLanguageCodes.forEach((code) => {
    upsertLink(`link[rel="alternate"][hreflang="${code}"]`, {
      rel: "alternate",
      hreflang: code,
      href: canonicalUrlForLanguage(code),
    });
  });
  upsertLink('link[rel="alternate"][hreflang="x-default"]', {
    rel: "alternate",
    hreflang: "x-default",
    href: canonicalUrlForLanguage("en"),
  });

  upsertMetaByProperty("og:title", copy.title);
  upsertMetaByProperty("og:description", copy.description);
  upsertMetaByProperty("og:type", "website");
  upsertMetaByProperty("og:url", canonicalUrl);
  upsertMetaByProperty("og:image", imageUrl);
  upsertMetaByProperty("og:image:width", "1200");
  upsertMetaByProperty("og:image:height", "630");
  upsertMetaByProperty("og:image:alt", `${copy.title} social preview`);

  upsertMetaByName("twitter:card", "summary_large_image");
  upsertMetaByName("twitter:title", copy.title);
  upsertMetaByName("twitter:description", copy.description);
  upsertMetaByName("twitter:image", imageUrl);
  upsertMetaByName("twitter:image:alt", `${copy.title} social preview`);

  let script = document.getElementById(jsonLdScriptId) as HTMLScriptElement | null;
  if (!script) {
    script = document.createElement("script");
    script.id = jsonLdScriptId;
    script.type = "application/ld+json";
    document.head.append(script);
  }
  const nextJsonLd = JSON.stringify(buildJsonLd(language));
  if (script.textContent !== nextJsonLd) {
    script.textContent = nextJsonLd;
  }
}

export function setupLandingSeoMetadata() {
  if (!isLandingPage()) return;

  const apply = () => {
    updateLandingSeoMetadata();
    [0, 160, 360, 760, 1260].forEach((delay) => {
      window.setTimeout(updateLandingSeoMetadata, delay);
    });
  };
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", apply, { once: true });
  } else {
    apply();
  }
  window.addEventListener("poliglot-language-change", apply);

  let queued = false;
  const observer = new MutationObserver(() => {
    if (queued) return;
    queued = true;
    window.setTimeout(() => {
      queued = false;
      updateLandingSeoMetadata();
    }, 0);
  });
  observer.observe(document.head, {
    attributeFilter: ["content", "href", "hreflang", "name", "property", "rel"],
    attributes: true,
    characterData: true,
    childList: true,
    subtree: true,
  });
}
