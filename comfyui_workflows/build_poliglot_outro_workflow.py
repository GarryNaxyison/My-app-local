from __future__ import annotations

import copy
import json
import math
import shutil
import uuid
import wave
from pathlib import Path

from PIL import Image, ImageDraw, ImageFilter, ImageFont

from build_ad_shorts_workflow import (
    COMFY_INPUT_DIR,
    COMFY_WORKFLOW_DIR,
    OUTPUT_DIR,
    build_workflow,
)


PROJECT_ROOT = Path(__file__).resolve().parents[1]
OUTRO_SLUG = "poliglot_outro_2s"
OUTRO_WORKFLOW_NAME = f"{OUTRO_SLUG}_native.json"
OUTRO_API_NAME = f"{OUTRO_SLUG}_wan_stage_api.json"
OUTRO_REFERENCE_NAME = "poliglot_outro_reference_1080x1920.png"
OUTRO_AUDIO_NAME = "poliglot_outro_2s_silent.wav"

WIDTH = 1080
HEIGHT = 1920
WAN_WIDTH = 480
WAN_HEIGHT = 832
WAN_LENGTH = 33
FPS = 32
RIFE_MULTIPLIER = 2
OUTRO_DURATION_SECONDS = WAN_LENGTH * RIFE_MULTIPLIER / FPS


FLUX_OUTRO_PROMPT = """Vertical 9:16 premium social outro keyframe for Poliglot AI.
Use the supplied Poliglot AI brand card as the visual identity lock: central 3D sapphire logo, mint neon language-wave orbit, silver/glass highlights, dark luxury tech background.
Composition: clean TikTok/Reels/Shorts end-card, centered logo reveal, exact URL zone below the logo, high contrast, crisp premium product identity.
Keep the supplied URL layout stable. Do not invent new text, subtitles, people, app screenshots, badges, icons, flags, or extra brand marks.
Mood: fast modern TikTok end screen, glossy light sweep, polished SaaS ad finish, elite but readable."""


WAN_OUTRO_PROMPT = """2-second vertical TikTok/Reels/Shorts outro for Poliglot AI, starting from the supplied branded end-card.
Motion: quick premium logo pop-in feel, subtle push-in, sapphire gem glint, mint orbit ribbon gently waves, soft light sweep passes behind the URL, tiny particle sparkle around the logo, smooth energetic finish.
Keep the exact logo and the exact visible text "poliglotAI.online" stable and readable. No new captions, no extra words, no people, no UI screens, no watermark."""


OUTRO_NEGATIVE_PROMPT = """misspelled URL, wrong text, extra text, subtitles, captions, watermark, duplicated logo, distorted logo, deformed letters, fake letters, unreadable text, random words, people, hands, faces, flags, app screenshots, low quality, jpeg artifacts, blur, flicker, noisy text, overexposed, bad composition, black frames"""


OUTRO_NOTES = f"""## Poliglot AI 2s outro

Use this as the final insert for TikTok/Reels/Shorts.

Expected output:
- vertical 9:16 MP4;
- about {OUTRO_DURATION_SECONDS:.2f}s;
- exact visible text: poliglotAI.online;
- brand logo stays central and readable;
- no voiceover by default.

Default route:
- the exact reference card is sent directly into Wan 2.1 I2V to protect the URL;
- FLUX.2 nodes are kept as a disabled creative keyframe option.

If the URL drifts during generation, use the generated clip as the moving background and overlay the reference card/text in an editor for the last 0.5s."""


def cover_resize(image: Image.Image, width: int, height: int) -> Image.Image:
    scale = max(width / image.width, height / image.height)
    size = (math.ceil(image.width * scale), math.ceil(image.height * scale))
    resized = image.resize(size, Image.Resampling.LANCZOS)
    left = (resized.width - width) // 2
    top = (resized.height - height) // 2
    return resized.crop((left, top, left + width, top + height))


def load_font(size: int, bold: bool = False) -> ImageFont.FreeTypeFont:
    font_name = "segoeuib.ttf" if bold else "segoeui.ttf"
    font_path = Path("C:/Windows/Fonts") / font_name
    if font_path.exists():
        return ImageFont.truetype(str(font_path), size=size)
    return ImageFont.load_default(size=size)


