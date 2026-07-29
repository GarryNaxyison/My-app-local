package ru.neriva.app.data.repo

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.runBlocking
import org.junit.Assert.assertEquals
import org.junit.Test
import ru.neriva.app.data.api.OkResponse
import ru.neriva.app.data.api.PhraseSaveRequest
import ru.neriva.app.data.api.PhrasebookResponse
import ru.neriva.app.data.db.PhrasebookDao
import ru.neriva.app.data.db.PhrasebookEntity
import ru.neriva.app.data.model.PhrasebookItem

class PhrasebookRepositoryTest {
    @Test
    fun refreshReplacesCachedItemsWithServerItems() = runBlocking {
        val dao = FakePhrasebookDao(
            PhrasebookEntity("stale", "Stale", null, null, null, null, null)
        )
        val remote = FakePhrasebookRemote(
            items = listOf(PhrasebookItem(id = "server-1", phrase = "Where is the station?", translation = "Где станция?"))
        )
        val repository = PhrasebookRepository(remote, dao)

        repository.refresh()

        assertEquals(listOf("server-1"), dao.getAll().map { it.id })
        assertEquals("Где станция?", dao.getAll().single().translation)
    }

    @Test
    fun savePostsPhraseThenRefreshesCachedItems() = runBlocking {
        val dao = FakePhrasebookDao()
        val remote = FakePhrasebookRemote(
            items = listOf(PhrasebookItem(id = "server-2", phrase = "Thank you", source = "manual"))
        )
        val repository = PhrasebookRepository(remote, dao)

        repository.save(phrase = "Thank you", translation = "Спасибо", note = "Polite", source = "manual", language = "en")

        assertEquals(
            PhraseSaveRequest("Thank you", "Спасибо", "Polite", "manual", "en"),
            remote.savedRequest
        )
        assertEquals(listOf("server-2"), dao.getAll().map { it.id })
    }

    @Test
    fun deleteSendsServerIdThenRefreshesCachedItems() = runBlocking {
        val dao = FakePhrasebookDao(PhrasebookEntity("old", "Old", null, null, null, null, null))
        val remote = FakePhrasebookRemote(items = emptyList())
        val repository = PhrasebookRepository(remote, dao)

        repository.delete("old")

        assertEquals("old", remote.deletedId)
        assertEquals(emptyList<PhrasebookEntity>(), dao.getAll())
    }
}

private class FakePhrasebookRemote(
    var items: List<PhrasebookItem>,
) : PhrasebookRemoteDataSource {
    var savedRequest: PhraseSaveRequest? = null
    var deletedId: String? = null

    override suspend fun getPhrasebook(): PhrasebookResponse = PhrasebookResponse(items)

    override suspend fun savePhrase(request: PhraseSaveRequest): OkResponse {
        savedRequest = request
        return OkResponse()
    }

    override suspend fun deletePhrase(id: String): OkResponse {
        deletedId = id
        return OkResponse()
    }
}

private class FakePhrasebookDao(
    vararg initial: PhrasebookEntity,
) : PhrasebookDao {
    private val stored = initial.toMutableList()
    private val changes = MutableStateFlow(stored.toList())

    override fun observeAll(): Flow<List<PhrasebookEntity>> = changes

    override suspend fun getAll(): List<PhrasebookEntity> = stored.toList()

    override suspend fun insert(item: PhrasebookEntity) {
        stored.removeAll { it.id == item.id }
        stored += item
        changes.value = stored.toList()
    }

    override suspend fun insertAll(items: List<PhrasebookEntity>) {
        for (item in items) insert(item)
    }

    override suspend fun delete(id: String) {
        stored.removeAll { it.id == id }
        changes.value = stored.toList()
    }

    override suspend fun clear() {
        stored.clear()
        changes.value = emptyList()
    }
}
