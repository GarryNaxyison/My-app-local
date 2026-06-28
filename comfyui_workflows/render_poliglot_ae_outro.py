from __future__ import annotations

import json
import math
from pathlib import Path

import av
from PIL import Image, ImageDraw, ImageFilter, ImageFont


PROJECT_ROOT = Path(__file__).resolve().parents[1]
COMFY_OUTPUT = Path(r"C:\Users\Admin\Documents\ComfyUI\output\outro")
SOURCE_CANDIDATES = [
    COMFY_OUTPUT / "poliglot_outro_2s_api_00001_.mp4",
    COMFY_OUTPUT / "poliglot_outro_2s_api_00002_.mp4",
]
OUTPUT_VIDEO = COMFY_OUTPUT / "poliglot_outro_ae_reveal_v9_clean.mp4"
PREVIEW_DIR = PROJECT_ROOT / "tmp" / "comfy-outro"
MANROPE_FONT_PATH = PROJECT_ROOT / "comfyui_workflows" / "assets" / "fonts" / "Manrope-wght.ttf"

WIDTH = 1080
HEIGHT = 1920
FPS = 32
DURATION = 2.0
TOTAL_FRAMES = int(FPS * DURATION)
TEXT = "poliglotAI.online"
TEXT_PREFIX = "poliglot"
TEXT_ACCENT = "AI"


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
    if MANROPE_FONT_PATH.exists():
        font = ImageFont.truetype(str(MANROPE_FONT_PATH), size=size)
        try:
            font.set_variation_by_name("ExtraBold" if bold else "Medium")
        except OSError:
            pass
        return font

    names = ["segoeui.ttf", "corbel.ttf", "trebuc.ttf", "arial.ttf"]
    if bold:
        names = ["seguisb.ttf", "segoeuib.ttf", "corbelb.ttf", "trebucbd.ttf", "arialbd.ttf"]
    for name in names:
        font_path = Path("C:/Windows/Fonts") / name
        if font_path.exists():
            return ImageFont.truetype(str(font_path), size=size)
    return ImageFont.load_default(size=size)


def draw_brand_text(
    draw: ImageDraw.ImageDraw,
    xy: tuple[int, int],
    font: ImageFont.FreeTypeFont,
    *,
    alpha: float = 1.0,
    stroke_width: int = 0,
    stroke_fill: tuple[int, int, int, int] | None = None,
    base_color: tuple[int, int, int, int] = (245, 255, 252, 255),
    accent_color: tuple[int, int, int, int] = (112, 255, 218, 255),
) -> None:
    x, y = xy
    base_fill = (base_color[0], base_color[1], base_color[2], int(base_color[3] * alpha))
    accent_fill = (accent_color[0], accent_color[1], accent_color[2], int(accent_color[3] * alpha))
    draw.text((x, y), TEXT, font=font, fill=base_fill, stroke_width=stroke_width, stroke_fill=stroke_fill)

    accent_x = x + int(round(draw.textlength(TEXT_PREFIX, font=font)))
    draw.text(
        (accent_x, y),
        TEXT_ACCENT,
        font=font,
        fill=accent_fill,
        stroke_width=stroke_width,
        stroke_fill=stroke_fill,
    )


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

    return band


def text_layer(font: ImageFont.FreeTypeFont) -> tuple[Image.Image, tuple[int, int, int, int]]:
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    bbox = draw.textbbox((0, 0), TEXT, font=font, stroke_width=2)
    text_width = int(math.ceil(draw.textlength(TEXT, font=font)))
    text_height = bbox[3] - bbox[1]
    x = (WIDTH - text_width) // 2
    y = 1300

    pedestal = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    pedestal_draw = ImageDraw.Draw(pedestal)
    pedestal_draw.rounded_rectangle(
        (x - 58, y - 18, x + text_width + 58, y + text_height + 44),
        radius=22,
        fill=(0, 16, 20, 126),
    )
    layer.alpha_composite(pedestal.filter(ImageFilter.GaussianBlur(5)))

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    draw_brand_text(glow_draw, (x, y), font, alpha=0.42, stroke_width=1, stroke_fill=(92, 255, 218, 62))
    layer.alpha_composite(glow.filter(ImageFilter.GaussianBlur(4)))

    shadow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    shadow_draw = ImageDraw.Draw(shadow)
    draw_brand_text(
        shadow_draw,
        (x, y + 6),
        font,
        stroke_width=2,
        stroke_fill=(0, 0, 0, 240),
        base_color=(0, 0, 0, 246),
        accent_color=(0, 0, 0, 246),
    )
    shadow = shadow.filter(ImageFilter.GaussianBlur(7))
    layer.alpha_composite(shadow)

    draw = ImageDraw.Draw(layer)
    draw_brand_text(
        draw,
        (x, y),
        font,
        stroke_width=1,
        stroke_fill=(0, 24, 28, 210),
        base_color=(250, 255, 253, 255),
        accent_color=(97, 255, 216, 255),
    )
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
    font = load_font(94, bold=True)
    prepared_text, text_box = text_layer(font)
    band = gradient_band()

    rendered: list[Image.Image] = []
    for index, frame in enumerate(source_frames):
        time = index / FPS
        base = frame.convert("RGBA")
        base.alpha_composite(band)
        draw_text_reveal(base, prepared_text, text_box, time)
        add_vignette(base)
        rendered.append(base.convert("RGB"))

    encode_video(rendered, OUTPUT_VIDEO)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / "outro_ae_reveal_v9_clean_first.png"
    last = PREVIEW_DIR / "outro_ae_reveal_v9_clean_last.png"
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
    (PREVIEW_DIR / "poliglot_outro_ae_reveal_v9_clean_result.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
