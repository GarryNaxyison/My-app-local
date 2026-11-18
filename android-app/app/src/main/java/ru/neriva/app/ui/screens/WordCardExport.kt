package ru.neriva.app.ui.screens

import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import ru.neriva.app.R
import ru.neriva.app.data.repo.VocabularyRepository
import ru.neriva.app.data.writeWordCards

@Composable
internal fun WordCardExport(repo: VocabularyRepository) {
    val context = LocalContext.current
    val scope = rememberCoroutineScope()
    var choosing by remember { mutableStateOf(false) }
    var fold by rememberSaveable { mutableStateOf(true) }
    var working by remember { mutableStateOf(false) }
    var failure by remember { mutableStateOf(false) }
    var saved by remember { mutableStateOf(false) }
    val exporter = rememberLauncherForActivityResult(ActivityResultContracts.CreateDocument("application/pdf")) { uri ->
        if (uri != null) {
            val foldLayout = fold
            working = true
            scope.launch {
                try {
                    val first = repo.getList(0)
                    val items = first.items.orEmpty().toMutableList()
                    for (page in 1 until (first.totalPages ?: 1)) {
                        items += repo.getList(page).items.orEmpty()
                    }
                    withContext(Dispatchers.IO) { writeWordCards(context, uri, items, foldLayout) }
                    saved = true
                } catch (e: CancellationException) {
                    throw e
                } catch (_: Exception) {
                    failure = true
                } finally {
                    working = false
                }
            }
        }
    }
    OutlinedButton(onClick = { choosing = true; saved = false; failure = false }, enabled = !working) {
        Text(stringResource(R.string.word_cards_export))
    }
    if (working) LinearProgressIndicator()
    if (saved) Text(stringResource(R.string.word_cards_saved))
    if (failure) Text(stringResource(R.string.word_cards_failed), color = MaterialTheme.colorScheme.error)
    if (choosing) {
        AlertDialog(
            onDismissRequest = { choosing = false },
            title = { Text(stringResource(R.string.word_cards_export)) },
            text = { Text(stringResource(R.string.word_cards_layout)) },
            confirmButton = {
                TextButton(onClick = {
                    fold = true
                    choosing = false
                    exporter.launch(context.getString(R.string.word_cards_fold_filename))
                }) { Text(stringResource(R.string.word_cards_fold)) }
            },
            dismissButton = {
                TextButton(onClick = {
                    fold = false
                    choosing = false
                    exporter.launch(context.getString(R.string.word_cards_cut_filename))
                }) { Text(stringResource(R.string.word_cards_cut)) }
            },
        )
    }
}
