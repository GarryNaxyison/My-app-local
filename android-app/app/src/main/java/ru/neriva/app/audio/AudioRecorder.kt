package ru.neriva.app.audio

import android.content.Context
import android.media.AudioFormat
import android.media.AudioRecord
import android.media.MediaRecorder
import android.os.Environment
import java.io.File
import java.io.FileOutputStream
import java.io.RandomAccessFile

/**
 * Records mono 16-bit PCM from the microphone into a WAV file.
 *
 * Must be constructed with an Android [Context] so it can resolve an app-private
 * output directory without instantiating an Activity (which would crash).
 */
class AudioRecorder(private val context: Context) {
    private var audioRecord: AudioRecord? = null
    private var isRecording = false
    private var outputFile: File? = null

    companion object {
        private const val SAMPLE_RATE = 16000
        private const val CHANNEL_CONFIG = AudioFormat.CHANNEL_IN_MONO
        private const val AUDIO_FORMAT = AudioFormat.ENCODING_PCM_16BIT
    }

    fun startRecording(): File {
        val bufferSize = AudioRecord.getMinBufferSize(SAMPLE_RATE, CHANNEL_CONFIG, AUDIO_FORMAT)
        val outputDir = File(context.getExternalFilesDir(Environment.DIRECTORY_MUSIC), "recordings")
        outputDir.mkdirs()
        outputDir.mkdirs()
        val file = File(outputDir, "rec_${System.currentTimeMillis()}.wav")
        outputFile = file

        audioRecord = AudioRecord(MediaRecorder.AudioSource.MIC, SAMPLE_RATE, CHANNEL_CONFIG, AUDIO_FORMAT, bufferSize)
        audioRecord?.startRecording()
        isRecording = true

        Thread {
            val buffer = ByteArray(bufferSize)
            FileOutputStream(file).use { fos ->
                while (isRecording) {
                    val read = audioRecord?.read(buffer, 0, buffer.size) ?: 0
                    if (read > 0) fos.write(buffer, 0, read)
                }
            }
        }.start()

        return file
    }

    fun stopRecording(): File? {
        isRecording = false
        audioRecord?.stop()
        audioRecord?.release()
        audioRecord = null

        outputFile?.let { file ->
            // Write WAV header
            writeWavHeader(file)
        }
        return outputFile
    }

    fun isRecording(): Boolean = isRecording

    private fun writeWavHeader(file: File) {
        val fileSize = file.length()
        val dataSize = fileSize - 44
        RandomAccessFile(file, "rw").use { raf ->
            raf.seek(0)
            raf.writeBytes("RIFF")
            raf.write(intToLittleEndian((dataSize + 36).toInt()))
            raf.writeBytes("WAVE")
            raf.writeBytes("fmt ")
            raf.write(intToLittleEndian(16))
            raf.write(shortToLittleEndian(1)) // PCM
            raf.write(shortToLittleEndian(1)) // mono
            raf.write(intToLittleEndian(SAMPLE_RATE))
            raf.write(intToLittleEndian(SAMPLE_RATE * 2)) // byte rate
            raf.write(shortToLittleEndian(2)) // block align
            raf.write(shortToLittleEndian(16)) // bits per sample
            raf.writeBytes("data")
            raf.write(intToLittleEndian(dataSize.toInt()))
        }
    }

    private fun intToLittleEndian(value: Int): ByteArray {
        return byteArrayOf(
            (value and 0xFF).toByte(),
            (value shr 8 and 0xFF).toByte(),
            (value shr 16 and 0xFF).toByte(),
            (value shr 24 and 0xFF).toByte()
        )
    }

    private fun shortToLittleEndian(value: Int): ByteArray {
        return byteArrayOf(
            (value and 0xFF).toByte(),
            (value shr 8 and 0xFF).toByte()
        )
    }
}
