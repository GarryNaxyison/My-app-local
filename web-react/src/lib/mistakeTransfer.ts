export type MistakeTransferTask = {
  prompt: { ru: string; en: string };
  hint: { ru: string; en: string };
  example: string;
  answers: string[];
};

export function buildMistakeTransfer(original: string, correction: string, language: string): MistakeTransferTask | null {
  if (language !== "en") return null;
  const before = original.trim();
  const after = correction.trim();
  if (/^i go\b/i.test(before) && /^i went\b/i.test(after)) {
    const alreadyUsedStation = /\bstation\b/i.test(after);
    const place = alreadyUsedStation ? "park" : "station";
    return {
      prompt: {
        ru: `Расскажите по-английски, что вчера вы дошли пешком ${alreadyUsedStation ? "до парка" : "до вокзала"}.`,
        en: `Tell a friend you went on foot to the ${place} yesterday. Use walk.`,
      },
      hint: { ru: "Используйте глагол walk в прошедшем времени.", en: "Use the past tense of walk." },
      example: `I walked to the ${place} yesterday.`,
      answers: [`i walked to the ${place} yesterday`, `yesterday i walked to the ${place}`],
    };
  }
  if (/^she go\b/i.test(before) && /^she goes\b/i.test(after)) {
    return {
      prompt: {
        ru: "Расскажите по-английски, что она читает книги каждый вечер.",
        en: "Say that she spends every evening reading books. Use read in the present simple.",
      },
      hint: { ru: "Добавьте окончание к глаголу read после she.", en: "Add the ending to read after she." },
      example: "She reads books every evening.",
      answers: ["she reads books every evening", "every evening she reads books"],
    };
  }
  return null;
}

export function checkMistakeTransfer(task: MistakeTransferTask, answer: string): boolean {
  const normalized = answer.trim().toLowerCase().replace(/[,]/g, " ").replace(/[.!?]+$/u, "").trim().replace(/\s+/gu, " ");
  return task.answers.includes(normalized);
}
