# NERIVA Social Links Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add branded YouTube, Instagram, and TikTok links to the public landing/footer and the web app Settings screen, then deploy the verified build to production.

**Architecture:** Keep social link data local to each frontend to avoid coupling `site-react` and `web-react` builds. Use inline SVG/React icons derived from stable brand icon sources so the app has no runtime dependency on external icon URLs. Add tests first, implement the minimum UI, verify builds/e2e, commit/push, then deploy static site and app assets.

**Tech Stack:** React 19, TypeScript, Vite, CSS, Playwright, Nx/npm scripts, Git Bash, SSH/SCP or tar-over-ssh deployment to `root@186.246.45.123`.

---

## File Structure

- Create `site-react/src/components/SocialLinks.tsx`: public-site social URL model, branded icon components, compact hero row, and footer link block.
- Modify `site-react/src/EnglishSparkLanding.tsx`: render the compact social proof row directly after `.hero-proof`.
- Modify `site-react/src/PublicSiteApp.tsx`: render social links in both English and Russian footers.
- Modify `site-react/src/styles.css`: style the landing social row, shared icon buttons, and footer social links.
- Modify `site-react/e2e/public-site.spec.ts`: add failing assertions for landing and footer social links.
- Modify `web-react/src/components/BrandIcons.tsx`: add YouTube, Instagram, and TikTok icon components beside the existing Telegram icon.
- Modify `web-react/src/App.tsx`: add a local social link model and a `Follow NERIVA` card inside `SettingsView`.
- Modify `web-react/src/styles/app.css`: style the Settings social card and icon row.
- Modify `web-react/e2e/web-smoke.spec.ts`: add failing assertions for Settings social links.

---

### Task 1: Public Site RED Test

**Files:**
- Modify: `site-react/e2e/public-site.spec.ts`

- [ ] **Step 1: Add failing public-site assertions**

In `test("landing presents the approved English spark hero product site", ...)`, after the existing `.hero-proof` checks, add:

```ts
  const heroSocialLinks = page.locator(".landing-social-proof a");
  await expect(page.locator(".landing-social-proof")).toContainText("Follow NERIVA");
  await expect(page.locator(".landing-social-proof")).toContainText("Short lessons, product updates, and learning tips.");
  await expect(heroSocialLinks).toHaveCount(3);
  await expect(heroSocialLinks.nth(0)).toHaveAttribute("href", "https://www.youtube.com/@NERIVA");
  await expect(heroSocialLinks.nth(0)).toHaveAttribute("aria-label", "Open NERIVA on YouTube");
  await expect(heroSocialLinks.nth(0)).toHaveAttribute("target", "_blank");
  await expect(heroSocialLinks.nth(0)).toHaveAttribute("rel", "noreferrer");
  await expect(heroSocialLinks.nth(1)).toHaveAttribute("href", "https://www.instagram.com/poliglotai.online/");
  await expect(heroSocialLinks.nth(1)).toHaveAttribute("aria-label", "Open NERIVA on Instagram");
  await expect(heroSocialLinks.nth(2)).toHaveAttribute("href", "https://www.tiktok.com/@poliglotai");
  await expect(heroSocialLinks.nth(2)).toHaveAttribute("aria-label", "Open NERIVA on TikTok");
```

In the same file, add this focused test after the landing hero test:

```ts
test("public footer exposes NERIVA social channels", async ({ page }) => {
  await page.goto("/poliglot-ai.html?lang=en");

  const footerSocial = page.locator(".site-footer .site-footer-social");
  await expect(footerSocial).toBeVisible();
  await expect(footerSocial).toContainText("Social");
  await expect(footerSocial.locator('a[href="https://www.youtube.com/@NERIVA"]')).toHaveAttribute("aria-label", "Open NERIVA on YouTube");
  await expect(footerSocial.locator('a[href="https://www.instagram.com/poliglotai.online/"]')).toHaveAttribute("aria-label", "Open NERIVA on Instagram");
  await expect(footerSocial.locator('a[href="https://www.tiktok.com/@poliglotai"]')).toHaveAttribute("aria-label", "Open NERIVA on TikTok");

  await page.goto("/poliglot-ai.html?lang=ru");
  const ruFooterSocial = page.locator(".site-footer .site-footer-social");
  await expect(ruFooterSocial).toContainText("Соцсети");
  await expect(ruFooterSocial.locator('a[href="https://www.youtube.com/@NERIVA"]')).toBeVisible();
  await expect(ruFooterSocial.locator('a[href="https://www.instagram.com/poliglotai.online/"]')).toBeVisible();
  await expect(ruFooterSocial.locator('a[href="https://www.tiktok.com/@poliglotai"]')).toBeVisible();
});
```

