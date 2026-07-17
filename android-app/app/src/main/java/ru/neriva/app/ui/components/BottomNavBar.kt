package ru.neriva.app.ui.components

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.size
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Dashboard
import androidx.compose.material.icons.filled.Forum
import androidx.compose.material.icons.filled.MenuBook
import androidx.compose.material.icons.filled.MoreHoriz
import androidx.compose.material.icons.filled.School
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.foundation.layout.BoxScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil.compose.AsyncImage
import coil.compose.AsyncImagePainter
import coil.compose.SubcomposeAsyncImage
import coil.request.ImageRequest
import androidx.compose.runtime.getValue

data class NavItem(val label: String, val icon: ImageVector, val viewId: String)

val navItems = listOf(
    NavItem("Home", Icons.Filled.Dashboard, "home"),
    NavItem("Tutor", Icons.Filled.School, "lesson"),
    NavItem("Practice", Icons.Filled.Forum, "practice"),
    NavItem("Words", Icons.Filled.MenuBook, "words"),
    NavItem("More", Icons.Filled.MoreHoriz, "settings"),
)

/**
 * Bottom navigation with 4 tabs + More button, using server PNG icons
 * (same assets as the web app) with Material icon fallbacks.
 */
@Composable
fun BottomNavBar(
    modifier: Modifier = Modifier,
    selectedIndex: Int = 0,
    onSelect: (Int) -> Unit = {},
) {
    val darkTheme = isSystemInDarkTheme()
    NavigationBar(
        modifier = modifier,
        containerColor = MaterialTheme.colorScheme.surface,
        tonalElevation = 0.dp,
    ) {
        navItems.forEachIndexed { index, item ->
            NavigationBarItem(
                icon = {
                    NavIcon(
                        viewId = item.viewId,
                        fallback = item.icon,
                        darkTheme = darkTheme,
                        selected = index == selectedIndex,
                        contentDescription = item.label,
                    )
                },
                label = {
                    Text(
                        text = item.label,
                        fontSize = 11.sp,
                        fontWeight = if (index == selectedIndex) FontWeight.Bold else FontWeight.Medium,
                    )
                },
                selected = index == selectedIndex,
                onClick = { onSelect(index) },
                colors = NavigationBarItemDefaults.colors(
                    selectedIconColor = MaterialTheme.colorScheme.primary,
                    selectedTextColor = MaterialTheme.colorScheme.primary,
                    unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                    unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant,
                    indicatorColor = MaterialTheme.colorScheme.primary.copy(alpha = 0.10f),
                ),
            )
        }
    }
}

/** Server PNG icon with Material fallback while loading / on error. */
@Composable
private fun NavIcon(
    viewId: String,
    fallback: ImageVector,
    darkTheme: Boolean,
    selected: Boolean,
    contentDescription: String?,
) {
    val themeSuffix = if (darkTheme) "-dark" else "-light"
    val url = "https://neriva.ru/app/assets/icon-$viewId$themeSuffix.png"
    val ctx = LocalContext.current
    val tint = if (selected) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.onSurfaceVariant
    Box(modifier = Modifier.size(24.dp), contentAlignment = Alignment.Center) {
        // Fallback Material icon underneath (shown until server icon loads / on error)
        Icon(imageVector = fallback, contentDescription = contentDescription, tint = tint, modifier = Modifier.size(22.dp))
        SubcomposeAsyncImage(
            model = ImageRequest.Builder(ctx).data(url).crossfade(true).build(),
            contentDescription = contentDescription,
            contentScale = ContentScale.Fit,
            modifier = Modifier.size(24.dp),
            loading = { /* keep fallback visible */ },
            error = { /* keep fallback visible */ },
            success = { state ->
                AsyncImage(
                    model = state.painter,
                    contentDescription = contentDescription,
                    contentScale = ContentScale.Fit,
                    modifier = Modifier.size(24.dp),
                )
            },
        )
    }
}
