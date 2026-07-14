export function canonicalPhraseKey(value: string): string {
  return String(value || "")
    .normalize("NFKC")
    .replace(/\s+/gu, " ")
    .trim()
    .replace(/[.!?…]+$/gu, "")
    .trim()
    .toLocaleLowerCase();
}
