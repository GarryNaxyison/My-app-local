package ru.neriva.app.ui.components

import androidx.compose.foundation.layout.size
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.ColorFilter
import androidx.compose.ui.graphics.ColorMatrix
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.dp
import coil.compose.AsyncImage
import coil.request.ImageRequest

/**
 * Loads an icon PNG from the NERIVA web assets CDN.
 * Web path: /app/assets/icon-{view}-{light|dark}.png
 */
@Composable
fun ServerIcon(
    viewId: String,
    modifier: Modifier = Modifier,
    darkTheme: Boolean = false,
    size: Dp = 28.dp,
    contentDescription: String? = null,
) {
    val themeSuffix = if (darkTheme) "-dark" else "-light"
    val url = "https://neriva.ru/app/assets/icon-$viewId$themeSuffix.png"
    val ctx = LocalContext.current
    AsyncImage(
        model = ImageRequest.Builder(ctx).data(url).crossfade(true).build(),
        contentDescription = contentDescription,
        contentScale = ContentScale.Fit,
        modifier = modifier.size(size),
    )
}

/** Load the award image for a level (1..20) from /app/assets/award-NN.png */
@Composable
fun AwardIcon(
    level: Int,
    modifier: Modifier = Modifier,
    size: Dp = 64.dp,
    grayscale: Boolean = false,
) {
    val normalized = level.coerceIn(1, 20)
    val key = normalized.toString().padStart(2, '0')
    val url = "https://neriva.ru/app/assets/award-$key.png"
    val ctx = LocalContext.current
    AsyncImage(
        model = ImageRequest.Builder(ctx).data(url).crossfade(true).build(),
        contentDescription = "Award level $normalized",
        contentScale = ContentScale.Fit,
        colorFilter = if (grayscale) ColorFilter.colorMatrix(ColorMatrix().apply { setToSaturation(0f) }) else null,
        modifier = modifier.size(size),
    )
}
