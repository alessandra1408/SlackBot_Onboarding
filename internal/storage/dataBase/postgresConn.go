package database

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// parametro: log *nlog.NLogger
func NewConnection(postgresURL string) (*sqlx.DB, error) {

	db, sErr := sqlx.Open("postgres", postgresURL)
	if sErr != nil {
		log.Printf("error on PostgreSQL connection: %q", sErr)
		return nil, sErr
	}
	pErr := db.Ping()
	if pErr != nil {
		log.Printf("error on pingDB : %q", sErr)
		return nil, pErr
	}
	return db, nil
}
