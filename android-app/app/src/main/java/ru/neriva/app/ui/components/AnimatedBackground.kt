package ru.neriva.app.ui.components

import androidx.compose.animation.core.InfiniteRepeatableSpec
import androidx.compose.animation.core.LinearEasing
import androidx.compose.animation.core.RepeatMode
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.BoxScope
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.graphicsLayer
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp

/**
 * Animated background with 3 floating gradient circles (blue, teal, gold)
 * that drift vertically — mirrors the web app's paper-shader-bg animation.
 */
@Composable
fun AnimatedBackground(
    modifier: Modifier = Modifier,
    content: @Composable BoxScope.() -> Unit,
) {
    Box(modifier = modifier.fillMaxSize()) {
        // Circle 1 — large blue, top-left, slow
        FloatingCircle(
            color = Color(0xFF2D5BFF).copy(alpha = 0.22f),
            size = 280.dp,
            alignment = Alignment.TopStart,
            offsetYFraction = 0.08f,
            animDurationMs = 14_000,
        )

        // Circle 2 — teal, top-right, medium
        FloatingCircle(
            color = Color(0xFF00A88F).copy(alpha = 0.18f),
            size = 220.dp,
            alignment = Alignment.TopEnd,
            offsetYFraction = 0.06f,
            animDurationMs = 11_000,
        )

        // Circle 3 — gold, bottom-center, fast
        FloatingCircle(
            color = Color(0xFFF3B84B).copy(alpha = 0.15f),
            size = 260.dp,
            alignment = Alignment.BottomCenter,
            offsetYFraction = 0.10f,
            animDurationMs = 9_000,
        )

        content()
    }
}

@Composable
private fun FloatingCircle(
    color: Color,
    size: Dp,
    alignment: Alignment,
    offsetYFraction: Float,
    animDurationMs: Int,
) {
    val infiniteTransition = rememberInfiniteTransition(label = "bg-float")
    val y by infiniteTransition.animateFloat(
        initialValue = -offsetYFraction,
        targetValue = offsetYFraction,
        animationSpec = infiniteRepeatable(
            animation = tween(animDurationMs, easing = LinearEasing),
            repeatMode = RepeatMode.Reverse,
        ),
        label = "circle-y",
    )

    Box(
        modifier = Modifier
            .fillMaxSize(),
        contentAlignment = alignment,
    ) {
        Box(
            modifier = Modifier
                .size(size)
                .graphicsLayer { translationY = size.toPx() * y }
                .background(
                    Brush.radialGradient(
                        colors = listOf(color, Color.Transparent),
                    ),
                    CircleShape,
                ),
        )
    }
}
