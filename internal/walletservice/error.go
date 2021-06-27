package walletservice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type serviceError struct {
	code    int
	Message string `json:"error"`
}

// Error returns a string representation of the error.
func (e *serviceError) Error() string {
	return fmt.Sprintf("status %d: %s", e.code, e.Message)
}

// Encode encodes the error using the given HTTP response writer.
func (e *serviceError) Encode(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.code)
	json.NewEncoder(w).Encode(e)
}

// Decode decodes the error from the given HTTP response.
func (e *serviceError) Decode(r *http.Response) {
	e.code = r.StatusCode
	var res struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(r.Body).Decode(&res); err == nil && res.Error != "" {
		e.Message = res.Error
	} else {
		e.Message = http.StatusText(r.StatusCode)
	}
}

// ErrBadRequest creates a BadRequest service error.
func errBadRequest(format string, v ...interface{}) error {
	return &serviceError{
		code:    http.StatusBadRequest,
		Message: fmt.Sprintf(format, v...),
	}
}

// ErrNotFound creates a NotFound service error.
func errNotFound(format string, v ...interface{}) error {
	return &serviceError{
		code:    http.StatusNotFound,
		Message: fmt.Sprintf(format, v...),
	}
}

// ErrInternal creates an Internal service error.
func errInternal(format string, v ...interface{}) error {
	return &serviceError{
		code:    http.StatusInternalServerError,
		Message: fmt.Sprintf(format, v...),
	}
}
