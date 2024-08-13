package main

import (
	"context"
	"os"
	"os/signal"
	_ "strconv"
	"syscall"

	db "about-me-bot/database"
	cfg "about-me-bot/internal/config"
	l "about-me-bot/internal/logger"
	"about-me-bot/internal/server"
	"about-me-bot/internal/worker"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	cfg.Load(ctx)

	if err := db.Init(ctx, cfg.Global.Db.URI, "sub-db"); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, err.Error())
	}

	switch cfg.Global.TgBot.Mode {
	case "server":
		go server.Run(ctx)
	case "worker":
		go worker.Run(ctx)
	}

	inputChan := make(chan os.Signal, 1)
	signal.Notify(inputChan, syscall.SIGINT, syscall.SIGTERM)

	<-inputChan
	cancel()
}
