package ru.neriva.app.ui.components

import androidx.compose.animation.core.LinearEasing
import androidx.compose.animation.core.RepeatMode
import androidx.compose.animation.core.animateFloat
import androidx.compose.animation.core.infiniteRepeatable
import androidx.compose.animation.core.rememberInfiniteTransition
import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.PlayArrow
import androidx.compose.material.icons.filled.VolumeUp
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.graphicsLayer
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import ru.neriva.app.audio.AudioPlayer

private enum class PlaybackState { IDLE, LOADING, PLAYING }

/**
 * Compact TTS playback button mirroring the web `AudioWaveButton`.
 *
 * On tap, it POSTs JSON to [endpoint] (authenticated via the shared cookie jar),
 * receives a WAV blob, and plays it. Shows animated equalizer bars while playing,
 * a spinner while loading, and a speaker icon when idle.
 *
 * @param endpoint API path, e.g. "/api/tools/translator-speech"
 * @param jsonBody request body, e.g. {"text":"hi","target_language":"en"}
 */
@Composable
fun AudioWaveButton(
    endpoint: String,
    jsonBody: String,
    modifier: Modifier = Modifier,
    label: String? = null,
) {
    val scope = rememberCoroutineScope()
    var state by remember { mutableStateOf(PlaybackState.IDLE) }
    val player = remember { AudioPlayer.get() }

    Row(
        modifier = modifier.clip(CircleShape).clickable {
            when (state) {
                PlaybackState.IDLE, PlaybackState.LOADING -> {
                    state = PlaybackState.LOADING
                    scope.launch {
                        AudioPlayer.playBlob(endpoint, jsonBody) { state = PlaybackState.IDLE }
                    }
                }
                PlaybackState.PLAYING -> {
                    player.stop()
                    state = PlaybackState.IDLE
                }
            }
        }.padding(horizontal = 6.dp, vertical = 4.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(6.dp),
    ) {
        Box(
            modifier = Modifier
                .size(34.dp)
                .clip(CircleShape)
                .background(MaterialTheme.colorScheme.secondary.copy(alpha = 0.14f)),
            contentAlignment = Alignment.Center,
        ) {
            when (state) {
                PlaybackState.IDLE -> Icon(
                    Icons.Filled.VolumeUp,
                    contentDescription = "Play audio",
                    tint = MaterialTheme.colorScheme.secondary,
                    modifier = Modifier.size(18.dp),
                )
                PlaybackState.LOADING -> CircularProgressIndicator(
                    modifier = Modifier.size(18.dp),
                    strokeWidth = 2.dp,
                )
                PlaybackState.PLAYING -> EqualizerBars()
            }
        }
        if (label != null) {
            Text(
                label,
                style = MaterialTheme.typography.labelMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant,
            )
        }
    }
}

/** Convenience: play a single text snippet via translator-speech. */
@Composable
fun PlayTextButton(
    text: String,
    targetLanguage: String,
    modifier: Modifier = Modifier,
) {
    val json = """{"text":"$text","target_language":"$targetLanguage"}"""
    AudioWaveButton(
        endpoint = "/api/tools/translator-speech",
        jsonBody = json,
        modifier = modifier,
    )
}

@Composable
private fun EqualizerBars() {
    val transition = rememberInfiniteTransition(label = "eq")
    Row(
        horizontalArrangement = Arrangement.spacedBy(2.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        repeat(4) { i ->
            val phase by transition.animateFloat(
                initialValue = 0.35f,
                targetValue = 1f,
                animationSpec = infiniteRepeatable(
                    animation = tween(360 + i * 90, easing = LinearEasing),
                    repeatMode = RepeatMode.Reverse,
                ),
                label = "bar-$i",
            )
            Box(
                modifier = Modifier
                    .size(width = 3.dp, height = 16.dp)
                    .graphicsLayer { scaleY = phase.coerceIn(0.3f, 1f) }
                    .background(MaterialTheme.colorScheme.secondary, RoundedCornerShape(2.dp)),
            )
        }
    }
}
