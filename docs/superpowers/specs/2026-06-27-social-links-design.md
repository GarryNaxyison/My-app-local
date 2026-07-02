# NERIVA Social Links Design

## Goal

Add YouTube, Instagram, and TikTok links so the landing page shows NERIVA as an active public brand and gently sends interested visitors to social content. The primary conversion on the landing page remains the web app or Telegram start flow.

## Channels

- YouTube: https://www.youtube.com/@NERIVA
- Instagram: https://www.instagram.com/poliglotai.online/
- TikTok: https://www.tiktok.com/@poliglotai.online

## Placement

### Landing Page

Add a compact social proof row to the English landing page near the hero proof area. It should sit below the existing product proof metrics or immediately before the first route section. The row uses the label `Follow NERIVA` and three icon links.

Duplicate the same destinations in the site footer as a dedicated `Social` column. Legal pages should keep their current footer structure, with social links added in a way that does not push legal links below the fold on common mobile widths.

### Web App

Add a small `Follow NERIVA` card inside `SettingsView`, near the existing Telegram and activation cards. Do not add these links to the main navigation, function ribbon, or core learning surfaces because they are not learning actions.

## Visual Design

Use recognizable platform icons:

- YouTube in brand red and white.
- Instagram in its gradient app-icon style where practical.
- TikTok in dark/white with cyan and red accent treatment where practical.

Use official or stable icon sources during implementation:

- YouTube official brand resources: https://brand.youtube/
- Instagram official brand page: https://about.instagram.com/brand
- TikTok brand hub or developer brand guidelines: https://www.tiktokbrandhub.com/visual-identity/logo and https://developers.tiktok.com/doc/getting-started-design-guidelines
- Simple Icons may be used as a consistent SVG fallback when official downloadable assets are not convenient for the app build.

Do not hotlink third-party image URLs at runtime. Store the icons as local SVG/React components or local asset files so the landing page and web app remain fast, stable, and privacy-friendly.

## Interaction

Each social link opens in a new tab with `target="_blank"` and `rel="noreferrer"`. Links need descriptive accessible labels:

- `Open NERIVA on YouTube`
- `Open NERIVA on Instagram`
- `Open NERIVA on TikTok`

The icon buttons should have visible hover, focus, and active states. Keyboard users must be able to tab through the links in a predictable order: YouTube, Instagram, TikTok.

## Components

Create a shared social link model with platform name, URL, accessible label, and icon. Reuse it in both frontends if the existing project boundaries allow it cleanly. If sharing would add unnecessary build coupling between `site-react` and `web-react`, duplicate a small local constant in each app and keep the URLs identical.

Extend the existing `BrandIcons` pattern in the web app with YouTube, Instagram, and TikTok icons. The site app can either use the same inline SVG definitions or a local `SocialIcons` helper.

## Copy

Landing row:

- Label: `Follow NERIVA`
- Supporting text: `Short lessons, product updates, and learning tips.`

Web app settings card:

- Title: `Follow NERIVA`
- Body: `Short lessons, updates, and product tips.`

If Russian localization is added in the same implementation, use:

- Title: `Соцсети NERIVA`
- Body: `Короткие уроки, обновления и советы по обучению.`

## Responsive Behavior

On desktop, show the social links as a horizontal icon row. On mobile, keep the same row if it fits; otherwise wrap to a second line without changing icon size or causing layout shift.

Icon tap targets should be at least 44 by 44 CSS pixels. Text must not overlap hero proof metrics, footer columns, or settings cards.

## Error Handling

There is no runtime API dependency. If a social network blocks a URL or the user is offline, the browser handles it as a normal external link failure. The app should not show custom error UI for these links.

## Testing

Verify:

- Landing build succeeds.
- Web app build succeeds.
- The three links render on the landing page and settings page.
- Links have correct `href`, `target`, `rel`, and accessible labels.
- Keyboard focus is visible.
- Desktop and mobile layouts do not overlap or shift when icons load.

Use Playwright or the in-app browser for visual validation when available. If Playwright MCP remains blocked by local Chrome policy, use the local project test runner or another available browser path and report the limitation.
