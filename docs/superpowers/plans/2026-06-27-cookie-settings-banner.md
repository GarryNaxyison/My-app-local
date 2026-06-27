# Cookie Settings Banner Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a real `Настроить` cookie-settings flow to the Poliglot AI public-site cookie banner.

**Architecture:** Keep `CookieConsentBanner` self-contained in `site-react/src/PublicSiteApp.tsx`, with small helper functions for consent parsing and serialization. Extend the existing CSS block in `site-react/src/styles.css` so the banner keeps its current product visual language while supporting an inline settings panel. Cover the behavior through focused Playwright tests in `site-react/e2e/public-site.spec.ts`.

**Tech Stack:** React 19, TypeScript, Vite, CSS, Playwright.

---

## File Structure

- Modify `site-react/e2e/public-site.spec.ts`: add failing behavioral coverage for the settings flow and legacy storage compatibility.
- Modify `site-react/src/PublicSiteApp.tsx`: add consent helpers, expanded settings state, optional category switch, and save handlers.
- Modify `site-react/src/styles.css`: style the inline settings rows, switch, and third action without changing the public-site design system.
- No backend files change.

## Task 1: Write Red Playwright Coverage

**Files:**
- Modify: `site-react/e2e/public-site.spec.ts`
- Test: `site-react/e2e/public-site.spec.ts`

- [ ] **Step 1: Add the failing settings-flow test**

Add this test after `landing shows Russian legal links and centers the cookie banner on desktop`:

```ts
test("cookie banner lets users configure optional cookies", async ({ page }) => {
  test.setTimeout(60_000);

  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  await page.evaluate(() => localStorage.removeItem("poliglot-cookie-consent"));
  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });

  const banner = page.locator(".cookie-consent-banner");
  await expect(banner).toBeVisible();
  await expect(banner.getByRole("button", { name: "Настроить" })).toBeVisible();

  await banner.getByRole("button", { name: "Настроить" }).click();
  await expect(banner.locator(".cookie-consent-banner__settings")).toBeVisible();
  await expect(banner).toContainText("Необходимые");
  await expect(banner).toContainText("Аналитические/маркетинговые");
  await expect.poll(() => page.evaluate(() => localStorage.getItem("poliglot-cookie-consent"))).toBeNull();

  const optionalSwitch = banner.getByRole("switch", { name: "Аналитические/маркетинговые" });
  await expect(optionalSwitch).toHaveAttribute("aria-checked", "false");
  await optionalSwitch.click();
  await expect(optionalSwitch).toHaveAttribute("aria-checked", "true");

  await banner.getByRole("button", { name: "Сохранить выбор" }).click();
  await expect(banner).toHaveCount(0);

  const storedConsent = await page.evaluate(() => localStorage.getItem("poliglot-cookie-consent"));
  expect(JSON.parse(storedConsent || "{}")).toEqual({
    version: 1,
    necessary: true,
    analyticsMarketing: true,
  });

  await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
  await expect(page.locator(".cookie-consent-banner")).toHaveCount(0);
});
```

- [ ] **Step 2: Add the legacy-storage compatibility test**

Add this test near the settings-flow test:

```ts
test("cookie banner respects legacy saved consent values", async ({ page }) => {
  for (const legacyValue of ["necessary", "accepted"]) {
    await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
    await page.evaluate((value) => localStorage.setItem("poliglot-cookie-consent", value), legacyValue);
    await page.goto("/poliglot-ai.html?lang=ru", { waitUntil: "domcontentloaded" });
    await expect(page.locator(".cookie-consent-banner")).toHaveCount(0);
  }
});
```

- [ ] **Step 3: Strengthen the existing opt-in test**

In `legal document package exposes operator details and cookie opt-in`, before clicking accept, add:

```ts
await expect(banner.getByRole("button", { name: "Настроить" })).toBeVisible();
await expect(banner.getByRole("button", { name: "Только необходимые" })).toBeVisible();
await expect(banner.getByRole("button", { name: "Принять все cookie" })).toBeVisible();
```

Change the accept click to target the exact button:

```ts
await banner.getByRole("button", { name: "Принять все cookie" }).click();
```

- [ ] **Step 4: Run RED**

Run:

```bash
npm --prefix site-react run e2e -- e2e/public-site.spec.ts -g "cookie"
```

Expected before implementation: FAIL because `Настроить`, the settings panel, and the switch do not exist yet.

