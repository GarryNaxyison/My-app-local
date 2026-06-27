from __future__ import annotations

import json
import math
import random
from pathlib import Path

import av
from PIL import Image, ImageDraw, ImageFilter, ImageFont


PROJECT_ROOT = Path(__file__).resolve().parents[1]
COMFY_OUTPUT = Path(r"C:\Users\Admin\Documents\ComfyUI\output\outro")
SOURCE_CANDIDATES = [
    COMFY_OUTPUT / "poliglot_outro_2s_api_00001_.mp4",
    COMFY_OUTPUT / "poliglot_outro_2s_api_00002_.mp4",
]
OUTPUT_VIDEO = COMFY_OUTPUT / "poliglot_outro_ae_reveal_v2.mp4"
PREVIEW_DIR = PROJECT_ROOT / "tmp" / "comfy-outro"
LOGO_PATH = PROJECT_ROOT / "web" / "assets" / "brand-logo.png"

WIDTH = 1080
HEIGHT = 1920
FPS = 32
DURATION = 2.0
TOTAL_FRAMES = int(FPS * DURATION)
TEXT = "poliglotAI.online"


def smoothstep(edge0: float, edge1: float, value: float) -> float:
    if edge0 == edge1:
        return 1.0 if value >= edge1 else 0.0
    x = max(0.0, min(1.0, (value - edge0) / (edge1 - edge0)))
    return x * x * (3.0 - 2.0 * x)


def ease_out_back(value: float) -> float:
    x = max(0.0, min(1.0, value))
    c1 = 1.70158
    c3 = c1 + 1
    return 1 + c3 * pow(x - 1, 3) + c1 * pow(x - 1, 2)


def cover_resize(image: Image.Image, width: int, height: int) -> Image.Image:
    scale = max(width / image.width, height / image.height)
    resized = image.resize((math.ceil(image.width * scale), math.ceil(image.height * scale)), Image.Resampling.LANCZOS)
    left = (resized.width - width) // 2
    top = (resized.height - height) // 2
    return resized.crop((left, top, left + width, top + height))


def load_font(size: int, bold: bool = False) -> ImageFont.FreeTypeFont:
    font_name = "segoeuib.ttf" if bold else "segoeui.ttf"
    font_path = Path("C:/Windows/Fonts") / font_name
    if font_path.exists():
        return ImageFont.truetype(str(font_path), size=size)
    return ImageFont.load_default(size=size)


def pick_source_video() -> Path:
    for candidate in SOURCE_CANDIDATES:
        if candidate.exists():
            return candidate
    raise FileNotFoundError("No ComfyUI outro source MP4 found. Render the ComfyUI outro first.")


def decode_source_frames(path: Path) -> list[Image.Image]:
    container = av.open(str(path))
    stream = container.streams.video[0]
    frames: list[Image.Image] = []
    for frame in container.decode(stream):
        frames.append(cover_resize(frame.to_image().convert("RGB"), WIDTH, HEIGHT))
        if len(frames) >= TOTAL_FRAMES:
            break
    if not frames:
        raise RuntimeError(f"No frames decoded from {path}")
    while len(frames) < TOTAL_FRAMES:
        frames.append(frames[-1].copy())
    return frames[:TOTAL_FRAMES]


def gradient_band() -> Image.Image:
    band = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(band)

    top_feather = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    top_pixels = top_feather.load()
    for y in range(930, 1190):
        alpha = int(250 * smoothstep(930, 1190, y))
        for x in range(WIDTH):
            top_pixels[x, y] = (0, 6, 10, alpha)
    band.alpha_composite(top_feather.filter(ImageFilter.GaussianBlur(18)))

    draw.rectangle((0, 1120, WIDTH, 1648), fill=(0, 6, 10, 255))

    bottom_feather = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    bottom_pixels = bottom_feather.load()
    for y in range(1600, 1765):
        alpha = int(240 * (1.0 - smoothstep(1600, 1765, y)))
        for x in range(WIDTH):
            bottom_pixels[x, y] = (0, 6, 10, alpha)
    band.alpha_composite(bottom_feather.filter(ImageFilter.GaussianBlur(16)))

    for y in range(1120, 1648, 28):
        draw.line((90, y, WIDTH - 90, y), fill=(32, 255, 214, 9), width=1)
    draw.line((130, 1222, WIDTH - 130, 1222), fill=(105, 255, 220, 34), width=2)
    draw.line((170, 1518, WIDTH - 170, 1518), fill=(105, 255, 220, 26), width=2)
    return band


