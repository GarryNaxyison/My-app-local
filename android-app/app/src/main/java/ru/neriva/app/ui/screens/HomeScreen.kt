package ru.neriva.app.ui.screens

import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.outlined.MenuBook
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavHostController
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.launch
import ru.neriva.app.NERIVAApp
import ru.neriva.app.R
import ru.neriva.app.ui.SessionViewModel
import ru.neriva.app.ui.components.NerivaAction
import ru.neriva.app.ui.components.nerivaWingCover

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(navController: NavHostController) {
    val vm: SessionViewModel = viewModel()
    val state by vm.state.collectAsState()
    val user = state.session?.user
    val colors = MaterialTheme.colorScheme
    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Row(horizontalArrangement = Arrangement.spacedBy(10.dp), verticalAlignment = Alignment.CenterVertically) {
                        Image(painterResource(R.drawable.neriva_brand_symbol), null, Modifier.width(36.dp))
                        Text(stringResource(R.string.today), style = MaterialTheme.typography.headlineMedium)
                    }
                },
                actions = {
                    user?.let { LearningTopicAction(it.learningFocus.orEmpty(), vm::updateLearningFocus) }
                    IconButton(onClick = { navController.navigate("bug-report") }) {
                        Icon(Icons.Outlined.BugReport, stringResource(R.string.send_report))
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(containerColor = colors.background),
            )
        },
    ) { padding ->
        Box(
            Modifier.fillMaxSize().padding(padding).verticalScroll(rememberScrollState()),
            contentAlignment = Alignment.TopCenter,
        ) {
            Column(
                Modifier.widthIn(max = 860.dp).fillMaxWidth().padding(horizontal = 18.dp).padding(top = 8.dp, bottom = 24.dp),
                verticalArrangement = Arrangement.spacedBy(20.dp),
            ) {
                TutorCover { navController.navigate("tutor") }
                LearningShortcuts { navController.navigate(it) }
                Surface(shape = MaterialTheme.shapes.large, color = colors.surface) {
                    Column(Modifier.fillMaxWidth().padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
                        Text(stringResource(R.string.progress), style = MaterialTheme.typography.titleLarge)
                        Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(10.dp)) {
                            StatItem("XP", user?.xp ?: 0, Modifier.weight(1f))
                            StatItem(stringResource(R.string.level_label), user?.xpLevel ?: 0, Modifier.weight(1f))
                            StatItem(stringResource(R.string.words), user?.learnedWords ?: 0, Modifier.weight(1f))
                        }
                        LearningProgress(user?.xpCurrent ?: 0, user?.xpNeeded ?: 0)
                        HorizontalDivider(color = colors.outline)
                        Text(stringResource(R.string.learning_activity), style = MaterialTheme.typography.titleMedium)
                        LimitRow(stringResource(R.string.lessons_today), user?.lessonsToday ?: 0, user?.lessonLimit ?: 0)
                        LimitRow(stringResource(R.string.practice_today), user?.practiceToday ?: 0, user?.practiceLimit ?: 0)
                        if ((user?.voiceLimit ?: 0) > 0) {
                            LimitRow(stringResource(R.string.voices_today), user?.voiceToday ?: 0, user?.voiceLimit ?: 0)
                        }
                    }
                }
                DailyBonus(onClaimed = { vm.refresh() })
            }
        }
    }
}

