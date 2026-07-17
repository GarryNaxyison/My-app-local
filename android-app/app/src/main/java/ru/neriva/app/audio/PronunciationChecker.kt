package ru.neriva.app.audio

import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.asRequestBody
import ru.neriva.app.NERIVAApp
import java.io.File

class PronunciationChecker {

    data class Result(
        val score: Int?,
        val feedback: String?,
        val weakWords: List<String>?,
        val xp: Int?
    )

    suspend fun check(audioFile: File, sessionId: String? = null): Result {
        val app = NERIVAApp.instance
        val requestFile = audioFile.asRequestBody("audio/wav".toMediaType())
        val body = MultipartBody.Builder()
            .setType(MultipartBody.FORM)
            .addFormDataPart("session_id", sessionId ?: "")
            .addFormDataPart("audio", audioFile.name, requestFile)
            .build()

        // Use the practice endpoint for pronunciation checking
        val response = app.apiClient.pronunciationCheck(
            ru.neriva.app.data.api.PronunciationCheckRequest(
                sessionId = sessionId,
                audio = audioFile.readBytes().let { android.util.Base64.encodeToString(it, android.util.Base64.NO_WRAP) }
            )
        )
        return Result(
            score = response.score,
            feedback = response.feedback,
            weakWords = response.weakWords,
            xp = response.xp
        )
    }
}