def draw_logo(base: Image.Image, logo: Image.Image, time: float) -> None:
    progress = ease_out_back(smoothstep(0.05, 0.52, time))
    size = int(360 + 235 * progress)
    logo_frame = logo.resize((size, size), Image.Resampling.LANCZOS)
    x = (WIDTH - size) // 2
    y = int(520 - 72 * progress)

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    glow_draw.ellipse((x - 80, y - 55, x + size + 80, y + size + 95), fill=(92, 245, 205, int(42 + 35 * progress)))
    glow_draw.ellipse((x + 20, y + 20, x + size - 20, y + size - 20), fill=(36, 94, 255, int(30 + 22 * progress)))
    glow = glow.filter(ImageFilter.GaussianBlur(55))
    base.alpha_composite(glow)

    shadow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    shadow.alpha_composite(logo_frame, (x, y))
    shadow = shadow.filter(ImageFilter.GaussianBlur(20))
    base.alpha_composite(shadow)
    base.alpha_composite(logo_frame, (x, y))

    sweep = smoothstep(0.48, 0.92, time) * (1 - smoothstep(0.92, 1.22, time))
    if sweep > 0:
        draw = ImageDraw.Draw(base)
        sx = int((time - 0.48) / 0.74 * (WIDTH + 280) - 140)
        draw.line((sx, y - 40, sx + 260, y + size + 70), fill=(235, 255, 250, int(135 * sweep)), width=12)


def text_layer(font: ImageFont.FreeTypeFont) -> tuple[Image.Image, tuple[int, int, int, int]]:
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    bbox = draw.textbbox((0, 0), TEXT, font=font)
    text_width = bbox[2] - bbox[0]
    text_height = bbox[3] - bbox[1]
    x = (WIDTH - text_width) // 2
    y = 1308

    shadow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    shadow_draw = ImageDraw.Draw(shadow)
    shadow_draw.text((x, y), TEXT, font=font, fill=(0, 0, 0, 230))
    shadow = shadow.filter(ImageFilter.GaussianBlur(13))
    layer.alpha_composite(shadow)

    draw = ImageDraw.Draw(layer)
    draw.text((x, y), TEXT, font=font, fill=(239, 255, 250, 255))
    draw.text((x, y + 5), TEXT, font=font, fill=(103, 255, 214, 96))
    return layer, (x, y, x + text_width, y + text_height)


def draw_text_reveal(base: Image.Image, prepared_text: Image.Image, text_box: tuple[int, int, int, int], time: float) -> None:
    progress = smoothstep(0.7, 1.24, time)
    if progress <= 0:
        return

    reveal_width = int((text_box[2] - text_box[0]) * min(1.0, progress * 1.08))
    mask = Image.new("L", (WIDTH, HEIGHT), 0)
    mask_draw = ImageDraw.Draw(mask)
    mask_draw.rounded_rectangle((text_box[0] - 28, text_box[1] - 32, text_box[0] + reveal_width + 28, text_box[3] + 34), radius=20, fill=int(255 * progress))
    mask = mask.filter(ImageFilter.GaussianBlur(max(0.2, 7 * (1 - progress))))

    text = prepared_text.copy()
    if progress < 0.92:
        text = text.filter(ImageFilter.GaussianBlur(5 * (1 - progress)))

    y_offset = int((1 - ease_out_back(progress)) * 145)
    x_offset = int((1 - progress) * -18)
    moved = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    moved.alpha_composite(text, (x_offset, y_offset))
    moved.putalpha(Image.composite(moved.getchannel("A"), Image.new("L", (WIDTH, HEIGHT), 0), mask))
    base.alpha_composite(moved)

    sweep = smoothstep(1.12, 1.46, time) * (1 - smoothstep(1.46, 1.72, time))
    if sweep > 0:
        sweep_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        draw = ImageDraw.Draw(sweep_layer)
        sx = int(text_box[0] - 180 + ((time - 1.08) / 0.60) * ((text_box[2] - text_box[0]) + 360))
        draw.polygon(
            [(sx, text_box[1] - 42), (sx + 95, text_box[1] - 42), (sx + 25, text_box[3] + 70), (sx - 70, text_box[3] + 70)],
            fill=(255, 255, 255, int(118 * sweep)),
        )
        sweep_layer.putalpha(Image.composite(sweep_layer.getchannel("A"), Image.new("L", (WIDTH, HEIGHT), 0), mask))
        base.alpha_composite(sweep_layer)

    line_progress = smoothstep(0.98, 1.42, time)
    if line_progress > 0:
        line_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        draw = ImageDraw.Draw(line_layer)
        center = WIDTH // 2
        half = int((text_box[2] - text_box[0]) * 0.44 * line_progress)
        y = text_box[3] + 34 + y_offset
        draw.line((center - half, y, center + half, y), fill=(119, 255, 216, int(155 * line_progress)), width=5)
        draw.line((center - half, y + 10, center + half, y + 10), fill=(34, 122, 102, int(70 * line_progress)), width=2)
        base.alpha_composite(line_layer.filter(ImageFilter.GaussianBlur(0.6)))


