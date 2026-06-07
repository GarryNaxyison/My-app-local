import json
import re
import xml.etree.ElementTree as ET
from pathlib import Path

TEI_PATH = Path("freedict-ita-rus-src/ita-rus/ita-rus.tei")
FREQ_PATH = Path("it_50k.txt")
OUT_PATH = Path("data/vocabulary/vocabulary_words_it.json")

WORD_RE = re.compile(r"^[A-Za-zÀ-ÖØ-öø-ÿ]+$")
CYRILLIC_RE = re.compile(r"[А-Яа-яЁё]")
SKIP_POS = {"suffix", "prefix", "letter", "symbol", "phrase", "proverb"}

CURATED_CORE = [
    ("casa", "дом", "A1", "noun"),
    ("acqua", "вода", "A1", "noun"),
    ("pane", "хлеб", "A1", "noun"),
    ("cibo", "еда", "A1", "noun"),
    ("caffè", "кофе", "A1", "noun"),
    ("tè", "чай", "A1", "noun"),
    ("latte", "молоко", "A1", "noun"),
    ("mela", "яблоко", "A1", "noun"),
    ("giorno", "день", "A1", "noun"),
    ("notte", "ночь", "A1", "noun"),
    ("mattina", "утро", "A1", "noun"),
    ("sera", "вечер", "A1", "noun"),
    ("tempo", "время; погода", "A1", "noun"),
    ("anno", "год", "A1", "noun"),
    ("oggi", "сегодня", "A1", "adverb"),
    ("domani", "завтра", "A1", "adverb"),
    ("ieri", "вчера", "A1", "adverb"),
    ("persona", "человек", "A1", "noun"),
    ("uomo", "мужчина", "A1", "noun"),
    ("donna", "женщина", "A1", "noun"),
    ("bambino", "ребёнок", "A1", "noun"),
    ("amico", "друг", "A1", "noun"),
    ("famiglia", "семья", "A1", "noun"),
    ("madre", "мать", "A1", "noun"),
    ("padre", "отец", "A1", "noun"),
    ("scuola", "школа", "A1", "noun"),
    ("lavoro", "работа", "A1", "noun"),
    ("città", "город", "A1", "noun"),
    ("strada", "улица; дорога", "A1", "noun"),
    ("paese", "страна; деревня", "A1", "noun"),
    ("lingua", "язык", "A1", "noun"),
    ("parola", "слово", "A1", "noun"),
    ("libro", "книга", "A1", "noun"),
    ("telefono", "телефон", "A1", "noun"),
    ("macchina", "машина", "A1", "noun"),
    ("treno", "поезд", "A1", "noun"),
    ("soldi", "деньги", "A1", "noun"),
    ("nome", "имя", "A1", "noun"),
    ("numero", "номер; число", "A1", "noun"),
    ("grande", "большой", "A1", "adjective"),
    ("piccolo", "маленький", "A1", "adjective"),
    ("buono", "хороший; вкусный", "A1", "adjective"),
    ("cattivo", "плохой", "A1", "adjective"),
    ("nuovo", "новый", "A1", "adjective"),
    ("vecchio", "старый", "A1", "adjective"),
    ("bello", "красивый", "A1", "adjective"),
    ("facile", "лёгкий", "A1", "adjective"),
    ("difficile", "трудный", "A1", "adjective"),
    ("caldo", "тёплый; жаркий", "A1", "adjective"),
    ("freddo", "холодный", "A1", "adjective"),
    ("aperto", "открытый", "A1", "adjective"),
    ("chiuso", "закрытый", "A1", "adjective"),
    ("essere", "быть", "A1", "verb"),
    ("avere", "иметь", "A1", "verb"),
    ("fare", "делать", "A1", "verb"),
    ("andare", "идти; ехать", "A1", "verb"),
    ("venire", "приходить; приезжать", "A1", "verb"),
    ("parlare", "говорить", "A1", "verb"),
    ("capire", "понимать", "A1", "verb"),
    ("sapere", "знать", "A1", "verb"),
    ("vedere", "видеть", "A1", "verb"),
    ("volere", "хотеть", "A1", "verb"),
    ("potere", "мочь", "A1", "verb"),
    ("dovere", "долженствовать", "A1", "verb"),
    ("mangiare", "есть", "A1", "verb"),
    ("bere", "пить", "A1", "verb"),
    ("dormire", "спать", "A1", "verb"),
    ("comprare", "покупать", "A1", "verb"),
    ("pagare", "платить", "A1", "verb"),
    ("vivere", "жить", "A1", "verb"),
    ("studiare", "учиться; изучать", "A1", "verb"),
    ("scrivere", "писать", "A1", "verb"),
    ("leggere", "читать", "A1", "verb"),
    ("ascoltare", "слушать", "A1", "verb"),
    ("molto", "очень; много", "A1", "adverb"),
    ("poco", "мало; немного", "A1", "adverb"),
    ("qui", "здесь", "A1", "adverb"),
    ("là", "там", "A1", "adverb"),
    ("sempre", "всегда", "A2", "adverb"),
    ("mai", "никогда", "A2", "adverb"),
    ("presto", "скоро; рано", "A2", "adverb"),
    ("tardi", "поздно", "A2", "adverb"),
    ("perché", "почему; потому что", "A2", "adverb"),
    ("quando", "когда", "A2", "adverb"),
    ("dove", "где; куда", "A2", "adverb"),
    ("come", "как", "A2", "adverb"),
    ("quanto", "сколько; насколько", "A2", "adverb"),
    ("problema", "проблема", "A2", "noun"),
    ("domanda", "вопрос", "A2", "noun"),
    ("risposta", "ответ", "A2", "noun"),
    ("viaggio", "путешествие", "A2", "noun"),
    ("biglietto", "билет", "A2", "noun"),
    ("stazione", "станция; вокзал", "A2", "noun"),
    ("aeroporto", "аэропорт", "A2", "noun"),
    ("albergo", "отель", "A2", "noun"),
    ("negozio", "магазин", "A2", "noun"),
    ("ristorante", "ресторан", "A2", "noun"),
    ("medico", "врач", "A2", "noun"),
    ("aiuto", "помощь", "A2", "noun"),
    ("salute", "здоровье", "A2", "noun"),
    ("felice", "счастливый", "A2", "adjective"),
    ("stanco", "усталый", "A2", "adjective"),
    ("importante", "важный", "A2", "adjective"),
    ("possibile", "возможный", "A2", "adjective"),
    ("diverso", "разный; другой", "A2", "adjective"),
    ("chiedere", "спрашивать; просить", "A2", "verb"),
    ("rispondere", "отвечать", "A2", "verb"),
    ("aspettare", "ждать", "A2", "verb"),
    ("trovare", "находить", "A2", "verb"),
    ("perdere", "терять", "A2", "verb"),
    ("prendere", "брать", "A2", "verb"),
    ("lasciare", "оставлять", "A2", "verb"),
    ("portare", "нести; приносить", "A2", "verb"),
    ("aprire", "открывать", "A2", "verb"),
    ("chiudere", "закрывать", "A2", "verb"),
    ("iniziare", "начинать", "A2", "verb"),
    ("finire", "заканчивать", "A2", "verb"),
]


