package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SagaStore struct {
	collection *mongo.Collection
}

func NewSagaStore(client *mongo.Client) *SagaStore {
	return &SagaStore{
		collection: client.Database("saga_orchestrator").Collection("saga"),
	}
}

// CreateSaga inserts a new saga into the database.
func (s *SagaStore) CreateSaga(ctx context.Context, saga *Saga) error {
	saga.CreatedAt = time.Now().UTC()
	saga.UpdatedAt = saga.CreatedAt

	_, err := s.collection.InsertOne(ctx, saga)
	return err
}

func (s *SagaStore) GetSaga(ctx context.Context, id string) (*Saga, error) {
	var saga Saga
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&saga)
	if err != nil {
		return nil, err
	}

	return &saga, nil
}

func (s *SagaStore) UpdateStatus(ctx context.Context, id string, status SagaState, errMsg string) error {
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