- [ ] **Step 2: Run RED test**

Run:

```bash
npm --prefix site-react run e2e -- public-site.spec.ts -g "landing presents the approved English spark hero product site|public footer exposes NERIVA social channels"
```

Expected: FAIL because `.landing-social-proof` and `.site-footer-social` do not exist yet.

- [ ] **Step 3: Commit test only if project policy allows intermediate commits**

Run:

```bash
git add site-react/e2e/public-site.spec.ts
git commit -m "test public social links"
```

Expected: commit succeeds, or skip the commit if the working tree contains unrelated staged changes from another worker.

---

### Task 2: Public Site Implementation

**Files:**
- Create: `site-react/src/components/SocialLinks.tsx`
- Modify: `site-react/src/EnglishSparkLanding.tsx`
- Modify: `site-react/src/PublicSiteApp.tsx`
- Modify: `site-react/src/styles.css`

- [ ] **Step 1: Create public social components**

Create `site-react/src/components/SocialLinks.tsx`:

```tsx
import type { ComponentType } from "react";

type SocialIconProps = {
  className?: string;
};

type SocialLink = {
  name: string;
  href: string;
  label: string;
  Icon: ComponentType<SocialIconProps>;
};

const YouTubeIcon = ({ className }: SocialIconProps) => (
  <svg className={className} fill="currentColor" role="img" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z" />
  </svg>
);

const InstagramIcon = ({ className }: SocialIconProps) => (
  <svg className={className} fill="currentColor" role="img" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    <path d="M7.0301.084c-1.2768.0602-2.1487.264-2.911.5634-.7888.3075-1.4575.72-2.1228 1.3877-.6652.6677-1.075 1.3368-1.3802 2.127-.2954.7638-.4956 1.6365-.552 2.914-.0564 1.2775-.0689 1.6882-.0626 4.947.0062 3.2586.0206 3.6671.0825 4.9473.061 1.2765.264 2.1482.5635 2.9107.308.7889.72 1.4573 1.388 2.1228.6679.6655 1.3365 1.0743 2.1285 1.38.7632.295 1.6361.4961 2.9134.552 1.2773.056 1.6884.069 4.9462.0627 3.2578-.0062 3.668-.0207 4.9478-.0814 1.28-.0607 2.147-.2652 2.9098-.5633.7889-.3086 1.4578-.72 2.1228-1.3881.665-.6682 1.0745-1.3378 1.3795-2.1284.2957-.7632.4966-1.636.552-2.9124.056-1.2809.0692-1.6898.063-4.948-.0063-3.2583-.021-3.6668-.0817-4.9465-.0607-1.2797-.264-2.1487-.5633-2.9117-.3084-.7889-.72-1.4568-1.3876-2.1228C21.2982 1.33 20.628.9208 19.8378.6165 19.074.321 18.2017.1197 16.9244.0645 15.6471.0093 15.236-.005 11.977.0014 8.718.0076 8.31.0215 7.0301.0839m.1402 21.6932c-1.17-.0509-1.8053-.2453-2.2287-.408-.5606-.216-.96-.4771-1.3819-.895-.422-.4178-.6811-.8186-.9-1.378-.1644-.4234-.3624-1.058-.4171-2.228-.0595-1.2645-.072-1.6442-.079-4.848-.007-3.2037.0053-3.583.0607-4.848.05-1.169.2456-1.805.408-2.2282.216-.5613.4762-.96.895-1.3816.4188-.4217.8184-.6814 1.3783-.9003.423-.1651 1.0575-.3614 2.227-.4171 1.2655-.06 1.6447-.072 4.848-.079 3.2033-.007 3.5835.005 4.8495.0608 1.169.0508 1.8053.2445 2.228.408.5608.216.96.4754 1.3816.895.4217.4194.6816.8176.9005 1.3787.1653.4217.3617 1.056.4169 2.2263.0602 1.2655.0739 1.645.0796 4.848.0058 3.203-.0055 3.5834-.061 4.848-.051 1.17-.245 1.8055-.408 2.2294-.216.5604-.4763.96-.8954 1.3814-.419.4215-.8181.6811-1.3783.9-.4224.1649-1.0577.3617-2.2262.4174-1.2656.0595-1.6448.072-4.8493.079-3.2045.007-3.5825-.006-4.848-.0608M16.953 5.5864A1.44 1.44 0 1 0 18.39 4.144a1.44 1.44 0 0 0-1.437 1.4424M5.8385 12.012c.0067 3.4032 2.7706 6.1557 6.173 6.1493 3.4026-.0065 6.157-2.7701 6.1506-6.1733-.0065-3.4032-2.771-6.1565-6.174-6.1498-3.403.0067-6.156 2.771-6.1496 6.1738M8 12.0077a4 4 0 1 1 4.008 3.9921A3.9996 3.9996 0 0 1 8 12.0077" />
  </svg>
);

const TikTokIcon = ({ className }: SocialIconProps) => (
  <svg className={className} fill="currentColor" role="img" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    <path d="M12.525.02c1.31-.02 2.61-.01 3.91-.02.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z" />
  </svg>
);

export const socialLinks: SocialLink[] = [
  { name: "YouTube", href: "https://www.youtube.com/@NERIVA", label: "Open NERIVA on YouTube", Icon: YouTubeIcon },
  { name: "Instagram", href: "https://www.instagram.com/poliglotai.online/", label: "Open NERIVA on Instagram", Icon: InstagramIcon },
  { name: "TikTok", href: "https://www.tiktok.com/@poliglotai", label: "Open NERIVA on TikTok", Icon: TikTokIcon },
];

export function SocialIconLinks({ className = "" }: { className?: string }) {
  return (
    <div className={`social-icon-links ${className}`.trim()} aria-label="NERIVA social channels">
      {socialLinks.map(({ name, href, label, Icon }) => (
        <a key={name} className={`social-icon-link social-icon-link--${name.toLowerCase()}`} href={href} aria-label={label} target="_blank" rel="noreferrer">
          <Icon />
          <span>{name}</span>
        </a>
      ))}
    </div>
  );
}

export function LandingSocialProof() {
  return (
    <div className="landing-social-proof">
      <div>
        <strong>Follow NERIVA</strong>
        <span>Short lessons, product updates, and learning tips.</span>
      </div>
      <SocialIconLinks />
    </div>
  );
}

export function SiteFooterSocial({ title }: { title: string }) {
  return (
    <nav className="site-footer-social" aria-label={title}>
      <strong>{title}</strong>
      <SocialIconLinks />
    </nav>
  );
}
```