def local_name(tag):
    return tag.split("}")[-1]


def norm_text(value):
    return re.sub(r"\s+", " ", (value or "").strip())


def make_id(word):
    return "it:" + word.strip().lower().replace(" ", "-").replace("'", "").replace("/", "-").replace(".", "").replace(",", "")


def is_learnable_italian(word):
    word = norm_text(word)
    return 2 <= len(word) <= 32 and WORD_RE.match(word) is not None and not word.isupper()


def clean_translation(text):
    text = norm_text(text)
    text = re.sub(r"\s*\([^)]{0,60}\)", "", text).strip()
    text = re.sub(r"\s*\[[^]]{0,60}\]", "", text).strip()
    text = text.strip(" ;,")
    if not CYRILLIC_RE.search(text):
        return ""
    if len(text) > 90:
        return ""
    if re.search(r"[<>{}\[\]|=]", text):
        return ""
    return text


def level_for_rank(rank):
    if rank <= 1200:
        return "A2"
    if rank <= 3500:
        return "B1"
    if rank <= 7500:
        return "B2"
    if rank <= 12000:
        return "C1"
    return "C2"


def load_frequency_ranks():
    ranks = {}
    if not FREQ_PATH.exists():
        return ranks
    for index, line in enumerate(FREQ_PATH.read_text(encoding="utf-8", errors="ignore").splitlines(), start=1):
        word = line.split(" ", 1)[0].strip().lower()
        if is_learnable_italian(word) and word not in ranks:
            ranks[word] = index
    return ranks


def main():
    freq_rank = load_frequency_ranks()
    entries = []
    seen = set()

    for rank, (word, russian, level, pos) in enumerate(CURATED_CORE, start=1):
        key = word.lower()
        if key in seen:
            continue
        seen.add(key)
        entries.append({
            "english": word,
            "russian": russian,
            "level": level,
            "topic": "curated-core",
            "part_of_speech": pos,
            "source": "Curated learning core for Russian-speaking learners",
            "frequency_rank": rank,
            "id": make_id(word),
        })

    parsed = []
    order = 0
    for event, elem in ET.iterparse(TEI_PATH, events=("end",)):
        if local_name(elem.tag) != "entry":
            continue
        order += 1
        orth = ""
        pos = ""
        translations = []
        for node in elem.iter():
            name = local_name(node.tag)
            if name == "orth" and not orth:
                orth = norm_text("".join(node.itertext()))
            elif name == "pos" and not pos:
                pos = norm_text("".join(node.itertext())).lower()
            elif name == "quote":
                translation = clean_translation("".join(node.itertext()))
                if translation and translation not in translations:
                    translations.append(translation)
        elem.clear()

        word = orth.lower()
        if pos in SKIP_POS or not is_learnable_italian(word) or word in seen or not translations:
            continue
        russian = "; ".join(translations[:2])
        if len(russian) > 100:
            continue
        rank = freq_rank.get(word, 1_000_000 + order)
        parsed.append((rank, order, word, russian, pos))

    parsed.sort(key=lambda item: (item[0], item[1]))
    for rank, order, word, russian, pos in parsed:
        if word in seen:
            continue
        seen.add(word)
        known_rank = rank if rank < 1_000_000 else 0
        entry = {
            "english": word,
            "russian": russian,
            "level": level_for_rank(known_rank) if known_rank else "C2",
            "topic": "freedict-filtered",
            "part_of_speech": pos,
            "source": "FreeDict/WikDict ita-rus 2025.11.23 (CC BY-SA 3.0)",
            "id": make_id(word),
        }
        if known_rank:
            entry["frequency_rank"] = known_rank
        entries.append(entry)

    OUT_PATH.write_text(json.dumps(entries, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"wrote {len(entries)} entries")
    print(f"curated {len(CURATED_CORE)}; filtered freedict {len(entries) - len(CURATED_CORE)}")


if __name__ == "__main__":
    main()
