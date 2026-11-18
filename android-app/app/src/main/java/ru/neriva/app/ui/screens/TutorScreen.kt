package ru.neriva.app.ui.screens

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.outlined.Bookmark
import androidx.compose.material.icons.outlined.Check
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.semantics.LiveRegionMode
import androidx.compose.ui.semantics.heading
import androidx.compose.ui.semantics.liveRegion
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.dp
import androidx.navigation.NavHostController
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.launch
import ru.neriva.app.LanguageManager
import ru.neriva.app.NERIVAApp
import ru.neriva.app.R
import ru.neriva.app.data.model.AiTutorResponse
import ru.neriva.app.data.model.LessonData
import ru.neriva.app.data.model.WordData
import ru.neriva.app.data.repo.TutorRepository
import ru.neriva.app.ui.components.NerivaAction

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TutorScreen(navController: NavHostController) {
    val app = NERIVAApp.instance
    val repo = remember { TutorRepository(app.apiClient) }
    val scope = rememberCoroutineScope()
    val context = LocalContext.current
    val snackbar = remember { SnackbarHostState() }
    val phrases = remember(app) { app.phrasebookRepo.observeAll() }
    val savedPhrases by phrases.collectAsState(initial = emptyList())
    val saved = remember(savedPhrases) { savedPhrases.map { it.phrase.lowercase() }.toSet() }
    var savingPhrase by remember { mutableStateOf(false) }
    var loading by remember { mutableStateOf(true) }
    var response by remember { mutableStateOf<AiTutorResponse?>(null) }
    var sessionId by remember { mutableStateOf<String?>(null) }
    var lesson by remember { mutableStateOf<LessonData?>(null) }
    var failed by remember { mutableStateOf(false) }
    var answer by rememberSaveable { mutableStateOf("") }

    suspend fun request(input: String? = null) {
        loading = true
        failed = false
        try {
            val result = if (input == null) repo.start() else repo.answer(input, sessionId)
            sessionId = result.session?.id ?: sessionId
            response = result
            lesson = result.lesson ?: result.nextStep?.lesson ?: lesson
            if (input != null) answer = ""
        } catch (e: CancellationException) {
            throw e
        } catch (_: Exception) {
            failed = true
        } finally {
            loading = false
        }
    }

    LaunchedEffect(Unit) { request() }

    fun submit(input: String) {
        if (loading || input.isBlank()) return
        scope.launch { request(input.trim()) }
    }

    fun saveWord(word: WordData) {
        if (savingPhrase) return
        savingPhrase = true
        scope.launch {
            val target = word.target?.trim().orEmpty()
            val translation = word.interfaceTranslation?.trim().orEmpty()
            val example = word.exampleSentenceTarget?.trim().orEmpty()
            val message = try {
                app.phrasebookRepo.save(
                    phrase = example.ifEmpty { target },
                    translation = if (example.isEmpty()) translation else word.exampleTranslationInterface,
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
    }

    Scaffold(
        snackbarHost = { SnackbarHost(snackbar) },
        topBar = {
            TopAppBar(
                title = { Text(stringResource(R.string.ai_tutor), style = MaterialTheme.typography.titleLarge) },
                navigationIcon = {
                    IconButton(onClick = { navController.popBackStack() }) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, stringResource(R.string.back))
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.background),
            )
        },
    ) { padding ->
        Box(
            Modifier.fillMaxSize().padding(padding).verticalScroll(rememberScrollState()),
            contentAlignment = Alignment.TopCenter,
        ) {
            Column(
                Modifier.widthIn(max = 1080.dp).fillMaxWidth().padding(start = 18.dp, end = 18.dp, bottom = 28.dp),
                verticalArrangement = Arrangement.spacedBy(16.dp),
            ) {
                val current = response
                if (current == null) {
                    TutorPanel {
                        if (loading) {
                            CircularProgressIndicator(Modifier.size(28.dp), strokeWidth = 3.dp)
                            Text(stringResource(R.string.tutor_loading), style = MaterialTheme.typography.titleMedium)
                        } else if (failed) {
                            Text(stringResource(R.string.tutor_start_failed), color = MaterialTheme.colorScheme.error)
                            NerivaAction(stringResource(R.string.retry), onClick = { scope.launch { request() } })
                        }
                    }
                } else {
                    val material = lesson
                    material?.let { TutorLessonCover(it) }
                    BoxWithConstraints(Modifier.fillMaxWidth()) {
                        val wide = maxWidth >= 680.dp && material != null
                        val reading: @Composable () -> Unit = {
                            if (material != null) TutorReading(material, saved, savingPhrase, ::saveWord)
                        }
                        val exercise: @Composable () -> Unit = {
                            TutorExercise(current, answer, { answer = it }, loading, failed, ::submit)
                        }
                        if (wide) {
                            Row(horizontalArrangement = Arrangement.spacedBy(16.dp)) {
                                Column(Modifier.weight(1.15f), verticalArrangement = Arrangement.spacedBy(12.dp)) { reading() }
                                Column(Modifier.weight(1f)) { exercise() }
                            }
                        } else {
                            Column(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                                reading()
                                exercise()
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun TutorLessonCover(lesson: LessonData) {
    Surface(shape = MaterialTheme.shapes.extraLarge) {
        Column(
            Modifier.fillMaxWidth()
                .background(Brush.horizontalGradient(listOf(Color(0xFF191F58), Color(0xFF3037A4))))
                .padding(20.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp),
        ) {
            Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(10.dp)) {
                Image(painterResource(R.drawable.neriva_brand_symbol), null, Modifier.width(37.dp))
                val metadata = listOfNotNull(lesson.theme, lesson.level?.takeUnless { lesson.theme?.contains(it, ignoreCase = true) == true })
                    .filter { it.isNotBlank() }.joinToString(" · ")
                if (metadata.isNotBlank()) Text(metadata, color = Color(0xFFE8EBFF), style = MaterialTheme.typography.labelLarge)
            }
            Text(
                lesson.title?.takeIf { it.isNotBlank() } ?: stringResource(R.string.ai_tutor),
                modifier = Modifier.semantics { heading() },
                color = Color.White,
                style = MaterialTheme.typography.displaySmall,
            )
            lesson.lessonGoal?.takeIf { it.isNotBlank() }?.let {
                Text(it, color = Color(0xFFE8EBFF), style = MaterialTheme.typography.bodyMedium)
            }
        }
    }
}

@Composable
private fun TutorReading(lesson: LessonData, saved: Set<String>, saving: Boolean, onSave: (WordData) -> Unit) {
    val colors = MaterialTheme.colorScheme
    lesson.story?.textTarget?.takeIf { it.isNotBlank() }?.let { story ->
        TutorPanel {
            TutorHeading(stringResource(R.string.tutor_step_story))
            Text(story, style = MaterialTheme.typography.bodyLarge)
        }
    }
    if (lesson.words.any { !it.target.isNullOrBlank() }) {
        TutorPanel {
            TutorHeading(stringResource(R.string.tutor_words))
            lesson.words.forEach { word ->
                val target = word.target?.trim().orEmpty()
                if (target.isEmpty()) return@forEach
                val example = word.exampleSentenceTarget?.trim().orEmpty()
                val translation = word.interfaceTranslation?.trim().orEmpty()
                val isSaved = example.ifEmpty { target }.lowercase() in saved
                Surface(
                    onClick = { onSave(word) },
                    enabled = !saving && !isSaved,
                    shape = MaterialTheme.shapes.medium,
                    color = colors.surfaceVariant,
                ) {
                    Row(
                        Modifier.fillMaxWidth().padding(12.dp),
                        horizontalArrangement = Arrangement.spacedBy(10.dp),
                        verticalAlignment = Alignment.CenterVertically,
                    ) {
                        Column(Modifier.weight(1f), verticalArrangement = Arrangement.spacedBy(3.dp)) {
                            Text(target, style = MaterialTheme.typography.titleMedium)
                            if (translation.isNotEmpty()) Text(translation, style = MaterialTheme.typography.bodyMedium)
                            if (example.isNotEmpty()) Text(example, style = MaterialTheme.typography.bodySmall, color = colors.onSurfaceVariant)
                        }
                        Icon(
                            if (isSaved) Icons.Outlined.Check else Icons.Outlined.Bookmark,
                            stringResource(if (isSaved) R.string.tutor_note_saved else R.string.save_to_phrasebook),
                            tint = colors.primary,
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun TutorExercise(
    response: AiTutorResponse,
    answer: String,
    onAnswerChange: (String) -> Unit,
    loading: Boolean,
    failed: Boolean,
    onSubmit: (String) -> Unit,
) {
    val step = response.nextStep
    val colors = MaterialTheme.colorScheme
    TutorPanel {
        TutorHeading(step?.title?.takeIf { it.isNotBlank() } ?: stringResource(R.string.tutor_user_turn))
        val feedback = response.feedback
        if (feedback != null && (feedback.ok == true || !feedback.message.isNullOrBlank())) {
            Surface(
                shape = MaterialTheme.shapes.medium,
                color = if (feedback.ok == false) colors.errorContainer else colors.secondaryContainer,
                contentColor = if (feedback.ok == false) colors.onErrorContainer else colors.onSecondaryContainer,
                modifier = Modifier.fillMaxWidth().semantics { liveRegion = LiveRegionMode.Polite },
            ) {
                Column(Modifier.padding(12.dp), verticalArrangement = Arrangement.spacedBy(4.dp)) {
                    if (feedback.ok == true) Text(stringResource(R.string.correct), style = MaterialTheme.typography.titleMedium)
                    feedback.message?.takeIf { it.isNotBlank() }?.let { Text(it) }
                }
            }
        }
        step?.instruction?.takeIf { it.isNotBlank() }?.let { Text(it, color = colors.onSurfaceVariant) }
        step?.question?.questionTarget?.takeIf { it.isNotBlank() }?.let { Text(it, style = MaterialTheme.typography.headlineSmall) }
        step?.options.orEmpty().filter { !it.text.isNullOrBlank() }.forEachIndexed { index, option ->
            val value = option.text.orEmpty()
            NerivaAction(
                label = value,
                onClick = { onSubmit(value) },
                enabled = !loading,
                modifier = Modifier.fillMaxWidth(),
                optionNumber = index + 1,
            )
        }
        OutlinedTextField(
            value = answer,
            onValueChange = onAnswerChange,
            label = { Text(stringResource(R.string.lesson_answer)) },
            modifier = Modifier.fillMaxWidth(),
            enabled = !loading,
            maxLines = 5,
        )
        NerivaAction(
            label = stringResource(R.string.send),
            onClick = { onSubmit(answer) },
            enabled = !loading && answer.isNotBlank(),
            modifier = Modifier.fillMaxWidth(),
        )
        AnimatedVisibility(loading) { LinearProgressIndicator(Modifier.fillMaxWidth()) }
        if (failed) {
            Text(stringResource(R.string.try_again), color = colors.error, modifier = Modifier.semantics { liveRegion = LiveRegionMode.Polite })
        }
    }
}

@Composable
private fun TutorPanel(content: @Composable ColumnScope.() -> Unit) {
    Surface(Modifier.fillMaxWidth(), shape = MaterialTheme.shapes.large, color = MaterialTheme.colorScheme.surface) {
        Column(Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp), content = content)
    }
}

@Composable
private fun TutorHeading(label: String) {
    Text(label, Modifier.semantics { heading() }, style = MaterialTheme.typography.headlineSmall)
}
