package ru.neriva.app

import org.junit.Assert.assertEquals
import org.junit.Test
import ru.neriva.app.ui.nav.Screen

class LanguageManagerTest {
    @Test
    fun exposesEverySupportedInterfaceLocale() {
        assertEquals(
            listOf(
                "ru", "en", "es", "de", "fr", "it", "zh", "ja", "ko",
                "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt",
                "ar", "bn", "cs", "el", "hi", "hu", "id", "nl", "sv",
                "ta", "te", "th", "tl", "tr", "vi",
            ),
            LanguageManager.getSupportedLocales(),
        )
    }

    @Test
    fun navigationUsesLocalizedResourceIds() {
        assertEquals(R.string.today, Screen.Home.titleRes)
        assertEquals(R.string.settings, Screen.Settings.titleRes)
        assertEquals(R.string.learn_words, Screen.Words.titleRes)
    }
}
