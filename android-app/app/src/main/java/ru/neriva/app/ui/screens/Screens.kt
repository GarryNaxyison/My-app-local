package ru.neriva.app.ui.screens

import androidx.activity.compose.rememberLauncherForActivityResult
import android.app.Activity
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.itemsIndexed
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody.Companion.toRequestBody
import ru.neriva.app.NERIVAApp
import ru.neriva.app.data.api.*
import ru.neriva.app.data.model.*

// ---- Shared helpers ----

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun ScreenScaffold(
    title: String,
    navController: androidx.navigation.NavHostController,
    content: @Composable (PaddingValues) -> Unit
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(title) },
                navigationIcon = { IconButton(onClick = { navController.popBackStack() }) { Icon(Icons.AutoMirrored.Filled.ArrowBack, "Back") } }
            )
        }
    ) { content(it) }
}

@Composable
private fun LoadingBox() {
    Box(Modifier.fillMaxSize(), Alignment.Center) { CircularProgressIndicator() }
}

// ---- Lesson ----

@Composable
fun LessonScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { ru.neriva.app.data.repo.PracticeRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var prompt by remember { mutableStateOf<String?>(null) }
    var message by remember { mutableStateOf<String?>(null) }
    var answer by remember { mutableStateOf("") }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(Unit) {
        loading = true
        try { val r = repo.startLesson(); prompt = r.prompt; message = r.message } catch (e: Exception) { error = e.message } finally { loading = false }
    }

    ScreenScaffold("Lesson", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            if (loading) return@Column
            prompt?.let { Text(it, style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold) }
            message?.let { Text(it, style = MaterialTheme.typography.bodyLarge) }
            error?.let { Text(it, color = MaterialTheme.colorScheme.error) }
            Spacer(Modifier.weight(1f))
            OutlinedTextField(
                value = answer, onValueChange = { answer = it },
                label = { Text("Your answer") },
                modifier = Modifier.fillMaxWidth()
            )
            Button(
                onClick = {
                    scope.launch {
                        loading = true; error = null
                        try { val r = repo.answerLesson(answer); message = r.message; answer = "" }
                        catch (e: Exception) { error = e.message } finally { loading = false }
                    }
                },
                enabled = !loading && answer.isNotBlank(),
                modifier = Modifier.fillMaxWidth()
            ) { Text("Submit") }
        }
    }
}

// ---- Onboarding ----

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun OnboardingScreen(onComplete: () -> Unit) {
    val app = NERIVAApp.instance
    val scope = rememberCoroutineScope()
    var learningLang by remember { mutableStateOf("en") }
    var level by remember { mutableStateOf("A1") }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }

    Scaffold(topBar = { CenterAlignedTopAppBar(title = { Text("Setup") }) }) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(24.dp).verticalScroll(rememberScrollState()), verticalArrangement = Arrangement.spacedBy(16.dp)) {
            Text("Choose your learning language", style = MaterialTheme.typography.titleMedium)
            val langs = listOf("en", "es", "de", "fr", "it", "zh", "ja", "ko")
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                langs.forEach { l ->
                    FilterChip(selected = learningLang == l, onClick = { learningLang = l }, label = { Text(l.uppercase()) })
                }
            }
            Text("Your level", style = MaterialTheme.typography.titleMedium)
            val levels = listOf("A1", "A2", "B1", "B2", "C1", "C2")
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                levels.forEach { lv ->
                    FilterChip(selected = level == lv, onClick = { level = lv }, label = { Text(lv) })
                }
            }
            error?.let { Text(it, color = MaterialTheme.colorScheme.error) }
            Button(
                onClick = {
                    scope.launch {
                        loading = true; error = null
                        try { app.sessionRepo.updateSettings(learningLanguage = learningLang, level = level); onComplete() }
                        catch (e: Exception) { error = e.message } finally { loading = false }
                    }
                },
                enabled = !loading,
                modifier = Modifier.fillMaxWidth()
            ) { if (loading) CircularProgressIndicator(Modifier.size(20.dp), strokeWidth = 2.dp) else Text("Start Learning") }
        }
    }
}

// ---- Roleplay ----

@Composable
fun RoleplayScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { ru.neriva.app.data.repo.PracticeRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    val scenarios = listOf("Travel", "Work", "Exam", "Speaking")
    var selected by remember { mutableStateOf<String?>(null) }
    var messages by remember { mutableStateOf(listOf<Pair<String, String>>()) }
    var input by remember { mutableStateOf("") }
    var loading by remember { mutableStateOf(false) }

    ScreenScaffold("Roleplay", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp)) {
            if (selected == null) {
                Text("Choose scenario", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                Spacer(Modifier.height(12.dp))
                scenarios.forEach { s ->
                    Card(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp), onClick = { selected = s }) {
                        Text(s, modifier = Modifier.padding(16.dp), style = MaterialTheme.typography.titleMedium)
                    }
                }
            } else {
                LazyColumn(Modifier.weight(1f), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    items(messages) { (role, text) ->
                        Card(modifier = Modifier.fillMaxWidth()) {
                            Column(Modifier.padding(12.dp)) {
                                Text(if (role == "user") "You" else "AI", style = MaterialTheme.typography.labelSmall)
                                Text(text, style = MaterialTheme.typography.bodyMedium)
                            }
                        }
                    }
                }
                Row(Modifier.fillMaxWidth().padding(top = 8.dp)) {
                    OutlinedTextField(value = input, onValueChange = { input = it }, label = { Text("Reply") }, modifier = Modifier.weight(1f))
                    IconButton(onClick = {
                        if (input.isNotBlank() && !loading) {
                            val msg = input; input = ""
                            messages = messages + ("user" to msg)
                            scope.launch {
                                loading = true
                                try { val r = repo.send(message = msg); messages = messages + ("ai" to (r.message ?: "...")) }
                                catch (_: Exception) {} finally { loading = false }
                            }
                        }
                    }) { Icon(Icons.Default.Send, "Send") }
                }
            }
        }
    }
}

// ---- Shadowing ----

@Composable
fun ShadowingScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { ru.neriva.app.data.repo.PracticeRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var text by remember { mutableStateOf<String?>(null) }
    var loading by remember { mutableStateOf(true) }
    var isRecording by remember { mutableStateOf(false) }
    var feedback by remember { mutableStateOf<String?>(null) }
    var xp by remember { mutableIntStateOf(0) }
    val recorder = remember { ru.neriva.app.audio.AudioRecorder(app) }

    LaunchedEffect(Unit) { try { text = repo.startShadowing().text } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Listening", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), horizontalAlignment = Alignment.CenterHorizontally, verticalArrangement = Arrangement.Center) {
            if (loading) CircularProgressIndicator()
            else {
                Text(text ?: "Listen and repeat", style = MaterialTheme.typography.headlineSmall, fontWeight = FontWeight.Bold)
                Spacer(Modifier.height(24.dp))
                feedback?.let {
                    Card(Modifier.fillMaxWidth()) {
                        Column(Modifier.padding(16.dp), horizontalAlignment = Alignment.CenterHorizontally) {
                            Text(it, style = MaterialTheme.typography.bodyMedium)
                            if (xp > 0) Text("+$xp XP", color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.Bold)
                        }
                    }
                    Spacer(Modifier.height(16.dp))
                }
                FilledIconButton(
                    onClick = {
                        if (isRecording) {
                            val file = recorder.stopRecording()
                            isRecording = false
                            if (file != null) {
                                loading = true
                                scope.launch {
                                    try {
                                        val r = repo.answerShadowingVoice(file, text = text)
                                        feedback = r.text ?: "Good attempt!"
                                        xp += r.xp ?: 0
                                    } catch (_: Exception) {} finally { loading = false }
                                }
                            }
                        } else {
                            recorder.startRecording()
                            isRecording = true
                        }
                    },
                    modifier = Modifier.size(72.dp),
                    colors = IconButtonDefaults.filledIconButtonColors(
                        containerColor = if (isRecording) MaterialTheme.colorScheme.error else MaterialTheme.colorScheme.primary
                    )
                ) {
                    Icon(Icons.Default.Mic, if (isRecording) "Stop" else "Record", modifier = Modifier.size(36.dp))
                }
                Text(if (isRecording) "Recording... tap to stop" else "Tap to record", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
            }
        }
    }
}

