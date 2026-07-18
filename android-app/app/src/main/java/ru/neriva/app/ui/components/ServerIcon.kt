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
 * Centralised asset URL helpers — mirror the web app's `menuAssetUrl()`
 * (web-react/src/App.tsx:9154) and ComfyUI asset generator output.
 *
 * URL pattern: /app/assets/icon-{slug}-{light|dark}.png?v={MENU_ASSET_VERSION}
 *   - tutor → brand-logo-mini.png (no theme suffix, no icon- prefix)
 *   - dashboard → progress slug
 */
internal const val MENU_ASSET_VERSION = "flux2-20260526-pronunciation-phrases-offline-roleplay"
internal const val ASSETS_BASE = "https://neriva.ru/app/assets"

/** Same slug mapping as the web app's `assetSlug()`. */
fun assetSlug(viewId: String): String =
    if (viewId == "dashboard") "progress" else viewId

/**
 * Build the server URL for a menu/nav icon — matches the web app exactly.
 * `tutor` uses the brand mini-logo; everything else uses icon-{slug}-{theme}.
 */
fun menuAssetUrl(viewId: String, darkTheme: Boolean): String {
    if (viewId == "tutor") return "$ASSETS_BASE/brand-logo-mini.png?v=$MENU_ASSET_VERSION"
    val theme = if (darkTheme) "dark" else "light"
    val slug = assetSlug(viewId)
    return "$ASSETS_BASE/icon-$slug-$theme.png?v=$MENU_ASSET_VERSION"
}

/**
 * Loads an icon PNG from the NERIVA web assets CDN.
 * Web path: /app/assets/icon-{view}-{light|dark}.png?v=...
 */
@Composable
fun ServerIcon(
    viewId: String,
    modifier: Modifier = Modifier,
    darkTheme: Boolean = false,
    size: Dp = 28.dp,
    contentDescription: String? = null,
) {
    val url = menuAssetUrl(viewId, darkTheme)
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
