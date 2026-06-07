from __future__ import annotations

import copy
import json
import shutil
import uuid
import wave
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[1]
COMFY_BASE = Path(r"C:\Users\Admin\Documents\ComfyUI")
COMFY_BLUEPRINTS = Path(r"C:\Users\Admin\AppData\Local\Programs\ComfyUI\resources\ComfyUI\blueprints")
OUTPUT_DIR = PROJECT_ROOT / "comfyui_workflows"
WORKFLOW_NAME = "ad_shorts_reels_tiktok_native.json"
COMFY_WORKFLOW_DIR = COMFY_BASE / "user" / "default" / "workflows"
COMFY_INPUT_DIR = COMFY_BASE / "input"


FLUX_SUBGRAPH_ID = "41b0c117-7470-454c-914e-b8742dc06d62"


MODEL_URLS = {
    "wan_i2v": "https://huggingface.co/Comfy-Org/Wan_2.1_ComfyUI_repackaged/resolve/main/split_files/diffusion_models/wan2.1_i2v_480p_14B_fp8_scaled.safetensors",
    "wan_clip": "https://huggingface.co/Comfy-Org/Wan_2.1_ComfyUI_repackaged/resolve/main/split_files/text_encoders/umt5_xxl_fp8_e4m3fn_scaled.safetensors",
    "wan_vae": "https://huggingface.co/Comfy-Org/Wan_2.1_ComfyUI_repackaged/resolve/main/split_files/vae/wan_2.1_vae.safetensors",
    "wan_clip_vision": "https://huggingface.co/Comfy-Org/Wan_2.1_ComfyUI_repackaged/resolve/main/split_files/clip_vision/clip_vision_h.safetensors",
    "rife": "https://huggingface.co/Comfy-Org/frame_interpolation/resolve/main/frame_interpolation/rife_v4.26.safetensors",
}


FLUX_PROMPT = """Vertical 9:16 luxury social ad keyframe for Poliglot AI, an AI language learning product.
Use the project web assets as the product identity lock: Poliglot AI logo, premium web app look, browser/web version, language learning dashboard, pronunciation waveform, progress and practice cards.
Scene: a warm confident young adult woman, age 18 or older, wearing an elegant business suit, presenting the Poliglot AI web app on a laptop/tablet in a high-end modern office studio.
Premium technology advertisement, cinematic studio lighting, tasteful elite brand mood, clean composition for TikTok/Reels/Shorts, sharp details, product-first.
No invented UI labels, no fake small text, no watermark, no duplicated logo, no flags, no sensual pose."""


WAN_PROMPT = """20-second elite vertical product ad concept for Poliglot AI, generated as a premium 5-second shot in a 4-shot sequence. Smooth cinematic push-in and subtle parallax. A friendly young adult woman, age 18 or older, in a business suit confidently presents the Poliglot AI web version on a laptop/tablet. Show the website/web app visually, premium language-learning dashboard, voice waveform, pronunciation feedback, progress cards, Telegram and web continuity as elegant abstract product visuals. Luxury tech commercial, clean studio lighting, stable logo/product identity, polished social media advertisement, sharp focus, no captions inside the image."""


NEGATIVE_PROMPT = """minor, underage, childlike, sensual pose, fake letters, unreadable text, random subtitles, watermark, duplicated logo, distorted logo, broken UI, cluttered background, overexposed, low quality, jpeg artifacts, static frozen scene, deformed hands, deformed face, extra limbs, noisy details, blurry product, bad composition"""


XTTS_SCRIPT_PROMPT = """XTTS voice direction:
Young adult female speaker, age 18 or older. Elegant business tone, warm and confident, premium English tech advertisement. Clear pronunciation, calm luxury pace, no sensual delivery.

20-second English voiceover:
"Meet Poliglot AI, your private language coach on the web and in Telegram. Practice real conversations, hear clear corrections, improve pronunciation, and watch your progress grow every day. Learn faster, speak with confidence, and make every lesson feel personal. Poliglot AI. Your premium way to master languages."

Visual beat plan:
0-5s: premium web app hero on laptop/tablet, presenter enters frame.
5-10s: voice practice and pronunciation waveform.
10-15s: progress, vocabulary, mistakes, and Telegram/web continuity.
15-20s: confident presenter, logo/product reveal, premium CTA mood."""


