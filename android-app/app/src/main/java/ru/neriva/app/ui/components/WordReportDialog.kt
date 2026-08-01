package ru.neriva.app.ui.components

import ru.neriva.app.R

import androidx.compose.ui.res.stringResource

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ReportProblem
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp

/**
 * Dialog for reporting a wrong word/translation to the review flow,
 * mirroring the web app's AI Tutor word-report button.
 */
@Composable
fun WordReportDialog(
    word: String,
    translation: String,
    onDismiss: () -> Unit,
    onSubmit: (comment: String) -> Unit,
    pending: Boolean = false,
    error: String? = null,
) {
    var comment by remember { mutableStateOf("") }
    AlertDialog(
        onDismissRequest = { if (!pending) onDismiss() },
        icon = { Icon(Icons.Default.ReportProblem, contentDescription = null, tint = MaterialTheme.colorScheme.error) },
        title = { Text(stringResource(R.string.ai_tutor_word_report_button), fontWeight = FontWeight.SemiBold) },
        text = {
            Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                Text(word, style = MaterialTheme.typography.titleSmall, fontWeight = FontWeight.SemiBold)
                if (translation.isNotBlank()) {
                    Text(translation, style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant)
                }
                OutlinedTextField(
                    value = comment,
                    onValueChange = { comment = it },
                    label = { Text(stringResource(R.string.ai_tutor_word_report_comment_label)) },
                    modifier = Modifier.fillMaxWidth(),
                    minLines = 2,
                )
                error?.let { Text(it, color = MaterialTheme.colorScheme.error, style = MaterialTheme.typography.bodySmall) }
            }
        },
        confirmButton = {
            Button(
                onClick = { onSubmit(comment) },
                enabled = !pending,
                colors = ButtonDefaults.buttonColors(containerColor = MaterialTheme.colorScheme.error),
            ) {
                if (pending) {
                    CircularProgressIndicator(modifier = Modifier.padding(end = 8.dp), strokeWidth = 2.dp, color = MaterialTheme.colorScheme.onError)
                }
                Text(stringResource(R.string.ai_tutor_word_report_submit))
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss, enabled = !pending) { Text(stringResource(R.string.cancel)) }
        },
    )
}
