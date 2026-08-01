package ru.neriva.app.ui.components

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.size
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.BookOnline
import androidx.compose.material.icons.filled.Dashboard
import androidx.compose.material.icons.filled.Forum
import androidx.compose.material.icons.filled.MenuBook
import androidx.compose.material.icons.filled.MoreHoriz
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil.compose.AsyncImage
import coil.compose.SubcomposeAsyncImage
import coil.request.ImageRequest
import androidx.compose.runtime.getValue
import ru.neriva.app.ui.nav.Screen
import ru.neriva.app.R

/**
 * Asset URL builder — mirrors web app's `menuAssetUrl()` (App.tsx:9154).
 *
 * - tutor → /app/assets/brand-logo-mini.png (no theme suffix, no icon- prefix)
 * - dashboard → progress slug
 * - everything else → /app/assets/icon-{slug}-{light|dark}.png
 *
 * Cache-buster `?v=` matches MENU_ASSET_VERSION from the web app so the
 * client fetches the same revision of the generated ComfyUI assets.
 */
private fun fallbackIcon(viewId: String): ImageVector = when (viewId) {
    "home" -> Icons.Filled.Dashboard
    "tutor" -> Icons.Filled.BookOnline
    "practice" -> Icons.Filled.Forum
    "words" -> Icons.Filled.MenuBook
    else -> Icons.Filled.MoreHoriz
}

/**
 * Bottom navigation bar using server PNG assets generated via ComfyUI
 * (identical to the web app's bottom nav). Falls back to Material icons
 * while loading or if the network image fails.
 */
@Composable
fun BottomNavBar(
    modifier: Modifier = Modifier,
    items: List<Screen>,
    currentRoute: String,
    onNavigate: (Screen) -> Unit,
    onMore: () -> Unit,
    moreSelected: Boolean = false,
) {
    val darkTheme = isSystemInDarkTheme()
    NavigationBar(
        modifier = modifier,
        containerColor = MaterialTheme.colorScheme.surface,
        tonalElevation = 0.dp,
    ) {
        items.forEach { screen ->
            val selected = currentRoute == screen.route
            val label = stringResource(screen.titleRes)
            NavigationBarItem(
                icon = {
                    NavIcon(
                        viewId = screen.viewId,
                        fallback = fallbackIcon(screen.viewId),
                        darkTheme = darkTheme,
                        selected = selected,
                        contentDescription = label,
                    )
                },
                label = {
                    Text(
                        text = label,
                        fontSize = 11.sp,
                        fontWeight = if (selected) FontWeight.Bold else FontWeight.Medium,
                    )
                },
                selected = selected,
                onClick = { onNavigate(screen) },
                colors = NavigationBarItemDefaults.colors(
                    selectedIconColor = MaterialTheme.colorScheme.primary,
                    selectedTextColor = MaterialTheme.colorScheme.primary,
                    unselectedIconColor = MaterialTheme.colorScheme.onSurfaceVariant,
                    unselectedTextColor = MaterialTheme.colorScheme.onSurfaceVariant,
                    indicatorColor = MaterialTheme.colorScheme.primary.copy(alpha = 0.10f),
                ),
            )
        }
        // ru.neriva.app.NERIVAApp.instance.getString(R.string.more) button — opens the MoreSheet (always uses a Material icon)
        NavigationBarItem(
            icon = {
                Icon(
                    Icons.Filled.MoreHoriz,
                    contentDescription = stringResource(R.string.more),
                    tint = if (moreSelected) MaterialTheme.colorScheme.primary
                           else MaterialTheme.colorScheme.onSurfaceVariant,
                )
            },
            label = {
                Text(
                    text = stringResource(R.string.more),
                    fontSize = 11.sp,
                    fontWeight = if (moreSelected) FontWeight.Bold else FontWeight.Medium,
                )
            },
            selected = moreSelected,
            onClick = onMore,
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

/** Server PNG icon overlaid on a Material fallback (visible while loading / on error). */
@Composable
private fun NavIcon(
    viewId: String,
    fallback: ImageVector,
    darkTheme: Boolean,
    selected: Boolean,
    contentDescription: String?,
) {
    val url = menuAssetUrl(viewId, darkTheme)
    val ctx = LocalContext.current
    val tint = if (selected) MaterialTheme.colorScheme.primary
               else MaterialTheme.colorScheme.onSurfaceVariant
    Box(modifier = Modifier.size(26.dp), contentAlignment = Alignment.Center) {
        // Fallback Material icon underneath — visible until the PNG loads / on error
        Icon(
            imageVector = fallback,
            contentDescription = contentDescription,
            tint = tint,
            modifier = Modifier.size(22.dp),
        )
        SubcomposeAsyncImage(
            model = ImageRequest.Builder(ctx).data(url).crossfade(true).build(),
            contentDescription = contentDescription,
            contentScale = ContentScale.Fit,
            modifier = Modifier.size(26.dp),
            loading = { /* keep fallback visible */ },
            error = { /* keep fallback visible */ },
            success = { state ->
                androidx.compose.foundation.Image(
                    painter = state.painter,
                    contentDescription = contentDescription,
                    contentScale = ContentScale.Fit,
                    modifier = Modifier.size(26.dp),
                )
            },
        )
    }
}
