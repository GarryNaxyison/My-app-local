package ru.neriva.app.audio

import android.media.AudioAttributes
import android.media.MediaPlayer
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import ru.neriva.app.NERIVAApp
import java.io.File

/**
 * Plays TTS / pronunciation audio returned by the backend.
 *
 * The web app fetches audio blobs from `/api/tools/translator-speech` and
 * `/api/words/pronunciation`. On Android we POST JSON to the same endpoint,
 * receive the WAV bytes, write them to a cache file, and play via MediaPlayer.
 */
class AudioPlayer {

    private var player: MediaPlayer? = null

    fun stop() {
        player?.let {
            try { if (it.isPlaying) it.stop() } catch (_: Exception) {}
            it.release()
        }
        player = null
    }

    fun isPlaying(): Boolean = player?.isPlaying == true

    /** Play a local audio file (e.g. a cached WAV blob). */
    fun playFile(file: File, onCompletion: (() -> Unit)? = null) {
        stop()
        try {
            val mp = MediaPlayer().apply {
                setAudioAttributes(
                    AudioAttributes.Builder()
                        .setContentType(AudioAttributes.CONTENT_TYPE_SPEECH)
                        .setUsage(AudioAttributes.USAGE_MEDIA)
                        .build()
                )
                setDataSource(file.absolutePath)
                setOnCompletionListener {
                    it.release()
                    player = null
                    onCompletion?.invoke()
                }
                setOnErrorListener { _, _, _ ->
                    player = null
                    onCompletion?.invoke()
                    true
                }
                prepare()
                start()
            }
            player = mp
        } catch (_: Exception) {
            onCompletion?.invoke()
        }
    }

    companion object {
        /**
         * Fetch an authenticated WAV blob from the backend and play it.
         *
         * Mirrors the web `apiBlob()` helper: POSTs JSON, receives bytes,
         * then plays them. Used for `/api/tools/translator-speech` and
         * `/api/words/pronunciation`.
         */
        suspend fun playBlob(
            endpoint: String,
            jsonBody: String,
            onCompletion: (() -> Unit)? = null
        ): Unit = withContext(Dispatchers.IO) {
            try {
                val baseUrl = "https://api.neriva.ru"
                val body = jsonBody.toRequestBody("application/json; charset=utf-8".toMediaType())
                val request = Request.Builder()
                    .url("$baseUrl$endpoint")
                    .post(body)
                    .build()
                NERIVAApp.instance.sharedHttpClient.newCall(request).execute().use { resp ->
                    if (!resp.isSuccessful) {
                        onCompletion?.invoke()
                        return@use
                    }
                    val bytes = resp.body?.bytes() ?: run {
                        onCompletion?.invoke()
                        return@use
                    }
                    val tmp = File.createTempFile("tts", ".wav", NERIVAApp.instance.cacheDir)
                    tmp.writeBytes(bytes)
                    get().playFile(tmp) { tmp.delete() ; onCompletion?.invoke() }
                }
            } catch (_: Exception) {
                onCompletion?.invoke()
            }
        }

        @Volatile private var _instance: AudioPlayer? = null
        fun get(): AudioPlayer = _instance ?: synchronized(this) {
            _instance ?: AudioPlayer().also { _instance = it }
        }

        /**
         * Download an authenticated audio file from the backend (relative or
         * absolute URL) and play it. Used for AI voice replies that include an
         * `audio_url` field.
         */
        suspend fun playUrl(url: String, onCompletion: (() -> Unit)? = null): Unit = withContext(Dispatchers.IO) {
            try {
                val fullUrl = if (url.startsWith("http")) url else "https://api.neriva.ru${if (url.startsWith("/")) "" else "/"}$url"
                val request = Request.Builder().url(fullUrl).get().build()
                NERIVAApp.instance.sharedHttpClient.newCall(request).execute().use { resp ->
                    if (!resp.isSuccessful) {
                        onCompletion?.invoke()
                        return@use
                    }
                    val bytes = resp.body?.bytes() ?: run {
                        onCompletion?.invoke()
                        return@use
                    }
                    val ext = if (fullUrl.contains(".mp3", ignoreCase = true)) ".mp3" else ".wav"
                    val tmp = File.createTempFile("reply", ext, NERIVAApp.instance.cacheDir)
                    tmp.writeBytes(bytes)
                    get().playFile(tmp) { tmp.delete(); onCompletion?.invoke() }
                }
            } catch (_: Exception) {
                onCompletion?.invoke()
            }
        }
    }
}
