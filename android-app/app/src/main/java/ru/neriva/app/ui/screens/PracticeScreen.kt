package ru.neriva.app.ui.screens

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.itemsIndexed
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.ui.graphics.Color
import androidx.compose.foundation.shape.RoundedCornerShape
import ru.neriva.app.ui.components.nerivaWingCover
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.Send
import androidx.compose.material.icons.filled.Mic
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import ru.neriva.app.NERIVAApp
import ru.neriva.app.audio.AudioPlayer
import ru.neriva.app.audio.AudioRecorder
import ru.neriva.app.data.repo.PracticeRepository

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PracticeScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { PracticeRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var input by remember { mutableStateOf("") }
    var messages by remember { mutableStateOf(listOf<Pair<String, String>>()) }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }

    val listState = rememberLazyListState()
    LaunchedEffect(messages.size) {
        if (messages.isNotEmpty()) listState.animateScrollToItem(messages.size)
    }
    ScreenScaffold(stringResource(R.string.practice), navController) { padding ->
        Column(modifier = Modifier.padding(padding).fillMaxSize()) {
            LazyColumn(
                state = listState,
                modifier = Modifier.weight(1f).fillMaxWidth().padding(horizontal = 16.dp),
                verticalArrangement = Arrangement.spacedBy(8.dp),
                contentPadding = PaddingValues(vertical = 16.dp)
            ) {
                item(key = "practice-intro", contentType = "intro") {
                    Surface(shape = MaterialTheme.shapes.extraLarge, contentColor = Color.White) {
                        Text(
                            stringResource(R.string.practice_prompt),
                            Modifier.fillMaxWidth().nerivaWingCover().padding(20.dp),
                            style = MaterialTheme.typography.headlineSmall,
                        )
                    }
                }
                itemsIndexed(messages, key = { index, _ -> index }, contentType = { _, _ -> "message" }) { _, (role, text) ->
                    val isUser = role == "user"
                    val colors = MaterialTheme.colorScheme
                    Box(Modifier.fillMaxWidth().animateItem(), contentAlignment = if (isUser) Alignment.CenterEnd else Alignment.CenterStart) {
                        Surface(
                            modifier = Modifier.fillMaxWidth(.87f),
                            shape = if (isUser) RoundedCornerShape(20.dp, 20.dp, 5.dp, 20.dp)
                                else RoundedCornerShape(20.dp, 20.dp, 20.dp, 5.dp),
                            color = if (isUser) colors.primaryContainer else colors.surface,
                            contentColor = if (isUser) colors.onPrimaryContainer else colors.onSurface,
                        ) {
                            Column(Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(4.dp)) {
                                Text(if (isUser) stringResource(R.string.you) else stringResource(R.string.ai_tutor), style = MaterialTheme.typography.labelMedium)
                                Text(text, style = MaterialTheme.typography.bodyMedium)
                            }
                        }
                    }
                }
                if (loading) item { Box(Modifier.fillMaxWidth(), Alignment.Center) { CircularProgressIndicator(modifier = Modifier.size(24.dp)) } }
                error?.let { item { Text(it, color = MaterialTheme.colorScheme.error) } }
            }
            Row(
                modifier = Modifier.fillMaxWidth().imePadding().padding(16.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                OutlinedTextField(
                    value = input, onValueChange = { input = it },
                    label = { Text(stringResource(R.string.input)) },
                    modifier = Modifier.weight(1f),
                    maxLines = 4
                )
                Spacer(Modifier.width(8.dp))
                // Voice recording button
                val recorder = remember { AudioRecorder(app) }
                var isRecording by remember { mutableStateOf(false) }
                FilledIconButton(
                    enabled = !loading,
                    onClick = {
                        if (isRecording) {
                            val file = recorder.stopRecording()
                            isRecording = false
                            if (file != null) {
                                messages = messages + ("user" to ru.neriva.app.NERIVAApp.instance.getString(R.string.voice_message))
                                scope.launch {
                                    loading = true; error = null
                                    try {
                                        val resp = repo.sendVoice(file)
                                        val reply = resp.message ?: "…"
                                        messages = messages + ("ai" to reply)
                                        resp.audioUrl?.let { url ->
                                            AudioPlayer.playUrl(url)
                                        }
                                    } catch (e: Exception) { error = e.message }
                                    finally { loading = false }
                                }
                            }
                        } else {
                            recorder.startRecording()
                            isRecording = true
                        }
                    },
                    colors = IconButtonDefaults.filledIconButtonColors(
                        containerColor = if (isRecording) MaterialTheme.colorScheme.error else MaterialTheme.colorScheme.primary
                    )
                ) { Icon(Icons.Default.Mic, if (isRecording) ru.neriva.app.NERIVAApp.instance.getString(R.string.stop_recording) else ru.neriva.app.NERIVAApp.instance.getString(R.string.record_voice)) }
                Spacer(Modifier.width(8.dp))
                FilledIconButton(
                    enabled = input.isNotBlank() && !loading && !isRecording,
                    onClick = {
                        if (input.isNotBlank() && !loading) {
                            val msg = input; input = ""
                            messages = messages + ("user" to msg)
                            scope.launch {
                                loading = true; error = null
                                try {
                                    val resp = repo.send(message = msg)
                                    messages = messages + ("ai" to (resp.message ?: "..."))
                                } catch (e: Exception) { error = e.message }
                                finally { loading = false }
                            }
                        }
                    }
                ) { Icon(Icons.AutoMirrored.Filled.Send, ru.neriva.app.NERIVAApp.instance.getString(R.string.send)) }
            }
        }
    }
}
