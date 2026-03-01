package main

import (
	"log"

	"github.com/gocql/gocql"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	cluster := gocql.NewCluster(cfg.ScyllaHosts...)
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Error creating scylla session: %s", err)
	}
	defer session.Close()

	repo := NewRepository(session)
	_ = repo
}
