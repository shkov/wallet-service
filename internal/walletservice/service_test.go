package walletservice

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/shkov/wallet-service/internal/account"
	"github.com/shkov/wallet-service/internal/storage"
)

type storageMock struct {
	onGetAccount      func(ctx context.Context, id int64) (*account.Account, error)
	onGetAccounts     func(ctx context.Context, ids []int64) ([]*account.Account, error)
	onGetPayments     func(ctx context.Context, accountID int64) ([]*account.Payment, error)
	onInsertPayment   func(ctx context.Context, p *account.Payment) error
	onReplaceAccounts func(ctx context.Context, aa []*account.Account) error
	onClose           func() error
	onExecTx          func(ctx context.Context, fn func(context.Context, storage.Storage) error) error
}

func (m *storageMock) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	return m.onGetAccount(ctx, id)
}

func (m *storageMock) GetAccounts(ctx context.Context, ids []int64) ([]*account.Account, error) {
	return m.onGetAccounts(ctx, ids)
}

func (m *storageMock) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	return m.onGetPayments(ctx, accountID)
}

func (m *storageMock) InsertPayment(ctx context.Context, p *account.Payment) error {
	return m.onInsertPayment(ctx, p)
}

func (m *storageMock) ReplaceAccounts(ctx context.Context, aa []*account.Account) error {
	return m.onReplaceAccounts(ctx, aa)
}

func (m *storageMock) Close() error {
	return m.onClose()
}

func (m *storageMock) ExecTx(ctx context.Context, fn func(context.Context, storage.Storage) error) error {
	return m.onExecTx(ctx, fn)
}

func TestService_ApplyPayment(t *testing.T) {
	testCases := []struct {
		name           string
		paymentRequest *account.PaymentRequest
		wantPayment    *account.Payment
		wantErr        error
	}{
		{
			name:           "normal response",
			paymentRequest: makePaymentRequest(t, nil),
			wantPayment:    makePayment(t, nil),
			wantErr:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := &storageMock{
				onGetAccounts: func(ctx context.Context, got []int64) ([]*account.Account, error) {
					want := []int64{
						1,
						2,
					}
					if ok := assert.Equal(t, want, got); !ok {
						t.Fatal()
					}
					return []*account.Account{
						makeAccount(t, nil),
						makeAccount(t, func(a *account.Account) {
							a.ID = 2
						}),
					}, nil
				},
				onInsertPayment: func(ctx context.Context, got *account.Payment) error {
					want := makePayment(t, nil)
					if ok := assert.Equal(t, want, got); !ok {
						t.Fatal()
					}
					return nil
				},

				onReplaceAccounts: func(ctx context.Context, got []*account.Account) error {
					want := []*account.Account{
						makeAccount(t, func(a *account.Account) {
							a.ID = 1
							a.Balance = "500.00"
						}),
						makeAccount(t, func(a *account.Account) {
							a.ID = 2
							a.Balance = "1500.00"
						}),
					}
					if ok := assert.Equal(t, want, got); !ok {
						t.Fatal()
					}
					return nil
				},
			}

			mock.onExecTx = func(ctx context.Context, fn func(context.Context, storage.Storage) error) error {
				return fn(ctx, mock)
			}

			svc := &serviceImpl{
				logger:  log.NewNopLogger(),
				storage: mock,
				now: func() time.Time {
					return parseTime(t, "2001-01-02T11:22:33+03:00")
				},
			}

			gotResp, gotErr := svc.ApplyPayment(context.Background(), tc.paymentRequest)
			assert.Equal(t, tc.wantPayment, gotResp)
			if tc.wantErr == nil && gotErr != nil {
				t.Fatalf("unexpected error: %v", gotErr)
			}
			if tc.wantErr != nil && (gotErr == nil || tc.wantErr != gotErr) {
				assert.Equal(t, tc.wantErr, gotErr)
			}
		})
	}
}

func makePayment(t *testing.T, fn func(*account.Payment)) *account.Payment {
	p := &account.Payment{
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

func makeAccount(t *testing.T, fn func(*account.Account)) *account.Account {
	a := &account.Account{
		ID:        1,
		Balance:   "1000",
		CreatedAt: parseTime(t, "2001-01-02T11:22:33+03:00"),
	}
	if fn != nil {
		fn(a)
	}
	return a
}

func makePaymentRequest(t *testing.T, fn func(*account.PaymentRequest)) *account.PaymentRequest {
	pr := &account.PaymentRequest{
		From:   1,
		To:     2,
		Amount: "500",
	}
	if fn != nil {
		fn(pr)
	}
	return pr
}

func parseTime(t *testing.T, timestr string) time.Time {
	tm, err := time.Parse(time.RFC3339, timestr)
	if err != nil {
		t.Fatalf("parseTime: failed to parse %q: %v", timestr, err)
	}
	return tm
}
