import json
import re
import time
import unicodedata
import xml.etree.ElementTree as ET
from pathlib import Path

from wordfreq import zipf_frequency


ROOT = Path(__file__).resolve().parent.parent
SOURCES = ROOT / "tools" / "sources"
VOCABULARY_DIR = ROOT / "data" / "vocabulary"

NS = {"tei": "http://www.tei-c.org/ns/1.0"}
TEI_ENTRY = "{http://www.tei-c.org/ns/1.0}entry"

LANGUAGE_ORDER = [
    "en", "ru", "es", "de", "fr", "it", "zh", "ja", "ko", "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "pl", "ro", "pt",
    "ar", "bn", "cs", "el", "hi", "hu", "id", "nl", "sv", "ta", "te", "th", "tl", "tr", "vi",
]
INTERFACE_LANGUAGES = ["ru"] + [code for code in LANGUAGE_ORDER if code != "ru"]

BASE_EN_RU_TEI = "eng-rus.tei"

WIKDICT_TEI = {
    "de": "deu-eng.tei",
    "es": "spa-eng.tei",
    "fr": "fra-eng.tei",
    "it": "ita-eng.tei",
    "zh": "zho-eng.tei",
    "ja": "jpn-eng.tei",
    "pl": "pol-eng.tei",
    "pt": "por-eng.tei",
}

WIKDICT_EN_TARGET = {
    "de": "eng-deu.tei",
    "es": "eng-spa.tei",
    "fr": "eng-fra.tei",
    "it": "eng-ita.tei",
    "zh": "eng-zho.tei",
    "ja": "eng-jpn.tei",
    "pl": "eng-pol.tei",
    "pt": "eng-por.tei",
}

KAIKKI_JSONL = {
    "ru": "kaikki-ru.jsonl",
    "ko": "kaikki-ko.jsonl",
    "tg": "kaikki-tg.jsonl",
    "uz": "kaikki-uz.jsonl",
    "tt": "kaikki-tt.jsonl",
    "hy": "kaikki-hy.jsonl",
    "kk": "kaikki-kk.jsonl",
    "ky": "kaikki-ky.jsonl",
    "ka": "kaikki-ka.jsonl",
    "uk": "kaikki-uk.jsonl",
    "ro": "kaikki-ro.jsonl",
    "ar": "kaikki-ar.jsonl",
    "bn": "kaikki-bn.jsonl",
    "cs": "kaikki-cs.jsonl",
    "el": "kaikki-el.jsonl",
    "hi": "kaikki-hi.jsonl",
    "hu": "kaikki-hu.jsonl",
    "id": "kaikki-id.jsonl",
    "nl": "kaikki-nl.jsonl",
    "sv": "kaikki-sv.jsonl",
    "ta": "kaikki-ta.jsonl",
    "te": "kaikki-te.jsonl",
    "th": "kaikki-th.jsonl",
    "tl": "kaikki-tl.jsonl",
    "tr": "kaikki-tr.jsonl",
    "vi": "kaikki-vi.jsonl",
}

KAIKKI_LIMITS = {
    "ru": 52000,
}

WIKDICT_ONLY = ["de", "es", "fr", "it", "zh", "ja", "pl", "pt"]
KAIKKI_ONLY = ["ko", "tg", "uz", "tt", "hy", "kk", "ky", "ka", "uk", "ro", "ar", "bn", "cs", "el", "hi", "hu", "id", "nl", "sv", "ta", "te", "th", "tl", "tr", "vi"]

COMMON_CORE = [
    ("house", "noun", "A1"),
    ("water", "noun", "A1"),
    ("bread", "noun", "A1"),
    ("food", "noun", "A1"),
    ("coffee", "noun", "A1"),
    ("tea", "noun", "A1"),
    ("milk", "noun", "A1"),
    ("apple", "noun", "A1"),
    ("day", "noun", "A1"),
    ("night", "noun", "A1"),
    ("morning", "noun", "A1"),
    ("evening", "noun", "A1"),
    ("time", "noun", "A1"),
    ("year", "noun", "A1"),
    ("today", "adverb", "A1"),
    ("tomorrow", "adverb", "A1"),
    ("yesterday", "adverb", "A1"),
    ("person", "noun", "A1"),
    ("man", "noun", "A1"),
    ("woman", "noun", "A1"),
    ("child", "noun", "A1"),
    ("friend", "noun", "A1"),
    ("family", "noun", "A1"),
    ("mother", "noun", "A1"),
    ("father", "noun", "A1"),
    ("school", "noun", "A1"),
    ("work", "noun", "A1"),
    ("city", "noun", "A1"),
    ("street", "noun", "A1"),
    ("country", "noun", "A1"),
    ("language", "noun", "A1"),
    ("word", "noun", "A1"),
    ("book", "noun", "A1"),
    ("phone", "noun", "A1"),
    ("car", "noun", "A1"),
    ("train", "noun", "A1"),
    ("money", "noun", "A1"),
    ("name", "noun", "A1"),
    ("big", "adjective", "A1"),
    ("small", "adjective", "A1"),
    ("good", "adjective", "A1"),
    ("bad", "adjective", "A1"),
    ("new", "adjective", "A1"),
    ("old", "adjective", "A1"),
    ("beautiful", "adjective", "A1"),
    ("easy", "adjective", "A1"),
    ("difficult", "adjective", "A1"),
    ("warm", "adjective", "A1"),
    ("cold", "adjective", "A1"),
    ("open", "adjective", "A1"),
    ("closed", "adjective", "A1"),
    ("be", "verb", "A1"),
    ("have", "verb", "A1"),
    ("do", "verb", "A1"),
    ("go", "verb", "A1"),
    ("come", "verb", "A1"),
    ("speak", "verb", "A1"),
    ("understand", "verb", "A1"),
    ("know", "verb", "A1"),
    ("see", "verb", "A1"),
    ("want", "verb", "A1"),
    ("can", "verb", "A1"),
    ("must", "verb", "A1"),
    ("eat", "verb", "A1"),
    ("drink", "verb", "A1"),
    ("sleep", "verb", "A1"),
    ("buy", "verb", "A1"),
    ("pay", "verb", "A1"),
    ("live", "verb", "A1"),
    ("study", "verb", "A1"),
    ("write", "verb", "A1"),
    ("read", "verb", "A1"),
    ("listen", "verb", "A1"),
]

COMMON_CORE_KEYS = {item[0] for item in COMMON_CORE}

COMMON_TRANSLATION_OVERRIDES = {
    "house": {
        "de": "Haus",
        "es": "casa",
        "fr": "maison",
        "it": "casa",
        "tg": "хона",
        "uz": "uy",
        "tt": "өй",
        "hy": "տուն",
        "kk": "үй",
        "ky": "үй",
        "ka": "სახლი",
        "uk": "дім",
        "zh": "房子",
        "ja": "家",
        "ko": "집",
        "pl": "dom",
        "ro": "casă",
        "pt": "casa",
    },
    "water": {
        "de": "Wasser",
        "es": "agua",
        "fr": "eau",
        "it": "acqua",
        "tg": "об",
        "uz": "suv",
        "tt": "су",
        "hy": "ջուր",
        "kk": "су",
        "ky": "суу",
        "ka": "წყალი",
        "uk": "вода",
        "pl": "woda",
        "ro": "apă",
        "pt": "água",
    },
    "do": {
        "de": "machen",
        "es": "hacer",
        "fr": "faire",
        "it": "fare",
        "tg": "кардан",
        "uz": "qilmoq",
        "tt": "эшләү",
        "hy": "անել",
        "kk": "жасау",
        "ky": "кылуу",
        "ka": "კეთება",
        "uk": "робити",
        "pl": "robić",
        "ro": "face",
        "pt": "fazer",
    },
    "make": {
        "de": "machen",
        "es": "hacer",
        "fr": "faire",
        "it": "fare",
        "tg": "кардан",
        "uz": "qilmoq",
        "tt": "эшләү",
        "hy": "անել",
        "kk": "жасау",
        "ky": "кылуу",
        "ka": "კეთება",
        "uk": "робити",
        "pl": "robić",
        "ro": "face",
        "pt": "fazer",
    },
}

