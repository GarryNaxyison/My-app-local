from __future__ import annotations

import json
import math
from pathlib import Path

from PIL import Image, ImageDraw, ImageFilter

from render_poliglot_tech_outro_v11 import (
    COMFY_OUTPUT,
    FPS,
    HEIGHT,
    MANROPE_FONT_PATH,
    PREVIEW_DIR,
    TOTAL_FRAMES,
    WIDTH,
    clamp,
    encode_video,
    smoothstep,
)


TEXT = "NERIVA.RU"
TEXT_PREFIX = "NERIVA."
TEXT_ACCENT = "RU"

FONT_CANDIDATES = (
    MANROPE_FONT_PATH,
    Path(r"C:\Windows\Fonts\segoeuib.ttf"),
    Path(r"C:\Windows\Fonts\arialbd.ttf"),
)

OUTPUTS = {
    "01_noir": COMFY_OUTPUT / "neriva_luxury_minimal_v18_01_noir.mp4",
    "02_ivory": COMFY_OUTPUT / "neriva_luxury_minimal_v18_02_ivory.mp4",
    "03_platinum": COMFY_OUTPUT / "neriva_luxury_minimal_v18_03_platinum.mp4",
    "04_gallery": COMFY_OUTPUT / "neriva_luxury_minimal_v18_04_gallery.mp4",
    "05_emerald": COMFY_OUTPUT / "neriva_luxury_minimal_v18_05_emerald.mp4",
}


def ease_out_cubic(value: float) -> float:
    x = clamp(value)
    return 1 - pow(1 - x, 3)


def ease_in_out(value: float) -> float:
    x = clamp(value)
    return x * x * (3 - 2 * x)


def load_clean_font(size: int):
    for path in FONT_CANDIDATES:
        if path.exists():
            font = ImageFont.truetype(str(path), size=size)
            if path == MANROPE_FONT_PATH:
                try:
                    font.set_variation_by_name("ExtraBold")
                except OSError:
                    pass
            return font, path
    raise FileNotFoundError("No TrueType font found for Neriva outro")


try:
    from PIL import ImageFont
except ImportError as exc:  # pragma: no cover
    raise RuntimeError("Pillow ImageFont is required") from exc


def apply_opacity(layer: Image.Image, alpha: float) -> Image.Image:
    if alpha >= 0.999:
        return layer
    out = layer.copy()
    out.putalpha(out.getchannel("A").point(lambda value: int(value * clamp(alpha))))
    return out


def add_grain(base: Image.Image, frame_index: int, alpha: int = 12) -> None:
    noise = Image.effect_noise((WIDTH, HEIGHT), 9.0 + (frame_index % 5) * 0.6).convert("L")
    grain = Image.merge("RGBA", (noise, noise, noise, noise.point(lambda value: alpha)))
    base.alpha_composite(grain)


def soft_ellipse(
    base: Image.Image,
    box: tuple[int, int, int, int],
    color: tuple[int, int, int, int],
    blur: int,
) -> None:
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    draw.ellipse(box, fill=color)
    base.alpha_composite(layer.filter(ImageFilter.GaussianBlur(blur)))


def text_width(draw: ImageDraw.ImageDraw, text: str, font, tracking: float) -> float:
    if not text:
        return 0.0
    return sum(draw.textlength(char, font=font) for char in text) + tracking * (len(text) - 1)


def text_bbox_tracked(font, tracking: float) -> tuple[int, int, int, int]:
    probe = Image.new("RGBA", (1, 1), (0, 0, 0, 0))
    draw = ImageDraw.Draw(probe)
    bbox = draw.textbbox((0, 0), TEXT, font=font, stroke_width=0)
    width = int(math.ceil(text_width(draw, TEXT, font, tracking)))
    return 0, bbox[1], width, bbox[3]


