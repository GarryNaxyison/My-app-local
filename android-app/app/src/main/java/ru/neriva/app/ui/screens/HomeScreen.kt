package ru.neriva.app.ui.screens

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
                title = { Text("Today", fontWeight = FontWeight.SemiBold) },
                actions = {
                    IconButton(onClick = { navController.navigate("bug-report") }) {
                        Icon(Icons.Default.BugReport, contentDescription = "Report")
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
                        Text("Your progress", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold, color = MaterialTheme.colorScheme.onSurface)
                        Row(horizontalArrangement = Arrangement.spacedBy(24.dp)) {
                            StatItem("XP", "${user?.xp ?: 0}")
                            StatItem("Level", "${user?.xpLevel ?: 0}")
                            StatItem("Words", "${user?.learnedWords ?: 0}")
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
                        Text("Daily activity", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold, color = MaterialTheme.colorScheme.onSurface)
                        LimitRow("Lessons", user?.lessonsToday ?: 0, user?.lessonLimit ?: 0)
                        LimitRow("Practice", user?.practiceToday ?: 0, user?.practiceLimit ?: 0)
                        if ((user?.voiceLimit ?: 0) > 0) {
                            LimitRow("Voice", user?.voiceToday ?: 0, user?.voiceLimit ?: 0)
                        }
                    }
                }

                // Daily bonus
                var claimingBonus by remember { mutableStateOf(false) }
                val app = ru.neriva.app.NERIVAApp.instance
                val scope = androidx.compose.runtime.rememberCoroutineScope()
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
                            Text("Daily Bonus", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold, color = MaterialTheme.colorScheme.onSurface)
                            Text(
                                bonusResult ?: "Claim your daily XP bonus",
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
                                                bonusResult = if (res.xp != null) "+${res.xp} XP!" else "Claimed!"
                                                vm.refresh()
                                            } catch (e: Exception) {
                                                bonusResult = e.message
                                            } finally { claimingBonus = false }
                                        }
                                    }
                                },
                                text = "Claim",
                            )
                        }
                    }
                }

                // Quick actions
                Text(
                    "Learning Lab",
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
        LabTile("tutor", "AI Tutor", "tutor"),
        LabTile("practice", "Practice", "practice"),
        LabTile("words", "Learn Words", "words"),
        LabTile("roleplay", "Roleplay", "roleplay"),
        LabTile("pronunciation", "Pronunciation", "pronunciation"),
        LabTile("tools", "Tools", "tools"),
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
