package ru.neriva.app

import android.content.Context
import android.content.SharedPreferences

object ThemeManager {
    private const val PREFS_NAME = "neriva_theme"
    private const val KEY_THEME = "theme_mode"

    enum class ThemeMode { LIGHT, DARK, SYSTEM }

    private fun prefs(context: Context): SharedPreferences =
        context.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE)

    fun getThemeMode(context: Context): ThemeMode {
        val value = prefs(context).getString(KEY_THEME, "system") ?: "system"
        return when (value) {
            "light" -> ThemeMode.LIGHT
            "dark" -> ThemeMode.DARK
            else -> ThemeMode.SYSTEM
        }
    }

    fun setThemeMode(context: Context, mode: ThemeMode) {
        prefs(context).edit().putString(KEY_THEME, mode.name.lowercase()).apply()
    }

    fun isDarkTheme(context: Context): Boolean {
        return when (getThemeMode(context)) {
            ThemeMode.DARK -> true
            ThemeMode.LIGHT -> false
            ThemeMode.SYSTEM -> false // Will be resolved by isSystemInDarkTheme()
        }
    }
}