def draw_tracked_text(
    draw: ImageDraw.ImageDraw,
    xy: tuple[float, float],
    font,
    tracking: float,
    fill: tuple[int, int, int, int],
    accent: tuple[int, int, int, int] | None = None,
    stroke_width: int = 0,
    stroke_fill: tuple[int, int, int, int] | None = None,
) -> None:
    x, y = xy
    accent_start = len(TEXT_PREFIX)
    accent_end = accent_start + len(TEXT_ACCENT)
    for index, char in enumerate(TEXT):
        char_fill = accent if accent is not None and accent_start <= index < accent_end else fill
        draw.text((x, y), char, font=font, fill=char_fill, stroke_width=stroke_width, stroke_fill=stroke_fill)
        x += draw.textlength(char, font=font) + tracking


def make_text_layer(
    *,
    font_size: int,
    tracking: float,
    fill: tuple[int, int, int, int],
    accent: tuple[int, int, int, int] | None,
    y: int,
    stroke_width: int = 0,
    stroke_fill: tuple[int, int, int, int] | None = None,
    shadow: tuple[int, int, tuple[int, int, int, int], int] | None = None,
) -> tuple[Image.Image, tuple[int, int, int, int], Path]:
    font, font_path = load_clean_font(font_size)
    bbox = text_bbox_tracked(font, tracking)
    width = bbox[2] - bbox[0]
    height = bbox[3] - bbox[1]
    x = int((WIDTH - width) / 2)
    top = y
    draw_y = top - bbox[1]

    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    if shadow is not None:
        dx, dy, shadow_fill, blur = shadow
        shadow_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        shadow_draw = ImageDraw.Draw(shadow_layer)
        draw_tracked_text(shadow_draw, (x + dx, draw_y + dy), font, tracking, shadow_fill, shadow_fill, stroke_width, shadow_fill)
        layer.alpha_composite(shadow_layer.filter(ImageFilter.GaussianBlur(blur)))
    draw_tracked_text(draw, (x, draw_y), font, tracking, fill, accent, stroke_width, stroke_fill)
    return layer, (x, top, x + width, top + height), font_path


def reveal_text(
    base: Image.Image,
    text_layer: Image.Image,
    box: tuple[int, int, int, int],
    time: float,
    *,
    start: float,
    end: float,
    mode: str,
    y_shift: int = 24,
) -> None:
    progress = smoothstep(start, end, time)
    if progress <= 0:
        return

    layer = text_layer.copy()
    if mode == "wipe":
        mask = Image.new("L", (WIDTH, HEIGHT), 0)
        draw = ImageDraw.Draw(mask)
        x0, y0, x1, y1 = box
        pad = 42
        reveal_x = x0 - pad + int((x1 - x0 + pad * 2) * ease_out_cubic(progress))
        draw.rectangle((x0 - pad, y0 - pad, reveal_x, y1 + pad), fill=255)
        layer.putalpha(Image.composite(layer.getchannel("A"), Image.new("L", (WIDTH, HEIGHT), 0), mask))
    elif mode == "center":
        mask = Image.new("L", (WIDTH, HEIGHT), 0)
        draw = ImageDraw.Draw(mask)
        x0, y0, x1, y1 = box
        cx = (x0 + x1) // 2
        half = int((x1 - x0 + 80) * ease_out_cubic(progress) / 2)
        draw.rectangle((cx - half, y0 - 48, cx + half, y1 + 48), fill=255)
        layer.putalpha(Image.composite(layer.getchannel("A"), Image.new("L", (WIDTH, HEIGHT), 0), mask))
    elif mode == "soft":
        if progress < 0.76:
            layer = layer.filter(ImageFilter.GaussianBlur(1.8 * (1 - progress / 0.76)))

    if mode in {"wipe", "center", "soft"}:
        alpha = progress
    else:
        alpha = progress

    offset = int((1 - ease_out_cubic(progress)) * y_shift)
    moved = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    moved.alpha_composite(layer, (0, offset))
    base.alpha_composite(apply_opacity(moved, alpha))


def add_light_sweep(
    base: Image.Image,
    box: tuple[int, int, int, int],
    time: float,
    *,
    start: float = 1.18,
    end: float = 1.78,
    color: tuple[int, int, int] = (255, 255, 255),
    strength: int = 70,
) -> None:
    progress = smoothstep(start, end, time)
    if progress <= 0 or progress >= 1:
        return
    x0, y0, x1, y1 = box
    w = x1 - x0
    x = x0 - 150 + int((w + 300) * ease_in_out(progress))
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    alpha = int(strength * math.sin(math.pi * progress))
    draw.polygon(
        [(x - 36, y0 - 34), (x + 22, y0 - 34), (x + 88, y1 + 46), (x + 30, y1 + 46)],
        fill=(color[0], color[1], color[2], alpha),
    )
    base.alpha_composite(layer.filter(ImageFilter.GaussianBlur(1.0)))


