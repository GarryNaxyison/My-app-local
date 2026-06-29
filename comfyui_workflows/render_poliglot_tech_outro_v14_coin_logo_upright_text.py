from __future__ import annotations

import json
import math
import random

from PIL import Image, ImageDraw, ImageFilter, ImageFont, ImageOps

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
    ease_out_cubic,
    encode_video,
    load_font,
    smoothstep,
)
from render_poliglot_tech_outro_v12_comet_text import decode_source, make_contact_sheet, source_background
from render_poliglot_tech_outro_v13_motion_logo_text import bezier, draw_comet_path, make_char_sprite


SOURCE_VIDEO = COMFY_OUTPUT / "poliglot_outro_2s_tech_warp_v11_api_00001_.mp4"
OUTPUT_VIDEO = COMFY_OUTPUT / "poliglot_outro_tech_warp_v14_coin_logo_upright_text.mp4"


def draw_coin_spin_logo(base: Image.Image, time: float, logo: Image.Image) -> None:
    if time < 0.94:
        draw_logo_shards(base, time, logo)
        return

    center_x = WIDTH // 2
    center_y = 662
    size = 512
    spin = smoothstep(0.94, 1.36, time)
    theta = ease_out_cubic(spin) * math.tau
    face = math.cos(theta)
    width_scale = 1.0 if spin >= 1 else max(0.16, abs(face))
    sprite = logo.resize((max(1, int(size * width_scale)), size), Image.Resampling.LANCZOS)
    if spin < 1 and face < 0:
        sprite = ImageOps.mirror(sprite)

    x = center_x - sprite.width // 2
    y = center_y - sprite.height // 2

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(glow)
    draw.ellipse((center_x - 395, center_y - 360, center_x + 395, center_y + 405), fill=(64, 255, 215, 46))
    draw.ellipse((center_x - 255, center_y - 230, center_x + 255, center_y + 285), fill=(48, 102, 255, 42))
    base.alpha_composite(glow.filter(ImageFilter.GaussianBlur(58)))
    draw_orbits(base, center_x, center_y, time, 1.0)
    base.alpha_composite(sprite, (x, y))

    if spin < 1:
        glint = Image.new("RGBA", sprite.size, (0, 0, 0, 0))
        glint_draw = ImageDraw.Draw(glint)
        sx = int((spin * 1.35 - 0.28) * sprite.width)
        glint_draw.polygon(
            [(sx, -20), (sx + 70, -20), (sx + 20, sprite.height + 20), (sx - 50, sprite.height + 20)],
            fill=(255, 255, 255, int(82 * (1 - abs(spin - 0.5) * 0.7))),
        )
        alpha_mask = sprite.getchannel("A")
        glint.putalpha(Image.composite(glint.getchannel("A"), Image.new("L", sprite.size, 0), alpha_mask))
        base.alpha_composite(glint.filter(ImageFilter.GaussianBlur(1.0)), (x, y))


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
    glow_draw.ellipse((x - 70, y - 10, x + text_width + 70, y + text_height + 58), fill=(0, 18, 22, 92))
    glow_draw.line((x - 54, y + text_height + 22, x + text_width + 54, y + text_height + 22), fill=(72, 255, 220, 34), width=5)
    glow = glow.filter(ImageFilter.GaussianBlur(16))

    glyphs: list[dict[str, object]] = []
    accent_start = len(TEXT_PREFIX)
    accent_end = accent_start + len(TEXT_ACCENT)
    for index, char in enumerate(TEXT):
        final_x = x + int(round(draw.textlength(TEXT[:index], font=font))) - 17
        accent = accent_start <= index < accent_end
        rng = random.Random(123300 + index)
        glyphs.append(
            {
                "sprite": make_char_sprite(char, font, accent),
                "final_x": final_x,
                "final_y": y + (32 if char == "." else -18),
                "delay": 0.72 + index * 0.024,
                "tilt": rng.uniform(-8, 8),
                "start_jitter": rng.uniform(-45, 45),
            }
        )
    return glow, glyphs, (x, y, x + text_width, y + text_height)


def draw_url_motion(base: Image.Image, glow: Image.Image, glyphs: list[dict[str, object]], text_box: tuple[int, int, int, int], time: float) -> None:
    draw_comet_path(base, text_box, time)

    glow_alpha = smoothstep(1.18, 1.58, time)
    if glow_alpha > 0:
        base.alpha_composite(apply_alpha(glow, glow_alpha))

    particle_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    particle_draw = ImageDraw.Draw(particle_layer)

    for index, glyph in enumerate(glyphs):
        delay = float(glyph["delay"])
        local = smoothstep(delay, delay + 0.50, time)
        if local <= 0:
            continue
        p = ease_out_cubic(local)
        final_x = float(glyph["final_x"])
        final_y = float(glyph["final_y"])
        start = (-180 + index * 18, text_box[1] - 210 + float(glyph["start_jitter"]))
        c1 = (130 + index * 10, text_box[1] - 315)
        c2 = (final_x - 140, final_y + 72)
        end = (final_x, final_y)
        x, y = bezier(start, c1, c2, end, p)
        sprite = glyph["sprite"]
        assert isinstance(sprite, Image.Image)
        scale = 0.52 + 0.48 * min(1.0, p)
        char = sprite.resize((max(1, int(sprite.width * scale)), max(1, int(sprite.height * scale))), Image.Resampling.LANCZOS)
        angle = float(glyph["tilt"]) * (1 - local)
        char = char.rotate(angle, resample=Image.Resampling.BICUBIC, expand=True)
        if local < 0.78:
            char = char.filter(ImageFilter.GaussianBlur(2.6 * (1 - local)))
        base.alpha_composite(apply_alpha(char, clamp(local * 1.35)), (int(x), int(y)))

        if local < 0.94:
            rng = random.Random(42200 + index)
            for _ in range(5):
                px = x + rng.uniform(-24, 28) - (1 - local) * rng.uniform(20, 115)
                py = y + rng.uniform(12, 74) + (1 - local) * rng.uniform(-28, 28)
                size = rng.uniform(1.2, 4.4)
                alpha = int(128 * local * (1 - local * 0.55))
                particle_draw.ellipse((px - size, py - size, px + size, py + size), fill=(142, 255, 232, alpha))

    base.alpha_composite(particle_layer.filter(ImageFilter.GaussianBlur(0.45)))


def render() -> dict:
    logo = Image.open(LOGO_PATH).convert("RGBA")
    font = load_font(92, bold=True)
    glow, glyphs, text_box = make_url_assets(font)
    source_frames = decode_source(SOURCE_VIDEO)

    frames: list[Image.Image] = []
    for index, source_frame in enumerate(source_frames):
        time = index / FPS
        base = source_background(source_frame, index)
        draw_coin_spin_logo(base, time, logo)
        draw_url_motion(base, glow, glyphs, text_box, time)
        add_finish_sweep(base, time)
        add_vignette(base)
        frames.append(base.convert("RGB"))

    encode_video(frames, OUTPUT_VIDEO)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / "outro_tech_warp_v14_first.png"
    mid = PREVIEW_DIR / "outro_tech_warp_v14_mid.png"
    last = PREVIEW_DIR / "outro_tech_warp_v14_last.png"
    sheet = PREVIEW_DIR / "outro_tech_warp_v14_contact_sheet.jpg"
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
    (PREVIEW_DIR / "poliglot_outro_tech_warp_v14_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
