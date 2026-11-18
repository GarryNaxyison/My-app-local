package ru.neriva.app.ui.components

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.ui.draw.drawBehind
import androidx.compose.ui.geometry.CornerRadius
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.unit.LayoutDirection
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.selection.selectable
import androidx.compose.foundation.selection.selectableGroup
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.outlined.MenuBook
import androidx.compose.material.icons.outlined.Forum
import androidx.compose.material.icons.outlined.Home
import androidx.compose.material.icons.outlined.MoreHoriz
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.text.style.TextOverflow
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
    val active = if (moreSelected) items.size else items.indexOfFirst { it.route == currentRoute }
    val activePosition by animateFloatAsState(
        active.coerceAtLeast(0).toFloat(), tween(NerivaTransitionMillis), label = "dockActivePosition",
    )
    val activeAlpha by animateFloatAsState(if (active >= 0) 1f else 0f, label = "dockActiveAlpha")
    Surface(modifier, color = colors.surface) {
        Column {
            HorizontalDivider(color = colors.outline)
            Box(Modifier.fillMaxWidth().windowInsetsPadding(NavigationBarDefaults.windowInsets), contentAlignment = Alignment.Center) {
                Row(
                    Modifier.widthIn(max = 720.dp).fillMaxWidth().height(IntrinsicSize.Min)
                        .padding(horizontal = 6.dp, vertical = 4.dp)
                        .drawBehind {
                            val column = size.width / (items.size + 1)
                            val start = column * (activePosition + .5f)
                            val center = if (layoutDirection == LayoutDirection.Rtl) size.width - start else start
                            val width = 19.dp.toPx()
                            val height = 2.dp.toPx()
                            drawRoundRect(
                                color = colors.primary.copy(alpha = activeAlpha),
                                topLeft = Offset(center - width / 2, size.height - height),
                                size = Size(width, height),
                                cornerRadius = CornerRadius(height / 2),
                            )
                        }.selectableGroup(),
                ) {
                    items.forEach { screen ->
                        DockTab(
                            stringResource(screen.titleRes),
                            !moreSelected && currentRoute == screen.route,
                            { onNavigate(screen) },
                            Modifier.weight(1f).fillMaxHeight(),
                        ) { tint ->
                            if (screen == Screen.Tutor) {
                                Image(painterResource(R.drawable.neriva_brand_symbol), null, Modifier.size(34.dp, 22.dp))
                            } else {
                                val vector = when (screen) {
                                    Screen.Home -> Icons.Outlined.Home
                                    Screen.Practice -> Icons.Outlined.Forum
                                    Screen.Words -> Icons.AutoMirrored.Outlined.MenuBook
                                    else -> Icons.Outlined.MoreHoriz
                                }
                                Icon(vector, null, Modifier.size(22.dp), tint = tint)
                            }
                        }
                    }
                    DockTab(stringResource(R.string.more), moreSelected, onMore, Modifier.weight(1f).fillMaxHeight()) { tint ->
                        Icon(Icons.Outlined.MoreHoriz, null, Modifier.size(22.dp), tint = tint)
                    }
                }
            }
        }
    }
}

@Composable
private fun DockTab(label: String, selected: Boolean, onClick: () -> Unit, modifier: Modifier, icon: @Composable (Color) -> Unit) {
    val colors = MaterialTheme.colorScheme
    val tint by animateColorAsState(if (selected) colors.primary else colors.onSurfaceVariant, label = "dockTint")
    val indicator by animateColorAsState(
        if (selected) colors.primary.copy(alpha = 0.13f) else Color.Transparent,
        label = "dockIndicator",
    )
    Column(
        modifier.heightIn(min = 65.dp).clip(RoundedCornerShape(14.dp))
            .selectable(selected, role = Role.Tab, onClick = onClick).padding(horizontal = 2.dp, vertical = 4.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(4.dp, Alignment.Top),
    ) {
        Box(
            Modifier.size(width = 44.dp, height = 30.dp).clip(RoundedCornerShape(10.dp)).background(indicator),
            contentAlignment = Alignment.Center,
        ) { icon(tint) }
        Text(label, style = MaterialTheme.typography.labelSmall, color = tint, textAlign = TextAlign.Center,
            maxLines = 2, overflow = TextOverflow.Ellipsis)
    }
}
