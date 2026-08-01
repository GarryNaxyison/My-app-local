package ru.neriva.app.ui.screens

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

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
        title = { Text(stringResource(R.string.ai_tutor_word_report_title)) },
        text = {
            Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                Text(stringResource(R.string.ai_tutor_word_report_body))
                OutlinedTextField(
                    value = word, onValueChange = { word = it },
                    label = { Text(stringResource(R.string.ai_tutor_word_report_word_label)) },
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = translation, onValueChange = { translation = it },
                    label = { Text(stringResource(R.string.ai_tutor_word_report_translation_label)) },
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = comment, onValueChange = { comment = it },
                    label = { Text(stringResource(R.string.ai_tutor_word_report_comment_label)) },
                    modifier = Modifier.fillMaxWidth()
                )
            }
        },
        confirmButton = {
            Button(
                onClick = { onSubmit(word, translation, comment); onDismiss() },
                enabled = word.isNotBlank() && translation.isNotBlank()
            ) { Text(stringResource(R.string.ai_tutor_word_report_submit)) }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text(stringResource(R.string.cancel)) }
        }
    )
}
