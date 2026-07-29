package ru.neriva.app.data.api

import android.content.SharedPreferences
import kotlinx.serialization.json.Json
import okhttp3.Cookie
import okhttp3.CookieJar
import okhttp3.HttpUrl
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import com.jakewharton.retrofit2.converter.kotlinx.serialization.asConverterFactory
import java.util.concurrent.TimeUnit

class PersistentCookieJar(private val prefs: SharedPreferences) : CookieJar {

    private val cookies = mutableMapOf<String, MutableList<Cookie>>()

    init {
        load()
    }

    override fun saveFromResponse(url: HttpUrl, cookies: List<Cookie>) {
        val host = url.host
        this.cookies[host] = cookies.filter { it.persistent }.toMutableList()
        save()
    }

    override fun loadForRequest(url: HttpUrl): List<Cookie> {
        return cookies[url.host]?.filterNot { it.expiresAt < System.currentTimeMillis() } ?: emptyList()
    }

    private fun save() {
        prefs.edit().apply {
            val all = cookies.flatMap { (host, list) ->
                list.map { c -> host to c }
            }
            val jsonArr = all.map { (host, c) ->
                "${c.name}=${c.value}|${c.path}|${c.hostOnly}|${c.secure}|${c.httpOnly}|${c.persistent}|${c.expiresAt}|$host"
            }
            putString("cookies_json", jsonArr.joinToString(separator = "§"))
            apply()
        }
    }

    private fun load() {
        val raw = prefs.getString("cookies_json", null) ?: return
        try {
            val list = raw.split("§")
            for (item in list) {
                if (item.isBlank()) continue
                val parts = item.split("|")
                if (parts.size < 8) continue
                val nv = parts[0].split("=", limit = 2)
                if (nv.size < 2) continue
                val domain = parts[7]
                val builder = Cookie.Builder()
                    .name(nv[0]).value(nv[1])
                    .domain(domain)
                    .path(parts[1])
                if (parts[3].toBoolean()) builder.secure()
                if (parts[4].toBoolean()) builder.httpOnly()
                builder.expiresAt(parts[6].toLong())
                val cookie = builder.build()
                cookies.getOrPut(domain) { mutableListOf() }.add(cookie)
            }
        } catch (_: Exception) {}
    }
}

object NerivaApiClient {

    fun normalizeBaseUrl(baseUrl: String): String = baseUrl.trimEnd('/') + "/"

    /** Build the shared OkHttp client from a pre-built [cookieJar]. */
    fun buildClient(cookieJar: PersistentCookieJar): OkHttpClient {
        val logging = HttpLoggingInterceptor().apply {
            level = HttpLoggingInterceptor.Level.BASIC
        }
        return OkHttpClient.Builder()
            .cookieJar(cookieJar)
            .addInterceptor(logging)
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(60, TimeUnit.SECONDS)
            .writeTimeout(60, TimeUnit.SECONDS)
            .build()
    }

    /** Build a Retrofit API backed by [cookieJar]. */
    fun create(baseUrl: String, cookieJar: PersistentCookieJar): NerivaApi {
        val json = Json {
            ignoreUnknownKeys = true
            isLenient = true
            encodeDefaults = true
        }

        val client = buildClient(cookieJar)

        return Retrofit.Builder()
            .baseUrl(normalizeBaseUrl(baseUrl))
            .client(client)
            .addConverterFactory(json.asConverterFactory("application/json".toMediaType()))
            .build()
            .create(NerivaApi::class.java)
    }
}
