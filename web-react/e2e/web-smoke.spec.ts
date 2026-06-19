import { expect, test, type Page } from "@playwright/test";
import { readFileSync } from "node:fs";
import { appCopy, appLocaleCodes } from "../src/lib/i18n";

const png = Buffer.from(
  "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/p9sAAAAASUVORK5CYII=",
  "base64",
);

const ru = (key: string, fallback = key) => appCopy("ru", key, fallback);
const mojibakePattern = /пїЅ|Гђ|Г‘|\u0420[\u201a\u0403\u201e\u2026\u2020\u2021\u20ac\u2030\u0409\u2039\u040a\u040b\u040f]|\u0421[\u201a\u0403\u201e\u2026\u2020\u2021\u20ac\u2030\u0409\u2039\u040a\u040b\u040f]/u;

const broadMojibakePattern = /РїС—Р…|Р“С’|Р“вЂ|РЎвЂ№|Р±В»|Р±С”|Р±С—/u;

const vocabularyItems = Array.from({ length: 12 }, (_, index) => ({
  id: `word-${index + 1}`,
  word: `travel-word-${index + 1}`,
  translation: index === 2 ? "самолётостроение" : `слово ${index + 1}`,
  example: `Example sentence ${index + 1}`,
}));

const referralInvitees = Array.from({ length: 11 }, (_, index) => ({
  id: 9000 + index,
  name: `ref-user-${index + 1}`,
  level: index % 3 === 0 ? "B1" : "A2",
  xp_level: index + 1,
  xp: 120 + index * 35,
  reached_level_3: index % 2 === 0,
  level_rewarded: index % 2 === 0,
  earned: index % 2 === 0 ? "150 RUB" : "0 RUB",
  earned_usdt: index % 2 === 0 ? 1.5 : 0,
  joined_at: "2026-05-26T10:00:00Z",
}));

const mistakeItems = Array.from({ length: 18 }, (_, index) => ({
  index,
  word: `wrong phrase ${index + 1}`,
  correction: `Correct phrase ${index + 1}.`,
  explanation: index % 2 === 0 ? "Use the verb and article." : "Word order needs correction.",
}));

const roleplayScenarioIds = [
  "restaurant",
  "work",
  "travel",
  "exam",
  "small-talk",
  "hotel",
  "shopping",
  "doctor",
  "job-interview",
  "bank",
];

const leaderboardItems = Array.from({ length: 12 }, (_, index) => ({
  name: index === 0 ? "Irina" : `Learner ${index + 1}`,
  xp: 520 - index * 23,
  words: 42 - index,
  mistakes: index % 3,
  level: index % 2 === 0 ? "B1" : "A2",
  languages: ["en", "es", "de"],
}));

const phrasebookSeed = [
  {
    id: "seed-phrase",
    phrase: "I have a reservation under the name Ivan Petrov.",
    translation: "У меня бронь на имя Иван Петров.",
    source: "lesson",
    language: "en",
    createdAt: new Date().toISOString(),
  },
];

const sessionPayload = {
  authenticated: true,
  account: { login: "demor22", password_set: true, profile_ready: true },
  user: {
    interface_language: "ru",
    learning_language: "en",
    telegram_linked: true,
    telegram_account: { id: 185156683, name: "demor22" },
    level: "A2",
    plan: "free",
    premium: false,
    xp: 635,
    xp_level: 9,
    xp_title: "Алхимик фраз",
    xp_current: 58,
    xp_needed: 100,
    lessons_today: 1,
    lesson_limit: 1,
    practice_today: 1,
    practice_limit: 1,
    voice_today: 0,
    voice_limit: 0,
    lesson_count: 12,
    practice_count: 9,
    word_game_count: 5,
    learned_words: 24,
    mistakes: 2,
    invited_count: 11,
    referral_balance_kopecks: 30000,
    referral_code: "DEMO22",
    referral_link: "https://t.me/poliglot_ai_bot?start=ref_DEMO22",
    referral_invitees: referralInvitees,
    phrasebook: phrasebookSeed,
  },
  interface_languages: appLocaleCodes.map((code) => ({ code, native_name: code.toUpperCase() })),
  learning_languages: [{ code: "en", native_name: "English" }],
  premium_plans: [
    { product: "premium_month", title: "Premium 30 дней", days_label: "6 AI audio actions per day", rub_price: "300", crypto_enabled: true },
    { product: "platinum_month", title: "Platinum 30 дней", days_label: "18 AI audio actions per day", rub_price: "590", crypto_enabled: true },
  ],
  yookassa_enabled: true,
  crypto_enabled: true,
  telegram_login_bot: "poliglot_ai_bot",
};

const legacyPhrasebookSeed = [
  {
    id: "seed-phrase",
    phrase: "I have a reservation under the name Ivan Petrov.",
    translation: "У меня бронь на имя Иван Петров.",
    source: "lesson",
    language: "en",
    createdAt: new Date().toISOString(),
  },
];

const paymentHistorySeed = [
  {
    id: "pending-ton",
    date: "2026-05-26T09:00:00Z",
    plan: "Premium 30 дней",
    period: "Месяц",
    amount: "1.5 TON",
    method: "TON",
    status: "pending",
  },
  {
    id: "created-card",
    date: "2026-05-26T09:05:00Z",
    plan: "Platinum 30 дней",
    period: "Месяц",
    amount: "590 RUB",
    method: "Bank card",
    status: "created",
  },
  {
    id: "paid-card",
    date: "2026-05-26T09:10:00Z",
    plan: "Premium 30 дней",
    period: "Месяц",
    amount: "300 RUB",
    method: "Bank card",
    status: "paid",
  },
];

let testSessionPayloadOverride: typeof sessionPayload | null = null;
let testPhrasebookItemsOverride: unknown[] | null = null;

async function mockApi(page: Page) {
  await page.addInitScript(
    ({ phrases, payments }) => {
      if ("serviceWorker" in navigator) {
        navigator.serviceWorker.getRegistrations().then((registrations) => {
          registrations.forEach((registration) => registration.unregister());
        });
      }
      if ("caches" in window) {
        caches.keys().then((keys) => {
          keys.forEach((key) => caches.delete(key));
        });
      }
      if (localStorage.getItem("poliglot-test-show-onboarding") === "1") {
        localStorage.removeItem("poliglot-onboarding-v2:demor22");
      } else {
        localStorage.setItem("poliglot-onboarding-v2:demor22", JSON.stringify({ completed_at: new Date().toISOString() }));
      }
      localStorage.setItem("poliglot-phrasebook-v2:demor22", JSON.stringify(phrases));
      localStorage.setItem("poliglot-payment-history-v2:demor22", JSON.stringify(payments));
      localStorage.removeItem("poliglot-offline-deck-v2");
      if (sessionStorage.getItem("poliglot-test-nav-reset-v2") !== "1") {
        localStorage.removeItem("poliglot-mobile-nav-v2:demor22");
        localStorage.removeItem("poliglot-mobile-nav-more-v2:demor22");
        localStorage.removeItem("poliglot-mobile-nav-rail-v2:demor22");
        sessionStorage.setItem("poliglot-test-nav-reset-v2", "1");
      }
    },
    { phrases: phrasebookSeed, payments: paymentHistorySeed },
  );

  const wordRounds = [
    {
      prompt: "яблоко",
      correct: "en:apple",
      options: [
        { id: "en:apple", text: "apple" },
        { id: "en:station", text: "station" },
        { id: "en:ticket", text: "ticket" },
        { id: "en:coffee", text: "coffee" },
      ],
    },
    {
      prompt: "поезд",
      correct: "en:train",
      options: [
        { id: "en:train", text: "train" },
        { id: "en:window", text: "window" },
        { id: "en:meeting", text: "meeting" },
        { id: "en:market", text: "market" },
      ],
    },
  ];
  let wordRoundIndex = 0;
  let tutorCompleteCount = 0;

  await page.route("**/app/assets/**", (route) => {
    const request = route.request();
    const assetPath = new URL(request.url()).pathname.toLowerCase();
    if (request.resourceType() === "image" || /\.(avif|gif|ico|jpe?g|png|svg|webp)$/.test(assetPath)) {
      return route.fulfill({ status: 200, contentType: "image/png", body: png });
    }
    return route.continue();
  });
  await page.route("**/api/session", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(testSessionPayloadOverride || sessionPayload) }));
  await page.route("**/api/settings", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(sessionPayload) }));
  await page.route("**/api/navigation-layout", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(sessionPayload) }));
  const aiTutorLesson = {
    title: "AI Репетитор",
    level: "A1",
    target_language: "en",
    interface_language: "ru",
    theme: "Утренняя рутина",
    lesson_goal: "Понять короткую историю и составить 2-3 предложения о своей рутине.",
    story: {
      text_target: "Mia wakes up at seven. She gets ready quickly. She eats breakfast. She leaves home at eight. She starts work at nine.",
      audio_text_target: "Mia wakes up at seven. She gets ready quickly. She eats breakfast. She leaves home at eight. She starts work at nine.",
    },
    words: [
      { id: "w1", target: "wake up", interface_translation: "просыпаться", example_sentence_target: "I wake up early on weekdays.", audio_text_target: "wake up", example_audio_text_target: "I wake up early on weekdays." },
      { id: "w2", target: "get ready", interface_translation: "собираться", example_sentence_target: "I get ready in ten minutes.", audio_text_target: "get ready", example_audio_text_target: "I get ready in ten minutes." },
      { id: "w3", target: "routine", interface_translation: "рутина", example_sentence_target: "My morning routine is simple.", audio_text_target: "routine", example_audio_text_target: "My morning routine is simple." },
      { id: "w4", target: "leave", interface_translation: "уходить", example_sentence_target: "I leave home at eight.", audio_text_target: "leave", example_audio_text_target: "I leave home at eight." },
      { id: "w5", target: "breakfast", interface_translation: "завтрак", example_sentence_target: "I eat breakfast at home.", audio_text_target: "breakfast", example_audio_text_target: "I eat breakfast at home." },
      { id: "w6", target: "work", interface_translation: "работа", example_sentence_target: "I start work at nine.", audio_text_target: "work", example_audio_text_target: "I start work at nine." },
    ],
    comprehension_questions: [
      { id: "q1", question_target: "What time does the person wake up?", expected_points: ["seven"] },
      { id: "q2", question_target: "What does the person eat?", expected_points: ["breakfast"] },
      { id: "q3", question_target: "What time does work start?", expected_points: ["nine"] },
    ],
    production_task: {
      instruction_interface: "Напишите 2-3 предложения о своей утренней рутине.",
      required_word_count: 3,
      sentence_count: "2-3",
      recommendations_interface: ["Повторите wake up, get ready и leave завтра."],
    },
  };
  const aiTutorStages = [
    "story_intro",
    "retell",
    "question_1",
    "question_2",
    "question_3",
    "word_learn_1",
    "word_learn_2",
    "word_learn_3",
    "word_learn_4",
    "word_learn_5",
    "word_learn_6",
    "word_recall_1",
    "word_recall_2",
    "word_recall_3",
    "word_recall_4",
    "word_recall_5",
    "word_recall_6",
    "production",
    "lesson_feedback",
    "review_schedule",
    "complete",
  ];
  let aiTutorStage = "story_intro";
  const aiTutorStep = (stage: string) => {
    const wordMatch = stage.match(/^word_(learn|recall)_(\d+)$/);
    const questionMatch = stage.match(/^question_(\d+)$/);
    const base: Record<string, unknown> = { stage, lesson: aiTutorLesson };
    if (stage === "story_intro") return { ...base, kind: "story", title: aiTutorLesson.title, instruction: aiTutorLesson.lesson_goal };
    if (stage === "retell") return { ...base, kind: "free_text", title: "Retell", instruction: "Расскажите историю своими словами." };
    if (questionMatch) {
      const question = aiTutorLesson.comprehension_questions[Number(questionMatch[1]) - 1];
      return { ...base, kind: "free_text", title: `Question ${questionMatch[1]}`, instruction: question.question_target, question };
    }
    if (wordMatch) {
      const index = Number(wordMatch[2]) - 1;
      const word = aiTutorLesson.words[index];
      if (wordMatch[1] === "learn") return { ...base, kind: "word_learn", title: word.target, instruction: word.interface_translation, word };
      const options = [
        aiTutorLesson.words[(index + 2) % aiTutorLesson.words.length],
        word,
        aiTutorLesson.words[(index + 3) % aiTutorLesson.words.length],
        aiTutorLesson.words[(index + 1) % aiTutorLesson.words.length],
      ].map((item) => ({ id: item.id, text: item.target }));
      return { ...base, kind: "word_recall", title: word.interface_translation, instruction: "Recall the target word.", word, options };
    }
    if (stage === "production") return { ...base, kind: "free_text", title: "Your sentences", instruction: aiTutorLesson.production_task.instruction_interface };
    if (stage === "lesson_feedback") return { ...base, kind: "rating", title: "Lesson feedback", instruction: "How did this lesson feel?", options: ["easy", "good", "hard", "bad"].map((id) => ({ id, text: id.replace("_", " ") })) };
    if (stage === "review_schedule") return { ...base, kind: "review", title: "Review", instruction: "Choose review time.", options: ["tomorrow", "3_days", "1_week", "no_review"].map((id) => ({ id, text: id.replace("_", " ") })) };
    return { ...base, kind: "complete", title: "Complete", instruction: "Lesson complete." };
  };
  const aiTutorResponse = (feedback: Record<string, unknown> = {}) => {
    const complete = aiTutorStage === "complete";
    const baseUser = (testSessionPayloadOverride || sessionPayload).user;
    const awarded = complete ? 40 : 0;
    return {
      session: { id: "session-ai-tutor", current_stage: aiTutorStage, status: complete ? "complete" : "active" },
      lesson: aiTutorLesson,
      lesson_status: "approved",
      current_stage: aiTutorStage,
      next_step: aiTutorStep(aiTutorStage),
      feedback,
      ...(awarded ? { xp: awarded, reward_xp: awarded, reward_title: "AI Tutor" } : {}),
      user: awarded ? { ...baseUser, xp: Number(baseUser.xp || 0) + awarded, xp_current: Number(baseUser.xp_current || 0) + awarded } : baseUser,
    };
  };
  await page.route("**/api/ai-tutor/start", (route) => {
    aiTutorStage = "story_intro";
    return route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(aiTutorResponse()) });
  });
  await page.route("**/api/ai-tutor/completed", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        items: [
          {
            session_id: "session-ai-tutor-completed",
            lesson_id: "lesson-ai-tutor-morning",
            title: aiTutorLesson.title,
            topic: aiTutorLesson.theme,
            level: aiTutorLesson.level,
            completed_at: "2026-06-12T12:00:00Z",
            lesson: aiTutorLesson,
          },
        ],
        total: 1,
      }),
    }),
  );
  await page.route("**/api/ai-tutor/answer", async (route) => {
    const body = route.request().postDataJSON() as { text?: string; choice?: string };
    if (aiTutorStage === "production" && String(body.text || "").split(/[.!?]+/).map((part) => part.trim()).filter(Boolean).length < 2) {
      return route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(aiTutorResponse({ ok: false, message: "Напишите минимум два полных предложения, чтобы продолжить." })),
      });
    }
    const currentIndex = aiTutorStages.indexOf(aiTutorStage);
    aiTutorStage = aiTutorStages[Math.min(currentIndex + 1, aiTutorStages.length - 1)];
    const feedback = aiTutorStage === "lesson_feedback"
      ? { ok: true, message: "Повторите wake up завтра.", json: JSON.stringify({ corrected_version_target: "I wake up at seven. I get ready quickly.", recommendations_interface: ["Repeat wake up tomorrow."], mistakes: [{ correction: "I wake up at seven.", explanation: "Use at + time." }] }) }
      : { ok: true, message: "Ответ сохранён" };
    return route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(aiTutorResponse(feedback)) });
  });
  await page.route("**/api/ai-tutor/review", async (route) => {
    aiTutorStage = "complete";
    return route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(aiTutorResponse({ ok: true, message: "Повторение запланировано" })) });
  });
  await page.route("**/api/ai-tutor/finish", async (route) => {
    aiTutorStage = "complete";
    return route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(aiTutorResponse({ ok: true, message: "Урок завершён" })) });
  });
  await page.route("**/api/words/next", (route) => {
    const round = wordRounds[wordRoundIndex % wordRounds.length];
    wordRoundIndex += 1;
    const correctOption = round.options.find((option) => option.id === round.correct);
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        empty: false,
        prompt: round.prompt,
        context: "Pick the target-language word.",
        direction: "RU -> EN",
        options: round.options,
        correct_answer_id: round.correct,
        instruction: "Choose answer",
        word_id: round.correct,
        word: correctOption?.text,
        translation: round.prompt,
        reportable: true,
      }),
    });
  });
  await page.route("**/api/words/report", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ ok: true, word_report: { id: "report-word-test" } }),
    }),
  );
  await page.route("**/api/words/answer", async (route) => {
    const answer = route.request().postDataJSON();
    const answerId = String(answer.answer_id || "");
    const isTrain = answerId.includes("train");
    const correct = answerId.includes("apple") || isTrain;
    const word = isTrain ? "train" : "apple";
    const translation = isTrain ? "\u043f\u043e\u0435\u0437\u0434" : "\u044f\u0431\u043b\u043e\u043a\u043e";
    const example = isTrain ? "I caught the train." : "I bought an apple.";
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        correct,
        learned: correct,
        total: 1,
        word_id: answerId,
        word: correct ? word : "",
        translation,
        context: "Saved.",
        example,
        xp: correct ? 15 : 0,
        user: correct ? { ...sessionPayload.user, xp: sessionPayload.user.xp + 15 } : sessionPayload.user,
      }),
    });
  });
  await page.route("**/api/vocabulary**", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ items: vocabularyItems, page: { page: 0, page_size: 20, total: vocabularyItems.length, total_pages: 1 } }),
    }),
  );
  await page.route("**/api/level-test/start", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        question: {
          index: 2,
          total: 36,
          question: "What does 'to bring up a topic' mean?",
          options: ["to mention it", "to lift it", "to forget it", "to prove it", "I do not know"],
        },
      }),
    }),
  );
  await page.route("**/api/level-test/answer", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ complete: true, level: "B1", score: 12, max_score: 36 }),
    }),
  );
  await page.route("**/api/phrasebook**", async (route) => {
    const request = route.request();
    const phrasebookItems = testPhrasebookItemsOverride || phrasebookSeed;
    if (request.method() === "GET") {
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: phrasebookItems }) });
      return;
    }
    if (request.method() === "POST") {
      const item = request.postDataJSON();
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, items: [item, ...phrasebookItems] }) });
      return;
    }
    if (request.method() === "DELETE") {
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, items: [] }) });
      return;
    }
    await route.fulfill({ status: 405, contentType: "application/json", body: JSON.stringify({ error: "method not allowed" }) });
  });
  await page.route("**/api/mistakes", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        items: mistakeItems,
      }),
    }),
  );
  await page.route("**/api/mistakes/practice/start", async (route) => {
    const body = route.request().postDataJSON();
    const index = Number(body.index ?? 0);
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ empty: false, mistake: mistakeItems[Math.max(0, Math.min(mistakeItems.length - 1, index))] }),
    });
  });
  await page.route("**/api/mistakes/practice/answer", async (route) => {
    const body = route.request().postDataJSON();
    const correct = String(body.text || "").toLowerCase().includes("correct phrase");
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        correct,
        xp: correct ? 8 : 0,
        mistake: mistakeItems[0],
        correction: "Correct phrase 1.",
        explanation: correct ? "Mistake repaired." : "Try the corrected phrase.",
      }),
    });
  });
  await page.route("**/api/spelling/start", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ empty: false, word_id: "spell-clear", prompt: "ясный", context: "Used for something easy to understand.", direction: "RU -> EN" }),
    }),
  );
  await page.route("**/api/spelling/answer", async (route) => {
    const body = route.request().postDataJSON();
    const gaveUp = Boolean(body.give_up);
    const correct = !gaveUp && String(body.text || "").trim().toLowerCase() === "clear";
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        correct,
        gave_up: gaveUp,
        word_id: "spell-clear",
        word: "clear",
        correct_answer: "clear",
        translation: "ясный",
        context: "Used for something easy to understand.",
        example: "The singer has a clear voice.",
        message: correct ? "Great spelling!" : "Correct answer shown without XP.",
      }),
    });
  });
  await page.route("**/api/lesson/start", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ prompt: "Say that you have a reservation.", example: "I have a reservation under the name Ivan Petrov.", question_audio_text: "Could you please show me your passport?" }),
    }),
  );
  await page.route("**/api/lesson/answer", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ correction: "I have a reservation under the name Ivan Petrov.", explanation: "Use under the name for reservations.", correction_audio_text: "I have a reservation under the name Ivan Petrov." }),
    }),
  );
  await page.route("**/api/practice", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ correction: "Could you repeat that, please?", explanation: "Short polite request.", question_audio_text: "?" }),
    }),
  );
  await page.route("**/api/shadowing/start", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ phrase: "Could you repeat that, please?", target: "Could you repeat that, please?" }),
    }),
  );
  let pronunciationSampleIndex = 0;
  const pronunciationSamples = [
    "Could you repeat that, please?",
    "That meeting starts at five.",
    "We are leaving on Friday.",
  ];
  await page.route("**/api/pronunciation/start", (route) => {
    const phrase = pronunciationSamples[pronunciationSampleIndex % pronunciationSamples.length];
    pronunciationSampleIndex += 1;
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ phrase, target: phrase, user: (testSessionPayloadOverride || sessionPayload).user }),
    });
  });
  await page.route("**/api/shadowing/answer", (route) => {
    const body = route.request().postData() || "";
    const isTutorListening = body.includes("Good morning") || body.includes("coffee");
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        transcript: isTutorListening ? "menu coffee breakfast" : "Could you repeat that please",
        pronunciation: {
          score: 67,
          accent_strength: 42,
          fluency: 74,
          stress: "Stress is late on coffee.",
          rhythm: "Rhythm has one long pause.",
          intonation: "Intonation rises naturally.",
          feedback: "Good rhythm. Work on endings.",
          problem_words: [
            { word: "coffee", confidence: 0.42, issue: "low_confidence", tip: "Keep the first syllable clear." },
            { word: "breakfast", confidence: 0.39, issue: "low_confidence", tip: "Do not drop the final sound." },
            { word: "please", confidence: 0.28, issue: "missing", tip: "" },
          ],
          phoneme_issues: [
            { word: "coffee", expected_sound: "/f/", heard_sound: "/v/", confidence: 0.58, tip: "Keep /f/ unvoiced." },
          ],
          corrected_text: "Good morning. Could I see the menu? I'd like coffee for breakfast, please.",
        },
        correction_audio_text: "Good morning. Could I see the menu? I'd like coffee for breakfast, please.",
      }),
    });
  });
  await page.route("**/api/pronunciation/check", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        transcript: "menu coffee breakfast",
        pronunciation: {
          score: 67,
          accent_strength: 42,
          fluency: 74,
          stress: "Stress is late on coffee.",
          rhythm: "Rhythm has one long pause.",
          intonation: "Intonation rises naturally.",
          feedback: "Good rhythm. Work on endings.",
          problem_words: [
            { word: "coffee", confidence: 0.42, issue: "low_confidence", tip: "Keep the first syllable clear." },
            { word: "breakfast", confidence: 0.39, issue: "low_confidence", tip: "Do not drop the final sound." },
            { word: "please", confidence: 0.28, issue: "missing", tip: "" },
          ],
          phoneme_issues: [
            { word: "coffee", expected_sound: "/f/", heard_sound: "/v/", confidence: 0.58, tip: "Keep /f/ unvoiced." },
          ],
          corrected_text: "Good morning. Could I see the menu? I'd like coffee for breakfast, please.",
        },
        correction_audio_text: "Good morning. Could I see the menu? I'd like coffee for breakfast, please.",
      }),
    }),
  );
  await page.route("**/api/daily/claim", (route) =>
    route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, claimed: true, xp: 25, streak: 1, user: { ...sessionPayload.user, xp: 660 } }) }),
  );
  await page.route("**/api/bug-report", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, id: "test-report" }) }));
  await page.route("**/api/premium/plans", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ plans: sessionPayload.premium_plans, enabled: true, crypto_enabled: true }) }));
  await page.route("**/api/premium/crypto/payment", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        payment_id: "pay-ton-1",
        amount: "1.5",
        currency: "TON",
        network: "TON",
        address: "UQCDkENqCLcFLvAzPHX8LxfK6wdlPEmpFQ5DQ5UhG9BIZQ3i",
        comment: "POLIGLOT:test",
        status: "pending",
        method: "TON",
        expires_at: "2026-05-27T01:38:00+03:00",
      }),
    }),
  );
  await page.route("**/api/premium/crypto/check", (route) =>
    route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ payment_id: "pay-ton-1", status: "pending", paid: false }) }),
  );
  await page.route("**/api/leaderboard**", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: leaderboardItems, overall: leaderboardItems, meta: { language: "en" } }) }));
}

