package models

import (
	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	StartMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Help", "help"),
		),
	)
	HelpMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("About", "about"),
			tgbotapi.NewInlineKeyboardButtonData("Holidays", "holidays"),
			tgbotapi.NewInlineKeyboardButtonData("Links", "links"),
		),
	)
	BackToHelpMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<< Back", "back"),
		),
	)
	BackToHolidaysMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<< Back", "holidayBack"),
		),
	)
	HolidaysMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇱🇻", "LV"),
			tgbotapi.NewInlineKeyboardButtonData("🇳🇱", "NL"),
			tgbotapi.NewInlineKeyboardButtonData("🇺🇦", "UA"),
			tgbotapi.NewInlineKeyboardButtonData("🇧🇷", "BR"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇫🇯", "FJ"),
			tgbotapi.NewInlineKeyboardButtonData("🇯🇵", "JP"),
			tgbotapi.NewInlineKeyboardButtonData("🇨🇺", "CU"),
			tgbotapi.NewInlineKeyboardButtonData("🇹🇩", "TD"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<< Back", "back"),
		),
	)

	Cyan     = color.New(color.FgCyan).SprintFunc()
	Yellow   = color.New(color.FgYellow).SprintFunc()
	Red      = color.New(color.FgRed).SprintFunc()
	DarkRed  = color.New(color.FgHiRed).SprintFunc()
	DbgColor = color.New(color.FgMagenta).SprintFunc()
)

// Responses
const (
	WelcomeText  string = "*Welcome!*\n_This_ is a bot which works as an interactive mini About me.\nAlso, it shows Holidays in some countries for a *very good* reason.\n"
	HelpText     string = "_These_ are also commands. But I would press the *buttons* instead."
	AboutText    string = "My name's Valera. I am student, 19 y.o., born in Riga, Latvia.\nA *soon-to-be* Second Year CS faculty student at Radboud University in The Netherlands, and right now I live there most of the year.\nI spend my time doing sports, and, of course, programming."
	LinksText    string = "_These_ are some of my socials:\n\n• *Github:* ```https://github.com/durned```\n• *Telegram:* ```t.me/kozlikov```"
	HolidaysText string = "Press a *flag* of the country you want to see today's holiday in!"

	WrongFmt          string = "I don't support this type of communication, try text."
	SendCommands      string = "I only support commands, type _/help_ to see them"
	UnknownCommand    string = "I don't support this command."
	CallbackQueryFail string = "Something went wrong. Try again later."
)