// ---- Pronunciation ----

@Composable
fun PronunciationScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { ru.neriva.app.data.repo.PracticeRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var text by remember { mutableStateOf<String?>(null) }
    var loading by remember { mutableStateOf(true) }
    var isRecording by remember { mutableStateOf(false) }
    var score by remember { mutableIntStateOf(0) }
    var feedback by remember { mutableStateOf<String?>(null) }
    var weakWords by remember { mutableStateOf<List<String>>(emptyList()) }
    var checking by remember { mutableStateOf(false) }
    val recorder = remember { ru.neriva.app.audio.AudioRecorder(app) }

    LaunchedEffect(Unit) { try { text = repo.startPronunciation().text } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Pronunciation", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), horizontalAlignment = Alignment.CenterHorizontally, verticalArrangement = Arrangement.Center) {
            if (loading) CircularProgressIndicator()
            else {
                Text("Repeat:", style = MaterialTheme.typography.labelLarge)
                Text(text ?: "", style = MaterialTheme.typography.headlineSmall, fontWeight = FontWeight.Bold)
                Spacer(Modifier.height(24.dp))

                if (score > 0) {
                    Card(Modifier.fillMaxWidth()) {
                        Column(Modifier.padding(16.dp), horizontalAlignment = Alignment.CenterHorizontally) {
                            Text("Score: $score/100", style = MaterialTheme.typography.headlineMedium, fontWeight = FontWeight.Bold,
                                color = if (score >= 80) MaterialTheme.colorScheme.secondary else if (score >= 60) MaterialTheme.colorScheme.tertiary else MaterialTheme.colorScheme.error)
                            feedback?.let { Text(it, style = MaterialTheme.typography.bodyMedium) }
                            if (weakWords.isNotEmpty()) {
                                Text("Weak words:", style = MaterialTheme.typography.labelLarge)
                                Text(weakWords.joinToString(", "), style = MaterialTheme.typography.bodySmall)
                            }
                        }
                    }
                    Spacer(Modifier.height(16.dp))
                }

                if (checking) {
                    CircularProgressIndicator()
                } else {
                    FilledIconButton(
                        onClick = {
                            if (isRecording) {
                                val file = recorder.stopRecording()
                                isRecording = false
                                if (file != null) {
                                    checking = true
                                    scope.launch {
                                        try {
                                            val result = repo.checkPronunciationVoice(file)
                                            score = result.score ?: 0
                                            feedback = result.feedback
                                            weakWords = result.weakWords ?: emptyList()
                                        } catch (_: Exception) {} finally { checking = false }
                                    }
                                }
                            } else {
                                recorder.startRecording()
                                isRecording = true
                                score = 0; feedback = null; weakWords = emptyList()
                            }
                        },
                        modifier = Modifier.size(72.dp),
                        colors = IconButtonDefaults.filledIconButtonColors(
                            containerColor = if (isRecording) MaterialTheme.colorScheme.error else MaterialTheme.colorScheme.primary
                        )
                    ) {
                        Icon(Icons.Default.Mic, if (isRecording) "Stop" else "Record", modifier = Modifier.size(36.dp))
                    }
                    Text(if (isRecording) "Recording... tap to stop" else "Tap to record", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                }
            }
        }
    }
}

// ---- Word Game ----

@Composable
fun WordGameScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { ru.neriva.app.data.repo.VocabularyRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var challenge by remember { mutableStateOf<WordChallenge?>(null) }
    var result by remember { mutableStateOf<String?>(null) }
    var xp by remember { mutableIntStateOf(0) }
    var loading by remember { mutableStateOf(true) }
    LaunchedEffect(Unit) { try { challenge = repo.nextGame() } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Review Game", navController) { padding ->
        Box(Modifier.padding(padding).fillMaxSize().padding(16.dp)) {
            if (loading) CircularProgressIndicator(Modifier.align(Alignment.Center))
            else challenge?.let { ch ->
                Column(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                    Text(ch.prompt ?: ch.word ?: "", style = MaterialTheme.typography.titleLarge, fontWeight = FontWeight.Bold)
                    ch.context?.let { Text(it, style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant) }
                    result?.let { Text(it, color = if (it.startsWith("Correct")) MaterialTheme.colorScheme.secondary else MaterialTheme.colorScheme.error) }
                    if (xp > 0) Text("Total: +$xp XP", color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.Bold)
                    Spacer(Modifier.height(8.dp))
                    ch.options.forEach { opt ->
                        Button(
                            onClick = {
                                scope.launch {
                                    loading = true; result = null
                                    try {
                                        val resp = repo.answerGame(opt.text, wordId = ch.wordId)
                                        result = if (resp.correct) "Correct! +${resp.xp ?: 10} XP" else "Wrong: ${resp.correctAnswer ?: ""}"
                                        if (resp.correct) xp += resp.xp ?: 10
                                    } catch (_: Exception) {} finally { loading = false }
                                }
                            },
                            enabled = !loading,
                            modifier = Modifier.fillMaxWidth()
                        ) { Text(opt.text) }
                    }
                }
            }
        }
    }
}

// ---- Spelling ----

