package repository

import (
	"context"
	"kafka-lab/internal/domain"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3/table"
)

var locationTable = table.New(table.Metadata{
	Name:    "location",
	Columns: []string{"id", "latitude", "longitude", "timestamp", "user_id"},
	PartKey: []string{"user_id"},
	SortKey: []string{"timestamp"},
})

type LocationRepository struct {
	scylla *gocql.Session
}

func NewRepository(scylla *gocql.Session) *LocationRepository {
	return &LocationRepository{
		scylla: scylla,
	}
}

func (r *LocationRepository) SaveLocation(ctx context.Context, location *domain.Location) error {
	q := r.scylla.Query(locationTable.Insert()).Bind(location)
	if err := q.Exec(); err != nil {
		return err
	}

	return nil
}

func (r *LocationRepository) GetLatestLocation(userID string) (*domain.Location, error) {
	return nil, nil
}