COMMON_TRANSLATION_OVERRIDES.update(
    {
        "bread": {"de": "Brot", "es": "pan", "fr": "pain", "it": "pane", "pl": "chleb", "pt": "pão", "ro": "pâine", "uk": "хліб", "zh": "面包", "ja": "パン", "ko": "빵"},
        "food": {"de": "Essen", "es": "comida", "fr": "nourriture", "it": "cibo", "pl": "jedzenie", "pt": "comida", "ro": "mâncare", "uk": "їжа", "zh": "食物", "ja": "食べ物", "ko": "음식"},
        "coffee": {"de": "Kaffee", "es": "café", "fr": "café", "it": "caffè", "pl": "kawa", "pt": "café", "ro": "cafea", "uk": "кава", "zh": "咖啡", "ja": "コーヒー", "ko": "커피"},
        "tea": {"de": "Tee", "es": "té", "fr": "thé", "it": "tè", "pl": "herbata", "pt": "chá", "ro": "ceai", "uk": "чай", "zh": "茶", "ja": "お茶", "ko": "차"},
        "milk": {"de": "Milch", "es": "leche", "fr": "lait", "it": "latte", "pl": "mleko", "pt": "leite", "ro": "lapte", "uk": "молоко", "zh": "牛奶", "ja": "牛乳", "ko": "우유"},
        "apple": {"de": "Apfel", "es": "manzana", "fr": "pomme", "it": "mela", "pl": "jabłko", "pt": "maçã", "ro": "măr", "uk": "яблуко", "zh": "苹果", "ja": "りんご", "ko": "사과"},
        "day": {"de": "Tag", "es": "día", "fr": "jour", "it": "giorno", "pl": "dzień", "pt": "dia", "ro": "zi", "uk": "день", "zh": "天", "ja": "日", "ko": "날"},
        "night": {"de": "Nacht", "es": "noche", "fr": "nuit", "it": "notte", "pl": "noc", "pt": "noite", "ro": "noapte", "uk": "ніч", "zh": "夜晚", "ja": "夜", "ko": "밤"},
        "morning": {"de": "Morgen", "es": "mañana", "fr": "matin", "it": "mattina", "pl": "rano", "pt": "manhã", "ro": "dimineață", "uk": "ранок", "zh": "早上", "ja": "朝", "ko": "아침"},
        "evening": {"de": "Abend", "es": "tarde", "fr": "soir", "it": "sera", "pl": "wieczór", "pt": "noite", "ro": "seară", "uk": "вечір", "zh": "晚上", "ja": "夕方", "ko": "저녁"},
        "time": {"de": "Zeit", "es": "tiempo", "fr": "temps", "it": "tempo", "pl": "czas", "pt": "tempo", "ro": "timp", "uk": "час", "zh": "时间", "ja": "時間", "ko": "시간"},
        "year": {"de": "Jahr", "es": "año", "fr": "année", "it": "anno", "pl": "rok", "pt": "ano", "ro": "an", "uk": "рік", "zh": "年", "ja": "年", "ko": "해"},
        "today": {"de": "heute", "es": "hoy", "fr": "aujourd'hui", "it": "oggi", "pl": "dzisiaj", "pt": "hoje", "ro": "astăzi", "uk": "сьогодні", "zh": "今天", "ja": "今日", "ko": "오늘"},
        "tomorrow": {"de": "morgen", "es": "mañana", "fr": "demain", "it": "domani", "pl": "jutro", "pt": "amanhã", "ro": "mâine", "uk": "завтра", "zh": "明天", "ja": "明日", "ko": "내일"},
        "yesterday": {"de": "gestern", "es": "ayer", "fr": "hier", "it": "ieri", "pl": "wczoraj", "pt": "ontem", "ro": "ieri", "uk": "вчора", "zh": "昨天", "ja": "昨日", "ko": "어제"},
        "person": {"de": "Person", "es": "persona", "fr": "personne", "it": "persona", "pl": "osoba", "pt": "pessoa", "ro": "persoană", "uk": "людина", "zh": "人", "ja": "人", "ko": "사람"},
        "man": {"de": "Mann", "es": "hombre", "fr": "homme", "it": "uomo", "pl": "mężczyzna", "pt": "homem", "ro": "bărbat", "uk": "чоловік", "zh": "男人", "ja": "男性", "ko": "남자"},
        "woman": {"de": "Frau", "es": "mujer", "fr": "femme", "it": "donna", "pl": "kobieta", "pt": "mulher", "ro": "femeie", "uk": "жінка", "zh": "女人", "ja": "女性", "ko": "여자"},
        "child": {"de": "Kind", "es": "niño", "fr": "enfant", "it": "bambino", "pl": "dziecko", "pt": "criança", "ro": "copil", "uk": "дитина", "zh": "孩子", "ja": "子供", "ko": "아이"},
        "friend": {"de": "Freund", "es": "amigo", "fr": "ami", "it": "amico", "pl": "przyjaciel", "pt": "amigo", "ro": "prieten", "uk": "друг", "zh": "朋友", "ja": "友達", "ko": "친구"},
        "family": {"de": "Familie", "es": "familia", "fr": "famille", "it": "famiglia", "pl": "rodzina", "pt": "família", "ro": "familie", "uk": "родина", "zh": "家庭", "ja": "家族", "ko": "가족"},
        "mother": {"de": "Mutter", "es": "madre", "fr": "mère", "it": "madre", "pl": "matka", "pt": "mãe", "ro": "mamă", "uk": "мати", "zh": "母亲", "ja": "母", "ko": "어머니"},
        "father": {"de": "Vater", "es": "padre", "fr": "père", "it": "padre", "pl": "ojciec", "pt": "pai", "ro": "tată", "uk": "батько", "zh": "父亲", "ja": "父", "ko": "아버지"},
        "school": {"de": "Schule", "es": "escuela", "fr": "école", "it": "scuola", "pl": "szkoła", "pt": "escola", "ro": "școală", "uk": "школа", "zh": "学校", "ja": "学校", "ko": "학교"},
        "work": {"de": "Arbeit", "es": "trabajo", "fr": "travail", "it": "lavoro", "pl": "praca", "pt": "trabalho", "ro": "muncă", "uk": "робота", "zh": "工作", "ja": "仕事", "ko": "일"},
        "city": {"de": "Stadt", "es": "ciudad", "fr": "ville", "it": "città", "pl": "miasto", "pt": "cidade", "ro": "oraș", "uk": "місто", "zh": "城市", "ja": "町", "ko": "도시"},
        "street": {"de": "Straße", "es": "calle", "fr": "rue", "it": "strada", "pl": "ulica", "pt": "rua", "ro": "stradă", "uk": "вулиця", "zh": "街道", "ja": "通り", "ko": "거리"},
        "country": {"de": "Land", "es": "país", "fr": "pays", "it": "paese", "pl": "kraj", "pt": "país", "ro": "țară", "uk": "країна", "zh": "国家", "ja": "国", "ko": "나라"},
        "language": {"de": "Sprache", "es": "idioma", "fr": "langue", "it": "lingua", "pl": "język", "pt": "idioma", "ro": "limbă", "uk": "мова", "zh": "语言", "ja": "言語", "ko": "언어"},
        "word": {"de": "Wort", "es": "palabra", "fr": "mot", "it": "parola", "pl": "słowo", "pt": "palavra", "ro": "cuvânt", "uk": "слово", "zh": "词", "ja": "言葉", "ko": "단어"},
        "book": {"de": "Buch", "es": "libro", "fr": "livre", "it": "libro", "pl": "książka", "pt": "livro", "ro": "carte", "uk": "книга", "zh": "书", "ja": "本", "ko": "책"},
        "phone": {"de": "Telefon", "es": "teléfono", "fr": "téléphone", "it": "telefono", "pl": "telefon", "pt": "telefone", "ro": "telefon", "uk": "телефон", "zh": "电话", "ja": "電話", "ko": "전화"},
        "car": {"de": "Auto", "es": "coche", "fr": "voiture", "it": "auto", "pl": "samochód", "pt": "carro", "ro": "mașină", "uk": "машина", "zh": "汽车", "ja": "車", "ko": "차"},
        "train": {"de": "Zug", "es": "tren", "fr": "train", "it": "treno", "pl": "pociąg", "pt": "trem", "ro": "tren", "uk": "поїзд", "zh": "火车", "ja": "電車", "ko": "기차"},
        "money": {"de": "Geld", "es": "dinero", "fr": "argent", "it": "denaro", "pl": "pieniądze", "pt": "dinheiro", "ro": "bani", "uk": "гроші", "zh": "钱", "ja": "お金", "ko": "돈"},
        "name": {"de": "Name", "es": "nombre", "fr": "nom", "it": "nome", "pl": "imię", "pt": "nome", "ro": "nume", "uk": "ім'я", "zh": "名字", "ja": "名前", "ko": "이름"},
        "big": {"de": "groß", "es": "grande", "fr": "grand", "it": "grande", "pl": "duży", "pt": "grande", "ro": "mare", "uk": "великий", "zh": "大", "ja": "大きい", "ko": "크다"},
        "small": {"de": "klein", "es": "pequeño", "fr": "petit", "it": "piccolo", "pl": "mały", "pt": "pequeno", "ro": "mic", "uk": "малий", "zh": "小", "ja": "小さい", "ko": "작다"},
        "good": {"de": "gut", "es": "bueno", "fr": "bon", "it": "buono", "pl": "dobry", "pt": "bom", "ro": "bun", "uk": "добрий", "zh": "好", "ja": "良い", "ko": "좋다"},
        "bad": {"de": "schlecht", "es": "malo", "fr": "mauvais", "it": "cattivo", "pl": "zły", "pt": "ruim", "ro": "rău", "uk": "поганий", "zh": "坏", "ja": "悪い", "ko": "나쁘다"},
        "new": {"de": "neu", "es": "nuevo", "fr": "nouveau", "it": "nuovo", "pl": "nowy", "pt": "novo", "ro": "nou", "uk": "новий", "zh": "新", "ja": "新しい", "ko": "새롭다"},
        "old": {"de": "alt", "es": "viejo", "fr": "vieux", "it": "vecchio", "pl": "stary", "pt": "velho", "ro": "vechi", "uk": "старий", "zh": "旧", "ja": "古い", "ko": "오래되다"},
        "beautiful": {"de": "schön", "es": "bonito", "fr": "beau", "it": "bello", "pl": "piękny", "pt": "bonito", "ro": "frumos", "uk": "красивий", "zh": "美丽", "ja": "美しい", "ko": "아름답다"},
        "easy": {"de": "einfach", "es": "fácil", "fr": "facile", "it": "facile", "pl": "łatwy", "pt": "fácil", "ro": "ușor", "uk": "легкий", "zh": "容易", "ja": "簡単", "ko": "쉽다"},
        "difficult": {"de": "schwierig", "es": "difícil", "fr": "difficile", "it": "difficile", "pl": "trudny", "pt": "difícil", "ro": "dificil", "uk": "важкий", "zh": "困难", "ja": "難しい", "ko": "어렵다"},
        "warm": {"de": "warm", "es": "cálido", "fr": "chaud", "it": "caldo", "pl": "ciepły", "pt": "quente", "ro": "cald", "uk": "теплий", "zh": "暖和", "ja": "暖かい", "ko": "따뜻하다"},
        "cold": {"de": "kalt", "es": "frío", "fr": "froid", "it": "freddo", "pl": "zimny", "pt": "frio", "ro": "rece", "uk": "холодний", "zh": "冷", "ja": "寒い", "ko": "춥다"},
        "open": {"de": "offen", "es": "abierto", "fr": "ouvert", "it": "aperto", "pl": "otwarty", "pt": "aberto", "ro": "deschis", "uk": "відкритий", "zh": "开", "ja": "開いている", "ko": "열려 있다"},
        "closed": {"de": "geschlossen", "es": "cerrado", "fr": "fermé", "it": "chiuso", "pl": "zamknięty", "pt": "fechado", "ro": "închis", "uk": "закритий", "zh": "关", "ja": "閉まっている", "ko": "닫혀 있다"},
        "be": {"de": "sein", "es": "ser; estar", "fr": "être", "it": "essere; stare", "pl": "być", "pt": "ser; estar", "ro": "fi", "uk": "бути", "zh": "是", "ja": "ある", "ko": "이다"},
        "have": {"de": "haben", "es": "tener", "fr": "avoir", "it": "avere", "pl": "mieć", "pt": "ter", "ro": "avea", "uk": "мати", "zh": "有", "ja": "持つ", "ko": "가지다"},
        "do": {"de": "machen; tun", "es": "hacer", "fr": "faire", "it": "fare", "pl": "robić", "pt": "fazer", "ro": "face", "uk": "робити", "zh": "做", "ja": "する", "ko": "하다"},
        "make": {"de": "machen", "es": "hacer", "fr": "faire", "it": "fare", "pl": "robić", "pt": "fazer", "ro": "face", "uk": "робити", "zh": "做", "ja": "作る", "ko": "만들다"},
        "go": {"de": "gehen", "es": "ir", "fr": "aller", "it": "andare", "pl": "iść", "pt": "ir", "ro": "merge", "uk": "іти", "zh": "去", "ja": "行く", "ko": "가다"},
        "come": {"de": "kommen", "es": "venir", "fr": "venir", "it": "venire", "pl": "przychodzić", "pt": "vir", "ro": "veni", "uk": "приходити", "zh": "来", "ja": "来る", "ko": "오다"},
        "speak": {"de": "sprechen", "es": "hablar", "fr": "parler", "it": "parlare", "pl": "mówić", "pt": "falar", "ro": "vorbi", "uk": "говорити", "zh": "说", "ja": "話す", "ko": "말하다"},
        "understand": {"de": "verstehen", "es": "entender", "fr": "comprendre", "it": "capire", "pl": "rozumieć", "pt": "entender", "ro": "înțelege", "uk": "розуміти", "zh": "懂", "ja": "分かる", "ko": "이해하다"},
        "know": {"de": "wissen", "es": "saber", "fr": "savoir", "it": "sapere", "pl": "wiedzieć", "pt": "saber", "ro": "ști", "uk": "знати", "zh": "知道", "ja": "知る", "ko": "알다"},
        "see": {"de": "sehen", "es": "ver", "fr": "voir", "it": "vedere", "pl": "widzieć", "pt": "ver", "ro": "vedea", "uk": "бачити", "zh": "看见", "ja": "見る", "ko": "보다"},
        "want": {"de": "wollen", "es": "querer", "fr": "vouloir", "it": "volere", "pl": "chcieć", "pt": "querer", "ro": "vrea", "uk": "хотіти", "zh": "想要", "ja": "欲しい", "ko": "원하다"},
        "can": {"de": "können", "es": "poder", "fr": "pouvoir", "it": "potere", "pl": "móc", "pt": "poder", "ro": "putea", "uk": "могти", "zh": "能", "ja": "できる", "ko": "할 수 있다"},
        "must": {"de": "müssen", "es": "deber", "fr": "devoir", "it": "dovere", "pl": "musieć", "pt": "dever", "ro": "trebui", "uk": "мусити", "zh": "必须", "ja": "しなければならない", "ko": "해야 하다"},
        "eat": {"de": "essen", "es": "comer", "fr": "manger", "it": "mangiare", "pl": "jeść", "pt": "comer", "ro": "mânca", "uk": "їсти", "zh": "吃", "ja": "食べる", "ko": "먹다"},
        "drink": {"de": "trinken", "es": "beber", "fr": "boire", "it": "bere", "pl": "pić", "pt": "beber", "ro": "bea", "uk": "пити", "zh": "喝", "ja": "飲む", "ko": "마시다"},
        "sleep": {"de": "schlafen", "es": "dormir", "fr": "dormir", "it": "dormire", "pl": "spać", "pt": "dormir", "ro": "dormi", "uk": "спати", "zh": "睡觉", "ja": "眠る", "ko": "자다"},
        "buy": {"de": "kaufen", "es": "comprar", "fr": "acheter", "it": "comprare", "pl": "kupić", "pt": "comprar", "ro": "cumpăra", "uk": "купувати", "zh": "买", "ja": "買う", "ko": "사다"},
        "pay": {"de": "bezahlen", "es": "pagar", "fr": "payer", "it": "pagare", "pl": "płacić", "pt": "pagar", "ro": "plăti", "uk": "платити", "zh": "支付", "ja": "払う", "ko": "내다"},
        "live": {"de": "leben", "es": "vivir", "fr": "vivre", "it": "vivere", "pl": "żyć", "pt": "viver", "ro": "trăi", "uk": "жити", "zh": "住", "ja": "住む", "ko": "살다"},
        "study": {"de": "lernen", "es": "estudiar", "fr": "étudier", "it": "studiare", "pl": "uczyć się", "pt": "estudar", "ro": "studia", "uk": "вчитися", "zh": "学习", "ja": "勉強する", "ko": "공부하다"},
        "write": {"de": "schreiben", "es": "escribir", "fr": "écrire", "it": "scrivere", "pl": "pisać", "pt": "escrever", "ro": "scrie", "uk": "писати", "zh": "写", "ja": "書く", "ko": "쓰다"},
        "read": {"de": "lesen", "es": "leer", "fr": "lire", "it": "leggere", "pl": "czytać", "pt": "ler", "ro": "citi", "uk": "читати", "zh": "读", "ja": "読む", "ko": "읽다"},
        "listen": {"de": "hören", "es": "escuchar", "fr": "écouter", "it": "ascoltare", "pl": "słuchać", "pt": "ouvir", "ro": "asculta", "uk": "слухати", "zh": "听", "ja": "聞く", "ko": "듣다"},
        "bread": {"zh": "面包", "ja": "パン", "ko": "빵"},
        "food": {"zh": "食物", "ja": "食べ物", "ko": "음식"},
        "coffee": {"zh": "咖啡", "ja": "コーヒー", "ko": "커피"},
        "tea": {"zh": "茶", "ja": "お茶", "ko": "차"},
        "milk": {"zh": "牛奶", "ja": "牛乳", "ko": "우유"},
        "apple": {"zh": "苹果", "ja": "りんご", "ko": "사과"},
        "day": {"zh": "天", "ja": "日", "ko": "날"},
        "night": {"zh": "夜晚", "ja": "夜", "ko": "밤"},
        "morning": {"zh": "早上", "ja": "朝", "ko": "아침"},
        "evening": {"zh": "晚上", "ja": "夕方", "ko": "저녁"},
        "time": {"zh": "时间", "ja": "時間", "ko": "시간"},
        "year": {"zh": "年", "ja": "年", "ko": "해"},
        "today": {"zh": "今天", "ja": "今日", "ko": "오늘"},
        "tomorrow": {"zh": "明天", "ja": "明日", "ko": "내일"},
        "yesterday": {"zh": "昨天", "ja": "昨日", "ko": "어제"},
        "person": {"zh": "人", "ja": "人", "ko": "사람"},
        "man": {"zh": "男人", "ja": "男性", "ko": "남자"},
        "woman": {"zh": "女人", "ja": "女性", "ko": "여자"},
        "child": {"zh": "孩子", "ja": "子供", "ko": "아이"},
        "friend": {"zh": "朋友", "ja": "友達", "ko": "친구"},
        "family": {"zh": "家庭", "ja": "家族", "ko": "가족"},
        "mother": {"zh": "母亲", "ja": "母", "ko": "어머니"},
        "father": {"zh": "父亲", "ja": "父", "ko": "아버지"},
        "school": {"zh": "学校", "ja": "学校", "ko": "학교"},
        "work": {"zh": "工作", "ja": "仕事", "ko": "일"},
        "city": {"zh": "城市", "ja": "町", "ko": "도시"},
        "street": {"zh": "街道", "ja": "通り", "ko": "거리"},
        "country": {"zh": "国家", "ja": "国", "ko": "나라"},
        "language": {"zh": "语言", "ja": "言語", "ko": "언어"},
        "word": {"zh": "词", "ja": "言葉", "ko": "단어"},
        "book": {"zh": "书", "ja": "本", "ko": "책"},
        "phone": {"zh": "电话", "ja": "電話", "ko": "전화"},
        "car": {"zh": "汽车", "ja": "車", "ko": "차"},
        "train": {"zh": "火车", "ja": "電車", "ko": "기차"},
        "money": {"zh": "钱", "ja": "お金", "ko": "돈"},
        "name": {"zh": "名字", "ja": "名前", "ko": "이름"},
        "big": {"zh": "大", "ja": "大きい", "ko": "크다"},
        "small": {"zh": "小", "ja": "小さい", "ko": "작다"},
        "good": {"zh": "好", "ja": "良い", "ko": "좋다"},
        "bad": {"zh": "坏", "ja": "悪い", "ko": "나쁘다"},
        "new": {"zh": "新", "ja": "新しい", "ko": "새롭다"},
        "old": {"zh": "旧", "ja": "古い", "ko": "오래되다"},
        "beautiful": {"zh": "美丽", "ja": "美しい", "ko": "아름답다"},
        "easy": {"zh": "容易", "ja": "簡単", "ko": "쉽다"},
        "difficult": {"zh": "困难", "ja": "難しい", "ko": "어렵다"},
        "warm": {"zh": "暖和", "ja": "暖かい", "ko": "따뜻하다"},
        "cold": {"zh": "冷", "ja": "寒い", "ko": "춥다"},
        "open": {"zh": "开", "ja": "開いている", "ko": "열려 있다"},
        "closed": {"zh": "关", "ja": "閉まっている", "ko": "닫혀 있다"},
        "be": {"zh": "是", "ja": "ある", "ko": "이다"},
        "have": {"zh": "有", "ja": "持つ", "ko": "가지다"},
        "do": {"zh": "做", "ja": "する", "ko": "하다"},
        "make": {"zh": "做", "ja": "作る", "ko": "만들다"},
        "go": {"zh": "去", "ja": "行く", "ko": "가다"},
        "come": {"zh": "来", "ja": "来る", "ko": "오다"},
        "speak": {"zh": "说", "ja": "話す", "ko": "말하다"},
        "understand": {"zh": "懂", "ja": "分かる", "ko": "이해하다"},
        "know": {"zh": "知道", "ja": "知る", "ko": "알다"},
        "see": {"zh": "看见", "ja": "見る", "ko": "보다"},
        "want": {"zh": "想要", "ja": "欲しい", "ko": "원하다"},
        "can": {"zh": "能", "ja": "できる", "ko": "할 수 있다"},
        "must": {"zh": "必须", "ja": "しなければならない", "ko": "해야 한다"},
        "eat": {"zh": "吃", "ja": "食べる", "ko": "먹다"},
        "drink": {"zh": "喝", "ja": "飲む", "ko": "마시다"},
        "sleep": {"zh": "睡觉", "ja": "眠る", "ko": "자다"},
        "buy": {"zh": "买", "ja": "買う", "ko": "사다"},
        "pay": {"zh": "支付", "ja": "払う", "ko": "내다"},
        "live": {"zh": "住", "ja": "住む", "ko": "살다"},
        "study": {"zh": "学习", "ja": "勉強する", "ko": "공부하다"},
        "write": {"zh": "写", "ja": "書く", "ko": "쓰다"},
        "read": {"zh": "读", "ja": "読む", "ko": "읽다"},
        "listen": {"zh": "听", "ja": "聞く", "ko": "듣다"},
    }
)

