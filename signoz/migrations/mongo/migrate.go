package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tuanta7/k6noz/services/pkg/slient"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		panic("MONGO_URI is required")
	}

	database := os.Getenv("MONGO_DATABASE")
	if database == "" {
		database = "test"
	}

	collection := "drivers"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	slient.PanicOnErr(err, "failed to connect to mongodb")
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	slient.PanicOnErr(err, "failed to ping mongodb")

	data, err := readMockData("./migrations/data/drivers.json")
	slient.PanicOnErr(err, "failed to read mock data")

	coll := client.Database(database).Collection(collection)

	err = deleteAll(ctx, coll)
	slient.PanicOnErr(err, "failed to delete all data")

	err = insertMockData(ctx, coll, data)
	slient.PanicOnErr(err, "failed to insert mock data")
}

func readMockData(filePath string) ([]map[string]any, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data []map[string]any
	if err = json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return data, nil
}

func deleteAll(ctx context.Context, collections *mongo.Collection) error {
	_, err := collections.DeleteMany(ctx, map[string]any{})
	return err
}

func insertMockData(ctx context.Context, collections *mongo.Collection, data []map[string]any) error {
	_, err := collections.InsertMany(ctx, data)
	return err
}
