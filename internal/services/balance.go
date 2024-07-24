package services

import (
	"github.com/desulaidovich/balance/internal/models"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Service {
	return &Service{
		db,
	}
}

func (s *Service) GetLimitByID(id string) (*models.LimitLaw, error) {
	limit := models.LimitLaw{}
	err := s.db.Get(&limit, "SELECT * FROM limit_law WHERE id=$1;", id)

	if err != nil {
		return nil, err
	}

	return &limit, nil
}

func (s *Service) CreateWallet(w *models.Wallet) error {
	rows, err := s.db.NamedQuery(`INSERT INTO balance (balance, hold, identification_level, created_at) 
		VALUES (:balance, :hold, :identification_level, CURRENT_DATE) RETURNING id;`, &w)

	if err != nil {
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&w.ID); err != nil {
			return err
		}
	}

	return nil
}
