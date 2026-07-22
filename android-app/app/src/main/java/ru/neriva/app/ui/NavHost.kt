package ru.neriva.app.ui

import androidx.compose.foundation.layout.padding
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import ru.neriva.app.ui.nav.Screen
import ru.neriva.app.ui.nav.bottomNavItems
import ru.neriva.app.ui.nav.moreNavItems
import ru.neriva.app.ui.screens.*

@Composable
fun NavHost(onThemeChange: (ru.neriva.app.ThemeManager.ThemeMode) -> Unit = {}, currentThemeMode: ru.neriva.app.ThemeManager.ThemeMode = ru.neriva.app.ThemeManager.ThemeMode.SYSTEM) {
    val navController = rememberNavController()
    val sessionViewModel: SessionViewModel = viewModel()
    val sessionState by sessionViewModel.state.collectAsState()

    val startRoute = if (sessionState.isLoggedIn == true) Screen.Home.route else Screen.Auth.route

    NavHost(
        navController = navController,
        startDestination = startRoute
    ) {
        composable(Screen.Auth.route) {
            AuthScreen(
                onAuthSuccess = {
                    sessionViewModel.refresh()
                    navController.navigate(Screen.Home.route) {
                        popUpTo(Screen.Auth.route) { inclusive = true }
                    }
                },
                onThemeChange = onThemeChange,
                currentThemeMode = currentThemeMode,
            )
        }

        composable(Screen.Onboarding.route) {
            OnboardingScreen(
                onComplete = {
                    sessionViewModel.refresh()
                    navController.navigate(Screen.Home.route) {
                        popUpTo(0)
                    }
                }
            )
        }

        composable(Screen.Home.route) { MainScaffold(navController, Screen.Home.route) { HomeScreen(navController) } }
        composable(Screen.Tutor.route) { MainScaffold(navController, Screen.Tutor.route) { TutorScreen(navController) } }
        composable(Screen.Lesson.route) { MainScaffold(navController, Screen.Lesson.route) { LessonScreen(navController) } }
        composable(Screen.Practice.route) { MainScaffold(navController, Screen.Practice.route) { PracticeScreen(navController) } }
        composable(Screen.Roleplay.route) { MainScaffold(navController, Screen.Roleplay.route) { RoleplayScreen(navController) } }
        composable(Screen.Shadowing.route) { MainScaffold(navController, Screen.Shadowing.route) { ShadowingScreen(navController) } }
        composable(Screen.Pronunciation.route) { MainScaffold(navController, Screen.Pronunciation.route) { PronunciationScreen(navController) } }
        composable(Screen.Words.route) { MainScaffold(navController, Screen.Words.route) { WordsScreen(navController) } }
        composable(Screen.WordGame.route) { MainScaffold(navController, Screen.WordGame.route) { WordGameScreen(navController) } }
        composable(Screen.Spelling.route) { MainScaffold(navController, Screen.Spelling.route) { SpellingScreen(navController) } }
        composable(Screen.Vocabulary.route) { MainScaffold(navController, Screen.Vocabulary.route) { VocabularyScreen(navController) } }
        composable(Screen.Phrasebook.route) { MainScaffold(navController, Screen.Phrasebook.route) { PhrasebookScreen(navController) } }
        composable(Screen.Offline.route) { MainScaffold(navController, Screen.Offline.route) { OfflineScreen(navController) } }
        composable(Screen.LevelTest.route) { MainScaffold(navController, Screen.LevelTest.route) { LevelTestScreen(navController) } }
        composable(Screen.Progress.route) { MainScaffold(navController, Screen.Progress.route) { ProgressScreen(navController) } }
        composable(Screen.Awards.route) { MainScaffold(navController, Screen.Awards.route) { AwardsScreen(navController) } }
        composable(Screen.Leaderboard.route) { MainScaffold(navController, Screen.Leaderboard.route) { LeaderboardScreen(navController) } }
        composable(Screen.Limits.route) { MainScaffold(navController, Screen.Limits.route) { LimitsScreen(navController) } }
        composable(Screen.Mistakes.route) { MainScaffold(navController, Screen.Mistakes.route) { MistakesScreen(navController) } }
        composable(Screen.Tools.route) { MainScaffold(navController, Screen.Tools.route) { ToolsScreen(navController) } }
        composable(Screen.Premium.route) { MainScaffold(navController, Screen.Premium.route) { PremiumScreen(navController) } }
        composable(Screen.Settings.route) { MainScaffold(navController, Screen.Settings.route) { SettingsScreen(navController, onThemeChange, currentThemeMode) } }
        composable(Screen.Referral.route) { MainScaffold(navController, Screen.Referral.route) { ReferralScreen(navController) } }
        composable(Screen.BugReport.route) { MainScaffold(navController, Screen.BugReport.route) { BugReportScreen(navController) } }
        composable(Screen.CompletedLessons.route) { MainScaffold(navController, Screen.CompletedLessons.route) { CompletedLessonsScreen(navController) } }
        composable(Screen.MistakePractice.route) { MainScaffold(navController, Screen.MistakePractice.route) { MistakePracticeScreen(navController) } }
        composable(Screen.RoleplayScenarios.route) { MainScaffold(navController, Screen.RoleplayScenarios.route) { RoleplayScenariosScreen(navController) } }
        composable(Screen.Dashboard.route) { MainScaffold(navController, Screen.Dashboard.route) { DashboardScreen(navController) } }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MainScaffold(
    navController: NavHostController,
    currentRoute: String,
    content: @Composable () -> Unit
) {
    var showMoreSheet by remember { mutableStateOf(false) }
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val current = navBackStackEntry?.destination?.route ?: currentRoute

    Scaffold(
        bottomBar = {
            ru.neriva.app.ui.components.BottomNavBar(
                items = bottomNavItems,
                currentRoute = current,
                onNavigate = { screen ->
                    navController.navigate(screen.route) {
                        popUpTo(Screen.Home.route) { saveState = true }
                        launchSingleTop = true
                        restoreState = true
                    }
                },
                onMore = { showMoreSheet = true },
                moreSelected = showMoreSheet,
            )
        }
    ) { padding ->
        androidx.compose.foundation.layout.Box(modifier = Modifier.padding(padding)) {
            content()
        }
    }

    if (showMoreSheet) {
        MoreSheet(
            items = moreNavItems,
            onDismiss = { showMoreSheet = false },
            onSelect = { screen ->
                showMoreSheet = false
                navController.navigate(screen.route) {
                    popUpTo(Screen.Home.route) { saveState = true }
                    launchSingleTop = true
                }
            }
        )
    }
}

