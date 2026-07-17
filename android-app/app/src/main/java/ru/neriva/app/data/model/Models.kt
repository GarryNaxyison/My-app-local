package ru.neriva.app.data.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

// ---- Session / Profile ----

@Serializable
data class SessionData(
    val authenticated: Boolean = false,
    val account: AccountInfo? = null,
    val user: UserProfile? = null,
    val copy: Map<String, String> = emptyMap(),
    val tool_copy: Map<String, String> = emptyMap(),
    val system_copy: Map<String, String> = emptyMap(),
    val premium_copy: Map<String, String> = emptyMap(),
    val interface_languages: List<LanguageOption>? = null,
    val learning_languages: List<LanguageOption>? = null,
    val premium_plans: List<PremiumPlan>? = null,
    val yookassa_enabled: Boolean = false,
    val telegram_login_bot: String? = null,
    val web_app_url: String? = null,
    val captcha: CaptchaInfo? = null,
)

@Serializable
data class AccountInfo(
    val login: String? = null,
    val password_set: Boolean = false,
    val profile_ready: Boolean = false,
)

@Serializable
data class CaptchaInfo(
    val enabled: Boolean = false,
    val site_key: String? = null,
    val provider: String? = null,
)

@Serializable
data class UserProfile(
    @SerialName("interface_language") val interfaceLanguage: String? = null,
    @SerialName("learning_language") val learningLanguage: String? = null,
    @SerialName("telegram_linked") val telegramLinked: Boolean = false,
    @SerialName("telegram_account") val telegramAccount: TelegramAccount? = null,
    @SerialName("created_at") val createdAt: String? = null,
    val level: String? = null,
    @SerialName("learning_focus") val learningFocus: String? = null,
    val plan: String? = null,
    val premium: Boolean = false,
    @SerialName("premium_until") val premiumUntil: String? = null,
    @SerialName("referral_code") val referralCode: String? = null,
    @SerialName("referral_count") val referralCount: Int = 0,
    @SerialName("referral_balance") val referralBalance: String? = null,
    @SerialName("referral_withdraw_min") val referralWithdrawMin: String? = null,
    @SerialName("referral_invitees") val referralInvitees: List<ReferralInvitee> = emptyList(),
    val xp: Int = 0,
    @SerialName("xp_level") val xpLevel: Int = 0,
    @SerialName("xp_title") val xpTitle: String? = null,
    @SerialName("xp_current") val xpCurrent: Int = 0,
    @SerialName("xp_needed") val xpNeeded: Int = 0,
    @SerialName("lessons_today") val lessonsToday: Int = 0,
    @SerialName("lesson_limit") val lessonLimit: Int = 0,
    @SerialName("practice_today") val practiceToday: Int = 0,
    @SerialName("practice_limit") val practiceLimit: Int = 0,
    @SerialName("voice_today") val voiceToday: Int = 0,
    @SerialName("voice_limit") val voiceLimit: Int = 0,
    @SerialName("daily_bonus_claims") val dailyBonusClaims: List<String> = emptyList(),
    @SerialName("active_lesson_prompt") val activeLessonPrompt: String? = null,
    @SerialName("active_lesson_instruction") val activeLessonInstruction: String? = null,
    val phrasebook: List<PhrasebookItem> = emptyList(),
    @SerialName("lesson_count") val lessonCount: Int = 0,
    @SerialName("practice_count") val practiceCount: Int = 0,
    @SerialName("word_game_count") val wordGameCount: Int = 0,
    @SerialName("learned_words") val learnedWords: Int = 0,
    val mistakes: Int = 0,
)

@Serializable
data class TelegramAccount(
    val id: Long? = null,
    val name: String? = null,
    val label: String? = null,
)

@Serializable
data class LanguageOption(
    val code: String = "",
    val name: String? = null,
    @SerialName("native_name") val nativeName: String? = null,
    @SerialName("interface_name") val interfaceName: String? = null,
)

@Serializable
data class ReferralInvitee(
    val id: Long? = null,
    val name: String? = null,
    val level: String? = null,
    @SerialName("xp_level") val xpLevel: Int = 0,
    val xp: Int = 0,
    @SerialName("reached_level_3") val reachedLevel3: Boolean = false,
    val earned: String? = null,
    @SerialName("joined_at") val joinedAt: String? = null,
)

// ---- Premium ----

@Serializable
data class PremiumPlan(
    val product: String = "",
    val tier: String? = null,
    val current: Boolean = false,
    val active: Boolean = false,
    @SerialName("is_current") val isCurrent: Boolean = false,
    val title: String? = null,
    val label: String? = null,
    val body: String? = null,
    val features: List<String> = emptyList(),
    @SerialName("locked_features") val lockedFeatures: List<String> = emptyList(),
    val note: String? = null,
    @SerialName("days_label") val daysLabel: String? = null,
    @SerialName("rub_price") val rubPrice: String? = null,
    @SerialName("stars_price") val starsPrice: String? = null,
)

// ---- AI Tutor ----

@Serializable
data class AiTutorResponse(
    val session: AiTutorSession? = null,
    val lesson: LessonData? = null,
    @SerialName("current_stage") val currentStage: String? = null,
    @SerialName("next_step") val nextStep: AiTutorStep? = null,
    val feedback: FeedbackInfo? = null,
    val xp: Int? = null,
    @SerialName("reward_xp") val rewardXp: Int? = null,
    @SerialName("reward_title") val rewardTitle: String? = null,
)