TEST_VARIANTS = [
    {
        "slug": "01_web_hero_presenter",
        "title": "Desktop web app hero + presenter",
        "audio_file": "poliglot_test_01_web_hero_presenter.wav",
        "voiceover": "Meet Poliglot AI. Your private language coach works beautifully on the web and in Telegram, helping you practice, improve pronunciation, and keep progress moving every day.",
        "flux_prompt": """Vertical 9:16 luxury ad keyframe for Poliglot AI desktop web app. A warm confident young adult woman, age 18 or older, in an elegant business suit presents the Poliglot AI web version on a laptop in a premium modern office studio. The laptop screen clearly shows a beautiful desktop web interface: sidebar, learning cards, progress panel, practice area, premium dashboard, voice waveform. Use the project web assets as brand identity: logo, dark premium dashboard mood, app background. Elite SaaS advertisement, product-first, clean cinematic lighting, tasteful and expensive, no fake text, no sensual pose.""",
        "wan_prompt": """5-second vertical luxury SaaS ad shot. Smooth cinematic push-in toward a laptop showing the Poliglot AI desktop web interface while a friendly young adult woman, age 18 or older, in a business suit presents it with confident calm energy. Premium office studio, polished functional UI visible: sidebar, learning cards, progress, practice, voice waveform, stable logo/product identity, no captions inside the image.""",
    },
    {
        "slug": "02_voice_pronunciation",
        "title": "Mobile web voice practice + pronunciation",
        "audio_file": "poliglot_test_02_voice_pronunciation.wav",
        "voiceover": "Speak, listen, and improve with instant pronunciation feedback. Poliglot AI helps you hear mistakes clearly and turn every practice session into confident progress.",
        "flux_prompt": """Vertical 9:16 premium ad keyframe for Poliglot AI mobile web voice practice. A warm confident young adult woman, age 18 or older, in a tailored business suit holds a smartphone with the Poliglot AI mobile web interface clearly visible. The phone screen shows a polished mobile app shell: compact header, voice practice card, microphone control, waveform, pronunciation feedback, progress chip. Luxury tech commercial, high-end studio lighting, functional product UI visible, no fake words, no random subtitles, no sensual pose.""",
        "wan_prompt": """5-second vertical premium product shot. Slow parallax around the presenter and smartphone, showing the Poliglot AI mobile web interface with voice practice, waveform motion, pronunciation feedback rings, and progress energy. Young adult business presenter, age 18 or older, friendly and professional. Elite language-learning SaaS commercial, stable product identity, sharp polished lighting, no captions inside the image.""",
    },
]


NOTE_MAIN = """## Ad shorts pipeline

Default chain:
Load product reference -> 480x832 crop -> FLUX.2 Dev keyframe -> Wan 2.1 I2V motion -> VAE Decode -> native RIFE 4.26 frame interpolation -> Create Video + voice -> Save MP4.

The 1080x1920 branch is optional and uses the existing RealESRGAN x4 model. Keep it disabled if VRAM or time is tight.

XTTS and Whisper are kept as external/pre-post lanes by default so this workflow loads without custom nodes. Put the generated voice file in ComfyUI/input/ad_voiceover.wav before running. Generate subtitles with Whisper after export.

For a 20-second ad, generate 4 approved 5-second shots from the beat plan, then stitch them. This is usually cleaner and lighter than one very long Wan generation."""


NOTE_CONTROL = """## Locking strategy

Product/UI lock:
- FLUX.2 uses the input reference through its native reference latent.
- Wan 2.1 receives the FLUX keyframe as both start_image and CLIP Vision conditioning.

Structure lock:
- Use ad_structure_reference.png as the composition reference when preparing the product image.
- A native ControlNet loader is parked here for SD3.5/SDXL sidecar work, but the main Flux.2 + Wan path stays custom-node-free and disk-conscious."""


NOTE_MODELS = """## Required model files

Already present on this machine:
- diffusion_models/flux2_dev_fp8mixed.safetensors
- text_encoders/mistral_3_small_flux2_bf16.safetensors
- vae/flux2-vae.safetensors
- upscale_models/RealESRGAN_x4plus.pth

Install missing compact video files with:
PowerShell -ExecutionPolicy Bypass -File .\\comfyui_workflows\\install_ad_short_models.ps1

This uses Wan 2.1 480p 14B FP8 scaled instead of fp16 to save about 16 GB."""


def load_flux_definition() -> dict:
    source = COMFY_BLUEPRINTS / "Image Edit (Flux.2 Dev).json"
    data = json.loads(source.read_text(encoding="utf-8"))
    subgraph = copy.deepcopy(data["definitions"]["subgraphs"][0])
    for node in subgraph["nodes"]:
        widgets = node.get("widgets_values")
        if not widgets:
            continue
        for index, value in enumerate(widgets):
            if value == "full_encoder_small_decoder.safetensors":
                widgets[index] = "flux2-vae.safetensors"
        if node.get("type") == "VAELoader":
            models = node.setdefault("properties", {}).setdefault("models", [])
            models[:] = [
                {
                    "name": "flux2-vae.safetensors",
                    "url": "https://huggingface.co/Comfy-Org/flux2-dev/resolve/main/split_files/vae/flux2-vae.safetensors",
                    "directory": "vae",
                }
            ]
    return subgraph


