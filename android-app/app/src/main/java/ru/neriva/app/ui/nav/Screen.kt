package ru.neriva.app.ui.nav

sealed class Screen(val route: String, val title: String, val viewId: String = route) {
    data object Auth : Screen("auth", "Sign In")
    data object Onboarding : Screen("onboarding", "Setup")
    data object Home : Screen("home", "Today")
    data object Tutor : Screen("tutor", "AI Tutor")
    data object Lesson : Screen("lesson", "Lesson")
    data object Practice : Screen("practice", "Practice")
    data object Roleplay : Screen("roleplay", "Roleplay")
    data object Shadowing : Screen("shadowing", "Shadowing")
    data object Pronunciation : Screen("pronunciation", "Pronunciation")
    data object Words : Screen("words", "Learn Words")
    data object WordGame : Screen("word-game", "Word Game")
    data object Spelling : Screen("spelling", "Spelling")
    data object Vocabulary : Screen("vocabulary", "Vocabulary")
    data object Phrasebook : Screen("phrasebook", "Phrasebook")
    data object Offline : Screen("offline", "Offline")
    data object LevelTest : Screen("level", "Level Test")
    data object Progress : Screen("progress", "Progress")
    data object Awards : Screen("awards", "Awards")
    data object Leaderboard : Screen("leaderboard", "Leaderboard")
    data object Limits : Screen("limits", "Limits")
    data object Mistakes : Screen("mistakes", "Mistakes")
    data object Tools : Screen("tools", "Tools")
    data object Premium : Screen("premium", "Premium")
    data object Settings : Screen("settings", "Settings")
    data object Referral : Screen("referral", "Referral")
    data object BugReport : Screen("bug-report", "Bug Report")
    data object CompletedLessons : Screen("completed-lessons", "Completed Lessons", viewId = "progress")
    data object MistakePractice : Screen("mistake-practice", "Mistake Practice", viewId = "mistakes")
    data object RoleplayScenarios : Screen("roleplay-scenarios", "Roleplay", viewId = "roleplay")
    data object Dashboard : Screen("dashboard", "Dashboard", viewId = "progress")
}

val bottomNavItems = listOf(
    Screen.Home,
    Screen.Tutor,
    Screen.Practice,
    Screen.Words,
)

val moreNavItems = listOf(
    Screen.RoleplayScenarios,
    Screen.Pronunciation,
    Screen.Shadowing,
    Screen.Spelling,
    Screen.Vocabulary,
    Screen.Phrasebook,
    Screen.Offline,
    Screen.Mistakes,
    Screen.MistakePractice,
    Screen.Tools,
    Screen.LevelTest,
    Screen.Progress,
    Screen.Awards,
    Screen.Leaderboard,
    Screen.Limits,
    Screen.CompletedLessons,
    Screen.Dashboard,
    Screen.Premium,
    Screen.Referral,
    Screen.Settings,
)
