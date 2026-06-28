from __future__ import annotations

import json
import math
import random
import shutil
from pathlib import Path

import av
from PIL import Image, ImageDraw, ImageFilter, ImageFont


PROJECT_ROOT = Path(__file__).resolve().parents[1]
COMFY_INPUT = Path(r"C:\Users\Admin\Documents\ComfyUI\input")
COMFY_OUTPUT = Path(r"C:\Users\Admin\Documents\ComfyUI\output\outro")
PREVIEW_DIR = PROJECT_ROOT / "tmp" / "comfy-outro"
LOGO_PATH = PROJECT_ROOT / "web" / "assets" / "brand-logo.png"
MANROPE_FONT_PATH = PROJECT_ROOT / "comfyui_workflows" / "assets" / "fonts" / "Manrope-wght.ttf"

BASE_API_PROMPT = PROJECT_ROOT / "comfyui_workflows" / "poliglot_outro_2s_wan_stage_api.json"
API_PROMPT = PROJECT_ROOT / "comfyui_workflows" / "poliglot_outro_2s_tech_warp_v11_api.json"
PROMPT_DOC = PROJECT_ROOT / "comfyui_workflows" / "poliglot_outro_2s_tech_warp_v11_prompt.md"
REFERENCE_NAME = "poliglot_outro_tech_warp_v11_reference_1080x1920.png"
REFERENCE_PATH = PROJECT_ROOT / "comfyui_workflows" / REFERENCE_NAME
OUTPUT_VIDEO = COMFY_OUTPUT / "poliglot_outro_tech_warp_v11.mp4"

WIDTH = 1080
HEIGHT = 1920
FPS = 32
DURATION = 2.0
TOTAL_FRAMES = int(FPS * DURATION)
TEXT = "poliglotAI.online"
TEXT_PREFIX = "poliglot"
TEXT_ACCENT = "AI"

WAN_PROMPT = """2-second vertical Poliglot AI technology outro for TikTok/Reels/Shorts.
Motion design style: high-speed sapphire and mint holographic data tunnel, parallax depth, camera push-in, logo assembling from luminous shards, glass glints, orbit ribbons, clean cinematic bloom, premium AI SaaS finish.
The URL moment is a light-object arrival: a bright cyan data comet streaks in and resolves into the exact readable text "poliglotAI.online" near the end.
Keep the Poliglot AI logo and URL stable, premium, sharp, centered, and readable. No people, no app UI, no captions, no extra words, no watermark."""

NEGATIVE_PROMPT = """flat static logo, primitive empty space, boring starfield, cheap neon, low detail, blurry text, misspelled URL, wrong URL, extra text, subtitles, captions, watermark, duplicated logo, distorted logo, deformed letters, random letters, unreadable typography, people, hands, faces, flags, app screenshots, UI panels, bad composition, black frames, flicker, jpeg artifacts"""


def clamp(value: float, low: float = 0.0, high: float = 1.0) -> float:
    return max(low, min(high, value))


def smoothstep(edge0: float, edge1: float, value: float) -> float:
    if edge0 == edge1:
        return 1.0 if value >= edge1 else 0.0
    x = clamp((value - edge0) / (edge1 - edge0))
    return x * x * (3.0 - 2.0 * x)


def ease_out_cubic(value: float) -> float:
    x = clamp(value)
    return 1 - pow(1 - x, 3)


def ease_out_back(value: float) -> float:
    x = clamp(value)
    c1 = 1.70158
    c3 = c1 + 1
    return 1 + c3 * pow(x - 1, 3) + c1 * pow(x - 1, 2)


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
    base_color: tuple[int, int, int, int] = (246, 255, 253, 255),
    accent_color: tuple[int, int, int, int] = (92, 255, 218, 255),
) -> None:
    x, y = xy
    base_fill = (base_color[0], base_color[1], base_color[2], int(base_color[3] * alpha))
    accent_fill = (accent_color[0], accent_color[1], accent_color[2], int(accent_color[3] * alpha))
    draw.text((x, y), TEXT, font=font, fill=base_fill, stroke_width=stroke_width, stroke_fill=stroke_fill)
    accent_x = x + int(round(draw.textlength(TEXT_PREFIX, font=font)))
    draw.text((accent_x, y), TEXT_ACCENT, font=font, fill=accent_fill, stroke_width=stroke_width, stroke_fill=stroke_fill)


