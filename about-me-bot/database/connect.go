package database

import (
	"context"
	"time"

	l "about-me-bot/internal/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DBClient               *mongo.Client
	SubscriptionCollection *mongo.Collection
	TimezoneCollection     *mongo.Collection
)

func Init(ctx context.Context, URI string, db string) error {
	var err error
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(URI)
	DBClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = DBClient.Ping(ctx, nil)
	if err != nil {
		return err
	}

	l.SimpleLogger.Info("successfully established a connection with db")
	SubscriptionCollection = DBClient.Database(db).Collection("Subscriptions")
	TimezoneCollection = DBClient.Database(db).Collection("Timezones")
	return nil
}