class WorkflowBuilder:
    def __init__(self) -> None:
        self.nodes: list[dict] = []
        self.links: list[list] = []
        self.groups: list[dict] = []
        self.last_node_id = 0
        self.last_link_id = 0

    def node(
        self,
        node_id: int,
        node_type: str,
        pos: tuple[int, int],
        size: tuple[int, int],
        inputs: list[dict],
        outputs: list[dict],
        widgets_values: list | None = None,
        title: str | None = None,
        properties: dict | None = None,
        color: str | None = None,
        bgcolor: str | None = None,
        flags: dict | None = None,
        mode: int = 0,
    ) -> dict:
        n = {
            "id": node_id,
            "type": node_type,
            "pos": list(pos),
            "size": list(size),
            "flags": flags or {},
            "order": len(self.nodes),
            "mode": mode,
            "inputs": inputs,
            "outputs": outputs,
            "properties": properties if properties is not None else {"Node name for S&R": node_type},
            "widgets_values": widgets_values or [],
        }
        if title:
            n["title"] = title
        if color:
            n["color"] = color
        if bgcolor:
            n["bgcolor"] = bgcolor
        self.nodes.append(n)
        self.last_node_id = max(self.last_node_id, node_id)
        return n

    def group(self, group_id: int, title: str, x: int, y: int, w: int, h: int, color: str = "#3f789e") -> None:
        self.groups.append(
            {
                "id": group_id,
                "title": title,
                "bounding": [x, y, w, h],
                "color": color,
                "font_size": 24,
                "flags": {},
            }
        )

    def connect(self, source: int, source_slot: int, target: int, target_slot: int, link_type: str) -> None:
        self.last_link_id += 1
        link_id = self.last_link_id
        self.links.append([link_id, source, source_slot, target, target_slot, link_type])

        source_node = next(n for n in self.nodes if n["id"] == source)
        target_node = next(n for n in self.nodes if n["id"] == target)
        source_links = source_node["outputs"][source_slot].setdefault("links", [])
        if source_links is None:
            source_node["outputs"][source_slot]["links"] = [link_id]
        else:
            source_links.append(link_id)
        target_node["inputs"][target_slot]["link"] = link_id


def input_slot(name: str, slot_type: str, widget: bool = False, optional: bool = False) -> dict:
    item = {"name": name, "type": slot_type, "link": None}
    if widget:
        item["widget"] = {"name": name}
    if optional:
        item["shape"] = 7
    return item


def output_slot(name: str, slot_type: str, slot_index: int | None = None) -> dict:
    item = {"name": name, "type": slot_type, "links": []}
    if slot_index is not None:
        item["slot_index"] = slot_index
    return item


def model_prop(node_name: str, models: list[dict]) -> dict:
    return {"Node name for S&R": node_name, "models": models}


