package models

import (
	"errors"
	"time"
)

const (
	INDENTIFICATION_LEVEL_ANONYMOUS  = 1
	INDENTIFICATION_LEVEL_SIMPLIFIED = 2
	INDENTIFICATION_LEVEL_FULL       = 3
)

type Wallet struct {
	ID                  int       `db:"id"`
	Balance             int       `db:"balance"`
	Hold                int       `db:"hold"`
	IdentificationLevel int       `db:"identification_level"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
}

type LimitLaw struct {
	ID                  int    `db:"id"`
	IdentificationLevel string `db:"identifiaction_level"`
	BalanceMin          int    `db:"balance_min"`
	BalanceMax          int    `db:"balance_max"`
}

func (w *Wallet) HoldBalance(hold int) error {
	if hold > w.Balance {
		return errors.New("вы не можете захолдировать больше, чем у вас есть")
	}

	if w.Hold > w.Balance {
		return errors.New("денег нет уже аааа")
	}

	if hold <= 0 {
		return errors.New("вы не можете захолдировать ноль или меньше")
	}

	w.Hold += hold
	return nil
}

func (w *Wallet) LimitLawCheck(l *LimitLaw) bool {
	if w.IdentificationLevel == l.ID && w.Balance >= l.BalanceMax || w.Balance <= l.BalanceMin {
		return false
	}
	return true
}

func (l *LimitLaw) LimitLevelCheck(value int) bool {
	switch l.ID {
	case INDENTIFICATION_LEVEL_FULL:
		return true
	case INDENTIFICATION_LEVEL_SIMPLIFIED:
		return true
	case INDENTIFICATION_LEVEL_ANONYMOUS:
		return true
	}
	return false
}