async function expectNoMojibake(page: Page) {
  const text = await page.locator("body").innerText();
  expect(text).not.toMatch(mojibakePattern);
  expect(text).not.toMatch(broadMojibakePattern);
}

async function useInterfaceLanguage(page: Page, code: string) {
  const base = testSessionPayloadOverride || sessionPayload;
  const payload = {
    ...base,
    user: { ...base.user, interface_language: code },
  };
  testSessionPayloadOverride = payload;
  await page.unroute("**/api/session").catch(() => undefined);
  await page.unroute("**/api/settings").catch(() => undefined);
  await page.unroute("**/api/navigation-layout").catch(() => undefined);
  await page.route("**/api/session", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(payload) }));
  await page.route("**/api/settings", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(payload) }));
  await page.route("**/api/navigation-layout", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(payload) }));
}

async function mockAnonymousAuth(page: Page) {
  await page.unroute("**/api/session");
  await page.addInitScript(() => {
    const popup = { opener: null, location: { href: "" }, close: () => undefined };
    Object.defineProperty(window, "open", { value: () => popup, writable: true });
    (window as unknown as { turnstile: unknown }).turnstile = {
      render: (_target: HTMLElement, options: { callback?: (token: string) => void }) => {
        window.setTimeout(() => options.callback?.("test-captcha-token"), 0);
        return "test-widget";
      },
      reset: () => undefined,
      remove: () => undefined,
    };
  });
  await page.route("**/api/session", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        authenticated: false,
        interface_languages: sessionPayload.interface_languages,
        learning_languages: sessionPayload.learning_languages,
        captcha: { enabled: true, provider: "turnstile", site_key: "1x00000000000000000000AA" },
        telegram_login_bot: "poliglot_ai_bot",
      }),
    }),
  );
}

function usePremiumSession(overrides: Partial<typeof sessionPayload.user> = {}) {
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      plan: "premium",
      premium: true,
      premium_until: "2026-07-12T12:00:00Z",
      lesson_limit: 50,
      practice_limit: 200,
      voice_limit: 20,
      ...overrides,
    },
  };
}

test.beforeEach(async ({ page }) => {
  testSessionPayloadOverride = null;
  testPhrasebookItemsOverride = null;
  await mockApi(page);
});

test("app shell embeds branded loader before JavaScript hydrates", () => {
  const html = readFileSync("index.html", "utf8");
  expect(html).toContain('class="poliglot-boot"');
  expect(html).toContain("brand-logo-mini.png");
  expect(html).toContain("Poliglot AI");
  expect(html).toContain("poliglotBootSpin");
  expect(html).toContain("poliglotBootOrbit");
  expect(html).toContain("poliglotBootBackground");
  expect(html).toContain("poliglotBootDrift");
});

test("PWA service worker cache is bumped for the AI Router desktop release", () => {
  const worker = readFileSync("public/offline-deck-sw.js", "utf8");
  expect(worker).toContain('const CACHE_NAME = "poliglot-v2-offline-decks-20260616-ai-router"');
});

test("v2 required labels are localized for all 35 interface languages", () => {
  const keys = [
    "roleplay",
    "phrasebook",
    "offline_decks",
    "dashboard",
    "referrals",
    "daily_quests",
    "weekly_plan",
    "daily_route",
    "v2_learning_lab",
    "v2_learning_lab_title",
    "tutor_short_hint",
    "offline_decks_title",
    "offline_short_hint",
    "send_code",
    "refresh_deck",
    "payment_history_title",
    "date",
    "plan",
    "method",
    "status",
    "nav_lesson_desc",
    "nav_practice_desc",
    "nav_roleplay_desc",
    "nav_phrasebook_desc",
    "nav_offline_desc",
    "nav_tools_desc",
    "nav_mistakes_desc",
    "nav_spelling_desc",
    "nav_progress_desc",
    "choose_another_scenario",
    "choose_tool",
    "skip",
    "cancel",
    "save_to_phrasebook",
    "phrasebook_title",
    "empty_phrasebook",
    "auth_fill_required",
    "auth_telegram_title",
    "auth_login_hint",
    "auth_register_hint",
    "auth_support",
    "auth_subtitle",
    "auth_referral_placeholder",
    "auth_finish_account",
    "auth_words",
    "auth_card_lessons",
    "auth_card_words",
    "auth_card_telegram",
    "auth_privacy_consent",
    "auth_privacy_policy",
    "auth_privacy_required",
    "choose_mistake",
    "back_to_mistake_list",
    "correct_spelling",
    "train_similar",
    "confirm_clear_mistakes",
    "yes_clear",
    "mistake_groups",
    "no_mistakes_loaded",
    "referral_invitee_terms",
    "referral_inviter_terms",
    "referral_direct_terms",
    "referral_indirect_terms",
    "referral_withdraw_terms",
    "referral_invitees",
    "referral_invitees_title",
    "level_3_status",
    "earned",
    "reached",
    "not_yet",
    "no_referral_invitees",
    "roleplay_loading",
    "roleplay_answer_placeholder",
    "logout",
    "lessons_today",
    "practice_today",
    "voices_today",
    "smart_training_plan",
    "today_plan_title",
    "today_plan_body_active",
    "today_plan_body_steady",
    "today_plan_words_title",
    "today_plan_words_body",
    "today_plan_speak_title",
    "today_plan_speak_body",
    "today_plan_listen_title",
    "today_plan_listen_body",
    "today_plan_repair_title",
    "today_plan_repair_body",
    "monday_short",
    "tuesday_short",
    "wednesday_short",
    "thursday_short",
    "friday_short",
    "saturday_short",
    "sunday_short",
    "plan_words",
    "plan_words_count",
    "plan_practice",
    "plan_practice_count",
    "plan_listening",
    "plan_listening_count",
    "plan_roleplay",
    "plan_roleplay_count",
    "plan_pronunciation",
    "plan_pronunciation_count",
    "plan_offline",
    "plan_offline_count",
    "plan_dashboard",
    "weekly_review",
    "daily_quests_title",
    "weekly_plan_title",
    "habit_calendar_title",
    "roleplay_title",
    "roleplay_subtitle",
    ...roleplayScenarioIds.flatMap((id) => [`roleplay_scenario_${id}`, `roleplay_scenario_${id}_description`]),
    "mistake_group_grammar",
    "mistake_group_word_order",
    "mistake_group_vocabulary",
    "mistake_group_politeness",
    "mistake_group_spelling",
    "security",
    "change_password",
    "current_password",
    "new_password",
    "confirm_password",
    "password_min_hint",
    "passwords_do_not_match",
    "passwords_match",
    "save_password",
    "auth_password_hint",
    "auth_password_confirm_hint",
    "auth_password_placeholder",
    "auth_password_confirm_placeholder",
    "pay_stars",
    "payment_stars_instruction",
    "payment_crypto_prepare_instruction",
    "payment_crypto_instruction",
    "payment_usdt_trc20_instruction",
    "free_plan_title",
    "free_plan_tier",
    "free_plan_body",
    "free_plan_price",
    "current_plan",
    "free_feature_daily",
    "free_feature_phrasebook",
    "free_feature_no_ai_tutor",
    "free_feature_no_audio",
    "report_bug",
    "report_bug_body",
    "problem_description",
    "problem_description_placeholder",
    "attach_screenshot",
    "screenshots_attached",
    "bug_report_close",
    "send_report",
    "bug_report_image_only",
    "bug_report_sent",
    "guide",
    "app_guide_title",
    "app_guide_body",
    "app_guide_step_1_title",
    "app_guide_step_1_body",
    "app_guide_step_2_title",
    "app_guide_step_2_body",
    "app_guide_step_3_title",
    "app_guide_step_3_body",
    "app_guide_step_4_title",
    "app_guide_step_4_body",
    "app_guide_step_5_title",
    "app_guide_step_5_body",
    "xp_gained_label",
    "tutor_reward_saved",
    "tutor_completed_lessons",
    "tutor_completed_lessons_body",
    "tutor_completed_lessons_empty",
    "tutor_step_pronunciation",
    "tutor_step_final_check",
    "tutor_step_assessment",
    "tutor_final_check_instruction",
    "tutor_final_check_progress",
    "tutor_final_check_done",
    "tutor_final_check_results",
    "tutor_pronunciation_instruction",
    "tutor_pronunciation_target",
    "tutor_need_voice",
    "tutor_voice_ready",
    "tutor_voice_hint",
    "tutor_pronunciation_retry",
    "tutor_pronunciation_feedback",
    "tutor_audio_word",
    "tutor_audio_example",
    "tutor_audio_question",
    "tutor_audio_system",
    "tutor_final_assessment",
    "tutor_final_assessment_instruction",
    "tutor_final_assessment_review",
    "tutor_final_assessment_empty",
  ];
  expect(appLocaleCodes).toHaveLength(35);
  for (const code of appLocaleCodes) {
    for (const key of keys) {
      const value = appCopy(code, key);
      expect(value, `${code}.${key}`).toBeTruthy();
      expect(value, `${code}.${key}`).not.toMatch(mojibakePattern);
    }
  }
  for (const code of appLocaleCodes) {
    expect(appCopy(code, "pay_stars"), `${code}.pay_stars keeps branded payment label`).toBe("Telegram Stars");
    expect(appCopy(code, "ai_router", "AI Router"), `${code}.ai_router keeps branded tool label`).toBe("AI Router");
  }
  expect(appCopy("vi", "roleplay_scenario_restaurant")).toBe("Nhà hàng");
  expect(appCopy("vi", "roleplay_title")).toBe("Kịch bản nhập vai AI");
  expect(appCopy("ka", "logout")).toBe("გასვლა");
  expect(appCopy("ka", "roleplay_scenario_restaurant")).toBe("რესტორანი");
  for (const code of ["vi", "ka"] as const) {
    const visibleCritical = [
      "logout",
      "pronunciation",
      "phrasebook",
      "offline_decks",
      "dashboard",
      "referrals",
      "awards",
      "roleplay_title",
      "roleplay_subtitle",
      "lessons_today",
      "practice_today",
      "daily_route",
      "weekly_plan_title",
      "today_plan_title",
      "payment_stars_instruction",
      "payment_crypto_instruction",
      "payment_usdt_trc20_instruction",
      "change_password",
      "password_min_hint",
      "roleplay_scenario_restaurant",
      "roleplay_scenario_work",
      "roleplay_scenario_description",
    ];
    for (const key of visibleCritical) {
      const value = appCopy(code, key);
      expect(value, `${code}.${key}`).not.toMatch(/^(Section|Раздел|Mục):/);
      expect(value, `${code}.${key}`).not.toMatch(/AI roleplay scenarios|Choose a situation|Lessons today|Practice today|Daily route|Weekly plan|Dashboard|Invite|Notes|Pronunciation|Offline|Restaurant|Work call/);
    }
  }
  const noEnglishFallbackKeys = [
    "ready_title",
    "ready_body",
    "request_failed",
    "lesson_ready",
    "practice_feedback",
    "language_cockpit",
    "learning_group",
    "words_group",
    "growth_group",
    "account_group",
    "training_signal",
    "practice_prompt",
    "lesson_answer_placeholder",
    "no_vocabulary_loaded",
    "language_leaderboard",
    "global_top",
    "payment_options",
    "learning_settings",
    "text_translator_body",
    "empty_panel",
    "view_home_title",
    "view_home_subtitle",
    "view_lesson_title",
    "view_lesson_subtitle",
    "view_practice_title",
    "view_practice_subtitle",
    "view_settings_title",
    "view_settings_subtitle",
    "nav_lesson_desc",
    "nav_practice_desc",
    "nav_tools_desc",
    "learning_activity",
    "progress_without_leaderboard",
    "weekly_rhythm",
    "training_balance",
    "words_learned",
    "phrases_practiced",
    "lessons_completed",
    "review_rounds",
    "voice_attempts",
    "mistakes_repaired",
    "xp_level",
    "lessons_total",
    "practices_total",
    "lessons_today",
    "practice_today",
    "voices_today",
    "phrasebook_note_placeholder",
    "premium_month_title",
    "platinum_month_title",
    "premium_year_title",
    "today_plan_title",
    "today_plan_body_active",
    "today_plan_body_steady",
    "change_password",
    "password_min_hint",
    "payment_stars_instruction",
    "payment_crypto_instruction",
    "payment_usdt_trc20_instruction",
    "tutor_step_pronunciation",
    "tutor_step_final_check",
    "tutor_step_assessment",
    "tutor_final_check_instruction",
    "tutor_final_check_progress",
    "tutor_final_check_done",
    "tutor_final_check_results",
    "tutor_pronunciation_instruction",
    "tutor_pronunciation_target",
    "tutor_need_voice",
    "tutor_voice_ready",
    "tutor_voice_hint",
    "tutor_pronunciation_retry",
    "tutor_pronunciation_feedback",
    "tutor_audio_word",
    "tutor_audio_example",
    "tutor_audio_question",
    "tutor_audio_system",
    "tutor_final_assessment",
    "tutor_final_assessment_instruction",
    "tutor_final_assessment_review",
    "tutor_final_assessment_empty",
  ];
  for (const code of appLocaleCodes.filter((locale) => locale !== "en")) {
    for (const key of noEnglishFallbackKeys) {
      expect(appCopy(code, key), `${code}.${key} should not fall back to English`).not.toBe(appCopy("en", key));
    }
  }
  expect(appCopy("ru", "roleplay")).toBe("Ролевая");
  expect(appCopy("ru", "offline_decks")).toBe("Оффлайн");
  expect(appCopy("ru", "dashboard")).toBe("Статистика");
  expect(appCopy("ru", "referrals")).toBe("Приглашения");
  expect(appCopy("ru", "phrasebook")).toBe("Заметки");
  expect(appCopy("ru", "cancel")).toBe("Отмена");
  expect(appCopy("en", "cancel")).toBe("Cancel");
  expect(appCopy("ru", "global_top")).toBe("Общий");
  expect(appCopy("en", "phrasebook")).toBe("Notes");
  expect(appCopy("en", "save_to_phrasebook")).toBe("Save to notes");
  expect(appCopy("ru", "auth_login_hint")).toBe("Введите логин и пароль или подтвердите вход через Telegram.");
  expect(appCopy("ru", "auth_subtitle")).toBe("Вход или регистрация");
  expect(appCopy("ru", "v2_learning_lab")).toBe("Новые функции");
  expect(appCopy("ru", "v2_learning_lab_title")).toBe("Из последнего крупного обновления");
  const noSectionLeakKeys = [
    "auth_login_hint",
    "auth_register_hint",
    "auth_subtitle",
    "auth_telegram_title",
    "auth_privacy_consent",
    "auth_privacy_policy",
    "monday_short",
    "tuesday_short",
    "wednesday_short",
    "thursday_short",
    "friday_short",
    "saturday_short",
    "sunday_short",
    "v2_learning_lab",
    "v2_learning_lab_title",
    "tutor_short_hint",
    "plan_words",
    "plan_words_count",
    "plan_practice",
    "plan_practice_count",
    "plan_listening",
    "plan_listening_count",
    "plan_roleplay",
    "plan_roleplay_count",
    "plan_pronunciation",
    "plan_pronunciation_count",
    "plan_offline",
    "plan_offline_count",
    "plan_dashboard",
    "weekly_review",
    "report_bug",
    "report_bug_body",
    "problem_description",
    "problem_description_placeholder",
    "attach_screenshot",
    "screenshots_attached",
    "bug_report_close",
    "send_report",
    "bug_report_image_only",
    "bug_report_sent",
    "guide",
    "app_guide_title",
    "app_guide_body",
    "app_guide_step_1_title",
    "app_guide_step_1_body",
    "app_guide_step_2_title",
    "app_guide_step_2_body",
    "app_guide_step_3_title",
    "app_guide_step_3_body",
    "app_guide_step_4_title",
    "app_guide_step_4_body",
    "app_guide_step_5_title",
    "app_guide_step_5_body",
    "xp_gained_label",
    "tutor_reward_saved",
    "tutor_completed_lessons",
    "tutor_completed_lessons_body",
    "tutor_completed_lessons_empty",
  ];
  for (const code of appLocaleCodes) {
    for (const key of noSectionLeakKeys) {
      const value = appCopy(code, key);
      expect(value, `${code}.${key} must not expose a generic section fallback`).not.toMatch(/^(Section|Раздел|Mục)(:|$)/);
      expect(value, `${code}.${key} must not expose the i18n key`).not.toBe(key);
    }
  }
  const collapsedLabelKeys = [
    "learning_activity",
    "progress_without_leaderboard",
    "weekly_rhythm",
    "training_balance",
    "words_learned",
    "phrases_practiced",
    "lessons_completed",
    "review_rounds",
    "voice_attempts",
    "mistakes_repaired",
    "lessons_today",
    "practice_today",
    "voices_today",
  ];
  for (const code of appLocaleCodes) {
    const dashboard = appCopy(code, "dashboard");
    const premium = appCopy(code, "premium");
    const phrasebook = appCopy(code, "phrasebook");
    for (const key of collapsedLabelKeys) {
      expect(appCopy(code, key), `${code}.${key} must not collapse to dashboard`).not.toBe(dashboard);
    }
    expect(appCopy(code, "phrasebook_note_placeholder"), `${code}.phrasebook_note_placeholder must not collapse to notes`).not.toBe(phrasebook);
    expect(appCopy(code, "premium_month_title"), `${code}.premium_month_title must not collapse to premium`).not.toBe(premium);
    expect(appCopy(code, "platinum_month_title"), `${code}.platinum_month_title must not collapse to premium`).not.toBe(premium);
  }
});

