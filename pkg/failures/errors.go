package failures

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrCanceled           = New(Canceled, "canceled")
	ErrUnknown            = New(Unknown, "unknown")
	ErrInvalidArgument    = New(InvalidArgument, "invalid argument")
	ErrDeadlineExceeded   = New(DeadlineExceeded, "deadline exceeded")
	ErrNotFound           = New(NotFound, "not found")
	ErrAlreadyExists      = New(AlreadyExists, "already exists")
	ErrPermissionDenied   = New(PermissionDenied, "permission denied")
	ErrResourceExhausted  = New(ResourceExhausted, "resource exhausted")
	ErrFailedPrecondition = New(FailedPrecondition, "failed precondition")
	ErrAborted            = New(Aborted, "aborted")
	ErrOutOfRange         = New(OutOfRange, "out of range")
	ErrUnimplemented      = New(Unimplemented, "unimplemented")
	ErrInternal           = New(Internal, "internal")
	ErrUnavailable        = New(Unavailable, "unavailable")
	ErrDataLoss           = New(DataLoss, "data loss")
	ErrUnauthenticated    = New(Unauthenticated, "unauthenticated")
)

func New(code Code, text string) error {
	return &errorStatus{
		Status: NewStatus(code, text),
		Err:    nil,
	}
}

func Newf(code Code, format string, a ...any) error {
	err := fmt.Errorf(format, a...)

	return &errorStatus{
		Status: NewStatus(code, err.Error()),
		Err:    errors.Unwrap(err),
	}
}

func Wrap(code Code, text string, err error) error {
	return &errorStatus{
		Status: NewStatus(code, text),
		Err:    err,
	}
}

type errorStatus struct {
	Err    error
	Status *Status
}

func (e *errorStatus) Error() string { return e.Status.Message }

func (e *errorStatus) Unwrap() error { return e.Err }

func (e *errorStatus) Status_() *Status { return e.Status }

func (e *errorStatus) String() string { return e.Error() }

func StatusFromError(err error) (*Status, bool) {
	if err == nil {
		return nil, true
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return NewStatus(DeadlineExceeded, err.Error()), true
	}

	if errors.Is(err, context.Canceled) {
		return NewStatus(Canceled, err.Error()), true
	}

	if se, ok := err.(interface{ Status_() *Status }); ok {
		return se.Status_(), true
	}

	return NewStatus(Unknown, err.Error()), false
}

func ConvertError(err error) *Status {
	s, _ := StatusFromError(err)

	return s
}

func CodeFromError(err error) Code {
	if err == nil {
		return OK
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return DeadlineExceeded
	}

	if errors.Is(err, context.Canceled) {
		return Canceled
	}

	if se, ok := err.(interface{ Status_() *Status }); ok {
		return se.Status_().Code
	}

	return Unknown
}
