package database

import (
	"context"
	"fmt"
	"time"

	l "about-me-bot/internal/logger"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbName string = "test-db"
)

func createMongoContainer(ctx context.Context) (testcontainers.Container, *mongo.Client, string, error) {
	env := map[string]string{
		"MONGO_INITDB_DATABASE": dbName,
	}
	port := "27017/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{port},
			Env:          env,
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	/*
		p, err := container.MappedPort(ctx, "27017")
		if err != nil {
			return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
		}

		l.SimpleLogger.Info("mongo container ready and running at port: " + p.Port())
	*/

	uri := fmt.Sprintf("mongodb://localhost:%s", port)
	errInit := Init(ctx, uri, dbName)
	if errInit != nil {
		return container, DBClient, uri, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, DBClient, uri, nil
}

type TestDatabase struct {
	DbInstance *mongo.Client
	DbURI      string
	container  testcontainers.Container
}

func SetupTestDatabase(ctx context.Context) *TestDatabase {
	ctx, _ = context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := createMongoContainer(ctx)
	if err != nil {
		l.SimpleLogger.Log(ctx, l.LevelFatal, "failed to setup test "+err.Error())
	}

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbURI:      dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	_ = tdb.container.Terminate(context.Background())
}

func dropAll(dbName string) error {
	// Get a handle for the database
	ctx := context.Background()
	db := DBClient.Database(dbName)

	// List all collections in the database
	collections, err := db.ListCollectionNames(ctx, struct{}{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}

	// Drop each collection
	for _, collName := range collections {
		coll := db.Collection(collName)
		if err := coll.Drop(ctx); err != nil {
			return fmt.Errorf("failed to drop collection %s: %v", collName, err)
		}
		l.SimpleLogger.Info("Dropped collection: " + collName)
	}

	l.SimpleLogger.Info("All collections have been dropped.")
	return nil
}
