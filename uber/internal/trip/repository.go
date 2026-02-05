package trip

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/tuanta7/k6noz/services/internal/domain"
	"github.com/tuanta7/k6noz/services/pkg/clickhouse"
)

type Repository struct {
	db clickhouse.DB
}

func NewRepository(db clickhouse.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetTripByID(ctx context.Context, tripID string) (*domain.Trip, error) {
	return nil, nil
}

func (r *Repository) BatchInsertLocations(ctx context.Context, locations []*domain.Location) error {
	query, _, err := sq.Insert("trips").Columns().ToSql()
	if err != nil {
		return err
	}

	batch, err := r.db.PrepareBatch(ctx, query)
	if err != nil {
		return err
	}
	defer batch.Close()

	for _, location := range locations {
		err = batch.Append(
			location.TripID,
			location.Latitude,
			location.Longitude,
			location.Timestamp,
		)
	}

	return batch.Send()
}
