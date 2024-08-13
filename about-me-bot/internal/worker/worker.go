package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	db "about-me-bot/database"
	"about-me-bot/internal/handler"
	l "about-me-bot/internal/logger"
	"about-me-bot/internal/server"
	"about-me-bot/internal/timezones"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	bot server.Bot
	wg  sync.WaitGroup
	mu  sync.Mutex
)

func Run(ctx context.Context) {
	utcLoc, _ := time.LoadLocation("UTC")

	bot = server.InitBot()

	list, err := db.GetForecastListUTC(ctx)
	if err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "couldn't receive fundamental data from db: "+err.Error())
	}

	select {
	case <-ctx.Done():
		break
	default:
		for i := range list {
			// less boilerplate typecasting this way
			// will get the same date but with different time
			trgtTime, err := timezones.HHMMStringToUTC(list[i].Time, utcLoc)
			if err != nil {
				l.SimpleLogger.Error("skipping a query: " + err.Error())
				continue
			}

			nowUTC := time.Now().UTC()
			if trgtTime.Before(nowUTC) {
				// then we would have to send it tomorrow
				trgtTime = trgtTime.Add(time.Hour * 24)
			}

			go func(ctx context.Context, duration time.Duration, chatID int64, city string) {
				defer wg.Done()

				l.SimpleLogger.Info("queued a query, will send in " + duration.String())
				select {
				case <-ctx.Done():
					break
				case <-time.After(duration):
					// to not overload the bot
					mu.Lock()
					defer mu.Unlock()

					_, err := bot.MySend(handler.HandleMessage(&tgbotapi.Message{Chat: &tgbotapi.Chat{ID: chatID}, Text: city}))
					if err != nil {
						l.SimpleLogger.Error(fmt.Sprintf("error sending regular forecast to chat [%d]: %s", chatID, err.Error()))
					}
					l.SimpleLogger.Info("sent a forecast to chat [%d]")
				}
			}(ctx, trgtTime.Sub(nowUTC), int64(list[i].ChatID), list[i].Location)
			wg.Add(1)
		}
		wg.Wait()
	}
}