- [ ] **Step 2: Render on the landing page**

In `site-react/src/EnglishSparkLanding.tsx`, add:

```tsx
import { LandingSocialProof } from "./components/SocialLinks";
```

Render after the closing `</div>` for `.hero-proof`:

```tsx
            <LandingSocialProof />
```

- [ ] **Step 3: Render in both public footers**

In `site-react/src/PublicSiteApp.tsx`, add:

```tsx
import { SiteFooterSocial } from "./components/SocialLinks";
```

Inside `SiteFooterEnglish`, place this between `Navigation` and `Documents`:

```tsx
      <SiteFooterSocial title="Social" />
```

Inside `SiteFooter`, place this between `Навигация` and `Документы`:

```tsx
      <SiteFooterSocial title="Соцсети" />
```

- [ ] **Step 4: Add public-site styles**

Append to the social/hero/footer section of `site-react/src/styles.css`:

```css
.landing-social-proof {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  max-width: 660px;
  margin-top: 14px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.08);
  padding: 12px;
  color: rgba(255, 255, 255, 0.76);
}

.landing-social-proof > div:first-child {
  display: grid;
  gap: 3px;
  min-width: min(260px, 100%);
}

.landing-social-proof strong {
  color: #fff;
  font-size: 0.95rem;
}

.landing-social-proof span {
  font-size: 0.82rem;
  line-height: 1.35;
}

.social-icon-links {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.social-icon-link {
  min-width: 44px;
  min-height: 44px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 999px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 9px;
  color: white;
  text-decoration: none;
  transition: transform 180ms ease, border-color 180ms ease, background 180ms ease;
}

.social-icon-link svg {
  width: 20px;
  height: 20px;
  flex: 0 0 auto;
}

.social-icon-link span {
  width: 1px;
  height: 1px;
  clip: rect(0 0 0 0);
  clip-path: inset(50%);
  overflow: hidden;
  position: absolute;
  white-space: nowrap;
}

.social-icon-link:hover,
.social-icon-link:focus-visible {
  border-color: rgba(255, 255, 255, 0.38);
  background: rgba(255, 255, 255, 0.12);
  transform: translateY(-1px);
  outline: none;
}

.social-icon-link--youtube {
  color: #ff0000;
}

.social-icon-link--instagram {
  color: #ff0069;
}

.social-icon-link--tiktok {
  color: #ffffff;
  text-shadow: 1px 0 #25f4ee, -1px 0 #fe2c55;
}

.site-footer-social .social-icon-links {
  margin-top: 2px;
}
```