def build_workflow() -> dict:
    b = WorkflowBuilder()

    b.group(1, "1 Product reference + prompts", -980, -520, 920, 850)
    b.group(2, "2 FLUX.2 Dev keyframe", -20, -520, 520, 610)
    b.group(3, "3 Wan 2.1 I2V motion", 560, -520, 980, 850)
    b.group(4, "4 RIFE smoothing + voice", 1600, -520, 720, 640)
    b.group(5, "5 Optional 1080x1920 upscale", 2380, -520, 800, 640)
    b.group(6, "Notes + structure/control lane", -980, 400, 2260, 720, "#7a5b2d")

    b.node(
        1,
        "LoadImage",
        (-920, -430),
        (330, 340),
        [input_slot("image", "COMBO", True), input_slot("upload", "IMAGEUPLOAD", True)],
        [output_slot("IMAGE", "IMAGE", 0), output_slot("MASK", "MASK")],
        ["poliglot_web_product_reference.png", "image"],
        "Product / UI reference",
    )
    b.node(
        2,
        "ImageScale",
        (-540, -405),
        (315, 130),
        [
            input_slot("image", "IMAGE"),
            input_slot("upscale_method", "COMBO", True),
            input_slot("width", "INT", True),
            input_slot("height", "INT", True),
            input_slot("crop", "COMBO", True),
        ],
        [output_slot("IMAGE", "IMAGE", 0)],
        ["lanczos", 480, 832, "center"],
        "Vertical crop 480x832",
    )
    b.node(
        3,
        "PrimitiveStringMultiline",
        (-920, -40),
        (695, 230),
        [input_slot("value", "STRING", True)],
        [output_slot("STRING", "STRING", 0)],
        [FLUX_PROMPT],
        "FLUX keyframe prompt",
        color="#232",
        bgcolor="#353",
    )
    b.node(
        4,
        FLUX_SUBGRAPH_ID,
        (60, -430),
        (400, 260),
        [
            input_slot("pixels", "IMAGE"),
            input_slot("text", "STRING", True),
            input_slot("unet_name", "COMBO", True),
            input_slot("clip_name", "COMBO", True),
            input_slot("vae_name", "COMBO", True),
            input_slot("value", "BOOLEAN", True),
            input_slot("lora_name", "COMBO", True),
        ],
        [output_slot("IMAGE", "IMAGE", 0)],
        [],
        "FLUX.2 Dev keyframe",
        properties={
            "proxyWidgets": [
                ["123", "text"],
                ["129", "unet_name"],
                ["124", "clip_name"],
                ["121", "vae_name"],
                ["138", "value"],
                ["128", "lora_name"],
                ["125", "noise_seed"],
                ["125", "control_after_generate"],
            ],
            "Node name for S&R": "Image Edit (Flux.2 Dev)",
        },
    )
    b.node(
        5,
        "PreviewImage",
        (65, -115),
        (320, 270),
        [input_slot("images", "IMAGE")],
        [],
        [],
        "Preview keyframe",
    )

    b.node(
        6,
        "PrimitiveStringMultiline",
        (610, -430),
        (500, 180),
        [input_slot("value", "STRING", True)],
        [output_slot("STRING", "STRING", 0)],
        [WAN_PROMPT],
        "Wan motion prompt",
        color="#232",
        bgcolor="#353",
    )
    b.node(
        7,
        "CLIPLoader",
        (610, -190),
        (360, 106),
        [input_slot("clip_name", "COMBO", True), input_slot("type", "COMBO", True), input_slot("device", "COMBO", True, True)],
        [output_slot("CLIP", "CLIP", 0)],
        ["umt5_xxl_fp8_e4m3fn_scaled.safetensors", "wan", "default"],
        "Load Wan text encoder",
        model_prop(
            "CLIPLoader",
            [
                {
                    "name": "umt5_xxl_fp8_e4m3fn_scaled.safetensors",
                    "url": MODEL_URLS["wan_clip"],
                    "directory": "text_encoders",
                }
            ],
        ),
    )
    b.node(
        8,
        "VAELoader",
        (610, -45),
        (360, 60),
        [input_slot("vae_name", "COMBO", True)],
        [output_slot("VAE", "VAE", 0)],
        ["wan_2.1_vae.safetensors"],
        "Load Wan VAE",
        model_prop("VAELoader", [{"name": "wan_2.1_vae.safetensors", "url": MODEL_URLS["wan_vae"], "directory": "vae"}]),
    )
    b.node(
        9,
        "UNETLoader",
        (610, 80),
        (360, 82),
        [input_slot("unet_name", "COMBO", True), input_slot("weight_dtype", "COMBO", True)],
        [output_slot("MODEL", "MODEL", 0)],
        ["wan2.1_i2v_480p_14B_fp8_scaled.safetensors", "default"],
        "Load Wan 2.1 I2V FP8",
        model_prop(
            "UNETLoader",
            [
                {
                    "name": "wan2.1_i2v_480p_14B_fp8_scaled.safetensors",
                    "url": MODEL_URLS["wan_i2v"],
                    "directory": "diffusion_models",
                }
            ],
        ),
    )
    b.node(
        10,
        "ModelSamplingSD3",
        (1015, 80),
        (245, 58),
        [input_slot("model", "MODEL"), input_slot("shift", "FLOAT", True)],
        [output_slot("MODEL", "MODEL", 0)],
        [5.0],
        "Wan sampling shift",
    )
    b.node(
        11,
        "CLIPTextEncode",
        (1015, -430),
        (480, 120),
        [input_slot("clip", "CLIP"), input_slot("text", "STRING", True)],
        [output_slot("CONDITIONING", "CONDITIONING", 0)],
        [""],
        "Positive conditioning",
        color="#232",
        bgcolor="#353",
    )
    b.node(
        12,
        "CLIPTextEncode",
        (1015, -250),
        (480, 140),
        [input_slot("clip", "CLIP"), input_slot("text", "STRING", True)],
        [output_slot("CONDITIONING", "CONDITIONING", 0)],
        [NEGATIVE_PROMPT],
        "Negative conditioning",
        color="#322",
        bgcolor="#533",
    )
    b.node(
        13,
        "CLIPVisionLoader",
        (1015, -45),
        (315, 58),
        [input_slot("clip_name", "COMBO", True)],
        [output_slot("CLIP_VISION", "CLIP_VISION", 0)],
        ["clip_vision_h.safetensors"],
        "IPAdapter-style image lock",
        model_prop(
            "CLIPVisionLoader",
            [{"name": "clip_vision_h.safetensors", "url": MODEL_URLS["wan_clip_vision"], "directory": "clip_vision"}],
        ),
    )
    b.node(
        14,
        "CLIPVisionEncode",
        (1015, 180),
        (330, 78),
        [input_slot("clip_vision", "CLIP_VISION"), input_slot("image", "IMAGE"), input_slot("crop", "COMBO", True)],
        [output_slot("CLIP_VISION_OUTPUT", "CLIP_VISION_OUTPUT", 0)],
        ["center"],
        "Encode keyframe lock",
    )
    b.node(
        15,
        "WanImageToVideo",
        (1350, -40),
        (360, 230),
        [
            input_slot("positive", "CONDITIONING"),
            input_slot("negative", "CONDITIONING"),
            input_slot("vae", "VAE"),
            input_slot("clip_vision_output", "CLIP_VISION_OUTPUT", False, True),
            input_slot("start_image", "IMAGE", False, True),
            input_slot("width", "INT", True),
            input_slot("height", "INT", True),
            input_slot("length", "INT", True),
            input_slot("batch_size", "INT", True),
        ],
        [output_slot("positive", "CONDITIONING", 0), output_slot("negative", "CONDITIONING", 1), output_slot("latent", "LATENT", 2)],
        [480, 832, 81, 1],
        "Wan 2.1 I2V motion",
    )
    b.node(
        16,
        "KSampler",
        (1740, -40),
        (315, 262),
        [
            input_slot("model", "MODEL"),
            input_slot("positive", "CONDITIONING"),
            input_slot("negative", "CONDITIONING"),
            input_slot("latent_image", "LATENT"),
            input_slot("seed", "INT", True),
            input_slot("steps", "INT", True),
            input_slot("cfg", "FLOAT", True),
            input_slot("sampler_name", "COMBO", True),
            input_slot("scheduler", "COMBO", True),
            input_slot("denoise", "FLOAT", True),
        ],
        [output_slot("LATENT", "LATENT", 0)],
        [431231234567890, "randomize", 20, 6.0, "uni_pc", "simple", 1.0],
        "Wan sampler",
    )
    b.node(
        17,
        "VAEDecode",
        (2090, -15),
        (210, 46),
        [input_slot("samples", "LATENT"), input_slot("vae", "VAE")],
        [output_slot("IMAGE", "IMAGE", 0)],
        [],
        "VAE Decode",
    )

    b.node(
        18,
        "FrameInterpolationModelLoader",
        (1645, -430),
        (360, 58),
        [input_slot("model_name", "COMBO", True)],
        [output_slot("INTERP_MODEL", "INTERP_MODEL", 0)],
        ["rife_v4.26.safetensors"],
        "Load RIFE",
        model_prop("FrameInterpolationModelLoader", [{"name": "rife_v4.26.safetensors", "url": MODEL_URLS["rife"], "directory": "frame_interpolation"}]),
    )
    b.node(
        19,
        "FrameInterpolate",
        (1645, -310),
        (315, 98),
        [input_slot("interp_model", "INTERP_MODEL"), input_slot("images", "IMAGE"), input_slot("multiplier", "INT", True)],
        [output_slot("IMAGE", "IMAGE", 0)],
        [2],
        "RIFE smoothing x2",
    )
    b.node(
        20,
        "LoadAudio",
        (1645, -165),
        (315, 58),
        [input_slot("audio", "COMBO", True)],
        [output_slot("AUDIO", "AUDIO", 0)],
        ["ad_voiceover.wav"],
        "XTTS voiceover input",
    )
    b.node(
        21,
        "CreateVideo",
        (1645, -45),
        (270, 78),
        [input_slot("images", "IMAGE"), input_slot("audio", "AUDIO", False, True), input_slot("fps", "FLOAT", True)],
        [output_slot("VIDEO", "VIDEO", 0)],
        [32],
        "Video Combine 480x832",
    )
    b.node(
        22,
        "SaveVideo",
        (1980, -45),
        (315, 112),
        [input_slot("video", "VIDEO"), input_slot("filename_prefix", "STRING", True), input_slot("format", "COMBO", True), input_slot("codec", "COMBO", True)],
        [],
        ["ad_shorts/poliglot_reel_480_%date:yyyy-MM-dd%", "mp4", "auto"],
        "Save MP4 480x832",
    )

    b.node(
        23,
        "UpscaleModelLoader",
        (2425, -430),
        (360, 58),
        [input_slot("model_name", "COMBO", True)],
        [output_slot("UPSCALE_MODEL", "UPSCALE_MODEL", 0)],
        ["RealESRGAN_x4plus.pth"],
        "Load upscale model",
        model_prop("UpscaleModelLoader", [{"name": "RealESRGAN_x4plus.pth", "directory": "upscale_models"}]),
        mode=2,
    )
    b.node(
        24,
        "ImageUpscaleWithModel",
        (2425, -310),
        (315, 78),
        [input_slot("upscale_model", "UPSCALE_MODEL"), input_slot("image", "IMAGE")],
        [output_slot("IMAGE", "IMAGE", 0)],
        [],
        "Upscale frames",
        mode=2,
    )
    b.node(
        25,
        "ImageScale",
        (2425, -185),
        (315, 130),
        [
            input_slot("image", "IMAGE"),
            input_slot("upscale_method", "COMBO", True),
            input_slot("width", "INT", True),
            input_slot("height", "INT", True),
            input_slot("crop", "COMBO", True),
        ],
        [output_slot("IMAGE", "IMAGE", 0)],
        ["lanczos", 1080, 1920, "center"],
        "Final 1080x1920",
        mode=2,
    )
    b.node(
        26,
        "CreateVideo",
        (2810, -185),
        (270, 78),
        [input_slot("images", "IMAGE"), input_slot("audio", "AUDIO", False, True), input_slot("fps", "FLOAT", True)],
        [output_slot("VIDEO", "VIDEO", 0)],
        [32],
        "Video Combine 1080p",
        mode=2,
    )
    b.node(
        27,
        "SaveVideo",
        (2810, -45),
        (315, 112),
        [input_slot("video", "VIDEO"), input_slot("filename_prefix", "STRING", True), input_slot("format", "COMBO", True), input_slot("codec", "COMBO", True)],
        [],
        ["ad_shorts/poliglot_reel_1080_%date:yyyy-MM-dd%", "mp4", "auto"],
        "Save MP4 1080x1920",
        mode=2,
    )

    b.node(
        28,
        "LoadImage",
        (-920, 485),
        (330, 220),
        [input_slot("image", "COMBO", True), input_slot("upload", "IMAGEUPLOAD", True)],
        [output_slot("IMAGE", "IMAGE", 0), output_slot("MASK", "MASK")],
        ["poliglot_web_structure_reference.png", "image"],
        "Structure reference",
    )
    b.node(
        29,
        "ControlNetLoader",
        (-540, 485),
        (360, 58),
        [input_slot("control_net_name", "COMBO", True)],
        [output_slot("CONTROL_NET", "CONTROL_NET", 0)],
        ["sd3.5_large_controlnet_depth.safetensors"],
        "Optional ControlNet sidecar",
        model_prop(
            "ControlNetLoader",
            [{"name": "sd3.5_large_controlnet_depth.safetensors", "directory": "controlnet"}],
        ),
        mode=2,
    )
    b.node(
        30,
        "MarkdownNote",
        (-130, 485),
        (640, 300),
        [],
        [],
        [NOTE_CONTROL],
        "Locking notes",
        {},
        "#432",
        "#000",
    )
    b.node(
        31,
        "MarkdownNote",
        (-920, 735),
        (640, 210),
        [],
        [],
        [NOTE_MODELS],
        "Model notes",
        {},
        "#432",
        "#000",
    )
    b.node(
        32,
        "MarkdownNote",
        (545, 485),
        (720, 300),
        [],
        [],
        [NOTE_MAIN],
        "Workflow notes",
        {},
        "#432",
        "#000",
    )
    b.node(
        33,
        "PrimitiveStringMultiline",
        (545, 825),
        (720, 250),
        [input_slot("value", "STRING", True)],
        [output_slot("STRING", "STRING", 0)],
        [XTTS_SCRIPT_PROMPT],
        "20s English ad script for XTTS",
        color="#232",
        bgcolor="#353",
    )

    b.connect(1, 0, 2, 0, "IMAGE")
    b.connect(2, 0, 4, 0, "IMAGE")
    b.connect(3, 0, 4, 1, "STRING")
    b.connect(4, 0, 5, 0, "IMAGE")
    b.connect(6, 0, 11, 1, "STRING")
    b.connect(7, 0, 11, 0, "CLIP")
    b.connect(7, 0, 12, 0, "CLIP")
    b.connect(9, 0, 10, 0, "MODEL")
    b.connect(13, 0, 14, 0, "CLIP_VISION")
    b.connect(4, 0, 14, 1, "IMAGE")
    b.connect(11, 0, 15, 0, "CONDITIONING")
    b.connect(12, 0, 15, 1, "CONDITIONING")
    b.connect(8, 0, 15, 2, "VAE")
    b.connect(14, 0, 15, 3, "CLIP_VISION_OUTPUT")
    b.connect(4, 0, 15, 4, "IMAGE")
    b.connect(10, 0, 16, 0, "MODEL")
    b.connect(15, 0, 16, 1, "CONDITIONING")
    b.connect(15, 1, 16, 2, "CONDITIONING")
    b.connect(15, 2, 16, 3, "LATENT")
    b.connect(16, 0, 17, 0, "LATENT")
    b.connect(8, 0, 17, 1, "VAE")
    b.connect(18, 0, 19, 0, "INTERP_MODEL")
    b.connect(17, 0, 19, 1, "IMAGE")
    b.connect(19, 0, 21, 0, "IMAGE")
    b.connect(20, 0, 21, 1, "AUDIO")
    b.connect(21, 0, 22, 0, "VIDEO")
    b.connect(23, 0, 24, 0, "UPSCALE_MODEL")
    b.connect(19, 0, 24, 1, "IMAGE")
    b.connect(24, 0, 25, 0, "IMAGE")
    b.connect(25, 0, 26, 0, "IMAGE")
    b.connect(20, 0, 26, 1, "AUDIO")
    b.connect(26, 0, 27, 0, "VIDEO")

    return {
        "id": str(uuid.uuid4()),
        "revision": 0,
        "last_node_id": b.last_node_id,
        "last_link_id": b.last_link_id,
        "nodes": b.nodes,
        "links": b.links,
        "groups": b.groups,
        "config": {},
        "extra": {
            "ds": {"scale": 0.55, "offset": [1040, 520]},
            "frontendVersion": "1.29.2",
            "workflowRendererVersion": "LG",
        },
        "version": 0.4,
        "definitions": {"subgraphs": [load_flux_definition()]},
    }


