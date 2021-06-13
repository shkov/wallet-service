package walletservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shkov/wallet-service/internal/account"
)

type getAccountRequest struct {
	id int64
}

type getAccountResponse struct {
	account *account.Account
}

func encodeGetAccountRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := request.(getAccountRequest)
	r.URL.Path = "/api/v1/account/s" + strconv.FormatInt(req.id, 10)
	return nil
}

func decodeGetAccountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return nil, errBadRequest("failed to parse product id: %w", err)
	}
	return getAccountRequest{id: id}, nil
}

func encodeGetAccountResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(getAccountResponse)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.account); err != nil {
		return errInternal("failed to encode json response: %w", err)
	}
	return nil
}

func decodeGetAccountResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, decodeError(r)
	}
	resp := getAccountResponse{}
	if err := json.NewDecoder(r.Body).Decode(&resp.account); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	return resp, nil
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	e, ok := err.(*serviceError)
	if !ok {
		e = &serviceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	e.Encode(w)
}

func decodeError(r *http.Response) error {
	e := &serviceError{}
	e.Decode(r)
	return e
}
