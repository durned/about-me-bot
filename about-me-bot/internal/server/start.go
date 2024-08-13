package server

import (
	"context"
	"fmt"

	cfg "about-me-bot/internal/config"
	"about-me-bot/internal/handler"
	l "about-me-bot/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	API *tgbotapi.BotAPI
}

func InitBot() Bot {
	var (
		myBot     Bot
		newBotErr error
	)

	myBot.API, newBotErr = tgbotapi.NewBotAPI(cfg.Global.TgBot.Token)
	if newBotErr != nil { // Context with cancel() will not affect this anyway
		l.SimpleLogger.Log(context.Background(), l.LevelFatal, fmt.Sprintf("error creating a new bot API: %s", newBotErr.Error()))
	}

	// myBot.API.Debug = false; goes to false by default
	return myBot
}

func (b *Bot) MySend(in tgbotapi.Chattable) (tgbotapi.Message, error) {
	return b.API.Send(in)
}

func (b *Bot) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive an update from updates channel and then handle it
		case update := <-updates:
			toSend := handler.HandleUpdate(update)
			if toSend == nil {
				l.SimpleLogger.Info("update was handled, but nothing was sent")
				break
			}

			if _, err := b.MySend(toSend); err != nil {
				l.SimpleLogger.Error(fmt.Sprintf("could not send a message to @%s: %s", update.Message.From.UserName, err.Error()))
			} else {
				l.SimpleLogger.Info("successfully handled an update")
			}
		}
	}
}

func Run(ctx context.Context) {
	myBot := InitBot()

	var u tgbotapi.UpdateConfig = tgbotapi.NewUpdate(cfg.Global.TgBot.UpdOffset)
	u.Timeout = cfg.Global.TgBot.UpdTimeout

	updates, err := myBot.API.GetUpdatesChan(u)
	if err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, fmt.Sprintf("error getting updates: %s", err.Error()))
	}

	go myBot.receiveUpdates(ctx, updates)
	l.SimpleLogger.Info(fmt.Sprintf("Authorized as @%s", myBot.API.Self.UserName))
}
