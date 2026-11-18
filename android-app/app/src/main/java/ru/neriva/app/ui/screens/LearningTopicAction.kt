package ru.neriva.app.ui.screens

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ExperimentalLayoutApi
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.MenuBook
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.launch
import ru.neriva.app.R

@OptIn(ExperimentalLayoutApi::class)
@Composable
internal fun LearningTopicAction(
    current: String,
    save: suspend (String) -> Unit,
    iconOnly: Boolean = true,
) {
    val title = stringResource(R.string.learning_topic_title)
    var open by rememberSaveable { mutableStateOf(false) }
    var draft by rememberSaveable { mutableStateOf("") }
    var saving by remember { mutableStateOf(false) }
    var failed by remember { mutableStateOf(false) }
    val scope = rememberCoroutineScope()
    val show = {
        draft = current
        failed = false
        open = true
    }
    if (iconOnly) {
        IconButton(onClick = show) { Icon(Icons.AutoMirrored.Filled.MenuBook, title) }
    } else {
        OutlinedButton(onClick = show, modifier = Modifier.fillMaxWidth()) { Text(title) }
    }
    if (!open) return

    AlertDialog(
        onDismissRequest = { if (!saving) open = false },
        title = { Text(title) },
        text = {
            Column(
                modifier = Modifier.verticalScroll(rememberScrollState()),
                verticalArrangement = Arrangement.spacedBy(8.dp),
            ) {
                FlowRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    FilterChip(draft.isBlank(), { draft = "" }, {
                        Text(stringResource(R.string.learning_topic_mixed))
                    }, enabled = !saving)
                    listOf(
                        R.string.roleplay_scenario_restaurant,
                        R.string.roleplay_scenario_work,
                        R.string.roleplay_scenario_travel,
                        R.string.roleplay_scenario_exam,
                    ).forEach { preset ->
                        val label = stringResource(preset)
                        FilterChip(draft.trim() == label, { draft = label }, { Text(label) }, enabled = !saving)
                    }
                }
                OutlinedTextField(
                    value = draft,
                    onValueChange = { draft = it },
                    label = { Text(stringResource(R.string.learning_topic_custom)) },
                    enabled = !saving,
                    modifier = Modifier.fillMaxWidth(),
                    maxLines = 3,
                )
                if (saving) LinearProgressIndicator(Modifier.fillMaxWidth())
                if (failed) Text(stringResource(R.string.request_failed), color = MaterialTheme.colorScheme.error)
            }
        },
        confirmButton = {
            TextButton(onClick = {
                val next = draft.trim()
                if (next == current.trim()) {
                    open = false
                } else {
                    saving = true
                    failed = false
                    scope.launch {
                        try {
                            save(next)
                            open = false
                        } catch (e: CancellationException) {
                            throw e
                        } catch (_: Exception) {
                            failed = true
                        } finally {
                            saving = false
                        }
                    }
                }
            }, enabled = !saving) { Text(stringResource(R.string.learning_topic_save)) }
        },
        dismissButton = {
            TextButton(onClick = { open = false }, enabled = !saving) { Text(stringResource(R.string.cancel)) }
        },
    )
}