@Composable
fun SpellingScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { ru.neriva.app.data.repo.VocabularyRepository(app.apiClient) }
    var ch by remember { mutableStateOf<SpellingChallenge?>(null) }
    var answer by remember { mutableStateOf("") }
    var result by remember { mutableStateOf<String?>(null) }
    var attempts by remember { mutableIntStateOf(0) }
    var hint by remember { mutableStateOf<String?>(null) }
    var showHint by remember { mutableStateOf(false) }
    var gaveUp by remember { mutableStateOf(false) }
    var loading by remember { mutableStateOf(true) }
    var xp by remember { mutableIntStateOf(0) }
    val scope = rememberCoroutineScope()
    LaunchedEffect(Unit) { try { ch = repo.startSpelling() } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Spelling", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            if (loading) { CircularProgressIndicator(); return@Column }
            ch?.let { c ->
                Text(c.prompt ?: c.message ?: "", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                c.context?.let { Text(it, style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant) }
                if (attempts > 0) Text("Attempts: $attempts", style = MaterialTheme.typography.bodySmall)
                result?.let { Text(it, color = if (it?.startsWith("Correct") == true) MaterialTheme.colorScheme.secondary else MaterialTheme.colorScheme.error) }
                if (xp > 0) Text("+$xp XP", color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.Bold)
                if (showHint && !gaveUp) {
                    Card(Modifier.fillMaxWidth(), colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.tertiaryContainer)) {
                        Text("Hint: starts with '${c.wordId?.firstOrNull() ?: "?"}'  (${c.wordId?.length ?: 0} letters)", modifier = Modifier.padding(12.dp))
                    }
                }
                if (gaveUp) {
                    Card(Modifier.fillMaxWidth(), colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.errorContainer)) {
                        Text("Answer: ${c.wordId ?: ""}", modifier = Modifier.padding(12.dp), fontWeight = FontWeight.Bold)
                    }
                }
                Spacer(Modifier.weight(1f))
                OutlinedTextField(value = answer, onValueChange = { answer = it }, label = { Text("Spell the word") }, modifier = Modifier.fillMaxWidth())
                Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    Button(onClick = {
                        scope.launch {
                            loading = true; attempts++
                            try {
                                val r = repo.answerSpelling(answer, wordId = c.wordId)
                                if (r.correct) {
                                    result = "Correct! +${r.xp ?: 7} XP"
                                    xp += r.xp ?: 7
                                } else {
                                    result = "Wrong: ${r.correctAnswer ?: ""}"
                                }
                                answer = ""
                            } catch (_: Exception) {} finally { loading = false }
                        }
                    }, enabled = !loading && answer.isNotBlank() && !gaveUp, modifier = Modifier.weight(1f)) { Text("Check") }
                    if (!showHint && !gaveUp) {
                        OutlinedButton(onClick = { showHint = true }) { Text("Hint") }
                    }
                    if (!gaveUp && attempts >= 2) {
                        OutlinedButton(onClick = { gaveUp = true; result = "The word was: ${c.wordId ?: ""}" }) { Text("Give up") }
                    }
                }
            }
        }
    }
}

// ---- Vocabulary ----

@Composable
fun VocabularyScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { ru.neriva.app.data.repo.VocabularyRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var words by remember { mutableStateOf<List<VocabularyItem>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var reportTarget by remember { mutableStateOf<VocabularyItem?>(null) }
    var reportPending by remember { mutableStateOf(false) }
    var reportError by remember { mutableStateOf<String?>(null) }

    LaunchedEffect(Unit) { try { words = repo.getList(1).items ?: emptyList() } catch (_: Exception) {} finally { loading = false } }

    reportTarget?.let { word ->
        ru.neriva.app.ui.components.WordReportDialog(
            word = word.word,
            translation = word.translation ?: "",
            pending = reportPending,
            error = reportError,
            onDismiss = {
                reportTarget = null
                reportError = null
            },
            onSubmit = { comment ->
                reportPending = true
                reportError = null
                scope.launch {
                    try {
                        // Report via word id if available; otherwise send word as id.
                        val wid = (word as? ru.neriva.app.data.model.VocabularyItem)?.id ?: word.word
                        repo.reportWord(wid)
                        reportTarget = null
                    } catch (e: Exception) {
                        reportError = e.message ?: "Failed to send report"
                    } finally { reportPending = false }
                }
            },
        )
    }

    ScreenScaffold("Vocabulary", navController) { padding ->
        if (loading) { Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) { CircularProgressIndicator() } }
        else LazyColumn(Modifier.padding(padding).fillMaxSize().padding(horizontal = 16.dp), verticalArrangement = Arrangement.spacedBy(8.dp), contentPadding = PaddingValues(vertical = 16.dp)) {
            items(words) { w ->
                Card(Modifier.fillMaxWidth()) {
                    Column(Modifier.padding(12.dp)) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Text(w.word, fontWeight = FontWeight.SemiBold, modifier = Modifier.weight(1f))
                            IconButton(onClick = { reportTarget = w }) {
                                Icon(Icons.Default.ReportProblem, contentDescription = "Report word", tint = MaterialTheme.colorScheme.onSurfaceVariant)
                            }
                        }
                        w.translation?.let { Text(it, style = MaterialTheme.typography.bodyMedium) }
                        w.example?.let { Text(it, style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant) }
                    }
                }
            }
        }
    }
}

// ---- Phrasebook ----

@Composable
fun PhrasebookScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val scope = rememberCoroutineScope()
    val phrases by app.phrasebookRepo.observeAll().collectAsState(initial = emptyList())
    var phrase by remember { mutableStateOf("") }
    var translation by remember { mutableStateOf("") }
    var note by remember { mutableStateOf("") }
    var loading by remember { mutableStateOf(true) }
    var saving by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }

    fun refresh() {
        scope.launch {
            loading = true
            error = null
            try {
                app.phrasebookRepo.refresh()
            } catch (e: Exception) {
                error = e.message ?: "Could not update the phrasebook."
            } finally {
                loading = false
            }
        }
    }

    LaunchedEffect(Unit) { refresh() }

    ScreenScaffold("Phrasebook", navController) { padding ->
        LazyColumn(
            Modifier.padding(padding).fillMaxSize().padding(horizontal = 16.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp),
            contentPadding = PaddingValues(vertical = 16.dp),
        ) {
            item {
                Card(Modifier.fillMaxWidth()) {
                    Column(Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(10.dp)) {
                        Text("Add phrase", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                        OutlinedTextField(
                            value = phrase,
                            onValueChange = { phrase = it },
                            label = { Text("Phrase") },
                            modifier = Modifier.fillMaxWidth(),
                            singleLine = true,
                        )
                        OutlinedTextField(
                            value = translation,
                            onValueChange = { translation = it },
                            label = { Text("Translation") },
                            modifier = Modifier.fillMaxWidth(),
                            singleLine = true,
                        )
                        OutlinedTextField(
                            value = note,
                            onValueChange = { note = it },
                            label = { Text("Note") },
                            modifier = Modifier.fillMaxWidth(),
                        )
                        Button(
                            onClick = {
                                scope.launch {
                                    saving = true
                                    error = null
                                    try {
                                        app.phrasebookRepo.save(
                                            phrase = phrase.trim(),
                                            translation = translation.trim().ifBlank { null },
                                            note = note.trim().ifBlank { null },
                                            source = "manual",
                                            language = ru.neriva.app.LanguageManager.getLearningLanguage(NERIVAApp.context()),
                                        )
                                        phrase = ""
                                        translation = ""
                                        note = ""
                                    } catch (e: Exception) {
                                        error = e.message ?: "Could not save the phrase."
                                    } finally {
                                        saving = false
                                    }
                                }
                            },
                            enabled = phrase.isNotBlank() && !saving,
                            modifier = Modifier.fillMaxWidth(),
                        ) { Text(if (saving) "Saving..." else "Save phrase") }
                    }
                }
            }
            if (loading && phrases.isEmpty()) {
                item { Box(Modifier.fillMaxWidth().padding(24.dp), Alignment.Center) { CircularProgressIndicator() } }
            }
            error?.let { message ->
                item {
                    Card(colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.errorContainer)) {
                        Row(Modifier.padding(12.dp), verticalAlignment = Alignment.CenterVertically) {
                            Text(message, color = MaterialTheme.colorScheme.onErrorContainer, modifier = Modifier.weight(1f))
                            TextButton(onClick = ::refresh) { Text("Retry") }
                        }
                    }
                }
            }
            item {
                Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween, verticalAlignment = Alignment.CenterVertically) {
                    Text("Saved phrases", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                    TextButton(onClick = ::refresh, enabled = !loading) { Text("Refresh") }
                }
            }
            if (!loading && phrases.isEmpty()) {
                item { Text("No phrases yet. Add one above to keep it on all your devices.", color = MaterialTheme.colorScheme.onSurfaceVariant) }
            }
            items(phrases, key = { it.id }) { item ->
                Card(Modifier.fillMaxWidth()) {
                    Row(Modifier.padding(12.dp), verticalAlignment = Alignment.Top) {
                        Column(Modifier.weight(1f)) {
                            Text(item.phrase, fontWeight = FontWeight.SemiBold)
                            item.translation?.takeIf { it.isNotBlank() }?.let { Text(it, style = MaterialTheme.typography.bodyMedium) }
                            item.note?.takeIf { it.isNotBlank() }?.let { Text(it, style = MaterialTheme.typography.bodySmall) }
                        }
                        IconButton(
                            onClick = {
                                scope.launch {
                                    saving = true
                                    error = null
                                    try {
                                        app.phrasebookRepo.delete(item.id)
                                    } catch (e: Exception) {
                                        error = e.message ?: "Could not delete the phrase."
                                    } finally {
                                        saving = false
                                    }
                                }
                            },
                            enabled = !saving,
                        ) { Icon(Icons.Default.Delete, contentDescription = "Delete phrase") }
                    }
                }
            }
        }
    }
}

