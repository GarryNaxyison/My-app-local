package ru.neriva.app.data

import android.content.Context
import android.graphics.BitmapFactory
import android.graphics.Canvas
import android.graphics.Color
import android.graphics.DashPathEffect
import android.graphics.Paint
import android.graphics.Typeface
import android.graphics.pdf.PdfDocument
import android.net.Uri
import android.text.Layout
import android.text.StaticLayout
import android.text.TextPaint
import ru.neriva.app.data.model.VocabularyItem
import ru.neriva.app.R

private data class CardLayout(val word: StaticLayout, val translation: StaticLayout)

private fun cardText(text: String, size: Float, bold: Boolean, width: Int): StaticLayout {
    val paint = TextPaint(Paint.ANTI_ALIAS_FLAG).apply {
        color = Color.rgb(25, 35, 59)
        textSize = size
        typeface = if (bold) Typeface.DEFAULT_BOLD else Typeface.DEFAULT
    }
    return StaticLayout.Builder.obtain(text, 0, text.length, paint, width)
        .setAlignment(Layout.Alignment.ALIGN_CENTER)
        .build()
}

internal fun writeWordCards(context: Context, uri: Uri, items: List<VocabularyItem>, fold: Boolean) {
    val columns = if (fold) 1 else 2
    val pageWidth = 595
    val pageHeight = 842
    val margin = 32
    val cardWidth = (pageWidth - 2 * margin) / columns
    val textWidth = cardWidth / (if (fold) 2 else 1) - 24
    val logo = BitmapFactory.decodeResource(context.resources, R.drawable.neriva_brand_symbol)
    val cards = items.filter { it.word.isNotBlank() }.map {
        CardLayout(cardText(it.word, 17f, true, textWidth), cardText(it.translation.orEmpty(), 13f, false, textWidth))
    }
    require(cards.isNotEmpty())
    val border = Paint(Paint.ANTI_ALIAS_FLAG).apply {
        color = Color.rgb(75, 85, 255)
        style = Paint.Style.STROKE
        strokeWidth = 1f
    }
    val divider = Paint(border).apply { pathEffect = DashPathEffect(floatArrayOf(4f, 4f), 0f) }
    val header = Paint(Paint.ANTI_ALIAS_FLAG).apply {
        color = Color.rgb(25, 35, 59)
        typeface = Typeface.DEFAULT_BOLD
        textSize = 20f
    }
    val top = 86
    val bottom = pageHeight - margin
    val pdf = PdfDocument()
    try {
        var pageNumber = 0
        fun newPage(): PdfDocument.Page {
            val page = pdf.startPage(PdfDocument.PageInfo.Builder(pageWidth, pageHeight, ++pageNumber).create())
            page.canvas.drawBitmap(logo, null, android.graphics.RectF(32f, 26f, 92f, 59.34f), null)
            page.canvas.drawText(context.getString(R.string.app_name), 104f, 50f, header)
            return page
        }

        var page = newPage()
        var y = top
        for (row in cards.chunked(columns)) {
            val height = row.maxOf { card ->
                val content = if (fold) maxOf(card.word.height, card.translation.height)
                    else card.word.height + 8 + card.translation.height
                maxOf(104, content + 28)
            }
            if (height > bottom - top) {
                pdf.finishPage(page)
                error("Card exceeds printable page height")
            }
            if (y + height > bottom) {
                pdf.finishPage(page)
                page = newPage()
                y = top
            }
            row.forEachIndexed { column, card ->
                val x = (margin + cardWidth * column).toFloat()
                val middle = x + cardWidth / 2f
                val canvas = page.canvas
                canvas.drawRect(x, y.toFloat(), x + cardWidth, (y + height).toFloat(), border)
                if (fold) {
                    canvas.drawLine(middle, y.toFloat(), middle, (y + height).toFloat(), divider)
                    canvas.drawCentered(card.word, x + cardWidth / 4f, y + (height - card.word.height) / 2f)
                    canvas.drawCentered(card.translation, x + cardWidth * 3f / 4f, y + (height - card.translation.height) / 2f)
                } else {
                    val textY = y + (height - card.word.height - 8 - card.translation.height) / 2f
                    canvas.drawCentered(card.word, middle, textY)
                    canvas.drawCentered(card.translation, middle, textY + card.word.height + 8)
                }
            }
            y += height
        }
        pdf.finishPage(page)
        requireNotNull(context.contentResolver.openOutputStream(uri, "wt")).use(pdf::writeTo)
    } finally {
        pdf.close()
    }
}

private fun Canvas.drawCentered(layout: StaticLayout, centerX: Float, top: Float) {
    save()
    translate(centerX - layout.width / 2f, top)
    layout.draw(this)
    restore()
}
