package main

import (
	"context"
	"log"
	"time"
)

type transitionBot struct {
	telegram  *telegramClient
	targetBot string
	offset    int64
}

func runTransitionBot(ctx context.Context, telegram *telegramClient, targetBot string) error {
	bot := &transitionBot{
		telegram:  telegram,
		targetBot: normalizeBotUsername(targetBot),
	}
	return bot.run(ctx)
}

func (b *transitionBot) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		updates, err := b.telegram.getUpdates(ctx, b.offset)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("transition bot getUpdates failed: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		for _, update := range updates {
			if update.UpdateID >= b.offset {
				b.offset = update.UpdateID + 1
			}
			if err := b.handleUpdate(ctx, update); err != nil {
				log.Printf("transition bot update %d failed: %v", update.UpdateID, err)
			}
		}
	}
}

func (b *transitionBot) handleUpdate(ctx context.Context, update telegramUpdate) error {
	if update.CallbackQuery != nil {
		_ = b.telegram.answerCallbackQuery(ctx, update.CallbackQuery.ID, "")
		if update.CallbackQuery.Message == nil {
			return nil
		}
		return b.sendTransition(ctx, update.CallbackQuery.Message.Chat.ID)
	}
	if update.Message == nil {
		return nil
	}
	return b.sendTransition(ctx, update.Message.Chat.ID)
}

func (b *transitionBot) sendTransition(ctx context.Context, chatID int64) error {
	if chatID == 0 {
		return nil
	}
	return b.telegram.sendInlineMessage(ctx, chatID, transitionBotText(b.targetBot), transitionBotInlineKeyboard(b.targetBot))
}
