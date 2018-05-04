package gofaker

type Fatality interface {
	Fatalf(format string, args ...interface{})
}

type AssertionFailedError struct {
	Message string
}

func (cf *AssertionFailedError) Error() string {
	return cf.Message
}
