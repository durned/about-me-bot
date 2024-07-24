package config

import (
	"log"

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
	BotCfg = Config{}

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading environmental values")
	}

	if errParse := env.Parse(&BotCfg); errParse != nil {
		log.Fatalf("error parsing environmental values into a struct.")
	}
}
