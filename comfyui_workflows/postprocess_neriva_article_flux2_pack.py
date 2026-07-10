from __future__ import annotations

import json
from pathlib import Path
from textwrap import wrap

from PIL import Image, ImageDraw, ImageFont


ROOT = Path(__file__).resolve().parents[1]
PACK_DIR = ROOT / "comfyui_workflows" / "generated" / "neriva_article_flux2_dev_v1"
FONT_PATH = ROOT / "comfyui_workflows" / "assets" / "fonts" / "Manrope-wght.ttf"


def font(size: int) -> ImageFont.FreeTypeFont:
    return ImageFont.truetype(str(FONT_PATH), size=size)


def composite_panel(
    image: Image.Image,
    box: tuple[int, int, int, int],
    *,
    fill: tuple[int, int, int, int] = (8, 18, 28, 188),
    outline: tuple[int, int, int, int] = (150, 235, 255, 155),
    radius: int = 28,
) -> ImageDraw.ImageDraw:
    overlay = Image.new("RGBA", image.size, (0, 0, 0, 0))
    draw = ImageDraw.Draw(overlay)
    draw.rounded_rectangle(box, radius=radius, fill=fill, outline=outline, width=2)
    image.alpha_composite(overlay)
    return ImageDraw.Draw(image)


def draw_text(draw: ImageDraw.ImageDraw, xy: tuple[int, int], text: str, size: int, fill=(245, 252, 255, 255), weight=False) -> None:
    draw.text(xy, text, font=font(size), fill=fill)


def draw_wrapped(
    draw: ImageDraw.ImageDraw,
    xy: tuple[int, int],
    text: str,
    size: int,
    max_chars: int,
    fill=(198, 221, 229, 255),
    line_gap: int = 8,
) -> int:
    x, y = xy
    for line in wrap(text, width=max_chars):
        draw.text((x, y), line, font=font(size), fill=fill)
        y += size + line_gap
    return y


def decorate_translation(path: Path) -> None:
    image = Image.open(path).convert("RGBA")
    w, h = image.size
    x1, y1, x2, y2 = int(w * 0.398), int(h * 0.205), int(w * 0.616), int(h * 0.690)
    draw = composite_panel(image, (x1, y1, x2, y2), fill=(246, 253, 255, 255), outline=(112, 210, 255, 185), radius=22)

    pad = int(w * 0.017)
    x = x1 + pad
    y = y1 + int(h * 0.034)
    draw_text(draw, (x, y), "RU", 20, (28, 122, 170, 255))
    y += 30
    draw_text(draw, (x, y), "Таким Запомни", 34, (12, 20, 28, 255))
    y += 58
    draw.line((x, y, x2 - pad, y), fill=(76, 175, 220, 180), width=2)
    y += 24
    draw_text(draw, (x, y), "EN", 20, (190, 116, 34, 255))
    y += 30
    draw_text(draw, (x, y), "Remember this", 30, (12, 20, 28, 255))
    y += 34
    draw_text(draw, (x, y), "moment", 30, (12, 20, 28, 255))
    y += 54

    chip_h = 34
    for label, color in [
        ("Emotional interpretation", (44, 205, 166, 230)),
        ("Metaphor check", (224, 151, 52, 230)),
    ]:
        draw.rounded_rectangle((x, y, x2 - pad, y + chip_h), radius=14, fill=(219, 238, 244, 235), outline=color, width=2)
        draw_text(draw, (x + 10, y + 6), label, 16, (13, 37, 49, 255))
        y += chip_h + 12

    image.convert("RGB").save(path, quality=96)


def decorate_music(path: Path) -> None:
    image = Image.open(path).convert("RGBA")
    w, h = image.size
    x1, y1, x2, y2 = int(w * 0.08), int(h * 0.12), int(w * 0.48), int(h * 0.62)
    draw = composite_panel(image, (x1, y1, x2, y2), fill=(10, 12, 18, 198), outline=(120, 214, 255, 150))
    x = x1 + 34
    y = y1 + 34
    draw_text(draw, (x, y), "Lyrics Analysis", 42)
    y += 52
    draw_text(draw, (x, y), "Thus Remember", 30, (255, 199, 121, 255))
    y += 58
    y = draw_wrapped(draw, (x, y), "Line focus: 'Таким Запомни'", 24, 33)
    y += 12
    y = draw_wrapped(draw, (x, y), "Context: emotional memory, direct address, metaphor check", 22, 38, fill=(196, 226, 236, 255))
    y += 14
    draw.line((x, y, x2 - 34, y), fill=(90, 200, 255, 115), width=2)
    y += 24
    draw_text(draw, (x, y), "slang: clean", 22, (118, 255, 196, 255))
    draw_text(draw, (x + 190, y), "tone: cinematic", 22, (118, 255, 196, 255))
    image.convert("RGB").save(path, quality=96)


