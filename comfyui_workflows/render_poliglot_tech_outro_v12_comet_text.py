from __future__ import annotations

import json
import math
import random
from pathlib import Path

import av
from PIL import Image, ImageDraw, ImageEnhance, ImageFilter, ImageFont

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
    ease_out_back,
    ease_out_cubic,
    encode_video,
    load_font,
    render_background,
    smoothstep,
)


SOURCE_VIDEO = COMFY_OUTPUT / "poliglot_outro_2s_tech_warp_v11_api_00001_.mp4"
OUTPUT_VIDEO = COMFY_OUTPUT / "poliglot_outro_tech_warp_v12_comet_text.mp4"


def cover_resize(image: Image.Image, width: int, height: int) -> Image.Image:
    scale = max(width / image.width, height / image.height)
    resized = image.resize((math.ceil(image.width * scale), math.ceil(image.height * scale)), Image.Resampling.LANCZOS)
    left = (resized.width - width) // 2
    top = (resized.height - height) // 2
    return resized.crop((left, top, left + width, top + height))


def decode_source(path: Path) -> list[Image.Image]:
    if not path.exists():
        raise FileNotFoundError(f"Missing ComfyUI source video: {path}")
    container = av.open(str(path))
    stream = container.streams.video[0]
    frames = [cover_resize(frame.to_image().convert("RGB"), WIDTH, HEIGHT) for frame in container.decode(stream)]
    if not frames:
        raise RuntimeError(f"No frames decoded from {path}")
    while len(frames) < TOTAL_FRAMES:
        frames.append(frames[-1].copy())
    return frames[:TOTAL_FRAMES]


def build_veil() -> Image.Image:
    veil = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    pixels = veil.load()
    for y in range(HEIGHT):
        strength = int(38 + 92 * (y / HEIGHT))
        for x in range(WIDTH):
            pixels[x, y] = (0, 8, 16, strength)
    return veil


VEIL = build_veil()


def hide_model_text_artifacts(base: Image.Image) -> None:
    mask = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(mask)
    draw.rounded_rectangle((205, 650, 875, 1530), radius=160, fill=(0, 8, 16, 218))
    draw.ellipse((225, 540, 855, 1300), fill=(7, 72, 88, 78))
    draw.ellipse((120, 910, 960, 1640), fill=(0, 12, 18, 150))
    base.alpha_composite(mask.filter(ImageFilter.GaussianBlur(72)))


def source_background(frame: Image.Image, index: int) -> Image.Image:
    time = index / FPS
    source = frame.convert("RGBA").filter(ImageFilter.GaussianBlur(10))
    source = ImageEnhance.Brightness(source).enhance(0.42)
    source = ImageEnhance.Color(source).enhance(1.22)
    source = ImageEnhance.Contrast(source).enhance(1.12)
    hide_model_text_artifacts(source)

    procedural = render_background(time)
    source.alpha_composite(apply_alpha(procedural, 0.56))
    source.alpha_composite(VEIL)
    return source


def draw_brand_char(
    draw: ImageDraw.ImageDraw,
    xy: tuple[int, int],
    char: str,
    font: ImageFont.FreeTypeFont,
    *,
    accent: bool,
    alpha: float = 1.0,
) -> None:
    base_color = (248, 255, 253, int(255 * alpha))
    accent_color = (86, 255, 216, int(255 * alpha))
    shadow = (0, 18, 23, int(232 * alpha))
    draw.text(xy, char, font=font, fill=accent_color if accent else base_color, stroke_width=1, stroke_fill=shadow)


