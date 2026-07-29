package ru.neriva.app.data.db

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import androidx.room.Transaction
import kotlinx.coroutines.flow.Flow

@Dao
interface PhrasebookDao {
    @Query("SELECT * FROM phrasebook ORDER BY createdAt DESC")
    fun observeAll(): Flow<List<PhrasebookEntity>>

    @Query("SELECT * FROM phrasebook ORDER BY createdAt DESC")
    suspend fun getAll(): List<PhrasebookEntity>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(item: PhrasebookEntity)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(items: List<PhrasebookEntity>)

    @Query("DELETE FROM phrasebook WHERE id = :id")
    suspend fun delete(id: String)

    @Query("DELETE FROM phrasebook")
    suspend fun clear()

    @Transaction
    suspend fun replaceAll(items: List<PhrasebookEntity>) {
        clear()
        insertAll(items)
    }
}
