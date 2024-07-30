package main

import (
	"about-me-bot/internal/config"
	"about-me-bot/internal/server"
)

func main() {
	config.Load()
	server.Run()
}