test("roleplay scenario copy and mistake category labels stay localized for all 35 languages", () => {
  const englishFallbacks = [
    "Hotel check-in",
    "Shopping",
    "Doctor visit",
    "Job interview",
    "Bank and payment",
    "Ask for size",
    "Describe symptoms",
    "Introduce experience",
    "Ask about a card",
  ];
  const mistakeGroupKeys = [
    "mistake_group_grammar",
    "mistake_group_word_order",
    "mistake_group_vocabulary",
    "mistake_group_politeness",
    "mistake_group_spelling",
  ];

  for (const code of appLocaleCodes) {
    for (const id of roleplayScenarioIds) {
      const title = appCopy(code, `roleplay_scenario_${id}`);
      const description = appCopy(code, `roleplay_scenario_${id}_description`);
      expect(title, `${code}.${id}.title`).toBeTruthy();
      expect(description, `${code}.${id}.description`).toBeTruthy();
      expect(title, `${code}.${id}.title key leak`).not.toBe(`roleplay_scenario_${id}`);
      expect(description, `${code}.${id}.description key leak`).not.toBe(`roleplay_scenario_${id}_description`);
      expect(title, `${code}.${id}.title mojibake`).not.toMatch(mojibakePattern);
      expect(description, `${code}.${id}.description mojibake`).not.toMatch(mojibakePattern);
      if (code !== "en") {
        expect(description, `${code}.${id}.description should be localized`).not.toBe(appCopy("en", `roleplay_scenario_${id}_description`));
      }
      if (code === "ru") {
        for (const fragment of englishFallbacks) {
          expect(`${title} ${description}`, `ru.${id} leaks ${fragment}`).not.toContain(fragment);
        }
      }
    }

    const groupLabels = mistakeGroupKeys.map((key) => appCopy(code, key));
    expect(new Set(groupLabels).size, `${code}.mistake groups are distinct`).toBe(groupLabels.length);
    for (const [index, label] of groupLabels.entries()) {
      expect(label, `${code}.${mistakeGroupKeys[index]}`).toBeTruthy();
      expect(label, `${code}.${mistakeGroupKeys[index]} key leak`).not.toBe(mistakeGroupKeys[index]);
      expect(label, `${code}.${mistakeGroupKeys[index]} mojibake`).not.toMatch(mojibakePattern);
      if (code !== "en") {
        expect(label, `${code}.${mistakeGroupKeys[index]} should be localized`).not.toBe(appCopy("en", mistakeGroupKeys[index]));
      }
    }
  }
});

test("main web screens do not expose mojibake or fallback labels across interface languages", async ({ page }) => {
  test.setTimeout(600_000);
  const views = ["home", "roleplay", "pronunciation", "offline", "mistakes", "leaderboard", "spelling", "tools", "dashboard"] as const;
  const fallbackFragments: Array<string | RegExp> = [
    "AI roleplay scenarios",
    "Choose a situation",
    "Lessons today",
    "Practice today",
    "Voice messages today",
    "Daily route",
    "Weekly plan",
    "Dashboard",
    /\bInvite\b/,
    "Pronunciation",
    "Offline decks",
    /\bWork call\b/,
    "Умный план",
    "Квесты",
    "План обучения",
    "Новые обучающие",
    "Раздел:",
    "Mục:",
  ];
  for (const code of appLocaleCodes.filter((locale) => locale !== "en")) {
    await useInterfaceLanguage(page, code);
    for (const view of views) {
      await page.goto(`/app/?view=${view}`);
      await expect(page.locator(".context-display")).toBeVisible();
      const text = await page.locator("body").innerText();
      if (code !== "ru") expect(text, `${code}.${view} mojibake`).not.toMatch(broadMojibakePattern);
      for (const fragment of fallbackFragments) {
        if (code === "ru" && typeof fragment === "string" && /[\u0400-\u04ff]/.test(fragment)) continue;
        if (typeof fragment === "string") {
          expect(text, `${code}.${view} fallback ${fragment}`).not.toContain(fragment);
        } else {
          expect(text, `${code}.${view} fallback ${fragment}`).not.toMatch(fragment);
        }
      }
    }
  }
});

test("menu labels for spelling progress and mistakes are localized in every non-English UI language", () => {
  const criticalMenuKeys = [
    "ai_tutor",
    "nav_tutor_desc",
    "view_tutor_title",
    "spelling",
    "progress",
    "mistakes",
    "nav_spelling_desc",
    "nav_progress_desc",
    "nav_mistakes_desc",
    "view_spelling_title",
    "view_progress_title",
    "view_mistakes_title",
  ];
  const english = Object.fromEntries(criticalMenuKeys.map((key) => [key, appCopy("en", key)]));

  for (const code of appLocaleCodes) {
    for (const key of criticalMenuKeys) {
      const value = appCopy(code, key);
      expect(value, `${code}.${key}`).toBeTruthy();
      expect(value, `${code}.${key}`).not.toMatch(mojibakePattern);
      if (code !== "en") {
        expect(value.toLowerCase(), `${code}.${key} must not leak English`).not.toBe(String(english[key]).toLowerCase());
      }
    }
  }
});

test("today plan and settings password helper copy stay localized", async ({ page }) => {
  await page.goto("/app/?view=home");
  await expect(page.locator(".home-plan-v2")).toContainText("План обучения на сегодня");
  await expect(page.locator(".home-plan-v2")).toContainText("Сначала закрываем одну ошибку");
  await expect(page.locator(".weekly-plan-v2")).toContainText("Словарная база");
  await expect(page.locator(".weekly-plan-v2")).toContainText("Одна тема");
  await expect(page.locator(".weekly-plan-v2")).not.toContainText("Раздел");
  await expect(page.locator(".learning-lab-v2")).toContainText("Новые функции");
  await expect(page.locator(".learning-lab-v2")).toContainText("Из последнего крупного обновления");
  await expect(page.locator(".learning-lab-v2")).toContainText("AI Репетитор");
  await expect(page.locator(".learning-lab-v2")).not.toContainText("Раздел");

  await page.goto("/app/?view=settings");
  await expect(page.locator(".password-card-v2")).toContainText("Смена пароля");
  await page.locator(".password-card-v2 input").nth(0).fill("old-password");
  await page.locator(".password-card-v2 input").nth(1).fill("short");
  await expect(page.locator(".password-card-v2")).toContainText("Минимум 8 символов");
  await page.locator(".password-card-v2 input").nth(1).fill("new-password");
  await page.locator(".password-card-v2 input").nth(2).fill("other-password");
  await expect(page.locator(".password-card-v2")).toContainText("Пароли не совпадают");
  await page.locator(".password-card-v2 input").nth(2).fill("new-password");
  await expect(page.locator(".password-card-v2")).toContainText("Пароли совпадают");
});

test("tools spacing and free plan Russian copy stay compact and localized", async ({ page, isMobile }) => {
  test.skip(isMobile, "Desktop spacing regression uses the desktop two-column tools layout.");

  await page.goto("/app/?view=premium");
  const freePlan = page.locator(".plan-card-v2.is-free");
  const freeLabel = freePlan.locator("> span");
  await expect(freePlan).toBeVisible();
  await expect(freeLabel).toHaveText("Для начала");
  await expect(freeLabel).toHaveCSS("text-transform", "none");
  await expect(freePlan.locator("h2")).toHaveText("Бесплатный");
  await expect(freePlan).toContainText("Включено");
  await expect(freePlan).toContainText("Базовое текстовое обучение");
  await expect(freePlan).toContainText("Заметки и обзор прогресса");
  await expect(freePlan).toContainText("AI Tutor закрыт до Premium");
  await expect(freePlan).toContainText("Аудирование и произношение закрыты до Premium");
  await expect(freePlan).toContainText("Голосовые проверки и фото-инструменты закрыты до Premium");
  await expect(freePlan).not.toContainText("Попробовать маршрут");
  await expect(freePlan).not.toContainText("Статистика");
  await expect(freePlan).not.toContainText("Раздел");
  await expect(freePlan).not.toContainText("AI Tutor guided lessons");
  await expect(freePlan).not.toContainText("Listening и pronunciation");
  await expect(freePlan).not.toContainText("voice checks и photo tools");

  await page.goto("/app/?view=tools");
  await expect(page.locator(".tool-switch-v2")).toBeVisible();
  await expect(page.locator(".composer-panel-v2")).toBeVisible();
  const switchBox = await page.locator(".tool-switch-v2").boundingBox();
  const composerBox = await page.locator(".composer-panel-v2").boundingBox();
  expect(switchBox).not.toBeNull();
  expect(composerBox).not.toBeNull();
  expect(composerBox!.y - (switchBox!.y + switchBox!.height)).toBeLessThanOrEqual(28);
});

test("bug report paste keeps one clipboard image and no generic section text", async ({ page, isMobile }) => {
  await page.goto("/app/?view=home");
  const reportButton = page.locator(isMobile ? ".mobile-report-button-v2" : ".v2-report-button");
  await expect(reportButton).toBeVisible();
  if (!isMobile) {
    await expect(reportButton).toHaveAccessibleName(ru("report_bug", "Сообщить об ошибке"));
    const box = await reportButton.boundingBox();
    const viewport = page.viewportSize();
    expect(box).not.toBeNull();
    expect(viewport).not.toBeNull();
    expect(box!.x + box!.width).toBeLessThanOrEqual(viewport!.width - 8);
  }
  await reportButton.click();
  const dialog = page.locator(".bug-report-dialog-v2");
  await expect(dialog).toBeVisible();
  await expect(dialog).not.toContainText("Раздел");
  await expect(dialog.locator("textarea")).toHaveAttribute("placeholder", "Опишите Вашу проблему");
  await dialog.locator("textarea").evaluate((node) => {
    const file = new File([new Uint8Array([137, 80, 78, 71])], "paste.png", { type: "image/png", lastModified: 123 });
    const transfer = new DataTransfer();
    transfer.items.add(file);
    const event = new ClipboardEvent("paste", { bubbles: true, cancelable: true });
    Object.defineProperty(event, "clipboardData", { value: { files: transfer.files, items: transfer.items } });
    node.dispatchEvent(event);
  });
  await expect(dialog.locator(".bug-report-dialog-v2__upload")).toContainText("paste.png");
  await expect(dialog.locator(".bug-report-dialog-v2__upload")).not.toContainText("2 ");
});

test("free users see AI Tutor premium paywall instead of starting a lesson", async ({ page }) => {
  let startCalled = false;
  await page.unroute("**/api/ai-tutor/start").catch(() => undefined);
  await page.route("**/api/ai-tutor/start", (route) => {
    startCalled = true;
    return route.fulfill({
      status: 402,
      contentType: "application/json",
      body: JSON.stringify({ error: { code: "premium_required", message: "AI Tutor is available with Premium." } }),
    });
  });

  await page.goto("/app/?view=tutor");
  await expect(page.locator(".tutor-paywall-v2")).toContainText(ru("tutor_premium_title", "AI Tutor is included with Premium"));
  await expect(page.locator(".tutor-session-v2")).toHaveCount(0);
  expect(startCalled).toBe(false);
});

test("AI Tutor server-driven lesson blocks old local flow and awards XP", async ({ page, isMobile }) => {
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      plan: "premium",
      premium: true,
      premium_until: "2026-07-12T12:00:00Z",
      lesson_limit: 50,
      practice_limit: 200,
      voice_limit: 20,
    },
  };
  await page.goto("/app/?view=tutor");
  await expect(page.locator('.function-ribbon [data-view="tutor"]')).toContainText("AI Репетитор");
  await expect(page.locator(".tutor-workspace")).toContainText("Утренняя рутина");
  await expect(page.locator(".tutor-session-v2")).toBeVisible();
  await expect(page.locator(".tutor-step-v2")).toHaveCount(21);
  await expect(page.locator(".tutor-context-v2")).not.toContainText("Раздел");
  await expect(page.locator(".tutor-context-v2__head > span")).not.toContainText("story_intro");
  await expect(page.locator(".tutor-context-v2")).toContainText("Mia wakes up at seven");
  await expect(page.locator(".tutor-context-v2 .audio-wave-button-v2")).toContainText(appCopy("ru", "tutor_story_audio"));

  const tutorSubmit = page.locator(".tutor-context-v2 .tutor-composer-v2 button").last();
  await tutorSubmit.click();
  await expect(page.locator(".tutor-context-v2__head h2")).toContainText(appCopy("ru", "tutor_step_retell"));
  await page.locator(".tutor-context-v2 textarea").fill("Mia wakes up at seven and goes to work.");
  await tutorSubmit.click();
  await expect(page.locator(".tutor-context-v2__head h2")).toContainText("Вопрос 1");
  await expect(page.locator(".tutor-message-v2.is-active")).toContainText("What time does the person wake up?");
  const questionCardText = await page.locator(".tutor-message-v2.is-active").innerText();
  expect((questionCardText.match(/What time does the person wake up\?/g) || []).length).toBe(1);
  expect(questionCardText).not.toContain("Question 1");
  await page.locator(".tutor-context-v2 textarea").fill("At seven.");
  await tutorSubmit.click();
  await expect(page.locator(".tutor-context-v2__head h2")).toContainText("Вопрос 2");
  await page.locator(".tutor-context-v2 textarea").fill("Breakfast.");
  await tutorSubmit.click();
  await expect(page.locator(".tutor-context-v2__head h2")).toContainText("Вопрос 3");
  await page.locator(".tutor-context-v2 textarea").fill("At nine.");
  await tutorSubmit.click();

  await expect(page.locator(".tutor-context-v2__head h2")).toContainText("Новые слова 1");
  await expect(page.locator(".tutor-word-check-v2")).toContainText("wake up");
  await expect(page.locator(".tutor-word-check-v2")).toContainText("просыпаться");
  await expect(page.locator(".tutor-word-check-v2")).not.toContainText("Раздел");
  expect(await page.locator(".tutor-word-check-v2 .phrase-quick-save-v2 button").count()).toBeGreaterThan(1);
  for (let index = 0; index < 6; index += 1) {
    await tutorSubmit.click();
  }

  await expect(page.locator(".tutor-context-v2__head h2")).toContainText("Проверка слов 1");
  await expect(page.locator(".tutor-context-v2")).not.toContainText("Recall the target word.");
  await expect(page.locator(".tutor-word-check-v2")).toContainText("просыпаться");
  await expect(page.locator(".tutor-word-check-v2")).not.toContainText("wake up");
  await expect(page.locator(".tutor-word-check-v2 .phrase-quick-save-v2")).toHaveCount(0);
  const firstRecallOptions = await page.locator(".tutor-srs button").allInnerTexts();
  expect(firstRecallOptions).not.toEqual(["wake up", "get ready", "routine", "leave"]);
  await page.locator(".tutor-srs button").filter({ hasText: "wake up" }).click();
  await tutorSubmit.click();
  for (const word of ["get ready", "routine", "leave", "breakfast", "work"]) {
    await page.locator(".tutor-srs button").filter({ hasText: word }).click();
    await tutorSubmit.click();
  }

  await expect(page.locator(".tutor-context-v2__head h2")).toContainText("Письмо");
  await expect(tutorSubmit).toBeDisabled();
  await page.locator(".tutor-context-v2 textarea").fill("I wake up at seven.");
  await expect(tutorSubmit).toBeDisabled();
  await expect(page.locator(".tutor-feedback-v2.is-error")).toContainText("минимум два");
  await page.locator(".tutor-context-v2 textarea").fill("I wake up at seven. I get ready quickly.");
  await tutorSubmit.click();
  await expect(page.locator(".tutor-context-v2__head h2")).toContainText("Разбор репетитора");
  await expect(tutorSubmit).toContainText(appCopy("ru", "tutor_continue"));
  await expect(tutorSubmit).not.toContainText(appCopy("ru", "tutor_check_answer"));
  await expect(tutorSubmit).toBeEnabled();
  const tutorNotes = page.locator(".tutor-context-v2 .tutor-note-strip-v2");
  await expect(tutorNotes).toHaveCount(1);
  await expect(tutorNotes).toContainText(appCopy("ru", "save_to_phrasebook"));
  await expect(tutorNotes.locator("button").filter({ hasText: "I wake up at seven" }).first()).toBeVisible();
  await expect(tutorNotes.locator("button").filter({ hasText: "Repeat wake up tomorrow" })).toHaveCount(1);
  await expect(page.locator(".tutor-context-v2 .phrase-quick-save-v2")).toHaveCount(0);
  await expect(page.locator(".tutor-context-v2 .tutor-notes-suggestions-v2")).toHaveCount(0);
  await expect(page.locator(".tutor-context-v2")).not.toContainText(/Ошибки:\s*Сохранить в заметки/);
  const tutorErrors = page.locator(".tutor-context-v2 .tutor-actual-errors-v2");
  await expect(tutorErrors).toHaveCount(1);
  await expect(tutorErrors).toContainText("I wake up at seven.: Use at + time.");
  await expect(tutorErrors).not.toContainText(appCopy("ru", "save_to_phrasebook"));
  const composerBeforeNotes = await page.locator(".tutor-context-v2").evaluate((node) => {
    const composer = node.querySelector(".tutor-composer-v2");
    const saveBlock = node.querySelector(".tutor-note-strip-v2");
    return Boolean(composer && saveBlock && (composer.compareDocumentPosition(saveBlock) & Node.DOCUMENT_POSITION_FOLLOWING));
  });
  expect(composerBeforeNotes).toBe(true);

  await page.locator(".tutor-srs button").filter({ hasText: appCopy("ru", "ai_tutor_option_good") }).click();
  await tutorSubmit.click();
  await expect(page.locator(".tutor-context-v2__head h2")).toContainText("Повторение");
  await expect(page.locator(".tutor-context-v2")).not.toContainText("Review");
  await page.locator(".tutor-srs button").filter({ hasText: appCopy("ru", "ai_tutor_option_tomorrow") }).click();
  await tutorSubmit.click();
  await expect(page.locator(".tutor-context-v2__head h2")).toContainText("Урок завершён");
  await expect(page.locator(".tutor-context-v2")).not.toContainText("Complete");
  await expect(page.locator(".xp-gain-pop-v2")).toContainText("+40 XP");
  await expect(page.locator(".xp-gain-pop-v2")).toContainText("AI Репетитор");

  await page.locator(".tutor-completed-lessons-button-v2").click();
  const completedDialog = page.locator(".tutor-completed-lessons-dialog-v2");
  await expect(completedDialog).toBeVisible();
  await expect(completedDialog).toContainText("Пройденные уроки");
  await expect(completedDialog).toContainText("Утренняя рутина");
  await expect(completedDialog).toContainText("A1");
  await completedDialog.locator(".tutor-completed-lessons-back-v2").click();
  await expect(completedDialog).toHaveCount(0);

  await useInterfaceLanguage(page, "de");
  await page.goto("/app/?view=tutor");
  await expect(page.locator(".tutor-context-v2")).toContainText("Geschichte");
  await expect(page.locator(".tutor-context-v2")).not.toContainText("Новые слова");

  await useInterfaceLanguage(page, "ru");
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/app/?view=home");
  await expect(page.locator('.mobile-bottom-nav-v2 [data-view="tutor"]')).toContainText("AI Репетитор");
  await page.locator('.mobile-bottom-nav-v2 [data-view="tutor"]').click();
  await expect(page.locator(".tutor-session-v2")).toBeVisible();
  const mobileTutorMetrics = await page.evaluate(() => {
    const rect = (selector: string) => {
      const box = document.querySelector(selector)?.getBoundingClientRect();
      return box ? { left: box.left, right: box.right, top: box.top, width: box.width } : null;
    };
    return {
      innerWidth: window.innerWidth,
      scrollWidth: document.documentElement.scrollWidth,
      context: rect(".tutor-context-v2"),
      plan: rect(".tutor-plan-v2"),
    };
  });
  expect(mobileTutorMetrics.context).not.toBeNull();
  expect(mobileTutorMetrics.plan).not.toBeNull();
  expect(mobileTutorMetrics.context!.top).toBeLessThan(mobileTutorMetrics.plan!.top);
  expect(mobileTutorMetrics.scrollWidth).toBeLessThanOrEqual(mobileTutorMetrics.innerWidth + 1);
  expect(mobileTutorMetrics.context!.right).toBeLessThanOrEqual(mobileTutorMetrics.innerWidth + 1);
  expect(mobileTutorMetrics.plan!.right).toBeLessThanOrEqual(mobileTutorMetrics.innerWidth + 1);
});

