package ru.neriva.app.data.repo

import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.RequestBody.Companion.asRequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import ru.neriva.app.data.api.*
import ru.neriva.app.data.db.PhrasebookDao
import ru.neriva.app.data.db.PhrasebookEntity
import ru.neriva.app.data.model.*
import android.content.SharedPreferences

class AuthRepository(
    private val api: NerivaApi,
    private val prefs: SharedPreferences
) {
    suspend fun login(login: String, password: String, captchaToken: String? = null): SessionData {
        val resp = api.login(LoginRequest(login = login, password = password, captchaToken = captchaToken))
        prefs.edit().putBoolean("logged_in", true).apply()
        return resp
    }

    suspend fun register(login: String, password: String, passwordConfirm: String, referralCode: String? = null, captchaToken: String? = null): SessionData {
        val resp = api.register(
            RegisterRequest(
                login = login, password = password,
                passwordConfirm = passwordConfirm, referralCode = referralCode,
                captchaToken = captchaToken
            )
        )
        prefs.edit().putBoolean("logged_in", true).apply()
        return resp
    }

    suspend fun logout() {
        try { api.logout() } catch (_: Exception) {}
        prefs.edit().putBoolean("logged_in", false).apply()
    }

    suspend fun changePassword(oldPassword: String?, newPassword: String, confirmPassword: String) {
        api.changePassword(PasswordRequest(oldPassword, newPassword, confirmPassword))
    }

    val isLoggedIn: Boolean get() = prefs.getBoolean("logged_in", false)
}

class SessionRepository(private val api: NerivaApi) {
    suspend fun getSession(): SessionData = api.getSession()

    suspend fun updateSettings(
        interfaceLanguage: String? = null,
        learningLanguage: String? = null,
        level: String? = null,
        learningFocus: String? = null
    ): SessionData = api.updateSettings(
        SettingsRequest(
            interfaceLanguage = interfaceLanguage,
            learningLanguage = learningLanguage,
            level = level,
            learningFocus = learningFocus
        )
    )

    suspend fun claimDailyBonus() = api.claimDailyBonus()

    suspend fun getProgress() = api.getProgress()

    suspend fun getLeaderboard(language: String? = null) = api.getLeaderboard(language)

    suspend fun getPremiumPlans() = api.getPremiumPlans()

    suspend fun createPayment(tier: String, period: String, method: String) =
        api.createPayment(PaymentRequest(tier, period, method))

    suspend fun activateKey(key: String) = api.activateKey(ActivationKeyRequest(key))
}

class TutorRepository(private val api: NerivaApi) {
    suspend fun start() = api.aiTutorStart()
    suspend fun session() = api.aiTutorSession()
    suspend fun restart() = api.aiTutorRestart()
    suspend fun answer(answer: String, sessionId: String? = null, wordId: String? = null) =
        api.aiTutorAnswer(AnswerRequest(answer, sessionId, wordId))
    suspend fun finish() = api.aiTutorFinish()
    suspend fun completed(page: Int = 1) = api.aiTutorCompleted(page)
}

class VocabularyRepository(private val api: NerivaApi) {
    suspend fun getList(page: Int = 1, pageSize: Int = 20) = api.vocabulary(page, pageSize)
    suspend fun nextWord() = api.wordsNext()
    suspend fun answerWord(answer: String, sessionId: String? = null, wordId: String? = null) =
        api.wordsAnswer(AnswerRequest(answer, sessionId, wordId))
    suspend fun reportWord(wordId: String, action: String = "wrong") = api.wordsReport(WordReportRequest(wordId, action))
    suspend fun nextGame() = api.wordGameNext()
    suspend fun answerGame(answer: String, sessionId: String? = null, wordId: String? = null) =
        api.wordGameAnswer(AnswerRequest(answer, sessionId, wordId))
    suspend fun startSpelling() = api.spellingStart()
    suspend fun answerSpelling(answer: String, sessionId: String? = null, wordId: String? = null) =
        api.spellingAnswer(AnswerRequest(answer, sessionId, wordId))
}

class PracticeRepository(private val api: NerivaApi) {
    suspend fun send(message: String? = null, sessionId: String? = null, action: String? = null) =
        api.practice(PracticeRequest(message, sessionId, action))

