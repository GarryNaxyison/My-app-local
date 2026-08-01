package ru.neriva.app.ui.screens

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.automirrored.filled.Send
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import ru.neriva.app.NERIVAApp
import ru.neriva.app.data.repo.MistakesRepository

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MistakePracticeScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { MistakesRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var word by remember { mutableStateOf<String?>(null) }
    var correction by remember { mutableStateOf<String?>(null) }
    var context by remember { mutableStateOf<String?>(null) }
    var answer by remember { mutableStateOf("") }
    var result by remember { mutableStateOf<String?>(null) }
    var loading by remember { mutableStateOf(true) }
    var done by remember { mutableStateOf(false) }
    var xp by remember { mutableIntStateOf(0) }

    LaunchedEffect(Unit) {
        try {
            val r = repo.startPractice()
            if (r.empty) { done = true } else {
                word = r.word; correction = r.correction; context = r.context
            }
        } catch (_: Exception) {} finally { loading = false }
    }

    Scaffold(topBar = {
        TopAppBar(
            title = { Text(stringResource(R.string.mistake_repair_plan)) },
            navigationIcon = { IconButton(onClick = { navController.popBackStack() }) { Icon(Icons.AutoMirrored.Filled.ArrowBack, ru.neriva.app.NERIVAApp.instance.getString(R.string.back)) } }
        )
    }) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            if (loading) { Box(Modifier.fillMaxSize(), Alignment.Center) { CircularProgressIndicator() }; return@Column }
            if (done) {
                Text(stringResource(R.string.done), style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.Bold)
                if (xp > 0) Text("+$xp XP", color = MaterialTheme.colorScheme.primary)
                return@Column
            }
            word?.let {
                Text(stringResource(R.string.mistake_original), style = MaterialTheme.typography.labelLarge)
                Text(it, style = MaterialTheme.typography.titleLarge, fontWeight = FontWeight.Bold)
            }
            correction?.let {
                    Text(stringResource(R.string.correct), style = MaterialTheme.typography.labelLarge)
                Text(it, style = MaterialTheme.typography.titleMedium, color = MaterialTheme.colorScheme.primary)
            }
            context?.let { Text(it, style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant) }
            result?.let { Text(it, color = MaterialTheme.colorScheme.secondary) }
            Spacer(Modifier.weight(1f))
            OutlinedTextField(
                value = answer, onValueChange = { answer = it },
                label = { Text(stringResource(R.string.type_answer)) },
                modifier = Modifier.fillMaxWidth()
            )
            Button(
                onClick = {
                    scope.launch {
                        loading = true
                        try {
                            val r = repo.answerPractice(answer)
                            if (r.empty) { done = true } else {
                                word = r.word; correction = r.correction; context = r.context
                                result = if (r.xp != null && r.xp > 0) "+${r.xp} XP" else ru.neriva.app.NERIVAApp.instance.getString(R.string.retry)
                                xp += r.xp ?: 0
                            }
                            answer = ""
                        } catch (_: Exception) {} finally { loading = false }
                    }
                },
                enabled = !loading && answer.isNotBlank() && !done,
                modifier = Modifier.fillMaxWidth()
            ) { Text(stringResource(R.string.send)) }
        }
    }
}
