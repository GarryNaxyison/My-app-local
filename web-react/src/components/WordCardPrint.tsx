import { useEffect, useMemo, useState } from "react";
import type { VocabularyItem } from "../lib/types";
import { renderWordCardSheets, type WordCardLayout } from "../lib/wordCardSheets";
import "./word-card-print.css";

interface Props {
  items: VocabularyItem[];
  russian: boolean;
}

function printSheets(images: string[]) {
  const frame = document.createElement("iframe");
  frame.title = "NERIVA";
  frame.className = "word-card-print__frame";
  frame.onload = () => {
    const target = frame.contentWindow;
    if (!target) {
      frame.remove();
      return;
    }
    target.addEventListener("afterprint", () => frame.remove(), { once: true });
    target.print();
  };
  const sheets = images.map((image) => `<img src="${image}" alt="">`).join("");
  frame.srcdoc = `<!doctype html><meta charset="utf-8"><title>NERIVA A4</title><style>
    @page { size: A4; margin: 0; }
    body { margin: 0; }
    img { display: block; width: 209.5mm; height: auto; margin-inline: auto; break-inside: avoid; }
    img + img { break-before: page; }
  </style>${sheets}`;
  document.body.append(frame);
}

export function WordCardPrint({ items, russian }: Props) {
  const [open, setOpen] = useState(false);
  const [layout, setLayout] = useState<WordCardLayout>("fold");
  const [omitted, setOmitted] = useState<Set<number>>(() => new Set());
  const [logo, setLogo] = useState<HTMLImageElement | null | undefined>(undefined);
  useEffect(() => {
    if (!open || logo !== undefined) return;
    const image = new Image();
    image.onload = () => setLogo(image);
    image.onerror = () => setLogo(null);
    image.src = "/app/assets/neriva-wordmark.png";
    return () => {
      image.onload = null;
      image.onerror = null;
    };
  }, [open, logo]);
  const fold = russian ? "Разрезать и сложить" : "Cut and fold";
  const cut = russian ? "Вырезать" : "Cut out";
  const { sheets: images, error, oversized } = useMemo(() => {
    if (!open || logo === undefined) return { sheets: [], error: false, oversized: [] };
    try {
      const instructions = layout === "fold" ? fold : cut;
      return { ...renderWordCardSheets(items.filter((_, index) => !omitted.has(index)), layout, instructions, logo), error: false };
    } catch {
      return { sheets: [], error: true, oversized: [] };
    }
  }, [items, omitted, open, layout, fold, cut, logo]);

  function toggleCard(index: number) {
    setOmitted((current) => {
      const next = new Set(current);
      if (next.has(index)) next.delete(index);
      else next.add(index);
      return next;
    });
  }

  return (
    <div className="word-card-print">
      <button type="button" aria-expanded={open} onClick={() => setOpen(!open)}>
        {russian ? "Печатные карточки" : "Print word cards"}
      </button>
      {open && (
        <>
          <div className="word-card-print__controls">
            <label>
              {russian ? "Раскладка" : "Layout"}
              <select value={layout} onChange={(event) => setLayout(event.target.value === "cut" ? "cut" : "fold")}>
                <option value="fold">{fold}</option>
                <option value="cut">{cut}</option>
              </select>
            </label>
            <button type="button" disabled={!images.length} onClick={() => printSheets(images)}>
              {russian ? "Печать / PDF" : "Print / PDF"}
            </button>
          </div>
          <p>
            {layout === "fold"
              ? russian ? "Режьте по сплошной линии, сгибайте по пунктиру печатной стороной наружу." : "Cut on the solid line and fold on the dashed line with the printed sides facing out."
              : russian ? "Режьте по сплошной линии." : "Cut on the solid line."}
          </p>
          <details>
            <summary>{russian ? "Выбрать слова с этой страницы" : "Choose words on this page"}</summary>
            <ul className="word-card-print__selection">
              {items.map((item, index) => (
                <li key={index}>
                  <label>
                    <input type="checkbox" checked={!omitted.has(index)} onChange={() => toggleCard(index)} />
                    <span dir="auto">{item.word}{item.translation && <> · {item.translation}</>}</span>
                  </label>
                </li>
              ))}
            </ul>
          </details>
          {oversized.length > 0 && <p role="status">{russian ? "Не поместились даже на отдельном листе" : "Too long even on a separate sheet"} — {oversized.join(", ")}</p>}
          {error && <p role="alert">{russian ? "Не удалось разложить набор. Снимите отметку с самых длинных карточек и повторите." : "The set could not be laid out. Deselect the longest cards and try again."}</p>}
          <div className="word-card-print__sheets">
            {images.map((src, index) => (
              <figure className="word-card-print__sheet" key={index}>
                <img data-word-card-sheet src={src} width="1240" height="1754" alt={`NERIVA A4 · ${index + 1}`} />
                <figcaption>
                  <a href={src} download={`NERIVA-A4-${index + 1}.png`}>PNG {index + 1}</a>
                </figcaption>
              </figure>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
