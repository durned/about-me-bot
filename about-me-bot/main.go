package main

import (
	"context"

	"about-me-bot/internal/config"
	l "about-me-bot/internal/logger"
	"about-me-bot/internal/server"
)

func main() {
	config.Load()

	l.SimpleLogger.Debug("I am a debug message")
	l.SimpleLogger.Info("I am an info message")
	l.SimpleLogger.Warn("I am a warning message")
	l.SimpleLogger.Error("I am an error message\n")

	server.Run()
	l.SimpleLogger.Log(context.Background(), l.LevelFatal, "I am a fatal message. Bye!")
}
