package ru.neriva.app.ui.components

import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.interaction.MutableInteractionSource
import androidx.compose.foundation.interaction.collectIsPressedAsState
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.outlined.ArrowForward
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.graphicsLayer
import androidx.compose.ui.graphics.luminance
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.semantics.role
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.unit.dp

@Composable
fun NerivaAction(
    label: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    onCover: Boolean = false,
    optionNumber: Int? = null,
) {
    val colors = MaterialTheme.colorScheme
    val isOption = optionNumber != null
    val source = remember { MutableInteractionSource() }
    val pressed by source.collectIsPressedAsState()
    val scale by animateFloatAsState(if (pressed) 0.97f else 1f, label = "actionPress")
    val foreground = when {
        !enabled -> colors.onSurfaceVariant
        onCover -> Color(0xFF202775)
        isOption -> colors.onSurface
        else -> if (colors.primary.luminance() > 0.2f) Color(0xFF0B1020) else Color.White
    }
    val fill = when {
        !enabled -> colors.surfaceVariant
        onCover -> Color(0xFFF2F4FF)
        isOption -> colors.surfaceVariant
        else -> colors.primary
    }
    Surface(
        onClick = onClick,
        modifier = modifier.heightIn(min = 52.dp).semantics { role = Role.Button }
            .graphicsLayer { scaleX = scale; scaleY = scale },
        shape = RoundedCornerShape(14.dp),
        color = fill,
        contentColor = foreground,
        border = if (isOption) BorderStroke(1.dp, colors.outline) else null,
        enabled = enabled,
        interactionSource = source,
    ) {
        Row(
            Modifier.padding(horizontal = 13.dp, vertical = 9.dp),
            horizontalArrangement = Arrangement.spacedBy(10.dp),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            if (optionNumber != null) {
                ActionInset(colors.surface, colors.primary) {
                    Text(optionNumber.toString(), style = MaterialTheme.typography.labelMedium)
                }
            }
            Text(label, style = MaterialTheme.typography.labelLarge, modifier = Modifier.weight(1f, fill = !onCover))
            if (!isOption) {
                ActionInset(foreground.copy(alpha = 0.14f), foreground) {
                    Icon(Icons.AutoMirrored.Outlined.ArrowForward, null, Modifier.size(16.dp))
                }
            }
        }
    }
}

@Composable
private fun ActionInset(color: Color, contentColor: Color, content: @Composable () -> Unit) {
    Surface(shape = RoundedCornerShape(8.dp), color = color, contentColor = contentColor) {
        Box(Modifier.size(28.dp), contentAlignment = Alignment.Center) { content() }
    }
}