def make_particles(count: int = 520) -> list[dict[str, float]]:
    rng = random.Random(91142)
    particles: list[dict[str, float]] = []
    for _ in range(count):
        angle = rng.uniform(0, math.tau)
        radius = pow(rng.uniform(0.08, 1.0), 0.62)
        particles.append(
            {
                "angle": angle,
                "radius": radius,
                "depth": rng.uniform(0.05, 1.0),
                "speed": rng.uniform(0.55, 1.8),
                "size": rng.uniform(0.8, 3.1),
                "hue": rng.uniform(0.0, 1.0),
                "twist": rng.uniform(-0.7, 0.7),
            }
        )
    return particles


def make_shards(count: int = 44) -> list[dict[str, float]]:
    rng = random.Random(20497)
    shards: list[dict[str, float]] = []
    for _ in range(count):
        shards.append(
            {
                "x": rng.uniform(-0.45, 0.45),
                "y": rng.uniform(-0.34, 0.34),
                "dx": rng.uniform(-580, 580),
                "dy": rng.uniform(-470, 420),
                "size": rng.uniform(8, 28),
                "rot": rng.uniform(0, math.tau),
                "delay": rng.uniform(0.0, 0.22),
            }
        )
    return shards


PARTICLES = make_particles()
SHARDS = make_shards()


def render_background(time: float) -> Image.Image:
    progress = time / DURATION
    small_w = WIDTH // 4
    small_h = HEIGHT // 4
    base = Image.new("RGB", (small_w, small_h))
    pixels = base.load()
    cx = 0.5 + math.sin(time * 1.7) * 0.018
    cy = 0.39 + math.cos(time * 1.2) * 0.018

    for y in range(small_h):
        ny = y / small_h - cy
        for x in range(small_w):
            nx = x / small_w - cx
            dist = math.sqrt(nx * nx * 1.25 + ny * ny * 0.72)
            angle = math.atan2(ny, nx)
            tunnel = max(0.0, 1.0 - dist * 1.72)
            ring = math.sin(dist * 58.0 - time * 12.0 + angle * 2.4) * 0.5 + 0.5
            scan = math.sin((x / small_w) * 22.0 + (y / small_h) * 12.0 - time * 9.0) * 0.5 + 0.5
            data = math.sin((angle + progress * 5.4) * 8.0 + dist * 16.0) * 0.5 + 0.5
            r = int(3 + 10 * tunnel + 18 * ring * tunnel + 24 * data * max(0.0, 0.9 - dist))
            g = int(8 + 38 * tunnel + 78 * ring * max(0.0, 0.75 - dist) + 35 * scan * max(0.0, 0.55 - dist))
            b = int(24 + 62 * tunnel + 58 * data * max(0.0, 0.85 - dist))
            if 0.32 < dist < 0.48:
                r += int(18 * ring)
                b += int(30 * ring)
            pixels[x, y] = (min(255, r), min(255, g), min(255, b))

    frame = base.resize((WIDTH, HEIGHT), Image.Resampling.BICUBIC).convert("RGBA")

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(glow)
    shift = math.sin(time * 2.1) * 80
    draw.ellipse((110 + shift, 220, 970 + shift, 1110), fill=(45, 245, 211, 46))
    draw.ellipse((-240 - shift * 0.3, 700, 520, 1690), fill=(86, 68, 210, 32))
    draw.ellipse((560 + shift * 0.4, 140, 1380, 980), fill=(34, 125, 255, 34))
    frame.alpha_composite(glow.filter(ImageFilter.GaussianBlur(105)))

    draw_tunnel_grid(frame, time)
    draw_particles(frame, time)
    draw_edge_circuits(frame, time)
    return frame


def draw_tunnel_grid(frame: Image.Image, time: float) -> None:
    grid = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(grid)
    cx = WIDTH / 2
    cy = 760
    pulse = time * 1.25

    for index in range(18):
        z = ((index / 18.0) + pulse * 0.42) % 1.0
        scale = 0.18 + z * z * 1.82
        w = int(270 * scale)
        h = int(470 * scale)
        alpha = int((1.0 - z) * 42 + 12)
        color = (68, 255, 222, max(0, min(70, alpha)))
        draw.ellipse((cx - w, cy - h, cx + w, cy + h), outline=color, width=max(1, int(5 - z * 3)))

    for index in range(28):
        angle = index / 28.0 * math.tau + time * 0.32
        length = 1390
        x = cx + math.cos(angle) * length
        y = cy + math.sin(angle) * length * 1.35
        alpha = 26 + int(18 * (math.sin(index + time * 4.0) * 0.5 + 0.5))
        draw.line((cx, cy, x, y), fill=(66, 218, 255, alpha), width=1)

    frame.alpha_composite(grid.filter(ImageFilter.GaussianBlur(0.35)))


