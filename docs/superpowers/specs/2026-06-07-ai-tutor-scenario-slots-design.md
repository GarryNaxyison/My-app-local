# AI Tutor Scenario Slots Design

## Goal

Перевести AI Tutor с генерации заданий из первых слов урока на сценарные шаблоны со слотами, чтобы мини-объяснение, выбор ответа и мини-диалог были связаны с реальной сценой.

## Problem

`course_tutor.go` сейчас смешивает словарь урока с обязательными деталями диалога:

- `tutorSuccessCriteria()` возвращает видимые служебные английские строки вроде `Use the lesson goal`, `Include ...`, `Variant focus ...`.
- `tutorScenarioChecks()` строит вопросы из `words[0..3]`, а не из задачи сцены.
- `tutorCheckDetailLine()` показывает формат `word: translation`, хотя пользователю нужна диалоговая деталь.
- `tutorAnswerVariants()` и `tutorScenarioDialogue()` используют первые слова массива, поэтому варианты ответа могут не отвечать на последнюю реплику репетитора.

Главный регрессионный пример: в кафе нерелевантное слово вроде `football` не должно становиться заказом, деталью или вариантом ответа.

## Recommended Approach

Использовать curated scenario templates with slots. Это предсказуемо, тестируемо и сохраняет текущую локальную генерацию уроков без LLM.

Rejected alternatives:

- LLM lesson generation: гибче, но нестабильно, дороже и сложнее покрывается регрессиями.
- Local patches in current functions: быстрее, но оставляет источник ошибки, потому что первые слова урока продолжают управлять сценой.

## Architecture

Добавить в lesson model поле `scenario_slots`, а в Go-коде ввести внутреннюю структуру сценария. Она должна описывать роль, ситуацию, обязательное действие, предмет, деталь, модельный ответ и последнюю реплику репетитора.

`scenario_slots` становится источником истины для:

- `success_criteria`;
- `mini_explanation`;
- `choice` and `checks`;
- `writing_task` and `writing_expected`;
- `answer_variants`;
- `listening_text`, `listening_expected`;
- `dialogue`, `dialogue_prompt`, `dialogue_goal`;
- `dialogue_variants`.

`words` остаётся словарным блоком урока. Слова делятся логически:

- scene-required words: используются в слотах и в диалоге;
- extra vocabulary words: показываются в словаре, но не обязательны для реплики;
- review words: не участвуют в сцене.

На первом шаге явное API-поле для этих категорий не требуется: достаточно не брать обязательные слоты из произвольного `words[0..3]`.

## Scenario Slots

Public JSON field:

```go
type tutorScenarioSlots struct {
    Role           string   `json:"role"`
    Situation      string   `json:"situation"`
    RequiredAction string   `json:"required_action"`
    Item           string   `json:"item"`
    Detail         string   `json:"detail"`
    Politeness     []string `json:"politeness,omitempty"`
    ModelAnswer    string   `json:"model_answer"`
}
```

Internal scenario template can include extra teaching fields that do not need JSON exposure:

```go
type tutorScenarioTemplate struct {
    TopicCode          string
    Role              string
    SituationRU       string
    SituationEN       string
    RequiredAction    string
    Item              string
    Detail            string
    Politeness        []string
    ModelAnswer       string
    LastTutorLine     string
    MiniExplanationRU string
    MiniExplanationEN string
}
```

## Curated Templates

Required templates for this change:

- `food`: cafe order. Slots: `ask for the menu`, `coffee`, `for breakfast`, politeness `Could I...` and `I'd like..., please`.
- `hotel`: reception. Slots: confirm reservation, room, passport or one-night detail.
- `doctor`: clinic. Slots: symptom, doctor/help, today or morning.
- `travel`: station/airport. Slots: ticket, destination or time, tomorrow.
- `services`: service desk. Slots: help/problem, number/passport, today or next step.

Other topics may use a generic service-like fallback, but their required slots must still come from the template, not from arbitrary vocabulary order.

## Cafe Expected Output

For `food` and Russian interface, mini explanation should read like:

```text
В кафе нужен короткий порядок: попросить меню -> заказать напиток -> уточнить время еды.

Паттерн:
Could I see the menu?
I'd like coffee for breakfast, please.

Скажите одну живую реплику: вежливое начало + заказ + одна деталь.
```

It must not include:

- `Use the lesson goal`;
- `Include cup`;
- `Variant focus`;
- unrelated words such as `football` as an order or detail.

Choice question 1 should compare real replies:

- too short: `Coffee.`;
- wrong scene: `I need a ticket.`;
- incomplete but plausible: `Could I see the menu, please?`;
- correct: `Could I see the menu and have coffee for breakfast, please?`.

Choice question 2 should ask what the order clarifies:

- correct: `for breakfast`;
- wrong scene: `at the airport`;
- unnecessary: `my passport`;
- unrelated: `football`.

Mini-dialogue should answer the last tutor line:

```text
Tutor: Good morning. Would you like breakfast or dinner?
Your turn:
A1: Coffee for breakfast, please.
A2: I'd like coffee for breakfast, please.
B1: Could I see the menu first? I'd like coffee for breakfast, please.
Avoid: cup / football.
```

Feedback should be specific enough to name the item and detail:

```text
Хорошо: есть заказ coffee и деталь for breakfast.
Добавьте вежливое начало, если хотите звучать естественнее.
```

## Data Flow

`buildTutorLessonForSequence()` should:

1. Select topic, function, variant and vocabulary as today.
2. Resolve a scenario template by topic and interface language.
3. Build checks, explanation, dialogue and variants from the template.
4. Keep vocabulary selection independent from required scenario slots.
5. Include `ScenarioSlots` in the returned lesson JSON.

For non-English learning languages, current fallback may remain simpler, but visible tasks should still avoid English service strings and random `word / word` answers where a scenario template exists.

## UI Impact

`web-react` already renders lesson fields from `/api/tutor/start`. No new UI layout is required for the first implementation.

If `scenario_slots` is added to TypeScript types, it is optional and can remain hidden. The Playwright test should verify rendered user-facing content, not internal JSON.

## Tests

Go tests:

- `TestTutorCafeScenarioUsesSlotsNotRandomWords`: with a `food` lesson or direct helper, `football` must not appear as order/model answer/dialogue variant while `coffee` and `for breakfast` do.
- `TestTutorRussianMiniExplanationDoesNotLeakServiceCriteria`: RU mini explanation and success criteria must not contain `Use the lesson goal`, `Include`, or `Variant focus`.
- `TestTutorScenarioChoicesAreDialogueReplies`: choice options must contain spoken replies and scenario details, not `word: translation` or `word / word` as correct answers.

Playwright test:

- Mock AI Tutor cafe lesson with `football` in `words`.
- Open tutor view.
- Step through words and explanation.
- Assert mini explanation contains cafe pattern and does not contain service strings.
- Assert choice shows the logical cafe reply and detail question.
- Assert dialogue variants answer the last tutor line and do not use `football` as the order.

## Completion Criteria

- Cafe AI Tutor no longer generates `football` as a cafe order, detail or model answer.
- RU mini explanation and criteria do not leak English internal service strings.
- Choice answers are conversational replies or meaningful details.
- Mini-dialogue variants are tied to the last tutor line.
- Go tests and relevant Playwright test pass.