COMMON_CORE_TRANSLATION_OVERRIDES = {
    "bread": {"de": "Brot", "es": "pan", "fr": "pain", "it": "pane", "pl": "chleb", "pt": "pão", "ro": "pâine", "uk": "хліб"},
    "food": {"de": "Essen", "es": "comida", "fr": "nourriture", "it": "cibo", "pl": "jedzenie", "pt": "comida", "ro": "mâncare", "uk": "їжа"},
    "coffee": {"de": "Kaffee", "es": "café", "fr": "café", "it": "caffè", "pl": "kawa", "pt": "café", "ro": "cafea", "uk": "кава"},
    "tea": {"de": "Tee", "es": "té", "fr": "thé", "it": "tè", "pl": "herbata", "pt": "chá", "ro": "ceai", "uk": "чай"},
    "milk": {"de": "Milch", "es": "leche", "fr": "lait", "it": "latte", "pl": "mleko", "pt": "leite", "ro": "lapte", "uk": "молоко"},
    "apple": {"de": "Apfel", "es": "manzana", "fr": "pomme", "it": "mela", "pl": "jabłko", "pt": "maçã", "ro": "măr", "uk": "яблуко"},
    "day": {"de": "Tag", "es": "día", "fr": "jour", "it": "giorno", "pl": "dzień", "pt": "dia", "ro": "zi", "uk": "день"},
    "night": {"de": "Nacht", "es": "noche", "fr": "nuit", "it": "notte", "pl": "noc", "pt": "noite", "ro": "noapte", "uk": "ніч"},
    "morning": {"de": "Morgen", "es": "mañana", "fr": "matin", "it": "mattina", "pl": "rano", "pt": "manhã", "ro": "dimineață", "uk": "ранок"},
    "evening": {"de": "Abend", "es": "tarde", "fr": "soir", "it": "sera", "pl": "wieczór", "pt": "noite", "ro": "seară", "uk": "вечір"},
    "time": {"de": "Zeit", "es": "tiempo", "fr": "temps", "it": "tempo", "pl": "czas", "pt": "tempo", "ro": "timp", "uk": "час"},
    "year": {"de": "Jahr", "es": "año", "fr": "année", "it": "anno", "pl": "rok", "pt": "ano", "ro": "an", "uk": "рік"},
    "today": {"de": "heute", "es": "hoy", "fr": "aujourd'hui", "it": "oggi", "pl": "dzisiaj", "pt": "hoje", "ro": "astăzi", "uk": "сьогодні"},
    "tomorrow": {"de": "morgen", "es": "mañana", "fr": "demain", "it": "domani", "pl": "jutro", "pt": "amanhã", "ro": "mâine", "uk": "завтра"},
    "yesterday": {"de": "gestern", "es": "ayer", "fr": "hier", "it": "ieri", "pl": "wczoraj", "pt": "ontem", "ro": "ieri", "uk": "вчора"},
    "person": {"de": "Person", "es": "persona", "fr": "personne", "it": "persona", "pl": "osoba", "pt": "pessoa", "ro": "persoană", "uk": "людина"},
    "man": {"de": "Mann", "es": "hombre", "fr": "homme", "it": "uomo", "pl": "mężczyzna", "pt": "homem", "ro": "bărbat", "uk": "чоловік"},
    "woman": {"de": "Frau", "es": "mujer", "fr": "femme", "it": "donna", "pl": "kobieta", "pt": "mulher", "ro": "femeie", "uk": "жінка"},
    "child": {"de": "Kind", "es": "niño", "fr": "enfant", "it": "bambino", "pl": "dziecko", "pt": "criança", "ro": "copil", "uk": "дитина"},
    "friend": {"de": "Freund", "es": "amigo", "fr": "ami", "it": "amico", "pl": "przyjaciel", "pt": "amigo", "ro": "prieten", "uk": "друг"},
    "family": {"de": "Familie", "es": "familia", "fr": "famille", "it": "famiglia", "pl": "rodzina", "pt": "família", "ro": "familie", "uk": "родина"},
    "mother": {"de": "Mutter", "es": "madre", "fr": "mère", "it": "madre", "pl": "matka", "pt": "mãe", "ro": "mamă", "uk": "мати"},
    "father": {"de": "Vater", "es": "padre", "fr": "père", "it": "padre", "pl": "ojciec", "pt": "pai", "ro": "tată", "uk": "батько"},
    "school": {"de": "Schule", "es": "escuela", "fr": "école", "it": "scuola", "pl": "szkoła", "pt": "escola", "ro": "școală", "uk": "школа"},
    "work": {"de": "Arbeit", "es": "trabajo", "fr": "travail", "it": "lavoro", "pl": "praca", "pt": "trabalho", "ro": "muncă", "uk": "робота"},
    "city": {"de": "Stadt", "es": "ciudad", "fr": "ville", "it": "città", "pl": "miasto", "pt": "cidade", "ro": "oraș", "uk": "місто"},
    "street": {"de": "Straße", "es": "calle", "fr": "rue", "it": "strada", "pl": "ulica", "pt": "rua", "ro": "stradă", "uk": "вулиця"},
    "country": {"de": "Land", "es": "país", "fr": "pays", "it": "paese", "pl": "kraj", "pt": "país", "ro": "țară", "uk": "країна"},
    "language": {"de": "Sprache", "es": "idioma", "fr": "langue", "it": "lingua", "pl": "język", "pt": "idioma", "ro": "limbă", "uk": "мова"},
    "word": {"de": "Wort", "es": "palabra", "fr": "mot", "it": "parola", "pl": "słowo", "pt": "palavra", "ro": "cuvânt", "uk": "слово"},
    "book": {"de": "Buch", "es": "libro", "fr": "livre", "it": "libro", "pl": "książka", "pt": "livro", "ro": "carte", "uk": "книга"},
    "phone": {"de": "Telefon", "es": "teléfono", "fr": "téléphone", "it": "telefono", "pl": "telefon", "pt": "telefone", "ro": "telefon", "uk": "телефон"},
    "car": {"de": "Auto", "es": "coche", "fr": "voiture", "it": "auto", "pl": "samochód", "pt": "carro", "ro": "mașină", "uk": "машина"},
    "train": {"de": "Zug", "es": "tren", "fr": "train", "it": "treno", "pl": "pociąg", "pt": "trem", "ro": "tren", "uk": "поїзд"},
    "money": {"de": "Geld", "es": "dinero", "fr": "argent", "it": "denaro", "pl": "pieniądze", "pt": "dinheiro", "ro": "bani", "uk": "гроші"},
    "name": {"de": "Name", "es": "nombre", "fr": "nom", "it": "nome", "pl": "imię", "pt": "nome", "ro": "nume", "uk": "ім'я"},
    "big": {"de": "groß", "es": "grande", "fr": "grand", "it": "grande", "pl": "duży", "pt": "grande", "ro": "mare", "uk": "великий"},
    "small": {"de": "klein", "es": "pequeño", "fr": "petit", "it": "piccolo", "pl": "mały", "pt": "pequeno", "ro": "mic", "uk": "малий"},
    "good": {"de": "gut", "es": "bueno", "fr": "bon", "it": "buono", "pl": "dobry", "pt": "bom", "ro": "bun", "uk": "добрий"},
    "bad": {"de": "schlecht", "es": "malo", "fr": "mauvais", "it": "cattivo", "pl": "zły", "pt": "ruim", "ro": "rău", "uk": "поганий"},
    "new": {"de": "neu", "es": "nuevo", "fr": "nouveau", "it": "nuovo", "pl": "nowy", "pt": "novo", "ro": "nou", "uk": "новий"},
    "old": {"de": "alt", "es": "viejo", "fr": "vieux", "it": "vecchio", "pl": "stary", "pt": "velho", "ro": "vechi", "uk": "старий"},
    "beautiful": {"de": "schön", "es": "bonito", "fr": "beau", "it": "bello", "pl": "piękny", "pt": "bonito", "ro": "frumos", "uk": "красивий"},
    "easy": {"de": "einfach", "es": "fácil", "fr": "facile", "it": "facile", "pl": "łatwy", "pt": "fácil", "ro": "ușor", "uk": "легкий"},
    "difficult": {"de": "schwierig", "es": "difícil", "fr": "difficile", "it": "difficile", "pl": "trudny", "pt": "difícil", "ro": "dificil", "uk": "важкий"},
    "warm": {"de": "warm", "es": "cálido", "fr": "chaud", "it": "caldo", "pl": "ciepły", "pt": "quente", "ro": "cald", "uk": "теплий"},
    "cold": {"de": "kalt", "es": "frío", "fr": "froid", "it": "freddo", "pl": "zimny", "pt": "frio", "ro": "rece", "uk": "холодний"},
    "open": {"de": "offen", "es": "abierto", "fr": "ouvert", "it": "aperto", "pl": "otwarty", "pt": "aberto", "ro": "deschis", "uk": "відкритий"},
    "closed": {"de": "geschlossen", "es": "cerrado", "fr": "fermé", "it": "chiuso", "pl": "zamknięty", "pt": "fechado", "ro": "închis", "uk": "закритий"},
    "be": {"de": "sein", "es": "ser; estar", "fr": "être", "it": "essere; stare", "tg": "будан", "uz": "bo'lmoq", "tt": "булу", "hy": "լինել", "kk": "болу", "ky": "болуу", "ka": "ყოფნა", "pl": "być", "pt": "ser; estar", "ro": "fi", "uk": "бути"},
    "have": {"de": "haben", "es": "tener", "fr": "avoir", "it": "avere", "pl": "mieć", "pt": "ter", "ro": "avea", "uk": "мати"},
    "do": {"de": "machen; tun", "es": "hacer", "fr": "faire", "it": "fare", "pl": "robić", "pt": "fazer", "ro": "face", "uk": "робити"},
    "make": {"de": "machen", "es": "hacer", "fr": "faire", "it": "fare", "pl": "robić", "pt": "fazer", "ro": "face", "uk": "робити"},
    "go": {"de": "gehen", "es": "ir", "fr": "aller", "it": "andare", "pl": "iść", "pt": "ir", "ro": "merge", "uk": "іти"},
    "come": {"de": "kommen", "es": "venir", "fr": "venir", "it": "venire", "pl": "przychodzić", "pt": "vir", "ro": "veni", "uk": "приходити"},
    "speak": {"de": "sprechen", "es": "hablar", "fr": "parler", "it": "parlare", "pl": "mówić", "pt": "falar", "ro": "vorbi", "uk": "говорити"},
    "understand": {"de": "verstehen", "es": "entender", "fr": "comprendre", "it": "capire", "pl": "rozumieć", "pt": "entender", "ro": "înțelege", "uk": "розуміти"},
    "know": {"de": "wissen", "es": "saber", "fr": "savoir", "it": "sapere", "pl": "wiedzieć", "pt": "saber", "ro": "ști", "uk": "знати"},
    "see": {"de": "sehen", "es": "ver", "fr": "voir", "it": "vedere", "pl": "widzieć", "pt": "ver", "ro": "vedea", "uk": "бачити"},
    "want": {"de": "wollen", "es": "querer", "fr": "vouloir", "it": "volere", "pl": "chcieć", "pt": "querer", "ro": "vrea", "uk": "хотіти"},
    "can": {"de": "können", "es": "poder", "fr": "pouvoir", "it": "potere", "pl": "móc", "pt": "poder", "ro": "putea", "uk": "могти"},
    "must": {"de": "müssen", "es": "deber", "fr": "devoir", "it": "dovere", "pl": "musieć", "pt": "dever", "ro": "trebui", "uk": "мусити"},
    "eat": {"de": "essen", "es": "comer", "fr": "manger", "it": "mangiare", "pl": "jeść", "pt": "comer", "ro": "mânca", "uk": "їсти"},
    "drink": {"de": "trinken", "es": "beber", "fr": "boire", "it": "bere", "pl": "pić", "pt": "beber", "ro": "bea", "uk": "пити"},
    "sleep": {"de": "schlafen", "es": "dormir", "fr": "dormir", "it": "dormire", "pl": "spać", "pt": "dormir", "ro": "dormi", "uk": "спати"},
    "buy": {"de": "kaufen", "es": "comprar", "fr": "acheter", "it": "comprare", "pl": "kupić", "pt": "comprar", "ro": "cumpăra", "uk": "купувати"},
    "pay": {"de": "bezahlen", "es": "pagar", "fr": "payer", "it": "pagare", "pl": "płacić", "pt": "pagar", "ro": "plăti", "uk": "платити"},
    "live": {"de": "leben", "es": "vivir", "fr": "vivre", "it": "vivere", "pl": "żyć", "pt": "viver", "ro": "trăi", "uk": "жити"},
    "study": {"de": "lernen", "es": "estudiar", "fr": "étudier", "it": "studiare", "pl": "studiować", "pt": "estudar", "ro": "studia", "uk": "вчитися"},
    "write": {"de": "schreiben", "es": "escribir", "fr": "écrire", "it": "scrivere", "pl": "pisać", "pt": "escrever", "ro": "scrie", "uk": "писати"},
    "read": {"de": "lesen", "es": "leer", "fr": "lire", "it": "leggere", "pl": "czytać", "pt": "ler", "ro": "citi", "uk": "читати"},
    "listen": {"de": "hören", "es": "escuchar", "fr": "écouter", "it": "ascoltare", "pl": "słuchać", "pt": "ouvir", "ro": "asculta", "uk": "слухати"},
}

