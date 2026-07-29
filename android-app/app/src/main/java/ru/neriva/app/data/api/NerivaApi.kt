package ru.neriva.app.data.api

import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.toRequestBody
import retrofit2.http.*
import ru.neriva.app.data.model.*

// ---- API Service interface ----

interface NerivaApi {

    // Auth
    @POST("api/auth/register")
    suspend fun register(@Body body: RegisterRequest): SessionData

    @POST("api/auth/login")
    suspend fun login(@Body body: LoginRequest): SessionData

    @POST("api/auth/logout")
    suspend fun logout(): OkResponse

    @POST("api/auth/telegram/start")
    suspend fun telegramStart(): TelegramStartResponse

    @POST("api/auth/telegram/status")
    suspend fun telegramStatus(@Body body: TelegramStatusRequest): SessionData

    @POST("api/auth/password")
    suspend fun changePassword(@Body body: PasswordRequest): OkResponse

    // Session
    @GET("api/session")
    suspend fun getSession(): SessionData

    @POST("api/settings")
    suspend fun updateSettings(@Body body: SettingsRequest): SessionData

    @POST("api/navigation-layout")
    suspend fun saveNavigationLayout(@Body body: NavigationLayoutRequest): OkResponse

    @POST("api/daily/claim")
    suspend fun claimDailyBonus(): DailyBonusResponse

    // AI Tutor
    @POST("api/ai-tutor/start")
    suspend fun aiTutorStart(): AiTutorResponse

    @POST("api/ai-tutor/session")
    suspend fun aiTutorSession(): AiTutorResponse

    @POST("api/ai-tutor/restart")
    suspend fun aiTutorRestart(): AiTutorResponse

    @POST("api/ai-tutor/answer")
    suspend fun aiTutorAnswer(@Body body: AnswerRequest): AiTutorResponse

    @POST("api/ai-tutor/finish")
    suspend fun aiTutorFinish(): AiTutorResponse

    @GET("api/ai-tutor/completed")
    suspend fun aiTutorCompleted(@Query("page") page: Int = 1, @Query("limit") limit: Int = 20): CompletedLessonsResponse

    // Lesson
    @POST("api/lesson/start")
    suspend fun lessonStart(): LessonResponse

    @POST("api/lesson/answer")
    suspend fun lessonAnswer(@Body body: AnswerRequest): LessonResponse

    // Practice
    @POST("api/practice")
    suspend fun practice(@Body body: PracticeRequest): PracticeMessage

    @Multipart
    @POST("api/practice")
    suspend fun practiceVoice(
        @Part voice: okhttp3.MultipartBody.Part,
        @Part("message") message: okhttp3.RequestBody? = null,
        @Part("session_id") sessionId: okhttp3.RequestBody? = null,
    ): PracticeMessage

    // Shadowing
    @POST("api/shadowing/start")
    suspend fun shadowingStart(): ShadowingResponse

    @POST("api/shadowing/answer")
    suspend fun shadowingAnswer(@Body body: AnswerRequest): ShadowingResponse

    @Multipart
    @POST("api/shadowing/answer")
    suspend fun shadowingAnswerVoice(
        @Part voice: okhttp3.MultipartBody.Part,
        @Part("target") target: okhttp3.RequestBody? = null,
        @Part("text") text: okhttp3.RequestBody? = null,
    ): ShadowingResponse

    // Pronunciation
    @POST("api/pronunciation/start")
    suspend fun pronunciationStart(): PronunciationResponse

    @POST("api/pronunciation/check")
    suspend fun pronunciationCheck(@Body body: PronunciationCheckRequest): PronunciationCheckResponse

    @Multipart
    @POST("api/pronunciation/check")
    suspend fun pronunciationCheckVoice(
        @Part voice: okhttp3.MultipartBody.Part,
        @Part("session_id") sessionId: okhttp3.RequestBody? = null,
    ): PronunciationCheckResponse

    // Words
    @POST("api/words/next")
    suspend fun wordsNext(): WordChallenge

    @POST("api/words/answer")
    suspend fun wordsAnswer(@Body body: AnswerRequest): WordAnswerResponse

    @POST("api/words/report")
    suspend fun wordsReport(@Body body: WordReportRequest): OkResponse

    // Word Game
    @POST("api/word-game/next")
    suspend fun wordGameNext(): WordChallenge

