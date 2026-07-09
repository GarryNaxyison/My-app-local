# Neriva Dzen Article Visual Pack

Article: `Почему вы учите английский годами, но не говорите? Как ИИ меняет правила игры`

Goal: 9 ComfyUI-ready image directions for a Yandex Zen article about Neriva.ru. Each of the 3 article blocks gets 3 materially different variants, so the final publication can choose the strongest visual rhythm instead of small color tweaks.

## Recommended Generation Setup

- Workflow role: still image keyframe / text-to-image pass.
- Preferred model chain: existing local ComfyUI FLUX.2 still-image branch.
- Draft size: `1536x864`.
- Final size: `1920x1080` or higher 16:9 export after the favorite variant is selected.
- Crop safety: keep the main subject inside the central 70% of the frame; leave some calm negative space for Yandex Zen title/crop behavior.
- Text policy: do not ask the model to render readable words, app names, subtitles, or UI text. Add article headline and brand text outside the generated image if needed.
- People policy: when a person appears, make the subject an adult, age 25+, with natural proportions and no sexualized styling.

## Global Negative Prompt

```text
child, minor, teenager, schoolchild, underage, sensual pose, erotic, nudity, readable text, fake letters, distorted words, subtitles, watermark, logo text, brand text, duplicated phone, broken smartphone, broken UI, messy interface, deformed hands, extra fingers, extra limbs, distorted face, uncanny skin, bad anatomy, low quality, jpeg artifacts, blurry subject, overexposed highlights, flat lighting, cluttered composition, cartoon, anime, illustration, plastic skin, AI artifacts
```

## Block Plan

| Block | Article placement | Visual job | Best-fit ratio |
| --- | --- | --- | --- |
| 1 | Very beginning / cover | Contrast old language learning vs AI future | 16:9 landscape |
| 2 | Before `Почему старые подходы дают сбой?` | Fear, confusion, barrier before live speech | 16:9 landscape |
| 3 | Before `Будущее, которое уже наступило` | Confidence, ease, AI language practice with Neriva mood | 16:9 landscape |

## 1. Cover - Main Visual Hook

### 1A - Neon Desk Contrast

Output filename: `neriva_dzen_01_cover_neon_desk_v01.png`

```text
Cinematic photorealistic editorial image for a Yandex Zen article cover, 16:9 landscape. A dark moody room with a wooden desk in the foreground. On the left side: a pile of old dusty English language textbooks, worn notebooks, paper flashcards, and a closed paper dictionary. On the right side: a sleek modern smartphone glowing with an abstract futuristic AI tutor interface, soft blue-white light reflecting on the desk. Cold blue neon light enters from a rain-streaked window, dramatic shadows, sharp focus on the contrast between obsolete study materials and modern AI learning. Premium tech mood, high detail, realistic dust, wood grain, glass reflections, cinematic 35mm lens, shallow depth of field. Leave clean negative space in the upper center for article cropping. No readable text, no logos, no watermark.
```

### 1B - Old Classroom, New Signal

Output filename: `neriva_dzen_01_cover_classroom_signal_v02.png`

```text
Photorealistic cinematic article cover, 16:9 landscape. An abandoned old language classroom at night: scratched wooden desks, chalk dust, faded grammar charts without readable words, stacked outdated textbooks, a cold rain window in the background. In the center foreground, a modern smartphone stands upright on the desk, projecting a subtle holographic AI conversation interface made of abstract chat shapes, waveform lines, and glowing translation particles. The room feels heavy and old, while the phone feels clean, precise, and alive. Teal-blue neon rim light, warm amber desk highlights, rich contrast, realistic texture, premium editorial photography, 8k detail. No readable text, no brand lettering, no watermark.
```

### 1C - Top-Down Editorial Split

Output filename: `neriva_dzen_01_cover_editorial_split_v03.png`

```text
High-end photorealistic top-down editorial composition for a Yandex Zen cover, 16:9 landscape. A dark wooden study table is visually split between two worlds: on one side, dusty grammar textbooks, a crumpled vocabulary list, red correction marks, and old paper flashcards; on the other side, a clean modern smartphone with a glowing abstract AI language coach interface, luminous blue neural lines, and soft glass reflections. Strong cinematic contrast between analog frustration and digital clarity. Controlled premium lighting, cold blue neon mixed with a small warm desk lamp, highly detailed paper texture and phone glass. Balanced composition with safe empty space for headline crop. No readable text, no fake letters, no logo, no watermark.
```

## 2. Old Approaches Fail

### 2A - Street With Dictionary

Output filename: `neriva_dzen_02_old_methods_street_dictionary_v01.png`

