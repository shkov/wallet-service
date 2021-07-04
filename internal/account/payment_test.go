package account

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidatePaymentRequest(t *testing.T) {
	testCases := []struct {
		name           string
		paymentRequest *PaymentRequest
		wantErr        error
	}{
		{
			name:           "normal response",
			paymentRequest: makePaymentRequest(t, nil),
			wantErr:        nil,
		},
		{
			name: "invalid amount",
			paymentRequest: makePaymentRequest(t, func(pr *PaymentRequest) {
				pr.Amount = "321,13"
			}),
			wantErr: errors.New("can't convert 321,13 to decimal"),
		},
		{
			name: "amount must be positive",
			paymentRequest: makePaymentRequest(t, func(pr *PaymentRequest) {
				pr.Amount = "-1"
			}),
			wantErr: errors.New("payment amount is not positive"),
		},
		{
			name: "from is not specified",
			paymentRequest: makePaymentRequest(t, func(pr *PaymentRequest) {
				pr.From = 0
			}),
			wantErr: errors.New("account from must be positive"),
		},
		{
			name: "to is not specified",
			paymentRequest: makePaymentRequest(t, func(pr *PaymentRequest) {
				pr.To = 0
			}),
			wantErr: errors.New("account to must be positive"),
		},
		{
			name: "from and to must be different",
			paymentRequest: makePaymentRequest(t, func(pr *PaymentRequest) {
				pr.To = 1
				pr.From = 1
			}),
			wantErr: errors.New("from and to must be different"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := ValidatePaymentRequest(tc.paymentRequest)
			if tc.wantErr == nil && gotErr != nil {
				t.Fatalf("unexpected error: %v", gotErr)
			}
			if tc.wantErr != nil && (gotErr == nil || tc.wantErr != gotErr) {
				assert.Equal(t, tc.wantErr, gotErr)
			}
		})
	}
}

func TestPaymentRequest_ToPayment(t *testing.T) {
	testCases := []struct {
		name           string
		paymentRequest *PaymentRequest
		payment        *Payment
		at             time.Time
	}{
		{
			name:           "normal response",
			paymentRequest: makePaymentRequest(t, nil),
			payment:        makePayment(t, nil),
			at:             parseTime(t, "2001-01-02T11:22:33+03:00"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.paymentRequest.ToPayment(tc.at)
			assert.Equal(t, tc.payment, got)
		})
	}
}

func makePaymentRequest(t *testing.T, fn func(*PaymentRequest)) *PaymentRequest {
	pr := &PaymentRequest{
		From:   1,
		To:     2,
		Amount: "500",
	}
	if fn != nil {
		fn(pr)
	}
	return pr
}

func makePayment(t *testing.T, fn func(*Payment)) *Payment {
	p := &Payment{
		ID:        0,
		From:      1,
		To:        2,
		Amount:    "500",
		CreatedAt: parseTime(t, "2001-01-02T11:22:33+03:00"),
	}
	if fn != nil {
		fn(p)
	}
	return p
}

func parseTime(t *testing.T, timestr string) time.Time {
	tm, err := time.Parse(time.RFC3339, timestr)
	if err != nil {
		t.Fatalf("parseTime: failed to parse %q: %v", timestr, err)
	}
	return tm
}