// ---- Offline ----

@Composable
fun OfflineScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    var cached by remember { mutableStateOf<List<ru.neriva.app.data.db.PhrasebookEntity>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    LaunchedEffect(Unit) { try { cached = app.phrasebookRepo.getCached() } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Offline Decks", navController) { padding ->
        if (loading) { Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) { CircularProgressIndicator() } }
        else LazyColumn(Modifier.padding(padding).fillMaxSize().padding(horizontal = 16.dp), verticalArrangement = Arrangement.spacedBy(8.dp), contentPadding = PaddingValues(vertical = 16.dp)) {
            items(cached) { p ->
                Card(Modifier.fillMaxWidth()) {
                    Column(Modifier.padding(12.dp)) {
                        Text(p.phrase, fontWeight = FontWeight.SemiBold)
                        p.translation?.let { Text(it, style = MaterialTheme.typography.bodyMedium) }
                    }
                }
            }
            if (cached.isEmpty()) item { Text("No offline cards yet. Save phrases to learn offline.", color = MaterialTheme.colorScheme.onSurfaceVariant) }
        }
    }
}

// ---- Level Test ----

@Composable
fun LevelTestScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val scope = rememberCoroutineScope()
    var question by remember { mutableStateOf<LevelQuestion?>(null) }
    var loading by remember { mutableStateOf(true) }
    var finished by remember { mutableStateOf(false) }
    var resultLevel by remember { mutableStateOf<String?>(null) }
    var resultScore by remember { mutableIntStateOf(0) }
    LaunchedEffect(Unit) { try { question = app.apiClient.levelTestStart().question } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Level Test", navController) { padding ->
        Box(Modifier.padding(padding).fillMaxSize().padding(16.dp)) {
            if (loading) CircularProgressIndicator(Modifier.align(Alignment.Center))
            else if (finished) {
                Column(horizontalAlignment = Alignment.CenterHorizontally, verticalArrangement = Arrangement.spacedBy(12.dp), modifier = Modifier.align(Alignment.Center)) {
                    Text("Test Complete!", style = MaterialTheme.typography.headlineMedium, fontWeight = FontWeight.Bold)
                    resultLevel?.let { Text("Your level: $it", style = MaterialTheme.typography.titleLarge, color = MaterialTheme.colorScheme.primary) }
                    if (resultScore > 0) Text("Score: $resultScore", style = MaterialTheme.typography.titleMedium)
                }
            }
            else question?.let { q ->
                Column(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                    Text("Question ${q.index}/${q.total}", style = MaterialTheme.typography.labelLarge)
                    Text(q.question, style = MaterialTheme.typography.titleMedium)
                    q.options.forEach { opt ->
                        Button(
                            onClick = {
                                scope.launch {
                                    loading = true
                                    try {
                                        val resp = app.apiClient.levelTestAnswer(ru.neriva.app.data.api.AnswerRequest(answer = opt))
                                        question = resp.question
                                        resp.level?.let { resultLevel = it }
                                        resp.score?.let { resultScore = it }
                                        if (resp.finished) finished = true
                                    } catch (_: Exception) {} finally { loading = false }
                                }
                            },
                            modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp)
                        ) { Text(opt) }
                    }
                }
            }
        }
    }
}

// ---- Progress ----

@Composable
fun ProgressScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    var progress by remember { mutableStateOf<ProgressResponse?>(null) }
    var loading by remember { mutableStateOf(true) }
    LaunchedEffect(Unit) { try { progress = app.sessionRepo.getProgress() } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Progress", navController) { padding ->
        if (loading) { Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) { CircularProgressIndicator() } }
        else Column(Modifier.padding(padding).fillMaxSize().padding(16.dp).verticalScroll(rememberScrollState()), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            progress?.let { p ->
                MetricCard("XP", "${p.xp}")
                MetricCard("Level", "${p.xpLevel}")
                MetricCard("Title", p.xpTitle ?: "-")
                MetricCard("Lessons", "${p.lessonCount}")
                MetricCard("Practice", "${p.practiceCount}")
                MetricCard("Words learned", "${p.learnedWords}")
                MetricCard("Mistakes", "${p.mistakes}")
                p.streak?.let { MetricCard("Streak", "$it days") }
                if (p.premium) MetricCard("Status", "Premium")
            }
        }
    }
}