test("AI Tutor completed lessons modal paginates finished lesson history", async ({ page }) => {
  await page.unroute("**/api/ai-tutor/completed").catch(() => undefined);
  await page.route("**/api/ai-tutor/completed", (route) =>
    route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [], total: 0 }) }),
  );
  await page.addInitScript(() => {
    const lessons = Array.from({ length: 12 }, (_, index) => ({
      id: `completed-${index + 1}`,
      title: "AI Репетитор",
      topic: `Cafe ${index + 1}`,
      level: index % 2 === 0 ? "A1" : "A2",
      completedAt: new Date(Date.now() - index * 1000).toISOString(),
      lesson: {
        id: `completed-${index + 1}`,
        title: "AI Репетитор",
        topic: `Cafe ${index + 1}`,
        level: index % 2 === 0 ? "A1" : "A2",
        learning_language: "en",
        interface_language: "ru",
        words: [],
      },
    }));
    localStorage.setItem("poliglot-tutor-completed-v3:demor22", JSON.stringify(lessons));
  });
  await page.goto("/app/?view=tutor");
  await page.locator(".tutor-completed-lessons-button-v2").click();
  const dialog = page.locator(".tutor-completed-lessons-dialog-v2");
  const lessonTitles = dialog.locator(".tutor-completed-lesson-v2 strong");
  await expect(dialog).toBeVisible();
  await expect(dialog.locator(".tutor-completed-lesson-v2")).toHaveCount(10);
  await expect(lessonTitles).toHaveText(["Cafe 1", "Cafe 2", "Cafe 3", "Cafe 4", "Cafe 5", "Cafe 6", "Cafe 7", "Cafe 8", "Cafe 9", "Cafe 10"]);
  await dialog.locator(".tutor-completed-lessons-next-v2").click();
  await expect(lessonTitles).toHaveText(["Cafe 11", "Cafe 12"]);
  await expect(dialog.locator(".tutor-completed-lesson-v2")).toHaveCount(2);
  await dialog.locator(".tutor-completed-lessons-back-v2").click();
  await expect(dialog).toHaveCount(0);
});

test("guide button opens localized quick start guide and mobile home shows level XP panel", async ({ page, isMobile }) => {
  await page.goto("/app/?view=home");
  await page.locator(isMobile ? ".mobile-quick-controls-v2 .app-guide-button-v2" : ".v2-topbar .app-guide-button-v2").click();
  const guide = page.locator(".app-guide-dialog-v2");
  await expect(guide).toBeVisible();
  await expect(guide.locator(".app-guide-step-v2")).toHaveCount(5);
  await expect(guide).toContainText("Как начать");
  await guide.locator(".app-guide-back-v2").click();
  await expect(guide).toHaveCount(0);

  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/app/?view=home");
  await expect(page.locator(".home-mobile-level-v2")).toBeVisible();
  await expect(page.locator(".home-mobile-level-v2")).toContainText("A2");
  await expect(page.locator(".home-mobile-level-v2")).toContainText("XP");
  await page.locator(".mobile-quick-controls-v2 .app-guide-button-v2").click();
  await expect(page.locator(".app-guide-dialog-v2")).toBeVisible();
});

test("word trainer shuffles options away from API order", async ({ page }) => {
  await page.goto("/app/?view=words");
  await expect(page.locator(".choice-grid-v2 button")).toHaveCount(4);
  const prompt = await page.locator(".trainer-display h2").innerText();
  const optionTexts = await page.locator(".choice-grid-v2 button").allInnerTexts();
  if (prompt.includes("яблоко")) {
    expect(optionTexts).not.toEqual(["apple", "station", "ticket", "coffee"]);
  } else {
    expect(optionTexts).not.toEqual(["train", "window", "meeting", "market"]);
  }
});

test("word trainer reports the active word with its word id", async ({ page }) => {
  let reportPayload: Record<string, unknown> | null = null;
  await page.unroute("**/api/words/next").catch(() => undefined);
  await page.route("**/api/words/next", (route) => {
    const round = reportPayload
      ? { id: "en:train", prompt: "поезд", word: "train" }
      : { id: "en:apple", prompt: "яблоко", word: "apple" };
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        empty: false,
        prompt: round.prompt,
        context: "Pick the target-language word.",
        direction: "RU -> EN",
        options: [
          { id: round.id, text: round.word },
          { id: "en:station", text: "station" },
          { id: "en:ticket", text: "ticket" },
          { id: "en:coffee", text: "coffee" },
        ],
        correct_answer_id: round.id,
        instruction: "Choose answer",
        word_id: round.id,
        word: round.word,
        translation: round.prompt,
        reportable: true,
      }),
    });
  });
  await page.unroute("**/api/words/report").catch(() => undefined);
  await page.route("**/api/words/report", async (route) => {
    reportPayload = route.request().postDataJSON();
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ ok: true, word_report: { id: "report-word-test" } }),
    });
  });

  await page.goto("/app/?view=words");
  await expect(page.locator(".trainer-display h2")).toContainText("яблоко");
  const reportButton = page.locator(".context-display--words .word-report-action-v2");
  await expect(reportButton).toBeVisible();
  await reportButton.click();

  const dialog = page.locator(".tutor-word-report-dialog-v2");
  await expect(dialog).toBeVisible();
  await expect(dialog).not.toContainText("Раздел");
  await expect(dialog.getByRole("button", { name: "Отмена" })).toBeVisible();
  await dialog.locator("input").nth(0).fill("apple tree");
  await dialog.locator("input").nth(1).fill("яблоня");
  await dialog.locator("textarea").fill("bad dictionary row");
  await dialog.getByRole("button", { name: /Отправить|Send report/ }).click();

  await expect.poll(() => reportPayload).toMatchObject({
    word_id: "en:apple",
    proposed_word: "apple tree",
    proposed_translation: "яблоня",
    comment: "bad dictionary row",
  });
  await expect(page.locator(".trainer-display h2")).toContainText("поезд");
});

test("standalone listening hides the target text and leaves only audio playback", async ({ page }) => {
  usePremiumSession();
  await page.goto("/app/?view=shadowing");
  const listeningPanel = page.locator(".chat-workspace--shadowing .task-box-v2");
  await expect(listeningPanel.locator(".audio-wave-button-v2")).toBeVisible();
  await expect(listeningPanel).not.toContainText("Could you repeat that, please?");
  await expect(listeningPanel.locator("strong")).toHaveCount(0);
});

test("standalone listening reveals the heard phrase after the answer is checked", async ({ page }) => {
  usePremiumSession();
  await page.unroute("**/api/shadowing/answer").catch(() => undefined);
  await page.route("**/api/shadowing/answer", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        feedback: "Почти верно: проверь запятую и please.",
        transcript: "Could you repeat that please",
        target: "Could you repeat that, please?",
        correction_audio_text: "Could you repeat that, please?",
      }),
    }),
  );
  await page.goto("/app/?view=shadowing");
  const listeningPanel = page.locator(".chat-workspace--shadowing .task-box-v2");
  await expect(listeningPanel.locator(".audio-wave-button-v2")).toBeVisible();
  await expect(listeningPanel).not.toContainText("Could you repeat that, please?");

  await page.locator(".composer-panel-v2 textarea").fill("Could you repeat that please");
  await page.locator(".composer-panel-v2").getByRole("button", { name: /Send|Отправить|Надіслати/ }).click();

  const feedback = page.locator(".chat-workspace--shadowing .chat-interface");
  await expect(feedback).toContainText("Почти верно");
  await expect(feedback).toContainText("Could you repeat that, please?");
});

test("pronunciation shows the text-to-pronounce block before the pronunciation summary", async ({ page }) => {
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      plan: "premium",
      premium: true,
      premium_until: "2026-07-12T12:00:00Z",
      voice_limit: 20,
    },
  };
  await page.goto("/app/?view=pronunciation");
  await expect(page.locator(".context-display--pronunciation")).toBeVisible();
  const targetBox = await page.locator(".pronunciation-target-primary-v2").boundingBox();
  const summaryBox = await page.locator(".pronunciation-hero-v2").boundingBox();
  expect(targetBox, "text-to-pronounce block").not.toBeNull();
  expect(summaryBox, "pronunciation summary block").not.toBeNull();
  expect(targetBox!.y).toBeLessThan(summaryBox!.y);
});

test("auth registration shows localized password rules", async ({ page }) => {
  await mockAnonymousAuth(page);
  await page.goto("/app/login");
  await page.locator(".sign-in-page-v2 p a").last().click();
  await expect(page.locator(".sign-in-page-v2")).toContainText(appCopy("en", "auth_password_hint"));
  await expect(page.locator(".sign-in-page-v2")).toContainText(appCopy("en", "auth_password_confirm_hint"));
});

test("premium payment modal explains Stars and exact USDT TRC20 crypto payment", async ({ page }) => {
  await page.route("**/api/premium/stars", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ bot_url: "https://t.me/poliglot_ai_bot?start=stars_invoice", status: "created", message: "Stars invoice created." }),
    }),
  );
  await page.unroute("**/api/premium/crypto/payment").catch(() => undefined);
  await page.route("**/api/premium/crypto/payment", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        payment_id: "pay-usdt-trc20",
        amount: "12.34",
        currency: "USDT",
        network: "TRC20",
        address: "TXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
        comment: "POLIGLOT:test",
        status: "pending",
        method: "USDT_TRC20",
        expires_at: "2026-05-27T01:38:00+03:00",
      }),
    }),
  );

  await page.goto("/app/?view=premium");
  await page.locator(".plan-card-v2:not(.is-free) button").first().click();
  await expect(page.locator(".v2-payment-modal")).toBeVisible();
  await page.locator(".payment-method-button-v2").filter({ hasText: "Telegram Stars" }).click();
  await expect(page.locator(".v2-payment-modal")).toContainText("Оплата Telegram Stars откроется в Telegram");
  await page.locator(".payment-method-button-v2").filter({ hasText: /TON|USDT/ }).click();
  await expect(page.locator(".v2-payment-modal")).toContainText(/Отправьте ровно 12[,.]34 USDT/);
  await expect(page.locator(".v2-payment-modal")).toContainText("USDT TRC20");
  await expect(page.locator(".v2-payment-modal")).toContainText("сеть Tron");
  await expect(page.locator(".v2-payment-modal")).toContainText("Сеть");
  await expect(page.locator(".v2-payment-modal")).toContainText("Комментарий");
  await expect(page.locator(".v2-payment-modal")).toContainText("Истекает");
  await expect(page.locator(".payment-invoice-v2")).not.toContainText("Раздел");
});

test("premium plans show AI Tutor as a paid-only upgrade ladder", async ({ page }) => {
  await page.goto("/app/?view=premium");
  const cards = page.locator(".plan-card-v2");
  await expect(cards).toHaveCount(3);
  await expect(cards.nth(0)).toContainText(ru("free_feature_no_ai_tutor", "AI Tutor is locked until Premium"));
  await expect(cards.nth(0).locator("button")).toBeDisabled();
  await expect(cards.nth(1)).toContainText(ru("premium_feature_ai_tutor", "AI Tutor guided lessons included"));
  await expect(cards.nth(2)).toContainText(ru("platinum_feature_priority", "Best tier for heavy daily learning"));
});

test("auth login page uses React sign-in component, Cloudflare slot and current login API", async ({ page }) => {
  await mockAnonymousAuth(page);
  let loginPayload: Record<string, unknown> | null = null;
  await page.route("**/api/auth/login", async (route) => {
    loginPayload = route.request().postDataJSON();
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(sessionPayload) });
  });

  await page.goto("/app/login");
  await expect(page.locator(".sign-in-page-v2")).toBeVisible();
  await expect(page.locator(".sign-in-page-v2")).not.toContainText("Раздел");
  await expect(page.locator(".sign-in-page-v2")).toContainText(appCopy("en", "auth_welcome"));
  await expect(page.locator(".sign-in-page-v2")).toContainText(appCopy("en", "auth_fill_required"));
  expect(await page.locator(".auth-generative-scene-v2 canvas").count()).toBeGreaterThan(0);
  await expect(page.locator(".sign-in-page-v2__hero-wrap, .sign-in-page-v2__hero, .sign-in-page-v2__side-panel, .sign-in-page-v2__testimonials")).toHaveCount(0);
  const layout = await page.evaluate(() => {
    const toRect = (element: Element | null) => {
      const rect = element?.getBoundingClientRect();
      return rect ? { x: rect.x, y: rect.y, width: rect.width, height: rect.height } : null;
    };
    return {
      viewport: { width: window.innerWidth, height: window.innerHeight },
      card: toRect(document.querySelector(".sign-in-page-v2__form-card")),
      section: toRect(document.querySelector(".sign-in-page-v2__form-section")),
      scene: toRect(document.querySelector(".auth-generative-scene-v2")),
      canvas: toRect(document.querySelector(".auth-generative-scene-v2 canvas")),
    };
  });
  expect(layout.card).not.toBeNull();
  expect(layout.section).not.toBeNull();
  expect(layout.scene).not.toBeNull();
  expect(layout.canvas).not.toBeNull();
  const cardCenterX = layout.card!.x + layout.card!.width / 2;
  const cardCenterY = layout.card!.y + layout.card!.height / 2;
  const sceneCenterX = layout.scene!.x + layout.scene!.width / 2;
  const sceneCenterY = layout.scene!.y + layout.scene!.height / 2;
  expect(Math.abs(cardCenterX - layout.viewport.width / 2)).toBeLessThan(layout.viewport.width * 0.08);
  expect(Math.abs(cardCenterY - layout.viewport.height / 2)).toBeLessThan(layout.viewport.height * 0.14);
  expect(Math.abs(sceneCenterX - layout.viewport.width / 2)).toBeLessThan(2);
  expect(Math.abs(sceneCenterY - layout.viewport.height / 2)).toBeLessThan(2);
  expect(layout.canvas!.width).toBeGreaterThanOrEqual(layout.viewport.width - 2);
  expect(layout.canvas!.height).toBeGreaterThanOrEqual(layout.viewport.height - 2);
  await expect(page.locator('[data-testid="auth-captcha"]')).toBeVisible();
  await page.locator('input[name="login"]').fill("demor22");
  await page.locator('input[name="password"]').fill("strong-password");
  await expect.poll(() => loginPayload).toBeNull();
  await page.locator('button[type="submit"]').click();
  await expect.poll(() => loginPayload).not.toBeNull();
  expect(loginPayload).toMatchObject({ login: "demor22", password: "strong-password", captcha_token: "test-captcha-token" });
  await expect(page).toHaveURL(/\/app$/);
});

test("auth login and registration do not expose Telegram entry", async ({ page }) => {
  await mockAnonymousAuth(page);
  let telegramStartPayload: Record<string, unknown> | null = null;
  await page.route("**/api/auth/telegram/start", (route) => {
    telegramStartPayload = route.request().postDataJSON();
    return route.fulfill({ status: 409, contentType: "application/json", body: JSON.stringify({ error: "telegram disabled" }) });
  });

  await page.goto("/app/login");
  await expect(page.getByRole("button", { name: /Telegram/ })).toHaveCount(0);
  await expect(page.locator(".auth-social-button-v2")).toHaveCount(0);
  await expect(page.locator(".otp-dialog-v2")).toHaveCount(0);
  await page.getByRole("link", { name: appCopy("en", "auth_create_account") }).click();
  await expect(page.getByRole("button", { name: /Telegram/ })).toHaveCount(0);
  expect(telegramStartPayload).toBeNull();
});

test("desktop Today has no top quick buttons, has compact quests and a green streak", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop layout assertion");
  await page.setViewportSize({ width: 1180, height: 825 });
  await page.goto("/app/?view=home");
  await expect(page.locator(".context-display--home")).toBeVisible();
  await expect(page.locator(".v2-topbar .logout-button-v2")).toBeVisible();
  await expectNoMojibake(page);

  const appBox = await page.locator(".v2-app").boundingBox();
  const topbarBox = await page.locator(".v2-topbar").boundingBox();
  const logoutBox = await page.locator(".v2-topbar .logout-button-v2").boundingBox();
  expect(appBox).not.toBeNull();
  expect(topbarBox).not.toBeNull();
  expect(logoutBox).not.toBeNull();
  expect(appBox!.y).toBeLessThanOrEqual(6);
  expect(topbarBox!.y).toBeGreaterThanOrEqual(appBox!.y - 1);
  expect(logoutBox!.x + logoutBox!.width).toBeLessThanOrEqual(appBox!.x + appBox!.width - 6);
  expect(logoutBox!.y + logoutBox!.height).toBeLessThanOrEqual(topbarBox!.y + topbarBox!.height + 1);

  await expect(page.locator(".action-card-v2")).toHaveCount(0);
  await expect(page.locator(".daily-quests-v2__list button")).toHaveCount(4);
  await expect(page.locator(".weekly-plan-v2__week button")).toHaveCount(7);
  await expect(page.getByText("Next features")).toHaveCount(0);

  const calendar = await page.locator(".calendar-rac--habit").boundingBox();
  const side = await page.locator(".habit-calendar-v2__side").boundingBox();
  expect(calendar).not.toBeNull();
  expect(side).not.toBeNull();
  expect(side!.x).toBeGreaterThan(calendar!.x + calendar!.width - 2);
  expect(Math.abs(side!.y - calendar!.y)).toBeLessThan(12);

  const streakColor = await page.locator(".habit-calendar-v2__side strong").evaluate((node) => getComputedStyle(node).color);
  const channels = streakColor.match(/\d+/g)?.map(Number) || [];
  expect(channels[1]).toBeGreaterThan(channels[0]);
});

test("desktop ribbon right arrow stays inside the app frame", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop layout assertion");
  await page.goto("/app/?view=home");
  const app = await page.locator(".v2-app").boundingBox();
  const ribbon = await page.locator(".function-ribbon").boundingBox();
  const leftArrow = await page.locator(".function-ribbon-shell .ribbon-morph-arrow").first().boundingBox();
  const rightArrow = await page.locator(".function-ribbon-shell .ribbon-morph-arrow").last().boundingBox();
  expect(app).not.toBeNull();
  expect(ribbon).not.toBeNull();
  expect(leftArrow).not.toBeNull();
  expect(rightArrow).not.toBeNull();
  expect(leftArrow!.width).toBeLessThanOrEqual(60);
  expect(rightArrow!.width).toBeLessThanOrEqual(60);
  expect(leftArrow!.height).toBeLessThanOrEqual(48);
  expect(rightArrow!.height).toBeLessThanOrEqual(48);
  expect(leftArrow!.x + leftArrow!.width).toBeLessThanOrEqual(ribbon!.x);
  expect(rightArrow!.x).toBeGreaterThanOrEqual(ribbon!.x + ribbon!.width);
  expect(rightArrow!.x + rightArrow!.width).toBeLessThanOrEqual(app!.x + app!.width - 6);

  for (const view of ["pronunciation", "phrasebook", "offline", "roleplay"]) {
    const background = await page.locator(`.function-chip[data-view="${view}"] .function-chip__art`).evaluate((node) => getComputedStyle(node).backgroundImage);
    expect(background).toContain(`icon-${view}-light.png?v=flux2-20260526-pronunciation-phrases-offline-roleplay`);
  }
});

test("browser back navigates inside app history instead of leaving the app", async ({ page, isMobile }) => {
  await page.goto("/app/?view=home");
  await expect(page.locator(".context-display--home")).toBeVisible();
  const practiceNav = isMobile ? page.locator('.mobile-bottom-nav-v2 [data-view="practice"]') : page.locator('.function-chip[data-view="practice"]');
  await practiceNav.click();
  await expect(page.locator(".context-display--practice")).toBeVisible();
  await expect(page).toHaveURL(/view=practice/);
  await page.goBack();
  await expect(page.locator(".context-display--home")).toBeVisible();
  await expect(page).toHaveURL(/view=home/);
});

test("unknown app route renders localized 404 without leaving shell", async ({ page }) => {
  await page.goto("/app/?view=unknown-section");
  await expect(page.locator(".context-display--not-found")).toBeVisible();
  await expect(page.getByText("Страница не найдена")).toBeVisible();
  await page.getByRole("button", { name: /На главный экран|Go home/ }).click();
  await expect(page.locator(".context-display--home")).toBeVisible();
});

test("daily bonus button stays locked when last claim is within 24 hours", async ({ page }) => {
  const today = new Date().toISOString().slice(0, 10);
  let dailyClaimCalled = false;
  await page.route("**/api/session", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        ...sessionPayload,
        user: {
          ...sessionPayload.user,
          daily_bonus_claims: [today],
          daily_bonus_last_claimed_at: new Date().toISOString(),
        },
      }),
    }),
  );
  await page.route("**/api/daily/claim", (route) => {
    dailyClaimCalled = true;
    return route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, claimed: false }) });
  });
  await page.goto("/app/?view=home");
  const claimButton = page.locator(".habit-calendar-v2__side button").last();
  await expect(claimButton).toBeDisabled();
  await expect(claimButton).toContainText(/Бонус получен|Bonus claimed/);
  await claimButton.click({ force: true });
  expect(dailyClaimCalled).toBe(false);
});

test("daily bonus is claimed automatically when today's available goal is complete", async ({ page }) => {
  let dailyClaimCalls = 0;
  await page.route("**/api/daily/claim", (route) => {
    dailyClaimCalls += 1;
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ ok: true, claimed: true, claimed_at: new Date().toISOString(), xp: 25, streak: 1, user: { ...sessionPayload.user, xp: 660 } }),
    });
  });

  await page.goto("/app/?view=home");
  await expect(page.locator(".context-display--home")).toBeVisible();
  await expect.poll(() => dailyClaimCalls).toBe(1);
  await expect(page.locator(".calendar-rac-cell--today")).toHaveAttribute("data-complete", "true");
  await expect(page.locator(".habit-calendar-v2__side button").last()).toContainText(/Бонус получен|Bonus claimed/);
});