def build_api_prompt(variant: dict | None = None) -> dict:
    wan_prompt = (variant or {}).get("wan_prompt", WAN_PROMPT)
    audio_file = (variant or {}).get("audio_file", "ad_voiceover.wav")
    save_slug = (variant or {}).get("slug", "api")
    return {
        "1": {"class_type": "LoadImage", "inputs": {"image": "poliglot_web_product_reference.png"}},
        "2": {
            "class_type": "ImageScale",
            "inputs": {"image": ["1", 0], "upscale_method": "lanczos", "width": 480, "height": 832, "crop": "center"},
        },
        "3": {"class_type": "CLIPLoader", "inputs": {"clip_name": "umt5_xxl_fp8_e4m3fn_scaled.safetensors", "type": "wan", "device": "default"}},
        "4": {"class_type": "VAELoader", "inputs": {"vae_name": "wan_2.1_vae.safetensors"}},
        "5": {"class_type": "UNETLoader", "inputs": {"unet_name": "wan2.1_i2v_480p_14B_fp8_scaled.safetensors", "weight_dtype": "default"}},
        "6": {"class_type": "ModelSamplingSD3", "inputs": {"model": ["5", 0], "shift": 5.0}},
        "7": {"class_type": "CLIPTextEncode", "inputs": {"clip": ["3", 0], "text": wan_prompt}},
        "8": {"class_type": "CLIPTextEncode", "inputs": {"clip": ["3", 0], "text": NEGATIVE_PROMPT}},
        "9": {"class_type": "CLIPVisionLoader", "inputs": {"clip_name": "clip_vision_h.safetensors"}},
        "10": {"class_type": "CLIPVisionEncode", "inputs": {"clip_vision": ["9", 0], "image": ["2", 0], "crop": "center"}},
        "11": {
            "class_type": "WanImageToVideo",
            "inputs": {
                "positive": ["7", 0],
                "negative": ["8", 0],
                "vae": ["4", 0],
                "clip_vision_output": ["10", 0],
                "start_image": ["2", 0],
                "width": 480,
                "height": 832,
                "length": 81,
                "batch_size": 1,
            },
        },
        "12": {
            "class_type": "KSampler",
            "inputs": {
                "model": ["6", 0],
                "positive": ["11", 0],
                "negative": ["11", 1],
                "latent_image": ["11", 2],
                "seed": 431231234567890,
                "steps": 20,
                "cfg": 6.0,
                "sampler_name": "uni_pc",
                "scheduler": "simple",
                "denoise": 1.0,
            },
        },
        "13": {"class_type": "VAEDecode", "inputs": {"samples": ["12", 0], "vae": ["4", 0]}},
        "14": {"class_type": "FrameInterpolationModelLoader", "inputs": {"model_name": "rife_v4.26.safetensors"}},
        "15": {"class_type": "FrameInterpolate", "inputs": {"interp_model": ["14", 0], "images": ["13", 0], "multiplier": 2}},
        "16": {"class_type": "LoadAudio", "inputs": {"audio": audio_file}},
        "17": {"class_type": "CreateVideo", "inputs": {"images": ["15", 0], "audio": ["16", 0], "fps": 32}},
        "18": {
            "class_type": "SaveVideo",
            "inputs": {"video": ["17", 0], "filename_prefix": f"ad_shorts/{save_slug}_api_%date:yyyy-MM-dd%", "format": "mp4", "codec": "auto"},
        },
    }