@Composable
private fun MetricCard(label: String, value: String) {
    Card(Modifier.fillMaxWidth()) {
        Row(Modifier.padding(16.dp).fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
            Text(label, style = MaterialTheme.typography.bodyMedium)
            Text(value, style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
        }
    }
}

// ---- Awards ----

@Composable
fun AwardsScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    var user by remember { mutableStateOf<UserProfile?>(null) }
    var loading by remember { mutableStateOf(true) }
    LaunchedEffect(Unit) { try { user = app.sessionRepo.getSession().user } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Awards", navController) { padding ->
        if (loading) { Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) { CircularProgressIndicator() } }
        else {
            val currentLevel = (user?.xpLevel ?: 1).coerceIn(1, 20)
            LazyColumn(
                modifier = Modifier.padding(padding).fillMaxSize().padding(16.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(12.dp),
            ) {
                // Summary card
                item {
                    Column(horizontalAlignment = Alignment.CenterHorizontally, verticalArrangement = Arrangement.spacedBy(8.dp)) {
                        ru.neriva.app.ui.components.AwardIcon(level = currentLevel, size = 120.dp)
                        Text("Level $currentLevel / 20", style = MaterialTheme.typography.headlineMedium, fontWeight = FontWeight.Bold)
                        Text(user?.xpTitle ?: "", style = MaterialTheme.typography.titleMedium, color = MaterialTheme.colorScheme.onSurfaceVariant, textAlign = androidx.compose.ui.text.style.TextAlign.Center)
                        Spacer(Modifier.height(8.dp))
                        Text("${user?.xp ?: 0} XP total", style = MaterialTheme.typography.bodyLarge)
                    }
                }
                // Grid of 20 award tiles
                item {
                    Text("All trophies", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold, modifier = Modifier.fillMaxWidth())
                }
                items((1..20).toList()) { level ->
                    val unlocked = currentLevel >= level
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.spacedBy(12.dp),
                    ) {
                        ru.neriva.app.ui.components.AwardIcon(
                            level = level,
                            size = 56.dp,
                            grayscale = !unlocked,
                        )
                        Column(modifier = Modifier.weight(1f)) {
                            Text("Level $level", style = MaterialTheme.typography.titleSmall, fontWeight = FontWeight.SemiBold)
                            Text(
                                if (unlocked) "Unlocked" else "Locked",
                                style = MaterialTheme.typography.bodySmall,
                                color = if (unlocked) MaterialTheme.colorScheme.secondary else MaterialTheme.colorScheme.onSurfaceVariant,
                            )
                        }
                    }
                }
            }
        }
    }
}

// ---- Leaderboard ----

@Composable
fun LeaderboardScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    var entries by remember { mutableStateOf<List<LeaderboardEntry>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    LaunchedEffect(Unit) { try { entries = app.sessionRepo.getLeaderboard().items ?: emptyList() } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Leaderboard", navController) { padding ->
        if (loading) { Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) { CircularProgressIndicator() } }
        else LazyColumn(Modifier.padding(padding).fillMaxSize().padding(horizontal = 16.dp), verticalArrangement = Arrangement.spacedBy(8.dp), contentPadding = PaddingValues(vertical = 16.dp)) {
            items(entries) { e ->
                Card(Modifier.fillMaxWidth()) {
                    Row(Modifier.padding(12.dp).fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                        Text(e.name ?: "—", fontWeight = FontWeight.SemiBold)
                        Text("${e.xp ?: 0} XP", color = MaterialTheme.colorScheme.primary)
                    }
                }
            }
        }
    }
}

// ---- Limits ----

@Composable
fun LimitsScreen(navController: androidx.navigation.NavHostController) {
    val vm: ru.neriva.app.ui.SessionViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
    val state by vm.state.collectAsState()
    val u = state.session?.user

    ScreenScaffold("Limits", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            MetricCard("Lessons", "${u?.lessonsToday ?: 0} / ${u?.lessonLimit ?: 0}")
            MetricCard("Practice", "${u?.practiceToday ?: 0} / ${u?.practiceLimit ?: 0}")
            if ((u?.voiceLimit ?: 0) > 0) MetricCard("Voice", "${u?.voiceToday ?: 0} / ${u?.voiceLimit ?: 0}")
            if (u?.premium != true) {
                Button(onClick = { navController.navigate("premium") }, modifier = Modifier.fillMaxWidth()) { Text("Upgrade to Premium") }
            }
        }
    }
}

// ---- Mistakes ----

@Composable
fun MistakesScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { ru.neriva.app.data.repo.MistakesRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    var mistakes by remember { mutableStateOf<List<MistakeItem>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    var busy by remember { mutableStateOf(false) }
    var confirmClear by remember { mutableStateOf(false) }

    suspend fun reload() {
        loading = true
        try { mistakes = repo.getAll().items ?: emptyList() } catch (_: Exception) {}
        loading = false
    }
    LaunchedEffect(Unit) { reload() }

    if (confirmClear) {
        AlertDialog(
            onDismissRequest = { confirmClear = false },
            title = { Text("Clear all mistakes?") },
            text = { Text("This will remove all ${mistakes.size} saved mistakes.") },
            confirmButton = {
                Button(
                    onClick = {
                        confirmClear = false
                        scope.launch {
                            busy = true
                            try { repo.clearAll(); reload() } catch (_: Exception) {} finally { busy = false }
                        }
                    },
                    colors = ButtonDefaults.buttonColors(containerColor = MaterialTheme.colorScheme.error),
                ) { Text("Yes, clear") }
            },
            dismissButton = { TextButton(onClick = { confirmClear = false }) { Text("Cancel") } },
        )
    }

    ScreenScaffold("Mistakes", navController) { padding ->
        if (loading) { Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) { CircularProgressIndicator() } }
        else LazyColumn(Modifier.padding(padding).fillMaxSize().padding(horizontal = 16.dp), verticalArrangement = Arrangement.spacedBy(8.dp), contentPadding = PaddingValues(vertical = 16.dp)) {
            item {
                Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    if (mistakes.isNotEmpty()) {
                        Button(
                            onClick = { navController.navigate("mistake-practice") },
                            modifier = Modifier.weight(1f),
                            enabled = !busy,
                        ) { Text("Practice") }
                        OutlinedButton(
                            onClick = { confirmClear = true },
                            modifier = Modifier.weight(1f),
                            enabled = !busy,
                        ) { Text("Clear all") }
                    }
                }
            }
            itemsIndexed(mistakes) { index, m ->
                Card(Modifier.fillMaxWidth()) {
                    Column(Modifier.padding(12.dp)) {
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Column(Modifier.weight(1f)) {
                                Text(m.word ?: "", fontWeight = FontWeight.SemiBold)
                                Text("→ ${m.correction ?: ""}", color = MaterialTheme.colorScheme.error)
                                m.explanation?.let { Text(it, style = MaterialTheme.typography.bodySmall) }
                            }
                            IconButton(onClick = {
                                scope.launch {
                                    busy = true
                                    try { repo.delete((index + 1).toString()); reload() } catch (_: Exception) {} finally { busy = false }
                                }
                            }) {
                                Icon(Icons.Default.Delete, contentDescription = "Delete", tint = MaterialTheme.colorScheme.onSurfaceVariant)
                            }
                        }
                    }
                }
            }
            if (mistakes.isEmpty()) {
                item { Text("No mistakes yet. Practice to fill this list.", color = MaterialTheme.colorScheme.onSurfaceVariant) }
            }
        }
    }
}

