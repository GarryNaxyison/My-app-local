# Telegram AI Tutor And Pronunciation Stage Design

## Status

Approved for implementation on 2026-06-07.

## Problem

Telegram already shows an `AI Tutor` entry in the main and learning menus, but the callback `menu_tutor` is not handled by `handleCallbackQuery`. A user pressing the button falls through to the generic unknown-button response, so AI Tutor is effectively unavailable in the Telegram bot.

The same button is hardcoded as Russian text, so it bypasses the existing `uiCopy` localization pipeline used by the rest of the menu. New visible buttons must be localized for every interface language returned by `interfaceLanguages()`.

On the web desktop AI Tutor pronunciation stage, the most important working area is visually too small. The target phrase, audio, recording/upload controls, and report compete with each other instead of giving the phrase a clear primary surface.

## Goals

- Make AI Tutor start from Telegram via `menu_tutor`.
- Reuse the same stored tutor lesson generator used by web `/api/tutor/start`.
- Add a localized AI Tutor button label through the shared Telegram UI copy path for all interface languages.
- Render a useful Telegram lesson summary with scenario, mini-explanation, target phrase, choice, dialogue, and navigation.
- Preserve existing onboarding, menu, premium, tools, word practice, and web auth behavior.
- Redesign only the web AI Tutor pronunciation stage inside `.tutor-context-v2` on desktop.
- Keep mobile pronunciation as a single top-to-bottom stack.
- Do not change words, choice, writing, or dialogue stage layouts.

## Non-Goals

- Do not move every web feature into Telegram in this change.
- Do not add a multi-step interactive Telegram tutor engine yet.
- Do not change tutor lesson generation rules beyond the already implemented scenario-slot generator.
- Do not redesign the whole AI Tutor web layout.
- Do not change onboarding order, privacy confirmation, language selection, timezone selection, or level assessment.

## Telegram AI Tutor Design

`handleCallbackQuery` will handle `menu_tutor` and call a new bot method, `startTutorLesson`.

`startTutorLesson` will:

- clear previous pronunciation audio like the other learning flows;
- call `store.nextTutorLesson(user, factory)` with `buildTutorLessonForSequence(tutorReusableLessonUser(user), sequence)`;
- render the returned `tutorLesson` into a compact Telegram MarkdownV2 message;
- send an inline keyboard with a localized "next tutor lesson" button and a localized menu/back button.

The first Telegram version is a read-only lesson card, not a full web clone. It should still be useful in chat: the user sees what to say, why, and how a good answer looks. Follow-up interactive scoring belongs in a separate tutor-chat spec after this bugfix ships.

The message will include:

- localized title and lesson number;
- goal;
- scenario;
- mini-explanation;
- target phrase from `PronunciationText` or `ScenarioSlots.ModelAnswer`;
- one scenario choice with options and correct-answer feedback;
- mini-dialogue line from the tutor and one recommended user reply.

Telegram callback actions:

- `menu_tutor`: start or advance to a tutor lesson;
- `back_menu`: unchanged existing return to main menu.

## Telegram Localization Design

Add `AITutor` to `uiCopy`. Existing explicit UI copies receive localized labels where the project already maintains explicit translations. Generated compact copies derive a stable label from a helper `aiTutorButtonLabel(code)`.

The helper will return a localized phrase for all supported interface languages. If a language is not explicitly mapped, it falls back to `AI Tutor`, which is acceptable because "AI" and "Tutor" are product-like terms and keeps the button non-empty.

Both `mainMenuInlineKeyboard` and `learningMenuKeyboard` will use `copy.AITutor` instead of hardcoded `AI Репетитор`.

Tests must verify that every `interfaceLanguages()` copy has a non-empty `AITutor` label and that no menu uses the Russian hardcoded label for non-Russian UI languages.

## Web Pronunciation Stage Design

Only the visual layout of the pronunciation stage inside `.tutor-context-v2` changes.

Desktop behavior:

- the pronunciation block becomes the main working screen for that stage;
- target phrase moves into a large primary panel with larger text and a stable height;
- audio controls sit directly under the target phrase;
- recording/upload controls and the checking result move into a right-side tools column;
- after checking, the report is visible and prominent, but it does not push the target phrase out of the primary position.

Mobile behavior:

- no two-column layout;
- content remains stacked top-to-bottom;
- target phrase stays first, then audio, controls, and report.

The CSS change must be scoped to the pronunciation stage so words, choice, writing, and dialogue stages keep their current layout.

## Testing

Go tests:

- `menu_tutor` callback sends a tutor lesson instead of the unknown-button message.
- AI Tutor menu labels are localized/non-empty for every interface language.
- Existing `/start` privacy/onboarding behavior remains unchanged.
- Tutor Telegram rendering contains scenario content and does not leak internal service criteria.

Frontend tests:

- existing AI Tutor cafe scenario smoke test remains passing.
- add or extend a Playwright check for desktop pronunciation stage so the primary target phrase panel is the dominant area and the report can appear without replacing it.
- add or preserve a mobile check that pronunciation remains a single stack.

Full verification:

- `go test ./...`
- relevant `npm --prefix web-react run e2e` AI Tutor tests

## Rollout

This change is safe to ship behind existing menu access because AI Tutor lesson generation is local and already used by web. If lesson generation fails, Telegram should send the localized runtime error instead of silently falling through to UnknownButton.

Full web-to-Telegram parity will be handled in a separate spec after this bugfix and pronunciation visual pass are verified.
