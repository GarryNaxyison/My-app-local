package ru.neriva.app

import java.io.File

/**
 * Generates clean Android strings.xml files for all 35 locales
 * from the web app's i18n.ts source.
 * Run: cd android-app && ../gradlew -q app:run
 * Or just run this as a standalone Kotlin script.
 */
object I18nGenerator {
    
    // Core keys extracted from web-react/src/lib/i18n.ts
    // These are the most important UI keys
    val enBase = mapOf(
        "app_name" to "NERIVA",
        "ready_title" to "Ready",
        "ready_body" to "Choose a mode, start a lesson, or send a practice message.",
        "loading" to "Loading...",
        "cancel" to "Cancel",
        "done" to "Done.",
        "you" to "You",
        "try_again" to "Try again.",
        "ok" to "OK",
        "retry" to "Retry",
        "send" to "Send",
        "close" to "Close",
        "save" to "Save",
        "delete" to "Delete",
        "back" to "Back",
        "next" to "Next",
        "submit" to "Submit",
        "error_generic" to "Something went wrong. Try again.",
        "today" to "Today",
        "ai_tutor" to "AI Tutor",
        "practice" to "Practice",
        "learn_words" to "Learn Words",
        "more" to "More",
        "home" to "Home",
        "lesson" to "Lesson",
        "words" to "Words",
        "roleplay" to "Roleplay",
        "shadowing" to "Listening",
        "pronunciation" to "Pronunciation",
        "spelling" to "Spelling",
        "vocabulary" to "Vocabulary",
        "phrasebook" to "Notes",
        "offline" to "Offline Decks",
        "level_test" to "Level Test",
        "progress" to "Progress",
        "awards" to "Awards",
        "leaderboard" to "Leaderboard",
        "limits" to "Limits",
        "mistakes" to "Mistakes",
        "tools" to "Tools",
        "premium" to "Premium",
        "settings" to "Settings",
        "referral" to "Referrals",
        "bug_report" to "Bug Report",
        "completed_lessons" to "Completed Lessons",
        "dashboard" to "Dashboard",
        "sign_in" to "Sign In",
        "register" to "Register",
        "login" to "Login",
        "password" to "Password",
        "confirm_password" to "Confirm Password",
        "referral_code" to "Referral code (optional)",
        "create_account" to "Create Account",
        "ai_language_tutor" to "AI Language Tutor",
        "premium_daily_practice" to "Premium daily practice",
        "your_progress" to "Your progress",
        "daily_activity" to "Daily activity",
        "daily_bonus" to "Daily Bonus",
        "claim_daily_xp" to "Claim your daily XP bonus",
        "learning_lab" to "Learning Lab",
        "lessons" to "Lessons",
        "practice_sessions" to "Practice",
        "voice" to "Voice",
        "words_learned" to "Words",
        "interface_language" to "Interface Language",
        "learning_language" to "Learning Language",
        "level" to "Level",
        "log_out" to "Log Out",
        "your_code" to "Your code",
        "invited" to "Invited",
        "balance" to "Balance",
        "invitees" to "Invitees",
        "report_a_bug" to "Report a bug",
        "describe_issue" to "Describe the issue",
        "send_report" to "Send report",
        "report_sent" to "Report sent. Thank you!",
        "translator" to "Translator",
        "text_to_translate" to "Text to translate",
        "translate" to "Translate",
        "upgrade_to_premium" to "Upgrade to Premium",
        "spell_the_word" to "Spell the word",
        "check" to "Check",
        "hint" to "Hint",
        "give_up" to "Give up",
        "attempts" to "Attempts",
        "correct" to "Correct!",
        "wrong" to "Wrong:",
        "practice_mistakes" to "Practice Mistakes",
        "all_mistakes_practiced" to "All mistakes practiced!",
        "type_correct_form" to "Type the correct form",
        "repeat_phrase" to "Repeat:",
        "record" to "Record",
        "choose_scenario" to "Choose scenario",
        "choose_level" to "Your level",
        "start_learning" to "Start Learning",
        "overall_progress" to "Overall Progress",
        "total_xp" to "Total XP",
        "streak" to "Streak",
        "days" to "days",
        "premium_active" to "Premium Active",
        "when_review" to "When do you want to review this lesson?",
        "tomorrow" to "Tomorrow",
        "in_3_days" to "In 3 days",
        "in_1_week" to "In 1 week",
        "no_review" to "No review",
        "report_mistake" to "Report mistake",
        "correct_word" to "Correct word",
        "correct_translation" to "Correct translation",
        "comment_optional" to "Comment (optional)",
        "no_completed_lessons" to "No completed lessons yet.",
        "lesson" to "Lesson",
        "topic" to "Topic",
        "completed_at" to "Completed",
        "xp_gained" to "XP gained",
        "xp_awarded_tutor" to "+40 XP awarded for this AI Tutor lesson.",
        "tutor_completed" to "Tutor lesson completed.",
        "start_tutor_lesson" to "Start tutor lesson",
        "guided_course" to "Guided course",
        "tutor_subtitle" to "Words, grammar, listening, writing, dialogue, and review in one context window.",
        "word_game" to "Review Game",
        "word_game_subtitle" to "Refresh saved vocabulary through fast multiple-choice rounds.",
        "review_game" to "Review Game",
        "memory_bank" to "Memory bank",
        "vocabulary_subtitle" to "Review learned words, context, and training progress.",
        "phrasebook_subtitle" to "Saved lesson notes and useful phrases.",
        "offline_subtitle" to "Downloaded decks for offline practice.",
        "leaderboard_subtitle" to "Compare your progress with other learners.",
        "limits_subtitle" to "Daily usage limits and remaining counts.",
        "settings_subtitle" to "Language, level, account, and preferences.",
        "referral_subtitle" to "Invite friends and earn rewards.",
        "premium_subtitle" to "Plans, pricing, and payment options.",
        "level_test_subtitle" to "CEFR level assessment test.",
        "awards_subtitle" to "20 trophy ranks and achievements.",
        "progress_subtitle" to "XP, streak, words learned, and stats.",
        "tools_subtitle" to "Translate text and convert voice.",
        "mistakes_subtitle" to "Words and grammar you got wrong.",
        "dashboard_subtitle" to "Detailed learning statistics.",
        "completed_subtitle" to "Lessons you have finished.",
        "voice_subtitle" to "Practice pronunciation with AI.",
        "listening_subtitle" to "Shadowing and listening exercises.",
        "pronunciation_subtitle" to "Voice recording and scoring.",
        "spelling_subtitle" to "Type words from memory.",
        "roleplay_subtitle" to "AI conversation scenarios.",
        "no_words_available" to "No words available",
        "correct_xp" to "Correct! +%d XP",
        "wrong_answer" to "Wrong: %s",
        "spell_word_prompt" to "Spell the word",
        "attempts_count" to "Attempts: %d",
        "hint_shows_first_letter" to "Hint: starts with '%c'  (%d letters)",
        "answer_was" to "Answer: %s",
        "practice_mistakes_button" to "Practice Mistakes",
        "no_offline_cards" to "No offline cards yet. Save phrases to learn offline.",
        "translate_button" to "Translate",
        "send_button" to "Send",
        "submit_button" to "Submit",
        "start_button" to "Start",
        "finish_button" to "Finish",
        "retry_button" to "Retry",
        "demo_mode" to "Demo Mode (test/test)",
        "telegram_login" to "Telegram login: @NERIVAapp_bot",
        "loading_lesson" to "Building a local A1-A2 lesson from the course vocabulary.",
        "theme_vocabulary" to "Theme vocabulary",
        "new_lesson" to "New lesson",
        "no_lesson_active" to "No active lesson. Start one from the menu.",
        "repeat_after_ai" to "Listen and repeat the phrase.",
        "your_answer" to "Your answer",
        "message" to "Message",
        "reply" to "Reply",
        "choose_answer" to "Choose the correct answer:",
        "report_word_button" to "Report word",
        "word_report_title" to "Report a word mistake",
        "word_report_sent" to "Report sent.",
        "easy" to "Easy",
        "good" to "Good",
        "hard" to "Hard",
        "bad" to "Bad",
        "question_instruction" to "Answer the question before moving on.",
        "start_practice" to "Start practice",
        "stop_practice" to "Stop",
        "limit_reached" to "Daily limit reached. Upgrade to Premium for more.",
        "premium_required" to "Premium required for this feature.",
        "activation_key" to "Activation key",
        "activate_key" to "Activate",
        "key_accepted" to "Activation key accepted.",
        "monthly_price" to "per month",
        "yearly_price" to "per year",
        "active_until" to "Active until",
        "plans" to "Plans",
        "free_plan" to "Free",
        "premium_plan" to "Premium",
        "platinum_plan" to "Platinum",
        "interface" to "Interface",
        "referral_count" to "Invited: %d",
        "referral_balance" to "Balance: %s",
        "no_invitees" to "No invitees yet."
    )
    