test("learning actions show a visible XP gain after profile XP increases", async ({ page }) => {
  const today = new Date().toISOString().slice(0, 10);
  const claimedSession = {
    ...sessionPayload,
    user: { ...sessionPayload.user, daily_bonus_claims: [today], daily_bonus_last_claimed_at: new Date().toISOString() },
  };
  let practiceCalled = false;
  await page.route("**/api/session", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(claimedSession) }));
  await page.route("**/api/practice", (route) => {
    practiceCalled = true;
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        correction: "Could you repeat that, please?",
        explanation: "Short polite request.",
        xp: 5,
        user: { ...claimedSession.user, xp: Number(sessionPayload.user.xp) + 5, xp_current: Number(sessionPayload.user.xp_current) + 5 },
      }),
    });
  });

  await page.goto("/app/?view=practice");
  await expect(page.locator(".context-display--practice")).toBeVisible();
  await page.locator(".composer-panel-v2 textarea").fill("Can you repeat?");
  await page.locator(".composer-panel-v2 textarea").press("Enter");
  await expect.poll(() => practiceCalled).toBe(true);
  await expect(page.locator(".xp-gain-pop-v2")).toContainText("+5 XP");
});

test("payment history shows only confirmed payments", async ({ page }) => {
  await page.goto("/app/?view=premium");
  const history = page.locator(".payment-history-v2");
  await expect(history).toBeVisible();
  await expect(history.getByText("300 RUB")).toBeVisible();
  await expect(history.getByText("paid")).toBeVisible();
  await expect(history.getByText("1.5 TON")).toHaveCount(0);
  await expect(history.getByText("590 RUB")).toHaveCount(0);
  await expect(history.getByText("pending")).toHaveCount(0);
  await expect(history.getByText("created")).toHaveCount(0);
});

test("mobile payment requisites modal stays above bottom navigation and scrolls", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile payment layout assertion");
  await page.goto("/app/?view=premium");
  await page.locator(".plan-card-v2:not(.is-free) button").first().click();
  await expect(page.locator(".v2-payment-modal")).toBeVisible();
  await page.getByRole("button", { name: /TON|USDT/ }).click();
  await expect(page.locator(".payment-invoice-v2")).toBeVisible();
  await expect(page.getByText("UQCDkENqCLcFLvAzPHX8LxfK6wdlPEmpFQ5DQ5UhG9BIZQ3i")).toBeVisible();

  const backdropZ = await page.locator(".modal-backdrop-v2").evaluate((node) => Number(getComputedStyle(node).zIndex));
  const navZ = await page.locator(".mobile-bottom-nav-v2").evaluate((node) => Number(getComputedStyle(node).zIndex));
  expect(backdropZ).toBeGreaterThan(navZ);

  const metrics = await page.locator(".modal-backdrop-v2").evaluate((node) => ({
    clientHeight: node.clientHeight,
    scrollHeight: node.scrollHeight,
    overflowY: getComputedStyle(node).overflowY,
  }));
  expect(metrics.overflowY).toMatch(/auto|scroll/);
  expect(metrics.scrollHeight).toBeGreaterThanOrEqual(metrics.clientHeight);
});

test("mobile onboarding dialog scrolls and keeps actions reachable", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.setViewportSize({ width: 390, height: 640 });
  await page.addInitScript(() => {
    localStorage.setItem("poliglot-test-show-onboarding", "1");
    localStorage.removeItem("poliglot-onboarding-v2:demor22");
  });
  await page.goto("/app/?view=home");
  const dialog = page.locator(".onboarding-dialog-v2");
  await expect(dialog).toBeVisible();
  await expect(dialog.locator('select').first()).toContainText("C2");
  const viewport = page.viewportSize();
  const dialogBox = await dialog.boundingBox();
  const bodyMetrics = await dialog.locator(".onboarding-dialog-v2__body").evaluate((node) => ({
    clientHeight: node.clientHeight,
    scrollHeight: node.scrollHeight,
  }));
  expect(viewport).not.toBeNull();
  expect(dialogBox).not.toBeNull();
  expect(dialogBox!.height).toBeLessThanOrEqual(viewport!.height - 8);
  expect(bodyMetrics.scrollHeight).toBeGreaterThan(bodyMetrics.clientHeight);

  await dialog.locator(".onboarding-dialog-v2__body").evaluate((node) => node.scrollTo(0, node.scrollHeight));
  const submit = dialog.locator('.onboarding-dialog-v2__footer button[type="submit"]');
  await expect(submit).toBeVisible();
  const submitBox = await submit.boundingBox();
  expect(submitBox).not.toBeNull();
  expect(submitBox!.y + submitBox!.height).toBeLessThanOrEqual(viewport!.height);
  await submit.click();
  await expect(dialog).toHaveCount(0);
});

test("offline deck paginates by 10 and exports the full deck", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop download assertion");
  await page.goto("/app/?view=offline");
  await expect(page.getByText("Оффлайн-карточки для повторения")).toBeVisible();
  await page.getByRole("button", { name: /Обновить колоду/ }).click();
  await expect(page.locator(".offline-deck-grid-v2 article")).toHaveCount(10);
  await expect(page.getByText(/1\/4/)).toBeVisible();
  await expect(page.getByText("travel-word-12")).toHaveCount(0);

  await page.locator(".offline-pagination-v2 button").nth(1).click();
  await expect(page.getByText("travel-word-12")).toBeVisible();

  const downloadPromise = page.waitForEvent("download");
  await page.getByRole("button", { name: "TXT" }).click();
  const download = await downloadPromise;
  const path = await download.path();
  expect(path).toBeTruthy();
  const content = readFileSync(path!, "utf8");
  expect(content).toContain("travel-word-1");
  expect(content).toContain("travel-word-12");
  expect(content).toContain("I have a reservation under the name Ivan Petrov.");
});

test("offline notes group paginates by 10 like vocabulary", async ({ page }) => {
  await page.addInitScript(() => {
    const notes = Array.from({ length: 13 }, (_, index) => ({
      id: `offline-note-${index + 1}`,
      title: `Offline note phrase ${index + 1}`,
      answer: index === 0 ? "перевод заметки 1" : `Offline note translation ${index + 1}`,
      example: "practice",
      source: "phrasebook",
    }));
    const mixedDeck = [
      { id: "offline-word-1", title: "Downloaded word", answer: "слово", example: "vocabulary", source: "vocabulary" },
      ...notes,
      { id: "offline-mistake-1", title: "Downloaded mistake", answer: "correction", example: "mistake", source: "mistakes" },
    ];
    localStorage.setItem("poliglot-offline-deck-v2", JSON.stringify(mixedDeck));
  });

  await page.goto("/app/?view=offline");
  await expect(page.locator(".context-display--offline")).toBeVisible();
  await page.locator(".offline-group-tabs-v2 button").filter({ hasText: ru("phrasebook", "Заметки") }).click();
  await expect(page.locator(".offline-deck-grid-v2 article")).toHaveCount(10);
  await expect(page.getByText(/1\/2/)).toBeVisible();
  await expect(page.getByText("Offline note phrase 11")).toHaveCount(0);

  await page.locator(".offline-pagination-v2 button").nth(1).click();
  await expect(page.locator(".offline-deck-grid-v2 article")).toHaveCount(3);
  await expect(page.getByText(/2\/2/)).toBeVisible();
  await expect(page.getByText("Offline note phrase 11")).toBeVisible();
  await expect(page.getByText(/^Offline note phrase 1$/)).toHaveCount(0);

  const downloadPromise = page.waitForEvent("download");
  await page.getByRole("button", { name: "TXT" }).click();
  const download = await downloadPromise;
  const path = await download.path();
  expect(path).toBeTruthy();
  const bytes = readFileSync(path!);
  expect(bytes[0]).toBe(0xef);
  expect(bytes[1]).toBe(0xbb);
  expect(bytes[2]).toBe(0xbf);
  const content = bytes.toString("utf8");
  expect(content).toContain("Offline note phrase 1");
  expect(content).toContain("перевод заметки 1");
  expect(content).toContain("Offline note phrase 13");
  expect(content).not.toContain("Downloaded word");
  expect(content).not.toContain("Downloaded mistake");
  expect(content).not.toMatch(mojibakePattern);

  await page.locator(".offline-card-delete-v2").first().click();
  await expect(page.getByText("Offline note phrase 11")).toHaveCount(0);
  await expect(page.locator(".offline-deck-grid-v2 article")).toHaveCount(10);
  await expect(page.getByText(/1\/2/)).toBeVisible();
  await expect
    .poll(() =>
      page.evaluate(() => JSON.parse(localStorage.getItem("poliglot-offline-deck-v2") || "[]").some((item: { id: string }) => item.id === "offline-note-11")),
    )
    .toBe(false);
});

test("referral block shows invited users, level 3 status, earnings and pagination", async ({ page }) => {
  await page.goto("/app/?view=referral");
  await expect(page.locator(".context-display--referral")).toBeVisible();
  await expect(page.getByText("Кто перешёл по ссылке")).toBeVisible();
  await expect(page.locator(".referral-invitees-v2__rows article")).toHaveCount(10);
  await expect(page.locator(".referral-invitees-v2__rows article").first().locator("strong")).toContainText("ref-user-1");
  await expect(page.getByText("ref-user-11")).toHaveCount(0);
  await expect(page.getByText("150 RUB").first()).toBeVisible();
  await page.locator(".referral-pagination-v2 button").nth(1).click();
  await expect(page.getByText("ref-user-11")).toBeVisible();
});

test("linked Telegram account hides Send code in settings", async ({ page }) => {
  await page.goto("/app/?view=settings");
  await expect(page.locator(".context-display--settings")).toBeVisible();
  await expect(page.getByText("Telegram ID 185156683")).toBeVisible();
  await expect(page.getByRole("button", { name: /Send code|Отправить код/ })).toHaveCount(0);
});

test("pronunciation map is compact and removes repeated advice", async ({ page }) => {
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      plan: "premium",
      premium: true,
      premium_until: "2026-07-12T12:00:00Z",
      voice_limit: 20,
    },
  };
  await page.goto("/app/?view=shadowing");
  await expect(page.locator(".context-display--shadowing")).toBeVisible();
  await page.locator("textarea").fill("Could you repeat that please");
  await page.locator("textarea").press("Enter");
  await page.locator('[data-view="pronunciation"]:visible').first().click();
  await expect(page.locator(".context-display--pronunciation")).toBeVisible();
  await expect(page.getByRole("heading", { name: "Карта произношения" })).toBeVisible();
  await expect(page.locator(".heatmap-token-v2")).toHaveCount(2);
  await expect(page.locator(".heatmap-token-v2 em", { hasText: "Уверенность" })).toHaveCount(1);
  const stored = await page.evaluate(() => localStorage.getItem("poliglot-pronunciation-v2:demor22") || "");
  expect(stored).toContain("coffee");
  await page.reload();
  await expect(page.locator(".context-display--pronunciation")).toBeVisible();
  await expect(page.locator(".heatmap-token-v2")).toHaveCount(2);
  await expect(page.locator(".pronunciation-history-v2__rows > div")).toHaveCount(1);
});

test("pronunciation is a standalone sample, voice input, result and next sample flow", async ({ page }) => {
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      plan: "premium",
      premium: true,
      premium_until: "2026-07-12T12:00:00Z",
      voice_limit: 20,
    },
  };
  await page.goto("/app/?view=pronunciation");
  await expect(page.locator(".context-display--pronunciation")).toBeVisible();
  await expect(page.locator(".context-display--shadowing")).toHaveCount(0);
  await expect(page.locator(".pronunciation-workbench-v2")).toBeVisible();
  await expect(page.locator(".pronunciation-target-primary-v2 .audio-wave-button-v2")).toBeVisible();
  await expect(page.locator(".pronunciation-tools-v2 .file-controls-v2")).toBeVisible();
  await expect(page.locator(".pronunciation-result-window-v2")).toHaveCount(0);

  await page.locator('.pronunciation-tools-v2 input[type="file"]').setInputFiles({
    name: "pronunciation.webm",
    mimeType: "audio/webm",
    buffer: Buffer.from("test-audio"),
  });
  await page.locator(".pronunciation-tools-v2").getByRole("button", { name: /Record and check|Записать|Проверить/ }).click();
  await expect(page.locator(".pronunciation-target-primary-v2")).toBeVisible();
  await expect(page.locator(".pronunciation-tools-v2 .pronunciation-report-v2")).toContainText("67/100");
  await expect(page.locator(".pronunciation-tools-v2")).toContainText("Stress is late");
  await expect(page.locator(".pronunciation-result-window-v2")).toHaveCount(0);

  await page.locator(".pronunciation-tools-v2").getByRole("button", { name: /Next phrase|Следующая|Следующий/ }).click();
  await expect(page.locator(".pronunciation-workbench-v2")).toBeVisible();
  await expect(page.locator(".pronunciation-tools-v2 .pronunciation-report-v2")).toHaveCount(0);
});

test("regression: lesson tab keeps one active task until the learner submits", async ({ page }) => {
  let lessonStarts = 0;
  await page.unroute("**/api/lesson/start").catch(() => undefined);
  await page.route("**/api/lesson/start", (route) => {
    lessonStarts += 1;
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        prompt: `Lesson task ${lessonStarts}`,
        example: "I have a reservation under the name Ivan Petrov.",
        question_audio_text: "Could you please show me your passport?",
      }),
    });
  });

  await page.goto("/app/?view=lesson");
  await page.locator(".lesson-new-button-v2").click();
  await expect(page.locator(".context-display--lesson")).toContainText("Lesson task 1");
  await page.locator('[data-view="words"]:visible').first().click();
  await expect(page.locator(".context-display--words")).toBeVisible();
  await page.locator('[data-view="lesson"]:visible').first().click();
  await expect(page.locator(".context-display--lesson")).toContainText("Lesson task 1");
  await expect(page.locator(".context-display--lesson")).not.toContainText("Lesson task 2");
  expect(lessonStarts).toBe(1);

  await page.locator(".context-display--lesson textarea").fill("I have a reservation under the name Ivan Petrov.");
  await page.locator(".context-display--lesson").getByRole("button", { name: /Отправить|Send/ }).click();
  await expect(page.locator(".context-display--lesson")).toContainText("Use under the name for reservations.");
  await expect(page.locator(".context-display--lesson textarea")).toHaveCount(0);
  const panelNextButton = page.locator(".context-display--lesson .lesson-panel-next-v2");
  await expect(panelNextButton).toBeVisible();
  await expect(panelNextButton).toContainText("Следующий урок");
  await expect(page.locator(".lesson-new-button-v2")).toContainText("Следующий урок");
  await page.locator(".lesson-new-button-v2").click();
  await expect(page.locator(".context-display--lesson")).toContainText("Lesson task 2");
  expect(lessonStarts).toBe(2);
});

test("regression: mobile composer and recording controls expose clear labels", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      plan: "premium",
      premium: true,
      premium_until: "2026-07-12T12:00:00Z",
      voice_limit: 20,
    },
  };
  await page.goto("/app/?view=lesson");
  await page.locator(".lesson-empty-v2 button").click();
  await expect(page.getByText("Say that you have a reservation.")).toBeVisible();
  const sendButton = page.locator(".context-display--lesson .composer-submit-v2");
  await expect(sendButton).toBeVisible();
  await expect(sendButton).toContainText("Отправить");
  const sendFits = await sendButton.evaluate((node) => {
    const element = node as HTMLElement;
    return element.scrollWidth <= element.clientWidth + 1 && element.scrollHeight <= element.clientHeight + 1;
  });
  expect(sendFits).toBe(true);

  await page.goto("/app/?view=pronunciation");
  const checkButton = page.locator(".pronunciation-tools-v2").getByRole("button", { name: /Проверить|Check/ }).first();
  await expect(checkButton).toBeVisible();
  await expect(checkButton).toContainText("Проверить");
});

test("regression: mobile XP gain notice stays compact", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=home");
  const wordsNav = page.locator('[data-view="words"]:visible').first();
  await expect(wordsNav).toBeVisible();
  await wordsNav.click();
  await expect(page.locator(".context-display--words")).toBeVisible();
  await expect(page.locator(".choice-grid-v2 button")).toHaveCount(4);
  await page.locator(".choice-grid-v2 button", { hasText: "apple" }).click();
  const notice = page.locator(".xp-gain-pop-v2");
  await expect(notice).toBeVisible();
  const box = await notice.boundingBox();
  expect(box).not.toBeNull();
  expect(box!.width).toBeLessThanOrEqual(360);
  expect(box!.height).toBeLessThanOrEqual(112);
});

test("desktop pronunciation workspace makes the target phrase the primary panel", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop layout assertion");
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      plan: "premium",
      premium: true,
      premium_until: "2026-07-12T12:00:00Z",
      voice_limit: 20,
    },
  };
  await page.goto("/app/?view=pronunciation");
  await expect(page.locator(".context-display--pronunciation")).toBeVisible();
  await expect(page.locator(".pronunciation-workbench-v2")).toBeVisible();

  const targetPanel = page.locator(".pronunciation-target-primary-v2");
  const toolsPanel = page.locator(".pronunciation-tools-v2");
  await expect(targetPanel).toBeVisible();
  await expect(toolsPanel).toBeVisible();
  await expect(targetPanel.locator(".audio-wave-button-v2")).toBeVisible();
  await expect(toolsPanel.locator(".file-controls-v2")).toBeVisible();

  const targetBox = await targetPanel.boundingBox();
  const toolsBox = await toolsPanel.boundingBox();
  expect(targetBox).not.toBeNull();
  expect(toolsBox).not.toBeNull();
  expect(targetBox!.width).toBeGreaterThan(toolsBox!.width);
  expect(targetBox!.height).toBeGreaterThan(240);
  expect(toolsBox!.x).toBeGreaterThan(targetBox!.x + targetBox!.width - 4);

  await toolsPanel.locator('input[type="file"]').setInputFiles({
    name: "pronunciation.webm",
    mimeType: "audio/webm",
    buffer: Buffer.from("test-audio"),
  });
  await toolsPanel.getByRole("button", { name: /Record and check|Записать|Проверить/ }).click();
  await expect(targetPanel).toBeVisible();
  await expect(toolsPanel.locator(".pronunciation-report-v2")).toContainText("67/100");
  await expect(page.locator(".pronunciation-result-window-v2")).toHaveCount(0);
});

test("regression: desktop function ribbon arrows stay inside the menu frame", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop layout assertion");
  await page.goto("/app/?view=pronunciation");
  const shell = page.locator(".function-ribbon-shell");
  const ribbon = page.locator(".function-ribbon");
  const leftArrow = shell.locator(".ribbon-morph-arrow").first();
  const rightArrow = shell.locator(".ribbon-morph-arrow").last();
  await expect(shell).toBeVisible();

  const metrics = await shell.evaluate((node) => {
    const shellBox = node.getBoundingClientRect();
    const ribbonBox = node.querySelector(".function-ribbon")?.getBoundingClientRect();
    const arrows = Array.from(node.querySelectorAll(".ribbon-morph-arrow")).map((arrow) => arrow.getBoundingClientRect());
    const leftArrowCenter = { x: (arrows[0]?.left || 0) + (arrows[0]?.width || 0) / 2, y: (arrows[0]?.top || 0) + (arrows[0]?.height || 0) / 2 };
    const rightArrowCenter = { x: (arrows[1]?.left || 0) + (arrows[1]?.width || 0) / 2, y: (arrows[1]?.top || 0) + (arrows[1]?.height || 0) / 2 };
    const leftHit = document.elementFromPoint(leftArrowCenter.x, leftArrowCenter.y);
    const rightHit = document.elementFromPoint(rightArrowCenter.x, rightArrowCenter.y);
    return {
      shellLeft: shellBox.left,
      shellRight: shellBox.right,
      ribbonLeft: ribbonBox?.left || 0,
      ribbonRight: ribbonBox?.right || 0,
      leftArrowLeft: arrows[0]?.left || 0,
      leftArrowRight: arrows[0]?.right || 0,
      rightArrowLeft: arrows[1]?.left || 0,
      rightArrowRight: arrows[1]?.right || 0,
      leftArrowGap: (ribbonBox?.left || 0) - (arrows[0]?.right || 0),
      rightArrowGap: (arrows[1]?.left || 0) - (ribbonBox?.right || 0),
      leftHitIsChip: Boolean(leftHit?.closest(".function-chip")),
      rightHitIsChip: Boolean(rightHit?.closest(".function-chip")),
      leftHitIsArrow: Boolean(leftHit?.closest(".ribbon-morph-arrow")),
      rightHitIsArrow: Boolean(rightHit?.closest(".ribbon-morph-arrow")),
    };
  });

  await expect(ribbon).toBeVisible();
  await expect(leftArrow).toBeVisible();
  await expect(rightArrow).toBeVisible();
  expect(metrics.leftArrowLeft).toBeGreaterThanOrEqual(metrics.shellLeft + 8);
  expect(metrics.rightArrowRight).toBeLessThanOrEqual(metrics.shellRight - 8);
  expect(metrics.leftArrowGap).toBeGreaterThanOrEqual(10);
  expect(metrics.rightArrowGap).toBeGreaterThanOrEqual(10);
  expect(metrics.leftHitIsChip).toBe(false);
  expect(metrics.rightHitIsChip).toBe(false);
  expect(metrics.leftHitIsArrow).toBe(true);
  expect(metrics.rightHitIsArrow).toBe(true);
});

