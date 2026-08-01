package ru.neriva.app.ui.screens

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp

data class RoleplayScenario(
    val id: String,
    val title: String,
    val description: String,
    val icon: androidx.compose.ui.graphics.vector.ImageVector
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun RoleplayScenariosScreen(navController: androidx.navigation.NavHostController) {
    val scenarios = listOf(
        RoleplayScenario("restaurant", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_restaurant), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_restaurant_description), Icons.Default.Restaurant),
        RoleplayScenario("work", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_work), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_work_description), Icons.Default.Work),
        RoleplayScenario("travel", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_travel), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_travel_description), Icons.Default.Flight),
        RoleplayScenario("exam", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_exam), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_exam_description), Icons.Default.Quiz),
        RoleplayScenario("small-talk", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_small_talk), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_small_talk_description), Icons.Default.Chat),
        RoleplayScenario("hotel", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_hotel), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_hotel_description), Icons.Default.Hotel),
        RoleplayScenario("shopping", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_shopping), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_shopping_description), Icons.Default.ShoppingBag),
        RoleplayScenario("doctor", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_doctor), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_doctor_description), Icons.Default.LocalHospital),
        RoleplayScenario("job-interview", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_job_interview), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_job_interview_description), Icons.Default.BusinessCenter),
        RoleplayScenario("bank", ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_bank), ru.neriva.app.NERIVAApp.instance.getString(R.string.roleplay_scenario_bank_description), Icons.Default.AccountBalance),
    )

    Scaffold(topBar = {
        TopAppBar(title = { Text(stringResource(R.string.roleplay)) })
    }) { padding ->
        LazyColumn(
            Modifier.padding(padding).fillMaxSize().padding(horizontal = 16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp),
            contentPadding = PaddingValues(vertical = 16.dp)
        ) {
            items(scenarios) { scenario ->
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    onClick = { navController.navigate("roleplay/${scenario.id}") }
                ) {
                    Row(Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                        Icon(scenario.icon, contentDescription = null, tint = MaterialTheme.colorScheme.primary, modifier = Modifier.size(32.dp))
                        Spacer(Modifier.width(16.dp))
                        Column {
                            Text(scenario.title, style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                            Text(scenario.description, style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                        }
                    }
                }
            }
        }
    }
}
