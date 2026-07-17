package ru.neriva.app.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch

@Composable
fun WordReportDialog(
    onDismiss: () -> Unit,
    onSubmit: (word: String, translation: String, comment: String) -> Unit
) {
    var word by remember { mutableStateOf("") }
    var translation by remember { mutableStateOf("") }
    var comment by remember { mutableStateOf("") }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Report a word mistake") },
        text = {
            Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                Text("Send the corrected word and translation.")
                OutlinedTextField(
                    value = word, onValueChange = { word = it },
                    label = { Text("Correct word") },
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = translation, onValueChange = { translation = it },
                    label = { Text("Correct translation") },
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = comment, onValueChange = { comment = it },
                    label = { Text("Comment (optional)") },
                    modifier = Modifier.fillMaxWidth()
                )
            }
        },
        confirmButton = {
            Button(
                onClick = { onSubmit(word, translation, comment); onDismiss() },
                enabled = word.isNotBlank() && translation.isNotBlank()
            ) { Text("Send report") }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text("Cancel") }
        }
    )
}
