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
	WelcomeText string = `*Welcome!*
_This_ is a bot which works as an interactive mini *About me*.
Also, it shows *Holidays* in some countries for a *very good* reason.
*But that's not it!* You can also sort of *chat* with the bot. 
Click the _Help_ button to find out.`
	HelpText string = `_These_ down here are also _commands_, but I'd press *buttons*.
The chatting part is when you send _messages_ with some *location*.
The bot will then send you a *weather forecast* for it.
Preferably, this location should be some _city_, for it to make sense.
Right now, the more populated locations will be chosen,
i.e. *Paris in France*, not Paris in Texas, USA.
Try it out, the bot will give you feedback as you go!
\*You can also make _some_ mistakes in the spelling`
	AboutText string = `My name's Valera. I am student, 19 y.o., born in Riga, Latvia.
A *soon-to-be* Second Year CS faculty student at Radboud University in The Netherlands, and right now I live there most of the year.
I spend my time doing sports, and, of course, programming.`
	LinksText    string = "_These_ are some of my socials:\n\n• *Github:* ```https://github.com/durned```\n• *Telegram:* ```t.me/kozlikov```"
	HolidaysText string = "Press a *flag* of the country you want to see today's holiday in!"

	WrongFmt       string = "I don't support this type of communication, try text."
	UndefinedMsg   string = "This isn't a _command_, but it is also not some *City*. Try again."
	UnknownCommand string = "I don't support this command."
	APIFail        string = "Something went wrong. Try again later."
)
