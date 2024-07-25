package server

import (
	"bufio"
	"context"
	"fmt"
	"os"

	cfg "about-me-bot/internal/config"
	"about-me-bot/internal/handler"
	l "about-me-bot/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	Bot struct {
		API *tgbotapi.BotAPI
	}
)

func (b *Bot) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive an update from updates channel and then handle it
		case update := <-updates:
			if _, err := b.API.Send(handler.HandleUpdate(update)); err != nil {
				l.SimpleLogger.Error(fmt.Sprintf("could not send a message to @%s: %s", update.Message.From.UserName, err.Error()))
			} else {
				l.SimpleLogger.Info("successfully handled an update")
			}
		}
	}
}

func Run() {
	if !(len(cfg.BotCfg.Token) > 0) {
		l.SimpleLogger.Log(context.Background(), l.LevelFatal, "bot token has not been initialized")
	}

	var (
		myBot     Bot
		newBotErr error
		u         tgbotapi.UpdateConfig = tgbotapi.NewUpdate(cfg.BotCfg.UpdOffset)
		ctx                             = context.Background()
	)
	ctx, cancel := context.WithCancel(ctx)

	myBot.API, newBotErr = tgbotapi.NewBotAPI(cfg.BotCfg.Token)
	if newBotErr != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, fmt.Sprintf("error creating a new bot API: %s", newBotErr.Error()))
	}

	u.Timeout = cfg.BotCfg.UpdTimeout
	updates, err := myBot.API.GetUpdatesChan(u)
	if err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, fmt.Sprintf("error getting updates: %s", err.Error()))
	}

	// myBot.API.Debug = false; goes to false by default
	go myBot.receiveUpdates(ctx, updates)
	l.SimpleLogger.Info(fmt.Sprintf("Authorized as @%s", myBot.API.Self.UserName))

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	l.SimpleLogger.Info("Pressed Enter key. Exiting.")
	cancel()
}
