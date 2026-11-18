package ru.neriva.app.ui.components

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.size
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.outlined.MenuBook
import androidx.compose.material.icons.outlined.Forum
import androidx.compose.material.icons.outlined.Home
import androidx.compose.material.icons.outlined.MoreHoriz
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.NavigationBarItemDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import ru.neriva.app.R
import ru.neriva.app.ui.nav.Screen

@Composable
fun BottomNavBar(
    modifier: Modifier = Modifier,
    items: List<Screen>,
    currentRoute: String,
    onNavigate: (Screen) -> Unit,
    onMore: () -> Unit,
    moreSelected: Boolean = false,
) {
    val colors = MaterialTheme.colorScheme
    val itemColors = NavigationBarItemDefaults.colors(
        selectedIconColor = colors.primary,
        selectedTextColor = colors.primary,
        unselectedIconColor = colors.onSurfaceVariant,
        unselectedTextColor = colors.onSurfaceVariant,
        indicatorColor = colors.primary.copy(alpha = 0.13f),
    )
    NavigationBar(modifier = modifier, containerColor = colors.surface, tonalElevation = 0.dp) {
        items.forEach { screen ->
            NavigationBarItem(
                selected = currentRoute == screen.route,
                onClick = { onNavigate(screen) },
                icon = {
                    if (screen == Screen.Tutor) {
                        Image(painterResource(R.drawable.neriva_brand_symbol), null, Modifier.size(28.dp, 24.dp))
                    } else {
                        val icon = when (screen) {
                            Screen.Home -> Icons.Outlined.Home
                            Screen.Practice -> Icons.Outlined.Forum
                            Screen.Words -> Icons.AutoMirrored.Outlined.MenuBook
                            else -> Icons.Outlined.MoreHoriz
                        }
                        Icon(icon, null, Modifier.size(22.dp))
                    }
                },
                label = { NavLabel(stringResource(screen.titleRes)) },
                colors = itemColors,
            )
        }
        NavigationBarItem(
            selected = moreSelected,
            onClick = onMore,
            icon = { Icon(Icons.Outlined.MoreHoriz, null, Modifier.size(22.dp)) },
            label = { NavLabel(stringResource(R.string.more)) },
            colors = itemColors,
        )
    }
}

@Composable
private fun NavLabel(text: String) {
    Text(text, style = MaterialTheme.typography.labelSmall, textAlign = TextAlign.Center)
}
