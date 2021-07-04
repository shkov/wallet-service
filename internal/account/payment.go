package account

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	tableName struct{}  `pq:"payments"`
	ID        int64     `pq:"id"`
	From      int64     `pg:"from_account_id"`
	To        int64     `pg:"to_account_id"`
	Amount    string    `pq:"amount"`
	CreatedAt time.Time `pq:"create_at"`
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
	if r.From == r.To {
		return ErrFromAndToMustBeDifferent
	}
	return nil
}
