package ru.neriva.app.ui.screens

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import ru.neriva.app.NERIVAApp
import ru.neriva.app.data.model.WordChallenge
import ru.neriva.app.data.repo.VocabularyRepository

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun WordsScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { VocabularyRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var challenge by remember { mutableStateOf<WordChallenge?>(null) }
    var result by remember { mutableStateOf<String?>(null) }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(Unit) {
        loading = true
        try { challenge = repo.nextWord() } catch (e: Exception) { error = e.message } finally { loading = false }
    }

    Scaffold(topBar = {
        TopAppBar(
            title = { Text(stringResource(R.string.learn_words)) },
            navigationIcon = { IconButton(onClick = { navController.popBackStack() }) { Icon(Icons.AutoMirrored.Filled.ArrowBack, "Back") } }
        )
    }) { padding ->
        Box(Modifier.padding(padding).fillMaxSize().padding(16.dp)) {
            if (loading) {
                CircularProgressIndicator(modifier = Modifier.align(Alignment.Center))
            } else {
                challenge?.let { ch ->
                    if (ch.empty) {
                        Text(ch.message ?: "No words available", modifier = Modifier.align(Alignment.Center))
                    } else {
                        Column(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                            Text(ch.prompt ?: ch.word ?: "", style = MaterialTheme.typography.headlineSmall, fontWeight = FontWeight.Bold)
                            ch.context?.let { Text(it, style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant) }
                            result?.let { Text(it, color = MaterialTheme.colorScheme.primary) }
                            ch.options.forEach { opt ->
                                Button(
                                    onClick = {
                                        scope.launch {
                                            loading = true; error = null; result = null
                                            try {
                                                val resp = repo.answerWord(opt.text, wordId = ch.wordId)
                                                result = if (resp.correct) "Correct! +${resp.xp ?: 0} XP" else "Wrong: ${resp.correctAnswer ?: ""}"
                                                challenge = repo.nextWord()
                                            } catch (e: Exception) { error = e.message } finally { loading = false }
                                        }
                                    },
                                    enabled = !loading,
                                    modifier = Modifier.fillMaxWidth()
                                ) { Text(opt.text) }
                            }
                            error?.let { Text(it, color = MaterialTheme.colorScheme.error) }
                        }
                    }
                }
            }
        }
    }
}
