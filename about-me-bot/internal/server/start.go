package server

import (
	"bufio"
	"context"
	"log"
	"os"

	cfg "about-me-bot/internal/config"
	"about-me-bot/internal/handler"

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
				log.Printf("error sending a message: %s", err.Error())
			}
		}
	}
}

func Run() {
	if !(len(cfg.BotCfg.Token) > 0) {
		log.Fatal("bot token has not been initialized")
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
		log.Fatalf("error creating a new bot API: %s", newBotErr.Error())
	}

	u.Timeout = cfg.BotCfg.UpdTimeout
	updates, err := myBot.API.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("error getting updates: %v", err)
	}

	// myBot.API.Debug = false; goes to false by default
	go myBot.receiveUpdates(ctx, updates)
	log.Printf("Authorized as @%s", myBot.API.Self.UserName)

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	log.Print("Pressed Enter key. Exiting.")
	cancel()
}