for key, values in COMMON_CORE_TRANSLATION_OVERRIDES.items():
    COMMON_TRANSLATION_OVERRIDES.setdefault(key, {}).update(values)

COMMON_RUSSIAN_OVERRIDES = {
    "house": "дом",
    "water": "вода",
    "bread": "хлеб",
    "food": "еда",
    "coffee": "кофе",
    "tea": "чай",
    "milk": "молоко",
    "apple": "яблоко",
    "day": "день",
    "night": "ночь",
    "morning": "утро",
    "evening": "вечер",
    "time": "время",
    "year": "год",
    "today": "сегодня",
    "tomorrow": "завтра",
    "yesterday": "вчера",
    "person": "человек",
    "man": "мужчина",
    "woman": "женщина",
    "child": "ребенок",
    "friend": "друг",
    "family": "семья",
    "mother": "мать",
    "father": "отец",
    "school": "школа",
    "work": "работа",
    "city": "город",
    "street": "улица",
    "country": "страна",
    "language": "язык",
    "word": "слово",
    "book": "книга",
    "phone": "телефон",
    "car": "машина",
    "train": "поезд",
    "money": "деньги",
    "name": "имя",
    "big": "большой",
    "small": "маленький",
    "good": "хороший",
    "bad": "плохой",
    "new": "новый",
    "old": "старый",
    "beautiful": "красивый",
    "easy": "легкий",
    "difficult": "трудный",
    "warm": "теплый",
    "cold": "холодный",
    "open": "открытый",
    "closed": "закрытый",
    "be": "быть",
    "have": "иметь",
    "do": "делать",
    "go": "идти",
    "come": "приходить",
    "speak": "говорить",
    "understand": "понимать",
    "know": "знать",
    "see": "видеть",
    "want": "хотеть",
    "can": "мочь",
    "must": "должен",
    "eat": "есть",
    "drink": "пить",
    "sleep": "спать",
    "buy": "покупать",
    "pay": "платить",
    "live": "жить",
    "study": "учиться",
    "write": "писать",
    "read": "читать",
    "listen": "слушать",
}

