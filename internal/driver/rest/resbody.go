package rest

import (
	"net/http"
	"time"
)

func ResponseErrorUnauthorized(msg string) error {
	return &httpError{
		StatusCode: http.StatusUnauthorized,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_UNAUTHORIZED",
	}
}

func ResponseErrorBadRequest(msg string) error {
	return &httpError{
		StatusCode: http.StatusBadRequest,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_BAD_REQUEST",
	}
}

func ResponseErrorReadTimeout(msg string) error {
	return &httpError{
		StatusCode: http.StatusRequestTimeout,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_READ_TIMEOUT",
	}
}

func ResponseErrorInvalidAccessToken(msg string) error {
	return &httpError{
		StatusCode: http.StatusUnauthorized,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_INVALID_ACCESS_TOKEN",
	}
}

func ResponseErrorForbiddenAccess(msg string) error {
	return &httpError{
		StatusCode: http.StatusForbidden,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_FORBIDDEN_ACCESS",
	}
}

func ResponseErrorNotFound(msg string) error {
	return &httpError{
		StatusCode: http.StatusNotFound,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_NOT_FOUND",
	}
}

func ResponseErrorInternalServerError(msg string) error {
	return &httpError{
		StatusCode: http.StatusInternalServerError,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_INTERNAL_SERVER_ERROR",
	}
}

type httpError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"msg"`
	Err        string `json:"err"`
	Ts         int64  `json:"ts"`
}

func (e *httpError) Error() string {
	return e.Message
}
