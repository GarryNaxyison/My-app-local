package ru.neriva.app.ui

import androidx.compose.foundation.clickable
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import ru.neriva.app.ui.nav.Screen
import androidx.compose.material3.Icon
import coil.compose.AsyncImage
import coil.compose.AsyncImagePainter
import coil.compose.SubcomposeAsyncImage
import coil.request.ImageRequest
import androidx.compose.runtime.getValue

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MoreSheet(
    items: List<Screen>,
    onDismiss: () -> Unit,
    onSelect: (Screen) -> Unit
) {
    val darkTheme = isSystemInDarkTheme()
    ModalBottomSheet(onDismissRequest = onDismiss) {
        Column(modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text("More", style = MaterialTheme.typography.headlineSmall, fontWeight = FontWeight.SemiBold)
                IconButton(onClick = onDismiss) {
                    Icon(Icons.Default.Close, contentDescription = "Close")
                }
            }
            Spacer(Modifier.height(12.dp))
            LazyVerticalGrid(
                columns = GridCells.Fixed(3),
                contentPadding = PaddingValues(bottom = 24.dp),
                horizontalArrangement = Arrangement.spacedBy(10.dp),
                verticalArrangement = Arrangement.spacedBy(10.dp),
                modifier = Modifier.heightIn(max = 480.dp)
            ) {
                items(items) { screen ->
                    MoreTile(
                        title = screen.title,
                        viewId = screen.viewId,
                        fallback = iconForScreen(screen.viewId),
                        darkTheme = darkTheme,
                    ) { onSelect(screen) }
                }
            }
        }
    }
}

/** Material fallback icons keyed by web app viewId. */
private fun iconForScreen(viewId: String): ImageVector = when (viewId) {
    "roleplay" -> Icons.Default.TheaterComedy
    "pronunciation" -> Icons.Default.Mic
    "shadowing" -> Icons.Default.RecordVoiceOver
    "spelling" -> Icons.Default.Spellcheck
    "word-game" -> Icons.Default.Games
    "vocabulary" -> Icons.Default.MenuBook
    "phrasebook" -> Icons.Default.Bookmarks
    "offline" -> Icons.Default.WifiOff
    "level" -> Icons.Default.Assessment
    "progress" -> Icons.Default.BarChart
    "awards" -> Icons.Default.EmojiEvents
    "leaderboard" -> Icons.Default.Leaderboard
    "limits" -> Icons.Default.Speed
    "mistakes" -> Icons.Default.ReportProblem
    "tools" -> Icons.Default.Build
    "premium" -> Icons.Default.WorkspacePremium
    "settings" -> Icons.Default.Settings
    "referral" -> Icons.Default.PersonAdd
    "bug-report" -> Icons.Default.BugReport
    else -> Icons.Default.Circle
}

@Composable
private fun MoreTile(
    title: String,
    viewId: String,
    fallback: ImageVector,
    darkTheme: Boolean,
    onClick: () -> Unit,
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .height(96.dp)
            .clickable(onClick = onClick),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
    ) {
        Column(
            modifier = Modifier.fillMaxSize().padding(8.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center,
        ) {
            ServerOrFallbackIcon(
                viewId = viewId,
                fallback = fallback,
                darkTheme = darkTheme,
                size = 28.dp,
                contentDescription = title,
                tint = MaterialTheme.colorScheme.primary,
            )
            Spacer(Modifier.height(6.dp))
            Text(
                title,
                style = MaterialTheme.typography.titleSmall,
                maxLines = 2,
                textAlign = TextAlign.Center,
            )
        }
    }
}

/** Server PNG icon with Material fallback while loading / on error. */
@Composable
fun ServerOrFallbackIcon(
    viewId: String,
    fallback: ImageVector,
    darkTheme: Boolean,
    size: androidx.compose.ui.unit.Dp,
    contentDescription: String?,
    tint: androidx.compose.ui.graphics.Color = MaterialTheme.colorScheme.primary,
) {
    val url = ru.neriva.app.ui.components.menuAssetUrl(viewId, darkTheme)
    val ctx = LocalContext.current
    Box(modifier = Modifier.size(size), contentAlignment = Alignment.Center) {
        Icon(imageVector = fallback, contentDescription = contentDescription, tint = tint, modifier = Modifier.size(size))
        SubcomposeAsyncImage(
            model = ImageRequest.Builder(ctx).data(url).crossfade(true).build(),
            contentDescription = contentDescription,
            contentScale = ContentScale.Fit,
            modifier = Modifier.size(size),
            loading = { },
            error = { },
            success = { state ->
                androidx.compose.foundation.Image(
                    painter = state.painter,
                    contentDescription = contentDescription,
                    contentScale = ContentScale.Fit,
                    modifier = Modifier.size(size),
                )
            },
        )
    }
}