SKIP_POS = {
    "abbr",
    "abbrev",
    "article",
    "character",
    "circumfix",
    "classifier",
    "ideophone",
    "interfix",
    "letter",
    "name",
    "particle",
    "phrase",
    "prefix",
    "prep",
    "preposition",
    "proper",
    "proper noun",
    "proverb",
    "punct",
    "punctuation",
    "romanization",
    "suffix",
    "symbol",
}

BAD_GLOSS_RE = re.compile(
    r"\b("
    r"abbreviation|acronym|alternative|archaic|borrowed|dated|dialectal|ellipsis|"
    r"feminine|form of|inflection|letter|masculine|misspelling|neuter|obsolete|"
    r"plural|prefix|proper noun|romanization|suffix|superlative|comparative"
    r")\b",
    re.IGNORECASE,
)
WORD_RE = re.compile(r"^[^\W\d_][^\d_]{0,39}$", re.UNICODE)
BAD_WORD_RE = re.compile(r"[\s/@#*_+=\\|<>\[\]{}()0-9.,;:!?\"“”„]")
SCRIPT_RE = {
    "ar": re.compile(r"[\u0600-\u06ff]"),
    "bn": re.compile(r"[\u0980-\u09ff]"),
    "el": re.compile(r"[\u0370-\u03ff]"),
    "hi": re.compile(r"[\u0900-\u097f]"),
    "ta": re.compile(r"[\u0b80-\u0bff]"),
    "te": re.compile(r"[\u0c00-\u0c7f]"),
    "th": re.compile(r"[\u0e00-\u0e7f]"),
}
PRESERVE_MARK_LANGUAGES = {"ar", "bn", "el", "hi", "ta", "te", "th", "vi"}
BAD_ENGLISH_GLOSSES = {
    "a",
    "an",
    "the",
    "to",
    "of",
    "in",
    "on",
    "at",
    "by",
    "for",
    "from",
    "with",
    "without",
    "and",
    "or",
    "but",
    "that",
    "this",
    "these",
    "those",
    "it",
    "he",
    "she",
    "we",
    "they",
    "you",
    "i",
    "is",
    "are",
    "was",
    "were",
    "as",
    "not",
    "my",
    "your",
    "his",
    "her",
    "their",
    "our",
    "me",
    "him",
    "them",
    "if",
    "so",
    "all",
    "just",
    "like",
    "will",
}


def strip_combining_marks(value: str) -> str:
    return "".join(ch for ch in value if unicodedata.category(ch) != "Mn")


def normalize_headword(code: str, value: str) -> str:
    value = clean_space(value)
    if code not in PRESERVE_MARK_LANGUAGES:
        value = strip_combining_marks(value)
    return value


def clean_space(value: str) -> str:
    return re.sub(r"\s+", " ", value or "").strip()


def clean_text(node) -> str:
    if node is None:
        return ""
    return clean_space("".join(node.itertext()))


UZBEK_CYRILLIC_MAP = {
    "а": "a",
    "б": "b",
    "в": "v",
    "г": "g",
    "д": "d",
    "е": "e",
    "ё": "yo",
    "ж": "j",
    "з": "z",
    "и": "i",
    "й": "y",
    "к": "k",
    "л": "l",
    "м": "m",
    "н": "n",
    "о": "o",
    "п": "p",
    "р": "r",
    "с": "s",
    "т": "t",
    "у": "u",
    "ф": "f",
    "х": "x",
    "ц": "s",
    "ч": "ch",
    "ш": "sh",
    "ъ": "",
    "ь": "",
    "э": "e",
    "ю": "yu",
    "я": "ya",
    "ў": "o'",
    "қ": "q",
    "ғ": "g'",
    "ҳ": "h",
}


def transliterate_uzbek_cyrillic(value: str) -> str:
    result: list[str] = []
    for char in value:
        mapped = UZBEK_CYRILLIC_MAP.get(char.lower())
        if mapped is None:
            result.append(char)
            continue
        if char.isupper() and mapped:
            mapped = mapped[:1].upper() + mapped[1:]
        result.append(mapped)
    return clean_space("".join(result))


def normalize_script_value(code: str, value: str) -> str:
    value = clean_space(value)
    if code == "uz":
        return transliterate_uzbek_cyrillic(value)
    return value


def english_key(value: str) -> str:
    value = clean_space(strip_combining_marks(value).lower())
    value = re.sub(r"\([^)]*\)", " ", value)
    value = re.sub(r"^(to|a|an|the)\s+", "", value)
    value = value.replace("’", "'")
    value = re.sub(r"[^a-z0-9 '+-]", " ", value)
    value = clean_space(value)
    return value


def target_key(value: str) -> str:
    value = clean_space(strip_combining_marks(value).lower())
    value = value.replace("’", "'")
    value = re.sub(r"\s+", " ", value)
    return value


def is_learning_word(value: str, code: str) -> bool:
    value = normalize_headword(code, value)
    if not value or len(value) < 2 or len(value) > 40:
        return False
    script = SCRIPT_RE.get(code)
    if script and not script.search(value):
        return False
    if code == "en" and english_key(value) in BAD_ENGLISH_GLOSSES:
        return False
    if BAD_WORD_RE.search(value):
        return False
    if not WORD_RE.match(value):
        return False
    if code != "de" and value[:1].isupper():
        return False
    return True


def clean_gloss(value: str) -> str:
    value = clean_space(value)
    if not value:
        return ""
    value = re.sub(r"\([^)]*\)", " ", value)
    value = re.split(r";|/|\bor\b|\band\b|,", value, maxsplit=1, flags=re.IGNORECASE)[0]
    value = clean_space(value)
    value = re.sub(r"^(to|a|an|the)\s+", "", value, flags=re.IGNORECASE)
    value = clean_space(value)
    key = english_key(value)
    if not key or key in BAD_ENGLISH_GLOSSES or BAD_GLOSS_RE.search(value):
        return ""
    if len(key) < 2 or len(key) > 60:
        return ""
    if len(key.split()) > 3:
        return ""
    if not re.search(r"[a-z]", key):
        return ""
    return key


def clean_context(value: str) -> str:
    value = clean_space(value)
    if not value:
        return ""
    value = re.sub(r"\[\[([^\]|]+)\|([^\]]+)\]\]", r"\2", value)
    value = re.sub(r"\[\[([^\]]+)\]\]", r"\1", value)
    value = re.sub(r"\s+([,;:.!?])", r"\1", value)
    value = clean_space(value)
    if len(value) > 420:
        value = value[:417].rstrip(" ,;/") + "..."
    return value


def has_context(value: str) -> bool:
    value = clean_context(value)
    return bool(value and (";" in value or "/" in value or "(" in value or ")" in value or len(value.split()) > 3))


def contextual_entry_value(raw: dict, primary: str = "") -> str:
    candidates = [raw.get("context", ""), raw.get("russian", "")]
    translations = raw.get("translations", {})
    if isinstance(translations, dict):
        candidates.extend([translations.get("ru", ""), translations.get("en", "")])
    for candidate in candidates:
        context = clean_context(str(candidate))
        if not context:
            continue
        if context != clean_context(primary) and (candidate == raw.get("context", "") or has_context(context)):
            return context
    return ""


def cefr_level(rank: int) -> str:
    if rank < 1500:
        return "A1"
    if rank < 4500:
        return "A2"
    if rank < 9000:
        return "B1"
    if rank < 14000:
        return "B2"
    if rank < 22000:
        return "C1"
    return "C2"


def cefr_level_from_frequency(score: float) -> str:
    if score >= 5.15:
        return "A1"
    if score >= 4.65:
        return "A2"
    if score >= 4.05:
        return "B1"
    if score >= 3.45:
        return "B2"
    if score >= 2.85:
        return "C1"
    return "C2"