def draw_particles(frame: Image.Image, time: float) -> None:
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    cx = WIDTH / 2
    cy = 760

    for particle in PARTICLES:
        depth = (particle["depth"] - time * particle["speed"] * 0.56) % 1.0
        prev_depth = min(1.0, depth + 0.045)
        twist = particle["twist"] * (1.0 - depth) + time * 0.62
        angle = particle["angle"] + twist
        scale = 1.0 / (depth * 1.65 + 0.045)
        prev_scale = 1.0 / (prev_depth * 1.65 + 0.045)
        x = cx + math.cos(angle) * particle["radius"] * scale * 510
        y = cy + math.sin(angle) * particle["radius"] * scale * 850
        px = cx + math.cos(angle - 0.028) * particle["radius"] * prev_scale * 510
        py = cy + math.sin(angle - 0.028) * particle["radius"] * prev_scale * 850
        if x < -120 or x > WIDTH + 120 or y < -140 or y > HEIGHT + 140:
            continue
        brightness = int(245 * (1 - depth) ** 0.68)
        if brightness <= 8:
            continue
        hue = particle["hue"]
        color = (
            int(88 + 92 * hue),
            int(224 + 30 * (1.0 - hue)),
            int(255 - 38 * hue),
            min(235, brightness),
        )
        width = max(1, int(particle["size"] * (1.9 - depth)))
        draw.line((px, py, x, y), fill=color, width=width)
        if depth < 0.22:
            s = particle["size"] * (2.4 - depth * 3)
            draw.ellipse((x - s, y - s, x + s, y + s), fill=(210, 255, 247, min(170, brightness)))

    frame.alpha_composite(layer.filter(ImageFilter.GaussianBlur(0.2)))


def draw_edge_circuits(frame: Image.Image, time: float) -> None:
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    rng = random.Random(5819)

    for side in (0, 1):
        base_x = 34 if side == 0 else WIDTH - 34
        direction = 1 if side == 0 else -1
        for index in range(15):
            y = 190 + index * 108 + math.sin(time * 2.0 + index) * 12
            length = rng.uniform(72, 180)
            alpha = int(18 + 30 * (math.sin(time * 4.0 + index * 0.7) * 0.5 + 0.5))
            x2 = base_x + direction * length
            draw.line((base_x, y, x2, y), fill=(90, 255, 221, alpha), width=2)
            draw.line((x2, y, x2 + direction * rng.uniform(12, 55), y + rng.choice([-1, 1]) * rng.uniform(18, 46)), fill=(90, 255, 221, alpha), width=2)
            if index % 3 == 0:
                draw.rounded_rectangle((min(base_x, x2) - 8, y - 8, max(base_x, x2) + 8, y + 8), radius=5, outline=(70, 190, 255, alpha), width=1)

    frame.alpha_composite(layer.filter(ImageFilter.GaussianBlur(0.25)))


def apply_alpha(image: Image.Image, alpha: float) -> Image.Image:
    if alpha >= 0.999:
        return image
    result = image.copy()
    result.putalpha(result.getchannel("A").point(lambda value: int(value * alpha)))
    return result


def tint_alpha(image: Image.Image, color: tuple[int, int, int], alpha: float) -> Image.Image:
    result = Image.new("RGBA", image.size, (color[0], color[1], color[2], 0))
    result.putalpha(image.getchannel("A").point(lambda value: int(value * alpha)))
    return result


