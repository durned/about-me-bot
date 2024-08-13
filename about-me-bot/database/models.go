package database

import (
	"errors"
	"fmt"
)

type Subscription struct {
	ChatID   int64  `bson:"chat_id"`
	Time     string `bson:"time"`
	Location string `bson:"location"`
}

func (s Subscription) String() string {
	return fmt.Sprintf("%d,%s,%s", s.ChatID, s.Time, s.Location) // basically a csv fmt
}

type Timezone struct {
	ChatID   int64  `bson:"chat_id"`
	Timezone string `bson:"timezone"`
}

var (
	ErrExiAlready   error = errors.New("this subscription already exists")
	ErrNoUserRecord error = errors.New("there is no record for such ChatID")
	ErrNoSubRecord  error = errors.New("there is no record for such subscription")
)
