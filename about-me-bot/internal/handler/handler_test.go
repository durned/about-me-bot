package handler

import (
	"reflect"
	"testing"

	"about-me-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var chatID int64 = 13

func msgHelper(text string) *tgbotapi.Message {
	return &tgbotapi.Message{
		From: &tgbotapi.User{UserName: "kozlikov"},
		Chat: &tgbotapi.Chat{ID: chatID},
		Text: text,
	}
}

// Callbacks are not tested.

func TestHandleUpdate(t *testing.T) {
	dummy := &tgbotapi.InlineKeyboardMarkup{}

	type args struct {
		update tgbotapi.Update
	}

	tests := []struct {
		name string
		args args
		want tgbotapi.Chattable
	}{
		{
			"wrong fmt",
			args{tgbotapi.Update{
				Message: msgHelper(""),
			}},
			CreateMsg(msgHelper(""), models.WrongFmt, false, dummy, false),
		},
		{
			"send commands",
			args{tgbotapi.Update{
				Message: msgHelper("hey"),
			}},
			CreateMsg(msgHelper("hey"), models.SendCommands, false, dummy, false),
		},
		{
			"unknown command",
			args{tgbotapi.Update{
				Message: msgHelper("/tryme"),
			}},
			CreateMsg(msgHelper("/tryme"), models.UnknownCommand, false, dummy, false),
		},
		{
			"start",
			args{tgbotapi.Update{
				Message: msgHelper("/start"),
			}},
			CreateMsg(msgHelper("/start"), models.WelcomeText, false, &models.StartMarkup, true),
		},
		{
			"help",
			args{tgbotapi.Update{
				Message: msgHelper("/help"),
			}},
			CreateMsg(msgHelper("/help"), models.HelpText, false, &models.HelpMarkup, true),
		},
		{
			"about",
			args{tgbotapi.Update{
				Message: msgHelper("/about"),
			}},
			CreateMsg(msgHelper("/about"), models.AboutText, false, &models.AboutMarkup, true),
		},
		{
			"links",
			args{tgbotapi.Update{
				Message: msgHelper("/links"),
			}},
			CreateMsg(msgHelper("/links"), models.LinksText, false, &models.LinksMarkup, true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := HandleUpdate(tt.args.update); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
