package failures

import "fmt"

type Status struct {
	Code    Code
	Message string
}

func NewStatus(code Code, text string) *Status {
	return &Status{Code: code, Message: text}
}

func NewStatusf(code Code, format string, a ...any) *Status {
	return &Status{Code: code, Message: fmt.Sprintf(format, a...)}
}
