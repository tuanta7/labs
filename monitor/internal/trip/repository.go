package trip

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	Database          = "monitor"
	DriversCollection = "drivers"
	TripsCollection   = "trips"
)

type Repository struct {
	driverCollection *mongo.Collection
	tripCollection   *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		driverCollection: client.Database(Database).Collection(DriversCollection),
		tripCollection:   client.Database(Database).Collection(TripsCollection),
	}
}

func (r *Repository) CreateTrip(ctx context.Context, trip *Trip) error {
	time.Sleep(10 * time.Millisecond)
	_, err := r.tripCollection.InsertOne(ctx, trip)
	return err
}

func (r *Repository) AcceptTrip(ctx context.Context, trip *Trip) error {
	filter := bson.M{
		"_id":    trip.ID,
		"status": StatusPending,
	}

	update := bson.M{
		"$set": bson.M{
			"status":   StatusActive,
			"driverId": trip.DriverID,
		},
	}

	res, err := r.tripCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		// already accepted/cancelled/not found
	}

	return nil
}
