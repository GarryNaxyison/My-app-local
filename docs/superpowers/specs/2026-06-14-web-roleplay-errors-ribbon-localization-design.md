# Web Roleplay, Errors, And Ribbon Localization Design

## Goal

Fix three UI issues in the web app: duplicated-looking mistake categories, untranslated roleplay cards, and awkward desktop ribbon arrow placement.

## Requirements

- Mistake category labels must be distinct and meaningful across all supported interface languages.
- Roleplay screen title, subtitle, scenario titles, and scenario descriptions must be localized for all 35 interface languages.
- Russian roleplay cards must not show English fallback titles or descriptions.
- Desktop function ribbon previous/next buttons must stay inside the app frame, be compact enough for the ribbon, and align cleanly with the scrollable menu.

## Design

Localization stays in `web-react/src/lib/i18n.ts`, which is already the source of `appCopy`. The app will read roleplay scenario titles and descriptions from localized keys before falling back to existing local dictionaries. Mistake category labels will also come from `appCopy` keys, with explicit values for every app locale.

The ribbon keeps the existing `MorphingArrowButton` component, but the component gets a compact mode for the function ribbon. CSS will remove conflicting negative margins and scale overrides so the two arrow buttons occupy predictable grid columns.

## Testing

The existing Playwright spec will be expanded with static localization assertions for all 35 locale codes and a stronger desktop geometry assertion for the ribbon arrows.
