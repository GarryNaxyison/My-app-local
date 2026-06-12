# Daily Premium Tutor Overflow Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix daily bonus state, Premium-only audio/pronunciation access, AI Tutor history labels, overflow, and repeated tutor lesson generation.

**Architecture:** Keep the existing Go API, SQLite/json store, and React app structure. Add small backend helpers for AI Tutor used-sequence counting and focused frontend guards for labels, gating, and layout containment.

**Tech Stack:** Go tests, SQLite-backed store, React 19 + TypeScript + Vite, Nx wrapper commands, Playwright/browser verification when available.

---

### Task 1: Daily Bonus Regression

**Files:**
- Modify: `web_api_feature_test.go`
- Modify: `web-react/src/App.tsx`

- [ ] **Step 1: Write failing daily API test**

Add `TestWebDailyClaimInsideTwentyFourHoursDoesNotAwardXP` in `web_api_feature_test.go`. The test should call `/api/daily/claim` twice for the same authenticated user and assert the second response has `claimed=false` and no XP increase.

- [ ] **Step 2: Verify RED**

Run:

```bash
npm run test -- --runInBand web_api_feature_test.go -run TestWebDailyClaimInsideTwentyFourHoursDoesNotAwardXP
```

Expected before frontend changes: backend likely passes. If it passes, keep it as regression evidence and continue to the frontend state fix.

- [ ] **Step 3: Fix frontend rejected-claim handling**

In `claimDailyBonus`, when `record.claimed === false`, show `daily_bonus_already_claimed`, refresh session if possible, and do not set today's `habitLog` to `claimed: true` or `complete: true`.

- [ ] **Step 4: Fix button copy**

In `HomeView`, replace the claim button label logic with:

- locked: `bonus_claimed` fallback `"Бонус уже получен"`;
- not complete: `complete_daily_first` fallback `"Сначала выполните дейлик"`;
- claimable: `claim_daily_bonus` fallback `"Забрать XP"`.

### Task 2: Premium-Only Listening And Pronunciation

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/lib/i18n.ts`
- Modify: `web_app_shell_test.go`

- [ ] **Step 1: Write frontend contract test**

In `TestReactFrontendKeepsGoAPIContracts` or a nearby React source test, assert the source contains `premium_audio_required`, `premium_pronunciation_required`, and Premium feature copy mentioning listening/pronunciation as Premium.

- [ ] **Step 2: Verify RED**

Run:

```bash
npm run test -- --runInBand web_app_shell_test.go -run TestReactFrontendKeepsGoAPIContracts
```

Expected: fail until the new copy/actions exist.

- [ ] **Step 3: Add premium routing helper**

Add a helper near `startTutor`:

```ts
const requirePremiumFeature = (messageKey: string, fallback: string) => {
  if (user.premium) return true;
  const text = copy(messageKey, fallback);
  setStatus({ kind: "info", text });
  setView("premium");
  return false;
};
```

Use it before `startShadowing`, before opening `pronunciation`, and in home/learning-lab/route actions.

- [ ] **Step 4: Update plan copy**

Update `premiumPlanFeatures`, `premiumPlanBody`, and i18n defaults so Free excludes AI Tutor, listening, and pronunciation, while Premium/Platinum include them.

### Task 3: AI Tutor History Cleanup And Overflow

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`
- Modify: `web_app_shell_test.go`

- [ ] **Step 1: Write source assertions**

Add assertions that `web-react/src/App.tsx` contains `cleanTutorHistoryLabel` and a guard for `To talk about`.

- [ ] **Step 2: Verify RED**

Run:

```bash
npm run test -- --runInBand web_app_shell_test.go -run TestReactFrontendKeepsGoAPIContracts
```

Expected: fail until the helper exists.

- [ ] **Step 3: Add history label helper**

Add `cleanTutorHistoryLabel(value, fallback)` near `cleanTutorPrompt`. It should return fallback for empty text, generic section text, prompt labels, `Topic/theme seed`, `Goal:`, and strings starting with `To talk about`.

- [ ] **Step 4: Use helper in completed lessons**

Use the helper in `tutorCompletedLessonFromAPI`, `writeTutorCompletedLesson`, and completed card render so history never prefers system-like goal text.

- [ ] **Step 5: Add CSS containment**

Add targeted wrap and containment rules for tutor panels, message bodies, feedback, option buttons, completed lesson dialog rows, and textareas.

### Task 4: Non-Repeating AI Tutor Lessons

**Files:**
- Modify: `storage.go`
- Modify: `sqlite_store.go`
- Modify: `ai_tutor_engine.go`
- Modify: `ai_tutor_engine_test.go`

- [ ] **Step 1: Write failing test**

Add `TestAITutorStartUsesUnusedGeneratedSequenceForRepeatedStarts` in `ai_tutor_engine_test.go`. Start AI Tutor twice for the same user without completion and assert the second result has a different `Lesson.ID`.

- [ ] **Step 2: Verify RED**

Run:

```bash
npm run test -- --runInBand ai_tutor_engine_test.go -run TestAITutorStartUsesUnusedGeneratedSequenceForRepeatedStarts
```

Expected: fail because both starts use the same generated lesson ID.

- [ ] **Step 3: Add store helper**

Add `aiTutorSessionCountForContext(telegramID int64, language string, interfaceLanguage string, levelBand string) (int, error)` to `store`. SQLite joins `ai_tutor_sessions` to `ai_tutor_lessons`; jsonStore counts matching sessions in memory.

- [ ] **Step 4: Use helper in engine start**

When the AI Tutor engine must generate a local lesson, use `tutorLessonSequence(user) + usedCount` and increment until `getAITutorLesson(generated.ID)` is not already used by that user.

- [ ] **Step 5: Verify GREEN**

Run the focused AI Tutor test again and then:

```bash
npm run test -- --runInBand ai_tutor_engine_test.go -run TestAITutorStart
```

### Task 5: Verification And Sync

**Files:**
- All modified files

- [ ] **Step 1: Build frontend**

Run:

```bash
npm --prefix web-react run build
```

- [ ] **Step 2: Run project checks**

Run:

```bash
npm run check
```

- [ ] **Step 3: Browser verification**

Attempt Playwright/browser verification for desktop and mobile overflow. If the MCP browser remains blocked by Chrome policy, record that limitation and rely on build plus source/CSS checks.

- [ ] **Step 4: Commit and push**

Run:

```bash
git status --short
git add docs/superpowers/specs/2026-06-12-daily-premium-tutor-overflow-design.md docs/superpowers/plans/2026-06-12-daily-premium-tutor-overflow.md web_api_feature_test.go web_app_shell_test.go storage.go sqlite_store.go ai_tutor_engine.go ai_tutor_engine_test.go web-react/src/App.tsx web-react/src/lib/i18n.ts web-react/src/styles/app.css web
git commit -m "Fix daily premium tutor web regressions"
git push origin codex/ai-tutor-rebuild-fix
```
