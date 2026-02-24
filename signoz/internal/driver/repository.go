package driver

import (
	"context"

	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository struct {
	driverCollection mongo.Collection
	ratingCollection mongo.Collection
}

func NewRepository(mongo *mongo.Client) *Repository {
	return &Repository{
		driverCollection: mongo.Collection("drivers"),
		ratingCollection: mongo.Collection("ratings"),
	}
}

func (r *Repository) GetDriverByID(ctx context.Context, driverID string) (*domain.Driver, error) {
	var driver domain.Driver
	err := r.driverCollection.FindOne(ctx, bson.M{"_id": driverID}, &driver)
	if err != nil {
		return nil, err
	}

	return &driver, nil
}

func (r *Repository) UpdateDriverRating(ctx context.Context, driverID string, rating, score float64) error {
	return r.driverCollection.UpdateOne(ctx,
		bson.M{"_id": driverID},
		bson.M{"$set": bson.M{
			"rating": rating,
			"score":  score,
		}},
	)
}

func (r *Repository) InsertRating(ctx context.Context, rating *domain.Rating) error {
	return r.ratingCollection.InsertOne(ctx, rating)
}
