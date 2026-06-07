from __future__ import annotations

import json
import math
import random
import shutil
from dataclasses import dataclass
from pathlib import Path

from PIL import Image, ImageDraw, ImageFilter


OUT_DIR = Path("web/assets")
CONTACT_SHEET = Path("tmp/brand-assets-theme-contact.png")
THEMED_GROUPS = {"plans", "headers", "icons"}
THEMES = ("light", "dark")
PRIMARY_COMFY_MODEL = "sd3.5_large_fp8_scaled.safetensors"
COMPARISON_MODELS = [
    "Juggernaut-XL_v9_RunDiffusionPhoto_v2.safetensors",
    "sd_xl_base_1.0.safetensors",
    "FLUX.1-schnell / FLUX.1-dev candidates require license review before production use",
]
STYLE_PROMPT = "Poliglot AI premium language-learning asset, no text, no pseudo letters, no watermark"

INK = "#101827"
DEEP = "#10243A"
DEEP_2 = "#163D59"
BLUE = "#2D5BFF"
BLUE_2 = "#2457D6"
TEAL = "#00A88F"
TEAL_2 = "#28D7BC"
MINT = "#DDF8EF"
GOLD = "#F3B84B"
GOLD_2 = "#D89B28"
PLUM = "#7C5CFF"
ROSE = "#E45A6F"
WHITE = "#FFFFFF"
CLOUD = "#F6F8FC"


@dataclass(frozen=True)
class Asset:
    file: str
    group: str
    width: int
    height: int
    kind: str
    seed: int


def hex_rgb(value: str) -> tuple[int, int, int]:
    value = value.lstrip("#")
    return tuple(int(value[i : i + 2], 16) for i in (0, 2, 4))


def rgba(value: str, alpha: int = 255) -> tuple[int, int, int, int]:
    return (*hex_rgb(value), alpha)


def mix(a: str, b: str, t: float) -> tuple[int, int, int, int]:
    ar, ag, ab = hex_rgb(a)
    br, bg, bb = hex_rgb(b)
    return (
        round(ar + (br - ar) * t),
        round(ag + (bg - ag) * t),
        round(ab + (bb - ab) * t),
        255,
    )


def gradient(size: tuple[int, int], top: str, bottom: str) -> Image.Image:
    w, h = size
    strip = Image.new("RGBA", (1, h))
    pix = strip.load()
    for y in range(h):
        pix[0, y] = mix(top, bottom, y / max(h - 1, 1))
    return strip.resize((w, h), Image.Resampling.BICUBIC)


def add_glow(img: Image.Image, xy: tuple[float, float, float, float], color: str, alpha: int, blur: int) -> None:
    layer = Image.new("RGBA", img.size, (0, 0, 0, 0))
    d = ImageDraw.Draw(layer)
    d.ellipse(xy, fill=rgba(color, alpha))
    layer = layer.filter(ImageFilter.GaussianBlur(blur))
    img.alpha_composite(layer)


def rounded_shadow(
    img: Image.Image,
    xy: tuple[float, float, float, float],
    radius: float,
    fill: str,
    alpha: int = 255,
    shadow_alpha: int = 54,
    blur: int = 28,
    offset: tuple[int, int] = (0, 18),
    outline: str | None = None,
    outline_alpha: int = 96,
) -> None:
    x1, y1, x2, y2 = xy
    shadow = Image.new("RGBA", img.size, (0, 0, 0, 0))
    sd = ImageDraw.Draw(shadow)
    sd.rounded_rectangle((x1 + offset[0], y1 + offset[1], x2 + offset[0], y2 + offset[1]), radius, fill=(8, 18, 34, shadow_alpha))
    shadow = shadow.filter(ImageFilter.GaussianBlur(blur))
    img.alpha_composite(shadow)

    d = ImageDraw.Draw(img)
    d.rounded_rectangle(xy, radius, fill=rgba(fill, alpha), outline=rgba(outline or WHITE, outline_alpha), width=max(1, round(radius / 12)))
    inset = max(2, round(radius / 8))
    d.rounded_rectangle((x1 + inset, y1 + inset, x2 - inset, y1 + max(inset + 3, radius * 0.72)), radius * 0.72, fill=rgba(WHITE, 34))


