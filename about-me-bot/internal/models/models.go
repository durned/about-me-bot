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
			tgbotapi.NewInlineKeyboardButtonData("Weather", "weather"),
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
	WeatherMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Status", "substatus"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<< Back", "back"),
		),
	)
	SusbcriptionStatusMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Check", "timezonecheck"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<< Back", "weatherBack"),
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
Press the *Help* button (or type _/help_) to proceed.`

	HelpText string = `_These_ down here are also _commands_, but I'd press *buttons*.
💬 The chatting part is sending _messages_ with some *location*.
Press the *Weather* button (or type _/weather_) to find out.`

	WeatherText string = `Send a location ➡️ bot will send a *weather forecast* for it.
Preferably, this location should be some _city_.\*

_Right now, the more populated locations will be chosen._
(i.e. 🗼 *Paris in France*, not Paris in Texas, USA.)

_*You can also make _some_ mistakes in the spelling_

Instead of doing it manually, you can *subscribe* to this forecast.
*Before doing so, click the "Status" button!* ⬇️

Type _/subscribe_ to start. (or just try it out by sending some cities)
Typing _/unsubscribe_ will allow deletion of your subscriptions. 
Don't worry, the bot will give you feedback as you go!`

	SusbcriptionStatusText string = `‼️ *IMPORTANT:* For the bot to send you regular forecasts, you will have to share the *approximate location* you are in, so bot would memorize *your timezone*.
_Approximate_ means that you can just send the city or a place somewhere on the map 🗺️, which is also in your same *timezone*, you are free to choose.
This is *private* data, so you have to do it *manually*. Sending a location *attachment* 📎 will hence always *override* your timezone in bot's knowledge.
Pressing the *button* below will send your status: will you receive forecasts or not.`

	BeginSubscriptionText string = `You want to subscribe to receive weather forecasts, nice!

You have to send "_/subscribe TIME LOCATION_" to me, the bot 🤖, where:
_TIME_ is the HH:MM time you'd like to receive a regular forecast. 
A very good time would be one written like this 04*:*02, with ':' and zeros! 
_LOCATION_ is, well, a location in which you'd like to know the weather.

Again, the bot will coordinate you on your way. Try it out!`

	AboutText string = `My name's Valera. I am student, 19 y.o., born in Riga, Latvia 🇱🇻.
A *soon-to-be* Second Year CS faculty 👨‍💻 student at Radboud University in The Netherlands 🇳🇱, and right now I live there most of the year.
I spend my time doing sports 🏋, and, of course, programming.`
	LinksText    string = "_These_ are some of my socials:\n\n• *Github:* ```https://github.com/durned```\n• *Telegram:* ```t.me/kozlikov```"
	HolidaysText string = "Press a *flag* of the country you want to see today's holiday in!"

	WrongFmt       string = "I don't support this type of communication, try text."
	UndefinedMsg   string = "❓ This isn't a _command_, but it is also not some *City*. Try again."
	UnknownCommand string = "I don't support this command."
	NoUserRecord   string = "❌ There are no current forecast subscriptions for this account."
	APIFail        string = "Something went wrong. Try again later."

	SubHandleArgsCount string = "Try again with *exactly* two arguments: time and location."
	SubHandleTimeFmt   string = "The time format seems to be incorrect. Check _/subscribe_ doc again."
)
