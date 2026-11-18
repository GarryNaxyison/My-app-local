package ru.neriva.app.ui.screens

import ru.neriva.app.data.model.gameCards
import ru.neriva.app.data.model.spellingTiles
import org.junit.Assert.assertEquals
import org.junit.Test
import ru.neriva.app.data.model.VocabularyItem

class LearningGamesTest {
    @Test
    fun pairsKeepDistinctWordsAndTranslations() {
        val items = listOf(
            VocabularyItem(word = "Cat", translation = "cat translation"),
            VocabularyItem(word = " cat ", translation = "duplicate"),
            VocabularyItem(word = "Kitten", translation = "cat translation"),
            VocabularyItem(word = "Dog", translation = "dog translation"),
            VocabularyItem(word = "Missing"),
        )
        val cards = gameCards(items, spelling = false)
        assertEquals(listOf("Cat", "Dog"), cards.map { it.word })
        assertEquals(listOf("cat translation", "dog translation"), cards.map { it.translation })
    }

    @Test
    fun spellingNormalizesCombiningMarksAndSkipsPhrases() {
        val items = listOf(
            VocabularyItem(word = "cafe\u0301", translation = "coffee"),
            VocabularyItem(word = "two words", translation = "phrase"),
            VocabularyItem(word = "a", translation = "letter"),
        )
        val card = gameCards(items, spelling = true).single()
        assertEquals(listOf("c", "a", "f", "\u00e9"), spellingTiles(card.word))
    }
}
