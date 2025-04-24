package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"grpc.bank-api/internal/config"
)

type Store interface {
	// Querier
}

type SQLStore struct {
	// *Queries
	db *sql.DB
}

func NewDatabase(cfg *config.Config) (Store, error) {

	fmt.Printf("Connecting to %s, %s", cfg.DbSource, cfg.DbDriver)

	db, err := sql.Open(cfg.DbDriver, cfg.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("cannot ping db:", err)
	}
	return &SQLStore{
		db: db,
		// Queries: New(db),
	}, nil

}
