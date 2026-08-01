package ru.neriva.app.ui.screens

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.background
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.BugReport
import androidx.compose.material.icons.filled.Star
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.compose.material3.MaterialTheme
import kotlinx.coroutines.launch
import ru.neriva.app.ui.SessionViewModel
import ru.neriva.app.ui.components.AnimatedBackground
import ru.neriva.app.ui.components.GlassCard
import ru.neriva.app.ui.components.GradientButton
import ru.neriva.app.ui.components.XpProgressBar
import ru.neriva.app.ui.components.ServerIcon
import androidx.compose.foundation.isSystemInDarkTheme

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun HomeScreen(navController: androidx.navigation.NavHostController) {
    val vm: SessionViewModel = viewModel()
    val state by vm.state.collectAsState()
    val user = state.session?.user

    Scaffold(
        containerColor = MaterialTheme.colorScheme.background,
        topBar = {
            TopAppBar(
                title = { Text(stringResource(R.string.today), fontWeight = FontWeight.SemiBold) },
                actions = {
                    IconButton(onClick = { navController.navigate("bug-report") }) {
                        Icon(Icons.Default.BugReport, contentDescription = stringResource(R.string.send_report))
                    }
                }
            )
        }
    ) { padding ->
        AnimatedBackground(modifier = Modifier.padding(padding)) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .verticalScroll(rememberScrollState())
                    .padding(16.dp),
                verticalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                // Progress card
                GlassCard(modifier = Modifier.fillMaxWidth()) {
                    Column(modifier = Modifier.padding(20.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
                        Text(stringResource(R.string.progress), style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold, color = MaterialTheme.colorScheme.onSurface)
                        Row(horizontalArrangement = Arrangement.spacedBy(24.dp)) {
                            StatItem("XP", "${user?.xp ?: 0}")
                            StatItem(stringResource(R.string.level_label), "${user?.xpLevel ?: 0}")
                            StatItem(stringResource(R.string.words), "${user?.learnedWords ?: 0}")
                        }
                        XpProgressBar(
                            currentXp = user?.xpCurrent ?: 0,
                            maxXp = user?.xpNeeded ?: 0,
                            modifier = Modifier.fillMaxWidth(),
                        )
                    }
                }

                // Daily limits card
                GlassCard(modifier = Modifier.fillMaxWidth()) {
                    Column(modifier = Modifier.padding(20.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
                        Text(stringResource(R.string.learning_activity), style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold, color = MaterialTheme.colorScheme.onSurface)
                        LimitRow(stringResource(R.string.lessons_today), user?.lessonsToday ?: 0, user?.lessonLimit ?: 0)
                        LimitRow(stringResource(R.string.practice_today), user?.practiceToday ?: 0, user?.practiceLimit ?: 0)
                        if ((user?.voiceLimit ?: 0) > 0) {
                            LimitRow(stringResource(R.string.voices_today), user?.voiceToday ?: 0, user?.voiceLimit ?: 0)
                        }
                    }
                }

                // Daily bonus
                var claimingBonus by remember { mutableStateOf(false) }
                val app = ru.neriva.app.NERIVAApp.instance
                val scope = androidx.compose.runtime.rememberCoroutineScope()
                val claimDailyBonus = stringResource(R.string.claim_daily_bonus)
                val bonusClaimed = stringResource(R.string.bonus_claimed)
                var bonusResult by remember { mutableStateOf<String?>(null) }

                GlassCard(modifier = Modifier.fillMaxWidth()) {
                    Row(
                        modifier = Modifier
                            .padding(20.dp)
                            .clip(MaterialTheme.shapes.medium),
                        verticalAlignment = Alignment.CenterVertically,
                    ) {
                        Box(
                            modifier = Modifier
                                .size(40.dp)
                                .clip(CircleShape)
                                .background(MaterialTheme.colorScheme.tertiary.copy(alpha = 0.18f)),
                            contentAlignment = Alignment.Center,
                        ) {
                            Icon(Icons.Default.Star, contentDescription = null, tint = MaterialTheme.colorScheme.tertiary)
                        }
                        Spacer(Modifier.width(12.dp))
                        Column(modifier = Modifier.weight(1f)) {
                Text(stringResource(R.string.claim_daily_bonus), style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold, color = MaterialTheme.colorScheme.onSurface)
                            Text(
                                bonusResult ?: claimDailyBonus,
                                style = MaterialTheme.typography.bodySmall,
                                color = MaterialTheme.colorScheme.onSurfaceVariant,
                            )
                        }
                        if (claimingBonus) {
                            CircularProgressIndicator(modifier = Modifier.size(20.dp), strokeWidth = 2.dp)
                        } else {
                            GradientButton(
                                onClick = {
                                    if (!claimingBonus) {
                                        claimingBonus = true
                                        scope.launch {
                                            try {
                                                val res = app.sessionRepo.claimDailyBonus()
                                                bonusResult = if (res.xp != null) "+${res.xp} XP!" else bonusClaimed
                                                vm.refresh()
                                            } catch (e: Exception) {
                                                bonusResult = e.message
                                            } finally { claimingBonus = false }
                                        }
                                    }
                                },
                                text = claimDailyBonus,
                            )
                        }
                    }
                }

                // Quick actions
                Text(stringResource(R.string.v2_learning_lab),
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.SemiBold,
                    color = MaterialTheme.colorScheme.onBackground,
                )
                QuickActionGrid(navController)
            }
        }
    }
}