// ---- Tools ----

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ToolsScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { ru.neriva.app.data.repo.ToolsRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    val context = androidx.compose.ui.platform.LocalContext.current

    var mode by remember { mutableStateOf("translator") } // translator | voice | image
    var input by remember { mutableStateOf("") }
    var result by remember { mutableStateOf<String?>(null) }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }

    // Voice recorder for voice-to-text mode
    val recorder = remember { ru.neriva.app.audio.AudioRecorder(context) }
    var recording by remember { mutableStateOf(false) }
    var voiceFile by remember { mutableStateOf<java.io.File?>(null) }

    // Image picker for image-translate mode
    var imageUri by remember { mutableStateOf<android.net.Uri?>(null) }
    val imagePicker = rememberLauncherForActivityResult(
        contract = androidx.activity.result.contract.ActivityResultContracts.GetContent()
    ) { uri: android.net.Uri? -> imageUri = uri }

    val vm: ru.neriva.app.ui.SessionViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
    val state by vm.state.collectAsState()
    val session = state.session
    val langs: List<ru.neriva.app.data.model.LanguageOption> = remember(session) {
        session?.learning_languages?.takeIf { it.isNotEmpty() }
            ?: session?.interface_languages
            ?: emptyList()
    }
    val langCodes: List<String> = remember(langs) {
        langs.map { it.code }.ifEmpty { listOf("en", "ru", "es", "de", "fr") }
    }
    var targetLang by remember { mutableStateOf("ru") }

    fun submit() {
        scope.launch {
            loading = true; error = null
            try {
                when (mode) {
                    "translator" -> {
                        val r = repo.translate(input, sourceLanguage = "auto", targetLanguage = targetLang)
                        result = r.translation ?: r.result ?: r.sourceText
                    }
                    "voice" -> {
                        val f = voiceFile ?: throw IllegalStateException("Record a voice message first")
                        val r = repo.voiceToTextFile(f, targetLanguage = targetLang)
                        result = r.translation ?: r.result ?: r.sourceText
                        voiceFile = null
                    }
                    "image" -> {
                        val uri = imageUri ?: throw IllegalStateException("Pick an image first")
                        val tmp = copyUriToTempFile(context, uri, "tool-image")
                        val r = repo.imageTranslate(tmp, targetLang)
                        result = r.translation ?: r.result ?: r.sourceText
                        imageUri = null
                    }
                }
            } catch (e: Exception) {
                error = e.message ?: "Tool failed"
            } finally { loading = false }
        }
    }

    ScreenScaffold("Tools", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp).verticalScroll(rememberScrollState()), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            // Mode tabs (mirrors web tool-switch)
            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                listOf("translator" to "Text", "voice" to "Voice", "image" to "Photo").forEach { (m, label) ->
                    FilterChip(
                        selected = mode == m,
                        onClick = { mode = m; result = null; error = null },
                        label = { Text(label) },
                        modifier = Modifier.weight(1f),
                    )
                }
            }

            val title = when (mode) {
                "translator" -> "Quick translator"
                "voice" -> "Voice to text"
                else -> "Photo translation"
            }
            Text(title, style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)

            when (mode) {
                "translator" -> {
                    OutlinedTextField(
                        value = input, onValueChange = { input = it },
                        label = { Text("Paste text to translate") },
                        modifier = Modifier.fillMaxWidth(),
                        minLines = 3,
                    )
                }
                "voice" -> {
                    Text(
                        if (recording) "Recording… tap to stop" else (voiceFile?.let { "Voice ready (${it.name})" } ?: "Record your voice"),
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                    Button(
                        onClick = {
                            if (recording) {
                                try { voiceFile = recorder.stopRecording() } catch (_: Exception) {}
                                recording = false
                            } else {
                                try { recorder.startRecording(); recording = true } catch (_: Exception) {}
                            }
                        },
                        modifier = Modifier.fillMaxWidth(),
                        colors = ButtonDefaults.buttonColors(containerColor = if (recording) MaterialTheme.colorScheme.error else MaterialTheme.colorScheme.primary),
                    ) {
                        Icon(
                            androidx.compose.material.icons.Icons.Default.Mic,
                            contentDescription = null,
                            modifier = Modifier.size(20.dp),
                        )
                        Spacer(Modifier.width(8.dp))
                        Text(if (recording) "Stop" else "Record")
                    }
                }
                "image" -> {
                    imageUri?.let {
                        Text("Image selected", style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant)
                    } ?: Text("Pick a photo containing text to translate", style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant)
                    OutlinedButton(onClick = { imagePicker.launch("image/*") }, modifier = Modifier.fillMaxWidth()) {
                        Icon(androidx.compose.material.icons.Icons.Default.Add, contentDescription = null, modifier = Modifier.size(20.dp))
                        Spacer(Modifier.width(8.dp))
                        Text("Choose image")
                    }
                }
            }

            // Target language selector
            Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                Text("Target:", style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant)
                var langExpanded by remember { mutableStateOf(false) }
                ExposedDropdownMenuBox(expanded = langExpanded, onExpandedChange = { langExpanded = it }) {
                    OutlinedTextField(
                        value = targetLang.uppercase(),
                        onValueChange = {},
                        readOnly = true,
                        trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(langExpanded) },
                        modifier = Modifier.fillMaxWidth().menuAnchor(),
                    )
                    ExposedDropdownMenu(expanded = langExpanded, onDismissRequest = { langExpanded = false }) {
                        langCodes.forEach { code ->
                            DropdownMenuItem(text = { Text(code.uppercase()) }, onClick = { targetLang = code; langExpanded = false })
                        }
                    }
                }
            }

            Button(
                onClick = { submit() },
                enabled = when (mode) {
                    "translator" -> !loading && input.isNotBlank()
                    "voice" -> !loading && voiceFile != null
                    else -> !loading && imageUri != null
                },
                modifier = Modifier.fillMaxWidth(),
            ) {
                if (loading) CircularProgressIndicator(Modifier.size(20.dp), strokeWidth = 2.dp)
                else Text("Send")
            }

            error?.let { Text(it, color = MaterialTheme.colorScheme.error) }
            result?.let {
                Card(Modifier.fillMaxWidth()) {
                    Column(Modifier.padding(16.dp)) {
                        Text("Result", style = MaterialTheme.typography.labelLarge, color = MaterialTheme.colorScheme.onSurfaceVariant)
                        Spacer(Modifier.height(4.dp))
                        Text(it, style = MaterialTheme.typography.bodyLarge)
                    }
                }
            }
        }
    }
}

// ---- Premium ----

@Composable
fun PremiumScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    var plans by remember { mutableStateOf<List<PremiumPlan>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    LaunchedEffect(Unit) { try { plans = app.sessionRepo.getPremiumPlans().plans ?: emptyList() } catch (_: Exception) {} finally { loading = false } }

    ScreenScaffold("Premium", navController) { padding ->
        if (loading) { Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) { CircularProgressIndicator() } }
        else LazyColumn(Modifier.padding(padding).fillMaxSize().padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            items(plans) { plan ->
                Card(Modifier.fillMaxWidth(), colors = CardDefaults.cardColors(containerColor = if (plan.current) MaterialTheme.colorScheme.tertiaryContainer else MaterialTheme.colorScheme.surfaceVariant)) {
                    Column(Modifier.padding(16.dp)) {
                        Text(plan.title ?: plan.product, style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.Bold)
                        plan.body?.let { Text(it, style = MaterialTheme.typography.bodyMedium) }
                        Spacer(Modifier.height(8.dp))
                        plan.features.take(4).forEach { Text("• $it", style = MaterialTheme.typography.bodySmall) }
                        plan.rubPrice?.let { Text("$it RUB", style = MaterialTheme.typography.titleMedium, color = MaterialTheme.colorScheme.primary) }
                    }
                }
            }
        }
    }
}