def frequency_score(code: str, target: str, gloss: str = "") -> float:
    candidates = []
    target = clean_space(target)
    gloss = clean_space(gloss)
    if target:
        candidates.append((target, code))
    if gloss:
        candidates.append((gloss, "en"))
        key = english_key(gloss)
        if key and key != gloss:
            candidates.append((key, "en"))

    score = 0.0
    for value, language in candidates:
        try:
            score = max(score, float(zipf_frequency(value, language)))
        except Exception:
            continue
    return score


def level_for_entry(code: str, target: str, gloss: str, raw: dict, fallback_rank: int) -> str:
    if raw.get("topic") == "multilingual-core":
        return raw.get("level") or "A1"
    score = frequency_score(code, target, gloss)
    if score > 0:
        return cefr_level_from_frequency(score)
    return cefr_level(fallback_rank)


def sorted_by_learning_level(code: str, entries: list[dict]) -> list[dict]:
    level_rank = {level: index for index, level in enumerate(["A1", "A2", "B1", "B2", "C1", "C2"])}

    def sort_key(item: dict) -> tuple[int, float, int, str]:
        level = str(item.get("level", "C2"))
        score = frequency_score(code, str(item.get("english", "")), str(item.get("translations", {}).get("en", "")))
        topic_rank = 0 if item.get("topic") == "multilingual-core" else 1
        rank = int(item.get("frequency_rank") or 999999)
        return (level_rank.get(level, 5), topic_rank, -score, rank, target_key(str(item.get("english", ""))))

    result = sorted(entries, key=sort_key)
    for index, item in enumerate(result, start=1):
        translations = item.get("translations", {})
        if isinstance(translations, dict):
            key = english_key(str(translations.get("en", "") or item.get("english", "")))
            if key in COMMON_RUSSIAN_OVERRIDES:
                item["russian"] = COMMON_RUSSIAN_OVERRIDES[key]
                translations["ru"] = COMMON_RUSSIAN_OVERRIDES[key]
        item["frequency_rank"] = index
    return result


def pos_class(pos: str) -> str:
    pos = clean_space((pos or "").lower())
    if pos.startswith("v") or "verb" in pos:
        return "verb"
    if pos.startswith("n") or "noun" in pos:
        return "noun"
    if pos.startswith("adj") or "adjective" in pos:
        return "adjective"
    if pos.startswith("adv") or "adverb" in pos:
        return "adverb"
    return "word"


def make_vocab_id(code: str, word: str) -> str:
    value = target_key(word)
    value = value.replace(" ", "-").replace("/", "-").replace(".", "").replace(",", "")
    value = value.replace("'", "").replace("’", "")
    if code == "en":
        return value
    return f"{code}:{value}"


def load_existing(code: str) -> list[dict]:
    path = VOCABULARY_DIR / ("vocabulary_words.json" if code == "en" else f"vocabulary_words_{code}.json")
    if not path.exists():
        return []
    return json.loads(path.read_text(encoding="utf-8"))


def add_en_map(
    en_to_lang: dict[str, dict[str, str]],
    en_to_lang_by_pos: dict[str, dict[str, dict[str, str]]],
    code: str,
    gloss: str,
    target: str,
    pos: str,
) -> None:
    key = english_key(gloss)
    if not key or not is_learning_word(target, code):
        return
    en_to_lang.setdefault(code, {})
    en_to_lang[code].setdefault(key, target)
    en_to_lang_by_pos.setdefault(code, {}).setdefault(key, {})
    en_to_lang_by_pos[code][key].setdefault(pos_class(pos), target)


def parse_tei(code: str, path: Path) -> tuple[list[dict], dict[str, str]]:
    entries: list[dict] = []
    target_to_en: dict[str, str] = {}
    seen: set[str] = set()

    context = ET.iterparse(path, events=("end",))
    for _, entry in context:
        if entry.tag != TEI_ENTRY:
            continue
        source = clean_space(strip_combining_marks(clean_text(entry.find("./tei:form/tei:orth", NS))))
        pos = clean_space(clean_text(entry.find("./tei:gramGrp/tei:pos", NS)).lower()) or "word"
        if pos in SKIP_POS or not is_learning_word(source, code):
            entry.clear()
            continue

        gloss = ""
        full_context = ""
        quote_texts: list[str] = []
        for quote in entry.findall(".//tei:cit[@type='trans']/tei:quote", NS):
            text = clean_text(quote)
            if text:
                quote_texts.append(text)
            if not gloss:
                gloss = clean_gloss(text)
        if not gloss:
            entry.clear()
            continue
        full_context = clean_context("; ".join(dict.fromkeys(quote_texts)))

        key = target_key(source)
        if key not in seen:
            seen.add(key)
            rank = len(entries) + 1
            item = {
                "language": code,
                "russian": gloss,
                "english": source,
                "level": cefr_level(rank),
                "topic": "wikdict",
                "part_of_speech": pos,
                "source": "WikDict/FreeDict TEI CC BY-SA 3.0",
                "frequency_rank": rank,
                "translations": {"en": gloss},
            }
            if full_context and full_context != gloss:
                item["context"] = full_context
            entries.append(item)
        target_to_en.setdefault(key, gloss)
        entry.clear()

    return entries, target_to_en


def parse_tei_en_to_target(code: str, path: Path) -> tuple[dict[str, str], dict[str, dict[str, str]]]:
    result: dict[str, str] = {}
    by_pos: dict[str, dict[str, str]] = {}

    context = ET.iterparse(path, events=("end",))
    for _, entry in context:
        if entry.tag != TEI_ENTRY:
            continue
        source = clean_gloss(clean_text(entry.find("./tei:form/tei:orth", NS)))
        pos = clean_space(clean_text(entry.find("./tei:gramGrp/tei:pos", NS)).lower()) or "word"
        if pos in SKIP_POS:
            entry.clear()
            continue
        key = english_key(source)
        if not key:
            entry.clear()
            continue

        values: list[str] = []
        for quote in entry.findall(".//tei:cit[@type='trans']/tei:quote", NS):
            target = clean_space(strip_combining_marks(clean_text(quote)))
            if is_learning_word(target, code):
                append_unique(values, target, code)
                if len(values) >= 4:
                    break
        if values:
            value = "; ".join(values)
            result.setdefault(key, value)
            by_pos.setdefault(key, {}).setdefault(pos_class(pos), value)
        entry.clear()

    return result, by_pos


def parse_tei_english_russian(path: Path) -> tuple[list[dict], dict[str, str]]:
    entries: list[dict] = []
    en_to_ru: dict[str, str] = {}
    seen: set[str] = set()

    context = ET.iterparse(path, events=("end",))
    for _, entry in context:
        if entry.tag != TEI_ENTRY:
            continue
        source = clean_space(strip_combining_marks(clean_text(entry.find("./tei:form/tei:orth", NS))))
        pos = clean_space(clean_text(entry.find("./tei:gramGrp/tei:pos", NS)).lower()) or "word"
        if pos in SKIP_POS or not is_learning_word(source, "en"):
            entry.clear()
            continue

        quote_texts: list[str] = []
        primary_ru = ""
        for quote in entry.findall(".//tei:cit[@type='trans']/tei:quote", NS):
            text = clean_context(clean_text(quote))
            if not text:
                continue
            quote_texts.append(text)
            if not primary_ru:
                primary_ru = clean_space(re.split(r";|/|,", text, maxsplit=1)[0])
        full_context = clean_context("; ".join(dict.fromkeys(quote_texts)))
        if not primary_ru and full_context:
            primary_ru = full_context
        if not primary_ru:
            entry.clear()
            continue

        key = english_key(source)
        target_id = make_vocab_id("en", source)
        if key and full_context:
            en_to_ru.setdefault(key, full_context)
        if target_id not in seen:
            seen.add(target_id)
            rank = len(entries) + 1
            item = {
                "id": target_id,
                "language": "en",
                "russian": full_context or primary_ru,
                "english": source,
                "level": cefr_level(rank),
                "topic": "wikdict",
                "part_of_speech": pos,
                "source": "WikDict/FreeDict eng-rus TEI CC BY-SA 3.0",
                "frequency_rank": rank,
                "translations": {"en": source, "ru": full_context or primary_ru},
            }
            if full_context and full_context != primary_ru:
                item["context"] = full_context
            entries.append(item)
        entry.clear()

    return entries, en_to_ru


def sense_gloss(sense: dict) -> str:
    if sense.get("form_of") or sense.get("alt_of"):
        return ""
    tags = {str(tag).lower() for tag in sense.get("tags", [])}
    if tags.intersection({"abbreviation", "alt-of", "form-of", "misspelling", "obsolete", "plural", "romanization"}):
        return ""
    for link in sense.get("links", []):
        if link and isinstance(link[0], str):
            gloss = clean_gloss(link[0])
            if gloss:
                return gloss
    for gloss in sense.get("glosses", []) or sense.get("raw_glosses", []):
        gloss = clean_gloss(gloss)
        if gloss:
            return gloss
    return ""


def sense_context(sense: dict) -> str:
    if sense.get("form_of") or sense.get("alt_of"):
        return ""
    tags = {str(tag).lower() for tag in sense.get("tags", [])}
    if tags.intersection({"abbreviation", "alt-of", "form-of", "misspelling", "obsolete", "plural", "romanization"}):
        return ""
    values: list[str] = []
    for link in sense.get("links", []):
        if link and isinstance(link[0], str):
            values.append(link[0])
    for gloss in sense.get("glosses", []) or sense.get("raw_glosses", []):
        if isinstance(gloss, str):
            values.append(gloss)
    return clean_context("; ".join(values))


def parse_kaikki(code: str, path: Path, limit: int = 30000) -> tuple[list[dict], dict[str, str]]:
    entries: list[dict] = []
    target_to_en: dict[str, str] = {}
    seen: set[str] = set()

    with path.open("r", encoding="utf-8") as handle:
        for line in handle:
            if len(entries) >= limit:
                break
            try:
                item = json.loads(line)
            except json.JSONDecodeError:
                continue
            if item.get("lang_code") != code:
                continue
            pos = clean_space(str(item.get("pos", "")).lower()) or "word"
            if pos in SKIP_POS:
                continue
            word = normalize_headword(code, str(item.get("word", "")))
            if not is_learning_word(word, code):
                continue
            key = target_key(word)
            if key in seen:
                continue

            gloss = ""
            context_values: list[str] = []
            for sense in item.get("senses", []):
                context = sense_context(sense)
                if context:
                    context_values.append(context)
                if not gloss:
                    gloss = sense_gloss(sense)
            if not gloss:
                continue
            full_context = clean_context("; ".join(dict.fromkeys(context_values)))

            seen.add(key)
            rank = len(entries) + 1
            item = {
                "language": code,
                "russian": gloss,
                "english": word,
                "level": cefr_level(rank),
                "topic": "wiktionary",
                "part_of_speech": pos,
                "source": "Kaikki/Wiktextract CC BY-SA 4.0",
                "frequency_rank": rank,
                "translations": {"en": gloss},
            }
            if full_context and full_context != gloss:
                item["context"] = full_context
            entries.append(item)
            target_to_en.setdefault(key, gloss)

    return entries, target_to_en