def make_url_assets(font: ImageFont.FreeTypeFont) -> tuple[Image.Image, list[dict[str, object]], tuple[int, int, int, int]]:
    probe = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    probe_draw = ImageDraw.Draw(probe)
    bbox = probe_draw.textbbox((0, 0), TEXT, font=font, stroke_width=1)
    text_width = int(math.ceil(probe_draw.textlength(TEXT, font=font)))
    text_height = bbox[3] - bbox[1]
    x = (WIDTH - text_width) // 2
    y = 1256

    plate = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    plate_draw = ImageDraw.Draw(plate)
    plate_draw.rounded_rectangle((x - 80, y - 28, x + text_width + 80, y + text_height + 52), radius=32, fill=(0, 13, 18, 118))
    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    glow_draw.rounded_rectangle((x - 86, y - 32, x + text_width + 86, y + text_height + 56), radius=34, fill=(72, 255, 218, 20))
    plate.alpha_composite(glow.filter(ImageFilter.GaussianBlur(14)))
    plate = plate.filter(ImageFilter.GaussianBlur(3.2))

    glyphs: list[dict[str, object]] = []
    accent_start = len(TEXT_PREFIX)
    accent_end = accent_start + len(TEXT_ACCENT)
    for index, char in enumerate(TEXT):
        char_x = x + int(round(probe_draw.textlength(TEXT[:index], font=font)))
        glyph = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        glyph_draw = ImageDraw.Draw(glyph)
        accent = accent_start <= index < accent_end
        glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        glow_draw = ImageDraw.Draw(glow)
        draw_brand_char(glow_draw, (char_x, y), char, font, accent=accent, alpha=0.48)
        glyph.alpha_composite(glow.filter(ImageFilter.GaussianBlur(4)))
        draw_brand_char(glyph_draw, (char_x, y), char, font, accent=accent)
        rng = random.Random(84100 + index)
        glyphs.append(
            {
                "image": glyph,
                "index": index,
                "start_dx": rng.uniform(-620, -210) - index * 10,
                "start_dy": rng.uniform(-210, 80) - math.sin(index / max(1, len(TEXT) - 1) * math.pi) * 70,
                "spark": rng.uniform(0, 1),
            }
        )

    return plate, glyphs, (x, y, x + text_width, y + text_height)


def draw_comet_head(base: Image.Image, text_box: tuple[int, int, int, int], time: float) -> None:
    raw = smoothstep(0.68, 1.30, time)
    if raw <= 0 or time > 1.54:
        return
    progress = ease_out_cubic(raw)
    head_x = int(-240 + progress * (text_box[2] + 285))
    head_y = int((text_box[1] + text_box[3]) / 2 - 135 * math.sin(progress * math.pi) - 28 * (1 - progress))
    tail = 1 - smoothstep(1.22, 1.54, time)

    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    for index in range(15):
        length = 190 + index * 54
        width = 10 + index * 5
        alpha = int(104 * tail * (1 - index / 16))
        draw.polygon(
            [
                (head_x - 10 - index * 2, head_y - 17),
                (head_x - length, head_y - width),
                (head_x - length - 45, head_y + width * 0.2),
                (head_x - 10 - index * 2, head_y + 17),
            ],
            fill=(72, 255, 220, alpha),
        )
    draw.ellipse((head_x - 46, head_y - 46, head_x + 46, head_y + 46), fill=(174, 255, 237, int(118 * tail)))
    draw.ellipse((head_x - 19, head_y - 19, head_x + 19, head_y + 19), fill=(255, 255, 255, int(210 * tail)))
    base.alpha_composite(layer.filter(ImageFilter.GaussianBlur(8)))


