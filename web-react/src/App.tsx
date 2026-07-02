import type { ComponentType, CSSProperties, FormEvent, ReactNode } from "react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { ClipboardEvent } from "react";
import type { PointerEvent as ReactPointerEvent } from "react";
import { createPortal } from "react-dom";
import type { LucideProps } from "lucide-react";
import {
  Award,
  AlertCircle,
  Activity,
  BarChart3,
  Bookmark,
  BookOpen,
  Bot,
  Brain,
  Bug,
  CalendarDays,
  Check,
  CheckCircle,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  CircleHelp,
  CircleDollarSign,
  Clock,
  Download,
  Flame,
  Crown,
  Diamond,
  FileImage,
  FileAudio,
  Image as ImageIcon,
  Languages,
  LayoutDashboard,
  ListChecks,
  MessageCircle,
  Mic,
  PenLine,
  Play,
  RefreshCw,
  Repeat2,
  Search,
  Send,
  Settings,
  ShieldCheck,
  Sparkles,
  Star,
  Plus,
  Target,
  Trash2,
  Trophy,
  Users,
  Volume2,
  Wand2,
  WifiOff,
  X,
} from "lucide-react";

declare global {
  interface Window {
    dataLayer?: unknown[];
    ym?: (...args: unknown[]) => void;
  }
}
import { ApiError, api, apiForm, loadSession, logout } from "./lib/api";
import type {
  AiTutorResponse,
  AiTutorCompletedLesson,
  AiTutorCompletedLessonsResponse,
  AiTutorStep,
  ChatMessage,
  ChoiceOption,
  LanguageOption,
  LeaderboardEntry,
  LevelQuestion,
  MistakeItem,
  NavItem,
  PhrasebookItem,
  PremiumPlan,
  SessionData,
  SpellingChallenge,
  TranslatorResult,
  UserProfile,
  ViewId,
  VocabularyItem,
  VocabularyPage,
  WordChallenge,
} from "./lib/types";
import { asText, compactNumber, percent, prettyDate } from "./lib/format";
import { appCopy, appNavDescription, cleanAppText, translateViewDetails } from "./lib/i18n";
import { cn } from "./lib/utils";
import { Button } from "./components/ui/button";
import ChatComponent, { type ChatConfig } from "./components/ui/chat-interface";
import { InstagramIcon, TelegramIcon, TikTokIcon, YouTubeIcon } from "./components/BrandIcons";
import { VoiceInput } from "./components/ui/voice-input";
import { AudioWaveButton } from "./components/ui/audio-wave-button";
import { AnimatedThemeToggle } from "./components/ui/animated-theme-toggle";
import { Slider } from "./components/ui/interfaces-slider";
import MorphingArrowButton from "./components/ui/morphing-arrow-button";
import { MorphButton } from "./components/ui/morph-button";
import LogoutButton from "./components/ui/logout-button";
import OTPDialog from "./components/ui/otpdialog";
import SignInPage, { type SignInMode, type Testimonial } from "./components/ui/sign-in";
import NotFoundPage from "./components/ui/404-page-not-found";
import { AudioUploadCard } from "./components/ui/audio-upload-card";
import { PasswordInputField } from "./components/ui/password-input";
import { Spinner } from "./components/ui/spinner";
import { useFileUpload } from "./components/ui/file-upload";
import { Calendar } from "./components/ui/calendar-rac";
import {
  Dialog,
  DialogBody,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
} from "./components/ui/dialog";

type Theme = "light" | "dark";
type ApiRecord = Record<string, unknown>;
type StatusKind = "ok" | "error" | "info";
type LucideIcon = ComponentType<LucideProps>;
type ToolMode = "translator" | "voice" | "image";
type TrainerTone = "success" | "warning" | "info";

const aiRouterTelegramURL = "https://t.me/AiRouterRu_bot";

const poliglotSocialLinks = [
  { name: "YouTube", href: "https://www.youtube.com/@neriva_app", label: "Open NERIVA on YouTube", Icon: YouTubeIcon },
  { name: "Instagram", href: "https://www.instagram.com/neriva.ru", label: "Open NERIVA on Instagram", Icon: InstagramIcon },
  { name: "TikTok", href: "https://tiktok.com/@nerivaru", label: "Open NERIVA on TikTok", Icon: TikTokIcon },
  { name: "Telegram", href: "https://t.me/NERIVAapp_bot", label: "Open NERIVA on Telegram", Icon: TelegramIcon },
] as const;

const dailyQuestTarget = {
  lesson: 2,
  practice: 3,
  roleplay: 1,
  pronunciation: 1,
  vocabulary: 1,
  mistakes: 1,
  listening: 2,
} as const;

function dailyQuestProgress(value: unknown, target: number) {
  const count = Math.floor(Number(value || 0));
  if (!Number.isFinite(count) || count <= 0) return 0;
  return Math.min(count, target);
}

type PronunciationProblem = {
  word?: string;
  spoken?: string;
  confidence?: number;
  issue?: string;
  tip?: string;
  expected_sound?: string;
  heard_sound?: string;
};

type PronunciationAssessment = {
  expected?: string;
  spoken?: string;
  score?: number;
  accent_strength?: number;
  fluency?: number;
  similarity?: number;
  average_confidence?: number;
  speech_rate_wpm?: number;
  problem_words?: PronunciationProblem[];
  phoneme_issues?: PronunciationProblem[];
  tips?: string[];
  feedback?: string;
  stress?: string;
  rhythm?: string;
  intonation?: string;
  corrected_text?: string;
  audio_model?: string;
};

type StoredPronunciationAssessment = PronunciationAssessment & {
  createdAt?: string;
  id?: string;
};

type TrainerResult = {
  tone: TrainerTone;
  title: string;
  body: string;
  selectedId?: string;
  correctId?: string;
  record?: ApiRecord;
  mode?: "words" | "word-game" | "spelling" | "level" | "mistakes";
};

type AwardInfo = {
  level: number;
  rank: string;
  story: string;
};

type LeaderboardMeta = {
  language: string;
  languageName: string;
};

type PaymentState = {
  plan: PremiumPlan;
  info?: ApiRecord | null;
};

type PremiumPlanFeature = {
  text: string;
  locked?: boolean;
  highlight?: boolean;
};

type PaymentNotice = {
  kind: "success" | "not-found";
  title: string;
  description: string;
  planTitle?: string;
  period?: string;
  premiumUntil?: string;
  amount?: string;
};

type PaymentHistoryItem = {
  id: string;
  date: string;
  plan: string;
  period: string;
  amount: string;
  method: string;
  status: string;
};

type PhrasebookSource = NonNullable<PhrasebookItem["source"]>;

type HabitDay = {
  login?: boolean;
  complete?: boolean;
  claimed?: boolean;
  claimedAt?: string;
  lessons?: number;
  practice?: number;
  voice?: number;
};

type DailyBonusNotice = {
  xp: number;
  streak: number;
};

type XPGainNotice = {
  xp: number;
  total?: number;
  title?: string;
};

type MistakeCategory = "grammar" | "word-order" | "vocabulary" | "politeness" | "spelling";
type AiTutorFeedbackState = NonNullable<AiTutorResponse["feedback"]>;

type TelegramCodeRequest = {
  token: string;
  botURL: string;
  expiresAt: string;
};

type TutorLessonWord = {
  id: string;
  word: string;
  translation: string;
  context?: string;
  example?: string;
  example_translation?: string;
  part_of_speech?: string;
  topic?: string;
  level?: string;
  audio_text?: string;
};

type TutorChoiceOption = {
  id: string;
  text: string;
  label?: string;
  why?: string;
  skill?: string;
  quality?: string;
  avoid?: boolean;
};

type TutorAnswerVariant = {
  id: string;
  label?: string;
  text: string;
  why?: string;
  level?: string;
  use_case?: string;
  avoid?: boolean;
};

type TutorLessonChoice = {
  prompt?: string;
  options?: TutorChoiceOption[];
  correct_answer_id?: string;
  feedback?: string;
};

type TutorLessonStep = {
  code: string;
  title: string;
  summary?: string;
};

type TutorTeachingPoint = {
  title?: string;
  pattern?: string;
  scene_order?: string[];
  explanation?: string;
  model_answer?: string;
};

type TutorFinalWordCheck = {
  prompt?: string;
  items?: TutorLessonChoice[];
  required_correct?: number;
  summary_pass?: string;
  summary_retry?: string;
};

type TutorSummaryBlock = {
  can_say?: string;
  strong_items?: string[];
  weak_items?: string[];
  next_review?: string;
};

type TutorLesson = {
  id: string;
  title: string;
  level: string;
  topic: string;
  variant_code?: string;
  variant_title?: string;
  focus_skill?: string;
  success_criteria?: string[];
  lesson_number?: number;
  course_size?: number;
  learning_language: string;
  interface_language: string;
  duration_minutes?: number;
  goal?: string;
  can_do?: string;
  scenario?: string;
  teaching_point?: TutorTeachingPoint;
  final_word_check?: TutorFinalWordCheck;
  tutor_summary?: TutorSummaryBlock;
  steps?: TutorLessonStep[];
  words?: TutorLessonWord[];
  grammar_title?: string;
  grammar?: string;
  mini_explanation?: string;
  choice?: TutorLessonChoice;
  checks?: TutorLessonChoice[];
  writing_task?: string;
  writing_expected?: string[];
  answer_variants?: TutorAnswerVariant[];
  listening_task?: string;
  listening_text?: string;
  listening_question?: string;
  listening_expected?: string[];
  pronunciation_text?: string;
  dialogue?: string[];
  dialogue_prompt?: string;
  dialogue_goal?: string;
  dialogue_variants?: TutorAnswerVariant[];
  review?: TutorLessonWord[];
  review_summary?: string[];
  srs?: {
    again?: string;
    hard?: string;
    good?: string;
    easy?: string;
  };
  next_actions?: string[];
  source?: string;
};

type TutorCompletedLessonRecord = {
  id: string;
  lessonId?: string;
  sessionId?: string;
  title: string;
  topic: string;
  level: string;
  completedAt: string;
  aiLesson?: AiTutorStep["lesson"];
  lesson?: TutorLesson;
};

function tutorStableHash(value: string) {
  let hash = 2166136261;
  for (let index = 0; index < value.length; index += 1) {
    hash ^= value.charCodeAt(index);
    hash = Math.imul(hash, 16777619);
  }
  return hash >>> 0;
}

function tutorStableShuffle<T>(items: T[], seed: string) {
  const result = [...items];
  let state = tutorStableHash(seed) || 1;
  for (let index = result.length - 1; index > 0; index -= 1) {
    state = Math.imul(state ^ (state >>> 13), 1103515245) + 12345;
    const swapIndex = Math.abs(state) % (index + 1);
    [result[index], result[swapIndex]] = [result[swapIndex], result[index]];
  }
  return result;
}

type TutorStageId = "words" | "explain" | "choice" | "writing" | "listening" | "pronunciation" | "dialogue" | "final-check" | "assessment" | "review";
type TutorFeedbackTone = "idle" | "success" | "error";
type TutorProgressSnapshot = {
  lessonId: string;
  currentStageIndex: number;
  unlockedStageIndex: number;
  wordQuizIndex: number;
  wordQuizSelection: string;
  wordQuizAttempts: number;
  wordQuizResults: Array<{ id: string; word: string; answer: string; attempts: number }>;
  choiceCheckIndex: number;
  choiceCheckResults: Array<{ id: string; answer: string }>;
  choiceSelection: string;
  finalCheckIndex: number;
  finalCheckSelection: string;
  finalCheckResults: Array<{ id: string; answer: string }>;
  tutorDraft: string;
  srsSelection: string;
  feedback: string;
  feedbackTone: TutorFeedbackTone;
  tutorPronunciation: PronunciationAssessment | null;
  tutorPronunciationCorrection: string;
  pendingAdvanceNote: string;
  completedNotes: Record<string, string>;
};

type TurnstileAPI = {
  render: (
    target: HTMLElement,
    options: {
      sitekey: string;
      callback?: (token: string) => void;
      "expired-callback"?: () => void;
      "error-callback"?: () => void;
    },
  ) => string;
  reset?: (widgetID?: string) => void;
  remove?: (widgetID?: string) => void;
};

type RoleplayScenario = {
  id: string;
  title: string;
  description: string;
  prompt?: string;
  goal?: string;
  learnerRole?: string;
  aiRole?: string;
  difficulty: string;
  icon: LucideIcon;
};

type OfflineDeckItem = {
  id: string;
  title: string;
  answer?: string;
  example?: string;
  source: "vocabulary" | "mistakes" | "phrasebook";
};

const legacyAwardRanks: Record<string, string[]> = {
  en: [
    "Alphabet Spark", "Word Collector", "Phrase Hunter", "Conversation Scout", "Polyglot Apprentice",
    "Lexicon Keeper", "Small Talk Master", "Grammar Pathfinder", "Phrase Alchemist", "Confident Speaker",
    "Speech Traveler", "Dialogue Builder", "Language Navigator", "Cultural Envoy", "Polyglot Strategist",
    "Accent Hunter", "Living Speech Master", "Language Architect", "Translation Legend", "Supreme Polyglot",
  ],
  ru: [
    "Искра алфавита", "Собиратель слов", "Охотник за фразами", "Разведчик диалога", "Ученик полиглота",
    "Хранитель словаря", "Мастер живой беседы", "Следопыт грамматики", "Алхимик фраз", "Уверенный спикер",
    "Путешественник речи", "Строитель диалогов", "Языковой навигатор", "Посол культур", "Полиглот-стратег",
    "Ловец акцентов", "Мастер живой речи", "Архитектор языка", "Легенда перевода", "Верховный полиглот",
  ],
};

const legacyAwardStories: Record<string, string[]> = {
  en: [
    "The Rosetta Stone carried one message in three scripts, helping scholars hear ancient Egypt again.",
    "In Nineveh, Ashurbanipal gathered thousands of clay tablets; his library shows how memory survives when words are collected.",
    "On the Silk Road, merchants learned greetings, prices, and promises in many languages; one useful phrase could open a city gate.",
    "Interpreters in ancient ports stood between ships and markets, turning suspicion into trade one sentence at a time.",
    "In Baghdad's House of Wisdom, scholars translated Greek, Persian, Syriac, and Sanskrit works, letting knowledge sit at one table.",
    "Medieval glossaries placed difficult words beside familiar ones, like small maps for readers crossing into another language world.",
    "Travelers such as Ibn Battuta learned how to ask for water, shelter, and news; everyday speech carried them farther than royal letters.",
    "Al-Kindi studied letter frequencies and helped create cryptanalysis, proving that grammar and patterns can unlock hidden messages.",
    "Renaissance translators connected old books with new readers; a good translation could turn distant philosophy into a living conversation.",
    "Gutenberg's press made books travel faster than any teacher could walk, giving learners a steady voice to return to.",
    "Zheng He's fleets carried interpreters across the Indian Ocean, turning harbors into meeting places of stories, goods, and names.",
    "King Sejong supported Hangul so more people could read and write; an alphabet can be an act of care for ordinary voices.",
    "When Champollion read Egyptian hieroglyphs, silent monuments began to speak again, bringing centuries back into the room.",
    "L. L. Zamenhof created Esperanto hoping strangers might meet more gently; even an invented language can carry a dream of peace.",
    "The Navajo Code Talkers used a living language to protect messages in World War II, showing that culture can become strategy.",
    "Swahili grew through coastal trade with Arabic, Bantu, Persian, and Portuguese echoes; an accent can be a shore where histories meet.",
    "Nushu, the women's script of Hunan, carried songs, letters, and friendship; some languages survive because people need their own door to speak.",
    "Sequoyah's Cherokee syllabary let a nation print newspapers and letters in its own signs; one alphabet can become a house for memory.",
    "The Universal Declaration of Human Rights became one of the most translated documents on Earth, helping a shared promise reach many homes.",
    "Unicode gave digital space to scripts from across the planet; today many alphabets can stand together on one screen.",
  ],
  ru: [
    "Розеттский камень сохранил одно послание в трех письменностях и помог ученым снова услышать Древний Египет.",
    "В Ниневии Ашшурбанипал собрал тысячи глиняных табличек; его библиотека показывает, как память живет, когда слова берегут.",
    "На Шелковом пути купцы учили приветствия, цены и обещания на разных языках; одна полезная фраза могла открыть городские ворота.",
    "Переводчики древних портов стояли между кораблями и рынками, превращая недоверие в торговлю по одному предложению за раз.",
    "В багдадском Доме мудрости переводили греческие, персидские, сирийские и санскритские тексты, сажая знания за один стол.",
    "Средневековые глоссарии ставили трудные слова рядом с знакомыми, как маленькие карты для перехода в другой языковой мир.",
    "Путешественники вроде Ибн Баттуты учились просить воду, ночлег и новости; обычная речь несла их дальше царских писем.",
    "Аль-Кинди изучал частоты букв и помог создать криптоанализ, доказав, что грамматика и узоры открывают скрытые сообщения.",
    "Переводчики Возрождения связывали старые книги с новыми читателями; хороший перевод превращал далекую философию в живой разговор.",
    "Печатный станок Гутенберга заставил книги путешествовать быстрее любого учителя и дал ученикам голос, к которому можно возвращаться.",
    "Флотилии Чжэн Хэ брали переводчиков через Индийский океан, превращая гавани во встречи историй, товаров и имен.",
    "Король Седжон поддержал хангыль, чтобы больше людей могли читать и писать; алфавит может быть заботой об обычных голосах.",
    "Когда Шампольон прочитал египетские иероглифы, молчавшие памятники снова заговорили и вернули века в нашу комнату.",
    "Л. Л. Заменгоф создал эсперанто с надеждой, что незнакомцы будут встречаться мягче; даже придуманный язык может нести мечту о мире.",
    "Навахо-кодовые радисты использовали живой язык для защиты сообщений во Второй мировой войне и показали, что культура может стать стратегией.",
    "Суахили рос на прибрежной торговле с отголосками арабских, банту, персидских и португальских слов; акцент бывает берегом встречи историй.",
    "Нюйшу, женское письмо Хунани, хранило песни, письма и дружбу; некоторые языки выживают, потому что людям нужна своя дверь для речи.",
    "Слоговая азбука чероки Секвойи позволила народу печатать газеты и письма своими знаками; один алфавит может стать домом памяти.",
    "Всеобщая декларация прав человека стала одним из самых переводимых документов на Земле и помогла общей надежде дойти до многих домов.",
    "Unicode дал цифровое место письменностям всей планеты; сегодня многие алфавиты могут стоять рядом на одном экране.",
  ],
};

const navItems: NavItem[] = [
  { id: "home", label: "Сегодня", shortLabel: "Сегодня", group: "learn", description: "Daily command center", icon: LayoutDashboard, accent: "blue" },
  { id: "tutor", label: "AI Репетитор", shortLabel: "Tutor", group: "learn", description: "One guided lesson", icon: Sparkles, accent: "gold" },
  { id: "words", label: "Learn words", shortLabel: "Words", group: "words", description: "New vocabulary", icon: Brain, accent: "violet" },
  { id: "word-game", label: "Review game", shortLabel: "Review", group: "words", description: "Multiple choice", icon: Repeat2, accent: "blue" },
  { id: "pronunciation", label: "Pronunciation", shortLabel: "Pronounce", group: "learn", description: "Voice heatmap", icon: Activity, accent: "blue" },
  { id: "shadowing", label: "Listening", shortLabel: "Listen", group: "learn", description: "Repeat and score", icon: Volume2, accent: "mint" },
  { id: "lesson", label: "Lesson", shortLabel: "Lesson", group: "learn", description: "AI micro lesson", icon: BookOpen, accent: "indigo" },
  { id: "practice", label: "Practice", shortLabel: "Practice", group: "learn", description: "Conversation coach", icon: MessageCircle, accent: "teal" },
  { id: "roleplay", label: "Roleplay", shortLabel: "Roleplay", group: "learn", description: "AI scenarios", icon: Sparkles, accent: "violet" },
  { id: "spelling", label: "Spelling", shortLabel: "Spell", group: "words", description: "Type from memory", icon: PenLine, accent: "rose" },
  { id: "vocabulary", label: "Vocabulary", shortLabel: "Vocab", group: "words", description: "Saved words", icon: ListChecks, accent: "teal" },
  { id: "phrasebook", label: "Notes", shortLabel: "Notes", group: "words", description: "Saved lesson notes", icon: Bookmark, accent: "gold" },
  { id: "offline", label: "Offline decks", shortLabel: "Offline", group: "words", description: "PWA mini decks", icon: WifiOff, accent: "mint" },
  { id: "level", label: "Level test", shortLabel: "Level", group: "growth", description: "CEFR assessment", icon: Target, accent: "gold" },
  { id: "progress", label: "Progress", shortLabel: "Stats", group: "growth", description: "XP and streak", icon: BarChart3, accent: "indigo" },
  { id: "awards", label: "Awards", shortLabel: "Awards", group: "growth", description: "20 trophy ranks", icon: Trophy, accent: "gold" },
  { id: "leaderboard", label: "Leaderboard", shortLabel: "Ranks", group: "growth", description: "Compare learners", icon: Star, accent: "violet" },
  { id: "limits", label: "Limits", shortLabel: "Limits", group: "account", description: "Daily usage", icon: ShieldCheck, accent: "rose" },
  { id: "mistakes", label: "Mistakes", shortLabel: "Errors", group: "account", description: "Repair drills", icon: Search, accent: "rose" },
  { id: "tools", label: "Tools", shortLabel: "Tools", group: "account", description: "Translate and convert", icon: Wand2, accent: "teal" },
  { id: "dashboard", label: "Teacher dashboard", shortLabel: "Dashboard", group: "account", description: "Progress and ops", icon: BarChart3, accent: "indigo" },
  { id: "referral", label: "Referrals", shortLabel: "Invite", group: "account", description: "Invite rewards", icon: Users, accent: "mint" },
  { id: "premium", label: "Premium", shortLabel: "Plans", group: "account", description: "Plans and payments", icon: Crown, accent: "gold" },
  { id: "settings", label: "Settings", shortLabel: "Settings", group: "account", description: "Language and account", icon: Settings, accent: "blue" },
];

const roleplayScenarios: RoleplayScenario[] = [
  {
    id: "restaurant",
    title: "Ресторан",
    description: "Заказ, уточнения, вежливые просьбы и реакция официанта.",
    prompt: "Roleplay: you are a waiter in a restaurant. Ask me what I would like to order, correct my phrasing, and keep the dialogue at my CEFR level.",
    difficulty: "A1-B1",
    icon: MessageCircle,
  },
  {
    id: "work",
    title: "Работа",
    description: "Созвон, задача, сроки, короткие деловые ответы.",
    prompt: "Roleplay: you are my colleague at work. Ask about a task deadline, correct my answer, and ask one follow-up question.",
    difficulty: "A2-B2",
    icon: Users,
  },
  {
    id: "travel",
    title: "Путешествие",
    description: "Аэропорт, отель, маршрут и бытовые уточнения.",
    prompt: "Roleplay: you are a travel assistant. Ask me about my hotel booking and help me answer naturally.",
    difficulty: "A1-B2",
    icon: Sparkles,
  },
  {
    id: "exam",
    title: "Экзамен",
    description: "Короткий устный ответ, аргументы и аккуратные исправления.",
    prompt: "Roleplay: you are an oral exam interviewer. Ask me one CEFR-level speaking question, evaluate my answer, and give a better version.",
    difficulty: "A2-C1",
    icon: Trophy,
  },
  {
    id: "small-talk",
    title: "Small talk",
    description: "Лёгкая беседа, вопросы и естественные реакции.",
    prompt: "Roleplay: start a casual small-talk conversation. Keep it natural, ask one question at a time, and correct me gently.",
    difficulty: "A1-B2",
    icon: MessageCircle,
  },
  {
    id: "hotel",
    title: "Hotel check-in",
    description: "Check in, ask about breakfast, Wi-Fi, room issues, and late checkout.",
    difficulty: "A1-B2",
    icon: BookOpen,
  },
  {
    id: "shopping",
    title: "Shopping",
    description: "Ask for size, color, price, return policy, and compare options.",
    difficulty: "A1-B1",
    icon: Star,
  },
  {
    id: "doctor",
    title: "Doctor visit",
    description: "Describe symptoms, ask for advice, understand simple instructions.",
    difficulty: "A2-B2",
    icon: ShieldCheck,
  },
  {
    id: "job-interview",
    title: "Job interview",
    description: "Introduce experience, answer strengths, goals, and follow-up questions.",
    difficulty: "B1-C1",
    icon: Users,
  },
  {
    id: "bank",
    title: "Bank and payment",
    description: "Ask about a card, transfer, fee, failed payment, and confirmation.",
    difficulty: "A2-B2",
    icon: CircleDollarSign,
  },
];

const chatViews: ViewId[] = ["lesson", "practice", "shadowing", "tools"];
const hiddenRibbonViews = new Set<ViewId>(["limits", "progress"]);
const clientCopyPreferredKeys = new Set([
  "phrasebook",
  "phrases",
  "nav_phrasebook_desc",
  "save_to_phrasebook",
  "saved_phrases",
  "remove",
  "add_sample",
  "phrase_source_manual",
  "phrase_source_lesson",
  "phrase_source_practice",
  "phrase_source_roleplay",
  "phrasebook_title",
  "phrasebook_body",
  "empty_phrasebook",
  "phrasebook_empty_hint",
  "phrase_saved",
  "not_found_kicker",
  "not_found_title",
  "not_found_body",
  "not_found_home",
  "daily_bonus_already_claimed",
  "bonus_claimed",
  "choose_another_scenario",
]);

function isGenericSectionCopy(value: string) {
  const text = value.trim();
  return /^[^\s:]+:\s*[a-z0-9_ -]+$/i.test(text) || /^(Section|Раздел|Mục)(?:\s*:.*)?$/i.test(text);
}

function isTutorWordLearnStage(stage?: string | null) {
  return Boolean(stage && /^word_learn_\d+$/.test(stage));
}

function isTutorWordRecallStage(stage?: string | null) {
  return Boolean(stage && /^word_recall_\d+$/.test(stage));
}

function isTutorKnownInstructionStage(stage?: string | null) {
  return Boolean(
    stage &&
      (stage === "story_intro" ||
        stage === "retell" ||
        stage === "production" ||
        stage === "lesson_feedback" ||
        stage === "review_schedule" ||
        stage === "complete" ||
        /^question_\d+$/.test(stage) ||
        isTutorWordLearnStage(stage) ||
        isTutorWordRecallStage(stage)),
  );
}

function localizedTutorOptionText(id: string, text: string, copy: (key: string, fallback: string) => string) {
  const normalized = cleanAppText(id || text).trim();
  const fallback = cleanAppText(text || id).trim().replaceAll("_", " ");
  const keys: Record<string, [string, string]> = {
    easy: ["ai_tutor_option_easy", "Easy"],
    good: ["ai_tutor_option_good", "Good"],
    hard: ["ai_tutor_option_hard", "Hard"],
    bad: ["ai_tutor_option_bad", "Bad"],
    tomorrow: ["ai_tutor_option_tomorrow", "Tomorrow"],
    "3_days": ["ai_tutor_option_3_days", "In 3 days"],
    "1_week": ["ai_tutor_option_1_week", "In 1 week"],
    no_review: ["ai_tutor_option_no_review", "No review"],
  };
  const entry = keys[normalized];
  return entry ? copy(entry[0], entry[1]) : fallback;
}

function cleanTutorPrompt(value: unknown) {
  const text = cleanAppText(value).trim();
  if (!text) return "";
  const lines = text
    .split(/\n+/)
    .map((line) => line.trim())
    .filter(Boolean)
    .filter((line) => !/^Question\s+\d+\s*$/i.test(line))
    .filter((line) => !isGenericSectionCopy(line))
    .filter((line) => !/^(Next|Continue|Дальше)$/i.test(line))
    .filter((line) => !/Напишите\s+(свой\s+)?ответ к активному уроку/i.test(line));
  const seen = new Set<string>();
  return lines
    .filter((line) => {
      const key = line.toLowerCase();
      if (seen.has(key)) return false;
      seen.add(key);
      return true;
    })
    .join("\n");
}

function aiTutorHistoryBody(lesson: AiTutorStep["lesson"] | undefined, step: AiTutorStep | null) {
  const story = cleanAppText(lesson?.story?.text_target).trim();
  if (story) return story;
  return cleanTutorPrompt(step?.instruction);
}

function sameTutorText(left: string, right: string) {
  const normalize = (value: string) => cleanTutorPrompt(value).replace(/\s+/g, " ").trim().toLowerCase();
  return Boolean(normalize(left) && normalize(left) === normalize(right));
}

function hasMojibakeText(value: string) {
  return /пїЅ|Гђ|Г‘|Ð|Ñ|Р[\u0400-\u045f]|С[\u0400-\u045f]|Г[\u0400-\u045f]|б»|бє|бї|аё|а№|а¤|аҐ|а¦|а§|а®|а°/u.test(value);
}

function isKnownViewValue(value: string | null) {
  return value === "progress" || Boolean(value && navItems.some((item) => item.id === value));
}

function normalizeView(value: string | null): ViewId {
  if (value === "progress") return "dashboard";
  if (value && navItems.some((item) => item.id === value)) return value as ViewId;
  return "home";
}

const appLoginPath = "/app/login";

function isAuthStandaloneRoute() {
  const path = location.pathname.replace(/\/+$/, "").toLowerCase();
  return document.documentElement.dataset.authStandalone === "true" || path === "/login" || path === appLoginPath;
}

function currentNotFoundPath() {
  if (isAuthStandaloneRoute()) return "";
  const normalizedPath = location.pathname.replace(/\/+$/, "").toLowerCase() || "/";
  const allowedPaths = new Set(["/", "/app", "/app/mobile", "/app/desktop"]);
  const view = new URLSearchParams(location.search).get("view");
  if (!allowedPaths.has(normalizedPath) || (view && !isKnownViewValue(view))) {
    return `${location.pathname}${location.search}`;
  }
  return "";
}

function useMobileUiLayout() {
  const queryText = "(max-width: 760px)";
  const getMobileState = () => {
    if (window.matchMedia(queryText).matches) return true;
    if (window.innerWidth <= 760) return true;
    if (window.matchMedia("(pointer: coarse)").matches) return true;
    return /Android|iPhone|iPad|iPod|Mobile/i.test(navigator.userAgent);
  };
  const [isMobile, setIsMobile] = useState(getMobileState);

  useEffect(() => {
    const query = window.matchMedia(queryText);
    const pointerQuery = window.matchMedia("(pointer: coarse)");
    const update = () => setIsMobile(getMobileState());
    update();
    query.addEventListener("change", update);
    pointerQuery.addEventListener("change", update);
    window.addEventListener("resize", update);
    return () => {
      query.removeEventListener("change", update);
      pointerQuery.removeEventListener("change", update);
      window.removeEventListener("resize", update);
    };
  }, []);

  return isMobile;
}

function navCopyKey(view: ViewId) {
  const map: Partial<Record<ViewId, string>> = {
    home: "today",
    tutor: "ai_tutor",
    lesson: "new_lesson",
    practice: "practice",
    roleplay: "roleplay",
    shadowing: "shadowing",
    pronunciation: "pronunciation",
    words: "learn_words",
    "word-game": "word_game",
    spelling: "spelling",
    vocabulary: "vocabulary",
    phrasebook: "phrasebook",
    offline: "offline_decks",
    level: "level_test",
    progress: "progress",
    awards: "awards",
    leaderboard: "leaders",
    limits: "limits",
    mistakes: "mistakes",
    tools: "tools",
    dashboard: "dashboard",
    referral: "referrals",
    premium: "premium",
    settings: "settings",
  };
  return map[view] || view;
}

const roleplayTitles: Record<string, Record<string, string>> = {
  restaurant: { en: "Restaurant", ru: "Ресторан", es: "Restaurante", fr: "Restaurant", de: "Restaurant", it: "Ristorante", pt: "Restaurante", tr: "Restoran", ar: "Ш§Щ„Щ…Ш·Ш№Щ…", zh: "й¤ђеЋ…", ja: "гѓ¬г‚№гѓ€гѓ©гѓі", ko: "м‹ќл‹№", hi: "а¤°аҐ‡а¤ёаҐЌа¤¤а¤°а¤ѕа¤‚", id: "Restoran", vi: "NhГ  hГ ng", th: "аёЈа№‰аёІаё™аё­аёІаё«аёІаёЈ", pl: "Restauracja", nl: "Restaurant", uk: "Ресторан", kk: "Мейрамхана" },
  work: { en: "Work call", ru: "Рабочий звонок", es: "Llamada de trabajo", fr: "Appel professionnel", de: "ArbeitsgesprГ¤ch", it: "Chiamata di lavoro", pt: "Chamada de trabalho", tr: "Д°Еџ gГ¶rГјЕџmesi", ar: "Щ…ЩѓШ§Щ„Щ…Ш© Ш№Щ…Щ„", zh: "е·ҐдЅњйЂљиЇќ", ja: "д»•дє‹гЃ®й›»и©±", ko: "м—…л¬ґ н†µн™”", hi: "а¤•а¤ѕа¤® а¤•аҐЂ а¤•аҐ‰а¤І", id: "Panggilan kerja", vi: "Cuб»™c gб»Ќi cГґng viб»‡c", th: "аёЄаёІаёўаё‡аёІаё™", pl: "Rozmowa sЕ‚uЕјbowa", nl: "Werkgesprek", uk: "Робочий дзвінок", kk: "ЖТ±мыс Т›оТЈырауы" },
  travel: { en: "Travel", ru: "Путешествие", es: "Viaje", fr: "Voyage", de: "Reise", it: "Viaggio", pt: "Viagem", tr: "Seyahat", ar: "Ш§Щ„ШіЩЃШ±", zh: "ж—…иЎЊ", ja: "ж—…иЎЊ", ko: "м—¬н–‰", hi: "а¤Їа¤ѕа¤¤аҐЌа¤°а¤ѕ", id: "Perjalanan", vi: "Du lб»‹ch", th: "аёЃаёІаёЈа№Ђаё”аёґаё™аё—аёІаё‡", pl: "PodrГіЕј", nl: "Reis", uk: "Подорож", kk: "Саяхат" },
  exam: { en: "Oral exam", ru: "Устный экзамен", es: "Examen oral", fr: "Examen oral", de: "MГјndliche PrГјfung", it: "Esame orale", pt: "Exame oral", tr: "SГ¶zlГј sД±nav", ar: "Ш§Ш®ШЄШЁШ§Ш± ШґЩЃЩ‡ЩЉ", zh: "еЏЈиЇ­иЂѓиЇ•", ja: "еЏЈй ­и©¦йЁ“", ko: "кµ¬м€  м‹њн—", hi: "а¤®аҐЊа¤–а¤їа¤• а¤Єа¤°аҐЂа¤•аҐЌа¤·а¤ѕ", id: "Ujian lisan", vi: "Thi nГіi", th: "аёЄаё­аёљаёћаё№аё”", pl: "Egzamin ustny", nl: "Mondeling examen", uk: "Усний іспит", kk: "Ауызша емтихан" },
  "small-talk": { en: "Small talk", ru: "Лёгкая беседа", es: "Charla informal", fr: "Petite conversation", de: "Small Talk", it: "Conversazione leggera", pt: "Conversa casual", tr: "KД±sa sohbet", ar: "Ш­ШЇЩЉШ« Щ‚ШµЩЉШ±", zh: "й—ІиЃЉ", ja: "й›‘и«‡", ko: "мЉ¤лЄ°н† нЃ¬", hi: "а¤№а¤ІаҐЌа¤•аҐЂ а¤¬а¤ѕа¤¤а¤љаҐЂа¤¤", id: "Obrolan ringan", vi: "TrГІ chuyб»‡n xГЈ giao", th: "аё„аёёаёўаёЄаё±а№‰аё™а№†", pl: "Small talk", nl: "Smalltalk", uk: "Невимушена бесіда", kk: "ЖеТЈіл У™ТЈгіме" },
  hotel: { en: "Hotel check-in", ru: "Заселение в отель", es: "Registro en hotel", fr: "ArrivГ©e Г  l'hГґtel", de: "Hotel-Check-in", it: "Check-in in hotel", pt: "Check-in no hotel", tr: "Otele giriЕџ", ar: "ШЄШіШ¬ЩЉЩ„ ШЇШ®Щ€Щ„ Ш§Щ„ЩЃЩ†ШЇЩ‚", zh: "й…’еє—е…ҐдЅЏ", ja: "гѓ›гѓ†гѓ«гѓЃг‚§гѓѓг‚Їг‚¤гѓі", ko: "нён…” мІґнЃ¬мќё", hi: "а¤№аҐ‹а¤џа¤І а¤љаҐ‡а¤•-а¤‡а¤Ё", id: "Check-in hotel", vi: "Nhбє­n phГІng khГЎch sбєЎn", th: "а№ЂаёЉа№‡аёЃаё­аёґаё™а№‚аёЈаё‡а№ЃаёЈаёЎ", pl: "Meldunek w hotelu", nl: "Hotel inchecken", uk: "Заселення в готель", kk: "ТљонаТ›ТЇйге кіру" },
  shopping: { en: "Shopping", ru: "Покупки", es: "Compras", fr: "Achats", de: "Einkaufen", it: "Shopping", pt: "Compras", tr: "AlД±ЕџveriЕџ", ar: "Ш§Щ„ШЄШіЩ€Щ‚", zh: "иґ­з‰©", ja: "иІ·гЃ„з‰©", ko: "м‡јн•‘", hi: "а¤–а¤°аҐЂа¤¦а¤ѕа¤°аҐЂ", id: "Belanja", vi: "Mua sбєЇm", th: "аёЉа№‰аё­аё›аё›аёґа№‰аё‡", pl: "Zakupy", nl: "Winkelen", uk: "Покупки", kk: "Сауда" },
  doctor: { en: "Doctor visit", ru: "Визит к врачу", es: "Visita al mГ©dico", fr: "Rendez-vous mГ©dical", de: "Arztbesuch", it: "Visita medica", pt: "Consulta mГ©dica", tr: "Doktor ziyareti", ar: "ШІЩЉШ§Ш±Ш© Ш§Щ„Ш·ШЁЩЉШЁ", zh: "зњ‹еЊ»з”џ", ja: "з—…й™ўгЃ§гЃ®дјљи©±", ko: "лі‘м›ђ л°©л¬ё", hi: "а¤ЎаҐ‰а¤•аҐЌа¤џа¤° а¤ёаҐ‡ а¤®а¤їа¤Іа¤Ёа¤ѕ", id: "Kunjungan dokter", vi: "KhГЎm bГЎc sД©", th: "а№„аё›аёћаёља№Ѓаёћаё—аёўа№Њ", pl: "Wizyta u lekarza", nl: "Doktersbezoek", uk: "Візит до лікаря", kk: "ДУ™рігерге бару" },
  "job-interview": { en: "Job interview", ru: "Собеседование", es: "Entrevista laboral", fr: "Entretien d'embauche", de: "VorstellungsgesprГ¤ch", it: "Colloquio di lavoro", pt: "Entrevista de emprego", tr: "Д°Еџ gГ¶rГјЕџmesi", ar: "Щ…Щ‚Ш§ШЁЩ„Ш© Ш№Щ…Щ„", zh: "ж±‚иЃЊйќўиЇ•", ja: "йќўжЋҐ", ko: "л©ґм ‘", hi: "а¤ЁаҐЊа¤•а¤°аҐЂ а¤ёа¤ѕа¤•аҐЌа¤·а¤ѕа¤¤аҐЌа¤•а¤ѕа¤°", id: "Wawancara kerja", vi: "Phб»Џng vбєҐn xin viб»‡c", th: "аёЄаё±аёЎаё аёІаё©аё“а№Њаё‡аёІаё™", pl: "Rozmowa kwalifikacyjna", nl: "Sollicitatiegesprek", uk: "Співбесіда", kk: "ЖТ±мыс сТ±хбаты" },
  bank: { en: "Bank and payment", ru: "Банк и платеж", es: "Banco y pago", fr: "Banque et paiement", de: "Bank und Zahlung", it: "Banca e pagamento", pt: "Banco e pagamento", tr: "Banka ve Г¶deme", ar: "Ш§Щ„ШЁЩ†Щѓ Щ€Ш§Щ„ШЇЩЃШ№", zh: "й“¶иЎЊе’Њд»ж¬ѕ", ja: "йЉЂиЎЊгЃЁж”Їж‰•гЃ„", ko: "мќЂн–‰кіј кІ°м њ", hi: "а¤¬аҐ€а¤‚а¤• а¤”а¤° а¤­аҐЃа¤—а¤¤а¤ѕа¤Ё", id: "Bank dan pembayaran", vi: "NgГўn hГ ng vГ  thanh toГЎn", th: "аёаё™аёІаё„аёІаёЈа№ЃаёҐаё°аёЃаёІаёЈаёЉаёіаёЈаё°а№Ђаё‡аёґаё™", pl: "Bank i pЕ‚atnoЕ›Д‡", nl: "Bank en betaling", uk: "Банк і платіж", kk: "Банк жУ™не тУ©лем" },
};

const roleplayDescriptions: Record<string, Record<string, string>> = {
  ru: {
    restaurant: "Заказ, уточнения, вежливые просьбы и реакция официанта.",
    work: "Задача, срок, статус и короткие деловые ответы.",
    travel: "Аэропорт, отель, маршрут и бытовые уточнения.",
    exam: "Короткий устный ответ, аргументы и аккуратное исправление.",
    "small-talk": "Лёгкая беседа, вопросы и естественные реакции.",
    hotel: "Заселение, завтрак, Wi-Fi, проблема с номером и поздний выезд.",
    shopping: "Размер, цвет, цена, возврат и сравнение вариантов.",
    doctor: "Симптомы, совет, запись и простые инструкции.",
    "job-interview": "Опыт, сильные стороны, цели и уточняющие вопросы.",
    bank: "Карта, перевод, комиссия, неуспешный платеж и подтверждение.",
  },
};

const concreteRoleplayTitles: Record<string, Record<string, string>> = {
  ru: {
    restaurant: "Ресторан",
    work: "Рабочий звонок",
    travel: "Путешествие",
    exam: "Устный экзамен",
    "small-talk": "Легкая беседа",
    hotel: "Заселение в отель",
    shopping: "Покупки",
    doctor: "Визит к врачу",
    "job-interview": "Собеседование",
    bank: "Банк и платеж",
  },
  en: { restaurant: "Restaurant", work: "Work call", travel: "Travel", exam: "Oral exam", "small-talk": "Small talk", hotel: "Hotel check-in", shopping: "Shopping", doctor: "Doctor visit", "job-interview": "Job interview", bank: "Bank and payment" },
  es: { restaurant: "Restaurante", work: "Llamada de trabajo", travel: "Viaje", exam: "Examen oral", "small-talk": "Charla informal", hotel: "Registro en hotel", shopping: "Compras", doctor: "Visita al medico", "job-interview": "Entrevista laboral", bank: "Banco y pago" },
  de: { restaurant: "Restaurant", work: "Arbeitsgesprach", travel: "Reise", exam: "Mundliche Prufung", "small-talk": "Small Talk", hotel: "Hotel-Check-in", shopping: "Einkaufen", doctor: "Arztbesuch", "job-interview": "Vorstellungsgesprach", bank: "Bank und Zahlung" },
  fr: { restaurant: "Restaurant", work: "Appel professionnel", travel: "Voyage", exam: "Examen oral", "small-talk": "Petite conversation", hotel: "Arrivee a l'hotel", shopping: "Achats", doctor: "Rendez-vous medical", "job-interview": "Entretien d'embauche", bank: "Banque et paiement" },
  it: { restaurant: "Ristorante", work: "Chiamata di lavoro", travel: "Viaggio", exam: "Esame orale", "small-talk": "Conversazione leggera", hotel: "Check-in in hotel", shopping: "Shopping", doctor: "Visita medica", "job-interview": "Colloquio di lavoro", bank: "Banca e pagamento" },
  zh: { restaurant: "餐厅", work: "工作电话", travel: "旅行", exam: "口语考试", "small-talk": "闲聊", hotel: "酒店入住", shopping: "购物", doctor: "看医生", "job-interview": "求职面试", bank: "银行和付款" },
  ja: { restaurant: "レストラン", work: "仕事の電話", travel: "旅行", exam: "口頭試験", "small-talk": "雑談", hotel: "ホテルのチェックイン", shopping: "買い物", doctor: "病院での相談", "job-interview": "就職面接", bank: "銀行と支払い" },
  ko: { restaurant: "식당", work: "업무 통화", travel: "여행", exam: "구술 시험", "small-talk": "가벼운 대화", hotel: "호텔 체크인", shopping: "쇼핑", doctor: "병원 방문", "job-interview": "취업 면접", bank: "은행과 결제" },
  tg: { restaurant: "Тарабхона", work: "Занги корӣ", travel: "Сафар", exam: "Имтиҳони шифоҳӣ", "small-talk": "Суҳбати кӯтоҳ", hotel: "Қайд дар меҳмонхона", shopping: "Харид", doctor: "Назди духтур", "job-interview": "Мусоҳибаи корӣ", bank: "Бонк ва пардохт" },
  uz: { restaurant: "Restoran", work: "Ish qo'ng'irog'i", travel: "Sayohat", exam: "Og'zaki imtihon", "small-talk": "Qisqa suhbat", hotel: "Mehmonxonaga joylashish", shopping: "Xarid", doctor: "Shifokorga borish", "job-interview": "Ish suhbati", bank: "Bank va to'lov" },
  tt: { restaurant: "Ресторан", work: "Эш шалтыратуы", travel: "Сәяхәт", exam: "Телдән имтихан", "small-talk": "Кыска сөйләшү", hotel: "Кунакханәгә урнашу", shopping: "Кибеттә сатып алу", doctor: "Табибка бару", "job-interview": "Эш әңгәмәсе", bank: "Банк һәм түләү" },
  hy: { restaurant: "Ռեստորան", work: "Աշխատանքային զանգ", travel: "Ճանապարհորդություն", exam: "Բանավոր քննություն", "small-talk": "Կարճ զրույց", hotel: "Հյուրանոցի գրանցում", shopping: "Գնումներ", doctor: "Այց բժշկին", "job-interview": "Աշխատանքի հարցազրույց", bank: "Բանկ և վճարում" },
  kk: { restaurant: "Мейрамхана", work: "Жұмыс қоңырауы", travel: "Саяхат", exam: "Ауызша емтихан", "small-talk": "Жеңіл әңгіме", hotel: "Қонақүйге тіркелу", shopping: "Сауда", doctor: "Дәрігерге бару", "job-interview": "Жұмыс сұхбаты", bank: "Банк және төлем" },
  ky: { restaurant: "Ресторан", work: "Жумуш чалуу", travel: "Саякат", exam: "Оозеки экзамен", "small-talk": "Кыска сүйлөшүү", hotel: "Мейманканага катталуу", shopping: "Сатып алуу", doctor: "Дарыгерге баруу", "job-interview": "Жумуш маеги", bank: "Банк жана төлөм" },
  ka: { restaurant: "რესტორანი", work: "სამუშაო ზარი", travel: "მოგზაურობა", exam: "ზეპირი გამოცდა", "small-talk": "მოკლე საუბარი", hotel: "სასტუმროში რეგისტრაცია", shopping: "შოპინგი", doctor: "ექიმთან ვიზიტი", "job-interview": "სამუშაო გასაუბრება", bank: "ბანკი და გადახდა" },
  uk: { restaurant: "Ресторан", work: "Робочий дзвінок", travel: "Подорож", exam: "Усний іспит", "small-talk": "Невимушена бесіда", hotel: "Заселення в готель", shopping: "Покупки", doctor: "Візит до лікаря", "job-interview": "Співбесіда", bank: "Банк і платіж" },
  pl: { restaurant: "Restauracja", work: "Rozmowa służbowa", travel: "Podróż", exam: "Egzamin ustny", "small-talk": "Luźna rozmowa", hotel: "Meldunek w hotelu", shopping: "Zakupy", doctor: "Wizyta u lekarza", "job-interview": "Rozmowa kwalifikacyjna", bank: "Bank i płatność" },
  ro: { restaurant: "Restaurant", work: "Apel de lucru", travel: "Călătorie", exam: "Examen oral", "small-talk": "Conversație scurtă", hotel: "Check-in la hotel", shopping: "Cumpărături", doctor: "Vizită la medic", "job-interview": "Interviu de angajare", bank: "Bancă și plată" },
  pt: { restaurant: "Restaurante", work: "Chamada de trabalho", travel: "Viagem", exam: "Exame oral", "small-talk": "Conversa casual", hotel: "Check-in no hotel", shopping: "Compras", doctor: "Consulta médica", "job-interview": "Entrevista de emprego", bank: "Banco e pagamento" },
  ar: { restaurant: "مطعم", work: "مكالمة عمل", travel: "سفر", exam: "امتحان شفهي", "small-talk": "حديث قصير", hotel: "تسجيل دخول الفندق", shopping: "تسوق", doctor: "زيارة الطبيب", "job-interview": "مقابلة عمل", bank: "بنك ودفع" },
  bn: { restaurant: "রেস্তোরাঁ", work: "কাজের কল", travel: "ভ্রমণ", exam: "মৌখিক পরীক্ষা", "small-talk": "হালকা আলাপ", hotel: "হোটেল চেক-ইন", shopping: "কেনাকাটা", doctor: "ডাক্তারের কাছে", "job-interview": "চাকরির সাক্ষাৎকার", bank: "ব্যাংক ও পেমেন্ট" },
  cs: { restaurant: "Restaurace", work: "Pracovní hovor", travel: "Cestování", exam: "Ústní zkouška", "small-talk": "Krátký rozhovor", hotel: "Check-in v hotelu", shopping: "Nakupování", doctor: "Návštěva lékaře", "job-interview": "Pracovní pohovor", bank: "Banka a platba" },
  el: { restaurant: "Εστιατόριο", work: "Επαγγελματική κλήση", travel: "Ταξίδι", exam: "Προφορική εξέταση", "small-talk": "Σύντομη συζήτηση", hotel: "Check-in σε ξενοδοχείο", shopping: "Αγορές", doctor: "Επίσκεψη σε γιατρό", "job-interview": "Συνέντευξη εργασίας", bank: "Τράπεζα και πληρωμή" },
  hi: { restaurant: "रेस्तरां", work: "काम की कॉल", travel: "यात्रा", exam: "मौखिक परीक्षा", "small-talk": "हल्की बातचीत", hotel: "होटल चेक-इन", shopping: "खरीदारी", doctor: "डॉक्टर से मिलना", "job-interview": "नौकरी इंटरव्यू", bank: "बैंक और भुगतान" },
  hu: { restaurant: "Étterem", work: "Munkahelyi hívás", travel: "Utazás", exam: "Szóbeli vizsga", "small-talk": "Könnyed beszélgetés", hotel: "Szállodai bejelentkezés", shopping: "Vásárlás", doctor: "Orvosi látogatás", "job-interview": "Állásinterjú", bank: "Bank és fizetés" },
  id: { restaurant: "Restoran", work: "Panggilan kerja", travel: "Perjalanan", exam: "Ujian lisan", "small-talk": "Obrolan ringan", hotel: "Check-in hotel", shopping: "Belanja", doctor: "Kunjungan dokter", "job-interview": "Wawancara kerja", bank: "Bank dan pembayaran" },
  nl: { restaurant: "Restaurant", work: "Werkgesprek", travel: "Reis", exam: "Mondeling examen", "small-talk": "Smalltalk", hotel: "Hotel inchecken", shopping: "Winkelen", doctor: "Doktersbezoek", "job-interview": "Sollicitatiegesprek", bank: "Bank en betaling" },
  sv: { restaurant: "Restaurang", work: "Jobbsamtal", travel: "Resa", exam: "Muntligt prov", "small-talk": "Småprat", hotel: "Hotellincheckning", shopping: "Shopping", doctor: "Läkarbesök", "job-interview": "Jobbintervju", bank: "Bank och betalning" },
  ta: { restaurant: "உணவகம்", work: "வேலை அழைப்பு", travel: "பயணம்", exam: "வாய்மொழி தேர்வு", "small-talk": "சிறு உரையாடல்", hotel: "ஹோட்டல் செக்-இன்", shopping: "கடைசெய்தல்", doctor: "மருத்துவர் சந்திப்பு", "job-interview": "வேலை நேர்காணல்", bank: "வங்கி மற்றும் கட்டணம்" },
  te: { restaurant: "రెస్టారెంట్", work: "పని కాల్", travel: "ప్రయాణం", exam: "మౌఖిక పరీక్ష", "small-talk": "చిన్న సంభాషణ", hotel: "హోటల్ చెక్-ఇన్", shopping: "షాపింగ్", doctor: "డాక్టర్ సందర్శన", "job-interview": "ఉద్యోగ ఇంటర్వ్యూ", bank: "బ్యాంక్ మరియు చెల్లింపు" },
  th: { restaurant: "ร้านอาหาร", work: "สายงาน", travel: "การเดินทาง", exam: "สอบพูด", "small-talk": "คุยสั้นๆ", hotel: "เช็กอินโรงแรม", shopping: "ช็อปปิ้ง", doctor: "ไปพบแพทย์", "job-interview": "สัมภาษณ์งาน", bank: "ธนาคารและการชำระเงิน" },
  tl: { restaurant: "Restawran", work: "Tawag sa trabaho", travel: "Paglalakbay", exam: "Oral exam", "small-talk": "Maikling usapan", hotel: "Hotel check-in", shopping: "Pamimili", doctor: "Pagbisita sa doktor", "job-interview": "Job interview", bank: "Bangko at bayad" },
  tr: { restaurant: "Restoran", work: "İş görüşmesi", travel: "Seyahat", exam: "Sözlü sınav", "small-talk": "Kısa sohbet", hotel: "Otele giriş", shopping: "Alışveriş", doctor: "Doktor ziyareti", "job-interview": "İş mülakatı", bank: "Banka ve ödeme" },
  vi: { restaurant: "Nhà hàng", work: "Cuộc gọi công việc", travel: "Du lịch", exam: "Thi nói", "small-talk": "Trò chuyện xã giao", hotel: "Nhận phòng khách sạn", shopping: "Mua sắm", doctor: "Khám bác sĩ", "job-interview": "Phỏng vấn xin việc", bank: "Ngân hàng và thanh toán" },
};

const roleplayGenericDescriptions: Record<string, string> = {
  en: "Practice this real-life situation with a clear goal, one question at a time, and short corrections.",
  ru: "Отработайте реальную ситуацию: цель, один вопрос за раз и короткие исправления.",
  es: "Practica esta situaciГіn real con un objetivo claro, una pregunta por vez y correcciones breves.",
  fr: "Travaillez cette situation rГ©elle avec un objectif clair, une question Г  la fois et des corrections courtes.",
  de: "Гњbe diese Alltagssituation mit klarem Ziel, einer Frage nach der anderen und kurzen Korrekturen.",
  it: "Esercita questa situazione reale con un obiettivo chiaro, una domanda alla volta e correzioni brevi.",
  pt: "Pratique esta situaГ§ГЈo real com uma meta clara, uma pergunta por vez e correГ§Гµes curtas.",
  tr: "Bu gerГ§ek durumu net bir hedefle, tek tek sorularla ve kД±sa dГјzeltmelerle Г§alД±Еџ.",
  ar: "ШЄШЇШ±Щ‘ШЁ Ш№Щ„Щ‰ Щ‡Ш°Ш§ Ш§Щ„Щ…Щ€Щ‚ЩЃ Ш§Щ„Щ€Ш§Щ‚Ш№ЩЉ ШЁЩ‡ШЇЩЃ Щ€Ш§Ш¶Ш­ШЊ ШіШ¤Ш§Щ„ Щ€Ш§Ш­ШЇ ЩЃЩЉ ЩѓЩ„ Щ…Ш±Ш© Щ€ШЄШµШ­ЩЉШ­Ш§ШЄ Щ‚ШµЩЉШ±Ш©.",
  zh: "з”ЁжЋзЎ®з›®ж ‡з»ѓд№ иї™дёЄзњџе®ћењєж™ЇпјЊдёЂж¬ЎдёЂдёЄй—®йўпјЊе№¶з»™е‡єз®Ђзџ­зє ж­ЈгЂ‚",
  ja: "жЋзўєгЃЄз›®жЁ™гЃ§гЃ“гЃ®е®џз”Ёе ґйќўг‚’з·ґзї’гЃ—гЂЃдёЂе•ЏгЃљгЃ¤зџ­гЃЏдї®ж­ЈгЃ—гЃѕгЃ™гЂ‚",
  ko: "лЄ…н™•н•њ лЄ©н‘њлЎњ м‹¤м њ мѓЃн™©мќ„ м—°мЉµн•кі , н•њ лІ€м—ђ н•њ м§€л¬ём”© м§§кІЊ кµђм •н•©л‹€л‹¤.",
  hi: "а¤ёаҐЌа¤Єа¤·аҐЌа¤џ а¤Іа¤•аҐЌа¤·аҐЌа¤Ї, а¤Џа¤• а¤¬а¤ѕа¤° а¤®аҐ‡а¤‚ а¤Џа¤• а¤ЄаҐЌа¤°а¤¶аҐЌа¤Ё а¤”а¤° а¤›аҐ‹а¤џаҐ‡ а¤ёаҐЃа¤§а¤ѕа¤°аҐ‹а¤‚ а¤•аҐ‡ а¤ёа¤ѕа¤Ґ а¤Їа¤№ а¤µа¤ѕа¤ёаҐЌа¤¤а¤µа¤їа¤• а¤ёаҐЌа¤Ґа¤їа¤¤а¤ї а¤…а¤­аҐЌа¤Їа¤ѕа¤ё а¤•а¤°аҐ‡а¤‚аҐ¤",
  id: "Latih situasi nyata ini dengan tujuan jelas, satu pertanyaan sekali, dan koreksi singkat.",
  vi: "Luyб»‡n tГ¬nh huб»‘ng thб»±c tбєї nГ y vб»›i mб»Ґc tiГЄu rГµ rГ ng, tб»«ng cГўu hб»Џi mб»™t vГ  sб»­a lб»—i ngбєЇn.",
  th: "аёќаё¶аёЃаёЄаё–аёІаё™аёЃаёІаёЈаё“а№Њаё€аёЈаёґаё‡аё™аёµа№‰аё”а№‰аё§аёўа№Ђаё›а№‰аёІаё«аёЎаёІаёўаёЉаё±аё”а№Ђаё€аё™ аё—аёµаёҐаё°аё„аёіаё–аёІаёЎ а№ЃаёҐаё°а№ЃаёЃа№‰аёЄаё±а№‰аё™а№†",
  tg: "Ин Тіолати воТ›еиро бо Тіадафи равшан, як савол дар як ваТ›т ва ислоТіи кУЇтоТі машТ› кунед.",
  uz: "Bu real vaziyatni aniq maqsad, bitta savol va qisqa tuzatishlar bilan mashq qiling.",
  tt: "Бу чын хУ™лне ачык максат, берьюлы бер сорау Т»У™м кыска тУ©зУ™тТЇлУ™р белУ™н кТЇнегегез.",
  hy: "ХЉХЎЦЂХЎХєХҐЦ„ ХЎХµХЅ Х«ЦЂХЎХЇХЎХ¶ Х«ЦЂХЎХѕХ«ХіХЎХЇХЁ Х°ХЅХїХЎХЇ Х¶ХєХЎХїХЎХЇХёХѕ, ХґХҐХЇ Х°ХЎЦЂЦЃХёХѕ Ц‡ ХЇХЎЦЂХі ХёЦ‚ХІХІХёЦ‚ХґХ¶ХҐЦЂХёХѕЦ‰",
  ky: "Бул реалдуу кырдаалды так максат, бирден суроо жана кыска оТЈдоолор менен машыгыТЈыз.",
  ka: "бѓбѓ•бѓђбѓ бѓЇбѓбѓЁбѓ”бѓ— бѓ”бѓЎ бѓ бѓ”бѓђбѓљбѓЈбѓ бѓ бѓЎбѓбѓўбѓЈбѓђбѓЄбѓбѓђ бѓ›бѓ™бѓђбѓ¤бѓбѓќ бѓ›бѓбѓ–бѓњбѓбѓ—, бѓ—бѓбѓ—бѓќ бѓ™бѓбѓ—бѓ®бѓ•бѓбѓ— бѓ“бѓђ бѓ›бѓќбѓ™бѓљбѓ” бѓЁбѓ”бѓЎбѓ¬бѓќбѓ бѓ”бѓ‘бѓ”бѓ‘бѓбѓ—.",
  pl: "Д†wicz tД™ realnД… sytuacjД™ z jasnym celem, jednym pytaniem naraz i krГіtkД… korektД….",
  ro: "ExerseazДѓ aceastДѓ situaИ›ie realДѓ cu un obiectiv clar, cГўte o Г®ntrebare И™i corecturi scurte.",
  nl: "Oefen deze echte situatie met een duidelijk doel, Г©Г©n vraag tegelijk en korte correcties.",
  uk: "Відпрацюйте реальну ситуацію з чіткою метою, одним питанням за раз і короткими виправленнями.",
  kk: "НаТ›ты жаТ“дайды айТ›ын маТ›сатпен, бір сТ±раТ›тан жУ™не Т›ысТ›а тТЇзетумен жаттыТ›тырыТЈыз.",
};

const pronunciationLocales: Record<string, Record<string, string>> = {
  en: { title: "Pronunciation", body: "A standalone pronunciation workout: listen to the model, record or upload your voice, and get a score with weak sounds.", heatmap: "Pronunciation heatmap", weak: "Weak words and sounds", history: "Score history", noData: "No voice data yet", hint: "Record or upload the phrase above.", latest: "latest score", repeat: "Try again" },
  ru: { title: "Произношение", body: "Самостоятельная тренировка произношения: послушайте образец, запишите или загрузите голос и получите оценку по словам и звукам.", heatmap: "Карта произношения", weak: "Слабые слова и звуки", history: "История прогресса", noData: "Пока нет голосовых данных", hint: "Запишите или загрузите фразу выше.", latest: "последняя оценка", repeat: "Повторить" },
  es: { title: "PronunciaciГіn", body: "Palabras dГ©biles, sonidos difГ­ciles e historial de puntuaciГіn de escucha y respuestas de voz.", heatmap: "Mapa de pronunciaciГіn", weak: "Palabras y sonidos dГ©biles", history: "Historial de puntuaciГіn", noData: "AГєn no hay datos de voz", hint: "Inicia escucha y envГ­a una respuesta de voz.", latest: "Гєltima puntuaciГіn", repeat: "Repetir" },
  de: { title: "Aussprache", body: "Schwache WГ¶rter, schwierige Laute und Bewertungsverlauf aus Listening und Sprachantworten.", heatmap: "Aussprache-Karte", weak: "Schwache WГ¶rter und Laute", history: "Bewertungsverlauf", noData: "Noch keine Sprachdaten", hint: "Starte Listening und sende eine Sprachantwort.", latest: "letzte Bewertung", repeat: "Wiederholen" },
  fr: { title: "Prononciation", body: "Mots faibles, sons difficiles et historique des scores issus de l'Г©coute et des rГ©ponses vocales.", heatmap: "Carte de prononciation", weak: "Mots et sons faibles", history: "Historique des scores", noData: "Pas encore de donnГ©es vocales", hint: "Lance l'Г©coute et envoie une rГ©ponse vocale.", latest: "dernier score", repeat: "RГ©essayer" },
  it: { title: "Pronuncia", body: "Parole deboli, suoni difficili e cronologia dei punteggi da ascolto e risposte vocali.", heatmap: "Mappa pronuncia", weak: "Parole e suoni deboli", history: "Cronologia punteggi", noData: "Nessun dato vocale ancora", hint: "Avvia l'ascolto e invia una risposta vocale.", latest: "ultimo punteggio", repeat: "Riprova" },
  zh: { title: "еЏ‘йџі", body: "жќҐи‡Єеђ¬еЉ›е’ЊиЇ­йџіе›ћз­”зљ„и–„еј±иЇЌгЂЃйљѕз‚№еЏ‘йџіе’Ње€†ж•°еЋ†еЏІгЂ‚", heatmap: "еЏ‘йџізѓ­еЉ›е›ѕ", weak: "и–„еј±иЇЌе’ЊеЈ°йџі", history: "е€†ж•°еЋ†еЏІ", noData: "иїжІЎжњ‰иЇ­йџіж•°жЌ®", hint: "ејЂе§‹еђ¬еЉ›е№¶еЏ‘йЂЃиЇ­йџіе›ћз­”гЂ‚", latest: "жњЂж–°е€†ж•°", repeat: "е†ЌиЇ•дёЂж¬Ў" },
  ja: { title: "з™єйџі", body: "гѓЄг‚№гѓ‹гѓіг‚°гЃЁйџіеЈ°е›ћз­”гЃ‹г‚‰еј±гЃ„еЌиЄћгЂЃй›ЈгЃ—гЃ„йџігЂЃг‚№г‚іг‚ўе±Ґж­ґг‚’иЎЁз¤єгЃ—гЃѕгЃ™гЂ‚", heatmap: "з™єйџігѓ’гѓјгѓ€гѓћгѓѓгѓ—", weak: "еј±гЃ„еЌиЄћгЃЁйџі", history: "г‚№г‚іг‚ўе±Ґж­ґ", noData: "йџіеЈ°гѓ‡гѓјг‚їгЃЇгЃѕгЃ гЃ‚г‚ЉгЃѕгЃ›г‚“", hint: "гѓЄг‚№гѓ‹гѓіг‚°г‚’й–‹е§‹гЃ—гЃ¦йџіеЈ°е›ћз­”г‚’йЂЃгЃЈгЃ¦гЃЏгЃ гЃ•гЃ„гЂ‚", latest: "жњЂж–°г‚№г‚іг‚ў", repeat: "г‚‚гЃ†дёЂеє¦" },
  ko: { title: "л°њмќЊ", body: "л“Јкё°м™Ђ мќЊм„± л‹µліЂм—ђм„њ м•Ѕн•њ л‹Ём–ґ, м–ґл ¤мљґ м†Њл¦¬, м ђм€ кё°лЎќмќ„ ліґм—¬м¤Ќл‹€л‹¤.", heatmap: "л°њмќЊ нћ€нЉёл§µ", weak: "м•Ѕн•њ л‹Ём–ґм™Ђ м†Њл¦¬", history: "м ђм€ кё°лЎќ", noData: "м•„м§Ѓ мќЊм„± лЌ°мќґн„°к°Ђ м—†мЉµл‹€л‹¤", hint: "л“Јкё°лҐј м‹њмћ‘н•кі  мќЊм„± л‹µліЂмќ„ ліґл‚ґм„ёмљ”.", latest: "мµњк·ј м ђм€", repeat: "л‹¤м‹њ м‹њлЏ„" },
  tg: { title: "Талаффуз", body: "КалимаТіои заиф, садоТіои душвор ва таърихи холТіо аз шунидан ва Т·авобТіои овозУЈ.", heatmap: "Харитаи талаффуз", weak: "КалимаТіо ва садоТіои заиф", history: "Таърихи холТіо", noData: "ТІоло маълумоти овозУЈ нест", hint: "Listening-ро оТ“оз кунед ва Т·авоби овозУЈ фиристед.", latest: "холи охирин", repeat: "Такрор" },
  uz: { title: "Talaffuz", body: "Listening va ovozli javoblardan zaif so'zlar, qiyin tovushlar va ball tarixi.", heatmap: "Talaffuz xaritasi", weak: "Zaif so'zlar va tovushlar", history: "Ball tarixi", noData: "Hali ovoz ma'lumoti yo'q", hint: "Listeningni boshlang va ovozli javob yuboring.", latest: "so'nggi ball", repeat: "Takrorlash" },
  tt: { title: "Уйтелеш", body: "ТыТЈлау Т»У™м тавыш Т—авапларыннан зУ™гыйфь сТЇзлУ™р, катлаулы авазлар Т»У™м баллар тарихы.", heatmap: "Уйтелеш картасы", weak: "ЗУ™гыйфь сТЇзлУ™р Т»У™м авазлар", history: "Баллар тарихы", noData: "УлегУ™ тавыш мУ™гълТЇматы юк", hint: "Listening башлап тавыш Т—авабы Т—ибУ™регез.", latest: "соТЈгы балл", repeat: "Кабатлау" },
  hy: { title: "Ф±ЦЂХїХЎХЅХЎХ¶ХёЦ‚Х©ХµХёЦ‚Х¶", body: "Ф№ХёЦ‚ХµХ¬ ХўХЎХјХҐЦЂ, Х¤ХЄХѕХЎЦЂ Х°Х¶Х№ХµХёЦ‚Х¶Х¶ХҐЦЂ Ц‡ ХЈХ¶ХЎХ°ХЎХїХЎХЇХЎХ¶Х¶ХҐЦЂХ« ХєХЎХїХґХёЦ‚Х©ХµХёЦ‚Х¶ Х¬ХЅХҐХ¬ХёЦ‚ЦЃ ХёЦ‚ Х±ХЎХµХ¶ХЎХµХ«Х¶ ХєХЎХїХЎХЅХ­ХЎХ¶Х¶ХҐЦЂХ«ЦЃЦ‰", heatmap: "Ф±ЦЂХїХЎХЅХЎХ¶ХёЦ‚Х©ХµХЎХ¶ Ц„ХЎЦЂХїХҐХ¦", weak: "Ф№ХёЦ‚ХµХ¬ ХўХЎХјХҐЦЂ Ц‡ Х°Х¶Х№ХµХёЦ‚Х¶Х¶ХҐЦЂ", history: "ФіХ¶ХЎХ°ХЎХїХЎХЇХЎХ¶Х¶ХҐЦЂХ« ХєХЎХїХґХёЦ‚Х©ХµХёЦ‚Х¶", noData: "ХЃХЎХµХ¶ХЎХµХ«Х¶ ХїХѕХµХЎХ¬Х¶ХҐЦЂ Х¤ХҐХј Х№ХЇХЎХ¶", hint: "ХЌХЇХЅХҐЦ„ Listening-ХЁ Ц‡ ХёЦ‚ХІХЎЦЂХЇХҐЦ„ Х±ХЎХµХ¶ХЎХµХ«Х¶ ХєХЎХїХЎХЅХ­ХЎХ¶Ц‰", latest: "ХѕХҐЦЂХ»Х«Х¶ ХЈХ¶ХЎХ°ХЎХїХЎХЇХЎХ¶", repeat: "ФїЦЂХЇХ¶ХҐХ¬" },
  kk: { title: "Айтылым", body: "ТыТЈдау мен дауыстыТ› жауаптардан У™лсіз сУ©здер, кТЇрделі дыбыстар жУ™не балл тарихы.", heatmap: "Айтылым картасы", weak: "Улсіз сУ©здер мен дыбыстар", history: "Балл тарихы", noData: "Узірге дауыс деректері жоТ›", hint: "Listening бастаТЈыз жУ™не дауыстыТ› жауап жіберіТЈіз.", latest: "соТЈТ“ы балл", repeat: "Тљайталау" },
  ky: { title: "Айтылыш", body: "Угуудан жана ТЇн жоопторунан алсыз сУ©здУ©р, татаал ТЇндУ©р жана упай тарыхы.", heatmap: "Айтылыш картасы", weak: "Алсыз сУ©здУ©р жана ТЇндУ©р", history: "Упай тарыхы", noData: "Азырынча ТЇн маалыматы жок", hint: "Listening баштап ТЇн жообун жУ©нУ©тТЇТЈТЇз.", latest: "акыркы упай", repeat: "Кайталоо" },
  ka: { title: "бѓ’бѓђбѓ›бѓќбѓ—бѓҐбѓ›бѓђ", body: "бѓЎбѓЈбѓЎбѓўбѓ бѓЎбѓбѓўбѓ§бѓ•бѓ”бѓ‘бѓ, бѓ бѓ—бѓЈбѓљбѓ бѓ‘бѓ’бѓ”бѓ бѓ”бѓ‘бѓ бѓ“бѓђ бѓҐбѓЈбѓљбѓ”бѓ‘бѓбѓЎ бѓбѓЎбѓўбѓќбѓ бѓбѓђ бѓ›бѓќбѓЎбѓ›бѓ”бѓњбѓбѓ“бѓђбѓњ бѓ“бѓђ бѓ®бѓ›бѓќбѓ•бѓђбѓњбѓ бѓћбѓђбѓЎбѓЈбѓ®бѓ”бѓ‘бѓбѓ“бѓђбѓњ.", heatmap: "бѓ’бѓђбѓ›бѓќбѓ—бѓҐбѓ›бѓбѓЎ бѓ бѓЈбѓ™бѓђ", weak: "бѓЎбѓЈбѓЎбѓўбѓ бѓЎбѓбѓўбѓ§бѓ•бѓ”бѓ‘бѓ бѓ“бѓђ бѓ‘бѓ’бѓ”бѓ бѓ”бѓ‘бѓ", history: "бѓҐбѓЈбѓљбѓ”бѓ‘бѓбѓЎ бѓбѓЎбѓўбѓќбѓ бѓбѓђ", noData: "бѓ®бѓ›бѓќбѓ•бѓђбѓњбѓ бѓ›бѓќбѓњбѓђбѓЄбѓ”бѓ›бѓ”бѓ‘бѓ бѓЇбѓ”бѓ  бѓђбѓ  бѓђбѓ бѓбѓЎ", hint: "бѓ“бѓђбѓбѓ¬бѓ§бѓ”бѓ— Listening бѓ“бѓђ бѓ’бѓђбѓ’бѓ–бѓђбѓ•бѓњбѓ”бѓ— бѓ®бѓ›бѓќбѓ•бѓђбѓњбѓ бѓћбѓђбѓЎбѓЈбѓ®бѓ.", latest: "бѓ‘бѓќбѓљбѓќ бѓҐбѓЈбѓљбѓђ", repeat: "бѓ’бѓђбѓ›бѓ”бѓќбѓ бѓ”бѓ‘бѓђ" },
  uk: { title: "Вимова", body: "Слабкі слова, складні звуки й історія оцінок із аудіювання та голосових відповідей.", heatmap: "Карта вимови", weak: "Слабкі слова й звуки", history: "Історія оцінок", noData: "Поки немає голосових даних", hint: "Запустіть аудіювання і надішліть голосову відповідь.", latest: "остання оцінка", repeat: "Повторити" },
  pl: { title: "Wymowa", body: "SЕ‚abe sЕ‚owa, trudne dЕєwiД™ki i historia wynikГіw z listeningu oraz odpowiedzi gЕ‚osowych.", heatmap: "Mapa wymowy", weak: "SЕ‚abe sЕ‚owa i dЕєwiД™ki", history: "Historia wynikГіw", noData: "Brak danych gЕ‚osowych", hint: "Uruchom Listening i wyЕ›lij odpowiedЕє gЕ‚osowД….", latest: "ostatni wynik", repeat: "PowtГіrz" },
  ro: { title: "PronunИ›ie", body: "Cuvinte slabe, sunete dificile И™i istoricul scorurilor din listening И™i rДѓspunsuri vocale.", heatmap: "HartДѓ de pronunИ›ie", weak: "Cuvinte И™i sunete slabe", history: "Istoric scoruri", noData: "Nu existДѓ Г®ncДѓ date vocale", hint: "PorneИ™te Listening И™i trimite un rДѓspuns vocal.", latest: "ultimul scor", repeat: "RepetДѓ" },
  pt: { title: "PronГєncia", body: "Palavras fracas, sons difГ­ceis e histГіrico de pontuaГ§ГЈo de listening e respostas de voz.", heatmap: "Mapa de pronГєncia", weak: "Palavras e sons fracos", history: "HistГіrico de pontuaГ§ГЈo", noData: "Ainda nГЈo hГЎ dados de voz", hint: "Inicie o Listening e envie uma resposta de voz.", latest: "Гєltima pontuaГ§ГЈo", repeat: "Repetir" },
};

const roleplayUiFallbacks: Record<string, { title: string; subtitle: string; result: string }> = {
  en: { title: "AI roleplay scenarios", subtitle: "Choose a situation: AI starts the dialogue, corrects phrases, and asks the next question.", result: "Roleplay result" },
  ru: { title: "AI-сценарии", subtitle: "Выберите ситуацию: AI начнет диалог, исправит фразы и задаст следующий вопрос.", result: "Результат сценария" },
  es: { title: "Escenarios AI", subtitle: "Elige una situaciГіn: la IA inicia el diГЎlogo, corrige frases y hace la siguiente pregunta.", result: "Resultado del escenario" },
  de: { title: "AI-Rollenspiele", subtitle: "WГ¤hle eine Situation: Die KI startet den Dialog, korrigiert SГ¤tze und stellt die nГ¤chste Frage.", result: "Rollenspiel-Ergebnis" },
  fr: { title: "ScГ©narios IA", subtitle: "Choisissez une situation : l'IA lance le dialogue, corrige les phrases et pose la question suivante.", result: "RГ©sultat du scГ©nario" },
  it: { title: "Scenari AI", subtitle: "Scegli una situazione: l'AI avvia il dialogo, corregge le frasi e fa la domanda successiva.", result: "Risultato scenario" },
  zh: { title: "AI жѓ…ж™ЇеЇ№иЇќ", subtitle: "йЂ‰ж‹©дёЂдёЄењєж™ЇпјљAI ејЂе§‹еЇ№иЇќгЂЃзє ж­ЈеЏҐе­ђе№¶жЏђе‡єдё‹дёЂдёЄй—®йўгЂ‚", result: "жѓ…ж™Їз»“жћњ" },
  ja: { title: "AIгѓ­гѓјгѓ«гѓ—гѓ¬г‚¤", subtitle: "е ґйќўг‚’йЃёгЃ¶гЃЁгЂЃAIгЃЊдјљи©±г‚’е§‹г‚ЃгЂЃиЎЁзЏѕг‚’з›ґгЃ—гЂЃж¬ЎгЃ®иіЄе•Џг‚’гЃ—гЃѕгЃ™гЂ‚", result: "гѓ­гѓјгѓ«гѓ—гѓ¬г‚¤зµђжћњ" },
  ko: { title: "AI лЎ¤н”Њл €мќґ", subtitle: "мѓЃн™©мќ„ м„ нѓќн•л©ґ AIк°Ђ лЊЂн™”лҐј м‹њмћ‘н•кі  л¬ёмћҐмќ„ кµђм •н•л©° л‹¤мќЊ м§€л¬ёмќ„ н•©л‹€л‹¤.", result: "лЎ¤н”Њл €мќґ кІ°кіј" },
  uk: { title: "AI-сценарії", subtitle: "Оберіть ситуацію: AI почне діалог, виправить фрази й поставить наступне питання.", result: "Результат сценарію" },
  pl: { title: "Scenariusze AI", subtitle: "Wybierz sytuacjД™: AI zacznie dialog, poprawi frazy i zada nastД™pne pytanie.", result: "Wynik scenariusza" },
  pt: { title: "CenГЎrios AI", subtitle: "Escolha uma situaГ§ГЈo: a IA inicia o diГЎlogo, corrige frases e faz a prГіxima pergunta.", result: "Resultado do cenГЎrio" },
};

function getRecord(payload: unknown): ApiRecord {
  return payload && typeof payload === "object" ? (payload as ApiRecord) : {};
}

function getPayloadUser(payload: unknown): UserProfile | undefined {
  if (payload && typeof payload === "object" && "user" in payload) {
    const user = (payload as ApiRecord).user;
    if (user && typeof user === "object") return user as UserProfile;
  }
  return undefined;
}

function toChoiceOptions(value: unknown): ChoiceOption[] {
  return Array.isArray(value)
    ? value
        .map((item) => (item && typeof item === "object" ? (item as ChoiceOption) : null))
        .filter((item): item is ChoiceOption => Boolean(item?.id && item?.text))
    : [];
}

function asTutorLesson(value: unknown): TutorLesson | null {
  if (!value || typeof value !== "object") return null;
  const lesson = value as TutorLesson;
  if (!lesson.id || !lesson.title) return null;
  return lesson;
}

function cleanTutorHistoryLabel(value: unknown, fallback = "AI Tutor") {
  const cleaned = cleanAppText(value).replace(/\s+/g, " ").trim();
  if (!cleaned) return fallback;
  const lower = cleaned.toLowerCase();
  const blockedPrefixes = ["to talk about", "topic/theme seed", "goal:", "lesson goal:", "return strict json", "create one complete"];
  if (blockedPrefixes.some((prefix) => lower.startsWith(prefix))) return fallback;
  if (/^(story|retell|question|new words|word check|writing|memory review|complete)$/i.test(cleaned)) return fallback;
  if (/^(topic|theme|level|section|раздел):\s*/i.test(cleaned)) return fallback;
  return cleaned;
}

function readTutorCompletedLessons(key: string): TutorCompletedLessonRecord[] {
  if (!key) return [];
  try {
    const raw = JSON.parse(localStorage.getItem(key) || "[]");
    if (!Array.isArray(raw)) return [];
    return raw
      .map((item) => {
        const record = getRecord(item);
        const lesson = asTutorLesson(record.lesson);
        const id = cleanAppText(record.id || lesson?.id).trim();
        if (!id) return null;
        const aiLessonRecord = getRecord(record.aiLesson || record.ai_lesson || record.lesson_payload);
        const aiLesson = Object.keys(aiLessonRecord).length ? (aiLessonRecord as AiTutorStep["lesson"]) : undefined;
        const completed: TutorCompletedLessonRecord = {
          id,
          lessonId: cleanAppText(record.lessonId || record.lesson_id || lesson?.id).trim(),
          sessionId: cleanAppText(record.sessionId || record.session_id).trim(),
          title: cleanTutorHistoryLabel(record.title || lesson?.title),
          topic: cleanTutorHistoryLabel(record.topic || lesson?.topic, ""),
          level: cleanAppText(record.level || lesson?.level).trim(),
          completedAt: cleanAppText(record.completedAt || record.completed_at).trim() || new Date().toISOString(),
          aiLesson,
          lesson: lesson || undefined,
        };
        return completed;
      })
      .filter((item): item is TutorCompletedLessonRecord => Boolean(item))
      .sort((left, right) => right.completedAt.localeCompare(left.completedAt));
  } catch {
    return [];
  }
}

function writeTutorCompletedLesson(key: string, lesson: { id: string; title: string; topic: string; level: string; lessonId?: string; sessionId?: string; aiLesson?: AiTutorStep["lesson"]; lesson?: TutorLesson }) {
  if (!key || !lesson?.id) return;
  const record: TutorCompletedLessonRecord = {
    id: lesson.id,
    lessonId: cleanAppText(lesson.lessonId).trim(),
    sessionId: cleanAppText(lesson.sessionId).trim(),
    title: cleanTutorHistoryLabel(lesson.title),
    topic: cleanTutorHistoryLabel(lesson.topic, ""),
    level: cleanAppText(lesson.level).trim(),
    completedAt: new Date().toISOString(),
    aiLesson: lesson.aiLesson,
    lesson: lesson.lesson,
  };
  const current = readTutorCompletedLessons(key).filter((item) => item.id !== record.id);
  localStorage.setItem(key, JSON.stringify([record, ...current].slice(0, 100)));
}

function tutorCompletedLessonFromAPI(item: AiTutorCompletedLesson): TutorCompletedLessonRecord | null {
  const lessonId = cleanAppText(item.lesson_id).trim();
  const sessionId = cleanAppText(item.session_id).trim();
  const id = sessionId || lessonId;
  if (!id) return null;
  const aiLessonRecord = getRecord(item.lesson);
  const aiLesson = Object.keys(aiLessonRecord).length ? (aiLessonRecord as AiTutorStep["lesson"]) : undefined;
  const title = cleanTutorHistoryLabel(item.title || aiLesson?.title || item.topic);
  return {
    id,
    lessonId,
    sessionId,
    title,
    topic: cleanTutorHistoryLabel(item.topic || aiLesson?.theme || title, ""),
    level: cleanAppText(item.level || aiLesson?.level).trim(),
    completedAt: cleanAppText(item.completed_at).trim() || new Date().toISOString(),
    aiLesson,
  };
}

function mergeTutorCompletedLessons(primary: TutorCompletedLessonRecord[], fallback: TutorCompletedLessonRecord[]) {
  const result: TutorCompletedLessonRecord[] = [];
  const seen = new Set<string>();
  const add = (item: TutorCompletedLessonRecord) => {
    const keys = [item.sessionId, item.lessonId, item.id].filter(Boolean) as string[];
    if (keys.some((key) => seen.has(key))) return;
    keys.forEach((key) => seen.add(key));
    result.push(item);
  };
  primary.forEach(add);
  fallback.forEach(add);
  return result.sort((left, right) => right.completedAt.localeCompare(left.completedAt));
}

function panelMessage(message: string, tone: ChatMessage["tone"] = "default", title = "NERIVA", details?: ApiRecord, meta?: string): ChatMessage {
  return {
    id: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
    role: tone === "danger" ? "system" : "assistant",
    title,
    body: message,
    tone,
    details,
    meta,
  };
}

function userMessage(message: string, title: string, meta?: string): ChatMessage {
  return {
    id: `${Date.now()}-user-${Math.random().toString(36).slice(2)}`,
    role: "user",
    title,
    body: message,
    meta,
  };
}

function challengeDetails(challenge: WordChallenge | SpellingChallenge | null): ApiRecord {
  if (!challenge) return {};
  return {
    prompt: challenge.prompt,
    context: challenge.context,
    direction: challenge.direction,
    question_audio_text: challenge.prompt,
  };
}

function challengeMessage(challenge: WordChallenge | SpellingChallenge, copy: (key: string, fallback: string) => string) {
  const prompt = cleanAppText(challenge.prompt || challenge.message);
  const context = cleanAppText(challenge.context);
  const direction = cleanAppText(challenge.direction);
  const lines = [prompt, context ? `${copy("context", "Context")}: ${context}` : "", direction ? `${copy("route", "Route")}: ${direction}` : ""];
  return lines.filter(Boolean).join("\n\n") || copy("loading", "Loading...");
}

function updateSessionUser(session: SessionData | null, user: UserProfile): SessionData {
  return {
    ...(session || { authenticated: true }),
    authenticated: true,
    user,
  };
}

function recordField(record: ApiRecord, keys: string[]) {
  for (const key of keys) {
    const value = cleanAppText(record[key]);
    if (value) return value;
  }
  return "";
}

function recordList(value: unknown): string[] {
  if (!Array.isArray(value)) return [];
  return value.map((item) => cleanAppText(item)).filter(Boolean);
}

function recordMistakes(value: unknown) {
  if (!Array.isArray(value)) return [];
  return value
    .map((item) => (item && typeof item === "object" ? (item as ApiRecord) : null))
    .filter((item): item is ApiRecord => Boolean(item));
}

function clampPercentScore(value: unknown, scaleUnitValue = false) {
  const numeric = Number(value);
  if (!Number.isFinite(numeric)) return 0;
  const scaled = scaleUnitValue && numeric > 0 && numeric <= 1 ? numeric * 100 : numeric;
  return Math.max(0, Math.min(100, Math.round(scaled)));
}

function pronunciationFeedbackScore(value: unknown) {
  const text = cleanAppText(value);
  const preferred = text.match(/(?:score|оценка|балл)[^0-9]{0,16}(\d{1,3})\s*\/\s*100/i);
  if (preferred) return clampPercentScore(preferred[1]);
  const matches = Array.from(text.matchAll(/(\d{1,3})\s*\/\s*100/g));
  const last = matches[matches.length - 1];
  return last ? clampPercentScore(last[1]) : 0;
}

function normalizedPronunciationScore(report: Partial<PronunciationAssessment> | null | undefined) {
  if (!report) return 0;
  const feedbackScore = pronunciationFeedbackScore(report.feedback);
  const rawScore = Number(report.score);
  if (Number.isFinite(rawScore)) {
    if (rawScore > 100 && feedbackScore) return feedbackScore;
    return clampPercentScore(rawScore);
  }
  const similarity = Number(report.similarity);
  if (Number.isFinite(similarity)) {
    if (similarity > 100 && feedbackScore) return feedbackScore;
    return clampPercentScore(similarity, true);
  }
  const averageConfidence = Number(report.average_confidence);
  if (Number.isFinite(averageConfidence)) {
    if (averageConfidence > 100 && feedbackScore) return feedbackScore;
    return clampPercentScore(averageConfidence, true);
  }
  return feedbackScore;
}

function pronunciationFrom(record: ApiRecord): PronunciationAssessment | null {
  const value = record.pronunciation;
  return value && typeof value === "object" ? (value as PronunciationAssessment) : null;
}

function pronunciationStorageKey(session: SessionData, user: UserProfile) {
  return `poliglot-pronunciation-v2:${asText(session.account?.login || user.telegram_account?.id || "guest", "guest")}`;
}

function pronunciationFingerprint(report: PronunciationAssessment) {
  const words = (report.problem_words || [])
    .map((item) => `${cleanPronunciationWord(item.word)}:${cleanPronunciationWord(item.spoken)}:${asText(item.issue, "")}`)
    .join("|");
  return [
    cleanAppText(report.expected),
    cleanAppText(report.spoken),
    normalizedPronunciationScore(report),
    cleanAppText(report.feedback),
    words,
  ].join("::");
}

function readPronunciationHistory(key: string): StoredPronunciationAssessment[] {
  try {
    const parsed = JSON.parse(localStorage.getItem(key) || "[]");
    return Array.isArray(parsed) ? parsed.filter((item) => item && typeof item === "object") as StoredPronunciationAssessment[] : [];
  } catch {
    return [];
  }
}

function mergePronunciationHistory(current: PronunciationAssessment[], stored: StoredPronunciationAssessment[]) {
  const seen = new Set<string>();
  const merged: StoredPronunciationAssessment[] = [];
  [...current.map((report) => ({ ...report, createdAt: new Date().toISOString() })), ...stored].forEach((report) => {
    const fingerprint = pronunciationFingerprint(report);
    if (seen.has(fingerprint)) return;
    seen.add(fingerprint);
    merged.push(report);
  });
  return merged.slice(0, 16);
}

function formatPronunciationPlain(record: ApiRecord, copy: (key: string, fallback: string) => string) {
  const p = pronunciationFrom(record);
  if (!p) return "";
  const lines = [
    `${copy("pronunciation_score", "Pronunciation")}: ${normalizedPronunciationScore(p)}/100`,
    `${copy("accent_strength", "Accent")}: ${clampPercentScore(p.accent_strength, true)}/100`,
    `${copy("fluency", "Fluency")}: ${clampPercentScore(p.fluency, true)}/100`,
  ];
  if (p.feedback) lines.push(cleanAppText(p.feedback));
  const problemWords = Array.isArray(p.problem_words)
    ? p.problem_words
      .slice(0, 3)
      .map((item) => cleanPronunciationWord(item.word || item.spoken))
      .filter(Boolean)
    : [];
  if (problemWords.length) {
    lines.push(`${copy("problem_words", "Weak words")}: ${problemWords.join(", ")}`);
  }
  return lines.join("\n");
}

function pronunciationHistoryText(report: PronunciationAssessment, copy: (key: string, fallback: string) => string) {
  const feedback = cleanAppText(report.feedback || "");
  const scorePattern = /^(?:Оценка|Score|Pronunciation|Произношение):\s*\d{1,3}\s*\/\s*100\.?\s*/i;
  const cleaned = feedback.replace(scorePattern, "").trim();
  return cleaned || cleanAppText(report.expected) || copy("pronunciation_history_empty", "Нет текстового комментария.");
}

type LearningRecordFormatOptions = {
  includeContext?: boolean;
  includeExample?: boolean;
  includeMistakes?: boolean;
};

function formatLearningRecord(
  record: ApiRecord,
  copy: (key: string, fallback: string) => string,
  fallback: string,
  options: LearningRecordFormatOptions = {},
) {
  const includeContext = options.includeContext ?? true;
  const includeExample = options.includeExample ?? true;
  const includeMistakes = options.includeMistakes ?? true;
  const lines: string[] = [];
  const main = recordField(record, ["feedback", "reply", "message", "result", "translation", "source_text"]);
  if (main) lines.push(main);

  const correction = recordField(record, ["correction", "model_phrase"]);
  if (correction && !lines.some((line) => line.includes(correction))) lines.push(`${copy("correct_variant", "Correct variant")}: ${correction}`);

  const explanation = recordField(record, ["explanation", "tip"]);
  if (explanation && !lines.some((line) => line.includes(explanation))) lines.push(explanation);

  const context = recordField(record, ["context"]);
  if (includeContext && context && !lines.some((line) => line.includes(context))) lines.push(`${copy("context", "Context")}: ${context}`);

  const example = recordField(record, ["example"]);
  if (includeExample && example) lines.push(`${copy("example", "Example")}: ${example}`);

  const mistakes = includeMistakes ? recordMistakes(record.mistakes) : [];
  if (mistakes.length) {
    lines.push(
      `${copy("mistakes", "Mistakes")}: ` +
        mistakes
          .slice(0, 4)
          .map((item) => `${cleanAppText(item.word)} → ${cleanAppText(item.correction)}${item.explanation ? ` (${cleanAppText(item.explanation)})` : ""}`)
          .join("; "),
    );
  }

  const pronunciation = formatPronunciationPlain(record, copy);
  if (pronunciation) lines.push(pronunciation);

  const promoted = recordField(record, ["promoted_to"]);
  if (promoted) lines.push(`${copy("level_updated", "Level updated")}: ${promoted}`);

  return lines.filter(Boolean).join("\n\n") || fallback;
}

function makeTrainerResult(
  record: ApiRecord,
  copy: (key: string, fallback: string) => string,
  selectedId?: string,
  mode?: TrainerResult["mode"],
): TrainerResult {
  const correct = Boolean(record.correct);
  const title = correct ? copy("correct", "Correct") : record.gave_up ? copy("correct_answer", "Correct answer") : copy("try_again", "Try again");
  const compactTrainer = mode === "words" || mode === "word-game" || mode === "spelling";
  const hideSpellingAnswer = mode === "spelling" && !correct && !record.gave_up;
  const body = formatLearningRecord(record, copy, recordField(record, ["word", "translation", "message"]) || title, {
    includeContext: !compactTrainer,
    includeExample: mode !== "spelling" || !hideSpellingAnswer,
    includeMistakes: !hideSpellingAnswer,
  });
  const correctId = recordField(record, ["correct_answer_id", "word_id"]) || (correct ? selectedId : undefined);
  return {
    tone: correct ? "success" : "warning",
    title,
    body,
    selectedId,
    correctId,
    record,
    mode,
  };
}

function awardLevel(user: UserProfile) {
  return Math.max(1, Math.min(20, Number(user.xp_level || 1)));
}

function awardFile(level: number) {
  return `/app/assets/award-${String(Math.max(1, Math.min(20, level))).padStart(2, "0")}.png`;
}

function awardInfo(level: number, copy: (key: string, fallback: string) => string, interfaceLanguage = "en"): AwardInfo {
  const normalized = Math.max(1, Math.min(20, Number(level) || 1));
  const key = String(normalized).padStart(2, "0");
  const language = (interfaceLanguage || "en").toLowerCase().split("-")[0];
  const fallbackRanks = legacyAwardRanks[language] || legacyAwardRanks.en;
  const fallbackStories = legacyAwardStories[language] || legacyAwardStories.en;
  const genericRank = copy("award_rank_template", `Rank {level}`).replace("{level}", String(normalized));
  const genericStory = copy("award_story_template", fallbackStories[normalized - 1]).replace("{level}", String(normalized));
  const rankFromCopy = copy(`award_rank_${key}`, "");
  const storyFromCopy = copy(`award_story_${key}`, "");
  const isPlaceholder = (value: string) => /^award_(rank|story)_\d{2}$/i.test(value.trim()) || /^award_(rank|story)_\d+$/i.test(value.trim());
  const rank = rankFromCopy && rankFromCopy !== genericRank && !isPlaceholder(rankFromCopy) ? rankFromCopy : fallbackRanks[normalized - 1];
  const story = storyFromCopy && storyFromCopy !== genericStory && !isPlaceholder(storyFromCopy) ? storyFromCopy : fallbackStories[normalized - 1];
  return {
    level: normalized,
    rank: cleanAppText(rank),
    story: cleanAppText(story),
  };
}

function premiumPeriodLabel(plan: PremiumPlan | undefined, copy: (key: string, fallback: string) => string) {
  const source = `${plan?.product || ""} ${plan?.title || ""} ${plan?.days_label || ""}`.toLowerCase();
  if (/year|annual|365|год/.test(source)) return copy("premium_period_year", "Год");
  if (/month|30|31|месяц/.test(source)) return copy("premium_period_month", "Месяц");
  if (/year|annual|365|год/.test(source)) return copy("premium_period_year", "Год");
  if (/month|30|31|месяц/.test(source)) return copy("premium_period_month", "Месяц");
  return cleanAppText(plan?.days_label) || copy("premium_period", "Premium");
}

function premiumPlanSource(plan: PremiumPlan | undefined) {
  return `${plan?.product || ""} ${plan?.title || ""} ${plan?.tier || ""} ${plan?.days_label || ""}`.toLowerCase();
}

function premiumPlanIsFree(plan: PremiumPlan | undefined) {
  const source = premiumPlanSource(plan);
  return source.includes("free") || source.includes("basic") || plan?.product === "free";
}

function premiumPlanIsPlatinum(plan: PremiumPlan | undefined) {
  return premiumPlanSource(plan).includes("platinum");
}

function premiumPlanIsPremium(plan: PremiumPlan | undefined) {
  const source = premiumPlanSource(plan);
  return !premiumPlanIsPlatinum(plan) && (source.includes("premium") || source.includes("month") || source.includes("year"));
}

function premiumPlanTierKey(plan: PremiumPlan | undefined) {
  if ((plan?.current || plan?.active || plan?.is_current) && premiumPlanIsFree(plan)) return "free";
  if (premiumPlanIsPlatinum(plan)) return "platinum";
  if (premiumPlanIsPremium(plan)) return "premium";
  if (premiumPlanIsFree(plan)) return "free";
  return cleanAppText(plan?.tier || plan?.product).toLowerCase();
}

function currentUserPlanKey(user: UserProfile) {
  const source = `${user.plan || ""} ${user.premium_until || ""}`.toLowerCase();
  if (source.includes("platinum")) return "platinum";
  if (source.includes("premium")) return "premium";
  if (!user.premium) return "free";
  return "premium";
}

function premiumPlanIsCurrent(plan: PremiumPlan | undefined, user: UserProfile) {
  if (plan?.current || plan?.active || plan?.is_current) return true;
  return premiumPlanTierKey(plan) === currentUserPlanKey(user);
}

function premiumPlanTitle(plan: PremiumPlan, copy: (key: string, fallback: string) => string) {
  const source = premiumPlanSource(plan);
  const title = cleanAppText(plan.title);
  const isYear = source.includes("year") || source.includes("365");
  if (premiumPlanIsFree(plan)) return title || copy("free_plan_title", "Free");
  if (premiumPlanIsPlatinum(plan)) return title || copy(isYear ? "platinum_year_title" : "platinum_month_title", isYear ? "Platinum for a year" : "Platinum");
  if (premiumPlanIsPremium(plan)) return title || copy(isYear ? "premium_year_title" : "premium_month_title", isYear ? "Premium for a year" : "Premium");
  return cleanAppText(plan.title || plan.product || copy("premium", "Premium"));
}

function premiumPlanTier(plan: PremiumPlan, copy: (key: string, fallback: string) => string) {
  const tier = cleanAppText(plan.tier);
  if (tier && !["free", "basic", "premium", "platinum"].includes(tier.toLowerCase())) return tier;
  if (premiumPlanIsFree(plan)) return copy("free_plan_tier", "Free");
  if (premiumPlanIsPlatinum(plan)) return copy("platinum_month_tier", "Platinum");
  if (premiumPlanIsPremium(plan)) return copy("premium_month_tier", "Premium");
  return cleanAppText(plan.tier || plan.product || copy("premium", "Premium"));
}

function premiumPlanLabel(plan: PremiumPlan, copy: (key: string, fallback: string) => string) {
  const label = cleanAppText(plan.label);
  if (label) return label;
  if (premiumPlanIsFree(plan)) return copy("free_plan_label", "Для начала");
  if (premiumPlanIsPlatinum(plan)) return copy("platinum_plan_label", "Intensive");
  if (premiumPlanIsPremium(plan)) return copy("premium_plan_label", "Regular study");
  return premiumPlanTier(plan, copy);
}

function premiumPlanBody(plan: PremiumPlan, copy: (key: string, fallback: string) => string) {
  const source = premiumPlanSource(plan);
  const body = cleanAppText(plan.body);
  if (body) return body;
  if (premiumPlanIsFree(plan)) return copy("free_plan_body", "Basic text learning, word training, Phrasebook, and progress overview. AI Tutor, listening practice, pronunciation scoring, and voice review open in Premium.");
  if (premiumPlanIsPlatinum(plan)) return copy("platinum_month_body", "AI Tutor with maximum daily limits, deeper role-play practice, intensive review practice, and maximum voice and pronunciation practice.");
  if (source.includes("year") || source.includes("365")) return copy("premium_year_body", "Same daily AI audio limits, paid yearly.");
  if (premiumPlanIsPremium(plan)) return copy("premium_month_body", "The main mode for daily practice: AI Tutor, listening practice, pronunciation scoring, AI-guided tutor lessons, voice review, image tools, and expanded daily limits.");
  return cleanAppText(plan.days_label) || copy("premium_month_body", "AI Tutor, listening practice, pronunciation scoring, voice tools, image practice, and focused daily training.");
}

function premiumPlanCatalog(plans: PremiumPlan[], copy: (key: string, fallback: string) => string) {
  const fallback = fallbackPlans(copy);
  const source = plans.length ? plans : fallback;
  const catalog = source.some((plan) => premiumPlanIsFree(plan)) ? source : [fallback[0], ...source];
  const seen = new Set<string>();
  return catalog.filter((plan) => {
    const key = String(plan.product || plan.title || premiumPlanTier(plan, copy)).toLowerCase();
    if (seen.has(key)) return false;
    seen.add(key);
    return true;
  });
}

function premiumPlanTextList(value: unknown): string[] {
  if (!Array.isArray(value)) return [];
  return value.map((item) => cleanAppText(item)).filter(Boolean);
}

function premiumPlanFeatures(plan: PremiumPlan, copy: (key: string, fallback: string) => string): PremiumPlanFeature[] {
  const features = premiumPlanTextList(plan.features);
  const lockedFeatures = premiumPlanTextList(plan.locked_features);
  if (features.length || lockedFeatures.length) {
    return [
      ...features.map((text, index) => ({ text, highlight: index === 0 && !premiumPlanIsFree(plan) })),
      ...lockedFeatures.map((text) => ({ text, locked: true })),
    ];
  }
  if (premiumPlanIsFree(plan)) {
    return [
      { text: copy("free_feature_daily", "Daily habit, starter lessons, and basic word training") },
      { text: copy("free_feature_phrasebook", "Phrasebook and progress overview") },
      { text: copy("free_feature_phrase_audio", "Phrase audio is available") },
      { text: copy("free_feature_no_ai_tutor", "AI Tutor is locked until Premium"), locked: true },
      { text: copy("free_feature_no_audio", "Listening practice and pronunciation scoring need Premium"), locked: true },
    ];
  }
  if (premiumPlanIsPlatinum(plan)) {
    return [
      { text: copy("platinum_feature_ai_tutor", "Maximum daily limits"), highlight: true },
      { text: copy("platinum_feature_voice", "More voice or image-context practice") },
      { text: copy("platinum_feature_priority", "Best mode for dense daily study") },
    ];
  }
  return [
    { text: copy("premium_feature_ai_tutor", "AI-guided tutor lessons included"), highlight: true },
    { text: copy("premium_feature_voice", "Listening practice, pronunciation scoring, image, and translator tools") },
    { text: copy("premium_feature_limits", "Expanded daily limits for steady learning") },
  ];
}

function paymentAmountLine(info: ApiRecord) {
  const currency = asText(info.currency || info.asset || info.network, "");
  const labeled = asText(info.amount_label || info.price_label, "");
  if (labeled) return labeled;
  const raw = asText(info.amount || info.price || info.amount_units, "");
  const numeric = Number(String(raw).replace(",", "."));
  let amount = raw;
  if (Number.isFinite(numeric) && numeric > 0) {
    amount = new Intl.NumberFormat("ru-RU", { maximumFractionDigits: 6 }).format(numeric);
  }
  return amount && currency && !amount.toLowerCase().includes(currency.toLowerCase()) ? `${amount} ${currency}` : amount;
}

function isConfirmedPaymentStatus(status: string) {
  return /paid|confirmed|success|succeed|complete|completed/i.test(cleanAppText(status));
}

function readPaymentHistory(key: string): PaymentHistoryItem[] {
  try {
    const value = localStorage.getItem(key);
    if (!value) return [];
    const parsed = JSON.parse(value);
    return Array.isArray(parsed) ? parsed.filter((item) => isConfirmedPaymentStatus(asText(item?.status, ""))).slice(0, 20) : [];
  } catch {
    return [];
  }
}

function readPhrasebook(key: string): PhrasebookItem[] {
  try {
    const value = localStorage.getItem(key);
    if (!value) return [];
    const parsed = JSON.parse(value);
    return normalizePhrasebookItems(parsed);
  } catch {
    return [];
  }
}

function normalizePhrasebookItems(value: unknown): PhrasebookItem[] {
  if (!Array.isArray(value)) return [];
  const seen = new Set<string>();
  const items: PhrasebookItem[] = [];
  for (const raw of value) {
    if (!raw || typeof raw !== "object") continue;
    const record = raw as Record<string, unknown>;
    const phrase = cleanAppText(asText(record.phrase, "")).trim();
    if (!phrase) continue;
    const key = phrase.toLowerCase();
    if (seen.has(key)) continue;
    seen.add(key);
    items.push({
      id: cleanAppText(asText(record.id, "")) || `${Date.now()}-${items.length}`,
      phrase,
      translation: cleanAppText(asText(record.translation, "")),
      note: cleanAppText(asText(record.note, "")),
      source: (["lesson", "practice", "roleplay", "mistake", "manual", "word", "tool"].includes(asText(record.source, "")) ? asText(record.source, "") : "manual") as PhrasebookSource,
      language: cleanAppText(asText(record.language, "")),
      createdAt: cleanAppText(asText(record.createdAt || record.created_at, "")) || new Date().toISOString(),
    });
  }
  return items.slice(0, 80);
}

function readHabitLog(key: string): Record<string, HabitDay> {
  try {
    const value = localStorage.getItem(key);
    if (!value) return {};
    const parsed = JSON.parse(value);
    return parsed && typeof parsed === "object" && !Array.isArray(parsed) ? parsed as Record<string, HabitDay> : {};
  } catch {
    return {};
  }
}

function normalizeHabitLogValue(value: unknown): Record<string, HabitDay> {
  if (!value || typeof value !== "object" || Array.isArray(value)) return {};
  const result: Record<string, HabitDay> = {};
  Object.entries(value as Record<string, unknown>).forEach(([date, raw]) => {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(date) || !raw || typeof raw !== "object" || Array.isArray(raw)) return;
    const record = raw as Record<string, unknown>;
    result[date] = {
      login: Boolean(record.login),
      complete: Boolean(record.complete),
      claimed: Boolean(record.claimed),
      claimedAt: asText(record.claimedAt || record.claimed_at || "", "") || undefined,
      lessons: Number(record.lessons || 0),
      practice: Number(record.practice || 0),
      voice: Number(record.voice || 0),
    };
  });
  return result;
}

function localDateKey(date = new Date()) {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`;
}

function monthDateKeys(offset = 0) {
  const now = new Date();
  const start = new Date(now.getFullYear(), now.getMonth() + offset, 1);
  const end = new Date(now.getFullYear(), now.getMonth() + offset + 1, 0);
  const keys: string[] = [];
  for (let day = 1; day <= end.getDate(); day += 1) {
    keys.push(localDateKey(new Date(start.getFullYear(), start.getMonth(), day)));
  }
  return keys;
}

function recentDateKeys(days: number, end = new Date()) {
  const keys: string[] = [];
  const cursor = new Date(end.getFullYear(), end.getMonth(), end.getDate());
  for (let index = 0; index < days; index += 1) {
    keys.push(localDateKey(cursor));
    cursor.setDate(cursor.getDate() - 1);
  }
  return keys;
}

function habitWeekTotal(log: Record<string, HabitDay> | undefined, field: "lessons" | "practice" | "voice") {
  return recentDateKeys(7).reduce((total, key) => total + Number(log?.[key]?.[field] || 0), 0);
}

function weeklyGoalFromDailyLimit(limit: unknown, fallback: number) {
  const daily = Number(limit || 0);
  return Math.max(1, Math.round((daily > 0 ? daily : fallback) * 7));
}

function dailyGoalComplete(user: UserProfile) {
  const voiceLimit = Number(user.voice_limit || 0);
  const voiceReady = voiceLimit > 0 ? Number(user.voice_today || 0) >= 1 : true;
  return Number(user.lessons_today || 0) >= 1 && Number(user.practice_today || 0) >= 1 && voiceReady;
}

function dailyBonusLastClaimedAt(user: UserProfile) {
  return asText(user.daily_bonus_last_claimed_at || "", "");
}

function userClaimedDailyBonusToday(user: UserProfile, todayKey = localDateKey()) {
  return Array.isArray(user.daily_bonus_claims) && user.daily_bonus_claims.includes(todayKey);
}

function dailyBonusLocked(day: HabitDay | undefined, user: UserProfile) {
  if (day?.claimed || userClaimedDailyBonusToday(user)) return true;
  const claimedAt = day?.claimedAt || dailyBonusLastClaimedAt(user);
  if (!claimedAt) return false;
  const timestamp = Date.parse(claimedAt);
  if (!Number.isFinite(timestamp)) return false;
  return Date.now() - timestamp < 24 * 60 * 60 * 1000;
}

function updateHabitForToday(current: Record<string, HabitDay>, user: UserProfile) {
  const todayKey = localDateKey();
  const previous = current[todayKey] || {};
  const claimedToday = previous.claimed || userClaimedDailyBonusToday(user, todayKey);
  const claimedAt = previous.claimedAt || dailyBonusLastClaimedAt(user) || undefined;
  return {
    ...current,
    [todayKey]: {
      ...previous,
      login: true,
      complete: previous.complete || dailyGoalComplete(user),
      claimed: claimedToday,
      claimedAt,
      lessons: Number(user.lessons_today || 0),
      practice: Number(user.practice_today || 0),
      voice: Number(user.voice_today || 0),
    },
  };
}

function habitStreak(log: Record<string, HabitDay>) {
  let streak = 0;
  const cursor = new Date();
  for (let index = 0; index < 31; index += 1) {
    const key = localDateKey(cursor);
    if (!log[key]?.complete) break;
    streak += 1;
    cursor.setDate(cursor.getDate() - 1);
  }
  return streak;
}

function formatStreakLabel(days: number, user: UserProfile, copy: (key: string, fallback: string) => string) {
  const value = Math.max(0, Math.round(Number(days || 0)));
  const lang = languageCode(user.interface_language);
  if (lang === "ru" || lang === "uk") {
    const mod10 = value % 10;
    const mod100 = value % 100;
    const key = mod10 === 1 && mod100 !== 11 ? "days_streak_one" : mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14) ? "days_streak_few" : "days_streak_many";
    const fallback = key === "days_streak_one" ? "день подряд" : key === "days_streak_few" ? "дня подряд" : "дней подряд";
    return `${value} ${copy(key, fallback)}`;
  }
  const key = value === 1 ? "days_streak_one" : "days_streak_many";
  return `${value} ${copy(key, "day streak")}`;
}

function phraseCandidatesFromMessages(messages: ChatMessage[]) {
  const candidates: Array<{ phrase: string; source: PhrasebookSource; details?: ApiRecord }> = [];
  const pushCandidate = (phrase: string, source: PhrasebookSource, details?: ApiRecord) => {
    const cleaned = cleanAppText(phrase).replace(/\s+/g, " ").trim();
    if (isUsefulPhraseCandidate(cleaned) && cleaned.length <= 220) candidates.push({ phrase: cleaned, source, details });
  };
  for (const message of messages) {
    if (message.role === "user") continue;
    const details = getRecord(message.details);
    const source = message.meta === "roleplay" ? "roleplay" : message.meta === "practice" ? "practice" : message.meta?.startsWith("tools:") ? "tool" : "lesson";
    const fields = ["correction", "model_phrase", "correction_audio_text", "corrected_text", "example"];
    for (const field of fields) {
      pushCandidate(recordField(details, [field]), source, details);
    }
    recordList(details.answer_variants)
      .concat(recordList(details.suggestions), recordList(details.examples), recordList(details.phrases))
      .forEach((phrase) => pushCandidate(phrase, source, details));
    recordMistakes(details.mistakes).forEach((mistake) => {
      pushCandidate(recordField(mistake, ["correction", "example"]), "lesson", {
        ...details,
        note: recordField(mistake, ["word", "explanation", "context"]),
      });
    });
  }
  const seen = new Set<string>();
  return candidates.filter((item) => {
    const key = item.phrase.toLowerCase();
    if (seen.has(key)) return false;
    seen.add(key);
    return true;
  }).slice(0, 6);
}

function phrasebookKey(value: string) {
  return cleanAppText(value).replace(/\s+/g, " ").trim().toLowerCase();
}

function samePhrasebookText(left: string | undefined, right: string | undefined) {
  const a = phrasebookKey(left || "");
  const b = phrasebookKey(right || "");
  return Boolean(a && b && a === b);
}

function isUsefulPhraseCandidate(value: string) {
  const phrase = cleanAppText(value).replace(/\s+/g, " ").trim();
  if (phrase.length < 4 || phrase.length > 220) return false;
  if (!/[\p{L}\p{N}]/u.test(phrase)) return false;
  if (/^[\p{P}\p{S}\s]+$/u.test(phrase)) return false;
  if (/^[?.!,;:]+$/.test(phrase)) return false;
  return true;
}

function mistakeCategory(item: MistakeItem): MistakeCategory {
  const source = `${item.word || ""} ${item.correction || ""} ${item.context || ""} ${item.explanation || ""}`.toLowerCase();
  if (/(please|polite|rude|would|could|thank|tone|request|order|вежлив|просьб|заказ|тон|пожалуйста|спасибо)/i.test(source)) return "politeness";
  if (/(word order|order|position|место слова|порядок|позици|перестав|поставь|поставил)/i.test(source)) return "word-order";
  if (/(spelling|typo|letter|orthograph|правопис|орфограф|букв|опечат|написан)/i.test(source)) return "spelling";
  if (/(vocab|lexic|word choice|meaning|словар|лексик|значени|не то слово|выбор слова)/i.test(source)) return "vocabulary";
  if (/(please|polite|rude|would|could|thank|вежлив|просьб|заказ|тон)/i.test(source)) return "politeness";
  if (/(word order|order|position|порядок|позици|перестав|место слова)/i.test(source)) return "word-order";
  if (/(spelling|typo|letter|orthograph|правопис|орфограф|букв|опечат)/i.test(source)) return "spelling";
  if (/(vocab|lexic|word choice|meaning|словар|лексик|значени|не то слово)/i.test(source)) return "vocabulary";
  return "grammar";
}

function mistakeCategoryLabel(category: MistakeCategory, copy: (key: string, fallback: string) => string) {
  const labels: Record<MistakeCategory, string> = {
    grammar: copy("mistake_group_grammar", "Grammar"),
    "word-order": copy("mistake_group_word_order", "Word order"),
    vocabulary: copy("mistake_group_vocabulary", "Vocabulary"),
    politeness: copy("mistake_group_politeness", "Politeness"),
    spelling: copy("mistake_group_spelling", "Spelling"),
  };
  return labels[category];
}

export function App() {
  const [session, setSession] = useState<SessionData | null>(null);
  const [sessionLoading, setSessionLoading] = useState(true);
  const [activeView, setActiveView] = useState<ViewId>(() => normalizeView(new URLSearchParams(location.search).get("view")));
  const [notFoundPath, setNotFoundPath] = useState(() => currentNotFoundPath());
  const [theme, setTheme] = useState<Theme>(() => (localStorage.getItem("poliglot-theme") === "dark" ? "dark" : "light"));
  const [brightness, setBrightness] = useState(() => {
    const saved = Number(localStorage.getItem("poliglot-brightness") || 100);
    return Number.isFinite(saved) ? Math.min(125, Math.max(65, saved)) : 100;
  });
  const [busy, setBusy] = useState<string | null>(null);
  const [status, setStatus] = useState<{ kind: StatusKind; text: string } | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([
    panelMessage(appCopy("ru", "ready_body"), "default", appCopy("ru", "ready_title")),
  ]);
  const [activeLessonTaskId, setActiveLessonTaskId] = useState("");
  const [draft, setDraft] = useState("");
  const [voiceFile, setVoiceFile] = useState<File | null>(null);
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [wordChallenge, setWordChallenge] = useState<WordChallenge | null>(null);
  const [wordResult, setWordResult] = useState<TrainerResult | null>(null);
  const [wordGameResult, setWordGameResult] = useState<TrainerResult | null>(null);
  const [aiTutorSessionId, setAiTutorSessionId] = useState("");
  const [aiTutorStep, setAiTutorStep] = useState<AiTutorStep | null>(null);
  const [aiTutorFeedback, setAiTutorFeedback] = useState<AiTutorFeedbackState | null>(null);
  const [tutorLoadError, setTutorLoadError] = useState("");
  const [spellingChallenge, setSpellingChallenge] = useState<SpellingChallenge | null>(null);
  const [spellingResult, setSpellingResult] = useState<TrainerResult | null>(null);
  const [levelQuestion, setLevelQuestion] = useState<LevelQuestion | null>(null);
  const [levelResult, setLevelResult] = useState<TrainerResult | null>(null);
  const [vocabulary, setVocabulary] = useState<VocabularyItem[]>([]);
  const [vocabularyMeta, setVocabularyMeta] = useState<VocabularyPage>({ page: 0, page_size: 10, total: 0, total_pages: 1 });
  const [mistakes, setMistakes] = useState<MistakeItem[]>([]);
  const [mistakePractice, setMistakePractice] = useState<MistakeItem | null>(null);
  const [mistakePracticeIndex, setMistakePracticeIndex] = useState<number | null>(null);
  const [mistakeAnswer, setMistakeAnswer] = useState("");
  const [mistakeResult, setMistakeResult] = useState<TrainerResult | null>(null);
  const [shadowingTarget, setShadowingTarget] = useState("");
  const [pronunciationTarget, setPronunciationTarget] = useState("");
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [leaderboardMeta, setLeaderboardMeta] = useState<LeaderboardMeta>({ language: "", languageName: "" });
  const [leaderboardLanguage, setLeaderboardLanguage] = useState("global");
  const [premiumPlans, setPremiumPlans] = useState<PremiumPlan[]>([]);
  const [payment, setPayment] = useState<PaymentState | null>(null);
  const [paymentNotice, setPaymentNotice] = useState<PaymentNotice | null>(null);
  const [paymentHistory, setPaymentHistory] = useState<PaymentHistoryItem[]>([]);
  const [phrasebook, setPhrasebook] = useState<PhrasebookItem[]>([]);
  const [habitLog, setHabitLog] = useState<Record<string, HabitDay>>({});
  const [dailyBonus, setDailyBonus] = useState<DailyBonusNotice | null>(null);
  const [xpGain, setXPGain] = useState<XPGainNotice | null>(null);
  const [bugReportOpen, setBugReportOpen] = useState(false);
  const [guideOpen, setGuideOpen] = useState(false);
  const [onboardingOpen, setOnboardingOpen] = useState(false);
  const [selectedAwardLevel, setSelectedAwardLevel] = useState<number | null>(null);
  const [activationKey, setActivationKey] = useState("");
  const [toolMode, setToolMode] = useState<ToolMode>("translator");
  const [toolSourceLanguage, setToolSourceLanguage] = useState("auto");
  const [toolTargetLanguage, setToolTargetLanguage] = useState("ru");
  const [toolVoiceFile, setToolVoiceFile] = useState<File | null>(null);
  const [toolImageFile, setToolImageFile] = useState<File | null>(null);
  const [toolResult, setToolResult] = useState<TranslatorResult | null>(null);
  const [roleplayResult, setRoleplayResult] = useState<ChatMessage | null>(null);
  const bootstrappedRef = useRef(false);
  const lastPremiumUntilRef = useRef("");
  const viewHistoryReadyRef = useRef(false);
  const suppressHistoryPushRef = useRef(false);
  const autoDailyBonusClaimRef = useRef("");
  const xpGainTimerRef = useRef<number | null>(null);

  useEffect(() => {
    document.documentElement.dataset.theme = theme;
    localStorage.setItem("poliglot-theme", theme);
  }, [theme]);

  useEffect(() => {
    localStorage.setItem("poliglot-brightness", String(brightness));
  }, [brightness]);

  useEffect(() => {
    if (!("serviceWorker" in navigator)) return;
    let refreshing = false;
    const hadController = !!navigator.serviceWorker.controller;
    const reloadWhenControlled = () => {
      if (!hadController || refreshing) return;
      refreshing = true;
      window.location.reload();
    };
    navigator.serviceWorker.addEventListener("controllerchange", reloadWhenControlled);
      navigator.serviceWorker
        .register("/app/offline-deck-sw.js", { updateViaCache: "none" })
      .then((registration) => {
        registration.update().catch(() => undefined);
        if (registration.waiting) {
          registration.waiting.postMessage({ type: "SKIP_WAITING" });
        }
        registration.addEventListener("updatefound", () => {
          const worker = registration.installing;
          if (!worker) return;
          worker.addEventListener("statechange", () => {
            if (worker.state === "installed" && navigator.serviceWorker.controller) {
              worker.postMessage({ type: "SKIP_WAITING" });
            }
          });
        });
      })
      .catch(() => undefined);
    return () => {
      navigator.serviceWorker.removeEventListener("controllerchange", reloadWhenControlled);
    };
  }, []);

  useEffect(() => {
    let cancelled = false;
    loadSession()
      .then((data) => {
        if (cancelled) return;
        setSession(data);
        const activeLessonPrompt = cleanAppText(data.user?.active_lesson_prompt).trim();
        if (activeLessonPrompt) {
          const taskId = "session-active-lesson";
          const instruction = cleanAppText(data.user?.active_lesson_instruction).trim()
            || appCopy(data.user?.interface_language || "ru", "new_lesson_title", "New lesson");
          setActiveLessonTaskId(taskId);
          setMessages((current) => [
            panelMessage(activeLessonPrompt, "default", instruction, { lesson: activeLessonPrompt, instruction, task_id: taskId }, "lesson"),
            ...current.filter((message) => message.meta !== "lesson"),
          ]);
        }
      })
      .catch((error) => {
        if (!cancelled) setStatus({ kind: "error", text: error instanceof Error ? error.message : appCopy("ru", "session_failed") });
      })
      .finally(() => {
        if (!cancelled) setSessionLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    if (isAuthStandaloneRoute() || notFoundPath) return;
    const url = new URL(location.href);
    url.searchParams.set("view", activeView);
    const state = { poliglotApp: true, view: activeView };
    if (!viewHistoryReadyRef.current) {
      history.replaceState(state, "", url);
      viewHistoryReadyRef.current = true;
      return;
    }
    if (suppressHistoryPushRef.current) {
      suppressHistoryPushRef.current = false;
      history.replaceState(state, "", url);
      return;
    }
    history.pushState(state, "", url);
  }, [activeView, notFoundPath]);

  useEffect(() => {
    const handlePopState = () => {
      const missing = currentNotFoundPath();
      if (missing) {
        setNotFoundPath(missing);
        return;
      }
      suppressHistoryPushRef.current = true;
      setNotFoundPath("");
      setActiveView(normalizeView(new URLSearchParams(location.search).get("view")));
    };
    window.addEventListener("popstate", handlePopState);
    return () => window.removeEventListener("popstate", handlePopState);
  }, []);

  const user = session?.user || {};

  useEffect(() => {
    if (user.interface_language && toolTargetLanguage === "ru") {
      setToolTargetLanguage(user.interface_language);
    }
  }, [toolTargetLanguage, user.interface_language]);

  const copy = (key: string, fallback: string) => {
    const clientValue = appCopy(user.interface_language, key, fallback);
    if (clientValue && clientValue !== key && !isGenericSectionCopy(clientValue)) return clientValue;
    const serverValue =
      cleanAppText(session?.copy?.[key]) ||
      cleanAppText(session?.tool_copy?.[key]) ||
      cleanAppText(session?.premium_copy?.[key]) ||
      cleanAppText(session?.system_copy?.[key]) ||
      "";
    if (serverValue && !isGenericSectionCopy(serverValue)) return serverValue;
    return (
      cleanAppText(pronunciationFallbackForKey(user.interface_language, key)) ||
      cleanAppText(roleplayFallbackForKey(user.interface_language, key)) ||
      clientValue ||
      cleanAppText(fallback)
    );
  };

  const accountKey = asText(session?.account?.login || user.telegram_account?.id || "guest", "guest");
  const onboardingKey = `poliglot-onboarding-v2:${accountKey}`;
  const paymentHistoryKey = `poliglot-payment-history-v2:${accountKey}`;
  const habitKey = `poliglot-habit-v2:${accountKey}`;

  const showXPGain = (xp: number, total?: number, title?: string) => {
    const amount = Math.max(0, Math.round(Number(xp || 0)));
    if (!amount) return;
    if (xpGainTimerRef.current) window.clearTimeout(xpGainTimerRef.current);
    setXPGain({ xp: amount, total: Number.isFinite(total) ? Math.round(Number(total)) : undefined, title });
    xpGainTimerRef.current = window.setTimeout(() => setXPGain(null), 4000);
  };

  useEffect(() => {
    return () => {
      if (xpGainTimerRef.current) window.clearTimeout(xpGainTimerRef.current);
    };
  }, []);

  useEffect(() => {
    if (!session?.authenticated) return;
    setPaymentHistory(readPaymentHistory(paymentHistoryKey));
    const serverPhrasebook = normalizePhrasebookItems(user.phrasebook);
    setPhrasebook(serverPhrasebook);
    const serverHabitLog = normalizeHabitLogValue(user.habit_log);
    const localHabitLog = readHabitLog(habitKey);
    const mergedHabitLog = { ...localHabitLog, ...serverHabitLog };
    setHabitLog(mergedHabitLog);
    localStorage.setItem(habitKey, JSON.stringify(mergedHabitLog));
    if (!localStorage.getItem(onboardingKey)) setOnboardingOpen(true);
  }, [session?.authenticated, paymentHistoryKey, habitKey, onboardingKey, user.phrasebook, user.habit_log]);

  useEffect(() => {
    if (!session?.authenticated) return;
    setHabitLog((current) => {
      const next = updateHabitForToday(current, user);
      localStorage.setItem(habitKey, JSON.stringify(next));
      return next;
    });
  }, [session?.authenticated, habitKey, user.lessons_today, user.practice_today, user.voice_today]);

  const rememberPaymentHistory = useCallback(
    (entry: PaymentHistoryItem) => {
      if (!isConfirmedPaymentStatus(entry.status)) return;
      setPaymentHistory((current) => {
        const next = [entry, ...current.filter((item) => item.id !== entry.id)].slice(0, 20);
        localStorage.setItem(paymentHistoryKey, JSON.stringify(next));
        return next;
      });
    },
    [paymentHistoryKey],
  );

  const savePhrase = useCallback(
    (phrase: string, source: PhrasebookSource, details?: ApiRecord) => {
      const cleaned = cleanAppText(phrase).trim();
      if (!cleaned) return;
      const translation = cleanAppText(recordField(details || {}, ["translation", "prompt", "message"]));
      const rawNote = cleanAppText(recordField(details || {}, ["note", "context"]));
      const note = samePhrasebookText(rawNote, translation) || samePhrasebookText(rawNote, cleaned) ? "" : rawNote;
      const item: PhrasebookItem = {
        id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
        phrase: cleaned,
        translation: samePhrasebookText(translation, cleaned) ? "" : translation,
        note,
        source,
        language: user.learning_language,
        createdAt: new Date().toISOString(),
      };
      setPhrasebook((current) => {
        const next = normalizePhrasebookItems([item, ...current.filter((entry) => phrasebookKey(entry.phrase) !== phrasebookKey(cleaned))]);
        setStatus({ kind: "ok", text: copy("phrase_saved", "Фраза добавлена в phrasebook") });
        return next;
      });
      void api<ApiRecord>("/api/phrasebook", { method: "POST", body: item })
        .then((payload) => {
          const next = normalizePhrasebookItems(payload.items);
          setPhrasebook(next);
        })
        .catch(() => undefined);
    },
    [copy, user.learning_language],
  );

  const savedPhraseKeys = useMemo(() => new Set(phrasebook.map((item) => phrasebookKey(item.phrase)).filter(Boolean)), [phrasebook]);
  const isPhraseSaved = useCallback((phrase: string) => savedPhraseKeys.has(phrasebookKey(phrase)), [savedPhraseKeys]);

  const removePhrase = useCallback(
    (id: string) => {
      setPhrasebook((current) => {
        const next = current.filter((item) => item.id !== id);
        return next;
      });
      void api<ApiRecord>(`/api/phrasebook?id=${encodeURIComponent(id)}`, { method: "DELETE" })
        .then((payload) => {
          const next = normalizePhrasebookItems(payload.items);
          setPhrasebook(next);
        })
        .catch(() => undefined);
    },
    [],
  );

  const claimDailyBonus = async () => {
    const todayKey = localDateKey();
    const current = habitLog[todayKey] || updateHabitForToday({}, user)[todayKey];
    if (!current?.complete || dailyBonusLocked(current, user)) return;
    const payload = await runAction("daily-bonus", () => api<ApiRecord>("/api/daily/claim", { method: "POST", body: { date: todayKey } }));
    if (!payload) return;
    const record = getRecord(payload);
    const claimedAt = asText(record.claimed_at || record.claimedAt || new Date().toISOString(), new Date().toISOString());
    if (record.claimed === false) {
      setStatus({ kind: "info", text: copy("daily_bonus_already_claimed", "Daily bonus is available once every 24 hours.") });
      const payloadUser = getPayloadUser(payload);
      if (payloadUser) setSession((currentSession) => updateSessionUser(currentSession, payloadUser));
      return;
    }
    const xp = Number(record.xp || 25);
    const streak = Number(record.streak || habitStreak(habitLog));
    setDailyBonus({ xp, streak });
    setTimeout(() => setDailyBonus(null), 2800);
    setHabitLog((currentLog) => {
      const next = { ...currentLog, [todayKey]: { ...(currentLog[todayKey] || {}), claimed: true, claimedAt, complete: true, login: true } };
      localStorage.setItem(habitKey, JSON.stringify(next));
      return next;
    });
    const payloadUser = getPayloadUser(payload);
    if (payloadUser) {
      setSession((currentSession) => updateSessionUser(currentSession, {
        ...payloadUser,
        interface_language: user.interface_language || payloadUser.interface_language,
        learning_language: user.learning_language || payloadUser.learning_language,
      }));
    }
  };

  useEffect(() => {
    const currentUntil = asText(user.premium_until, "");
    if (!lastPremiumUntilRef.current) {
      lastPremiumUntilRef.current = currentUntil;
      return;
    }
    if (currentUntil && currentUntil !== lastPremiumUntilRef.current && user.premium) {
      setPaymentNotice({
        kind: "success",
        title: copy("payment_success_title", "Оплата подтверждена"),
        description: copy("payment_success_body", "Подписка обновлена. Premium уже активен в вашем аккаунте."),
        planTitle: user.plan || copy("premium", "Premium"),
        period: user.plan || copy("premium", "Premium"),
        premiumUntil: prettyDate(currentUntil),
      });
      rememberPaymentHistory({
        id: `premium-${currentUntil}`,
        date: new Date().toISOString(),
        plan: user.plan || copy("premium", "Premium"),
        period: user.plan || copy("premium", "Premium"),
        amount: "",
        method: copy("webhook_or_auto_check", "Webhook / auto check"),
        status: copy("paid", "Paid"),
      });
    }
    lastPremiumUntilRef.current = currentUntil;
  }, [user.premium, user.premium_until, user.plan, rememberPaymentHistory]);

  const nav = useMemo(
    () =>
      navItems
        .filter((item) => !hiddenRibbonViews.has(item.id))
        .map((item) => ({
          ...item,
          label: copy(navCopyKey(item.id), item.label),
          shortLabel: copy(navCopyKey(item.id), item.shortLabel || item.label),
          description: appNavDescription(item.id, copy, item.description),
        })),
    [session],
  );

  const activeItem = nav.find((item) => item.id === activeView) || nav[0];
  const details = translateViewDetails(activeView, copy);

  const setView = (view: ViewId) => {
    setNotFoundPath("");
    setActiveView(normalizeView(view));
  };

  const refreshSession = async () => {
    const next = await loadSession();
    setSession(next);
    return next;
  };

  const absorbPayload = (payload: unknown, options?: { preserveLanguage?: boolean; suppressXPGain?: boolean; rewardTitle?: string }) => {
    const payloadUser = getPayloadUser(payload);
    if (!payloadUser) return;
    const nextUser = options?.preserveLanguage
      ? {
          ...payloadUser,
          interface_language: user.interface_language || payloadUser.interface_language,
          learning_language: user.learning_language || payloadUser.learning_language,
        }
      : payloadUser;
    const previousXP = Number(user.xp || 0);
    const nextXP = Number(nextUser.xp || previousXP);
    const record = getRecord(payload);
    const explicitRewardXP = Number(recordField(record, ["xp", "reward_xp", "awarded_xp"]));
    const hasExplicitReward = Number.isFinite(explicitRewardXP) && explicitRewardXP > 0;
    if (!options?.suppressXPGain && ((Number.isFinite(nextXP) && nextXP > previousXP) || hasExplicitReward)) {
      const rewardXP = hasExplicitReward ? explicitRewardXP : nextXP - previousXP;
      showXPGain(rewardXP, nextXP, options?.rewardTitle || recordField(record, ["reward_title", "source_title"]));
    }
    setSession((current) => updateSessionUser(current, nextUser));
  };

  const runAction = async (label: string, action: () => Promise<unknown>, success?: string) => {
    setBusy(label);
    setStatus(null);
    try {
      const payload = await action();
      absorbPayload(payload, {
        preserveLanguage: label === "daily-bonus",
        suppressXPGain: label === "daily-bonus",
        rewardTitle: label === "tutor" ? copy("ai_tutor", "AI Tutor") : undefined,
      });
      if (success) setStatus({ kind: "ok", text: success });
      return payload;
    } catch (error) {
      const text = localizedAPIError(error, copy);
      if (label === "tutor") setTutorLoadError(text);
      setStatus({ kind: "error", text });
      setMessages((current) => [panelMessage(text, "danger", copy("request_failed", "Request failed")), ...current].slice(0, 12));
      return null;
    } finally {
      setBusy(null);
    }
  };

  const startLesson = async () => {
    const payload = await runAction("lesson", () => api<ApiRecord>("/api/lesson/start", { method: "POST", body: {} }), copy("lesson_ready", "Lesson is ready."));
    if (!payload) return;
    const record = getRecord(payload);
    const body = recordField(record, ["lesson", "prompt", "question", "task", "message"]) || copy("lesson_generated", "Lesson generated.");
    const taskId = `${Date.now()}-${Math.random().toString(36).slice(2)}`;
    setActiveLessonTaskId(taskId);
    setMessages((current) => [
      panelMessage(body, "default", asText(record.instruction, copy("new_lesson_title", "New lesson")), { ...record, task_id: taskId }, "lesson"),
      ...current.filter((message) => message.meta !== "lesson"),
    ]);
    setDraft("");
    setVoiceFile(null);
    setImageFile(null);
    setView("lesson");
  };

  const startTutor = async () => {
    if (!user.premium) {
      const text = copy("tutor_premium_required", "AI Tutor is available with Premium.");
      setAiTutorSessionId("");
      setAiTutorStep(null);
      setAiTutorFeedback(null);
      setTutorLoadError(text);
      setStatus({ kind: "info", text });
      setView("premium");
      return;
    }
    setAiTutorSessionId("");
    setAiTutorStep(null);
    setAiTutorFeedback(null);
    setTutorLoadError("");
    const payload = await runAction("tutor", async () => {
      const controller = new AbortController();
      const timeout = window.setTimeout(() => controller.abort(), 30000);
      try {
        return await api<AiTutorResponse>("/api/ai-tutor/start", { method: "POST", body: {}, signal: controller.signal });
      } catch (error) {
        if (error instanceof DOMException && error.name === "AbortError") {
          throw new Error(copy("tutor_timeout", "Tutor lesson is taking too long. Try again."));
        }
        throw error;
      } finally {
        window.clearTimeout(timeout);
      }
    }, copy("tutor_ready", "Tutor lesson is ready."));
    if (!payload) return;
    const response = payload as AiTutorResponse;
    const step = response.next_step || null;
    const sessionID = String(response.session?.id || response.session?.ID || "");
    const lessonPayload = response.lesson || step?.lesson;
    if (!step || !sessionID) {
      const text = copy("request_failed", "Request failed");
      setTutorLoadError(text);
      setStatus({ kind: "error", text });
      return;
    }
    setTutorLoadError("");
    setAiTutorSessionId(sessionID);
    setAiTutorStep(step);
    setAiTutorFeedback(response.feedback || null);
    setMessages((current) => [
      panelMessage(
        aiTutorHistoryBody(lessonPayload, step),
        "default",
        cleanAppText(lessonPayload?.title || step.title || copy("ai_tutor", "AI Tutor")),
        getRecord(payload),
        "tutor",
      ),
      ...current,
    ].slice(0, 12));
    setView("tutor");
  };

  const requirePremiumFeature = (messageKey: string, fallback: string) => {
    if (user.premium) return true;
    const text = copy(messageKey, fallback);
    setStatus({ kind: "info", text });
    setView("premium");
    void loadPremiumPlans();
    return false;
  };

  const restartTutorLesson = async (lessonId: string) => {
    if (!user.premium) {
      const text = copy("tutor_premium_required", "AI Tutor is available with Premium.");
      setTutorLoadError(text);
      setStatus({ kind: "info", text });
      setView("premium");
      return false;
    }
    const cleanLessonId = cleanAppText(lessonId).trim();
    if (!cleanLessonId) return false;
    setTutorLoadError("");
    const payload = await runAction(
      "tutor",
      () => api<AiTutorResponse>("/api/ai-tutor/restart", { method: "POST", body: { lesson_id: cleanLessonId } }),
      copy("tutor_ready", "Tutor lesson is ready."),
    );
    if (!payload) return false;
    const response = payload as AiTutorResponse;
    const step = response.next_step || null;
    const sessionID = String(response.session?.id || response.session?.ID || "");
    const lessonPayload = response.lesson || step?.lesson;
    if (!step || !sessionID) {
      const text = copy("request_failed", "Request failed");
      setTutorLoadError(text);
      setStatus({ kind: "error", text });
      return false;
    }
    setAiTutorSessionId(sessionID);
    setAiTutorStep(step);
    setAiTutorFeedback(response.feedback || null);
    setMessages((current) => [
      panelMessage(
        aiTutorHistoryBody(lessonPayload, step),
        "default",
        cleanAppText(lessonPayload?.title || step.title || copy("ai_tutor", "AI Tutor")),
        getRecord(payload),
        "tutor",
      ),
      ...current,
    ].slice(0, 12));
    setView("tutor");
    return true;
  };

  const submitAiTutorStep = async (text: string, choice = "") => {
    if (!user.premium) {
      const notice = copy("tutor_premium_required", "AI Tutor is available with Premium.");
      setTutorLoadError(notice);
      setStatus({ kind: "info", text: notice });
      setView("premium");
      return;
    }
    if (!aiTutorSessionId) return;
    const payload = await runAction(
      "tutor",
      () => api<AiTutorResponse>("/api/ai-tutor/answer", { method: "POST", body: { session_id: aiTutorSessionId, text, choice } }),
      copy("tutor_answer_saved", "Answer saved"),
    );
    if (!payload) return;
    const response = payload as AiTutorResponse;
    const sessionID = String(response.session?.id || response.session?.ID || aiTutorSessionId);
    if (sessionID) setAiTutorSessionId(sessionID);
    setAiTutorStep(response.next_step || null);
    setAiTutorFeedback(response.feedback || null);
  };

  const reportAiTutorWord = async (input: { stage: string; proposedWord: string; proposedTranslation: string; comment: string }) => {
    if (!user.premium) {
      const notice = copy("tutor_premium_required", "AI Tutor is available with Premium.");
      setTutorLoadError(notice);
      setStatus({ kind: "info", text: notice });
      setView("premium");
      return false;
    }
    if (!aiTutorSessionId || !input.stage) return false;
    const payload = await runAction(
      "tutor-word-report",
      () => api<AiTutorResponse>("/api/ai-tutor/word-report", {
        method: "POST",
        body: {
          session_id: aiTutorSessionId,
          stage: input.stage,
          proposed_word: input.proposedWord,
          proposed_translation: input.proposedTranslation,
          comment: input.comment,
        },
      }),
      copy("ai_tutor_word_report_sent", "Report sent."),
    );
    if (!payload) return false;
    const response = payload as AiTutorResponse;
    const sessionID = String(response.session?.id || response.session?.ID || aiTutorSessionId);
    if (sessionID) setAiTutorSessionId(sessionID);
    setAiTutorStep(response.next_step || null);
    setAiTutorFeedback(response.feedback || null);
    return true;
  };

  const reportWord = async (input: { wordId: string; proposedWord: string; proposedTranslation: string; comment: string }) => {
    if (!input.wordId) return false;
    const payload = await runAction(
      "word-report",
      () => api<ApiRecord>("/api/words/report", {
        method: "POST",
        body: {
          word_id: input.wordId,
          proposed_word: input.proposedWord,
          proposed_translation: input.proposedTranslation,
          comment: input.comment,
        },
      }),
      copy("ai_tutor_word_report_sent", "Report sent."),
    );
    if (!payload) return false;
    const notice = copy("ai_tutor_word_report_sent", "Report sent.");
    setStatus({ kind: "ok", text: notice });
    setWordChallenge(null);
    setWordResult(null);
    const record = getRecord(payload);
    setMessages((current) => [
      panelMessage(notice, "default", copy("ai_tutor_word_report_title", "Word report"), record, "words"),
      ...current,
    ].slice(0, 12));
    return true;
  };

  const submitLearningAnswer = async (mode: "lesson" | "practice") => {
    const text = draft.trim();
    if (!text && !voiceFile && !(mode === "practice" && imageFile)) {
      setStatus({ kind: "error", text: copy("add_input_first", "Add text, voice, or a practice photo first.") });
      return;
    }
    const endpoint = mode === "lesson" ? "/api/lesson/answer" : "/api/practice";
    const hasFile = Boolean(voiceFile || (mode === "practice" && imageFile));
    const payload = await runAction(mode, async () => {
      if (!hasFile) return api<ApiRecord>(endpoint, { method: "POST", body: { text } });
      const form = new FormData();
      if (text) form.append("text", text);
      if (voiceFile) form.append("voice", voiceFile, voiceFile.name || `${mode}.webm`);
      if (mode === "practice" && imageFile) form.append("image", imageFile, imageFile.name || "practice.jpg");
      return apiForm<ApiRecord>(endpoint, form);
    });
    if (!payload) return;
    const record = getRecord(payload);
    const responseText = formatLearningRecord(record, copy, copy("done", "Done."));
    const transcript = asText(record.transcript || text || (voiceFile ? copy("voice_message", "Voice message") : imageFile ? copy("image_message", "Image message") : ""), "");
    const nextMessages = [
      panelMessage(responseText, "success", mode === "lesson" ? copy("lesson_feedback", "Lesson feedback") : copy("practice_feedback", "Practice feedback"), mode === "lesson" ? { ...record, task_id: activeLessonTaskId } : record, mode),
    ];
    if (transcript) nextMessages.push({ id: `${Date.now()}-user`, role: "user", title: copy("you", "You"), body: transcript, meta: mode });
    setMessages((current) => [
      ...nextMessages,
      ...current,
    ].slice(0, 12));
    setDraft("");
    setVoiceFile(null);
    setImageFile(null);
    if (mode === "lesson") setActiveLessonTaskId("");
  };

  const startRoleplayScenario = async (scenario: RoleplayScenario) => {
    const title = roleplayScenarioTitle(scenario, user);
    const description = roleplayScenarioDescription(scenario, user);
    const aiLine = copy("roleplay_opening_line", "Hi! Do you have a minute?");
    const body = [
      description,
      `${copy("roleplay_ai_line", "AI")}: ${aiLine}`,
      copy("roleplay_your_turn", "Твоя очередь: ответь одной фразой, и я продолжу сцену."),
    ].map((line) => cleanAppText(line)).filter(Boolean).join("\n\n");
    const result = panelMessage(body, "default", title, { scenario_id: scenario.id, source: "roleplay" }, "roleplay");
    setRoleplayResult(result);
    setMessages((current) => [
      result,
      ...current.filter((message) => message.meta !== "roleplay"),
    ].slice(0, 12));
    setDraft("");
    setView("roleplay");
  };

  const submitRoleplayAnswer = async (scenario: RoleplayScenario, text: string, voice?: File | null) => {
    const learnerText = text.trim();
    if (!learnerText && !voice) {
      setStatus({ kind: "error", text: copy("add_input_first", "Add text, voice, or a practice photo first.") });
      return;
    }
    const prompt = buildRoleplayScenarioPrompt(scenario, user, session as SessionData, copy, learnerText || copy("voice_message", "Voice message"));
    const title = roleplayScenarioTitle(scenario, user);
    const payload = await runAction("roleplay", async () => {
      if (!voice) return api<ApiRecord>("/api/practice", { method: "POST", body: { text: prompt } });
      const form = new FormData();
      form.append("text", prompt);
      form.append("voice", voice, voice.name || "roleplay.webm");
      return apiForm<ApiRecord>("/api/practice", form);
    });
    if (!payload) return;
    const record = getRecord(payload);
    const transcript = asText(learnerText || (voice ? record.transcript || copy("voice_message", "Voice message") : ""), "");
    const result = panelMessage(formatLearningRecord(record, copy, copy("roleplay_ready", "Roleplay is ready.")), "success", title, record, "roleplay");
    setRoleplayResult(result);
    setMessages((current) => [
      result,
      transcript ? userMessage(transcript, copy("you", "You"), "roleplay") : null,
      ...current,
    ].filter((message): message is ChatMessage => Boolean(message)).slice(0, 12));
    setDraft("");
    setVoiceFile(null);
    setView("roleplay");
  };

  const startShadowing = async () => {
    if (!requirePremiumFeature("premium_audio_required", "Listening is available with Premium.")) return;
    const previous = cleanAppText(shadowingTarget);
    const payload = await runAction("shadowing", () => api<ApiRecord>("/api/shadowing/start", { method: "POST", body: { previous } }), copy("phrase_ready", "Phrase is ready."));
    if (!payload) return;
    const record = getRecord(payload);
    const target = asText(record.target || record.phrase || record.text || record.message, "");
    setShadowingTarget(target);
    setView("shadowing");
  };

  const submitShadowing = async () => {
    if (!requirePremiumFeature("premium_audio_required", "Listening is available with Premium.")) return;
    const text = draft.trim();
    const submittedTarget = cleanAppText(shadowingTarget);
    if (!text && !voiceFile) {
      setStatus({ kind: "error", text: copy("add_input_first", "Add text, voice, or a practice photo first.") });
      return;
    }
    const payload = await runAction("shadowing", () => apiForm<ApiRecord>("/api/shadowing/answer", buildShadowingForm(text, voiceFile, submittedTarget)));
    if (!payload) return;
    const record = getRecord(payload);
    const heardPhrase = submittedTarget || cleanAppText(recordField(record, ["target", "phrase", "correction_audio_text", "correction"]));
    const displayRecord = heardPhrase
      ? { ...record, target: heardPhrase, phrase: heardPhrase, correction_audio_text: heardPhrase }
      : record;
    let responseText = formatLearningRecord(record, copy, copy("answer_checked", "Answer checked."));
    if (heardPhrase && !responseText.includes(heardPhrase)) {
      responseText = [responseText, `${copy("listening_heard_phrase", "Heard phrase")}: ${heardPhrase}`].filter(Boolean).join("\n\n");
    }
    const transcript = asText(record.transcript || text || (voiceFile ? copy("voice_message", "Voice message") : ""), "");
    const nextMessages = [panelMessage(responseText, "success", copy("listening_feedback", "Listening feedback"), displayRecord, "shadowing")];
    if (transcript) nextMessages.push({ id: `${Date.now()}-shadowing-user`, role: "user", title: copy("you", "You"), body: transcript, meta: "shadowing" });
    setMessages((current) => [...nextMessages, ...current].slice(0, 12));
    setDraft("");
    setVoiceFile(null);
    const nextPayload = await runAction("shadowing-next", () => api<ApiRecord>("/api/shadowing/start", { method: "POST", body: { previous: heardPhrase } }), copy("phrase_ready", "Phrase is ready."));
    if (nextPayload) {
      const nextRecord = getRecord(nextPayload);
      const nextTarget = cleanAppText(nextRecord.target || nextRecord.phrase || nextRecord.text || nextRecord.message);
      if (nextTarget) setShadowingTarget(nextTarget);
    }
  };

  const startPronunciation = async () => {
    if (!requirePremiumFeature("premium_pronunciation_required", "Pronunciation practice is available with Premium.")) return;
    const previous = cleanAppText(pronunciationTarget || shadowingTarget);
    const payload = await runAction("pronunciation-start", () => api<ApiRecord>("/api/pronunciation/start", { method: "POST", body: { previous } }), copy("phrase_ready", "Phrase is ready."));
    if (!payload) return;
    const record = getRecord(payload);
    const target = cleanAppText(recordField(record, ["target", "phrase", "text"]));
    if (target) {
      setPronunciationTarget(target);
      setMessages((current) => current.filter((message) => message.meta !== "pronunciation"));
    }
    const refreshed = record.user && typeof record.user === "object" ? record.user as UserProfile : null;
    if (refreshed) setSession((current) => current ? updateSessionUser(current, refreshed) : current);
  };

  const submitPronunciation = async (target: string) => {
    if (!requirePremiumFeature("premium_pronunciation_required", "Pronunciation practice is available with Premium.")) return null;
    const cleanTarget = cleanAppText(target).trim();
    if (!cleanTarget) {
      setStatus({ kind: "error", text: copy("pronunciation_target_missing", "Нет фразы для проверки произношения.") });
      return null;
    }
    if (!voiceFile) {
      setStatus({ kind: "error", text: copy("tutor_need_voice", "Record your voice first.") });
      return null;
    }
    const form = buildShadowingForm("", voiceFile, cleanTarget);
    const payload = await runAction("pronunciation", () => apiForm<ApiRecord>("/api/pronunciation/check", form), copy("answer_checked", "Answer checked."));
    if (!payload) return null;
    const record = getRecord(payload);
    const responseText = formatLearningRecord(record, copy, copy("pronunciation_checked", "Произношение проверено."));
    const transcript = asText(record.transcript || copy("voice_message", "Voice message"), "");
    setMessages((current) => [
      panelMessage(responseText, "success", copy("pronunciation_feedback", "Проверка произношения"), record, "pronunciation"),
      userMessage(transcript, copy("you", "You"), "pronunciation"),
      ...current.filter((message) => message.meta !== "pronunciation"),
    ].slice(0, 12));
    const refreshed = record.user && typeof record.user === "object" ? record.user as UserProfile : null;
    if (refreshed) setSession((current) => current ? updateSessionUser(current, refreshed) : current);
    setVoiceFile(null);
    return record;
  };

  const startWord = async () => {
    setWordChallenge(null);
    setWordResult(null);
    const payload = await runAction("word", () => api<WordChallenge>("/api/words/next", { method: "POST", body: {} }));
    if (!payload) return;
    setWordChallenge(payload as WordChallenge);
    setMessages((current) => [
      panelMessage(challengeMessage(payload as WordChallenge, copy), "default", copy("learn_words", "Learn words"), challengeDetails(payload as WordChallenge), "words"),
      ...current,
    ].slice(0, 12));
    setView("words");
  };

  const answerWord = async (answerId: string) => {
    const selectedText = toChoiceOptions(wordChallenge?.options).find((option) => option.id === answerId)?.text || answerId;
    const payload = await runAction("word-answer", () => api<ApiRecord>("/api/words/answer", { method: "POST", body: { answer_id: answerId } }));
    if (!payload) return;
    const record = getRecord(payload);
    setWordResult(makeTrainerResult(record, copy, answerId, "words"));
    setMessages((current) => [
      panelMessage(
        formatLearningRecord(record, copy, copy("answer_checked", "Answer checked."), { includeContext: false }),
        record.correct ? "success" : "warning",
        record.correct ? copy("saved_to_vocabulary", "Saved to vocabulary") : copy("word_trainer", "Word trainer"),
        record,
        "words",
      ),
      userMessage(selectedText, copy("you", "You"), "words"),
      ...current,
    ]);
    if (record.correct) setWordChallenge(null);
  };

  const startWordGame = async () => {
    setWordChallenge(null);
    setWordGameResult(null);
    const payload = await runAction("word-game", () => api<WordChallenge>("/api/word-game/next", { method: "POST", body: {} }));
    if (!payload) return;
    setWordChallenge(payload as WordChallenge);
    setMessages((current) => [
      panelMessage(challengeMessage(payload as WordChallenge, copy), "default", copy("word_game", "Review game"), challengeDetails(payload as WordChallenge), "word-game"),
      ...current,
    ].slice(0, 12));
    setView("word-game");
  };

  const answerWordGame = async (answerId: string) => {
    const selectedText = toChoiceOptions(wordChallenge?.options).find((option) => option.id === answerId)?.text || answerId;
    const payload = await runAction("word-game-answer", () => api<ApiRecord>("/api/word-game/answer", { method: "POST", body: { answer_id: answerId } }));
    if (!payload) return;
    const record = getRecord(payload);
    setWordGameResult(makeTrainerResult(record, copy, answerId, "word-game"));
    setMessages((current) => [
      panelMessage(
        formatLearningRecord(record, copy, copy("answer_checked", "Answer checked."), { includeContext: false }),
        record.correct ? "success" : "warning",
        copy("word_game", "Review game"),
        record,
        "word-game",
      ),
      userMessage(selectedText, copy("you", "You"), "word-game"),
      ...current,
    ]);
    if (record.correct) setWordChallenge(null);
  };

  const startSpelling = async () => {
    setSpellingChallenge(null);
    setSpellingResult(null);
    setDraft("");
    const payload = await runAction("spelling", () => api<SpellingChallenge>("/api/spelling/start", { method: "POST", body: {} }));
    if (!payload) return;
    setSpellingChallenge(payload as SpellingChallenge);
    setMessages((current) => [
      panelMessage(challengeMessage(payload as SpellingChallenge, copy), "default", copy("spelling", "Spelling"), challengeDetails(payload as SpellingChallenge), "spelling"),
      ...current,
    ].slice(0, 12));
    setView("spelling");
  };

  const answerSpelling = async (giveUp = false) => {
    const typed = draft.trim();
    const payload = await runAction("spelling-answer", () =>
      api<ApiRecord>("/api/spelling/answer", { method: "POST", body: giveUp ? { give_up: true } : { text: typed } }),
    );
    if (!payload) return;
    const record = getRecord(payload);
    setSpellingResult(makeTrainerResult(record, copy, undefined, "spelling"));
    const hideSpellingAnswer = !record.correct && !record.gave_up;
    setMessages((current) => [
      panelMessage(
        formatLearningRecord(record, copy, copy("spelling_checked", "Spelling checked."), { includeContext: false, includeMistakes: !hideSpellingAnswer }),
        record.correct ? "success" : "warning",
        copy("spelling", "Spelling"),
        record,
        "spelling",
      ),
      userMessage(giveUp ? copy("give_up", "I do not know") : typed, copy("you", "You"), "spelling"),
      ...current,
    ]);
    setDraft("");
    if (record.correct || giveUp) setSpellingChallenge(null);
  };

  const startLevel = async () => {
    setLevelQuestion(null);
    setLevelResult(null);
    const payload = await runAction("level", () => api<ApiRecord>("/api/level-test/start", { method: "POST", body: {} }));
    if (!payload) return;
    const record = getRecord(payload);
    setLevelQuestion(record.question as LevelQuestion);
    setView("level");
  };

  const answerLevel = async (answer: number) => {
    if (!levelQuestion) return;
    const payload = await runAction("level-answer", () =>
      api<ApiRecord>("/api/level-test/answer", { method: "POST", body: { index: levelQuestion.index, answer } }),
    );
    if (!payload) return;
    const record = getRecord(payload);
    if (record.complete) {
      setLevelQuestion(null);
      const result = {
        tone: "success" as const,
        title: `${copy("level_updated", "Level updated")}: ${asText(record.level, "A1")}`,
        body: "",
        record,
        mode: "level" as const,
      };
      setLevelResult(result);
      setMessages((current) => [
        panelMessage(`${result.title}\n${copy("score", "Score")}: ${asText(record.score, "0")} / ${asText(record.max_score, "0")}`, "success", copy("level_test", "Level test"), record, "level"),
        ...current,
      ]);
    } else {
      setLevelQuestion(record.question as LevelQuestion);
    }
  };

  const ensureLevelStarted = () => {
    if (levelQuestion || levelResult || busy === "level") return;
    void startLevel();
  };

  const loadVocabulary = async (page = vocabularyMeta.page, options: { navigate?: boolean } = {}): Promise<VocabularyItem[]> => {
    const payload = await runAction("vocabulary", () => api<ApiRecord>(`/api/vocabulary?page=${encodeURIComponent(String(page))}`));
    if (!payload) return [];
    const record = getRecord(payload);
    const items = Array.isArray(record.items) ? (record.items as VocabularyItem[]) : [];
    setVocabulary(items);
    setVocabularyMeta({
      page: Number(record.page || 0),
      page_size: Number(record.page_size || 10),
      total: Number(record.total || items.length),
      total_pages: Math.max(1, Number(record.total_pages || 1)),
    });
    if (options.navigate !== false) setView("vocabulary");
    return items;
  };

  const loadMistakes = async (options: { navigate?: boolean } = {}): Promise<MistakeItem[]> => {
    const payload = await runAction("mistakes", () => api<ApiRecord>("/api/mistakes"));
    if (!payload) return [];
    const record = getRecord(payload);
    const items = Array.isArray(record.items) ? (record.items as MistakeItem[]) : [];
    setMistakes(items);
    if (options.navigate !== false) setView("mistakes");
    return items;
  };

  const startMistakePractice = async (index?: number) => {
    const body = typeof index === "number" ? { index } : {};
    const payload = await runAction("mistake-practice", () => api<ApiRecord>("/api/mistakes/practice/start", { method: "POST", body }));
    if (!payload) return;
    const record = getRecord(payload);
    if (record.empty) {
      setMistakePractice(null);
      setMessages((current) => [panelMessage(asText(record.message, copy("no_mistakes_loaded", "No mistakes loaded")), "default", copy("mistakes", "Mistakes"), record, "mistakes"), ...current]);
      setView("mistakes");
      return;
    }
    const selectedMistake = (record.mistake || null) as MistakeItem | null;
    setMistakePractice(selectedMistake);
    const selectedIndex = typeof selectedMistake?.index === "number" ? selectedMistake.index : typeof index === "number" ? index : null;
    setMistakePracticeIndex(selectedIndex);
    setMistakeAnswer("");
    setMistakeResult(null);
    setView("mistakes");
  };

  const submitMistakeAnswer = async () => {
    const text = mistakeAnswer.trim();
    if (!text) {
      setStatus({ kind: "error", text: copy("type_answer", "Type the answer") });
      return;
    }
    const payload = await runAction("mistake-answer", () => api<ApiRecord>("/api/mistakes/practice/answer", { method: "POST", body: { text } }));
    if (!payload) return;
    const record = getRecord(payload);
    const mistake = (record.mistake || mistakePractice || {}) as MistakeItem;
    setMistakeResult(makeTrainerResult(record, copy, undefined));
    if (record.correct) {
      setMistakePractice(null);
      setMistakePracticeIndex(null);
      setMistakeAnswer("");
      setMessages((current) => [
        panelMessage(formatLearningRecord(record, copy, `${asText(mistake.word || mistake.correction, copy("answer_checked", "Answer checked."))} +${asText(record.xp, "8")} XP`), "success", copy("mistake_repaired", "Mistake repaired"), record, "mistakes"),
        userMessage(text, copy("you", "You"), "mistakes"),
        ...current,
      ]);
      await loadMistakes({ navigate: false });
      setStatus({ kind: "ok", text: copy("mistake_repaired", "Mistake repaired") });
      return;
    }
    setMistakePractice(mistake);
    setMistakePracticeIndex(typeof mistake.index === "number" ? mistake.index : mistakePracticeIndex);
    setMessages((current) => [panelMessage(formatLearningRecord(record, copy, copy("try_again", "Try again.")), "warning", copy("mistake_practice", "Mistake practice"), record, "mistakes"), userMessage(text, copy("you", "You"), "mistakes"), ...current]);
  };

  const cancelMistakePractice = () => {
    setMistakePractice(null);
    setMistakePracticeIndex(null);
    setMistakeAnswer("");
    setMistakeResult(null);
  };

  const deleteMistake = async (index: number) => {
    const payload = await runAction("mistake-delete", () =>
      api<ApiRecord>("/api/mistakes/delete", { method: "POST", body: { index } }),
      copy("mistake_deleted", "Mistake deleted."),
    );
    if (!payload) return;
    const record = getRecord(payload);
    if (Array.isArray(record.items)) {
      setMistakes(record.items as MistakeItem[]);
    } else {
      await loadMistakes({ navigate: false });
    }
    if (mistakePracticeIndex === index) {
      setMistakePractice(null);
      setMistakePracticeIndex(null);
      setMistakeAnswer("");
      setMistakeResult(null);
    }
  };

  const clearMistakes = async () => {
    const payload = await runAction("mistakes-clear", () => api<ApiRecord>("/api/mistakes/clear", { method: "POST", body: {} }), copy("mistakes_cleared", "Mistakes cleared."));
    if (!payload) return;
    setMistakes([]);
    setMistakePractice(null);
    setMistakePracticeIndex(null);
    setMistakeAnswer("");
    setMistakeResult(null);
  };

  const submitTool = async () => {
    const text = draft.trim();
    if (toolMode === "translator" && !text) {
      setStatus({ kind: "error", text: copy("paste_text_translate", "Paste text to translate") });
      return;
    }
    if (toolMode === "voice" && !toolVoiceFile) {
      setStatus({ kind: "error", text: copy("voice_file", "Voice file") });
      return;
    }
    if (toolMode === "image" && !toolImageFile) {
      setStatus({ kind: "error", text: copy("practice_photo", "Practice photo") });
      return;
    }
    const form = new FormData();
    form.append("source_language", toolSourceLanguage || "auto");
    form.append("target_language", toolTargetLanguage || user.interface_language || "ru");
    if (text) form.append("text", text);
    if (toolVoiceFile) form.append("voice", toolVoiceFile, toolVoiceFile.name || "tool-voice.webm");
    if (toolImageFile) form.append("image", toolImageFile, toolImageFile.name || "tool-image.jpg");
    const endpoint = toolMode === "voice" ? "/api/tools/voice-text" : toolMode === "image" ? "/api/tools/image-translate" : "/api/tools/translator";
    const payload = await runAction(`tools-${toolMode}`, () => apiForm<ApiRecord>(endpoint, form), copy("tool_result_ready", "Tool result is ready."));
    if (!payload) return;
    const record = getRecord(payload);
    const result: TranslatorResult = {
      source_text: asText(record.source_text || record.transcript || text, ""),
      translation: asText(record.translation, ""),
      result: asText(record.result || record.translation || record.transcript, copy("done", "Done.")),
      source_language: asText(record.source_language || toolSourceLanguage, ""),
      target_language: asText(record.target_language || toolTargetLanguage, ""),
    };
    setToolResult(result);
    const toolInput = text || (toolVoiceFile ? copy("voice_message", "Voice message") : toolImageFile ? copy("image_message", "Image message") : "");
    setMessages((current) => [
      panelMessage(formatLearningRecord(record, copy, result.result || result.translation || result.source_text || copy("done", "Done.")), "success", copy("tools", "Tools"), record, `tools:${toolMode}`),
      ...(toolInput ? [userMessage(toolInput, copy("you", "You"), `tools:${toolMode}`)] : []),
      ...current,
    ]);
    if (toolMode === "translator") setDraft("");
    if (toolMode === "voice") setToolVoiceFile(null);
    if (toolMode === "image") setToolImageFile(null);
  };

  const loadLeaderboard = async (language = leaderboardLanguage || "global", options: { navigate?: boolean } = {}) => {
    const normalizedLanguage = language === "global" ? "all" : language || user.learning_language || "en";
    const payload = await runAction("leaderboard", () => api<ApiRecord>(`/api/leaderboard?language=${encodeURIComponent(normalizedLanguage)}`));
    if (!payload) return;
    const record = getRecord(payload);
    const global = language === "global" || normalizedLanguage === "all";
    const entries = global
      ? Array.isArray(record.overall)
        ? (record.overall as LeaderboardEntry[])
        : []
      : Array.isArray(record.language_entries)
        ? (record.language_entries as LeaderboardEntry[])
        : Array.isArray(record.overall)
          ? (record.overall as LeaderboardEntry[])
          : [];
    setLeaderboard(entries.slice(0, 10));
    setLeaderboardLanguage(global ? "global" : asText(record.language || language, language));
    setLeaderboardMeta({
      language: global ? "global" : asText(record.language || language, language),
      languageName: global ? copy("global_top", "Global top") : asText(record.language_name || language, language),
    });
    if (options.navigate !== false) setView("leaderboard");
  };

  const loadPremiumPlans = async () => {
    const payload = await runAction("premium-plans", () => api<ApiRecord>("/api/premium/plans"));
    if (!payload) return;
    const record = getRecord(payload);
    setPremiumPlans(Array.isArray(record.plans) ? (record.plans as PremiumPlan[]) : session?.premium_plans || []);
    setView("premium");
  };

  const startPremiumPayment = async (plan: PremiumPlan, method: "card" | "stars", yooKassaMethod?: string) => {
    const endpoint =
      method === "card"
        ? "/api/premium/payment"
        : "/api/premium/stars";
    const payload = await runAction(`premium-${method}`, () =>
      api<ApiRecord>(endpoint, {
        method: "POST",
        body:
          method === "card"
            ? { product: plan.product, payment_method: yooKassaMethod || "" }
            : { product: plan.product },
      }),
    );
    if (!payload) return;
    const record = getRecord(payload);
    setPayment({ plan, info: record });
    const url = asText(record.confirmation_url || record.bot_url || record.invoice_url || record.url, "");
    if (url) window.open(url, "_blank", "noopener,noreferrer");
    setStatus({ kind: "ok", text: asText(record.message, copy("payment_window_ready", "Payment window is ready.")) });
  };

  const showPaymentSuccessNotice = async (record: ApiRecord, plan = payment?.plan) => {
    const nextSession = await refreshSession().catch(() => null);
    const nextUser = nextSession?.user || user;
    const premiumUntil = asText(record.premium_until || record.expires_at || nextUser.premium_until, "");
    setPaymentNotice({
      kind: "success",
      title: copy("payment_success_title", "Оплата подтверждена"),
      description: copy("payment_success_body", "Подписка обновлена. Premium уже активен в вашем аккаунте."),
      planTitle: plan?.title || nextUser.plan || copy("premium", "Premium"),
      period: premiumPeriodLabel(plan, copy),
      premiumUntil: premiumUntil ? prettyDate(premiumUntil) : "",
      amount: paymentAmountLine(record),
    });
  };

  useEffect(() => {
    const timer = window.setInterval(() => {
      void refreshSession().catch(() => undefined);
    }, 15 * 60 * 1000);
    return () => window.clearInterval(timer);
  }, []);

  const saveSettings = async (payload: ApiRecord) => {
    await runAction("settings", () => api<SessionData>("/api/settings", { method: "POST", body: payload }), copy("settings_saved", "Settings saved."));
  };

  const saveNavigationLayout = (patch: NonNullable<UserProfile["navigation_layout"]>) => {
    const nextLayout = { ...(session?.user?.navigation_layout || {}), ...patch };
    setSession((current) => current?.user ? updateSessionUser(current, { ...current.user, navigation_layout: nextLayout }) : current);
    void api<SessionData>("/api/navigation-layout", { method: "POST", body: nextLayout })
      .then((nextSession) => setSession((current) => {
        const source = nextSession || current;
        return source?.user ? updateSessionUser(source, { ...source.user, navigation_layout: nextLayout }) : source;
      }))
      .catch(() => undefined);
  };

  const changePassword = async (payload: { current_password: string; new_password: string; new_password_confirm: string }) => {
    await runAction("password", () => api<ApiRecord>("/api/auth/password", { method: "POST", body: payload }), copy("password_changed", "Password changed."));
    await refreshSession().catch(() => undefined);
  };

  const startTelegramCode = async (body: ApiRecord = {}): Promise<TelegramCodeRequest | null> => {
    let telegramWindow: Window | null = null;
    try {
      telegramWindow = window.open("about:blank", "_blank");
      if (telegramWindow) telegramWindow.opener = null;
    } catch {
      telegramWindow = null;
    }
    const payload = await runAction(
      "telegram-code",
      () => api<ApiRecord>("/api/auth/telegram/start", { method: "POST", body }),
      copy("telegram_code_sent", "Код отправлен в Telegram."),
    );
    if (!payload) {
      telegramWindow?.close();
      return null;
    }
    const record = getRecord(payload);
    const token = asText(record.token, "");
    const botURL = asText(record.bot_url, "");
    const expiresAt = asText(record.expires_at, "");
    if (!token || !botURL) {
      telegramWindow?.close();
      setStatus({ kind: "error", text: copy("telegram_code_start_failed", "Не удалось подготовить Telegram-код.") });
      return null;
    }
    if (telegramWindow) {
      telegramWindow.location.href = botURL;
    } else {
      window.open(botURL, "_blank", "noopener,noreferrer");
    }
    return { token, botURL, expiresAt };
  };

  const verifyTelegramCode = async (token: string, code: string): Promise<boolean> => {
    if (!token) {
      setStatus({ kind: "error", text: copy("telegram_code_missing", "Сначала отправьте код в Telegram.") });
      return false;
    }
    const payload = await runAction(
      "telegram-code-verify",
      () => api<SessionData | ApiRecord>("/api/auth/telegram/status", { method: "POST", body: { token, code } }),
      copy("telegram_linked_success", "Telegram подключен."),
    );
    if (!payload) return false;
    const nextSession = payload as SessionData;
    if (nextSession.authenticated) {
      setSession(nextSession);
    } else {
      await refreshSession().catch(() => undefined);
    }
    return true;
  };

  const activateKey = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const key = activationKey.trim();
    if (!key) return;
    const payload = await runAction("activation", () => api<SessionData>("/api/premium/activation-key", { method: "POST", body: { key } }), copy("activation_accepted", "Activation key accepted."));
    if (payload) {
      setActivationKey("");
      await refreshSession().catch(() => undefined);
    }
  };

  const completeOnboarding = async (payload: { goal: string; level: string; learning_language: string; format: string }) => {
    localStorage.setItem(onboardingKey, JSON.stringify({ ...payload, completed_at: new Date().toISOString() }));
    setOnboardingOpen(false);
    await saveSettings({
      interface_language: user.interface_language || "ru",
      learning_language: payload.learning_language,
      level: payload.level,
      onboarding_goal: payload.goal,
      onboarding_format: payload.format,
      learning_focus: `${payload.goal}; ${payload.format}`,
    });
  };

  const onLogout = async () => {
    setBusy("logout");
    try {
      await logout();
    } finally {
      setSession(null);
      location.replace(appLoginPath);
    }
  };

  useEffect(() => {
    if (!session?.authenticated || busy === "daily-bonus") return;
    const todayKey = localDateKey();
    const current = habitLog[todayKey] || updateHabitForToday({}, user)[todayKey];
    if (!current?.complete || dailyBonusLocked(current, user)) return;
    if (autoDailyBonusClaimRef.current === todayKey) return;
    autoDailyBonusClaimRef.current = todayKey;
    void claimDailyBonus();
  }, [
    session?.authenticated,
    busy,
    habitLog,
    user.lessons_today,
    user.practice_today,
    user.voice_today,
    user.lesson_limit,
    user.practice_limit,
    user.voice_limit,
    user.daily_bonus_claims,
    user.daily_bonus_last_claimed_at,
  ]);

  const activateView = (view: ViewId) => {
    view = normalizeView(view);
    if (view === "shadowing" && !requirePremiumFeature("premium_audio_required", "Listening is available with Premium.")) return;
    if (view === "pronunciation" && !requirePremiumFeature("premium_pronunciation_required", "Pronunciation practice is available with Premium.")) return;
    setNotFoundPath("");
    setActiveView(view);
    if (view === "shadowing") void startShadowing();
    if (view === "words") void startWord();
    if (view === "word-game") void startWordGame();
    if (view === "spelling") void startSpelling();
    if (view === "level") ensureLevelStarted();
    if (view === "vocabulary") void loadVocabulary(0);
    if (view === "offline") {
      void loadVocabulary(0, { navigate: false });
      void loadMistakes({ navigate: false });
    }
    if (view === "mistakes") void loadMistakes();
    if (view === "leaderboard") void loadLeaderboard();
    if (view === "dashboard") {
      void loadVocabulary(0, { navigate: false });
      void loadMistakes({ navigate: false });
      void loadLeaderboard(undefined, { navigate: false });
    }
    if (view === "premium") void loadPremiumPlans();
  };

  useEffect(() => {
    if (!session?.authenticated || bootstrappedRef.current) return;
    bootstrappedRef.current = true;
    if (activeView === "shadowing") void startShadowing();
    if (activeView === "pronunciation" && !user.premium) {
      setStatus({ kind: "info", text: copy("premium_pronunciation_required", "Pronunciation practice is available with Premium.") });
      setActiveView("premium");
      void loadPremiumPlans();
      return;
    }
    if (activeView === "words") void startWord();
    if (activeView === "word-game") void startWordGame();
    if (activeView === "spelling") void startSpelling();
    if (activeView === "level") void startLevel();
    if (activeView === "vocabulary") void loadVocabulary(0);
    if (activeView === "offline") {
      void loadVocabulary(0, { navigate: false });
      void loadMistakes({ navigate: false });
    }
    if (activeView === "mistakes") void loadMistakes();
    if (activeView === "leaderboard") void loadLeaderboard();
    if (activeView === "dashboard") {
      void loadVocabulary(0, { navigate: false });
      void loadMistakes({ navigate: false });
      void loadLeaderboard(undefined, { navigate: false });
    }
    if (activeView === "premium") void loadPremiumPlans();
  }, [session?.authenticated]);

  const primaryAction = () => {
    if (activeView === "tutor") return startTutor();
    if (activeView === "lesson") return startLesson();
    if (activeView === "practice") return submitLearningAnswer("practice");
    if (activeView === "roleplay") return setActiveView("roleplay");
    if (activeView === "shadowing") return startShadowing();
    if (activeView === "pronunciation") return requirePremiumFeature("premium_pronunciation_required", "Pronunciation practice is available with Premium.") ? setActiveView("pronunciation") : undefined;
    if (activeView === "words") return startWord();
    if (activeView === "word-game") return startWordGame();
    if (activeView === "spelling") return startSpelling();
    if (activeView === "level") return startLevel();
    if (activeView === "mistakes") return loadMistakes();
    if (activeView === "leaderboard") return loadLeaderboard();
    if (activeView === "vocabulary") return loadVocabulary();
    if (activeView === "offline") {
      void loadVocabulary(vocabularyMeta.page, { navigate: false });
      return loadMistakes({ navigate: false });
    }
    if (activeView === "dashboard") return loadLeaderboard(undefined, { navigate: false });
    if (activeView === "premium") return loadPremiumPlans();
    return startLesson();
  };

  const authStandalone = isAuthStandaloneRoute();
  if (sessionLoading) return <LoadingScreen theme={theme} />;
  if (authStandalone) {
    if (session?.authenticated && session.user) {
      location.replace("/app");
      return <LoadingScreen theme={theme} />;
    }
    return <AuthStandaloneView session={session} theme={theme} setTheme={setTheme} onAuthenticated={setSession} />;
  }
  if (!session?.authenticated || !session.user) {
    location.replace(appLoginPath);
    return <LoadingScreen theme={theme} />;
  }

  const viewProps: ViewRendererProps = {
    activeView,
    user,
    copy,
    busy,
    messages,
    activeLessonTaskId,
    draft,
    setDraft,
    voiceFile,
    imageFile,
    setVoiceFile,
    setImageFile,
    setView: activateView,
    aiTutorStep,
    aiTutorFeedback,
    submitAiTutorStep,
    reportAiTutorWord,
    tutorLoadError,
    startTutor,
    restartTutorLesson,
    startLesson,
    submitLesson: () => submitLearningAnswer("lesson"),
    submitPractice: () => submitLearningAnswer("practice"),
    startRoleplayScenario,
    submitRoleplayAnswer,
    roleplayResult,
    startShadowing,
    submitShadowing,
    pronunciationTarget,
    setPronunciationTarget,
    startPronunciation,
    submitPronunciation,
    startWord,
    wordChallenge,
    answerWord,
    reportWord,
    wordResult,
    startWordGame,
    answerWordGame,
    wordGameResult,
    spellingChallenge,
    startSpelling,
    answerSpelling,
    spellingResult,
    levelQuestion,
    startLevel,
    answerLevel,
    levelResult,
    vocabulary,
    vocabularyMeta,
    loadVocabulary,
    mistakes,
    loadMistakes,
    mistakePractice,
    mistakePracticeIndex,
    mistakeAnswer,
    setMistakeAnswer,
    startMistakePractice,
    submitMistakeAnswer,
    cancelMistakePractice,
    mistakeResult,
    deleteMistake,
    clearMistakes,
    shadowingTarget,
    leaderboard,
    leaderboardMeta,
    leaderboardLanguage,
    setLeaderboardLanguage,
    loadLeaderboard,
    premiumPlans: premiumPlans.length ? premiumPlans : session.premium_plans || [],
    loadPremiumPlans,
    paymentHistory,
    setPayment,
    saveSettings,
    changePassword,
    startTelegramCode,
    verifyTelegramCode,
    session,
    activationKey,
    setActivationKey,
    activateKey,
    toolMode,
    setToolMode,
    toolSourceLanguage,
    setToolSourceLanguage,
    toolTargetLanguage,
    setToolTargetLanguage,
    toolVoiceFile,
    setToolVoiceFile,
    toolImageFile,
    setToolImageFile,
    toolResult,
    submitTool,
    selectedAwardLevel,
    setSelectedAwardLevel,
    phrasebook,
    isPhraseSaved,
    savePhrase,
    removePhrase,
    habitLog,
    dailyBonus,
    claimDailyBonus,
    openBugReport: () => setBugReportOpen(true),
    openGuide: () => setGuideOpen(true),
  };
  const handleContextPaste = (event: ClipboardEvent<HTMLElement>) => {
    if (activeView !== "tools" || toolMode !== "image" || event.defaultPrevented) return;
    const pastedImage = firstClipboardImage(event);
    if (!pastedImage) return;
    event.preventDefault();
    setToolImageFile(nameClipboardImage(pastedImage));
  };

  return (
    <div className="v2-shell" data-v2-shell="ios-function-ribbon" style={{ "--v2-brightness": brightness / 100 } as CSSProperties}>
      <ThemeBackground theme={theme} />
      <AppCookieConsentBanner language={user.interface_language || "ru"} />
      <div className="v2-app">
        <TopBar
          user={user}
          session={session}
          theme={theme}
          setTheme={setTheme}
          brightness={brightness}
          setBrightness={setBrightness}
          copy={copy}
          onLogout={onLogout}
          onLanguageChange={(code) => saveSettings({ interface_language: code })}
          setView={activateView}
          busy={busy}
          onBugReport={() => setBugReportOpen(true)}
          onGuide={() => setGuideOpen(true)}
        />
        {activeView === "home" ? (
          <MobileQuickControls
            user={user}
            session={session}
            theme={theme}
            setTheme={setTheme}
            copy={copy}
            onLogout={onLogout}
            onLanguageChange={(code) => saveSettings({ interface_language: code })}
            onBugReport={() => setBugReportOpen(true)}
            onGuide={() => setGuideOpen(true)}
            busy={busy}
          />
        ) : null}
        {status ? <StatusBanner kind={status.kind} text={status.text} onClose={() => setStatus(null)} /> : null}
        {xpGain ? (
          <div className="xp-gain-pop-v2" role="status" aria-live="polite">
            <span className="xp-gain-pop-v2__burst" aria-hidden="true"><Sparkles size={26} /></span>
            <span className="xp-gain-pop-v2__body">
              <span className="xp-gain-pop-v2__label">{xpGain.title || copy("xp_gained_label", "XP gained")}</span>
              <strong>+{xpGain.xp} XP</strong>
              {xpGain.total ? <span className="xp-gain-pop-v2__total">XP {compactNumber(xpGain.total)}</span> : null}
            </span>
          </div>
        ) : null}
        <FunctionRibbon
          nav={nav}
          activeView={activeView}
          setView={activateView}
          theme={theme}
          copy={copy}
          accountKey={accountKey}
          navigationLayout={user.navigation_layout || {}}
          onNavigationLayoutChange={(patch) => saveNavigationLayout(patch)}
        />
        {payment ? (
          <PaymentModal
            payment={payment}
            copy={copy}
            busy={busy}
            onClose={() => setPayment(null)}
            onPay={(method, yooKassaMethod) => void startPremiumPayment(payment.plan, method, yooKassaMethod)}
          />
        ) : null}
        {paymentNotice ? <PaymentNoticeDialog notice={paymentNotice} copy={copy} onClose={() => setPaymentNotice(null)} /> : null}
        {session ? (
          <OnboardingDialog
            open={onboardingOpen}
            session={session}
            user={user}
            copy={copy}
            onComplete={completeOnboarding}
            onSkip={() => {
              localStorage.setItem(onboardingKey, JSON.stringify({ skipped_at: new Date().toISOString() }));
              setOnboardingOpen(false);
            }}
          />
        ) : null}
        {bugReportOpen ? (
          <BugReportDialog
            activeView={activeView}
            copy={copy}
            busy={busy}
            onClose={() => setBugReportOpen(false)}
            onSubmit={(form) =>
              runAction("bug-report", () => apiForm<ApiRecord>("/api/bug-report", form), copy("bug_report_sent", "Bug report saved."))
                .then((payload) => {
                  if (payload) setBugReportOpen(false);
              })
            }
          />
        ) : null}
        {guideOpen ? <AppGuideDialog copy={copy} onClose={() => setGuideOpen(false)} /> : null}
        <main
          key={notFoundPath || activeView}
          className={cn("context-display", `context-display--${notFoundPath ? "not-found" : activeView}`)}
          data-view={notFoundPath ? "not-found" : activeView}
          data-mode={notFoundPath ? "not-found" : activeView}
          data-mobile-controls={activeView === "home" ? "true" : "false"}
          onPaste={handleContextPaste}
        >
          {notFoundPath ? (
            <NotFoundPage
              kicker={copy("not_found_kicker", "404")}
              title={copy("not_found_title", "Страница не найдена")}
              description={copy("not_found_body", "Такого раздела нет или ссылка устарела. Вернитесь на главный экран приложения.")}
              actionLabel={copy("not_found_home", "На главный экран")}
              onHome={() => {
                setNotFoundPath("");
                activateView("home");
              }}
            />
          ) : (
            <ViewRenderer {...viewProps} />
          )}
        </main>
        <MobileBottomNav
          nav={nav}
          activeView={activeView}
          setView={activateView}
          theme={theme}
          copy={copy}
          accountKey={accountKey}
          navigationLayout={user.navigation_layout || {}}
          onNavigationLayoutChange={(patch) => saveNavigationLayout(patch)}
        />
      </div>
    </div>
  );
}

function ThemeBackground({ theme }: { theme: Theme }) {
  return (
    <div className={cn("v2-bg", theme === "dark" ? "v2-bg--dark" : "v2-bg--light")} aria-hidden="true" />
  );
}

function LoadingScreen({ theme }: { theme: Theme }) {
  return (
    <div className="v2-loading" data-theme="dark">
      <ThemeBackground theme="dark" />
      <div className="v2-loading__card">
        <BrandMark />
        <Spinner size="small" />
      </div>
    </div>
  );
}

function apiErrorCode(error: unknown): string {
  if (!(error instanceof ApiError) || !error.payload || typeof error.payload !== "object") return "";
  const payload = error.payload as Record<string, unknown>;
  const apiError = payload.error;
  if (!apiError || typeof apiError !== "object") return "";
  return asText((apiError as Record<string, unknown>).code, "").trim();
}

function localizedAuthError(error: unknown, copy: (key: string, fallback: string) => string): string {
  const code = apiErrorCode(error);
  if (code) {
    return copy(`auth_error_${code}`, copy("auth_error_failed", "Authentication failed."));
  }
  const message = cleanAppText(error instanceof Error ? error.message : "");
  if (!message || /^(Section|Раздел|Mục):/i.test(message)) return copy("auth_error_failed", "Authentication failed.");
  return message;
}

function keepShortAuthPrefix(text: string): string {
  return text.replace(/^([A-Za-zА-Яа-яЁёІіЇїЄєҐґЎў])\s+(?=\S)/u, "$1\u00a0");
}

function localizedAPIError(error: unknown, copy: (key: string, fallback: string) => string): string {
  const code = apiErrorCode(error);
  if (code) {
    const apiKey = `api_error_${code}`;
    const apiMessage = cleanAppText(copy(apiKey, ""));
    if (apiMessage && apiMessage !== apiKey) return apiMessage;
    return copy(`auth_error_${code}`, copy("request_failed", "Request failed"));
  }
  const message = cleanAppText(error instanceof Error ? error.message : "");
  if (!message || isGenericSectionCopy(message)) return copy("request_failed", "Request failed");
  return message;
}

function AuthStandaloneView({
  session,
  theme,
  setTheme,
  onAuthenticated,
}: {
  session: SessionData | null;
  theme: Theme;
  setTheme: (theme: Theme) => void;
  onAuthenticated: (session: SessionData) => void;
}) {
  const storedLanguage = localStorage.getItem("poliglot-auth-language") || session?.user?.interface_language || languageCode(navigator.language || "ru");
  const [language, setLanguage] = useState(storedLanguage);
  const [mode, setMode] = useState<SignInMode>("login");
  const [status, setStatus] = useState<{ kind: StatusKind; text: string } | null>(null);
  const [busy, setBusy] = useState<string | null>(null);
  const [captchaToken, setCaptchaToken] = useState("");
  const [privacyAccepted, setPrivacyAccepted] = useState(false);
  const captchaRef = useRef<HTMLDivElement | null>(null);
  const captchaWidgetRef = useRef("");
  const interfaceLanguages = session?.interface_languages?.length
    ? session.interface_languages
    : [{ code: language || "ru", native_name: (language || "ru").toUpperCase() }];
  const copy = useCallback((key: string, fallback: string) => appCopy(language, key, fallback), [language]);
  const captcha = session?.captcha;
  const captchaEnabled = Boolean(captcha?.enabled && captcha.provider === "turnstile" && captcha.site_key);
  const legalLanguage = languageCode(language || "ru");
  const privacyURL = legalDocumentURL("privacy.html", legalLanguage);
  const consentURL = legalDocumentURL("consent.html", legalLanguage);
  const agreementURL = legalDocumentURL("agreement.html", legalLanguage);
  const heroImageSrc = theme === "dark" ? "/app/assets/auth-login-hero-dark.png?v=flux2-20260527-auth" : "/app/assets/auth-login-hero-light.png?v=flux2-20260527-auth";
  const testimonials = useMemo<Testimonial[]>(
    () => [
      {
        avatarSrc: "/app/assets/award-09.png",
        name: copy("auth_lessons", "Lessons"),
        handle: "@neriva",
        text: copy("auth_card_lessons", "Short AI lessons adapt to your level."),
      },
      {
        avatarSrc: "/app/assets/award-12.png",
        name: copy("auth_words", "Words"),
        handle: "@daily",
        text: copy("auth_card_words", "Daily route keeps words, practice, and review together."),
      },
      {
        avatarSrc: "/app/assets/award-16.png",
        name: copy("telegram", "Telegram"),
        handle: "@account",
        text: copy("auth_card_telegram", "Website and Telegram use one learning profile."),
      },
    ],
    [copy],
  );

  useEffect(() => {
    localStorage.setItem("poliglot-auth-language", language);
  }, [language]);

  useEffect(() => {
    if (!captchaEnabled || !captchaRef.current || captchaWidgetRef.current) return;
    let cancelled = false;
    const currentWindow = window as Window & { turnstile?: TurnstileAPI };
    const renderCaptcha = () => {
      if (cancelled || !captchaRef.current || !currentWindow.turnstile || captchaWidgetRef.current) return;
      captchaWidgetRef.current = currentWindow.turnstile.render(captchaRef.current, {
        sitekey: captcha?.site_key || "",
        callback: (token) => setCaptchaToken(token || ""),
        "expired-callback": () => setCaptchaToken(""),
        "error-callback": () => setCaptchaToken(""),
      });
    };
    if (currentWindow.turnstile) {
      renderCaptcha();
      return () => {
        cancelled = true;
        if (captchaWidgetRef.current) currentWindow.turnstile?.remove?.(captchaWidgetRef.current);
        captchaWidgetRef.current = "";
      };
    }
    const scriptID = "cf-turnstile-v2-auth";
    const existing = document.getElementById(scriptID) as HTMLScriptElement | null;
    const script = existing || document.createElement("script");
    script.id = scriptID;
    script.src = "https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit";
    script.async = true;
    script.defer = true;
    script.onload = renderCaptcha;
    script.onerror = () => setStatus({ kind: "error", text: copy("captcha_load_failed", "Captcha did not load. Refresh the page and try again.") });
    if (!existing) document.head.appendChild(script);
    return () => {
      cancelled = true;
      if (captchaWidgetRef.current) currentWindow.turnstile?.remove?.(captchaWidgetRef.current);
      captchaWidgetRef.current = "";
    };
  }, [captcha?.site_key, captchaEnabled, copy]);

  const afterAuth = (payload: unknown) => {
    const nextSession = payload as SessionData;
    if (!nextSession?.authenticated) {
      setStatus({ kind: "error", text: copy("auth_failed", "Authentication failed") });
      return;
    }
    onAuthenticated(nextSession);
    location.replace("/app");
  };

  const submitAuth = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    const login = asText(form.get("login"), "").trim();
    const password = asText(form.get("password"), "");
    const passwordConfirm = asText(form.get("password_confirm"), "");
    const referralCode = asText(form.get("referral_code"), new URLSearchParams(location.search).get("ref") || "").trim();
    if (!login || !password) {
      setStatus({ kind: "error", text: copy("auth_fill_required", "Enter login and password.") });
      return;
    }
    if (mode === "register" && password !== passwordConfirm) {
      setStatus({ kind: "error", text: copy("auth_password_mismatch", "Passwords do not match.") });
      return;
    }
    if (mode === "register" && !privacyAccepted) {
      setStatus({ kind: "error", text: copy("auth_privacy_required", "Accept the personal data policy to continue.") });
      return;
    }
    if (captchaEnabled && !captchaToken) {
      setStatus({ kind: "error", text: copy("captcha_required", "Confirm you are not a robot.") });
      return;
    }
    setBusy("auth");
    setStatus(null);
    try {
      const endpoint = mode === "register" ? "/api/auth/register" : "/api/auth/login";
      const payload = await api<SessionData>(endpoint, {
        method: "POST",
        body: {
          login,
          password,
          password_confirm: mode === "register" ? passwordConfirm : undefined,
          referral_code: mode === "register" ? referralCode : undefined,
          interface_language: language,
          captcha_token: captchaToken,
          privacy_consent: mode === "register" ? privacyAccepted : undefined,
          privacy_policy_url: mode === "register" ? privacyURL : undefined,
          personal_data_consent_url: mode === "register" ? consentURL : undefined,
          user_agreement_url: mode === "register" ? agreementURL : undefined,
        },
      });
      afterAuth(payload);
    } catch (error) {
      setStatus({ kind: "error", text: localizedAuthError(error, copy) });
      const currentWindow = window as Window & { turnstile?: TurnstileAPI };
      if (captchaWidgetRef.current) currentWindow.turnstile?.reset?.(captchaWidgetRef.current);
      setCaptchaToken("");
    } finally {
      setBusy(null);
    }
  };

  const captchaSlot = captchaEnabled ? (
    <div className="auth-turnstile-slot-v2" ref={captchaRef} data-testid="auth-captcha" aria-label={copy("captcha_required", "Confirm you are not a robot.")} />
  ) : null;
  const authTitle = mode === "register" ? copy("auth_create_account", "Create account") : keepShortAuthPrefix(copy("auth_welcome", "Welcome back"));

  const extraFields = (
    <>
      <div className="auth-page-controls-v2 animate-element animate-delay-500">
        <AnimatedSelect
          className="auth-page-language-v2"
          value={language}
          options={interfaceLanguages}
          icon={<Languages size={15} />}
          ariaLabel={copy("interface_language", "Interface language")}
          disabled={busy === "auth"}
          onChange={setLanguage}
        />
        <AnimatedThemeToggle className="auth-page-theme-v2" theme={theme} onChange={setTheme} />
      </div>
      <label className="auth-privacy-v2 animate-element animate-delay-500">
        <input
          type="checkbox"
          name="privacy_consent"
          checked={privacyAccepted}
          onChange={(event) => setPrivacyAccepted(event.target.checked)}
        />
        <span>
          {copy("auth_privacy_consent", "I agree to personal data processing under the policy.")}{" "}
          <a href={privacyURL} target="_blank" rel="noreferrer">{copy("auth_privacy_policy", "Privacy policy")}</a>
          {" · "}
          <a href={consentURL} target="_blank" rel="noreferrer">{copy("auth_personal_data_consent", "Personal data consent")}</a>
          {" · "}
          <a href={agreementURL} target="_blank" rel="noreferrer">{copy("auth_user_agreement", "User agreement")}</a>
        </span>
      </label>
    </>
  );

  return (
    <>
      <SignInPage
        mode={mode}
        title={<span>{authTitle}</span>}
        description={mode === "register" ? copy("auth_register_hint", "Choose a login and password.") : copy("auth_fill_required", "Enter login and password.")}
        heroImageSrc={heroImageSrc}
        testimonials={testimonials}
        status={status ? <span className={`is-${status.kind}`}>{status.text}</span> : null}
        captchaSlot={captchaSlot}
        extraFields={extraFields}
        submitting={busy === "auth"}
        labels={{
          login: copy("auth_username", "Login"),
          loginPlaceholder: copy("auth_username", "Login"),
          password: copy("auth_password", "Password"),
          passwordPlaceholder: copy("auth_password_placeholder", "Enter your password"),
          passwordHint: copy("auth_password_hint", "Use at least 8 characters."),
          passwordConfirm: copy("auth_password_repeat", "Repeat password"),
          passwordConfirmPlaceholder: copy("auth_password_confirm_placeholder", "Repeat the same password"),
          passwordConfirmHint: copy("auth_password_confirm_hint", "Repeat the password exactly."),
          remember: copy("keep_signed_in", "Keep me signed in"),
          resetPassword: copy("auth_support", "Account recovery"),
          submit: mode === "register" ? copy("auth_create", "Create") : copy("auth_login", "Sign in"),
          divider: copy("auth_subtitle", "Sign in or create an account"),
          social: copy("auth_telegram_title", "Sign in with Telegram"),
          createPrompt: copy("auth_new_account", "New account"),
          createAction: copy("auth_create_account", "Create account"),
          loginPrompt: copy("already_have_account", "Already have an account?"),
          loginAction: copy("auth_login", "Sign in"),
          referral: copy("auth_referral", "Referral code"),
          referralPlaceholder: copy("auth_referral_placeholder", "8 characters"),
          hidePassword: copy("hide_password", "Hide password"),
          showPassword: copy("show_password", "Show password"),
        }}
        onSignIn={submitAuth}
        onResetPassword={() => window.open("https://t.me/NERIVAapp_bot", "_blank", "noopener,noreferrer")}
        onCreateAccount={() => setMode("register")}
        onLoginMode={() => setMode("login")}
      />
      <AppCookieConsentBanner language={language} />
    </>
  );
}

function BrandMark({ compact = false, subtitle }: { compact?: boolean; subtitle?: string }) {
  return (
    <div className={cn("v2-brand", compact && "v2-brand--compact")}>
      <img src="/app/assets/brand-logo-mini.png" alt="" />
      <span>
        <strong>NERIVA</strong>
        {!compact ? <small>{subtitle || "AI Tutor"}</small> : null}
      </span>
    </div>
  );
}

function optionLabel(option: LanguageOption | { code: string; name?: string; native_name?: string }) {
  return option.native_name || option.name || option.code;
}

function languageCode(code?: string) {
  return (code || "en").split("-")[0].toLowerCase();
}

function legalDocumentOrigin() {
  const fallback = "https://neriva.ru";
  if (typeof window === "undefined") return fallback;
  const hostname = window.location.hostname.toLowerCase();
  if (hostname === "neriva.ru" || hostname === "www.neriva.ru") {
    return "https://neriva.ru";
  }
  return fallback;
}

function legalDocumentURL(path: "privacy.html" | "consent.html" | "agreement.html" | "terms.html", language: string) {
  return `${legalDocumentOrigin()}/${path}?lang=${encodeURIComponent(languageCode(language || "ru"))}`;
}

type AppCookieConsentChoice = {
  version: 1;
  necessary: true;
  analyticsMarketing: boolean;
};

type MetrikaFunction = ((...args: unknown[]) => void) & {
  a?: unknown[][];
  l?: number;
};

const appCookieConsentStorageKey = "poliglot-app-cookie-consent";
const appYandexMetrikaId = 109326597;

function serializeAppCookieConsent(analyticsMarketing: boolean) {
  const choice: AppCookieConsentChoice = {
    version: 1,
    necessary: true,
    analyticsMarketing,
  };

  return JSON.stringify(choice);
}

function parseCookieConsent(value: string | null): { saved: boolean; analyticsMarketing: boolean } {
  if (!value) return { saved: false, analyticsMarketing: false };
  if (value === "accepted") return { saved: true, analyticsMarketing: true };
  if (value === "necessary") return { saved: true, analyticsMarketing: false };

  try {
    const parsed = JSON.parse(value) as Partial<AppCookieConsentChoice>;
    if (parsed.version === 1 && parsed.necessary === true) {
      return { saved: true, analyticsMarketing: parsed.analyticsMarketing === true };
    }
  } catch {
    return { saved: false, analyticsMarketing: false };
  }

  return { saved: false, analyticsMarketing: false };
}

function loadYandexMetrika(counterId: number) {
  if (typeof window === "undefined" || typeof document === "undefined") return;
  const existing = document.querySelector(`script[data-yandex-metrika-id="${counterId}"]`);
  if (existing) return;

  window.dataLayer = window.dataLayer || [];
  const metrika = window as Window & { ym?: MetrikaFunction };
  metrika.ym =
    metrika.ym ||
    function ymStub(...args: unknown[]) {
      const queue = (metrika.ym as MetrikaFunction);
      queue.a = queue.a || [];
      queue.a.push(args);
    };
  metrika.ym.l = Date.now();

  const script = document.createElement("script");
  script.async = true;
  script.src = `https://mc.yandex.ru/metrika/tag.js?id=${counterId}`;
  script.dataset.yandexMetrikaId = String(counterId);
  document.head.appendChild(script);

  metrika.ym(counterId, "init", {
    ssr: true,
    webvisor: true,
    clickmap: true,
    ecommerce: "dataLayer",
    referrer: document.referrer,
    url: location.href,
    accurateTrackBounce: true,
    trackLinks: true,
  });
}

function AppCookieConsentBanner({ language }: { language: string }) {
  const copy = useCallback((key: string, fallback: string) => appCopy(language, key, fallback), [language]);
  const privacyURL = legalDocumentURL("privacy.html", language);
  const consentURL = legalDocumentURL("consent.html", language);
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [analyticsMarketing, setAnalyticsMarketing] = useState(false);
  const [visible, setVisible] = useState(() => {
    if (typeof window === "undefined") return false;
    return !parseCookieConsent(localStorage.getItem(appCookieConsentStorageKey)).saved;
  });

  useEffect(() => {
    const consent = parseCookieConsent(localStorage.getItem(appCookieConsentStorageKey));
    setVisible(!consent.saved);
    if (consent.analyticsMarketing) loadYandexMetrika(appYandexMetrikaId);
  }, []);

  const saveChoice = (value: string, allowAnalytics: boolean) => {
    try {
      localStorage.setItem(appCookieConsentStorageKey, value);
    } catch {
      // The banner stays available even when storage is unavailable.
    }
    setVisible(false);
    if (allowAnalytics) loadYandexMetrika(appYandexMetrikaId);
  };

  if (!visible) return null;

  return (
    <section className="app-cookie-consent-banner" aria-label={copy("cookie_consent_title", "Cookie consent")}>
      <div>
        <strong>{copy("cookie_consent_title", "Cookie and technical data")}</strong>
        <p>{copy("cookie_consent_body", "We use necessary cookies and local storage for login, language, theme, security, saved consents, and stable app operation. Optional cookies are used only after consent.")}</p>
        <nav>
          <a href={privacyURL} target="_blank" rel="noreferrer">{copy("auth_privacy_policy", "Privacy policy")}</a>
          <a href={consentURL} target="_blank" rel="noreferrer">{copy("auth_personal_data_consent", "Personal data consent")}</a>
        </nav>
        {settingsOpen && (
          <div className="app-cookie-consent-banner__settings" aria-label={copy("cookie_settings", "Cookie settings")}>
            <div className="app-cookie-consent-banner__setting-row">
              <div>
                <strong>{copy("cookie_necessary_title", "Necessary")}</strong>
                <span>{copy("cookie_necessary_body", "Always active for login, security, language, and saved consents.")}</span>
              </div>
              <span>{copy("cookie_always_active", "Always active")}</span>
            </div>
            <div className="app-cookie-consent-banner__setting-row">
              <div>
                <strong>{copy("cookie_analytics_title", "Analytics and marketing")}</strong>
                <span>{copy("cookie_analytics_body", "Yandex Metrika and similar tools help measure visits and improve the app after consent.")}</span>
              </div>
              <button
                type="button"
                className="app-cookie-consent-banner__switch"
                role="switch"
                aria-checked={analyticsMarketing}
                aria-label={copy("cookie_analytics_title", "Analytics and marketing")}
                onClick={() => setAnalyticsMarketing((value) => !value)}
              >
                <span />
              </button>
            </div>
          </div>
        )}
      </div>
      <div className="app-cookie-consent-banner__actions">
        <button type="button" onClick={() => saveChoice(serializeAppCookieConsent(false), false)}>{copy("cookie_necessary", "Necessary only")}</button>
        <button type="button" data-app-cookie="settings" onClick={() => setSettingsOpen((value) => !value)}>{copy("cookie_settings", "Settings")}</button>
        {settingsOpen && <button type="button" onClick={() => saveChoice(serializeAppCookieConsent(analyticsMarketing), analyticsMarketing)}>{copy("cookie_save_selected", "Save choice")}</button>}
        <button type="button" onClick={() => saveChoice(serializeAppCookieConsent(true), true)}>{copy("cookie_accept", "Accept all cookies")}</button>
      </div>
    </section>
  );
}

function languageDisplayName(code: string | undefined, options?: LanguageOption[]) {
  const normalized = languageCode(code);
  const match = options?.find((option) => languageCode(option.code) === normalized);
  return match ? optionLabel(match) : normalized.toUpperCase();
}

function pronunciationLocale(user: UserProfile) {
  const lang = languageCode(user.interface_language);
  return pronunciationLocales[lang] || pronunciationLocales.en;
}

function cleanPronunciationWord(value: unknown) {
  return cleanAppText(asText(value, ""))
    .replace(/(?:low_confidence|substituted_or_unclear|missing)$/i, "")
    .replace(/[_-]+$/g, "")
    .trim();
}

function localizePronunciationCoachText(value: unknown, user: UserProfile, copy: (key: string, fallback: string) => string) {
  const text = cleanAppText(asText(value, "")).trim();
  if (!text || languageCode(user.interface_language) !== "ru") return text;
  const word = (match: RegExpMatchArray, index: number) => cleanAppText(match[index] || "").replace(/^['"]|['"]$/g, "");
  let match = text.match(/Your pronunciation is understandable/i);
  if (match) return copy("pronunciation_tip_understandable_low_confidence", "Произношение понятно, но оценка уверенности низкая. Говори чуть медленнее и четче.");
  match = text.match(/Focus on clear articulation of the word ['"]([^'"]+)['"]/i);
  if (match) return `${copy("pronunciation_tip_articulate_word", "Четко проговори слово")}: ${word(match, 1)}.`;
  match = text.match(/Ensure your voice volume is consistent/i);
  if (match) return copy("pronunciation_tip_consistent_volume", "Держи ровную громкость на протяжении всей фразы.");
  match = text.match(/Ensure the ['"]([^'"]+)['"] sound is clear and slightly elongated/i);
  if (match) return `${copy("pronunciation_tip_clear_sound", "Сделай звук четким и немного более протяжным")}: ${word(match, 1)}.`;
  match = text.match(/Focus on the stress:\s*(.+)$/i);
  if (match) return `${copy("pronunciation_tip_focus_stress", "Поставь ударение по схеме")}: ${cleanAppText(match[1])}`;
  match = text.match(/Make sure to pronounce ['"]([^'"]+)['"] clearly as a long ['"]ai['"] sound/i);
  if (match) return `${copy("pronunciation_tip_long_ai", "Произнеси четко как долгий звук ай")}: ${word(match, 1)}.`;
  match = text.match(/Try to speak at a steady pace/i);
  if (match) return copy("pronunciation_tip_steady_pace", "Говори ровнее и отделяй каждое слово.");
  match = text.match(/Pay attention to the ['"]([^'"]+)['"] sound at the end of ['"]([^'"]+)['"]/i);
  if (match) return `${copy("pronunciation_tip_ending_sound", "Обрати внимание на звук в конце слова")}: ${word(match, 1)} / ${word(match, 2)}.`;
  return text;
}

function pronunciationProblemText(item: PronunciationProblem, user: UserProfile, copy: (key: string, fallback: string) => string) {
  const issue = asText(item.issue, "").trim().toLowerCase();
  const tip = localizePronunciationCoachText(item.tip || "", user, copy);
  const ru = languageCode(user.interface_language) === "ru";
  if (tip && !/low_confidence|substituted_or_unclear|missing/i.test(tip)) return tip;
  if (issue.includes("missing")) {
    return ru
      ? copy("pronunciation_missing_word", "Слово пропущено или не распознано. Скажи его отдельно, затем повтори всю фразу.")
      : copy("pronunciation_missing_word", "The word was missed or not recognized. Say it alone, then repeat the full phrase.");
  }
  if (issue.includes("substituted") || issue.includes("unclear")) {
    const spoken = cleanPronunciationWord(item.spoken);
    return spoken
      ? (ru
        ? copy("pronunciation_substituted_word", "Ожидалось другое слово. Сравни с услышанным вариантом и повтори медленнее.") + ` ${copy("heard", "Услышано")}: ${spoken}.`
        : copy("pronunciation_substituted_word", "A different word was recognized. Compare it with what was heard and repeat slower.") + ` ${copy("heard", "Heard")}: ${spoken}.`)
      : (ru
        ? copy("pronunciation_unclear_word", "Слово прозвучало неясно. Повтори его с ровным темпом и четким окончанием.")
        : copy("pronunciation_unclear_word", "The word sounded unclear. Repeat it with steady rhythm and a clear ending."));
  }
  if (issue.includes("low") || issue.includes("confidence")) {
    return ru
      ? copy("pronunciation_low_confidence", "Распознавание не уверено в этом слове. Выдели ударение и не проглатывай окончание.")
      : copy("pronunciation_low_confidence", "Recognition was unsure about this word. Stress it clearly and avoid swallowing the ending.");
  }
  return tip || (ru
    ? copy("pronunciation_practice_word", "Потренируй слово отдельно и затем произнеси всю фразу.")
    : copy("pronunciation_practice_word", "Practice the word alone, then say the full phrase."));
}

function pronunciationProblemLabel(item: PronunciationProblem, copy: (key: string, fallback: string) => string) {
  const word = cleanPronunciationWord(item.word) || cleanPronunciationWord(item.spoken) || copy("phrase", "Phrase");
  const confidence = Number(item.confidence);
  return {
    word,
    confidence: Number.isFinite(confidence) && confidence > 0 ? `${Math.round(confidence * 100)}%` : "",
  };
}

function uniquePronunciationProblems(items: PronunciationProblem[]) {
  const seen = new Set<string>();
  const out: PronunciationProblem[] = [];
  for (const item of items) {
    const key = cleanPronunciationWord(item.word || item.spoken).toLowerCase();
    if (!key || seen.has(key)) continue;
    seen.add(key);
    out.push(item);
    if (out.length >= 12) break;
  }
  return out;
}

function compactPronunciationMap(items: PronunciationProblem[], user: UserProfile, copy: (key: string, fallback: string) => string) {
  const seenText = new Set<string>();
  const compact: PronunciationProblem[] = [];
  for (const item of items) {
    const label = cleanPronunciationWord(item.word || item.spoken);
    if (!label) continue;
    const issue = asText(item.issue, "").toLowerCase();
    const tip = pronunciationProblemText(item, user, copy)
      .replace(/\s+/g, " ")
      .trim();
    const signal = issue.includes("missing")
      ? copy("pronunciation_map_missing", "Пропуск")
      : issue.includes("substituted") || issue.includes("unclear")
        ? copy("pronunciation_map_unclear", "Нечетко")
        : issue.includes("low") || issue.includes("confidence")
          ? copy("pronunciation_map_confidence", "Уверенность")
          : tip || copy("pronunciation_map_focus", "Фокус");
    const normalized = `${label.toLowerCase()}|${signal.toLowerCase()}`;
    if (seenText.has(normalized)) continue;
    seenText.add(normalized);
    compact.push({ ...item, issue: signal, tip: tip && tip !== signal ? tip : "" });
    if (compact.length >= 8) break;
  }
  return compact;
}

function pronunciationFocusItems(items: PronunciationProblem[], copy: (key: string, fallback: string) => string) {
  const focus = uniquePronunciationProblems(items).slice(0, 6).map((item) => {
    const word = cleanPronunciationWord(item.word || item.spoken);
    const lower = word.toLowerCase();
    const target = lower.includes("th")
      ? "th"
      : lower.match(/[rl]/)
        ? "r / l"
        : lower.match(/[vw]/)
          ? "v / w"
          : lower.length > 4
            ? copy("word_ending", "окончание")
            : copy("word_stress", "ударение");
    return { word, target };
  });
  return focus.length ? focus : [{ word: copy("no_voice_data", "Пока нет голосовых данных"), target: copy("start_pronunciation_hint", "запишите фразу выше") }];
}

function pronunciationPracticeFallback(user: UserProfile, copy: (key: string, fallback: string) => string, rotate = false) {
  const phrases = [
    copy("pronunciation_sample_1", "Could you say that a little slower, please?"),
    copy("pronunciation_sample_2", "I would like a quiet room for one night."),
    copy("pronunciation_sample_3", "The reservation is under my name."),
    copy("pronunciation_sample_4", "Can you help me with this sentence?"),
  ];
  const seed = rotate ? Date.now() : tutorStableHash(`${user.learning_language || "en"}:${user.level || "A1"}`);
  return phrases[Math.abs(seed) % phrases.length];
}

function pronunciationFallbackForKey(interfaceLanguage: string | undefined, key: string) {
  const p = pronunciationLocales[languageCode(interfaceLanguage)] || pronunciationLocales.en;
  const map: Record<string, string> = {
    pronunciation: p.title,
    pronunciation_dashboard: p.title,
    pronunciation_dashboard_body: p.body,
    pronunciation_heatmap: p.heatmap,
    weak_words: p.weak,
    score_history: p.history,
    no_voice_data: p.noData,
    start_listening_hint: p.hint,
    latest_score: p.latest,
    try_again: p.repeat,
  };
  return map[key];
}

function roleplayFallbackForKey(interfaceLanguage: string | undefined, key: string) {
  const labels = roleplayUiFallbacks[languageCode(interfaceLanguage)] || roleplayUiFallbacks.en;
  const map: Record<string, string> = {
    roleplay_title: labels.title,
    roleplay_subtitle: labels.subtitle,
    roleplay_result: labels.result,
  };
  return map[key];
}

function roleplayScenarioTitle(scenario: RoleplayScenario, user: UserProfile) {
  const lang = languageCode(user.interface_language);
  const fallbackLang = ["tg", "uz", "tt", "ky", "uk", "kk"].includes(lang) ? "ru" : "en";
  const genericRoleplay = cleanAppText(appCopy(lang, "roleplay", "Roleplay"));
  const copyCandidate = cleanAppText(appCopy(lang, `roleplay_scenario_${scenario.id}`, ""));
  if (copyCandidate && copyCandidate !== `roleplay_scenario_${scenario.id}` && copyCandidate !== genericRoleplay && !hasMojibakeText(copyCandidate)) return copyCandidate;
  const concrete = cleanAppText(concreteRoleplayTitles[lang]?.[scenario.id] || "");
  if (concrete) return concrete;
  const localized = cleanAppText(roleplayTitles[scenario.id]?.[lang] || "");
  const english = cleanAppText(roleplayTitles[scenario.id]?.en || scenario.title);
  if (localized && !hasMojibakeText(localized) && localized !== genericRoleplay) return localized;
  return cleanAppText(roleplayTitles[scenario.id]?.[fallbackLang] || english || scenario.title);
}

function roleplayScenarioDescription(scenario: RoleplayScenario, user: UserProfile) {
  const lang = languageCode(user.interface_language);
  const genericRoleplay = cleanAppText(appCopy(lang, "roleplay", "Roleplay"));
  const specificCopy = cleanAppText(appCopy(lang, `roleplay_scenario_${scenario.id}_description`, ""));
  if (specificCopy && specificCopy !== `roleplay_scenario_${scenario.id}_description` && specificCopy !== genericRoleplay && !hasMojibakeText(specificCopy)) return specificCopy;
  const localized = cleanAppText(roleplayDescriptions[lang]?.[scenario.id] || roleplayGenericDescriptions[lang] || "");
  if (localized && localized !== genericRoleplay && !hasMojibakeText(localized)) return localized;
  const copyCandidate = cleanAppText(appCopy(lang, "roleplay_scenario_description", ""));
  if (copyCandidate && copyCandidate !== "roleplay_scenario_description" && copyCandidate !== genericRoleplay && !hasMojibakeText(copyCandidate)) return copyCandidate;
  return cleanAppText(scenario.description || roleplayDescriptions.en?.[scenario.id] || roleplayGenericDescriptions.en);
}

function leaderboardLanguages(entry: LeaderboardEntry, languages: LanguageOption[], copy: (key: string, fallback: string) => string) {
  const codes = Array.isArray(entry.languages) ? entry.languages.filter(Boolean) : [];
  if (!codes.length) return copy("languages_not_listed", "Languages not listed");
  const names = codes.map((code) => languageDisplayName(code, languages));
  return `${copy("languages", "Languages")}: ${names.join(", ")}`;
}

function buildRoleplayScenarioPrompt(
  scenario: RoleplayScenario,
  user: UserProfile,
  session: SessionData,
  copy: (key: string, fallback: string) => string,
  learnerLine = "",
) {
  const interfaceName = languageDisplayName(user.interface_language || "ru", session.interface_languages);
  const learningName = languageDisplayName(user.learning_language || "en", session.learning_languages);
  const title = roleplayScenarioTitle(scenario, user);
  const description = roleplayScenarioDescription(scenario, user);
  const goal = scenario.goal || title;
  const learnerRole = scenario.learnerRole || copy("learner", "learner");
  const aiRole = scenario.aiRole || title;
  return [
    "ROLEPLAY_TOOL_V2",
    `Interface language: ${interfaceName}`,
    `Target learning language: ${learningName}`,
    `CEFR level: ${user.level || "A1"}`,
    `Scenario: ${title}`,
    `Scenario description: ${description}`,
    `Learner role: ${learnerRole}`,
    `AI role: ${aiRole}`,
    `Goal: ${goal}`,
    learnerLine ? `Learner line: ${learnerLine}` : "Learner line: <none yet>",
    "",
    "Run this as a roleplay, not as a normal practice chat.",
    `Use ${interfaceName} for setup, task text, and short correction notes.`,
    `Use ${learningName} only for the character dialogue, model phrase, corrected phrase, and the next character line.`,
    learnerLine
      ? "The learner has already answered. First give one compact correction if needed, then continue the scene with one natural character line/question. Do not show a generic practice prompt."
      : "This is the first turn. Start with a short scene setup, one AI character line in the target language, the learner's task, and one optional example answer. Do not ask a separate 'next question' block before the learner answers.",
    "Do not invent corrections. If you name a problem, quote the exact learner fragment and the exact replacement. If a word or preposition is already present, do not say to add it.",
    "Keep labels learner-friendly. Do not write raw prompt names, system words, or technical wording. Keep it concise.",
  ].join("\n");
}

function AnimatedSelect({
  value,
  options,
  onChange,
  icon,
  label,
  ariaLabel,
  className,
  disabled,
}: {
  value: string;
  options: Array<LanguageOption | { code: string; name?: string; native_name?: string }>;
  onChange: (value: string) => void;
  icon?: ReactNode;
  label?: string;
  ariaLabel?: string;
  className?: string;
  disabled?: boolean;
}) {
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement | null>(null);
  const selected = options.find((option) => option.code === value) || options[0];

  useEffect(() => {
    if (!open) return;
    const close = (event: MouseEvent) => {
      if (!rootRef.current?.contains(event.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, [open]);

  return (
    <div ref={rootRef} className={cn("animated-select-v2", open && "is-open", disabled && "is-disabled", className)}>
      {label ? <span className="animated-select-v2__label">{label}</span> : null}
      <button
        type="button"
        aria-label={ariaLabel || label}
        aria-expanded={open}
        disabled={disabled}
        onClick={() => setOpen((current) => !current)}
      >
        {icon}
        <strong>{selected ? optionLabel(selected) : value}</strong>
        <ChevronDown size={15} />
      </button>
      <div className="animated-select-v2__menu" role="listbox">
        {options.map((option) => (
          <button
            key={option.code}
            type="button"
            role="option"
            aria-selected={option.code === value}
            className={option.code === value ? "is-selected" : ""}
            onClick={() => {
              onChange(option.code);
              setOpen(false);
            }}
          >
            <span>{optionLabel(option)}</span>
            {option.code === value ? <Check size={14} /> : null}
          </button>
        ))}
      </div>
    </div>
  );
}

function TopBar({
  user,
  session,
  theme,
  setTheme,
  brightness,
  setBrightness,
  copy,
  onLogout,
  onLanguageChange,
  setView,
  busy,
  onBugReport,
  onGuide,
}: {
  user: UserProfile;
  session: SessionData;
  theme: Theme;
  setTheme: (theme: Theme) => void;
  brightness: number;
  setBrightness: (value: number) => void;
  copy: (key: string, fallback: string) => string;
  onLogout: () => void;
  onLanguageChange: (code: string) => Promise<void>;
  setView: (view: ViewId) => void;
  busy: string | null;
  onBugReport: () => void;
  onGuide: () => void;
}) {
  const interfaceLanguages = session.interface_languages?.length ? session.interface_languages : [{ code: user.interface_language || "ru", native_name: user.interface_language || "ru" }];
  const accountLogin = cleanAppText(session.account?.login) || cleanAppText(user.telegram_account?.name) || copy("learner", "Learner");
  return (
    <header className="v2-topbar">
      <button className="v2-brand-button" type="button" onClick={() => setView("settings")} aria-label={copy("settings", "Settings")}>
        <BrandMark subtitle={accountLogin} />
      </button>
      <div className="v2-topbar__profile">
        <button className="v2-level-button" type="button" onClick={() => setView("awards")} aria-label={copy("awards", "Awards")}>
          <LevelProgress user={user} copy={copy} />
        </button>
      </div>
      <div className="v2-topbar__actions">
        <button className="app-guide-button-v2" type="button" onClick={onGuide} aria-label={copy("app_guide_title", "Quick start guide")}>
          <CircleHelp size={15} />
          <span>{copy("guide", "Guide")}</span>
        </button>
        <button className="v2-report-button" type="button" onClick={onBugReport}>
          <Bug size={15} />
          <span>{copy("report_bug", "Сообщить об ошибке")}</span>
        </button>
        <label className="v2-brightness-control" aria-label={copy("brightness", "Brightness")}>
          <Sparkles size={14} />
          <Slider value={[brightness]} min={65} max={125} step={1} onValueChange={(value) => setBrightness(value[0] || 100)} />
        </label>
        <AnimatedSelect
          className="v2-language-select"
          value={user.interface_language || "ru"}
          options={interfaceLanguages}
          icon={<Languages size={15} />}
          ariaLabel={copy("interface_language", "Interface language")}
          disabled={busy === "settings"}
          onChange={(code) => void onLanguageChange(code)}
        />
        <AnimatedThemeToggle theme={theme} onChange={setTheme} />
        <LogoutButton onConfirm={onLogout} disabled={busy === "logout"} label={copy("logout", "Log out")} cancelLabel={copy("cancel", "Cancel")} />
      </div>
    </header>
  );
}

function MobileQuickControls({
  user,
  session,
  theme,
  setTheme,
  copy,
  onLogout,
  onLanguageChange,
  onBugReport,
  onGuide,
  busy,
}: {
  user: UserProfile;
  session: SessionData;
  theme: Theme;
  setTheme: (theme: Theme) => void;
  copy: (key: string, fallback: string) => string;
  onLogout: () => void;
  onLanguageChange: (code: string) => Promise<void>;
  onBugReport: () => void;
  onGuide: () => void;
  busy: string | null;
}) {
  const interfaceLanguages = session.interface_languages?.length ? session.interface_languages : [{ code: user.interface_language || "ru", native_name: user.interface_language || "ru" }];
  return (
    <div className="mobile-quick-controls-v2" aria-label={copy("mobile_quick_controls", "Mobile quick controls")}>
      <button className="app-guide-button-v2" type="button" onClick={onGuide} aria-label={copy("app_guide_title", "Quick start guide")}>
        <CircleHelp size={17} />
      </button>
      <button className="mobile-report-button-v2" type="button" onClick={onBugReport} aria-label={copy("report_bug", "Сообщить об ошибке")}>
        <Bug size={17} />
      </button>
      <AnimatedSelect
        className="mobile-language-select-v2"
        value={user.interface_language || "ru"}
        options={interfaceLanguages}
        icon={<Languages size={16} />}
        ariaLabel={copy("interface_language", "Interface language")}
        disabled={busy === "settings"}
        onChange={(code) => void onLanguageChange(code)}
      />
      <AnimatedThemeToggle className="mobile-theme-toggle-v2" theme={theme} onChange={setTheme} />
      <LogoutButton onConfirm={onLogout} disabled={busy === "logout"} label={copy("logout", "Log out")} cancelLabel={copy("cancel", "Cancel")} iconOnly />
    </div>
  );
}

function LevelProgress({ user, copy }: { user: UserProfile; copy: (key: string, fallback: string) => string }) {
  const current = Number(user.xp_current ?? user.xp ?? 0);
  const needed = Number(user.xp_needed || 100);
  const pct = percent(current, needed);
  const currentAward = awardLevel(user);
  return (
    <div className="v2-level-progress">
      <span className="v2-current-award" title={`${copy("current_award", "Current award")} ${currentAward}`}>
        <img src={awardFile(currentAward)} alt="" />
      </span>
      <div>
        <strong>{user.level || "A1"} - LVL {user.xp_level || 1}</strong>
        <span>{user.xp_title || copy("current_award", "Current award")}</span>
      </div>
      <i aria-label={`${copy("progress", "Progress")} ${pct}%`}><b style={{ width: `${pct}%` }} /></i>
      <small>{formatLimit(current, needed)}</small>
    </div>
  );
}

function FunctionRibbon({
  nav,
  activeView,
  setView,
  theme,
  copy,
  accountKey,
  navigationLayout,
  onNavigationLayoutChange,
}: {
  nav: NavItem[];
  activeView: ViewId;
  setView: (view: ViewId) => void;
  theme: Theme;
  copy: (key: string, fallback: string) => string;
  accountKey: string;
  navigationLayout: NonNullable<UserProfile["navigation_layout"]>;
  onNavigationLayoutChange: (patch: NonNullable<UserProfile["navigation_layout"]>) => void;
}) {
  const ribbonRef = useRef<HTMLElement | null>(null);
  const [draggedId, setDraggedId] = useState<ViewId | null>(null);
  const [dragOverId, setDragOverId] = useState<ViewId | null>(null);
  const [dragAfterTarget, setDragAfterTarget] = useState(false);
  const orderStorageKey = `poliglot-function-ribbon-v2:${accountKey || "guest"}`;
  const readOrder = useCallback(() => {
    const serverOrder = navigationLayout.function_ribbon;
    if (Array.isArray(serverOrder) && serverOrder.length) {
      const ids = serverOrder.filter((id): id is ViewId => nav.some((item) => item.id === id));
      if (ids.length) return ids;
    }
    try {
      const parsed = JSON.parse(localStorage.getItem(orderStorageKey) || "[]");
      if (!Array.isArray(parsed)) return [] as ViewId[];
      return parsed.filter((id): id is ViewId => nav.some((item) => item.id === id));
    } catch {
      return [] as ViewId[];
    }
  }, [nav, navigationLayout.function_ribbon, orderStorageKey]);
  const [orderIds, setOrderIds] = useState<ViewId[]>(readOrder);
  useEffect(() => {
    setOrderIds(readOrder());
  }, [readOrder]);
  const orderedNav = useMemo(() => {
    const ids = [...orderIds, ...nav.map((item) => item.id)].filter((id, index, list) => list.indexOf(id) === index && nav.some((item) => item.id === id));
    return ids.map((id) => nav.find((item) => item.id === id)).filter(Boolean) as NavItem[];
  }, [nav, orderIds]);
  const persistOrder = (next: ViewId[]) => {
    const normalized = next.filter((id, index, list) => nav.some((item) => item.id === id) && list.indexOf(id) === index);
    setOrderIds(normalized);
    localStorage.setItem(orderStorageKey, JSON.stringify(normalized));
    onNavigationLayoutChange({ function_ribbon: normalized });
  };
  const reorderRibbon = (target: ViewId, after = false) => {
    if (!draggedId || draggedId === target) return;
    const current = orderedNav.map((item) => item.id).filter((id) => id !== draggedId);
    const targetIndex = current.indexOf(target);
    current.splice(targetIndex < 0 ? current.length : targetIndex + (after ? 1 : 0), 0, draggedId);
    persistOrder(current);
  };
  const updateRibbonDropTarget = (target: ViewId, clientX: number, rect: DOMRect) => {
    if (!draggedId) return;
    setDragOverId(target);
    setDragAfterTarget(clientX > rect.left + rect.width / 2);
    const node = ribbonRef.current;
    if (!node) return;
    const rail = node.getBoundingClientRect();
    if (clientX < rail.left + 74) node.scrollBy({ left: -32, behavior: "auto" });
    if (clientX > rail.right - 74) node.scrollBy({ left: 32, behavior: "auto" });
  };
  const scrollByTile = (direction: -1 | 1) => {
    ribbonRef.current?.scrollBy({ left: direction * 360, behavior: "smooth" });
  };
  return (
    <div className="function-ribbon-shell">
      <MorphingArrowButton direction="left" onClick={() => scrollByTile(-1)} label={copy("back", "Back")} className="ribbon-morph-arrow" size="compact" />
      <nav
        ref={ribbonRef}
        className="function-ribbon"
        aria-label="Web function ribbon"
        onWheel={(event) => {
          if (Math.abs(event.deltaY) > Math.abs(event.deltaX)) {
            event.currentTarget.scrollLeft += event.deltaY;
          }
        }}
      >
        {orderedNav.map((item) => (
          <button
            key={item.id}
            type="button"
            className={cn(
              "function-chip",
              activeView === item.id && "is-active",
              draggedId === item.id && "is-dragging",
              dragOverId === item.id && draggedId !== item.id && "is-drop-target",
              dragOverId === item.id && draggedId !== item.id && !dragAfterTarget && "is-drop-before",
              dragOverId === item.id && draggedId !== item.id && dragAfterTarget && "is-drop-after",
            )}
            data-accent={item.accent}
            data-view={item.id}
            draggable
            aria-pressed={activeView === item.id}
            onDragStart={(event) => {
              setDraggedId(item.id);
              event.dataTransfer.effectAllowed = "move";
              event.dataTransfer.setData("text/plain", item.id);
            }}
            onDragOver={(event) => {
              event.preventDefault();
              updateRibbonDropTarget(item.id, event.clientX, event.currentTarget.getBoundingClientRect());
            }}
            onDrop={(event) => {
              event.preventDefault();
              const rect = event.currentTarget.getBoundingClientRect();
              reorderRibbon(item.id, event.clientX > rect.left + rect.width / 2);
              setDraggedId(null);
              setDragOverId(null);
              setDragAfterTarget(false);
            }}
            onDragEnd={() => {
              if (dragOverId) reorderRibbon(dragOverId, dragAfterTarget);
              setDraggedId(null);
              setDragOverId(null);
              setDragAfterTarget(false);
            }}
            onClick={() => {
              if (draggedId) return;
              setView(item.id);
            }}
          >
            <span className="function-chip__art" style={{ backgroundImage: `url("${menuAssetUrl(item.id, theme)}")` }} />
            <span className="function-chip__caption">
              <strong>{item.shortLabel || item.label}</strong>
            </span>
          </button>
        ))}
      </nav>
      <MorphingArrowButton direction="right" onClick={() => scrollByTile(1)} label={copy("next", "Next")} className="ribbon-morph-arrow" size="compact" />
    </div>
  );
}

function MobileBottomNav({
  nav,
  activeView,
  setView,
  theme,
  copy,
  accountKey,
  navigationLayout,
  onNavigationLayoutChange,
}: {
  nav: NavItem[];
  activeView: ViewId;
  setView: (view: ViewId) => void;
  theme: Theme;
  copy: (key: string, fallback: string) => string;
  accountKey: string;
  navigationLayout: NonNullable<UserProfile["navigation_layout"]>;
  onNavigationLayoutChange: (patch: NonNullable<UserProfile["navigation_layout"]>) => void;
}) {
  const [editMode, setEditMode] = useState(false);
  const [draggedId, setDraggedId] = useState<ViewId | null>(null);
  const [dragOverId, setDragOverId] = useState<ViewId | null>(null);
  const [dragAfterTarget, setDragAfterTarget] = useState(false);
  const bottomNavRef = useRef<HTMLElement | null>(null);
  const draggedIdRef = useRef<ViewId | null>(null);
  const dragOverIdRef = useRef<ViewId | null>(null);
  const dragAfterTargetRef = useRef(false);
  const longPressTimerRef = useRef<number | null>(null);
  const longPressTriggeredRef = useRef(false);
  const pointerStartRef = useRef<{ id: ViewId; x: number; y: number } | null>(null);
  const dragGhostRef = useRef<HTMLDivElement | null>(null);
  const dragGhostPointRef = useRef<{ x: number; y: number }>({ x: 0, y: 0 });
  const dragGhostFrameRef = useRef<number | null>(null);
  const storageKey = `poliglot-mobile-nav-v2:${accountKey || "guest"}`;
  const defaultPinned: ViewId[] = ["home", "tutor", "words", "word-game", "pronunciation"];
  const longPressEditDelayMs = 650;
  const readPinned = useCallback(() => {
    const serverPinned = navigationLayout.mobile_pinned;
    if (Array.isArray(serverPinned) && serverPinned.length) {
      const ids = serverPinned.filter((id): id is ViewId => nav.some((item) => item.id === id));
      return ids.length ? ids.slice(0, 5) : defaultPinned;
    }
    try {
      const parsed = JSON.parse(localStorage.getItem(storageKey) || "[]");
      if (!Array.isArray(parsed)) return defaultPinned;
      const ids = parsed.filter((id): id is ViewId => nav.some((item) => item.id === id));
      return ids.length ? ids.slice(0, 5) : defaultPinned;
    } catch {
      return defaultPinned;
    }
  }, [storageKey, nav, navigationLayout.mobile_pinned]);
  const [pinnedIds, setPinnedIds] = useState<ViewId[]>(readPinned);
  const moreStorageKey = `poliglot-mobile-nav-more-v2:${accountKey || "guest"}`;
  const readMoreOrder = useCallback(() => {
    const serverMore = navigationLayout.mobile_more;
    if (Array.isArray(serverMore) && serverMore.length) {
      return serverMore.filter((id): id is ViewId => nav.some((item) => item.id === id));
    }
    try {
      const parsed = JSON.parse(localStorage.getItem(moreStorageKey) || "[]");
      if (!Array.isArray(parsed)) return [] as ViewId[];
      return parsed.filter((id): id is ViewId => nav.some((item) => item.id === id));
    } catch {
      return [] as ViewId[];
    }
  }, [moreStorageKey, nav, navigationLayout.mobile_more]);
  const [moreOrderIds, setMoreOrderIds] = useState<ViewId[]>(readMoreOrder);
  const railStorageKey = `poliglot-mobile-nav-rail-v2:${accountKey || "guest"}`;
  const readRailOrder = useCallback(() => {
    const serverRail = navigationLayout.mobile_rail;
    if (Array.isArray(serverRail) && serverRail.length) {
      return serverRail.filter((id): id is ViewId => nav.some((item) => item.id === id));
    }
    try {
      const parsed = JSON.parse(localStorage.getItem(railStorageKey) || "[]");
      if (!Array.isArray(parsed)) return [] as ViewId[];
      return parsed.filter((id): id is ViewId => nav.some((item) => item.id === id));
    } catch {
      return [] as ViewId[];
    }
  }, [railStorageKey, nav, navigationLayout.mobile_rail]);
  const [railOrderIds, setRailOrderIds] = useState<ViewId[]>(readRailOrder);
  const [draftRailOrderIds, setDraftRailOrderIds] = useState<ViewId[] | null>(null);
  const draftRailOrderIdsRef = useRef<ViewId[] | null>(null);
  useEffect(() => {
    setPinnedIds(readPinned());
    setMoreOrderIds(readMoreOrder());
    setRailOrderIds(readRailOrder());
  }, [readPinned, readMoreOrder, readRailOrder]);
  useEffect(() => {
    setDraftRailOrderIds(null);
    draftRailOrderIdsRef.current = null;
  }, [railStorageKey]);
  const normalizeRailOrder = (next: ViewId[]) => {
    return next.filter((id, index, list) => nav.some((item) => item.id === id) && list.indexOf(id) === index);
  };
  const persistRailOrder = (next: ViewId[]) => {
    const normalized = normalizeRailOrder(next);
    draftRailOrderIdsRef.current = null;
    setDraftRailOrderIds(null);
    setRailOrderIds((current) => current.length === normalized.length && current.every((id, index) => id === normalized[index]) ? current : normalized);
    localStorage.setItem(railStorageKey, JSON.stringify(normalized));
    onNavigationLayoutChange({ mobile_pinned: pinnedIds, mobile_more: moreOrderIds, mobile_rail: normalized });
  };
  const setActiveDrag = (view: ViewId | null) => {
    draggedIdRef.current = view;
    setDraggedId(view);
  };
  const updateDragGhost = (x: number, y: number) => {
    dragGhostPointRef.current = { x, y };
    if (dragGhostFrameRef.current !== null) return;
    dragGhostFrameRef.current = window.requestAnimationFrame(() => {
      dragGhostFrameRef.current = null;
      const point = dragGhostPointRef.current;
      if (dragGhostRef.current) {
        dragGhostRef.current.style.transform = `translate3d(${point.x}px, ${point.y}px, 0) translate(-50%, -50%)`;
      }
    });
  };
  const enterEditMode = (view?: ViewId) => {
    longPressTriggeredRef.current = true;
    longPressTimerRef.current = null;
    setEditMode(true);
    if (view) {
      setActiveDrag(view);
      setDropTarget(view, false);
      if (pointerStartRef.current) updateDragGhost(pointerStartRef.current.x, pointerStartRef.current.y);
    }
  };
  const startLongPress = (view?: ViewId) => {
    longPressTriggeredRef.current = false;
    if (longPressTimerRef.current) window.clearTimeout(longPressTimerRef.current);
    longPressTimerRef.current = window.setTimeout(() => enterEditMode(view), longPressEditDelayMs);
  };
  const cancelLongPress = () => {
    if (longPressTimerRef.current) window.clearTimeout(longPressTimerRef.current);
    longPressTimerRef.current = null;
  };
  const resetPointerStart = () => {
    pointerStartRef.current = null;
  };
  const setDropTarget = (id: ViewId | null, after = false) => {
    if (dragOverIdRef.current === id && dragAfterTargetRef.current === after) return;
    dragOverIdRef.current = id;
    dragAfterTargetRef.current = after;
    setDragOverId(id);
    setDragAfterTarget(after);
  };
  const railNav = useMemo(() => {
    const orderedIds = [
      ...(draftRailOrderIds ?? railOrderIds),
      ...pinnedIds,
      ...moreOrderIds,
      ...nav.map((item) => item.id),
    ].filter((id, index, list) => list.indexOf(id) === index && nav.some((item) => item.id === id));
    return orderedIds.map((id) => nav.find((item) => item.id === id)).filter(Boolean) as NavItem[];
  }, [draftRailOrderIds, moreOrderIds, nav, pinnedIds, railOrderIds]);
  const railNavRef = useRef<NavItem[]>(railNav);
  useEffect(() => {
    railNavRef.current = railNav;
  }, [railNav]);
  const chooseView = (view: ViewId) => {
    if (editMode) return;
    setView(view);
  };
  const reorderedRailIds = (target: ViewId, after = false) => {
    const activeDrag = draggedIdRef.current || draggedId;
    if (!activeDrag || activeDrag === target) return null;
    const currentSource = draftRailOrderIdsRef.current ?? railNavRef.current.map((item) => item.id);
    const current = currentSource.filter((id) => id !== activeDrag);
    const targetIndex = current.indexOf(target);
    current.splice(targetIndex < 0 ? current.length : targetIndex + (after ? 1 : 0), 0, activeDrag);
    return normalizeRailOrder(current);
  };
  const previewRailReorder = (target: ViewId, after = false) => {
    const next = reorderedRailIds(target, after);
    if (!next) return;
    const current = draftRailOrderIdsRef.current ?? railNavRef.current.map((item) => item.id);
    if (current.length === next.length && current.every((id, index) => id === next[index])) return;
    draftRailOrderIdsRef.current = next;
    setDraftRailOrderIds(next);
    railNavRef.current = next.map((id) => nav.find((item) => item.id === id)).filter(Boolean) as NavItem[];
  };
  const commitRailReorder = () => {
    const draft = draftRailOrderIdsRef.current;
    if (draft?.length) {
      persistRailOrder(draft);
      return;
    }
    if (dragOverIdRef.current) {
      const next = reorderedRailIds(dragOverIdRef.current, dragAfterTargetRef.current);
      if (next) persistRailOrder(next);
    }
  };
  const navDropTargetFromPoint = (clientX: number, clientY: number) => {
    const activeDrag = draggedIdRef.current || draggedId;
    const node = bottomNavRef.current;
    if (!node || !activeDrag) return null;
    const buttons = Array.from(node.querySelectorAll<HTMLButtonElement>("button[data-view]"))
      .filter((button) => button.dataset.view && button.dataset.view !== activeDrag);
    if (!buttons.length) return null;
    let best: { id: ViewId; after: boolean; distance: number } | null = null;
    for (const button of buttons) {
      const view = button.dataset.view as ViewId;
      if (!nav.some((item) => item.id === view)) continue;
      const rect = button.getBoundingClientRect();
      const centerX = rect.left + rect.width / 2;
      const centerY = rect.top + rect.height / 2;
      const dx = clientX - centerX;
      const dy = clientY - centerY;
      const distance = Math.hypot(dx, dy);
      const after = Math.abs(dy) > rect.height * 0.55 ? clientY > centerY : clientX > centerX;
      const candidate = { id: view, after, distance };
      if (!best || candidate.distance < best.distance) best = candidate;
    }
    return best ? { id: best.id, after: best.after } : null;
  };
  const navDropTargetFromPointer = (event: ReactPointerEvent<HTMLElement>) => navDropTargetFromPoint(event.clientX, event.clientY);
  const autoScrollBottomNav = (clientX: number) => {
    const node = bottomNavRef.current;
    if (!node || node.scrollWidth <= node.clientWidth) return;
    const rect = node.getBoundingClientRect();
    const edge = 42;
    const maxStep = 8;
    const leftPressure = Math.max(0, rect.left + edge - clientX) / edge;
    const rightPressure = Math.max(0, clientX - (rect.right - edge)) / edge;
    if (leftPressure > 0) node.scrollBy({ left: -Math.max(2, Math.round(maxStep * Math.min(1, leftPressure))), behavior: "auto" });
    if (rightPressure > 0) node.scrollBy({ left: Math.max(2, Math.round(maxStep * Math.min(1, rightPressure))), behavior: "auto" });
  };
  const handleBottomPointerMove = (event: ReactPointerEvent<HTMLElement>) => {
    const activeDrag = draggedIdRef.current || draggedId;
    const pointerStart = pointerStartRef.current;
    if (!editMode && pointerStart) {
      const dx = event.clientX - pointerStart.x;
      const dy = event.clientY - pointerStart.y;
      if (Math.hypot(dx, dy) > 12) {
        cancelLongPress();
        return;
      }
    }
    if (!editMode || !activeDrag) return;
    event.preventDefault();
    updateDragGhost(event.clientX, event.clientY);
    autoScrollBottomNav(event.clientX);
    const next = navDropTargetFromPointer(event);
    if (next) {
      setDropTarget(next.id, next.after);
      previewRailReorder(next.id, next.after);
    }
  };
  const handleBottomPointerUp = (event?: ReactPointerEvent<HTMLElement>) => {
    cancelLongPress();
    resetPointerStart();
    if (!editMode) {
      setActiveDrag(null);
      setDropTarget(null, false);
      return;
    }
    event?.preventDefault();
    const activeDrag = draggedIdRef.current || draggedId;
    if (editMode && activeDrag) commitRailReorder();
    setActiveDrag(null);
    setDropTarget(null, false);
    longPressTriggeredRef.current = false;
  };
  useEffect(() => {
    if (!editMode) return;
    const handleMove = (event: PointerEvent) => {
      const activeDrag = draggedIdRef.current;
      if (!activeDrag) return;
      event.preventDefault();
      updateDragGhost(event.clientX, event.clientY);
      autoScrollBottomNav(event.clientX);
      const next = navDropTargetFromPoint(event.clientX, event.clientY);
      if (!next) return;
      setDropTarget(next.id, next.after);
      previewRailReorder(next.id, next.after);
    };
    const handleUp = (event: PointerEvent) => {
      event.preventDefault();
      const activeDrag = draggedIdRef.current;
      if (activeDrag) commitRailReorder();
      setActiveDrag(null);
      setDropTarget(null, false);
      cancelLongPress();
      resetPointerStart();
      longPressTriggeredRef.current = false;
    };
    window.addEventListener("pointermove", handleMove, { passive: false });
    window.addEventListener("pointerup", handleUp, { passive: false });
    window.addEventListener("pointercancel", handleUp, { passive: false });
    return () => {
      window.removeEventListener("pointermove", handleMove);
      window.removeEventListener("pointerup", handleUp);
      window.removeEventListener("pointercancel", handleUp);
      if (dragGhostFrameRef.current !== null) {
        window.cancelAnimationFrame(dragGhostFrameRef.current);
        dragGhostFrameRef.current = null;
      }
    };
  }, [editMode]);
  const finishEditMode = () => {
    cancelLongPress();
    resetPointerStart();
    setEditMode(false);
    setActiveDrag(null);
    setDropTarget(null, false);
    draftRailOrderIdsRef.current = null;
    setDraftRailOrderIds(null);
    longPressTriggeredRef.current = false;
  };
  const draggedItem = draggedId ? railNav.find((item) => item.id === draggedId) : null;
  return (
    <>
      {editMode && draggedItem ? (
        <div
          ref={dragGhostRef}
          className="mobile-bottom-nav-drag-ghost-v2"
          style={{ backgroundImage: `url("${menuAssetUrl(draggedItem.id, theme)}")` }}
          aria-hidden="true"
        >
          <strong>{draggedItem.shortLabel || draggedItem.label}</strong>
        </div>
      ) : null}
      {editMode ? (
        <button type="button" className="mobile-menu-edit-done-v2" onClick={finishEditMode}>
          <Check size={14} />
          {copy("done", "Готово")}
        </button>
      ) : null}
      <nav
        ref={bottomNavRef}
        className={cn("mobile-bottom-nav-v2", editMode && "is-editing")}
        aria-label={copy("mobile_navigation", "Mobile navigation")}
        onPointerMove={handleBottomPointerMove}
        onPointerUp={handleBottomPointerUp}
        onDragOver={(event) => {
          if (!editMode) return;
          event.preventDefault();
          autoScrollBottomNav(event.clientX);
        }}
        onContextMenu={(event) => event.preventDefault()}
        onPointerCancel={() => { cancelLongPress(); resetPointerStart(); setActiveDrag(null); setDropTarget(null, false); longPressTriggeredRef.current = false; }}
      >
        {railNav.map((item) => (
          <button
            key={item.id}
            type="button"
            className={cn(
              activeView === item.id && "is-active",
              draggedId === item.id && "is-dragging",
              dragOverId === item.id && draggedId !== item.id && "is-drop-target",
              dragOverId === item.id && draggedId !== item.id && !dragAfterTarget && "is-drop-before",
              dragOverId === item.id && draggedId !== item.id && dragAfterTarget && "is-drop-after",
            )}
            data-view={item.id}
            draggable={false}
            onDragStart={() => setActiveDrag(item.id)}
            onDragEnd={() => { setActiveDrag(null); setDropTarget(null, false); longPressTriggeredRef.current = false; }}
            onDragOver={(event) => {
              if (!editMode) return;
              event.preventDefault();
              autoScrollBottomNav(event.clientX);
              const rect = event.currentTarget.getBoundingClientRect();
              setDropTarget(item.id, event.clientX > rect.left + rect.width / 2);
            }}
            onDrop={(event) => {
              const rect = event.currentTarget.getBoundingClientRect();
              previewRailReorder(item.id, event.clientX > rect.left + rect.width / 2);
              commitRailReorder();
              setDropTarget(null, false);
            }}
            onPointerDown={(event) => {
              try {
                event.currentTarget.setPointerCapture?.(event.pointerId);
              } catch {
                // Synthetic tests and some mobile webviews can provide a non-capturable pointer.
              }
              if (editMode) event.preventDefault();
              pointerStartRef.current = { id: item.id, x: event.clientX, y: event.clientY };
              updateDragGhost(event.clientX, event.clientY);
              startLongPress(item.id);
              if (editMode) {
                setActiveDrag(item.id);
                setDropTarget(item.id, false);
              }
            }}
            onPointerCancel={() => { cancelLongPress(); resetPointerStart(); setActiveDrag(null); setDropTarget(null, false); longPressTriggeredRef.current = false; }}
            onClick={() => {
              if (longPressTriggeredRef.current) {
                longPressTriggeredRef.current = false;
                return;
              }
              chooseView(item.id);
            }}
            aria-pressed={activeView === item.id}
          >
            <span style={{ backgroundImage: `url("${menuAssetUrl(item.id, theme)}")` }} />
            <strong>{item.shortLabel || item.label}</strong>
          </button>
        ))}
      </nav>
    </>
  );
}

function ContextHeader({
  item,
  details,
  busy,
  onPrimary,
}: {
  item: NavItem;
  details: { kicker: string; title: string; subtitle: string; action?: string };
  busy: string | null;
  onPrimary: () => void;
}) {
  const Icon = item.icon;
  return (
    <section className="context-header">
      <div className="context-header__icon"><Icon size={24} /></div>
      <div>
        <span className="eyebrow">{details.kicker}</span>
        <h1>{details.title}</h1>
        <p>{details.subtitle}</p>
      </div>
      {details.action ? (
        <Button className="context-header__action" onClick={onPrimary} disabled={Boolean(busy)}>
          {busy ? <Spinner size="small" className="button-spinner-v2" /> : <Play size={18} />}
          {details.action}
        </Button>
      ) : null}
    </section>
  );
}

function StatusBanner({ kind, text, onClose }: { kind: StatusKind; text: string; onClose: () => void }) {
  useEffect(() => {
    const timer = window.setTimeout(onClose, 3000);
    return () => window.clearTimeout(timer);
  }, [onClose, text]);

  return (
    <div className={cn("v2-status", `v2-status--${kind}`)}>
      <span>{kind === "ok" ? <Check size={16} /> : kind === "error" ? <Search size={16} /> : <Sparkles size={16} />}</span>
      <p>{text}</p>
      <button type="button" onClick={onClose} aria-label="Close"><X size={14} /></button>
    </div>
  );
}

function PaymentNoticeDialog({
  notice,
  copy,
  onClose,
}: {
  notice: PaymentNotice;
  copy: (key: string, fallback: string) => string;
  onClose: () => void;
}) {
  const success = notice.kind === "success";
  return (
    <Dialog open onOpenChange={(open) => !open && onClose()}>
      <DialogContent className={cn("v2-dialog-content payment-notice-dialog-v2", success ? "is-success" : "is-warning")}>
        <DialogBody className="payment-notice-dialog-v2__body">
          {success ? <CheckCircle className="payment-notice-dialog-v2__icon" /> : <AlertCircle className="payment-notice-dialog-v2__icon" />}
          <DialogTitle className="payment-notice-dialog-v2__title">{notice.title}</DialogTitle>
          <DialogDescription className="payment-notice-dialog-v2__description">{notice.description}</DialogDescription>
          {success ? (
            <div className="payment-notice-dialog-v2__summary">
              <p><span>{copy("plan", "Plan")}</span><strong>{notice.planTitle || copy("premium", "Premium")}</strong></p>
              <p><span>{copy("period", "Period")}</span><strong>{notice.period || copy("premium", "Premium")}</strong></p>
              {notice.amount ? <p><span>{copy("amount", "Amount")}</span><strong>{notice.amount}</strong></p> : null}
              {notice.premiumUntil ? <p><span>{copy("premium_until", "Premium until")}</span><strong>{notice.premiumUntil}</strong></p> : null}
            </div>
          ) : null}
        </DialogBody>
        <DialogFooter className="payment-notice-dialog-v2__footer">
          <DialogClose asChild>
            <Button className="w-full" onClick={onClose}>{success ? copy("continue", "Continue") : copy("ok", "OK")}</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function BugReportDialog({
  activeView,
  copy,
  busy,
  onClose,
  onSubmit,
}: {
  activeView: ViewId;
  copy: (key: string, fallback: string) => string;
  busy: string | null;
  onClose: () => void;
  onSubmit: (form: FormData) => Promise<void>;
}) {
  const [message, setMessage] = useState("");
  const [screenshots, setScreenshots] = useState<File[]>([]);
  const canSend = message.trim().length >= 8 || screenshots.length > 0;
  const screenshotKey = (file: File) => `${file.name || "clipboard"}:${file.type || "application/octet-stream"}:${file.size}`;
  const addScreenshotFiles = (files: File[]) => {
    const imageFiles = files.filter((file) => file.type.startsWith("image/") || /\.(png|jpe?g|webp|gif|bmp|svg|heic|heif)$/i.test(file.name));
    if (!imageFiles.length) return;
    setScreenshots((current) => {
      const seen = new Set(current.map(screenshotKey));
      const next = [...current];
      for (const file of imageFiles) {
        const key = screenshotKey(file);
        if (seen.has(key)) continue;
        seen.add(key);
        next.push(file);
      }
      return next.slice(0, 5);
    });
  };
  const addClipboardScreenshots = (event: ClipboardEvent<HTMLTextAreaElement>) => {
    const files = Array.from(event.clipboardData.files || []);
    const itemFiles = Array.from(event.clipboardData.items || [])
      .filter((item) => item.kind === "file")
      .map((item) => item.getAsFile())
      .filter((file): file is File => Boolean(file));
    addScreenshotFiles([...files, ...itemFiles]);
  };
  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!canSend) return;
    const form = new FormData();
    form.set("message", message.trim() || copy("bug_report_image_only", "Скриншот или изображение без текста."));
    form.set("view", activeView);
    form.set("user_agent", navigator.userAgent);
    screenshots.forEach((file) => form.append("screenshot", file));
    await onSubmit(form);
  };
  const screenshotLabel = screenshots.length
    ? screenshots.length === 1
      ? screenshots[0].name
      : `${screenshots.length} ${copy("screenshots_attached", "скриншота прикреплено")}`
    : copy("attach_screenshot", "Прикрепить скриншот");
  return (
    <Dialog open onOpenChange={(open) => { if (!open) onClose(); }}>
      <DialogContent className="v2-dialog-content bug-report-dialog-v2">
        <form onSubmit={submit}>
          <DialogBody className="bug-report-dialog-v2__body">
            <Bug className="bug-report-dialog-v2__icon" size={44} />
            <DialogTitle>{copy("report_bug", "Сообщить об ошибке")}</DialogTitle>
            <DialogDescription>
              {copy("report_bug_body", "Опишите, что произошло. Скриншот можно прикрепить сразу, отчет сохранится в файл баг-репортов.")}
            </DialogDescription>
            <label className="bug-report-dialog-v2__field">
              <span>{copy("problem_description", "Описание проблемы")}</span>
              <textarea
                value={message}
                onChange={(event) => setMessage(event.target.value)}
                onPaste={addClipboardScreenshots}
                placeholder={copy("problem_description_placeholder", "Опишите Вашу проблему")}
                rows={5}
              />
            </label>
            <label className="bug-report-dialog-v2__upload">
              <FileImage size={18} />
              <span>{screenshotLabel}</span>
              <input
                type="file"
                accept="image/*,.png,.jpg,.jpeg,.webp,.gif,.bmp,.svg,.heic,.heif"
                multiple
                onChange={(event) => addScreenshotFiles(Array.from(event.target.files || []))}
              />
            </label>
          </DialogBody>
          <DialogFooter>
            <DialogClose asChild>
              <Button type="button" variant="outline" onClick={onClose}>{copy("bug_report_close", "Закрыть")}</Button>
            </DialogClose>
            <Button type="submit" disabled={!canSend || busy === "bug-report"}>
              {busy === "bug-report" ? <Spinner size="small" className="button-spinner-v2" /> : <Send size={16} />}
              {copy("send_report", "Отправить отчет")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function AppGuideDialog({
  copy,
  onClose,
}: {
  copy: (key: string, fallback: string) => string;
  onClose: () => void;
}) {
  const [activeStep, setActiveStep] = useState(0);
  const steps = [
    {
      title: copy("app_guide_step_1_title", "AI Tutor and bot"),
      body: copy("app_guide_step_1_body", "Use the website and Telegram as one learning profile. The AI Tutor runs a guided lesson: story, words, answer checks, writing, listening, pronunciation, dialogue, final word check, and review."),
      details: [
        copy("app_guide_step_1_detail_1", "Start from Today when you want a ready route, or open AI Tutor for a complete lesson."),
        copy("app_guide_step_1_detail_2", "Send text, voice, or a practice photo where the mode supports it; the bot and web app keep the same account progress."),
        copy("app_guide_step_1_detail_3", "Completed lessons stay in history, so you can return to the topic, level, summary, and review."),
      ],
      result: copy("app_guide_step_1_result", "Best first session: Today -> AI Tutor -> one weak-spot repair."),
      icon: Bot,
    },
    {
      title: copy("app_guide_step_2_title", "How plans differ"),
      body: copy("app_guide_step_2_body", "Free keeps the starter learning loop. Premium opens AI Tutor with listening practice, pronunciation scoring, voice review, image tools, and higher daily limits. Platinum is for dense study with the maximum limits."),
      details: [
        copy("app_guide_step_2_detail_1", "Free is enough to test the route and keep a small daily habit."),
        copy("app_guide_step_2_detail_2", "Premium is the normal daily mode when you need guided lessons, voice, listening, and pronunciation feedback."),
        copy("app_guide_step_2_detail_3", "Platinum is useful before travel, work, exams, or any period where you practice heavily every day."),
      ],
      result: copy("app_guide_step_2_result", "Open Premium to compare current prices, payment methods, activation keys, and active status."),
      icon: CircleDollarSign,
    },
    {
      title: copy("app_guide_step_3_title", "Global theme and settings"),
      body: copy("app_guide_step_3_body", "The theme toggle changes the whole web interface: background, menu icons, panels, dialogs, login visuals, and contrast. Brightness adjusts the app surface without changing your learning data."),
      details: [
        copy("app_guide_step_3_detail_1", "Use the sun/moon control in the top bar or mobile quick controls to switch light and dark themes."),
        copy("app_guide_step_3_detail_2", "Settings also control interface language, learning language, level, and global learning focus."),
        copy("app_guide_step_3_detail_3", "Global learning focus guides new lessons and practice topics, so the app stops jumping to random themes."),
      ],
      result: copy("app_guide_step_3_result", "Theme is visual; language, level, and learning focus affect future content."),
      icon: Settings,
    },
    {
      title: copy("app_guide_step_4_title", "How words enter vocabulary"),
      body: copy("app_guide_step_4_body", "A word does not become learned after one click. New words first go to review and spelling practice; they enter the learned vocabulary after enough correct answers, including the 10-correct-answer mastery rule used by the bot."),
      details: [
        copy("app_guide_step_4_detail_1", "Use Learn words for new cards, Review game for choices, and Spelling when you need to type from memory."),
        copy("app_guide_step_4_detail_2", "If a word or translation is wrong, report it from the AI Tutor word report button so the correction reaches the review flow."),
        copy("app_guide_step_4_detail_3", "Mistakes are grouped by grammar, word order, vocabulary, politeness, and spelling, then trained separately."),
      ],
      result: copy("app_guide_step_4_result", "Do not chase a huge list. Repeat fewer words correctly until they become stable."),
      icon: BookOpen,
    },
    {
      title: copy("app_guide_step_5_title", "Notes, audio, and reports"),
      body: copy("app_guide_step_5_body", "After lessons and practice, save useful answers, corrections, and translations to Notes. Audio buttons let you listen to messages and examples again; pronunciation and listening modes compare what you heard or said."),
      details: [
        copy("app_guide_step_5_detail_1", "Save lesson answers to Notes when they are phrases you will actually reuse."),
        copy("app_guide_step_5_detail_2", "Use Listening to hear a phrase first, repeat it, then check the weak words instead of reading the target immediately."),
        copy("app_guide_step_5_detail_3", "Use bug reports for interface issues and word reports for wrong vocabulary, translations, or lesson words."),
      ],
      result: copy("app_guide_step_5_result", "A strong loop is: listen, answer, save the best phrase, report bad data, repeat the weak spot."),
      icon: Volume2,
    },
  ];
  const totalSteps = steps.length;
  const safeStep = Math.min(activeStep, totalSteps - 1);
  const step = steps[safeStep];
  const StepIcon = step.icon;
  const atStart = safeStep === 0;
  const atEnd = safeStep === totalSteps - 1;
  const goBack = () => setActiveStep((current) => Math.max(0, current - 1));
  const goNext = () => {
    if (atEnd) {
      onClose();
      return;
    }
    setActiveStep((current) => Math.min(totalSteps - 1, current + 1));
  };

  return (
    <Dialog open onOpenChange={(open) => { if (!open) onClose(); }}>
      <DialogContent className="v2-dialog-content app-guide-dialog-v2">
        <DialogBody className="app-guide-dialog-v2__body">
          <div className="app-guide-dialog-v2__head">
            <CircleHelp className="app-guide-dialog-v2__icon" size={42} />
            <span className="app-guide-page-v2" aria-live="polite">
              {copy("app_guide_page_label", "Page")} {safeStep + 1}/{totalSteps}
            </span>
          </div>
          <DialogTitle className="app-visually-hidden-v2">{copy("app_guide_title", "Quick start guide")}</DialogTitle>
          <DialogDescription className="app-visually-hidden-v2">{copy("app_guide_body", "Five pages for the main app flows: lessons, plans, theme, vocabulary, notes, audio, and reports.")}</DialogDescription>
          <div className="app-guide-steps-v2">
            <article className="app-guide-step-v2" key={step.title}>
              <strong>{safeStep + 1}</strong>
              <span><b><StepIcon size={18} />{step.title}</b><small>{step.body}</small></span>
              <ul className="app-guide-points-v2">
                {step.details.map((detail) => (
                  <li key={detail}><CheckCircle size={15} /><span>{detail}</span></li>
                ))}
              </ul>
              <p className="app-guide-result-v2"><Target size={15} />{step.result}</p>
            </article>
            <div className="app-guide-dots-v2" aria-label={copy("app_guide_title", "Quick start guide")}>
              {steps.map((item, index) => (
                <button
                  type="button"
                  key={item.title}
                  className={cn("app-guide-dot-v2", index === safeStep && "is-active")}
                  aria-current={index === safeStep ? "step" : undefined}
                  aria-label={`${copy("app_guide_page_label", "Page")} ${index + 1}`}
                  onClick={() => setActiveStep(index)}
                />
              ))}
            </div>
          </div>
        </DialogBody>
        <DialogFooter className="app-guide-footer-v2">
          <Button className="app-guide-close-v2" type="button" variant="ghost" onClick={onClose}>
            <X size={16} />
            {copy("bug_report_close", "Close")}
          </Button>
          <div className="app-guide-footer-v2__pages">
            <Button className="app-guide-back-v2" type="button" variant="outline" onClick={goBack} disabled={atStart}>
              <ChevronLeft size={16} />
              {copy("back", "Back")}
            </Button>
            <Button className="app-guide-next-v2" type="button" onClick={goNext}>
              {atEnd ? <Check size={16} /> : <ChevronRight size={16} />}
              {atEnd ? copy("done", "Done.") : copy("next", "Next")}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function OnboardingDialog({
  open,
  session,
  user,
  copy,
  onComplete,
  onSkip,
}: {
  open: boolean;
  session: SessionData;
  user: UserProfile;
  copy: (key: string, fallback: string) => string;
  onComplete: (payload: { goal: string; level: string; learning_language: string; format: string }) => Promise<void>;
  onSkip: () => void;
}) {
  const languages = session.learning_languages?.length
    ? session.learning_languages
    : [
        { code: "en", native_name: "English" },
        { code: "es", native_name: "EspaГ±ol" },
        { code: "de", native_name: "Deutsch" },
        { code: "fr", native_name: "FranГ§ais" },
      ];
  const goals = [
    { id: "speaking", label: copy("onboarding_goal_speaking", "Свободно говорить"), icon: MessageCircle },
    { id: "travel", label: copy("onboarding_goal_travel", "Путешествия"), icon: Sparkles },
    { id: "work", label: copy("onboarding_goal_work", "Работа"), icon: Users },
    { id: "exam", label: copy("onboarding_goal_exam", "Экзамен"), icon: Trophy },
  ];
  const formats = [
    { id: "voice", label: copy("onboarding_format_voice", "Голосовая практика") },
    { id: "daily", label: copy("onboarding_format_daily", "10 минут каждый день") },
    { id: "words", label: copy("onboarding_format_words", "Слова + повторение") },
    { id: "roleplay", label: copy("onboarding_format_roleplay", "AI roleplay") },
  ];
  const [goal, setGoal] = useState("speaking");
  const [level, setLevel] = useState(user.level || "A1");
  const [learningLanguage, setLearningLanguage] = useState(user.learning_language || languages[0]?.code || "en");
  const [format, setFormat] = useState("daily");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!open) return;
    setLevel(user.level || "A1");
    setLearningLanguage(user.learning_language || languages[0]?.code || "en");
  }, [open, user.level, user.learning_language]);

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSaving(true);
    try {
      await onComplete({ goal, level, learning_language: learningLanguage, format });
    } finally {
      setSaving(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(nextOpen) => !nextOpen && onSkip()}>
      <DialogContent className="v2-dialog-content onboarding-dialog-v2">
        <form onSubmit={submit}>
          <DialogBody className="onboarding-dialog-v2__body">
            <span className="eyebrow">{copy("onboarding_kicker", "Старт обучения")}</span>
            <DialogTitle className="onboarding-dialog-v2__title">{copy("onboarding_title", "Соберём ваш план на неделю")}</DialogTitle>
            <DialogDescription className="onboarding-dialog-v2__description">
              {copy("onboarding_description", "Выберите цель, уровень, язык и формат занятий. Web сразу покажет главный шаг на сегодня.")}
            </DialogDescription>

            <div className="onboarding-dialog-v2__section">
              <strong>{copy("learning_goal", "Цель")}</strong>
              <div className="onboarding-choice-grid-v2">
                {goals.map((item) => {
                  const Icon = item.icon;
                  return (
                    <button key={item.id} type="button" className={goal === item.id ? "is-selected" : ""} onClick={() => setGoal(item.id)}>
                      <Icon size={17} />
                      <span>{item.label}</span>
                    </button>
                  );
                })}
              </div>
            </div>

            <div className="onboarding-dialog-v2__fields">
              <label>
                <span>{copy("level", "Уровень")}</span>
                <select value={level} onChange={(event) => setLevel(event.target.value)}>
                  {["A1", "A2", "B1", "B2", "C1", "C2"].map((item) => <option key={item} value={item}>{item}</option>)}
                </select>
              </label>
              <label>
                <span>{copy("learning_language", "Язык обучения")}</span>
                <select value={learningLanguage} onChange={(event) => setLearningLanguage(event.target.value)}>
                  {languages.map((item) => <option key={item.code} value={item.code}>{cleanAppText(item.native_name || item.code)}</option>)}
                </select>
              </label>
            </div>

            <div className="onboarding-dialog-v2__section">
              <strong>{copy("training_format", "Формат")}</strong>
              <div className="onboarding-format-grid-v2">
                {formats.map((item) => (
                  <button key={item.id} type="button" className={format === item.id ? "is-selected" : ""} onClick={() => setFormat(item.id)}>
                    {item.label}
                  </button>
                ))}
              </div>
            </div>
          </DialogBody>
          <DialogFooter className="onboarding-dialog-v2__footer">
            <Button type="button" variant="outline" onClick={onSkip}>{copy("later", "Позже")}</Button>
            <Button type="submit" disabled={saving}>{saving ? <Spinner size="small" className="button-spinner-v2" /> : <Check size={16} />}{copy("create_plan", "Создать план")}</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function PaymentModal({
  payment,
  copy,
  busy,
  onClose,
  onPay,
}: {
  payment: PaymentState;
  copy: (key: string, fallback: string) => string;
  busy: string | null;
  onClose: () => void;
  onPay: (method: "card" | "stars", yooKassaMethod?: string) => void;
}) {
  const [selectedMethod, setSelectedMethod] = useState("");
  const plan = payment.plan;
  const info = payment.info || {};
  const invoiceURL = asText(info.invoice_url || info.url || info.confirmation_url || info.bot_url, "");
  const amountLine = paymentAmountLine(info);
  const address = asText(info.address || info.wallet || info.recipient, "");
  const memo = asText(info.memo || info.comment || info.payload || info.payment_comment, "");
  const status = asText(info.status || info.payment_status || "", "");
  const paymentID = asText(info.payment_id || info.id, "");
  const network = asText(info.network || info.chain, "");
  const expiresRaw = asText(info.expires_at, "");
  const expiresLabel = (() => {
    if (!expiresRaw) return "";
    const date = new Date(expiresRaw);
    if (Number.isNaN(date.getTime())) return prettyDate(expiresRaw);
    return new Intl.DateTimeFormat(undefined, { dateStyle: "medium", timeStyle: "short" }).format(date);
  })();
  const paymentInstruction = "";
  const isConfirmed = /paid|confirmed|success|succeed|complete/i.test(status);
  const chooseMethod = (key: string, method: "card" | "stars", yooKassaMethod?: string) => {
    setSelectedMethod(key);
    onPay(method, yooKassaMethod);
  };
  const requisites = [
    [copy("status", "Status"), status],
    [copy("amount", "Amount"), amountLine],
    [copy("network", "Network"), network],
    [copy("wallet", "Wallet"), address],
    [copy("comment", "Comment"), memo],
    [copy("expires_at", "Expires at"), expiresLabel],
    [copy("tx_hash", "Transaction"), asText(info.tx_hash || info.hash, "")],
  ].filter(([, value]) => value);
  const methodNote = (() => {
    if (selectedMethod === "stars") {
      return copy("payment_stars_instruction", "Telegram Stars payment opens in Telegram. Finish it in the bot window; Premium will turn on automatically after Telegram confirms the invoice.");
    }
    return "";
  })();
  const yooKassaMethods = [
    { id: "any", label: copy("pay_yookassa_choose", "Choose in YooKassa") },
    { id: "bank_card", label: copy("pay_card", "Bank card") },
    { id: "sbp", label: copy("pay_sbp", "SBP") },
    { id: "yoo_money", label: "ЮMoney" },
  ];
  return (
    <div className="modal-backdrop-v2" role="dialog" aria-modal="true">
      <section className="v2-payment-modal">
        <div className="panel-head">
          <div>
            <span className="eyebrow">{copy("payment_options", "Payment options")}</span>
            <h2>{premiumPlanTitle(plan, copy)}</h2>
          </div>
          <Button variant="outline" size="sm" onClick={onClose}>x</Button>
        </div>
        <p>{premiumPlanBody(plan, copy)}</p>
        <div className="payment-methods-v2">
          {yooKassaMethods.map((method) => (
            <MorphButton
              key={method.id}
              text={method.label}
              icon={<CircleDollarSign size={18} />}
              isLoading={selectedMethod === `yookassa-${method.id}` && !isConfirmed}
              onClick={() => chooseMethod(`yookassa-${method.id}`, "card", method.id)}
              className={cn("payment-method-button-v2", selectedMethod === `yookassa-${method.id}` && "is-selected")}
            />
          ))}
          <MorphButton
            text={copy("pay_stars", "Telegram Stars")}
            icon={<Star size={18} />}
            variant="secondary"
            isLoading={selectedMethod === "stars" && !isConfirmed}
            onClick={() => chooseMethod("stars", "stars")}
            className={cn("payment-method-button-v2", selectedMethod === "stars" && "is-selected")}
          />
        </div>
        {methodNote ? <div className="payment-method-note-v2">{methodNote}</div> : null}
        {requisites.length || invoiceURL ? (
          <div className="payment-invoice-v2">
            <strong>{copy("payment_requisites", "Payment requisites")}</strong>
            {paymentInstruction ? <div className="payment-instruction-v2">{paymentInstruction}</div> : null}
            {requisites.map(([label, value]) => (
              <p key={label}>
                <span>{label}</span>
                <code>{value}</code>
              </p>
            ))}
            <div className="payment-actions-v2">
              {invoiceURL ? <a href={invoiceURL} target="_blank" rel="noreferrer">{copy("open_payment", "Open payment")}</a> : null}
            </div>
          </div>
        ) : null}
      </section>
    </div>
  );
}

type ViewRendererProps = {
  activeView: ViewId;
  user: UserProfile;
  copy: (key: string, fallback: string) => string;
  busy: string | null;
  messages: ChatMessage[];
  activeLessonTaskId: string;
  draft: string;
  setDraft: (value: string) => void;
  voiceFile: File | null;
  imageFile: File | null;
  setVoiceFile: (file: File | null) => void;
  setImageFile: (file: File | null) => void;
  setView: (view: ViewId) => void;
  aiTutorStep: AiTutorStep | null;
  aiTutorFeedback: AiTutorFeedbackState | null;
  submitAiTutorStep: (text: string, choice?: string) => Promise<void>;
  reportAiTutorWord: (input: { stage: string; proposedWord: string; proposedTranslation: string; comment: string }) => Promise<boolean>;
  tutorLoadError: string;
  startTutor: () => Promise<void>;
  restartTutorLesson: (lessonId: string) => Promise<boolean>;
  startLesson: () => Promise<void>;
  submitLesson: () => Promise<void>;
  submitPractice: () => Promise<void>;
  startRoleplayScenario: (scenario: RoleplayScenario) => Promise<void>;
  submitRoleplayAnswer: (scenario: RoleplayScenario, text: string, voice?: File | null) => Promise<void>;
  roleplayResult: ChatMessage | null;
  startShadowing: () => Promise<void>;
  submitShadowing: () => Promise<void>;
  pronunciationTarget: string;
  setPronunciationTarget: (value: string) => void;
  startPronunciation: () => Promise<void>;
  submitPronunciation: (target: string) => Promise<ApiRecord | null>;
  startWord: () => Promise<void>;
  wordChallenge: WordChallenge | null;
  answerWord: (answerId: string) => Promise<void>;
  reportWord: (input: { wordId: string; proposedWord: string; proposedTranslation: string; comment: string }) => Promise<boolean>;
  wordResult: TrainerResult | null;
  startWordGame: () => Promise<void>;
  answerWordGame: (answerId: string) => Promise<void>;
  wordGameResult: TrainerResult | null;
  spellingChallenge: SpellingChallenge | null;
  startSpelling: () => Promise<void>;
  answerSpelling: (giveUp?: boolean) => Promise<void>;
  spellingResult: TrainerResult | null;
  levelQuestion: LevelQuestion | null;
  startLevel: () => Promise<void>;
  answerLevel: (answer: number) => Promise<void>;
  levelResult: TrainerResult | null;
  vocabulary: VocabularyItem[];
  vocabularyMeta: VocabularyPage;
  loadVocabulary: (page?: number, options?: { navigate?: boolean }) => Promise<VocabularyItem[]>;
  mistakes: MistakeItem[];
  loadMistakes: (options?: { navigate?: boolean }) => Promise<MistakeItem[]>;
  mistakePractice: MistakeItem | null;
  mistakePracticeIndex: number | null;
  mistakeAnswer: string;
  setMistakeAnswer: (value: string) => void;
  startMistakePractice: (index?: number) => Promise<void>;
  submitMistakeAnswer: () => Promise<void>;
  cancelMistakePractice: () => void;
  mistakeResult: TrainerResult | null;
  deleteMistake: (index: number) => Promise<void>;
  clearMistakes: () => Promise<void>;
  shadowingTarget: string;
  leaderboard: LeaderboardEntry[];
  leaderboardMeta: LeaderboardMeta;
  leaderboardLanguage: string;
  setLeaderboardLanguage: (language: string) => void;
  loadLeaderboard: (language?: string, options?: { navigate?: boolean }) => Promise<void>;
  premiumPlans: PremiumPlan[];
  loadPremiumPlans: () => Promise<void>;
  paymentHistory: PaymentHistoryItem[];
  setPayment: (payment: PaymentState | null) => void;
  saveSettings: (payload: ApiRecord) => Promise<void>;
  changePassword: (payload: { current_password: string; new_password: string; new_password_confirm: string }) => Promise<void>;
  startTelegramCode: (body?: ApiRecord) => Promise<TelegramCodeRequest | null>;
  verifyTelegramCode: (token: string, code: string) => Promise<boolean>;
  session: SessionData;
  activationKey: string;
  setActivationKey: (value: string) => void;
  activateKey: (event: FormEvent<HTMLFormElement>) => Promise<void>;
  toolMode: ToolMode;
  setToolMode: (mode: ToolMode) => void;
  toolSourceLanguage: string;
  setToolSourceLanguage: (value: string) => void;
  toolTargetLanguage: string;
  setToolTargetLanguage: (value: string) => void;
  toolVoiceFile: File | null;
  setToolVoiceFile: (file: File | null) => void;
  toolImageFile: File | null;
  setToolImageFile: (file: File | null) => void;
  toolResult: TranslatorResult | null;
  submitTool: () => Promise<void>;
  selectedAwardLevel: number | null;
  setSelectedAwardLevel: (level: number | null) => void;
  phrasebook: PhrasebookItem[];
  isPhraseSaved: (phrase: string) => boolean;
  savePhrase: (phrase: string, source: PhrasebookSource, details?: ApiRecord) => void;
  removePhrase: (id: string) => void;
  habitLog: Record<string, HabitDay>;
  dailyBonus: DailyBonusNotice | null;
  claimDailyBonus: () => Promise<void>;
  openBugReport: () => void;
  openGuide: () => void;
};

function ViewRenderer(props: ViewRendererProps) {
  if (props.activeView === "home") return <HomeView {...props} />;
  if (props.activeView === "tutor") return <TutorView {...props} />;
  if (props.activeView === "lesson") return <ChatWorkView mode="lesson" {...props} />;
  if (props.activeView === "practice") return <ChatWorkView mode="practice" {...props} />;
  if (props.activeView === "roleplay") return <RoleplayView {...props} />;
  if (props.activeView === "shadowing") return <ChatWorkView mode="shadowing" {...props} />;
  if (props.activeView === "pronunciation") return <PronunciationDashboardView {...props} />;
  if (props.activeView === "tools") return <ToolsView {...props} />;
  if (props.activeView === "words") return <ChoiceTrainer mode="words" {...props} />;
  if (props.activeView === "word-game") return <ChoiceTrainer mode="word-game" {...props} />;
  if (props.activeView === "spelling") return <SpellingView {...props} />;
  if (props.activeView === "level") return <LevelView {...props} />;
  if (props.activeView === "vocabulary") return <VocabularyView {...props} />;
  if (props.activeView === "phrasebook") return <PhrasebookView {...props} />;
  if (props.activeView === "offline") return <OfflineDecksView {...props} />;
  if (props.activeView === "mistakes") return <MistakesView {...props} />;
  if (props.activeView === "leaderboard") return <LeaderboardView {...props} />;
  if (props.activeView === "premium") return <PremiumView {...props} />;
  if (props.activeView === "settings") return <SettingsView {...props} />;
  if (props.activeView === "dashboard" || props.activeView === "progress") return <TeacherDashboardView {...props} />;
  return <MetricsView {...props} />;
}

function TutorView({ user, session, aiTutorStep, aiTutorFeedback, submitAiTutorStep, reportAiTutorWord, tutorLoadError, startTutor, restartTutorLesson, busy, copy, savePhrase, isPhraseSaved, setView }: ViewRendererProps) {
  const display = (value: unknown) => cleanAppText(value).trim();
  const tutorPremiumLocked = !user.premium;
  const [tutorDraft, setTutorDraft] = useState("");
  const [selectedChoice, setSelectedChoice] = useState("");
  const [wordReportOpen, setWordReportOpen] = useState(false);
  const [wordReportWord, setWordReportWord] = useState("");
  const [wordReportTranslation, setWordReportTranslation] = useState("");
  const [wordReportComment, setWordReportComment] = useState("");
  const [wordReportError, setWordReportError] = useState("");
  const [wordReportPending, setWordReportPending] = useState(false);
  const [completedLessonsOpen, setCompletedLessonsOpen] = useState(false);
  const [completedLessonPage, setCompletedLessonPage] = useState(0);
  const [completedLessonsVersion, setCompletedLessonsVersion] = useState(0);
  const [remoteCompletedLessons, setRemoteCompletedLessons] = useState<TutorCompletedLessonRecord[]>([]);
  const [completedLessonsLoading, setCompletedLessonsLoading] = useState(false);
  const [completedLessonsError, setCompletedLessonsError] = useState("");
  const [restartingLessonId, setRestartingLessonId] = useState("");
  const [previewStage, setPreviewStage] = useState("");
  const completedRecordedRef = useRef("");
  const accountKey = asText(session.account?.login || user.telegram_account?.id || user.created_at || "guest", "guest");
  const completedLessonsKey = `poliglot-tutor-completed-v3:${accountKey}`;

  useEffect(() => {
    if (tutorPremiumLocked) return;
    if (!aiTutorStep && !busy && !tutorLoadError) void startTutor();
  }, [tutorPremiumLocked, Boolean(aiTutorStep), Boolean(busy), tutorLoadError]);

  useEffect(() => {
    setTutorDraft("");
    setSelectedChoice("");
    setPreviewStage("");
    setWordReportOpen(false);
    setWordReportError("");
  }, [aiTutorStep?.stage]);

  const serverStages = useMemo(() => {
    const stages = [
      "story_intro",
      "retell",
      "question_1",
      "question_2",
      "question_3",
      ...Array.from({ length: 6 }, (_, index) => `word_learn_${index + 1}`),
      ...Array.from({ length: 6 }, (_, index) => `word_recall_${index + 1}`),
      "production",
      "lesson_feedback",
      "review_schedule",
      "complete",
    ];
    return stages;
  }, []);

  const currentStageIndex = Math.max(0, serverStages.indexOf(aiTutorStep?.stage || "story_intro"));
  const completedCount = aiTutorStep?.stage === "complete" ? serverStages.length : currentStageIndex;
  const progressPercent = Math.round((completedCount / Math.max(serverStages.length, 1)) * 100);
  const localCompletedLessons = useMemo(
    () => readTutorCompletedLessons(completedLessonsKey),
    [completedLessonsKey, completedLessonsVersion, completedLessonsOpen],
  );
  const completedLessons = useMemo(
    () => mergeTutorCompletedLessons(remoteCompletedLessons, localCompletedLessons),
    [remoteCompletedLessons, localCompletedLessons],
  );
  const completedLessonPageSize = 10;
  const completedLessonPageCount = Math.max(1, Math.ceil(completedLessons.length / completedLessonPageSize));
  const safeCompletedLessonPage = Math.min(completedLessonPage, completedLessonPageCount - 1);
  const visibleCompletedLessons = completedLessons.slice(
    safeCompletedLessonPage * completedLessonPageSize,
    safeCompletedLessonPage * completedLessonPageSize + completedLessonPageSize,
  );

  useEffect(() => {
    if (!completedLessonsOpen) return;
    let cancelled = false;
    setCompletedLessonsLoading(true);
    setCompletedLessonsError("");
    api<AiTutorCompletedLessonsResponse>("/api/ai-tutor/completed")
      .then((payload) => {
        if (cancelled) return;
        const items = (payload.items || [])
          .map(tutorCompletedLessonFromAPI)
          .filter((item): item is TutorCompletedLessonRecord => Boolean(item));
        setRemoteCompletedLessons(items);
      })
      .catch((error) => {
        if (cancelled) return;
        setCompletedLessonsError(localizedAPIError(error, copy));
      })
      .finally(() => {
        if (!cancelled) setCompletedLessonsLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [completedLessonsOpen]);

  const stageLabel = (stage: string) => {
    if (stage === "story_intro") return copy("tutor_step_story", "Story");
    if (stage === "retell") return copy("tutor_step_retell", "Retell");
    if (/^question_\d+$/.test(stage)) return `${copy("tutor_choice_progress", "Question")} ${stage.replace("question_", "")}`;
    if (isTutorWordLearnStage(stage)) return `${copy("tutor_step_words", "New words")} ${stage.replace("word_learn_", "")}`;
    if (isTutorWordRecallStage(stage)) return `${copy("tutor_final_check_progress", "Word check")} ${stage.replace("word_recall_", "")}`;
    if (stage === "production") return copy("tutor_step_writing", "Writing");
    if (stage === "lesson_feedback") return copy("tutor_step_assessment", "Tutor assessment");
    if (stage === "review_schedule") return copy("tutor_step_review", "Memory review");
    if (stage === "complete") return copy("tutor_lesson_complete", "Lesson complete");
    return copy("ai_tutor", "AI Tutor");
  };

  const stageInstruction = (stage: string) => {
    if (stage === "story_intro") return copy("tutor_story_instruction", "Read the story, remember the main plot, then continue.");
    if (stage === "retell") return copy("tutor_retell_instruction", "Retell the story in 2-3 target-language sentences.");
    if (/^question_\d+$/.test(stage)) return copy("tutor_question_instruction", "Answer the question before moving on.");
    if (isTutorWordLearnStage(stage)) return copy("tutor_words_instruction", "Study this word and the example, then continue.");
    if (isTutorWordRecallStage(stage)) return copy("tutor_final_check_instruction", "Choose the target-language word.");
    if (stage === "production") return copy("tutor_writing_instruction", "Write 2-3 target-language sentences about the story and use at least 3 lesson words.");
    if (stage === "lesson_feedback") return copy("tutor_final_assessment_instruction", "Choose how this lesson felt.");
    if (stage === "review_schedule") return copy("tutor_review_instruction", "Choose when to review this lesson.");
    if (stage === "complete") return "";
    return "";
  };

  const currentLesson = aiTutorStep?.lesson;
  const previewStep = useMemo(
    () => previewStage && currentLesson ? buildAITutorPreviewStep(previewStage, currentLesson) : null,
    [previewStage, currentLesson],
  );
  const visibleTutorStep = previewStep || aiTutorStep;
  const isPreviewing = Boolean(previewStep);
  const lesson = visibleTutorStep?.lesson || currentLesson;
  const words = lesson?.words || [];
  const options = (visibleTutorStep?.options || [])
    .map((option) => {
      const id = typeof option === "string" ? display(option) : display(option.id || option.text);
      const rawText = typeof option === "string" ? display(option) : display(option.text || option.id);
      const text = localizedTutorOptionText(id, rawText, copy);
      return id && text ? { id, text } : null;
    })
    .filter((option): option is { id: string; text: string } => Boolean(option));
  const feedbackMessage = isPreviewing ? "" : display(aiTutorFeedback?.message);
  const feedbackOK = aiTutorFeedback?.ok !== false;
  const parsedFeedback = useMemo(() => parseAITutorFeedbackJSON(aiTutorFeedback?.json), [aiTutorFeedback?.json]);
  const phraseCandidates = useMemo(() => aiTutorPhraseCandidates(parsedFeedback, savePhrase), [parsedFeedback, savePhrase]);
  const mistakeCandidates = useMemo(() => aiTutorMistakeCandidates(parsedFeedback, savePhrase), [parsedFeedback, savePhrase]);
  const tutorNoteCandidates = useMemo(() => mergeTutorNoteCandidates(phraseCandidates, mistakeCandidates), [phraseCandidates, mistakeCandidates]);
  const tutorActualErrorLines = useMemo(() => aiTutorActualErrorLines(parsedFeedback), [parsedFeedback]);
  const heroTopic = display(currentLesson?.theme || currentLesson?.title || aiTutorStep?.title);
  const heroLevel = display(currentLesson?.level || user.level || "A1");
  const heroGoal = display(currentLesson?.lesson_goal) || copy("tutor_subtitle", "One server-guided AI Tutor lesson with checked answers.");
  const needsText = visibleTutorStep?.kind === "free_text" || ((visibleTutorStep?.kind === "word_recall" || isTutorWordRecallStage(visibleTutorStep?.stage)) && options.length === 0);
  const canContinue = visibleTutorStep?.kind === "story" || visibleTutorStep?.kind === "word_learn" || isTutorWordLearnStage(visibleTutorStep?.stage);
  const isChoiceStage = options.length > 0 && !canContinue;
  const isReviewChoiceStage = Boolean(aiTutorStep && ["rating", "review"].includes(aiTutorStep.kind));
  const isComplete = aiTutorStep?.kind === "complete" || aiTutorStep?.stage === "complete";
  const isVisibleComplete = visibleTutorStep?.kind === "complete" || visibleTutorStep?.stage === "complete";
  const isProduction = visibleTutorStep?.stage === "production";
  const defaultReviewChoice = aiTutorStep?.stage === "review_schedule" ? "no_review" : aiTutorStep?.stage === "lesson_feedback" ? "good" : "";
  const productionSentenceCount = countTutorSentences(tutorDraft);
  const textBlocked = needsText && (!tutorDraft.trim() || (isProduction && productionSentenceCount < 2));
  const submitIcon = busy === "tutor" ? <Spinner size="small" className="button-spinner-v2" /> : isChoiceStage && !isReviewChoiceStage ? <Check size={16} /> : needsText ? <Send size={16} /> : <ChevronRight size={16} />;
  const submitLabel = isChoiceStage && !isReviewChoiceStage
    ? copy("tutor_check_answer", "Check answer")
    : needsText
      ? copy("send", "Send")
      : copy("tutor_continue", "Continue");

  useEffect(() => {
    if (!isComplete || !aiTutorStep) return;
    const completedID = display((aiTutorStep as AiTutorStep).lesson?.title || heroTopic || aiTutorStep.stage);
    const key = `${accountKey}:${completedID}:${heroLevel}`;
    if (completedRecordedRef.current === key) return;
    completedRecordedRef.current = key;
    writeTutorCompletedLesson(completedLessonsKey, {
      id: key,
      title: cleanTutorHistoryLabel(lesson?.title || copy("ai_tutor", "AI Tutor")),
      topic: cleanTutorHistoryLabel(heroTopic, copy("ai_tutor", "AI Tutor")),
      level: heroLevel,
      aiLesson: lesson,
    });
    setCompletedLessonsVersion((version) => version + 1);
  }, [isComplete, aiTutorStep?.stage, accountKey, completedLessonsKey, heroLevel, heroTopic, lesson?.title]);

  const submitCurrentStep = () => {
    if (!aiTutorStep || busy === "tutor") return;
    if (canContinue) {
      void submitAiTutorStep("", "continue");
      return;
    }
    if (isChoiceStage) {
      if (!selectedChoice && !isReviewChoiceStage) return;
      void submitAiTutorStep("", selectedChoice || defaultReviewChoice);
      return;
    }
    if (needsText) {
      if (!tutorDraft.trim()) return;
      if (isProduction && productionSentenceCount < 2) return;
      void submitAiTutorStep(tutorDraft.trim(), "");
    }
  };

  const openAiTutorWordReport = () => {
    const word = aiTutorStep?.word || visibleTutorStep?.word;
    setWordReportWord(display(word?.target));
    setWordReportTranslation(display(word?.interface_translation));
    setWordReportComment("");
    setWordReportError("");
    setWordReportOpen(true);
  };

  const submitAiTutorWordReport = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!aiTutorStep || !isTutorWordLearnStage(aiTutorStep.stage)) return;
    const proposedWord = wordReportWord.trim();
    const proposedTranslation = wordReportTranslation.trim();
    if (!proposedWord || !proposedTranslation) {
      setWordReportError(copy("ai_tutor_word_report_required", "Add the corrected word and translation."));
      return;
    }
    setWordReportPending(true);
    setWordReportError("");
    const ok = await reportAiTutorWord({
      stage: aiTutorStep.stage,
      proposedWord,
      proposedTranslation,
      comment: wordReportComment.trim(),
    });
    setWordReportPending(false);
    if (!ok) {
      setWordReportError(copy("request_failed", "Request failed"));
      return;
    }
    setWordReportOpen(false);
    setWordReportWord("");
    setWordReportTranslation("");
    setWordReportComment("");
  };

  const restartCompletedLesson = async (item: TutorCompletedLessonRecord) => {
    const lessonId = display(item.lessonId);
    if (!lessonId || busy === "tutor" || restartingLessonId) {
      if (!lessonId) setCompletedLessonsError(copy("tutor_completed_restart_unavailable", "This saved lesson can only be viewed."));
      return;
    }
    setCompletedLessonsError("");
    setRestartingLessonId(lessonId);
    const ok = await restartTutorLesson(lessonId);
    setRestartingLessonId("");
    if (ok) setCompletedLessonsOpen(false);
  };

  const renderWord = () => {
    const word = visibleTutorStep?.word;
    if (!word) return null;
    const isRecall = visibleTutorStep?.kind === "word_recall" || isTutorWordRecallStage(visibleTutorStep?.stage);
    if (isRecall) {
      const prompt = display(word.interface_translation || visibleTutorStep?.title);
      return prompt ? (
        <div className="tutor-word-check-v2 tutor-word-check-v2--recall">
          <strong>{prompt}</strong>
        </div>
      ) : null;
    }
    const clips = [
      { label: copy("tutor_audio_word", "Word audio"), text: display(word.audio_text_target || word.target), targetLanguage: user.learning_language },
      { label: copy("tutor_audio_example", "Example audio"), text: display(word.example_audio_text_target || word.example_sentence_target), targetLanguage: user.learning_language },
    ].filter((clip) => clip.text);
    return (
      <div className="tutor-word-check-v2">
        <strong>{display(word.target)}</strong>
        {word.interface_translation ? <span>{display(word.interface_translation)}</span> : null}
        {word.example_sentence_target ? <p>{display(word.example_sentence_target)}</p> : null}
        {clips.length ? <AudioActionRow clips={clips} /> : null}
        <PhraseQuickSave
          candidates={aiTutorWordPhraseCandidates(word, savePhrase)}
          savePhrase={savePhrase}
          isPhraseSaved={isPhraseSaved}
          copy={copy}
        />
        {!isPreviewing && isTutorWordLearnStage(aiTutorStep?.stage) ? (
          <Button className="tutor-word-report-action-v2" type="button" variant="outline" size="sm" onClick={openAiTutorWordReport}>
            <Bug size={15} />
            <span>{copy("ai_tutor_word_report_button", "Report mistake")}</span>
          </Button>
        ) : null}
      </div>
    );
  };

  const renderMaterial = () => {
    if (!visibleTutorStep) return null;
    const storyText = display(lesson?.story?.text_target);
    const storyAudioText = display(lesson?.story?.audio_text_target || lesson?.story?.text_target);
    const storyTitle = cleanTutorPrompt(lesson?.story?.title_interface || lesson?.story?.story_title || visibleTutorStep.title);
    const localInstruction = stageInstruction(visibleTutorStep.stage);
    const questionText = cleanTutorPrompt(visibleTutorStep.question?.question_target);
    const serverInstruction = isVisibleComplete ? "" : cleanTutorPrompt(visibleTutorStep.instruction);
    const instruction = isVisibleComplete ? "" : cleanTutorPrompt(serverInstruction || localInstruction);
    const showQuestionText = Boolean(questionText && !sameTutorText(questionText, instruction));
    const cardTitle = visibleTutorStep.kind === "story" && storyTitle ? storyTitle : stageLabel(visibleTutorStep.stage);
    return (
      <article className="tutor-message-v2 is-active">
        <div className="tutor-message-v2__avatar"><Sparkles size={16} /></div>
        <div>
          <strong>{cardTitle}</strong>
          {instruction ? <p>{instruction}</p> : null}
          {visibleTutorStep.kind === "story" && storyText ? (
            <>
              <p className="tutor-task-copy-v2">{storyText}</p>
              {storyAudioText ? <AudioActionRow clips={[{ label: copy("tutor_story_audio", "Story audio"), text: storyAudioText, targetLanguage: user.learning_language }]} /> : null}
            </>
          ) : null}
          {renderWord()}
          {showQuestionText ? <p className="tutor-task-copy-v2">{questionText}</p> : null}
          {isVisibleComplete ? (
            <div className="tutor-review-summary-v2">
              <span>{copy("tutor_complete_body", "Lesson finished: story, retell, questions, words, writing, assessment, and review.")}</span>
              <span>{copy("ai_tutor_xp_awarded", "+40 XP awarded for this AI Tutor lesson.")}</span>
            </div>
          ) : null}
          {options.length ? (
            <div className="tutor-srs" role="group" aria-label={copy("tutor_options", "Options")}>
              {options.map((option) => (
                <button
                  key={option.id}
                  type="button"
                  className={cn(selectedChoice === option.id && "is-selected")}
                  onClick={() => setSelectedChoice(option.id)}
                >
                  {option.text}
                </button>
              ))}
            </div>
          ) : null}
          {words.length && visibleTutorStep.kind === "story" ? (
            <div className="tutor-review-list">
              {words.slice(0, 6).map((word) => <span key={word.id || word.target}>{display(word.target)} - {display(word.interface_translation)}</span>)}
            </div>
          ) : null}
        </div>
      </article>
    );
  };

  return (
    <div className="tutor-workspace">
      <section className="v2-panel tutor-hero">
        <div className="tutor-hero__brand" aria-hidden="true">
          <img src="/app/assets/brand-logo-mini.png" alt="" />
        </div>
        <div className="tutor-hero__copy">
          <span className="eyebrow">{copy("ai_tutor", "AI Tutor")}</span>
          <h2>{tutorPremiumLocked ? copy("tutor_premium_title", "AI Tutor is included with Premium") : aiTutorStep ? `${heroLevel} · ${heroTopic || copy("tutor_words_title", "Theme vocabulary")}` : copy("tutor_loading", "Preparing a guided lesson")}</h2>
          <p>{tutorPremiumLocked ? copy("tutor_premium_body", "Free keeps basic text practice. Upgrade to Premium to use guided AI Tutor lessons.") : heroGoal}</p>
          <div className="tutor-hero__meta">
            <span>{copy("tutor_duration", "5-10 minutes")}</span>
            <span>{copy("tutor_course_label", "Course")}: AI</span>
          </div>
        </div>
        <div className="tutor-hero__actions">
          <Button className="tutor-completed-lessons-button-v2" type="button" variant="outline" onClick={() => { setCompletedLessonPage(0); setCompletedLessonsVersion((version) => version + 1); setCompletedLessonsOpen(true); }}>
            <BookOpen size={17} />
            <span>{copy("tutor_completed_lessons", "Completed lessons")}</span>
          </Button>
          <Button type="button" onClick={() => tutorPremiumLocked ? setView("premium") : void startTutor()} disabled={busy === "tutor"}>
            {busy === "tutor" ? <Spinner size="small" className="button-spinner-v2" /> : tutorPremiumLocked ? <Crown size={18} /> : <Sparkles size={18} />}
            <span>{tutorPremiumLocked ? copy("tutor_premium_cta", "Upgrade to Premium") : copy("tutor_new_lesson", "New tutor lesson")}</span>
          </Button>
        </div>
      </section>

      {tutorPremiumLocked ? (
        <section className="v2-panel tutor-paywall-v2">
          <div className="tutor-paywall-v2__icon"><Crown size={24} /></div>
          <div>
            <span className="eyebrow">{copy("premium", "Premium")}</span>
            <h2>{copy("tutor_premium_title", "AI Tutor is included with Premium")}</h2>
            <p>{copy("tutor_premium_body", "Free keeps basic text practice. Upgrade to Premium to use guided AI Tutor lessons.")}</p>
          </div>
          <div className="tutor-paywall-v2__features">
            <span><ShieldCheck size={16} />{copy("free_feature_daily", "Daily habit, starter lessons, and basic word training")}</span>
            <span><Sparkles size={16} />{copy("premium_feature_ai_tutor", "AI-guided tutor lessons included")}</span>
            <span><Crown size={16} />{copy("platinum_feature_priority", "Best mode for dense daily study")}</span>
          </div>
          <Button type="button" onClick={() => setView("premium")}>
            <Crown size={18} />
            <span>{copy("tutor_premium_cta", "Upgrade to Premium")}</span>
          </Button>
        </section>
      ) : !aiTutorStep ? (
        <section className={cn("v2-panel tutor-loading-panel", tutorLoadError && "is-error")}>
          {tutorLoadError ? (
            <>
              <AlertCircle size={34} />
              <h3>{copy("tutor_start_failed", "Could not build the lesson")}</h3>
              <p>{tutorLoadError}</p>
              <Button type="button" onClick={() => void startTutor()} disabled={busy === "tutor"}>
                {busy === "tutor" ? <Spinner size="small" className="button-spinner-v2" /> : <Repeat2 size={16} />}
                {copy("retry", "Retry")}
              </Button>
            </>
          ) : (
            <>
              <Spinner size="large" show />
              <p>{copy("tutor_loading_body", "Building a checked AI Tutor lesson.")}</p>
            </>
          )}
        </section>
      ) : (
        <div className="tutor-session-v2">
          <aside className="v2-panel tutor-plan-v2" aria-label={copy("tutor_flow", "Lesson flow")}>
            <div className="panel-head">
              <span className="eyebrow">{copy("tutor_progress_label", "Progress")}</span>
              <h2>{progressPercent}%</h2>
            </div>
            <div className="tutor-progress-v2" aria-hidden="true"><span style={{ width: `${progressPercent}%` }} /></div>
            <div className="tutor-step-list-v2">
              {serverStages.map((stage, index) => (
                <button
                  key={stage}
                  type="button"
                  className={cn("tutor-step-v2", index < currentStageIndex && "is-done", stage === aiTutorStep.stage && "is-active", stage === previewStage && "is-preview")}
                  disabled={index > currentStageIndex}
                  onClick={() => setPreviewStage(index < currentStageIndex ? stage : "")}
                >
                  <span>{index < currentStageIndex ? <Check size={13} /> : index + 1}</span>
                  <strong>{stageLabel(stage)}</strong>
                </button>
              ))}
            </div>
          </aside>

          <section className="v2-panel tutor-context-v2">
            <div className="tutor-context-v2__head">
              <div>
                <span className="eyebrow">{copy("tutor_context_title", "Tutor context")}</span>
                <h2>{stageLabel(visibleTutorStep?.stage || aiTutorStep.stage)}</h2>
              </div>
              {isPreviewing ? (
                <Button type="button" variant="outline" size="sm" onClick={() => setPreviewStage("")}>
                  <ChevronRight size={15} />
                  {copy("tutor_back_to_current", "Back to current step")}
                </Button>
              ) : (
                <span>{completedCount}/{serverStages.length}</span>
              )}
            </div>

            {feedbackMessage ? <p className={cn("tutor-feedback-v2", feedbackOK ? "is-success" : "is-error")}>{feedbackMessage}</p> : null}
            <div className="tutor-transcript-v2" aria-live="polite">
              {renderMaterial()}
            </div>

            {isPreviewing ? (
              <div className="tutor-composer-v2 tutor-composer-v2--preview">
                <Button type="button" variant="outline" onClick={() => setPreviewStage("")}>
                  <ChevronRight size={16} />
                  <span>{copy("tutor_back_to_current", "Back to current step")}</span>
                </Button>
              </div>
            ) : !isComplete ? (
              <form
                className="tutor-composer-v2"
                onSubmit={(event) => {
                  event.preventDefault();
                  submitCurrentStep();
                }}
              >
                {needsText ? (
                  <div className="composer-textarea-shell-v2">
                    <textarea
                      value={tutorDraft}
                      onChange={(event) => setTutorDraft(event.target.value)}
                      placeholder={isProduction ? copy("tutor_production_placeholder", "Write 2-3 complete sentences.") : copy("tutor_answer_field_placeholder", "Enter your answer in this field.")}
                      rows={3}
                    />
                  </div>
                ) : null}
                {isProduction && tutorDraft.trim() && productionSentenceCount < 2 ? (
                  <p className="tutor-feedback-v2 is-error">{copy("tutor_need_two_sentences", "Add at least two complete sentences to continue.")}</p>
                ) : null}
                <Button
                  type="submit"
                  disabled={busy === "tutor" || (isChoiceStage && !isReviewChoiceStage && !selectedChoice) || textBlocked}
                >
                  {submitIcon}
                  <span>{submitLabel}</span>
                </Button>
              </form>
            ) : (
              <div className="tutor-composer-v2">
                <Button type="button" onClick={() => void startTutor()} disabled={busy === "tutor"}>
                  {busy === "tutor" ? <Spinner size="small" className="button-spinner-v2" /> : <Sparkles size={16} />}
                  <span>{copy("tutor_restart", "Start another tutor lesson")}</span>
                </Button>
              </div>
            )}

            {tutorNoteCandidates.length ? <TutorNoteStrip candidates={tutorNoteCandidates} savePhrase={savePhrase} isPhraseSaved={isPhraseSaved} copy={copy} /> : null}
            {tutorActualErrorLines.length ? (
              <div className="tutor-answer-variants-v2 tutor-actual-errors-v2" aria-label={copy("mistakes", "Mistakes")}>
                <strong>{copy("mistakes", "Mistakes")}</strong>
                {tutorActualErrorLines.map((line) => <span key={line}>{line}</span>)}
              </div>
            ) : null}
          </section>
        </div>
      )}

      {wordReportOpen ? (
        <Dialog open onOpenChange={(open) => { if (!open) setWordReportOpen(false); }}>
          <DialogContent className="v2-dialog-content tutor-word-report-dialog-v2">
            <form onSubmit={submitAiTutorWordReport}>
              <DialogBody className="tutor-word-report-dialog-v2__body">
                <Bug className="tutor-word-report-dialog-v2__icon" size={34} />
                <DialogTitle>{copy("ai_tutor_word_report_title", "Report a word mistake")}</DialogTitle>
                <DialogDescription>{copy("ai_tutor_word_report_body", "Send the corrected word and translation. The lesson will continue with the next word.")}</DialogDescription>
                <label className="tutor-word-report-field-v2">
                  <span>{copy("ai_tutor_word_report_word_label", "Correct word")}</span>
                  <input
                    type="text"
                    value={wordReportWord}
                    maxLength={80}
                    onChange={(event) => setWordReportWord(event.target.value)}
                    placeholder={copy("ai_tutor_word_report_word_placeholder", "Word in the learning language")}
                    autoFocus
                  />
                </label>
                <label className="tutor-word-report-field-v2">
                  <span>{copy("ai_tutor_word_report_translation_label", "Correct translation")}</span>
                  <input
                    type="text"
                    value={wordReportTranslation}
                    maxLength={160}
                    onChange={(event) => setWordReportTranslation(event.target.value)}
                    placeholder={copy("ai_tutor_word_report_translation_placeholder", "Translation")}
                  />
                </label>
                <label className="tutor-word-report-field-v2">
                  <span>{copy("ai_tutor_word_report_comment_label", "Comment")}</span>
                  <textarea
                    value={wordReportComment}
                    maxLength={500}
                    rows={3}
                    onChange={(event) => setWordReportComment(event.target.value)}
                    placeholder={copy("ai_tutor_word_report_comment_placeholder", "Optional short note")}
                  />
                </label>
                {wordReportError ? <p className="tutor-feedback-v2 is-error">{wordReportError}</p> : null}
              </DialogBody>
              <DialogFooter className="tutor-word-report-dialog-v2__footer">
                <DialogClose asChild>
                  <Button type="button" variant="outline" disabled={wordReportPending}>
                    {copy("cancel", "Cancel")}
                  </Button>
                </DialogClose>
                <Button type="submit" disabled={wordReportPending || !wordReportWord.trim() || !wordReportTranslation.trim()}>
                  {wordReportPending ? <Spinner size="small" className="button-spinner-v2" /> : <Send size={15} />}
                  <span>{copy("ai_tutor_word_report_submit", "Send report")}</span>
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      ) : null}

      {completedLessonsOpen ? (
        <Dialog open onOpenChange={(open) => { if (!open) setCompletedLessonsOpen(false); }}>
          <DialogContent className="v2-dialog-content tutor-completed-lessons-dialog-v2">
            <DialogBody className="tutor-completed-lessons-dialog-v2__body">
              <BookOpen className="tutor-completed-lessons-dialog-v2__icon" size={42} />
              <DialogTitle>{copy("tutor_completed_lessons", "Completed lessons")}</DialogTitle>
              <DialogDescription>{copy("tutor_completed_lessons_body", "Choose a finished tutor lesson to review its topic, level, and summary.")}</DialogDescription>
              {completedLessonsLoading ? <p className="tutor-completed-lessons-empty-v2">{copy("loading", "Loading...")}</p> : null}
              {completedLessonsError ? <p className="tutor-feedback-v2 is-error">{completedLessonsError}</p> : null}
              <div className="tutor-completed-lessons-list-v2">
                {visibleCompletedLessons.length ? visibleCompletedLessons.map((item) => (
                  <button className="tutor-completed-lesson-v2" key={item.id} type="button" disabled={Boolean(restartingLessonId) || !item.lessonId} onClick={() => void restartCompletedLesson(item)}>
                    <span>
                      <strong>{item.topic || item.title}</strong>
                      <small>{prettyDate(item.completedAt)}</small>
                    </span>
                    <em>{restartingLessonId === item.lessonId ? copy("loading", "Loading...") : item.level || copy("level_label", "Level")}</em>
                  </button>
                )) : (
                  <p className="tutor-completed-lessons-empty-v2">{copy("tutor_completed_lessons_empty", "Completed tutor lessons will appear here after the final review.")}</p>
                )}
              </div>
            </DialogBody>
            <DialogFooter className="tutor-completed-lessons-footer-v2">
              <Button className="tutor-completed-lessons-back-v2" type="button" variant="outline" onClick={() => setCompletedLessonsOpen(false)}>
                <ChevronLeft size={16} />
                {copy("back", "Back")}
              </Button>
              <div className="tutor-completed-lessons-pages-v2">
                <Button className="tutor-completed-lessons-prev-v2" type="button" variant="outline" size="sm" disabled={safeCompletedLessonPage <= 0} onClick={() => setCompletedLessonPage((page) => Math.max(0, page - 1))}>
                  <ChevronLeft size={15} />
                  {copy("back", "Back")}
                </Button>
                <span>{safeCompletedLessonPage + 1}/{completedLessonPageCount}</span>
                <Button className="tutor-completed-lessons-next-v2" type="button" variant="outline" size="sm" disabled={safeCompletedLessonPage + 1 >= completedLessonPageCount} onClick={() => setCompletedLessonPage((page) => Math.min(completedLessonPageCount - 1, page + 1))}>
                  {copy("next", "Next")}
                  <ChevronRight size={15} />
                </Button>
              </div>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      ) : null}
    </div>
  );
}

function countTutorSentences(text: string) {
  const parts = cleanAppText(text).split(/[.!????]+/).map((part) => part.trim()).filter((part) => part.length > 1);
  if (parts.length) return parts.length;
  return cleanAppText(text).trim() ? 1 : 0;
}

function buildAITutorPreviewStep(stage: string, lesson: NonNullable<AiTutorStep["lesson"]>): AiTutorStep {
  if (stage === "story_intro") {
    return { stage, kind: "story", title: lesson.title || "Story", instruction: "", lesson };
  }
  if (stage === "retell") {
    return { stage, kind: "free_text", title: "Retell", instruction: "", lesson };
  }
  const questionMatch = /^question_(\d+)$/.exec(stage);
  if (questionMatch) {
    const index = Number(questionMatch[1]) - 1;
    return {
      stage,
      kind: "free_text",
      title: `Question ${index + 1}`,
      instruction: "",
      lesson,
      question: lesson.comprehension_questions?.[index],
    };
  }
  const wordLearnMatch = /^word_learn_(\d+)$/.exec(stage);
  if (wordLearnMatch) {
    const index = Number(wordLearnMatch[1]) - 1;
    return { stage, kind: "word_learn", title: "New words", instruction: "", lesson, word: lesson.words?.[index] };
  }
  const wordRecallMatch = /^word_recall_(\d+)$/.exec(stage);
  if (wordRecallMatch) {
    const index = Number(wordRecallMatch[1]) - 1;
    return { stage, kind: "word_recall", title: "Word check", instruction: "", lesson, word: lesson.words?.[index] };
  }
  if (stage === "production") {
    return { stage, kind: "free_text", title: "Writing", instruction: lesson.production_task?.instruction_interface || "", lesson };
  }
  if (stage === "lesson_feedback") {
    return { stage, kind: "rating", title: "Tutor assessment", instruction: "", lesson };
  }
  if (stage === "review_schedule") {
    return { stage, kind: "review", title: "Memory review", instruction: "", lesson };
  }
  return { stage, kind: "complete", title: lesson.title || "AI Tutor", instruction: "", lesson };
}

function parseAITutorFeedbackJSON(raw?: string) {
  if (!raw) return {} as ApiRecord;
  try {
    return getRecord(JSON.parse(raw));
  } catch {
    return {} as ApiRecord;
  }
}

function aiTutorPhraseCandidates(feedback: ApiRecord, savePhrase: ViewRendererProps["savePhrase"]) {
  const candidates: Array<{ phrase: string; source: PhrasebookSource; details?: ApiRecord }> = [];
  const seen = new Set<string>();
  const add = (phrase: unknown, note = "") => {
    const text = cleanAppText(phrase).trim();
    if (!isUsefulPhraseCandidate(text)) return;
    const key = text.toLowerCase();
    if (seen.has(key)) return;
    seen.add(key);
    candidates.push({ phrase: text, source: "lesson", details: { note } });
  };
  add(feedback.corrected_answer_target, "AI Tutor correction");
  add(feedback.corrected_version_target, "AI Tutor correction");
  recordList(feedback.recommendations_interface)
    .concat(recordList(feedback.recommendations), recordList(feedback.suggestions), recordList(feedback.tips))
    .forEach((item) => add(item, "AI Tutor recommendation"));
  return candidates.slice(0, 4);
}

function aiTutorMistakeCandidates(feedback: ApiRecord, savePhrase: ViewRendererProps["savePhrase"]) {
  void savePhrase;
  const candidates: Array<{ phrase: string; source: PhrasebookSource; details?: ApiRecord }> = [];
  const seen = new Set<string>();
  const add = (phrase: unknown, note = "") => {
    const text = cleanAppText(phrase).trim();
    if (!text || text.length < 2) return;
    const key = text.toLowerCase();
    if (seen.has(key)) return;
    seen.add(key);
    candidates.push({ phrase: text, source: "mistake", details: { note: note || "AI Tutor mistake" } });
  };
  if (Array.isArray(feedback.mistakes)) {
    feedback.mistakes.forEach((item) => {
      const record = getRecord(item);
      add(record.correction || record.corrected || record.expected, cleanAppText(record.explanation || record.issue || "AI Tutor mistake"));
    });
  }
  add(feedback.corrected_answer_target, "AI Tutor correction");
  add(feedback.corrected_version_target, "AI Tutor correction");
  return candidates.slice(0, 4);
}

function mergeTutorNoteCandidates(...groups: Array<Array<{ phrase: string; source: PhrasebookSource; details?: ApiRecord }>>) {
  const seen = new Set<string>();
  const merged: Array<{ phrase: string; source: PhrasebookSource; details?: ApiRecord }> = [];
  for (const group of groups) {
    for (const item of group) {
      const text = cleanAppText(item.phrase).replace(/\s+/g, " ").trim();
      if (!isUsefulPhraseCandidate(text)) continue;
      const key = text.toLowerCase();
      if (seen.has(key)) continue;
      seen.add(key);
      merged.push({ ...item, phrase: text });
    }
  }
  return merged.slice(0, 6);
}

function aiTutorActualErrorLines(feedback: ApiRecord) {
  if (!Array.isArray(feedback.mistakes)) return [];
  const seen = new Set<string>();
  const lines: string[] = [];
  feedback.mistakes.forEach((item) => {
    const record = getRecord(item);
    const original = cleanAppText(record.word || record.original || record.issue).trim();
    const correction = cleanAppText(record.correction || record.corrected || record.expected).trim();
    const explanation = cleanAppText(record.explanation || record.context || record.reason).trim();
    const line = [
      original && correction ? `${original} -> ${correction}` : correction || original,
      explanation,
    ].filter(Boolean).join(": ");
    if (!line || /save|phrasebook|notes?|сохран/i.test(line)) return;
    const key = line.toLowerCase();
    if (seen.has(key)) return;
    seen.add(key);
    lines.push(line);
  });
  return lines.slice(0, 4);
}

function aiTutorWordPhraseCandidates(word: NonNullable<AiTutorStep["word"]>, savePhrase: ViewRendererProps["savePhrase"]) {
  void savePhrase;
  const candidates: Array<{ phrase: string; source: PhrasebookSource; details?: ApiRecord }> = [];
  const example = cleanAppText(word.example_sentence_target).trim();
  const target = cleanAppText(word.target).trim();
  if (isUsefulPhraseCandidate(example)) {
    candidates.push({ phrase: example, source: "lesson", details: { translation: cleanAppText(word.interface_translation), note: target } });
  }
  if (target) {
    candidates.push({ phrase: target, source: "lesson", details: { translation: cleanAppText(word.interface_translation), note: cleanAppText(word.example_sentence_target) } });
  }
  return candidates;
}

function HomeView(props: ViewRendererProps) {
  const { user, copy, startLesson, startWord, startShadowing, setView, loadVocabulary, loadMistakes, loadPremiumPlans, busy, habitLog, dailyBonus, claimDailyBonus } = props;
  const todayKey = localDateKey();
  const completedDates = Object.entries(habitLog).filter(([, day]) => day.complete).map(([date]) => date);
  const loginDates = Object.entries(habitLog).filter(([, day]) => day.login).map(([date]) => date);
  const currentHabit = habitLog[todayKey] || updateHabitForToday({}, user)[todayKey];
  const bonusLocked = dailyBonusLocked(currentHabit, user);
  const canClaimBonus = Boolean(currentHabit?.complete && !bonusLocked);
  const bonusButtonLabel = bonusLocked
    ? copy("bonus_claimed", "Bonus already claimed")
    : currentHabit?.complete
      ? copy("claim_daily_bonus", "Claim XP")
      : copy("complete_daily_first", "Complete daily first");
  const streak = habitStreak(habitLog);
  const quests = [
    { label: copy("new_lesson", "Новый урок"), value: dailyQuestProgress(user.lessons_today, dailyQuestTarget.lesson), max: dailyQuestTarget.lesson, icon: BookOpen, action: startLesson, detail: copy("quest_lesson_detail", "2 коротких урока") },
    { label: copy("practice", "Практика"), value: dailyQuestProgress(user.practice_today, dailyQuestTarget.practice), max: dailyQuestTarget.practice, icon: MessageCircle, action: () => setView("practice"), detail: copy("quest_practice_detail", "3 живые фразы") },
    { label: copy("roleplay", "Ролевая"), value: dailyQuestProgress(0, dailyQuestTarget.roleplay), max: dailyQuestTarget.roleplay, icon: Sparkles, action: () => setView("roleplay"), detail: copy("quest_roleplay_detail", "1 сцена на 5 минут") },
    { label: copy("pronunciation", "Произношение"), value: dailyQuestProgress(user.voice_today, dailyQuestTarget.pronunciation), max: dailyQuestTarget.pronunciation, icon: Activity, action: () => setView("pronunciation"), detail: copy("quest_pronunciation_detail", "1 проверка произношения") },
    { label: copy("vocabulary", "Словарик"), value: dailyQuestProgress(0, dailyQuestTarget.vocabulary), max: dailyQuestTarget.vocabulary, icon: ListChecks, action: () => void loadVocabulary(), detail: copy("quest_vocabulary_detail", "5 слов + 1 пример") },
    { label: copy("mistakes", "Словарик ошибок"), value: dailyQuestProgress(Number(user.mistakes || 0) > 0 ? 0 : 1, dailyQuestTarget.mistakes), max: dailyQuestTarget.mistakes, icon: Search, action: loadMistakes, detail: copy("quest_mistakes_detail", "1-2 ошибки в ремонт") },
    { label: copy("shadowing", "Аудирование"), value: dailyQuestProgress(user.voice_today, dailyQuestTarget.listening), max: dailyQuestTarget.listening, icon: Volume2, action: startShadowing, detail: copy("quest_listening_detail", "2 аудио-повтора") },
  ];
  const routeSteps = [
    { title: copy("today_plan_words_title", "Words"), body: copy("today_plan_words_body", "Warm up with 5 words from the current level and keep only the hard ones for review."), icon: Brain, action: () => void startWord() },
    { title: copy("today_plan_speak_title", "Speak"), body: copy("today_plan_speak_body", "Turn the words into 3 useful phrases for your real situations."), icon: MessageCircle, action: () => setView("practice") },
    { title: copy("today_plan_listen_title", "Listen"), body: copy("today_plan_listen_body", "Repeat one short audio phrase and watch the weak sounds."), icon: Volume2, action: () => void startShadowing() },
    { title: copy("today_plan_repair_title", "Repair"), body: copy("today_plan_repair_body", "Close one mistake or spelling gap before adding new material."), icon: Search, action: () => void loadMistakes() },
  ];
  const weekPlan = [
    { day: copy("monday_short", "Mon"), title: copy("plan_words", "Vocabulary base"), meta: copy("plan_words_count", "Build one theme: meaning, example, review."), view: "words" as ViewId },
    { day: copy("tuesday_short", "Tue"), title: copy("plan_practice", "Phrase output"), meta: copy("plan_practice_count", "Make useful phrases and save the best one."), view: "practice" as ViewId },
    { day: copy("wednesday_short", "Wed"), title: copy("plan_listening", "Listening loop"), meta: copy("plan_listening_count", "Hear, repeat, compare, then retry once."), view: "shadowing" as ViewId },
    { day: copy("thursday_short", "Thu"), title: copy("plan_roleplay", "Dialogue"), meta: copy("plan_roleplay_count", "Run one practical scene to use the week theme."), view: "roleplay" as ViewId },
    { day: copy("friday_short", "Fri"), title: copy("plan_pronunciation", "Pronunciation"), meta: copy("plan_pronunciation_count", "Repair weak words and one recurring sound."), view: "pronunciation" as ViewId },
    { day: copy("saturday_short", "Sat"), title: copy("plan_offline", "Offline review"), meta: copy("plan_offline_count", "Refresh cards and export a compact TXT pack."), view: "offline" as ViewId },
    { day: copy("sunday_short", "Sun"), title: copy("plan_dashboard", "Weekly review"), meta: copy("weekly_review", "Check progress and choose next week's focus."), view: "dashboard" as ViewId },
  ];
  const weakSpot = Number(user.mistakes || 0) > 0 ? copy("today_plan_body_active", "Today starts with one repair, then a short speak-listen loop so new material does not pile up over old mistakes.") : copy("today_plan_body_steady", "Today keeps the rhythm: words first, then output, then listening, then one quick check.");
  return (
    <div className="home-grid">
      <section className="v2-panel home-focus">
        <span className="eyebrow">{copy("today", "Today")}</span>
        <h2>{copy("training_signal", "Training signal")}</h2>
        <Meter label={copy("lessons_today", "Lessons today")} value={user.lessons_today} max={user.lesson_limit} />
        <Meter label={copy("practice_today", "Practice today")} value={user.practice_today} max={user.practice_limit} />
        <Meter label={copy("voices_today", "Voice messages today")} value={user.voice_today} max={user.voice_limit} />
      </section>
      <section className="v2-panel home-mobile-level-v2" aria-label={copy("level_label", "Level")}>
        <span className="home-mobile-level-v2__xp-label">XP</span>
        <LevelProgress user={user} copy={copy} />
      </section>
      <section className="v2-panel home-premium">
        <span className="eyebrow">{copy("premium", "Premium")}</span>
        <h2>{user.premium ? prettyDate(user.premium_until) : copy("base_plan", "Basic")}</h2>
        <p>{copy("premium_month_body", "Voice, photo, mistakes and vocabulary")}</p>
        <Button onClick={() => void loadPremiumPlans()}><Crown size={17} />{copy("payment_options", "Payment options")}</Button>
      </section>
      <section className="v2-panel home-plan-v2">
        <span className="eyebrow">{copy("daily_route", "Daily route")}</span>
        <h2>{copy("today_plan_title", "Focused study plan")}</h2>
        <p>{weakSpot}</p>
        <div className="home-plan-v2__steps">
          {routeSteps.map((step) => {
            const Icon = step.icon;
            return (
              <button key={step.title} type="button" onClick={step.action}>
                <Icon size={16} />
                <span><b>{step.title}</b><small>{step.body}</small></span>
              </button>
            );
          })}
        </div>
      </section>
      <section className="v2-panel daily-quests-v2">
        <span className="eyebrow"><Flame size={15} />{copy("daily_quests", "Daily quests")}</span>
        <h2>{copy("daily_quests_title", "Задания на сегодня")}</h2>
        <div className="daily-quests-v2__list">
          {quests.slice(0, 4).map((quest) => {
            const Icon = quest.icon;
            const done = Number(quest.value) >= Number(quest.max);
            return (
              <button key={quest.label} type="button" className={done ? "is-done" : ""} onClick={() => void quest.action()}>
                <Icon size={18} />
                <span><b>{quest.label}</b><small>{quest.detail}</small></span>
                <strong>{done ? copy("done", "Готово") : formatLimit(quest.value, quest.max)}</strong>
              </button>
            );
          })}
        </div>
      </section>
      <section className="v2-panel weekly-plan-v2">
        <span className="eyebrow"><CalendarDays size={15} />{copy("weekly_plan", "Weekly plan")}</span>
        <h2>{copy("weekly_plan_title", "План обучения на неделю")}</h2>
        <div className="weekly-plan-v2__week">
          {weekPlan.map((item) => (
            <button key={item.day} type="button" onClick={() => setView(item.view)}>
              <strong>{item.day}</strong>
              <span>{item.title}</span>
              <small>{item.meta}</small>
            </button>
          ))}
        </div>
      </section>
      <section className="v2-panel habit-calendar-v2">
        <span className="eyebrow"><CalendarDays size={15} />{copy("habit_calendar", "Календарь привычки")}</span>
        <h2>{copy("habit_calendar_title", "Streak и дейлики")}</h2>
        <div className="habit-calendar-v2__layout">
          <Calendar className="calendar-rac--habit" completedDates={completedDates} loginDates={loginDates} aria-label={copy("habit_calendar", "Календарь привычки")} />
          <div className="habit-calendar-v2__side">
            <strong>{formatStreakLabel(streak, user, copy)}</strong>
            <span>{copy("habit_green_hint", "Зеленые дни - вход и выполненный дейлик.")}</span>
            <small>{copy("today_goal", "Сегодня")}: {currentHabit?.complete ? copy("done", "Готово") : copy("in_progress", "В процессе")}</small>
            <Button type="button" onClick={() => void claimDailyBonus()} disabled={!canClaimBonus || busy === "daily-bonus"}>
              {busy === "daily-bonus" ? <Spinner size="small" className="button-spinner-v2" /> : <Flame size={16} />}
              {bonusButtonLabel}
            </Button>
          </div>
        </div>
        {dailyBonus ? (
          <div className="daily-bonus-pop-v2">
            <Sparkles size={18} />
            <strong>+{dailyBonus.xp} XP</strong>
            <span>{formatStreakLabel(dailyBonus.streak, user, copy)}</span>
          </div>
        ) : null}
      </section>
      <section className="v2-panel learning-lab-v2">
        <span className="eyebrow"><Sparkles size={15} />{copy("v2_learning_lab", "Новые функции")}</span>
        <h2>{copy("v2_learning_lab_title", "Новые обучающие режимы")}</h2>
        <div className="learning-lab-v2__grid">
          <button type="button" onClick={() => setView("tutor")}><Sparkles size={18} /><strong>{copy("ai_tutor", "AI Tutor")}</strong><span>{copy("tutor_short_hint", "One guided lesson: words, grammar, listening, writing, dialogue, and review")}</span></button>
          <button type="button" onClick={() => setView("roleplay")}><MessageCircle size={18} /><strong>{copy("ai_roleplay", "AI roleplay")}</strong><span>{copy("roleplay_short_hint", "Ресторан, работа, путешествие, экзамен")}</span></button>
          <button type="button" onClick={() => setView("pronunciation")}><Activity size={18} /><strong>{copy("pronunciation", "Pronunciation")}</strong><span>{copy("pronunciation_short_hint", "Heatmap слов и звуков")}</span></button>
          <button type="button" onClick={() => setView("offline")}><WifiOff size={18} /><strong>{copy("offline_decks", "Offline decks")}</strong><span>{copy("offline_short_hint", "Mini-decks без Telegram")}</span></button>
          <button type="button" onClick={() => setView("dashboard")}><BarChart3 size={18} /><strong>{copy("dashboard", "Dashboard")}</strong><span>{copy("dashboard_short_hint", "Прогресс, темы, активность")}</span></button>
        </div>
      </section>
    </div>
  );
}

function RoleplayView({
  startRoleplayScenario,
  submitRoleplayAnswer,
  messages,
  draft,
  setDraft,
  voiceFile,
  imageFile,
  setVoiceFile,
  setImageFile,
  busy,
  copy,
  user,
  roleplayResult,
  savePhrase,
  isPhraseSaved,
}: ViewRendererProps) {
  const [activeScenario, setActiveScenario] = useState<RoleplayScenario | null>(null);
  const isMobileUi = useMobileUiLayout();
  const beginScenario = async (scenario: RoleplayScenario) => {
    setActiveScenario(scenario);
    await startRoleplayScenario(scenario);
  };
  const resetScenario = () => {
    setActiveScenario(null);
    setDraft("");
    setVoiceFile(null);
  };
  const roleplayMessages = messages.filter((message) => message.meta === "roleplay");
  const sessionMessages = roleplayResult
    ? [roleplayResult, ...roleplayMessages.filter((message) => message.id !== roleplayResult.id)].slice(0, 8)
    : roleplayMessages.slice(0, 8);
  const phraseCandidates = phraseCandidatesFromMessages(sessionMessages);
  useEffect(() => {
    if (!activeScenario) return;
    requestAnimationFrame(() => {
      document.querySelector<HTMLElement>(".context-display--roleplay")?.scrollTo({ top: 0, behavior: "auto" });
    });
  }, [activeScenario?.id]);

  if (activeScenario) {
    return (
      <section className={cn("v2-panel roleplay-view-v2 roleplay-view-v2--session", isMobileUi ? "roleplay-view-v2--mobile-session" : "roleplay-view-v2--desktop-session")}>
        <div className="chat-workspace chat-workspace--roleplay" data-work-mode="roleplay">
          <div className="chat-workspace__output">
          <div className="roleplay-dialog-card-v2">
            <div className="panel-head roleplay-dialog-head-v2">
              <Button variant="outline" type="button" onClick={resetScenario}>
                <ChevronRight size={16} />
                {copy("choose_another_scenario", "Choose another scenario")}
              </Button>
              {busy === "roleplay" ? <Spinner size="small" /> : null}
            </div>
            <div className="roleplay-dialog-scroll-v2">
              {roleplayResult || sessionMessages.length ? (
                <ChatPanel messages={sessionMessages} copy={copy} targetLanguage={user.learning_language} />
              ) : (
                <div className="roleplay-empty-v2">
                  <Spinner size="small" show={busy === "roleplay"} />
                  <p>{copy("roleplay_loading", "Сцена готовится. Сейчас AI начнет диалог.")}</p>
                </div>
              )}
            </div>
            {phraseCandidates.length ? <PhraseQuickSave candidates={phraseCandidates} savePhrase={savePhrase} isPhraseSaved={isPhraseSaved} copy={copy} /> : null}
          </div>
          </div>
          <section className="v2-panel composer-panel-v2">
            <div className="panel-head">
              <span className="eyebrow">{copy("input", "Input")}</span>
              <h2>{copy("roleplay", "Roleplay")}</h2>
            </div>
            <div className="composer-textarea-shell-v2">
              <textarea
                value={draft}
                onChange={(event) => setDraft(event.target.value)}
                placeholder={copy("roleplay_answer_placeholder", "Ответьте в рамках роли. AI продолжит сцену и исправит фразу.")}
                onKeyDown={(event) => {
                  if (event.key === "Enter" && !event.shiftKey) {
                    event.preventDefault();
                    void submitRoleplayAnswer(activeScenario, draft, voiceFile);
                  }
                }}
              />
              <Button className="composer-submit-v2 roleplay-submit-v2" type="button" onClick={() => void submitRoleplayAnswer(activeScenario, draft, voiceFile)} disabled={busy === "roleplay"} aria-label={copy("send", "Send")}>
                {busy === "roleplay" ? <Spinner size="small" className="button-spinner-v2" /> : <Send size={17} />}
                <span>{copy("send", "Send")}</span>
              </Button>
            </div>
            <FileControls voiceFile={voiceFile} imageFile={imageFile} setVoiceFile={setVoiceFile} setImageFile={setImageFile} allowImage={false} copy={copy} />
          </section>
        </div>
      </section>
    );
  }

  return (
    <section className="v2-panel roleplay-view-v2">
      <div className="panel-head">
        <span className="eyebrow">{copy("roleplay", "Roleplay")}</span>
        <h2>{copy("roleplay_title", "AI roleplay-сценарии")}</h2>
        <p>{copy("roleplay_subtitle", "Выберите ситуацию: AI начнёт диалог, будет исправлять фразы и задавать следующий вопрос.")}</p>
      </div>
      <div className="roleplay-grid-v2">
        {roleplayScenarios.map((scenario) => {
          const Icon = scenario.icon;
          const title = roleplayScenarioTitle(scenario, user);
          const description = roleplayScenarioDescription(scenario, user);
          return (
            <button key={scenario.id} type="button" onClick={() => void beginScenario(scenario)} disabled={busy === "roleplay"}>
              <Icon size={22} />
              <strong>{title}</strong>
              <span>{description}</span>
              <small>{scenario.difficulty}</small>
            </button>
          );
        })}
      </div>
    </section>
  );
}

function PronunciationDashboardView({ messages, mistakes, user, session, shadowingTarget, pronunciationTarget, startPronunciation, submitPronunciation, voiceFile, imageFile, setVoiceFile, setImageFile, busy, copy }: ViewRendererProps) {
  const p = pronunciationLocale(user);
  const isMobileUi = useMobileUiLayout();
  const liveReports = useMemo(
    () => messages
      .map((message) => pronunciationFrom(getRecord(message.details)))
      .filter((item): item is PronunciationAssessment => Boolean(item)),
    [messages],
  );
  const storageKey = pronunciationStorageKey(session, user);
  const [storedReports, setStoredReports] = useState<StoredPronunciationAssessment[]>(() => readPronunciationHistory(storageKey));
  useEffect(() => {
    setStoredReports(readPronunciationHistory(storageKey));
  }, [storageKey]);
  useEffect(() => {
    if (!liveReports.length) return;
    setStoredReports((current) => {
      const merged = mergePronunciationHistory(liveReports, current);
      localStorage.setItem(storageKey, JSON.stringify(merged));
      return merged;
    });
  }, [storageKey, liveReports]);
  const reports = mergePronunciationHistory(liveReports, storedReports);
  const latest = reports[0];
  const score = normalizedPronunciationScore(latest);
  const defaultTarget = cleanAppText(pronunciationTarget || pronunciationPracticeFallback(user, copy));
  const problemWords = uniquePronunciationProblems(reports.flatMap((report) => report.problem_words || []));
  const fallbackWords: PronunciationProblem[] = mistakes.slice(0, 8).map((item) => ({ word: item.word || item.correction || "", issue: item.explanation || copy("mistake", "Mistake") }));
  const tokens = compactPronunciationMap(uniquePronunciationProblems(problemWords.length ? problemWords : fallbackWords), user, copy);
  const focusItems = pronunciationFocusItems(tokens, copy);
  const [practiceResult, setPracticeResult] = useState<PronunciationAssessment | null>(null);
  useEffect(() => {
    if (!pronunciationTarget && !busy) void startPronunciation();
  }, [Boolean(pronunciationTarget), Boolean(busy)]);
  const activeTarget = cleanAppText(pronunciationTarget || defaultTarget);
  useEffect(() => {
    setPracticeResult(null);
    setVoiceFile(null);
  }, [pronunciationTarget]);
  const checkPronunciation = async () => {
    const record = await submitPronunciation(activeTarget);
    const result = pronunciationFrom(record || {});
    if (result) {
      const resultWithTarget = { ...result, expected: result.expected || activeTarget };
      setPracticeResult(resultWithTarget);
      setStoredReports((current) => {
        const merged = mergePronunciationHistory([resultWithTarget], current);
        localStorage.setItem(storageKey, JSON.stringify(merged));
        return merged;
      });
    }
  };
  const nextPronunciationSample = () => {
    setPracticeResult(null);
    setVoiceFile(null);
    void startPronunciation();
  };
  return (
    <div className={cn("pronunciation-dashboard-v2", isMobileUi && "pronunciation-dashboard-v2--mobile")}>
      <section className={cn("v2-panel pronunciation-workbench-v2 pronunciation-practice-v2 pronunciation-work-window-v2", isMobileUi && "pronunciation-workbench-v2--mobile")}>
        <div className="pronunciation-target-primary-v2">
          <span className="eyebrow"><Volume2 size={15} />{copy("pronunciation_target", "Text to pronounce")}</span>
          <div>
            <h2>{activeTarget}</h2>
            <p>{copy("pronunciation_practice_body", "Послушайте пример, произнесите эту же фразу и отправьте запись на проверку.")}</p>
          </div>
          <AudioActionRow clips={[{ label: copy("spoken_model", "Spoken model"), text: activeTarget, targetLanguage: user.learning_language }]} />
        </div>
        <aside className="pronunciation-tools-v2">
          <div className="pronunciation-tool-card-v2">
            <span className="eyebrow"><Mic size={15} />{copy("record_and_check", "Record and check")}</span>
            <FileControls voiceFile={voiceFile} imageFile={imageFile} setVoiceFile={setVoiceFile} setImageFile={setImageFile} allowImage={false} copy={copy} />
            <div className="home-insights-v2__actions pronunciation-tool-actions-v2">
              <Button onClick={() => void checkPronunciation()} disabled={busy === "pronunciation"}>
                {busy === "pronunciation" ? <Spinner size="small" className="button-spinner-v2" /> : <CheckCircle size={16} />}
                {copy("check", "Check")}
              </Button>
              <Button type="button" variant="outline" size="sm" onClick={nextPronunciationSample} disabled={busy === "pronunciation-start"}>
                {busy === "pronunciation-start" ? <Spinner size="small" className="button-spinner-v2" /> : <ChevronRight size={16} />}
                {copy("next_phrase", "Next phrase")}
              </Button>
            </div>
          </div>
          {practiceResult ? (
            <div className="pronunciation-tool-card-v2 pronunciation-report-window-v2">
              <span className="eyebrow"><CheckCircle size={15} />{copy("pronunciation_result", "Статистика произношения")}</span>
              <p>{copy("pronunciation_result_body", "Разбор записи готов: оценка, слабые слова, звуки и следующий шаг.")}</p>
              <PronunciationReport pronunciation={practiceResult} user={user} copy={copy} />
            </div>
          ) : null}
        </aside>
      </section>
      <section className="v2-panel pronunciation-hero-v2">
        <span className="eyebrow"><Activity size={15} />{copy("pronunciation", p.title)}</span>
        <h2>{copy("pronunciation_dashboard", p.title)}</h2>
        <p>{copy("pronunciation_dashboard_body", "Карта слабых слов и звуков собирается из ваших голосовых ответов и проверок произношения.")}</p>
        <div className="pronunciation-score-v2">
          <strong>{score}/100</strong>
          <span>{user.level || "A1"} · {copy("latest_score", "latest score")}</span>
        </div>
      </section>
      <section className="v2-panel heatmap-panel-v2">
        <span className="eyebrow">{copy("pronunciation_heatmap", "Heatmap")}</span>
        <h2>{copy("weak_words", "Карта произношения")}</h2>
        <div className="heatmap-grid-v2">
          {(tokens.length ? tokens : [{ word: copy("no_voice_data", "Пока нет голосовых данных"), issue: copy("start_pronunciation_hint", "Запишите фразу выше и отправьте голос") }]).map((item, index) => (
            <PronunciationHeatmapToken key={`${item.word || item.spoken || index}-${index}`} item={item} index={index} user={user} copy={copy} />
          ))}
        </div>
        <div className="pronunciation-focus-v2">
          {focusItems.map((item) => (
            <span key={`${item.word}-${item.target}`}>
              <strong>{item.word}</strong>
              <small>{item.target}</small>
            </span>
          ))}
        </div>
      </section>
      <section className="v2-panel pronunciation-history-v2">
        <span className="eyebrow"><Clock size={15} />{copy("history", "History")}</span>
        <h2>{copy("score_history", "История прогресса")}</h2>
        <div className="pronunciation-history-v2__rows">
          {reports.slice(0, 5).map((report, index) => (
            <div key={index}>
              <span>{copy("attempt", "Попытка")} {index + 1}</span>
              <strong>{normalizedPronunciationScore(report)}/100</strong>
              <small>{pronunciationHistoryText(report, copy)}</small>
            </div>
          ))}
          {!reports.length ? <p>{copy("no_pronunciation_history", "История появится после голосовых ответов.")}</p> : null}
        </div>
      </section>
    </div>
  );
}

function PronunciationHeatmapToken({
  item,
  index,
  user,
  copy,
}: {
  item: PronunciationProblem;
  index: number;
  user: UserProfile;
  copy: (key: string, fallback: string) => string;
}) {
  const label = pronunciationProblemLabel(item, copy);
  const issue = cleanAppText(item.issue || "");
  const tip = cleanAppText(item.tip || "");
  return (
    <div className="heatmap-token-v2" data-intensity={index % 3}>
      <strong>{label.word}</strong>
      {label.confidence ? <small>{label.confidence}</small> : null}
      {issue ? <em>{issue}</em> : null}
      {tip ? <span>{tip}</span> : null}
    </div>
  );
}

function buildOfflineDeck(vocabulary: VocabularyItem[], mistakes: MistakeItem[], phrasebook: PhrasebookItem[] = []): OfflineDeckItem[] {
  const words = vocabulary.slice(0, 20).map((item) => ({
    id: `vocab-${item.id || item.word}`,
    title: item.word,
    answer: item.translation,
    example: item.example || item.context,
    source: "vocabulary" as const,
  }));
  const repairs = mistakes.slice(0, 20).map((item, index) => ({
    id: `mistake-${item.index ?? index}`,
    title: item.word || item.context || "",
    answer: item.correction,
    example: item.explanation,
    source: "mistakes" as const,
  }));
  const phrases = phrasebook.slice(0, 30).map((item) => ({
    id: `phrase-${item.id}`,
    title: item.phrase,
    answer: item.translation || item.note,
    example: item.source,
    source: "phrasebook" as const,
  }));
  return [...words, ...phrases, ...repairs].filter((item) => item.title);
}

function offlineSourceExportLabel(source: OfflineDeckItem["source"]) {
  if (source === "vocabulary") return "Words";
  if (source === "phrasebook") return "Notes";
  return "Mistakes";
}

function offlineSourceUiLabel(source: OfflineDeckItem["source"], copy: (key: string, fallback: string) => string) {
  if (source === "vocabulary") return copy("words", "Слова");
  if (source === "phrasebook") return copy("phrasebook", "Заметки");
  return copy("mistakes", "Ошибки");
}

function formatOfflineDeckTxt(deck: OfflineDeckItem[], copy: (key: string, fallback: string) => string) {
  if (!deck.length) return copy("no_offline_cards", "No offline cards yet.");
  return deck
    .map((item, index) => {
      const lines = [
        `${index + 1}. ${item.title}`,
        `Source: ${offlineSourceExportLabel(item.source)}`,
      ];
      if (item.answer) lines.push(`Answer: ${item.answer}`);
      if (item.example) lines.push(`Example: ${item.example}`);
      return lines.join("\n");
    })
    .join("\n\n");
}

function OfflineDecksView({ vocabulary, mistakes, phrasebook, loadVocabulary, loadMistakes, busy, copy }: ViewRendererProps) {
  const [deck, setDeck] = useState<OfflineDeckItem[]>(() => {
    try {
      return JSON.parse(localStorage.getItem("poliglot-offline-deck-v2") || "[]") as OfflineDeckItem[];
    } catch {
      return [];
    }
  });
  const [page, setPage] = useState(0);
  const sourceDeck = deck.length ? deck : buildOfflineDeck(vocabulary, mistakes, phrasebook);
  const visibleDeck = sourceDeck;
  const pageSize = 10;
  const totalPages = Math.max(1, Math.ceil(visibleDeck.length / pageSize));
  const safePage = Math.min(page, totalPages - 1);
  const pagedDeck = visibleDeck.slice(safePage * pageSize, safePage * pageSize + pageSize);
  useEffect(() => {
    setPage(0);
  }, [deck.length]);
  const refreshDeck = async () => {
    const [nextVocabulary, nextMistakes] = await Promise.all([loadVocabulary(0, { navigate: false }), loadMistakes({ navigate: false })]);
    const next = buildOfflineDeck(nextVocabulary.length ? nextVocabulary : vocabulary, nextMistakes.length ? nextMistakes : mistakes, phrasebook);
    setDeck(next);
    localStorage.setItem("poliglot-offline-deck-v2", JSON.stringify(next));
  };
  const removeDeckItem = (itemId: string) => {
    const next = sourceDeck.filter((item) => item.id !== itemId);
    setDeck(next);
    localStorage.setItem("poliglot-offline-deck-v2", JSON.stringify(next));
  };
  const downloadDeck = (format: "txt" | "json") => {
    const exportDeck = visibleDeck;
    const content = format === "txt" ? formatOfflineDeckTxt(exportDeck, copy) : JSON.stringify(exportDeck, null, 2);
    const blob = format === "txt"
      ? new Blob(["\uFEFF", content], { type: "text/plain;charset=utf-8" })
      : new Blob([content], { type: "application/json;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = `poliglot-offline-deck.${format}`;
    anchor.click();
    URL.revokeObjectURL(url);
  };
  return (
    <div className="offline-decks-v2">
      <section className="v2-panel offline-deck-list-v2">
        <div className="offline-deck-head-v2">
          <div className="offline-deck-summary-v2">
            <span className="eyebrow">{copy("saved_offline", "Saved offline")}</span>
            <h2>{visibleDeck.length ? `${visibleDeck.length} ${copy("cards", "cards")}` : copy("no_offline_cards", "Пока нет карточек.")}</h2>
          </div>
          <div className="offline-hero-v2">
            <span className="eyebrow"><WifiOff size={15} />{copy("offline_decks", "Offline decks")}</span>
            <h2>{copy("offline_decks_title", "PWA mini-decks для повторения")}</h2>
            <p>{copy("offline_decks_body", "Мини-колоды сохраняются в браузере и доступны без Telegram. Service worker кэширует web shell.")}</p>
            <div className="home-insights-v2__actions">
              <Button onClick={() => void refreshDeck()} disabled={busy === "vocabulary" || busy === "mistakes"}><Repeat2 size={16} />{copy("refresh_deck", "Обновить колоду")}</Button>
              <Button variant="outline" onClick={() => downloadDeck("txt")} disabled={!visibleDeck.length}><Download size={16} />TXT</Button>
              <Button variant="outline" onClick={() => downloadDeck("json")} disabled={!visibleDeck.length}><Download size={16} />JSON</Button>
            </div>
          </div>
        </div>
        <div className="offline-pagination-v2">
          <span>{copy("showing_cards_page", "Показано по 10")}: {visibleDeck.length ? `${safePage + 1}/${totalPages}` : "0/0"}</span>
          <div>
            <Button variant="outline" size="sm" disabled={safePage <= 0} onClick={() => setPage((current) => Math.max(0, current - 1))}>{copy("back", "Назад")}</Button>
            <Button variant="outline" size="sm" disabled={safePage + 1 >= totalPages} onClick={() => setPage((current) => Math.min(totalPages - 1, current + 1))}>{copy("offline_next_page", "Вперёд")}</Button>
          </div>
        </div>
        <div className="offline-deck-grid-v2">
          {pagedDeck.map((item) => (
            <article key={item.id} data-source={item.source}>
              <small>{offlineSourceUiLabel(item.source, copy)}</small>
              <strong>{item.title}</strong>
              {item.answer ? <span>{item.answer}</span> : null}
              {item.example ? <p>{item.example}</p> : null}
              <Button
                className="offline-card-delete-v2"
                variant="ghost"
                size="sm"
                type="button"
                onClick={() => removeDeckItem(item.id)}
                aria-label={copy("delete_card", "Удалить карточку")}
              >
                <Trash2 size={15} />
              </Button>
            </article>
          ))}
          {!visibleDeck.length ? <p className="empty-copy">{copy("no_offline_cards", "Пока нет карточек.")}</p> : null}
        </div>
      </section>
    </div>
  );
}

function TeacherDashboardView({ user, session, vocabulary, mistakes, copy }: ViewRendererProps) {
  const learnedWords = Number(user.learned_words || vocabulary.length || 0);
  const lessonCount = Number(user.lesson_count || user.lessons_today || 0);
  const practiceCount = Number(user.practice_count || user.practice_today || 0);
  const reviewCount = Number(user.word_game_count || 0);
  const voiceCount = Number(user.voice_today || 0);
  const mistakeCount = Number(user.mistakes || mistakes.length || 0);
  const weekLessons = habitWeekTotal(user.habit_log, "lessons") || Number(user.lessons_today || 0);
  const weekPractice = habitWeekTotal(user.habit_log, "practice") || Number(user.practice_today || 0);
  const weekVoice = habitWeekTotal(user.habit_log, "voice") || Number(user.voice_today || 0);
  const weeklyBalance = [
    { label: copy("words_learned", "Слов изучено"), value: Math.min(learnedWords, 70), goal: 70 },
    { label: copy("phrases_practiced", "Фраз отработано"), value: weekPractice, goal: weeklyGoalFromDailyLimit(user.practice_limit, 2) },
    { label: copy("lessons_completed", "Уроков пройдено"), value: weekLessons, goal: weeklyGoalFromDailyLimit(user.lesson_limit, 1) },
    { label: copy("review_rounds", "Повторений"), value: Math.min(reviewCount, 56), goal: 56 },
    { label: copy("voice_attempts", "Озвучено"), value: weekVoice, goal: weeklyGoalFromDailyLimit(user.voice_limit, 1) },
  ].map((metric) => ({
    ...metric,
    value: Math.max(0, Math.round(Number(metric.value || 0))),
    goal: Math.max(1, Math.round(Number(metric.goal || 1))),
  }));
  const learningMetrics = [
    { label: copy("words_learned", "Слов изучено"), value: learnedWords, icon: Brain },
    { label: copy("phrases_practiced", "Фраз отработано"), value: practiceCount, icon: MessageCircle },
    { label: copy("lessons_completed", "Уроков пройдено"), value: lessonCount, icon: BookOpen },
    { label: copy("review_rounds", "Повторений"), value: reviewCount, icon: Repeat2 },
    { label: copy("voice_attempts", "Озвучено"), value: voiceCount, icon: Volume2 },
    { label: copy("mistakes_repaired", "Ошибок в работе"), value: mistakeCount, icon: Search },
  ];
  const activityRows = [
    { label: copy("lessons_today", "Lessons today"), value: user.lessons_today, max: user.lesson_limit },
    { label: copy("practice_today", "Practice today"), value: user.practice_today, max: user.practice_limit },
    { label: copy("voices_today", "Voice messages today"), value: user.voice_today, max: user.voice_limit },
  ];
  return (
    <div className="teacher-dashboard-v2">
      <section className="v2-panel dashboard-hero-v2">
        <span className="eyebrow"><BarChart3 size={15} />{copy("dashboard", "Dashboard")}</span>
        <h2>{session.account?.login || "NERIVA learner"}</h2>
        <p>{copy("dashboard_learning_body", "Сводка только по обучению и активности: слова, уроки, разговорная практика, голосовые попытки и повторения.")}</p>
      </section>
      <section className="v2-panel dashboard-grid-v2 dashboard-grid-v2--learning">
        <Metric label="XP" value={compactNumber(user.xp || 0)} />
        <Metric label={copy("level_label", "Level")} value={`${user.level || "A1"} / ${user.xp_level || 1}`} />
        <Metric label={copy("premium", "Premium")} value={user.premium ? prettyDate(user.premium_until) : copy("base_plan", "Basic")} />
        <Metric label={copy("referrals", "Referrals")} value={user.referral_count || 0} />
      </section>
      <section className="v2-panel dashboard-learning-v2">
        <span className="eyebrow">{copy("learning_activity", "Learning activity")}</span>
        <h2>{copy("progress_without_leaderboard", "Учебный прогресс")}</h2>
        <div className="dashboard-learning-v2__grid">
          {learningMetrics.map((metric) => {
            const Icon = metric.icon;
            return (
              <article key={metric.label}>
                <Icon size={18} />
                <span>{metric.label}</span>
                <strong>{compactNumber(metric.value)}</strong>
              </article>
            );
          })}
        </div>
      </section>
      <section className="v2-panel dashboard-activity-v2">
        <span className="eyebrow">{copy("activity", "Activity")}</span>
        <h2>{copy("today", "Today")}</h2>
        {activityRows.map((row) => <Meter key={row.label} label={row.label} value={row.value} max={row.max} />)}
      </section>
      <section className="v2-panel dashboard-rhythm-v2">
        <span className="eyebrow">{copy("weekly_rhythm", "Weekly rhythm")}</span>
        <h2>{copy("training_balance", "Баланс тренировки")}</h2>
        <div className="dashboard-rhythm-v2__bars">
          {weeklyBalance.map((metric) => (
            <span key={metric.label} style={{ "--bar": `${Math.min(100, Math.round((metric.value / metric.goal) * 100))}%` } as CSSProperties}>
              <small>{metric.label}<em>{metric.value}/{metric.goal}</em></small>
              <b />
            </span>
          ))}
        </div>
      </section>
    </div>
  );
}
function ChatWorkView({
  mode,
  messages,
  activeLessonTaskId,
  draft,
  setDraft,
  voiceFile,
  imageFile,
  setVoiceFile,
  setImageFile,
  startLesson,
  submitLesson,
  submitPractice,
  startShadowing,
  submitShadowing,
  shadowingTarget,
  busy,
  copy,
  user,
  savePhrase,
  isPhraseSaved,
}: ViewRendererProps & { mode: "lesson" | "practice" | "shadowing" }) {
  const isLesson = mode === "lesson";
  const isPractice = mode === "practice";
  const isShadowing = mode === "shadowing";
  const submit = isLesson ? submitLesson : isPractice ? submitPractice : submitShadowing;
  const scopedMessages = messages.filter((message) => message.meta === mode);
  const shadowingTitle = copy("listening_phrase", "Listening phrase");
  const visibleMessages = isShadowing
    ? scopedMessages.filter((message) => !(message.title === shadowingTitle && cleanAppText(message.body).trim() === cleanAppText(shadowingTarget).trim()))
    : scopedMessages;
  const showOutput = !isShadowing || visibleMessages.length > 0;
  const practiceEmpty = isPractice && !visibleMessages.length;
  const phraseCandidates = isLesson || isPractice ? phraseCandidatesFromMessages(scopedMessages) : [];
  const lessonHasActiveTask = isLesson && Boolean(activeLessonTaskId);
  const lessonHasCompletedAnswer = isLesson && !lessonHasActiveTask && scopedMessages.some((message) => message.tone === "success" && Boolean(message.details?.task_id));
  const lessonComposerLocked = isLesson && !lessonHasActiveTask;
  const lessonActionLabel = lessonHasCompletedAnswer ? copy("ai_tutor_next_lesson", "Next lesson") : copy("new_lesson", "New lesson");
  return (
    <div className={cn("chat-workspace", `chat-workspace--${mode}`, isShadowing && !visibleMessages.length && "chat-workspace--single")} data-work-mode={mode}>
      {showOutput ? (
        <div className="chat-workspace__output">
          {isLesson && !visibleMessages.length ? (
            <section className="v2-panel lesson-empty-v2">
              <span className="eyebrow">{copy("new_lesson_title", "New lesson")}</span>
              <h2>{copy("lesson_waiting_title", "Нажмите «Новый урок»")}</h2>
              <p>{copy("lesson_waiting_body", "Задание появится здесь и не будет перезапускаться при переходе между вкладками.")}</p>
              <Button type="button" onClick={() => void startLesson()} disabled={busy === "lesson"}>
                {busy === "lesson" ? <Spinner size="small" className="button-spinner-v2" /> : <BookOpen size={16} />}
                {copy("new_lesson", "Новый урок")}
              </Button>
            </section>
          ) : practiceEmpty ? (
            <section className="v2-panel practice-empty-v2">
              <span className="eyebrow">{copy("practice", "Практика")}</span>
              <h2>{copy("practice_empty_title", "Напишите фразу для проверки")}</h2>
              <p>{copy("practice_empty_body", "Отправьте короткий текст, вопрос или фото, и AI предложит исправление и пример.")}</p>
            </section>
          ) : (
            <>
              <ChatPanel messages={visibleMessages} copy={copy} targetLanguage={user.learning_language} />
              {lessonHasCompletedAnswer ? (
                <Button className="lesson-panel-next-v2" type="button" onClick={() => void startLesson()} disabled={busy === "lesson"}>
                  {busy === "lesson" ? <Spinner size="small" className="button-spinner-v2" /> : <BookOpen size={16} />}
                  {copy("ai_tutor_next_lesson", "Next lesson")}
                </Button>
              ) : null}
            </>
          )}
        </div>
      ) : null}
      <section className="v2-panel composer-panel-v2">
        {isLesson ? (
          <Button className="lesson-new-button-v2" type="button" variant="outline" size="sm" onClick={() => void startLesson()} disabled={busy === "lesson"}>
            {busy === "lesson" ? <Spinner size="small" className="button-spinner-v2" /> : <BookOpen size={15} />}
            {lessonActionLabel}
          </Button>
        ) : null}
        {isShadowing && shadowingTarget ? (
          <div className="task-box-v2">
            <span>{copy("spoken_model", "Spoken model")}</span>
            <AudioActionRow clips={[{ label: copy("spoken_model", "Spoken model"), text: shadowingTarget }]} />
            <Button className="task-box-v2__next" variant="outline" size="sm" type="button" onClick={() => void startShadowing()} disabled={busy === "shadowing"}>
              <ChevronRight size={16} />
              {copy("next_phrase", "Next phrase")}
            </Button>
          </div>
        ) : null}
        {!lessonComposerLocked ? (
          <>
            <div className="panel-head">
              <span className="eyebrow">{copy("input", "Input")}</span>
              <h2>{isLesson ? copy("lesson_answer", "Lesson answer") : isPractice ? copy("practice", "Practice") : copy("listening_answer", "Listening answer")}</h2>
            </div>
            <div className="composer-textarea-shell-v2">
              <textarea
                value={draft}
                onChange={(event) => setDraft(event.target.value)}
                onKeyDown={(event) => {
                  if (event.key === "Enter" && !event.shiftKey) {
                    event.preventDefault();
                    void submit();
                  }
                }}
                placeholder={isPractice ? copy("practice_prompt", "Write a phrase, question, or short text for practice.") : isShadowing ? copy("listening_placeholder", "Type the phrase if voice is unavailable.") : copy("lesson_answer_placeholder", "Write your answer to the active lesson.")}
              />
              <Button className="composer-submit-v2" onClick={() => void submit()} disabled={busy === mode} aria-label={copy("send", "Send")}>
                {busy === mode ? <Spinner size="small" className="button-spinner-v2" /> : <Send size={18} />}
                <span>{copy("send", "Send")}</span>
              </Button>
            </div>
            <FileControls voiceFile={voiceFile} imageFile={imageFile} setVoiceFile={setVoiceFile} setImageFile={setImageFile} allowImage={isPractice} copy={copy} />
            {phraseCandidates.length ? <PhraseQuickSave candidates={phraseCandidates} savePhrase={savePhrase} isPhraseSaved={isPhraseSaved} copy={copy} /> : null}
          </>
        ) : null}
      </section>
    </div>
  );
}

function audioClipsFrom(details: ApiRecord | undefined, body: string, copy: (key: string, fallback: string) => string) {
  if (!details) return [];
  const spokenModelLabel = copy("spoken_model", "Spoken model");
  const questionAudioLabel = copy("question_audio", "Question audio");
  const clips = [
    { label: spokenModelLabel, text: recordField(details, ["correction_audio_text", "target", "phrase"]) },
    { label: questionAudioLabel, text: recordField(details, ["question_audio_text", "prompt"]) },
  ].filter((clip) => clip.text && !(clip.label === spokenModelLabel && body.includes(`${spokenModelLabel}: ${clip.text}`)));
  return clips.filter((clip, index, all) => all.findIndex((other) => other.text === clip.text) === index);
}

function ChatPanel({ messages, copy, targetLanguage }: { messages: ChatMessage[]; copy: (key: string, fallback: string) => string; targetLanguage?: string }) {
  const chatConfig: ChatConfig = {
    leftPerson: { name: "NERIVA", avatar: "/app/assets/brand-logo-mini.png" },
    rightPerson: { name: copy("you", "You") },
    targetLanguage,
    messages: messages
      .map((message) => ({
        id: message.id,
        sender: message.role === "user" ? "right" : "left",
        type: "text",
        content: message.body,
        audio: audioClipsFrom(message.details as ApiRecord | undefined, message.body, copy).map((clip) => ({ ...clip, targetLanguage })),
      })),
  };
  return (
    <div className="terminal-chat-v2">
      <ChatComponent config={chatConfig} uiConfig={{ className: "chat-panel-v2" }} />
    </div>
  );
}

function PhraseQuickSave({
  candidates,
  savePhrase,
  isPhraseSaved,
  copy,
}: {
  candidates: Array<{ phrase: string; source: PhrasebookSource; details?: ApiRecord }>;
  savePhrase: (phrase: string, source: PhrasebookSource, details?: ApiRecord) => void;
  isPhraseSaved?: (phrase: string) => boolean;
  copy: (key: string, fallback: string) => string;
}) {
  return (
    <section className="phrase-quick-save-v2">
      <span className="eyebrow"><Bookmark size={14} />{copy("save_to_phrasebook", "Сохранить в заметки")}</span>
      <div className="phrase-quick-save-v2__chips">
        {candidates.slice(0, 4).map((item) => {
          const saved = Boolean(isPhraseSaved?.(item.phrase));
          return (
            <button key={item.phrase} className={cn(saved && "is-saved")} type="button" aria-pressed={saved} onClick={() => savePhrase(item.phrase, item.source, item.details)}>
              <Bookmark size={14} fill={saved ? "currentColor" : "none"} />
              <span>{item.phrase}</span>
            </button>
          );
        })}
      </div>
    </section>
  );
}

function TutorNoteStrip({
  candidates,
  savePhrase,
  isPhraseSaved,
  copy,
}: {
  candidates: Array<{ phrase: string; source: PhrasebookSource; details?: ApiRecord }>;
  savePhrase: (phrase: string, source: PhrasebookSource, details?: ApiRecord) => void;
  isPhraseSaved?: (phrase: string) => boolean;
  copy: (key: string, fallback: string) => string;
}) {
  return (
    <section className="tutor-note-strip-v2" aria-label={copy("save_to_phrasebook", "Save to notes")}>
      <span className="tutor-note-strip-v2__label"><Bookmark size={14} />{copy("save_to_phrasebook", "Save to notes")}</span>
      <div className="tutor-note-strip-v2__chips">
        {candidates.map((item) => {
          const saved = Boolean(isPhraseSaved?.(item.phrase));
          return (
            <button key={`${item.source}-${item.phrase}`} className={cn(saved && "is-saved")} type="button" aria-pressed={saved} onClick={() => savePhrase(item.phrase, item.source, item.details)}>
              <Bookmark size={14} fill={saved ? "currentColor" : "none"} />
              <span>{item.phrase}</span>
            </button>
          );
        })}
      </div>
    </section>
  );
}

function RecordDetails({
  record,
  copy,
  user,
  showContext = true,
  showExampleText = true,
  showWordAudio = true,
  showAudio = true,
  showMistakes = true,
}: {
  record: ApiRecord;
  copy: (key: string, fallback: string) => string;
  user?: UserProfile;
  showContext?: boolean;
  showExampleText?: boolean;
  showWordAudio?: boolean;
  showAudio?: boolean;
  showMistakes?: boolean;
}) {
  const p = pronunciationFrom(record);
  const mistakes = showMistakes ? recordMistakes(record.mistakes) : [];
  const example = recordField(record, ["example"]);
  const translation = recordField(record, ["translation"]);
  const exampleTranslation = recordField(record, ["example_translation", "exampleTranslation", "context_translation", "translation_ru"]);
  const context = showContext ? recordField(record, ["context"]) : "";
  const wordId = recordField(record, ["word_id"]);
  const word = recordField(record, ["word"]);
  const correctionAudio = recordField(record, ["correction_audio_text", "correction"]);
  const promoted = recordField(record, ["promoted_to"]);
  const score = recordField(record, ["score"]);
  const maxScore = recordField(record, ["max_score", "max"]) || "100";
  const audioClips = showAudio
    ? ([
        showWordAudio && wordId && word ? { label: copy("word", "Word"), text: word, wordId } : null,
        example ? { label: copy("example", "Example"), text: example } : null,
        correctionAudio && !word ? { label: copy("correct_variant", "Correct variant"), text: correctionAudio } : null,
      ].filter(Boolean) as Array<{ label: string; text: string; wordId?: string }>)
    : [];
  const showWordTranslation = Boolean(word && translation && !samePhrasebookText(word, translation));
  const showExampleTranslation = Boolean(example && exampleTranslation && !samePhrasebookText(example, exampleTranslation));
  const reportUser = user || ({ interface_language: "ru" } as UserProfile);
  if (!p && !mistakes.length && !example && !context && !promoted && !score && !audioClips.length) return null;
  return (
    <section className="record-details-v2">
      {score ? (
        <div className="detail-chip-v2 is-score">
          <span>{copy("score", "Score")}</span>
          <strong>{score}/{maxScore}</strong>
        </div>
      ) : null}
      {promoted ? (
        <div className="detail-chip-v2">
          <span>{copy("level_updated", "Level updated")}</span>
          <strong>{promoted}</strong>
        </div>
      ) : null}
      {context ? <p><strong>{copy("context", "Context")}:</strong> {context}</p> : null}
      {showWordTranslation || (showExampleText && example) ? (
        <div className="record-audio-text-v2">
          {showWordTranslation ? (
            <p>
              <strong>{copy("word", "Word")}:</strong>
              <span>{word} — {translation}</span>
            </p>
          ) : null}
          {showExampleText && example ? (
            <p>
              <strong>{copy("example", "Example")}:</strong>
              <span>{example}</span>
              {showExampleTranslation ? <em>{exampleTranslation}</em> : null}
            </p>
          ) : null}
        </div>
      ) : null}
      <AudioActionRow clips={audioClips} />
      {mistakes.length ? (
        <div className="detail-list-v2">
          <strong>{copy("mistakes", "Mistakes")}</strong>
          {mistakes.slice(0, 4).map((item, index) => (
            <span key={index}>{cleanAppText(item.word)} → {cleanAppText(item.correction)}</span>
          ))}
        </div>
      ) : null}
      {p ? <PronunciationReport pronunciation={p} user={reportUser} copy={copy} compact /> : null}
    </section>
  );
}

function AudioActionRow({ clips }: { clips: Array<{ label: string; text: string; wordId?: string; targetLanguage?: string }> }) {
  const uniqueClips = clips
    .map((clip) => {
      const generic = /^(Раздел|Section|Mục)$/i.test(cleanAppText(clip.label));
      if (generic && !clip.text && !clip.wordId) return null;
      return generic ? { ...clip, label: "Audio" } : clip;
    })
    .filter((clip): clip is { label: string; text: string; wordId?: string; targetLanguage?: string } => clip !== null)
    .filter((clip) => Boolean(clip.text || clip.wordId))
    .filter((clip, index, all) => all.findIndex((other) => other.text === clip.text && other.wordId === clip.wordId) === index);
  if (!uniqueClips.length) return null;
  return (
    <div className="audio-action-row-v2">
      {uniqueClips.map((clip, index) => (
        <AudioWaveButton key={`${clip.label}-${clip.wordId || clip.text || index}-${clip.targetLanguage || ""}`} label={clip.label} text={clip.text} wordId={clip.wordId} targetLanguage={clip.targetLanguage} />
      ))}
    </div>
  );
}

function PronunciationReport({ pronunciation, user, copy, compact = false }: { pronunciation: PronunciationAssessment; user: UserProfile; copy: (key: string, fallback: string) => string; compact?: boolean }) {
  const problems = Array.isArray(pronunciation.problem_words) ? pronunciation.problem_words.slice(0, compact ? 3 : 5) : [];
  const phonemes = Array.isArray(pronunciation.phoneme_issues) ? pronunciation.phoneme_issues.slice(0, compact ? 2 : 5) : [];
  const tips = recordList(pronunciation.tips).map((item) => localizePronunciationCoachText(item, user, copy)).filter(Boolean).slice(0, compact ? 2 : 4);
  const diagnostics = [pronunciation.stress, pronunciation.rhythm, pronunciation.intonation].map((item) => localizePronunciationCoachText(item, user, copy)).filter(Boolean).slice(0, compact ? 2 : 3);
  return (
    <div className={cn("pronunciation-report-v2", compact && "is-compact")}>
      <div className="pronunciation-metrics-v2">
        <Metric label={copy("pronunciation_score", "Pronunciation")} value={`${normalizedPronunciationScore(pronunciation)}/100`} />
        <Metric label={copy("accent_strength", "Accent")} value={`${clampPercentScore(pronunciation.accent_strength, true)}/100`} />
        <Metric label={copy("fluency", "Fluency")} value={`${clampPercentScore(pronunciation.fluency, true)}/100`} />
      </div>
      {pronunciation.feedback ? <p>{localizePronunciationCoachText(pronunciation.feedback, user, copy)}</p> : null}
      {diagnostics.length ? (
        <div className="pronunciation-problems-v2">
          <strong>{copy("pronunciation_diagnostics", "Stress / rhythm / intonation")}</strong>
          {diagnostics.map((item) => <span key={item}>{item}</span>)}
        </div>
      ) : null}
      {problems.length ? (
        <div className="pronunciation-problems-v2">
          <strong>{copy("problem_words", "Weak words")}</strong>
          {problems.map((item, index) => {
            const label = pronunciationProblemLabel(item, copy);
            return (
              <span key={`${item.word || item.spoken || index}-${index}`}>
                {label.word}
                {label.confidence ? ` ${label.confidence}` : ""}
                {` - ${pronunciationProblemText(item, user, copy)}`}
              </span>
            );
          })}
        </div>
      ) : null}
      {phonemes.length ? (
        <div className="pronunciation-problems-v2">
          <strong>{copy("pronunciation_phonemes", "Sound details")}</strong>
          {phonemes.map((item, index) => {
            const sounds = [cleanAppText(item.expected_sound), cleanAppText(item.heard_sound)].filter(Boolean).join(" -> ");
            const label = [cleanAppText(item.word), sounds].filter(Boolean).join(": ");
            const tip = localizePronunciationCoachText(item.tip || "", user, copy);
            return <span key={`${label}-${index}`}>{label || copy("pronunciation_sound_focus", "Sound focus")}{tip ? ` - ${tip}` : ""}</span>;
          })}
        </div>
      ) : null}
      {tips.length ? (
        <div className="pronunciation-problems-v2">
          <strong>{copy("pronunciation_tips", "Tips")}</strong>
          {tips.map((tip) => <span key={tip}>{tip}</span>)}
        </div>
      ) : null}
    </div>
  );
}

function ChoiceTrainer({ mode, wordChallenge, wordResult, wordGameResult, startWord, startWordGame, answerWord, answerWordGame, reportWord, savePhrase, isPhraseSaved, busy, copy, user }: ViewRendererProps & { mode: "words" | "word-game" }) {
  const challenge = wordChallenge;
  const [wordReportOpen, setWordReportOpen] = useState(false);
  const [wordReportWord, setWordReportWord] = useState("");
  const [wordReportTranslation, setWordReportTranslation] = useState("");
  const [wordReportComment, setWordReportComment] = useState("");
  const [wordReportError, setWordReportError] = useState("");
  const options = useMemo(
    () => tutorStableShuffle(toChoiceOptions(challenge?.options), `${mode}:${challenge?.prompt || ""}:${challenge?.correct_answer_id || ""}`),
    [mode, challenge?.prompt, challenge?.correct_answer_id, challenge?.options],
  );
  const start = mode === "words" ? startWord : startWordGame;
  const answer = mode === "words" ? answerWord : answerWordGame;
  const result = mode === "words" ? wordResult : wordGameResult;
  const title = mode === "words" ? copy("learn_words", "Learn words") : copy("word_game", "Review game");
  const description = mode === "words"
    ? copy("words_section_hint", "Choose the correct translation, listen to examples, and save useful phrases from new words.")
    : copy("word_game_section_hint", "Review learned words from memory and reinforce the ones that are easy to forget.");
  useEffect(() => {
    if (!challenge && !result && !busy) void start();
  }, [mode, Boolean(challenge), Boolean(result), Boolean(busy)]);
  useEffect(() => {
    setWordReportOpen(false);
    setWordReportError("");
  }, [mode, challenge?.word_id]);
  const canReportWord = mode === "words" && !result && Boolean(challenge?.word_id && challenge?.reportable && !challenge?.empty);
  const openWordReport = () => {
    setWordReportWord(challenge?.word || "");
    setWordReportTranslation(challenge?.translation || challenge?.prompt || "");
    setWordReportComment("");
    setWordReportError("");
    setWordReportOpen(true);
  };
  const submitWordReport = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const wordId = challenge?.word_id || "";
    const proposedWord = wordReportWord.trim();
    const proposedTranslation = wordReportTranslation.trim();
    if (!wordId || !proposedWord || !proposedTranslation) {
      setWordReportError(copy("ai_tutor_word_report_required", "Fill in the word and translation."));
      return;
    }
    const ok = await reportWord({
      wordId,
      proposedWord,
      proposedTranslation,
      comment: wordReportComment.trim(),
    });
    if (!ok) return;
    setWordReportOpen(false);
    setWordReportWord("");
    setWordReportTranslation("");
    setWordReportComment("");
    void start();
  };
  return (
    <>
      <section className="v2-panel trainer-display">
        <span className="eyebrow">{title}</span>
        <h2>{challenge?.empty ? copy("no_words_ready", "No words ready") : challenge?.prompt || (busy ? copy("loading", "Loading...") : copy("start_new_round", "Start a new round"))}</h2>
        {!result ? <p>{challenge?.context || description}</p> : null}
        {canReportWord ? (
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="word-report-action-v2 tutor-word-report-action-v2"
            onClick={openWordReport}
            disabled={busy === "word-report"}
          >
            <AlertCircle className="h-4 w-4" />
            {copy("ai_tutor_word_report_button", "Сообщить об ошибке")}
          </Button>
        ) : null}
        {result ? <TrainerResultBox result={result} copy={copy} user={user} onNext={() => void start()} savePhrase={savePhrase} isPhraseSaved={isPhraseSaved} /> : null}
        {options.length ? (
          <div className="choice-grid-v2">
            {options.map((option) => (
              <button
                key={option.id}
                type="button"
                className={cn(
                  result?.selectedId === option.id && "is-selected",
                  result?.correctId === option.id && result.tone === "success" && "is-correct",
                  result?.selectedId === option.id && result.tone === "warning" && "is-wrong",
                )}
                onClick={() => void answer(option.id)}
                disabled={Boolean(busy) || result?.tone === "success"}
              >
                {option.text}
              </button>
            ))}
          </div>
        ) : (
          <div className="trainer-loading-v2">
            <Spinner size="small" show={Boolean(busy || (!challenge && !result))} />
          </div>
        )}
      </section>
      <Dialog open={wordReportOpen} onOpenChange={setWordReportOpen}>
        <DialogContent className="tutor-word-report-dialog-v2">
          <DialogTitle>
            <AlertCircle className="h-5 w-5 tutor-word-report-dialog-v2__icon" />
            {copy("ai_tutor_word_report_title", "Сообщить об ошибке в слове")}
          </DialogTitle>
          <DialogDescription>
            {copy("ai_tutor_word_report_body", "Модератор проверит исправление перед заменой слова.")}
          </DialogDescription>
          <form id="word-report-form-v2" className="tutor-word-report-dialog-v2__body" onSubmit={submitWordReport}>
            <label className="tutor-word-report-field-v2">
              <span>{copy("ai_tutor_word_report_word_label", "Correct word")}</span>
              <input
                value={wordReportWord}
                onChange={(event) => setWordReportWord(event.target.value)}
                maxLength={80}
                autoFocus
              />
            </label>
            <label className="tutor-word-report-field-v2">
              <span>{copy("ai_tutor_word_report_translation_label", "Correct translation")}</span>
              <input
                value={wordReportTranslation}
                onChange={(event) => setWordReportTranslation(event.target.value)}
                maxLength={160}
              />
            </label>
            <label className="tutor-word-report-field-v2">
              <span>{copy("ai_tutor_word_report_comment_label", "Comment")}</span>
              <textarea
                value={wordReportComment}
                onChange={(event) => setWordReportComment(event.target.value)}
                rows={3}
                maxLength={500}
              />
            </label>
            {wordReportError ? <p className="form-error-v2">{wordReportError}</p> : null}
          </form>
          <DialogFooter className="tutor-word-report-dialog-v2__footer">
            <DialogClose asChild>
              <Button type="button" variant="ghost">{copy("cancel", "Cancel")}</Button>
            </DialogClose>
            <Button type="submit" form="word-report-form-v2" disabled={busy === "word-report"}>
              {busy === "word-report" ? copy("loading", "Sending...") : copy("ai_tutor_word_report_submit", "Send report")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

function TrainerResultBox({
  result,
  copy,
  user,
  onNext,
  savePhrase,
  isPhraseSaved,
}: {
  result: TrainerResult;
  copy: (key: string, fallback: string) => string;
  user?: UserProfile;
  onNext?: () => void;
  savePhrase?: (phrase: string, source: PhrasebookSource, details?: ApiRecord) => void;
  isPhraseSaved?: (phrase: string) => boolean;
}) {
  const canAdvance = Boolean(result.mode !== "level" && onNext && (result.tone === "success" || (result.record && recordField(result.record, ["gave_up", "word_id"]))));
  const compactWordResult = result.mode === "words" || result.mode === "word-game";
  const compactSpellingResult = result.mode === "spelling";
  const hideSpellingAnswer = compactSpellingResult && result.tone !== "success" && !result.record?.gave_up;
  const spellingAnswer = compactSpellingResult && !hideSpellingAnswer && result.record ? recordField(result.record, ["correct_answer", "word"]) : "";
  const learnedWord = compactWordResult && result.tone === "success" && result.record ? recordField(result.record, ["word", "correct_answer"]) : "";
  const learnedTranslation = compactWordResult && result.record ? recordField(result.record, ["translation", "prompt", "message"]) : "";
  const phrasebookLabel = [learnedWord, learnedTranslation].filter(Boolean).join(" — ");
  const learnedWordSaved = Boolean(learnedWord && isPhraseSaved?.(learnedWord));
  return (
    <article className={cn("trainer-result-v2", `is-${result.tone}`)}>
      <div>
        <span className="eyebrow">{copy("result", "Result")}</span>
        <h3>{result.title}</h3>
      </div>
      {result.body ? <p>{result.body}</p> : null}
      {learnedWord ? (
        <p className="trainer-correct-word-v2">
          <strong>{copy("word", "Word")}:</strong>
          <span>{learnedWord}</span>
        </p>
      ) : null}
      {spellingAnswer ? (
        <p className="trainer-correct-word-v2">
          <strong>{copy("correct_spelling", "Correct spelling")}:</strong>
          <span>{spellingAnswer}</span>
        </p>
      ) : null}
      {result.record ? (
        <RecordDetails
          record={result.record}
          copy={copy}
          user={user}
          showContext={!compactWordResult && !compactSpellingResult}
          showExampleText={!compactSpellingResult}
          showWordAudio={!compactSpellingResult}
          showMistakes={!hideSpellingAnswer}
          showAudio={!(compactWordResult && result.tone !== "success") && !hideSpellingAnswer}
        />
      ) : null}
      {canAdvance ? (
        <Button type="button" onClick={onNext}><ChevronRight size={17} />{copy("next_word", "Next")}</Button>
      ) : null}
      {learnedWord && savePhrase ? (
        <section className="phrase-quick-save-v2 trainer-word-save-v2">
          <div className="phrase-quick-save-v2__chips">
            <button
              type="button"
              className={cn(learnedWordSaved && "is-saved")}
              aria-pressed={learnedWordSaved}
              aria-label={copy("save_to_phrasebook", "Сохранить в заметки")}
              onClick={() => savePhrase(learnedWord, "word", { ...result.record, translation: learnedTranslation, note: "" })}
            >
              <Bookmark size={14} fill={learnedWordSaved ? "currentColor" : "none"} />
              <span>{phrasebookLabel || learnedWord}</span>
            </button>
          </div>
        </section>
      ) : null}
    </article>
  );
}

function SpellingView({ spellingChallenge, spellingResult, startSpelling, answerSpelling, draft, setDraft, busy, copy, user }: ViewRendererProps) {
  useEffect(() => {
    if (!spellingChallenge && !spellingResult && !busy) void startSpelling();
  }, [Boolean(spellingChallenge), Boolean(spellingResult), Boolean(busy)]);
  return (
    <section className="v2-panel trainer-display">
      <span className="eyebrow">{spellingChallenge?.direction || copy("spelling", "Spelling")}</span>
      <h2>{spellingChallenge?.prompt || (busy ? copy("loading", "Loading...") : copy("start_spelling_drill", "Start a spelling drill"))}</h2>
      {spellingChallenge?.context ? <p>{spellingChallenge.context}</p> : null}
      {spellingChallenge ? (
        <div className="inline-form-v2">
          <input
            value={draft}
            onChange={(event) => setDraft(event.target.value)}
            placeholder={copy("type_answer", "Type the answer")}
            onKeyDown={(event) => {
              if (event.key === "Enter") void answerSpelling(false);
            }}
          />
          <Button variant="approve" onClick={() => void answerSpelling(false)} disabled={Boolean(busy)}><Check size={17} />{copy("check", "Check")}</Button>
          <Button variant="reject" onClick={() => void answerSpelling(true)} disabled={Boolean(busy)}>{copy("give_up", "I do not know")}</Button>
        </div>
      ) : (
        <div className="trainer-loading-v2">
          <Spinner size="small" show={Boolean(busy || (!spellingChallenge && !spellingResult))} />
        </div>
      )}
      {spellingResult ? <TrainerResultBox result={spellingResult} copy={copy} user={user} onNext={() => void startSpelling()} /> : null}
    </section>
  );
}

function LevelView({ levelQuestion, levelResult, startLevel, answerLevel, busy, copy, user }: ViewRendererProps) {
  const levelOptions = levelQuestion
    ? levelQuestion.options
      .map((option, index) => ({ option, index }))
      .filter(({ option }) => {
        const normalized = cleanAppText(option).trim().toLowerCase();
        return !["i do not know", "i don't know", "idk", "я не знаю", "не знаю"].includes(normalized);
      })
    : [];
  return (
    <section className="v2-panel trainer-display">
      <span className="eyebrow">{levelQuestion ? `${levelQuestion.index + 1} / ${levelQuestion.total}` : copy("assessment", "Assessment")}</span>
      <h2>{levelQuestion?.question || copy("start_level_test", "Start the level test")}</h2>
      <p>{copy("level_help", "Answers update the Go backend profile and keep the active question order stable.")}</p>
      {levelQuestion ? (
        <div className="level-answer-sheet-v2">
          {levelOptions.map(({ option, index }) => (
            <button key={`${option}-${index}`} type="button" onClick={() => void answerLevel(index)} disabled={Boolean(busy)}>
              <span>{String.fromCharCode(65 + index)}</span>
              {option}
            </button>
          ))}
          <Button className="level-answer-skip-v2" variant="outline" type="button" onClick={() => void answerLevel(-1)} disabled={Boolean(busy)}>
            {copy("skip", "Skip")}
          </Button>
        </div>
      ) : (
        <Button onClick={startLevel} disabled={Boolean(busy)}><Play size={18} />{copy("view_level_action", "Start test")}</Button>
      )}
      {levelResult ? <TrainerResultBox result={levelResult} copy={copy} user={user} /> : null}
    </section>
  );
}

function VocabularyView({ vocabulary, vocabularyMeta, loadVocabulary, busy, copy }: ViewRendererProps) {
  useEffect(() => {
    if (!vocabulary.length && busy !== "vocabulary") void loadVocabulary(0);
  }, []);
  return (
    <section className="v2-panel vocabulary-panel-v2">
      <div className="panel-head">
        <div>
          <span className="eyebrow">{copy("vocabulary", "Vocabulary")}</span>
          <h2>
            {vocabulary.length
              ? `${vocabularyMeta.total} ${copy("words_count", "words")} - ${vocabularyMeta.page + 1}/${vocabularyMeta.total_pages}`
              : copy("no_vocabulary_loaded", "No vocabulary loaded")}
          </h2>
        </div>
        <div className="pagination-v2">
          <MorphingArrowButton
            direction="left"
            label={copy("back", "Back")}
            onClick={() => void loadVocabulary(Math.max(0, vocabularyMeta.page - 1))}
            disabled={busy === "vocabulary" || vocabularyMeta.page <= 0}
            className="vocabulary-arrow-v2"
          />
          <MorphingArrowButton
            direction="right"
            label={copy("next", "Next")}
            onClick={() => void loadVocabulary(vocabularyMeta.page + 1)}
            disabled={busy === "vocabulary" || vocabularyMeta.page + 1 >= vocabularyMeta.total_pages}
            className="vocabulary-arrow-v2"
          />
        </div>
      </div>
      <div className="vocabulary-grid-v2">
        {vocabulary.length ? (
          vocabulary.map((item) => (
            <article className="vocabulary-card-v2" key={item.id || item.word}>
              <div>
                <strong>{item.word}</strong>
                {item.translation ? <span>{item.translation}</span> : null}
                {item.example ? <p className="vocabulary-example-v2">{item.example}</p> : item.context ? <p className="vocabulary-example-v2">{item.context}</p> : null}
              </div>
              <AudioWaveButton label={copy("listen", "Listen")} wordId={item.id} text={item.word} compact />
            </article>
          ))
        ) : (
          <p className="empty-copy">
            {busy === "vocabulary"
              ? copy("loading", "Loading...")
              : copy("empty_panel", "No items yet. Start the related practice mode to fill this panel.")}
          </p>
        )}
      </div>
    </section>
  );
}

function PhrasebookView({ phrasebook, savePhrase, removePhrase, copy, user }: ViewRendererProps) {
  const [manualPhrase, setManualPhrase] = useState("");
  const [manualNote, setManualNote] = useState("");
  const [page, setPage] = useState(0);
  const pageSize = 10;
  const totalPages = Math.max(1, Math.ceil(phrasebook.length / pageSize));
  const safePage = Math.min(page, totalPages - 1);
  const pagedPhrasebook = phrasebook.slice(safePage * pageSize, safePage * pageSize + pageSize);
  useEffect(() => {
    setPage(0);
  }, [phrasebook.length]);
  const addManualPhrase = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const phrase = manualPhrase.trim();
    const note = manualNote.trim();
    const text = phrase || note;
    if (!text) return;
    savePhrase(text, "manual", phrase && note ? { note } : {});
    setManualPhrase("");
    setManualNote("");
  };
  return (
    <div className="phrasebook-view-v2">
      <section className="v2-panel phrasebook-hero-v2">
        <span className="eyebrow"><Bookmark size={15} />{copy("phrasebook", "Phrasebook")}</span>
        <h2>{copy("phrasebook_title", "Избранные фразы из уроков")}</h2>
        <p>{copy("phrasebook_body", "Сохраняйте удачные фразы из уроков и практики, чтобы повторять их как личный разговорник.")}</p>
        <form className="phrasebook-add-v2" onSubmit={addManualPhrase}>
          <input
            value={manualPhrase}
            onChange={(event) => setManualPhrase(event.target.value)}
            placeholder={copy("phrasebook_phrase_placeholder", "Фраза для личного разговорника")}
          />
          <input
            value={manualNote}
            onChange={(event) => setManualNote(event.target.value)}
            placeholder={copy("phrasebook_note_placeholder", "Заметка или перевод")}
          />
          <Button type="submit"><Plus size={16} />{copy("add", "Добавить")}</Button>
        </form>
      </section>
      <section className="v2-panel phrasebook-list-v2">
        <div className="panel-head">
          <div>
            <span className="eyebrow">{copy("saved_phrases", "Saved phrases")}</span>
            <h2>{phrasebook.length ? `${phrasebook.length} ${copy("phrases", "phrases")}` : copy("empty_phrasebook", "Phrasebook is empty")}</h2>
          </div>
          <Button variant="outline" onClick={() => savePhrase(copy("sample_phrase", "Could you say that again, please?"), "manual")}>
            <Bookmark size={16} />
            {copy("add_sample", "Добавить пример")}
          </Button>
        </div>
        {phrasebook.length > pageSize ? (
          <div className="offline-pagination-v2 phrasebook-pagination-v2">
            <span>{copy("showing_cards_page", "Показано по 10")}: {phrasebook.length ? `${safePage + 1}/${totalPages}` : "0/0"}</span>
            <div>
              <Button variant="outline" size="sm" disabled={safePage <= 0} onClick={() => setPage((current) => Math.max(0, current - 1))}>{copy("back", "Назад")}</Button>
              <Button variant="outline" size="sm" disabled={safePage + 1 >= totalPages} onClick={() => setPage((current) => Math.min(totalPages - 1, current + 1))}>{copy("next", "Вперёд")}</Button>
            </div>
          </div>
        ) : null}
        <div className="phrasebook-grid-v2">
          {phrasebook.length ? (
            pagedPhrasebook.map((item) => (
              <article className="phrasebook-card-v2" key={item.id}>
                <small>{copy(`phrase_source_${item.source || "manual"}`, item.source || "manual")} · {prettyDate(item.createdAt || item.created_at)}</small>
                <strong>{item.phrase}</strong>
                <small className="phrasebook-card-v2__group">{(item.source || "manual") === "manual" ? copy("phrasebook_group_my", languageCode(user.interface_language) === "ru" ? "Мои заметки" : "My notes") : copy("phrasebook_group_learning", languageCode(user.interface_language) === "ru" ? "Из разделов" : "From sections")}</small>
                {item.translation && !samePhrasebookText(item.translation, item.phrase) ? <p>{item.translation}</p> : null}
                {item.note && !samePhrasebookText(item.note, item.phrase) && !samePhrasebookText(item.note, item.translation) ? <p>{item.note}</p> : null}
                <div className="phrasebook-card-v2__actions">
                  <AudioWaveButton label={copy("listen", "Listen")} text={item.phrase} targetLanguage={item.language || user.learning_language} compact />
                  <Button type="button" variant="outline" size="sm" onClick={() => removePhrase(item.id)}>
                    <X size={14} />
                    {copy("remove", "Remove")}
                  </Button>
                </div>
              </article>
            ))
          ) : (
            <p className="empty-copy">{copy("phrasebook_empty_hint", "Во время урока нажмите кнопку сохранения рядом с полезной фразой.")}</p>
          )}
        </div>
      </section>
    </div>
  );
}

function MistakesView({ mistakes, loadMistakes, startMistakePractice, submitMistakeAnswer, cancelMistakePractice, deleteMistake, clearMistakes, mistakePractice, mistakePracticeIndex, mistakeAnswer, setMistakeAnswer, mistakeResult, busy, copy, user }: ViewRendererProps) {
  const [confirmClear, setConfirmClear] = useState(false);
  const [activeCategory, setActiveCategory] = useState<MistakeCategory | "all">("all");
  const [page, setPage] = useState(0);
  const categories: MistakeCategory[] = ["grammar", "word-order", "vocabulary", "politeness", "spelling"];
  const grouped = useMemo(() => {
    const seed = categories.reduce((acc, category) => ({ ...acc, [category]: [] }), {} as Record<MistakeCategory, Array<{ item: MistakeItem; index: number }>>);
    mistakes.forEach((item, index) => {
      seed[mistakeCategory(item)].push({ item, index });
    });
    return seed;
  }, [mistakes]);
  const visibleCategories = categories.filter((category) => grouped[category].length > 0);
  const visibleMistakes = activeCategory === "all"
    ? mistakes.map((item, index) => ({ item, index }))
    : grouped[activeCategory];
  const pageSize = 10;
  const totalPages = Math.max(1, Math.ceil(visibleMistakes.length / pageSize));
  const safePage = Math.min(page, totalPages - 1);
  const pagedMistakes = visibleMistakes.slice(safePage * pageSize, safePage * pageSize + pageSize);

  useEffect(() => {
    if (!mistakes.length && busy !== "mistakes") void loadMistakes({ navigate: false });
  }, []);
  useEffect(() => {
    setPage(0);
  }, [activeCategory, mistakes.length]);
  useEffect(() => {
    if (activeCategory !== "all" && grouped[activeCategory].length === 0) setActiveCategory("all");
  }, [activeCategory, grouped]);

  const itemIndex = (item: MistakeItem, index: number) => item.index ?? index;
  const startPracticeFor = (item: MistakeItem, index: number) => startMistakePractice(itemIndex(item, index));
  const trainSimilar = () => {
    const category = mistakePractice ? mistakeCategory(mistakePractice) : activeCategory === "all" ? "grammar" : activeCategory;
    const activeIndex = mistakePracticeIndex ?? (typeof mistakePractice?.index === "number" ? mistakePractice.index : null);
    const pool = grouped[category]?.length ? grouped[category] : visibleMistakes;
    const candidate = pool.find(({ item, index }) => itemIndex(item, index) !== activeIndex)
      || visibleMistakes.find(({ item, index }) => itemIndex(item, index) !== activeIndex)
      || pool[0];
    if (candidate) void startPracticeFor(candidate.item, candidate.index);
  };

  const clearAction = confirmClear ? (
    <div className="confirm-clear-v2">
      <span>{copy("confirm_clear_mistakes", "Clear all saved mistakes?")}</span>
      <Button variant="destructive" size="sm" onClick={() => { setConfirmClear(false); void clearMistakes(); }} disabled={!mistakes.length || busy === "mistakes-clear"}>{copy("yes_clear", "Yes, clear")}</Button>
      <Button variant="outline" size="sm" onClick={() => setConfirmClear(false)}>{copy("cancel", "Cancel")}</Button>
    </div>
  ) : (
    <Button variant="destructive" size="sm" onClick={() => setConfirmClear(true)} disabled={!mistakes.length || busy === "mistakes-clear"}><Trash2 size={15} />{copy("clear", "Clear")}</Button>
  );

  if (mistakePractice) {
    return (
      <div className="mistake-layout-v2 mistake-layout-v2--practice">
        <section className="v2-panel mistake-dictionary-v2 mistake-dictionary-v2--practice mistake-practice-v2">
        <div className="panel-head">
          <div>
            <span className="eyebrow">{copy("mistake_practice", "Mistake practice")}</span>
            <h2>{mistakePractice.word || copy("choose_mistake", "Choose a saved mistake")}</h2>
          </div>
          <div className="panel-actions-v2">
            <Button variant="outline" size="sm" onClick={cancelMistakePractice}>
              <ChevronRight className="flip-icon-v2" size={15} />
              {copy("back_to_mistake_list", "Back to mistakes")}
            </Button>
            <Button variant="outline" size="sm" onClick={trainSimilar} disabled={!mistakes.length || busy === "mistake-practice"}>
              <Repeat2 size={15} />
              {copy("train_similar", "Practice similar")}
            </Button>
          </div>
        </div>
          <div className="mistake-card-v2">
            <span>{mistakeCategoryLabel(mistakeCategory(mistakePractice), copy)}</span>
            <strong>{mistakePractice.word || "-"}</strong>
            {mistakePractice.correction ? <p>{mistakePractice.correction}</p> : null}
            {mistakePractice.explanation ? <small>{mistakePractice.explanation}</small> : null}
            {mistakePractice.correction ? <AudioActionRow clips={[{ label: copy("correct_variant", "Correct variant"), text: mistakePractice.correction }]} /> : null}
          </div>
          <div className="inline-form-v2">
            <input value={mistakeAnswer} onChange={(event) => setMistakeAnswer(event.target.value)} placeholder={copy("type_corrected_variant", "Type the corrected variant")} onKeyDown={(event) => { if (event.key === "Enter") void submitMistakeAnswer(); }} />
            <Button onClick={() => void submitMistakeAnswer()} disabled={busy === "mistake-answer"}><Check size={17} />{copy("check", "Check")}</Button>
          </div>
        {mistakeResult ? <TrainerResultBox result={mistakeResult} copy={copy} user={user} onNext={() => void startMistakePractice()} /> : null}
        </section>
      </div>
    );
  }

  return (
    <div className="mistake-layout-v2 mistake-layout-v2--list">
      <section className="v2-panel mistake-dictionary-v2">
        <div className="panel-head">
          <div>
            <span className="eyebrow">{copy("mistakes", "Mistakes")}</span>
            <h2>{mistakes.length ? `${mistakes.length} ${copy("items", "items")}` : copy("no_mistakes_loaded", "No mistakes loaded")}</h2>
          </div>
          <div className="panel-actions-v2 mistake-header-actions-v2">
            <p className="mistake-click-hint-v2">{copy("mistake_click_hint", "Нажмите на ошибку, чтобы проработать её")}</p>
            <Button variant="outline" size="sm" onClick={() => void loadMistakes({ navigate: false })} disabled={busy === "mistakes"}>
              {busy === "mistakes" ? <Spinner size="small" className="button-spinner-v2" /> : <Search size={15} />}
              {copy("refresh", "Refresh")}
            </Button>
            {clearAction}
          </div>
        </div>
        <div className="mistake-groups-v2" role="tablist" aria-label={copy("mistake_groups", "Mistake groups")}>
          <button type="button" className={activeCategory === "all" ? "is-active" : ""} onClick={() => setActiveCategory("all")}>
            {copy("all", "All")} <span>{mistakes.length}</span>
          </button>
          {visibleCategories.map((category) => (
            <button key={category} type="button" className={activeCategory === category ? "is-active" : ""} onClick={() => setActiveCategory(category)}>
              {mistakeCategoryLabel(category, copy)} <span>{grouped[category].length}</span>
            </button>
          ))}
        </div>
        {visibleMistakes.length > pageSize ? (
          <div className="offline-pagination-v2 mistake-pagination-v2">
            <span>{copy("showing_cards_page", "Показано по 10")}: {visibleMistakes.length ? `${safePage + 1}/${totalPages}` : "0/0"}</span>
            <div>
              <Button variant="outline" size="sm" disabled={safePage <= 0} onClick={() => setPage((current) => Math.max(0, current - 1))}>{copy("back", "Назад")}</Button>
              <Button variant="outline" size="sm" disabled={safePage + 1 >= totalPages} onClick={() => setPage((current) => Math.min(totalPages - 1, current + 1))}>{copy("next", "Вперёд")}</Button>
            </div>
          </div>
        ) : null}
        <div className="mistake-list-v2">
          {visibleMistakes.length ? pagedMistakes.map(({ item, index }) => {
            const category = mistakeCategory(item);
            const indexValue = itemIndex(item, index);
            return (
              <article key={`${item.word || "mistake"}-${index}`} className="mistake-list-card-v2">
                <button className="mistake-card-main-v2" type="button" onClick={() => void startPracticeFor(item, index)}>
                  <span>{mistakeCategoryLabel(category, copy)}</span>
                  <strong>{item.word || item.context || copy("mistake", "Mistake")}</strong>
                  {item.correction ? <em>{item.correction}</em> : null}
                  {item.explanation ? <p>{item.explanation}</p> : item.context ? <p>{item.context}</p> : null}
                </button>
                <div className="mistake-card-actions-v2">
                  {item.correction ? <AudioActionRow clips={[{ label: copy("correct_variant", "Correct variant"), text: item.correction }]} /> : null}
                  <Button className="mistake-delete-v2" variant="outline" size="sm" type="button" onClick={() => void deleteMistake(indexValue)} disabled={busy === "mistake-delete"}>
                    <Trash2 size={15} />
                    {copy("remove", "Remove")}
                  </Button>
                </div>
              </article>
            );
          }) : (
            <p className="empty-copy">{busy === "mistakes" ? copy("loading", "Loading...") : copy("no_group_mistakes", "В этой группе пока нет ошибок.")}</p>
          )}
        </div>
      </section>
    </div>
  );
}

function LeaderboardView({ leaderboard, leaderboardMeta, leaderboardLanguage, setLeaderboardLanguage, loadLeaderboard, busy, copy, session }: ViewRendererProps) {
  const languages = session.learning_languages?.length ? session.learning_languages : session.interface_languages || [];
  useEffect(() => {
    if (!leaderboard.length && busy !== "leaderboard") void loadLeaderboard(leaderboardLanguage);
  }, []);
  return (
    <section className="v2-panel data-display">
      <div className="panel-head">
        <div>
          <span className="eyebrow">{copy("language_leaderboard", "Language leaderboard")}</span>
          <h2>{leaderboardMeta.languageName || copy("top_learners", "Top learners")}</h2>
        </div>
        <div className="panel-actions-v2">
          <Button
            variant={leaderboardLanguage === "global" ? "default" : "outline"}
            onClick={() => {
              setLeaderboardLanguage("global");
              void loadLeaderboard("global");
            }}
            disabled={busy === "leaderboard"}
          >
            <Trophy size={16} />
            {copy("global_top", "Global top")}
          </Button>
          <AnimatedSelect
            className="leaderboard-select-v2"
            value={leaderboardLanguage || leaderboardMeta.language || "global"}
            options={[{ code: "global", native_name: copy("global_top", "Global top") }, ...languages]}
            onChange={(value) => {
              setLeaderboardLanguage(value);
              void loadLeaderboard(value);
            }}
            ariaLabel={copy("language_leaderboard", "Language leaderboard")}
          />
          <Button variant="outline" onClick={() => void loadLeaderboard(leaderboardLanguage || "global")} disabled={busy === "leaderboard"}>{copy("refresh", "Refresh")}</Button>
        </div>
      </div>
      <div className="leaderboard-v2">
        {(leaderboard.length ? leaderboard : [{ name: copy("load_leaderboard", "Load leaderboard"), score: 0 }]).slice(0, 10).map((entry, index) => (
          <div className="leaderboard-row-v2" key={`${entry.name || index}-${index}`}>
            <span>{index + 1}</span>
            <strong>{entry.name || entry.title || copy("learner", "Learner")}<small>{entry.level || ""}</small></strong>
            <em>{compactNumber(entry.rating_points || entry.score || 0)} {copy("rating_points", "баллов")} - {compactNumber(entry.words || 0)} {copy("words_count", "words")} - {compactNumber(entry.mistakes || 0)} {copy("mistakes", "Mistakes")}</em>
            {leaderboardLanguage === "global" ? <small className="leaderboard-languages-v2">{leaderboardLanguages(entry, languages, copy)}</small> : null}
          </div>
        ))}
      </div>
    </section>
  );
}

function PremiumView({ user, premiumPlans, loadPremiumPlans, paymentHistory, busy, activationKey, setActivationKey, activateKey, copy, setPayment }: ViewRendererProps) {
  useEffect(() => {
    void loadPremiumPlans();
  }, [user.interface_language]);
  const plans = premiumPlanCatalog(premiumPlans, copy);
  const premiumUntil = asText(user.premium_until, "");
  const currentPlan = asText(user.plan, "");
  const displayPlan = currentPlan || copy("base_plan", "Basic");
  const isPremiumActive = Boolean(user.premium);
  return (
    <div className="premium-grid-v2">
      <section className="v2-panel premium-status-v2">
        <span className="eyebrow">{copy("current_plan", "Current plan")}</span>
        <h2>{isPremiumActive ? copy("premium", "Premium") : copy("base_plan", "Basic")}</h2>
        <div className="premium-status-v2__meta">
          <p>
            <span>{copy("tariff", "Tariff")}</span>
            <strong>{displayPlan}</strong>
          </p>
          <p>
            <span>{copy("premium_until", "Premium until")}</span>
            <strong>{premiumUntil ? prettyDate(premiumUntil) : copy("not_active", "Not active")}</strong>
          </p>
        </div>
      </section>
      {plans.map((plan) => {
        const isFree = premiumPlanIsFree(plan);
        const isCurrent = premiumPlanIsCurrent(plan, user);
        const canPay = !isFree;
        const planNote = cleanAppText(plan.note);
        return (
        <section className={cn("v2-panel plan-card-v2", isFree && "is-free", isCurrent && "is-current", premiumPlanIsPlatinum(plan) && "is-platinum")} key={plan.product} aria-current={isCurrent ? "true" : undefined}>
          <span>{premiumPlanLabel(plan, copy)}</span>
          <h2>{premiumPlanTitle(plan, copy)}</h2>
          {isCurrent ? (
            <span className="plan-current-badge-v2">
              <CheckCircle size={14} />
              {copy("current_plan", "Current plan")}
            </span>
          ) : null}
          <strong>{isFree ? copy("free_plan_price", "Included") : plan.rub_price ? `${plan.rub_price} RUB` : copy("available", "Available")}</strong>
          <p>{premiumPlanBody(plan, copy)}</p>
          <ul className="plan-features-v2">
            {premiumPlanFeatures(plan, copy).map((feature) => (
              <li key={feature.text} className={cn(feature.locked && "is-locked", feature.highlight && "is-highlight")}>
                {feature.locked ? <X size={15} /> : <Check size={15} />}
                <span>{feature.text}</span>
              </li>
            ))}
          </ul>
          {planNote && <p className="plan-note-v2">{planNote}</p>}
          <Button variant={canPay ? "default" : "outline"} onClick={() => canPay && setPayment({ plan })} disabled={isFree || busy === "premium-plans"}>
            {canPay ? <CircleDollarSign size={18} /> : <ShieldCheck size={18} />}
            {isCurrent ? copy("extend_plan", "Extend") : isFree ? copy("free_plan_price", "Included") : copy("payment_options", "Payment options")}
          </Button>
        </section>
        );
      })}
      <form className="v2-panel activation-v2" onSubmit={activateKey}>
        <span className="eyebrow">{copy("activation_key", "Activation key")}</span>
        <h2>{copy("serial_premium_access", "Serial Premium access")}</h2>
        <div className="inline-form-v2">
          <input value={activationKey} onChange={(event) => setActivationKey(event.target.value)} placeholder="XXXX-XXXX-XXXX-XXXX" />
          <Button className="activation-button-v2" variant="shimmer" type="submit" disabled={busy === "activation"}><Diamond size={17} />{copy("activate", "Activate")}</Button>
        </div>
      </form>
      <section className="v2-panel payment-history-v2">
        <span className="eyebrow">{copy("payment_history", "Payment history")}</span>
        <h2>{copy("payment_history_title", "История платежей")}</h2>
        <div className="payment-history-v2__table">
          <div className="payment-history-v2__head">
            <span>{copy("date", "Дата")}</span>
            <span>{copy("plan", "План")}</span>
            <span>{copy("amount", "Сумма")}</span>
            <span>{copy("method", "Способ")}</span>
            <span>{copy("status", "Статус")}</span>
          </div>
          {paymentHistory.length ? paymentHistory.map((item) => (
            <div className="payment-history-v2__row" key={item.id}>
              <span>{prettyDate(item.date)}</span>
              <strong>{item.plan}<small>{item.period}</small></strong>
              <span>{item.amount || "-"}</span>
              <span>{item.method}</span>
              <em>{item.status}</em>
            </div>
          )) : (
            <p className="empty-copy">{copy("no_payment_history", "Подтверждённых платежей пока нет. Созданные и ожидающие оплаты счета здесь не показываются.")}</p>
          )}
        </div>
      </section>
    </div>
  );
}

function SettingsView({
  session,
  user,
  saveSettings,
  changePassword,
  startTelegramCode,
  verifyTelegramCode,
  busy,
  activationKey,
  setActivationKey,
  activateKey,
  copy,
}: ViewRendererProps) {
  const [interfaceLanguage, setInterfaceLanguage] = useState(user.interface_language || "ru");
  const [learningLanguage, setLearningLanguage] = useState(user.learning_language || "en");
  const [level, setLevel] = useState(user.level || "A1");
  const [learningFocus, setLearningFocus] = useState(user.learning_focus || "");
  const [telegramOtpOpen, setTelegramOtpOpen] = useState(false);
  const [telegramRequest, setTelegramRequest] = useState<TelegramCodeRequest | null>(null);
  const [telegramPrivacyAccepted, setTelegramPrivacyAccepted] = useState(false);
  const [telegramPrivacyError, setTelegramPrivacyError] = useState("");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [newPasswordConfirm, setNewPasswordConfirm] = useState("");
  const telegramLinked = Boolean(user.telegram_account?.id || user.telegram_linked);
  const legalLanguage = languageCode(interfaceLanguage || user.interface_language || "ru");
  const privacyURL = legalDocumentURL("privacy.html", legalLanguage);
  const consentURL = legalDocumentURL("consent.html", legalLanguage);
  const agreementURL = legalDocumentURL("agreement.html", legalLanguage);
  const passwordTooShort = newPassword.length > 0 && newPassword.length < 8;
  const passwordsMatch = newPassword.length > 0 && newPasswordConfirm.length > 0 && newPassword === newPasswordConfirm;
  const passwordsMismatch = newPasswordConfirm.length > 0 && newPassword !== newPasswordConfirm;

  useEffect(() => {
    setInterfaceLanguage(user.interface_language || "ru");
  }, [user.interface_language]);

  useEffect(() => {
    setLearningLanguage(user.learning_language || "en");
  }, [user.learning_language]);

  useEffect(() => {
    setLevel(user.level || "A1");
  }, [user.level]);

  useEffect(() => {
    setLearningFocus(user.learning_focus || "");
  }, [user.learning_focus]);

  const submitPassword = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (passwordTooShort || passwordsMismatch || !currentPassword || !newPassword) return;
    await changePassword({
      current_password: currentPassword,
      new_password: newPassword,
      new_password_confirm: newPasswordConfirm,
    });
    setCurrentPassword("");
    setNewPassword("");
    setNewPasswordConfirm("");
  };
  const openTelegramTwoFactor = async () => {
    if (!telegramPrivacyAccepted) {
      setTelegramPrivacyError(copy("auth_privacy_required", "Accept the personal data policy to continue."));
      return;
    }
    setTelegramPrivacyError("");
    const request = await startTelegramCode({
      privacy_consent: true,
      privacy_policy_url: privacyURL,
      personal_data_consent_url: consentURL,
      user_agreement_url: agreementURL,
    });
    if (!request) return;
    setTelegramRequest(request);
    setTelegramOtpOpen(true);
  };
  const resendTelegramTwoFactor = async () => {
    const request = await startTelegramCode({
      privacy_consent: true,
      privacy_policy_url: privacyURL,
      personal_data_consent_url: consentURL,
      user_agreement_url: agreementURL,
    });
    if (!request) return false;
    setTelegramRequest(request);
    return true;
  };
  const verifyTelegramTwoFactor = async (code: string) => {
    const ok = await verifyTelegramCode(telegramRequest?.token || "", code);
    if (ok) {
      setTelegramOtpOpen(false);
      setTelegramRequest(null);
    }
    return ok;
  };
  const telegramCodeExpiresAt = telegramRequest?.expiresAt ? prettyDate(telegramRequest.expiresAt) : "";
  return (
    <div className="settings-grid-v2">
      <section className="v2-panel settings-card-v2">
        <span className="eyebrow">{copy("profile", "Profile")}</span>
        <h2>{copy("learning_settings", "Learning settings")}</h2>
        <AnimatedSelect label={copy("interface_language", "Interface language")} value={interfaceLanguage} options={session.interface_languages || []} onChange={setInterfaceLanguage} />
        <AnimatedSelect label={copy("learning_language", "Learning language")} value={learningLanguage} options={session.learning_languages || []} onChange={setLearningLanguage} />
        <AnimatedSelect label={copy("level_label", "Level")} value={level} options={["A1", "A2", "B1", "B2", "C1", "C2"].map((code) => ({ code, native_name: code }))} onChange={setLevel} />
        <label className="settings-focus-v2">
          <span>{copy("learning_focus", "Global learning focus")}</span>
          <textarea
            value={learningFocus}
            onChange={(event) => setLearningFocus(event.target.value)}
            placeholder={copy("learning_focus_placeholder", "For example: daily travel, hotel, restaurant, work calls. Leave empty for mixed topics.")}
          />
          <small>{copy("learning_focus_hint", "This changes the topic focus used after onboarding, so new lessons and practice stop repeating one random theme.")}</small>
        </label>
        <Button variant="approve" disabled={busy === "settings"} onClick={() => void saveSettings({ interface_language: interfaceLanguage, learning_language: learningLanguage, level, learning_focus: learningFocus })}><Check size={18} />{copy("save_settings", "Save settings")}</Button>
      </section>
      <form className="v2-panel settings-card-v2 password-card-v2" onSubmit={submitPassword}>
        <span className="eyebrow">{copy("security", "Security")}</span>
        <h2>{copy("change_password", "Change password")}</h2>
        <PasswordInputField
          label={copy("current_password", "Current password")}
          value={currentPassword}
          onChange={(event) => setCurrentPassword(event.target.value)}
          autoComplete="current-password"
        />
        <PasswordInputField
          label={copy("new_password", "New password")}
          value={newPassword}
          onChange={(event) => setNewPassword(event.target.value)}
          autoComplete="new-password"
          invalid={passwordTooShort}
          errorText={copy("password_min_hint", "At least 8 characters")}
          helperText={!passwordTooShort ? copy("password_min_hint", "At least 8 characters") : undefined}
        />
        <PasswordInputField
          label={copy("confirm_password", "Confirm password")}
          value={newPasswordConfirm}
          onChange={(event) => setNewPasswordConfirm(event.target.value)}
          autoComplete="new-password"
          invalid={passwordsMismatch}
          valid={Boolean(passwordsMatch)}
          errorText={copy("passwords_do_not_match", "Passwords do not match")}
          helperText={passwordsMatch ? copy("passwords_match", "Passwords match") : undefined}
        />
        <Button variant="approve" type="submit" disabled={busy === "password" || passwordTooShort || passwordsMismatch || !currentPassword || !newPassword || !newPasswordConfirm}>
          {busy === "password" ? <Spinner size="small" className="button-spinner-v2" /> : <ShieldCheck size={18} />}
          {copy("save_password", "Save password")}
        </Button>
      </form>
      <div className="settings-side-stack-v2">
        <section className="v2-panel settings-card-v2">
          <span className="eyebrow">{copy("telegram", "Telegram")}</span>
          <h2>{user.telegram_account?.name || copy("not_linked", "Not linked")}</h2>
          <p>{user.telegram_account?.id ? `Telegram ID ${user.telegram_account.id}` : copy("telegram_link_hint", "Link Telegram to share progress between bot and web app.")}</p>
          {!telegramLinked ? (
            <label className="auth-privacy-v2 telegram-privacy-v2">
              <input
                type="checkbox"
                checked={telegramPrivacyAccepted}
                onChange={(event) => {
                  setTelegramPrivacyAccepted(event.target.checked);
                  if (event.target.checked) setTelegramPrivacyError("");
                }}
              />
              <span>
                {copy("auth_privacy_consent", "I agree to personal data processing under the policy.")}{" "}
                <a href={privacyURL} target="_blank" rel="noreferrer">{copy("auth_privacy_policy", "Privacy policy")}</a>
                {" · "}
                <a href={consentURL} target="_blank" rel="noreferrer">{copy("auth_personal_data_consent", "Personal data consent")}</a>
                {" · "}
                <a href={agreementURL} target="_blank" rel="noreferrer">{copy("auth_user_agreement", "User agreement")}</a>
              </span>
            </label>
          ) : null}
          {telegramPrivacyError ? <p className="settings-error-v2">{telegramPrivacyError}</p> : null}
          <div className="telegram-actions-v2">
            {session.telegram_login_bot ? <a className="telegram-link-v2" href={`https://t.me/${session.telegram_login_bot}`} target="_blank" rel="noreferrer"><TelegramIcon />{copy("open_telegram", "Open Telegram")}</a> : null}
            {!telegramLinked ? (
              <Button variant="outline" onClick={() => void openTelegramTwoFactor()} disabled={busy === "telegram-code"}>
                {busy === "telegram-code" ? <Spinner size="small" className="button-spinner-v2" /> : <Send size={16} />}
                {copy("send_code", "Send code")}
              </Button>
            ) : null}
          </div>
          <OTPDialog
            open={telegramOtpOpen}
            title={copy("telegram_code_title", "Telegram two-factor code")}
            description={copy("telegram_code_description", "We opened Telegram with the same code flow as Version 1. Enter the 6-digit code from the bot to link this web account.")}
            telegramUrl={telegramRequest?.botURL}
            expiresAt={telegramCodeExpiresAt}
            verifying={busy === "telegram-code-verify"}
            resending={busy === "telegram-code"}
            onVerify={verifyTelegramTwoFactor}
            onResend={resendTelegramTwoFactor}
            onClose={() => {
              setTelegramOtpOpen(false);
              setTelegramRequest(null);
            }}
          />
        </section>
        <section className="v2-panel settings-card-v2 settings-social-card-v2">
          <span className="eyebrow">{copy("social_channels", "Social channels")}</span>
          <h2>{copy("poliglot_social_title", "Follow NERIVA")}</h2>
          <p>{copy("poliglot_social_body", "Short lessons, updates, and product tips.")}</p>
          <div className="settings-social-links-v2" aria-label={copy("poliglot_social_title", "Follow NERIVA")}>
            {poliglotSocialLinks.map(({ name, href, label, Icon }) => (
              <a key={name} className={`settings-social-link-v2 settings-social-link-v2--${name.toLowerCase()}`} href={href} aria-label={label} target="_blank" rel="noreferrer">
                <Icon />
                <span>{name}</span>
              </a>
            ))}
          </div>
        </section>
        <form className="v2-panel settings-card-v2" onSubmit={activateKey}>
          <span className="eyebrow">{copy("activation", "Activation")}</span>
          <h2>{copy("premium_key", "Premium key")}</h2>
          <input value={activationKey} onChange={(event) => setActivationKey(event.target.value)} placeholder="XXXX-XXXX-XXXX-XXXX" />
          <Button className="activation-button-v2" variant="shimmer" type="submit" disabled={busy === "activation"}><Crown size={18} />{copy("activate", "Activate")}</Button>
        </form>
      </div>
    </div>
  );
}

function MetricsView({ activeView, user, copy, selectedAwardLevel, setSelectedAwardLevel }: ViewRendererProps) {
  if (activeView === "awards") {
    const currentLevel = awardLevel(user);
    const selected = selectedAwardLevel ? awardInfo(selectedAwardLevel, copy, user.interface_language) : null;
    return (
      <section className="awards-stage-v2">
        <div className="awards-summary-v2">
          <img src={awardFile(currentLevel)} alt="" />
          <div>
            <span className="eyebrow">{copy("awards", "Awards")}</span>
            <h2>{currentLevel}/20</h2>
            <p>{awardInfo(currentLevel, copy, user.interface_language).rank}</p>
          </div>
          <Meter label={copy("awards_unlocked", "Unlocked")} value={currentLevel} max={20} />
        </div>
        <div className="awards-grid-v2">
          {Array.from({ length: 20 }, (_, index) => index + 1).map((level) => {
            const unlocked = currentLevel >= level;
            const info = awardInfo(level, copy, user.interface_language);
            return (
              <button
                className={cn("award-tile-v2", unlocked && "is-unlocked", selectedAwardLevel === level && "is-selected")}
                key={level}
                type="button"
                disabled={!unlocked}
                onClick={() => setSelectedAwardLevel(level)}
                aria-label={unlocked ? info.rank : copy("award_secret", "Secret trophy")}
              >
                <img src={awardFile(level)} alt="" />
                <span>{copy("level_label", "Level")} {level}</span>
                <strong>{unlocked ? info.rank : copy("awards_locked", "Locked")}</strong>
              </button>
            );
          })}
        </div>
        {selected ? <AwardModal selected={selected} copy={copy} onClose={() => setSelectedAwardLevel(null)} /> : null}
      </section>
    );
  }
  if (activeView === "referral") {
    return <ReferralView user={user} copy={copy} />;
  }
  return (
    <section className="metrics-grid-v2">
      <Metric label={copy("level_label", "Level")} value={user.level || "A1"} />
      <Metric label={copy("xp_level", "XP level")} value={user.xp_level || 1} />
      <Metric label={copy("lessons_total", "Lessons total")} value={user.lesson_count || 0} />
      <Metric label={copy("practices_total", "Practice total")} value={user.practice_count || 0} />
      <section className="v2-panel metric-wide-v2">
        <Meter label={copy("lessons_today", "Lessons today")} value={user.lessons_today} max={user.lesson_limit} />
        <Meter label={copy("practice_today", "Practice today")} value={user.practice_today} max={user.practice_limit} />
        <Meter label={copy("voices_today", "Voice messages today")} value={user.voice_today} max={user.voice_limit} />
      </section>
    </section>
  );
}

function ReferralView({ user, copy }: { user: UserProfile; copy: (key: string, fallback: string) => string }) {
  const [page, setPage] = useState(0);
  const link = user.referral_code ? `${location.origin}/app?ref=${user.referral_code}` : "";
  const shareTemplate = copy("referral_share", "Join NERIVA with my invite and get a free week of Premium: {link}");
  const shareText = link ? shareTemplate.replace("{link}", link).replace("%s", link) : "";
  const telegramShare = link ? `https://t.me/share/url?url=${encodeURIComponent(link)}&text=${encodeURIComponent(shareText)}` : "";
  const invitees = Array.isArray(user.referral_invitees) ? user.referral_invitees : [];
  const pageSize = 10;
  const totalPages = Math.max(1, Math.ceil(invitees.length / pageSize));
  const safePage = Math.min(page, totalPages - 1);
  const visibleInvitees = invitees.slice(safePage * pageSize, safePage * pageSize + pageSize);
  return (
    <div className="referral-view-v2">
      <section className="v2-panel data-display">
        <span className="eyebrow">{copy("referral_code", "Referral code")}</span>
        <h2>{user.referral_code || copy("no_code_yet", "No code yet")}</h2>
        <p>{link || copy("referral_unavailable", "Referral links are available after account setup.")}</p>
        <div className="referral-terms-v2">
          <div><strong>7</strong><span>{copy("referral_invitee_terms", "free days of Premium for every invited learner")}</span></div>
          <div><strong>7</strong><span>{copy("referral_inviter_terms", "days Premium when the invited learner reaches XP level 3")}</span></div>
          <div><strong>20%</strong><span>{copy("referral_direct_terms", "direct purchase reward")}</span></div>
          <div><strong>5%</strong><span>{copy("referral_indirect_terms", "second-line purchase reward")}</span></div>
          <div><strong>{user.referral_withdraw_min || "1000 RUB"}</strong><span>{copy("referral_withdraw_terms", "minimum withdrawal")}</span></div>
        </div>
        <div className="metric-grid-v2">
          <Metric label={copy("invites", "Invites")} value={user.referral_count || 0} />
          <Metric label={copy("balance", "Balance")} value={user.referral_balance || "0 RUB"} />
        </div>
        <div className="referral-actions-v2">
          <Button variant="outline" disabled={!link} onClick={() => link && navigator.clipboard?.writeText(shareText)}>{copy("copy_invite", "Copy invite")}</Button>
          {telegramShare ? <a className="telegram-link-v2" href={telegramShare} target="_blank" rel="noreferrer"><TelegramIcon />{copy("invite_telegram", "Invite in Telegram")}</a> : null}
          {navigator.share ? <Button variant="outline" disabled={!link} onClick={() => link && navigator.share({ title: "NERIVA", text: shareText, url: link })}>{copy("share", "Share")}</Button> : null}
        </div>
      </section>
      <section className="v2-panel referral-invitees-v2">
        <div className="panel-head">
          <div>
            <span className="eyebrow">{copy("referral_invitees", "Приглашённые")}</span>
            <h2>{copy("referral_invitees_title", "Кто перешёл по ссылке")}</h2>
          </div>
          <span>{invitees.length ? `${safePage + 1}/${totalPages}` : "0/0"}</span>
        </div>
        <div className="referral-invitees-v2__head">
          <span>{copy("learner", "Ученик")}</span>
          <span>{copy("level_label", "Уровень")}</span>
          <span>{copy("level_3_status", "3 уровень")}</span>
          <span>{copy("earned", "Заработано")}</span>
        </div>
        <div className="referral-invitees-v2__rows">
          {visibleInvitees.map((invitee) => (
            <article key={String(invitee.id || invitee.name || invitee.joined_at)}>
              <strong>{cleanAppText(invitee.name) || copy("learner", "Ученик")}<small>{invitee.joined_at ? prettyDate(invitee.joined_at) : ""}</small></strong>
              <span>{invitee.level || "A1"} / LVL {invitee.xp_level || 1}</span>
              <em className={invitee.reached_level_3 ? "is-reached" : ""}>{invitee.reached_level_3 ? copy("reached", "Достиг") : copy("not_yet", "Пока нет")}</em>
              <span>{invitee.earned || invitee.earned_usdt || "0"}</span>
            </article>
          ))}
          {!visibleInvitees.length ? <p className="empty-copy">{copy("no_referral_invitees", "Пока никто не перешёл по вашей ссылке.")}</p> : null}
        </div>
        {invitees.length > pageSize ? (
          <div className="pagination-v2 referral-pagination-v2">
            <Button variant="outline" size="sm" disabled={safePage <= 0} onClick={() => setPage((current) => Math.max(0, current - 1))}>{copy("back", "Назад")}</Button>
            <Button variant="outline" size="sm" disabled={safePage + 1 >= totalPages} onClick={() => setPage((current) => Math.min(totalPages - 1, current + 1))}>{copy("next", "Вперёд")}</Button>
          </div>
        ) : null}
      </section>
    </div>
  );
}

function AwardModal({ selected, copy, onClose }: { selected: AwardInfo; copy: (key: string, fallback: string) => string; onClose: () => void }) {
  if (typeof document === "undefined") return null;
  return createPortal(
    <div className="modal-backdrop-v2 modal-backdrop-v2--award" role="dialog" aria-modal="true">
      <section className="award-modal-v2">
        <img src={awardFile(selected.level)} alt="" />
        <div>
          <span className="eyebrow">{copy("level_label", "Level")} {selected.level}</span>
          <h3>{selected.rank}</h3>
          <p>{selected.story}</p>
        </div>
        <Button variant="outline" size="sm" onClick={onClose}>x</Button>
      </section>
    </div>,
    document.body,
  );
}

function ToolsView({ draft, setDraft, busy, copy, session, toolMode, setToolMode, toolSourceLanguage, setToolSourceLanguage, toolTargetLanguage, setToolTargetLanguage, toolVoiceFile, setToolVoiceFile, toolImageFile, setToolImageFile, submitTool, messages, savePhrase, isPhraseSaved }: ViewRendererProps) {
  const languageOptions = session.learning_languages?.length ? session.learning_languages : session.interface_languages || [];
  const fallbackTargetLanguage = languageOptions.find((language) => language.code !== toolTargetLanguage)?.code || toolTargetLanguage || "en";
  const toolMessages = messages.filter((message) => message.meta === `tools:${toolMode}`);
  const toolPhraseCandidates = useMemo(() => {
    const candidates = [...phraseCandidatesFromMessages(toolMessages)];
    const draftPhrase = cleanAppText(draft).replace(/\s+/g, " ").trim();
    if (isUsefulPhraseCandidate(draftPhrase)) candidates.unshift({ phrase: draftPhrase, source: "tool" as PhrasebookSource, details: { note: copy("tools", "Tools") } });
    const seen = new Set<string>();
    return candidates.filter((item) => {
      const key = phrasebookKey(item.phrase);
      if (!key || seen.has(key)) return false;
      seen.add(key);
      return true;
    }).slice(0, 4);
  }, [copy, draft, toolMessages]);
  const [mobilePickerOpen, setMobilePickerOpen] = useState(true);
  const swapToolLanguages = () => {
    if (toolSourceLanguage === "auto") {
      setToolSourceLanguage(toolTargetLanguage || fallbackTargetLanguage);
      setToolTargetLanguage(fallbackTargetLanguage);
      return;
    }
    setToolSourceLanguage(toolTargetLanguage || "auto");
    setToolTargetLanguage(toolSourceLanguage || fallbackTargetLanguage);
  };
  const chooseTool = (mode: ToolMode) => {
    setToolMode(mode);
    setMobilePickerOpen(false);
  };
  const toolText = cleanAppText(draft).trim();
  const toolsBusy = Boolean(busy?.startsWith("tools-"));
  const toolSubmitDisabled = toolsBusy
    || (toolMode === "translator" && !toolText)
    || (toolMode === "voice" && !toolVoiceFile)
    || (toolMode === "image" && !toolImageFile);
  const handleToolsPaste = useCallback((event: ClipboardEvent<HTMLDivElement>) => {
    if (toolMode !== "image" || event.defaultPrevented) return;
    const pastedImage = firstClipboardImage(event);
    if (!pastedImage) return;
    event.preventDefault();
    setToolImageFile(nameClipboardImage(pastedImage));
  }, [setToolImageFile, toolMode]);
  useEffect(() => {
    if (toolMode !== "image") return;
    const handleDocumentPaste = (event: globalThis.ClipboardEvent) => {
      if (event.defaultPrevented) return;
      const pastedImage = firstClipboardImageData(event.clipboardData);
      if (!pastedImage) return;
      event.preventDefault();
      setToolImageFile(nameClipboardImage(pastedImage));
    };
    document.addEventListener("paste", handleDocumentPaste);
    return () => document.removeEventListener("paste", handleDocumentPaste);
  }, [setToolImageFile, toolMode]);
  return (
    <div className={cn("tools-layout-v2", mobilePickerOpen && "is-picker-open")} onPaste={handleToolsPaste}>
      <div className="tool-switch-v2">
        <ToolButton active={!mobilePickerOpen && toolMode === "translator"} icon={Languages} title={copy("text_translator", "Text translator")} onClick={() => chooseTool("translator")} />
        <ToolButton active={!mobilePickerOpen && toolMode === "voice"} icon={FileAudio} title={copy("voice_to_text", "Voice to text")} onClick={() => chooseTool("voice")} />
        <ToolButton active={!mobilePickerOpen && toolMode === "image"} icon={ImageIcon} title={copy("photo_translation", "Photo translation")} onClick={() => chooseTool("image")} />
        <a className="tool-button-v2 tools-ai-router-v2" href={aiRouterTelegramURL} target="_blank" rel="noreferrer" aria-label={copy("ai_router", "AI Router")}>
          <Bot size={18} />
          {copy("ai_router", "AI Router")}
        </a>
      </div>
      <div className="tools-work-v2">
      <Button className="tools-change-v2" variant="outline" size="sm" type="button" onClick={() => setMobilePickerOpen(true)}>
        <Wand2 size={16} />
        {copy("choose_tool", "Выбрать инструмент")}
      </Button>
      <ChatPanel messages={toolMessages} copy={copy} targetLanguage={toolTargetLanguage} />
      <section className="v2-panel composer-panel-v2">
        <span className="eyebrow">{copy("tools", "Tools")}</span>
        <h2>{toolMode === "translator" ? copy("quick_translator", "Quick translator") : toolMode === "voice" ? copy("voice_to_text", "Voice to text") : copy("photo_translation", "Photo translation")}</h2>
        <div className="tool-input-shell-v2">
          <textarea
            value={draft}
            onChange={(event) => setDraft(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === "Enter" && !event.shiftKey) {
                event.preventDefault();
                if (toolSubmitDisabled) return;
                void submitTool();
              }
            }}
            placeholder={toolMode === "translator" ? copy("paste_text_translate", "Paste text to translate") : copy("optional_context", "Optional context for the tool")}
          />
          <Button className="tools-submit-v2" onClick={() => void submitTool()} disabled={toolSubmitDisabled}>
            {toolsBusy ? <Spinner size="small" className="button-spinner-v2" /> : <Send size={18} />}
            <span>{copy("send", "Send")}</span>
          </Button>
        </div>
        {toolMode === "translator" ? (
          <div className="language-row-v2 language-row-v2--with-swap">
            <AnimatedSelect
              label={copy("source_language", "Source language")}
              value={toolSourceLanguage}
              options={[{ code: "auto", native_name: copy("auto_detect", "Auto detect") }, ...languageOptions]}
              onChange={setToolSourceLanguage}
            />
            <button className="language-swap-v2" type="button" onClick={swapToolLanguages} aria-label={copy("swap_languages", "Swap languages")}>
              <Repeat2 size={18} />
            </button>
            <AnimatedSelect
              label={copy("target_language", "Target language")}
              value={toolTargetLanguage}
              options={languageOptions}
              onChange={setToolTargetLanguage}
            />
          </div>
        ) : (
          <FileControls voiceFile={toolVoiceFile} imageFile={toolImageFile} setVoiceFile={setToolVoiceFile} setImageFile={setToolImageFile} allowImage={toolMode === "image"} allowVoice={toolMode === "voice"} copy={copy} />
        )}
        {toolPhraseCandidates.length ? <PhraseQuickSave candidates={toolPhraseCandidates} savePhrase={savePhrase} isPhraseSaved={isPhraseSaved} copy={copy} /> : null}
      </section>
      </div>
    </div>
  );
}

function ListPanel({ title, action, onAction, busy, items, copy, extraAction, showAction = true }: { title: string; action: string; onAction: () => Promise<void>; busy: boolean; items: Array<{ key: string; title?: string; meta?: string; body?: string; audio?: Array<{ label: string; text: string }>; action?: () => Promise<void> }>; copy: (key: string, fallback: string) => string; extraAction?: ReactNode; showAction?: boolean }) {
  return (
    <section className="v2-panel list-panel-v2">
      <div className="panel-head">
        <div>
          <span className="eyebrow">{copy("mistakes", "Mistakes")}</span>
          <h2>{title}</h2>
        </div>
        <div className="panel-actions-v2">
          {showAction ? <Button variant="outline" size="sm" onClick={onAction} disabled={busy}>{busy ? <Spinner size="small" className="button-spinner-v2" /> : null}{action}</Button> : null}
          {extraAction}
        </div>
      </div>
      <div className="list-items-v2">
        {items.length ? items.map((item) => (
          <article className="list-item-v2" key={item.key}>
            <div>
              <strong>{item.title}</strong>
              {item.meta ? <span>{item.meta}</span> : null}
              {item.body ? <p>{item.body}</p> : null}
              {item.audio?.length ? <AudioActionRow clips={item.audio} /> : null}
            </div>
            {item.action ? <Button variant="outline" size="sm" onClick={() => void item.action?.()}>{copy("practice", "Practice")}</Button> : null}
          </article>
        )) : <p className="empty-copy">{copy("empty_panel", "No items yet. Start the related practice mode to fill this panel.")}</p>}
      </div>
    </section>
  );
}

function ActionCard({ icon: Icon, title, meta, onClick, busy }: { icon: LucideIcon; title: string; meta: string; onClick: () => Promise<void> | void; busy?: boolean }) {
  return (
    <button className="action-card-v2" type="button" onClick={() => void onClick()} disabled={busy}>
      <span>{busy ? <Spinner size="small" className="button-spinner-v2" /> : <Icon size={21} />}</span>
      <strong>{title}</strong>
      <small>{meta}</small>
    </button>
  );
}

function ToolButton({ active, icon: Icon, title, onClick }: { active: boolean; icon: LucideIcon; title: string; onClick: () => void }) {
  return (
    <button className={cn("tool-button-v2", active && "is-active")} type="button" onClick={onClick} aria-pressed={active}>
      <Icon size={18} />
      {title}
    </button>
  );
}

function FileControls({ voiceFile, imageFile, setVoiceFile, setImageFile, allowImage, allowVoice = true, copy }: { voiceFile: File | null; imageFile: File | null; setVoiceFile: (file: File | null) => void; setImageFile: (file: File | null) => void; allowImage: boolean; allowVoice?: boolean; copy: (key: string, fallback: string) => string }) {
  const recorderRef = useRef<MediaRecorder | null>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const chunksRef = useRef<BlobPart[]>([]);
  const [recording, setRecording] = useState(false);
  const [recordError, setRecordError] = useState("");

  const stopRecording = useCallback(() => {
    recorderRef.current?.stop();
  }, []);

  const startRecording = useCallback(async () => {
    if (!navigator.mediaDevices?.getUserMedia) {
      setRecordError(copy("mic_unavailable", "Microphone is unavailable"));
      return;
    }
    setRecordError("");
    let stream: MediaStream;
    try {
      stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    } catch {
      setRecordError(copy("mic_unavailable", "Microphone is unavailable"));
      return;
    }
    streamRef.current = stream;
    chunksRef.current = [];
    const recorder = new MediaRecorder(stream);
    recorderRef.current = recorder;
    recorder.ondataavailable = (event) => {
      if (event.data.size > 0) chunksRef.current.push(event.data);
    };
    recorder.onstop = () => {
      const blob = new Blob(chunksRef.current, { type: recorder.mimeType || "audio/webm" });
      setVoiceFile(new File([blob], `voice-${Date.now()}.webm`, { type: blob.type || "audio/webm" }));
      streamRef.current?.getTracks().forEach((track) => track.stop());
      streamRef.current = null;
      recorderRef.current = null;
      setRecording(false);
    };
    recorder.start();
    setRecording(true);
  }, [copy, setVoiceFile]);

  return (
    <div className="file-controls-v2">
      {allowVoice ? (
        <>
          <VoiceInput
            className={cn("record-button-v2", recording && "is-recording")}
            isListening={recording}
            onToggle={() => recording ? stopRecording() : void startRecording()}
          />
          <AudioUploadCard
            className="voice-upload-control-v2"
            title={copy("voice_file", "Voice file")}
            description={voiceFile ? voiceFile.name : copy("drop_or_upload_audio", "Drop or upload audio")}
            file={voiceFile}
            onFileChange={setVoiceFile}
            compact
          />
          {recordError ? <small className="file-error-v2">{recordError}</small> : null}
        </>
      ) : null}
      {allowImage ? <ImageUploadControl imageFile={imageFile} setImageFile={setImageFile} copy={copy} /> : null}
    </div>
  );
}

function ImageUploadControl({ imageFile, setImageFile, copy }: { imageFile: File | null; setImageFile: (file: File | null) => void; copy: (key: string, fallback: string) => string }) {
  const [{ files, isDragging }, { addFiles, removeFile, clearFiles, openFileDialog, getInputProps, handleDragEnter, handleDragLeave, handleDragOver, handleDrop }] = useFileUpload({
    accept: "image/*",
    onFilesAdded: (added) => {
      const first = added[0]?.file;
      if (first instanceof File) setImageFile(first);
    },
    onFilesChange: (nextFiles) => {
      const first = nextFiles[0]?.file;
      setImageFile(first instanceof File ? first : null);
    },
  });
  const fallbackPreviewUrl = useMemo(() => (imageFile ? URL.createObjectURL(imageFile) : ""), [imageFile]);
  const previewUrl = files[0]?.preview || fallbackPreviewUrl;
  const fileName = files[0]?.file.name || imageFile?.name || "";
  const handlePasteImage = useCallback((event: ClipboardEvent<HTMLDivElement>) => {
    const pastedImage = firstClipboardImage(event);
    if (!pastedImage) return;
    event.preventDefault();
    addFiles([nameClipboardImage(pastedImage)]);
  }, [addFiles]);

  useEffect(() => {
    if (!imageFile && files.length) clearFiles();
  }, [clearFiles, files.length, imageFile]);

  useEffect(() => {
    return () => {
      if (fallbackPreviewUrl) URL.revokeObjectURL(fallbackPreviewUrl);
    };
  }, [fallbackPreviewUrl]);

  return (
    <div
      className={cn("image-upload-control-v2", isDragging && "is-dragging")}
      onDragEnter={handleDragEnter}
      onDragLeave={handleDragLeave}
      onDragOver={handleDragOver}
      onDrop={handleDrop}
      onPaste={handlePasteImage}
      tabIndex={0}
    >
      <button
        type="button"
        className="attachment-button-v2 image-upload-control-v2__preview"
        onClick={openFileDialog}
        aria-label={fileName ? copy("change_image", "Change image") : copy("upload_image", "Upload image")}
      >
        {previewUrl ? <img src={previewUrl} alt="" /> : <FileImage size={18} />}
      </button>
      <input {...getInputProps()} className="sr-only" aria-label={copy("upload_image", "Upload image")} tabIndex={-1} />
      {fileName ? (
        <div className="file-chip-v2">
          <FileImage size={15} />
          <span>{fileName}</span>
          <button type="button" onClick={() => { removeFile(files[0]?.id || ""); setImageFile(null); }} aria-label={copy("remove", "Remove")}>
            <X size={13} />
          </button>
        </div>
      ) : null}
    </div>
  );
}

function isClipboardImageFile(file: File) {
  return file.type.startsWith("image/") || /\.(png|jpe?g|webp|gif|bmp|svg|heic|heif)$/i.test(file.name || "");
}

function firstClipboardImage(event: ClipboardEvent<HTMLElement>) {
  return firstClipboardImageData(event.clipboardData);
}

function firstClipboardImageData(clipboardData: DataTransfer | null) {
  const clipboardFiles = Array.from(clipboardData?.files || []).filter(isClipboardImageFile);
  const itemFiles = Array.from(clipboardData?.items || [])
    .filter((item) => item.kind === "file" && item.type.startsWith("image/"))
    .map((item) => item.getAsFile())
    .filter((file): file is File => file instanceof File && isClipboardImageFile(file));
  return [...clipboardFiles, ...itemFiles][0] || null;
}

function nameClipboardImage(file: File) {
  if (file.name) return file;
  return new File([file], "clipboard-image.png", { type: file.type || "image/png" });
}

function Meter({ label, value, max }: { label: string; value: unknown; max: unknown }) {
  const pct = percent(value, max);
  return (
    <div className="meter-v2">
      <div><span>{label}</span><strong>{formatLimit(value, max)}</strong></div>
      <i><b style={{ width: `${pct}%` }} /></i>
    </div>
  );
}

function Metric({ label, value }: { label: string; value: unknown }) {
  return (
    <div className="metric-v2">
      <span>{label}</span>
      <strong>{typeof value === "number" ? compactNumber(value) : asText(value, "0")}</strong>
    </div>
  );
}

function buildShadowingForm(text: string, voice: File | null, target: string) {
  const form = new FormData();
  if (target.trim()) form.append("target", target.trim());
  if (text.trim()) form.append("text", text.trim());
  if (voice) form.append("voice", voice, voice.name || "repeat.webm");
  return form;
}

function formatLimit(value: unknown, max: unknown) {
  const current = Number(value || 0);
  const total = Number(max || 0);
  if (!total) return compactNumber(current);
  return `${compactNumber(current)} / ${compactNumber(total)}`;
}

function fallbackPlans(copy: (key: string, fallback: string) => string): PremiumPlan[] {
  return [
    {
      product: "free",
      title: copy("free_plan_title", "Free"),
      tier: "free",
      label: copy("free_plan_label", "Для начала"),
      body: copy("free_plan_body", "Basic text learning, word training, notes, and progress overview. AI Tutor, listening practice, pronunciation scoring, and voice review open in Premium."),
      features: [
        copy("free_feature_daily", "ежедневная привычка и стартовые уроки"),
        copy("free_feature_words", "базовая тренировка слов"),
        copy("free_feature_phrasebook", "notes, Phrasebook, and progress overview"),
        copy("free_feature_phrase_audio", "Phrase audio is available"),
      ],
      locked_features: [
        copy("free_locked_ai_tutor", "AI-guided tutor lessons are locked until Premium"),
        copy("free_locked_listening", "Listening practice and pronunciation scoring need Premium"),
        copy("free_locked_voice_photo", "Voice review and image tools are locked until Premium"),
      ],
      note: copy("free_plan_note", "Подходит для знакомства с продуктом без оплаты."),
      rub_price: "0",
    },
    {
      product: "premium_month",
      title: copy("premium_month_title", "Premium"),
      tier: "premium",
      label: copy("premium_plan_label", "Регулярная учеба"),
      body: copy("premium_month_body", "The main mode for daily practice: AI Tutor, listening practice, pronunciation scoring, AI-guided tutor lessons, voice review, image tools, and expanded daily limits."),
      features: [
        copy("premium_feature_voice_text", "голос в текст и перевод услышанного"),
        copy("premium_feature_image_text", "перевод текста с картинки"),
        copy("premium_feature_context", "практика по контексту голоса или фото"),
        copy("premium_feature_errors", "mistake dictionary, notes, XP, and streak"),
      ],
      note: copy("premium_plan_note", "Лучший выбор для стабильного ежедневного обучения."),
      rub_price: "300",
    },
    {
      product: "platinum_month",
      title: copy("platinum_month_title", "Platinum"),
      tier: "platinum",
      label: copy("platinum_plan_label", "Интенсив"),
      body: copy("platinum_month_body", "AI Tutor with maximum daily limits, deeper role-play practice, intensive review practice, and maximum voice and pronunciation practice."),
      features: [
        copy("platinum_feature_limits", "максимальные дневные лимиты"),
        copy("platinum_feature_voice", "more voice or image-context practice"),
        copy("platinum_feature_review", "intensive weak-spot review"),
        copy("platinum_feature_priority", "best mode for dense daily study"),
      ],
      note: copy("platinum_plan_note", "Для поездки, работы, экзамена или очень плотного темпа."),
      rub_price: "590",
    },
  ];
}

const MENU_ASSET_VERSION = "flux2-20260526-pronunciation-phrases-offline-roleplay";

function assetSlug(view: ViewId) {
  if (view === "word-game") return "word-game";
  if (view === "dashboard") return "progress";
  return view;
}

function menuAssetUrl(view: ViewId, theme: Theme) {
  if (view === "tutor") return `/app/assets/brand-logo-mini.png?v=${MENU_ASSET_VERSION}`;
  return `/app/assets/icon-${assetSlug(view)}-${theme}.png?v=${MENU_ASSET_VERSION}`;
}
