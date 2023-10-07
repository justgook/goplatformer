package util

import (
	"fmt"
	"runtime"
	"strings"
)

type TraceError struct {
	Source   error
	Package  string
	Function string
	Code     ErrorCode
	Line     int
}

func Catch(err error) *TraceError {
	if err == nil {
		return nil
	}

	p, fn, line := getErrorContextInfo()

	return &TraceError{
		Package:  p,
		Function: fn,
		Line:     line,
		Code:     ErrorCodeWrapper,
		Source:   err,
	}
}

type ErrorCode int

const (
	ErrorCodeUnauthorized        ErrorCode = iota // 401
	ErrorCodeTimeOut                              // 408
	ErrorCodeConflict                             // 409  - user already exists
	ErrorCodeUnprocessableEntity                  // 422 - failed to process input
	ErrorCodeFailedDependency                     // 424 - record is null in state
	ErrorCodeNotFoundInState
	ErrorCodeNotFoundInBase
	ErrorCodeNotFound
	ErrorCodeWrapper
)

func (e *TraceError) Unwrap() error {
	return e.Source
}
func (e *TraceError) Error() string {
	if e.Source != nil {
		return fmt.Sprintf(
			"(1)%s.%s[%d]: %s",
			e.Package,
			e.Function,
			e.Line,
			e.Source.Error(),
		)
	}

	return fmt.Sprintf(
		"(2)%s.%s[%d]: %s",
		e.Package,
		e.Function,
		e.Line,
		e.Code.codeText(),
	)
}

// codeText returns a text for the status code. It returns the empty
// string if the code is unknown.
func (c ErrorCode) codeText() string {
	return map[ErrorCode]string{
		ErrorCodeUnauthorized:        "Unauthorized",
		ErrorCodeTimeOut:             "Request Timeout",
		ErrorCodeConflict:            "Conflict",
		ErrorCodeUnprocessableEntity: "Unprocessable Entity",
		ErrorCodeFailedDependency:    "Failed Dependency",
		ErrorCodeNotFoundInState:     "Not Found In State",
		ErrorCodeNotFoundInBase:      "Not Found In Base",
		ErrorCodeWrapper:             "Something went wrong",
	}[c]
}
func getErrorContextInfo() (string, string, int) {
	const (
		sep  = "."
		skip = 2
	)

	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", "unknown", 0
	}

	parts := strings.Split(runtime.FuncForPC(pc).Name(), sep)
	pl := len(parts)
	pName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + sep + funcName
		pName = strings.Join(parts[0:pl-2], sep)
	} else {
		pName = strings.Join(parts[0:pl-1], sep)
	}

	return pName, funcName, line
}