def ellipse_shadow(
    img: Image.Image,
    xy: tuple[float, float, float, float],
    fill: str,
    alpha: int = 255,
    shadow_alpha: int = 58,
    blur: int = 24,
    offset: tuple[int, int] = (0, 14),
    outline: str | None = None,
) -> None:
    x1, y1, x2, y2 = xy
    shadow = Image.new("RGBA", img.size, (0, 0, 0, 0))
    sd = ImageDraw.Draw(shadow)
    sd.ellipse((x1 + offset[0], y1 + offset[1], x2 + offset[0], y2 + offset[1]), fill=(8, 18, 34, shadow_alpha))
    shadow = shadow.filter(ImageFilter.GaussianBlur(blur))
    img.alpha_composite(shadow)

    d = ImageDraw.Draw(img)
    d.ellipse(xy, fill=rgba(fill, alpha), outline=rgba(outline or WHITE, 95), width=max(2, round((x2 - x1) / 80)))
    d.arc((x1 + 7, y1 + 7, x2 - 7, y2 - 7), 205, 318, fill=rgba(WHITE, 108), width=max(2, round((x2 - x1) / 48)))


def draw_background(img: Image.Image, seed: int, dark: bool = False) -> None:
    rng = random.Random(seed)
    w, h = img.size
    if dark:
        img.alpha_composite(gradient(img.size, "#0E1826", "#153F56"))
        add_glow(img, (-w * 0.12, -h * 0.16, w * 0.52, h * 0.58), BLUE, 70, round(w * 0.09))
        add_glow(img, (w * 0.58, -h * 0.1, w * 1.15, h * 0.62), TEAL, 62, round(w * 0.08))
    else:
        img.alpha_composite(gradient(img.size, "#F8FCFF", "#E9F3F8"))
        add_glow(img, (-w * 0.14, -h * 0.18, w * 0.48, h * 0.6), BLUE, 45, round(w * 0.08))
        add_glow(img, (w * 0.58, -h * 0.08, w * 1.12, h * 0.68), TEAL, 44, round(w * 0.08))
        add_glow(img, (w * 0.62, h * 0.52, w * 1.08, h * 1.12), GOLD, 30, round(w * 0.08))

    route = Image.new("RGBA", img.size, (0, 0, 0, 0))
    d = ImageDraw.Draw(route)
    points = []
    for i in range(9):
        px = rng.uniform(w * 0.05, w * 0.95)
        py = rng.uniform(h * 0.12, h * 0.9)
        points.append((px, py))
    for a, b in zip(points, points[1:]):
        d.line((a, b), fill=rgba(TEAL_2 if dark else BLUE_2, 28), width=max(1, round(w / 420)))
    for px, py in points:
        r = max(3, round(w / 260))
        d.ellipse((px - r, py - r, px + r, py + r), fill=rgba(WHITE if dark else TEAL_2, 120))
    img.alpha_composite(route.filter(ImageFilter.GaussianBlur(0.4)))


