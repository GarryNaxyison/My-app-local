package ru.neriva.app.data.api

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class ContextGameRequest(val preset: String? = null, val text: String? = null)

@Serializable
data class ContextGameResponse(val questions: List<ContextGameQuestion>)

@Serializable
data class ContextGameQuestion(
    val context: String,
    val options: List<String>,
    @SerialName("correct_index") val correctIndex: Int,
    val explanation: String,
)
