package main

import (
	"context"
	"kafka/publisher"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3/table"
)

var locationTable = table.New(table.Metadata{
	Name:    "location",
	Columns: []string{"id", "latitude", "longitude", "timestamp", "user_id"},
	PartKey: []string{"user_id"},
	SortKey: []string{"timestamp"},
})

type Repository struct {
	scylla *gocql.Session
}

func NewRepository(scylla *gocql.Session) *Repository {
	return &Repository{
		scylla: scylla,
	}
}

func (r *Repository) SaveLocation(ctx context.Context, location *publisher.Location) error {
	q := r.scylla.Query(locationTable.Insert()).Bind(location)
	if err := q.Exec(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLatestLocation(userID string) (*publisher.Location, error) {
	return nil, nil
}
