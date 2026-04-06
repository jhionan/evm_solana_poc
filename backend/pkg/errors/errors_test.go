package errors_test

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
	apperrors "github.com/jhionan/multichain-staking/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAppError_ErrorString(t *testing.T) {
	err := apperrors.AppError{Code: "NOT_FOUND", Message: "record not found"}
	assert.Equal(t, "NOT_FOUND: record not found", err.Error())
}

func TestAppError_Is_SameCode(t *testing.T) {
	err := apperrors.AppError{Code: "NOT_FOUND", Message: "record not found"}
	assert.True(t, errors.Is(err, apperrors.ErrNotFound))
}

func TestAppError_Is_DifferentCode(t *testing.T) {
	err := apperrors.AppError{Code: "NOT_FOUND", Message: "record not found"}
	assert.False(t, errors.Is(err, apperrors.ErrUnauthorized))
}

func TestAppError_Wrap(t *testing.T) {
	wrapped := apperrors.ErrNotFound.Wrap("user not found")
	assert.Equal(t, "NOT_FOUND", wrapped.Code)
	assert.Equal(t, "user not found", wrapped.Message)
	assert.True(t, errors.Is(wrapped, apperrors.ErrNotFound))
}

func TestToConnectError_NotFound(t *testing.T) {
	err := apperrors.ErrNotFound.Wrap("user not found")
	connectErr := apperrors.ToConnectError(err)
	var connErr *connect.Error
	assert.True(t, errors.As(connectErr, &connErr))
	assert.Equal(t, connect.CodeNotFound, connErr.Code())
}

func TestToConnectError_Validation(t *testing.T) {
	err := apperrors.ErrValidation.Wrap("field is required")
	connectErr := apperrors.ToConnectError(err)
	var connErr *connect.Error
	assert.True(t, errors.As(connectErr, &connErr))
	assert.Equal(t, connect.CodeInvalidArgument, connErr.Code())
}

func TestToConnectError_GenericError(t *testing.T) {
	err := errors.New("something internal happened")
	connectErr := apperrors.ToConnectError(err)
	var connErr *connect.Error
	assert.True(t, errors.As(connectErr, &connErr))
	assert.Equal(t, connect.CodeInternal, connErr.Code())
	assert.Equal(t, "internal error", connErr.Message())
}

func TestToConnectError_AllSentinels(t *testing.T) {
	cases := []struct {
		err      apperrors.AppError
		wantCode connect.Code
	}{
		{apperrors.ErrNotFound, connect.CodeNotFound},
		{apperrors.ErrUnauthorized, connect.CodeUnauthenticated},
		{apperrors.ErrForbidden, connect.CodePermissionDenied},
		{apperrors.ErrValidation, connect.CodeInvalidArgument},
		{apperrors.ErrConflict, connect.CodeAlreadyExists},
		{apperrors.ErrInternal, connect.CodeInternal},
		{apperrors.ErrBadRequest, connect.CodeInvalidArgument},
		{apperrors.ErrUnavailable, connect.CodeUnavailable},
	}
	for _, tc := range cases {
		t.Run(tc.err.Code, func(t *testing.T) {
			connectErr := apperrors.ToConnectError(tc.err)
			var connErr *connect.Error
			assert.True(t, errors.As(connectErr, &connErr))
			assert.Equal(t, tc.wantCode, connErr.Code())
		})
	}
}
