import { test, expect } from '@playwright/test';

test.beforeEach(async ({ page }) => { await page.goto('/'); });

test('shows the workspace, configured providers, and idle agents', async ({ page }) => {
  await expect(page.getByRole('heading', { name: /Пусть агенты/ })).toBeVisible();
  await expect(page.locator('.profile-card')).toContainText('codex-orchestrator');
  await expect(page.locator('.agent-card[data-agent-id="free_qwen_coder"]')).toContainText('OpenRouter');
  await expect(page.locator('.agent-card[data-agent-id="local_qwen"]')).toContainText('Ollama · локально');
  await expect(page.getByText('В ожидании задачи')).toBeVisible();
});

test('keeps the agent flow inside laptop and mobile viewport widths', async ({ page }) => {
  for (const viewport of [{ width: 872, height: 744 }, { width: 390, height: 844 }]) {
    await page.setViewportSize(viewport);
    const metrics = await page.evaluate(() => ({ viewport: window.innerWidth, document: document.documentElement.scrollWidth, flow: document.querySelector('#flow-panel').scrollWidth }));
    expect(metrics.document, `document overflow at ${viewport.width}px`).toBeLessThanOrEqual(metrics.viewport);
    expect(metrics.flow).toBeLessThanOrEqual(metrics.viewport);
    await expect(page.locator('.agent-card[data-agent-id="local_qwen"]')).toBeVisible();
    await expect(page.locator('.agent-card[data-agent-id="local_qwen"]')).toContainText('Ожидает');
  }
});

test('rejects an invalid project path without starting a task', async ({ page }) => {
  await page.getByLabel('Опиши задачу').fill('Проверь запуск приложения');
  await page.getByRole('button', { name: 'Изменить' }).click();
  await page.getByLabel('Папка проекта').fill('Z:\\folder-that-does-not-exist');
  await page.getByRole('button', { name: 'Запустить агентов' }).click();
  await expect(page.getByRole('alert')).toContainText('Папка проекта не найдена');
  await expect(page.locator('#session-status')).toHaveText('Готов к работе');
});

test('starts a task, streams a real-time subagent, and shows the final result', async ({ page }) => {
  await page.getByLabel('Опиши задачу').fill('Проверь регистрацию и предложи исправление');
  await page.getByRole('button', { name: 'Запустить агентов' }).click();
  await expect(page.locator('#session-status')).toHaveText('Идёт работа', { timeout: 5_000 });
  await expect(page.locator('.activity-row').filter({ hasText: 'В работе · Запуск субагента' })).toBeVisible();
  await expect(page.locator('.agent-card[data-agent-id="free_qwen_coder"]')).toHaveClass(/running/);
  await expect(page.locator('#answer-text')).toContainText('Qwen Coder проверил изменения', { timeout: 5_000 });
  await expect(page.locator('#session-status')).toHaveText('Задача выполнена');
  await expect(page.locator('.agent-card[data-agent-id="free_qwen_coder"]')).toHaveClass(/done/);
});

test('stops an active turn and can reset the session', async ({ page }) => {
  await page.getByLabel('Опиши задачу').fill('hold: проверь вручную');
  await page.getByRole('button', { name: 'Запустить агентов' }).click();
  await expect(page.locator('#stop-button')).toBeVisible();
  await page.getByRole('button', { name: 'Остановить задачу' }).click();
  await expect(page.locator('#session-status')).toHaveText('Задача остановлена');
  await page.getByRole('button', { name: 'Сбросить сессию' }).click();
  await expect(page.locator('#session-status')).toHaveText('Готов к работе');
  await expect(page.getByLabel('Опиши задачу')).toHaveValue('');
});