def draw_wave(img: Image.Image, cx: float, cy: float, width: float, height: float, color: str, alpha: int = 180, waves: int = 2) -> None:
    d = ImageDraw.Draw(img)
    pts = []
    steps = 96
    for i in range(steps + 1):
        t = i / steps
        x = cx - width / 2 + width * t
        y = cy + math.sin(t * math.tau * waves) * height * (0.35 + 0.65 * math.sin(math.pi * t))
        pts.append((x, y))
    d.line(pts, fill=rgba(color, alpha), width=max(4, round(width / 70)), joint="curve")
    d.line([(x, y + height * 0.35) for x, y in pts], fill=rgba(WHITE, max(32, alpha // 4)), width=max(2, round(width / 140)))


def draw_gem(img: Image.Image, cx: float, cy: float, r: float, color: str, accent: str = GOLD) -> None:
    shadow = Image.new("RGBA", img.size, (0, 0, 0, 0))
    sd = ImageDraw.Draw(shadow)
    poly = [(cx, cy - r), (cx + r * 0.9, cy), (cx, cy + r), (cx - r * 0.9, cy)]
    sd.polygon([(x, y + r * 0.22) for x, y in poly], fill=(8, 18, 34, 56))
    shadow = shadow.filter(ImageFilter.GaussianBlur(round(r * 0.35)))
    img.alpha_composite(shadow)
    d = ImageDraw.Draw(img)
    d.polygon(poly, fill=rgba(color, 238), outline=rgba(accent, 145))
    d.line((cx - r * 0.55, cy - r * 0.1, cx, cy - r * 0.62, cx + r * 0.55, cy - r * 0.1), fill=rgba(WHITE, 130), width=max(2, round(r / 12)))


def draw_ring(img: Image.Image, cx: float, cy: float, r: float, color: str, width: float, start: int = 0, end: int = 360) -> None:
    d = ImageDraw.Draw(img)
    box = (cx - r, cy - r, cx + r, cy + r)
    d.arc(box, start, end, fill=rgba(color, 210), width=max(3, round(width)))
    d.arc((cx - r + width * 1.7, cy - r + width * 1.7, cx + r - width * 1.7, cy + r - width * 1.7), start + 20, end - 20, fill=rgba(WHITE, 84), width=max(2, round(width / 2)))


def draw_speech(img: Image.Image, x: float, y: float, w: float, h: float, color: str, flip: bool = False) -> None:
    rounded_shadow(img, (x, y, x + w, y + h), h * 0.32, color, alpha=235, shadow_alpha=42, blur=max(14, round(w * 0.08)))
    d = ImageDraw.Draw(img)
    if flip:
        tail = [(x + w * 0.72, y + h * 0.84), (x + w * 0.84, y + h * 1.12), (x + w * 0.52, y + h * 0.86)]
    else:
        tail = [(x + w * 0.28, y + h * 0.84), (x + w * 0.16, y + h * 1.12), (x + w * 0.48, y + h * 0.86)]
    d.polygon(tail, fill=rgba(color, 235))
    d.ellipse((x + w * 0.18, y + h * 0.34, x + w * 0.28, y + h * 0.52), fill=rgba(WHITE, 180))
    d.ellipse((x + w * 0.43, y + h * 0.34, x + w * 0.53, y + h * 0.52), fill=rgba(WHITE, 142))
    d.ellipse((x + w * 0.68, y + h * 0.34, x + w * 0.78, y + h * 0.52), fill=rgba(WHITE, 180))


def draw_tiles(img: Image.Image, cx: float, cy: float, scale: float, colors: list[str]) -> None:
    for i, color in enumerate(colors):
        row, col = divmod(i, 3)
        x = cx + (col - 1) * scale * 1.12
        y = cy + (row - 1) * scale * 1.05
        rounded_shadow(img, (x - scale * 0.42, y - scale * 0.36, x + scale * 0.42, y + scale * 0.36), scale * 0.16, color, alpha=232, shadow_alpha=32, blur=round(scale * 0.2))
        if i in (1, 4, 7):
            draw_gem(img, x, y, scale * 0.13, WHITE, color)


def draw_podium(img: Image.Image, cx: float, cy: float, scale: float) -> None:
    heights = [0.62, 1.0, 0.78]
    colors = [TEAL, GOLD, BLUE_2]
    for i, hmul in enumerate(heights):
        x = cx + (i - 1) * scale * 0.86
        h = scale * hmul
        rounded_shadow(img, (x - scale * 0.28, cy - h, x + scale * 0.28, cy), scale * 0.12, colors[i], alpha=232, shadow_alpha=38, blur=round(scale * 0.18))
        ellipse_shadow(img, (x - scale * 0.18, cy - h - scale * 0.42, x + scale * 0.18, cy - h - scale * 0.06), WHITE if i != 1 else GOLD, alpha=238, shadow_alpha=32, blur=round(scale * 0.12), outline=colors[i])


def draw_crown(img: Image.Image, cx: float, cy: float, scale: float) -> None:
    d = ImageDraw.Draw(img)
    shadow = Image.new("RGBA", img.size, (0, 0, 0, 0))
    sd = ImageDraw.Draw(shadow)
    pts = [
        (cx - scale * 0.55, cy + scale * 0.25),
        (cx - scale * 0.38, cy - scale * 0.32),
        (cx - scale * 0.08, cy + scale * 0.02),
        (cx, cy - scale * 0.48),
        (cx + scale * 0.08, cy + scale * 0.02),
        (cx + scale * 0.38, cy - scale * 0.32),
        (cx + scale * 0.55, cy + scale * 0.25),
    ]
    sd.polygon([(x, y + scale * 0.16) for x, y in pts], fill=(8, 18, 34, 56))
    shadow = shadow.filter(ImageFilter.GaussianBlur(round(scale * 0.16)))
    img.alpha_composite(shadow)
    d.polygon(pts, fill=rgba(GOLD, 238), outline=rgba(WHITE, 120))
    d.rounded_rectangle((cx - scale * 0.52, cy + scale * 0.18, cx + scale * 0.52, cy + scale * 0.42), scale * 0.09, fill=rgba(GOLD_2, 245), outline=rgba(WHITE, 112))
    d.line((cx - scale * 0.35, cy + scale * 0.08, cx + scale * 0.33, cy + scale * 0.08), fill=rgba(WHITE, 130), width=max(2, round(scale * 0.04)))


def draw_check(img: Image.Image, cx: float, cy: float, scale: float, color: str = TEAL_2) -> None:
    d = ImageDraw.Draw(img)
    pts = [(cx - scale * 0.46, cy), (cx - scale * 0.13, cy + scale * 0.3), (cx + scale * 0.5, cy - scale * 0.38)]
    d.line(pts, fill=rgba(color, 230), width=max(5, round(scale * 0.12)), joint="curve")
    d.line([(x, y - scale * 0.04) for x, y in pts], fill=rgba(WHITE, 88), width=max(2, round(scale * 0.04)), joint="curve")


def draw_object(img: Image.Image, kind: str, center: tuple[float, float], scale: float, seed: int) -> None:
    rng = random.Random(seed)
    cx, cy = center
    if kind in {"home", "hero"}:
        draw_ring(img, cx, cy, scale * 0.72, TEAL_2, scale * 0.06, 200, 520)
        draw_speech(img, cx - scale * 0.78, cy - scale * 0.32, scale * 0.8, scale * 0.48, BLUE_2)
        draw_speech(img, cx - scale * 0.04, cy - scale * 0.12, scale * 0.76, scale * 0.44, TEAL, True)
        draw_wave(img, cx, cy + scale * 0.56, scale * 1.34, scale * 0.22, GOLD, 190, 2)
        for a in range(5):
            angle = a * math.tau / 5 + 0.4
            draw_gem(img, cx + math.cos(angle) * scale * 0.95, cy + math.sin(angle) * scale * 0.55, scale * 0.08, MINT if a % 2 else WHITE, TEAL)
    elif kind in {"lesson", "free"}:
        rounded_shadow(img, (cx - scale * 0.72, cy - scale * 0.3, cx + scale * 0.1, cy + scale * 0.32), scale * 0.12, TEAL, 238)
        rounded_shadow(img, (cx - scale * 0.05, cy - scale * 0.38, cx + scale * 0.78, cy + scale * 0.24), scale * 0.12, WHITE, 245, outline=TEAL)
        d = ImageDraw.Draw(img)
        d.line((cx - scale * 0.02, cy - scale * 0.34, cx - scale * 0.02, cy + scale * 0.26), fill=rgba(GOLD, 190), width=max(4, round(scale * 0.05)))
        draw_gem(img, cx + scale * 0.52, cy - scale * 0.42, scale * 0.12, GOLD, WHITE)
        draw_check(img, cx - scale * 0.42, cy + scale * 0.48, scale * 0.32)
    elif kind in {"practice"}:
        draw_speech(img, cx - scale * 0.88, cy - scale * 0.34, scale * 0.9, scale * 0.55, BLUE_2)
        draw_speech(img, cx - scale * 0.08, cy - scale * 0.08, scale * 0.88, scale * 0.52, TEAL, True)
        draw_wave(img, cx, cy + scale * 0.62, scale * 1.35, scale * 0.18, GOLD, 200, 3)
    elif kind in {"words", "vocabulary", "word-game", "premium"}:
        colors = [BLUE_2, TEAL, WHITE, MINT, GOLD, BLUE_2, WHITE, TEAL_2, MINT]
        draw_tiles(img, cx, cy, scale * (0.5 if kind != "word-game" else 0.44), colors)
        if kind == "premium":
            draw_crown(img, cx, cy - scale * 0.72, scale * 0.34)
    elif kind in {"spelling", "mistakes"}:
        draw_ring(img, cx, cy, scale * 0.62, ROSE if kind == "mistakes" else BLUE_2, scale * 0.07, -40, 250)
        draw_check(img, cx, cy, scale * 0.82, TEAL_2)
        for i in range(4):
            draw_gem(img, cx + (i - 1.5) * scale * 0.34, cy + scale * 0.55 + rng.uniform(-8, 8), scale * 0.08, WHITE if i % 2 else MINT, TEAL)
    elif kind in {"progress", "limits"}:
        draw_ring(img, cx, cy, scale * 0.62, BLUE_2, scale * 0.09, 205, 540)
        draw_ring(img, cx, cy, scale * 0.42, TEAL_2, scale * 0.07, 110, 390)
        draw_ring(img, cx, cy, scale * 0.24, GOLD, scale * 0.05, -20, 235)
        if kind == "limits":
            for i, color in enumerate([TEAL, BLUE_2, GOLD]):
                x = cx - scale * 0.5 + i * scale * 0.5
                rounded_shadow(img, (x - scale * 0.08, cy + scale * 0.68 - scale * 0.18 * i, x + scale * 0.08, cy + scale * 0.88), scale * 0.06, color, 230, blur=round(scale * 0.09))
    elif kind in {"awards", "platinum"}:
        draw_ring(img, cx, cy, scale * 0.68, GOLD, scale * 0.07, 210, 520)
        draw_crown(img, cx, cy - scale * 0.08, scale * 0.62)
        for i in range(5):
            draw_gem(img, cx + (i - 2) * scale * 0.33, cy + scale * 0.62 + rng.uniform(-8, 8), scale * 0.08, TEAL_2 if i % 2 else WHITE, GOLD)
    elif kind == "leaderboard":
        draw_podium(img, cx, cy + scale * 0.55, scale * 0.9)
        draw_ring(img, cx, cy - scale * 0.28, scale * 0.35, TEAL_2, scale * 0.05, 0, 320)
    elif kind == "level":
        for i in range(5):
            x = cx - scale * 0.62 + i * scale * 0.31
            y = cy + scale * 0.46 - i * scale * 0.22
            ellipse_shadow(img, (x - scale * 0.1, y - scale * 0.1, x + scale * 0.1, y + scale * 0.1), [TEAL, BLUE_2, GOLD, TEAL_2, PLUM][i], 240)
            if i:
                ImageDraw.Draw(img).line((cx - scale * 0.62 + (i - 1) * scale * 0.31, cy + scale * 0.46 - (i - 1) * scale * 0.22, x, y), fill=rgba(WHITE, 155), width=max(3, round(scale * 0.04)))
    elif kind == "tools":
        draw_ring(img, cx - scale * 0.28, cy - scale * 0.02, scale * 0.42, TEAL_2, scale * 0.06, 0, 360)
        rounded_shadow(img, (cx + scale * 0.03, cy - scale * 0.4, cx + scale * 0.58, cy + scale * 0.15), scale * 0.16, WHITE, 238, outline=TEAL)
        draw_wave(img, cx - scale * 0.1, cy + scale * 0.5, scale * 1.25, scale * 0.2, GOLD, 205, 3)
        draw_gem(img, cx + scale * 0.53, cy - scale * 0.42, scale * 0.12, BLUE_2, WHITE)
    elif kind == "referral":
        ellipse_shadow(img, (cx - scale * 0.42, cy - scale * 0.42, cx + scale * 0.42, cy + scale * 0.42), WHITE, 242, outline=TEAL)
        d = ImageDraw.Draw(img)
        d.rounded_rectangle((cx - scale * 0.48, cy - scale * 0.08, cx + scale * 0.48, cy + scale * 0.18), scale * 0.06, fill=rgba(GOLD, 235))
        d.rounded_rectangle((cx - scale * 0.08, cy - scale * 0.48, cx + scale * 0.12, cy + scale * 0.46), scale * 0.04, fill=rgba(BLUE_2, 235))
        draw_ring(img, cx - scale * 0.18, cy - scale * 0.46, scale * 0.18, GOLD, scale * 0.04, 210, 520)
        draw_ring(img, cx + scale * 0.18, cy - scale * 0.46, scale * 0.18, GOLD, scale * 0.04, 20, 330)
    elif kind == "settings":
        for i, color in enumerate([BLUE_2, TEAL, GOLD]):
            y = cy - scale * 0.36 + i * scale * 0.34
            ImageDraw.Draw(img).line((cx - scale * 0.6, y, cx + scale * 0.6, y), fill=rgba(WHITE, 138), width=max(4, round(scale * 0.055)))
            ellipse_shadow(img, (cx - scale * 0.36 + i * scale * 0.34, y - scale * 0.12, cx - scale * 0.12 + i * scale * 0.34, y + scale * 0.12), color, 238)
    elif kind == "languages":
        draw_ring(img, cx, cy, scale * 0.7, TEAL_2, scale * 0.05, 195, 535)
        for i in range(12):
            a = i * math.tau / 12
            x = cx + math.cos(a) * scale * 0.72
            y = cy + math.sin(a) * scale * 0.43
            draw_gem(img, x, y, scale * 0.055, [WHITE, MINT, TEAL_2, GOLD][i % 4], TEAL)
        draw_speech(img, cx - scale * 0.5, cy - scale * 0.22, scale * 1.0, scale * 0.5, BLUE_2)
    elif kind == "legal":
        draw_ring(img, cx, cy, scale * 0.65, TEAL_2, scale * 0.05, 200, 520)
        shield = [(cx, cy - scale * 0.62), (cx + scale * 0.46, cy - scale * 0.4), (cx + scale * 0.36, cy + scale * 0.3), (cx, cy + scale * 0.66), (cx - scale * 0.36, cy + scale * 0.3), (cx - scale * 0.46, cy - scale * 0.4)]
        shadow = Image.new("RGBA", img.size, (0, 0, 0, 0))
        sd = ImageDraw.Draw(shadow)
        sd.polygon([(x, y + scale * 0.12) for x, y in shield], fill=(8, 18, 34, 58))
        img.alpha_composite(shadow.filter(ImageFilter.GaussianBlur(round(scale * 0.15))))
        d = ImageDraw.Draw(img)
        d.polygon(shield, fill=rgba(WHITE, 238), outline=rgba(TEAL, 145))
        draw_check(img, cx, cy - scale * 0.02, scale * 0.58, TEAL)
    else:
        draw_ring(img, cx, cy, scale * 0.62, TEAL_2, scale * 0.06, 200, 520)
        draw_gem(img, cx, cy, scale * 0.26, WHITE, GOLD)


def themed_file(file_name: str, theme: str) -> str:
    path = Path(file_name)
    return f"{path.stem}-{theme}{path.suffix}"


def render_asset(asset: Asset, theme: str | None = None) -> Path:
    img = Image.new("RGBA", (asset.width, asset.height), (0, 0, 0, 0))
    dark = theme == "dark" if theme else asset.group in {"site"} or asset.kind in {"home", "practice", "tools", "premium", "platinum", "languages"}
    draw_background(img, asset.seed, dark=dark)

    w, h = asset.width, asset.height
    if asset.group == "icons":
        rounded_shadow(img, (w * 0.16, h * 0.16, w * 0.84, h * 0.84), w * 0.14, "#122B42" if dark else WHITE, 210, shadow_alpha=46, blur=28, outline=TEAL_2)
        draw_object(img, asset.kind, (w * 0.5, h * 0.5), min(w, h) * 0.42, asset.seed)
    elif asset.group == "headers":
        draw_object(img, asset.kind, (w * 0.68, h * 0.52), min(w, h) * 0.58, asset.seed)
        add_glow(img, (w * 0.0, h * 0.1, w * 0.56, h * 0.9), WHITE if not dark else BLUE, 18, round(w * 0.05))
    elif asset.group == "plans":
        draw_object(img, asset.kind, (w * 0.56, h * 0.52), min(w, h) * 0.62, asset.seed)
        draw_wave(img, w * 0.55, h * 0.78, w * 0.55, h * 0.07, TEAL_2, 92, 2)
    elif asset.file == "site-hero-cockpit.png":
        draw_object(img, "hero", (w * 0.68, h * 0.52), h * 0.44, asset.seed)
        add_glow(img, (w * -0.08, h * 0.04, w * 0.58, h * 0.96), DEEP, 70, round(w * 0.08))
    else:
        draw_object(img, asset.kind, (w * 0.58, h * 0.52), h * 0.42, asset.seed)
        add_glow(img, (w * 0.02, h * 0.12, w * 0.52, h * 0.88), WHITE if not dark else BLUE, 18, round(w * 0.06))

    OUT_DIR.mkdir(parents=True, exist_ok=True)
    output_file = themed_file(asset.file, theme) if theme else asset.file
    output_path = OUT_DIR / output_file
    img.convert("RGB").save(output_path, optimize=True)
    return output_path


ASSETS = [
    Asset("plan-free.png", "plans", 1280, 768, "free", 2601),
    Asset("plan-premium.png", "plans", 1280, 768, "premium", 2602),
    Asset("plan-platinum.png", "plans", 1280, 768, "platinum", 2603),
    Asset("header-home.png", "headers", 1536, 640, "home", 2610),
    Asset("header-lesson.png", "headers", 1536, 640, "lesson", 2611),
    Asset("header-practice.png", "headers", 1536, 640, "practice", 2612),
    Asset("header-progress.png", "headers", 1536, 640, "progress", 2613),
    Asset("header-awards.png", "headers", 1536, 640, "awards", 2614),
    Asset("header-vocabulary.png", "headers", 1536, 640, "vocabulary", 2615),
    Asset("header-words.png", "headers", 1536, 640, "words", 2621),
    Asset("header-word-game.png", "headers", 1536, 640, "word-game", 2622),
    Asset("header-spelling.png", "headers", 1536, 640, "spelling", 2623),
    Asset("header-level.png", "headers", 1536, 640, "level", 2616),
    Asset("header-leaderboard.png", "headers", 1536, 640, "leaderboard", 2624),
    Asset("header-premium.png", "headers", 1536, 640, "premium", 2617),
    Asset("header-limits.png", "headers", 1536, 640, "limits", 2625),
    Asset("header-mistakes.png", "headers", 1536, 640, "mistakes", 2618),
    Asset("header-tools.png", "headers", 1536, 640, "tools", 2619),
    Asset("header-referral.png", "headers", 1536, 640, "referral", 2620),
    Asset("header-settings.png", "headers", 1536, 640, "settings", 2626),
    Asset("icon-home.png", "icons", 512, 512, "home", 2701),
    Asset("icon-lesson.png", "icons", 512, 512, "lesson", 2702),
    Asset("icon-practice.png", "icons", 512, 512, "practice", 2703),
    Asset("icon-words.png", "icons", 512, 512, "words", 2704),
    Asset("icon-word-game.png", "icons", 512, 512, "word-game", 2705),
    Asset("icon-spelling.png", "icons", 512, 512, "spelling", 2706),
    Asset("icon-vocabulary.png", "icons", 512, 512, "vocabulary", 2707),
    Asset("icon-level.png", "icons", 512, 512, "level", 2708),
    Asset("icon-progress.png", "icons", 512, 512, "progress", 2709),
    Asset("icon-awards.png", "icons", 512, 512, "awards", 2710),
    Asset("icon-leaderboard.png", "icons", 512, 512, "leaderboard", 2711),
    Asset("icon-mistakes.png", "icons", 512, 512, "mistakes", 2712),
    Asset("icon-referral.png", "icons", 512, 512, "referral", 2713),
    Asset("icon-premium.png", "icons", 512, 512, "premium", 2714),
    Asset("icon-limits.png", "icons", 512, 512, "limits", 2715),
    Asset("icon-tools.png", "icons", 512, 512, "tools", 2716),
    Asset("icon-settings.png", "icons", 512, 512, "settings", 2717),
    Asset("site-hero-cockpit.png", "site", 1920, 1080, "hero", 2801),
    Asset("site-languages-map.png", "site", 1600, 900, "languages", 2802),
    Asset("site-feature-tools.png", "site", 1600, 900, "tools", 2803),
    Asset("site-progress-awards.png", "site", 1600, 900, "awards", 2804),
    Asset("site-pricing-premium.png", "site", 1600, 900, "premium", 2805),
    Asset("site-legal-shield.png", "site", 1600, 900, "legal", 2806),
]


def manifest_entry(asset: Asset, file_name: str, theme: str | None = None, alias_for: str | None = None) -> dict[str, object]:
    return {
        "file": file_name,
        "group": asset.group,
        "theme": theme or "fixed",
        "width": asset.width,
        "height": asset.height,
        "seed": asset.seed,
        "prompt": STYLE_PROMPT,
        "generator": "tools/create_poliglot_premium_assets.py",
        "model_primary": PRIMARY_COMFY_MODEL,
        "comparison_models": COMPARISON_MODELS,
        "production_source": "procedural-safe-theme-pack",
        "compatibility_alias_for": alias_for,
        "qa_status": "generated_no_text_visual_review_required",
    }


def write_contact_sheet(files: list[Path]) -> None:
    if not files:
        return
    CONTACT_SHEET.parent.mkdir(parents=True, exist_ok=True)
    cell_w, cell_h = 260, 170
    cols = 4
    rows = math.ceil(len(files) / cols)
    sheet = Image.new("RGB", (cols * cell_w, rows * cell_h), "#f3f6fb")
    for index, file_path in enumerate(files):
        if not file_path.exists():
            continue
        img = Image.open(file_path).convert("RGB")
        img.thumbnail((cell_w - 24, cell_h - 24), Image.Resampling.LANCZOS)
        x = index % cols * cell_w + (cell_w - img.width) // 2
        y = index // cols * cell_h + (cell_h - img.height) // 2
        sheet.paste(img, (x, y))
    sheet.save(CONTACT_SHEET, optimize=True)


def main() -> None:
    manifest = []
    themed_outputs: list[Path] = []
    for asset in ASSETS:
        if asset.group in THEMED_GROUPS:
            light_alias = None
            for theme in THEMES:
                output_path = render_asset(asset, theme)
                themed_outputs.append(output_path)
                manifest.append(manifest_entry(asset, output_path.name, theme))
                if theme == "light":
                    light_alias = output_path
                print(f"saved {output_path}")
            if light_alias:
                legacy_path = OUT_DIR / asset.file
                shutil.copyfile(light_alias, legacy_path)
                manifest.append(manifest_entry(asset, asset.file, "light", alias_for=light_alias.name))
                print(f"saved {legacy_path}")
            continue

        output_path = render_asset(asset)
        manifest.append(manifest_entry(asset, asset.file))
        print(f"saved {output_path}")

    write_contact_sheet(themed_outputs)
    (OUT_DIR / "brand-assets-manifest.json").write_text(json.dumps(manifest, indent=2), encoding="utf-8")
    print(f"saved {CONTACT_SHEET}")


if __name__ == "__main__":
    main()
