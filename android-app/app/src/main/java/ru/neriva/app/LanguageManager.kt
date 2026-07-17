package ru.neriva.app

import android.app.Activity
import android.content.Context
import android.content.SharedPreferences
import android.os.Build
import java.util.Locale

object LanguageManager {
    private const val PREFS_NAME = "neriva_language"
    private const val KEY_INTERFACE_LANG = "interface_language"
    private const val KEY_LEARNING_LANG = "learning_language"

    private val supportedLocales = listOf(
        "ru", "en", "es", "de", "fr", "it", "zh", "ja", "ko",
        "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt",
        "ar", "bn", "cs", "el", "hi", "hu", "id", "nl", "sv",
        "ta", "te", "th", "tl", "tr", "vi"
    )

    private val languageNames = mapOf(
        "ru" to "Русский",
        "en" to "English",
        "es" to "Espanol",
        "de" to "Deutsch",
        "fr" to "Francais",
        "it" to "Italiano",
        "zh" to "中文",
        "ja" to "日本語",
        "ko" to "한국어",
        "tg" to "Тоҷикӣ",
        "uz" to "O'zbekcha",
        "tt" to "Татарча",
        "hy" to "Հայերեն",
        "kk" to "Қазақша",
        "ky" to "Кыргызча",
        "ka" to "ქართული",
        "uk" to "Українська",
        "pl" to "Polski",
        "ro" to "Romana",
        "pt" to "Portugues",
        "ar" to "العربية",
        "bn" to "বাংলা",
        "cs" to "Cestina",
        "el" to "Ελληνικά",
        "hi" to "हिंदी",
        "hu" to "Magyar",
        "id" to "Bahasa Indonesia",
        "nl" to "Nederlands",
        "sv" to "svenska",
        "ta" to "தமிழ்",
        "te" to "తెలుగు",
        "th" to "ภาษาไทย",
        "tl" to "Tagalog",
        "tr" to "Turkce",
        "vi" to "Tieng Viet"
    )

    private fun prefs(context: Context): SharedPreferences =
        context.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)

    fun getInterfaceLanguage(context: Context): String =
        prefs(context).getString(KEY_INTERFACE_LANG, "ru") ?: "ru"

    fun getLearningLanguage(context: Context): String =
        prefs(context).getString(KEY_LEARNING_LANG, "en") ?: "en"

    fun setInterfaceLanguage(context: Context, lang: String) {
        prefs(context).edit().putString(KEY_INTERFACE_LANG, lang).apply()
    }

    fun setLearningLanguage(context: Context, lang: String) {
        prefs(context).edit().putString(KEY_LEARNING_LANG, lang).apply()
    }

    fun getLanguageName(code: String): String =
        languageNames[code] ?: code.uppercase()

    fun getSupportedLocales(): List<String> = supportedLocales

    fun applyLocale(activity: Activity, langCode: String) {
        val locale = Locale(langCode)
        Locale.setDefault(locale)
        val config = activity.resources.configuration
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.JELLY_BEAN_MR1) {
            config.setLocale(locale)
        } else {
            @Suppress("DEPRECATION")
            config.locale = locale
        }
        activity.resources.updateConfiguration(config, activity.resources.displayMetrics)
    }
}
