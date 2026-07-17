package ru.neriva.app

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Surface
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import ru.neriva.app.ui.NavHost
import ru.neriva.app.ui.theme.NerivaTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        // Apply language on startup
        val lang = LanguageManager.getInterfaceLanguage(this)
        LanguageManager.applyLocale(this, lang)

        setContent {
            val themeMode = remember { mutableStateOf(ThemeManager.getThemeMode(this@MainActivity)) }
            val isDark = when (themeMode.value) {
                ThemeManager.ThemeMode.DARK -> true
                ThemeManager.ThemeMode.LIGHT -> false
                ThemeManager.ThemeMode.SYSTEM -> isSystemInDarkTheme()
            }

            NerivaTheme(darkTheme = isDark) {
                Surface(modifier = Modifier.fillMaxSize()) {
                    NavHost(onThemeChange = { newMode ->
                        themeMode.value = newMode
                        ThemeManager.setThemeMode(this@MainActivity, newMode)
                    }, currentThemeMode = themeMode.value)
                }
            }
        }
    }
}
