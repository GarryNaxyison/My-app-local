package ru.neriva.app.ui.screens

import ru.neriva.app.data.model.GameCard
import ru.neriva.app.data.model.gameCards
import ru.neriva.app.data.model.spellingTiles
import androidx.annotation.StringRes
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.Check
import androidx.compose.material.icons.outlined.Extension
import androidx.compose.material.icons.outlined.Forum
import androidx.compose.material.icons.outlined.Spellcheck
import androidx.compose.material.icons.outlined.Refresh
import androidx.compose.ui.Alignment
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.LiveRegionMode
import androidx.compose.ui.semantics.liveRegion
import androidx.compose.ui.semantics.semantics
import ru.neriva.app.ui.components.NerivaAction
import ru.neriva.app.ui.components.nerivaWingCover
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.navigation.NavHostController
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.launch
import ru.neriva.app.NERIVAApp
import ru.neriva.app.R
import ru.neriva.app.data.api.ContextGameQuestion
import ru.neriva.app.data.api.ContextGameRequest
import ru.neriva.app.ui.nav.Screen

private enum class Mode(@StringRes val title: Int) {
    Context(R.string.game_context), Pairs(R.string.game_pairs), Spelling(R.string.game_build_word),
}

@OptIn(ExperimentalLayoutApi::class)
@Composable
fun LearningGamesScreen(navController: NavHostController) {
    var mode by rememberSaveable { mutableStateOf(Mode.Context) }
    ScreenScaffold(stringResource(R.string.learning_games), navController) { padding ->
        Box(Modifier.padding(padding).fillMaxSize(), contentAlignment = Alignment.TopCenter) {
            Column(
                Modifier.widthIn(max = 680.dp).fillMaxWidth().verticalScroll(rememberScrollState()).padding(18.dp),
                verticalArrangement = Arrangement.spacedBy(12.dp),
            ) {
                FlowRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    Mode.entries.forEach { item ->
                        FilterChip(mode == item, { mode = item }, { Text(stringResource(item.title)) })
                    }
                }
                GameCover(mode)
                key(mode) {
                    if (mode == Mode.Context) {
                        ContextGame()
                    } else {
                        VocabularyGame(mode == Mode.Spelling) { navController.navigate(Screen.Vocabulary.route) }
                    }
                }
            }
        }
    }
}

@Composable
private fun GameCover(mode: Mode) {
    val accent by animateColorAsState(
        when (mode) {
            Mode.Context -> Color(0xFF4F46BF)
            Mode.Pairs -> Color(0xFF1A636C)
            Mode.Spelling -> Color(0xFF8B416B)
        },
        label = "gameModeColor",
    )
    val icon = when (mode) {
        Mode.Context -> Icons.Outlined.Forum
        Mode.Pairs -> Icons.Outlined.Extension
        Mode.Spelling -> Icons.Outlined.Spellcheck
    }
    Surface(shape = MaterialTheme.shapes.extraLarge, contentColor = Color.White) {
        Row(
            Modifier.fillMaxWidth().nerivaWingCover(accent).padding(20.dp),
            horizontalArrangement = Arrangement.spacedBy(14.dp),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Icon(icon, null, Modifier.size(38.dp))
            Text(stringResource(mode.title), style = MaterialTheme.typography.displaySmall)
        }
    }
}