def variant_script_prompt(variant: dict) -> str:
    return f"""XTTS voice direction:
Young adult female speaker, age 18 or older. Elegant business tone, warm and confident, premium English tech advertisement. Clear pronunciation, calm luxury pace, no sensual delivery.

5-second English voiceover:
"{variant["voiceover"]}"

Visual focus:
{variant["title"]}
"""


def apply_variant(workflow: dict, variant: dict) -> dict:
    workflow["id"] = str(uuid.uuid4())
    for node in workflow["nodes"]:
        title = node.get("title", "")
        if title == "FLUX keyframe prompt":
            node["widgets_values"] = [variant["flux_prompt"]]
        elif title == "Wan motion prompt":
            node["widgets_values"] = [variant["wan_prompt"]]
        elif title == "XTTS voiceover input":
            node["widgets_values"] = [variant["audio_file"]]
        elif title == "20s English ad script for XTTS":
            node["title"] = f"5s script: {variant['title']}"
            node["widgets_values"] = [variant_script_prompt(variant)]
        elif title == "Save MP4 480x832":
            node["widgets_values"][0] = f"ad_shorts/{variant['slug']}_480_%date:yyyy-MM-dd%"
        elif title == "Save MP4 1080x1920":
            node["widgets_values"][0] = f"ad_shorts/{variant['slug']}_1080_%date:yyyy-MM-dd%"
    return workflow


