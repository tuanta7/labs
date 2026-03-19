package main

import (
	"kafka-lab/internal/config"
	"kafka-lab/internal/repository"
	"log"

	"github.com/gocql/gocql"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	cluster := gocql.NewCluster(cfg.ScyllaHosts...)
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Error creating scylla session: %s", err)
	}
	defer session.Close()

	repo := repository.NewRepository(session)
	_ = repo
}