Add light-theme overrides near existing light hero/footer overrides:

```css
html[data-site-theme="light"] .landing-social-proof {
  border-color: rgba(36, 64, 108, 0.16);
  background: rgba(255, 255, 255, 0.72);
  color: rgba(20, 34, 58, 0.72);
}

html[data-site-theme="light"] .landing-social-proof strong {
  color: #111827;
}

html[data-site-theme="light"] .landing-social-proof .social-icon-link {
  border-color: rgba(36, 64, 108, 0.16);
  background: rgba(255, 255, 255, 0.72);
}
```

- [ ] **Step 5: Run GREEN public-site test**

Run:

```bash
npm --prefix site-react run e2e -- public-site.spec.ts -g "landing presents the approved English spark hero product site|public footer exposes NERIVA social channels"
```

Expected: PASS.

- [ ] **Step 6: Run public site build**

Run:

```bash
npm --prefix site-react run build
```

Expected: exit code 0 and Vite build output.

- [ ] **Step 7: Commit public site slice**

Run:

```bash
git add site-react/src/components/SocialLinks.tsx site-react/src/EnglishSparkLanding.tsx site-react/src/PublicSiteApp.tsx site-react/src/styles.css site-react/e2e/public-site.spec.ts
git commit -m "Add public social links"
```

Expected: commit succeeds.

---

### Task 3: Web App RED Test

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add failing Settings assertions**

In `test("today plan and settings password helper copy stay localized", ...)`, after the existing `await page.goto("/app/?view=settings");`, add:

```ts
  const settingsSocial = page.locator(".settings-social-card-v2");
  await expect(settingsSocial).toBeVisible();
  await expect(settingsSocial).toContainText("Соцсети NERIVA");
  await expect(settingsSocial).toContainText("Короткие уроки, обновления и советы по обучению.");
  await expect(settingsSocial.locator('a[href="https://www.youtube.com/@NERIVA"]')).toHaveAttribute("aria-label", "Open NERIVA on YouTube");
  await expect(settingsSocial.locator('a[href="https://www.instagram.com/poliglotai.online/"]')).toHaveAttribute("aria-label", "Open NERIVA on Instagram");
  await expect(settingsSocial.locator('a[href="https://www.tiktok.com/@poliglotai"]')).toHaveAttribute("aria-label", "Open NERIVA on TikTok");
```

- [ ] **Step 2: Run RED web app test**

Run:

```bash
npm --prefix web-react run e2e -- web-smoke.spec.ts -g "today plan and settings password helper copy stay localized"
```

Expected: FAIL because `.settings-social-card-v2` does not exist yet.

- [ ] **Step 3: Commit test only if project policy allows intermediate commits**

Run:

```bash
git add web-react/e2e/web-smoke.spec.ts
git commit -m "test settings social links"
```

Expected: commit succeeds, or skip the commit if the working tree contains unrelated staged changes from another worker.

---

### Task 4: Web App Implementation

**Files:**
- Modify: `web-react/src/components/BrandIcons.tsx`
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add brand icons to the app**

Append these exports to `web-react/src/components/BrandIcons.tsx`:

```tsx
export const YouTubeIcon: React.FC<{ className?: string }> = ({ className }) => (
  <svg className={className} fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z" />
  </svg>
);

export const InstagramIcon: React.FC<{ className?: string }> = ({ className }) => (
  <svg className={className} fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    <path d="M7.0301.084c-1.2768.0602-2.1487.264-2.911.5634-.7888.3075-1.4575.72-2.1228 1.3877-.6652.6677-1.075 1.3368-1.3802 2.127-.2954.7638-.4956 1.6365-.552 2.914-.0564 1.2775-.0689 1.6882-.0626 4.947.0062 3.2586.0206 3.6671.0825 4.9473.061 1.2765.264 2.1482.5635 2.9107.308.7889.72 1.4573 1.388 2.1228.6679.6655 1.3365 1.0743 2.1285 1.38.7632.295 1.6361.4961 2.9134.552 1.2773.056 1.6884.069 4.9462.0627 3.2578-.0062 3.668-.0207 4.9478-.0814 1.28-.0607 2.147-.2652 2.9098-.5633.7889-.3086 1.4578-.72 2.1228-1.3881.665-.6682 1.0745-1.3378 1.3795-2.1284.2957-.7632.4966-1.636.552-2.9124.056-1.2809.0692-1.6898.063-4.948-.0063-3.2583-.021-3.6668-.0817-4.9465-.0607-1.2797-.264-2.1487-.5633-2.9117-.3084-.7889-.72-1.4568-1.3876-2.1228C21.2982 1.33 20.628.9208 19.8378.6165 19.074.321 18.2017.1197 16.9244.0645 15.6471.0093 15.236-.005 11.977.0014 8.718.0076 8.31.0215 7.0301.0839m.1402 21.6932c-1.17-.0509-1.8053-.2453-2.2287-.408-.5606-.216-.96-.4771-1.3819-.895-.422-.4178-.6811-.8186-.9-1.378-.1644-.4234-.3624-1.058-.4171-2.228-.0595-1.2645-.072-1.6442-.079-4.848-.007-3.2037.0053-3.583.0607-4.848.05-1.169.2456-1.805.408-2.2282.216-.5613.4762-.96.895-1.3816.4188-.4217.8184-.6814 1.3783-.9003.423-.1651 1.0575-.3614 2.227-.4171 1.2655-.06 1.6447-.072 4.848-.079 3.2033-.007 3.5835.005 4.8495.0608 1.169.0508 1.8053.2445 2.228.408.5608.216.96.4754 1.3816.895.4217.4194.6816.8176.9005 1.3787.1653.4217.3617 1.056.4169 2.2263.0602 1.2655.0739 1.645.0796 4.848.0058 3.203-.0055 3.5834-.061 4.848-.051 1.17-.245 1.8055-.408 2.2294-.216.5604-.4763.96-.8954 1.3814-.419.4215-.8181.6811-1.3783.9-.4224.1649-1.0577.3617-2.2262.4174-1.2656.0595-1.6448.072-4.8493.079-3.2045.007-3.5825-.006-4.848-.0608M16.953 5.5864A1.44 1.44 0 1 0 18.39 4.144a1.44 1.44 0 0 0-1.437 1.4424M5.8385 12.012c.0067 3.4032 2.7706 6.1557 6.173 6.1493 3.4026-.0065 6.157-2.7701 6.1506-6.1733-.0065-3.4032-2.771-6.1565-6.174-6.1498-3.403.0067-6.156 2.771-6.1496 6.1738M8 12.0077a4 4 0 1 1 4.008 3.9921A3.9996 3.9996 0 0 1 8 12.0077" />
  </svg>
);

export const TikTokIcon: React.FC<{ className?: string }> = ({ className }) => (
  <svg className={className} fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    <path d="M12.525.02c1.31-.02 2.61-.01 3.91-.02.08 1.53.63 3.09 1.75 4.17 1.12 1.11 2.7 1.62 4.24 1.79v4.03c-1.44-.05-2.89-.35-4.2-.97-.57-.26-1.1-.59-1.62-.93-.01 2.92.01 5.84-.02 8.75-.08 1.4-.54 2.79-1.35 3.94-1.31 1.92-3.58 3.17-5.91 3.21-1.43.08-2.86-.31-4.08-1.03-2.02-1.19-3.44-3.37-3.65-5.71-.02-.5-.03-1-.01-1.49.18-1.9 1.12-3.72 2.58-4.96 1.66-1.44 3.98-2.13 6.15-1.72.02 1.48-.04 2.96-.04 4.44-.99-.32-2.15-.23-3.02.37-.63.41-1.11 1.04-1.36 1.75-.21.51-.15 1.07-.14 1.61.24 1.64 1.82 3.02 3.5 2.87 1.12-.01 2.19-.66 2.77-1.61.19-.33.4-.67.41-1.06.1-1.79.06-3.57.07-5.36.01-4.03-.01-8.05.02-12.07z" />
  </svg>
);
```

