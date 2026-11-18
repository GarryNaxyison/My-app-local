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
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.launch
import androidx.compose.ui.platform.LocalContext
import ru.neriva.app.LanguageManager
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
    val context = LocalContext.current
    val snackbar = remember { SnackbarHostState() }
    val savedPhrases by app.phrasebookRepo.observeAll().collectAsState(initial = emptyList())
    var savingPhrase by remember { mutableStateOf(false) }
    var answer by remember { mutableStateOf("") }

    LaunchedEffect(Unit) {
        loading = true
        try { response = repo.start() } catch (e: Exception) { error = e.message } finally { loading = false }
    }

    Scaffold(snackbarHost = { SnackbarHost(snackbar) }, topBar = {
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
                            val target = w.target?.trim().orEmpty()
                            if (target.isEmpty()) return@forEach
                            val translation = w.interfaceTranslation?.trim().orEmpty()
                            val example = w.exampleSentenceTarget?.trim().orEmpty()
                            val phrase = example.ifEmpty { target }
                            val saved = savedPhrases.any { it.phrase.equals(phrase, ignoreCase = true) }
                            Card(
                                onClick = {
                                    savingPhrase = true
                                    scope.launch {
                                        val message = try {
                                            app.phrasebookRepo.save(
                                                phrase = phrase,
                                                translation = if (example.isEmpty()) translation else w.exampleTranslationInterface,
                                                note = when {
                                                    example.isEmpty() -> null
                                                    translation.isEmpty() -> target
                                                    else -> context.getString(R.string.tutor_note_word_translation, target, translation)
                                                },
                                                source = "lesson",
                                                language = LanguageManager.getLearningLanguage(context),
                                            )
                                            R.string.tutor_note_saved
                                        } catch (e: CancellationException) {
                                            throw e
                                        } catch (_: Exception) {
                                            R.string.tutor_note_save_failed
                                        } finally {
                                            savingPhrase = false
                                        }
                                        snackbar.showSnackbar(context.getString(message))
                                    }
                                },
                                enabled = !savingPhrase && !saved,
                                modifier = Modifier.fillMaxWidth(),
                            ) {
                                Column(Modifier.padding(12.dp)) {
                                    Text(target, fontWeight = FontWeight.SemiBold)
                                    if (translation.isNotEmpty()) Text(translation, style = MaterialTheme.typography.bodySmall)
                                    if (example.isNotEmpty()) Text(example, style = MaterialTheme.typography.bodySmall)
                                    Text(
                                        stringResource(if (saved) R.string.tutor_note_saved else R.string.save_to_phrasebook),
                                        style = MaterialTheme.typography.labelMedium,
                                        color = MaterialTheme.colorScheme.primary,
                                    )
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