    @POST("api/word-game/answer")
    suspend fun wordGameAnswer(@Body body: AnswerRequest): WordAnswerResponse

    // Spelling
    @POST("api/spelling/start")
    suspend fun spellingStart(): SpellingChallenge

    @POST("api/spelling/answer")
    suspend fun spellingAnswer(@Body body: AnswerRequest): SpellingAnswerResponse

    // Vocabulary
    @GET("api/vocabulary")
    suspend fun vocabulary(@Query("page") page: Int = 1, @Query("page_size") pageSize: Int = 20): VocabularyResponse

    // Phrasebook
    @GET("api/phrasebook")
    suspend fun getPhrasebook(): PhrasebookResponse

    @POST("api/phrasebook")
    suspend fun savePhrase(@Body body: PhraseSaveRequest): OkResponse

    @DELETE("api/phrasebook")
    suspend fun deletePhrase(@Query("id") id: String): OkResponse

    // Mistakes
    @GET("api/mistakes")
    suspend fun getMistakes(): MistakesResponse

    @POST("api/mistakes/practice/start")
    suspend fun mistakePracticeStart(): MistakePracticeResponse

    @POST("api/mistakes/practice/answer")
    suspend fun mistakePracticeAnswer(@Body body: AnswerRequest): MistakePracticeResponse

    @POST("api/mistakes/delete")
    suspend fun deleteMistake(@Body body: IdRequest): OkResponse

    @POST("api/mistakes/clear")
    suspend fun clearMistakes(): OkResponse

    // Tools
    @POST("api/tools/translator")
    suspend fun translate(@Body body: TranslateRequest): TranslatorResult

    @POST("api/tools/translator-speech")
    suspend fun translateSpeech(@Body body: TranslateRequest): TranslatorResult

    @POST("api/tools/voice-text")
    suspend fun voiceToText(@Body body: VoiceTextRequest): VoiceTextResponse

    @Multipart
    @POST("api/tools/voice-text")
    suspend fun voiceToTextMultipart(
        @Part voice: okhttp3.MultipartBody.Part,
        @Part("language") language: okhttp3.RequestBody? = null,
    ): TranslatorResult

    @Multipart
    @POST("api/tools/image-translate")
    suspend fun imageTranslate(
        @Part image: okhttp3.MultipartBody.Part,
        @Part("language") language: okhttp3.RequestBody,
    ): TranslatorResult

    // Progress
    @GET("api/progress")
    suspend fun getProgress(): ProgressResponse

    @GET("api/leaderboard")
    suspend fun getLeaderboard(@Query("language") language: String? = null): LeaderboardResponse

    // Premium
    @GET("api/premium/plans")
    suspend fun getPremiumPlans(): PremiumPlansResponse

    @POST("api/premium/payment")
    suspend fun createPayment(@Body body: PaymentRequest): PaymentResponse

    @POST("api/premium/activation-key")
    suspend fun activateKey(@Body body: ActivationKeyRequest): SessionData

    // Bug report
    @Multipart
    @POST("api/bug-report")
    suspend fun bugReport(
        @Part("description") description: okhttp3.RequestBody,
        @Part screenshot: okhttp3.MultipartBody.Part?
    ): OkResponse

    // Level test
    @POST("api/level-test/start")
    suspend fun levelTestStart(): LevelTestStartResponse

    @POST("api/level-test/answer")
    suspend fun levelTestAnswer(@Body body: AnswerRequest): LevelTestAnswerResponse
}

// ---- Request / Response DTOs ----

@kotlinx.serialization.Serializable
data class RegisterRequest(
    val login: String,
    val password: String,
    @kotlinx.serialization.SerialName("password_confirm") val passwordConfirm: String,
    @kotlinx.serialization.SerialName("referral_code") val referralCode: String? = null,
    @kotlinx.serialization.SerialName("captcha_token") val captchaToken: String? = null,
)

@kotlinx.serialization.Serializable
data class LoginRequest(
    val login: String,
    val password: String,
    @kotlinx.serialization.SerialName("captcha_token") val captchaToken: String? = null,
)

@kotlinx.serialization.Serializable
data class TelegramStartResponse(
    val token: String,
    val status: String,
    @kotlinx.serialization.SerialName("bot_url") val botUrl: String,
    @kotlinx.serialization.SerialName("expires_at") val expiresAt: String,
)