```text
Cinematic photorealistic scene for an article section about fear of speaking English, 16:9 landscape. An adult person, age 25+, stands alone on a dark city street at night, looking confused and frozen before speaking. Wet asphalt reflects cold neon signs as abstract colored light, heavy fog in the background. The person holds a crumpled paper dictionary and a few old study notes, symbolizing outdated learning methods. The city around them suggests real conversation pressure without showing readable signs or text. Moody gritty atmosphere, 35mm lens, shallow depth of field, teal and muted red reflections, realistic rain texture, emotional but not melodramatic. No readable text, no logos, no watermark.
```

### 2B - Platform Silence

Output filename: `neriva_dzen_02_old_methods_platform_silence_v02.png`

```text
Photorealistic cinematic image, 16:9 landscape. An adult person, age 25+, stands on a late-night metro platform, tense and unable to speak, holding an old pocket dictionary and paper grammar notes. Around the person, abstract blurred silhouettes pass by like fluent speakers, with soft translucent speech bubbles made only of shapes, no words. The platform lights are cold and clinical, the floor is wet and reflective, the atmosphere suggests isolation inside a busy city. Premium editorial look, 50mm lens, realistic skin and fabric, foggy depth, sharp subject, subtle motion blur in background. No readable text, no subtitles, no logo, no watermark.
```

### 2C - Red Marks Memory

Output filename: `neriva_dzen_02_old_methods_red_marks_v03.png`

```text
Cinematic photorealistic editorial scene, 16:9 landscape. A dim room with an adult learner, age 25+, seated at a desk, hands hovering over an open notebook filled with vague red correction marks and crossed-out lines with no readable words. In the background, a dark classroom silhouette and an old blackboard appear as a memory of school pressure. The learner looks hesitant, as if afraid to make a mistake before speaking. Lighting is dramatic: cold blue window light, a small warm lamp, deep shadows, detailed paper texture, realistic anxious posture, premium magazine photography. The image should show fear of mistakes without looking like horror. No readable text, no fake letters, no watermark.
```

## 3. The Future Has Already Arrived

### 3A - Urban Holographic Coach

Output filename: `neriva_dzen_03_future_holographic_coach_v01.png`

```text
Photorealistic cinematic image for an article section about Neriva and AI language learning, 16:9 landscape. A confident adult young professional, age 25+, walks through a modern urban night city after rain, relaxed posture, speaking naturally while interacting with a glowing holographic AI language coach interface in the air. The interface shows abstract chat cards, pronunciation waveform shapes, and progress arcs with no readable text. Wet asphalt reflects teal and orange lighting, city bokeh in the background, premium high-tech mood, clean futuristic realism, 35mm lens, dynamic but elegant composition. The person feels calm, capable, and supported by technology. No readable text, no logos, no watermark.
```

### 3B - Smartphone Conversation Flow

Output filename: `neriva_dzen_03_future_phone_conversation_v02.png`

```text
High-end photorealistic lifestyle-tech scene, 16:9 landscape. An adult person, age 25+, sits by a large city window in a modern evening cafe or coworking space, practicing spoken English through a smartphone AI tutor. The phone emits a soft clean glow, with abstract floating UI elements: waveform, gentle correction marks, conversation bubbles, and progress rings, all without readable text. The person looks relaxed and slightly smiling, as if speaking freely without fear. Warm interior light contrasts with cool blue city reflections outside, polished but natural, premium editorial photography, realistic glass, fabric, skin, and phone reflections. No visible brand text, no fake interface words, no watermark.
```

### 3C - Crossing Into Fluency

Output filename: `neriva_dzen_03_future_crossing_fluency_v03.png`

```text
Cinematic photorealistic editorial image, 16:9 landscape. A confident adult learner, age 25+, crosses a rain-wet city intersection at night while a subtle AI assistant layer flows around them as luminous language particles, abstract waveform paths, and clean holographic cards. The scene feels like moving from uncertainty into fluency: old paper notes are barely visible in the background, while the foreground is modern, clear, and energetic. Premium teal, white, and warm amber lighting, realistic reflections, crisp subject, deep city bokeh, elegant high-tech Neriva mood without literal brand text. Leave balanced space for article layout. No readable text, no logo lettering, no watermark.
```

## Suggested Selection Logic

- Use `1A` if the article needs the clearest old-school vs AI contrast.
- Use `1C` if the Yandex Zen cover will receive a large headline overlay.
- Use `2A` for the most emotional section divider.
- Use `2B` if you want the old-methods block to feel more urban and social.
- Use `3A` for the strongest "future already here" impression.
- Use `3B` if the article should feel more accessible and product-like instead of cyberpunk.
