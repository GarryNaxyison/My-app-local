# Desktop Web App Stability Design

## Goal

Remove the reported desktop-only layout failures and prevent learning flows from contradicting their own examples, without changing Android or the mobile web experience.

## Scope

- Desktop web only: viewports wider than 760px, with a primary verification viewport of 2048×1050.
- `web-react` desktop shell, desktop workspaces, login form, phrasebook suggestions, and the Go web API prompt routing.
- No Android files and no changes inside the existing `@media (max-width: 760px)` mobile rules.

## Design

1. The desktop application shell remains a fixed-height three-row grid. Its content row and every desktop workspace will explicitly use `min-width: 0`, `width: 100%`, and `min-height: 0`, so a wide viewport is filled by the active view rather than sized by its contents. The desktop content surface owns vertical scrolling; nested panels do not create competing full-height scroll regions.
2. Desktop pronunciation, roleplay, and chat workspaces use bounded grid tracks and document-flow controls. Their lower regions remain reachable through the content-surface scrollbar rather than being clipped by an inherited `height: 100%`/`overflow: hidden` combination.
3. Login owns its vertical scroll only. The form card will hide horizontal overflow, and the Turnstile slot will constrain third-party iframe width to its available content width.
4. Phrase suggestion and saved-state comparison use one canonical key: whitespace/case normalized and terminal decorative punctuation removed. This treats `Hello`, `Hello!`, and `Hello.` as the same note while retaining the original display text.
5. A roleplay submission is identified by its `ROLEPLAY_TOOL_V2` envelope and evaluated with the roleplay prompt, not the generic practice-chat prompt. The roleplay prompt states that only the exact learner line may be evaluated. A lesson answer that matches an offered model phrase after case/whitespace/punctuation normalization gets an explicit successful feedback response and an empty mistake list.
6. The logout control gets a stable desktop footprint and immediate visible base styling; text animation is decorative only and cannot make the action disappear during initial render.

## Validation

- Go tests cover roleplay prompt selection and model-answer equivalence.
- A focused web regression test verifies canonical phrase deduplication.
- Production web build and targeted Go tests pass.
- Desktop browser inspection at 2048×1050 verifies full-width workspaces, accessible lower content, visible logout, and no horizontal scroll in the auth card. A 760px mobile smoke check confirms no mobile CSS selector was changed.
