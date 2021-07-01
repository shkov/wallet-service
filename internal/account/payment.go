package account

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	ID        int64
	From      int64
	To        int64
	Amount    string
	CreatedAt time.Time
}

type PaymentRequest struct {
	From   int64
	To     int64
	Amount string
}

func (r *PaymentRequest) ToPayment(createdAt time.Time) *Payment {
	return &Payment{
		ID:        0,
		From:      r.From,
		To:        r.To,
		Amount:    r.Amount,
		CreatedAt: createdAt,
	}
}

func ValidatePaymentRequest(r *PaymentRequest) error {
	amount, err := decimal.NewFromString(r.Amount)
	if err != nil {
		return err
	}
	if !amount.IsPositive() {
		return ErrNotPositiveAmount
	}
	if r.From <= 0 {
		return ErrAccountFromMustBePositive
	}
	if r.To <= 0 {
		return ErrAccountToMustBePositive
	}
	return nil
}