@Serializable
data class AiTutorSession(
    val id: String? = null,
    @SerialName("current_stage") val currentStage: String? = null,
    val status: String? = null,
)

@Serializable
data class AiTutorStep(
    val stage: String = "",
    val kind: String = "",
    val title: String? = null,
    val instruction: String? = null,
    val lesson: LessonData? = null,
    val word: WordData? = null,
    val question: QuestionData? = null,
    val options: List<StepOption>? = null,
)

@Serializable
data class LessonData(
    val title: String? = null,
    val level: String? = null,
    val theme: String? = null,
    @SerialName("lesson_goal") val lessonGoal: String? = null,
    val story: StoryData? = null,
    val words: List<WordData> = emptyList(),
    @SerialName("comprehension_questions") val comprehensionQuestions: List<QuestionData> = emptyList(),
    @SerialName("production_task") val productionTask: ProductionTask? = null,
)

@Serializable
data class StoryData(
    @SerialName("title_interface") val titleInterface: String? = null,
    @SerialName("story_title") val storyTitle: String? = null,
    @SerialName("text_target") val textTarget: String? = null,
    @SerialName("audio_text_target") val audioTextTarget: String? = null,
)

@Serializable
data class WordData(
    val id: String? = null,
    val target: String? = null,
    @SerialName("interface_translation") val interfaceTranslation: String? = null,
    @SerialName("example_sentence_target") val exampleSentenceTarget: String? = null,
    @SerialName("audio_text_target") val audioTextTarget: String? = null,
)

@Serializable
data class QuestionData(
    val id: String? = null,
    @SerialName("question_target") val questionTarget: String? = null,
)

@Serializable
data class ProductionTask(
    @SerialName("instruction_interface") val instructionInterface: String? = null,
    @SerialName("required_word_count") val requiredWordCount: Int? = null,
    @SerialName("sentence_count") val sentenceCount: String? = null,
)

@Serializable
data class StepOption(
    val id: String? = null,
    val text: String? = null,
    val correct: Boolean = false,
)

@Serializable
data class FeedbackInfo(
    val ok: Boolean? = null,
    val message: String? = null,
)

// ---- Vocabulary / Words ----

@Serializable
data class WordChallenge(
    val empty: Boolean = false,
    val message: String? = null,
    val prompt: String? = null,
    val context: String? = null,
    val direction: String? = null,
    val options: List<ChoiceOption> = emptyList(),
    val instruction: String? = null,
    @SerialName("correct_answer_id") val correctAnswerId: String? = null,
    @SerialName("word_id") val wordId: String? = null,
    val word: String? = null,
    val translation: String? = null,
    val reportable: Boolean = false,
)

@Serializable
data class ChoiceOption(
    val id: String = "",
    val text: String = "",
)

@Serializable
data class SpellingChallenge(
    val empty: Boolean = false,
    val message: String? = null,
    @SerialName("word_id") val wordId: String? = null,
    val prompt: String? = null,
    val context: String? = null,
    val direction: String? = null,
)

@Serializable
data class VocabularyItem(
    val id: String = "",
    val word: String = "",
    val translation: String? = null,
    val context: String? = null,
    val example: String? = null,
    @SerialName("review_correct_count") val reviewCorrectCount: Int = 0,
    @SerialName("spelling_correct_count") val spellingCorrectCount: Int = 0,
    @SerialName("training_count") val trainingCount: Int = 0,
)

// ---- Mistakes ----

@Serializable
data class MistakeItem(
    val index: Int? = null,
    val word: String? = null,
    val correction: String? = null,
    val context: String? = null,
    val explanation: String? = null,
    @SerialName("created_at") val createdAt: String? = null,
    val count: Int? = null,
)

// ---- Phrasebook ----

@Serializable
data class PhrasebookItem(
    val id: String = "",
    val phrase: String = "",
    val translation: String? = null,
    val note: String? = null,
    val source: String? = null,
    val language: String? = null,
    @SerialName("created_at") val createdAt: String? = null,
)

// ---- Leaderboard ----

@Serializable
data class LeaderboardEntry(
    val name: String? = null,
    val score: Int? = null,
    @SerialName("rating_points") val ratingPoints: Int? = null,
    val xp: Int? = null,
    val words: Int? = null,
    val mistakes: Int? = null,
    val languages: List<String>? = null,
    val level: String? = null,
    val title: String? = null,
)

// ---- Level Test ----

@Serializable
data class LevelQuestion(
    val index: Int = 0,
    val total: Int = 0,
    val question: String = "",
    val options: List<String> = emptyList(),
    val score: Int? = null,
)

// ---- Translator ----

@Serializable
data class TranslatorResult(
    @SerialName("source_text") val sourceText: String? = null,
    val translation: String? = null,
    val result: String? = null,
    @SerialName("source_language") val sourceLanguage: String? = null,
    @SerialName("target_language") val targetLanguage: String? = null,
)

// ---- Practice ----

@Serializable
data class PracticeResponse(
    val message: String? = null,
    val correction: String? = null,
    @SerialName("audio_url") val audioUrl: String? = null,
    val xp: Int? = null,
)

// ---- API error ----

@Serializable
data class ApiError(
    val error: String? = null,
    val message: String? = null,
)
