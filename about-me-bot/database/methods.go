package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	l "about-me-bot/internal/logger"
	"about-me-bot/internal/timezones"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func existsSubscription(ctx context.Context, sub Subscription) (bool, error) {
	var res Subscription
	filter := bson.D{{Key: "chat_id", Value: sub.ChatID}, {Key: "time", Value: sub.Time}, {Key: "location", Value: sub.Location}}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := SubscriptionCollection.FindOne(ctx, filter).Decode(&res)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else if err != nil {
		l.SimpleLogger.Error(err.Error())
		return false, err
	}

	return true, nil
}

func Subscribe(ctx context.Context, ChatID int64, sTime string, location string) error {
	sub := Subscription{ChatID, sTime, location}
	if exi, err := existsSubscription(ctx, sub); exi && err == nil {
		return ErrExiAlready
	} else if err != nil {
		return err
	}

	if _, err := SubscriptionCollection.InsertOne(ctx, sub); err != nil {
		l.SimpleLogger.Error(err.Error())
		return err
	}

	l.SimpleLogger.Info(fmt.Sprintf("added user's [%d] query to the subscription list", ChatID))
	return nil
}

func existRecords(ctx context.Context, ChatID int64) (bool, error) {
	var res Subscription
	filter := bson.M{"chat_id": ChatID}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := SubscriptionCollection.FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		l.SimpleLogger.Error(err.Error())
		return false, err
	}

	return true, nil
}

func Unsubscribe(ctx context.Context, ChatID int64, sTime string, location string) error {
	sub := Subscription{ChatID, sTime, location}
	if exi, err := existsSubscription(ctx, sub); !exi && err == nil {
		return ErrNoSubRecord
	} else if err != nil {
		return err
	}

	if _, err := SubscriptionCollection.DeleteOne(ctx, sub); err != nil {
		l.SimpleLogger.Error(err.Error())
		return err
	}

	l.SimpleLogger.Info(fmt.Sprintf("deleted user's [%d] query from the subscription list", ChatID))
	return nil
}

func GetSubscriptions(ctx context.Context, ChatID int64) (map[Subscription]string, error) {
	if exiUser, err := existRecords(ctx, ChatID); !exiUser && err == nil {
		return nil, ErrNoUserRecord
	} else if err != nil {
		return nil, err
	}

	var userSubs []Subscription
	filter := bson.D{{Key: "chat_id", Value: ChatID}}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := SubscriptionCollection.Find(ctx, filter)
	if err != nil {
		l.SimpleLogger.Error(err.Error())
		return nil, err
	}

	if queryHandleErr := cursor.All(ctx, &userSubs); queryHandleErr != nil {
		l.SimpleLogger.Error(queryHandleErr.Error())
		return nil, err
	}

	res := make(map[Subscription]string, len(userSubs))
	for i := range userSubs {
		el := userSubs[i]
		res[el] = fmt.Sprint(el.Time, ", ", el.Location)
	}

	return res, nil
}

func AddTimezone(ctx context.Context, ChatID int64, tz string) error {
	var res Timezone

	timezone := Timezone{ChatID, tz}
	filter := bson.M{"chat_id": ChatID}
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := TimezoneCollection.FindOneAndReplace(ctxTimeout, filter, timezone).Decode(&res)
	switch err {
	case mongo.ErrNoDocuments:
		if _, err := TimezoneCollection.InsertOne(ctx, timezone); err != nil {
			l.SimpleLogger.Error("error adding a user's timezone: " + err.Error())
			return err
		}
	case nil:
		break
	default:
		l.SimpleLogger.Error("error handling user's timezone addition: " + err.Error())
		return err
	}

	l.SimpleLogger.Info(fmt.Sprintf("added user's [%d] timezone to the db", ChatID))
	return nil
}

func UserTimezoneKnown(ctx context.Context, ChatID int64) (bool, error) {
	var res Timezone
	filter := bson.M{"chat_id": ChatID}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := TimezoneCollection.FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		l.SimpleLogger.Error("error fetching info from Timezone collection: " + err.Error())
		return false, err
	}

	return true, nil
}

func GetForecastListUTC(ctx context.Context) ([]Subscription, error) {
	var userTimezones []Timezone

	ctxTz, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// get the whole collection
	cursor, err := TimezoneCollection.Find(ctxTz, bson.M{})
	if err != nil {
		l.SimpleLogger.Error("couldn't get whole Timezones collection: " + err.Error())
		return nil, err
	}

	if queryHandleErr := cursor.All(ctx, &userTimezones); queryHandleErr != nil {
		// if a collection is empty worker is not needed
		l.SimpleLogger.Error("error decoding a mongo document: " + queryHandleErr.Error())
		return nil, err
	}

	var res []Subscription

	for i := range userTimezones {
		chid := userTimezones[i].ChatID
		location, err := time.LoadLocation(userTimezones[i].Timezone)
		if err != nil {
			l.SimpleLogger.Error("failed loading location for a timezone: " + err.Error())
			return nil, err
		}

		subMap, err := GetSubscriptions(ctx, int64(chid))
		switch err {
		case ErrNoUserRecord:
			l.SimpleLogger.Error(fmt.Sprintf("there are no subscriptions for ChatID [%d]: %s", chid, err.Error()))
			l.SimpleLogger.Info("skipping a chat")
			continue
		case nil:
		default:
			l.SimpleLogger.Error("error communicating with the db: " + err.Error())
			return nil, err
		}

		// not interested in a .csv format, so only key is inspected
		for key := range subMap {
			userUTC, err := timezones.HHMMStringToUTC(key.Time, location)
			if err != nil {
				l.SimpleLogger.Error("failed to obtain a full list because of time conversion, aborting: " + err.Error())
				return nil, err
			}
			res = append(res, Subscription{key.ChatID, userUTC.Format("15:04"), key.Location})
		}
	}

	return res, nil
}
