package walletservice

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/shkov/wallet-service/internal/account"
)

type mockService struct {
	onApplyPayment func(ctx context.Context, p *account.PaymentRequest) (*account.Payment, error)
	onGetPayments  func(ctx context.Context, accountID int64) ([]*account.Payment, error)
	onGetAccount   func(ctx context.Context, id int64) (*account.Account, error)
}

func (m *mockService) ApplyPayment(ctx context.Context, p *account.PaymentRequest) (*account.Payment, error) {
	return m.onApplyPayment(ctx, p)
}

func (m *mockService) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	return m.onGetPayments(ctx, accountID)
}

func (m *mockService) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	return m.onGetAccount(ctx, id)
}

// returns mocked server and http client and mocked service for transport testing.
func initTransportTest(t *testing.T) (*httptest.Server, Service, *mockService) {
	svc := &mockService{}
	handler := makeHandler(svc)
	server := httptest.NewServer(handler)
	client, err := NewClient(ClientConfig{
		ServiceURL: server.URL,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return server, client, svc
}

func TestTransportApplyPayment(t *testing.T) {
	server, client, svc := initTransportTest(t)
	defer server.Close()

	testCases := []struct {
		name     string
		request  *account.PaymentRequest
		response *account.Payment
		err      error
	}{
		{
			name:     "ok",
			request:  makePaymentRequest(t, nil),
			response: makePayment(t, nil),
			err:      nil,
		},
		{
			name:     "some err",
			request:  makePaymentRequest(t, nil),
			response: nil,
			err:      &serviceError{code: 500, Message: "kek some err occurs"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc.onApplyPayment = func(ctx context.Context, p *account.PaymentRequest) (*account.Payment, error) {
				return tc.response, tc.err
			}

			gotResp, gotErr := client.ApplyPayment(context.Background(), tc.request)
			assert.Equal(t, tc.err, gotErr)
			assert.Equal(t, tc.response, gotResp)
		})
	}
}

func TestTransportGetAccount(t *testing.T) {
	server, client, svc := initTransportTest(t)
	defer server.Close()

	testCases := []struct {
		name     string
		id       int64
		response *account.Account
		err      error
	}{
		{
			name:     "ok",
			id:       1,
			response: makeAccount(t, nil),
			err:      nil,
		},
		{
			name:     "some err",
			id:       1,
			response: nil,
			err: &serviceError{
				code:    500,
				Message: "kek some err occurs",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc.onGetAccount = func(ctx context.Context, id int64) (*account.Account, error) {
				return tc.response, tc.err
			}

			gotResp, gotErr := client.GetAccount(context.Background(), tc.id)
			assert.Equal(t, tc.err, gotErr)
			assert.Equal(t, tc.response, gotResp)
		})
	}
}

func TestTransportGetPayments(t *testing.T) {
	server, client, svc := initTransportTest(t)
	defer server.Close()

	testCases := []struct {
		name      string
		accountID int64
		response  []*account.Payment
		err       error
	}{
		{
			name:      "ok",
			accountID: 1,
			response: []*account.Payment{
				makePayment(t, nil),
			},
			err: nil,
		},
		{
			name:      "some err",
			accountID: 1,
			response:  nil,
			err: &serviceError{
				code:    500,
				Message: "kek some err occurs",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc.onGetPayments = func(ctx context.Context, accountID int64) ([]*account.Payment, error) {
				return tc.response, tc.err
			}

			gotResp, gotErr := client.GetPayments(context.Background(), tc.accountID)
			assert.Equal(t, tc.err, gotErr)
			assert.Equal(t, tc.response, gotResp)
		})
	}
}