## Task 2: Implement Consent State and Settings UI

**Files:**
- Modify: `site-react/src/PublicSiteApp.tsx`
- Test: `site-react/e2e/public-site.spec.ts`

- [ ] **Step 1: Replace the current `CookieConsentBanner` implementation**

Replace `CookieConsentBanner` with this implementation:

```tsx
type CookieConsentChoice = {
  version: 1;
  necessary: true;
  analyticsMarketing: boolean;
};

const cookieConsentStorageKey = "poliglot-cookie-consent";

function hasSavedCookieConsent(value: string | null) {
  if (!value) return false;
  if (value === "necessary" || value === "accepted") return true;

  try {
    const parsed = JSON.parse(value) as Partial<CookieConsentChoice>;
    return parsed.version === 1 && parsed.necessary === true && typeof parsed.analyticsMarketing === "boolean";
  } catch {
    return false;
  }
}

function serializeCookieConsent(analyticsMarketing: boolean) {
  const choice: CookieConsentChoice = {
    version: 1,
    necessary: true,
    analyticsMarketing,
  };

  return JSON.stringify(choice);
}

function CookieConsentBanner() {
  const [choice, setChoice] = useState<string | null>(() => {
    try {
      return localStorage.getItem(cookieConsentStorageKey);
    } catch {
      return null;
    }
  });
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [analyticsMarketing, setAnalyticsMarketing] = useState(false);

  const saveChoice = (value: string) => {
    try {
      localStorage.setItem(cookieConsentStorageKey, value);
    } catch {
      // Ignore storage failures; the banner can still close for this session.
    }
    setChoice(value);
  };

  const saveNecessary = () => saveChoice("necessary");
  const saveAccepted = () => saveChoice(serializeCookieConsent(true));
  const saveSelected = () => saveChoice(serializeCookieConsent(analyticsMarketing));

  if (hasSavedCookieConsent(choice)) return null;

  return (
    <section className="cookie-consent-banner" aria-label="Cookie consent">
      <div className="cookie-consent-banner__copy">
        <strong data-legal-cookie="title">Cookie и технические данные</strong>
        <p data-legal-cookie="body">Мы используем необходимые cookie и локальное хранилище для входа, языка, темы, безопасности, сохранения согласий и корректной работы сайта. Необязательные cookie применяются только после согласия.</p>
        <nav>
          <a href="/privacy.html" data-legal-cookie="privacy">Политика</a>
          <a href="/consent.html" data-legal-cookie="consent">Согласие</a>
        </nav>

        {settingsOpen && (
          <div className="cookie-consent-banner__settings" aria-label="Настройки cookie">
            <div className="cookie-consent-banner__setting-row">
              <div>
                <strong>Необходимые</strong>
                <span>Всегда активны для входа, безопасности, языка и сохранения согласий.</span>
              </div>
              <span className="cookie-consent-banner__required">Всегда активны</span>
            </div>

            <div className="cookie-consent-banner__setting-row">
              <div>
                <strong>Аналитические/маркетинговые</strong>
                <span>Помогают понимать работу сайта и улучшать продвижение, если такие инструменты включены.</span>
              </div>
              <button
                type="button"
                className="cookie-consent-banner__switch"
                role="switch"
                aria-checked={analyticsMarketing}
                aria-label="Аналитические/маркетинговые"
                onClick={() => setAnalyticsMarketing((value) => !value)}
              >
                <span />
              </button>
            </div>
          </div>
        )}
      </div>

      <div className="cookie-consent-banner__actions">
        <button type="button" onClick={saveNecessary} data-legal-cookie="necessary">Только необходимые</button>
        <button type="button" onClick={() => setSettingsOpen((value) => !value)} data-legal-cookie="settings">Настроить</button>
        {settingsOpen && <button type="button" onClick={saveSelected} data-legal-cookie="save-selected">Сохранить выбор</button>}
        <button type="button" onClick={saveAccepted} data-legal-cookie="accept">Принять все cookie</button>
      </div>
    </section>
  );
}
```

- [ ] **Step 2: Run focused tests**

Run:

```bash
npm --prefix site-react run e2e -- e2e/public-site.spec.ts -g "cookie"
```

Expected after implementation but before CSS polish: behavioral tests should pass or expose only selector/layout issues.

## Task 3: Style Inline Settings

**Files:**
- Modify: `site-react/src/styles.css`
- Test: `site-react/e2e/public-site.spec.ts`

