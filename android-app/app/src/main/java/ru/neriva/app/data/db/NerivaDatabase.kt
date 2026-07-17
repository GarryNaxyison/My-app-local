package ru.neriva.app.data.db

import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase

@Database(entities = [PhrasebookEntity::class], version = 1, exportSchema = false)
abstract class NerivaDatabase : RoomDatabase() {
    abstract fun phrasebookDao(): PhrasebookDao

    companion object {
        @Volatile private var instance: NerivaDatabase? = null

        fun getInstance(context: Context): NerivaDatabase {
            return instance ?: synchronized(this) {
                instance ?: Room.databaseBuilder(
                    context.applicationContext,
                    NerivaDatabase::class.java,
                    "neriva.db"
                ).build().also { instance = it }
            }
        }
    }
}