// ---- Settings ----

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SettingsScreen(navController: androidx.navigation.NavHostController, onThemeChange: (ru.neriva.app.ThemeManager.ThemeMode) -> Unit = {}, currentThemeMode: ru.neriva.app.ThemeManager.ThemeMode = ru.neriva.app.ThemeManager.ThemeMode.SYSTEM) {
    val app = NERIVAApp.instance
    val context = androidx.compose.ui.platform.LocalContext.current
    val scope = rememberCoroutineScope()
    val vm: ru.neriva.app.ui.SessionViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
    val state by vm.state.collectAsState()
    val user = state.session?.user

    var interfaceLangExpanded by remember { mutableStateOf(false) }
    var learningLangExpanded by remember { mutableStateOf(false) }

    val interfaceLang = ru.neriva.app.LanguageManager.getInterfaceLanguage(context)
    val learningLang = ru.neriva.app.LanguageManager.getLearningLanguage(context)
    val interfaceLocales = ru.neriva.app.LanguageManager.getSupportedLocales()
    val learningLocales = ru.neriva.app.LanguageManager.getSupportedLearningLanguages()

    ScreenScaffold("Settings", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp).verticalScroll(rememberScrollState()), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            // Theme toggle
            Text("Theme", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                FilterChip(selected = currentThemeMode == ru.neriva.app.ThemeManager.ThemeMode.LIGHT, onClick = { onThemeChange(ru.neriva.app.ThemeManager.ThemeMode.LIGHT) }, label = { Text("Light") }, leadingIcon = { Icon(Icons.Default.LightMode, null, modifier = Modifier.size(18.dp)) })
                FilterChip(selected = currentThemeMode == ru.neriva.app.ThemeManager.ThemeMode.DARK, onClick = { onThemeChange(ru.neriva.app.ThemeManager.ThemeMode.DARK) }, label = { Text("Dark") }, leadingIcon = { Icon(Icons.Default.DarkMode, null, modifier = Modifier.size(18.dp)) })
                FilterChip(selected = currentThemeMode == ru.neriva.app.ThemeManager.ThemeMode.SYSTEM, onClick = { onThemeChange(ru.neriva.app.ThemeManager.ThemeMode.SYSTEM) }, label = { Text("System") }, leadingIcon = { Icon(Icons.Default.PhoneAndroid, null, modifier = Modifier.size(18.dp)) })
            }

            HorizontalDivider()

            // Interface language
            Text("Interface Language", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
            ExposedDropdownMenuBox(expanded = interfaceLangExpanded, onExpandedChange = { interfaceLangExpanded = it }) {
                OutlinedTextField(
                    value = ru.neriva.app.LanguageManager.getLanguageName(interfaceLang),
                    onValueChange = {},
                    readOnly = true,
                    trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(interfaceLangExpanded) },
                    modifier = Modifier.fillMaxWidth().menuAnchor()
                )
                ExposedDropdownMenu(expanded = interfaceLangExpanded, onDismissRequest = { interfaceLangExpanded = false }) {
                    interfaceLocales.forEach { code ->
                        DropdownMenuItem(
                            text = { Text(ru.neriva.app.LanguageManager.getLanguageName(code)) },
                            onClick = {
                                ru.neriva.app.LanguageManager.setInterfaceLanguage(context, code)
                                interfaceLangExpanded = false
                                app.appScope.launch {
                                    try { app.sessionRepo.updateSettings(interfaceLanguage = code) } catch (_: Exception) {}
                                }
                                (context as? Activity)?.let {
                                    ru.neriva.app.LanguageManager.applyLocale(it, code)
                                    it.recreate()
                                }
                            }
                        )
                    }
                }
            }

            // Learning language
            Text("Learning Language", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
            ExposedDropdownMenuBox(expanded = learningLangExpanded, onExpandedChange = { learningLangExpanded = it }) {
                OutlinedTextField(
                    value = ru.neriva.app.LanguageManager.getLanguageName(learningLang),
                    onValueChange = {},
                    readOnly = true,
                    trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(learningLangExpanded) },
                    modifier = Modifier.fillMaxWidth().menuAnchor()
                )
                ExposedDropdownMenu(expanded = learningLangExpanded, onDismissRequest = { learningLangExpanded = false }) {
                    learningLocales.forEach { code ->
                        DropdownMenuItem(
                            text = { Text(ru.neriva.app.LanguageManager.getLanguageName(code)) },
                            onClick = {
                                ru.neriva.app.LanguageManager.setLearningLanguage(context, code)
                                learningLangExpanded = false
                                scope.launch {
                                    try { app.sessionRepo.updateSettings(learningLanguage = code); vm.refresh() } catch (_: Exception) {}
                                }
                            }
                        )
                    }
                }
            }

            HorizontalDivider()

            // Level selector
            var levelExpanded by remember { mutableStateOf(false) }
            val levels = listOf("A1", "A2", "B1", "B2", "C1", "C2")
            Text("Level", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
            ExposedDropdownMenuBox(expanded = levelExpanded, onExpandedChange = { levelExpanded = it }) {
                OutlinedTextField(
                    value = user?.level ?: "A1",
                    onValueChange = {},
                    readOnly = true,
                    trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(levelExpanded) },
                    modifier = Modifier.fillMaxWidth().menuAnchor()
                )
                ExposedDropdownMenu(expanded = levelExpanded, onDismissRequest = { levelExpanded = false }) {
                    levels.forEach { lv ->
                        DropdownMenuItem(
                            text = { Text(lv) },
                            onClick = {
                                levelExpanded = false
                                scope.launch {
                                    try { app.sessionRepo.updateSettings(level = lv); vm.refresh() } catch (_: Exception) {}
                                }
                            }
                        )
                    }
                }
            }

            // Learning focus
            var focusText by remember { mutableStateOf("") }
            Text("Learning Focus", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
            OutlinedTextField(
                value = focusText, onValueChange = { focusText = it },
                label = { Text("e.g. travel, business, exam prep") },
                modifier = Modifier.fillMaxWidth()
            )
            Button(onClick = {
                scope.launch {
                    try { app.sessionRepo.updateSettings(learningFocus = focusText.ifBlank { null }); vm.refresh() } catch (_: Exception) {}
                }
            }, modifier = Modifier.fillMaxWidth()) { Text("Save Focus") }

            HorizontalDivider()

            // Activation key
            var activationKey by remember { mutableStateOf("") }
            var activationResult by remember { mutableStateOf<String?>(null) }
            Text("Activation Key", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
            OutlinedTextField(
                value = activationKey, onValueChange = { activationKey = it },
                label = { Text("XXXX-XXXX-XXXX-XXXX") },
                modifier = Modifier.fillMaxWidth()
            )
            Button(onClick = {
                scope.launch {
                    try { app.sessionRepo.activateKey(activationKey); activationResult = "Key activated!"; vm.refresh() }
                    catch (e: Exception) { activationResult = e.message }
                }
            }, enabled = activationKey.isNotBlank(), modifier = Modifier.fillMaxWidth()) { Text("Activate") }
            activationResult?.let { Text(it, color = MaterialTheme.colorScheme.primary) }

            HorizontalDivider()

            // Change password
            var oldPass by remember { mutableStateOf("") }
            var newPass by remember { mutableStateOf("") }
            var confirmPass by remember { mutableStateOf("") }
            var passResult by remember { mutableStateOf<String?>(null) }
            Text("Change Password", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
            OutlinedTextField(value = oldPass, onValueChange = { oldPass = it }, label = { Text("Current password") }, modifier = Modifier.fillMaxWidth())
            OutlinedTextField(value = newPass, onValueChange = { newPass = it }, label = { Text("New password") }, modifier = Modifier.fillMaxWidth())
            OutlinedTextField(value = confirmPass, onValueChange = { confirmPass = it }, label = { Text("Confirm new password") }, modifier = Modifier.fillMaxWidth())
            Button(onClick = {
                scope.launch {
                    try { app.authRepo.changePassword(oldPass, newPass, confirmPass); passResult = "Password changed!" }
                    catch (e: Exception) { passResult = e.message }
                }
            }, enabled = oldPass.isNotBlank() && newPass.isNotBlank() && newPass == confirmPass, modifier = Modifier.fillMaxWidth()) { Text("Change Password") }
            passResult?.let { Text(it, color = MaterialTheme.colorScheme.primary) }

            HorizontalDivider()
            if (user?.premium == true) Text("Premium active until: ${user.premiumUntil ?: "-"}")
            Text("Referral code: ${user?.referralCode ?: "-"}")

            HorizontalDivider()

            // Telegram linking
            val telegramAccount = user?.telegramAccount
            val telegramLinked = telegramAccount != null || user?.telegramLinked == true
            val bot = state.session?.telegram_login_bot ?: "NERIVAapp_bot"
            val telegramIntent = remember(bot) {
                android.content.Intent(android.content.Intent.ACTION_VIEW, android.net.Uri.parse("https://t.me/$bot"))
            }
            Text("Telegram", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
            if (telegramLinked) {
                telegramAccount?.name?.let { Text("Linked: $it", style = MaterialTheme.typography.bodyMedium) }
                telegramAccount?.id?.let { Text("Telegram ID: $it", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant) }
            } else {
                Text(
                    "Link Telegram to share progress between bot and app.",
                    style = MaterialTheme.typography.bodySmall,
                    color = MaterialTheme.colorScheme.onSurfaceVariant,
                )
            }
            Row(horizontalArrangement = Arrangement.spacedBy(8.dp), modifier = Modifier.fillMaxWidth()) {
                OutlinedButton(
                    onClick = { context.startActivity(telegramIntent) },
                    modifier = Modifier.weight(1f),
                ) {
                    Icon(Icons.Default.Send, contentDescription = null, modifier = Modifier.size(18.dp))
                    Spacer(Modifier.width(6.dp))
                    Text("Open")
                }
                if (!telegramLinked) {
                    Button(
                        onClick = {
                            // Opens the Telegram bot to begin the two-factor linking flow
                            context.startActivity(telegramIntent)
                        },
                        modifier = Modifier.weight(1f),
                    ) { Text("Send code") }
                }
            }

            HorizontalDivider()

            // Social channels
            Text("Social channels", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
            Text("Follow NERIVA — short lessons, updates, and product tips.", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
            SocialLinkButton("Telegram", "https://t.me/NERIVAapp_bot", Icons.Default.Send, context)
            SocialLinkButton("YouTube", "https://www.youtube.com/@neriva_app", Icons.Default.PlayCircle, context)
            SocialLinkButton("Instagram", "https://www.instagram.com/neriva.ru", Icons.Default.CameraAlt, context)
            SocialLinkButton("TikTok", "https://tiktok.com/@nerivaru", Icons.Default.MusicNote, context)

            Spacer(Modifier.weight(1f))
            Button(
                onClick = { scope.launch { app.authRepo.logout(); vm.logout() } },
                colors = ButtonDefaults.buttonColors(containerColor = MaterialTheme.colorScheme.error),
                modifier = Modifier.fillMaxWidth()
            ) { Text("Log Out") }
        }
    }
}

// ---- Referral ----

@Composable
fun ReferralScreen(navController: androidx.navigation.NavHostController) {
    val vm: ru.neriva.app.ui.SessionViewModel = androidx.lifecycle.viewmodel.compose.viewModel()
    val state by vm.state.collectAsState()
    val user = state.session?.user

    ScreenScaffold("Referral", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            Text("Your code", style = MaterialTheme.typography.labelLarge)
            Text(user?.referralCode ?: "-", style = MaterialTheme.typography.headlineMedium, fontWeight = FontWeight.Bold)
            Text("Invited: ${user?.referralCount ?: 0}")
            Text("Balance: ${user?.referralBalance ?: "0"}")
            user?.referralInvitees?.takeIf { it.isNotEmpty() }?.let {
                Text("Invitees", style = MaterialTheme.typography.titleMedium)
                it.forEach { inv -> Text("• ${inv.name} — ${inv.xp} XP") }
            }
        }
    }
}

// ---- Bug Report ----

@Composable
fun BugReportScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    val scope = rememberCoroutineScope()
    var desc by remember { mutableStateOf("") }
    var sent by remember { mutableStateOf(false) }
    var loading by remember { mutableStateOf(false) }

    ScreenScaffold("Bug Report", navController) { padding ->
        Column(Modifier.padding(padding).fillMaxSize().padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            if (sent) {
                Text("Report sent. Thank you!", color = MaterialTheme.colorScheme.secondary, style = MaterialTheme.typography.titleMedium)
            }
            OutlinedTextField(value = desc, onValueChange = { desc = it }, label = { Text("Describe the issue") }, modifier = Modifier.fillMaxWidth(), minLines = 4)
            Button(onClick = {
                scope.launch {
                    loading = true
                    try {
                        val body = desc.toRequestBody("text/plain".toMediaType())
                        app.apiClient.bugReport(body, null)
                        sent = true; desc = ""
                    } catch (_: Exception) {} finally { loading = false }
                }
            }, enabled = !loading && desc.isNotBlank(), modifier = Modifier.fillMaxWidth()) {
                if (loading) CircularProgressIndicator(Modifier.size(20.dp), strokeWidth = 2.dp) else Text("Send Report")
            }
        }
    }
}

/** Copies a content:// Uri into a temporary file so it can be uploaded as multipart. */
private fun copyUriToTempFile(context: android.content.Context, uri: android.net.Uri, prefix: String): java.io.File {
    val suffix = when (context.contentResolver.getType(uri)) {
        "image/png" -> ".png"
        "image/webp" -> ".webp"
        else -> ".jpg"
    }
    val tmp = java.io.File.createTempFile(prefix, suffix, context.cacheDir)
    context.contentResolver.openInputStream(uri)?.use { input ->
        tmp.outputStream().use { output -> input.copyTo(output) }
    } ?: throw IllegalStateException("Cannot read selected image")
    return tmp
}

/** Outlined button that opens an external social link. */
@Composable
private fun SocialLinkButton(
    name: String,
    url: String,
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    context: android.content.Context,
) {
    OutlinedButton(
        onClick = {
            context.startActivity(android.content.Intent(android.content.Intent.ACTION_VIEW, android.net.Uri.parse(url)))
        },
        modifier = Modifier.fillMaxWidth(),
    ) {
        Icon(icon, contentDescription = name, modifier = Modifier.size(18.dp))
        Spacer(Modifier.width(8.dp))
        Text(name)
    }
}