def russian_key(value: str) -> str:
    value = clean_space(value).lower()
    value = re.split(r";|,|/", value, maxsplit=1)[0]
    return clean_space(value)


def is_curated_item(item: dict) -> bool:
    return item.get("topic") == "curated-core" or str(item.get("source", "")).startswith("Curated")


def apply_existing_ru_pivot_maps(
    existing: dict[str, list[dict]],
    en_to_ru: dict[str, str],
    en_to_lang: dict[str, dict[str, str]],
    en_to_lang_by_pos: dict[str, dict[str, dict[str, str]]],
) -> None:
    ru_to_en: dict[str, str] = {}
    for key, russian in en_to_ru.items():
        ru_key = russian_key(russian)
        if ru_key:
            ru_to_en.setdefault(ru_key, key)

    for code in EXISTING_BASE:
        if code == "en":
            continue
        for item in existing.get(code, []):
            if not is_curated_item(item):
                continue
            target = normalize_script_value(code, clean_space(str(item.get("english", ""))))
            ru_key = russian_key(str(item.get("russian", "")))
            key = ru_to_en.get(ru_key)
            if not key or not is_learning_word(target, code):
                continue
            en_to_lang.setdefault(code, {})[key] = target
            en_to_lang_by_pos.setdefault(code, {}).setdefault(key, {})[pos_class(item.get("part_of_speech", ""))] = target


def build_ru_to_lang(existing: dict[str, list[dict]]) -> dict[str, dict[str, str]]:
    result: dict[str, dict[str, str]] = {}
    for code in EXISTING_BASE:
        if code == "en":
            continue
        for item in existing.get(code, []):
            if not is_curated_item(item):
                continue
            target = normalize_script_value(code, clean_space(str(item.get("english", ""))))
            ru_key = russian_key(str(item.get("russian", "")))
            if ru_key and is_learning_word(target, code):
                result.setdefault(code, {}).setdefault(ru_key, target)
    return result


def first_dictionary_value(value: str) -> str:
    original = clean_space(value)
    if not original:
        return ""
    value = re.sub(r"\[\[([^\]|]+)\|([^\]]+)\]\]", r"\2", original)
    value = re.sub(r"\[\[([^\]]+)\]\]", r"\1", value)
    value = re.sub(r"\([^)]*\)", " ", value)
    value = clean_space(re.split(r";|/|,", clean_space(value), maxsplit=1)[0])
    if value:
        return value
    return clean_space(re.sub(r"[()]", " ", re.split(r";|/|,", original, maxsplit=1)[0]))


def target_for_core(
    code: str,
    key: str,
    pos: str,
    ru: str,
    en_to_lang: dict[str, dict[str, str]],
    en_to_lang_by_pos: dict[str, dict[str, dict[str, str]]],
    ru_to_lang: dict[str, dict[str, str]],
) -> str:
    if code == "ru":
        value = en_to_lang_by_pos.get(code, {}).get(key, {}).get(pos_class(pos), "")
        if value:
            return value
        value = en_to_lang.get(code, {}).get(key, "")
        if value:
            return value
        return clean_russian_target(ru)
    value = COMMON_TRANSLATION_OVERRIDES.get(key, {}).get(code, "")
    if value:
        return first_dictionary_value(normalize_script_value(code, value))
    ru_key = russian_key(ru)
    if ru_key:
        value = ru_to_lang.get(code, {}).get(ru_key, "")
        if value:
            return first_dictionary_value(normalize_script_value(code, value))
    value = en_to_lang_by_pos.get(code, {}).get(key, {}).get(pos_class(pos), "")
    if value:
        return first_dictionary_value(normalize_script_value(code, value))
    return ""


def prepend_common_core(
    raw_by_code: dict[str, list[dict]],
    en_to_ru: dict[str, str],
    en_to_lang: dict[str, dict[str, str]],
    en_to_lang_by_pos: dict[str, dict[str, dict[str, str]]],
    ru_to_lang: dict[str, dict[str, str]],
) -> None:
    for code in LANGUAGE_ORDER:
        if code == "en":
            continue
        core: list[dict] = []
        seen: set[str] = set()
        for index, (gloss, pos, level) in enumerate(COMMON_CORE, start=1):
            key = english_key(gloss)
            ru = COMMON_RUSSIAN_OVERRIDES.get(key) or en_to_ru.get(key, gloss)
            target = target_for_core(code, key, pos, ru, en_to_lang, en_to_lang_by_pos, ru_to_lang)
            if not target or not is_learning_word(target, code):
                continue
            vocab_id = make_vocab_id(code, target)
            if vocab_id in seen:
                continue
            seen.add(vocab_id)
            core.append(
                {
                    "id": vocab_id,
                    "language": code,
                    "russian": key if code == "ru" else ru,
                    "english": target,
                    "translations": {"en": key},
                    "context": f"{ru}; {key}" if ru != key else key,
                    "level": level,
                    "topic": "multilingual-core",
                    "part_of_speech": pos,
                    "source": "Open multilingual core from WikDict/Kaikki pivot",
                    "frequency_rank": index,
                }
            )
        raw_by_code[code] = core + raw_by_code.get(code, [])


def prepend_english_common_core(raw_by_code: dict[str, list[dict]]) -> None:
    core: list[dict] = []
    seen: set[str] = set()
    for index, (word, pos, level) in enumerate(COMMON_CORE, start=1):
        key = english_key(word)
        if not key or key in seen:
            continue
        seen.add(key)
        core.append(
            {
                "id": make_vocab_id("en", key),
                "language": "en",
                "russian": COMMON_RUSSIAN_OVERRIDES.get(key, key),
                "english": key,
                "translations": {"en": key, "ru": COMMON_RUSSIAN_OVERRIDES.get(key, key)},
                "context": f"{COMMON_RUSSIAN_OVERRIDES.get(key, key)}; {key}",
                "level": level,
                "topic": "multilingual-core",
                "part_of_speech": pos,
                "source": "Open multilingual core from WikDict/Kaikki pivot",
                "frequency_rank": index,
            }
        )
    raw_by_code["en"] = core + raw_by_code.get("en", [])


def build_translations(
    code: str,
    target: str,
    gloss: str,
    ru_fallback: str,
    pos: str,
    en_to_ru: dict[str, str],
    en_to_lang: dict[str, dict[str, str]],
    en_to_lang_by_pos: dict[str, dict[str, dict[str, str]]],
    ru_to_lang: dict[str, dict[str, str]],
) -> dict[str, str]:
    key = english_key(gloss)
    translations: dict[str, str] = {}

    if key:
        translations["en"] = first_dictionary_value(gloss)
    if key and COMMON_RUSSIAN_OVERRIDES.get(key):
        translations["ru"] = COMMON_RUSSIAN_OVERRIDES[key]
    elif key and en_to_ru.get(key):
        translations["ru"] = first_dictionary_value(en_to_ru[key])
    elif clean_space(ru_fallback):
        translations["ru"] = first_dictionary_value(ru_fallback)

    for lang in LANGUAGE_ORDER:
        if lang == code:
            continue
        value = COMMON_TRANSLATION_OVERRIDES.get(key, {}).get(lang, "")
        ru_key = russian_key(ru_fallback)
        if not value and ru_key:
            value = ru_to_lang.get(lang, {}).get(ru_key, "")
        if not value:
            value = en_to_lang_by_pos.get(lang, {}).get(key, {}).get(pos_class(pos))
        if not value:
            value = en_to_lang.get(lang, {}).get(key)
        if value:
            translations[lang] = first_dictionary_value(normalize_script_value(lang, value))

    cleaned = {}
    for lang in INTERFACE_LANGUAGES:
        value = first_dictionary_value(translations.get(lang, ""))
        if value and target_key(value) != target_key(target):
            cleaned[lang] = value
    return cleaned


def split_dictionary_values(value: str) -> list[str]:
    value = clean_context(value)
    if not value:
        return []
    return [clean_space(part) for part in re.split(r";|/|,", value) if clean_space(part)]


def append_unique(values: list[str], value: str, code: str) -> None:
    value = clean_space(value)
    if not value or not is_learning_word(value, code):
        return
    key = target_key(value)
    if any(target_key(existing) == key for existing in values):
        return
    values.append(value)


def build_translation_banks(entries_by_code: dict[str, list[dict]]) -> dict[str, dict[str, dict[str, list[str]]]]:
    banks: dict[str, dict[str, dict[str, list[str]]]] = {}
    for code, entries in entries_by_code.items():
        for item in entries:
            translations = item.get("translations", {})
            if not isinstance(translations, dict):
                continue
            english = clean_space(str(translations.get("en", "") or (item.get("english", "") if code == "en" else "")))
            key = english_key(english)
            if not key:
                continue
            term = clean_space(str(item.get("english", "")))
            if not is_learning_word(term, code):
                continue
            pos = pos_class(str(item.get("part_of_speech", "")))
            values = banks.setdefault(code, {}).setdefault(key, {}).setdefault(pos, [])
            append_unique(values, term, code)
    return banks


def bank_values(
    banks: dict[str, dict[str, dict[str, list[str]]]],
    code: str,
    key: str,
    pos: str,
) -> list[str]:
    by_key = banks.get(code, {}).get(key, {})
    values = list(by_key.get(pos_class(pos), []))
    if len(values) < 2 and pos_class(pos) != "word":
        values.extend(by_key.get("word", []))
    result: list[str] = []
    for value in values:
        append_unique(result, value, code)
    return result


def enrich_dictionary_prompts(entries_by_code: dict[str, list[dict]]) -> None:
    banks = build_translation_banks(entries_by_code)
    for code, entries in entries_by_code.items():
        for item in entries:
            translations = item.get("translations", {})
            if not isinstance(translations, dict):
                continue
            key = english_key(str(translations.get("en", "") or (item.get("english", "") if code == "en" else "")))
            if not key:
                continue
            if key in COMMON_CORE_KEYS:
                continue
            pos = str(item.get("part_of_speech", ""))
            target = target_key(str(item.get("english", "")))
            for lang in INTERFACE_LANGUAGES:
                if lang == code:
                    continue
                existing = clean_context(str(translations.get(lang, "")))
                if existing and has_context(existing):
                    continue
                values: list[str] = []
                for part in split_dictionary_values(existing):
                    append_unique(values, part, lang)
                for value in bank_values(banks, lang, key, pos):
                    if target_key(value) == target:
                        continue
                    append_unique(values, value, lang)
                    if len(values) >= 4:
                        break
                if len(values) > 1:
                    translations[lang] = "; ".join(values[:4])
                elif not existing and values:
                    translations[lang] = values[0]


def clean_russian_target(value: str) -> str:
    value = clean_space(value)
    value = re.sub(r"\[\[([^\]]+)\]\]", r"\1", value)
    value = re.sub(r"\([^)]*\)", " ", value)
    value = re.split(r";|/|,", value, maxsplit=1)[0]
    value = clean_space(strip_combining_marks(value))
    return value