- [ ] **Step 1: Extend the cookie banner CSS**

Add these rules near the existing `.cookie-consent-banner` block:

```css
.cookie-consent-banner {
  grid-template-columns: minmax(0, 1fr);
}

.cookie-consent-banner__copy {
  min-width: 0;
}

.cookie-consent-banner__settings {
  display: grid;
  gap: 10px;
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid var(--line);
}

.cookie-consent-banner__setting-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
}

.cookie-consent-banner__setting-row strong {
  margin: 0 0 3px;
  font-size: 0.92rem;
}

.cookie-consent-banner__setting-row span {
  display: block;
  color: var(--muted);
  font-size: 0.82rem;
  line-height: 1.4;
}

.cookie-consent-banner__required {
  max-width: 132px;
  border: 1px solid var(--line);
  border-radius: 999px;
  padding: 7px 10px;
  text-align: center;
  font-weight: 850;
}

.cookie-consent-banner__switch {
  position: relative;
  width: 48px;
  min-height: 28px;
  border: 0;
  border-radius: 999px;
  padding: 3px;
  background: color-mix(in srgb, var(--muted) 35%, var(--line));
}

.cookie-consent-banner__switch span {
  display: block;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: var(--card);
  transition: transform 160ms ease;
}

.cookie-consent-banner__switch[aria-checked="true"] {
  background: linear-gradient(135deg, var(--gold), #8ff1d0);
}

.cookie-consent-banner__switch[aria-checked="true"] span {
  transform: translateX(20px);
}
```

Change the primary-action selector from:

```css
.cookie-consent-banner button:last-child {
```

to:

```css
.cookie-consent-banner button[data-legal-cookie="accept"] {
```

Inside the mobile media query, add:

```css
.cookie-consent-banner__setting-row {
  grid-template-columns: 1fr;
}

.cookie-consent-banner__required,
.cookie-consent-banner__switch {
  justify-self: start;
}
```

The one-column banner grid is required because three default actions and four expanded actions make the old `minmax(0, 1fr) auto` layout let the actions column consume max-content width and collapse the copy/settings column.

- [ ] **Step 2: Run the focused Playwright tests**

Run:

```bash
npm --prefix site-react run e2e -- e2e/public-site.spec.ts -g "cookie"
```

Expected: PASS for desktop and mobile projects.

## Task 4: Final Verification, Browser Check, Commit, Push

**Files:**
- Verify: `site-react/src/PublicSiteApp.tsx`
- Verify: `site-react/src/styles.css`
- Verify: `site-react/e2e/public-site.spec.ts`
- Verify: `docs/superpowers/plans/2026-06-27-cookie-settings-banner.md`

- [ ] **Step 1: Run build**

Run:

```bash
npm --prefix site-react run build
```

Expected: exit 0.

- [ ] **Step 2: Run focused e2e**

Run:

```bash
npm --prefix site-react run e2e -- e2e/public-site.spec.ts -g "cookie"
```

Expected: all selected tests pass in desktop and mobile projects.

- [ ] **Step 3: Verify in browser**

Open `http://127.0.0.1:4175/poliglot-ai.html?lang=ru` with Playwright/browser MCP after the preview server is running. Check:

- banner appears when `poliglot-cookie-consent` is absent;
- `Настроить` opens settings;
- switch toggles visibly;
- saving hides the banner;
- mobile viewport keeps all actions visible.

- [ ] **Step 4: Stage only owned files**

Run:

```bash
git add docs/superpowers/plans/2026-06-27-cookie-settings-banner.md site-react/src/PublicSiteApp.tsx site-react/src/styles.css site-react/e2e/public-site.spec.ts
git status --short
```

Expected staged files only include the plan and the three implementation files. Existing unrelated dirty files remain unstaged.

- [ ] **Step 5: Commit and push**

Run:

```bash
git commit -m "Add cookie settings banner"
git push
```

Expected: commit succeeds and pushes to the configured GitHub remote.

## Self-Review

- Spec coverage: the plan covers the `Настроить` action, inline settings panel, optional category switch, JSON storage, legacy value compatibility, quick actions, mobile layout, Playwright tests, browser verification, commit, and push.
- Placeholder scan: the plan contains no `TBD`, `TODO`, or unspecified implementation steps.
- Type consistency: `CookieConsentChoice`, `analyticsMarketing`, `cookieConsentStorageKey`, and `serializeCookieConsent` are named consistently across tasks.