def draw_monogram_n(
    base: Image.Image,
    *,
    cx: int,
    cy: int,
    size: int,
    color: tuple[int, int, int, int],
    alpha: float,
    width: int = 6,
    blur: float = 0.0,
) -> None:
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    h = size
    w = int(size * 0.62)
    x0 = cx - w // 2
    x1 = cx + w // 2
    y0 = cy - h // 2
    y1 = cy + h // 2
    fill = (color[0], color[1], color[2], int(color[3] * clamp(alpha)))
    draw.line((x0, y1, x0, y0), fill=fill, width=width)
    draw.line((x0, y0, x1, y1), fill=fill, width=width)
    draw.line((x1, y1, x1, y0), fill=fill, width=width)
    if blur > 0:
        layer = layer.filter(ImageFilter.GaussianBlur(blur))
    base.alpha_composite(layer)


def add_vignette(base: Image.Image, color: tuple[int, int, int] = (0, 0, 0), strength: int = 170) -> None:
    vignette = Image.new("L", (WIDTH, HEIGHT), 0)
    draw = ImageDraw.Draw(vignette)
    draw.ellipse((-260, -210, WIDTH + 260, HEIGHT + 210), fill=255)
    vignette = Image.eval(vignette.filter(ImageFilter.GaussianBlur(140)), lambda value: int((255 - value) * strength / 255))
    overlay = Image.new("RGBA", (WIDTH, HEIGHT), (color[0], color[1], color[2], 0))
    overlay.putalpha(vignette)
    base.alpha_composite(overlay)


