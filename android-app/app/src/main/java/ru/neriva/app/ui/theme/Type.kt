package ru.neriva.app.ui.theme

import androidx.compose.material3.Typography
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.Font
import androidx.compose.ui.text.font.FontFamily
import ru.neriva.app.R
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp

/** Bundled Latin, Greek and Cyrillic; Android supplies fallback for other scripts. */
val InterFamily = FontFamily(
    Font(R.font.inter_400, FontWeight.Normal),
    Font(R.font.inter_500, FontWeight.Medium),
    Font(R.font.inter_600, FontWeight.SemiBold),
    Font(R.font.inter_700, FontWeight.Bold),
)

val NerivaTypography = Typography(
    displaySmall = TextStyle(fontFamily = InterFamily, fontWeight = FontWeight.Bold, fontSize = 28.sp, lineHeight = 34.sp),
    displayLarge = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Bold,
        fontSize = 34.sp,
        lineHeight = 42.sp,
        letterSpacing = (-0.25).sp,
    ),
    displayMedium = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Bold,
        fontSize = 28.sp,
        lineHeight = 36.sp,
    ),
    headlineLarge = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.SemiBold,
        fontSize = 24.sp,
        lineHeight = 32.sp,
    ),
    headlineMedium = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.SemiBold,
        fontSize = 20.sp,
        lineHeight = 28.sp,
    ),
    headlineSmall = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.SemiBold,
        fontSize = 18.sp,
        lineHeight = 24.sp,
    ),
    titleLarge = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.SemiBold,
        fontSize = 16.sp,
        lineHeight = 22.sp,
    ),
    titleMedium = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Medium,
        fontSize = 14.sp,
        lineHeight = 20.sp,
    ),
    titleSmall = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Medium,
        fontSize = 13.sp,
        lineHeight = 18.sp,
    ),
    bodyLarge = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Normal,
        fontSize = 16.sp,
        lineHeight = 26.sp,   // 1.55 line-height like web
    ),
    bodyMedium = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Normal,
        fontSize = 14.sp,
        lineHeight = 22.sp,
    ),
    bodySmall = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Normal,
        fontSize = 12.sp,
        lineHeight = 17.sp,
    ),
    labelLarge = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Medium,
        fontSize = 14.sp,
        lineHeight = 18.sp,
    ),
    labelMedium = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Bold,
        fontSize = 12.sp,
        lineHeight = 16.sp,
    ),
    labelSmall = TextStyle(
        fontFamily = InterFamily,
        fontWeight = FontWeight.Bold,
        fontSize = 11.sp,
        lineHeight = 14.sp,
        letterSpacing = 0.3.sp,
    ),
)
