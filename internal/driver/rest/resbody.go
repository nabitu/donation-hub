package rest

import (
	"encoding/json"
	"net/http"
	"time"
)

type internalServerError struct {
	Status     int64  `json:"status"`
	StatusCode int    `json:"-"`
	Message    string `json:"msg"`
	Err        string `json:"err"`
	Ts         int64  `json:"ts"`
}
type httpError struct {
	Ok         bool   `json:"ok"`
	StatusCode int    `json:"-"`
	Message    string `json:"msg"`
	Err        string `json:"err"`
	Ts         int64  `json:"ts"`
}

type httpSuccess struct {
	Ok   bool        `json:"ok"`
	Data interface{} `json:"data,omitempty"`
	Ts   int64       `json:"ts"`
}

func ResponseSuccess(w http.ResponseWriter, data interface{}) {
	res := &httpSuccess{
		Ok:   true,
		Data: data,
		Ts:   time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func ResponseErrorUnauthorized(w http.ResponseWriter, msg string) {
	res := &httpError{
		Ok:         false,
		StatusCode: http.StatusUnauthorized,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_UNAUTHORIZED",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(res)
}

func ResponseErrorBadRequest(w http.ResponseWriter, msg string) {
	res := &httpError{
		Ok:         false,
		StatusCode: http.StatusBadRequest,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_BAD_REQUEST",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(res)
}

func ResponseErrorReadTimeout(w http.ResponseWriter, msg string) {
	res := &httpError{
		Ok:         false,
		StatusCode: http.StatusRequestTimeout,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_READ_TIMEOUT",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusRequestTimeout)
	json.NewEncoder(w).Encode(res)
}

func ResponseErrorInvalidAccessToken(w http.ResponseWriter, msg string) {
	res := &httpError{
		Ok:         false,
		StatusCode: http.StatusUnauthorized,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_INVALID_ACCESS_TOKEN",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(res)
}

func ResponseErrorForbiddenAccess(w http.ResponseWriter, msg string) {
	res := &httpError{
		Ok:         false,
		StatusCode: http.StatusForbidden,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_FORBIDDEN_ACCESS",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(res)
}

func ResponseErrorNotFound(w http.ResponseWriter, msg string) {
	res := &httpError{
		Ok:         false,
		StatusCode: http.StatusNotFound,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_NOT_FOUND",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(res)
}

func ResponseErrorInternalServerError(w http.ResponseWriter, msg string) {
	res := &internalServerError{
		Status:     http.StatusInternalServerError,
		StatusCode: http.StatusInternalServerError,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_INTERNAL_SERVER_ERROR",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(res)
}

func ResponseErrorInvalidCredentials(w http.ResponseWriter, msg string) {
	res := &httpError{
		Ok:         false,
		StatusCode: http.StatusUnauthorized,
		Message:    msg,
		Ts:         time.Now().Unix(),
		Err:        "ERR_INVALID_CREDENTIALS",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(res)
}
