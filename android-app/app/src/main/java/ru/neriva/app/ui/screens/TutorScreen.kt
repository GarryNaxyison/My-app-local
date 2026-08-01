package ru.neriva.app.ui.screens

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
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
import ru.neriva.app.data.model.AiTutorResponse
import ru.neriva.app.data.repo.TutorRepository

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TutorScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { TutorRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var loading by remember { mutableStateOf(false) }
    var response by remember { mutableStateOf<AiTutorResponse?>(null) }
    var error by remember { mutableStateOf<String?>(null) }
    var answer by remember { mutableStateOf("") }

    LaunchedEffect(Unit) {
        loading = true
        try { response = repo.start() } catch (e: Exception) { error = e.message } finally { loading = false }
    }

    Scaffold(topBar = {
        TopAppBar(
            title = { Text(stringResource(R.string.ai_tutor)) },
            navigationIcon = { IconButton(onClick = { navController.popBackStack() }) { Icon(Icons.AutoMirrored.Filled.ArrowBack, ru.neriva.app.NERIVAApp.instance.getString(R.string.back)) } }
        )
    }) { padding ->
        Column(
            modifier = Modifier.padding(padding).fillMaxSize().padding(16.dp).verticalScroll(rememberScrollState()),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            if (loading) {
                Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) { CircularProgressIndicator() }
                return@Column
            }
            error?.let {
                Text(it, color = MaterialTheme.colorScheme.error)
            }
            response?.let { resp ->
                resp.lesson?.let { lesson ->
                    lesson.title?.let { Text(it, style = MaterialTheme.typography.headlineSmall, fontWeight = FontWeight.Bold) }
                    lesson.theme?.let { Text(it, style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant) }
                    lesson.story?.textTarget?.let {
                        Card(modifier = Modifier.fillMaxWidth()) {
                            Text(it, modifier = Modifier.padding(16.dp), style = MaterialTheme.typography.bodyLarge)
                        }
                    }
                    if (lesson.words.isNotEmpty()) {
                        Text(stringResource(R.string.tutor_words), style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                        lesson.words.forEach { w ->
                            Card(modifier = Modifier.fillMaxWidth()) {
                                Column(Modifier.padding(12.dp)) {
                                    Text(w.target ?: "", fontWeight = FontWeight.SemiBold)
                                    Text(w.interfaceTranslation ?: "", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                                    w.exampleSentenceTarget?.let { Text("Example: $it", style = MaterialTheme.typography.bodySmall) }
                                }
                            }
                        }
                    }
                }
                resp.nextStep?.let { step ->
                    step.instruction?.let { Text(it, style = MaterialTheme.typography.bodyLarge) }
                    step.question?.questionTarget?.let { Text(it, style = MaterialTheme.typography.titleMedium) }
                    step.options?.let { opts ->
                        opts.forEach { opt ->
                            Button(
                                onClick = {
                                    scope.launch {
                                        loading = true
                                        try { response = repo.answer(opt.text ?: "", response?.session?.id) }
                                        catch (e: Exception) { error = e.message } finally { loading = false }
                                    }
                                },
                                modifier = Modifier.fillMaxWidth()
                            ) { Text(opt.text ?: "") }
                        }
                    }
                }
                resp.feedback?.let {
                    if (it.ok == true) Text(stringResource(R.string.correct), color = MaterialTheme.colorScheme.secondary)
                    it.message?.let { msg -> Text(msg) }
                }
            }

            OutlinedTextField(
                value = answer, onValueChange = { answer = it },
                label = { Text(stringResource(R.string.lesson_answer)) },
                modifier = Modifier.fillMaxWidth(),
                trailingIcon = {
                    IconButton(onClick = {
                        scope.launch {
                            loading = true; error = null
                            try { response = repo.answer(answer, response?.session?.id); answer = "" }
                            catch (e: Exception) { error = e.message } finally { loading = false }
                        }
                    }) { Icon(Icons.AutoMirrored.Filled.Send, ru.neriva.app.NERIVAApp.instance.getString(R.string.send)) }
                }
            )
        }
    }
}
