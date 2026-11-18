import { useEffect, useMemo, useRef, useState } from "react";
import { ArrowLeft, Grid2X2, Search } from "lucide-react";
import type { NavItem, UserProfile, ViewId } from "../lib/types";

type LearningNavigationProps = {
  nav: NavItem[];
  activeView: ViewId;
  setView: (view: ViewId) => void;
  theme: "light" | "dark";
  copy: (key: string, fallback: string) => string;
  platform: "desktop" | "mobile";
  navigationLayout: NonNullable<UserProfile["navigation_layout"]>;
  accountKey: string;
};

const groups = [
  { id: "learn", key: "learning_group", label: "Learning" },
  { id: "words", key: "words_group", label: "Words" },
  { id: "growth", key: "growth_group", label: "Growth" },
  { id: "account", key: "account_group", label: "Account" },
] as const;
const defaultPinned: ViewId[] = ["home", "tutor", "practice", "words"];
const assetVersion = "flux2-20260526-pronunciation-phrases-offline-roleplay";

/** Reuse the existing generated artwork; changing navigation never migrates assets. */
function navigationAsset(view: ViewId, theme: "light" | "dark") {
  if (view === "tutor") return "/app/assets/neriva-symbol.png";
  const slug = view === "dashboard" ? "progress" : view;
  return `/app/assets/icon-${slug}-${theme}.png?v=${assetVersion}`;
}

function storedOrder(key: string): ViewId[] {
  try {
    const value: unknown = JSON.parse(localStorage.getItem(key) || "[]");
    return Array.isArray(value) ? value.filter((id): id is ViewId => typeof id === "string") : [];
  } catch {
    return [];
  }
}

function orderItems(nav: NavItem[], ids: ViewId[]) {
  return [...new Set([...ids, ...nav.map((item) => item.id)])]
    .map((id) => nav.find((item) => item.id === id))
    .filter((item): item is NavItem => Boolean(item));
}

