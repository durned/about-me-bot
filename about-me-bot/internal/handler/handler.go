package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	db "about-me-bot/database"
	"about-me-bot/external/holidays"
	weather "about-me-bot/external/openweather"
	cfg "about-me-bot/internal/config"
	l "about-me-bot/internal/logger"
	"about-me-bot/internal/models"
	"about-me-bot/internal/timezones"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	Markdown string = "Markdown"
)

// HTMLImgAppender all formatting must be done with HTMl in text
func HTMLImgAppender(m *tgbotapi.Message, text string, url string) tgbotapi.Chattable {
	text += fmt.Sprint("<a href=\"", url, "\">", "\u200e", "</a>")

	msg := tgbotapi.NewMessage(m.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}

func CreateMsg(m *tgbotapi.Message, text string, edit bool, markUp *tgbotapi.InlineKeyboardMarkup, applyMarkup bool) tgbotapi.Chattable {
	if edit {
		msg := tgbotapi.NewEditMessageText(m.Chat.ID, m.MessageID, text)
		msg.ParseMode = Markdown
		if applyMarkup {
			msg.ReplyMarkup = markUp
		}

		return msg
	}

	msg := tgbotapi.NewMessage(m.Chat.ID, text)
	msg.ParseMode = Markdown
	if applyMarkup {
		msg.ReplyMarkup = markUp
	}

	return msg
}

func HandleUpdate(update tgbotapi.Update) tgbotapi.Chattable {
	switch {
	case update.CallbackQuery != nil:
		l.SimpleLogger.Info(fmt.Sprintf("[@%s]: pressed %s button", update.CallbackQuery.From.UserName, update.CallbackQuery.Data))
		return handleCallbackQuery(update.CallbackQuery)
	case update.Message != nil:
		msgNotification := fmt.Sprintf("[@%s]: %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.Location != nil {
			loc := update.Message.Location
			msgNotification += fmt.Sprintf("%.2f %.2f", loc.Latitude, loc.Longitude)
		}
		l.SimpleLogger.Info(msgNotification)
		return HandleMessage(update.Message)
	}
	return nil
}

func HandleMessage(m *tgbotapi.Message) tgbotapi.Chattable {
	dummy := &tgbotapi.InlineKeyboardMarkup{}

	if len(m.Text) > 0 {
		if m.Text[0] == '/' {
			return handleCommand(m)
		}

		coords, err := weather.Geocode(m.Text)
		switch err {
		case weather.ErrNoMatchesFound:
			return CreateMsg(m, models.UndefinedMsg, false, dummy, false)
		case nil:
			break
		default:
			return CreateMsg(m, models.APIFail, false, dummy, false)
		}

		wForecast, wIconUrl, err := weather.Forecast(coords)
		switch err {
		case weather.ErrNoData:
			return CreateMsg(m, weather.ErrNoData.Error(), false, dummy, false)
		case nil:
			break
		default:
			return CreateMsg(m, models.APIFail, false, dummy, false)
		}

		return HTMLImgAppender(m, wForecast, wIconUrl)
	}

	if m.Location != nil {
		err := db.AddTimezone(cfg.Global.Ctx, m.Chat.ID, timezones.LatLonToTz(m.Location.Latitude, m.Location.Longitude))
		if err != nil {
			return CreateMsg(m, models.APIFail, false, dummy, false)
		}
		return CreateMsg(m, "Successfully updated your timezone", false, dummy, false)
	}

	return CreateMsg(m, models.WrongFmt, false, dummy, false)
}

func HandleSubscription(m *tgbotapi.Message) string {
	l.SimpleLogger.Info(fmt.Sprintf("Trying to add user's %v subscription", m.From.ID))

	split := strings.Fields(m.Text)
	if len(split) != 3 {
		return models.SubHandleArgsCount
	}

	if _, err := time.Parse("15:04", split[1]); err != nil {
		l.SimpleLogger.Error(err.Error())
		return models.SubHandleTimeFmt
	}

	res, err := weather.Geocode(split[2])
	switch err {
	case weather.ErrNoMatchesFound:
		return "The last argument does not seem to be a valid location, try again."
	case nil:
		break
	default:
		return models.APIFail
	}

	errAdd := db.Subscribe(cfg.Global.Ctx, m.Chat.ID, split[1], res.Name)
	switch errAdd {
	case db.ErrExiAlready:
		return "There seems to be such a record already, try a different one."
	case nil:
		return "Successfully subscribed!"
	default:
		return models.APIFail
	}
}

func handleUnsub(ChatID int64) (string, tgbotapi.InlineKeyboardMarkup, bool) {
	var apply bool
	dummy := tgbotapi.InlineKeyboardMarkup{}

	subs, err := db.GetSubscriptions(cfg.Global.Ctx, ChatID)
	switch err {
	case db.ErrNoUserRecord:
		return models.NoUserRecord, dummy, apply
	case nil:
		break
	default:
		return models.APIFail, dummy, apply
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	if len(subs) > 0 {
		apply = true
		for key, val := range subs {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(val, key.String())))
		}
	}

	return "Those are your current subscriptions.\nPress one to unsubscribe from it.", tgbotapi.NewInlineKeyboardMarkup(rows...), apply
}

