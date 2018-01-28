package gofaker

import (
	"fmt"
	"strings"
)

type FailTrigger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type FailTriggerTest struct {
	FailedAsExpected bool
	FailMessage      string
}

type ErrTestFailed struct {
	Msg string
}

func (ft *FailTriggerTest) Fatal(args ...interface{}) {
	ft.FailedAsExpected = true
	ft.FailMessage = strings.TrimRight(fmt.Sprintln(args...), "\r\n")
}

func (ft *FailTriggerTest) Fatalf(format string, args ...interface{}) {
	ft.FailedAsExpected = true
	ft.FailMessage = fmt.Sprintf(format, args...)
}

func (ft *FailTriggerTest) Reset() {
	ft.FailedAsExpected = false
	ft.FailMessage = ""
}

func (err *ErrTestFailed) Error() string {
	return err.Msg
}
