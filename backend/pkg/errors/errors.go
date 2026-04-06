// Package errors defines application-level error types and connect-RPC mapping.
package errors

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
)

// AppError is a typed application error with a stable Code and a human-readable Message.
type AppError struct {
	Code    string
	Message string
}

// Error implements the error interface. Format: "CODE: message".
func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Is reports whether target has the same Code as e.
// This enables errors.Is(someWrappedErr, ErrNotFound) to work.
func (e AppError) Is(target error) bool {
	var t AppError
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}

// Wrap returns a new AppError with the same Code but a different Message.
func (e AppError) Wrap(msg string) AppError {
	return AppError{Code: e.Code, Message: msg}
}

// Sentinel errors — compare with errors.Is.
var (
	ErrNotFound    = AppError{Code: "NOT_FOUND", Message: "not found"}
	ErrUnauthorized = AppError{Code: "UNAUTHORIZED", Message: "unauthorized"}
	ErrForbidden   = AppError{Code: "FORBIDDEN", Message: "forbidden"}
	ErrValidation  = AppError{Code: "VALIDATION", Message: "validation error"}
	ErrConflict    = AppError{Code: "CONFLICT", Message: "conflict"}
	ErrInternal    = AppError{Code: "INTERNAL", Message: "internal error"}
	ErrBadRequest  = AppError{Code: "BAD_REQUEST", Message: "bad request"}
	ErrUnavailable = AppError{Code: "UNAVAILABLE", Message: "service unavailable"}
)

// codeToConnect maps AppError.Code values to connect-RPC codes.
var codeToConnect = map[string]connect.Code{
	"NOT_FOUND":    connect.CodeNotFound,
	"UNAUTHORIZED": connect.CodeUnauthenticated,
	"FORBIDDEN":    connect.CodePermissionDenied,
	"VALIDATION":   connect.CodeInvalidArgument,
	"CONFLICT":     connect.CodeAlreadyExists,
	"INTERNAL":     connect.CodeInternal,
	"BAD_REQUEST":  connect.CodeInvalidArgument,
	"UNAVAILABLE":  connect.CodeUnavailable,
}

// ToConnectError converts any error to a *connect.Error.
// AppErrors are mapped by their Code. Unknown / generic errors become
// CodeInternal with the message "internal error" — details are never exposed.
func ToConnectError(err error) error {
	var appErr AppError
	if errors.As(err, &appErr) {
		code, ok := codeToConnect[appErr.Code]
		if !ok {
			code = connect.CodeInternal
		}
		return connect.NewError(code, errors.New(appErr.Message))
	}

	// Generic error — hide details.
	return connect.NewError(connect.CodeInternal, errors.New("internal error"))
}
