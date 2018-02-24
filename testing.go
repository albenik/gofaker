package gofaker

import (
	"fmt"
	"strings"
)

type FailTrigger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type TestTrigger struct {
	FailedAsExpected bool
	FailMessage      string
}

type ErrTestFailed struct {
	Msg string
}

func (ft *TestTrigger) Fatal(args ...interface{}) {
	ft.FailedAsExpected = true
	ft.FailMessage = strings.TrimRight(fmt.Sprintln(args...), "\r\n")
}

func (ft *TestTrigger) Fatalf(format string, args ...interface{}) {
	ft.FailedAsExpected = true
	ft.FailMessage = fmt.Sprintf(format, args...)
}

func (ft *TestTrigger) Reset() {
	ft.FailedAsExpected = false
	ft.FailMessage = ""
}

func (err *ErrTestFailed) Error() string {
	return err.Msg
}
