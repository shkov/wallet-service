package account

import "time"

type Payment struct {
	ID        int64
	From      int64
	To        int64
	Amount    string
	CreatedAt time.Duration
}

type PaymentRequest struct {
	From   int64
	To     int64
	Amount string
}
