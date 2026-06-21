package saga

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(client *mongo.Client) *Store {
	return &Store{
		collection: client.Database("saga_orchestrator").Collection("saga"),
	}
}

func (s *Store) CreateSaga(ctx context.Context, data *Saga) error {
	data.CreatedAt = time.Now().UTC()
	data.UpdatedAt = data.CreatedAt

	_, err := s.collection.InsertOne(ctx, data)
	return err
}

func (s *Store) GetSaga(ctx context.Context, id string) (*Saga, error) {
	var saga Saga
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&saga)
	if err != nil {
		return nil, err
	}

	return &saga, nil
}

func (s *Store) UpdateStatus(ctx context.Context, id string, status State, errMsg string) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"error":      errMsg,
			"updated_at": time.Now(),
		},
	}

	_, err := s.collection.UpdateByID(ctx, id, update)
	return err
}
