package ru.neriva.app.ui.screens

import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.togetherWith
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.semantics.LiveRegionMode
import androidx.compose.ui.semantics.liveRegion
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.dp
import androidx.navigation.NavHostController
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.launch
import ru.neriva.app.NERIVAApp
import ru.neriva.app.R
import ru.neriva.app.data.api.WordAnswerResponse
import ru.neriva.app.data.model.WordChallenge
import ru.neriva.app.data.repo.VocabularyRepository
import ru.neriva.app.ui.components.NerivaAction
import ru.neriva.app.ui.components.nerivaWingCover

@Composable
fun WordsScreen(navController: NavHostController) {
    val repo = remember { VocabularyRepository(NERIVAApp.instance.apiClient) }
    val scope = rememberCoroutineScope()
    var challenge by remember { mutableStateOf<WordChallenge?>(null) }
    var result by remember { mutableStateOf<WordAnswerResponse?>(null) }
    var loading by remember { mutableStateOf(true) }
    var failed by remember { mutableStateOf(false) }

    suspend fun load() {
        loading = true
        failed = false
        try {
            challenge = repo.nextWord()
            result = null
        } catch (cancelled: CancellationException) {
            throw cancelled
        } catch (_: Exception) {
            failed = true
        } finally {
            loading = false
        }
    }
    LaunchedEffect(Unit) { load() }

    val colors = MaterialTheme.colorScheme
    ScreenScaffold(stringResource(R.string.learn_words), navController) { padding ->
        Column(
            Modifier.padding(padding).fillMaxSize().verticalScroll(rememberScrollState()).padding(18.dp),
            verticalArrangement = Arrangement.spacedBy(14.dp),
        ) {
            if (loading) LinearProgressIndicator(Modifier.fillMaxWidth())
            if (failed) {
                Text(stringResource(R.string.request_failed), color = colors.error)
                if (challenge == null) NerivaAction(stringResource(R.string.retry), { scope.launch { load() } })
            }
            AnimatedContent(
                challenge,
                transitionSpec = { fadeIn() togetherWith fadeOut() },
                label = "wordCard",
            ) { current ->
                if (current != null && current.empty) {
                    Text(current.message ?: stringResource(R.string.empty_panel))
                } else if (current != null) {
                    Column(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                        Surface(shape = MaterialTheme.shapes.extraLarge, contentColor = Color.White) {
                            Column(
                                Modifier.fillMaxWidth().nerivaWingCover().padding(22.dp),
                                verticalArrangement = Arrangement.spacedBy(8.dp),
                            ) {
                                Text(current.prompt ?: current.word.orEmpty(), style = MaterialTheme.typography.displaySmall)
                                current.context?.takeIf { it.isNotBlank() }?.let {
                                    Text(it, style = MaterialTheme.typography.bodyLarge, color = Color(0xFFE8EBFF))
                                }
                            }
                        }
                        current.options.forEach { option ->
                            NerivaAction(
                                option.text,
                                onClick = {
                                    if (!loading && result == null) {
                                        loading = true
                                        failed = false
                                        scope.launch {
                                            try {
                                                result = repo.answerWord(option.text, wordId = current.wordId)
                                            } catch (cancelled: CancellationException) {
                                                throw cancelled
                                            } catch (_: Exception) {
                                                failed = true
                                            } finally {
                                                loading = false
                                            }
                                        }
                                    }
                                },
                                enabled = !loading && result == null,
                                modifier = Modifier.fillMaxWidth(),
                            )
                        }
                    }
                }
            }
            AnimatedContent(result, label = "wordFeedback") { feedback ->
                if (feedback != null) {
                    Surface(
                        modifier = Modifier.fillMaxWidth().semantics { liveRegion = LiveRegionMode.Polite },
                        shape = MaterialTheme.shapes.large,
                        color = if (feedback.correct) colors.primaryContainer else colors.errorContainer,
                        contentColor = if (feedback.correct) colors.onPrimaryContainer else colors.onErrorContainer,
                    ) {
                        Column(Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                            Text(
                                stringResource(if (feedback.correct) R.string.correct else R.string.game_try_again),
                                style = MaterialTheme.typography.titleLarge,
                            )
                            feedback.correctAnswer?.takeIf { !feedback.correct && it.isNotBlank() }?.let {
                                Text(stringResource(R.string.game_correct_word, it))
                            }
                            feedback.xp?.takeIf { it > 0 }?.let { Text("+$it XP") }
                            NerivaAction(
                                stringResource(R.string.tutor_continue),
                                { scope.launch { load() } },
                                enabled = !loading,
                                modifier = Modifier.fillMaxWidth(),
                            )
                        }
                    }
                }
            }
        }
    }
}