def write_test_prompt_pack() -> Path:
    lines = ["# Poliglot AI 5s Test Ad Variants", ""]
    lines.append("Generate these first before building the full 12-clip pack.")
    lines.append("")
    for index, variant in enumerate(TEST_VARIANTS, 1):
        lines.extend(
            [
                f"## {index}. {variant['title']}",
                "",
                f"- Workflow: `test_variants/ad_test_{variant['slug']}.json`",
                f"- Voice file placeholder: `{variant['audio_file']}`",
                "",
                "### Voiceover",
                "",
                variant["voiceover"],
                "",
                "### FLUX Keyframe Prompt",
                "",
                variant["flux_prompt"],
                "",
                "### Wan Motion Prompt",
                "",
                variant["wan_prompt"],
                "",
            ]
        )
    path = OUTPUT_DIR / "poliglot_2x5s_test_ad_variants.md"
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")
    return path


def ensure_silent_audio(filename: str, seconds: int) -> None:
    audio_path = COMFY_INPUT_DIR / filename
    if audio_path.exists():
        return
    sample_rate = 44100
    with wave.open(str(audio_path), "wb") as wav:
        wav.setnchannels(2)
        wav.setsampwidth(2)
        wav.setframerate(sample_rate)
        wav.writeframes(b"\x00\x00" * 2 * sample_rate * seconds)


