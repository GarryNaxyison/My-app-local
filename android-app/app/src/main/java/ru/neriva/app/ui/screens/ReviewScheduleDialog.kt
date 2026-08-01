package ru.neriva.app.ui.screens

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp

@Composable
fun ReviewScheduleDialog(
    onDismiss: () -> Unit,
    onSelect: (schedule: String) -> Unit
) {
    val options = listOf(
        "tomorrow" to "Tomorrow",
        "3_days" to "In 3 days",
        "1_week" to "In 1 week",
        "no_review" to "No review"
    )

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text(stringResource(R.string.tutor_review)) },
        text = {
            Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                Text(stringResource(R.string.tutor_srs))
                options.forEach { (value, label) ->
                    OutlinedButton(
                        onClick = { onSelect(value); onDismiss() },
                        modifier = Modifier.fillMaxWidth()
                    ) { Text(label, fontWeight = FontWeight.Medium) }
                }
            }
        },
        confirmButton = {},
        dismissButton = {
            TextButton(onClick = onDismiss) { Text(stringResource(R.string.cancel)) }
        }
    )
}
