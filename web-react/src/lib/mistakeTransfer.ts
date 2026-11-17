export type MistakeTransferTask = {
  prompt: { ru: string; en: string };
  hint: { ru: string; en: string };
  example: string;
  answers: string[];
};

const pastToBase: Record<string, string> = {
  went: "go",
  saw: "see",
  called: "call",
  bought: "buy",
  ate: "eat",
  had: "have",
  made: "make",
  took: "take",
};

function normalizeSentence(value: string): string {
  return value.trim().toLowerCase().replace(/[’‘]/gu, "'").replace(/,/g, " ").replace(/[.!?]+$/u, "").trim().replace(/\s+/gu, " ");
}

export function buildMistakeTransfer(original: string, correction: string, language: string): MistakeTransferTask | null {
  if (language !== "en") return null;
  const before = normalizeSentence(original);
  const after = normalizeSentence(correction);
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
      answers: [`i walked to the ${place} yesterday`, `yesterday i walked to the ${place}`, `i walked yesterday to the ${place}`],
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
      answers: ["she reads books every evening", "every evening she reads books", "she reads books each evening"],
    };
  }
  const afterWords = after.split(" ");
  const addedArticle = afterWords.length === before.split(" ").length + 1 && afterWords.some((word, index) =>
    (word === "a" || word === "an") && index < afterWords.length - 1 &&
    afterWords.filter((_, position) => position !== index).join(" ") === before,
  );
  if (addedArticle) {
    const useApple = /\bumbrella\b/u.test(after);
    const noun = useApple ? "apple" : "umbrella";
    return {
      prompt: {
        ru: useApple ? "Скажите по-английски, что вам нужно яблоко." : "Скажите по-английски, что вам нужен зонтик.",
        en: `Say that you need one ${noun}. Use an indefinite article.`,
      },
      hint: { ru: "Перед гласным звуком используйте an.", en: "Use an before a vowel sound." },
      example: `I need an ${noun}.`,
      answers: [`i need an ${noun}`, `i'll need an ${noun}`, `i will need an ${noun}`],
    };
  }
  if (/^(we|they|you) is\b/u.test(before) && after === before.replace(/^(we|they|you) is\b/u, "$1 are")) {
    const useOutside = /\bready\b/u.test(after);
    return {
      prompt: {
        ru: useOutside ? "Скажите по-английски, что они сейчас на улице." : "Скажите по-английски, что вы все готовы.",
        en: useOutside ? "Tell someone where the others are. Say they are outdoors, using outside." : "Tell someone that everyone in your group is ready. Start with we.",
      },
      hint: { ru: "После we и they используйте are.", en: "Use are after we and they." },
      example: useOutside ? "They are outside." : "We are ready.",
      answers: useOutside ? ["they are outside", "they're outside", "they are outside now", "they're outside now"] : ["we are ready", "we're ready", "we are all ready", "we're all ready"],
    };
  }
  const pastNegative = before.match(/^((?:i|you|he|she|we|they) )(?:didn't|did not) (\w+)(.*)$/u);
  if (pastNegative) {
    const [, subject, verb, rest] = pastNegative;
    const base = pastToBase[verb];
    if (base && after.replace(/\bdid not\b/u, "didn't") === `${subject}didn't ${base}${rest}`) {
      const useText = /\bcall mia\b/u.test(after);
      const action = useText ? "text leo" : "call mia";
      return {
        prompt: {
          ru: useText ? "Скажите по-английски, что вчера вы не написали Лео." : "Скажите по-английски, что вчера вы не позвонили Мие.",
          en: useText ? "Explain that yesterday you sent Leo no text. Use text in a past negative." : "Explain that yesterday you made no call to Mia. Use call in a past negative.",
        },
        hint: { ru: "После didn't ставьте глагол в начальной форме.", en: "Use the base form after didn’t." },
        example: useText ? "I didn’t text Leo yesterday." : "I didn’t call Mia yesterday.",
        answers: [
          `i didn't ${action} yesterday`, `i did not ${action} yesterday`,
          `yesterday i didn't ${action}`, `yesterday i did not ${action}`,
        ],
      };
    }
  }
  return null;
}

export function checkMistakeTransfer(task: MistakeTransferTask, answer: string): boolean {
  return task.answers.includes(normalizeSentence(answer));
}
