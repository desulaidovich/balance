package models

import (
	"errors"
	"time"
)

const (
	TYPE_ID_DEBIT   = 1
	TYPE_ID_DEPOSIT = 2
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

func (w *Wallet) LimitLawCheck(l *LimitLaw) error {
	if w.Balance >= l.BalanceMax {
		return errors.New("превышение допустимого лимита")
	}
	if w.Balance <= l.BalanceMin {
		return errors.New("ниже допустимого лимита")
	}
	return nil
}

func (w *Wallet) HoldBalance(money int) error {
	if money <= 0 {
		return errors.New("отрицательное значение")
	}
	if w.Balance < (w.Hold + money) {
		return errors.New("недостаточно денег")
	}

	w.Hold += money
	return nil
}

func (w *Wallet) DisholdBalance(money int) error {
	if money <= 0 {
		return errors.New("отрицательное значение")
	}
	if w.Hold < money {
		return errors.New("недостаточно денег")
	}

	w.Hold -= money
	return nil
}

func (w *Wallet) EditWithType(l *LimitLaw, typeID int, money int) error {
	if money <= 0 {
		return errors.New("отрицательное значение")
	}

	switch typeID {
	case TYPE_ID_DEBIT:
		{
			if err := w.DebitBalance(money); err != nil {
				return err
			}
			return nil
		}
	case TYPE_ID_DEPOSIT:
		{
			if err := w.DepositBalance(l, money); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("неизвестный тип")
}

func (w *Wallet) DebitBalance(money int) error {
	if err := w.DisholdBalance(money); err != nil {
		return err
	}

	w.Balance -= money
	return nil
}

func (w *Wallet) DepositBalance(l *LimitLaw, money int) error {
	w.Balance += money
	if err := w.LimitLawCheck(l); err != nil {
		return err
	}
	return nil
}

func (w *Wallet) GetDates() (string, string) {
	createAt := w.CreatedAt.Format(time.DateOnly)
	updatedAt := w.UpdatedAt.Format(time.DateOnly)
	return createAt, updatedAt
}
