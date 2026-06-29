from __future__ import annotations

import json
import math
import random
from pathlib import Path

from PIL import Image, ImageDraw, ImageFilter, ImageFont

from render_poliglot_tech_outro_v11 import (
    COMFY_OUTPUT,
    FPS,
    HEIGHT,
    LOGO_PATH,
    MANROPE_FONT_PATH,
    PREVIEW_DIR,
    TEXT,
    TEXT_ACCENT,
    TEXT_PREFIX,
    TOTAL_FRAMES,
    WIDTH,
    add_finish_sweep,
    add_vignette,
    apply_alpha,
    clamp,
    ease_out_cubic,
    encode_video,
    smoothstep,
)
from render_poliglot_tech_outro_v12_comet_text import decode_source, make_contact_sheet, source_background
from render_poliglot_tech_outro_v13_motion_logo_text import bezier, draw_comet_path
from render_poliglot_tech_outro_v14_coin_logo_upright_text import draw_coin_spin_logo


SOURCE_VIDEO = COMFY_OUTPUT / "poliglot_outro_2s_tech_warp_v11_api_00001_.mp4"
OUTPUT_VIDEO = COMFY_OUTPUT / "poliglot_outro_tech_warp_v15_clean_vector_text.mp4"

CLEAN_FONT_CANDIDATES = (
    Path(r"C:\Windows\Fonts\segoeuib.ttf"),
    Path(r"C:\Windows\Fonts\arialbd.ttf"),
    MANROPE_FONT_PATH,
)


def load_clean_font(size: int) -> tuple[ImageFont.FreeTypeFont, Path]:
    for path in CLEAN_FONT_CANDIDATES:
        if path.exists():
            return ImageFont.truetype(str(path), size=size), path
    raise FileNotFoundError("No clean TrueType font found for outro text")


def draw_segmented_url(
    draw: ImageDraw.ImageDraw,
    xy: tuple[float, float],
    font: ImageFont.FreeTypeFont,
    alpha: int,
    stroke_alpha: int,
) -> None:
    x, y = xy
    stroke = (1, 18, 24, stroke_alpha)
    white = (246, 253, 255, alpha)
    mint = (75, 255, 214, alpha)

    draw.text((x, y), TEXT, font=font, fill=white, stroke_width=2, stroke_fill=stroke)
    accent_x = x + draw.textlength(TEXT_PREFIX, font=font)
    draw.text((accent_x, y), TEXT_ACCENT, font=font, fill=mint, stroke_width=2, stroke_fill=stroke)


def make_clean_url_layer(font: ImageFont.FreeTypeFont) -> tuple[Image.Image, Image.Image, tuple[int, int, int, int]]:
    probe = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(probe)
    bbox = draw.textbbox((0, 0), TEXT, font=font, stroke_width=2)
    text_width = bbox[2] - bbox[0]
    target_top = 1260
    x = (WIDTH - text_width) // 2 - bbox[0]
    y = target_top - bbox[1]
    text_box = (x + bbox[0], y + bbox[1], x + bbox[2], y + bbox[3])

    mask = Image.new("L", (WIDTH, HEIGHT), 0)
    mask_draw = ImageDraw.Draw(mask)
    mask_draw.text((x, y), TEXT, font=font, fill=255, stroke_width=2, stroke_fill=255)

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    glow_draw.text((x, y), TEXT, font=font, fill=(64, 255, 225, 82), stroke_width=2, stroke_fill=(64, 255, 225, 70))
    glow = glow.filter(ImageFilter.GaussianBlur(10))

    shadow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    shadow_draw = ImageDraw.Draw(shadow)
    shadow_draw.text((x + 4, y + 8), TEXT, font=font, fill=(0, 7, 12, 238), stroke_width=2, stroke_fill=(0, 7, 12, 238))
    shadow = shadow.filter(ImageFilter.GaussianBlur(6))

    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    layer.alpha_composite(glow)
    layer.alpha_composite(shadow)
    layer_draw = ImageDraw.Draw(layer)
    draw_segmented_url(layer_draw, (x, y), font, 255, 220)
    return layer, mask, text_box


