package ru.neriva.app

import org.junit.Assert.assertEquals
import org.junit.Test

class LanguageManagerTest {
    @Test
    fun exposesOnlyReleaseOneInterfaceLocales() {
        assertEquals(listOf("ru", "en"), LanguageManager.getSupportedLocales())
    }
}