def draw_url_comet_assemble(
    base: Image.Image,
    plate: Image.Image,
    glyphs: list[dict[str, object]],
    text_box: tuple[int, int, int, int],
    time: float,
) -> None:
    plate_alpha = smoothstep(0.94, 1.34, time)
    if plate_alpha > 0:
        base.alpha_composite(apply_alpha(plate, plate_alpha))

    draw_comet_head(base, text_box, time)

    particle_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    particle_draw = ImageDraw.Draw(particle_layer)

    for glyph in glyphs:
        index = int(glyph["index"])
        delay = 0.80 + index * 0.026
        local = smoothstep(delay, delay + 0.46, time)
        if local <= 0:
            continue
        settle = ease_out_back(local)
        dx = int(float(glyph["start_dx"]) * (1 - ease_out_cubic(local)))
        dy = int(float(glyph["start_dy"]) * (1 - settle))
        char_layer = glyph["image"]
        assert isinstance(char_layer, Image.Image)
        if local < 0.82:
            char_layer = char_layer.filter(ImageFilter.GaussianBlur(4.2 * (1 - local)))
        moved = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        moved.alpha_composite(apply_alpha(char_layer, clamp(local * 1.25)), (dx, dy))
        base.alpha_composite(moved)

        if local < 0.98:
            rng = random.Random(99000 + index)
            for _ in range(5):
                px = text_box[0] + rng.random() * (text_box[2] - text_box[0]) + dx * rng.uniform(0.1, 0.75)
                py = text_box[1] + rng.uniform(-28, 72) + dy * rng.uniform(0.1, 0.75)
                size = rng.uniform(1.2, 4.3)
                alpha = int(125 * local * (1 - local * 0.55))
                particle_draw.ellipse((px - size, py - size, px + size, py + size), fill=(130, 255, 231, alpha))

    base.alpha_composite(particle_layer.filter(ImageFilter.GaussianBlur(0.45)))

    shock = smoothstep(1.36, 1.58, time) * (1 - smoothstep(1.72, 1.96, time))
    if shock > 0:
        layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        draw = ImageDraw.Draw(layer)
        center_y = int((text_box[1] + text_box[3]) / 2)
        draw.rounded_rectangle((text_box[0] - 72, center_y - 5, text_box[2] + 72, center_y + 5), radius=5, fill=(101, 255, 222, int(96 * shock)))
        for index in range(18):
            rng = random.Random(6710 + index)
            sx = rng.uniform(text_box[0] - 40, text_box[2] + 40)
            sy = rng.uniform(text_box[1] - 42, text_box[3] + 60)
            size = rng.uniform(1.4, 4.6)
            draw.ellipse((sx - size, sy - size, sx + size, sy + size), fill=(196, 255, 240, int(132 * shock)))
        base.alpha_composite(layer.filter(ImageFilter.GaussianBlur(3)))


def make_contact_sheet(frames: list[Image.Image], path: Path) -> None:
    picks = [0, 10, 20, 32, 44, 56, 63]
    thumb_w = 216
    thumb_h = 384
    sheet = Image.new("RGB", (thumb_w * len(picks), thumb_h), (0, 0, 0))
    for col, frame_index in enumerate(picks):
        thumb = frames[frame_index].resize((thumb_w, thumb_h), Image.Resampling.LANCZOS)
        sheet.paste(thumb, (col * thumb_w, 0))
    sheet.save(path, quality=92)


def render() -> dict:
    logo = Image.open(LOGO_PATH).convert("RGBA")
    font = load_font(92, bold=True)
    plate, glyphs, text_box = make_url_assets(font)
    source_frames = decode_source(SOURCE_VIDEO)

    frames: list[Image.Image] = []
    for index, source_frame in enumerate(source_frames):
        time = index / FPS
        base = source_background(source_frame, index)
        draw_logo_shards(base, time, logo)
        draw_url_comet_assemble(base, plate, glyphs, text_box, time)
        add_finish_sweep(base, time)
        add_vignette(base)
        frames.append(base.convert("RGB"))

    encode_video(frames, OUTPUT_VIDEO)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / "outro_tech_warp_v12_first.png"
    mid = PREVIEW_DIR / "outro_tech_warp_v12_mid.png"
    last = PREVIEW_DIR / "outro_tech_warp_v12_last.png"
    sheet = PREVIEW_DIR / "outro_tech_warp_v12_contact_sheet.jpg"
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
    (PREVIEW_DIR / "poliglot_outro_tech_warp_v12_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