def derive_russian_entries(raw_entries: list[dict]) -> list[dict]:
    result: list[dict] = []
    seen: set[str] = set()
    for raw in raw_entries:
        target = clean_russian_target(str(raw.get("translations", {}).get("ru", "") or raw.get("russian", "")))
        gloss = clean_space(str(raw.get("english", "")))
        if not target or not gloss or not is_learning_word(target, "ru"):
            continue
        target_id = make_vocab_id("ru", target)
        if target_id in seen:
            continue
        seen.add(target_id)
        translations = {"en": gloss}
        for code, value in raw.get("translations", {}).items():
            if code in {"en", "ru"}:
                continue
            value = clean_space(str(value))
            if value and target_key(value) != target_key(target):
                translations[code] = value
        entry = {
            "id": target_id,
            "language": "ru",
            "russian": gloss,
            "english": target,
            "translations": translations,
            "level": level_for_entry("ru", target, gloss, raw, len(result) + 1),
            "topic": raw.get("topic") or "frequency",
            "part_of_speech": raw.get("part_of_speech") or "word",
            "source": "Derived Russian learning dictionary from FreeDict/WikDict eng-rus base",
            "frequency_rank": len(result) + 1,
        }
        context = contextual_entry_value(raw, gloss)
        if context:
            entry["context"] = context
        result.append(entry)
    return result


def extend_english_entries(base_entries: list[dict], raw_by_code: dict[str, list[dict]]) -> list[dict]:
    result = list(base_entries)
    seen = {make_vocab_id("en", str(item.get("english", ""))) for item in result}
    seen.discard("")

    for code in LANGUAGE_ORDER:
        if code == "en":
            continue
        for raw in raw_by_code.get(code, []):
            gloss = clean_gloss(str(raw.get("translations", {}).get("en", "")))
            if not gloss or not is_learning_word(gloss, "en"):
                continue
            target_id = make_vocab_id("en", gloss)
            if target_id in seen:
                continue
            seen.add(target_id)
            rank = len(result) + 1
            result.append(
                {
                    "id": target_id,
                    "language": "en",
                    "russian": "",
                    "english": gloss,
                    "translations": {"en": gloss},
                    "level": cefr_level(rank),
                    "topic": "open-pivot",
                    "part_of_speech": raw.get("part_of_speech") or "word",
                    "source": f"Derived English learning dictionary from {code} open dictionary pivot",
                    "frequency_rank": rank,
                }
            )

    return result


def finalize_entries(
    code: str,
    raw_entries: list[dict],
    target_to_en: dict[str, str],
    en_to_ru: dict[str, str],
    en_to_lang: dict[str, dict[str, str]],
    en_to_lang_by_pos: dict[str, dict[str, dict[str, str]]],
    ru_to_lang: dict[str, dict[str, str]],
) -> list[dict]:
    result: list[dict] = []
    seen: set[str] = set()
    for raw in raw_entries:
        target = normalize_script_value(code, normalize_headword(code, str(raw.get("english", ""))))
        if not is_learning_word(target, code):
            continue
        target_id = make_vocab_id(code, target)
        if target_id in seen:
            continue
        seen.add(target_id)

        gloss = clean_space(raw.get("translations", {}).get("en", ""))
        if not gloss:
            gloss = target_to_en.get(target_key(target), "")
        if code == "en":
            gloss = target
        gloss = clean_gloss(gloss) or english_key(gloss)

        ru_fallback = clean_space(str(raw.get("russian", "")))
        pos = raw.get("part_of_speech") or "word"
        translations = build_translations(code, target, gloss, ru_fallback, pos, en_to_ru, en_to_lang, en_to_lang_by_pos, ru_to_lang)
        russian = translations.get("ru") or ru_fallback or gloss or target

        rank = len(result) + 1
        entry = {
            "id": target_id,
            "language": code,
            "russian": russian,
            "english": target,
            "translations": translations,
            "level": level_for_entry(code, target, gloss, raw, rank),
            "topic": raw.get("topic") or "open-dictionary",
            "part_of_speech": pos,
            "source": raw.get("source") or "Open dictionary import",
            "frequency_rank": raw.get("frequency_rank") or rank,
        }
        context = contextual_entry_value(raw, gloss)
        if context:
            entry["context"] = context
        result.append(entry)
    return result


def finalize_english_entries(
    raw_entries: list[dict],
    en_to_ru: dict[str, str],
    en_to_lang: dict[str, dict[str, str]],
    en_to_lang_by_pos: dict[str, dict[str, dict[str, str]]],
    ru_to_lang: dict[str, dict[str, str]],
) -> list[dict]:
    result: list[dict] = []
    seen: set[str] = set()
    for raw in raw_entries:
        target = clean_space(str(raw.get("english", "")))
        if not target:
            continue
        target_id = make_vocab_id("en", target)
        if not target_id or target_id in seen:
            continue
        seen.add(target_id)

        ru_fallback = clean_space(str(raw.get("russian", "")))
        pos = raw.get("part_of_speech") or "word"
        translations = build_translations("en", target, target, ru_fallback, pos, en_to_ru, en_to_lang, en_to_lang_by_pos, ru_to_lang)
        russian = translations.get("ru") or ru_fallback or target
        rank = len(result) + 1
        entry = {
            "id": target_id,
            "language": "en",
            "russian": russian,
            "english": target,
            "translations": translations,
            "level": level_for_entry("en", target, target, raw, rank),
            "topic": raw.get("topic") or "open-dictionary",
            "part_of_speech": pos,
            "source": raw.get("source") or "Open dictionary import",
            "frequency_rank": raw.get("frequency_rank") or rank,
        }
        context = contextual_entry_value(raw, target)
        if context:
            entry["context"] = context
        result.append(entry)
    return result


def write_vocab(code: str, entries: list[dict]) -> None:
    entries = [sanitize_output_entry(code, entry) for entry in sorted_by_learning_level(code, entries)]
    VOCABULARY_DIR.mkdir(parents=True, exist_ok=True)
    path = VOCABULARY_DIR / ("vocabulary_words.json" if code == "en" else f"vocabulary_words_{code}.json")
    tmp_path = path.with_suffix(path.suffix + ".tmp")
    tmp_path.write_text(json.dumps(entries, ensure_ascii=False, separators=(",", ":")) + "\n", encoding="utf-8")
    for attempt in range(8):
        try:
            tmp_path.replace(path)
            break
        except PermissionError:
            if attempt == 7:
                raise
            time.sleep(0.5)
            try:
                path.unlink(missing_ok=True)
            except PermissionError:
                time.sleep(0.5)
    print(f"{code}: {len(entries):,}")


def sanitize_output_entry(code: str, entry: dict) -> dict:
    entry = dict(entry)
    entry.pop("context", None)
    entry.pop("contexts", None)
    target = normalize_script_value(code, normalize_headword(code, str(entry.get("english", ""))))
    entry["english"] = target
    entry["russian"] = first_dictionary_value(str(entry.get("russian", "")))
    translations = {}
    for lang, value in (entry.get("translations") or {}).items():
        lang = normalize_language_code(str(lang))
        value = first_dictionary_value(normalize_script_value(lang, str(value)))
        if value and target_key(value) != target_key(target):
            translations[lang] = value
    if entry["russian"]:
        translations["ru"] = entry["russian"]
    if code == "en":
        translations["en"] = target
    elif translations.get("en"):
        translations["en"] = first_dictionary_value(translations["en"])
    entry["translations"] = {lang: translations[lang] for lang in INTERFACE_LANGUAGES if translations.get(lang)}
    return entry


def normalize_language_code(code: str) -> str:
    code = clean_space(code).lower()
    if code in {"cn", "zh-cn", "zh_hans"}:
        return "zh"
    if code in {"jp"}:
        return "ja"
    return code


def main() -> None:
    base_entries, en_to_ru = parse_tei_english_russian(SOURCES / BASE_EN_RU_TEI)

    en_to_lang: dict[str, dict[str, str]] = {"en": {}}
    en_to_lang_by_pos: dict[str, dict[str, dict[str, str]]] = {"en": {}}
    for item in base_entries:
        word = clean_space(str(item.get("english", "")))
        ru = clean_space(str(item.get("russian", "")))
        key = english_key(word)
        if key:
            en_to_lang["en"][key] = word
            en_to_lang_by_pos.setdefault("en", {}).setdefault(key, {})[pos_class(item.get("part_of_speech", ""))] = word
        if key and ru:
            en_to_ru.setdefault(key, ru)

    raw_by_code: dict[str, list[dict]] = {"en": base_entries}
    raw_by_code["ru"] = derive_russian_entries(base_entries)
    target_to_en_by_code: dict[str, dict[str, str]] = {}
    ru_to_lang: dict[str, dict[str, str]] = {}

    for code, filename in WIKDICT_TEI.items():
        raw_entries, target_to_en = parse_tei(code, SOURCES / filename)
        target_to_en_by_code[code] = target_to_en
        for entry in raw_entries:
            add_en_map(en_to_lang, en_to_lang_by_pos, code, entry["translations"]["en"], entry["english"], entry.get("part_of_speech", ""))
        raw_by_code[code] = raw_entries

    for code, filename in WIKDICT_EN_TARGET.items():
        direct_map, direct_by_pos = parse_tei_en_to_target(code, SOURCES / filename)
        en_to_lang.setdefault(code, {}).update(direct_map)
        en_to_lang_by_pos.setdefault(code, {})
        for key, pos_values in direct_by_pos.items():
            en_to_lang_by_pos[code].setdefault(key, {}).update(pos_values)

    for code, filename in KAIKKI_JSONL.items():
        raw_entries, target_to_en = parse_kaikki(code, SOURCES / filename, limit=KAIKKI_LIMITS.get(code, 30000))
        target_to_en_by_code[code] = target_to_en
        for entry in raw_entries:
            add_en_map(en_to_lang, en_to_lang_by_pos, code, entry["translations"]["en"], entry["english"], entry.get("part_of_speech", ""))
        if code == "ru":
            raw_by_code[code] = raw_by_code.get(code, []) + raw_entries
        else:
            raw_by_code[code] = raw_entries

    raw_by_code["en"] = extend_english_entries(raw_by_code.get("en", []), raw_by_code)

    prepend_english_common_core(raw_by_code)
    prepend_common_core(raw_by_code, en_to_ru, en_to_lang, en_to_lang_by_pos, ru_to_lang)

    entries_by_code: dict[str, list[dict]] = {}
    for code in LANGUAGE_ORDER:
        if code == "en":
            entries = finalize_english_entries(raw_by_code.get(code, []), en_to_ru, en_to_lang, en_to_lang_by_pos, ru_to_lang)
        else:
            entries = finalize_entries(
                code,
                raw_by_code.get(code, []),
                target_to_en_by_code.get(code, {}),
                en_to_ru,
                en_to_lang,
                en_to_lang_by_pos,
                ru_to_lang,
            )
        entries_by_code[code] = sorted_by_learning_level(code, entries)

    for code in LANGUAGE_ORDER:
        write_vocab(code, entries_by_code.get(code, []))


if __name__ == "__main__":
    main()