export function LearningNavigation({ nav, activeView, setView, theme, copy, platform, navigationLayout, accountKey }: LearningNavigationProps) {
  const [catalogOpen, setCatalogOpen] = useState(false);
  const [query, setQuery] = useState("");
  const moreButton = useRef<HTMLButtonElement>(null);
  const catalogHeading = useRef<HTMLHeadingElement>(null);
  const desktopItems = useMemo(
    () => orderItems(nav, navigationLayout.function_ribbon?.length ? navigationLayout.function_ribbon : storedOrder(`poliglot-function-ribbon-v2:${accountKey || "guest"}`)),
    [nav, navigationLayout.function_ribbon, accountKey],
  );
  const pinnedItems = useMemo(() => {
    const saved = navigationLayout.mobile_rail?.length
      ? navigationLayout.mobile_rail
      : navigationLayout.mobile_pinned?.length
        ? navigationLayout.mobile_pinned
        : storedOrder(`poliglot-mobile-nav-rail-v2:${accountKey || "guest"}`);
    const legacy = saved.length ? saved : storedOrder(`poliglot-mobile-nav-v2:${accountKey || "guest"}`);
    return orderItems(nav, [...legacy, ...defaultPinned]).slice(0, 4);
  }, [nav, navigationLayout.mobile_rail, navigationLayout.mobile_pinned, accountKey]);
  const catalogItems = useMemo(
    () => orderItems(nav, navigationLayout.mobile_more || []).filter((item) => `${item.label} ${item.description}`.toLocaleLowerCase().includes(query.trim().toLocaleLowerCase())),
    [nav, navigationLayout.mobile_more, query],
  );

  useEffect(() => {
    setCatalogOpen(false);
    setQuery("");
  }, [activeView, platform, accountKey]);

  useEffect(() => {
    if (catalogOpen) catalogHeading.current?.focus();
  }, [catalogOpen]);

  const closeCatalog = () => {
    setCatalogOpen(false);
    setQuery("");
    moreButton.current?.focus();
  };
  const chooseView = (view: ViewId) => {
    setCatalogOpen(false);
    setQuery("");
    setView(view);
  };
  const artwork = (item: NavItem) => {
    if (item.id === "games") {
      const Icon = item.icon;
      return <span className="learning-nav-art function-chip__art learning-nav-art--games" aria-hidden="true"><Icon size={18} /></span>;
    }
    return <span className="learning-nav-art function-chip__art" aria-hidden="true" style={{ backgroundImage: `url("${navigationAsset(item.id, theme)}")` }} />;
  };

  if (platform === "desktop") {
    return (
      <nav className="learning-sidebar" aria-label={copy("open_menu", "Open menu")}>
        <div className="learning-sidebar__brand">
          <img src="/app/assets/neriva-wordmark.png" alt="NERIVA" width="164" height="30" />
          <small>{copy("learning_group", "Learning")}</small>
        </div>
        {groups.map((group) => {
          const items = desktopItems.filter((item) => item.group === group.id);
          if (!items.length) return null;
          return (
            <section className="learning-sidebar__group" key={group.id} aria-labelledby={`learning-group-${group.id}`}>
              <h2 id={`learning-group-${group.id}`}>{copy(group.key, group.label)}</h2>
              {items.map((item) => (
                <button
                  className={`learning-sidebar__item function-chip${item.id === activeView ? " is-active" : ""}`}
                  key={item.id}
                  type="button"
                  data-view={item.id}
                  aria-current={item.id === activeView ? "page" : undefined}
                  aria-pressed={item.id === activeView}
                  title={item.description || undefined}
                  onClick={() => chooseView(item.id)}
                >
                  {artwork(item)}
                  <span className="learning-sidebar__label">{item.label}</span>
                </button>
              ))}
            </section>
          );
        })}
      </nav>
    );
  }

  const secondaryActive = !pinnedItems.some((item) => item.id === activeView);
  return (
    <>
      {catalogOpen ? (
        <section
          className="learning-catalog"
          id="learning-catalog"
          aria-labelledby="learning-catalog-title"
          onKeyDown={(event) => {
            if (event.key === "Escape") {
              event.preventDefault();
              closeCatalog();
            }
          }}
        >
          <header className="learning-catalog__header">
            <button type="button" className="learning-catalog__back" aria-label={copy("close_menu", "Close menu")} onClick={closeCatalog}>
              <ArrowLeft size={22} aria-hidden="true" />
            </button>
            <h1 id="learning-catalog-title" ref={catalogHeading} tabIndex={-1}>{copy("more", "More")}</h1>
          </header>
          <label className="learning-catalog__search">
            <Search size={19} aria-hidden="true" />
            <input
              type="search"
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder={copy("navigation_search", "Search sections")}
              aria-label={copy("navigation_search", "Search sections")}
            />
          </label>
          <nav aria-label={copy("more", "More")}>
            {groups.map((group) => {
              const items = catalogItems.filter((item) => item.group === group.id);
              if (!items.length) return null;
              return (
                <section className="learning-catalog__group" key={group.id} aria-labelledby={`catalog-group-${group.id}`}>
                  <h2 id={`catalog-group-${group.id}`}>{copy(group.key, group.label)}</h2>
                  <div className="learning-catalog__items">
                    {items.map((item) => (
                      <button type="button" key={item.id} data-view={item.id} aria-current={item.id === activeView ? "page" : undefined} onClick={() => chooseView(item.id)}>
                        {artwork(item)}
                        <span><strong>{item.label}</strong>{item.description ? <small>{item.description}</small> : null}</span>
                      </button>
                    ))}
                  </div>
                </section>
              );
            })}
          </nav>
          {!catalogItems.length ? <p className="learning-catalog__empty" role="status">{copy("navigation_empty", "No sections found")}</p> : null}
        </section>
      ) : null}
      <nav className="mobile-bottom-nav-v2 learning-bottom-nav" aria-label={copy("mobile_navigation", "Mobile navigation")}>
        {pinnedItems.map((item) => {
          const Icon = item.icon;
          const selected = !catalogOpen && item.id === activeView;
          return (
            <button
              type="button"
              key={item.id}
              className={selected ? "is-active" : undefined}
              data-view={item.id}
              aria-current={selected ? "page" : undefined}
              aria-pressed={selected}
              onClick={() => chooseView(item.id)}
            >
              <Icon size={21} aria-hidden="true" />
              <strong>{item.shortLabel || item.label}</strong>
            </button>
          );
        })}
        <button
          ref={moreButton}
          type="button"
          data-view="more"
          className={catalogOpen || secondaryActive ? "is-active" : undefined}
          aria-expanded={catalogOpen}
          aria-controls={catalogOpen ? "learning-catalog" : undefined}
          onClick={() => catalogOpen ? closeCatalog() : setCatalogOpen(true)}
        >
          <Grid2X2 size={21} aria-hidden="true" />
          <strong>{copy("more", "More")}</strong>
        </button>
      </nav>
    </>
  );
}
