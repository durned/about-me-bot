package config

import (
	"context"
	_ "log"

	l "about-me-bot/internal/logger"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Holidays struct {
	APIUrl   string `env:"HOLIDAYS_API_URL"`
	APIKey   string `env:"HOLIDAYS_API_KEY"`
	Endpoint string
}

type Config struct {
	Token      string `env:"TELEGRAM_BOT_TOKEN"`
	UpdOffset  int    `env:"UPDATE_OFFSET"`
	UpdTimeout int    `env:"UPDATE_TIMEOUT"`

	Holidays Holidays
}

var (
	HDayCfg Holidays
	BotCfg  Config
)

func Load() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error loading environmental values")
	}

	if err := env.Parse(&HDayCfg); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error loading environmental values into the Holidays struct.")
	}

	HDayCfg.Endpoint = HDayCfg.APIUrl + HDayCfg.APIKey

	if err := env.Parse(&BotCfg); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error loading environmental values into the Bot struct.")
	}

	BotCfg.Holidays = HDayCfg

	if !(len(BotCfg.Token) > 0) || !(len(BotCfg.Holidays.APIUrl) > 0) || !(len(BotCfg.Holidays.APIKey) > 0) {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error parsing environmental values.")
	}
}