def make_flight_char_sprite(char: str, font: ImageFont.FreeTypeFont, accent: bool) -> Image.Image:
    probe = Image.new("RGBA", (1, 1), (0, 0, 0, 0))
    probe_draw = ImageDraw.Draw(probe)
    bbox = probe_draw.textbbox((0, 0), char, font=font, stroke_width=1)
    width = max(30, bbox[2] - bbox[0] + 34)
    height = max(92, bbox[3] - bbox[1] + 46)
    sprite = Image.new("RGBA", (width, height), (0, 0, 0, 0))
    draw = ImageDraw.Draw(sprite)
    x = 17 - bbox[0]
    y = 20 - bbox[1]
    fill = (78, 255, 218, 248) if accent else (246, 253, 255, 248)
    stroke = (2, 18, 28, 168)
    draw.text((x + 2, y + 4), char, font=font, fill=(0, 8, 16, 120), stroke_width=1, stroke_fill=(0, 8, 16, 120))
    draw.text((x, y), char, font=font, fill=fill, stroke_width=1, stroke_fill=stroke)
    return sprite


def make_motion_url_assets(
    final_font: ImageFont.FreeTypeFont,
    flight_font: ImageFont.FreeTypeFont,
) -> tuple[Image.Image, Image.Image, list[dict[str, object]], tuple[int, int, int, int]]:
    final_layer, final_mask, text_box = make_clean_url_layer(final_font)

    probe = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(probe)
    accent_start = len(TEXT_PREFIX)
    accent_end = accent_start + len(TEXT_ACCENT)
    glyphs: list[dict[str, object]] = []
    for index, char in enumerate(TEXT):
        accent = accent_start <= index < accent_end
        rng = random.Random(124800 + index)
        final_x = text_box[0] + int(round(draw.textlength(TEXT[:index], font=final_font))) - 12
        glyphs.append(
            {
                "sprite": make_flight_char_sprite(char, flight_font, accent),
                "final_x": final_x,
                "final_y": text_box[1] - 12 + (30 if char == "." else 0),
                "delay": 0.72 + index * 0.023,
                "tilt": rng.uniform(-5.5, 5.5),
                "start_jitter": rng.uniform(-42, 42),
            }
        )
    return final_layer, final_mask, glyphs, text_box


def draw_text_glint(base: Image.Image, final_mask: Image.Image, text_box: tuple[int, int, int, int], time: float) -> None:
    sweep = smoothstep(1.44, 1.78, time)
    if sweep <= 0 or sweep >= 1:
        return

    x0, y0, x1, y1 = text_box
    width = x1 - x0
    band_x = x0 - 155 + int((width + 310) * ease_out_cubic(sweep))
    band = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(band)
    draw.polygon(
        [
            (band_x - 42, y0 - 24),
            (band_x + 16, y0 - 24),
            (band_x + 84, y1 + 34),
            (band_x + 25, y1 + 34),
        ],
        fill=(255, 255, 255, int(105 * math.sin(math.pi * sweep))),
    )
    band_alpha = Image.composite(band.getchannel("A"), Image.new("L", (WIDTH, HEIGHT), 0), final_mask)
    band.putalpha(band_alpha)
    base.alpha_composite(band.filter(ImageFilter.GaussianBlur(0.6)))


