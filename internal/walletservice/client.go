package walletservice

import (
	"context"
	"errors"
	"time"

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
}

func (c *client) ApplyPayment(ctx context.Context, p *account.Payment) error {
	return nil
}

func (c *client) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	return nil, nil
}

func (c *client) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	return nil, nil
}

// NewClient creates a new client.
func NewClient(cfg ClientConfig) (Service, error) {
	err := cfg.validate()
	if err != nil {
		return nil, err
	}

	c := &client{}

	return c, nil
}
