from __future__ import annotations

import json
import math
from pathlib import Path

import av
from PIL import Image, ImageEnhance, ImageFilter

from render_poliglot_tech_outro_v11 import (
    COMFY_OUTPUT,
    FPS,
    HEIGHT,
    LOGO_PATH,
    MANROPE_FONT_PATH,
    PREVIEW_DIR,
    TOTAL_FRAMES,
    WIDTH,
    add_finish_sweep,
    add_vignette,
    apply_alpha,
    draw_logo_shards,
    draw_url_comet,
    encode_video,
    load_font,
    make_text_layer,
    render_background,
)


SOURCE_VIDEO = COMFY_OUTPUT / "poliglot_outro_2s_tech_warp_v11_api_00001_.mp4"
OUTPUT_VIDEO = COMFY_OUTPUT / "poliglot_outro_tech_warp_v11_comfy_composite.mp4"


def build_veil() -> Image.Image:
    veil = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    veil_pixels = veil.load()
    for y in range(HEIGHT):
        strength = int(26 + 76 * (y / HEIGHT))
        for x in range(WIDTH):
            veil_pixels[x, y] = (0, 8, 16, strength)
    return veil


VEIL = build_veil()


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


def source_background(frame: Image.Image, index: int) -> Image.Image:
    time = index / FPS
    source = frame.convert("RGBA").filter(ImageFilter.GaussianBlur(6))
    source = ImageEnhance.Brightness(source).enhance(0.52)
    source = ImageEnhance.Color(source).enhance(1.28)
    source = ImageEnhance.Contrast(source).enhance(1.18)

    procedural = render_background(time)
    source.alpha_composite(apply_alpha(procedural, 0.34))

    source.alpha_composite(VEIL)
    return source


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
    prepared_text, text_box = make_text_layer(font)
    source_frames = decode_source(SOURCE_VIDEO)

    frames: list[Image.Image] = []
    for index, source_frame in enumerate(source_frames):
        time = index / FPS
        base = source_background(source_frame, index)
        draw_logo_shards(base, time, logo)
        draw_url_comet(base, prepared_text, text_box, time)
        add_finish_sweep(base, time)
        add_vignette(base)
        frames.append(base.convert("RGB"))

    encode_video(frames, OUTPUT_VIDEO)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / "outro_tech_warp_v11_composite_first.png"
    mid = PREVIEW_DIR / "outro_tech_warp_v11_composite_mid.png"
    last = PREVIEW_DIR / "outro_tech_warp_v11_composite_last.png"
    sheet = PREVIEW_DIR / "outro_tech_warp_v11_composite_contact_sheet.jpg"
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
    (PREVIEW_DIR / "poliglot_outro_tech_warp_v11_composite_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
