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

    private val supportedLearningLanguages = listOf(
        "ru", "en", "es", "de", "fr", "it", "zh", "ja", "ko",
        "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt",
        "ar", "bn", "cs", "el", "hi", "hu", "id", "nl", "sv",
        "ta", "te", "th", "tl", "tr", "vi",
    )
    private val languageNames = mapOf(
        "ru" to "Русский",
        "en" to "English",
        "es" to "Spanish",
        "de" to "German",
        "fr" to "French",
        "it" to "Italian",
        "zh" to "Chinese",
        "ja" to "Japanese",
        "ko" to "Korean",
        "tg" to "Tajik",
        "uz" to "Uzbek",
        "tt" to "Tatar",
        "hy" to "Armenian",
        "kk" to "Kazakh",
        "ky" to "Kyrgyz",
        "ka" to "Georgian",
        "uk" to "Ukrainian",
        "pl" to "Polish",
        "ro" to "Romanian",
        "pt" to "Portuguese",
        "ar" to "Arabic",
        "bn" to "Bengali",
        "cs" to "Czech",
        "el" to "Greek",
        "hi" to "Hindi",
        "hu" to "Hungarian",
        "id" to "Indonesian",
        "nl" to "Dutch",
        "sv" to "Swedish",
        "ta" to "Tamil",
        "te" to "Telugu",
        "th" to "Thai",
        "tl" to "Tagalog",
        "tr" to "Turkish",
        "vi" to "Vietnamese",
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

    fun getLanguageName(code: String): String = languageNames[code] ?: code.uppercase()

    fun getSupportedLocales(): List<String> = supportedLearningLanguages

    fun getSupportedLearningLanguages(): List<String> = supportedLearningLanguages

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
