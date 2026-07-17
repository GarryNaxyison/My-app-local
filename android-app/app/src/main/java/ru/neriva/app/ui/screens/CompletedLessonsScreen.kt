package ru.neriva.app.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import ru.neriva.app.NERIVAApp
import ru.neriva.app.data.api.CompletedLesson

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CompletedLessonsScreen(navController: androidx.navigation.NavHostController) {
    val app = NERIVAApp.instance
    var lessons by remember { mutableStateOf<List<CompletedLesson>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    LaunchedEffect(Unit) {
        try { lessons = app.apiClient.aiTutorCompleted().items ?: emptyList() }
        catch (_: Exception) {} finally { loading = false }
    }

    Scaffold(topBar = {
        TopAppBar(
            title = { Text("Completed Lessons") },
            navigationIcon = { IconButton(onClick = { navController.popBackStack() }) { Icon(Icons.AutoMirrored.Filled.ArrowBack, "Back") } }
        )
    }) { padding ->
        if (loading) {
            Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) { CircularProgressIndicator() }
        } else if (lessons.isEmpty()) {
            Box(Modifier.padding(padding).fillMaxSize(), Alignment.Center) {
                Text("No completed lessons yet.", color = MaterialTheme.colorScheme.onSurfaceVariant)
            }
        } else {
            LazyColumn(
                Modifier.padding(padding).fillMaxSize().padding(horizontal = 16.dp),
                verticalArrangement = Arrangement.spacedBy(8.dp),
                contentPadding = PaddingValues(vertical = 16.dp)
            ) {
                items(lessons) { lesson ->
                    Card(Modifier.fillMaxWidth()) {
                        Column(Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(4.dp)) {
                            Text(lesson.title ?: "Lesson", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                            lesson.topic?.let { Text(it, style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant) }
                            Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                                lesson.level?.let { Text("Level: $it", style = MaterialTheme.typography.bodySmall) }
                                lesson.completedAt?.let { Text(it.take(10), style = MaterialTheme.typography.bodySmall) }
                            }
                        }
                    }
                }
            }
        }
    }
}