    /** Upload a recorded voice message (WAV/AAC) and receive the AI reply. */
    suspend fun sendVoice(audioFile: java.io.File, message: String? = null, sessionId: String? = null): PracticeMessage {
        val reqFile = audioFile.asRequestBody()
        val voicePart = okhttp3.MultipartBody.Part.createFormData(
            "voice", audioFile.name, reqFile
        )
        val messagePart = message?.takeIf { it.isNotBlank() }
            ?.toRequestBody("text/plain".toMediaTypeOrNull())
        val sessionPart = sessionId?.takeIf { it.isNotBlank() }
            ?.toRequestBody("text/plain".toMediaTypeOrNull())
        return api.practiceVoice(voicePart, messagePart, sessionPart)
    }

    suspend fun startShadowing() = api.shadowingStart()
    suspend fun answerShadowing(answer: String, sessionId: String? = null) =
        api.shadowingAnswer(AnswerRequest(answer, sessionId))

    /** Upload a recording for the shadowing "listen & repeat" exercise. */
    suspend fun answerShadowingVoice(audioFile: java.io.File, target: String? = null, text: String? = null): ShadowingResponse {
        val reqFile = audioFile.asRequestBody()
        val voicePart = okhttp3.MultipartBody.Part.createFormData(
            "voice", audioFile.name, reqFile
        )
        val targetPart = target?.takeIf { it.isNotBlank() }?.toRequestBody("text/plain".toMediaTypeOrNull())
        val textPart = text?.takeIf { it.isNotBlank() }?.toRequestBody("text/plain".toMediaTypeOrNull())
        return api.shadowingAnswerVoice(voicePart, targetPart, textPart)
    }

    suspend fun startPronunciation() = api.pronunciationStart()
    suspend fun checkPronunciation(sessionId: String? = null, audio: String? = null) =
        api.pronunciationCheck(PronunciationCheckRequest(sessionId, audio))

    /** Upload a recording and receive a pronunciation score. */
    suspend fun checkPronunciationVoice(audioFile: java.io.File, sessionId: String? = null): PronunciationCheckResponse {
        val reqFile = audioFile.asRequestBody()
        val voicePart = okhttp3.MultipartBody.Part.createFormData(
            "voice", audioFile.name, reqFile
        )
        val sessionPart = sessionId?.takeIf { it.isNotBlank() }
            ?.toRequestBody("text/plain".toMediaTypeOrNull())
        return api.pronunciationCheckVoice(voicePart, sessionPart)
    }

    suspend fun startLesson() = api.lessonStart()
    suspend fun answerLesson(answer: String, sessionId: String? = null) =
        api.lessonAnswer(AnswerRequest(answer, sessionId))
}

class MistakesRepository(private val api: NerivaApi) {
    suspend fun getAll() = api.getMistakes()
    suspend fun startPractice() = api.mistakePracticeStart()
    suspend fun answerPractice(answer: String, sessionId: String? = null) =
        api.mistakePracticeAnswer(AnswerRequest(answer, sessionId))
    suspend fun delete(id: String) = api.deleteMistake(IdRequest(id))
    suspend fun clearAll() = api.clearMistakes()
}

class PhrasebookRepository(private val dao: PhrasebookDao) {
    suspend fun getCached(): List<PhrasebookEntity> = dao.getAll()
    suspend fun insert(item: PhrasebookEntity) = dao.insert(item)
    suspend fun delete(id: String) = dao.delete(id)
    fun observeAll() = dao.observeAll()
}

class ToolsRepository(private val api: NerivaApi) {
    suspend fun translate(text: String, sourceLanguage: String? = null, targetLanguage: String? = null) =
        api.translate(TranslateRequest(text, sourceLanguage, targetLanguage))
    suspend fun translateSpeech(text: String, sourceLanguage: String? = null) =
        api.translateSpeech(TranslateRequest(text, sourceLanguage))
    suspend fun voiceToText(text: String, sourceLanguage: String? = null) =
        api.voiceToText(VoiceTextRequest(text, sourceLanguage))

    /** Transcribe an uploaded voice file and translate it (tools voice-to-text mode). */
    suspend fun voiceToTextFile(audioFile: java.io.File, targetLanguage: String? = null): TranslatorResult {
        val reqFile = audioFile.asRequestBody()
        val voicePart = okhttp3.MultipartBody.Part.createFormData("voice", audioFile.name, reqFile)
        val langPart = targetLanguage?.toRequestBody("text/plain".toMediaTypeOrNull())
        return api.voiceToTextMultipart(voicePart, langPart)
    }

    /** Translate text from an uploaded image (tools photo-translation mode). */
    suspend fun imageTranslate(imageFile: java.io.File, targetLanguage: String): TranslatorResult {
        val reqFile = imageFile.asRequestBody()
        val imagePart = okhttp3.MultipartBody.Part.createFormData("image", imageFile.name, reqFile)
        return api.imageTranslate(imagePart, targetLanguage)
    }
}
