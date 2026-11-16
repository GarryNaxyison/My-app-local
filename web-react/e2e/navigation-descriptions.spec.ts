import { expect, test } from '@playwright/test';
import { appCopy, appLocaleCodes, appNavDescription } from '../src/lib/i18n';
import type { ViewId } from '../src/lib/types';

const views: Array<[ViewId, string]> = [['home', 'today'], ['lesson', 'new_lesson'], ['practice', 'practice'], ['words', 'learn_words'], ['vocabulary', 'vocabulary'], ['level', 'level_test'], ['awards', 'awards'], ['premium', 'premium'], ['settings', 'settings']];
const comparable = (text: string) => text.toLocaleLowerCase().replace(/[\p{P}\p{Z}\s]/gu, '');
test('catalogue descriptions never repeat headings or consist of punctuation', () => {
  for (const locale of appLocaleCodes) {
    const copy = (key: string, fallback: string) => appCopy(locale, key, fallback);
    for (const [view, labelKey] of views) {
      const description = appNavDescription(view, locale, '');
      if (!description) continue;
      expect(comparable(description), `${locale}.${view} meaningful`).not.toBe('');
      expect(comparable(description), `${locale}.${view} duplicate`).not.toBe(comparable(copy(labelKey, '')));
      expect(description, `${locale}.${view} placeholder`).not.toMatch(/[:：]\s*$/);
    }
  }
});
test('catalogue Russian descriptions explain the destination', () => {
  expect(appNavDescription('words', 'ru', '')).toBe('Новые слова и примеры');
  expect(appNavDescription('lesson', 'ru', '')).toBe('Короткое задание с проверкой');
  expect(appNavDescription('level', 'ru', '')).toBe('Проверка уровня языка');
  expect(appNavDescription('words', 'en', '')).toBe('New vocabulary');
  expect(appNavDescription('home', 'de', '')).toBe('');
  expect(appNavDescription('premium', 'tr', '')).toBe('');
});
