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
        RoleplayScenario("restaurant", "Restaurant", "Order food, ask about menu, handle complaints", Icons.Default.Restaurant),
        RoleplayScenario("work", "Work", "Office meetings, presentations, small talk with colleagues", Icons.Default.Work),
        RoleplayScenario("travel", "Travel", "Airport, hotel, asking for directions, sightseeing", Icons.Default.Flight),
        RoleplayScenario("exam", "Exam Prep", "IELTS/TOEFL speaking practice, timed responses", Icons.Default.Quiz),
        RoleplayScenario("small-talk", "Small Talk", "Casual conversations, weather, hobbies, weekend plans", Icons.Default.Chat),
        RoleplayScenario("hotel", "Hotel", "Check-in, room service, complaints, checkout", Icons.Default.Hotel),
        RoleplayScenario("shopping", "Shopping", "Ask for sizes, bargain, returns, compare products", Icons.Default.ShoppingBag),
        RoleplayScenario("doctor", "Doctor", "Describe symptoms, understand prescriptions, follow-up", Icons.Default.LocalHospital),
        RoleplayScenario("job-interview", "Job Interview", "Introduce yourself, describe experience, salary negotiation", Icons.Default.BusinessCenter),
        RoleplayScenario("bank", "Bank", "Open account, transfer money, explain issues", Icons.Default.AccountBalance),
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
