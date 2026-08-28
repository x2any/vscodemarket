package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorCode is a stable identifier for a client-facing error.
type ErrorCode string

const (
	CodeVersionNotFound       ErrorCode = "VERSION_NOT_FOUND"
	CodeExtensionVersionNF    ErrorCode = "EXTENSION_VERSION_NOT_FOUND"
	CodeInvalidPlatformArch   ErrorCode = "INVALID_PLATFORM_ARCH"
	CodeUpstreamFailure       ErrorCode = "UPSTREAM_FAILURE"
	CodeInvalidRequest        ErrorCode = "INVALID_REQUEST"
	CodeInternal              ErrorCode = "INTERNAL"
	CodeNonOfficialURLBlocked ErrorCode = "NON_OFFICIAL_URL_BLOCKED"
)

// APIError carries bilingual messages for a single error response.
type APIError struct {
	Code      ErrorCode `json:"code"`
	MessageZH string    `json:"message_zh"`
	MessageEN string    `json:"message_en"`
}

// ErrorBody is the wire envelope.
type ErrorBody struct {
	Error APIError `json:"error"`
}

// WriteError emits a uniform error JSON body.
func WriteError(w http.ResponseWriter, status int, e APIError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorBody{Error: e})
}

// Err is a convenience constructor.
func Err(code ErrorCode, zh, en string) APIError {
	return APIError{Code: code, MessageZH: zh, MessageEN: en}
}
