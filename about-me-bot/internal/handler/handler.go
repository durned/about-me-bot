package handler

import (
	"fmt"

	"about-me-bot/external/holidays"
	weather "about-me-bot/external/openweather"
	l "about-me-bot/internal/logger"
	"about-me-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	Markdown string = "Markdown"
)

// all formatting must be done with HTMl in text
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
		l.SimpleLogger.Info(fmt.Sprintf("[@%s]: %s", update.Message.From.UserName, update.Message.Text))
		return handleMessage(update.Message)
	}

	return nil
}

func handleMessage(m *tgbotapi.Message) tgbotapi.Chattable {
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

	return CreateMsg(m, models.WrongFmt, false, dummy, false)
}

func handleCommand(m *tgbotapi.Message) tgbotapi.Chattable {
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
	}
	return CreateMsg(m, models.UnknownCommand, false, &tgbotapi.InlineKeyboardMarkup{}, false)
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
	}

	_, exi := holidays.SupportedCountries[query.Data]
	fail := CreateMsg(query.Message, models.APIFail, false, &tgbotapi.InlineKeyboardMarkup{}, false) // nil or a message with an empty string as text will cause a panic

	if exi {
		hDay, err := holidays.Holiday(query.Data)
		if err != nil {
			l.SimpleLogger.Error(err.Error())
			return fail
		}
		return CreateMsg(query.Message, hDay, true, &models.BackToHolidaysMarkup, true)
	}

	return fail
}
