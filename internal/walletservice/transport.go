package walletservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shkov/wallet-service/internal/account"
)

type applyPaymentRequest struct {
	input *account.Payment
}

type applyPaymentResponse struct {
}

func encodeApplyPaymentRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := request.(applyPaymentRequest)
	r.URL.Path = "/api/v1/payments"
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req.input); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeApplyPaymentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	input := &account.Payment{}
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		return nil, errBadRequest("failed to decode json request: %w", err)
	}
	return applyPaymentRequest{input: input}, nil
}

func encodeApplyPaymentResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func decodeApplyPaymentResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return nil, decodeError(r)
	}
	return applyPaymentResponse{}, nil
}

type getAccountRequest struct {
	id int64
}

type getAccountResponse struct {
	account *account.Account
}

func encodeGetAccountRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := request.(getAccountRequest)
	r.URL.Path = "/api/v1/accounts/" + strconv.FormatInt(req.id, 10)
	return nil
}

func decodeGetAccountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return nil, errBadRequest("failed to parse account id: %w", err)
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

type getPaymentsRequest struct {
	accountID int64
}

type getPaymentsResponse struct {
	payments []*account.Payment
}

func encodeGetPaymentsRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := request.(getPaymentsRequest)
	r.URL.Path = "/api/v1/payments/" + strconv.FormatInt(req.accountID, 10)
	return nil
}

func decodeGetPaymentsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	accountID, err := strconv.ParseInt(mux.Vars(r)["account_id"], 10, 64)
	if err != nil {
		return nil, errBadRequest("failed to parse account id: %w", err)
	}
	return getPaymentsRequest{accountID: accountID}, nil
}

func encodeGetPaymentsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(getPaymentsResponse)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp.payments); err != nil {
		return errInternal("failed to encode json response: %w", err)
	}
	return nil
}

func decodeGetPaymentsResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, decodeError(r)
	}
	resp := getPaymentsResponse{}
	if err := json.NewDecoder(r.Body).Decode(&resp.payments); err != nil {
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