    val ruBase = mapOf(
        "app_name" to "NERIVA",
        "today" to "Сегодня",
        "ai_tutor" to "AI Tutor",
        "practice" to "Практика",
        "learn_words" to "Учить слова",
        "more" to "Ещё",
        "home" to "Главная",
        "lesson" to "Урок",
        "words" to "Слова",
        "roleplay" to "Рольplay",
        "shadowing" to "Аудирование",
        "pronunciation" to "Произношение",
        "spelling" to "Орфография",
        "vocabulary" to "Словарь",
        "phrasebook" to "Заметки",
        "offline" to "Оффлайн",
        "level_test" to "Тест уровня",
        "progress" to "Прогресс",
        "awards" to "Награды",
        "leaderboard" to "Рейтинг",
        "limits" to "Лимиты",
        "mistakes" to "Ошибки",
        "tools" to "Инструменты",
        "premium" to "Премиум",
        "settings" to "Настройки",
        "referral" to "Рефералы",
        "bug_report" to "Сообщить об ошибке",
        "completed_lessons" to "Пройденные уроки",
        "dashboard" to "Панель",
        "sign_in" to "Войти",
        "register" to "Регистрация",
        "login" to "Логин",
        "password" to "Пароль",
        "confirm_password" to "Подтвердите пароль",
        "create_account" to "Создать аккаунт",
        "ai_language_tutor" to "AI Языковой Tutor",
        "premium_daily_practice" to "Премиум ежедневная практика",
        "your_progress" to "Ваш прогресс",
        "daily_activity" to "Дневная активность",
        "daily_bonus" to "Дневной бонус",
        "learning_lab" to "Учебная лаборатория",
        "lessons" to "Уроки",
        "interface_language" to "Язык интерфейса",
        "learning_language" to "Язык обучения",
        "level" to "Уровень",
        "log_out" to "Выйти",
        "your_code" to "Ваш код",
        "invited" to "Приглашено",
        "balance" to "Баланс",
        "send_report" to "Отправить",
        "report_sent" to "Отчёт отправлен. Спасибо!",
        "translator" to "Переводчик",
        "text_to_translate" to "Текст для перевода",
        "translate" to "Перевести",
        "upgrade_to_premium" to "Обновить до Премиум",
        "spell_the_word" to "Напишите слово",
        "check" to "Проверить",
        "hint" to "Подсказка",
        "give_up" to "Сдаться",
        "correct" to "Правильно!",
        "wrong" to "Неправильно:",
        "practice_mistakes" to "Практика ошибок",
        "type_correct_form" to "Введите правильную форму",
        "repeat_phrase" to "Повторите:",
        "choose_scenario" to "Выберите сценарий",
        "choose_level" to "Ваш уровень",
        "start_learning" to "Начать обучение",
        "overall_progress" to "Общий прогресс",
        "total_xp" to "Всего XP",
        "streak" to "Серия",
        "days" to "дней",
        "tomorrow" to "Завтра",
        "in_3_days" to "Через 3 дня",
        "in_1_week" to "Через неделю",
        "no_review" to "Без повторения",
        "demo_mode" to "Демо (test/test)",
        "telegram_login" to "Telegram: @NERIVAapp_bot",
        "no_words_available" to "Нет доступных слов",
        "your_answer" to "Ваш ответ",
        "message" to "Сообщение",
        "easy" to "Легко",
        "good" to "Хорошо",
        "hard" to "Сложно",
        "bad" to "Плохо",
        "monthly_price" to "в месяц",
        "yearly_price" to "в год",
        "plans" to "Тарифы",
        "free_plan" to "Бесплатный",
        "premium_plan" to "Премиум",
        "platinum_plan" to "Платинум",
        "interface" to "Интерфейс"
    )
    
