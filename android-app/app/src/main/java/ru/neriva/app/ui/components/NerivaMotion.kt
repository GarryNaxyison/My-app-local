package ru.neriva.app.ui.components

import androidx.compose.animation.core.FastOutSlowInEasing
import androidx.compose.animation.core.animateFloatAsState
import androidx.compose.animation.core.tween
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.drawWithCache
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.graphics.drawscope.scale
import androidx.compose.ui.graphics.drawscope.translate
import androidx.compose.ui.unit.LayoutDirection
import androidx.compose.ui.unit.dp

internal const val NerivaTransitionMillis = 260
internal val NerivaAmber = Color(0xFFFFC47D)

@Composable
internal fun Modifier.nerivaWingCover(accent: Color = Color(0xFF343DB0)): Modifier {
    var revealed by remember { mutableStateOf(false) }
    val reveal by animateFloatAsState(
        if (revealed) 1f else 0f,
        tween(durationMillis = 560, easing = FastOutSlowInEasing),
        label = "coverWingReveal",
    )
    LaunchedEffect(Unit) { revealed = true }
    return drawWithCache {
        val backdrop = Brush.horizontalGradient(listOf(Color(0xFF161B48), accent))
        val w = size.width
        val h = size.height
        val wings = List(6) { index ->
            val feather = index.toFloat()
            Path().apply {
                moveTo(w * (.64f + .022f * feather), h * (.95f - .052f * feather))
                cubicTo(
                    w * (.81f + .014f * feather), h * (.92f - .09f * feather),
                    w * (.91f + .016f * feather), h * (.55f - .085f * feather),
                    w * 1.035f, h * (.13f - .085f * feather),
                )
            }
        }
        val line = Stroke(1.2.dp.toPx())
        val rtl = layoutDirection == LayoutDirection.Rtl
        onDrawBehind {
            if (rtl) scale(-1f, 1f, pivot = center) { drawRect(backdrop) } else drawRect(backdrop)
            val drawWing: androidx.compose.ui.graphics.drawscope.DrawScope.() -> Unit = {
                val fraction = reveal
                translate(left = w * .10f * (1f - fraction)) {
                    wings.forEachIndexed { index, path ->
                        val color = if (index == 2) NerivaAmber else Color(0xFFA9B7FF)
                        drawPath(path, color.copy(alpha = (if (index == 2) .38f else .21f) * fraction), style = line)
                    }
                    drawCircle(NerivaAmber.copy(alpha = .8f * fraction), 2.4.dp.toPx(), Offset(w * .89f, h * .33f))
                }
            }
            if (rtl) scale(-1f, 1f, pivot = center, block = drawWing) else drawWing()
        }
    }
}
