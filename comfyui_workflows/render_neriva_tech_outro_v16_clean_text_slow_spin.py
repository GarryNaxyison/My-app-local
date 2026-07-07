from __future__ import annotations

import json
import math

from PIL import Image, ImageDraw, ImageFilter, ImageOps

import render_poliglot_tech_outro_v15_clean_vector_text as v15
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
    draw_logo_shards,
    draw_orbits,
    ease_out_cubic,
    encode_video,
    smoothstep,
)
from render_poliglot_tech_outro_v12_comet_text import decode_source, make_contact_sheet, source_background


SOURCE_VIDEO = COMFY_OUTPUT / "poliglot_outro_2s_tech_warp_v11_api_00001_.mp4"
VARIANTS = (
    {
        "name": "lower",
        "brand_text": "neriva.ru",
        "prefix": "neriva.",
        "accent": "ru",
        "output": COMFY_OUTPUT / "neriva_outro_tech_warp_v16_lower_clean_text_slow_spin.mp4",
        "preview_stem": "neriva_outro_tech_warp_v16_lower",
        "font_size": 148,
        "flight_font_size": 96,
    },
    {
        "name": "upper",
        "brand_text": "NERIVA.RU",
        "prefix": "NERIVA.",
        "accent": "RU",
        "output": COMFY_OUTPUT / "neriva_outro_tech_warp_v16_upper_clean_text_slow_spin.mp4",
        "preview_stem": "neriva_outro_tech_warp_v16_upper",
        "font_size": 138,
        "flight_font_size": 92,
    },
)


def draw_slow_coin_spin_logo(base: Image.Image, time: float, logo: Image.Image) -> None:
    spin_start = 0.38
    spin_end = 1.98
    if time < spin_start:
        draw_logo_shards(base, min(0.96, time * 2.45), logo)
        return

    center_x = WIDTH // 2
    center_y = 662
    size = 512
    spin = smoothstep(spin_start, spin_end, time)
    theta = spin * math.tau
    face = math.cos(theta)
    width_scale = 1.0 if spin >= 1 else max(0.16, abs(face))
    sprite = logo.resize((max(1, int(size * width_scale)), size), Image.Resampling.LANCZOS)
    if spin < 1 and face < 0:
        sprite = ImageOps.mirror(sprite)

    x = center_x - sprite.width // 2
    y = center_y - sprite.height // 2

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(glow)
    glow_pulse = 0.92 + 0.08 * math.sin(time * 7.0)
    draw.ellipse((center_x - 395, center_y - 360, center_x + 395, center_y + 405), fill=(64, 255, 215, int(44 * glow_pulse)))
    draw.ellipse((center_x - 255, center_y - 230, center_x + 255, center_y + 285), fill=(48, 102, 255, int(40 * glow_pulse)))
    base.alpha_composite(glow.filter(ImageFilter.GaussianBlur(58)))

    draw_orbits(base, center_x, center_y, time, 1.0)
    base.alpha_composite(sprite, (x, y))

    if spin < 1:
        glint = Image.new("RGBA", sprite.size, (0, 0, 0, 0))
        glint_draw = ImageDraw.Draw(glint)
        sx = int((spin * 1.12 - 0.08) * sprite.width)
        glint_draw.polygon(
            [(sx, -20), (sx + 64, -20), (sx + 18, sprite.height + 20), (sx - 46, sprite.height + 20)],
            fill=(255, 255, 255, int(66 * (1 - abs(spin - 0.5) * 0.55))),
        )
        alpha_mask = sprite.getchannel("A")
        glint.putalpha(Image.composite(glint.getchannel("A"), Image.new("L", sprite.size, 0), alpha_mask))
        base.alpha_composite(glint.filter(ImageFilter.GaussianBlur(1.0)), (x, y))


def render_variant(variant: dict[str, object], logo: Image.Image, source_frames: list[Image.Image]) -> dict:
    brand_text = str(variant["brand_text"])
    v15.TEXT = brand_text
    v15.TEXT_PREFIX = str(variant["prefix"])
    v15.TEXT_ACCENT = str(variant["accent"])

    final_font, final_font_path = v15.load_clean_font(int(variant["font_size"]))
    flight_font, _ = v15.load_clean_font(int(variant["flight_font_size"]))
    final_layer, final_mask, glyphs, text_box = v15.make_motion_url_assets(final_font, flight_font)

    frames: list[Image.Image] = []
    for index, source_frame in enumerate(source_frames):
        time = index / FPS
        base = source_background(source_frame, index)
        draw_slow_coin_spin_logo(base, time, logo)
        v15.draw_url_motion(base, final_layer, final_mask, glyphs, text_box, time)
        add_finish_sweep(base, time)
        add_vignette(base)
        frames.append(base.convert("RGB"))

    output_video = variant["output"]
    assert isinstance(output_video, type(COMFY_OUTPUT))
    encode_video(frames, output_video)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    preview_stem = str(variant["preview_stem"])
    first = PREVIEW_DIR / f"{preview_stem}_first.png"
    mid = PREVIEW_DIR / f"{preview_stem}_mid.png"
    last = PREVIEW_DIR / f"{preview_stem}_last.png"
    sheet = PREVIEW_DIR / f"{preview_stem}_contact_sheet.jpg"
    frames[0].save(first)
    frames[TOTAL_FRAMES // 2].save(mid)
    frames[-1].save(last)
    make_contact_sheet(frames, sheet)

    result = {
        "source": str(SOURCE_VIDEO),
        "name": str(variant["name"]),
        "video": str(output_video),
        "first_frame": str(first),
        "mid_frame": str(mid),
        "last_frame": str(last),
        "contact_sheet": str(sheet),
        "font": str(final_font_path),
        "brand_text": brand_text,
        "width": WIDTH,
        "height": HEIGHT,
        "fps": FPS,
        "frames": TOTAL_FRAMES,
        "duration": TOTAL_FRAMES / FPS,
    }
    return result


def render() -> dict:
    logo = Image.open(LOGO_PATH).convert("RGBA")
    source_frames = decode_source(SOURCE_VIDEO)
    variants = [render_variant(variant, logo, source_frames) for variant in VARIANTS]
    result = {"source": str(SOURCE_VIDEO), "variants": variants}
    (PREVIEW_DIR / "neriva_outro_tech_warp_v16_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
