package ru.neriva.app.ui.nav

import androidx.annotation.StringRes
import ru.neriva.app.R

sealed class Screen(val route: String, @StringRes val titleRes: Int, val viewId: String = route) {
    data object Auth : Screen("auth", R.string.auth_login)
    data object Onboarding : Screen("onboarding", R.string.auth_finish_account)
    data object Home : Screen("home", R.string.today)
    data object Tutor : Screen("tutor", R.string.ai_tutor)
    data object Lesson : Screen("lesson", R.string.new_lesson)
    data object Practice : Screen("practice", R.string.practice)
    data object Roleplay : Screen("roleplay", R.string.roleplay)
    data object Shadowing : Screen("shadowing", R.string.shadowing)
    data object Pronunciation : Screen("pronunciation", R.string.pronunciation)
    data object Words : Screen("words", R.string.learn_words)
    data object WordGame : Screen("word-game", R.string.word_game)
    data object Spelling : Screen("spelling", R.string.spelling)
    data object Vocabulary : Screen("vocabulary", R.string.vocabulary)
    data object Phrasebook : Screen("phrasebook", R.string.phrasebook)
    data object Offline : Screen("offline", R.string.offline_decks)
    data object LevelTest : Screen("level", R.string.level_test)
    data object Progress : Screen("progress", R.string.progress)
    data object Awards : Screen("awards", R.string.awards)
    data object Leaderboard : Screen("leaderboard", R.string.leaders)
    data object Limits : Screen("limits", R.string.limits)
    data object Mistakes : Screen("mistakes", R.string.mistakes)
    data object Tools : Screen("tools", R.string.tools)
    data object Premium : Screen("premium", R.string.premium)
    data object Settings : Screen("settings", R.string.settings)
    data object Referral : Screen("referral", R.string.referrals)
    data object BugReport : Screen("bug-report", R.string.bug_report_sent)
    data object CompletedLessons : Screen("completed-lessons", R.string.tutor_completed_lessons, viewId = "progress")
    data object MistakePractice : Screen("mistake-practice", R.string.mistakes, viewId = "mistakes")
    data object RoleplayScenarios : Screen("roleplay-scenarios", R.string.roleplay, viewId = "roleplay")
    data object Dashboard : Screen("dashboard", R.string.dashboard, viewId = "progress")
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
