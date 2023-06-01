package storage

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type postgresStorage struct {
	DB *sqlx.DB
	//incluir log *nlog.NLogger
}

//adicionar parametro do nlogger futuramente
func NewPostgresStorage(db *sqlx.DB) *postgresStorage {
	return &postgresStorage{
		DB: db,
	}
}

//Queryx queries
func (p *postgresStorage) Queryx(query string) (*sqlx.Rows, error) {
	return p.DB.Queryx(query)
}

//Exec queries
func (p *postgresStorage) Exec(query string) (sql.Result, error) {
	return p.DB.Exec(query)
}

func (p *postgresStorage) MustExec(query string, args ...any) sql.Result {
	return p.DB.MustExec(query, args...)
}

//Close database connections
func (p *postgresStorage) Close() error {
	return p.DB.Close()
}
