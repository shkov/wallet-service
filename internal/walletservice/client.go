package walletservice

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/shkov/wallet-service/internal/account"
)

// ClientConfig is an Client configuration.
type ClientConfig struct {
	ServiceURL string
	Timeout    time.Duration
}

func (cfg ClientConfig) validate() error {
	if cfg.ServiceURL == "" {
		return errors.New("must provide ServiceURL")
	}
	if cfg.Timeout <= 0 {
		return errors.New("invalid Timeout")
	}
	return nil
}

var _ Service = (*client)(nil)

// Client is a wallet-service client.
type client struct {
	getPaymentsEndpoint  endpoint.Endpoint
	getAccountEndpoint   endpoint.Endpoint
	applyPaymentEndpoint endpoint.Endpoint
}

// NewClient creates a new client.
func NewClient(cfg ClientConfig) (Service, error) {
	err := cfg.validate()
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(cfg.ServiceURL)
	if err != nil {
		return nil, err
	}

	options := []kithttp.ClientOption{
		kithttp.SetClient(&http.Client{
			Timeout: cfg.Timeout,
		}),
	}

	c := &client{
		getAccountEndpoint: kithttp.NewClient(
			http.MethodGet,
			baseURL,
			encodeGetAccountRequest,
			decodeGetAccountResponse,
			options...,
		).Endpoint(),
		getPaymentsEndpoint: kithttp.NewClient(
			http.MethodGet,
			baseURL,
			encodeGetPaymentsRequest,
			decodeGetPaymentsResponse,
			options...,
		).Endpoint(),
		applyPaymentEndpoint: kithttp.NewClient(
			http.MethodPost,
			baseURL,
			encodeApplyPaymentRequest,
			decodeApplyPaymentResponse,
			options...,
		).Endpoint(),
	}

	return c, nil
}

func (c *client) ApplyPayment(ctx context.Context, p *account.Payment) error {
	_, err := c.applyPaymentEndpoint(ctx, applyPaymentRequest{input: p})
	if err != nil {
		return err
	}
	return nil

}

func (c *client) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	response, err := c.getPaymentsEndpoint(ctx, getPaymentsRequest{accountID: accountID})
	if err != nil {
		return nil, err
	}
	return response.(getPaymentsResponse).payments, nil
}

func (c *client) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	response, err := c.getAccountEndpoint(ctx, getAccountRequest{id: id})
	if err != nil {
		return nil, err
	}
	return response.(getAccountResponse).account, nil
}