test("regression: mobile pronunciation blocks keep vertical order after checking audio", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      plan: "premium",
      premium: true,
      premium_until: "2026-07-12T12:00:00Z",
      voice_limit: 20,
    },
  };
  await page.goto("/app/?view=pronunciation");
  await expect(page.locator(".context-display--pronunciation")).toBeVisible();
  await page.locator(".pronunciation-tools-v2 input[type='file']").setInputFiles({
    name: "pronunciation.webm",
    mimeType: "audio/webm",
    buffer: Buffer.from("pronunciation-test"),
  });
  await page.getByRole("button", { name: /Check|Проверить/ }).click();
  await expect(page.locator(".pronunciation-report-v2")).toContainText("67/100");
  const boxes = await page.locator(".pronunciation-workbench-v2").evaluate((node) => {
    const rect = (element: Element | null) => {
      if (!element) return null;
      const box = element.getBoundingClientRect();
      return { top: box.top, bottom: box.bottom, left: box.left, right: box.right, width: box.width, height: box.height };
    };
    return {
      target: rect(node.querySelector(".pronunciation-target-primary-v2")),
      tools: rect(node.querySelector(".pronunciation-tools-v2")),
      controls: rect(node.querySelector(".pronunciation-tool-card-v2")),
      report: rect(node.querySelector(".pronunciation-report-window-v2")),
      workbench: rect(node),
    };
  });
  expect(boxes.target).not.toBeNull();
  expect(boxes.tools).not.toBeNull();
  expect(boxes.controls).not.toBeNull();
  expect(boxes.report).not.toBeNull();
  expect(boxes.tools!.top).toBeGreaterThanOrEqual(boxes.target!.bottom - 1);
  expect(boxes.report!.top).toBeGreaterThanOrEqual(boxes.controls!.bottom - 1);
  expect(boxes.workbench!.bottom).toBeGreaterThanOrEqual(boxes.report!.bottom - 1);
});

test("regression: desktop tools selector buttons stay compact", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop layout assertion");
  await page.goto("/app/?view=tools");
  await expect(page.locator(".context-display--tools")).toBeVisible();
  await expect(page.locator(".tools-ai-router-v2")).toHaveText("AI Router");
  await expect(page.locator(".tools-ai-router-v2")).not.toContainText("Раздел");
  const buttons = await page.locator(".tool-switch-v2 .tool-button-v2").evaluateAll((nodes) =>
    nodes.map((node) => {
      const box = (node as HTMLElement).getBoundingClientRect();
      return { width: box.width, height: box.height, scrollWidth: (node as HTMLElement).scrollWidth };
    }),
  );
  expect(buttons).toHaveLength(4);
  for (const button of buttons) {
    expect(button.height).toBeLessThanOrEqual(40);
    expect(button.scrollWidth).toBeLessThanOrEqual(button.width + 1);
  }
});

test("roleplay scenario cards show concrete roles instead of the generic roleplay label", async ({ page }) => {
  await page.goto("/app/?view=roleplay");
  await expect(page.locator(".context-display--roleplay")).toBeVisible();
  await expect(page.locator(".roleplay-grid-v2 button").nth(0)).toContainText("Ресторан");
  await expect(page.locator(".roleplay-grid-v2 button").nth(1)).toContainText("Рабочий");
  await expect(page.locator(".roleplay-grid-v2 button").nth(2)).toContainText("Путешествие");
  const genericCards = await page.locator(".roleplay-grid-v2 button").evaluateAll((cards) =>
    cards.filter((card) => {
      const parts = Array.from(card.querySelectorAll("strong, span")).map((node) => (node.textContent || "").trim());
      return parts.filter(Boolean).every((part) => part === "Ролевая");
    }).length,
  );
  expect(genericCards).toBe(0);
});

test("AI Tutor failed start stops loading loop and shows retry", async ({ page }) => {
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      plan: "premium",
      premium: true,
      premium_until: "2026-07-12T12:00:00Z",
      lesson_limit: 50,
      practice_limit: 200,
      voice_limit: 20,
    },
  };
  let tutorStartCalls = 0;
  await page.unroute("**/api/ai-tutor/start");
  await page.route("**/api/ai-tutor/start", (route) => {
    tutorStartCalls += 1;
    return route.fulfill({
      status: 500,
      contentType: "application/json",
      body: JSON.stringify({ error: "No tutor lesson available" }),
    });
  });

  await page.goto("/app/?view=tutor");
  await expect(page.locator(".tutor-loading-panel")).toContainText(/No tutor lesson available|Не удалось/);
  await expect(page.locator(".tutor-loading-panel").getByRole("button", { name: /Retry|Повторить/ })).toBeVisible();
  expect(tutorStartCalls).toBe(1);
});

test("learn words shows varied wrong answer options across rounds", async ({ page }) => {
  await page.goto("/app/?view=words");
  await expect(page.locator(".context-display--words")).toBeVisible();
  await expect(page.locator(".choice-grid-v2 button")).toHaveCount(4);
  const firstOptions = await page.locator(".choice-grid-v2 button").allInnerTexts();
  const firstCorrect = firstOptions.includes("apple") ? "apple" : "train";
  await page.locator(".choice-grid-v2 button", { hasText: firstCorrect }).click();
  await expect(page.locator(".trainer-result-v2")).toBeVisible();
  const resultBox = await page.locator(".trainer-result-v2").boundingBox();
  const headingBox = await page.locator(".trainer-display h2").boundingBox();
  expect(resultBox).not.toBeNull();
  expect(headingBox).not.toBeNull();
  expect(resultBox!.y).toBeGreaterThanOrEqual(headingBox!.y + headingBox!.height - 2);
  expect(resultBox!.y).toBeLessThan(headingBox!.y + headingBox!.height + 24);
  await page.locator(".trainer-result-v2 > button").click();
  await expect(page.locator(".context-display--words")).toBeVisible();
  await expect(page.locator(".choice-grid-v2 button")).toHaveCount(4);
  const secondOptions = await page.locator(".choice-grid-v2 button").allInnerTexts();
  const firstWrong = firstOptions.filter((text) => !["apple", "train"].includes(text)).sort();
  const secondWrong = secondOptions.filter((text) => !["apple", "train"].includes(text)).sort();
  expect(firstWrong).not.toEqual(secondWrong);
});

test("learn words result shows target word and saves word pair to notes", async ({ page, isMobile }) => {
  let savedItem: Record<string, unknown> | null = null;
  await page.route("**/api/phrasebook**", async (route) => {
    const request = route.request();
    if (request.method() === "POST") {
      savedItem = request.postDataJSON();
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, items: [savedItem, ...phrasebookSeed] }) });
      return;
    }
    await route.fallback();
  });

  await page.goto("/app/?view=words");
  await expect(page.locator(".context-display--words")).toBeVisible();
  await expect(page.locator(".choice-grid-v2 button")).toHaveCount(4);
  const options = await page.locator(".choice-grid-v2 button").allInnerTexts();
  const targetWord = options.includes("apple") ? "apple" : "train";
  const targetTranslation = targetWord === "apple" ? "\u044f\u0431\u043b\u043e\u043a\u043e" : "\u043f\u043e\u0435\u0437\u0434";
  await page.locator(".choice-grid-v2 button", { hasText: targetWord }).click();

  const result = page.locator(".trainer-result-v2");
  await expect(result).toBeVisible();
  await expect(result.locator(".trainer-correct-word-v2")).toContainText(targetWord);
  await expect(result).toContainText(targetTranslation);

  const nextButton = result.getByRole("button", { name: new RegExp(ru("next_word", "Next")) });
  const saveChip = result.locator(".trainer-word-save-v2 .phrase-quick-save-v2__chips button");
  await expect(nextButton).toBeVisible();
  await expect(saveChip).toContainText(targetWord);
  await expect(saveChip).toContainText(targetTranslation);

  const nextBox = await nextButton.boundingBox();
  const saveBox = await saveChip.boundingBox();
  expect(nextBox).not.toBeNull();
  expect(saveBox).not.toBeNull();
  expect(saveBox!.y).toBeGreaterThan(nextBox!.y);
  if (isMobile) {
    const resultBox = await result.boundingBox();
    const navBox = await page.locator(".mobile-bottom-nav-v2").boundingBox();
    expect(resultBox).not.toBeNull();
    expect(navBox).not.toBeNull();
    expect(resultBox!.y + resultBox!.height).toBeLessThanOrEqual(navBox!.y - 4);
  }

  await saveChip.click();
  await expect.poll(() => String(savedItem?.phrase || "")).toContain(targetWord);
  expect(String(savedItem?.translation || savedItem?.note || "")).toContain(targetTranslation);
});

test("spelling result shows the correct target-language word", async ({ page }) => {
  await page.goto("/app/?view=spelling");
  await expect(page.locator(".context-display--spelling")).toBeVisible();
  await page.locator(".trainer-display input").fill("clear");
  await page.locator(".trainer-display button", { hasText: ru("check", "Check") }).click();
  await expect(page.locator(".trainer-correct-word-v2")).toContainText("clear");
  await expect(page.locator(".trainer-result-v2")).toContainText("The singer has a clear voice.");
});

test("spelling wrong attempt does not reveal the correct answer", async ({ page }) => {
  await page.goto("/app/?view=spelling");
  await expect(page.locator(".context-display--spelling")).toBeVisible();
  await page.locator(".trainer-display input").fill("cler");
  await page.locator(".trainer-display button", { hasText: ru("check", "Check") }).click();
  const result = page.locator(".trainer-result-v2");
  await expect(result).toBeVisible();
  await expect(result.locator(".trainer-correct-word-v2")).toHaveCount(0);
  await expect(result).not.toContainText("clear");
  await expect(result).not.toContainText("The singer has a clear voice.");
});

test("mobile spelling result stays readable above bottom menu", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=spelling");
  await expect(page.locator(".context-display--spelling")).toBeVisible();
  await page.locator(".trainer-display input").fill("clear");
  await page.locator(".trainer-display button", { hasText: ru("check", "Check") }).click();
  await expect(page.locator(".trainer-correct-word-v2")).toContainText("clear");
  await page.locator(".context-display").evaluate((node) => node.scrollTo(0, node.scrollHeight));
  const resultBox = await page.locator(".trainer-result-v2").boundingBox();
  const navBox = await page.locator(".mobile-bottom-nav-v2").boundingBox();
  expect(resultBox).not.toBeNull();
  expect(navBox).not.toBeNull();
  expect(resultBox!.y + resultBox!.height).toBeLessThanOrEqual(navBox!.y - 4);
});

test("mobile header controls and editable nav rail reorder work", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=home");
  await page.evaluate(() => localStorage.removeItem("poliglot-mobile-nav-rail-v2:demor22"));
  await page.reload();
  await expect(page.locator(".v2-topbar")).toBeHidden();
  await expect(page.locator(".function-ribbon-shell")).toBeHidden();
  await expect(page.locator(".mobile-quick-controls-v2 .mobile-report-button-v2")).toBeVisible();
  await expect(page.locator(".mobile-quick-controls-v2 .mobile-language-select-v2")).toBeVisible();
  await expect(page.locator(".mobile-quick-controls-v2 .mobile-theme-toggle-v2")).toBeVisible();
  await expect(page.locator(".mobile-quick-controls-v2 .logout-button-v2")).toBeVisible();
  await expect(page.locator(".mobile-quick-controls-v2").getByText("Выйти")).toHaveCount(0);

  const quick = await page.locator(".mobile-quick-controls-v2").boundingBox();
  const frame = await page.locator(".context-display").boundingBox();
  expect(quick).not.toBeNull();
  expect(frame).not.toBeNull();
  expect(quick!.y).toBeGreaterThan(frame!.y + 8);

  await expect(page.locator('.mobile-bottom-nav-v2 [data-view="more"]')).toHaveCount(0);
  await expect(page.locator(".mobile-more-sheet-v2")).toHaveCount(0);
  await expect(page.locator('.mobile-bottom-nav-v2 [data-view="settings"]')).toBeVisible();
  const railMetrics = await page.locator(".mobile-bottom-nav-v2").evaluate((node) => ({
    clientWidth: node.clientWidth,
    scrollWidth: node.scrollWidth,
    overflowX: getComputedStyle(node).overflowX,
  }));
  expect(railMetrics.scrollWidth).toBeGreaterThan(railMetrics.clientWidth);
  expect(railMetrics.overflowX).toMatch(/auto|scroll/);
  const pronunciationBg = await page.locator('.mobile-bottom-nav-v2 [data-view="pronunciation"] span').evaluate((node) => getComputedStyle(node).backgroundImage);
  expect(pronunciationBg).toContain("icon-pronunciation-light.png?v=flux2-20260526-pronunciation-phrases-offline-roleplay");
  for (const view of ["phrasebook", "offline"]) {
    const background = await page.locator(`.mobile-bottom-nav-v2 [data-view="${view}"] span`).evaluate((node) => getComputedStyle(node).backgroundImage);
    expect(background).toContain(`icon-${view}-light.png?v=flux2-20260526-pronunciation-phrases-offline-roleplay`);
  }
  const railText = await page.locator(".mobile-bottom-nav-v2").innerText();
  expect(railText).not.toMatch(/AI micro lesson|PWA mini decks|Translate and convert|Repair drills/);
  const defaultRailViews = await page.locator(".mobile-bottom-nav-v2").evaluate((node) =>
    Array.from(node.querySelectorAll<HTMLButtonElement>("button[data-view]")).map((button) => button.dataset.view),
  );
  expect(defaultRailViews.slice(0, 6)).toEqual(["home", "tutor", "words", "word-game", "pronunciation", "shadowing"]);
  expect(defaultRailViews.at(-1)).toBe("settings");
  const initialScrollLeft = await page.locator(".mobile-bottom-nav-v2").evaluate((node) => {
    node.scrollLeft = 0;
    return node.scrollLeft;
  });
  await expect(page.locator(".mobile-bottom-nav-v2 button").first()).toHaveCSS("touch-action", "pan-x");
  await page.locator(".mobile-bottom-nav-v2").evaluate((node) => {
    node.scrollBy({ left: 220, behavior: "auto" });
  });
  await expect.poll(() => page.locator(".mobile-bottom-nav-v2").evaluate((node) => node.scrollLeft)).toBeGreaterThan(initialScrollLeft);
  await page.locator(".mobile-bottom-nav-v2").evaluate((node) => {
    node.scrollLeft = 0;
  });

  const touchEvent = (clientX = 0, clientY = 0, buttons = 1) => ({
    bubbles: true,
    cancelable: true,
    pointerType: "touch",
    pointerId: 7,
    isPrimary: true,
    button: 0,
    buttons,
    clientX,
    clientY,
  });
  const lessonButton = page.locator('.mobile-bottom-nav-v2 [data-view="lesson"]');
  const lessonBox = await lessonButton.boundingBox();
  expect(lessonBox).not.toBeNull();
  await lessonButton.dispatchEvent("pointerdown", touchEvent(lessonBox!.x + lessonBox!.width / 2, lessonBox!.y + lessonBox!.height / 2));
  await page.waitForTimeout(800);
  await expect(page.locator(".mobile-menu-edit-done-v2")).toBeVisible();
  const dragGhost = page.locator(".mobile-bottom-nav-drag-ghost-v2");
  await expect(dragGhost).toBeVisible();
  await expect(dragGhost).toHaveCSS("will-change", "transform");
  await lessonButton.dispatchEvent("pointerup", touchEvent(lessonBox!.x + lessonBox!.width / 2, lessonBox!.y + lessonBox!.height / 2, 0));
  await expect(page.locator(".mobile-bottom-nav-v2")).toHaveClass(/is-editing/);
  await expect(dragGhost).toHaveCount(0);
  await page.locator(".mobile-bottom-nav-v2").evaluate((node) => {
    node.scrollLeft = 0;
  });
  const scrollProbe = page.locator(".mobile-bottom-nav-v2 button[data-view]").nth(1);
  const scrollProbeBox = await scrollProbe.boundingBox();
  const railBox = await page.locator(".mobile-bottom-nav-v2").boundingBox();
  expect(scrollProbeBox).not.toBeNull();
  expect(railBox).not.toBeNull();
  await scrollProbe.dispatchEvent("pointerdown", touchEvent(scrollProbeBox!.x + scrollProbeBox!.width / 2, scrollProbeBox!.y + scrollProbeBox!.height / 2));
  await page.locator(".mobile-bottom-nav-v2").dispatchEvent("pointermove", touchEvent(railBox!.x + railBox!.width - 4, railBox!.y + railBox!.height / 2));
  await expect(dragGhost).toBeVisible();
  await expect.poll(() => page.locator(".mobile-bottom-nav-v2").evaluate((node) => node.scrollLeft)).toBeGreaterThan(0);
  await page.locator(".mobile-bottom-nav-v2").dispatchEvent("pointerup", touchEvent(railBox!.x + railBox!.width - 4, railBox!.y + railBox!.height / 2, 0));
  await expect(page.locator(".mobile-bottom-nav-v2")).toHaveClass(/is-editing/);
  await page.locator(".mobile-bottom-nav-v2").evaluate((node) => {
    node.scrollLeft = 0;
  });
  const firstRail = page.locator(".mobile-bottom-nav-v2 button[data-view]").first();
  const secondRail = page.locator(".mobile-bottom-nav-v2 button[data-view]").nth(1);
  const secondView = await secondRail.getAttribute("data-view");
  const firstBox = await firstRail.boundingBox();
  const secondBox = await secondRail.boundingBox();
  expect(firstBox).not.toBeNull();
  expect(secondBox).not.toBeNull();
  const storedBeforeReorderMove = await page.evaluate(() => localStorage.getItem("poliglot-mobile-nav-rail-v2:demor22"));
  await secondRail.dispatchEvent("pointerdown", touchEvent(secondBox!.x + secondBox!.width / 2, secondBox!.y + secondBox!.height / 2));
  await page.locator(".mobile-bottom-nav-v2").dispatchEvent("pointermove", touchEvent(firstBox!.x + firstBox!.width / 2, firstBox!.y + firstBox!.height / 2));
  await expect(dragGhost).toBeVisible();
  await expect(page.locator(".mobile-bottom-nav-v2 button.is-drop-target")).toHaveCount(1);
  await expect(page.locator(".mobile-bottom-nav-v2 button[data-view]").first()).toHaveAttribute("data-view", secondView || "");
  await expect.poll(() => page.evaluate(() => localStorage.getItem("poliglot-mobile-nav-rail-v2:demor22"))).toBe(storedBeforeReorderMove);
  await page.locator(".mobile-bottom-nav-v2").dispatchEvent("pointerup", touchEvent(firstBox!.x + firstBox!.width / 2, firstBox!.y + firstBox!.height / 2, 0));
  await expect(page.locator(".mobile-bottom-nav-v2 button[data-view]").first()).toHaveAttribute("data-view", secondView || "");
  await expect(page.locator(".mobile-menu-edit-done-v2")).toBeVisible();
  await expect(page.locator(".mobile-bottom-nav-v2")).toHaveClass(/is-editing/);
  const firstAfterDrop = page.locator(".mobile-bottom-nav-v2 button[data-view]").first();
  const secondAfterDrop = page.locator(".mobile-bottom-nav-v2 button[data-view]").nth(1);
  const firstAfterDropView = await firstAfterDrop.getAttribute("data-view");
  const firstAfterDropBox = await firstAfterDrop.boundingBox();
  const secondAfterDropBox = await secondAfterDrop.boundingBox();
  expect(firstAfterDropBox).not.toBeNull();
  expect(secondAfterDropBox).not.toBeNull();
  await firstAfterDrop.dispatchEvent("pointerdown", touchEvent(firstAfterDropBox!.x + firstAfterDropBox!.width / 2, firstAfterDropBox!.y + firstAfterDropBox!.height / 2));
  await page.locator(".mobile-bottom-nav-v2").dispatchEvent("pointermove", touchEvent(secondAfterDropBox!.x + secondAfterDropBox!.width * 0.75, secondAfterDropBox!.y + secondAfterDropBox!.height / 2));
  await page.locator(".mobile-bottom-nav-v2").dispatchEvent("pointerup", touchEvent(secondAfterDropBox!.x + secondAfterDropBox!.width * 0.75, secondAfterDropBox!.y + secondAfterDropBox!.height / 2, 0));
  await expect(page.locator(".mobile-bottom-nav-v2")).toHaveClass(/is-editing/);
  await expect(page.locator(".mobile-bottom-nav-v2 button[data-view]").nth(1)).toHaveAttribute("data-view", firstAfterDropView || "");
  await page.locator(".mobile-menu-edit-done-v2").click();
  await expect(page.locator(".mobile-bottom-nav-v2")).not.toHaveClass(/is-editing/);
  await expect(page.locator(".mobile-menu-edit-done-v2")).toHaveCount(0);
  await page.reload();
  await expect(page.locator(".mobile-bottom-nav-v2 button[data-view]").nth(1)).toHaveAttribute("data-view", firstAfterDropView || "");
});

test("mobile nav rail order reloads from the per-user server layout", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  let serverSession = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      navigation_layout: {
        mobile_rail: ["settings", "home", "tutor", "words", "word-game", "pronunciation", "shadowing"],
      },
    },
  };
  await page.unroute("**/api/session");
  await page.unroute("**/api/navigation-layout");
  await page.route("**/api/session", (route) => route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(serverSession) }));
  await page.route("**/api/navigation-layout", async (route) => {
    const payload = route.request().postDataJSON() as Record<string, unknown>;
    serverSession = { ...serverSession, user: { ...serverSession.user, navigation_layout: payload } };
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(serverSession) });
  });
  await page.goto("/app/?view=home");
  await expect(page.locator(".mobile-bottom-nav-v2 button[data-view]").first()).toHaveAttribute("data-view", "settings");
  await page.evaluate(() => localStorage.removeItem("poliglot-mobile-nav-rail-v2:demor22"));
  await page.reload();
  await expect(page.locator(".mobile-bottom-nav-v2 button[data-view]").first()).toHaveAttribute("data-view", "settings");
});