def draw_url_motion(
    base: Image.Image,
    final_layer: Image.Image,
    final_mask: Image.Image,
    glyphs: list[dict[str, object]],
    text_box: tuple[int, int, int, int],
    time: float,
) -> None:
    draw_comet_path(base, text_box, time)

    particle_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    particle_draw = ImageDraw.Draw(particle_layer)
    piece_fade = 1.0 - smoothstep(1.28, 1.52, time)

    for index, glyph in enumerate(glyphs):
        delay = float(glyph["delay"])
        local = smoothstep(delay, delay + 0.50, time)
        if local <= 0:
            continue

        p = ease_out_cubic(local)
        final_x = float(glyph["final_x"])
        final_y = float(glyph["final_y"])
        start = (-190 + index * 18, text_box[1] - 214 + float(glyph["start_jitter"]))
        c1 = (125 + index * 10, text_box[1] - 318)
        c2 = (final_x - 130, final_y + 70)
        x, y = bezier(start, c1, c2, (final_x, final_y), p)

        sprite = glyph["sprite"]
        assert isinstance(sprite, Image.Image)
        scale = 0.54 + 0.46 * min(1.0, p)
        char = sprite.resize((max(1, int(sprite.width * scale)), max(1, int(sprite.height * scale))), Image.Resampling.LANCZOS)
        angle = float(glyph["tilt"]) * (1 - local)
        char = char.rotate(angle, resample=Image.Resampling.BICUBIC, expand=True)
        if local < 0.80:
            char = char.filter(ImageFilter.GaussianBlur(2.4 * (1 - local)))
        alpha = clamp(local * 1.32) * piece_fade
        if alpha > 0:
            base.alpha_composite(apply_alpha(char, alpha), (int(x), int(y)))

        if local < 0.94:
            rng = random.Random(43300 + index)
            for _ in range(5):
                px = x + rng.uniform(-24, 28) - (1 - local) * rng.uniform(20, 115)
                py = y + rng.uniform(12, 72) + (1 - local) * rng.uniform(-28, 28)
                size = rng.uniform(1.2, 4.2)
                particle_alpha = int(124 * local * (1 - local * 0.55))
                particle_draw.ellipse((px - size, py - size, px + size, py + size), fill=(142, 255, 232, particle_alpha))

    base.alpha_composite(particle_layer.filter(ImageFilter.GaussianBlur(0.45)))

    final_alpha = smoothstep(1.18, 1.48, time)
    if final_alpha > 0:
        base.alpha_composite(apply_alpha(final_layer, final_alpha))
        draw_text_glint(base, final_mask, text_box, time)


def render() -> dict:
    logo = Image.open(LOGO_PATH).convert("RGBA")
    final_font, final_font_path = load_clean_font(100)
    flight_font, _ = load_clean_font(74)
    final_layer, final_mask, glyphs, text_box = make_motion_url_assets(final_font, flight_font)
    source_frames = decode_source(SOURCE_VIDEO)

    frames: list[Image.Image] = []
    for index, source_frame in enumerate(source_frames):
        time = index / FPS
        base = source_background(source_frame, index)
        draw_coin_spin_logo(base, time, logo)
        draw_url_motion(base, final_layer, final_mask, glyphs, text_box, time)
        add_finish_sweep(base, time)
        add_vignette(base)
        frames.append(base.convert("RGB"))

    encode_video(frames, OUTPUT_VIDEO)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / "outro_tech_warp_v15_first.png"
    mid = PREVIEW_DIR / "outro_tech_warp_v15_mid.png"
    last = PREVIEW_DIR / "outro_tech_warp_v15_last.png"
    sheet = PREVIEW_DIR / "outro_tech_warp_v15_contact_sheet.jpg"
    frames[0].save(first)
    frames[TOTAL_FRAMES // 2].save(mid)
    frames[-1].save(last)
    make_contact_sheet(frames, sheet)

    result = {
        "source": str(SOURCE_VIDEO),
        "video": str(OUTPUT_VIDEO),
        "first_frame": str(first),
        "mid_frame": str(mid),
        "last_frame": str(last),
        "contact_sheet": str(sheet),
        "font": str(final_font_path),
        "width": WIDTH,
        "height": HEIGHT,
        "fps": FPS,
        "frames": TOTAL_FRAMES,
        "duration": TOTAL_FRAMES / FPS,
    }
    (PREVIEW_DIR / "poliglot_outro_tech_warp_v15_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