@Composable
private fun TutorCover(onStart: () -> Unit) {
    val shape = RoundedCornerShape(28.dp)
    Box(
        Modifier.fillMaxWidth().clip(shape)
            .nerivaWingCover(),
    ) {
        Column(Modifier.padding(22.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            Image(painterResource(R.drawable.neriva_brand_symbol), null, Modifier.width(58.dp))
            Text(stringResource(R.string.ai_tutor), color = Color.White, style = MaterialTheme.typography.displaySmall)
            Text(
                stringResource(R.string.home_tutor_summary),
                modifier = Modifier.widthIn(max = 294.dp),
                color = Color(0xFFE8EBFF), style = MaterialTheme.typography.bodyMedium,
            )
            NerivaAction(
                label = stringResource(R.string.home_start_lesson),
                onClick = onStart,
                onCover = true,
            )
        }
    }
}

private data class Shortcut(val route: String, val label: Int, val icon: ImageVector)

private val homeShortcuts = listOf(
    Shortcut("practice", R.string.practice, Icons.Outlined.Forum),
    Shortcut("words", R.string.learn_words, Icons.AutoMirrored.Outlined.MenuBook),
    Shortcut("games", R.string.learning_games, Icons.Outlined.Extension),
    Shortcut("roleplay-scenarios", R.string.roleplay, Icons.Outlined.TheaterComedy),
    Shortcut("pronunciation", R.string.pronunciation, Icons.Outlined.RecordVoiceOver),
    Shortcut("tools", R.string.tools, Icons.Outlined.Tune),
)

@Composable
private fun LearningShortcuts(onOpen: (String) -> Unit) {
    val colors = MaterialTheme.colorScheme
    Column(verticalArrangement = Arrangement.spacedBy(10.dp)) {
        Text(stringResource(R.string.home_learning_shortcuts), style = MaterialTheme.typography.headlineSmall)
        BoxWithConstraints(Modifier.fillMaxWidth()) {
            val columns = if (maxWidth >= 600.dp) 3 else 2
            Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                homeShortcuts.chunked(columns).forEach { tiles ->
                    Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                        tiles.forEach { tile ->
                            Surface(
                                onClick = { onOpen(tile.route) },
                                modifier = Modifier.weight(1f).heightIn(min = 64.dp),
                                shape = MaterialTheme.shapes.medium,
                                color = colors.surface,
                            ) {
                                Row(Modifier.padding(12.dp), verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(10.dp)) {
                                    Icon(tile.icon, null, Modifier.size(20.dp), tint = colors.primary)
                                    Text(stringResource(tile.label), style = MaterialTheme.typography.titleSmall)
                                }
                            }
                        }
                        repeat(columns - tiles.size) { Spacer(Modifier.weight(1f)) }
                    }
                }
            }
        }
    }
}

@Composable
private fun StatItem(label: String, value: Number, modifier: Modifier) {
    Column(modifier, verticalArrangement = Arrangement.spacedBy(2.dp)) {
        Text(value.toString(), style = MaterialTheme.typography.headlineSmall)
        Text(label, style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
    }
}

@Composable
private fun LearningProgress(used: Int, limit: Int) {
    val fraction by animateFloatAsState(
        targetValue = if (limit > 0) (used.toFloat() / limit).coerceIn(0f, 1f) else 0f,
        label = "learningProgress",
    )
    LinearProgressIndicator(
        progress = { fraction }, modifier = Modifier.fillMaxWidth().height(5.dp).clip(RoundedCornerShape(4.dp)),
        trackColor = MaterialTheme.colorScheme.surfaceVariant,
    )
}

@Composable
private fun LimitRow(label: String, used: Int, limit: Int) {
    Column(verticalArrangement = Arrangement.spacedBy(5.dp)) {
        Row(Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(10.dp)) {
            Text(label, Modifier.weight(1f), style = MaterialTheme.typography.bodyMedium)
            Text("$used / $limit", style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant)
        }
        LearningProgress(used, limit)
    }
}

@Composable
private fun DailyBonus(onClaimed: () -> Unit) {
    var claiming by remember { mutableStateOf(false) }
    var result by remember { mutableStateOf<String?>(null) }
    val scope = rememberCoroutineScope()
    val claimed = stringResource(R.string.bonus_claimed)
    Surface(shape = MaterialTheme.shapes.large, color = MaterialTheme.colorScheme.surface) {
        Column(Modifier.fillMaxWidth().padding(16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
            Row(Modifier.fillMaxWidth(), verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                Icon(Icons.Outlined.Star, null, tint = MaterialTheme.colorScheme.primary)
                Text(stringResource(R.string.claim_daily_bonus), Modifier.weight(1f), style = MaterialTheme.typography.titleMedium)
                FilledTonalButton(
                    onClick = {
                        claiming = true
                        scope.launch {
                            try {
                                val res = NERIVAApp.instance.sessionRepo.claimDailyBonus()
                                result = res.xp?.let { "+$it XP" } ?: claimed
                                onClaimed()
                            } catch (e: CancellationException) {
                                throw e
                            } catch (e: Exception) {
                                result = e.message
                            } finally {
                                claiming = false
                            }
                        }
                    },
                    enabled = !claiming,
                    modifier = Modifier.heightIn(min = 48.dp),
                ) {
                    AnimatedContent(claiming, label = "bonusStatus") { busy ->
                        if (busy) CircularProgressIndicator(Modifier.size(18.dp), strokeWidth = 2.dp)
                        else Text(stringResource(R.string.home_claim_bonus))
                    }
                }
            }
            result?.let { Text(it, style = MaterialTheme.typography.bodySmall) }
        }
    }
}
