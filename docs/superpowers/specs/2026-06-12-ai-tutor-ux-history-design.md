# AI Tutor UX History Design

## Goal

Исправить пять UX-проблем AI Tutor: разделитель слов в истории, плейсхолдер ответа, историю пройденных уроков с повторным запуском урока, просмотр уже отвеченных шагов и блоки "Как начать".

## Scope

- AI Tutor story word list shows `target - translation`, not `target ? translation`.
- Answer textarea placeholder is a localized "enter your answer in this field" phrase instead of generic "Section".
- Completed AI Tutor lessons come from backend history, not only browser `localStorage`.
- Completed lesson history lets the learner start the same lesson again.
- Already completed AI Tutor steps are clickable and open as read-only material.
- The active current step remains the only answerable step.
- "Quick start guide" remains five visual blocks, but each block describes a product capability.

## Out Of Scope

- Editing old answers.
- Showing every historical answer and feedback inside completed lesson history.
- Replacing the AI Tutor lesson generation flow.
- Adding a full course archive outside AI Tutor.

## Backend Design

Add a small completed-history API:

- `GET /api/ai-tutor/completed`
  - Requires the current authenticated web user.
  - Returns newest completed AI Tutor sessions for the user.
  - Each item contains `session_id`, `lesson_id`, `title`, `topic`, `level`, `completed_at`, and a compact `lesson` payload needed by the frontend history card.

Add a repeat API:

- `POST /api/ai-tutor/restart`
  - Body: `{ "lesson_id": "..." }`.
  - Requires the current authenticated web user.
  - Verifies the user has completed that lesson before, or the lesson belongs to an existing user session.
  - Creates a new AI Tutor session for the same lesson at `story_intro`.
  - Returns the existing `aiTutorDTO` shape so the React state can switch into the repeated lesson immediately.

The store interface gets two focused methods:

- `completedAITutorLessons(telegramID int64, limit int) ([]aiTutorCompletedLessonRecord, error)`
- `userCanRestartAITutorLesson(telegramID int64, lessonID string) (bool, error)`

SQLite reads from `ai_tutor_sessions` joined to `ai_tutor_lessons`, filtering `status=complete` or non-empty `completed_at`. JSON store uses its in-memory maps.

## Frontend Design

`TutorView` keeps the existing completed-lessons dialog, but loads backend history when the dialog opens. If the backend fails or returns nothing, it can still show local `localStorage` records.

Completed lesson cards become action buttons:

- click a card to call `/api/ai-tutor/restart`;
- close the dialog;
- replace the current `aiTutorStep` with the returned repeated lesson step;
- show loading/disabled state while restart is running.

The left lesson-flow buttons become navigable for indexes lower than the current stage index. Clicking a completed step sets a local `previewStage`. The main panel renders `aiTutorBuildStep`-equivalent material from the current lesson data for that previous stage. Preview mode is read-only and shows a "Back to current step" action. Current and future steps do not become editable through preview.

For the story word list, render words as `target - translation`.

For placeholder copy, add/derive a `tutor_answer_field_placeholder` key and use it for non-production free-text stages. Production keeps its specific "write 2-3 sentences" placeholder.

For "How to start", keep five cards but update the copy:

1. AI Tutor builds a full mini-lesson: story, words, questions, writing, review.
2. Vocabulary and audio keep words tied to examples and pronunciation.
3. Practice modes repair weak spots: mistakes, spelling, pronunciation.
4. Phrasebook saves useful corrections and phrases for reuse.
5. Progress shows XP, streak, level and completed lessons.

## Error Handling

- If completed-history loading fails, show existing local records and keep the modal usable.
- If restart fails, keep the dialog open and show a localized error message.
- If a lesson from history no longer exists, backend returns `404` and frontend does not change the active lesson.
- If a user tries to restart a lesson they never completed or opened, backend returns `403`.

## Testing

Backend:

- `GET /api/ai-tutor/completed` returns completed sessions newest first with lesson summaries.
- `GET /api/ai-tutor/completed` does not return another user's lessons.
- `POST /api/ai-tutor/restart` creates a new session for a completed lesson and returns `story_intro`.
- `POST /api/ai-tutor/restart` rejects unrelated lessons.

Frontend/build:

- TypeScript build passes.
- React source contains `/api/ai-tutor/completed` and `/api/ai-tutor/restart` API calls.
- AI Tutor sidebar completed buttons are not disabled and future steps stay disabled.
- Story word separator renders as ` - `.

Browser verification:

- Open the React app.
- Confirm "Как начать" renders as five separate capability cards.
- Confirm AI Tutor answer placeholder is no longer "Раздел".
- Confirm completed/history modal shows backend items when available.

