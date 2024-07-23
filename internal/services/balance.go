package services

import (
	"github.com/jmoiron/sqlx"
)

type Service struct {
	*sqlx.DB
}

func New(db *sqlx.DB) *Service {
	return &Service{
		db,
	}
}

func (s *Service) CreateWallet() {
	// Создание кошелька в бд
}
