BEGIN IMMEDIATE;

DELETE FROM ai_tutor_answers
WHERE session_id IN (SELECT id FROM ai_tutor_sessions);

DELETE FROM mistakes;
DELETE FROM learned_words;
DELETE FROM phrasebook_entries;
DELETE FROM daily_bonus_claims;
DELETE FROM habit_days;
DELETE FROM tutor_user_lessons;
DELETE FROM ai_tutor_reviews;
DELETE FROM ai_tutor_word_reports;
DELETE FROM ai_tutor_sessions;

UPDATE users
SET mode = 'idle',
    level = 'A2',
    lesson_count = 0,
    practice_count = 0,
    voice_count = 0,
    word_lesson_count = 0,
    word_game_count = 0,
    xp = 0,
    daily_date = '',
    lessons_today = 0,
    practice_today = 0,
    voice_today = 0,
    lesson_history = '',
    practice_history = '',
    last_lesson_prompt = '',
    updated_at = datetime('now');

COMMIT;
