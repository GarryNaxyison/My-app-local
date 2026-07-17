package ru.neriva.app.data.db

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "phrasebook")
data class PhrasebookEntity(
    @PrimaryKey val id: String,
    val phrase: String,
    val translation: String?,
    val note: String?,
    val source: String?,
    val language: String?,
    val createdAt: String?,
)
