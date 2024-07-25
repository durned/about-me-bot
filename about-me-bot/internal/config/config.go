package config

import (
	"context"

	l "about-me-bot/internal/logger"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Token      string `env:"TELEGRAM_SECRET"`
	UpdOffset  int    `env:"UPDATE_OFFSET"`
	UpdTimeout int    `env:"UPDATE_TIMEOUT"`
}

var BotCfg Config

func Load() {
	ctx := context.Background()

	BotCfg = Config{}

	if err := godotenv.Load(); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error loading environmental values")
	}

	if errParse := env.Parse(&BotCfg); errParse != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error parsing environmental values into a struct.")
	}
}
