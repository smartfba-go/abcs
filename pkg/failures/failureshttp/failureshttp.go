package failureshttp

import (
	"encoding/json"
	"net/http"

	"go.smartfba.io/abcs/pkg/failures"
)

type Status struct {
	Code    failures.Code `json:"code"`
	Message string        `json:"message"`
}

func WriteErrorJSON(w http.ResponseWriter, err error) {
	status, _ := failures.StatusFromError(err)

	data, err := json.Marshal(Status{
		Code:    status.Code,
		Message: status.Message,
	})
	if err != nil {
		panic(err)
	}

	statusCode := StatusCode(status.Code)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func StatusCode(code failures.Code) int {
	switch code {
	case failures.OK:
		return http.StatusOK
	case failures.Canceled:
		return 499 // StatusClientClosedRequest
	case failures.Unknown:
		return http.StatusInternalServerError
	case failures.InvalidArgument:
		return http.StatusBadRequest
	case failures.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case failures.NotFound:
		return http.StatusNotFound
	case failures.AlreadyExists:
		return http.StatusConflict
	case failures.PermissionDenied:
		return http.StatusForbidden
	case failures.ResourceExhausted:
		return http.StatusTooManyRequests
	case failures.FailedPrecondition:
		return http.StatusBadRequest
	case failures.Aborted:
		return http.StatusConflict
	case failures.OutOfRange:
		return http.StatusBadRequest
	case failures.Unimplemented:
		return http.StatusNotImplemented
	case failures.Internal:
		return http.StatusInternalServerError
	case failures.Unavailable:
		return http.StatusServiceUnavailable
	case failures.DataLoss:
		return http.StatusInternalServerError
	case failures.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