@Composable
fun StatItem(label: String, value: String) {
    Column(horizontalAlignment = Alignment.CenterHorizontally) {
        Text(value, style = MaterialTheme.typography.headlineSmall, fontWeight = FontWeight.Bold, color = MaterialTheme.colorScheme.onBackground)
        Text(label, style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
    }
}

@Composable
fun LimitRow(label: String, used: Int, limit: Int) {
    Column {
        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
            Text(label, style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurface)
            Text("$used / $limit", style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant)
        }
        LinearProgressIndicator(
            progress = { if (limit > 0) used.toFloat() / limit else 0f },
            modifier = Modifier.fillMaxWidth().height(6.dp),
        )
    }
}

private data class LabTile(val route: String, val label: String, val viewId: String)

@Composable
fun QuickActionGrid(navController: androidx.navigation.NavHostController) {
    val darkTheme = isSystemInDarkTheme()
    val tiles = listOf(
        LabTile("tutor", stringResource(R.string.ai_tutor), "tutor"),
        LabTile("practice", stringResource(R.string.practice), "practice"),
        LabTile("words", stringResource(R.string.learn_words), "words"),
        LabTile("roleplay", stringResource(R.string.roleplay), "roleplay"),
        LabTile("pronunciation", stringResource(R.string.pronunciation), "pronunciation"),
        LabTile("tools", stringResource(R.string.tools), "tools"),
    )
    Column(verticalArrangement = Arrangement.spacedBy(10.dp)) {
        tiles.chunked(2).forEach { row ->
            Row(horizontalArrangement = Arrangement.spacedBy(10.dp)) {
                row.forEach { tile ->
                    GlassCard(modifier = Modifier.weight(1f).height(92.dp)) {
                        androidx.compose.foundation.layout.Box(
                            modifier = Modifier
                                .fillMaxSize()
                                .padding(12.dp),
                            contentAlignment = Alignment.Center,
                        ) {
                            Column(
                                horizontalAlignment = Alignment.CenterHorizontally,
                                verticalArrangement = Arrangement.Center,
                            ) {
                                ServerIcon(
                                    viewId = tile.viewId,
                                    darkTheme = darkTheme,
                                    size = 30.dp,
                                    contentDescription = tile.label,
                                )
                                Spacer(Modifier.height(6.dp))
                                Text(
                                    tile.label,
                                    style = MaterialTheme.typography.titleSmall,
                                    color = MaterialTheme.colorScheme.onSurface,
                                )
                            }
                        }
                    }
                }
                if (row.size == 1) Spacer(Modifier.weight(1f))
            }
        }
    }
}
