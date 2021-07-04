package account

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_ApplyPayment(t *testing.T) {
	testCases := []struct {
		name        string
		account     *Account
		payment     *Payment
		wantAccount *Account
		wantErr     error
	}{
		{
			name: "from: normal response",
			account: &Account{
				ID:      1,
				Balance: "1000",
			},
			payment: &Payment{
				From:   1,
				To:     2,
				Amount: "501",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "499.00",
			},
			wantErr: nil,
		},
		{
			name: "from: not enough funds",
			account: &Account{
				ID:      1,
				Balance: "1000",
			},
			payment: &Payment{
				From:   1,
				To:     2,
				Amount: "1001",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "1000",
			},
			wantErr: errors.New("not enough funds in account"),
		},
		{
			name: "from: use all balance",
			account: &Account{
				ID:      1,
				Balance: "1000",
			},
			payment: &Payment{
				From:   1,
				To:     2,
				Amount: "1000",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "0.00",
			},
			wantErr: nil,
		},
		{
			name: "from: test rounding",
			account: &Account{
				ID:      1,
				Balance: "1000.6",
			},
			payment: &Payment{
				From:   1,
				To:     2,
				Amount: "100.8",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "899.80",
			},
			wantErr: nil,
		},

		{
			name: "to: normal response",
			account: &Account{
				ID:      1,
				Balance: "1000",
			},
			payment: &Payment{
				From:   2,
				To:     1,
				Amount: "501",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "1501.00",
			},
			wantErr: nil,
		},
		{
			name: "to: initial zero balance",
			account: &Account{
				ID:      1,
				Balance: "0",
			},
			payment: &Payment{
				From:   2,
				To:     1,
				Amount: "501",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "501.00",
			},
			wantErr: nil,
		},
		{
			name: "to: with cents",
			account: &Account{
				ID:      1,
				Balance: "10.1",
			},
			payment: &Payment{
				From:   2,
				To:     1,
				Amount: "15.9",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "26.00",
			},
			wantErr: nil,
		},
		{
			name: "to: test rounding",
			account: &Account{
				ID:      1,
				Balance: "10.13",
			},
			payment: &Payment{
				From:   2,
				To:     1,
				Amount: "15.98",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "26.11",
			},
			wantErr: nil,
		},
		{
			name: "invalid format of balance",
			account: &Account{
				ID:      1,
				Balance: "10,1",
			},
			payment: &Payment{
				From:   2,
				To:     1,
				Amount: "1,1",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "10,1",
			},
			wantErr: errors.New("can't convert 10,1 to decimal"),
		},
		{
			name: "invalid format of amount",
			account: &Account{
				ID:      1,
				Balance: "10.1",
			},
			payment: &Payment{
				From:   2,
				To:     1,
				Amount: "1,1",
			},
			wantAccount: &Account{
				ID:      1,
				Balance: "10.1",
			},
			wantErr: errors.New("can't convert 1,1 to decimal"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := tc.account.ApplyPayment(tc.payment)
			assert.Equal(t, tc.wantAccount, tc.account)
			if tc.wantErr == nil && gotErr != nil {
				t.Fatalf("unexpected error: %v", gotErr)
			}
			if tc.wantErr != nil && (gotErr == nil || tc.wantErr != gotErr) {
				assert.Equal(t, tc.wantErr, gotErr)
			}
		})
	}
}
