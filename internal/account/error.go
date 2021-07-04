package account

import (
	"errors"
)

// domain-level errors.
var (
	ErrNotFound                  = errors.New("account is not found")
	ErrMismatchPayment           = errors.New("mismatch of payment to account")
	ErrNotEnoughFunds            = errors.New("not enough funds in account")
	ErrNotPositiveAmount         = errors.New("payment amount is not positive")
	ErrAccountFromMustBePositive = errors.New("account from must be positive")
	ErrAccountToMustBePositive   = errors.New("account to must be positive")
	ErrMustBePositive            = errors.New("must be positive")
	ErrFromAndToMustBeDifferent  = errors.New("from and to must be different")
)