def draw_centered_text(draw: ImageDraw.ImageDraw, text: str, y: int, font: ImageFont.FreeTypeFont) -> None:
    bbox = draw.textbbox((0, 0), text, font=font)
    x = (WIDTH - (bbox[2] - bbox[0])) // 2
    shadow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    shadow_draw = ImageDraw.Draw(shadow)
    shadow_draw.text((x, y), text, font=font, fill=(0, 0, 0, 210))
    shadow = shadow.filter(ImageFilter.GaussianBlur(10))
    draw.bitmap((0, 0), shadow)
    draw.text((x, y), text, font=font, fill=(238, 255, 248, 255))
    draw.text((x, y + 3), text, font=font, fill=(126, 239, 207, 82))


def render_reference_card() -> Path:
    bg_path = PROJECT_ROOT / "web" / "assets" / "brand-logo-hero-bg-dark.png"
    logo_path = PROJECT_ROOT / "web" / "assets" / "brand-logo.png"
    if not bg_path.exists():
        bg_path = PROJECT_ROOT / "web" / "assets" / "brand-logo-hero-dark.png"
    if not logo_path.exists():
        raise FileNotFoundError(f"Missing logo asset: {logo_path}")

    bg = cover_resize(Image.open(bg_path).convert("RGB"), WIDTH, HEIGHT)
    bg = bg.filter(ImageFilter.GaussianBlur(5))
    canvas = bg.convert("RGBA")

    overlay = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    overlay_pixels = overlay.load()
    for y in range(HEIGHT):
        darkness = int(82 + 80 * (y / HEIGHT))
        mint = int(18 * max(0, 1 - abs(y - 900) / 900))
        for x in range(WIDTH):
            overlay_pixels[x, y] = (3, 6 + mint, 16, darkness)
    canvas.alpha_composite(overlay)

    glow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    glow_draw = ImageDraw.Draw(glow)
    glow_draw.ellipse((180, 330, 900, 1050), fill=(74, 238, 199, 38))
    glow_draw.ellipse((300, 390, 780, 870), fill=(35, 92, 255, 56))
    glow = glow.filter(ImageFilter.GaussianBlur(58))
    canvas.alpha_composite(glow)

    draw = ImageDraw.Draw(canvas)
    for index in range(11):
        phase = index / 10
        y = int(465 + math.sin(phase * math.pi) * 255)
        x1 = int(100 + phase * 880)
        x2 = int(980 - phase * 880)
        color = (147, 255, 220, 36 + index * 4)
        draw.arc((min(x1, x2), y - 125, max(x1, x2), y + 125), 190, 350, fill=color, width=3)

    logo = Image.open(logo_path).convert("RGBA")
    logo_size = 520
    logo.thumbnail((logo_size, logo_size), Image.Resampling.LANCZOS)
    logo_x = (WIDTH - logo.width) // 2
    logo_y = 535
    logo_shadow = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    logo_shadow.alpha_composite(logo, (logo_x, logo_y))
    logo_shadow = logo_shadow.filter(ImageFilter.GaussianBlur(24))
    shadow_tint = Image.new("RGBA", (WIDTH, HEIGHT), (0, 0, 0, 0))
    shadow_tint.alpha_composite(logo_shadow)
    canvas.alpha_composite(shadow_tint)
    canvas.alpha_composite(logo, (logo_x, logo_y))

    font = load_font(92, bold=True)
    draw_centered_text(draw, "poliglotAI.online", 1220, font)

    footer_font = load_font(28, bold=False)
    footer = " "
    footer_bbox = draw.textbbox((0, 0), footer, font=footer_font)
    draw.text(((WIDTH - (footer_bbox[2] - footer_bbox[0])) // 2, 1520), footer, font=footer_font, fill=(255, 255, 255, 0))

    output_path = OUTPUT_DIR / OUTRO_REFERENCE_NAME
    output_path.parent.mkdir(parents=True, exist_ok=True)
    canvas.convert("RGB").save(output_path, quality=96)
    return output_path


def ensure_silent_audio() -> Path:
    COMFY_INPUT_DIR.mkdir(parents=True, exist_ok=True)
    audio_path = COMFY_INPUT_DIR / OUTRO_AUDIO_NAME
    sample_rate = 44100
    frame_count = int(sample_rate * OUTRO_DURATION_SECONDS)
    with wave.open(str(audio_path), "wb") as wav:
        wav.setnchannels(2)
        wav.setsampwidth(2)
        wav.setframerate(sample_rate)
        wav.writeframes(b"\x00\x00" * 2 * frame_count)
    return audio_path


def build_native_workflow() -> dict:
    workflow = copy.deepcopy(build_workflow())
    workflow["id"] = str(uuid.uuid4())
    for node in workflow["nodes"]:
        title = node.get("title", "")
        if title == "Product / UI reference":
            node["widgets_values"] = [OUTRO_REFERENCE_NAME, "image"]
        elif title == "FLUX keyframe prompt":
            node["widgets_values"] = [FLUX_OUTRO_PROMPT]
            node["mode"] = 2
        elif title == "Wan motion prompt":
            node["widgets_values"] = [WAN_OUTRO_PROMPT]
        elif title == "Negative conditioning":
            node["widgets_values"] = [OUTRO_NEGATIVE_PROMPT]
        elif title == "FLUX.2 Dev keyframe":
            node["mode"] = 2
        elif title == "Preview keyframe":
            node["mode"] = 2
        elif title == "Wan 2.1 I2V motion":
            node["widgets_values"] = [WAN_WIDTH, WAN_HEIGHT, WAN_LENGTH, 1]
        elif title == "Wan sampler":
            node["widgets_values"] = [624042002, "fixed", 16, 5.4, "uni_pc", "simple", 1.0]
        elif title == "XTTS voiceover input":
            node["widgets_values"] = [OUTRO_AUDIO_NAME]
        elif title == "Save MP4 480x832":
            node["widgets_values"][0] = "outro/poliglot_outro_2s_480_%date:yyyy-MM-dd%"
        elif title == "Save MP4 1080x1920":
            node["widgets_values"][0] = "outro/poliglot_outro_2s_1080_%date:yyyy-MM-dd%"
        elif title == "20s English ad script for XTTS":
            node["title"] = "2s outro direction"
            node["widgets_values"] = [OUTRO_NOTES]
        elif title == "Workflow notes":
            node["widgets_values"] = [OUTRO_NOTES]
        elif title == "Locking notes":
            node["widgets_values"] = [
                "## Outro locking strategy\n\n"
                f"- Start image: `{OUTRO_REFERENCE_NAME}` with exact logo and URL.\n"
                "- Default route sends the exact reference card directly into Wan 2.1 I2V.\n"
                "- FLUX.2 prompt/keyframe nodes are parked disabled for a creative-background variant.\n"
                "- Wan 2.1 I2V animates only a short 33-frame motion pass.\n"
                "- RIFE doubles frames for a smooth ~2-second vertical MP4.\n"
                "- If exact URL readability matters more than motion, overlay the reference URL in the editor after export."
            ]
    route_native_wan_from_reference(workflow)
    return workflow


def route_native_wan_from_reference(workflow: dict) -> None:
    """Bypass FLUX for the default native route so the URL remains readable."""
    nodes_by_id = {node["id"]: node for node in workflow["nodes"]}
    for link in workflow["links"]:
        if link[0] in {10, 15}:
            link[1] = 2

    node2_links = nodes_by_id[2]["outputs"][0].setdefault("links", [])
    for link_id in (10, 15):
        if link_id not in node2_links:
            node2_links.append(link_id)

    node4_links = nodes_by_id[4]["outputs"][0].setdefault("links", [])
    nodes_by_id[4]["outputs"][0]["links"] = [link_id for link_id in node4_links if link_id not in {10, 15}]


def build_api_prompt() -> dict:
    return {
        "1": {"class_type": "LoadImage", "inputs": {"image": OUTRO_REFERENCE_NAME}},
        "2": {
            "class_type": "ImageScale",
            "inputs": {"image": ["1", 0], "upscale_method": "lanczos", "width": WAN_WIDTH, "height": WAN_HEIGHT, "crop": "center"},
        },
        "17": {"class_type": "CLIPLoader", "inputs": {"clip_name": "umt5_xxl_fp8_e4m3fn_scaled.safetensors", "type": "wan", "device": "default"}},
        "18": {"class_type": "VAELoader", "inputs": {"vae_name": "wan_2.1_vae.safetensors"}},
        "19": {"class_type": "UNETLoader", "inputs": {"unet_name": "wan2.1_i2v_480p_14B_fp8_scaled.safetensors", "weight_dtype": "default"}},
        "20": {"class_type": "ModelSamplingSD3", "inputs": {"model": ["19", 0], "shift": 5.0}},
        "21": {"class_type": "CLIPTextEncode", "inputs": {"clip": ["17", 0], "text": WAN_OUTRO_PROMPT}},
        "22": {"class_type": "CLIPTextEncode", "inputs": {"clip": ["17", 0], "text": OUTRO_NEGATIVE_PROMPT}},
        "23": {"class_type": "CLIPVisionLoader", "inputs": {"clip_name": "clip_vision_h.safetensors"}},
        "24": {"class_type": "CLIPVisionEncode", "inputs": {"clip_vision": ["23", 0], "image": ["2", 0], "crop": "center"}},
        "25": {
            "class_type": "WanImageToVideo",
            "inputs": {
                "positive": ["21", 0],
                "negative": ["22", 0],
                "vae": ["18", 0],
                "clip_vision_output": ["24", 0],
                "start_image": ["2", 0],
                "width": WAN_WIDTH,
                "height": WAN_HEIGHT,
                "length": WAN_LENGTH,
                "batch_size": 1,
            },
        },
        "26": {
            "class_type": "KSampler",
            "inputs": {
                "model": ["20", 0],
                "positive": ["25", 0],
                "negative": ["25", 1],
                "latent_image": ["25", 2],
                "seed": 624042101,
                "steps": 12,
                "cfg": 5.2,
                "sampler_name": "uni_pc",
                "scheduler": "simple",
                "denoise": 1.0,
            },
        },
        "27": {"class_type": "VAEDecode", "inputs": {"samples": ["26", 0], "vae": ["18", 0]}},
        "28": {"class_type": "FrameInterpolationModelLoader", "inputs": {"model_name": "rife_v4.26.safetensors"}},
        "29": {"class_type": "FrameInterpolate", "inputs": {"interp_model": ["28", 0], "images": ["27", 0], "multiplier": RIFE_MULTIPLIER}},
        "30": {"class_type": "EmptyAudio", "inputs": {"duration": round(OUTRO_DURATION_SECONDS, 3), "sample_rate": 44100, "channels": 2}},
        "31": {"class_type": "CreateVideo", "inputs": {"images": ["29", 0], "fps": FPS, "audio": ["30", 0]}},
        "32": {
            "class_type": "SaveVideo",
            "inputs": {"video": ["31", 0], "filename_prefix": "outro/poliglot_outro_2s_api", "format": "mp4", "codec": "auto"},
        },
        "33": {"class_type": "SaveImage", "inputs": {"images": ["2", 0], "filename_prefix": "outro/poliglot_outro_2s_reference"}},
    }


def build_flux_wan_api_prompt() -> dict:
    return {
        "1": {"class_type": "LoadImage", "inputs": {"image": OUTRO_REFERENCE_NAME}},
        "2": {
            "class_type": "ImageScale",
            "inputs": {"image": ["1", 0], "upscale_method": "lanczos", "width": WAN_WIDTH, "height": WAN_HEIGHT, "crop": "center"},
        },
        "3": {"class_type": "UNETLoader", "inputs": {"unet_name": "flux2_dev_fp8mixed.safetensors", "weight_dtype": "default"}},
        "4": {"class_type": "CLIPLoader", "inputs": {"clip_name": "mistral_3_small_flux2_bf16.safetensors", "type": "flux2", "device": "default"}},
        "5": {"class_type": "VAELoader", "inputs": {"vae_name": "flux2-vae.safetensors"}},
        "6": {"class_type": "CLIPTextEncode", "inputs": {"clip": ["4", 0], "text": FLUX_OUTRO_PROMPT}},
        "7": {"class_type": "FluxGuidance", "inputs": {"conditioning": ["6", 0], "guidance": 3.8}},
        "8": {"class_type": "VAEEncode", "inputs": {"pixels": ["2", 0], "vae": ["5", 0]}},
        "9": {"class_type": "ReferenceLatent", "inputs": {"conditioning": ["7", 0], "latent": ["8", 0]}},
        "10": {"class_type": "BasicGuider", "inputs": {"model": ["3", 0], "conditioning": ["9", 0]}},
        "11": {"class_type": "RandomNoise", "inputs": {"noise_seed": 624042002}},
        "12": {"class_type": "KSamplerSelect", "inputs": {"sampler_name": "euler"}},
        "13": {"class_type": "Flux2Scheduler", "inputs": {"steps": 14, "width": WAN_WIDTH, "height": WAN_HEIGHT}},
        "14": {"class_type": "EmptyFlux2LatentImage", "inputs": {"width": WAN_WIDTH, "height": WAN_HEIGHT, "batch_size": 1}},
        "15": {
            "class_type": "SamplerCustomAdvanced",
            "inputs": {"noise": ["11", 0], "guider": ["10", 0], "sampler": ["12", 0], "sigmas": ["13", 0], "latent_image": ["14", 0]},
        },
        "16": {"class_type": "VAEDecode", "inputs": {"samples": ["15", 0], "vae": ["5", 0]}},
        "17": {"class_type": "CLIPLoader", "inputs": {"clip_name": "umt5_xxl_fp8_e4m3fn_scaled.safetensors", "type": "wan", "device": "default"}},
        "18": {"class_type": "VAELoader", "inputs": {"vae_name": "wan_2.1_vae.safetensors"}},
        "19": {"class_type": "UNETLoader", "inputs": {"unet_name": "wan2.1_i2v_480p_14B_fp8_scaled.safetensors", "weight_dtype": "default"}},
        "20": {"class_type": "ModelSamplingSD3", "inputs": {"model": ["19", 0], "shift": 5.0}},
        "21": {"class_type": "CLIPTextEncode", "inputs": {"clip": ["17", 0], "text": WAN_OUTRO_PROMPT}},
        "22": {"class_type": "CLIPTextEncode", "inputs": {"clip": ["17", 0], "text": OUTRO_NEGATIVE_PROMPT}},
        "23": {"class_type": "CLIPVisionLoader", "inputs": {"clip_name": "clip_vision_h.safetensors"}},
        "24": {"class_type": "CLIPVisionEncode", "inputs": {"clip_vision": ["23", 0], "image": ["16", 0], "crop": "center"}},
        "25": {
            "class_type": "WanImageToVideo",
            "inputs": {
                "positive": ["21", 0],
                "negative": ["22", 0],
                "vae": ["18", 0],
                "clip_vision_output": ["24", 0],
                "start_image": ["16", 0],
                "width": WAN_WIDTH,
                "height": WAN_HEIGHT,
                "length": WAN_LENGTH,
                "batch_size": 1,
            },
        },
        "26": {
            "class_type": "KSampler",
            "inputs": {
                "model": ["20", 0],
                "positive": ["25", 0],
                "negative": ["25", 1],
                "latent_image": ["25", 2],
                "seed": 624042101,
                "steps": 12,
                "cfg": 5.4,
                "sampler_name": "uni_pc",
                "scheduler": "simple",
                "denoise": 1.0,
            },
        },
        "27": {"class_type": "VAEDecode", "inputs": {"samples": ["26", 0], "vae": ["18", 0]}},
        "28": {"class_type": "FrameInterpolationModelLoader", "inputs": {"model_name": "rife_v4.26.safetensors"}},
        "29": {"class_type": "FrameInterpolate", "inputs": {"interp_model": ["28", 0], "images": ["27", 0], "multiplier": RIFE_MULTIPLIER}},
        "30": {"class_type": "EmptyAudio", "inputs": {"duration": round(OUTRO_DURATION_SECONDS, 3), "sample_rate": 44100, "channels": 2}},
        "31": {"class_type": "CreateVideo", "inputs": {"images": ["29", 0], "fps": FPS, "audio": ["30", 0]}},
        "32": {
            "class_type": "SaveVideo",
            "inputs": {"video": ["31", 0], "filename_prefix": "outro/poliglot_outro_2s_api", "format": "mp4", "codec": "auto"},
        },
        "33": {"class_type": "SaveImage", "inputs": {"images": ["16", 0], "filename_prefix": "outro/poliglot_outro_2s_keyframe"}},
    }


def write_prompt_doc() -> Path:
    path = OUTPUT_DIR / "poliglot_outro_2s_prompt.md"
    path.write_text(
        "\n".join(
            [
                "# Poliglot AI 2s Outro Prompt",
                "",
                "Use this insert at the end of TikTok, Reels, and Shorts videos.",
                "",
                "## Exact On-Screen Text",
                "",
                "`poliglotAI.online`",
                "",
                "## FLUX.2 Keyframe Prompt",
                "",
                "Optional creative-background variant. The default workflow keeps these FLUX.2 nodes disabled so the exact URL is not redrawn by the image model.",
                "",
                FLUX_OUTRO_PROMPT,
                "",
                "## Wan 2.1 I2V Motion Prompt",
                "",
                "Default text-safe route: animate the exact reference card directly with Wan 2.1 I2V.",
                "",
                WAN_OUTRO_PROMPT,
                "",
                "## Negative Prompt",
                "",
                OUTRO_NEGATIVE_PROMPT,
                "",
                "## Timing",
                "",
                f"- Wan length: `{WAN_LENGTH}` frames.",
                f"- RIFE multiplier: `{RIFE_MULTIPLIER}`.",
                f"- Output fps: `{FPS}`.",
                f"- Expected duration: `{OUTRO_DURATION_SECONDS:.2f}s`.",
                "",
                "## Files",
                "",
                f"- Native ComfyUI workflow: `{OUTRO_WORKFLOW_NAME}`.",
                f"- API prompt: `{OUTRO_API_NAME}`.",
                f"- Reference card: `{OUTRO_REFERENCE_NAME}`.",
                f"- Silent audio placeholder: `{OUTRO_AUDIO_NAME}` in ComfyUI input.",
                "",
                "## Production Note",
                "",
                "For the cleanest final brand insert, keep the generated motion subtle. The default API/native route animates the exact reference card directly. If the URL becomes less readable, use the generated clip as animated background and overlay the exact logo/URL from the reference card in the editor.",
            ]
        )
        + "\n",
        encoding="utf-8",
    )
    return path


def main() -> None:
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    COMFY_WORKFLOW_DIR.mkdir(parents=True, exist_ok=True)
    COMFY_INPUT_DIR.mkdir(parents=True, exist_ok=True)

    reference_path = render_reference_card()
    shutil.copyfile(reference_path, COMFY_INPUT_DIR / OUTRO_REFERENCE_NAME)
    audio_path = ensure_silent_audio()

    workflow = build_native_workflow()
    workflow_path = OUTPUT_DIR / OUTRO_WORKFLOW_NAME
    workflow_path.write_text(json.dumps(workflow, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    shutil.copyfile(workflow_path, COMFY_WORKFLOW_DIR / OUTRO_WORKFLOW_NAME)

    api_path = OUTPUT_DIR / OUTRO_API_NAME
    api_path.write_text(json.dumps(build_api_prompt(), ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    prompt_path = write_prompt_doc()

    print(f"Wrote {workflow_path}")
    print(f"Wrote {api_path}")
    print(f"Wrote {prompt_path}")
    print(f"Wrote {reference_path}")
    print(f"Copied reference to {COMFY_INPUT_DIR / OUTRO_REFERENCE_NAME}")
    print(f"Wrote {audio_path}")
    print(f"Copied workflow to {COMFY_WORKFLOW_DIR / OUTRO_WORKFLOW_NAME}")


if __name__ == "__main__":
    main()
