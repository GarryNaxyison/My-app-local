package ru.neriva.app.data.api

import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import org.junit.Assert.assertEquals
import org.junit.Test

class MistakeDeleteRequestTest {
    @Test
    fun serializesTheZeroBasedServerIndex() {
        assertEquals("{\"index\":0}", Json.encodeToString(MistakeDeleteRequest(index = 0)))
    }
}