@kotlinx.serialization.Serializable
data class TelegramStatusRequest(
    val token: String,
    val code: String? = null,
)

@kotlinx.serialization.Serializable
data class PasswordRequest(
    @kotlinx.serialization.SerialName("old_password") val oldPassword: String? = null,
    val password: String,
    @kotlinx.serialization.SerialName("password_confirm") val passwordConfirm: String,
)

@kotlinx.serialization.Serializable
data class SettingsRequest(
    @kotlinx.serialization.SerialName("interface_language") val interfaceLanguage: String? = null,
    @kotlinx.serialization.SerialName("learning_language") val learningLanguage: String? = null,
    val level: String? = null,
    @kotlinx.serialization.SerialName("learning_focus") val learningFocus: String? = null,
)

@kotlinx.serialization.Serializable
data class NavigationLayoutRequest(
    @kotlinx.serialization.SerialName("function_ribbon") val functionRibbon: List<String>? = null,
    @kotlinx.serialization.SerialName("mobile_pinned") val mobilePinned: List<String>? = null,
    @kotlinx.serialization.SerialName("mobile_more") val mobileMore: List<String>? = null,
    @kotlinx.serialization.SerialName("mobile_rail") val mobileRail: List<String>? = null,
)

@kotlinx.serialization.Serializable
data class AnswerRequest(
    val answer: String,
    @kotlinx.serialization.SerialName("session_id") val sessionId: String? = null,
    @kotlinx.serialization.SerialName("word_id") val wordId: String? = null,
)

@kotlinx.serialization.Serializable
data class PracticeRequest(
    val message: String? = null,
    @kotlinx.serialization.SerialName("session_id") val sessionId: String? = null,
    val action: String? = null,
)

@kotlinx.serialization.Serializable
data class PronunciationCheckRequest(
    @kotlinx.serialization.SerialName("session_id") val sessionId: String? = null,
    val audio: String? = null,
)

@kotlinx.serialization.Serializable
data class TranslateRequest(
    val text: String,
    @kotlinx.serialization.SerialName("source_language") val sourceLanguage: String? = null,
    @kotlinx.serialization.SerialName("target_language") val targetLanguage: String? = null,
)

@kotlinx.serialization.Serializable
data class VoiceTextRequest(
    val text: String,
    @kotlinx.serialization.SerialName("source_language") val sourceLanguage: String? = null,
)

@kotlinx.serialization.Serializable
data class PhraseSaveRequest(
    val phrase: String,
    val translation: String? = null,
    val note: String? = null,
    val source: String? = null,
    val language: String? = null,
    @kotlinx.serialization.SerialName("phrase_id") val phraseId: String? = null,
    val action: String = "save",
)

@kotlinx.serialization.Serializable
data class WordReportRequest(
    @kotlinx.serialization.SerialName("word_id") val wordId: String,
    val action: String = "wrong",
)

@kotlinx.serialization.Serializable
data class IdRequest(val id: String)

@kotlinx.serialization.Serializable
data class PaymentRequest(
    val tier: String,
    val period: String,
    val method: String,
)

@kotlinx.serialization.Serializable
data class ActivationKeyRequest(
    val key: String,
)

// ---- Simple response wrappers ----

@kotlinx.serialization.Serializable
data class OkResponse(val ok: Boolean = true)

@kotlinx.serialization.Serializable
data class DailyBonusResponse(val ok: Boolean = true, val xp: Int? = null)

@kotlinx.serialization.Serializable
data class LessonResponse(
    val message: String? = null,
    val prompt: String? = null,
    val instruction: String? = null,
    @kotlinx.serialization.SerialName("audio_url") val audioUrl: String? = null,
    @kotlinx.serialization.SerialName("model_url") val modelUrl: String? = null,
    val xp: Int? = null,
)

@kotlinx.serialization.Serializable
data class PracticeMessage(
    val message: String? = null,
    val correction: String? = null,
    @kotlinx.serialization.SerialName("audio_url") val audioUrl: String? = null,
    val xp: Int? = null,
)

@kotlinx.serialization.Serializable
data class ShadowingResponse(
    val text: String? = null,
    @kotlinx.serialization.SerialName("audio_url") val audioUrl: String? = null,
    val xp: Int? = null,
)

