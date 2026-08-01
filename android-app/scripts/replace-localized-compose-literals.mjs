import fs from "node:fs";
import path from "node:path";

const uiRoot = path.resolve("app/src/main/java/ru/neriva/app/ui");
const stringsPath = path.resolve("app/src/main/res/values/strings.xml");
const valueToKeys = new Map();
const literalToResource = {
  "Report word": "ai_tutor_word_report_button",
  "What's wrong? (optional)": "ai_tutor_word_report_comment_label",
  "Premium daily practice": "premium_month_body",
  "AI Language Tutor": "ai_tutor",
  "Sign In": "auth_login",
  "Register": "auth_create",
  "Confirm Password": "confirm_password",
  "Referral code (optional)": "auth_referral",
  "Completed Lessons": "tutor_completed_lessons",
  "No completed lessons yet.": "tutor_completed_lessons_empty",
  "Overall Progress": "progress_without_leaderboard",
  "Premium Active": "premium",
  "Your progress": "progress",
  "Daily activity": "learning_activity",
  "Daily Bonus": "claim_daily_bonus",
  "Learning Lab": "v2_learning_lab",
  "Mistake Practice": "mistake_repair_plan",
  "All mistakes practiced!": "done",
  "Original:": "mistake_original",
  "Correct form:": "correct",
  "Type the correct form": "type_answer",
  "Submit": "send",
  "Message": "input",
  "Schedule Review": "tutor_review",
  "When do you want to review this lesson?": "tutor_srs",
  "Roleplay Scenarios": "roleplay",
  "Your answer": "lesson_answer",
  "Setup": "auth_finish_account",
  "Choose your learning language": "learning_language",
  "Your level": "level_label",
  "Start Learning": "start",
  "Choose scenario": "roleplay",
  "Reply": "send",
  "Repeat:": "shadowing_target",
  "Weak words:": "weak_words",
  "Spell the word": "start_spelling",
  "Hint": "spelling_hint",
  "Give up": "give_up",
  "Add phrase": "add_sample",
  "Phrase": "phrasebook_phrase_placeholder",
  "Translation": "phrasebook_note_placeholder",
  "Note": "phrasebook",
  "Saved phrases": "saved_phrases",
  "No phrases yet. Add one above to keep it on all your devices.": "empty_phrasebook",
  "No offline cards yet. Save phrases to learn offline.": "offline_short_hint",
  "No mistakes yet. Practice to fill this list.": "no_mistakes_loaded",
  "Test Complete!": "done",
  "All trophies": "awards",
  "Clear all mistakes?": "confirm_clear_mistakes",
  "Practice": "practice",
  "Clear all": "clear",
  "Image selected": "image_message",
  "Pick a photo containing text to translate": "image_message",
  "Choose image": "attach_screenshot",
  "Target:": "target_language",
  "Theme": "toggle_theme",
  "Light": "toggle_theme",
  "Dark": "toggle_theme",
  "System": "toggle_theme",
  "Interface Language": "interface_language",
  "Learning Language": "learning_language",
  "Level": "level_label",
  "Learning Focus": "learning_focus",
  "e.g. travel, business, exam prep": "learning_focus_placeholder",
  "Save Focus": "save_settings",
  "Activation Key": "activation_key",
  "Change Password": "change_password",
  "Confirm new password": "confirm_password",
  "Link Telegram to share progress between bot and app.": "telegram_link_hint",
  "Open": "open_telegram",
  "Log Out": "logout",
  "Follow NERIVA — short lessons, updates, and product tips.": "poliglot_social_body",
  "Your code": "referral_code",
  "Invitees": "referral_invitees",
  "Describe the issue": "problem_description",
  "Send Report": "send_report",
  "Report sent. Thank you!": "bug_report_sent",
  "AI Tutor": "ai_tutor",
  "Learn Words": "learn_words",
  "New words": "tutor_words",
  "Correct!": "correct",
  "Report a word mistake": "ai_tutor_word_report_title",
  "Send the corrected word and translation.": "ai_tutor_word_report_body",
  "Correct word": "ai_tutor_word_report_word_label",
  "Correct translation": "ai_tutor_word_report_translation_label",
  "Comment (optional)": "ai_tutor_word_report_comment_label",
  "Send report": "ai_tutor_word_report_submit",
  "Cancel": "cancel",
};

function decodeXml(value) {
  return value
    .replaceAll("&quot;", '"')
    .replaceAll("&apos;", "'")
    .replaceAll("&lt;", "<")
    .replaceAll("&gt;", ">")
    .replaceAll("&amp;", "&")
    .replaceAll("\\'", "'")
    .replaceAll("\\\\", "\\")
    .replaceAll("\\n", "\n");
}

for (const match of fs.readFileSync(stringsPath, "utf8").matchAll(/<string name="([a-z0-9_]+)">([\s\S]*?)<\/string>/g)) {
  const [_, key, rawValue] = match;
  const value = decodeXml(rawValue);
  const keys = valueToKeys.get(value) || [];
  keys.push(key);
  valueToKeys.set(value, keys);
}

function visit(directory) {
  for (const entry of fs.readdirSync(directory, { withFileTypes: true })) {
    const file = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      visit(file);
      continue;
    }
    if (!entry.name.endsWith(".kt")) continue;

    const original = fs.readFileSync(file, "utf8");
    let changed = false;
    const source = original.replace(/\bText\(\s*"((?:\\.|[^"\\])*)"/g, (full, rawValue) => {
      let value;
      try {
        value = JSON.parse(`"${rawValue}"`);
      } catch {
        return full;
      }
      if (value.includes("$") || /^X+(?:-X+)+$/.test(value) || value.startsWith("Demo Mode") || value.startsWith("Telegram login:")) return full;
      const mappedKey = literalToResource[value];
      const keys = valueToKeys.get(value);
      const key = mappedKey || (keys?.length === 1 ? keys[0] : null);
      if (!key) return full;
      changed = true;
      return `Text(stringResource(R.string.${key})`;
    });
    if (!changed) continue;

    let localized = source;
    if (!localized.includes("import androidx.compose.ui.res.stringResource")) {
      localized = localized.replace(/(package [^\n]+\n)/, "$1\nimport androidx.compose.ui.res.stringResource\n");
    }
    if (!localized.includes("import ru.neriva.app.R")) {
      localized = localized.replace(/(package [^\n]+\n)/, "$1\nimport ru.neriva.app.R\n");
    }
    fs.writeFileSync(file, localized);
    console.log(path.relative(process.cwd(), file));
  }
}

visit(uiRoot);
