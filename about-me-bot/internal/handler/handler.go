package handler

import (
	"log"

	"about-me-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	Markdown string = "Markdown"
)

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
		return handleCallbackQuery(update.CallbackQuery)
	case update.Message != nil:
		log.Printf("[@%s]: %s", update.Message.From.UserName, update.Message.Text)
		return handleMessage(update.Message)
	}

	return nil
}

func handleMessage(m *tgbotapi.Message) tgbotapi.Chattable {
	dummy := &tgbotapi.InlineKeyboardMarkup{}

	if len(m.Text) > 0 {
		if m.Text[0] == '/' {
			return handleCommand(m)
		} else {
			return CreateMsg(m, models.SendCommands, false, dummy, false)
		}
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
		return CreateMsg(m, models.AboutText, false, &models.AboutMarkup, true)
	case "/links":
		return CreateMsg(m, models.LinksText, false, &models.LinksMarkup, true)
	}
	return CreateMsg(m, models.UnknownCommand, false, &tgbotapi.InlineKeyboardMarkup{}, false)
}

func handleCallbackQuery(query *tgbotapi.CallbackQuery) tgbotapi.Chattable {
	switch query.Data {
	case "help":
		helpMessage := tgbotapi.Message{Chat: &tgbotapi.Chat{ID: query.Message.Chat.ID}, Text: "/help"}
		return handleCommand(&helpMessage)
	case "about":
		return CreateMsg(query.Message, models.AboutText, true, &models.AboutMarkup, true)
	case "links":
		return CreateMsg(query.Message, models.LinksText, true, &models.LinksMarkup, true)
	case "back":
		return CreateMsg(query.Message, models.HelpText, true, &models.HelpMarkup, true)
	}

	return nil
}
