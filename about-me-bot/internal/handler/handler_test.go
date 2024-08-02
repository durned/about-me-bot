package handler

import (
	"log"
	"reflect"
	"testing"
	"time"

	"about-me-bot/external/holidays"
	"about-me-bot/external/openweather"
	"about-me-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var chatID int64 = 13

func msgHelper(text string) *tgbotapi.Message {
	return &tgbotapi.Message{
		From:      &tgbotapi.User{UserName: "kozlikov"},
		Chat:      &tgbotapi.Chat{ID: chatID},
		Text:      text,
		MessageID: 1,
	}
}

func cQueryHelper(data string) *tgbotapi.CallbackQuery {
	return &tgbotapi.CallbackQuery{
		Data:    data,
		Message: msgHelper(""),
		From:    &tgbotapi.User{UserName: "kozlikov"},
	}
}

func TestHandleMessages(t *testing.T) {
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
			"send commands or cities",
			args{tgbotapi.Update{
				Message: msgHelper("["),
			}},
			CreateMsg(msgHelper("["), models.UndefinedMsg, false, dummy, false),
		},
		{
			"unknown command",
			args{tgbotapi.Update{
				Message: msgHelper("/tryme"),
			}},
			CreateMsg(msgHelper("/tryme"), models.UnknownCommand, false, dummy, false),
		},
		{
			"/start",
			args{tgbotapi.Update{
				Message: msgHelper("/start"),
			}},
			CreateMsg(msgHelper("/start"), models.WelcomeText, false, &models.StartMarkup, true),
		},
		{
			"/help",
			args{tgbotapi.Update{
				Message: msgHelper("/help"),
			}},
			CreateMsg(msgHelper("/help"), models.HelpText, false, &models.HelpMarkup, true),
		},
		{
			"/about",
			args{tgbotapi.Update{
				Message: msgHelper("/about"),
			}},
			CreateMsg(msgHelper("/about"), models.AboutText, false, &models.BackToHelpMarkup, true),
		},
		{
			"/links",
			args{tgbotapi.Update{
				Message: msgHelper("/links"),
			}},
			CreateMsg(msgHelper("/links"), models.LinksText, false, &models.BackToHelpMarkup, true),
		},
		{
			"/holidays",
			args{tgbotapi.Update{
				Message: msgHelper("/holidays"),
			}},
			CreateMsg(msgHelper("/holidays"), models.HolidaysText, false, &models.HolidaysMarkup, true),
		},
	}

	openweather.GeocodeClient = openweather.NoMatchesFoundClient

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := HandleUpdate(tt.args.update); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleCallbackQueries(t *testing.T) {
	type args struct {
		update tgbotapi.Update
	}

	holidays.HolidaysClient = holidays.LatvianHolidaysClient
	LVHolidays, err := holidays.Holiday("LV")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second) // to avoid rate limits

	tests := []struct {
		name string
		args args
		want tgbotapi.Chattable
	}{
		{
			"help",
			args{
				tgbotapi.Update{
					CallbackQuery: cQueryHelper("help"),
				},
			},
			CreateMsg(msgHelper(""), models.HelpText, false, &models.HelpMarkup, true),
		},
		{
			"about",
			args{
				tgbotapi.Update{
					CallbackQuery: cQueryHelper("about"),
				},
			},
			CreateMsg(msgHelper(""), models.AboutText, true, &models.BackToHelpMarkup, true),
		},
		{
			"links",
			args{
				tgbotapi.Update{
					CallbackQuery: cQueryHelper("links"),
				},
			},
			CreateMsg(msgHelper(""), models.LinksText, true, &models.BackToHelpMarkup, true),
		},
		{
			"back",
			args{
				tgbotapi.Update{
					CallbackQuery: cQueryHelper("back"),
				},
			},
			CreateMsg(msgHelper(""), models.HelpText, true, &models.HelpMarkup, true),
		},
		{
			"holidays",
			args{
				tgbotapi.Update{
					CallbackQuery: cQueryHelper("holidays"),
				},
			},
			CreateMsg(msgHelper(""), models.HolidaysText, true, &models.HolidaysMarkup, true),
		},
		{
			"holidayBack",
			args{
				tgbotapi.Update{
					CallbackQuery: cQueryHelper("holidayBack"),
				},
			},
			CreateMsg(msgHelper(""), models.HolidaysText, true, &models.HolidaysMarkup, true),
		},
		{
			"Latvian holidays today",
			args{
				tgbotapi.Update{
					CallbackQuery: cQueryHelper("LV"),
				},
			},
			CreateMsg(msgHelper(""), LVHolidays, true, &models.BackToHolidaysMarkup, true),
		},
		{
			"definetly a country",
			args{
				tgbotapi.Update{
					CallbackQuery: cQueryHelper("COUNTRY"),
				},
			},
			CreateMsg(msgHelper(""), models.APIFail, false, &tgbotapi.InlineKeyboardMarkup{}, false),
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