    // Minimal overrides for other core languages
    val langOverrides = mapOf(
        "es" to mapOf("today" to "Hoy", "sign_in" to "Iniciar sesion", "settings" to "Ajustes", "premium" to "Premium", "practice" to "Practica", "lesson" to "Leccion", "words" to "Palabras", "vocabulary" to "Vocabulario", "progress" to "Progreso", "mistakes" to "Errores", "tools" to "Herramientas", "limits" to "Limites", "leaderboard" to "Clasificacion", "awards" to "Premios", "level_test" to "Test de nivel", "referral" to "Referidos", "bug_report" to "Reportar error", "completed_lessons" to "Lecciones completadas", "log_out" to "Cerrar sesion", "interface_language" to "Idioma de interfaz", "learning_language" to "Idioma de aprendizaje", "level" to "Nivel", "daily_bonus" to "Bonus diario", "daily_activity" to "Actividad diaria", "your_progress" to "Tu progreso", "learning_lab" to "Laboratorio de aprendizaje", "ai_tutor" to "Tutor IA", "roleplay" to "Juego de roles", "shadowing" to "Escucha", "pronunciation" to "Pronunciacion", "spelling" to "Ortografia", "phrasebook" to "Notas", "offline" to "Sin conexion", "dashboard" to "Panel", "translator" to "Traductor", "translate" to "Traducir", "check" to "Verificar", "hint" to "Pista", "submit" to "Enviar", "cancel" to "Cancelar", "done" to "Listo", "loading" to "Cargando...", "demo_mode" to "Modo demo (test/test)"),
        "de" to mapOf("today" to "Heute", "sign_in" to "Anmelden", "settings" to "Einstellungen", "premium" to "Premium", "practice" to "Ubung", "lesson" to "Lektion", "words" to "Worter", "vocabulary" to "Vokabular", "progress" to "Fortschritt", "mistakes" to "Fehler", "tools" to "Werkzeuge", "limits" to "Limits", "leaderboard" to "Rangliste", "awards" to "Auszeichnungen", "level_test" to "Leveltest", "referral" to "Empfehlungen", "bug_report" to "Fehler melden", "completed_lessons" to "Abgeschlossene Lektionen", "log_out" to "Abmelden", "interface_language" to "Oberflachensprache", "learning_language" to "Lernsprache", "level" to "Level", "daily_bonus" to "Tagesbonus", "daily_activity" to "Tagliche Aktivitat", "your_progress" to "Dein Fortschritt", "ai_tutor" to "AI Tutor", "roleplay" to "Rollenspiel", "shadowing" to "Horen", "pronunciation" to "Aussprache", "spelling" to "Rechtschreibung", "check" to "Prufen", "hint" to "Hinweis", "submit" to "Absenden", "cancel" to "Abbrechen", "demo_mode" to "Demo-Modus (test/test)"),
        "fr" to mapOf("today" to "Aujourd hui", "sign_in" to "Se connecter", "settings" to "Parametres", "premium" to "Premium", "practice" to "Pratique", "lesson" to "Lecon", "words" to "Mots", "vocabulary" to "Vocabulaire", "progress" to "Progres", "mistakes" to "Erreurs", "tools" to "Outils", "limits" to "Limites", "leaderboard" to "Classement", "awards" to "Recompenses", "level_test" to "Test de niveau", "referral" to "Parrainage", "bug_report" to "Signaler un bug", "completed_lessons" to "Lecons terminees", "log_out" to "Deconnexion", "interface_language" to "Langue de l interface", "learning_language" to "Langue d apprentissage", "level" to "Niveau", "daily_bonus" to "Bonus quotidien", "daily_activity" to "Activite quotidienne", "your_progress" to "Votre progres", "ai_tutor" to "Tuteur IA", "roleplay" to "Jeu de role", "shadowing" to "Ecoute", "pronunciation" to "Prononciation", "spelling" to "Orthographe", "check" to "Verifier", "hint" to "Indice", "submit" to "Envoyer", "cancel" to "Annuler", "demo_mode" to "Mode demo (test/test)"),
        "it" to mapOf("today" to "Oggi", "sign_in" to "Accedi", "settings" to "Impostazioni", "premium" to "Premium", "practice" to "Pratica", "lesson" to "Lezione", "words" to "Parole", "vocabulary" to "Vocabolario", "progress" to "Progresso", "mistakes" to "Errori", "check" to "Verifica", "demo_mode" to "Modalita demo (test/test)"),
        "zh" to mapOf("today" to "今天", "sign_in" to "登录", "settings" to "设置", "premium" to "高级版", "practice" to "练习", "lesson" to "课程", "words" to "单词", "vocabulary" to "词汇", "progress" to "进度", "mistakes" to "错误", "tools" to "工具", "check" to "检查", "demo_mode" to "演示模式 (test/test)"),
        "ja" to mapOf("today" to "今日", "sign_in" to "ログイン", "settings" to "設定", "premium" to "プレミアム", "practice" to "練習", "lesson" to "レッスン", "words" to "単語", "vocabulary" to "語彙", "progress" to "進捗", "mistakes" to "間違い", "check" to "確認", "demo_mode" to "デモモード (test/test)"),
        "ko" to mapOf("today" to "오늘", "sign_in" to "로그인", "settings" to "설정", "premium" to "프리미엄", "practice" to "연습", "lesson" to "수업", "words" to "단어", "vocabulary" to "어휘", "progress" to "진행", "mistakes" to "오류", "check" to "확인", "demo_mode" to "데모 모드 (test/test)")
    )
    