def draw_particles(base: Image.Image, time: float) -> None:
    rng = random.Random(90210)
    draw = ImageDraw.Draw(base)
    appear = smoothstep(0.25, 0.9, time)
    for _ in range(34):
        start = rng.uniform(0.15, 1.55)
        life = rng.uniform(0.45, 0.9)
        p = smoothstep(start, start + life * 0.45, time) * (1 - smoothstep(start + life * 0.55, start + life, time))
        if p <= 0:
            continue
        angle = rng.uniform(0, math.tau)
        radius = rng.uniform(230, 490) + time * rng.uniform(12, 38)
        x = WIDTH / 2 + math.cos(angle) * radius
        y = 810 + math.sin(angle) * radius * 0.48 - time * rng.uniform(10, 32)
        size = rng.uniform(2.0, 5.5)
        alpha = int(210 * p * appear)
        draw.ellipse((x - size, y - size, x + size, y + size), fill=(190, 255, 235, alpha))


def add_vignette(base: Image.Image) -> None:
    vignette = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(vignette)
    for i in range(90):
        alpha = int(i * 1.45)
        draw.rectangle((i, i, WIDTH - i, HEIGHT - i), outline=(0, 0, 0, max(0, 130 - alpha)), width=2)
    base.alpha_composite(vignette.filter(ImageFilter.GaussianBlur(16)))


def encode_video(frames: list[Image.Image], output_path: Path) -> None:
    output_path.parent.mkdir(parents=True, exist_ok=True)
    container = av.open(str(output_path), "w")
    stream = container.add_stream("libx264", rate=FPS)
    stream.width = WIDTH
    stream.height = HEIGHT
    stream.pix_fmt = "yuv420p"
    stream.options = {"crf": "18", "preset": "medium"}

    for image in frames:
        frame = av.VideoFrame.from_image(image.convert("RGB"))
        for packet in stream.encode(frame):
            container.mux(packet)
    for packet in stream.encode():
        container.mux(packet)
    container.close()


def render() -> dict:
    source = pick_source_video()
    source_frames = decode_source_frames(source)
    logo = Image.open(LOGO_PATH).convert("RGBA")
    font = load_font(92, bold=True)
    prepared_text, text_box = text_layer(font)
    band = gradient_band()

    rendered: list[Image.Image] = []
    for index, frame in enumerate(source_frames):
        time = index / FPS
        base = frame.convert("RGBA")
        base.alpha_composite(band)
        draw_logo(base, logo, time)
        draw_particles(base, time)
        draw_text_reveal(base, prepared_text, text_box, time)
        add_vignette(base)
        rendered.append(base.convert("RGB"))

    encode_video(rendered, OUTPUT_VIDEO)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / "outro_ae_reveal_v2_first.png"
    last = PREVIEW_DIR / "outro_ae_reveal_v2_last.png"
    rendered[0].save(first)
    rendered[-1].save(last)

    result = {
        "source": str(source),
        "video": str(OUTPUT_VIDEO),
        "first_frame": str(first),
        "last_frame": str(last),
        "width": WIDTH,
        "height": HEIGHT,
        "fps": FPS,
        "frames": TOTAL_FRAMES,
        "duration": DURATION,
    }
    (PREVIEW_DIR / "poliglot_outro_ae_reveal_v2_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
