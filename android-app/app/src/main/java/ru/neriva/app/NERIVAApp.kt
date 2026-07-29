package ru.neriva.app

import android.app.Application
import android.content.Context
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey
import okhttp3.OkHttpClient
import ru.neriva.app.data.api.NerivaApi
import ru.neriva.app.data.api.NerivaApiClient
import ru.neriva.app.data.api.PersistentCookieJar
import ru.neriva.app.data.db.NerivaDatabase
import ru.neriva.app.data.repo.AuthRepository
import ru.neriva.app.data.repo.PhrasebookRepository
import ru.neriva.app.data.repo.SessionRepository
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers

class NERIVAApp : Application() {

    lateinit var apiClient: NerivaApi
    lateinit var authRepo: AuthRepository
    lateinit var sessionRepo: SessionRepository
    lateinit var phrasebookRepo: PhrasebookRepository
    lateinit var database: NerivaDatabase
    lateinit var encryptedPrefs: EncryptedSharedPreferences
    lateinit var cookieJar: PersistentCookieJar

    /**
     * Shared OkHttp client that reuses [cookieJar], so TTS/voice/image requests
     * stay authenticated with the same session cookie as the Retrofit API.
     */
    lateinit var sharedHttpClient: OkHttpClient
        private set

    val appScope = CoroutineScope(Dispatchers.IO)

    override fun onCreate() {
        super.onCreate()
        instance = this

        val masterKey = MasterKey.Builder(this)
            .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
            .build()

        encryptedPrefs = EncryptedSharedPreferences.create(
            this,
            "neriva_secure_prefs",
            masterKey,
            EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
            EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
        ) as EncryptedSharedPreferences

        database = NerivaDatabase.getInstance(this)

        cookieJar = PersistentCookieJar(encryptedPrefs)
        sharedHttpClient = NerivaApiClient.buildClient(cookieJar)
        apiClient = NerivaApiClient.create(
            baseUrl = "https://api.neriva.ru/",
            cookieJar = cookieJar
        )

        authRepo = AuthRepository(apiClient, encryptedPrefs)
        sessionRepo = SessionRepository(apiClient)
        phrasebookRepo = PhrasebookRepository(database.phrasebookDao())
    }

    companion object {
        lateinit var instance: NERIVAApp
            private set

        fun context(): Context = instance.applicationContext
    }
}
