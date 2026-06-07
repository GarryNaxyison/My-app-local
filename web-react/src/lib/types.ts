import type { ComponentType } from "react";
import type { LucideProps } from "lucide-react";

export type ViewId =
  | "home"
  | "tutor"
  | "lesson"
  | "practice"
  | "roleplay"
  | "shadowing"
  | "pronunciation"
  | "words"
  | "word-game"
  | "spelling"
  | "vocabulary"
  | "phrasebook"
  | "offline"
  | "level"
  | "progress"
  | "awards"
  | "leaderboard"
  | "limits"
  | "mistakes"
  | "tools"
  | "dashboard"
  | "referral"
  | "premium"
  | "settings";

export type LanguageOption = {
  code: string;
  name?: string;
  native_name?: string;
  nativeName?: string;
  interface_name?: string;
};

export type UserProfile = {
  interface_language?: string;
  learning_language?: string;
  telegram_linked?: boolean;
  telegram_account?: { id?: number; name?: string; label?: string } | null;
  created_at?: string;
  level?: string;
  learning_focus?: string;
  plan?: string;
  premium?: boolean;
  premium_until?: string;
  referral_code?: string;
  referral_count?: number;
  referral_balance?: string;
  referral_balance_usdt?: string;
  referral_withdraw_min?: string;
  referral_withdraw_min_usdt?: string;
  referral_invitees?: ReferralInvitee[];
  phrasebook?: PhrasebookItem[];
  usdt_rub_rate?: string;
  invited_by?: string;
  xp?: number;
  xp_level?: number;
  xp_title?: string;
  xp_current?: number;
  xp_needed?: number;
  lessons_today?: number;
  lesson_limit?: number;
  practice_today?: number;
  practice_limit?: number;
  voice_today?: number;
  voice_limit?: number;
  daily_bonus_claims?: string[];
  daily_bonus_last_claimed_at?: string;
  habit_log?: Record<string, {
    login?: boolean;
    complete?: boolean;
    claimed?: boolean;
    claimedAt?: string;
    claimed_at?: string;
    lessons?: number;
    practice?: number;
    voice?: number;
  }>;
  navigation_layout?: {
    function_ribbon?: ViewId[];
    mobile_pinned?: ViewId[];
    mobile_more?: ViewId[];
    mobile_rail?: ViewId[];
  };
  lesson_count?: number;
  practice_count?: number;
  word_game_count?: number;
  learned_words?: number;
  mistakes?: number;
};

export type ReferralInvitee = {
  id?: number | string;
  name?: string;
  level?: string;
  xp_level?: number;
  xp?: number;
  reached_level_3?: boolean;
  level_rewarded?: boolean;
  earned?: string;
  earned_usdt?: string;
  joined_at?: string;
};

export type PhrasebookItem = {
  id: string;
  phrase: string;
  translation?: string;
  note?: string;
  source?: "lesson" | "practice" | "roleplay" | "manual";
  language?: string;
  createdAt?: string;
  created_at?: string;
};

export type SessionData = {
  authenticated: boolean;
  account?: {
    login?: string;
    password_set?: boolean;
    profile_ready?: boolean;
  };
  user?: UserProfile;
  copy?: Record<string, string>;
  tool_copy?: Record<string, string>;
  system_copy?: Record<string, string>;
  premium_copy?: Record<string, string>;
  interface_languages?: LanguageOption[];
  learning_languages?: LanguageOption[];
  premium_plans?: PremiumPlan[];
  yookassa_enabled?: boolean;
  crypto_enabled?: boolean;
  telegram_login_bot?: string;
  web_app_url?: string;
  support?: {
    telegram?: string;
    email?: string;
  };
  captcha?: {
    enabled?: boolean;
    site_key?: string;
    provider?: string;
  };
  activation?: {
    duration_days?: number;
    tier?: string;
    premium_until?: string;
  };
};

export type PremiumPlan = {
  product: string;
  tier?: string;
  title?: string;
  days_label?: string;
  rub_price?: string | number;
  stars_price?: string | number;
  usdt_price?: string;
  crypto_enabled?: boolean;
  crypto_methods?: CryptoMethod[];
};

export type CryptoMethod = {
  id: string;
  label?: string;
  currency?: string;
  network?: string;
};

export type NavItem = {
  id: ViewId;
  label: string;
  shortLabel?: string;
  group: "learn" | "words" | "growth" | "account";
  description: string;
  icon: ComponentType<LucideProps>;
  accent: string;
};

export type ChatMessage = {
  id: string;
  role: "system" | "user" | "assistant";
  title?: string;
  body: string;
  tone?: "default" | "success" | "warning" | "danger";
  meta?: string;
  details?: Record<string, unknown>;
};

export type ChoiceOption = {
  id: string;
  text: string;
};

export type WordChallenge = {
  empty?: boolean;
  message?: string;
  prompt?: string;
  context?: string;
  direction?: string;
  options?: ChoiceOption[];
  instruction?: string;
  correct_answer_id?: string;
};

export type SpellingChallenge = {
  empty?: boolean;
  message?: string;
  word_id?: string;
  prompt?: string;
  context?: string;
  direction?: string;
};

export type LevelQuestion = {
  index: number;
  total: number;
  question: string;
  options: string[];
  score?: number;
};

export type VocabularyItem = {
  id: string;
  word: string;
  translation?: string;
  context?: string;
  example?: string;
  review_correct_count?: number;
  spelling_correct_count?: number;
  training_count?: number;
};

export type VocabularyPage = {
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
};

export type MistakeItem = {
  index?: number;
  word?: string;
  correction?: string;
  context?: string;
  explanation?: string;
  created_at?: string;
  count?: number;
};

export type LeaderboardEntry = {
  name?: string;
  score?: number;
  xp?: number;
  words?: number;
  mistakes?: number;
  languages?: string[];
  level?: string;
  title?: string;
};

export type TranslatorResult = {
  source_text?: string;
  translation?: string;
  result?: string;
  source_language?: string;
  target_language?: string;
};
