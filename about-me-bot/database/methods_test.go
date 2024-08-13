package database

import (
	"context"
	"fmt"
	"testing"

	l "about-me-bot/internal/logger"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	testDatabase *TestDatabase
}

type SetupAllSuite interface {
	SetupSuite()
}

type TearDownAllSuite interface {
	TearDownSuite()
}

func (suite *TestSuite) SetupSuite() {
	suite.testDatabase = SetupTestDatabase(context.Background())
}

func (suite *TestSuite) TearDownSuite() {
	err := dropAll(dbName)
	if err != nil {
		l.SimpleLogger.Error(err.Error())
	}
	suite.testDatabase.container.Terminate(context.Background())
}

// TestAll All methods that begin with "Test" are run as tests within a suite.
func (suite *TestSuite) TestAll() {
	ctx := context.Background()
	suite.Run("Subscribe() basic test", func() {
		sub := Subscription{
			ChatID:   101,
			Time:     "17:57",
			Location: "Riga",
		}
		errSub := Subscribe(ctx, sub.ChatID, sub.Time, sub.Location)
		suite.Nil(errSub)

		successBool, errExi := existsSubscription(ctx, sub)
		suite.Nil(errExi)
		suite.True(successBool)

		recordsBool, errRecords := existRecords(ctx, sub.ChatID)
		suite.True(recordsBool)
		suite.Nil(errRecords)
	})

	suite.Run("Subscribe() with same entry test", func() {
		sub := Subscription{
			ChatID:   101,
			Time:     "17:57",
			Location: "Riga",
		}
		errSub := Subscribe(ctx, sub.ChatID, sub.Time, sub.Location)

		suite.Equal(ErrExiAlready, errSub)
	})

	suite.Run("Does not exist", func() {
		sub := Subscription{
			ChatID:   102,
			Time:     "17:57",
			Location: "Riga",
		}
		successBool, errExi := existsSubscription(ctx, sub)
		suite.False(successBool)
		suite.Nil(errExi)

		recordsBool, errRecords := existRecords(ctx, sub.ChatID)
		suite.False(recordsBool)
		suite.Nil(errRecords)
	})

	suite.Run("Unsubscribe() basic test", func() {
		sub := Subscription{
			ChatID:   101,
			Time:     "17:57",
			Location: "Riga",
		}
		err := Unsubscribe(ctx, sub.ChatID, sub.Time, sub.Location)
		suite.Nil(err)

		errUnsub := Unsubscribe(ctx, sub.ChatID, sub.Time, sub.Location)
		suite.Equal(ErrNoSubRecord, errUnsub)
	})

	suite.Run("GetSubscriptions() test", func() {
		sub := Subscription{
			ChatID:   101,
			Time:     "17:57",
			Location: "Riga",
		}
		nilMap, err := GetSubscriptions(ctx, sub.ChatID)
		suite.Nil(nilMap)
		suite.Equal(ErrNoUserRecord, err)

		// resubscribe
		Subscribe(ctx, sub.ChatID, sub.Time, sub.Location)
		resMap, errRes := GetSubscriptions(ctx, sub.ChatID)
		suite.Nil(errRes)
		temp := make(map[Subscription]string, 1)
		temp[sub] = fmt.Sprint(sub.Time, ", ", sub.Location)
		suite.Equal(temp, resMap)
	})

	suite.Run("GetForecastListUTC() test", func() {
		sub := Subscription{
			ChatID:   101,
			Time:     "17:57",
			Location: "Riga",
		}
		AddTimezone(ctx, sub.ChatID, "UTC")

		uTzKnown, err := UserTimezoneKnown(ctx, 101)
		suite.True(uTzKnown)
		suite.Nil(err)

		resSlice, err := GetForecastListUTC(ctx)
		suite.Nil(err)
		suite.Equal([]Subscription{sub}, resSlice)
	})

	suite.Run("User timezone not known test", func() {
		uTzKnown, err := UserTimezoneKnown(ctx, 102)
		suite.False(uTzKnown)
		suite.Nil(err)
	})
}

// TestMyDatabase In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMyDatabase(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
