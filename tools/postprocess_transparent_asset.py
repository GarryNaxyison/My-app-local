#!/usr/bin/env python3
"""Convert a chroma-key UI asset into a normalized transparent PNG."""

from __future__ import annotations

import argparse
import json
from pathlib import Path
from statistics import median

from PIL import Image, ImageFilter


KINDS = {
    "icon": {
        "size": 512,
        "fill": 0.92,
        "max_density": 0.78,
        "min_short_ratio": 0.72,
        "max_edge_coverage": 0.18,
        "edge_density_threshold": 0.55,
    },
    "award": {
        "size": 768,
        "fill": 0.91,
        "max_density": 0.78,
        "min_short_ratio": 0.58,
        "max_edge_coverage": 0.20,
        "edge_density_threshold": 0.60,
    },
}


def border_key(image: Image.Image) -> tuple[int, int, int]:
    width, height = image.size
    px = image.load()
    band = max(4, min(width, height) // 96)
    step = max(1, min(width, height) // 256)
    samples: list[tuple[int, int, int]] = []
    for x in range(0, width, step):
        for y in range(band):
            samples.append(px[x, y][:3])
            samples.append(px[x, height - 1 - y][:3])
    for y in range(0, height, step):
        for x in range(band):
            samples.append(px[x, y][:3])
            samples.append(px[width - 1 - x, y][:3])
    return (
        int(round(median(channel[0] for channel in samples))),
        int(round(median(channel[1] for channel in samples))),
        int(round(median(channel[2] for channel in samples))),
    )


def smoothstep(value: float) -> float:
    value = max(0.0, min(1.0, value))
    return value * value * (3.0 - 2.0 * value)


def apply_chroma_alpha(image: Image.Image, key: tuple[int, int, int]) -> Image.Image:
    rgba = image.convert("RGBA")
    px = rgba.load()
    width, height = rgba.size
    transparent_threshold = 20
    opaque_threshold = 118

    for y in range(height):
        for x in range(width):
            red, green, blue, original_alpha = px[x, y]
            distance = max(abs(red - key[0]), abs(green - key[1]), abs(blue - key[2]))
            green_dominance = green - max(red, blue)
            key_like = distance <= 64 or green_dominance >= 20
            if key_like:
                if distance <= transparent_threshold:
                    alpha = 0
                elif distance >= opaque_threshold:
                    alpha = 255
                else:
                    alpha = int(round(255 * smoothstep((distance - transparent_threshold) / (opaque_threshold - transparent_threshold))))
                if green_dominance > 0 and alpha < 252:
                    cap = max(red, blue)
                    green = min(green, cap)
            else:
                alpha = 255
            alpha = int(round(alpha * (original_alpha / 255)))
            if alpha <= 6:
                px[x, y] = (0, 0, 0, 0)
            else:
                px[x, y] = (red, green, blue, alpha)

    alpha = rgba.getchannel("A").filter(ImageFilter.MinFilter(3)).filter(ImageFilter.GaussianBlur(0.25))
    rgba.putalpha(alpha)
    return rgba


def normalize_canvas(image: Image.Image, size: int, fill: float) -> tuple[Image.Image, tuple[int, int, int, int] | None]:
    alpha = image.getchannel("A")
    bbox = alpha.point(lambda value: 255 if value > 18 else 0).getbbox()
    canvas = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    if not bbox:
        return canvas, None

    crop = image.crop(bbox)
    target = int(round(size * fill))
    scale = target / max(crop.size)
    resized = crop.resize(
        (max(1, int(round(crop.width * scale))), max(1, int(round(crop.height * scale)))),
        Image.Resampling.LANCZOS,
    )
    canvas.alpha_composite(resized, ((size - resized.width) // 2, (size - resized.height) // 2))
    final_bbox = canvas.getchannel("A").point(lambda value: 255 if value > 18 else 0).getbbox()
    return canvas, final_bbox


def edge_coverage(alpha: Image.Image, bbox: tuple[int, int, int, int] | None) -> float:
    if not bbox:
        return 1.0
    left, top, right, bottom = bbox
    px = alpha.load()
    samples = []
    for x in range(left, right):
        samples.append(px[x, top] > 18)
        samples.append(px[x, bottom - 1] > 18)
    for y in range(top, bottom):
        samples.append(px[left, y] > 18)
        samples.append(px[right - 1, y] > 18)
    return sum(samples) / max(1, len(samples))


def validate(image: Image.Image, bbox: tuple[int, int, int, int] | None, settings: dict[str, float]) -> dict[str, object]:
    alpha = image.getchannel("A")
    size = image.width
    target = int(round(size * settings["fill"]))
    messages: list[str] = []
    ok = True

    if image.mode != "RGBA":
        ok = False
        messages.append("output is not RGBA")
    if image.size != (size, size):
        ok = False
        messages.append("output is not square")
    if not bbox:
        ok = False
        messages.append("no visible object after chroma removal")
        return {"ok": ok, "messages": messages}

    left, top, right, bottom = bbox
    bbox_width = right - left
    bbox_height = bottom - top
    max_dim = max(bbox_width, bbox_height)
    min_dim = min(bbox_width, bbox_height)
    histogram = alpha.histogram()
    visible = sum(histogram[19:])
    density = visible / max(1, bbox_width * bbox_height)
    edge = edge_coverage(alpha, bbox)
    corners = [
        alpha.getpixel((0, 0)),
        alpha.getpixel((size - 1, 0)),
        alpha.getpixel((0, size - 1)),
        alpha.getpixel((size - 1, size - 1)),
    ]

    if abs(max_dim - target) > 3:
        ok = False
        messages.append(f"object max dimension {max_dim}px does not match target {target}px")
    min_target = int(round(target * settings["min_short_ratio"]))
    if min_dim < min_target:
        ok = False
        messages.append(f"object short dimension {min_dim}px is too small; expected at least {min_target}px")
    if max(corners) > 4:
        ok = False
        messages.append("transparent canvas corners are not fully transparent")
    if density > settings["max_density"]:
        ok = False
        messages.append("object still looks like a filled square/tile instead of a transparent cutout")
    if edge > settings["max_edge_coverage"] and density > settings["edge_density_threshold"]:
        ok = False
        messages.append("object edge coverage suggests a tile/backplate instead of a freestanding cutout")
    if visible / (size * size) < 0.12:
        ok = False
        messages.append("object occupies too little of the canvas")

    return {
        "ok": ok,
        "messages": messages,
        "bbox": [left, top, right, bottom],
        "bbox_width": bbox_width,
        "bbox_height": bbox_height,
        "max_dim": max_dim,
        "min_dim": min_dim,
        "target_dim": target,
        "visible_density_in_bbox": round(density, 4),
        "bbox_edge_coverage": round(edge, 4),
        "visible_canvas_ratio": round(visible / (size * size), 4),
        "corner_alpha_max": max(corners),
    }


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", required=True)
    parser.add_argument("--out", required=True)
    parser.add_argument("--kind", choices=sorted(KINDS), required=True)
    parser.add_argument("--report", required=True)
    args = parser.parse_args()

    settings = KINDS[args.kind]
    source = Image.open(args.input).convert("RGBA")
    key = border_key(source)
    keyed = apply_chroma_alpha(source, key)
    normalized, bbox = normalize_canvas(keyed, int(settings["size"]), float(settings["fill"]))
    report = validate(normalized, bbox, settings)
    report["input"] = str(Path(args.input))
    report["out"] = str(Path(args.out))
    report["kind"] = args.kind
    report["sampled_key"] = "#{:02x}{:02x}{:02x}".format(*key)

    report_path = Path(args.report)
    report_path.parent.mkdir(parents=True, exist_ok=True)
    report_path.write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")

    print(json.dumps(report, ensure_ascii=False))
    if not report["ok"]:
        raise SystemExit(2)

    out = Path(args.out)
    out.parent.mkdir(parents=True, exist_ok=True)
    temp_out = out.with_name(f"{out.stem}.tmp{out.suffix}")
    normalized.save(temp_out)
    temp_out.replace(out)


if __name__ == "__main__":
    main()