def background_noir(frame_index: int, time: float) -> Image.Image:
    base = Image.new("RGBA", (WIDTH, HEIGHT), (4, 5, 6, 255))
    soft_ellipse(base, (110, 260, 970, 1300), (210, 220, 220, 18), 120)
    soft_ellipse(base, (330, 480, 750, 1050), (255, 255, 245, 16), 80)
    draw = ImageDraw.Draw(base)
    y = 955 + int(math.sin(time * 1.4) * 10)
    draw.line((270, y, 810, y), fill=(238, 232, 212, 32), width=1)
    draw.line((402, y + 16, 678, y + 16), fill=(67, 255, 218, 22), width=1)
    draw_monogram_n(base, cx=WIDTH // 2, cy=700, size=430, color=(245, 242, 232, 96), alpha=smoothstep(0.08, 1.1, time), width=5, blur=0.4)
    add_grain(base, frame_index, 8)
    add_vignette(base, strength=160)
    return base


def background_ivory(frame_index: int, time: float) -> Image.Image:
    base = Image.new("RGBA", (WIDTH, HEIGHT), (232, 229, 218, 255))
    soft_ellipse(base, (-120, 80, 760, 1140), (255, 255, 255, 96), 100)
    soft_ellipse(base, (420, 660, 1260, 1850), (188, 176, 150, 34), 120)
    draw = ImageDraw.Draw(base)
    line_y = 1038
    sweep = smoothstep(0.1, 1.7, time)
    draw.line((250, line_y, 250 + int(580 * sweep), line_y), fill=(34, 34, 32, 66), width=1)
    draw_monogram_n(base, cx=WIDTH // 2, cy=710, size=380, color=(24, 25, 24, 54), alpha=smoothstep(0.2, 1.0, time), width=5, blur=0.2)
    add_grain(base, frame_index, 5)
    return base


def background_platinum(frame_index: int, time: float) -> Image.Image:
    base = Image.new("RGBA", (WIDTH, HEIGHT), (12, 13, 15, 255))
    soft_ellipse(base, (20, 160, 1060, 1280), (165, 176, 178, 32), 110)
    soft_ellipse(base, (285, 445, 795, 965), (255, 255, 255, 24), 80)
    ring = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(ring)
    rotation = int(time * 34)
    for radius, alpha in ((310, 60), (394, 30)):
        bbox = (WIDTH // 2 - radius, 694 - int(radius * 0.36), WIDTH // 2 + radius, 694 + int(radius * 0.36))
        draw.arc(bbox, rotation, rotation + 128, fill=(230, 232, 226, alpha), width=2)
        draw.arc(bbox, rotation + 190, rotation + 235, fill=(69, 255, 220, int(alpha * 0.45)), width=2)
    base.alpha_composite(ring.filter(ImageFilter.GaussianBlur(0.45)))
    draw_monogram_n(base, cx=WIDTH // 2, cy=690, size=350, color=(235, 236, 226, 68), alpha=smoothstep(0.1, 0.9, time), width=4, blur=0.4)
    add_grain(base, frame_index, 9)
    add_vignette(base, strength=145)
    return base


def background_gallery(frame_index: int, time: float) -> Image.Image:
    base = Image.new("RGBA", (WIDTH, HEIGHT), (2, 2, 3, 255))
    beam = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(beam)
    cx = WIDTH // 2 + int(math.sin(time * 0.9) * 20)
    draw.polygon((cx - 140, -40, cx + 140, -40, cx + 320, HEIGHT, cx - 320, HEIGHT), fill=(255, 255, 248, 18))
    base.alpha_composite(beam.filter(ImageFilter.GaussianBlur(44)))
    draw = ImageDraw.Draw(base)
    for x in (248, 832):
        alpha = int(18 + 12 * math.sin(time * 2.0 + x))
        draw.line((x, 320, x, 1500), fill=(255, 255, 255, alpha), width=1)
    draw_monogram_n(base, cx=WIDTH // 2, cy=720, size=390, color=(255, 255, 248, 44), alpha=smoothstep(0.2, 1.1, time), width=4, blur=0.7)
    add_grain(base, frame_index, 7)
    add_vignette(base, strength=185)
    return base


def background_emerald(frame_index: int, time: float) -> Image.Image:
    base = Image.new("RGBA", (WIDTH, HEIGHT), (0, 9, 9, 255))
    soft_ellipse(base, (60, 240, 1020, 1360), (27, 154, 128, 42), 120)
    soft_ellipse(base, (420, 560, 1180, 1520), (26, 255, 200, 24), 120)
    draw = ImageDraw.Draw(base)
    for idx in range(8):
        y = 430 + idx * 122 + int(math.sin(time * 1.5 + idx) * 12)
        draw.line((230, y, 850, y), fill=(88, 255, 218, 12 + idx % 2 * 8), width=1)
    sweep = smoothstep(0.2, 1.82, time)
    x = -220 + int((WIDTH + 440) * sweep)
    draw.line((x - 260, 1500, x + 220, 420), fill=(80, 255, 218, 50), width=2)
    draw_monogram_n(base, cx=WIDTH // 2, cy=705, size=410, color=(80, 255, 218, 54), alpha=smoothstep(0.1, 1.0, time), width=4, blur=0.4)
    add_grain(base, frame_index, 8)
    add_vignette(base, color=(0, 8, 8), strength=160)
    return base


def render_variant(name: str) -> dict[str, object]:
    if name == "01_noir":
        text_layer, text_box, font_path = make_text_layer(
            font_size=102,
            tracking=10,
            fill=(245, 244, 236, 255),
            accent=(245, 244, 236, 255),
            y=1120,
            shadow=(0, 8, (0, 0, 0, 132), 8),
        )
        bg = background_noir
        reveal = "center"
        sweep_color = (255, 255, 248)
    elif name == "02_ivory":
        text_layer, text_box, font_path = make_text_layer(
            font_size=104,
            tracking=9,
            fill=(18, 20, 19, 255),
            accent=(18, 20, 19, 255),
            y=1126,
            shadow=(0, 10, (178, 170, 150, 68), 12),
        )
        bg = background_ivory
        reveal = "soft"
        sweep_color = (255, 255, 255)
    elif name == "03_platinum":
        text_layer, text_box, font_path = make_text_layer(
            font_size=112,
            tracking=6,
            fill=(238, 240, 236, 255),
            accent=(92, 255, 218, 255),
            y=1134,
            shadow=(0, 9, (0, 0, 0, 150), 9),
        )
        bg = background_platinum
        reveal = "wipe"
        sweep_color = (255, 255, 255)
    elif name == "04_gallery":
        text_layer, text_box, font_path = make_text_layer(
            font_size=96,
            tracking=16,
            fill=(250, 250, 244, 255),
            accent=(250, 250, 244, 255),
            y=1138,
            shadow=(0, 8, (0, 0, 0, 145), 7),
        )
        bg = background_gallery
        reveal = "center"
        sweep_color = (255, 255, 248)
    elif name == "05_emerald":
        text_layer, text_box, font_path = make_text_layer(
            font_size=110,
            tracking=7,
            fill=(245, 252, 249, 255),
            accent=(73, 255, 218, 255),
            y=1130,
            shadow=(0, 9, (0, 8, 8, 164), 9),
        )
        bg = background_emerald
        reveal = "wipe"
        sweep_color = (84, 255, 218)
    else:
        raise ValueError(f"Unknown variant {name}")

    frames: list[Image.Image] = []
    for index in range(TOTAL_FRAMES):
        time = index / FPS
        frame = bg(index, time)
        reveal_text(frame, text_layer, text_box, time, start=0.42, end=1.18, mode=reveal, y_shift=18)
        add_light_sweep(frame, text_box, time, start=1.18, end=1.72, color=sweep_color, strength=58)
        frames.append(frame.convert("RGB"))

    output = OUTPUTS[name]
    encode_video(frames, output)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / f"neriva_luxury_minimal_v18_{name}_first.png"
    mid = PREVIEW_DIR / f"neriva_luxury_minimal_v18_{name}_mid.png"
    last = PREVIEW_DIR / f"neriva_luxury_minimal_v18_{name}_last.png"
    sheet = PREVIEW_DIR / f"neriva_luxury_minimal_v18_{name}_contact_sheet.jpg"
    frames[0].save(first)
    frames[TOTAL_FRAMES // 2].save(mid)
    frames[-1].save(last)
    make_contact_sheet(frames, sheet)

    return {
        "name": name,
        "video": str(output),
        "first_frame": str(first),
        "mid_frame": str(mid),
        "last_frame": str(last),
        "contact_sheet": str(sheet),
        "font": str(font_path),
        "brand_text": TEXT,
        "width": WIDTH,
        "height": HEIGHT,
        "fps": FPS,
        "frames": TOTAL_FRAMES,
        "duration": TOTAL_FRAMES / FPS,
    }


def make_contact_sheet(frames: list[Image.Image], path: Path) -> None:
    picks = [0, 10, 20, 32, 44, 56, 63]
    thumb_w = 216
    thumb_h = 384
    sheet = Image.new("RGB", (thumb_w * len(picks), thumb_h), (0, 0, 0))
    for col, frame_index in enumerate(picks):
        thumb = frames[frame_index].resize((thumb_w, thumb_h), Image.Resampling.LANCZOS)
        sheet.paste(thumb, (col * thumb_w, 0))
    sheet.save(path, quality=92)


def make_variant_overview(results: list[dict[str, object]]) -> Path:
    thumb_w = 216
    thumb_h = 384
    overview = Image.new("RGB", (thumb_w * len(results), thumb_h), (0, 0, 0))
    for index, result in enumerate(results):
        image = Image.open(str(result["last_frame"])).convert("RGB").resize((thumb_w, thumb_h), Image.Resampling.LANCZOS)
        draw = ImageDraw.Draw(image)
        draw.rectangle((0, 0, thumb_w, 34), fill=(0, 0, 0))
        draw.text((10, 9), str(result["name"]), fill=(255, 255, 255))
        overview.paste(image, (index * thumb_w, 0))
    path = PREVIEW_DIR / "neriva_luxury_minimal_v18_overview.jpg"
    overview.save(path, quality=94)
    return path


def render() -> dict[str, object]:
    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    results = [render_variant(name) for name in OUTPUTS]
    overview = make_variant_overview(results)
    result = {"variants": results, "overview": str(overview)}
    (PREVIEW_DIR / "neriva_luxury_minimal_v18_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
