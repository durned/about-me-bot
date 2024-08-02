package config

import (
	"context"

	l "about-me-bot/internal/logger"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Token      string `env:"TELEGRAM_BOT_TOKEN"`
	UpdOffset  int    `env:"TGBOT_UPDATE_OFFSET"`
	UpdTimeout int    `env:"TGBOT_UPDATE_TIMEOUT"`

	Holidays Holidays
	Weather  WeatherForecast
}

type Holidays struct {
	APIUrl   string `env:"HOLIDAYS_API_URL"`
	APIKey   string `env:"HOLIDAYS_API_KEY"`
	Endpoint string
}

type WeatherForecast struct {
	GeocodingAPIUrl string `env:"OPENWEATHER_GEOCODING"`
	ForecastAPIUrl  string `env:"OPENWEATHER_FORECAST"`
	IconsUrl        string `env:"OPENWEATHER_ICONS"`
	APIKey          string `env:"OPENWEATHER_API_KEY"`
}

var (
	HolidaysCfg Holidays
	WeatherCfg  WeatherForecast
	BotCfg      Config
)

func Load() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error loading environmental values")
	}

	if err := env.Parse(&HolidaysCfg); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error loading environmental values into the Holidays struct.")
	}
	HolidaysCfg.Endpoint = HolidaysCfg.APIUrl + HolidaysCfg.APIKey

	if err := env.Parse(&WeatherCfg); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error loading environmental values into the WeatherForecast struct.")
	}

	if err := env.Parse(&BotCfg); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error loading environmental values into the Bot struct.")
	}

	BotCfg.Holidays = HolidaysCfg
	BotCfg.Weather = WeatherCfg

	if !(len(BotCfg.Token) > 0) ||
		!(len(BotCfg.Holidays.APIUrl) > 0) || !(len(BotCfg.Holidays.APIKey) > 0) ||
		!(len(BotCfg.Weather.GeocodingAPIUrl) > 0) || !(len(BotCfg.Weather.ForecastAPIUrl) > 0) || !(len(BotCfg.Weather.IconsUrl) > 0) || !(len(BotCfg.Weather.APIKey) > 0) {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error parsing environmental values.")
	}
}
