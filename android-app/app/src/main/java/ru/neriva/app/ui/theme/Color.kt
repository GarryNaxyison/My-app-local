package ru.neriva.app.ui.theme

import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color

// ── Light theme ──────────────────────────────────────────────────────────────
val LightBg = Color(0xFFEDF7FF)
val LightSurface = Color(0xFFFFFFFF)
val LightSurfaceMuted = Color(0xFFF2F6FC)
val LightText = Color(0xFF07111F)
val LightMuted = Color(0xFF5C6C82)
val LightFaint = Color(0xFF8391A5)
val LightLine = Color(0x1F0F172A)
val LightPrimary = Color(0xFF2D5BFF)
val LightTeal = Color(0xFF00A88F)
val LightGold = Color(0xFFF3B84B)
val LightRose = Color(0xFFE45A6F)
val LightPlum = Color(0xFF7C5CFF)
val LightMint = Color(0xFFDDF8EF)

// ── Dark theme ───────────────────────────────────────────────────────────────
val DarkBg = Color(0xFF020511)
val DarkSurface = Color(0xFF0D1627)
val DarkSurfaceMuted = Color(0xFF152137)
val DarkText = Color(0xFFF6F9FF)
val DarkMuted = Color(0xFFA9B7CD)
val DarkFaint = Color(0xFF7789A6)
val DarkLine = Color(0x24B8CDEF)
val DarkPrimary = Color(0xFF7AA2FF)
val DarkTeal = Color(0xFF45D8BB)
val DarkGold = Color(0xFFF4C96D)
val DarkRose = Color(0xFFFF8395)
val DarkPlum = Color(0xFF9D85FF)
val DarkMint = Color(0xFF1A3D33)

// ── Shared accent colors ─────────────────────────────────────────────────────
val AccentBlue = Color(0xFF2D5BFF)
val AccentTeal = Color(0xFF00A88F)
val AccentGold = Color(0xFFF3B84B)
val AccentMint = Color(0xFFDDF8EF)
val AccentRose = Color(0xFFE45A6F)

// ── Gradient brushes ─────────────────────────────────────────────────────────

/** CTA button gradient: gold → mint */
val GradientButtonBrush = Brush.horizontalGradient(
    colors = listOf(AccentGold, AccentMint),
)

/** Level / XP progress bar: primary → teal → gold */
val GradientProgressBrush = Brush.horizontalGradient(
    colors = listOf(AccentBlue, AccentTeal, AccentGold),
)

/** Subtle gradient for glass-card borders or highlights */
val GradientGlassBorder = Brush.linearGradient(
    colors = listOf(
        Color.White.copy(alpha = 0.24f),
        Color.White.copy(alpha = 0.08f),
    ),
)
