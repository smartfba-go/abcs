package failurestwilio

import (
	"errors"

	twilio "github.com/twilio/twilio-go/client"
	"go.smartfba.io/abcs/failures"
)

func ErrorFromTwilioError(err error) error {
	if err == nil {
		return nil
	}

	if s, ok := StatusFromTwilioError(err); ok {
		return &errorStatus{
			Status: s,
			Err:    err,
		}
	}

	return err
}

func StatusFromTwilioError(err error) (*failures.Status, bool) {
	if err == nil {
		return nil, true
	}

	var restErr *twilio.TwilioRestError

	if errors.As(err, &restErr) {
		switch restErr.Code {
		case 400:
			return failures.NewStatusf(failures.InvalidArgument, "%w", err), true
		case 403, 20003:
			return failures.NewStatusf(failures.PermissionDenied, "%w", err), true
		case 404:
			return failures.NewStatusf(failures.NotFound, "%w", err), true
		case 410:
			return failures.NewStatusf(failures.Unknown, "%w", err), true
		case 503:
			return failures.NewStatusf(failures.Internal, "%w", err), true
		case 10001, 20005:
			return failures.NewStatusf(failures.FailedPrecondition, "account is not active: %w", err), true
		default:
			return failures.NewStatusf(failures.Unknown, "%w", err), true
		}
	}

	return nil, false
}

type errorStatus struct {
	Err    error
	Status *failures.Status
}

func (e *errorStatus) Error() string { return e.Status.Message }

func (e *errorStatus) Unwrap() error { return e.Err }

func (e *errorStatus) Status_() *failures.Status { return e.Status }

func (e *errorStatus) String() string { return e.Error() }
