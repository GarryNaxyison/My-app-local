from __future__ import annotations

import json
import random

from PIL import Image, ImageDraw, ImageFilter, ImageFont

import render_poliglot_tech_outro_v15_clean_vector_text as v15
from render_neriva_tech_outro_v16_clean_text_slow_spin import draw_slow_coin_spin_logo
from render_poliglot_tech_outro_v11 import (
    COMFY_OUTPUT,
    FPS,
    HEIGHT,
    LOGO_PATH,
    PREVIEW_DIR,
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


TEXT = "NERIVA.RU"
TEXT_PREFIX = "NERIVA."
TEXT_ACCENT = "RU"
SOURCE_VIDEO = COMFY_OUTPUT / "poliglot_outro_2s_tech_warp_v11_api_00001_.mp4"
OUTPUT_VIDEO = COMFY_OUTPUT / "neriva_outro_tech_warp_v17_upper_comet_text_lockup.mp4"


def make_comet_char_sprite(char: str, font: ImageFont.FreeTypeFont, accent: bool) -> tuple[Image.Image, int, int]:
    probe = Image.new("RGBA", (1, 1), (0, 0, 0, 0))
    probe_draw = ImageDraw.Draw(probe)
    bbox = probe_draw.textbbox((0, 0), char, font=font, stroke_width=2)
    pad_x = 20
    pad_y = 24
    width = max(34, bbox[2] - bbox[0] + pad_x * 2)
    height = max(96, bbox[3] - bbox[1] + pad_y * 2)

    glow = Image.new("RGBA", (width, height), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    x = pad_x - bbox[0]
    y = pad_y - bbox[1]
    fill = (80, 255, 218, 244) if accent else (246, 253, 255, 248)
    stroke = (1, 18, 24, 210)
    shadow = (0, 8, 14, 170)

    glow_draw.text((x, y), char, font=font, fill=(66, 255, 224, 62), stroke_width=2, stroke_fill=(66, 255, 224, 58))
    glow = glow.filter(ImageFilter.GaussianBlur(4.5))

    sprite = Image.new("RGBA", (width, height), (0, 0, 0, 0))
    sprite.alpha_composite(glow)
    draw = ImageDraw.Draw(sprite)
    draw.text((x + 4, y + 7), char, font=font, fill=shadow, stroke_width=2, stroke_fill=shadow)
    draw.text((x, y), char, font=font, fill=fill, stroke_width=2, stroke_fill=stroke)
    return sprite, pad_x, pad_y


def make_comet_text_assets(font: ImageFont.FreeTypeFont) -> tuple[list[dict[str, object]], tuple[int, int, int, int]]:
    probe = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(probe)
    bbox = draw.textbbox((0, 0), TEXT, font=font, stroke_width=2)
    text_width = bbox[2] - bbox[0]
    x = (WIDTH - text_width) // 2 - bbox[0]
    y = 1260 - bbox[1]
    text_box = (x + bbox[0], y + bbox[1], x + bbox[2], y + bbox[3])

    glyphs: list[dict[str, object]] = []
    accent_start = len(TEXT_PREFIX)
    accent_end = accent_start + len(TEXT_ACCENT)
    for index, char in enumerate(TEXT):
        accent = accent_start <= index < accent_end
        sprite, pad_x, pad_y = make_comet_char_sprite(char, font, accent)
        char_bbox = draw.textbbox((0, 0), char, font=font, stroke_width=2)
        draw_x = x + int(round(draw.textlength(TEXT[:index], font=font)))
        final_x = draw_x + char_bbox[0] - pad_x
        final_y = y + char_bbox[1] - pad_y
        rng = random.Random(925000 + index)
        glyphs.append(
            {
                "sprite": sprite,
                "final_x": final_x,
                "final_y": final_y,
                "delay": 0.72 + index * 0.025,
                "tilt": rng.uniform(-5.0, 5.0),
                "start_jitter": rng.uniform(-42, 42),
            }
        )
    return glyphs, text_box


def draw_comet_text_lockup(base: Image.Image, glyphs: list[dict[str, object]], text_box: tuple[int, int, int, int], time: float) -> None:
    draw_comet_path(base, text_box, time)

    particle_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    particle_draw = ImageDraw.Draw(particle_layer)

    for index, glyph in enumerate(glyphs):
        delay = float(glyph["delay"])
        local = smoothstep(delay, delay + 0.52, time)
        if local <= 0:
            continue

        p = ease_out_cubic(local)
        final_x = float(glyph["final_x"])
        final_y = float(glyph["final_y"])
        start = (-190 + index * 22, text_box[1] - 214 + float(glyph["start_jitter"]))
        c1 = (120 + index * 14, text_box[1] - 322)
        c2 = (final_x - 138, final_y + 72)
        x, y = bezier(start, c1, c2, (final_x, final_y), p)

        sprite = glyph["sprite"]
        assert isinstance(sprite, Image.Image)
        scale = 0.55 + 0.45 * min(1.0, p)
        char = sprite.resize((max(1, int(sprite.width * scale)), max(1, int(sprite.height * scale))), Image.Resampling.LANCZOS)
        angle = float(glyph["tilt"]) * (1 - local)
        char = char.rotate(angle, resample=Image.Resampling.BICUBIC, expand=True)
        if local < 0.78:
            char = char.filter(ImageFilter.GaussianBlur(2.2 * (1 - local)))
        base.alpha_composite(apply_alpha(char, clamp(local * 1.34)), (int(x), int(y)))

        if local < 0.94:
            rng = random.Random(44000 + index)
            for _ in range(5):
                px = x + rng.uniform(-24, 28) - (1 - local) * rng.uniform(20, 115)
                py = y + rng.uniform(12, 72) + (1 - local) * rng.uniform(-28, 28)
                size = rng.uniform(1.2, 4.2)
                alpha = int(124 * local * (1 - local * 0.55))
                particle_draw.ellipse((px - size, py - size, px + size, py + size), fill=(142, 255, 232, alpha))

    base.alpha_composite(particle_layer.filter(ImageFilter.GaussianBlur(0.45)))


def render() -> dict:
    logo = Image.open(LOGO_PATH).convert("RGBA")
    font, font_path = v15.load_clean_font(138)
    glyphs, text_box = make_comet_text_assets(font)
    source_frames = decode_source(SOURCE_VIDEO)

    frames: list[Image.Image] = []
    for index, source_frame in enumerate(source_frames):
        time = index / FPS
        base = source_background(source_frame, index)
        draw_slow_coin_spin_logo(base, time, logo)
        draw_comet_text_lockup(base, glyphs, text_box, time)
        add_finish_sweep(base, time)
        add_vignette(base)
        frames.append(base.convert("RGB"))

    encode_video(frames, OUTPUT_VIDEO)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / "neriva_outro_tech_warp_v17_upper_first.png"
    mid = PREVIEW_DIR / "neriva_outro_tech_warp_v17_upper_mid.png"
    last = PREVIEW_DIR / "neriva_outro_tech_warp_v17_upper_last.png"
    sheet = PREVIEW_DIR / "neriva_outro_tech_warp_v17_upper_contact_sheet.jpg"
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
        "font": str(font_path),
        "brand_text": TEXT,
        "text_mode": "comet glyphs remain as final lockup; no full-line overlay",
        "width": WIDTH,
        "height": HEIGHT,
        "fps": FPS,
        "frames": TOTAL_FRAMES,
        "duration": TOTAL_FRAMES / FPS,
    }
    (PREVIEW_DIR / "neriva_outro_tech_warp_v17_upper_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
