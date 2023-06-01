package core

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	ServiceSPS SPostgresStorage
	ServiceSS  SSlack
}

func NewService(serviceSPS SPostgresStorage, serviceSS SSlack) *Service {
	return &Service{
		ServiceSPS: serviceSPS,
		ServiceSS:  serviceSS,
	}
}

//Close database connection
func (s *Service) Close() error {
	if err := s.ServiceSPS.Close(); err != nil {
		return err
	}
	return nil
}

//Queryx queries to run
func (s *Service) Queryx(query string) (*sqlx.Rows, error) {
	return s.ServiceSPS.Queryx(query)
}

//Exec queries to run
func (s *Service) Exec(query string) (sql.Result, error) {
	return s.ServiceSPS.Exec(query)
}

//MustExec Keep Event Sourcing queries
func (s *Service) MustExec(query string, args ...any) sql.Result {
	return s.ServiceSPS.MustExec(query, args...)
}

func (s *Service) Running() error {
	//adicionar aqui as regras de negocio

	fmt.Println("Running rodando!")
	//adicionar o que a parte do slack precisa fazer aqui
	return nil
}

//criar mais metodos depois