def ensure_input_assets() -> None:
    COMFY_INPUT_DIR.mkdir(parents=True, exist_ok=True)
    asset_dir = PROJECT_ROOT / "web" / "assets"
    product_source = asset_dir / "brand-logo-hero-dark.png"
    structure_source = asset_dir / "app-background-dark.png"
    fallback = asset_dir / "brand-logo.png"

    for source, filename in [
        (product_source if product_source.exists() else fallback, "poliglot_web_product_reference.png"),
        (structure_source if structure_source.exists() else fallback, "poliglot_web_structure_reference.png"),
        (fallback, "ad_product_reference.png"),
        (fallback, "ad_structure_reference.png"),
    ]:
        if source.exists():
            target = COMFY_INPUT_DIR / filename
            if not target.exists():
                shutil.copyfile(source, target)

    ensure_silent_audio("ad_voiceover.wav", 8)
    for variant in TEST_VARIANTS:
        ensure_silent_audio(variant["audio_file"], 6)


def main() -> None:
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    COMFY_WORKFLOW_DIR.mkdir(parents=True, exist_ok=True)
    ensure_input_assets()

    workflow = build_workflow()
    workflow_path = OUTPUT_DIR / WORKFLOW_NAME
    workflow_path.write_text(json.dumps(workflow, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    api_prompt_path = OUTPUT_DIR / "ad_shorts_reels_tiktok_wan_stage_api.json"
    api_prompt_path.write_text(json.dumps(build_api_prompt(), ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    comfy_copy = COMFY_WORKFLOW_DIR / WORKFLOW_NAME
    shutil.copyfile(workflow_path, comfy_copy)

    test_dir = OUTPUT_DIR / "test_variants"
    test_dir.mkdir(parents=True, exist_ok=True)
    prompt_pack_path = write_test_prompt_pack()
    for variant in TEST_VARIANTS:
        variant_workflow = apply_variant(copy.deepcopy(workflow), variant)
        variant_path = test_dir / f"ad_test_{variant['slug']}.json"
        variant_path.write_text(json.dumps(variant_workflow, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
        shutil.copyfile(variant_path, COMFY_WORKFLOW_DIR / variant_path.name)

    print(f"Wrote {workflow_path}")
    print(f"Wrote {api_prompt_path}")
    print(f"Wrote {prompt_pack_path}")
    print(f"Wrote {len(TEST_VARIANTS)} test workflows in {test_dir}")
    print(f"Copied {comfy_copy}")
    print(f"Ensured input assets in {COMFY_INPUT_DIR}")


if __name__ == "__main__":
    main()
