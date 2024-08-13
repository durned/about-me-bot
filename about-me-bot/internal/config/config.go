package config

import (
	"context"
	"os"

	l "about-me-bot/internal/logger"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	TgBot    TgBot
	Db       Database
	Holidays Holidays
	Weather  WeatherForecast

	Ctx context.Context
}

type TgBot struct {
	Mode       string `env:"MODE"`
	Token      string `env:"TELEGRAM_BOT_TOKEN"`
	UpdOffset  int    `env:"TGBOT_UPDATE_OFFSET"`
	UpdTimeout int    `env:"TGBOT_UPDATE_TIMEOUT"`
}

type Database struct {
	URI  string
	Port string
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
	BotCfg      TgBot
	DbCfg       Database
	HolidaysCfg Holidays
	WeatherCfg  WeatherForecast
	Global      Config
)

func Load(ctx context.Context) {
	if err := godotenv.Load(); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error loading env vars")
	}

	if err := env.Parse(&BotCfg); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error parsing env vars into the TgBot struct")
	}

	if DbCfg.Port = os.Getenv("DB_PORT"); DbCfg.Port == "" {
		l.SimpleLogger.Info("DB_PORT env var not provided, defaulting to 27017")
		DbCfg.Port = "27017"
	}

	// makes it possible to set docker provided env vars
	if DbCfg.URI = os.Getenv("DB_URI"); DbCfg.URI == "" {
		l.SimpleLogger.Info("DB_URI env var not provided, defaulting to localholst:PORT")
		DbCfg.URI = "mongodb://localhost:" + DbCfg.Port
	}

	if err := env.Parse(&HolidaysCfg); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error parsing env vars into the Holidays struct")
	}
	HolidaysCfg.Endpoint = HolidaysCfg.APIUrl + HolidaysCfg.APIKey

	if err := env.Parse(&WeatherCfg); err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error parsing environmental variables into the WeatherForecast struct")
	}

	Global = Config{BotCfg, DbCfg, HolidaysCfg, WeatherCfg, ctx}

	if !(len(Global.TgBot.Mode) > 0) || !(len(BotCfg.Token) > 0) ||
		!(len(Global.Holidays.APIUrl) > 0) || !(len(Global.Holidays.APIKey) > 0) ||
		!(len(Global.Weather.GeocodingAPIUrl) > 0) || !(len(Global.Weather.ForecastAPIUrl) > 0) || !(len(Global.Weather.IconsUrl) > 0) || !(len(Global.Weather.APIKey) > 0) {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "error parsing env vars: initialization went wrong or some fields were deliberately left empty")
	}

	switch Global.TgBot.Mode {
	case "worker", "server":
		break
	default:
		l.SimpleLogger.Log(ctx, l.LevelFatal, "undefined bot mode (neither 'worker' nor 'server')")
	}
}
