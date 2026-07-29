package ru.neriva.app.data.api

import org.junit.Assert.assertEquals
import org.junit.Test

class NerivaApiClientTest {
    @Test
    fun normalizesBaseUrlWithTrailingSlash() {
        assertEquals("https://api.neriva.ru/", NerivaApiClient.normalizeBaseUrl("https://api.neriva.ru"))
    }
}
