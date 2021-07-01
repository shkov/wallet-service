package account

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID        int64
	Balance   string
	CreatedAt time.Time
}

func Create(id int64, createdAt time.Time) *Account {
	return &Account{
		ID:        id,
		Balance:   "1000",
		CreatedAt: createdAt,
	}
}

func (a *Account) ApplyPayment(p *Payment) error {
	balance, err := decimal.NewFromString(a.Balance)
	if err != nil {
		return err
	}
	amount, err := decimal.NewFromString(p.Amount)
	if err != nil {
		return err
	}

	switch a.ID {
	case p.From:
		if balance.LessThan(amount) {
			return ErrNotEnoughFunds
		}
		a.Balance = balance.Sub(amount).StringFixed(2)

	case p.To:
		a.Balance = balance.Add(amount).StringFixed(2)

	default:
		return ErrMismatchPayment
	}

	return nil
}

func ValidateAccountID(id int64) error {
	if id <= 0 {
		return ErrMustBePositive
	}
	return nil
}