def decorate_top_list(path: Path) -> None:
    image = Image.open(path).convert("RGBA")
    w, h = image.size
    x1, y1, x2, y2 = int(w * 0.18), int(h * 0.10), int(w * 0.84), int(h * 0.86)
    draw = composite_panel(image, (x1, y1, x2, y2), fill=(5, 13, 25, 232), outline=(95, 205, 255, 175))
    x = x1 + 46
    y = y1 + 38
    draw.rounded_rectangle((x - 18, y - 18, x2 - 46, y + 92), radius=20, fill=(3, 9, 17, 242), outline=(75, 175, 220, 105), width=1)
    draw_text(draw, (x, y), "ТОП-5 ИИ ДЛЯ ЯЗЫКОВ 2026", 42)
    y += 54
    draw_text(draw, (x, y), "BEST LANGUAGE BOTS 2026", 26, (159, 211, 234, 255))
    y += 55
    rows = [
        ("1", "Neriva AI", (70, 255, 164, 230), True),
        ("2", "Andy", (31, 52, 70, 205), False),
        ("3", "Translate", (31, 52, 70, 205), False),
        ("4", "Lingvo", (31, 52, 70, 205), False),
        ("5", "Voice Tutor", (31, 52, 70, 205), False),
    ]
    row_h = 72
    for rank, name, fill, highlight in rows:
        outline = (87, 255, 174, 205) if highlight else (91, 158, 190, 95)
        draw.rounded_rectangle((x, y, x2 - 46, y + row_h), radius=20, fill=fill, outline=outline, width=2)
        draw_text(draw, (x + 22, y + 18), rank, 30, (245, 252, 255, 255))
        draw_text(draw, (x + 86, y + 17), name, 32, (2, 18, 16, 255) if highlight else (231, 244, 250, 255))
        if highlight:
            draw_text(draw, (x2 - 178, y + 21), "BEST", 24, (2, 34, 24, 255))
        y += row_h + 16
    image.convert("RGB").save(path, quality=96)


def make_contact_sheet(paths: list[Path]) -> None:
    thumbs = []
    for path in paths:
        img = Image.open(path).convert("RGB")
        img.thumbnail((480, 270), Image.Resampling.LANCZOS)
        thumbs.append((path, img.copy()))

    margin = 28
    label_h = 48
    cell_w, cell_h = 480, 270 + label_h
    cols = 2
    rows = (len(thumbs) + cols - 1) // cols
    sheet = Image.new("RGB", (cols * cell_w + (cols + 1) * margin, rows * cell_h + (rows + 1) * margin), (10, 13, 18))
    draw = ImageDraw.Draw(sheet)
    for i, (path, img) in enumerate(thumbs):
        col = i % cols
        row = i // cols
        x = margin + col * (cell_w + margin)
        y = margin + row * (cell_h + margin)
        sheet.paste(img, (x, y))
        draw.text((x, y + 282), path.name, font=font(20), fill=(232, 240, 244))
    sheet.save(PACK_DIR / "neriva_article_flux2_dev_v1_contact_sheet.png", quality=94)


def main() -> None:
    manifest_path = PACK_DIR / "manifest.json"
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    paths = [PACK_DIR / item["project_output"].split("neriva_article_flux2_dev_v1", 1)[1].lstrip("\\/") for item in manifest["results"]]

    decorators = {
        "neriva_02_translation_context_phone_flux2.png": decorate_translation,
        "neriva_04_creative_music_analysis_flux2.png": decorate_music,
        "neriva_05_top_language_bots_flux2.png": decorate_top_list,
    }
    for path in paths:
        decorator = decorators.get(path.name)
        if decorator:
            decorator(path)
    make_contact_sheet(paths)
    manifest["postprocessed"] = True
    manifest["contact_sheet"] = str(PACK_DIR / "neriva_article_flux2_dev_v1_contact_sheet.png")
    manifest_path.write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main()