- [ ] **Step 2: Add Settings social card**

In `web-react/src/App.tsx`, change the icon import:

```tsx
import { InstagramIcon, TelegramIcon, TikTokIcon, YouTubeIcon } from "./components/BrandIcons";
```

Add near `const aiRouterTelegramURL`:

```tsx
const poliglotSocialLinks = [
  { name: "YouTube", href: "https://www.youtube.com/@NERIVA", label: "Open NERIVA on YouTube", Icon: YouTubeIcon },
  { name: "Instagram", href: "https://www.instagram.com/poliglotai.online/", label: "Open NERIVA on Instagram", Icon: InstagramIcon },
  { name: "TikTok", href: "https://www.tiktok.com/@poliglotai", label: "Open NERIVA on TikTok", Icon: TikTokIcon },
] as const;
```

Inside `SettingsView`, in `.settings-side-stack-v2` after the Telegram section and before the activation form, render:

```tsx
        <section className="v2-panel settings-card-v2 settings-social-card-v2">
          <span className="eyebrow">{copy("social_channels", "Social channels")}</span>
          <h2>{copy("poliglot_social_title", "Соцсети NERIVA")}</h2>
          <p>{copy("poliglot_social_body", "Короткие уроки, обновления и советы по обучению.")}</p>
          <div className="settings-social-links-v2" aria-label={copy("poliglot_social_title", "Соцсети NERIVA")}>
            {poliglotSocialLinks.map(({ name, href, label, Icon }) => (
              <a key={name} className={`settings-social-link-v2 settings-social-link-v2--${name.toLowerCase()}`} href={href} aria-label={label} target="_blank" rel="noreferrer">
                <Icon />
                <span>{name}</span>
              </a>
            ))}
          </div>
        </section>
```

- [ ] **Step 3: Add app styles**

In `web-react/src/styles/app.css`, near Settings styles, add:

```css
.settings-social-card-v2 p {
  margin: 0;
  color: var(--muted);
  line-height: 1.45;
}

.settings-social-links-v2 {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.settings-social-link-v2 {
  min-width: 44px;
  min-height: 44px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: var(--surface-strong);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 9px 12px;
  color: var(--text);
  text-decoration: none;
  transition: transform 180ms ease, border-color 180ms ease, background 180ms ease;
}

.settings-social-link-v2 svg {
  width: 19px;
  height: 19px;
  flex: 0 0 auto;
}

.settings-social-link-v2:hover,
.settings-social-link-v2:focus-visible {
  border-color: color-mix(in srgb, var(--primary) 44%, var(--line));
  background: color-mix(in srgb, var(--primary) 10%, var(--surface-strong));
  outline: none;
  transform: translateY(-1px);
}

.settings-social-link-v2--youtube svg {
  color: #ff0000;
}

.settings-social-link-v2--instagram svg {
  color: #ff0069;
}

.settings-social-link-v2--tiktok svg {
  color: var(--text);
  filter: drop-shadow(1px 0 #25f4ee) drop-shadow(-1px 0 #fe2c55);
}
```

Change `.settings-side-stack-v2` from two rows to three rows:

```css
  grid-template-rows: auto auto auto;
```

- [ ] **Step 4: Run GREEN web app test**

Run:

```bash
npm --prefix web-react run e2e -- web-smoke.spec.ts -g "today plan and settings password helper copy stay localized"
```

Expected: PASS.

- [ ] **Step 5: Run web app build**

Run:

```bash
npm --prefix web-react run build
```

Expected: exit code 0 and Vite build output.

- [ ] **Step 6: Commit web app slice**

Run:

```bash
git add web-react/src/components/BrandIcons.tsx web-react/src/App.tsx web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts
git commit -m "Add settings social links"
```

Expected: commit succeeds.

---

### Task 5: Full Verification, Commit Cleanup, And GitHub Sync

**Files:**
- Read: changed files from Tasks 1-4
- Modify: no source files unless verification finds a defect

- [ ] **Step 1: Run deploy regression contract**

Run:

```bash
go test ./...
npm --prefix web-react run build
npm --prefix web-react run e2e
npm --prefix site-react run build
npm --prefix site-react run e2e
node tools/check_encoding_artifacts.mjs
```

