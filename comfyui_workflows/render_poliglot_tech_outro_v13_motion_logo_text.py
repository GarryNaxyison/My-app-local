from __future__ import annotations

import json
import math
import random
from pathlib import Path

import av
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
    draw_logo_shards,
    draw_orbits,
    ease_out_back,
    ease_out_cubic,
    encode_video,
    load_font,
    smoothstep,
)
from render_poliglot_tech_outro_v12_comet_text import decode_source, make_contact_sheet, source_background


SOURCE_VIDEO = COMFY_OUTPUT / "poliglot_outro_2s_tech_warp_v11_api_00001_.mp4"
OUTPUT_VIDEO = COMFY_OUTPUT / "poliglot_outro_tech_warp_v13_motion_logo_text.mp4"


def draw_rotating_logo(base: Image.Image, time: float, logo: Image.Image) -> None:
    if time < 0.94:
        draw_logo_shards(base, time, logo)
        return

    size = 512
    center_x = WIDTH // 2
    center_y = 690
    bob = math.sin((time - 0.94) * math.tau * 0.85) * 10
    yaw = math.sin((time - 0.94) * math.tau * 0.92) * 0.38
    width = int(size * (1.0 - abs(yaw) * 0.24))
    height = int(size * (1.0 + abs(yaw) * 0.035))
    sprite = logo.resize((width, height), Image.Resampling.LANCZOS)
    sprite = sprite.rotate(yaw * 8.5, resample=Image.Resampling.BICUBIC, expand=True)
    x = center_x - sprite.width // 2
    y = int(center_y - sprite.height // 2 - 28 + bob)

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(glow)
    draw.ellipse((x - 150, y - 128, x + sprite.width + 150, y + sprite.height + 155), fill=(65, 255, 215, 48))
    draw.ellipse((x + 26, y + 46, x + sprite.width - 26, y + sprite.height + 62), fill=(46, 102, 255, 44))
    base.alpha_composite(glow.filter(ImageFilter.GaussianBlur(58)))
    draw_orbits(base, center_x, y + sprite.height // 2, time, 1.0)

    base.alpha_composite(sprite, (x, y))

    sweep = ((time - 0.94) * 0.72) % 1.0
    glint = Image.new("RGBA", sprite.size, (0, 0, 0, 0))
    glint_draw = ImageDraw.Draw(glint)
    sx = int(-sprite.width * 0.55 + sweep * sprite.width * 1.85)
    glint_draw.polygon(
        [(sx, -20), (sx + 72, -20), (sx + 16, sprite.height + 20), (sx - 56, sprite.height + 20)],
        fill=(255, 255, 255, 64),
    )
    glint_draw.polygon(
        [(sx + 76, -20), (sx + 118, -20), (sx + 56, sprite.height + 20), (sx + 14, sprite.height + 20)],
        fill=(96, 255, 220, 42),
    )
    alpha_mask = sprite.getchannel("A")
    glint.putalpha(Image.composite(glint.getchannel("A"), Image.new("L", sprite.size, 0), alpha_mask))
    base.alpha_composite(glint.filter(ImageFilter.GaussianBlur(1.2)), (x, y))


def make_char_sprite(char: str, font: ImageFont.FreeTypeFont, accent: bool) -> Image.Image:
    probe = Image.new("RGBA", (160, 180), (0, 0, 0, 0))
    draw = ImageDraw.Draw(probe)
    bbox = draw.textbbox((0, 0), char, font=font, stroke_width=1)
    width = max(32, int(math.ceil(draw.textlength(char, font=font))) + 34)
    height = max(88, bbox[3] - bbox[1] + 42)
    sprite = Image.new("RGBA", (width, height), (0, 0, 0, 0))
    glow = Image.new("RGBA", (width, height), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    sprite_draw = ImageDraw.Draw(sprite)
    color = (86, 255, 216, 255) if accent else (248, 255, 253, 255)
    glow_draw.text((17 - bbox[0], 18 - bbox[1]), char, font=font, fill=(color[0], color[1], color[2], 118), stroke_width=1, stroke_fill=(80, 255, 220, 56))
    sprite.alpha_composite(glow.filter(ImageFilter.GaussianBlur(4)))
    sprite_draw.text((17 - bbox[0], 18 - bbox[1]), char, font=font, fill=color, stroke_width=1, stroke_fill=(0, 20, 25, 238))
    return sprite


def make_url_assets(font: ImageFont.FreeTypeFont) -> tuple[Image.Image, list[dict[str, object]], tuple[int, int, int, int]]:
    probe = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(probe)
    bbox = draw.textbbox((0, 0), TEXT, font=font, stroke_width=1)
    text_width = int(math.ceil(draw.textlength(TEXT, font=font)))
    text_height = bbox[3] - bbox[1]
    x = (WIDTH - text_width) // 2
    y = 1256

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    glow_draw.rounded_rectangle((x - 92, y - 34, x + text_width + 92, y + text_height + 62), radius=38, fill=(0, 13, 18, 104))
    glow_draw.rounded_rectangle((x - 96, y - 38, x + text_width + 96, y + text_height + 66), radius=40, fill=(76, 255, 220, 18))
    glow = glow.filter(ImageFilter.GaussianBlur(8))

    glyphs: list[dict[str, object]] = []
    accent_start = len(TEXT_PREFIX)
    accent_end = accent_start + len(TEXT_ACCENT)
    for index, char in enumerate(TEXT):
        final_x = x + int(round(draw.textlength(TEXT[:index], font=font))) - 17
        accent = accent_start <= index < accent_end
        sprite = make_char_sprite(char, font, accent)
        rng = random.Random(123300 + index)
        glyphs.append(
            {
                "sprite": sprite,
                "final_x": final_x,
                "final_y": y + (32 if char == "." else -18),
                "delay": 0.72 + index * 0.024,
                "spin": rng.uniform(-26, 26),
                "start_jitter": rng.uniform(-45, 45),
            }
        )
    return glow, glyphs, (x, y, x + text_width, y + text_height)


def bezier(p0: tuple[float, float], p1: tuple[float, float], p2: tuple[float, float], p3: tuple[float, float], t: float) -> tuple[float, float]:
    inv = 1 - t
    x = inv**3 * p0[0] + 3 * inv * inv * t * p1[0] + 3 * inv * t * t * p2[0] + t**3 * p3[0]
    y = inv**3 * p0[1] + 3 * inv * inv * t * p1[1] + 3 * inv * t * t * p2[1] + t**3 * p3[1]
    return x, y


def draw_comet_path(base: Image.Image, text_box: tuple[int, int, int, int], time: float) -> None:
    raw = smoothstep(0.62, 1.32, time)
    if raw <= 0 or time > 1.58:
        return
    p = ease_out_cubic(raw)
    start = (-260, text_box[1] - 160)
    c1 = (130, text_box[1] - 335)
    c2 = (text_box[0] + 260, text_box[1] + 90)
    end = (text_box[2] + 90, text_box[1] + 20)
    hx, hy = bezier(start, c1, c2, end, p)
    tail = 1 - smoothstep(1.22, 1.58, time)

    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    for index in range(16):
        lag = max(0.0, p - index * 0.018)
        tx, ty = bezier(start, c1, c2, end, lag)
        alpha = int(118 * tail * (1 - index / 17))
        size = 15 + index * 4
        draw.ellipse((tx - size, ty - size * 0.45, tx + size, ty + size * 0.45), fill=(68, 255, 220, alpha))
    draw.ellipse((hx - 46, hy - 46, hx + 46, hy + 46), fill=(180, 255, 238, int(126 * tail)))
    draw.ellipse((hx - 18, hy - 18, hx + 18, hy + 18), fill=(255, 255, 255, int(220 * tail)))
    base.alpha_composite(layer.filter(ImageFilter.GaussianBlur(7)))


def draw_url_motion(base: Image.Image, glow: Image.Image, glyphs: list[dict[str, object]], text_box: tuple[int, int, int, int], time: float) -> None:
    draw_comet_path(base, text_box, time)

    glow_alpha = smoothstep(1.12, 1.55, time)
    if glow_alpha > 0:
        base.alpha_composite(apply_alpha(glow, glow_alpha))

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
        start = (-180 + index * 18, text_box[1] - 210 + float(glyph["start_jitter"]))
        c1 = (120 + index * 12, text_box[1] - 350)
        c2 = (final_x - 160, final_y + 80)
        end = (final_x, final_y)
        x, y = bezier(start, c1, c2, end, p)
        sprite = glyph["sprite"]
        assert isinstance(sprite, Image.Image)
        scale = 0.45 + 0.62 * min(1.0, ease_out_back(local))
        scale_w = max(1, int(sprite.width * scale))
        scale_h = max(1, int(sprite.height * scale))
        char = sprite.resize((scale_w, scale_h), Image.Resampling.LANCZOS)
        angle = float(glyph["spin"]) * (1 - local) + math.sin((time + index * 0.06) * 9) * 1.5 * (1 - smoothstep(1.42, 1.8, time))
        char = char.rotate(angle, resample=Image.Resampling.BICUBIC, expand=True)
        if local < 0.82:
            char = char.filter(ImageFilter.GaussianBlur(3.2 * (1 - local)))
        base.alpha_composite(apply_alpha(char, clamp(local * 1.35)), (int(x), int(y)))

        if local < 0.96:
            rng = random.Random(42200 + index)
            for _ in range(5):
                px = x + rng.uniform(-24, 28) - (1 - local) * rng.uniform(20, 130)
                py = y + rng.uniform(12, 74) + (1 - local) * rng.uniform(-36, 36)
                size = rng.uniform(1.2, 4.6)
                alpha = int(132 * local * (1 - local * 0.55))
                particle_draw.ellipse((px - size, py - size, px + size, py + size), fill=(142, 255, 232, alpha))

    base.alpha_composite(particle_layer.filter(ImageFilter.GaussianBlur(0.45)))

    lock = smoothstep(1.44, 1.66, time) * (1 - smoothstep(1.78, 2.0, time) * 0.35)
    if lock > 0:
        layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        draw = ImageDraw.Draw(layer)
        y = int((text_box[1] + text_box[3]) / 2 + 6)
        draw.line((text_box[0] - 44, y, text_box[2] + 44, y), fill=(95, 255, 222, int(74 * lock)), width=4)
        base.alpha_composite(layer.filter(ImageFilter.GaussianBlur(7)))


def render() -> dict:
    logo = Image.open(LOGO_PATH).convert("RGBA")
    font = load_font(92, bold=True)
    glow, glyphs, text_box = make_url_assets(font)
    source_frames = decode_source(SOURCE_VIDEO)

    frames: list[Image.Image] = []
    for index, source_frame in enumerate(source_frames):
        time = index / FPS
        base = source_background(source_frame, index)
        draw_rotating_logo(base, time, logo)
        draw_url_motion(base, glow, glyphs, text_box, time)
        add_finish_sweep(base, time)
        add_vignette(base)
        frames.append(base.convert("RGB"))

    encode_video(frames, OUTPUT_VIDEO)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / "outro_tech_warp_v13_first.png"
    mid = PREVIEW_DIR / "outro_tech_warp_v13_mid.png"
    last = PREVIEW_DIR / "outro_tech_warp_v13_last.png"
    sheet = PREVIEW_DIR / "outro_tech_warp_v13_contact_sheet.jpg"
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
        "font": str(MANROPE_FONT_PATH),
        "width": WIDTH,
        "height": HEIGHT,
        "fps": FPS,
        "frames": TOTAL_FRAMES,
        "duration": TOTAL_FRAMES / FPS,
    }
    (PREVIEW_DIR / "poliglot_outro_tech_warp_v13_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
