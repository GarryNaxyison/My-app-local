import json
import re
import xml.etree.ElementTree as ET
from pathlib import Path

NS = {"tei": "http://www.tei-c.org/ns/1.0"}
ROOT = Path(__file__).resolve().parent.parent
VOCABULARY_DIR = ROOT / "data" / "vocabulary"
LETTER_RE = re.compile(
    r"[A-Za-zÀ-ÖØ-öø-ÿĀ-žА-Яа-яЁёІіЇїЄєҐґĄąĆćĘęŁłŃńÓóŚśŹźŻżÇçÃãÕõÁáÉéÍíÓóÚúÂâÊêÔôÀàÜüÄäÖöß]+"
)
SKIP_POS = {"suffix", "prefix", "interj", "abbr", "symbol", "letter", "num", "pn", "prep", "conj", "pron", "det"}


def clean_text(node):
    if node is None:
        return ""
    return re.sub(r"\s+", " ", "".join(node.itertext())).strip()


def is_single_word(text):
    return bool(LETTER_RE.fullmatch(text))


def cefr_level(rank):
    if rank < 1500:
        return "A1"
    if rank < 4500:
        return "A2"
    if rank < 9000:
        return "B1"
    if rank < 14000:
        return "B2"
    return "C1"


def parse(path, language, reverse=False):
    root = ET.parse(path).getroot()
    words = []
    seen = set()
    for entry in root.findall(".//tei:entry", NS):
        source = clean_text(entry.find("./tei:form/tei:orth", NS))
        pos = clean_text(entry.find("./tei:gramGrp/tei:pos", NS)).lower()
        if pos in SKIP_POS:
            continue
        for quote in entry.findall(".//tei:cit[@type='trans']/tei:quote", NS):
            target = clean_text(quote)
            learning, russian = (target, source) if reverse else (source, target)
            if not is_single_word(learning) or not is_single_word(russian):
                continue
            if learning[:1].isupper() or russian[:1].isupper():
                continue
            key = learning.lower()
            if key in seen:
                continue
            seen.add(key)
            rank = len(words) + 1
            words.append(
                {
                    "language": language,
                    "russian": russian,
                    "english": learning,
                    "level": cefr_level(rank),
                    "topic": "freedict",
                    "part_of_speech": pos or "word",
                    "source": "FreeDict/WikDict CC BY-SA 3.0",
                    "frequency_rank": rank,
                }
            )
            break
    return words


SOURCES = [
    ("freedict-pol-rus-2025.11.23.tar/pol-rus/pol-rus.tei", "pl", False),
    ("freedict-rus-por-2025.11.23.tar/rus-por/rus-por.tei", "pt", True),
]

for source_path, language, reverse in SOURCES:
    words = parse(source_path, language, reverse=reverse)
    VOCABULARY_DIR.mkdir(parents=True, exist_ok=True)
    (VOCABULARY_DIR / f"vocabulary_words_{language}.json").write_text(
        json.dumps(words, ensure_ascii=False, indent=2) + "\n",
        encoding="utf-8",
    )
    print(language, len(words))
