package ru.neriva.app.ui.screens

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import ru.neriva.app.NERIVAApp
import ru.neriva.app.data.api.ProgressResponse

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DashboardScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    var progress by remember { mutableStateOf<ProgressResponse?>(null) }
    var loading by remember { mutableStateOf(true) }
    LaunchedEffect(Unit) { try { progress = app.sessionRepo.getProgress() } catch (_: Exception) {} finally { loading = false } }

    Scaffold(topBar = {
        TopAppBar(
            title = { Text(stringResource(R.string.dashboard)) },
            navigationIcon = { IconButton(onClick = { navController.popBackStack() }) { Icon(Icons.AutoMirrored.Filled.ArrowBack, ru.neriva.app.NERIVAApp.instance.getString(R.string.back)) } }
        )
    }) { padding ->
        if (loading) { Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) { CircularProgressIndicator() } }
        else Column(Modifier.padding(padding).fillMaxSize().padding(16.dp).verticalScroll(rememberScrollState()), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            progress?.let { p ->
                Text(stringResource(R.string.progress_without_leaderboard), style = MaterialTheme.typography.headlineSmall, fontWeight = FontWeight.Bold)
                HorizontalDivider()
                StatRow(ru.neriva.app.NERIVAApp.instance.getString(R.string.xp_gained_label), "${p.xp}")
                StatRow(ru.neriva.app.NERIVAApp.instance.getString(R.string.level_label), "${p.xpLevel} — ${p.xpTitle ?: ""}")
                HorizontalDivider()
                Text(stringResource(R.string.activity), style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                StatRow(ru.neriva.app.NERIVAApp.instance.getString(R.string.lessons_completed), "${p.lessonCount}")
                StatRow(ru.neriva.app.NERIVAApp.instance.getString(R.string.practices_total), "${p.practiceCount}")
                StatRow(ru.neriva.app.NERIVAApp.instance.getString(R.string.word_game), "${p.wordGameCount}")
                StatRow(ru.neriva.app.NERIVAApp.instance.getString(R.string.words_learned), "${p.learnedWords}")
                StatRow(ru.neriva.app.NERIVAApp.instance.getString(R.string.mistakes), "${p.mistakes}")
                p.streak?.let { StatRow(ru.neriva.app.NERIVAApp.instance.getString(R.string.habit_calendar), "$it") }
                HorizontalDivider()
                if (p.premium) {
                    Card(Modifier.fillMaxWidth(), colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.tertiaryContainer)) {
                        Text(stringResource(R.string.premium), Modifier.padding(16.dp), fontWeight = FontWeight.Bold)
                    }
                }
            }
        }
    }
}

@Composable
private fun StatRow(label: String, value: String) {
    Row(Modifier.fillMaxWidth().padding(vertical = 4.dp), horizontalArrangement = Arrangement.SpaceBetween) {
        Text(label, style = MaterialTheme.typography.bodyMedium)
        Text(value, style = MaterialTheme.typography.bodyMedium, fontWeight = FontWeight.SemiBold)
    }
}
