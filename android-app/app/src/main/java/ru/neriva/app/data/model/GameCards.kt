package ru.neriva.app.data.model

import java.text.Normalizer

internal data class GameCard(val word: String, val translation: String)

private val spellingWord = Regex("\\p{L}[\\p{L}\\p{M}]+")
private val spellingTile = Regex("\\P{M}\\p{M}*|\\p{M}+")

internal fun spellingTiles(word: String): List<String> =
    spellingTile.findAll(Normalizer.normalize(word, Normalizer.Form.NFC)).map { it.value }.toList()

internal fun gameCards(items: List<VocabularyItem>, spelling: Boolean): List<GameCard> {
    val words = mutableSetOf<String>()
    val translations = mutableSetOf<String>()
    return items.asSequence().mapNotNull { item ->
        val word = Normalizer.normalize(item.word.trim(), Normalizer.Form.NFC)
        val translation = item.translation?.trim().orEmpty()
        if (word.isEmpty() || translation.isEmpty()) return@mapNotNull null
        if (spelling && (!spellingWord.matches(word) || spellingTiles(word).size !in 2..20)) {
            return@mapNotNull null
        }
        val wordKey = word.lowercase()
        val translationKey = translation.lowercase()
        if (wordKey in words || translationKey in translations) return@mapNotNull null
        words.add(wordKey)
        translations.add(translationKey)
        GameCard(word, translation)
    }.take(5).toList()
}