    fun generate() {
        val resDir = File("app/src/main/res")
        
        // English (default)
        writeStringsXml(File(resDir, "values/strings.xml"), enBase)
        println("Generated values/strings.xml (${enBase.size} keys)")
        
        // Russian
        writeStringsXml(File(resDir, "values-ru/strings.xml"), ruBase)
        println("Generated values-ru/strings.xml (${ruBase.size} keys)")
        
        // Core languages with overrides
        for ((lang, overrides) in langOverrides) {
            val merged = enBase.toMutableMap()
            merged.putAll(overrides)
            writeStringsXml(File(resDir, "values-$lang/strings.xml"), merged)
            println("Generated values-$lang/strings.xml (${merged.size} keys)")
        }
        
        // Fallback languages (English)
        val fallbackLangs = listOf("tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt", "ar", "bn", "cs", "el", "hi", "hu", "id", "nl", "sv", "ta", "te", "th", "tl", "tr", "vi")
        for (lang in fallbackLangs) {
            if (langOverrides.containsKey(lang)) continue // already generated
            writeStringsXml(File(resDir, "values-$lang/strings.xml"), enBase)
            println("Generated values-$lang/strings.xml (English fallback, ${enBase.size} keys)")
        }
        
        println("\nDone! Generated strings.xml for 35 languages.")
    }
    
    private fun writeStringsXml(file: File, strings: Map<String, String>) {
        file.parentFile.mkdirs()
        val sb = StringBuilder()
        sb.appendLine("<?xml version=\"1.0\" encoding=\"utf-8\"?>")
        sb.appendLine("<resources>")
        for ((key, value) in strings.toSortedMap()) {
            val escaped = value
                .replace("&", "&amp;")
                .replace("<", "&lt;")
                .replace(">", "&gt;")
                .replace("'", "\\'")
                .replace("\"", "&quot;")
            sb.appendLine("    <string name=\"$key\">$escaped</string>")
        }
        sb.appendLine("</resources>")
        file.writeText(sb.toString(), Charsets.UTF_8)
    }
}

fun main() {
    I18nGenerator.generate()
}