func handleCommand(m *tgbotapi.Message) tgbotapi.Chattable {
	dummy := &tgbotapi.InlineKeyboardMarkup{}

	switch m.Text {
	case "/start":
		return CreateMsg(m, models.WelcomeText, false, &models.StartMarkup, true)
	case "/help":
		return CreateMsg(m, models.HelpText, false, &models.HelpMarkup, true)
	case "/about":
		return CreateMsg(m, models.AboutText, false, &models.BackToHelpMarkup, true)
	case "/links":
		return CreateMsg(m, models.LinksText, false, &models.BackToHelpMarkup, true)
	case "/holidays":
		return CreateMsg(m, models.HolidaysText, false, &models.HolidaysMarkup, true)
	case "/weather":
		return CreateMsg(m, models.WeatherText, false, &models.BackToHelpMarkup, true)
	case "/unsubscribe", "/unsub":
		text, markup, apply := handleUnsub(m.Chat.ID)
		return CreateMsg(m, text, false, &markup, apply)
	case "/subscribe", "/sub":
		return CreateMsg(m, models.BeginSubscriptionText, false, dummy, false)
	case "/status":
		return CreateMsg(m, models.SusbcriptionStatusText, false, &models.SusbcriptionStatusMarkup, true)
	case "/check":
		return handleCallbackQuery(&tgbotapi.CallbackQuery{Data: "timezonecheck", Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: m.Chat.ID}}})
	}

	if attrs := strings.Fields(m.Text); attrs[0] == "/subscribe" && len(attrs) > 1 {
		return CreateMsg(m, HandleSubscription(m), false, dummy, false)
	}

	return CreateMsg(m, models.UnknownCommand, false, dummy, false)
}

func handleCallbackQuery(query *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	switch query.Data {
	case "help":
		helpMessage := tgbotapi.Message{Chat: &tgbotapi.Chat{ID: query.Message.Chat.ID}, Text: "/help"}
		return handleCommand(&helpMessage)
	case "about":
		return CreateMsg(query.Message, models.AboutText, true, &models.BackToHelpMarkup, true)
	case "links":
		return CreateMsg(query.Message, models.LinksText, true, &models.BackToHelpMarkup, true)
	case "back":
		return CreateMsg(query.Message, models.HelpText, true, &models.HelpMarkup, true)
	case "holidays", "holidayBack":
		return CreateMsg(query.Message, models.HolidaysText, true, &models.HolidaysMarkup, true)
	case "weather", "weatherBack":
		return CreateMsg(query.Message, models.WeatherText, true, &models.WeatherMarkup, true)
	case "substatus":
		return CreateMsg(query.Message, models.SusbcriptionStatusText, true, &models.SusbcriptionStatusMarkup, true)
	case "timezonecheck":
		text := ""

		tzKnown, err := db.UserTimezoneKnown(cfg.Global.Ctx, query.Message.Chat.ID)

		if tzKnown {
			text = "✅ *Yes*, you *will*"
		} else {
			text = "⛔ *No*, you *won't*"
		}

		if err != nil {
			text = models.APIFail
		}

		return CreateMsg(query.Message, text, false, &tgbotapi.InlineKeyboardMarkup{}, false)
	}

	fail := CreateMsg(query.Message, models.APIFail, false, &tgbotapi.InlineKeyboardMarkup{}, false) // nil or a message with an empty string as text will cause a panic

	if len(query.Data) > 2 {
		unsubAttempt := strings.Split(query.Data, ",")
		if len(unsubAttempt) == 3 {
			num, err := strconv.Atoi(unsubAttempt[0])
			if err != nil {
				return fail
			}
			if errUnsub := db.Unsubscribe(cfg.Global.Ctx, int64(num), unsubAttempt[1], unsubAttempt[2]); errUnsub != nil {
				return fail
			}
			tempText, tempMarkup, apply := handleUnsub(query.Message.Chat.ID)
			return CreateMsg(query.Message, tempText, true, &tempMarkup, apply)
		}
	} else {
		_, exi := holidays.SupportedCountries[query.Data]

		if exi {
			hDay, err := holidays.Holiday(query.Data)
			if err != nil {
				l.SimpleLogger.Error(err.Error())
				return fail
			}
			return CreateMsg(query.Message, hDay, true, &models.BackToHolidaysMarkup, true)
		}
	}
	return fail
}
