package failures

import (
	"errors"

	"go.temporal.io/sdk/temporal"
)

var (
	temporalApplicationErrorType = "StatusError"
)

func TemporalError(err error) error {
	if err == nil {
		return nil
	}

	if se, ok := err.(interface{ Status_() *Status }); ok {
		s := se.Status_()

		if s.Code.CanRetry() {
			return temporal.NewApplicationErrorWithCause(s.Message, temporalApplicationErrorType, err, s)
		}

		return temporal.NewNonRetryableApplicationError(s.Message, temporalApplicationErrorType, err, s)
	}

	return err
}

func ErrorFromTemporalError(err error) error {
	if err == nil {
		return nil
	}

	if s, ok := StatusFromTemporalError(err); ok {
		return &errorStatus{
			Status: s,
			Err:    err,
		}
	}

	return err
}

func StatusFromTemporalError(err error) (*Status, bool) {
	if err == nil {
		return nil, true
	}

	var (
		applicationError *temporal.ApplicationError
		canceledError    *temporal.CanceledError
		timeoutError     *temporal.TimeoutError
		panicError       *temporal.PanicError
	)

	if errors.As(err, &applicationError) {
		if applicationError.Type() == temporalApplicationErrorType {
			var s Status

			if applicationError.Details(&s) == nil {
				return &s, true
			}
		}
	}

	if errors.As(err, &canceledError) {
		if !canceledError.HasDetails() {
			return NewStatus(Canceled, canceledError.Error()), true
		}

		var s Status

		if canceledError.Details(&s) == nil {
			return &s, true
		}
	}

	if errors.As(err, &timeoutError) {
		return NewStatus(DeadlineExceeded, timeoutError.Error()), true
	}

	if errors.As(err, &panicError) {
		return NewStatus(Internal, panicError.Error()), true
	}

	return nil, false
}
