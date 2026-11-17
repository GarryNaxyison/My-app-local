import type { VocabularyItem } from "./types";

export type WordCardLayout = "fold" | "cut";

const width = 1240;
const height = 1754;
const margin = 64;
const top = 180;
const gap = 20;
const cardWidth = (width - margin * 2 - gap) / 2;
const cardHeight = 460;
const ink = "#19233b";
const accent = "#4b55ff";

function contextForSheet() {
  const canvas = document.createElement("canvas");
  canvas.width = width;
  canvas.height = height;
  return canvas.getContext("2d")!;
}

function textLines(context: CanvasRenderingContext2D, text: string, maxWidth: number, font: string) {
  context.font = font;
  const lines: string[] = [];
  let line = "";
  for (const word of text.trim().split(/\s+/u)) {
    const candidate = line ? `${line} ${word}` : word;
    if (context.measureText(candidate).width <= maxWidth) {
      line = candidate;
      continue;
    }
    if (line) lines.push(line);
    line = "";
    for (const character of word) {
      if (line && context.measureText(line + character).width > maxWidth) {
        lines.push(line);
        line = "";
      }
      line += character;
    }
  }
  if (line) lines.push(line);
  return lines;
}

function prepareCard(context: CanvasRenderingContext2D, item: VocabularyItem, columns: number, layout: WordCardLayout) {
  const panelWidth = (columns === 2 ? cardWidth : width - margin * 2) / (layout === "fold" ? 2 : 1);
  const usableWidth = panelWidth - 56;
  const word = textLines(context, item.word, usableWidth, "bold 43px sans-serif");
  const answer = textLines(context, item.translation ?? "", usableWidth, "bold 32px sans-serif");
  const example = textLines(context, item.example || item.context || "", usableWidth, "27px sans-serif");
  const answerHeight = answer.length * 40 + example.length * 34 + (example.length ? 20 : 0);
  const wordHeight = word.length * 53;
  const contentHeight = layout === "fold" ? Math.max(wordHeight, answerHeight) : wordHeight + answerHeight + 20;
  return { word, answer, example, height: contentHeight + 128 };
}

function drawLines(context: CanvasRenderingContext2D, lines: string[], x: number, y: number, font: string, spacing: number) {
  context.font = font;
  lines.forEach((line, index) => context.fillText(line, x, y + index * spacing));
}

function drawCard(context: CanvasRenderingContext2D, card: ReturnType<typeof prepareCard>, x: number, y: number, w: number, h: number, layout: WordCardLayout) {
  context.strokeStyle = ink;
  context.lineWidth = 2;
  context.strokeRect(x, y, w, h);
  if (layout === "fold") {
    context.beginPath();
    context.setLineDash([10, 8]);
    context.moveTo(x + w / 2, y);
    context.lineTo(x + w / 2, y + h);
    context.stroke();
    context.setLineDash([]);
  }
  context.fillStyle = accent;
  context.fillRect(x + 28, y + 30, 30, 5);
  const answerX = layout === "fold" ? x + w / 2 + 28 : x + 28;
  const answerY = layout === "fold" ? y + 72 : y + 92 + card.word.length * 53;
  context.fillStyle = ink;
  drawLines(context, card.word, x + 28, y + 72, "bold 43px sans-serif", 53);
  drawLines(context, card.answer, answerX, answerY, "bold 32px sans-serif", 40);
  drawLines(context, card.example, answerX, answerY + card.answer.length * 40 + 20, "27px sans-serif", 34);
}

export function renderWordCardSheets(items: VocabularyItem[], layout: WordCardLayout, instructions: string, logo: HTMLImageElement | null) {
  const measure = contextForSheet();
  const sheets: string[] = [];
  const oversized: string[] = [];
  for (let offset = 0; offset < items.length;) {
    const pageStart = offset;
    const context = contextForSheet();
    context.fillStyle = "#fff";
    context.fillRect(0, 0, width, height);
    context.textBaseline = "top";
    context.fillStyle = ink;
    if (logo) {
      const scale = Math.min(210 / logo.naturalWidth, 60 / logo.naturalHeight);
      context.drawImage(logo, margin, 50, logo.naturalWidth * scale, logo.naturalHeight * scale);
    } else {
      drawLines(context, ["NERIVA"], margin, 59, "bold 49px sans-serif", 50);
    }
    const label = `${instructions}   ·   ${sheets.length + 1}`;
    drawLines(context, [label], margin, 120, "27px sans-serif", 34);
    let y = top;
    for (let row = 0; row < 3 && offset < items.length; row++) {
      const pair = items.slice(offset, offset + 2).map((item) => prepareCard(measure, item, 2, layout));
      if (pair.every((card) => card.height <= cardHeight)) {
        pair.forEach((card, column) => drawCard(context, card, margin + column * (cardWidth + gap), y, cardWidth, cardHeight, layout));
        offset += pair.length;
        y += cardHeight + gap;
        continue;
      }
      if (row > 0) break;
      const card = prepareCard(measure, items[offset], 1, layout);
      const fullHeight = height - top - margin;
      if (card.height > fullHeight) {
        oversized.push(items[offset].word);
        offset++;
        break;
      }
      drawCard(context, card, margin, y, width - margin * 2, Math.max(cardHeight, card.height), layout);
      y += cardHeight;
      offset++;
      break;
    }
    if (offset > pageStart && y > top) {
      sheets.push(context.canvas.toDataURL("image/png"));
    }
  }
  return { sheets, oversized };
}