@OptIn(ExperimentalLayoutApi::class)
@Composable
private fun ContextGame() {
    val scope = rememberCoroutineScope()
    var source by rememberSaveable { mutableStateOf("") }
    var questions by remember { mutableStateOf(emptyList<ContextGameQuestion>()) }
    var index by remember { mutableIntStateOf(0) }
    var score by remember { mutableIntStateOf(0) }
    var selected by remember { mutableIntStateOf(-1) }
    var busy by remember { mutableStateOf(false) }
    var failed by remember { mutableStateOf(false) }

    fun load(request: ContextGameRequest) {
        if (busy) return
        busy = true
        failed = false
        scope.launch {
            try {
                val result = NERIVAApp.instance.apiClient.contextGame(request).questions
                require(result.isNotEmpty() && result.all { it.options.size == 3 && it.correctIndex in it.options.indices })
                questions = result
                index = 0
                score = 0
                selected = -1
            } catch (e: CancellationException) {
                throw e
            } catch (_: Exception) {
                failed = true
            } finally {
                busy = false
            }
        }
    }

    Text(stringResource(R.string.game_context_prompt), style = MaterialTheme.typography.titleMedium)
    FlowRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
        listOf("travel" to R.string.game_travel, "restaurant" to R.string.game_restaurant, "work" to R.string.game_work)
            .forEach { (preset, label) ->
                OutlinedButton(onClick = { load(ContextGameRequest(preset = preset)) }, enabled = !busy) {
                    Text(stringResource(label))
                }
            }
    }
    OutlinedTextField(
        value = source,
        onValueChange = { if (it.codePointCount(0, it.length) <= 1200) source = it },
        label = { Text(stringResource(R.string.game_own_text)) },
        modifier = Modifier.fillMaxWidth(),
        minLines = 2,
        maxLines = 5,
        enabled = !busy,
    )
    NerivaAction(
        stringResource(R.string.game_generate),
        { load(ContextGameRequest(text = source.trim())) },
        enabled = !busy && source.isNotBlank(),
        modifier = Modifier.fillMaxWidth(),
    )
    if (busy) LinearProgressIndicator(Modifier.fillMaxWidth())
    if (failed) Text(stringResource(R.string.game_load_error), color = MaterialTheme.colorScheme.error)
    if (busy) return
    val question = questions.getOrNull(index)
    if (question != null) {
        Text(stringResource(R.string.game_round_progress, index + 1, questions.size))
        val progress by animateFloatAsState((index + 1).toFloat() / questions.size, label = "gameProgress")
        LinearProgressIndicator(progress = { progress }, modifier = Modifier.fillMaxWidth())
        Surface(shape = MaterialTheme.shapes.large, color = MaterialTheme.colorScheme.surface) {
            Text(question.context, Modifier.fillMaxWidth().padding(18.dp), style = MaterialTheme.typography.headlineSmall)
        }
        question.options.forEachIndexed { option, word ->
            FilterChip(selected == option, {
                selected = option
                if (option == question.correctIndex) score++
            }, { Text(word) }, enabled = selected == -1, modifier = Modifier.fillMaxWidth())
        }
        AnimatedVisibility(selected != -1) {
            val correct = selected == question.correctIndex
            val colors = MaterialTheme.colorScheme
            Surface(
                shape = MaterialTheme.shapes.large,
                color = if (correct) colors.primaryContainer else colors.errorContainer,
                contentColor = if (correct) colors.onPrimaryContainer else colors.onErrorContainer,
                modifier = Modifier.fillMaxWidth().semantics { liveRegion = LiveRegionMode.Polite },
            ) {
                Column(Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                        Icon(if (correct) Icons.Outlined.Check else Icons.Outlined.Refresh, null)
                        Text(stringResource(if (correct) R.string.game_correct else R.string.game_try_again), style = MaterialTheme.typography.titleLarge)
                    }
                    Text(stringResource(R.string.game_correct_word, question.options[question.correctIndex]))
                    Text(question.explanation)
                    NerivaAction(stringResource(R.string.game_continue), { index++; selected = -1 })
                }
            }
        }
    } else if (questions.isNotEmpty()) {
        Text(stringResource(R.string.game_score, score, questions.size), style = MaterialTheme.typography.titleLarge)
    }
}