test("desktop function ribbon reorder persists after reload", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop ribbon assertion");
  await page.goto("/app/?view=home");
  await expect(page.locator(".function-ribbon-shell")).toBeVisible();
  await page.evaluate(() => localStorage.removeItem("poliglot-function-ribbon-v2:demor22"));
  await page.reload();
  const first = page.locator(".function-ribbon button[data-view]").first();
  const second = page.locator(".function-ribbon button[data-view]").nth(1);
  const secondView = await second.getAttribute("data-view");
  const firstBox = await first.boundingBox();
  const secondBox = await second.boundingBox();
  expect(firstBox).not.toBeNull();
  expect(secondBox).not.toBeNull();
  await second.dragTo(first, {
    sourcePosition: { x: secondBox!.width / 2, y: secondBox!.height / 2 },
    targetPosition: { x: firstBox!.width * 0.25, y: firstBox!.height / 2 },
  });
  await expect(page.locator(".function-ribbon button[data-view]").first()).toHaveAttribute("data-view", secondView || "");
  const stored = await page.evaluate(() => JSON.parse(localStorage.getItem("poliglot-function-ribbon-v2:demor22") || "[]"));
  expect(stored[0]).toBe(secondView);
  await page.reload();
  await expect(page.locator(".function-ribbon button[data-view]").first()).toHaveAttribute("data-view", secondView || "");
});

test("mobile lesson keeps output readable and phrase save inside input controls", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=lesson");
  await expect(page.locator(".context-display--lesson")).toBeVisible();
  await page.locator(".lesson-empty-v2 button").click();
  await expect(page.getByText("Say that you have a reservation.")).toBeVisible();
  await expect(page.locator(".phrase-quick-save-v2")).toBeVisible();

  const output = await page.locator(".chat-workspace__output").boundingBox();
  const composer = await page.locator(".composer-panel-v2").boundingBox();
  const fileControls = await page.locator(".file-controls-v2").boundingBox();
  const sendButton = await page.locator(".composer-submit-v2").boundingBox();
  await expect(page.locator(".composer-panel-v2 .phrase-quick-save-v2")).toBeVisible();
  let phraseSave = await page.locator(".phrase-quick-save-v2").boundingBox();
  expect(output).not.toBeNull();
  expect(composer).not.toBeNull();
  expect(fileControls).not.toBeNull();
  expect(sendButton).not.toBeNull();
  expect(phraseSave).not.toBeNull();
  expect(output!.height).toBeGreaterThan(composer!.height * 1.2);
  expect(composer!.y).toBeGreaterThanOrEqual(output!.y + output!.height - 2);
  expect(phraseSave!.y).toBeGreaterThanOrEqual(fileControls!.y + fileControls!.height - 2);
  expect(sendButton!.x + sendButton!.width).toBeLessThanOrEqual((await page.locator(".composer-textarea-shell-v2").boundingBox())!.x + (await page.locator(".composer-textarea-shell-v2").boundingBox())!.width + 2);
  expect(phraseSave!.height).toBeLessThanOrEqual(130);
  await page.locator(".composer-panel-v2").evaluate((node) => node.scrollTo(0, node.scrollHeight));
  phraseSave = await page.locator(".phrase-quick-save-v2").boundingBox();
  const viewport = page.viewportSize();
  expect(viewport).not.toBeNull();
  expect(phraseSave).not.toBeNull();
  expect(phraseSave!.y + phraseSave!.height).toBeLessThanOrEqual(viewport!.height - 64);

  const newLessonButton = await page.locator(".lesson-new-button-v2").boundingBox();
  const mobileNav = await page.locator(".mobile-bottom-nav-v2").boundingBox();
  expect(newLessonButton).not.toBeNull();
  expect(mobileNav).not.toBeNull();
  expect(newLessonButton!.y).toBeGreaterThanOrEqual(0);
  expect(newLessonButton!.y + newLessonButton!.height).toBeLessThanOrEqual(mobileNav!.y - 8);

  await page.locator(".composer-panel-v2 textarea").fill("I have a reservation.");
  await page.locator(".composer-panel-v2 textarea").press("Enter");
  await expect(page.getByText("Use under the name for reservations.")).toBeVisible();
  await expect(page.locator(".chat-interface__audio .audio-wave-button-v2").first()).toBeVisible();
});

test("mobile roleplay session has readable scenario, dialogue and input without overlap", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=roleplay");
  await expect(page.locator(".context-display--roleplay")).toBeVisible();
  await page.locator(".roleplay-grid-v2 button").first().click();
  await expect(page.locator(".roleplay-view-v2--session")).toBeVisible();
  await expect(page.locator(".roleplay-view-v2--mobile-session")).toBeVisible();
  await expect(page.locator(".roleplay-view-v2--session .roleplay-brief-v2")).toBeHidden();
  await expect(page.getByText(/ROLEPLAY DIALOGUE|Диалоговая сцена/)).toHaveCount(0);
  await expect(page.locator(".roleplay-dialog-head-v2 button")).toBeVisible();
  await expect(page.locator(".roleplay-dialog-scroll-v2")).toBeVisible();
  await expect(page.locator(".roleplay-dialog-scroll-v2")).not.toContainText("Short polite request.");
  await expect(page.locator(".roleplay-view-v2--session .composer-panel-v2")).toBeVisible();
  await expect(page.locator(".roleplay-view-v2--session .composer-textarea-shell-v2 textarea")).toBeVisible();
  await expect(page.locator(".roleplay-composer-v2")).toHaveCount(0);

  const dialog = await page.locator(".roleplay-dialog-scroll-v2").boundingBox();
  const composer = await page.locator(".roleplay-view-v2--session .composer-panel-v2").boundingBox();
  expect(dialog).not.toBeNull();
  expect(composer).not.toBeNull();
  expect(composer!.y).toBeGreaterThan(dialog!.y + dialog!.height - 4);

  await page.locator(".roleplay-view-v2--session .composer-textarea-shell-v2 textarea").fill("Could you repeat that, please?");
  await page.locator(".roleplay-view-v2--session .composer-textarea-shell-v2 textarea").press("Enter");
  await expect(page.locator(".roleplay-dialog-scroll-v2").getByText("Short polite request.").first()).toBeVisible();
});

test("regression: roleplay accepts a voice answer and shows the transcript", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop interaction assertion");
  let roleplayUploadSeen = false;
  await page.unroute("**/api/practice");
  await page.route("**/api/practice", async (route) => {
    const request = route.request();
    const contentType = request.headers()["content-type"] || "";
    const body = request.postData() || "";
    roleplayUploadSeen = contentType.includes("multipart/form-data") && body.includes("roleplay-voice.webm");
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        transcript: "Could you repeat that, please?",
        correction: "Could you repeat that, please?",
        explanation: "Short polite request.",
        correction_audio_text: "Could you repeat that, please?",
      }),
    });
  });
  await page.goto("/app/?view=roleplay");
  await expect(page.locator(".context-display--roleplay")).toBeVisible();
  await page.locator(".roleplay-grid-v2 button").first().click();
  await expect(page.locator(".roleplay-view-v2--desktop-session")).toBeVisible();
  await page.locator(".roleplay-view-v2--session .file-controls-v2 input[type='file']").setInputFiles({
    name: "roleplay-voice.webm",
    mimeType: "audio/webm",
    buffer: Buffer.from("voice-test"),
  });
  await page.locator(".roleplay-submit-v2").click();
  await expect(page.locator(".roleplay-dialog-scroll-v2").getByText("Could you repeat that, please?").first()).toBeVisible();
  expect(roleplayUploadSeen).toBe(true);
});

test("regression: image tool accepts a pasted clipboard image", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop clipboard assertion");
  let uploadedImageSeen = false;
  await page.unroute("**/api/tools/image-translate");
  await page.route("**/api/tools/image-translate", async (route) => {
    const request = route.request();
    const contentType = request.headers()["content-type"] || "";
    const body = request.postData() || "";
    uploadedImageSeen = contentType.includes("multipart/form-data") && body.includes("clipboard-image-from-page.png");
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        source_text: "Cafe menu",
        translation: "Меню кафе",
        result: "Меню кафе",
        source_language: "auto",
        target_language: "ru",
      }),
    });
  });

  await page.goto("/app/?view=tools");
  await expect(page.locator(".context-display--tools")).toBeVisible();
  await page.locator(".tool-switch-v2 button").nth(2).click();
  await expect(page.locator(".image-upload-control-v2")).toBeVisible();
  await page.locator(".image-upload-control-v2").evaluate((node, bytes) => {
    const file = new File([new Uint8Array(bytes)], "clipboard-image.png", { type: "image/png" });
    const dataTransfer = new DataTransfer();
    dataTransfer.items.add(file);
    node.dispatchEvent(new ClipboardEvent("paste", { clipboardData: dataTransfer, bubbles: true, cancelable: true }));
  }, Array.from(png));
  await expect(page.locator(".file-chip-v2")).toContainText("clipboard-image.png");
  await page.locator(".file-chip-v2 button").click();
  await expect(page.locator(".file-chip-v2")).toHaveCount(0);
  await page.locator(".context-display--tools").evaluate((node, bytes) => {
    const file = new File([new Uint8Array(bytes)], "clipboard-image-from-page.png", { type: "image/png" });
    const dataTransfer = new DataTransfer();
    dataTransfer.items.add(file);
    node.dispatchEvent(new ClipboardEvent("paste", { clipboardData: dataTransfer, bubbles: true, cancelable: true }));
  }, Array.from(png));
  await expect(page.locator(".file-chip-v2")).toContainText("clipboard-image-from-page.png");
  await page.locator(".tools-submit-v2").click();
  await expect(page.locator(".chat-interface").getByText("Меню кафе").first()).toBeVisible();
  expect(uploadedImageSeen).toBe(true);
});

test("desktop roleplay session uses the same right input window as practice", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop layout assertion");
  await page.goto("/app/?view=roleplay");
  await expect(page.locator(".context-display--roleplay")).toBeVisible();
  await page.locator(".roleplay-grid-v2 button").first().click();
  await expect(page.locator(".roleplay-view-v2--desktop-session")).toBeVisible();
  await expect(page.locator(".roleplay-view-v2--session .roleplay-brief-v2")).toHaveCount(0);
  await expect(page.locator(".roleplay-view-v2--session .chat-workspace--roleplay")).toBeVisible();
  await expect(page.locator(".roleplay-view-v2--session .composer-panel-v2")).toBeVisible();
  await expect(page.locator(".roleplay-view-v2--session .composer-textarea-shell-v2")).toBeVisible();
  await expect(page.locator(".roleplay-composer-v2")).toHaveCount(0);

  const panel = await page.locator(".roleplay-view-v2--session").boundingBox();
  const dialog = await page.locator(".roleplay-dialog-card-v2").boundingBox();
  const composer = await page.locator(".roleplay-view-v2--session .composer-panel-v2").boundingBox();
  expect(panel).not.toBeNull();
  expect(dialog).not.toBeNull();
  expect(composer).not.toBeNull();
  expect(dialog!.width).toBeLessThan(panel!.width * 0.76);
  expect(composer!.x).toBeGreaterThan(dialog!.x + dialog!.width - 4);
});

test("regression: quick save bookmark becomes active after saving to notes", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop state assertion");
  await page.goto("/app/?view=practice");
  await expect(page.locator(".context-display--practice")).toBeVisible();
  await page.locator(".composer-panel-v2 textarea").fill("Can you repeat?");
  await page.locator(".composer-panel-v2 textarea").press("Enter");
  const saveButton = page.locator(".phrase-quick-save-v2 button").filter({ hasText: "Could you repeat that, please?" }).first();
  await expect(saveButton).toBeVisible();
  await expect(saveButton).toHaveAttribute("aria-pressed", "false");
  const colorBefore = await saveButton.locator("svg").evaluate((node) => getComputedStyle(node).color);
  await saveButton.click();
  await expect(saveButton).toHaveAttribute("aria-pressed", "true");
  await expect(saveButton).toHaveClass(/is-saved/);
  const colorAfter = await saveButton.locator("svg").evaluate((node) => getComputedStyle(node).color);
  expect(colorAfter).not.toBe(colorBefore);
});

test("mobile roleplay scenario cards expand and scroll above bottom menu", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=roleplay");
  await expect(page.locator(".context-display--roleplay")).toBeVisible();
  await page.locator(".context-display").evaluate((node) => node.scrollTo(0, node.scrollHeight));

  const cards = await page.locator(".roleplay-grid-v2 button").evaluateAll((nodes) =>
    nodes.map((node) => {
      const element = node as HTMLElement;
      return {
        height: element.getBoundingClientRect().height,
        offsetHeight: element.offsetHeight,
        scrollHeight: element.scrollHeight,
      };
    }),
  );
  expect(cards.length).toBeGreaterThan(8);
  expect(cards.filter((card) => card.scrollHeight > card.offsetHeight + 2)).toEqual([]);
  const clippedText = await page.locator(".roleplay-grid-v2 button").evaluateAll((nodes) =>
    nodes.flatMap((node, cardIndex) => {
      const card = (node as HTMLElement).getBoundingClientRect();
      return Array.from((node as HTMLElement).querySelectorAll<HTMLElement>("strong, span, small"))
        .filter((child) => {
          const box = child.getBoundingClientRect();
          return box.bottom > card.bottom + 2 || box.right > card.right + 2 || box.left < card.left - 2;
        })
        .map((child) => ({ cardIndex, text: child.textContent || "", cardBottom: card.bottom, childBottom: child.getBoundingClientRect().bottom }));
    }),
  );
  expect(clippedText).toEqual([]);

  const lastCard = page.locator(".roleplay-grid-v2 button").last();
  await expect(lastCard).toBeVisible();
  const cardBox = await lastCard.boundingBox();
  const navBox = await page.locator(".mobile-bottom-nav-v2").boundingBox();
  const listBox = await page.locator(".roleplay-grid-v2").boundingBox();
  expect(cardBox).not.toBeNull();
  expect(navBox).not.toBeNull();
  expect(listBox).not.toBeNull();
  expect(listBox!.y + listBox!.height).toBeGreaterThanOrEqual(cardBox!.y + cardBox!.height - 2);
  expect(cardBox!.y + cardBox!.height).toBeLessThanOrEqual(navBox!.y - 4);
});

test("shadowing uses one compact work panel with the sample audio inside the task", async ({ page }) => {
  usePremiumSession();
  const phrases = ["Could you repeat that, please?", "The train leaves at nine."];
  let phraseIndex = 0;
  const spokenTexts: string[] = [];
  await page.addInitScript(() => {
    Object.defineProperty(HTMLMediaElement.prototype, "play", {
      configurable: true,
      value() {
        setTimeout(() => this.dispatchEvent(new Event("ended")), 0);
        return Promise.resolve();
      },
    });
    Object.defineProperty(HTMLMediaElement.prototype, "pause", {
      configurable: true,
      value() {},
    });
  });
  await page.unroute("**/api/shadowing/start");
  await page.route("**/api/shadowing/start", (route) => {
    const phrase = phrases[Math.min(phraseIndex, phrases.length - 1)];
    phraseIndex += 1;
    return route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ phrase, target: phrase }),
    });
  });
  await page.route("**/api/tools/translator-speech", async (route) => {
    const payload = route.request().postDataJSON() as Record<string, unknown>;
    spokenTexts.push(String(payload.text || ""));
    await route.fulfill({
      status: 200,
      contentType: "audio/mpeg",
      body: Buffer.from("test-audio"),
    });
  });
  await page.goto("/app/?view=home");
  await page.locator(".home-plan-v2__steps button").filter({ hasText: appCopy("ru", "today_plan_listen_title") }).click();
  await expect(page.locator(".context-display--shadowing")).toBeVisible();
  await expect(page.locator(".chat-workspace--single")).toBeVisible();
  await expect(page.locator(".chat-workspace__output")).toHaveCount(0);
  await expect(page.locator(".task-box-v2 .audio-wave-button-v2")).toBeVisible();
  await page.locator(".task-box-v2 .audio-wave-button-v2").click();
  await expect.poll(() => spokenTexts).toEqual([phrases[0]]);
  await expect(page.locator(".task-box-v2 .audio-wave-button-v2")).not.toHaveClass(/is-playing/);
  await expect(page.locator(".task-box-v2")).not.toContainText("Could you repeat that, please?");
  await expect(page.locator(".task-box-v2 strong")).toHaveCount(0);
  await page.locator(".task-box-v2__next").click();
  await expect.poll(() => phraseIndex).toBeGreaterThan(1);
  await expect(page.locator(".chat-workspace--single")).toBeVisible();
  await expect(page.locator(".chat-workspace__output")).toHaveCount(0);
  await expect(page.locator(".task-box-v2 .audio-wave-button-v2")).toHaveCount(1);
  await expect(page.locator(".task-box-v2 .audio-wave-button-v2")).toBeEnabled();
  await expect(page.locator(".task-box-v2 .audio-wave-button-v2")).not.toHaveClass(/is-playing/);
  await page.locator(".task-box-v2 .audio-wave-button-v2").click();
  await expect.poll(() => spokenTexts).toEqual([phrases[0], phrases[1]]);
});

test("mobile listening panel stays above bottom menu", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  usePremiumSession();
  await page.goto("/app/?view=shadowing");
  await expect(page.locator(".context-display--shadowing")).toBeVisible();
  await expect(page.locator(".task-box-v2 .audio-wave-button-v2")).toBeVisible();
  await page.locator(".context-display").evaluate((node) => node.scrollTo(0, node.scrollHeight));
  const taskBox = await page.locator(".task-box-v2").boundingBox();
  const navBox = await page.locator(".mobile-bottom-nav-v2").boundingBox();
  expect(taskBox).not.toBeNull();
  expect(navBox).not.toBeNull();
  expect(taskBox!.y + taskBox!.height).toBeLessThanOrEqual(navBox!.y - 4);
});

test("mobile practice and tools keep input controls visible and tools selectable", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=practice");
  await expect(page.locator(".context-display--practice")).toBeVisible();
  await expect(page.locator(".mobile-quick-controls-v2")).toHaveCount(0);
  await page.locator(".context-display").evaluate((node) => node.scrollTo(0, node.scrollHeight));
  const practiceOutput = await page.locator(".chat-workspace__output").boundingBox();
  const practiceComposer = await page.locator(".composer-panel-v2").boundingBox();
  const practiceInput = await page.locator(".composer-panel-v2 textarea").boundingBox();
  const viewport = page.viewportSize();
  expect(practiceOutput).not.toBeNull();
  expect(practiceComposer).not.toBeNull();
  expect(practiceInput).not.toBeNull();
  expect(viewport).not.toBeNull();
  expect(practiceOutput!.height).toBeGreaterThan(practiceComposer!.height * 1.2);
  expect(practiceInput!.y + practiceInput!.height).toBeLessThanOrEqual(viewport!.height - 64);
  await page.locator(".composer-panel-v2 textarea").fill("Can you repeat?");
  await page.locator(".composer-panel-v2 textarea").press("Enter");
  await expect(page.locator(".phrase-quick-save-v2")).toBeVisible();
  const quickSaveText = await page.locator(".phrase-quick-save-v2").innerText();
  expect(quickSaveText).toContain("Could you repeat that, please?");
  expect(quickSaveText.split(/\s+/)).not.toContain("?");
  const practiceFileControls = await page.locator(".file-controls-v2").boundingBox();
  const practiceSend = await page.locator(".composer-submit-v2").boundingBox();
  const practiceSave = await page.locator(".phrase-quick-save-v2").boundingBox();
  expect(practiceFileControls).not.toBeNull();
  expect(practiceSend).not.toBeNull();
  expect(practiceSave).not.toBeNull();
  expect(practiceSave!.y).toBeGreaterThanOrEqual(practiceFileControls!.y + practiceFileControls!.height - 2);
  expect(practiceSend!.x + practiceSend!.width).toBeLessThanOrEqual((await page.locator(".composer-textarea-shell-v2").boundingBox())!.x + (await page.locator(".composer-textarea-shell-v2").boundingBox())!.width + 2);

  await page.goto("/app/?view=tools");
  await expect(page.locator(".context-display--tools")).toBeVisible();
  await expect(page.locator(".tool-switch-v2 button")).toHaveCount(3);
  await expect(page.locator(".tools-work-v2 .composer-panel-v2 textarea:visible")).toHaveCount(0);
  await page.locator(".tool-switch-v2 button").nth(1).click();
  await expect(page.locator(".tool-switch-v2 button:visible")).toHaveCount(0);
  await expect(page.locator(".tools-change-v2")).toBeVisible();
  await expect(page.locator(".tools-submit-v2")).toBeVisible();
  const inputShell = await page.locator(".tool-input-shell-v2").boundingBox();
  const submitButton = await page.locator(".tools-submit-v2").boundingBox();
  const toolsInput = await page.locator(".composer-panel-v2 textarea").boundingBox();
  expect(inputShell).not.toBeNull();
  expect(submitButton).not.toBeNull();
  expect(toolsInput).not.toBeNull();
  expect(toolsInput!.y + toolsInput!.height).toBeLessThanOrEqual(viewport!.height - 64);
  expect(submitButton!.x + submitButton!.width).toBeLessThanOrEqual(inputShell!.x + inputShell!.width + 2);
  expect(submitButton!.y + submitButton!.height).toBeLessThanOrEqual(inputShell!.y + inputShell!.height + 2);
});

