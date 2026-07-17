package ru.neriva.app.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
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
import ru.neriva.app.data.api.PracticeMessage
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

    Scaffold(topBar = {
        TopAppBar(
            title = { Text("Practice") },
            navigationIcon = { IconButton(onClick = { navController.popBackStack() }) { Icon(Icons.AutoMirrored.Filled.ArrowBack, "Back") } }
        )
    }) { padding ->
        Column(modifier = Modifier.padding(padding).fillMaxSize()) {
            LazyColumn(
                modifier = Modifier.weight(1f).fillMaxWidth().padding(horizontal = 16.dp),
                verticalArrangement = Arrangement.spacedBy(8.dp),
                contentPadding = PaddingValues(vertical = 16.dp)
            ) {
                items(messages) { (role, text) ->
                    val isUser = role == "user"
                    Card(
                        modifier = Modifier.fillMaxWidth(),
                        colors = CardDefaults.cardColors(
                            containerColor = if (isUser) MaterialTheme.colorScheme.primaryContainer else MaterialTheme.colorScheme.surfaceVariant
                        )
                    ) {
                        Column(Modifier.padding(12.dp)) {
                            Text(if (isUser) "You" else "AI", style = MaterialTheme.typography.labelSmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                            Text(text, style = MaterialTheme.typography.bodyMedium)
                        }
                    }
                }
                if (loading) item { Box(Modifier.fillMaxWidth(), Alignment.Center) { CircularProgressIndicator(modifier = Modifier.size(24.dp)) } }
                error?.let { item { Text(it, color = MaterialTheme.colorScheme.error) } }
            }
            Row(
                modifier = Modifier.fillMaxWidth().padding(16.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                OutlinedTextField(
                    value = input, onValueChange = { input = it },
                    label = { Text("Message") },
                    modifier = Modifier.weight(1f),
                    maxLines = 4
                )
                Spacer(Modifier.width(8.dp))
                // Voice recording button
                val recorder = remember { AudioRecorder(app) }
                var isRecording by remember { mutableStateOf(false) }
                FilledIconButton(
                    onClick = {
                        if (isRecording) {
                            val file = recorder.stopRecording()
                            isRecording = false
                            if (file != null) {
                                messages = messages + ("user" to "🎤 Voice message")
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
                ) { Icon(Icons.Default.Mic, if (isRecording) "Stop" else "Record") }
                Spacer(Modifier.width(8.dp))
                IconButton(
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
                ) { Icon(Icons.AutoMirrored.Filled.Send, "Send") }
            }
        }
    }
}