@OptIn(ExperimentalLayoutApi::class)
@Composable
private fun VocabularyGame(spelling: Boolean, openVocabulary: () -> Unit) {
    val scope = rememberCoroutineScope()
    var cards by remember { mutableStateOf<List<GameCard>?>(null) }
    var translations by remember { mutableStateOf(emptyList<Int>()) }
    var word by remember { mutableIntStateOf(-1) }
    var matched by remember { mutableStateOf(emptySet<Int>()) }
    var puzzle by remember { mutableIntStateOf(0) }
    var tiles by remember { mutableStateOf(emptyList<String>()) }
    var picked by remember { mutableStateOf(emptyList<Int>()) }
    var failed by remember { mutableStateOf(false) }
    var incorrect by remember { mutableStateOf(false) }
    var busy by remember { mutableStateOf(false) }
    val round = cards
    val enough = round != null && round.size >= if (spelling) 1 else 2
    val complete = enough && round != null && if (spelling) puzzle == round.size else matched.size == round.size

    Text(stringResource(if (spelling) R.string.game_spelling_prompt else R.string.game_pairs_prompt))
    if (!enough || complete || failed) {
        Button(onClick = {
            busy = true
            failed = false
            scope.launch {
                try {
                    val items = NERIVAApp.instance.apiClient.vocabulary(page = 0).items.orEmpty()
                    val next = gameCards(items.shuffled(), spelling)
                    cards = next
                    translations = next.indices.shuffled()
                    word = -1
                    matched = emptySet()
                    puzzle = 0
                    picked = emptyList()
                    tiles = if (spelling) next.firstOrNull()?.let { spellingTiles(it.word).shuffled() }.orEmpty() else emptyList()
                    incorrect = false
                } catch (e: CancellationException) {
                    throw e
                } catch (_: Exception) {
                    failed = true
                } finally {
                    busy = false
                }
            }
        }, enabled = !busy) { Text(stringResource(R.string.game_new_round)) }
    }
    if (busy) {
        LinearProgressIndicator(Modifier.fillMaxWidth())
        return
    }
    if (failed) Text(stringResource(R.string.game_load_error), color = MaterialTheme.colorScheme.error)
    if (!enough || round == null) {
        if (round != null) Text(stringResource(R.string.game_not_enough_words))
        OutlinedButton(onClick = openVocabulary) { Text(stringResource(R.string.vocabulary)) }
        return
    }
    if (complete) {
        Text(stringResource(R.string.game_score, round.size, round.size), style = MaterialTheme.typography.titleLarge)
        return
    }

    if (spelling) {
        val card = round[puzzle]
        Text(stringResource(R.string.game_round_progress, puzzle + 1, round.size))
        Text(card.translation, style = MaterialTheme.typography.titleLarge)
        Text(picked.joinToString("") { tiles[it] }.ifEmpty { "…" })
        FlowRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            tiles.forEachIndexed { index, tile ->
                OutlinedButton(onClick = {
                    picked = picked + index
                    incorrect = false
                }, enabled = index !in picked) { Text(tile) }
            }
        }
        Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            OutlinedButton(onClick = { picked = picked.dropLast(1); incorrect = false }, enabled = picked.isNotEmpty()) {
                Text(stringResource(R.string.game_undo))
            }
            Button(onClick = {
                if (picked.map(tiles::get) == spellingTiles(card.word)) {
                    puzzle++
                    tiles = round.getOrNull(puzzle)?.let { spellingTiles(it.word).shuffled() }.orEmpty()
                    picked = emptyList()
                } else {
                    incorrect = true
                    picked = emptyList()
                }
            }, enabled = picked.size == tiles.size) { Text(stringResource(R.string.game_check)) }
        }
    } else {
        Text(stringResource(R.string.game_score, matched.size, round.size))
        Text(stringResource(R.string.game_words), style = MaterialTheme.typography.titleMedium)
        FlowRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            round.forEachIndexed { index, card ->
                FilterChip(word == index || index in matched, { word = index; incorrect = false }, {
                    Text(card.word)
                }, enabled = index !in matched)
            }
        }
        Text(stringResource(R.string.game_translations), style = MaterialTheme.typography.titleMedium)
        FlowRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            translations.forEach { index ->
                OutlinedButton(onClick = {
                    incorrect = word != index
                    if (!incorrect) {
                        matched = matched + index
                        word = -1
                    }
                }, enabled = word >= 0 && index !in matched) { Text(round[index].translation) }
            }
        }
    }
    if (incorrect) Text(stringResource(R.string.game_try_again), color = MaterialTheme.colorScheme.error)
}