test("mobile mistakes dictionary paginates after ten cards", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=mistakes");
  await expect(page.locator(".context-display--mistakes")).toBeVisible();
  await expect(page.locator(".mistake-pagination-v2")).toBeVisible();
  await expect(page.locator(".mistake-list-card-v2")).toHaveCount(10);
  await expect(page.getByText("wrong phrase 18")).toHaveCount(0);
  await page.locator(".mistake-pagination-v2 button").nth(1).click();
  await expect(page.locator(".mistake-list-card-v2")).toHaveCount(8);
  await expect(page.getByText("wrong phrase 18")).toBeVisible();
});

test("mobile phrasebook list can scroll below the bottom menu", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=phrasebook");
  await expect(page.locator(".context-display--phrasebook")).toBeVisible();
  await page.getByRole("button", { name: /Добавить пример|Add example/ }).click();
  await page.locator(".context-display").evaluate((node) => node.scrollTo(0, node.scrollHeight));
  const lastCard = page.locator(".phrasebook-card-v2").last();
  await expect(lastCard).toBeVisible();
  const cardBox = await lastCard.boundingBox();
  const navBox = await page.locator(".mobile-bottom-nav-v2").boundingBox();
  expect(cardBox).not.toBeNull();
  expect(navBox).not.toBeNull();
  expect(cardBox!.y + cardBox!.height).toBeLessThanOrEqual(navBox!.y - 4);
});

test("mobile phrasebook paginates notes after ten cards", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  const manyPhrases = Array.from({ length: 12 }, (_, index) => ({
    id: `phrase-page-${index + 1}`,
    phrase: `Saved phrase ${index + 1}`,
    translation: `Translation ${index + 1}`,
    source: "manual",
    language: "en",
    createdAt: new Date().toISOString(),
  }));
  await page.addInitScript((phrases) => {
    localStorage.setItem("poliglot-phrasebook-v2:demor22", JSON.stringify(phrases));
  }, manyPhrases);
  testSessionPayloadOverride = { ...sessionPayload, user: { ...sessionPayload.user, phrasebook: manyPhrases } };
  testPhrasebookItemsOverride = manyPhrases;
  await page.unrouteAll({ behavior: "ignoreErrors" });
  await page.route("**/app/assets/**", (route) => {
    const request = route.request();
    const assetPath = new URL(request.url()).pathname.toLowerCase();
    if (request.resourceType() === "image" || /\.(avif|gif|ico|jpe?g|png|svg|webp)$/.test(assetPath)) {
      return route.fulfill({ status: 200, contentType: "image/png", body: png });
    }
    return route.continue();
  });
  await page.route("**/api/session", (route) =>
    route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ ...sessionPayload, user: { ...sessionPayload.user, phrasebook: manyPhrases } }),
    }),
  );
  await page.route("**/api/daily/claim", (route) =>
    route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, already_claimed: true }) }),
  );
  await page.route("**/api/phrasebook**", async (route) => {
    if (route.request().method() === "GET") {
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: manyPhrases }) });
      return;
    }
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, items: manyPhrases }) });
  });
  await page.goto("/app/?view=phrasebook");
  await expect(page.locator(".context-display--phrasebook")).toBeVisible();
  await expect(page.locator(".phrasebook-pagination-v2")).toBeVisible();
  await expect(page.locator(".phrasebook-card-v2")).toHaveCount(10);
  await expect(page.getByText("Saved phrase 11")).toHaveCount(0);
  await page.locator(".phrasebook-pagination-v2 button").nth(1).click();
  await expect(page.locator(".phrasebook-card-v2")).toHaveCount(2);
  await expect(page.getByText("Saved phrase 11")).toBeVisible();
});

test("mobile award details open in the current viewport", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=awards");
  await expect(page.locator(".context-display--awards")).toBeVisible();
  await page.locator(".context-display").evaluate((node) => node.scrollTo(0, 420));
  await page.locator(".award-tile-v2.is-unlocked").nth(4).click();
  await expect(page.locator(".award-modal-v2")).toBeVisible();
  const viewport = page.viewportSize();
  const modal = await page.locator(".award-modal-v2").boundingBox();
  expect(viewport).not.toBeNull();
  expect(modal).not.toBeNull();
  expect(modal!.y).toBeGreaterThanOrEqual(0);
  expect(modal!.y + modal!.height).toBeLessThanOrEqual(viewport!.height);
});

test("phrasebook loads from session and save calls persistent API", async ({ page }) => {
  let savedPhrase = "";
  await page.route("**/api/phrasebook**", async (route) => {
    const request = route.request();
    if (request.method() === "POST") {
      const item = request.postDataJSON();
      savedPhrase = String(item.phrase || "");
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, items: [item, ...phrasebookSeed] }) });
      return;
    }
    await route.fallback();
  });
  await page.goto("/app/?view=phrasebook");
  await expect(page.getByText("I have a reservation under the name Ivan Petrov.")).toBeVisible();
  await page.getByRole("button", { name: /Добавить пример|Add sample/ }).click();
  await expect.poll(() => savedPhrase).toContain("Could you say that again");
});

test("regression: paid Platinum subscription is marked as the current plan", async ({ page }) => {
  usePremiumSession({ plan: "Platinum 2026-07-12", premium_until: "2026-07-12T12:00:00Z" });
  await page.goto("/app/?view=premium");
  const freeCard = page.locator(".plan-card-v2.is-free");
  const platinumCard = page.locator(".plan-card-v2").filter({ hasText: /Platinum 30/ }).first();
  await expect(freeCard).toBeVisible();
  await expect(platinumCard).toBeVisible();
  await expect(freeCard.getByRole("button", { name: /Текущий план|Current plan/ })).toHaveCount(0);
  await expect(platinumCard.getByRole("button", { name: /Текущий план|Current plan/ })).toBeVisible();
});

test("regression: phrasebook note field can add a manual note without a phrase", async ({ page }) => {
  let savedItem: Record<string, unknown> | null = null;
  await page.route("**/api/phrasebook**", async (route) => {
    const request = route.request();
    if (request.method() === "POST") {
      savedItem = request.postDataJSON() as Record<string, unknown>;
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ ok: true, items: [savedItem, ...phrasebookSeed] }) });
      return;
    }
    await route.fallback();
  });
  await page.goto("/app/?view=phrasebook");
  await page.getByPlaceholder(ru("phrasebook_note_placeholder", "Заметка или перевод")).fill("Today is a good day.");
  await page.locator(".phrasebook-add-v2 button[type='submit']").click();
  await expect.poll(() => String(savedItem?.phrase || "")).toBe("Today is a good day.");
  expect(String(savedItem?.note || savedItem?.translation || "")).toBe("");
});

test("regression: audio buttons show a loading status while speech is generated", async ({ page }) => {
  let resolveSpeech: (() => void) | null = null;
  await page.addInitScript(() => {
    Object.defineProperty(HTMLMediaElement.prototype, "play", {
      configurable: true,
      value() {
        setTimeout(() => this.dispatchEvent(new Event("ended")), 0);
        return Promise.resolve();
      },
    });
  });
  await page.unroute("**/api/tools/translator-speech");
  await page.route("**/api/tools/translator-speech", async (route) => {
    await new Promise<void>((resolve) => {
      resolveSpeech = resolve;
    });
    await route.fulfill({ status: 200, contentType: "audio/mpeg", body: Buffer.from("test-audio") });
  });
  await page.goto("/app/?view=phrasebook");
  const audioButton = page.locator(".phrasebook-card-v2 .audio-wave-button-v2").first();
  await audioButton.click();
  await expect(audioButton).toContainText(/Загрузка|Loading/);
  resolveSpeech?.();
  await expect(audioButton).not.toContainText(/Загрузка|Loading/);
});

test("regression: pronunciation history normalizes over-limit score", async ({ page }) => {
  usePremiumSession();
  await page.addInitScript(() => {
    localStorage.setItem("poliglot-pronunciation-v2:demor22", JSON.stringify([
      {
        score: 222,
        feedback: "Оценка: 22/100. Балл ограничен: сравниваем распознанный текст с образцом.",
        expected: "Do you think this shirt looks good on me",
        spoken: "Do you think this shirt",
        createdAt: "2026-06-19T12:00:00.000Z",
      },
    ]));
  });
  await page.goto("/app/?view=pronunciation");
  await expect(page.locator(".pronunciation-score-v2")).toContainText("22/100");
  await expect(page.locator(".pronunciation-history-v2")).toContainText("22/100");
  await expect(page.locator(".pronunciation-history-v2")).not.toContainText("222/100");
});

test("regression: auth language menu is layered above privacy and captcha blocks", async ({ page }) => {
  await mockAnonymousAuth(page);
  await page.goto("/app/login");
  await page.locator(".auth-page-language-v2 > button").click();
  const menu = page.locator(".auth-page-language-v2 .animated-select-v2__menu");
  await expect(menu).toBeVisible();
  const zIndex = await menu.evaluate((node) => Number(getComputedStyle(node).zIndex || 0));
  expect(zIndex).toBeGreaterThanOrEqual(1200);
  const menuBox = await menu.boundingBox();
  const privacyBox = await page.locator(".auth-privacy-v2").boundingBox();
  expect(menuBox).not.toBeNull();
  expect(privacyBox).not.toBeNull();
  expect(menuBox!.y + menuBox!.height).toBeGreaterThan(privacyBox!.y);
});

test("main mobile and desktop views have scrollable output without mojibake", async ({ page, isMobile }) => {
  const views = isMobile ? ["home", "offline", "referral", "settings", "tools"] : ["home", "offline", "premium", "referral", "settings"];
  for (const view of views) {
    await page.goto(`/app/?view=${view}`);
    await expect(page.locator(`.context-display--${view}`)).toBeVisible();
    await expectNoMojibake(page);
    const metrics = await page.locator(".context-display").evaluate((node) => ({
      clientHeight: node.clientHeight,
      scrollHeight: node.scrollHeight,
    }));
    expect(metrics.scrollHeight).toBeGreaterThanOrEqual(metrics.clientHeight);
  }
});

test("regression: desktop offline controls are folded into the deck panel on the right", async ({ page, isMobile }) => {
  test.skip(isMobile, "desktop layout assertion");
  await page.goto("/app/?view=offline");
  await expect(page.locator(".context-display--offline")).toBeVisible();
  const panelBox = await page.locator(".offline-deck-list-v2").boundingBox();
  const summaryBox = await page.locator(".offline-deck-summary-v2").boundingBox();
  const controlsBox = await page.locator(".offline-hero-v2").boundingBox();
  expect(panelBox).not.toBeNull();
  expect(summaryBox).not.toBeNull();
  expect(controlsBox).not.toBeNull();
  expect(controlsBox!.x).toBeGreaterThan(summaryBox!.x + summaryBox!.width - 8);
  expect(controlsBox!.x + controlsBox!.width).toBeLessThanOrEqual(panelBox!.x + panelBox!.width + 2);
  expect(controlsBox!.y).toBeGreaterThanOrEqual(panelBox!.y + 12);
  expect(controlsBox!.y + controlsBox!.height).toBeLessThan(panelBox!.y + panelBox!.height);
});

test("regression: dashboard training balance uses the current week and does not max every bar", async ({ page }) => {
  const today = new Date();
  const dayKey = (offset: number) => {
    const date = new Date(today);
    date.setDate(today.getDate() - offset);
    return date.toISOString().slice(0, 10);
  };
  testSessionPayloadOverride = {
    ...sessionPayload,
    user: {
      ...sessionPayload.user,
      lesson_limit: 2,
      practice_limit: 2,
      voice_limit: 2,
      lessons_today: 1,
      practice_today: 0,
      voice_today: 0,
      lesson_count: 42,
      practice_count: 42,
      word_game_count: 42,
      learned_words: 42,
      habit_log: {
        [dayKey(0)]: { login: true, complete: false, lessons: 1, practice: 0, voice: 0 },
        [dayKey(1)]: { login: true, complete: true, lessons: 1, practice: 1, voice: 0 },
        [dayKey(2)]: { login: true, complete: false, lessons: 0, practice: 1, voice: 0 },
        [dayKey(8)]: { login: true, complete: true, lessons: 7, practice: 7, voice: 7 },
      },
    },
  };
  await page.goto("/app/?view=dashboard");
  await expect(page.locator(".dashboard-rhythm-v2")).toBeVisible();
  const barRatios = await page.locator(".dashboard-rhythm-v2__bars span").evaluateAll((nodes) =>
    nodes.map((node) => {
      const row = node as HTMLElement;
      const bar = row.querySelector("b") as HTMLElement | null;
      const rowWidth = row.getBoundingClientRect().width;
      const barWidth = bar?.getBoundingClientRect().width || 0;
      return Number((barWidth / Math.max(1, rowWidth)).toFixed(2));
    }),
  );
  expect(barRatios.length).toBeGreaterThanOrEqual(5);
  expect(barRatios.filter((ratio) => ratio < 0.7).length).toBeGreaterThanOrEqual(3);
  expect(Math.max(...barRatios)).toBeLessThan(0.95);
});

test("regression: chat messages show sender first and full-width text below it", async ({ page }) => {
  await page.goto("/app/?view=lesson");
  await page.locator(".lesson-empty-v2 button").click();
  await expect(page.locator(".chat-interface__group").first()).toBeVisible();
  const metrics = await page.locator(".chat-interface__group").first().evaluate((group) => {
    const meta = group.querySelector(".chat-interface__meta") as HTMLElement | null;
    const bubble = group.querySelector(".chat-interface__bubble") as HTMLElement | null;
    const groupBox = group.getBoundingClientRect();
    const metaBox = meta?.getBoundingClientRect();
    const bubbleBox = bubble?.getBoundingClientRect();
    return {
      groupWidth: groupBox.width,
      metaY: metaBox?.y || 0,
      bubbleY: bubbleBox?.y || 0,
      bubbleWidth: bubbleBox?.width || 0,
    };
  });
  expect(metrics.bubbleY).toBeGreaterThan(metrics.metaY);
  expect(metrics.bubbleWidth).toBeGreaterThan(metrics.groupWidth * 0.92);
});

test("regression: mobile awards header stays compact and mistakes open as list then practice screen", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");

  await page.goto("/app/?view=awards");
  await expect(page.locator(".awards-summary-v2")).toBeVisible();
  const awardsHeader = await page.locator(".awards-summary-v2").boundingBox();
  expect(awardsHeader).not.toBeNull();
  expect(awardsHeader!.height).toBeLessThanOrEqual(180);

  await page.goto("/app/?view=mistakes");
  await expect(page.locator(".mistake-layout-v2")).toBeVisible();
  await expect(page.locator(".mistake-dictionary-v2")).toBeVisible();
  await expect(page.locator(".mistake-practice-v2")).toHaveCount(0);
  await page.locator(".mistake-list-card-v2 button").first().click();
  await expect(page.locator(".mistake-practice-v2")).toBeVisible();
  await expect(page.locator(".mistake-dictionary-v2")).toHaveCount(1);
  await expect(page.locator(".mistake-dictionary-v2--practice")).toBeVisible();
  const practice = await page.locator(".mistake-practice-v2").boundingBox();
  expect(practice).not.toBeNull();
  expect(practice!.height).toBeGreaterThan(340);
  await page.locator(".mistake-practice-v2 input").fill("Correct phrase 1.");
  await page.locator(".mistake-practice-v2 button", { hasText: ru("check", "Check") }).click();
  await expect(page.locator(".mistake-dictionary-v2")).toBeVisible();
  await expect(page.locator(".mistake-dictionary-v2--practice")).toHaveCount(0);
  await expect(page.locator(".mistake-practice-v2")).toHaveCount(0);
});

test("regression: mobile cards expand to fit text in leaderboard notes and offline decks", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");

  const expectNoCardOverflow = async (selector: string) => {
    const overflows = await page.locator(selector).evaluateAll((nodes) =>
      nodes.slice(0, 6).map((node) => {
        const element = node as HTMLElement;
        return {
          text: element.innerText.slice(0, 80),
          height: element.getBoundingClientRect().height,
          scrollHeight: element.scrollHeight,
        };
      }).filter((item) => item.scrollHeight > item.height + 2),
    );
    expect(overflows).toEqual([]);
  };

  const expectPanelWrapsLastItem = async (panelSelector: string, itemSelector: string) => {
    const panelBox = await page.locator(panelSelector).boundingBox();
    const lastItemBox = await page.locator(itemSelector).last().boundingBox();
    expect(panelBox).not.toBeNull();
    expect(lastItemBox).not.toBeNull();
    expect(panelBox!.y + panelBox!.height).toBeGreaterThanOrEqual(lastItemBox!.y + lastItemBox!.height - 2);
  };

  await page.goto("/app/?view=leaderboard");
  await expect(page.locator(".context-display--leaderboard")).toBeVisible();
  await expect(page.locator(".context-display--leaderboard h2")).toContainText("Общий");
  await expect(page.locator(".leaderboard-select-v2")).toContainText("Общий");
  await expectNoCardOverflow(".leaderboard-row-v2");
  const leaderboardClippedText = await page.locator(".leaderboard-row-v2").evaluateAll((nodes) =>
    nodes.flatMap((node, rowIndex) => {
      const row = (node as HTMLElement).getBoundingClientRect();
      return Array.from((node as HTMLElement).querySelectorAll<HTMLElement>("strong, em, small"))
        .filter((child) => {
          const box = child.getBoundingClientRect();
          return box.bottom > row.bottom + 2 || box.right > row.right + 2 || box.left < row.left - 2;
        })
        .map((child) => ({ rowIndex, text: child.textContent || "", rowBottom: row.bottom, childBottom: child.getBoundingClientRect().bottom }));
    }),
  );
  expect(leaderboardClippedText).toEqual([]);
  await expectPanelWrapsLastItem(".context-display--leaderboard .data-display", ".leaderboard-row-v2");

  await page.goto("/app/?view=roleplay");
  await expect(page.locator(".context-display--roleplay")).toBeVisible();
  await expectNoCardOverflow(".roleplay-grid-v2 button");
  await expectPanelWrapsLastItem(".context-display--roleplay .roleplay-view-v2", ".roleplay-grid-v2 button");

  await page.goto("/app/?view=phrasebook");
  await expect(page.locator(".context-display--phrasebook")).toBeVisible();
  await expectNoCardOverflow(".phrasebook-card-v2");

  await page.goto("/app/?view=vocabulary");
  await expect(page.locator(".context-display--vocabulary")).toBeVisible();
  await expect(page.getByText("самолётостроение")).toBeVisible();
  await expectNoCardOverflow(".vocabulary-card-v2");

  await page.goto("/app/?view=level");
  await expect(page.locator(".context-display--level")).toBeVisible();
  await expect(page.locator(".level-answer-skip-v2")).toHaveText("Skip");
  await expect(page.locator(".level-answer-skip-v2")).not.toHaveText("Раздел");
  const skipTextFits = await page.locator(".level-answer-skip-v2 > span").last().evaluate((node) => {
    const element = node as HTMLElement;
    return element.scrollWidth <= element.clientWidth + 1 && element.scrollHeight <= element.clientHeight + 1;
  });
  expect(skipTextFits).toBe(true);
  await expectNoCardOverflow(".trainer-display");

  await page.goto("/app/?view=offline");
  await expect(page.locator(".context-display--offline")).toBeVisible();
  await expectNoCardOverflow(".offline-deck-grid-v2 article");
});

test("regression: leaderboard language selector opens above the panel", async ({ page }) => {
  await page.goto("/app/?view=leaderboard");
  await expect(page.locator(".context-display--leaderboard")).toBeVisible();
  const select = page.locator(".context-display--leaderboard .leaderboard-select-v2").first();
  await expect(select).toBeVisible();
  await select.click();
  const listbox = select.locator(".animated-select-v2__menu");
  await expect(listbox).toBeVisible();
  const zIndex = await listbox.evaluate((node) => Number(getComputedStyle(node).zIndex || 0));
  expect(zIndex).toBeGreaterThanOrEqual(40);
});

test("regression: mobile leaderboard stays readable above bottom menu", async ({ page, isMobile }) => {
  test.skip(!isMobile, "mobile layout assertion");
  await page.goto("/app/?view=leaderboard");
  await expect(page.locator(".context-display--leaderboard")).toBeVisible();
  await expect(page.locator(".leaderboard-row-v2")).toHaveCount(10);
  await page.locator(".context-display").evaluate((node) => node.scrollTo(0, node.scrollHeight));
  const lastRow = page.locator(".leaderboard-row-v2").last();
  await expect(lastRow).toBeVisible();
  const rowBox = await lastRow.boundingBox();
  const navBox = await page.locator(".mobile-bottom-nav-v2").boundingBox();
  expect(rowBox).not.toBeNull();
  expect(navBox).not.toBeNull();
  expect(rowBox!.y + rowBox!.height).toBeLessThanOrEqual(navBox!.y - 4);
});

