package services

import (
	"sync"

	"github.com/desulaidovich/balance/internal/models"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
	mu sync.Mutex
}

func New(db *sqlx.DB) *Service {
	service := new(Service)
	service.db = db
	service.mu = sync.Mutex{}

	return service
}

func (s *Service) GetLimitByID(id int) (*models.LimitLaw, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var limit models.LimitLaw
	err := s.db.Get(&limit, "SELECT * FROM limit_law WHERE id=$1;", id)

	if err != nil {
		return nil, err
	}

	return &limit, nil
}

func (s *Service) GetWalletByID(id int) (*models.Wallet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var wallet models.Wallet
	err := s.db.Get(&wallet, "SELECT * FROM balance WHERE id=$1;", id)

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (s *Service) UpdateWallet(w *models.Wallet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.NamedQuery(`UPDATE balance SET balance=:balance, hold=:hold, identification_level=:identification_level, 
		updated_at=CURRENT_DATE WHERE id=:id;`, &w)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateWallet(w *models.Wallet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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