@kotlinx.serialization.Serializable
data class PronunciationResponse(
    val text: String? = null,
    @kotlinx.serialization.SerialName("audio_url") val audioUrl: String? = null,
    val xp: Int? = null,
)

@kotlinx.serialization.Serializable
data class PronunciationCheckResponse(
    val score: Int? = null,
    val feedback: String? = null,
    @kotlinx.serialization.SerialName("weak_words") val weakWords: List<String>? = null,
    val xp: Int? = null,
)

@kotlinx.serialization.Serializable
data class WordAnswerResponse(
    val correct: Boolean = false,
    @kotlinx.serialization.SerialName("correct_answer") val correctAnswer: String? = null,
    @kotlinx.serialization.SerialName("correct_translation") val correctTranslation: String? = null,
    val xp: Int? = null,
)

@kotlinx.serialization.Serializable
data class SpellingAnswerResponse(
    val correct: Boolean = false,
    @kotlinx.serialization.SerialName("correct_answer") val correctAnswer: String? = null,
    val xp: Int? = null,
)

@kotlinx.serialization.Serializable
data class VocabularyResponse(
    val items: List<VocabularyItem>? = null,
    val page: Int? = null,
    @kotlinx.serialization.SerialName("total") val total: Int? = null,
    @kotlinx.serialization.SerialName("total_pages") val totalPages: Int? = null,
)

@kotlinx.serialization.Serializable
data class PhrasebookResponse(
    val items: List<PhrasebookItem>? = null,
)

@kotlinx.serialization.Serializable
data class MistakesResponse(
    val items: List<MistakeItem>? = null,
    val total: Int? = null,
)

@kotlinx.serialization.Serializable
data class MistakePracticeResponse(
    val word: String? = null,
    val correction: String? = null,
    val context: String? = null,
    val xp: Int? = null,
    val empty: Boolean = false,
)

@kotlinx.serialization.Serializable
data class VoiceTextResponse(
    val text: String? = null,
    val translation: String? = null,
)

@kotlinx.serialization.Serializable
data class ProgressResponse(
    @kotlinx.serialization.SerialName("lesson_count") val lessonCount: Int = 0,
    @kotlinx.serialization.SerialName("practice_count") val practiceCount: Int = 0,
    @kotlinx.serialization.SerialName("word_game_count") val wordGameCount: Int = 0,
    @kotlinx.serialization.SerialName("learned_words") val learnedWords: Int = 0,
    val mistakes: Int = 0,
    val xp: Int = 0,
    @kotlinx.serialization.SerialName("xp_level") val xpLevel: Int = 0,
    @kotlinx.serialization.SerialName("xp_title") val xpTitle: String? = null,
    val streak: Int? = null,
    @kotlinx.serialization.SerialName("premium") val premium: Boolean = false,
)

@kotlinx.serialization.Serializable
data class LeaderboardResponse(
    val items: List<LeaderboardEntry>? = null,
)

@kotlinx.serialization.Serializable
data class PremiumPlansResponse(
    val plans: List<PremiumPlan>? = null,
    @kotlinx.serialization.SerialName("yookassa_enabled") val yookassaEnabled: Boolean = false,
)

@kotlinx.serialization.Serializable
data class PaymentResponse(
    val url: String? = null,
    @kotlinx.serialization.SerialName("invoice_id") val invoiceId: String? = null,
)

@kotlinx.serialization.Serializable
data class CompletedLessonsResponse(
    val items: List<CompletedLesson>? = null,
    val total: Int? = null,
)

@kotlinx.serialization.Serializable
data class CompletedLesson(
    @kotlinx.serialization.SerialName("session_id") val sessionId: String? = null,
    @kotlinx.serialization.SerialName("lesson_id") val lessonId: String? = null,
    val title: String? = null,
    val topic: String? = null,
    val level: String? = null,
    @kotlinx.serialization.SerialName("completed_at") val completedAt: String? = null,
)

@kotlinx.serialization.Serializable
data class LevelTestStartResponse(
    val question: LevelQuestion? = null,
    @kotlinx.serialization.SerialName("total_questions") val totalQuestions: Int = 0,
)

@kotlinx.serialization.Serializable
data class LevelTestAnswerResponse(
    val question: LevelQuestion? = null,
    val score: Int? = null,
    val level: String? = null,
    val finished: Boolean = false,
)