def draw_logo_shards(base: Image.Image, time: float, logo: Image.Image) -> None:
    raw = smoothstep(0.06, 0.82, time)
    progress = ease_out_back(raw)
    final_size = 512
    start_size = 330
    size = int(start_size + (final_size - start_size) * clamp(progress))
    if raw > 0.85:
        size += int(math.sin((time - 0.82) * 18) * 7 * (1 - smoothstep(0.82, 1.2, time)))

    logo_frame = logo.resize((size, size), Image.Resampling.LANCZOS)
    center_x = WIDTH // 2
    center_y = 690
    x = center_x - size // 2
    y = int(center_y - size // 2 - 24 * smoothstep(0.9, 1.45, time))

    glow_strength = smoothstep(0.0, 0.62, time) * (1.0 + 0.4 * math.sin(time * 18))
    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(glow)
    draw.ellipse((x - 145, y - 130, x + size + 145, y + size + 150), fill=(73, 255, 215, int(42 * glow_strength)))
    draw.ellipse((x + 10, y + 30, x + size - 10, y + size + 60), fill=(49, 98, 255, int(44 * glow_strength)))
    base.alpha_composite(glow.filter(ImageFilter.GaussianBlur(58)))

    draw_orbits(base, center_x, y + size // 2, time, raw)

    if raw < 0.98:
        strip_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        strip_count = 12
        for index in range(strip_count):
            left = int(index * logo_frame.width / strip_count)
            right = int((index + 1) * logo_frame.width / strip_count)
            crop = logo_frame.crop((left, 0, right, logo_frame.height))
            direction = -1 if index % 2 == 0 else 1
            local = clamp((raw - index * 0.018) / 0.82)
            offset = (1 - ease_out_cubic(local)) ** 2
            dx = int(direction * (310 + index * 12) * offset + math.sin(time * 22 + index) * 9 * (1 - local))
            dy = int(((-150 if index % 3 == 0 else 180) + index * 8) * offset)
            crop = apply_alpha(crop, clamp(raw * 1.5))
            if local < 0.7:
                crop = crop.filter(ImageFilter.GaussianBlur(max(0.2, 3.8 * (1 - local))))
            strip_layer.alpha_composite(crop, (x + left + dx, y + dy))

        glitch = smoothstep(0.18, 0.58, time) * (1 - smoothstep(0.72, 0.92, time))
        if glitch > 0:
            alpha_mask = logo_frame.getchannel("A")
            cyan = tint_alpha(logo_frame, (75, 255, 220), 0.18 * glitch)
            blue = tint_alpha(logo_frame, (46, 105, 255), 0.15 * glitch)
            magenta = tint_alpha(logo_frame, (155, 82, 255), 0.12 * glitch)
            cyan.putalpha(alpha_mask.point(lambda value: int(value * 0.22 * glitch)))
            blue.putalpha(alpha_mask.point(lambda value: int(value * 0.18 * glitch)))
            magenta.putalpha(alpha_mask.point(lambda value: int(value * 0.14 * glitch)))
            strip_layer.alpha_composite(cyan, (x - 11, y + 4))
            strip_layer.alpha_composite(blue, (x + 9, y - 3))
            strip_layer.alpha_composite(magenta, (x + 2, y + 11))

        base.alpha_composite(strip_layer)
    else:
        base.alpha_composite(logo_frame, (x, y))

    settle = smoothstep(0.82, 1.16, time)
    if settle > 0:
        clean = apply_alpha(logo_frame, settle)
        base.alpha_composite(clean, (x, y))

    draw_assembly_shards(base, center_x, y + size // 2, time, raw)


def draw_orbits(base: Image.Image, cx: int, cy: int, time: float, logo_progress: float) -> None:
    orbit = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(orbit)
    alpha_base = int(105 * smoothstep(0.18, 0.95, time) * (1 - smoothstep(1.75, 2.0, time) * 0.35))
    for index, radius in enumerate((330, 386, 450)):
        start = int((time * (118 + index * 28) + index * 84) % 360)
        width = 3 if index != 1 else 4
        bbox = (cx - radius, cy - int(radius * 0.36), cx + radius, cy + int(radius * 0.36))
        draw.arc(bbox, start, start + 78 + index * 16, fill=(86, 255, 218, alpha_base), width=width)
        draw.arc(bbox, start + 180, start + 218, fill=(83, 142, 255, int(alpha_base * 0.72)), width=width)
    if logo_progress < 0.88:
        draw.ellipse((cx - 246, cy - 246, cx + 246, cy + 246), outline=(188, 255, 239, int(44 * logo_progress)), width=2)
    base.alpha_composite(orbit.filter(ImageFilter.GaussianBlur(0.45)))


def draw_assembly_shards(base: Image.Image, cx: int, cy: int, time: float, logo_progress: float) -> None:
    intensity = smoothstep(0.0, 0.42, time) * (1 - smoothstep(0.78, 1.12, time))
    if intensity <= 0:
        return
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    for shard in SHARDS:
        local = clamp((logo_progress - shard["delay"]) / (1.0 - shard["delay"]))
        x = cx + shard["x"] * 430 + shard["dx"] * (1 - ease_out_cubic(local))
        y = cy + shard["y"] * 430 + shard["dy"] * (1 - ease_out_cubic(local))
        s = shard["size"] * (0.45 + local)
        angle = shard["rot"] + time * 3.5
        points = []
        for offset in (0, 2.1, 4.2):
            points.append((x + math.cos(angle + offset) * s, y + math.sin(angle + offset) * s))
        alpha = int(138 * intensity * (1 - local * 0.35))
        draw.polygon(points, fill=(92, 255, 220, alpha))
        draw.line((x, y, cx + shard["x"] * 130, cy + shard["y"] * 130), fill=(76, 198, 255, int(alpha * 0.45)), width=1)
    base.alpha_composite(layer.filter(ImageFilter.GaussianBlur(0.35)))


def make_text_layer(font: ImageFont.FreeTypeFont) -> tuple[Image.Image, tuple[int, int, int, int]]:
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    bbox = draw.textbbox((0, 0), TEXT, font=font, stroke_width=1)
    text_width = int(math.ceil(draw.textlength(TEXT, font=font)))
    text_height = bbox[3] - bbox[1]
    x = (WIDTH - text_width) // 2
    y = 1256

    plate = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    plate_draw = ImageDraw.Draw(plate)
    plate_draw.rounded_rectangle((x - 64, y - 22, x + text_width + 64, y + text_height + 46), radius=24, fill=(0, 16, 20, 152))
    layer.alpha_composite(plate.filter(ImageFilter.GaussianBlur(4)))

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    draw_brand_text(glow_draw, (x, y), font, alpha=0.52, stroke_width=1, stroke_fill=(100, 255, 222, 80))
    layer.alpha_composite(glow.filter(ImageFilter.GaussianBlur(5)))

    shadow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    shadow_draw = ImageDraw.Draw(shadow)
    draw_brand_text(
        shadow_draw,
        (x, y + 7),
        font,
        stroke_width=3,
        stroke_fill=(0, 0, 0, 245),
        base_color=(0, 0, 0, 246),
        accent_color=(0, 0, 0, 246),
    )
    layer.alpha_composite(shadow.filter(ImageFilter.GaussianBlur(7)))

    draw = ImageDraw.Draw(layer)
    draw_brand_text(
        draw,
        (x, y),
        font,
        stroke_width=1,
        stroke_fill=(0, 31, 34, 230),
        base_color=(248, 255, 253, 255),
        accent_color=(86, 255, 216, 255),
    )
    return layer, (x, y, x + text_width, y + text_height)


def draw_url_comet(base: Image.Image, prepared_text: Image.Image, text_box: tuple[int, int, int, int], time: float) -> None:
    raw = smoothstep(0.74, 1.36, time)
    if raw <= 0:
        return
    progress = ease_out_cubic(raw)
    x_offset = int((1 - progress) * -910 + math.sin(raw * math.pi) * 18)
    y_offset = int((1 - progress) * -188 - math.sin(raw * math.pi) * 22)

    tail = smoothstep(0.74, 0.98, time) * (1 - smoothstep(1.23, 1.53, time))
    if tail > 0:
        tail_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        draw = ImageDraw.Draw(tail_layer)
        head_x = text_box[2] + x_offset - 12
        head_y = (text_box[1] + text_box[3]) / 2 + y_offset + 2
        for index in range(14):
            length = 230 + index * 64
            width = 14 + index * 6
            alpha = int(98 * tail * (1 - index / 15))
            draw.polygon(
                [
                    (head_x - 12 - index * 3, head_y - 18),
                    (head_x - length, head_y - width),
                    (head_x - length - 58, head_y + width * 0.18),
                    (head_x - 12 - index * 3, head_y + 18),
                ],
                fill=(74, 255, 220, alpha),
            )
        draw.ellipse((head_x - 42, head_y - 42, head_x + 42, head_y + 42), fill=(235, 255, 250, int(150 * tail)))
        draw.ellipse((head_x - 18, head_y - 18, head_x + 18, head_y + 18), fill=(255, 255, 255, int(195 * tail)))
        base.alpha_composite(tail_layer.filter(ImageFilter.GaussianBlur(9)))

    mask = Image.new("L", (WIDTH, HEIGHT), 0)
    mask_draw = ImageDraw.Draw(mask)
    reveal_width = int((text_box[2] - text_box[0] + 150) * min(1.0, raw * 1.08))
    mask_draw.rounded_rectangle((text_box[0] - 80, text_box[1] - 42, text_box[0] - 80 + reveal_width, text_box[3] + 58), radius=24, fill=int(255 * min(1, raw * 1.22)))
    mask = mask.filter(ImageFilter.GaussianBlur(max(0.2, 7 * (1 - raw))))

    text = prepared_text.copy()
    if raw < 0.85:
        text = text.filter(ImageFilter.GaussianBlur(5 * (1 - raw)))
    moved = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    moved.alpha_composite(apply_alpha(text, clamp(raw * 1.38)), (x_offset, y_offset))
    moved.putalpha(Image.composite(moved.getchannel("A"), Image.new("L", (WIDTH, HEIGHT), 0), mask))
    base.alpha_composite(moved)

    sweep = smoothstep(1.17, 1.46, time) * (1 - smoothstep(1.50, 1.78, time))
    if sweep > 0:
        sweep_layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        draw = ImageDraw.Draw(sweep_layer)
        sx = int(text_box[0] - 210 + ((time - 1.15) / 0.58) * ((text_box[2] - text_box[0]) + 420))
        draw.polygon(
            [(sx, text_box[1] - 50), (sx + 118, text_box[1] - 50), (sx + 34, text_box[3] + 78), (sx - 84, text_box[3] + 78)],
            fill=(255, 255, 255, int(122 * sweep)),
        )
        sweep_layer.putalpha(Image.composite(sweep_layer.getchannel("A"), Image.new("L", (WIDTH, HEIGHT), 0), mask))
        base.alpha_composite(sweep_layer)

    sparkle = smoothstep(1.34, 1.62, time) * (1 - smoothstep(1.78, 2.0, time) * 0.45)
    if sparkle > 0:
        dots = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
        draw = ImageDraw.Draw(dots)
        rng = random.Random(74421)
        for _ in range(26):
            sx = rng.uniform(text_box[0] - 40, text_box[2] + 40)
            sy = rng.uniform(text_box[1] - 42, text_box[3] + 58)
            size = rng.uniform(1.4, 4.8)
            phase = rng.uniform(0.0, 0.22)
            alpha = int(150 * sparkle * (1 - smoothstep(1.62 + phase, 2.02 + phase, time)))
            if alpha <= 0:
                continue
            draw.ellipse((sx - size, sy - size, sx + size, sy + size), fill=(164, 255, 232, alpha))
        base.alpha_composite(dots.filter(ImageFilter.GaussianBlur(0.35)))


def add_finish_sweep(base: Image.Image, time: float) -> None:
    sweep = smoothstep(1.24, 1.54, time) * (1 - smoothstep(1.64, 1.92, time))
    if sweep <= 0:
        return
    layer = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(layer)
    sx = int(-240 + sweep * (WIDTH + 420))
    draw.polygon([(sx, 150), (sx + 180, 150), (sx - 120, HEIGHT - 70), (sx - 300, HEIGHT - 70)], fill=(255, 255, 255, int(38 * sweep)))
    base.alpha_composite(layer.filter(ImageFilter.GaussianBlur(15)))


def add_vignette(base: Image.Image) -> None:
    vignette = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    draw = ImageDraw.Draw(vignette)
    for index in range(118):
        alpha = max(0, int(152 - index * 1.28))
        draw.rectangle((index, index, WIDTH - index, HEIGHT - index), outline=(0, 0, 0, alpha), width=2)
    base.alpha_composite(vignette.filter(ImageFilter.GaussianBlur(18)))


def render_frame(index: int, logo: Image.Image, prepared_text: Image.Image, text_box: tuple[int, int, int, int]) -> Image.Image:
    time = index / FPS
    frame = render_background(time)
    draw_logo_shards(frame, time, logo)
    draw_url_comet(frame, prepared_text, text_box, time)
    add_finish_sweep(frame, time)
    add_vignette(frame)
    return frame.convert("RGB")


def encode_video(frames: list[Image.Image], output_path: Path) -> None:
    output_path.parent.mkdir(parents=True, exist_ok=True)
    container = av.open(str(output_path), "w")
    stream = container.add_stream("libx264", rate=FPS)
    stream.width = WIDTH
    stream.height = HEIGHT
    stream.pix_fmt = "yuv420p"
    stream.options = {"crf": "16", "preset": "medium", "profile": "high"}

    for image in frames:
        frame = av.VideoFrame.from_image(image)
        for packet in stream.encode(frame):
            container.mux(packet)
    for packet in stream.encode():
        container.mux(packet)
    container.close()


def write_comfy_assets(reference: Image.Image) -> None:
    REFERENCE_PATH.parent.mkdir(parents=True, exist_ok=True)
    COMFY_INPUT.mkdir(parents=True, exist_ok=True)
    reference.save(REFERENCE_PATH, quality=96)
    shutil.copyfile(REFERENCE_PATH, COMFY_INPUT / REFERENCE_NAME)

    prompt = json.loads(BASE_API_PROMPT.read_text(encoding="utf-8"))
    prompt["1"]["inputs"]["image"] = REFERENCE_NAME
    prompt["21"]["inputs"]["text"] = WAN_PROMPT
    prompt["22"]["inputs"]["text"] = NEGATIVE_PROMPT
    prompt["26"]["inputs"]["seed"] = 624042911
    prompt["26"]["inputs"]["steps"] = 14
    prompt["26"]["inputs"]["cfg"] = 5.6
    prompt["32"]["inputs"]["filename_prefix"] = "outro/poliglot_outro_2s_tech_warp_v11_api"
    prompt["33"]["inputs"]["filename_prefix"] = "outro/poliglot_outro_2s_tech_warp_v11_reference"
    API_PROMPT.write_text(json.dumps(prompt, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    PROMPT_DOC.write_text(
        "\n".join(
            [
                "# Poliglot AI 2s Tech Warp Outro v11",
                "",
                "A more kinetic replacement for the previous static/flat space outro.",
                "",
                "## On-Screen Text",
                "",
                "`poliglotAI.online`",
                "",
                "## Wan 2.1 I2V Prompt",
                "",
                WAN_PROMPT,
                "",
                "## Negative Prompt",
                "",
                NEGATIVE_PROMPT,
                "",
                "## Files",
                "",
                f"- Deterministic final MP4: `{OUTPUT_VIDEO}`.",
                f"- ComfyUI API prompt: `{API_PROMPT.name}`.",
                f"- ComfyUI input reference: `{REFERENCE_NAME}`.",
                "- Final URL is composited deterministically so social compression cannot distort the domain.",
            ]
        )
        + "\n",
        encoding="utf-8",
    )


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
    if not LOGO_PATH.exists():
        raise FileNotFoundError(f"Missing logo asset: {LOGO_PATH}")

    logo = Image.open(LOGO_PATH).convert("RGBA")
    font = load_font(92, bold=True)
    prepared_text, text_box = make_text_layer(font)

    frames = [render_frame(index, logo, prepared_text, text_box) for index in range(TOTAL_FRAMES)]
    encode_video(frames, OUTPUT_VIDEO)

    PREVIEW_DIR.mkdir(parents=True, exist_ok=True)
    first = PREVIEW_DIR / "outro_tech_warp_v11_first.png"
    mid = PREVIEW_DIR / "outro_tech_warp_v11_mid.png"
    last = PREVIEW_DIR / "outro_tech_warp_v11_last.png"
    sheet = PREVIEW_DIR / "outro_tech_warp_v11_contact_sheet.jpg"
    frames[0].save(first)
    frames[TOTAL_FRAMES // 2].save(mid)
    frames[-1].save(last)
    make_contact_sheet(frames, sheet)

    reference = frames[18].copy()
    write_comfy_assets(reference)

    result = {
        "video": str(OUTPUT_VIDEO),
        "api_prompt": str(API_PROMPT),
        "prompt_doc": str(PROMPT_DOC),
        "reference": str(REFERENCE_PATH),
        "comfy_input_reference": str(COMFY_INPUT / REFERENCE_NAME),
        "first_frame": str(first),
        "mid_frame": str(mid),
        "last_frame": str(last),
        "contact_sheet": str(sheet),
        "width": WIDTH,
        "height": HEIGHT,
        "fps": FPS,
        "frames": TOTAL_FRAMES,
        "duration": DURATION,
    }
    result_path = PREVIEW_DIR / "poliglot_outro_tech_warp_v11_result.json"
    result_path.write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    return result


if __name__ == "__main__":
    print(json.dumps(render(), indent=2))