Expected: each command exits 0.

- [ ] **Step 2: Inspect changed files**

Run:

```bash
git status --short --branch
git diff --check
git diff --stat HEAD
```

Expected: no whitespace errors. Unrelated pre-existing changes such as `site-react/playwright-report/index.html` or `english-coach-bot` remain unstaged unless this task intentionally updates them.

- [ ] **Step 3: Push the branch**

Run:

```bash
git push origin codex/ai-tutor-rebuild-fix
```

Expected: push succeeds to `https://github.com/GarryNaxyison/My-app-local.git`.

---

### Task 6: Production Deploy

**Files:**
- Deploy local build outputs from `web/`
- Deploy public-site output from `Сайт полиглота для бота/`
- Do not commit server-only backup files

- [ ] **Step 1: Confirm deploy target**

Run:

```bash
ssh root@186.246.45.123 'hostname && date && systemctl is-active aibot.service && systemctl is-active caddy'
```

Expected: SSH connects, hostname/date print, and both services report `active`. If password auth is requested, use the credential from the private project instructions and do not write it to logs, docs, or commits.

- [ ] **Step 2: Prepare deploy archives**

Run:

```bash
mkdir -p tmp
tar -czf tmp/poliglot-social-web.tar.gz web/index.html web/assets web/manifest.webmanifest web/offline-deck-sw.js
tar -czf tmp/poliglot-social-site.tar.gz -C "Сайт полиглота для бота" poliglot-ai.html privacy.html terms.html agreement.html consent.html assets
```

Expected: both archives are created under `tmp/`.

- [ ] **Step 3: Upload archives**

Run:

```bash
scp tmp/poliglot-social-web.tar.gz tmp/poliglot-social-site.tar.gz root@186.246.45.123:/tmp/
```

Expected: both uploads complete.

- [ ] **Step 4: Back up and extract on server**

Run:

```bash
ssh root@186.246.45.123 'set -euo pipefail
stamp=$(date +%Y%m%d-%H%M%S-social-links)
mkdir -p /opt/aibot/deploy-backups/$stamp/web /opt/aibot/deploy-backups/$stamp/site
cp -a /opt/aibot/web/. /opt/aibot/deploy-backups/$stamp/web/
cp -a /var/www/poliglotai/. /opt/aibot/deploy-backups/$stamp/site/
tar -xzf /tmp/poliglot-social-web.tar.gz -C /opt/aibot
tar -xzf /tmp/poliglot-social-site.tar.gz -C /var/www/poliglotai
chown -R root:root /opt/aibot/web /var/www/poliglotai
find /opt/aibot/web /var/www/poliglotai -type d -exec chmod 755 {} \;
find /opt/aibot/web /var/www/poliglotai -type f -exec chmod 644 {} \;
systemctl reload caddy
systemctl is-active aibot.service
systemctl is-active caddy
'
```

Expected: extraction succeeds, Caddy reloads, and both service checks report `active`.

- [ ] **Step 5: Verify production**

Run:

```bash
curl -fsS https://poliglotai.ru/healthz
curl -fsS https://poliglotai.ru/poliglot-ai.html | grep -E 'Follow NERIVA|youtube.com/@NERIVA|instagram.com/poliglotai.online|tiktok.com/@poliglotai'
curl -fsS https://poliglotai.ru/app/ | grep -E 'assets/.*\.js|assets/.*\.css'
```

Expected:

```text
ok
```

The landing HTML check prints the social-link markers. The app shell check prints current asset references.

- [ ] **Step 6: Record deploy result**

Append one line to `docs/tracking/PROJECT_TRACKING.md` under the deployment history section with the date, target paths, backup stamp, commands verified, and production endpoints checked.

Commit and push:

```bash
git add docs/tracking/PROJECT_TRACKING.md
git commit -m "Track social links deploy"
git push origin codex/ai-tutor-rebuild-fix
```

Expected: deployment history is stored in GitHub.

---

## Self-Review

- Spec coverage: landing hero row, footer links, Settings card, branded local SVG icons, accessibility labels, new-tab behavior, no hotlinking, responsive tap targets, build/e2e checks, and deploy are all covered.
- Placeholder scan: no incomplete task marker or vague test instruction is present.
- Type consistency: both frontend link models use `name`, `href`, `label`, and `Icon`; CSS class names match the test selectors and JSX.
