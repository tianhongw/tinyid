package errdef

import (
	"net/http"
)

var (
	HttpErrorInternal     = NewHttpError("internal server error", HttpErrorCodeUniversal)
	HttpErrorNotFound     = NewHttpError("not found", HttpErrorCodeNotFound)
	HttpErrorBadRequest   = NewHttpError("bad request", HttpErrorCodeBadRequest)
	HttpErrorUnauthorized = NewHttpError("unauthorized", HttpErrorCodeUnauthorized)
	HttpErrorForbidden    = NewHttpError("forbidden", HttpErrorCodeForbidden)
)

type HttpErrorCode int

const (
	HttpErrorCodeUniversal          = HttpErrorCode(-1)
	HttpErrorCodeUserBlocked        = HttpErrorCode(-11)
	HttpErrorCodeInvalidAccountType = HttpErrorCode(-12)
	HttpErrorCodeBadRequest         = HttpErrorCode(-400)
	HttpErrorCodeUnauthorized       = HttpErrorCode(-401)
	HttpErrorCodeForbidden          = HttpErrorCode(-403)
	HttpErrorCodeNotFound           = HttpErrorCode(-404)
)

func (c HttpErrorCode) StatusCode() int {
	switch c {
	case HttpErrorCodeBadRequest:
		return http.StatusBadRequest
	case HttpErrorCodeNotFound:
		return http.StatusNotFound
	case HttpErrorCodeUnauthorized, HttpErrorCodeUserBlocked:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

type HttpError struct {
	Code    HttpErrorCode `json:"code"`
	Message string        `json:"message"`
	Details []string      `json:"details,omitempty"`
}

func (h HttpError) Error() string {
	if h.Message == "" {
		return http.StatusText(http.StatusInternalServerError)
	}
	return h.Message
}

func NewHttpError(message string, code ...HttpErrorCode) HttpError {
	return NewHttpErrorWithDetails(message, nil, code...)
}

func NewHttpErrorWithDetails(message string, details []string, code ...HttpErrorCode) HttpError {
	var c HttpErrorCode
	if len(code) > 0 {
		c = code[0]
	} else {
		c = HttpErrorCodeUniversal
	}

	return HttpError{
		Code:    c,
		Message: message,
		Details: details,
	}
}
